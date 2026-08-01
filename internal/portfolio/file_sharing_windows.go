//go:build windows

package portfolio

import (
	"errors"
	"syscall"
)

func isTransientFileSharingError(err error) bool {
	const errorSharingViolation = syscall.Errno(32)
	return errors.Is(err, errorSharingViolation) || errors.Is(err, syscall.ERROR_ACCESS_DENIED)
}

func catalogTargetMustCloseBeforeReplace() bool { return true }
