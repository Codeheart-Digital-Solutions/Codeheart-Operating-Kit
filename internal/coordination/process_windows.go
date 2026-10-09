//go:build windows

package coordination

import (
	"os"
	"os/signal"
	"syscall"
)

func processAlive(pid int) (bool, error) {
	const (
		processQueryLimitedInformation = 0x1000
		stillActive                    = 259
		errorInvalidParameter          = syscall.Errno(87)
	)
	handle, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err == errorInvalidParameter {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer syscall.CloseHandle(handle)
	var code uint32
	if err := syscall.GetExitCodeProcess(handle, &code); err != nil {
		return false, err
	}
	return code == stillActive, nil
}

func notifyInterrupt(ch chan<- os.Signal) {
	signal.Notify(ch, os.Interrupt)
}
