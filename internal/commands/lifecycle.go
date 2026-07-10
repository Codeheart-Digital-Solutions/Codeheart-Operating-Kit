package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/capabilities"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

type lifecycleRequest struct {
	command       string
	root          string
	now           time.Time
	dryRun        bool
	observed      state.Observed
	graph         state.Graph
	desiredLock   map[string]any
	desiredConfig map[string]any
	extraFiles    map[string][]byte
	ensureIgnore  bool
	expectedAfter []state.Classification
}

func runLifecycle(request lifecycleRequest) (reconcile.Result, error) {
	plan, err := reconcile.BuildPlan(reconcile.Request{
		Command:       request.command,
		Root:          request.root,
		Observed:      request.observed,
		Graph:         request.graph,
		DesiredLock:   request.desiredLock,
		DesiredConfig: request.desiredConfig,
		ExtraFiles:    request.extraFiles,
		EnsureIgnore:  request.ensureIgnore,
		ExpectedAfter: request.expectedAfter,
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	if request.dryRun {
		return reconcile.Preview(plan), nil
	}
	return reconcile.Apply(plan, reconcile.ApplyOptions{Now: request.now})
}

func blockedLifecycle(command string, observed state.Observed, code, message, remediation, retry string) reconcile.Result {
	result := reconcile.NewResult(command)
	result.Status = reconcile.StatusBlocked
	result.StateBefore = string(observed.Classification)
	result.Blockers = append(result.Blockers, reconcile.Blocker{
		Code:         code,
		Message:      message,
		Remediation:  remediation,
		RetryCommand: retry,
	})
	return result
}

func resultPayload(result reconcile.Result) map[string]any {
	return map[string]any{"result": result.ToMap()}
}

func compileObservedGraph(observed state.Observed) (state.Graph, error) {
	if observed.Graph.ProfileID != "" {
		return observed.Graph, nil
	}
	profileID := state.AsString(observed.Lock["selected_profile"])
	if profileID == "" {
		profileID = "standard"
	}
	return state.CompileGraph(profileID)
}

func initialConfig(projectName, purpose, selectedFolder string) map[string]any {
	config := map[string]any{
		"schema_version":        1,
		"selected_profile":      "standard",
		"project_display_name":  projectName,
		"selected_setup_folder": selectedFolder,
		"local_consumer_layer": map[string]any{
			"repo_docs_path":           "docs/repo/",
			"agent_memory_path":        "docs/agent-memory/",
			"user_layer_path":          ".codeheart/user/",
			"local_machine_layer_path": ".codeheart/local/",
		},
		"component_settings": map[string]any{},
	}
	if purpose != "" {
		config["setup_purpose"] = purpose
	}
	return config
}

func desiredLifecycleLock(command string, observed state.Observed, graph state.Graph, now time.Time) (map[string]any, error) {
	now = now.UTC().Truncate(time.Second)
	var lock map[string]any
	if observed.Lock == nil || len(observed.Lock) == 0 {
		lock = map[string]any{
			"schema_version":   2,
			"kit_version":      version.Version,
			"state_generation": 1,
			"release": map[string]any{
				"asset_url":       "local-source",
				"checksum_sha256": strings.Repeat("0", 64),
			},
			"release_provenance": map[string]any{
				"verification_status": "local-source",
				"source":              "embedded-running-binary",
			},
			"cli_repair": map[string]any{
				"installed_cli_path":     "codeheart-operating-kit",
				"repair_source_url":      "local-source",
				"repair_checksum_sha256": strings.Repeat("0", 64),
			},
			"update_check": map[string]any{
				"last_update_check_at":  now.Format(time.RFC3339),
				"next_update_check_due": now.Add(7 * 24 * time.Hour).Format(time.RFC3339),
				"latest_seen_version":   version.Version,
				"update_status":         "current",
			},
			"native_capabilities": capabilities.UnknownNativeCapabilityState(now),
		}
	} else if state.AsInt(observed.Lock["schema_version"]) == 1 {
		migrated, _, err := state.MigrateLockV1(observed.Lock, now, "migration-preview", command)
		if err != nil {
			return nil, err
		}
		lock = migrated
	} else if state.AsInt(observed.Lock["schema_version"]) == 2 {
		lock = state.DeepCopy(observed.Lock)
	} else {
		return nil, fmt.Errorf("unsupported lock schema version %d", state.AsInt(observed.Lock["schema_version"]))
	}

	lock["schema_version"] = 2
	lock["selected_profile"] = graph.ProfileID
	lock["selected_components"] = stringsAsAny(graph.SelectedComponents)
	lock["managed_paths"] = reconcile.ManagedPathRecords(graph)
	lock["managed_sections"] = reconcile.ManagedSectionRecords(graph)
	lock["generated_surfaces"] = reconcile.GeneratedSurfaceRecords(graph)
	lock["last_operation"] = map[string]any{
		"transaction_id":      "planning",
		"command":             command,
		"completed_at":        now.Format(time.RFC3339),
		"previous_generation": 0,
	}
	return lock, nil
}

func repairConfig(root string, observed state.Observed) map[string]any {
	if observed.Config != nil {
		return nil
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.ConfigPath))); err == nil {
		return nil
	}
	name := filepath.Base(root)
	if name == "." || name == string(filepath.Separator) || name == "" {
		name = "Operating-Kit-Project"
	}
	return initialConfig(name, "", root)
}

func adoptionReport(now time.Time, findings []string) []byte {
	lines := []string{
		"Last updated: " + now.UTC().Truncate(time.Second).Format(time.RFC3339) + " (UTC)",
		"",
		"# Adoption Cleanup Report",
		"",
		"Existing project guidance was preserved. Review these overlapping surfaces before cleanup:",
		"",
	}
	for _, finding := range findings {
		lines = append(lines, "- "+finding)
	}
	return []byte(strings.Join(lines, "\n") + "\n")
}

func changeContains(result reconcile.Result, path string) bool {
	for _, change := range result.Changes {
		if change.Path == path {
			return true
		}
	}
	return false
}
