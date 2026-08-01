package plancatalog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

func TestLegacyReconciliationAndMixedViewHaveOneRowPerCanonicalPath(t *testing.T) {
	root := migrationRepository(t)
	snapshot, err := LoadRepositorySnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Settings.Mode != ModeMixed || snapshot.Settings.RepositoryID != "example" || len(snapshot.Records) != 2 {
		t.Fatalf("snapshot=%#v", snapshot)
	}
	if problemExists(snapshot.Problems, "mixed_new_plan_metadata_missing", SeverityError) || len(snapshot.Reconciliation.Unpaired) != 1 {
		t.Fatalf("mixed compatibility problems=%#v unpaired=%#v", snapshot.Problems, snapshot.Reconciliation.Unpaired)
	}
	view := BuildView(snapshot)
	if len(view.Rows) != 2 || view.Rows[0].CanonicalPath == view.Rows[1].CanonicalPath {
		t.Fatalf("view rows=%#v", view.Rows)
	}
	ids := []string{view.Rows[0].ID, view.Rows[1].ID}
	if strings.Join(ids, ",") != "legacy:PR-001,legacy:PR-002" {
		t.Fatalf("legacy identities=%#v", ids)
	}

	canonical := metadataForTest("example.discovery.alpha", KindDiscovery, "Canonical purpose")
	data := mustReadFile(t, filepath.Join(root, "docs/repo/plans/alpha/alpha_discovery_doc.md"))
	updated, err := InsertMetadata(data, canonical)
	if err != nil {
		t.Fatal(err)
	}
	record, err := ParseDocument("docs/repo/plans/alpha/alpha_discovery_doc.md", updated, KindDiscovery)
	if err != nil {
		t.Fatal(err)
	}
	snapshot.Records[0] = record
	view = BuildView(snapshot)
	if view.Rows[0].Purpose != "Canonical purpose" || view.Rows[0].ID != "example.discovery.alpha" {
		t.Fatalf("legacy evidence overrode canonical metadata: %#v", view.Rows[0])
	}
}

func TestListViewKeepsDiscoveryImplementationAndFamilyAsSeparateRows(t *testing.T) {
	root := filepath.Join("..", "..", "tests", "fixtures", "plans", "family-repository")
	snapshot, err := LoadRepositorySnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	view := BuildView(snapshot)
	if len(view.Rows) != 3 {
		t.Fatalf("rows=%#v problems=%#v", view.Rows, view.Problems)
	}
	kinds := map[Kind]int{}
	for _, row := range view.Rows {
		kinds[row.Kind]++
	}
	if kinds[KindDiscovery] != 1 || kinds[KindImplementation] != 1 || kinds[KindFamily] != 1 {
		t.Fatalf("typed rows=%#v", view.Rows)
	}
}

func TestDisplayTitleUsesReviewedLegacyTitleOnlyForOneCompatibilityMatch(t *testing.T) {
	record := Record{Header: Header{Title: "Overview", CompatibilityTitle: true}}
	match := LegacyEntry{ID: "PR-001", Title: "Semantic Plan Title"}
	if title := DisplayTitle(record, []LegacyEntry{match}); title != "Semantic Plan Title" {
		t.Fatalf("compatibility display title = %q", title)
	}
	record.Header.CompatibilityTitle = false
	if title := DisplayTitle(record, []LegacyEntry{match}); title != "Overview" {
		t.Fatalf("direct title was overridden = %q", title)
	}
	record.Header.CompatibilityTitle = true
	if title := DisplayTitle(record, []LegacyEntry{match, {ID: "PR-002", Title: "Competing Title"}}); title != "Overview" {
		t.Fatalf("ambiguous legacy evidence overrode title = %q", title)
	}
}

func TestQuotedMetadataMarkersDoNotBecomeCanonicalMetadata(t *testing.T) {
	data := []byte("Last updated: 2026-07-31T10:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: draft\n\n# Marker Example\n\n```md\n<!-- BEGIN CODEHEART PLAN METADATA -->\n```\n")
	record, err := ParseDocument("marker_discovery_doc.md", data, KindDiscovery)
	if err == nil || ErrorCode(err) != "metadata_missing" || !record.Legacy {
		t.Fatalf("record=%#v error=%v", record, err)
	}
}

func TestInsertMetadataIgnoresMarkersInsideFencedExamples(t *testing.T) {
	data := []byte("Last updated: 2026-07-31T10:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: draft\n\n# Marker Example\n\n````md\n<!-- BEGIN CODEHEART PLAN METADATA -->\n```yaml\nplan: example\n```\n<!-- END CODEHEART PLAN METADATA -->\n````\n")
	updated, err := InsertMetadata(data, metadataForTest("example.discovery.marker-example", KindDiscovery, "Document fenced marker examples"))
	if err != nil {
		t.Fatal(err)
	}
	record, err := ParseDocument("docs/repo/plans/marker-example/marker-example_discovery_doc.md", updated, KindDiscovery)
	if err != nil || record.Metadata == nil || record.Metadata.ID != "example.discovery.marker-example" {
		t.Fatalf("record=%#v error=%v\n%s", record, err, updated)
	}
	if bytes.Count(updated, []byte(MetadataBeginMarker)) != 2 || !bytes.Contains(updated, data[len(data)-len("````\n"):]) {
		t.Fatalf("fenced example was not preserved:\n%s", updated)
	}
}

func TestHistoricalImplementationHandoffStatusMigratesWithoutChangingHeader(t *testing.T) {
	root := migrationRepository(t)
	beta := "docs/repo/plans/beta/beta_implementation_doc.md"
	path := filepath.Join(root, filepath.FromSlash(beta))
	before := mustReadFile(t, path)
	before = bytes.Replace(before, []byte("Status: active"), []byte("Status: implementation-handoff-ready"), 1)
	if err := os.WriteFile(path, before, 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", beta)
	runGitTest(t, root, "commit", "-m", "use historical producer status")
	ledger := migrationLedger(t, root)
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil || len(plan.FilePlan.Actions) != 2 || skipExists(plan.Skips, "malformed_metadata") {
		t.Fatalf("migration plan=%#v skips=%#v err=%v", plan.FilePlan, plan.Skips, err)
	}
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)})
	if err != nil || outcome.Result.Status != reconcile.StatusSucceeded {
		t.Fatalf("migration outcome=%#v err=%v", outcome, err)
	}
	after := mustReadFile(t, path)
	if headerLine(after, "Status:") != "Status: implementation-handoff-ready" {
		t.Fatalf("historical status changed:\n%s", after)
	}
	record, parseErr := ParseDocument(beta, after, KindImplementation)
	problems := ValidateRecords([]Record{record}, ModeMixed, "example")
	if parseErr != nil || record.Header.Lifecycle != LifecycleDraft || record.Header.LegacyStatus != "implementation-handoff-ready" || !problemExists(problems, "legacy_header_status", SeverityWarning) {
		t.Fatalf("record=%#v problems=%#v error=%v", record, problems, parseErr)
	}
}

func TestUnsupportedStatusWithCanonicalMetadataStillFails(t *testing.T) {
	data := []byte("Last updated: 2026-07-31T10:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: draft\n\n# Unsupported Status\n\nBody.\n")
	canonical, err := InsertMetadata(data, metadataForTest("example.discovery.unsupported-status", KindDiscovery, "Reject unknown lifecycle values"))
	if err != nil {
		t.Fatal(err)
	}
	canonical = bytes.Replace(canonical, []byte("Status: draft"), []byte("Status: custom-workflow-state"), 1)
	if _, err := ParseDocument("docs/repo/plans/unsupported/unsupported_discovery_doc.md", canonical, KindDiscovery); err == nil || ErrorCode(err) != "header_invalid" {
		t.Fatalf("unsupported status error=%v", err)
	}
}

func TestLegacyHeaderCompatibilityNeverHidesAnUnmatchedMetadataMarker(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "docs", "repo", "plans", "invalid", "invalid_discovery_doc.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	data := []byte("Last updated: 2026-07-31T10:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: implementation-handoff-ready\n\n# Invalid Legacy Marker\n\n<!-- END CODEHEART PLAN METADATA -->\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := Discover(root, ModeMixed, "example")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 0 || !problemExists(result.Problems, "metadata_marker_count", SeverityError) {
		t.Fatalf("unmatched marker fell back to legacy: records=%#v problems=%#v", result.Records, result.Problems)
	}
}

func TestPlanAndRegisterSourcesMustBeContainedRegularFiles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("ordinary Windows test users may not have symlink privilege")
	}
	t.Run("formal plan symlink", func(t *testing.T) {
		root := migrationRepository(t)
		external := filepath.Join(t.TempDir(), "external_discovery_doc.md")
		externalData := []byte("Last updated: 2026-07-31T10:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: draft\n\n# External Plan Must Not Be Cataloged\n")
		if err := os.WriteFile(external, externalData, 0o644); err != nil {
			t.Fatal(err)
		}
		linked := "docs/repo/plans/linked/linked_discovery_doc.md"
		linkedPath := filepath.Join(root, filepath.FromSlash(linked))
		if err := os.MkdirAll(filepath.Dir(linkedPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(external, linkedPath); err != nil {
			t.Fatal(err)
		}
		snapshot, err := LoadRepositorySnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemExists(snapshot.Problems, "plan_source_unsafe", SeverityError) {
			t.Fatalf("unsafe plan problems=%#v", snapshot.Problems)
		}
		for _, record := range snapshot.Records {
			if record.Path == linked || record.Header.Title == "External Plan Must Not Be Cataloged" {
				t.Fatalf("external plan was cataloged: %#v", record)
			}
		}
		inventory, err := BuildInventory(root, time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatal(err)
		}
		record := inventoryByPath(inventory)[linked]
		if inventory.Coverage.FormalRecords != 3 || inventory.Coverage.UnsafeRecords != 1 || record.ParseStatus != "unsafe" || record.SourceSHA256 != "" {
			t.Fatalf("unsafe inventory coverage=%#v record=%#v", inventory.Coverage, record)
		}
	})

	t.Run("legacy register symlink", func(t *testing.T) {
		root := migrationRepository(t)
		register := filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))
		external := filepath.Join(t.TempDir(), "external-register.md")
		if err := os.WriteFile(external, mustReadFile(t, register), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(register); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(external, register); err != nil {
			t.Fatal(err)
		}
		snapshot, err := LoadRepositorySnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemExists(snapshot.Problems, "legacy_register_unsafe", SeverityError) || len(snapshot.LegacyEntries) != 0 {
			t.Fatalf("unsafe register entries=%#v problems=%#v", snapshot.LegacyEntries, snapshot.Problems)
		}
	})

	t.Run("legacy register symlinked parent", func(t *testing.T) {
		root := migrationRepository(t)
		plans := filepath.Join(root, "docs", "repo", "plans")
		externalPlans := filepath.Join(t.TempDir(), "external-plans")
		if err := os.Rename(plans, externalPlans); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(externalPlans, plans); err != nil {
			t.Fatal(err)
		}
		snapshot, err := LoadRepositorySnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemExists(snapshot.Problems, "legacy_register_unsafe", SeverityError) || len(snapshot.Records) != 0 {
			t.Fatalf("symlink-parent records=%#v problems=%#v", snapshot.Records, snapshot.Problems)
		}
	})

	t.Run("formal fifo", func(t *testing.T) {
		root := migrationRepository(t)
		fifo := filepath.Join(root, "docs", "repo", "plans", "fifo", "fifo_discovery_doc.md")
		if err := os.MkdirAll(filepath.Dir(fifo), 0o755); err != nil {
			t.Fatal(err)
		}
		if output, err := exec.Command("mkfifo", fifo).CombinedOutput(); err != nil {
			t.Fatalf("mkfifo: %v\n%s", err, output)
		}
		snapshot, err := LoadRepositorySnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemExists(snapshot.Problems, "plan_source_unsafe", SeverityError) {
			t.Fatalf("FIFO problems=%#v", snapshot.Problems)
		}
	})

	t.Run("identity replacement", func(t *testing.T) {
		root := migrationRepository(t)
		relative := "docs/repo/plans/alpha/alpha_discovery_doc.md"
		_, err := readRegularSourceWithHook(root, relative, func(target string) error {
			if err := os.Rename(target, target+".replaced"); err != nil {
				return err
			}
			return os.WriteFile(target, []byte("replacement inode\n"), 0o644)
		})
		if err == nil || ErrorCode(err) != "source_unsafe" {
			t.Fatalf("identity replacement error=%v", err)
		}
	})
}

func TestMixedModeRequiresPreCutoverPlanAndFrozenRegisterEvidence(t *testing.T) {
	root := migrationRepository(t)
	newPlan := "docs/repo/plans/post-cutover/post-cutover_discovery_doc.md"
	path := filepath.Join(root, filepath.FromSlash(newPlan))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("Last updated: 2026-07-31T12:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: draft\n\n# Post-Cutover Legacy Plan\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	register := filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))
	appendFile(t, register, "\n## PR-004 - Post-Cutover Legacy Plan\n\nType: discovery-plan\nPurpose: Attempt a forbidden mixed-mode register append.\nStatus: draft\nOwner / repository: Example\nCanonical docs: "+newPlan+"\nCreated: 2026-07-31\nLast updated: 2026-07-31T12:00:00Z (UTC)\n")
	runGitTest(t, root, "add", newPlan, LegacyRegisterPath)
	runGitTest(t, root, "commit", "-m", "attempt legacy authoring after cutover")
	snapshot, err := LoadRepositorySnapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if !problemExists(snapshot.Problems, "legacy_register_modified_after_cutover", SeverityError) || !problemForPathExists(snapshot.Problems, "mixed_new_plan_metadata_missing", newPlan, SeverityError) {
		t.Fatalf("mixed cutover problems=%#v", snapshot.Problems)
	}
}

func TestMigrationBlocksBrokenMixedCutoverInvariants(t *testing.T) {
	root := migrationRepository(t)
	ledger := migrationLedger(t, root)
	before := migrationPlanBytes(t, root)
	register := filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))
	appendFile(t, register, "\nForbidden post-cutover register mutation.\n")
	runGitTest(t, root, "add", LegacyRegisterPath)
	runGitTest(t, root, "commit", "-m", "mutate frozen register")
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)})
	if err != nil || outcome.Result.Status != reconcile.StatusBlocked || !reconcileBlockerExists(outcome.Result, "legacy_register_modified_after_cutover") {
		t.Fatalf("broken cutover outcome=%#v err=%v", outcome, err)
	}
	assertPlanBytes(t, root, before)
}

func TestMixedCutoverRequiresCommitAncestorAndExistingRegister(t *testing.T) {
	t.Run("configured object is a commit ancestor", func(t *testing.T) {
		root := migrationRepository(t)
		configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
		config := mustReadFile(t, configPath)
		baseline := cutoverRevisionFromConfig(t, config)
		blob := runGitTest(t, root, "rev-parse", "HEAD:"+LegacyRegisterPath)
		config = bytes.Replace(config, []byte(baseline), []byte(blob), 1)
		if err := os.WriteFile(configPath, config, 0o644); err != nil {
			t.Fatal(err)
		}
		snapshot, err := LoadRepositorySnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemExists(snapshot.Problems, "mixed_cutover_revision_invalid", SeverityError) {
			t.Fatalf("non-commit cutover problems=%#v", snapshot.Problems)
		}
	})

	t.Run("register remains present", func(t *testing.T) {
		root := migrationRepository(t)
		if err := os.Remove(filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))); err != nil {
			t.Fatal(err)
		}
		snapshot, err := LoadRepositorySnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemExists(snapshot.Problems, "legacy_register_missing_after_cutover", SeverityError) {
			t.Fatalf("missing register problems=%#v", snapshot.Problems)
		}
	})

	t.Run("revision is an ancestor", func(t *testing.T) {
		root := migrationRepository(t)
		tree := runGitTest(t, root, "rev-parse", "HEAD^{tree}")
		unrelated := runGitTest(t, root, "commit-tree", tree, "-m", "unrelated cutover")
		configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
		config := mustReadFile(t, configPath)
		baseline := cutoverRevisionFromConfig(t, config)
		config = bytes.Replace(config, []byte(baseline), []byte(unrelated), 1)
		if err := os.WriteFile(configPath, config, 0o644); err != nil {
			t.Fatal(err)
		}
		snapshot, err := LoadRepositorySnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemExists(snapshot.Problems, "mixed_cutover_revision_not_ancestor", SeverityError) {
			t.Fatalf("non-ancestor cutover problems=%#v", snapshot.Problems)
		}
	})

	t.Run("revision predates mixed mode", func(t *testing.T) {
		root := migrationRepository(t)
		mixedCommit := runGitTest(t, root, "rev-parse", "HEAD")
		configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
		config := mustReadFile(t, configPath)
		baseline := cutoverRevisionFromConfig(t, config)
		config = bytes.Replace(config, []byte(baseline), []byte(mixedCommit), 1)
		if err := os.WriteFile(configPath, config, 0o644); err != nil {
			t.Fatal(err)
		}
		snapshot, err := LoadRepositorySnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemExists(snapshot.Problems, "mixed_cutover_revision_not_legacy", SeverityError) {
			t.Fatalf("mixed baseline problems=%#v", snapshot.Problems)
		}
	})

	t.Run("baseline configuration is schema valid", func(t *testing.T) {
		root := migrationRepository(t)
		configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
		config := mustReadFile(t, configPath)
		config = bytes.Replace(config, []byte("schema_version: 1"), []byte("schema_version: 99"), 1)
		config = bytes.Replace(config, []byte("plan_catalog_mode: mixed\n    plan_catalog_cutover_revision: "+cutoverRevisionFromConfig(t, config)), []byte("plan_catalog_mode: legacy"), 1)
		if err := os.WriteFile(configPath, config, 0o644); err != nil {
			t.Fatal(err)
		}
		runGitTest(t, root, "add", state.ConfigPath)
		runGitTest(t, root, "commit", "-m", "invalid legacy baseline config")
		baseline := runGitTest(t, root, "rev-parse", "HEAD")
		config = bytes.Replace(config, []byte("schema_version: 99"), []byte("schema_version: 1"), 1)
		config = bytes.Replace(config, []byte("plan_catalog_mode: legacy"), []byte("plan_catalog_mode: mixed\n    plan_catalog_cutover_revision: "+baseline), 1)
		if err := os.WriteFile(configPath, config, 0o644); err != nil {
			t.Fatal(err)
		}
		runGitTest(t, root, "add", state.ConfigPath)
		runGitTest(t, root, "commit", "-m", "point at invalid baseline config")
		snapshot, err := LoadRepositorySnapshot(root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemExists(snapshot.Problems, "mixed_cutover_baseline_config_invalid", SeverityError) {
			t.Fatalf("invalid baseline config problems=%#v", snapshot.Problems)
		}
	})
}

func TestRootPlansReadmeNeverQualifiesAsFamilyRecord(t *testing.T) {
	root := migrationRepository(t)
	path := filepath.Join(root, "docs", "repo", "plans", "readme.md")
	if err := os.WriteFile(path, []byte("Last updated: 2026-07-31T12:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: draft\n\n# Root Plan Index\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{"docs/repo/plans/README.md", "docs/repo/plans/readme.md"} {
		if kind, formal := FormalPathKind(root, relative); formal || kind != "" {
			t.Fatalf("root plans README %s classified as %q formal=%t", relative, kind, formal)
		}
	}
	candidates, err := Enumerate(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, candidate := range candidates {
		if strings.EqualFold(candidate.Path, "docs/repo/plans/README.md") {
			t.Fatalf("root plans README was enumerated: %#v", candidate)
		}
	}
}

func TestInventoryCapturesRevisionCoverageUnpairedEvidenceBranchTouchesAndDirtyOverlap(t *testing.T) {
	root := migrationRepository(t)
	alpha := "docs/repo/plans/alpha/alpha_discovery_doc.md"
	beta := "docs/repo/plans/beta/beta_implementation_doc.md"
	runGitTest(t, root, "checkout", "-b", "feature/alpha")
	appendFile(t, filepath.Join(root, filepath.FromSlash(alpha)), "\nFeature branch edit.\n")
	runGitTest(t, root, "add", alpha)
	runGitTest(t, root, "commit", "-m", "touch alpha on feature")
	runGitTest(t, root, "checkout", "main")
	runGitTest(t, root, "update-index", "--chmod=+x", beta)

	inventory, err := BuildInventory(root, time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if inventory.SourceRevision == "" || inventory.Coverage.FormalRecords != 2 || inventory.Coverage.LegacyRecords != 2 || inventory.Coverage.UnpairedLegacyEvidence != 1 {
		t.Fatalf("coverage=%#v revision=%q", inventory.Coverage, inventory.SourceRevision)
	}
	records := inventoryByPath(inventory)
	if strings.Join(records[alpha].ActiveBranchTouch, ",") != "refs/heads/feature/alpha" {
		t.Fatalf("alpha branch touches=%#v", records[alpha].ActiveBranchTouch)
	}
	if !records[beta].DirtyOverlap {
		t.Fatalf("beta dirty overlap not captured: %#v", records[beta])
	}
}

func TestInventoryRetainsReadableMalformedFormalRecords(t *testing.T) {
	root := migrationRepository(t)
	malformed := "docs/repo/plans/malformed/malformed_discovery_doc.md"
	path := filepath.Join(root, filepath.FromSlash(malformed))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	data := []byte("readable but malformed planning authority\n")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	register := filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))
	appendFile(t, register, "\n## PR-003 - Malformed Discovery Evidence\n\nType: discovery-plan\nPurpose: Preserve evidence for a malformed formal plan.\nStatus: draft\nOwner / repository: Example\nCanonical docs: "+malformed+"\nCreated: 2026-07-31\nLast updated: 2026-07-31T10:00:00Z (UTC)\n")
	runGitTest(t, root, "add", malformed, LegacyRegisterPath)
	runGitTest(t, root, "commit", "-m", "add malformed plan evidence")

	inventory, err := BuildInventory(root, time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	records := inventoryByPath(inventory)
	record, exists := records[malformed]
	if !exists || inventory.Coverage.FormalRecords != 3 || inventory.Coverage.InvalidRecords != 1 {
		t.Fatalf("coverage=%#v record=%#v exists=%t", inventory.Coverage, record, exists)
	}
	if record.ParseStatus != "invalid" || record.MetadataCoverage != "invalid" || record.SourceSHA256 != sha256Text(data) || !stringSliceContains(record.ProblemCodes, "header_missing") {
		t.Fatalf("malformed inventory record=%#v", record)
	}
	if len(record.LegacyEvidence) != 1 || record.LegacyEvidence[0].ID != "PR-003" || legacyEntryExists(inventory.UnpairedLegacy, "PR-003") {
		t.Fatalf("malformed legacy evidence record=%#v unpaired=%#v", record, inventory.UnpairedLegacy)
	}
}

func TestActiveBranchTouchesIncludeDeletionsAndBothRenamePaths(t *testing.T) {
	alpha := "docs/repo/plans/alpha/alpha_discovery_doc.md"

	t.Run("deletion", func(t *testing.T) {
		root := migrationRepository(t)
		runGitTest(t, root, "checkout", "-b", "feature/delete-alpha")
		runGitTest(t, root, "rm", alpha)
		runGitTest(t, root, "commit", "-m", "delete alpha")
		runGitTest(t, root, "checkout", "main")
		touches, problems := activeBranchTouches(root)
		if len(problems) != 0 || strings.Join(touches[alpha], ",") != "refs/heads/feature/delete-alpha" {
			t.Fatalf("touches=%#v problems=%#v", touches, problems)
		}
		plan, err := BuildMigrationPlan(root, migrationLedger(t, root))
		if err != nil || len(plan.FilePlan.Actions) != 1 || !skipExists(plan.Skips, "active_branch_ownership") {
			t.Fatalf("plan=%#v skips=%#v err=%v", plan.FilePlan, plan.Skips, err)
		}
	})

	t.Run("rename", func(t *testing.T) {
		root := migrationRepository(t)
		renamed := "docs/repo/plans/alpha/renamed-alpha_discovery_doc.md"
		runGitTest(t, root, "checkout", "-b", "feature/rename-alpha")
		runGitTest(t, root, "mv", alpha, renamed)
		runGitTest(t, root, "commit", "-m", "rename alpha")
		runGitTest(t, root, "checkout", "main")
		touches, problems := activeBranchTouches(root)
		if len(problems) != 0 || strings.Join(touches[alpha], ",") != "refs/heads/feature/rename-alpha" || strings.Join(touches[renamed], ",") != "refs/heads/feature/rename-alpha" {
			t.Fatalf("touches=%#v problems=%#v", touches, problems)
		}
		plan, err := BuildMigrationPlan(root, migrationLedger(t, root))
		if err != nil || len(plan.FilePlan.Actions) != 1 || !skipExists(plan.Skips, "active_branch_ownership") {
			t.Fatalf("plan=%#v skips=%#v err=%v", plan.FilePlan, plan.Skips, err)
		}
	})
}

func TestMigrationDryRunApplyChronologyIdempotencyAndCanonicalCoverage(t *testing.T) {
	root := migrationRepository(t)
	setCatalogModeAndCommit(t, root, "canonical")
	ledger := migrationLedger(t, root)
	before := migrationPlanBytes(t, root)
	beforeSnapshot, err := LoadRepositorySnapshot(root)
	if err != nil || !HasErrors(beforeSnapshot.Problems) {
		t.Fatalf("canonical mode should reject legacy coverage: err=%v problems=%#v", err, beforeSnapshot.Problems)
	}

	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	preview, err := ExecuteMigration(plan, MigrationApplyOptions{DryRun: true, Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)})
	if err != nil || preview.Result.Status != reconcile.StatusPlanned || len(preview.Result.Changes) != 2 {
		t.Fatalf("preview=%#v err=%v", preview, err)
	}
	assertPlanBytes(t, root, before)

	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 7, 31, 12, 1, 0, 0, time.UTC)})
	if err != nil || outcome.Result.Status != reconcile.StatusSucceeded || len(outcome.Result.Changes) != 2 {
		t.Fatalf("outcome=%#v err=%v", outcome, err)
	}
	for path, original := range before {
		after := mustReadFile(t, filepath.Join(root, filepath.FromSlash(path)))
		if headerLine(original, "Last updated:") != headerLine(after, "Last updated:") {
			t.Fatalf("content chronology changed for %s", path)
		}
		if !bytes.Contains(after, []byte(MetadataBeginMarker)) {
			t.Fatalf("metadata missing after migration for %s", path)
		}
		t.Logf("simulation path=%s before_sha256=%s after_sha256=%s preserved_last_updated=%q", path, byteDigest(original), byteDigest(after), headerLine(after, "Last updated:"))
	}
	afterSnapshot, err := LoadRepositorySnapshot(root)
	if err != nil || HasErrors(afterSnapshot.Problems) {
		t.Fatalf("canonical coverage did not close: err=%v problems=%#v", err, afterSnapshot.Problems)
	}

	secondPlan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	second, err := ExecuteMigration(secondPlan, MigrationApplyOptions{Now: time.Date(2026, 7, 31, 12, 2, 0, 0, time.UTC)})
	if err != nil || second.Result.Status != reconcile.StatusSucceeded || len(second.Result.Changes) != 0 || HasMaterialMigrationSkips(second.Skips) {
		t.Fatalf("second apply was not idempotent: outcome=%#v err=%v", second, err)
	}
	t.Logf("simulation second_run_status=%s changes=%d skips=%d", second.Result.Status, len(second.Result.Changes), len(second.Skips))
}

func TestReviewedLedgerRoundTripsThroughDurableSchema(t *testing.T) {
	root := migrationRepository(t)
	ledger := migrationLedger(t, root)
	data, err := state.EncodeYAML(ledger)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadMigrationLedger(data)
	if err != nil || loaded.RepositoryID != ledger.RepositoryID || len(loaded.Records) != 2 {
		t.Fatalf("loaded=%#v err=%v\n%s", loaded, err, data)
	}
	invalid := bytes.Replace(data, []byte("schema_version: 1"), []byte("schema_version: 9"), 1)
	if _, err := LoadMigrationLedger(invalid); err == nil || !strings.Contains(err.Error(), "invalid_ledger") {
		t.Fatalf("invalid ledger error=%v", err)
	}
}

func TestMigrationChangedHashAndDirtyOverlapProduceZeroWrites(t *testing.T) {
	t.Run("changed hash", func(t *testing.T) {
		root := migrationRepository(t)
		ledger := migrationLedger(t, root)
		alphaPath := filepath.Join(root, "docs/repo/plans/alpha/alpha_discovery_doc.md")
		appendFile(t, alphaPath, "\nConcurrent edit.\n")
		before := mustReadFile(t, alphaPath)
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil {
			t.Fatal(err)
		}
		if len(plan.FilePlan.Actions) != 1 || !skipExists(plan.Skips, "source_mismatch") {
			t.Fatalf("plan actions=%#v skips=%#v", plan.FilePlan.Actions, plan.Skips)
		}
		if string(mustReadFile(t, alphaPath)) != string(before) {
			t.Fatal("planning overwrote changed source")
		}
	})

	t.Run("dirty overlap", func(t *testing.T) {
		root := migrationRepository(t)
		ledger := migrationLedger(t, root)
		alpha := "docs/repo/plans/alpha/alpha_discovery_doc.md"
		runGitTest(t, root, "update-index", "--chmod=+x", alpha)
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil {
			t.Fatal(err)
		}
		if len(plan.FilePlan.Actions) != 1 || !skipExists(plan.Skips, "dirty_plan_overlap") {
			t.Fatalf("plan actions=%#v skips=%#v", plan.FilePlan.Actions, plan.Skips)
		}
	})

	t.Run("unrelated dirty file", func(t *testing.T) {
		root := migrationRepository(t)
		ledger := migrationLedger(t, root)
		if err := os.WriteFile(filepath.Join(root, "unrelated.txt"), []byte("unrelated\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil || len(plan.FilePlan.Actions) != 2 || len(plan.Skips) != 0 {
			t.Fatalf("unrelated dirty state affected migration: actions=%#v skips=%#v err=%v", plan.FilePlan.Actions, plan.Skips, err)
		}
	})
}

func TestMigrationAcceptsCleanCheckoutLineEndingConversion(t *testing.T) {
	root := migrationRepository(t)
	alpha := "docs/repo/plans/alpha/alpha_discovery_doc.md"
	runGitTest(t, root, "config", "core.autocrlf", "true")
	runGitTest(t, root, "config", "core.eol", "crlf")
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(alpha))); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "checkout", "--", alpha)
	if data := mustReadFile(t, filepath.Join(root, filepath.FromSlash(alpha))); !bytes.Contains(data, []byte("\r\n")) {
		t.Skip("Git did not apply the configured checkout line-ending conversion")
	}
	if status := runGitTest(t, root, "status", "--porcelain=v1", "--", alpha); status != "" {
		t.Fatalf("line-ending conversion unexpectedly dirtied the checkout: %q", status)
	}
	inventory, err := BuildInventory(root, time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	record := inventoryRecordByPath(t, inventory, alpha)
	sourceBytes, err := gitBytes(root, "show", inventory.SourceRevision+":"+alpha)
	if err != nil {
		t.Fatal(err)
	}
	if record.SourceSHA256 != sha256Text(sourceBytes) || record.SourceSHA256 == sha256Text(mustReadFile(t, filepath.Join(root, filepath.FromSlash(alpha)))) {
		t.Fatalf("inventory did not retain the committed blob identity across checkout conversion: %#v", record)
	}

	plan, err := BuildMigrationPlan(root, migrationLedger(t, root))
	if err != nil || len(plan.FilePlan.Actions) != 2 || len(plan.Skips) != 0 {
		t.Fatalf("clean converted checkout was not migratable: actions=%#v skips=%#v err=%v", plan.FilePlan.Actions, plan.Skips, err)
	}
}

func TestMigrationActiveBranchSkipAndTransactionalRollback(t *testing.T) {
	t.Run("active branch", func(t *testing.T) {
		root := migrationRepository(t)
		ledger := migrationLedger(t, root)
		alpha := "docs/repo/plans/alpha/alpha_discovery_doc.md"
		runGitTest(t, root, "checkout", "-b", "feature/alpha")
		appendFile(t, filepath.Join(root, filepath.FromSlash(alpha)), "\nOwned elsewhere.\n")
		runGitTest(t, root, "add", alpha)
		runGitTest(t, root, "commit", "-m", "own alpha")
		runGitTest(t, root, "checkout", "main")
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil {
			t.Fatal(err)
		}
		if len(plan.FilePlan.Actions) != 1 || !skipExists(plan.Skips, "active_branch_ownership") {
			t.Fatalf("plan actions=%#v skips=%#v", plan.FilePlan.Actions, plan.Skips)
		}
	})

	t.Run("rollback", func(t *testing.T) {
		root := migrationRepository(t)
		ledger := migrationLedger(t, root)
		before := migrationPlanBytes(t, root)
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil {
			t.Fatal(err)
		}
		outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
			if phase == "commit:docs/repo/plans/beta/beta_implementation_doc.md" {
				return errors.New("injected migration failure")
			}
			return nil
		}})
		if err != nil || outcome.Result.Status != reconcile.StatusRolledBack || !outcome.Result.Rollback.Succeeded {
			t.Fatalf("rollback outcome=%#v err=%v", outcome, err)
		}
		assertPlanBytes(t, root, before)
	})

	t.Run("rollback failure", func(t *testing.T) {
		root := migrationRepository(t)
		ledger := migrationLedger(t, root)
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil {
			t.Fatal(err)
		}
		outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
			if phase != "commit:docs/repo/plans/beta/beta_implementation_doc.md" {
				return nil
			}
			backup := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.FilePlan.ID, "backup", "docs", "repo", "plans", "alpha", "alpha_discovery_doc.md")
			if err := os.Remove(backup); err != nil {
				return err
			}
			return errors.New("injected failure after backup loss")
		}})
		if err != nil || outcome.Result.Status != reconcile.StatusRecoveryRequired || !reconcileBlockerExists(outcome.Result, "recovery_required") {
			t.Fatalf("rollback failure outcome=%#v err=%v", outcome, err)
		}
	})

	t.Run("partial commit restores the current action", func(t *testing.T) {
		root := migrationRepository(t)
		ledger := migrationLedger(t, root)
		before := migrationPlanBytes(t, root)
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil {
			t.Fatal(err)
		}
		outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
			if phase != "commit:docs/repo/plans/beta/beta_implementation_doc.md" {
				return nil
			}
			stage := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.FilePlan.ID, "stage", "docs", "repo", "plans", "beta", "beta_implementation_doc.md")
			return os.Remove(stage)
		}})
		if err != nil || outcome.Result.Status != reconcile.StatusRolledBack || !outcome.Result.Rollback.Succeeded {
			t.Fatalf("partial rollback outcome=%#v err=%v", outcome, err)
		}
		assertPlanBytes(t, root, before)
	})

	t.Run("rollback preserves a concurrent edit to an earlier action", func(t *testing.T) {
		root := migrationRepository(t)
		ledger := migrationLedger(t, root)
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil {
			t.Fatal(err)
		}
		alpha := filepath.Join(root, "docs", "repo", "plans", "alpha", "alpha_discovery_doc.md")
		concurrent := []byte("concurrent branch-owner bytes\n")
		outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
			if phase != "commit:docs/repo/plans/beta/beta_implementation_doc.md" {
				return nil
			}
			if err := os.WriteFile(alpha, concurrent, 0o644); err != nil {
				return err
			}
			return errors.New("injected failure after concurrent alpha edit")
		}})
		if err != nil || outcome.Result.Status != reconcile.StatusRecoveryRequired || !reconcileBlockerExists(outcome.Result, "recovery_required") {
			t.Fatalf("concurrent rollback outcome=%#v err=%v", outcome, err)
		}
		if actual := mustReadFile(t, alpha); !bytes.Equal(actual, concurrent) {
			t.Fatalf("rollback overwrote concurrent bytes: %q", actual)
		}
	})

	t.Run("rollback restore never replaces a target recreated after quarantine", func(t *testing.T) {
		root := migrationRepository(t)
		ledger := migrationLedger(t, root)
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil {
			t.Fatal(err)
		}
		alphaRelative := "docs/repo/plans/alpha/alpha_discovery_doc.md"
		alpha := filepath.Join(root, filepath.FromSlash(alphaRelative))
		concurrent := []byte("concurrent bytes recreated during rollback\n")
		outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
			switch phase {
			case "commit:docs/repo/plans/beta/beta_implementation_doc.md":
				return errors.New("trigger rollback after alpha commit")
			case "rollback-quarantined:" + alphaRelative:
				return os.WriteFile(alpha, concurrent, 0o644)
			default:
				return nil
			}
		}})
		if err != nil || outcome.Result.Status != reconcile.StatusRecoveryRequired || !reconcileBlockerExists(outcome.Result, "recovery_required") {
			t.Fatalf("rollback no-replace outcome=%#v err=%v", outcome, err)
		}
		if actual := mustReadFile(t, alpha); !bytes.Equal(actual, concurrent) {
			t.Fatalf("rollback restore overwrote recreated target: %q", actual)
		}
	})
}

func TestMigrationSemanticAmbiguityAndInvalidPlacementAreStableSkips(t *testing.T) {
	root := migrationRepository(t)
	ledger := migrationLedger(t, root)
	ledger.Records[0].Ambiguity = []string{"register title conflicts with plan evidence"}
	ledger.Records[1].Path = "docs/repo/plans/plan-register.md"
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.FilePlan.Actions) != 0 || !skipExists(plan.Skips, "ambiguous_legacy_evidence") || !skipExists(plan.Skips, "invalid_target_placement") {
		t.Fatalf("actions=%#v skips=%#v", plan.FilePlan.Actions, plan.Skips)
	}
}

func TestMigrationRejectsSymlinkTraversalDuringDryRunPlanning(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("ordinary Windows test users may not have symlink privilege")
	}
	root := migrationRepository(t)
	outside := t.TempDir()
	outsidePlan := filepath.Join(outside, "escaped_discovery_doc.md")
	if err := os.WriteFile(outsidePlan, []byte("outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "docs", "repo", "plans", "escaped")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}
	ledger := migrationLedger(t, root)
	ledger.Records = ledger.Records[:1]
	ledger.Records[0].Path = "docs/repo/plans/escaped/escaped_discovery_doc.md"
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.FilePlan.Actions) != 0 || !skipExists(plan.Skips, "unsafe_target") {
		t.Fatalf("symlink traversal actions=%#v skips=%#v", plan.FilePlan.Actions, plan.Skips)
	}
}

func TestMigrationCommitPreconditionPreservesConcurrentEdit(t *testing.T) {
	root := migrationRepository(t)
	ledger := migrationLedger(t, root)
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	alpha := filepath.Join(root, "docs/repo/plans/alpha/alpha_discovery_doc.md")
	concurrent := []byte("concurrent owner bytes\n")
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
		if phase == "commit:docs/repo/plans/alpha/alpha_discovery_doc.md" {
			return os.WriteFile(alpha, concurrent, 0o644)
		}
		return nil
	}})
	if err != nil || outcome.Result.Status != reconcile.StatusRolledBack || string(mustReadFile(t, alpha)) != string(concurrent) {
		t.Fatalf("concurrent edit outcome=%#v bytes=%q err=%v", outcome, mustReadFile(t, alpha), err)
	}
}

func TestMigrationBackedUpPreconditionClosesCheckRenameRace(t *testing.T) {
	root := migrationRepository(t)
	ledger := migrationLedger(t, root)
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	alpha := filepath.Join(root, "docs", "repo", "plans", "alpha", "alpha_discovery_doc.md")
	concurrent := []byte("concurrent bytes written after the first precondition check\n")
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
		if phase == "pre-commit:docs/repo/plans/alpha/alpha_discovery_doc.md" {
			return os.WriteFile(alpha, concurrent, 0o644)
		}
		return nil
	}})
	if err != nil || outcome.Result.Status != reconcile.StatusRolledBack || !outcome.Result.Rollback.Succeeded {
		t.Fatalf("check/rename race outcome=%#v err=%v", outcome, err)
	}
	if actual := mustReadFile(t, alpha); !bytes.Equal(actual, concurrent) {
		t.Fatalf("check/rename race overwrote concurrent bytes: %q", actual)
	}
}

func TestMigrationInstallNeverReplacesARecreatedTarget(t *testing.T) {
	root := migrationRepository(t)
	ledger := migrationLedger(t, root)
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	alphaRelative := "docs/repo/plans/alpha/alpha_discovery_doc.md"
	alpha := filepath.Join(root, filepath.FromSlash(alphaRelative))
	original := mustReadFile(t, alpha)
	concurrent := []byte("concurrent target recreated after backup validation\n")
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
		if phase == "backed-up:"+alphaRelative {
			return os.WriteFile(alpha, concurrent, 0o644)
		}
		return nil
	}})
	if err != nil || outcome.Result.Status != reconcile.StatusRecoveryRequired || !reconcileBlockerExists(outcome.Result, "recovery_required") {
		t.Fatalf("no-replace outcome=%#v err=%v", outcome, err)
	}
	if actual := mustReadFile(t, alpha); !bytes.Equal(actual, concurrent) {
		t.Fatalf("installation overwrote recreated target: %q", actual)
	}
	backup := filepath.Join(root, ".codeheart", "local", "kit-transactions", plan.FilePlan.ID, "backup", filepath.FromSlash(alphaRelative))
	if actual := mustReadFile(t, backup); !bytes.Equal(actual, original) {
		t.Fatalf("recovery backup changed: %q", actual)
	}
}

func TestCanonicalMigrationBlocksProjectedIncompleteCoverage(t *testing.T) {
	root := migrationRepository(t)
	setCatalogModeAndCommit(t, root, "canonical")
	ledger := migrationLedger(t, root)
	ledger.Records = ledger.Records[:1]
	before := migrationPlanBytes(t, root)
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.FilePlan.Actions) != 1 {
		t.Fatalf("expected one projected action, got %#v", plan.FilePlan.Actions)
	}
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)})
	if err != nil || outcome.Result.Status != reconcile.StatusBlocked || !reconcileBlockerExists(outcome.Result, "incomplete_coverage") {
		t.Fatalf("incomplete canonical outcome=%#v err=%v", outcome, err)
	}
	assertPlanBytes(t, root, before)
}

func migrationRepository(t *testing.T) string {
	t.Helper()
	source := filepath.Join("..", "..", "tests", "fixtures", "plans", "migration-repository")
	root := t.TempDir()
	err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil || relative == "." {
			return err
		}
		target := filepath.Join(root, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "init")
	runGitTest(t, root, "checkout", "-b", "main")
	runGitTest(t, root, "config", "user.email", "test@example.invalid")
	runGitTest(t, root, "config", "user.name", "Plan Catalog Test")
	runGitTest(t, root, "add", ".")
	runGitTest(t, root, "commit", "-m", "fixture baseline")
	adoptMixedCatalogTest(t, root)
	return root
}

func adoptMixedCatalogTest(t *testing.T, root string) {
	t.Helper()
	baseline := runGitTest(t, root, "rev-parse", "HEAD")
	path := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	data := mustReadFile(t, path)
	data = bytes.Replace(data, []byte("plan_catalog_mode: legacy"), []byte("plan_catalog_mode: mixed\n    plan_catalog_cutover_revision: "+baseline), 1)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", state.ConfigPath)
	runGitTest(t, root, "commit", "-m", "adopt mixed plan catalog")
}

func migrationLedger(t *testing.T, root string) MigrationLedger {
	t.Helper()
	revision := runGitTest(t, root, "rev-parse", "HEAD")
	record := func(path, id string, kind Kind, purpose, alias string) MigrationRecord {
		data, err := gitBytes(root, "show", revision+":"+path)
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(data)
		return MigrationRecord{
			Path: path, SourceRevision: revision, SourceSHA256: hex.EncodeToString(digest[:]),
			Decision:   MigrationDecision{ID: id, Kind: kind, Purpose: purpose, FirstCataloged: "2026-07-31T12:00:00Z", CatalogMetadataUpdated: "2026-07-31T12:00:00Z"},
			Confidence: "high", Ambiguity: []string{}, Evidence: []string{"reviewed fixture plan and register"}, LegacyAliases: []string{alias}, Conflicts: []string{}, Deferred: false, BranchOwner: "none",
		}
	}
	return MigrationLedger{
		SchemaVersion: 1, RepositoryID: "example", InventoryRevision: revision, ReviewedAt: "2026-07-31T12:00:00Z",
		Records: []MigrationRecord{
			record("docs/repo/plans/alpha/alpha_discovery_doc.md", "example.discovery.alpha", KindDiscovery, "Discover alpha", "PR-001"),
			record("docs/repo/plans/beta/beta_implementation_doc.md", "example.implementation.beta", KindImplementation, "Implement beta", "PR-002"),
		},
	}
}

func metadataForTest(id string, kind Kind, purpose string) Metadata {
	return Metadata{SchemaVersion: 1, ID: id, Kind: kind, Purpose: purpose, FirstCataloged: "2026-07-31T12:00:00Z", CatalogMetadataUpdated: "2026-07-31T12:00:00Z"}
}

func setCatalogModeAndCommit(t *testing.T, root, mode string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	data := mustReadFile(t, path)
	data = bytes.Replace(data, []byte("plan_catalog_mode: mixed"), []byte("plan_catalog_mode: "+mode), 1)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", state.ConfigPath)
	runGitTest(t, root, "commit", "-m", "set catalog mode "+mode)
}

func migrationPlanBytes(t *testing.T, root string) map[string][]byte {
	t.Helper()
	result := map[string][]byte{}
	for _, path := range []string{"docs/repo/plans/alpha/alpha_discovery_doc.md", "docs/repo/plans/beta/beta_implementation_doc.md"} {
		result[path] = mustReadFile(t, filepath.Join(root, filepath.FromSlash(path)))
	}
	return result
}

func assertPlanBytes(t *testing.T, root string, expected map[string][]byte) {
	t.Helper()
	for path, data := range expected {
		if actual := mustReadFile(t, filepath.Join(root, filepath.FromSlash(path))); !bytes.Equal(actual, data) {
			t.Fatalf("%s changed\nwant:\n%s\ngot:\n%s", path, data, actual)
		}
	}
}

func inventoryByPath(inventory Inventory) map[string]InventoryRecord {
	result := map[string]InventoryRecord{}
	for _, record := range inventory.Records {
		result[record.Path] = record
	}
	return result
}

func stringSliceContains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func legacyEntryExists(entries []LegacyEntry, id string) bool {
	for _, entry := range entries {
		if entry.ID == id {
			return true
		}
	}
	return false
}

func skipExists(skips []MigrationSkip, code string) bool {
	for _, skip := range skips {
		if skip.Code == code {
			return true
		}
	}
	return false
}

func problemForPathExists(problems []Problem, code, path string, severity Severity) bool {
	for _, problem := range problems {
		if problem.Code == code && problem.Path == path && problem.Severity == severity {
			return true
		}
	}
	return false
}

func cutoverRevisionFromConfig(t *testing.T, data []byte) string {
	t.Helper()
	for _, line := range strings.Split(string(data), "\n") {
		if value := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "plan_catalog_cutover_revision:")); value != strings.TrimSpace(line) && value != "" {
			return value
		}
	}
	t.Fatal("plan_catalog_cutover_revision missing from fixture config")
	return ""
}

func reconcileBlockerExists(result reconcile.Result, code string) bool {
	for _, blocker := range result.Blockers {
		if blocker.Code == code {
			return true
		}
	}
	return false
}

func runGitTest(t *testing.T, root string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func inventoryRecordByPath(t *testing.T, inventory Inventory, path string) InventoryRecord {
	t.Helper()
	for _, record := range inventory.Records {
		if record.Path == path {
			return record
		}
	}
	t.Fatalf("inventory record missing for %s", path)
	return InventoryRecord{}
}

func appendFile(t *testing.T, path, content string) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		t.Fatal(err)
	}
}

func headerLine(data []byte, prefix string) string {
	for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	return ""
}

func byteDigest(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
