Last updated: 2026-08-10T19:08:24Z (UTC)
Created: 2026-08-10
Status: active
Execution log: activation-ancestry-hotfix_execution_log.md

# Activation Ancestry Resolution Hotfix Implementation Plan

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.activation-ancestry-hotfix
  kind: implementation
  purpose: Accept safely incorporated activation checkpoints across normal merge topology while preserving exact migration chronology and fail-closed evidence.
  first_cataloged: 2026-08-10T19:08:24Z
  catalog_metadata_updated: 2026-08-10T19:08:24Z
  products:
    - codeheart-operating-kit
  capabilities:
    - semantic-plan-catalog
  strategic_themes:
    - ownership-aware-discovery
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.branch-ownership-reconciliation
```
<!-- END CODEHEART PLAN METADATA -->

## Outcome

Ship v0.1.29 as a focused validator-only correction. Persisted discovery-v2 activation evidence
must remain valid when the exact reviewed activation checkpoint is incorporated by a normal merge
commit or retained in a later descendant, including when that checkpoint is not on the current
branch's first-parent chain. Missing, rewritten, unrelated, stale, ambiguous, or unincorporated
evidence must continue to fail closed.

The authorized workflow includes implementation, validation, a focused pull request, exact-head
merge, annotated patch tag, reproducible public release assets, live verification, and isolated
Kit-only consumer upgrades. Consumer plan, register, ledger, catalog, and config bytes are not
migration targets for this correction.

## Authority And Impact

- Base: exact released v0.1.28 merge `38e9b85180c94f06dcddf2fd69c14db6bae85996`.
- Branch: `codex/activation-ancestry-v0.1.29`.
- Producer authority: `internal/plancatalog`, `internal/commands`, `internal/portfolio`, tests,
  release surfaces, and this plan bundle.
- Consumer impact: `validator-only change`; release notes are required; no plan-catalog migration,
  repair, sync, reactivation, register rewrite, or consumer metadata change is required.
- Public-core boundary: committed fixtures use synthetic repositories and generic revisions only.
  Live rollout evidence stays outside public artifacts except for public repository release and PR
  identities intentionally required by the release workflow.

## Security Invariant

The validator resolves one immutable activation checkpoint `A` without weakening chronology:

1. `L^1 = E`; `E..L` contains only the reviewed ledger path and exact ledger blob.
2. `A^1 = L`; `L..A` contains exactly the guarded config plus reviewed migration actions, and its
   action digest and blobs match the persisted binding.
3. `A` must be reachable from current `HEAD`; current `HEAD` may equal `A`, be a linear descendant,
   or incorporate `A` through any merge parent.
4. Resolution searches the complete `L..HEAD` ancestry path for direct first-parent children of
   `L`, filters by the bound migration action digest, and requires one unambiguous match.
5. Current config, candidate authority, deferred-owner integration, remote overlay, index, and
   worktree checks remain unchanged after `A` is resolved.

Tree equality alone is never incorporation proof. A matching tree on an unrelated sibling,
reparented commit, or branch that is not an ancestor of current `HEAD` cannot satisfy the binding.

## EP-01 - Topology Regression And Generic Resolver

Outcome: local and remote-aware validation resolve the exact incorporated activation checkpoint
across merge graphs while rejecting absent or ambiguous evidence.

- [ ] Add synthetic `E -> L -> A` fixtures for exact `A`, a normal merge with `A` as second parent,
  and safe later descendants.
- [ ] Add negative fixtures for unreachable `A`, reparented/tree-only coincidence, wrong ledger,
  unrelated sibling, action-digest mismatch, and ambiguous matching activation children.
- [ ] Replace first-parent-from-`HEAD` discovery with one shared ancestry-path resolver.
- [ ] Use the resolver consistently in persisted local validation, command-side remote target
  reconstruction, portfolio target reconstruction, and compatibility projection.
- [ ] Preserve exact checkpoint parent/delta/blob/config/action and current-authority guards.

## EP-02 - Compatibility And Release Proof

Outcome: the correction is released as v0.1.29 with no migration and no unrelated behavior change.

- [ ] Run focused plan-catalog, command, and portfolio suites first.
- [ ] Run full Go, race, vet, Python/schema/resource/routing/release, backward-compatibility,
  installer, and transaction validation.
- [ ] Complete fresh read-only review and resolve every High or Medium finding.
- [ ] Update all authoritative v0.1.29 version, compatibility, manifest, fixture, and release-note
  surfaces.
- [ ] Build macOS universal and Windows x64 assets twice, verify byte equality, sidecars, catalog,
  platform identity, and catalog-to-binary chain.
- [ ] Push a focused PR, pass exact-head CI, merge with the expected-head guard, create an annotated
  v0.1.29 tag, publish the exact assets, and pass live public installer and post-publication gates.

## EP-03 - Isolated Consumer Confirmation

Outcome: representative activated repositories adopt only the new Kit release and validate their
existing plan-catalog state without semantic migration work.

- [ ] Use clean isolated branches from exact live default branches.
- [ ] Preview and apply only the v0.1.29 Kit upgrade; preserve consumer-owned plan, register,
  ledger, catalog, and config bytes.
- [ ] Validate Kit health and repository-specific plan-catalog state after upgrade.
- [ ] Keep further plan metadata migration and coordination refresh outside this corrective task.

## Acceptance Contract

- Exact `A`, normal merges incorporating `A`, and later descendants validate identically when
  bound authority has not moved.
- Absent, unreachable, reparented, wrong-ledger, unrelated, digest-mismatched, or ambiguous
  activation evidence produces a stable blocking result.
- Local, command-side remote-aware, and portfolio reconstruction use the same topology semantics.
- Existing linear chronology, transaction, ownership, candidate, overlay, compatibility, and
  no-write tests remain green.
- v0.1.29 is reproducibly merged, tagged, published, live-verified, and installed through isolated
  Kit-only consumer changes without plan-catalog migration.
