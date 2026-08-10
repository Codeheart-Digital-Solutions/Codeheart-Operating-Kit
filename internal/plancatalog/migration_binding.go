package plancatalog

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

func boundDeferredPaths(root string, binding *MigrationEvidenceBinding) map[string]bool {
	paths := map[string]bool{}
	if binding == nil {
		return paths
	}
	data, err := readRegularSource(root, binding.LedgerPath)
	if err != nil || sha256Text(data) != binding.LedgerSHA256 {
		return paths
	}
	ledger, err := LoadMigrationLedger(data)
	if err != nil || ledger.SchemaVersion != 3 {
		return paths
	}
	for _, record := range ledger.Records {
		if record.Deferred && !deferredRecordHasIntegratedMetadata(root, record) {
			paths[record.CurrentPath] = true
		}
	}
	return paths
}

func deferredRecordHasIntegratedMetadata(root string, record MigrationRecord) bool {
	head, headErr := gitText(root, "rev-parse", "--verify", "HEAD^{commit}")
	if headErr != nil {
		return false
	}
	data, err := gitBytes(root, "show", head+":"+record.CurrentPath)
	state, blocker := ReadBranchPathState(root, head, record.CurrentPath)
	if err != nil || blocker != nil {
		return false
	}
	ownerTip, matched := IntegratedDeferredOwnerTip(record, state, data)
	if !matched {
		return false
	}
	runner, runnerBlocker := newBranchGitRunner(root)
	if runnerBlocker != nil {
		return false
	}
	defer runner.close()
	merged, ancestryBlocker := runner.isAncestor(ownerTip, head)
	return ancestryBlocker == nil && merged
}

// IntegratedDeferredOwnerTip verifies the exact committed path state and
// canonical metadata required to resolve a deferred owner. Callers must also
// prove that the returned immutable owner tip is reachable from their target
// revision.
func IntegratedDeferredOwnerTip(record MigrationRecord, state BranchPathState, data []byte) (string, bool) {
	var owner *OwnerTipCandidate
	var ownerTip string
	for _, review := range record.BranchTouchReviews {
		if review.Disposition != BranchDispositionDeferredActiveOwner {
			continue
		}
		if owner != nil || review.OwnerTipCandidate == nil {
			return "", false
		}
		owner = review.OwnerTipCandidate
		ownerTip = review.RefTip
	}
	if owner == nil || owner.Path != record.CurrentPath || state.State != BranchPathPresent || state.Mode != owner.Mode || state.ObjectID != owner.ObjectID || state.SHA256 != owner.SHA256 || sha256Text(data) != owner.SHA256 {
		return "", false
	}
	kind, ok := kindForFilename(record.CurrentPath)
	if !ok {
		return "", false
	}
	parsed, err := ParseDocument(record.CurrentPath, data, kind)
	return ownerTip, err == nil && parsed.Metadata != nil && metadataEqual(*parsed.Metadata, metadataFromDecision(record))
}

func BoundRemoteOverlayDigest(root string, binding *MigrationEvidenceBinding) string {
	if binding == nil || binding.EvidenceScope != "remote-aware" {
		return ""
	}
	data, err := readRegularSource(root, binding.LedgerPath)
	if err != nil || sha256Text(data) != binding.LedgerSHA256 {
		return ""
	}
	ledger, err := LoadMigrationLedger(data)
	if err != nil || ledger.BranchEvidence == nil {
		return ""
	}
	return ledger.BranchEvidence.RemoteOverlayDigest
}

func validatePersistedMigrationEvidence(root string, binding *MigrationEvidenceBinding, observedRemoteOverlayDigest string, observedRemoteBranchEvidence []BranchCandidateEvidence) []Problem {
	if binding == nil {
		return []Problem{{Code: "migration_evidence_binding_missing", Message: "config schema v2 requires a migration evidence binding", Path: state.ConfigPath, Severity: SeverityError}}
	}
	problems := []Problem{}
	ledgerData, err := readRegularSource(root, binding.LedgerPath)
	if err != nil || sha256Text(ledgerData) != binding.LedgerSHA256 {
		message := "bound migration ledger is missing or changed"
		if err != nil {
			message = err.Error()
		}
		return []Problem{{Code: "migration_evidence_ledger_mismatch", Message: message, Path: binding.LedgerPath, Severity: SeverityError, Remediation: "restore the exact reviewed ledger or perform a new reviewed migration"}}
	}
	ledger, err := LoadMigrationLedger(ledgerData)
	if err != nil || ledger.SchemaVersion != 3 {
		message := "bound migration ledger is not a valid schema-v3 artifact"
		if err != nil {
			message = err.Error()
		}
		return []Problem{{Code: "migration_evidence_ledger_invalid", Message: message, Path: binding.LedgerPath, Severity: SeverityError}}
	}
	ledger.ArtifactPath = binding.LedgerPath
	ledger.ArtifactSHA256 = binding.LedgerSHA256
	ledger.ObservedRemoteOverlayDigest = observedRemoteOverlayDigest
	ledger.ObservedRemoteBranchEvidence = append([]BranchCandidateEvidence{}, observedRemoteBranchEvidence...)
	if ledger.EvidenceRevision != binding.EvidenceRevision || ledger.BranchEvidence == nil || ledger.BranchEvidence.Digest != binding.BranchEvidenceDigest || ledger.BranchEvidence.EvidenceScope != binding.EvidenceScope {
		problems = append(problems, Problem{Code: "migration_evidence_binding_mismatch", Message: "config evidence fields do not match the reviewed ledger", Path: state.ConfigPath, Severity: SeverityError})
	}
	if binding.EvidenceScope == "remote-aware" {
		expected := ""
		if ledger.BranchEvidence != nil {
			expected = ledger.BranchEvidence.RemoteOverlayDigest
		}
		switch {
		case observedRemoteOverlayDigest == "":
			problems = append(problems, Problem{Code: "remote_overlay_missing", Message: "remote-aware migration evidence requires a fresh target-member overlay", Severity: SeverityError, Remediation: "rerun the command with --remote-overlays"})
		case expected == "" || observedRemoteOverlayDigest != expected:
			problems = append(problems, Problem{Code: "remote_overlay_stale", Message: "fresh target-member overlay differs from the bound reviewed evidence", Severity: SeverityError, Remediation: "reinventory and review the current remote evidence"})
		}
	}
	parent, parentErr := gitText(root, "rev-parse", "--verify", binding.ActivationBaseRevision+"^1")
	if parentErr != nil || parent != binding.EvidenceRevision {
		problems = append(problems, Problem{Code: "ledger_checkpoint_parent_mismatch", Message: "activation base is not the exact first-parent ledger checkpoint", Path: binding.LedgerPath, Severity: SeverityError})
	}
	ledgerDelta, deltaErr := gitNULFields(root, "diff-tree", "--no-commit-id", "--name-only", "-z", "-r", binding.EvidenceRevision, binding.ActivationBaseRevision)
	if deltaErr != nil || !equalStrings(ledgerDelta, []string{binding.LedgerPath}) {
		problems = append(problems, Problem{Code: "ledger_checkpoint_delta_mismatch", Message: "ledger checkpoint contains paths other than the exact bound ledger", Path: binding.LedgerPath, Severity: SeverityError})
	}
	if committedLedger, readErr := gitBytes(root, "show", binding.ActivationBaseRevision+":"+binding.LedgerPath); readErr != nil || sha256Text(committedLedger) != binding.LedgerSHA256 {
		problems = append(problems, Problem{Code: "ledger_checkpoint_blob_mismatch", Message: "ledger checkpoint blob differs from the bound ledger", Path: binding.LedgerPath, Severity: SeverityError})
	}
	activationRevision, activationErr := BoundActivationCheckpointRevision(root, binding.ActivationBaseRevision, "HEAD", binding.MigrationActionDigest)
	if activationErr != nil {
		code := "activation_checkpoint_missing"
		if strings.HasPrefix(activationErr.Error(), "activation_checkpoint_ambiguous:") {
			code = "activation_checkpoint_ambiguous"
		}
		problems = append(problems, Problem{Code: code, Message: strings.TrimSpace(strings.TrimPrefix(activationErr.Error(), code+":")), Path: state.ConfigPath, Severity: SeverityError, Remediation: "incorporate exactly one reviewed activation checkpoint without rewriting its E-to-L-to-A chronology"})
		SortProblems(problems)
		return problems
	}
	activationParent, err := gitText(root, "rev-parse", "--verify", activationRevision+"^1")
	if err != nil || activationParent != binding.ActivationBaseRevision {
		problems = append(problems, Problem{Code: "activation_checkpoint_parent_mismatch", Message: "activation checkpoint is not the immediate first-parent child of the ledger checkpoint", Severity: SeverityError})
	}
	if configErr := ValidateGuardedActivationConfig(root, ledger, binding.ActivationBaseRevision, activationRevision, binding.MigrationActionDigest); configErr != nil {
		problems = append(problems, Problem{Code: "activation_checkpoint_config_mismatch", Message: strings.TrimSpace(strings.TrimPrefix(configErr.Error(), "activation_checkpoint_config_mismatch:")), Path: state.ConfigPath, Severity: SeverityError})
	}
	actions, projection, _, _, projectionProblems, projectErr := projectV3MigrationActions(root, ledger)
	if projectErr != nil {
		problems = append(problems, Problem{Code: "migration_evidence_recompute_failed", Message: projectErr.Error(), Severity: SeverityError})
	} else {
		problems = append(problems, projectionProblems...)
		if migrationActionDigest(actions) != binding.MigrationActionDigest {
			problems = append(problems, Problem{Code: "migration_action_digest_mismatch", Message: "current reviewed projection differs from the bound activation action digest", Severity: SeverityError})
		}
		expectedPaths := []string{state.ConfigPath}
		for _, action := range actions {
			expectedPaths = append(expectedPaths, action.Target)
		}
		expectedPaths = uniqueSortedStrings(expectedPaths)
		actualPaths, diffErr := gitNULFields(root, "diff-tree", "--no-commit-id", "--name-only", "-z", "-r", binding.ActivationBaseRevision, activationRevision)
		sort.Strings(actualPaths)
		if diffErr != nil || !equalStrings(actualPaths, expectedPaths) {
			problems = append(problems, Problem{Code: "activation_checkpoint_delta_mismatch", Message: fmt.Sprintf("activation checkpoint paths %v differ from bound projection %v", actualPaths, expectedPaths), Severity: SeverityError})
		}
		for _, action := range actions {
			if action.Kind == "remove" {
				if _, showErr := gitBytes(root, "show", activationRevision+":"+action.Target); showErr == nil {
					problems = append(problems, Problem{Code: "activation_checkpoint_blob_mismatch", Message: "reviewed removal remains present in the activation checkpoint", Path: action.Target, Severity: SeverityError})
				}
				continue
			}
			checkpointBytes, showErr := gitBytes(root, "show", activationRevision+":"+action.Target)
			if showErr != nil || !bytes.Equal(checkpointBytes, action.Content) {
				problems = append(problems, Problem{Code: "activation_checkpoint_blob_mismatch", Message: "activation checkpoint plan bytes differ from the reviewed projection", Path: action.Target, Severity: SeverityError})
			}
		}
		if !projection.MixedCoverageComplete {
			problems = append(problems, Problem{Code: "mixed_coverage_incomplete", Message: "bound migration no longer covers every reviewed candidate", Severity: SeverityError})
		}
	}
	currentConfig, configErr := readRegularSource(root, state.ConfigPath)
	activationConfig, activationConfigErr := gitBytes(root, "show", activationRevision+":"+state.ConfigPath)
	if configErr != nil || activationConfigErr != nil || !bytes.Equal(currentConfig, activationConfig) {
		problems = append(problems, Problem{Code: "migration_evidence_config_moved", Message: "current config differs from the guarded activation checkpoint", Path: state.ConfigPath, Severity: SeverityError})
	}
	settings, settingProblems := DecodeRepositorySettings(currentConfig)
	problems = append(problems, settingProblems...)
	head, headErr := gitText(root, "rev-parse", "--verify", "HEAD^{commit}")
	if headErr != nil {
		problems = append(problems, Problem{Code: "migration_evidence_revision_unavailable", Message: "current committed revision is unavailable", Severity: SeverityError})
	} else if !HasErrors(settingProblems) {
		resolvedDeferred := map[string]bool{}
		for _, record := range ledger.Records {
			if !record.Deferred {
				continue
			}
			data, readErr := gitBytes(root, "show", head+":"+record.CurrentPath)
			pathState, blocker := ReadBranchPathState(root, head, record.CurrentPath)
			ownerTip, integrated := IntegratedDeferredOwnerTip(record, pathState, data)
			if readErr == nil && blocker == nil && integrated {
				runner, runnerBlocker := newBranchGitRunner(root)
				if runnerBlocker == nil {
					merged, ancestryBlocker := runner.isAncestor(ownerTip, head)
					runner.close()
					if ancestryBlocker == nil && merged {
						resolvedDeferred[record.CurrentPath] = true
					}
				}
			}
		}
		activationClassification, activationClassifyErr := ClassifyCommitTree(root, activationRevision, settings, settings.Mode, settings.RepositoryID)
		currentClassification, currentClassifyErr := ClassifyCommitTree(root, head, settings, settings.Mode, settings.RepositoryID)
		if activationClassifyErr != nil || currentClassifyErr != nil || !activationClassification.Complete || !currentClassification.Complete {
			problems = append(problems, Problem{Code: "migration_evidence_candidate_authority_unavailable", Message: "activated or current candidate authority is incomplete", Severity: SeverityError})
		} else if activationClassification.PolicyDigest != ledger.PolicyDigest || currentClassification.PolicyDigest != ledger.PolicyDigest || candidateDigest(candidatesWithoutPaths(activationClassification.Candidates, resolvedDeferred)) != candidateDigest(candidatesWithoutPaths(currentClassification.Candidates, resolvedDeferred)) {
			problems = append(problems, Problem{Code: "migration_evidence_candidate_set_moved", Message: "current candidate authority differs from the activation checkpoint", Severity: SeverityError, Remediation: "reinventory and perform the required incremental migration"})
		}
	}
	if clean, cleanErr := gitIndexMatchesHead(root); cleanErr != nil || !clean {
		problems = append(problems, Problem{Code: "migration_evidence_index_dirty", Message: "index differs from the validated current revision", Severity: SeverityError})
	}
	if clean, cleanErr := gitWorktreeClean(root); cleanErr != nil || !clean {
		problems = append(problems, Problem{Code: "migration_evidence_worktree_dirty", Message: "worktree differs from the validated current revision", Severity: SeverityError, Remediation: "commit or move unrelated work before validating persisted migration evidence"})
	}
	SortProblems(problems)
	return problems
}

func candidatesWithoutPaths(candidates []Candidate, excluded map[string]bool) []Candidate {
	result := make([]Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		if !excluded[candidate.Path] {
			result = append(result, candidate)
		}
	}
	return result
}
