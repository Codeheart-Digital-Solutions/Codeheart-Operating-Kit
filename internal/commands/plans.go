package commands

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/plancatalog"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/portfolio"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type plansValidationOutput struct {
	SchemaVersion              int                                   `json:"schema_version"`
	Mode                       plancatalog.CatalogMode               `json:"mode"`
	RepositoryID               string                                `json:"repository_id,omitempty"`
	Valid                      bool                                  `json:"valid"`
	RecordCount                int                                   `json:"record_count"`
	Problems                   []plancatalog.Problem                 `json:"problems"`
	RemoteOverlays             *portfolio.ScanResult                 `json:"remote_overlays,omitempty"`
	DiscoveryVersion           plancatalog.DiscoveryVersion          `json:"discovery_version,omitempty"`
	ConfiguredDiscoveryVersion plancatalog.DiscoveryVersion          `json:"configured_discovery_version,omitempty"`
	ConfiguredCatalogMode      plancatalog.CatalogMode               `json:"configured_catalog_mode,omitempty"`
	TargetCatalogMode          plancatalog.CatalogMode               `json:"target_catalog_mode,omitempty"`
	PolicyDigest               string                                `json:"policy_digest,omitempty"`
	CandidateSetDigest         string                                `json:"candidate_set_digest,omitempty"`
	Complete                   *bool                                 `json:"complete,omitempty"`
	MixedCoverageComplete      *bool                                 `json:"mixed_coverage_complete,omitempty"`
	CanonicalReady             *bool                                 `json:"canonical_ready,omitempty"`
	MigrationEvidence          *plancatalog.MigrationEvidenceBinding `json:"migration_evidence,omitempty"`
	Candidates                 []plancatalog.Candidate               `json:"candidates,omitempty"`
	PreviewCandidates          []plancatalog.Candidate               `json:"preview_candidates,omitempty"`
	PreviewProblems            []plancatalog.Problem                 `json:"preview_problems,omitempty"`
	PreviewValid               *bool                                 `json:"preview_valid,omitempty"`
}

func RunPlans(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		return writeArgError(stderr, "plans", fmt.Errorf("the following arguments are required: subcommand"))
	}
	switch args[0] {
	case "validate":
		return runPlansValidate(args[1:], stdout, stderr)
	case "list":
		return runPlansList(args[1:], stdout, stderr)
	case "inventory":
		return runPlansInventory(args[1:], stdout, stderr)
	case "migrate":
		return runPlansMigrate(args[1:], stdout, stderr)
	case "catalog-activate":
		return runPlansCatalogActivate(args[1:], stdout, stderr)
	default:
		return writeArgError(stderr, "plans", fmt.Errorf("invalid subcommand %q (choose from validate, list, inventory, migrate, catalog-activate)", args[0]))
	}
}

func runPlansValidate(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--json": false, "--remote-overlays": false, "--target-discovery-version": true, "--target-catalog-mode": true, "--include-untracked": false})
	if err != nil {
		return writeArgError(stderr, "plans validate", err)
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "plans validate", err)
	}
	options, err := planReadSnapshotOptions(values, bools)
	if err != nil {
		return writeArgError(stderr, "plans validate", err)
	}
	var remoteResult *portfolio.ScanResult
	configuredSettings, _ := plancatalog.LoadRepositorySettings(root)
	if bools["--remote-overlays"] && configuredSettings.MigrationEvidence != nil && configuredSettings.MigrationEvidence.EvidenceScope == "remote-aware" {
		target, targetErr := remotePlanTargetForBinding(root, configuredSettings.MigrationEvidence)
		if targetErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans validate: error: %v\n", targetErr)
			return 1
		}
		remote, scanErr := scanRemoteOverlaysForTarget(root, target)
		if scanErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans validate: error: %v\n", scanErr)
			return 1
		}
		remoteResult = &remote
		overlay := remoteOverlayEvidence(remote, configuredSettings.RepositoryID)
		options.ObservedRemoteOverlayDigest = overlay.Digest
		options.ObservedRemoteBranchEvidence = append([]plancatalog.BranchCandidateEvidence{}, overlay.BranchTouchCandidates...)
	}
	snapshot, err := plancatalog.LoadRepositorySnapshotWithOptions(root, options)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans validate: error: %v\n", err)
		return 1
	}
	if bools["--remote-overlays"] && remoteResult == nil {
		remote, scanErr := scanRemoteOverlays(root)
		if scanErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans validate: error: %v\n", scanErr)
			return 1
		}
		remoteResult = &remote
	}
	output := plansValidationOutput{SchemaVersion: 1, Mode: snapshot.Settings.Mode, RepositoryID: snapshot.Settings.RepositoryID, Valid: snapshot.Complete && !plancatalog.HasErrors(snapshot.Problems), RecordCount: len(snapshot.Records), Problems: snapshot.Problems}
	if snapshot.Settings.DiscoveryVersion == plancatalog.DiscoveryV2 || snapshot.Targeted {
		complete := snapshot.Complete
		output.SchemaVersion = 3
		output.DiscoveryVersion = snapshot.Settings.DiscoveryVersion
		output.ConfiguredDiscoveryVersion = snapshot.ConfiguredSettings.DiscoveryVersion
		output.ConfiguredCatalogMode = snapshot.ConfiguredSettings.Mode
		output.TargetCatalogMode = snapshot.Settings.Mode
		output.PolicyDigest = snapshot.PolicyDigest
		output.CandidateSetDigest = snapshot.CandidateSetDigest
		output.Complete = &complete
		mixedComplete := complete && !plancatalog.HasErrors(snapshot.Problems)
		canonicalReady := mixedComplete && len(snapshot.DeferredPaths) == 0
		output.MixedCoverageComplete = &mixedComplete
		output.CanonicalReady = &canonicalReady
		output.MigrationEvidence = snapshot.Settings.MigrationEvidence
		output.Candidates = append([]plancatalog.Candidate{}, snapshot.Candidates...)
		output.PreviewCandidates = append([]plancatalog.Candidate{}, snapshot.PreviewCandidates...)
		output.PreviewProblems = append([]plancatalog.Problem{}, snapshot.PreviewProblems...)
		if options.IncludeUntracked {
			previewValid := !plancatalog.HasErrors(snapshot.PreviewProblems)
			output.PreviewValid = &previewValid
		}
	}
	if remoteResult != nil {
		remote := *remoteResult
		output.RemoteOverlays = &remote
		if !catalogMemberComplete(remote.Catalog, snapshot.Settings.RepositoryID) {
			output.Valid = false
		}
		if snapshot.Settings.MigrationEvidence != nil && snapshot.Settings.MigrationEvidence.EvidenceScope == "remote-aware" {
			overlay := remoteOverlayEvidence(remote, snapshot.Settings.RepositoryID)
			expectedOverlay := plancatalog.BoundRemoteOverlayDigest(root, snapshot.Settings.MigrationEvidence)
			if overlay.Status != "complete" || overlay.Digest == "" || overlay.Digest != expectedOverlay {
				output.Valid = false
				output.Problems = append(output.Problems, plancatalog.Problem{Code: "remote_overlay_stale", Message: "fresh target-member overlay does not match bound reviewed evidence", Severity: plancatalog.SeverityError})
			}
		}
	}
	if bools["--json"] {
		if err := writeJSON(stdout, output); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans validate: error: %v\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(stdout, "Plan validation (%s mode): %d record(s); valid=%t.\n", output.Mode, output.RecordCount, output.Valid)
		if output.SchemaVersion >= 2 {
			owned, excluded, blocked, unowned := validationCandidateCounts(output.Candidates)
			fmt.Fprintf(stdout, "Discovery v%d (configured v%d); target mode=%s; candidates included=%d excluded=%d blocked=%d unowned=%d preview=%d; complete=%t; mixed coverage complete=%t; canonical ready=%t.\n", output.DiscoveryVersion, output.ConfiguredDiscoveryVersion, output.TargetCatalogMode, owned, excluded, blocked, unowned, len(output.PreviewCandidates), output.Complete != nil && *output.Complete, output.MixedCoverageComplete != nil && *output.MixedCoverageComplete, output.CanonicalReady != nil && *output.CanonicalReady)
		}
		if output.PreviewValid != nil {
			fmt.Fprintf(stdout, "Preview valid=%t; preview evidence is non-authoritative.\n", *output.PreviewValid)
		}
		for _, problem := range output.Problems {
			fmt.Fprintf(stdout, "- %s %s: %s", problem.Severity, problem.Code, problem.Message)
			if problem.Path != "" {
				fmt.Fprintf(stdout, " (%s)", problem.Path)
			}
			fmt.Fprintln(stdout)
		}
		for _, problem := range output.PreviewProblems {
			fmt.Fprintf(stdout, "- preview %s %s: %s", problem.Severity, problem.Code, problem.Message)
			if problem.Path != "" {
				fmt.Fprintf(stdout, " (%s)", problem.Path)
			}
			fmt.Fprintln(stdout)
		}
		if output.RemoteOverlays != nil {
			fmt.Fprintf(stdout, "Remote overlays: complete=%t; %d pushed observation(s) across %d member(s). Local heads, worktree changes, and unpushed commits are omitted.\n", output.RemoteOverlays.Catalog.Complete, len(output.RemoteOverlays.Catalog.Observations), len(output.RemoteOverlays.Catalog.Members))
		}
	}
	if output.Valid && (output.PreviewValid == nil || *output.PreviewValid) {
		return 0
	}
	return 1
}

func runPlansList(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--format": true, "--json": false, "--target-discovery-version": true, "--remote-overlays": false})
	if err != nil {
		return writeArgError(stderr, "plans list", err)
	}
	format := values["--format"]
	switch format {
	case "", "text", "json":
	default:
		return writeArgError(stderr, "plans list", fmt.Errorf("invalid --format %q (choose from text, json)", format))
	}
	if bools["--json"] {
		if format == "text" {
			return writeArgError(stderr, "plans list", fmt.Errorf("--json conflicts with --format text"))
		}
		format = "json"
	}
	if format == "" {
		format = "text"
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "plans list", err)
	}
	options, err := planReadSnapshotOptions(values, bools)
	if err != nil {
		return writeArgError(stderr, "plans list", err)
	}
	if bools["--remote-overlays"] {
		settings, _ := plancatalog.LoadRepositorySettings(root)
		var remote portfolio.ScanResult
		var scanErr error
		if settings.MigrationEvidence != nil && settings.MigrationEvidence.EvidenceScope == "remote-aware" {
			target, targetErr := remotePlanTargetForBinding(root, settings.MigrationEvidence)
			if targetErr != nil {
				fmt.Fprintf(stderr, "codeheart-operating-kit plans list: error: %v\n", targetErr)
				return 1
			}
			remote, scanErr = scanRemoteOverlaysForTarget(root, target)
		} else {
			remote, scanErr = scanRemoteOverlays(root)
		}
		if scanErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans list: error: %v\n", scanErr)
			return 1
		}
		if settings.MigrationEvidence != nil && settings.MigrationEvidence.EvidenceScope == "remote-aware" {
			overlay := remoteOverlayEvidence(remote, settings.RepositoryID)
			options.ObservedRemoteOverlayDigest = overlay.Digest
			options.ObservedRemoteBranchEvidence = append([]plancatalog.BranchCandidateEvidence{}, overlay.BranchTouchCandidates...)
		}
	}
	snapshot, err := plancatalog.LoadRepositorySnapshotWithOptions(root, options)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans list: error: %v\n", err)
		return 1
	}
	view := plancatalog.BuildView(snapshot)
	if format == "json" {
		if err := plancatalog.WriteViewJSON(stdout, view); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans list: error: %v\n", err)
			return 1
		}
	} else {
		plancatalog.WriteViewText(stdout, view)
	}
	if plancatalog.HasErrors(view.Problems) {
		return 1
	}
	return 0
}

func runPlansInventory(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--output": true, "--json": false, "--remote-overlays": false, "--target-discovery-version": true, "--target-catalog-mode": true, "--include-untracked": false})
	if err != nil {
		return writeArgError(stderr, "plans inventory", err)
	}
	if values["--output"] == "" {
		return writeArgError(stderr, "plans inventory", fmt.Errorf("option --output requires a value"))
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "plans inventory", err)
	}
	options, err := planReadSnapshotOptions(values, bools)
	if err != nil {
		return writeArgError(stderr, "plans inventory", err)
	}
	var remoteResult *portfolio.ScanResult
	configuredSettings, _ := plancatalog.LoadRepositorySettings(root)
	if bools["--remote-overlays"] && configuredSettings.MigrationEvidence != nil && configuredSettings.MigrationEvidence.EvidenceScope == "remote-aware" {
		target, targetErr := remotePlanTargetForBinding(root, configuredSettings.MigrationEvidence)
		if targetErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", targetErr)
			return 1
		}
		remote, scanErr := scanRemoteOverlaysForTarget(root, target)
		if scanErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", scanErr)
			return 1
		}
		remoteResult = &remote
		overlay := remoteOverlayEvidence(remote, configuredSettings.RepositoryID)
		options.ObservedRemoteOverlayDigest = overlay.Digest
		options.ObservedRemoteBranchEvidence = append([]plancatalog.BranchCandidateEvidence{}, overlay.BranchTouchCandidates...)
	}
	inventory, err := plancatalog.BuildInventoryWithOptions(root, time.Now(), options)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", err)
		return 1
	}
	if bools["--remote-overlays"] && remoteResult == nil {
		var target *portfolio.RemotePlanTarget
		if inventory.SchemaVersion == 3 {
			target = &portfolio.RemotePlanTarget{
				RepositoryID: inventory.RepositoryID, EvidenceRevision: inventory.EvidenceRevision,
				DiscoveryVersion: plancatalog.DiscoveryV2, TargetCatalogMode: inventory.TargetCatalogMode,
				PolicyDigest: inventory.PolicyDigest, CandidateSetDigest: inventory.CandidateSetDigest,
			}
		}
		remote, scanErr := scanRemoteOverlaysForTarget(root, target)
		if scanErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", scanErr)
			return 1
		}
		remoteResult = &remote
	}
	if remoteResult != nil {
		remote := *remoteResult
		overlay := remoteOverlayEvidence(remote, inventory.RepositoryID)
		plancatalog.AttachRemoteOverlay(root, &inventory, overlay)
		if overlay.Status != "complete" {
			inventory.Problems = append(inventory.Problems, plancatalog.Problem{Code: "remote_overlay_coverage_incomplete", Message: "pushed default and unmerged-branch evidence is incomplete for this repository", Severity: plancatalog.SeverityError})
		}
		for _, scanErr := range remote.Catalog.Errors {
			if scanErr.RepositoryID == "" || scanErr.RepositoryID == inventory.RepositoryID {
				inventory.Problems = append(inventory.Problems, plancatalog.Problem{Code: scanErr.Code, Message: scanErr.Message, Path: scanErr.Path, Severity: plancatalog.SeverityError})
			}
		}
		plancatalog.SortProblems(inventory.Problems)
	}
	if err := writeInventoryArtifact(root, values["--output"], inventory); err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, inventory); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans inventory: error: %v\n", err)
			return 1
		}
	} else {
		fmt.Fprintf(stdout, "Plan inventory: %d formal record(s), %d canonical, %d legacy, %d invalid, %d unreadable, %d unsafe, %d unpaired legacy; revision %s.\n", inventory.Coverage.FormalRecords, inventory.Coverage.CanonicalMetadata, inventory.Coverage.LegacyRecords, inventory.Coverage.InvalidRecords, inventory.Coverage.UnreadableRecords, inventory.Coverage.UnsafeRecords, inventory.Coverage.UnpairedLegacyEvidence, inventory.SourceRevision)
		if inventory.SchemaVersion >= 2 {
			fmt.Fprintf(stdout, "Discovery v%d (configured v%d); target mode=%s; candidates included=%d excluded=%d blocked=%d unowned=%d preview=%d; mixed coverage complete=%t; canonical ready=%t.\n", inventory.DiscoveryVersion, inventory.ConfiguredDiscoveryVersion, inventory.TargetCatalogMode, inventory.Coverage.IncludedCandidates, inventory.Coverage.ExcludedCandidates, inventory.Coverage.BlockedCandidates, inventory.Coverage.UnownedCandidates, inventory.Coverage.PreviewCandidates, inventory.MixedCoverageComplete != nil && *inventory.MixedCoverageComplete, inventory.CanonicalReady != nil && *inventory.CanonicalReady)
			if inventory.BranchEvidence != nil {
				fmt.Fprintf(stdout, "Branch evidence: algorithm=%s scope=%s overlay=%s digest=%s; touches=%d blocking=%d.\n", inventory.BranchEvidence.Algorithm, inventory.BranchEvidence.EvidenceScope, inventory.BranchEvidence.RemoteOverlayStatus, inventory.BranchEvidence.Digest, inventory.Coverage.BranchTouchCandidates, inventory.Coverage.BlockingBranchTouches)
			}
			for _, record := range inventory.Records {
				for _, evidence := range record.BranchTouchCandidates {
					fmt.Fprintf(stdout, "- branch %s | %s | %s", evidence.Identity.LogicalRef, evidence.Identity.CandidatePath, evidence.ProposedDisposition)
					for _, blocker := range evidence.Blockers {
						fmt.Fprintf(stdout, " | blocker=%s", blocker.Code)
					}
					fmt.Fprintln(stdout)
				}
			}
		}
		fmt.Fprintf(stdout, "Wrote inventory: %s\n", values["--output"])
	}
	if plancatalog.HasErrors(inventory.Problems) || plancatalog.HasErrors(inventory.PreviewProblems) {
		return 1
	}
	return 0
}

func planReadSnapshotOptions(values map[string]string, bools map[string]bool) (plancatalog.SnapshotOptions, error) {
	options := plancatalog.SnapshotOptions{IncludeUntracked: bools["--include-untracked"]}
	if value := values["--target-discovery-version"]; value != "" {
		if value != "2" {
			return options, fmt.Errorf("invalid --target-discovery-version %q (only 2 is supported)", value)
		}
		options.TargetDiscoveryVersion = plancatalog.DiscoveryV2
	}
	if value := values["--target-catalog-mode"]; value != "" {
		if value != string(plancatalog.ModeCanonical) {
			return options, fmt.Errorf("invalid --target-catalog-mode %q (only canonical is supported)", value)
		}
		options.TargetCatalogMode = plancatalog.ModeCanonical
	}
	return options, nil
}

func validationCandidateCounts(candidates []plancatalog.Candidate) (owned, excluded, blocked, unowned int) {
	for _, candidate := range candidates {
		switch candidate.Ownership {
		case plancatalog.OwnershipOwned:
			owned++
		case plancatalog.OwnershipExcluded:
			excluded++
		case plancatalog.OwnershipProspectiveBlocked:
			blocked++
		case plancatalog.OwnershipHardUnowned:
			unowned++
		}
	}
	return
}

func scanRemoteOverlays(root string) (portfolio.ScanResult, error) {
	return scanRemoteOverlaysForTarget(root, nil)
}

func scanRemoteOverlaysForTarget(root string, target *portfolio.RemotePlanTarget) (portfolio.ScanResult, error) {
	config, err := portfolio.LoadConfig(root)
	if err != nil {
		return portfolio.ScanResult{}, err
	}
	runner := portfolio.ExecRunner{}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	result, err := portfolio.Scan(ctx, portfolio.ScanOptions{Root: root, Config: config, Sources: portfolio.BuildSources(config, runner), Runner: runner, WriteCache: false, PlanTarget: target})
	if err != nil {
		return portfolio.ScanResult{}, err
	}
	return publicRemoteOverlayEvidence(result), nil
}

// publicRemoteOverlayEvidence removes machine-local discovery paths from plan
// validation and migration artifacts. The portfolio cache remains ignored local
// evidence and retains the scanner's exact locator for operational recovery.
func publicRemoteOverlayEvidence(result portfolio.ScanResult) portfolio.ScanResult {
	for index := range result.Catalog.Members {
		member := &result.Catalog.Members[index]
		member.SourceLocator = publicSourceLocator(member.SourceLocator, member.RepositoryID)
	}
	for index := range result.Catalog.Candidates {
		candidate := &result.Catalog.Candidates[index]
		candidate.SourceLocator = publicSourceLocator(candidate.SourceLocator, candidate.RepositoryID)
	}
	for index := range result.Catalog.Errors {
		scanError := &result.Catalog.Errors[index]
		scanError.SourceLocator = publicSourceLocator(scanError.SourceLocator, scanError.RepositoryID)
		scanError.Message = publicScanErrorMessage(scanError.Code)
	}
	return result
}

func publicScanErrorMessage(code string) string {
	switch code {
	case "source_discovery_failed":
		return "configured portfolio source discovery failed"
	case "self_remote_unavailable":
		return "coordination repository remote evidence is unavailable"
	case "mirror_refresh_failed":
		return "remote repository evidence could not be refreshed"
	case "membership_evidence_unavailable", "membership_validation_failed":
		return "default-branch membership evidence is unavailable or invalid"
	case "github_auth_unavailable", "github_access_denied", "github_rate_limit_exhausted", "github_retry_exhausted":
		return "configured GitHub source could not be read completely"
	default:
		return "remote portfolio scan reported an error"
	}
}

func publicSourceLocator(locator, repositoryID string) string {
	trimmed := strings.TrimSpace(locator)
	if trimmed == "" || (!filepath.IsAbs(trimmed) && !strings.HasPrefix(strings.ToLower(trimmed), "file:")) {
		return locator
	}
	if repositoryID != "" {
		return "local-git:" + repositoryID
	}
	return "local-git"
}

func catalogContainsMember(catalog portfolio.Catalog, repositoryID string) bool {
	for _, member := range catalog.Members {
		if member.RepositoryID == repositoryID {
			return true
		}
	}
	return false
}

func catalogMemberComplete(catalog portfolio.Catalog, repositoryID string) bool {
	for _, member := range catalog.Members {
		if member.RepositoryID == repositoryID {
			return member.Complete != nil && *member.Complete
		}
	}
	return false
}

func remoteOverlayEvidence(result portfolio.ScanResult, repositoryID string) plancatalog.RemoteOverlayEvidence {
	overlay := plancatalog.RemoteOverlayEvidence{Remote: repositoryID, Status: "incomplete", ObservedAt: result.Catalog.CompletedAt, Refs: []plancatalog.RemoteOverlayRef{}, BranchTouchCandidates: []plancatalog.BranchCandidateEvidence{}}
	found := false
	for _, member := range result.Catalog.Members {
		if member.RepositoryID != repositoryID {
			continue
		}
		found = true
		overlay.SourceIdentitySHA256 = member.SourceIdentitySHA256
		if member.Complete != nil && *member.Complete {
			overlay.Status = "complete"
		}
		overlay.Refs = append(overlay.Refs, plancatalog.RemoteOverlayRef{LogicalRef: remoteLogicalRef(repositoryID, member.DefaultRef), Tip: member.SourceRevision})
	}
	if !found {
		overlay.Status = "unavailable"
	}
	for _, scanErr := range result.Catalog.Errors {
		if scanErr.RepositoryID == repositoryID {
			overlay.Status = "incomplete"
		}
	}
	for _, observation := range result.Catalog.Observations {
		if observation.RepositoryID == repositoryID && observation.Visibility == "unmerged-branch" && strings.HasPrefix(observation.Ref, "refs/") {
			overlay.Refs = append(overlay.Refs, plancatalog.RemoteOverlayRef{LogicalRef: remoteLogicalRef(repositoryID, observation.Ref), Tip: observation.Commit})
		}
	}
	for _, observation := range result.Catalog.PlanCandidates {
		if observation.RepositoryID == repositoryID && observation.Visibility == "unmerged-branch" && strings.HasPrefix(observation.Ref, "refs/") {
			overlay.Refs = append(overlay.Refs, plancatalog.RemoteOverlayRef{LogicalRef: remoteLogicalRef(repositoryID, observation.Ref), Tip: observation.Commit})
		}
	}
	for _, evidence := range result.Catalog.RemoteBranchEvidence {
		if evidence.Identity.RepositoryID == repositoryID {
			overlay.BranchTouchCandidates = append(overlay.BranchTouchCandidates, evidence)
		}
	}
	sort.SliceStable(overlay.BranchTouchCandidates, func(i, j int) bool {
		left := overlay.BranchTouchCandidates[i].Identity.LogicalRef + "\x00" + overlay.BranchTouchCandidates[i].Identity.CandidatePath
		right := overlay.BranchTouchCandidates[j].Identity.LogicalRef + "\x00" + overlay.BranchTouchCandidates[j].Identity.CandidatePath
		return left < right
	})
	sort.SliceStable(overlay.Refs, func(i, j int) bool {
		if overlay.Refs[i].LogicalRef != overlay.Refs[j].LogicalRef {
			return overlay.Refs[i].LogicalRef < overlay.Refs[j].LogicalRef
		}
		return overlay.Refs[i].Tip < overlay.Refs[j].Tip
	})
	unique := overlay.Refs[:0]
	for _, ref := range overlay.Refs {
		if len(unique) == 0 || unique[len(unique)-1] != ref {
			unique = append(unique, ref)
		}
	}
	overlay.Refs = unique
	overlay.Digest = plancatalog.CanonicalRemoteOverlayDigest(overlay)
	return overlay
}

func remoteLogicalRef(repositoryID, ref string) string {
	ref = strings.TrimSpace(ref)
	switch {
	case strings.HasPrefix(ref, "refs/heads/"):
	case strings.HasPrefix(ref, "refs/remotes/origin/"):
		ref = "refs/heads/" + strings.TrimPrefix(ref, "refs/remotes/origin/")
	case strings.HasPrefix(ref, "refs/"):
	default:
		ref = "refs/heads/" + ref
	}
	return "remote:" + repositoryID + ":" + ref
}

func remotePlanTargetForLedger(root string, ledger plancatalog.MigrationLedger) (*portfolio.RemotePlanTarget, error) {
	activationBase, err := plancatalog.LedgerActivationBaseRevision(root, ledger)
	if err != nil {
		return nil, err
	}
	activationRevision, _ := plancatalog.ActivationCheckpointRevision(root, activationBase)
	target := &portfolio.RemotePlanTarget{
		RepositoryID: ledger.RepositoryID, EvidenceRevision: ledger.EvidenceRevision,
		ActivationBaseRevision: activationBase, ActivationRevision: activationRevision, LedgerPath: ledger.ArtifactPath, LedgerSHA256: ledger.ArtifactSHA256,
		DiscoveryVersion: plancatalog.DiscoveryV2, TargetCatalogMode: ledger.TargetCatalogMode,
		PolicyDigest: ledger.PolicyDigest, CandidateSetDigest: ledger.CandidateSetDigest,
	}
	if ledger.BranchEvidence != nil {
		target.BranchEvidenceDigest = ledger.BranchEvidence.Digest
		target.EvidenceScope = ledger.BranchEvidence.EvidenceScope
		target.RemoteOverlayDigest = ledger.BranchEvidence.RemoteOverlayDigest
		target.RemoteSourceIdentity = ledger.BranchEvidence.RemoteSourceIdentitySHA256
	}
	return target, nil
}

func remotePlanTargetForBinding(root string, binding *plancatalog.MigrationEvidenceBinding) (*portfolio.RemotePlanTarget, error) {
	if binding == nil {
		return nil, fmt.Errorf("migration_evidence_binding_missing: remote-aware evidence binding is required")
	}
	ledger, err := plancatalog.LoadMigrationLedgerArtifact(root, binding.LedgerPath)
	if err != nil {
		return nil, err
	}
	if ledger.EvidenceRevision != binding.EvidenceRevision || ledger.ArtifactSHA256 != binding.LedgerSHA256 || ledger.BranchEvidence == nil || ledger.BranchEvidence.Digest != binding.BranchEvidenceDigest || ledger.BranchEvidence.EvidenceScope != binding.EvidenceScope {
		return nil, fmt.Errorf("migration_evidence_binding_mismatch: config binding differs from the reviewed ledger")
	}
	if binding.EvidenceScope == "remote-aware" && (ledger.BranchEvidence.RemoteOverlayStatus != "complete" || ledger.BranchEvidence.RemoteOverlayDigest == "" || ledger.BranchEvidence.RemoteSourceIdentitySHA256 == "") {
		return nil, fmt.Errorf("migration_evidence_binding_mismatch: remote-aware config binding lacks complete ledger overlay authority")
	}
	activationRevision, err := plancatalog.BoundActivationCheckpointRevision(root, binding.ActivationBaseRevision, "HEAD", binding.MigrationActionDigest)
	if err != nil {
		return nil, err
	}
	return &portfolio.RemotePlanTarget{
		RepositoryID: ledger.RepositoryID, EvidenceRevision: binding.EvidenceRevision,
		ActivationBaseRevision: binding.ActivationBaseRevision, ActivationRevision: activationRevision,
		LedgerPath: binding.LedgerPath, LedgerSHA256: binding.LedgerSHA256,
		DiscoveryVersion: plancatalog.DiscoveryV2, TargetCatalogMode: ledger.TargetCatalogMode,
		PolicyDigest: ledger.PolicyDigest, CandidateSetDigest: ledger.CandidateSetDigest,
		MigrationActionDigest: binding.MigrationActionDigest,
		BranchEvidenceDigest:  binding.BranchEvidenceDigest, EvidenceScope: binding.EvidenceScope,
		RemoteOverlayDigest: ledger.BranchEvidence.RemoteOverlayDigest, RemoteSourceIdentity: ledger.BranchEvidence.RemoteSourceIdentitySHA256,
	}, nil
}

func runPlansMigrate(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--ledger": true, "--dry-run": false, "--yes": false, "--json": false, "--remote-overlays": false})
	if err != nil {
		return writeArgError(stderr, "plans migrate", err)
	}
	if values["--ledger"] == "" {
		return writeArgError(stderr, "plans migrate", fmt.Errorf("option --ledger requires a value"))
	}
	if bools["--dry-run"] == bools["--yes"] {
		return writeArgError(stderr, "plans migrate", fmt.Errorf("choose exactly one of --dry-run or --yes"))
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "plans migrate", err)
	}
	ledger, err := plancatalog.LoadMigrationLedgerArtifact(root, values["--ledger"])
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", err)
		return 1
	}
	var remoteCheck func() (string, error)
	if bools["--remote-overlays"] {
		target, targetErr := remotePlanTargetForLedger(root, ledger)
		if targetErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", targetErr)
			return 1
		}
		var observedOverlay plancatalog.RemoteOverlayEvidence
		remoteCheck = func() (string, error) {
			remote, scanErr := scanRemoteOverlaysForTarget(root, target)
			if scanErr != nil {
				return "", scanErr
			}
			overlay := remoteOverlayEvidence(remote, ledger.RepositoryID)
			if overlay.Status != "complete" {
				return "", fmt.Errorf("remote_overlay_incomplete: target member overlay is not complete")
			}
			observedOverlay = overlay
			return overlay.Digest, nil
		}
		digest, scanErr := remoteCheck()
		if scanErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", scanErr)
			return 1
		}
		ledger.ObservedRemoteOverlayDigest = digest
		ledger.ObservedRemoteBranchEvidence = append([]plancatalog.BranchCandidateEvidence{}, observedOverlay.BranchTouchCandidates...)
	}
	plan, err := plancatalog.BuildMigrationPlan(root, ledger)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", err)
		return 1
	}
	outcome, err := plancatalog.ExecuteMigration(plan, plancatalog.MigrationApplyOptions{DryRun: bools["--dry-run"], Now: time.Now(), RemoteOverlayCheck: remoteCheck})
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, outcome); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans migrate: error: %v\n", err)
			return 1
		}
	} else {
		reconcile.WriteText(stdout, outcome.Result)
		if outcome.Projection != nil {
			fmt.Fprintf(stdout, "Projected %s readiness: candidates=%d canonical=%d grandfathered=%d gaps=%d ready=%t; activation performed=false.\n", outcome.TargetCatalogMode, outcome.Projection.Candidates, outcome.Projection.Canonical, outcome.Projection.Grandfathered, outcome.Projection.Gaps, outcome.Projection.Ready)
		}
		for _, skip := range outcome.Skips {
			fmt.Fprintf(stdout, "- skipped %s: %s (%s)\n", skip.Code, skip.Message, skip.Path)
		}
		for _, review := range outcome.BranchReviews {
			fmt.Fprintf(stdout, "- branch review %s: disposition=%s verified=%t", review.Identity.LogicalRef, review.Disposition, review.Verified)
			for _, blocker := range review.Blockers {
				fmt.Fprintf(stdout, " blocker=%s", blocker.Code)
			}
			fmt.Fprintln(stdout)
		}
		for _, followUp := range outcome.FollowUps {
			fmt.Fprintf(stdout, "- pending incremental migration: %s for %s when %s\n", followUp.Action, followUp.PlanPath, followUp.Trigger)
		}
	}
	if !outcome.Result.OK() || plancatalog.HasMaterialMigrationSkips(outcome.Skips) {
		return 1
	}
	return 0
}

func runPlansCatalogActivate(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--ledger": true, "--dry-run": false, "--yes": false, "--json": false, "--remote-overlays": false})
	if err != nil {
		return writeArgError(stderr, "plans catalog-activate", err)
	}
	if values["--ledger"] == "" {
		return writeArgError(stderr, "plans catalog-activate", fmt.Errorf("option --ledger requires a value"))
	}
	if bools["--dry-run"] == bools["--yes"] {
		return writeArgError(stderr, "plans catalog-activate", fmt.Errorf("choose exactly one of --dry-run or --yes"))
	}
	root, err := singleRoot(positionals)
	if err != nil {
		return writeArgError(stderr, "plans catalog-activate", err)
	}
	ledger, err := plancatalog.LoadMigrationLedgerArtifact(root, values["--ledger"])
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans catalog-activate: error: %v\n", err)
		return 1
	}
	var remoteCheck func() (string, error)
	if bools["--remote-overlays"] {
		target, targetErr := remotePlanTargetForLedger(root, ledger)
		if targetErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans catalog-activate: error: %v\n", targetErr)
			return 1
		}
		var observedOverlay plancatalog.RemoteOverlayEvidence
		remoteCheck = func() (string, error) {
			remote, scanErr := scanRemoteOverlaysForTarget(root, target)
			if scanErr != nil {
				return "", scanErr
			}
			overlay := remoteOverlayEvidence(remote, ledger.RepositoryID)
			if overlay.Status != "complete" {
				return "", fmt.Errorf("remote_overlay_incomplete: target member overlay is not complete")
			}
			observedOverlay = overlay
			return overlay.Digest, nil
		}
		digest, scanErr := remoteCheck()
		if scanErr != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans catalog-activate: error: %v\n", scanErr)
			return 1
		}
		ledger.ObservedRemoteOverlayDigest = digest
		ledger.ObservedRemoteBranchEvidence = append([]plancatalog.BranchCandidateEvidence{}, observedOverlay.BranchTouchCandidates...)
	}
	plan, err := plancatalog.BuildCatalogActivationPlan(root, ledger)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans catalog-activate: error: %v\n", err)
		return 1
	}
	outcome, err := plancatalog.ExecuteCatalogActivation(plan, plancatalog.CatalogActivationOptions{DryRun: bools["--dry-run"], Now: time.Now(), RemoteOverlayCheck: remoteCheck})
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit plans catalog-activate: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, outcome); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit plans catalog-activate: error: %v\n", err)
			return 1
		}
	} else {
		reconcile.WriteText(stdout, outcome.Result)
		fmt.Fprintf(stdout, "Catalog activation: mixed coverage complete=%t canonical ready=%t; activation performed=%t; checkpoint required=%t.\n", outcome.ProjectedCoverage.MixedCoverageComplete, outcome.ProjectedCoverage.CanonicalReady, outcome.ActivationPerformed, outcome.ActivationCheckpointRequired)
		for _, followUp := range outcome.IncrementalFollowUps {
			fmt.Fprintf(stdout, "- pending follow-up %s for %s after %s\n", followUp.Action, followUp.PlanPath, followUp.Trigger)
		}
	}
	if !outcome.Result.OK() {
		return 1
	}
	return 0
}

func singleRoot(positionals []string) (string, error) {
	root := "."
	if len(positionals) > 0 {
		root = expandPath(positionals[0])
	}
	if len(positionals) > 1 {
		return "", fmt.Errorf("unexpected argument %q", positionals[1])
	}
	return root, nil
}

func writeInventoryArtifact(root, destination string, inventory plancatalog.Inventory) error {
	return writeInventoryArtifactWithHook(root, destination, inventory, nil)
}

func writeInventoryArtifactWithHook(root, destination string, inventory plancatalog.Inventory, hook func(string) error) error {
	destination = expandPath(destination)
	absDestination, err := filepath.Abs(destination)
	if err != nil {
		return err
	}
	rootAbsolute, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	parentCandidate := filepath.Dir(absDestination)
	parentCandidateInfo, err := os.Lstat(parentCandidate)
	if os.IsNotExist(err) {
		return fmt.Errorf("inventory_target_parent_missing: create the explicit output directory before writing inventory")
	}
	if err != nil {
		return err
	}
	if !parentCandidateInfo.IsDir() || parentCandidateInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("inventory_target_unsafe: output parent is not a regular directory")
	}
	parentPath, err := filepath.EvalSymlinks(parentCandidate)
	if err != nil {
		return err
	}
	parentRoot, err := os.OpenRoot(parentPath)
	if err != nil {
		return err
	}
	defer parentRoot.Close()
	parentInfo, err := parentRoot.Stat(".")
	if err != nil {
		return err
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return err
	}
	if err := invokeInventoryHook(hook, "inventory-parent-bound-before-protection"); err != nil {
		return err
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return err
	}
	absDestination = filepath.Join(parentPath, filepath.Base(absDestination))
	if err := rejectRepoLocalSymlinkTraversal(rootAbsolute, absDestination); err != nil {
		return err
	}
	if info, err := os.Lstat(absDestination); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("inventory_target_unsafe: output target is a symbolic link")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	resolvedDestination, err := resolvePathAllowMissing(absDestination)
	if err != nil {
		return fmt.Errorf("inventory_target_unsafe: resolve output target: %w", err)
	}
	resolvedRoot, err := filepath.EvalSymlinks(rootAbsolute)
	if err != nil {
		return err
	}
	if relative, relErr := filepath.Rel(resolvedRoot, resolvedDestination); relErr == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		repositoryPath := filepath.ToSlash(relative)
		if strings.EqualFold(repositoryPath, "docs/repo/plans/README.md") || strings.EqualFold(repositoryPath, plancatalog.LegacyRegisterPath) {
			return fmt.Errorf("inventory_target_is_plan: output cannot create or overwrite planning authority %s", repositoryPath)
		}
		_, formal := plancatalog.FormalPathKind(resolvedRoot, repositoryPath)
		if inventory.DiscoveryVersion == plancatalog.DiscoveryV2 {
			_, formal = plancatalog.ProspectiveV2PathKind(repositoryPath)
			signal, signalErr := plancatalog.WorktreeV2PlanSignal(resolvedRoot, repositoryPath)
			if signalErr != nil {
				return fmt.Errorf("inventory_target_unsafe: classify existing target: %w", signalErr)
			}
			formal = formal || signal
		}
		if formal {
			return fmt.Errorf("inventory_target_is_plan: output cannot create or overwrite formal plan path %s", repositoryPath)
		}
	}
	protected := []string{"docs/repo/plans/README.md", "docs/repo/plans/plan-register.md"}
	candidates := append([]plancatalog.Candidate{}, inventory.Candidates...)
	candidates = append(candidates, inventory.PreviewCandidates...)
	if inventory.DiscoveryVersion != plancatalog.DiscoveryV2 {
		candidates, err = plancatalog.Enumerate(root)
		if err != nil {
			return fmt.Errorf("enumerate protected plan targets: %w", err)
		}
	}
	for _, candidate := range candidates {
		protected = append(protected, candidate.Path)
	}
	for _, path := range protected {
		canonical, err := filepath.Abs(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			return err
		}
		resolvedProtected, err := resolvePathAllowMissing(canonical)
		if err != nil {
			return err
		}
		if sameFileOrPath(resolvedDestination, resolvedProtected) {
			return fmt.Errorf("inventory_target_is_plan: output cannot overwrite planning authority %s", path)
		}
	}
	absDestination = resolvedDestination
	if inventory.SchemaVersion == 3 {
		encoded, marshalErr := json.Marshal(inventory)
		if marshalErr != nil {
			return marshalErr
		}
		var value any
		if unmarshalErr := json.Unmarshal(encoded, &value); unmarshalErr != nil {
			return unmarshalErr
		}
		schemaPath, schemaErr := state.SchemaForPlanInventoryVersion(inventory.SchemaVersion)
		if schemaErr != nil {
			return schemaErr
		}
		if schemaErr = state.Validate(schemaPath, value); schemaErr != nil {
			return fmt.Errorf("inventory_schema_invalid: %w", schemaErr)
		}
	}
	var data []byte
	if strings.EqualFold(filepath.Ext(absDestination), ".json") {
		data, err = json.MarshalIndent(inventory, "", "  ")
		data = append(data, '\n')
	} else {
		data, err = state.EncodeYAML(inventory)
	}
	if err != nil {
		return err
	}
	if !samePath(filepath.Join(parentPath, filepath.Base(absDestination)), absDestination) {
		return fmt.Errorf("inventory_target_unsafe: destination parent identity changed during creation")
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return err
	}
	targetName := filepath.Base(absDestination)
	if targetName == "." || targetName == string(filepath.Separator) {
		return fmt.Errorf("inventory_target_unsafe: output target must be a file")
	}
	existingInfo, statErr := parentRoot.Lstat(targetName)
	if statErr == nil && (!existingInfo.Mode().IsRegular() || existingInfo.Mode()&os.ModeSymlink != 0) {
		return fmt.Errorf("inventory_target_unsafe: existing output is not a regular file")
	}
	if statErr != nil && !os.IsNotExist(statErr) {
		return statErr
	}
	var existingFile *os.File
	var existingData []byte
	if statErr == nil {
		existingFile, err = parentRoot.Open(targetName)
		if err != nil {
			return err
		}
		defer func() {
			if existingFile != nil {
				_ = existingFile.Close()
			}
		}()
		openedInfo, err := existingFile.Stat()
		if err != nil || !os.SameFile(existingInfo, openedInfo) {
			return fmt.Errorf("inventory output identity changed while binding existing bytes")
		}
		existingData, err = io.ReadAll(existingFile)
		if err != nil {
			return err
		}
		currentInfo, err := parentRoot.Lstat(targetName)
		if err != nil || !os.SameFile(existingInfo, currentInfo) {
			return fmt.Errorf("inventory output identity changed while reading existing bytes")
		}
	}
	temporaryName, temporary, err := createInventoryTemp(parentRoot)
	if err != nil {
		return err
	}
	temporaryInfo, err := temporary.Stat()
	if err != nil {
		_ = temporary.Close()
		return err
	}
	defer parentRoot.Remove(temporaryName)
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := invokeInventoryHook(hook, "inventory-parent-bound"); err != nil {
		return err
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return err
	}
	var quarantine string
	if statErr == nil {
		quarantine, err = inventorySibling(parentRoot, targetName, "inventory-replaced")
		if err != nil {
			return err
		}
		if err := parentRoot.Rename(targetName, quarantine); err != nil {
			return err
		}
		capturedInfo, captureErr := parentRoot.Lstat(quarantine)
		if captureErr != nil || !os.SameFile(existingInfo, capturedInfo) {
			return fmt.Errorf("inventory output identity changed during replacement; retained captured entry at %s", filepath.Join(parentPath, quarantine))
		}
		if _, err := existingFile.Seek(0, io.SeekStart); err != nil {
			return err
		}
		capturedData, err := io.ReadAll(existingFile)
		if err != nil || !bytes.Equal(existingData, capturedData) {
			return fmt.Errorf("inventory output bytes changed during replacement; retained captured entry at %s", filepath.Join(parentPath, quarantine))
		}
		openedInfo, err := existingFile.Stat()
		if err != nil || !os.SameFile(existingInfo, openedInfo) {
			return fmt.Errorf("inventory output descriptor identity changed during replacement; retained captured entry at %s", filepath.Join(parentPath, quarantine))
		}
		if err := invokeInventoryHook(hook, "inventory-output-quarantined"); err != nil {
			return err
		}
	}
	if err := parentRoot.Link(temporaryName, targetName); err != nil {
		if quarantine != "" {
			if _, targetErr := parentRoot.Lstat(targetName); os.IsNotExist(targetErr) {
				_ = parentRoot.Link(quarantine, targetName)
			}
		}
		return fmt.Errorf("install inventory output without replacement: %w", err)
	}
	installedInfo, err := parentRoot.Lstat(targetName)
	if err != nil || !os.SameFile(temporaryInfo, installedInfo) {
		return fmt.Errorf("inventory output identity changed during installation")
	}
	installedData, err := parentRoot.ReadFile(targetName)
	if err != nil || !bytes.Equal(installedData, data) {
		return fmt.Errorf("inventory output bytes changed during installation")
	}
	if err := parentRoot.Remove(temporaryName); err != nil {
		return err
	}
	if quarantine != "" {
		capturedInfo, err := parentRoot.Lstat(quarantine)
		if err != nil || !os.SameFile(existingInfo, capturedInfo) {
			return fmt.Errorf("replaced inventory evidence changed before cleanup; retained at %s", filepath.Join(parentPath, quarantine))
		}
		if _, err := existingFile.Seek(0, io.SeekStart); err != nil {
			return err
		}
		capturedData, err := io.ReadAll(existingFile)
		if err != nil || !bytes.Equal(existingData, capturedData) {
			return fmt.Errorf("replaced inventory bytes changed before cleanup; retained at %s", filepath.Join(parentPath, quarantine))
		}
		openedInfo, err := existingFile.Stat()
		if err != nil || !os.SameFile(existingInfo, openedInfo) {
			return fmt.Errorf("replaced inventory descriptor identity changed before cleanup; retained at %s", filepath.Join(parentPath, quarantine))
		}
		if err := parentRoot.Remove(quarantine); err != nil {
			return err
		}
		if err := existingFile.Close(); err != nil {
			return err
		}
		existingFile = nil
	}
	if err := syncInventoryDirectory(parentRoot); err != nil {
		return err
	}
	if err := validateBoundInventoryParent(parentPath, parentRoot, parentInfo); err != nil {
		return fmt.Errorf("inventory output was installed through the bound parent but requested parent authority changed: %w", err)
	}
	return nil
}

func createInventoryTemp(root *os.Root) (string, *os.File, error) {
	for attempt := 0; attempt < 8; attempt++ {
		name, err := inventoryRandomName(".plan-inventory-")
		if err != nil {
			return "", nil, err
		}
		file, err := root.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
		if os.IsExist(err) {
			continue
		}
		return name, file, err
	}
	return "", nil, fmt.Errorf("could not reserve an inventory temporary file")
}

func inventorySibling(root *os.Root, target, purpose string) (string, error) {
	for attempt := 0; attempt < 8; attempt++ {
		name, err := inventoryRandomName("." + target + "." + purpose + "-")
		if err != nil {
			return "", err
		}
		if _, err := root.Lstat(name); os.IsNotExist(err) {
			return name, nil
		} else if err != nil {
			return "", err
		}
	}
	return "", fmt.Errorf("could not reserve an inventory replacement quarantine")
}

func inventoryRandomName(prefix string) (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(nonce), nil
}

func validateBoundInventoryParent(path string, root *os.Root, identity os.FileInfo) error {
	pathInfo, err := os.Lstat(path)
	if err != nil || identity == nil || !os.SameFile(identity, pathInfo) {
		return fmt.Errorf("inventory_target_unsafe: output parent path identity changed")
	}
	boundInfo, err := root.Stat(".")
	if err != nil || !os.SameFile(identity, boundInfo) {
		return fmt.Errorf("inventory_target_unsafe: output parent descriptor identity changed")
	}
	return nil
}

func syncInventoryDirectory(root *os.Root) error {
	directory, err := root.Open(".")
	if err != nil {
		return err
	}
	defer directory.Close()
	if err := directory.Sync(); err != nil && runtime.GOOS != "windows" {
		return err
	}
	return nil
}

func invokeInventoryHook(hook func(string) error, phase string) error {
	if hook == nil {
		return nil
	}
	return hook(phase)
}

func sameFileOrPath(left, right string) bool {
	leftInfo, leftErr := os.Stat(left)
	rightInfo, rightErr := os.Stat(right)
	if leftErr == nil && rightErr == nil && os.SameFile(leftInfo, rightInfo) {
		return true
	}
	return samePath(left, right)
}

func rejectRepoLocalSymlinkTraversal(root, target string) error {
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return nil
	}
	current := root
	for _, part := range strings.Split(relative, string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("inventory_target_unsafe: output traverses symbolic link %s", current)
		}
	}
	return nil
}

func resolvePathAllowMissing(path string) (string, error) {
	current := filepath.Clean(path)
	suffix := []string{}
	for {
		if _, err := os.Lstat(current); err == nil {
			break
		} else if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("no existing ancestor for %s", path)
		}
		suffix = append(suffix, filepath.Base(current))
		current = parent
	}
	resolved, err := filepath.EvalSymlinks(current)
	if err != nil {
		return "", err
	}
	for index := len(suffix) - 1; index >= 0; index-- {
		resolved = filepath.Join(resolved, suffix[index])
	}
	return filepath.Clean(resolved), nil
}

func samePath(left, right string) bool {
	if filepath.Separator == '\\' {
		return strings.EqualFold(filepath.Clean(left), filepath.Clean(right))
	}
	return filepath.Clean(left) == filepath.Clean(right)
}
