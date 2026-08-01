package plancatalog

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const LegacyRegisterPath = "docs/repo/plans/plan-register.md"

var legacyEntryHeading = regexp.MustCompile(`^#{2,6}\s+([A-Za-z0-9][A-Za-z0-9._:-]*)\s+-\s+(.+?)\s*$`)

type LegacyEntry struct {
	ID            string     `json:"id" yaml:"id"`
	Title         string     `json:"title" yaml:"title"`
	Type          string     `json:"type,omitempty" yaml:"type,omitempty"`
	Purpose       string     `json:"purpose,omitempty" yaml:"purpose,omitempty"`
	Lifecycle     Lifecycle  `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`
	Owner         string     `json:"owner,omitempty" yaml:"owner,omitempty"`
	CanonicalDocs []string   `json:"canonical_docs" yaml:"canonical_docs"`
	Created       string     `json:"created,omitempty" yaml:"created,omitempty"`
	LastUpdated   string     `json:"last_updated,omitempty" yaml:"last_updated,omitempty"`
	Relations     []Relation `json:"relations,omitempty" yaml:"relations,omitempty"`
}

type LegacyReconciliation struct {
	ByCanonicalPath map[string][]LegacyEntry `json:"by_canonical_path"`
	Unpaired        []LegacyEntry            `json:"unpaired"`
	Problems        []Problem                `json:"problems"`
}

func ParseLegacyRegister(data []byte) ([]LegacyEntry, []Problem) {
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
	entries := []LegacyEntry{}
	problems := []Problem{}
	seen := map[string]bool{}
	for index, item := range boundaries {
		end := len(lines)
		if index+1 < len(boundaries) {
			end = boundaries[index+1].line
		}
		entry := LegacyEntry{
			ID:            item.id,
			Title:         item.title,
			Type:          readLegacyField(lines, item.line+1, end, "Type"),
			Purpose:       readLegacyField(lines, item.line+1, end, "Purpose"),
			Lifecycle:     Lifecycle(readLegacyField(lines, item.line+1, end, "Status")),
			Owner:         readLegacyField(lines, item.line+1, end, "Owner / repository"),
			CanonicalDocs: readLegacyPaths(lines, item.line+1, end),
			Created:       readLegacyField(lines, item.line+1, end, "Created"),
			LastUpdated:   strings.TrimSuffix(readLegacyField(lines, item.line+1, end, "Last updated"), " (UTC)"),
			Relations:     readLegacyRelations(lines, item.line+1, end),
		}
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

func readLegacyField(lines []string, start, end int, name string) string {
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
			if value == "" || legacyFieldBoundary(value) {
				break
			}
			parts = append(parts, value)
		}
		return strings.Join(parts, " ")
	}
	return ""
}

func readLegacyPaths(lines []string, start, end int) []string {
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
			if value == "" || legacyFieldBoundary(value) {
				break
			}
			values = append(values, value)
		}
		for _, value := range values {
			value = strings.Trim(strings.TrimSpace(strings.TrimPrefix(value, "-")), "`")
			if normalized, ok := normalizeLegacyCanonicalPath(value); ok {
				paths = append(paths, normalized)
			}
		}
		break
	}
	sort.Strings(paths)
	return paths
}

func normalizeLegacyCanonicalPath(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" || filepath.IsAbs(filepath.FromSlash(value)) {
		return "", false
	}
	if colon := strings.Index(value, ":"); colon >= 0 && (strings.Index(value, "/") < 0 || colon < strings.Index(value, "/")) {
		return "", false
	}
	cleaned := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", false
	}
	return cleaned, true
}

func readLegacyRelations(lines []string, start, end int) []Relation {
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

func legacyFieldBoundary(line string) bool {
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
