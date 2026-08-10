package plancatalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
)

type activationTopologyFixture struct {
	Root             string
	EvidenceRevision string
	LedgerRevision   string
	Activation       string
	ActionDigest     string
}

func newActivationTopologyFixture(t *testing.T) activationTopologyFixture {
	t.Helper()
	root := t.TempDir()
	runGitTest(t, root, "init", "-b", "main")
	runGitTest(t, root, "config", "user.name", "Topology Test")
	runGitTest(t, root, "config", "user.email", "topology@example.invalid")
	writeCheckpointFixture(t, root, "base.txt", []byte("evidence\n"))
	runGitTest(t, root, "add", "base.txt")
	runGitTest(t, root, "commit", "-m", "record evidence")
	evidence := runGitTest(t, root, "rev-parse", "HEAD")
	writeCheckpointFixture(t, root, "reviewed-ledger.yaml", []byte("schema_version: 3\n"))
	runGitTest(t, root, "add", "reviewed-ledger.yaml")
	runGitTest(t, root, "commit", "-m", "record ledger")
	ledger := runGitTest(t, root, "rev-parse", "HEAD")
	writeCheckpointFixture(t, root, "docs/repo/plans/example/example_implementation_doc.md", []byte("reviewed activation\n"))
	runGitTest(t, root, "add", "docs/repo/plans/example/example_implementation_doc.md")
	runGitTest(t, root, "commit", "-m", "record activation")
	activation := runGitTest(t, root, "rev-parse", "HEAD")
	digest, _, err := MigrationActionDigestForCheckpoint(root, ledger, activation)
	if err != nil {
		t.Fatal(err)
	}
	return activationTopologyFixture{Root: root, EvidenceRevision: evidence, LedgerRevision: ledger, Activation: activation, ActionDigest: digest}
}

func TestBoundActivationCheckpointRevisionAcrossMergeTopologies(t *testing.T) {
	t.Run("exact activation", func(t *testing.T) {
		fixture := newActivationTopologyFixture(t)
		actual, err := BoundActivationCheckpointRevision(fixture.Root, fixture.LedgerRevision, fixture.Activation, fixture.ActionDigest)
		if err != nil || actual != fixture.Activation {
			t.Fatalf("exact activation=%s want=%s err=%v", actual, fixture.Activation, err)
		}
	})

	t.Run("normal merge and descendant", func(t *testing.T) {
		fixture := newActivationTopologyFixture(t)
		runGitTest(t, fixture.Root, "branch", "reviewed-activation", fixture.Activation)
		runGitTest(t, fixture.Root, "switch", "-c", "integration-main", fixture.EvidenceRevision)
		runGitTest(t, fixture.Root, "merge", "--no-ff", "reviewed-activation", "-m", "merge reviewed activation")
		merge := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
		actual, err := BoundActivationCheckpointRevision(fixture.Root, fixture.LedgerRevision, merge, fixture.ActionDigest)
		if err != nil || actual != fixture.Activation {
			t.Fatalf("merged activation=%s want=%s err=%v", actual, fixture.Activation, err)
		}
		writeCheckpointFixture(t, fixture.Root, "later.txt", []byte("safe descendant\n"))
		runGitTest(t, fixture.Root, "add", "later.txt")
		runGitTest(t, fixture.Root, "commit", "-m", "record safe descendant")
		descendant := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
		actual, err = BoundActivationCheckpointRevision(fixture.Root, fixture.LedgerRevision, descendant, fixture.ActionDigest)
		if err != nil || actual != fixture.Activation {
			t.Fatalf("descendant activation=%s want=%s err=%v", actual, fixture.Activation, err)
		}
	})

	t.Run("merge first parent is ledger and second parent is activation", func(t *testing.T) {
		fixture := newActivationTopologyFixture(t)
		tree := runGitTest(t, fixture.Root, "rev-parse", fixture.Activation+"^{tree}")
		merged := runGitTest(t, fixture.Root, "commit-tree", tree, "-p", fixture.LedgerRevision, "-p", fixture.Activation, "-m", "merge reviewed activation from ledger main")
		actual, err := BoundActivationCheckpointRevision(fixture.Root, fixture.LedgerRevision, merged, fixture.ActionDigest)
		if err != nil || actual != fixture.Activation {
			t.Fatalf("ledger-first-parent merge activation=%s want=%s err=%v", actual, fixture.Activation, err)
		}
	})

	t.Run("unrelated sibling", func(t *testing.T) {
		fixture := newActivationTopologyFixture(t)
		runGitTest(t, fixture.Root, "switch", "-c", "unrelated-sibling", fixture.LedgerRevision)
		writeCheckpointFixture(t, fixture.Root, "unrelated.txt", []byte("not activation\n"))
		runGitTest(t, fixture.Root, "add", "unrelated.txt")
		runGitTest(t, fixture.Root, "commit", "-m", "record unrelated sibling")
		sibling := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
		if _, err := BoundActivationCheckpointRevision(fixture.Root, fixture.LedgerRevision, sibling, fixture.ActionDigest); err == nil || !strings.HasPrefix(err.Error(), "activation_checkpoint_missing:") {
			t.Fatalf("unrelated sibling was accepted: %v", err)
		}
	})

	t.Run("tree-only reparented coincidence", func(t *testing.T) {
		fixture := newActivationTopologyFixture(t)
		tree := runGitTest(t, fixture.Root, "rev-parse", fixture.Activation+"^{tree}")
		forged := runGitTest(t, fixture.Root, "commit-tree", tree, "-p", fixture.EvidenceRevision, "-m", "reparent activation tree")
		if _, err := BoundActivationCheckpointRevision(fixture.Root, fixture.LedgerRevision, forged, fixture.ActionDigest); err == nil || !strings.HasPrefix(err.Error(), "activation_checkpoint_missing:") {
			t.Fatalf("tree-only reparented coincidence was accepted: %v", err)
		}
	})

	t.Run("wrong action digest", func(t *testing.T) {
		fixture := newActivationTopologyFixture(t)
		if _, err := BoundActivationCheckpointRevision(fixture.Root, fixture.LedgerRevision, fixture.Activation, strings.Repeat("f", 64)); err == nil || !strings.HasPrefix(err.Error(), "activation_checkpoint_missing:") {
			t.Fatalf("wrong action digest was accepted: %v", err)
		}
	})

	t.Run("ambiguous matching children", func(t *testing.T) {
		fixture := newActivationTopologyFixture(t)
		tree := runGitTest(t, fixture.Root, "rev-parse", fixture.Activation+"^{tree}")
		duplicate := runGitTest(t, fixture.Root, "commit-tree", tree, "-p", fixture.LedgerRevision, "-m", "duplicate activation")
		merged := runGitTest(t, fixture.Root, "commit-tree", tree, "-p", fixture.Activation, "-p", duplicate, "-m", "incorporate both activations")
		if _, err := BoundActivationCheckpointRevision(fixture.Root, fixture.LedgerRevision, merged, fixture.ActionDigest); err == nil || !strings.HasPrefix(err.Error(), "activation_checkpoint_ambiguous:") {
			t.Fatalf("ambiguous matching activation children were accepted: %v", err)
		}
	})
}

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
