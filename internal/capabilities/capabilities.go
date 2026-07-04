package capabilities

import "time"

var BaselineCapabilities = []string{"documents", "spreadsheets", "presentations", "browser", "pdf"}

func UnknownNativeCapabilityState(checkedAt time.Time) map[string]any {
	state := map[string]any{}
	for _, capability := range BaselineCapabilities {
		state[capability] = map[string]any{
			"status":                  "unknown",
			"checked_at":              checkedAt.UTC().Truncate(time.Second).Format(time.RFC3339),
			"profile_applicability":   "standard",
			"command_result_category": "not-checked",
		}
	}
	return state
}
