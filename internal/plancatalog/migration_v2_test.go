package plancatalog

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

func TestV2MigrationDirectLegacyToCanonicalIsBoundReadinessOnlyAndIdempotent(t *testing.T) {
	root := legacyV2MigrationRepository(t)
	ledger := v2MigrationLedger(t, root, ModeCanonical, nil, nil)
	protected := map[string][]byte{
		state.ConfigPath:   mustReadFile(t, filepath.Join(root, filepath.FromSlash(state.ConfigPath))),
		LegacyRegisterPath: mustReadFile(t, filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))),
	}
	headers := map[string][]byte{}
	for _, item := range ledger.Records {
		headers[item.CurrentPath] = mustReadFile(t, filepath.Join(root, filepath.FromSlash(item.CurrentPath)))
	}

	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil || len(plan.FilePlan.Actions) != 2 || len(plan.FilePlan.Blockers) != 0 || !plan.Projection.Ready || plan.Projection.Gaps != 0 {
		t.Fatalf("v2 plan=%#v projection=%#v problems=%#v skips=%#v err=%v", plan.FilePlan, plan.Projection, plan.Problems, plan.Skips, err)
	}
	preview, err := ExecuteMigration(plan, MigrationApplyOptions{DryRun: true, Now: time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)})
	if err != nil || preview.SchemaVersion != 2 || preview.Result.Status != reconcile.StatusPlanned || preview.ActivationPerformed == nil || *preview.ActivationPerformed || preview.Projection == nil || !preview.Projection.Ready {
		t.Fatalf("v2 preview=%#v err=%v", preview, err)
	}
	for path, expected := range protected {
		if actual := mustReadFile(t, filepath.Join(root, filepath.FromSlash(path))); !bytes.Equal(actual, expected) {
			t.Fatalf("dry-run changed protected %s", path)
		}
	}

	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 8, 6, 12, 1, 0, 0, time.UTC)})
	if err != nil || outcome.Result.Status != reconcile.StatusSucceeded || len(outcome.Result.Changes) != 2 {
		t.Fatalf("v2 outcome=%#v err=%v", outcome, err)
	}
	for path, before := range headers {
		after := mustReadFile(t, filepath.Join(root, filepath.FromSlash(path)))
		if headerLine(before, "Created:") != headerLine(after, "Created:") || headerLine(before, "Last updated:") != headerLine(after, "Last updated:") || headerLine(before, "Status:") != headerLine(after, "Status:") {
			t.Fatalf("migration changed lifecycle chronology for %s", path)
		}
	}
	for path, expected := range protected {
		if actual := mustReadFile(t, filepath.Join(root, filepath.FromSlash(path))); !bytes.Equal(actual, expected) {
			t.Fatalf("migration changed protected %s", path)
		}
	}
	active, err := LoadRepositorySnapshot(root)
	if err != nil || active.Settings.DiscoveryVersion != DiscoveryV1 || active.Settings.Mode != ModeLegacy {
		t.Fatalf("migration silently activated catalog settings: snapshot=%#v err=%v", active, err)
	}
	prospective, err := LoadRepositorySnapshotWithOptions(root, SnapshotOptions{TargetCatalogMode: ModeCanonical})
	if err != nil || HasErrors(prospective.Problems) || !prospective.Complete {
		t.Fatalf("projected canonical readiness did not close: complete=%t problems=%#v err=%v", prospective.Complete, prospective.Problems, err)
	}

	secondPlan, err := BuildMigrationPlan(root, ledger)
	if err != nil || len(secondPlan.FilePlan.Actions) != 0 || len(secondPlan.FilePlan.Blockers) != 0 || HasMaterialMigrationSkips(secondPlan.Skips) {
		t.Fatalf("v2 reapply was not idempotent: plan=%#v blockers=%#v skips=%#v err=%v", secondPlan.FilePlan.Actions, secondPlan.FilePlan.Blockers, secondPlan.Skips, err)
	}
	second, err := ExecuteMigration(secondPlan, MigrationApplyOptions{Now: time.Date(2026, 8, 6, 12, 2, 0, 0, time.UTC)})
	if err != nil || second.Result.Status != reconcile.StatusSucceeded || len(second.Result.Changes) != 0 {
		t.Fatalf("v2 second outcome=%#v err=%v", second, err)
	}
}

func TestV2MigrationFinalStateFastPathRequiresExactReviewedTransformation(t *testing.T) {
	tests := []struct {
		name   string
		tamper func(t *testing.T, root string, ledger *MigrationLedger)
	}{
		{
			name: "record source hash",
			tamper: func(t *testing.T, root string, ledger *MigrationLedger) {
				ledger.Records[0].SourceSHA256 = strings.Repeat("a", 64)
			},
		},
		{
			name: "record source revision",
			tamper: func(t *testing.T, root string, ledger *MigrationLedger) {
				ledger.Records[0].SourceRevision = strings.Repeat("b", 40)
			},
		},
		{
			name: "body bytes with matching metadata",
			tamper: func(t *testing.T, root string, ledger *MigrationLedger) {
				appendFile(t, filepath.Join(root, filepath.FromSlash(ledger.Records[0].TargetPath)), "\nUnreviewed body drift.\n")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := legacyV2MigrationRepository(t)
			ledger := v2MigrationLedger(t, root, ModeCanonical, nil, nil)
			plan, err := BuildMigrationPlan(root, ledger)
			if err != nil || len(plan.FilePlan.Blockers) != 0 {
				t.Fatalf("initial plan blockers=%#v err=%v", plan.FilePlan.Blockers, err)
			}
			outcome, err := ExecuteMigration(plan, MigrationApplyOptions{})
			if err != nil || outcome.Result.Status != reconcile.StatusSucceeded {
				t.Fatalf("initial outcome=%#v err=%v", outcome, err)
			}
			test.tamper(t, root, &ledger)
			reapply, err := BuildMigrationPlan(root, ledger)
			if err != nil {
				t.Fatal(err)
			}
			if reapply.Projection.Ready || len(reapply.FilePlan.Blockers) == 0 || len(reapply.FilePlan.Actions) != 0 {
				t.Fatalf("tampered final state was accepted: projection=%#v actions=%#v blockers=%#v", reapply.Projection, reapply.FilePlan.Actions, reapply.FilePlan.Blockers)
			}
		})
	}
}

func TestV2MigrationRenamesMetadataOnlyCandidateAtomically(t *testing.T) {
	root := legacyV2MigrationRepository(t)
	current := "products/widget/docs/notes/idea.md"
	target := "products/widget/docs/notes/idea_discovery_doc.md"
	writeClassifierFile(t, root, current, canonicalClassifierDocument("example.discovery.idea", KindDiscovery, "Idea"))
	runGitTest(t, root, "add", current)
	runGitTest(t, root, "commit", "-m", "add metadata-only plan")
	ledger := v2MigrationLedger(t, root, ModeCanonical, map[string]string{current: target}, nil)
	before := mustReadFile(t, filepath.Join(root, filepath.FromSlash(current)))
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil || len(plan.FilePlan.Actions) != 4 || len(plan.FilePlan.Blockers) != 0 || !plan.Projection.Ready {
		t.Fatalf("rename plan=%#v projection=%#v problems=%#v err=%v", plan.FilePlan, plan.Projection, plan.Problems, err)
	}
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
		if phase == "commit:"+target {
			return errors.New("injected rename failure")
		}
		return nil
	}})
	if err != nil || outcome.Result.Status != reconcile.StatusRolledBack || !outcome.Result.Rollback.Succeeded {
		t.Fatalf("rename rollback=%#v err=%v", outcome, err)
	}
	if actual := mustReadFile(t, filepath.Join(root, filepath.FromSlash(current))); !bytes.Equal(actual, before) {
		t.Fatalf("rollback changed source bytes")
	}
	if _, statErr := os.Stat(filepath.Join(root, filepath.FromSlash(target))); !os.IsNotExist(statErr) {
		t.Fatalf("rollback left rename target: %v", statErr)
	}

	outcome, err = ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 8, 6, 13, 0, 0, 0, time.UTC)})
	if err != nil || outcome.Result.Status != reconcile.StatusSucceeded {
		t.Fatalf("rename apply=%#v err=%v", outcome, err)
	}
	if _, statErr := os.Stat(filepath.Join(root, filepath.FromSlash(current))); !os.IsNotExist(statErr) {
		t.Fatalf("rename source remains: %v", statErr)
	}
	after := mustReadFile(t, filepath.Join(root, filepath.FromSlash(target)))
	if !bytes.Equal(after, before) {
		t.Fatalf("metadata-only rename changed document bytes")
	}
	secondPlan, err := BuildMigrationPlan(root, ledger)
	if err != nil || len(secondPlan.FilePlan.Actions) != 0 || len(secondPlan.FilePlan.Blockers) != 0 || HasMaterialMigrationSkips(secondPlan.Skips) {
		t.Fatalf("immediate rename reapply was not a safe no-op: actions=%#v blockers=%#v skips=%#v err=%v", secondPlan.FilePlan.Actions, secondPlan.FilePlan.Blockers, secondPlan.Skips, err)
	}
	second, err := ExecuteMigration(secondPlan, MigrationApplyOptions{})
	if err != nil || second.Result.Status != reconcile.StatusSucceeded || len(second.Result.Changes) != 0 {
		t.Fatalf("immediate rename reapply outcome=%#v err=%v", second, err)
	}
}

func TestV2MigrationRevalidatesAuthorityInsideTransaction(t *testing.T) {
	root := legacyV2MigrationRepository(t)
	ledger := v2MigrationLedger(t, root, ModeCanonical, nil, nil)
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil || len(plan.FilePlan.Blockers) != 0 {
		t.Fatalf("plan blockers=%#v err=%v", plan.FilePlan.Blockers, err)
	}
	before := migrationPlanBytes(t, root)
	registerPath := filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Hook: func(phase string) error {
		if phase == "validated" {
			appendFile(t, registerPath, "\nConcurrent frozen-register drift.\n")
		}
		return nil
	}})
	if err != nil || outcome.Result.Status != reconcile.StatusRolledBack || !reconcileBlockerExists(outcome.Result, "migration_authority_drift") {
		t.Fatalf("authority drift outcome=%#v err=%v", outcome, err)
	}
	assertPlanBytes(t, root, before)
}

func TestMigrationLedgerVersionsCannotCrossDiscoveryAuthority(t *testing.T) {
	root := migrationRepository(t)
	configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	config := mustReadFile(t, configPath)
	config = bytes.Replace(config, []byte("    plan_catalog_mode: mixed\n"), []byte("    plan_catalog_mode: mixed\n    plan_catalog_discovery_version: 2\n    plan_catalog_ownership:\n      excluded_roots: []\n"), 1)
	if err := os.WriteFile(configPath, config, 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", state.ConfigPath)
	runGitTest(t, root, "commit", "-m", "activate discovery v2")
	plan, err := BuildMigrationPlan(root, migrationLedger(t, root))
	if err != nil || !reconcilePlanBlockerExists(plan.FilePlan, "migration_ledger_version_incompatible") {
		t.Fatalf("v1 ledger crossed v2 authority: blockers=%#v err=%v", plan.FilePlan.Blockers, err)
	}
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{})
	if err != nil || outcome.Result.Status != reconcile.StatusBlocked {
		t.Fatalf("incompatible ledger outcome=%#v err=%v", outcome, err)
	}
}

func TestV2MigrationProtectsFrozenRegisterEvenWhenItHasPlanMetadata(t *testing.T) {
	root := legacyV2MigrationRepository(t)
	registerPath := filepath.Join(root, filepath.FromSlash(LegacyRegisterPath))
	register := mustReadFile(t, registerPath)
	register = bytes.Replace(register, []byte("Last updated: 2026-07-31T10:00:00Z (UTC)\n"), []byte("Last updated: 2026-07-31T10:00:00Z (UTC)\nCreated: 2026-07-31\nStatus: draft\n"), 1)
	updated, err := InsertMetadata(register, metadataForTest("example.discovery.frozen-register", KindDiscovery, "Hostile register signal"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(registerPath, updated, 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", LegacyRegisterPath)
	runGitTest(t, root, "commit", "-m", "add hostile register metadata")
	target := "docs/archive/frozen-register_discovery_doc.md"
	ledger := v2MigrationLedger(t, root, ModeCanonical, map[string]string{LegacyRegisterPath: target}, nil)
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil || !reconcilePlanBlockerExists(plan.FilePlan, "migration_protected_path") {
		t.Fatalf("protected register blockers=%#v problems=%#v err=%v", plan.FilePlan.Blockers, plan.Problems, err)
	}
	before := mustReadFile(t, registerPath)
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{})
	if err != nil || outcome.Result.Status != reconcile.StatusBlocked || !bytes.Equal(before, mustReadFile(t, registerPath)) {
		t.Fatalf("protected register outcome=%#v err=%v", outcome, err)
	}
}

func TestV2MigrationPreservesExecutableMarkdownMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Git executable mode is not represented by ordinary Windows worktree permissions")
	}
	root := legacyV2MigrationRepository(t)
	alpha := "docs/repo/plans/alpha/alpha_discovery_doc.md"
	alphaPath := filepath.Join(root, filepath.FromSlash(alpha))
	if err := os.Chmod(alphaPath, 0o755); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", alpha)
	runGitTest(t, root, "commit", "-m", "make plan executable")
	ledger := v2MigrationLedger(t, root, ModeCanonical, nil, nil)
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil || len(plan.FilePlan.Blockers) != 0 {
		t.Fatalf("executable plan blockers=%#v err=%v", plan.FilePlan.Blockers, err)
	}
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{})
	if err != nil || outcome.Result.Status != reconcile.StatusSucceeded {
		t.Fatalf("executable migration=%#v err=%v", outcome, err)
	}
	info, err := os.Stat(alphaPath)
	if err != nil || info.Mode().Perm() != 0o755 {
		t.Fatalf("executable mode changed: mode=%v err=%v", info.Mode().Perm(), err)
	}
}

func TestGitModeMatchingUsesExecutableClassAcrossPlatforms(t *testing.T) {
	tests := []struct {
		mode        GitMode
		permissions os.FileMode
		windows     bool
		want        bool
	}{
		{mode: GitModeRegular, permissions: 0o600, want: true},
		{mode: GitModeRegular, permissions: 0o666, want: true},
		{mode: GitModeRegular, permissions: 0o755, want: false},
		{mode: GitModeExecutable, permissions: 0o700, want: true},
		{mode: GitModeExecutable, permissions: 0o644, want: false},
		{mode: GitModeRegular, permissions: 0o666, windows: true, want: true},
		{mode: GitModeExecutable, permissions: 0o666, windows: true, want: true},
		{mode: GitModeGitlink, permissions: 0o777, windows: true, want: false},
	}
	for _, test := range tests {
		if got := gitModeMatchesPermissions(test.mode, test.permissions, test.windows); got != test.want {
			t.Fatalf("mode=%s permissions=%#o windows=%t got=%t want=%t", test.mode, test.permissions, test.windows, got, test.want)
		}
	}
}

func TestV2MigrationBlocksTargetCollisionStaleCoverageAndDirtySourceWithoutWrites(t *testing.T) {
	t.Run("target collision", func(t *testing.T) {
		root := legacyV2MigrationRepository(t)
		current := "products/widget/docs/notes/collision.md"
		target := "products/widget/docs/notes/collision_discovery_doc.md"
		writeClassifierFile(t, root, current, canonicalClassifierDocument("example.discovery.collision", KindDiscovery, "Collision"))
		runGitTest(t, root, "add", current)
		runGitTest(t, root, "commit", "-m", "add collision source")
		ledger := v2MigrationLedger(t, root, ModeCanonical, map[string]string{current: target}, nil)
		writeClassifierFile(t, root, target, "occupied\n")
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil || !reconcilePlanBlockerExists(plan.FilePlan, "target_precondition_mismatch") {
			t.Fatalf("collision blockers=%#v problems=%#v err=%v", plan.FilePlan.Blockers, plan.Problems, err)
		}
		before := mustReadFile(t, filepath.Join(root, filepath.FromSlash(current)))
		outcome, err := ExecuteMigration(plan, MigrationApplyOptions{})
		if err != nil || outcome.Result.Status != reconcile.StatusBlocked || !bytes.Equal(before, mustReadFile(t, filepath.Join(root, filepath.FromSlash(current)))) {
			t.Fatalf("collision execution=%#v err=%v", outcome, err)
		}
	})

	t.Run("incomplete ledger", func(t *testing.T) {
		root := legacyV2MigrationRepository(t)
		ledger := v2MigrationLedger(t, root, ModeCanonical, nil, nil)
		ledger.Records = ledger.Records[:1]
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil || !reconcilePlanBlockerExists(plan.FilePlan, "migration_candidate_unreviewed") || plan.Projection.Gaps != 1 || plan.Projection.Ready {
			t.Fatalf("incomplete projection=%#v blockers=%#v err=%v", plan.Projection, plan.FilePlan.Blockers, err)
		}
	})

	t.Run("dirty source", func(t *testing.T) {
		root := legacyV2MigrationRepository(t)
		ledger := v2MigrationLedger(t, root, ModeCanonical, nil, nil)
		path := ledger.Records[0].CurrentPath
		appendFile(t, filepath.Join(root, filepath.FromSlash(path)), "\nconcurrent edit\n")
		before := mustReadFile(t, filepath.Join(root, filepath.FromSlash(path)))
		plan, err := BuildMigrationPlan(root, ledger)
		if err != nil || (!reconcilePlanBlockerExists(plan.FilePlan, "source_mismatch") && !reconcilePlanBlockerExists(plan.FilePlan, "migration_candidate_set_mismatch")) {
			t.Fatalf("dirty blockers=%#v problems=%#v err=%v", plan.FilePlan.Blockers, plan.Problems, err)
		}
		outcome, err := ExecuteMigration(plan, MigrationApplyOptions{})
		if err != nil || outcome.Result.Status != reconcile.StatusBlocked || !bytes.Equal(before, mustReadFile(t, filepath.Join(root, filepath.FromSlash(path)))) {
			t.Fatalf("dirty execution=%#v err=%v", outcome, err)
		}
	})
}

func TestV2MigrationSupportsOptionalMixedGrandfatheringWithExactCutoverBytes(t *testing.T) {
	root := migrationRepository(t)
	alpha := "docs/repo/plans/alpha/alpha_discovery_doc.md"
	ledger := v2MigrationLedger(t, root, ModeMixed, nil, map[string]bool{alpha: true})
	plan, err := BuildMigrationPlan(root, ledger)
	if err != nil || len(plan.FilePlan.Blockers) != 0 || !plan.Projection.Ready || plan.Projection.Grandfathered != 1 || len(plan.FilePlan.Actions) != 1 || HasMaterialMigrationSkips(plan.Skips) {
		t.Fatalf("mixed plan=%#v projection=%#v problems=%#v skips=%#v err=%v", plan.FilePlan, plan.Projection, plan.Problems, plan.Skips, err)
	}

	path := filepath.Join(root, filepath.FromSlash(alpha))
	data := mustReadFile(t, path)
	if err := os.WriteFile(path, bytes.ReplaceAll(data, []byte("\n"), []byte("\r\n")), 0o644); err != nil {
		t.Fatal(err)
	}
	snapshot, err := LoadRepositorySnapshot(root)
	if err != nil || !problemForPathExists(snapshot.Problems, "mixed_grandfathered_plan_modified", alpha, SeverityError) {
		t.Fatalf("line-ending-only cutover drift was accepted: problems=%#v err=%v", snapshot.Problems, err)
	}
}

func TestV2LedgerRoundTripsThroughDurableSchema(t *testing.T) {
	root := legacyV2MigrationRepository(t)
	ledger := v2MigrationLedger(t, root, ModeCanonical, nil, nil)
	data, err := state.EncodeYAML(ledger)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadMigrationLedger(data)
	if err != nil || loaded.SchemaVersion != 2 || loaded.DiscoveryVersion != DiscoveryV2 || loaded.CandidateSetDigest != ledger.CandidateSetDigest || len(loaded.Records) != len(ledger.Records) {
		t.Fatalf("loaded=%#v err=%v\n%s", loaded, err, data)
	}
}

func legacyV2MigrationRepository(t *testing.T) string {
	t.Helper()
	root := migrationRepository(t)
	configPath := filepath.Join(root, filepath.FromSlash(state.ConfigPath))
	config := mustReadFile(t, configPath)
	lines := strings.Split(string(config), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(line, "plan_catalog_cutover_revision:") {
			continue
		}
		filtered = append(filtered, strings.Replace(line, "plan_catalog_mode: mixed", "plan_catalog_mode: legacy", 1))
	}
	if err := os.WriteFile(configPath, []byte(strings.Join(filtered, "\n")), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, root, "add", state.ConfigPath)
	runGitTest(t, root, "commit", "-m", "return fixture to legacy mode")
	return root
}

func v2MigrationLedger(t *testing.T, root string, targetMode CatalogMode, renames map[string]string, deferred map[string]bool) MigrationLedger {
	t.Helper()
	inventory, err := BuildInventoryWithOptions(root, time.Date(2026, 8, 6, 10, 0, 0, 0, time.UTC), SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, TargetCatalogMode: targetMode})
	if err != nil {
		t.Fatal(err)
	}
	decisions := map[string]MigrationDecision{
		"docs/repo/plans/alpha/alpha_discovery_doc.md": {
			ID: "example.discovery.alpha", Kind: KindDiscovery, Purpose: "Discover alpha", FirstCataloged: "2026-08-06T10:00:00Z", CatalogMetadataUpdated: "2026-08-06T10:00:00Z",
		},
		"docs/repo/plans/beta/beta_implementation_doc.md": {
			ID: "example.implementation.beta", Kind: KindImplementation, Purpose: "Implement beta", FirstCataloged: "2026-08-06T10:00:00Z", CatalogMetadataUpdated: "2026-08-06T10:00:00Z",
		},
		"products/widget/docs/notes/idea.md": {
			ID: "example.discovery.idea", Kind: KindDiscovery, Purpose: "Exercise repository-wide discovery.", FirstCataloged: "2026-08-06T00:00:00Z", CatalogMetadataUpdated: "2026-08-06T00:00:00Z",
		},
		"products/widget/docs/notes/collision.md": {
			ID: "example.discovery.collision", Kind: KindDiscovery, Purpose: "Exercise repository-wide discovery.", FirstCataloged: "2026-08-06T00:00:00Z", CatalogMetadataUpdated: "2026-08-06T00:00:00Z",
		},
		LegacyRegisterPath: {
			ID: "example.discovery.frozen-register", Kind: KindDiscovery, Purpose: "Hostile register signal", FirstCataloged: "2026-07-31T12:00:00Z", CatalogMetadataUpdated: "2026-07-31T12:00:00Z",
		},
	}
	records := []MigrationRecord{}
	for _, candidate := range inventory.Candidates {
		if candidate.Ownership != OwnershipOwned {
			continue
		}
		decision, ok := decisions[candidate.Path]
		if !ok {
			t.Fatalf("missing fixture decision for %s", candidate.Path)
		}
		target := candidate.Path
		state := "same-path"
		if renamed := renames[candidate.Path]; renamed != "" {
			target = renamed
			state = "absent"
		}
		records = append(records, MigrationRecord{
			CurrentPath: candidate.Path, TargetPath: target, SourceRevision: inventory.SourceRevision,
			SourceSHA256: candidate.Provenance.Source.ContentSHA256, TargetState: &TargetPrecondition{State: state}, Ownership: OwnershipOwned,
			Decision: decision, Confidence: "high", Ambiguity: []string{}, Evidence: []string{"reviewed fixture source"}, LegacyAliases: []string{}, Conflicts: []string{}, Deferred: deferred[candidate.Path], DeferralReason: "retained by reviewed mixed cutover", BranchOwner: "none",
		})
	}
	return MigrationLedger{
		SchemaVersion: 2, RepositoryID: "example", DiscoveryVersion: DiscoveryV2, TargetCatalogMode: targetMode,
		PolicyDigest: inventory.PolicyDigest, CandidateSetDigest: inventory.CandidateSetDigest, InventoryRevision: inventory.SourceRevision,
		ReviewedAt: "2026-08-06T10:00:00Z", Records: records,
	}
}

func reconcilePlanBlockerExists(plan reconcile.Plan, code string) bool {
	for _, blocker := range plan.Blockers {
		if blocker.Code == code {
			return true
		}
	}
	return false
}
