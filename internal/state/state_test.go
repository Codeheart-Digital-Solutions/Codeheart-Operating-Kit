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

func TestVersionedSchemaDispatchPreservesHistoricalContracts(t *testing.T) {
	cases := []struct {
		name string
		got  string
		err  error
		want string
	}{
		{"config-v1", mustSchema(SchemaForConfigVersion(1)), nil, ConfigV1Schema},
		{"config-v2", mustSchema(SchemaForConfigVersion(2)), nil, ConfigV2Schema},
		{"catalog-v2", mustSchema(SchemaForPlanCatalogVersion(2)), nil, PlanCatalogV2Schema},
		{"catalog-v3", mustSchema(SchemaForPlanCatalogVersion(3)), nil, PlanCatalogV3Schema},
		{"ledger-v2", mustSchema(SchemaForPlanMigrationVersion(2)), nil, PlanMigrationV2Schema},
		{"ledger-v3", mustSchema(SchemaForPlanMigrationVersion(3)), nil, PlanMigrationV3Schema},
		{"inventory-v3", mustSchema(SchemaForPlanInventoryVersion(3)), nil, PlanInventoryV3Schema},
	}
	for _, tc := range cases {
		if tc.err != nil || tc.got != tc.want {
			t.Fatalf("%s: got %q err=%v, want %q", tc.name, tc.got, tc.err, tc.want)
		}
	}
	if _, err := SchemaForConfigVersion(3); err == nil {
		t.Fatal("config version 3 unexpectedly accepted")
	}
	if _, err := SchemaForPlanCatalogVersion(4); err == nil {
		t.Fatal("catalog version 4 unexpectedly accepted")
	}
	if _, err := SchemaForPlanMigrationVersion(4); err == nil {
		t.Fatal("migration version 4 unexpectedly accepted")
	}
}

func TestConfigValidationDispatchesRawVersionAndClosesV2Binding(t *testing.T) {
	v1 := representativeConfig(1)
	if err := ValidateConfig(v1); err != nil {
		t.Fatalf("config v1: %v", err)
	}
	v2 := representativeConfig(2)
	planning := map[string]any{
		"plan_catalog_mode":              "mixed",
		"plan_catalog_discovery_version": 2,
		"plan_catalog_cutover_revision":  strings.Repeat("a", 40),
		"plan_catalog_migration_evidence": map[string]any{
			"ledger_path":              "docs/repo/plans/example/attachments/ledger-v3.yaml",
			"ledger_sha256":            strings.Repeat("b", 64),
			"branch_evidence_digest":   strings.Repeat("c", 64),
			"evidence_revision":        strings.Repeat("d", 40),
			"activation_base_revision": strings.Repeat("e", 40),
			"migration_action_digest":  strings.Repeat("f", 64),
			"evidence_scope":           "local",
		},
	}
	v2["component_settings"] = map[string]any{"planning-workflows": planning}
	encoded, err := EncodeYAML(v2)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeAndValidateConfigYAML(encoded); err != nil {
		t.Fatalf("config v2: %v", err)
	}

	for name, mutate := range map[string]func(map[string]any){
		"canonical": func(value map[string]any) {
			Map(Map(value["component_settings"])["planning-workflows"])["plan_catalog_mode"] = "canonical"
		},
		"missing-evidence": func(value map[string]any) {
			delete(Map(Map(value["component_settings"])["planning-workflows"]), "plan_catalog_migration_evidence")
		},
		"unknown-evidence": func(value map[string]any) {
			Map(Map(Map(value["component_settings"])["planning-workflows"])["plan_catalog_migration_evidence"])["reviewed"] = true
		},
		"future-version": func(value map[string]any) { value["schema_version"] = 99 },
	} {
		invalid := DeepCopy(v2)
		mutate(invalid)
		if err := ValidateConfig(invalid); err == nil {
			t.Fatalf("%s config unexpectedly accepted", name)
		}
	}
	for _, invalid := range []map[string]any{{}, {"schema_version": "2"}} {
		if err := ValidateConfig(invalid); err == nil {
			t.Fatalf("raw version %#v unexpectedly accepted", invalid)
		}
	}
}

func TestLedgerV3CandidateScopedProofAndDeferralConditionals(t *testing.T) {
	ledger := representativeLedgerV3()
	if err := Validate(PlanMigrationV3Schema, ledger); err != nil {
		t.Fatalf("valid same-content ledger: %v", err)
	}
	remoteAware := DeepCopy(ledger)
	remoteSummary := Map(remoteAware["branch_evidence"])
	remoteSummary["evidence_scope"] = "remote-aware"
	remoteSummary["remote_overlay_status"] = "complete"
	remoteSummary["remote_overlay_digest"] = strings.Repeat("6", 64)
	remoteSummary["remote_source_identity_sha256"] = strings.Repeat("7", 64)
	if err := Validate(PlanMigrationV3Schema, remoteAware); err != nil {
		t.Fatalf("valid remote-aware ledger: %v", err)
	}
	missingRemoteIdentity := DeepCopy(remoteAware)
	delete(Map(missingRemoteIdentity["branch_evidence"]), "remote_source_identity_sha256")
	if err := Validate(PlanMigrationV3Schema, missingRemoteIdentity); err == nil {
		t.Fatal("remote-aware ledger without source identity unexpectedly accepted")
	}
	localWithRemoteBinding := DeepCopy(ledger)
	Map(localWithRemoteBinding["branch_evidence"])["remote_source_identity_sha256"] = strings.Repeat("7", 64)
	if err := Validate(PlanMigrationV3Schema, localWithRemoteBinding); err == nil {
		t.Fatal("local ledger with remote source identity unexpectedly accepted")
	}
	for _, logicalRef := range []string{
		"tracking:origin:refs/heads/history/example",
		"remote:github:refs/heads/history/example",
	} {
		versionedRef := DeepCopy(ledger)
		review := Map(AnySlice(Map(AnySlice(versionedRef["records"])[0])["branch_touch_reviews"])[0])
		Map(review["identity"])["logical_ref"] = logicalRef
		if err := Validate(PlanMigrationV3Schema, versionedRef); err != nil {
			t.Fatalf("logical ref %q: %v", logicalRef, err)
		}
	}
	invalidLogicalRef := DeepCopy(ledger)
	reviewWithInvalidRef := Map(AnySlice(Map(AnySlice(invalidLogicalRef["records"])[0])["branch_touch_reviews"])[0])
	Map(reviewWithInvalidRef["identity"])["logical_ref"] = "tracking:refs/remotes/origin/history/example"
	if err := Validate(PlanMigrationV3Schema, invalidLogicalRef); err == nil {
		t.Fatal("tracking ref without remote identity unexpectedly accepted")
	}

	missingProof := DeepCopy(ledger)
	delete(Map(AnySlice(Map(AnySlice(missingProof["records"])[0])["branch_touch_reviews"])[0]), "proof")
	if err := Validate(PlanMigrationV3Schema, missingProof); err == nil {
		t.Fatal("same-content review without proof unexpectedly accepted")
	}

	incorporated := DeepCopy(ledger)
	incorporatedReview := Map(AnySlice(Map(AnySlice(incorporated["records"])[0])["branch_touch_reviews"])[0])
	incorporatedReview["disposition"] = "incorporated-history-non-owner"
	incorporatedReview["proof"] = map[string]any{
		"kind": "incorporated-history-v1", "transition_digest": strings.Repeat("1", 64),
		"stable_patch_id": strings.Repeat("2", 40), "incorporated_parent": strings.Repeat("3", 40),
		"incorporated_commit":            strings.Repeat("4", 40),
		"incorporated_transition_digest": strings.Repeat("1", 64),
		"incorporated_stable_patch_id":   strings.Repeat("2", 40),
	}
	if err := Validate(PlanMigrationV3Schema, incorporated); err != nil {
		t.Fatalf("valid incorporated-history proof: %v", err)
	}
	delete(Map(incorporatedReview["proof"]), "incorporated_commit")
	if err := Validate(PlanMigrationV3Schema, incorporated); err == nil {
		t.Fatal("partial incorporated-history proof unexpectedly accepted")
	}

	blocking := DeepCopy(ledger)
	blockingReview := Map(AnySlice(Map(AnySlice(blocking["records"])[0])["branch_touch_reviews"])[0])
	blockingReview["disposition"] = "blocking"
	blockingReview["blocking_reasons"] = []any{"multiple merge bases"}
	for _, field := range []string{"ref_tip", "merge_base", "path_transitions", "proof"} {
		delete(blockingReview, field)
	}
	if err := Validate(PlanMigrationV3Schema, blocking); err != nil {
		t.Fatalf("blocking review with unavailable proof inputs: %v", err)
	}

	deferred := representativeLedgerV3()
	record := Map(AnySlice(deferred["records"])[0])
	review := Map(AnySlice(record["branch_touch_reviews"])[0])
	review["disposition"] = "deferred-active-owner"
	delete(review, "proof")
	review["owner_tip_candidate"] = map[string]any{
		"path": "docs/repo/plans/example/example_implementation_doc.md",
		"mode": "100644", "object_id": strings.Repeat("1", 40),
		"sha256": strings.Repeat("2", 64), "lifecycle": "active",
	}
	review["incremental_follow_up"] = map[string]any{
		"action":                "reinventory-and-incremental-migrate",
		"trigger":               "owner-ref-integrated-or-target-path-changed",
		"cutover_revision":      strings.Repeat("3", 40),
		"plan_path":             "docs/repo/plans/example/example_implementation_doc.md",
		"cutover_source_sha256": strings.Repeat("4", 64),
		"owner_ref":             "local:refs/heads/feature/example",
		"owner_tip":             strings.Repeat("5", 40),
		"owner_path":            "docs/repo/plans/example/example_implementation_doc.md",
		"owner_mode":            "100644", "owner_object_id": strings.Repeat("1", 40),
		"owner_sha256": strings.Repeat("2", 64), "state": "pending",
		"success_conditions": []any{
			"owner-supplied-canonical-metadata-at-merge",
			"new-reviewed-ledger-applied-to-merged-target-bytes",
		},
	}
	record["deferred"] = true
	record["deferral_reason"] = "reviewed active owner remains branch-owned"
	if err := Validate(PlanMigrationV3Schema, deferred); err != nil {
		t.Fatalf("valid mixed deferral: %v", err)
	}
	conflictingOwner := DeepCopy(deferred)
	conflictingRecord := Map(AnySlice(conflictingOwner["records"])[0])
	reviews := AnySlice(conflictingRecord["branch_touch_reviews"])
	activeReview := DeepCopy(Map(reviews[0]))
	activeReview["disposition"] = "active-owner"
	delete(activeReview, "incremental_follow_up")
	conflictingRecord["branch_touch_reviews"] = append(reviews, activeReview)
	if err := Validate(PlanMigrationV3Schema, conflictingOwner); err == nil {
		t.Fatal("deferral with a second active owner unexpectedly accepted")
	}
	deferred["target_mode"] = "canonical"
	if err := Validate(PlanMigrationV3Schema, deferred); err == nil {
		t.Fatal("direct canonical deferral unexpectedly accepted")
	}
}

func TestInventoryAndCatalogV3SchemasValidateReadinessSurfaces(t *testing.T) {
	digest := strings.Repeat("b", 64)
	commit := strings.Repeat("a", 40)
	timestamp := "2026-08-09T10:00:00Z"
	inventory := map[string]any{
		"schema_version": 3, "repository_id": "example", "catalog_mode": "mixed",
		"source_revision": commit, "evidence_revision": commit, "generated_at": timestamp,
		"records": []any{}, "unpaired_legacy_evidence": []any{}, "problems": []any{},
		"coverage": map[string]any{
			"formal_records": 0, "canonical_metadata": 0, "legacy_records": 0,
			"invalid_records": 0, "unreadable_records": 0, "unsafe_records": 0,
			"unpaired_legacy_evidence": 0, "branch_touch_candidates": 0,
			"blocking_branch_touches": 0, "dirty_overlaps": 0, "candidates": 0,
			"included_candidates": 0, "excluded_candidates": 0,
			"blocked_candidates": 0, "unowned_candidates": 0,
		},
		"discovery_version": 2, "configured_discovery_version": 2,
		"configured_catalog_mode": "mixed", "target_catalog_mode": "mixed",
		"policy_digest": digest, "candidate_set_digest": strings.Repeat("c", 64),
		"target_config_precondition_sha256": strings.Repeat("d", 64),
		"complete":                          true, "mixed_coverage_complete": true, "canonical_ready": false,
		"candidates": []any{},
		"branch_evidence": map[string]any{
			"algorithm": "git-candidate-proof-v1", "git_version": "2.47.1",
			"evidence_scope": "local", "remote_overlay_status": "not-requested",
			"digest": strings.Repeat("e", 64),
		},
	}
	if err := Validate(PlanInventoryV3Schema, inventory); err != nil {
		t.Fatalf("inventory v3: %v", err)
	}
	remoteInventory := DeepCopy(inventory)
	remoteInventorySummary := Map(remoteInventory["branch_evidence"])
	remoteInventorySummary["evidence_scope"] = "remote-aware"
	remoteInventorySummary["remote_overlay_status"] = "complete"
	remoteInventorySummary["remote_overlay_digest"] = strings.Repeat("f", 64)
	remoteInventorySummary["remote_source_identity_sha256"] = strings.Repeat("1", 64)
	remoteInventory["remote_overlays"] = []any{map[string]any{
		"remote": "example", "source_identity_sha256": strings.Repeat("1", 64),
		"status": "complete", "observed_at": timestamp, "digest": strings.Repeat("f", 64),
		"refs": []any{}, "branch_touch_candidates": []any{},
	}}
	if err := Validate(PlanInventoryV3Schema, remoteInventory); err != nil {
		t.Fatalf("valid remote-aware inventory: %v", err)
	}
	missingInventorySourceIdentity := DeepCopy(remoteInventory)
	delete(Map(missingInventorySourceIdentity["branch_evidence"]), "remote_source_identity_sha256")
	if err := Validate(PlanInventoryV3Schema, missingInventorySourceIdentity); err == nil {
		t.Fatal("complete remote-aware inventory without summary source identity unexpectedly accepted")
	}
	missingOverlaySourceIdentity := DeepCopy(remoteInventory)
	delete(Map(AnySlice(missingOverlaySourceIdentity["remote_overlays"])[0]), "source_identity_sha256")
	if err := Validate(PlanInventoryV3Schema, missingOverlaySourceIdentity); err == nil {
		t.Fatal("complete remote overlay without source identity unexpectedly accepted")
	}
	catalog := map[string]any{
		"schema_version": 3, "discovery_version": 2, "policy_digest": digest,
		"candidate_set_digest": strings.Repeat("c", 64), "coordination_home_id": "example-home",
		"started_at": timestamp, "completed_at": timestamp, "complete": true,
		"mixed_coverage_complete": true, "canonical_ready": true,
		"members": []any{}, "observations": []any{}, "compatibility_observations": []any{},
		"candidates": []any{}, "errors": []any{},
		"metrics": map[string]any{
			"duration_ms": 0, "source_count": 0, "member_count": 0,
			"candidate_count": 0, "observation_count": 0,
			"compatibility_observation_count": 0, "stale_count": 0,
			"api_call_count": 0, "max_concurrency": 1,
		},
	}
	catalog["members"] = []any{map[string]any{
		"repository_id": "example", "source_kind": "github",
		"source_identity_sha256": strings.Repeat("1", 64),
		"default_ref":            "refs/remotes/origin/main", "source_revision": commit,
		"self_member": false, "discovery_version": 2, "policy_digest": digest,
		"candidate_set_digest": strings.Repeat("c", 64), "complete": true,
		"mixed_coverage_complete": true, "canonical_ready": true,
	}}
	if err := Validate(PlanCatalogV3Schema, catalog); err != nil {
		t.Fatalf("catalog v3: %v", err)
	}
	missingCatalogSourceIdentity := DeepCopy(catalog)
	delete(Map(AnySlice(missingCatalogSourceIdentity["members"])[0]), "source_identity_sha256")
	if err := Validate(PlanCatalogV3Schema, missingCatalogSourceIdentity); err == nil {
		t.Fatal("catalog member without source identity unexpectedly accepted")
	}
	catalog["compatibility_observations"] = []any{map[string]any{
		"repository_id": "example", "plan_id": "example.implementation.deferred",
		"title": "Deferred", "kind": "implementation", "purpose": "Retain active work",
		"lifecycle": "active", "canonical_path": "docs/repo/plans/deferred_implementation_doc.md",
		"ref": "refs/heads/main", "commit": commit, "content_sha256": digest,
		"coverage_disposition": "mixed-grandfathered", "incremental_migration_required": true,
		"migration_evidence": map[string]any{
			"ledger_path":   "docs/repo/plans/example/attachments/ledger-v3.yaml",
			"ledger_sha256": digest, "branch_evidence_digest": strings.Repeat("c", 64),
			"evidence_revision": commit, "activation_base_revision": strings.Repeat("d", 40),
			"migration_action_digest": strings.Repeat("e", 64), "evidence_scope": "local",
		},
		"verification": "verified", "observed_at": timestamp, "stale": false, "conflict": false,
	}}
	if err := Validate(PlanCatalogV3Schema, catalog); err == nil {
		t.Fatal("catalog with compatibility row and canonical_ready true unexpectedly accepted")
	}
	catalog["canonical_ready"] = false
	if err := Validate(PlanCatalogV3Schema, catalog); err != nil {
		t.Fatalf("catalog compatibility row: %v", err)
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

func mustSchema(path string, err error) string {
	if err != nil {
		panic(err)
	}
	return path
}

func representativeConfig(version int) map[string]any {
	return map[string]any{
		"schema_version":        version,
		"selected_profile":      "standard",
		"project_display_name":  "Example",
		"selected_setup_folder": "/tmp/example",
		"local_consumer_layer": map[string]any{
			"repo_docs_path":           "docs/repo/",
			"agent_memory_path":        "docs/agent-memory/",
			"user_layer_path":          ".codeheart/user/",
			"local_machine_layer_path": ".codeheart/local/",
		},
		"component_settings": map[string]any{},
	}
}

func representativeLedgerV3() map[string]any {
	digest := strings.Repeat("b", 64)
	commit := strings.Repeat("a", 40)
	path := "docs/repo/plans/example/example_implementation_doc.md"
	state := map[string]any{
		"state": "present", "mode": "100644",
		"object_id": strings.Repeat("1", 40), "sha256": digest,
	}
	return map[string]any{
		"schema_version": 3, "repository_id": "example-repository",
		"evidence_revision": commit, "target_mode": "mixed",
		"inventory_revision": commit, "policy_digest": digest,
		"candidate_set_digest":              strings.Repeat("c", 64),
		"target_config_precondition_sha256": strings.Repeat("d", 64),
		"branch_evidence": map[string]any{
			"algorithm": "git-candidate-proof-v1", "git_version": "2.47.1",
			"evidence_scope": "local", "remote_overlay_status": "not-requested",
			"digest": strings.Repeat("e", 64),
		},
		"records": []any{map[string]any{
			"current_path": path, "target_path": path,
			"source_revision": commit, "source_sha256": digest,
			"target_precondition": "same-path", "ownership_disposition": "owned",
			"decision": map[string]any{
				"id": "example-repository.implementation.example", "kind": "implementation",
				"purpose": "Implement the example", "first_cataloged": "2026-08-09T10:00:00Z",
				"catalog_metadata_updated": "2026-08-09T10:00:00Z",
			},
			"confidence": "high", "ambiguity": []any{},
			"evidence": []any{"reviewed source"}, "legacy_aliases": []any{},
			"conflicts": []any{}, "deferred": false,
			"branch_touch_reviews": []any{map[string]any{
				"identity": map[string]any{
					"evidence_scope": "local", "repository_id": "example-repository",
					"logical_ref": "local:refs/heads/history/example", "candidate_path": path,
				},
				"disposition": "same-content-non-owner", "ref_tip": strings.Repeat("f", 40),
				"merge_base": commit,
				"path_transitions": []any{map[string]any{
					"path": path, "before": DeepCopy(state), "after": DeepCopy(state),
				}},
				"proof": map[string]any{
					"kind": "same-content-v1", "target_path_state": DeepCopy(state),
					"ref_path_state": DeepCopy(state),
				},
			}},
		}},
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
