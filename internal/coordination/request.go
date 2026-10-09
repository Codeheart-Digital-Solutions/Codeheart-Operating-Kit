// Package coordination implements the narrow cross-tool invocation helper: it launches one
// coordinator-prepared Claude CLI invocation, captures its original output and writes compact
// attempt, delivery and native-message records. It never interprets results, schedules work,
// retries, messages another tool or changes permissions.
package coordination

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// RequestSchemaVersion is the only accepted request schema version.
const RequestSchemaVersion = 1

// InlineReplyLimit caps the original reply copied into the native message argument file.
const InlineReplyLimit = 2000

var (
	identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)
	uuidPattern       = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	modelPattern      = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/\[\]-]{0,127}$`)
)

// permissionModes are passed through unchanged. Newer CLIs name the default mode "manual" and
// still accept "default".
var permissionModes = map[string]bool{
	"default":     true,
	"manual":      true,
	"acceptEdits": true,
	"auto":        true,
	"plan":        true,
	"dontAsk":     true,
}

var effortLevels = map[string]bool{"low": true, "medium": true, "high": true, "xhigh": true, "max": true}

var permissionPromptValues = map[string]bool{"none": true}

// Request is the complete coordinator-authored invocation request.
type Request struct {
	SchemaVersion     int         `json:"schema_version"`
	AssignmentID      string      `json:"assignment_id"`
	AttemptID         string      `json:"attempt_id"`
	ReturnDestination Destination `json:"return_destination"`
	WorkingDirectory  string      `json:"working_directory"`
	BriefFile         string      `json:"brief_file"`
	Executable        string      `json:"executable"`
	Model             string      `json:"model"`
	Effort            string      `json:"effort,omitempty"`
	Session           SessionSpec `json:"session"`
	Permissions       Permissions `json:"permissions"`
	AuthorityRefs     []string    `json:"authority_refs"`
	StateRoot         string      `json:"state_root"`
}

// Destination is the exact ordinary return chat prepared by the coordinator.
type Destination struct {
	ThreadID string `json:"thread_id"`
	HostID   string `json:"host_id"`
}

// SessionSpec selects a new CLI-owned session or the retained session to resume.
type SessionSpec struct {
	Mode string `json:"mode"`
	ID   string `json:"id,omitempty"`
}

// Permissions holds explicit, concrete permission and tool settings. Nothing is defaulted.
type Permissions struct {
	Mode              string   `json:"mode"`
	PermissionPrompts string   `json:"permission_prompts,omitempty"`
	Tools             []string `json:"tools,omitempty"`
	AllowedTools      []string `json:"allowed_tools,omitempty"`
	DisallowedTools   []string `json:"disallowed_tools,omitempty"`
	SettingsFile      string   `json:"settings_file,omitempty"`
	AddDirs           []string `json:"add_dirs,omitempty"`
	DisableMCPServers bool     `json:"disable_mcp_servers,omitempty"`
	DisableChrome     bool     `json:"disable_chrome,omitempty"`
}

// RequestError reports an invalid or unsupported request. It is returned before any launch.
type RequestError struct {
	Problems []string
}

func (e *RequestError) Error() string {
	return "invalid request: " + strings.Join(e.Problems, "; ")
}

// LoadRequest strictly decodes a request file, rejecting unknown fields and trailing data.
func LoadRequest(path string) (Request, []byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Request{}, nil, &RequestError{Problems: []string{fmt.Sprintf("cannot read request: %v", err)}}
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var request Request
	if err := decoder.Decode(&request); err != nil {
		return Request{}, raw, &RequestError{Problems: []string{fmt.Sprintf("malformed request JSON: %v", err)}}
	}
	// The file must hold exactly one JSON value followed only by whitespace.
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); err != io.EOF {
		return Request{}, raw, &RequestError{Problems: []string{"request contains trailing data after the JSON object"}}
	}
	if err := request.Validate(); err != nil {
		return request, raw, err
	}
	return request, raw, nil
}

// Validate checks that every required field is present, concrete and supported.
func (r Request) Validate() error {
	var problems []string
	add := func(format string, args ...any) { problems = append(problems, fmt.Sprintf(format, args...)) }

	if r.SchemaVersion != RequestSchemaVersion {
		add("schema_version must be %d", RequestSchemaVersion)
	}
	if !identifierPattern.MatchString(r.AssignmentID) {
		add("assignment_id must match %s", identifierPattern.String())
	}
	if !identifierPattern.MatchString(r.AttemptID) {
		add("attempt_id must match %s", identifierPattern.String())
	}
	if !plainValue(r.ReturnDestination.ThreadID) {
		add("return_destination.thread_id is required")
	}
	if !plainValue(r.ReturnDestination.HostID) {
		add("return_destination.host_id is required")
	}
	checkDir := func(field, value string) {
		if !filepath.IsAbs(value) {
			add("%s must be an absolute path", field)
			return
		}
		info, err := os.Stat(value)
		if err != nil || !info.IsDir() {
			add("%s must be an existing directory", field)
		}
	}
	checkFile := func(field, value string) {
		if !filepath.IsAbs(value) {
			add("%s must be an absolute path", field)
			return
		}
		info, err := os.Stat(value)
		if err != nil || !info.Mode().IsRegular() {
			add("%s must be an existing regular file", field)
		}
	}
	checkDir("working_directory", r.WorkingDirectory)
	checkFile("brief_file", r.BriefFile)
	checkFile("executable", r.Executable)
	if !modelPattern.MatchString(r.Model) {
		add("model is required and must be a single concrete model value")
	}
	if r.Effort != "" && !effortLevels[r.Effort] {
		add("effort %q is unsupported", r.Effort)
	}
	switch r.Session.Mode {
	case "new":
		if r.Session.ID != "" {
			add("session.id must be omitted for a new session; the helper allocates and records it")
		}
	case "resume":
		if !uuidPattern.MatchString(r.Session.ID) {
			add("session.id must be the retained session UUID for resume")
		}
	default:
		add("session.mode must be new or resume")
	}
	if r.Permissions.Mode == "bypassPermissions" {
		add("permissions.mode bypassPermissions is not supported by this helper")
	} else if !permissionModes[r.Permissions.Mode] {
		add("permissions.mode must be one of default, manual, acceptEdits, auto, plan, dontAsk")
	}
	if r.Permissions.PermissionPrompts != "" && !permissionPromptValues[r.Permissions.PermissionPrompts] {
		add("permissions.permission_prompts %q is unsupported", r.Permissions.PermissionPrompts)
	}
	for field, list := range map[string][]string{
		"permissions.tools":            r.Permissions.Tools,
		"permissions.allowed_tools":    r.Permissions.AllowedTools,
		"permissions.disallowed_tools": r.Permissions.DisallowedTools,
	} {
		for _, entry := range list {
			if !plainValue(entry) || strings.HasPrefix(entry, "-") {
				add("%s contains an empty, flag-like or control-character entry", field)
				break
			}
		}
	}
	if r.Permissions.SettingsFile != "" {
		checkFile("permissions.settings_file", r.Permissions.SettingsFile)
	}
	for _, dir := range r.Permissions.AddDirs {
		checkDir("permissions.add_dirs", dir)
	}
	if len(r.AuthorityRefs) == 0 {
		add("authority_refs must name at least one authority or source reference")
	}
	for _, ref := range r.AuthorityRefs {
		if !plainValue(ref) {
			add("authority_refs entries must be non-empty text")
			break
		}
	}
	if !filepath.IsAbs(r.StateRoot) {
		add("state_root must be an absolute path")
	}
	if len(problems) > 0 {
		return &RequestError{Problems: problems}
	}
	return nil
}

func plainValue(value string) bool {
	if strings.TrimSpace(value) == "" || len(value) > 512 {
		return false
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

// Args returns the fixed, documented CLI argument vector (excluding the executable). The brief
// is supplied on stdin, never as an argument.
func (r Request) Args(sessionID string) []string {
	args := []string{"--print", "--output-format", "stream-json", "--verbose", "--model", r.Model}
	if r.Session.Mode == "new" {
		args = append(args, "--session-id", sessionID)
	} else {
		args = append(args, "--resume", sessionID)
	}
	args = append(args, "--permission-mode", r.Permissions.Mode)
	if r.Permissions.PermissionPrompts != "" {
		args = append(args, "--permission-prompts", r.Permissions.PermissionPrompts)
	}
	if r.Effort != "" {
		args = append(args, "--effort", r.Effort)
	}
	if r.Permissions.SettingsFile != "" {
		args = append(args, "--settings", r.Permissions.SettingsFile)
	}
	if r.Permissions.DisableMCPServers {
		args = append(args, "--strict-mcp-config", "--mcp-config", `{"mcpServers":{}}`)
	}
	if r.Permissions.DisableChrome {
		args = append(args, "--no-chrome")
	}
	// Variadic list options come last, each followed by another option or the end of argv.
	if len(r.Permissions.Tools) > 0 {
		args = append(args, "--tools", strings.Join(r.Permissions.Tools, ","))
	}
	if len(r.Permissions.AllowedTools) > 0 {
		args = append(append(args, "--allowedTools"), r.Permissions.AllowedTools...)
	}
	if len(r.Permissions.DisallowedTools) > 0 {
		args = append(append(args, "--disallowedTools"), r.Permissions.DisallowedTools...)
	}
	if len(r.Permissions.AddDirs) > 0 {
		args = append(append(args, "--add-dir"), r.Permissions.AddDirs...)
	}
	return args
}

// AttemptDir is the deterministic evidence directory the coordinator knows before dispatch.
func AttemptDir(stateRoot, assignmentID, attemptID string) string {
	return filepath.Join(stateRoot, "attempts", assignmentID, attemptID)
}

// SessionLockPath is the shared host-local lock for one CLI session.
func SessionLockPath(stateRoot, sessionID string) string {
	return filepath.Join(stateRoot, "sessions", sessionID+".lock")
}
