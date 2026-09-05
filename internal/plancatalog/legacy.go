package plancatalog

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const LegacyRegisterPath = "docs/repo/plans/plan-register.md"

var legacyEntryHeading = regexp.MustCompile(`^#{2,6}\s+([A-Za-z0-9][A-Za-z0-9._:-]*)\s+-\s+(.+?)\s*$`)

type LegacyStatus string

const LegacyStatusImplementationHandoffReady LegacyStatus = "implementation-handoff-ready"

type LegacyEntry struct {
	ID            string       `json:"id" yaml:"id"`
	Title         string       `json:"title" yaml:"title"`
	Type          string       `json:"type,omitempty" yaml:"type,omitempty"`
	Purpose       string       `json:"purpose,omitempty" yaml:"purpose,omitempty"`
	Lifecycle     Lifecycle    `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`
	LegacyStatus  LegacyStatus `json:"legacy_status,omitempty" yaml:"legacy_status,omitempty"`
	Owner         string       `json:"owner,omitempty" yaml:"owner,omitempty"`
	CanonicalDocs []string     `json:"canonical_docs" yaml:"canonical_docs"`
	Created       string       `json:"created,omitempty" yaml:"created,omitempty"`
	LastUpdated   string       `json:"last_updated,omitempty" yaml:"last_updated,omitempty"`
	Completed     string       `json:"completed,omitempty" yaml:"completed,omitempty"`
	Relations     []Relation   `json:"relations,omitempty" yaml:"relations,omitempty"`
	legacyV1      *legacyEntryV1
}

type legacyEntryV1 struct {
	Type          string
	Purpose       string
	Lifecycle     Lifecycle
	Owner         string
	CanonicalDocs []string
	Created       string
	LastUpdated   string
	Relations     []Relation
}

type legacyRegisterField struct {
	Values []string
}

// legacyRegisterObservation is intentionally not a wire type. It retains raw,
// possibly duplicated historical fields until projection can either produce a
// closed LegacyEntry or emit a structured blocker without choosing an
// ambiguous value.
type legacyRegisterObservation struct {
	ID     string
	Title  string
	Fields map[string][]legacyRegisterField
	Lines  []string
}

type LegacyReconciliation struct {
	ByCanonicalPath map[string][]LegacyEntry `json:"by_canonical_path"`
	Unpaired        []LegacyEntry            `json:"unpaired"`
	Problems        []Problem                `json:"problems"`
}

func ParseLegacyRegister(data []byte) ([]LegacyEntry, []Problem) {
	observations := parseLegacyRegisterObservations(data)
	entries := make([]LegacyEntry, 0, len(observations))
	problems := []Problem{}
	seen := map[string]bool{}
	for _, observation := range observations {
		entry, projectionProblems := projectLegacyRegisterObservation(observation)
		problems = append(problems, projectionProblems...)
		if seen[entry.ID] {
			problems = append(problems, Problem{Code: "legacy_duplicate_id", Message: fmt.Sprintf("legacy register ID %q occurs more than once", entry.ID), Path: LegacyRegisterPath, Severity: SeverityError, Remediation: "preserve both entries as evidence and resolve the duplicate during semantic review"})
		}
		seen[entry.ID] = true
		if len(entry.CanonicalDocs) == 0 {
			problems = append(problems, Problem{Code: "legacy_canonical_path_missing", Message: fmt.Sprintf("legacy register entry %q has no local canonical document path", entry.ID), Path: LegacyRegisterPath, Severity: SeverityWarning, Remediation: "review the entry as unpaired legacy evidence"})
		}
		entries = append(entries, entry)
	}
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })
	SortProblems(problems)
	return entries, problems
}

func parseLegacyRegisterObservations(data []byte) []legacyRegisterObservation {
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	type boundary struct {
		line  int
		id    string
		title string
	}
	boundaries := []boundary{}
	for index, line := range lines {
		match := legacyEntryHeading.FindStringSubmatch(strings.TrimSpace(line))
		if match != nil {
			boundaries = append(boundaries, boundary{line: index, id: match[1], title: strings.TrimSpace(match[2])})
		}
	}
	observations := []legacyRegisterObservation{}
	for index, item := range boundaries {
		end := len(lines)
		if index+1 < len(boundaries) {
			end = boundaries[index+1].line
		}
		observation := legacyRegisterObservation{
			ID:     item.id,
			Title:  item.title,
			Fields: map[string][]legacyRegisterField{},
			Lines:  append([]string{}, lines[item.line+1:end]...),
		}
		for _, name := range []string{"Type", "Purpose", "Status", "Owner / repository", "Canonical docs", "Created", "Last updated", "Completed", "Relations"} {
			observation.Fields[name] = readLegacyFields(lines, item.line+1, end, name)
		}
		observations = append(observations, observation)
	}
	return observations
}

func projectLegacyRegisterObservation(observation legacyRegisterObservation) (LegacyEntry, []Problem) {
	entry := LegacyEntry{ID: observation.ID, Title: observation.Title, CanonicalDocs: []string{}, legacyV1: projectLegacyRegisterV1Observation(observation)}
	problems := []Problem{}
	entry.Type, problems = projectLegacyScalar(observation, "Type", problems)
	entry.Purpose, problems = projectLegacyScalar(observation, "Purpose", problems)
	status, next := projectLegacyScalar(observation, "Status", problems)
	problems = next
	if status != "" {
		switch Lifecycle(status) {
		case LifecycleDraft, LifecycleActive, LifecycleCompleted, LifecycleSuperseded, LifecycleArchived:
			entry.Lifecycle = Lifecycle(status)
		case Lifecycle(LegacyStatusImplementationHandoffReady):
			entry.LegacyStatus = LegacyStatusImplementationHandoffReady
		default:
			problems = append(problems, legacyProjectionProblem("legacy_status_unsupported", observation.ID, "Status", "uses an unsupported historical label", "use a canonical lifecycle or a recognized pre-canonical status"))
		}
	}
	entry.Owner, problems = projectLegacyScalar(observation, "Owner / repository", problems)
	entry.CanonicalDocs, problems = projectLegacyPaths(observation, problems)
	created, next := projectLegacyScalar(observation, "Created", problems)
	problems = next
	entry.Created, problems = projectLegacyDate(observation.ID, "Created", created, false, problems)
	lastUpdated, next := projectLegacyScalar(observation, "Last updated", problems)
	problems = next
	entry.LastUpdated, problems = projectLegacyDate(observation.ID, "Last updated", lastUpdated, true, problems)
	completed, next := projectLegacyScalar(observation, "Completed", problems)
	problems = next
	entry.Completed, problems = projectLegacyDate(observation.ID, "Completed", completed, false, problems)
	entry.Relations, problems = projectLegacyRelations(observation, problems)
	return entry, problems
}

func projectLegacyRegisterV1Observation(observation legacyRegisterObservation) *legacyEntryV1 {
	lines := observation.Lines
	return &legacyEntryV1{
		Type:          readLegacyFieldV1(lines, 0, len(lines), "Type"),
		Purpose:       readLegacyFieldV1(lines, 0, len(lines), "Purpose"),
		Lifecycle:     Lifecycle(readLegacyFieldV1(lines, 0, len(lines), "Status")),
		Owner:         readLegacyFieldV1(lines, 0, len(lines), "Owner / repository"),
		CanonicalDocs: readLegacyPathsV1(lines, 0, len(lines)),
		Created:       readLegacyFieldV1(lines, 0, len(lines), "Created"),
		LastUpdated:   strings.TrimSuffix(readLegacyFieldV1(lines, 0, len(lines), "Last updated"), " (UTC)"),
		Relations:     readLegacyRelationsV1(lines, 0, len(lines)),
	}
}

func legacyEntriesForInventoryV1(entries []LegacyEntry) []LegacyEntry {
	projected := make([]LegacyEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.legacyV1 != nil {
			entry.Type = entry.legacyV1.Type
			entry.Purpose = entry.legacyV1.Purpose
			entry.Lifecycle = entry.legacyV1.Lifecycle
			entry.Owner = entry.legacyV1.Owner
			entry.CanonicalDocs = append([]string{}, entry.legacyV1.CanonicalDocs...)
			entry.Created = entry.legacyV1.Created
			entry.LastUpdated = entry.legacyV1.LastUpdated
			entry.Relations = append([]Relation{}, entry.legacyV1.Relations...)
		}
		entry.LegacyStatus = ""
		entry.Completed = ""
		projected = append(projected, entry)
	}
	return projected
}

func projectLegacyScalar(observation legacyRegisterObservation, name string, problems []Problem) (string, []Problem) {
	fields := observation.Fields[name]
	if len(fields) == 0 {
		return "", problems
	}
	if len(fields) > 1 {
		problems = append(problems, legacyProjectionProblem("legacy_field_duplicate", observation.ID, name, "occurs more than once", "retain one unambiguous field value before regenerating inventory"))
		return "", problems
	}
	value := strings.TrimSpace(strings.Join(fields[0].Values, " "))
	if value == "" {
		problems = append(problems, legacyProjectionProblem("legacy_field_malformed", observation.ID, name, "has no value", "supply one well-formed value or remove the empty field"))
		return "", problems
	}
	return value, problems
}

func projectLegacyPaths(observation legacyRegisterObservation, problems []Problem) ([]string, []Problem) {
	fields := observation.Fields["Canonical docs"]
	if len(fields) > 1 {
		problems = append(problems, legacyProjectionProblem("legacy_field_duplicate", observation.ID, "Canonical docs", "occurs more than once", "retain one unambiguous Canonical docs field before regenerating inventory"))
		return []string{}, problems
	}
	paths := []string{}
	seen := map[string]bool{}
	valid := true
	for _, field := range fields {
		if len(field.Values) == 0 {
			problems = append(problems, legacyProjectionProblem("legacy_field_malformed", observation.ID, "Canonical docs", "has no value", "supply contained repository-relative canonical paths"))
			valid = false
			continue
		}
		for _, raw := range field.Values {
			value := strings.Trim(strings.TrimSpace(strings.TrimPrefix(raw, "-")), "`")
			if legacyExternalCanonicalReference(value) {
				continue
			}
			normalized, ok := normalizeLegacyCanonicalPath(value)
			if !ok {
				problems = append(problems, legacyProjectionProblem("legacy_canonical_path_malformed", observation.ID, "Canonical docs", "contains an unsafe or malformed path", "retain only contained repository-relative canonical paths"))
				valid = false
				continue
			}
			if seen[normalized] {
				problems = append(problems, legacyProjectionProblem("legacy_canonical_path_duplicate", observation.ID, "Canonical docs", "contains a duplicate canonical path", "retain one copy of each unambiguous canonical path"))
				valid = false
				continue
			}
			seen[normalized] = true
			paths = append(paths, normalized)
		}
	}
	if !valid {
		return []string{}, problems
	}
	sort.Strings(paths)
	return paths, problems
}

func projectLegacyRelations(observation legacyRegisterObservation, problems []Problem) ([]Relation, []Problem) {
	fields := observation.Fields["Relations"]
	if len(fields) == 0 {
		return nil, problems
	}
	if len(fields) > 1 {
		problems = append(problems, legacyProjectionProblem("legacy_field_duplicate", observation.ID, "Relations", "occurs more than once", "retain one unambiguous Relations block before regenerating inventory"))
		return nil, problems
	}
	if len(fields[0].Values) == 0 {
		problems = append(problems, legacyProjectionProblem("legacy_field_malformed", observation.ID, "Relations", "has no value", "supply one well-formed relation list or remove the empty field"))
		return nil, problems
	}
	relations := []Relation{}
	seen := map[Relation]bool{}
	valid := true
	for _, raw := range fields[0].Values {
		line := strings.TrimSpace(raw)
		if !strings.HasPrefix(line, "-") {
			problems = append(problems, legacyProjectionProblem("legacy_relation_malformed", observation.ID, "Relations", "contains a non-list value", "use '- <supported-kind>: <target>' for every relation"))
			valid = false
			continue
		}
		kind, target, ok := strings.Cut(strings.TrimSpace(strings.TrimPrefix(line, "-")), ":")
		kind = strings.TrimSpace(kind)
		if !ok || kind == "" {
			problems = append(problems, legacyProjectionProblem("legacy_relation_malformed", observation.ID, "Relations", "contains a relation without a kind and target", "use '- <supported-kind>: <target>' for every relation"))
			valid = false
			continue
		}
		if !relationKinds[kind] {
			problems = append(problems, legacyProjectionProblem("legacy_relation_unsupported", observation.ID, "Relations", "contains an unsupported relation kind", "use a supported canonical relation kind"))
			valid = false
			continue
		}
		target = strings.TrimSpace(target)
		values := strings.Fields(target)
		// Historical targets may carry a ' - title' suffix. Free-form prose
		// does not identify a target: truncating it to its first word invents
		// identities and can incorrectly promote malformed evidence to a
		// duplicate blocker. Keep the schema-v1 projection separately frozen.
		if len(values) > 1 && (len(values) < 3 || values[1] != "-") {
			problems = append(problems, legacyProjectionProblem("legacy_relation_malformed", observation.ID, "Relations", "contains prose without an unambiguous target", "use one target optionally followed by ' - title'"))
			valid = false
			continue
		}
		if len(values) > 0 {
			target = values[0]
		} else {
			target = ""
		}
		if target == "" {
			problems = append(problems, legacyProjectionProblem("legacy_relation_malformed", observation.ID, "Relations", "contains a relation without a target", "supply one non-empty relation target"))
			valid = false
			continue
		}
		relation := Relation{Kind: kind, Target: target}
		if seen[relation] {
			problems = append(problems, legacyProjectionProblem("legacy_relation_duplicate", observation.ID, "Relations", "contains a duplicate relation", "retain one copy of each relation"))
			valid = false
			continue
		}
		seen[relation] = true
		relations = append(relations, relation)
	}
	if !valid {
		return nil, problems
	}
	return relations, problems
}

func projectLegacyDate(id, name, value string, timestamp bool, problems []Problem) (string, []Problem) {
	if value == "" {
		return "", problems
	}
	if timestamp {
		if !strings.HasSuffix(value, " (UTC)") {
			problems = append(problems, legacyProjectionProblem("legacy_date_malformed", id, name, "must end in (UTC)", "use one UTC RFC3339 timestamp ending in Z followed by (UTC)"))
			return "", problems
		}
		value = strings.TrimSuffix(value, " (UTC)")
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil || !strings.HasSuffix(value, "Z") || parsed.Format(time.RFC3339) != value {
			problems = append(problems, legacyProjectionProblem("legacy_date_malformed", id, name, "is not a canonical UTC RFC3339 timestamp", "use one UTC RFC3339 timestamp with second precision ending in Z followed by (UTC)"))
			return "", problems
		}
		return value, problems
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil || parsed.Format("2006-01-02") != value {
		problems = append(problems, legacyProjectionProblem("legacy_date_malformed", id, name, "is not a valid YYYY-MM-DD date", "use one valid calendar date in YYYY-MM-DD form"))
		return "", problems
	}
	return value, problems
}

func legacyProjectionProblem(code, id, field, detail, remediation string) Problem {
	return Problem{Code: code, Message: fmt.Sprintf("legacy register entry %q field %q %s", id, field, detail), Path: LegacyRegisterPath, Severity: SeverityError, Remediation: remediation}
}

func ReconcileLegacy(records []Record, candidates []Candidate, entries []LegacyEntry) LegacyReconciliation {
	result := LegacyReconciliation{ByCanonicalPath: map[string][]LegacyEntry{}, Unpaired: []LegacyEntry{}, Problems: []Problem{}}
	recordPaths := map[string]bool{}
	for _, record := range records {
		recordPaths[filepath.ToSlash(filepath.Clean(filepath.FromSlash(record.Path)))] = true
	}
	for _, candidate := range candidates {
		recordPaths[filepath.ToSlash(filepath.Clean(filepath.FromSlash(candidate.Path)))] = true
	}
	for _, entry := range entries {
		paired := false
		for _, path := range entry.CanonicalDocs {
			if !recordPaths[path] {
				continue
			}
			result.ByCanonicalPath[path] = append(result.ByCanonicalPath[path], entry)
			paired = true
		}
		if !paired {
			result.Unpaired = append(result.Unpaired, entry)
		}
	}
	for path, matches := range result.ByCanonicalPath {
		sort.SliceStable(matches, func(i, j int) bool { return matches[i].ID < matches[j].ID })
		result.ByCanonicalPath[path] = matches
		if len(matches) > 1 {
			ids := make([]string, len(matches))
			for index, match := range matches {
				ids[index] = match.ID
			}
			result.Problems = append(result.Problems, Problem{Code: "legacy_evidence_ambiguous", Message: fmt.Sprintf("canonical path %q is referenced by legacy entries %s", path, strings.Join(ids, ", ")), Path: path, Severity: SeverityError, Remediation: "select the correct legacy evidence in the reviewed migration ledger"})
		}
	}
	sort.SliceStable(result.Unpaired, func(i, j int) bool { return result.Unpaired[i].ID < result.Unpaired[j].ID })
	SortProblems(result.Problems)
	return result
}

func readLegacyFields(lines []string, start, end int, name string) []legacyRegisterField {
	prefix := name + ":"
	fields := []legacyRegisterField{}
	for index := start; index < end; index++ {
		line := strings.TrimSpace(lines[index])
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		values := []string{}
		if value := strings.TrimSpace(strings.TrimPrefix(line, prefix)); value != "" {
			values = append(values, value)
		}
		for nested := index + 1; nested < end; nested++ {
			value := strings.TrimSpace(lines[nested])
			if value == "" || legacyFieldBoundary(value) {
				break
			}
			values = append(values, value)
		}
		fields = append(fields, legacyRegisterField{Values: values})
	}
	return fields
}

func readLegacyFieldV1(lines []string, start, end int, name string) string {
	prefix := name + ":"
	for index := start; index < end; index++ {
		line := strings.TrimSpace(lines[index])
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		parts := []string{}
		if value := strings.TrimSpace(strings.TrimPrefix(line, prefix)); value != "" {
			parts = append(parts, value)
		}
		for nested := index + 1; nested < end; nested++ {
			value := strings.TrimSpace(lines[nested])
			if value == "" || legacyFieldBoundaryV1(value) {
				break
			}
			parts = append(parts, value)
		}
		return strings.Join(parts, " ")
	}
	return ""
}

func readLegacyPathsV1(lines []string, start, end int) []string {
	prefix := "Canonical docs:"
	paths := []string{}
	for index := start; index < end; index++ {
		line := strings.TrimSpace(lines[index])
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		values := []string{}
		if value := strings.TrimSpace(strings.TrimPrefix(line, prefix)); value != "" {
			values = append(values, value)
		}
		for nested := index + 1; nested < end; nested++ {
			value := strings.TrimSpace(lines[nested])
			if value == "" || legacyFieldBoundaryV1(value) {
				break
			}
			values = append(values, value)
		}
		for _, value := range values {
			value = strings.Trim(strings.TrimSpace(strings.TrimPrefix(value, "-")), "`")
			if normalized, ok := normalizeLegacyCanonicalPathV1(value); ok {
				paths = append(paths, normalized)
			}
		}
		break
	}
	sort.Strings(paths)
	return paths
}

func readLegacyRelationsV1(lines []string, start, end int) []Relation {
	relations := []Relation{}
	inRelations := false
	for index := start; index < end; index++ {
		line := strings.TrimSpace(lines[index])
		if line == "Relations:" {
			inRelations = true
			continue
		}
		if !inRelations {
			continue
		}
		if line == "" {
			if len(relations) > 0 {
				break
			}
			continue
		}
		if !strings.HasPrefix(line, "-") {
			break
		}
		kind, target, ok := strings.Cut(strings.TrimSpace(strings.TrimPrefix(line, "-")), ":")
		if !ok || !relationKinds[strings.TrimSpace(kind)] {
			continue
		}
		target = strings.TrimSpace(target)
		if fields := strings.Fields(target); len(fields) > 0 {
			target = fields[0]
		}
		relations = append(relations, Relation{Kind: strings.TrimSpace(kind), Target: target})
	}
	return relations
}

func normalizeLegacyCanonicalPath(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" || strings.Contains(value, "\\") || path.IsAbs(value) || filepath.IsAbs(filepath.FromSlash(value)) {
		return "", false
	}
	if legacyExternalCanonicalReference(value) {
		return "", false
	}
	cleaned := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", false
	}
	return cleaned, true
}

// normalizeLegacyCanonicalPathV1 preserves the v0.1.27 schema-v1 projection.
// Schema-v3 uses normalizeLegacyCanonicalPath so malformed backslash paths
// cannot escape into its closed path wire format.
func normalizeLegacyCanonicalPathV1(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" || filepath.IsAbs(filepath.FromSlash(value)) {
		return "", false
	}
	if legacyExternalCanonicalReference(value) {
		return "", false
	}
	cleaned := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", false
	}
	return cleaned, true
}

func legacyExternalCanonicalReference(value string) bool {
	colon := strings.Index(value, ":")
	slash := strings.Index(value, "/")
	return colon >= 0 && (slash < 0 || colon < slash)
}

func legacyFieldBoundary(line string) bool {
	for _, prefix := range []string{
		"Type:", "Purpose:", "Status:", "Owner / repository:", "Canonical docs:", "Created:",
		"Last updated:", "Completed:", "Priority / ordering note:", "Relations:", "Session refs:",
		"Coordination note:", "Coverage note:", "Last reviewed:", "Sync state:",
	} {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return strings.HasPrefix(line, "#")
}

func legacyFieldBoundaryV1(line string) bool {
	for _, prefix := range []string{
		"Type:", "Purpose:", "Status:", "Owner / repository:", "Canonical docs:", "Created:",
		"Last updated:", "Priority / ordering note:", "Relations:", "Session refs:",
		"Coordination note:", "Coverage note:", "Last reviewed:", "Sync state:",
	} {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return strings.HasPrefix(line, "#")
}
