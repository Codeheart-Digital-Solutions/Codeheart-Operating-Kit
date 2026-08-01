package portfolio

import (
	"fmt"
	"os"
	"path/filepath"
)

func openRetainedRenameRoot(rootPath string, retainedRoot *os.Root, oldName, newName string) (*os.File, error) {
	for _, name := range []string{oldName, newName} {
		if !filepath.IsLocal(name) || name == "." {
			return nil, fmt.Errorf("mirror_replace_failed: non-local rename path %q", name)
		}
	}
	rootFile, err := os.Open(rootPath)
	if err != nil {
		return nil, err
	}
	retainedInfo, retainedErr := retainedRoot.Stat(".")
	fileInfo, fileErr := rootFile.Stat()
	if retainedErr != nil || fileErr != nil || !os.SameFile(retainedInfo, fileInfo) {
		rootFile.Close()
		return nil, fmt.Errorf("mirror_root_changed: no-replace rename root identity changed")
	}
	return rootFile, nil
}
