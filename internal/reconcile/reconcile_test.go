package reconcile

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

func TestBuildPlanPreviewAndApplyInitializesTransactionally(t *testing.T) {
	root := t.TempDir()
	observed, err := state.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	graph, err := state.CompileGraph("standard")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 7, 9, 20, 0, 0, 0, time.UTC)
	request := Request{
		Command:       "init",
		Root:          root,
		Observed:      observed,
		Graph:         graph,
		DesiredLock:   testLock(graph, nil, now, "init"),
		DesiredConfig: testConfig(root),
		EnsureIgnore:  true,
	}
	plan, err := BuildPlan(request)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Actions) == 0 || len(plan.Blockers) != 0 {
		t.Fatalf("plan actions=%d blockers=%#v", len(plan.Actions), plan.Blockers)
	}
	preview := Preview(plan)
	if preview.Status != StatusPlanned || !preview.DryRun {
		t.Fatalf("preview = %#v", preview)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.TransactionPath))); !os.IsNotExist(err) {
		t.Fatalf("preview created transaction marker: %v", err)
	}
	result, err := Apply(plan, ApplyOptions{Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusSucceeded || result.StateAfter != string(state.StateCurrent) {
		t.Fatalf("apply = %#v", result)
	}
	for _, path := range []string{state.TransactionPath, ".codeheart/local/kit-transactions"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(path))); !os.IsNotExist(err) {
			t.Fatalf("successful apply retained %s: %v", path, err)
		}
	}
	final, err := state.Inspect(root)
	if err != nil || final.Classification != state.StateCurrent {
		t.Fatalf("final = %#v, err=%v", final, err)
	}

	noOpLock := testLock(graph, final.Lock, now.Add(time.Minute), "sync")
	noOpPlan, err := BuildPlan(Request{Command: "sync", Root: root, Observed: final, Graph: graph, DesiredLock: noOpLock, EnsureIgnore: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(noOpPlan.Actions) != 0 {
		t.Fatalf("idempotent plan has actions: %#v", noOpPlan.Actions)
	}
}

func TestOperationResultJSONGolden(t *testing.T) {
	result := NewResult("repair")
	result.Status = StatusBlocked
	result.DryRun = true
	result.StateBefore = string(state.StateDrifted)
	result.Changes = append(result.Changes, Change{Action: "replace", Path: ".codeheart/kit/example.md", Owner: "managed"})
	result.Blockers = append(result.Blockers, Blocker{Code: "managed_path_modified", Message: "retired managed path contains changes", Path: ".codeheart/kit/retired.md", Remediation: "preserve the file", RetryCommand: "repair"})
	result.Validations = append(result.Validations, Validation{Name: "change-plan", Status: "passed"})
	var output bytes.Buffer
	if err := WriteJSON(&output, result); err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join("testdata", "operation-result.json"))
	if err != nil {
		t.Fatal(err)
	}
	if output.String() != string(expected) {
		t.Fatalf("operation result changed\nwant:\n%s\ngot:\n%s", expected, output.String())
	}
}

func TestRecoverStaleTransactionRequiresDeadVerifiedPreCommitIdentity(t *testing.T) {
	root := t.TempDir()
	canonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	transactionID := "stale-transaction"
	recoveryPath := filepath.ToSlash(filepath.Join(".codeheart", "local", "kit-transactions", transactionID))
	marker := transactionMarker{
		SchemaVersion: 1,
		TransactionID: transactionID,
		Command:       "repair",
		Phase:         "staged",
		PID:           2_000_000_000,
		TargetRoot:    canonical,
		RecoveryPath:  recoveryPath,
	}
	markerData, _ := json.Marshal(marker)
	markerPath := filepath.Join(root, filepath.FromSlash(state.TransactionPath))
	if err := os.MkdirAll(filepath.Dir(markerPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(markerPath, markerData, 0o600); err != nil {
		t.Fatal(err)
	}
	recoveryRoot := filepath.Join(root, filepath.FromSlash(recoveryPath))
	if err := os.MkdirAll(recoveryRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	planData, _ := json.Marshal(Plan{ID: transactionID, Command: "repair", Root: canonical})
	if err := os.WriteFile(filepath.Join(recoveryRoot, "plan.json"), planData, 0o600); err != nil {
		t.Fatal(err)
	}
	if recoverable, err := RecoverStaleTransaction(root, true); err != nil || !recoverable {
		t.Fatalf("dry recovery = %v, %v", recoverable, err)
	}
	if _, err := os.Stat(markerPath); err != nil {
		t.Fatalf("dry recovery changed marker: %v", err)
	}
	if recovered, err := RecoverStaleTransaction(root, false); err != nil || !recovered {
		t.Fatalf("recovery = %v, %v", recovered, err)
	}
	if _, err := os.Stat(markerPath); !os.IsNotExist(err) {
		t.Fatalf("recovery retained marker: %v", err)
	}
}

func TestApplyDetectsParentIdentityReplacement(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Join(root, "nested")
	if err := os.Mkdir(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            "parent-identity",
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		ExpectedAfter: []state.Classification{state.StateCurrent},
		Actions:       []Action{{Kind: "create", Target: "nested/file.txt", Owner: "managed", Content: []byte("x\n")}},
	}
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "validated" {
			return nil
		}
		if err := os.Rename(parent, parent+"-old"); err != nil {
			return err
		}
		return os.Mkdir(parent, 0o755)
	}})
	if err != nil || result.Status != StatusRolledBack || !result.Rollback.Succeeded {
		t.Fatalf("identity replacement result = %#v, err=%v", result, err)
	}
}

func TestFailedRollbackRetainsRecoveryState(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("before\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            "rollback-failure",
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions: []Action{
			{Kind: "replace", Target: "a.txt", Owner: "managed", Content: []byte("after-a\n")},
			{Kind: "replace", Target: "b.txt", Owner: "managed", Content: []byte("after-b\n")},
		},
	}
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "commit:b.txt" {
			return nil
		}
		backup := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID, "backup", "a.txt")
		if err := os.Remove(backup); err != nil {
			return err
		}
		return errors.New("injected failure after backup loss")
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("failed rollback result = %#v, err=%v", result, err)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.TransactionPath))); err != nil {
		t.Fatalf("recovery marker missing: %v", err)
	}
}

func TestApplyRollsBackCommittedChangesAfterInjectedFailure(t *testing.T) {
	root := t.TempDir()
	sentinel := filepath.Join(root, "sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("before\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	observed, _ := state.Inspect(root)
	graph, _ := state.CompileGraph("standard")
	now := time.Date(2026, 7, 9, 20, 0, 0, 0, time.UTC)
	plan, err := BuildPlan(Request{
		Command:       "init",
		Root:          root,
		Observed:      observed,
		Graph:         graph,
		DesiredLock:   testLock(graph, nil, now, "init"),
		DesiredConfig: testConfig(root),
		EnsureIgnore:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	commits := 0
	result, err := Apply(plan, ApplyOptions{Now: now, Hook: func(phase string) error {
		if strings.HasPrefix(phase, "commit:") {
			commits++
			if commits == 3 {
				return errors.New("injected commit failure")
			}
		}
		return nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != StatusRolledBack || !result.Rollback.Succeeded {
		t.Fatalf("rollback result = %#v", result)
	}
	data, err := os.ReadFile(sentinel)
	if err != nil || string(data) != "before\n" {
		t.Fatalf("sentinel changed: %q err=%v", data, err)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.LockPath))); !os.IsNotExist(err) {
		t.Fatalf("rollback left lock: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.TransactionPath))); !os.IsNotExist(err) {
		t.Fatalf("rollback left marker: %v", err)
	}
}

func TestBuildPlanPreservesModifiedRetiredManagedPath(t *testing.T) {
	root := t.TempDir()
	retired := filepath.Join(root, ".codeheart", "kit", "retired.md")
	if err := os.MkdirAll(filepath.Dir(retired), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(retired, []byte("user change\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	graph, _ := state.CompileGraph("standard")
	observed := state.Observed{
		Root:           root,
		Classification: state.StateDrifted,
		Lock: map[string]any{
			"managed_paths": []any{map[string]any{
				"path":            ".codeheart/kit/retired.md",
				"ownership":       "managed",
				"checksum_sha256": strings.Repeat("a", 64),
			}},
		},
	}
	plan, err := BuildPlan(Request{Command: "sync", Root: root, Observed: observed, Graph: graph})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, blocker := range plan.Blockers {
		if blocker.Code == "managed_path_modified" && blocker.Path == ".codeheart/kit/retired.md" {
			found = true
		}
	}
	if !found {
		t.Fatalf("plan did not preserve modified retired path: %#v", plan.Blockers)
	}
}

func TestBuildPlanRemovesOnlyUnmodifiedRetiredManagedPath(t *testing.T) {
	root := t.TempDir()
	content := []byte("prior managed bytes\n")
	retired := filepath.Join(root, ".codeheart", "kit", "retired.md")
	if err := os.MkdirAll(filepath.Dir(retired), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(retired, content, 0o644); err != nil {
		t.Fatal(err)
	}
	graph, _ := state.CompileGraph("standard")
	observed := state.Observed{Classification: state.StateDrifted, Lock: map[string]any{
		"managed_paths": []any{map[string]any{
			"path":            ".codeheart/kit/retired.md",
			"ownership":       "managed",
			"checksum_sha256": digest(content),
		}},
	}}
	plan, err := BuildPlan(Request{Command: "sync", Root: root, Observed: observed, Graph: graph})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, action := range plan.Actions {
		if action.Kind == "remove" && action.Target == ".codeheart/kit/retired.md" {
			found = true
		}
	}
	if !found || len(plan.Blockers) != 0 {
		t.Fatalf("safe removal plan = %#v, blockers=%#v", plan.Actions, plan.Blockers)
	}
}

func TestStagedValidationFailureRollsBackBeforeCommit(t *testing.T) {
	root := t.TempDir()
	observed, _ := state.Inspect(root)
	graph, _ := state.CompileGraph("standard")
	now := time.Date(2026, 7, 9, 20, 0, 0, 0, time.UTC)
	plan, err := BuildPlan(Request{
		Command:       "init",
		Root:          root,
		Observed:      observed,
		Graph:         graph,
		DesiredLock:   testLock(graph, nil, now, "init"),
		DesiredConfig: testConfig(root),
		EnsureIgnore:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := Apply(plan, ApplyOptions{Now: now, Hook: func(phase string) error {
		if phase != "staged" {
			return nil
		}
		stagedLock := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID, "stage", filepath.FromSlash(state.LockPath))
		return os.WriteFile(stagedLock, []byte("schema_version: invalid\n"), 0o600)
	}})
	if err != nil || result.Status != StatusRolledBack || !result.Rollback.Succeeded {
		t.Fatalf("staged validation result = %#v, err=%v", result, err)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.LockPath))); !os.IsNotExist(err) {
		t.Fatalf("failed validation committed lock: %v", err)
	}
}

func TestApplyRejectsSymlinkTraversalAndConcurrentMarker(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows reparse behavior is covered by the Windows validation job")
	}
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            "unsafe",
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions:       []Action{{Kind: "create", Target: "linked/file.txt", Owner: "managed", Content: []byte("x\n")}},
	}
	result, err := Apply(plan, ApplyOptions{})
	if err != nil || result.Status != StatusBlocked || len(result.Blockers) == 0 || result.Blockers[0].Code != "unsafe_target" {
		t.Fatalf("symlink result = %#v, err=%v", result, err)
	}
	if _, err := os.Stat(filepath.Join(outside, "file.txt")); !os.IsNotExist(err) {
		t.Fatalf("symlink escape wrote outside target: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(root, ".codeheart"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(state.TransactionPath)), []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan = Plan{SchemaVersion: 1, ID: "concurrent", Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{{Kind: "create", Target: "safe.txt", Owner: "managed", Content: []byte("x\n")}}}
	result, err = Apply(plan, ApplyOptions{})
	if err != nil || result.Status != StatusBlocked || result.Blockers[0].Code != "transaction_in_progress" {
		t.Fatalf("concurrent result = %#v, err=%v", result, err)
	}
}

func testLock(graph state.Graph, existing map[string]any, now time.Time, command string) map[string]any {
	update := map[string]any{
		"last_update_check_at":  now.UTC().Format(time.RFC3339),
		"next_update_check_due": now.Add(7 * 24 * time.Hour).UTC().Format(time.RFC3339),
		"latest_seen_version":   version.Version,
		"update_status":         "current",
	}
	native := map[string]any{}
	release := map[string]any{"asset_url": "local-source", "checksum_sha256": strings.Repeat("0", 64)}
	provenance := map[string]any{"verification_status": "local-source", "source": "local-source"}
	cliRepair := map[string]any{"installed_cli_path": "codeheart-operating-kit", "repair_source_url": "local-source", "repair_checksum_sha256": strings.Repeat("0", 64)}
	if existing != nil {
		if value := state.Map(existing["update_check"]); value != nil {
			update = value
		}
		if value := state.Map(existing["native_capabilities"]); value != nil {
			native = value
		}
		if value := state.Map(existing["release"]); value != nil {
			release = value
		}
		if value := state.Map(existing["release_provenance"]); value != nil {
			provenance = value
		}
		if value := state.Map(existing["cli_repair"]); value != nil {
			cliRepair = value
		}
	}
	return map[string]any{
		"schema_version":      2,
		"kit_version":         version.Version,
		"state_generation":    1,
		"selected_profile":    graph.ProfileID,
		"selected_components": stringsToAny(graph.SelectedComponents),
		"release":             release,
		"release_provenance":  provenance,
		"managed_paths":       ManagedPathRecords(graph),
		"managed_sections":    ManagedSectionRecords(graph),
		"generated_surfaces":  GeneratedSurfaceRecords(graph),
		"cli_repair":          cliRepair,
		"update_check":        update,
		"native_capabilities": native,
		"last_operation": map[string]any{
			"transaction_id":      "pending",
			"command":             command,
			"completed_at":        now.UTC().Format(time.RFC3339),
			"previous_generation": 0,
		},
	}
}

func testConfig(root string) map[string]any {
	return map[string]any{
		"schema_version":        1,
		"selected_profile":      "standard",
		"project_display_name":  "Example",
		"selected_setup_folder": root,
		"local_consumer_layer": map[string]any{
			"repo_docs_path":           "docs/repo/",
			"agent_memory_path":        "docs/agent-memory/",
			"user_layer_path":          ".codeheart/user/",
			"local_machine_layer_path": ".codeheart/local/",
		},
		"component_settings": map[string]any{},
	}
}

func stringsToAny(values []string) []any {
	result := make([]any, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}

func digest(data []byte) string {
	value := sha256.Sum256(data)
	return hex.EncodeToString(value[:])
}
