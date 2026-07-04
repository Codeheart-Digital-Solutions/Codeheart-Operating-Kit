package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/capabilities"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/components"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/drift"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/lockfile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

var routeTargetPattern = regexp.MustCompile(`\.codeheart/kit/[A-Za-z0-9._/\-]+\.md`)

func RunCheck(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{"--json": false})
	_ = values
	if err != nil {
		return writeArgError(stderr, "check", err)
	}
	path := "."
	if len(positionals) > 0 {
		path = positionals[0]
	}
	if len(positionals) > 1 {
		return writeArgError(stderr, "check", fmt.Errorf("unexpected argument %q", positionals[1]))
	}
	result, err := CheckRepository(path)
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit check: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, result); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit check: error: %v\n", err)
			return 1
		}
	} else {
		fmt.Fprintln(stdout, "Operating Kit check")
		fmt.Fprintf(stdout, "OK: %v\n", result["ok"])
		fmt.Fprintf(stdout, "Missing CLI: %v\n", result["missing_cli"])
		fmt.Fprintf(stdout, "Stale CLI: %v\n", result["stale_cli"])
		fmt.Fprintf(stdout, "Missing routing: %v\n", result["missing_routing"])
		fmt.Fprintf(stdout, "Missing route targets: %d\n", len(result["missing_route_targets"].([]any)))
		fmt.Fprintf(stdout, "Drift findings: %d\n", len(result["drift"].([]any)))
	}
	if ok, _ := result["ok"].(bool); ok {
		return 0
	}
	return 1
}

func CheckRepository(path string) (map[string]any, error) {
	root := expandPath(path)
	lock, err := lockfile.ReadLock(root)
	if err != nil {
		return nil, err
	}
	missingLockMetadata := lockfile.MissingRequiredLockMetadata(lock)
	missingCLI := !CLIAvailable()
	agentsText := ""
	if data, err := os.ReadFile(filepath.Join(root, "AGENTS.md")); err == nil {
		agentsText = string(data)
	}
	missingTargets := MissingRouteTargets(root, agentsText)
	missingRouting := !stringsContains(agentsText, components.BeginMarker) || !stringsContains(agentsText, components.EndMarker) || len(missingTargets) > 0
	staleCLI := valueString(lock["kit_version"]) != "" && valueString(lock["kit_version"]) != version.Version
	driftFindings := drift.Report(root, lock)
	native := lock["native_capabilities"]
	if native == nil {
		now := lockfile.UTCNow()
		native = capabilities.UnknownNativeCapabilityState(now)
	}
	ok := len(missingLockMetadata) == 0 && !missingRouting && len(driftFindings) == 0 && !missingCLI && !staleCLI
	return map[string]any{
		"ok":                    ok,
		"missing_cli":           missingCLI,
		"stale_cli":             staleCLI,
		"missing_routing":       missingRouting,
		"missing_route_targets": stringsAsAny(missingTargets),
		"missing_lock_metadata": stringsAsAny(missingLockMetadata),
		"drift":                 mapStringListAsAny(driftFindings),
		"native_capabilities":   native,
	}, nil
}

func MissingRouteTargets(root string, agentsText string) []string {
	targets := routeTargetPattern.FindAllString(agentsText, -1)
	seen := map[string]bool{}
	missing := []string{}
	for _, target := range targets {
		if seen[target] {
			continue
		}
		seen[target] = true
		if !exists(filepath.Join(root, filepath.FromSlash(target))) {
			missing = append(missing, target)
		}
	}
	sort.Strings(missing)
	return missing
}

func CLIAvailable() bool {
	if os.Getenv("CODEHEART_OPERATING_KIT_CLI") == "1" {
		return true
	}
	if _, err := exec.LookPath("codeheart-operating-kit"); err == nil {
		return true
	}
	executable := filepath.Base(os.Args[0])
	return (executable == "codeheart-operating-kit" || executable == "codeheart-operating-kit.exe") && exists(os.Args[0])
}

func mapStringListAsAny(values []map[string]string) []any {
	result := make([]any, len(values))
	for index, value := range values {
		mapping := map[string]any{}
		for key, nested := range value {
			mapping[key] = nested
		}
		result[index] = mapping
	}
	return result
}

func stringsContains(value string, needle string) bool {
	return strings.Contains(value, needle)
}
