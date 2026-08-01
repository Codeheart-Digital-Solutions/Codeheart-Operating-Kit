package portfolio

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type CommandResult struct {
	Stdout []byte
	Stderr []byte
}

type CommandRunner interface {
	Run(context.Context, string, ...string) (CommandResult, error)
}

type BoundCommandRunner interface {
	RunBound(context.Context, *os.File, string, string, ...string) (CommandResult, error)
}

type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, name string, args ...string) (CommandResult, error) {
	command := exec.CommandContext(ctx, name, args...)
	if name == "git" {
		command.Env = sanitizedGitEnvironment(os.Environ())
	}
	return captureCommand(command)
}

func captureCommand(command *exec.Cmd) (CommandResult, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	return CommandResult{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}, err
}

func sanitizedGitEnvironment(environment []string) []string {
	result := make([]string, 0, len(environment)+1)
	for _, entry := range environment {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(strings.ToUpper(key), "GIT_") {
			continue
		}
		result = append(result, entry)
	}
	result = append(result, "GIT_TERMINAL_PROMPT=0")
	return result
}

func runBoundCommand(ctx context.Context, runner CommandRunner, directory *os.File, directoryPath, name string, args ...string) (CommandResult, error) {
	bound, ok := runner.(BoundCommandRunner)
	if !ok {
		return CommandResult{}, fmt.Errorf("bound_command_runner_required: scanner Git runner cannot use retained directory authority")
	}
	return bound.RunBound(ctx, directory, directoryPath, name, args...)
}

type LocalGitSource struct {
	root   string
	runner CommandRunner
}

func NewLocalGitSource(root string, runner CommandRunner) *LocalGitSource {
	if runner == nil {
		runner = ExecRunner{}
	}
	return &LocalGitSource{root: root, runner: runner}
}

func (source *LocalGitSource) Kind() string { return "local-git" }

func (source *LocalGitSource) Discover(ctx context.Context) (DiscoveryResult, error) {
	root, err := filepath.Abs(source.root)
	if err != nil {
		return DiscoveryResult{}, err
	}
	repositories := []string{}
	if isGitRepository(root) {
		repositories = append(repositories, root)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return DiscoveryResult{}, fmt.Errorf("local_source_unavailable: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(root, entry.Name())
		if isGitRepository(candidate) {
			repositories = append(repositories, candidate)
		}
	}
	sort.Strings(repositories)
	result := DiscoveryResult{Repositories: []RepositorySource{}, Complete: true, Errors: []ScanError{}}
	seen := map[string]bool{}
	for _, repository := range repositories {
		if seen[repository] {
			continue
		}
		seen[repository] = true
		remote, runErr := source.runner.Run(ctx, "git", safeGitArgs("-C", repository, "config", "--get", "remote.origin.url")...)
		if runErr != nil || strings.TrimSpace(string(remote.Stdout)) == "" {
			result.Complete = false
			result.Errors = append(result.Errors, ScanError{Code: "local_remote_unavailable", Message: "repository has no readable origin remote", SourceLocator: repository})
			continue
		}
		cloneURL := normalizeLocalRemote(repository, strings.TrimSpace(string(remote.Stdout)))
		result.Repositories = append(result.Repositories, RepositorySource{
			Kind:     source.Kind(),
			Locator:  repository,
			CloneURL: cloneURL,
			NameHint: filepath.Base(repository),
		})
	}
	return result, nil
}

func normalizeLocalRemote(repository, value string) string {
	if value == "" || filepath.IsAbs(value) || strings.Contains(value, "://") {
		return value
	}
	// SCP-like remotes such as git@example.invalid:owner/repo.git are not filesystem paths.
	if colon := strings.Index(value, ":"); colon > 0 && strings.Contains(value[:colon], "@") {
		return value
	}
	return filepath.Clean(filepath.Join(repository, value))
}

func isGitRepository(path string) bool {
	info, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil && (info.IsDir() || info.Mode().IsRegular())
}
