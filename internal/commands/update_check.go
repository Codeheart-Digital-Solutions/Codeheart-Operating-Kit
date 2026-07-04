package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/lockfile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

const defaultLatestReleaseURL = "https://api.github.com/repos/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/latest"

func RunUpdateCheck(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{
		"--latest-version":     true,
		"--metadata-url":       true,
		"--now":                true,
		"--agent-notification": false,
		"--json":               false,
	})
	if err != nil {
		return writeArgError(stderr, "update-check", err)
	}
	path := "."
	if len(positionals) > 0 {
		path = positionals[0]
	}
	if len(positionals) > 1 {
		return writeArgError(stderr, "update-check", fmt.Errorf("unexpected argument %q", positionals[1]))
	}
	result, err := UpdateCheck(path, values["--latest-version"], values["--now"], values["--metadata-url"])
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit update-check: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, result); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit update-check: error: %v\n", err)
			return 1
		}
		return 0
	}
	status := valueString(result["status"])
	if bools["--agent-notification"] && status == "current" {
		return 0
	}
	switch status {
	case "update-available":
		fmt.Fprintf(stdout, "Operating Kit update available: %s. Apply it only if the user asks.\n", result["latest_seen_version"])
	case "failed":
		fmt.Fprintf(stdout, "Operating Kit update check failed: %s\n", result["error"])
	default:
		fmt.Fprintln(stdout, "Operating Kit is current.")
	}
	return 0
}

func UpdateCheck(path string, latestVersion string, nowText string, metadataURL string) (map[string]any, error) {
	root := expandPath(path)
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		return nil, err
	}
	now := lockfile.UTCNow()
	if nowText != "" {
		parsed, err := lockfile.ParseTime(nowText)
		if err != nil {
			return nil, err
		}
		now = parsed
	}
	current := valueString(lock["kit_version"])
	if current == "" {
		current = version.Version
	}
	ensureLockDefaults(lock, current)

	latest := latestVersion
	if latest == "" {
		latest, err = latestVersionFromMetadata(defaultString(metadataURL, defaultLatestReleaseURL))
		if err != nil {
			updateState := mapValue(lock["update_check"])
			if updateState["last_update_check_at"] == nil {
				updateState["last_update_check_at"] = lockfile.FormatTime(now)
			}
			if updateState["next_update_check_due"] == nil {
				updateState["next_update_check_due"] = lockfile.FormatTime(now)
			}
			if updateState["latest_seen_version"] == nil {
				updateState["latest_seen_version"] = current
			}
			updateState["update_status"] = "failed"
			lock["update_check"] = updateState
			if writeErr := lockfile.WriteLock(root, lock); writeErr != nil {
				return nil, writeErr
			}
			return map[string]any{
				"status":                "failed",
				"latest_seen_version":   updateState["latest_seen_version"],
				"next_update_check_due": updateState["next_update_check_due"],
				"error":                 err.Error(),
			}, nil
		}
	}
	status := "current"
	if compareVersions(latest, current) > 0 {
		status = "update-available"
	}
	next := now.UTC().Truncate(time.Second).Add(7 * 24 * time.Hour)
	lock["update_check"] = map[string]any{
		"last_update_check_at":  lockfile.FormatTime(now),
		"next_update_check_due": lockfile.FormatTime(next),
		"latest_seen_version":   latest,
		"update_status":         status,
	}
	if err := lockfile.WriteLock(root, lock); err != nil {
		return nil, err
	}
	return map[string]any{"status": status, "latest_seen_version": latest, "next_update_check_due": lockfile.FormatTime(next)}, nil
}

func latestVersionFromMetadata(metadataURL string) (string, error) {
	data, err := readMetadataURL(metadataURL)
	if err != nil {
		return "", err
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", err
	}
	for _, key := range []string{"tag_name", "latest_version", "version", "name"} {
		if value := valueString(payload[key]); value != "" {
			return value, nil
		}
	}
	return "", fmt.Errorf("latest-version metadata did not contain tag_name, latest_version, version, or name")
}

func readMetadataURL(metadataURL string) ([]byte, error) {
	parsed, err := url.Parse(metadataURL)
	if err == nil && parsed.Scheme == "file" {
		return os.ReadFile(localFileURLPath(parsed))
	}
	request, err := http.NewRequest(http.MethodGet, metadataURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/vnd.github+json, application/json")
	request.Header.Set("User-Agent", "codeheart-operating-kit/"+version.Version)
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("latest-version metadata returned HTTP %d", response.StatusCode)
	}
	return io.ReadAll(response.Body)
}

func ensureLockDefaults(lock map[string]any, current string) {
	if lock["schema_version"] == nil {
		lock["schema_version"] = 1
	}
	if lock["kit_version"] == nil {
		lock["kit_version"] = current
	}
	if lock["selected_profile"] == nil {
		lock["selected_profile"] = "standard"
	}
	if lock["selected_components"] == nil {
		lock["selected_components"] = []any{}
	}
	if lock["release"] == nil {
		lock["release"] = map[string]any{"asset_url": "unknown", "checksum_sha256": strings.Repeat("0", 64)}
	}
	if lock["managed_paths"] == nil {
		lock["managed_paths"] = []any{}
	}
	if lock["generated_surfaces"] == nil {
		lock["generated_surfaces"] = []any{}
	}
	if lock["cli_repair"] == nil {
		lock["cli_repair"] = map[string]any{"installed_cli_path": "codeheart-operating-kit", "repair_source_url": "unknown", "repair_checksum_sha256": strings.Repeat("0", 64)}
	}
}

func compareVersions(left string, right string) int {
	leftParts := versionParts(left)
	rightParts := versionParts(right)
	maxLen := len(leftParts)
	if len(rightParts) > maxLen {
		maxLen = len(rightParts)
	}
	for len(leftParts) < maxLen {
		leftParts = append(leftParts, 0)
	}
	for len(rightParts) < maxLen {
		rightParts = append(rightParts, 0)
	}
	for index := 0; index < maxLen; index++ {
		if leftParts[index] > rightParts[index] {
			return 1
		}
		if leftParts[index] < rightParts[index] {
			return -1
		}
	}
	return 0
}

func versionParts(value string) []int {
	value = strings.TrimPrefix(value, "v")
	parts := strings.Split(value, ".")
	result := make([]int, len(parts))
	for index, part := range parts {
		parsed, err := strconv.Atoi(part)
		if err != nil {
			parsed = 0
		}
		result[index] = parsed
	}
	return result
}
