package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/capabilities"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/hash"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/lockfile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/manifest"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

const (
	BeginMarker = "<!-- BEGIN CODEHEART OPERATING KIT MANAGED BLOCK -->"
	EndMarker   = "<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->"
)

var localGitignoreLines = []string{
	"# Codeheart Operating Kit local user layer",
	".codeheart/user/preferences.yaml",
	".codeheart/user/*.local.yaml",
	".codeheart/user/feedback/",
	"# Codeheart Operating Kit local machine layer",
	".codeheart/local/",
}

type DefaultState struct {
	ManagedPaths      []map[string]any
	GeneratedSurfaces []map[string]any
	AgentsStatus      string
	GitignoreChanged  bool
}

func CopyManagedFiles(root string, profileID string) ([]map[string]any, error) {
	records := []map[string]any{}
	files, err := manifest.IterComponentFiles(profileID)
	if err != nil {
		return nil, err
	}
	for _, entry := range files {
		if entry.Ownership != "managed" {
			continue
		}
		targetPath := filepath.Join(root, filepath.FromSlash(entry.Target))
		if err := copyResource(entry.Source, targetPath); err != nil {
			return nil, err
		}
		checksum, err := hash.FileSHA256(targetPath)
		if err != nil {
			return nil, err
		}
		records = append(records, map[string]any{
			"path":            entry.Target,
			"component":       entry.Component,
			"source":          entry.Source,
			"ownership":       "managed",
			"checksum_sha256": checksum,
		})
	}
	return records, nil
}

func RenderAgents(root string) (string, error) {
	template, err := kitfs.ReadText("templates/agents/AGENTS.managed-block.md")
	if err != nil {
		return "", err
	}
	path := filepath.Join(root, "AGENTS.md")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "created", os.WriteFile(path, []byte(template), 0o644)
	}
	existing, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	text := string(existing)
	if strings.Contains(text, BeginMarker) && strings.Contains(text, EndMarker) {
		refreshed, err := replaceManagedBlock(text, template)
		if err != nil {
			return "", err
		}
		return "repaired-managed-block", os.WriteFile(path, []byte(refreshed), 0o644)
	}
	body := template + "\n\n# Existing Instructions Preserved\n\n" + text
	return "added-managed-block", os.WriteFile(path, []byte(body), 0o644)
}

func RefreshAgentsManagedBlock(root string) (string, error) {
	template, err := kitfs.ReadText("templates/agents/AGENTS.managed-block.md")
	if err != nil {
		return "", err
	}
	path := filepath.Join(root, "AGENTS.md")
	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "missing", nil
	}
	if err != nil {
		return "", err
	}
	text := string(existing)
	if !strings.Contains(text, BeginMarker) || !strings.Contains(text, EndMarker) {
		return "unchanged-no-managed-block", nil
	}
	refreshed, err := replaceManagedBlock(text, template)
	if err != nil {
		return "", err
	}
	return "refreshed-managed-block", os.WriteFile(path, []byte(refreshed), 0o644)
}

func ScaffoldConsumerFiles(root string, profileID string) ([]map[string]any, error) {
	created := []map[string]any{}
	files, err := manifest.IterComponentFiles(profileID)
	if err != nil {
		return nil, err
	}
	for _, entry := range files {
		if entry.Ownership != "scaffold" {
			continue
		}
		targetPath := filepath.Join(root, filepath.FromSlash(entry.Target))
		if exists(targetPath) {
			continue
		}
		if err := copyResource(entry.Source, targetPath); err != nil {
			return nil, err
		}
		created = append(created, map[string]any{"path": entry.Target, "ownership": "scaffold"})
	}

	userDir := filepath.Join(root, ".codeheart", "user")
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		return nil, err
	}
	userReadme := filepath.Join(userDir, "README.md")
	if !exists(userReadme) {
		if err := copyResource("templates/user-layer/README.md", userReadme); err != nil {
			return nil, err
		}
		created = append(created, map[string]any{"path": ".codeheart/user/README.md", "ownership": "local-user"})
	}
	examplesDir := filepath.Join(userDir, "examples")
	if err := os.MkdirAll(examplesDir, 0o755); err != nil {
		return nil, err
	}
	examplePref := filepath.Join(examplesDir, "preferences.yaml")
	if !exists(examplePref) {
		if err := copyResource("templates/user-layer/example.preferences.yaml", examplePref); err != nil {
			return nil, err
		}
		created = append(created, map[string]any{"path": ".codeheart/user/examples/preferences.yaml", "ownership": "local-user"})
	}
	return created, nil
}

func EnsureGitignore(root string) (bool, error) {
	path := filepath.Join(root, ".gitignore")
	var lines []string
	if data, err := os.ReadFile(path); err == nil {
		lines = strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
	} else if !os.IsNotExist(err) {
		return false, err
	}
	present := map[string]bool{}
	for _, line := range lines {
		present[line] = true
	}
	missing := []string{}
	for _, line := range localGitignoreLines {
		if !present[line] {
			missing = append(missing, line)
		}
	}
	if len(missing) == 0 {
		return false, nil
	}
	additions := []string{}
	if len(lines) > 0 {
		additions = append(additions, "")
	}
	additions = append(additions, missing...)
	body := strings.Join(append(lines, additions...), "\n")
	body = strings.TrimLeft(body, "\n") + "\n"
	return true, os.WriteFile(path, []byte(body), 0o644)
}

func WriteDefaultState(root string, projectName string, purpose string, selectedFolder string, now time.Time) (DefaultState, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return DefaultState{}, err
	}
	profile, err := manifest.LoadProfile("standard")
	if err != nil {
		return DefaultState{}, err
	}
	managedRecords, err := CopyManagedFiles(root, "standard")
	if err != nil {
		return DefaultState{}, err
	}
	scaffoldRecords, err := ScaffoldConsumerFiles(root, "standard")
	if err != nil {
		return DefaultState{}, err
	}
	agentsStatus, err := RenderAgents(root)
	if err != nil {
		return DefaultState{}, err
	}
	gitignoreChanged, err := EnsureGitignore(root)
	if err != nil {
		return DefaultState{}, err
	}

	generatedSurfaces := []map[string]any{
		{"path": ".codeheart/kit/", "ownership": "managed"},
		{"path": ".codeheart/kit.lock.yaml", "ownership": "generated-surface"},
		{"path": ".codeheart/kit.config.yaml", "ownership": "generated-surface"},
		{"path": "AGENTS.md", "ownership": "template"},
	}
	generatedSurfaces = append(generatedSurfaces, scaffoldRecords...)

	now = now.UTC().Truncate(time.Second)
	nextCheck := now.Add(7 * 24 * time.Hour)
	lock := map[string]any{
		"schema_version":      1,
		"kit_version":         version.Version,
		"selected_profile":    "standard",
		"selected_components": stringSliceAsAny(profile.SelectedComponents),
		"release":             map[string]any{"asset_url": "local-source", "checksum_sha256": strings.Repeat("0", 64)},
		"managed_paths":       mapsAsAny(managedRecords),
		"generated_surfaces":  mapsAsAny(generatedSurfaces),
		"cli_repair":          map[string]any{"installed_cli_path": "codeheart-operating-kit", "repair_source_url": "local-source", "repair_checksum_sha256": strings.Repeat("0", 64)},
		"update_check":        map[string]any{"last_update_check_at": lockfile.FormatTime(now), "next_update_check_due": lockfile.FormatTime(nextCheck), "latest_seen_version": version.Version, "update_status": "current"},
		"native_capabilities": capabilities.UnknownNativeCapabilityState(now),
	}
	if err := lockfile.WriteLock(root, lock); err != nil {
		return DefaultState{}, err
	}

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
	if err := lockfile.WriteConfig(root, config); err != nil {
		return DefaultState{}, err
	}

	return DefaultState{
		ManagedPaths:      managedRecords,
		GeneratedSurfaces: generatedSurfaces,
		AgentsStatus:      agentsStatus,
		GitignoreChanged:  gitignoreChanged,
	}, nil
}

func WriteAdoptionCleanupReport(root string, findings []string) (string, error) {
	path := filepath.Join(root, ".codeheart", "reports", "adoption-cleanup-report.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	lines := []string{
		"Last updated: 2026-06-13T00:00:00Z (UTC)",
		"",
		"# Adoption Cleanup Report",
		"",
		"Existing project guidance was preserved. Review these overlapping surfaces before cleanup:",
		"",
	}
	for _, finding := range findings {
		lines = append(lines, "- "+finding)
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func copyResource(source string, targetPath string) error {
	data, err := kitfs.ReadFile(source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(targetPath, data, 0o644)
}

func replaceManagedBlock(existing string, template string) (string, error) {
	before, rest, ok := strings.Cut(existing, BeginMarker)
	if !ok {
		return "", fmt.Errorf("existing text has no begin marker")
	}
	_, after, ok := strings.Cut(rest, EndMarker)
	if !ok {
		return "", fmt.Errorf("existing text has no end marker")
	}
	_, managedRest, ok := strings.Cut(template, BeginMarker)
	if !ok {
		return "", fmt.Errorf("template has no begin marker")
	}
	managed, _, ok := strings.Cut(managedRest, EndMarker)
	if !ok {
		return "", fmt.Errorf("template has no end marker")
	}
	return before + BeginMarker + managed + EndMarker + after, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func mapsAsAny(values []map[string]any) []any {
	result := make([]any, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}

func stringSliceAsAny(values []string) []any {
	result := make([]any, len(values))
	for index, value := range values {
		result[index] = value
	}
	return result
}
