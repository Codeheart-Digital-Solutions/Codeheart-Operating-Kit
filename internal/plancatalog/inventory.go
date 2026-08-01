package plancatalog

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type InventoryRecord struct {
	Path              string        `json:"path" yaml:"path"`
	Kind              Kind          `json:"kind" yaml:"kind"`
	Title             string        `json:"title" yaml:"title"`
	Lifecycle         Lifecycle     `json:"lifecycle" yaml:"lifecycle"`
	SourceRevision    string        `json:"source_revision" yaml:"source_revision"`
	SourceSHA256      string        `json:"source_sha256" yaml:"source_sha256"`
	MetadataCoverage  string        `json:"metadata_coverage" yaml:"metadata_coverage"`
	ParseStatus       string        `json:"parse_status" yaml:"parse_status"`
	ProblemCodes      []string      `json:"problem_codes" yaml:"problem_codes"`
	PlanID            string        `json:"plan_id,omitempty" yaml:"plan_id,omitempty"`
	LegacyEvidence    []LegacyEntry `json:"legacy_evidence" yaml:"legacy_evidence"`
	LegacyAmbiguous   bool          `json:"legacy_ambiguous" yaml:"legacy_ambiguous"`
	ActiveBranchTouch []string      `json:"active_branch_touches" yaml:"active_branch_touches"`
	DirtyOverlap      bool          `json:"dirty_overlap" yaml:"dirty_overlap"`
}

type InventoryCoverage struct {
	FormalRecords          int `json:"formal_records" yaml:"formal_records"`
	CanonicalMetadata      int `json:"canonical_metadata" yaml:"canonical_metadata"`
	LegacyRecords          int `json:"legacy_records" yaml:"legacy_records"`
	InvalidRecords         int `json:"invalid_records" yaml:"invalid_records"`
	UnreadableRecords      int `json:"unreadable_records" yaml:"unreadable_records"`
	UnsafeRecords          int `json:"unsafe_records" yaml:"unsafe_records"`
	UnpairedLegacyEvidence int `json:"unpaired_legacy_evidence" yaml:"unpaired_legacy_evidence"`
	ActiveBranchTouches    int `json:"active_branch_touches" yaml:"active_branch_touches"`
	DirtyOverlaps          int `json:"dirty_overlaps" yaml:"dirty_overlaps"`
}

type Inventory struct {
	SchemaVersion  int               `json:"schema_version" yaml:"schema_version"`
	RepositoryID   string            `json:"repository_id,omitempty" yaml:"repository_id,omitempty"`
	CatalogMode    CatalogMode       `json:"catalog_mode" yaml:"catalog_mode"`
	SourceRevision string            `json:"source_revision" yaml:"source_revision"`
	GeneratedAt    string            `json:"generated_at" yaml:"generated_at"`
	Records        []InventoryRecord `json:"records" yaml:"records"`
	UnpairedLegacy []LegacyEntry     `json:"unpaired_legacy_evidence" yaml:"unpaired_legacy_evidence"`
	Coverage       InventoryCoverage `json:"coverage" yaml:"coverage"`
	Problems       []Problem         `json:"problems" yaml:"problems"`
	RemoteOverlays any               `json:"remote_overlays,omitempty" yaml:"remote_overlays,omitempty"`
}

func BuildInventory(root string, now time.Time) (Inventory, error) {
	snapshot, err := LoadRepositorySnapshot(root)
	if err != nil {
		return Inventory{}, err
	}
	revision, err := gitText(root, "rev-parse", "--verify", "HEAD")
	if err != nil {
		return Inventory{}, fmt.Errorf("git_revision_unavailable: %w", err)
	}
	touches, touchProblems := activeBranchTouches(root)
	candidates, err := Enumerate(root)
	if err != nil {
		return Inventory{}, err
	}
	if now.IsZero() {
		now = time.Now()
	}
	inventory := Inventory{
		SchemaVersion:  1,
		RepositoryID:   snapshot.Settings.RepositoryID,
		CatalogMode:    snapshot.Settings.Mode,
		SourceRevision: revision,
		GeneratedAt:    now.UTC().Truncate(time.Second).Format(time.RFC3339),
		Records:        []InventoryRecord{},
		UnpairedLegacy: append([]LegacyEntry{}, snapshot.Reconciliation.Unpaired...),
		Problems:       append([]Problem{}, snapshot.Problems...),
	}
	inventory.Problems = append(inventory.Problems, touchProblems...)
	recordsByPath := map[string]Record{}
	for _, record := range snapshot.Records {
		recordsByPath[record.Path] = record
	}
	problemCodesByPath := map[string][]string{}
	for _, problem := range inventory.Problems {
		if problem.Path != "" {
			problemCodesByPath[problem.Path] = append(problemCodesByPath[problem.Path], problem.Code)
		}
	}
	for _, candidate := range candidates {
		record, parsed := recordsByPath[candidate.Path]
		coverage := "invalid"
		parseStatus := "invalid"
		planID := ""
		title := ""
		lifecycle := Lifecycle("")
		sourceSHA := ""
		if parsed && record.Metadata != nil {
			coverage = "canonical"
			parseStatus = "canonical"
			planID = record.Metadata.ID
			title = record.Header.Title
			lifecycle = record.Header.Lifecycle
			sourceSHA = record.ContentSHA256
			inventory.Coverage.CanonicalMetadata++
		} else if parsed {
			coverage = "legacy"
			parseStatus = "legacy"
			title = record.Header.Title
			lifecycle = record.Header.Lifecycle
			sourceSHA = record.ContentSHA256
			inventory.Coverage.LegacyRecords++
		} else if data, readErr := readRegularSource(root, candidate.Path); readErr != nil {
			if ErrorCode(readErr) == "source_unsafe" {
				coverage = "unsafe"
				parseStatus = "unsafe"
				inventory.Coverage.UnsafeRecords++
			} else {
				coverage = "unreadable"
				parseStatus = "unreadable"
				inventory.Coverage.UnreadableRecords++
			}
		} else {
			sourceSHA = sha256Text(data)
			inventory.Coverage.InvalidRecords++
		}
		dirty, dirtyErr := gitPathDirty(root, candidate.Path)
		if dirtyErr != nil {
			inventory.Problems = append(inventory.Problems, Problem{Code: "git_dirty_check_failed", Message: dirtyErr.Error(), Path: candidate.Path, Severity: SeverityError})
		}
		branchTouches := append([]string{}, touches[candidate.Path]...)
		if dirty {
			inventory.Coverage.DirtyOverlaps++
		}
		if len(branchTouches) > 0 {
			inventory.Coverage.ActiveBranchTouches++
		}
		matches := append([]LegacyEntry{}, snapshot.Reconciliation.ByCanonicalPath[candidate.Path]...)
		problemCodes := uniqueSortedStrings(problemCodesByPath[candidate.Path])
		inventory.Records = append(inventory.Records, InventoryRecord{
			Path:              candidate.Path,
			Kind:              candidate.ExpectedKind,
			Title:             title,
			Lifecycle:         lifecycle,
			SourceRevision:    revision,
			SourceSHA256:      sourceSHA,
			MetadataCoverage:  coverage,
			ParseStatus:       parseStatus,
			ProblemCodes:      problemCodes,
			PlanID:            planID,
			LegacyEvidence:    matches,
			LegacyAmbiguous:   len(matches) > 1,
			ActiveBranchTouch: branchTouches,
			DirtyOverlap:      dirty,
		})
	}
	inventory.Coverage.FormalRecords = len(candidates)
	inventory.Coverage.UnpairedLegacyEvidence = len(inventory.UnpairedLegacy)
	sort.SliceStable(inventory.Records, func(i, j int) bool { return inventory.Records[i].Path < inventory.Records[j].Path })
	SortProblems(inventory.Problems)
	return inventory, nil
}

func uniqueSortedStrings(values []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func activeBranchTouches(root string) (map[string][]string, []Problem) {
	result := map[string][]string{}
	seen := map[string]map[string]bool{}
	problems := []Problem{}
	current, _ := gitText(root, "symbolic-ref", "-q", "HEAD")
	refsText, err := gitText(root, "for-each-ref", "--format=%(refname)", "refs/heads", "refs/remotes")
	if err != nil {
		return result, []Problem{{Code: "branch_enumeration_failed", Message: err.Error(), Severity: SeverityError}}
	}
	refs := strings.Fields(refsText)
	sort.Strings(refs)
	for _, ref := range refs {
		if ref == current || strings.HasSuffix(ref, "/HEAD") {
			continue
		}
		merged := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", ref, "HEAD")
		if err := merged.Run(); err == nil {
			continue
		}
		paths, diffErr := gitText(root, "diff", "--name-status", "--find-renames", "--diff-filter=ACDMRT", "HEAD..."+ref, "--", "docs/repo/plans")
		if diffErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: diffErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair the ref or merge-base evidence before migration"})
			continue
		}
		for _, line := range strings.Split(strings.TrimSpace(paths), "\n") {
			fields := strings.Split(line, "\t")
			if len(fields) < 2 {
				continue
			}
			affected := fields[1:2]
			if (strings.HasPrefix(fields[0], "R") || strings.HasPrefix(fields[0], "C")) && len(fields) >= 3 {
				affected = fields[1:3]
			}
			for _, path := range affected {
				path = strings.TrimSpace(path)
				if path == "" {
					continue
				}
				if seen[path] == nil {
					seen[path] = map[string]bool{}
				}
				if seen[path][ref] {
					continue
				}
				seen[path][ref] = true
				result[path] = append(result[path], ref)
			}
		}
	}
	for path := range result {
		sort.Strings(result[path])
	}
	SortProblems(problems)
	return result, problems
}

func gitPathDirty(root, path string) (bool, error) {
	output, err := gitBytes(root, "status", "--porcelain=v1", "--untracked-files=all", "--", path)
	if err != nil {
		return false, err
	}
	return len(bytes.TrimSpace(output)) > 0, nil
}

func gitText(root string, args ...string) (string, error) {
	output, err := gitBytes(root, args...)
	return strings.TrimSpace(string(output)), err
}

func gitBytes(root string, args ...string) ([]byte, error) {
	commandArgs := append([]string{"-C", root}, args...)
	command := exec.Command("git", commandArgs...)
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return output, nil
}
