package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/capabilities"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/lockfile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

var routeTargetPattern = regexp.MustCompile(`\.codeheart/kit/[A-Za-z0-9._/\-]+\.md`)

func RunCheck(args []string, stdout io.Writer, stderr io.Writer) int {
	_, bools, positionals, err := parseValueArgs(args, map[string]bool{"--json": false})
	if err != nil {
		return writeArgError(stderr, "check", err)
	}
	path := "."
	if len(positionals) > 0 {
		path = positionals[0]
	}
	if len(positionals) > 1 {
		return writeArgError(stderr, "check", fmt.Errorf("unexpected argument %q", positionals[1]))
	}
	result, err := CheckRepository(path)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit check: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, result); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit check: error: %v\n", err)
			return 1
		}
	} else if result["ok"] == true {
		fmt.Fprintf(stdout, "Operating Kit check: healthy (%s). OK: true.\n", result["state"])
	} else {
		operation := operationFromPayload(result)
		reconcile.WriteText(stdout, operation)
	}
	if ok, _ := result["ok"].(bool); ok {
		return 0
	}
	return 1
}

func CheckRepository(path string) (map[string]any, error) {
	root := expandPath(path)
	observed, err := state.Inspect(root)
	if err != nil {
		return nil, err
	}
	operation := reconcile.NewResult("check")
	operation.StateBefore = string(observed.Classification)
	operation.StateAfter = string(observed.Classification)
	operation.Validations = append(operation.Validations, reconcile.Validation{Name: "installed-state", Status: "passed"})

	if observed.Classification != state.StateCurrent {
		operation.Blockers = append(operation.Blockers, reconcile.Blocker{
			Code:         "state_" + string(observed.Classification),
			Message:      stateMessage(observed),
			Remediation:  stateRemediation(observed.Classification),
			RetryCommand: stateRetry(observed.Classification),
		})
	}
	if state.AsInt(observed.Lock["schema_version"]) == 1 {
		operation.Blockers = append(operation.Blockers, reconcile.Blocker{
			Code:         "legacy_lock_requires_repair",
			Message:      "lock v1 lacks completed-operation and release-provenance evidence",
			Remediation:  "preview the bounded migration, then run repair",
			RetryCommand: "repair --dry-run",
		})
	}
	missingCLI := !CLIAvailable()
	if missingCLI {
		operation.Blockers = append(operation.Blockers, reconcile.Blocker{
			Code:         "cli_unavailable",
			Message:      "codeheart-operating-kit is not available to this process",
			Remediation:  "restore the recorded CLI path before lifecycle work",
			RetryCommand: "check",
		})
	}
	agentsText := ""
	if data, err := os.ReadFile(filepath.Join(root, "AGENTS.md")); err == nil {
		agentsText = string(data)
	}
	missingTargets := MissingRouteTargets(root, agentsText)
	missingRouting := !containsManagedMarkers(agentsText) || len(missingTargets) > 0
	if missingRouting && observed.Classification == state.StateCurrent {
		operation.Blockers = append(operation.Blockers, reconcile.Blocker{
			Code:         "routing_incomplete",
			Message:      "root routing markers or managed route targets are missing",
			Remediation:  "run repair to restore the managed routing block and targets",
			RetryCommand: "repair",
		})
	}
	if len(operation.Blockers) == 0 {
		operation.Status = reconcile.StatusSucceeded
	} else {
		operation.Status = reconcile.StatusBlocked
	}

	driftFindings := []any{}
	for _, missing := range observed.MissingPaths {
		driftFindings = append(driftFindings, map[string]any{"path": missing, "status": "missing"})
	}
	for _, changed := range observed.DriftedPaths {
		driftFindings = append(driftFindings, map[string]any{"path": changed, "status": "modified"})
	}
	for _, stateError := range observed.Errors {
		driftFindings = append(driftFindings, map[string]any{"path": state.LockPath, "status": "invalid", "detail": stateError})
	}
	missingLockMetadata := []string{}
	if observed.Lock != nil {
		missingLockMetadata = lockfile.MissingRequiredLockMetadata(observed.Lock)
	}
	native := observed.Lock["native_capabilities"]
	if native == nil {
		native = capabilities.UnknownNativeCapabilityState(lockfile.UTCNow())
	}
	payload := resultPayload(operation)
	payload["ok"] = operation.Status == reconcile.StatusSucceeded
	payload["state"] = string(observed.Classification)
	payload["traits"] = stringsAsAny(observed.Traits)
	payload["missing_cli"] = missingCLI
	payload["stale_cli"] = state.AsString(observed.Lock["kit_version"]) != "" && state.AsString(observed.Lock["kit_version"]) != version.Version
	payload["missing_routing"] = missingRouting
	payload["missing_route_targets"] = stringsAsAny(missingTargets)
	payload["missing_lock_metadata"] = stringsAsAny(missingLockMetadata)
	payload["drift"] = driftFindings
	payload["native_capabilities"] = native
	return payload, nil
}

func operationFromPayload(payload map[string]any) reconcile.Result {
	result := state.Map(payload["result"])
	operation := reconcile.NewResult(state.AsString(result["command"]))
	operation.Status = reconcile.Status(state.AsString(result["status"]))
	operation.StateBefore = state.AsString(result["state_before"])
	for _, item := range state.AnySlice(result["blockers"]) {
		blocker := state.Map(item)
		operation.Blockers = append(operation.Blockers, reconcile.Blocker{
			Code:         state.AsString(blocker["code"]),
			Message:      state.AsString(blocker["message"]),
			Path:         state.AsString(blocker["path"]),
			Remediation:  state.AsString(blocker["remediation"]),
			RetryCommand: state.AsString(blocker["retry_command"]),
		})
	}
	return operation
}

func stateMessage(observed state.Observed) string {
	if len(observed.Errors) > 0 {
		return observed.Errors[0]
	}
	if len(observed.MissingPaths) > 0 {
		return fmt.Sprintf("installation is %s; first missing path is %s", observed.Classification, observed.MissingPaths[0])
	}
	if len(observed.DriftedPaths) > 0 {
		return fmt.Sprintf("installation is %s; first modified path is %s", observed.Classification, observed.DriftedPaths[0])
	}
	return "installation state is " + string(observed.Classification)
}

func stateRemediation(classification state.Classification) string {
	switch classification {
	case state.StateDrifted, state.StatePartial, state.StateLegacyV1Compatible:
		return "preview repair, then repair the compatible installation"
	case state.StateStaleCLI:
		return "use upgrade with explicit approval for a version change"
	case state.StateTransactionInProgress:
		return "wait for the active lifecycle operation to finish"
	case state.StateRecoveryRequired:
		return "preserve recovery files and run repair"
	case state.StateAbsent, state.StateAdoptable:
		return "run init if this repository should use the Operating Kit"
	default:
		return "resolve the invalid or unsupported state before mutation"
	}
}

func stateRetry(classification state.Classification) string {
	switch classification {
	case state.StateDrifted, state.StatePartial, state.StateLegacyV1Compatible, state.StateRecoveryRequired:
		return "repair --dry-run"
	case state.StateAbsent, state.StateAdoptable:
		return "init --dry-run"
	case state.StateStaleCLI:
		return "upgrade --dry-run"
	default:
		return "check"
	}
}

func containsManagedMarkers(text string) bool {
	return strings.Contains(text, reconcile.BeginMarker) && strings.Contains(text, reconcile.EndMarker)
}

func MissingRouteTargets(root string, agentsText string) []string {
	targets := routeTargetPattern.FindAllString(agentsText, -1)
	seen := map[string]bool{}
	missing := []string{}
	for _, target := range targets {
		if seen[target] {
			continue
		}
		seen[target] = true
		if !exists(filepath.Join(root, filepath.FromSlash(target))) {
			missing = append(missing, target)
		}
	}
	sort.Strings(missing)
	return missing
}

func CLIAvailable() bool {
	if os.Getenv("CODEHEART_OPERATING_KIT_CLI") == "1" {
		return true
	}
	if _, err := exec.LookPath("codeheart-operating-kit"); err == nil {
		return true
	}
	executable := filepath.Base(os.Args[0])
	return (executable == "codeheart-operating-kit" || executable == "codeheart-operating-kit.exe") && exists(os.Args[0])
}
