package reconcile

import (
	"encoding/json"
	"fmt"
	"io"
)

type Status string

const (
	StatusPlanned          Status = "planned"
	StatusSucceeded        Status = "succeeded"
	StatusBlocked          Status = "blocked"
	StatusFailed           Status = "failed"
	StatusRolledBack       Status = "rolled-back"
	StatusRecoveryRequired Status = "recovery-required"
)

type Change struct {
	Action string `json:"action"`
	Path   string `json:"path"`
	Owner  string `json:"owner,omitempty"`
}

type Blocker struct {
	Code         string `json:"code"`
	Message      string `json:"message"`
	Path         string `json:"path,omitempty"`
	Remediation  string `json:"remediation,omitempty"`
	RetryCommand string `json:"retry_command,omitempty"`
}

type Validation struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type Rollback struct {
	Attempted bool   `json:"attempted"`
	Succeeded bool   `json:"succeeded"`
	Detail    string `json:"detail,omitempty"`
}

type Result struct {
	SchemaVersion int            `json:"schema_version"`
	Command       string         `json:"command"`
	Status        Status         `json:"status"`
	DryRun        bool           `json:"dry_run"`
	StateBefore   string         `json:"state_before"`
	StateAfter    string         `json:"state_after,omitempty"`
	TransactionID string         `json:"transaction_id,omitempty"`
	Changes       []Change       `json:"changes"`
	Blockers      []Blocker      `json:"blockers"`
	Validations   []Validation   `json:"validations"`
	Provenance    map[string]any `json:"provenance,omitempty"`
	Rollback      Rollback       `json:"rollback"`
}

func NewResult(command string) Result {
	return Result{
		SchemaVersion: 1,
		Command:       command,
		Changes:       []Change{},
		Blockers:      []Blocker{},
		Validations:   []Validation{},
		Provenance:    map[string]any{},
	}
}

func Preview(plan Plan) Result {
	result := NewResult(plan.Command)
	result.DryRun = true
	result.StateBefore = plan.StateBefore
	result.TransactionID = plan.ID
	result.Changes = planChanges(plan)
	if len(plan.Blockers) > 0 {
		result.Status = StatusBlocked
		result.Blockers = append(result.Blockers, plan.Blockers...)
		return result
	}
	result.Status = StatusPlanned
	result.StateAfter = plan.StateBefore
	result.Validations = append(result.Validations, Validation{Name: "change-plan", Status: "passed"})
	return result
}

func (result Result) OK() bool {
	return result.Status == StatusPlanned || result.Status == StatusSucceeded
}

func (result Result) ToMap() map[string]any {
	data, _ := json.Marshal(result)
	var mapping map[string]any
	_ = json.Unmarshal(data, &mapping)
	return mapping
}

func WriteJSON(writer io.Writer, result Result) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func WriteText(writer io.Writer, result Result) {
	if result.Status == StatusBlocked || result.Status == StatusFailed || result.Status == StatusRecoveryRequired {
		fmt.Fprintf(writer, "%s: %s\n", result.Command, result.Status)
		for _, blocker := range result.Blockers {
			fmt.Fprintf(writer, "- %s: %s", blocker.Code, blocker.Message)
			if blocker.Path != "" {
				fmt.Fprintf(writer, " (%s)", blocker.Path)
			}
			fmt.Fprintln(writer)
			if blocker.Remediation != "" {
				fmt.Fprintf(writer, "  Remediation: %s\n", blocker.Remediation)
			}
		}
		return
	}
	mode := "Applied"
	if result.DryRun {
		mode = "Planned"
	}
	fmt.Fprintf(writer, "%s %d change(s); state %s", mode, len(result.Changes), result.StateBefore)
	if result.StateAfter != "" {
		fmt.Fprintf(writer, " -> %s", result.StateAfter)
	}
	fmt.Fprintln(writer, ".")
}
