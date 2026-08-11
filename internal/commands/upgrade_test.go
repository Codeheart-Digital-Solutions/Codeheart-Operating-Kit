package commands

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/lockfile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	releasekit "github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/release"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

func TestMain(m *testing.M) {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version":
			fmt.Printf("codeheart-operating-kit %s\n", version.Version)
			os.Exit(0)
		case "__upgrade-reconcile":
			if os.Getenv("CODEHEART_UPGRADE_TEST_RECOVERY_RECONCILE") == "1" {
				os.Exit(3)
			}
			if os.Getenv("CODEHEART_UPGRADE_TEST_FAIL_RECONCILE") == "1" {
				if signal := os.Getenv("CODEHEART_UPGRADE_TEST_FAILURE_SIGNAL"); signal != "" {
					_ = os.WriteFile(signal, []byte("target reconcile failed\n"), 0o600)
				}
				os.Exit(23)
			}
			os.Exit(RunUpgradeReconcile(os.Args[2:], os.Stdout, os.Stderr))
		case "__upgrade-handoff":
			os.Exit(RunUpgradeHandoff(os.Args[2:], os.Stderr))
		case "__cleanup-upgrade-handoff":
			os.Exit(RunCleanupUpgradeHandoff(os.Args[2:], os.Stderr))
		case "__test-upgrade-apply":
			os.Exit(RunUpgrade(os.Args[2:], os.Stdout, os.Stderr))
		}
	}
	os.Exit(m.Run())
}

func TestLegacyV1UpgradeDryRunAndAtomicApply(t *testing.T) {
	if version.Version != "0.1.30" {
		t.Fatalf("test requires release version 0.1.30, got %s", version.Version)
	}
	root := materializeLegacyV119Consumer(t)
	installed := writePreviousBinary(t)
	catalog, evidence := writeUpgradePack(t)

	configBefore := mustRead(t, filepath.Join(root, filepath.FromSlash(state.ConfigPath)))
	treeBefore := snapshotTree(t, root)
	binaryBefore := mustRead(t, installed)
	var preview bytes.Buffer
	code := RunUpgrade([]string{root, "--version", "0.1.30", "--catalog", catalog, "--installed-binary", installed, "--dry-run", "--json"}, &preview, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("legacy dry-run exit = %d; output=%s", code, preview.String())
	}
	if !equalTreeSnapshot(treeBefore, snapshotTree(t, root)) || !bytes.Equal(binaryBefore, mustRead(t, installed)) {
		t.Fatal("legacy dry-run mutated repository or installed binary")
	}
	var payload map[string]any
	if err := json.Unmarshal(preview.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	result := state.Map(payload["result"])
	if result["status"] != "planned" || payload["source_lock_schema"] != float64(1) || payload["target_lock_schema"] != float64(2) || payload["installed_version"] != "0.1.19" || payload["target_version"] != "0.1.30" || payload["migrated_lock"] != true {
		t.Fatalf("legacy preview = %#v", payload)
	}
	if !resultHasChange(result, state.LockPath, "replace") || !resultHasValidation(result, "schema-and-version-migration") {
		t.Fatalf("legacy preview result = %#v", result)
	}

	code, applyOutput := runApprovedUpgrade(t, []string{root, "--version", "0.1.30", "--catalog", catalog, "--installed-binary", installed, "--yes", "--json"})
	if code != 0 {
		t.Fatalf("legacy apply exit = %d; output=%s", code, applyOutput)
	}
	waitForUpgrade(t, root, installed)
	observed, err := state.Inspect(root)
	if err != nil || observed.Classification != state.StateCurrent {
		t.Fatalf("upgraded state = %#v, err=%v", observed, err)
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatal(err)
	}
	if state.AsInt(lock["schema_version"]) != 2 || state.AsString(lock["kit_version"]) != "0.1.30" || state.AsInt(lock["state_generation"]) != 1 {
		t.Fatalf("upgraded lock identity = %#v", lock)
	}
	operation := state.Map(lock["last_operation"])
	if state.AsString(operation["command"]) != "upgrade" || state.AsInt(operation["previous_generation"]) != 0 || state.AsString(operation["transaction_id"]) == "" || state.AsString(operation["transaction_id"]) == "migration-preview" || state.AsString(operation["transaction_id"]) == "planning" {
		t.Fatalf("upgrade operation evidence = %#v", operation)
	}
	provenance := state.Map(lock["release_provenance"])
	for key, expected := range evidence {
		if state.AsString(provenance[key]) != expected {
			t.Fatalf("provenance %s = %q, want %q; all=%#v", key, state.AsString(provenance[key]), expected, provenance)
		}
	}
	capabilities := state.Map(lock["native_capabilities"])
	for _, capability := range []string{"documents", "spreadsheets", "presentations", "browser", "pdf"} {
		record := state.Map(capabilities[capability])
		if state.AsString(record["status"]) != "unknown" || state.AsString(record["profile_applicability"]) != "standard" || state.AsString(record["command_result_category"]) != "not-checked" {
			t.Fatalf("migrated native capability %s = %#v", capability, record)
		}
	}
	if !bytes.Equal(configBefore, mustRead(t, filepath.Join(root, filepath.FromSlash(state.ConfigPath)))) {
		t.Fatal("upgrade changed consumer configuration")
	}
	for relative, expected := range map[string]string{
		"docs/repo/plans/plan-register.md":    "consumer plan register\n",
		"docs/agent-memory/session-ledger.md": "consumer memory ledger\n",
		".codeheart/user/preferences.yaml":    "language: de\n",
	} {
		if string(mustRead(t, filepath.Join(root, filepath.FromSlash(relative)))) != expected {
			t.Fatalf("upgrade changed create-once or consumer-owned surface %s", relative)
		}
	}
	agents := string(mustRead(t, filepath.Join(root, "AGENTS.md")))
	if !strings.Contains(agents, "# Consumer-Owned Guidance\n\npreserve this suffix\n") {
		t.Fatal("upgrade did not preserve AGENTS content outside the managed block")
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(legacyMissingManagedPath))); err != nil {
		t.Fatalf("upgrade did not create current managed addition: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.TransactionPath))); !os.IsNotExist(err) {
		t.Fatalf("successful upgrade left transaction evidence: %v", err)
	}
}

func TestLegacyV1UpgradeRejectsInvalidAmbiguousAndDirtySources(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, string)
	}{
		{name: "invalid config", mutate: func(t *testing.T, root string) {
			path := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
			if err := os.WriteFile(path, append(mustRead(t, path), []byte("unexpected: true\n")...), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "malformed lock", mutate: func(t *testing.T, root string) {
			if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(state.LockPath)), []byte("schema_version: [\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "unsupported lock", mutate: func(t *testing.T, root string) {
			path := filepath.Join(root, filepath.FromSlash(state.LockPath))
			data := bytes.Replace(mustRead(t, path), []byte("schema_version: 1"), []byte("schema_version: 3"), 1)
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "modified managed path", mutate: func(t *testing.T, root string) {
			lock, _ := lockfile.ReadLock(root)
			record := state.Map(state.AnySlice(lock["managed_paths"])[0])
			if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(state.AsString(record["path"]))), []byte("dirty\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "missing managed path", mutate: func(t *testing.T, root string) {
			lock, _ := lockfile.ReadLock(root)
			record := state.Map(state.AnySlice(lock["managed_paths"])[0])
			if err := os.Remove(filepath.Join(root, filepath.FromSlash(state.AsString(record["path"])))); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "duplicate managed authority", mutate: func(t *testing.T, root string) {
			lock, _ := lockfile.ReadLock(root)
			managed := state.AnySlice(lock["managed_paths"])
			lock["managed_paths"] = append(managed, state.DeepCopy(state.Map(managed[0])))
			if err := lockfile.WriteLock(root, lock); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "unsafe managed authority", mutate: func(t *testing.T, root string) {
			lock, _ := lockfile.ReadLock(root)
			state.Map(state.AnySlice(lock["managed_paths"])[0])["path"] = "../outside.md"
			if err := lockfile.WriteLock(root, lock); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "component mismatch", mutate: func(t *testing.T, root string) {
			lock, _ := lockfile.ReadLock(root)
			lock["selected_components"] = []any{"agent-interface"}
			if err := lockfile.WriteLock(root, lock); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "portable case-colliding managed authority", mutate: func(t *testing.T, root string) {
			lock, _ := lockfile.ReadLock(root)
			managed := state.AnySlice(lock["managed_paths"])
			duplicate := state.DeepCopy(state.Map(managed[0]))
			state.Map(duplicate)["path"] = strings.ToUpper(state.AsString(state.Map(duplicate)["path"]))
			lock["managed_paths"] = append(managed, duplicate)
			if err := lockfile.WriteLock(root, lock); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "reversed managed markers", mutate: func(t *testing.T, root string) {
			content := "# Legacy Fixture\n\n" + reconcile.EndMarker + "\nlegacy\n" + reconcile.BeginMarker + "\n"
			if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "target-added managed path overlap", mutate: func(t *testing.T, root string) {
			path := filepath.Join(root, filepath.FromSlash(legacyMissingManagedPath))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte("consumer-owned overlap\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "transaction in progress", mutate: func(t *testing.T, root string) {
			marker := `{"schema_version":1,"transaction_id":"0123456789abcdef0123456789abcdef","command":"upgrade","phase":"staged","pid":999999,"started_at":"2026-08-08T00:00:00Z","target_root":"/fixture","recovery_path":".codeheart/local/kit-transactions/0123456789abcdef0123456789abcdef"}`
			if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(state.TransactionPath)), []byte(marker), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "recovery required", mutate: func(t *testing.T, root string) {
			marker := `{"schema_version":1,"transaction_id":"0123456789abcdef0123456789abcdef","command":"upgrade","phase":"recovery-required","pid":999999,"started_at":"2026-08-08T00:00:00Z","target_root":"/fixture","recovery_path":".codeheart/local/kit-transactions/0123456789abcdef0123456789abcdef"}`
			if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(state.TransactionPath)), []byte(marker), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
	}
	if runtime.GOOS != "windows" {
		tests = append(tests, struct {
			name   string
			mutate func(*testing.T, string)
		}{name: "symlinked managed path", mutate: func(t *testing.T, root string) {
			lock, _ := lockfile.ReadLock(root)
			record := state.Map(state.AnySlice(lock["managed_paths"])[0])
			path := filepath.Join(root, filepath.FromSlash(state.AsString(record["path"])))
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(filepath.Join(root, "AGENTS.md"), path); err != nil {
				t.Fatal(err)
			}
		}})
	} else {
		tests = append(tests, struct {
			name   string
			mutate func(*testing.T, string)
		}{name: "junctioned managed root", mutate: func(t *testing.T, root string) {
			kitRoot := filepath.Join(root, ".codeheart", "kit")
			outside := filepath.Join(t.TempDir(), "legacy-kit")
			if err := os.Rename(kitRoot, outside); err != nil {
				t.Fatal(err)
			}
			output, err := exec.Command("cmd", "/c", "mklink", "/J", kitRoot, outside).CombinedOutput()
			if err != nil {
				t.Fatalf("create junction: %v: %s", err, output)
			}
		}})
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := materializeLegacyV119Consumer(t)
			test.mutate(t, root)
			for _, dryRun := range []bool{true, false} {
				mode := "apply"
				if dryRun {
					mode = "dry-run"
				}
				t.Run(mode, func(t *testing.T) {
					before := snapshotTree(t, root)
					_, result, err := upgradeOperation(root, "0.1.30", filepath.Join(t.TempDir(), "unused-catalog.json"), writePreviousBinary(t), dryRun)
					if err != nil || result.Status != reconcile.StatusBlocked {
						t.Fatalf("invalid legacy source status=%s err=%v result=%#v", result.Status, err, result)
					}
					if !equalTreeSnapshot(before, snapshotTree(t, root)) {
						t.Fatal("rejected legacy source was mutated")
					}
				})
			}
		})
	}
}

func TestLegacyV1UpgradeForwardAndTargetAuthorityPreconditions(t *testing.T) {
	root := materializeLegacyV119Consumer(t)
	before := snapshotTree(t, root)
	for _, target := range []string{"0.1.19", "0.1.18", "invalid"} {
		_, result, err := upgradeOperation(root, target, filepath.Join(t.TempDir(), "unused-catalog.json"), writePreviousBinary(t), true)
		if err != nil || result.Status != reconcile.StatusBlocked || !hasBlocker(result, "upgrade_direction_invalid") {
			t.Fatalf("target %s status=%s err=%v result=%#v", target, result.Status, err, result)
		}
	}
	authority, err := inspectUpgradeSource(mustInspect(t, root))
	if err != nil {
		t.Fatal(err)
	}
	base := reconcileArguments(root, authority)
	for name, mutate := range map[string]func([]string) []string{
		"wrong target":           func(args []string) []string { return replaceArgument(args, "--target-version", "0.1.29") },
		"wrong previous version": func(args []string) []string { return replaceArgument(args, "--previous-version", "0.1.18") },
		"wrong schema":           func(args []string) []string { return replaceArgument(args, "--previous-lock-schema", "2") },
		"wrong lock digest": func(args []string) []string {
			return replaceArgument(args, "--previous-lock-sha256", strings.Repeat("f", 64))
		},
		"malformed release digest": func(args []string) []string {
			return replaceArgument(args, "--catalog-sha256", "invalid")
		},
		"partial bound protocol": func(args []string) []string {
			return removeArguments(args, "--target-version")
		},
	} {
		t.Run(name, func(t *testing.T) {
			var stderr bytes.Buffer
			if code := RunUpgradeReconcile(mutate(append([]string(nil), base...)), &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "precondition failed") {
				t.Fatalf("target reconcile code=%d stderr=%q", code, stderr.String())
			}
			if !equalTreeSnapshot(before, snapshotTree(t, root)) {
				t.Fatal("target precondition failure mutated repository")
			}
		})
	}
	var legacyStderr bytes.Buffer
	legacyArgs := removeArguments(append([]string(nil), base...), "--target-version", "--previous-lock-schema", "--previous-lock-sha256")
	if code := RunUpgradeReconcile(legacyArgs, &bytes.Buffer{}, &legacyStderr); code != 1 || !strings.Contains(legacyStderr.String(), "precondition failed") {
		t.Fatalf("legacy protocol accepted lock-v1 source: code=%d stderr=%q", code, legacyStderr.String())
	}
	if !equalTreeSnapshot(before, snapshotTree(t, root)) {
		t.Fatal("legacy protocol rejection mutated lock-v1 repository")
	}
}

func TestLegacyV1UpgradeSourceAuthorityCheckDetectsChange(t *testing.T) {
	root := materializeLegacyV119Consumer(t)
	observed := mustInspect(t, root)
	authority, err := inspectUpgradeSource(observed)
	if err != nil {
		t.Fatal(err)
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatal(err)
	}
	state.Map(lock["update_check"])["latest_seen_version"] = "0.1.20"
	if err := lockfile.WriteLock(root, lock); err != nil {
		t.Fatal(err)
	}
	blockers, err := boundUpgradeSourceAuthorityCheck(observed, authority)()
	if err != nil || len(blockers) != 1 || blockers[0].Code != "upgrade_source_lock_changed" {
		t.Fatalf("authority check blockers=%#v err=%v", blockers, err)
	}
}

func TestLegacyV1UpgradeTargetTransactionRollsBackOnAuthorityChange(t *testing.T) {
	root := materializeLegacyV119Consumer(t)
	observed := mustInspect(t, root)
	graph, err := compileObservedGraph(observed)
	if err != nil {
		t.Fatal(err)
	}
	digest := strings.Repeat("a", 64)
	now := time.Date(2026, 8, 8, 0, 0, 0, 0, time.UTC)
	lock, err := desiredUpgradeLock(observed, graph, now, version.Version, "fixture.zip", "catalog.json", digest, digest, digest, digest, digest)
	if err != nil {
		t.Fatal(err)
	}
	before := snapshotTree(t, root)
	result, err := runLifecycle(lifecycleRequest{
		command: "upgrade", root: root, now: now, observed: observed, graph: graph,
		desiredLock: lock, ensureIgnore: true,
		authorityCheck: func() ([]reconcile.Blocker, error) {
			return []reconcile.Blocker{{Code: "upgrade_source_lock_changed", Message: "fixture authority changed"}}, nil
		},
	})
	if err != nil || result.Status != reconcile.StatusRolledBack || !result.Rollback.Attempted || !result.Rollback.Succeeded {
		t.Fatalf("target transaction result=%#v err=%v", result, err)
	}
	if !equalTreeSnapshot(before, snapshotTree(t, root)) {
		t.Fatal("target transaction authority failure changed repository bytes")
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.TransactionPath))); !os.IsNotExist(err) {
		t.Fatalf("target transaction rollback left marker: %v", err)
	}
}

func TestLegacyV1UpgradeTargetTransactionPostCheckAndRollbackFailure(t *testing.T) {
	for _, test := range []struct {
		name             string
		failRollback     bool
		expectedStatus   reconcile.Status
		expectedRestored bool
	}{
		{name: "post-check failure rolls back", expectedStatus: reconcile.StatusRolledBack, expectedRestored: true},
		{name: "rollback failure requires recovery", failRollback: true, expectedStatus: reconcile.StatusRecoveryRequired},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := materializeLegacyV119Consumer(t)
			observed := mustInspect(t, root)
			authority, authorityErr := inspectUpgradeSource(observed)
			if authorityErr != nil {
				t.Fatal(authorityErr)
			}
			graph, err := compileObservedGraph(observed)
			if err != nil {
				t.Fatal(err)
			}
			digest := strings.Repeat("a", 64)
			now := time.Date(2026, 8, 8, 0, 0, 0, 0, time.UTC)
			lock, err := desiredUpgradeLock(observed, graph, now, version.Version, "fixture.zip", "catalog.json", digest, digest, digest, digest, digest)
			if err != nil {
				t.Fatal(err)
			}
			before := snapshotTree(t, root)
			rollbackFailureInjected := false
			result, err := runLifecycle(lifecycleRequest{
				command: "upgrade", root: root, now: now, observed: observed, graph: graph,
				desiredLock: lock, ensureIgnore: true,
				hook: func(phase string) error {
					if phase == "post-check" {
						return fmt.Errorf("fixture post-check failure")
					}
					if test.failRollback && !rollbackFailureInjected && strings.HasPrefix(phase, "rollback-quarantined:") {
						rollbackFailureInjected = true
						return fmt.Errorf("fixture rollback failure")
					}
					return nil
				},
			})
			if err != nil || result.Status != test.expectedStatus || !result.Rollback.Attempted || result.Rollback.Succeeded != test.expectedRestored {
				t.Fatalf("transaction result=%#v err=%v", result, err)
			}
			if test.expectedRestored {
				if !equalTreeSnapshot(before, snapshotTree(t, root)) {
					t.Fatal("post-check rollback did not restore repository bytes")
				}
				if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.TransactionPath))); !os.IsNotExist(err) {
					t.Fatalf("successful rollback left marker: %v", err)
				}
			} else {
				observedAfter, inspectErr := state.Inspect(root)
				if inspectErr != nil || observedAfter.Classification != state.StateRecoveryRequired || !hasBlocker(result, "recovery_required") || upgradeSourceRestored(root, authority) {
					t.Fatalf("rollback failure evidence observed=%#v err=%v result=%#v", observedAfter, inspectErr, result)
				}
			}
		})
	}
}

func TestLegacyV1UpgradeFailedTargetReconcileRestoresBinaryAndState(t *testing.T) {
	root := materializeLegacyV119Consumer(t)
	installed := writePreviousBinary(t)
	catalog, _ := writeUpgradePack(t)
	treeBefore := snapshotTree(t, root)
	binaryBefore := mustRead(t, installed)
	t.Setenv("CODEHEART_UPGRADE_TEST_FAIL_RECONCILE", "1")
	signal := filepath.Join(t.TempDir(), "target-reconcile-failed")
	t.Setenv("CODEHEART_UPGRADE_TEST_FAILURE_SIGNAL", signal)
	code, output := runApprovedUpgrade(t, []string{root, "--version", "0.1.30", "--catalog", catalog, "--installed-binary", installed, "--yes", "--json"})
	if runtime.GOOS == "windows" {
		if code != 0 {
			t.Fatalf("deferred failed reconcile scheduling exit=%d output=%s", code, output)
		}
		waitForFile(t, signal)
		waitForRestoredBinary(t, installed, binaryBefore)
	} else if code != 1 {
		t.Fatalf("failed reconcile exit=%d output=%s", code, output)
	}
	if runtime.GOOS != "windows" {
		var payload map[string]any
		if err := json.Unmarshal(output, &payload); err != nil {
			t.Fatal(err)
		}
		result := state.Map(payload["result"])
		rollback := state.Map(result["rollback"])
		if result["status"] != "rolled-back" || rollback["attempted"] != true || rollback["succeeded"] != true {
			t.Fatalf("rollback result = %#v", result)
		}
	}
	if !equalTreeSnapshot(treeBefore, snapshotTree(t, root)) || !bytes.Equal(binaryBefore, mustRead(t, installed)) {
		t.Fatal("failed reconcile did not restore prior binary and repository state")
	}
}

func TestLegacyV1UpgradeRecoveryRequiredTargetIsNotReportedRolledBack(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the initiating Windows process exits before deferred target reconciliation")
	}
	root := materializeLegacyV119Consumer(t)
	installed := writePreviousBinary(t)
	catalog, _ := writeUpgradePack(t)
	treeBefore := snapshotTree(t, root)
	binaryBefore := mustRead(t, installed)
	t.Setenv("CODEHEART_UPGRADE_TEST_RECOVERY_RECONCILE", "1")
	code, output := runApprovedUpgrade(t, []string{root, "--version", "0.1.30", "--catalog", catalog, "--installed-binary", installed, "--yes", "--json"})
	if code != 1 {
		t.Fatalf("recovery-required reconcile exit=%d output=%s", code, output)
	}
	var payload map[string]any
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatal(err)
	}
	result := state.Map(payload["result"])
	rollback := state.Map(result["rollback"])
	if result["status"] != "recovery-required" || rollback["attempted"] != true || rollback["succeeded"] != false {
		t.Fatalf("recovery-required result = %#v", result)
	}
	if !equalTreeSnapshot(treeBefore, snapshotTree(t, root)) || !bytes.Equal(binaryBefore, mustRead(t, installed)) {
		t.Fatal("typed recovery-required handoff did not restore source bytes")
	}
}

func TestCurrentLockV2UpgradeBehaviorAndNoLegacyRepairSyncBypass(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "V2 Upgrade"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("init=%d", code)
	}
	lock, _ := lockfile.ReadLock(root)
	lock["kit_version"] = "0.1.25"
	if err := lockfile.WriteLock(root, lock); err != nil {
		t.Fatal(err)
	}
	configBefore := mustRead(t, filepath.Join(root, filepath.FromSlash(state.ConfigPath)))
	installed := writePreviousBinary(t)
	catalog, _ := writeUpgradePack(t)
	if code, output := runApprovedUpgrade(t, []string{root, "--version", "0.1.30", "--catalog", catalog, "--installed-binary", installed, "--yes", "--json"}); code != 0 {
		t.Fatalf("v2 upgrade exit=%d output=%s", code, output)
	}
	waitForUpgrade(t, root, installed)
	upgraded, _ := lockfile.ReadLock(root)
	if state.AsInt(upgraded["schema_version"]) != 2 || state.AsString(upgraded["kit_version"]) != "0.1.30" || state.AsInt(upgraded["state_generation"]) != 2 || state.AsInt(state.Map(upgraded["last_operation"])["previous_generation"]) != 1 {
		t.Fatalf("v2 upgraded lock = %#v", upgraded)
	}
	if !bytes.Equal(configBefore, mustRead(t, filepath.Join(root, filepath.FromSlash(state.ConfigPath)))) {
		t.Fatal("v2 upgrade changed config")
	}

	legacy := materializeLegacyV119Consumer(t)
	before := snapshotTree(t, legacy)
	_, repairResult, err := repairOperation(legacy, time.Now(), false)
	if err != nil || repairResult.Status != reconcile.StatusBlocked || !hasBlocker(repairResult, "version_change_requires_upgrade") {
		t.Fatalf("repair bypass result=%#v err=%v", repairResult, err)
	}
	_, syncResult, err := syncOperation(legacy, "", time.Now(), false)
	if err != nil || syncResult.Status != reconcile.StatusBlocked {
		t.Fatalf("sync bypass result=%#v err=%v", syncResult, err)
	}
	if !equalTreeSnapshot(before, snapshotTree(t, legacy)) {
		t.Fatal("repair or sync changed cross-version legacy state")
	}
}

func TestCurrentLockV2TargetAllowsOnlyTargetRelativeAdditions(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "V2 Target Addition"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("init=%d", code)
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatal(err)
	}
	lock["kit_version"] = "0.1.25"
	managed := []any{}
	for _, item := range state.AnySlice(lock["managed_paths"]) {
		if state.AsString(state.Map(item)["path"]) != legacyMissingManagedPath {
			managed = append(managed, item)
		}
	}
	lock["managed_paths"] = managed
	if err := lockfile.WriteLock(root, lock); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(legacyMissingManagedPath))); err != nil {
		t.Fatal(err)
	}
	observed := mustInspect(t, root)
	if observed.Classification != state.StatePartial {
		t.Fatalf("target-relative addition state=%s missing=%#v", observed.Classification, observed.MissingPaths)
	}
	if _, err := inspectUpgradeSource(observed); err == nil {
		t.Fatal("public lock-v2 preflight accepted a partial installation")
	}
	authority, err := inspectTargetUpgradeSource(observed)
	if err != nil {
		t.Fatalf("target-side lock-v2 validation rejected a target-relative addition: %v", err)
	}
	args := replaceArgument(reconcileArguments(root, authority), "--previous-version", "0.1.25")
	args = removeArguments(args, "--target-version", "--previous-lock-schema", "--previous-lock-sha256")
	var stderr bytes.Buffer
	if code := RunUpgradeReconcile(args, &bytes.Buffer{}, &stderr); code != 0 {
		t.Fatalf("target reconcile exit=%d stderr=%s", code, stderr.String())
	}
	upgraded := mustInspect(t, root)
	if upgraded.Classification != state.StateCurrent || state.AsString(upgraded.Lock["kit_version"]) != version.Version {
		t.Fatalf("target-relative v2 upgrade state=%#v", upgraded)
	}
}

func TestCurrentLockV2TargetRejectsMissingDeclaredPath(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "V2 Missing Declared Path"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("init=%d", code)
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatal(err)
	}
	lock["kit_version"] = "0.1.25"
	if err := lockfile.WriteLock(root, lock); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(legacyMissingManagedPath))); err != nil {
		t.Fatal(err)
	}
	observed := mustInspect(t, root)
	if observed.Classification != state.StatePartial {
		t.Fatalf("missing declared path state=%s", observed.Classification)
	}
	if _, err := inspectTargetUpgradeSource(observed); err == nil {
		t.Fatal("target-side lock-v2 validation accepted a missing lock-declared path")
	}
}

func TestCurrentLockV2TargetRejectsTargetAddedManagedOverlap(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "V2 Target Overlap"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("init=%d", code)
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatal(err)
	}
	lock["kit_version"] = "0.1.25"
	managed := []any{}
	for _, item := range state.AnySlice(lock["managed_paths"]) {
		if state.AsString(state.Map(item)["path"]) != legacyMissingManagedPath {
			managed = append(managed, item)
		}
	}
	lock["managed_paths"] = managed
	if err := lockfile.WriteLock(root, lock); err != nil {
		t.Fatal(err)
	}
	overlapPath := filepath.Join(root, filepath.FromSlash(legacyMissingManagedPath))
	if err := os.WriteFile(overlapPath, []byte("consumer-owned overlap\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := snapshotTree(t, root)
	observed := mustInspect(t, root)
	if observed.Classification != state.StateDrifted {
		t.Fatalf("target overlap state=%s", observed.Classification)
	}
	if _, err := inspectUpgradeSource(observed); err == nil {
		t.Fatal("public lock-v2 validation accepted target-added managed overlap")
	}
	lockAuthority, err := readUpgradeSourceAuthority(observed.CanonicalRoot, 2)
	if err != nil {
		t.Fatal(err)
	}
	args := replaceArgument(reconcileArguments(root, lockAuthority), "--previous-version", "0.1.25")
	var stderr bytes.Buffer
	if code := RunUpgradeReconcile(args, &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "precondition failed") {
		t.Fatalf("target overlap reconcile exit=%d stderr=%s", code, stderr.String())
	}
	if !equalTreeSnapshot(before, snapshotTree(t, root)) {
		t.Fatal("target-added lock-v2 overlap rejection mutated repository bytes")
	}
}

const legacyMissingManagedPath = ".codeheart/kit/docs/agent-interface/runbooks/maintain-operating-kit-installation.md"

func materializeLegacyV119Consumer(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	skeletonData := mustRead(t, filepath.Join("..", "..", "tests", "fixtures", "state", "lock-v1-v0.1.19.yaml"))
	skeleton, err := state.DecodeAndValidateYAML(state.LockV1Schema, skeletonData)
	if err != nil {
		t.Fatal(err)
	}
	configData := mustRead(t, filepath.Join("..", "..", "tests", "fixtures", "state", "config-v1-v0.1.19.yaml"))
	config, err := state.DecodeAndValidateYAML(state.ConfigV1Schema, configData)
	if err != nil {
		t.Fatal(err)
	}
	config["selected_setup_folder"] = root
	encodedConfig, err := state.EncodeYAML(config)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range state.AnySlice(skeleton["managed_paths"]) {
		record := state.Map(item)
		target := state.AsString(record["path"])
		content := legacyFixtureManagedContent(target)
		if shaBytes(content) != state.AsString(record["checksum_sha256"]) {
			t.Fatalf("legacy fixture checksum does not match generic content for %s", target)
		}
		path := filepath.Join(root, filepath.FromSlash(target))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, ".codeheart"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(state.ConfigPath)), encodedConfig, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := lockfile.WriteLock(root, skeleton); err != nil {
		t.Fatal(err)
	}

	for relative, content := range map[string]string{
		"docs/repo/README.md":                          "legacy consumer repository documentation\n",
		"docs/repo/plans/plan-register.md":             "consumer plan register\n",
		"docs/repo/plans/coordination-sync-pending.md": "legacy consumer coordination surface\n",
		"docs/agent-memory/README.md":                  "legacy consumer memory documentation\n",
		"docs/agent-memory/goal-register.md":           "legacy consumer goal register\n",
		"docs/agent-memory/session-ledger.md":          "consumer memory ledger\n",
		"docs/agent-memory/untriaged-sessions.md":      "legacy consumer untriaged sessions\n",
		".codeheart/user/README.md":                    "legacy consumer user layer\n",
		".codeheart/user/examples/preferences.yaml":    "language: en\n",
		".codeheart/user/preferences.yaml":             "language: de\n",
	} {
		path := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	agentsPath := filepath.Join(root, "AGENTS.md")
	agents := "# Legacy Consumer Instructions\n\n" + reconcile.BeginMarker + "\n\nlegacy managed guidance\n\n" + reconcile.EndMarker + "\n\n# Consumer-Owned Guidance\n\npreserve this suffix\n"
	if err := os.WriteFile(agentsPath, []byte(agents), 0o644); err != nil {
		t.Fatal(err)
	}
	observed := mustInspect(t, root)
	if observed.LockSchemaVersion != 1 || observed.Classification != state.StatePartial {
		t.Fatalf("legacy fixture classification = %#v", observed)
	}
	if _, err := inspectUpgradeSource(observed); err != nil {
		t.Fatalf("legacy fixture is not upgrade-compatible: %v", err)
	}
	return root
}

func legacyFixtureManagedContent(target string) []byte {
	return []byte("legacy v0.1.19 managed fixture: " + target + "\n")
}

type packEvidence map[string]string

func writeUpgradePack(t *testing.T) (string, packEvidence) {
	t.Helper()
	root := t.TempDir()
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	binary := mustRead(t, executable)
	platform, err := releasekit.PlatformForRuntime(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		t.Skip(err)
	}
	binaryPath := "bin/codeheart-operating-kit"
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}
	content := mustRead(t, filepath.Join("..", "..", "manifest.yaml"))
	files := map[string][]byte{binaryPath: binary, "content-manifest.yaml": content}
	paths := []string{binaryPath, "content-manifest.yaml"}
	sort.Strings(paths)
	var checksumText strings.Builder
	for _, name := range paths {
		fmt.Fprintf(&checksumText, "%s  %s\n", shaBytes(files[name]), name)
	}
	files["checksums.txt"] = []byte(checksumText.String())
	manifest := releasekit.PackManifest{
		SchemaVersion: 1, Version: version.Version, Platform: platform, Command: "codeheart-operating-kit",
		BinaryPath: binaryPath, BinarySHA256: shaBytes(binary), ContentManifestPath: "content-manifest.yaml",
		ContentManifestSHA256: shaBytes(content), PayloadChecksumsPath: "checksums.txt", PayloadChecksumsSHA256: shaBytes(files["checksums.txt"]),
	}
	manifestData, _ := json.MarshalIndent(manifest, "", "  ")
	manifestData = append(manifestData, '\n')
	files["pack-manifest.json"] = manifestData
	archive := filepath.Join(root, fmt.Sprintf("codeheart-operating-kit-%s-%s.zip", version.Version, platform))
	archiveFile, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(archiveFile)
	archivePaths := make([]string, 0, len(files))
	for name := range files {
		archivePaths = append(archivePaths, name)
	}
	sort.Strings(archivePaths)
	payloadRoot := fmt.Sprintf("codeheart-operating-kit-%s-%s", version.Version, platform)
	for _, name := range archivePaths {
		entry, err := writer.Create(payloadRoot + "/" + name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write(files[name]); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := archiveFile.Close(); err != nil {
		t.Fatal(err)
	}
	asset := releasekit.CatalogAsset{Name: filepath.Base(archive), Version: version.Version, Platform: platform, URL: archive, ArchiveSHA256: shaFile(t, archive), PackManifestSHA256: shaBytes(manifestData)}
	catalog := releasekit.Catalog{SchemaVersion: 1, Version: version.Version, Assets: []releasekit.CatalogAsset{asset}}
	catalogData, _ := json.MarshalIndent(catalog, "", "  ")
	catalogData = append(catalogData, '\n')
	catalogPath := filepath.Join(root, "release-catalog-"+version.Version+".json")
	if err := os.WriteFile(catalogPath, catalogData, 0o644); err != nil {
		t.Fatal(err)
	}
	return catalogPath, packEvidence{
		"catalog_sha256": shaBytes(catalogData), "archive_sha256": asset.ArchiveSHA256,
		"pack_manifest_sha256": asset.PackManifestSHA256, "content_manifest_sha256": manifest.ContentManifestSHA256,
		"binary_sha256": manifest.BinarySHA256,
	}
}

func writePreviousBinary(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "codeheart-operating-kit")
	if runtime.GOOS == "windows" {
		path += ".exe"
	}
	if err := os.WriteFile(path, []byte("previous binary\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func runApprovedUpgrade(t *testing.T, args []string) (int, []byte) {
	t.Helper()
	if runtime.GOOS != "windows" {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunUpgrade(args, &stdout, &stderr)
		return code, append(stdout.Bytes(), stderr.Bytes()...)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	commandArgs := append([]string{"__test-upgrade-apply"}, args...)
	output, err := exec.Command(executable, commandArgs...).CombinedOutput()
	if err == nil {
		return 0, output
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), output
	}
	t.Fatalf("start short-lived upgrade helper: %v", err)
	return -1, output
}

func waitForFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(35 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("target reconcile failure signal was not observed")
}

func waitForRestoredBinary(t *testing.T, installed string, expected []byte) {
	t.Helper()
	deadline := time.Now().Add(35 * time.Second)
	for time.Now().Before(deadline) {
		if actual, err := os.ReadFile(installed); err == nil && bytes.Equal(actual, expected) {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("deferred handoff did not restore the previous binary")
}

func waitForUpgrade(t *testing.T, root, installed string) {
	t.Helper()
	deadline := time.Now().Add(35 * time.Second)
	for time.Now().Before(deadline) {
		output, err := exec.Command(installed, "--version").CombinedOutput()
		observed, inspectErr := state.Inspect(root)
		if err == nil && strings.TrimSpace(string(output)) == "codeheart-operating-kit "+version.Version && inspectErr == nil && observed.Classification == state.StateCurrent {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("upgrade handoff did not reach current state")
}

func reconcileArguments(root string, authority upgradeSourceAuthority) []string {
	digest := strings.Repeat("a", 64)
	return []string{
		"--repository", root, "--target-version", version.Version, "--previous-version", "0.1.19",
		"--previous-lock-schema", fmt.Sprint(authority.LockSchema), "--previous-lock-sha256", authority.LockSHA256,
		"--asset-url", "fixture.zip", "--catalog-location", "catalog.json", "--catalog-sha256", digest,
		"--archive-sha256", digest, "--pack-manifest-sha256", digest, "--content-manifest-sha256", digest, "--binary-sha256", digest,
	}
}

func replaceArgument(args []string, name, value string) []string {
	for index := 0; index+1 < len(args); index++ {
		if args[index] == name {
			args[index+1] = value
			return args
		}
	}
	return args
}

func removeArguments(args []string, names ...string) []string {
	remove := map[string]bool{}
	for _, name := range names {
		remove[name] = true
	}
	result := make([]string, 0, len(args))
	for index := 0; index < len(args); index++ {
		if remove[args[index]] && index+1 < len(args) {
			index++
			continue
		}
		result = append(result, args[index])
	}
	return result
}

func mustInspect(t *testing.T, root string) state.Observed {
	t.Helper()
	observed, err := state.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	return observed
}

type treeEntry struct {
	Mode fs.FileMode
	Data string
}

func snapshotTree(t *testing.T, root string) map[string]treeEntry {
	t.Helper()
	result := map[string]treeEntry{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil || relative == "." {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		item := treeEntry{Mode: info.Mode()}
		if info.Mode().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			item.Data = string(data)
		} else if info.Mode()&os.ModeSymlink != 0 {
			item.Data, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}
		result[filepath.ToSlash(relative)] = item
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func equalTreeSnapshot(left, right map[string]treeEntry) bool {
	if len(left) != len(right) {
		return false
	}
	for path, entry := range left {
		if right[path] != entry {
			return false
		}
	}
	return true
}

func resultHasChange(result map[string]any, path, action string) bool {
	for _, item := range state.AnySlice(result["changes"]) {
		change := state.Map(item)
		if state.AsString(change["path"]) == path && state.AsString(change["action"]) == action {
			return true
		}
	}
	return false
}

func resultHasValidation(result map[string]any, name string) bool {
	for _, item := range state.AnySlice(result["validations"]) {
		if state.AsString(state.Map(item)["name"]) == name {
			return true
		}
	}
	return false
}

func hasBlocker(result reconcile.Result, code string) bool {
	for _, blocker := range result.Blockers {
		if blocker.Code == code {
			return true
		}
	}
	return false
}

func shaBytes(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
func shaFile(t *testing.T, path string) string { t.Helper(); return shaBytes(mustRead(t, path)) }
