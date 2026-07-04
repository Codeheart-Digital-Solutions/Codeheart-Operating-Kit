package kitfs

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	operatingkit "github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit"
)

// FS returns the embedded Operating Kit source-resource filesystem.
func FS() fs.FS {
	return operatingkit.EmbeddedResources
}

// Clean validates and normalizes a resource path for use with the embedded FS.
func Clean(name string) (string, error) {
	cleaned := path.Clean(strings.TrimPrefix(name, "/"))
	if cleaned == "." || strings.HasPrefix(cleaned, "../") || cleaned == ".." || !fs.ValidPath(cleaned) {
		return "", fmt.Errorf("invalid resource path %q", name)
	}
	return cleaned, nil
}

// ReadFile reads an embedded Operating Kit resource.
func ReadFile(name string) ([]byte, error) {
	cleaned, err := Clean(name)
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(FS(), cleaned)
}

// ReadText reads an embedded Operating Kit resource as UTF-8 text.
func ReadText(name string) (string, error) {
	data, err := ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Exists reports whether an embedded resource path exists.
func Exists(name string) bool {
	cleaned, err := Clean(name)
	if err != nil {
		return false
	}
	_, err = fs.Stat(FS(), cleaned)
	return err == nil
}

// WalkDir walks embedded Operating Kit resources below root.
func WalkDir(root string, fn fs.WalkDirFunc) error {
	cleaned, err := Clean(root)
	if err != nil {
		return err
	}
	return fs.WalkDir(FS(), cleaned, fn)
}
