package plancatalog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
)

type RepositorySettings struct {
	Mode             CatalogMode      `json:"mode"`
	RepositoryID     string           `json:"repository_id,omitempty"`
	CutoverRevision  string           `json:"cutover_revision,omitempty"`
	DiscoveryVersion DiscoveryVersion `json:"discovery_version"`
	ExcludedRoots    []string         `json:"excluded_roots"`
	PolicyDigest     string           `json:"policy_digest"`
}

type RepositorySnapshot struct {
	Settings           RepositorySettings   `json:"settings"`
	ConfiguredSettings RepositorySettings   `json:"-"`
	Targeted           bool                 `json:"-"`
	Complete           bool                 `json:"-"`
	Candidates         []Candidate          `json:"-"`
	PreviewCandidates  []Candidate          `json:"-"`
	PreviewRecords     []Record             `json:"-"`
	PreviewProblems    []Problem            `json:"-"`
	PolicyDigest       string               `json:"-"`
	CandidateSetDigest string               `json:"-"`
	Records            []Record             `json:"records"`
	LegacyEntries      []LegacyEntry        `json:"legacy_entries"`
	Reconciliation     LegacyReconciliation `json:"legacy_reconciliation"`
	Problems           []Problem            `json:"problems"`
}

type SnapshotOptions struct {
	TargetDiscoveryVersion DiscoveryVersion
	TargetCatalogMode      CatalogMode
	IncludeUntracked       bool
}

type ViewRow struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Kind           Kind      `json:"kind"`
	Purpose        string    `json:"purpose,omitempty"`
	Lifecycle      Lifecycle `json:"lifecycle"`
	Family         string    `json:"family,omitempty"`
	CanonicalPath  string    `json:"canonical_path"`
	Legacy         bool      `json:"legacy"`
	LegacyEvidence []string  `json:"legacy_evidence"`
}

type View struct {
	SchemaVersion              int              `json:"schema_version"`
	Mode                       CatalogMode      `json:"mode"`
	RepositoryID               string           `json:"repository_id,omitempty"`
	DiscoveryVersion           DiscoveryVersion `json:"discovery_version,omitempty"`
	ConfiguredDiscoveryVersion DiscoveryVersion `json:"configured_discovery_version,omitempty"`
	ConfiguredCatalogMode      CatalogMode      `json:"configured_catalog_mode,omitempty"`
	TargetCatalogMode          CatalogMode      `json:"target_catalog_mode,omitempty"`
	PolicyDigest               string           `json:"policy_digest,omitempty"`
	CandidateSetDigest         string           `json:"candidate_set_digest,omitempty"`
	Complete                   *bool            `json:"complete,omitempty"`
	Candidates                 []Candidate      `json:"candidates,omitempty"`
	PreviewCandidates          []Candidate      `json:"preview_candidates,omitempty"`
	PreviewProblems            []Problem        `json:"preview_problems,omitempty"`
	Rows                       []ViewRow        `json:"rows"`
	Problems                   []Problem        `json:"problems"`
}

func LoadRepositorySettings(root string) (RepositorySettings, []Problem) {
	settings := defaultRepositorySettings()
	configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return settings, nil
	}
	if err != nil {
		return settings, []Problem{{Code: "catalog_config_unreadable", Message: err.Error(), Path: state.ConfigPath, Severity: SeverityError}}
	}
	return DecodeRepositorySettings(data)
}

func defaultRepositorySettings() RepositorySettings {
	settings := RepositorySettings{Mode: ModeLegacy, DiscoveryVersion: DiscoveryV1, ExcludedRoots: []string{}}
	settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	return settings
}

func (settings RepositorySettings) DiscoveryPolicy(includeUntracked bool) DiscoveryPolicy {
	return DiscoveryPolicy{
		Version:               settings.DiscoveryVersion,
		ExcludedRoots:         append([]string{}, settings.ExcludedRoots...),
		AmbiguitySegments:     append([]string{}, conventionalAmbiguitySegments...),
		IncludeUntracked:      includeUntracked,
		AuthoritativeUniverse: "git-regular-blobs",
	}
}

// DecodeRepositorySettings decodes config bytes without consulting a checkout. Remote
// membership evaluation uses this exact decoder for default-branch policy authority.
func DecodeRepositorySettings(data []byte) (RepositorySettings, []Problem) {
	settings := defaultRepositorySettings()
	config, err := state.DecodeYAMLMap(data)
	legacyNullComponents := err == nil && config["component_settings"] == nil
	if legacyNullComponents {
		config["component_settings"] = map[string]any{}
	}
	if err == nil {
		err = state.Validate(state.ConfigV1Schema, config)
	}
	if err != nil {
		return settings, []Problem{{Code: "catalog_config_invalid", Message: err.Error(), Path: state.ConfigPath, Severity: SeverityError, Remediation: "repair the shared Operating Kit configuration before plan catalog work"}}
	}
	components := state.Map(config["component_settings"])
	planning := state.Map(components["planning-workflows"])
	mode, ok := ParseCatalogMode(state.AsString(planning["plan_catalog_mode"]))
	if !ok {
		return settings, []Problem{{Code: "catalog_mode_invalid", Message: fmt.Sprintf("plan catalog mode %q is invalid", state.AsString(planning["plan_catalog_mode"])), Path: state.ConfigPath, Severity: SeverityError}}
	}
	settings.Mode = mode
	settings.CutoverRevision = state.AsString(planning["plan_catalog_cutover_revision"])
	if rawVersion := state.AsInt(planning["plan_catalog_discovery_version"]); rawVersion != 0 {
		version, valid := ParseDiscoveryVersion(rawVersion)
		if !valid {
			return settings, []Problem{{Code: "catalog_discovery_version_invalid", Message: fmt.Sprintf("plan catalog discovery version %d is invalid", rawVersion), Path: state.ConfigPath, Severity: SeverityError}}
		}
		settings.DiscoveryVersion = version
	}
	ownership := state.Map(planning["plan_catalog_ownership"])
	settings.ExcludedRoots = stringValues(ownership["excluded_roots"])
	portfolio := state.Map(config["portfolio"])
	settings.RepositoryID = state.AsString(portfolio["member_repository_id"])
	problems := []Problem{}
	if legacyNullComponents {
		problems = append(problems, Problem{Code: "legacy_null_component_settings", Message: "legacy null component_settings was interpreted as an empty object without modifying config bytes", Path: state.ConfigPath, Severity: SeverityInfo, Remediation: "the next reviewed config migration may write an explicit empty mapping"})
	}
	if settings.Mode != ModeLegacy && settings.RepositoryID == "" {
		problems = append(problems, Problem{Code: "repository_identity_missing", Message: "mixed and canonical catalog modes require portfolio.member_repository_id", Path: state.ConfigPath, Severity: SeverityError, Remediation: "configure the stable repository identity before adopting semantic plan IDs"})
	}
	if settings.Mode == ModeMixed && settings.CutoverRevision == "" {
		problems = append(problems, Problem{Code: "mixed_cutover_revision_missing", Message: "mixed catalog mode requires plan_catalog_cutover_revision", Path: state.ConfigPath, Severity: SeverityError, Remediation: "record the exact pre-cutover Git revision that contains the frozen register and grandfathered plans"})
	}
	problems = append(problems, validateExcludedRoots(settings.ExcludedRoots)...)
	sort.Strings(settings.ExcludedRoots)
	settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	SortProblems(problems)
	return settings, problems
}

func stringValues(value any) []string {
	items, ok := value.([]any)
	if !ok {
		return []string{}
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if text, ok := item.(string); ok {
			result = append(result, text)
		}
	}
	return result
}

func validateExcludedRoots(roots []string) []Problem {
	problems := []Problem{}
	portable := map[string]string{}
	for _, root := range roots {
		code, reason := validateExcludedRoot(root)
		if code != "" {
			problems = append(problems, Problem{Code: code, Message: reason, Path: state.ConfigPath, Severity: SeverityError, Remediation: "use a unique slash-normalized repository-relative directory ending in /"})
			continue
		}
		key := cases.Fold().String(norm.NFC.String(root))
		if prior, exists := portable[key]; exists {
			problems = append(problems, Problem{Code: "excluded_root_portable_collision", Message: fmt.Sprintf("excluded roots %q and %q collide under portable path comparison", prior, root), Path: state.ConfigPath, Severity: SeverityError, Remediation: "keep one portable spelling"})
			continue
		}
		portable[key] = root
	}
	keys := make([]string, 0, len(portable))
	for key := range portable {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, left := range keys {
		for _, right := range keys[i+1:] {
			if strings.HasPrefix(right, left) {
				problems = append(problems, Problem{Code: "excluded_roots_overlap", Message: fmt.Sprintf("excluded roots %q and %q overlap", portable[left], portable[right]), Path: state.ConfigPath, Severity: SeverityError, Remediation: "retain only the shallowest intended exclusion"})
			}
		}
	}
	return problems
}

func validateExcludedRoot(root string) (string, string) {
	if root == "" || strings.TrimSpace(root) != root || !strings.HasSuffix(root, "/") {
		return "excluded_root_invalid", fmt.Sprintf("excluded root %q must be non-empty, trimmed, and end in /", root)
	}
	if strings.HasPrefix(root, "/") || strings.Contains(root, "\\") || windowsDriveAbsolute(root) {
		return "excluded_root_invalid", fmt.Sprintf("excluded root %q must be slash-normalized and repository-relative", root)
	}
	if strings.ContainsAny(root, "*?[]{}") {
		return "excluded_root_glob", fmt.Sprintf("excluded root %q must not contain glob syntax", root)
	}
	segments := strings.Split(strings.TrimSuffix(root, "/"), "/")
	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return "excluded_root_escape", fmt.Sprintf("excluded root %q contains an empty or escaping segment", root)
		}
		for _, char := range segment {
			if unicode.IsControl(char) {
				return "excluded_root_invalid", fmt.Sprintf("excluded root %q contains control characters", root)
			}
		}
	}
	if norm.NFC.String(root) != root {
		return "excluded_root_not_normalized", fmt.Sprintf("excluded root %q is not Unicode NFC-normalized", root)
	}
	return "", ""
}

func windowsDriveAbsolute(path string) bool {
	return len(path) >= 3 && ((path[0] >= 'A' && path[0] <= 'Z') || (path[0] >= 'a' && path[0] <= 'z')) && path[1] == ':' && path[2] == '/'
}

func LoadRepositorySnapshot(root string) (RepositorySnapshot, error) {
	return LoadRepositorySnapshotWithOptions(root, SnapshotOptions{})
}

func LoadRepositorySnapshotWithOptions(root string, options SnapshotOptions) (RepositorySnapshot, error) {
	configuredSettings, settingsProblems := LoadRepositorySettings(root)
	settings := configuredSettings
	// Canonical readiness is defined by the discovery-v2 contract. Treat a
	// target-mode-only request as the complete prospective v2 lens rather than
	// emitting a schema-v2 wrapper around discovery-v1 evidence.
	if options.TargetCatalogMode == ModeCanonical && options.TargetDiscoveryVersion == 0 {
		options.TargetDiscoveryVersion = DiscoveryV2
	}
	if options.TargetDiscoveryVersion != 0 {
		settings.DiscoveryVersion = options.TargetDiscoveryVersion
		settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	}
	if options.TargetCatalogMode != "" {
		settings.Mode = options.TargetCatalogMode
	}
	targeted := settings.DiscoveryVersion != configuredSettings.DiscoveryVersion || settings.Mode != configuredSettings.Mode
	if settings.Mode != ModeLegacy && settings.RepositoryID == "" && !problemCodePresent(settingsProblems, "repository_identity_missing") {
		settingsProblems = append(settingsProblems, Problem{Code: "repository_identity_missing", Message: "mixed and canonical catalog modes require portfolio.member_repository_id", Path: state.ConfigPath, Severity: SeverityError, Remediation: "configure the stable repository identity before adopting semantic plan IDs"})
	}
	if options.IncludeUntracked && settings.DiscoveryVersion != DiscoveryV2 {
		return RepositorySnapshot{}, fmt.Errorf("include_untracked_requires_discovery_v2: use --target-discovery-version 2 or activate discovery v2")
	}
	var discovery Discovery
	var candidates []Candidate
	complete := true
	policyDigest := settings.PolicyDigest
	candidateSetDigest := ""
	if settings.DiscoveryVersion == DiscoveryV2 {
		classification, err := ClassifyLocalIndex(root, settings, settings.Mode, settings.RepositoryID)
		if err != nil {
			return RepositorySnapshot{}, err
		}
		discovery = classification.Discovery
		candidates = append(candidates, classification.Candidates...)
		complete = classification.Complete
		policyDigest = classification.PolicyDigest
		candidateSetDigest = classification.CandidateSetDigest
	} else {
		var err error
		discovery, err = Discover(root, settings.Mode, settings.RepositoryID)
		if err != nil {
			return RepositorySnapshot{}, err
		}
		candidates, err = Enumerate(root)
		if err != nil {
			return RepositorySnapshot{}, err
		}
	}
	authoritativeCandidates := make([]Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		if settings.DiscoveryVersion != DiscoveryV2 || candidate.Ownership == OwnershipOwned {
			authoritativeCandidates = append(authoritativeCandidates, candidate)
		}
	}
	previewCandidates := []Candidate{}
	previewRecords := []Record{}
	previewProblems := []Problem{}
	if options.IncludeUntracked {
		preview, err := ClassifyUntrackedPreview(root, settings, settings.Mode, settings.RepositoryID)
		if err != nil {
			return RepositorySnapshot{}, err
		}
		previewCandidates = preview.Candidates
		previewRecords = preview.Discovery.Records
		previewProblems = ValidatePreviewContext(discovery.Records, candidates, preview, settings.Mode, settings.RepositoryID)
	}
	entries := []LegacyEntry{}
	legacyProblems := []Problem{}
	if data, readErr := readRegularSource(root, LegacyRegisterPath); readErr == nil {
		entries, legacyProblems = ParseLegacyRegister(data)
	} else if !os.IsNotExist(readErr) {
		code := "legacy_register_unreadable"
		remediation := "repair register readability before using legacy evidence"
		if ErrorCode(readErr) == "source_unsafe" {
			code = "legacy_register_unsafe"
			remediation = "replace the symlink or non-regular source with a contained regular register file"
		}
		legacyProblems = append(legacyProblems, Problem{Code: code, Message: readErr.Error(), Path: LegacyRegisterPath, Severity: SeverityError, Remediation: remediation})
	}
	reconciliation := ReconcileLegacy(discovery.Records, authoritativeCandidates, entries)
	problems := append([]Problem{}, settingsProblems...)
	problems = append(problems, discovery.Problems...)
	problems = append(problems, legacyProblems...)
	problems = append(problems, reconciliation.Problems...)
	for _, entry := range reconciliation.Unpaired {
		problems = append(problems, Problem{Code: "legacy_evidence_unpaired", Message: fmt.Sprintf("legacy register entry %q does not identify an enumerated local canonical document", entry.ID), Path: LegacyRegisterPath, Severity: SeverityWarning, Remediation: "retain it in inventory for semantic reconciliation"})
	}
	if settings.Mode == ModeMixed {
		baselinePaths, baselineProblems := loadMixedBaseline(root, settings.CutoverRevision)
		problems = append(problems, baselineProblems...)
		for _, record := range discovery.Records {
			if record.Metadata != nil {
				continue
			}
			if len(reconciliation.ByCanonicalPath[record.Path]) > 0 && baselinePaths[record.Path] {
				continue
			}
			problems = append(problems, Problem{Code: "mixed_new_plan_metadata_missing", Message: "mixed mode requires metadata for a formal plan without grandfathered register evidence", Path: record.Path, Severity: SeverityError, Remediation: "author the new plan with canonical metadata or complete reviewed migration evidence"})
		}
	}
	SortProblems(problems)
	return RepositorySnapshot{
		Settings:           settings,
		ConfiguredSettings: configuredSettings,
		Targeted:           targeted,
		Complete:           complete,
		Candidates:         candidates,
		PreviewCandidates:  previewCandidates,
		PreviewRecords:     previewRecords,
		PreviewProblems:    previewProblems,
		PolicyDigest:       policyDigest,
		CandidateSetDigest: candidateSetDigest,
		Records:            discovery.Records,
		LegacyEntries:      entries,
		Reconciliation:     reconciliation,
		Problems:           problems,
	}, nil
}

func problemCodePresent(problems []Problem, code string) bool {
	for _, problem := range problems {
		if problem.Code == code {
			return true
		}
	}
	return false
}

func loadMixedBaseline(root, revision string) (map[string]bool, []Problem) {
	paths := map[string]bool{}
	if revision == "" {
		return paths, nil
	}
	objectType, err := gitText(root, "cat-file", "-t", revision)
	if err != nil || objectType != "commit" {
		if err == nil {
			err = fmt.Errorf("configured cutover object is %q, not a commit", objectType)
		}
		return paths, []Problem{{Code: "mixed_cutover_revision_invalid", Message: err.Error(), Path: state.ConfigPath, Severity: SeverityError, Remediation: "record a reachable pre-cutover commit containing the frozen register and legacy plans"}}
	}
	if _, err := gitBytes(root, "merge-base", "--is-ancestor", revision, "HEAD"); err != nil {
		return paths, []Problem{{Code: "mixed_cutover_revision_not_ancestor", Message: err.Error(), Path: state.ConfigPath, Severity: SeverityError, Remediation: "record a pre-cutover commit that is an ancestor of the current revision"}}
	}
	if _, err := gitBytes(root, "cat-file", "-e", revision+":"+state.ConfigPath); err == nil {
		if !gitRevisionHasRegularFile(root, revision, state.ConfigPath) {
			return paths, []Problem{{Code: "mixed_cutover_baseline_config_invalid", Message: "cutover configuration is not a regular Git blob", Path: state.ConfigPath, Severity: SeverityError, Remediation: "select a pre-cutover commit with absent or valid legacy catalog configuration"}}
		}
		configData, readErr := gitBytes(root, "show", revision+":"+state.ConfigPath)
		config, decodeErr := state.DecodeAndValidateYAML(state.ConfigV1Schema, configData)
		if readErr != nil || decodeErr != nil {
			if readErr == nil {
				readErr = decodeErr
			}
			return paths, []Problem{{Code: "mixed_cutover_baseline_config_invalid", Message: readErr.Error(), Path: state.ConfigPath, Severity: SeverityError, Remediation: "select a pre-cutover commit with absent or valid legacy catalog configuration"}}
		}
		planning := state.Map(state.Map(config["component_settings"])["planning-workflows"])
		mode, validMode := ParseCatalogMode(state.AsString(planning["plan_catalog_mode"]))
		if !validMode || mode != ModeLegacy {
			return paths, []Problem{{Code: "mixed_cutover_revision_not_legacy", Message: "configured cutover revision already has a non-legacy plan catalog mode", Path: state.ConfigPath, Severity: SeverityError, Remediation: "record the last legacy-mode commit before mixed adoption"}}
		}
	}
	if !gitRevisionHasRegularFile(root, revision, LegacyRegisterPath) {
		return paths, []Problem{{Code: "mixed_cutover_revision_unavailable", Message: "cutover revision does not contain a regular frozen plan register", Path: LegacyRegisterPath, Severity: SeverityError, Remediation: "record a pre-cutover commit containing the regular frozen register and legacy plans"}}
	}
	registerData, err := gitBytes(root, "show", revision+":"+LegacyRegisterPath)
	if err != nil {
		return paths, []Problem{{Code: "mixed_cutover_revision_unavailable", Message: err.Error(), Path: state.ConfigPath, Severity: SeverityError, Remediation: "record a reachable pre-cutover commit containing the frozen register and legacy plans"}}
	}
	currentRegister, currentErr := readRegularSource(root, LegacyRegisterPath)
	if currentErr != nil {
		return paths, []Problem{{Code: "legacy_register_missing_after_cutover", Message: currentErr.Error(), Path: LegacyRegisterPath, Severity: SeverityError, Remediation: "restore the contained regular frozen register from the cutover revision"}}
	}
	if !bytes.Equal(currentRegister, registerData) {
		return paths, []Problem{{Code: "legacy_register_modified_after_cutover", Message: "plan-register.md differs from the frozen mixed-mode cutover revision", Path: LegacyRegisterPath, Severity: SeverityError, Remediation: "restore the frozen register and author new plans with canonical metadata"}}
	}
	entries, _ := ParseLegacyRegister(registerData)
	for _, entry := range entries {
		for _, path := range entry.CanonicalDocs {
			if gitRevisionHasRegularFile(root, revision, path) {
				paths[path] = true
			}
		}
	}
	return paths, nil
}

func gitRevisionHasRegularFile(root, revision, path string) bool {
	output, err := gitText(root, "ls-tree", revision, "--", path)
	if err != nil || output == "" {
		return false
	}
	fields := strings.Fields(output)
	return len(fields) >= 3 && (fields[0] == "100644" || fields[0] == "100755") && fields[1] == "blob"
}

func BuildView(snapshot RepositorySnapshot) View {
	rows := make([]ViewRow, 0, len(snapshot.Records))
	for _, record := range snapshot.Records {
		matches := snapshot.Reconciliation.ByCanonicalPath[record.Path]
		row := ViewRow{
			Title:          DisplayTitle(record, matches),
			Kind:           record.ExpectedKind,
			Lifecycle:      record.Header.Lifecycle,
			CanonicalPath:  record.Path,
			Legacy:         record.Metadata == nil,
			LegacyEvidence: []string{},
		}
		for _, match := range matches {
			row.LegacyEvidence = append(row.LegacyEvidence, match.ID)
		}
		if record.Metadata != nil {
			row.ID = record.Metadata.ID
			row.Kind = record.Metadata.Kind
			row.Purpose = record.Metadata.Purpose
			row.Family = record.Metadata.Family
		} else if len(matches) == 1 {
			row.ID = "legacy:" + matches[0].ID
			row.Purpose = matches[0].Purpose
		} else if len(matches) > 1 {
			row.ID = "legacy-ambiguous:" + record.Path
		} else {
			row.ID = "legacy-path:" + record.Path
		}
		rows = append(rows, row)
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].ID != rows[j].ID {
			return rows[i].ID < rows[j].ID
		}
		return rows[i].CanonicalPath < rows[j].CanonicalPath
	})
	view := View{SchemaVersion: 1, Mode: snapshot.Settings.Mode, RepositoryID: snapshot.Settings.RepositoryID, Rows: rows, Problems: append([]Problem{}, snapshot.Problems...)}
	if snapshot.Settings.DiscoveryVersion == DiscoveryV2 || snapshot.Targeted {
		complete := snapshot.Complete
		view.SchemaVersion = 2
		view.DiscoveryVersion = snapshot.Settings.DiscoveryVersion
		view.ConfiguredDiscoveryVersion = snapshot.ConfiguredSettings.DiscoveryVersion
		view.ConfiguredCatalogMode = snapshot.ConfiguredSettings.Mode
		view.TargetCatalogMode = snapshot.Settings.Mode
		view.PolicyDigest = snapshot.PolicyDigest
		view.CandidateSetDigest = snapshot.CandidateSetDigest
		view.Complete = &complete
		view.Candidates = append([]Candidate{}, snapshot.Candidates...)
		view.PreviewCandidates = append([]Candidate{}, snapshot.PreviewCandidates...)
		view.PreviewProblems = append([]Problem{}, snapshot.PreviewProblems...)
	}
	return view
}

// DisplayTitle preserves historical document headings while giving every catalog
// surface the same reviewed semantic title for legacy compatibility layouts.
func DisplayTitle(record Record, matches []LegacyEntry) string {
	title := record.Header.Title
	if record.Header.CompatibilityTitle && (title == "Document Header" || title == "Overview") && len(matches) == 1 && matches[0].Title != "" {
		return matches[0].Title
	}
	return title
}

func WriteViewJSON(writer io.Writer, view View) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(view)
}

func WriteViewText(writer io.Writer, view View) {
	fmt.Fprintf(writer, "Plans (%s mode): %d\n", view.Mode, len(view.Rows))
	if view.SchemaVersion >= 2 {
		owned, excluded, blocked, hard := candidateOwnershipCounts(view.Candidates)
		fmt.Fprintf(writer, "Discovery v%d (configured v%d); target mode=%s; complete=%t; candidates included=%d excluded=%d blocked=%d unowned=%d preview=%d.\n", view.DiscoveryVersion, view.ConfiguredDiscoveryVersion, view.TargetCatalogMode, view.Complete != nil && *view.Complete, owned, excluded, blocked, hard, len(view.PreviewCandidates))
	}
	for _, row := range view.Rows {
		family := "-"
		if row.Family != "" {
			family = row.Family
		}
		fmt.Fprintf(writer, "- %s | %s | %s | %s | family=%s | %s\n", row.Title, row.Kind, row.ID, row.Lifecycle, family, row.CanonicalPath)
	}
	if len(view.Problems) > 0 {
		fmt.Fprintf(writer, "Problems: %d\n", len(view.Problems))
		for _, problem := range view.Problems {
			location := strings.TrimSpace(problem.Path)
			if location != "" {
				location = " (" + location + ")"
			}
			fmt.Fprintf(writer, "- %s %s: %s%s\n", problem.Severity, problem.Code, problem.Message, location)
		}
	}
}

func candidateOwnershipCounts(candidates []Candidate) (owned, excluded, blocked, hard int) {
	for _, candidate := range candidates {
		switch candidate.Ownership {
		case OwnershipOwned:
			owned++
		case OwnershipExcluded:
			excluded++
		case OwnershipProspectiveBlocked:
			blocked++
		case OwnershipHardUnowned:
			hard++
		}
	}
	return
}
