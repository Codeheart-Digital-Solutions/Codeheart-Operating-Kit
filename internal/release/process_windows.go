//go:build windows

package release

import "syscall"

func processExists(pid int) (bool, error) {
	const processQueryLimitedInformation = 0x1000
	handle, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err == nil {
		_ = syscall.CloseHandle(handle)
		return true, nil
	}
	if err == syscall.Errno(87) {
		return false, nil
	}
	return false, err
}
