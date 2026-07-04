package commands

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/lockfile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/yamlmini"
)

func TestInspectFolderModes(t *testing.T) {
	root := t.TempDir()
	if got := InspectFolder(filepath.Join(root, "new"))["mode"]; got != "new-folder-setup" {
		t.Fatalf("missing folder mode = %v", got)
	}
	if got := InspectFolder(root)["mode"]; got != "new-folder-setup" {
		t.Fatalf("empty folder mode = %v", got)
	}
	if err := os.WriteFile(filepath.Join(root, "pyproject.toml"), []byte("[project]\nname = 'x'\n"), 0o644); err != nil {
		t.Fatalf("write marker: %v", err)
	}
	technical := InspectFolder(root)
	if technical["mode"] != "existing-technical-project-adoption" {
		t.Fatalf("technical folder mode = %v", technical["mode"])
	}
	markers := technical["markers"].([]any)
	if len(markers) != 1 || markers[0] != "pyproject.toml" {
		t.Fatalf("technical markers = %#v", markers)
	}
}

func TestRunInspectJSON(t *testing.T) {
	var stdout bytes.Buffer
	code := RunInspect([]string{"--json", "."}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunInspect exit = %d", code)
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse inspect JSON: %v\n%s", err, stdout.String())
	}
	if payload["mode"] == "" {
		t.Fatalf("inspect JSON missing mode: %#v", payload)
	}
}

func TestInitWritesStandardSurfaces(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	code := RunInit([]string{root, "--project-name", "Example-Automation", "--purpose", "company-automation"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunInit exit = %d; stdout: %s", code, stdout.String())
	}
	if !strings.Contains(stdout.String(), "Operating Kit initialized.") {
		t.Fatalf("init text output missing success line: %q", stdout.String())
	}
	for _, relative := range []string{
		".codeheart/kit",
		".codeheart/kit/README.md",
		".codeheart/kit.lock.yaml",
		".codeheart/kit.config.yaml",
		".codeheart/user/README.md",
		".codeheart/user/examples/preferences.yaml",
		"AGENTS.md",
		"docs/repo/README.md",
		"docs/repo/plans/plan-register.md",
		"docs/repo/plans/coordination-sync-pending.md",
		"docs/agent-memory/README.md",
		"docs/agent-memory/goal-register.md",
		"docs/agent-memory/session-ledger.md",
		"docs/agent-memory/untriaged-sessions.md",
		".gitignore",
	} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(relative))); err != nil {
			t.Fatalf("expected %s to exist: %v", relative, err)
		}
	}
	config := readYAMLFile(t, filepath.Join(root, ".codeheart", "kit.config.yaml"))
	if config["setup_purpose"] != "company-automation" {
		t.Fatalf("setup_purpose = %#v", config["setup_purpose"])
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	if set := mapKeys(mapValue(lock["native_capabilities"])); !sameStringSet(set, []string{"browser", "documents", "pdf", "presentations", "spreadsheets"}) {
		t.Fatalf("native capabilities = %#v", set)
	}
}

func TestInitJSONOutput(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	code := RunInit([]string{root, "--project-name", "Example-Automation", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunInit --json exit = %d; stdout: %s", code, stdout.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse init JSON: %v\n%s", err, stdout.String())
	}
	if payload["inspection"] == nil || payload["state"] == nil {
		t.Fatalf("init JSON missing expected keys: %#v", payload)
	}
}

func TestValueFlagsSupportEqualsFormAndRejectMissingValues(t *testing.T) {
	root := t.TempDir()
	equalStyleTarget := filepath.Join(root, "equals")
	code := RunInit([]string{equalStyleTarget, "--project-name=EqualStyle"}, &bytes.Buffer{}, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunInit with equals-style flag exit = %d", code)
	}
	config := readYAMLFile(t, filepath.Join(equalStyleTarget, ".codeheart", "kit.config.yaml"))
	if config["project_display_name"] != "EqualStyle" {
		t.Fatalf("project_display_name = %#v", config["project_display_name"])
	}

	malformedTarget := filepath.Join(root, "malformed")
	var stderr bytes.Buffer
	code = RunInit([]string{malformedTarget, "--project-name", "--json"}, &bytes.Buffer{}, &stderr)
	if code != 2 {
		t.Fatalf("RunInit malformed flag exit = %d, want 2; stderr: %s", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(malformedTarget, ".codeheart")); !os.IsNotExist(err) {
		t.Fatalf("malformed init should not create .codeheart, stat err: %v", err)
	}

	updateTarget := filepath.Join(root, "update")
	if code := RunInit([]string{updateTarget, "--project-name=Update"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunInit update target exit = %d", code)
	}
	stderr.Reset()
	code = RunUpdateCheck([]string{updateTarget, "--latest-version", "--json"}, &bytes.Buffer{}, &stderr)
	if code != 2 {
		t.Fatalf("RunUpdateCheck malformed flag exit = %d, want 2; stderr: %s", code, stderr.String())
	}
	lock, err := lockfile.ReadLock(updateTarget)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	if mapValue(lock["update_check"])["latest_seen_version"] == "--json" {
		t.Fatalf("malformed update-check wrote --json as latest version")
	}
}

func TestLocalFileURLPathForWindowsDriveAndHost(t *testing.T) {
	driveURL, err := url.Parse("file:///C:/Users/Example%20User/latest-release.json")
	if err != nil {
		t.Fatalf("parse drive URL: %v", err)
	}
	if got, want := localFileURLPathForGOOS(driveURL, "windows"), `C:\Users\Example User\latest-release.json`; got != want {
		t.Fatalf("windows drive file URL path = %q, want %q", got, want)
	}

	uncURL, err := url.Parse("file://server/share/latest-release.json")
	if err != nil {
		t.Fatalf("parse UNC URL: %v", err)
	}
	if got, want := localFileURLPathForGOOS(uncURL, "windows"), `\\server\share\latest-release.json`; got != want {
		t.Fatalf("windows UNC file URL path = %q, want %q", got, want)
	}
}

func TestOnboardYesWritesBaseSetupWithoutNativeCapabilityOfferOrCheck(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	code := RunOnboard([]string{"--target", root, "--project-name", "Companyname-Automation", "--yes", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunOnboard exit = %d; stdout: %s", code, stdout.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse onboard JSON: %v\n%s", err, stdout.String())
	}
	if payload["written"] != true {
		t.Fatalf("onboard did not report written: %#v", payload)
	}
	if _, exists := payload["native_capabilities"]; exists {
		t.Fatalf("onboard --yes should not report native capability checks: %#v", payload["native_capabilities"])
	}
	joinedScript := strings.Join(anyStrings(payload["script"]), "\n")
	for _, forbidden := range []string{
		"Should I check and set up these tools now?",
		"After setup writes complete, ask whether to check native Codex capabilities.",
		"documents, spreadsheets, presentations, browser work, and PDFs",
	} {
		if strings.Contains(joinedScript, forbidden) {
			t.Fatalf("onboarding script contains forbidden optional capability prompt %q:\n%s", forbidden, joinedScript)
		}
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	for capability, value := range mapValue(lock["native_capabilities"]) {
		record := mapValue(value)
		if record["status"] != "unknown" || record["command_result_category"] != "not-checked" {
			t.Fatalf("%s capability was checked or installed: %#v", capability, record)
		}
	}
}

func TestOnboardYesRequiresTargetAndProjectName(t *testing.T) {
	var stdout bytes.Buffer
	code := RunOnboard([]string{"--project-name", "Companyname-Automation", "--yes", "--json"}, &stdout, &bytes.Buffer{})
	if code != 2 {
		t.Fatalf("missing target exit = %d", code)
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse missing target JSON: %v", err)
	}
	missing := anyStrings(payload["required_user_decisions_missing"])
	if len(missing) != 1 || missing[0] != "target_folder" {
		t.Fatalf("missing decisions = %#v", missing)
	}
}

func TestSyncRepairsManagedDriftAndCheckPassesWithCLIMarker(t *testing.T) {
	t.Setenv("CODEHEART_OPERATING_KIT_CLI", "1")
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "Example-Automation"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunInit exit = %d", code)
	}
	managed := filepath.Join(root, ".codeheart", "kit", "docs", "agent-interface", "README.md")
	if err := os.WriteFile(managed, []byte("drift\n"), 0o644); err != nil {
		t.Fatalf("write drift: %v", err)
	}
	before, err := CheckRepository(root)
	if err != nil {
		t.Fatalf("CheckRepository before: %v", err)
	}
	if len(before["drift"].([]any)) == 0 {
		t.Fatalf("expected drift before sync")
	}
	if code := RunSync([]string{root}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunSync exit = %d", code)
	}
	var syncJSON bytes.Buffer
	if code := RunSync([]string{root, "--json"}, &syncJSON, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunSync --json exit = %d; stdout: %s", code, syncJSON.String())
	}
	var syncPayload map[string]any
	if err := json.Unmarshal(syncJSON.Bytes(), &syncPayload); err != nil {
		t.Fatalf("parse sync JSON: %v\n%s", err, syncJSON.String())
	}
	if syncPayload["kit_version"] == "" || syncPayload["synced_managed_paths"] == nil {
		t.Fatalf("sync JSON missing expected keys: %#v", syncPayload)
	}
	after, err := CheckRepository(root)
	if err != nil {
		t.Fatalf("CheckRepository after: %v", err)
	}
	if after["ok"] != true || len(after["drift"].([]any)) != 0 {
		t.Fatalf("check after sync = %#v", after)
	}
	var checkText bytes.Buffer
	if code := RunCheck([]string{root}, &checkText, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunCheck text exit = %d; stdout: %s", code, checkText.String())
	}
	if !strings.Contains(checkText.String(), "Operating Kit check") || !strings.Contains(checkText.String(), "OK: true") {
		t.Fatalf("check text output unexpected: %q", checkText.String())
	}
}

func TestCLIAvailableRecognizesInstalledExeName(t *testing.T) {
	t.Setenv("CODEHEART_OPERATING_KIT_CLI", "")
	root := t.TempDir()
	t.Setenv("PATH", root)
	exe := filepath.Join(root, "codeheart-operating-kit.exe")
	if err := os.WriteFile(exe, []byte(""), 0o755); err != nil {
		t.Fatalf("write exe marker: %v", err)
	}
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()
	os.Args = []string{exe}
	if !CLIAvailable() {
		t.Fatalf("CLIAvailable should recognize codeheart-operating-kit.exe as installed")
	}
}

func TestUpdateCheckWritesCadenceAndFailurePreservesDueDate(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "Example-Automation"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunInit exit = %d", code)
	}
	var stdout bytes.Buffer
	code := RunUpdateCheck([]string{root, "--latest-version", "0.2.0", "--now", "2026-06-13T00:00:00Z", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunUpdateCheck exit = %d; stdout: %s", code, stdout.String())
	}
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	update := mapValue(lock["update_check"])
	if update["update_status"] != "update-available" || update["next_update_check_due"] != "2026-06-20T00:00:00Z" {
		t.Fatalf("unexpected update state: %#v", update)
	}
	beforeDue := update["next_update_check_due"]
	stdout.Reset()
	missingURL := (&urlBuilder{path: filepath.Join(root, "missing-release.json")}).String()
	code = RunUpdateCheck([]string{root, "--metadata-url", missingURL, "--now", "2026-06-14T00:00:00Z", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunUpdateCheck failure case exit = %d; stdout: %s", code, stdout.String())
	}
	lock, err = lockfile.ReadLock(root)
	if err != nil {
		t.Fatalf("ReadLock after failure: %v", err)
	}
	update = mapValue(lock["update_check"])
	if update["update_status"] != "failed" || update["next_update_check_due"] != beforeDue {
		t.Fatalf("failed update state = %#v, previous due %v", update, beforeDue)
	}
	var text bytes.Buffer
	code = RunUpdateCheck([]string{root, "--latest-version", "0.1.19"}, &text, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunUpdateCheck text exit = %d; stdout: %s", code, text.String())
	}
	if !strings.Contains(text.String(), "Operating Kit is current.") {
		t.Fatalf("update-check text output unexpected: %q", text.String())
	}
}

func TestUpdateCheckEqualsStyleLatestVersion(t *testing.T) {
	root := t.TempDir()
	if code := RunInit([]string{root, "--project-name", "Example-Automation"}, &bytes.Buffer{}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("RunInit exit = %d", code)
	}
	var stdout bytes.Buffer
	code := RunUpdateCheck([]string{root, "--latest-version=0.2.0", "--now=2026-06-13T00:00:00Z", "--json"}, &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("RunUpdateCheck equals-style exit = %d; stdout: %s", code, stdout.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("parse update JSON: %v\n%s", err, stdout.String())
	}
	if payload["status"] != "update-available" {
		t.Fatalf("update status = %#v", payload)
	}
}

type urlBuilder struct {
	path string
}

func (builder *urlBuilder) String() string {
	return "file://" + filepath.ToSlash(builder.path)
}

func readYAMLFile(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read YAML %s: %v", path, err)
	}
	parsed, err := yamlmini.MustMap(string(data))
	if err != nil {
		t.Fatalf("parse YAML %s: %v", path, err)
	}
	return parsed
}

func anyStrings(value any) []string {
	result := []string{}
	for _, item := range anyList(value) {
		result = append(result, valueString(item))
	}
	return result
}

func mapKeys(value map[string]any) []string {
	keys := make([]string, 0, len(value))
	for key := range value {
		keys = append(keys, key)
	}
	return keys
}

func sameStringSet(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	seen := map[string]int{}
	for _, value := range left {
		seen[value]++
	}
	for _, value := range right {
		seen[value]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}
