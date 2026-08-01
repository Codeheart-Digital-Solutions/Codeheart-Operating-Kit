//go:build linux && (amd64 || arm64)

package portfolio

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

const linuxRenameNoReplace = 1

func renameDirectoryNoReplace(rootPath string, retainedRoot *os.Root, oldName, newName string) error {
	rootFile, err := openRetainedRenameRoot(rootPath, retainedRoot, oldName, newName)
	if err != nil {
		return err
	}
	defer rootFile.Close()
	oldPointer, err := syscall.BytePtrFromString(oldName)
	if err != nil {
		return err
	}
	newPointer, err := syscall.BytePtrFromString(newName)
	if err != nil {
		return err
	}
	_, _, errno := syscall.Syscall6(
		linuxRenameat2Syscall,
		rootFile.Fd(), uintptr(unsafe.Pointer(oldPointer)),
		rootFile.Fd(), uintptr(unsafe.Pointer(newPointer)),
		linuxRenameNoReplace, 0,
	)
	runtime.KeepAlive(rootFile)
	if errno != 0 {
		return errno
	}
	return nil
}
