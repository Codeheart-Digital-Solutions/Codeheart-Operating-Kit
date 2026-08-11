Last updated: 2026-08-11T14:45:48Z (UTC)
Created: 2026-08-11
Status: active
Execution log: frozen-register-canonical-migration-hotfix_execution_log.md

# Document Header

## Frozen Register Canonical Migration Hotfix Implementation Plan

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.frozen-register-canonical-migration-hotfix
  kind: implementation
  purpose: Allow complete reviewed canonical migrations to retain malformed frozen-register fields as non-authoritative historical evidence without weakening legacy or mixed authority.
  first_cataloged: 2026-08-11T14:05:49Z
  catalog_metadata_updated: 2026-08-11T14:45:48Z
  products:
    - codeheart-operating-kit
  capabilities:
    - semantic-plan-catalog
  strategic_themes:
    - ownership-aware-discovery
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.legacy-register-inventory-hotfix
    - kind: related
      target: codeheart-operating-kit.implementation.activation-ancestry-hotfix
```
<!-- END CODEHEART PLAN METADATA -->

Overview: Release v0.1.30 as a focused compatibility correction for repositories whose frozen
legacy plan register contains unsupported historical status labels or relation syntax. A complete,
high-confidence schema-v3 ledger remains the sole authority for canonical metadata. The Kit must
retain the malformed legacy fields as warnings after canonical targeting or activation, while
legacy and mixed modes continue to fail on the same evidence and every authority-critical register
problem remains blocking.

The authorized workflow includes implementation, validation, version bump, ready pull request,
expected-head merge, annotated tag, reproducible release assets, live release verification,
installation into Foundry and HQ as required, completion of the Foundry semantic migration, and an
HQ portfolio refresh. It never rewrites a frozen register or changes a consumer plan merely to make
old register syntax valid.

Essential context:

| Source | Why it matters |
| --- | --- |
| `docs/repo/runbooks/change-operating-kit.md` | Defines producer change and validation gates. |
| `docs/repo/runbooks/release-operating-kit.md` | Defines the exact release and publication workflow. |
| `docs/repo/reference/consumer-impact-classification.md` | Classifies the fix as validator-only plus managed instruction clarification. |
| `components/planning-workflows/managed/runbooks/migrate-plan-catalog.md` | Defines frozen-register, reviewed-ledger, migration, and activation authority. |
| `components/planning-workflows/managed/reference/plan-register-format.md` | Defines the frozen legacy evidence boundary. |
| `internal/plancatalog/legacy.go` | Parses historical register fields into structured problems. |
| `internal/plancatalog/view.go` | Builds target and configured catalog snapshots. |
| `internal/plancatalog/migrate.go` | Applies the fail-closed reviewed-ledger migration gates. |
| `internal/plancatalog/legacy_projection_test.go` | Protects strict legacy parser behavior. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Permit a canonical target inventory and an activated canonical catalog to treat only these frozen,
non-authoritative historical-field problems as warnings:

- `legacy_status_unsupported`;
- `legacy_relation_malformed`;
- `legacy_relation_unsupported`.

Completion requires all of the following:

- the same problem codes remain errors in configured legacy and mixed modes;
- malformed or duplicate IDs, canonical paths, dates, fields, and authority mappings remain errors
  in every mode;
- a migration still requires complete discovery-v2 evidence, a clean checkpoint, a coherent
  schema-v3 ledger, unique IDs and aliases, reviewed branch dispositions, and all existing binding
  digests;
- the frozen register is neither modified nor normalized;
- target-canonical dry-runs can reach `ready` when these three historical-field problems are the
  only remaining register problems;
- configured canonical validation retains the observations as warnings and remains idempotently
  valid;
- v0.1.30 is reproducibly published and verified before consumer installation;
- Foundry is migrated and activated from a fresh post-upgrade evidence revision;
- HQ is upgraded as required and its portfolio catalog is refreshed from current default branches.

## 1.2 Problem Context

Foundry reached an exact reviewed ledger checkpoint with 90 owned records: 78 legacy documents to
receive metadata and 12 already canonical. The guarded dry-run was otherwise coherent, but it
stopped on 29 errors originating only from frozen register status and relation text. Rewriting the
register would violate the post-cutover evidence contract, while accepting those fields as current
authority would contradict the schema-v3 ledger.

The defect is therefore contextual severity, not missing migration review. Historical parser
observations are useful and must remain visible, but those three fields no longer select identity,
path, ownership, lifecycle, or relations once the target catalog is canonical and the reviewed
ledger supplies the canonical decisions.

## 1.3 Scope And Impact

- Producer base: exact v0.1.29 release `996036893faa251e7f38c3391f94efe7f4d1273d`.
- Target release: v0.1.30, the next unused patch tag at activation.
- Primary change: validator-only contextual severity correction.
- Documentation change: instruction-only clarification of the frozen-register compatibility rule,
  mirrored through managed resources.
- Consumer upgrade effect: managed Kit and lock/version surfaces only. Upgrade must not change
  consumer plans, registers, ledgers, catalog metadata, or portfolio content.
- Migration effect: subsequent explicit migration and activation commands remain separate guarded
  consumer writes.
- Security boundary: no secrets, consumer identities, private plan text, or private topology
  identifiers enter public fixtures or release artifacts.

# Section 2 - Strategy

## 2.1 Implementation Strategy

Centralize the contextual severity decision at repository snapshot assembly, after the legacy
register has been parsed and reconciled and after the effective target mode has been selected.
Downgrade only a fixed allowlist, only when the effective mode is canonical, and only for problems
whose path is the frozen legacy register. Preserve the original code, message, path, and
remediation so operators can still inspect the historical defect.

This placement gives the same result to:

- prospective canonical inventory used by `plans migrate`;
- configured canonical `plans inventory`, `plans validate`, and compatibility views;
- local and remote-aware consumers of the shared snapshot;
- portfolio reconstruction using the same catalog semantics.

Legacy and mixed behavior is unchanged because their effective mode never enters the compatibility
projection. Migration code does not receive a broader remediable-error allowlist; it receives
warnings from the correctly contextualized inventory.

## 2.2 Validation Strategy

Add focused tests that prove:

1. the three historical-field codes are errors in legacy mode;
2. the same codes are errors in mixed mode;
3. the same codes become warnings for a target-canonical snapshot;
4. the same codes remain warnings after configured canonical activation;
5. authority-critical legacy problems remain errors in canonical mode;
6. the reviewed migration dry-run is ready when only the permitted frozen-field evidence remains;
7. existing malformed-register, mixed-cutover, migration binding, and activation tests remain
   unchanged and green.

Run focused Go tests first, then the complete Go, race, vet, Python/schema/resource/routing,
backward-compatibility, installer, release, reproducibility, and public-core matrix required by the
maintainer and release runbooks. Exact-head CI must pass on supported platforms.

## 2.3 Release And Rollout Strategy

Update every authoritative version, compatibility, manifest, fixture, checksum, graph, packaged
resource, and release-note surface through repository-owned scripts and tests. Build macOS universal
and Windows x64 assets twice and require byte identity. Merge only the reviewed exact PR head, tag
the exact merge commit with annotated `v0.1.30`, publish the exact assets, and run live installer
verification.

After publication:

1. upgrade the isolated Foundry migration branch to v0.1.30 without changing consumer-owned plan,
   register, ledger, catalog, or config bytes;
2. establish a fresh evidence revision and rebuild the reviewed ledger because managed Kit bytes
   and source revision changed;
3. apply the enumerated metadata actions, activate discovery v2, validate, and merge Foundry;
4. upgrade HQ as required through isolated Kit-only handling;
5. refresh the HQ portfolio catalog and verify that Foundry's default branch is represented.

# Section 3 - Execution Plan

## EP-01 - Plan Activation And Regression Contract

Outcome: the bounded correction has an authoritative plan checkpoint before product writes.

- [x] Reproduce the Foundry blocker from the exact reviewed dry-run and identify the 29 problem
  instances and three problem codes.
- [x] Verify v0.1.29 and current producer main contain no supported bypass or later fix.
- [x] Confirm the frozen register must not be edited and the schema-v3 ledger is complete.
- [x] Validate this plan bundle, commit only the plan, log, and nearest plan index, and push the
  checkpoint branch.

## EP-02 - Contextual Severity Correction

Outcome: canonical mode retains non-authoritative malformed historical fields as warnings without
weakening any authority-bearing gate.

- [x] Add the narrow canonical-mode compatibility projection.
- [x] Add positive canonical-target and configured-canonical regressions.
- [x] Add negative legacy, mixed, wrong-path, and authority-critical regression coverage.
- [x] Clarify the managed migration and legacy-register references and update packaged mirrors.
- [x] Run focused tests and fresh main-thread security review; resolve every High or Medium finding.

## EP-03 - v0.1.30 Release

Outcome: the correction is merged, tagged, published, and live-verified as one reproducible patch
release.

- [x] Update version, compatibility, release-note, manifest, checksum, graph, fixture, and packaged
  resource surfaces.
- [x] Run the full producer validation matrix and deterministic asset build twice.
- [ ] Push the implementation, open a ready PR, pass exact-head CI, and merge with expected-head
  protection.
- [ ] Create and push annotated `v0.1.30`; publish the exact assets and verify public checksums,
  catalog identity, CLI version, installer flow, and Kit health.

## EP-04 - Foundry Migration And HQ Refresh

Outcome: Foundry reaches canonical discovery-v2 authority and HQ observes it from the current
default branch.

- [ ] Upgrade the isolated Foundry migration branch to v0.1.30 and prove the upgrade changed no
  consumer-owned catalog content.
- [ ] Establish fresh `E`, inventory all candidates and branch evidence, and commit a rebuilt
  ledger-only `L`.
- [ ] Apply only the reviewed metadata actions sequentially; activate discovery v2 separately;
  prove idempotence and merge the exact Foundry PR head.
- [ ] Upgrade HQ as required without colliding with unrelated work.
- [ ] Refresh the HQ portfolio catalog, validate member/home state, and record final identities and
  any deliberate deferrals.

## Acceptance Contract

- The exact three non-authoritative frozen-field codes are warnings only in effective canonical
  mode and only on the legacy register path.
- No authority-critical problem changes severity.
- No migration can proceed without all existing discovery, ledger, branch, digest, worktree,
  chronology, and activation guards.
- No frozen register is rewritten by the producer fix or consumer rollout.
- v0.1.30 assets are deterministic, exact-head CI is green, and public verification passes.
- Foundry main ends in validated canonical discovery-v2 mode.
- HQ ends with a validated refresh from current member default branches, or reports an exact
  evidence-backed blocker without hiding or discarding work.

# Section 4 - Future Planning

This hotfix does not redesign the legacy parser, remove register compatibility, broaden accepted
historical syntax, migrate AWS or other consumers, or introduce automatic catalog writes during a
Kit upgrade. Any future retirement of frozen-register parsing requires separate discovery,
compatibility analysis, and a major-enough lifecycle decision.

# Revision Notes

- 2026-08-11: Activated the focused successor plan after Foundry's complete reviewed migration
  exposed a canonical-mode severity defect in v0.1.29.
