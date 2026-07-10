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
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

const defaultLatestReleaseURL = "https://api.github.com/repos/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/latest"

func RunUpdateCheck(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{
		"--latest-version":     true,
		"--metadata-url":       true,
		"--now":                true,
		"--agent-notification": false,
		"--dry-run":            false,
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
	payload, operation, err := updateCheckOperation(path, values["--latest-version"], values["--now"], values["--metadata-url"], bools["--dry-run"])
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit update-check: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, payload); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit update-check: error: %v\n", err)
			return 1
		}
	} else if !operation.OK() || bools["--dry-run"] {
		reconcile.WriteText(stdout, operation)
	} else {
		status := valueString(payload["status"])
		if bools["--agent-notification"] && status == "current" {
			return 0
		}
		switch status {
		case "update-available":
			fmt.Fprintf(stdout, "Operating Kit update available: %s. Apply it only if the user asks.\n", payload["latest_seen_version"])
		case "failed":
			fmt.Fprintf(stdout, "Operating Kit update check failed: %s\n", payload["error"])
		default:
			fmt.Fprintln(stdout, "Operating Kit is current.")
		}
	}
	if operation.OK() {
		return 0
	}
	return 1
}

func UpdateCheck(path string, latestVersion string, nowText string, metadataURL string) (map[string]any, error) {
	payload, _, err := updateCheckOperation(path, latestVersion, nowText, metadataURL, false)
	return payload, err
}

func updateCheckOperation(path, latestVersion, nowText, metadataURL string, dryRun bool) (map[string]any, reconcile.Result, error) {
	root := expandPath(path)
	observed, err := state.Inspect(root)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	if (observed.Classification != state.StateCurrent && observed.Classification != state.StateDrifted) || state.AsInt(observed.Lock["schema_version"]) != 2 {
		result := blockedLifecycle("update-check", observed, "update_check_requires_valid_v2_installation", "update-check requires a valid lock-v2 installation", "run check and repair a compatible legacy or drifted installation first", "repair")
		payload := resultPayload(result)
		payload["status"] = "blocked"
		return payload, result, nil
	}
	now := lockfile.UTCNow()
	if nowText != "" {
		parsed, err := lockfile.ParseTime(nowText)
		if err != nil {
			return nil, reconcile.Result{}, err
		}
		now = parsed
	}
	current := state.AsString(observed.Lock["kit_version"])
	lock := state.DeepCopy(observed.Lock)
	updateState := state.Map(lock["update_check"])
	domainStatus := "current"
	latest := latestVersion
	metadataError := ""
	if latest == "" {
		latest, err = latestVersionFromMetadata(defaultString(metadataURL, defaultLatestReleaseURL))
		if err != nil {
			metadataError = err.Error()
			domainStatus = "failed"
			latest = state.AsString(updateState["latest_seen_version"])
			updateState["update_status"] = "failed"
		} else if compareVersions(latest, current) > 0 {
			domainStatus = "update-available"
		}
	} else if compareVersions(latest, current) > 0 {
		domainStatus = "update-available"
	}
	if metadataError == "" {
		next := now.UTC().Truncate(time.Second).Add(7 * 24 * time.Hour)
		lock["update_check"] = map[string]any{
			"last_update_check_at":  lockfile.FormatTime(now),
			"next_update_check_due": lockfile.FormatTime(next),
			"latest_seen_version":   latest,
			"update_status":         domainStatus,
		}
	}
	lock["last_operation"] = map[string]any{
		"transaction_id":      "planning",
		"command":             "update-check",
		"completed_at":        now.UTC().Truncate(time.Second).Format(time.RFC3339),
		"previous_generation": state.AsInt(lock["state_generation"]),
	}
	result, err := runLifecycle(lifecycleRequest{
		command:       "update-check",
		root:          root,
		now:           now,
		dryRun:        dryRun,
		observed:      observed,
		desiredLock:   lock,
		expectedAfter: []state.Classification{observed.Classification},
	})
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	updateState = state.Map(lock["update_check"])
	payload := resultPayload(result)
	payload["status"] = domainStatus
	payload["latest_seen_version"] = updateState["latest_seen_version"]
	payload["next_update_check_due"] = updateState["next_update_check_due"]
	if metadataError != "" {
		payload["error"] = metadataError
	}
	return payload, result, nil
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
