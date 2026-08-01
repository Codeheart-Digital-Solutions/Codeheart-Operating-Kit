//go:build !darwin && !linux && !windows

package portfolio

import (
	"fmt"
	"os"
)

func renameDirectoryNoReplace(_ string, _ *os.Root, _, _ string) error {
	return fmt.Errorf("mirror_replace_failed: atomic no-replace directory publication is unsupported on this platform")
}
