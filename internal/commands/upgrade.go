package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
	authority, authorityErr := inspectUpgradeSource(observed)
	if authorityErr != nil {
		code := "upgrade_requires_valid_v2_installation"
		message := "upgrade requires a valid lock-v2 installation"
		if state.AsInt(observed.Lock["schema_version"]) == 1 {
			code = "upgrade_requires_clean_legacy_v1_installation"
			message = authorityErr.Error()
		}
		result := blockedLifecycle("upgrade", observed, code, message, "run check and resolve the reported installation state before retrying", "check")
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
	lockAction := "reconcile"
	if authority.LockSchema == 1 {
		lockAction = "migrate"
	}
	result.Changes = append(result.Changes,
		reconcile.Change{Action: "replace", Path: installedBinary, Owner: "installed-binary"},
		reconcile.Change{Action: lockAction, Path: state.LockPath, Owner: "generated-surface"},
	)
	result.Provenance = upgradeProvenance(prepared)
	if dryRun {
		if authority.LockSchema == 1 {
			graph, err := compileObservedGraph(observed)
			if err != nil {
				return nil, reconcile.Result{}, err
			}
			now := time.Now().UTC().Truncate(time.Second)
			lock, err := desiredUpgradeLock(observed, graph, now, targetVersion, prepared.Asset.URL, prepared.Catalog.Location, prepared.Catalog.DigestSHA256, prepared.Pack.ArchiveSHA256, prepared.Pack.PackManifestSHA256, prepared.Pack.ContentManifestSHA256, prepared.Pack.Manifest.BinarySHA256)
			if err != nil {
				return nil, reconcile.Result{}, err
			}
			preview, err := runLifecycle(lifecycleRequest{command: "upgrade", root: root, now: now, dryRun: true, observed: observed, graph: graph, desiredLock: lock, ensureIgnore: true})
			if err != nil {
				return nil, reconcile.Result{}, err
			}
			preview.Changes = append([]reconcile.Change{{Action: "replace", Path: installedBinary, Owner: "installed-binary"}}, preview.Changes...)
			preview.StateAfter = string(state.StateCurrent)
			preview.Provenance = upgradeProvenance(prepared)
			preview.Validations = append(preview.Validations,
				reconcile.Validation{Name: "legacy-lock-v1-integrity", Status: "passed"},
				reconcile.Validation{Name: "source-lock-identity", Status: "passed"},
				reconcile.Validation{Name: "schema-and-version-migration", Status: "passed", Detail: fmt.Sprintf("lock schema 1 -> 2; kit %s -> %s", currentVersion, targetVersion)},
				reconcile.Validation{Name: "catalog-to-archive", Status: "passed"},
				reconcile.Validation{Name: "pack-to-binary", Status: "passed"},
				reconcile.Validation{Name: "staged-version", Status: "passed"},
			)
			payload := resultPayload(preview)
			payload["installed_version"] = currentVersion
			payload["target_version"] = targetVersion
			payload["source_lock_schema"] = 1
			payload["target_lock_schema"] = 2
			payload["migrated_lock"] = true
			return payload, preview, nil
		}
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
	previousBinarySHA, err := regularFileSHA256(installedBinary)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	handoff, err := releasekit.NewHandoff(prepared, root, installedBinary, currentVersion, authority.LockSchema, authority.LockSHA256)
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
			restoredSHA, restoreErr := regularFileSHA256(installedBinary)
			targetRecoveryRequired := releasekit.IsHandoffRecoveryRequired(err)
			restored := !targetRecoveryRequired && restoreErr == nil && restoredSHA == previousBinarySHA && upgradeSourceRestored(root, authority)
			result.Rollback = reconcile.Rollback{Attempted: true, Succeeded: restored, Detail: err.Error()}
			if restored {
				result.Status = reconcile.StatusRolledBack
				result.Blockers = append(result.Blockers, reconcile.Blocker{Code: "upgrade_rolled_back", Message: err.Error(), RetryCommand: "upgrade --dry-run"})
			} else {
				result.Status = reconcile.StatusRecoveryRequired
				result.Blockers = append(result.Blockers, reconcile.Blocker{Code: "recovery_required", Message: "upgrade failed and the previous binary plus repository source authority could not be verified as restored", Remediation: "preserve the handoff and repository transaction evidence before recovery", RetryCommand: "check"})
			}
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

func upgradeSourceRestored(root string, expected upgradeSourceAuthority) bool {
	observed, err := state.Inspect(root)
	if err != nil {
		return false
	}
	actual, err := inspectTargetUpgradeSource(observed)
	return err == nil && actual == expected
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

func regularFileSHA256(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("installed binary is not a regular file")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

func RunUpgradeReconcile(args []string, stdout io.Writer, stderr io.Writer) int {
	values, _, positionals, err := parseValueArgs(args, map[string]bool{
		"--repository":              true,
		"--target-version":          true,
		"--previous-version":        true,
		"--previous-lock-schema":    true,
		"--previous-lock-sha256":    true,
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
	boundProtocol := values["--target-version"] != "" && values["--previous-lock-schema"] != "" && values["--previous-lock-sha256"] != ""
	legacyV2Protocol := values["--target-version"] == "" && values["--previous-lock-schema"] == "" && values["--previous-lock-sha256"] == ""
	previousLockSchema := 0
	var schemaErr error
	if boundProtocol {
		previousLockSchema, schemaErr = strconv.Atoi(values["--previous-lock-schema"])
	}
	observed, err := state.Inspect(root)
	authority, authorityErr := inspectTargetUpgradeSource(observed)
	digestValues := []string{values["--catalog-sha256"], values["--archive-sha256"], values["--pack-manifest-sha256"], values["--content-manifest-sha256"], values["--binary-sha256"]}
	if boundProtocol {
		digestValues = append(digestValues, values["--previous-lock-sha256"])
	}
	digestsValid := true
	for _, digest := range digestValues {
		if !validUpgradeSHA256(strings.ToLower(digest)) {
			digestsValid = false
			break
		}
	}
	protocolAuthorityValid := false
	if boundProtocol {
		protocolAuthorityValid = schemaErr == nil && authorityErr == nil && authority.LockSchema == previousLockSchema && authority.LockSHA256 == strings.ToLower(values["--previous-lock-sha256"])
	} else if legacyV2Protocol {
		protocolAuthorityValid = authorityErr == nil && authority.LockSchema == 2
	}
	if err != nil || (!boundProtocol && !legacyV2Protocol) || (boundProtocol && values["--target-version"] != version.Version) || values["--asset-url"] == "" || values["--catalog-location"] == "" || !digestsValid || !protocolAuthorityValid || state.AsString(observed.Lock["kit_version"]) != values["--previous-version"] || releasekit.RequireForwardUpgrade(values["--previous-version"], version.Version) != nil {
		fmt.Fprintf(stderr, "upgrade reconcile precondition failed\n")
		return 1
	}
	graph, err := compileObservedGraph(observed)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	now := time.Now().UTC().Truncate(time.Second)
	lock, err := desiredUpgradeLock(observed, graph, now, version.Version, values["--asset-url"], values["--catalog-location"], values["--catalog-sha256"], values["--archive-sha256"], values["--pack-manifest-sha256"], values["--content-manifest-sha256"], values["--binary-sha256"])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	result, err := runLifecycle(lifecycleRequest{command: "upgrade", root: root, now: now, observed: observed, graph: graph, desiredLock: lock, ensureIgnore: true, authorityCheck: boundUpgradeSourceAuthorityCheck(observed, authority)})
	if err != nil {
		fmt.Fprintf(stderr, "upgrade reconcile failed: %v\n", err)
		return 1
	}
	if !result.OK() {
		fmt.Fprintf(stderr, "upgrade reconcile failed with status %s\n", result.Status)
		if result.Status == reconcile.StatusRecoveryRequired {
			return 3
		}
		return 1
	}
	return 0
}

func desiredUpgradeLock(observed state.Observed, graph state.Graph, now time.Time, targetVersion, assetURL, catalogLocation, catalogSHA256, archiveSHA256, packManifestSHA256, contentManifestSHA256, binarySHA256 string) (map[string]any, error) {
	lock, err := desiredLifecycleLock("upgrade", observed, graph, now)
	if err != nil {
		return nil, err
	}
	lock["kit_version"] = targetVersion
	lock["release"] = map[string]any{"asset_url": assetURL, "checksum_sha256": archiveSHA256}
	lock["release_provenance"] = map[string]any{
		"verification_status":     "verified",
		"source":                  catalogLocation,
		"catalog_url":             catalogLocation,
		"catalog_sha256":          catalogSHA256,
		"archive_sha256":          archiveSHA256,
		"pack_manifest_sha256":    packManifestSHA256,
		"content_manifest_sha256": contentManifestSHA256,
		"binary_sha256":           binarySHA256,
		"verified_at":             now.UTC().Truncate(time.Second).Format(time.RFC3339),
	}
	return lock, nil
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
