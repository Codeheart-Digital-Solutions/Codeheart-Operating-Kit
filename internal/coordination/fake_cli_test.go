package coordination

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The test binary doubles as a fake Claude CLI and as a separately killable helper process.
// Fake output is built from sanitized real CLI shapes in testdata.
const (
	envFakeCLI         = "CH_COORD_FAKE_CLI"
	envFakeFixtures    = "CH_COORD_FAKE_FIXTURES"
	envFakeRecordDir   = "CH_COORD_FAKE_RECORD_DIR"
	envFakeReleaseFile = "CH_COORD_FAKE_RELEASE_FILE"
	envHelperRequest   = "CH_COORD_HELPER_REQUEST"
	envHelperExitEarly = "CH_COORD_HELPER_EXIT_AFTER_START"
	envFakeMkdir       = "CH_COORD_FAKE_MKDIR"
)

func TestMain(m *testing.M) {
	if mode := os.Getenv(envFakeCLI); mode != "" {
		os.Exit(runFakeCLI(mode))
	}
	if request := os.Getenv(envHelperRequest); request != "" {
		// The helper itself is not a fake CLI; its child inherits the fake mode.
		os.Setenv(envFakeCLI, os.Getenv("CH_COORD_CHILD_MODE"))
		if os.Getenv(envHelperExitEarly) != "" {
			afterStartHook = func() { os.Exit(9) }
		}
		outcome, err := Invoke(Options{RequestPath: request})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		_ = json.NewEncoder(os.Stdout).Encode(outcome)
		os.Exit(0)
	}
	os.Exit(m.Run())
}

type fakeRecord struct {
	Args  []string `json:"args"`
	Cwd   string   `json:"cwd"`
	Stdin string   `json:"stdin"`
	PID   int      `json:"pid"`
}

func runFakeCLI(mode string) int {
	args := os.Args[1:]
	stdin, _ := io.ReadAll(os.Stdin)
	cwd, _ := os.Getwd()
	if dir := os.Getenv(envFakeRecordDir); dir != "" {
		data, _ := json.Marshal(fakeRecord{Args: args, Cwd: cwd, Stdin: string(stdin), PID: os.Getpid()})
		_ = os.WriteFile(filepath.Join(dir, fmt.Sprintf("fake-%d.json", os.Getpid())), data, 0o600)
	}
	session := ""
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "--session-id" || args[i] == "--resume" {
			session = args[i+1]
		}
	}
	fixtures := os.Getenv(envFakeFixtures)
	line := func(name string, edit func(map[string]any)) {
		data, err := os.ReadFile(filepath.Join(fixtures, name))
		if err != nil {
			panic(err)
		}
		var fields map[string]any
		if err := json.Unmarshal(data, &fields); err != nil {
			panic(err)
		}
		fields["session_id"] = session
		if edit != nil {
			edit(fields)
		}
		out, _ := json.Marshal(fields)
		os.Stdout.Write(append(out, '\n'))
	}
	if dir := os.Getenv(envFakeMkdir); dir != "" {
		_ = os.MkdirAll(dir, 0o700)
	}
	reply := "echo:" + string(stdin)
	switch mode {
	case "success":
		line("claude-stream-init.json", nil)
		line("claude-result-success.json", func(f map[string]any) { f["result"] = reply })
	case "long":
		line("claude-stream-init.json", nil)
		line("claude-result-success.json", func(f map[string]any) { f["result"] = strings.Repeat("x", InlineReplyLimit+500) })
	case "null_reply", "number_reply", "missing_reply", "empty_reply":
		line("claude-stream-init.json", nil)
		line("claude-result-success.json", func(f map[string]any) {
			switch mode {
			case "null_reply":
				f["result"] = nil
			case "number_reply":
				f["result"] = 42
			case "missing_reply":
				delete(f, "result")
			default:
				f["result"] = ""
			}
		})
	case "denial":
		line("claude-stream-init.json", nil)
		line("claude-result-denial.json", nil)
	case "multi_result":
		// Observed shape: the substantive reply and its denial arrive in a first result record,
		// then a later background completion emits another result without denials.
		line("claude-stream-init.json", nil)
		line("claude-result-denial.json", func(f map[string]any) { f["result"] = "Substantive handoff: the Write call was denied." })
		line("claude-result-success.json", func(f map[string]any) { f["result"] = "Background task finished." })
	case "wrong_session":
		line("claude-stream-init.json", func(f map[string]any) { f["session_id"] = "99999999-9999-4999-8999-999999999999" })
		line("claude-result-success.json", func(f map[string]any) { f["session_id"] = "99999999-9999-4999-8999-999999999999" })
	case "empty_session":
		line("claude-result-success.json", func(f map[string]any) { f["session_id"] = "" })
	case "malformed":
		os.Stdout.WriteString("this is not json\n")
	case "empty":
	case "error":
		line("claude-stream-init.json", nil)
		line("claude-result-success.json", func(f map[string]any) {
			f["subtype"] = "error_during_execution"
			f["is_error"] = true
			delete(f, "result")
		})
		return 1
	case "nonzero":
		os.Stderr.WriteString("fatal: simulated execution failure\n")
		return 3
	case "usage":
		data, _ := os.ReadFile(filepath.Join(fixtures, "claude-usage-error.txt"))
		os.Stderr.Write(data)
		return 1
	case "slow":
		line("claude-stream-init.json", nil)
		release := os.Getenv(envFakeReleaseFile)
		for i := 0; ; i++ {
			if _, err := os.Stat(release); err == nil {
				break
			}
			fmt.Fprintf(os.Stdout, "{\"type\":\"assistant\",\"progress\":%d}\n", i)
			time.Sleep(20 * time.Millisecond)
		}
		line("claude-result-success.json", func(f map[string]any) { f["result"] = reply })
	case "selfkill":
		line("claude-stream-init.json", nil)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Kill()
		time.Sleep(time.Second)
	default:
		fmt.Fprintf(os.Stderr, "unknown fake mode %s\n", mode)
		return 4
	}
	return 0
}
