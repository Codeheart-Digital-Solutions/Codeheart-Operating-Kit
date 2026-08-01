package plancatalog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"go.yaml.in/yaml/v3"
)

const (
	MetadataBeginMarker = "<!-- BEGIN CODEHEART PLAN METADATA -->"
	MetadataEndMarker   = "<!-- END CODEHEART PLAN METADATA -->"
)

var (
	lastUpdatedPattern  = regexp.MustCompile(`^Last updated: (.+) \(UTC\)$`)
	createdPattern      = regexp.MustCompile(`^Created: (\d{4}-\d{2}-\d{2})$`)
	statusPattern       = regexp.MustCompile(`^Status: (draft|active|completed|superseded|archived)$`)
	legacyStatusPattern = regexp.MustCompile(`^Status:\s*(\S.*?)\s*$`)
)

type CodedError struct {
	Code string
	Err  error
}

func (err *CodedError) Error() string { return err.Err.Error() }
func (err *CodedError) Unwrap() error { return err.Err }

func ErrorCode(err error) string {
	var coded *CodedError
	if errors.As(err, &coded) {
		return coded.Code
	}
	return "plan_parse_failed"
}

func ParseDocument(path string, data []byte, expectedKind Kind) (Record, error) {
	header, lines, err := parseHeaderAndTitle(data)
	if err != nil {
		return Record{}, err
	}
	record := Record{
		Path:          path,
		Header:        header,
		ContentSHA256: sha256Text(data),
		ExpectedKind:  expectedKind,
	}
	beginIndexes := lineIndexesOutsideFences(lines, MetadataBeginMarker)
	endIndexes := lineIndexesOutsideFences(lines, MetadataEndMarker)
	if len(beginIndexes) == 0 && len(endIndexes) == 0 {
		record.Legacy = true
		return record, &CodedError{Code: "metadata_missing", Err: fmt.Errorf("%s has no bounded Codeheart plan metadata block", path)}
	}
	if len(beginIndexes) != 1 || len(endIndexes) != 1 {
		return Record{}, &CodedError{Code: "metadata_marker_count", Err: fmt.Errorf("%s must contain exactly one begin marker and one end marker", path)}
	}
	begin, end := beginIndexes[0], endIndexes[0]
	if end <= begin {
		return Record{}, &CodedError{Code: "metadata_marker_order", Err: fmt.Errorf("%s metadata end marker must follow its begin marker", path)}
	}
	if !markerImmediatelyFollowsTitle(lines, header.TitleLine-1, begin) {
		return Record{}, &CodedError{Code: "metadata_misplaced", Err: fmt.Errorf("%s metadata block must appear immediately below the canonical title", path)}
	}
	if end-begin < 3 || strings.TrimSpace(lines[begin+1]) != "```yaml" || strings.TrimSpace(lines[end-1]) != "```" {
		return Record{}, &CodedError{Code: "metadata_fence_invalid", Err: fmt.Errorf("%s metadata markers must contain one visible fenced yaml block", path)}
	}
	yamlBytes := []byte(strings.Join(lines[begin+2:end-1], "\n") + "\n")
	metadata, err := decodeMetadata(yamlBytes)
	if err != nil {
		return Record{}, err
	}
	record.Metadata = &metadata
	record.MetadataStart = begin + 1
	record.MetadataEnd = end + 1
	return record, nil
}

func ParseHeader(data []byte) (Header, error) {
	header, _, err := parseHeaderAndTitle(data)
	return header, err
}

func parseLegacyHeader(data []byte) (Header, string, error) {
	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	if len(lines) < 5 {
		return Header{}, "", &CodedError{Code: "header_missing", Err: fmt.Errorf("planning document is missing the required header and title")}
	}
	lastUpdated := lastUpdatedPattern.FindStringSubmatch(lines[0])
	created := createdPattern.FindStringSubmatch(lines[1])
	status := legacyStatusPattern.FindStringSubmatch(lines[2])
	if lastUpdated == nil || created == nil || status == nil {
		return Header{}, "", &CodedError{Code: "header_invalid", Err: fmt.Errorf("legacy planning document must retain readable Last updated, Created, and Status lines")}
	}
	if _, err := time.Parse(time.RFC3339, lastUpdated[1]); err != nil || !strings.HasSuffix(lastUpdated[1], "Z") {
		return Header{}, "", &CodedError{Code: "header_last_updated_invalid", Err: fmt.Errorf("Last updated must be a UTC RFC3339 timestamp ending in Z")}
	}
	if _, err := time.Parse("2006-01-02", created[1]); err != nil {
		return Header{}, "", &CodedError{Code: "header_created_invalid", Err: fmt.Errorf("Created must be a valid YYYY-MM-DD date")}
	}
	title, titleLine, compatibilityTitle := canonicalTitle(lines)
	if title == "" {
		return Header{}, "", &CodedError{Code: "title_missing", Err: fmt.Errorf("planning document has no canonical Markdown title")}
	}
	lifecycle := LifecycleDraft
	switch Lifecycle(status[1]) {
	case LifecycleDraft, LifecycleActive, LifecycleCompleted, LifecycleSuperseded, LifecycleArchived:
		lifecycle = Lifecycle(status[1])
	}
	legacyStatus := ""
	if string(lifecycle) != status[1] {
		legacyStatus = status[1]
	}
	return Header{LastUpdated: lastUpdated[1], Created: created[1], Lifecycle: lifecycle, LegacyStatus: legacyStatus, Title: title, TitleLine: titleLine + 1, CompatibilityTitle: compatibilityTitle}, status[1], nil
}

func decodeMetadata(data []byte) (Metadata, error) {
	var document struct {
		Plan Metadata `yaml:"plan"`
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&document); err != nil {
		return Metadata{}, &CodedError{Code: "metadata_yaml_invalid", Err: fmt.Errorf("decode plan metadata: %w", err)}
	}
	if document.Plan.SchemaVersion != 0 && document.Plan.SchemaVersion != 1 {
		return Metadata{}, &CodedError{Code: "metadata_schema_unsupported", Err: fmt.Errorf("metadata schema version %d is unsupported", document.Plan.SchemaVersion)}
	}
	for field, value := range map[string]string{
		"first_cataloged":          document.Plan.FirstCataloged,
		"catalog_metadata_updated": document.Plan.CatalogMetadataUpdated,
	} {
		if value == "" {
			continue
		}
		if _, err := time.Parse(time.RFC3339, value); err != nil || !strings.HasSuffix(value, "Z") {
			return Metadata{}, &CodedError{Code: "catalog_timestamp_invalid", Err: fmt.Errorf("%s must be a UTC RFC3339 timestamp ending in Z", field)}
		}
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != nil && !errors.Is(err, io.EOF) {
		return Metadata{}, &CodedError{Code: "metadata_yaml_invalid", Err: fmt.Errorf("decode trailing plan metadata: %w", err)}
	} else if err == nil {
		return Metadata{}, &CodedError{Code: "metadata_yaml_multiple_documents", Err: fmt.Errorf("plan metadata must contain one YAML document")}
	}
	value, err := state.DecodeYAMLMap(data)
	if err != nil {
		return Metadata{}, &CodedError{Code: "metadata_yaml_invalid", Err: err}
	}
	if err := state.Validate(state.PlanMetadataSchema, value); err != nil {
		return Metadata{}, &CodedError{Code: "metadata_schema_invalid", Err: err}
	}
	return document.Plan, nil
}

func parseHeaderAndTitle(data []byte) (Header, []string, error) {
	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	if len(lines) < 5 {
		return Header{}, nil, &CodedError{Code: "header_missing", Err: fmt.Errorf("planning document is missing the required header and title")}
	}
	lastUpdated := lastUpdatedPattern.FindStringSubmatch(lines[0])
	created := createdPattern.FindStringSubmatch(lines[1])
	status := statusPattern.FindStringSubmatch(lines[2])
	legacyStatus := ""
	if status == nil {
		candidate := legacyStatusPattern.FindStringSubmatch(lines[2])
		if candidate != nil && candidate[1] == "implementation-handoff-ready" {
			legacyStatus = candidate[1]
			status = []string{candidate[0], string(LifecycleDraft)}
		}
	}
	if lastUpdated == nil || created == nil || status == nil {
		return Header{}, nil, &CodedError{Code: "header_invalid", Err: fmt.Errorf("planning document must start with Last updated, Created, and Status header lines")}
	}
	parsedUpdate, err := time.Parse(time.RFC3339, lastUpdated[1])
	if err != nil || !strings.HasSuffix(lastUpdated[1], "Z") {
		return Header{}, nil, &CodedError{Code: "header_last_updated_invalid", Err: fmt.Errorf("Last updated must be a UTC RFC3339 timestamp ending in Z")}
	}
	if _, err := time.Parse("2006-01-02", created[1]); err != nil {
		return Header{}, nil, &CodedError{Code: "header_created_invalid", Err: fmt.Errorf("Created must be a valid YYYY-MM-DD date")}
	}
	_ = parsedUpdate
	title, titleLine, compatibilityTitle := canonicalTitle(lines)
	if title == "" {
		return Header{}, nil, &CodedError{Code: "title_missing", Err: fmt.Errorf("planning document has no canonical Markdown title")}
	}
	return Header{
		LastUpdated:        lastUpdated[1],
		Created:            created[1],
		Lifecycle:          Lifecycle(status[1]),
		LegacyStatus:       legacyStatus,
		Title:              title,
		TitleLine:          titleLine + 1,
		CompatibilityTitle: compatibilityTitle,
	}, lines, nil
}

func canonicalTitle(lines []string) (string, int, bool) {
	for index := 3; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])
		if !strings.HasPrefix(line, "# ") {
			continue
		}
		title := strings.TrimSpace(strings.TrimPrefix(line, "# "))
		if title != "Document Header" {
			return title, index, false
		}
		for nested := index + 1; nested < len(lines); nested++ {
			nestedLine := strings.TrimSpace(lines[nested])
			if strings.HasPrefix(nestedLine, "## ") {
				return strings.TrimSpace(strings.TrimPrefix(nestedLine, "## ")), nested, true
			}
			if strings.HasPrefix(nestedLine, "# ") {
				break
			}
		}
		return title, index, true
	}
	return "", -1, false
}

func lineIndexesOutsideFences(lines []string, target string) []int {
	indexes := []int{}
	fenceCharacter := byte(0)
	fenceLength := 0
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if character, length, ok := markdownFence(trimmed); ok {
			if fenceCharacter == 0 {
				fenceCharacter, fenceLength = character, length
			} else if character == fenceCharacter && length >= fenceLength {
				fenceCharacter, fenceLength = 0, 0
			}
			continue
		}
		if fenceCharacter == 0 && trimmed == target {
			indexes = append(indexes, index)
		}
	}
	return indexes
}

func markdownFence(line string) (byte, int, bool) {
	if len(line) < 3 || (line[0] != '`' && line[0] != '~') {
		return 0, 0, false
	}
	character := line[0]
	length := 0
	for length < len(line) && line[length] == character {
		length++
	}
	return character, length, length >= 3
}

func markerImmediatelyFollowsTitle(lines []string, titleIndex, markerIndex int) bool {
	if markerIndex <= titleIndex {
		return false
	}
	for index := titleIndex + 1; index < markerIndex; index++ {
		if strings.TrimSpace(lines[index]) != "" {
			return false
		}
	}
	return true
}

func sha256Text(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
