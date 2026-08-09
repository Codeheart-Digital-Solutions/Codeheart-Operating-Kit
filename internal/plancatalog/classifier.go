package plancatalog

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
)

type Candidate struct {
	Path            string              `json:"path"`
	ExpectedKind    Kind                `json:"expected_kind,omitempty"`
	FamilyQualified bool                `json:"family_qualified,omitempty"`
	Signal          CandidateSignal     `json:"signal"`
	Ownership       OwnershipClass      `json:"ownership"`
	Provenance      CandidateProvenance `json:"provenance"`
}

type ContentReader func(GitBlob) ([]byte, error)

// BoundaryClassifier contributes local worktree-only evidence, such as a nested
// repository boundary. Remote commit trees pass nil because that state is not
// represented by tree objects.
type BoundaryClassifier func(GitBlob) (boundary string, hardUnowned bool)

type Classification struct {
	Discovery          Discovery   `json:"discovery"`
	Candidates         []Candidate `json:"candidates"`
	PolicyDigest       string      `json:"policy_digest"`
	CandidateSetDigest string      `json:"candidate_set_digest"`
	Complete           bool        `json:"complete"`
}

func ClassifyGitBlobs(blobs []GitBlob, reader ContentReader, policy DiscoveryPolicy, mode CatalogMode, repositoryID string, boundary BoundaryClassifier) (Classification, error) {
	result := Classification{
		Discovery:    Discovery{Mode: mode, Repository: repositoryID, Records: []Record{}, Problems: []Problem{}},
		Candidates:   []Candidate{},
		PolicyDigest: policy.Digest(),
		Complete:     true,
	}
	ordered := append([]GitBlob{}, blobs...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Path != ordered[j].Path {
			return ordered[i].Path < ordered[j].Path
		}
		if ordered[i].Stage != ordered[j].Stage {
			return ordered[i].Stage < ordered[j].Stage
		}
		return ordered[i].ObjectID < ordered[j].ObjectID
	})
	seenPath := map[string]bool{}
	for _, blob := range ordered {
		pathProblem := validateGitPath(blob.Path)
		if pathProblem != nil {
			result.Complete = false
			result.Discovery.Problems = append(result.Discovery.Problems, *pathProblem)
			continue
		}
		markdownLike := strings.EqualFold(path.Ext(blob.Path), ".md")
		if !markdownLike {
			continue
		}
		if seenPath[blob.Path] {
			result.Complete = false
			result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: "git_index_stage_ambiguous", Message: "tracked path has multiple index stages", Path: blob.Path, Severity: SeverityError, Remediation: "resolve the Git index conflict before catalog discovery"})
			continue
		}
		seenPath[blob.Path] = true
		filenameKind, filenameSignal := kindForFilename(path.Base(blob.Path))
		foldFilenameSignal := filenameSignalFold(path.Base(blob.Path))
		markdown := path.Ext(blob.Path) == ".md"
		caseVariantMarkdown := !markdown && strings.EqualFold(path.Ext(blob.Path), ".md")
		eligible := markdown && hasExactDocsSegment(blob.Path)
		hardBoundary, hardUnowned := pathHardBoundary(blob.Path)
		if !hardUnowned {
			hardBoundary, hardUnowned = modeBoundary(blob, boundary)
		}

		if !blob.Preview && (blob.Stage != 0 || strings.Trim(blob.ObjectID, "0") == "") {
			hardUnowned = true
			hardBoundary = "non-stage-zero-index-entry"
			result.Complete = false
		}
		caseMismatchReported := false
		if caseVariantMarkdown && foldFilenameSignal {
			result.Complete = false
			result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: "plan_path_case_mismatch", Message: "plan-like Markdown extension must be exact lowercase .md", Path: blob.Path, Severity: SeverityError, Remediation: "rename the file using the exact supported lowercase extension"})
			caseMismatchReported = true
		} else if foldFilenameSignal && !filenameSignal {
			result.Complete = false
			result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: "plan_path_case_mismatch", Message: "plan-like filename must use the exact lowercase reserved suffix", Path: blob.Path, Severity: SeverityError, Remediation: "rename the file using the exact supported lowercase suffix"})
			caseMismatchReported = true
		}

		var data []byte
		metadataSignal := false
		contentUnreadable := false
		var unsafeSource error
		if markdownLike && (!hardUnowned || indexedContentReadableAcrossLocalBoundary(blob, hardBoundary)) {
			if reader == nil {
				return Classification{}, fmt.Errorf("classifier content reader is required for regular Markdown blob %s", blob.Path)
			}
			var err error
			data, err = reader(blob)
			var fallback *localSourceFallbackError
			if err != nil && errors.As(err, &fallback) {
				hardUnowned = true
				hardBoundary = fallback.Boundary
				unsafeSource = fallback.Err
				blob.ContentSHA256 = sha256Text(data)
				metadataSignal = hasGenuineMetadataMarker(data)
			} else if err != nil {
				result.Complete = false
				result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: "plan_source_unsafe", Message: err.Error(), Path: blob.Path, Severity: SeverityError, Remediation: "restore the indexed regular blob as a contained non-symlink worktree file"})
				contentUnreadable = true
			} else {
				blob.ContentSHA256 = sha256Text(data)
				metadataSignal = hasGenuineMetadataMarker(data)
			}
		}
		if !filenameSignal && !metadataSignal {
			continue
		}
		if metadataSignal && !caseMismatchReported && (!markdown || (strings.EqualFold(path.Base(blob.Path), "README.md") && path.Base(blob.Path) != "README.md")) {
			result.Complete = false
			result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: "plan_path_case_mismatch", Message: "metadata-bearing plan path must use exact lowercase .md and exact README.md family spelling", Path: blob.Path, Severity: SeverityError, Remediation: "rename the path to the exact supported spelling"})
		}

		signal := SignalBoth
		if filenameSignal && !metadataSignal {
			signal = SignalFilename
		} else if metadataSignal && !filenameSignal {
			signal = SignalMetadata
		}
		expectedKind := filenameKind
		familyQualified := false
		ownership, evidence := classifyOwnership(blob.Path, eligible, hardUnowned, hardBoundary, policy)
		candidate := Candidate{
			Path:            blob.Path,
			ExpectedKind:    expectedKind,
			FamilyQualified: familyQualified,
			Signal:          signal,
			Ownership:       ownership,
			Provenance: CandidateProvenance{
				Signal:         signal,
				Ownership:      ownership,
				PolicyDigest:   result.PolicyDigest,
				Source:         blob,
				ExclusionRoot:  evidence.exclusion,
				Boundary:       evidence.boundary,
				AmbiguousUnder: evidence.ambiguous,
			},
		}
		result.Candidates = append(result.Candidates, candidate)
		if unsafeSource != nil {
			result.Complete = false
			result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: "plan_source_unsafe", Message: unsafeSource.Error(), Path: blob.Path, Severity: SeverityError, Remediation: "restore the indexed regular blob as a contained non-symlink worktree file"})
		}

		switch ownership {
		case OwnershipExcluded:
			result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: "plan_candidate_excluded", Message: fmt.Sprintf("plan candidate is excluded by %s", evidence.exclusion), Path: blob.Path, Severity: SeverityInfo, Remediation: "retain this exclusion as reviewed inventory evidence"})
			continue
		case OwnershipProspectiveBlocked:
			result.Complete = false
			result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: "plan_documentation_root_ambiguous", Message: fmt.Sprintf("plan candidate is beneath conventional ambiguity root %s", evidence.ambiguous), Path: blob.Path, Severity: SeverityError, Remediation: "move the genuine plan to a truthful owned docs location or list the misleading root in excluded_roots"})
			continue
		case OwnershipHardUnowned:
			result.Complete = false
			if metadataSignal && path.Base(blob.Path) != "README.md" && !filenameSignal {
				result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: "canonical_filename_missing", Message: "metadata-bearing plan requires a supported plan filename", Path: blob.Path, Severity: SeverityError, Remediation: "rename the file to a supported discovery or implementation filename, or use exact README.md for a family"})
			}
			if hardBoundary == "" && hasExactDocsSegment(blob.Path) && !markdown && strings.EqualFold(path.Ext(blob.Path), ".md") {
				continue
			}
			code := "plan_path_unowned"
			message := "plan signal is outside an eligible repository-owned docs tree"
			if hardBoundary != "" {
				message = "plan signal crosses hard-unowned boundary " + hardBoundary
			} else if hasCaseVariantDocsSegment(blob.Path) {
				code = "plan_path_case_mismatch"
				message = "plan signal is beneath a case-variant docs segment; eligibility requires exact lowercase docs"
			} else if metadataSignal {
				code = "plan_metadata_misplaced"
				message = "genuine plan metadata is outside an exact lowercase docs directory segment"
			}
			result.Discovery.Problems = append(result.Discovery.Problems, Problem{Code: code, Message: message, Path: blob.Path, Severity: SeverityError, Remediation: "move the plan to a repository-owned path containing an exact lowercase docs segment"})
			continue
		}
		if contentUnreadable {
			continue
		}

		record, parseProblems := parseCandidate(candidate, data, mode)
		result.Discovery.Problems = append(result.Discovery.Problems, parseProblems...)
		if record.Path != "" {
			if record.FamilyQualified && qualifiesAsCanonicalFamily(record) {
				record.ExpectedKind = KindFamily
				record.FamilyQualified = true
				candidate.ExpectedKind = KindFamily
				candidate.FamilyQualified = true
				result.Candidates[len(result.Candidates)-1] = candidate
			}
			result.Discovery.Records = append(result.Discovery.Records, record)
		}
	}
	portableProblems := portableCandidateProblems(result.Candidates)
	result.Discovery.Problems = append(result.Discovery.Problems, portableProblems...)
	if len(portableProblems) > 0 {
		result.Complete = false
	}
	result.Discovery.Problems = append(result.Discovery.Problems, ValidateRecordsForDiscovery(result.Discovery.Records, mode, repositoryID, DiscoveryV2)...)
	SortRecords(result.Discovery.Records)
	SortProblems(result.Discovery.Problems)
	sort.SliceStable(result.Candidates, func(i, j int) bool { return result.Candidates[i].Path < result.Candidates[j].Path })
	result.CandidateSetDigest = candidateDigest(result.Candidates)
	return result, nil
}

type ownershipEvidence struct {
	exclusion string
	boundary  string
	ambiguous string
}

func classifyOwnership(candidatePath string, eligible, hard bool, hardBoundary string, policy DiscoveryPolicy) (OwnershipClass, ownershipEvidence) {
	evidence := ownershipEvidence{boundary: hardBoundary}
	if hard {
		return OwnershipHardUnowned, evidence
	}
	for _, root := range policy.ExcludedRoots {
		if strings.HasPrefix(candidatePath, root) {
			evidence.exclusion = root
			return OwnershipExcluded, evidence
		}
	}
	if !eligible {
		return OwnershipHardUnowned, evidence
	}
	segments := strings.Split(candidatePath, "/")
	ambiguous := map[string]bool{}
	for _, segment := range policy.AmbiguitySegments {
		ambiguous[segment] = true
	}
	for index, segment := range segments[:len(segments)-1] {
		if ambiguous[segment] {
			evidence.ambiguous = strings.Join(segments[:index+1], "/") + "/"
			return OwnershipProspectiveBlocked, evidence
		}
	}
	return OwnershipOwned, evidence
}

func modeBoundary(blob GitBlob, boundary BoundaryClassifier) (string, bool) {
	switch blob.Mode {
	case GitModeRegular, GitModeExecutable:
		if boundary != nil {
			return boundary(blob)
		}
		return "", false
	case GitModeSymlink:
		return "symlink", true
	case GitModeGitlink:
		return "gitlink", true
	default:
		return "non-regular-git-mode-" + string(blob.Mode), true
	}
}

func indexedContentReadableAcrossLocalBoundary(blob GitBlob, boundary string) bool {
	if blob.Mode != GitModeRegular && blob.Mode != GitModeExecutable {
		return false
	}
	return strings.HasPrefix(boundary, "worktree-") || strings.HasPrefix(boundary, "nested-repository")
}

func pathHardBoundary(value string) (string, bool) {
	for _, root := range []string{".codeheart/kit/", ".codeheart/local/", ".codeheart/user/"} {
		if strings.HasPrefix(value, root) {
			return root, true
		}
	}
	return "", false
}

func validateGitPath(value string) *Problem {
	if !utf8.ValidString(value) {
		return &Problem{Code: "plan_path_invalid", Message: "Git path is not valid UTF-8", Path: value, Severity: SeverityError, Remediation: "rename the tracked path to a valid UTF-8 portable spelling"}
	}
	if value == "" || strings.HasPrefix(value, "/") || strings.Contains(value, "\\") || strings.ContainsRune(value, 0) || path.Clean(value) != value || value == "." || value == ".." || strings.HasPrefix(value, "../") || windowsDriveAbsolute(value) {
		return &Problem{Code: "plan_path_invalid", Message: "Git path is not a normalized contained slash path", Path: value, Severity: SeverityError, Remediation: "repair the tracked path before discovery"}
	}
	for _, char := range value {
		if unicode.IsControl(char) {
			return &Problem{Code: "plan_path_invalid", Message: "Git path contains a control character", Path: value, Severity: SeverityError, Remediation: "rename the tracked path to a portable spelling"}
		}
	}
	if norm.NFC.String(value) != value {
		return &Problem{Code: "plan_path_unicode_not_normalized", Message: "Git path is not Unicode NFC-normalized", Path: value, Severity: SeverityError, Remediation: "rename the path to its NFC-normalized spelling"}
	}
	return nil
}

func portableCandidateProblems(candidates []Candidate) []Problem {
	seen := map[string]string{}
	problems := []Problem{}
	for _, candidate := range candidates {
		key := cases.Fold().String(norm.NFC.String(candidate.Path))
		if prior, exists := seen[key]; exists && prior != candidate.Path {
			paths := []string{prior, candidate.Path}
			sort.Strings(paths)
			message := fmt.Sprintf("candidate paths collide under portable comparison: %s", strings.Join(paths, ", "))
			for _, collisionPath := range paths {
				problems = append(problems, Problem{Code: "plan_path_portable_collision", Message: message, Path: collisionPath, Severity: SeverityError, Remediation: "retain one portable path spelling"})
			}
			continue
		}
		seen[key] = candidate.Path
	}
	return problems
}

// ValidatePreviewContext evaluates non-authoritative authoring candidates as
// though they were added beside the tracked catalog while keeping every new
// finding in the preview evidence channel. The authoritative discovery,
// completeness, and candidate digest are never changed by this function.
func ValidatePreviewContext(authoritativeRecords []Record, authoritativeCandidates []Candidate, preview Classification, mode CatalogMode, repositoryID string) []Problem {
	previewPaths := map[string]bool{}
	previewIDs := map[string]bool{}
	for _, candidate := range preview.Candidates {
		previewPaths[candidate.Path] = true
	}
	for _, record := range preview.Discovery.Records {
		if record.Metadata != nil {
			previewIDs[record.Metadata.ID] = true
		}
	}

	problems := []Problem{}
	for _, problem := range preview.Discovery.Problems {
		// These catalog-context findings are recomputed against the combined
		// tracked plus preview record set below.
		if problem.Code == "duplicate_plan_id" || problem.Code == "family_record_missing" {
			continue
		}
		problems = append(problems, problem)
	}
	combinedRecords := append([]Record{}, authoritativeRecords...)
	combinedRecords = append(combinedRecords, preview.Discovery.Records...)
	for _, problem := range ValidateRecordsForDiscovery(combinedRecords, mode, repositoryID, DiscoveryV2) {
		if previewPaths[problem.Path] || (problem.Code == "duplicate_plan_id" && previewIDs[problem.PlanID]) {
			problems = append(problems, problem)
		}
	}
	combinedCandidates := append([]Candidate{}, authoritativeCandidates...)
	combinedCandidates = append(combinedCandidates, preview.Candidates...)
	for _, problem := range portableCandidateProblems(combinedCandidates) {
		if previewPaths[problem.Path] {
			problems = append(problems, problem)
		}
	}
	return uniqueProblems(problems)
}

func uniqueProblems(problems []Problem) []Problem {
	seen := map[string]bool{}
	result := make([]Problem, 0, len(problems))
	for _, problem := range problems {
		key := string(problem.Severity) + "\x00" + problem.Code + "\x00" + problem.Path + "\x00" + problem.PlanID + "\x00" + problem.Message + "\x00" + problem.Remediation
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, problem)
	}
	SortProblems(result)
	return result
}

func hasExactDocsSegment(value string) bool {
	segments := strings.Split(value, "/")
	for _, segment := range segments[:len(segments)-1] {
		if segment == "docs" {
			return true
		}
	}
	return false
}

func hasCaseVariantDocsSegment(value string) bool {
	segments := strings.Split(value, "/")
	for _, segment := range segments[:len(segments)-1] {
		if strings.EqualFold(segment, "docs") && segment != "docs" {
			return true
		}
	}
	return false
}

func filenameSignalFold(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, "_discovery_doc.md") || strings.HasSuffix(lower, "_implementation_doc.md")
}

func candidateDigest(candidates []Candidate) string {
	canonical := make([]Candidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.Provenance.Source.Preview {
			continue
		}
		// A selected commit/index revision is surrounding inventory evidence,
		// not candidate identity. Clearing it keeps the same candidate set
		// digest portable across local-index and exact commit-tree adapters.
		candidate.Provenance.Source.Revision = ""
		canonical = append(canonical, candidate)
	}
	data, _ := json.Marshal(canonical)
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

// CandidateSetDigestForCandidates returns the versioned, source-revision-
// independent authority digest for an explicit candidate subset.
func CandidateSetDigestForCandidates(candidates []Candidate) string {
	return candidateDigest(candidates)
}

func parseCandidate(candidate Candidate, data []byte, mode CatalogMode) (Record, []Problem) {
	record, parseErr := ParseDocument(candidate.Path, data, candidate.ExpectedKind)
	record.FamilyQualified = candidate.FamilyQualified
	if parseErr == nil {
		problems := []Problem{}
		validFamily := qualifiesAsCanonicalFamily(record)
		if validFamily {
			record.ExpectedKind = KindFamily
			record.FamilyQualified = true
		} else if candidate.ExpectedKind == "" {
			problems = append(problems, Problem{Code: "canonical_filename_missing", Message: "metadata-bearing plan requires a supported plan filename", Path: candidate.Path, PlanID: record.Metadata.ID, Severity: SeverityError, Remediation: "rename the file to a supported discovery or implementation filename, or use exact README.md for a family"})
		}
		return record, problems
	}
	code := ErrorCode(parseErr)
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if code == "header_invalid" && len(lineIndexesOutsideFences(lines, MetadataBeginMarker)) == 0 && len(lineIndexesOutsideFences(lines, MetadataEndMarker)) == 0 {
		legacyHeader, rawStatus, legacyErr := parseLegacyHeader(data)
		if legacyErr == nil {
			record = Record{Path: candidate.Path, Header: legacyHeader, Legacy: true, ContentSHA256: sha256Text(data), ExpectedKind: candidate.ExpectedKind, FamilyQualified: candidate.FamilyQualified}
			headerSeverity := SeverityWarning
			metadataSeverity := SeverityWarning
			if mode == ModeCanonical {
				headerSeverity = SeverityError
				metadataSeverity = SeverityError
			}
			return record, []Problem{
				{Code: "legacy_header_status", Message: fmt.Sprintf("legacy status %q is represented as lifecycle %q", rawStatus, legacyHeader.Lifecycle), Path: candidate.Path, Severity: headerSeverity, Remediation: "semantically review the lifecycle before canonical cutover"},
				{Code: "metadata_missing", Message: fmt.Sprintf("%s has no bounded Codeheart plan metadata block", candidate.Path), Path: candidate.Path, Severity: metadataSeverity, Remediation: "add reviewed canonical metadata through the migration workflow"},
			}
		}
	}
	if code == "metadata_missing" {
		severity := SeverityWarning
		if mode == ModeCanonical {
			severity = SeverityError
		}
		return record, []Problem{{Code: code, Message: parseErr.Error(), Path: candidate.Path, Severity: severity, Remediation: "add reviewed canonical metadata through the migration workflow"}}
	}
	problems := []Problem{{Code: code, Message: parseErr.Error(), Path: candidate.Path, Severity: SeverityError}}
	if candidate.Signal == SignalMetadata && candidate.ExpectedKind == "" {
		problems = append(problems, Problem{Code: "canonical_filename_missing", Message: "metadata-bearing plan requires a supported plan filename", Path: candidate.Path, Severity: SeverityError, Remediation: "rename the file to a supported discovery or implementation filename, or use exact README.md for a family"})
	}
	if metadata, metadataErr := parseMetadataBlockAnywhere(candidate.Path, data); metadataErr == nil {
		partial := Record{Path: candidate.Path, Metadata: &metadata, ContentSHA256: sha256Text(data), ExpectedKind: candidate.ExpectedKind}
		return partial, problems
	}
	return Record{}, problems
}
