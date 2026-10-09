package coordination

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// usageErrorPattern recognizes the CLI's own argument/usage rejection on stderr.
var usageErrorPattern = regexp.MustCompile(`(?i)(unknown option|unknown command|error: option|too many arguments|invalid choice|is not a valid choice|unexpected argument|requires --verbose)`)

// Options configures one invocation.
type Options struct {
	RequestPath string
	// Interrupt receives helper interruption signals. When nil, Invoke installs a handler for
	// the platform's ordinary termination signals.
	Interrupt <-chan os.Signal
}

// Outcome is the compact metadata printed to stdout; it never contains the reply text.
// MessageFile and DeliveryFile are set only when this invocation wrote them. A relay loads the
// MessageFile named here, never a fixed path. FallbackMessage carries a diagnostic addressed to
// the coordinator-prepared return destination when no new message file could be written.
type Outcome struct {
	Status            string            `json:"status"`
	Detail            string            `json:"detail,omitempty"`
	Problems          []string          `json:"problems,omitempty"`
	AssignmentID      string            `json:"assignment_id,omitempty"`
	AttemptID         string            `json:"attempt_id,omitempty"`
	SessionID         string            `json:"session_id,omitempty"`
	AttemptDir        string            `json:"attempt_dir,omitempty"`
	AttemptFile       string            `json:"attempt_file,omitempty"`
	MessageFile       string            `json:"message_file,omitempty"`
	DeliveryFile      string            `json:"delivery_file,omitempty"`
	FallbackMessage   *MessageArguments `json:"fallback_message,omitempty"`
	PersistenceErrors []string          `json:"persistence_errors,omitempty"`
	PermissionDenials int               `json:"permission_denials"`
	ChildExitCode     *int              `json:"child_exit_code,omitempty"`
}

// Outcome statuses that exist only in helper output, never in an attempt record.
const (
	StatusInvalidRequest = "invalid_request"
	StatusAttemptExists  = "attempt_exists"
	StatusHelperError    = "helper_error"
)

// afterStartHook lets tests simulate helper death before child identity is recorded.
var afterStartHook func()

// Invoke validates the request, launches the CLI once and records original evidence. A
// RequestError is returned before any attempt state exists. Every outcome after the attempt
// directory is created passes through one finish step that writes the attempt record, message
// arguments and a pending delivery record, except a refused replay, which changes nothing.
func Invoke(opts Options) (Outcome, error) {
	requestPath, err := filepath.Abs(opts.RequestPath)
	if err != nil {
		return Outcome{Status: StatusInvalidRequest, Problems: []string{err.Error()}}, &RequestError{Problems: []string{err.Error()}}
	}
	request, raw, err := LoadRequest(requestPath)
	if err != nil {
		outcome := Outcome{Status: StatusInvalidRequest, AssignmentID: request.AssignmentID, AttemptID: request.AttemptID}
		var requestErr *RequestError
		if errors.As(err, &requestErr) {
			outcome.Problems = requestErr.Problems
		}
		outcome.FallbackMessage = fallbackMessage(request, outcome.Status, "the request was rejected before launch: "+strings.Join(outcome.Problems, "; "))
		return outcome, err
	}
	digest := sha256.Sum256(raw)

	// Capture launcher identity before creating the attempt or the session lock.
	launcher := currentIdentity()

	attemptDir := AttemptDir(request.StateRoot, request.AssignmentID, request.AttemptID)
	outcome := Outcome{AssignmentID: request.AssignmentID, AttemptID: request.AttemptID, AttemptDir: attemptDir}
	helperError := func(detail string) (Outcome, error) {
		outcome.Status, outcome.Detail = StatusHelperError, detail
		outcome.FallbackMessage = fallbackMessage(request, outcome.Status, detail+"; no attempt record was written")
		return outcome, nil
	}

	sessionID := request.Session.ID
	if request.Session.Mode == "new" {
		if sessionID, err = NewSessionID(); err != nil {
			return helperError(fmt.Sprintf("allocate session id: %v", err))
		}
	}
	outcome.SessionID = sessionID

	if err := os.MkdirAll(filepath.Dir(attemptDir), 0o700); err != nil {
		return helperError(fmt.Sprintf("create attempt parent: %v", err))
	}
	if err := os.Mkdir(attemptDir, 0o700); err != nil {
		if errors.Is(err, os.ErrExist) {
			// A refused replay must not touch the earlier attempt, its message or its delivery.
			outcome.Status = StatusAttemptExists
			outcome.SessionID = ""
			outcome.AttemptFile = filepath.Join(attemptDir, AttemptFile)
			outcome.Detail = "this attempt was already started; nothing was launched and the existing attempt, message and delivery records were left unchanged"
			outcome.FallbackMessage = fallbackMessage(request, outcome.Status, outcome.Detail+"; this is a new refusal notice, not the earlier result. Existing record: "+outcome.AttemptFile)
			return outcome, nil
		}
		return helperError(fmt.Sprintf("create attempt directory: %v", err))
	}
	attemptFile := filepath.Join(attemptDir, AttemptFile)

	record := &AttemptRecord{
		SchemaVersion:     RequestSchemaVersion,
		AssignmentID:      request.AssignmentID,
		AttemptID:         request.AttemptID,
		Status:            StatusPreparing,
		RequestFile:       requestPath,
		RequestSHA256:     hex.EncodeToString(digest[:]),
		SessionMode:       request.Session.Mode,
		RequestedSession:  sessionID,
		ReturnDestination: request.ReturnDestination,
		RequestedModel:    request.Model,
		Executable:        request.Executable,
		Args:              request.Args(sessionID),
		WorkingDirectory:  request.WorkingDirectory,
		Launcher:          launcher,
		PermissionDenials: []Denial{},
	}
	save := func() error {
		record.UpdatedAt = timestamp()
		return writeJSONAtomic(attemptFile, record)
	}
	if err := save(); err != nil {
		return helperError(fmt.Sprintf("write attempt record: %v", err))
	}
	outcome.AttemptFile = attemptFile

	lockPath := SessionLockPath(request.StateRoot, sessionID)
	var parsed streamResult
	// finish writes the final record, message arguments and pending delivery for this attempt.
	// It releases the session lock only when this launcher owns it and no child can still run.
	finish := func(releaseLock bool) (Outcome, error) {
		outcome.Status, outcome.Detail = record.Status, record.Detail
		outcome.PermissionDenials = len(record.PermissionDenials)
		if err := save(); err != nil {
			outcome.PersistenceErrors = append(outcome.PersistenceErrors, "attempt record: "+err.Error())
		}
		message := MessageArguments{
			ThreadID: request.ReturnDestination.ThreadID,
			HostID:   request.ReturnDestination.HostID,
			Prompt:   composeMessage(record, parsed, attemptDir, sessionID),
		}
		messagePath := filepath.Join(attemptDir, MessageFile)
		if err := writeJSONAtomic(messagePath, message); err != nil {
			outcome.PersistenceErrors = append(outcome.PersistenceErrors, "message file: "+err.Error())
			outcome.FallbackMessage = &message
		} else {
			outcome.MessageFile = messagePath
		}
		now := timestamp()
		deliveryPath := filepath.Join(attemptDir, DeliveryFile)
		// Persist pending before any send is attempted.
		if err := writeJSONAtomic(deliveryPath, DeliveryRecord{
			Status:            DeliveryPending,
			PreparedRecipient: request.ReturnDestination,
			CreatedAt:         now,
			UpdatedAt:         now,
		}); err != nil {
			outcome.PersistenceErrors = append(outcome.PersistenceErrors, "delivery record: "+err.Error())
		} else {
			outcome.DeliveryFile = deliveryPath
		}
		if releaseLock {
			if err := releaseOwnLock(lockPath, launcher); err != nil {
				outcome.PersistenceErrors = append(outcome.PersistenceErrors, "session lock not released: "+err.Error())
			}
		}
		return outcome, nil
	}

	lock := SessionLock{
		SessionID:    sessionID,
		AssignmentID: request.AssignmentID,
		AttemptID:    request.AttemptID,
		AttemptDir:   attemptDir,
		Launcher:     launcher,
		CreatedAt:    timestamp(),
	}
	if owner, err := acquireLock(lockPath, lock); err != nil {
		if owner != nil {
			record.Status = StatusSessionLocked
			record.Detail = "another participating attempt holds this session; no second writer was launched"
			record.LockOwner = owner
		} else {
			record.Status = StatusLaunchFailed
			record.Detail = fmt.Sprintf("session lock could not be created: %v", err)
		}
		return finish(false)
	}

	stdoutPath := filepath.Join(attemptDir, StdoutFile)
	stderrPath := filepath.Join(attemptDir, StderrFile)
	launchFail := func(detail string) (Outcome, error) {
		// No child exists, and this launcher owns the lock: releasing it is safe.
		record.Status, record.Detail = StatusLaunchFailed, detail
		return finish(true)
	}
	stdoutFile, err := os.OpenFile(stdoutPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return launchFail(fmt.Sprintf("open stdout file: %v", err))
	}
	defer stdoutFile.Close()
	stderrFile, err := os.OpenFile(stderrPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return launchFail(fmt.Sprintf("open stderr file: %v", err))
	}
	defer stderrFile.Close()
	brief, err := os.Open(request.BriefFile)
	if err != nil {
		return launchFail(fmt.Sprintf("open brief: %v", err))
	}
	defer brief.Close()

	record.Status = StatusLaunching
	if err := save(); err != nil {
		return launchFail(fmt.Sprintf("persist launch state: %v", err))
	}

	// The child writes directly to the attempt files, not to a helper-owned pipe.
	cmd := exec.Command(request.Executable, record.Args...)
	cmd.Dir = request.WorkingDirectory
	cmd.Stdin = brief
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		return launchFail(fmt.Sprintf("start CLI: %v", err))
	}
	if afterStartHook != nil {
		afterStartHook()
	}
	host, _ := os.Hostname()
	record.Child = &ChildRecord{PID: cmd.Process.Pid, Host: host, StartedAt: timestamp()}
	record.Status = StatusRunning
	if err := save(); err != nil {
		outcome.PersistenceErrors = append(outcome.PersistenceErrors, "running state: "+err.Error())
	}

	interrupt := opts.Interrupt
	if interrupt == nil {
		ch := make(chan os.Signal, 1)
		notifyInterrupt(ch)
		interrupt = ch
	}
	waitDone := make(chan error, 1)
	go func() { waitDone <- cmd.Wait() }()

	var waitErr error
	select {
	case waitErr = <-waitDone:
	case sig := <-interrupt:
		// Do not kill or abandon evidence: the child may still be writing. Keep the lock.
		alive, aliveErr := processAlive(record.Child.PID)
		state := "child still running"
		if aliveErr != nil {
			state = fmt.Sprintf("child state unknown: %v", aliveErr)
		} else if !alive {
			state = "child no longer running"
		}
		record.Status = StatusHelperInterrupted
		record.Detail = fmt.Sprintf("helper received %v; %s; effects are uncertain and the session lock is retained", sig, state)
		return finish(false)
	}

	record.Child.EndedAt = timestamp()
	exitCode := -1
	signaled := ""
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
		if exitCode == -1 {
			signaled = cmd.ProcessState.String()
		}
	}
	record.Child.ExitCode = &exitCode
	record.Child.Signal = signaled
	_ = stdoutFile.Sync()
	_ = stderrFile.Sync()

	parsed = parseStream(stdoutPath)
	applyParsed(record, parsed, stdoutPath)
	record.Status, record.Detail = classify(record, parsed, exitCode, signaled, waitErr, sessionID, stderrPath)
	outcome.ChildExitCode = &exitCode
	// The child has exited and this launcher owns the lock.
	return finish(true)
}

// fallbackMessage addresses a diagnostic to the coordinator-prepared return destination when no
// attempt message file exists for this invocation. It returns nil when the request supplied no
// usable destination; the relay then reports only to its commissioning coordinator.
func fallbackMessage(request Request, status string, detail string) *MessageArguments {
	if !plainValue(request.ReturnDestination.ThreadID) || !plainValue(request.ReturnDestination.HostID) {
		return nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Claude CLI attempt %s for assignment %s: %s.\n", request.AttemptID, request.AssignmentID, status)
	fmt.Fprintf(&b, "Detail: %s\n", detail)
	b.WriteString("Transport status only; no CLI reply is attached and nothing was retried.")
	return &MessageArguments{ThreadID: request.ReturnDestination.ThreadID, HostID: request.ReturnDestination.HostID, Prompt: b.String()}
}

func acquireLock(path string, lock SessionLock) (*SessionLock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			var owner SessionLock
			if readErr := readJSON(path, &owner); readErr != nil {
				owner = SessionLock{SessionID: lock.SessionID}
			}
			return &owner, err
		}
		return nil, err
	}
	data, _ := json.MarshalIndent(lock, "", "  ")
	_, writeErr := file.Write(append(data, '\n'))
	syncErr := file.Sync()
	closeErr := file.Close()
	if err := errors.Join(writeErr, syncErr, closeErr); err != nil {
		// This launcher created the lock exclusively and no child exists: do not leave it behind.
		_ = os.Remove(path)
		return nil, err
	}
	return nil, nil
}

func releaseOwnLock(path string, launcher ProcessIdentity) error {
	var owner SessionLock
	if err := readJSON(path, &owner); err != nil {
		return err
	}
	if owner.Launcher.Token != launcher.Token || owner.Launcher.PID != launcher.PID {
		return fmt.Errorf("lock owner changed")
	}
	return os.Remove(path)
}

type streamResult struct {
	initSession   string
	initModel     string
	cliVersion    string
	hasResult     bool
	resultLine    int
	result        map[string]json.RawMessage
	reply         string
	replyIsString bool
	// earlier holds text replies from result records before the last one, and denials gathers
	// permission denials from every result record, so a later record cannot hide them.
	earlier    []earlierResult
	denials    []Denial
	unparsed   int
	totalBytes int64
}

type earlierResult struct {
	line  int
	chars int
}

func parseStream(path string) streamResult {
	var parsed streamResult
	file, err := os.Open(path)
	if err != nil {
		return parsed
	}
	defer file.Close()
	if info, err := file.Stat(); err == nil {
		parsed.totalBytes = info.Size()
	}
	reader := bufio.NewReader(file)
	line := 0
	for {
		text, readErr := reader.ReadString('\n')
		if len(text) > 0 {
			line++
			trimmed := strings.TrimSpace(text)
			if trimmed != "" {
				var fields map[string]json.RawMessage
				if json.Unmarshal([]byte(trimmed), &fields) != nil {
					parsed.unparsed++
				} else {
					var kind, subtype string
					_ = json.Unmarshal(fields["type"], &kind)
					_ = json.Unmarshal(fields["subtype"], &subtype)
					if kind == "system" && subtype == "init" && parsed.initSession == "" {
						_ = json.Unmarshal(fields["session_id"], &parsed.initSession)
						_ = json.Unmarshal(fields["model"], &parsed.initModel)
						_ = json.Unmarshal(fields["claude_code_version"], &parsed.cliVersion)
					}
					if kind == "result" {
						if parsed.hasResult && parsed.replyIsString {
							parsed.earlier = append(parsed.earlier, earlierResult{line: parsed.resultLine, chars: len([]rune(parsed.reply))})
						}
						var denials []Denial
						if json.Unmarshal(fields["permission_denials"], &denials) == nil {
							parsed.denials = appendDenials(parsed.denials, denials)
						}
						parsed.hasResult = true
						parsed.resultLine = line
						parsed.result = fields
						parsed.reply = ""
						parsed.replyIsString = isJSONString(fields["result"]) && json.Unmarshal(fields["result"], &parsed.reply) == nil
					}
				}
			}
		}
		if readErr != nil {
			break
		}
	}
	return parsed
}

func applyParsed(record *AttemptRecord, parsed streamResult, stdoutPath string) {
	record.ActualModel = parsed.initModel
	record.CLIVersion = parsed.cliVersion
	record.ActualSession = parsed.initSession
	record.UnparsedLines = parsed.unparsed
	if !parsed.hasResult {
		return
	}
	fields := parsed.result
	var session string
	if json.Unmarshal(fields["session_id"], &session) == nil && session != "" {
		record.ActualSession = session
	}
	var subtype, terminal string
	_ = json.Unmarshal(fields["subtype"], &subtype)
	_ = json.Unmarshal(fields["terminal_reason"], &terminal)
	record.ResultSubtype, record.TerminalReason = subtype, terminal
	var isError bool
	if json.Unmarshal(fields["is_error"], &isError) == nil {
		record.IsError = &isError
	}
	if parsed.denials != nil {
		record.PermissionDenials = parsed.denials
	}
	for _, earlier := range parsed.earlier {
		record.EarlierResponses = append(record.EarlierResponses, ResponseLocator{
			File:      stdoutPath,
			Line:      earlier.line,
			JSONField: "result",
			Chars:     earlier.chars,
		})
	}
	if raw, ok := fields["usage"]; ok && string(raw) != "null" {
		record.Usage = raw
	}
	if raw, ok := fields["modelUsage"]; ok && string(raw) != "null" {
		record.ModelUsage = raw
	}
	var cost float64
	if json.Unmarshal(fields["total_cost_usd"], &cost) == nil {
		record.TotalCostUSD = &cost
	}
	var duration int64
	if json.Unmarshal(fields["duration_ms"], &duration) == nil {
		record.DurationMS = &duration
	}
	var turns int
	if json.Unmarshal(fields["num_turns"], &turns) == nil {
		record.NumTurns = &turns
	}
	if parsed.replyIsString {
		record.Response = &ResponseLocator{
			File:      stdoutPath,
			Line:      parsed.resultLine,
			JSONField: "result",
			Chars:     len([]rune(parsed.reply)),
		}
	}
}

func classify(record *AttemptRecord, parsed streamResult, exitCode int, signaled string, waitErr error, expectedSession string, stderrPath string) (string, string) {
	if signaled != "" {
		return StatusInterrupted, "CLI process ended by " + signaled + "; inspect partial output before resuming"
	}
	if !parsed.hasResult {
		if exitCode != 0 {
			stderr, _ := os.ReadFile(stderrPath)
			if usageErrorPattern.Match(stderr) {
				return StatusUnsupportedInvocation, "the CLI rejected the fixed invocation flags; see original stderr"
			}
			detail := fmt.Sprintf("CLI exited %d without a final result record", exitCode)
			if waitErr != nil {
				detail += ": " + waitErr.Error()
			}
			return StatusExecutionFailed, detail
		}
		if parsed.totalBytes == 0 {
			return StatusIncompleteOutput, "CLI exited 0 without output"
		}
		return StatusIncompleteOutput, "CLI output has no final result record; treat completion as uncertain"
	}
	if record.ActualSession != expectedSession {
		return StatusSessionMismatch, fmt.Sprintf("expected session %s, CLI reported %q", expectedSession, record.ActualSession)
	}
	if (record.IsError != nil && *record.IsError) || record.ResultSubtype != "success" || exitCode != 0 {
		return StatusCLIError, fmt.Sprintf("CLI result subtype %q, is_error %v, exit %d", record.ResultSubtype, record.IsError != nil && *record.IsError, exitCode)
	}
	if !parsed.replyIsString {
		return StatusIncompleteOutput, "final result record has no text reply"
	}
	return StatusResponseCaptured, ""
}

func composeMessage(record *AttemptRecord, parsed streamResult, attemptDir string, sessionID string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Claude CLI attempt %s for assignment %s: %s.\n", record.AttemptID, record.AssignmentID, record.Status)
	fmt.Fprintf(&b, "Session: %s. Permission denials: %d.\n", sessionID, len(record.PermissionDenials))
	if record.Detail != "" {
		fmt.Fprintf(&b, "Detail: %s\n", record.Detail)
	}
	fmt.Fprintf(&b, "Attempt record: %s\n", filepath.Join(attemptDir, AttemptFile))
	if len(parsed.earlier) > 0 {
		lines := make([]string, 0, len(parsed.earlier))
		for _, earlier := range parsed.earlier {
			lines = append(lines, strconv.Itoa(earlier.line))
		}
		fmt.Fprintf(&b, "Earlier result records with original replies: %s, line(s) %s, JSON field \"result\". Read them too; only the last reply is quoted below.\n", filepath.Join(attemptDir, StdoutFile), strings.Join(lines, ", "))
	}
	b.WriteString("Transport status only; this is not task acceptance.\n")
	switch {
	case parsed.replyIsString && len([]rune(parsed.reply)) <= InlineReplyLimit:
		b.WriteString("Original reply (unchanged):\n")
		b.WriteString(parsed.reply)
	case parsed.replyIsString:
		fmt.Fprintf(&b, "Original reply (%d characters): %s, line %d, JSON field \"result\".", len([]rune(parsed.reply)), filepath.Join(attemptDir, StdoutFile), parsed.resultLine)
	default:
		b.WriteString("No final reply was captured.")
		if record.Child == nil {
			b.WriteString(" No CLI process was recorded as started for this attempt.")
		} else {
			fmt.Fprintf(&b, " Original CLI output: %s; stderr: %s.", filepath.Join(attemptDir, StdoutFile), filepath.Join(attemptDir, StderrFile))
		}
		if record.LockOwner != nil {
			fmt.Fprintf(&b, " Session held by attempt %s of assignment %s.", record.LockOwner.AttemptID, record.LockOwner.AssignmentID)
		}
	}
	return b.String()
}

// appendDenials adds denials not already present; a tool use ID identifies a repeated report of
// the same denial across result records. Denials without an ID are always kept.
func appendDenials(existing []Denial, more []Denial) []Denial {
	if existing == nil {
		existing = []Denial{}
	}
	for _, denial := range more {
		duplicate := false
		if denial.ToolUseID != "" {
			for _, known := range existing {
				if known.ToolUseID == denial.ToolUseID && known.ToolName == denial.ToolName {
					duplicate = true
					break
				}
			}
		}
		if !duplicate {
			existing = append(existing, denial)
		}
	}
	return existing
}

// isJSONString reports whether raw is a JSON string value. Missing, null and other values are
// not a captured reply; an actual empty string is.
func isJSONString(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) > 0 && trimmed[0] == '"'
}
