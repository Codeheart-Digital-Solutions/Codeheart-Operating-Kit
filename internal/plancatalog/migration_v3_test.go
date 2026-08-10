package plancatalog

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

const (
	v3DeferredPlanPath = "docs/repo/plans/alpha/alpha_discovery_doc.md"
	v3MigratedPlanPath = "docs/repo/plans/beta/beta_implementation_doc.md"
	v3OwnerRef         = "refs/heads/packaging-active"
	v3LedgerPath       = "docs/repo/plans/migrations/branch-ownership-ledger-v3.yaml"
)

type v3DeferralFixture struct {
	Root             string
	Ledger           MigrationLedger
	EvidenceRevision string
	LedgerRevision   string
	OwnerTip         string
	CutoverRevision  string
}

func TestV3MixedActiveOwnerDeferralAndELAActivation(t *testing.T) {
	fixture := newV3DeferralFixture(t)
	plan := requireReadyV3MigrationPlan(t, fixture)
	if plan.Projection.Candidates != 2 || plan.Projection.Canonical != 1 || plan.Projection.Grandfathered != 1 || plan.Projection.Gaps != 0 || !plan.Projection.MixedCoverageComplete || plan.Projection.CanonicalReady {
		t.Fatalf("unexpected mixed projection: %#v", plan.Projection)
	}
	if len(plan.BranchReviews) != 1 || !plan.BranchReviews[0].Verified || plan.BranchReviews[0].Disposition != BranchDispositionDeferredActiveOwner {
		t.Fatalf("active owner was not verified as an exact deferral: %#v", plan.BranchReviews)
	}
	if len(plan.FollowUps) != 1 || plan.FollowUps[0].Action != "reinventory-and-incremental-migrate" || plan.FollowUps[0].State != "pending" {
		t.Fatalf("mandatory incremental follow-up missing: %#v", plan.FollowUps)
	}
	deferredBefore := mustReadFile(t, filepath.Join(fixture.Root, filepath.FromSlash(v3DeferredPlanPath)))

	migration, err := ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 8, 9, 12, 0, 0, 0, time.UTC)})
	if err != nil || migration.Result.Status != reconcile.StatusSucceeded || len(migration.Result.Changes) != 1 || migration.ActivationPerformed == nil || *migration.ActivationPerformed {
		t.Fatalf("schema-v3 migration failed: outcome=%#v err=%v", migration, err)
	}
	if actual := mustReadFile(t, filepath.Join(fixture.Root, filepath.FromSlash(v3DeferredPlanPath))); !bytes.Equal(actual, deferredBefore) {
		t.Fatal("mixed migration changed the deferred owner path")
	}

	activationPlan, err := BuildCatalogActivationPlan(fixture.Root, fixture.Ledger)
	if err != nil || len(activationPlan.FilePlan.Blockers) != 0 || activationPlan.Checkpoint.ActivationBaseRevision != fixture.LedgerRevision || len(activationPlan.FollowUps) != 1 {
		t.Fatalf("activation plan blockers=%#v problems=%#v followups=%#v err=%v", activationPlan.FilePlan.Blockers, activationPlan.Problems, activationPlan.FollowUps, err)
	}
	activation, err := ExecuteCatalogActivation(activationPlan, CatalogActivationOptions{Now: time.Date(2026, 8, 9, 12, 1, 0, 0, time.UTC)})
	if err != nil || activation.Result.Status != reconcile.StatusSucceeded || !activation.ActivationPerformed || !activation.ActivationCheckpointRequired {
		t.Fatalf("activation failed: outcome=%#v err=%v", activation, err)
	}
	runGitTest(t, fixture.Root, "add", "--all")
	runGitTest(t, fixture.Root, "commit", "-m", "activate reviewed mixed catalog")
	activationRevision := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
	if parent := runGitTest(t, fixture.Root, "rev-parse", "HEAD^1"); parent != fixture.LedgerRevision {
		t.Fatalf("activation checkpoint parent=%s want ledger checkpoint=%s", parent, fixture.LedgerRevision)
	}
	if firstParent := runGitTest(t, fixture.Root, "rev-parse", fixture.LedgerRevision+"^1"); firstParent != fixture.EvidenceRevision {
		t.Fatalf("ledger checkpoint parent=%s want evidence=%s", firstParent, fixture.EvidenceRevision)
	}
	changed := runGitTest(t, fixture.Root, "diff-tree", "--no-commit-id", "--name-only", "-r", fixture.LedgerRevision, activationRevision)
	if !strings.Contains(changed, state.ConfigPath) || !strings.Contains(changed, v3MigratedPlanPath) || strings.Contains(changed, v3DeferredPlanPath) {
		t.Fatalf("activation delta did not preserve the deferred plan: %q", changed)
	}

	snapshot, err := LoadRepositorySnapshot(fixture.Root)
	if err != nil || HasErrors(snapshot.Problems) || snapshot.Settings.ConfigSchemaVersion != 2 || snapshot.Settings.DiscoveryVersion != DiscoveryV2 || snapshot.Settings.Mode != ModeMixed || !snapshot.DeferredPaths[v3DeferredPlanPath] {
		t.Fatalf("activated E->L->A state invalid: settings=%#v deferred=%#v problems=%#v err=%v", snapshot.Settings, snapshot.DeferredPaths, snapshot.Problems, err)
	}
	view := BuildView(snapshot)
	deferredVisible := false
	for _, row := range view.Rows {
		if row.CanonicalPath == v3DeferredPlanPath && row.CoverageDisposition == "mixed-grandfathered" && row.IncrementalMigrationRequired {
			deferredVisible = true
		}
	}
	if !deferredVisible {
		t.Fatalf("activated canonical view silently omitted the deferred plan: %#v", view.Rows)
	}
}

func TestV3PersistedActivationEvidenceSurvivesNormalMerge(t *testing.T) {
	fixture, _ := activateV3DeferralFixture(t)
	activationRevision := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
	runGitTest(t, fixture.Root, "branch", "reviewed-activation", activationRevision)
	runGitTest(t, fixture.Root, "switch", "-c", "integration-main", fixture.EvidenceRevision)
	runGitTest(t, fixture.Root, "merge", "--no-ff", "reviewed-activation", "-m", "merge reviewed activation")
	mergeRevision := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
	if secondParent := runGitTest(t, fixture.Root, "rev-parse", mergeRevision+"^2"); secondParent != activationRevision {
		t.Fatalf("normal merge second parent=%s want activation=%s", secondParent, activationRevision)
	}
	if mergeTree, activationTree := runGitTest(t, fixture.Root, "rev-parse", mergeRevision+"^{tree}"), runGitTest(t, fixture.Root, "rev-parse", activationRevision+"^{tree}"); mergeTree != activationTree {
		t.Fatalf("normal merge tree=%s want activation tree=%s", mergeTree, activationTree)
	}

	snapshot, err := LoadRepositorySnapshot(fixture.Root)
	if err != nil || HasErrors(snapshot.Problems) {
		t.Fatalf("normal merge incorporating reviewed activation was rejected: problems=%#v err=%v", snapshot.Problems, err)
	}
	writeCheckpointFixture(t, fixture.Root, "later.txt", []byte("safe descendant\n"))
	runGitTest(t, fixture.Root, "add", "later.txt")
	runGitTest(t, fixture.Root, "commit", "-m", "record safe descendant")
	snapshot, err = LoadRepositorySnapshot(fixture.Root)
	if err != nil || HasErrors(snapshot.Problems) {
		t.Fatalf("safe descendant of merged activation was rejected: problems=%#v err=%v", snapshot.Problems, err)
	}
}

func TestV3PersistedActivationEvidenceRejectsTreeOnlyCoincidence(t *testing.T) {
	fixture, _ := activateV3DeferralFixture(t)
	activationRevision := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
	activationTree := runGitTest(t, fixture.Root, "rev-parse", activationRevision+"^{tree}")
	forged := runGitTest(t, fixture.Root, "commit-tree", activationTree, "-p", fixture.EvidenceRevision, "-m", "reparent activation tree")
	runGitTest(t, fixture.Root, "branch", "tree-only-coincidence", forged)
	runGitTest(t, fixture.Root, "switch", "tree-only-coincidence")
	if currentTree := runGitTest(t, fixture.Root, "rev-parse", "HEAD^{tree}"); currentTree != activationTree {
		t.Fatalf("tree-only fixture tree=%s want activation tree=%s", currentTree, activationTree)
	}
	snapshot, err := LoadRepositorySnapshot(fixture.Root)
	if err != nil {
		t.Fatal(err)
	}
	if !problemCodePresent(snapshot.Problems, "activation_checkpoint_missing") {
		t.Fatalf("tree-only coincidence without activation ancestry was accepted: %#v", snapshot.Problems)
	}
}

func TestV3PersistedActivationEvidenceRejectsForgedSiblingConfig(t *testing.T) {
	fixture, activation := activateV3DeferralFixture(t)
	reviewedActivation := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
	reviewedConfig, err := gitBytes(fixture.Root, "show", reviewedActivation+":"+state.ConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	config, err := state.DecodeYAMLMap(reviewedConfig)
	if err != nil {
		t.Fatal(err)
	}
	config["project_display_name"] = "Forged sibling activation"
	forgedConfig, err := state.EncodeYAML(config)
	if err != nil {
		t.Fatal(err)
	}
	reviewedPlan, err := gitBytes(fixture.Root, "show", reviewedActivation+":"+v3MigratedPlanPath)
	if err != nil {
		t.Fatal(err)
	}
	runGitTest(t, fixture.Root, "switch", "-c", "forged-config-sibling", fixture.LedgerRevision)
	writeCheckpointFixture(t, fixture.Root, state.ConfigPath, forgedConfig)
	writeCheckpointFixture(t, fixture.Root, v3MigratedPlanPath, reviewedPlan)
	runGitTest(t, fixture.Root, "add", state.ConfigPath, v3MigratedPlanPath)
	runGitTest(t, fixture.Root, "commit", "-m", "forge sibling with copied action delta")
	forged := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
	if _, ancestryErr := gitText(fixture.Root, "merge-base", "--is-ancestor", reviewedActivation, forged); ancestryErr == nil {
		t.Fatal("forged sibling unexpectedly incorporates the reviewed activation")
	}
	if digest, _, digestErr := MigrationActionDigestForCheckpoint(fixture.Root, fixture.LedgerRevision, forged); digestErr != nil || digest != activation.MigrationActionDigest {
		t.Fatalf("forged sibling did not preserve the reviewed action digest: digest=%s want=%s err=%v", digest, activation.MigrationActionDigest, digestErr)
	}
	snapshot, err := LoadRepositorySnapshot(fixture.Root)
	if err != nil {
		t.Fatal(err)
	}
	if !problemCodePresent(snapshot.Problems, "activation_checkpoint_config_mismatch") {
		t.Fatalf("forged non-binding config bytes were accepted: %#v", snapshot.Problems)
	}
}

func TestV3PersistedActivationEvidenceRejectsWrongBaseAndPlanBlob(t *testing.T) {
	t.Run("wrong activation base binding", func(t *testing.T) {
		fixture, _ := activateV3DeferralFixture(t)
		reviewedActivation := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
		reviewedConfig, err := gitBytes(fixture.Root, "show", reviewedActivation+":"+state.ConfigPath)
		if err != nil {
			t.Fatal(err)
		}
		config, err := state.DecodeYAMLMap(reviewedConfig)
		if err != nil {
			t.Fatal(err)
		}
		planning := state.Map(state.Map(config["component_settings"])["planning-workflows"])
		binding := state.Map(planning["plan_catalog_migration_evidence"])
		binding["activation_base_revision"] = fixture.EvidenceRevision
		wrongBinding, err := state.EncodeYAML(config)
		if err != nil {
			t.Fatal(err)
		}
		reviewedPlan, err := gitBytes(fixture.Root, "show", reviewedActivation+":"+v3MigratedPlanPath)
		if err != nil {
			t.Fatal(err)
		}
		runGitTest(t, fixture.Root, "switch", "-c", "wrong-activation-base", fixture.LedgerRevision)
		writeCheckpointFixture(t, fixture.Root, state.ConfigPath, wrongBinding)
		writeCheckpointFixture(t, fixture.Root, v3MigratedPlanPath, reviewedPlan)
		runGitTest(t, fixture.Root, "add", state.ConfigPath, v3MigratedPlanPath)
		runGitTest(t, fixture.Root, "commit", "-m", "bind activation to wrong ledger base")
		snapshot, err := LoadRepositorySnapshot(fixture.Root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemCodePresent(snapshot.Problems, "ledger_checkpoint_parent_mismatch") || !problemCodePresent(snapshot.Problems, "activation_checkpoint_missing") {
			t.Fatalf("wrong ledger activation base was accepted: %#v", snapshot.Problems)
		}
	})

	t.Run("wrong projected plan blob", func(t *testing.T) {
		fixture, _ := activateV3DeferralFixture(t)
		reviewedActivation := runGitTest(t, fixture.Root, "rev-parse", "HEAD")
		reviewedConfig, err := gitBytes(fixture.Root, "show", reviewedActivation+":"+state.ConfigPath)
		if err != nil {
			t.Fatal(err)
		}
		reviewedPlan, err := gitBytes(fixture.Root, "show", reviewedActivation+":"+v3MigratedPlanPath)
		if err != nil {
			t.Fatal(err)
		}
		runGitTest(t, fixture.Root, "switch", "-c", "wrong-plan-blob", fixture.LedgerRevision)
		writeCheckpointFixture(t, fixture.Root, state.ConfigPath, reviewedConfig)
		writeCheckpointFixture(t, fixture.Root, v3MigratedPlanPath, append(reviewedPlan, []byte("\nforged projection\n")...))
		runGitTest(t, fixture.Root, "add", state.ConfigPath, v3MigratedPlanPath)
		runGitTest(t, fixture.Root, "commit", "-m", "forge projected plan bytes")
		snapshot, err := LoadRepositorySnapshot(fixture.Root)
		if err != nil {
			t.Fatal(err)
		}
		if !problemCodePresent(snapshot.Problems, "activation_checkpoint_missing") {
			t.Fatalf("wrong projected blob retained reviewed activation authority: %#v", snapshot.Problems)
		}
	})
}

func TestV3DirectCanonicalRejectsActiveOwnerDeferral(t *testing.T) {
	fixture := newV3DeferralFixture(t)
	ledger := fixture.Ledger
	ledger.TargetMode = ModeCanonical
	ledger.TargetCatalogMode = ModeCanonical
	plan, err := BuildMigrationPlan(fixture.Root, ledger)
	if err != nil {
		t.Fatal(err)
	}
	if !reconcilePlanBlockerExists(plan.FilePlan, "canonical_deferral_forbidden") || plan.Projection.CanonicalReady {
		t.Fatalf("direct canonical deferral was not rejected: actions=%#v projection=%#v blockers=%#v problems=%#v", plan.FilePlan.Actions, plan.Projection, plan.FilePlan.Blockers, plan.Problems)
	}
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{})
	if err != nil || outcome.Result.Status != reconcile.StatusBlocked {
		t.Fatalf("direct canonical execution=%#v err=%v", outcome, err)
	}
}

func TestV3MovedRefAndTamperedReviewedEvidenceFailClosed(t *testing.T) {
	t.Run("moved ref after review", func(t *testing.T) {
		fixture := newV3DeferralFixture(t)
		runGitTest(t, fixture.Root, "switch", "packaging-active")
		appendFile(t, filepath.Join(fixture.Root, filepath.FromSlash(v3DeferredPlanPath)), "\nOwner work moved after review.\n")
		commitBranchEvidenceAll(t, fixture.Root, "advance reviewed owner")
		runGitTest(t, fixture.Root, "switch", "main")
		plan, err := BuildMigrationPlan(fixture.Root, fixture.Ledger)
		if err != nil {
			t.Fatal(err)
		}
		if !reconcilePlanBlockerExists(plan.FilePlan, "branch_evidence_ref_moved") && !reconcilePlanBlockerExists(plan.FilePlan, "branch_ref_moved") {
			t.Fatalf("moved ref was accepted: blockers=%#v problems=%#v", plan.FilePlan.Blockers, plan.Problems)
		}
		outcome, err := ExecuteMigration(plan, MigrationApplyOptions{})
		if err != nil || outcome.Result.Status != reconcile.StatusBlocked {
			t.Fatalf("moved-ref execution=%#v err=%v", outcome, err)
		}
	})

	t.Run("tampered review transition", func(t *testing.T) {
		fixture := newV3DeferralFixture(t)
		ledger := fixture.Ledger
		ledger.Records[0].BranchTouchReviews[0].PathTransitions[0].After.SHA256 = strings.Repeat("f", 64)
		plan, err := BuildMigrationPlan(fixture.Root, ledger)
		if err != nil {
			t.Fatal(err)
		}
		if !reconcilePlanBlockerExists(plan.FilePlan, "branch_path_state_moved") && !reconcilePlanBlockerExists(plan.FilePlan, "branch_evidence_digest_mismatch") {
			t.Fatalf("tampered transition was accepted: blockers=%#v problems=%#v", plan.FilePlan.Blockers, plan.Problems)
		}
	})

	t.Run("forbidden proof on deferred review", func(t *testing.T) {
		fixture := newV3DeferralFixture(t)
		data := mustReadFile(t, filepath.Join(fixture.Root, filepath.FromSlash(v3LedgerPath)))
		value, err := state.DecodeYAMLMap(data)
		if err != nil {
			t.Fatal(err)
		}
		reviews := state.AnySlice(state.Map(state.AnySlice(value["records"])[0])["branch_touch_reviews"])
		review := state.Map(reviews[0])
		if state.AsString(review["disposition"]) != BranchDispositionDeferredActiveOwner {
			t.Fatalf("fixture did not contain the deferred disposition: %#v", review)
		}
		review["proof"] = map[string]any{
			"kind":              "same-content-v1",
			"target_path_state": map[string]any{"state": "present", "mode": "100644", "object_id": strings.Repeat("1", 40), "sha256": strings.Repeat("2", 64)},
			"ref_path_state":    map[string]any{"state": "present", "mode": "100644", "object_id": strings.Repeat("1", 40), "sha256": strings.Repeat("2", 64)},
		}
		tampered, err := state.EncodeYAML(value)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := LoadMigrationLedger(tampered); err == nil {
			t.Fatal("schema-validity route accepted a proof forbidden for deferred-active-owner")
		}
	})
}

func TestV3PersistedConfigAndLedgerBindingsDetectTampering(t *testing.T) {
	fixture, _ := activateV3DeferralFixture(t)
	configPath := filepath.Join(fixture.Root, filepath.FromSlash(state.ConfigPath))
	configBefore := mustReadFile(t, configPath)
	tamperedConfig := bytes.Replace(configBefore, []byte("branch_evidence_digest: "+fixture.Ledger.BranchEvidence.Digest), []byte("branch_evidence_digest: "+strings.Repeat("f", 64)), 1)
	if bytes.Equal(tamperedConfig, configBefore) {
		t.Fatal("activated config did not contain the branch evidence binding")
	}
	if err := os.WriteFile(configPath, tamperedConfig, 0o644); err != nil {
		t.Fatal(err)
	}
	snapshot, err := LoadRepositorySnapshot(fixture.Root)
	if err != nil || (!problemCodePresent(snapshot.Problems, "migration_evidence_binding_mismatch") && !problemCodePresent(snapshot.Problems, "migration_evidence_config_moved")) {
		t.Fatalf("tampered config binding was accepted: problems=%#v err=%v", snapshot.Problems, err)
	}
	if err := os.WriteFile(configPath, configBefore, 0o644); err != nil {
		t.Fatal(err)
	}

	ledgerPath := filepath.Join(fixture.Root, filepath.FromSlash(v3LedgerPath))
	appendFile(t, ledgerPath, "\n# unauthorized ledger mutation\n")
	snapshot, err = LoadRepositorySnapshot(fixture.Root)
	if err != nil || !problemCodePresent(snapshot.Problems, "migration_evidence_ledger_mismatch") {
		t.Fatalf("tampered ledger binding was accepted: problems=%#v err=%v", snapshot.Problems, err)
	}
}

func TestV3OwnerMergeWithReviewedMetadataResolvesIncrementalFollowUp(t *testing.T) {
	fixture, activation := activateV3DeferralFixture(t)
	if len(activation.IncrementalFollowUps) != 1 || activation.IncrementalFollowUps[0].State != "pending" {
		t.Fatalf("activation did not surface pending incremental work: %#v", activation.IncrementalFollowUps)
	}
	runGitTest(t, fixture.Root, "merge", "--no-ff", "packaging-active", "-m", "integrate active plan owner")
	snapshot, err := LoadRepositorySnapshot(fixture.Root)
	if err != nil || HasErrors(snapshot.Problems) {
		t.Fatalf("reviewed owner metadata did not resolve cleanly: problems=%#v err=%v", snapshot.Problems, err)
	}
	if snapshot.DeferredPaths[v3DeferredPlanPath] {
		t.Fatal("integrated canonical owner remained marked as a pending deferral")
	}
	view := BuildView(snapshot)
	if view.CanonicalReady == nil || !*view.CanonicalReady {
		t.Fatalf("integrated owner did not close canonical readiness: %#v", view)
	}
}

func TestV3IntegratedOwnerMetadataTamperingFailsClosed(t *testing.T) {
	fixture, _ := activateV3DeferralFixture(t)
	runGitTest(t, fixture.Root, "merge", "--no-ff", "packaging-active", "-m", "integrate active plan owner")
	path := filepath.Join(fixture.Root, filepath.FromSlash(v3DeferredPlanPath))
	data := mustReadFile(t, path)
	data = bytes.Replace(data, []byte("purpose: Discover alpha"), []byte("purpose: Tampered purpose"), 1)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	commitBranchEvidenceAll(t, fixture.Root, "tamper integrated owner metadata")
	snapshot, err := LoadRepositorySnapshot(fixture.Root)
	if err != nil {
		t.Fatal(err)
	}
	if !HasErrors(snapshot.Problems) || (!problemCodePresent(snapshot.Problems, "deferred_owner_merge_mismatch") && !problemCodePresent(snapshot.Problems, "deferred_owner_metadata_mismatch")) {
		t.Fatalf("tampered integrated owner metadata was accepted: %#v", snapshot.Problems)
	}
}

func TestV3UnstagedOwnerOverlayCannotSatisfyCommittedAuthority(t *testing.T) {
	fixture, _ := activateV3DeferralFixture(t)
	runGitTest(t, fixture.Root, "merge", "--no-ff", "packaging-active", "-m", "integrate active plan owner")
	path := filepath.Join(fixture.Root, filepath.FromSlash(v3DeferredPlanPath))
	reviewed := mustReadFile(t, path)
	if err := os.WriteFile(path, append(append([]byte{}, reviewed...), []byte("\nCommitted drift.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	commitBranchEvidenceAll(t, fixture.Root, "commit incorrect integrated owner bytes")
	if err := os.WriteFile(path, reviewed, 0o644); err != nil {
		t.Fatal(err)
	}
	snapshot, err := LoadRepositorySnapshot(fixture.Root)
	if err != nil {
		t.Fatal(err)
	}
	if !HasErrors(snapshot.Problems) || (!problemCodePresent(snapshot.Problems, "deferred_owner_merge_mismatch") && !problemCodePresent(snapshot.Problems, "migration_evidence_worktree_dirty")) {
		t.Fatalf("unstaged reviewed bytes masked incorrect committed authority: %#v", snapshot.Problems)
	}
}

func TestV3PostActivationCandidateAdditionRequiresIncrementalMigration(t *testing.T) {
	fixture, _ := activateV3DeferralFixture(t)
	path := "docs/repo/plans/new/new_implementation_doc.md"
	writeCheckpointFixture(t, fixture.Root, path, branchEvidencePlanDocument("active", "example.implementation.new"))
	commitBranchEvidenceAll(t, fixture.Root, "add post-activation candidate")
	snapshot, err := LoadRepositorySnapshot(fixture.Root)
	if err != nil {
		t.Fatal(err)
	}
	if !problemCodePresent(snapshot.Problems, "migration_evidence_candidate_set_moved") {
		t.Fatalf("post-activation candidate addition did not force incremental migration: %#v", snapshot.Problems)
	}
}

func TestV3PostWriteActionDigestDriftRollsBack(t *testing.T) {
	fixture := newV3DeferralFixture(t)
	plan := requireReadyV3MigrationPlan(t, fixture)
	if len(plan.FilePlan.Actions) != 1 {
		t.Fatalf("expected one migration action, got %#v", plan.FilePlan.Actions)
	}
	before := mustReadFile(t, filepath.Join(fixture.Root, filepath.FromSlash(plan.FilePlan.Actions[0].Target)))
	plan.FilePlan.Actions[0].Content = append(append([]byte{}, plan.FilePlan.Actions[0].Content...), []byte("\nTampered after review.\n")...)
	outcome, err := ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 8, 9, 13, 0, 0, 0, time.UTC)})
	if err != nil || outcome.Result.Status != reconcile.StatusRolledBack || !outcome.Result.Rollback.Succeeded || !reconcileBlockerExists(outcome.Result, "post_write_authority_drift") {
		t.Fatalf("post-write drift outcome=%#v err=%v", outcome, err)
	}
	after := mustReadFile(t, filepath.Join(fixture.Root, filepath.FromSlash(plan.FilePlan.Actions[0].Target)))
	if !bytes.Equal(after, before) {
		t.Fatal("post-write authority rollback did not restore exact plan bytes")
	}
}

func TestV3ActivationPostWriteCRLFDriftWithoutCheckoutAuthorityRollsBack(t *testing.T) {
	fixture := newV3DeferralFixture(t)
	plan := requireReadyV3MigrationPlan(t, fixture)
	migration, err := ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 8, 9, 13, 10, 0, 0, time.UTC)})
	if err != nil || migration.Result.Status != reconcile.StatusSucceeded {
		t.Fatalf("migration before activation=%#v err=%v", migration, err)
	}
	activationPlan, err := BuildCatalogActivationPlan(fixture.Root, fixture.Ledger)
	if err != nil || len(activationPlan.FilePlan.Blockers) != 0 || len(activationPlan.MigrationActions) != 1 {
		t.Fatalf("activation plan blockers=%#v problems=%#v err=%v", activationPlan.FilePlan.Blockers, activationPlan.Problems, err)
	}
	configBefore := mustReadFile(t, filepath.Join(fixture.Root, filepath.FromSlash(state.ConfigPath)))
	target := activationPlan.MigrationActions[0].Target
	tampered := false
	activation, err := ExecuteCatalogActivation(activationPlan, CatalogActivationOptions{
		Now: time.Date(2026, 8, 9, 13, 11, 0, 0, time.UTC),
		Hook: func(phase string) error {
			if phase != "commit:"+state.ConfigPath || tampered {
				return nil
			}
			tampered = true
			current := mustReadFile(t, filepath.Join(fixture.Root, filepath.FromSlash(target)))
			writeCheckpointFixture(t, fixture.Root, target, bytes.ReplaceAll(current, []byte("\n"), []byte("\r\n")))
			runGitTest(t, fixture.Root, "config", "core.autocrlf", "false")
			return nil
		},
	})
	if err != nil || !tampered || activation.Result.Status != reconcile.StatusRolledBack || !activation.Result.Rollback.Succeeded || !reconcileBlockerExists(activation.Result, "catalog_activation_authority_drift") {
		t.Fatalf("activation drift outcome=%#v tampered=%t err=%v", activation, tampered, err)
	}
	if configAfter := mustReadFile(t, filepath.Join(fixture.Root, filepath.FromSlash(state.ConfigPath))); !bytes.Equal(configAfter, configBefore) {
		t.Fatal("activation rollback did not restore exact config bytes")
	}
}

func newV3DeferralFixture(t *testing.T) v3DeferralFixture {
	t.Helper()
	root := migrationRepository(t)
	writeCheckpointFixture(t, root, ".gitignore", []byte(".codeheart/local/\n.codeheart/kit.transaction.json\n"))
	runGitTest(t, root, "add", ".gitignore")
	runGitTest(t, root, "commit", "-m", "ignore local transaction state")
	evidenceRevision := runGitTest(t, root, "rev-parse", "HEAD")
	settings, problems := LoadRepositorySettings(root)
	if HasErrors(problems) || settings.CutoverRevision == "" {
		t.Fatalf("fixture mixed settings invalid: %#v", problems)
	}

	runGitTest(t, root, "switch", "-c", "packaging-active")
	ownerPath := filepath.Join(root, filepath.FromSlash(v3DeferredPlanPath))
	ownerData := mustReadFile(t, ownerPath)
	ownerData = bytes.Replace(ownerData, []byte("Status: draft"), []byte("Status: active"), 1)
	ownerData, err := InsertMetadata(ownerData, Metadata{
		SchemaVersion: 1, ID: "example.discovery.alpha", Kind: KindDiscovery, Purpose: "Discover alpha",
		FirstCataloged: "2026-08-09T10:00:00Z", CatalogMetadataUpdated: "2026-08-09T10:00:00Z",
		LegacyAliases: []string{"PR-001"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ownerPath, ownerData, 0o644); err != nil {
		t.Fatal(err)
	}
	commitBranchEvidenceAll(t, root, "continue active plan on owner branch")
	ownerTip := runGitTest(t, root, "rev-parse", "HEAD")
	runGitTest(t, root, "switch", "main")
	if head := runGitTest(t, root, "rev-parse", "HEAD"); head != evidenceRevision {
		t.Fatalf("fixture main moved while creating owner branch: %s != %s", head, evidenceRevision)
	}

	inventory, err := BuildInventoryWithOptions(root, time.Date(2026, 8, 9, 10, 0, 0, 0, time.UTC), SnapshotOptions{TargetDiscoveryVersion: DiscoveryV2, TargetCatalogMode: ModeMixed})
	if err != nil {
		t.Fatal(err)
	}
	candidates := map[string]Candidate{}
	for _, candidate := range inventory.Candidates {
		if candidate.Ownership == OwnershipOwned {
			candidates[candidate.Path] = candidate
		}
	}
	alpha, alphaOK := candidates[v3DeferredPlanPath]
	beta, betaOK := candidates[v3MigratedPlanPath]
	if !alphaOK || !betaOK || len(candidates) != 2 {
		t.Fatalf("unexpected fixture candidates: %#v", candidates)
	}
	evidence := proposeBranchCandidateEvidence(BranchEvidenceRequest{
		Root: root, EvidenceScope: "local", RepositoryID: inventory.RepositoryID,
		LogicalRef: "local:" + v3OwnerRef, Ref: v3OwnerRef, EvidenceRevision: evidenceRevision,
		PolicyDigest: inventory.PolicyDigest, CandidateSetDigest: inventory.CandidateSetDigest,
		CandidatePath: v3DeferredPlanPath, Paths: []string{v3DeferredPlanPath},
	})
	if evidence.ProposedDisposition != BranchDispositionActiveOwner || len(evidence.Blockers) != 0 {
		t.Fatalf("fixture owner was not classified active: %#v", evidence)
	}
	ownerState, blocker := ReadBranchPathState(root, ownerTip, v3DeferredPlanPath)
	if blocker != nil {
		t.Fatal(blocker)
	}

	decision := func(id string, kind Kind, purpose string) MigrationDecision {
		return MigrationDecision{ID: id, Kind: kind, Purpose: purpose, FirstCataloged: "2026-08-09T10:00:00Z", CatalogMetadataUpdated: "2026-08-09T10:00:00Z"}
	}
	record := func(candidate Candidate, migrationDecision MigrationDecision, alias string) MigrationRecord {
		return MigrationRecord{
			CurrentPath: candidate.Path, TargetPath: candidate.Path, SourceRevision: evidenceRevision,
			SourceSHA256: candidate.Provenance.Source.ContentSHA256, TargetState: &TargetPrecondition{State: "same-path"},
			Ownership: OwnershipOwned, Decision: migrationDecision, Confidence: "high",
			Ambiguity: []string{}, Evidence: []string{"reviewed synthetic public-safe fixture"}, LegacyAliases: []string{alias}, Conflicts: []string{}, BranchTouchReviews: []BranchTouchReview{},
		}
	}
	alphaRecord := record(alpha, decision("example.discovery.alpha", KindDiscovery, "Discover alpha"), "PR-001")
	alphaRecord.Deferred = true
	alphaRecord.DeferralReason = "reviewed active owner remains branch-owned"
	alphaRecord.BranchTouchReviews = []BranchTouchReview{{
		Identity: evidence.Identity, Disposition: BranchDispositionDeferredActiveOwner,
		RefTip: evidence.RefTip, MergeBase: evidence.MergeBase, PathTransitions: evidence.PathTransitions,
		OwnerTipCandidate: &OwnerTipCandidate{Path: v3DeferredPlanPath, Mode: ownerState.Mode, ObjectID: ownerState.ObjectID, SHA256: ownerState.SHA256, Lifecycle: LifecycleActive},
		IncrementalFollowUp: &IncrementalFollowUp{
			Action: "reinventory-and-incremental-migrate", Trigger: "owner-ref-integrated-or-target-path-changed",
			CutoverRevision: settings.CutoverRevision, PlanPath: v3DeferredPlanPath, CutoverSourceSHA256: alphaRecord.SourceSHA256,
			OwnerRef: evidence.Identity.LogicalRef, OwnerTip: ownerTip, OwnerPath: v3DeferredPlanPath,
			OwnerMode: ownerState.Mode, OwnerObjectID: ownerState.ObjectID, OwnerSHA256: ownerState.SHA256, State: "pending",
			SuccessConditions: []string{"owner-supplied-canonical-metadata-at-merge", "new-reviewed-ledger-applied-to-merged-target-bytes"},
		},
	}}
	betaRecord := record(beta, decision("example.implementation.beta", KindImplementation, "Implement beta"), "PR-002")
	configSHA, err := configPrecondition(root)
	if err != nil {
		t.Fatal(err)
	}
	ledger := MigrationLedger{
		SchemaVersion: 3, RepositoryID: inventory.RepositoryID, EvidenceRevision: evidenceRevision,
		TargetMode: ModeMixed, PolicyDigest: inventory.PolicyDigest, CandidateSetDigest: inventory.CandidateSetDigest,
		InventoryRevision: evidenceRevision, TargetConfigPreconditionSHA256: configSHA,
		BranchEvidence: &BranchEvidenceSummary{
			Algorithm: BranchEvidenceAlgorithm, GitVersion: evidence.GitVersion, EvidenceScope: "local",
			RemoteOverlayStatus: "not-requested", Digest: branchEvidenceSetDigest([]BranchCandidateEvidence{evidence}, ""),
		},
		Records: []MigrationRecord{alphaRecord, betaRecord},
	}
	ledgerData := encodeV3LedgerFixture(t, ledger)
	writeCheckpointFixture(t, root, v3LedgerPath, ledgerData)
	runGitTest(t, root, "add", v3LedgerPath)
	runGitTest(t, root, "commit", "-m", "record reviewed branch ownership ledger")
	ledgerRevision := runGitTest(t, root, "rev-parse", "HEAD")
	loaded, err := LoadMigrationLedgerArtifact(root, filepath.Join(root, filepath.FromSlash(v3LedgerPath)))
	if err != nil {
		t.Fatalf("load committed v3 ledger: %v\n%s", err, ledgerData)
	}
	for _, item := range loaded.Records {
		if item.BranchOwner != "" {
			t.Fatalf("schema-v3 fixture unexpectedly serialized retired branch_owner=%q", item.BranchOwner)
		}
	}
	return v3DeferralFixture{Root: root, Ledger: loaded, EvidenceRevision: evidenceRevision, LedgerRevision: ledgerRevision, OwnerTip: ownerTip, CutoverRevision: settings.CutoverRevision}
}

func encodeV3LedgerFixture(t *testing.T, ledger MigrationLedger) []byte {
	t.Helper()
	data, err := state.EncodeYAML(ledger)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := LoadMigrationLedger(data); err != nil {
		t.Fatalf("generated v3 ledger is invalid: %v\n%s", err, data)
	}
	return data
}

func TestV3LedgerMarshalsDirectlyToStrictWireContract(t *testing.T) {
	fixture := newV3DeferralFixture(t)
	data, err := state.EncodeYAML(fixture.Ledger)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("target_precondition:\n")) || !bytes.Contains(data, []byte("target_precondition: same-path")) {
		t.Fatalf("v3 target precondition did not use the scalar wire contract:\n%s", data)
	}
	if _, err := LoadMigrationLedger(data); err != nil {
		t.Fatalf("directly encoded v3 ledger did not round-trip: %v\n%s", err, data)
	}
	exact := fixture.Ledger
	exact.Records[0].TargetState = &TargetPrecondition{State: "exact"}
	exact.Records[0].TargetPreconditionSHA256 = strings.Repeat("a", 64)
	exactYAML, err := state.EncodeYAML(exact)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadMigrationLedger(exactYAML)
	if err != nil || loaded.Records[0].TargetState == nil || loaded.Records[0].TargetState.SHA256 != strings.Repeat("a", 64) {
		t.Fatalf("exact YAML precondition did not round-trip: ledger=%#v err=%v\n%s", loaded, err, exactYAML)
	}
	exactJSON, err := json.Marshal(exact)
	if err != nil {
		t.Fatal(err)
	}
	var exactValue map[string]any
	if err := json.Unmarshal(exactJSON, &exactValue); err != nil {
		t.Fatal(err)
	}
	if err := state.Validate(state.PlanMigrationV3Schema, exactValue); err != nil {
		t.Fatalf("exact JSON wire contract is invalid: %v\n%s", err, exactJSON)
	}
}

func requireReadyV3MigrationPlan(t *testing.T, fixture v3DeferralFixture) MigrationPlan {
	t.Helper()
	plan, err := BuildMigrationPlan(fixture.Root, fixture.Ledger)
	if err != nil || len(plan.FilePlan.Blockers) != 0 || !plan.Projection.Ready {
		t.Fatalf("v3 plan not ready: blockers=%#v problems=%#v projection=%#v err=%v", plan.FilePlan.Blockers, plan.Problems, plan.Projection, err)
	}
	return plan
}

func activateV3DeferralFixture(t *testing.T) (v3DeferralFixture, CatalogActivationOutcome) {
	t.Helper()
	fixture := newV3DeferralFixture(t)
	plan := requireReadyV3MigrationPlan(t, fixture)
	migration, err := ExecuteMigration(plan, MigrationApplyOptions{Now: time.Date(2026, 8, 9, 14, 0, 0, 0, time.UTC)})
	if err != nil || migration.Result.Status != reconcile.StatusSucceeded {
		t.Fatalf("migration before activation=%#v err=%v", migration, err)
	}
	activationPlan, err := BuildCatalogActivationPlan(fixture.Root, fixture.Ledger)
	if err != nil || len(activationPlan.FilePlan.Blockers) != 0 {
		t.Fatalf("activation plan blockers=%#v problems=%#v err=%v", activationPlan.FilePlan.Blockers, activationPlan.Problems, err)
	}
	activation, err := ExecuteCatalogActivation(activationPlan, CatalogActivationOptions{Now: time.Date(2026, 8, 9, 14, 1, 0, 0, time.UTC)})
	if err != nil || activation.Result.Status != reconcile.StatusSucceeded {
		t.Fatalf("activation=%#v err=%v", activation, err)
	}
	runGitTest(t, fixture.Root, "add", "--all")
	runGitTest(t, fixture.Root, "commit", "-m", "activate reviewed migration checkpoint")
	return fixture, activation
}
