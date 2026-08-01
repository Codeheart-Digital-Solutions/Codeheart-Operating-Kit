package plancatalog

import (
	"sort"
	"strings"
)

type Kind string

const (
	KindDiscovery      Kind = "discovery"
	KindImplementation Kind = "implementation"
	KindFamily         Kind = "family"
)

type Lifecycle string

const (
	LifecycleDraft      Lifecycle = "draft"
	LifecycleActive     Lifecycle = "active"
	LifecycleCompleted  Lifecycle = "completed"
	LifecycleSuperseded Lifecycle = "superseded"
	LifecycleArchived   Lifecycle = "archived"
)

type CatalogMode string

const (
	ModeLegacy    CatalogMode = "legacy"
	ModeMixed     CatalogMode = "mixed"
	ModeCanonical CatalogMode = "canonical"
)

type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Relation struct {
	Kind   string `json:"kind" yaml:"kind"`
	Target string `json:"target" yaml:"target"`
}

type Metadata struct {
	SchemaVersion          int        `json:"schema_version" yaml:"schema_version"`
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
	LegacyAliases          []string   `json:"legacy_aliases,omitempty" yaml:"legacy_aliases,omitempty"`
}

type Header struct {
	LastUpdated        string    `json:"last_updated"`
	Created            string    `json:"created"`
	Lifecycle          Lifecycle `json:"lifecycle"`
	LegacyStatus       string    `json:"legacy_status,omitempty"`
	Title              string    `json:"title"`
	TitleLine          int       `json:"title_line"`
	CompatibilityTitle bool      `json:"compatibility_title"`
}

type Record struct {
	Path            string    `json:"path"`
	Header          Header    `json:"header"`
	Metadata        *Metadata `json:"metadata,omitempty"`
	Legacy          bool      `json:"legacy"`
	ContentSHA256   string    `json:"content_sha256"`
	MetadataStart   int       `json:"metadata_start,omitempty"`
	MetadataEnd     int       `json:"metadata_end,omitempty"`
	ExpectedKind    Kind      `json:"expected_kind"`
	FamilyQualified bool      `json:"family_qualified,omitempty"`
}

type Problem struct {
	Code        string   `json:"code"`
	Message     string   `json:"message"`
	Path        string   `json:"path,omitempty"`
	PlanID      string   `json:"plan_id,omitempty"`
	Severity    Severity `json:"severity"`
	Remediation string   `json:"remediation,omitempty"`
}

type Discovery struct {
	Mode       CatalogMode `json:"mode"`
	Repository string      `json:"repository_id,omitempty"`
	Records    []Record    `json:"records"`
	Problems   []Problem   `json:"problems"`
}

type SourceObservation struct {
	RepositoryID    string     `json:"repository_id"`
	PlanID          string     `json:"plan_id"`
	Title           string     `json:"title"`
	Kind            Kind       `json:"kind"`
	Purpose         string     `json:"purpose"`
	Lifecycle       Lifecycle  `json:"lifecycle"`
	Family          string     `json:"family,omitempty"`
	Products        []string   `json:"products,omitempty"`
	Capabilities    []string   `json:"capabilities,omitempty"`
	StrategicThemes []string   `json:"strategic_themes,omitempty"`
	Relations       []Relation `json:"relations,omitempty"`
	LegacyAliases   []string   `json:"legacy_aliases,omitempty"`
	CanonicalPath   string     `json:"canonical_path"`
	Ref             string     `json:"ref"`
	Commit          string     `json:"commit"`
	PullRequest     string     `json:"pull_request,omitempty"`
	ContentSHA256   string     `json:"content_sha256"`
	Visibility      string     `json:"visibility"`
	Verification    string     `json:"verification"`
	ObservedAt      string     `json:"observed_at"`
	SourceChangedAt string     `json:"source_changed_at,omitempty"`
	Stale           bool       `json:"stale"`
	Conflict        bool       `json:"conflict"`
}

type CatalogMember struct {
	RepositoryID   string `json:"repository_id"`
	SourceKind     string `json:"source_kind"`
	SourceLocator  string `json:"source_locator,omitempty"`
	DefaultRef     string `json:"default_ref"`
	SourceRevision string `json:"source_revision"`
	SelfMember     bool   `json:"self_member"`
}

func CoordinationHomeSelfMember(repositoryID, defaultRef, sourceRevision string) CatalogMember {
	return CatalogMember{
		RepositoryID:   repositoryID,
		SourceKind:     "local-git",
		DefaultRef:     defaultRef,
		SourceRevision: sourceRevision,
		SelfMember:     true,
	}
}

func ParseCatalogMode(value string) (CatalogMode, bool) {
	switch CatalogMode(strings.TrimSpace(value)) {
	case "":
		return ModeLegacy, true
	case ModeLegacy:
		return ModeLegacy, true
	case ModeMixed:
		return ModeMixed, true
	case ModeCanonical:
		return ModeCanonical, true
	default:
		return "", false
	}
}

func DerivedFamilyMembers(records []Record) map[string][]string {
	members := map[string][]string{}
	for _, record := range records {
		if record.Metadata == nil || record.Metadata.Family == "" {
			continue
		}
		members[record.Metadata.Family] = append(members[record.Metadata.Family], record.Metadata.ID)
	}
	for family := range members {
		sort.Strings(members[family])
	}
	return members
}

func SortRecords(records []Record) {
	sort.SliceStable(records, func(i, j int) bool {
		left, right := records[i], records[j]
		leftID, rightID := "", ""
		if left.Metadata != nil {
			leftID = left.Metadata.ID
		}
		if right.Metadata != nil {
			rightID = right.Metadata.ID
		}
		if leftID != rightID {
			return leftID < rightID
		}
		return left.Path < right.Path
	})
}

func SortProblems(problems []Problem) {
	sort.SliceStable(problems, func(i, j int) bool {
		if problems[i].Path != problems[j].Path {
			return problems[i].Path < problems[j].Path
		}
		if problems[i].Code != problems[j].Code {
			return problems[i].Code < problems[j].Code
		}
		return problems[i].Message < problems[j].Message
	})
}

func HasErrors(problems []Problem) bool {
	for _, problem := range problems {
		if problem.Severity == SeverityError {
			return true
		}
	}
	return false
}
