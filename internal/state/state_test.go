package state

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/kitfs"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

func TestYAMLRoundTripPreservesNumericLookingStrings(t *testing.T) {
	input := map[string]any{
		"checksum":  strings.Repeat("0", 64),
		"short":     "0",
		"timestamp": "2026-07-09T20:00:00Z",
		"nested":    map[string]any{"value": "00123"},
	}
	data, err := EncodeYAML(input)
	if err != nil {
		t.Fatalf("EncodeYAML: %v", err)
	}
	decoded, err := DecodeYAMLMap(data)
	if err != nil {
		t.Fatalf("DecodeYAMLMap: %v", err)
	}
	if !reflect.DeepEqual(input, decoded) {
		t.Fatalf("round trip changed values\ninput: %#v\noutput: %#v\nYAML:\n%s", input, decoded, data)
	}
}

func TestCompileGraphIsDeterministicAndUsesStateDefaults(t *testing.T) {
	first, err := CompileGraph("standard")
	if err != nil {
		t.Fatalf("CompileGraph first: %v", err)
	}
	second, err := CompileGraph("standard")
	if err != nil {
		t.Fatalf("CompileGraph second: %v", err)
	}
	if first.DigestSHA256 != second.DigestSHA256 || !reflect.DeepEqual(first.Nodes, second.Nodes) {
		t.Fatalf("graph is not deterministic")
	}
	foundManaged := false
	foundScaffold := false
	foundRoot := false
	for _, node := range first.Nodes {
		switch node.Target {
		case ".codeheart/kit/docs/agent-interface/README.md":
			foundManaged = node.Ownership == OwnershipManaged && node.Update == UpdateReplace && node.Removal == RemovalReconcile
		case "docs/agent-memory/README.md":
			foundScaffold = node.Ownership == OwnershipScaffold && node.Update == UpdatePreserve
		case "AGENTS.md":
			foundRoot = node.Update == UpdateManagedSection && node.RouteID == "root-agents"
		}
	}
	if !foundManaged || !foundScaffold || !foundRoot {
		t.Fatalf("graph missing expected semantics: managed=%v scaffold=%v root=%v", foundManaged, foundScaffold, foundRoot)
	}
}

func TestNormalizeAndMigrateKnownLegacyChecksumPlaceholders(t *testing.T) {
	lock := representativeV1Lock()
	Map(lock["release"])["checksum_sha256"] = 0
	Map(lock["cli_repair"])["repair_checksum_sha256"] = "0"
	normalized, anomalies := NormalizeLegacyV1(lock)
	if len(anomalies) != 2 {
		t.Fatalf("anomalies = %#v, want two", anomalies)
	}
	if err := Validate(LockV1Schema, normalized); err != nil {
		t.Fatalf("normalized lock v1: %v", err)
	}
	migrated, migratedAnomalies, err := MigrateLockV1(lock, time.Date(2026, 7, 9, 20, 0, 0, 0, time.UTC), "tx-1", "repair")
	if err != nil {
		t.Fatalf("MigrateLockV1: %v", err)
	}
	if len(migratedAnomalies) != 2 || AsInt(migrated["schema_version"]) != 2 {
		t.Fatalf("migration result = %#v; anomalies=%#v", migrated, migratedAnomalies)
	}
	if Map(migrated["release_provenance"])["verification_status"] != "unverified-legacy" {
		t.Fatalf("release provenance = %#v", migrated["release_provenance"])
	}
}

func TestUnrelatedInvalidV1FieldFailsClosed(t *testing.T) {
	lock := representativeV1Lock()
	lock["unexpected"] = true
	if _, _, err := MigrateLockV1(lock, time.Now(), "tx", "repair"); err == nil {
		t.Fatalf("expected unrelated invalid v1 field to fail")
	}
}

func TestRepositoryStateFixturesCoverLegacyInvalidFutureAndOptionalConfig(t *testing.T) {
	legacy := loadStateFixture(t, "lock-v1-legacy-zero.yaml")
	migrated, anomalies, err := MigrateLockV1(legacy, time.Date(2026, 7, 9, 20, 0, 0, 0, time.UTC), "fixture", "repair")
	if err != nil || len(anomalies) != 2 || AsInt(migrated["schema_version"]) != 2 {
		t.Fatalf("legacy fixture migration = %#v, anomalies=%#v, err=%v", migrated, anomalies, err)
	}
	invalid := loadStateFixture(t, "lock-v1-invalid-unknown.yaml")
	if _, _, err := MigrateLockV1(invalid, time.Now(), "fixture", "repair"); err == nil {
		t.Fatalf("invalid fixture unexpectedly migrated")
	}
	futureData := readStateFixture(t, "lock-future.yaml")
	future, err := DecodeYAMLMap(futureData)
	if err != nil || AsInt(future["schema_version"]) != 99 {
		t.Fatalf("future fixture = %#v, err=%v", future, err)
	}
	configData := readStateFixture(t, "config-v1-optional.yaml")
	config, err := DecodeAndValidateYAML(ConfigV1Schema, configData)
	if err != nil {
		t.Fatalf("optional config fixture: %v", err)
	}
	encoded, err := EncodeYAML(config)
	if err != nil {
		t.Fatal(err)
	}
	roundTrip, err := DecodeAndValidateYAML(ConfigV1Schema, encoded)
	if err != nil || !reflect.DeepEqual(config, roundTrip) {
		t.Fatalf("optional config round trip changed: %#v %#v err=%v", config, roundTrip, err)
	}
}

func TestInspectAbsentAndAdoptable(t *testing.T) {
	root := t.TempDir()
	absent, err := Inspect(filepath.Join(root, "missing"))
	if err != nil || absent.Classification != StateAbsent {
		t.Fatalf("absent = %#v, err=%v", absent, err)
	}
	adoptable, err := Inspect(root)
	if err != nil || adoptable.Classification != StateAdoptable {
		t.Fatalf("adoptable = %#v, err=%v", adoptable, err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("example\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	adoptable, err = Inspect(root)
	if err != nil || adoptable.Classification != StateAdoptable {
		t.Fatalf("existing folder adoptable = %#v, err=%v", adoptable, err)
	}
}

func TestInspectClassifiesPartialInvalidFutureTransactionAndRecovery(t *testing.T) {
	partialRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(partialRoot, ".codeheart", "kit"), 0o755); err != nil {
		t.Fatal(err)
	}
	partial, err := Inspect(partialRoot)
	if err != nil || partial.Classification != StatePartial {
		t.Fatalf("partial = %#v, err=%v", partial, err)
	}

	invalidRoot := t.TempDir()
	materializeRequiredState(t, invalidRoot, representativeV1Lock())
	if err := os.WriteFile(filepath.Join(invalidRoot, filepath.FromSlash(LockPath)), []byte("schema_version: invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	invalid, err := Inspect(invalidRoot)
	if err != nil || invalid.Classification != StateSchemaInvalid {
		t.Fatalf("invalid = %#v, err=%v", invalid, err)
	}

	futureRoot := t.TempDir()
	future := representativeV1Lock()
	future["schema_version"] = 99
	materializeRequiredState(t, futureRoot, future)
	futureState, err := Inspect(futureRoot)
	if err != nil || futureState.Classification != StateUnsupportedFutureVersion {
		t.Fatalf("future = %#v, err=%v", futureState, err)
	}
	incompleteFutureRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(incompleteFutureRoot, ".codeheart"), 0o755); err != nil {
		t.Fatal(err)
	}
	futureData, err := EncodeYAML(future)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(incompleteFutureRoot, filepath.FromSlash(LockPath)), futureData, 0o644); err != nil {
		t.Fatal(err)
	}
	incompleteFuture, err := Inspect(incompleteFutureRoot)
	if err != nil || incompleteFuture.Classification != StateUnsupportedFutureVersion {
		t.Fatalf("incomplete future = %#v, err=%v", incompleteFuture, err)
	}

	transactionRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(transactionRoot, ".codeheart"), 0o755); err != nil {
		t.Fatal(err)
	}
	marker := []byte(`{"schema_version":1,"transaction_id":"tx","phase":"staging"}`)
	if err := os.WriteFile(filepath.Join(transactionRoot, filepath.FromSlash(TransactionPath)), marker, 0o600); err != nil {
		t.Fatal(err)
	}
	transaction, err := Inspect(transactionRoot)
	if err != nil || transaction.Classification != StateTransactionInProgress {
		t.Fatalf("transaction = %#v, err=%v", transaction, err)
	}
	marker = []byte(`{"schema_version":1,"transaction_id":"tx","phase":"recovery-required"}`)
	if err := os.WriteFile(filepath.Join(transactionRoot, filepath.FromSlash(TransactionPath)), marker, 0o600); err != nil {
		t.Fatal(err)
	}
	recovery, err := Inspect(transactionRoot)
	if err != nil || recovery.Classification != StateRecoveryRequired {
		t.Fatalf("recovery = %#v, err=%v", recovery, err)
	}
}

func TestInspectClassifiesCurrentDriftedStaleAndLegacy(t *testing.T) {
	currentRoot := t.TempDir()
	currentLock := representativeV1Lock()
	currentLock["kit_version"] = version.Version
	Map(currentLock["release"])["asset_url"] = "https://example.invalid/kit.zip"
	Map(currentLock["release"])["checksum_sha256"] = strings.Repeat("a", 64)
	Map(currentLock["cli_repair"])["repair_checksum_sha256"] = strings.Repeat("b", 64)
	materializeRequiredState(t, currentRoot, currentLock)
	current, err := Inspect(currentRoot)
	if err != nil || current.Classification != StateCurrent {
		t.Fatalf("current = %#v, err=%v", current, err)
	}

	driftTarget := filepath.Join(currentRoot, ".codeheart", "kit", "docs", "agent-interface", "README.md")
	if err := os.WriteFile(driftTarget, []byte("drift\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	drifted, err := Inspect(currentRoot)
	if err != nil || drifted.Classification != StateDrifted {
		t.Fatalf("drifted = %#v, err=%v", drifted, err)
	}

	staleRoot := t.TempDir()
	staleLock := DeepCopy(currentLock)
	staleLock["kit_version"] = "0.0.1"
	materializeRequiredState(t, staleRoot, staleLock)
	stale, err := Inspect(staleRoot)
	if err != nil || stale.Classification != StateStaleCLI {
		t.Fatalf("stale = %#v, err=%v", stale, err)
	}

	legacyRoot := t.TempDir()
	legacyLock := DeepCopy(currentLock)
	Map(legacyLock["cli_repair"])["repair_checksum_sha256"] = 0
	materializeRequiredState(t, legacyRoot, legacyLock)
	legacy, err := Inspect(legacyRoot)
	if err != nil || legacy.Classification != StateLegacyV1Compatible {
		t.Fatalf("legacy = %#v, err=%v", legacy, err)
	}
}

func representativeV1Lock() map[string]any {
	return map[string]any{
		"schema_version":      1,
		"kit_version":         "0.1.21",
		"selected_profile":    "standard",
		"selected_components": []any{"agent-interface"},
		"release": map[string]any{
			"asset_url":       "local-source",
			"checksum_sha256": strings.Repeat("0", 64),
		},
		"managed_paths":      []any{},
		"generated_surfaces": []any{},
		"cli_repair": map[string]any{
			"installed_cli_path":     "codeheart-operating-kit",
			"repair_source_url":      "local-source",
			"repair_checksum_sha256": strings.Repeat("0", 64),
		},
		"update_check": map[string]any{
			"last_update_check_at":  "2026-07-09T20:00:00Z",
			"next_update_check_due": "2026-07-16T20:00:00Z",
			"latest_seen_version":   "0.1.21",
			"update_status":         "current",
		},
		"native_capabilities": map[string]any{},
	}
}

func materializeRequiredState(t *testing.T, root string, lock map[string]any) {
	t.Helper()
	graph, err := CompileGraph("standard")
	if err != nil {
		t.Fatal(err)
	}
	for _, node := range graph.Nodes {
		if node.Presence != PresenceRequired || node.Target == LockPath || node.Target == ConfigPath {
			continue
		}
		target := filepath.Join(root, filepath.FromSlash(strings.TrimSuffix(node.Target, "/")))
		if node.DirectoryTarget {
			if err := os.MkdirAll(target, 0o755); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if node.Source == "" || !kitfs.Exists(node.Source) {
			continue
		}
		data, err := kitfs.ReadFile(node.Source)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	lockData, err := EncodeYAML(lock)
	if err != nil {
		t.Fatal(err)
	}
	lockPath := filepath.Join(root, filepath.FromSlash(LockPath))
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lockPath, lockData, 0o644); err != nil {
		t.Fatal(err)
	}
	config := map[string]any{
		"schema_version":        1,
		"selected_profile":      "standard",
		"project_display_name":  "Example",
		"selected_setup_folder": root,
		"local_consumer_layer": map[string]any{
			"repo_docs_path":           "docs/repo/",
			"agent_memory_path":        "docs/agent-memory/",
			"user_layer_path":          ".codeheart/user/",
			"local_machine_layer_path": ".codeheart/local/",
		},
		"component_settings": map[string]any{},
	}
	configData, err := EncodeYAML(config)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(ConfigPath)), configData, 0o644); err != nil {
		t.Fatal(err)
	}
}

func readStateFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "tests", "fixtures", "state", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func loadStateFixture(t *testing.T, name string) map[string]any {
	t.Helper()
	value, err := DecodeYAMLMap(readStateFixture(t, name))
	if err != nil {
		t.Fatal(err)
	}
	return value
}
