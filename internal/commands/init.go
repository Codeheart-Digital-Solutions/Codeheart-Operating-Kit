package commands

import (
	"fmt"
	"io"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/components"
)

func RunInit(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{
		"--project-name":    true,
		"--purpose":         true,
		"--selected-folder": true,
		"--json":            false,
	})
	if err != nil {
		return writeArgError(stderr, "init", err)
	}
	if err := validatePurpose(values["--purpose"]); err != nil {
		return writeArgError(stderr, "init", err)
	}
	path := "."
	if len(positionals) > 0 {
		path = positionals[0]
	}
	if len(positionals) > 1 {
		return writeArgError(stderr, "init", fmt.Errorf("unexpected argument %q", positionals[1]))
	}
	projectName := values["--project-name"]
	if projectName == "" {
		projectName = "Example-Automation"
	}
	selectedFolder := values["--selected-folder"]
	if selectedFolder == "" {
		selectedFolder = expandPath(path)
	}
	result, err := Initialize(path, projectName, values["--purpose"], selectedFolder, time.Now())
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit init: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, result); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit init: error: %v\n", err)
			return 1
		}
		return 0
	}
	fmt.Fprintln(stdout, "Operating Kit initialized.")
	fmt.Fprintf(stdout, "Mode: %s\n", result["inspection"].(map[string]any)["mode"])
	if report, ok := result["adoption_cleanup_report"].(string); ok && report != "" {
		fmt.Fprintf(stdout, "Adoption cleanup report: %s\n", report)
	}
	return 0
}

func Initialize(path string, projectName string, purpose string, selectedFolder string, now time.Time) (map[string]any, error) {
	root := expandPath(path)
	if err := osMkdirAll(root); err != nil {
		return nil, err
	}
	inspection := InspectFolder(root)
	preexisting := []string{}
	for _, candidate := range []string{"AGENTS.md", "docs/repo/README.md", "docs/agent-memory/README.md"} {
		if exists(joinRoot(root, candidate)) {
			preexisting = append(preexisting, candidate)
		}
	}
	state, err := components.WriteDefaultState(root, projectName, purpose, selectedFolder, now)
	if err != nil {
		return nil, err
	}
	reportPath := ""
	if inspection["mode"] == "existing-technical-project-adoption" && len(preexisting) > 0 {
		reportPath, err = components.WriteAdoptionCleanupReport(root, preexisting)
		if err != nil {
			return nil, err
		}
	}
	return map[string]any{
		"inspection":              inspection,
		"state":                   stateToMap(state),
		"adoption_cleanup_report": reportPath,
	}, nil
}

func stateToMap(state components.DefaultState) map[string]any {
	return map[string]any{
		"managed_paths":      mapsToAny(state.ManagedPaths),
		"generated_surfaces": mapsToAny(state.GeneratedSurfaces),
		"agents_status":      state.AgentsStatus,
		"gitignore_changed":  state.GitignoreChanged,
	}
}
