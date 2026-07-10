package lockfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

const (
	LockPath   = ".codeheart/kit.lock.yaml"
	ConfigPath = ".codeheart/kit.config.yaml"
)

func UTCNow() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func ParseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, strings.Replace(value, "Z", "+00:00", 1))
}

func FormatTime(value time.Time) string {
	return value.UTC().Truncate(time.Second).Format(time.RFC3339)
}

func ReadLock(root string) (map[string]any, error) {
	value, err := readYAML(filepath.Join(root, filepath.FromSlash(LockPath)))
	if err != nil || len(value) == 0 {
		return value, err
	}
	version := state.AsInt(value["schema_version"])
	schemaPath, err := state.SchemaForLockVersion(version)
	if err != nil {
		return nil, err
	}
	if version == 1 {
		value, _ = state.NormalizeLegacyV1(value)
	}
	if err := state.Validate(schemaPath, value); err != nil {
		return nil, err
	}
	return value, nil
}

func WriteLock(root string, lock map[string]any) error {
	schemaPath, err := state.SchemaForLockVersion(state.AsInt(lock["schema_version"]))
	if err != nil {
		return err
	}
	if err := state.Validate(schemaPath, lock); err != nil {
		return err
	}
	return writeYAML(filepath.Join(root, filepath.FromSlash(LockPath)), lock)
}

func ReadConfig(root string) (map[string]any, error) {
	value, err := readYAML(filepath.Join(root, filepath.FromSlash(ConfigPath)))
	if err != nil || len(value) == 0 {
		return value, err
	}
	if err := state.Validate(state.ConfigV1Schema, value); err != nil {
		return nil, err
	}
	return value, nil
}

func WriteConfig(root string, config map[string]any) error {
	if err := state.Validate(state.ConfigV1Schema, config); err != nil {
		return err
	}
	return writeYAML(filepath.Join(root, filepath.FromSlash(ConfigPath)), config)
}

func RequiredLockKeys() map[string]bool {
	keys := map[string]bool{}
	for _, key := range []string{
		"schema_version",
		"kit_version",
		"selected_profile",
		"selected_components",
		"release",
		"managed_paths",
		"generated_surfaces",
		"cli_repair",
		"update_check",
		"native_capabilities",
	} {
		keys[key] = true
	}
	return keys
}

func RequiredLockMetadataPaths() []string {
	return []string{
		"schema_version",
		"kit_version",
		"selected_profile",
		"selected_components",
		"release",
		"release.asset_url",
		"release.checksum_sha256",
		"managed_paths",
		"generated_surfaces",
		"cli_repair",
		"cli_repair.installed_cli_path",
		"cli_repair.repair_source_url",
		"update_check",
		"update_check.last_update_check_at",
		"update_check.next_update_check_due",
		"update_check.latest_seen_version",
		"update_check.update_status",
		"native_capabilities",
	}
}

func MissingRequiredLockMetadata(lock map[string]any) []string {
	missing := []string{}
	for _, metadataPath := range RequiredLockMetadataPaths() {
		current := any(lock)
		found := true
		for _, part := range strings.Split(metadataPath, ".") {
			mapping, ok := current.(map[string]any)
			if !ok {
				found = false
				break
			}
			current, ok = mapping[part]
			if !ok {
				found = false
				break
			}
		}
		if !found {
			missing = append(missing, metadataPath)
		}
	}
	return missing
}

func readYAML(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]any{}, nil
	}
	if err != nil {
		return nil, err
	}
	return state.DecodeYAMLMap(data)
}

func writeYAML(path string, value map[string]any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := state.EncodeYAML(value)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	return os.WriteFile(path, data, 0o644)
}
