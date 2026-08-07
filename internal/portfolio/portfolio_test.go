package portfolio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/plancatalog"
)

func TestPortfolioScanUsesRemoteBaselinesAndIndependentChangedBranchObservations(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	overlayPath := filepath.Join(fixture.home, filepath.FromSlash(OverlayPath))
	overlayBefore := mustReadPortfolioFile(t, overlayPath)
	memberRefsBefore := portfolioGitOutput(t, fixture.member, "for-each-ref", "--format=%(refname):%(objectname)", "refs/heads", "refs/remotes")
	memberStatusBefore := portfolioGitOutput(t, fixture.member, "status", "--porcelain=v1")
	now := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config,
		Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{}, Now: now, WriteCache: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Catalog.Complete || !result.CacheUpdated {
		t.Fatalf("scan result = %+v", result)
	}
	if len(result.Catalog.Members) != 2 {
		t.Fatalf("members = %+v", result.Catalog.Members)
	}
	selfMembers := 0
	for _, member := range result.Catalog.Members {
		if member.SelfMember {
			selfMembers++
		}
	}
	if selfMembers != 1 {
		t.Fatalf("self member count = %d", selfMembers)
	}
	if len(result.Catalog.Observations) != 3 {
		t.Fatalf("observations = %+v", result.Catalog.Observations)
	}
	memberBaseline := 0
	memberBranch := 0
	staleBranch := 0
	conflicting := 0
	for _, observation := range result.Catalog.Observations {
		if observation.RepositoryID != "example-member" {
			continue
		}
		if observation.Visibility == "default" {
			memberBaseline++
		}
		if observation.Visibility == "unmerged-branch" {
			memberBranch++
			if observation.Stale {
				staleBranch++
			}
		}
		if observation.Conflict {
			conflicting++
		}
		if strings.Contains(observation.CanonicalPath, "local-only") {
			t.Fatalf("local-only work appeared in remote catalog: %+v", observation)
		}
	}
	if memberBaseline != 1 || memberBranch != 1 || staleBranch != 1 || conflicting != 2 {
		t.Fatalf("baseline=%d branch=%d stale=%d conflicting=%d", memberBaseline, memberBranch, staleBranch, conflicting)
	}
	if len(result.Catalog.Candidates) != 1 || result.Catalog.Candidates[0].Reason != "coordination_home_id_mismatch" {
		t.Fatalf("candidates = %+v", result.Catalog.Candidates)
	}
	if actual := mustReadPortfolioFile(t, overlayPath); !bytes.Equal(actual, overlayBefore) {
		t.Fatal("strategic overlay changed during scan")
	}
	if refsAfter := portfolioGitOutput(t, fixture.member, "for-each-ref", "--format=%(refname):%(objectname)", "refs/heads", "refs/remotes"); refsAfter != memberRefsBefore {
		t.Fatal("scanner changed developer repository refs")
	}
	if statusAfter := portfolioGitOutput(t, fixture.member, "status", "--porcelain=v1"); statusAfter != memberStatusBefore {
		t.Fatal("scanner changed developer worktree or index bytes")
	}
	runPortfolioGit(t, fixture.member, "push", "origin", "--delete", "feature/changed-plan")
	retired, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config,
		Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{}, Now: now.Add(30 * time.Minute), WriteCache: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !retired.Catalog.Complete || len(retired.Catalog.Observations) != 2 {
		t.Fatalf("deleted remote ref was not retired: %+v", retired.Catalog.Observations)
	}
	for _, observation := range retired.Catalog.Observations {
		if observation.Visibility == "unmerged-branch" {
			t.Fatalf("deleted branch observation survived complete refresh: %+v", observation)
		}
	}
	cacheBefore := mustReadPortfolioFile(t, filepath.Join(fixture.home, filepath.FromSlash(CatalogPath)))
	if _, _, err := storeCompleteCatalogWithHook(fixture.home, retired.Catalog, func(phase string) error {
		if phase == "catalog-staged" {
			return fmt.Errorf("injected interruption")
		}
		return nil
	}); err == nil || !strings.Contains(err.Error(), "injected interruption") {
		t.Fatalf("interrupted cache write error = %v", err)
	}
	if cacheAfter := mustReadPortfolioFile(t, filepath.Join(fixture.home, filepath.FromSlash(CatalogPath))); !bytes.Equal(cacheAfter, cacheBefore) {
		t.Fatal("interrupted staged write changed the complete cache")
	}
	offlineRemote := filepath.Join(fixture.parent, "member-remote-offline.git")
	if err := os.Rename(filepath.Join(fixture.parent, "member-remote.git"), offlineRemote); err != nil {
		t.Fatal(err)
	}
	failure, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{}, Now: now.Add(time.Hour), WriteCache: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if failure.Catalog.Complete || failure.CacheUpdated || !failure.PreviousPreserved {
		t.Fatalf("incomplete scan result = %+v", failure)
	}
	if failure.Catalog.LastCompleteScanAt != retired.Catalog.CompletedAt {
		t.Fatalf("incomplete scan did not label preserved v2 cache: last_complete=%q want=%q", failure.Catalog.LastCompleteScanAt, retired.Catalog.CompletedAt)
	}
	if cacheAfter := mustReadPortfolioFile(t, filepath.Join(fixture.home, filepath.FromSlash(CatalogPath))); !bytes.Equal(cacheAfter, cacheBefore) {
		t.Fatal("incomplete scan replaced the prior complete cache")
	}
	if actual := mustReadPortfolioFile(t, overlayPath); !bytes.Equal(actual, overlayBefore) {
		t.Fatal("failed scan changed strategic overlay bytes")
	}
}

func TestRemoteClassifierMatchesLocalCommitTreeAcrossNestedDocsAndOwnership(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	configData := strings.Replace(portfolioConfig("member", "example-member", "example-coordination"), "    plan_catalog_discovery_version: 2\n", "    plan_catalog_discovery_version: 2\n    plan_catalog_ownership:\n      excluded_roots:\n        - third_party/\n", 1)
	writePortfolioFile(t, fixture.member, ConfigPath, []byte(configData))
	writePortfolioPlanAt(t, fixture.member, "docs/root/root_discovery_doc.md", "Root Docs", "example-member.discovery.root", plancatalog.KindDiscovery, "root docs candidate")
	writePortfolioPlanAt(t, fixture.member, "products/widget/source/deep/docs/domain/plans/deep_implementation_doc.md", "Deep Docs", "example-member.implementation.deep", plancatalog.KindImplementation, "deep docs candidate")
	writePortfolioPlanAt(t, fixture.member, "products/widget/docs/domain/unconventional.md", "Metadata Only", "example-member.discovery.metadata-only", plancatalog.KindDiscovery, "metadata-only candidate")
	writePortfolioPlanAt(t, fixture.member, "products/widget/docs/domain/README.md", "Domain Family", "example-member.family.domain", plancatalog.KindFamily, "metadata-qualified family")
	writePortfolioFile(t, fixture.member, "products/widget/docs/router/README.md", []byte("# Ordinary router\n\nNo plan metadata.\n"))
	writePortfolioPlanAt(t, fixture.member, "third_party/copied/docs/plans/copied_discovery_doc.md", "Excluded Copy", "example-member.discovery.excluded", plancatalog.KindDiscovery, "excluded candidate")
	runPortfolioGit(t, fixture.member, "add", ".")
	runPortfolioGit(t, fixture.member, "commit", "-m", "Add repository-wide plan discovery fixture")
	runPortfolioGit(t, fixture.member, "push", "origin", "main")

	settings, problems := plancatalog.DecodeRepositorySettings([]byte(configData))
	if plancatalog.HasErrors(problems) {
		t.Fatalf("settings problems: %+v", problems)
	}
	local, err := plancatalog.ClassifyCommitTree(fixture.member, "main", settings, settings.Mode, settings.RepositoryID)
	if err != nil {
		t.Fatal(err)
	}
	scanRoot := t.TempDir()
	manager := MirrorManager{RepositoryRoot: scanRoot, Root: filepath.Join(scanRoot, filepath.FromSlash(MirrorRootPath)), Runner: ExecRunner{}}
	repository, err := manager.Refresh(context.Background(), RepositorySource{Kind: "local-git", Locator: fixture.member, CloneURL: filepath.Join(fixture.parent, "member-remote.git"), DefaultBranch: "main", NameHint: "member"})
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	remote, err := classifyRemoteTree(context.Background(), repository, settings.RepositoryID, repository.DefaultRef, settings)
	if err != nil {
		t.Fatal(err)
	}
	if local.PolicyDigest != remote.PolicyDigest || local.CandidateSetDigest != remote.CandidateSetDigest {
		t.Fatalf("classifier digests differ local=%s/%s remote=%s/%s", local.PolicyDigest, local.CandidateSetDigest, remote.PolicyDigest, remote.CandidateSetDigest)
	}
	for index := range local.Candidates {
		local.Candidates[index].Provenance.Source.Revision = ""
	}
	for index := range remote.Candidates {
		remote.Candidates[index].Provenance.Source.Revision = ""
	}
	localCandidates, _ := json.Marshal(local.Candidates)
	remoteCandidates, _ := json.Marshal(remote.Candidates)
	localDiscovery, _ := json.Marshal(local.Discovery)
	remoteDiscovery, _ := json.Marshal(remote.Discovery)
	if !bytes.Equal(localCandidates, remoteCandidates) || !bytes.Equal(localDiscovery, remoteDiscovery) {
		t.Fatalf("remote classifier diverged from local commit tree\nlocal candidates=%s\nremote candidates=%s\nlocal discovery=%s\nremote discovery=%s", localCandidates, remoteCandidates, localDiscovery, remoteDiscovery)
	}
	paths := map[string]bool{}
	for _, candidate := range remote.Candidates {
		paths[candidate.Path] = true
	}
	for _, expected := range []string{"docs/root/root_discovery_doc.md", "products/widget/source/deep/docs/domain/plans/deep_implementation_doc.md", "products/widget/docs/domain/unconventional.md", "products/widget/docs/domain/README.md", "third_party/copied/docs/plans/copied_discovery_doc.md"} {
		if !paths[expected] {
			t.Fatalf("remote candidate missing: %s", expected)
		}
	}
	if paths["products/widget/docs/router/README.md"] {
		t.Fatal("ordinary README gained family authority")
	}
}

func TestRemoteGitlinkPlanSignalIsHardUnowned(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	head := strings.TrimSpace(portfolioGitOutput(t, fixture.member, "rev-parse", "HEAD"))
	gitlinkPath := "products/widget/docs/embedded/embedded_discovery_doc.md"
	runPortfolioGit(t, fixture.member, "update-index", "--add", "--cacheinfo", "160000,"+head+","+gitlinkPath)
	runPortfolioGit(t, fixture.member, "commit", "-m", "Add plan-like gitlink")
	runPortfolioGit(t, fixture.member, "push", "origin", "main")

	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{}, Now: time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	if result.Catalog.Complete {
		t.Fatal("plan-like gitlink was accepted as remote outer-commit authority")
	}
	found := false
	for _, scanError := range result.Catalog.Errors {
		if scanError.RepositoryID == "example-member" && scanError.Code == "plan_path_unowned" && strings.Contains(scanError.Message, "gitlink") {
			found = true
		}
	}
	if !found {
		t.Fatalf("gitlink boundary evidence missing: %+v", result.Catalog.Errors)
	}
}

func TestV1RequiredMemberAndHistoricalV1CacheRemainIncompleteEvidence(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	v1Config := strings.Replace(portfolioConfig("member", "example-member", "example-coordination"), "    plan_catalog_discovery_version: 2\n", "", 1)
	writePortfolioFile(t, fixture.member, ConfigPath, []byte(v1Config))
	commitAndPushPortfolioFile(t, fixture.member, ConfigPath, "Retain discovery v1")
	historical := Catalog{
		SchemaVersion: 1, CoordinationHomeID: "example-coordination", StartedAt: "2026-07-01T12:00:00Z", CompletedAt: "2026-07-01T12:00:00Z", LastCompleteScanAt: "2026-07-01T12:00:00Z", Complete: true,
		Members: []plancatalog.CatalogMember{}, Observations: []plancatalog.SourceObservation{}, Candidates: []Candidate{}, Errors: []ScanError{}, Metrics: Metrics{MaxConcurrency: 1},
	}
	data, err := json.Marshal(historical)
	if err != nil {
		t.Fatal(err)
	}
	writePortfolioFile(t, fixture.home, CatalogPath, append(data, '\n'))
	before := mustReadPortfolioFile(t, filepath.Join(fixture.home, filepath.FromSlash(CatalogPath)))
	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{}, Now: time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC), WriteCache: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.Catalog.Complete || result.Catalog.LastCompleteScanAt != "" || result.CacheUpdated || !result.PreviousPreserved {
		t.Fatalf("v1 evidence satisfied v2 completeness: %+v", result)
	}
	found := false
	for _, member := range result.Catalog.Members {
		if member.RepositoryID == "example-member" && member.DiscoveryVersion == plancatalog.DiscoveryV1 && member.Complete != nil && !*member.Complete {
			found = true
		}
	}
	if !found {
		t.Fatalf("v1 member evidence missing: %+v", result.Catalog.Members)
	}
	if after := mustReadPortfolioFile(t, filepath.Join(fixture.home, filepath.FromSlash(CatalogPath))); !bytes.Equal(before, after) {
		t.Fatal("incomplete v2 scan replaced historical v1 cache")
	}
}

func TestDiscoveryV2LegacyModeStillRequiresCanonicalRemoteEvidence(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	legacyConfig := strings.Replace(portfolioConfig("member", "example-member", "example-coordination"), "    plan_catalog_mode: canonical\n", "    plan_catalog_mode: legacy\n", 1)
	writePortfolioFile(t, fixture.member, ConfigPath, []byte(legacyConfig))
	planPath := "docs/repo/plans/member-plan/member-plan_discovery_doc.md"
	writePortfolioFile(t, fixture.member, planPath, []byte("Last updated: 2026-06-01T12:00:00Z (UTC)\nCreated: 2026-06-01\nStatus: active\n\n# Filename Only Legacy Record\n"))
	runPortfolioGit(t, fixture.member, "add", ConfigPath, planPath)
	runPortfolioGit(t, fixture.member, "commit", "-m", "Retain filename-only legacy record under discovery v2")
	runPortfolioGit(t, fixture.member, "push", "origin", "main")

	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{}, Now: time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC), WriteCache: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.Catalog.Complete || result.CacheUpdated {
		t.Fatalf("filename-only v2 member was reported complete: %+v", result)
	}
	metadataMissing := false
	for _, scanError := range result.Catalog.Errors {
		if scanError.RepositoryID == "example-member" && scanError.Ref == "refs/remotes/origin/main" && scanError.Path == planPath && scanError.Code == "metadata_missing" {
			metadataMissing = true
		}
	}
	if !metadataMissing {
		t.Fatalf("canonical remote coverage blocker missing: %+v", result.Catalog.Errors)
	}
	for _, observation := range result.Catalog.Observations {
		if observation.RepositoryID == "example-member" && observation.Ref == "refs/remotes/origin/main" && observation.CanonicalPath == planPath {
			t.Fatal("filename-only record gained a verified stable-ID observation")
		}
	}
	candidateFound := false
	for _, candidate := range result.Catalog.PlanCandidates {
		if candidate.RepositoryID == "example-member" && candidate.Ref == "refs/remotes/origin/main" && candidate.Path == planPath && candidate.Signal == plancatalog.SignalFilename {
			candidateFound = true
		}
	}
	if !candidateFound {
		t.Fatalf("filename-only candidate evidence missing: %+v", result.Catalog.PlanCandidates)
	}
}

func TestRemoteCatalogInventoryAndListShareReviewedCompatibilityTitle(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	planPath := "docs/repo/plans/member-plan/member-plan_discovery_doc.md"
	planData := mustReadPortfolioFile(t, filepath.Join(fixture.member, filepath.FromSlash(planPath)))
	lineEnding := "\n"
	if bytes.Contains(planData, []byte("\r\n")) {
		lineEnding = "\r\n"
	}
	originalHeading := []byte("# Member Plan" + lineEnding)
	if !bytes.Contains(planData, originalHeading) {
		t.Fatalf("fixture plan heading missing from %s", planPath)
	}
	compatibilityHeading := []byte("# Document Header" + lineEnding + lineEnding + "## Overview" + lineEnding)
	planData = bytes.Replace(planData, originalHeading, compatibilityHeading, 1)
	writePortfolioFile(t, fixture.member, planPath, planData)
	register := "## MEMBER-PR-001 - Member Semantic Plan Title\n\nCanonical docs:\n- `" + planPath + "`\n"
	writePortfolioFile(t, fixture.member, plancatalog.LegacyRegisterPath, []byte(register))
	runPortfolioGit(t, fixture.member, "add", planPath, plancatalog.LegacyRegisterPath)
	runPortfolioGit(t, fixture.member, "commit", "-m", "Use compatibility plan title")
	runPortfolioGit(t, fixture.member, "push", "origin", "main")

	snapshot, err := plancatalog.LoadRepositorySnapshot(fixture.member)
	if err != nil {
		t.Fatal(err)
	}
	view := plancatalog.BuildView(snapshot)
	if len(view.Rows) != 1 || view.Rows[0].Title != "Member Semantic Plan Title" {
		t.Fatalf("local semantic view = %+v", view.Rows)
	}
	inventory, err := plancatalog.BuildInventory(fixture.member, time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory.Records) != 1 || inventory.Records[0].Title != view.Rows[0].Title {
		t.Fatalf("inventory title = %+v, view title = %+v", inventory.Records, view.Rows)
	}

	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{},
		Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Catalog.Complete {
		t.Fatalf("remote catalog incomplete: %+v", result.Catalog.Errors)
	}
	for _, observation := range result.Catalog.Observations {
		if observation.RepositoryID == "example-member" && observation.PlanID == "example-member.discovery.member-plan" && observation.Visibility == "default" {
			if observation.Title != view.Rows[0].Title {
				t.Fatalf("remote title = %q, local title = %q", observation.Title, view.Rows[0].Title)
			}
			return
		}
	}
	t.Fatal("member default observation missing")
}

func TestMembershipRequiresDefaultBranchKitLockV2IdentityAndMatchingHome(t *testing.T) {
	validConfig := []byte(portfolioConfig("member", "example-member", "example-home"))
	valid := MembershipInput{ConfigData: validConfig, LockData: []byte(testLockYAML), KitMarker: []byte("kit\n"), HomeID: "example-home"}
	if decision := EvaluateMembership(valid); !decision.Member || decision.RepositoryID != "example-member" || decision.PlanSettings.DiscoveryVersion != plancatalog.DiscoveryV2 || plancatalog.HasErrors(decision.PlanProblems) {
		t.Fatalf("valid membership = %+v", decision)
	}
	wrong := valid
	wrong.HomeID = "other-home"
	if decision := EvaluateMembership(wrong); decision.Member || decision.Reason != "coordination_home_id_mismatch" {
		t.Fatalf("mismatched membership = %+v", decision)
	}
	missingMarker := valid
	missingMarker.KitMarker = nil
	if decision := EvaluateMembership(missingMarker); decision.Member || decision.Reason != "default_branch_kit_marker_missing" {
		t.Fatalf("missing-marker membership = %+v", decision)
	}
	sourceRepository := valid
	sourceRepository.KitMarker = nil
	sourceRepository.LockData = nil
	sourceRepository.KitSourceMarker = []byte(testKitSourceMarkerYAML)
	if decision := EvaluateMembership(sourceRepository); !decision.Member || decision.RepositoryID != "example-member" {
		t.Fatalf("Kit source membership = %+v", decision)
	}
	invalidSourceRepository := sourceRepository
	invalidSourceRepository.KitSourceMarker = []byte(strings.Replace(testKitSourceMarkerYAML, "schema_version: 1", "schema_version: 99", 1))
	if decision := EvaluateMembership(invalidSourceRepository); decision.Member || !decision.Incomplete || !strings.HasPrefix(decision.Reason, "default_branch_kit_source_marker_invalid:") {
		t.Fatalf("invalid Kit source membership = %+v", decision)
	}
	legacy := valid
	legacyConfig := strings.Replace(portfolioConfig("member", "example-member", "example-home"), "  schema_version: 2\n", "", 1)
	legacyConfig = strings.Replace(legacyConfig, "  coordination_home_id: example-home\n", "  coordination_home_id: example-home\n  coordination_home_path: ../home\n  coordination_home_register_path: docs/repo/plans/plan-register.md\n", 1)
	legacy.ConfigData = []byte(legacyConfig)
	if decision := EvaluateMembership(legacy); decision.Member || decision.Reason != "portfolio_v2_required_for_enrollment" {
		t.Fatalf("v1 membership = %+v", decision)
	}
}

func TestPortfolioScanEnrollsKitSourceRepositoryWithoutConsumerInstallation(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	for _, relative := range []string{LockPath, KitMarkerPath} {
		if err := os.Remove(filepath.Join(fixture.member, filepath.FromSlash(relative))); err != nil {
			t.Fatal(err)
		}
	}
	writePortfolioFile(t, fixture.member, KitSourceMarkerPath, []byte(testKitSourceMarkerYAML))
	runPortfolioGit(t, fixture.member, "add", "-A", LockPath, KitMarkerPath, KitSourceMarkerPath)
	runPortfolioGit(t, fixture.member, "commit", "-m", "Use producer source marker")
	runPortfolioGit(t, fixture.member, "push", "origin", "main")

	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})},
		Runner: ExecRunner{}, Now: time.Date(2026, 8, 1, 8, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Catalog.Complete {
		t.Fatalf("source-repository scan incomplete: %+v", result.Catalog.Errors)
	}
	found := false
	for _, member := range result.Catalog.Members {
		if member.RepositoryID == "example-member" {
			found = true
		}
	}
	if !found {
		t.Fatalf("source repository was not enrolled: members=%+v candidates=%+v", result.Catalog.Members, result.Catalog.Candidates)
	}
}

func TestMalformedMembershipAndFailedHomeSelfMembershipPreserveCompleteCache(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, portfolioGitFixture)
	}{
		{
			name: "malformed member config",
			mutate: func(t *testing.T, fixture portfolioGitFixture) {
				writePortfolioFile(t, fixture.member, ConfigPath, []byte("portfolio: [\n"))
				commitAndPushPortfolioFile(t, fixture.member, ConfigPath, "Corrupt member config")
			},
		},
		{
			name: "malformed member lock",
			mutate: func(t *testing.T, fixture portfolioGitFixture) {
				writePortfolioFile(t, fixture.member, LockPath, []byte("schema_version: nope\n"))
				commitAndPushPortfolioFile(t, fixture.member, LockPath, "Corrupt member lock")
			},
		},
		{
			name: "failed coordination home self membership",
			mutate: func(t *testing.T, fixture portfolioGitFixture) {
				writePortfolioFile(t, fixture.home, ConfigPath, []byte(portfolioConfig("member", "example-home", "example-coordination")))
				commitAndPushPortfolioFile(t, fixture.home, ConfigPath, "Break home self membership")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newPortfolioGitFixture(t)
			now := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
			config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
			initial, err := Scan(context.Background(), ScanOptions{
				Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})},
				Runner: ExecRunner{}, Now: now, WriteCache: true,
			})
			if err != nil || !initial.Catalog.Complete || !initial.CacheUpdated {
				t.Fatalf("initial scan = %+v, err = %v", initial, err)
			}
			cachePath := filepath.Join(fixture.home, filepath.FromSlash(CatalogPath))
			before := mustReadPortfolioFile(t, cachePath)
			test.mutate(t, fixture)
			failed, err := Scan(context.Background(), ScanOptions{
				Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})},
				Runner: ExecRunner{}, Now: now.Add(time.Hour), WriteCache: true,
			})
			if err != nil {
				t.Fatal(err)
			}
			if failed.Catalog.Complete || failed.CacheUpdated || !failed.PreviousPreserved {
				t.Fatalf("failed membership scan = %+v", failed)
			}
			found := false
			for _, scanError := range failed.Catalog.Errors {
				if scanError.Code == "membership_validation_failed" {
					found = true
				}
			}
			if !found {
				t.Fatalf("membership blocker missing: %+v", failed.Catalog.Errors)
			}
			if after := mustReadPortfolioFile(t, cachePath); !bytes.Equal(before, after) {
				t.Fatal("failed membership scan replaced the prior complete cache")
			}
		})
	}
}

func TestLoadConfigNormalizesV1V2CommittedAndLocalSources(t *testing.T) {
	fixtureRoot := filepath.Join("..", "..", "tests", "fixtures", "portfolio")
	t.Run("v2 coordination home", func(t *testing.T) {
		root := t.TempDir()
		writePortfolioFile(t, root, ConfigPath, mustReadPortfolioFile(t, filepath.Join(fixtureRoot, "kit-config-v2-home.yaml")))
		writePortfolioFile(t, root, LocalSourcesPath, mustReadPortfolioFile(t, filepath.Join(fixtureRoot, "local-sources.yaml")))
		config, err := LoadConfig(root)
		if err != nil {
			t.Fatal(err)
		}
		if config.SchemaVersion != 2 || config.Role != RoleCoordinationHome || config.RepositoryID != "example-home-repository" || config.CoordinationHomeID != "example-home" || len(config.ProviderSources) != 1 || len(config.LocalSources) != 1 {
			t.Fatalf("normalized v2 config = %+v", config)
		}
	})
	t.Run("v1 compatibility", func(t *testing.T) {
		root := t.TempDir()
		writePortfolioFile(t, root, ConfigPath, mustReadPortfolioFile(t, filepath.Join(fixtureRoot, "kit-config-v1-member.yaml")))
		config, err := LoadConfig(root)
		if err != nil {
			t.Fatal(err)
		}
		if config.SchemaVersion != 1 || !config.CompatibilityV1 || config.Role != RoleMember || config.RepositoryID != "Legacy-Member" {
			t.Fatalf("normalized v1 config = %+v", config)
		}
	})
	t.Run("v1 path-only coordination home compatibility", func(t *testing.T) {
		root := t.TempDir()
		writePortfolioFile(t, root, ConfigPath, mustReadPortfolioFile(t, filepath.Join(fixtureRoot, "kit-config-v1-home.yaml")))
		config, err := LoadConfig(root)
		if err != nil {
			t.Fatal(err)
		}
		if config.SchemaVersion != 1 || !config.CompatibilityV1 || config.Role != RoleCoordinationHome || config.RepositoryID != "" || config.CoordinationHomeID != "" {
			t.Fatalf("normalized v1 coordination home = %+v", config)
		}
	})
	t.Run("missing home repository identity", func(t *testing.T) {
		root := t.TempDir()
		writePortfolioFile(t, root, ConfigPath, mustReadPortfolioFile(t, filepath.Join(fixtureRoot, "kit-config-v2-home-missing-repository-id.yaml")))
		if _, err := LoadConfig(root); err == nil || !strings.Contains(err.Error(), "portfolio_config_invalid") {
			t.Fatalf("missing identity error = %v", err)
		}
	})
	t.Run("member ignores discovery sources", func(t *testing.T) {
		root := t.TempDir()
		writePortfolioFile(t, root, ConfigPath, []byte(portfolioConfig("member", "example-member", "example-home")))
		writePortfolioFile(t, root, LocalSourcesPath, []byte("not: valid portfolio sources\n"))
		config, err := LoadConfig(root)
		if err != nil {
			t.Fatal(err)
		}
		config.ProviderSources = []ProviderSource{{Kind: "github-owner", Owner: "must-not-run"}}
		if sources := BuildSources(config, ExecRunner{}); len(sources) != 0 {
			t.Fatalf("member authorized discovery sources: %+v", sources)
		}
	})
}

func TestEquivalentRemoteURLsDeduplicateToTheCoordinationHomeSelfSource(t *testing.T) {
	enricher := &staticRefEnricher{}
	sources := []RepositorySource{
		{Kind: "github", Locator: "Example/Home", CloneURL: "https://github.com/Example/Home.git", NameHint: "Home", Enricher: enricher},
		{Kind: "local-git", Locator: "/coordination/home", CloneURL: "git@github.com:example/home.git", NameHint: "Home", Self: true},
	}
	result := deduplicateSources(sources)
	if len(result) != 1 || !result[0].Self || result[0].Locator != "/coordination/home" || result[0].Enricher == nil {
		t.Fatalf("equivalent remote deduplication = %+v", result)
	}
	for _, remote := range []string{
		"https://github.com/Example/Home.git",
		"ssh://git@github.com/example/home.git",
		"git@github.com:EXAMPLE/HOME.git",
	} {
		if identity := canonicalRemoteIdentity(remote); identity != "github.com/example/home" {
			t.Fatalf("canonical identity for %q = %q", remote, identity)
		}
	}
}

func TestCoordinationHomeSelfDiscoveryNormalizesRelativeOriginRemote(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	runPortfolioGit(t, fixture.home, "remote", "set-url", "origin", "../home-remote.git")
	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Runner: ExecRunner{},
		Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Catalog.Complete || len(result.Catalog.Members) != 1 || !result.Catalog.Members[0].SelfMember || result.Catalog.Members[0].RepositoryID != "example-home" {
		t.Fatalf("relative-origin self scan = %+v", result.Catalog)
	}
}

func TestGitHubSourcePaginatesRetriesRedactsAndHonorsCancellation(t *testing.T) {
	runner := &fakeGHRunner{}
	source := NewGitHubSource([]string{"example-org"}, runner)
	result, err := source.Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Complete || len(result.Repositories) != 101 || runner.pageOneAttempts != 2 {
		t.Fatalf("GitHub discovery = %+v attempts=%d", result, runner.pageOneAttempts)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cancelled := NewGitHubSource([]string{"example-org"}, cancelledRunner{})
	result, err = cancelled.Discover(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if result.Complete || len(result.Errors) != 1 || result.Errors[0].Code != "github_command_cancelled" {
		t.Fatalf("cancelled discovery = %+v", result)
	}
	if strings.Contains(fmt.Sprintf("%+v", result.Errors), "super-secret-token") {
		t.Fatal("provider diagnostics exposed a secret-bearing command value")
	}
	for _, test := range []struct {
		name         string
		mode         string
		expectedCode string
		expectedRuns int
	}{
		{name: "rate limit", mode: "rate", expectedCode: "github_rate_limit_exhausted", expectedRuns: 1},
		{name: "permission", mode: "permission", expectedCode: "github_permission_denied", expectedRuns: 1},
		{name: "retry exhaustion", mode: "retry", expectedCode: "github_retry_exhausted", expectedRuns: 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			runner := &ghFailureRunner{mode: test.mode}
			result, err := NewGitHubSource([]string{"example-org"}, runner).Discover(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if result.Complete || len(result.Errors) != 1 || result.Errors[0].Code != test.expectedCode || runner.apiRuns != test.expectedRuns {
				t.Fatalf("failure result=%+v runs=%d", result, runner.apiRuns)
			}
			if strings.Contains(fmt.Sprintf("%+v", result.Errors), "super-secret-token") {
				t.Fatal("failure diagnostics exposed provider stderr")
			}
		})
	}
	t.Run("truncation", func(t *testing.T) {
		previousLimit := githubPageLimit
		githubPageLimit = 2
		defer func() { githubPageLimit = previousLimit }()
		runner := &ghFailureRunner{mode: "full-pages"}
		result, err := NewGitHubSource([]string{"example-org"}, runner).Discover(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if result.Complete || len(result.Errors) != 1 || result.Errors[0].Code != "github_response_truncated" || runner.apiRuns != 2 {
			t.Fatalf("truncated result=%+v runs=%d", result, runner.apiRuns)
		}
	})
}

func TestGitHubPullRequestEnrichmentPaginatesAndCorrelatesOnlyRequestedRefs(t *testing.T) {
	runner := &prPaginationRunner{}
	source := NewGitHubSource([]string{"example-org"}, runner)
	ref := "refs/remotes/origin/feature/changed-plan"
	forkCollision := "refs/remotes/origin/feature/fork-collision"
	result, err := source.EnrichRefs(context.Background(), RepositorySource{Kind: "github", Locator: "example-org/example-repo"}, []string{ref, forkCollision})
	if err != nil {
		t.Fatal(err)
	}
	if result.APICalls != 2 || runner.apiRuns != 2 {
		t.Fatalf("PR pagination calls result=%d runner=%d", result.APICalls, runner.apiRuns)
	}
	if result.PullRequests[ref] != "https://github.com/example-org/example-repo/pull/101" || len(result.PullRequests) != 1 {
		t.Fatalf("PR correlation = %+v", result.PullRequests)
	}
	if result.PullRequests[forkCollision] != "" {
		t.Fatalf("fork PR with colliding head ref enriched origin branch: %+v", result.PullRequests)
	}
}

func TestOptionalPullRequestFactsEnrichWithoutSelectingRemoteBranches(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	runPortfolioGit(t, fixture.member, "checkout", "-b", "feature/no-pr", "main")
	writePortfolioPlan(t, fixture.member, "no-pr/no-pr_discovery_doc.md", "No PR Plan", "example-member.discovery.no-pr", "branch without pull request")
	runPortfolioGit(t, fixture.member, "add", "docs/repo/plans")
	runPortfolioGit(t, fixture.member, "commit", "-m", "Add branch without pull request")
	runPortfolioGit(t, fixture.member, "push", "origin", "feature/no-pr")
	runPortfolioGit(t, fixture.member, "checkout", "main")
	enricher := &staticRefEnricher{pullRequests: map[string]string{
		"refs/remotes/origin/feature/changed-plan": "https://github.com/example-org/example-member/pull/45",
	}, apiCalls: 2}
	memberRemote := filepath.Join(fixture.parent, "member-remote.git")
	source := staticPortfolioSource{repository: RepositorySource{
		Kind: "github", Locator: "example-org/example-member", CloneURL: memberRemote, DefaultBranch: "main", NameHint: "example-member", Enricher: enricher,
	}}
	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Sources: []Source{source}, Runner: ExecRunner{},
		Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Catalog.Complete || result.Catalog.Metrics.APICallCount != 2 {
		t.Fatalf("enriched scan = %+v", result.Catalog)
	}
	changedSeen := false
	noPRSeen := false
	for _, observation := range result.Catalog.Observations {
		if observation.Ref == "refs/remotes/origin/feature/changed-plan" {
			changedSeen = true
			if observation.PullRequest != "https://github.com/example-org/example-member/pull/45" {
				t.Fatalf("matching branch PR = %q", observation.PullRequest)
			}
		}
		if observation.Ref == "refs/remotes/origin/feature/no-pr" {
			noPRSeen = true
			if observation.PullRequest != "" {
				t.Fatalf("branch without PR was enriched: %q", observation.PullRequest)
			}
		}
	}
	if !changedSeen || !noPRSeen {
		t.Fatalf("PR facts selected branch coverage: changed=%v no-pr=%v observations=%+v", changedSeen, noPRSeen, result.Catalog.Observations)
	}
	if !enricher.seen["refs/remotes/origin/feature/changed-plan"] || !enricher.seen["refs/remotes/origin/feature/no-pr"] {
		t.Fatalf("enricher did not receive all independently selected refs: %+v", enricher.seen)
	}
	t.Logf(
		"portfolio evidence repositories=%d branches=%d plans=%d api_calls=%d duration_ms=%d max_concurrency=%d",
		len(result.Catalog.Members), 2, len(result.Catalog.Observations),
		result.Catalog.Metrics.APICallCount, result.Catalog.Metrics.DurationMS,
		result.Catalog.Metrics.MaxConcurrency,
	)
}

func TestLocalGitSourceIgnoresNonRepositoriesAndReportsMissingRemotesAsCandidates(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	runPortfolioGit(t, "", "init", repo)
	if err := os.Mkdir(filepath.Join(root, "ordinary"), 0o755); err != nil {
		t.Fatal(err)
	}
	result, err := NewLocalGitSource(root, ExecRunner{}).Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Complete || len(result.Repositories) != 0 || len(result.Errors) != 1 || result.Errors[0].Code != "local_remote_unavailable" {
		t.Fatalf("local discovery = %+v", result)
	}
}

func TestConcurrentPortfolioScanIsRejectedBeforeMirrorOrCacheWork(t *testing.T) {
	root := t.TempDir()
	lock, err := acquireScanLock(root)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	_, err = Scan(context.Background(), ScanOptions{Root: root, Config: Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "home", CoordinationHomeID: "example-home", Root: root}, Runner: ExecRunner{}, Now: time.Now()})
	if err == nil || !strings.Contains(err.Error(), "portfolio_scan_in_progress") {
		t.Fatalf("concurrent scan error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(CatalogPath))); !os.IsNotExist(err) {
		t.Fatalf("concurrent scan changed cache state: %v", err)
	}
}

func TestCacheReadsRejectIntermediateSymlinksAndStableLockSurvivesCacheParentReplacement(t *testing.T) {
	t.Run("intermediate symlink", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("symlink creation requires elevated Windows privileges on some runners")
		}
		root := t.TempDir()
		inside := filepath.Join(root, "redirected-local")
		if err := os.MkdirAll(filepath.Join(inside, "portfolio"), 0o700); err != nil {
			t.Fatal(err)
		}
		data, err := json.Marshal(portfolioCatalogFixture("2026-07-31T12:00:00Z"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(inside, "portfolio", "catalog.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(filepath.Join(root, ".codeheart"), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("../redirected-local", filepath.Join(root, ".codeheart", "local")); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadCachedCatalog(root); err == nil || !strings.Contains(err.Error(), "portfolio_local_path_unsafe") {
			t.Fatalf("intermediate symlink cache read error = %v", err)
		}
	})

	t.Run("cache parent replacement", func(t *testing.T) {
		root := t.TempDir()
		if updated, _, err := StoreCompleteCatalog(root, portfolioCatalogFixture("2026-07-31T12:00:00Z")); err != nil || !updated {
			t.Fatalf("initial cache write updated=%v err=%v", updated, err)
		}
		lock, err := acquireScanLock(root)
		if err != nil {
			t.Fatal(err)
		}
		defer lock.Close()
		parent := filepath.Join(root, filepath.FromSlash(filepath.Dir(CatalogPath)))
		displaced := filepath.Join(root, ".displaced-cache-parent")
		if err := os.Rename(parent, displaced); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(parent, 0o700); err != nil {
			t.Fatal(err)
		}
		sentinel := []byte("replacement cache parent\n")
		if err := os.WriteFile(filepath.Join(parent, "sentinel"), sentinel, 0o600); err != nil {
			t.Fatal(err)
		}
		_, err = Scan(context.Background(), ScanOptions{
			Root:   root,
			Config: Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "home", CoordinationHomeID: "example-home", Root: root},
			Runner: ExecRunner{}, Now: time.Now(), WriteCache: true,
		})
		if err == nil || !strings.Contains(err.Error(), "portfolio_scan_in_progress") {
			t.Fatalf("second scan after cache-parent replacement error = %v", err)
		}
		if actual := mustReadPortfolioFile(t, filepath.Join(parent, "sentinel")); !bytes.Equal(actual, sentinel) {
			t.Fatal("second scan changed replacement cache parent")
		}
		entries, err := os.ReadDir(parent)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 1 || entries[0].Name() != "sentinel" {
			t.Fatalf("second scan mutated replacement cache parent: %+v", entries)
		}
	})

	t.Run("local namespace replacement", func(t *testing.T) {
		root := t.TempDir()
		if updated, _, err := StoreCompleteCatalog(root, portfolioCatalogFixture("2026-07-31T12:00:00Z")); err != nil || !updated {
			t.Fatalf("initial cache write updated=%v err=%v", updated, err)
		}
		lock, err := acquireScanLock(root)
		if err != nil {
			t.Fatal(err)
		}
		defer lock.Close()
		localNamespace := filepath.Join(root, ".codeheart")
		displaced := filepath.Join(root, ".displaced-codeheart")
		if err := os.Rename(localNamespace, displaced); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(localNamespace, 0o700); err != nil {
			t.Fatal(err)
		}
		sentinel := []byte("replacement local namespace\n")
		if err := os.WriteFile(filepath.Join(localNamespace, "sentinel"), sentinel, 0o600); err != nil {
			t.Fatal(err)
		}
		_, err = Scan(context.Background(), ScanOptions{
			Root:   root,
			Config: Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "home", CoordinationHomeID: "example-home", Root: root},
			Runner: ExecRunner{}, Now: time.Now(), WriteCache: true,
		})
		if err == nil || !strings.Contains(err.Error(), "portfolio_scan_in_progress") {
			t.Fatalf("second scan after local-namespace replacement error = %v", err)
		}
		if actual := mustReadPortfolioFile(t, filepath.Join(localNamespace, "sentinel")); !bytes.Equal(actual, sentinel) {
			t.Fatal("second scan changed replacement local namespace")
		}
		entries, err := os.ReadDir(localNamespace)
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 1 || entries[0].Name() != "sentinel" {
			t.Fatalf("second scan mutated replacement local namespace: %+v", entries)
		}
	})
}

func TestStrategicOverlayRejectsSymlinkAuthority(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated Windows privileges on some runners")
	}
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside-overlay.yaml")
	data := []byte("schema_version: 1\ncoordination_home_id: example-home\nfamilies: []\nthemes: []\nrelations: []\npriorities: []\nanalyses: []\n")
	if err := os.WriteFile(outside, data, 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, filepath.FromSlash(OverlayPath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, path); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadOverlay(root, "example-home"); err == nil || !strings.Contains(err.Error(), "portfolio_overlay_unsafe") {
		t.Fatalf("symlink overlay error = %v", err)
	}
	if actual := mustReadPortfolioFile(t, outside); !bytes.Equal(actual, data) {
		t.Fatal("overlay safety check changed outside bytes")
	}
}

func TestMirrorRootRejectsLocalSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated Windows privileges on some runners")
	}
	repositoryRoot := t.TempDir()
	outside := t.TempDir()
	if err := os.Mkdir(filepath.Join(repositoryRoot, ".codeheart"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(repositoryRoot, ".codeheart", "local")); err != nil {
		t.Fatal(err)
	}
	manager := MirrorManager{
		RepositoryRoot: repositoryRoot,
		Root:           filepath.Join(repositoryRoot, filepath.FromSlash(MirrorRootPath)),
		Runner:         ExecRunner{},
	}
	_, err := manager.Refresh(context.Background(), RepositorySource{Kind: "local-git", Locator: "fixture", CloneURL: filepath.Join(t.TempDir(), "remote.git")})
	if err == nil || !strings.Contains(err.Error(), "mirror_root_unsafe") {
		t.Fatalf("symlink mirror root error = %v", err)
	}
	entries, readErr := os.ReadDir(outside)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("mirror escape wrote outside repository: %+v", entries)
	}
}

func TestMirrorStagingCleanupPreservesSubstitutedBytes(t *testing.T) {
	parent := t.TempDir()
	createPortfolioRepository(t, parent, "staging-cleanup", "member", "staging-cleanup-member", "example-home", "Staging Cleanup Plan", "staging-cleanup-member.discovery.staging-cleanup-plan")
	remote := filepath.Join(parent, "staging-cleanup-remote.git")
	sentinel := []byte("replacement staging authority remains untouched\n")
	scanRoot := t.TempDir()
	mirrorRoot := filepath.Join(scanRoot, filepath.FromSlash(MirrorRootPath))
	manager := MirrorManager{
		RepositoryRoot: scanRoot,
		Root:           mirrorRoot,
		Runner:         ExecRunner{},
		installHook: func(phase, stagingPath string) error {
			if phase != "mirror-before-staging-cleanup" {
				return nil
			}
			if err := os.Rename(stagingPath, stagingPath+".displaced"); err != nil {
				return err
			}
			if err := os.Mkdir(stagingPath, 0o700); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(stagingPath, "sentinel"), sentinel, 0o600)
		},
	}
	repository, err := manager.Refresh(context.Background(), RepositorySource{
		Kind: "local-git", Locator: "staging-cleanup", CloneURL: remote, DefaultBranch: "main", NameHint: "staging-cleanup",
	})
	if repository != nil || err == nil || !strings.Contains(err.Error(), "mirror_staging_cleanup_failed") {
		t.Fatalf("staging substitution repository=%v error=%v", repository, err)
	}
	preserved := false
	entries, readErr := os.ReadDir(mirrorRoot)
	if readErr != nil {
		t.Fatal(readErr)
	}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".cleanup-refresh-") {
			continue
		}
		actual := mustReadPortfolioFile(t, filepath.Join(mirrorRoot, entry.Name(), "sentinel"))
		preserved = bytes.Equal(actual, sentinel)
	}
	if !preserved {
		t.Fatal("identity-safe staging cleanup did not preserve substituted bytes")
	}
}

func TestMirrorRefreshRejectsChildSubstitutionBeforeExternalMutation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated Windows privileges on some runners")
	}
	parent := t.TempDir()
	createPortfolioRepository(t, parent, "substitution", "member", "substitution-member", "example-home", "Substitution Plan", "substitution-member.discovery.substitution-plan")
	remote := filepath.Join(parent, "substitution-remote.git")
	outside := filepath.Join(t.TempDir(), "outside.git")
	runPortfolioGit(t, "", "init", "--bare", outside)
	outsideConfig := mustReadPortfolioFile(t, filepath.Join(outside, "config"))
	runner := &mirrorSubstitutionRunner{delegate: ExecRunner{}, outside: outside}
	scanRoot := t.TempDir()
	manager := MirrorManager{
		RepositoryRoot: scanRoot,
		Root:           filepath.Join(scanRoot, filepath.FromSlash(MirrorRootPath)),
		Runner:         runner,
	}
	_, err := manager.Refresh(context.Background(), RepositorySource{
		Kind: "local-git", Locator: "substitution", CloneURL: remote, DefaultBranch: "main", NameHint: "substitution",
	})
	if err == nil {
		t.Fatal("refresh accepted a substituted mirror child")
	}
	if runner.hookErr != nil {
		t.Fatalf("substitution fixture failed: %v", runner.hookErr)
	}
	if runner.gitRuns != 2 {
		t.Fatalf("Git ran %d time(s); expected initialization plus one descriptor-bound command", runner.gitRuns)
	}
	if actual := mustReadPortfolioFile(t, filepath.Join(outside, "config")); !bytes.Equal(actual, outsideConfig) {
		t.Fatal("refresh mutated the substituted external repository")
	}
}

func TestMirrorInstallRejectsDifferentStagedIdentityAndRestoresPriorMirror(t *testing.T) {
	parent := t.TempDir()
	createPortfolioRepository(t, parent, "install-identity", "member", "install-identity-member", "example-home", "Install Identity Plan", "install-identity-member.discovery.install-identity-plan")
	remote := filepath.Join(parent, "install-identity-remote.git")
	source := RepositorySource{
		Kind: "local-git", Locator: "install-identity", CloneURL: remote, DefaultBranch: "main", NameHint: "install-identity",
	}
	scanRoot := t.TempDir()
	mirrorRoot := filepath.Join(scanRoot, filepath.FromSlash(MirrorRootPath))
	manager := MirrorManager{RepositoryRoot: scanRoot, Root: mirrorRoot, Runner: ExecRunner{}}
	repository, err := manager.Refresh(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	repository.Close()
	finalMirror := filepath.Join(mirrorRoot, mirrorDirectoryName(source))
	priorSnapshot := mustSnapshotMirrorTree(t, finalMirror)

	manager.installHook = func(phase, _ string) error {
		if phase == "mirror-before-retained-authority-check" {
			return fmt.Errorf("injected retained-authority comparison failure")
		}
		return nil
	}
	if _, err := manager.Refresh(context.Background(), source); err == nil || !strings.Contains(err.Error(), "retained refresh authority check failed") {
		t.Fatalf("retained-authority fault error = %v", err)
	}
	if actual := mustSnapshotMirrorTree(t, finalMirror); !bytes.Equal(actual, priorSnapshot) {
		t.Fatal("late retained-authority failure did not restore the full prior mirror tree")
	}

	manager.installHook = func(phase, _ string) error {
		if phase == "mirror-before-prior-cleanup" {
			return fmt.Errorf("injected prior-cleanup failure")
		}
		return nil
	}
	if _, err := manager.Refresh(context.Background(), source); err == nil || !strings.Contains(err.Error(), "prior mirror cleanup precondition failed") {
		t.Fatalf("prior-cleanup fault error = %v", err)
	}
	if actual := mustSnapshotMirrorTree(t, finalMirror); !bytes.Equal(actual, priorSnapshot) {
		t.Fatal("prior-cleanup failure did not restore the full prior mirror tree")
	}

	replacementSentinel := []byte("untrusted replacement remains untouched\n")
	manager.installHook = func(phase, installedPath string) error {
		if phase != "mirror-before-retained-authority-check" {
			return nil
		}
		if err := os.Rename(installedPath, installedPath+".late-displaced"); err != nil {
			return err
		}
		if err := os.Mkdir(installedPath, 0o700); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(installedPath, "sentinel"), replacementSentinel, 0o600)
	}
	if _, err := manager.Refresh(context.Background(), source); err == nil || !strings.Contains(err.Error(), "mirror_rollback_failed") {
		t.Fatalf("late installed identity rollback error = %v", err)
	}
	if actual := mustSnapshotMirrorTree(t, finalMirror); !bytes.Equal(actual, priorSnapshot) {
		t.Fatal("late canonical substitution did not restore the prior mirror")
	}
	preservedReplacement := false
	entries, err := os.ReadDir(mirrorRoot)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".rollback-suspect-") {
			continue
		}
		actual := mustReadPortfolioFile(t, filepath.Join(mirrorRoot, entry.Name(), "sentinel"))
		preservedReplacement = bytes.Equal(actual, replacementSentinel)
	}
	if !preservedReplacement {
		t.Fatal("identity-safe rollback did not preserve the substituted canonical bytes in quarantine")
	}
}

func TestMirrorNoReplaceInstallAndRestorePreserveUnexpectedOccupants(t *testing.T) {
	for _, scenario := range []string{"install", "restore"} {
		t.Run(scenario, func(t *testing.T) {
			parent := t.TempDir()
			name := "no-replace-" + scenario
			createPortfolioRepository(t, parent, name, "member", name+"-member", "example-home", "No Replace Plan", name+"-member.discovery.no-replace-plan")
			remote := filepath.Join(parent, name+"-remote.git")
			source := RepositorySource{Kind: "local-git", Locator: name, CloneURL: remote, DefaultBranch: "main", NameHint: name}
			scanRoot := t.TempDir()
			mirrorRoot := filepath.Join(scanRoot, filepath.FromSlash(MirrorRootPath))
			manager := MirrorManager{RepositoryRoot: scanRoot, Root: mirrorRoot, Runner: ExecRunner{}}
			repository, err := manager.Refresh(context.Background(), source)
			if err != nil {
				t.Fatal(err)
			}
			repository.Close()
			finalMirror := filepath.Join(mirrorRoot, mirrorDirectoryName(source))
			sentinel := []byte("unexpected no-replace occupant remains authoritative\n")
			manager.installHook = func(phase, canonicalPath string) error {
				if scenario == "restore" && phase == "mirror-before-retained-authority-check" {
					return fmt.Errorf("inject rollback before restore")
				}
				wantedPhase := "mirror-before-no-replace-install"
				if scenario == "restore" {
					wantedPhase = "mirror-before-no-replace-restore"
				}
				if phase != wantedPhase {
					return nil
				}
				if err := os.Mkdir(canonicalPath, 0o700); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(canonicalPath, "sentinel"), sentinel, 0o600)
			}
			if _, err := manager.Refresh(context.Background(), source); err == nil || !strings.Contains(err.Error(), "mirror_rollback_failed") {
				t.Fatalf("no-replace %s error = %v", scenario, err)
			}
			if actual := mustReadPortfolioFile(t, filepath.Join(finalMirror, "sentinel")); !bytes.Equal(actual, sentinel) {
				t.Fatalf("no-replace %s overwrote unexpected occupant", scenario)
			}
		})
	}
}

func TestMembershipReadFailureMakesScanIncompleteInsteadOfRetiringMembers(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	now := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})},
		Runner: membershipFailureRunner{delegate: ExecRunner{}}, Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Catalog.Complete {
		t.Fatalf("membership read failure produced a complete catalog: %+v", result.Catalog)
	}
	found := false
	for _, scanError := range result.Catalog.Errors {
		if scanError.Code == "membership_evidence_unavailable" {
			found = true
		}
	}
	if !found {
		t.Fatalf("membership evidence blocker missing: %+v", result.Catalog.Errors)
	}
}

func TestMissingPlanSpecificFreshnessEvidenceMakesScanIncomplete(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	now := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})},
		Runner: freshnessFailureRunner{delegate: ExecRunner{}}, Now: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Catalog.Complete {
		t.Fatalf("missing freshness evidence produced a complete catalog: %+v", result.Catalog)
	}
	found := false
	for _, scanError := range result.Catalog.Errors {
		if scanError.Code == "observation_freshness_unavailable" {
			found = true
		}
	}
	if !found {
		t.Fatalf("freshness blocker missing: %+v", result.Catalog.Errors)
	}
}

func TestCachePublicationRejectsParentSubstitutionAndSymlinkReads(t *testing.T) {
	t.Run("parent substitution", func(t *testing.T) {
		root := t.TempDir()
		oldCatalog := portfolioCatalogFixture("2026-07-31T12:00:00Z")
		if updated, _, err := StoreCompleteCatalog(root, oldCatalog); err != nil || !updated {
			t.Fatalf("initial cache write updated=%v err=%v", updated, err)
		}
		parentPath := filepath.Join(root, filepath.FromSlash(filepath.Dir(CatalogPath)))
		displaced := filepath.Join(root, ".displaced-portfolio-cache")
		sentinel := []byte("replacement authority\n")
		_, _, err := storeCompleteCatalogWithHook(root, portfolioCatalogFixture("2026-07-31T13:00:00Z"), func(phase string) error {
			if phase != "catalog-staged" {
				return nil
			}
			if err := os.Rename(parentPath, displaced); err != nil {
				return err
			}
			if err := os.Mkdir(parentPath, 0o700); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(parentPath, "sentinel"), sentinel, 0o600)
		})
		if err == nil || !strings.Contains(err.Error(), "portfolio_cache_parent_changed") {
			t.Fatalf("parent substitution error = %v", err)
		}
		if actual := mustReadPortfolioFile(t, filepath.Join(parentPath, "sentinel")); !bytes.Equal(actual, sentinel) {
			t.Fatal("cache writer changed replacement-parent bytes")
		}
		if _, err := os.Stat(filepath.Join(parentPath, filepath.Base(CatalogPath))); !os.IsNotExist(err) {
			t.Fatalf("cache writer published through replacement parent: %v", err)
		}
	})

	t.Run("symlink read", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("symlink creation requires elevated Windows privileges on some runners")
		}
		root := t.TempDir()
		outsideRoot := t.TempDir()
		outside := filepath.Join(outsideRoot, "outside-catalog.json")
		data, err := json.Marshal(portfolioCatalogFixture("2026-07-31T12:00:00Z"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(outside, data, 0o600); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, filepath.FromSlash(CatalogPath))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, path); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadCachedCatalog(root); err == nil || !strings.Contains(err.Error(), "not a contained regular file") {
			t.Fatalf("symlink cache read error = %v", err)
		}
	})
}

func TestConcurrentCacheReadersSeeOnlyOldOrNewCompleteBytes(t *testing.T) {
	root := t.TempDir()
	oldCatalog := portfolioCatalogFixture("2026-07-31T12:00:00Z")
	newCatalog := portfolioCatalogFixture("2026-07-31T13:00:00Z")
	if updated, _, err := StoreCompleteCatalog(root, oldCatalog); err != nil || !updated {
		t.Fatalf("initial cache write updated=%v err=%v", updated, err)
	}
	cachePath := filepath.Join(root, filepath.FromSlash(CatalogPath))
	oldBytes := mustReadPortfolioFile(t, cachePath)
	staged := make(chan struct{})
	release := make(chan struct{})
	storeDone := make(chan error, 1)
	go func() {
		_, _, err := storeCompleteCatalogWithHook(root, newCatalog, func(phase string) error {
			if phase == "catalog-staged" {
				close(staged)
				<-release
			}
			return nil
		})
		storeDone <- err
	}()
	<-staged
	if actual := mustReadPortfolioFile(t, cachePath); !bytes.Equal(actual, oldBytes) {
		t.Fatal("staging changed the published cache")
	}
	readerStop := make(chan struct{})
	readerDone := make(chan error, 1)
	go func() {
		for {
			select {
			case <-readerStop:
				readerDone <- nil
				return
			default:
			}
			catalog, err := ReadCachedCatalog(root)
			if err != nil {
				readerDone <- err
				return
			}
			if catalog.CompletedAt != oldCatalog.CompletedAt && catalog.CompletedAt != newCatalog.CompletedAt {
				readerDone <- fmt.Errorf("reader observed neither complete cache")
				return
			}
		}
	}()
	close(release)
	if err := <-storeDone; err != nil {
		t.Fatal(err)
	}
	close(readerStop)
	if err := <-readerDone; err != nil {
		t.Fatalf("concurrent reader observed publication gap: %v", err)
	}
	if actual := mustReadPortfolioFile(t, cachePath); bytes.Equal(actual, oldBytes) {
		t.Fatal("atomic publication did not install the new cache")
	}
}

type membershipFailureRunner struct{ delegate CommandRunner }

func (runner membershipFailureRunner) Run(ctx context.Context, name string, args ...string) (CommandResult, error) {
	if membershipFailure(name, args) {
		return CommandResult{}, fmt.Errorf("injected membership object failure")
	}
	return runner.delegate.Run(ctx, name, args...)
}

func (runner membershipFailureRunner) RunBound(ctx context.Context, directory *os.File, directoryPath, name string, args ...string) (CommandResult, error) {
	if membershipFailure(name, args) {
		return CommandResult{}, fmt.Errorf("injected membership object failure")
	}
	return runner.delegate.(BoundCommandRunner).RunBound(ctx, directory, directoryPath, name, args...)
}

func membershipFailure(name string, args []string) bool {
	if name != "git" || len(args) == 0 || args[len(args)-1] != ConfigPath {
		return false
	}
	for _, arg := range args {
		if arg == "ls-tree" {
			return true
		}
	}
	return false
}

type freshnessFailureRunner struct{ delegate CommandRunner }

func (runner freshnessFailureRunner) Run(ctx context.Context, name string, args ...string) (CommandResult, error) {
	if freshnessFailure(name, args) {
		return CommandResult{}, fmt.Errorf("injected freshness history failure")
	}
	return runner.delegate.Run(ctx, name, args...)
}

func (runner freshnessFailureRunner) RunBound(ctx context.Context, directory *os.File, directoryPath, name string, args ...string) (CommandResult, error) {
	if freshnessFailure(name, args) {
		return CommandResult{}, fmt.Errorf("injected freshness history failure")
	}
	return runner.delegate.(BoundCommandRunner).RunBound(ctx, directory, directoryPath, name, args...)
}

func freshnessFailure(name string, args []string) bool {
	if name == "git" {
		for index := 0; index+2 < len(args); index++ {
			if args[index] == "log" && args[index+1] == "-1" && args[index+2] == "--format=%cI" {
				return true
			}
		}
	}
	return false
}

type fakeGHRunner struct{ pageOneAttempts int }

func (runner *fakeGHRunner) Run(_ context.Context, name string, args ...string) (CommandResult, error) {
	if name != "gh" {
		return CommandResult{}, fmt.Errorf("unexpected command")
	}
	if len(args) >= 2 && args[0] == "auth" {
		return CommandResult{}, nil
	}
	endpoint := args[len(args)-1]
	if strings.Contains(endpoint, "&page=1&type=") {
		runner.pageOneAttempts++
		if runner.pageOneAttempts == 1 {
			return CommandResult{Stderr: []byte("temporary failure super-secret-token")}, fmt.Errorf("exit")
		}
		items := make([]string, 100)
		for index := range items {
			items[index] = fmt.Sprintf(`{"full_name":"example-org/repo-%03d","clone_url":"https://example.invalid/repo-%03d.git","default_branch":"main"}`, index, index)
		}
		return CommandResult{Stdout: []byte("[" + strings.Join(items, ",") + "]")}, nil
	}
	return CommandResult{Stdout: []byte(`[{"full_name":"example-org/repo-100","clone_url":"https://example.invalid/repo-100.git","default_branch":"main"}]`)}, nil
}

type cancelledRunner struct{}

func (cancelledRunner) Run(ctx context.Context, _ string, args ...string) (CommandResult, error) {
	if len(args) >= 2 && args[0] == "auth" {
		return CommandResult{}, nil
	}
	return CommandResult{Stderr: []byte("super-secret-token")}, ctx.Err()
}

type ghFailureRunner struct {
	mode    string
	apiRuns int
}

type prPaginationRunner struct{ apiRuns int }

func (runner *prPaginationRunner) Run(_ context.Context, name string, args ...string) (CommandResult, error) {
	if name != "gh" || len(args) == 0 || args[0] != "api" {
		return CommandResult{}, fmt.Errorf("unexpected PR enrichment command")
	}
	runner.apiRuns++
	endpoint := args[len(args)-1]
	if strings.HasSuffix(endpoint, "&page=1") {
		items := make([]string, githubPageSize)
		for index := range items {
			items[index] = fmt.Sprintf(`{"number":%d,"html_url":"https://github.com/example-org/example-repo/pull/%d","head":{"ref":"feature/unrelated-%d","repo":{"full_name":"example-org/example-repo"}}}`, index+1, index+1, index+1)
		}
		return CommandResult{Stdout: []byte("[" + strings.Join(items, ",") + "]")}, nil
	}
	return CommandResult{Stdout: []byte(`[{"number":101,"html_url":"https://github.com/example-org/example-repo/pull/101","head":{"ref":"feature/changed-plan","repo":{"full_name":"example-org/example-repo"}}},{"number":102,"html_url":"https://github.com/fork-owner/fork-repo/pull/102","head":{"ref":"feature/fork-collision","repo":{"full_name":"fork-owner/fork-repo"}}}]`)}, nil
}

type staticPortfolioSource struct{ repository RepositorySource }

func (source staticPortfolioSource) Kind() string { return source.repository.Kind }

func (source staticPortfolioSource) Discover(context.Context) (DiscoveryResult, error) {
	return DiscoveryResult{Repositories: []RepositorySource{source.repository}, Complete: true, Errors: []ScanError{}}, nil
}

type staticRefEnricher struct {
	pullRequests map[string]string
	apiCalls     int
	seen         map[string]bool
}

func (enricher *staticRefEnricher) EnrichRefs(_ context.Context, _ RepositorySource, refs []string) (RefEnrichment, error) {
	enricher.seen = map[string]bool{}
	for _, ref := range refs {
		enricher.seen[ref] = true
	}
	return RefEnrichment{PullRequests: enricher.pullRequests, APICalls: enricher.apiCalls}, nil
}

func (runner *ghFailureRunner) Run(_ context.Context, _ string, args ...string) (CommandResult, error) {
	if len(args) >= 2 && args[0] == "auth" {
		return CommandResult{}, nil
	}
	runner.apiRuns++
	switch runner.mode {
	case "rate":
		return CommandResult{Stderr: []byte("API rate limit super-secret-token")}, fmt.Errorf("exit")
	case "retry":
		return CommandResult{Stderr: []byte("temporary super-secret-token")}, fmt.Errorf("exit")
	case "permission":
		return CommandResult{Stderr: []byte("HTTP 403 resource not accessible super-secret-token")}, fmt.Errorf("exit")
	case "full-pages":
		item := `{"full_name":"example-org/repository","clone_url":"https://example.invalid/repository.git","default_branch":"main"}`
		items := make([]string, githubPageSize)
		for index := range items {
			items[index] = item
		}
		return CommandResult{Stdout: []byte("[" + strings.Join(items, ",") + "]")}, nil
	default:
		return CommandResult{}, fmt.Errorf("unexpected mode")
	}
}

type portfolioGitFixture struct {
	parent string
	home   string
	member string
}

func newPortfolioGitFixture(t *testing.T) portfolioGitFixture {
	t.Helper()
	parent := t.TempDir()
	home := createPortfolioRepository(t, parent, "home", "coordination-home", "example-home", "example-coordination", "Home Plan", "example-home.discovery.home-plan")
	member := createPortfolioRepository(t, parent, "member", "member", "example-member", "example-coordination", "Member Plan", "example-member.discovery.member-plan")
	unrelated := createPortfolioRepository(t, parent, "unrelated", "member", "unrelated-member", "different-home", "Unrelated Plan", "unrelated-member.discovery.unrelated-plan")
	runPortfolioGit(t, unrelated, "checkout", "-b", "feature/branch-only-enrollment")
	writePortfolioFile(t, unrelated, ConfigPath, []byte(portfolioConfig("member", "unrelated-member", "example-coordination")))
	runPortfolioGit(t, unrelated, "add", ConfigPath)
	runPortfolioGit(t, unrelated, "commit", "-m", "Attempt branch-only portfolio enrollment")
	runPortfolioGit(t, unrelated, "push", "origin", "feature/branch-only-enrollment")
	runPortfolioGit(t, unrelated, "checkout", "main")
	runPortfolioGit(t, member, "checkout", "-b", "feature/changed-plan")
	writePortfolioPlan(t, member, "member-plan/member-plan_discovery_doc.md", "Member Plan On Branch", "example-member.discovery.member-plan", "branch purpose")
	runPortfolioGitWithDate(t, member, "2026-06-01T12:00:00Z", "add", "docs/repo/plans")
	runPortfolioGitWithDate(t, member, "2026-06-01T12:00:00Z", "commit", "-m", "Change member plan on remote branch")
	runPortfolioGit(t, member, "push", "origin", "feature/changed-plan")
	runPortfolioGit(t, member, "checkout", "main")
	runPortfolioGit(t, member, "checkout", "-b", "local-only")
	writePortfolioPlan(t, member, "local-only/local-only_discovery_doc.md", "Local Only", "example-member.discovery.local-only", "not pushed")
	runPortfolioGit(t, member, "add", "docs/repo/plans")
	runPortfolioGit(t, member, "commit", "-m", "Unpushed local plan")
	runPortfolioGit(t, member, "checkout", "main")
	overlayPath := filepath.Join(home, filepath.FromSlash(OverlayPath))
	if err := os.MkdirAll(filepath.Dir(overlayPath), 0o755); err != nil {
		t.Fatal(err)
	}
	overlay := "schema_version: 1\ncoordination_home_id: example-coordination\nfamilies: []\nthemes: []\nrelations: []\npriorities: []\nanalyses: []\n"
	if err := os.WriteFile(overlayPath, []byte(overlay), 0o644); err != nil {
		t.Fatal(err)
	}
	return portfolioGitFixture{parent: parent, home: home, member: member}
}

func createPortfolioRepository(t *testing.T, parent, name, role, repositoryID, homeID, title, planID string) string {
	t.Helper()
	work := filepath.Join(parent, name)
	remote := filepath.Join(parent, name+"-remote.git")
	runPortfolioGit(t, "", "init", "--bare", remote)
	runPortfolioGit(t, "", "init", "-b", "main", work)
	runPortfolioGit(t, work, "config", "user.email", "test@example.invalid")
	runPortfolioGit(t, work, "config", "user.name", "Portfolio Test")
	writePortfolioFile(t, work, ConfigPath, []byte(portfolioConfig(role, repositoryID, homeID)))
	writePortfolioFile(t, work, LockPath, []byte(testLockYAML))
	writePortfolioFile(t, work, KitMarkerPath, []byte("# Installed Kit\n"))
	writePortfolioPlan(t, work, name+"-plan/"+name+"-plan_discovery_doc.md", title, planID, "baseline purpose")
	runPortfolioGit(t, work, "add", ".")
	runPortfolioGit(t, work, "commit", "-m", "Initial remote baseline")
	runPortfolioGit(t, work, "remote", "add", "origin", remote)
	runPortfolioGit(t, work, "push", "-u", "origin", "main")
	runPortfolioGit(t, remote, "symbolic-ref", "HEAD", "refs/heads/main")
	return work
}

func portfolioConfig(role, repositoryID, homeID string) string {
	return fmt.Sprintf("schema_version: 1\nselected_profile: standard\nproject_display_name: %s\nselected_setup_folder: .\nlocal_consumer_layer:\n  repo_docs_path: docs/repo/\n  agent_memory_path: docs/agent-memory/\n  user_layer_path: .codeheart/user/\n  local_machine_layer_path: .codeheart/local/\ncomponent_settings:\n  planning-workflows:\n    plan_catalog_mode: canonical\n    plan_catalog_discovery_version: 2\nportfolio:\n  schema_version: 2\n  role: %s\n  member_repository_id: %s\n  coordination_home_id: %s\n", repositoryID, role, repositoryID, homeID)
}

const testLockYAML = `schema_version: 1
kit_version: 0.1.23
selected_profile: standard
selected_components: []
release:
  asset_url: local-test
  checksum_sha256: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
managed_paths: []
generated_surfaces: []
cli_repair:
  installed_cli_path: codeheart-operating-kit
  repair_source_url: local-test
update_check:
  last_update_check_at: "2026-07-01T00:00:00Z"
  next_update_check_due: "2026-08-01T00:00:00Z"
  latest_seen_version: 0.1.23
  update_status: current
native_capabilities: {}
`

const testKitSourceMarkerYAML = `schema_version: 1
version: 0.1.23
compatibility:
  lock_schema_versions: [1]
  config_schema_versions: [1]
  operation_result_schema_version: 1
  platforms: [macos-universal, windows-x64]
  commands: [init, repair, sync, update-check, upgrade, check]
components:
  - id: planning-workflows
    version: 0.1.23
    manifest_path: components/planning-workflows/component.yaml
    checksum_sha256: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
    consumer_impact: []
profiles:
  - id: standard
    version: 0.1.23
    manifest_path: profiles/standard.yaml
    checksum_sha256: bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
    graph_sha256: cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc
    components: [planning-workflows]
consumer_impact: []
`

func writePortfolioPlan(t *testing.T, root, relative, title, id, purpose string) {
	t.Helper()
	writePortfolioPlanAt(t, root, filepath.ToSlash(filepath.Join("docs/repo/plans", relative)), title, id, plancatalog.KindDiscovery, purpose)
}

func writePortfolioPlanAt(t *testing.T, root, relative, title, id string, kind plancatalog.Kind, purpose string) {
	t.Helper()
	content := fmt.Sprintf("Last updated: 2026-06-01T12:00:00Z (UTC)\nCreated: 2026-06-01\nStatus: active\n\n# %s\n<!-- BEGIN CODEHEART PLAN METADATA -->\n```yaml\nplan:\n  schema_version: 1\n  id: %s\n  kind: %s\n  purpose: %s\n  first_cataloged: 2026-06-01T12:00:00Z\n  catalog_metadata_updated: 2026-06-01T12:00:00Z\n```\n<!-- END CODEHEART PLAN METADATA -->\n\n## Scope\n\nTest plan.\n", title, id, kind, purpose)
	writePortfolioFile(t, root, relative, []byte(content))
}

func writePortfolioFile(t *testing.T, root, relative string, data []byte) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func runPortfolioGit(t *testing.T, root string, args ...string) {
	t.Helper()
	commandArgs := append([]string{}, args...)
	if root != "" {
		commandArgs = append([]string{"-C", root}, commandArgs...)
	}
	command := exec.Command("git", commandArgs...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(commandArgs, " "), err, output)
	}
}

func runPortfolioGitWithDate(t *testing.T, root, date string, args ...string) {
	t.Helper()
	commandArgs := append([]string{"-C", root}, args...)
	command := exec.Command("git", commandArgs...)
	command.Env = append(os.Environ(), "GIT_AUTHOR_DATE="+date, "GIT_COMMITTER_DATE="+date)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(commandArgs, " "), err, output)
	}
}

func portfolioGitOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	commandArgs := append([]string{"-C", root}, args...)
	output, err := exec.Command("git", commandArgs...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(commandArgs, " "), err, output)
	}
	return string(output)
}

func mustReadPortfolioFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func mustSnapshotMirrorTree(t *testing.T, path string) []byte {
	t.Helper()
	root, err := os.OpenRoot(path)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	snapshot, err := snapshotMirrorTree(root)
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func commitAndPushPortfolioFile(t *testing.T, repository, relative, message string) {
	t.Helper()
	runPortfolioGit(t, repository, "add", relative)
	runPortfolioGit(t, repository, "commit", "-m", message)
	runPortfolioGit(t, repository, "push", "origin", "main")
}

func portfolioCatalogFixture(completedAt string) Catalog {
	emptyDigest := aggregateMemberDigest(nil, func(plancatalog.CatalogMember) string { return "" })
	return Catalog{
		SchemaVersion: 2, DiscoveryVersion: plancatalog.DiscoveryV2, PolicyDigest: emptyDigest, CandidateSetDigest: emptyDigest, CoordinationHomeID: "example-home",
		StartedAt: completedAt, CompletedAt: completedAt, LastCompleteScanAt: completedAt, Complete: true,
		Members: []plancatalog.CatalogMember{}, Observations: []plancatalog.SourceObservation{}, Candidates: []Candidate{}, Errors: []ScanError{},
		Metrics: Metrics{MaxConcurrency: 1},
	}
}

func TestMirrorNamesArePortableAndStable(t *testing.T) {
	sources := []RepositorySource{{Kind: "github", Locator: "Org/Repo", CloneURL: "https://example.invalid/repo.git", NameHint: "Repo Name"}, {Kind: "local-git", Locator: "/tmp/repo", CloneURL: "/tmp/remote.git", NameHint: "Repo Name"}}
	names := []string{mirrorDirectoryName(sources[0]), mirrorDirectoryName(sources[1])}
	sort.Strings(names)
	for _, name := range names {
		if strings.ContainsAny(name, `/\\ :`) || !strings.HasSuffix(name, ".git") {
			t.Fatalf("non-portable mirror name %q on %s", name, runtime.GOOS)
		}
	}
	if names[0] == names[1] {
		t.Fatal("distinct sources received the same mirror name")
	}
	for _, remote := range []string{"https://token@example.invalid/repository.git", "https://example.invalid/repository.git?token=secret", "token@example.invalid:repository.git"} {
		if err := validateRemoteURL(remote); err == nil {
			t.Fatalf("credential-bearing remote was accepted: %s", remote)
		}
	}
	for _, remote := range []string{"https://example.invalid/repository.git", "ssh://git@example.invalid/repository.git", "git@example.invalid:repository.git"} {
		if err := validateRemoteURL(remote); err != nil {
			t.Fatalf("ordinary remote %s rejected: %v", remote, err)
		}
	}
}

func TestGitTransportPolicyRejectsRemoteHelpersAndPinsSupportedProtocols(t *testing.T) {
	for _, remote := range []string{
		"ext::sh -c touch-owned",
		"custom::repository",
		"git://example.invalid/repository.git",
		"http://example.invalid/repository.git",
		"helper+ssh://example.invalid/repository.git",
	} {
		if err := validateRemoteURL(remote); err == nil || !strings.Contains(err.Error(), "remote_url_protocol_forbidden") {
			t.Fatalf("unsafe remote %q error = %v", remote, err)
		}
	}
	for _, remote := range []string{
		"https://example.invalid/repository.git",
		"ssh://git@example.invalid/repository.git",
		"git@example.invalid:repository.git",
		"file:///tmp/repository.git",
		filepath.Join(t.TempDir(), "repository.git"),
	} {
		if err := validateRemoteURL(remote); err != nil {
			t.Fatalf("supported remote %q rejected: %v", remote, err)
		}
	}
	remote := "https://example.invalid/repository.git"
	arguments := safeGitRemoteArgs(remote, "fetch", "origin")
	joined := strings.Join(arguments, "\x00")
	for _, expected := range []string{
		"protocol.allow=never",
		"protocol.https.allow=always",
		"protocol.ssh.allow=always",
		"protocol.file.allow=always",
		"protocol.ext.allow=never",
		"core.sshCommand=ssh",
		"core.hooksPath=" + os.DevNull,
		"url." + remote + ".insteadOf=" + remote,
		"remote.origin.uploadpack=git-upload-pack",
	} {
		if !strings.Contains(joined, expected) {
			t.Fatalf("safe Git arguments omit %q: %+v", expected, arguments)
		}
	}
	environment := sanitizedGitEnvironment([]string{
		"PATH=/usr/bin", "GIT_ALLOW_PROTOCOL=ext", "GIT_CONFIG_COUNT=1",
		"GIT_CONFIG_KEY_0=protocol.ext.allow", "GIT_CONFIG_VALUE_0=always", "GIT_CONFIG_PARAMETERS='protocol.ext.allow=always'", "GIT_EXEC_PATH=/unsafe", "GIT_SSH_COMMAND=unsafe-command",
		"git_config_parameters='remote.origin.uploadpack=unsafe'", "Git_Exec_Path=/mixed-case-unsafe", "git_ssh_command=mixed-case-unsafe",
	})
	if strings.Join(environment, "\n") != "PATH=/usr/bin\nGIT_TERMINAL_PROMPT=0\nGIT_LITERAL_PATHSPECS=1\nGIT_NO_LAZY_FETCH=1\nGIT_NO_REPLACE_OBJECTS=1\nLC_ALL=C" {
		t.Fatalf("unsafe Git environment survived sanitization: %+v", environment)
	}
	if runtime.GOOS != "windows" {
		home := t.TempDir()
		sentinel := filepath.Join(t.TempDir(), "unsafe-helper-ran")
		helper := filepath.Join(t.TempDir(), "unsafe-helper.sh")
		script := []byte("#!/bin/sh\nprintf unsafe > \"" + sentinel + "\"\nexit 1\n")
		if err := os.WriteFile(helper, script, 0o700); err != nil {
			t.Fatal(err)
		}
		globalConfig := fmt.Sprintf("[url \"ext::%s\"]\n\tinsteadOf = %s\n[protocol \"ext\"]\n\tallow = always\n", helper, remote)
		if err := os.WriteFile(filepath.Join(home, ".gitconfig"), []byte(globalConfig), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("HOME", home)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = ExecRunner{}.Run(ctx, "git", safeGitRemoteArgs(remote, "ls-remote", remote)...)
		if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
			t.Fatalf("global URL rewrite executed forbidden remote helper: %v", err)
		}

		parent := t.TempDir()
		createPortfolioRepository(t, parent, "transport", "member", "transport-member", "example-home", "Transport Plan", "transport-member.discovery.transport-plan")
		uploadPackSentinel := filepath.Join(t.TempDir(), "uploadpack-helper-ran")
		uploadPackHelper := filepath.Join(t.TempDir(), "uploadpack-helper.sh")
		uploadPackScript := []byte("#!/bin/sh\nprintf unsafe > \"" + uploadPackSentinel + "\"\nexit 1\n")
		if err := os.WriteFile(uploadPackHelper, uploadPackScript, 0o700); err != nil {
			t.Fatal(err)
		}
		t.Setenv("GIT_CONFIG_PARAMETERS", "'remote.origin.uploadpack="+uploadPackHelper+"'")
		injected, err := exec.Command("git", "config", "--get", "remote.origin.uploadpack").CombinedOutput()
		if err != nil || strings.TrimSpace(string(injected)) != uploadPackHelper {
			t.Fatalf("uploadpack injection fixture was not active: %v %q", err, injected)
		}
		scanRoot := t.TempDir()
		manager := MirrorManager{RepositoryRoot: scanRoot, Root: filepath.Join(scanRoot, filepath.FromSlash(MirrorRootPath)), Runner: ExecRunner{}}
		repository, err := manager.Refresh(context.Background(), RepositorySource{
			Kind: "local-git", Locator: "transport", CloneURL: filepath.Join(parent, "transport-remote.git"), DefaultBranch: "main", NameHint: "transport",
		})
		if err != nil {
			t.Fatalf("safe fetch with injected uploadpack failed: %v", err)
		}
		repository.Close()
		if _, err := os.Stat(uploadPackSentinel); !os.IsNotExist(err) {
			t.Fatalf("command-scope uploadpack injection executed helper: %v", err)
		}
	} else {
		parent := t.TempDir()
		createPortfolioRepository(t, parent, "windows-transport", "member", "windows-transport-member", "example-home", "Windows Transport Plan", "windows-transport-member.discovery.windows-transport-plan")
		uploadPackSentinel := filepath.Join(t.TempDir(), "mixed-case-uploadpack-helper-ran")
		uploadPackHelper := filepath.Join(t.TempDir(), "mixed-case-uploadpack-helper.cmd")
		uploadPackScript := []byte("@echo off\r\necho unsafe>\"" + uploadPackSentinel + "\"\r\nexit /b 1\r\n")
		if err := os.WriteFile(uploadPackHelper, uploadPackScript, 0o700); err != nil {
			t.Fatal(err)
		}
		t.Setenv("git_config_parameters", "'remote.origin.uploadpack="+uploadPackHelper+"'")
		injected, err := exec.Command("git", "config", "--get", "remote.origin.uploadpack").CombinedOutput()
		if err != nil || !strings.EqualFold(strings.TrimSpace(string(injected)), uploadPackHelper) {
			t.Fatalf("mixed-case Windows uploadpack injection fixture was not active: %v %q", err, injected)
		}
		scanRoot := t.TempDir()
		manager := MirrorManager{RepositoryRoot: scanRoot, Root: filepath.Join(scanRoot, filepath.FromSlash(MirrorRootPath)), Runner: ExecRunner{}}
		repository, err := manager.Refresh(context.Background(), RepositorySource{
			Kind: "local-git", Locator: "windows-transport", CloneURL: filepath.Join(parent, "windows-transport-remote.git"), DefaultBranch: "main", NameHint: "windows-transport",
		})
		if err != nil {
			t.Fatalf("safe Windows fetch with mixed-case injection failed: %v", err)
		}
		repository.Close()
		if _, err := os.Stat(uploadPackSentinel); !os.IsNotExist(err) {
			t.Fatalf("mixed-case Windows uploadpack injection executed helper: %v", err)
		}
	}
}

func TestEveryScannerOwnedGitCommandUsesCentralSafetyPolicy(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	runner := &policyRecordingRunner{delegate: ExecRunner{}}
	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, runner)}, Runner: runner,
		Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC),
	})
	if err != nil || !result.Catalog.Complete {
		t.Fatalf("recorded policy scan complete=%v err=%v errors=%+v", result.Catalog.Complete, err, result.Catalog.Errors)
	}
	count, batches, failures := runner.snapshot()
	if count == 0 || batches == 0 || len(failures) != 0 {
		t.Fatalf("scanner Git policy count=%d batches=%d failures=%+v", count, batches, failures)
	}
}

func TestScannerTreatsAdversarialBranchContentAsInertData(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	sentinel := filepath.Join(t.TempDir(), "branch-content-executed")
	runPortfolioGit(t, fixture.member, "checkout", "-b", "feature/adversarial-content", "main")

	executable := filepath.Join(fixture.member, "docs", "repo", "plans", "adversarial", "run-me.sh")
	writePortfolioFile(t, fixture.member, "docs/repo/plans/adversarial/run-me.sh", []byte("#!/bin/sh\nprintf executed > \""+sentinel+"\"\n"))
	if err := os.Chmod(executable, 0o755); err != nil {
		t.Fatal(err)
	}
	writePortfolioFile(t, fixture.member, ".githooks/post-checkout", []byte("#!/bin/sh\nprintf hook > \""+sentinel+"\"\n"))
	if err := os.Chmod(filepath.Join(fixture.member, ".githooks", "post-checkout"), 0o755); err != nil {
		t.Fatal(err)
	}
	writePortfolioFile(t, fixture.member, "docs/repo/plans/adversarial/adversarial_discovery_doc.md", []byte("Last updated: 2026-07-31T12:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: active\n\n# Adversarial\n<!-- BEGIN CODEHEART PLAN METADATA -->\n```yaml\nplan: [\n"))
	oversizedPath := "docs/repo/plans/adversarial/oversized_implementation_doc.md"
	writePortfolioFile(t, fixture.member, oversizedPath, []byte(strings.Repeat("x", plancatalog.MaxPlanSourceBytes+1)))
	writePortfolioFile(t, fixture.member, "docs/repo/plans/adversarial/secret-placeholder.txt", []byte("Authorization: Bearer EXAMPLE_TOKEN_PLACEHOLDER\n"))
	runPortfolioGit(t, fixture.member, "add", "docs/repo/plans/adversarial", ".githooks/post-checkout")
	runPortfolioGit(t, fixture.member, "commit", "-m", "Add adversarial inert branch content")
	runPortfolioGit(t, fixture.member, "push", "origin", "feature/adversarial-content")
	runPortfolioGit(t, fixture.member, "checkout", "main")

	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{},
		Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sentinel); !os.IsNotExist(err) {
		t.Fatalf("scanner executed adversarial branch content: %v", err)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(encoded, []byte("EXAMPLE_TOKEN_PLACEHOLDER")) || bytes.Contains(encoded, []byte(sentinel)) {
		t.Fatalf("scanner captured ignored adversarial bytes: %s", encoded)
	}
	if result.Catalog.Complete {
		t.Fatal("malformed adversarial formal plans unexpectedly produced a complete scan")
	}
	oversizedBlocked := false
	for _, scanError := range result.Catalog.Errors {
		if scanError.Path == oversizedPath && scanError.Code == "plan_source_unsafe" {
			oversizedBlocked = true
		}
	}
	if !oversizedBlocked {
		t.Fatalf("oversized remote blob did not produce bounded-read evidence: %+v", result.Catalog.Errors)
	}
}

func TestBranchOverlaySuppressesPureRenameAndRetainsChangedRename(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	original := "docs/repo/plans/member-plan/member-plan_discovery_doc.md"

	runPortfolioGit(t, fixture.member, "checkout", "-b", "feature/pure-rename", "main")
	pureDestination := "docs/repo/plans/pure-renamed/pure-renamed_discovery_doc.md"
	if err := os.MkdirAll(filepath.Join(fixture.member, filepath.Dir(pureDestination)), 0o755); err != nil {
		t.Fatal(err)
	}
	runPortfolioGit(t, fixture.member, "mv", original, pureDestination)
	runPortfolioGit(t, fixture.member, "commit", "-m", "Purely rename member plan")
	runPortfolioGit(t, fixture.member, "push", "origin", "feature/pure-rename")
	runPortfolioGit(t, fixture.member, "checkout", "main")

	if runtime.GOOS != "windows" {
		runPortfolioGit(t, fixture.member, "checkout", "-b", "feature/mode-only", "main")
		if err := os.Chmod(filepath.Join(fixture.member, filepath.FromSlash(original)), 0o755); err != nil {
			t.Fatal(err)
		}
		runPortfolioGit(t, fixture.member, "add", original)
		runPortfolioGit(t, fixture.member, "commit", "-m", "Change only member plan mode")
		runPortfolioGit(t, fixture.member, "push", "origin", "feature/mode-only")
		runPortfolioGit(t, fixture.member, "checkout", "main")
	}

	runPortfolioGit(t, fixture.member, "checkout", "-b", "feature/changed-rename", "main")
	changedDestination := "docs/repo/plans/changed-renamed/changed-renamed_discovery_doc.md"
	if err := os.MkdirAll(filepath.Join(fixture.member, filepath.Dir(changedDestination)), 0o755); err != nil {
		t.Fatal(err)
	}
	runPortfolioGit(t, fixture.member, "mv", original, changedDestination)
	writePortfolioPlan(t, fixture.member, "changed-renamed/changed-renamed_discovery_doc.md", "Changed Renamed Member Plan", "example-member.discovery.member-plan", "renamed plan with changed canonical bytes")
	runPortfolioGit(t, fixture.member, "add", "docs/repo/plans")
	runPortfolioGit(t, fixture.member, "commit", "-m", "Rename and change member plan")
	runPortfolioGit(t, fixture.member, "push", "origin", "feature/changed-rename")
	runPortfolioGit(t, fixture.member, "checkout", "main")

	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{
		Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{},
		Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Catalog.Complete {
		t.Fatalf("rename scan incomplete: %+v", result.Catalog.Errors)
	}
	pureCount := 0
	changedCount := 0
	modeCount := 0
	for _, observation := range result.Catalog.Observations {
		switch observation.Ref {
		case "refs/remotes/origin/feature/pure-rename":
			pureCount++
		case "refs/remotes/origin/feature/changed-rename":
			changedCount++
			if observation.CanonicalPath != changedDestination {
				t.Fatalf("changed rename path = %q", observation.CanonicalPath)
			}
		case "refs/remotes/origin/feature/mode-only":
			modeCount++
		}
	}
	if pureCount != 0 || modeCount != 0 || changedCount != 1 {
		t.Fatalf("rename overlays pure=%d mode=%d changed=%d observations=%+v", pureCount, modeCount, changedCount, result.Catalog.Observations)
	}
}

func TestBranchOverlayUsesDefaultPolicyAndDoesNotInferFamilyAuthority(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	defaultConfig := strings.Replace(portfolioConfig("member", "example-member", "example-coordination"), "    plan_catalog_discovery_version: 2\n", "    plan_catalog_discovery_version: 2\n    plan_catalog_ownership:\n      excluded_roots:\n        - vendor/\n", 1)
	writePortfolioFile(t, fixture.member, ConfigPath, []byte(defaultConfig))
	writePortfolioFile(t, fixture.member, "products/widget/docs/router/README.md", []byte("# Domain router\n\nNo plan metadata.\n"))
	movable := "products/widget/docs/movable/movable_discovery_doc.md"
	movedExcluded := "vendor/moved/docs/plans/movable_discovery_doc.md"
	writePortfolioPlanAt(t, fixture.member, movable, "Movable", "example-member.discovery.movable", plancatalog.KindDiscovery, "ownership rename fixture")
	runPortfolioGit(t, fixture.member, "add", ConfigPath, "products/widget/docs/router/README.md", movable)
	runPortfolioGit(t, fixture.member, "commit", "-m", "Define default branch exclusions and router")
	runPortfolioGit(t, fixture.member, "push", "origin", "main")

	runPortfolioGit(t, fixture.member, "checkout", "-b", "feature/cannot-broaden-policy", "main")
	branchConfig := strings.Replace(defaultConfig, "    plan_catalog_ownership:\n      excluded_roots:\n        - vendor/\n", "    plan_catalog_ownership:\n      excluded_roots: []\n", 1)
	writePortfolioFile(t, fixture.member, ConfigPath, []byte(branchConfig))
	writePortfolioPlanAt(t, fixture.member, "vendor/embedded/docs/plans/embedded_discovery_doc.md", "Embedded", "example-member.discovery.embedded", plancatalog.KindDiscovery, "must remain excluded")
	writePortfolioPlanAt(t, fixture.member, "products/widget/docs/router/child/child_discovery_doc.md", "Router Child", "example-member.discovery.router-child", plancatalog.KindDiscovery, "ordinary child")
	if err := os.MkdirAll(filepath.Join(fixture.member, filepath.Dir(movedExcluded)), 0o755); err != nil {
		t.Fatal(err)
	}
	runPortfolioGit(t, fixture.member, "mv", movable, movedExcluded)
	runPortfolioGit(t, fixture.member, "add", ConfigPath, "vendor", "products/widget/docs/router/child")
	runPortfolioGit(t, fixture.member, "commit", "-m", "Attempt branch policy broadening")
	runPortfolioGit(t, fixture.member, "push", "origin", "feature/cannot-broaden-policy")
	runPortfolioGit(t, fixture.member, "checkout", "main")

	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{}, Now: time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Catalog.Complete {
		t.Fatalf("reviewed exclusion made scan incomplete: %+v", result.Catalog.Errors)
	}
	childFound := false
	excludedRenameFound := false
	for _, observation := range result.Catalog.Observations {
		if observation.CanonicalPath == "vendor/embedded/docs/plans/embedded_discovery_doc.md" {
			t.Fatal("feature branch broadened default-branch ownership authority")
		}
		if observation.CanonicalPath == "products/widget/docs/router/README.md" {
			t.Fatal("unchanged ordinary README gained family authority from a new child")
		}
		if observation.CanonicalPath == "products/widget/docs/router/child/child_discovery_doc.md" && observation.Visibility == "unmerged-branch" {
			childFound = true
		}
	}
	for _, candidate := range result.Catalog.PlanCandidates {
		if candidate.RepositoryID == "example-member" && candidate.Ref == "refs/remotes/origin/feature/cannot-broaden-policy" && candidate.Path == movedExcluded && candidate.Ownership == plancatalog.OwnershipExcluded && candidate.ExclusionRoot == "vendor/" && candidate.ContentSHA256 != "" && candidate.Commit != "" {
			excludedRenameFound = true
		}
	}
	if !childFound {
		t.Fatalf("valid branch child observation missing: %+v", result.Catalog.Observations)
	}
	if !excludedRenameFound {
		t.Fatalf("ownership-changing excluded rename evidence missing: %+v", result.Catalog.PlanCandidates)
	}
}

func TestBranchOverlayKeepsEligibilityAndKindChangingRenamesVisible(t *testing.T) {
	fixture := newPortfolioGitFixture(t)
	original := "docs/repo/plans/member-plan/member-plan_discovery_doc.md"

	runPortfolioGit(t, fixture.member, "checkout", "-b", "feature/kind-changing-rename", "main")
	kindDestination := "docs/repo/plans/member-plan/member-plan_implementation_doc.md"
	runPortfolioGit(t, fixture.member, "mv", original, kindDestination)
	runPortfolioGit(t, fixture.member, "commit", "-m", "Rename plan across filename kinds")
	runPortfolioGit(t, fixture.member, "push", "origin", "feature/kind-changing-rename")
	runPortfolioGit(t, fixture.member, "checkout", "main")

	outside := "archive/incoming_discovery_doc.md"
	inside := "products/widget/source/deep/docs/plans/incoming_discovery_doc.md"
	writePortfolioPlanAt(t, fixture.member, outside, "Incoming", "example-member.discovery.incoming", plancatalog.KindDiscovery, "outside before rename")
	runPortfolioGit(t, fixture.member, "add", outside)
	runPortfolioGit(t, fixture.member, "commit", "-m", "Add candidate outside docs on default branch")
	runPortfolioGit(t, fixture.member, "push", "origin", "main")
	runPortfolioGit(t, fixture.member, "checkout", "-b", "feature/rename-into-docs", "main")
	if err := os.MkdirAll(filepath.Join(fixture.member, filepath.Dir(inside)), 0o755); err != nil {
		t.Fatal(err)
	}
	runPortfolioGit(t, fixture.member, "mv", outside, inside)
	runPortfolioGit(t, fixture.member, "commit", "-m", "Move candidate into deep docs")
	runPortfolioGit(t, fixture.member, "push", "origin", "feature/rename-into-docs")
	runPortfolioGit(t, fixture.member, "checkout", "main")

	config := Config{SchemaVersion: 2, Role: RoleCoordinationHome, RepositoryID: "example-home", CoordinationHomeID: "example-coordination", Root: fixture.home}
	result, err := Scan(context.Background(), ScanOptions{Root: fixture.home, Config: config, Sources: []Source{NewLocalGitSource(fixture.parent, ExecRunner{})}, Runner: ExecRunner{}, Now: time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatal(err)
	}
	kindMismatch := false
	for _, scanError := range result.Catalog.Errors {
		if scanError.Ref == "refs/remotes/origin/feature/kind-changing-rename" && scanError.Code == "record_kind_path_mismatch" {
			kindMismatch = true
		}
	}
	if !kindMismatch {
		t.Fatalf("same-byte kind-changing rename was suppressed: %+v", result.Catalog.Errors)
	}
	insideFound := false
	for _, observation := range result.Catalog.Observations {
		if observation.Ref == "refs/remotes/origin/feature/rename-into-docs" && observation.CanonicalPath == inside {
			insideFound = true
		}
	}
	if !insideFound {
		t.Fatalf("eligibility-changing rename was not observed: %+v", result.Catalog.Observations)
	}
}

type policyRecordingRunner struct {
	mu       sync.Mutex
	delegate CommandRunner
	count    int
	batches  int
	failures []string
}

type mirrorSubstitutionRunner struct {
	delegate CommandRunner
	outside  string
	gitRuns  int
	hookErr  error
}

func (runner *mirrorSubstitutionRunner) Run(ctx context.Context, name string, args ...string) (CommandResult, error) {
	return runner.delegate.Run(ctx, name, args...)
}

func (runner *mirrorSubstitutionRunner) RunBound(ctx context.Context, directory *os.File, directoryPath, name string, args ...string) (CommandResult, error) {
	if name == "git" {
		runner.gitRuns++
	}
	if name == "git" && runner.gitRuns == 2 {
		displaced := directoryPath + ".displaced"
		if renameErr := os.Rename(directoryPath, displaced); renameErr != nil {
			runner.hookErr = renameErr
			return CommandResult{}, renameErr
		}
		if symlinkErr := os.Symlink(runner.outside, directoryPath); symlinkErr != nil {
			runner.hookErr = symlinkErr
			return CommandResult{}, symlinkErr
		}
	}
	// The retained descriptor, not directoryPath, must select the Git repository. The
	// post-command namespace check then rejects the substitution before refresh continues.
	bound := runner.delegate.(BoundCommandRunner)
	return bound.RunBound(ctx, directory, directoryPath, name, args...)
}

func (runner *policyRecordingRunner) Run(ctx context.Context, name string, args ...string) (CommandResult, error) {
	runner.record(name, args)
	return runner.delegate.Run(ctx, name, args...)
}

func (runner *policyRecordingRunner) RunBound(ctx context.Context, directory *os.File, directoryPath, name string, args ...string) (CommandResult, error) {
	runner.record(name, args)
	return runner.delegate.(BoundCommandRunner).RunBound(ctx, directory, directoryPath, name, args...)
}

func (runner *policyRecordingRunner) StartBound(ctx context.Context, directory *os.File, directoryPath, name string, args ...string) (*BoundCommandStream, error) {
	runner.record(name, args)
	return runner.delegate.(BoundStreamingRunner).StartBound(ctx, directory, directoryPath, name, args...)
}

func (runner *policyRecordingRunner) record(name string, args []string) {
	if name == "git" {
		joined := strings.Join(args, "\x00")
		missing := []string{}
		for _, expected := range []string{
			"protocol.allow=never", "protocol.https.allow=always", "protocol.ssh.allow=always",
			"protocol.file.allow=always", "protocol.ext.allow=never", "core.hooksPath=" + os.DevNull,
		} {
			if !strings.Contains(joined, expected) {
				missing = append(missing, expected)
			}
		}
		remoteOperation := false
		for _, argument := range args {
			if argument == "fetch" || argument == "ls-remote" {
				remoteOperation = true
			}
		}
		if remoteOperation && (!strings.Contains(joined, ".insteadOf=") || !strings.Contains(joined, "remote.origin.uploadpack=git-upload-pack")) {
			missing = append(missing, "remote URL/uploadpack pin")
		}
		runner.mu.Lock()
		runner.count++
		for index, argument := range args {
			if argument == "cat-file" && index+1 < len(args) && args[index+1] == "--batch" {
				runner.batches++
			}
		}
		if len(missing) != 0 {
			runner.failures = append(runner.failures, fmt.Sprintf("%v missing %v", args, missing))
		}
		runner.mu.Unlock()
	}
}

func (runner *policyRecordingRunner) snapshot() (int, int, []string) {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return runner.count, runner.batches, append([]string{}, runner.failures...)
}
