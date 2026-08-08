package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type upgradeSourceAuthority struct {
	LockSchema int
	LockSHA256 string
}

func inspectUpgradeSource(observed state.Observed) (upgradeSourceAuthority, error) {
	schema := state.AsInt(observed.Lock["schema_version"])
	switch schema {
	case 1:
		if err := validateLegacyV1UpgradeSource(observed); err != nil {
			return upgradeSourceAuthority{}, err
		}
	case 2:
		if observed.Classification != state.StateCurrent && observed.Classification != state.StateDrifted && observed.Classification != state.StateStaleCLI {
			return upgradeSourceAuthority{}, fmt.Errorf("upgrade requires a valid lock-v2 installation")
		}
		if err := validateTargetManagedAdditions(observed); err != nil {
			return upgradeSourceAuthority{}, err
		}
	default:
		return upgradeSourceAuthority{}, fmt.Errorf("upgrade requires a supported lock installation")
	}
	return readUpgradeSourceAuthority(observed.CanonicalRoot, schema)
}

// inspectTargetUpgradeSource validates the bound source installation after the staged target
// binary starts. A valid lock-v2 installation can be partial only because the target graph added
// paths that were not present in the previous lock; paths declared by that lock must still exist.
func inspectTargetUpgradeSource(observed state.Observed) (upgradeSourceAuthority, error) {
	schema := state.AsInt(observed.Lock["schema_version"])
	if schema == 1 {
		if err := validateLegacyV1UpgradeSource(observed); err != nil {
			return upgradeSourceAuthority{}, err
		}
		return readUpgradeSourceAuthority(observed.CanonicalRoot, schema)
	}
	if schema != 2 || observed.LockSchemaVersion != 2 || observed.Config == nil || len(observed.Errors) > 0 {
		return upgradeSourceAuthority{}, fmt.Errorf("upgrade requires a valid lock-v2 installation")
	}
	switch observed.Classification {
	case state.StateCurrent, state.StateDrifted, state.StateStaleCLI:
	case state.StatePartial:
		declared := map[string]bool{}
		for _, field := range []string{"managed_paths", "managed_sections", "generated_surfaces"} {
			for _, item := range state.AnySlice(observed.Lock[field]) {
				if target := state.AsString(state.Map(item)["path"]); target != "" {
					declared[target] = true
				}
			}
		}
		for _, target := range observed.MissingPaths {
			if declared[target] {
				return upgradeSourceAuthority{}, fmt.Errorf("upgrade source is missing a lock-declared path")
			}
		}
	default:
		return upgradeSourceAuthority{}, fmt.Errorf("upgrade requires a valid lock-v2 installation")
	}
	if state.AsString(observed.Lock["selected_profile"]) != state.AsString(observed.Config["selected_profile"]) {
		return upgradeSourceAuthority{}, fmt.Errorf("upgrade source profile identity is incompatible")
	}
	if err := validateTargetManagedAdditions(observed); err != nil {
		return upgradeSourceAuthority{}, err
	}
	return readUpgradeSourceAuthority(observed.CanonicalRoot, schema)
}

func readUpgradeSourceAuthority(canonicalRoot string, schema int) (upgradeSourceAuthority, error) {
	lockPath := filepath.Join(canonicalRoot, filepath.FromSlash(state.LockPath))
	info, err := os.Lstat(lockPath)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return upgradeSourceAuthority{}, fmt.Errorf("upgrade source lock is not a regular file")
	}
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return upgradeSourceAuthority{}, fmt.Errorf("read upgrade source lock: %w", err)
	}
	digest := sha256.Sum256(data)
	return upgradeSourceAuthority{LockSchema: schema, LockSHA256: hex.EncodeToString(digest[:])}, nil
}

func boundUpgradeSourceAuthorityCheck(observed state.Observed, expected upgradeSourceAuthority) func() ([]reconcile.Blocker, error) {
	return func() ([]reconcile.Blocker, error) {
		actual, err := readUpgradeSourceAuthority(observed.CanonicalRoot, expected.LockSchema)
		legacyErr := error(nil)
		if expected.LockSchema == 1 {
			legacyErr = validateLegacyV1UpgradeSource(observed)
		}
		if err == nil && legacyErr == nil && actual == expected {
			return nil, nil
		}
		return []reconcile.Blocker{{
			Code:         "upgrade_source_lock_changed",
			Message:      "upgrade source authority changed after handoff validation",
			Path:         state.LockPath,
			Remediation:  "preserve the changed lock and run check before retrying",
			RetryCommand: "check",
		}}, nil
	}
}

func validateLegacyV1UpgradeSource(observed state.Observed) error {
	if observed.LockSchemaVersion != 1 || state.AsInt(observed.Lock["schema_version"]) != 1 || observed.Config == nil || observed.Graph.ProfileID == "" || len(observed.Errors) > 0 {
		return fmt.Errorf("legacy lock-v1 upgrade requires a schema-valid installation and config")
	}
	switch observed.Classification {
	case state.StateCurrent, state.StateDrifted, state.StateStaleCLI, state.StatePartial, state.StateLegacyV1Compatible:
	default:
		return fmt.Errorf("legacy lock-v1 installation state is not upgrade-compatible")
	}
	profile := state.AsString(observed.Lock["selected_profile"])
	if profile == "" || profile != observed.Graph.ProfileID || profile != state.AsString(observed.Config["selected_profile"]) {
		return fmt.Errorf("legacy lock-v1 profile identity is incompatible")
	}
	if !sameLegacyStringSet(state.AnySlice(observed.Lock["selected_components"]), observed.Graph.SelectedComponents) {
		return fmt.Errorf("legacy lock-v1 component identity is incompatible")
	}

	managed := state.AnySlice(observed.Lock["managed_paths"])
	if len(managed) == 0 {
		return fmt.Errorf("legacy lock-v1 managed-path authority is empty")
	}
	managedPaths := map[string]bool{}
	for _, item := range managed {
		record := state.Map(item)
		target := state.AsString(record["path"])
		if record == nil || state.AsString(record["ownership"]) != "managed" || state.AsString(record["component"]) == "" || state.AsString(record["source"]) == "" {
			return fmt.Errorf("legacy lock-v1 managed-path authority is incomplete")
		}
		key := portableLegacyPathKey(target)
		if !safeLegacyManagedPath(target) || managedPaths[key] || reconcile.ValidateFileTarget(observed.CanonicalRoot, target) != nil {
			return fmt.Errorf("legacy lock-v1 managed-path authority is unsafe or ambiguous")
		}
		managedPaths[key] = true
		if !containsLegacyString(observed.Graph.SelectedComponents, state.AsString(record["component"])) {
			return fmt.Errorf("legacy lock-v1 managed-path component is incompatible")
		}
		expected := strings.ToLower(state.AsString(record["checksum_sha256"]))
		if !validUpgradeSHA256(expected) || strings.Trim(expected, "0") == "" {
			return fmt.Errorf("legacy lock-v1 managed-path checksum is unusable")
		}
		actual, err := readContainedRegularFile(observed.CanonicalRoot, target)
		if err != nil {
			return fmt.Errorf("legacy lock-v1 managed path is missing, unsafe, or unreadable")
		}
		digest := sha256.Sum256(actual)
		if hex.EncodeToString(digest[:]) != expected {
			return fmt.Errorf("legacy lock-v1 managed path is modified")
		}
	}
	if err := validateTargetManagedAdditions(observed); err != nil {
		return err
	}

	requiredSurfaces := map[string]bool{
		".codeheart/kit/": false,
		state.LockPath:    false,
		state.ConfigPath:  false,
		"AGENTS.md":       false,
	}
	generatedPaths := map[string]bool{}
	for _, item := range state.AnySlice(observed.Lock["generated_surfaces"]) {
		record := state.Map(item)
		target := state.AsString(record["path"])
		key := portableLegacyPathKey(target)
		if record == nil || state.AsString(record["ownership"]) == "" || !safeLegacySurfacePath(target) || generatedPaths[key] || managedPaths[key] {
			return fmt.Errorf("legacy lock-v1 generated-surface authority is unsafe or ambiguous")
		}
		generatedPaths[key] = true
		if _, required := requiredSurfaces[target]; required {
			requiredSurfaces[target] = true
		}
	}
	for _, present := range requiredSurfaces {
		if !present {
			return fmt.Errorf("legacy lock-v1 generated-surface authority is incomplete")
		}
	}
	agents, err := readContainedRegularFile(observed.CanonicalRoot, "AGENTS.md")
	agentsText := string(agents)
	begin := strings.Index(agentsText, reconcile.BeginMarker)
	end := strings.Index(agentsText, reconcile.EndMarker)
	if err != nil || strings.Count(agentsText, reconcile.BeginMarker) != 1 || strings.Count(agentsText, reconcile.EndMarker) != 1 || begin < 0 || end < 0 || begin >= end {
		return fmt.Errorf("legacy lock-v1 managed section is missing or ambiguous")
	}
	return nil
}

func validateTargetManagedAdditions(observed state.Observed) error {
	declared := map[string]bool{}
	for _, item := range state.AnySlice(observed.Lock["managed_paths"]) {
		if target := state.AsString(state.Map(item)["path"]); target != "" {
			declared[portableLegacyPathKey(target)] = true
		}
	}
	for _, node := range observed.Graph.Nodes {
		if node.DirectoryTarget || node.Ownership != state.OwnershipManaged || node.Update != state.UpdateReplace || declared[portableLegacyPathKey(node.Target)] {
			continue
		}
		if err := reconcile.ValidateFileTarget(observed.CanonicalRoot, node.Target); err != nil {
			return fmt.Errorf("target-added managed path is unsafe")
		}
		target := filepath.Join(observed.CanonicalRoot, filepath.FromSlash(node.Target))
		if _, err := os.Lstat(target); err == nil {
			return fmt.Errorf("target-added managed path overlaps existing consumer content")
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("target-added managed path cannot be inspected safely")
		}
	}
	return nil
}

func sameLegacyStringSet(values []any, expected []string) bool {
	actual := make([]string, 0, len(values))
	for _, value := range values {
		text, ok := value.(string)
		if !ok || text == "" {
			return false
		}
		actual = append(actual, text)
	}
	if len(actual) != len(expected) {
		return false
	}
	sort.Strings(actual)
	want := append([]string(nil), expected...)
	sort.Strings(want)
	for index := range actual {
		if actual[index] != want[index] {
			return false
		}
	}
	return true
}

func containsLegacyString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func safeLegacyManagedPath(target string) bool {
	return strings.HasPrefix(target, ".codeheart/kit/") && !strings.HasSuffix(target, "/") && safeLegacySurfacePath(target)
}

func safeLegacySurfacePath(target string) bool {
	if target == "" || strings.Contains(target, "\\") || strings.HasPrefix(target, "/") {
		return false
	}
	trimmed := strings.TrimSuffix(target, "/")
	return trimmed != "" && trimmed != "." && pathpkg.Clean(trimmed) == trimmed && !strings.HasPrefix(trimmed, "../")
}

func portableLegacyPathKey(target string) string {
	return strings.ToLower(strings.TrimSuffix(pathpkg.Clean(target), "/"))
}

func readContainedRegularFile(root, target string) ([]byte, error) {
	if err := reconcile.ValidateFileTarget(root, target); err != nil {
		return nil, err
	}
	current := root
	parts := strings.Split(filepath.FromSlash(target), string(filepath.Separator))
	for _, part := range parts {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("unsafe path")
		}
	}
	info, err := os.Lstat(current)
	if err != nil || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("path is not a regular file")
	}
	return os.ReadFile(current)
}

func validUpgradeSHA256(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
