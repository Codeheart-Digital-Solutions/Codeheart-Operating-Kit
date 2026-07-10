package commands

import (
	"fmt"
	"io"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

func RunRepair(args []string, stdout io.Writer, stderr io.Writer) int {
	_, bools, positionals, err := parseValueArgs(args, map[string]bool{
		"--dry-run": false,
		"--json":    false,
	})
	if err != nil {
		return writeArgError(stderr, "repair", err)
	}
	path := "."
	if len(positionals) > 0 {
		path = positionals[0]
	}
	if len(positionals) > 1 {
		return writeArgError(stderr, "repair", fmt.Errorf("unexpected argument %q", positionals[1]))
	}
	payload, operation, err := repairOperation(path, time.Now(), bools["--dry-run"])
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit repair: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, payload); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit repair: error: %v\n", err)
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

func Repair(path string, now time.Time) (map[string]any, error) {
	payload, _, err := repairOperation(path, now, false)
	return payload, err
}

func repairOperation(path string, now time.Time, dryRun bool) (map[string]any, reconcile.Result, error) {
	root := expandPath(path)
	observed, err := state.Inspect(root)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	if observed.Classification == state.StateTransactionInProgress || observed.Classification == state.StateRecoveryRequired {
		recoverable, recoveryErr := reconcile.RecoverStaleTransaction(root, dryRun)
		if recoveryErr != nil {
			result := blockedLifecycle("repair", observed, "transaction_recovery_blocked", recoveryErr.Error(), "wait for an active process or preserve recovery evidence for manual diagnosis", "check")
			return resultPayload(result), result, nil
		}
		if recoverable && dryRun {
			result := reconcile.NewResult("repair")
			result.Status = reconcile.StatusPlanned
			result.DryRun = true
			result.StateBefore = string(observed.Classification)
			result.StateAfter = string(observed.Classification)
			result.Changes = append(result.Changes, reconcile.Change{Action: "remove-stale", Path: state.TransactionPath, Owner: "local-machine"})
			result.Validations = append(result.Validations, reconcile.Validation{Name: "stale-transaction-identity", Status: "passed"})
			return resultPayload(result), result, nil
		}
		if recoverable {
			observed, err = state.Inspect(root)
			if err != nil {
				return nil, reconcile.Result{}, err
			}
		}
	}
	allowed := observed.Classification == state.StateCurrent ||
		observed.Classification == state.StateDrifted ||
		observed.Classification == state.StateLegacyV1Compatible ||
		(observed.Classification == state.StatePartial && observed.Lock != nil)
	if !allowed {
		result := blockedLifecycle("repair", observed, "repair_requires_compatible_installation", "repair requires a readable current-version lock and compatible installation state", "run check and resolve the reported state before retrying", "check")
		return resultPayload(result), result, nil
	}
	if installed := state.AsString(observed.Lock["kit_version"]); installed != "" && installed != version.Version {
		result := blockedLifecycle("repair", observed, "version_change_requires_upgrade", fmt.Sprintf("installed kit version %s does not match running CLI version %s", installed, version.Version), "use upgrade with explicit version-change approval", "upgrade --dry-run")
		return resultPayload(result), result, nil
	}
	graph, err := compileObservedGraph(observed)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	lock, err := desiredLifecycleLock("repair", observed, graph, now)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	result, err := runLifecycle(lifecycleRequest{
		command:       "repair",
		root:          root,
		now:           now,
		dryRun:        dryRun,
		observed:      observed,
		graph:         graph,
		desiredLock:   lock,
		desiredConfig: repairConfig(root, observed),
		ensureIgnore:  true,
	})
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	payload := resultPayload(result)
	payload["migrated_lock"] = state.AsInt(observed.Lock["schema_version"]) == 1 && changeContains(result, state.LockPath)
	payload["repaired_managed_paths"] = recordsAsMaps(reconcile.ManagedPathRecords(graph))
	payload["kit_version"] = lock["kit_version"]
	return payload, result, nil
}
