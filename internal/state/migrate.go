package state

import (
	"fmt"
	"strings"
	"time"
)

const zeroSHA256 = "0000000000000000000000000000000000000000000000000000000000000000"

// NormalizeLegacyV1 applies only the released checksum-placeholder compatibility rules.
func NormalizeLegacyV1(lock map[string]any) (map[string]any, []string) {
	normalized := DeepCopy(lock)
	anomalies := []string{}
	for _, item := range []struct {
		parent string
		field  string
		path   string
	}{
		{parent: "release", field: "checksum_sha256", path: "release.checksum_sha256"},
		{parent: "cli_repair", field: "repair_checksum_sha256", path: "cli_repair.repair_checksum_sha256"},
	} {
		parent := Map(normalized[item.parent])
		if parent == nil {
			continue
		}
		value, exists := parent[item.field]
		if !exists {
			continue
		}
		if legacyPlaceholder(value) {
			parent[item.field] = zeroSHA256
			anomalies = append(anomalies, item.path)
		}
	}
	return normalized, anomalies
}

// MigrateLockV1 creates a validated lock-v2 value without writing it.
func MigrateLockV1(lock map[string]any, now time.Time, transactionID, command string) (map[string]any, []string, error) {
	if AsInt(lock["schema_version"]) != 1 {
		return nil, nil, fmt.Errorf("lock schema version is not 1")
	}
	normalized, anomalies := NormalizeLegacyV1(lock)
	if err := Validate(LockV1Schema, normalized); err != nil {
		return nil, anomalies, err
	}
	release := Map(normalized["release"])
	checksum := AsString(release["checksum_sha256"])
	assetURL := AsString(release["asset_url"])
	status := "verified"
	if assetURL == "" || assetURL == "local-source" {
		status = "local-source"
	}
	if len(anomalies) > 0 || checksum == zeroSHA256 {
		status = "unverified-legacy"
	}
	provenance := map[string]any{
		"verification_status": status,
		"source":              firstNonEmpty(assetURL, "legacy-lock-v1"),
	}
	if status == "verified" {
		provenance["archive_sha256"] = checksum
	}
	if transactionID == "" {
		transactionID = "migration-preview"
	}
	if command == "" {
		command = "repair"
	}
	now = now.UTC().Truncate(time.Second)
	result := map[string]any{
		"schema_version":      2,
		"kit_version":         AsString(normalized["kit_version"]),
		"state_generation":    1,
		"selected_profile":    AsString(normalized["selected_profile"]),
		"selected_components": normalized["selected_components"],
		"release":             release,
		"release_provenance":  provenance,
		"managed_paths":       normalized["managed_paths"],
		"managed_sections":    []any{},
		"generated_surfaces":  normalized["generated_surfaces"],
		"cli_repair":          normalized["cli_repair"],
		"update_check":        normalized["update_check"],
		"native_capabilities": normalized["native_capabilities"],
		"last_operation": map[string]any{
			"transaction_id":      transactionID,
			"command":             command,
			"completed_at":        now.Format(time.RFC3339),
			"previous_generation": 0,
		},
	}
	if err := Validate(LockV2Schema, result); err != nil {
		return nil, anomalies, fmt.Errorf("validate migrated lock v2: %w", err)
	}
	return result, anomalies, nil
}

func legacyPlaceholder(value any) bool {
	switch item := value.(type) {
	case int:
		return item == 0
	case int64:
		return item == 0
	case float64:
		return item == 0
	case string:
		return item == "0" || item == zeroSHA256 || (len(item) == 64 && strings.Trim(item, "0") == "")
	default:
		return false
	}
}

func AsString(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}

func AsInt(value any) int {
	switch item := value.(type) {
	case int:
		return item
	case int64:
		return int(item)
	case uint64:
		return int(item)
	case float64:
		return int(item)
	default:
		return 0
	}
}

func Map(value any) map[string]any {
	if mapping, ok := value.(map[string]any); ok {
		return mapping
	}
	return nil
}

func AnySlice(value any) []any {
	if values, ok := value.([]any); ok {
		return values
	}
	return nil
}
