package lockfile

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/yamlmini"
)

func TestWriteAndReadRepresentativeLockAndConfig(t *testing.T) {
	root := t.TempDir()
	lock := loadFixtureMap(t, "tests/fixtures/kit-lock.yaml")
	if err := WriteLock(root, lock); err != nil {
		t.Fatalf("WriteLock: %v", err)
	}
	readLock, err := ReadLock(root)
	if err != nil {
		t.Fatalf("ReadLock: %v", err)
	}
	if !reflect.DeepEqual(lock, readLock) {
		t.Fatalf("lock changed after round trip\nbefore: %#v\nafter: %#v", lock, readLock)
	}
	if missing := MissingRequiredLockMetadata(readLock); len(missing) != 0 {
		t.Fatalf("complete lock reported missing metadata: %#v", missing)
	}

	config := loadFixtureMap(t, "tests/fixtures/kit-config.yaml")
	if err := WriteConfig(root, config); err != nil {
		t.Fatalf("WriteConfig: %v", err)
	}
	readConfig, err := ReadConfig(root)
	if err != nil {
		t.Fatalf("ReadConfig: %v", err)
	}
	if !reflect.DeepEqual(config, readConfig) {
		t.Fatalf("config changed after round trip\nbefore: %#v\nafter: %#v", config, readConfig)
	}
}

func TestWriteAndReadConfigV2UsesDeclaredSchemaAndFailsClosed(t *testing.T) {
	root := t.TempDir()
	config := loadFixtureMap(t, "tests/fixtures/kit-config.yaml")
	config["schema_version"] = 2
	config["component_settings"] = map[string]any{"planning-workflows": map[string]any{
		"plan_catalog_mode":              "mixed",
		"plan_catalog_discovery_version": 2,
		"plan_catalog_cutover_revision":  strings.Repeat("a", 40),
		"plan_catalog_migration_evidence": map[string]any{
			"ledger_path":              "docs/repo/plans/example/attachments/ledger-v3.yaml",
			"ledger_sha256":            strings.Repeat("b", 64),
			"branch_evidence_digest":   strings.Repeat("c", 64),
			"evidence_revision":        strings.Repeat("d", 40),
			"activation_base_revision": strings.Repeat("e", 40),
			"migration_action_digest":  strings.Repeat("f", 64),
			"evidence_scope":           "local",
		},
	}}
	if err := WriteConfig(root, config); err != nil {
		t.Fatalf("WriteConfig v2: %v", err)
	}
	read, err := ReadConfig(root)
	if err != nil {
		t.Fatalf("ReadConfig v2: %v", err)
	}
	if !reflect.DeepEqual(config, read) {
		t.Fatalf("config v2 changed after round trip\nbefore: %#v\nafter: %#v", config, read)
	}
	config["schema_version"] = 99
	if err := WriteConfig(root, config); err == nil {
		t.Fatal("future config unexpectedly written")
	}
}

func TestMissingRequiredLockMetadata(t *testing.T) {
	lock := loadFixtureMap(t, "tests/fixtures/kit-lock.yaml")
	delete(lock["release"].(map[string]any), "asset_url")
	delete(lock, "native_capabilities")
	missing := MissingRequiredLockMetadata(lock)
	if !contains(missing, "release.asset_url") || !contains(missing, "native_capabilities") {
		t.Fatalf("missing metadata = %#v", missing)
	}
}

func TestTimeFormattingUsesUTCSecondPrecision(t *testing.T) {
	value := time.Date(2026, 7, 4, 20, 43, 3, 900, time.FixedZone("CEST", 2*60*60))
	formatted := FormatTime(value)
	if formatted != "2026-07-04T18:43:03Z" {
		t.Fatalf("FormatTime = %q", formatted)
	}
	parsed, err := ParseTime(formatted)
	if err != nil {
		t.Fatalf("ParseTime: %v", err)
	}
	if parsed.UTC().Format(time.RFC3339) != formatted {
		t.Fatalf("ParseTime round trip = %s", parsed.UTC().Format(time.RFC3339))
	}
}

func loadFixtureMap(t *testing.T, relative string) map[string]any {
	t.Helper()
	path := filepath.Join("..", "..", relative)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", relative, err)
	}
	parsed, err := yamlmini.MustMap(string(data))
	if err != nil {
		t.Fatalf("parse fixture %s: %v", relative, err)
	}
	return parsed
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
