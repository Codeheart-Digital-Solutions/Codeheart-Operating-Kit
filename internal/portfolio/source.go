package portfolio

import (
	"context"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/plancatalog"
)

func canonicalRemoteIdentity(value string) string {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		if at := strings.LastIndex(raw, "@"); at >= 0 {
			if colon := strings.Index(raw[at+1:], ":"); colon >= 0 {
				hostStart := at + 1
				hostEnd := hostStart + colon
				return canonicalHostPath(raw[hostStart:hostEnd], raw[hostEnd+1:])
			}
		}
	}
	parsed, err := url.Parse(raw)
	if err == nil && parsed.Scheme != "" {
		if parsed.Host != "" {
			host := parsed.Hostname()
			if parsed.Port() != "" {
				host += ":" + parsed.Port()
			}
			return canonicalHostPath(host, parsed.EscapedPath())
		}
		if parsed.Scheme == "file" {
			return "file:" + filepath.Clean(parsed.Path)
		}
	}
	if filepath.IsAbs(raw) {
		return "file:" + filepath.Clean(raw)
	}
	return strings.TrimSuffix(raw, ".git")
}

func canonicalHostPath(host, repositoryPath string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	repositoryPath = strings.Trim(strings.TrimSpace(repositoryPath), "/")
	repositoryPath = strings.TrimSuffix(repositoryPath, ".git")
	if host == "github.com" {
		repositoryPath = strings.ToLower(repositoryPath)
	}
	if host == "" || repositoryPath == "" {
		return ""
	}
	return host + "/" + repositoryPath
}

const (
	ConfigPath          = ".codeheart/kit.config.yaml"
	LockPath            = ".codeheart/kit.lock.yaml"
	KitMarkerPath       = ".codeheart/kit/README.md"
	KitSourceMarkerPath = "manifest.yaml"
	LocalSourcesPath    = ".codeheart/local/portfolio/sources.yaml"
	CatalogPath         = ".codeheart/local/portfolio/catalog.json"
	MirrorRootPath      = ".codeheart/local/portfolio/git"
	ScanLockPath        = ".codeheart-portfolio.scan.lock"
	OverlayPath         = "docs/repo/portfolio/strategic-overlay.yaml"
)

// Source discovers remote repositories. It does not read a developer worktree or interpret plans.
type Source interface {
	Kind() string
	Discover(context.Context) (DiscoveryResult, error)
}

// RefEnricher contributes optional provider facts for already-selected remote
// refs. Enrichment must never decide which refs the scanner evaluates.
type RefEnricher interface {
	EnrichRefs(context.Context, RepositorySource, []string) (RefEnrichment, error)
}

type RefEnrichment struct {
	PullRequests map[string]string
	APICalls     int
}

type RepositorySource struct {
	Kind          string      `json:"kind"`
	Locator       string      `json:"locator"`
	CloneURL      string      `json:"clone_url"`
	DefaultBranch string      `json:"default_branch,omitempty"`
	NameHint      string      `json:"name_hint,omitempty"`
	Self          bool        `json:"self,omitempty"`
	Enricher      RefEnricher `json:"-"`
}

type DiscoveryResult struct {
	Repositories []RepositorySource `json:"repositories"`
	APICalls     int                `json:"api_calls"`
	Complete     bool               `json:"complete"`
	Errors       []ScanError        `json:"errors"`
}

type Candidate struct {
	SourceLocator string `json:"source_locator"`
	RepositoryID  string `json:"repository_id,omitempty"`
	Reason        string `json:"reason"`
}

type ScanError struct {
	Code          string `json:"code"`
	Message       string `json:"message"`
	SourceLocator string `json:"source_locator,omitempty"`
	RepositoryID  string `json:"repository_id,omitempty"`
	Ref           string `json:"ref,omitempty"`
	Retryable     bool   `json:"retryable,omitempty"`
}

type Metrics struct {
	DurationMS       int64 `json:"duration_ms"`
	SourceCount      int   `json:"source_count"`
	MemberCount      int   `json:"member_count"`
	CandidateCount   int   `json:"candidate_count"`
	ObservationCount int   `json:"observation_count"`
	StaleCount       int   `json:"stale_count"`
	APICallCount     int   `json:"api_call_count"`
	MaxConcurrency   int   `json:"max_concurrency"`
}

type Catalog struct {
	SchemaVersion      int                             `json:"schema_version"`
	DiscoveryVersion   plancatalog.DiscoveryVersion    `json:"discovery_version,omitempty"`
	PolicyDigest       string                          `json:"policy_digest,omitempty"`
	CandidateSetDigest string                          `json:"candidate_set_digest,omitempty"`
	CoordinationHomeID string                          `json:"coordination_home_id"`
	StartedAt          string                          `json:"started_at"`
	CompletedAt        string                          `json:"completed_at"`
	LastCompleteScanAt string                          `json:"last_complete_scan_at,omitempty"`
	Complete           bool                            `json:"complete"`
	Members            []plancatalog.CatalogMember     `json:"members"`
	Observations       []plancatalog.SourceObservation `json:"observations"`
	Candidates         []Candidate                     `json:"candidates"`
	Errors             []ScanError                     `json:"errors"`
	Metrics            Metrics                         `json:"metrics"`
}

type ScanResult struct {
	Catalog           Catalog `json:"catalog"`
	CacheUpdated      bool    `json:"cache_updated"`
	PreviousPreserved bool    `json:"previous_complete_cache_preserved"`
	LocalOnlyOmitted  bool    `json:"local_only_evidence_omitted"`
}

func sortCatalog(catalog *Catalog) {
	sort.SliceStable(catalog.Members, func(i, j int) bool {
		return catalog.Members[i].RepositoryID < catalog.Members[j].RepositoryID
	})
	sort.SliceStable(catalog.Observations, func(i, j int) bool {
		left, right := catalog.Observations[i], catalog.Observations[j]
		if left.RepositoryID != right.RepositoryID {
			return left.RepositoryID < right.RepositoryID
		}
		if left.PlanID != right.PlanID {
			return left.PlanID < right.PlanID
		}
		if left.Ref != right.Ref {
			return left.Ref < right.Ref
		}
		return left.CanonicalPath < right.CanonicalPath
	})
	sort.SliceStable(catalog.Candidates, func(i, j int) bool {
		if catalog.Candidates[i].SourceLocator != catalog.Candidates[j].SourceLocator {
			return catalog.Candidates[i].SourceLocator < catalog.Candidates[j].SourceLocator
		}
		return catalog.Candidates[i].Reason < catalog.Candidates[j].Reason
	})
	sort.SliceStable(catalog.Errors, func(i, j int) bool {
		left, right := catalog.Errors[i], catalog.Errors[j]
		if left.SourceLocator != right.SourceLocator {
			return left.SourceLocator < right.SourceLocator
		}
		if left.Ref != right.Ref {
			return left.Ref < right.Ref
		}
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		return left.Message < right.Message
	})
}
