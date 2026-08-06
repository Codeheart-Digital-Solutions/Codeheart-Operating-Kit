package plancatalog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"go.yaml.in/yaml/v3"
)

type MigrationDecision struct {
	ID                     string     `json:"id" yaml:"id"`
	Kind                   Kind       `json:"kind" yaml:"kind"`
	Purpose                string     `json:"purpose" yaml:"purpose"`
	FirstCataloged         string     `json:"first_cataloged" yaml:"first_cataloged"`
	CatalogMetadataUpdated string     `json:"catalog_metadata_updated" yaml:"catalog_metadata_updated"`
	Family                 string     `json:"family,omitempty" yaml:"family,omitempty"`
	Products               []string   `json:"products,omitempty" yaml:"products,omitempty"`
	Capabilities           []string   `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	StrategicThemes        []string   `json:"strategic_themes,omitempty" yaml:"strategic_themes,omitempty"`
	Relations              []Relation `json:"relations,omitempty" yaml:"relations,omitempty"`
}

type MigrationRecord struct {
	Path           string              `json:"path,omitempty" yaml:"path,omitempty"`
	CurrentPath    string              `json:"current_path,omitempty" yaml:"current_path,omitempty"`
	TargetPath     string              `json:"target_path,omitempty" yaml:"target_path,omitempty"`
	SourceRevision string              `json:"source_revision" yaml:"source_revision"`
	SourceSHA256   string              `json:"source_sha256" yaml:"source_sha256"`
	TargetState    *TargetPrecondition `json:"target_precondition,omitempty" yaml:"target_precondition,omitempty"`
	Ownership      OwnershipClass      `json:"ownership_disposition,omitempty" yaml:"ownership_disposition,omitempty"`
	Decision       MigrationDecision   `json:"decision" yaml:"decision"`
	Confidence     string              `json:"confidence" yaml:"confidence"`
	Ambiguity      []string            `json:"ambiguity" yaml:"ambiguity"`
	Evidence       []string            `json:"evidence" yaml:"evidence"`
	LegacyAliases  []string            `json:"legacy_aliases" yaml:"legacy_aliases"`
	Conflicts      []string            `json:"conflicts" yaml:"conflicts"`
	Deferred       bool                `json:"deferred" yaml:"deferred"`
	DeferralReason string              `json:"deferral_reason,omitempty" yaml:"deferral_reason,omitempty"`
	BranchOwner    string              `json:"branch_owner" yaml:"branch_owner"`
}

type TargetPrecondition struct {
	State  string `json:"state" yaml:"state"`
	SHA256 string `json:"sha256,omitempty" yaml:"sha256,omitempty"`
}

type MigrationLedger struct {
	SchemaVersion      int               `json:"schema_version" yaml:"schema_version"`
	RepositoryID       string            `json:"repository_id" yaml:"repository_id"`
	DiscoveryVersion   DiscoveryVersion  `json:"discovery_version,omitempty" yaml:"discovery_version,omitempty"`
	TargetCatalogMode  CatalogMode       `json:"target_catalog_mode,omitempty" yaml:"target_catalog_mode,omitempty"`
	PolicyDigest       string            `json:"policy_digest,omitempty" yaml:"policy_digest,omitempty"`
	CandidateSetDigest string            `json:"candidate_set_digest,omitempty" yaml:"candidate_set_digest,omitempty"`
	InventoryRevision  string            `json:"inventory_revision" yaml:"inventory_revision"`
	ReviewedAt         string            `json:"reviewed_at" yaml:"reviewed_at"`
	Records            []MigrationRecord `json:"records" yaml:"records"`
}

type MigrationSkip struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Path        string `json:"path"`
	Remediation string `json:"remediation,omitempty"`
}

type MigrationPlan struct {
	Ledger          MigrationLedger
	FilePlan        reconcile.Plan
	Skips           []MigrationSkip
	Problems        []Problem
	Projection      MigrationProjection
	AuthorityDigest string
}

type MigrationProjection struct {
	Candidates    int  `json:"candidates"`
	Canonical     int  `json:"canonical"`
	Grandfathered int  `json:"grandfathered"`
	Gaps          int  `json:"gaps"`
	Ready         bool `json:"ready"`
}

type MigrationOutcome struct {
	SchemaVersion       int                  `json:"schema_version"`
	DiscoveryVersion    DiscoveryVersion     `json:"discovery_version,omitempty"`
	TargetCatalogMode   CatalogMode          `json:"target_catalog_mode,omitempty"`
	Projection          *MigrationProjection `json:"projected_coverage,omitempty"`
	ActivationPerformed *bool                `json:"activation_performed,omitempty"`
	Result              reconcile.Result     `json:"result"`
	Skips               []MigrationSkip      `json:"skips"`
}

type MigrationApplyOptions struct {
	DryRun bool
	Now    time.Time
	Hook   reconcile.PhaseHook
}

func LoadMigrationLedger(data []byte) (MigrationLedger, error) {
	value, err := state.DecodeYAMLMap(data)
	if err != nil {
		return MigrationLedger{}, fmt.Errorf("invalid_ledger: %w", err)
	}
	schemaPath, err := state.SchemaForPlanMigrationVersion(state.AsInt(value["schema_version"]))
	if err != nil {
		return MigrationLedger{}, fmt.Errorf("invalid_ledger: %w", err)
	}
	if err := state.Validate(schemaPath, value); err != nil {
		return MigrationLedger{}, fmt.Errorf("invalid_ledger: %w", err)
	}
	var ledger MigrationLedger
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&ledger); err != nil {
		return MigrationLedger{}, fmt.Errorf("invalid_ledger: %w", err)
	}
	for index := range ledger.Records {
		if ledger.Records[index].CurrentPath == "" {
			ledger.Records[index].CurrentPath = ledger.Records[index].Path
		}
		if ledger.Records[index].TargetPath == "" {
			ledger.Records[index].TargetPath = ledger.Records[index].CurrentPath
		}
	}
	return ledger, nil
}

func BuildMigrationPlan(root string, ledger MigrationLedger) (MigrationPlan, error) {
	if ledger.SchemaVersion == 2 {
		return buildV2MigrationPlan(root, ledger)
	}
	return buildV1MigrationPlan(root, ledger)
}

func buildV1MigrationPlan(root string, ledger MigrationLedger) (MigrationPlan, error) {
	snapshot, err := LoadRepositorySnapshot(root)
	if err != nil {
		return MigrationPlan{}, err
	}
	settings := snapshot.Settings
	problems := []Problem{}
	if settings.DiscoveryVersion == DiscoveryV2 {
		problems = append(problems, Problem{Code: "migration_ledger_version_incompatible", Message: "schema-v1 migration ledgers are historical evidence and cannot mutate a discovery-v2 catalog", Path: state.ConfigPath, Severity: SeverityError, Remediation: "generate and review a schema-v2 inventory and migration ledger"})
	}
	for _, problem := range snapshot.Problems {
		if problem.Severity != SeverityError || problem.Code == "metadata_missing" || problem.Code == "mixed_new_plan_metadata_missing" || problem.Code == "mixed_grandfathered_plan_modified" {
			continue
		}
		problems = append(problems, problem)
	}
	if settings.Mode == ModeLegacy {
		problems = append(problems, Problem{Code: "catalog_mode_legacy", Message: "migration writes require mixed or canonical catalog mode", Path: state.ConfigPath, Severity: SeverityError, Remediation: "review the inventory, then switch the repository to mixed mode before applying metadata"})
	}
	if settings.RepositoryID != "" && ledger.RepositoryID != settings.RepositoryID {
		problems = append(problems, Problem{Code: "repository_identity_mismatch", Message: fmt.Sprintf("ledger repository %q does not match configured repository %q", ledger.RepositoryID, settings.RepositoryID), Severity: SeverityError})
	}
	if _, err := time.Parse(time.RFC3339, ledger.ReviewedAt); err != nil || !strings.HasSuffix(ledger.ReviewedAt, "Z") {
		problems = append(problems, Problem{Code: "ledger_review_timestamp_invalid", Message: "reviewed_at must be a UTC RFC3339 timestamp ending in Z", Severity: SeverityError})
	}

	observed, err := state.Inspect(root)
	if err != nil {
		return MigrationPlan{}, err
	}
	touches, touchProblems := activeBranchTouches(root)
	problems = append(problems, touchProblems...)
	existingIDs := map[string]string{}
	for _, record := range snapshot.Records {
		if record.Metadata != nil {
			existingIDs[record.Metadata.ID] = record.Path
		}
	}
	decisionIDs := map[string]string{}
	actions := []reconcile.Action{}
	skips := []MigrationSkip{}
	for _, item := range ledger.Records {
		metadata := metadataFromDecision(item)
		if previous, exists := decisionIDs[metadata.ID]; exists && previous != item.Path {
			problems = append(problems, Problem{Code: "duplicate_plan_id", Message: fmt.Sprintf("reviewed ledger assigns %q to both %s and %s", metadata.ID, previous, item.Path), PlanID: metadata.ID, Severity: SeverityError})
			continue
		}
		decisionIDs[metadata.ID] = item.Path
		if existing, exists := existingIDs[metadata.ID]; exists && existing != item.Path {
			skips = append(skips, MigrationSkip{Code: "duplicate_plan_id", Message: fmt.Sprintf("semantic ID already belongs to %s", existing), Path: item.Path, Remediation: "review the identity conflict without overwriting either plan"})
			continue
		}
		if skip := semanticReviewSkip(item); skip != nil {
			skips = append(skips, *skip)
			continue
		}
		if item.SourceRevision != ledger.InventoryRevision {
			skips = append(skips, MigrationSkip{Code: "inventory_revision_mismatch", Message: "record source revision does not match the reviewed inventory revision", Path: item.Path, Remediation: "regenerate one coherent inventory and semantic ledger"})
			continue
		}
		if err := reconcile.ValidateFileTarget(root, item.Path); err != nil {
			skips = append(skips, MigrationSkip{Code: "unsafe_target", Message: err.Error(), Path: item.Path, Remediation: "select a contained non-symlink canonical plan path"})
			continue
		}
		expectedKind, placementOK := migrationPathKind(root, item.Path)
		if !placementOK || expectedKind != metadata.Kind {
			skips = append(skips, MigrationSkip{Code: "invalid_target_placement", Message: fmt.Sprintf("target does not qualify as a canonical %s record", metadata.Kind), Path: item.Path, Remediation: "select the canonical discovery, implementation, or qualifying family path"})
			continue
		}
		identityProblems := ValidateIdentity(metadata.ID, settings.RepositoryID, metadata.Kind)
		if len(identityProblems) > 0 {
			skips = append(skips, MigrationSkip{Code: identityProblems[0].Code, Message: identityProblems[0].Message, Path: item.Path, Remediation: identityProblems[0].Remediation})
			continue
		}
		current, readErr := readRegularSource(root, item.Path)
		if readErr != nil {
			code := "source_unreadable"
			if ErrorCode(readErr) == "source_unsafe" {
				code = "unsafe_target"
			}
			skips = append(skips, MigrationSkip{Code: code, Message: readErr.Error(), Path: item.Path})
			continue
		}
		if existing, parseErr := ParseDocument(item.Path, current, expectedKind); parseErr == nil && existing.Metadata != nil {
			if metadataEqual(*existing.Metadata, metadata) {
				skips = append(skips, MigrationSkip{Code: "already_applied", Message: "canonical metadata already matches the reviewed decision", Path: item.Path})
				continue
			}
			skips = append(skips, MigrationSkip{Code: "metadata_conflict", Message: "target already contains different canonical metadata", Path: item.Path, Remediation: "reevaluate the latest metadata instead of replacing it"})
			continue
		} else if parseErr != nil && ErrorCode(parseErr) != "metadata_missing" {
			skips = append(skips, MigrationSkip{Code: "malformed_metadata", Message: parseErr.Error(), Path: item.Path, Remediation: "repair malformed metadata explicitly; never infer it as legacy absence"})
			continue
		}
		currentDigest := sha256.Sum256(current)
		currentSHA256 := hex.EncodeToString(currentDigest[:])
		checkoutBytesMatchReviewedSource := currentSHA256 == item.SourceSHA256
		dirty, dirtyErr := gitPathDirty(root, item.Path)
		if dirtyErr != nil {
			skips = append(skips, MigrationSkip{Code: "dirty_check_failed", Message: dirtyErr.Error(), Path: item.Path})
			continue
		}
		if dirty {
			if !checkoutBytesMatchReviewedSource {
				skips = append(skips, MigrationSkip{Code: "source_mismatch", Message: "current plan bytes no longer match the reviewed source", Path: item.Path, Remediation: "inventory and semantically reevaluate the latest plan"})
			} else {
				skips = append(skips, MigrationSkip{Code: "dirty_plan_overlap", Message: "target plan has uncommitted changes", Path: item.Path, Remediation: "have the current branch owner migrate or stabilize the plan first"})
			}
			continue
		}
		if refs := touches[item.Path]; len(refs) > 0 {
			skips = append(skips, MigrationSkip{Code: "active_branch_ownership", Message: "target plan is changed on unmerged refs: " + strings.Join(refs, ", "), Path: item.Path, Remediation: "assign migration to the branch owner or wait for reconciliation"})
			continue
		}
		sourceBytes, gitErr := gitBytes(root, "show", item.SourceRevision+":"+item.Path)
		if gitErr != nil {
			skips = append(skips, MigrationSkip{Code: "source_revision_unavailable", Message: gitErr.Error(), Path: item.Path, Remediation: "retain the source commit and rerun inventory"})
			continue
		}
		sourceDigest := sha256.Sum256(sourceBytes)
		if hex.EncodeToString(sourceDigest[:]) != item.SourceSHA256 {
			skips = append(skips, MigrationSkip{Code: "source_revision_mismatch", Message: "reviewed source revision does not contain the reviewed plan blob", Path: item.Path, Remediation: "regenerate inventory from an exact committed source"})
			continue
		}
		if !checkoutBytesMatchReviewedSource && !bytes.Equal(bytes.ReplaceAll(current, []byte("\r\n"), []byte("\n")), bytes.ReplaceAll(sourceBytes, []byte("\r\n"), []byte("\n"))) {
			skips = append(skips, MigrationSkip{Code: "source_mismatch", Message: "clean checkout bytes differ from the reviewed source beyond line-ending conversion", Path: item.Path, Remediation: "remove checkout filters or regenerate inventory from the exact working bytes"})
			continue
		}
		sourceBlob, gitErr := gitText(root, "rev-parse", "--verify", item.SourceRevision+":"+item.Path)
		if gitErr != nil {
			skips = append(skips, MigrationSkip{Code: "source_revision_unavailable", Message: gitErr.Error(), Path: item.Path, Remediation: "retain the source commit and rerun inventory"})
			continue
		}
		indexBlob, gitErr := gitText(root, "rev-parse", "--verify", ":"+item.Path)
		if gitErr != nil {
			skips = append(skips, MigrationSkip{Code: "source_revision_unavailable", Message: gitErr.Error(), Path: item.Path, Remediation: "restore the reviewed plan to the Git index and rerun inventory"})
			continue
		}
		if sourceBlob != indexBlob {
			skips = append(skips, MigrationSkip{Code: "source_revision_mismatch", Message: "reviewed source revision does not contain the indexed plan blob", Path: item.Path, Remediation: "regenerate inventory from an exact committed source"})
			continue
		}
		updated, insertErr := InsertMetadata(current, metadata)
		if insertErr != nil {
			skips = append(skips, MigrationSkip{Code: ErrorCode(insertErr), Message: insertErr.Error(), Path: item.Path})
			continue
		}
		parsed, parseErr := ParseDocument(item.Path, updated, expectedKind)
		if parseErr != nil {
			skips = append(skips, MigrationSkip{Code: ErrorCode(parseErr), Message: parseErr.Error(), Path: item.Path})
			continue
		}
		validation := ValidateRecords([]Record{parsed}, settings.Mode, settings.RepositoryID)
		if HasErrors(validation) {
			first := firstError(validation)
			skips = append(skips, MigrationSkip{Code: first.Code, Message: first.Message, Path: item.Path, Remediation: first.Remediation})
			continue
		}
		actions = append(actions, reconcile.Action{Kind: "replace", Target: item.Path, Owner: "repo-plan", Content: updated, Mode: 0o644, ExpectedSHA256: currentSHA256})
	}

	filePlan, err := reconcile.BuildFilePlan("plans migrate", root, string(observed.Classification), []state.Classification{observed.Classification}, actions)
	if err != nil {
		return MigrationPlan{}, err
	}
	if settings.Mode == ModeCanonical {
		covered := map[string]bool{}
		for _, record := range snapshot.Records {
			if record.Metadata != nil {
				covered[record.Path] = true
			}
		}
		for _, action := range actions {
			covered[action.Target] = true
		}
		candidates, enumerateErr := Enumerate(root)
		if enumerateErr != nil {
			return MigrationPlan{}, enumerateErr
		}
		for _, candidate := range candidates {
			if covered[candidate.Path] {
				continue
			}
			problems = append(problems, Problem{Code: "incomplete_coverage", Message: "canonical migration would leave a formal plan without valid metadata", Path: candidate.Path, Severity: SeverityError, Remediation: "review and include the plan in this coherent migration or defer canonical cutover"})
		}
	}
	actionPaths := map[string]bool{}
	for _, action := range actions {
		actionPaths[action.Target] = true
	}
	for _, problem := range snapshot.Problems {
		if problem.Severity == SeverityError && (problem.Code == "mixed_new_plan_metadata_missing" || problem.Code == "mixed_grandfathered_plan_modified") && !actionPaths[problem.Path] {
			problems = append(problems, problem)
		}
	}
	for _, problem := range problems {
		if problem.Severity != SeverityError {
			continue
		}
		filePlan.Blockers = append(filePlan.Blockers, reconcile.Blocker{Code: problem.Code, Message: problem.Message, Path: problem.Path, Remediation: problem.Remediation, RetryCommand: "plans migrate --dry-run"})
	}
	sort.SliceStable(skips, func(i, j int) bool {
		if skips[i].Path != skips[j].Path {
			return skips[i].Path < skips[j].Path
		}
		return skips[i].Code < skips[j].Code
	})
	SortProblems(problems)
	return MigrationPlan{Ledger: ledger, FilePlan: filePlan, Skips: skips, Problems: problems}, nil
}

func buildV2MigrationPlan(root string, ledger MigrationLedger) (MigrationPlan, error) {
	configuredSnapshot, err := LoadRepositorySnapshot(root)
	if err != nil {
		return MigrationPlan{}, err
	}
	options := SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, TargetCatalogMode: ledger.TargetCatalogMode}
	inventory, err := BuildInventoryWithOptions(root, time.Time{}, options)
	if err != nil {
		return MigrationPlan{}, err
	}
	observed, err := state.Inspect(root)
	if err != nil {
		return MigrationPlan{}, err
	}
	problems := []Problem{}
	skips := []MigrationSkip{}
	actions := []reconcile.Action{}
	projection := MigrationProjection{}
	for _, problem := range configuredSnapshot.Problems {
		if problem.Severity == SeverityError && mixedCutoverInvariantProblem(problem.Code) {
			problems = append(problems, problem)
		}
	}

	if ledger.DiscoveryVersion != DiscoveryV2 {
		problems = append(problems, Problem{Code: "migration_discovery_version_invalid", Message: "schema-v2 migration requires discovery_version 2", Severity: SeverityError})
	}
	if ledger.TargetCatalogMode != ModeMixed && ledger.TargetCatalogMode != ModeCanonical {
		problems = append(problems, Problem{Code: "migration_target_mode_invalid", Message: "schema-v2 migration target mode must be mixed or canonical", Severity: SeverityError})
	}
	if ledger.RepositoryID != inventory.RepositoryID {
		problems = append(problems, Problem{Code: "repository_identity_mismatch", Message: fmt.Sprintf("ledger repository %q does not match configured repository %q", ledger.RepositoryID, inventory.RepositoryID), Severity: SeverityError})
	}
	if _, parseErr := time.Parse(time.RFC3339, ledger.ReviewedAt); parseErr != nil || !strings.HasSuffix(ledger.ReviewedAt, "Z") {
		problems = append(problems, Problem{Code: "ledger_review_timestamp_invalid", Message: "reviewed_at must be a UTC RFC3339 timestamp ending in Z", Severity: SeverityError})
	}
	if ledger.PolicyDigest != inventory.PolicyDigest {
		problems = append(problems, Problem{Code: "migration_policy_mismatch", Message: "reviewed policy digest does not match the prospective discovery-v2 policy", Path: state.ConfigPath, Severity: SeverityError, Remediation: "regenerate inventory and review the current exclusions and ownership policy"})
	}
	authorityDigest, authorityErr := migrationAuthorityDigest(root, configuredSnapshot, inventory)
	if authorityErr != nil {
		return MigrationPlan{}, authorityErr
	}

	ownedCandidates := map[string]Candidate{}
	for _, candidate := range inventory.Candidates {
		if candidate.Ownership == OwnershipOwned {
			ownedCandidates[candidate.Path] = candidate
		}
	}
	projection.Candidates = len(ownedCandidates)
	finalApplied := v2LedgerAlreadyApplied(root, ledger, inventory, ownedCandidates)
	if !finalApplied {
		for _, problem := range inventory.Problems {
			if problem.Severity == SeverityError && !v2MigrationRemediableProblem(problem.Code) {
				problems = append(problems, problem)
			}
		}
		if inventory.Complete == nil || !*inventory.Complete {
			problems = append(problems, Problem{Code: "migration_inventory_incomplete", Message: "prospective discovery-v2 inventory is incomplete", Severity: SeverityError, Remediation: "resolve unsafe, ambiguous, or unavailable candidate evidence before migration"})
		}
	}

	if finalApplied {
		for _, item := range ledger.Records {
			skips = append(skips, MigrationSkip{Code: "already_applied", Message: "canonical target already matches the reviewed decision", Path: item.TargetPath})
		}
		projection.Candidates = len(ledger.Records)
		projection.Canonical = len(ledger.Records)
		projection.Ready = len(problems) == 0
		filePlan, buildErr := reconcile.BuildFilePlan("plans migrate", root, string(observed.Classification), []state.Classification{observed.Classification}, nil)
		if buildErr != nil {
			return MigrationPlan{}, buildErr
		}
		appendMigrationBlockers(&filePlan, problems, nil)
		SortProblems(problems)
		sortMigrationSkips(skips)
		return MigrationPlan{Ledger: ledger, FilePlan: filePlan, Skips: skips, Problems: problems, Projection: projection, AuthorityDigest: authorityDigest}, nil
	}

	if ledger.InventoryRevision != inventory.SourceRevision {
		problems = append(problems, Problem{Code: "inventory_revision_mismatch", Message: "ledger inventory revision does not match the current repository revision", Severity: SeverityError, Remediation: "regenerate one coherent inventory and semantic ledger"})
	}
	if ledger.CandidateSetDigest != inventory.CandidateSetDigest {
		problems = append(problems, Problem{Code: "migration_candidate_set_mismatch", Message: "reviewed candidate-set digest does not match the prospective discovery-v2 inventory", Severity: SeverityError, Remediation: "regenerate inventory and review every current candidate"})
	}

	inventoryRecords := map[string]InventoryRecord{}
	for _, record := range inventory.Records {
		if record.Authoritative != nil && *record.Authoritative {
			inventoryRecords[record.Path] = record
		}
	}
	currentPaths := map[string]bool{}
	targetPaths := map[string]bool{}
	decisionIDs := map[string]string{}
	projectedRecords := []Record{}
	for _, item := range ledger.Records {
		currentPath := item.CurrentPath
		targetPath := item.TargetPath
		if currentPath == "" {
			currentPath = item.Path
		}
		if targetPath == "" {
			targetPath = currentPath
		}
		item.CurrentPath = currentPath
		item.TargetPath = targetPath
		if protectedMigrationPath(currentPath) || protectedMigrationPath(targetPath) {
			problems = append(problems, Problem{Code: "migration_protected_path", Message: "migration must never create, replace, rename, or remove catalog configuration or the frozen register", Path: currentPath, Severity: SeverityError})
			continue
		}
		if currentPaths[currentPath] {
			problems = append(problems, Problem{Code: "migration_current_path_duplicate", Message: "ledger contains the current path more than once", Path: currentPath, Severity: SeverityError})
			continue
		}
		currentPaths[currentPath] = true
		if targetPaths[targetPath] {
			problems = append(problems, Problem{Code: "migration_target_path_duplicate", Message: "ledger assigns more than one record to the target path", Path: targetPath, Severity: SeverityError})
			continue
		}
		targetPaths[targetPath] = true

		metadata := metadataFromDecision(item)
		if prior, exists := decisionIDs[metadata.ID]; exists {
			problems = append(problems, Problem{Code: "duplicate_plan_id", Message: fmt.Sprintf("reviewed ledger assigns %q to both %s and %s", metadata.ID, prior, currentPath), PlanID: metadata.ID, Severity: SeverityError})
			continue
		}
		decisionIDs[metadata.ID] = currentPath
		candidate, exists := ownedCandidates[currentPath]
		if !exists {
			problems = append(problems, Problem{Code: "migration_candidate_missing", Message: "ledger current path is not an authoritative discovery-v2 candidate", Path: currentPath, Severity: SeverityError, Remediation: "review the current prospective inventory and do not migrate excluded or unowned paths"})
			continue
		}
		if item.Ownership != OwnershipOwned {
			problems = append(problems, Problem{Code: "migration_ownership_invalid", Message: "migration record must retain reviewed owned disposition", Path: currentPath, Severity: SeverityError})
			continue
		}
		if item.SourceRevision != ledger.InventoryRevision {
			problems = append(problems, Problem{Code: "inventory_revision_mismatch", Message: "record source revision does not match the reviewed inventory revision", Path: currentPath, Severity: SeverityError})
			continue
		}
		if item.BranchOwner != "none" {
			problems = append(problems, Problem{Code: "active_branch_ownership", Message: "reviewed ledger assigns the plan to branch owner " + item.BranchOwner, Path: currentPath, Severity: SeverityError})
			continue
		}
		if len(item.Ambiguity) > 0 || len(item.Conflicts) > 0 || item.Confidence == "low" {
			problems = append(problems, Problem{Code: "ambiguous_legacy_evidence", Message: "semantic review retains ambiguity, conflicts, or low confidence", Path: currentPath, Severity: SeverityError})
			continue
		}
		if err := reconcile.ValidateFileTarget(root, currentPath); err != nil {
			problems = append(problems, Problem{Code: "unsafe_target", Message: err.Error(), Path: currentPath, Severity: SeverityError})
			continue
		}
		if err := reconcile.ValidateFileTarget(root, targetPath); err != nil {
			problems = append(problems, Problem{Code: "unsafe_target", Message: err.Error(), Path: targetPath, Severity: SeverityError})
			continue
		}
		targetKind, targetOK := migrationV2TargetKind(targetPath, metadata)
		if !targetOK || targetKind != metadata.Kind {
			problems = append(problems, Problem{Code: "invalid_target_placement", Message: fmt.Sprintf("target does not qualify as a canonical %s record", metadata.Kind), Path: targetPath, Severity: SeverityError})
			continue
		}
		if identityProblems := ValidateIdentity(metadata.ID, inventory.RepositoryID, metadata.Kind); len(identityProblems) > 0 {
			problem := identityProblems[0]
			problem.Path = currentPath
			problems = append(problems, problem)
			continue
		}
		current, readErr := readRegularSource(root, currentPath)
		if readErr != nil {
			problems = append(problems, Problem{Code: "source_unreadable", Message: readErr.Error(), Path: currentPath, Severity: SeverityError})
			continue
		}
		currentSHA := sha256Text(current)
		if currentSHA != item.SourceSHA256 || candidate.Provenance.Source.ContentSHA256 != item.SourceSHA256 {
			problems = append(problems, Problem{Code: "source_mismatch", Message: "current candidate bytes no longer match the reviewed source hash", Path: currentPath, Severity: SeverityError})
			continue
		}
		sourceBytes, gitErr := gitBytes(root, "show", item.SourceRevision+":"+currentPath)
		if gitErr != nil || sha256Text(sourceBytes) != item.SourceSHA256 || !bytes.Equal(sourceBytes, current) {
			message := "reviewed source revision, index, and worktree bytes must match exactly"
			if gitErr != nil {
				message = gitErr.Error()
			}
			problems = append(problems, Problem{Code: "source_revision_mismatch", Message: message, Path: currentPath, Severity: SeverityError})
			continue
		}
		inventoryRecord := inventoryRecords[currentPath]
		if inventoryRecord.DirtyOverlap {
			problems = append(problems, Problem{Code: "dirty_plan_overlap", Message: "candidate has uncommitted changes", Path: currentPath, Severity: SeverityError})
			continue
		}
		if len(inventoryRecord.ActiveBranchTouch) > 0 {
			problems = append(problems, Problem{Code: "active_branch_ownership", Message: "candidate is changed on unmerged refs: " + strings.Join(inventoryRecord.ActiveBranchTouch, ", "), Path: currentPath, Severity: SeverityError})
			continue
		}

		if item.Deferred {
			if ledger.TargetCatalogMode != ModeMixed || inventory.ConfiguredCatalogMode != ModeMixed || currentPath != targetPath {
				problems = append(problems, Problem{Code: "mixed_cutover_proof_missing", Message: "deferred filename-only records require an already-active mixed catalog and unchanged same-path cutover proof", Path: currentPath, Severity: SeverityError})
				continue
			}
			baseline, baselineProblems := loadMixedBaseline(root, loadCutoverRevision(root))
			problems = append(problems, baselineProblems...)
			if baseline[currentPath] == "" || baseline[currentPath] != currentSHA {
				problems = append(problems, Problem{Code: "mixed_grandfathered_plan_modified", Message: "deferred plan does not match its exact register-proven cutover blob", Path: currentPath, Severity: SeverityError})
				continue
			}
			if parsed, parseErr := ParseDocument(currentPath, current, candidate.ExpectedKind); parseErr == nil && parsed.Metadata != nil {
				problems = append(problems, Problem{Code: "migration_deferral_invalid", Message: "canonical metadata-bearing records must not be grandfathered as filename-only", Path: currentPath, Severity: SeverityError})
				continue
			}
			skips = append(skips, MigrationSkip{Code: "grandfathered_mixed_record", Message: item.DeferralReason, Path: currentPath})
			projection.Grandfathered++
			continue
		}

		if skip := semanticReviewSkip(item); skip != nil {
			problems = append(problems, Problem{Code: skip.Code, Message: skip.Message, Path: currentPath, Severity: SeverityError, Remediation: skip.Remediation})
			continue
		}
		if preconditionProblem := validateMigrationTargetPrecondition(root, item, currentSHA); preconditionProblem != nil {
			problems = append(problems, *preconditionProblem)
			continue
		}

		currentKind, _ := kindForFilename(path.Base(currentPath))
		existing, parseErr := ParseDocument(currentPath, current, currentKind)
		updated := current
		if parseErr == nil && existing.Metadata != nil {
			if !metadataEqual(*existing.Metadata, metadata) {
				problems = append(problems, Problem{Code: "metadata_conflict", Message: "source contains canonical metadata that differs from the reviewed decision", Path: currentPath, Severity: SeverityError})
				continue
			}
		} else if parseErr != nil && ErrorCode(parseErr) == "metadata_missing" {
			updated, err = InsertMetadata(current, metadata)
			if err != nil {
				problems = append(problems, Problem{Code: ErrorCode(err), Message: err.Error(), Path: currentPath, Severity: SeverityError})
				continue
			}
		} else if parseErr != nil {
			problems = append(problems, Problem{Code: "malformed_metadata", Message: parseErr.Error(), Path: currentPath, Severity: SeverityError})
			continue
		}
		beforeHeader, beforeErr := ParseHeader(current)
		afterHeader, afterErr := ParseHeader(updated)
		if beforeErr != nil || afterErr != nil || beforeHeader.Created != afterHeader.Created || beforeHeader.LastUpdated != afterHeader.LastUpdated || beforeHeader.Lifecycle != afterHeader.Lifecycle {
			problems = append(problems, Problem{Code: "migration_header_changed", Message: "metadata migration must preserve Created, Last updated, and lifecycle header values", Path: currentPath, Severity: SeverityError})
			continue
		}
		projected, parseErr := ParseDocument(targetPath, updated, targetKind)
		if parseErr != nil {
			problems = append(problems, Problem{Code: ErrorCode(parseErr), Message: parseErr.Error(), Path: targetPath, Severity: SeverityError})
			continue
		}
		projected.ExpectedKind = targetKind
		if targetKind == KindFamily {
			projected.FamilyQualified = true
		}
		projectedRecords = append(projectedRecords, projected)
		projection.Canonical++
		fileMode := migrationActionMode(candidate.Provenance.Source.Mode)
		if currentPath == targetPath {
			if bytes.Equal(current, updated) {
				skips = append(skips, MigrationSkip{Code: "already_applied", Message: "canonical metadata already matches the reviewed decision", Path: currentPath})
				continue
			}
			actions = append(actions, reconcile.Action{Kind: "replace", Target: currentPath, Owner: "repo-plan", Content: updated, Mode: fileMode, ExpectedSHA256: currentSHA})
		} else {
			actions = append(actions,
				reconcile.Action{Kind: "create", Target: targetPath, Owner: "repo-plan", Content: updated, Mode: fileMode},
				reconcile.Action{Kind: "remove", Target: currentPath, Owner: "repo-plan", ExpectedSHA256: currentSHA},
			)
		}
	}

	for candidatePath := range ownedCandidates {
		if !currentPaths[candidatePath] {
			problems = append(problems, Problem{Code: "migration_candidate_unreviewed", Message: "authoritative candidate is absent from the reviewed migration ledger", Path: candidatePath, Severity: SeverityError, Remediation: "review every discovered candidate in one coherent ledger"})
		}
	}
	projection.Gaps = projection.Candidates - projection.Canonical - projection.Grandfathered
	if projection.Gaps < 0 {
		projection.Gaps = 0
	}
	for _, problem := range ValidateRecordsForDiscovery(projectedRecords, ledger.TargetCatalogMode, inventory.RepositoryID, DiscoveryV2) {
		if problem.Severity == SeverityError {
			problems = append(problems, problem)
		}
	}
	if ledger.TargetCatalogMode == ModeCanonical && projection.Grandfathered > 0 {
		problems = append(problems, Problem{Code: "incomplete_coverage", Message: "canonical migration cannot retain grandfathered filename-only records", Severity: SeverityError})
	}
	if projection.Gaps > 0 {
		problems = append(problems, Problem{Code: "incomplete_coverage", Message: fmt.Sprintf("projected migration leaves %d authoritative candidate(s) without canonical metadata or reviewed mixed grandfathering", projection.Gaps), Severity: SeverityError})
	}
	projection.Ready = !HasErrors(problems) && projection.Gaps == 0

	filePlan, err := reconcile.BuildFilePlan("plans migrate", root, string(observed.Classification), []state.Classification{observed.Classification}, actions)
	if err != nil {
		return MigrationPlan{}, err
	}
	appendMigrationBlockers(&filePlan, problems, skips)
	SortProblems(problems)
	sortMigrationSkips(skips)
	return MigrationPlan{Ledger: ledger, FilePlan: filePlan, Skips: skips, Problems: problems, Projection: projection, AuthorityDigest: authorityDigest}, nil
}

func v2MigrationRemediableProblem(code string) bool {
	switch code {
	case "metadata_missing", "canonical_filename_missing", "record_kind_path_mismatch", "family_placement_invalid", "mixed_new_plan_metadata_missing", "mixed_grandfathered_plan_modified":
		return true
	default:
		return false
	}
}

func mixedCutoverInvariantProblem(code string) bool {
	return strings.HasPrefix(code, "mixed_cutover_") || code == "legacy_register_modified_after_cutover" || code == "legacy_register_missing_after_cutover"
}

func v2LedgerAlreadyApplied(root string, ledger MigrationLedger, inventory Inventory, owned map[string]Candidate) bool {
	if ledger.TargetCatalogMode != ModeCanonical {
		return false
	}
	settings, settingProblems := LoadRepositorySettings(root)
	if HasErrors(settingProblems) || settings.RepositoryID != inventory.RepositoryID {
		return false
	}
	settings.DiscoveryVersion = DiscoveryV2
	settings.Mode = ledger.TargetCatalogMode
	settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	if settings.PolicyDigest != ledger.PolicyDigest {
		return false
	}
	classification, err := ClassifyCommitTree(root, ledger.InventoryRevision, settings, ledger.TargetCatalogMode, inventory.RepositoryID)
	if err != nil || !classification.Complete || classification.CandidateSetDigest != ledger.CandidateSetDigest {
		return false
	}
	originalOwned := map[string]Candidate{}
	for _, candidate := range classification.Candidates {
		if candidate.Ownership == OwnershipOwned {
			originalOwned[candidate.Path] = candidate
		}
	}
	if len(originalOwned) != len(ledger.Records) {
		return false
	}
	seenCurrent := map[string]bool{}
	seenTarget := map[string]bool{}
	for _, item := range ledger.Records {
		if !migrationRecordCanBeFinal(item) || seenCurrent[item.CurrentPath] || seenTarget[item.TargetPath] {
			return false
		}
		seenCurrent[item.CurrentPath] = true
		seenTarget[item.TargetPath] = true
		candidate, exists := originalOwned[item.CurrentPath]
		if !exists || !finalMigrationTargetMatchesReviewedSource(root, ledger, item, candidate, inventory.RepositoryID) {
			return false
		}
	}

	head, err := gitText(root, "rev-parse", "--verify", "HEAD")
	if err == nil && head == ledger.InventoryRevision {
		staged, stagedErr := gitBytes(root, "diff", "--cached", "--name-only", "-z", "HEAD")
		if stagedErr == nil && len(staged) == 0 {
			return true
		}
	}
	if inventory.Complete == nil || !*inventory.Complete || HasErrors(inventory.Problems) || len(owned) != len(seenTarget) {
		return false
	}
	for target := range seenTarget {
		if _, exists := owned[target]; !exists {
			return false
		}
	}
	return true
}

func migrationRecordCanBeFinal(item MigrationRecord) bool {
	if item.Deferred || item.CurrentPath == "" || item.TargetPath == "" || item.Ownership != OwnershipOwned || item.BranchOwner != "none" || item.Confidence == "low" || len(item.Ambiguity) > 0 || len(item.Conflicts) > 0 || protectedMigrationPath(item.CurrentPath) || protectedMigrationPath(item.TargetPath) {
		return false
	}
	return true
}

func finalMigrationTargetMatchesReviewedSource(root string, ledger MigrationLedger, item MigrationRecord, candidate Candidate, repositoryID string) bool {
	if item.SourceRevision != ledger.InventoryRevision || item.SourceSHA256 == "" || item.SourceSHA256 != candidate.Provenance.Source.ContentSHA256 {
		return false
	}
	if item.CurrentPath == item.TargetPath {
		if item.TargetState == nil || (item.TargetState.State == "same-path" && item.TargetState.SHA256 != "") || (item.TargetState.State != "same-path" && (item.TargetState.State != "exact" || item.TargetState.SHA256 != item.SourceSHA256)) {
			return false
		}
	} else if item.TargetState == nil || item.TargetState.State != "absent" || item.TargetState.SHA256 != "" {
		return false
	}
	if item.CurrentPath != item.TargetPath {
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(item.CurrentPath))); !os.IsNotExist(err) {
			return false
		}
	}
	metadata := metadataFromDecision(item)
	targetKind, targetOK := migrationV2TargetKind(item.TargetPath, metadata)
	if !targetOK || metadata.Kind != targetKind || len(ValidateIdentity(metadata.ID, repositoryID, metadata.Kind)) > 0 {
		return false
	}
	source, err := gitBytes(root, "show", item.SourceRevision+":"+item.CurrentPath)
	if err != nil || sha256Text(source) != item.SourceSHA256 {
		return false
	}
	sourceRecord, parseErr := ParseDocument(item.CurrentPath, source, candidate.ExpectedKind)
	expected := source
	if parseErr == nil && sourceRecord.Metadata != nil {
		if !metadataEqual(*sourceRecord.Metadata, metadata) {
			return false
		}
	} else if parseErr != nil && ErrorCode(parseErr) == "metadata_missing" {
		expected, err = InsertMetadata(source, metadata)
		if err != nil {
			return false
		}
	} else if parseErr != nil {
		return false
	}
	actual, err := readRegularSource(root, item.TargetPath)
	if err != nil || !bytes.Equal(actual, expected) {
		return false
	}
	info, err := os.Stat(filepath.Join(root, filepath.FromSlash(item.TargetPath)))
	if err != nil || !info.Mode().IsRegular() || !gitModeMatchesPermissions(candidate.Provenance.Source.Mode, info.Mode().Perm(), runtime.GOOS == "windows") {
		return false
	}
	record, parseErr := ParseDocument(item.TargetPath, actual, targetKind)
	return parseErr == nil && record.Metadata != nil && metadataEqual(*record.Metadata, metadata)
}

func protectedMigrationPath(value string) bool {
	return value == LegacyRegisterPath || value == state.ConfigPath
}

func migrationActionMode(mode GitMode) uint32 {
	if mode == GitModeExecutable {
		return 0o755
	}
	return 0o644
}

func gitModeMatchesPermissions(mode GitMode, permissions os.FileMode, windows bool) bool {
	if windows {
		return mode == GitModeRegular || mode == GitModeExecutable
	}
	executable := permissions&0o111 != 0
	if mode == GitModeExecutable {
		return executable
	}
	return mode == GitModeRegular && !executable
}

type migrationAuthorityRecord struct {
	Path                string   `json:"path"`
	SourceSHA256        string   `json:"source_sha256"`
	DirtyOverlap        bool     `json:"dirty_overlap"`
	ActiveBranchTouches []string `json:"active_branch_touches"`
}

type migrationAuthorityEvidence struct {
	ConfiguredSettings RepositorySettings         `json:"configured_settings"`
	ConfiguredProblems []Problem                  `json:"configured_problems"`
	InventoryRevision  string                     `json:"inventory_revision"`
	PolicyDigest       string                     `json:"policy_digest"`
	CandidateSetDigest string                     `json:"candidate_set_digest"`
	Complete           bool                       `json:"complete"`
	Candidates         []Candidate                `json:"candidates"`
	Problems           []Problem                  `json:"problems"`
	Records            []migrationAuthorityRecord `json:"records"`
	ConfigSHA256       string                     `json:"config_sha256"`
	RegisterSHA256     string                     `json:"register_sha256"`
}

func migrationAuthorityDigest(root string, configured RepositorySnapshot, inventory Inventory) (string, error) {
	configSHA, err := regularSourceDigest(root, state.ConfigPath)
	if err != nil {
		return "", err
	}
	registerSHA, err := regularSourceDigest(root, LegacyRegisterPath)
	if err != nil {
		return "", err
	}
	complete := inventory.Complete != nil && *inventory.Complete
	evidence := migrationAuthorityEvidence{
		ConfiguredSettings: configured.Settings,
		ConfiguredProblems: append([]Problem{}, configured.Problems...),
		InventoryRevision:  inventory.SourceRevision,
		PolicyDigest:       inventory.PolicyDigest,
		CandidateSetDigest: inventory.CandidateSetDigest,
		Complete:           complete,
		Candidates:         append([]Candidate{}, inventory.Candidates...),
		Problems:           append([]Problem{}, inventory.Problems...),
		ConfigSHA256:       configSHA,
		RegisterSHA256:     registerSHA,
	}
	for _, record := range inventory.Records {
		if record.Authoritative == nil || !*record.Authoritative {
			continue
		}
		evidence.Records = append(evidence.Records, migrationAuthorityRecord{Path: record.Path, SourceSHA256: record.SourceSHA256, DirtyOverlap: record.DirtyOverlap, ActiveBranchTouches: append([]string{}, record.ActiveBranchTouch...)})
	}
	data, err := json.Marshal(evidence)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

func currentMigrationAuthorityDigest(root string, ledger MigrationLedger) (string, error) {
	configured, err := LoadRepositorySnapshot(root)
	if err != nil {
		return "", err
	}
	inventory, err := BuildInventoryWithOptions(root, time.Time{}, SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, TargetCatalogMode: ledger.TargetCatalogMode})
	if err != nil {
		return "", err
	}
	return migrationAuthorityDigest(root, configured, inventory)
}

func regularSourceDigest(root, relative string) (string, error) {
	data, err := readRegularSource(root, relative)
	if os.IsNotExist(err) {
		return "absent", nil
	}
	if err != nil {
		return "", err
	}
	return sha256Text(data), nil
}

func migrationV2TargetKind(target string, metadata Metadata) (Kind, bool) {
	if validateGitPath(target) != nil || !hasExactDocsSegment(target) || path.Ext(target) != ".md" {
		return "", false
	}
	if kind, ok := kindForFilename(path.Base(target)); ok {
		return kind, true
	}
	if path.Base(target) == "README.md" && metadata.Kind == KindFamily {
		return KindFamily, true
	}
	return "", false
}

func validateMigrationTargetPrecondition(root string, item MigrationRecord, currentSHA string) *Problem {
	if item.TargetState == nil {
		return &Problem{Code: "target_precondition_missing", Message: "schema-v2 migration requires an explicit target precondition", Path: item.TargetPath, Severity: SeverityError}
	}
	if item.CurrentPath == item.TargetPath {
		switch item.TargetState.State {
		case "same-path":
			if item.TargetState.SHA256 != "" {
				return &Problem{Code: "target_precondition_invalid", Message: "same-path precondition must not carry a second digest", Path: item.TargetPath, Severity: SeverityError}
			}
			return nil
		case "exact":
			if item.TargetState.SHA256 == currentSHA {
				return nil
			}
		}
		return &Problem{Code: "target_precondition_mismatch", Message: "same-path target precondition does not match current source bytes", Path: item.TargetPath, Severity: SeverityError}
	}
	if item.TargetState.State != "absent" || item.TargetState.SHA256 != "" {
		return &Problem{Code: "target_precondition_invalid", Message: "filename correction requires an absent no-replace target precondition", Path: item.TargetPath, Severity: SeverityError}
	}
	if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(item.TargetPath))); err == nil {
		return &Problem{Code: "target_precondition_mismatch", Message: "reviewed no-replace target already exists", Path: item.TargetPath, Severity: SeverityError}
	} else if !os.IsNotExist(err) {
		return &Problem{Code: "target_precondition_unavailable", Message: err.Error(), Path: item.TargetPath, Severity: SeverityError}
	}
	return nil
}

func loadCutoverRevision(root string) string {
	settings, _ := LoadRepositorySettings(root)
	return settings.CutoverRevision
}

func appendMigrationBlockers(plan *reconcile.Plan, problems []Problem, skips []MigrationSkip) {
	for _, problem := range problems {
		if problem.Severity == SeverityError {
			plan.Blockers = append(plan.Blockers, reconcile.Blocker{Code: problem.Code, Message: problem.Message, Path: problem.Path, Remediation: problem.Remediation, RetryCommand: "plans migrate --dry-run"})
		}
	}
	for _, skip := range skips {
		if skip.Code == "already_applied" || skip.Code == "grandfathered_mixed_record" {
			continue
		}
		plan.Blockers = append(plan.Blockers, reconcile.Blocker{Code: skip.Code, Message: skip.Message, Path: skip.Path, Remediation: skip.Remediation, RetryCommand: "plans migrate --dry-run"})
	}
}

func sortMigrationSkips(skips []MigrationSkip) {
	sort.SliceStable(skips, func(i, j int) bool {
		if skips[i].Path != skips[j].Path {
			return skips[i].Path < skips[j].Path
		}
		return skips[i].Code < skips[j].Code
	})
}

func ExecuteMigration(plan MigrationPlan, options MigrationApplyOptions) (MigrationOutcome, error) {
	var result reconcile.Result
	var err error
	authorityCheck := func() ([]reconcile.Blocker, error) {
		if plan.Ledger.SchemaVersion != 2 {
			return nil, nil
		}
		current, digestErr := currentMigrationAuthorityDigest(plan.FilePlan.Root, plan.Ledger)
		if digestErr != nil {
			return nil, digestErr
		}
		if current != plan.AuthorityDigest {
			return []reconcile.Blocker{{Code: "migration_authority_drift", Message: "revision, policy, candidate, branch, dirty-state, config, or frozen-register authority changed after migration planning", Remediation: "rebuild and review the migration plan from a fresh prospective inventory", RetryCommand: "plans migrate --dry-run"}}, nil
		}
		return nil, nil
	}
	if options.DryRun {
		filePlan := plan.FilePlan
		blockers, checkErr := authorityCheck()
		if checkErr != nil {
			return MigrationOutcome{}, checkErr
		}
		filePlan.Blockers = append(filePlan.Blockers, blockers...)
		result = reconcile.Preview(filePlan)
	} else {
		result, err = reconcile.Apply(plan.FilePlan, reconcile.ApplyOptions{Now: options.Now, Hook: options.Hook, AuthorityCheck: authorityCheck})
	}
	outcome := MigrationOutcome{SchemaVersion: 1, Result: result, Skips: append([]MigrationSkip{}, plan.Skips...)}
	if plan.Ledger.SchemaVersion == 2 {
		activated := false
		projection := plan.Projection
		outcome.SchemaVersion = 2
		outcome.DiscoveryVersion = plan.Ledger.DiscoveryVersion
		outcome.TargetCatalogMode = plan.Ledger.TargetCatalogMode
		outcome.Projection = &projection
		outcome.ActivationPerformed = &activated
	}
	return outcome, err
}

func InsertMetadata(data []byte, metadata Metadata) ([]byte, error) {
	header, err := ParseHeader(data)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if len(lineIndexesOutsideFences(lines, MetadataBeginMarker)) > 0 || len(lineIndexesOutsideFences(lines, MetadataEndMarker)) > 0 {
		return nil, &CodedError{Code: "metadata_conflict", Err: fmt.Errorf("document already contains a plan metadata marker")}
	}
	value := map[string]any{}
	encoded, _ := json.Marshal(map[string]any{"plan": metadata})
	if err := json.Unmarshal(encoded, &value); err != nil {
		return nil, err
	}
	yamlData, err := state.EncodeYAML(value)
	if err != nil {
		return nil, err
	}
	lineEnding := []byte("\n")
	if bytes.Contains(data, []byte("\r\n")) {
		lineEnding = []byte("\r\n")
	}
	yamlData = bytes.ReplaceAll(yamlData, []byte("\n"), lineEnding)
	block := bytes.Join([][]byte{
		[]byte(MetadataBeginMarker),
		[]byte("```yaml"),
		bytes.TrimSuffix(yamlData, lineEnding),
		[]byte("```"),
		[]byte(MetadataEndMarker),
	}, lineEnding)
	offset := offsetAfterLine(data, header.TitleLine)
	if offset < 0 {
		return nil, &CodedError{Code: "title_missing", Err: fmt.Errorf("cannot locate canonical title line")}
	}
	prefix := append([]byte{}, data[:offset]...)
	suffix := data[offset:]
	if offset == len(data) && (len(prefix) == 0 || (!bytes.HasSuffix(prefix, []byte("\n")) && !bytes.HasSuffix(prefix, []byte("\r")))) {
		prefix = append(prefix, lineEnding...)
	}
	result := append(prefix, lineEnding...)
	result = append(result, block...)
	result = append(result, lineEnding...)
	result = append(result, suffix...)
	return result, nil
}

func metadataFromDecision(record MigrationRecord) Metadata {
	return Metadata{
		SchemaVersion:          1,
		ID:                     record.Decision.ID,
		Kind:                   record.Decision.Kind,
		Purpose:                record.Decision.Purpose,
		FirstCataloged:         record.Decision.FirstCataloged,
		CatalogMetadataUpdated: record.Decision.CatalogMetadataUpdated,
		Family:                 record.Decision.Family,
		Products:               append([]string{}, record.Decision.Products...),
		Capabilities:           append([]string{}, record.Decision.Capabilities...),
		StrategicThemes:        append([]string{}, record.Decision.StrategicThemes...),
		Relations:              append([]Relation{}, record.Decision.Relations...),
		LegacyAliases:          append([]string{}, record.LegacyAliases...),
	}
}

func semanticReviewSkip(record MigrationRecord) *MigrationSkip {
	if record.Deferred {
		return &MigrationSkip{Code: "migration_deferred", Message: record.DeferralReason, Path: record.Path, Remediation: "resolve the recorded deferral with the plan owner"}
	}
	if record.BranchOwner != "none" {
		return &MigrationSkip{Code: "active_branch_ownership", Message: "reviewed ledger assigns the plan to branch owner " + record.BranchOwner, Path: record.Path, Remediation: "let the recorded branch owner apply or refresh metadata"}
	}
	if len(record.Ambiguity) > 0 || len(record.Conflicts) > 0 || record.Confidence == "low" {
		return &MigrationSkip{Code: "ambiguous_legacy_evidence", Message: "semantic review retains ambiguity, conflicts, or low confidence", Path: record.Path, Remediation: "reevaluate the source evidence before migration"}
	}
	return nil
}

func migrationPathKind(root, path string) (Kind, bool) {
	return FormalPathKind(root, path)
}

func metadataEqual(left, right Metadata) bool {
	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)
	return bytes.Equal(leftJSON, rightJSON)
}

func firstError(problems []Problem) Problem {
	for _, problem := range problems {
		if problem.Severity == SeverityError {
			return problem
		}
	}
	return Problem{Code: "migration_validation_failed", Message: "migration validation failed", Severity: SeverityError}
}

func offsetAfterLine(data []byte, oneBasedLine int) int {
	if oneBasedLine < 1 {
		return -1
	}
	line := 1
	for index := 0; index < len(data); index++ {
		if data[index] != '\n' {
			continue
		}
		if line == oneBasedLine {
			return index + 1
		}
		line++
	}
	if line == oneBasedLine {
		return len(data)
	}
	return -1
}

func HasMaterialMigrationSkips(skips []MigrationSkip) bool {
	for _, skip := range skips {
		if skip.Code != "already_applied" && skip.Code != "grandfathered_mixed_record" {
			return true
		}
	}
	return false
}
