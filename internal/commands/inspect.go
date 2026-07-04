package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

var technicalMarkers = map[string]bool{
	".git":           true,
	"pyproject.toml": true,
	"package.json":   true,
	"Cargo.toml":     true,
	"go.mod":         true,
	"Makefile":       true,
	"AGENTS.md":      true,
	"src":            true,
	"tests":          true,
}

func RunInspect(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--json": false})
	_ = values
	if err != nil {
		return writeArgError(stderr, "inspect", err)
	}
	path := "."
	if len(positionals) > 0 {
		path = positionals[0]
	}
	if len(positionals) > 1 {
		return writeArgError(stderr, "inspect", fmt.Errorf("unexpected argument %q", positionals[1]))
	}
	result := InspectFolder(path)
	if bools["--json"] {
		if err := writeJSON(stdout, result); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit inspect: error: %v\n", err)
			return 1
		}
		return 0
	}
	fmt.Fprintf(stdout, "%s: %s\n", result["mode"], result["reason"])
	return 0
}

func InspectFolder(input string) map[string]any {
	path := expandPath(input)
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return map[string]any{"path": path, "mode": "ambiguous-folder-stop", "reason": "target is not a folder"}
	}
	if os.IsNotExist(err) {
		return map[string]any{"path": path, "mode": "new-folder-setup", "reason": "folder does not exist"}
	}
	entriesList, err := os.ReadDir(path)
	if err != nil {
		return map[string]any{"path": path, "mode": "ambiguous-folder-stop", "reason": err.Error()}
	}
	entries := map[string]bool{}
	for _, entry := range entriesList {
		entries[entry.Name()] = true
	}
	if entries[".codeheart"] && (exists(filepath.Join(path, ".codeheart", "kit.lock.yaml")) || exists(filepath.Join(path, ".codeheart", "kit"))) {
		return map[string]any{"path": path, "mode": "existing-operating-kit-repair", "reason": "Operating Kit metadata exists"}
	}
	if len(entries) == 0 {
		return map[string]any{"path": path, "mode": "new-folder-setup", "reason": "folder is empty"}
	}
	markers := []string{}
	for entry := range entries {
		if technicalMarkers[entry] {
			markers = append(markers, entry)
		}
	}
	if len(markers) > 0 {
		sort.Strings(markers)
		return map[string]any{
			"path":    path,
			"mode":    "existing-technical-project-adoption",
			"reason":  "technical project markers found",
			"markers": stringsAsAny(markers),
		}
	}
	if len(entries) > 100 {
		return map[string]any{"path": path, "mode": "ambiguous-folder-stop", "reason": "folder has many existing files"}
	}
	return map[string]any{"path": path, "mode": "existing-folder-setup", "reason": "folder contains existing files"}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func stringsAsAny(values []string) []any {
	result := make([]any, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}
