package coordination

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

type harness struct {
	t         *testing.T
	root      string
	stateRoot string
	workDir   string
	recordDir string
	release   string
	exe       string
}

func newHarness(t *testing.T, mode string) *harness {
	t.Helper()
	root := filepath.Join(t.TempDir(), "coordination root with spaces")
	h := &harness{
		t:         t,
		root:      root,
		stateRoot: filepath.Join(root, "state root"),
		workDir:   filepath.Join(root, "work tree"),
		recordDir: filepath.Join(root, "fake records"),
		release:   filepath.Join(root, "release child"),
	}
	for _, dir := range []string{h.workDir, h.recordDir, filepath.Join(root, "bin dir")} {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	h.exe = filepath.Join(root, "bin dir", "fake claude")
	if runtime.GOOS == "windows" {
		h.exe += ".exe"
	}
	copyFile(t, testBinary(t), h.exe)
	fixtures, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv(envFakeCLI, mode)
	t.Setenv(envFakeFixtures, fixtures)
	t.Setenv(envFakeRecordDir, h.recordDir)
	t.Setenv(envFakeReleaseFile, h.release)
	return h
}

func testBinary(t *testing.T) string {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return exe
}

func copyFile(t *testing.T, from, to string) {
	t.Helper()
	in, err := os.Open(from)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o700)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(out, in); err != nil {
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
}

func (h *harness) request(assignment, attempt, brief string, mutate func(map[string]any)) string {
	h.t.Helper()
	briefPath := filepath.Join(h.root, assignment+"-"+attempt+" brief.md")
	if err := os.WriteFile(briefPath, []byte(brief), 0o600); err != nil {
		h.t.Fatal(err)
	}
	request := map[string]any{
		"schema_version":     1,
		"assignment_id":      assignment,
		"attempt_id":         attempt,
		"return_destination": map[string]any{"thread_id": "thread-" + assignment, "host_id": "local"},
		"working_directory":  h.workDir,
		"brief_file":         briefPath,
		"executable":         h.exe,
		"model":              "claude-example-model",
		"session":            map[string]any{"mode": "new"},
		"permissions": map[string]any{
			"mode":             "default",
			"allowed_tools":    []string{"Read", "Bash(git status *)"},
			"disallowed_tools": []string{"WebFetch"},
		},
		"authority_refs": []string{"docs/repo/plans/example_implementation_doc.md"},
		"state_root":     h.stateRoot,
	}
	if mutate != nil {
		mutate(request)
	}
	data, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		h.t.Fatal(err)
	}
	path := filepath.Join(h.root, assignment+"-"+attempt+" request.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		h.t.Fatal(err)
	}
	return path
}

func (h *harness) fakeRecords() []fakeRecord {
	h.t.Helper()
	entries, _ := os.ReadDir(h.recordDir)
	var records []fakeRecord
	for _, entry := range entries {
		var record fakeRecord
		if err := readJSON(filepath.Join(h.recordDir, entry.Name()), &record); err == nil {
			records = append(records, record)
		}
	}
	return records
}

func mustInvoke(t *testing.T, path string) Outcome {
	t.Helper()
	outcome, err := Invoke(Options{RequestPath: path, Interrupt: make(chan os.Signal)})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	return outcome
}

func readMessage(t *testing.T, outcome Outcome) MessageArguments {
	t.Helper()
	var message MessageArguments
	if err := readJSON(outcome.MessageFile, &message); err != nil {
		t.Fatal(err)
	}
	return message
}

func TestInvokeCapturesOriginalReplyWithExactArgsStdinAndRecipient(t *testing.T) {
	h := newHarness(t, "success")
	brief := "/goal Implement the plan.\nQuotes \"double\" 'single' $(not a command) `tick` ünïcode\n\nLast line"
	outcome := mustInvoke(t, h.request("assign-a", "attempt-1", brief, nil))
	if outcome.Status != StatusResponseCaptured {
		t.Fatalf("status = %s (%s)", outcome.Status, outcome.Detail)
	}
	records := h.fakeRecords()
	if len(records) != 1 {
		t.Fatalf("fake CLI ran %d times", len(records))
	}
	fake := records[0]
	if fake.Stdin != brief {
		t.Fatalf("stdin not preserved:\n%q\n%q", fake.Stdin, brief)
	}
	wantCwd, _ := filepath.EvalSymlinks(h.workDir)
	gotCwd, _ := filepath.EvalSymlinks(fake.Cwd)
	if gotCwd != wantCwd {
		t.Fatalf("cwd = %q, want %q", gotCwd, wantCwd)
	}
	want := []string{"--print", "--output-format", "stream-json", "--verbose", "--model", "claude-example-model",
		"--session-id", outcome.SessionID, "--permission-mode", "default",
		"--allowedTools", "Read", "Bash(git status *)", "--disallowedTools", "WebFetch"}
	if !reflect.DeepEqual(fake.Args, want) {
		t.Fatalf("argv = %#v\nwant %#v", fake.Args, want)
	}
	record, err := ReadAttempt(outcome.AttemptDir)
	if err != nil {
		t.Fatal(err)
	}
	if record.ActualSession != outcome.SessionID || record.RequestedSession != outcome.SessionID || record.CLIVersion != "2.1.286" || record.ActualModel != "claude-example-model" {
		t.Fatalf("record identity = %+v", record)
	}
	if record.Child == nil || record.Child.PID != fake.PID || record.Child.ExitCode == nil || *record.Child.ExitCode != 0 {
		t.Fatalf("child record = %+v, fake pid %d", record.Child, fake.PID)
	}
	if record.Response == nil || record.Response.JSONField != "result" || record.Usage == nil {
		t.Fatalf("response locator/usage missing: %+v", record)
	}
	message := readMessage(t, outcome)
	if message.ThreadID != "thread-assign-a" || message.HostID != "local" {
		t.Fatalf("message recipient = %+v", message)
	}
	if !strings.Contains(message.Prompt, "Original reply (unchanged):\necho:"+brief) || !strings.Contains(message.Prompt, "not task acceptance") {
		t.Fatalf("message prompt does not carry the original reply:\n%s", message.Prompt)
	}
	var delivery DeliveryRecord
	if err := readJSON(outcome.DeliveryFile, &delivery); err != nil || delivery.Status != DeliveryPending {
		t.Fatalf("delivery = %+v, %v", delivery, err)
	}
	if _, err := os.Stat(SessionLockPath(h.stateRoot, outcome.SessionID)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("session lock should be released after completion: %v", err)
	}
	encoded, _ := json.Marshal(outcome)
	if strings.Contains(string(encoded), "echo:") || strings.Contains(string(encoded), "Implement the plan") {
		t.Fatalf("normal output leaked reply or brief: %s", encoded)
	}
	if runtime.GOOS != "windows" {
		info, _ := os.Stat(outcome.AttemptDir)
		if info.Mode().Perm() != 0o700 {
			t.Fatalf("attempt dir mode = %v", info.Mode().Perm())
		}
	}
}

func TestInvokeResumeUsesRetainedSession(t *testing.T) {
	h := newHarness(t, "success")
	session := "0a1b2c3d-1111-4222-8333-444455556666"
	outcome := mustInvoke(t, h.request("assign-a", "attempt-2", "answer", func(r map[string]any) {
		r["session"] = map[string]any{"mode": "resume", "id": session}
	}))
	if outcome.Status != StatusResponseCaptured || outcome.SessionID != session {
		t.Fatalf("outcome = %+v", outcome)
	}
	args := h.fakeRecords()[0].Args
	if strings.Join(args, " ") == "" || !containsPair(args, "--resume", session) || containsPair(args, "--session-id", session) {
		t.Fatalf("resume argv = %#v", args)
	}
}

func containsPair(args []string, flag, value string) bool {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == flag && args[i+1] == value {
			return true
		}
	}
	return false
}

func TestInvokeLongReplyUsesExactFileReference(t *testing.T) {
	h := newHarness(t, "long")
	outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", nil))
	message := readMessage(t, outcome)
	if strings.Contains(message.Prompt, strings.Repeat("x", 100)) {
		t.Fatal("long reply was copied inline")
	}
	if !strings.Contains(message.Prompt, filepath.Join(outcome.AttemptDir, StdoutFile)) || !strings.Contains(message.Prompt, `JSON field "result"`) {
		t.Fatalf("long reply reference missing:\n%s", message.Prompt)
	}
}

func TestInvokeReportsFailuresTruthfully(t *testing.T) {
	cases := []struct {
		mode   string
		status string
	}{
		{"wrong_session", StatusSessionMismatch},
		{"empty_session", StatusSessionMismatch},
		{"malformed", StatusIncompleteOutput},
		{"empty", StatusIncompleteOutput},
		{"error", StatusCLIError},
		{"nonzero", StatusExecutionFailed},
		{"usage", StatusUnsupportedInvocation},
	}
	for _, tc := range cases {
		t.Run(tc.mode, func(t *testing.T) {
			h := newHarness(t, tc.mode)
			outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", nil))
			if outcome.Status != tc.status {
				t.Fatalf("status = %s (%s), want %s", outcome.Status, outcome.Detail, tc.status)
			}
			record, err := ReadAttempt(outcome.AttemptDir)
			if err != nil || record.Status != tc.status {
				t.Fatalf("record = %+v, %v", record, err)
			}
			message := readMessage(t, outcome)
			if !strings.Contains(message.Prompt, tc.status) {
				t.Fatalf("message hides status:\n%s", message.Prompt)
			}
			if _, err := os.Stat(filepath.Join(outcome.AttemptDir, StderrFile)); err != nil {
				t.Fatal("stderr evidence missing")
			}
		})
	}
}

func TestInvokeRecordsChildKilledBySignal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("signal termination shape is Unix-specific")
	}
	h := newHarness(t, "selfkill")
	outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", nil))
	if outcome.Status != StatusInterrupted {
		t.Fatalf("status = %s (%s)", outcome.Status, outcome.Detail)
	}
	data, _ := os.ReadFile(filepath.Join(outcome.AttemptDir, StdoutFile))
	if !strings.Contains(string(data), `"subtype":"init"`) {
		t.Fatal("partial output was not retained")
	}
}

func TestInvokeKeepsEarlierResultDenialsAndLocations(t *testing.T) {
	h := newHarness(t, "multi_result")
	outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", nil))
	if outcome.Status != StatusResponseCaptured || outcome.PermissionDenials != 1 {
		t.Fatalf("outcome = %+v", outcome)
	}
	record, err := ReadAttempt(outcome.AttemptDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(record.PermissionDenials) != 1 || record.PermissionDenials[0].ToolName != "Write" {
		t.Fatalf("earlier denial lost: %+v", record.PermissionDenials)
	}
	// Stream lines: 1 init, 2 first result with the substantive reply, 3 later completion.
	if record.Response == nil || record.Response.Line != 3 {
		t.Fatalf("final response locator = %+v", record.Response)
	}
	want := ResponseLocator{File: filepath.Join(outcome.AttemptDir, StdoutFile), Line: 2, JSONField: "result", Chars: len([]rune("Substantive handoff: the Write call was denied."))}
	if len(record.EarlierResponses) != 1 || record.EarlierResponses[0] != want {
		t.Fatalf("earlier responses = %+v, want %+v", record.EarlierResponses, want)
	}
	message := readMessage(t, outcome)
	if !strings.Contains(message.Prompt, "Permission denials: 1.") || !strings.Contains(message.Prompt, "line(s) 2, JSON field \"result\"") {
		t.Fatalf("message hides the earlier result or denial:\n%s", message.Prompt)
	}
	if strings.Contains(message.Prompt, "Substantive handoff") {
		t.Fatalf("message copied an earlier reply instead of locating it:\n%s", message.Prompt)
	}
	if !strings.Contains(message.Prompt, "Original reply (unchanged):\nBackground task finished.") {
		t.Fatalf("final reply not preserved:\n%s", message.Prompt)
	}
}

func TestAppendDenialsDeduplicatesRepeatedToolUse(t *testing.T) {
	got := appendDenials(nil, []Denial{{ToolName: "Write", ToolUseID: "a"}, {ToolName: "Bash"}})
	got = appendDenials(got, []Denial{{ToolName: "Write", ToolUseID: "a"}, {ToolName: "Bash"}, {ToolName: "Edit", ToolUseID: "b"}})
	want := []Denial{{ToolName: "Write", ToolUseID: "a"}, {ToolName: "Bash"}, {ToolName: "Bash"}, {ToolName: "Edit", ToolUseID: "b"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("denials = %+v, want %+v", got, want)
	}
}

func TestInvokePreservesRealDenialMetadata(t *testing.T) {
	h := newHarness(t, "denial")
	outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", nil))
	if outcome.Status != StatusResponseCaptured || outcome.PermissionDenials != 1 {
		t.Fatalf("outcome = %+v", outcome)
	}
	record, _ := ReadAttempt(outcome.AttemptDir)
	if len(record.PermissionDenials) != 1 || record.PermissionDenials[0].ToolName != "Write" {
		t.Fatalf("denials = %+v", record.PermissionDenials)
	}
	message := readMessage(t, outcome)
	if !strings.Contains(message.Prompt, "Permission denials: 1") || !strings.Contains(message.Prompt, "permission was denied") {
		t.Fatalf("denial not returned:\n%s", message.Prompt)
	}
}

func TestInvokeRejectsInvalidRequestsWithoutLaunching(t *testing.T) {
	h := newHarness(t, "success")
	cases := map[string]func(map[string]any){
		"unknown field":     func(r map[string]any) { r["extra_flags"] = "--dangerously-skip-permissions" },
		"bypass":            func(r map[string]any) { r["permissions"] = map[string]any{"mode": "bypassPermissions"} },
		"missing mode":      func(r map[string]any) { r["permissions"] = map[string]any{} },
		"relative workdir":  func(r map[string]any) { r["working_directory"] = "relative/dir" },
		"missing model":     func(r map[string]any) { delete(r, "model") },
		"shell model":       func(r map[string]any) { r["model"] = "opus; rm -rf /" },
		"resume without id": func(r map[string]any) { r["session"] = map[string]any{"mode": "resume"} },
		"new with id": func(r map[string]any) {
			r["session"] = map[string]any{"mode": "new", "id": "0a1b2c3d-1111-4222-8333-444455556666"}
		},
		"no authority": func(r map[string]any) { r["authority_refs"] = []string{} },
		"no recipient": func(r map[string]any) { r["return_destination"] = map[string]any{"thread_id": "", "host_id": "local"} },
		"flag-like tool": func(r map[string]any) {
			r["permissions"] = map[string]any{"mode": "default", "allowed_tools": []string{"--dangerously-skip-permissions"}}
		},
		"bad prompt setting": func(r map[string]any) {
			r["permissions"] = map[string]any{"mode": "default", "permission_prompts": "always"}
		},
		"wrong schema": func(r map[string]any) { r["schema_version"] = 2 },
		"traversal id": func(r map[string]any) { r["attempt_id"] = "../escape" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Invoke(Options{RequestPath: h.request("assign-x", "attempt-1", "brief", mutate)})
			var requestErr *RequestError
			if !errors.As(err, &requestErr) {
				t.Fatalf("err = %v, want RequestError", err)
			}
		})
	}
	trailing := filepath.Join(h.root, "trailing.json")
	data, _ := os.ReadFile(h.request("assign-x", "attempt-1", "brief", nil))
	_ = os.WriteFile(trailing, append(data, []byte(`{"second":true}`)...), 0o600)
	if _, err := Invoke(Options{RequestPath: trailing}); err == nil {
		t.Fatal("trailing data accepted")
	}
	if len(h.fakeRecords()) != 0 {
		t.Fatal("invalid request launched the CLI")
	}
	if _, err := os.Stat(filepath.Join(h.stateRoot, "attempts")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("invalid request created attempt state")
	}
}

func TestReplayedAttemptIsRefusedAndOriginalPreserved(t *testing.T) {
	h := newHarness(t, "success")
	path := h.request("assign-a", "attempt-1", "brief", nil)
	first := mustInvoke(t, path)
	before, _ := os.ReadFile(first.AttemptFile)
	second := mustInvoke(t, path)
	if second.Status != "attempt_exists" {
		t.Fatalf("replay status = %s", second.Status)
	}
	after, _ := os.ReadFile(first.AttemptFile)
	if string(before) != string(after) || len(h.fakeRecords()) != 1 {
		t.Fatal("replay changed the original attempt or launched again")
	}
}

func TestReplayRefusalDoesNotResendOrResetEarlierDelivery(t *testing.T) {
	h := newHarness(t, "success")
	path := h.request("assign-a", "attempt-1", "brief", nil)
	first := mustInvoke(t, path)
	if _, err := RecordDelivery(DeliveryInput{AttemptDir: first.AttemptDir, Status: DeliverySent, ReportedThreadID: "thread-assign-a"}); err != nil {
		t.Fatal(err)
	}
	messageBefore, _ := os.ReadFile(first.MessageFile)
	deliveryBefore, _ := os.ReadFile(first.DeliveryFile)
	second := mustInvoke(t, path)
	if second.Status != StatusAttemptExists || second.MessageFile != "" || second.DeliveryFile != "" {
		t.Fatalf("replay outcome = %+v", second)
	}
	messageAfter, _ := os.ReadFile(first.MessageFile)
	deliveryAfter, _ := os.ReadFile(first.DeliveryFile)
	if string(messageBefore) != string(messageAfter) || string(deliveryBefore) != string(deliveryAfter) {
		t.Fatal("replay regenerated the message or reset the delivery receipt")
	}
	fallback := second.FallbackMessage
	if fallback == nil || fallback.ThreadID != "thread-assign-a" || fallback.HostID != "local" {
		t.Fatalf("refusal notice not addressed to the prepared destination: %+v", fallback)
	}
	if strings.Contains(fallback.Prompt, "echo:") || !strings.Contains(fallback.Prompt, "not the earlier result") {
		t.Fatalf("refusal notice resends or misdescribes the old result:\n%s", fallback.Prompt)
	}
}

func assertNotificationFiles(t *testing.T, outcome Outcome, status string) MessageArguments {
	t.Helper()
	if outcome.Status != status {
		t.Fatalf("status = %s (%s), want %s", outcome.Status, outcome.Detail, status)
	}
	if outcome.MessageFile != filepath.Join(outcome.AttemptDir, MessageFile) || outcome.DeliveryFile != filepath.Join(outcome.AttemptDir, DeliveryFile) || len(outcome.PersistenceErrors) != 0 {
		t.Fatalf("notification files not reported: %+v", outcome)
	}
	message := readMessage(t, outcome)
	if !strings.Contains(message.Prompt, status) {
		t.Fatalf("message hides status %s:\n%s", status, message.Prompt)
	}
	var delivery DeliveryRecord
	if err := readJSON(outcome.DeliveryFile, &delivery); err != nil || delivery.Status != DeliveryPending {
		t.Fatalf("delivery = %+v, %v", delivery, err)
	}
	record, err := ReadAttempt(outcome.AttemptDir)
	if err != nil || record.Status != status {
		t.Fatalf("record = %+v, %v", record, err)
	}
	return message
}

func TestLaunchFailureWritesNotificationAndReleasesOwnLock(t *testing.T) {
	h := newHarness(t, "success")
	notExecutable := filepath.Join(h.root, "not executable.txt")
	writeFile(t, notExecutable, "plain text")
	outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", func(r map[string]any) { r["executable"] = notExecutable }))
	message := assertNotificationFiles(t, outcome, StatusLaunchFailed)
	if !strings.Contains(message.Prompt, "No CLI process was recorded as started") || strings.Contains(message.Prompt, StdoutFile) {
		t.Fatalf("launch failure message:\n%s", message.Prompt)
	}
	if _, err := os.Stat(SessionLockPath(h.stateRoot, outcome.SessionID)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("launch failure left its own lock")
	}
}

func TestInvalidRequestFallbackUsesOnlyThePreparedDestination(t *testing.T) {
	h := newHarness(t, "success")
	outcome, err := Invoke(Options{RequestPath: h.request("assign-a", "attempt-1", "brief", func(r map[string]any) { r["model"] = "" })})
	if err == nil || outcome.Status != StatusInvalidRequest || outcome.FallbackMessage == nil || outcome.FallbackMessage.ThreadID != "thread-assign-a" || outcome.MessageFile != "" {
		t.Fatalf("outcome = %+v, err = %v", outcome, err)
	}
	outcome, _ = Invoke(Options{RequestPath: h.request("assign-b", "attempt-1", "brief", func(r map[string]any) {
		r["model"] = ""
		r["return_destination"] = map[string]any{"thread_id": "", "host_id": "local"}
	})})
	if outcome.FallbackMessage != nil {
		t.Fatalf("fallback invented a destination: %+v", outcome.FallbackMessage)
	}
}

func TestPersistenceFailureIsReportedNotClaimed(t *testing.T) {
	h := newHarness(t, "success")
	// Occupy the message path with a directory so the atomic rename fails.
	blocked := filepath.Join(AttemptDir(h.stateRoot, "assign-a", "attempt-1"), MessageFile)
	t.Setenv(envFakeMkdir, blocked)
	outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", nil))
	if outcome.MessageFile != "" || len(outcome.PersistenceErrors) == 0 || outcome.FallbackMessage == nil {
		t.Fatalf("persistence failure hidden: %+v", outcome)
	}
	if !strings.Contains(outcome.FallbackMessage.Prompt, "Original reply (unchanged):\necho:brief") || outcome.FallbackMessage.ThreadID != "thread-assign-a" {
		t.Fatalf("fallback does not carry the composed message:\n%+v", outcome.FallbackMessage)
	}
	if outcome.DeliveryFile == "" {
		t.Fatal("pending delivery should still be recorded")
	}
}

func TestDistinctConcurrentAssignmentsStayIsolated(t *testing.T) {
	h := newHarness(t, "slow")
	paths := []string{
		h.request("assign-a", "attempt-1", "brief A", nil),
		h.request("assign-b", "attempt-1", "brief B", nil),
	}
	outcomes := make([]Outcome, 2)
	var wg sync.WaitGroup
	for i, path := range paths {
		wg.Add(1)
		go func(i int, path string) {
			defer wg.Done()
			outcomes[i], _ = Invoke(Options{RequestPath: path, Interrupt: make(chan os.Signal)})
		}(i, path)
	}
	waitFor(t, func() bool { return len(h.fakeRecords()) == 2 })
	writeFile(t, h.release, "go")
	wg.Wait()
	if outcomes[0].SessionID == outcomes[1].SessionID || outcomes[0].AttemptDir == outcomes[1].AttemptDir {
		t.Fatalf("assignments shared identity: %+v", outcomes)
	}
	for i, outcome := range outcomes {
		if outcome.Status != StatusResponseCaptured {
			t.Fatalf("outcome %d = %+v", i, outcome)
		}
		message := readMessage(t, outcome)
		want := []string{"thread-assign-a", "thread-assign-b"}[i]
		brief := []string{"echo:brief A", "echo:brief B"}[i]
		if message.ThreadID != want || !strings.Contains(message.Prompt, brief) {
			t.Fatalf("message %d crossed assignments: %+v", i, message)
		}
	}
}

func TestSameSessionCollisionLaunchesNoSecondWriter(t *testing.T) {
	h := newHarness(t, "slow")
	session := "0a1b2c3d-1111-4222-8333-444455556666"
	resume := func(r map[string]any) { r["session"] = map[string]any{"mode": "resume", "id": session} }
	pathA := h.request("assign-a", "attempt-1", "first", resume)
	pathB := h.request("assign-a", "attempt-2", "second", resume)
	done := make(chan Outcome, 1)
	go func() {
		outcome, _ := Invoke(Options{RequestPath: pathA, Interrupt: make(chan os.Signal)})
		done <- outcome
	}()
	waitFor(t, func() bool { return len(h.fakeRecords()) == 1 })
	second := mustInvoke(t, pathB)
	if second.Status != StatusSessionLocked {
		t.Fatalf("second status = %s", second.Status)
	}
	record, _ := ReadAttempt(second.AttemptDir)
	if record.LockOwner == nil || record.LockOwner.AttemptID != "attempt-1" {
		t.Fatalf("lock owner not recorded: %+v", record.LockOwner)
	}
	if len(h.fakeRecords()) != 1 {
		t.Fatal("second writer launched")
	}
	message := assertNotificationFiles(t, second, StatusSessionLocked)
	if !strings.Contains(message.Prompt, "Session held by attempt attempt-1") {
		t.Fatalf("locked message does not name the holder:\n%s", message.Prompt)
	}
	if _, err := os.Stat(SessionLockPath(h.stateRoot, session)); err != nil {
		t.Fatal("refused attempt released another attempt's lock")
	}
	writeFile(t, h.release, "go")
	first := <-done
	if first.Status != StatusResponseCaptured {
		t.Fatalf("first = %+v", first)
	}
	if _, err := os.Stat(SessionLockPath(h.stateRoot, session)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("lock retained after clean completion")
	}
}

func TestHelperInterruptRetainsLockAndEvidence(t *testing.T) {
	h := newHarness(t, "slow")
	interrupt := make(chan os.Signal, 1)
	done := make(chan Outcome, 1)
	path := h.request("assign-a", "attempt-1", "brief", nil)
	go func() {
		outcome, _ := Invoke(Options{RequestPath: path, Interrupt: interrupt})
		done <- outcome
	}()
	waitFor(t, func() bool { return len(h.fakeRecords()) == 1 })
	interrupt <- os.Interrupt
	outcome := <-done
	if outcome.Status != StatusHelperInterrupted || !strings.Contains(outcome.Detail, "child still running") {
		t.Fatalf("outcome = %+v", outcome)
	}
	if _, err := os.Stat(SessionLockPath(h.stateRoot, outcome.SessionID)); err != nil {
		t.Fatal("lock not retained after helper interruption")
	}
	assertNotificationFiles(t, outcome, StatusHelperInterrupted)
	_, err := ReleaseLock(ReleaseInput{StateRoot: h.stateRoot, SessionID: outcome.SessionID, AssignmentID: "assign-a", AttemptID: "attempt-1", ManualVerification: "statement cannot override a live launcher"})
	var refusal *ReleaseRefusal
	if !errors.As(err, &refusal) || refusal.Reason != "launcher_alive_or_unknown" {
		t.Fatalf("release while launcher/child alive = %v", err)
	}
	writeFile(t, h.release, "go")
	pid := h.fakeRecords()[0].PID
	waitFor(t, func() bool { alive, _ := processAlive(pid); return !alive })
}

func TestKilledHelperLeavesChildOutputAndGuardedRelease(t *testing.T) {
	h := newHarness(t, "slow")
	path := h.request("assign-a", "attempt-1", "brief", nil)
	helper := startHelperProcess(t, path, false)
	attemptDir := AttemptDir(h.stateRoot, "assign-a", "attempt-1")
	var record AttemptRecord
	waitFor(t, func() bool {
		var err error
		record, err = ReadAttempt(attemptDir)
		return err == nil && record.Status == StatusRunning && record.Child != nil
	})
	stdoutPath := filepath.Join(attemptDir, StdoutFile)
	waitFor(t, func() bool { data, _ := os.ReadFile(stdoutPath); return strings.Count(string(data), "progress") >= 3 })
	if err := helper.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = helper.Wait()
	before, _ := os.ReadFile(stdoutPath)
	waitFor(t, func() bool { data, _ := os.ReadFile(stdoutPath); return len(data) > len(before) })
	if alive, _ := processAlive(record.Child.PID); !alive {
		t.Fatal("child should survive helper termination in this fixture")
	}
	session := record.RequestedSession
	_, err := ReleaseLock(ReleaseInput{StateRoot: h.stateRoot, SessionID: session, AssignmentID: "assign-a", AttemptID: "attempt-1", ManualVerification: "statement cannot override a live child"})
	var refusal *ReleaseRefusal
	if !errors.As(err, &refusal) || refusal.Reason != "child_alive_or_unknown" {
		t.Fatalf("release while child alive = %v", err)
	}
	_, err = ReleaseLock(ReleaseInput{StateRoot: h.stateRoot, SessionID: session, AssignmentID: "assign-a", AttemptID: "attempt-2"})
	if !errors.As(err, &refusal) || refusal.Reason != "owner_mismatch" {
		t.Fatalf("release for another attempt = %v", err)
	}
	writeFile(t, h.release, "go")
	waitFor(t, func() bool { alive, _ := processAlive(record.Child.PID); return !alive })
	after, _ := os.ReadFile(stdoutPath)
	if !strings.HasPrefix(string(after), string(before)) || !strings.Contains(string(after), `"type":"result"`) {
		t.Fatal("child output was not retained intact")
	}
	if current, _ := ReadAttempt(attemptDir); current.Status != StatusRunning {
		t.Fatalf("record should remain running/uncertain, got %s", current.Status)
	}
	released, err := ReleaseLock(ReleaseInput{StateRoot: h.stateRoot, SessionID: session, AssignmentID: "assign-a", AttemptID: "attempt-1"})
	if err != nil || released.LockRelease == nil {
		t.Fatalf("release after verified exit = %v", err)
	}
	if _, err := os.Stat(SessionLockPath(h.stateRoot, session)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("lock still present")
	}
}

func TestInterruptionBeforeChildIdentityRequiresManualVerification(t *testing.T) {
	h := newHarness(t, "slow")
	path := h.request("assign-a", "attempt-1", "brief", nil)
	helper := startHelperProcess(t, path, true)
	_ = helper.Wait()
	attemptDir := AttemptDir(h.stateRoot, "assign-a", "attempt-1")
	record, err := ReadAttempt(attemptDir)
	if err != nil || record.Status != StatusLaunching || record.Child != nil {
		t.Fatalf("record = %+v, %v", record, err)
	}
	session := record.RequestedSession
	_, err = ReleaseLock(ReleaseInput{StateRoot: h.stateRoot, SessionID: session, AssignmentID: "assign-a", AttemptID: "attempt-1"})
	var refusal *ReleaseRefusal
	if !errors.As(err, &refusal) || refusal.Reason != "child_identity_unknown" {
		t.Fatalf("release with unknown child = %v", err)
	}
	waitFor(t, func() bool { return len(h.fakeRecords()) == 1 })
	writeFile(t, h.release, "go")
	pid := h.fakeRecords()[0].PID
	waitFor(t, func() bool { alive, _ := processAlive(pid); return !alive })
	released, err := ReleaseLock(ReleaseInput{StateRoot: h.stateRoot, SessionID: session, AssignmentID: "assign-a", AttemptID: "attempt-1", ManualVerification: "Process list checked; no CLI process for this session remains."})
	if err != nil || released.LockRelease == nil || released.LockRelease.ManualVerification == "" {
		t.Fatalf("manual release = %+v, %v", released.LockRelease, err)
	}
}

func TestRecordDeliveryStates(t *testing.T) {
	newAttempt := func(t *testing.T) string {
		h := newHarness(t, "success")
		return mustInvoke(t, h.request("assign-a", "attempt-1", "brief", nil)).AttemptDir
	}
	t.Run("sent with matching reported recipient", func(t *testing.T) {
		dir := newAttempt(t)
		record, err := RecordDelivery(DeliveryInput{AttemptDir: dir, Status: DeliverySent, ReportedThreadID: "thread-assign-a", ReportedHostID: "local"})
		if err != nil || record.Status != DeliverySent || !record.RecipientConfirmed {
			t.Fatalf("record = %+v, %v", record, err)
		}
		if _, err := RecordDelivery(DeliveryInput{AttemptDir: dir, Status: DeliveryRejected}); err == nil {
			t.Fatal("second delivery update accepted")
		}
	})
	t.Run("omitted recipient is not confirmation", func(t *testing.T) {
		record, err := RecordDelivery(DeliveryInput{AttemptDir: newAttempt(t), Status: DeliverySent})
		if err != nil || record.RecipientConfirmed {
			t.Fatalf("record = %+v, %v", record, err)
		}
	})
	t.Run("mismatched recipient is uncertain", func(t *testing.T) {
		record, err := RecordDelivery(DeliveryInput{AttemptDir: newAttempt(t), Status: DeliverySent, ReportedThreadID: "thread-other"})
		if err != nil || record.Status != DeliveryUncertain || !record.RecipientMismatch || record.RecipientConfirmed {
			t.Fatalf("record = %+v, %v", record, err)
		}
	})
	t.Run("injected rejection is retained and found by known path", func(t *testing.T) {
		dir := newAttempt(t)
		receipt, _ := filepath.Abs(filepath.Join("testdata", "native-send-rejection.txt"))
		if _, err := RecordDelivery(DeliveryInput{AttemptDir: dir, Status: DeliveryRejected, ReceiptFile: receipt}); err != nil {
			t.Fatal(err)
		}
		var delivery DeliveryRecord
		if err := readJSON(filepath.Join(dir, DeliveryFile), &delivery); err != nil || delivery.Status != DeliveryRejected || !strings.Contains(delivery.Receipt, "Message rejected") {
			t.Fatalf("delivery = %+v, %v", delivery, err)
		}
		record, err := ReadAttempt(dir)
		if err != nil || record.Status != StatusResponseCaptured || record.Response == nil {
			t.Fatalf("retained result not found: %+v, %v", record, err)
		}
	})
}

func startHelperProcess(t *testing.T, requestPath string, exitAfterStart bool) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(testBinary(t))
	cmd.Env = append(os.Environ(), envHelperRequest+"="+requestPath)
	cmd.Env = removeEnv(cmd.Env, envFakeCLI)
	cmd.Env = append(cmd.Env, "CH_COORD_CHILD_MODE="+os.Getenv(envFakeCLI))
	if exitAfterStart {
		cmd.Env = append(cmd.Env, envHelperExitEarly+"=1")
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	return cmd
}

func removeEnv(env []string, name string) []string {
	out := env[:0:0]
	for _, entry := range env {
		if !strings.HasPrefix(entry, name+"=") {
			out = append(out, entry)
		}
	}
	return out
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition not reached before deadline")
}

func TestUnreadableLockStaysRefusedEvenWithStatement(t *testing.T) {
	stateRoot := filepath.Join(t.TempDir(), "state")
	session := "0a1b2c3d-1111-4222-8333-444455556666"
	lockPath := SessionLockPath(stateRoot, session)
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o700); err != nil {
		t.Fatal(err)
	}
	writeFile(t, lockPath, "")
	_, err := ReleaseLock(ReleaseInput{StateRoot: stateRoot, SessionID: session, AssignmentID: "assign-a", AttemptID: "attempt-1", ManualVerification: "statement cannot establish ownership"})
	var refusal *ReleaseRefusal
	if !errors.As(err, &refusal) || refusal.Reason != "owner_unknown" {
		t.Fatalf("release of unreadable lock = %v", err)
	}
	if _, err := os.Stat(lockPath); err != nil {
		t.Fatal("unreadable lock was removed")
	}
}

func TestInvokePassesOptionalSettingsExactly(t *testing.T) {
	h := newHarness(t, "success")
	extra := filepath.Join(h.root, "extra dir")
	if err := os.MkdirAll(extra, 0o700); err != nil {
		t.Fatal(err)
	}
	settings := filepath.Join(h.root, "role profile.json")
	writeFile(t, settings, "{}")
	outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", func(r map[string]any) {
		r["effort"] = "high"
		r["permissions"] = map[string]any{
			"mode":                "auto",
			"permission_prompts":  "none",
			"tools":               []string{"Read", "Write"},
			"allowed_tools":       []string{"Read"},
			"settings_file":       settings,
			"add_dirs":            []string{extra},
			"disable_mcp_servers": true,
			"disable_chrome":      true,
		}
	}))
	if outcome.Status != StatusResponseCaptured {
		t.Fatalf("outcome = %+v", outcome)
	}
	want := []string{"--print", "--output-format", "stream-json", "--verbose", "--model", "claude-example-model",
		"--session-id", outcome.SessionID, "--permission-mode", "auto", "--permission-prompts", "none",
		"--effort", "high", "--settings", settings, "--strict-mcp-config", "--mcp-config", `{"mcpServers":{}}`,
		"--no-chrome", "--tools", "Read,Write", "--allowedTools", "Read", "--add-dir", extra}
	if got := h.fakeRecords()[0].Args; !reflect.DeepEqual(got, want) {
		t.Fatalf("argv = %#v\nwant %#v", got, want)
	}
}

func TestRequestMustBeOneJSONValue(t *testing.T) {
	h := newHarness(t, "success")
	valid, err := os.ReadFile(h.request("assign-a", "attempt-1", "brief", nil))
	if err != nil {
		t.Fatal(err)
	}
	for name, suffix := range map[string]string{
		"closing bracket": "]",
		"closing brace":   "}",
		"second value":    "\n{\"schema_version\":1}",
		"junk":            " trailing-junk",
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(h.root, "trailing request.json")
			writeFile(t, path, string(valid)+suffix)
			outcome, err := Invoke(Options{RequestPath: path})
			var requestErr *RequestError
			if !errors.As(err, &requestErr) || outcome.Status != StatusInvalidRequest {
				t.Fatalf("outcome = %+v, err = %v", outcome, err)
			}
		})
	}
	if len(h.fakeRecords()) != 0 {
		t.Fatal("trailing data launched the CLI")
	}
	if _, err := os.Stat(filepath.Join(h.stateRoot, "attempts")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("trailing data created attempt state")
	}
	path := filepath.Join(h.root, "whitespace request.json")
	writeFile(t, path, string(valid)+"\n\t  \r\n")
	if outcome := mustInvoke(t, path); outcome.Status != StatusResponseCaptured {
		t.Fatalf("trailing whitespace rejected: %+v", outcome)
	}
}

func TestFinalReplyMustBeAJSONString(t *testing.T) {
	for _, mode := range []string{"null_reply", "number_reply", "missing_reply"} {
		t.Run(mode, func(t *testing.T) {
			h := newHarness(t, mode)
			outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", nil))
			if outcome.Status != StatusIncompleteOutput {
				t.Fatalf("status = %s (%s)", outcome.Status, outcome.Detail)
			}
			record, _ := ReadAttempt(outcome.AttemptDir)
			if record.Response != nil {
				t.Fatalf("non-string reply recorded as captured: %+v", record.Response)
			}
			if message := readMessage(t, outcome); strings.Contains(message.Prompt, "Original reply (unchanged)") {
				t.Fatalf("message claims an original reply:\n%s", message.Prompt)
			}
		})
	}
	t.Run("empty string", func(t *testing.T) {
		h := newHarness(t, "empty_reply")
		outcome := mustInvoke(t, h.request("assign-a", "attempt-1", "brief", nil))
		record, _ := ReadAttempt(outcome.AttemptDir)
		if outcome.Status != StatusResponseCaptured || record.Response == nil || record.Response.Chars != 0 {
			t.Fatalf("empty string reply = %+v, %+v", outcome, record.Response)
		}
	})
}
