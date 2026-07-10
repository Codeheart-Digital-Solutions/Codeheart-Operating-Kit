package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	releasekit "github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/release"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

func RunUpgrade(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{
		"--version":          true,
		"--catalog":          true,
		"--installed-binary": true,
		"--dry-run":          false,
		"--yes":              false,
		"--json":             false,
	})
	if err != nil {
		return writeArgError(stderr, "upgrade", err)
	}
	if values["--version"] == "" {
		return writeArgError(stderr, "upgrade", fmt.Errorf("option --version requires a value"))
	}
	if bools["--dry-run"] == bools["--yes"] {
		return writeArgError(stderr, "upgrade", fmt.Errorf("choose exactly one of --dry-run or --yes"))
	}
	root := "."
	if len(positionals) > 0 {
		root = positionals[0]
	}
	if len(positionals) > 1 {
		return writeArgError(stderr, "upgrade", fmt.Errorf("unexpected argument %q", positionals[1]))
	}
	payload, operation, err := upgradeOperation(root, values["--version"], values["--catalog"], values["--installed-binary"], bools["--dry-run"])
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit upgrade: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, payload); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit upgrade: error: %v\n", err)
			return 1
		}
	} else {
		reconcile.WriteText(stdout, operation)
	}
	if operation.OK() {
		return 0
	}
	return 1
}

func upgradeOperation(root, targetVersion, catalogLocation, installedBinary string, dryRun bool) (map[string]any, reconcile.Result, error) {
	root = expandPath(root)
	observed, err := state.Inspect(root)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	if state.AsInt(observed.Lock["schema_version"]) != 2 || (observed.Classification != state.StateCurrent && observed.Classification != state.StateDrifted && observed.Classification != state.StateStaleCLI) {
		result := blockedLifecycle("upgrade", observed, "upgrade_requires_valid_v2_installation", "upgrade requires a valid lock-v2 installation", "run check and repair compatible state first", "repair --dry-run")
		return resultPayload(result), result, nil
	}
	currentVersion := state.AsString(observed.Lock["kit_version"])
	if err := releasekit.RequireForwardUpgrade(currentVersion, targetVersion); err != nil {
		result := blockedLifecycle("upgrade", observed, "upgrade_direction_invalid", err.Error(), "select a release newer than the installed version", "upgrade --dry-run")
		return resultPayload(result), result, nil
	}
	if catalogLocation == "" {
		catalogLocation = fmt.Sprintf("https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v%s/release-catalog-%s.json", targetVersion, targetVersion)
	}
	workDir, err := os.MkdirTemp("", "codeheart-operating-kit-upgrade-")
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	defer os.RemoveAll(workDir)
	prepared, err := releasekit.PrepareUpgrade(catalogLocation, targetVersion, workDir, true)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	if installedBinary == "" {
		installedBinary, err = releasekit.DefaultInstalledBinary()
		if err != nil {
			return nil, reconcile.Result{}, err
		}
	}
	result := reconcile.NewResult("upgrade")
	result.StateBefore = string(observed.Classification)
	result.TransactionID = prepared.Pack.PackManifestSHA256[:32]
	result.Changes = append(result.Changes,
		reconcile.Change{Action: "replace", Path: installedBinary, Owner: "installed-binary"},
		reconcile.Change{Action: "reconcile", Path: state.LockPath, Owner: "generated-surface"},
	)
	result.Provenance = upgradeProvenance(prepared)
	if dryRun {
		result.Status = reconcile.StatusPlanned
		result.DryRun = true
		result.StateAfter = string(observed.Classification)
		result.Validations = append(result.Validations,
			reconcile.Validation{Name: "catalog-to-archive", Status: "passed"},
			reconcile.Validation{Name: "pack-to-binary", Status: "passed"},
			reconcile.Validation{Name: "staged-version", Status: "passed"},
		)
		payload := resultPayload(result)
		payload["target_version"] = targetVersion
		return payload, result, nil
	}
	handoff, err := releasekit.NewHandoff(prepared, root, installedBinary, currentVersion)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	if runtime.GOOS == "windows" {
		if err := releasekit.StartDeferredHandoff(handoff); err != nil {
			return nil, reconcile.Result{}, err
		}
		result.Status = reconcile.StatusSucceeded
		result.StateAfter = string(observed.Classification)
		result.Validations = append(result.Validations, reconcile.Validation{Name: "deferred-handoff", Status: "passed", Detail: "new binary will replace the installed binary after this process exits"})
	} else {
		if err := releasekit.ApplyHandoff(handoff); err != nil {
			result.Status = reconcile.StatusRolledBack
			result.Rollback = reconcile.Rollback{Attempted: true, Succeeded: true, Detail: err.Error()}
			result.Blockers = append(result.Blockers, reconcile.Blocker{Code: "upgrade_rolled_back", Message: err.Error(), RetryCommand: "upgrade --dry-run"})
		} else {
			result.Status = reconcile.StatusSucceeded
			result.StateAfter = string(state.StateCurrent)
			result.Validations = append(result.Validations, reconcile.Validation{Name: "binary-and-state-handoff", Status: "passed"})
		}
	}
	payload := resultPayload(result)
	payload["target_version"] = targetVersion
	return payload, result, nil
}

func upgradeProvenance(prepared releasekit.PreparedUpgrade) map[string]any {
	return map[string]any{
		"catalog_location":        prepared.Catalog.Location,
		"catalog_sha256":          prepared.Catalog.DigestSHA256,
		"archive_sha256":          prepared.Pack.ArchiveSHA256,
		"pack_manifest_sha256":    prepared.Pack.PackManifestSHA256,
		"content_manifest_sha256": prepared.Pack.ContentManifestSHA256,
		"binary_sha256":           prepared.Pack.Manifest.BinarySHA256,
	}
}

func RunUpgradeReconcile(args []string, stdout io.Writer, stderr io.Writer) int {
	values, _, positionals, err := parseValueArgs(args, map[string]bool{
		"--repository":              true,
		"--previous-version":        true,
		"--asset-url":               true,
		"--catalog-location":        true,
		"--catalog-sha256":          true,
		"--archive-sha256":          true,
		"--pack-manifest-sha256":    true,
		"--content-manifest-sha256": true,
		"--binary-sha256":           true,
	})
	if err != nil || len(positionals) > 0 {
		fmt.Fprintf(stderr, "upgrade reconcile arguments are invalid\n")
		return 2
	}
	root := values["--repository"]
	observed, err := state.Inspect(root)
	if err != nil || state.AsInt(observed.Lock["schema_version"]) != 2 || state.AsString(observed.Lock["kit_version"]) != values["--previous-version"] {
		fmt.Fprintf(stderr, "upgrade reconcile precondition failed\n")
		return 1
	}
	graph, err := compileObservedGraph(observed)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	now := time.Now().UTC().Truncate(time.Second)
	lock, err := desiredLifecycleLock("upgrade", observed, graph, now)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	lock["kit_version"] = version.Version
	lock["release"] = map[string]any{"asset_url": values["--asset-url"], "checksum_sha256": values["--archive-sha256"]}
	lock["release_provenance"] = map[string]any{
		"verification_status":  "verified",
		"source":               values["--catalog-location"],
		"catalog_url":          values["--catalog-location"],
		"catalog_sha256":       values["--catalog-sha256"],
		"archive_sha256":       values["--archive-sha256"],
		"pack_manifest_sha256": values["--pack-manifest-sha256"],
		"binary_sha256":        values["--binary-sha256"],
		"verified_at":          now.Format(time.RFC3339),
	}
	result, err := runLifecycle(lifecycleRequest{command: "upgrade", root: root, now: now, observed: observed, graph: graph, desiredLock: lock, ensureIgnore: true})
	if err != nil || !result.OK() {
		fmt.Fprintf(stderr, "upgrade reconcile failed: %v\n", err)
		return 1
	}
	check, err := CheckRepository(root)
	if err != nil || check["ok"] != true {
		fmt.Fprintf(stderr, "upgrade post-check failed\n")
		return 1
	}
	return 0
}

func RunUpgradeHandoff(args []string, stderr io.Writer) int {
	values, _, positionals, err := parseValueArgs(args, map[string]bool{"--file": true})
	if err != nil || len(positionals) > 0 || values["--file"] == "" {
		fmt.Fprintln(stderr, "upgrade handoff arguments are invalid")
		return 2
	}
	if err := releasekit.ExecuteDeferredHandoff(filepath.Clean(values["--file"])); err != nil {
		fmt.Fprintf(stderr, "upgrade handoff failed: %v\n", err)
		return 1
	}
	return 0
}

func RunVerifyContentIdentity(args []string, stderr io.Writer) int {
	values, _, positionals, err := parseValueArgs(args, map[string]bool{"--path": true, "--version": true})
	if err != nil || len(positionals) > 0 || values["--path"] == "" || values["--version"] == "" {
		fmt.Fprintln(stderr, "content identity arguments are invalid")
		return 2
	}
	data, err := os.ReadFile(values["--path"])
	if err != nil {
		fmt.Fprintln(stderr, "content identity is unreadable")
		return 1
	}
	content, err := state.DecodeAndValidateYAML(state.ContentManifestSchema, data)
	if err != nil || state.AsString(content["version"]) != values["--version"] {
		fmt.Fprintln(stderr, "content identity does not match the staged version")
		return 1
	}
	return 0
}

func RunVerifyReleaseEvidence(args []string, stderr io.Writer) int {
	values, _, positionals, err := parseValueArgs(args, map[string]bool{"--catalog": true, "--version": true})
	if err != nil || len(positionals) > 0 || values["--catalog"] == "" || values["--version"] == "" {
		fmt.Fprintln(stderr, "release evidence arguments are invalid")
		return 2
	}
	workDir, err := os.MkdirTemp("", "codeheart-release-evidence-")
	if err != nil {
		fmt.Fprintln(stderr, "release evidence temporary storage failed")
		return 1
	}
	defer os.RemoveAll(workDir)
	if _, err := releasekit.PrepareUpgrade(values["--catalog"], values["--version"], workDir, false); err != nil {
		fmt.Fprintf(stderr, "release evidence verification failed: %v\n", err)
		return 1
	}
	return 0
}

func RunCleanupUpgradeHandoff(args []string, stderr io.Writer) int {
	values, _, positionals, err := parseValueArgs(args, map[string]bool{"--path": true, "--parent-pid": true})
	pid, parseErr := strconv.Atoi(values["--parent-pid"])
	if err != nil || parseErr != nil || pid <= 0 || len(positionals) > 0 || values["--path"] == "" {
		fmt.Fprintln(stderr, "upgrade cleanup arguments are invalid")
		return 2
	}
	if err := releasekit.CleanupDeferredHandoff(filepath.Clean(values["--path"]), pid); err != nil {
		fmt.Fprintf(stderr, "upgrade cleanup failed: %v\n", err)
		return 1
	}
	return 0
}
