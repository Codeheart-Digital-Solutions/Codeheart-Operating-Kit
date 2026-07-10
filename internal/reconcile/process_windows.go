//go:build windows

package reconcile

import (
	"syscall"
)

func processAlive(pid int) (bool, error) {
	const (
		processQueryLimitedInformation = 0x1000
		errorInvalidParameter          = syscall.Errno(87)
	)
	handle, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err == nil {
		_ = syscall.CloseHandle(handle)
		return true, nil
	}
	if err == errorInvalidParameter {
		return false, nil
	}
	return false, err
}
