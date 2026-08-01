package reconcile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type PhaseHook func(phase string) error

type ApplyOptions struct {
	Now  time.Time
	Hook PhaseHook
}

type transactionMarker struct {
	SchemaVersion int    `json:"schema_version"`
	TransactionID string `json:"transaction_id"`
	Command       string `json:"command"`
	Phase         string `json:"phase"`
	PID           int    `json:"pid"`
	StartedAt     string `json:"started_at"`
	TargetRoot    string `json:"target_root"`
	RecoveryPath  string `json:"recovery_path"`
	Error         string `json:"error,omitempty"`
}

type committedAction struct {
	action         Action
	parentRoot     *os.Root
	backupFile     *os.File
	targetName     string
	backup         string
	parentIdentity pathIdentity
	targetInfo     os.FileInfo
	originalInfo   os.FileInfo
	backupMode     os.FileMode
	backupSHA256   string
	backupInfo     os.FileInfo
	originalStash  string
	hadTarget      bool
	mutated        bool
	installed      bool
	committed      bool
}

type pathIdentity struct {
	path string
	info os.FileInfo
}

type boundFileEvidence struct {
	path string
	file *os.File
	info os.FileInfo
	data []byte
}

var errRecoveryPathAuthority = errors.New("recovery path authority changed")

// Apply executes a fully planned transaction. Expected operation failures are returned as a
// structured result; setup and programming errors are returned as errors.
func Apply(plan Plan, options ApplyOptions) (Result, error) {
	result := NewResult(plan.Command)
	result.StateBefore = plan.StateBefore
	result.TransactionID = plan.ID
	result.Changes = planChanges(plan)
	if len(plan.Blockers) > 0 {
		result.Status = StatusBlocked
		result.Blockers = append(result.Blockers, plan.Blockers...)
		return result, nil
	}
	if len(plan.Actions) == 0 {
		result.Status = StatusSucceeded
		result.StateAfter = plan.StateBefore
		result.Validations = append(result.Validations, Validation{Name: "no-op", Status: "passed"})
		return result, nil
	}
	now := options.Now
	if now.IsZero() {
		now = time.Now()
	}
	now = now.UTC().Truncate(time.Second)
	rootCreated := false
	if _, err := os.Stat(plan.Root); os.IsNotExist(err) {
		if err := os.MkdirAll(plan.Root, 0o755); err != nil {
			return result, err
		}
		rootCreated = true
	}
	canonicalRoot, err := filepath.EvalSymlinks(plan.Root)
	if err != nil {
		return result, err
	}
	canonicalRoot, err = filepath.Abs(canonicalRoot)
	if err != nil {
		return result, err
	}
	repositoryRoot, err := os.OpenRoot(canonicalRoot)
	if err != nil {
		return result, err
	}
	defer repositoryRoot.Close()
	identities := map[string]pathIdentity{}
	for _, action := range plan.Actions {
		if (action.Kind == "replace" || action.Kind == "remove") && !validSHA256(action.ExpectedSHA256) {
			result.Status = StatusBlocked
			result.Blockers = append(result.Blockers, Blocker{Code: "missing_precondition", Message: "replace and remove actions require an exact source SHA-256 precondition", Path: action.Target})
			return result, nil
		}
		if _, err := safeTarget(canonicalRoot, action.Target); err != nil {
			result.Status = StatusBlocked
			result.Blockers = append(result.Blockers, Blocker{Code: "unsafe_target", Message: err.Error(), Path: action.Target})
			if rootCreated {
				_ = os.Remove(plan.Root)
			}
			return result, nil
		}
		identity, err := nearestExistingParent(canonicalRoot, action.Target)
		if err != nil {
			return result, err
		}
		identities[action.Target] = identity
	}
	transactionRelative := filepath.ToSlash(filepath.Join(".codeheart", "local", "kit-transactions", plan.ID))
	markerRelative := filepath.ToSlash(filepath.FromSlash(state.TransactionPath))
	if _, err := safeTarget(canonicalRoot, transactionRelative); err != nil {
		result.Status = StatusBlocked
		result.Blockers = append(result.Blockers, Blocker{Code: "unsafe_transaction_path", Message: err.Error(), Path: transactionRelative})
		return result, nil
	}
	marker := transactionMarker{
		SchemaVersion: 1,
		TransactionID: plan.ID,
		Command:       plan.Command,
		Phase:         "planning",
		PID:           os.Getpid(),
		StartedAt:     now.Format(time.RFC3339),
		TargetRoot:    canonicalRoot,
		RecoveryPath:  filepath.ToSlash(filepath.Join(".codeheart", "local", "kit-transactions", plan.ID)),
	}
	markerFile, markerInfo, err := acquireMarkerRoot(repositoryRoot, markerRelative, marker)
	if err != nil {
		result.Status = StatusBlocked
		result.Blockers = append(result.Blockers, Blocker{
			Code:         "transaction_in_progress",
			Message:      err.Error(),
			Path:         state.TransactionPath,
			Remediation:  "run check, then repair after the active process finishes",
			RetryCommand: plan.Command,
		})
		if rootCreated {
			_ = os.Remove(plan.Root)
		}
		return result, nil
	}
	cleanupMarker := true
	defer func() {
		if cleanupMarker {
			_ = removeBoundMarker(repositoryRoot, markerRelative, markerFile, markerInfo)
		} else {
			_ = markerFile.Close()
		}
	}()
	rollbackUnbound := func(cause error) (Result, error) {
		return rollbackBeforeCommitRoot(result, repositoryRoot, markerFile, markerInfo, markerRelative, transactionRelative, marker, cause, rootCreated, plan.Root, nil, nil)
	}
	if err := invokeHook(options.Hook, "marker"); err != nil {
		return rollbackUnbound(err)
	}
	if _, err := safeTarget(canonicalRoot, transactionRelative); err != nil {
		return rollbackUnbound(err)
	}
	transactionParent := filepath.Dir(filepath.FromSlash(transactionRelative))
	if err := repositoryRoot.MkdirAll(transactionParent, 0o700); err != nil {
		return rollbackUnbound(err)
	}
	if _, err := safeTarget(canonicalRoot, transactionRelative); err != nil {
		return rollbackUnbound(err)
	}
	if err := repositoryRoot.Mkdir(filepath.FromSlash(transactionRelative), 0o700); err != nil {
		if !os.IsExist(err) {
			return rollbackUnbound(err)
		}
		result.Status = StatusBlocked
		result.Blockers = append(result.Blockers, Blocker{
			Code:         "transaction_path_exists",
			Message:      fmt.Sprintf("exclusive transaction directory acquisition failed: %v", err),
			Path:         transactionRelative,
			Remediation:  "preserve and inspect the existing transaction directory, then run check and repair",
			RetryCommand: "repair",
		})
		return result, nil
	}
	transactionRoot, err := repositoryRoot.OpenRoot(transactionRelative)
	if err != nil {
		return rollbackUnbound(err)
	}
	defer transactionRoot.Close()
	transactionInfo, err := transactionRoot.Stat(".")
	if err != nil {
		return rollbackBeforeCommitRoot(result, repositoryRoot, markerFile, markerInfo, markerRelative, transactionRelative, marker, err, rootCreated, plan.Root, transactionRoot, nil)
	}
	rollbackBound := func(cause error) (Result, error) {
		return rollbackBeforeCommitRoot(result, repositoryRoot, markerFile, markerInfo, markerRelative, transactionRelative, marker, cause, rootCreated, plan.Root, transactionRoot, transactionInfo)
	}
	if _, err := safeTarget(canonicalRoot, transactionRelative); err != nil {
		return rollbackBound(err)
	}
	transactionPathInfo, err := os.Lstat(filepath.Join(canonicalRoot, filepath.FromSlash(transactionRelative)))
	if err != nil || !os.SameFile(transactionInfo, transactionPathInfo) {
		return rollbackBound(fmt.Errorf("transaction root identity changed during binding"))
	}
	if err := transactionRoot.Mkdir("stage", 0o700); err != nil {
		return rollbackBound(err)
	}
	if err := transactionRoot.Mkdir("backup", 0o700); err != nil {
		return rollbackBound(err)
	}
	planData, _ := json.MarshalIndent(plan, "", "  ")
	if err := transactionRoot.WriteFile("plan.json", append(planData, '\n'), 0o600); err != nil {
		return rollbackBound(err)
	}
	for _, action := range plan.Actions {
		if action.Kind == "remove" {
			continue
		}
		stagePath := filepath.ToSlash(filepath.Join("stage", filepath.FromSlash(action.Target)))
		if err := transactionRoot.MkdirAll(filepath.ToSlash(filepath.Dir(stagePath)), 0o700); err != nil {
			return rollbackBound(err)
		}
		mode := os.FileMode(action.Mode)
		if mode == 0 {
			mode = 0o644
		}
		if err := transactionRoot.WriteFile(stagePath, action.Content, mode); err != nil {
			return rollbackBound(err)
		}
	}
	marker.Phase = "staged"
	if err := writeMarkerBound(repositoryRoot, markerRelative, markerFile, markerInfo, marker); err != nil {
		return rollbackBound(err)
	}
	if err := invokeHook(options.Hook, "staged"); err != nil {
		return rollbackBound(err)
	}
	if err := validateStagedRoot(plan, transactionRoot); err != nil {
		return rollbackBound(err)
	}
	result.Validations = append(result.Validations, Validation{Name: "staged-state", Status: "passed"})
	marker.Phase = "validated"
	if err := writeMarkerBound(repositoryRoot, markerRelative, markerFile, markerInfo, marker); err != nil {
		return rollbackBound(err)
	}
	if err := invokeHook(options.Hook, "validated"); err != nil {
		return rollbackBound(err)
	}

	marker.Phase = "committing"
	if err := writeMarkerBound(repositoryRoot, markerRelative, markerFile, markerInfo, marker); err != nil {
		return rollbackBound(err)
	}
	committed := []committedAction{}
	defer closeCommittedRoots(&committed)
	for _, action := range plan.Actions {
		if err := revalidateParentIdentity(identities[action.Target]); err != nil {
			result, keepMarker := rollbackCommitted(result, repositoryRoot, transactionRoot, transactionInfo, markerFile, markerInfo, transactionRelative, markerRelative, marker, committed, err, options.Hook)
			cleanupMarker = !keepMarker
			return result, nil
		}
		if err := invokeHook(options.Hook, "commit:"+action.Target); err != nil {
			result, keepMarker := rollbackCommitted(result, repositoryRoot, transactionRoot, transactionInfo, markerFile, markerInfo, transactionRelative, markerRelative, marker, committed, err, options.Hook)
			cleanupMarker = !keepMarker
			return result, nil
		}
		if err := invokeHook(options.Hook, "pre-commit:"+action.Target); err != nil {
			result, keepMarker := rollbackCommitted(result, repositoryRoot, transactionRoot, transactionInfo, markerFile, markerInfo, transactionRelative, markerRelative, marker, committed, err, options.Hook)
			cleanupMarker = !keepMarker
			return result, nil
		}
		record, err := commitAction(canonicalRoot, repositoryRoot, transactionRoot, identities[action.Target], action, options.Hook)
		if err != nil {
			if record.mutated {
				committed = append(committed, record)
			}
			result, keepMarker := rollbackCommitted(result, repositoryRoot, transactionRoot, transactionInfo, markerFile, markerInfo, transactionRelative, markerRelative, marker, committed, err, options.Hook)
			cleanupMarker = !keepMarker
			return result, nil
		}
		committed = append(committed, record)
	}
	marker.Phase = "post-check"
	if err := writeMarkerBound(repositoryRoot, markerRelative, markerFile, markerInfo, marker); err != nil {
		result, keepMarker := rollbackCommitted(result, repositoryRoot, transactionRoot, transactionInfo, markerFile, markerInfo, transactionRelative, markerRelative, marker, committed, err, options.Hook)
		cleanupMarker = !keepMarker
		return result, nil
	}
	if err := invokeHook(options.Hook, "post-check"); err != nil {
		result, keepMarker := rollbackCommitted(result, repositoryRoot, transactionRoot, transactionInfo, markerFile, markerInfo, transactionRelative, markerRelative, marker, committed, err, options.Hook)
		cleanupMarker = !keepMarker
		return result, nil
	}
	observed, err := state.InspectIgnoringTransaction(canonicalRoot)
	if err != nil || !expectedClassification(plan.ExpectedAfter, observed.Classification) {
		if err == nil {
			err = fmt.Errorf("post-check classified state as %s", observed.Classification)
		}
		result, keepMarker := rollbackCommitted(result, repositoryRoot, transactionRoot, transactionInfo, markerFile, markerInfo, transactionRelative, markerRelative, marker, committed, err, options.Hook)
		cleanupMarker = !keepMarker
		return result, nil
	}
	result.Validations = append(result.Validations, Validation{Name: "post-check", Status: "passed"})
	result.StateAfter = string(observed.Classification)
	markerAuthorityErr := validateBoundMarker(repositoryRoot, markerRelative, markerFile, markerInfo)
	transactionAuthorityErr := validateBoundTransaction(repositoryRoot, transactionRoot, transactionRelative, transactionInfo)
	if markerAuthorityErr != nil || transactionAuthorityErr != nil {
		details := []string{"committed changes passed post-check but cleanup authority changed"}
		if markerAuthorityErr != nil {
			details = append(details, "marker authority: "+markerAuthorityErr.Error())
		}
		if transactionAuthorityErr != nil {
			details = append(details, "transaction authority: "+transactionAuthorityErr.Error())
		}
		marker.Phase = "recovery-required"
		marker.Error = strings.Join(details, "; ")
		if markerAuthorityErr == nil {
			if err := writeMarkerBound(repositoryRoot, markerRelative, markerFile, markerInfo, marker); err != nil {
				details = append(details, "marker: "+err.Error())
			}
		}
		result.Status = StatusRecoveryRequired
		result.Validations = append(result.Validations, Validation{Name: "cleanup-authority", Status: "failed", Detail: strings.Join(details, "; ")})
		result.Blockers = append(result.Blockers, Blocker{Code: "recovery_required", Message: strings.Join(details, "; "), Path: state.TransactionPath, Remediation: "preserve marker and transaction evidence, then run check before repair", RetryCommand: "repair"})
		cleanupMarker = false
		return result, nil
	}
	if err := removeBoundMarkerWithHook(repositoryRoot, markerRelative, markerFile, markerInfo, options.Hook); err != nil {
		result.Status = StatusRecoveryRequired
		result.Validations = append(result.Validations, Validation{Name: "marker-cleanup", Status: "failed", Detail: err.Error()})
		result.Blockers = append(result.Blockers, Blocker{Code: "recovery_required", Message: "committed changes passed post-check but marker cleanup failed: " + err.Error(), Path: state.TransactionPath, Remediation: "preserve marker and transaction evidence, then run check before repair", RetryCommand: "repair"})
		cleanupMarker = false
		return result, nil
	}
	cleanupMarker = false
	closeCommittedRoots(&committed)
	if err := cleanupBoundTransactionWithHook(repositoryRoot, transactionRoot, transactionRelative, transactionInfo, options.Hook); err != nil {
		result.Status = StatusRecoveryRequired
		result.Validations = append(result.Validations, Validation{Name: "transaction-cleanup", Status: "failed", Detail: err.Error()})
		result.Blockers = append(result.Blockers, Blocker{Code: "recovery_required", Message: "committed changes passed post-check but transaction cleanup failed: " + err.Error(), Path: transactionRelative, Remediation: "preserve transaction evidence and run check before repair", RetryCommand: "repair"})
		return result, nil
	} else {
		result.Validations = append(result.Validations, Validation{Name: "transaction-cleanup", Status: "passed"})
	}
	result.Validations = append(result.Validations, Validation{Name: "marker-cleanup", Status: "passed"})
	result.Status = StatusSucceeded
	_ = invokeHook(options.Hook, "cleanup")
	return result, nil
}

func expectedClassification(expected []state.Classification, actual state.Classification) bool {
	for _, candidate := range expected {
		if candidate == actual {
			return true
		}
	}
	return false
}

func planChanges(plan Plan) []Change {
	changes := make([]Change, len(plan.Actions))
	for index, action := range plan.Actions {
		changes[index] = Change{Action: action.Kind, Path: action.Target, Owner: action.Owner}
	}
	return changes
}

func acquireMarkerRoot(root *os.Root, path string, marker transactionMarker) (*os.File, os.FileInfo, error) {
	path = filepath.FromSlash(path)
	if err := root.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, nil, err
	}
	if _, err := root.Lstat(path); err == nil {
		return nil, nil, fmt.Errorf("transaction marker already exists; run check and explicit repair recovery")
	} else if !os.IsNotExist(err) {
		return nil, nil, err
	}
	file, err := root.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return nil, nil, err
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, nil, err
	}
	data, _ := json.MarshalIndent(marker, "", "  ")
	if _, err := file.Write(append(data, '\n')); err != nil {
		if cleanupErr := removeBoundMarker(root, path, file, info); cleanupErr != nil {
			return nil, nil, fmt.Errorf("write transaction marker: %v; preserve failed marker: %w", err, cleanupErr)
		}
		return nil, nil, err
	}
	if err := file.Sync(); err != nil {
		if cleanupErr := removeBoundMarker(root, path, file, info); cleanupErr != nil {
			return nil, nil, fmt.Errorf("sync transaction marker: %v; preserve failed marker: %w", err, cleanupErr)
		}
		return nil, nil, err
	}
	return file, info, nil
}

// RecoverStaleTransaction removes only a verified, dead, pre-commit transaction. It never
// takes over an active process or a transaction that might have started committing target bytes.
func RecoverStaleTransaction(root string, dryRun bool) (bool, error) {
	canonical, err := filepath.Abs(root)
	if err != nil {
		return false, err
	}
	canonical, err = filepath.EvalSymlinks(canonical)
	if err != nil {
		return false, err
	}
	repositoryRoot, err := os.OpenRoot(canonical)
	if err != nil {
		return false, err
	}
	defer repositoryRoot.Close()
	return recoverStaleMarkerRoot(repositoryRoot, filepath.FromSlash(state.TransactionPath), canonical, dryRun)
}

func recoverStaleMarkerRoot(root *os.Root, markerPath, canonicalRoot string, dryRun bool) (bool, error) {
	return recoverStaleMarkerRootWithHook(root, markerPath, canonicalRoot, dryRun, nil)
}

func recoverStaleMarkerRootWithHook(root *os.Root, markerPath, canonicalRoot string, dryRun bool, hook PhaseHook) (bool, error) {
	markerPath = filepath.FromSlash(markerPath)
	markerFile, markerInfo, data, err := openBoundRegularRootFile(root, markerPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer markerFile.Close()
	var existing transactionMarker
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&existing); err != nil {
		return false, fmt.Errorf("existing transaction marker is invalid: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return false, fmt.Errorf("existing transaction marker has trailing data")
	}
	if existing.SchemaVersion != 1 || !validTransactionID(existing.TransactionID) || existing.Command == "" || existing.PID <= 0 {
		return false, fmt.Errorf("existing transaction marker identity is incomplete")
	}
	if !sameCanonicalPath(existing.TargetRoot, canonicalRoot) {
		return false, fmt.Errorf("existing transaction marker targets a different root")
	}
	expectedRecovery := filepath.ToSlash(filepath.Join(".codeheart", "local", "kit-transactions", existing.TransactionID))
	if existing.RecoveryPath != expectedRecovery {
		return false, fmt.Errorf("existing transaction recovery path does not match its identity")
	}
	recoveryPath := filepath.FromSlash(existing.RecoveryPath)
	transactionParent := filepath.Join(".codeheart", "local", "kit-transactions")
	if filepath.Clean(recoveryPath) != filepath.Join(transactionParent, existing.TransactionID) || filepath.Dir(recoveryPath) != transactionParent || filepath.Base(recoveryPath) != existing.TransactionID {
		return false, fmt.Errorf("existing transaction recovery path is not one direct transaction child")
	}
	if _, err := safeTarget(root.Name(), recoveryPath); err != nil {
		return false, err
	}
	var transactionRoot *os.Root
	var transactionInfo os.FileInfo
	transactionRoot, err = root.OpenRoot(recoveryPath)
	if err == nil {
		defer transactionRoot.Close()
		transactionInfo, err = transactionRoot.Stat(".")
		if err != nil {
			return false, err
		}
		if err := validateBoundTransaction(root, transactionRoot, recoveryPath, transactionInfo); err != nil {
			return false, err
		}
	} else if !os.IsNotExist(err) {
		return false, err
	}
	alive, err := processAlive(existing.PID)
	if err != nil {
		return false, fmt.Errorf("cannot verify transaction process %d: %w", existing.PID, err)
	}
	if alive {
		return false, fmt.Errorf("transaction %s is still owned by live process %d", existing.TransactionID, existing.PID)
	}
	switch existing.Phase {
	case "planning", "staged", "validated", "failed":
	default:
		return false, fmt.Errorf("stale transaction %s reached phase %s and requires recovery", existing.TransactionID, existing.Phase)
	}
	planRequired := existing.Phase == "staged" || existing.Phase == "validated"
	if transactionRoot == nil && planRequired {
		return false, fmt.Errorf("stale transaction directory is missing")
	}
	var planFile *os.File
	var planInfo os.FileInfo
	var planData []byte
	if transactionRoot != nil {
		var planErr error
		planFile, planInfo, planData, planErr = openBoundRegularRootFile(transactionRoot, "plan.json")
		if planErr == nil {
			defer planFile.Close()
			var plan Plan
			if err := json.Unmarshal(planData, &plan); err != nil || plan.SchemaVersion != 1 || !validTransactionID(plan.ID) || plan.ID != existing.TransactionID || !sameCanonicalPath(plan.Root, canonicalRoot) || plan.Command != existing.Command {
				return false, fmt.Errorf("stale transaction plan identity does not match its marker")
			}
		} else if os.IsNotExist(planErr) {
			return false, fmt.Errorf("stale transaction plan is missing")
		} else if !os.IsNotExist(planErr) {
			return false, planErr
		}
	}
	if err := validateBoundMarker(root, markerPath, markerFile, markerInfo); err != nil {
		return false, err
	}
	if err := validateBoundRootFile(root, markerPath, markerFile, markerInfo, data); err != nil {
		return false, err
	}
	if transactionRoot != nil {
		if err := validateBoundTransaction(root, transactionRoot, recoveryPath, transactionInfo); err != nil {
			return false, err
		}
		if planFile != nil {
			if err := validateBoundRootFile(transactionRoot, "plan.json", planFile, planInfo, planData); err != nil {
				return false, err
			}
		}
	}
	if dryRun {
		return true, nil
	}
	if err := removeBoundMarkerWithEvidence(root, markerPath, markerFile, markerInfo, data, hook); err != nil {
		return false, err
	}
	if transactionRoot != nil {
		evidence := &boundFileEvidence{path: "plan.json", file: planFile, info: planInfo, data: planData}
		if err := cleanupBoundTransactionWithEvidence(root, transactionRoot, recoveryPath, transactionInfo, evidence, hook); err != nil {
			return false, err
		}
	}
	return true, nil
}

func openBoundRegularRootFile(root *os.Root, path string) (*os.File, os.FileInfo, []byte, error) {
	path = filepath.FromSlash(path)
	pathInfo, err := root.Lstat(path)
	if err != nil {
		return nil, nil, nil, err
	}
	if !pathInfo.Mode().IsRegular() || pathInfo.Mode()&os.ModeSymlink != 0 {
		return nil, nil, nil, fmt.Errorf("bound source is not a regular file: %s", filepath.ToSlash(path))
	}
	file, err := root.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	openInfo, err := file.Stat()
	if err != nil || !os.SameFile(pathInfo, openInfo) {
		_ = file.Close()
		return nil, nil, nil, fmt.Errorf("bound source identity changed while opening: %s", filepath.ToSlash(path))
	}
	data, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return nil, nil, nil, err
	}
	currentInfo, err := root.Lstat(path)
	if err != nil || !os.SameFile(pathInfo, currentInfo) {
		_ = file.Close()
		return nil, nil, nil, fmt.Errorf("bound source identity changed while reading: %s", filepath.ToSlash(path))
	}
	return file, pathInfo, data, nil
}

func validateBoundRootFile(root *os.Root, path string, file *os.File, identity os.FileInfo, expected []byte) error {
	path = filepath.FromSlash(path)
	pathInfo, err := root.Lstat(path)
	if err != nil || identity == nil || !os.SameFile(identity, pathInfo) {
		return fmt.Errorf("bound source path identity changed: %s", filepath.ToSlash(path))
	}
	openInfo, err := file.Stat()
	if err != nil || !os.SameFile(identity, openInfo) {
		return fmt.Errorf("bound source descriptor identity changed: %s", filepath.ToSlash(path))
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	current, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if !bytes.Equal(current, expected) {
		return fmt.Errorf("bound source bytes changed after validation: %s", filepath.ToSlash(path))
	}
	pathInfo, err = root.Lstat(path)
	if err != nil || !os.SameFile(identity, pathInfo) {
		return fmt.Errorf("bound source path identity changed while revalidating: %s", filepath.ToSlash(path))
	}
	return nil
}

func writeMarkerBound(root *os.Root, path string, file *os.File, identity os.FileInfo, marker transactionMarker) error {
	path = filepath.FromSlash(path)
	pathInfo, err := root.Lstat(path)
	if err != nil || identity == nil || !os.SameFile(identity, pathInfo) {
		return fmt.Errorf("transaction marker path identity changed")
	}
	openInfo, err := file.Stat()
	if err != nil || !os.SameFile(identity, openInfo) {
		return fmt.Errorf("transaction marker descriptor identity changed")
	}
	data, _ := json.MarshalIndent(marker, "", "  ")
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if _, err := file.Write(append(data, '\n')); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	pathInfo, err = root.Lstat(path)
	if err != nil || !os.SameFile(identity, pathInfo) {
		return fmt.Errorf("transaction marker path identity changed during update")
	}
	return nil
}

func removeBoundMarker(root *os.Root, path string, file *os.File, identity os.FileInfo) error {
	return removeBoundMarkerWithHook(root, path, file, identity, nil)
}

func removeBoundMarkerWithHook(root *os.Root, path string, file *os.File, identity os.FileInfo, hook PhaseHook) error {
	return removeBoundMarkerWithEvidence(root, path, file, identity, nil, hook)
}

func removeBoundMarkerWithEvidence(root *os.Root, path string, file *os.File, identity os.FileInfo, expected []byte, hook PhaseHook) error {
	path = filepath.FromSlash(path)
	if err := validateBoundMarker(root, path, file, identity); err != nil {
		return err
	}
	if err := invokeHook(hook, "marker-removal-authority-validated"); err != nil {
		return err
	}
	quarantine, err := unpredictableSibling(root, path, "marker-cleanup")
	if err != nil {
		return err
	}
	if err := root.Rename(path, quarantine); err != nil {
		return fmt.Errorf("atomically quarantine transaction marker: %w", err)
	}
	quarantineInfo, err := root.Lstat(quarantine)
	if err != nil || identity == nil || !os.SameFile(identity, quarantineInfo) {
		return fmt.Errorf("transaction marker authority changed during quarantine; retained captured entry at %s", filepath.ToSlash(quarantine))
	}
	openInfo, err := file.Stat()
	if err != nil || !os.SameFile(identity, openInfo) {
		return fmt.Errorf("transaction marker descriptor identity changed after quarantine; retained marker at %s", filepath.ToSlash(quarantine))
	}
	if expected != nil {
		if err := validateBoundRootFile(root, quarantine, file, identity, expected); err != nil {
			return fmt.Errorf("transaction marker bytes changed during quarantine; retained marker at %s: %w", filepath.ToSlash(quarantine), err)
		}
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close quarantined transaction marker at %s: %w", filepath.ToSlash(quarantine), err)
	}
	quarantineInfo, err = root.Lstat(quarantine)
	if err != nil || !os.SameFile(identity, quarantineInfo) {
		return fmt.Errorf("transaction marker quarantine identity changed before removal; retained evidence at %s", filepath.ToSlash(quarantine))
	}
	if err := root.Remove(quarantine); err != nil {
		return fmt.Errorf("remove quarantined transaction marker %s: %w", filepath.ToSlash(quarantine), err)
	}
	return syncRootDirectory(root, filepath.Dir(quarantine))
}

func validateStagedRoot(plan Plan, transactionRoot *os.Root) error {
	for _, action := range plan.Actions {
		if action.Kind == "remove" {
			continue
		}
		stagePath := filepath.Join("stage", filepath.FromSlash(action.Target))
		data, err := transactionRoot.ReadFile(stagePath)
		if err != nil {
			return err
		}
		if string(data) != string(action.Content) {
			return fmt.Errorf("staged content changed for %s", action.Target)
		}
		switch action.Target {
		case state.LockPath:
			lock, err := state.DecodeYAMLMap(data)
			if err != nil {
				return err
			}
			schemaPath, err := state.SchemaForLockVersion(state.AsInt(lock["schema_version"]))
			if err != nil {
				return err
			}
			if err := state.Validate(schemaPath, lock); err != nil {
				return err
			}
		case state.ConfigPath:
			config, err := state.DecodeYAMLMap(data)
			if err != nil {
				return err
			}
			if err := state.Validate(state.ConfigV1Schema, config); err != nil {
				return err
			}
		}
	}
	return nil
}

func commitAction(root string, repositoryRoot, transactionRoot *os.Root, expectedParent pathIdentity, action Action, hook PhaseHook) (committedAction, error) {
	target, err := safeTarget(root, action.Target)
	if err != nil {
		return committedAction{}, err
	}
	if err := invokeHook(hook, "bind-parent:"+action.Target); err != nil {
		return committedAction{}, err
	}
	if err := revalidateParentIdentity(expectedParent); err != nil {
		return committedAction{}, err
	}
	targetRelative := filepath.FromSlash(action.Target)
	parentRelative := filepath.Dir(targetRelative)
	if err := repositoryRoot.MkdirAll(parentRelative, 0o755); err != nil {
		return committedAction{}, err
	}
	parentRoot, err := repositoryRoot.OpenRoot(parentRelative)
	if err != nil {
		return committedAction{}, err
	}
	parentInfo, err := parentRoot.Stat(".")
	if err != nil {
		_ = parentRoot.Close()
		return committedAction{}, err
	}
	if _, err := safeTarget(root, action.Target); err != nil {
		_ = parentRoot.Close()
		return committedAction{}, err
	}
	if err := revalidateParentIdentity(expectedParent); err != nil {
		_ = parentRoot.Close()
		return committedAction{}, err
	}
	currentParentInfo, err := os.Lstat(filepath.Dir(target))
	if err != nil || !os.SameFile(parentInfo, currentParentInfo) {
		_ = parentRoot.Close()
		return committedAction{}, fmt.Errorf("opened parent identity does not match intended path for %s", action.Target)
	}
	record := committedAction{
		action:         action,
		parentRoot:     parentRoot,
		targetName:     filepath.Base(targetRelative),
		backup:         filepath.Join("backup", targetRelative),
		parentIdentity: pathIdentity{path: filepath.Dir(target), info: parentInfo},
	}
	closeUnmutated := func() {
		if !record.mutated && record.parentRoot != nil {
			_ = record.parentRoot.Close()
			record.parentRoot = nil
		}
	}
	targetInfo, statErr := parentRoot.Lstat(record.targetName)
	if statErr == nil {
		if action.Kind == "create" {
			closeUnmutated()
			return record, fmt.Errorf("create precondition changed for %s: target appeared", action.Target)
		}
		if !targetInfo.Mode().IsRegular() || targetInfo.Mode()&os.ModeSymlink != 0 {
			closeUnmutated()
			return record, fmt.Errorf("transaction target is not a regular file: %s", action.Target)
		}
		record.hadTarget = true
		record.originalInfo = targetInfo
		record.backupMode = targetInfo.Mode()
		if err := transactionRoot.MkdirAll(filepath.Dir(record.backup), 0o700); err != nil {
			closeUnmutated()
			return record, err
		}
		record.originalStash, err = unpredictableSibling(parentRoot, record.targetName, "original")
		if err != nil {
			closeUnmutated()
			return record, err
		}
		if err := parentRoot.Rename(record.targetName, record.originalStash); err != nil {
			closeUnmutated()
			return record, err
		}
		record.mutated = true
		movedInfo, err := parentRoot.Lstat(record.originalStash)
		if err != nil || !os.SameFile(targetInfo, movedInfo) {
			return record, fmt.Errorf("source identity changed while backing up %s", action.Target)
		}
		record.backupFile, record.backupInfo, record.backupSHA256, err = copyRootFileBound(parentRoot, record.originalStash, targetInfo, transactionRoot, record.backup, targetInfo.Mode())
		if err != nil {
			return record, err
		}
	} else if !os.IsNotExist(statErr) {
		closeUnmutated()
		return record, statErr
	} else if action.Kind != "create" {
		closeUnmutated()
		return record, fmt.Errorf("source precondition changed for %s: target disappeared", action.Target)
	}
	if action.ExpectedSHA256 != "" {
		if !record.hadTarget {
			closeUnmutated()
			return record, fmt.Errorf("source precondition changed for %s after validation: target disappeared", action.Target)
		}
		actual := record.backupSHA256
		if actual != action.ExpectedSHA256 {
			return record, fmt.Errorf("source precondition changed for %s before commit: expected %s, found %s", action.Target, action.ExpectedSHA256, actual)
		}
	}
	if record.originalStash != "" {
		stashInfo, err := parentRoot.Lstat(record.originalStash)
		if err != nil || record.originalInfo == nil || !os.SameFile(record.originalInfo, stashInfo) {
			return record, fmt.Errorf("original stash identity changed for %s", action.Target)
		}
		stashData, err := parentRoot.ReadFile(record.originalStash)
		if err != nil {
			return record, err
		}
		stashDigest := sha256.Sum256(stashData)
		if hex.EncodeToString(stashDigest[:]) != record.backupSHA256 {
			return record, fmt.Errorf("original stash bytes changed for %s", action.Target)
		}
		if err := parentRoot.Remove(record.originalStash); err != nil {
			return record, err
		}
		record.originalStash = ""
	}
	if action.Kind == "remove" {
		record.committed = true
		return record, nil
	}
	if err := invokeHook(hook, "backed-up:"+action.Target); err != nil {
		closeUnmutated()
		return record, err
	}
	mode := os.FileMode(action.Mode)
	if mode == 0 {
		mode = 0o644
	}
	installed, err := parentRoot.OpenFile(record.targetName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		closeUnmutated()
		return record, err
	}
	record.mutated = true
	record.installed = true
	record.targetInfo, err = installed.Stat()
	if err == nil {
		_, err = installed.Write(action.Content)
	}
	if err == nil {
		err = installed.Sync()
	}
	closeErr := installed.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return record, err
	}
	record.committed = true
	if err := transactionRoot.Remove(filepath.Join("stage", targetRelative)); err != nil {
		return record, err
	}
	return record, nil
}

func rollbackCommitted(result Result, repositoryRoot, transactionRoot *os.Root, transactionInfo os.FileInfo, markerFile *os.File, markerInfo os.FileInfo, transactionPath, markerPath string, marker transactionMarker, committed []committedAction, cause error, hook PhaseHook) (Result, bool) {
	result.Rollback.Attempted = true
	marker.Phase = "rolling-back"
	marker.Error = cause.Error()
	markerErr := writeMarkerBound(repositoryRoot, markerPath, markerFile, markerInfo, marker)
	rollbackErr := rollbackActions(transactionRoot, committed, hook)
	if rollbackErr != nil || markerErr != nil {
		marker.Phase = "recovery-required"
		details := []string{cause.Error()}
		if rollbackErr != nil {
			details = append(details, "rollback: "+rollbackErr.Error())
		}
		if markerErr != nil {
			details = append(details, "marker: "+markerErr.Error())
		}
		marker.Error = strings.Join(details, "; ")
		if markerErr == nil {
			if err := writeMarkerBound(repositoryRoot, markerPath, markerFile, markerInfo, marker); err != nil {
				details = append(details, "marker: "+err.Error())
				marker.Error = strings.Join(details, "; ")
			}
		}
		result.Status = StatusRecoveryRequired
		result.Rollback.Succeeded = false
		result.Rollback.Detail = strings.Join(details, "; ")
		result.Blockers = append(result.Blockers, Blocker{
			Code:         "recovery_required",
			Message:      marker.Error,
			Path:         state.TransactionPath,
			Remediation:  "run check, preserve the transaction directory, then run repair",
			RetryCommand: "repair",
		})
		return result, true
	}
	markerAuthorityErr := validateBoundMarker(repositoryRoot, markerPath, markerFile, markerInfo)
	transactionAuthorityErr := validateBoundTransaction(repositoryRoot, transactionRoot, transactionPath, transactionInfo)
	if markerAuthorityErr != nil || transactionAuthorityErr != nil {
		marker.Phase = "recovery-required"
		details := []string{cause.Error()}
		if markerAuthorityErr != nil {
			details = append(details, "marker authority: "+markerAuthorityErr.Error())
		}
		if transactionAuthorityErr != nil {
			details = append(details, "transaction authority: "+transactionAuthorityErr.Error())
		}
		marker.Error = strings.Join(details, "; ")
		if markerAuthorityErr == nil {
			if err := writeMarkerBound(repositoryRoot, markerPath, markerFile, markerInfo, marker); err != nil {
				details = append(details, "marker: "+err.Error())
				marker.Error = strings.Join(details, "; ")
			}
		}
		result.Status = StatusRecoveryRequired
		result.Rollback.Succeeded = false
		result.Rollback.Detail = marker.Error
		result.Blockers = append(result.Blockers, Blocker{Code: "recovery_required", Message: marker.Error, Path: state.TransactionPath, Remediation: "preserve marker and transaction evidence, then run check before repair", RetryCommand: "repair"})
		return result, true
	}
	if err := invokeHook(hook, "rollback-authority-validated"); err != nil {
		result.Status = StatusRecoveryRequired
		result.Rollback.Succeeded = false
		result.Rollback.Detail = cause.Error() + "; rollback authority hook: " + err.Error()
		result.Blockers = append(result.Blockers, Blocker{Code: "recovery_required", Message: result.Rollback.Detail, Path: state.TransactionPath, Remediation: "preserve marker and transaction evidence, then run check before repair", RetryCommand: "repair"})
		return result, true
	}
	// Remove the verified marker before deleting transaction evidence. If marker authority was
	// lost after validation, retain the transaction directory for explicit recovery.
	if err := removeBoundMarkerWithHook(repositoryRoot, markerPath, markerFile, markerInfo, hook); err != nil {
		result.Status = StatusRecoveryRequired
		result.Rollback.Succeeded = false
		result.Rollback.Detail = cause.Error() + "; marker cleanup: " + err.Error()
		result.Blockers = append(result.Blockers, Blocker{Code: "recovery_required", Message: result.Rollback.Detail, Path: state.TransactionPath, Remediation: "preserve marker and transaction evidence, then run check before repair", RetryCommand: "repair"})
		return result, true
	}
	closeCommittedRoots(&committed)
	cleanupErr := cleanupBoundTransactionWithHook(repositoryRoot, transactionRoot, transactionPath, transactionInfo, hook)
	if cleanupErr != nil {
		result.Status = StatusRecoveryRequired
		result.Rollback.Succeeded = false
		result.Rollback.Detail = cause.Error() + "; cleanup: " + cleanupErr.Error()
		result.Blockers = append(result.Blockers, Blocker{Code: "recovery_required", Message: result.Rollback.Detail, Path: transactionPath, Remediation: "preserve transaction evidence and run check before repair", RetryCommand: "repair"})
		return result, true
	}
	result.Status = StatusRolledBack
	result.Rollback.Succeeded = true
	result.Rollback.Detail = cause.Error()
	result.Blockers = append(result.Blockers, Blocker{Code: "operation_rolled_back", Message: cause.Error(), RetryCommand: result.Command})
	result.StateAfter = result.StateBefore
	return result, false
}

func rollbackActions(transactionRoot *os.Root, committed []committedAction, hook PhaseHook) error {
	errorsList := []string{}
	for index := len(committed) - 1; index >= 0; index-- {
		record := committed[index]
		if !record.mutated {
			continue
		}
		if record.parentIdentity.info == nil {
			detail := fmt.Sprintf("rollback parent identity unavailable for %s", record.action.Target)
			if err := preserveBoundBackup(transactionRoot, record, hook); err != nil {
				detail += "; preserve original backup: " + err.Error()
			}
			errorsList = append(errorsList, detail)
			continue
		}
		if err := revalidateParentIdentity(record.parentIdentity); err != nil {
			detail := fmt.Sprintf("rollback refused changed parent for %s: %v", record.action.Target, err)
			if preserveErr := preserveBoundBackup(transactionRoot, record, hook); preserveErr != nil {
				detail += "; preserve original backup: " + preserveErr.Error()
			}
			errorsList = append(errorsList, detail)
			continue
		}
		boundParentInfo, err := record.parentRoot.Stat(".")
		if err != nil || !os.SameFile(record.parentIdentity.info, boundParentInfo) {
			detail := fmt.Sprintf("rollback bound parent identity unavailable for %s", record.action.Target)
			if preserveErr := preserveBoundBackup(transactionRoot, record, hook); preserveErr != nil {
				detail += "; preserve original backup: " + preserveErr.Error()
			}
			errorsList = append(errorsList, detail)
			continue
		}
		if record.hadTarget {
			if record.backupFile == nil || record.backupInfo == nil {
				errorsList = append(errorsList, fmt.Sprintf("rollback backup unavailable for %s: %v", record.action.Target, err))
				continue
			}
		}
		if record.originalStash != "" {
			stashInfo, err := record.parentRoot.Lstat(record.originalStash)
			if err != nil || !stashInfo.Mode().IsRegular() || record.originalInfo == nil || !os.SameFile(record.originalInfo, stashInfo) {
				errorsList = append(errorsList, fmt.Sprintf("rollback original stash unavailable for %s: %v", record.action.Target, err))
				continue
			}
			if err := rootLinkNoReplace(record.parentRoot, record.originalStash, record.targetName); err != nil {
				errorsList = append(errorsList, fmt.Sprintf("restore original stash for %s: %v", record.action.Target, err))
				continue
			}
			if err := record.parentRoot.Remove(record.originalStash); err != nil {
				errorsList = append(errorsList, fmt.Sprintf("remove restored stash for %s: %v", record.action.Target, err))
			}
			continue
		}
		if record.installed {
			currentInfo, err := record.parentRoot.Lstat(record.targetName)
			if err != nil {
				errorsList = append(errorsList, fmt.Sprintf("rollback target identity unavailable for %s: %v", record.action.Target, err))
				continue
			}
			if record.targetInfo == nil || !os.SameFile(record.targetInfo, currentInfo) {
				errorsList = append(errorsList, fmt.Sprintf("rollback target identity changed concurrently: %s", record.action.Target))
				continue
			}
			quarantine, err := unpredictableSibling(record.parentRoot, record.targetName, "rollback-current")
			if err != nil {
				errorsList = append(errorsList, err.Error())
				continue
			}
			if err := record.parentRoot.Rename(record.targetName, quarantine); err != nil {
				errorsList = append(errorsList, fmt.Sprintf("quarantine rollback target %s: %v", record.action.Target, err))
				continue
			}
			quarantineInfo, err := record.parentRoot.Lstat(quarantine)
			if err != nil || !os.SameFile(record.targetInfo, quarantineInfo) {
				detail := fmt.Sprintf("rollback quarantined target identity mismatch: %s", record.action.Target)
				if err != nil {
					detail += ": " + err.Error()
				}
				errorsList = append(errorsList, detail)
				continue
			}
			matches, matchErr := regularRootFileEquals(record.parentRoot, quarantine, record.action.Content)
			if record.committed && (matchErr != nil || !matches) {
				restoreErr := rootLinkNoReplace(record.parentRoot, quarantine, record.targetName)
				detail := fmt.Sprintf("rollback target changed concurrently: %s", record.action.Target)
				if matchErr != nil {
					detail += ": " + matchErr.Error()
				}
				if restoreErr == nil {
					_ = record.parentRoot.Remove(quarantine)
				} else {
					detail += "; preserve quarantine: " + restoreErr.Error()
				}
				errorsList = append(errorsList, detail)
				continue
			}
			if err := invokeHook(hook, "rollback-quarantined:"+record.action.Target); err != nil {
				errorsList = append(errorsList, err.Error())
				continue
			}
			if err := revalidateParentIdentity(record.parentIdentity); err != nil {
				errorsList = append(errorsList, fmt.Sprintf("rollback parent changed after quarantine for %s: %v", record.action.Target, err))
				continue
			}
			postHookInfo, postHookInfoErr := record.parentRoot.Lstat(quarantine)
			postHookMatches, postHookMatchErr := regularRootFileEquals(record.parentRoot, quarantine, record.action.Content)
			if postHookInfoErr != nil || record.targetInfo == nil || !os.SameFile(record.targetInfo, postHookInfo) || (record.committed && (postHookMatchErr != nil || !postHookMatches)) {
				restoreErr := rootLinkNoReplace(record.parentRoot, quarantine, record.targetName)
				detail := fmt.Sprintf("rollback quarantine changed after hook: %s", record.action.Target)
				if postHookInfoErr != nil {
					detail += ": " + postHookInfoErr.Error()
				} else if postHookMatchErr != nil {
					detail += ": " + postHookMatchErr.Error()
				}
				if restoreErr == nil {
					_ = record.parentRoot.Remove(quarantine)
				} else {
					detail += "; preserve quarantine: " + restoreErr.Error()
				}
				errorsList = append(errorsList, detail)
				continue
			}
			if record.hadTarget {
				backupData, backupPathChanged, err := readVerifiedBackup(transactionRoot, record)
				if err == nil {
					err = writeRootNoReplace(record.parentRoot, record.targetName, backupData, record.backupMode)
				}
				if err != nil {
					errorsList = append(errorsList, fmt.Sprintf("restore rollback backup for %s without replacement: %v", record.action.Target, err))
					continue
				}
				if backupPathChanged {
					errorsList = append(errorsList, fmt.Sprintf("rollback restored bound backup for %s but preserved a substituted backup path", record.action.Target))
					continue
				}
			}
			if err := record.parentRoot.Remove(quarantine); err != nil {
				errorsList = append(errorsList, fmt.Sprintf("remove rollback quarantine for %s: %v", record.action.Target, err))
			}
			continue
		}
		if record.hadTarget {
			backupData, backupPathChanged, err := readVerifiedBackup(transactionRoot, record)
			if err == nil {
				err = writeRootNoReplace(record.parentRoot, record.targetName, backupData, record.backupMode)
			}
			if err != nil {
				errorsList = append(errorsList, fmt.Sprintf("restore rollback backup for %s without replacement: %v", record.action.Target, err))
			} else if backupPathChanged {
				errorsList = append(errorsList, fmt.Sprintf("rollback restored bound backup for %s but preserved a substituted backup path", record.action.Target))
			}
		}
	}
	if len(errorsList) > 0 {
		return errors.New(strings.Join(errorsList, "; "))
	}
	return nil
}

func regularRootFileEquals(root *os.Root, path string, expected []byte) (bool, error) {
	info, err := root.Lstat(path)
	if err != nil {
		return false, err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return false, fmt.Errorf("quarantined target is not a regular file")
	}
	current, err := root.ReadFile(path)
	if err != nil {
		return false, err
	}
	return bytes.Equal(current, expected), nil
}

func copyRootFileBound(sourceRoot *os.Root, source string, expectedSource os.FileInfo, targetRoot *os.Root, target string, mode os.FileMode) (*os.File, os.FileInfo, string, error) {
	sourceFile, err := sourceRoot.Open(source)
	if err != nil {
		return nil, nil, "", err
	}
	defer sourceFile.Close()
	info, err := sourceFile.Stat()
	if err != nil {
		return nil, nil, "", err
	}
	if !info.Mode().IsRegular() || expectedSource == nil || !os.SameFile(expectedSource, info) {
		return nil, nil, "", fmt.Errorf("source identity changed while copying: %s", source)
	}
	sourceData, err := io.ReadAll(sourceFile)
	if err != nil {
		return nil, nil, "", err
	}
	targetFile, err := targetRoot.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_EXCL, mode.Perm())
	if err != nil {
		return nil, nil, "", err
	}
	removeTarget := true
	defer func() {
		if removeTarget {
			_ = targetFile.Close()
			_ = targetRoot.Remove(target)
		}
	}()
	if _, err := targetFile.Write(sourceData); err != nil {
		return nil, nil, "", err
	}
	if err := targetFile.Sync(); err != nil {
		return nil, nil, "", err
	}
	targetInfo, err := targetFile.Stat()
	if err != nil {
		return nil, nil, "", err
	}
	pathInfo, err := targetRoot.Lstat(target)
	if err != nil || !os.SameFile(targetInfo, pathInfo) {
		removeTarget = false
		_ = targetFile.Close()
		return nil, nil, "", fmt.Errorf("backup path identity changed during creation")
	}
	digest := sha256.Sum256(sourceData)
	removeTarget = false
	return targetFile, targetInfo, hex.EncodeToString(digest[:]), nil
}

func rootLinkNoReplace(root *os.Root, source, target string) error {
	return root.Link(source, target)
}

func readVerifiedBackup(transactionRoot *os.Root, record committedAction) ([]byte, bool, error) {
	pathInfo, err := transactionRoot.Lstat(record.backup)
	pathChanged := err != nil || record.backupInfo == nil || !os.SameFile(record.backupInfo, pathInfo)
	if record.backupFile == nil || record.backupInfo == nil {
		return nil, pathChanged, fmt.Errorf("backup descriptor unavailable")
	}
	openInfo, err := record.backupFile.Stat()
	if err != nil || !os.SameFile(record.backupInfo, openInfo) {
		return nil, pathChanged, fmt.Errorf("backup descriptor identity changed")
	}
	if _, err := record.backupFile.Seek(0, io.SeekStart); err != nil {
		return nil, pathChanged, err
	}
	data, err := io.ReadAll(record.backupFile)
	if err != nil {
		return nil, pathChanged, err
	}
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != record.backupSHA256 {
		return nil, pathChanged, fmt.Errorf("backup bytes changed")
	}
	return data, pathChanged, nil
}

func preserveBoundBackup(transactionRoot *os.Root, record committedAction, hook PhaseHook) error {
	if !record.hadTarget {
		return nil
	}
	data, pathChanged, err := readVerifiedBackup(transactionRoot, record)
	if err != nil {
		return err
	}
	if !pathChanged {
		return nil
	}
	target := filepath.FromSlash(record.action.Target)
	recoveryDirectory := filepath.Join("recovery-original", filepath.Dir(target))
	if err := transactionRoot.MkdirAll(recoveryDirectory, 0o700); err != nil {
		return err
	}
	var lastErr error
	for attempt := 0; attempt < 8; attempt++ {
		nonce := make([]byte, 16)
		if _, err := rand.Read(nonce); err != nil {
			return fmt.Errorf("generate recovery path: %w", err)
		}
		recoveryPath := filepath.Join(recoveryDirectory, filepath.Base(target)+".original-"+hex.EncodeToString(nonce))
		if err := writeRootNoReplaceWithHook(transactionRoot, recoveryPath, data, record.backupMode, hook); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	return fmt.Errorf("could not preserve verified original bytes at an unpredictable recovery path for %s after bounded retries: %w", record.action.Target, lastErr)
}

func writeRootNoReplace(root *os.Root, target string, data []byte, mode os.FileMode) error {
	return writeRootNoReplaceWithHook(root, target, data, mode, nil)
}

func writeRootNoReplaceWithHook(root *os.Root, target string, data []byte, mode os.FileMode, hook PhaseHook) error {
	file, err := root.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_EXCL, mode.Perm())
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			_ = file.Close()
		}
	}()
	openedInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	written, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if !bytes.Equal(written, data) {
		return fmt.Errorf("%w: recovery bytes changed before close", errRecoveryPathAuthority)
	}
	pathInfo, err := root.Lstat(target)
	if err != nil || !os.SameFile(openedInfo, pathInfo) {
		return fmt.Errorf("%w during creation", errRecoveryPathAuthority)
	}
	if err := file.Close(); err != nil {
		return err
	}
	closed = true
	pathInfo, err = root.Lstat(target)
	if err != nil || !os.SameFile(openedInfo, pathInfo) {
		return fmt.Errorf("%w after creation", errRecoveryPathAuthority)
	}
	if err := syncRootDirectory(root, filepath.Dir(target)); err != nil {
		return err
	}
	if err := invokeHook(hook, "recovery-copy-after-directory-sync:"+filepath.ToSlash(target)); err != nil {
		return err
	}
	reopened, err := root.Open(target)
	if err != nil {
		return fmt.Errorf("%w while reopening: %v", errRecoveryPathAuthority, err)
	}
	reopenedInfo, err := reopened.Stat()
	if err != nil || !os.SameFile(openedInfo, reopenedInfo) {
		_ = reopened.Close()
		return fmt.Errorf("%w after directory sync", errRecoveryPathAuthority)
	}
	reopenedData, err := io.ReadAll(reopened)
	closeErr := reopened.Close()
	if err != nil {
		return fmt.Errorf("%w while rereading durable recovery bytes: %v", errRecoveryPathAuthority, err)
	}
	if closeErr != nil {
		return fmt.Errorf("%w while closing durable recovery verification: %v", errRecoveryPathAuthority, closeErr)
	}
	if !bytes.Equal(reopenedData, data) {
		return fmt.Errorf("%w: recovery bytes changed after directory sync", errRecoveryPathAuthority)
	}
	pathInfo, err = root.Lstat(target)
	if err != nil || !os.SameFile(openedInfo, pathInfo) {
		return fmt.Errorf("%w after durable verification", errRecoveryPathAuthority)
	}
	return nil
}

func syncRootDirectory(root *os.Root, directory string) error {
	directoryFile, err := root.Open(directory)
	if err != nil {
		return err
	}
	defer directoryFile.Close()
	if err := directoryFile.Sync(); err != nil && runtime.GOOS != "windows" {
		return err
	}
	return nil
}

func rollbackBeforeCommitRoot(result Result, repositoryRoot *os.Root, markerFile *os.File, markerInfo os.FileInfo, markerPath, transactionPath string, marker transactionMarker, cause error, rootCreated bool, root string, transactionRoot *os.Root, transactionInfo os.FileInfo) (Result, error) {
	marker.Phase = "failed"
	marker.Error = cause.Error()
	markerErr := writeMarkerBound(repositoryRoot, markerPath, markerFile, markerInfo, marker)
	var cleanupErr error
	var removeMarkerErr error
	if markerErr == nil && transactionRoot != nil {
		markerErr = validateBoundMarker(repositoryRoot, markerPath, markerFile, markerInfo)
		cleanupErr = validateBoundTransaction(repositoryRoot, transactionRoot, transactionPath, transactionInfo)
	}
	if markerErr == nil && cleanupErr == nil {
		removeMarkerErr = removeBoundMarker(repositoryRoot, markerPath, markerFile, markerInfo)
	}
	if markerErr == nil && cleanupErr == nil && removeMarkerErr == nil && transactionRoot != nil {
		cleanupErr = cleanupBoundTransaction(repositoryRoot, transactionRoot, transactionPath, transactionInfo)
	}
	if rootCreated {
		_ = os.Remove(root)
	}
	if markerErr != nil || cleanupErr != nil || removeMarkerErr != nil {
		details := []string{cause.Error()}
		for _, failure := range []error{markerErr, cleanupErr, removeMarkerErr} {
			if failure != nil {
				details = append(details, failure.Error())
			}
		}
		result.Status = StatusRecoveryRequired
		result.Rollback = Rollback{Attempted: true, Succeeded: false, Detail: strings.Join(details, "; ")}
		result.Blockers = append(result.Blockers, Blocker{Code: "recovery_required", Message: strings.Join(details, "; "), Path: state.TransactionPath, Remediation: "preserve concurrent transaction evidence and run check before repair", RetryCommand: "repair"})
		return result, nil
	}
	result.Status = StatusRolledBack
	result.Rollback = Rollback{Attempted: true, Succeeded: true, Detail: cause.Error()}
	result.Blockers = append(result.Blockers, Blocker{Code: "operation_rolled_back", Message: cause.Error(), RetryCommand: result.Command})
	return result, nil
}

func cleanupBoundTransaction(repositoryRoot, transactionRoot *os.Root, transactionPath string, transactionInfo os.FileInfo) error {
	return cleanupBoundTransactionWithHook(repositoryRoot, transactionRoot, transactionPath, transactionInfo, nil)
}

func cleanupBoundTransactionWithHook(repositoryRoot, transactionRoot *os.Root, transactionPath string, transactionInfo os.FileInfo, hook PhaseHook) error {
	return cleanupBoundTransactionWithEvidence(repositoryRoot, transactionRoot, transactionPath, transactionInfo, nil, hook)
}

func cleanupBoundTransactionWithEvidence(repositoryRoot, transactionRoot *os.Root, transactionPath string, transactionInfo os.FileInfo, evidence *boundFileEvidence, hook PhaseHook) error {
	if err := validateBoundTransaction(repositoryRoot, transactionRoot, transactionPath, transactionInfo); err != nil {
		return err
	}
	if err := invokeHook(hook, "transaction-cleanup-authority-validated"); err != nil {
		return err
	}
	transactionPath = filepath.FromSlash(transactionPath)
	quarantine, err := unpredictableSibling(repositoryRoot, transactionPath, "transaction-cleanup")
	if err != nil {
		return err
	}
	if err := repositoryRoot.Rename(transactionPath, quarantine); err != nil {
		return fmt.Errorf("atomically quarantine transaction directory: %w", err)
	}
	if err := validateBoundTransaction(repositoryRoot, transactionRoot, quarantine, transactionInfo); err != nil {
		return fmt.Errorf("transaction authority changed during quarantine; retained captured entry at %s: %w", filepath.ToSlash(quarantine), err)
	}
	if evidence != nil {
		if err := validateBoundRootFile(transactionRoot, evidence.path, evidence.file, evidence.info, evidence.data); err != nil {
			return fmt.Errorf("transaction evidence changed during quarantine; retained at %s: %w", filepath.ToSlash(quarantine), err)
		}
	}
	entries, err := fs.ReadDir(transactionRoot.FS(), ".")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := validateBoundTransaction(repositoryRoot, transactionRoot, quarantine, transactionInfo); err != nil {
			return err
		}
		if err := transactionRoot.RemoveAll(entry.Name()); err != nil {
			return err
		}
	}
	if err := validateBoundTransaction(repositoryRoot, transactionRoot, quarantine, transactionInfo); err != nil {
		return err
	}
	if err := transactionRoot.Close(); err != nil {
		return err
	}
	pathInfo, err := repositoryRoot.Lstat(quarantine)
	if os.IsNotExist(err) {
		return fmt.Errorf("transaction quarantine disappeared before removal: %s", filepath.ToSlash(quarantine))
	}
	if err != nil {
		return err
	}
	if transactionInfo == nil || !os.SameFile(transactionInfo, pathInfo) {
		return fmt.Errorf("transaction quarantine identity changed during cleanup; retained at %s", filepath.ToSlash(quarantine))
	}
	if err := repositoryRoot.Remove(quarantine); err != nil {
		return err
	}
	if err := syncRootDirectory(repositoryRoot, filepath.Dir(quarantine)); err != nil {
		return err
	}
	parent := filepath.Dir(transactionPath)
	if _, err := safeTarget(repositoryRoot.Name(), parent); err == nil {
		_ = repositoryRoot.Remove(parent)
	}
	return nil
}

func unpredictableSibling(root *os.Root, path, purpose string) (string, error) {
	directory := filepath.Dir(path)
	base := filepath.Base(path)
	for attempt := 0; attempt < 8; attempt++ {
		nonce := make([]byte, 16)
		if _, err := rand.Read(nonce); err != nil {
			return "", fmt.Errorf("generate %s quarantine path: %w", purpose, err)
		}
		candidate := filepath.Join(directory, "."+base+"."+purpose+"-"+hex.EncodeToString(nonce))
		if _, err := root.Lstat(candidate); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return "", err
		}
		return candidate, nil
	}
	return "", fmt.Errorf("could not reserve an unpredictable %s quarantine path", purpose)
}

func validateBoundMarker(root *os.Root, path string, file *os.File, identity os.FileInfo) error {
	pathInfo, err := root.Lstat(filepath.FromSlash(path))
	if err != nil || identity == nil || !os.SameFile(identity, pathInfo) {
		return fmt.Errorf("transaction marker path identity changed")
	}
	openInfo, err := file.Stat()
	if err != nil || !os.SameFile(identity, openInfo) {
		return fmt.Errorf("transaction marker descriptor identity changed")
	}
	return nil
}

func validateBoundTransaction(repositoryRoot, transactionRoot *os.Root, transactionPath string, transactionInfo os.FileInfo) error {
	if transactionInfo == nil {
		return fmt.Errorf("transaction root identity unavailable")
	}
	boundInfo, err := transactionRoot.Stat(".")
	if err != nil || !os.SameFile(transactionInfo, boundInfo) {
		return fmt.Errorf("transaction root descriptor identity changed")
	}
	pathInfo, err := repositoryRoot.Lstat(filepath.FromSlash(transactionPath))
	if err != nil || !os.SameFile(transactionInfo, pathInfo) {
		return fmt.Errorf("transaction path identity changed before cleanup")
	}
	return nil
}

func closeCommittedRoots(committed *[]committedAction) {
	for index := range *committed {
		record := &(*committed)[index]
		if record.backupFile != nil {
			_ = record.backupFile.Close()
			record.backupFile = nil
		}
		if record.parentRoot != nil {
			_ = record.parentRoot.Close()
			record.parentRoot = nil
		}
	}
}

func safeTarget(root, relative string) (string, error) {
	cleaned := filepath.Clean(filepath.FromSlash(relative))
	if cleaned == "." || filepath.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("target escapes root: %s", relative)
	}
	target := filepath.Join(root, cleaned)
	rootWithSep := root + string(filepath.Separator)
	if target != root && !strings.HasPrefix(target, rootWithSep) {
		return "", fmt.Errorf("target escapes root: %s", relative)
	}
	current := root
	parts := strings.Split(cleaned, string(filepath.Separator))
	for _, part := range parts {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("target traverses symbolic link or reparse point: %s", relative)
		}
		if runtime.GOOS == "windows" && info.Mode()&os.ModeIrregular != 0 {
			return "", fmt.Errorf("target traverses irregular Windows filesystem entry: %s", relative)
		}
		resolved, err := filepath.EvalSymlinks(current)
		if err != nil {
			return "", err
		}
		resolved, err = filepath.Abs(resolved)
		if err != nil {
			return "", err
		}
		currentAbsolute, err := filepath.Abs(current)
		if err != nil {
			return "", err
		}
		if !sameCanonicalPath(resolved, currentAbsolute) {
			return "", fmt.Errorf("target traverses symbolic link or Windows reparse point: %s", relative)
		}
	}
	return target, nil
}

func sameCanonicalPath(left, right string) bool {
	left = filepath.Clean(left)
	right = filepath.Clean(right)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(left, right)
	}
	return left == right
}

func validSHA256(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func validTransactionID(value string) bool {
	if len(value) != 32 {
		return false
	}
	for index := 0; index < len(value); index++ {
		if (value[index] < '0' || value[index] > '9') && (value[index] < 'a' || value[index] > 'f') {
			return false
		}
	}
	return true
}

func nearestExistingParent(root, relative string) (pathIdentity, error) {
	target, err := safeTarget(root, relative)
	if err != nil {
		return pathIdentity{}, err
	}
	current := filepath.Dir(target)
	for {
		info, err := os.Lstat(current)
		if err == nil {
			return pathIdentity{path: current, info: info}, nil
		}
		if !os.IsNotExist(err) {
			return pathIdentity{}, err
		}
		if current == root {
			return pathIdentity{}, fmt.Errorf("canonical root disappeared")
		}
		parent := filepath.Dir(current)
		if parent == current {
			return pathIdentity{}, fmt.Errorf("no existing parent for %s", relative)
		}
		current = parent
	}
}

func revalidateParentIdentity(identity pathIdentity) error {
	current, err := os.Lstat(identity.path)
	if err != nil {
		return fmt.Errorf("parent identity unavailable for %s: %w", identity.path, err)
	}
	if !os.SameFile(identity.info, current) {
		return fmt.Errorf("parent identity changed for %s", identity.path)
	}
	return nil
}

func invokeHook(hook PhaseHook, phase string) error {
	if hook == nil {
		return nil
	}
	return hook(phase)
}

func sortedActionTargets(actions []Action) []string {
	values := make([]string, len(actions))
	for index, action := range actions {
		values[index] = action.Target
	}
	sort.Strings(values)
	return values
}
