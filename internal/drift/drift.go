package drift

import (
	"path/filepath"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/hash"
)

func Report(root string, lock map[string]any) []map[string]string {
	findings := []map[string]string{}
	for _, record := range anySlice(lock["managed_paths"]) {
		mapping, ok := record.(map[string]any)
		if !ok {
			continue
		}
		relative := stringValue(mapping["path"])
		expected := stringValue(mapping["checksum_sha256"])
		path := filepath.Join(root, filepath.FromSlash(relative))
		actual, err := hash.FileSHA256(path)
		if err != nil {
			findings = append(findings, map[string]string{"path": relative, "status": "missing"})
			continue
		}
		if expected != "" && actual != expected {
			findings = append(findings, map[string]string{
				"path":     relative,
				"status":   "drift",
				"expected": expected,
				"actual":   actual,
			})
		}
	}
	return findings
}

func anySlice(value any) []any {
	switch item := value.(type) {
	case []any:
		return item
	default:
		return nil
	}
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	return filepath.ToSlash(filepath.Clean(valueToString(value)))
}

func valueToString(value any) string {
	switch item := value.(type) {
	case string:
		return item
	default:
		return ""
	}
}
