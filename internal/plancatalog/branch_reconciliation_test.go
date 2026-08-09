package plancatalog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const branchEvidencePlanPath = "docs/repo/plans/example/example_implementation_doc.md"

func TestBranchEvidenceSameContentLocalAndRemoteRefs(t *testing.T) {
	root := newBranchEvidenceRepository(t)
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("# Plan\n\nbase\n"))
	commitBranchEvidenceAll(t, root, "base")
	base := runGitTest(t, root, "rev-parse", "HEAD")

	runGitTest(t, root, "switch", "-c", "history")
	writeBranchEvidenceFile(t, root, "history-only.txt", []byte("historical\n"))
	commitBranchEvidenceAll(t, root, "historical side change")
	historyTip := runGitTest(t, root, "rev-parse", "HEAD")
	runGitTest(t, root, "update-ref", "refs/remotes/origin/history", historyTip)

	runGitTest(t, root, "switch", "main")
	writeBranchEvidenceFile(t, root, "target-only.txt", []byte("target\n"))
	commitBranchEvidenceAll(t, root, "target side change")
	target := runGitTest(t, root, "rev-parse", "HEAD")
	if target == base || target == historyTip {
		t.Fatal("fixture did not create distinct histories")
	}

	for _, test := range []struct {
		name       string
		ref        string
		logicalRef string
	}{
		{name: "local", ref: "refs/heads/history", logicalRef: "local:refs/heads/history"},
		{name: "remote", ref: "refs/remotes/origin/history", logicalRef: "tracking:origin:refs/heads/history"},
	} {
		t.Run(test.name, func(t *testing.T) {
			result := RecomputeBranchCandidateEvidence(branchEvidenceRequest(root, target, test.ref, test.logicalRef, branchEvidencePlanPath))
			if result.ProposedDisposition != BranchDispositionSameContentNonOwner || result.Proof == nil || result.Proof.Kind != BranchProofSameContentV1 {
				t.Fatalf("expected same-content proof, got disposition=%q proof=%#v blockers=%#v", result.ProposedDisposition, result.Proof, result.Blockers)
			}
			if result.RefTip != historyTip || result.MergeBase != base || result.Digest == "" {
				t.Fatalf("evidence bindings missing: %#v", result)
			}
			if !ExactBranchPathIdentity(*result.Proof.TargetPathState, *result.Proof.RefPathState) {
				t.Fatal("same-content proof does not carry exact target/ref identity")
			}
		})
	}
}

func TestBranchEvidenceProvesSingleSquashEquivalentCommit(t *testing.T) {
	root, target, incorporated := newBranchEvidenceSquashRepository(t, []byte("middle\n"), []byte("middle\n"))
	request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)
	request.IncorporatedCommit = incorporated
	result := RecomputeBranchCandidateEvidence(request)
	if result.ProposedDisposition != BranchDispositionIncorporatedHistoryNonOwner || result.Proof == nil {
		t.Fatalf("expected incorporated-history proof, got disposition=%q proof=%#v blockers=%#v", result.ProposedDisposition, result.Proof, result.Blockers)
	}
	if result.Proof.Kind != BranchProofIncorporatedHistoryV1 || result.Proof.TransitionDigest == "" || result.Proof.StablePatchID == "" {
		t.Fatalf("incorporated proof is incomplete: %#v", result.Proof)
	}
	if result.Proof.TransitionDigest != result.Proof.IncorporatedTransitionDigest || result.Proof.StablePatchID != result.Proof.IncorporatedStablePatchID {
		t.Fatalf("two-factor proof does not match: %#v", result.Proof)
	}
	if len(result.AttributeSnapshots) != 4 {
		t.Fatalf("expected four attribute endpoints, got %d", len(result.AttributeSnapshots))
	}
	second := RecomputeBranchCandidateEvidence(request)
	if result.Digest != second.Digest {
		t.Fatalf("evidence digest is not reproducible: %s != %s", result.Digest, second.Digest)
	}
	tampered := result
	tampered.RefTip = strings.Repeat("f", len(result.RefTip))
	if CanonicalBranchEvidenceDigest(tampered) == result.Digest {
		t.Fatal("canonical evidence digest did not bind the exact ref tip")
	}
	portable := result
	portable.GitVersion = "9.9.9"
	if CanonicalBranchEvidenceDigest(portable) != result.Digest {
		t.Fatal("canonical evidence digest varied across admitted Git implementation versions")
	}
}

func TestBranchEvidenceAggregateDigestIsPortableAcrossGitVersions(t *testing.T) {
	evidence := BranchCandidateEvidence{
		Algorithm: BranchEvidenceAlgorithm, GitVersion: "2.43.0",
		Identity:         BranchReviewIdentity{EvidenceScope: "local", RepositoryID: "example", LogicalRef: "local:refs/heads/history", CandidatePath: "docs/repo/plans/example/example_discovery_doc.md"},
		EvidenceRevision: strings.Repeat("1", 40), PolicyDigest: strings.Repeat("2", 64), CandidateSetDigest: strings.Repeat("3", 64),
		RefTip: strings.Repeat("4", 40), MergeBase: strings.Repeat("5", 40), ProposedDisposition: BranchDispositionBlocking,
		Blockers: []BranchEvidenceBlocker{{Code: "branch_clearance_unproven", Message: "synthetic blocker"}},
	}
	evidence.Digest = CanonicalBranchEvidenceDigest(evidence)
	first := branchEvidenceSetDigest([]BranchCandidateEvidence{evidence}, "")
	evidence.GitVersion = "9.9.9"
	evidence.Digest = CanonicalBranchEvidenceDigest(evidence)
	if second := branchEvidenceSetDigest([]BranchCandidateEvidence{evidence}, ""); second != first {
		t.Fatalf("aggregate evidence digest changed across admitted Git versions: %s != %s", first, second)
	}
}

func TestBranchEvidenceRejectsPartialAndPatchNormalizedCollisions(t *testing.T) {
	t.Run("exact transition rejects patch-normalized collision", func(t *testing.T) {
		root, target, incorporated := newBranchEvidenceSquashRepository(t, []byte("a b\n"), []byte("a  b\n"))
		request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)
		request.IncorporatedCommit = incorporated
		result := RecomputeBranchCandidateEvidence(request)
		assertBranchEvidenceBlocker(t, result, "branch_evidence_transition_mismatch")
	})

	t.Run("partial candidate path set cannot clear", func(t *testing.T) {
		root := newBranchEvidenceRepository(t)
		secondPath := "docs/repo/plans/example/example_execution_log.md"
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("plan base\n"))
		writeBranchEvidenceFile(t, root, secondPath, []byte("log base\n"))
		commitBranchEvidenceAll(t, root, "base")
		runGitTest(t, root, "switch", "-c", "history")
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("plan historical\n"))
		writeBranchEvidenceFile(t, root, secondPath, []byte("log historical\n"))
		commitBranchEvidenceAll(t, root, "historical complete transition")
		runGitTest(t, root, "switch", "main")
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("plan historical\n"))
		commitBranchEvidenceAll(t, root, "incorporate only plan")
		incorporated := runGitTest(t, root, "rev-parse", "HEAD")
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("plan current\n"))
		commitBranchEvidenceAll(t, root, "advance target")
		target := runGitTest(t, root, "rev-parse", "HEAD")

		request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)
		request.Paths = []string{secondPath, branchEvidencePlanPath}
		request.IncorporatedCommit = incorporated
		result := RecomputeBranchCandidateEvidence(request)
		assertBranchEvidenceBlocker(t, result, "branch_evidence_transition_mismatch")
	})
}

func TestBranchEvidenceRenameAndDeleteTransitions(t *testing.T) {
	t.Run("rename binds old and new paths", func(t *testing.T) {
		root := newBranchEvidenceRepository(t)
		oldPath := "docs/repo/plans/example/old_implementation_doc.md"
		newPath := branchEvidencePlanPath
		writeBranchEvidenceFile(t, root, oldPath, []byte("old\n"))
		commitBranchEvidenceAll(t, root, "base")
		runGitTest(t, root, "switch", "-c", "history")
		runGitTest(t, root, "mv", oldPath, newPath)
		commitBranchEvidenceAll(t, root, "historical rename")
		runGitTest(t, root, "switch", "main")
		runGitTest(t, root, "mv", oldPath, newPath)
		commitBranchEvidenceAll(t, root, "incorporated rename")
		incorporated := runGitTest(t, root, "rev-parse", "HEAD")
		writeBranchEvidenceFile(t, root, newPath, []byte("current\n"))
		commitBranchEvidenceAll(t, root, "advance renamed plan")
		target := runGitTest(t, root, "rev-parse", "HEAD")

		request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", newPath)
		request.Paths = []string{newPath, oldPath}
		request.IncorporatedCommit = incorporated
		result := RecomputeBranchCandidateEvidence(request)
		if result.ProposedDisposition != BranchDispositionIncorporatedHistoryNonOwner {
			t.Fatalf("rename proof failed: %#v", result.Blockers)
		}
		if len(result.PathTransitions) != 2 || result.PathTransitions[0].Before.State != BranchPathAbsent || result.PathTransitions[1].After.State != BranchPathAbsent {
			t.Fatalf("rename endpoints not retained: %#v", result.PathTransitions)
		}
	})

	t.Run("delete binds present to absent", func(t *testing.T) {
		root := newBranchEvidenceRepository(t)
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("base\n"))
		commitBranchEvidenceAll(t, root, "base")
		runGitTest(t, root, "switch", "-c", "history")
		runGitTest(t, root, "rm", branchEvidencePlanPath)
		commitBranchEvidenceAll(t, root, "historical delete")
		runGitTest(t, root, "switch", "main")
		runGitTest(t, root, "rm", branchEvidencePlanPath)
		commitBranchEvidenceAll(t, root, "incorporated delete")
		incorporated := runGitTest(t, root, "rev-parse", "HEAD")
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("recreated current\n"))
		commitBranchEvidenceAll(t, root, "recreate target plan")
		target := runGitTest(t, root, "rev-parse", "HEAD")

		request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)
		request.IncorporatedCommit = incorporated
		result := RecomputeBranchCandidateEvidence(request)
		if result.ProposedDisposition != BranchDispositionIncorporatedHistoryNonOwner {
			t.Fatalf("delete proof failed: %#v", result.Blockers)
		}
		transition := result.PathTransitions[0]
		if transition.Before.State != BranchPathPresent || transition.After.State != BranchPathAbsent {
			t.Fatalf("delete endpoints not retained: %#v", transition)
		}
	})
}

func TestBranchEvidenceGitVersionFloorAndCanonicalTransitionDigest(t *testing.T) {
	for _, test := range []struct {
		output      string
		wantVersion string
		wantBlock   bool
	}{
		{output: "git version 2.43.0", wantVersion: "2.43.0"},
		{output: "git version 2.47.1.windows.1", wantVersion: "2.47.1"},
		{output: "git version 2.42.9", wantVersion: "2.42.9", wantBlock: true},
		{output: "unexpected", wantBlock: true},
	} {
		version, blocker := parseBranchEvidenceGitVersion(test.output)
		if version != test.wantVersion || (blocker != nil) != test.wantBlock {
			t.Fatalf("parse %q = version %q blocker %#v", test.output, version, blocker)
		}
		if blocker != nil && blocker.Code != "branch_evidence_git_version_unsupported" {
			t.Fatalf("unexpected version blocker: %#v", blocker)
		}
	}

	present := BranchPathState{State: BranchPathPresent, Mode: GitModeRegular, ObjectID: strings.Repeat("1", 40), SHA256: strings.Repeat("a", 64)}
	absent := BranchPathState{State: BranchPathAbsent}
	left := []BranchPathTransition{{Path: "z.md", Before: absent, After: present}, {Path: "a.md", Before: present, After: absent}}
	right := []BranchPathTransition{left[1], left[0]}
	leftDigest, blocker := CanonicalBranchTransitionDigest(left)
	if blocker != nil {
		t.Fatal(blocker)
	}
	rightDigest, blocker := CanonicalBranchTransitionDigest(right)
	if blocker != nil {
		t.Fatal(blocker)
	}
	if leftDigest != rightDigest {
		t.Fatalf("transition digest depends on caller order: %s != %s", leftDigest, rightDigest)
	}
}

func TestBranchEvidenceMovedRefInvalidatesReview(t *testing.T) {
	root := newBranchEvidenceRepository(t)
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("base\n"))
	commitBranchEvidenceAll(t, root, "base")
	runGitTest(t, root, "switch", "-c", "history")
	writeBranchEvidenceFile(t, root, "history.txt", []byte("one\n"))
	commitBranchEvidenceAll(t, root, "history one")
	reviewedTip := runGitTest(t, root, "rev-parse", "HEAD")
	writeBranchEvidenceFile(t, root, "history.txt", []byte("two\n"))
	commitBranchEvidenceAll(t, root, "history two")

	blocker := VerifyBranchRefTip(root, "refs/heads/history", reviewedTip)
	if blocker == nil || blocker.Code != "branch_evidence_ref_moved" {
		t.Fatalf("expected moved-ref blocker, got %#v", blocker)
	}
	runGitTest(t, root, "switch", "main")
	target := runGitTest(t, root, "rev-parse", "HEAD")
	request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)
	request.ExpectedRefTip = reviewedTip
	assertBranchEvidenceBlocker(t, RecomputeBranchCandidateEvidence(request), "branch_evidence_ref_moved")
}

func TestBranchEvidenceIgnoresHostileConfigButBlocksUnsupportedAttributesAndBinary(t *testing.T) {
	t.Run("external diff and injected config are inert", func(t *testing.T) {
		root, target, incorporated := newBranchEvidenceSquashRepository(t, []byte("middle\n"), []byte("middle\n"))
		marker := filepath.Join(t.TempDir(), "external-diff-ran")
		script := filepath.Join(t.TempDir(), "hostile-diff.sh")
		writeBranchEvidenceFile(t, filepath.Dir(script), filepath.Base(script), []byte("#!/bin/sh\ntouch \""+marker+"\"\nexit 1\n"))
		if err := os.Chmod(script, 0o700); err != nil {
			t.Fatal(err)
		}
		runGitTest(t, root, "config", "diff.external", script)
		t.Setenv("GIT_EXTERNAL_DIFF", script)
		t.Setenv("GIT_CONFIG_COUNT", "1")
		t.Setenv("GIT_CONFIG_KEY_0", "diff.external")
		t.Setenv("GIT_CONFIG_VALUE_0", script)

		request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)
		request.IncorporatedCommit = incorporated
		result := RecomputeBranchCandidateEvidence(request)
		if result.ProposedDisposition != BranchDispositionIncorporatedHistoryNonOwner {
			t.Fatalf("hostile config affected proof: %#v", result.Blockers)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("external diff executed; stat error=%v", err)
		}
	})

	t.Run("custom attributes block incorporated proof", func(t *testing.T) {
		root := newBranchEvidenceRepository(t)
		writeBranchEvidenceFile(t, root, ".gitattributes", []byte("*.md diff=hostile\n"))
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("base\n"))
		commitBranchEvidenceAll(t, root, "base with custom attributes")
		target, incorporated := completeBranchEvidenceSquashHistory(t, root, []byte("middle\n"), []byte("middle\n"))
		request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)
		request.IncorporatedCommit = incorporated
		assertBranchEvidenceBlocker(t, RecomputeBranchCandidateEvidence(request), "branch_evidence_attributes_unsupported")
	})

	t.Run("same-content remains object-exact under custom attributes", func(t *testing.T) {
		root := newBranchEvidenceRepository(t)
		writeBranchEvidenceFile(t, root, ".gitattributes", []byte("*.md diff=hostile\n"))
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("same\x00content\n"))
		commitBranchEvidenceAll(t, root, "base exact content")
		runGitTest(t, root, "branch", "history")
		writeBranchEvidenceFile(t, root, "target.txt", []byte("advance\n"))
		commitBranchEvidenceAll(t, root, "advance target")
		target := runGitTest(t, root, "rev-parse", "HEAD")
		result := RecomputeBranchCandidateEvidence(branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath))
		if result.ProposedDisposition != BranchDispositionSameContentNonOwner {
			t.Fatalf("exact object/content identity should remain sufficient: %#v", result.Blockers)
		}
	})

	t.Run("NUL blobs block incorporated proof", func(t *testing.T) {
		root := newBranchEvidenceRepository(t)
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("base\x00value\n"))
		commitBranchEvidenceAll(t, root, "binary base")
		target, incorporated := completeBranchEvidenceSquashHistory(t, root, []byte("middle\x00value\n"), []byte("middle\x00value\n"))
		request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)
		request.IncorporatedCommit = incorporated
		assertBranchEvidenceBlocker(t, RecomputeBranchCandidateEvidence(request), "branch_evidence_binary_unsupported")
	})

	t.Run("invalid UTF-8 blobs block incorporated proof", func(t *testing.T) {
		root := newBranchEvidenceRepository(t)
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("base\n"))
		commitBranchEvidenceAll(t, root, "base")
		target, incorporated := completeBranchEvidenceSquashHistory(t, root, []byte{0xff, '\n'}, []byte{0xff, '\n'})
		request := branchEvidenceRequest(root, target, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)
		request.IncorporatedCommit = incorporated
		assertBranchEvidenceBlocker(t, RecomputeBranchCandidateEvidence(request), "branch_evidence_utf8_unsupported")
	})

	t.Run("unsupported Git mode blocks exact proof", func(t *testing.T) {
		root := newBranchEvidenceRepository(t)
		target := filepath.Join(root, filepath.FromSlash(branchEvidencePlanPath))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("target.md", target); err != nil {
			t.Fatal(err)
		}
		commitBranchEvidenceAll(t, root, "symlink candidate")
		runGitTest(t, root, "branch", "history")
		writeBranchEvidenceFile(t, root, "target.txt", []byte("advance\n"))
		commitBranchEvidenceAll(t, root, "advance target")
		revision := runGitTest(t, root, "rev-parse", "HEAD")
		assertBranchEvidenceBlocker(t, RecomputeBranchCandidateEvidence(branchEvidenceRequest(root, revision, "refs/heads/history", "local:refs/heads/history", branchEvidencePlanPath)), "branch_evidence_git_mode_unsupported")
	})
}

func TestBranchEvidenceInventoryBlocksAmbiguousMergeBaseBeforePathEnumeration(t *testing.T) {
	root := newBranchEvidenceRepository(t)
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("base\n"))
	commitBranchEvidenceAll(t, root, "root")
	rootCommit := runGitTest(t, root, "rev-parse", "HEAD")
	tree := runGitTest(t, root, "rev-parse", "HEAD^{tree}")
	left := runGitTest(t, root, "commit-tree", tree, "-p", rootCommit, "-m", "left base")
	right := runGitTest(t, root, "commit-tree", tree, "-p", rootCommit, "-m", "right base")
	target := runGitTest(t, root, "commit-tree", tree, "-p", left, "-p", right, "-m", "target merge")
	tip := runGitTest(t, root, "commit-tree", tree, "-p", right, "-p", left, "-m", "history merge")
	runGitTest(t, root, "update-ref", "refs/heads/main", target)
	runGitTest(t, root, "update-ref", "refs/heads/history", tip)
	settings := defaultRepositorySettings()
	settings.DiscoveryVersion = DiscoveryV2
	settings.RepositoryID = "example-repository"
	settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	touches, problems := activeBranchTouchesWithSettings(root, settings)
	if len(touches) != 0 || !problemCodePresent(problems, "branch_evidence_merge_base_ambiguous") {
		t.Fatalf("ambiguous merge base was not fail-closed before path enumeration: touches=%#v problems=%#v", touches, problems)
	}
}

func TestBranchEvidenceProofScopeFailureNeverFallsBackToCandidateOnly(t *testing.T) {
	root := newBranchEvidenceRepository(t)
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("base\n"))
	commitBranchEvidenceAll(t, root, "base")
	target := runGitTest(t, root, "rev-parse", "HEAD")
	paths, blocker := branchCandidateProofPaths(root, target, "refs/heads/missing-proof-scope", branchEvidencePlanPath)
	if blocker == nil || blocker.Code != "branch_evidence_scope_unavailable" || len(paths) != 0 {
		t.Fatalf("failed proof scope fell back: paths=%#v blocker=%#v", paths, blocker)
	}
}

func TestBranchEvidenceSanitizedAcceptanceCohorts(t *testing.T) {
	root := newBranchEvidenceRepository(t)
	ownerOnePath := "docs/repo/plans/owner-one/owner-one_implementation_doc.md"
	ownerTwoPath := "docs/repo/plans/owner-two/owner-two_implementation_doc.md"
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("base\n"))
	writeBranchEvidenceFile(t, root, ownerOnePath, branchEvidencePlanDocument("draft", "example.implementation.owner-one"))
	writeBranchEvidenceFile(t, root, ownerTwoPath, branchEvidencePlanDocument("draft", "example.implementation.owner-two"))
	commitBranchEvidenceAll(t, root, "acceptance base")
	base := runGitTest(t, root, "rev-parse", "HEAD")

	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("incorporated\n"))
	commitBranchEvidenceAll(t, root, "incorporate historical candidate transition")
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("current target\n"))
	commitBranchEvidenceAll(t, root, "advance reviewed target")
	evidenceRevision := runGitTest(t, root, "rev-parse", "HEAD")

	requests := make([]BranchEvidenceRequest, 0, 19)
	for index := 0; index < 12; index++ {
		name := fmt.Sprintf("same-content-%02d", index+1)
		runGitTest(t, root, "switch", "-c", name, base)
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("current target\n"))
		commitBranchEvidenceAll(t, root, "independent same-content candidate")
		requests = append(requests, branchEvidenceRequest(root, evidenceRevision, "refs/heads/"+name, "local:refs/heads/"+name, branchEvidencePlanPath))
		runGitTest(t, root, "switch", "main")
	}
	for index := 0; index < 6; index++ {
		name := fmt.Sprintf("incorporated-history-%02d", index+1)
		runGitTest(t, root, "switch", "-c", name, base)
		writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("incorporated\n"))
		commitBranchEvidenceAll(t, root, "independent incorporated candidate")
		requests = append(requests, branchEvidenceRequest(root, evidenceRevision, "refs/heads/"+name, "local:refs/heads/"+name, branchEvidencePlanPath))
		runGitTest(t, root, "switch", "main")
	}
	runGitTest(t, root, "switch", "-c", "active-owner-one", evidenceRevision)
	writeBranchEvidenceFile(t, root, ownerOnePath, branchEvidencePlanDocument("active", "example.implementation.owner-one"))
	commitBranchEvidenceAll(t, root, "continue first active owner")
	ownerRequest := branchEvidenceRequest(root, evidenceRevision, "refs/heads/active-owner-one", "local:refs/heads/active-owner-one", ownerOnePath)
	ownerRequest.ExpectedKind = KindImplementation
	requests = append(requests, ownerRequest)
	runGitTest(t, root, "switch", "main")

	counts := map[string]int{}
	for _, request := range requests {
		evidence := ProposeBranchCandidateEvidence(request)
		if len(evidence.Blockers) != 0 {
			t.Fatalf("acceptance evidence for %s blocked: %#v", request.LogicalRef, evidence.Blockers)
		}
		counts[evidence.ProposedDisposition]++
	}
	if counts[BranchDispositionSameContentNonOwner] != 12 || counts[BranchDispositionIncorporatedHistoryNonOwner] != 6 || counts[BranchDispositionActiveOwner] != 1 || len(requests) != 19 {
		t.Fatalf("scenario B disposition counts=%#v requests=%d", counts, len(requests))
	}

	runGitTest(t, root, "switch", "-c", "active-owner-two", evidenceRevision)
	writeBranchEvidenceFile(t, root, ownerTwoPath, branchEvidencePlanDocument("active", "example.implementation.owner-two"))
	commitBranchEvidenceAll(t, root, "continue second active owner")
	secondOwner := branchEvidenceRequest(root, evidenceRevision, "refs/heads/active-owner-two", "local:refs/heads/active-owner-two", ownerTwoPath)
	secondOwner.ExpectedKind = KindImplementation
	secondEvidence := ProposeBranchCandidateEvidence(secondOwner)
	if secondEvidence.ProposedDisposition != BranchDispositionActiveOwner || len(secondEvidence.Blockers) != 0 {
		t.Fatalf("scenario A second active owner was not retained: %#v", secondEvidence)
	}
}

func branchEvidencePlanDocument(status, id string) []byte {
	return []byte(fmt.Sprintf(`Last updated: 2026-08-09T10:00:00Z (UTC)
Created: 2026-08-09
Status: %s

# Document Header

## Sanitized Acceptance Implementation

<!-- BEGIN CODEHEART PLAN METADATA -->
%s
plan:
  schema_version: 1
  id: %s
  kind: implementation
  purpose: Exercise sanitized branch ownership acceptance evidence.
  first_cataloged: 2026-08-09T10:00:00Z
  catalog_metadata_updated: 2026-08-09T10:00:00Z
%s
<!-- END CODEHEART PLAN METADATA -->

# Section 1 - Foundation

Synthetic public-safe acceptance body.
`, status, "```yaml", id, "```"))
}

func newBranchEvidenceRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGitTest(t, root, "init", "-b", "main")
	runGitTest(t, root, "config", "user.name", "Branch Evidence Test")
	runGitTest(t, root, "config", "user.email", "branch-evidence@example.invalid")
	return root
}

func newBranchEvidenceSquashRepository(t *testing.T, branchContent, incorporatedContent []byte) (string, string, string) {
	t.Helper()
	root := newBranchEvidenceRepository(t)
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("base\n"))
	commitBranchEvidenceAll(t, root, "base")
	target, incorporated := completeBranchEvidenceSquashHistory(t, root, branchContent, incorporatedContent)
	return root, target, incorporated
}

func completeBranchEvidenceSquashHistory(t *testing.T, root string, branchContent, incorporatedContent []byte) (string, string) {
	t.Helper()
	runGitTest(t, root, "switch", "-c", "history")
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, branchContent)
	commitBranchEvidenceAll(t, root, "historical candidate transition")
	runGitTest(t, root, "switch", "main")
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, incorporatedContent)
	commitBranchEvidenceAll(t, root, "incorporated squash transition")
	incorporated := runGitTest(t, root, "rev-parse", "HEAD")
	writeBranchEvidenceFile(t, root, branchEvidencePlanPath, []byte("current target\n"))
	commitBranchEvidenceAll(t, root, "advance target candidate")
	return runGitTest(t, root, "rev-parse", "HEAD"), incorporated
}

func branchEvidenceRequest(root, target, ref, logicalRef, candidatePath string) BranchEvidenceRequest {
	return BranchEvidenceRequest{
		Root:               root,
		EvidenceScope:      "local",
		RepositoryID:       "example-repository",
		LogicalRef:         logicalRef,
		Ref:                ref,
		EvidenceRevision:   target,
		PolicyDigest:       strings.Repeat("a", 64),
		CandidateSetDigest: strings.Repeat("b", 64),
		CandidatePath:      candidatePath,
		Paths:              []string{candidatePath},
	}
}

func writeBranchEvidenceFile(t *testing.T, root, relative string, data []byte) {
	t.Helper()
	target := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitBranchEvidenceAll(t *testing.T, root, message string) {
	t.Helper()
	runGitTest(t, root, "add", "--all")
	runGitTest(t, root, "commit", "-m", message)
}

func assertBranchEvidenceBlocker(t *testing.T, result BranchCandidateEvidence, code string) {
	t.Helper()
	if result.ProposedDisposition != BranchDispositionBlocking || result.Proof != nil {
		t.Fatalf("expected fail-closed evidence, got disposition=%q proof=%#v", result.ProposedDisposition, result.Proof)
	}
	for _, blocker := range result.Blockers {
		if blocker.Code == code {
			return
		}
	}
	t.Fatalf("missing blocker %q in %#v", code, result.Blockers)
}
