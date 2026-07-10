package drift

import "github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"

// Report is the retained compatibility projection over the authoritative state classifier.
// The lock argument remains for callers of the v1 API; state.Inspect reads and validates the
// installed lock together with the declaration graph.
func Report(root string, _ map[string]any) []map[string]string {
	observed, err := state.Inspect(root)
	if err != nil {
		return []map[string]string{{"path": root, "status": "invalid", "detail": err.Error()}}
	}
	findings := []map[string]string{}
	for _, path := range observed.MissingPaths {
		findings = append(findings, map[string]string{"path": path, "status": "missing"})
	}
	for _, path := range observed.DriftedPaths {
		findings = append(findings, map[string]string{"path": path, "status": "drift"})
	}
	for _, detail := range observed.Errors {
		findings = append(findings, map[string]string{"path": state.LockPath, "status": "invalid", "detail": detail})
	}
	return findings
}
