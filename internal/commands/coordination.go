package commands

import (
	"errors"
	"fmt"
	"io"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/coordination"
)

// RunCoordination dispatches the narrow cross-tool invocation helper subcommands. Normal stdout
// carries only compact status and locators, never the reply text, brief or transcript.
func RunCoordination(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "codeheart-operating-kit coordination: error: the following arguments are required: subcommand")
		return 2
	}
	switch args[0] {
	case "invoke-claude":
		return runInvokeClaude(args[1:], stdout, stderr)
	case "record-delivery":
		return runRecordDelivery(args[1:], stdout, stderr)
	case "release-lock":
		return runReleaseLock(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "codeheart-operating-kit coordination: error: invalid subcommand %q (choose from invoke-claude, record-delivery, release-lock)\n", args[0])
		return 2
	}
}

func runInvokeClaude(args []string, stdout io.Writer, stderr io.Writer) int {
	values, _, positionals, err := parseValueArgs(args, map[string]bool{"--request": true, "--json": false})
	if err != nil {
		return writeArgError(stderr, "coordination invoke-claude", err)
	}
	if len(positionals) > 0 {
		return writeArgError(stderr, "coordination invoke-claude", fmt.Errorf("unexpected argument %q", positionals[0]))
	}
	if values["--request"] == "" {
		return writeArgError(stderr, "coordination invoke-claude", fmt.Errorf("--request is required"))
	}
	outcome, err := coordination.Invoke(coordination.Options{RequestPath: expandPath(values["--request"])})
	var requestErr *coordination.RequestError
	if err != nil && !errors.As(err, &requestErr) {
		outcome.Status = coordination.StatusHelperError
		outcome.PersistenceErrors = append(outcome.PersistenceErrors, err.Error())
	}
	_ = writeJSON(stdout, outcome)
	if requestErr != nil {
		return 2
	}
	if outcome.Status == coordination.StatusResponseCaptured && len(outcome.PersistenceErrors) == 0 {
		return 0
	}
	return 1
}

func runRecordDelivery(args []string, stdout io.Writer, stderr io.Writer) int {
	values, _, positionals, err := parseValueArgs(args, map[string]bool{
		"--attempt-dir":        true,
		"--status":             true,
		"--receipt-file":       true,
		"--reported-thread-id": true,
		"--reported-host-id":   true,
		"--json":               false,
	})
	if err != nil {
		return writeArgError(stderr, "coordination record-delivery", err)
	}
	if len(positionals) > 0 {
		return writeArgError(stderr, "coordination record-delivery", fmt.Errorf("unexpected argument %q", positionals[0]))
	}
	if values["--attempt-dir"] == "" || values["--status"] == "" {
		return writeArgError(stderr, "coordination record-delivery", fmt.Errorf("--attempt-dir and --status are required"))
	}
	record, err := coordination.RecordDelivery(coordination.DeliveryInput{
		AttemptDir:       expandPath(values["--attempt-dir"]),
		Status:           values["--status"],
		ReceiptFile:      expandPath(values["--receipt-file"]),
		ReportedThreadID: values["--reported-thread-id"],
		ReportedHostID:   values["--reported-host-id"],
	})
	if err != nil {
		_ = writeJSON(stdout, map[string]any{"status": "not_recorded", "detail": err.Error(), "delivery_status": record.Status})
		return 1
	}
	_ = writeJSON(stdout, map[string]any{
		"status":              "recorded",
		"delivery_status":     record.Status,
		"recipient_confirmed": record.RecipientConfirmed,
		"recipient_mismatch":  record.RecipientMismatch,
	})
	return 0
}

func runReleaseLock(args []string, stdout io.Writer, stderr io.Writer) int {
	values, _, positionals, err := parseValueArgs(args, map[string]bool{
		"--state-root":          true,
		"--session-id":          true,
		"--assignment-id":       true,
		"--attempt-id":          true,
		"--manual-verification": true,
		"--json":                false,
	})
	if err != nil {
		return writeArgError(stderr, "coordination release-lock", err)
	}
	if len(positionals) > 0 {
		return writeArgError(stderr, "coordination release-lock", fmt.Errorf("unexpected argument %q", positionals[0]))
	}
	for _, required := range []string{"--state-root", "--session-id", "--assignment-id", "--attempt-id"} {
		if values[required] == "" {
			return writeArgError(stderr, "coordination release-lock", fmt.Errorf("%s is required", required))
		}
	}
	_, err = coordination.ReleaseLock(coordination.ReleaseInput{
		StateRoot:          expandPath(values["--state-root"]),
		SessionID:          values["--session-id"],
		AssignmentID:       values["--assignment-id"],
		AttemptID:          values["--attempt-id"],
		ManualVerification: values["--manual-verification"],
	})
	if err != nil {
		var refusal *coordination.ReleaseRefusal
		if errors.As(err, &refusal) {
			_ = writeJSON(stdout, map[string]any{"status": "refused", "reason": refusal.Reason, "detail": refusal.Detail})
			return 1
		}
		_ = writeJSON(stdout, map[string]any{"status": "error", "detail": err.Error()})
		return 1
	}
	_ = writeJSON(stdout, map[string]any{"status": "released", "session_id": values["--session-id"], "attempt_id": values["--attempt-id"]})
	return 0
}
