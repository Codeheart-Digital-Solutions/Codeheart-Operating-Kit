package plancatalog

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type CatalogActivationPlan struct {
	Ledger                  MigrationLedger
	FilePlan                reconcile.Plan
	Problems                []Problem
	Projection              MigrationProjection
	BranchReviews           []BranchReviewResult
	FollowUps               []IncrementalFollowUp
	Checkpoint              ledgerCheckpoint
	ExternalAuthorityDigest string
	MigrationActionDigest   string
	ExpectedConfigSHA256    string
	MigrationActions        []reconcile.Action
	CheckpointActions       []reconcile.Action
}

type CatalogActivationOptions struct {
	DryRun             bool
	Now                time.Time
	Hook               reconcile.PhaseHook
	RemoteOverlayCheck func() (string, error)
}

type CatalogActivationOutcome struct {
	SchemaVersion                int                   `json:"schema_version"`
	EvidenceRevision             string                `json:"evidence_revision"`
	ActivationBaseRevision       string                `json:"activation_base_revision"`
	BranchEvidenceDigest         string                `json:"branch_evidence_digest"`
	MigrationActionDigest        string                `json:"migration_action_digest"`
	ProjectedCoverage            MigrationProjection   `json:"projected_coverage"`
	BranchReviews                []BranchReviewResult  `json:"branch_reviews"`
	IncrementalFollowUps         []IncrementalFollowUp `json:"incremental_follow_ups"`
	ActivationPerformed          bool                  `json:"activation_performed"`
	ActivationCheckpointRequired bool                  `json:"activation_checkpoint_required"`
	Result                       reconcile.Result      `json:"result"`
}

func BuildCatalogActivationPlan(root string, ledger MigrationLedger) (CatalogActivationPlan, error) {
	if ledger.SchemaVersion != 3 {
		return CatalogActivationPlan{}, fmt.Errorf("catalog_activation_requires_ledger_v3: guarded activation requires a reviewed schema-v3 ledger")
	}
	problems := []Problem{}
	checkpoint, checkpointProblems := verifyLedgerCheckpoint(root, ledger.EvidenceRevision, ledger.ArtifactPath, ledger.ArtifactSHA256)
	problems = append(problems, checkpointProblems...)
	externalDigest, externalProblems, err := currentV3ExternalAuthorityDigest(root, ledger)
	if err != nil {
		return CatalogActivationPlan{}, err
	}
	problems = append(problems, externalProblems...)
	migrationActions, projection, reviews, followUps, projectionProblems, err := projectV3MigrationActions(root, ledger)
	if err != nil {
		return CatalogActivationPlan{}, err
	}
	problems = append(problems, projectionProblems...)
	problems = append(problems, verifyMigrationWorktree(root, migrationActions)...)
	actionDigest := migrationActionDigest(migrationActions)
	configData, err := readRegularSource(root, state.ConfigPath)
	if err != nil {
		return CatalogActivationPlan{}, err
	}
	if sha256Text(configData) != ledger.TargetConfigPreconditionSHA256 {
		problems = append(problems, Problem{Code: "target_config_precondition_mismatch", Message: "config bytes changed after ledger review", Path: state.ConfigPath, Severity: SeverityError})
	}
	newConfig, err := buildActivatedConfig(configData, ledger, checkpoint.ActivationBaseRevision, actionDigest, projection)
	if err != nil {
		problems = append(problems, Problem{Code: "catalog_activation_config_invalid", Message: err.Error(), Path: state.ConfigPath, Severity: SeverityError})
	}
	configAction := reconcile.Action{Kind: "replace", Target: state.ConfigPath, Owner: "repo-config", Content: newConfig, Mode: 0o644, ExpectedSHA256: sha256Text(configData)}
	checkpointActions := append([]reconcile.Action{}, migrationActions...)
	checkpointActions = append(checkpointActions, configAction)
	observed, err := state.Inspect(root)
	if err != nil {
		return CatalogActivationPlan{}, err
	}
	filePlan, err := reconcile.BuildFilePlan("plans catalog-activate", root, string(observed.Classification), []state.Classification{observed.Classification}, []reconcile.Action{configAction})
	if err != nil {
		return CatalogActivationPlan{}, err
	}
	SortProblems(problems)
	for _, problem := range problems {
		if problem.Severity == SeverityError {
			filePlan.Blockers = append(filePlan.Blockers, reconcile.Blocker{Code: problem.Code, Message: problem.Message, Path: problem.Path, Remediation: problem.Remediation, RetryCommand: "plans catalog-activate --dry-run"})
		}
	}
	return CatalogActivationPlan{
		Ledger: ledger, FilePlan: filePlan, Problems: problems, Projection: projection,
		BranchReviews: reviews, FollowUps: followUps, Checkpoint: checkpoint,
		ExternalAuthorityDigest: externalDigest, MigrationActionDigest: actionDigest,
		ExpectedConfigSHA256: sha256Text(newConfig), MigrationActions: migrationActions,
		CheckpointActions: checkpointActions,
	}, nil
}

func ExecuteCatalogActivation(plan CatalogActivationPlan, options CatalogActivationOptions) (CatalogActivationOutcome, error) {
	remoteCheck := func() ([]reconcile.Blocker, error) {
		if plan.Ledger.BranchEvidence != nil && plan.Ledger.BranchEvidence.EvidenceScope == "remote-aware" {
			if options.RemoteOverlayCheck == nil {
				return []reconcile.Blocker{{Code: "remote_overlay_missing", Message: "remote-aware activation requires --remote-overlays"}}, nil
			}
			observed, overlayErr := options.RemoteOverlayCheck()
			if overlayErr != nil {
				return nil, overlayErr
			}
			if observed != plan.Ledger.ObservedRemoteOverlayDigest || observed != plan.Ledger.BranchEvidence.RemoteOverlayDigest {
				return []reconcile.Blocker{{Code: "remote_overlay_stale", Message: "fresh remote overlay differs from reviewed activation evidence"}}, nil
			}
		}
		return nil, nil
	}
	preCheck := func() ([]reconcile.Blocker, error) {
		if blockers, remoteErr := remoteCheck(); remoteErr != nil || len(blockers) > 0 {
			return blockers, remoteErr
		}
		digest, problems, err := currentV3ExternalAuthorityDigest(plan.FilePlan.Root, plan.Ledger)
		if err != nil {
			return nil, err
		}
		if HasErrors(problems) || digest != plan.ExternalAuthorityDigest || HasErrors(verifyMigrationWorktree(plan.FilePlan.Root, plan.MigrationActions)) {
			return []reconcile.Blocker{{Code: "catalog_activation_authority_drift", Message: "ledger, checkpoint, candidate, ref, proof, overlay, config, or migration output moved after activation planning", Remediation: "reinventory and review before activation", RetryCommand: "plans catalog-activate --dry-run"}}, nil
		}
		return nil, nil
	}
	postCheck := func() ([]reconcile.Blocker, error) {
		if blockers, remoteErr := remoteCheck(); remoteErr != nil || len(blockers) > 0 {
			return blockers, remoteErr
		}
		digest, problems, err := currentV3ExternalAuthorityDigestForConfig(plan.FilePlan.Root, plan.Ledger, plan.ExpectedConfigSHA256)
		if err != nil {
			return nil, err
		}
		if HasErrors(problems) || digest != plan.ExternalAuthorityDigest || HasErrors(verifyMigrationWorktree(plan.FilePlan.Root, plan.CheckpointActions)) {
			return []reconcile.Blocker{{Code: "catalog_activation_authority_drift", Message: "activation authority moved after the config write", Remediation: "preserve transaction evidence and reinventory", RetryCommand: "plans catalog-activate --dry-run"}}, nil
		}
		return nil, nil
	}
	var result reconcile.Result
	var err error
	if options.DryRun {
		filePlan := plan.FilePlan
		blockers, checkErr := preCheck()
		if checkErr != nil {
			return CatalogActivationOutcome{}, checkErr
		}
		filePlan.Blockers = append(filePlan.Blockers, blockers...)
		result = reconcile.Preview(filePlan)
	} else {
		result, err = reconcile.Apply(plan.FilePlan, reconcile.ApplyOptions{Now: options.Now, Hook: options.Hook, PreWriteAuthorityCheck: preCheck, PostWriteAuthorityCheck: postCheck})
	}
	digest := ""
	if plan.Ledger.BranchEvidence != nil {
		digest = plan.Ledger.BranchEvidence.Digest
	}
	return CatalogActivationOutcome{
		SchemaVersion: 3, EvidenceRevision: plan.Ledger.EvidenceRevision,
		ActivationBaseRevision: plan.Checkpoint.ActivationBaseRevision,
		BranchEvidenceDigest:   digest, MigrationActionDigest: plan.MigrationActionDigest,
		ProjectedCoverage: plan.Projection, BranchReviews: append([]BranchReviewResult{}, plan.BranchReviews...),
		IncrementalFollowUps:         append([]IncrementalFollowUp{}, plan.FollowUps...),
		ActivationPerformed:          !options.DryRun && result.Status == reconcile.StatusSucceeded,
		ActivationCheckpointRequired: !options.DryRun && result.Status == reconcile.StatusSucceeded,
		Result:                       result,
	}, err
}

func projectV3MigrationActions(root string, ledger MigrationLedger) ([]reconcile.Action, MigrationProjection, []BranchReviewResult, []IncrementalFollowUp, []Problem, error) {
	settings, settingProblems := LoadRepositorySettings(root)
	settings.DiscoveryVersion = DiscoveryV2
	settings.Mode = ledger.TargetCatalogMode
	settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	classification, err := ClassifyCommitTree(root, ledger.EvidenceRevision, settings, ledger.TargetCatalogMode, ledger.RepositoryID)
	if err != nil {
		return nil, MigrationProjection{}, nil, nil, settingProblems, err
	}
	candidates := map[string]Candidate{}
	for _, candidate := range classification.Candidates {
		if candidate.Ownership == OwnershipOwned {
			candidates[candidate.Path] = candidate
		}
	}
	touches, touchProblems := activeBranchTouchesWithSettings(root, settings)
	problems := append([]Problem{}, settingProblems...)
	problems = append(problems, touchProblems...)
	actions := []reconcile.Action{}
	reviews := []BranchReviewResult{}
	branchEvidence := []BranchCandidateEvidence{}
	followUps := []IncrementalFollowUp{}
	projection := MigrationProjection{Candidates: len(candidates)}
	seen := map[string]bool{}
	for _, item := range ledger.Records {
		candidate, ok := candidates[item.CurrentPath]
		if !ok {
			problems = append(problems, Problem{Code: "migration_candidate_missing", Message: "reviewed candidate is absent from the evidence revision", Path: item.CurrentPath, Severity: SeverityError})
			continue
		}
		seen[item.CurrentPath] = true
		verification := verifyV3BranchReviews(root, ledger, item, InventoryRecord{Path: item.CurrentPath, ActiveBranchTouch: touches[item.CurrentPath]}, candidate, settings)
		reviews = append(reviews, verification.Results...)
		branchEvidence = append(branchEvidence, verification.Evidence...)
		problems = append(problems, verification.Problems...)
		if verification.FollowUp != nil {
			followUps = append(followUps, *verification.FollowUp)
		}
		if verification.Deferred {
			projection.Grandfathered++
			continue
		}
		if verification.Resolved {
			projection.Canonical++
			continue
		}
		if item.Deferred || HasErrors(verification.Problems) {
			continue
		}
		source, sourceErr := gitBytes(root, "show", ledger.EvidenceRevision+":"+item.CurrentPath)
		if sourceErr != nil || sha256Text(source) != item.SourceSHA256 {
			message := "reviewed source bytes are unavailable at the evidence revision"
			if sourceErr != nil {
				message = sourceErr.Error()
			}
			problems = append(problems, Problem{Code: "source_revision_mismatch", Message: message, Path: item.CurrentPath, Severity: SeverityError})
			continue
		}
		metadata := metadataFromDecision(item)
		targetKind, ok := migrationV2TargetKind(item.TargetPath, metadata)
		if !ok || targetKind != metadata.Kind {
			problems = append(problems, Problem{Code: "invalid_target_placement", Message: "target does not match the reviewed metadata kind", Path: item.TargetPath, Severity: SeverityError})
			continue
		}
		updated := source
		existing, parseErr := ParseDocument(item.CurrentPath, source, candidate.ExpectedKind)
		if parseErr == nil && existing.Metadata != nil {
			if !metadataEqual(*existing.Metadata, metadata) {
				problems = append(problems, Problem{Code: "metadata_conflict", Message: "evidence source metadata differs from the reviewed decision", Path: item.CurrentPath, Severity: SeverityError})
				continue
			}
		} else if parseErr != nil && ErrorCode(parseErr) == "metadata_missing" {
			updated, err = InsertMetadata(source, metadata)
			if err != nil {
				return nil, projection, reviews, followUps, problems, err
			}
		} else if parseErr != nil {
			problems = append(problems, Problem{Code: "malformed_metadata", Message: parseErr.Error(), Path: item.CurrentPath, Severity: SeverityError})
			continue
		}
		mode := migrationActionMode(candidate.Provenance.Source.Mode)
		if item.CurrentPath == item.TargetPath {
			if !bytes.Equal(source, updated) {
				actions = append(actions, reconcile.Action{Kind: "replace", Target: item.CurrentPath, Owner: "repo-plan", Content: updated, Mode: mode, ExpectedSHA256: item.SourceSHA256})
			}
		} else {
			actions = append(actions, reconcile.Action{Kind: "create", Target: item.TargetPath, Owner: "repo-plan", Content: updated, Mode: mode}, reconcile.Action{Kind: "remove", Target: item.CurrentPath, Owner: "repo-plan", ExpectedSHA256: item.SourceSHA256})
		}
		projection.Canonical++
	}
	for candidatePath := range candidates {
		if !seen[candidatePath] {
			problems = append(problems, Problem{Code: "migration_candidate_unreviewed", Message: "authoritative candidate is absent from the reviewed ledger", Path: candidatePath, Severity: SeverityError})
		}
	}
	projection.Gaps = projection.Candidates - projection.Canonical - projection.Grandfathered
	projection.MixedCoverageComplete = projection.Gaps == 0 && !HasErrors(problems)
	projection.CanonicalReady = projection.MixedCoverageComplete && projection.Grandfathered == 0
	projection.Ready = projection.MixedCoverageComplete
	if ledger.BranchEvidence == nil || branchEvidenceSetDigest(branchEvidence, func() string {
		if ledger.BranchEvidence == nil {
			return ""
		}
		return ledger.BranchEvidence.RemoteOverlayDigest
	}()) != func() string {
		if ledger.BranchEvidence == nil {
			return ""
		}
		return ledger.BranchEvidence.Digest
	}() {
		problems = append(problems, Problem{Code: "branch_evidence_digest_mismatch", Message: "recomputed branch evidence differs from the reviewed aggregate", Severity: SeverityError})
		projection.Ready = false
		projection.MixedCoverageComplete = false
		projection.CanonicalReady = false
	}
	if ledger.TargetCatalogMode == ModeCanonical && projection.Grandfathered > 0 {
		problems = append(problems, Problem{Code: "canonical_deferral_forbidden", Message: "canonical activation cannot retain deferred owners", Severity: SeverityError})
	}
	sort.SliceStable(actions, func(i, j int) bool {
		if actions[i].Target != actions[j].Target {
			return actions[i].Target < actions[j].Target
		}
		return actions[i].Kind < actions[j].Kind
	})
	SortProblems(problems)
	return actions, projection, reviews, followUps, problems, nil
}

func buildActivatedConfig(configData []byte, ledger MigrationLedger, activationBase, actionDigest string, projection MigrationProjection) ([]byte, error) {
	config, err := state.DecodeAndValidateConfigYAML(configData)
	if err != nil {
		return nil, err
	}
	components := state.Map(config["component_settings"])
	if components == nil {
		components = map[string]any{}
		config["component_settings"] = components
	}
	planning := state.Map(components["planning-workflows"])
	if planning == nil {
		planning = map[string]any{}
		components["planning-workflows"] = planning
	}
	planning["plan_catalog_mode"] = string(ledger.TargetCatalogMode)
	planning["plan_catalog_discovery_version"] = 2
	if ledger.TargetCatalogMode == ModeMixed && ledger.BranchEvidence != nil && projection.Grandfathered > 0 {
		config["schema_version"] = 2
		planning["plan_catalog_migration_evidence"] = map[string]any{
			"ledger_path": ledger.ArtifactPath, "ledger_sha256": ledger.ArtifactSHA256,
			"branch_evidence_digest":   ledger.BranchEvidence.Digest,
			"evidence_revision":        ledger.EvidenceRevision,
			"activation_base_revision": activationBase,
			"migration_action_digest":  actionDigest,
			"evidence_scope":           ledger.BranchEvidence.EvidenceScope,
		}
	} else {
		config["schema_version"] = 1
		delete(planning, "plan_catalog_migration_evidence")
	}
	if err := state.ValidateConfig(config); err != nil {
		return nil, err
	}
	return state.EncodeYAML(config)
}
