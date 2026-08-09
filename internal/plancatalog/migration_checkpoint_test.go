package plancatalog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
)

func TestLedgerCheckpointRequiresExactFirstParentAndLedgerOnlyDelta(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	runGitTest(t, root, "config", "user.name", "Test")
	runGitTest(t, root, "config", "user.email", "test@example.com")
	writeCheckpointFixture(t, root, "README.md", []byte("baseline\n"))
	runGitTest(t, root, "add", "README.md")
	runGitTest(t, root, "commit", "-m", "evidence")
	evidence := runGitTest(t, root, "rev-parse", "HEAD")
	ledgerPath := "docs/repo/plans/example/attachments/ledger-v3.yaml"
	ledger := []byte("schema_version: 3\n")
	writeCheckpointFixture(t, root, ledgerPath, ledger)
	runGitTest(t, root, "add", ledgerPath)
	runGitTest(t, root, "commit", "-m", "ledger")
	_, digest, err := bindLedgerArtifact(root, filepath.Join(root, filepath.FromSlash(ledgerPath)), ledger)
	if err != nil {
		t.Fatal(err)
	}
	checkpoint, problems := verifyLedgerCheckpoint(root, evidence, ledgerPath, digest)
	if HasErrors(problems) || checkpoint.EvidenceRevision != evidence || checkpoint.ActivationBaseRevision == "" {
		t.Fatalf("checkpoint=%#v problems=%#v", checkpoint, problems)
	}

	writeCheckpointFixture(t, root, "unrelated.txt", []byte("unexpected\n"))
	runGitTest(t, root, "add", "unrelated.txt")
	if _, problems = verifyLedgerCheckpoint(root, evidence, ledgerPath, digest); !problemCodePresent(problems, "ledger_checkpoint_index_dirty") {
		t.Fatalf("dirty checkpoint problems=%#v", problems)
	}
}

func TestMigrationActionDigestAndExactWorktreeValidation(t *testing.T) {
	root := t.TempDir()
	runGitTest(t, root, "init")
	runGitTest(t, root, "config", "user.name", "Test")
	runGitTest(t, root, "config", "user.email", "test@example.com")
	writeCheckpointFixture(t, root, "replace.md", []byte("before\n"))
	writeCheckpointFixture(t, root, "remove.md", []byte("remove\n"))
	runGitTest(t, root, "add", "replace.md", "remove.md")
	runGitTest(t, root, "commit", "-m", "baseline")
	actions := []reconcile.Action{
		{Kind: "replace", Target: "replace.md", Content: []byte("after\n"), Mode: 0o644, ExpectedSHA256: sha256Text([]byte("before\n"))},
		{Kind: "remove", Target: "remove.md", ExpectedSHA256: sha256Text([]byte("remove\n"))},
		{Kind: "create", Target: "create.md", Content: []byte("created\n"), Mode: 0o644},
	}
	writeCheckpointFixture(t, root, "replace.md", []byte("after\n"))
	if err := os.Remove(filepath.Join(root, "remove.md")); err != nil {
		t.Fatal(err)
	}
	writeCheckpointFixture(t, root, "create.md", []byte("created\n"))
	if problems := verifyMigrationWorktree(root, actions); HasErrors(problems) {
		t.Fatalf("exact worktree problems=%#v", problems)
	}
	first := migrationActionDigest(actions)
	second := migrationActionDigest([]reconcile.Action{actions[2], actions[0], actions[1]})
	if first == "" || first != second {
		t.Fatalf("action digests first=%s second=%s", first, second)
	}
	writeCheckpointFixture(t, root, "extra.md", []byte("extra\n"))
	if problems := verifyMigrationWorktree(root, actions); !problemCodePresent(problems, "catalog_activation_worktree_mismatch") {
		t.Fatalf("extra path problems=%#v", problems)
	}
}

func TestMigrationWorktreeCRLFRequiresCorroboratingCheckoutPolicy(t *testing.T) {
	setup := func(t *testing.T, attributes []byte) (string, []reconcile.Action) {
		t.Helper()
		root := t.TempDir()
		runGitTest(t, root, "init")
		runGitTest(t, root, "config", "user.name", "Test")
		runGitTest(t, root, "config", "user.email", "test@example.com")
		if attributes != nil {
			writeCheckpointFixture(t, root, ".gitattributes", attributes)
		}
		writeCheckpointFixture(t, root, "replace.md", []byte("before\n"))
		runGitTest(t, root, "add", ".")
		runGitTest(t, root, "commit", "-m", "baseline")
		writeCheckpointFixture(t, root, "replace.md", []byte("after\r\n"))
		return root, []reconcile.Action{{Kind: "replace", Target: "replace.md", Content: []byte("after\n"), Mode: 0o644, ExpectedSHA256: sha256Text([]byte("before\n"))}}
	}

	t.Run("core autocrlf", func(t *testing.T) {
		root, actions := setup(t, nil)
		runGitTest(t, root, "config", "core.autocrlf", "false")
		if problems := verifyMigrationWorktree(root, actions); !problemCodePresent(problems, "catalog_activation_action_mismatch") {
			t.Fatalf("CRLF drift without checkout authority was accepted: %#v", problems)
		}
		runGitTest(t, root, "config", "core.autocrlf", "true")
		if problems := verifyMigrationWorktree(root, actions); HasErrors(problems) {
			t.Fatalf("corroborated CRLF checkout projection was rejected: %#v", problems)
		}
	})

	t.Run("text disabled", func(t *testing.T) {
		root, actions := setup(t, []byte("*.md -text\n"))
		runGitTest(t, root, "config", "core.autocrlf", "true")
		if problems := verifyMigrationWorktree(root, actions); !problemCodePresent(problems, "catalog_activation_action_mismatch") {
			t.Fatalf("CRLF drift under -text was accepted: %#v", problems)
		}
	})

	t.Run("autocrlf input overrides core eol", func(t *testing.T) {
		root, _ := setup(t, []byte("*.md text\n"))
		runGitTest(t, root, "config", "core.autocrlf", "input")
		runGitTest(t, root, "config", "core.eol", "crlf")
		if permitted, err := checkoutPolicyUsesCRLF(root, "replace.md", true); err != nil || permitted {
			t.Fatalf("core.autocrlf=input credited output conversion: permitted=%t err=%v", permitted, err)
		}
	})

	t.Run("native eol follows platform", func(t *testing.T) {
		root, _ := setup(t, []byte("*.md text\n"))
		runGitTest(t, root, "config", "core.autocrlf", "false")
		runGitTest(t, root, "config", "core.eol", "native")
		if permitted, err := checkoutPolicyUsesCRLF(root, "replace.md", true); err != nil || !permitted {
			t.Fatalf("Windows-native CRLF policy was not credited: permitted=%t err=%v", permitted, err)
		}
		if permitted, err := checkoutPolicyUsesCRLF(root, "replace.md", false); err != nil || permitted {
			t.Fatalf("LF-native policy credited CRLF conversion: permitted=%t err=%v", permitted, err)
		}
	})
}

func writeCheckpointFixture(t *testing.T, root, relative string, data []byte) {
	t.Helper()
	absolute := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolute, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
