package plancatalog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	Path           string            `json:"path" yaml:"path"`
	SourceRevision string            `json:"source_revision" yaml:"source_revision"`
	SourceSHA256   string            `json:"source_sha256" yaml:"source_sha256"`
	Decision       MigrationDecision `json:"decision" yaml:"decision"`
	Confidence     string            `json:"confidence" yaml:"confidence"`
	Ambiguity      []string          `json:"ambiguity" yaml:"ambiguity"`
	Evidence       []string          `json:"evidence" yaml:"evidence"`
	LegacyAliases  []string          `json:"legacy_aliases" yaml:"legacy_aliases"`
	Conflicts      []string          `json:"conflicts" yaml:"conflicts"`
	Deferred       bool              `json:"deferred" yaml:"deferred"`
	DeferralReason string            `json:"deferral_reason,omitempty" yaml:"deferral_reason,omitempty"`
	BranchOwner    string            `json:"branch_owner" yaml:"branch_owner"`
}

type MigrationLedger struct {
	SchemaVersion     int               `json:"schema_version" yaml:"schema_version"`
	RepositoryID      string            `json:"repository_id" yaml:"repository_id"`
	InventoryRevision string            `json:"inventory_revision" yaml:"inventory_revision"`
	ReviewedAt        string            `json:"reviewed_at" yaml:"reviewed_at"`
	Records           []MigrationRecord `json:"records" yaml:"records"`
}

type MigrationSkip struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	Path        string `json:"path"`
	Remediation string `json:"remediation,omitempty"`
}

type MigrationPlan struct {
	Ledger   MigrationLedger
	FilePlan reconcile.Plan
	Skips    []MigrationSkip
	Problems []Problem
}

type MigrationOutcome struct {
	SchemaVersion int              `json:"schema_version"`
	Result        reconcile.Result `json:"result"`
	Skips         []MigrationSkip  `json:"skips"`
}

type MigrationApplyOptions struct {
	DryRun bool
	Now    time.Time
	Hook   reconcile.PhaseHook
}

func LoadMigrationLedger(data []byte) (MigrationLedger, error) {
	if _, err := state.DecodeAndValidateYAML(state.PlanMigrationSchema, data); err != nil {
		return MigrationLedger{}, fmt.Errorf("invalid_ledger: %w", err)
	}
	var ledger MigrationLedger
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&ledger); err != nil {
		return MigrationLedger{}, fmt.Errorf("invalid_ledger: %w", err)
	}
	return ledger, nil
}

func BuildMigrationPlan(root string, ledger MigrationLedger) (MigrationPlan, error) {
	snapshot, err := LoadRepositorySnapshot(root)
	if err != nil {
		return MigrationPlan{}, err
	}
	settings := snapshot.Settings
	problems := []Problem{}
	for _, problem := range snapshot.Problems {
		if problem.Severity != SeverityError || problem.Code == "metadata_missing" || problem.Code == "mixed_new_plan_metadata_missing" {
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
		if hex.EncodeToString(currentDigest[:]) != item.SourceSHA256 {
			skips = append(skips, MigrationSkip{Code: "source_mismatch", Message: "current plan bytes no longer match the reviewed source hash", Path: item.Path, Remediation: "inventory and semantically reevaluate the latest plan"})
			continue
		}
		dirty, dirtyErr := gitPathDirty(root, item.Path)
		if dirtyErr != nil {
			skips = append(skips, MigrationSkip{Code: "dirty_check_failed", Message: dirtyErr.Error(), Path: item.Path})
			continue
		}
		if dirty {
			skips = append(skips, MigrationSkip{Code: "dirty_plan_overlap", Message: "target plan has uncommitted changes", Path: item.Path, Remediation: "have the current branch owner migrate or stabilize the plan first"})
			continue
		}
		if refs := touches[item.Path]; len(refs) > 0 {
			skips = append(skips, MigrationSkip{Code: "active_branch_ownership", Message: "target plan is changed on unmerged refs: " + strings.Join(refs, ", "), Path: item.Path, Remediation: "assign migration to the branch owner or wait for reconciliation"})
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
		actions = append(actions, reconcile.Action{Kind: "replace", Target: item.Path, Owner: "repo-plan", Content: updated, Mode: 0o644, ExpectedSHA256: item.SourceSHA256})
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
		if problem.Severity == SeverityError && problem.Code == "mixed_new_plan_metadata_missing" && !actionPaths[problem.Path] {
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

func ExecuteMigration(plan MigrationPlan, options MigrationApplyOptions) (MigrationOutcome, error) {
	var result reconcile.Result
	var err error
	if options.DryRun {
		result = reconcile.Preview(plan.FilePlan)
	} else {
		result, err = reconcile.Apply(plan.FilePlan, reconcile.ApplyOptions{Now: options.Now, Hook: options.Hook})
	}
	return MigrationOutcome{SchemaVersion: 1, Result: result, Skips: append([]MigrationSkip{}, plan.Skips...)}, err
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
		if skip.Code != "already_applied" {
			return true
		}
	}
	return false
}
