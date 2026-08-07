package portfolio

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/plancatalog"
)

var mirrorNamePattern = regexp.MustCompile(`[^a-z0-9-]+`)

type MirrorManager struct {
	RepositoryRoot string
	Root           string
	Runner         CommandRunner
	installHook    func(phase, stagedPath string) error
}

type GitRepository struct {
	Path          string
	DefaultBranch string
	DefaultRef    string
	Source        RepositorySource
	runner        CommandRunner
	rootPath      string
	rootHandle    *os.Root
	mirrorName    string
	mirrorHandle  *os.Root
	mirrorFile    *os.File
	refreshHandle *os.Root
}

type GitTreeFile struct {
	Path     string
	Mode     string
	ObjectID string
	Revision string
}

type GitTreeChange struct {
	Status      string
	OldPath     string
	NewPath     string
	OldMode     string
	NewMode     string
	OldObjectID string
	NewObjectID string
}

func (manager MirrorManager) Refresh(ctx context.Context, source RepositorySource) (returnedRepository *GitRepository, returnedError error) {
	runner := manager.Runner
	if runner == nil {
		runner = ExecRunner{}
	}
	if err := validateRemoteURL(source.CloneURL); err != nil {
		return nil, err
	}
	repositoryRoot, root, rootHandle, err := manager.openRoot()
	if err != nil {
		return nil, err
	}
	defer repositoryRoot.Close()
	retainRootHandle := false
	defer func() {
		if !retainRootHandle {
			_ = rootHandle.Close()
		}
	}()
	temporaryName, err := randomLocalName(".refresh-")
	if err != nil {
		return nil, err
	}
	if err := rootHandle.Mkdir(temporaryName, 0o700); err != nil {
		return nil, err
	}
	temporary := filepath.Join(root, temporaryName)
	temporaryHandle, err := rootHandle.OpenRoot(temporaryName)
	if err != nil {
		return nil, fmt.Errorf("mirror_root_unsafe: refresh directory could not be bound")
	}
	temporaryFile, err := rootHandle.OpenFile(temporaryName, os.O_RDONLY, 0)
	if err != nil {
		temporaryHandle.Close()
		return nil, fmt.Errorf("mirror_root_unsafe: refresh directory descriptor could not be retained")
	}
	defer func() {
		cleanupErr := cleanupStagingDirectory(manager, rootHandle, temporaryName, temporary, temporaryHandle, temporaryFile)
		if cleanupErr == nil {
			return
		}
		if returnedRepository != nil {
			returnedRepository.Close()
			returnedRepository = nil
		}
		if returnedError != nil {
			returnedError = fmt.Errorf("%v; %w", returnedError, cleanupErr)
		} else {
			returnedError = cleanupErr
		}
	}()
	hooks := filepath.Join(temporary, "disabled-hooks")
	if err := temporaryHandle.Mkdir("disabled-hooks", 0o700); err != nil {
		return nil, err
	}
	bare := filepath.Join(temporary, "mirror.git")
	if err := temporaryHandle.Mkdir("mirror.git", 0o700); err != nil {
		return nil, err
	}
	bareHandle, err := temporaryHandle.OpenRoot("mirror.git")
	if err != nil {
		return nil, fmt.Errorf("mirror_root_unsafe: refresh mirror could not be bound")
	}
	bareHandleOpen := true
	defer func() {
		if bareHandleOpen {
			_ = bareHandle.Close()
		}
	}()
	bareFile, err := temporaryHandle.OpenFile("mirror.git", os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("mirror_root_unsafe: refresh mirror descriptor could not be retained")
	}
	bareFileOpen := true
	defer func() {
		if bareFileOpen {
			_ = bareFile.Close()
		}
	}()
	runBoundGit := func(args []string) (CommandResult, error) {
		if err := verifyRefreshMirrorAuthority(root, rootHandle, temporaryName, temporaryHandle, bareHandle, bareFile); err != nil {
			return CommandResult{}, err
		}
		result, runErr := runBoundCommand(ctx, runner, bareFile, bare, "git", args...)
		if err := verifyRefreshMirrorAuthority(root, rootHandle, temporaryName, temporaryHandle, bareHandle, bareFile); err != nil {
			return result, err
		}
		return result, runErr
	}
	if _, err := runBoundGit(safeGitArgs("-c", "core.hooksPath="+hooks, "init", "--bare", ".")); err != nil {
		return nil, fmt.Errorf("mirror_init_failed: scanner-owned bare mirror initialization failed")
	}
	if _, err := runBoundGit(safeGitArgs("-c", "core.hooksPath="+hooks, "remote", "add", "origin", source.CloneURL)); err != nil {
		return nil, fmt.Errorf("mirror_remote_failed: scanner-owned remote configuration failed")
	}
	if _, err := runBoundGit(safeGitRemoteArgs(source.CloneURL, "-c", "core.hooksPath="+hooks, "fetch", "--prune", "--force", "--no-tags", "origin", "+refs/heads/*:refs/remotes/origin/*")); err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("mirror_refresh_cancelled: remote refresh was cancelled")
		}
		return nil, fmt.Errorf("mirror_refresh_failed: remote refs could not be refreshed completely")
	}
	defaultBranch := strings.TrimPrefix(strings.TrimSpace(source.DefaultBranch), "refs/heads/")
	if defaultBranch == "" {
		remote, runErr := runBoundGit(safeGitRemoteArgs(source.CloneURL, "-c", "core.hooksPath="+hooks, "ls-remote", "--symref", "origin", "HEAD"))
		if runErr != nil {
			return nil, fmt.Errorf("default_branch_unavailable: remote HEAD could not be resolved")
		}
		defaultBranch = parseDefaultBranch(remote.Stdout)
	}
	if defaultBranch == "" {
		return nil, fmt.Errorf("default_branch_unavailable: remote HEAD did not identify a branch")
	}
	defaultRef := "refs/remotes/origin/" + defaultBranch
	if _, err := runBoundGit(safeGitArgs("rev-parse", "--verify", defaultRef+"^{commit}")); err != nil {
		return nil, fmt.Errorf("default_branch_unavailable: fetched default ref is missing")
	}
	if err := verifyRefreshMirrorAuthority(root, rootHandle, temporaryName, temporaryHandle, bareHandle, bareFile); err != nil {
		return nil, err
	}
	refreshedInfo, err := bareFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("mirror_replace_failed: refreshed mirror identity became unavailable")
	}
	if manager.installHook != nil {
		if err := manager.installHook("mirror-ready-for-install", bare); err != nil {
			return nil, err
		}
	}
	finalName := mirrorDirectoryName(source)
	final := filepath.Join(root, finalName)
	oldName := ""
	var priorHandle *os.Root
	var priorFile *os.File
	var priorSnapshot []byte
	priorQuarantined := false
	if info, statErr := rootHandle.Lstat(finalName); statErr == nil {
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("mirror_replace_failed: existing mirror is not a regular directory")
		}
		priorHandle, err = rootHandle.OpenRoot(finalName)
		if err != nil {
			return nil, fmt.Errorf("mirror_replace_failed: existing mirror could not be bound")
		}
		priorFile, err = rootHandle.OpenFile(finalName, os.O_RDONLY, 0)
		if err != nil {
			priorHandle.Close()
			return nil, fmt.Errorf("mirror_replace_failed: existing mirror descriptor could not be retained")
		}
		if err := verifyRetainedMirrorEntry(rootHandle, finalName, priorHandle, priorFile); err != nil {
			priorFile.Close()
			priorHandle.Close()
			return nil, err
		}
		priorSnapshot, err = snapshotMirrorTree(priorHandle)
		if err != nil {
			priorFile.Close()
			priorHandle.Close()
			return nil, fmt.Errorf("mirror_replace_failed: existing mirror snapshot unavailable: %w", err)
		}
		oldName, err = randomLocalName(".previous-")
		if err != nil {
			priorFile.Close()
			priorHandle.Close()
			return nil, err
		}
		if err := renameDirectoryNoReplace(root, rootHandle, finalName, oldName); err != nil {
			priorFile.Close()
			priorHandle.Close()
			return nil, fmt.Errorf("mirror_replace_failed: existing mirror could not be quarantined")
		}
		priorQuarantined = true
	} else if !os.IsNotExist(statErr) {
		return nil, statErr
	}
	installed := false
	var finalHandle *os.Root
	rollback := func(primary error) error {
		failures := []string{}
		if finalHandle != nil {
			if rollbackErr := finalHandle.Close(); rollbackErr != nil {
				failures = append(failures, "installed mirror handle could not be closed: "+rollbackErr.Error())
			}
			finalHandle = nil
		}
		if installed {
			capturedName, captured, rollbackErr := captureRetainedMirrorEntry(rootHandle, finalName, ".rollback-suspect-", bareHandle, bareFile)
			if captured {
				installed = false
			}
			if rollbackErr != nil {
				failures = append(failures, "suspect installation capture failed: "+rollbackErr.Error())
			} else if rollbackErr := rootHandle.RemoveAll(capturedName); rollbackErr != nil {
				failures = append(failures, "suspect installation could not be removed: "+rollbackErr.Error())
			}
		}
		if priorQuarantined {
			capturedName, captured, rollbackErr := captureRetainedMirrorEntry(rootHandle, oldName, ".rollback-prior-", priorHandle, priorFile)
			if captured {
				priorQuarantined = false
			}
			if rollbackErr != nil {
				failures = append(failures, "prior mirror capture failed: "+rollbackErr.Error())
			} else if _, occupiedErr := rootHandle.Lstat(finalName); !os.IsNotExist(occupiedErr) {
				failures = append(failures, "prior mirror could not be restored because the canonical name is occupied")
			} else {
				restoreHookFailed := false
				if manager.installHook != nil {
					if hookErr := manager.installHook("mirror-before-no-replace-restore", final); hookErr != nil {
						failures = append(failures, "prior mirror restore hook failed: "+hookErr.Error())
						restoreHookFailed = true
					}
				}
				if !restoreHookFailed {
					if rollbackErr := renameDirectoryNoReplace(root, rootHandle, capturedName, finalName); rollbackErr != nil {
						failures = append(failures, "captured prior mirror could not be restored: "+rollbackErr.Error())
					} else if rollbackErr := verifyRetainedMirrorEntry(rootHandle, finalName, priorHandle, priorFile); rollbackErr != nil {
						failures = append(failures, "restored prior mirror identity changed: "+rollbackErr.Error())
					} else if restoredSnapshot, rollbackErr := snapshotMirrorTree(priorHandle); rollbackErr != nil || !bytes.Equal(priorSnapshot, restoredSnapshot) {
						failures = append(failures, "restored prior mirror bytes differ from the quarantined snapshot")
					} else {
						priorQuarantined = false
					}
				}
			}
		}
		if priorFile != nil {
			if rollbackErr := priorFile.Close(); rollbackErr != nil {
				failures = append(failures, "prior mirror descriptor could not be closed: "+rollbackErr.Error())
			}
			priorFile = nil
		}
		if priorHandle != nil {
			if rollbackErr := priorHandle.Close(); rollbackErr != nil {
				failures = append(failures, "prior mirror handle could not be closed: "+rollbackErr.Error())
			}
			priorHandle = nil
		}
		if len(failures) != 0 {
			return fmt.Errorf("%v; mirror_rollback_failed: %s", primary, strings.Join(failures, "; "))
		}
		return primary
	}
	if manager.installHook != nil {
		if hookErr := manager.installHook("mirror-before-no-replace-install", final); hookErr != nil {
			return nil, rollback(fmt.Errorf("mirror_replace_failed: no-replace install precondition failed: %w", hookErr))
		}
	}
	if err := renameDirectoryNoReplace(root, rootHandle, filepath.ToSlash(filepath.Join(temporaryName, "mirror.git")), finalName); err != nil {
		return nil, rollback(fmt.Errorf("mirror_replace_failed: refreshed mirror could not be installed"))
	}
	installed = true
	finalHandle, err = rootHandle.OpenRoot(finalName)
	if err != nil {
		return nil, rollback(fmt.Errorf("mirror_replace_failed: installed mirror could not be bound"))
	}
	if err := verifyMirrorAuthority(root, rootHandle, finalName, finalHandle); err != nil {
		return nil, rollback(err)
	}
	installedInfo, err := finalHandle.Stat(".")
	if err != nil || !os.SameFile(refreshedInfo, installedInfo) {
		return nil, rollback(fmt.Errorf("mirror_replace_failed: installed mirror identity differs from refreshed evidence"))
	}
	if manager.installHook != nil {
		if hookErr := manager.installHook("mirror-before-retained-authority-check", final); hookErr != nil {
			return nil, rollback(fmt.Errorf("mirror_replace_failed: retained refresh authority check failed: %w", hookErr))
		}
	}
	if err := verifyRetainedMirrorEntry(rootHandle, finalName, bareHandle, bareFile); err != nil {
		return nil, rollback(err)
	}
	boundInfo, boundErr := bareHandle.Stat(".")
	retainedInfo, retainedErr := bareFile.Stat()
	if boundErr != nil || retainedErr != nil || !os.SameFile(installedInfo, boundInfo) || !os.SameFile(installedInfo, retainedInfo) {
		return nil, rollback(fmt.Errorf("mirror_replace_failed: installed mirror differs from retained refresh authority"))
	}
	abortCommitted := func(primary error) error {
		failures := []string{}
		if finalHandle != nil {
			if closeErr := finalHandle.Close(); closeErr != nil {
				failures = append(failures, "installed mirror handle could not be closed: "+closeErr.Error())
			}
			finalHandle = nil
		}
		if priorFile != nil {
			if closeErr := priorFile.Close(); closeErr != nil {
				failures = append(failures, "prior mirror descriptor could not be closed: "+closeErr.Error())
			}
			priorFile = nil
		}
		if priorHandle != nil {
			if closeErr := priorHandle.Close(); closeErr != nil {
				failures = append(failures, "prior mirror handle could not be closed: "+closeErr.Error())
			}
			priorHandle = nil
		}
		if bareFileOpen {
			if closeErr := bareFile.Close(); closeErr != nil {
				failures = append(failures, "installed mirror descriptor could not be closed: "+closeErr.Error())
			}
			bareFileOpen = false
		}
		if bareHandleOpen {
			if closeErr := bareHandle.Close(); closeErr != nil {
				failures = append(failures, "refresh mirror handle could not be closed: "+closeErr.Error())
			}
			bareHandleOpen = false
		}
		if len(failures) != 0 {
			return fmt.Errorf("%v; mirror_cleanup_failed: %s", primary, strings.Join(failures, "; "))
		}
		return primary
	}
	if priorQuarantined {
		if err := verifyRetainedMirrorEntry(rootHandle, oldName, priorHandle, priorFile); err != nil {
			return nil, rollback(err)
		}
		if currentSnapshot, err := snapshotMirrorTree(priorHandle); err != nil || !bytes.Equal(priorSnapshot, currentSnapshot) {
			return nil, rollback(fmt.Errorf("mirror_replace_failed: quarantined prior mirror bytes changed before cleanup"))
		}
		if manager.installHook != nil {
			if hookErr := manager.installHook("mirror-before-prior-cleanup", filepath.Join(root, oldName)); hookErr != nil {
				return nil, rollback(fmt.Errorf("mirror_replace_failed: prior mirror cleanup precondition failed: %w", hookErr))
			}
		}
		if err := verifyRetainedMirrorEntry(rootHandle, finalName, bareHandle, bareFile); err != nil {
			return nil, rollback(err)
		}
		if err := verifyRetainedMirrorEntry(rootHandle, oldName, priorHandle, priorFile); err != nil {
			return nil, rollback(err)
		}
		capturedName, captured, err := captureRetainedMirrorEntry(rootHandle, oldName, ".cleanup-prior-", priorHandle, priorFile)
		if captured {
			priorQuarantined = false
		}
		if err != nil {
			return nil, abortCommitted(fmt.Errorf("mirror_cleanup_failed: prior mirror capture failed: %w", err))
		}
		if err := rootHandle.RemoveAll(capturedName); err != nil {
			return nil, abortCommitted(fmt.Errorf("mirror_cleanup_failed: quarantined prior mirror cleanup failed: %w", err))
		}
		if err := priorFile.Close(); err != nil {
			return nil, abortCommitted(fmt.Errorf("mirror_cleanup_failed: prior mirror descriptor could not be closed: %w", err))
		}
		priorFile = nil
		if err := priorHandle.Close(); err != nil {
			return nil, abortCommitted(fmt.Errorf("mirror_cleanup_failed: prior mirror handle could not be closed: %w", err))
		}
		priorHandle = nil
	}
	if err := verifyRetainedMirrorEntry(rootHandle, finalName, bareHandle, bareFile); err != nil {
		return nil, rollback(err)
	}
	bareHandleOpen = false
	bareFileOpen = false
	retainRootHandle = true
	return &GitRepository{
		Path: final, DefaultBranch: defaultBranch, DefaultRef: defaultRef, Source: source, runner: runner,
		rootPath: root, rootHandle: rootHandle, mirrorName: finalName, mirrorHandle: finalHandle, mirrorFile: bareFile, refreshHandle: bareHandle,
	}, nil
}

func (manager MirrorManager) openRoot() (*os.Root, string, *os.Root, error) {
	if strings.TrimSpace(manager.RepositoryRoot) == "" {
		return nil, "", nil, fmt.Errorf("mirror_root_unsafe: repository root is required")
	}
	repositoryPath, err := filepath.Abs(manager.RepositoryRoot)
	if err != nil {
		return nil, "", nil, err
	}
	rootPath, err := filepath.Abs(manager.Root)
	if err != nil {
		return nil, "", nil, err
	}
	relative, err := filepath.Rel(repositoryPath, rootPath)
	if err != nil || relative == "." || relative == "" || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return nil, "", nil, fmt.Errorf("mirror_root_unsafe: scanner mirror root must be below the repository root")
	}
	repositoryRoot, err := os.OpenRoot(repositoryPath)
	if err != nil {
		return nil, "", nil, err
	}
	relative = filepath.FromSlash(filepath.ToSlash(relative))
	if err := ensureLocalDirectoryChain(repositoryRoot, relative); err != nil {
		repositoryRoot.Close()
		return nil, "", nil, fmt.Errorf("mirror_root_unsafe: %w", err)
	}
	rootHandle, err := repositoryRoot.OpenRoot(relative)
	if err != nil {
		repositoryRoot.Close()
		return nil, "", nil, fmt.Errorf("mirror_root_unsafe: scanner mirror root could not be bound")
	}
	if err := verifyBoundDirectory(rootPath, rootHandle); err != nil {
		rootHandle.Close()
		repositoryRoot.Close()
		return nil, "", nil, err
	}
	return repositoryRoot, rootPath, rootHandle, nil
}

func verifyMirrorAuthority(rootPath string, root *os.Root, childName string, child *os.Root) error {
	if err := verifyBoundDirectory(rootPath, root); err != nil {
		return err
	}
	pathInfo, err := root.Lstat(childName)
	if err != nil || !pathInfo.IsDir() || pathInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("mirror_root_changed: refresh directory authority changed")
	}
	boundInfo, err := child.Stat(".")
	if err != nil || !os.SameFile(pathInfo, boundInfo) {
		return fmt.Errorf("mirror_root_changed: refresh directory identity changed")
	}
	return nil
}

func verifyRefreshMirrorAuthority(rootPath string, root *os.Root, temporaryName string, temporary, mirror *os.Root, mirrorFile *os.File) error {
	if err := verifyMirrorAuthority(rootPath, root, temporaryName, temporary); err != nil {
		return err
	}
	pathInfo, err := temporary.Lstat("mirror.git")
	if err != nil || !pathInfo.IsDir() || pathInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("mirror_root_changed: refresh mirror authority changed")
	}
	boundInfo, err := mirror.Stat(".")
	if err != nil || !os.SameFile(pathInfo, boundInfo) {
		return fmt.Errorf("mirror_root_changed: refresh mirror identity changed")
	}
	fileInfo, err := mirrorFile.Stat()
	if err != nil || !os.SameFile(pathInfo, fileInfo) {
		return fmt.Errorf("mirror_root_changed: refresh mirror descriptor identity changed")
	}
	return nil
}

func verifyRetainedMirrorEntry(root *os.Root, name string, retainedRoot *os.Root, retainedFile *os.File) error {
	if retainedRoot == nil || retainedFile == nil {
		return fmt.Errorf("mirror_root_changed: retained mirror authority is unavailable")
	}
	pathInfo, err := root.Lstat(name)
	if err != nil || !pathInfo.IsDir() || pathInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("mirror_root_changed: retained mirror namespace changed")
	}
	rootInfo, rootErr := retainedRoot.Stat(".")
	fileInfo, fileErr := retainedFile.Stat()
	if rootErr != nil || fileErr != nil || !os.SameFile(pathInfo, rootInfo) || !os.SameFile(pathInfo, fileInfo) {
		return fmt.Errorf("mirror_root_changed: retained mirror identity changed")
	}
	return nil
}

func captureRetainedMirrorEntry(root *os.Root, name, prefix string, retainedRoot *os.Root, retainedFile *os.File) (string, bool, error) {
	capturedName, err := randomLocalName(prefix)
	if err != nil {
		return "", false, err
	}
	if _, statErr := root.Lstat(capturedName); !os.IsNotExist(statErr) {
		if statErr == nil {
			return "", false, fmt.Errorf("mirror_capture_failed: random capture name already exists")
		}
		return "", false, statErr
	}
	if err := renameDirectoryNoReplace(root.Name(), root, name, capturedName); err != nil {
		return "", false, err
	}
	if err := verifyRetainedMirrorEntry(root, capturedName, retainedRoot, retainedFile); err != nil {
		return capturedName, true, err
	}
	return capturedName, true, nil
}

func cleanupStagingDirectory(manager MirrorManager, root *os.Root, temporaryName, temporaryPath string, retainedRoot *os.Root, retainedFile *os.File) error {
	failures := []string{}
	if manager.installHook != nil {
		if err := manager.installHook("mirror-before-staging-cleanup", temporaryPath); err != nil {
			failures = append(failures, "staging cleanup hook failed: "+err.Error())
		}
	}
	capturedName, _, captureErr := captureRetainedMirrorEntry(root, temporaryName, ".cleanup-refresh-", retainedRoot, retainedFile)
	if captureErr != nil {
		failures = append(failures, "staging directory capture failed: "+captureErr.Error())
	} else if err := root.RemoveAll(capturedName); err != nil {
		failures = append(failures, "captured staging directory could not be removed: "+err.Error())
	}
	if err := retainedFile.Close(); err != nil {
		failures = append(failures, "staging directory descriptor could not be closed: "+err.Error())
	}
	if err := retainedRoot.Close(); err != nil {
		failures = append(failures, "staging directory handle could not be closed: "+err.Error())
	}
	if len(failures) != 0 {
		return fmt.Errorf("mirror_staging_cleanup_failed: %s", strings.Join(failures, "; "))
	}
	return nil
}

func snapshotMirrorTree(root *os.Root) ([]byte, error) {
	var snapshot bytes.Buffer
	err := fs.WalkDir(root.FS(), ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(&snapshot, "%s\x00%d\x00", path, uint32(info.Mode()))
		switch {
		case info.Mode().IsRegular():
			data, err := root.ReadFile(path)
			if err != nil {
				return err
			}
			digest := sha256.Sum256(data)
			fmt.Fprintf(&snapshot, "%d\x00%s\n", len(data), hex.EncodeToString(digest[:]))
		case info.IsDir():
			snapshot.WriteByte('\n')
		case info.Mode()&os.ModeSymlink != 0:
			target, err := root.Readlink(path)
			if err != nil {
				return err
			}
			snapshot.WriteString(target)
			snapshot.WriteByte('\n')
		default:
			return fmt.Errorf("unsupported mirror entry type at %s", path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return snapshot.Bytes(), nil
}

func verifyBoundDirectory(path string, root *os.Root) error {
	pathInfo, err := os.Lstat(path)
	if err != nil || !pathInfo.IsDir() || pathInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("mirror_root_changed: scanner mirror root authority changed")
	}
	boundInfo, err := root.Stat(".")
	if err != nil || !os.SameFile(pathInfo, boundInfo) {
		return fmt.Errorf("mirror_root_changed: scanner mirror root identity changed")
	}
	return nil
}

func validateRemoteURL(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("remote_url_invalid: empty remote")
	}
	if filepath.IsAbs(value) {
		return nil
	}
	if !strings.Contains(value, "://") {
		at := strings.LastIndex(value, "@")
		colon := strings.Index(value, ":")
		if at > 0 && colon > at && value[:at] == "git" && strings.TrimSpace(value[at+1:colon]) != "" && strings.TrimSpace(value[colon+1:]) != "" {
			return nil
		}
		return fmt.Errorf("remote_url_protocol_forbidden: remote must use HTTPS, SSH, or an absolute file path")
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" {
		return fmt.Errorf("remote_url_invalid: remote URL could not be parsed")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return fmt.Errorf("remote_url_credential_forbidden: query-bearing remotes are not persisted")
	}
	switch strings.ToLower(parsed.Scheme) {
	case "https":
		if parsed.Hostname() == "" || parsed.Path == "" || parsed.User != nil {
			return fmt.Errorf("remote_url_credential_forbidden: credential-bearing remotes are not persisted")
		}
	case "ssh":
		if parsed.Hostname() == "" || parsed.Path == "" {
			return fmt.Errorf("remote_url_invalid: SSH remote requires host and path")
		}
		if parsed.User != nil {
			_, hasPassword := parsed.User.Password()
			if hasPassword || parsed.User.Username() != "git" {
				return fmt.Errorf("remote_url_credential_forbidden: SSH remotes may use only the git username without a password")
			}
		}
	case "file":
		if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.Path == "" {
			return fmt.Errorf("remote_url_invalid: file remote is invalid")
		}
	default:
		return fmt.Errorf("remote_url_protocol_forbidden: remote must use HTTPS, SSH, or file transport")
	}
	return nil
}

func safeGitArgs(args ...string) []string {
	policy := []string{
		"-c", "protocol.allow=never",
		"-c", "protocol.https.allow=always",
		"-c", "protocol.ssh.allow=always",
		"-c", "protocol.file.allow=always",
		"-c", "protocol.ext.allow=never",
		"-c", "core.longpaths=true",
		"-c", "core.sshCommand=ssh",
		"-c", "core.hooksPath=" + os.DevNull,
	}
	return append(policy, args...)
}

func safeGitRemoteArgs(remote string, args ...string) []string {
	pinned := []string{
		"-c", "url." + remote + ".insteadOf=" + remote,
		"-c", "remote.origin.uploadpack=git-upload-pack",
	}
	return safeGitArgs(append(pinned, args...)...)
}

func mirrorDirectoryName(source RepositorySource) string {
	base := strings.ToLower(source.NameHint)
	base = mirrorNamePattern.ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")
	if base == "" {
		base = "repository"
	}
	digest := sha256.Sum256([]byte(source.Kind + "\x00" + source.Locator + "\x00" + source.CloneURL))
	return base + "-" + hex.EncodeToString(digest[:6]) + ".git"
}

func parseDefaultBranch(data []byte) string {
	for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 3 && fields[0] == "ref:" && fields[2] == "HEAD" {
			return strings.TrimPrefix(fields[1], "refs/heads/")
		}
	}
	return ""
}

func (repository *GitRepository) git(ctx context.Context, args ...string) (CommandResult, error) {
	if err := repository.verifyAuthority(); err != nil {
		return CommandResult{}, err
	}
	result, err := runBoundCommand(ctx, repository.runner, repository.mirrorFile, repository.Path, "git", safeGitArgs(args...)...)
	if err != nil {
		return result, err
	}
	if err := repository.verifyAuthority(); err != nil {
		return result, err
	}
	return result, nil
}

func (repository *GitRepository) verifyAuthority() error {
	if repository.rootHandle == nil || repository.mirrorHandle == nil || repository.mirrorFile == nil {
		return nil
	}
	if err := verifyMirrorAuthority(repository.rootPath, repository.rootHandle, repository.mirrorName, repository.mirrorHandle); err != nil {
		return err
	}
	pathInfo, err := repository.rootHandle.Lstat(repository.mirrorName)
	if err != nil {
		return fmt.Errorf("mirror_root_changed: installed mirror authority changed")
	}
	fileInfo, err := repository.mirrorFile.Stat()
	if err != nil || !os.SameFile(pathInfo, fileInfo) {
		return fmt.Errorf("mirror_root_changed: installed mirror descriptor identity changed")
	}
	return nil
}

func (repository *GitRepository) Close() {
	if repository.refreshHandle != nil {
		_ = repository.refreshHandle.Close()
		repository.refreshHandle = nil
	}
	if repository.mirrorFile != nil {
		_ = repository.mirrorFile.Close()
		repository.mirrorFile = nil
	}
	if repository.mirrorHandle != nil {
		_ = repository.mirrorHandle.Close()
		repository.mirrorHandle = nil
	}
	if repository.rootHandle != nil {
		_ = repository.rootHandle.Close()
		repository.rootHandle = nil
	}
}

func (repository *GitRepository) Revision(ctx context.Context, ref string) (string, error) {
	result, err := repository.git(ctx, "rev-parse", "--verify", ref+"^{commit}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result.Stdout)), nil
}

func (repository *GitRepository) RemoteRefs(ctx context.Context) ([]string, error) {
	result, err := repository.git(ctx, "for-each-ref", "--format=%(refname)", "refs/remotes/origin")
	if err != nil {
		return nil, err
	}
	refs := []string{}
	for _, ref := range strings.Fields(string(result.Stdout)) {
		if ref == repository.DefaultRef || strings.HasSuffix(ref, "/HEAD") {
			continue
		}
		refs = append(refs, ref)
	}
	sort.Strings(refs)
	return refs, nil
}

func (repository *GitRepository) ListFiles(ctx context.Context, ref string) ([]GitTreeFile, error) {
	result, err := repository.git(ctx, "ls-tree", "-r", "-z", "--full-tree", ref)
	if err != nil {
		return nil, err
	}
	files := []GitTreeFile{}
	for _, entry := range strings.Split(string(result.Stdout), "\x00") {
		if entry == "" {
			continue
		}
		metadata, path, found := strings.Cut(entry, "\t")
		fields := strings.Fields(metadata)
		if !found || len(fields) != 3 {
			return nil, fmt.Errorf("git_tree_invalid: malformed ls-tree record")
		}
		files = append(files, GitTreeFile{Path: path, Mode: fields[0], ObjectID: fields[2], Revision: ref})
	}
	sort.SliceStable(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func (repository *GitRepository) ReadFile(ctx context.Context, ref, path string) ([]byte, error) {
	result, err := repository.git(ctx, "ls-tree", "-z", ref, "--", path)
	if err != nil {
		return nil, err
	}
	var selected *GitTreeFile
	for _, entry := range strings.Split(string(result.Stdout), "\x00") {
		if entry == "" {
			continue
		}
		metadata, candidatePath, found := strings.Cut(entry, "\t")
		fields := strings.Fields(metadata)
		if !found || len(fields) != 3 {
			return nil, fmt.Errorf("git_tree_invalid: malformed ls-tree record")
		}
		if candidatePath == path {
			file := GitTreeFile{Path: candidatePath, Mode: fields[0], ObjectID: fields[2], Revision: ref}
			selected = &file
		}
	}
	if selected == nil || (selected.Mode != string(plancatalog.GitModeRegular) && selected.Mode != string(plancatalog.GitModeExecutable)) {
		return nil, os.ErrNotExist
	}
	return repository.ReadBlob(ctx, *selected)
}

func (repository *GitRepository) ReadBlob(ctx context.Context, file GitTreeFile) ([]byte, error) {
	if file.Mode != string(plancatalog.GitModeRegular) && file.Mode != string(plancatalog.GitModeExecutable) {
		return nil, fmt.Errorf("remote_blob_mode_unreadable: %s has mode %s", file.Path, file.Mode)
	}
	sizeResult, err := repository.git(ctx, "cat-file", "-s", file.ObjectID)
	if err != nil {
		return nil, err
	}
	size, err := strconv.ParseInt(strings.TrimSpace(string(sizeResult.Stdout)), 10, 64)
	if err != nil || size < 0 {
		return nil, fmt.Errorf("remote_blob_invalid_size: %s", file.Path)
	}
	if size > plancatalog.MaxPlanSourceBytes {
		return nil, &plancatalog.CodedError{Code: "source_too_large", Err: fmt.Errorf("remote source exceeds %d-byte plan-catalog limit: %s", plancatalog.MaxPlanSourceBytes, file.Path)}
	}
	result, err := repository.git(ctx, "cat-file", "blob", file.ObjectID)
	if err != nil {
		return nil, err
	}
	if int64(len(result.Stdout)) != size {
		return nil, fmt.Errorf("remote_blob_size_changed: %s", file.Path)
	}
	return result.Stdout, nil
}

func (repository *GitRepository) MergeBase(ctx context.Context, ref string) (string, error) {
	result, err := repository.git(ctx, "merge-base", repository.DefaultRef, ref)
	if err != nil {
		return "", err
	}
	base := strings.TrimSpace(string(result.Stdout))
	if base == "" {
		return "", fmt.Errorf("empty merge base")
	}
	return base, nil
}

func (repository *GitRepository) IsMerged(ctx context.Context, ref, mergeBase string) (bool, error) {
	tip, err := repository.Revision(ctx, ref)
	if err != nil {
		return false, err
	}
	defaultRevision, err := repository.Revision(ctx, repository.DefaultRef)
	if err != nil {
		return false, err
	}
	if mergeBase == tip || tip == defaultRevision {
		return true, nil
	}
	return false, nil
}

func (repository *GitRepository) ChangedPaths(ctx context.Context, mergeBase, ref string) ([]GitTreeChange, error) {
	result, err := repository.git(ctx, "diff", "--raw", "-z", "--find-renames", "--diff-filter=AMR", mergeBase+".."+ref, "--")
	if err != nil {
		return nil, err
	}
	changes := []GitTreeChange{}
	fields := strings.Split(string(result.Stdout), "\x00")
	for index := 0; index < len(fields); {
		header := fields[index]
		index++
		if header == "" {
			continue
		}
		metadata := strings.Fields(strings.TrimPrefix(header, ":"))
		if len(metadata) != 5 {
			return nil, fmt.Errorf("branch_change_invalid: malformed raw change header")
		}
		status := metadata[4]
		change := GitTreeChange{Status: status, OldMode: metadata[0], NewMode: metadata[1], OldObjectID: metadata[2], NewObjectID: metadata[3]}
		if strings.HasPrefix(status, "R") {
			if index+1 >= len(fields) {
				return nil, fmt.Errorf("branch_change_invalid: rename evidence is truncated")
			}
			change.OldPath = fields[index]
			change.NewPath = fields[index+1]
			index += 2
			changes = append(changes, change)
			continue
		}
		if status != "A" && status != "M" {
			return nil, fmt.Errorf("branch_change_invalid: unexpected change status %q", status)
		}
		if index >= len(fields) {
			return nil, fmt.Errorf("branch_change_invalid: changed path evidence is truncated")
		}
		change.NewPath = fields[index]
		if status == "M" {
			change.OldPath = change.NewPath
		}
		index++
		if change.NewPath == "" {
			continue
		}
		changes = append(changes, change)
	}
	sort.SliceStable(changes, func(i, j int) bool {
		if changes[i].NewPath != changes[j].NewPath {
			return changes[i].NewPath < changes[j].NewPath
		}
		return changes[i].OldPath < changes[j].OldPath
	})
	return changes, nil
}

func (repository *GitRepository) CommitTime(ctx context.Context, ref string) (string, error) {
	result, err := repository.git(ctx, "show", "-s", "--format=%cI", ref)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result.Stdout)), nil
}

func (repository *GitRepository) PathCommitTime(ctx context.Context, ref, path string) (string, error) {
	result, err := repository.git(ctx, "log", "-1", "--format=%cI", ref, "--", filepath.ToSlash(path))
	if err != nil {
		return "", err
	}
	value := strings.TrimSpace(string(result.Stdout))
	if value == "" {
		return "", fmt.Errorf("path history is empty")
	}
	return value, nil
}
