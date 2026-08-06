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
	Path              string               `json:"path" yaml:"path"`
	Kind              Kind                 `json:"kind" yaml:"kind"`
	Title             string               `json:"title" yaml:"title"`
	Lifecycle         Lifecycle            `json:"lifecycle" yaml:"lifecycle"`
	SourceRevision    string               `json:"source_revision" yaml:"source_revision"`
	SourceSHA256      string               `json:"source_sha256" yaml:"source_sha256"`
	MetadataCoverage  string               `json:"metadata_coverage" yaml:"metadata_coverage"`
	ParseStatus       string               `json:"parse_status" yaml:"parse_status"`
	ProblemCodes      []string             `json:"problem_codes" yaml:"problem_codes"`
	PlanID            string               `json:"plan_id,omitempty" yaml:"plan_id,omitempty"`
	LegacyEvidence    []LegacyEntry        `json:"legacy_evidence" yaml:"legacy_evidence"`
	LegacyAmbiguous   bool                 `json:"legacy_ambiguous" yaml:"legacy_ambiguous"`
	ActiveBranchTouch []string             `json:"active_branch_touches" yaml:"active_branch_touches"`
	DirtyOverlap      bool                 `json:"dirty_overlap" yaml:"dirty_overlap"`
	Signal            CandidateSignal      `json:"signal,omitempty" yaml:"signal,omitempty"`
	Ownership         OwnershipClass       `json:"ownership,omitempty" yaml:"ownership,omitempty"`
	Provenance        *CandidateProvenance `json:"provenance,omitempty" yaml:"provenance,omitempty"`
	Authoritative     *bool                `json:"authoritative,omitempty" yaml:"authoritative,omitempty"`
	Preview           bool                 `json:"preview,omitempty" yaml:"preview,omitempty"`
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
	Candidates             int `json:"candidates,omitempty" yaml:"candidates,omitempty"`
	IncludedCandidates     int `json:"included_candidates,omitempty" yaml:"included_candidates,omitempty"`
	ExcludedCandidates     int `json:"excluded_candidates,omitempty" yaml:"excluded_candidates,omitempty"`
	BlockedCandidates      int `json:"blocked_candidates,omitempty" yaml:"blocked_candidates,omitempty"`
	UnownedCandidates      int `json:"unowned_candidates,omitempty" yaml:"unowned_candidates,omitempty"`
	PreviewCandidates      int `json:"preview_candidates,omitempty" yaml:"preview_candidates,omitempty"`
}

type Inventory struct {
	SchemaVersion              int               `json:"schema_version" yaml:"schema_version"`
	RepositoryID               string            `json:"repository_id,omitempty" yaml:"repository_id,omitempty"`
	CatalogMode                CatalogMode       `json:"catalog_mode" yaml:"catalog_mode"`
	SourceRevision             string            `json:"source_revision" yaml:"source_revision"`
	GeneratedAt                string            `json:"generated_at" yaml:"generated_at"`
	Records                    []InventoryRecord `json:"records" yaml:"records"`
	UnpairedLegacy             []LegacyEntry     `json:"unpaired_legacy_evidence" yaml:"unpaired_legacy_evidence"`
	Coverage                   InventoryCoverage `json:"coverage" yaml:"coverage"`
	Problems                   []Problem         `json:"problems" yaml:"problems"`
	RemoteOverlays             any               `json:"remote_overlays,omitempty" yaml:"remote_overlays,omitempty"`
	DiscoveryVersion           DiscoveryVersion  `json:"discovery_version,omitempty" yaml:"discovery_version,omitempty"`
	ConfiguredDiscoveryVersion DiscoveryVersion  `json:"configured_discovery_version,omitempty" yaml:"configured_discovery_version,omitempty"`
	ConfiguredCatalogMode      CatalogMode       `json:"configured_catalog_mode,omitempty" yaml:"configured_catalog_mode,omitempty"`
	TargetCatalogMode          CatalogMode       `json:"target_catalog_mode,omitempty" yaml:"target_catalog_mode,omitempty"`
	PolicyDigest               string            `json:"policy_digest,omitempty" yaml:"policy_digest,omitempty"`
	CandidateSetDigest         string            `json:"candidate_set_digest,omitempty" yaml:"candidate_set_digest,omitempty"`
	Complete                   *bool             `json:"complete,omitempty" yaml:"complete,omitempty"`
	Candidates                 []Candidate       `json:"candidates,omitempty" yaml:"candidates,omitempty"`
	PreviewCandidates          []Candidate       `json:"preview_candidates,omitempty" yaml:"preview_candidates,omitempty"`
	PreviewProblems            []Problem         `json:"preview_problems,omitempty" yaml:"preview_problems,omitempty"`
	PreviewValid               *bool             `json:"preview_valid,omitempty" yaml:"preview_valid,omitempty"`
}

func BuildInventory(root string, now time.Time) (Inventory, error) {
	return BuildInventoryWithOptions(root, now, SnapshotOptions{})
}

func BuildInventoryWithOptions(root string, now time.Time, options SnapshotOptions) (Inventory, error) {
	snapshot, err := LoadRepositorySnapshotWithOptions(root, options)
	if err != nil {
		return Inventory{}, err
	}
	revision, err := gitText(root, "rev-parse", "--verify", "HEAD")
	if err != nil {
		return Inventory{}, fmt.Errorf("git_revision_unavailable: %w", err)
	}
	touches, touchProblems := activeBranchTouchesWithSettings(root, snapshot.Settings)
	candidates := append([]Candidate{}, snapshot.Candidates...)
	if snapshot.Settings.DiscoveryVersion != DiscoveryV2 {
		candidates, err = Enumerate(root)
		if err != nil {
			return Inventory{}, err
		}
	}
	allCandidates := append([]Candidate{}, candidates...)
	allCandidates = append(allCandidates, snapshot.PreviewCandidates...)
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
	if snapshot.Settings.DiscoveryVersion == DiscoveryV2 || snapshot.Targeted {
		complete := snapshot.Complete
		inventory.SchemaVersion = 2
		inventory.DiscoveryVersion = snapshot.Settings.DiscoveryVersion
		inventory.ConfiguredDiscoveryVersion = snapshot.ConfiguredSettings.DiscoveryVersion
		inventory.ConfiguredCatalogMode = snapshot.ConfiguredSettings.Mode
		inventory.TargetCatalogMode = snapshot.Settings.Mode
		inventory.PolicyDigest = snapshot.PolicyDigest
		inventory.CandidateSetDigest = snapshot.CandidateSetDigest
		inventory.Complete = &complete
		inventory.Candidates = append([]Candidate{}, candidates...)
		inventory.PreviewCandidates = append([]Candidate{}, snapshot.PreviewCandidates...)
		inventory.PreviewProblems = append([]Problem{}, snapshot.PreviewProblems...)
		if options.IncludeUntracked {
			previewValid := !HasErrors(snapshot.PreviewProblems)
			inventory.PreviewValid = &previewValid
		}
	}
	inventory.Problems = append(inventory.Problems, touchProblems...)
	if inventory.Complete != nil && HasErrors(touchProblems) {
		complete := false
		inventory.Complete = &complete
	}
	recordsByPath := map[string]Record{}
	for _, record := range snapshot.Records {
		recordsByPath[record.Path] = record
	}
	for _, record := range snapshot.PreviewRecords {
		recordsByPath[record.Path] = record
	}
	problemCodesByPath := map[string][]string{}
	for _, problem := range inventory.Problems {
		if problem.Path != "" {
			problemCodesByPath[problem.Path] = append(problemCodesByPath[problem.Path], problem.Code)
		}
	}
	for _, problem := range inventory.PreviewProblems {
		if problem.Path != "" {
			problemCodesByPath[problem.Path] = append(problemCodesByPath[problem.Path], problem.Code)
		}
	}
	for _, candidate := range allCandidates {
		record, parsed := recordsByPath[candidate.Path]
		matches := append([]LegacyEntry{}, snapshot.Reconciliation.ByCanonicalPath[candidate.Path]...)
		coverage := "invalid"
		parseStatus := "invalid"
		planID := ""
		title := ""
		lifecycle := Lifecycle("")
		sourceSHA := ""
		preview := candidate.Provenance.Source.Preview
		authoritative := !preview && (snapshot.Settings.DiscoveryVersion != DiscoveryV2 || candidate.Ownership == OwnershipOwned)
		if preview {
			coverage = "preview"
			parseStatus = "preview"
			inventory.Coverage.PreviewCandidates++
		} else if snapshot.Settings.DiscoveryVersion == DiscoveryV2 && candidate.Ownership != OwnershipOwned {
			coverage = string(candidate.Ownership)
			parseStatus = string(candidate.Ownership)
			switch candidate.Ownership {
			case OwnershipExcluded:
				inventory.Coverage.ExcludedCandidates++
			case OwnershipProspectiveBlocked:
				inventory.Coverage.BlockedCandidates++
			case OwnershipHardUnowned:
				inventory.Coverage.UnownedCandidates++
			}
		} else if parsed && record.Metadata != nil {
			coverage = "canonical"
			parseStatus = "canonical"
			planID = record.Metadata.ID
			title = DisplayTitle(record, matches)
			lifecycle = record.Header.Lifecycle
			sourceSHA = record.ContentSHA256
			if authoritative {
				inventory.Coverage.CanonicalMetadata++
			}
		} else if parsed {
			coverage = "legacy"
			parseStatus = "legacy"
			title = DisplayTitle(record, matches)
			lifecycle = record.Header.Lifecycle
			sourceSHA = record.ContentSHA256
			if authoritative {
				inventory.Coverage.LegacyRecords++
			}
		} else if snapshot.Settings.DiscoveryVersion != DiscoveryV2 {
			data, readErr := readRegularSource(root, candidate.Path)
			if readErr != nil {
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
		} else {
			sourceSHA = candidate.Provenance.Source.ContentSHA256
			if authoritative {
				inventory.Coverage.InvalidRecords++
			}
		}
		if sourceSHA == "" {
			sourceSHA = candidate.Provenance.Source.ContentSHA256
		}
		if parsed {
			title = DisplayTitle(record, matches)
			lifecycle = record.Header.Lifecycle
			if record.Metadata != nil {
				planID = record.Metadata.ID
			}
		}
		dirty, dirtyErr := gitPathDirty(root, candidate.Path)
		if dirtyErr != nil {
			inventory.Problems = append(inventory.Problems, Problem{Code: "git_dirty_check_failed", Message: dirtyErr.Error(), Path: candidate.Path, Severity: SeverityError})
		} else if authoritative && !dirty && sourceSHA != "" {
			sourceBytes, sourceErr := gitBytes(root, "show", revision+":"+candidate.Path)
			if sourceErr != nil {
				inventory.Problems = append(inventory.Problems, Problem{Code: "source_revision_unavailable", Message: sourceErr.Error(), Path: candidate.Path, Severity: SeverityError, Remediation: "commit the reviewed plan or regenerate inventory from a coherent revision"})
			} else {
				sourceSHA = sha256Text(sourceBytes)
			}
		}
		branchTouches := append([]string{}, touches[candidate.Path]...)
		if dirty && authoritative {
			inventory.Coverage.DirtyOverlaps++
		}
		if len(branchTouches) > 0 && authoritative {
			inventory.Coverage.ActiveBranchTouches++
		}
		problemCodes := uniqueSortedStrings(problemCodesByPath[candidate.Path])
		// Filename/family qualification owns the human-readable record kind.
		// Metadata kind remains identity evidence and mismatches stay visible as
		// validation problems rather than rewriting this field.
		kind := candidate.ExpectedKind
		authority := authoritative
		recordRevision := revision
		if preview {
			recordRevision = ""
		}
		var provenance *CandidateProvenance
		if snapshot.Settings.DiscoveryVersion == DiscoveryV2 {
			copy := candidate.Provenance
			provenance = &copy
		}
		inventory.Records = append(inventory.Records, InventoryRecord{
			Path:              candidate.Path,
			Kind:              kind,
			Title:             title,
			Lifecycle:         lifecycle,
			SourceRevision:    recordRevision,
			SourceSHA256:      sourceSHA,
			MetadataCoverage:  coverage,
			ParseStatus:       parseStatus,
			ProblemCodes:      problemCodes,
			PlanID:            planID,
			LegacyEvidence:    matches,
			LegacyAmbiguous:   len(matches) > 1,
			ActiveBranchTouch: branchTouches,
			DirtyOverlap:      dirty && authoritative,
			Signal:            candidate.Signal,
			Ownership:         candidate.Ownership,
			Provenance:        provenance,
			Authoritative: func() *bool {
				if snapshot.Settings.DiscoveryVersion == DiscoveryV2 {
					return &authority
				}
				return nil
			}(),
			Preview: preview,
		})
		if snapshot.Settings.DiscoveryVersion == DiscoveryV2 && !preview {
			inventory.Coverage.Candidates++
			if authoritative {
				inventory.Coverage.IncludedCandidates++
			}
		}
	}
	if snapshot.Settings.DiscoveryVersion == DiscoveryV2 {
		inventory.Coverage.FormalRecords = inventory.Coverage.IncludedCandidates
	} else {
		inventory.Coverage.FormalRecords = len(candidates)
	}
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
	settings, problems := LoadRepositorySettings(root)
	if HasErrors(problems) {
		return map[string][]string{}, problems
	}
	return activeBranchTouchesWithSettings(root, settings)
}

func activeBranchTouchesWithSettings(root string, settings RepositorySettings) (map[string][]string, []Problem) {
	if settings.DiscoveryVersion != DiscoveryV2 {
		return legacyActiveBranchTouches(root)
	}
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
		merged.Env = append(merged.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1")
		if err := merged.Run(); err == nil {
			continue
		}
		diff, diffErr := gitBytes(root, "diff", "--name-status", "-z", "--find-renames", "--diff-filter=ACDMRT", "HEAD..."+ref)
		if diffErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: diffErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair the ref or merge-base evidence before migration"})
			continue
		}
		affected, parseErr := parseNameStatusZ(diff)
		if parseErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: parseErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair the ref diff evidence before migration"})
			continue
		}
		headCandidates, classifyErr := candidatePathsAtRevision(root, "HEAD", affected, settings)
		if classifyErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: classifyErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair local Git object evidence before migration"})
			continue
		}
		refCandidates, classifyErr := candidatePathsAtRevision(root, ref, affected, settings)
		if classifyErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: classifyErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair branch Git object evidence before migration"})
			continue
		}
		for _, candidatePath := range affected {
			if !headCandidates[candidatePath] && !refCandidates[candidatePath] {
				continue
			}
			if seen[candidatePath] == nil {
				seen[candidatePath] = map[string]bool{}
			}
			if seen[candidatePath][ref] {
				continue
			}
			seen[candidatePath][ref] = true
			result[candidatePath] = append(result[candidatePath], ref)
		}
	}
	for candidatePath := range result {
		sort.Strings(result[candidatePath])
	}
	SortProblems(problems)
	return result, problems
}

func legacyActiveBranchTouches(root string) (map[string][]string, []Problem) {
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

func parseNameStatusZ(data []byte) ([]string, error) {
	parts := bytes.Split(data, []byte{0})
	paths := []string{}
	for index := 0; index < len(parts); {
		if len(parts[index]) == 0 {
			index++
			continue
		}
		status := string(parts[index])
		index++
		count := 1
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			count = 2
		}
		if index+count > len(parts) {
			return nil, fmt.Errorf("invalid NUL-delimited name-status record for %q", status)
		}
		for offset := 0; offset < count; offset++ {
			if len(parts[index+offset]) == 0 {
				return nil, fmt.Errorf("empty path in NUL-delimited name-status record for %q", status)
			}
			paths = append(paths, string(parts[index+offset]))
		}
		index += count
	}
	return uniqueSortedStrings(paths), nil
}

func candidatePathsAtRevision(root, revision string, paths []string, settings RepositorySettings) (map[string]bool, error) {
	result := map[string]bool{}
	if len(paths) == 0 {
		return result, nil
	}
	blobs := []GitBlob{}
	const pathChunkSize = 256
	for start := 0; start < len(paths); start += pathChunkSize {
		end := start + pathChunkSize
		if end > len(paths) {
			end = len(paths)
		}
		args := []string{"-C", root, "ls-tree", "-z", revision, "--"}
		args = append(args, paths[start:end]...)
		command := exec.Command("git", args...)
		command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_LITERAL_PATHSPECS=1", "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1", "LC_ALL=C")
		output, err := command.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("git ls-tree %s: %w: %s", revision, err, strings.TrimSpace(string(output)))
		}
		for _, item := range bytes.Split(output, []byte{0}) {
			if len(item) == 0 {
				continue
			}
			header, rawPath, ok := bytes.Cut(item, []byte{'\t'})
			fields := strings.Fields(string(header))
			if !ok || len(fields) != 3 {
				return nil, fmt.Errorf("invalid ls-tree record at %s", revision)
			}
			blobs = append(blobs, GitBlob{Path: string(rawPath), Mode: GitMode(fields[0]), ObjectID: fields[2], Revision: revision})
		}
	}
	batch, err := newIndexBatchReader(root)
	if err != nil {
		return nil, err
	}
	classification, classifyErr := ClassifyGitBlobs(blobs, batch.Read, settings.DiscoveryPolicy(false), settings.Mode, settings.RepositoryID, nil)
	closeErr := batch.Close()
	if classifyErr != nil {
		return nil, classifyErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	for _, candidate := range classification.Candidates {
		result[candidate.Path] = true
	}
	return result, nil
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
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1")
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return output, nil
}
