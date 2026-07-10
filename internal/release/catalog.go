package release

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type Catalog struct {
	SchemaVersion int            `json:"schema_version"`
	Version       string         `json:"version"`
	GeneratedAt   string         `json:"generated_at,omitempty"`
	Assets        []CatalogAsset `json:"assets"`
}

type CatalogAsset struct {
	Name               string `json:"name"`
	Version            string `json:"version"`
	Platform           string `json:"platform"`
	URL                string `json:"url"`
	ArchiveSHA256      string `json:"archive_sha256"`
	PackManifestSHA256 string `json:"pack_manifest_sha256"`
}

type LoadedCatalog struct {
	Catalog      Catalog
	Location     string
	Base         string
	DigestSHA256 string
	Raw          []byte
}

func LoadCatalog(location string) (LoadedCatalog, error) {
	data, base, err := readLocation(location, "")
	if err != nil {
		return LoadedCatalog{}, err
	}
	var raw map[string]any
	if err := decodeStrictJSON(data, &raw); err != nil {
		return LoadedCatalog{}, fmt.Errorf("decode release catalog: %w", err)
	}
	if err := state.Validate(state.ReleaseCatalogSchema, raw); err != nil {
		return LoadedCatalog{}, err
	}
	var catalog Catalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return LoadedCatalog{}, err
	}
	digest := sha256.Sum256(data)
	return LoadedCatalog{Catalog: catalog, Location: location, Base: base, DigestSHA256: hex.EncodeToString(digest[:]), Raw: data}, nil
}

func (catalog LoadedCatalog) Select(version, platform string) (CatalogAsset, error) {
	if catalog.Catalog.Version != version {
		return CatalogAsset{}, fmt.Errorf("release catalog version %s does not match requested version %s", catalog.Catalog.Version, version)
	}
	matches := []CatalogAsset{}
	for _, asset := range catalog.Catalog.Assets {
		if asset.Version == version && asset.Platform == platform {
			matches = append(matches, asset)
		}
	}
	if len(matches) != 1 {
		return CatalogAsset{}, fmt.Errorf("release catalog has %d assets for version %s platform %s", len(matches), version, platform)
	}
	expectedName := fmt.Sprintf("codeheart-operating-kit-%s-%s.zip", version, platform)
	if matches[0].Name != expectedName || filepath.Base(matches[0].Name) != matches[0].Name {
		return CatalogAsset{}, fmt.Errorf("release catalog asset name %q does not match %q", matches[0].Name, expectedName)
	}
	return matches[0], nil
}

func FetchAsset(catalog LoadedCatalog, asset CatalogAsset, destination string) error {
	data, _, err := readLocation(asset.URL, catalog.Base)
	if err != nil {
		return err
	}
	digest := sha256.Sum256(data)
	actual := hex.EncodeToString(digest[:])
	if !strings.EqualFold(actual, asset.ArchiveSHA256) {
		return fmt.Errorf("archive checksum mismatch: expected %s, got %s", asset.ArchiveSHA256, actual)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return err
	}
	return os.WriteFile(destination, data, 0o600)
}

func readLocation(location, base string) ([]byte, string, error) {
	resolved, err := resolveLocation(location, base)
	if err != nil {
		return nil, "", err
	}
	parsed, err := url.Parse(resolved)
	if err != nil {
		return nil, "", err
	}
	switch parsed.Scheme {
	case "https":
		request, err := http.NewRequest(http.MethodGet, resolved, nil)
		if err != nil {
			return nil, "", err
		}
		request.Header.Set("Accept", "application/json, application/zip, application/octet-stream")
		request.Header.Set("User-Agent", "codeheart-operating-kit")
		client := http.Client{Timeout: 30 * time.Second}
		response, err := client.Do(request)
		if err != nil {
			return nil, "", err
		}
		defer response.Body.Close()
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return nil, "", fmt.Errorf("%s returned HTTP %d", resolved, response.StatusCode)
		}
		data, err := io.ReadAll(io.LimitReader(response.Body, 512<<20))
		return data, directoryLocation(resolved), err
	case "file":
		path, err := url.PathUnescape(parsed.Path)
		if err != nil {
			return nil, "", err
		}
		if runtime.GOOS == "windows" {
			if parsed.Host != "" && parsed.Host != "localhost" {
				path = `\\` + parsed.Host + filepath.FromSlash(path)
			} else if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
				path = path[1:]
			}
		}
		localPath := filepath.FromSlash(path)
		data, err := os.ReadFile(localPath)
		return data, filepath.Dir(localPath), err
	case "":
		absolute, err := filepath.Abs(resolved)
		if err != nil {
			return nil, "", err
		}
		data, err := os.ReadFile(absolute)
		return data, filepath.Dir(absolute), err
	default:
		return nil, "", fmt.Errorf("release locations must use HTTPS, file URLs, or local paths")
	}
}

func resolveLocation(location, base string) (string, error) {
	parsed, err := url.Parse(location)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "" || filepath.IsAbs(location) || base == "" {
		return location, nil
	}
	baseURL, err := url.Parse(base)
	if err == nil && baseURL.Scheme == "https" {
		reference, err := url.Parse(location)
		if err != nil {
			return "", err
		}
		return baseURL.ResolveReference(reference).String(), nil
	}
	return filepath.Join(base, filepath.FromSlash(location)), nil
}

func directoryLocation(location string) string {
	parsed, err := url.Parse(location)
	if err != nil {
		return location
	}
	parsed.Path = strings.TrimSuffix(parsed.Path, filepath.Base(parsed.Path))
	return parsed.String()
}

func decodeStrictJSON(data []byte, target any) error {
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	return nil
}
