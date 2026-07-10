//go:build !windows

package release

import "syscall"

func processExists(pid int) (bool, error) {
	err := syscall.Kill(pid, 0)
	if err == nil || err == syscall.EPERM {
		return true, nil
	}
	if err == syscall.ESRCH {
		return false, nil
	}
	return false, err
}
