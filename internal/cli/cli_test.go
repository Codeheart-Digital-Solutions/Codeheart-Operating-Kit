package cli

import (
	"bytes"
	"strings"
	"testing"
)

func runForTest(args ...string) (int, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestRootHelpListsCommands(t *testing.T) {
	code, stdout, stderr := runForTest("--help")
	if code != 0 {
		t.Fatalf("help exit code = %d, want 0; stderr: %s", code, stderr)
	}
	for _, command := range []string{"onboard", "inspect", "init", "repair", "sync", "check", "update-check", "upgrade"} {
		if !strings.Contains(stdout, command) {
			t.Fatalf("root help did not list %q:\n%s", command, stdout)
		}
	}
}

func TestNoArgsRequiresCommand(t *testing.T) {
	code, stdout, stderr := runForTest()
	if code != 2 {
		t.Fatalf("no-arg exit code = %d, want 2", code)
	}
	if stdout != "" {
		t.Fatalf("no-arg command wrote stdout: %q", stdout)
	}
	if !strings.Contains(stderr, "the following arguments are required: command") {
		t.Fatalf("no-arg stderr did not explain error:\n%s", stderr)
	}
}

func TestVersionOutput(t *testing.T) {
	code, stdout, stderr := runForTest("--version")
	if code != 0 {
		t.Fatalf("version exit code = %d, want 0; stderr: %s", code, stderr)
	}
	if !strings.HasPrefix(stdout, "codeheart-operating-kit ") {
		t.Fatalf("version output has unexpected prefix: %q", stdout)
	}
}

func TestSubcommandHelp(t *testing.T) {
	for _, command := range []string{"onboard", "inspect", "init", "repair", "sync", "check", "update-check", "upgrade"} {
		t.Run(command, func(t *testing.T) {
			code, stdout, stderr := runForTest(command, "--help")
			if code != 0 {
				t.Fatalf("%s help exit code = %d, want 0; stderr: %s", command, code, stderr)
			}
			if !strings.Contains(stdout, "codeheart-operating-kit "+command) {
				t.Fatalf("%s help missing usage:\n%s", command, stdout)
			}
			if !strings.Contains(stdout, "--json") {
				t.Fatalf("%s help missing --json option:\n%s", command, stdout)
			}
		})
	}
}

func TestMutatingLifecycleHelpDocumentsDryRun(t *testing.T) {
	for _, command := range []string{"init", "repair", "sync", "update-check", "upgrade"} {
		_, stdout, _ := runForTest(command, "--help")
		if !strings.Contains(stdout, "--dry-run") {
			t.Fatalf("%s help missing --dry-run:\n%s", command, stdout)
		}
	}
}

func TestUpgradeHelpRequiresExplicitApprovalAndHidesHandoff(t *testing.T) {
	_, stdout, _ := runForTest("upgrade", "--help")
	for _, expected := range []string{"--version VERSION", "--dry-run", "--yes", "--catalog CATALOG"} {
		if !strings.Contains(stdout, expected) {
			t.Fatalf("upgrade help missing %s:\n%s", expected, stdout)
		}
	}
	rootCode, rootHelp, _ := runForTest("--help")
	if rootCode != 0 || strings.Contains(rootHelp, "__upgrade-handoff") || strings.Contains(rootHelp, "__upgrade-reconcile") || strings.Contains(rootHelp, "__verify-content-identity") || strings.Contains(rootHelp, "__verify-release-evidence") || strings.Contains(rootHelp, "__cleanup-upgrade-handoff") {
		t.Fatalf("internal handoff leaked into help:\n%s", rootHelp)
	}
	code, _, _ := runForTest("upgrade", "--version", "0.2.0")
	if code != 2 {
		t.Fatalf("upgrade without explicit dry-run/yes exit = %d", code)
	}
}

func TestSubcommandHelpAfterOtherOptions(t *testing.T) {
	code, stdout, stderr := runForTest("init", "--json", "--help")
	if code != 0 {
		t.Fatalf("delayed help exit code = %d, want 0; stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "codeheart-operating-kit init") || !strings.Contains(stdout, "--project-name") {
		t.Fatalf("delayed help missing init usage:\n%s", stdout)
	}
	if stderr != "" {
		t.Fatalf("delayed help wrote stderr: %q", stderr)
	}
}

func TestSubcommandHelpInValuePositionIsMissingValue(t *testing.T) {
	tests := [][]string{
		{"init", "--project-name", "--help"},
		{"init", "--project-name", "--json", "--help"},
		{"init", "--project-name", "-h"},
		{"onboard", "--target", "--help"},
		{"onboard", "--target", "--yes", "--help"},
		{"update-check", "--latest-version", "--help"},
		{"update-check", "--latest-version", "--json", "--help"},
	}
	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			code, stdout, stderr := runForTest(args...)
			if code != 2 {
				t.Fatalf("exit code = %d, want 2; stdout: %s stderr: %s", code, stdout, stderr)
			}
			if stdout != "" {
				t.Fatalf("missing-value help wrote stdout: %q", stdout)
			}
			if !strings.Contains(stderr, "requires a value") {
				t.Fatalf("missing-value help stderr did not explain missing value:\n%s", stderr)
			}
		})
	}
}

func TestOnboardHelpDocumentsApprovedSetupBoundary(t *testing.T) {
	code, stdout, stderr := runForTest("onboard", "--help")
	if code != 0 {
		t.Fatalf("onboard help exit code = %d, want 0; stderr: %s", code, stderr)
	}
	normalized := strings.Join(strings.Fields(stdout), " ")
	for _, expected := range []string{
		"Non-interactive --yes writes require explicit --target and --project-name values",
		"Required with --yes",
		"does not select a different profile",
		"Base onboarding does not install, offer, or implicitly check optional native capabilities",
	} {
		if !strings.Contains(normalized, expected) {
			t.Fatalf("onboard help missing %q:\n%s", expected, stdout)
		}
	}
}

func TestMissingCommandErrors(t *testing.T) {
	code, stdout, stderr := runForTest("missing")
	if code != 2 {
		t.Fatalf("missing command exit code = %d, want 2", code)
	}
	if stdout != "" {
		t.Fatalf("missing command wrote stdout: %q", stdout)
	}
	if !strings.Contains(stderr, "invalid choice: 'missing'") {
		t.Fatalf("missing command stderr did not explain error:\n%s", stderr)
	}
}

func TestKnownCommandDispatches(t *testing.T) {
	code, stdout, stderr := runForTest("inspect", ".")
	if code != 0 {
		t.Fatalf("inspect exit code = %d, want 0; stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, ":") {
		t.Fatalf("inspect stdout did not contain mode output: %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("inspect wrote stderr: %q", stderr)
	}
}
