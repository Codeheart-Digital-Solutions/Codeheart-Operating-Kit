Last updated: 2026-08-09T16:34:26Z (UTC)

# Branch Ownership Reconciliation Consumer Impact Record

## Classification

| Impact class | Applicability | Reason |
| --- | --- | --- |
| `instruction-only change` | yes | Managed planning, migration, register, portfolio, and routing guidance describes reviewed branch clearance and mixed deferral. |
| `validator-only change` | yes | Inventory, validation, migration, activation, and portfolio scans recompute hash-bound branch evidence and fail closed on drift. |
| `consumer migration required` | only for repositories choosing the new migration route | A Kit upgrade does not activate discovery v2, change catalog mode, write a ledger, or migrate a plan. Repositories with branch-touch blockers must separately inventory, review, checkpoint, migrate, and activate. |
| `security or safety policy change` | yes | Clearance is proof-specific, candidate-scoped, ref-tip-bound, and transaction-rechecked. Missing, stale, ambiguous, or tampered evidence blocks. |
| `breaking placement-contract change` | no | No managed or consumer-owned path moves. Historical schema contracts remain available under versioned schema files. |
| `component addition` | no | No component or profile-selected component is added. |
| `backwards-compatible scaffold addition` | no | No consumer-owned scaffold or generated consumer path is added. |

## Affected Surfaces

- Plan-catalog evidence, inventory, migration, activation, view, CLI, and transaction behavior under
  `internal/`.
- Remote mirror, overlay, cache, compatibility-observation, and readiness behavior under
  `internal/portfolio/`.
- Strict config v2, migration-ledger v3, inventory v3, and portfolio-catalog v3 schemas, with
  historical config v1, ledger v2, and catalog v2 schemas retained.
- Managed planning-workflow and routing doctrine plus exact packaged resource mirrors.
- Release validation on Ubuntu, macOS, and Windows, including old-CLI compatibility probes.

## Compatibility And Consumer Action

Operating Kit v0.1.25 and v0.1.26 consumers retain their existing config-v1 and ledger-v2
behavior. Those CLIs reject schema-v3 ledgers and schema-v2 configs before writes. The new CLI
continues to load and execute existing schema-v1/v2 ledgers under their historical semantics.

A repository using reviewed clearance must create a schema-v3 ledger at a dedicated ledger-only
checkpoint. A migration apply remains separate from catalog activation. Mixed activation with a
genuine deferred owner persists a strict config-v2 binding to the evidence revision, ledger bytes,
activation base, branch evidence, action digest, and evidence scope. Direct canonical activation
remains zero-gap and uses config v1 because no durable deferral binding is needed.

Remote-aware review requires a fresh complete overlay. A matching local tracking observation and
source-stable mirror observation coalesce to one authority identity; disagreement is retained as a
blocking diagnostic. Optional merged-PR evidence remains corroboration only.

## Safety And Rollback

- Same-content clearance requires exact mode, object ID, and SHA-256 equality at the reviewed
  target and ref tip.
- Incorporated-history clearance requires a complete candidate transition plus stable patch
  equivalence to exactly one reachable target commit; partial collisions never clear.
- A deferred owner is limited to one exact active candidate and one mandatory bound incremental
  follow-up. Integration resolves only when the reviewed owner bytes and canonical metadata match.
- Ref, merge-base, path state, attribute, policy, candidate-set, source, target, ledger, config,
  action, and requested overlay drift invalidate authority.
- Pre-write and post-write authority checks preserve transaction rollback behavior, including
  no-op execution.

## Release And Rollout Requirements

- Consumer-facing release notes and migration guidance: required.
- Versioned machine-contract disclosure: required.
- Two byte-identical release builds, archive sidecars, cross-platform CI, and live release
  verification: required before handoff.
- HQ, Foundry, AWS, and other consumer upgrades remain separate isolated tasks after publication.
  This producer change performs no consumer install or semantic migration.

## Residual Boundary

The proof engine intentionally depends on Git 2.43 or newer. Tooling that does not meet that floor
must follow the managed tooling-readiness route; the CLI must not downgrade or bypass proof
semantics. Public examples and tests remain synthetic and contain no consumer repository evidence.
