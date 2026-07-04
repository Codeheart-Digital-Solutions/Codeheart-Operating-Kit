package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/components"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/lockfile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/manifest"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/platforms"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

func RunSync(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{
		"--release-manifest": true,
		"--json":             false,
	})
	if err != nil {
		return writeArgError(stderr, "sync", err)
	}
	path := "."
	if len(positionals) > 0 {
		path = positionals[0]
	}
	if len(positionals) > 1 {
		return writeArgError(stderr, "sync", fmt.Errorf("unexpected argument %q", positionals[1]))
	}
	result, err := Sync(path, values["--release-manifest"], time.Now())
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit sync: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, result); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit sync: error: %v\n", err)
			return 1
		}
		return 0
	}
	fmt.Fprintf(stdout, "Synced %d managed files under .codeheart/kit/.\n", len(result["synced_managed_paths"].([]map[string]any)))
	return 0
}

func Sync(path string, releaseManifestPath string, now time.Time) (map[string]any, error) {
	root := expandPath(path)
	var release map[string]any
	var err error
	if releaseManifestPath != "" {
		release, err = validateReleaseManifestFile(releaseManifestPath)
		if err != nil {
			return nil, err
		}
	}
	existingLock, err := lockfile.ReadLock(root)
	if err != nil {
		return nil, err
	}
	profileID := valueString(existingLock["selected_profile"])
	if profileID == "" {
		profileID = "standard"
	}
	managed, err := components.CopyManagedFiles(root, profileID)
	if err != nil {
		return nil, err
	}
	created, err := components.ScaffoldConsumerFiles(root, profileID)
	if err != nil {
		return nil, err
	}
	agentsStatus, err := components.RefreshAgentsManagedBlock(root)
	if err != nil {
		return nil, err
	}
	gitignoreChanged, err := components.EnsureGitignore(root)
	if err != nil {
		return nil, err
	}
	refreshed, err := refreshLock(root, existingLock, managed, release, created, now)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"synced_managed_paths":       managed,
		"created_generated_surfaces": created,
		"agents_status":              agentsStatus,
		"gitignore_changed":          gitignoreChanged,
		"release_manifest":           release != nil,
		"kit_version":                refreshed["kit_version"],
	}, nil
}

func refreshLock(root string, existing map[string]any, managed []map[string]any, releaseManifest map[string]any, created []map[string]any, now time.Time) (map[string]any, error) {
	profileID := valueString(existing["selected_profile"])
	if profileID == "" {
		profileID = "standard"
	}
	profile, err := manifest.LoadProfile(profileID)
	if err != nil {
		return nil, err
	}
	existingUpdate := mapValue(existing["update_check"])
	existingCLI := mapValue(existing["cli_repair"])
	existingNative := mapValue(existing["native_capabilities"])
	now = now.UTC().Truncate(time.Second)
	refreshed := map[string]any{
		"schema_version":      1,
		"kit_version":         version.Version,
		"selected_profile":    profileID,
		"selected_components": stringsAsAny(profile.SelectedComponents),
		"release":             releaseMetadata(existing, releaseManifest),
		"managed_paths":       mapsToAny(managed),
		"generated_surfaces":  mapsToAny(mergeGeneratedSurfaces(existing["generated_surfaces"], created)),
		"cli_repair":          map[string]any{"installed_cli_path": defaultString(existingCLI["installed_cli_path"], "codeheart-operating-kit"), "repair_source_url": defaultString(existingCLI["repair_source_url"], "local-source"), "repair_checksum_sha256": defaultChecksum(existingCLI["repair_checksum_sha256"])},
		"update_check":        map[string]any{"last_update_check_at": defaultString(existingUpdate["last_update_check_at"], lockfile.FormatTime(now)), "next_update_check_due": defaultString(existingUpdate["next_update_check_due"], lockfile.FormatTime(now.Add(7*24*time.Hour))), "latest_seen_version": version.Version, "update_status": "current"},
		"native_capabilities": existingNative,
	}
	if err := lockfile.WriteLock(root, refreshed); err != nil {
		return nil, err
	}
	return refreshed, nil
}

func releaseMetadata(existing map[string]any, releaseManifest map[string]any) map[string]any {
	if fromManifest := releaseAssetFromManifest(releaseManifest); fromManifest != nil {
		return fromManifest
	}
	bundled, err := manifest.LoadReleaseManifest()
	if err == nil {
		if fromManifest := releaseAssetFromManifest(bundled.Raw); fromManifest != nil {
			return fromManifest
		}
	}
	existingRelease := mapValue(existing["release"])
	checksum := valueString(existingRelease["checksum_sha256"])
	if !validSHA256(checksum) {
		checksum = strings.Repeat("0", 64)
	}
	return map[string]any{
		"asset_url":       defaultString(existingRelease["asset_url"], "local-source"),
		"checksum_sha256": checksum,
	}
}

func releaseAssetFromManifest(data map[string]any) map[string]any {
	if data == nil {
		return nil
	}
	assets, ok := data["assets"].([]any)
	if !ok {
		return nil
	}
	candidates := map[string]bool{}
	for _, platform := range platforms.CandidateAssetPlatforms(runtime.GOOS, runtime.GOARCH) {
		candidates[platform] = true
	}
	for _, item := range assets {
		asset, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name := valueString(asset["name"])
		if !strings.HasPrefix(name, "codeheart-operating-kit-"+version.Version) || strings.HasSuffix(name, ".sha256") {
			continue
		}
		if !candidates[valueString(asset["platform"])] || valueString(asset["url"]) == "" || !usableSHA256(asset["sha256"]) {
			continue
		}
		return map[string]any{
			"asset_url":       valueString(asset["url"]),
			"checksum_sha256": valueString(asset["sha256"]),
		}
	}
	return nil
}

func validateReleaseManifestFile(path string) (map[string]any, error) {
	data, err := readMaybeFileURL(path)
	if err != nil {
		return nil, err
	}
	var manifestData map[string]any
	if err := json.Unmarshal(data, &manifestData); err != nil {
		return nil, err
	}
	for _, item := range anyList(manifestData["assets"]) {
		asset, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if !validSHA256(valueString(asset["sha256"])) {
			return nil, fmt.Errorf("invalid asset checksum for %s", defaultString(asset["name"], "<unnamed>"))
		}
	}
	return manifestData, nil
}

func readMaybeFileURL(path string) ([]byte, error) {
	parsed, err := url.Parse(path)
	if err == nil && parsed.Scheme == "file" {
		return os.ReadFile(localFileURLPath(parsed))
	}
	return os.ReadFile(path)
}

func mergeGeneratedSurfaces(existing any, created []map[string]any) []map[string]any {
	merged := []map[string]any{}
	seen := map[string]bool{}
	for _, item := range anyList(existing) {
		mapping, ok := item.(map[string]any)
		if !ok {
			continue
		}
		copied := copyMap(mapping)
		merged = append(merged, copied)
		if path := valueString(copied["path"]); path != "" {
			seen[path] = true
		}
	}
	for _, item := range created {
		path := valueString(item["path"])
		if path == "" || seen[path] {
			continue
		}
		merged = append(merged, copyMap(item))
		seen[path] = true
	}
	return merged
}

func copyMap(value map[string]any) map[string]any {
	result := map[string]any{}
	for key, nested := range value {
		result[key] = nested
	}
	return result
}

func validSHA256(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, char := range value {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}
	return true
}

func usableSHA256(value any) bool {
	text := valueString(value)
	return validSHA256(text) && text != strings.Repeat("0", 64)
}

func defaultChecksum(value any) string {
	text := valueString(value)
	if !validSHA256(text) {
		return strings.Repeat("0", 64)
	}
	return text
}

func defaultString(value any, fallback string) string {
	text := valueString(value)
	if text == "" {
		return fallback
	}
	return text
}

func valueString(value any) string {
	switch item := value.(type) {
	case string:
		return item
	case nil:
		return ""
	default:
		return fmt.Sprint(item)
	}
}

func mapValue(value any) map[string]any {
	if mapping, ok := value.(map[string]any); ok {
		return mapping
	}
	return map[string]any{}
}

func anyList(value any) []any {
	if list, ok := value.([]any); ok {
		return list
	}
	return nil
}
