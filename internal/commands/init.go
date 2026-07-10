package commands

import (
	"fmt"
	"io"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

func RunInit(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{
		"--project-name":    true,
		"--purpose":         true,
		"--selected-folder": true,
		"--dry-run":         false,
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
	payload, operation, err := initializeOperation(path, projectName, values["--purpose"], selectedFolder, time.Now(), bools["--dry-run"])
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit init: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		if err := writeJSON(stdout, payload); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit init: error: %v\n", err)
			return 1
		}
	} else if bools["--dry-run"] || !operation.OK() {
		reconcile.WriteText(stdout, operation)
	} else {
		fmt.Fprintln(stdout, "Operating Kit initialized.")
		fmt.Fprintf(stdout, "Mode: %s\n", payload["inspection"].(map[string]any)["mode"])
		if report, ok := payload["adoption_cleanup_report"].(string); ok && report != "" {
			fmt.Fprintf(stdout, "Adoption cleanup report: %s\n", report)
		}
	}
	if operation.OK() {
		return 0
	}
	return 1
}

func Initialize(path string, projectName string, purpose string, selectedFolder string, now time.Time) (map[string]any, error) {
	payload, _, err := initializeOperation(path, projectName, purpose, selectedFolder, now, false)
	return payload, err
}

func initializeOperation(path, projectName, purpose, selectedFolder string, now time.Time, dryRun bool) (map[string]any, reconcile.Result, error) {
	root := expandPath(path)
	inspection := InspectFolder(root)
	observed, err := state.Inspect(root)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	if observed.Classification != state.StateAbsent && observed.Classification != state.StateAdoptable {
		result := blockedLifecycle("init", observed, "init_requires_absent_or_adoptable", "init cannot replace an existing or incomplete Operating Kit installation", "run check, then use repair for an existing installation", "repair")
		payload := resultPayload(result)
		payload["inspection"] = inspection
		return payload, result, nil
	}
	graph, err := state.CompileGraph("standard")
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	lock, err := desiredLifecycleLock("init", observed, graph, now)
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	preexisting := []string{}
	for _, candidate := range []string{"AGENTS.md", "docs/repo/README.md", "docs/agent-memory/README.md"} {
		if exists(joinRoot(root, candidate)) {
			preexisting = append(preexisting, candidate)
		}
	}
	extraFiles := map[string][]byte{}
	reportPath := ""
	if inspection["mode"] == "existing-technical-project-adoption" && len(preexisting) > 0 {
		reportPath = ".codeheart/reports/adoption-cleanup-report.md"
		extraFiles[reportPath] = adoptionReport(now, preexisting)
	}
	result, err := runLifecycle(lifecycleRequest{
		command:       "init",
		root:          root,
		now:           now,
		dryRun:        dryRun,
		observed:      observed,
		graph:         graph,
		desiredLock:   lock,
		desiredConfig: initialConfig(projectName, purpose, selectedFolder),
		extraFiles:    extraFiles,
		ensureIgnore:  true,
	})
	if err != nil {
		return nil, reconcile.Result{}, err
	}
	payload := resultPayload(result)
	payload["inspection"] = inspection
	payload["state"] = map[string]any{
		"managed_paths":      recordsAsMaps(reconcile.ManagedPathRecords(graph)),
		"generated_surfaces": recordsAsMaps(reconcile.GeneratedSurfaceRecords(graph)),
		"agents_status":      changeStatus(result, "AGENTS.md"),
		"gitignore_changed":  changeContains(result, ".gitignore"),
	}
	payload["adoption_cleanup_report"] = reportPath
	return payload, result, nil
}

func recordsAsMaps(records []any) []map[string]any {
	result := make([]map[string]any, 0, len(records))
	for _, record := range records {
		if mapping := state.Map(record); mapping != nil {
			result = append(result, mapping)
		}
	}
	return result
}

func changeStatus(result reconcile.Result, path string) string {
	if changeContains(result, path) {
		if result.DryRun {
			return "planned"
		}
		return "updated"
	}
	return "unchanged"
}
