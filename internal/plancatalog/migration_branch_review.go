package plancatalog

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type v3BranchReviewVerification struct {
	Results  []BranchReviewResult
	Evidence []BranchCandidateEvidence
	Problems []Problem
	Deferred bool
	Resolved bool
	FollowUp *IncrementalFollowUp
}

func verifyV3BranchReviews(root string, ledger MigrationLedger, item MigrationRecord, inventoryRecord InventoryRecord, candidate Candidate, settings RepositorySettings) v3BranchReviewVerification {
	verification := v3BranchReviewVerification{}
	reviews := map[string]BranchTouchReview{}
	for _, review := range item.BranchTouchReviews {
		key := branchReviewIdentityKey(review.Identity)
		if _, duplicate := reviews[key]; duplicate {
			verification.Problems = append(verification.Problems, Problem{Code: "branch_review_duplicate", Message: "ledger contains a duplicate branch-touch review identity", Path: item.CurrentPath, Severity: SeverityError})
			continue
		}
		reviews[key] = review
	}
	seen := map[string]bool{}
	deferredOwners := 0
	for _, rawRef := range inventoryRecord.ActiveBranchTouch {
		if remoteEvidence, mapped := matchingRemoteEvidenceForTracking(root, rawRef, item.CurrentPath, ledger); mapped {
			paths := make([]string, 0, len(remoteEvidence.PathTransitions))
			for _, transition := range remoteEvidence.PathTransitions {
				paths = append(paths, transition.Path)
			}
			request := BranchEvidenceRequest{
				Root: root, EvidenceScope: "remote-aware", RepositoryID: ledger.RepositoryID,
				LogicalRef: logicalRefForLocalTouch(rawRef), Ref: rawRef, EvidenceRevision: ledger.EvidenceRevision,
				PolicyDigest: ledger.PolicyDigest, CandidateSetDigest: ledger.CandidateSetDigest,
				CandidatePath: item.CurrentPath, Paths: paths, ExpectedRefTip: remoteEvidence.RefTip,
				ExpectedKind: candidate.ExpectedKind,
			}
			if remoteEvidence.Proof != nil && remoteEvidence.Proof.Kind == BranchProofIncorporatedHistoryV1 {
				request.IncorporatedCommit = remoteEvidence.Proof.IncorporatedCommit
			}
			localEvidence := proposeBranchCandidateEvidence(request)
			if EquivalentBranchCandidateEvidence(localEvidence, remoteEvidence) {
				// The source-stable remote review below authorizes this exact
				// local tracking observation; do not require a second review.
				continue
			}
			verification.Problems = append(verification.Problems, Problem{Code: "branch_tracking_overlay_mismatch", Message: "tracking ref and fresh mirror evidence differ", Path: item.CurrentPath, Severity: SeverityError, Remediation: "refresh tracking refs and reinventory remote-aware evidence"})
		}
		logicalRef := logicalRefForLocalTouch(rawRef)
		identity := BranchReviewIdentity{EvidenceScope: ledger.BranchEvidence.EvidenceScope, RepositoryID: ledger.RepositoryID, LogicalRef: logicalRef, CandidatePath: item.CurrentPath}
		key := branchReviewIdentityKey(identity)
		review, ok := reviews[key]
		if !ok {
			verification.Problems = append(verification.Problems, Problem{Code: "branch_touch_unreviewed", Message: "current touching ref has no exact reviewed disposition: " + logicalRef, Path: item.CurrentPath, Severity: SeverityError, Remediation: "regenerate inventory and review every current candidate/ref identity"})
			continue
		}
		seen[key] = true
		result := BranchReviewResult{Identity: review.Identity, Disposition: review.Disposition, Blockers: []BranchEvidenceBlocker{}}
		paths := make([]string, 0, len(review.PathTransitions))
		for _, transition := range review.PathTransitions {
			paths = append(paths, transition.Path)
		}
		request := BranchEvidenceRequest{
			Root:               root,
			EvidenceScope:      review.Identity.EvidenceScope,
			RepositoryID:       review.Identity.RepositoryID,
			LogicalRef:         review.Identity.LogicalRef,
			Ref:                rawRef,
			EvidenceRevision:   ledger.EvidenceRevision,
			PolicyDigest:       ledger.PolicyDigest,
			CandidateSetDigest: ledger.CandidateSetDigest,
			CandidatePath:      item.CurrentPath,
			Paths:              paths,
			ExpectedRefTip:     review.RefTip,
		}
		if review.Proof != nil && review.Proof.Kind == BranchProofIncorporatedHistoryV1 {
			request.IncorporatedCommit = review.Proof.IncorporatedCommit
		}
		evidence := proposeBranchCandidateEvidence(request)
		verification.Evidence = append(verification.Evidence, evidence)
		if review.Disposition != BranchDispositionDeferredActiveOwner {
			result.Blockers = append(result.Blockers, evidence.Blockers...)
		}
		if evidence.RefTip != review.RefTip {
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_ref_moved", Message: "reviewed ref tip changed", Path: item.CurrentPath, Ref: rawRef})
		}
		if evidence.MergeBase != review.MergeBase {
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_merge_base_moved", Message: "reviewed merge base changed", Path: item.CurrentPath, Ref: rawRef})
		}
		if !reflect.DeepEqual(evidence.PathTransitions, review.PathTransitions) {
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_path_state_moved", Message: "reviewed candidate path transition changed", Path: item.CurrentPath, Ref: rawRef})
		}
		switch review.Disposition {
		case BranchDispositionSameContentNonOwner, BranchDispositionIncorporatedHistoryNonOwner:
			if evidence.ProposedDisposition != review.Disposition || review.Proof == nil || evidence.Proof == nil || !reflect.DeepEqual(*evidence.Proof, *review.Proof) {
				result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_clearance_proof_mismatch", Message: "recomputed proof does not match the reviewed non-owner disposition", Path: item.CurrentPath, Ref: rawRef})
			}
		case BranchDispositionDeferredActiveOwner:
			deferredOwners++
			deferredProblems := verifyDeferredOwner(root, ledger, item, review, candidate, settings)
			for _, problem := range deferredProblems {
				result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: problem.Code, Message: problem.Message, Path: problem.Path, Ref: rawRef})
			}
			if review.IncrementalFollowUp != nil {
				followUp := *review.IncrementalFollowUp
				verification.FollowUp = &followUp
			}
		case BranchDispositionActiveOwner:
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "active_branch_ownership", Message: "reviewed active owner is not approved for mixed deferral", Path: item.CurrentPath, Ref: rawRef})
		case BranchDispositionBlocking:
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_review_blocking", Message: strings.Join(review.BlockingReasons, "; "), Path: item.CurrentPath, Ref: rawRef})
		default:
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_disposition_unknown", Message: "unknown reviewed branch disposition", Path: item.CurrentPath, Ref: rawRef})
		}
		result.Blockers = uniqueBranchBlockers(result.Blockers)
		result.Verified = len(result.Blockers) == 0
		verification.Results = append(verification.Results, result)
		if !result.Verified {
			for _, blocker := range result.Blockers {
				verification.Problems = append(verification.Problems, Problem{Code: blocker.Code, Message: blocker.Message, Path: item.CurrentPath, Severity: SeverityError, Remediation: "reinventory and review the current branch evidence"})
			}
		}
	}
	for _, evidence := range ledger.ObservedRemoteBranchEvidence {
		if evidence.Identity.RepositoryID != ledger.RepositoryID || evidence.Identity.CandidatePath != item.CurrentPath {
			continue
		}
		key := branchReviewIdentityKey(evidence.Identity)
		review, ok := reviews[key]
		if !ok {
			verification.Problems = append(verification.Problems, Problem{Code: "branch_touch_unreviewed", Message: "current remote touching ref has no exact reviewed disposition: " + evidence.Identity.LogicalRef, Path: item.CurrentPath, Severity: SeverityError, Remediation: "regenerate remote-aware inventory and review every current candidate/ref identity"})
			continue
		}
		seen[key] = true
		result := BranchReviewResult{Identity: review.Identity, Disposition: review.Disposition, Blockers: []BranchEvidenceBlocker{}}
		verification.Evidence = append(verification.Evidence, evidence)
		if !reflect.DeepEqual(evidence.Identity, review.Identity) || evidence.RefTip != review.RefTip {
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_ref_moved", Message: "reviewed remote ref identity or tip changed", Path: item.CurrentPath, Ref: evidence.Identity.LogicalRef})
		}
		if evidence.MergeBase != review.MergeBase {
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_merge_base_moved", Message: "reviewed remote merge base changed", Path: item.CurrentPath, Ref: evidence.Identity.LogicalRef})
		}
		if !reflect.DeepEqual(evidence.PathTransitions, review.PathTransitions) {
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_path_state_moved", Message: "reviewed remote candidate path transition changed", Path: item.CurrentPath, Ref: evidence.Identity.LogicalRef})
		}
		switch review.Disposition {
		case BranchDispositionSameContentNonOwner, BranchDispositionIncorporatedHistoryNonOwner:
			result.Blockers = append(result.Blockers, evidence.Blockers...)
			if evidence.ProposedDisposition != review.Disposition || review.Proof == nil || evidence.Proof == nil || !reflect.DeepEqual(*evidence.Proof, *review.Proof) {
				result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_clearance_proof_mismatch", Message: "recomputed remote proof does not match the reviewed non-owner disposition", Path: item.CurrentPath, Ref: evidence.Identity.LogicalRef})
			}
		case BranchDispositionDeferredActiveOwner:
			deferredOwners++
			for _, problem := range verifyDeferredRemoteOwner(root, ledger, item, review, evidence, candidate, settings) {
				result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: problem.Code, Message: problem.Message, Path: problem.Path, Ref: evidence.Identity.LogicalRef})
			}
			if review.IncrementalFollowUp != nil {
				followUp := *review.IncrementalFollowUp
				verification.FollowUp = &followUp
			}
		case BranchDispositionActiveOwner:
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "active_branch_ownership", Message: "reviewed remote active owner is not approved for mixed deferral", Path: item.CurrentPath, Ref: evidence.Identity.LogicalRef})
		case BranchDispositionBlocking:
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_review_blocking", Message: strings.Join(review.BlockingReasons, "; "), Path: item.CurrentPath, Ref: evidence.Identity.LogicalRef})
		default:
			result.Blockers = append(result.Blockers, BranchEvidenceBlocker{Code: "branch_disposition_unknown", Message: "unknown reviewed remote branch disposition", Path: item.CurrentPath, Ref: evidence.Identity.LogicalRef})
		}
		result.Blockers = uniqueBranchBlockers(result.Blockers)
		result.Verified = len(result.Blockers) == 0
		verification.Results = append(verification.Results, result)
		if !result.Verified {
			for _, blocker := range result.Blockers {
				verification.Problems = append(verification.Problems, branchReviewProblem(blocker, item.CurrentPath))
			}
		}
	}
	for key, review := range reviews {
		if seen[key] {
			continue
		}
		if review.Disposition == BranchDispositionDeferredActiveOwner {
			resolutionProblems := verifyMergedDeferredOwner(root, ledger, item, review, candidate, settings)
			if !HasErrors(resolutionProblems) {
				verification.Evidence = append(verification.Evidence, historicalDeferredOwnerEvidence(root, ledger, review, candidate))
				verification.Results = append(verification.Results, BranchReviewResult{Identity: review.Identity, Disposition: review.Disposition, Verified: true, Blockers: []BranchEvidenceBlocker{}})
				verification.Resolved = true
				continue
			}
			verification.Problems = append(verification.Problems, resolutionProblems...)
		}
		verification.Results = append(verification.Results, BranchReviewResult{Identity: review.Identity, Disposition: review.Disposition, Verified: false, Blockers: []BranchEvidenceBlocker{{Code: "branch_review_stale", Message: "reviewed ref/candidate identity is not a current touch", Path: item.CurrentPath, Ref: review.Identity.LogicalRef}}})
		verification.Problems = append(verification.Problems, Problem{Code: "branch_review_stale", Message: "reviewed ref/candidate identity is not a current touch", Path: item.CurrentPath, Severity: SeverityError, Remediation: "regenerate inventory rather than reusing stale clearance"})
	}
	if deferredOwners > 1 {
		verification.Problems = append(verification.Problems, Problem{Code: "multiple_active_owners", Message: "a candidate may retain at most one deferred active owner", Path: item.CurrentPath, Severity: SeverityError})
	}
	verification.Deferred = deferredOwners == 1 && !HasErrors(verification.Problems)
	sort.SliceStable(verification.Results, func(i, j int) bool {
		return branchReviewIdentityKey(verification.Results[i].Identity) < branchReviewIdentityKey(verification.Results[j].Identity)
	})
	SortProblems(verification.Problems)
	return verification
}

func historicalDeferredOwnerEvidence(root string, ledger MigrationLedger, review BranchTouchReview, candidate Candidate) BranchCandidateEvidence {
	paths := make([]string, 0, len(review.PathTransitions))
	for _, transition := range review.PathTransitions {
		paths = append(paths, transition.Path)
	}
	return proposeBranchCandidateEvidence(BranchEvidenceRequest{
		Root: root, EvidenceScope: review.Identity.EvidenceScope, RepositoryID: review.Identity.RepositoryID,
		LogicalRef: review.Identity.LogicalRef, Ref: review.RefTip, EvidenceRevision: ledger.EvidenceRevision,
		PolicyDigest: ledger.PolicyDigest, CandidateSetDigest: ledger.CandidateSetDigest,
		CandidatePath: review.Identity.CandidatePath, Paths: paths, ExpectedRefTip: review.RefTip,
		ExpectedKind: candidate.ExpectedKind,
	})
}

func verifyMergedDeferredOwner(root string, ledger MigrationLedger, item MigrationRecord, review BranchTouchReview, candidate Candidate, settings RepositorySettings) []Problem {
	problems := verifyDeferredOwner(root, ledger, item, review, candidate, settings)
	owner := review.OwnerTipCandidate
	if owner == nil || review.IncrementalFollowUp == nil {
		return problems
	}
	head, headErr := gitText(root, "rev-parse", "--verify", "HEAD^{commit}")
	if headErr != nil {
		return append(problems, Problem{Code: "deferred_owner_merge_unavailable", Message: "current revision is unavailable", Path: item.CurrentPath, Severity: SeverityError})
	}
	runner, blocker := newBranchGitRunner(root)
	if blocker != nil {
		return append(problems, Problem{Code: blocker.Code, Message: blocker.Message, Path: item.CurrentPath, Severity: SeverityError})
	}
	defer runner.close()
	merged, ancestryBlocker := runner.isAncestor(review.RefTip, head)
	if ancestryBlocker != nil || !merged {
		return append(problems, Problem{Code: "deferred_owner_not_integrated", Message: "reviewed owner tip is not reachable from the current revision", Path: item.CurrentPath, Severity: SeverityError})
	}
	stateAtHead, stateBlocker := ReadBranchPathState(root, head, item.CurrentPath)
	data, readErr := gitBytes(root, "show", head+":"+item.CurrentPath)
	if stateBlocker != nil || readErr != nil || owner.Path != item.CurrentPath || stateAtHead.State != BranchPathPresent || stateAtHead.Mode != owner.Mode || stateAtHead.ObjectID != owner.ObjectID || stateAtHead.SHA256 != owner.SHA256 || sha256Text(data) != owner.SHA256 {
		return append(problems, Problem{Code: "deferred_owner_merge_mismatch", Message: "current plan bytes do not exactly match the integrated reviewed owner", Path: item.CurrentPath, Severity: SeverityError})
	}
	record, parseErr := ParseDocument(item.CurrentPath, data, candidate.ExpectedKind)
	if parseErr != nil || record.Metadata == nil || !metadataEqual(*record.Metadata, metadataFromDecision(item)) {
		return append(problems, Problem{Code: "deferred_owner_metadata_mismatch", Message: "integrated owner did not supply the exact reviewed canonical metadata", Path: item.CurrentPath, Severity: SeverityError})
	}
	return problems
}

func matchingRemoteEvidenceForTracking(root, rawRef, candidatePath string, ledger MigrationLedger) (BranchCandidateEvidence, bool) {
	if ledger.BranchEvidence == nil || ledger.BranchEvidence.EvidenceScope != "remote-aware" {
		return BranchCandidateEvidence{}, false
	}
	remoteName, sourceMatched := uniqueTrackingRemoteForSource(root, ledger.BranchEvidence.RemoteSourceIdentitySHA256)
	prefix := "refs/remotes/" + remoteName + "/"
	if !sourceMatched || !strings.HasPrefix(rawRef, prefix) {
		return BranchCandidateEvidence{}, false
	}
	branch := strings.TrimPrefix(rawRef, prefix)
	logical := "remote:" + ledger.RepositoryID + ":refs/heads/" + branch
	var matchedEvidence BranchCandidateEvidence
	found := false
	for _, evidence := range ledger.ObservedRemoteBranchEvidence {
		if evidence.Identity.EvidenceScope != "remote-aware" || evidence.Identity.RepositoryID != ledger.RepositoryID || evidence.Identity.CandidatePath != candidatePath || evidence.Identity.LogicalRef != logical || !validSHA256(evidence.Digest) || CanonicalBranchEvidenceDigest(evidence) != evidence.Digest {
			continue
		}
		if found {
			return BranchCandidateEvidence{}, false
		}
		matchedEvidence = evidence
		found = true
	}
	return matchedEvidence, found
}

func verifyDeferredOwner(root string, ledger MigrationLedger, item MigrationRecord, review BranchTouchReview, candidate Candidate, settings RepositorySettings) []Problem {
	problems := []Problem{}
	owner := review.OwnerTipCandidate
	followUp := review.IncrementalFollowUp
	if ledger.TargetCatalogMode != ModeMixed {
		problems = append(problems, Problem{Code: "canonical_deferral_forbidden", Message: "direct canonical mode cannot defer an active owner", Path: item.CurrentPath, Severity: SeverityError})
	}
	if settings.Mode != ModeMixed || item.CurrentPath != item.TargetPath || !item.Deferred {
		problems = append(problems, Problem{Code: "mixed_cutover_proof_missing", Message: "active-owner deferral requires active mixed mode, a same-path record, and deferred=true", Path: item.CurrentPath, Severity: SeverityError})
	}
	if owner == nil || followUp == nil {
		return append(problems, Problem{Code: "deferred_follow_up_missing", Message: "deferred active owner requires an owner snapshot and incremental follow-up", Path: item.CurrentPath, Severity: SeverityError})
	}
	stateAtTip, blocker := ReadBranchPathState(root, review.RefTip, owner.Path)
	if blocker != nil || !reflect.DeepEqual(stateAtTip, BranchPathState{State: "present", Mode: owner.Mode, ObjectID: owner.ObjectID, SHA256: owner.SHA256}) {
		problems = append(problems, Problem{Code: "deferred_owner_snapshot_mismatch", Message: "owner-tip candidate no longer matches the reviewed path/blob", Path: item.CurrentPath, Severity: SeverityError})
	} else if data, err := gitBytes(root, "show", review.RefTip+":"+owner.Path); err != nil {
		problems = append(problems, Problem{Code: "deferred_owner_candidate_unavailable", Message: err.Error(), Path: owner.Path, Severity: SeverityError})
	} else if record, parseErr := ParseDocument(owner.Path, data, candidate.ExpectedKind); parseErr != nil || record.Header.Lifecycle != LifecycleActive {
		problems = append(problems, Problem{Code: "deferred_owner_lifecycle_invalid", Message: "owner-tip path must remain a valid active plan candidate", Path: owner.Path, Severity: SeverityError})
	}
	baseline, baselineProblems := loadMixedBaseline(root, settings.CutoverRevision)
	problems = append(problems, baselineProblems...)
	if baseline[item.CurrentPath] == "" || baseline[item.CurrentPath] != item.SourceSHA256 {
		problems = append(problems, Problem{Code: "mixed_grandfathered_plan_modified", Message: "deferred plan does not match its register-proven cutover blob", Path: item.CurrentPath, Severity: SeverityError})
	}
	expectedSuccess := []string{"new-reviewed-ledger-applied-to-merged-target-bytes", "owner-supplied-canonical-metadata-at-merge"}
	actualSuccess := append([]string{}, followUp.SuccessConditions...)
	sort.Strings(actualSuccess)
	if followUp.Action != "reinventory-and-incremental-migrate" || followUp.Trigger != "owner-ref-integrated-or-target-path-changed" || followUp.CutoverRevision != settings.CutoverRevision || followUp.PlanPath != item.CurrentPath || followUp.CutoverSourceSHA256 != item.SourceSHA256 || followUp.OwnerRef != review.Identity.LogicalRef || followUp.OwnerTip != review.RefTip || followUp.OwnerPath != owner.Path || followUp.OwnerMode != owner.Mode || followUp.OwnerObjectID != owner.ObjectID || followUp.OwnerSHA256 != owner.SHA256 || followUp.State != "pending" || !equalStrings(actualSuccess, expectedSuccess) {
		problems = append(problems, Problem{Code: "deferred_follow_up_stale", Message: "incremental follow-up does not bind the exact cutover and owner evidence", Path: item.CurrentPath, Severity: SeverityError})
	}
	return problems
}

func verifyDeferredRemoteOwner(root string, ledger MigrationLedger, item MigrationRecord, review BranchTouchReview, evidence BranchCandidateEvidence, candidate Candidate, settings RepositorySettings) []Problem {
	problems := []Problem{}
	owner := review.OwnerTipCandidate
	followUp := review.IncrementalFollowUp
	if ledger.TargetCatalogMode != ModeMixed {
		problems = append(problems, Problem{Code: "canonical_deferral_forbidden", Message: "direct canonical mode cannot defer an active remote owner", Path: item.CurrentPath, Severity: SeverityError})
	}
	if settings.Mode != ModeMixed || item.CurrentPath != item.TargetPath || !item.Deferred {
		problems = append(problems, Problem{Code: "mixed_cutover_proof_missing", Message: "remote active-owner deferral requires active mixed mode, a same-path record, and deferred=true", Path: item.CurrentPath, Severity: SeverityError})
	}
	if owner == nil || followUp == nil {
		return append(problems, Problem{Code: "deferred_follow_up_missing", Message: "deferred remote active owner requires an owner snapshot and incremental follow-up", Path: item.CurrentPath, Severity: SeverityError})
	}
	if evidence.ProposedDisposition != BranchDispositionActiveOwner || evidence.OwnerTipCandidate == nil || !reflect.DeepEqual(*evidence.OwnerTipCandidate, *owner) || owner.Lifecycle != LifecycleActive {
		problems = append(problems, Problem{Code: "deferred_owner_snapshot_mismatch", Message: "fresh remote owner evidence does not match the reviewed active candidate snapshot", Path: item.CurrentPath, Severity: SeverityError})
	}
	baseline, baselineProblems := loadMixedBaseline(root, settings.CutoverRevision)
	problems = append(problems, baselineProblems...)
	if baseline[item.CurrentPath] == "" || baseline[item.CurrentPath] != item.SourceSHA256 {
		problems = append(problems, Problem{Code: "mixed_grandfathered_plan_modified", Message: "deferred plan does not match its register-proven cutover blob", Path: item.CurrentPath, Severity: SeverityError})
	}
	expectedSuccess := []string{"new-reviewed-ledger-applied-to-merged-target-bytes", "owner-supplied-canonical-metadata-at-merge"}
	actualSuccess := append([]string{}, followUp.SuccessConditions...)
	sort.Strings(actualSuccess)
	if followUp.Action != "reinventory-and-incremental-migrate" || followUp.Trigger != "owner-ref-integrated-or-target-path-changed" || followUp.CutoverRevision != settings.CutoverRevision || followUp.PlanPath != item.CurrentPath || followUp.CutoverSourceSHA256 != item.SourceSHA256 || followUp.OwnerRef != review.Identity.LogicalRef || followUp.OwnerTip != review.RefTip || followUp.OwnerPath != owner.Path || followUp.OwnerMode != owner.Mode || followUp.OwnerObjectID != owner.ObjectID || followUp.OwnerSHA256 != owner.SHA256 || followUp.State != "pending" || !equalStrings(actualSuccess, expectedSuccess) {
		problems = append(problems, Problem{Code: "deferred_follow_up_stale", Message: "incremental follow-up does not bind the exact cutover and remote owner evidence", Path: item.CurrentPath, Severity: SeverityError})
	}
	_ = candidate
	return problems
}

func logicalRefForLocalTouch(ref string) string {
	if strings.HasPrefix(ref, "refs/heads/") {
		return "local:" + ref
	}
	if strings.HasPrefix(ref, "refs/remotes/") {
		parts := strings.SplitN(strings.TrimPrefix(ref, "refs/remotes/"), "/", 2)
		if len(parts) == 2 {
			return "tracking:" + parts[0] + ":refs/heads/" + parts[1]
		}
	}
	return "tracking:" + ref
}

func rawRefForLogicalRef(logical string) (string, bool) {
	if strings.HasPrefix(logical, "local:refs/heads/") {
		return strings.TrimPrefix(logical, "local:"), true
	}
	if strings.HasPrefix(logical, "tracking:") {
		remainder := strings.TrimPrefix(logical, "tracking:")
		parts := strings.SplitN(remainder, ":refs/heads/", 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			return "refs/remotes/" + parts[0] + "/" + parts[1], true
		}
		if strings.HasPrefix(remainder, "refs/") {
			return remainder, true
		}
	}
	return "", false
}

func branchReviewIdentityKey(identity BranchReviewIdentity) string {
	return identity.EvidenceScope + "\x00" + identity.RepositoryID + "\x00" + identity.LogicalRef + "\x00" + identity.CandidatePath
}

func branchEvidenceSetDigest(evidence []BranchCandidateEvidence, remoteOverlayDigest string) string {
	items := append([]BranchCandidateEvidence{}, evidence...)
	for index := range items {
		items[index].GitVersion = ""
		items[index].Digest = CanonicalBranchEvidenceDigest(items[index])
	}
	sort.SliceStable(items, func(i, j int) bool {
		return branchReviewIdentityKey(items[i].Identity) < branchReviewIdentityKey(items[j].Identity)
	})
	value := struct {
		Algorithm           string                    `json:"algorithm"`
		RemoteOverlayDigest string                    `json:"remote_overlay_digest,omitempty"`
		Evidence            []BranchCandidateEvidence `json:"evidence"`
	}{Algorithm: BranchEvidenceAlgorithm, RemoteOverlayDigest: remoteOverlayDigest, Evidence: items}
	data, _ := json.Marshal(value)
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func uniqueBranchBlockers(blockers []BranchEvidenceBlocker) []BranchEvidenceBlocker {
	seen := map[string]bool{}
	result := []BranchEvidenceBlocker{}
	for _, blocker := range blockers {
		key := blocker.Code + "\x00" + blocker.Message + "\x00" + blocker.Path + "\x00" + blocker.Ref
		if !seen[key] {
			seen[key] = true
			result = append(result, blocker)
		}
	}
	return result
}

func onlyExpectedActiveProofBlockers(blockers []BranchEvidenceBlocker) bool {
	for _, blocker := range blockers {
		switch blocker.Code {
		case "branch_candidate_active", "branch_clearance_unproven", "incorporated_commit_missing":
			continue
		default:
			return false
		}
	}
	return true
}

func branchReviewProblem(blocker BranchEvidenceBlocker, candidatePath string) Problem {
	path := blocker.Path
	if path == "" {
		path = candidatePath
	}
	return Problem{Code: blocker.Code, Message: blocker.Message, Path: path, Severity: SeverityError, Remediation: fmt.Sprintf("reinventory %s and review the current evidence", candidatePath)}
}

func currentV3ExternalAuthorityDigest(root string, ledger MigrationLedger) (string, []Problem, error) {
	return currentV3ExternalAuthorityDigestForConfig(root, ledger, ledger.TargetConfigPreconditionSHA256)
}

func currentV3ExternalAuthorityDigestForConfig(root string, ledger MigrationLedger, requiredConfigSHA string) (string, []Problem, error) {
	settings, settingProblems := LoadRepositorySettings(root)
	problems := append([]Problem{}, settingProblems...)
	settings.DiscoveryVersion = DiscoveryV2
	settings.Mode = ledger.TargetCatalogMode
	settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	if settings.RepositoryID != ledger.RepositoryID || settings.PolicyDigest != ledger.PolicyDigest {
		problems = append(problems, Problem{Code: "migration_policy_mismatch", Message: "repository identity or candidate policy moved after review", Severity: SeverityError})
	}
	classification, err := ClassifyCommitTree(root, ledger.EvidenceRevision, settings, ledger.TargetCatalogMode, ledger.RepositoryID)
	if err != nil {
		return "", problems, err
	}
	if !classification.Complete || classification.CandidateSetDigest != ledger.CandidateSetDigest {
		problems = append(problems, Problem{Code: "migration_candidate_set_mismatch", Message: "evidence-revision candidate set no longer matches the reviewed digest", Severity: SeverityError})
	}
	candidates := map[string]Candidate{}
	for _, candidate := range classification.Candidates {
		if candidate.Ownership == OwnershipOwned {
			candidates[candidate.Path] = candidate
		}
	}
	touches, touchProblems := activeBranchTouchesWithSettings(root, settings)
	problems = append(problems, touchProblems...)
	checkpoint, checkpointProblems := verifyLedgerCheckpoint(root, ledger.EvidenceRevision, ledger.ArtifactPath, ledger.ArtifactSHA256)
	problems = append(problems, checkpointProblems...)
	configSHA, configErr := configPrecondition(root)
	if configErr != nil {
		return "", problems, configErr
	}
	if configSHA != requiredConfigSHA {
		problems = append(problems, Problem{Code: "target_config_precondition_mismatch", Message: "config bytes moved after review", Path: ".codeheart/kit.config.yaml", Severity: SeverityError})
	}
	evidence := []BranchCandidateEvidence{}
	for _, item := range ledger.Records {
		candidate, ok := candidates[item.CurrentPath]
		if !ok {
			problems = append(problems, Problem{Code: "migration_candidate_missing", Message: "reviewed candidate is absent at the evidence revision", Path: item.CurrentPath, Severity: SeverityError})
			continue
		}
		verification := verifyV3BranchReviews(root, ledger, item, InventoryRecord{Path: item.CurrentPath, ActiveBranchTouch: append([]string{}, touches[item.CurrentPath]...)}, candidate, settings)
		problems = append(problems, verification.Problems...)
		evidence = append(evidence, verification.Evidence...)
	}
	branchDigest := ""
	remoteDigest := ""
	if ledger.BranchEvidence != nil {
		remoteDigest = ledger.BranchEvidence.RemoteOverlayDigest
		branchDigest = branchEvidenceSetDigest(evidence, remoteDigest)
		if branchDigest != ledger.BranchEvidence.Digest {
			problems = append(problems, Problem{Code: "branch_evidence_digest_mismatch", Message: "branch evidence aggregate moved after review", Severity: SeverityError})
		}
	}
	registerSHA, registerErr := regularSourceDigest(root, LegacyRegisterPath)
	if registerErr != nil {
		return "", problems, registerErr
	}
	envelope := struct {
		EvidenceRevision         string `json:"evidence_revision"`
		ActivationBaseRevision   string `json:"activation_base_revision"`
		LedgerPath               string `json:"ledger_path"`
		LedgerSHA256             string `json:"ledger_sha256"`
		PolicyDigest             string `json:"policy_digest"`
		CandidateSetDigest       string `json:"candidate_set_digest"`
		BranchEvidenceDigest     string `json:"branch_evidence_digest"`
		RemoteOverlayDigest      string `json:"remote_overlay_digest,omitempty"`
		ConfigPreconditionSHA256 string `json:"config_precondition_sha256"`
		RegisterSHA256           string `json:"register_sha256"`
	}{ledger.EvidenceRevision, checkpoint.ActivationBaseRevision, ledger.ArtifactPath, ledger.ArtifactSHA256, ledger.PolicyDigest, ledger.CandidateSetDigest, branchDigest, remoteDigest, ledger.TargetConfigPreconditionSHA256, registerSHA}
	data, marshalErr := json.Marshal(envelope)
	if marshalErr != nil {
		return "", problems, marshalErr
	}
	digest := sha256.Sum256(data)
	SortProblems(problems)
	return hex.EncodeToString(digest[:]), problems, nil
}
