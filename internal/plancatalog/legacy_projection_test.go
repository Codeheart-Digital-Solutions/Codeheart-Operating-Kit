package plancatalog

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

func TestLegacyRegisterProjectsRecognizedStatusesAndSeparateDates(t *testing.T) {
	data := legacyProjectionFixture(t, "recognized-statuses-and-dates.md")
	entries, problems := ParseLegacyRegister(data)
	if HasErrors(problems) || len(entries) != 6 {
		t.Fatalf("entries=%#v problems=%#v", entries, problems)
	}
	byID := legacyEntriesByID(entries)
	wantLifecycle := map[string]Lifecycle{
		"GENERIC-001": LifecycleDraft,
		"GENERIC-002": LifecycleActive,
		"GENERIC-003": LifecycleCompleted,
		"GENERIC-004": LifecycleSuperseded,
		"GENERIC-005": LifecycleArchived,
	}
	for id, lifecycle := range wantLifecycle {
		if byID[id].Lifecycle != lifecycle || byID[id].LegacyStatus != "" {
			t.Fatalf("%s projected status=%q legacy=%q", id, byID[id].Lifecycle, byID[id].LegacyStatus)
		}
	}
	handoff := byID["GENERIC-006"]
	if handoff.Lifecycle != "" || handoff.LegacyStatus != LegacyStatusImplementationHandoffReady {
		t.Fatalf("pre-canonical status projection=%#v", handoff)
	}
	completed := byID["GENERIC-003"]
	if completed.LastUpdated != "2026-03-03T03:04:05Z" || completed.Completed != "2026-03-02" {
		t.Fatalf("separate dates were not retained: %#v", completed)
	}
	encoded, err := json.Marshal(handoff)
	if err != nil || bytes.Contains(encoded, []byte(`"lifecycle"`)) || !bytes.Contains(encoded, []byte(`"legacy_status":"implementation-handoff-ready"`)) {
		t.Fatalf("closed handoff wire projection=%s err=%v", encoded, err)
	}
}

func TestLegacyRegisterMalformedAndAmbiguousFormsBlockWithoutInvalidWireValues(t *testing.T) {
	data := legacyProjectionFixture(t, "malformed-and-ambiguous.md")
	entries, problems := ParseLegacyRegister(data)
	for _, code := range []string{"legacy_date_malformed", "legacy_field_duplicate", "legacy_status_unsupported", "legacy_canonical_path_malformed", "legacy_canonical_path_duplicate", "legacy_relation_unsupported", "legacy_relation_malformed", "legacy_relation_duplicate"} {
		if !problemExists(problems, code, SeverityError) {
			t.Fatalf("missing %s in %#v", code, problems)
		}
	}
	byID := legacyEntriesByID(entries)
	if entry := byID["GENERIC-NEG-001"]; entry.LastUpdated != "" || entry.Completed != "" || entry.Lifecycle != LifecycleCompleted {
		t.Fatalf("combined date leaked into wire fields: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-002"]; entry.Lifecycle != "" || entry.Created != "" || len(entry.CanonicalDocs) != 0 {
		t.Fatalf("duplicated fields were selected arbitrarily: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-003"]; entry.Lifecycle != "" || entry.LegacyStatus != "" {
		t.Fatalf("unsupported status leaked into a closed field: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-004"]; len(entry.CanonicalDocs) != 0 || entry.Created != "" || entry.LastUpdated != "" || entry.Completed != "" {
		t.Fatalf("malformed values leaked into wire fields: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-007"]; len(entry.Relations) != 0 {
		t.Fatalf("malformed relations leaked into wire fields: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-008"]; len(entry.CanonicalDocs) != 0 {
		t.Fatalf("backslash path leaked into schema-v3 wire fields: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-009"]; len(entry.CanonicalDocs) != 0 {
		t.Fatalf("duplicate path retained reconciliation authority: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-010"]; len(entry.CanonicalDocs) != 0 {
		t.Fatalf("rooted path leaked into schema-v3 wire fields: %#v", entry)
	}
	legacyV1ByID := legacyEntriesByID(legacyEntriesForInventoryV1(entries))
	rawV1Path := `docs\repo\plans\generic-backslash\generic-backslash_discovery_doc.md`
	wantV1Path := filepath.ToSlash(filepath.Clean(filepath.FromSlash(rawV1Path)))
	if entry := legacyV1ByID["GENERIC-NEG-008"]; len(entry.CanonicalDocs) != 1 || entry.CanonicalDocs[0] != wantV1Path {
		t.Fatalf("schema-v1 backslash path compatibility changed: %#v", entry)
	}
	if entry := legacyV1ByID["GENERIC-NEG-009"]; len(entry.CanonicalDocs) != 2 || entry.CanonicalDocs[0] != entry.CanonicalDocs[1] {
		t.Fatalf("schema-v1 duplicate path compatibility changed: %#v", entry)
	}
	reconciliation := ReconcileLegacy([]Record{{Path: "docs/repo/plans/generic-shared/generic-shared_discovery_doc.md"}}, nil, entries)
	if !problemExists(reconciliation.Problems, "legacy_evidence_ambiguous", SeverityError) {
		t.Fatalf("ambiguous path evidence was not blocked: %#v", reconciliation.Problems)
	}
	secondEntries, secondProblems := ParseLegacyRegister(data)
	if !reflect.DeepEqual(entries, secondEntries) || !reflect.DeepEqual(problems, secondProblems) {
		t.Fatalf("legacy projection is nondeterministic\nfirst=%#v %#v\nsecond=%#v %#v", entries, problems, secondEntries, secondProblems)
	}

	root := migrationRepository(t)
	writeLegacyProjectionRegister(t, root, data)
	inventory, err := BuildInventoryWithOptions(root, time.Date(2026, 8, 10, 1, 2, 3, 0, time.UTC), SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, TargetCatalogMode: ModeCanonical})
	if err != nil {
		t.Fatal(err)
	}
	validateLegacyProjectionInventory(t, inventory)
}

func TestFrozenRegisterHistoricalFieldsAreWarningsOnlyForCanonicalAuthority(t *testing.T) {
	allowed := []string{"legacy_status_unsupported", "legacy_relation_malformed", "legacy_relation_unsupported"}
	problems := []Problem{
		{Code: allowed[0], Path: LegacyRegisterPath, Severity: SeverityError},
		{Code: allowed[1], Path: LegacyRegisterPath, Severity: SeverityError},
		{Code: allowed[2], Path: LegacyRegisterPath, Severity: SeverityError},
		{Code: "legacy_relation_duplicate", Path: LegacyRegisterPath, Severity: SeverityError},
		{Code: allowed[0], Path: "docs/repo/plans/not-the-register.md", Severity: SeverityError},
	}
	for _, mode := range []CatalogMode{ModeLegacy, ModeMixed} {
		if projected := projectFrozenRegisterCompatibility(problems, mode); !reflect.DeepEqual(projected, problems) {
			t.Fatalf("%s mode changed strict legacy evidence:\ngot=%#v\nwant=%#v", mode, projected, problems)
		}
	}
	projected := projectFrozenRegisterCompatibility(problems, ModeCanonical)
	for _, code := range allowed {
		if !problemExists(projected, code, SeverityWarning) {
			t.Fatalf("canonical mode did not retain %s as a warning: %#v", code, projected)
		}
	}
	if !problemExists(projected, "legacy_relation_duplicate", SeverityError) {
		t.Fatalf("canonical mode weakened authority-critical legacy evidence: %#v", projected)
	}
	if projected[4].Severity != SeverityError {
		t.Fatalf("canonical mode weakened a problem outside the frozen register: %#v", projected[4])
	}
	for index := range problems {
		if problems[index].Severity != SeverityError {
			t.Fatalf("projection mutated its input: %#v", problems)
		}
	}
}

func TestLegacyRelationProseDoesNotInventDuplicateTargets(t *testing.T) {
	cases := []struct {
		name, body           string
		malformed, duplicate bool
		want                 []Relation
	}{
		{"wrapped prose", "- related: first consumer discovery handoff -\n  <consumer>/docs/discovery.md\n- related: first consumer implementation handoff -\n  <consumer>/docs/implementation.md", true, false, nil},
		{"unwrapped prose", "- related: first consumer discovery\n- related: first consumer implementation", true, false, nil},
		{"duplicate IDs with different titles", "- related: PR-001 - Discovery\n- related: PR-001 - Renamed discovery", false, true, nil},
		{"duplicate paths", "- related: docs/discovery.md\n- related: docs/discovery.md", false, true, nil},
		{"targets and optional titles", "- related: PR-001 - Discovery\n- depends-on: other:PR-002\n- related: docs/discovery.md", false, false, []Relation{{Kind: "related", Target: "PR-001"}, {Kind: "depends-on", Target: "other:PR-002"}, {Kind: "related", Target: "docs/discovery.md"}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte("## PR-010 - Example\nCanonical docs: docs/example.md\nRelations:\n" + tc.body + "\n")
			entries, problems := ParseLegacyRegister(data)
			if len(entries) != 1 || !reflect.DeepEqual(entries[0].Relations, tc.want) {
				t.Fatalf("relations=%#v want=%#v", entries, tc.want)
			}
			if problemExists(problems, "legacy_relation_malformed", SeverityError) != tc.malformed || problemExists(problems, "legacy_relation_duplicate", SeverityError) != tc.duplicate {
				t.Fatalf("incorrect classification: %#v", problems)
			}
			for _, mode := range []CatalogMode{ModeLegacy, ModeMixed} {
				if !reflect.DeepEqual(projectFrozenRegisterCompatibility(problems, mode), problems) {
					t.Fatalf("%s weakened strict evidence", mode)
				}
			}
			canonical := projectFrozenRegisterCompatibility(problems, ModeCanonical)
			if HasErrors(canonical) != tc.duplicate {
				t.Fatalf("canonical duplicate guard changed: %#v", canonical)
			}
			// The schema-v1 compatibility projection deliberately retains its old
			// first-token behavior, including stopping at a wrapped continuation.
			v1 := legacyEntriesForInventoryV1(entries)
			lines := strings.Split(string(data), "\n")
			if !reflect.DeepEqual(v1[0].Relations, readLegacyRelationsV1(lines, 0, len(lines))) {
				t.Fatalf("schema-v1 relation compatibility changed: %#v", v1)
			}
			if tc.name == "wrapped prose" && !reflect.DeepEqual(v1[0].Relations, []Relation{{Kind: "related", Target: "first"}}) {
				t.Fatalf("schema-v1 wrapped projection changed: %#v", v1)
			}
		})
	}
}

func TestCanonicalMigrationRetainsMalformedFrozenFieldsWithoutRegisterWrite(t *testing.T) {
	root := legacyV2MigrationRepository(t)
	register := []byte(`Last updated: 2026-08-11T00:00:00Z (UTC)

# Plan Register

## Entries

## PR-001 - Alpha Discovery

Status: awaiting-review
Canonical docs: docs/repo/plans/alpha/alpha_discovery_doc.md
Created: 2026-07-30
Last updated: 2026-07-31T09:00:00Z (UTC)

Relations:
related: PR-002

## PR-002 - Beta Implementation

Status: active
Canonical docs: docs/repo/plans/beta/beta_implementation_doc.md
Created: 2026-07-31
Last updated: 2026-07-31T09:05:00Z (UTC)

Relations:
- obsolete-kind: PR-001
	`)
	writeLegacyProjectionRegister(t, root, register)
	adoptMixedCatalogTest(t, root)
	registerPath := filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))
	registerBefore := mustReadFile(t, registerPath)

	mixed, err := LoadRepositorySnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"legacy_status_unsupported", "legacy_relation_malformed", "legacy_relation_unsupported"} {
		if !problemExists(mixed.Problems, code, SeverityError) {
			t.Fatalf("mixed mode did not retain strict %s evidence: %#v", code, mixed.Problems)
		}
	}

	prospective, err := LoadRepositorySnapshotWithOptions(root, SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, TargetCatalogMode: ModeCanonical})
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"legacy_status_unsupported", "legacy_relation_malformed", "legacy_relation_unsupported"} {
		if !problemExists(prospective.Problems, code, SeverityWarning) || problemExists(prospective.Problems, code, SeverityError) {
			t.Fatalf("prospective canonical mode did not project %s as warning-only: %#v", code, prospective.Problems)
		}
	}

	ledger := v2MigrationLedger(t, root, ModeCanonical, nil, nil)
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"legacy_status_unsupported", "legacy_relation_malformed", "legacy_relation_unsupported"} {
		if problemExists(plan.Problems, code, SeverityError) {
			t.Fatalf("reviewed canonical migration remained blocked on %s: %#v", code, plan.Problems)
		}
	}
	if !plan.Projection.Ready {
		t.Fatalf("otherwise-complete reviewed canonical migration is not ready: %#v", plan.Problems)
	}
	if after := mustReadFile(t, registerPath); !bytes.Equal(after, registerBefore) {
		t.Fatal("prospective inventory or migration planning changed the frozen register")
	}

	configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	config := mustReadFile(t, configPath)
	cutover := cutoverRevisionFromConfig(t, config)
	old := "    plan_catalog_mode: mixed\n    plan_catalog_cutover_revision: " + cutover + "\n"
	replacement := "    plan_catalog_mode: canonical\n    plan_catalog_discovery_version: 2\n    plan_catalog_ownership:\n      excluded_roots: []\n"
	configured := bytes.Replace(config, []byte(old), []byte(replacement), 1)
	if bytes.Equal(configured, config) {
		t.Fatal("mixed fixture config did not contain the expected cutover block")
	}
	if err := os.WriteFile(configPath, configured, 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", state.ConfigPath)
	runGitTest(t, root, "commit", "-m", "configure canonical catalog")
	activated, err := LoadRepositorySnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"legacy_status_unsupported", "legacy_relation_malformed", "legacy_relation_unsupported"} {
		if !problemExists(activated.Problems, code, SeverityWarning) || problemExists(activated.Problems, code, SeverityError) {
			t.Fatalf("configured canonical mode did not retain %s as warning-only: %#v", code, activated.Problems)
		}
	}
	if after := mustReadFile(t, registerPath); !bytes.Equal(after, registerBefore) {
		t.Fatal("canonical configuration changed the frozen register")
	}
}

func TestLegacyRegisterSchemaV3ProjectionHasModeParityDeterministicHashingAndZeroWrites(t *testing.T) {
	root := migrationRepository(t)
	register := legacyProjectionFixture(t, "recognized-statuses-and-dates.md")
	writeLegacyProjectionRegister(t, root, register)
	fixedNow := time.Date(2026, 8, 10, 1, 2, 3, 0, time.UTC)
	legacyV1, err := BuildInventory(root, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	if legacyV1.SchemaVersion != 1 {
		t.Fatalf("legacy inventory schema=%d", legacyV1.SchemaVersion)
	}
	legacyV1ByID := legacyInventoryProjection(legacyV1)
	if entry := legacyV1ByID["GENERIC-006"]; entry.Lifecycle != Lifecycle(LegacyStatusImplementationHandoffReady) || entry.LegacyStatus != "" {
		t.Fatalf("schema-v1 pre-canonical status compatibility changed: %#v", entry)
	}
	if entry := legacyV1ByID["GENERIC-003"]; entry.Completed != "" || entry.LastUpdated != "2026-03-03T03:04:05Z (UTC) Completed: 2026-03-02" {
		t.Fatalf("schema-v1 combined date compatibility changed: %#v", entry)
	}
	legacyV1JSON, _ := json.Marshal(legacyV1)
	if bytes.Contains(legacyV1JSON, []byte(`"legacy_status":`)) || bytes.Contains(legacyV1JSON, []byte(`"completed":`)) {
		t.Fatalf("schema-v3 legacy fields leaked into schema-v1 inventory: %s", legacyV1JSON)
	}

	prospective, err := BuildInventoryWithOptions(root, fixedNow, SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, TargetCatalogMode: ModeCanonical})
	if err != nil {
		t.Fatal(err)
	}
	validateLegacyProjectionInventory(t, prospective)

	remoteAware, err := BuildInventoryWithOptions(root, fixedNow, SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, TargetCatalogMode: ModeCanonical})
	if err != nil {
		t.Fatal(err)
	}
	overlay := RemoteOverlayEvidence{
		Remote: "example", SourceIdentitySHA256: strings.Repeat("a", 64), Status: "complete",
		ObservedAt: fixedNow.Format(time.RFC3339), Refs: []RemoteOverlayRef{}, BranchTouchCandidates: []BranchCandidateEvidence{},
	}
	overlay.Digest = CanonicalRemoteOverlayDigest(overlay)
	AttachRemoteOverlay(root, &remoteAware, overlay)
	validateLegacyProjectionInventory(t, remoteAware)

	configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	config := mustReadFile(t, configPath)
	cutover := cutoverRevisionFromConfig(t, config)
	old := "    plan_catalog_mode: mixed\n    plan_catalog_cutover_revision: " + cutover + "\n"
	replacement := "    plan_catalog_mode: canonical\n    plan_catalog_discovery_version: 2\n    plan_catalog_ownership:\n      excluded_roots: []\n"
	updated := strings.Replace(string(config), old, replacement, 1)
	if updated == string(config) {
		t.Fatal("mixed fixture config did not contain the expected cutover block")
	}
	if err := os.WriteFile(configPath, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", state.ConfigPath)
	runGitTest(t, root, "commit", "-m", "configure local discovery v2")

	protected := map[string][]byte{
		LegacyRegisterPath: mustReadFile(t, filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))),
		"docs/repo/plans/alpha/alpha_discovery_doc.md":    mustReadFile(t, filepath.Join(root, "docs/repo/plans/alpha/alpha_discovery_doc.md")),
		"docs/repo/plans/beta/beta_implementation_doc.md": mustReadFile(t, filepath.Join(root, "docs/repo/plans/beta/beta_implementation_doc.md")),
	}
	local, err := BuildInventory(root, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	validateLegacyProjectionInventory(t, local)
	second, err := BuildInventory(root, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	firstJSON, _ := json.Marshal(local)
	secondJSON, _ := json.Marshal(second)
	if !bytes.Equal(firstJSON, secondJSON) || local.CandidateSetDigest != second.CandidateSetDigest || local.BranchEvidence.Digest != second.BranchEvidence.Digest {
		t.Fatal("equivalent inventory observations did not produce deterministic bytes and digests")
	}

	wantProjection := legacyInventoryProjection(prospective)
	for name, inventory := range map[string]Inventory{"local": local, "prospective": prospective, "remote-aware": remoteAware} {
		if got := legacyInventoryProjection(inventory); !reflect.DeepEqual(got, wantProjection) {
			t.Fatalf("%s legacy projection diverged\ngot=%#v\nwant=%#v", name, got, wantProjection)
		}
	}
	for relative, before := range protected {
		if after := mustReadFile(t, filepath.Join(root, filepath.FromSlash(relative))); !bytes.Equal(before, after) {
			t.Fatalf("inventory changed protected consumer bytes at %s", relative)
		}
	}
	if status := runGitTest(t, root, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("inventory left repository writes: %q", status)
	}

	ambiguous := prospective
	for recordIndex := range ambiguous.Records {
		for evidenceIndex := range ambiguous.Records[recordIndex].LegacyEvidence {
			if ambiguous.Records[recordIndex].LegacyEvidence[evidenceIndex].Lifecycle != "" {
				ambiguous.Records[recordIndex].LegacyEvidence[evidenceIndex].LegacyStatus = LegacyStatusImplementationHandoffReady
				if legacyProjectionInventoryValid(ambiguous) {
					t.Fatal("schema accepted simultaneous canonical lifecycle and pre-canonical status")
				}
				return
			}
		}
	}
	t.Fatal("fixture did not provide canonical legacy evidence for ambiguity check")
}

func legacyProjectionFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "legacy-register", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func writeLegacyProjectionRegister(t *testing.T, root string, data []byte) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", LegacyRegisterPath)
	runGitTest(t, root, "commit", "-m", "use generic legacy register evidence")
}

func legacyEntriesByID(entries []LegacyEntry) map[string]LegacyEntry {
	result := map[string]LegacyEntry{}
	for _, entry := range entries {
		result[entry.ID] = entry
	}
	return result
}

func legacyInventoryProjection(inventory Inventory) map[string]LegacyEntry {
	result := legacyEntriesByID(inventory.UnpairedLegacy)
	for _, record := range inventory.Records {
		for _, entry := range record.LegacyEvidence {
			result[entry.ID] = entry
		}
	}
	return result
}

func validateLegacyProjectionInventory(t *testing.T, inventory Inventory) {
	t.Helper()
	if inventory.SchemaVersion != 3 || !legacyProjectionInventoryValid(inventory) {
		data, _ := json.MarshalIndent(inventory, "", "  ")
		t.Fatalf("schema-v3 inventory projection is invalid:\n%s", data)
	}
}

func legacyProjectionInventoryValid(inventory Inventory) bool {
	data, err := json.Marshal(inventory)
	if err != nil {
		return false
	}
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return false
	}
	return state.Validate(state.PlanInventoryV3Schema, value) == nil
}
