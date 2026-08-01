package portfolio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

var githubOwnerPattern = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])?$`)

const (
	githubPageSize = 100
	githubAttempts = 3
)

var githubPageLimit = 100

type GitHubSource struct {
	owners []string
	runner CommandRunner
}

type githubRepository struct {
	FullName      string `json:"full_name"`
	CloneURL      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
}

type githubPullRequest struct {
	Number  int    `json:"number"`
	HTMLURL string `json:"html_url"`
	Head    struct {
		Ref  string `json:"ref"`
		Repo struct {
			FullName string `json:"full_name"`
		} `json:"repo"`
	} `json:"head"`
}

func NewGitHubSource(owners []string, runner CommandRunner) *GitHubSource {
	if runner == nil {
		runner = ExecRunner{}
	}
	copyOwners := append([]string{}, owners...)
	sort.Strings(copyOwners)
	return &GitHubSource{owners: copyOwners, runner: runner}
}

func (source *GitHubSource) Kind() string { return "github" }

func (source *GitHubSource) Discover(ctx context.Context) (DiscoveryResult, error) {
	if _, err := source.runner.Run(ctx, "gh", "auth", "status"); err != nil {
		return DiscoveryResult{}, fmt.Errorf("github_auth_unavailable: authenticated gh session required")
	}
	result := DiscoveryResult{Repositories: []RepositorySource{}, Complete: true, Errors: []ScanError{}, APICalls: 1}
	seen := map[string]bool{}
	for _, owner := range source.owners {
		if !githubOwnerPattern.MatchString(owner) {
			result.Complete = false
			result.Errors = append(result.Errors, ScanError{Code: "github_owner_invalid", Message: "GitHub owner has an invalid portable shape", SourceLocator: owner})
			continue
		}
		ownerComplete := false
		ownerFailed := false
		for page := 1; page <= githubPageLimit; page++ {
			endpoint := fmt.Sprintf("orgs/%s/repos?per_page=%d&page=%d&type=all", owner, githubPageSize, page)
			response, attempts, code, err := source.page(ctx, endpoint)
			result.APICalls += attempts
			if err != nil {
				result.Complete = false
				ownerFailed = true
				result.Errors = append(result.Errors, ScanError{Code: code, Message: err.Error(), SourceLocator: owner, Retryable: code == "github_retry_exhausted"})
				break
			}
			var repositories []githubRepository
			if err := json.Unmarshal(response, &repositories); err != nil {
				result.Complete = false
				ownerFailed = true
				result.Errors = append(result.Errors, ScanError{Code: "github_response_invalid", Message: "GitHub repository page was not a JSON array", SourceLocator: owner})
				break
			}
			for _, repository := range repositories {
				if repository.FullName == "" || repository.CloneURL == "" || repository.DefaultBranch == "" {
					result.Complete = false
					ownerFailed = true
					result.Errors = append(result.Errors, ScanError{Code: "github_repository_record_incomplete", Message: "GitHub repository evidence omitted a required locator, clone URL, or default branch", SourceLocator: owner})
					continue
				}
				if seen[repository.FullName] {
					continue
				}
				seen[repository.FullName] = true
				result.Repositories = append(result.Repositories, RepositorySource{
					Kind:          source.Kind(),
					Locator:       repository.FullName,
					CloneURL:      repository.CloneURL,
					DefaultBranch: repository.DefaultBranch,
					NameHint:      repository.FullName[strings.LastIndex(repository.FullName, "/")+1:],
					Enricher:      source,
				})
			}
			if len(repositories) < githubPageSize {
				ownerComplete = true
				break
			}
		}
		if !ownerComplete && !ownerFailed {
			result.Complete = false
			result.Errors = append(result.Errors, ScanError{Code: "github_response_truncated", Message: "GitHub pagination exceeded the configured complete-scan bound", SourceLocator: owner})
		}
	}
	sort.SliceStable(result.Repositories, func(i, j int) bool { return result.Repositories[i].Locator < result.Repositories[j].Locator })
	return result, nil
}

func (source *GitHubSource) EnrichRefs(ctx context.Context, repository RepositorySource, refs []string) (RefEnrichment, error) {
	result := RefEnrichment{PullRequests: map[string]string{}}
	if len(refs) == 0 {
		return result, nil
	}
	locator := repository.Locator
	if identity := canonicalRemoteIdentity(repository.CloneURL); strings.HasPrefix(identity, "github.com/") {
		locator = strings.TrimPrefix(identity, "github.com/")
	}
	owner, name, found := strings.Cut(locator, "/")
	if !found || !githubOwnerPattern.MatchString(owner) || strings.TrimSpace(name) == "" || strings.Contains(name, "/") {
		return result, fmt.Errorf("github_pr_repository_invalid: repository locator is not an owner/name pair")
	}
	wanted := map[string]bool{}
	for _, ref := range refs {
		wanted[ref] = true
	}
	for page := 1; page <= githubPageLimit; page++ {
		endpoint := fmt.Sprintf("repos/%s/%s/pulls?state=open&per_page=%d&page=%d", url.PathEscape(owner), url.PathEscape(name), githubPageSize, page)
		response, attempts, _, err := source.page(ctx, endpoint)
		result.APICalls += attempts
		if err != nil {
			return result, fmt.Errorf("github_pr_enrichment_unavailable: optional pull-request facts could not be read")
		}
		var pullRequests []githubPullRequest
		if err := json.Unmarshal(response, &pullRequests); err != nil {
			return result, fmt.Errorf("github_pr_response_invalid: pull-request page was not a JSON array")
		}
		for _, pullRequest := range pullRequests {
			ref := "refs/remotes/origin/" + strings.TrimPrefix(strings.TrimSpace(pullRequest.Head.Ref), "refs/heads/")
			if !wanted[ref] || !strings.EqualFold(strings.TrimSpace(pullRequest.Head.Repo.FullName), owner+"/"+name) || pullRequest.Number <= 0 || !validGitHubPullRequestURL(pullRequest.HTMLURL) {
				continue
			}
			if current := result.PullRequests[ref]; current == "" || pullRequest.HTMLURL < current {
				result.PullRequests[ref] = pullRequest.HTMLURL
			}
		}
		if len(pullRequests) < githubPageSize {
			return result, nil
		}
	}
	return result, fmt.Errorf("github_pr_response_truncated: optional pull-request pagination exceeded the configured bound")
}

func validGitHubPullRequestURL(value string) bool {
	parsed, err := url.Parse(strings.TrimSpace(value))
	return err == nil && parsed.Scheme == "https" && strings.EqualFold(parsed.Hostname(), "github.com") && parsed.RawQuery == "" && parsed.Fragment == ""
}

func (source *GitHubSource) page(ctx context.Context, endpoint string) ([]byte, int, string, error) {
	for attempt := 1; attempt <= githubAttempts; attempt++ {
		result, err := source.runner.Run(ctx, "gh", "api", "--method", "GET", endpoint)
		if err == nil {
			return result.Stdout, attempt, "", nil
		}
		if ctx.Err() != nil {
			return nil, attempt, "github_command_cancelled", fmt.Errorf("GitHub request was cancelled")
		}
		lower := strings.ToLower(string(result.Stderr))
		if strings.Contains(lower, "rate limit") || strings.Contains(lower, "secondary rate") {
			return nil, attempt, "github_rate_limit_exhausted", fmt.Errorf("GitHub rate limit prevented a complete scan")
		}
		if strings.Contains(lower, "http 403") || strings.Contains(lower, "permission denied") || strings.Contains(lower, "resource not accessible") {
			return nil, attempt, "github_permission_denied", fmt.Errorf("GitHub permissions prevented a complete scan")
		}
	}
	return nil, githubAttempts, "github_retry_exhausted", fmt.Errorf("GitHub request failed after three attempts")
}
