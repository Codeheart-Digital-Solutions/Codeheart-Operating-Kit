//go:build windows

package portfolio

import (
	"os"
	"path/filepath"
	"syscall"
)

func renameDirectoryNoReplace(rootPath string, retainedRoot *os.Root, oldName, newName string) error {
	rootFile, err := openRetainedRenameRoot(rootPath, retainedRoot, oldName, newName)
	if err != nil {
		return err
	}
	defer rootFile.Close()
	oldPointer, err := syscall.UTF16PtrFromString(filepath.Join(rootPath, filepath.FromSlash(oldName)))
	if err != nil {
		return err
	}
	newPointer, err := syscall.UTF16PtrFromString(filepath.Join(rootPath, filepath.FromSlash(newName)))
	if err != nil {
		return err
	}
	return syscall.MoveFile(oldPointer, newPointer)
}
