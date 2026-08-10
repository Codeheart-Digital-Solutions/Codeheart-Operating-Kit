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
	for _, code := range []string{"legacy_date_malformed", "legacy_field_duplicate", "legacy_status_unsupported", "legacy_canonical_path_malformed"} {
		if !problemExists(problems, code, SeverityError) {
			t.Fatalf("missing %s in %#v", code, problems)
		}
	}
	byID := legacyEntriesByID(entries)
	if entry := byID["GENERIC-NEG-001"]; entry.LastUpdated != "" || entry.Completed != "" || entry.Lifecycle != LifecycleCompleted {
		t.Fatalf("combined date leaked into wire fields: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-002"]; entry.Lifecycle != "" || entry.Created != "" {
		t.Fatalf("duplicated fields were selected arbitrarily: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-003"]; entry.Lifecycle != "" || entry.LegacyStatus != "" {
		t.Fatalf("unsupported status leaked into a closed field: %#v", entry)
	}
	if entry := byID["GENERIC-NEG-004"]; len(entry.CanonicalDocs) != 0 || entry.Created != "" || entry.LastUpdated != "" || entry.Completed != "" {
		t.Fatalf("malformed values leaked into wire fields: %#v", entry)
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

func TestLegacyRegisterSchemaV3ProjectionHasModeParityDeterministicHashingAndZeroWrites(t *testing.T) {
	root := migrationRepository(t)
	register := legacyProjectionFixture(t, "recognized-statuses-and-dates.md")
	writeLegacyProjectionRegister(t, root, register)
	fixedNow := time.Date(2026, 8, 10, 1, 2, 3, 0, time.UTC)

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
