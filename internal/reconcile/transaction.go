package reconcile

import (
	"encoding/json"
	"errors"
	"fmt"
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
	action    Action
	target    string
	backup    string
	hadTarget bool
}

type pathIdentity struct {
	path string
	info os.FileInfo
}

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
	identities := map[string]pathIdentity{}
	for _, action := range plan.Actions {
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
	transactionRoot := filepath.Join(canonicalRoot, ".codeheart", "local", "kit-transactions", plan.ID)
	markerPath := filepath.Join(canonicalRoot, filepath.FromSlash(state.TransactionPath))
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
	if err := acquireMarker(markerPath, marker); err != nil {
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
			_ = os.Remove(markerPath)
		}
	}()
	if err := invokeHook(options.Hook, "marker"); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}
	if err := os.MkdirAll(filepath.Join(transactionRoot, "stage"), 0o700); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}
	if err := os.MkdirAll(filepath.Join(transactionRoot, "backup"), 0o700); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}
	planData, _ := json.MarshalIndent(plan, "", "  ")
	if err := os.WriteFile(filepath.Join(transactionRoot, "plan.json"), append(planData, '\n'), 0o600); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}
	for _, action := range plan.Actions {
		if action.Kind == "remove" {
			continue
		}
		stagePath := filepath.Join(transactionRoot, "stage", filepath.FromSlash(action.Target))
		if err := os.MkdirAll(filepath.Dir(stagePath), 0o700); err != nil {
			return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
		}
		mode := os.FileMode(action.Mode)
		if mode == 0 {
			mode = 0o644
		}
		if err := os.WriteFile(stagePath, action.Content, mode); err != nil {
			return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
		}
	}
	marker.Phase = "staged"
	if err := writeMarker(markerPath, marker); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}
	if err := invokeHook(options.Hook, "staged"); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}
	if err := validateStaged(plan, transactionRoot); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}
	result.Validations = append(result.Validations, Validation{Name: "staged-state", Status: "passed"})
	marker.Phase = "validated"
	if err := writeMarker(markerPath, marker); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}
	if err := invokeHook(options.Hook, "validated"); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}

	marker.Phase = "committing"
	if err := writeMarker(markerPath, marker); err != nil {
		return rollbackBeforeCommit(result, markerPath, transactionRoot, marker, err, rootCreated, plan.Root)
	}
	committed := []committedAction{}
	for _, action := range plan.Actions {
		if err := revalidateParentIdentity(identities[action.Target]); err != nil {
			result, keepMarker := rollbackCommitted(result, canonicalRoot, transactionRoot, markerPath, marker, committed, err)
			cleanupMarker = !keepMarker
			return result, nil
		}
		if err := invokeHook(options.Hook, "commit:"+action.Target); err != nil {
			result, keepMarker := rollbackCommitted(result, canonicalRoot, transactionRoot, markerPath, marker, committed, err)
			cleanupMarker = !keepMarker
			return result, nil
		}
		record, err := commitAction(canonicalRoot, transactionRoot, action)
		if err != nil {
			result, keepMarker := rollbackCommitted(result, canonicalRoot, transactionRoot, markerPath, marker, committed, err)
			cleanupMarker = !keepMarker
			return result, nil
		}
		committed = append(committed, record)
	}
	marker.Phase = "post-check"
	if err := writeMarker(markerPath, marker); err != nil {
		result, keepMarker := rollbackCommitted(result, canonicalRoot, transactionRoot, markerPath, marker, committed, err)
		cleanupMarker = !keepMarker
		return result, nil
	}
	if err := invokeHook(options.Hook, "post-check"); err != nil {
		result, keepMarker := rollbackCommitted(result, canonicalRoot, transactionRoot, markerPath, marker, committed, err)
		cleanupMarker = !keepMarker
		return result, nil
	}
	observed, err := state.InspectIgnoringTransaction(canonicalRoot)
	if err != nil || !expectedClassification(plan.ExpectedAfter, observed.Classification) {
		if err == nil {
			err = fmt.Errorf("post-check classified state as %s", observed.Classification)
		}
		result, keepMarker := rollbackCommitted(result, canonicalRoot, transactionRoot, markerPath, marker, committed, err)
		cleanupMarker = !keepMarker
		return result, nil
	}
	result.Validations = append(result.Validations, Validation{Name: "post-check", Status: "passed"})
	result.StateAfter = string(observed.Classification)
	result.Status = StatusSucceeded
	if err := os.RemoveAll(transactionRoot); err != nil {
		result.Validations = append(result.Validations, Validation{Name: "transaction-cleanup", Status: "warning", Detail: err.Error()})
	} else {
		result.Validations = append(result.Validations, Validation{Name: "transaction-cleanup", Status: "passed"})
	}
	cleanupTransactionParents(canonicalRoot)
	if err := os.Remove(markerPath); err != nil && !os.IsNotExist(err) {
		result.Validations = append(result.Validations, Validation{Name: "marker-cleanup", Status: "warning", Detail: err.Error()})
	}
	cleanupMarker = false
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

func acquireMarker(path string, marker transactionMarker) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if _, err := os.Lstat(path); err == nil {
		cleared, recoveryErr := recoverStaleMarker(path, marker.TargetRoot, false)
		if recoveryErr != nil {
			return recoveryErr
		}
		if !cleared {
			return fmt.Errorf("active transaction marker already exists")
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	data, _ := json.MarshalIndent(marker, "", "  ")
	if _, err := file.Write(append(data, '\n')); err != nil {
		_ = os.Remove(path)
		return err
	}
	if err := file.Sync(); err != nil {
		_ = os.Remove(path)
		return err
	}
	return nil
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
	markerPath := filepath.Join(canonical, filepath.FromSlash(state.TransactionPath))
	return recoverStaleMarker(markerPath, canonical, dryRun)
}

func recoverStaleMarker(markerPath, canonicalRoot string, dryRun bool) (bool, error) {
	data, err := os.ReadFile(markerPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var existing transactionMarker
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&existing); err != nil {
		return false, fmt.Errorf("existing transaction marker is invalid: %w", err)
	}
	if existing.SchemaVersion != 1 || existing.TransactionID == "" || existing.Command == "" || existing.PID <= 0 {
		return false, fmt.Errorf("existing transaction marker identity is incomplete")
	}
	if !sameCanonicalPath(existing.TargetRoot, canonicalRoot) {
		return false, fmt.Errorf("existing transaction marker targets a different root")
	}
	expectedRecovery := filepath.ToSlash(filepath.Join(".codeheart", "local", "kit-transactions", existing.TransactionID))
	if existing.RecoveryPath != expectedRecovery {
		return false, fmt.Errorf("existing transaction recovery path does not match its identity")
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
	recoveryRoot := filepath.Join(canonicalRoot, filepath.FromSlash(existing.RecoveryPath))
	planRequired := existing.Phase == "staged" || existing.Phase == "validated"
	if planData, err := os.ReadFile(filepath.Join(recoveryRoot, "plan.json")); err == nil {
		var plan Plan
		if err := json.Unmarshal(planData, &plan); err != nil || plan.ID != existing.TransactionID || !sameCanonicalPath(plan.Root, canonicalRoot) || plan.Command != existing.Command {
			return false, fmt.Errorf("stale transaction plan identity does not match its marker")
		}
	} else if os.IsNotExist(err) && planRequired {
		return false, fmt.Errorf("stale transaction plan is missing")
	} else if !os.IsNotExist(err) {
		return false, err
	}
	if dryRun {
		return true, nil
	}
	if err := os.RemoveAll(recoveryRoot); err != nil {
		return false, err
	}
	if err := os.Remove(markerPath); err != nil {
		return false, err
	}
	cleanupTransactionParents(canonicalRoot)
	return true, nil
}

func writeMarker(path string, marker transactionMarker) error {
	data, _ := json.MarshalIndent(marker, "", "  ")
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func validateStaged(plan Plan, transactionRoot string) error {
	for _, action := range plan.Actions {
		if action.Kind == "remove" {
			continue
		}
		stagePath := filepath.Join(transactionRoot, "stage", filepath.FromSlash(action.Target))
		data, err := os.ReadFile(stagePath)
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

func commitAction(root, transactionRoot string, action Action) (committedAction, error) {
	target, err := safeTarget(root, action.Target)
	if err != nil {
		return committedAction{}, err
	}
	backup := filepath.Join(transactionRoot, "backup", filepath.FromSlash(action.Target))
	record := committedAction{action: action, target: target, backup: backup}
	if _, err := os.Lstat(target); err == nil {
		record.hadTarget = true
		if err := os.MkdirAll(filepath.Dir(backup), 0o700); err != nil {
			return record, err
		}
		if err := os.Rename(target, backup); err != nil {
			return record, err
		}
	} else if !os.IsNotExist(err) {
		return record, err
	}
	if action.Kind == "remove" {
		return record, nil
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return record, err
	}
	if _, err := safeTarget(root, action.Target); err != nil {
		return record, err
	}
	stage := filepath.Join(transactionRoot, "stage", filepath.FromSlash(action.Target))
	if err := os.Rename(stage, target); err != nil {
		return record, err
	}
	return record, nil
}

func rollbackCommitted(result Result, root, transactionRoot, markerPath string, marker transactionMarker, committed []committedAction, cause error) (Result, bool) {
	result.Rollback.Attempted = true
	marker.Phase = "rolling-back"
	marker.Error = cause.Error()
	_ = writeMarker(markerPath, marker)
	err := rollbackActions(committed)
	if err != nil {
		marker.Phase = "recovery-required"
		marker.Error = cause.Error() + "; rollback: " + err.Error()
		_ = writeMarker(markerPath, marker)
		result.Status = StatusRecoveryRequired
		result.Rollback.Succeeded = false
		result.Rollback.Detail = err.Error()
		result.Blockers = append(result.Blockers, Blocker{
			Code:         "recovery_required",
			Message:      marker.Error,
			Path:         state.TransactionPath,
			Remediation:  "run check, preserve the transaction directory, then run repair",
			RetryCommand: "repair",
		})
		return result, true
	}
	result.Status = StatusRolledBack
	result.Rollback.Succeeded = true
	result.Rollback.Detail = cause.Error()
	result.Blockers = append(result.Blockers, Blocker{Code: "operation_rolled_back", Message: cause.Error(), RetryCommand: result.Command})
	_ = os.RemoveAll(transactionRoot)
	cleanupTransactionParents(root)
	_ = os.Remove(markerPath)
	result.StateAfter = result.StateBefore
	return result, false
}

func rollbackActions(committed []committedAction) error {
	errorsList := []string{}
	for index := len(committed) - 1; index >= 0; index-- {
		record := committed[index]
		if record.action.Kind != "remove" {
			if err := os.RemoveAll(record.target); err != nil {
				errorsList = append(errorsList, err.Error())
			}
		}
		if record.hadTarget {
			if err := os.MkdirAll(filepath.Dir(record.target), 0o755); err != nil {
				errorsList = append(errorsList, err.Error())
				continue
			}
			if err := os.Rename(record.backup, record.target); err != nil {
				errorsList = append(errorsList, err.Error())
			}
		}
	}
	if len(errorsList) > 0 {
		return errors.New(strings.Join(errorsList, "; "))
	}
	return nil
}

func rollbackBeforeCommit(result Result, markerPath, transactionRoot string, marker transactionMarker, cause error, rootCreated bool, root string) (Result, error) {
	marker.Phase = "failed"
	marker.Error = cause.Error()
	_ = writeMarker(markerPath, marker)
	_ = os.RemoveAll(transactionRoot)
	_ = os.Remove(markerPath)
	cleanupTransactionParents(root)
	if rootCreated {
		_ = os.Remove(root)
	}
	result.Status = StatusRolledBack
	result.Rollback = Rollback{Attempted: true, Succeeded: true, Detail: cause.Error()}
	result.Blockers = append(result.Blockers, Blocker{Code: "operation_rolled_back", Message: cause.Error(), RetryCommand: result.Command})
	return result, nil
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

func cleanupTransactionParents(root string) {
	for _, path := range []string{
		filepath.Join(root, ".codeheart", "local", "kit-transactions"),
		filepath.Join(root, ".codeheart", "local"),
		filepath.Join(root, ".codeheart"),
	} {
		_ = os.Remove(path)
	}
}

func sortedActionTargets(actions []Action) []string {
	values := make([]string, len(actions))
	for index, action := range actions {
		values[index] = action.Target
	}
	sort.Strings(values)
	return values
}
