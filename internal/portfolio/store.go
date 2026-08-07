package portfolio

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/plancatalog"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

func StoreCompleteCatalog(root string, catalog Catalog) (updated bool, previousPreserved bool, err error) {
	return storeCompleteCatalogWithHook(root, catalog, nil)
}

func storeCompleteCatalogWithHook(root string, catalog Catalog, hook func(string) error) (updated bool, previousPreserved bool, err error) {
	if !catalog.Complete {
		_, readErr := readRootRegular(root, CatalogPath)
		return false, readErr == nil, nil
	}
	if catalog.SchemaVersion != 2 || catalog.DiscoveryVersion != plancatalog.DiscoveryV2 {
		_, readErr := readRootRegular(root, CatalogPath)
		return false, readErr == nil, fmt.Errorf("portfolio_catalog_incompatible: only complete discovery-v2 evidence may replace the current cache")
	}
	schemaPath, err := state.SchemaForPlanCatalogVersion(catalog.SchemaVersion)
	if err != nil {
		return false, false, fmt.Errorf("portfolio_catalog_invalid: %w", err)
	}
	if err := state.Validate(schemaPath, catalog); err != nil {
		return false, false, fmt.Errorf("portfolio_catalog_invalid: %w", err)
	}
	data, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return false, false, err
	}
	data = append(data, '\n')
	parentName := filepath.ToSlash(filepath.Dir(filepath.FromSlash(CatalogPath)))
	bound, err := bindLocalDirectory(root, parentName, true)
	if err != nil {
		return false, false, err
	}
	defer bound.Close()
	parent := bound.directory
	target := filepath.Base(CatalogPath)
	temporaryName, temporary, err := createExclusiveLocalFile(parent, ".catalog-")
	if err != nil {
		return false, false, err
	}
	cleanupTemporary := true
	defer func() {
		if cleanupTemporary {
			_ = parent.Remove(temporaryName)
		}
	}()
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return false, false, err
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return false, false, err
	}
	temporaryInfo, err := temporary.Stat()
	if err != nil {
		_ = temporary.Close()
		return false, false, err
	}
	if err := temporary.Close(); err != nil {
		return false, false, err
	}
	if hook != nil {
		if err := hook("catalog-staged"); err != nil {
			return false, false, err
		}
	}
	if err := bound.Verify(); err != nil {
		cleanupTemporary = false
		return false, false, err
	}
	var existing *os.File
	var existingInfo os.FileInfo
	var existingData []byte
	existingPresent := false
	if info, statErr := parent.Lstat(target); statErr == nil {
		if !info.Mode().IsRegular() {
			return false, false, fmt.Errorf("portfolio_cache_unsafe: existing cache is not a regular file")
		}
		existingPresent = true
		existingInfo = info
		existing, err = parent.Open(target)
		if err != nil {
			return false, false, err
		}
		defer existing.Close()
		opened, statErr := existing.Stat()
		if statErr != nil || !os.SameFile(info, opened) {
			return false, false, fmt.Errorf("portfolio_cache_changed: existing cache identity changed while binding")
		}
		existingData, err = io.ReadAll(existing)
		if err != nil {
			return false, false, err
		}
	} else if !os.IsNotExist(statErr) {
		return false, false, statErr
	}
	if hook != nil {
		if err := hook("before-catalog-publish"); err != nil {
			return false, existingPresent, err
		}
	}
	if err := bound.Verify(); err != nil {
		cleanupTemporary = false
		return false, existingPresent, err
	}
	if existingPresent {
		current, statErr := parent.Lstat(target)
		if statErr != nil || !os.SameFile(existingInfo, current) {
			return false, true, fmt.Errorf("portfolio_cache_changed: existing cache identity changed before publication")
		}
		opened, statErr := existing.Stat()
		if statErr != nil || !os.SameFile(current, opened) {
			return false, true, fmt.Errorf("portfolio_cache_changed: existing cache descriptor identity changed")
		}
		if _, err := existing.Seek(0, io.SeekStart); err != nil {
			return false, true, err
		}
		currentData, readErr := io.ReadAll(existing)
		if readErr != nil || !bytes.Equal(existingData, currentData) {
			return false, true, fmt.Errorf("portfolio_cache_changed: existing cache bytes changed before publication")
		}
	} else if _, statErr := parent.Lstat(target); !os.IsNotExist(statErr) {
		if statErr == nil {
			return false, false, fmt.Errorf("portfolio_cache_changed: cache appeared before publication")
		}
		return false, false, statErr
	}
	if existingPresent && catalogTargetMustCloseBeforeReplace() {
		if err := existing.Close(); err != nil {
			return false, true, fmt.Errorf("portfolio_cache_install_failed: close existing cache authority: %w", err)
		}
	}
	// Root.Rename performs one same-directory old-to-new replacement. Readers
	// therefore observe either the previous complete file or the new complete
	// file; there is no remove-then-install gap.
	var renameErr error
	for attempt := 0; attempt < fileShareRetryAttempts; attempt++ {
		renameErr = parent.Rename(temporaryName, target)
		if renameErr == nil || !isTransientFileSharingError(renameErr) {
			break
		}
		if err := bound.Verify(); err != nil {
			cleanupTemporary = false
			return false, existingPresent, err
		}
		if existingPresent {
			current, statErr := parent.Lstat(target)
			if statErr != nil && isTransientFileSharingError(statErr) {
				time.Sleep(fileShareRetryDelay)
				continue
			}
			if statErr != nil || !os.SameFile(existingInfo, current) {
				return false, true, fmt.Errorf("portfolio_cache_changed: existing cache identity changed while publication was waiting")
			}
		} else if _, statErr := parent.Lstat(target); !os.IsNotExist(statErr) {
			if statErr != nil && isTransientFileSharingError(statErr) {
				time.Sleep(fileShareRetryDelay)
				continue
			}
			if statErr == nil {
				return false, false, fmt.Errorf("portfolio_cache_changed: cache appeared while publication was waiting")
			}
			return false, false, statErr
		}
		time.Sleep(fileShareRetryDelay)
	}
	if renameErr != nil {
		return false, existingPresent, fmt.Errorf("portfolio_cache_install_failed: %w", renameErr)
	}
	cleanupTemporary = false
	installedInfo, err := parent.Lstat(target)
	if err != nil || !installedInfo.Mode().IsRegular() || !os.SameFile(temporaryInfo, installedInfo) {
		return false, existingPresent, fmt.Errorf("portfolio_cache_install_failed: installed identity differs")
	}
	installedFile, err := parent.Open(target)
	if err != nil {
		return false, existingPresent, fmt.Errorf("portfolio_cache_install_failed: %w", err)
	}
	defer installedFile.Close()
	openedInstalled, statErr := installedFile.Stat()
	if statErr != nil || !os.SameFile(installedInfo, openedInstalled) {
		return false, existingPresent, fmt.Errorf("portfolio_cache_install_failed: installed descriptor identity differs")
	}
	installed, readErr := io.ReadAll(installedFile)
	if readErr != nil || !bytes.Equal(installed, data) {
		return false, existingPresent, fmt.Errorf("portfolio_cache_install_failed: installed bytes differ")
	}
	if err := bound.Verify(); err != nil {
		return false, existingPresent, err
	}
	return true, false, nil
}

type boundLocalDirectory struct {
	absoluteRoot   string
	repositoryRoot *os.Root
	rootIdentity   os.FileInfo
	relative       string
	directory      *os.Root
	identity       os.FileInfo
}

func bindLocalDirectory(root, relative string, create bool) (*boundLocalDirectory, error) {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	pathRootInfo, err := os.Lstat(absoluteRoot)
	if err != nil || !pathRootInfo.IsDir() || pathRootInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("portfolio_local_root_unsafe: repository root is not a regular directory")
	}
	repositoryRoot, err := os.OpenRoot(absoluteRoot)
	if err != nil {
		return nil, err
	}
	rootIdentity, err := repositoryRoot.Stat(".")
	if err != nil || !os.SameFile(pathRootInfo, rootIdentity) {
		repositoryRoot.Close()
		return nil, fmt.Errorf("portfolio_local_root_changed: repository root identity changed while binding")
	}
	relative = filepath.Clean(filepath.FromSlash(relative))
	if relative == "" {
		relative = "."
	}
	if filepath.IsAbs(relative) || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		repositoryRoot.Close()
		return nil, fmt.Errorf("portfolio_local_path_unsafe: local directory must stay below the repository root")
	}
	if create {
		err = ensureLocalDirectoryChain(repositoryRoot, relative)
	} else {
		err = validateLocalDirectoryChain(repositoryRoot, relative)
	}
	if err != nil {
		repositoryRoot.Close()
		return nil, err
	}
	identity, err := repositoryRoot.Lstat(relative)
	if err != nil || !identity.IsDir() || identity.Mode()&os.ModeSymlink != 0 {
		repositoryRoot.Close()
		return nil, fmt.Errorf("portfolio_local_path_unsafe: local directory is not a regular directory")
	}
	directory, err := repositoryRoot.OpenRoot(relative)
	if err != nil {
		repositoryRoot.Close()
		return nil, err
	}
	opened, err := directory.Stat(".")
	if err != nil || !os.SameFile(identity, opened) {
		directory.Close()
		repositoryRoot.Close()
		return nil, fmt.Errorf("portfolio_local_path_changed: local directory identity changed while binding")
	}
	return &boundLocalDirectory{absoluteRoot: absoluteRoot, repositoryRoot: repositoryRoot, rootIdentity: rootIdentity, relative: relative, directory: directory, identity: identity}, nil
}

func (bound *boundLocalDirectory) Verify() error {
	pathRootInfo, err := os.Lstat(bound.absoluteRoot)
	if err != nil || !pathRootInfo.IsDir() || pathRootInfo.Mode()&os.ModeSymlink != 0 || !os.SameFile(bound.rootIdentity, pathRootInfo) {
		return fmt.Errorf("portfolio_local_root_changed: canonical repository root authority changed")
	}
	openedRoot, err := bound.repositoryRoot.Stat(".")
	if err != nil || !os.SameFile(bound.rootIdentity, openedRoot) {
		return fmt.Errorf("portfolio_local_root_changed: bound repository root identity changed")
	}
	if err := validateLocalDirectoryChain(bound.repositoryRoot, bound.relative); err != nil {
		return err
	}
	pathInfo, err := bound.repositoryRoot.Lstat(bound.relative)
	if err != nil || !pathInfo.IsDir() || pathInfo.Mode()&os.ModeSymlink != 0 || !os.SameFile(bound.identity, pathInfo) {
		return fmt.Errorf("portfolio_cache_parent_changed: canonical local directory authority changed")
	}
	opened, err := bound.directory.Stat(".")
	if err != nil || !os.SameFile(bound.identity, opened) {
		return fmt.Errorf("portfolio_cache_parent_changed: bound local directory identity changed")
	}
	return nil
}

func (bound *boundLocalDirectory) Close() {
	_ = bound.directory.Close()
	_ = bound.repositoryRoot.Close()
}

type scanLock struct {
	bound    *boundLocalDirectory
	file     *os.File
	identity os.FileInfo
	name     string
}

func acquireScanLock(root string) (*scanLock, error) {
	parentName := filepath.ToSlash(filepath.Dir(filepath.FromSlash(ScanLockPath)))
	bound, err := bindLocalDirectory(root, parentName, true)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(filepath.FromSlash(ScanLockPath))
	lock, err := bound.directory.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
	if os.IsExist(err) {
		bound.Close()
		return nil, fmt.Errorf("portfolio_scan_in_progress: another scanner owns the local portfolio lock")
	}
	if err != nil {
		bound.Close()
		return nil, err
	}
	identity, err := lock.Stat()
	if err != nil {
		lock.Close()
		bound.Close()
		return nil, err
	}
	result := &scanLock{bound: bound, file: lock, identity: identity, name: name}
	if err := result.Verify(); err != nil {
		result.Close()
		return nil, err
	}
	return result, nil
}

func (lock *scanLock) Verify() error {
	if err := lock.bound.Verify(); err != nil {
		return err
	}
	current, err := lock.bound.directory.Lstat(lock.name)
	if err != nil || !current.Mode().IsRegular() || !os.SameFile(lock.identity, current) {
		return fmt.Errorf("portfolio_scan_lock_changed: scan lock namespace identity changed")
	}
	opened, err := lock.file.Stat()
	if err != nil || !os.SameFile(current, opened) {
		return fmt.Errorf("portfolio_scan_lock_changed: scan lock descriptor identity changed")
	}
	return nil
}

func (lock *scanLock) Close() {
	if lock.Verify() == nil {
		_ = lock.bound.directory.Remove(lock.name)
	}
	_ = lock.file.Close()
	lock.bound.Close()
}

func validateLocalDirectoryChain(root *os.Root, relative string) error {
	current := ""
	for _, part := range strings.Split(filepath.ToSlash(relative), "/") {
		if part == "" || part == "." {
			continue
		}
		if current == "" {
			current = part
		} else {
			current += "/" + part
		}
		info, err := root.Lstat(filepath.FromSlash(current))
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("portfolio_local_path_unsafe: %s is not a regular directory", current)
		}
	}
	return nil
}

func ensureLocalDirectoryChain(root *os.Root, relative string) error {
	current := ""
	for _, part := range strings.Split(filepath.ToSlash(relative), "/") {
		if part == "" || part == "." {
			continue
		}
		if current == "" {
			current = part
		} else {
			current += "/" + part
		}
		info, err := root.Lstat(filepath.FromSlash(current))
		if os.IsNotExist(err) {
			if err := root.Mkdir(filepath.FromSlash(current), 0o700); err != nil && !os.IsExist(err) {
				return err
			}
			info, err = root.Lstat(filepath.FromSlash(current))
		}
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("portfolio_local_path_unsafe: %s is not a regular directory", current)
		}
	}
	return nil
}

func createExclusiveLocalFile(root *os.Root, prefix string) (string, *os.File, error) {
	for attempt := 0; attempt < 8; attempt++ {
		name, err := randomLocalName(prefix)
		if err != nil {
			return "", nil, err
		}
		file, err := root.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o600)
		if os.IsExist(err) {
			continue
		}
		return name, file, err
	}
	return "", nil, fmt.Errorf("could not reserve a local portfolio file")
}

func randomLocalName(prefix string) (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(value), nil
}
