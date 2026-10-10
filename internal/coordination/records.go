package coordination

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Attempt file names inside the deterministic attempt directory.
const (
	AttemptFile  = "attempt.json"
	StdoutFile   = "stdout.jsonl"
	StderrFile   = "stderr.txt"
	MessageFile  = "message.json"
	DeliveryFile = "delivery.json"
)

// Attempt statuses. Only response_captured means a final reply was captured; none of them
// means the assignment was accepted.
const (
	StatusPreparing             = "preparing"
	StatusLaunching             = "launching"
	StatusRunning               = "running"
	StatusResponseCaptured      = "response_captured"
	StatusCLIError              = "cli_error"
	StatusSessionMismatch       = "session_mismatch"
	StatusIncompleteOutput      = "incomplete_output"
	StatusUnsupportedInvocation = "unsupported_invocation"
	StatusExecutionFailed       = "execution_failed"
	StatusInterrupted           = "interrupted"
	StatusHelperInterrupted     = "helper_interrupted"
	StatusLaunchFailed          = "launch_failed"
	StatusSessionLocked         = "session_locked"
)

// Delivery statuses for the native return message.
const (
	DeliveryPending   = "pending"
	DeliverySent      = "sent"
	DeliveryRejected  = "rejected"
	DeliveryUncertain = "uncertain"
)

// ProcessIdentity records one process as soon as it is known.
type ProcessIdentity struct {
	PID       int    `json:"pid"`
	Host      string `json:"host"`
	Token     string `json:"token,omitempty"`
	StartedAt string `json:"started_at"`
}

// ChildRecord records the launched CLI process and how it ended.
type ChildRecord struct {
	PID       int    `json:"pid"`
	Host      string `json:"host"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at,omitempty"`
	ExitCode  *int   `json:"exit_code,omitempty"`
	Signal    string `json:"signal,omitempty"`
}

// ResponseLocator points at the original final reply without copying it.
type ResponseLocator struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	JSONField string `json:"json_field"`
	Chars     int    `json:"chars"`
}

// EarlierResult identifies one result record before the last, as the CLI reported it. Reply is
// set only when its result field is a JSON string; the text itself stays in stdout.jsonl.
type EarlierResult struct {
	Line    int              `json:"line"`
	Subtype string           `json:"subtype,omitempty"`
	IsError *bool            `json:"is_error,omitempty"`
	Reply   *ResponseLocator `json:"reply,omitempty"`
}

// Failed reports whether the CLI marked this earlier result as an error or not a success.
func (e EarlierResult) Failed() bool {
	return (e.IsError != nil && *e.IsError) || e.Subtype != "success"
}

// Denial keeps permission-denial metadata; the original tool input stays in stdout.jsonl.
type Denial struct {
	ToolName  string `json:"tool_name"`
	ToolUseID string `json:"tool_use_id,omitempty"`
}

// LockRelease records a deliberate, verified lock release.
type LockRelease struct {
	ReleasedAt         string          `json:"released_at"`
	ReleasedBy         ProcessIdentity `json:"released_by"`
	ManualVerification string          `json:"manual_verification,omitempty"`
}

// AttemptRecord is the single compact record for one invocation attempt.
type AttemptRecord struct {
	SchemaVersion     int              `json:"schema_version"`
	AssignmentID      string           `json:"assignment_id"`
	AttemptID         string           `json:"attempt_id"`
	Status            string           `json:"status"`
	Detail            string           `json:"detail,omitempty"`
	RequestFile       string           `json:"request_file"`
	RequestSHA256     string           `json:"request_sha256"`
	SessionMode       string           `json:"session_mode"`
	RequestedSession  string           `json:"requested_session_id"`
	ActualSession     string           `json:"actual_session_id,omitempty"`
	ReturnDestination Destination      `json:"return_destination"`
	RequestedModel    string           `json:"requested_model"`
	ActualModel       string           `json:"actual_model,omitempty"`
	CLIVersion        string           `json:"cli_version,omitempty"`
	Executable        string           `json:"executable"`
	Args              []string         `json:"args"`
	WorkingDirectory  string           `json:"working_directory"`
	Launcher          ProcessIdentity  `json:"launcher"`
	Child             *ChildRecord     `json:"child,omitempty"`
	LockOwner         *SessionLock     `json:"blocking_lock_owner,omitempty"`
	Response          *ResponseLocator `json:"response,omitempty"`
	// EarlierResults lists result records before the last one. Some CLI runs emit more than one;
	// an earlier record can hold the substantive reply or a failure a later record would hide.
	EarlierResults    []EarlierResult `json:"earlier_results,omitempty"`
	ResultSubtype     string          `json:"result_subtype,omitempty"`
	IsError           *bool           `json:"is_error,omitempty"`
	TerminalReason    string          `json:"terminal_reason,omitempty"`
	PermissionDenials []Denial        `json:"permission_denials"`
	Usage             json.RawMessage `json:"usage,omitempty"`
	ModelUsage        json.RawMessage `json:"model_usage,omitempty"`
	TotalCostUSD      *float64        `json:"total_cost_usd,omitempty"`
	DurationMS        *int64          `json:"duration_ms,omitempty"`
	NumTurns          *int            `json:"num_turns,omitempty"`
	UnparsedLines     int             `json:"unparsed_output_lines,omitempty"`
	LockRelease       *LockRelease    `json:"lock_release,omitempty"`
	UpdatedAt         string          `json:"updated_at"`
}

// SessionLock is the exclusive-create lock content for one session.
type SessionLock struct {
	SessionID    string          `json:"session_id"`
	AssignmentID string          `json:"assignment_id"`
	AttemptID    string          `json:"attempt_id"`
	AttemptDir   string          `json:"attempt_dir"`
	Launcher     ProcessIdentity `json:"launcher"`
	CreatedAt    string          `json:"created_at"`
}

// MessageArguments is loaded unchanged into the host's native send-message tool. Field names
// follow the tested Codex desktop send-message contract.
type MessageArguments struct {
	ThreadID string `json:"threadId"`
	HostID   string `json:"hostId"`
	Prompt   string `json:"prompt"`
}

// DeliveryRecord tracks the native return message separately from CLI execution.
type DeliveryRecord struct {
	Status             string       `json:"status"`
	PreparedRecipient  Destination  `json:"prepared_recipient"`
	ReportedRecipient  *Destination `json:"reported_recipient,omitempty"`
	RecipientConfirmed bool         `json:"recipient_confirmed"`
	RecipientMismatch  bool         `json:"recipient_mismatch,omitempty"`
	Receipt            string       `json:"receipt,omitempty"`
	ReceiptFile        string       `json:"receipt_file,omitempty"`
	CreatedAt          string       `json:"created_at"`
	UpdatedAt          string       `json:"updated_at"`
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func newToken() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("t%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buffer)
}

// NewSessionID allocates a random RFC 4122 version 4 UUID.
func NewSessionID() (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	buffer[6] = (buffer[6] & 0x0f) | 0x40
	buffer[8] = (buffer[8] & 0x3f) | 0x80
	h := hex.EncodeToString(buffer)
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32], nil
}

func currentIdentity() ProcessIdentity {
	host, _ := os.Hostname()
	return ProcessIdentity{PID: os.Getpid(), Host: host, Token: newToken(), StartedAt: timestamp()}
}

// writeJSONAtomic writes a private JSON file through a temporary sibling and rename.
func writeJSONAtomic(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	temp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	cleanup := func() { _ = os.Remove(tempName) }
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		cleanup()
		return err
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		cleanup()
		return err
	}
	if err := temp.Close(); err != nil {
		cleanup()
		return err
	}
	if err := os.Chmod(tempName, 0o600); err != nil {
		cleanup()
		return err
	}
	if err := os.Rename(tempName, path); err != nil {
		cleanup()
		return err
	}
	return nil
}

func readJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// ReadAttempt loads the compact attempt record from an attempt directory.
func ReadAttempt(attemptDir string) (AttemptRecord, error) {
	var record AttemptRecord
	err := readJSON(filepath.Join(attemptDir, AttemptFile), &record)
	return record, err
}
