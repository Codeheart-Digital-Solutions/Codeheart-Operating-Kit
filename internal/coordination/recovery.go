package coordination

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DeliveryInput is the relay's factual report of one native send attempt.
type DeliveryInput struct {
	AttemptDir       string
	Status           string
	ReceiptFile      string
	ReportedThreadID string
	ReportedHostID   string
}

// RecordDelivery moves a pending delivery record to sent, rejected or uncertain exactly once.
// A reported recipient that differs from the prepared one is retained as uncertain delivery;
// an omitted reported recipient never confirms the recipient.
func RecordDelivery(input DeliveryInput) (DeliveryRecord, error) {
	switch input.Status {
	case DeliverySent, DeliveryRejected, DeliveryUncertain:
	default:
		return DeliveryRecord{}, fmt.Errorf("status must be sent, rejected or uncertain")
	}
	path := filepath.Join(input.AttemptDir, DeliveryFile)
	var record DeliveryRecord
	if err := readJSON(path, &record); err != nil {
		return DeliveryRecord{}, fmt.Errorf("read delivery record: %w", err)
	}
	if record.Status != DeliveryPending {
		return record, fmt.Errorf("delivery already recorded as %s; the original record is preserved", record.Status)
	}
	if input.ReceiptFile != "" {
		receipt, err := os.ReadFile(input.ReceiptFile)
		if err != nil {
			return record, fmt.Errorf("read receipt: %w", err)
		}
		record.Receipt = string(receipt)
		record.ReceiptFile = input.ReceiptFile
	}
	record.Status = input.Status
	if input.ReportedThreadID != "" || input.ReportedHostID != "" {
		reported := Destination{ThreadID: input.ReportedThreadID, HostID: input.ReportedHostID}
		record.ReportedRecipient = &reported
		threadMatches := reported.ThreadID == record.PreparedRecipient.ThreadID
		hostMatches := reported.HostID == "" || reported.HostID == record.PreparedRecipient.HostID
		if threadMatches && hostMatches {
			record.RecipientConfirmed = input.Status == DeliverySent
		} else {
			record.RecipientMismatch = true
			if record.Status == DeliverySent {
				record.Status = DeliveryUncertain
			}
		}
	}
	record.UpdatedAt = timestamp()
	if err := writeJSONAtomic(path, record); err != nil {
		return record, err
	}
	return record, nil
}

// ReleaseInput names the exact lock the coordinator deliberately wants to release.
type ReleaseInput struct {
	StateRoot          string
	SessionID          string
	AssignmentID       string
	AttemptID          string
	ManualVerification string
}

// ReleaseRefusal explains why a lock was retained.
type ReleaseRefusal struct {
	Reason string
	Detail string
}

func (r *ReleaseRefusal) Error() string { return r.Reason + ": " + r.Detail }

// ReleaseLock removes a retained session lock only after verifying the named owner and that
// the recorded launcher and child have ended. When the child identity was never recorded, it
// refuses unless the coordinator supplies a manual verification statement. It never guesses
// staleness from age.
func ReleaseLock(input ReleaseInput) (AttemptRecord, error) {
	if !uuidPattern.MatchString(input.SessionID) || !identifierPattern.MatchString(input.AssignmentID) || !identifierPattern.MatchString(input.AttemptID) || !filepath.IsAbs(input.StateRoot) {
		return AttemptRecord{}, &ReleaseRefusal{Reason: "invalid_input", Detail: "state root must be absolute and session, assignment and attempt identifiers valid"}
	}
	lockPath := SessionLockPath(input.StateRoot, input.SessionID)
	var owner SessionLock
	if err := readJSON(lockPath, &owner); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return AttemptRecord{}, &ReleaseRefusal{Reason: "no_lock", Detail: "no lock exists for this session"}
		}
		return AttemptRecord{}, &ReleaseRefusal{Reason: "owner_unknown", Detail: fmt.Sprintf("lock content unreadable: %v", err)}
	}
	if owner.SessionID != input.SessionID || owner.AssignmentID != input.AssignmentID || owner.AttemptID != input.AttemptID {
		return AttemptRecord{}, &ReleaseRefusal{Reason: "owner_mismatch", Detail: fmt.Sprintf("lock belongs to assignment %q attempt %q", owner.AssignmentID, owner.AttemptID)}
	}
	attemptDir := AttemptDir(input.StateRoot, input.AssignmentID, input.AttemptID)
	record, err := ReadAttempt(attemptDir)
	if err != nil {
		return AttemptRecord{}, &ReleaseRefusal{Reason: "attempt_unknown", Detail: fmt.Sprintf("attempt record unreadable: %v", err)}
	}
	if record.Launcher.Token != owner.Launcher.Token || record.Launcher.PID != owner.Launcher.PID {
		return record, &ReleaseRefusal{Reason: "owner_mismatch", Detail: "attempt launcher does not match the lock owner"}
	}
	host, _ := os.Hostname()
	if owner.Launcher.Host != host {
		return record, &ReleaseRefusal{Reason: "other_host", Detail: fmt.Sprintf("lock was created on host %q; verify and release it there", owner.Launcher.Host)}
	}
	if alive, err := processAlive(owner.Launcher.PID); err != nil || alive {
		return record, &ReleaseRefusal{Reason: "launcher_alive_or_unknown", Detail: aliveDetail("launcher", owner.Launcher.PID, alive, err)}
	}
	manual := strings.TrimSpace(input.ManualVerification)
	if record.Child == nil {
		if manual == "" {
			return record, &ReleaseRefusal{Reason: "child_identity_unknown", Detail: "the launcher ended before recording the CLI process; verify manually that no CLI process for this session is running, then repeat with a manual verification statement or hand over to a new session"}
		}
	} else {
		if record.Child.Host != host {
			return record, &ReleaseRefusal{Reason: "other_host", Detail: "child process was recorded on another host"}
		}
		if alive, err := processAlive(record.Child.PID); err != nil || alive {
			return record, &ReleaseRefusal{Reason: "child_alive_or_unknown", Detail: aliveDetail("CLI child", record.Child.PID, alive, err)}
		}
	}
	record.LockRelease = &LockRelease{ReleasedAt: timestamp(), ReleasedBy: currentIdentity(), ManualVerification: manual}
	record.UpdatedAt = timestamp()
	if err := writeJSONAtomic(filepath.Join(attemptDir, AttemptFile), record); err != nil {
		return record, err
	}
	if err := os.Remove(lockPath); err != nil {
		return record, err
	}
	return record, nil
}

func aliveDetail(name string, pid int, alive bool, err error) string {
	if err != nil {
		return fmt.Sprintf("%s pid %d state unknown: %v", name, pid, err)
	}
	return fmt.Sprintf("%s pid %d is still running (a reused process id also counts as running)", name, pid)
}
