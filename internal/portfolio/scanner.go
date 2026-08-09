package portfolio

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/plancatalog"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type ScanOptions struct {
	Root           string
	Config         Config
	Sources        []Source
	Runner         CommandRunner
	Now            time.Time
	MaxConcurrency int
	WriteCache     bool
	PlanTarget     *RemotePlanTarget
}

type RemotePlanTarget struct {
	RepositoryID           string
	EvidenceRevision       string
	ActivationBaseRevision string
	ActivationRevision     string
	LedgerPath             string
	LedgerSHA256           string
	DiscoveryVersion       plancatalog.DiscoveryVersion
	TargetCatalogMode      plancatalog.CatalogMode
	PolicyDigest           string
	CandidateSetDigest     string
	MigrationActionDigest  string
	BranchEvidenceDigest   string
	EvidenceScope          string
	RemoteOverlayDigest    string
	RemoteSourceIdentity   string
}

type repositoryScan struct {
	Source                    RepositorySource
	Member                    *plancatalog.CatalogMember
	Observations              []plancatalog.SourceObservation
	CompatibilityObservations []plancatalog.CompatibilityObservation
	PlanCandidates            []PlanCandidateObservation
	RemoteBranchEvidence      []plancatalog.BranchCandidateEvidence
	Candidate                 *Candidate
	Errors                    []ScanError
	APICalls                  int
}

func Scan(ctx context.Context, options ScanOptions) (ScanResult, error) {
	started := options.Now
	if started.IsZero() {
		started = time.Now()
	}
	started = started.UTC().Truncate(time.Second)
	if options.MaxConcurrency <= 0 || options.MaxConcurrency > 4 {
		options.MaxConcurrency = 4
	}
	if options.Runner == nil {
		options.Runner = ExecRunner{}
	}
	if options.Config.Root == "" {
		config, err := LoadConfig(options.Root)
		if err != nil {
			return ScanResult{}, err
		}
		options.Config = config
	}
	if options.Config.CoordinationHomeID == "" {
		return ScanResult{}, fmt.Errorf("portfolio_home_identity_missing: coordination_home_id is required for remote scanning")
	}
	scanLock, err := acquireScanLock(options.Config.Root)
	if err != nil {
		return ScanResult{}, err
	}
	defer scanLock.Close()
	overlay, err := ReadOverlay(options.Config.Root, options.Config.CoordinationHomeID)
	if err != nil {
		return ScanResult{}, err
	}
	defer overlay.Close()
	catalog := Catalog{
		SchemaVersion:             3,
		DiscoveryVersion:          plancatalog.DiscoveryV2,
		CoordinationHomeID:        options.Config.CoordinationHomeID,
		StartedAt:                 started.Format(time.RFC3339),
		Complete:                  true,
		MixedCoverageComplete:     true,
		CanonicalReady:            true,
		Members:                   []plancatalog.CatalogMember{},
		Observations:              []plancatalog.SourceObservation{},
		CompatibilityObservations: []plancatalog.CompatibilityObservation{},
		PlanCandidates:            []PlanCandidateObservation{},
		Candidates:                []Candidate{},
		Errors:                    []ScanError{},
		Metrics:                   Metrics{MaxConcurrency: options.MaxConcurrency},
	}
	repositorySources := []RepositorySource{}
	for _, source := range options.Sources {
		catalog.Metrics.SourceCount++
		discovery, discoverErr := source.Discover(ctx)
		catalog.Metrics.APICallCount += discovery.APICalls
		if discoverErr != nil {
			catalog.Complete = false
			catalog.Errors = append(catalog.Errors, ScanError{Code: "source_discovery_failed", Message: discoverErr.Error(), SourceLocator: source.Kind(), Retryable: true})
			continue
		}
		repositorySources = append(repositorySources, discovery.Repositories...)
		catalog.Errors = append(catalog.Errors, discovery.Errors...)
		if !discovery.Complete {
			catalog.Complete = false
		}
	}
	self, selfErr := discoverSelfRepository(ctx, options.Config, options.Runner)
	catalog.Metrics.SourceCount++
	if selfErr != nil {
		catalog.Complete = false
		catalog.Errors = append(catalog.Errors, ScanError{Code: "self_remote_unavailable", Message: selfErr.Error(), RepositoryID: options.Config.RepositoryID})
	} else {
		repositorySources = append(repositorySources, self)
	}
	repositorySources = deduplicateSources(repositorySources)
	manager := MirrorManager{RepositoryRoot: options.Config.Root, Root: filepath.Join(options.Config.Root, filepath.FromSlash(MirrorRootPath)), Runner: options.Runner}
	results := make([]repositoryScan, len(repositorySources))
	semaphore := make(chan struct{}, options.MaxConcurrency)
	var wait sync.WaitGroup
	for index, source := range repositorySources {
		wait.Add(1)
		go func(index int, source RepositorySource) {
			defer wait.Done()
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				results[index] = repositoryScan{Source: source, Errors: []ScanError{{Code: "scan_cancelled", Message: "repository scan was cancelled", SourceLocator: source.Locator, Retryable: true}}}
				return
			}
			results[index] = scanRepository(ctx, manager, source, options.Config, options.PlanTarget, started)
		}(index, source)
	}
	wait.Wait()
	if err := scanLock.Verify(); err != nil {
		return ScanResult{}, err
	}
	seenMembers := map[string]bool{}
	for _, result := range results {
		catalog.Metrics.APICallCount += result.APICalls
		if len(result.Errors) > 0 {
			catalog.Complete = false
			catalog.Errors = append(catalog.Errors, result.Errors...)
		}
		if result.Candidate != nil {
			catalog.Candidates = append(catalog.Candidates, *result.Candidate)
		}
		if result.Member == nil {
			continue
		}
		if seenMembers[result.Member.RepositoryID] {
			catalog.Complete = false
			catalog.Errors = append(catalog.Errors, ScanError{Code: "duplicate_member_repository_id", Message: "multiple remote repositories declare the same stable member identity", SourceLocator: result.Source.Locator, RepositoryID: result.Member.RepositoryID})
			continue
		}
		seenMembers[result.Member.RepositoryID] = true
		if result.Member.Complete == nil || !*result.Member.Complete {
			catalog.Complete = false
		}
		if !result.Member.MixedCoverageComplete {
			catalog.MixedCoverageComplete = false
		}
		if !result.Member.CanonicalReady {
			catalog.CanonicalReady = false
		}
		catalog.Members = append(catalog.Members, *result.Member)
		catalog.Observations = append(catalog.Observations, result.Observations...)
		catalog.CompatibilityObservations = append(catalog.CompatibilityObservations, result.CompatibilityObservations...)
		catalog.PlanCandidates = append(catalog.PlanCandidates, result.PlanCandidates...)
		catalog.RemoteBranchEvidence = append(catalog.RemoteBranchEvidence, result.RemoteBranchEvidence...)
	}
	catalog.PolicyDigest = aggregateMemberDigest(catalog.Members, func(member plancatalog.CatalogMember) string { return member.PolicyDigest })
	catalog.CandidateSetDigest = aggregateMemberDigest(catalog.Members, func(member plancatalog.CatalogMember) string { return member.CandidateSetDigest })
	markObservationConflicts(catalog.Observations)
	completed := time.Now().UTC().Truncate(time.Second)
	if !options.Now.IsZero() {
		completed = started
	}
	catalog.CompletedAt = completed.Format(time.RFC3339)
	if catalog.Complete {
		catalog.LastCompleteScanAt = catalog.CompletedAt
	} else if previous, readErr := ReadCachedCatalog(options.Config.Root); readErr == nil && previous.SchemaVersion == 3 && previous.DiscoveryVersion == plancatalog.DiscoveryV2 && previous.Complete {
		catalog.LastCompleteScanAt = previous.LastCompleteScanAt
		if catalog.LastCompleteScanAt == "" {
			catalog.LastCompleteScanAt = previous.CompletedAt
		}
	}
	catalog.Metrics.DurationMS = completed.Sub(started).Milliseconds()
	if catalog.Metrics.DurationMS < 0 {
		catalog.Metrics.DurationMS = 0
	}
	catalog.Metrics.MemberCount = len(catalog.Members)
	catalog.Metrics.CandidateCount = len(catalog.Candidates)
	catalog.Metrics.ObservationCount = len(catalog.Observations)
	catalog.Metrics.CompatibilityObservationCount = len(catalog.CompatibilityObservations)
	if !catalog.Complete {
		catalog.MixedCoverageComplete = false
		catalog.CanonicalReady = false
	}
	for _, observation := range catalog.Observations {
		if observation.Stale {
			catalog.Metrics.StaleCount++
		}
	}
	sortCatalog(&catalog)
	schemaPath, err := state.SchemaForPlanCatalogVersion(catalog.SchemaVersion)
	if err != nil {
		return ScanResult{}, fmt.Errorf("portfolio_catalog_invalid: %w", err)
	}
	if err := state.Validate(schemaPath, catalog); err != nil {
		return ScanResult{}, fmt.Errorf("portfolio_catalog_invalid: %w", err)
	}
	if err := overlay.VerifyUnchanged(); err != nil {
		return ScanResult{}, err
	}
	if err := scanLock.Verify(); err != nil {
		return ScanResult{}, err
	}
	result := ScanResult{Catalog: catalog, LocalOnlyOmitted: true}
	if options.WriteCache {
		updated, preserved, err := StoreCompleteCatalog(options.Config.Root, catalog)
		if err != nil {
			return ScanResult{}, err
		}
		result.CacheUpdated = updated
		result.PreviousPreserved = preserved
	}
	if err := overlay.VerifyUnchanged(); err != nil {
		return ScanResult{}, err
	}
	if err := scanLock.Verify(); err != nil {
		return ScanResult{}, err
	}
	return result, nil
}

func discoverSelfRepository(ctx context.Context, config Config, runner CommandRunner) (RepositorySource, error) {
	result, err := runner.Run(ctx, "git", safeGitArgs("-C", config.Root, "config", "--get", "remote.origin.url")...)
	if err != nil || strings.TrimSpace(string(result.Stdout)) == "" {
		return RepositorySource{}, fmt.Errorf("coordination repository has no readable origin remote")
	}
	cloneURL := normalizeLocalRemote(config.Root, strings.TrimSpace(string(result.Stdout)))
	return RepositorySource{Kind: "local-git", Locator: config.Root, CloneURL: cloneURL, NameHint: config.RepositoryID, Self: true}, nil
}

func deduplicateSources(sources []RepositorySource) []RepositorySource {
	sort.SliceStable(sources, func(i, j int) bool {
		if sources[i].Self != sources[j].Self {
			return sources[i].Self
		}
		if sources[i].Kind != sources[j].Kind {
			return sources[i].Kind < sources[j].Kind
		}
		return sources[i].Locator < sources[j].Locator
	})
	seen := map[string]int{}
	result := []RepositorySource{}
	for _, source := range sources {
		key := canonicalRemoteIdentity(source.CloneURL)
		if key == "" {
			key = source.Kind + "\x00" + source.Locator
		}
		if index, found := seen[key]; found {
			if result[index].Enricher == nil && source.Enricher != nil {
				result[index].Enricher = source.Enricher
			}
			continue
		}
		seen[key] = len(result)
		result = append(result, source)
	}
	return result
}

func scanRepository(ctx context.Context, manager MirrorManager, source RepositorySource, home Config, target *RemotePlanTarget, now time.Time) repositoryScan {
	result := repositoryScan{Source: source, Observations: []plancatalog.SourceObservation{}, CompatibilityObservations: []plancatalog.CompatibilityObservation{}, PlanCandidates: []PlanCandidateObservation{}, RemoteBranchEvidence: []plancatalog.BranchCandidateEvidence{}, Errors: []ScanError{}}
	repository, err := manager.Refresh(ctx, source)
	if err != nil {
		result.Errors = append(result.Errors, ScanError{Code: "mirror_refresh_failed", Message: err.Error(), SourceLocator: source.Locator, Retryable: true})
		return result
	}
	defer repository.Close()
	configData, configErr := readMembershipEvidence(ctx, repository, ConfigPath)
	lockData, lockErr := readMembershipEvidence(ctx, repository, LockPath)
	markerData, markerErr := readMembershipEvidence(ctx, repository, KitMarkerPath)
	sourceMarkerData, sourceMarkerErr := readMembershipEvidence(ctx, repository, KitSourceMarkerPath)
	if configErr != nil || lockErr != nil || markerErr != nil || sourceMarkerErr != nil {
		result.Errors = append(result.Errors, ScanError{
			Code:          "membership_evidence_unavailable",
			Message:       "default-branch Kit membership evidence could not be read completely",
			SourceLocator: source.Locator,
			Retryable:     true,
		})
		return result
	}
	decision := EvaluateMembership(MembershipInput{ConfigData: configData, LockData: lockData, KitMarker: markerData, KitSourceMarker: sourceMarkerData, HomeID: home.CoordinationHomeID, Self: source.Self && home.Role == RoleCoordinationHome})
	if !decision.Member {
		if decision.Incomplete || source.Self {
			result.Errors = append(result.Errors, ScanError{
				Code:          "membership_validation_failed",
				Message:       "default-branch Kit membership evidence is invalid for a required repository",
				SourceLocator: source.Locator,
				RepositoryID:  decision.RepositoryID,
			})
			return result
		}
		result.Candidate = &Candidate{SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Reason: decision.Reason}
		return result
	}
	revision, err := repository.Revision(ctx, repository.DefaultRef)
	if err != nil {
		result.Errors = append(result.Errors, ScanError{Code: "default_revision_unavailable", Message: "default branch revision could not be resolved", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID})
		return result
	}
	settings, settingProblems := decision.PlanSettings, decision.PlanProblems
	var validatedTarget *RemotePlanTarget
	baselineRef := repository.DefaultRef
	targeted := target != nil && target.RepositoryID == decision.RepositoryID
	if targeted {
		if target.EvidenceRevision == "" {
			result.Errors = append(result.Errors, ScanError{Code: "remote_evidence_revision_missing", Message: "prospective remote evidence requires an exact evidence revision", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID})
			return result
		}
		if err := validateRemotePlanTarget(ctx, repository, revision, *target); err != nil {
			result.Errors = append(result.Errors, ScanError{Code: "remote_evidence_checkpoint_invalid", Message: err.Error(), SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: repository.DefaultRef})
			return result
		}
		validatedTarget = target
		baselineRef = target.EvidenceRevision
		settings.DiscoveryVersion = target.DiscoveryVersion
		settings.Mode = target.TargetCatalogMode
		settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
		if settings.DiscoveryVersion != plancatalog.DiscoveryV2 || settings.PolicyDigest != target.PolicyDigest {
			result.Errors = append(result.Errors, ScanError{Code: "remote_evidence_policy_mismatch", Message: "prospective remote candidate policy differs from the reviewed target", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID})
			return result
		}
	}
	memberComplete := false
	memberSourceRevision := revision
	if targeted {
		memberSourceRevision = target.EvidenceRevision
	}
	member := plancatalog.CatalogMember{
		RepositoryID: decision.RepositoryID, SourceKind: source.Kind, SourceLocator: source.Locator,
		SourceIdentitySHA256: remoteIdentitySHA256(source.CloneURL),
		DefaultRef:           repository.DefaultRef, SourceRevision: memberSourceRevision, SelfMember: source.Self,
		DiscoveryVersion: settings.DiscoveryVersion, PolicyDigest: settings.PolicyDigest, Complete: &memberComplete,
	}
	result.Member = &member
	for _, problem := range settingProblems {
		if problem.Severity == plancatalog.SeverityError {
			result.Errors = append(result.Errors, scanProblem(repository, decision.RepositoryID, repository.DefaultRef, problem))
		}
	}
	if plancatalog.HasErrors(settingProblems) {
		return result
	}
	if !targeted && settings.ConfigSchemaVersion == 2 {
		persistedTarget, targetErr := remotePlanTargetFromBinding(ctx, repository, decision.RepositoryID, revision, settings)
		if targetErr != nil || validateRemotePlanTarget(ctx, repository, revision, persistedTarget) != nil {
			message := "persisted migration evidence does not match the remote E/L/A checkpoints"
			if targetErr != nil {
				message = targetErr.Error()
			}
			result.Errors = append(result.Errors, ScanError{Code: "migration_evidence_checkpoint_invalid", Message: message, SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: repository.DefaultRef})
			return result
		}
		validatedTarget = &persistedTarget
	}
	if settings.DiscoveryVersion != plancatalog.DiscoveryV2 {
		result.Errors = append(result.Errors, ScanError{Code: "member_discovery_incompatible", Message: "required portfolio member has not activated plan-catalog discovery v2", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: repository.DefaultRef})
		return result
	}
	baselineClassification, classifyErr := classifyRemoteTree(ctx, repository, decision.RepositoryID, baselineRef, settings)
	if classifyErr != nil {
		result.Errors = append(result.Errors, ScanError{Code: "plan_tree_unavailable", Message: classifyErr.Error(), SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: repository.DefaultRef})
		return result
	}
	member.CandidateSetDigest = baselineClassification.CandidateSetDigest
	if targeted && target.CandidateSetDigest != "" && baselineClassification.CandidateSetDigest != target.CandidateSetDigest {
		result.Errors = append(result.Errors, ScanError{Code: "remote_evidence_candidate_set_mismatch", Message: "prospective remote candidate set differs from the reviewed target", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: baselineRef})
		return result
	}
	baseline, baselineCandidates, errors := collectObservations(ctx, repository, decision.RepositoryID, baselineRef, baselineClassification, nil, "default", "", now)
	result.Observations = append(result.Observations, baseline...)
	result.PlanCandidates = append(result.PlanCandidates, baselineCandidates...)
	result.Errors = append(result.Errors, errors...)
	memberComplete = baselineClassification.Complete && !plancatalog.HasErrors(baselineClassification.Discovery.Problems)
	compatibility := []plancatalog.CompatibilityObservation{}
	compatibilityErrors := []ScanError{}
	if !targeted {
		compatibility, compatibilityErrors = collectCompatibilityObservations(ctx, repository, decision.RepositoryID, repository.DefaultRef, revision, settings, baselineClassification, now)
	}
	result.CompatibilityObservations = append(result.CompatibilityObservations, compatibility...)
	result.Errors = append(result.Errors, compatibilityErrors...)
	legacyRecords := 0
	for _, record := range baselineClassification.Discovery.Records {
		if record.Metadata == nil {
			legacyRecords++
		}
	}
	if !targeted && (len(compatibilityErrors) > 0 || len(compatibility) != legacyRecords) {
		memberComplete = false
	}
	refs, err := repository.RemoteRefs(ctx)
	if err != nil {
		memberComplete = false
		result.Errors = append(result.Errors, ScanError{Code: "branch_enumeration_failed", Message: "remote branch enumeration failed", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID})
		return result
	}
	pullRequests := map[string]string{}
	if source.Enricher != nil {
		enrichment, _ := source.Enricher.EnrichRefs(ctx, source, refs)
		result.APICalls += enrichment.APICalls
		pullRequests = enrichment.PullRequests
	}
	branchBaselineRef := baselineRef
	if validatedTarget != nil {
		branchBaselineRef = validatedTarget.EvidenceRevision
	}
	for _, ref := range refs {
		if refTip, tipErr := repository.Revision(ctx, ref); tipErr == nil && refTip == revision {
			continue
		}
		currentMergeBase, mergeErr := repository.MergeBaseFrom(ctx, repository.DefaultRef, ref)
		if mergeErr != nil {
			memberComplete = false
			code := "branch_merge_base_unavailable"
			if strings.Contains(mergeErr.Error(), "expected one merge base") {
				code = "branch_merge_base_ambiguous"
			}
			result.Errors = append(result.Errors, ScanError{Code: code, Message: "branch merge-base evidence is incomplete or ambiguous", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: ref})
			continue
		}
		merged, mergedErr := repository.IsMergedInto(ctx, ref, repository.DefaultRef, currentMergeBase)
		if mergedErr != nil {
			memberComplete = false
			result.Errors = append(result.Errors, ScanError{Code: "branch_merge_state_unavailable", Message: "branch merge state could not be resolved", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: ref})
			continue
		}
		if merged {
			continue
		}
		mergeBase := currentMergeBase
		if branchBaselineRef != repository.DefaultRef {
			mergeBase, mergeErr = repository.MergeBaseFrom(ctx, branchBaselineRef, ref)
			if mergeErr != nil {
				memberComplete = false
				code := "branch_merge_base_unavailable"
				if strings.Contains(mergeErr.Error(), "expected one merge base") {
					code = "branch_merge_base_ambiguous"
				}
				result.Errors = append(result.Errors, ScanError{Code: code, Message: "reviewed baseline branch merge-base evidence is incomplete or ambiguous", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: ref})
				continue
			}
		}
		changes, changedErr := repository.ChangedPaths(ctx, mergeBase, ref)
		if changedErr != nil {
			memberComplete = false
			result.Errors = append(result.Errors, ScanError{Code: "branch_comparison_incomplete", Message: "changed plan paths could not be read", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: ref})
			continue
		}
		baseClassification, baseClassifyErr := classifyRemoteTree(ctx, repository, decision.RepositoryID, mergeBase, settings)
		branchClassification, branchClassifyErr := classifyRemoteTree(ctx, repository, decision.RepositoryID, ref, settings)
		if baseClassifyErr != nil || branchClassifyErr != nil {
			result.Errors = append(result.Errors, ScanError{Code: "branch_plan_tree_unavailable", Message: "branch or merge-base plan tree could not be classified", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: ref})
			memberComplete = false
			continue
		}
		selected := selectBranchCandidates(changes, baseClassification, branchClassification)
		branchBaselineClassification := baselineClassification
		if branchBaselineRef != baselineRef {
			branchBaselineClassification, baseClassifyErr = classifyRemoteTree(ctx, repository, decision.RepositoryID, branchBaselineRef, settings)
			if baseClassifyErr != nil {
				memberComplete = false
				result.Errors = append(result.Errors, ScanError{Code: "branch_plan_tree_unavailable", Message: "reviewed branch baseline could not be classified", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: ref})
				continue
			}
		}
		result.RemoteBranchEvidence = append(result.RemoteBranchEvidence, buildRemoteBranchEvidence(repository, decision.RepositoryID, branchBaselineRef, ref, mergeBase, changes, branchBaselineClassification)...)
		observations, planCandidates, branchErrors := collectObservations(ctx, repository, decision.RepositoryID, ref, branchClassification, selected, "unmerged-branch", pullRequests[ref], now)
		result.Observations = append(result.Observations, observations...)
		result.PlanCandidates = append(result.PlanCandidates, planCandidates...)
		result.Errors = append(result.Errors, branchErrors...)
		if len(branchErrors) > 0 {
			memberComplete = false
		}
	}
	member.MixedCoverageComplete = memberComplete && (!targeted && len(compatibility) == legacyRecords)
	member.CanonicalReady = memberComplete && legacyRecords == 0
	if validatedTarget != nil && validatedTarget.EvidenceScope == "remote-aware" {
		currentOverlay := remoteOverlayForRepositoryScan(result, decision.RepositoryID, member.SourceIdentitySHA256, validatedTarget.EvidenceRevision, memberComplete)
		overlayMatches := currentOverlay.Digest == validatedTarget.RemoteOverlayDigest
		if !overlayMatches && validatedTarget.RemoteOverlayDigest != "" {
			overlayMatches, _ = remoteOverlayMatchesReviewedOrResolved(ctx, repository, revision, *validatedTarget, currentOverlay)
		}
		switch {
		case validatedTarget.RemoteSourceIdentity == "" || member.SourceIdentitySHA256 != validatedTarget.RemoteSourceIdentity:
			result.Errors = append(result.Errors, ScanError{Code: "remote_overlay_source_mismatch", Message: "selected remote source identity differs from reviewed evidence", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: repository.DefaultRef})
		case validatedTarget.RemoteOverlayDigest == "" || !overlayMatches:
			result.Errors = append(result.Errors, ScanError{Code: "remote_overlay_binding_mismatch", Message: "current remote refs or branch evidence differ from the reviewed overlay", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: repository.DefaultRef})
		}
		if len(result.Errors) > 0 {
			result.CompatibilityObservations = nil
			memberComplete = false
			member.MixedCoverageComplete = false
			member.CanonicalReady = false
		}
	}
	return result
}

func remoteOverlayMatchesReviewedOrResolved(ctx context.Context, repository *GitRepository, defaultRevision string, target RemotePlanTarget, overlay plancatalog.RemoteOverlayEvidence) (bool, error) {
	if overlay.Status != "complete" || overlay.Remote != target.RepositoryID || overlay.SourceIdentitySHA256 != target.RemoteSourceIdentity {
		return false, nil
	}
	ledgerData, err := repository.ReadFile(ctx, target.ActivationBaseRevision, target.LedgerPath)
	if err != nil {
		return false, err
	}
	ledger, err := plancatalog.LoadMigrationLedger(ledgerData)
	if err != nil || ledger.SchemaVersion != 3 {
		return false, err
	}
	reviews := map[string]struct {
		record plancatalog.MigrationRecord
		review plancatalog.BranchTouchReview
	}{}
	for _, record := range ledger.Records {
		for _, review := range record.BranchTouchReviews {
			if !strings.HasPrefix(review.Identity.LogicalRef, "remote:"+target.RepositoryID+":refs/heads/") {
				continue
			}
			key := remoteReviewIdentityKey(review.Identity)
			if _, duplicate := reviews[key]; duplicate {
				return false, nil
			}
			reviews[key] = struct {
				record plancatalog.MigrationRecord
				review plancatalog.BranchTouchReview
			}{record: record, review: review}
		}
	}
	currentEvidence := map[string]plancatalog.BranchCandidateEvidence{}
	for _, evidence := range overlay.BranchTouchCandidates {
		key := remoteReviewIdentityKey(evidence.Identity)
		if _, duplicate := currentEvidence[key]; duplicate {
			return false, nil
		}
		currentEvidence[key] = evidence
	}
	refs := map[plancatalog.RemoteOverlayRef]bool{}
	for _, ref := range overlay.Refs {
		refs[ref] = true
	}
	defaultRef := plancatalog.RemoteOverlayRef{LogicalRef: remoteLogicalBranchRef(target.RepositoryID, repository.DefaultRef), Tip: target.EvidenceRevision}
	if !refs[defaultRef] {
		return false, nil
	}
	reconstructed := overlay
	for key, reviewed := range reviews {
		if evidence, present := currentEvidence[key]; present {
			if !remoteEvidenceMatchesReview(evidence, reviewed.review) {
				return false, nil
			}
			ref := plancatalog.RemoteOverlayRef{LogicalRef: evidence.Identity.LogicalRef, Tip: evidence.RefTip}
			if !refs[ref] {
				return false, nil
			}
			delete(currentEvidence, key)
			continue
		}
		if reviewed.review.Disposition != plancatalog.BranchDispositionDeferredActiveOwner || !remoteDeferredOwnerResolved(ctx, repository, defaultRevision, reviewed.record, reviewed.review) {
			return false, nil
		}
		evidence := plancatalog.BranchCandidateEvidence{
			Algorithm: plancatalog.BranchEvidenceAlgorithm, Identity: reviewed.review.Identity,
			EvidenceRevision: target.EvidenceRevision, PolicyDigest: target.PolicyDigest, CandidateSetDigest: target.CandidateSetDigest,
			RefTip: reviewed.review.RefTip, MergeBase: reviewed.review.MergeBase, PathTransitions: append([]plancatalog.BranchPathTransition{}, reviewed.review.PathTransitions...),
			ProposedDisposition: plancatalog.BranchDispositionActiveOwner, Proof: reviewed.review.Proof, OwnerTipCandidate: reviewed.review.OwnerTipCandidate,
		}
		evidence.Digest = plancatalog.CanonicalBranchEvidenceDigest(evidence)
		reconstructed.BranchTouchCandidates = append(reconstructed.BranchTouchCandidates, evidence)
		reconstructed.Refs = append(reconstructed.Refs, plancatalog.RemoteOverlayRef{LogicalRef: reviewed.review.Identity.LogicalRef, Tip: reviewed.review.RefTip})
	}
	if len(currentEvidence) != 0 {
		return false, nil
	}
	return plancatalog.CanonicalRemoteOverlayDigest(reconstructed) == target.RemoteOverlayDigest, nil
}

func remoteReviewIdentityKey(identity plancatalog.BranchReviewIdentity) string {
	return identity.EvidenceScope + "\x00" + identity.RepositoryID + "\x00" + identity.LogicalRef + "\x00" + identity.CandidatePath
}

func remoteEvidenceMatchesReview(evidence plancatalog.BranchCandidateEvidence, review plancatalog.BranchTouchReview) bool {
	disposition := review.Disposition
	if disposition == plancatalog.BranchDispositionDeferredActiveOwner {
		disposition = plancatalog.BranchDispositionActiveOwner
	}
	return evidence.ProposedDisposition == disposition && len(evidence.Blockers) == 0 && evidence.Digest == plancatalog.CanonicalBranchEvidenceDigest(evidence) &&
		reflect.DeepEqual(evidence.Identity, review.Identity) && evidence.RefTip == review.RefTip && evidence.MergeBase == review.MergeBase &&
		reflect.DeepEqual(evidence.PathTransitions, review.PathTransitions) && reflect.DeepEqual(evidence.Proof, review.Proof) &&
		reflect.DeepEqual(evidence.OwnerTipCandidate, review.OwnerTipCandidate)
}

func remoteDeferredOwnerResolved(ctx context.Context, repository *GitRepository, revision string, record plancatalog.MigrationRecord, review plancatalog.BranchTouchReview) bool {
	if !record.Deferred || review.Disposition != plancatalog.BranchDispositionDeferredActiveOwner || review.OwnerTipCandidate == nil || review.RefTip == "" {
		return false
	}
	file, present, stateErr := remoteTreePathState(ctx, repository, revision, record.CurrentPath)
	data, readErr := repository.ReadFile(ctx, revision, record.CurrentPath)
	if stateErr != nil || !present || readErr != nil {
		return false
	}
	state := plancatalog.BranchPathState{State: plancatalog.BranchPathPresent, Mode: plancatalog.GitMode(file.Mode), ObjectID: file.ObjectID, SHA256: compatibilitySHA256(data)}
	ownerTip, matched := plancatalog.IntegratedDeferredOwnerTip(record, state, data)
	if !matched || ownerTip != review.RefTip {
		return false
	}
	ancestry, ancestryErr := repository.git(ctx, "merge-base", "--is-ancestor", ownerTip, revision)
	return ancestryErr == nil && len(ancestry.Stderr) == 0
}

func remoteOverlayForRepositoryScan(result repositoryScan, repositoryID, sourceIdentity, evidenceRevision string, complete bool) plancatalog.RemoteOverlayEvidence {
	status := "incomplete"
	if complete {
		status = "complete"
	}
	overlay := plancatalog.RemoteOverlayEvidence{
		Remote: repositoryID, SourceIdentitySHA256: sourceIdentity, Status: status,
		Refs: []plancatalog.RemoteOverlayRef{}, BranchTouchCandidates: append([]plancatalog.BranchCandidateEvidence{}, result.RemoteBranchEvidence...),
	}
	if result.Member != nil {
		overlay.Refs = append(overlay.Refs, plancatalog.RemoteOverlayRef{LogicalRef: remoteLogicalBranchRef(repositoryID, result.Member.DefaultRef), Tip: evidenceRevision})
	}
	for _, observation := range result.Observations {
		if observation.RepositoryID == repositoryID && observation.Visibility == "unmerged-branch" && strings.HasPrefix(observation.Ref, "refs/") {
			overlay.Refs = append(overlay.Refs, plancatalog.RemoteOverlayRef{LogicalRef: remoteLogicalBranchRef(repositoryID, observation.Ref), Tip: observation.Commit})
		}
	}
	for _, observation := range result.PlanCandidates {
		if observation.RepositoryID == repositoryID && observation.Visibility == "unmerged-branch" && strings.HasPrefix(observation.Ref, "refs/") {
			overlay.Refs = append(overlay.Refs, plancatalog.RemoteOverlayRef{LogicalRef: remoteLogicalBranchRef(repositoryID, observation.Ref), Tip: observation.Commit})
		}
	}
	overlay.Digest = plancatalog.CanonicalRemoteOverlayDigest(overlay)
	return overlay
}

func remotePlanTargetFromBinding(ctx context.Context, repository *GitRepository, repositoryID, revision string, settings plancatalog.RepositorySettings) (RemotePlanTarget, error) {
	binding := settings.MigrationEvidence
	if binding == nil {
		return RemotePlanTarget{}, fmt.Errorf("migration_evidence_binding_missing: config schema v2 requires an exact binding")
	}
	history, err := repositoryGitText(ctx, repository, "rev-list", "--first-parent", "--reverse", binding.ActivationBaseRevision+".."+revision)
	if err != nil || history == "" {
		return RemotePlanTarget{}, fmt.Errorf("activation_checkpoint_missing: bound activation checkpoint is absent")
	}
	activation := strings.Fields(history)[0]
	ledgerData, err := repository.ReadFile(ctx, binding.ActivationBaseRevision, binding.LedgerPath)
	if err != nil {
		return RemotePlanTarget{}, fmt.Errorf("migration_evidence_ledger_unavailable: bound ledger is unavailable at L")
	}
	ledger, err := plancatalog.LoadMigrationLedger(ledgerData)
	ledgerDigest := sha256.Sum256(ledgerData)
	if err != nil || ledger.SchemaVersion != 3 || ledger.RepositoryID != repositoryID || hex.EncodeToString(ledgerDigest[:]) != binding.LedgerSHA256 ||
		ledger.EvidenceRevision != binding.EvidenceRevision || ledger.TargetCatalogMode != settings.Mode || ledger.BranchEvidence == nil ||
		ledger.BranchEvidence.Digest != binding.BranchEvidenceDigest || ledger.BranchEvidence.EvidenceScope != binding.EvidenceScope {
		return RemotePlanTarget{}, fmt.Errorf("migration_evidence_ledger_invalid: bound ledger is not the repository schema-v3 authority")
	}
	if binding.EvidenceScope == "remote-aware" && (ledger.BranchEvidence.RemoteOverlayStatus != "complete" || ledger.BranchEvidence.RemoteOverlayDigest == "" || ledger.BranchEvidence.RemoteSourceIdentitySHA256 == "") {
		return RemotePlanTarget{}, fmt.Errorf("migration_evidence_ledger_invalid: remote-aware ledger omits complete overlay authority")
	}
	return RemotePlanTarget{
		RepositoryID: repositoryID, EvidenceRevision: binding.EvidenceRevision,
		ActivationBaseRevision: binding.ActivationBaseRevision, ActivationRevision: activation,
		LedgerPath: binding.LedgerPath, LedgerSHA256: binding.LedgerSHA256,
		DiscoveryVersion: plancatalog.DiscoveryV2, TargetCatalogMode: settings.Mode,
		PolicyDigest: ledger.PolicyDigest, CandidateSetDigest: ledger.CandidateSetDigest,
		MigrationActionDigest: binding.MigrationActionDigest, BranchEvidenceDigest: binding.BranchEvidenceDigest,
		EvidenceScope: binding.EvidenceScope, RemoteOverlayDigest: ledger.BranchEvidence.RemoteOverlayDigest,
		RemoteSourceIdentity: ledger.BranchEvidence.RemoteSourceIdentitySHA256,
	}, nil
}

func validateRemotePlanTarget(ctx context.Context, repository *GitRepository, defaultRevision string, target RemotePlanTarget) error {
	evidenceRevision, err := repository.Revision(ctx, target.EvidenceRevision)
	if err != nil || evidenceRevision != target.EvidenceRevision {
		return fmt.Errorf("reviewed evidence revision is not available in the fresh mirror")
	}
	if defaultRevision == target.EvidenceRevision {
		return nil
	}
	if target.ActivationBaseRevision == "" || target.LedgerPath == "" || target.LedgerSHA256 == "" {
		return fmt.Errorf("remote target is missing the exact ledger checkpoint binding")
	}
	parent, err := repository.Revision(ctx, target.ActivationBaseRevision+"^1")
	if err != nil || parent != target.EvidenceRevision {
		return fmt.Errorf("remote ledger checkpoint does not have the evidence revision as first parent")
	}
	changes, err := repository.ChangedPaths(ctx, target.EvidenceRevision, target.ActivationBaseRevision)
	if err != nil || len(changes) != 1 || changes[0].NewPath != target.LedgerPath || (changes[0].OldPath != "" && changes[0].OldPath != target.LedgerPath) {
		return fmt.Errorf("remote ledger checkpoint tree delta is not ledger-only")
	}
	ledgerData, err := repository.ReadFile(ctx, target.ActivationBaseRevision, target.LedgerPath)
	if err != nil {
		return fmt.Errorf("remote ledger checkpoint blob is unavailable")
	}
	digest := sha256.Sum256(ledgerData)
	if hex.EncodeToString(digest[:]) != target.LedgerSHA256 {
		return fmt.Errorf("remote ledger checkpoint blob differs from the reviewed ledger")
	}
	if defaultRevision == target.ActivationBaseRevision {
		return nil
	}
	if target.ActivationRevision == "" || target.MigrationActionDigest == "" {
		return fmt.Errorf("remote default is past the ledger checkpoint without an exact activation binding")
	}
	activation, activationErr := repository.Revision(ctx, target.ActivationRevision)
	activationParent, parentErr := repository.Revision(ctx, target.ActivationRevision+"^1")
	if activationErr != nil || parentErr != nil || activation != target.ActivationRevision || activationParent != target.ActivationBaseRevision {
		return fmt.Errorf("remote activation checkpoint is not the immediate child of the ledger checkpoint")
	}
	merged, mergeErr := repository.git(ctx, "merge-base", "--is-ancestor", target.ActivationRevision, defaultRevision)
	if mergeErr != nil || len(merged.Stderr) != 0 {
		return fmt.Errorf("remote default does not descend from the exact activation checkpoint")
	}
	actionDigest, actionPaths, actionErr := plancatalog.MigrationActionDigestForCheckpoint(repository.Path, target.ActivationBaseRevision, target.ActivationRevision)
	if actionErr != nil || actionDigest != target.MigrationActionDigest {
		return fmt.Errorf("remote activation action digest differs from the reviewed projection")
	}
	configData, configErr := repository.ReadFile(ctx, target.ActivationRevision, state.ConfigPath)
	if configErr != nil {
		return fmt.Errorf("remote activation config is unavailable")
	}
	settings, settingProblems := plancatalog.DecodeRepositorySettings(configData)
	if plancatalog.HasErrors(settingProblems) {
		return fmt.Errorf("remote activation config is invalid: %s", settingProblems[0].Message)
	}
	if settings.ConfigSchemaVersion != 2 {
		return fmt.Errorf("remote activation config is not schema v2")
	}
	if settings.DiscoveryVersion != plancatalog.DiscoveryV2 || settings.Mode != target.TargetCatalogMode {
		return fmt.Errorf("remote activation catalog mode or discovery version differs from reviewed evidence")
	}
	if settings.MigrationEvidence == nil {
		return fmt.Errorf("remote activation config omits the reviewed migration evidence binding")
	}
	binding := settings.MigrationEvidence
	if binding.LedgerPath != target.LedgerPath || binding.LedgerSHA256 != target.LedgerSHA256 || binding.BranchEvidenceDigest != target.BranchEvidenceDigest || binding.EvidenceRevision != target.EvidenceRevision || binding.ActivationBaseRevision != target.ActivationBaseRevision || binding.MigrationActionDigest != target.MigrationActionDigest || binding.EvidenceScope != target.EvidenceScope {
		return fmt.Errorf("remote activation config binding differs from reviewed evidence")
	}
	ledger, ledgerErr := plancatalog.LoadMigrationLedger(ledgerData)
	if ledgerErr != nil || ledger.SchemaVersion != 3 || ledger.RepositoryID != target.RepositoryID || ledger.EvidenceRevision != target.EvidenceRevision ||
		ledger.TargetCatalogMode != target.TargetCatalogMode || ledger.PolicyDigest != target.PolicyDigest || ledger.CandidateSetDigest != target.CandidateSetDigest ||
		ledger.BranchEvidence == nil || ledger.BranchEvidence.Digest != target.BranchEvidenceDigest || ledger.BranchEvidence.EvidenceScope != target.EvidenceScope ||
		ledger.BranchEvidence.RemoteOverlayDigest != target.RemoteOverlayDigest || ledger.BranchEvidence.RemoteSourceIdentitySHA256 != target.RemoteSourceIdentity {
		return fmt.Errorf("remote ledger checkpoint is not a valid schema-v3 ledger")
	}
	resolvedDeferred := map[string]bool{}
	for _, record := range ledger.Records {
		if !record.Deferred {
			continue
		}
		file, present, stateErr := remoteTreePathState(ctx, repository, defaultRevision, record.CurrentPath)
		data, readErr := repository.ReadFile(ctx, defaultRevision, record.CurrentPath)
		if stateErr != nil || !present || readErr != nil {
			continue
		}
		state := plancatalog.BranchPathState{State: plancatalog.BranchPathPresent, Mode: plancatalog.GitMode(file.Mode), ObjectID: file.ObjectID, SHA256: compatibilitySHA256(data)}
		ownerTip, integrated := plancatalog.IntegratedDeferredOwnerTip(record, state, data)
		if !integrated {
			continue
		}
		ancestry, ancestryErr := repository.git(ctx, "merge-base", "--is-ancestor", ownerTip, defaultRevision)
		if ancestryErr == nil && len(ancestry.Stderr) == 0 {
			resolvedDeferred[record.CurrentPath] = true
		}
	}
	activationClassification, activationClassifyErr := classifyRemoteTree(ctx, repository, target.RepositoryID, target.ActivationRevision, settings)
	defaultClassification, defaultClassifyErr := classifyRemoteTree(ctx, repository, target.RepositoryID, defaultRevision, settings)
	if activationClassifyErr != nil || defaultClassifyErr != nil || !activationClassification.Complete || !defaultClassification.Complete {
		return fmt.Errorf("remote descendant candidate authority is unavailable")
	}
	if activationClassification.PolicyDigest != target.PolicyDigest || defaultClassification.PolicyDigest != target.PolicyDigest || plancatalog.CandidateSetDigestForCandidates(candidatesExcludingPaths(activationClassification.Candidates, resolvedDeferred)) != plancatalog.CandidateSetDigestForCandidates(candidatesExcludingPaths(defaultClassification.Candidates, resolvedDeferred)) {
		return fmt.Errorf("remote descendant changed the activated candidate set")
	}
	safetyPaths := append([]string{state.ConfigPath, target.LedgerPath, plancatalog.LegacyRegisterPath}, actionPaths...)
	for _, record := range ledger.Records {
		safetyPaths = append(safetyPaths, record.CurrentPath, record.TargetPath)
	}
	sort.Strings(safetyPaths)
	for index, candidatePath := range safetyPaths {
		if candidatePath == "" || (index > 0 && candidatePath == safetyPaths[index-1]) {
			continue
		}
		if resolvedDeferred[candidatePath] {
			continue
		}
		if equal, compareErr := remotePathStateEqual(ctx, repository, target.ActivationRevision, defaultRevision, candidatePath); compareErr != nil || !equal {
			return fmt.Errorf("remote descendant changed bound state at %s", candidatePath)
		}
	}
	return nil
}

func candidatesExcludingPaths(candidates []plancatalog.Candidate, excluded map[string]bool) []plancatalog.Candidate {
	result := make([]plancatalog.Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		if !excluded[candidate.Path] {
			result = append(result, candidate)
		}
	}
	return result
}

func remotePathStateEqual(ctx context.Context, repository *GitRepository, leftRevision, rightRevision, candidatePath string) (bool, error) {
	left, leftPresent, leftErr := remoteTreePathState(ctx, repository, leftRevision, candidatePath)
	right, rightPresent, rightErr := remoteTreePathState(ctx, repository, rightRevision, candidatePath)
	if leftErr != nil || rightErr != nil {
		return false, fmt.Errorf("bound path state is unavailable")
	}
	if !leftPresent || !rightPresent {
		return leftPresent == rightPresent, nil
	}
	return left.Mode == right.Mode && left.ObjectID == right.ObjectID, nil
}

func remoteTreePathState(ctx context.Context, repository *GitRepository, revision, candidatePath string) (GitTreeFile, bool, error) {
	files, err := repository.ListFiles(ctx, revision)
	if err != nil {
		return GitTreeFile{}, false, err
	}
	index := sort.Search(len(files), func(index int) bool { return files[index].Path >= candidatePath })
	if index >= len(files) || files[index].Path != candidatePath {
		return GitTreeFile{}, false, nil
	}
	return files[index], true, nil
}

func buildRemoteBranchEvidence(repository *GitRepository, repositoryID, evidenceRevision, ref, mergeBase string, changes []GitTreeChange, baseline plancatalog.Classification) []plancatalog.BranchCandidateEvidence {
	result := []plancatalog.BranchCandidateEvidence{}
	for _, candidate := range baseline.Candidates {
		if candidate.Ownership != plancatalog.OwnershipOwned {
			continue
		}
		paths := []string{}
		touched := false
		for _, change := range changes {
			if change.OldPath == candidate.Path || change.NewPath == candidate.Path {
				touched = true
				paths = append(paths, change.OldPath, change.NewPath)
			}
		}
		if !touched {
			continue
		}
		paths = compactRemoteEvidencePaths(candidate.Path, paths)
		logical := remoteLogicalBranchRef(repositoryID, ref)
		evidence := plancatalog.ProposeBranchCandidateEvidence(plancatalog.BranchEvidenceRequest{
			Root: repository.Path, EvidenceScope: "remote-aware", RepositoryID: repositoryID,
			LogicalRef: logical, Ref: ref, EvidenceRevision: evidenceRevision,
			PolicyDigest: baseline.PolicyDigest, CandidateSetDigest: baseline.CandidateSetDigest,
			CandidatePath: candidate.Path, Paths: paths, ExpectedKind: candidate.ExpectedKind,
		})
		if evidence.MergeBase != "" && evidence.MergeBase != mergeBase {
			evidence.ProposedDisposition = plancatalog.BranchDispositionBlocking
			evidence.Proof = nil
			evidence.OwnerTipCandidate = nil
			evidence.Blockers = append(evidence.Blockers, plancatalog.BranchEvidenceBlocker{Code: "branch_merge_base_moved", Message: "remote scanner merge base differs from branch proof merge base", Path: candidate.Path, Ref: logical})
			evidence.Digest = plancatalog.CanonicalBranchEvidenceDigest(evidence)
		}
		result = append(result, evidence)
	}
	sort.SliceStable(result, func(i, j int) bool {
		left := result[i].Identity.LogicalRef + "\x00" + result[i].Identity.CandidatePath
		right := result[j].Identity.LogicalRef + "\x00" + result[j].Identity.CandidatePath
		return left < right
	})
	return result
}

func compactRemoteEvidencePaths(candidatePath string, paths []string) []string {
	seen := map[string]bool{candidatePath: true}
	result := []string{candidatePath}
	for _, candidate := range paths {
		if candidate == "" || seen[candidate] {
			continue
		}
		seen[candidate] = true
		result = append(result, candidate)
	}
	sort.Strings(result)
	return result
}

func remoteLogicalBranchRef(repositoryID, ref string) string {
	branch := strings.TrimPrefix(ref, "refs/remotes/origin/")
	if strings.HasPrefix(ref, "refs/heads/") {
		branch = strings.TrimPrefix(ref, "refs/heads/")
	}
	return "remote:" + repositoryID + ":refs/heads/" + branch
}

func collectCompatibilityObservations(ctx context.Context, repository *GitRepository, repositoryID, ref, revision string, settings plancatalog.RepositorySettings, classification plancatalog.Classification, now time.Time) ([]plancatalog.CompatibilityObservation, []ScanError) {
	legacyPaths := map[string]plancatalog.Record{}
	currentRecords := map[string]plancatalog.Record{}
	for _, record := range classification.Discovery.Records {
		currentRecords[record.Path] = record
		if record.Metadata == nil {
			legacyPaths[record.Path] = record
		}
	}
	binding := settings.MigrationEvidence
	if binding == nil {
		if len(legacyPaths) == 0 {
			return []plancatalog.CompatibilityObservation{}, []ScanError{}
		}
		return nil, []ScanError{{Code: "migration_evidence_binding_missing", Message: "metadata-less mixed records require a persisted reviewed migration binding", RepositoryID: repositoryID, Ref: ref}}
	}
	ledgerData, err := repository.ReadFile(ctx, ref, binding.LedgerPath)
	if err != nil {
		return nil, []ScanError{{Code: "migration_evidence_ledger_unavailable", Message: "bound migration ledger could not be read", RepositoryID: repositoryID, Ref: ref, Path: binding.LedgerPath}}
	}
	digest := sha256.Sum256(ledgerData)
	if hex.EncodeToString(digest[:]) != binding.LedgerSHA256 {
		return nil, []ScanError{{Code: "migration_evidence_ledger_mismatch", Message: "bound migration ledger bytes changed", RepositoryID: repositoryID, Ref: ref, Path: binding.LedgerPath}}
	}
	ledger, err := plancatalog.LoadMigrationLedger(ledgerData)
	if err != nil || ledger.SchemaVersion != 3 || ledger.EvidenceRevision != binding.EvidenceRevision || ledger.BranchEvidence == nil || ledger.BranchEvidence.Digest != binding.BranchEvidenceDigest {
		return nil, []ScanError{{Code: "migration_evidence_ledger_invalid", Message: "bound migration ledger does not match config authority", RepositoryID: repositoryID, Ref: ref, Path: binding.LedgerPath}}
	}
	if parent, parentErr := repositoryGitText(ctx, repository, "rev-parse", "--verify", binding.ActivationBaseRevision+"^1"); parentErr != nil || parent != binding.EvidenceRevision {
		return nil, []ScanError{{Code: "ledger_checkpoint_parent_mismatch", Message: "remote ledger checkpoint first parent is not the evidence revision", RepositoryID: repositoryID, Ref: ref}}
	}
	children, historyErr := repositoryGitText(ctx, repository, "rev-list", "--first-parent", "--reverse", binding.ActivationBaseRevision+".."+revision)
	if historyErr != nil || children == "" {
		return nil, []ScanError{{Code: "activation_checkpoint_missing", Message: "remote default history does not contain the bound activation checkpoint", RepositoryID: repositoryID, Ref: ref}}
	}
	activation := strings.Fields(children)[0]
	if parent, parentErr := repositoryGitText(ctx, repository, "rev-parse", "--verify", activation+"^1"); parentErr != nil || parent != binding.ActivationBaseRevision {
		return nil, []ScanError{{Code: "activation_checkpoint_parent_mismatch", Message: "remote activation checkpoint is not the immediate first-parent child", RepositoryID: repositoryID, Ref: ref}}
	}
	observations := []plancatalog.CompatibilityObservation{}
	errors := []ScanError{}
	seen := map[string]bool{}
	treeFiles, treeErr := repository.ListFiles(ctx, ref)
	if treeErr != nil {
		return nil, []ScanError{{Code: "migration_evidence_tree_unavailable", Message: "current remote tree state is unavailable", RepositoryID: repositoryID, Ref: ref}}
	}
	treeByPath := map[string]GitTreeFile{}
	for _, file := range treeFiles {
		treeByPath[file.Path] = file
	}
	for _, record := range ledger.Records {
		if !record.Deferred {
			continue
		}
		legacy, exists := legacyPaths[record.CurrentPath]
		if !exists {
			current, currentExists := currentRecords[record.CurrentPath]
			file, fileExists := treeByPath[record.CurrentPath]
			data, readErr := repository.ReadFile(ctx, ref, record.CurrentPath)
			state := plancatalog.BranchPathState{}
			if fileExists && readErr == nil {
				state = plancatalog.BranchPathState{State: plancatalog.BranchPathPresent, Mode: plancatalog.GitMode(file.Mode), ObjectID: file.ObjectID, SHA256: compatibilitySHA256(data)}
			}
			ownerTip, integrated := plancatalog.IntegratedDeferredOwnerTip(record, state, data)
			if currentExists && current.Metadata != nil && integrated {
				ancestry, ancestryErr := repository.git(ctx, "merge-base", "--is-ancestor", ownerTip, revision)
				if ancestryErr == nil && len(ancestry.Stderr) == 0 {
					seen[record.CurrentPath] = true
					continue
				}
			}
			errors = append(errors, ScanError{Code: "deferred_owner_resolution_invalid", Message: "canonical metadata does not prove integration of the exact reviewed owner tip", RepositoryID: repositoryID, Ref: ref, Path: record.CurrentPath})
			continue
		}
		data, readErr := repository.ReadFile(ctx, ref, record.CurrentPath)
		if readErr != nil || compatibilitySHA256(data) != record.SourceSHA256 {
			errors = append(errors, ScanError{Code: "deferred_plan_moved", Message: "deferred plan bytes differ from the reviewed cutover source", RepositoryID: repositoryID, Ref: ref, Path: record.CurrentPath})
			continue
		}
		seen[record.CurrentPath] = true
		decision := record.Decision
		observations = append(observations, plancatalog.CompatibilityObservation{
			RepositoryID: repositoryID, PlanID: decision.ID, Title: plancatalog.DisplayTitle(legacy, nil), Kind: decision.Kind,
			Purpose: decision.Purpose, Lifecycle: legacy.Header.Lifecycle, Family: decision.Family,
			Products: append([]string{}, decision.Products...), Capabilities: append([]string{}, decision.Capabilities...),
			StrategicThemes: append([]string{}, decision.StrategicThemes...), Relations: append([]plancatalog.Relation{}, decision.Relations...),
			LegacyAliases: append([]string{}, record.LegacyAliases...), CanonicalPath: record.CurrentPath,
			Ref: ref, Commit: revision, ContentSHA256: record.SourceSHA256, CoverageDisposition: "mixed-grandfathered",
			IncrementalMigrationRequired: true, MigrationEvidence: *binding, Verification: "verified",
			ObservedAt: now.UTC().Truncate(time.Second).Format(time.RFC3339), Stale: false, Conflict: false,
		})
	}
	for legacyPath := range legacyPaths {
		if !seen[legacyPath] {
			errors = append(errors, ScanError{Code: "mixed_coverage_unreviewed", Message: "metadata-less mixed plan is absent from the bound deferred ledger", RepositoryID: repositoryID, Ref: ref, Path: legacyPath})
		}
	}
	return observations, errors
}

func repositoryGitText(ctx context.Context, repository *GitRepository, args ...string) (string, error) {
	result, err := repository.runner.Run(ctx, "git", safeGitArgs(append([]string{"-C", repository.Path}, args...)...)...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result.Stdout)), nil
}

func compatibilitySHA256(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func readMembershipEvidence(ctx context.Context, repository *GitRepository, path string) ([]byte, error) {
	data, err := repository.ReadFile(ctx, repository.DefaultRef, path)
	if stderrors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return data, err
}

func classifyRemoteTree(ctx context.Context, repository *GitRepository, repositoryID, ref string, settings plancatalog.RepositorySettings) (plancatalog.Classification, error) {
	files, err := repository.ListFiles(ctx, ref)
	if err != nil {
		return plancatalog.Classification{}, err
	}
	blobs := make([]plancatalog.GitBlob, 0, len(files))
	byPath := make(map[string]GitTreeFile, len(files))
	for _, file := range files {
		blobs = append(blobs, plancatalog.GitBlob{Path: file.Path, Mode: plancatalog.GitMode(file.Mode), ObjectID: file.ObjectID, Revision: ref})
		byPath[file.Path] = file
	}
	batch, err := newRemoteBatchReader(ctx, repository)
	if err != nil {
		return plancatalog.Classification{}, err
	}
	reader := func(blob plancatalog.GitBlob) ([]byte, error) {
		file, found := byPath[blob.Path]
		if !found || file.ObjectID != blob.ObjectID || file.Mode != string(blob.Mode) {
			return nil, fmt.Errorf("remote_blob_evidence_changed: %s", blob.Path)
		}
		if batch != nil {
			return batch.Read(blob)
		}
		return repository.ReadBlob(ctx, file)
	}
	classification, classifyErr := plancatalog.ClassifyGitBlobs(blobs, reader, settings.DiscoveryPolicy(false), settings.Mode, repositoryID, nil)
	if batch == nil {
		return classification, classifyErr
	}
	closeErr := batch.Close()
	if classifyErr != nil {
		return plancatalog.Classification{}, classifyErr
	}
	if closeErr != nil {
		return plancatalog.Classification{}, closeErr
	}
	return classification, nil
}

func collectObservations(ctx context.Context, repository *GitRepository, repositoryID, ref string, classification plancatalog.Classification, selected map[string]bool, visibility, pullRequest string, now time.Time) ([]plancatalog.SourceObservation, []PlanCandidateObservation, []ScanError) {
	errors := classificationScanErrors(repository, repositoryID, ref, classification, selected)
	invalidPaths := map[string]bool{}
	for _, problem := range classification.Discovery.Problems {
		if problem.Severity == plancatalog.SeverityError {
			invalidPaths[problem.Path] = true
		}
	}
	records := []plancatalog.Record{}
	for _, record := range classification.Discovery.Records {
		if record.Metadata == nil || invalidPaths[record.Path] || (selected != nil && !selected[record.Path]) {
			continue
		}
		records = append(records, record)
	}
	legacyEntries := []plancatalog.LegacyEntry{}
	if data, readErr := repository.ReadFile(ctx, ref, plancatalog.LegacyRegisterPath); readErr == nil {
		legacyEntries, _ = plancatalog.ParseLegacyRegister(data)
	} else if !stderrors.Is(readErr, os.ErrNotExist) {
		errors = append(errors, ScanError{Code: "legacy_title_evidence_unavailable", Message: "legacy title evidence could not be read", SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref})
	}
	legacyReconciliation := plancatalog.ReconcileLegacy(records, nil, legacyEntries)
	commit, err := repository.Revision(ctx, ref)
	if err != nil {
		return nil, nil, append(errors, ScanError{Code: "observation_revision_unavailable", Message: "observation commit could not be resolved", SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref})
	}
	planCandidates := []PlanCandidateObservation{}
	for _, candidate := range classification.Candidates {
		if selected != nil && !selected[candidate.Path] {
			continue
		}
		verification := "verified"
		if (candidate.Provenance.Source.Mode == plancatalog.GitModeRegular || candidate.Provenance.Source.Mode == plancatalog.GitModeExecutable) && candidate.Provenance.Source.ContentSHA256 == "" {
			verification = "unverified"
		}
		planCandidates = append(planCandidates, PlanCandidateObservation{
			RepositoryID: repositoryID, Path: candidate.Path, Ref: ref, Commit: commit, Visibility: visibility,
			Signal: candidate.Signal, Ownership: candidate.Ownership, ExpectedKind: candidate.ExpectedKind,
			GitMode: candidate.Provenance.Source.Mode, ObjectID: candidate.Provenance.Source.ObjectID,
			ContentSHA256: candidate.Provenance.Source.ContentSHA256, PolicyDigest: candidate.Provenance.PolicyDigest,
			ExclusionRoot: candidate.Provenance.ExclusionRoot, Boundary: candidate.Provenance.Boundary,
			AmbiguousUnder: candidate.Provenance.AmbiguousUnder, PullRequest: pullRequest, Verification: verification,
		})
	}
	changedAt := now
	if value, timeErr := repository.CommitTime(ctx, ref); timeErr == nil {
		if parsed, parseErr := time.Parse(time.RFC3339, value); parseErr == nil {
			changedAt = parsed.UTC()
		}
	}
	observations := []plancatalog.SourceObservation{}
	for _, record := range records {
		metadata := record.Metadata
		recordChangedAt := changedAt
		if value, historyErr := repository.PathCommitTime(ctx, ref, record.Path); historyErr == nil {
			if parsed, parseErr := time.Parse(time.RFC3339, value); parseErr == nil {
				recordChangedAt = parsed.UTC()
			} else {
				errors = append(errors, ScanError{Code: "observation_freshness_unavailable", Message: "plan-specific source-change time is invalid", SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref})
			}
		} else {
			errors = append(errors, ScanError{Code: "observation_freshness_unavailable", Message: "plan-specific source-change time could not be resolved", SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref})
		}
		observations = append(observations, plancatalog.SourceObservation{
			RepositoryID: repositoryID, PlanID: metadata.ID, Title: plancatalog.DisplayTitle(record, legacyReconciliation.ByCanonicalPath[record.Path]), Kind: metadata.Kind,
			Purpose: metadata.Purpose, Lifecycle: record.Header.Lifecycle, Family: metadata.Family,
			Products: append([]string{}, metadata.Products...), Capabilities: append([]string{}, metadata.Capabilities...),
			StrategicThemes: append([]string{}, metadata.StrategicThemes...), Relations: append([]plancatalog.Relation{}, metadata.Relations...),
			LegacyAliases: append([]string{}, metadata.LegacyAliases...), CanonicalPath: record.Path, Ref: ref, Commit: commit,
			PullRequest: pullRequest, ContentSHA256: record.ContentSHA256, Visibility: visibility, Verification: "verified",
			ObservedAt: now.UTC().Truncate(time.Second).Format(time.RFC3339), SourceChangedAt: recordChangedAt.UTC().Truncate(time.Second).Format(time.RFC3339),
			Stale: visibility == "unmerged-branch" && now.Sub(recordChangedAt) > 30*24*time.Hour,
		})
	}
	return observations, planCandidates, errors
}

func classificationScanErrors(repository *GitRepository, repositoryID, ref string, classification plancatalog.Classification, selected map[string]bool) []ScanError {
	errors := []ScanError{}
	for _, problem := range classification.Discovery.Problems {
		if problem.Severity != plancatalog.SeverityError || (selected != nil && problem.Path != "" && !selected[problem.Path]) {
			continue
		}
		errors = append(errors, scanProblem(repository, repositoryID, ref, problem))
	}
	return errors
}

func scanProblem(repository *GitRepository, repositoryID, ref string, problem plancatalog.Problem) ScanError {
	return ScanError{Code: problem.Code, Message: problem.Message, SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref, Path: problem.Path}
}

func selectBranchCandidates(changes []GitTreeChange, base, branch plancatalog.Classification) map[string]bool {
	baseCandidates := map[string]plancatalog.Candidate{}
	branchCandidates := map[string]plancatalog.Candidate{}
	for _, candidate := range base.Candidates {
		baseCandidates[candidate.Path] = candidate
	}
	for _, candidate := range branch.Candidates {
		branchCandidates[candidate.Path] = candidate
	}
	branchProblemPaths := map[string]bool{}
	for _, problem := range branch.Discovery.Problems {
		if problem.Severity == plancatalog.SeverityError && problem.Path != "" {
			branchProblemPaths[problem.Path] = true
		}
	}
	selected := map[string]bool{}
	for _, change := range changes {
		candidate, found := branchCandidates[change.NewPath]
		if !found {
			if branchProblemPaths[change.NewPath] {
				selected[change.NewPath] = true
			}
			continue
		}
		if strings.HasPrefix(change.Status, "R") {
			prior, priorFound := baseCandidates[change.OldPath]
			if change.OldObjectID == change.NewObjectID && priorFound && sameCandidateSemantics(prior, candidate) {
				continue
			}
		} else if change.Status == "M" && change.OldObjectID == change.NewObjectID {
			prior, priorFound := baseCandidates[change.OldPath]
			if priorFound && sameCandidateSemantics(prior, candidate) {
				continue
			}
		}
		selected[change.NewPath] = true
	}
	return selected
}

func sameCandidateSemantics(left, right plancatalog.Candidate) bool {
	return left.ExpectedKind == right.ExpectedKind && left.FamilyQualified == right.FamilyQualified && left.Signal == right.Signal && left.Ownership == right.Ownership
}

func aggregateMemberDigest(members []plancatalog.CatalogMember, selectDigest func(plancatalog.CatalogMember) string) string {
	type evidence struct {
		RepositoryID     string                       `json:"repository_id"`
		DiscoveryVersion plancatalog.DiscoveryVersion `json:"discovery_version"`
		Digest           string                       `json:"digest"`
	}
	values := make([]evidence, 0, len(members))
	for _, member := range members {
		values = append(values, evidence{RepositoryID: member.RepositoryID, DiscoveryVersion: member.DiscoveryVersion, Digest: selectDigest(member)})
	}
	sort.SliceStable(values, func(i, j int) bool { return values[i].RepositoryID < values[j].RepositoryID })
	data, _ := json.Marshal(values)
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func markObservationConflicts(observations []plancatalog.SourceObservation) {
	counts := map[string]int{}
	for _, observation := range observations {
		counts[observation.RepositoryID+"\x00"+observation.PlanID]++
	}
	for index := range observations {
		key := observations[index].RepositoryID + "\x00" + observations[index].PlanID
		observations[index].Conflict = counts[key] > 1
	}
}

func ReadCachedCatalog(root string) (Catalog, error) {
	var data []byte
	var err error
	for attempt := 0; attempt < fileShareRetryAttempts; attempt++ {
		data, err = readRootRegular(root, CatalogPath)
		if err == nil || (!stderrors.Is(err, errRootRegularChanged) && !isTransientFileSharingError(err)) {
			break
		}
		time.Sleep(fileShareRetryDelay)
	}
	if err != nil {
		return Catalog{}, err
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return Catalog{}, err
	}
	schemaVersion := state.AsInt(raw["schema_version"])
	schemaPath, err := state.SchemaForPlanCatalogVersion(schemaVersion)
	if err != nil {
		return Catalog{}, err
	}
	if err := state.Validate(schemaPath, raw); err != nil {
		return Catalog{}, err
	}
	var catalog Catalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return Catalog{}, err
	}
	return catalog, nil
}
