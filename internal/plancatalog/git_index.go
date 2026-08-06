package plancatalog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ListIndexBlobs(root string) ([]GitBlob, error) {
	command := exec.Command("git", "-C", root, "ls-files", "--stage", "-z")
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_LITERAL_PATHSPECS=1", "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1", "LC_ALL=C")
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git_index_unavailable: %w: %s", err, strings.TrimSpace(string(output)))
	}
	blobs := []GitBlob{}
	for _, record := range bytes.Split(output, []byte{0}) {
		if len(record) == 0 {
			continue
		}
		header, rawPath, ok := bytes.Cut(record, []byte{'\t'})
		if !ok {
			return nil, fmt.Errorf("git_index_invalid: stage record has no path separator")
		}
		fields := strings.Fields(string(header))
		if len(fields) != 3 {
			return nil, fmt.Errorf("git_index_invalid: stage record has %d header fields", len(fields))
		}
		stage, parseErr := strconv.Atoi(fields[2])
		if parseErr != nil {
			return nil, fmt.Errorf("git_index_invalid: parse stage %q: %w", fields[2], parseErr)
		}
		blobs = append(blobs, GitBlob{Path: string(rawPath), Mode: GitMode(fields[0]), ObjectID: fields[1], Stage: stage})
	}
	return blobs, nil
}

func ListUntrackedMarkdown(root string) ([]GitBlob, error) {
	command := exec.Command("git", "-C", root, "ls-files", "--others", "--exclude-standard", "-z")
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_LITERAL_PATHSPECS=1", "LC_ALL=C")
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git_untracked_unavailable: %w: %s", err, strings.TrimSpace(string(output)))
	}
	blobs := []GitBlob{}
	for _, rawPath := range bytes.Split(output, []byte{0}) {
		if len(rawPath) == 0 {
			continue
		}
		candidatePath := string(rawPath)
		if !strings.EqualFold(filepath.Ext(filepath.FromSlash(candidatePath)), ".md") {
			continue
		}
		mode := GitMode("untracked-nonregular")
		if validateGitPath(candidatePath) != nil {
			blobs = append(blobs, GitBlob{Path: candidatePath, Mode: mode, Preview: true})
			continue
		}
		info, statErr := os.Lstat(filepath.Join(root, filepath.FromSlash(candidatePath)))
		if statErr == nil {
			switch {
			case info.Mode()&os.ModeSymlink != 0:
				mode = GitModeSymlink
			case info.Mode().IsRegular() && info.Mode().Perm()&0o111 != 0:
				mode = GitModeExecutable
			case info.Mode().IsRegular():
				mode = GitModeRegular
			}
		}
		blobs = append(blobs, GitBlob{Path: candidatePath, Mode: mode, Preview: true})
	}
	return blobs, nil
}

func ClassifyUntrackedPreview(root string, settings RepositorySettings, mode CatalogMode, repositoryID string) (Classification, error) {
	blobs, err := ListUntrackedMarkdown(root)
	if err != nil {
		return Classification{}, err
	}
	return ClassifyGitBlobs(blobs, func(blob GitBlob) ([]byte, error) {
		return readBoundedRegularSource(root, blob.Path)
	}, settings.DiscoveryPolicy(true), mode, repositoryID, localWorktreeBoundary(root))
}

func WorktreeV2PlanSignal(root, candidatePath string) (bool, error) {
	if !strings.EqualFold(filepath.Ext(filepath.FromSlash(candidatePath)), ".md") {
		return false, nil
	}
	if problem := validateGitPath(candidatePath); problem != nil {
		return false, fmt.Errorf("%s: %s", problem.Code, problem.Message)
	}
	info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(candidatePath)))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	mode := GitMode("untracked-nonregular")
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		mode = GitModeSymlink
	case info.Mode().IsRegular() && info.Mode().Perm()&0o111 != 0:
		mode = GitModeExecutable
	case info.Mode().IsRegular():
		mode = GitModeRegular
	}
	settings := defaultRepositorySettings()
	settings.DiscoveryVersion = DiscoveryV2
	settings.PolicyDigest = settings.DiscoveryPolicy(false).Digest()
	classification, err := ClassifyGitBlobs([]GitBlob{{Path: candidatePath, Mode: mode, Preview: true}}, func(blob GitBlob) ([]byte, error) {
		return readBoundedRegularSource(root, blob.Path)
	}, settings.DiscoveryPolicy(true), ModeCanonical, "", localWorktreeBoundary(root))
	if err != nil {
		return false, err
	}
	if len(classification.Candidates) > 0 {
		return true, nil
	}
	if !classification.Complete {
		return false, fmt.Errorf("existing Markdown target could not be safely classified")
	}
	return false, nil
}

func ClassifyLocalIndex(root string, settings RepositorySettings, mode CatalogMode, repositoryID string) (Classification, error) {
	blobs, err := ListIndexBlobs(root)
	if err != nil {
		return Classification{}, err
	}
	batch, err := newIndexBatchReader(root)
	if err != nil {
		return Classification{}, err
	}
	worktreeBoundary := localWorktreeBoundary(root)
	localBoundaries := map[string]string{}
	for _, blob := range blobs {
		if validateGitPath(blob.Path) != nil {
			continue
		}
		if !strings.EqualFold(filepath.Ext(filepath.FromSlash(blob.Path)), ".md") || (blob.Mode != GitModeRegular && blob.Mode != GitModeExecutable) {
			continue
		}
		if _, hard := pathHardBoundary(blob.Path); hard {
			continue
		}
		if boundary, hard := worktreeBoundary(blob); hard {
			localBoundaries[blob.Path] = boundary
		}
	}
	reader := func(blob GitBlob) ([]byte, error) {
		boundary := localBoundaries[blob.Path]
		if boundary == "" {
			data, readErr := readBoundedRegularSource(root, blob.Path)
			if readErr == nil {
				return data, nil
			}
			boundary = "worktree-source-unsafe"
			indexData, indexErr := batch.Read(blob)
			if indexErr != nil {
				return nil, fmt.Errorf("worktree source unsafe (%v) and indexed blob unavailable: %w", readErr, indexErr)
			}
			return indexData, &localSourceFallbackError{Boundary: boundary, Err: readErr}
		}
		indexData, indexErr := batch.Read(blob)
		worktreeErr := fmt.Errorf("worktree path crosses %s", boundary)
		if indexErr != nil {
			return nil, fmt.Errorf("%v and indexed blob unavailable: %w", worktreeErr, indexErr)
		}
		return indexData, &localSourceFallbackError{Boundary: boundary, Err: worktreeErr}
	}
	boundary := func(blob GitBlob) (string, bool) {
		value := localBoundaries[blob.Path]
		return value, value != ""
	}
	classification, classifyErr := ClassifyGitBlobs(blobs, reader, settings.DiscoveryPolicy(false), mode, repositoryID, boundary)
	closeErr := batch.Close()
	if classifyErr != nil {
		return Classification{}, classifyErr
	}
	if closeErr != nil {
		return Classification{}, closeErr
	}
	return classification, nil
}

type localSourceFallbackError struct {
	Boundary string
	Err      error
}

func (err *localSourceFallbackError) Error() string { return err.Err.Error() }
func (err *localSourceFallbackError) Unwrap() error { return err.Err }

type indexBatchReader struct {
	command *exec.Cmd
	stdin   io.WriteCloser
	stdout  *bufio.Reader
	stderr  bytes.Buffer
}

func newIndexBatchReader(root string) (*indexBatchReader, error) {
	command := exec.Command("git", "-C", root, "cat-file", "--batch")
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_LITERAL_PATHSPECS=1", "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1", "LC_ALL=C")
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	reader := &indexBatchReader{command: command, stdin: stdin, stdout: bufio.NewReader(stdoutPipe)}
	command.Stderr = &reader.stderr
	if err := command.Start(); err != nil {
		return nil, err
	}
	return reader, nil
}

func (reader *indexBatchReader) Read(blob GitBlob) ([]byte, error) {
	if _, err := io.WriteString(reader.stdin, blob.ObjectID+"\n"); err != nil {
		return nil, err
	}
	header, err := reader.stdout.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(header)
	if len(fields) != 3 || fields[1] != "blob" {
		return nil, fmt.Errorf("git_index_blob_unavailable: %s", strings.TrimSpace(header))
	}
	size, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil || size < 0 {
		return nil, fmt.Errorf("git_index_blob_invalid_size: %s", strings.TrimSpace(header))
	}
	if size > maxPlanSourceBytes {
		if _, discardErr := io.CopyN(io.Discard, reader.stdout, size+1); discardErr != nil {
			return nil, discardErr
		}
		return nil, &CodedError{Code: "source_too_large", Err: fmt.Errorf("indexed source exceeds %d-byte plan-catalog limit: %s", maxPlanSourceBytes, blob.Path)}
	}
	data := make([]byte, size)
	if _, err := io.ReadFull(reader.stdout, data); err != nil {
		return nil, err
	}
	if separator, err := reader.stdout.ReadByte(); err != nil || separator != '\n' {
		return nil, fmt.Errorf("git_index_blob_invalid_terminator: %v", err)
	}
	return data, nil
}

func (reader *indexBatchReader) Close() error {
	_ = reader.stdin.Close()
	err := reader.command.Wait()
	if err != nil {
		return fmt.Errorf("git_index_batch_failed: %w: %s", err, strings.TrimSpace(reader.stderr.String()))
	}
	return nil
}

func localWorktreeBoundary(root string) BoundaryClassifier {
	return func(blob GitBlob) (string, bool) {
		parts := strings.Split(blob.Path, "/")
		current := root
		for index, part := range parts {
			current = filepath.Join(current, filepath.FromSlash(part))
			info, err := os.Lstat(current)
			if err == nil && info.Mode()&os.ModeSymlink != 0 {
				return "worktree-symlink", true
			}
			if index >= len(parts)-1 {
				continue
			}
			gitMarker := filepath.Join(current, ".git")
			if _, markerErr := os.Lstat(gitMarker); markerErr == nil {
				return "nested-repository", true
			} else if !os.IsNotExist(markerErr) {
				return "nested-repository-marker-unreadable", true
			}
		}
		return "", false
	}
}
