package release

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Handoff struct {
	SchemaVersion         int    `json:"schema_version"`
	TransactionID         string `json:"transaction_id"`
	ParentPID             int    `json:"parent_pid"`
	RepositoryRoot        string `json:"repository_root"`
	TargetBinary          string `json:"target_binary"`
	StagedBinary          string `json:"staged_binary"`
	BinarySHA256          string `json:"binary_sha256"`
	Version               string `json:"version"`
	PreviousVersion       string `json:"previous_version"`
	PreviousLockSchema    int    `json:"previous_lock_schema"`
	PreviousLockSHA256    string `json:"previous_lock_sha256"`
	AssetURL              string `json:"asset_url"`
	CatalogLocation       string `json:"catalog_location"`
	CatalogSHA256         string `json:"catalog_sha256"`
	ArchiveSHA256         string `json:"archive_sha256"`
	PackManifestSHA256    string `json:"pack_manifest_sha256"`
	ContentManifestSHA256 string `json:"content_manifest_sha256"`
	HandoffSHA256         string `json:"handoff_sha256"`
}

// legacyV2HandoffWire is the schema-1 handoff emitted by released CLIs through v0.1.25.
// It is accepted only as the transport for a target-side lock-v2 forward upgrade.
type legacyV2HandoffWire struct {
	SchemaVersion         int    `json:"schema_version"`
	TransactionID         string `json:"transaction_id"`
	ParentPID             int    `json:"parent_pid"`
	RepositoryRoot        string `json:"repository_root"`
	TargetBinary          string `json:"target_binary"`
	StagedBinary          string `json:"staged_binary"`
	BinarySHA256          string `json:"binary_sha256"`
	Version               string `json:"version"`
	PreviousVersion       string `json:"previous_version"`
	AssetURL              string `json:"asset_url"`
	CatalogLocation       string `json:"catalog_location"`
	CatalogSHA256         string `json:"catalog_sha256"`
	ArchiveSHA256         string `json:"archive_sha256"`
	PackManifestSHA256    string `json:"pack_manifest_sha256"`
	ContentManifestSHA256 string `json:"content_manifest_sha256"`
	HandoffSHA256         string `json:"handoff_sha256"`
}

type HandoffReconcileError struct {
	ExitStatus int
	Cause      error
	Output     string
}

func (err *HandoffReconcileError) Error() string {
	return fmt.Sprintf("new binary reconciliation failed with exit %d: %v: %s", err.ExitStatus, err.Cause, err.Output)
}

func (err *HandoffReconcileError) Unwrap() error { return err.Cause }

func IsHandoffRecoveryRequired(err error) bool {
	var reconcileErr *HandoffReconcileError
	return errors.As(err, &reconcileErr) && reconcileErr.ExitStatus == 3
}

func NewHandoff(prepared PreparedUpgrade, repositoryRoot, targetBinary, previousVersion string, previousLockSchema int, previousLockSHA256 string) (Handoff, error) {
	root, err := filepath.Abs(repositoryRoot)
	if err != nil {
		return Handoff{}, err
	}
	target, err := filepath.Abs(targetBinary)
	if err != nil {
		return Handoff{}, err
	}
	assetURL, err := resolveLocation(prepared.Asset.URL, prepared.Catalog.Base)
	if err != nil {
		return Handoff{}, err
	}
	digest := sha256.Sum256([]byte(prepared.Pack.BinaryPath + "\x00" + target + "\x00" + prepared.Asset.Version))
	handoff := Handoff{
		SchemaVersion: 1, TransactionID: hex.EncodeToString(digest[:16]), ParentPID: os.Getpid(),
		RepositoryRoot: root, TargetBinary: target, StagedBinary: prepared.Pack.BinaryPath,
		BinarySHA256: prepared.Pack.Manifest.BinarySHA256, Version: prepared.Asset.Version,
		PreviousVersion: previousVersion, PreviousLockSchema: previousLockSchema, PreviousLockSHA256: previousLockSHA256,
		AssetURL: assetURL, CatalogLocation: prepared.Catalog.Location,
		CatalogSHA256: prepared.Catalog.DigestSHA256, ArchiveSHA256: prepared.Pack.ArchiveSHA256,
		PackManifestSHA256:    prepared.Pack.PackManifestSHA256,
		ContentManifestSHA256: prepared.Pack.ContentManifestSHA256,
	}
	handoff.HandoffSHA256 = wireHandoffDigest(handoff)
	return handoff, nil
}

func ApplyHandoff(handoff Handoff) (resultErr error) {
	if err := validateHandoff(handoff); err != nil {
		return err
	}
	legacyProtocol := isLegacyV2Handoff(handoff)
	if !legacyProtocol {
		if err := validateSourceLockIdentity(handoff); err != nil {
			return err
		}
	}
	stageDir := filepath.Join(filepath.Dir(handoff.TargetBinary), ".upgrade-"+handoff.TransactionID)
	if err := os.MkdirAll(stageDir, 0o700); err != nil {
		return err
	}
	cleanupStage := true
	defer func() {
		if cleanupStage {
			_ = os.RemoveAll(stageDir)
		}
	}()
	stagedCopy := filepath.Join(stageDir, filepath.Base(handoff.TargetBinary))
	if err := copyFile(handoff.StagedBinary, stagedCopy, 0o755); err != nil {
		return err
	}
	handoff.StagedBinary = stagedCopy
	handoff.HandoffSHA256 = wireHandoffDigest(handoff)
	if err := validateHandoff(handoff); err != nil {
		return err
	}
	backup := filepath.Join(stageDir, "previous-binary")
	if err := os.Rename(handoff.TargetBinary, backup); err != nil {
		return fmt.Errorf("back up installed binary: %w", err)
	}
	restore := true
	defer func() {
		if restore {
			if err := restorePreviousBinary(backup, handoff.TargetBinary); err != nil {
				cleanupStage = false
				resultErr = fmt.Errorf("%v; restore previous binary failed: %w; backup preserved at %s", resultErr, err, backup)
			}
		}
	}()
	if err := copyFile(stagedCopy, handoff.TargetBinary, 0o755); err != nil {
		return err
	}
	arguments := []string{
		"__upgrade-reconcile", "--repository", handoff.RepositoryRoot,
		"--previous-version", handoff.PreviousVersion,
		"--asset-url", handoff.AssetURL,
		"--catalog-location", handoff.CatalogLocation,
		"--catalog-sha256", handoff.CatalogSHA256,
		"--archive-sha256", handoff.ArchiveSHA256,
		"--pack-manifest-sha256", handoff.PackManifestSHA256,
		"--content-manifest-sha256", handoff.ContentManifestSHA256,
		"--binary-sha256", handoff.BinarySHA256,
	}
	if !legacyProtocol {
		arguments = append(arguments,
			"--target-version", handoff.Version,
			"--previous-lock-schema", fmt.Sprint(handoff.PreviousLockSchema),
			"--previous-lock-sha256", handoff.PreviousLockSHA256,
		)
	}
	command := exec.Command(handoff.TargetBinary, arguments...)
	command.Env = append(os.Environ(), "CODEHEART_OPERATING_KIT_CLI=1")
	if output, err := command.CombinedOutput(); err != nil {
		exitStatus := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitStatus = exitErr.ExitCode()
		}
		return &HandoffReconcileError{ExitStatus: exitStatus, Cause: err, Output: strings.TrimSpace(string(output))}
	}
	restore = false
	return nil
}

func restorePreviousBinary(backup, target string) error {
	attempts := 1
	if runtime.GOOS == "windows" {
		attempts = 300
	}
	var removeErr, renameErr error
	for attempt := 0; attempt < attempts; attempt++ {
		removeErr = os.Remove(target)
		if os.IsNotExist(removeErr) {
			removeErr = nil
		}
		if removeErr == nil {
			renameErr = os.Rename(backup, target)
			if renameErr == nil {
				return nil
			}
		}
		if attempt+1 < attempts {
			time.Sleep(100 * time.Millisecond)
		}
	}
	return fmt.Errorf("remove target: %v; rename backup: %v", removeErr, renameErr)
}

func WriteHandoff(path string, handoff Handoff) error {
	data, err := json.MarshalIndent(handoff, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func ReadHandoff(path string) (Handoff, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Handoff{}, err
	}
	var handoff Handoff
	if err := decodeStrictJSON(data, &handoff); err != nil {
		return Handoff{}, err
	}
	if err := validateHandoff(handoff); err != nil {
		return Handoff{}, err
	}
	absolute, err := filepath.Abs(path)
	if err != nil || filepath.Base(absolute) != "handoff.json" || !strings.HasPrefix(filepath.Base(filepath.Dir(absolute)), ".upgrade-handoff-") || filepath.Dir(filepath.Dir(absolute)) != filepath.Dir(handoff.TargetBinary) {
		return Handoff{}, fmt.Errorf("upgrade handoff file is outside its target binary directory")
	}
	if filepath.Dir(handoff.StagedBinary) != filepath.Dir(absolute) {
		return Handoff{}, fmt.Errorf("upgrade handoff staged binary is outside its handoff directory")
	}
	return handoff, nil
}

func StartDeferredHandoff(handoff Handoff) error {
	stageDir := filepath.Join(filepath.Dir(handoff.TargetBinary), ".upgrade-handoff-"+handoff.TransactionID)
	if err := os.MkdirAll(stageDir, 0o700); err != nil {
		return err
	}
	staged := filepath.Join(stageDir, filepath.Base(handoff.TargetBinary))
	if err := copyFile(handoff.StagedBinary, staged, 0o755); err != nil {
		return err
	}
	handoff.StagedBinary = staged
	handoff.HandoffSHA256 = wireHandoffDigest(handoff)
	file := filepath.Join(stageDir, "handoff.json")
	if err := WriteHandoff(file, handoff); err != nil {
		return err
	}
	command := exec.Command(staged, "__upgrade-handoff", "--file", file)
	command.Env = os.Environ()
	if err := command.Start(); err != nil {
		return err
	}
	return command.Process.Release()
}

func ExecuteDeferredHandoff(path string) error {
	handoff, err := ReadHandoff(path)
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		parentExited := false
		for attempt := 0; attempt < 300; attempt++ {
			alive, err := processExists(handoff.ParentPID)
			if err != nil {
				return err
			}
			if !alive {
				parentExited = true
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if !parentExited {
			return fmt.Errorf("parent process %d did not exit before handoff timeout", handoff.ParentPID)
		}
	}
	if err := ApplyHandoff(handoff); err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		command := exec.Command(handoff.TargetBinary, "__cleanup-upgrade-handoff", "--path", filepath.Dir(path), "--parent-pid", fmt.Sprint(os.Getpid()))
		if err := command.Start(); err != nil {
			return err
		}
		return command.Process.Release()
	}
	return os.RemoveAll(filepath.Dir(path))
}

func CleanupDeferredHandoff(path string, parentPID int) error {
	if !strings.HasPrefix(filepath.Base(path), ".upgrade-handoff-") {
		return fmt.Errorf("cleanup path is not an upgrade handoff directory")
	}
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	path, err = filepath.Abs(path)
	if err != nil || filepath.Dir(path) != filepath.Dir(executable) {
		return fmt.Errorf("cleanup path is outside the installed binary directory")
	}
	for attempt := 0; attempt < 300; attempt++ {
		alive, err := processExists(parentPID)
		if err != nil {
			return err
		}
		if !alive {
			return os.RemoveAll(path)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("handoff process %d did not exit before cleanup timeout", parentPID)
}

func validateHandoff(handoff Handoff) error {
	legacyProtocol := isLegacyV2Handoff(handoff)
	newProtocol := (handoff.PreviousLockSchema == 1 || handoff.PreviousLockSchema == 2) && validDigest(handoff.PreviousLockSHA256)
	if handoff.SchemaVersion != 1 || handoff.TransactionID == "" || handoff.RepositoryRoot == "" || handoff.TargetBinary == "" || handoff.StagedBinary == "" || handoff.Version == "" || handoff.PreviousVersion == "" || (!legacyProtocol && !newProtocol) || handoff.AssetURL == "" || handoff.CatalogLocation == "" || RequireForwardUpgrade(handoff.PreviousVersion, handoff.Version) != nil {
		return fmt.Errorf("upgrade handoff identity is incomplete")
	}
	digests := []string{handoff.BinarySHA256, handoff.CatalogSHA256, handoff.ArchiveSHA256, handoff.PackManifestSHA256, handoff.ContentManifestSHA256}
	if !legacyProtocol {
		digests = append(digests, handoff.PreviousLockSHA256)
	}
	for _, digest := range digests {
		if !validDigest(digest) {
			return fmt.Errorf("upgrade handoff digest identity is incomplete")
		}
	}
	if handoff.HandoffSHA256 == "" || handoff.HandoffSHA256 != wireHandoffDigest(handoff) {
		return fmt.Errorf("upgrade handoff metadata identity mismatch")
	}
	for label, path := range map[string]string{"target": handoff.TargetBinary, "staged": handoff.StagedBinary} {
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("upgrade handoff %s binary is not a regular file", label)
		}
	}
	digest, err := fileSHA256(handoff.StagedBinary)
	if err != nil || digest != handoff.BinarySHA256 {
		return fmt.Errorf("upgrade handoff binary identity mismatch")
	}
	return nil
}

func validateSourceLockIdentity(handoff Handoff) error {
	lockPath := filepath.Join(handoff.RepositoryRoot, ".codeheart", "kit.lock.yaml")
	info, err := os.Lstat(lockPath)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("upgrade source lock is not a regular file")
	}
	digest, err := fileSHA256(lockPath)
	if err != nil || !strings.EqualFold(digest, handoff.PreviousLockSHA256) {
		return fmt.Errorf("upgrade source lock identity changed before handoff")
	}
	return nil
}

func validDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func handoffDigest(handoff Handoff) string {
	handoff.HandoffSHA256 = ""
	data, _ := json.Marshal(handoff)
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func isLegacyV2Handoff(handoff Handoff) bool {
	return handoff.PreviousLockSchema == 0 && handoff.PreviousLockSHA256 == ""
}

func wireHandoffDigest(handoff Handoff) string {
	if !isLegacyV2Handoff(handoff) {
		return handoffDigest(handoff)
	}
	legacy := legacyV2Wire(handoff)
	legacy.HandoffSHA256 = ""
	data, _ := json.Marshal(legacy)
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func legacyV2Wire(handoff Handoff) legacyV2HandoffWire {
	return legacyV2HandoffWire{
		SchemaVersion: handoff.SchemaVersion, TransactionID: handoff.TransactionID, ParentPID: handoff.ParentPID,
		RepositoryRoot: handoff.RepositoryRoot, TargetBinary: handoff.TargetBinary, StagedBinary: handoff.StagedBinary,
		BinarySHA256: handoff.BinarySHA256, Version: handoff.Version, PreviousVersion: handoff.PreviousVersion,
		AssetURL: handoff.AssetURL, CatalogLocation: handoff.CatalogLocation, CatalogSHA256: handoff.CatalogSHA256,
		ArchiveSHA256: handoff.ArchiveSHA256, PackManifestSHA256: handoff.PackManifestSHA256,
		ContentManifestSHA256: handoff.ContentManifestSHA256, HandoffSHA256: handoff.HandoffSHA256,
	}
}

func copyFile(source, target string, mode os.FileMode) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
