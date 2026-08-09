package reconcile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

func TestValidateStagedRootDispatchesConfigSchemaVersion(t *testing.T) {
	for _, version := range []int{1, 2} {
		t.Run("schema-v"+string(rune('0'+version)), func(t *testing.T) {
			config := testConfig(t.TempDir())
			config["schema_version"] = version
			if version == 2 {
				config["component_settings"] = map[string]any{
					"planning-workflows": map[string]any{
						"plan_catalog_mode":              "mixed",
						"plan_catalog_cutover_revision":  "1111111111111111111111111111111111111111",
						"plan_catalog_discovery_version": 2,
						"plan_catalog_migration_evidence": map[string]any{
							"ledger_path":              "docs/repo/plans/migrations/example.yaml",
							"ledger_sha256":            "2222222222222222222222222222222222222222222222222222222222222222",
							"branch_evidence_digest":   "3333333333333333333333333333333333333333333333333333333333333333",
							"evidence_revision":        "4444444444444444444444444444444444444444",
							"activation_base_revision": "5555555555555555555555555555555555555555",
							"migration_action_digest":  "6666666666666666666666666666666666666666666666666666666666666666",
							"evidence_scope":           "local",
						},
					},
				}
			}
			content, err := state.EncodeYAML(config)
			if err != nil {
				t.Fatal(err)
			}
			transactionPath := t.TempDir()
			stagePath := filepath.Join(transactionPath, "stage", filepath.FromSlash(state.ConfigPath))
			if err := os.MkdirAll(filepath.Dir(stagePath), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(stagePath, content, 0o600); err != nil {
				t.Fatal(err)
			}
			transactionRoot, err := os.OpenRoot(transactionPath)
			if err != nil {
				t.Fatal(err)
			}
			defer transactionRoot.Close()
			plan := Plan{Actions: []Action{{Kind: "replace", Target: state.ConfigPath, Content: content}}}
			if err := validateStagedRoot(plan, transactionRoot); err != nil {
				t.Fatalf("schema-v%d staged config validation failed: %v", version, err)
			}
		})
	}
}
