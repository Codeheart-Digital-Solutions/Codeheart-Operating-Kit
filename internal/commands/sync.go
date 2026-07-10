package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

func RunSync(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{
		"--release-manifest": true,
		"--dry-run":          false,
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
	payload, operation, err := syncOperation(path, values["--release-manifest"], time.Now(), bools["--dry-run"])
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit sync: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, payload); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit sync: error: %v\n", err)
			return 1
		}
	} else if bools["--dry-run"] || !operation.OK() {
		reconcile.WriteText(stdout, operation)
	} else {
		fmt.Fprintf(stdout, "Synced %d managed files under .codeheart/kit/.\n", len(payload["synced_managed_paths"].([]map[string]any)))
	}
	if operation.OK() {
		return 0
	}
	return 1
}

func Sync(path string, releaseManifestPath string, now time.Time) (map[string]any, error) {
	payload, _, err := syncOperation(path, releaseManifestPath, now, false)
	return payload, err
}

func syncOperation(path, releaseManifestPath string, now time.Time, dryRun bool) (map[string]any, reconcile.Result, error) {
	root := expandPath(path)
	releaseManifestProvided := releaseManifestPath != ""
	if releaseManifestProvided {
		if _, err := validateReleaseManifestFile(releaseManifestPath); err != nil {
			return nil, reconcile.Result{}, err
		}
	}
	observed, err := state.Inspect(root)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	if observed.Classification != state.StateCurrent && observed.Classification != state.StateDrifted {
		result := blockedLifecycle("sync", observed, "sync_requires_current_installation", "sync requires a compatible current-version installation", "run check; use repair for compatible drift or upgrade for a version change", "check")
		return syncPayload(result, state.Graph{}, nil, releaseManifestProvided), result, nil
	}
	installedVersion := state.AsString(observed.Lock["kit_version"])
	if installedVersion != version.Version {
		result := blockedLifecycle("sync", observed, "version_change_requires_upgrade", fmt.Sprintf("sync cannot change kit version from %s to %s", installedVersion, version.Version), "use upgrade with explicit version-change approval", "upgrade --dry-run")
		return syncPayload(result, state.Graph{}, observed.Lock, releaseManifestProvided), result, nil
	}
	graph, err := compileObservedGraph(observed)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	lock, err := desiredLifecycleLock("sync", observed, graph, now)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	// Sync refreshes only the running binary's embedded content. A supplied release
	// manifest remains a compatibility validation input and never changes provenance.
	result, err := runLifecycle(lifecycleRequest{
		command:      "sync",
		root:         root,
		now:          now,
		dryRun:       dryRun,
		observed:     observed,
		graph:        graph,
		desiredLock:  lock,
		ensureIgnore: true,
	})
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	return syncPayload(result, graph, lock, releaseManifestProvided), result, nil
}

func syncPayload(result reconcile.Result, graph state.Graph, lock map[string]any, releaseManifest bool) map[string]any {
	payload := resultPayload(result)
	payload["synced_managed_paths"] = recordsAsMaps(reconcile.ManagedPathRecords(graph))
	created := []map[string]any{}
	for _, change := range result.Changes {
		if change.Action == "create" && change.Owner != "managed" {
			created = append(created, map[string]any{"path": change.Path, "ownership": change.Owner})
		}
	}
	payload["created_generated_surfaces"] = created
	payload["agents_status"] = changeStatus(result, "AGENTS.md")
	payload["gitignore_changed"] = changeContains(result, ".gitignore")
	payload["release_manifest"] = releaseManifest
	payload["kit_version"] = lock["kit_version"]
	return payload
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
