package plancatalog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
)

const (
	BranchEvidenceAlgorithm = "git-candidate-proof-v1"

	BranchDispositionSameContentNonOwner         = "same-content-non-owner"
	BranchDispositionIncorporatedHistoryNonOwner = "incorporated-history-non-owner"
	BranchDispositionDeferredActiveOwner         = "deferred-active-owner"
	BranchDispositionActiveOwner                 = "active-owner"
	BranchDispositionBlocking                    = "blocking"

	BranchProofSameContentV1         = "same-content-v1"
	BranchProofIncorporatedHistoryV1 = "incorporated-history-v1"

	BranchPathPresent = "present"
	BranchPathAbsent  = "absent"

	MaxBranchEvidencePaths      = 256
	MaxBranchEvidenceDiffBytes  = 16 << 20
	maxBranchGitDiagnosticBytes = 64 << 10
)

var (
	branchGitVersionPattern   = regexp.MustCompile(`^git version ([0-9]+)\.([0-9]+)(?:\.([0-9]+))?`)
	branchRepositoryIDPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)
	branchLogicalRefPattern   = regexp.MustCompile(`^(?:local:refs/heads/|tracking:[a-z0-9](?:[a-z0-9._-]*[a-z0-9])?:refs/heads/|remote:[a-z0-9](?:[a-z0-9._-]*[a-z0-9])?:refs/heads/)[^\r\n]+$`)
)

// BranchReviewIdentity is the candidate-scoped review key. A review for one
// candidate never authorizes a disposition for another candidate touched by
// the same ref.
type BranchReviewIdentity struct {
	EvidenceScope string `json:"evidence_scope" yaml:"evidence_scope"`
	RepositoryID  string `json:"repository_id" yaml:"repository_id"`
	LogicalRef    string `json:"logical_ref" yaml:"logical_ref"`
	CandidatePath string `json:"candidate_path" yaml:"candidate_path"`
}

type BranchPathState struct {
	State    string  `json:"state" yaml:"state"`
	Mode     GitMode `json:"mode,omitempty" yaml:"mode,omitempty"`
	ObjectID string  `json:"object_id,omitempty" yaml:"object_id,omitempty"`
	SHA256   string  `json:"sha256,omitempty" yaml:"sha256,omitempty"`
}

type BranchPathTransition struct {
	Path   string          `json:"path" yaml:"path"`
	Before BranchPathState `json:"before" yaml:"before"`
	After  BranchPathState `json:"after" yaml:"after"`
}

// BranchPathAttributes records the effective attributes which can change
// Git's textual diff semantics. Values are Git's stable check-attr spellings.
type BranchPathAttributes struct {
	Path                string `json:"path" yaml:"path"`
	Diff                string `json:"diff" yaml:"diff"`
	Text                string `json:"text" yaml:"text"`
	WorkingTreeEncoding string `json:"working_tree_encoding" yaml:"working_tree_encoding"`
	Binary              string `json:"binary" yaml:"binary"`
}

type BranchAttributeFileState struct {
	Path  string          `json:"path" yaml:"path"`
	State BranchPathState `json:"state" yaml:"state"`
}

// BranchAttributeSnapshot binds both the effective attribute result and every
// repository .gitattributes blob that can govern one of the reviewed paths.
type BranchAttributeSnapshot struct {
	Revision       string                     `json:"revision" yaml:"revision"`
	PathAttributes []BranchPathAttributes     `json:"path_attributes" yaml:"path_attributes"`
	AttributeFiles []BranchAttributeFileState `json:"attribute_files" yaml:"attribute_files"`
	Digest         string                     `json:"digest" yaml:"digest"`
}

type BranchEvidenceProof struct {
	Kind                         string           `json:"kind" yaml:"kind"`
	TargetPathState              *BranchPathState `json:"target_path_state,omitempty" yaml:"target_path_state,omitempty"`
	RefPathState                 *BranchPathState `json:"ref_path_state,omitempty" yaml:"ref_path_state,omitempty"`
	TransitionDigest             string           `json:"transition_digest,omitempty" yaml:"transition_digest,omitempty"`
	StablePatchID                string           `json:"stable_patch_id,omitempty" yaml:"stable_patch_id,omitempty"`
	IncorporatedParent           string           `json:"incorporated_parent,omitempty" yaml:"incorporated_parent,omitempty"`
	IncorporatedCommit           string           `json:"incorporated_commit,omitempty" yaml:"incorporated_commit,omitempty"`
	IncorporatedTransitionDigest string           `json:"incorporated_transition_digest,omitempty" yaml:"incorporated_transition_digest,omitempty"`
	IncorporatedStablePatchID    string           `json:"incorporated_stable_patch_id,omitempty" yaml:"incorporated_stable_patch_id,omitempty"`
}

type BranchEvidenceBlocker struct {
	Code    string `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
	Path    string `json:"path,omitempty" yaml:"path,omitempty"`
	Ref     string `json:"ref,omitempty" yaml:"ref,omitempty"`
}

// BranchCandidateEvidence is a deterministic recomputation result suitable
// for comparison with a reviewed ledger entry. Blockers always force the
// proposed disposition to blocking and remove any clearance proof.
type BranchCandidateEvidence struct {
	Algorithm           string                    `json:"algorithm" yaml:"algorithm"`
	GitVersion          string                    `json:"git_version,omitempty" yaml:"git_version,omitempty"`
	Identity            BranchReviewIdentity      `json:"identity" yaml:"identity"`
	EvidenceRevision    string                    `json:"evidence_revision" yaml:"evidence_revision"`
	PolicyDigest        string                    `json:"policy_digest" yaml:"policy_digest"`
	CandidateSetDigest  string                    `json:"candidate_set_digest" yaml:"candidate_set_digest"`
	RefTip              string                    `json:"ref_tip,omitempty" yaml:"ref_tip,omitempty"`
	MergeBase           string                    `json:"merge_base,omitempty" yaml:"merge_base,omitempty"`
	PathTransitions     []BranchPathTransition    `json:"path_transitions,omitempty" yaml:"path_transitions,omitempty"`
	AttributeSnapshots  []BranchAttributeSnapshot `json:"attribute_snapshots,omitempty" yaml:"attribute_snapshots,omitempty"`
	ProposedDisposition string                    `json:"proposed_disposition" yaml:"proposed_disposition"`
	Proof               *BranchEvidenceProof      `json:"proof,omitempty" yaml:"proof,omitempty"`
	OwnerTipCandidate   *OwnerTipCandidate        `json:"owner_tip_candidate,omitempty" yaml:"owner_tip_candidate,omitempty"`
	Blockers            []BranchEvidenceBlocker   `json:"blockers,omitempty" yaml:"blockers,omitempty"`
	Digest              string                    `json:"digest" yaml:"digest"`
}

type BranchEvidenceRequest struct {
	Root               string
	EvidenceScope      string
	RepositoryID       string
	LogicalRef         string
	Ref                string
	EvidenceRevision   string
	PolicyDigest       string
	CandidateSetDigest string
	CandidatePath      string
	Paths              []string
	ExpectedRefTip     string
	IncorporatedCommit string
	ExpectedKind       Kind
}

// RecomputeBranchCandidateEvidence builds a fresh, object-only proof. It does
// not trust a previously serialized proof and does not fetch missing objects.
func RecomputeBranchCandidateEvidence(request BranchEvidenceRequest) BranchCandidateEvidence {
	evidence := BranchCandidateEvidence{
		Algorithm: BranchEvidenceAlgorithm,
		Identity: BranchReviewIdentity{
			EvidenceScope: request.EvidenceScope,
			RepositoryID:  request.RepositoryID,
			LogicalRef:    request.LogicalRef,
			CandidatePath: request.CandidatePath,
		},
		EvidenceRevision:    request.EvidenceRevision,
		PolicyDigest:        request.PolicyDigest,
		CandidateSetDigest:  request.CandidateSetDigest,
		ProposedDisposition: BranchDispositionBlocking,
	}
	block := func(blocker *BranchEvidenceBlocker) BranchCandidateEvidence {
		if blocker != nil {
			evidence.Blockers = append(evidence.Blockers, *blocker)
		}
		evidence.ProposedDisposition = BranchDispositionBlocking
		evidence.Proof = nil
		evidence.Digest = CanonicalBranchEvidenceDigest(evidence)
		return evidence
	}

	if blocker := validateBranchEvidenceRequest(request); blocker != nil {
		return block(blocker)
	}
	paths, blocker := normalizeBranchEvidencePaths(request.CandidatePath, request.Paths)
	if blocker != nil {
		return block(blocker)
	}
	runner, blocker := newBranchGitRunner(request.Root)
	if blocker != nil {
		return block(blocker)
	}
	defer runner.close()

	version, blocker := runner.checkVersion()
	if blocker != nil {
		return block(blocker)
	}
	evidence.GitVersion = version
	if blocker := runner.checkRepositoryIsolation(); blocker != nil {
		return block(blocker)
	}
	if blocker := runner.requireCommit(request.EvidenceRevision, "evidence revision"); blocker != nil {
		return block(blocker)
	}
	tip, blocker := runner.resolveRefTip(request.Ref)
	if blocker != nil {
		return block(blocker)
	}
	evidence.RefTip = tip
	if request.ExpectedRefTip != "" && request.ExpectedRefTip != tip {
		return block(&BranchEvidenceBlocker{Code: "branch_evidence_ref_moved", Message: fmt.Sprintf("reviewed ref tip %s moved to %s", request.ExpectedRefTip, tip), Ref: request.Ref})
	}
	base, blocker := runner.resolveUniqueMergeBase(request.EvidenceRevision, tip)
	if blocker != nil {
		return block(blocker)
	}
	evidence.MergeBase = base
	transitions, blocker := runner.buildTransitions(base, tip, paths)
	if blocker != nil {
		return block(blocker)
	}
	evidence.PathTransitions = transitions
	targetState, _, blocker := runner.readPathState(request.EvidenceRevision, request.CandidatePath)
	if blocker != nil {
		return block(blocker)
	}
	refState, _, blocker := runner.readPathState(tip, request.CandidatePath)
	if blocker != nil {
		return block(blocker)
	}
	if ExactBranchPathIdentity(targetState, refState) {
		evidence.ProposedDisposition = BranchDispositionSameContentNonOwner
		evidence.Proof = &BranchEvidenceProof{Kind: BranchProofSameContentV1, TargetPathState: &targetState, RefPathState: &refState}
		if blocker := runner.verifyRefTip(request.Ref, tip); blocker != nil {
			return block(blocker)
		}
		evidence.Digest = CanonicalBranchEvidenceDigest(evidence)
		return evidence
	}
	if request.IncorporatedCommit == "" {
		return block(&BranchEvidenceBlocker{Code: "branch_evidence_non_owner_unproven", Message: "candidate differs at the reviewed target and no incorporated commit was supplied", Path: request.CandidatePath, Ref: request.Ref})
	}
	proof, snapshots, blocker := runner.proveIncorporated(request.EvidenceRevision, base, tip, request.IncorporatedCommit, paths, transitions)
	if blocker != nil {
		return block(blocker)
	}
	evidence.AttributeSnapshots = snapshots
	evidence.ProposedDisposition = BranchDispositionIncorporatedHistoryNonOwner
	evidence.Proof = &proof
	if blocker := runner.verifyRefTip(request.Ref, tip); blocker != nil {
		return block(blocker)
	}
	evidence.Digest = CanonicalBranchEvidenceDigest(evidence)
	return evidence
}

func ResolveBranchRefTip(root, ref string) (string, *BranchEvidenceBlocker) {
	runner, blocker := newBranchGitRunner(root)
	if blocker != nil {
		return "", blocker
	}
	defer runner.close()
	if _, blocker = runner.checkVersion(); blocker != nil {
		return "", blocker
	}
	if blocker = runner.checkRepositoryIsolation(); blocker != nil {
		return "", blocker
	}
	return runner.resolveRefTip(ref)
}

func VerifyBranchRefTip(root, ref, expected string) *BranchEvidenceBlocker {
	runner, blocker := newBranchGitRunner(root)
	if blocker != nil {
		return blocker
	}
	defer runner.close()
	if _, blocker = runner.checkVersion(); blocker != nil {
		return blocker
	}
	if blocker = runner.checkRepositoryIsolation(); blocker != nil {
		return blocker
	}
	return runner.verifyRefTip(ref, expected)
}

func ResolveBranchUniqueMergeBase(root, target, tip string) (string, *BranchEvidenceBlocker) {
	runner, blocker := newBranchGitRunner(root)
	if blocker != nil {
		return "", blocker
	}
	defer runner.close()
	if _, blocker = runner.checkVersion(); blocker != nil {
		return "", blocker
	}
	if blocker = runner.checkRepositoryIsolation(); blocker != nil {
		return "", blocker
	}
	return runner.resolveUniqueMergeBase(target, tip)
}

func ReadBranchPathState(root, revision, candidatePath string) (BranchPathState, *BranchEvidenceBlocker) {
	runner, blocker := newBranchGitRunner(root)
	if blocker != nil {
		return BranchPathState{}, blocker
	}
	defer runner.close()
	if _, blocker = runner.checkVersion(); blocker != nil {
		return BranchPathState{}, blocker
	}
	if blocker = runner.checkRepositoryIsolation(); blocker != nil {
		return BranchPathState{}, blocker
	}
	state, _, blocker := runner.readPathState(revision, candidatePath)
	return state, blocker
}

func BuildBranchPathTransitions(root, before, after string, paths []string) ([]BranchPathTransition, *BranchEvidenceBlocker) {
	normalized, blocker := normalizeBranchEvidencePaths("", paths)
	if blocker != nil {
		return nil, blocker
	}
	runner, blocker := newBranchGitRunner(root)
	if blocker != nil {
		return nil, blocker
	}
	defer runner.close()
	if _, blocker = runner.checkVersion(); blocker != nil {
		return nil, blocker
	}
	if blocker = runner.checkRepositoryIsolation(); blocker != nil {
		return nil, blocker
	}
	return runner.buildTransitions(before, after, normalized)
}

func ExactBranchPathIdentity(target, ref BranchPathState) bool {
	return validateBranchPathState("", target) == nil && validateBranchPathState("", ref) == nil &&
		target.State == BranchPathPresent && ref.State == BranchPathPresent &&
		target.Mode == ref.Mode && target.ObjectID == ref.ObjectID && target.SHA256 == ref.SHA256
}

// ProveIncorporatedBranchTransition exposes the collision-resistant proof for
// callers which already resolved the target, merge base, and tip.
func ProveIncorporatedBranchTransition(root, target, mergeBase, refTip, incorporatedCommit string, paths []string) (BranchEvidenceProof, []BranchAttributeSnapshot, *BranchEvidenceBlocker) {
	normalized, blocker := normalizeBranchEvidencePaths("", paths)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	runner, blocker := newBranchGitRunner(root)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	defer runner.close()
	if _, blocker = runner.checkVersion(); blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	if blocker = runner.checkRepositoryIsolation(); blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	resolvedBase, blocker := runner.resolveUniqueMergeBase(target, refTip)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	if resolvedBase != mergeBase {
		return BranchEvidenceProof{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_merge_base_moved", Message: "supplied merge base no longer matches the unique target/ref merge base"}
	}
	transitions, blocker := runner.buildTransitions(mergeBase, refTip, normalized)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	return runner.proveIncorporated(target, mergeBase, refTip, incorporatedCommit, normalized, transitions)
}

func CheckBranchEvidenceGitVersion(root string) (string, *BranchEvidenceBlocker) {
	runner, blocker := newBranchGitRunner(root)
	if blocker != nil {
		return "", blocker
	}
	defer runner.close()
	return runner.checkVersion()
}

func CanonicalBranchTransitionDigest(transitions []BranchPathTransition) (string, *BranchEvidenceBlocker) {
	canonical := append([]BranchPathTransition(nil), transitions...)
	sort.Slice(canonical, func(left, right int) bool { return canonical[left].Path < canonical[right].Path })
	paths := make([]string, len(canonical))
	for index := range canonical {
		paths[index] = canonical[index].Path
	}
	if _, blocker := normalizeBranchEvidencePaths("", paths); blocker != nil {
		return "", blocker
	}
	for index := range canonical {
		if blocker := validateBranchPathState(canonical[index].Path, canonical[index].Before); blocker != nil {
			return "", blocker
		}
		if blocker := validateBranchPathState(canonical[index].Path, canonical[index].After); blocker != nil {
			return "", blocker
		}
	}
	return canonicalJSONDigest(canonical), nil
}

func CanonicalBranchEvidenceDigest(evidence BranchCandidateEvidence) string {
	canonical := evidence
	// The concrete Git implementation is recorded for diagnostics. The
	// versioned proof algorithm and minimum supported Git version own proof
	// semantics, so identical evidence stays portable across admitted runners.
	canonical.GitVersion = ""
	canonical.Digest = ""
	canonical.PathTransitions = append([]BranchPathTransition(nil), evidence.PathTransitions...)
	sort.Slice(canonical.PathTransitions, func(left, right int) bool {
		return canonical.PathTransitions[left].Path < canonical.PathTransitions[right].Path
	})
	canonical.AttributeSnapshots = append([]BranchAttributeSnapshot(nil), evidence.AttributeSnapshots...)
	sort.Slice(canonical.AttributeSnapshots, func(left, right int) bool {
		return canonical.AttributeSnapshots[left].Revision < canonical.AttributeSnapshots[right].Revision
	})
	canonical.Blockers = append([]BranchEvidenceBlocker(nil), evidence.Blockers...)
	sort.Slice(canonical.Blockers, func(left, right int) bool {
		leftKey := canonical.Blockers[left].Code + "\x00" + canonical.Blockers[left].Path + "\x00" + canonical.Blockers[left].Ref + "\x00" + canonical.Blockers[left].Message
		rightKey := canonical.Blockers[right].Code + "\x00" + canonical.Blockers[right].Path + "\x00" + canonical.Blockers[right].Ref + "\x00" + canonical.Blockers[right].Message
		return leftKey < rightKey
	})
	return canonicalJSONDigest(canonical)
}

// EquivalentBranchCandidateEvidence compares proof-bearing evidence while
// ignoring only the transport-specific scope/ref, diagnostic Git version,
// and derived digest. Repository and candidate identity must remain exact.
// Only blocker-free evidence can be coalesced.
func EquivalentBranchCandidateEvidence(left, right BranchCandidateEvidence) bool {
	if len(left.Blockers) != 0 || len(right.Blockers) != 0 {
		return false
	}
	if left.Identity.RepositoryID != right.Identity.RepositoryID || left.Identity.CandidatePath != right.Identity.CandidatePath {
		return false
	}
	left.Identity.EvidenceScope, right.Identity.EvidenceScope = "", ""
	left.Identity.LogicalRef, right.Identity.LogicalRef = "", ""
	left.GitVersion, right.GitVersion = "", ""
	left.Digest, right.Digest = "", ""
	return canonicalJSONDigest(left) == canonicalJSONDigest(right)
}

func validateBranchEvidenceRequest(request BranchEvidenceRequest) *BranchEvidenceBlocker {
	for _, field := range []struct{ label, value string }{
		{label: "evidence scope", value: request.EvidenceScope},
		{label: "repository ID", value: request.RepositoryID},
		{label: "logical ref", value: request.LogicalRef},
	} {
		label, value := field.label, field.value
		if value == "" || len(value) > 1024 || !utf8.ValidString(value) || strings.ContainsRune(value, 0) {
			return &BranchEvidenceBlocker{Code: "branch_evidence_identity_invalid", Message: label + " is missing or malformed"}
		}
		for _, char := range value {
			if unicode.IsControl(char) {
				return &BranchEvidenceBlocker{Code: "branch_evidence_identity_invalid", Message: label + " contains a control character"}
			}
		}
	}
	if request.EvidenceScope != "local" && request.EvidenceScope != "remote-aware" {
		return &BranchEvidenceBlocker{Code: "branch_evidence_identity_invalid", Message: "evidence scope must be local or remote-aware"}
	}
	if !branchRepositoryIDPattern.MatchString(request.RepositoryID) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_identity_invalid", Message: "repository ID is not a normalized public identifier"}
	}
	if !branchLogicalRefPattern.MatchString(request.LogicalRef) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_identity_invalid", Message: "logical ref is not a normalized local, tracking, or remote identity", Ref: request.LogicalRef}
	}
	if problem := validateGitPath(request.CandidatePath); problem != nil {
		return &BranchEvidenceBlocker{Code: "branch_evidence_path_invalid", Message: problem.Message, Path: request.CandidatePath}
	}
	if !validObjectID(request.EvidenceRevision) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_revision_invalid", Message: "evidence revision must be a full lowercase Git object ID"}
	}
	for _, field := range []struct{ label, value string }{
		{label: "policy digest", value: request.PolicyDigest},
		{label: "candidate-set digest", value: request.CandidateSetDigest},
	} {
		label, value := field.label, field.value
		if !validSHA256(value) {
			return &BranchEvidenceBlocker{Code: "branch_evidence_binding_invalid", Message: label + " must be a lowercase SHA-256 digest"}
		}
	}
	if request.ExpectedRefTip != "" && !validObjectID(request.ExpectedRefTip) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_ref_tip_invalid", Message: "expected ref tip must be a full lowercase Git object ID", Ref: request.Ref}
	}
	if request.IncorporatedCommit != "" && !validObjectID(request.IncorporatedCommit) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_incorporated_commit_invalid", Message: "incorporated commit must be a full lowercase Git object ID"}
	}
	// An immutable reviewed tip may be replayed after its named owner ref was
	// integrated or removed. This is only admissible when the request binds the
	// object name to the exact reviewed tip; ordinary discovery still requires
	// a full refs/ identity so movement is detectable.
	if validObjectID(request.Ref) && request.Ref == request.ExpectedRefTip {
		return nil
	}
	return validateBranchRef(request.Ref)
}

func validateBranchRef(ref string) *BranchEvidenceBlocker {
	if ref == "" || len(ref) > 1024 || !strings.HasPrefix(ref, "refs/") || !utf8.ValidString(ref) || strings.ContainsRune(ref, 0) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_ref_invalid", Message: "ref must be a full, valid refs/ name", Ref: ref}
	}
	for _, char := range ref {
		if unicode.IsControl(char) {
			return &BranchEvidenceBlocker{Code: "branch_evidence_ref_invalid", Message: "ref contains a control character", Ref: ref}
		}
	}
	return nil
}

func normalizeBranchEvidencePaths(candidatePath string, paths []string) ([]string, *BranchEvidenceBlocker) {
	if len(paths) == 0 && candidatePath != "" {
		paths = []string{candidatePath}
	}
	if len(paths) == 0 {
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_paths_missing", Message: "at least one candidate path is required"}
	}
	if len(paths) > MaxBranchEvidencePaths {
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_path_limit", Message: fmt.Sprintf("candidate proof exceeds the %d-path limit", MaxBranchEvidencePaths)}
	}
	normalized := append([]string(nil), paths...)
	portable := map[string]string{}
	fold := cases.Fold()
	containsCandidate := candidatePath == ""
	for _, candidate := range normalized {
		if problem := validateGitPath(candidate); problem != nil {
			return nil, &BranchEvidenceBlocker{Code: "branch_evidence_path_invalid", Message: problem.Message, Path: candidate}
		}
		if candidate == candidatePath {
			containsCandidate = true
		}
		key := fold.String(candidate)
		if prior, exists := portable[key]; exists {
			return nil, &BranchEvidenceBlocker{Code: "branch_evidence_path_collision", Message: fmt.Sprintf("candidate paths collide under portable comparison: %s and %s", prior, candidate), Path: candidate}
		}
		portable[key] = candidate
	}
	if !containsCandidate {
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_candidate_path_missing", Message: "candidate path is absent from the proof path set", Path: candidatePath}
	}
	sort.Strings(normalized)
	return normalized, nil
}

func validateBranchPathState(candidatePath string, state BranchPathState) *BranchEvidenceBlocker {
	switch state.State {
	case BranchPathAbsent:
		if state.Mode != "" || state.ObjectID != "" || state.SHA256 != "" {
			return &BranchEvidenceBlocker{Code: "branch_evidence_path_state_invalid", Message: "absent path state carries object fields", Path: candidatePath}
		}
	case BranchPathPresent:
		if state.Mode != GitModeRegular && state.Mode != GitModeExecutable {
			return &BranchEvidenceBlocker{Code: "branch_evidence_git_mode_unsupported", Message: "only regular and executable blobs are supported", Path: candidatePath}
		}
		if !validObjectID(state.ObjectID) || !validSHA256(state.SHA256) {
			return &BranchEvidenceBlocker{Code: "branch_evidence_path_state_invalid", Message: "present path state has an invalid object or content digest", Path: candidatePath}
		}
	default:
		return &BranchEvidenceBlocker{Code: "branch_evidence_path_state_invalid", Message: "path state must be present or absent", Path: candidatePath}
	}
	return nil
}

type branchGitRunner struct {
	root      string
	emptyFile string
}

func newBranchGitRunner(root string) (*branchGitRunner, *BranchEvidenceBlocker) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_repository_unavailable", Message: "cannot resolve repository root: " + err.Error()}
	}
	info, err := os.Stat(absolute)
	if err != nil || !info.IsDir() {
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_repository_unavailable", Message: "repository root is not an accessible directory"}
	}
	empty, err := os.CreateTemp("", "codeheart-branch-evidence-empty-*")
	if err != nil {
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_runtime_unavailable", Message: "cannot create isolated Git configuration: " + err.Error()}
	}
	name := empty.Name()
	if closeErr := empty.Close(); closeErr != nil {
		_ = os.Remove(name)
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_runtime_unavailable", Message: "cannot close isolated Git configuration: " + closeErr.Error()}
	}
	return &branchGitRunner{root: absolute, emptyFile: name}, nil
}

func (runner *branchGitRunner) close() { _ = os.Remove(runner.emptyFile) }

func (runner *branchGitRunner) command(args ...string) *exec.Cmd {
	global := []string{
		"--literal-pathspecs",
		"-c", "color.ui=false",
		"-c", "core.quotePath=true",
		"-c", "diff.algorithm=myers",
		"-c", "diff.renames=false",
		"-c", "diff.external=",
		"-c", "core.attributesFile=" + runner.emptyFile,
		"-C", runner.root,
	}
	command := exec.Command("git", append(global, args...)...)
	command.Env = branchGitEnvironment(runner.emptyFile)
	return command
}

func branchGitEnvironment(emptyConfig string) []string {
	environment := make([]string, 0, len(os.Environ())+8)
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		upper := strings.ToUpper(key)
		if strings.HasPrefix(upper, "GIT_") || upper == "LC_ALL" || upper == "LANG" {
			continue
		}
		environment = append(environment, entry)
	}
	environment = append(environment,
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL="+emptyConfig,
		"GIT_TERMINAL_PROMPT=0",
		"GIT_NO_LAZY_FETCH=1",
		"GIT_NO_REPLACE_OBJECTS=1",
		"GIT_LITERAL_PATHSPECS=1",
		"LC_ALL=C",
		"LANG=C",
	)
	return environment
}

func (runner *branchGitRunner) run(limit int, input []byte, args ...string) ([]byte, *BranchEvidenceBlocker) {
	command := runner.command(args...)
	if input != nil {
		command.Stdin = bytes.NewReader(input)
	}
	stdout := newBranchBoundedBuffer(limit)
	stderr := newBranchBoundedBuffer(maxBranchGitDiagnosticBytes)
	command.Stdout = stdout
	command.Stderr = stderr
	err := command.Run()
	if stdout.exceeded || stderr.exceeded {
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_output_limit", Message: "Git output exceeded a deterministic evidence limit"}
	}
	if err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_git_failed", Message: message}
	}
	return stdout.Bytes(), nil
}

func (runner *branchGitRunner) checkVersion() (string, *BranchEvidenceBlocker) {
	output, blocker := runner.run(1024, nil, "--version")
	if blocker != nil {
		blocker.Code = "branch_evidence_git_version_unsupported"
		return "", blocker
	}
	return parseBranchEvidenceGitVersion(strings.TrimSpace(string(output)))
}

func parseBranchEvidenceGitVersion(output string) (string, *BranchEvidenceBlocker) {
	match := branchGitVersionPattern.FindStringSubmatch(output)
	if match == nil {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_git_version_unsupported", Message: "Git version is unparsable"}
	}
	major, _ := strconv.Atoi(match[1])
	minor, _ := strconv.Atoi(match[2])
	patchLevel := 0
	if match[3] != "" {
		patchLevel, _ = strconv.Atoi(match[3])
	}
	version := fmt.Sprintf("%d.%d.%d", major, minor, patchLevel)
	if major < 2 || (major == 2 && minor < 43) {
		return version, &BranchEvidenceBlocker{Code: "branch_evidence_git_version_unsupported", Message: "git-candidate-proof-v1 requires Git >= 2.43.0"}
	}
	return version, nil
}

func (runner *branchGitRunner) checkRepositoryIsolation() *BranchEvidenceBlocker {
	gitPath, blocker := runner.run(4096, nil, "rev-parse", "--git-path", "objects/info/alternates")
	if blocker != nil {
		blocker.Code = "branch_evidence_repository_unavailable"
		return blocker
	}
	alternates := strings.TrimSpace(string(gitPath))
	if !filepath.IsAbs(alternates) {
		alternates = filepath.Join(runner.root, alternates)
	}
	if info, err := os.Lstat(alternates); err == nil {
		if !info.Mode().IsRegular() || info.Size() > 0 {
			return &BranchEvidenceBlocker{Code: "branch_evidence_alternates_unsupported", Message: "repository object alternates are not supported"}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_repository_unavailable", Message: "cannot inspect repository object alternates"}
	}
	infoAttributes, blocker := runner.run(4096, nil, "rev-parse", "--git-path", "info/attributes")
	if blocker != nil {
		blocker.Code = "branch_evidence_repository_unavailable"
		return blocker
	}
	attributePath := strings.TrimSpace(string(infoAttributes))
	if !filepath.IsAbs(attributePath) {
		attributePath = filepath.Join(runner.root, attributePath)
	}
	if info, err := os.Lstat(attributePath); err == nil {
		if !info.Mode().IsRegular() || info.Size() > 0 {
			return &BranchEvidenceBlocker{Code: "branch_evidence_attributes_unsupported", Message: "non-empty repository-local info/attributes is not supported"}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_repository_unavailable", Message: "cannot inspect repository-local attributes"}
	}
	command := runner.command("config", "--local", "--get-regexp", `^remote\..*\.promisor$`)
	stdout := newBranchBoundedBuffer(4096)
	stderr := newBranchBoundedBuffer(maxBranchGitDiagnosticBytes)
	command.Stdout = stdout
	command.Stderr = stderr
	err := command.Run()
	if err == nil && strings.TrimSpace(stdout.String()) != "" {
		return &BranchEvidenceBlocker{Code: "branch_evidence_promisor_unsupported", Message: "promisor remotes are not supported for offline branch evidence"}
	}
	var exitErr *exec.ExitError
	if err != nil && (!errors.As(err, &exitErr) || exitErr.ExitCode() != 1) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_repository_unavailable", Message: "cannot inspect promisor configuration"}
	}
	return nil

}

func (runner *branchGitRunner) requireCommit(objectID, label string) *BranchEvidenceBlocker {
	if !validObjectID(objectID) {
		return &BranchEvidenceBlocker{Code: "branch_evidence_object_invalid", Message: label + " is not a full lowercase object ID"}
	}
	output, blocker := runner.run(128, nil, "cat-file", "-t", objectID)
	if blocker != nil || strings.TrimSpace(string(output)) != "commit" {
		return &BranchEvidenceBlocker{Code: "branch_evidence_object_missing", Message: label + " is not an available commit object"}
	}
	return nil
}

func (runner *branchGitRunner) resolveRefTip(ref string) (string, *BranchEvidenceBlocker) {
	if validObjectID(ref) {
		if blocker := runner.requireCommit(ref, "reviewed tip"); blocker != nil {
			blocker.Ref = ref
			return "", blocker
		}
		return ref, nil
	}
	if blocker := validateBranchRef(ref); blocker != nil {
		return "", blocker
	}
	if _, blocker := runner.run(128, nil, "check-ref-format", ref); blocker != nil {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_ref_invalid", Message: "ref fails Git ref-format validation", Ref: ref}
	}
	output, blocker := runner.run(256, nil, "show-ref", "--verify", "--hash", ref)
	if blocker != nil {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_ref_missing", Message: "reviewed ref is missing or unavailable", Ref: ref}
	}
	lines := nonEmptyLines(output)
	if len(lines) != 1 || !validObjectID(lines[0]) {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_ref_ambiguous", Message: "reviewed ref did not resolve to exactly one full object ID", Ref: ref}
	}
	if blocker := runner.requireCommit(lines[0], "ref tip"); blocker != nil {
		blocker.Ref = ref
		return "", blocker
	}
	return lines[0], nil
}

func (runner *branchGitRunner) verifyRefTip(ref, expected string) *BranchEvidenceBlocker {
	observed, blocker := runner.resolveRefTip(ref)
	if blocker != nil {
		return blocker
	}
	if observed != expected {
		return &BranchEvidenceBlocker{Code: "branch_evidence_ref_moved", Message: fmt.Sprintf("reviewed ref tip %s moved to %s", expected, observed), Ref: ref}
	}
	return nil
}

func (runner *branchGitRunner) resolveUniqueMergeBase(target, tip string) (string, *BranchEvidenceBlocker) {
	if blocker := runner.requireCommit(target, "target revision"); blocker != nil {
		return "", blocker
	}
	if blocker := runner.requireCommit(tip, "ref tip"); blocker != nil {
		return "", blocker
	}
	output, blocker := runner.run(512, nil, "merge-base", "--all", target, tip)
	if blocker != nil {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_merge_base_missing", Message: "target and ref do not have an available merge base"}
	}
	bases := nonEmptyLines(output)
	if len(bases) != 1 || !validObjectID(bases[0]) {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_merge_base_ambiguous", Message: fmt.Sprintf("expected one merge base, found %d", len(bases))}
	}
	return bases[0], nil
}

func (runner *branchGitRunner) readPathState(revision, candidatePath string) (BranchPathState, []byte, *BranchEvidenceBlocker) {
	if !validObjectID(revision) {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_revision_invalid", Message: "path-state revision must be a full lowercase object ID", Path: candidatePath}
	}
	if problem := validateGitPath(candidatePath); problem != nil {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_path_invalid", Message: problem.Message, Path: candidatePath}
	}
	output, blocker := runner.run(4096, nil, "ls-tree", "-z", revision, "--", candidatePath)
	if blocker != nil {
		blocker.Code = "branch_evidence_path_state_unavailable"
		blocker.Path = candidatePath
		return BranchPathState{}, nil, blocker
	}
	records := splitNUL(output)
	if len(records) == 0 {
		return BranchPathState{State: BranchPathAbsent}, nil, nil
	}
	if len(records) != 1 {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_path_collision", Message: "literal path resolved to multiple tree entries", Path: candidatePath}
	}
	header, rawPath, ok := bytes.Cut(records[0], []byte{'\t'})
	fields := strings.Fields(string(header))
	if !ok || len(fields) != 3 || string(rawPath) != candidatePath {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_path_state_invalid", Message: "Git returned a malformed or non-exact path state", Path: candidatePath}
	}
	mode := GitMode(fields[0])
	if mode != GitModeRegular && mode != GitModeExecutable {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_git_mode_unsupported", Message: "only regular and executable blobs are supported", Path: candidatePath}
	}
	if fields[1] != "blob" || !validObjectID(fields[2]) {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_object_unsupported", Message: "path does not resolve to a regular blob object", Path: candidatePath}
	}
	sizeOutput, blocker := runner.run(128, nil, "cat-file", "-s", fields[2])
	if blocker != nil {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_object_missing", Message: "path blob is missing or unavailable", Path: candidatePath}
	}
	size, err := strconv.ParseInt(strings.TrimSpace(string(sizeOutput)), 10, 64)
	if err != nil || size < 0 {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_object_invalid", Message: "path blob has an invalid size", Path: candidatePath}
	}
	if size > MaxPlanSourceBytes {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_source_limit", Message: fmt.Sprintf("path blob exceeds the %d-byte evidence limit", MaxPlanSourceBytes), Path: candidatePath}
	}
	data, blocker := runner.run(int(size)+1, nil, "cat-file", "blob", fields[2])
	if blocker != nil || int64(len(data)) != size {
		return BranchPathState{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_object_missing", Message: "path blob could not be read completely", Path: candidatePath}
	}
	state := BranchPathState{State: BranchPathPresent, Mode: mode, ObjectID: fields[2], SHA256: sha256Text(data)}
	return state, data, nil
}

func (runner *branchGitRunner) buildTransitions(before, after string, paths []string) ([]BranchPathTransition, *BranchEvidenceBlocker) {
	if blocker := runner.requireCommit(before, "transition before revision"); blocker != nil {
		return nil, blocker
	}
	if blocker := runner.requireCommit(after, "transition after revision"); blocker != nil {
		return nil, blocker
	}
	transitions := make([]BranchPathTransition, 0, len(paths))
	for _, candidatePath := range paths {
		beforeState, _, blocker := runner.readPathState(before, candidatePath)
		if blocker != nil {
			return nil, blocker
		}
		afterState, _, blocker := runner.readPathState(after, candidatePath)
		if blocker != nil {
			return nil, blocker
		}
		transitions = append(transitions, BranchPathTransition{Path: candidatePath, Before: beforeState, After: afterState})
	}
	return transitions, nil
}

func (runner *branchGitRunner) proveIncorporated(target, mergeBase, refTip, incorporatedCommit string, paths []string, refTransitions []BranchPathTransition) (BranchEvidenceProof, []BranchAttributeSnapshot, *BranchEvidenceBlocker) {
	if blocker := runner.requireCommit(incorporatedCommit, "incorporated commit"); blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	reachable, blocker := runner.isAncestor(incorporatedCommit, target)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	if !reachable {
		return BranchEvidenceProof{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_incorporated_commit_unreachable", Message: "incorporated commit is not reachable from the evidence revision"}
	}
	parent, blocker := runner.singleParent(incorporatedCommit)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	incorporatedTransitions, blocker := runner.buildTransitions(parent, incorporatedCommit, paths)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	refDigest, blocker := CanonicalBranchTransitionDigest(refTransitions)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	incorporatedDigest, blocker := CanonicalBranchTransitionDigest(incorporatedTransitions)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	if refDigest != incorporatedDigest {
		return BranchEvidenceProof{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_transition_mismatch", Message: "candidate transition does not exactly match the incorporated commit transition"}
	}
	revisions := []string{mergeBase, refTip, parent, incorporatedCommit}
	snapshots := make([]BranchAttributeSnapshot, 0, len(revisions))
	for _, revision := range revisions {
		snapshot, blocker := runner.attributeSnapshot(revision, paths)
		if blocker != nil {
			return BranchEvidenceProof{}, nil, blocker
		}
		snapshots = append(snapshots, snapshot)
	}
	if blocker := validateAttributeConsistency(snapshots); blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	for _, scope := range []struct {
		before      string
		after       string
		transitions []BranchPathTransition
	}{
		{before: mergeBase, after: refTip, transitions: refTransitions},
		{before: parent, after: incorporatedCommit, transitions: incorporatedTransitions},
	} {
		if blocker := runner.validateTextTransitions(scope.before, scope.after, scope.transitions); blocker != nil {
			return BranchEvidenceProof{}, nil, blocker
		}
	}
	refPatch, blocker := runner.stablePatchID(mergeBase, refTip, paths)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	incorporatedPatch, blocker := runner.stablePatchID(parent, incorporatedCommit, paths)
	if blocker != nil {
		return BranchEvidenceProof{}, nil, blocker
	}
	if refPatch != incorporatedPatch {
		return BranchEvidenceProof{}, nil, &BranchEvidenceBlocker{Code: "branch_evidence_patch_mismatch", Message: "stable candidate patch does not match the incorporated commit patch"}
	}
	return BranchEvidenceProof{
		Kind:                         BranchProofIncorporatedHistoryV1,
		TransitionDigest:             refDigest,
		StablePatchID:                refPatch,
		IncorporatedParent:           parent,
		IncorporatedCommit:           incorporatedCommit,
		IncorporatedTransitionDigest: incorporatedDigest,
		IncorporatedStablePatchID:    incorporatedPatch,
	}, snapshots, nil
}

func (runner *branchGitRunner) validateTextTransitions(before, after string, transitions []BranchPathTransition) *BranchEvidenceBlocker {
	for _, transition := range transitions {
		for _, endpoint := range []struct {
			revision string
			state    BranchPathState
		}{
			{revision: before, state: transition.Before},
			{revision: after, state: transition.After},
		} {
			if endpoint.state.State != BranchPathPresent {
				continue
			}
			observed, data, blocker := runner.readPathState(endpoint.revision, transition.Path)
			if blocker != nil {
				return blocker
			}
			if observed != endpoint.state {
				return &BranchEvidenceBlocker{Code: "branch_evidence_path_state_changed", Message: "path state changed during proof recomputation", Path: transition.Path}
			}
			if !utf8.Valid(data) {
				return &BranchEvidenceBlocker{Code: "branch_evidence_utf8_unsupported", Message: "incorporated-history proof requires valid UTF-8 candidate blobs", Path: transition.Path}
			}
			if bytes.IndexByte(data, 0) >= 0 {
				return &BranchEvidenceBlocker{Code: "branch_evidence_binary_unsupported", Message: "incorporated-history proof does not support NUL-containing candidate blobs", Path: transition.Path}
			}
		}
	}
	return nil
}

func (runner *branchGitRunner) isAncestor(ancestor, descendant string) (bool, *BranchEvidenceBlocker) {
	command := runner.command("merge-base", "--is-ancestor", ancestor, descendant)
	stderr := newBranchBoundedBuffer(maxBranchGitDiagnosticBytes)
	command.Stderr = stderr
	err := command.Run()
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, &BranchEvidenceBlocker{Code: "branch_evidence_ancestry_unavailable", Message: "cannot prove incorporated commit ancestry"}
}

func (runner *branchGitRunner) singleParent(commit string) (string, *BranchEvidenceBlocker) {
	output, blocker := runner.run(512, nil, "rev-list", "--parents", "-n", "1", commit)
	if blocker != nil {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_incorporated_parent_unavailable", Message: "cannot resolve incorporated commit parent"}
	}
	fields := strings.Fields(strings.TrimSpace(string(output)))
	if len(fields) != 2 || fields[0] != commit || !validObjectID(fields[1]) {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_incorporated_parent_ambiguous", Message: "incorporated proof requires one non-root, single-parent commit"}
	}
	return fields[1], nil
}

func (runner *branchGitRunner) stablePatchID(before, after string, paths []string) (string, *BranchEvidenceBlocker) {
	args := []string{
		"diff-tree", "--no-commit-id", "-r", "-p", "--full-index", "--binary", "--no-renames",
		"--no-ext-diff", "--no-textconv", "--no-color", "--no-indent-heuristic",
		"--diff-algorithm=myers", "--unified=3", before, after, "--",
	}
	args = append(args, paths...)
	patch, blocker := runner.run(MaxBranchEvidenceDiffBytes, nil, args...)
	if blocker != nil {
		blocker.Code = "branch_evidence_diff_unavailable"
		return "", blocker
	}
	if len(patch) == 0 {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_patch_empty", Message: "candidate transition produced no patch"}
	}
	if !utf8.Valid(patch) || bytes.IndexByte(patch, 0) >= 0 || bytes.Contains(patch, []byte("GIT binary patch")) || bytes.Contains(patch, []byte("Binary files ")) {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_binary_unsupported", Message: "candidate transition produced unsupported binary or non-UTF-8 patch output"}
	}
	output, blocker := runner.run(1024, patch, "patch-id", "--stable")
	if blocker != nil {
		blocker.Code = "branch_evidence_patch_id_unavailable"
		return "", blocker
	}
	lines := nonEmptyLines(output)
	if len(lines) != 1 {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_patch_id_ambiguous", Message: fmt.Sprintf("expected one stable patch ID, found %d", len(lines))}
	}
	fields := strings.Fields(lines[0])
	if len(fields) < 1 || !validObjectID(fields[0]) {
		return "", &BranchEvidenceBlocker{Code: "branch_evidence_patch_id_invalid", Message: "Git returned an invalid stable patch ID"}
	}
	return fields[0], nil
}

func (runner *branchGitRunner) attributeSnapshot(revision string, paths []string) (BranchAttributeSnapshot, *BranchEvidenceBlocker) {
	input := make([]byte, 0)
	for _, candidatePath := range paths {
		input = append(input, []byte(candidatePath)...)
		input = append(input, 0)
	}
	args := []string{"check-attr", "--source=" + revision, "-z", "--stdin", "diff", "text", "working-tree-encoding", "binary"}
	output, blocker := runner.run(maxBranchGitDiagnosticBytes, input, args...)
	if blocker != nil {
		return BranchAttributeSnapshot{}, &BranchEvidenceBlocker{Code: "branch_evidence_attributes_unavailable", Message: "cannot resolve effective attributes at " + revision}
	}
	fields := splitNUL(output)
	if len(fields) != len(paths)*4*3 {
		return BranchAttributeSnapshot{}, &BranchEvidenceBlocker{Code: "branch_evidence_attributes_invalid", Message: "Git returned malformed effective attributes"}
	}
	byPath := make(map[string]*BranchPathAttributes, len(paths))
	for _, candidatePath := range paths {
		byPath[candidatePath] = &BranchPathAttributes{Path: candidatePath}
	}
	for index := 0; index < len(fields); index += 3 {
		candidatePath := string(fields[index])
		attribute := string(fields[index+1])
		value := string(fields[index+2])
		entry := byPath[candidatePath]
		if entry == nil {
			return BranchAttributeSnapshot{}, &BranchEvidenceBlocker{Code: "branch_evidence_attributes_invalid", Message: "Git returned attributes for an unexpected path", Path: candidatePath}
		}
		switch attribute {
		case "diff":
			entry.Diff = value
		case "text":
			entry.Text = value
		case "working-tree-encoding":
			entry.WorkingTreeEncoding = value
		case "binary":
			entry.Binary = value
		default:
			return BranchAttributeSnapshot{}, &BranchEvidenceBlocker{Code: "branch_evidence_attributes_invalid", Message: "Git returned an unexpected attribute", Path: candidatePath}
		}
	}
	entries := make([]BranchPathAttributes, 0, len(paths))
	for _, candidatePath := range paths {
		entry := *byPath[candidatePath]
		if entry.Diff != "unspecified" && entry.Diff != "set" && entry.Diff != "unset" {
			return BranchAttributeSnapshot{}, &BranchEvidenceBlocker{Code: "branch_evidence_attributes_unsupported", Message: "custom diff drivers are not supported", Path: candidatePath}
		}
		if entry.WorkingTreeEncoding != "unspecified" && entry.WorkingTreeEncoding != "unset" {
			return BranchAttributeSnapshot{}, &BranchEvidenceBlocker{Code: "branch_evidence_attributes_unsupported", Message: "working-tree-encoding is not supported", Path: candidatePath}
		}
		if entry.Binary != "unspecified" && entry.Binary != "unset" {
			return BranchAttributeSnapshot{}, &BranchEvidenceBlocker{Code: "branch_evidence_attributes_unsupported", Message: "binary attributes are not supported", Path: candidatePath}
		}
		entries = append(entries, entry)
	}
	attributeFiles := []BranchAttributeFileState{}
	for _, attributePath := range branchAttributePaths(paths) {
		state, _, blocker := runner.readPathState(revision, attributePath)
		if blocker != nil {
			return BranchAttributeSnapshot{}, blocker
		}
		attributeFiles = append(attributeFiles, BranchAttributeFileState{Path: attributePath, State: state})
	}
	snapshot := BranchAttributeSnapshot{Revision: revision, PathAttributes: entries, AttributeFiles: attributeFiles}
	snapshot.Digest = canonicalJSONDigest(struct {
		Revision       string                     `json:"revision"`
		PathAttributes []BranchPathAttributes     `json:"path_attributes"`
		AttributeFiles []BranchAttributeFileState `json:"attribute_files"`
	}{snapshot.Revision, snapshot.PathAttributes, snapshot.AttributeFiles})
	return snapshot, nil
}

func branchAttributePaths(paths []string) []string {
	set := map[string]bool{".gitattributes": true}
	for _, candidatePath := range paths {
		directory := path.Dir(candidatePath)
		if directory == "." {
			continue
		}
		segments := strings.Split(directory, "/")
		for index := range segments {
			set[strings.Join(segments[:index+1], "/")+"/.gitattributes"] = true
		}
	}
	result := make([]string, 0, len(set))
	for attributePath := range set {
		result = append(result, attributePath)
	}
	sort.Strings(result)
	return result
}

func validateAttributeConsistency(snapshots []BranchAttributeSnapshot) *BranchEvidenceBlocker {
	if len(snapshots) == 0 {
		return &BranchEvidenceBlocker{Code: "branch_evidence_attributes_missing", Message: "attribute snapshots are missing"}
	}
	baseline := map[string]BranchPathAttributes{}
	for _, entry := range snapshots[0].PathAttributes {
		baseline[entry.Path] = entry
	}
	for _, snapshot := range snapshots[1:] {
		if len(snapshot.PathAttributes) != len(baseline) {
			return &BranchEvidenceBlocker{Code: "branch_evidence_attributes_conflict", Message: "effective attribute path set changed across proof endpoints"}
		}
		for _, entry := range snapshot.PathAttributes {
			prior, ok := baseline[entry.Path]
			if !ok || prior.Diff != entry.Diff || prior.Text != entry.Text || prior.WorkingTreeEncoding != entry.WorkingTreeEncoding || prior.Binary != entry.Binary {
				return &BranchEvidenceBlocker{Code: "branch_evidence_attributes_conflict", Message: "effective text or binary attributes changed across proof endpoints", Path: entry.Path}
			}
		}
	}
	return nil
}

type branchBoundedBuffer struct {
	buffer   bytes.Buffer
	limit    int
	exceeded bool
}

func newBranchBoundedBuffer(limit int) *branchBoundedBuffer {
	return &branchBoundedBuffer{limit: limit}
}

func (buffer *branchBoundedBuffer) Write(data []byte) (int, error) {
	original := len(data)
	remaining := buffer.limit - buffer.buffer.Len()
	if remaining > 0 {
		if len(data) > remaining {
			_, _ = buffer.buffer.Write(data[:remaining])
		} else {
			_, _ = buffer.buffer.Write(data)
		}
	}
	if original > remaining {
		buffer.exceeded = true
	}
	return original, nil
}

func (buffer *branchBoundedBuffer) Bytes() []byte  { return buffer.buffer.Bytes() }
func (buffer *branchBoundedBuffer) String() string { return buffer.buffer.String() }

func splitNUL(data []byte) [][]byte {
	parts := bytes.Split(data, []byte{0})
	result := make([][]byte, 0, len(parts))
	for _, part := range parts {
		if len(part) > 0 {
			result = append(result, part)
		}
	}
	return result
}

func nonEmptyLines(data []byte) []string {
	lines := []string{}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func validObjectID(value string) bool {
	if len(value) != 40 && len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil && value == strings.ToLower(value)
}

func validSHA256(value string) bool {
	return len(value) == sha256.Size*2 && validObjectID(value)
}

func canonicalJSONDigest(value any) string {
	data, _ := json.Marshal(value)
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
