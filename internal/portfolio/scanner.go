package portfolio

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"os"
	"path/filepath"
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
}

type repositoryScan struct {
	Source       RepositorySource
	Member       *plancatalog.CatalogMember
	Observations []plancatalog.SourceObservation
	Candidate    *Candidate
	Errors       []ScanError
	APICalls     int
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
		SchemaVersion:      1,
		CoordinationHomeID: options.Config.CoordinationHomeID,
		StartedAt:          started.Format(time.RFC3339),
		Complete:           true,
		Members:            []plancatalog.CatalogMember{},
		Observations:       []plancatalog.SourceObservation{},
		Candidates:         []Candidate{},
		Errors:             []ScanError{},
		Metrics:            Metrics{MaxConcurrency: options.MaxConcurrency},
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
			results[index] = scanRepository(ctx, manager, source, options.Config, started)
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
		catalog.Members = append(catalog.Members, *result.Member)
		catalog.Observations = append(catalog.Observations, result.Observations...)
	}
	markObservationConflicts(catalog.Observations)
	completed := time.Now().UTC().Truncate(time.Second)
	if !options.Now.IsZero() {
		completed = started
	}
	catalog.CompletedAt = completed.Format(time.RFC3339)
	if catalog.Complete {
		catalog.LastCompleteScanAt = catalog.CompletedAt
	} else if previous, readErr := ReadCachedCatalog(options.Config.Root); readErr == nil && previous.Complete {
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
	for _, observation := range catalog.Observations {
		if observation.Stale {
			catalog.Metrics.StaleCount++
		}
	}
	sortCatalog(&catalog)
	if err := state.Validate(state.PlanCatalogSchema, catalog); err != nil {
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

func scanRepository(ctx context.Context, manager MirrorManager, source RepositorySource, home Config, now time.Time) repositoryScan {
	result := repositoryScan{Source: source, Observations: []plancatalog.SourceObservation{}, Errors: []ScanError{}}
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
	member := plancatalog.CatalogMember{RepositoryID: decision.RepositoryID, SourceKind: source.Kind, SourceLocator: source.Locator, DefaultRef: repository.DefaultRef, SourceRevision: revision, SelfMember: source.Self}
	result.Member = &member
	baseline, errors := collectObservations(ctx, repository, decision.RepositoryID, repository.DefaultRef, nil, "default", "", now)
	result.Observations = append(result.Observations, baseline...)
	result.Errors = append(result.Errors, errors...)
	refs, err := repository.RemoteRefs(ctx)
	if err != nil {
		result.Errors = append(result.Errors, ScanError{Code: "branch_enumeration_failed", Message: "remote branch enumeration failed", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID})
		return result
	}
	pullRequests := map[string]string{}
	if source.Enricher != nil {
		enrichment, _ := source.Enricher.EnrichRefs(ctx, source, refs)
		result.APICalls += enrichment.APICalls
		pullRequests = enrichment.PullRequests
	}
	for _, ref := range refs {
		mergeBase, mergeErr := repository.MergeBase(ctx, ref)
		if mergeErr != nil {
			result.Errors = append(result.Errors, ScanError{Code: "branch_merge_base_unavailable", Message: "branch merge-base evidence is incomplete", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: ref})
			continue
		}
		merged, mergedErr := repository.IsMerged(ctx, ref, mergeBase)
		if mergedErr != nil {
			result.Errors = append(result.Errors, ScanError{Code: "branch_merge_state_unavailable", Message: "branch merge state could not be resolved", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: ref})
			continue
		}
		if merged {
			continue
		}
		paths, changedErr := repository.ChangedPaths(ctx, mergeBase, ref)
		if changedErr != nil {
			result.Errors = append(result.Errors, ScanError{Code: "branch_comparison_incomplete", Message: "changed plan paths could not be read", SourceLocator: source.Locator, RepositoryID: decision.RepositoryID, Ref: ref})
			continue
		}
		observations, branchErrors := collectObservations(ctx, repository, decision.RepositoryID, ref, paths, "unmerged-branch", pullRequests[ref], now)
		result.Observations = append(result.Observations, observations...)
		result.Errors = append(result.Errors, branchErrors...)
	}
	return result
}

func readMembershipEvidence(ctx context.Context, repository *GitRepository, path string) ([]byte, error) {
	data, err := repository.ReadFile(ctx, repository.DefaultRef, path)
	if stderrors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return data, err
}

func collectObservations(ctx context.Context, repository *GitRepository, repositoryID, ref string, selected []string, visibility, pullRequest string, now time.Time) ([]plancatalog.SourceObservation, []ScanError) {
	files, err := repository.ListFiles(ctx, ref, "docs/repo/plans")
	if err != nil {
		return nil, []ScanError{{Code: "plan_tree_unavailable", Message: "plan tree could not be read", SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref}}
	}
	selectedSet := map[string]bool{}
	for _, path := range selected {
		selectedSet[filepath.ToSlash(path)] = true
	}
	records := []plancatalog.Record{}
	errors := []ScanError{}
	for _, file := range files {
		kind, formal := remoteFormalKind(file.Path, files)
		if !formal || (selected != nil && !selectedSet[file.Path]) {
			continue
		}
		data, readErr := repository.ReadFile(ctx, ref, file.Path)
		if readErr != nil {
			errors = append(errors, ScanError{Code: "plan_read_failed", Message: "canonical plan bytes could not be read", SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref})
			continue
		}
		record, parseErr := plancatalog.ParseDocument(file.Path, data, kind)
		record.FamilyQualified = kind != plancatalog.KindFamily || remoteFamilyQualified(filepath.ToSlash(filepath.Dir(file.Path)), files)
		if parseErr != nil || record.Metadata == nil {
			code := "remote_plan_invalid"
			if parseErr != nil {
				code = plancatalog.ErrorCode(parseErr)
			}
			errors = append(errors, ScanError{Code: code, Message: "canonical remote plan metadata is invalid", SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref})
			continue
		}
		records = append(records, record)
	}
	for _, problem := range plancatalog.ValidateRecords(records, plancatalog.ModeCanonical, repositoryID) {
		if problem.Severity == plancatalog.SeverityError {
			errors = append(errors, ScanError{Code: problem.Code, Message: problem.Message, SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref})
		}
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
		return nil, append(errors, ScanError{Code: "observation_revision_unavailable", Message: "observation commit could not be resolved", SourceLocator: repository.Source.Locator, RepositoryID: repositoryID, Ref: ref})
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
	return observations, errors
}

func remoteFormalKind(path string, files []GitTreeFile) (plancatalog.Kind, bool) {
	name := filepath.Base(path)
	switch {
	case strings.HasSuffix(name, "_discovery_doc.md"):
		return plancatalog.KindDiscovery, true
	case strings.HasSuffix(name, "_implementation_doc.md"):
		return plancatalog.KindImplementation, true
	case strings.EqualFold(name, "README.md") && filepath.ToSlash(path) != "docs/repo/plans/README.md" && remoteFamilyQualified(filepath.ToSlash(filepath.Dir(path)), files):
		return plancatalog.KindFamily, true
	default:
		return "", false
	}
}

func remoteFamilyQualified(directory string, files []GitTreeFile) bool {
	children := map[string]bool{}
	prefix := strings.TrimSuffix(directory, "/") + "/"
	for _, file := range files {
		if !strings.HasPrefix(file.Path, prefix) {
			continue
		}
		remainder := strings.TrimPrefix(file.Path, prefix)
		parts := strings.Split(remainder, "/")
		if len(parts) < 2 {
			continue
		}
		name := parts[len(parts)-1]
		if strings.HasSuffix(name, "_discovery_doc.md") || strings.HasSuffix(name, "_implementation_doc.md") {
			children[parts[0]] = true
		}
	}
	return len(children) >= 2
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
	var catalog Catalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return Catalog{}, err
	}
	if err := state.Validate(state.PlanCatalogSchema, catalog); err != nil {
		return Catalog{}, err
	}
	return catalog, nil
}
