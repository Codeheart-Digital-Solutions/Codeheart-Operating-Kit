package release

import (
	"archive/zip"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

const (
	maxPackBytes   = 256 << 20
	maxFileBytes   = 128 << 20
	maxPackEntries = 4096
)

type PackManifest struct {
	SchemaVersion          int    `json:"schema_version"`
	Version                string `json:"version"`
	Platform               string `json:"platform"`
	Command                string `json:"command"`
	BinaryPath             string `json:"binary_path"`
	BinarySHA256           string `json:"binary_sha256"`
	ContentManifestPath    string `json:"content_manifest_path"`
	ContentManifestSHA256  string `json:"content_manifest_sha256"`
	PayloadChecksumsPath   string `json:"payload_checksums_path"`
	PayloadChecksumsSHA256 string `json:"payload_checksums_sha256"`
}

type VerifyOptions struct {
	Version   string
	Platform  string
	Command   string
	SmokeTest bool
}

type VerifiedPack struct {
	Asset                 CatalogAsset
	Manifest              PackManifest
	ArchivePath           string
	PayloadRoot           string
	BinaryPath            string
	ArchiveSHA256         string
	PackManifestSHA256    string
	ContentManifestSHA256 string
}

func VerifyPack(archivePath, destination string, asset CatalogAsset, options VerifyOptions) (VerifiedPack, error) {
	if options.Command == "" {
		options.Command = "codeheart-operating-kit"
	}
	archiveDigest, err := fileSHA256(archivePath)
	if err != nil {
		return VerifiedPack{}, err
	}
	if !strings.EqualFold(archiveDigest, asset.ArchiveSHA256) {
		return VerifiedPack{}, fmt.Errorf("archive checksum mismatch")
	}
	payloadRoot, err := extractPack(archivePath, destination)
	if err != nil {
		return VerifiedPack{}, err
	}
	manifestPath := filepath.Join(payloadRoot, "pack-manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return VerifiedPack{}, fmt.Errorf("read pack manifest: %w", err)
	}
	manifestDigest := sha256.Sum256(manifestData)
	manifestSHA := hex.EncodeToString(manifestDigest[:])
	if !strings.EqualFold(manifestSHA, asset.PackManifestSHA256) {
		return VerifiedPack{}, fmt.Errorf("pack manifest checksum mismatch")
	}
	var raw map[string]any
	if err := decodeStrictJSON(manifestData, &raw); err != nil {
		return VerifiedPack{}, err
	}
	if err := state.Validate(state.PackManifestSchema, raw); err != nil {
		return VerifiedPack{}, err
	}
	var manifest PackManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return VerifiedPack{}, err
	}
	if manifest.Version != options.Version || manifest.Version != asset.Version {
		return VerifiedPack{}, fmt.Errorf("pack version %s does not match requested version %s", manifest.Version, options.Version)
	}
	if manifest.Platform != options.Platform || manifest.Platform != asset.Platform {
		return VerifiedPack{}, fmt.Errorf("pack platform %s does not match requested platform %s", manifest.Platform, options.Platform)
	}
	if manifest.Command != options.Command {
		return VerifiedPack{}, fmt.Errorf("pack command %s does not match %s", manifest.Command, options.Command)
	}
	checksumsPath := filepath.Join(payloadRoot, filepath.FromSlash(manifest.PayloadChecksumsPath))
	checksumsSHA, err := fileSHA256(checksumsPath)
	if err != nil || !strings.EqualFold(checksumsSHA, manifest.PayloadChecksumsSHA256) {
		return VerifiedPack{}, fmt.Errorf("payload checksum identity mismatch")
	}
	if err := verifyPayloadChecksums(payloadRoot, checksumsPath); err != nil {
		return VerifiedPack{}, err
	}
	binaryPath := filepath.Join(payloadRoot, filepath.FromSlash(manifest.BinaryPath))
	binarySHA, err := fileSHA256(binaryPath)
	if err != nil || !strings.EqualFold(binarySHA, manifest.BinarySHA256) {
		return VerifiedPack{}, fmt.Errorf("binary checksum mismatch")
	}
	contentPath := filepath.Join(payloadRoot, filepath.FromSlash(manifest.ContentManifestPath))
	contentSHA, err := fileSHA256(contentPath)
	if err != nil || !strings.EqualFold(contentSHA, manifest.ContentManifestSHA256) {
		return VerifiedPack{}, fmt.Errorf("content manifest checksum mismatch")
	}
	contentData, err := os.ReadFile(contentPath)
	if err != nil {
		return VerifiedPack{}, err
	}
	content, err := state.DecodeAndValidateYAML(state.ContentManifestSchema, contentData)
	if err != nil || state.AsString(content["version"]) != options.Version {
		return VerifiedPack{}, fmt.Errorf("content identity does not match pack version")
	}
	if options.SmokeTest {
		if runtime.GOOS != "windows" {
			_ = os.Chmod(binaryPath, 0o755)
		}
		output, err := exec.Command(binaryPath, "--version").CombinedOutput()
		if err != nil {
			return VerifiedPack{}, fmt.Errorf("staged binary smoke test failed: %w", err)
		}
		expected := options.Command + " " + options.Version
		if strings.TrimSpace(string(output)) != expected {
			return VerifiedPack{}, fmt.Errorf("staged binary reported %q, expected %q", strings.TrimSpace(string(output)), expected)
		}
	}
	return VerifiedPack{
		Asset: asset, Manifest: manifest, ArchivePath: archivePath, PayloadRoot: payloadRoot,
		BinaryPath: binaryPath, ArchiveSHA256: archiveDigest, PackManifestSHA256: manifestSHA,
		ContentManifestSHA256: contentSHA,
	}, nil
}

func extractPack(archivePath, destination string) (string, error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", fmt.Errorf("open release pack: %w", err)
	}
	defer reader.Close()
	if err := os.MkdirAll(destination, 0o700); err != nil {
		return "", err
	}
	roots := map[string]bool{}
	var declaredTotal uint64
	var actualTotal int64
	for index, file := range reader.File {
		if index >= maxPackEntries {
			return "", fmt.Errorf("archive exceeds entry limit")
		}
		name := file.Name
		clean := path.Clean(name)
		if name == "" || strings.Contains(name, "\\") || strings.HasPrefix(name, "/") || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
			return "", fmt.Errorf("unsafe archive path %q", name)
		}
		parts := strings.Split(clean, "/")
		if len(parts) < 2 {
			return "", fmt.Errorf("archive entry has no payload root: %s", name)
		}
		roots[parts[0]] = true
		if len(roots) != 1 {
			return "", fmt.Errorf("archive has multiple payload roots")
		}
		if file.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("archive contains symbolic link %s", name)
		}
		if file.Mode()&os.ModeType != 0 && !file.FileInfo().IsDir() {
			return "", fmt.Errorf("archive contains unsupported filesystem entry %s", name)
		}
		if file.UncompressedSize64 > maxFileBytes || declaredTotal+file.UncompressedSize64 > maxPackBytes {
			return "", fmt.Errorf("archive exceeds extraction limits")
		}
		declaredTotal += file.UncompressedSize64
		target := filepath.Join(destination, filepath.FromSlash(clean))
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return "", err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return "", err
		}
		source, err := file.Open()
		if err != nil {
			return "", err
		}
		mode := os.FileMode(0o644)
		if file.Mode()&0o111 != 0 {
			mode = 0o755
		}
		targetFile, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
		if err == nil {
			var written int64
			written, err = io.Copy(targetFile, io.LimitReader(source, maxFileBytes+1))
			if err == nil && (written > int64(maxFileBytes) || actualTotal+written > int64(maxPackBytes)) {
				err = fmt.Errorf("archive exceeds actual extraction limits")
			}
			actualTotal += written
			closeErr := targetFile.Close()
			if err == nil {
				err = closeErr
			}
		}
		_ = source.Close()
		if err != nil {
			return "", err
		}
	}
	if len(roots) != 1 {
		return "", fmt.Errorf("archive has no payload root")
	}
	for root := range roots {
		return filepath.Join(destination, filepath.FromSlash(root)), nil
	}
	panic("unreachable")
}

func verifyPayloadChecksums(root, checksumPath string) error {
	data, err := os.ReadFile(checksumPath)
	if err != nil {
		return err
	}
	expected := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		checksum, relative, ok := strings.Cut(line, "  ")
		clean := path.Clean(relative)
		if !ok || len(checksum) != 64 || clean != relative || clean == "." || strings.HasPrefix(clean, "../") || strings.Contains(relative, "\\") {
			return fmt.Errorf("invalid payload checksum line %q", line)
		}
		if expected[relative] != "" {
			return fmt.Errorf("duplicate payload checksum %s", relative)
		}
		expected[relative] = checksum
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	actualPaths := []string{}
	err = filepath.WalkDir(root, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		relative, err := filepath.Rel(root, current)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if relative == "pack-manifest.json" || relative == "checksums.txt" {
			return nil
		}
		actualPaths = append(actualPaths, relative)
		want := expected[relative]
		if want == "" {
			return fmt.Errorf("payload file %s is not checksummed", relative)
		}
		actual, err := fileSHA256(current)
		if err != nil || !strings.EqualFold(actual, want) {
			return fmt.Errorf("payload checksum mismatch for %s", relative)
		}
		return nil
	})
	if err != nil {
		return err
	}
	sort.Strings(actualPaths)
	if len(actualPaths) != len(expected) {
		return fmt.Errorf("payload checksum set does not match extracted files")
	}
	return nil
}

func fileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}
