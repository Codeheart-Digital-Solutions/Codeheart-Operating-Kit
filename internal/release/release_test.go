package release

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type packFixtureOptions struct {
	version           string
	platform          string
	command           string
	includeBinary     bool
	corruptPayload    bool
	traversal         bool
	symlink           bool
	reconcileSucceeds bool
	nativeBinary      bool
}

func TestMain(m *testing.M) {
	if os.Getenv("CODEHEART_RELEASE_TEST_RECONCILE_FAILURE") == "1" {
		os.Exit(23)
	}
	os.Exit(m.Run())
}

func TestCatalogAndPackVerification(t *testing.T) {
	root := t.TempDir()
	archive, asset := writePackFixture(t, root, packFixtureOptions{version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", includeBinary: true, reconcileSucceeds: true})
	catalogPath := writeCatalogFixture(t, root, asset)
	catalog, err := LoadCatalog(catalogPath)
	if err != nil {
		t.Fatal(err)
	}
	selected, err := catalog.Select("0.1.22", "macos-universal")
	if err != nil {
		t.Fatal(err)
	}
	fetched := filepath.Join(root, "fetched.zip")
	if err := FetchAsset(catalog, selected, fetched); err != nil {
		t.Fatal(err)
	}
	verified, err := VerifyPack(fetched, filepath.Join(root, "extract"), selected, VerifyOptions{Version: "0.1.22", Platform: "macos-universal", Command: "codeheart-operating-kit"})
	if err != nil {
		t.Fatal(err)
	}
	if verified.BinaryPath == "" || verified.ArchiveSHA256 != fileDigest(t, archive) {
		t.Fatalf("verified pack = %#v", verified)
	}
}

func TestPackVerificationRejectsInvalidEvidence(t *testing.T) {
	tests := []struct {
		name     string
		options  packFixtureOptions
		version  string
		platform string
		command  string
		contains string
	}{
		{name: "missing-binary", options: packFixtureOptions{version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit"}, version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", contains: "binary checksum"},
		{name: "payload-checksum", options: packFixtureOptions{version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", includeBinary: true, corruptPayload: true}, version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", contains: "payload checksum"},
		{name: "wrong-version", options: packFixtureOptions{version: "9.9.9", platform: "macos-universal", command: "codeheart-operating-kit", includeBinary: true}, version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", contains: "pack version"},
		{name: "wrong-platform", options: packFixtureOptions{version: "0.1.22", platform: "windows-x64", command: "codeheart-operating-kit", includeBinary: true}, version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", contains: "pack platform"},
		{name: "wrong-command", options: packFixtureOptions{version: "0.1.22", platform: "macos-universal", command: "other", includeBinary: true}, version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", contains: "jsonschema validation"},
		{name: "traversal", options: packFixtureOptions{version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", includeBinary: true, traversal: true}, version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", contains: "unsafe archive path"},
		{name: "symlink", options: packFixtureOptions{version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", includeBinary: true, symlink: true}, version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", contains: "symbolic link"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			archive, asset := writePackFixture(t, root, test.options)
			_, err := VerifyPack(archive, filepath.Join(root, "extract"), asset, VerifyOptions{Version: test.version, Platform: test.platform, Command: test.command})
			if err == nil || !strings.Contains(err.Error(), test.contains) {
				t.Fatalf("error = %v, want %q", err, test.contains)
			}
		})
	}
}

func TestCatalogRejectsChecksumMismatchDuplicateAndUnknownFields(t *testing.T) {
	root := t.TempDir()
	archive, asset := writePackFixture(t, root, packFixtureOptions{version: "0.1.22", platform: "macos-universal", command: "codeheart-operating-kit", includeBinary: true})
	asset.ArchiveSHA256 = strings.Repeat("0", 64)
	catalogPath := writeCatalogFixture(t, root, asset)
	catalog, err := LoadCatalog(catalogPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := FetchAsset(catalog, asset, filepath.Join(root, "copy.zip")); err == nil {
		t.Fatalf("checksum mismatch accepted for %s", archive)
	}
	catalog.Catalog.Assets = append(catalog.Catalog.Assets, asset)
	if _, err := catalog.Select(asset.Version, asset.Platform); err == nil {
		t.Fatalf("duplicate catalog selection accepted")
	}
	catalog.Catalog.Assets = catalog.Catalog.Assets[:1]
	catalog.Catalog.Version = "0.1.20"
	if _, err := catalog.Select(asset.Version, asset.Platform); err == nil {
		t.Fatalf("mismatched catalog version accepted")
	}
	catalog.Catalog.Version = asset.Version
	catalog.Catalog.Assets[0].Name = "../codeheart-operating-kit-0.1.22-macos-universal.zip"
	if _, err := catalog.Select(asset.Version, asset.Platform); err == nil {
		t.Fatalf("unsafe catalog asset name accepted")
	}
	invalid := filepath.Join(root, "invalid.json")
	if err := os.WriteFile(invalid, []byte(`{"schema_version":1,"version":"0.1.22","assets":[],"unexpected":true}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCatalog(invalid); err == nil {
		t.Fatalf("unknown catalog field accepted")
	}
}

func TestCommittedNegativeReleaseFixturesFailClosed(t *testing.T) {
	if _, err := LoadCatalog(filepath.Join("..", "..", "tests", "fixtures", "release", "catalog-invalid-checksum.json")); err == nil {
		t.Fatalf("invalid catalog fixture was accepted")
	}
	data, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "release", "pack-wrong-command.json"))
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if err := state.Validate(state.PackManifestSchema, raw); err == nil {
		t.Fatalf("wrong-command pack fixture was accepted")
	}
}

func TestHandoffTamperingAndFailedReconcileRestorePreviousBinary(t *testing.T) {
	root := t.TempDir()
	platform := "macos-universal"
	targetName := "installed"
	if runtime.GOOS == "windows" {
		platform = "windows-x64"
		targetName += ".exe"
	}
	target := filepath.Join(root, targetName)
	if err := os.WriteFile(target, []byte("previous binary\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	archive, asset := writePackFixture(t, root, packFixtureOptions{version: "0.1.22", platform: platform, command: "codeheart-operating-kit", includeBinary: true, nativeBinary: true})
	pack, err := VerifyPack(archive, filepath.Join(root, "extract"), asset, VerifyOptions{Version: "0.1.22", Platform: platform, Command: "codeheart-operating-kit"})
	if err != nil {
		t.Fatal(err)
	}
	prepared := PreparedUpgrade{Catalog: LoadedCatalog{Location: "catalog.json", DigestSHA256: strings.Repeat("a", 64)}, Asset: asset, Pack: pack}
	handoff, err := NewHandoff(prepared, root, target, "0.1.21")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("CODEHEART_RELEASE_TEST_RECONCILE_FAILURE", "1")
	if err := ApplyHandoff(handoff); err == nil {
		t.Fatalf("failed reconciliation unexpectedly succeeded")
	}
	if data, _ := os.ReadFile(target); string(data) != "previous binary\n" {
		t.Fatalf("previous binary was not restored: %q", data)
	}
	if err := os.WriteFile(handoff.StagedBinary, []byte("tampered\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := validateHandoff(handoff); err == nil {
		t.Fatalf("tampered handoff accepted")
	}
}

func TestVersionDirection(t *testing.T) {
	if err := RequireForwardUpgrade("0.1.21", "0.1.22"); err != nil {
		t.Fatal(err)
	}
	for _, target := range []string{"0.1.21", "0.1.20", "invalid"} {
		if err := RequireForwardUpgrade("0.1.21", target); err == nil {
			t.Fatalf("target %s accepted", target)
		}
	}
}

func writePackFixture(t *testing.T, root string, options packFixtureOptions) (string, CatalogAsset) {
	t.Helper()
	payloadRoot := "codeheart-operating-kit-" + options.version + "-" + options.platform
	files := map[string][]byte{}
	binaryName := "bin/codeheart-operating-kit"
	if options.platform == "windows-x64" {
		binaryName += ".exe"
	}
	if options.includeBinary {
		if options.nativeBinary {
			executable, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			files[binaryName], err = os.ReadFile(executable)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			exitCode := "1"
			if options.reconcileSucceeds {
				exitCode = "0"
			}
			files[binaryName] = []byte("#!/bin/sh\nif [ \"$1\" = \"--version\" ]; then echo \"codeheart-operating-kit " + options.version + "\"; exit 0; fi\nif [ \"$1\" = \"__upgrade-reconcile\" ]; then exit " + exitCode + "; fi\nexit 1\n")
		}
	}
	content, err := os.ReadFile(filepath.Join("..", "..", "manifest.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files["content-manifest.yaml"] = content
	files["README.md"] = []byte("fixture\n")
	checksums := ""
	paths := make([]string, 0, len(files))
	for name := range files {
		paths = append(paths, name)
	}
	sortStrings(paths)
	for _, name := range paths {
		digest := sha256.Sum256(files[name])
		checksums += hex.EncodeToString(digest[:]) + "  " + name + "\n"
	}
	files["checksums.txt"] = []byte(checksums)
	manifest := map[string]any{
		"schema_version": 1, "version": options.version, "platform": options.platform, "command": options.command,
		"binary_path": binaryName, "binary_sha256": digestBytes(files[binaryName]),
		"content_manifest_path": "content-manifest.yaml", "content_manifest_sha256": digestBytes(content),
		"payload_checksums_path": "checksums.txt", "payload_checksums_sha256": digestBytes(files["checksums.txt"]),
	}
	manifestData, _ := json.MarshalIndent(manifest, "", "  ")
	manifestData = append(manifestData, '\n')
	files["pack-manifest.json"] = manifestData
	if options.corruptPayload {
		files["README.md"] = []byte("changed after checksums\n")
	}
	archivePath := filepath.Join(root, "pack.zip")
	archive, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(archive)
	paths = paths[:0]
	for name := range files {
		paths = append(paths, name)
	}
	sortStrings(paths)
	for _, name := range paths {
		entry, err := writer.Create(payloadRoot + "/" + name)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = entry.Write(files[name])
	}
	if options.traversal {
		entry, _ := writer.Create("../escape")
		_, _ = entry.Write([]byte("escape"))
	}
	if options.symlink {
		header := &zip.FileHeader{Name: payloadRoot + "/linked"}
		header.SetMode(os.ModeSymlink | 0o777)
		entry, _ := writer.CreateHeader(header)
		_, _ = entry.Write([]byte("outside"))
	}
	_ = writer.Close()
	_ = archive.Close()
	asset := CatalogAsset{
		Name: fmt.Sprintf("codeheart-operating-kit-%s-%s.zip", options.version, options.platform), Version: options.version, Platform: options.platform, URL: archivePath,
		ArchiveSHA256: fileDigest(t, archivePath), PackManifestSHA256: digestBytes(manifestData),
	}
	return archivePath, asset
}

func writeCatalogFixture(t *testing.T, root string, asset CatalogAsset) string {
	t.Helper()
	data, _ := json.MarshalIndent(Catalog{SchemaVersion: 1, Version: asset.Version, Assets: []CatalogAsset{asset}}, "", "  ")
	path := filepath.Join(root, "catalog.json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func digestBytes(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func fileDigest(t *testing.T, path string) string {
	t.Helper()
	digest, err := fileSHA256(path)
	if err != nil {
		t.Fatal(err)
	}
	return digest
}

func sortStrings(values []string) {
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] < values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
}
