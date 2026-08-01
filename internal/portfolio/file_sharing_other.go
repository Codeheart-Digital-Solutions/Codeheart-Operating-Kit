//go:build !windows

package portfolio

func isTransientFileSharingError(error) bool { return false }

func catalogTargetMustCloseBeforeReplace() bool { return false }
