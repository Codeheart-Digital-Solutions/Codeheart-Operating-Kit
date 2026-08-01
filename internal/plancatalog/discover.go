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

type Candidate struct {
	Path            string
	ExpectedKind    Kind
	FamilyQualified bool
}

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
		kind, formal := kindForFilename(entry.Name())
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
		record, parseErr := ParseDocument(candidate.Path, data, candidate.ExpectedKind)
		record.FamilyQualified = candidate.FamilyQualified
		if parseErr != nil {
			code := ErrorCode(parseErr)
			lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
			if code == "header_invalid" && len(lineIndexesOutsideFences(lines, MetadataBeginMarker)) == 0 && len(lineIndexesOutsideFences(lines, MetadataEndMarker)) == 0 {
				legacyHeader, rawStatus, legacyErr := parseLegacyHeader(data)
				if legacyErr == nil {
					record = Record{Path: candidate.Path, Header: legacyHeader, Legacy: true, ContentSHA256: sha256Text(data), ExpectedKind: candidate.ExpectedKind, FamilyQualified: candidate.FamilyQualified}
					result.Records = append(result.Records, record)
					headerSeverity := SeverityWarning
					metadataSeverity := SeverityWarning
					if mode == ModeCanonical {
						headerSeverity = SeverityError
						metadataSeverity = SeverityError
					}
					result.Problems = append(result.Problems,
						Problem{Code: "legacy_header_status", Message: fmt.Sprintf("legacy status %q is represented as lifecycle %q", rawStatus, legacyHeader.Lifecycle), Path: candidate.Path, Severity: headerSeverity, Remediation: "semantically review the lifecycle before canonical cutover"},
						Problem{Code: "metadata_missing", Message: fmt.Sprintf("%s has no bounded Codeheart plan metadata block", candidate.Path), Path: candidate.Path, Severity: metadataSeverity, Remediation: "add reviewed canonical metadata through the migration workflow"},
					)
					continue
				}
			}
			if code == "metadata_missing" {
				severity := SeverityWarning
				if mode == ModeCanonical {
					severity = SeverityError
				}
				result.Records = append(result.Records, record)
				result.Problems = append(result.Problems, Problem{Code: code, Message: parseErr.Error(), Path: candidate.Path, Severity: severity, Remediation: "add reviewed canonical metadata through the migration workflow"})
				continue
			}
			result.Problems = append(result.Problems, Problem{Code: code, Message: parseErr.Error(), Path: candidate.Path, Severity: SeverityError})
			continue
		}
		result.Records = append(result.Records, record)
	}
	result.Problems = append(result.Problems, ValidateRecords(result.Records, mode, repositoryID)...)
	SortRecords(result.Records)
	SortProblems(result.Problems)
	return result, nil
}

func kindForFilename(name string) (Kind, bool) {
	switch {
	case strings.HasSuffix(name, "_discovery_doc.md"):
		return KindDiscovery, true
	case strings.HasSuffix(name, "_implementation_doc.md"):
		return KindImplementation, true
	default:
		return "", false
	}
}

// FormalPathKind classifies an existing or prospective repository-relative planning path.
func FormalPathKind(root, path string) (Kind, bool) {
	cleaned := filepath.Clean(filepath.FromSlash(path))
	plansRoot := filepath.Join("docs", "repo", "plans")
	prefix := filepath.Join("docs", "repo", "plans") + string(filepath.Separator)
	if filepath.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) || !strings.HasPrefix(cleaned, prefix) {
		return "", false
	}
	if kind, ok := kindForFilename(filepath.Base(cleaned)); ok {
		return kind, true
	}
	if !strings.EqualFold(filepath.ToSlash(cleaned), filepath.ToSlash(filepath.Join(plansRoot, "README.md"))) && strings.EqualFold(filepath.Base(cleaned), "README.md") && qualifiesAsFamily(filepath.Join(root, filepath.Dir(cleaned))) {
		return KindFamily, true
	}
	return "", false
}

func readRegularSource(root, relative string) ([]byte, error) {
	return readRegularSourceWithHook(root, relative, nil)
}

func readRegularSourceWithHook(root, relative string, afterLstat func(string) error) ([]byte, error) {
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
	return io.ReadAll(file)
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
			_, matched = kindForFilename(nested.Name())
			return nil
		})
		if matched {
			siblings++
		}
	}
	return siblings >= 2
}
