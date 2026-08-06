package plancatalog

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const maxPlanSourceBytes = 8 << 20

func Enumerate(root string) ([]Candidate, error) {
	plansRoot := filepath.Join(root, "docs", "repo", "plans")
	candidates := []Candidate{}
	err := filepath.WalkDir(plansRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		kind, formal := legacyKindForFilename(entry.Name())
		familyQualified := false
		if !formal && strings.EqualFold(entry.Name(), "README.md") && !strings.EqualFold(filepath.Clean(path), filepath.Join(plansRoot, "README.md")) {
			familyQualified = qualifiesAsFamily(filepath.Dir(path))
			if familyQualified {
				kind, formal = KindFamily, true
			}
		}
		if !formal {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		candidates = append(candidates, Candidate{Path: filepath.ToSlash(relative), ExpectedKind: kind, FamilyQualified: familyQualified})
		return nil
	})
	if os.IsNotExist(err) {
		return []Candidate{}, nil
	}
	if err != nil {
		return nil, err
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Path < candidates[j].Path })
	return candidates, nil
}

func EnumerateWithSettings(root string, settings RepositorySettings) ([]Candidate, error) {
	if settings.DiscoveryVersion != DiscoveryV2 {
		return Enumerate(root)
	}
	classification, err := ClassifyLocalIndex(root, settings, settings.Mode, settings.RepositoryID)
	if err != nil {
		return nil, err
	}
	candidates := []Candidate{}
	for _, candidate := range classification.Candidates {
		if candidate.Ownership == OwnershipOwned {
			candidates = append(candidates, candidate)
		}
	}
	return candidates, nil
}

func Discover(root string, mode CatalogMode, repositoryID string) (Discovery, error) {
	candidates, err := Enumerate(root)
	if err != nil {
		return Discovery{}, fmt.Errorf("enumerate canonical planning documents: %w", err)
	}
	result := Discovery{Mode: mode, Repository: repositoryID, Records: []Record{}, Problems: []Problem{}}
	for _, candidate := range candidates {
		data, readErr := readRegularSource(root, candidate.Path)
		if readErr != nil {
			code := "plan_unreadable"
			remediation := "repair file readability before cataloging the plan"
			if ErrorCode(readErr) == "source_unsafe" {
				code = "plan_source_unsafe"
				remediation = "replace the symlink or non-regular source with a contained regular plan file"
			}
			result.Problems = append(result.Problems, Problem{Code: code, Message: readErr.Error(), Path: candidate.Path, Severity: SeverityError, Remediation: remediation})
			continue
		}
		record, problems := parseCandidate(candidate, data, mode)
		result.Problems = append(result.Problems, problems...)
		if record.Path != "" {
			result.Records = append(result.Records, record)
		}
	}
	result.Problems = append(result.Problems, ValidateRecords(result.Records, mode, repositoryID)...)
	SortRecords(result.Records)
	SortProblems(result.Problems)
	return result, nil
}

func DiscoverWithSettings(root string, settings RepositorySettings) (Discovery, error) {
	if settings.DiscoveryVersion != DiscoveryV2 {
		return Discover(root, settings.Mode, settings.RepositoryID)
	}
	classification, err := ClassifyLocalIndex(root, settings, settings.Mode, settings.RepositoryID)
	if err != nil {
		return Discovery{}, err
	}
	return classification.Discovery, nil
}

func kindForFilename(name string) (Kind, bool) {
	switch {
	case strings.HasSuffix(name, "_discovery_doc.md") && len(strings.TrimSuffix(name, "_discovery_doc.md")) > 0:
		return KindDiscovery, true
	case strings.HasSuffix(name, "_implementation_doc.md") && len(strings.TrimSuffix(name, "_implementation_doc.md")) > 0:
		return KindImplementation, true
	default:
		return "", false
	}
}

// ProspectiveV2PathKind classifies a not-yet-tracked path using only the
// canonical discovery-v2 path and reserved-filename contract. Family authority
// still requires metadata and is therefore classified from content instead.
func ProspectiveV2PathKind(candidatePath string) (Kind, bool) {
	if validateGitPath(candidatePath) != nil || filepath.ToSlash(filepath.Clean(filepath.FromSlash(candidatePath))) != candidatePath || !hasExactDocsSegment(candidatePath) || filepath.Ext(filepath.FromSlash(candidatePath)) != ".md" {
		return "", false
	}
	return kindForFilename(filepath.Base(filepath.FromSlash(candidatePath)))
}

func legacyKindForFilename(name string) (Kind, bool) {
	switch {
	case strings.HasSuffix(name, "_discovery_doc.md"):
		return KindDiscovery, true
	case strings.HasSuffix(name, "_implementation_doc.md"):
		return KindImplementation, true
	default:
		return "", false
	}
}

func FormalPathKindWithSettings(root, candidatePath string, settings RepositorySettings) (Kind, bool) {
	if settings.DiscoveryVersion != DiscoveryV2 {
		return FormalPathKind(root, candidatePath)
	}
	blobs, err := ListIndexBlobs(root)
	if err != nil {
		return "", false
	}
	selected := []GitBlob{}
	for _, blob := range blobs {
		if blob.Path == candidatePath {
			selected = append(selected, blob)
		}
	}
	if len(selected) != 1 {
		return "", false
	}
	classification, err := ClassifyGitBlobs(selected, func(blob GitBlob) ([]byte, error) {
		return readBoundedRegularSource(root, blob.Path)
	}, settings.DiscoveryPolicy(false), settings.Mode, settings.RepositoryID, localWorktreeBoundary(root))
	if err != nil || len(classification.Candidates) != 1 || classification.Candidates[0].Ownership != OwnershipOwned {
		return "", false
	}
	candidate := classification.Candidates[0]
	return candidate.ExpectedKind, candidate.ExpectedKind != ""
}

func pathHardUnowned(value string) bool {
	_, hard := pathHardBoundary(value)
	return hard
}

// FormalPathKind classifies an existing or prospective repository-relative planning path.
func FormalPathKind(root, path string) (Kind, bool) {
	cleaned := filepath.Clean(filepath.FromSlash(path))
	plansRoot := filepath.Join("docs", "repo", "plans")
	prefix := filepath.Join("docs", "repo", "plans") + string(filepath.Separator)
	if filepath.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) || !strings.HasPrefix(cleaned, prefix) {
		return "", false
	}
	if kind, ok := legacyKindForFilename(filepath.Base(cleaned)); ok {
		return kind, true
	}
	if !strings.EqualFold(filepath.ToSlash(cleaned), filepath.ToSlash(filepath.Join(plansRoot, "README.md"))) && strings.EqualFold(filepath.Base(cleaned), "README.md") && qualifiesAsFamily(filepath.Join(root, filepath.Dir(cleaned))) {
		return KindFamily, true
	}
	return "", false
}

func readRegularSource(root, relative string) ([]byte, error) {
	return readRegularSourceWithLimit(root, relative, nil, 0)
}

func readRegularSourceWithHook(root, relative string, afterLstat func(string) error) ([]byte, error) {
	return readRegularSourceWithLimit(root, relative, afterLstat, 0)
}

func readBoundedRegularSource(root, relative string) ([]byte, error) {
	return readRegularSourceWithLimit(root, relative, nil, maxPlanSourceBytes)
}

func readRegularSourceWithLimit(root, relative string, afterLstat func(string) error, limit int64) ([]byte, error) {
	cleaned := filepath.Clean(filepath.FromSlash(relative))
	if filepath.IsAbs(cleaned) || cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return nil, &CodedError{Code: "source_unsafe", Err: fmt.Errorf("source escapes repository root: %s", relative)}
	}
	canonicalRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	canonicalRoot, err = filepath.EvalSymlinks(canonicalRoot)
	if err != nil {
		return nil, err
	}
	current := canonicalRoot
	parts := strings.Split(cleaned, string(filepath.Separator))
	var pathInfo os.FileInfo
	for index, part := range parts {
		current = filepath.Join(current, part)
		pathInfo, err = os.Lstat(current)
		if err != nil {
			return nil, err
		}
		if pathInfo.Mode()&os.ModeSymlink != 0 {
			return nil, &CodedError{Code: "source_unsafe", Err: fmt.Errorf("source traverses symbolic link: %s", relative)}
		}
		if index < len(parts)-1 && !pathInfo.IsDir() {
			return nil, &CodedError{Code: "source_unsafe", Err: fmt.Errorf("source parent is not a directory: %s", relative)}
		}
	}
	if pathInfo == nil || !pathInfo.Mode().IsRegular() {
		return nil, &CodedError{Code: "source_unsafe", Err: fmt.Errorf("source is not a regular file: %s", relative)}
	}
	if afterLstat != nil {
		if err := afterLstat(current); err != nil {
			return nil, err
		}
	}
	file, err := os.Open(current)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	openInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !openInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openInfo) {
		return nil, &CodedError{Code: "source_unsafe", Err: fmt.Errorf("source identity changed while opening: %s", relative)}
	}
	reader := io.Reader(file)
	if limit > 0 {
		reader = io.LimitReader(file, limit+1)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if limit > 0 && int64(len(data)) > limit {
		return nil, &CodedError{Code: "source_too_large", Err: fmt.Errorf("source exceeds %d-byte plan-catalog limit: %s", limit, relative)}
	}
	return data, nil
}

func qualifiesAsFamily(directory string) bool {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return false
	}
	siblings := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		matched := false
		_ = filepath.WalkDir(filepath.Join(directory, entry.Name()), func(path string, nested fs.DirEntry, walkErr error) error {
			if walkErr != nil || matched {
				return nil
			}
			if nested.IsDir() {
				return nil
			}
			_, matched = legacyKindForFilename(nested.Name())
			return nil
		})
		if matched {
			siblings++
		}
	}
	return siblings >= 2
}
