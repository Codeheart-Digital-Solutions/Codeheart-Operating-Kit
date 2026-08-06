package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/lockfile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/plancatalog"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/portfolio"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/yamlmini"
)

func TestInspectFolderModes(t *testing.T) {
	root := t.TempDir()
	if got := InspectFolder(filepath.Join(root, "new"))["mode"]; got != "new-folder-setup" {
		t.Fatalf("missing folder mode = %v", got)
	}
	if got := InspectFolder(root)["mode"]; got != "new-folder-setup" {
		t.Fatalf("empty folder mode = %v", got)
	}
	if err := os.WriteFile(filepath.Join(root, "pyproject.toml"), []byte("[project]\nname = 'x'\n"), 0o644); err != nil {
		t.Fatalf("write marker: %v", err)
	}
	technical := InspectFolder(root)
	if technical["mode"] != "existing-technical-project-adoption" {
		t.Fatalf("technical folder mode = %v", technical["mode"])
	}
	markers := technical["markers"].([]any)
	if len(markers) != 1 || markers[0] != "pyproject.toml" {
		t.Fatalf("technical markers = %#v", markers)
	}
}

func TestLifecycleDryRunRepairIdempotencyAndConsumerPreservation(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing-target")
	var preview bytes.Buffer
	if code := RunInit([]string{root, "--project-name", "Example", "--dry-run", "--json"}, &preview, &bytes.Buffer{}); code != 0 {
		t.Fatalf("init dry-run exit = %d; %s", code, preview.String())
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("init dry-run created target: %v", err)
	}
	var previewPayload map[string]any
	if err := json.Unmarshal(preview.Bytes(), &previewPayload); err != nil {
		t.Fatal(err)
	}
	if state.Map(previewPayload["result"])["status"] != "planned" {
		t.Fatalf("preview result = %#v", previewPayload["result"])
	}

	if code := RunInit([]string{root, "--project-name", "Example", "--purpose", "company-automation"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("init exit = %d", code)
	}
	configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	configBefore, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	preserved := map[string]string{
		".codeheart/user/preferences.yaml":    "language: de\n",
		"docs/repo/plans/plan-register.md":    "custom plan state\n",
		"docs/agent-memory/session-ledger.md": "custom memory state\n",
	}
	for relative, content := range preserved {
		path := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	managed := filepath.Join(root, ".codeheart", "kit", "docs", "agent-interface", "README.md")
	if err := os.WriteFile(managed, []byte("drift\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := RunRepair([]string{root, "--dry-run"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("repair dry-run exit = %d", code)
	}
	if data, _ := os.ReadFile(managed); string(data) != "drift\n" {
		t.Fatalf("repair dry-run changed managed file")
	}
	if code := RunRepair([]string{root}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("repair exit = %d", code)
	}
	configAfter, _ := os.ReadFile(configPath)
	if string(configAfter) != string(configBefore) {
		t.Fatalf("repair changed consumer config\nbefore:\n%s\nafter:\n%s", configBefore, configAfter)
	}
	for relative, content := range preserved {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil || string(data) != content {
			t.Fatalf("repair changed %s: %q err=%v", relative, data, err)
		}
	}
	var noOp bytes.Buffer
	if code := RunRepair([]string{root, "--json"}, &noOp, &bytes.Buffer{}); code != 0 {
		t.Fatalf("idempotent repair exit = %d; %s", code, noOp.String())
	}
	var noOpPayload map[string]any
	_ = json.Unmarshal(noOp.Bytes(), &noOpPayload)
	if changes := state.AnySlice(state.Map(noOpPayload["result"])["changes"]); len(changes) != 0 {
		t.Fatalf("idempotent repair changes = %#v", changes)
	}
	if code := RunInit([]string{root, "--project-name", "Again"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 1 {
		t.Fatalf("second init exit = %d, want blocked", code)
	}
}

func TestRepairMigratesOnlyCompatibleLockV1(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "Example"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("init exit = %d", code)
	}
	v2, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatal(err)
	}
	v1 := map[string]any{
		"schema_version":      1,
		"kit_version":         v2["kit_version"],
		"selected_profile":    v2["selected_profile"],
		"selected_components": v2["selected_components"],
		"release":             state.DeepCopy(state.Map(v2["release"])),
		"managed_paths":       v2["managed_paths"],
		"generated_surfaces":  v2["generated_surfaces"],
		"cli_repair":          state.DeepCopy(state.Map(v2["cli_repair"])),
		"update_check":        v2["update_check"],
		"native_capabilities": v2["native_capabilities"],
	}
	legacySurfaces := []any{}
	for _, item := range state.AnySlice(v1["generated_surfaces"]) {
		if state.AsString(state.Map(item)["ownership"]) != "local-machine" {
			legacySurfaces = append(legacySurfaces, item)
		}
	}
	v1["generated_surfaces"] = legacySurfaces
	state.Map(v1["release"])["checksum_sha256"] = 0
	state.Map(v1["cli_repair"])["repair_checksum_sha256"] = "0"
	data, err := state.EncodeYAML(v1)
	if err != nil {
		t.Fatal(err)
	}
	lockPath := filepath.Join(root, filepath.FromSlash(state.LockPath))
	if err := os.WriteFile(lockPath, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if code := RunRepair([]string{root, "--dry-run"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("migration preview exit = %d", code)
	}
	stillV1, _ := state.DecodeYAMLMap(mustRead(t, lockPath))
	if state.AsInt(stillV1["schema_version"]) != 1 {
		t.Fatalf("dry-run migrated lock on disk")
	}
	if code := RunRepair([]string{root}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("migration repair exit = %d", code)
	}
	migrated, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatal(err)
	}
	if state.AsInt(migrated["schema_version"]) != 2 || state.AsString(state.Map(migrated["release_provenance"])["verification_status"]) != "unverified-legacy" {
		t.Fatalf("migrated lock = %#v", migrated)
	}
}

func TestUpdateCheckChangesNoRepositoryBytesExceptLock(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "Example"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("init exit = %d", code)
	}
	paths := []string{state.ConfigPath, "AGENTS.md", ".codeheart/kit/README.md", "docs/repo/plans/plan-register.md"}
	before := map[string][]byte{}
	for _, relative := range paths {
		before[relative] = mustRead(t, filepath.Join(root, filepath.FromSlash(relative)))
	}
	if _, err := UpdateCheck(root, "9.0.0", time.Date(2026, 7, 9, 20, 0, 0, 0, time.UTC).Format(time.RFC3339), ""); err != nil {
		t.Fatal(err)
	}
	for _, relative := range paths {
		after := mustRead(t, filepath.Join(root, filepath.FromSlash(relative)))
		if string(after) != string(before[relative]) {
			t.Fatalf("update-check changed %s", relative)
		}
	}
}

func TestOnboardRoutesExistingInstallationToRepair(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "Example"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("init exit = %d", code)
	}
	managed := filepath.Join(root, ".codeheart", "kit", "docs", "agent-interface", "README.md")
	if err := os.WriteFile(managed, []byte("drift\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, _, code, err := Onboard(root, "Example", "", true, time.Now())
	if err != nil || code != 0 || result["written"] != true {
		t.Fatalf("existing onboard = %#v, code=%d, err=%v", result, code, err)
	}
	observed, err := state.Inspect(root)
	if err != nil || observed.Classification != state.StateCurrent {
		t.Fatalf("onboard repair state = %#v, err=%v", observed, err)
	}
}

func TestLifecycleStartingStatePreconditionMatrix(t *testing.T) {
	t.Setenv("CODEHEART_OPERATING_KIT_CLI", "1")
	type expected struct {
		init, repair, sync, update, healthy bool
	}
	cases := map[string]expected{
		"absent":      {init: true},
		"adoptable":   {init: true},
		"current":     {repair: true, sync: true, update: true, healthy: true},
		"drifted":     {repair: true, sync: true, update: true},
		"stale-cli":   {},
		"partial":     {repair: true},
		"invalid":     {},
		"future":      {},
		"transaction": {},
		"recovery":    {},
	}
	now := time.Date(2026, 7, 9, 20, 0, 0, 0, time.UTC)
	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			root := lifecycleStateFixture(t, name)
			_, initResult, err := initializeOperation(root, "Example", "", root, now, true)
			if err != nil || initResult.OK() != want.init {
				t.Fatalf("init OK=%v want=%v err=%v result=%#v", initResult.OK(), want.init, err, initResult)
			}
			_, repairResult, err := repairOperation(root, now, true)
			if err != nil || repairResult.OK() != want.repair {
				t.Fatalf("repair OK=%v want=%v err=%v result=%#v", repairResult.OK(), want.repair, err, repairResult)
			}
			_, syncResult, err := syncOperation(root, "", now, true)
			if err != nil || syncResult.OK() != want.sync {
				t.Fatalf("sync OK=%v want=%v err=%v result=%#v", syncResult.OK(), want.sync, err, syncResult)
			}
			_, updateResult, err := updateCheckOperation(root, "0.1.24", now.Format(time.RFC3339), "", true)
			if err != nil || updateResult.OK() != want.update {
				t.Fatalf("update OK=%v want=%v err=%v result=%#v", updateResult.OK(), want.update, err, updateResult)
			}
			check, err := CheckRepository(root)
			if err != nil || (check["ok"] == true) != want.healthy {
				t.Fatalf("check healthy=%v want=%v err=%v result=%#v", check["ok"], want.healthy, err, check)
			}
		})
	}
}

func lifecycleStateFixture(t *testing.T, kind string) string {
	t.Helper()
	base := t.TempDir()
	if kind == "absent" {
		return filepath.Join(base, "missing")
	}
	if kind == "adoptable" {
		return base
	}
	if code := RunInit([]string{base, "--project-name", "Example"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("fixture init exit = %d", code)
	}
	switch kind {
	case "current":
	case "drifted":
		path := filepath.Join(base, ".codeheart", "kit", "docs", "agent-interface", "README.md")
		if err := os.WriteFile(path, []byte("drift\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	case "stale-cli":
		lock, err := lockfile.ReadLock(base)
		if err != nil {
			t.Fatal(err)
		}
		lock["kit_version"] = "0.0.1"
		if err := lockfile.WriteLock(base, lock); err != nil {
			t.Fatal(err)
		}
	case "partial":
		if err := os.Remove(filepath.Join(base, "AGENTS.md")); err != nil {
			t.Fatal(err)
		}
	case "invalid":
		if err := os.WriteFile(filepath.Join(base, filepath.FromSlash(state.LockPath)), []byte("schema_version: invalid\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	case "future":
		if err := os.WriteFile(filepath.Join(base, filepath.FromSlash(state.LockPath)), []byte("schema_version: 99\nkit_version: future\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	case "transaction", "recovery":
		phase := "staged"
		if kind == "recovery" {
			phase = "recovery-required"
		}
		marker := fmt.Sprintf("{\"schema_version\":1,\"transaction_id\":\"fixture\",\"command\":\"sync\",\"phase\":%q,\"pid\":%d}\n", phase, os.Getpid())
		if err := os.WriteFile(filepath.Join(base, filepath.FromSlash(state.TransactionPath)), []byte(marker), 0o600); err != nil {
			t.Fatal(err)
		}
	default:
		t.Fatalf("unknown state fixture %s", kind)
	}
	return base
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestRunInspectJSON(t *testing.T) {
	var stdout bytes.Buffer
	code := RunInspect([]string{"--json", "."}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunInspect exit = %d", code)
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse inspect JSON: %v\n%s", err, stdout.String())
	}
	if payload["mode"] == "" {
		t.Fatalf("inspect JSON missing mode: %#v", payload)
	}
}

func TestPlansListJSONGolden(t *testing.T) {
	root := copyPlanCommandFixture(t)
	expected := mustRead(t, filepath.Join("testdata", "plans-list.json"))
	for _, args := range [][]string{{"list", "--format", "json", root}, {"list", "--json", root}} {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		if code := RunPlans(args, &stdout, &stderr); code != 0 {
			t.Fatalf("plans list %v exit=%d stderr=%s", args, code, stderr.String())
		}
		if stdout.String() != string(expected) {
			t.Fatalf("plans list %v output changed\nwant:\n%s\ngot:\n%s", args, expected, stdout.String())
		}
	}
}

func TestPlansListFormatValidation(t *testing.T) {
	root := copyPlanCommandFixture(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunPlans([]string{"list", "--format", "text", root}, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "Plans (mixed mode)") {
		t.Fatalf("text list exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	stderr.Reset()
	if code := RunPlans([]string{"list", "--format", "yaml", root}, &bytes.Buffer{}, &stderr); code != 2 || !strings.Contains(stderr.String(), "invalid --format") {
		t.Fatalf("invalid format exit=%d stderr=%s", code, stderr.String())
	}
	stderr.Reset()
	if code := RunPlans([]string{"list", "--format", "text", "--json", root}, &bytes.Buffer{}, &stderr); code != 2 || !strings.Contains(stderr.String(), "conflicts") {
		t.Fatalf("format conflict exit=%d stderr=%s", code, stderr.String())
	}
}

func TestPlansInventoryWritesOnlyTheExplicitArtifact(t *testing.T) {
	root := copyPlanCommandFixture(t)
	alpha := filepath.Join(root, "docs/repo/plans/alpha/alpha_discovery_doc.md")
	beta := filepath.Join(root, "docs/repo/plans/beta/beta_implementation_doc.md")
	beforeAlpha := mustRead(t, alpha)
	beforeBeta := mustRead(t, beta)
	destination := filepath.Join(root, "migration-inventory.yaml")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunPlans([]string{"inventory", "--output", destination, root}, &stdout, &stderr); code != 0 {
		t.Fatalf("plans inventory exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !exists(destination) || !bytes.Equal(beforeAlpha, mustRead(t, alpha)) || !bytes.Equal(beforeBeta, mustRead(t, beta)) {
		t.Fatal("inventory failed to preserve canonical plan bytes or explicit output")
	}
	if data := mustRead(t, destination); !bytes.Contains(data, []byte("source_revision:")) || !bytes.Contains(data, []byte("unpaired_legacy_evidence:")) {
		t.Fatalf("inventory artifact missing evidence:\n%s", data)
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunPlans([]string{"inventory", "--output", destination, root}, &stdout, &stderr); code != 0 {
		t.Fatalf("repeated plans inventory refresh exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if data := mustRead(t, destination); !bytes.Contains(data, []byte("source_revision:")) {
		t.Fatalf("repeated inventory refresh produced invalid output:\n%s", data)
	}
}

func TestPlansInventoryBoundParentRejectsConcurrentRedirection(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory symlink substitution regression is covered by Unix validation jobs")
	}
	root := copyPlanCommandFixture(t)
	inventory, err := plancatalog.BuildInventory(root, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	scratch := filepath.Join(root, "scratch")
	if err := os.Mkdir(scratch, 0o755); err != nil {
		t.Fatal(err)
	}
	formalDirectory := filepath.Join(root, "docs", "repo", "plans", "alpha")
	formalPlan := filepath.Join(formalDirectory, "alpha_discovery_doc.md")
	formalBefore := mustRead(t, formalPlan)
	destination := filepath.Join(scratch, "alpha_discovery_doc.md")
	ownedScratch := scratch + "-owned"
	err = writeInventoryArtifactWithHook(root, destination, inventory, func(phase string) error {
		if phase != "inventory-parent-bound" {
			return nil
		}
		if err := os.Rename(scratch, ownedScratch); err != nil {
			return err
		}
		return os.Symlink(formalDirectory, scratch)
	})
	if err == nil || !strings.Contains(err.Error(), "parent path identity changed") {
		t.Fatalf("redirected inventory parent error = %v", err)
	}
	if formalAfter := mustRead(t, formalPlan); !bytes.Equal(formalAfter, formalBefore) {
		t.Fatal("bound inventory writer changed redirected formal plan")
	}
	if _, err := os.Stat(filepath.Join(ownedScratch, "alpha_discovery_doc.md")); !os.IsNotExist(err) {
		t.Fatalf("failed bound inventory write retained an output target: %v", err)
	}
}

func TestPlansInventoryBindsParentBeforeProtectionChecks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("renaming an open directory is not portable to Windows")
	}
	root := copyPlanCommandFixture(t)
	inventory, err := plancatalog.BuildInventory(root, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	scratch := filepath.Join(root, "scratch")
	if err := os.Mkdir(scratch, 0o755); err != nil {
		t.Fatal(err)
	}
	formalDirectory := filepath.Join(root, "docs", "repo", "plans", "alpha")
	formalPlan := filepath.Join(formalDirectory, "alpha_discovery_doc.md")
	formalBefore := mustRead(t, formalPlan)
	destination := filepath.Join(scratch, "alpha_discovery_doc.md")
	ownedScratch := scratch + "-owned"
	err = writeInventoryArtifactWithHook(root, destination, inventory, func(phase string) error {
		if phase != "inventory-parent-bound-before-protection" {
			return nil
		}
		if err := os.Rename(scratch, ownedScratch); err != nil {
			return err
		}
		return os.Rename(formalDirectory, scratch)
	})
	if err == nil || !strings.Contains(err.Error(), "parent path identity changed") {
		t.Fatalf("exchanged inventory parent error = %v", err)
	}
	if formalAfter := mustRead(t, destination); !bytes.Equal(formalAfter, formalBefore) {
		t.Fatal("pre-protection parent exchange changed the formal plan")
	}
	if _, err := os.Stat(filepath.Join(ownedScratch, "alpha_discovery_doc.md")); !os.IsNotExist(err) {
		t.Fatalf("failed pre-protection write retained an output target: %v", err)
	}
}

func TestPlansInventoryRequiresExistingOutputParentWithoutCreatingAncestors(t *testing.T) {
	root := copyPlanCommandFixture(t)
	destinationParent := filepath.Join(root, "missing-output-parent", "nested")
	destination := filepath.Join(destinationParent, "inventory.yaml")
	var stderr bytes.Buffer
	if code := RunPlans([]string{"inventory", "--output", destination, root}, &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "inventory_target_parent_missing") {
		t.Fatalf("missing inventory parent exit=%d stderr=%s", code, stderr.String())
	}
	if _, err := os.Lstat(filepath.Join(root, "missing-output-parent")); !os.IsNotExist(err) {
		t.Fatalf("inventory created a missing output ancestor: %v", err)
	}
}

func TestPlansInventoryPreservesConcurrentOutputChanges(t *testing.T) {
	root := copyPlanCommandFixture(t)
	inventory, err := plancatalog.BuildInventory(root, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	t.Run("missing target appears", func(t *testing.T) {
		destination := filepath.Join(root, "concurrent-missing-inventory.yaml")
		concurrent := []byte("concurrent new output\n")
		err := writeInventoryArtifactWithHook(root, destination, inventory, func(phase string) error {
			if phase == "inventory-parent-bound" {
				return os.WriteFile(destination, concurrent, 0o600)
			}
			return nil
		})
		if err == nil {
			t.Fatal("inventory unexpectedly replaced a concurrently appearing output")
		}
		if actual := mustRead(t, destination); !bytes.Equal(actual, concurrent) {
			t.Fatalf("concurrently appearing output changed: %q", actual)
		}
	})
	t.Run("existing bytes change", func(t *testing.T) {
		destination := filepath.Join(root, "concurrent-existing-inventory.yaml")
		if err := os.WriteFile(destination, []byte("before\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		concurrent := []byte("concurrent existing output\n")
		err := writeInventoryArtifactWithHook(root, destination, inventory, func(phase string) error {
			if phase == "inventory-parent-bound" {
				return os.WriteFile(destination, concurrent, 0o600)
			}
			return nil
		})
		if err == nil || !strings.Contains(err.Error(), "bytes changed") {
			t.Fatalf("same-inode concurrent output error = %v", err)
		}
		captures, err := filepath.Glob(filepath.Join(root, "."+filepath.Base(destination)+".inventory-replaced-*"))
		if err != nil || len(captures) != 1 {
			t.Fatalf("concurrent output captures = %v, err=%v", captures, err)
		}
		if actual := mustRead(t, captures[0]); !bytes.Equal(actual, concurrent) {
			t.Fatalf("same-inode concurrent output capture changed: %q", actual)
		}
	})
	t.Run("quarantined bytes change", func(t *testing.T) {
		destination := filepath.Join(root, "concurrent-quarantined-inventory.yaml")
		if err := os.WriteFile(destination, []byte("before\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		concurrent := []byte("concurrent quarantined output\n")
		err := writeInventoryArtifactWithHook(root, destination, inventory, func(phase string) error {
			if phase != "inventory-output-quarantined" {
				return nil
			}
			captures, err := filepath.Glob(filepath.Join(root, "."+filepath.Base(destination)+".inventory-replaced-*"))
			if err != nil {
				return err
			}
			if len(captures) != 1 {
				return fmt.Errorf("expected one quarantined inventory output, found %d", len(captures))
			}
			return os.WriteFile(captures[0], concurrent, 0o600)
		})
		if err == nil || !strings.Contains(err.Error(), "bytes changed before cleanup") {
			t.Fatalf("post-quarantine same-inode output error = %v", err)
		}
		captures, err := filepath.Glob(filepath.Join(root, "."+filepath.Base(destination)+".inventory-replaced-*"))
		if err != nil || len(captures) != 1 {
			t.Fatalf("post-quarantine output captures = %v, err=%v", captures, err)
		}
		if actual := mustRead(t, captures[0]); !bytes.Equal(actual, concurrent) {
			t.Fatalf("post-quarantine concurrent output capture changed: %q", actual)
		}
	})
}

func TestPlansInventoryRequiresOutputAndProtectsAllPlanningAuthority(t *testing.T) {
	root := copyPlanCommandFixture(t)
	var stderr bytes.Buffer
	if code := RunPlans([]string{"inventory", root}, &bytes.Buffer{}, &stderr); code != 2 || !strings.Contains(stderr.String(), "--output requires a value") {
		t.Fatalf("missing output exit=%d stderr=%s", code, stderr.String())
	}

	t.Run("malformed formal plan", func(t *testing.T) {
		path := filepath.Join(root, "docs", "repo", "plans", "malformed", "malformed_discovery_doc.md")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		original := []byte("malformed planning authority\n")
		if err := os.WriteFile(path, original, 0o644); err != nil {
			t.Fatal(err)
		}
		stderr.Reset()
		if code := RunPlans([]string{"inventory", "--output", path, root}, &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "inventory_target_is_plan") {
			t.Fatalf("malformed target exit=%d stderr=%s", code, stderr.String())
		}
		if actual := mustRead(t, path); !bytes.Equal(actual, original) {
			t.Fatalf("malformed plan overwritten: %q", actual)
		}
	})

	t.Run("legacy register", func(t *testing.T) {
		path := filepath.Join(root, "docs", "repo", "plans", "plan-register.md")
		original := mustRead(t, path)
		stderr.Reset()
		if code := RunPlans([]string{"inventory", "--output", path, root}, &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "inventory_target_is_plan") {
			t.Fatalf("register target exit=%d stderr=%s", code, stderr.String())
		}
		if actual := mustRead(t, path); !bytes.Equal(actual, original) {
			t.Fatal("legacy register overwritten")
		}
	})
}

func TestPlansInventoryRejectsSymlinkAliasesToPlanningAuthority(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("ordinary Windows test users may not have symlink privilege")
	}
	root := copyPlanCommandFixture(t)
	cases := []struct {
		name       string
		linkTarget string
		linkName   string
		fileName   string
		realPath   string
	}{
		{
			name:       "formal plan",
			linkTarget: filepath.Join(root, "docs", "repo", "plans", "alpha"),
			linkName:   filepath.Join(root, "alpha-alias"),
			fileName:   "alpha_discovery_doc.md",
			realPath:   filepath.Join(root, "docs", "repo", "plans", "alpha", "alpha_discovery_doc.md"),
		},
		{
			name:       "legacy register",
			linkTarget: filepath.Join(root, "docs", "repo", "plans"),
			linkName:   filepath.Join(root, "plans-alias"),
			fileName:   "plan-register.md",
			realPath:   filepath.Join(root, "docs", "repo", "plans", "plan-register.md"),
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if err := os.Symlink(test.linkTarget, test.linkName); err != nil {
				t.Fatal(err)
			}
			before := mustRead(t, test.realPath)
			var stderr bytes.Buffer
			destination := filepath.Join(test.linkName, test.fileName)
			if code := RunPlans([]string{"inventory", "--output", destination, root}, &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "inventory_target_unsafe") {
				t.Fatalf("symlink alias exit=%d stderr=%s", code, stderr.String())
			}
			if after := mustRead(t, test.realPath); !bytes.Equal(after, before) {
				t.Fatalf("planning authority changed through alias: %s", test.realPath)
			}
		})
	}
}

func TestPlansInventoryRejectsCaseVariantExistingPlanIdentity(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("case-variant identity regression targets case-insensitive macOS filesystems")
	}
	root := copyPlanCommandFixture(t)
	realPath := filepath.Join(root, "docs", "repo", "plans", "alpha", "alpha_discovery_doc.md")
	variant := filepath.Join(root, "docs", "repo", "plans", "alpha", "ALPHA_DISCOVERY_DOC.MD")
	if _, err := os.Stat(variant); err != nil {
		t.Skip("current macOS volume is case-sensitive")
	}
	before := mustRead(t, realPath)
	var stderr bytes.Buffer
	if code := RunPlans([]string{"inventory", "--output", variant, root}, &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "inventory_target_is_plan") {
		t.Fatalf("case-variant target exit=%d stderr=%s", code, stderr.String())
	}
	if after := mustRead(t, realPath); !bytes.Equal(after, before) {
		t.Fatal("case-variant target overwrote the plan")
	}
}

func TestPlansInventoryRejectsProspectiveFormalPlanPaths(t *testing.T) {
	root := copyPlanCommandFixture(t)
	plansRoot := filepath.Join(root, "docs", "repo", "plans")
	familyRoot := filepath.Join(plansRoot, "prospective-family")
	if err := os.MkdirAll(filepath.Join(plansRoot, "new"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{
		filepath.Join("child-a", "child-a_discovery_doc.md"),
		filepath.Join("child-b", "child-b_implementation_doc.md"),
	} {
		path := filepath.Join(familyRoot, relative)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("Last updated: 2026-07-31T10:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: draft\n\n# Prospective Family Child\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	targets := []string{
		filepath.Join(plansRoot, "README.md"),
		filepath.Join(plansRoot, "readme.md"),
		filepath.Join(plansRoot, "new", "new_discovery_doc.md"),
		filepath.Join(plansRoot, "new", "new_implementation_doc.md"),
		filepath.Join(familyRoot, "README.md"),
	}
	for _, destination := range targets {
		t.Run(filepath.Base(destination), func(t *testing.T) {
			var stderr bytes.Buffer
			if code := RunPlans([]string{"inventory", "--output", destination, root}, &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "inventory_target_is_plan") {
				t.Fatalf("prospective target exit=%d stderr=%s", code, stderr.String())
			}
			if _, err := os.Lstat(destination); !os.IsNotExist(err) {
				t.Fatalf("prospective formal path was created: %s", destination)
			}
		})
	}
}

func TestPortfolioConfigureAppliesRoleSpecificIdentitySourcesAndScaffoldsIdempotently(t *testing.T) {
	root := t.TempDir()
	config := []byte("schema_version: 1\nselected_profile: standard\nproject_display_name: Portfolio Home\nselected_setup_folder: .\nlocal_consumer_layer:\n  repo_docs_path: docs/repo/\n  agent_memory_path: docs/agent-memory/\n  user_layer_path: .codeheart/user/\n  local_machine_layer_path: .codeheart/local/\ncomponent_settings: {}\n")
	configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, config, 0o644); err != nil {
		t.Fatal(err)
	}
	pristineOverlay, err := kitfs.ReadFile("components/planning-workflows/scaffolds/portfolio-strategic-overlay.yaml")
	if err != nil {
		t.Fatal(err)
	}
	writeCommandFixtureFile(t, root, "docs/repo/portfolio/strategic-overlay.yaml", pristineOverlay)
	localRoot := filepath.Join(root, "repositories")
	if err := os.Mkdir(localRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	args := []string{"configure", "--role", "coordination-home", "--member-repository-id", "portfolio-home", "--coordination-home-id", "example-home", "--github-owner", "zeta-org", "--github-owner", "example-org", "--github-owner", "example-org", "--local-root", localRoot, "--local-root", localRoot, "--yes", root}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunPortfolio(args, &stdout, &stderr); code != 0 {
		t.Fatalf("configure code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	configured := readYAMLFile(t, configPath)
	portfolioConfig := mapValue(configured["portfolio"])
	if valueString(portfolioConfig["role"]) != "coordination-home" || valueString(portfolioConfig["member_repository_id"]) != "portfolio-home" || valueString(portfolioConfig["coordination_home_id"]) != "example-home" {
		t.Fatalf("portfolio config = %#v", portfolioConfig)
	}
	sources := anyList(mapValue(portfolioConfig["discovery"])["sources"])
	if len(sources) != 2 || valueString(mapValue(sources[0])["owner"]) != "example-org" || valueString(mapValue(sources[1])["owner"]) != "zeta-org" {
		t.Fatalf("provider sources = %#v", sources)
	}
	localSources := readYAMLFile(t, filepath.Join(root, filepath.FromSlash(".codeheart/local/portfolio/sources.yaml")))
	if sources := anyList(localSources["sources"]); len(sources) != 1 || valueString(mapValue(sources[0])["root"]) != localRoot {
		t.Fatalf("local sources = %#v", sources)
	}
	for _, path := range []string{"docs/repo/portfolio/README.md", "docs/repo/portfolio/strategic-overlay.yaml"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(path))); err != nil {
			t.Fatalf("missing configured surface %s: %v", path, err)
		}
	}
	expectedReadme, err := kitfs.ReadFile("components/planning-workflows/scaffolds/portfolio-README.md")
	if err != nil {
		t.Fatal(err)
	}
	if actual := mustRead(t, filepath.Join(root, filepath.FromSlash("docs/repo/portfolio/README.md"))); !bytes.Equal(actual, expectedReadme) {
		t.Fatal("portfolio configure did not create the declared README scaffold")
	}
	overlay := readYAMLFile(t, filepath.Join(root, filepath.FromSlash("docs/repo/portfolio/strategic-overlay.yaml")))
	if valueString(overlay["coordination_home_id"]) != "example-home" {
		t.Fatalf("pristine overlay was not bound to approved home: %#v", overlay)
	}
	customOverlay := []byte("schema_version: 1\ncoordination_home_id: example-home\nfamilies: []\nthemes: []\nrelations: []\npriorities: []\nanalyses:\n  - id: retained-analysis\n    title: Retained analysis\n    updated_at: 2026-07-31T00:00:00Z\n    summary: Preserve repository-owned bytes.\n")
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash("docs/repo/portfolio/strategic-overlay.yaml")), customOverlay, 0o644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunPortfolio(args, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "Applied 0 change") {
		t.Fatalf("idempotent configure code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if actual := mustRead(t, filepath.Join(root, filepath.FromSlash("docs/repo/portfolio/strategic-overlay.yaml"))); !bytes.Equal(actual, customOverlay) {
		t.Fatal("matching reconfiguration overwrote repository-owned strategic overlay")
	}
	conflict := append([]string{}, args...)
	for index := range conflict {
		if conflict[index] == "portfolio-home" {
			conflict[index] = "different-home-repository"
			break
		}
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunPortfolio(conflict, &stdout, &stderr); code != 1 || !strings.Contains(stderr.String(), "portfolio_identity_conflict") {
		t.Fatalf("conflicting configure code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
}

func TestPortfolioConfigureMemberRejectsDiscoveryScopeAndDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	config := []byte("schema_version: 1\nselected_profile: standard\nproject_display_name: Member\nselected_setup_folder: .\nlocal_consumer_layer:\n  repo_docs_path: docs/repo/\n  agent_memory_path: docs/agent-memory/\n  user_layer_path: .codeheart/user/\ncomponent_settings: {}\n")
	configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, config, 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	base := []string{"configure", "--role", "member", "--member-repository-id", "member-one", "--coordination-home-id", "example-home"}
	if code := RunPortfolio(append(append([]string{}, base...), "--github-owner", "example-org", "--dry-run", root), &stdout, &stderr); code != 2 || !strings.Contains(stderr.String(), "only for role coordination-home") {
		t.Fatalf("member discovery code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := RunPortfolio(append(append([]string{}, base...), "--dry-run", root), &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "Planned 1 change") {
		t.Fatalf("member dry-run code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if actual := mustRead(t, configPath); !bytes.Equal(actual, config) {
		t.Fatal("portfolio dry-run changed shared config")
	}
}

func TestPortfolioConfigurePreservesUnsafeOverlayInsteadOfTreatingItAsPristine(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink authority regression runs on Unix")
	}
	root := t.TempDir()
	config := []byte("schema_version: 1\nselected_profile: standard\nproject_display_name: Portfolio Home\nselected_setup_folder: .\nlocal_consumer_layer:\n  repo_docs_path: docs/repo/\n  agent_memory_path: docs/agent-memory/\n  user_layer_path: .codeheart/user/\ncomponent_settings: {}\n")
	writeCommandFixtureFile(t, root, state.ConfigPath, config)
	pristineOverlay, err := kitfs.ReadFile("components/planning-workflows/scaffolds/portfolio-strategic-overlay.yaml")
	if err != nil {
		t.Fatal(err)
	}
	external := filepath.Join(t.TempDir(), "external-overlay.yaml")
	if err := os.WriteFile(external, pristineOverlay, 0o644); err != nil {
		t.Fatal(err)
	}
	overlayPath := filepath.Join(root, filepath.FromSlash(portfolio.OverlayPath))
	if err := os.MkdirAll(filepath.Dir(overlayPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, overlayPath); err != nil {
		t.Fatal(err)
	}
	args := []string{"configure", "--role", "coordination-home", "--member-repository-id", "portfolio-home", "--coordination-home-id", "example-home", "--yes", root}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunPortfolio(args, &stdout, &stderr); code != 0 {
		t.Fatalf("configure code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if actual := mustRead(t, external); !bytes.Equal(actual, pristineOverlay) {
		t.Fatal("portfolio configure changed bytes behind unsafe overlay symlink")
	}
	info, err := os.Lstat(overlayPath)
	if err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("portfolio configure replaced unsafe overlay: info=%v err=%v", info, err)
	}
}

func TestPlansRemoteOverlayValidationAndInventoryUsePushedRefsOnly(t *testing.T) {
	root := t.TempDir()
	remote := filepath.Join(t.TempDir(), "remote.git")
	if output, err := exec.Command("git", "init", "--bare", remote).CombinedOutput(); err != nil {
		t.Fatalf("init bare remote: %v\n%s", err, output)
	}
	runGitCommandTest(t, root, "init", "-b", "main")
	runGitCommandTest(t, root, "config", "user.email", "test@example.invalid")
	runGitCommandTest(t, root, "config", "user.name", "Remote Overlay Test")
	config := []byte("schema_version: 1\nselected_profile: standard\nproject_display_name: Remote Overlay\nselected_setup_folder: .\nlocal_consumer_layer:\n  repo_docs_path: docs/repo/\n  agent_memory_path: docs/agent-memory/\n  user_layer_path: .codeheart/user/\n  local_machine_layer_path: .codeheart/local/\ncomponent_settings:\n  planning-workflows:\n    plan_catalog_mode: canonical\nportfolio:\n  schema_version: 2\n  role: member\n  member_repository_id: remote-example\n  coordination_home_id: example-home\n")
	lock := []byte("schema_version: 1\nkit_version: 0.1.23\nselected_profile: standard\nselected_components: []\nrelease:\n  asset_url: local-test\n  checksum_sha256: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\nmanaged_paths: []\ngenerated_surfaces: []\ncli_repair:\n  installed_cli_path: codeheart-operating-kit\n  repair_source_url: local-test\nupdate_check:\n  last_update_check_at: \"2026-07-01T00:00:00Z\"\n  next_update_check_due: \"2026-08-01T00:00:00Z\"\n  latest_seen_version: 0.1.23\n  update_status: current\nnative_capabilities: {}\n")
	writeCommandFixtureFile(t, root, state.ConfigPath, config)
	writeCommandFixtureFile(t, root, state.LockPath, lock)
	writeCommandFixtureFile(t, root, ".codeheart/kit/README.md", []byte("# Installed Kit\n"))
	planPath := "docs/repo/plans/remote-example/remote-example_discovery_doc.md"
	writeCommandFixtureFile(t, root, planPath, remoteOverlayPlan("Baseline purpose"))
	runGitCommandTest(t, root, "add", ".")
	runGitCommandTest(t, root, "commit", "-m", "Remote baseline")
	runGitCommandTest(t, root, "remote", "add", "origin", remote)
	runGitCommandTest(t, root, "push", "-u", "origin", "main")
	runGitCommandTest(t, remote, "symbolic-ref", "HEAD", "refs/heads/main")
	runGitCommandTest(t, root, "checkout", "-b", "feature/remote-plan-change")
	writeCommandFixtureFile(t, root, planPath, remoteOverlayPlan("Pushed branch purpose"))
	runGitCommandTest(t, root, "add", planPath)
	runGitCommandTest(t, root, "commit", "-m", "Change pushed plan")
	runGitCommandTest(t, root, "push", "origin", "feature/remote-plan-change")
	runGitCommandTest(t, root, "checkout", "main")
	runGitCommandTest(t, root, "checkout", "-b", "local-only-plan")
	writeCommandFixtureFile(t, root, "docs/repo/plans/local-only/local-only_discovery_doc.md", bytes.ReplaceAll(remoteOverlayPlan("Unpushed purpose"), []byte("remote-example.discovery.remote-example"), []byte("remote-example.discovery.local-only")))
	runGitCommandTest(t, root, "add", "docs/repo/plans/local-only")
	runGitCommandTest(t, root, "commit", "-m", "Local-only plan")
	runGitCommandTest(t, root, "checkout", "main")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunPlans([]string{"validate", "--remote-overlays", root}, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), "2 pushed observation(s)") || !strings.Contains(stdout.String(), "Local heads, worktree changes, and unpushed commits are omitted") {
		t.Fatalf("remote validate code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	destination := filepath.Join(root, "remote-overlay-inventory.json")
	stdout.Reset()
	stderr.Reset()
	if code := RunPlans([]string{"inventory", "--remote-overlays", "--output", destination, root}, &stdout, &stderr); code != 0 {
		t.Fatalf("remote inventory code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	data := mustRead(t, destination)
	if !bytes.Contains(data, []byte(`"remote_overlays"`)) || !bytes.Contains(data, []byte(`"visibility": "unmerged-branch"`)) || bytes.Contains(data, []byte("local-only")) {
		t.Fatalf("remote inventory evidence is incomplete or includes local-only state:\n%s", data)
	}
}

func TestPlanRemoteOverlayEvidenceRedactsMachineLocalSourceLocators(t *testing.T) {
	root := filepath.Join(t.TempDir(), "repository")
	result := portfolio.ScanResult{Catalog: portfolio.Catalog{
		Members: []plancatalog.CatalogMember{
			{RepositoryID: "example-member", SourceLocator: root},
			{RepositoryID: "public-member", SourceLocator: "github.com/example/public-member"},
		},
		Candidates: []portfolio.Candidate{{RepositoryID: "candidate-member", SourceLocator: "file:" + root}},
		Errors:     []portfolio.ScanError{{Code: "example", Message: "cannot read " + root, SourceLocator: root}},
	}}
	public := publicRemoteOverlayEvidence(result)
	encoded, err := json.Marshal(public)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(encoded, []byte(root)) || bytes.Contains(encoded, []byte("file:"+root)) {
		t.Fatalf("public remote-overlay evidence retained a machine path: %s", encoded)
	}
	if public.Catalog.Members[0].SourceLocator != "local-git:example-member" ||
		public.Catalog.Candidates[0].SourceLocator != "local-git:candidate-member" ||
		public.Catalog.Errors[0].SourceLocator != "local-git" {
		t.Fatalf("redacted locators = %+v", public.Catalog)
	}
	if public.Catalog.Members[1].SourceLocator != "github.com/example/public-member" {
		t.Fatalf("public provider locator changed: %+v", public.Catalog.Members[1])
	}
	if public.Catalog.Errors[0].Message != "remote portfolio scan reported an error" {
		t.Fatalf("public scan error message = %q", public.Catalog.Errors[0].Message)
	}
}

func TestPlansRemoteOverlayInventoryRedactsFailingLocalSourceMessage(t *testing.T) {
	root := t.TempDir()
	remote := filepath.Join(t.TempDir(), "remote.git")
	if output, err := exec.Command("git", "init", "--bare", remote).CombinedOutput(); err != nil {
		t.Fatalf("init bare remote: %v\n%s", err, output)
	}
	runGitCommandTest(t, root, "init", "-b", "main")
	runGitCommandTest(t, root, "config", "user.email", "test@example.invalid")
	runGitCommandTest(t, root, "config", "user.name", "Remote Error Redaction Test")
	config := []byte("schema_version: 1\nselected_profile: standard\nproject_display_name: Error Redaction Home\nselected_setup_folder: .\nlocal_consumer_layer:\n  repo_docs_path: docs/repo/\n  agent_memory_path: docs/agent-memory/\n  user_layer_path: .codeheart/user/\n  local_machine_layer_path: .codeheart/local/\ncomponent_settings:\n  planning-workflows:\n    plan_catalog_mode: canonical\nportfolio:\n  schema_version: 2\n  role: coordination-home\n  member_repository_id: error-redaction-home\n  coordination_home_id: example-home\n")
	lock := []byte("schema_version: 1\nkit_version: 0.1.23\nselected_profile: standard\nselected_components: []\nrelease:\n  asset_url: local-test\n  checksum_sha256: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\nmanaged_paths: []\ngenerated_surfaces: []\ncli_repair:\n  installed_cli_path: codeheart-operating-kit\n  repair_source_url: local-test\nupdate_check:\n  last_update_check_at: \"2026-07-01T00:00:00Z\"\n  next_update_check_due: \"2026-08-01T00:00:00Z\"\n  latest_seen_version: 0.1.23\n  update_status: current\nnative_capabilities: {}\n")
	writeCommandFixtureFile(t, root, state.ConfigPath, config)
	writeCommandFixtureFile(t, root, state.LockPath, lock)
	writeCommandFixtureFile(t, root, ".codeheart/kit/README.md", []byte("# Installed Kit\n"))
	plan := bytes.ReplaceAll(remoteOverlayPlan("Error redaction baseline"), []byte("remote-example.discovery.remote-example"), []byte("error-redaction-home.discovery.baseline"))
	writeCommandFixtureFile(t, root, "docs/repo/plans/error-redaction/error-redaction_discovery_doc.md", plan)
	runGitCommandTest(t, root, "add", ".")
	runGitCommandTest(t, root, "commit", "-m", "Remote error redaction baseline")
	runGitCommandTest(t, root, "remote", "add", "origin", remote)
	runGitCommandTest(t, root, "push", "-u", "origin", "main")
	runGitCommandTest(t, remote, "symbolic-ref", "HEAD", "refs/heads/main")

	missingRoot := filepath.Join(t.TempDir(), "missing-local-source")
	localSources := []byte(fmt.Sprintf("schema_version: 1\nsources:\n  - kind: local-git-root\n    root: %s\n", missingRoot))
	writeCommandFixtureFile(t, root, portfolio.LocalSourcesPath, localSources)
	destination := filepath.Join(root, "remote-overlay-error.json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := RunPlans([]string{"inventory", "--remote-overlays", "--output", destination, root}, &stdout, &stderr); code != 1 {
		t.Fatalf("remote inventory code=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	data := mustRead(t, destination)
	if bytes.Contains(data, []byte(missingRoot)) || !bytes.Contains(data, []byte(`"code": "source_discovery_failed"`)) || !bytes.Contains(data, []byte(`"message": "configured portfolio source discovery failed"`)) {
		t.Fatalf("public failing-source artifact was not safely redacted:\n%s", data)
	}
}

func writeCommandFixtureFile(t *testing.T, root, relative string, data []byte) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func remoteOverlayPlan(purpose string) []byte {
	return []byte("Last updated: 2026-07-01T00:00:00Z (UTC)\nCreated: 2026-07-01\nStatus: active\n\n# Remote Example Discovery\n<!-- BEGIN CODEHEART PLAN METADATA -->\n```yaml\nplan:\n  schema_version: 1\n  id: remote-example.discovery.remote-example\n  kind: discovery\n  purpose: " + purpose + "\n  first_cataloged: 2026-07-01T00:00:00Z\n  catalog_metadata_updated: 2026-07-01T00:00:00Z\n```\n<!-- END CODEHEART PLAN METADATA -->\n\n## Scope\n\nRemote fixture.\n")
}

func TestPlansMigrateRequiresExplicitModeAndReviewedLedger(t *testing.T) {
	var stderr bytes.Buffer
	if code := RunPlans([]string{"migrate", "--ledger", "reviewed.yaml", "."}, &bytes.Buffer{}, &stderr); code != 2 || !strings.Contains(stderr.String(), "choose exactly one") {
		t.Fatalf("migrate approval error code=%d stderr=%s", code, stderr.String())
	}
	stderr.Reset()
	if code := RunPlans([]string{"migrate", "--dry-run", "."}, &bytes.Buffer{}, &stderr); code != 2 || !strings.Contains(stderr.String(), "--ledger requires a value") {
		t.Fatalf("migrate ledger error code=%d stderr=%s", code, stderr.String())
	}
}

func copyPlanCommandFixture(t *testing.T) string {
	t.Helper()
	source := filepath.Join("..", "..", "tests", "fixtures", "plans", "migration-repository")
	root := t.TempDir()
	if err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil || relative == "." {
			return err
		}
		target := filepath.Join(root, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return os.WriteFile(target, mustRead(t, path), 0o644)
	}); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"init"}, {"checkout", "-b", "main"}, {"config", "user.email", "test@example.invalid"}, {"config", "user.name", "Plan Command Test"}, {"add", "."}, {"commit", "-m", "fixture"}} {
		command := exec.Command("git", append([]string{"-C", root}, args...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
		}
	}
	baseline := strings.TrimSpace(runGitCommandTest(t, root, "rev-parse", "HEAD"))
	configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	configData := mustRead(t, configPath)
	configData = bytes.Replace(configData, []byte("plan_catalog_mode: legacy"), []byte("plan_catalog_mode: mixed\n    plan_catalog_cutover_revision: "+baseline), 1)
	if err := os.WriteFile(configPath, configData, 0o644); err != nil {
		t.Fatal(err)
	}
	runGitCommandTest(t, root, "add", state.ConfigPath)
	runGitCommandTest(t, root, "commit", "-m", "adopt mixed plan catalog")
	return root
}

func runGitCommandTest(t *testing.T, root string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func TestInitWritesStandardSurfaces(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	code := RunInit([]string{root, "--project-name", "Example-Automation", "--purpose", "company-automation"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunInit exit = %d; stdout: %s", code, stdout.String())
	}
	if !strings.Contains(stdout.String(), "Operating Kit initialized.") {
		t.Fatalf("init text output missing success line: %q", stdout.String())
	}
	for _, relative := range []string{
		".codeheart/kit",
		".codeheart/kit/README.md",
		".codeheart/kit.lock.yaml",
		".codeheart/kit.config.yaml",
		".codeheart/user/README.md",
		".codeheart/user/examples/preferences.yaml",
		"AGENTS.md",
		"docs/repo/README.md",
		"docs/repo/plans/plan-register.md",
		"docs/repo/portfolio/README.md",
		"docs/repo/portfolio/strategic-overlay.yaml",
		"docs/agent-memory/README.md",
		"docs/agent-memory/goal-register.md",
		"docs/agent-memory/session-ledger.md",
		"docs/agent-memory/untriaged-sessions.md",
		".gitignore",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(relative))); err != nil {
			t.Fatalf("expected %s to exist: %v", relative, err)
		}
	}
	config := readYAMLFile(t, filepath.Join(root, ".codeheart", "kit.config.yaml"))
	if config["setup_purpose"] != "company-automation" {
		t.Fatalf("setup_purpose = %#v", config["setup_purpose"])
	}
	planning := state.Map(state.Map(config["component_settings"])["planning-workflows"])
	if state.AsInt(planning["plan_catalog_discovery_version"]) != 2 {
		t.Fatalf("fresh discovery version = %#v", planning["plan_catalog_discovery_version"])
	}
	ownership := state.Map(planning["plan_catalog_ownership"])
	if roots, ok := ownership["excluded_roots"].([]any); !ok || len(roots) != 0 {
		t.Fatalf("fresh excluded roots = %#v", ownership["excluded_roots"])
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	if set := mapKeys(mapValue(lock["native_capabilities"])); !sameStringSet(set, []string{"browser", "documents", "pdf", "presentations", "spreadsheets"}) {
		t.Fatalf("native capabilities = %#v", set)
	}
}

func TestInitJSONOutput(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	code := RunInit([]string{root, "--project-name", "Example-Automation", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunInit --json exit = %d; stdout: %s", code, stdout.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse init JSON: %v\n%s", err, stdout.String())
	}
	if payload["inspection"] == nil || payload["state"] == nil {
		t.Fatalf("init JSON missing expected keys: %#v", payload)
	}
}

func TestValueFlagsSupportEqualsFormAndRejectMissingValues(t *testing.T) {
	root := t.TempDir()
	equalStyleTarget := filepath.Join(root, "equals")
	code := RunInit([]string{equalStyleTarget, "--project-name=EqualStyle"}, &bytes.Buffer{}, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunInit with equals-style flag exit = %d", code)
	}
	config := readYAMLFile(t, filepath.Join(equalStyleTarget, ".codeheart", "kit.config.yaml"))
	if config["project_display_name"] != "EqualStyle" {
		t.Fatalf("project_display_name = %#v", config["project_display_name"])
	}

	malformedTarget := filepath.Join(root, "malformed")
	var stderr bytes.Buffer
	code = RunInit([]string{malformedTarget, "--project-name", "--json"}, &bytes.Buffer{}, &stderr)
	if code != 2 {
		t.Fatalf("RunInit malformed flag exit = %d, want 2; stderr: %s", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(malformedTarget, ".codeheart")); !os.IsNotExist(err) {
		t.Fatalf("malformed init should not create .codeheart, stat err: %v", err)
	}

	updateTarget := filepath.Join(root, "update")
	if code := RunInit([]string{updateTarget, "--project-name=Update"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunInit update target exit = %d", code)
	}
	stderr.Reset()
	code = RunUpdateCheck([]string{updateTarget, "--latest-version", "--json"}, &bytes.Buffer{}, &stderr)
	if code != 2 {
		t.Fatalf("RunUpdateCheck malformed flag exit = %d, want 2; stderr: %s", code, stderr.String())
	}
	lock, err := lockfile.ReadLock(updateTarget)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	if mapValue(lock["update_check"])["latest_seen_version"] == "--json" {
		t.Fatalf("malformed update-check wrote --json as latest version")
	}
}

func TestLocalFileURLPathForWindowsDriveAndHost(t *testing.T) {
	driveURL, err := url.Parse("file:///C:/Users/Example%20User/latest-release.json")
	if err != nil {
		t.Fatalf("parse drive URL: %v", err)
	}
	if got, want := localFileURLPathForGOOS(driveURL, "windows"), `C:\Users\Example User\latest-release.json`; got != want {
		t.Fatalf("windows drive file URL path = %q, want %q", got, want)
	}

	uncURL, err := url.Parse("file://server/share/latest-release.json")
	if err != nil {
		t.Fatalf("parse UNC URL: %v", err)
	}
	if got, want := localFileURLPathForGOOS(uncURL, "windows"), `\\server\share\latest-release.json`; got != want {
		t.Fatalf("windows UNC file URL path = %q, want %q", got, want)
	}
}

func TestOnboardYesWritesBaseSetupWithoutNativeCapabilityOfferOrCheck(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	code := RunOnboard([]string{"--target", root, "--project-name", "Companyname-Automation", "--yes", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunOnboard exit = %d; stdout: %s", code, stdout.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse onboard JSON: %v\n%s", err, stdout.String())
	}
	if payload["written"] != true {
		t.Fatalf("onboard did not report written: %#v", payload)
	}
	if _, exists := payload["native_capabilities"]; exists {
		t.Fatalf("onboard --yes should not report native capability checks: %#v", payload["native_capabilities"])
	}
	joinedScript := strings.Join(anyStrings(payload["script"]), "\n")
	for _, forbidden := range []string{
		"Should I check and set up these tools now?",
		"After setup writes complete, ask whether to check native Codex capabilities.",
		"documents, spreadsheets, presentations, browser work, and PDFs",
	} {
		if strings.Contains(joinedScript, forbidden) {
			t.Fatalf("onboarding script contains forbidden optional capability prompt %q:\n%s", forbidden, joinedScript)
		}
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	for capability, value := range mapValue(lock["native_capabilities"]) {
		record := mapValue(value)
		if record["status"] != "unknown" || record["command_result_category"] != "not-checked" {
			t.Fatalf("%s capability was checked or installed: %#v", capability, record)
		}
	}
}

func TestOnboardYesRequiresTargetAndProjectName(t *testing.T) {
	var stdout bytes.Buffer
	code := RunOnboard([]string{"--project-name", "Companyname-Automation", "--yes", "--json"}, &stdout, &bytes.Buffer{})
	if code != 2 {
		t.Fatalf("missing target exit = %d", code)
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse missing target JSON: %v", err)
	}
	missing := anyStrings(payload["required_user_decisions_missing"])
	if len(missing) != 1 || missing[0] != "target_folder" {
		t.Fatalf("missing decisions = %#v", missing)
	}
}

func TestSyncRepairsManagedDriftAndCheckPassesWithCLIMarker(t *testing.T) {
	t.Setenv("CODEHEART_OPERATING_KIT_CLI", "1")
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "Example-Automation"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunInit exit = %d", code)
	}
	managed := filepath.Join(root, ".codeheart", "kit", "docs", "agent-interface", "README.md")
	if err := os.WriteFile(managed, []byte("drift\n"), 0o644); err != nil {
		t.Fatalf("write drift: %v", err)
	}
	before, err := CheckRepository(root)
	if err != nil {
		t.Fatalf("CheckRepository before: %v", err)
	}
	if len(before["drift"].([]any)) == 0 {
		t.Fatalf("expected drift before sync")
	}
	if code := RunSync([]string{root}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunSync exit = %d", code)
	}
	var syncJSON bytes.Buffer
	if code := RunSync([]string{root, "--json"}, &syncJSON, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunSync --json exit = %d; stdout: %s", code, syncJSON.String())
	}
	var syncPayload map[string]any
	if err := json.Unmarshal(syncJSON.Bytes(), &syncPayload); err != nil {
		t.Fatalf("parse sync JSON: %v\n%s", err, syncJSON.String())
	}
	if syncPayload["kit_version"] == "" || syncPayload["synced_managed_paths"] == nil {
		t.Fatalf("sync JSON missing expected keys: %#v", syncPayload)
	}
	after, err := CheckRepository(root)
	if err != nil {
		t.Fatalf("CheckRepository after: %v", err)
	}
	if after["ok"] != true || len(after["drift"].([]any)) != 0 {
		t.Fatalf("check after sync = %#v", after)
	}
	var checkText bytes.Buffer
	if code := RunCheck([]string{root}, &checkText, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunCheck text exit = %d; stdout: %s", code, checkText.String())
	}
	if !strings.Contains(checkText.String(), "Operating Kit check") || !strings.Contains(checkText.String(), "OK: true") {
		t.Fatalf("check text output unexpected: %q", checkText.String())
	}
}

func TestCLIAvailableRecognizesInstalledExeName(t *testing.T) {
	t.Setenv("CODEHEART_OPERATING_KIT_CLI", "")
	root := t.TempDir()
	t.Setenv("PATH", root)
	exe := filepath.Join(root, "codeheart-operating-kit.exe")
	if err := os.WriteFile(exe, []byte(""), 0o755); err != nil {
		t.Fatalf("write exe marker: %v", err)
	}
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()
	os.Args = []string{exe}
	if !CLIAvailable() {
		t.Fatalf("CLIAvailable should recognize codeheart-operating-kit.exe as installed")
	}
}

func TestUpdateCheckWritesCadenceAndFailurePreservesDueDate(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "Example-Automation"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunInit exit = %d", code)
	}
	var stdout bytes.Buffer
	code := RunUpdateCheck([]string{root, "--latest-version", "0.2.0", "--now", "2026-06-13T00:00:00Z", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunUpdateCheck exit = %d; stdout: %s", code, stdout.String())
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	update := mapValue(lock["update_check"])
	if update["update_status"] != "update-available" || update["next_update_check_due"] != "2026-06-20T00:00:00Z" {
		t.Fatalf("unexpected update state: %#v", update)
	}
	beforeDue := update["next_update_check_due"]
	stdout.Reset()
	missingURL := (&urlBuilder{path: filepath.Join(root, "missing-release.json")}).String()
	code = RunUpdateCheck([]string{root, "--metadata-url", missingURL, "--now", "2026-06-14T00:00:00Z", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunUpdateCheck failure case exit = %d; stdout: %s", code, stdout.String())
	}
	lock, err = lockfile.ReadLock(root)
	if err != nil {
		t.Fatalf("ReadLock after failure: %v", err)
	}
	update = mapValue(lock["update_check"])
	if update["update_status"] != "failed" || update["next_update_check_due"] != beforeDue {
		t.Fatalf("failed update state = %#v, previous due %v", update, beforeDue)
	}
	var text bytes.Buffer
	code = RunUpdateCheck([]string{root, "--latest-version", "0.1.24"}, &text, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunUpdateCheck text exit = %d; stdout: %s", code, text.String())
	}
	if !strings.Contains(text.String(), "Operating Kit is current.") {
		t.Fatalf("update-check text output unexpected: %q", text.String())
	}
}

func TestUpdateCheckEqualsStyleLatestVersion(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "Example-Automation"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunInit exit = %d", code)
	}
	var stdout bytes.Buffer
	code := RunUpdateCheck([]string{root, "--latest-version=0.2.0", "--now=2026-06-13T00:00:00Z", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunUpdateCheck equals-style exit = %d; stdout: %s", code, stdout.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse update JSON: %v\n%s", err, stdout.String())
	}
	if payload["status"] != "update-available" {
		t.Fatalf("update status = %#v", payload)
	}
}

type urlBuilder struct {
	path string
}

func (builder *urlBuilder) String() string {
	return "file://" + filepath.ToSlash(builder.path)
}

func readYAMLFile(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read YAML %s: %v", path, err)
	}
	parsed, err := yamlmini.MustMap(string(data))
	if err != nil {
		t.Fatalf("parse YAML %s: %v", path, err)
	}
	return parsed
}

func anyStrings(value any) []string {
	result := []string{}
	for _, item := range anyList(value) {
		result = append(result, valueString(item))
	}
	return result
}

func mapKeys(value map[string]any) []string {
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	return keys
}

func sameStringSet(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	seen := map[string]int{}
	for _, value := range left {
		seen[value]++
	}
	for _, value := range right {
		seen[value]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}
