package reconcile

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
	transactionID := strings.Repeat("1", 32)
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
	planData, _ := json.Marshal(Plan{SchemaVersion: 1, ID: transactionID, Command: "repair", Root: canonical})
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

func TestStaleRecoveryPreservesTransactionWhenMarkerChangesAfterValidation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("marker symlink stale-recovery regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	markerPath, transactionPath := writeStaleTransactionFixture(t, root, strings.Repeat("2", 32))
	victim := filepath.Join(root, "victim.txt")
	if err := os.WriteFile(victim, []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	repositoryRoot, err := os.OpenRoot(canonicalRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer repositoryRoot.Close()
	recovered, recoveryErr := recoverStaleMarkerRootWithHook(repositoryRoot, filepath.FromSlash(state.TransactionPath), canonicalRoot, false, func(phase string) error {
		if phase != "marker-removal-authority-validated" {
			return nil
		}
		if err := os.Remove(markerPath); err != nil {
			return err
		}
		return os.Symlink("../victim.txt", markerPath)
	})
	if recoveryErr == nil || recovered {
		t.Fatalf("stale marker substitution recovered=%t err=%v", recovered, recoveryErr)
	}
	victimData, err := os.ReadFile(victim)
	if err != nil || string(victimData) != "preserve\n" {
		t.Fatalf("stale marker substitution touched victim: %q err=%v", victimData, err)
	}
	if _, err := os.Stat(filepath.Join(transactionPath, "plan.json")); err != nil {
		t.Fatalf("stale marker substitution deleted transaction evidence: %v", err)
	}
}

func TestStaleRecoveryPreservesCapturedTransactionReplacementAfterValidation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("open-directory stale-recovery regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	markerPath, transactionPath := writeStaleTransactionFixture(t, root, strings.Repeat("3", 32))
	ownedPath := transactionPath + "-owned"
	victim := filepath.Join(root, "victim-transaction")
	if err := os.Mkdir(victim, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(victim, "sentinel.txt"), []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	repositoryRoot, err := os.OpenRoot(canonicalRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer repositoryRoot.Close()
	recovered, recoveryErr := recoverStaleMarkerRootWithHook(repositoryRoot, filepath.FromSlash(state.TransactionPath), canonicalRoot, false, func(phase string) error {
		if phase != "transaction-cleanup-authority-validated" {
			return nil
		}
		if err := os.Rename(transactionPath, ownedPath); err != nil {
			return err
		}
		return os.Rename(victim, transactionPath)
	})
	if recoveryErr == nil || recovered {
		t.Fatalf("stale transaction substitution recovered=%t err=%v", recovered, recoveryErr)
	}
	if _, err := os.Stat(markerPath); !os.IsNotExist(err) {
		t.Fatalf("verified stale marker was not removed before transaction cleanup: %v recovery_err=%v", err, recoveryErr)
	}
	quarantines, err := filepath.Glob(filepath.Join(filepath.Dir(transactionPath), "."+filepath.Base(transactionPath)+".transaction-cleanup-*"))
	if err != nil || len(quarantines) != 1 {
		t.Fatalf("captured stale transaction replacements = %v, err=%v", quarantines, err)
	}
	sentinelData, err := os.ReadFile(filepath.Join(quarantines[0], "sentinel.txt"))
	if err != nil || string(sentinelData) != "preserve\n" {
		t.Fatalf("stale transaction cleanup touched replacement bytes: %q err=%v", sentinelData, err)
	}
	if _, err := os.Stat(filepath.Join(ownedPath, "plan.json")); err != nil {
		t.Fatalf("stale transaction cleanup deleted bound evidence: %v", err)
	}
}

func TestStaleRecoveryPreservesSameInodeMarkerMutationAfterValidation(t *testing.T) {
	root := t.TempDir()
	markerPath, transactionPath := writeStaleTransactionFixture(t, root, strings.Repeat("4", 32))
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	repositoryRoot, err := os.OpenRoot(canonicalRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer repositoryRoot.Close()
	recovered, recoveryErr := recoverStaleMarkerRootWithHook(repositoryRoot, filepath.FromSlash(state.TransactionPath), canonicalRoot, false, func(phase string) error {
		if phase != "marker-removal-authority-validated" {
			return nil
		}
		return os.WriteFile(markerPath, []byte("concurrent marker bytes\n"), 0o600)
	})
	if recoveryErr == nil || recovered {
		t.Fatalf("same-inode stale marker mutation recovered=%t err=%v", recovered, recoveryErr)
	}
	quarantines, err := filepath.Glob(filepath.Join(filepath.Dir(markerPath), "."+filepath.Base(markerPath)+".marker-cleanup-*"))
	if err != nil || len(quarantines) != 1 {
		t.Fatalf("same-inode marker captures = %v, err=%v", quarantines, err)
	}
	markerData, err := os.ReadFile(quarantines[0])
	if err != nil || string(markerData) != "concurrent marker bytes\n" {
		t.Fatalf("same-inode marker evidence changed: %q err=%v", markerData, err)
	}
	if _, err := os.Stat(filepath.Join(transactionPath, "plan.json")); err != nil {
		t.Fatalf("same-inode marker mutation deleted transaction evidence: %v", err)
	}
}

func TestStaleRecoveryPreservesSameInodePlanMutationAfterValidation(t *testing.T) {
	root := t.TempDir()
	markerPath, transactionPath := writeStaleTransactionFixture(t, root, strings.Repeat("5", 32))
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	repositoryRoot, err := os.OpenRoot(canonicalRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer repositoryRoot.Close()
	planPath := filepath.Join(transactionPath, "plan.json")
	recovered, recoveryErr := recoverStaleMarkerRootWithHook(repositoryRoot, filepath.FromSlash(state.TransactionPath), canonicalRoot, false, func(phase string) error {
		if phase != "transaction-cleanup-authority-validated" {
			return nil
		}
		return os.WriteFile(planPath, []byte("concurrent plan bytes\n"), 0o600)
	})
	if recoveryErr == nil || recovered {
		t.Fatalf("same-inode stale plan mutation recovered=%t err=%v", recovered, recoveryErr)
	}
	if _, err := os.Stat(markerPath); !os.IsNotExist(err) {
		t.Fatalf("verified stale marker was not removed before plan evidence check: %v", err)
	}
	quarantines, err := filepath.Glob(filepath.Join(filepath.Dir(transactionPath), "."+filepath.Base(transactionPath)+".transaction-cleanup-*"))
	if err != nil || len(quarantines) != 1 {
		t.Fatalf("same-inode plan captures = %v, err=%v", quarantines, err)
	}
	planData, err := os.ReadFile(filepath.Join(quarantines[0], "plan.json"))
	if err != nil || string(planData) != "concurrent plan bytes\n" {
		t.Fatalf("same-inode plan evidence changed: %q err=%v", planData, err)
	}
}

func TestStaleRecoveryDoesNotClaimPlanningDirectoryWithoutBoundPlan(t *testing.T) {
	root := t.TempDir()
	markerPath, transactionPath := writeStaleTransactionFixture(t, root, strings.Repeat("6", 32))
	markerData, err := os.ReadFile(markerPath)
	if err != nil {
		t.Fatal(err)
	}
	var marker transactionMarker
	if err := json.Unmarshal(markerData, &marker); err != nil {
		t.Fatal(err)
	}
	marker.Phase = "planning"
	markerData, _ = json.Marshal(marker)
	if err := os.WriteFile(markerPath, markerData, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(transactionPath, "plan.json")); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(transactionPath, "sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("preserve unproved directory\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	recovered, recoveryErr := RecoverStaleTransaction(root, false)
	if recoveryErr == nil || recovered {
		t.Fatalf("unproved planning directory recovered=%t err=%v", recovered, recoveryErr)
	}
	sentinelData, err := os.ReadFile(sentinel)
	if err != nil || string(sentinelData) != "preserve unproved directory\n" {
		t.Fatalf("unproved planning directory was changed: %q err=%v", sentinelData, err)
	}
	if _, err := os.Stat(markerPath); err != nil {
		t.Fatalf("unproved planning marker was removed: %v", err)
	}
}

func TestStaleRecoveryRejectsTraversalTransactionIdentityBeforeOpeningDirectory(t *testing.T) {
	root := t.TempDir()
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	maliciousID := "../../../docs/repo/plans"
	recoveryPath := filepath.ToSlash(filepath.Join(".codeheart", "local", "kit-transactions", maliciousID))
	planDirectory := filepath.Join(root, filepath.FromSlash(recoveryPath))
	if err := os.MkdirAll(planDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	sentinelPath := filepath.Join(planDirectory, "sentinel.txt")
	sentinelBefore := []byte("preserve canonical plans directory\n")
	if err := os.WriteFile(sentinelPath, sentinelBefore, 0o600); err != nil {
		t.Fatal(err)
	}
	planPath := filepath.Join(planDirectory, "plan.json")
	planBefore, _ := json.Marshal(Plan{SchemaVersion: 1, ID: maliciousID, Command: "repair", Root: canonicalRoot})
	if err := os.WriteFile(planPath, planBefore, 0o600); err != nil {
		t.Fatal(err)
	}
	marker := transactionMarker{SchemaVersion: 1, TransactionID: maliciousID, Command: "repair", Phase: "staged", PID: 2_000_000_000, StartedAt: "2026-07-31T15:00:00Z", TargetRoot: canonicalRoot, RecoveryPath: recoveryPath}
	markerBefore, _ := json.Marshal(marker)
	markerPath := filepath.Join(root, filepath.FromSlash(state.TransactionPath))
	if err := os.MkdirAll(filepath.Dir(markerPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(markerPath, markerBefore, 0o600); err != nil {
		t.Fatal(err)
	}
	recovered, recoveryErr := RecoverStaleTransaction(root, false)
	if recoveryErr == nil || recovered {
		t.Fatalf("traversal transaction recovered=%t err=%v", recovered, recoveryErr)
	}
	for path, expected := range map[string][]byte{markerPath: markerBefore, planPath: planBefore, sentinelPath: sentinelBefore} {
		actual, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(actual, expected) {
			t.Fatalf("traversal recovery changed %s: %q err=%v", path, actual, err)
		}
	}
}

func writeStaleTransactionFixture(t *testing.T, root, transactionID string) (string, string) {
	t.Helper()
	canonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	recoveryPath := filepath.ToSlash(filepath.Join(".codeheart", "local", "kit-transactions", transactionID))
	marker := transactionMarker{SchemaVersion: 1, TransactionID: transactionID, Command: "repair", Phase: "staged", PID: 2_000_000_000, StartedAt: "2026-07-31T15:00:00Z", TargetRoot: canonical, RecoveryPath: recoveryPath}
	markerData, _ := json.Marshal(marker)
	markerPath := filepath.Join(root, filepath.FromSlash(state.TransactionPath))
	if err := os.MkdirAll(filepath.Dir(markerPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(markerPath, markerData, 0o600); err != nil {
		t.Fatal(err)
	}
	transactionPath := filepath.Join(root, filepath.FromSlash(recoveryPath))
	if err := os.MkdirAll(transactionPath, 0o700); err != nil {
		t.Fatal(err)
	}
	planData, _ := json.Marshal(Plan{SchemaVersion: 1, ID: transactionID, Command: "repair", Root: canonical})
	if err := os.WriteFile(filepath.Join(transactionPath, "plan.json"), planData, 0o600); err != nil {
		t.Fatal(err)
	}
	return markerPath, transactionPath
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
			{Kind: "replace", Target: "a.txt", Owner: "managed", Content: []byte("after-a\n"), ExpectedSHA256: digest([]byte("before\n"))},
			{Kind: "replace", Target: "b.txt", Owner: "managed", Content: []byte("after-b\n"), ExpectedSHA256: digest([]byte("before\n"))},
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

func TestRollbackRefusesReplacedParentWithoutTouchingExternalFiles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory symlink rollback regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	parent := filepath.Join(root, "plans")
	if err := os.Mkdir(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(parent, name), []byte("before\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	outside := t.TempDir()
	external := filepath.Join(outside, "a.txt")
	if err := os.WriteFile(external, []byte("after-a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            "rollback-parent-replacement",
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions: []Action{
			{Kind: "replace", Target: "plans/a.txt", Owner: "managed", Content: []byte("after-a\n"), ExpectedSHA256: digest([]byte("before\n"))},
			{Kind: "replace", Target: "plans/b.txt", Owner: "managed", Content: []byte("after-b\n"), ExpectedSHA256: digest([]byte("before\n"))},
		},
	}
	originalParent := parent + "-original"
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "commit:plans/b.txt" {
			return nil
		}
		if err := os.Rename(parent, originalParent); err != nil {
			return err
		}
		if err := os.Symlink(outside, parent); err != nil {
			return err
		}
		return errors.New("injected failure after parent replacement")
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("parent replacement result = %#v, err=%v", result, err)
	}
	externalData, err := os.ReadFile(external)
	if err != nil || string(externalData) != "after-a\n" {
		t.Fatalf("rollback touched external file: %q err=%v", externalData, err)
	}
	committedData, err := os.ReadFile(filepath.Join(originalParent, "a.txt"))
	if err != nil || string(committedData) != "after-a\n" {
		t.Fatalf("committed evidence was not preserved: %q err=%v", committedData, err)
	}
	markerPath := filepath.Join(root, filepath.FromSlash(state.TransactionPath))
	if _, err := os.Stat(markerPath); err != nil {
		t.Fatalf("recovery marker missing: %v", err)
	}
	backupPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID, "backup", "plans", "a.txt")
	backupData, err := os.ReadFile(backupPath)
	if err != nil || string(backupData) != "before\n" {
		t.Fatalf("rollback backup evidence missing: %q err=%v", backupData, err)
	}
}

func TestRollbackUsesBoundParentWhenPathChangesDuringRollback(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory symlink rollback regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	parent := filepath.Join(root, "plans")
	if err := os.Mkdir(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(parent, name), []byte("before\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	outside := t.TempDir()
	external := filepath.Join(outside, "a.txt")
	if err := os.WriteFile(external, []byte("after-a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            "rollback-parent-race",
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions: []Action{
			{Kind: "replace", Target: "plans/a.txt", Owner: "managed", Content: []byte("after-a\n"), ExpectedSHA256: digest([]byte("before\n"))},
			{Kind: "replace", Target: "plans/b.txt", Owner: "managed", Content: []byte("after-b\n"), ExpectedSHA256: digest([]byte("before\n"))},
		},
	}
	originalParent := parent + "-original"
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		switch phase {
		case "commit:plans/b.txt":
			return errors.New("injected failure after first commit")
		case "rollback-quarantined:plans/a.txt":
			if err := os.Rename(parent, originalParent); err != nil {
				return err
			}
			return os.Symlink(outside, parent)
		default:
			return nil
		}
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("during-rollback parent replacement result = %#v, err=%v", result, err)
	}
	externalData, err := os.ReadFile(external)
	if err != nil || string(externalData) != "after-a\n" {
		t.Fatalf("bound rollback touched external file: %q err=%v", externalData, err)
	}
	quarantines, err := filepath.Glob(filepath.Join(originalParent, ".a.txt.rollback-current-*"))
	if err != nil || len(quarantines) != 1 {
		t.Fatalf("rollback quarantine evidence = %v, err=%v", quarantines, err)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.TransactionPath))); err != nil {
		t.Fatalf("recovery marker missing: %v", err)
	}
}

func TestUnpredictableActionCapturesPreservePrecreatedLegacyStashNames(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("before\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	legacyName := func(suffix string) string {
		digest := sha256.Sum256([]byte("a.txt\x00" + suffix))
		return filepath.Join(root, ".a.txt.codeheart-"+suffix+"-"+hex.EncodeToString(digest[:6]))
	}
	legacyOriginal := legacyName("original")
	legacyRollback := legacyName("rollback-current")
	plan := Plan{SchemaVersion: 1, ID: "unpredictable-action-captures", Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{
		{Kind: "replace", Target: "a.txt", Owner: "managed", Content: []byte("after-a\n"), ExpectedSHA256: digest([]byte("before\n"))},
		{Kind: "replace", Target: "b.txt", Owner: "managed", Content: []byte("after-b\n"), ExpectedSHA256: digest([]byte("before\n"))},
	}}
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		switch phase {
		case "pre-commit:a.txt":
			return os.WriteFile(legacyOriginal, []byte("preserve original-name sentinel\n"), 0o600)
		case "commit:b.txt":
			if err := os.WriteFile(legacyRollback, []byte("preserve rollback-name sentinel\n"), 0o600); err != nil {
				return err
			}
			return errors.New("injected failure after first commit")
		default:
			return nil
		}
	}})
	if err != nil || result.Status != StatusRolledBack || !result.Rollback.Succeeded {
		t.Fatalf("unpredictable capture rollback result = %#v, err=%v", result, err)
	}
	for path, expected := range map[string]string{legacyOriginal: "preserve original-name sentinel\n", legacyRollback: "preserve rollback-name sentinel\n"} {
		data, err := os.ReadFile(path)
		if err != nil || string(data) != expected {
			t.Fatalf("precreated stash sentinel %s changed: %q err=%v", path, data, err)
		}
	}
	targetData, err := os.ReadFile(filepath.Join(root, "a.txt"))
	if err != nil || string(targetData) != "before\n" {
		t.Fatalf("rollback did not restore target: %q err=%v", targetData, err)
	}
}

func TestApplyRejectsExternalTransactionTreeSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory symlink containment regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".codeheart"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, ".codeheart", "local")); err != nil {
		t.Fatal(err)
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            "external-transaction-tree",
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions:       []Action{{Kind: "create", Target: "plan.txt", Owner: "managed", Content: []byte("plan\n")}},
	}
	result, err := Apply(plan, ApplyOptions{})
	if err != nil || result.Status != StatusBlocked || len(result.Blockers) == 0 || result.Blockers[0].Code != "unsafe_transaction_path" {
		t.Fatalf("external transaction tree result = %#v, err=%v", result, err)
	}
	entries, err := os.ReadDir(outside)
	if err != nil || len(entries) != 0 {
		t.Fatalf("transaction escaped repository: entries=%v err=%v", entries, err)
	}
	if _, err := os.Stat(filepath.Join(root, "plan.txt")); !os.IsNotExist(err) {
		t.Fatalf("blocked transaction changed target: %v", err)
	}
}

func TestApplyRejectsInternalTransactionTreeSymlinkWithoutCleanup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory symlink containment regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	victim := filepath.Join(root, "victim")
	transactionID := "internal-transaction-tree"
	victimTransaction := filepath.Join(victim, "kit-transactions", transactionID)
	if err := os.MkdirAll(victimTransaction, 0o755); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(victimTransaction, "sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("preserve\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".codeheart"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../victim", filepath.Join(root, ".codeheart", "local")); err != nil {
		t.Fatal(err)
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            transactionID,
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions:       []Action{{Kind: "create", Target: "plan.txt", Owner: "managed", Content: []byte("plan\n")}},
	}
	result, err := Apply(plan, ApplyOptions{})
	if err != nil || result.Status != StatusBlocked || len(result.Blockers) == 0 || result.Blockers[0].Code != "unsafe_transaction_path" {
		t.Fatalf("internal transaction tree result = %#v, err=%v", result, err)
	}
	data, err := os.ReadFile(sentinel)
	if err != nil || string(data) != "preserve\n" {
		t.Fatalf("internal transaction cleanup touched victim: %q err=%v", data, err)
	}
}

func TestApplyPreservesPreexistingTransactionDirectory(t *testing.T) {
	root := t.TempDir()
	transactionID := "preexisting-transaction-directory"
	transactionPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", transactionID)
	if err := os.MkdirAll(transactionPath, 0o700); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(transactionPath, "victim-sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            transactionID,
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions:       []Action{{Kind: "create", Target: "plan.txt", Owner: "managed", Content: []byte("plan\n")}},
	}
	result, err := Apply(plan, ApplyOptions{})
	if err != nil || result.Status != StatusBlocked || len(result.Blockers) == 0 || result.Blockers[0].Code != "transaction_path_exists" {
		t.Fatalf("preexisting transaction result = %#v, err=%v", result, err)
	}
	data, err := os.ReadFile(sentinel)
	if err != nil || string(data) != "preserve\n" {
		t.Fatalf("preexisting transaction evidence changed: %q err=%v", data, err)
	}
	if _, err := os.Stat(filepath.Join(root, "plan.txt")); !os.IsNotExist(err) {
		t.Fatalf("blocked transaction changed target: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.TransactionPath))); !os.IsNotExist(err) {
		t.Fatalf("blocked transaction retained marker: %v", err)
	}
}

func TestApplyDoesNotAutoRecoverStaleMarkerOrTransactionDirectory(t *testing.T) {
	root := t.TempDir()
	transactionID := "stale-marker-owned-explicitly"
	recoveryPath := filepath.ToSlash(filepath.Join(".codeheart", "local", "kit-transactions", transactionID))
	transactionPath := filepath.Join(root, filepath.FromSlash(recoveryPath))
	if err := os.MkdirAll(transactionPath, 0o700); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(transactionPath, "victim-sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	marker := transactionMarker{SchemaVersion: 1, TransactionID: transactionID, Command: "sync", Phase: "planning", PID: 99999999, StartedAt: "2026-07-31T14:00:00Z", TargetRoot: root, RecoveryPath: recoveryPath}
	markerData, err := json.Marshal(marker)
	if err != nil {
		t.Fatal(err)
	}
	markerPath := filepath.Join(root, filepath.FromSlash(state.TransactionPath))
	if err := os.MkdirAll(filepath.Dir(markerPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(markerPath, append(markerData, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := Plan{SchemaVersion: 1, ID: transactionID, Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{{Kind: "create", Target: "plan.txt", Owner: "managed", Content: []byte("plan\n")}}}
	result, err := Apply(plan, ApplyOptions{})
	if err != nil || result.Status != StatusBlocked || len(result.Blockers) == 0 || result.Blockers[0].Code != "transaction_in_progress" {
		t.Fatalf("stale marker apply result = %#v, err=%v", result, err)
	}
	data, err := os.ReadFile(sentinel)
	if err != nil || string(data) != "preserve\n" {
		t.Fatalf("implicit recovery changed transaction evidence: %q err=%v", data, err)
	}
	if _, err := os.Stat(markerPath); err != nil {
		t.Fatalf("implicit recovery removed stale marker: %v", err)
	}
}

func TestPreCommitRollbackUsesBoundTransactionRoot(t *testing.T) {
	root := t.TempDir()
	transactionID := "bound-precommit-cleanup"
	transactionPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", transactionID)
	victim := filepath.Join(root, "victim")
	if err := os.Mkdir(victim, 0o755); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(victim, "sentinel.txt")
	if err := os.WriteFile(sentinel, []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := Plan{SchemaVersion: 1, ID: transactionID, Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{{Kind: "create", Target: "plan.txt", Owner: "managed", Content: []byte("plan\n")}}}
	ownedPath := transactionPath + "-owned"
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "staged" {
			return nil
		}
		if err := os.Rename(transactionPath, ownedPath); err != nil {
			return err
		}
		if err := os.Rename(victim, transactionPath); err != nil {
			return err
		}
		return errors.New("injected staged failure after transaction path substitution")
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("bound precommit cleanup result = %#v, err=%v", result, err)
	}
	data, err := os.ReadFile(filepath.Join(transactionPath, "sentinel.txt"))
	if err != nil || string(data) != "preserve\n" {
		t.Fatalf("bound cleanup touched substituted directory: %q err=%v", data, err)
	}
}

func TestMarkerUpdatesRefuseInternalSymlinkSubstitution(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("marker symlink regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	victim := filepath.Join(root, "victim.txt")
	if err := os.WriteFile(victim, []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := Plan{SchemaVersion: 1, ID: "marker-symlink-substitution", Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{{Kind: "create", Target: "plan.txt", Owner: "managed", Content: []byte("plan\n")}}}
	markerPath := filepath.Join(root, filepath.FromSlash(state.TransactionPath))
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "marker" {
			return nil
		}
		if err := os.Remove(markerPath); err != nil {
			return err
		}
		return os.Symlink("../victim.txt", markerPath)
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("marker substitution result = %#v, err=%v", result, err)
	}
	data, err := os.ReadFile(victim)
	if err != nil || string(data) != "preserve\n" {
		t.Fatalf("marker update touched symlink target: %q err=%v", data, err)
	}
}

func TestPreAcquisitionFailureDoesNotCleanNewlyAppearingTransactionDirectory(t *testing.T) {
	root := t.TempDir()
	transactionID := "preacquisition-directory-race"
	transactionPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", transactionID)
	sentinel := filepath.Join(transactionPath, "sentinel.txt")
	plan := Plan{SchemaVersion: 1, ID: transactionID, Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{{Kind: "create", Target: "plan.txt", Owner: "managed", Content: []byte("plan\n")}}}
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "marker" {
			return nil
		}
		if err := os.MkdirAll(transactionPath, 0o700); err != nil {
			return err
		}
		if err := os.WriteFile(sentinel, []byte("preserve\n"), 0o600); err != nil {
			return err
		}
		return errors.New("injected pre-acquisition failure")
	}})
	if err != nil || result.Status != StatusRolledBack || !result.Rollback.Succeeded {
		t.Fatalf("pre-acquisition failure result = %#v, err=%v", result, err)
	}
	data, err := os.ReadFile(sentinel)
	if err != nil || string(data) != "preserve\n" {
		t.Fatalf("pre-acquisition cleanup touched new directory: %q err=%v", data, err)
	}
}

func TestCommittedRollbackTreatsMarkerSubstitutionAsRecoveryRequired(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("marker symlink regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("before\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	victim := filepath.Join(root, "victim.txt")
	if err := os.WriteFile(victim, []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := Plan{SchemaVersion: 1, ID: "committed-marker-substitution", Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{{Kind: "replace", Target: "a.txt", Owner: "managed", Content: []byte("after\n"), ExpectedSHA256: digest([]byte("before\n"))}}}
	markerPath := filepath.Join(root, filepath.FromSlash(state.TransactionPath))
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "post-check" {
			return nil
		}
		if err := os.Remove(markerPath); err != nil {
			return err
		}
		if err := os.Symlink("../victim.txt", markerPath); err != nil {
			return err
		}
		return errors.New("injected post-check failure after marker substitution")
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("committed marker substitution result = %#v, err=%v", result, err)
	}
	targetData, err := os.ReadFile(filepath.Join(root, "a.txt"))
	if err != nil || string(targetData) != "before\n" {
		t.Fatalf("target rollback failed: %q err=%v", targetData, err)
	}
	victimData, err := os.ReadFile(victim)
	if err != nil || string(victimData) != "preserve\n" {
		t.Fatalf("marker substitution touched victim: %q err=%v", victimData, err)
	}
	transactionPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID)
	if _, err := os.Stat(transactionPath); err != nil {
		t.Fatalf("marker authority failure did not retain transaction evidence: %v", err)
	}
}

func TestCommittedRollbackValidatesMarkerAuthorityBeforeDeletingEvidence(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("marker symlink rollback regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("before\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	victim := filepath.Join(root, "victim.txt")
	if err := os.WriteFile(victim, []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := Plan{SchemaVersion: 1, ID: "rollback-marker-authority", Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{
		{Kind: "replace", Target: "a.txt", Owner: "managed", Content: []byte("after-a\n"), ExpectedSHA256: digest([]byte("before\n"))},
		{Kind: "replace", Target: "b.txt", Owner: "managed", Content: []byte("after-b\n"), ExpectedSHA256: digest([]byte("before\n"))},
	}}
	markerPath := filepath.Join(root, filepath.FromSlash(state.TransactionPath))
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		switch phase {
		case "commit:b.txt":
			return errors.New("injected failure after first commit")
		case "marker-removal-authority-validated":
			if err := os.Remove(markerPath); err != nil {
				return err
			}
			return os.Symlink("../victim.txt", markerPath)
		default:
			return nil
		}
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("rollback marker authority result = %#v, err=%v", result, err)
	}
	targetData, err := os.ReadFile(filepath.Join(root, "a.txt"))
	if err != nil || string(targetData) != "before\n" {
		t.Fatalf("target rollback failed: %q err=%v", targetData, err)
	}
	victimData, err := os.ReadFile(victim)
	if err != nil || string(victimData) != "preserve\n" {
		t.Fatalf("marker authority loss touched victim: %q err=%v", victimData, err)
	}
	transactionPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID)
	if _, err := os.Stat(filepath.Join(transactionPath, "plan.json")); err != nil {
		t.Fatalf("marker authority loss deleted transaction evidence: %v", err)
	}
}

func TestCommittedRollbackValidatesTransactionAuthorityBeforeDeletingEvidence(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("open-directory rename rollback regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("before\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	plan := Plan{SchemaVersion: 1, ID: "rollback-transaction-authority", Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{
		{Kind: "replace", Target: "a.txt", Owner: "managed", Content: []byte("after-a\n"), ExpectedSHA256: digest([]byte("before\n"))},
		{Kind: "replace", Target: "b.txt", Owner: "managed", Content: []byte("after-b\n"), ExpectedSHA256: digest([]byte("before\n"))},
	}}
	transactionPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID)
	ownedPath := transactionPath + "-owned"
	victim := filepath.Join(root, "victim-transaction")
	if err := os.Mkdir(victim, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(victim, "sentinel.txt"), []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		switch phase {
		case "commit:b.txt":
			return errors.New("injected failure after first commit")
		case "transaction-cleanup-authority-validated":
			if err := os.Rename(transactionPath, ownedPath); err != nil {
				return err
			}
			return os.Rename(victim, transactionPath)
		default:
			return nil
		}
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("rollback transaction authority result = %#v, err=%v", result, err)
	}
	targetData, err := os.ReadFile(filepath.Join(root, "a.txt"))
	if err != nil || string(targetData) != "before\n" {
		t.Fatalf("target rollback failed: %q err=%v", targetData, err)
	}
	quarantines, err := filepath.Glob(filepath.Join(filepath.Dir(transactionPath), "."+filepath.Base(transactionPath)+".transaction-cleanup-*"))
	if err != nil || len(quarantines) != 1 {
		t.Fatalf("captured transaction replacements = %v, err=%v", quarantines, err)
	}
	sentinelData, err := os.ReadFile(filepath.Join(quarantines[0], "sentinel.txt"))
	if err != nil || string(sentinelData) != "preserve\n" {
		t.Fatalf("transaction authority loss touched substituted directory: %q err=%v", sentinelData, err)
	}
	if _, err := os.Stat(filepath.Join(ownedPath, "plan.json")); err != nil {
		t.Fatalf("transaction authority loss deleted bound evidence: %v", err)
	}
}

func TestSuccessfulCommitPreservesEvidenceWhenMarkerAuthorityChangesAfterValidation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("marker symlink finalization regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	victim := filepath.Join(root, "victim.txt")
	if err := os.WriteFile(victim, []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := Plan{SchemaVersion: 1, ID: "success-marker-authority", Command: "sync", Root: root, StateBefore: string(state.StateAdoptable), ExpectedAfter: []state.Classification{state.StateAdoptable}, Actions: []Action{{Kind: "create", Target: "plan.txt", Owner: "managed", Content: []byte("created\n")}}}
	markerPath := filepath.Join(root, filepath.FromSlash(state.TransactionPath))
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "marker-removal-authority-validated" {
			return nil
		}
		if err := os.Remove(markerPath); err != nil {
			return err
		}
		return os.Symlink("../victim.txt", markerPath)
	}})
	if err != nil || result.Status != StatusRecoveryRequired {
		t.Fatalf("successful commit marker authority result = %#v, err=%v", result, err)
	}
	created, err := os.ReadFile(filepath.Join(root, "plan.txt"))
	if err != nil || string(created) != "created\n" {
		t.Fatalf("committed target changed during finalization: %q err=%v", created, err)
	}
	victimData, err := os.ReadFile(victim)
	if err != nil || string(victimData) != "preserve\n" {
		t.Fatalf("marker finalization touched victim: %q err=%v", victimData, err)
	}
	transactionPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID)
	if _, err := os.Stat(filepath.Join(transactionPath, "plan.json")); err != nil {
		t.Fatalf("marker finalization deleted transaction evidence: %v", err)
	}
}

func TestSuccessfulCommitPreservesEvidenceWhenTransactionAuthorityChangesAfterValidation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("open-directory rename finalization regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	plan := Plan{SchemaVersion: 1, ID: "success-transaction-authority", Command: "sync", Root: root, StateBefore: string(state.StateAdoptable), ExpectedAfter: []state.Classification{state.StateAdoptable}, Actions: []Action{{Kind: "create", Target: "plan.txt", Owner: "managed", Content: []byte("created\n")}}}
	transactionPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID)
	ownedPath := transactionPath + "-owned"
	victim := filepath.Join(root, "victim-transaction-success")
	if err := os.Mkdir(victim, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(victim, "sentinel.txt"), []byte("preserve\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "transaction-cleanup-authority-validated" {
			return nil
		}
		if err := os.Rename(transactionPath, ownedPath); err != nil {
			return err
		}
		return os.Rename(victim, transactionPath)
	}})
	if err != nil || result.Status != StatusRecoveryRequired {
		t.Fatalf("successful commit transaction authority result = %#v, err=%v", result, err)
	}
	created, err := os.ReadFile(filepath.Join(root, "plan.txt"))
	if err != nil || string(created) != "created\n" {
		t.Fatalf("committed target changed during finalization: %q err=%v", created, err)
	}
	quarantines, err := filepath.Glob(filepath.Join(filepath.Dir(transactionPath), "."+filepath.Base(transactionPath)+".transaction-cleanup-*"))
	if err != nil || len(quarantines) != 1 {
		t.Fatalf("captured successful transaction replacements = %v, err=%v", quarantines, err)
	}
	sentinelData, err := os.ReadFile(filepath.Join(quarantines[0], "sentinel.txt"))
	if err != nil || string(sentinelData) != "preserve\n" {
		t.Fatalf("successful finalization touched replacement bytes: %q err=%v", sentinelData, err)
	}
	if _, err := os.Stat(filepath.Join(ownedPath, "plan.json")); err != nil {
		t.Fatalf("successful finalization deleted bound transaction evidence: %v", err)
	}
}

func TestUnsafeParentPreservesUnlinkedBoundBackup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory symlink recovery regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	parent := filepath.Join(root, "plans")
	if err := os.Mkdir(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(parent, name), []byte("before\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	outside := t.TempDir()
	plan := Plan{SchemaVersion: 1, ID: "unsafe-parent-bound-backup", Command: "sync", Root: root, StateBefore: string(state.StateCurrent), Actions: []Action{
		{Kind: "replace", Target: "plans/a.txt", Owner: "managed", Content: []byte("after-a\n"), ExpectedSHA256: digest([]byte("before\n"))},
		{Kind: "replace", Target: "plans/b.txt", Owner: "managed", Content: []byte("after-b\n"), ExpectedSHA256: digest([]byte("before\n"))},
	}}
	backupPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID, "backup", "plans", "a.txt")
	transactionRoot := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID)
	deterministicRecoveryPath := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID, "recovery-original", "plans", "a.txt.original")
	originalParent := parent + "-original"
	recoveryMutations := 0
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if strings.HasPrefix(phase, "recovery-copy-after-directory-sync:") {
			if recoveryMutations > 0 {
				return nil
			}
			recoveryMutations++
			relative := strings.TrimPrefix(phase, "recovery-copy-after-directory-sync:")
			return os.WriteFile(filepath.Join(transactionRoot, filepath.FromSlash(relative)), []byte("concurrent\n"), 0o644)
		}
		if phase != "commit:plans/b.txt" {
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(deterministicRecoveryPath), 0o700); err != nil {
			return err
		}
		if err := os.WriteFile(deterministicRecoveryPath, []byte("attacker\n"), 0o600); err != nil {
			return err
		}
		if err := os.Remove(backupPath); err != nil {
			return err
		}
		if err := os.Rename(parent, originalParent); err != nil {
			return err
		}
		if err := os.Symlink(outside, parent); err != nil {
			return err
		}
		return errors.New("injected parent and backup substitution")
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("combined parent/backup substitution result = %#v, err=%v", result, err)
	}
	recoveryPaths, err := filepath.Glob(filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID, "recovery-original", "plans", "a.txt.original-*"))
	if err != nil || len(recoveryPaths) != 2 {
		t.Fatalf("random recovery paths = %v, err=%v", recoveryPaths, err)
	}
	recoveryContents := map[string]int{}
	for _, recoveryPath := range recoveryPaths {
		recoveryData, readErr := os.ReadFile(recoveryPath)
		if readErr != nil {
			t.Fatal(readErr)
		}
		recoveryContents[string(recoveryData)]++
	}
	if recoveryContents["before\n"] != 1 || recoveryContents["concurrent\n"] != 1 {
		t.Fatalf("recovery retries did not preserve verified original and concurrent evidence: %#v", recoveryContents)
	}
	deterministicData, err := os.ReadFile(deterministicRecoveryPath)
	if err != nil || string(deterministicData) != "attacker\n" {
		t.Fatalf("pre-existing deterministic recovery path was changed: %q err=%v", deterministicData, err)
	}
	entries, err := os.ReadDir(outside)
	if err != nil || len(entries) != 0 {
		t.Fatalf("unsafe parent recovery touched external directory: entries=%v err=%v", entries, err)
	}
}

func TestRollbackRefusesSubstitutedBackupBytes(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("before\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            "rollback-backup-substitution",
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions: []Action{
			{Kind: "replace", Target: "a.txt", Owner: "managed", Content: []byte("after-a\n"), ExpectedSHA256: digest([]byte("before\n"))},
			{Kind: "replace", Target: "b.txt", Owner: "managed", Content: []byte("after-b\n"), ExpectedSHA256: digest([]byte("before\n"))},
		},
	}
	backup := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID, "backup", "a.txt")
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		switch phase {
		case "commit:b.txt":
			return errors.New("injected failure after first commit")
		case "rollback-quarantined:a.txt":
			if err := os.Remove(backup); err != nil {
				return err
			}
			return os.WriteFile(backup, []byte("attacker\n"), 0o644)
		default:
			return nil
		}
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("backup substitution result = %#v, err=%v", result, err)
	}
	targetData, err := os.ReadFile(filepath.Join(root, "a.txt"))
	if err != nil || string(targetData) != "before\n" {
		t.Fatalf("bound original backup was not restored: %q err=%v", targetData, err)
	}
	backupData, err := os.ReadFile(backup)
	if err != nil || string(backupData) != "attacker\n" {
		t.Fatalf("substituted backup evidence not preserved: %q err=%v", backupData, err)
	}
	quarantines, err := filepath.Glob(filepath.Join(root, ".a.txt.rollback-current-*"))
	if err != nil || len(quarantines) != 1 {
		t.Fatalf("installed-byte quarantines = %v, err=%v", quarantines, err)
	}
	quarantineData, err := os.ReadFile(quarantines[0])
	if err != nil || string(quarantineData) != "after-a\n" {
		t.Fatalf("installed-byte quarantine evidence not preserved: %q err=%v", quarantineData, err)
	}
}

func TestRollbackPreservesQuarantineChangedAfterHook(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("before\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            "rollback-quarantine-change",
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions: []Action{
			{Kind: "replace", Target: "a.txt", Owner: "managed", Content: []byte("after-a\n"), ExpectedSHA256: digest([]byte("before\n"))},
			{Kind: "replace", Target: "b.txt", Owner: "managed", Content: []byte("after-b\n"), ExpectedSHA256: digest([]byte("before\n"))},
		},
	}
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		switch phase {
		case "commit:b.txt":
			return errors.New("injected failure after first commit")
		case "rollback-quarantined:a.txt":
			quarantines, err := filepath.Glob(filepath.Join(root, ".a.txt.rollback-current-*"))
			if err != nil || len(quarantines) != 1 {
				return fmt.Errorf("find rollback quarantine: paths=%v err=%v", quarantines, err)
			}
			return os.WriteFile(quarantines[0], []byte("concurrent\n"), 0o644)
		default:
			return nil
		}
	}})
	if err != nil || result.Status != StatusRecoveryRequired || result.Rollback.Succeeded {
		t.Fatalf("quarantine mutation result = %#v, err=%v", result, err)
	}
	data, err := os.ReadFile(filepath.Join(root, "a.txt"))
	if err != nil || string(data) != "concurrent\n" {
		t.Fatalf("concurrent quarantine bytes not preserved at target: %q err=%v", data, err)
	}
	backup := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.ID, "backup", "a.txt")
	backupData, err := os.ReadFile(backup)
	if err != nil || string(backupData) != "before\n" {
		t.Fatalf("original backup evidence not preserved: %q err=%v", backupData, err)
	}
}

func TestCommitRejectsRedirectedInRepositoryParentBinding(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory symlink binding regression is covered by the Unix validation jobs")
	}
	root := t.TempDir()
	parent := filepath.Join(root, "plans")
	alternate := filepath.Join(root, "alternate")
	if err := os.Mkdir(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(alternate, 0o755); err != nil {
		t.Fatal(err)
	}
	plan := Plan{
		SchemaVersion: 1,
		ID:            "redirected-parent-binding",
		Command:       "sync",
		Root:          root,
		StateBefore:   string(state.StateCurrent),
		Actions:       []Action{{Kind: "create", Target: "plans/a.txt", Owner: "managed", Content: []byte("created\n")}},
	}
	result, err := Apply(plan, ApplyOptions{Hook: func(phase string) error {
		if phase != "bind-parent:plans/a.txt" {
			return nil
		}
		if err := os.Rename(parent, parent+"-original"); err != nil {
			return err
		}
		return os.Symlink("alternate", parent)
	}})
	if err != nil || result.Status != StatusRolledBack || !result.Rollback.Succeeded {
		t.Fatalf("redirected parent binding result = %#v, err=%v", result, err)
	}
	if _, err := os.Stat(filepath.Join(alternate, "a.txt")); !os.IsNotExist(err) {
		t.Fatalf("redirected parent received target: %v", err)
	}
}

func TestLifecycleCreatePreservesTargetAppearingAfterPlanning(t *testing.T) {
	root := t.TempDir()
	observed, err := state.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	graph, err := state.CompileGraph("standard")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 7, 31, 13, 0, 0, 0, time.UTC)
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
	userBytes := []byte("user-created-after-planning\n")
	result, err := Apply(plan, ApplyOptions{Now: now, Hook: func(phase string) error {
		if phase != "pre-commit:"+state.ConfigPath {
			return nil
		}
		target := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, userBytes, 0o644)
	}})
	if err != nil || result.Status != StatusRolledBack || !result.Rollback.Succeeded {
		t.Fatalf("lifecycle create race result = %#v, err=%v", result, err)
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(state.ConfigPath)))
	if err != nil || !bytes.Equal(data, userBytes) {
		t.Fatalf("lifecycle create overwrote concurrent file: %q err=%v", data, err)
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
