package lockfile

import (
	"os"
	"path/filepath"
	"reflect"
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
