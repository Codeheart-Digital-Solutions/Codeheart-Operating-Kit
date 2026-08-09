Last updated: 2026-08-09T14:33:58Z (UTC)
Created: 2026-08-09
Status: active
Execution log: branch-ownership-reconciliation_execution_log.md

# Document Header

## Evidence-Backed Branch Ownership Reconciliation And Active-Plan Deferral Implementation Plan

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.branch-ownership-reconciliation
  kind: implementation
  purpose: Implement hash-bound branch-owner reconciliation and safe mixed-mode deferral for semantic plan-catalog migration.
  first_cataloged: 2026-08-09T13:31:33Z
  catalog_metadata_updated: 2026-08-09T14:33:58Z
  products:
    - codeheart-operating-kit
  capabilities:
    - semantic-plan-catalog
    - portfolio-coordination
  strategic_themes:
    - ownership-aware-discovery
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.ownership-aware-plan-catalog-discovery
    - kind: related
      target: codeheart-operating-kit.implementation.semantic-plan-catalog-coordination
```
<!-- END CODEHEART PLAN METADATA -->

Overview: Implement a deterministic, reviewed reconciliation layer between discovery-v2 branch
touch detection and migration ownership enforcement. The layer must distinguish an actual active
plan owner from a same-content historical ref and from a divergent ref whose complete candidate
transition was already incorporated into a commit reachable from the reviewed target revision.
It must do so without deleting refs, weakening candidate coverage, trusting pull-request state as
sole evidence, or bypassing migration hashes and transaction preconditions.

The implementation also permits one exact genuine owner per deferred candidate to remain
filename-only through the existing mixed-mode cutover contract. Such a deferral is valid only when the plan remains at its
frozen cutover path and bytes on the target revision, the owner ref and active plan bytes are
bound, every other touching ref is reconciled, and the reviewed ledger contains a mandatory
incremental-migration follow-up. Direct canonical mode remains zero-gap. A guarded catalog
activation command must revalidate the exact ledger, local refs, requested remote overlay, and
candidate projection before it writes config, then persist a strict hash binding to the evidence
while any deferral remains. Discovery-v2 activation in mixed mode must keep the deferred plan
visible as explicit compatibility evidence; it must never represent the repository as
canonical-ready while that follow-up remains.

This document is planning authority only. It does not authorize implementation edits, dependency
changes, a version bump, release-note edits, a commit, push, pull request, tag, release, consumer
upgrade, or any write in another repository. The current worktree was verified clean and detached
at `7a6cd54b27b3f0d322d6b65e436542c677fa7e95`, which is simultaneously `origin/main` and the
`v0.1.26` release tag. No planning branch is required while this plan remains draft.

Essential context:

| Source | Why a future implementer must read it |
| --- | --- |
| `AGENTS.md` | Producer authority, public-core safety, managed-resource rules, and repository change boundaries. |
| `docs/repo/runbooks/change-operating-kit.md` | Consumer-impact, source/mirror parity, transaction, routing-probe, and validation gates. |
| `docs/repo/runbooks/release-operating-kit.md` | Separate version, release-candidate, cross-platform, reproducibility, publication, and rollback gates. |
| `components/planning-workflows/managed/reference/plan-catalog-format.md` | Discovery v2, candidate authority, mixed/canonical modes, branch overlays, and activation contract. |
| `components/planning-workflows/managed/reference/planning-document-lifecycle.md` | Lifecycle, plan bundle, attachment, index, and mixed-view behavior. |
| `components/planning-workflows/managed/reference/portfolio-coordination-format.md` | Remote ref, cache completeness, stale evidence, optional pull-request facts, and overlay authority. |
| `components/planning-workflows/managed/runbooks/migrate-plan-catalog.md` | Current inventory, schema-v2 ledger, branch-owner stop, mixed grandfathering, apply, activation, and recovery sequence. |
| `components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Frozen register and mixed compatibility behavior; the register must not be edited by this feature. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | `planning.migrate-catalog` route, approval boundary, evidence contract, and low-context probe requirement. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Agent-facing runbook intention, source, lane, precondition, stop, evidence, validation, and recovery shape. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | L1 migration orchestration, L3 command wrappers, structured blockers, validation tiers, and non-promotion decision. |
| `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` | Approval-gated recovery when the required Git version or another local validation tool is missing; provider service preflight remains portfolio-owned. |
| `internal/plancatalog/inventory.go` | Current local/tracking-ref enumeration and ref-name-only active-branch touch evidence. |
| `internal/plancatalog/migrate.go` | Current schema-v1/v2 loading, branch-owner ordering defect, mixed deferral, authority digest, projection, and transaction hook. |
| `internal/plancatalog/view.go` | Current config decoding, frozen cutover validation, mixed plan-row preservation, and structured local views. |
| `internal/plancatalog/classifier.go` and `git_index.go` | Candidate policy, Git-blob provenance, digests, safe Git object reads, and path normalization. |
| `internal/commands/plans.go` | `validate`, `list`, `inventory`, and `migrate` CLI flags, JSON output, remote overlays, and activation=false result. |
| `internal/portfolio/scanner.go`, `mirror.go`, `source.go`, and `store.go` | Remote ref enumeration, merge-base diffing, optional PR enrichment, cache schema handling, and atomic preservation. |
| `internal/reconcile/plan.go` and `transaction.go` | Exact file preconditions, staged apply, authority recheck, post-check, rollback, and recovery. |
| `schemas/plan-migration-ledger.schema.json` | Strict current ledger v2 contract that must be preserved as historical behavior. |
| `schemas/plan-catalog.schema.json` | Strict current portfolio catalog/cache v2 contract and compatibility boundary. |
| `schemas/kit-config.schema.json` and `manifest.yaml` | Released config-v1 contract and declared config compatibility boundary. |
| `.github/workflows/validate.yml` | Current Linux, macOS, Windows, Go, Python, schema, routing, packaging, and release validation lanes. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Deliver evidence-backed branch reconciliation and mixed-mode active-owner deferral as an additive
plan-catalog capability whose success is measurable by all of the following outcomes:

- prospective inventory emits one structured, deterministic touch record for every relevant local
  branch, tracking ref, and requested remote-overlay ref, including the exact ref tip, unique merge
  base, target revision, candidate paths, Git modes, blob object IDs, content SHA-256 values, and
  proof inputs;
- a strict reviewed ledger v3 classifies every touch as `same-content-non-owner`,
  `incorporated-history-non-owner`, `deferred-active-owner`, `active-owner`, or `blocking`;
- same-content clearance requires exact same-path, same-mode, same-blob, and same-SHA-256 identity
  at the reviewed target revision and touching ref tip;
- incorporated-history clearance requires a unique merge base, a specific incorporated commit and
  parent reachable from the target revision, identical ordered path-state transitions, and equal
  aggregate `git patch-id --stable` output; patch ID never clears a ref without exact state proof;
- optional merged-pull-request evidence remains corroboration and never replaces offline Git
  object, ancestry, path-state, blob, hash, and patch evidence;
- every reviewed ref, candidate, policy, source revision/hash, path state, target precondition,
  cutover hash, and evidence algorithm is recomputed during planning, immediately before
  transaction commit, after writes before transaction cleanup, and again by guarded activation,
  including zero-action/idempotent application;
- evidence revision `E`, ledger-only checkpoint `L`, and activation checkpoint `A` have
  non-self-referential first-parent/tree-delta contracts, with mandatory clean validation
  immediately after `A`;
- one genuine active owner per candidate can remain legacy only in an already-valid mixed catalog with exact
  frozen register/cutover proof, one active owner ref, a bound ref-tip candidate whose lifecycle is
  `active`, and a structured incremental follow-up;
- every other touching ref for that plan is proven non-owning, while multiple active owners,
  deletion-only owners, dirty/uncommitted owners, missing refs, ambiguous merge bases, and stale
  proof inputs remain blockers;
- mixed local and remote structured views retain explicit compatibility observations for deferred
  plans, report canonical readiness separately from scan completeness, and keep the plan visible;
- a later owner-branch merge that changes the grandfathered target bytes causes validation to fail
  until the owner supplied metadata or a reviewed incremental migration completes;
- direct canonical mode rejects every active-owner deferral and retains zero owned gaps;
- ledgers v1 and v2 and configs v1 retain their released behavior, while v0.1.25 and v0.1.26
  reject ledger v3 and config v2 before any mutation;
- dry-run/apply chronology, exact source preservation, activation separation, atomic rollback,
  remote-cache preservation, public-safe evidence, and idempotency remain intact; and
- the complete unit, CLI, schema, fixture, routing, packaging, security, Linux, macOS, and Windows
  matrix passes before a separately authorized release candidate is presented.

## 1.2 Project And Problem Context

### Evidence-backed problem

The v0.1.26 discovery-v2 migrator correctly treats branch ownership as a safety boundary, but its
current evidence is too coarse. `activeBranchTouchesWithSettings` enumerates local and tracking
refs, skips only refs whose tips are ancestors of `HEAD`, and records a path as touched from a
three-dot `HEAD...ref` diff. This answers whether the ref changed the candidate after its merge
base; it does not answer whether the ref tip now has the same candidate blob as the target, or
whether the ref's complete candidate patch was already incorporated through squash/merge history.

The inventory stores only ref names. The migration authority digest consequently binds the set of
touching ref names but not their tip commits, merge bases, candidate blobs, or proof results. A ref
can move while retaining the same touched path set without necessarily changing that digest.

The runtime ordering creates a second, independent defect. A schema-v2 record is rejected when
`branch_owner != "none"`, and any inventory `active_branch_touches` entry is rejected, before the
mixed-mode `deferred` contract is evaluated. The existing grandfathering contract is strong once
reached—already-active mixed mode, frozen register, same path, exact cutover bytes, filename-only
record—but a genuine owner can never reach it.

### Sanitized external acceptance scenarios

These scenarios are acceptance inputs only. Implementation and validation must use synthetic,
public-safe fixtures; no consumer repository is read or edited by this plan.

| Scenario | Verified evidence to reproduce generically | Required projection |
| --- | --- | --- |
| Scenario A | Historical and same-content branch refs coexist with two genuinely active owner refs. | Clear only proven non-owners, retain both active plans as exact mixed deferrals, preserve their rows, and require later incremental migration. |
| Scenario B | Nineteen branch-touch candidates comprise twelve same-content non-owners, six divergent refs whose exact candidate blobs and aggregate stable patch IDs match incorporated commits reachable from the target revision, and one genuine active owner. | Produce 18 reviewed non-owner clearances and one mixed deferral without deleting any ref, while direct canonical projection remains blocked by the active owner. |

### Goals

- Reconcile ref ownership from deterministic Git evidence instead of branch age, naming, or PR
  state.
- Keep historical branches valuable and visible without treating their continued existence as
  active ownership.
- Preserve branch-owner safety when evidence is incomplete, stale, ambiguous, tampered, or moved.
- Make mixed mode a safe incremental migration lane for exact active plans.
- Preserve canonical mode as a complete metadata catalog with no grandfathered gaps.
- Make offline evidence useful and honest, and make requested remote-overlay completeness a
  separately bound gate.
- Keep all evidence public-safe, stable in JSON/YAML, and reproducible across supported platforms.

### Non-goals

- Do not delete, prune, rename, merge, close, or rewrite historical or active branches as the
  primary solution.
- Do not infer non-ownership from branch age, stale labels, naming, a merged PR, a closed PR, or a
  missing provider response.
- Do not accept stable patch ID equality without exact path-state and blob/hash equality.
- Do not clear partial branch patches, split incorporation across multiple target commits, or
  ambiguous criss-cross merge bases in the first implementation.
- Do not migrate dirty worktree bytes, untracked plan authority, symlinks, gitlinks, unsafe paths,
  or another checkout's uncommitted work.
- Do not add a safety-bearing key to config schema v1. Its permissive
  `component_settings.planning-workflows` object would let v0.1.25/v0.1.26 silently ignore the
  key. Preserve config v1 and add a strict config v2 only for repositories that activate a
  deferral-bound mixed catalog.
- Do not edit the frozen plan register, add a second cutover revision, activate discovery v2
  automatically, or make `plans migrate` write shared config.
- Do not add an L2 script wrapper around the existing L3 plan commands.
- Do not operate consumer repositories, bump versions, edit release surfaces, publish, or install
  a release during this planning task.

### Consumer impact classification

The eventual producer change requires one consumer-impact record with these classes:

- `instruction-only change`: managed plan-catalog, migration, register, portfolio, routing, and
  lifecycle doctrine changes;
- `validator-only change`: new evidence, disposition, compatibility observation, canonical
  readiness, and stale-proof validation;
- `consumer migration required`: consumers must upgrade before creating or applying ledger v3 and
  must retain the reviewed ledger for mixed activation evidence;
- `security or safety policy change`: branch ownership can be cleared only through the new
  cryptographic and Git-history proof contract.

The placement contract does not change: plan documents, plan-scoped ledgers, inventories, config,
local cache, and managed doctrine remain in their existing ownership areas. No scaffold addition
is planned.

## 1.3 Current State Analysis

### Released baseline and compatibility

- `v0.1.26`, `origin/main`, and the inspected worktree all resolve to
  `7a6cd54b27b3f0d322d6b65e436542c677fa7e95`.
- The relevant ledger schema, inventory, migration runtime, config schema, and plan command files
  are byte-identical between v0.1.25 and v0.1.26.
- Ledgers v1 and v2 use strict `additionalProperties: false` schemas. Ledger loading validates the
  raw YAML map and then decodes with `KnownFields(true)`.
- Old CLIs support only ledger schema versions 1 and 2, so schema version 3 already fails closed.
- The repository config is top-level schema v1; portfolio contains a nested schema-v2 block. Old
  CLIs validate any config against the schema-v1 constant, so a top-level config schema v2 fails
  closed before settings are used. The implementation must preserve current bytes as
  `kit-config-v1.schema.json`, make the unsuffixed schema the strict v2 contract, and explicitly
  dispatch both versions in new CLIs.

### Current implementation map

| Surface | Current behavior | Required target behavior |
| --- | --- | --- |
| Branch touch discovery | `internal/plancatalog/inventory.go` records `map[path][]ref-name` from unmerged local/tracking refs and a three-dot diff. | Emit versioned per-ref evidence with tip, unique merge base, transition states, object IDs, hashes, proof scope, and a digest. |
| Same-content refs | A ref can remain a touch even when its tip candidate blob equals the target candidate blob. | Classify exact same-path/mode/blob/hash identity as reviewable non-owner evidence. |
| Squash/merge history | No `patch-id`, historical reachable-commit, or exact transition comparison exists. | Match one complete candidate transition to one exact reachable incorporated commit transition, then require stable patch equality as a second proof. |
| Stale tracking refs | Age and provider state do not affect the local hard-owner result. | Age never clears; exact Git proof may clear. Requested fresh remote evidence must agree with the reviewed snapshot. |
| Ref movement | Authority digest binds ref names, not tip/object state. | Bind and recheck every tip, base, path state, blob, hash, and proof result before commit, including idempotent apply. |
| Ledger schema | V2 contains flat `branch_owner`, free-text evidence, and a boolean deferral. | Preserve v2; add strict v3 per-ref reviews, evidence policy/digest, conditional proof fields, and structured follow-up. |
| Migration ordering | `branch_owner` and discovered touches block before deferral. | Reconcile every touch first; allow only proof-complete non-owner clearances and one exact mixed active-owner deferral. |
| Mixed deferral | Exact frozen cutover proof exists but only `branch_owner: none` can reach it; follow-up is free text. | Preserve frozen proof, require a real active owner and exact ref bytes, and bind a mandatory incremental follow-up. |
| Local views | Mixed filename-only records remain visible through legacy register reconciliation. | Add explicit `mixed-grandfathered` coverage and `incremental_migration_required` fields; never call the view canonical-ready. |
| Remote catalog | Metadata-less mixed records are omitted from canonical observations and remain only low-level candidates. | Add explicit compatibility observations, preserve scan completeness separately, and expose `canonical_ready: false`. |
| Remote overlays | Inventory may attach a scan, but migration rebuilds only local evidence. | Add migration options that bind a requested complete overlay snapshot without making online PR facts a prerequisite. |
| Activation | Metadata migration intentionally reports `activation_performed: false`, but no supported command binds the later config change back to the reviewed branch evidence. | Add a config-only guarded catalog activation that rechecks the exact ledger/ref/overlay/projection state and persists its digest while deferrals remain. |
| Transaction | Action plans recheck an authority digest; the no-action fast path can return before the transaction callback, and ref authority can move after the precommit check. | Recheck v3 evidence before every result, before writes, after the action batch, and immediately before cleanup; rollback on any post-write drift. |
| Cache | Catalog v2 is strict, but cached JSON is decoded into a struct before validation, which can discard unknown fields. | Preserve v2, add catalog v3, and validate raw JSON before versioned decode. |

### Existing proof foundations to preserve

Current tests already prove branch deletion and rename detection, dirty/hash blocking, exact target
preconditions, mixed cutover integrity, separate activation, migration idempotency, transaction
rollback, remote pushed-only observations, prior-cache preservation, PR enrichment without ref
selection, adversarial branch inertness, executable Markdown preservation, CRLF handling, and
Linux/macOS/Windows execution. The implementation extends these suites; it does not replace them.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Expected producer changes:

```text
schemas/
├── kit-config-v1.schema.json                            # create: frozen released config v1
├── kit-config.schema.json                               # modify: strict v2 with evidence binding
├── plan-migration-ledger-v2.schema.json                 # create: frozen released v2 contract
├── plan-migration-ledger.schema.json                    # modify: latest strict v3 contract
├── plan-inventory.schema.json                           # create: strict reviewed inventory v3
├── plan-catalog-v2.schema.json                          # create: frozen released catalog v2
└── plan-catalog.schema.json                             # modify: catalog/cache v3 compatibility observations
internal/
├── state/
│   ├── schema.go                                        # modify: explicit config/catalog/ledger dispatch
│   └── observed.go                                      # modify: raw-version config validation
├── lockfile/lockfile.go                                 # modify: config-version-aware validation
├── plancatalog/
│   ├── branch_reconciliation.go                         # create: ref snapshots, transitions, proof engine, digest
│   ├── branch_reconciliation_test.go                    # create: deterministic evidence unit matrix
│   ├── migration_v3_test.go                             # create: ledger v3 migration and deferral matrix
│   ├── model.go                                         # modify: evidence, dispositions, compatibility rows
│   ├── git_index.go                                     # modify: safe object/path-state reads
│   ├── inventory.go                                     # modify: structured local/tracking touch evidence
│   ├── migrate.go                                       # modify: explicit version dispatch and v3 reconciliation
│   ├── view.go                                          # modify: config dispatch, binding validation, compatibility
│   ├── activation.go                                    # create: guarded catalog-activation plan
│   ├── migration_v2_test.go                             # modify: frozen v2 compatibility assertions
│   ├── ep02_test.go                                     # modify: regression, atomicity, dirty, and rollback coverage
│   └── classifier_test.go                               # modify: rename/delete and path-state coverage
├── commands/
│   ├── plans.go                                         # modify: v3 UX, overlays, guarded catalog activation
│   ├── portfolio.go                                     # modify: config-version-aware validation
│   └── commands_test.go                                 # modify: CLI, output, redaction, and approval coverage
├── cli/
│   ├── cli.go                                           # modify: help for branch reconciliation fields/options
│   └── cli_test.go                                      # modify: stable command/help contract
├── portfolio/
│   ├── config.go                                        # modify: config-version-aware decoding
│   ├── membership.go                                    # modify: config-version-aware decoding
│   ├── source.go                                        # modify: remote touch and compatibility observations
│   ├── scanner.go                                       # modify: remote proof snapshots and raw cache validation
│   ├── mirror.go                                        # modify: exact ref/base/path transition reads
│   ├── store.go                                         # modify: v2 historical and v3 current cache handling
│   └── portfolio_test.go                                # modify: remote proof, stale, cache, and visibility matrix
└── reconcile/
    ├── plan.go                                          # modify: split authority callback contract
    ├── transaction.go                                   # modify: zero-action/pre/post-write checks
    └── reconcile_test.go                                # modify: drift injection and rollback matrix
components/
├── planning-workflows/
│   ├── component.yaml                                   # modify: changed managed-resource identities
│   └── managed/
│       ├── README.md                                    # modify: reviewed reconciliation discoverability
│       ├── reference/
│       │   ├── plan-catalog-format.md                   # modify: evidence/status/visibility contract
│       │   ├── planning-document-lifecycle.md           # modify: mixed follow-up lifecycle
│       │   └── portfolio-coordination-format.md         # modify: remote proof and readiness semantics
│       └── runbooks/
│           ├── migrate-plan-catalog.md                  # modify: reviewed clearance/deferral procedure
│           └── maintain-plan-register.md                # modify: visibility/follow-up without register edits
└── agent-interface/managed/reference/
    └── operation-routing-and-dispatch.md                # modify: migration preconditions/evidence/blockers
src/codeheart_operating_kit/resources/components/        # modify: byte-identical managed resource mirrors
tests/
├── fixtures/plans/                                      # modify: local evidence and migration fixtures
├── fixtures/portfolio/                                  # modify: remote branch and cache fixtures
├── test_json_schemas.py                                 # modify: v2/v3 positive and adversarial schemas
├── test_routing.py                                      # modify: clearance route and authority negatives
├── test_packaging_resources.py                          # modify: source-to-packaged byte equality
├── test_install_metadata.py                             # modify later: config-v1/v2 compatibility declaration
├── test_go_cli_parity.py                                # modify later: packaged command/schema parity
└── test_release_assets.py                               # modify later: config-schema packaging
.github/workflows/validate.yml                           # modify: deterministic cross-platform proof gates
docs/repo/plans/branch-ownership-reconciliation/         # current plan and later plan-scoped evidence
├── branch-ownership-reconciliation_implementation_doc.md
├── branch-ownership-reconciliation_execution_log.md     # create only after implementation activation
└── attachments/                                         # create during execution: sanitized impact/release evidence
release surfaces                                         # modify only in the separately authorized release epic
├── pyproject.toml
├── src/codeheart_operating_kit/__init__.py
├── internal/version/version.go
├── manifest.yaml
├── profiles/standard.yaml
├── bootstrap.md
├── install.sh
├── install.ps1
└── release-notes.md
```

No production script asset is added. The existing `plans inventory`, `plans migrate`,
`plans validate`, `plans list`, and `portfolio scan` L3 commands remain executable boundaries, and
one new L3 `plans catalog-activate` command owns the config-only activation transaction. The
migration runbook remains the L1 operator entry point.

## 2.2 Open Questions And Assumptions Requiring Clarification

### OQ-01 - Persisted deferral evidence binding

- `BLOCKER: no`
- `Affects: EP-03, EP-04, EP-05`
- Question: Must shared config bind the reviewed ledger after mixed activation?
- Decision: Yes. Config schema v2 adds a strict
  `component_settings.planning-workflows.plan_catalog_migration_evidence` object containing the
  tracked ledger path/hash, branch-evidence digest, evidence revision, activation-base revision,
  migration-action digest, and evidence scope.
  It is required whenever discovery v2 mixed mode contains a deferred active owner and forbidden
  when no follow-up remains.
- Rationale: otherwise a ref, overlay, or candidate can move after metadata apply and before the
  separate config change, and later validation cannot prove which approval justified the gap.
  Adding the field to config v1 is unsafe because released CLIs would ignore it; top-level config
  v2 makes v0.1.25/v0.1.26 reject the entire config before use.
- Closure: D-06 and EP-03 define guarded activation and ongoing validation; no further design
  choice is required.

### OQ-02 - Multi-commit incorporated history

- `BLOCKER: no`
- `Affects: EP-02, EP-03, EP-06`
- Question: Should one historical ref transition match a sequence of reachable target commits?
- Recommended default: No for the first implementation. Require one exact incorporated commit and
  parent whose complete scoped transition matches the ref aggregate.
- Rationale: the verified scenarios are satisfied by a specific reachable commit, while sequence
  partitioning introduces patch ordering and subset ambiguity.
- Reopen when: a reviewed consumer case proves exact incorporation split across commits and a
  collision-safe whole-sequence contract is approved.

### OQ-03 - Remote-overlay requirement by repository

- `BLOCKER: no`
- `Affects: EP-02, EP-04, EP-06`
- Question: Must every migration perform a live remote scan?
- Recommended default: No. Local and tracking refs form a complete offline evidence lane for the
  current checkout. A ledger created with requested remote-overlay evidence requires the same
  complete overlay scope at dry-run and apply. Offline results state that remote freshness and
  portfolio completeness remain unclaimed.
- Rationale: PR/provider facts cannot become an offline prerequisite, while configured activation
  and portfolio completeness still require their normal remote gates.

### OQ-04 - Release version

- `BLOCKER: no`
- `Affects: EP-07`
- Question: Which version carries the eventual change?
- Recommended default: Recheck the live release baseline at EP-07 and select the next available
  patch. From the verified v0.1.26 baseline, the current candidate is v0.1.27.
- Rationale: the feature is additive and fail-closed, existing configs and ledger v2 behavior are
  preserved, and repository history delivers comparable safety features as patch releases.

### Assumptions

- `A-01`: Git provides `patch-id --stable` on every supported validation platform. EP-02 must
  define the supported Git-version range from versions proven by the release matrix; any version
  outside that recorded range fails `branch_evidence_git_version_unsupported` before proof.
- `A-02`: fully qualified refs and selected commit objects exist locally for offline proof; lazy
  fetch remains disabled.
- `A-03`: each incorporated-history proof in the verified acceptance set maps to one non-root
  commit reachable from the target revision.
- `A-04`: plan-scoped inventory and ledger artifacts remain public-safe and contain no private
  branch names in committed producer fixtures.
- `A-05`: consumer rollout supplies its own sanitized ledger and never copies raw private evidence
  into this public repository.

## 2.3 Architectural Decisions With Reasoning

### D-01 - Preserve ledger v2 and introduce strict ledger v3

1. Problem: adding fields to the released v2 contract would relabel safety semantics and make old
   artifact behavior unclear.
2. Simplest working solution: copy the current schema to
   `plan-migration-ledger-v2.schema.json`, promote the unsuffixed latest schema to version 3, add
   explicit version dispatch, and preserve v1/v2 runtime paths unchanged.
3. Six-to-twelve-month pressure: later proof kinds may need their own version without weakening
   historical ledgers.
4. Rationale: old v0.1.25/v0.1.26 CLIs reject version 3 before mutation; new CLIs can reproduce
   current v2 blocking behavior exactly.
5. Alternatives rejected: mutating v2 in place; interpreting free-text `evidence`; adding a config
   escape hatch; accepting unknown fields.

### D-02 - Separate repository ownership from per-ref branch disposition

1. Problem: `ownership_disposition: owned` describes repository path authority, while flat
   `branch_owner` conflates that authority with ref-specific migration responsibility.
2. Simplest working solution: retain repository ownership and add a closed
   `branch_touch_reviews` array containing one review for every inventory touch.
3. Six-to-twelve-month pressure: multiple refs and remote scopes need independent review and
   invalidation.
4. Rationale: every ref can be reconciled without weakening the candidate's repository-owned
   status.
5. Alternatives rejected: a single branch owner string; one top-level clearance applying to all
   refs; automatic clearance without a reviewed ledger.

The disposition vocabulary is:

| Disposition | Runtime meaning |
| --- | --- |
| `same-content-non-owner` | Reviewable clearance only after exact target/ref tip candidate identity. |
| `incorporated-history-non-owner` | Reviewable clearance only after exact transition and reachable-commit patch proof. |
| `deferred-active-owner` | One genuine active owner preserved through exact mixed-mode grandfathering and mandatory follow-up. |
| `active-owner` | Genuine owner that is not approved for deferral; migration remains blocked. |
| `blocking` | Missing, stale, ambiguous, conflicting, unsafe, unsupported, or failed proof; migration remains blocked. |

Inventory may propose a classification, but only the reviewed ledger records the disposition used
by migration.

The review identity is exactly the four-tuple `(evidence_scope, repository_id, logical_ref,
candidate_path)`. `candidate_path` is the default-revision authoritative path; rename counterparts
belong in that review's ordered `path_transitions` and do not create a second review identity.
Two candidates touched by the same ref therefore have separate reviews, proof scopes, and failure
outcomes. Runtime requires a one-to-one set equality between inventory touch identities and ledger
review identities; it never lets one branch-level decision clear unrelated plans.

The following is the normative v3 shape. Angle-bracket values stand for the full lowercase object
ID or 64-character content digest required by the schema; field names and nesting are normative.

```yaml
schema_version: 3
repository_id: example-repository
evidence_revision: <git-commit-id>
target_mode: mixed
inventory_revision: <git-commit-id>
policy_digest: <sha256>
candidate_set_digest: <sha256>
target_config_precondition_sha256: <sha256>
branch_evidence:
  algorithm: git-candidate-proof-v1
  git_version: 2.47.1
  evidence_scope: local
  remote_overlay_status: not-requested
  digest: <sha256>
records:
  - current_path: docs/repo/plans/example/example_implementation_doc.md
    target_path: docs/repo/plans/example/example_implementation_doc.md
    source_revision: <git-commit-id>
    source_sha256: <sha256>
    target_precondition: absent
    ownership_disposition: owned
    deferred: false
    branch_touch_reviews:
      - identity:
          evidence_scope: local
          repository_id: example-repository
          logical_ref: local:refs/heads/example-history
          candidate_path: docs/repo/plans/example/example_implementation_doc.md
        disposition: incorporated-history-non-owner
        ref_tip: <git-commit-id>
        merge_base: <git-commit-id>
        path_transitions:
          - path: docs/repo/plans/example/example_implementation_doc.md
            before: {state: present, mode: "100644", object_id: <git-object-id>, sha256: <sha256>}
            after: {state: present, mode: "100644", object_id: <git-object-id>, sha256: <sha256>}
        proof:
          kind: incorporated-history-v1
          transition_digest: <sha256>
          stable_patch_id: <git-patch-id>
          incorporated_parent: <git-commit-id>
          incorporated_commit: <git-commit-id>
          incorporated_transition_digest: <sha256>
          incorporated_stable_patch_id: <git-patch-id>
```

Schema conditionals require `proof.kind: same-content-v1` to carry target/ref path states and no
incorporated-commit fields; `incorporated-history-v1` requires every field shown above;
`deferred-active-owner` requires an owner-tip candidate snapshot and
`incremental_follow_up`; `active-owner` and `blocking` forbid clearance proof. Array order is
canonical but semantic identity is the four-tuple, so reordering cannot hide a duplicate.

### D-03 - Bind exact ref and path-state evidence

1. Problem: ref-name-only evidence does not invalidate when a ref moves but still touches the same
   path.
2. Simplest working solution: bind a normalized evidence record containing scope, full ref name,
   exact tip, unique merge base, evidence revision, ordered candidate transition entries, Git modes,
   blob object IDs, SHA-256 values, and an algorithm/version string. Hash the canonical JSON form
   into `branch_evidence_digest`.
3. Six-to-twelve-month pressure: SHA-256 object-format repositories, renamed plans, remote mirrors,
   and future proof kinds must remain reproducible.
4. Rationale: a moved/missing ref, changed merge base, path, mode, blob, hash, source revision,
   candidate policy, or proof implementation changes the digest and blocks.
5. Alternatives rejected: ref age; reflog timestamps; branch names; PR state; path names without
   blobs; worktree diffs.

Each path-state entry uses this semantic shape:

```yaml
path: docs/example/example_implementation_doc.md
before:
  state: present
  mode: "100644"
  object_id: <git-object-id>
  sha256: <content-sha256>
after:
  state: present
  mode: "100644"
  object_id: <git-object-id>
  sha256: <content-sha256>
```

`state` is `present` or `absent`. Absent states forbid mode/object/hash fields. Paths are
slash-normalized, unique, sorted by raw UTF-8 bytes, and interpreted as literal Git paths. Rename
and deletion evidence retains both old and new paths. Multiple merge bases, missing objects,
unsafe Git modes, malformed refs, and path collisions are `blocking`.

### D-04 - Use two-factor incorporated-history proof

1. Problem: stable patch IDs normalize content and can collide or match only part of a larger
   historical change.
2. Simplest working solution: compare the complete candidate-scoped transition from the ref's
   unique merge base to its tip against one incorporated commit's parent-to-commit transition.
   Require identical ordered path/state/mode/blob/hash transition digests and equal aggregate
   stable patch IDs. Require the incorporated commit to be an ancestor of the target revision.
3. Six-to-twelve-month pressure: a future reviewed contract may support a sequence of incorporated
   commits.
4. Rationale: exact state equality prevents partial-patch and patch-ID collision clearance; stable
   patch equality proves the aggregate edit, including textual change semantics.
5. Alternatives rejected: patch ID alone; final blob alone; PR merged state; commit-message
   matching; `git cherry` without path-state bindings.

`git-candidate-proof-v1` fixes the complete object-only patch pipeline. For each side of a rename,
delete, or add it sorts the literal union of candidate paths, disables rename detection so a
rename is represented as delete-plus-add, and invokes the equivalent of:

```text
git --literal-pathspecs -c color.ui=false -c core.quotePath=true
  -c diff.algorithm=myers -c diff.renames=false
  diff-tree --no-commit-id -r -p --full-index --binary --no-renames
  --no-ext-diff --no-textconv --no-color --no-indent-heuristic
  --diff-algorithm=myers --unified=3 <before-tree> <after-tree> -- <ordered-paths...>
git -c color.ui=false patch-id --stable
```

These are separate child processes; the first process's exact stdout is bounded and passed to the
second process through stdin without a shell. The environment fixes `LC_ALL=C`,
`GIT_CONFIG_NOSYSTEM=1`, `GIT_CONFIG_GLOBAL` to an empty reviewed file,
`GIT_TERMINAL_PROMPT=0`, `GIT_NO_LAZY_FETCH=1`, `GIT_NO_REPLACE_OBJECTS=1`, and
`GIT_LITERAL_PATHSPECS=1`. The implementation rejects alternates or missing/promisor objects that
would require network access. It never invokes hooks, external diff commands, textconv, filters,
credential helpers, or checked-out branch code.

Before diffing, the proof engine binds the `.gitattributes` blob stack and effective attributes at
the base, ref tip, incorporated parent, and incorporated commit. A custom diff driver,
`working-tree-encoding`, conflicting effective text/binary state, invalid UTF-8, NUL byte,
`GIT binary patch`, or `Binary files differ` result is unsupported and blocks incorporated-history
clearance; same-content proof may still succeed from exact object/content identity. EOL attributes
never normalize source hashes or transition state. Rename thresholds cannot affect the result
because rename detection is disabled.

The first release supports Git `>=2.43.0` only after the identical fixtures pass on every declared
runner. Runtime records `git_version`, rejects an older or unparsable version with
`branch_evidence_git_version_unsupported`, and reinventory is required when the proof algorithm
version changes. Tests must cover the oldest admitted Git, each release-runner Git, global and
repository config hostility, attribute changes, CRLF/LF blobs, renames, and binary candidates.

### D-05 - Permit only exact mixed-mode active-owner deferral

1. Problem: current ordering blocks a genuine owner before the strong existing grandfathering
   proof can be credited.
2. Simplest working solution: evaluate every touch review before the ownership stop. Permit exactly
   one `deferred-active-owner` when the target catalog is already valid mixed mode, current path and
   target path remain the cutover path, current target bytes equal the register-proven cutover
   blob, the ref tip contains a valid active plan candidate, all other refs are proven non-owners,
   and `incremental_follow_up` is complete.
3. Six-to-twelve-month pressure: multiple active owners represent real branch coordination, not a
   schema feature.
4. Rationale: default-branch authority remains frozen and visible while active work stays owned by
   its branch. The later merge cannot silently change grandfathered bytes.
5. Alternatives rejected: a global ignore list; canonical gaps; multiple owner deferrals for one
   plan; deletion-only owner deferral; dirty worktree deferral; config mode weakening.

Required follow-up fields bind:

- action `reinventory-and-incremental-migrate`;
- trigger `owner-ref-integrated-or-target-path-changed`;
- cutover revision, plan path, cutover source SHA-256, owner ref, owner tip, and owner path/blob;
- state `pending` at initial deferral; and
- success conditions: owner-supplied canonical metadata at merge, or a new reviewed ledger applied
  to the merged target bytes.

Direct canonical mode rejects `deferred-active-owner` with a stable blocker and a nonzero gap.

### D-06 - Guard activation and persist the exact deferral authority

1. Problem: metadata migration and discovery-v2 activation are separate. A ledger cannot contain
   the commit ID of the same commit that first tracks that ledger, and a ref, overlay, candidate,
   ledger, or target can move between review, metadata apply, activation, and commit.
2. Simplest working solution: use a two-checkpoint chronology with distinct revisions. Ledger v3
   binds `evidence_revision` `E`, the clean commit inventoried and reviewed. The ledger alone is
   committed as checkpoint `L`, whose first parent must be `E` and whose only tree delta is the
   exact ledger path/blob. Metadata migration runs from clean `HEAD=L` while recomputing branch
   proof against `E`. Guarded activation then accepts only the exact projected migration outputs
   as dirty paths and writes config. The next commit `A` must have first parent `L` and a tree
   delta containing exactly the projected metadata actions plus config. This removes commit
   self-reference while keeping every intermediate state provable.
3. Six-to-twelve-month pressure: automated follow-up tooling and portfolio views can depend on one
   explicit authority rather than infer approval from historical files.
4. Rationale: config schema v2's strict `plan_catalog_migration_evidence` object contains exactly
   `ledger_path`, `ledger_sha256`, `branch_evidence_digest`, `evidence_revision`,
   `activation_base_revision`, `migration_action_digest`, and `evidence_scope`. Here
   `activation_base_revision` is `L`, which is known when config is written. Validation discovers
   `A` as the first first-parent child of `L`, requires `A^1=L`, verifies the exact `L..A` delta,
   then allows later descendants while continuously rechecking the ledger, refs, candidates,
   frozen cutover, and pending follow-ups.
5. Alternatives rejected: a self-referential ledger/activation commit; a permissive config-v1
   extension; an unbound retained ledger; manual config edits; modifying the frozen register;
   hidden local-only approval; letting `plans migrate` commit or activate config.

The exact chronology is:

1. Require clean index/worktree at evidence revision `E`; inventory and author the ledger against
   `E` without applying metadata.
2. Review the ledger, then create a separately authorized ledger-only commit `L`. Require
   `L^1=E`; `git diff-tree E L` contains exactly one normalized tracked regular ledger path, and
   its blob hash equals the reviewed ledger hash. Any amended/squashed/extra-path checkpoint
   requires new inventory and review.
3. From clean `HEAD=L`, run migration dry-run and apply. V3 planning permits the ledger-only delta
   but requires every candidate/config/register blob at `L` to equal its evidence-revision state.
   Apply may change the worktree only at the exact projected metadata action paths; the index must
   remain identical to `L`.
4. Run `plans catalog-activate --ledger <path> --dry-run`, then `--yes`. It rebuilds evidence
   against `E`, requires current `HEAD=L` and index tree `L`, verifies the complete worktree delta
   (including projected new/untracked and removed paths) byte-for-byte against the migration action
   digest, requires config still equal its precondition, and writes only config v2 plus the
   binding. Any staged or intent-to-add path, and any extra/mismatched unstaged, untracked,
   renamed, or deleted path, blocks.
5. Create separately authorized activation checkpoint `A` containing exactly the action paths and
   config. Require `A^1=L`; no ledger edit and no unrelated path are allowed.
6. Run `plans validate` immediately at clean `HEAD=A`. It locates `A` as the first first-parent
   child of `L`, verifies the exact action/config tree delta and binding, proves `A` is an ancestor
   of the current revision, and then performs normal current-ref/candidate/cutover validation.

The ledger path must be a normalized repository-relative path to a tracked regular file outside
all migration targets and protected configuration/register paths, normally the owning plan's
`attachments/` directory. Checkpoint commits require explicit repository-local commit authority;
neither migration nor catalog activation commits, stages, pushes, or broadens Git authority.
Squash/rebase of `L` or `A`, non-first-parent insertion before `A`, or an extra checkpoint path
invalidates the proof and requires reinventory rather than hash repair.

The config-v2 binding shape is normative:

```yaml
schema_version: 2
component_settings:
  planning-workflows:
    plan_catalog_mode: mixed
    plan_catalog_discovery_version: 2
    plan_catalog_cutover_revision: <git-commit-id>
    plan_catalog_migration_evidence:
      ledger_path: docs/repo/plans/example/attachments/plan-migration-ledger-v3.yaml
      ledger_sha256: <sha256>
      branch_evidence_digest: <sha256>
      evidence_revision: <git-commit-id>
      activation_base_revision: <git-commit-id>
      migration_action_digest: <sha256>
      evidence_scope: local
```

The omitted required config-v2 fields retain their released names and meanings. The schema closes
the evidence object, and runtime requires its values to equal the ledger and the proved `E/L`
history rather than merely accepting well-formed hashes.

A movement after metadata apply but before activation, including a ref that retains the same
touched path set, yields `catalog_activation_authority_drift` and zero config writes. Movement
after activation but before `A` makes the mandatory first post-commit validation fail closed.
After `A`, `plans validate`, `plans list`, migration, and member portfolio collection load the
binding before crediting a compatibility row; missing, changed, or stale evidence makes the
catalog invalid rather than silently dropping the deferred plan.

Config v1 remains accepted and byte-unchanged for legacy, mixed-without-v3-deferral, and canonical
repositories. New installs and ordinary upgrades are not automatically rewritten to v2. Only the
guarded activation command may introduce the v2 binding, and the release manifest advertises
config schema versions 1 and 2 only after upgrade/install parity tests pass.

Config-v2 conditionals require `plan_catalog_mode: mixed`,
`plan_catalog_discovery_version: 2`, a valid cutover revision, and the evidence object; canonical
or legacy config v2 is invalid. A zero-gap direct-canonical activation is allowed only from a new
reviewed projection with no deferrals and writes config schema v1 without the evidence object.
That guarded transition is the only normal path that removes the binding; rollback restores the
exact prior config bytes rather than synthesizing an unbound config.

### D-07 - Separate scan completeness, mixed coverage, and canonical readiness

1. Problem: remote portfolio observations currently omit metadata-less plans, and one `complete`
   boolean cannot express a complete scan with reviewed mixed compatibility gaps.
2. Simplest working solution: version local structured views and portfolio catalog/cache outputs to
   3, add explicit compatibility observations, and report `complete`,
   `mixed_coverage_complete`, and `canonical_ready` separately.
3. Six-to-twelve-month pressure: coordination must distinguish current facts from readiness for a
   mode transition.
4. Rationale: a deferred plan stays visible and a complete scan remains truthful without claiming
   canonical completeness.
5. Alternatives rejected: omitting the plan; treating low-level candidates as sufficient current
   plan rows; marking the entire scan unavailable; calling mixed state canonical.

Preserve catalog v2 as historical evidence. New cache loading validates raw JSON against the
versioned schema before struct decoding so unknown fields cannot be silently discarded.

### D-08 - Make offline and remote-overlay lanes explicit

1. Problem: local migrations need deterministic offline evidence, while requested remote overlays
   must not become stale or silently unbound.
2. Simplest working solution: inventory records its evidence scope. A local ledger binds local and
   tracking refs and states `remote_overlay_status: not-requested`. A remote-aware ledger binds the
   complete current-repository overlay digest and exact remote ref snapshots; both dry-run and
   apply require `--remote-overlays` and a matching complete scan.
3. Six-to-twelve-month pressure: additional provider adapters may supply corroboration but not
   ownership authority.
4. Rationale: offline application is honest about its limit, while remote-aware application cannot
   reuse a stale overlay.
5. Alternatives rejected: always-online migration; PR-only evidence; silently using the previous
   complete cache; treating a stale local tracking ref as cleared due to age.

Remote-aware inventory is prospective and member-scoped. Before config activation, the scanner
reads the remote default-branch config blob but substitutes only the ledger's target discovery
version and candidate policy in memory; it does not require or write discovery v2 remotely. The
remote default commit must equal ledger `evidence_revision` `E` when remote evidence is first
reviewed; at activation it may equal `L` only when `L^1=E` and the exact ledger-only delta is
present remotely. After activation checkpoint `A` is pushed, remote default may equal `A` or any
later commit `D` only when the mirror proves `E` is the first parent of `L`, `L` is the first
parent of `A`, both checkpoint tree deltas are exact, `A` is an ancestor of `D`, and current bound
config/candidate/register/cutover states remain valid. Branch proof stays anchored to `E`; remote
advancement never silently rebases the ledger. A squashed/rebased `L` or `A`, missing checkpoint
object, default history divergent from `A`, or changed bound state blocks and requires new
evidence. The scanner then scans every fresh
`refs/heads/*` ref from that source mirror. Completeness for migration means this repository
member's default ref, config, candidate discovery, required objects, and full ref namespace all
succeeded. An unrelated portfolio member may keep aggregate portfolio `complete: false`, but does
not block this member-scoped migration; any incomplete fact for the target member blocks.

Logical ref normalization is exact:

- a local branch is `local:refs/heads/<name>`;
- a tracking ref whose configured remote URL uniquely matches the selected portfolio source is
  `remote:<source-id>:refs/heads/<name>`;
- a tracking ref with no unique configured-source match remains
  `tracking:<remote-name>:refs/heads/<name>`; and
- a mirror ref is `remote:<source-id>:refs/heads/<name>`.

Local tracking and mirror evidence coalesce only when logical key, tip, merge base, and complete
candidate transition are byte-identical. A mismatch is two evidence records and normally a
blocking stale-tracking conflict; it is never resolved by age. Current authority always comes
from the fresh invocation, not the prior cache. The cache remains valuable historical evidence
and is preserved on failure, but cannot satisfy activation or apply.

### D-09 - Preserve atomicity and released behavior

1. Problem: new clearances could bypass existing transaction checks, and the no-action fast path
   needs explicit proof revalidation.
2. Simplest working solution: normalize v3 evidence into the migration authority digest, run the
   authority check before every preview/apply outcome, retain the in-transaction prewrite check,
   then recheck after the complete action batch and once more immediately before transaction
   cleanup. Any post-write drift invokes the existing rollback/recovery path. The guarded
   config-only activation transaction uses the same chronology.
3. Six-to-twelve-month pressure: further evidence sources must plug into one authority snapshot.
4. Rationale: a clearance changes only whether a touch blocks; it does not weaken file, config,
   candidate, or transaction authority.
5. Alternatives rejected: one-time ledger validation; skipping ref recheck on idempotent apply;
   accepting a prewrite-only check; writing config in the metadata transaction; partial success.

The transaction API must expose separate `PreWriteAuthorityCheck` and
`PostWriteAuthorityCheck` callbacks. For an action-bearing plan the post-write callback runs after
all actions and post-checks but before marker/stage cleanup; a mismatch rolls back every action.
For a zero-action/idempotent plan the authority callback still runs before success. Tests inject a
ref move after prewrite, between two file actions, after the action batch, and immediately before
cleanup. A ref move after successful cleanup becomes new repository state and is caught by the
persisted binding on the next view or activation operation; no finite file transaction can lock an
external ref indefinitely.

Authority checking is phase-aware. The immutable `external_authority_digest` covers `E`, `L`,
ledger/policy/ref/overlay/proof inputs and must be identical before and after writes. The
`pre_state_digest` covers exact source/config preconditions; the `post_state_digest` covers exact
planned target blobs, removals, and expected config output. Post-write validation compares against
the post state instead of incorrectly expecting pre-migration candidate bytes, while any external
movement still rolls back.

### CLI and structured output contract

| Command | Required UX after implementation |
| --- | --- |
| `plans inventory` | Emit schema v3 branch-touch records, evidence revision `E`, deterministic proposed classifications, evidence scope/digest, exact proof inputs, blockers, and remote-overlay status. It never authors or approves a ledger. |
| `plans migrate --ledger ... --dry-run` | Verify ledger-only checkpoint `L` against `E`, require v3 proof recomputation, show every ref disposition/blocker, projected actions/coverage/follow-ups, remote status, `activation_performed: false`, and zero writes. |
| `plans migrate --ledger ... --yes` | Apply only the exact reviewed non-deferred plan actions atomically; retain deferred rows as non-material skips; reject changed evidence before commit. |
| `plans catalog-activate --ledger ... --dry-run` | Require `HEAD=L`, index tree `L`, and exact applied metadata outputs as the only worktree delta; recompute authority against `E`, and show config-v2 binding, action digest, exact `A` path set, projection, and zero writes. |
| `plans catalog-activate --ledger ... --yes` | Atomically change only reviewed catalog config fields, persist `E`/`L`/ledger/evidence/action bindings, post-write recheck authority, roll back config on drift, and report that checkpoint `A` plus clean validation remain required. |
| `plans validate` | Verify/disclose evidence revision `E`, ledger checkpoint `L`, activation checkpoint `A`, scan/mixed/canonical readiness, deferred compatibility rows, and incremental-migration blockers. |
| `plans list` | Show each deferred plan once with its legacy ID/path/title/lifecycle and explicit `mixed-grandfathered` coverage. |
| `portfolio scan` | Keep default and branch facts separate, add compatibility observations, retain optional PR corroboration, preserve the last complete compatible cache, and never call a mixed member canonical-ready. |

Stable blocker codes must cover at least: unknown disposition, missing/unreviewed touch, ref missing,
ref moved, merge base missing/ambiguous/moved, path state moved, blob/hash mismatch, transition
digest mismatch, patch mismatch, incorporated commit unreachable, partial incorporation, remote
overlay missing/stale/incomplete, unsupported Git version/attributes/binary proof, multiple active
owners, invalid active lifecycle, dirty owner, mixed proof missing, follow-up missing/stale,
activation authority drift, evidence binding missing/mismatched, canonical deferral forbidden,
ledger checkpoint parent/delta mismatch, activation checkpoint missing/parent/delta mismatch,
post-write authority drift, and rollback failure.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Versioned v3 schemas and explicit v1/v2/v3 compatibility dispatch exist. | L | None |
| EP-02 | Deterministic local and remote branch-evidence engine classifies every touch without granting clearance. | XL | EP-01 |
| EP-03 | Migration runtime applies reviewed non-owner clearances and exact mixed deferrals atomically. | XL | EP-01, EP-02 |
| EP-04 | Local views, remote overlays, and portfolio cache preserve deferred visibility and readiness semantics. | L | EP-01, EP-02, EP-03 |
| EP-05 | Managed doctrine, routing, packaged resources, and consumer-impact evidence match the implementation. | L | EP-03, EP-04 |
| EP-06 | Exhaustive adversarial, cross-platform, acceptance-scenario, and rollback proof passes. | XL | EP-01 through EP-05 |
| EP-07 | A separately authorized patch release is reproducible, published through the release gate, and handed off for repository-owned rollout. | L | EP-06 |

## EP-01 - Versioned Evidence And Compatibility Contracts

### A) Epic ID, Title, And Outcome

`EP-01 - Versioned Evidence And Compatibility Contracts`

Outcome: strict ledger, inventory, catalog, and config-binding contracts describe the complete
evidence model, while released ledger/catalog v1/v2 and config v1 artifacts retain exact behavior
and old CLIs reject new versions before mutation.

### B) Scope

- Freeze ledger v2 and catalog v2 under explicit historical filenames.
- Define ledger v3 conditionals for dispositions, path states, proof types, overlay bindings, and
  incremental follow-up.
- Define inventory v3 and catalog/cache v3 structured output.
- Add explicit schema version dispatch and raw-before-struct cache validation.
- Preserve top-level config schema v1 under a historical filename and define strict config v2
  binding conditionals for deferred mixed activation.

### C) Files Touched

```text
schemas/
├── kit-config-v1.schema.json               # create
├── kit-config.schema.json                  # modify: latest strict v2
├── plan-migration-ledger-v2.schema.json    # create
├── plan-migration-ledger.schema.json       # modify
├── plan-inventory.schema.json              # create
├── plan-catalog-v2.schema.json             # create
└── plan-catalog.schema.json                # modify
internal/state/schema.go                    # modify
internal/state/observed.go                  # modify
internal/lockfile/lockfile.go               # modify
internal/plancatalog/model.go               # modify
internal/plancatalog/migrate.go             # modify: versioned decode/dispatch types only
internal/plancatalog/view.go                # modify: versioned config decode types only
internal/commands/portfolio.go              # modify
internal/portfolio/config.go                # modify
internal/portfolio/membership.go            # modify
internal/portfolio/source.go                # modify: v3 models
internal/portfolio/scanner.go               # modify: raw version validation
tests/test_json_schemas.py                  # modify
internal/plancatalog/migration_v2_test.go    # modify
internal/plancatalog/plancatalog_test.go     # modify
```

### D) Acceptance Criteria And Size

- Size: `L`.
- Current v2 schema bytes are preserved exactly in the new historical file.
- Version dispatch accepts only 1, 2, and 3 and rejects every other value.
- V3 schema conditionals require proof-specific fields and forbid irrelevant fields.
- Every touch disposition and follow-up field is closed to unknown values/properties.
- Old v0.1.25/v0.1.26 binaries reject v3 ledger fixtures without modifying targets.
- Old v0.1.25/v0.1.26 binaries reject config-v2 fixtures before plan settings are used.
- New binary reproduces current v2 active-branch blocking and mixed behavior.
- Config-v1 schema bytes and existing consumer config bytes remain unchanged; config v2 is valid
  only for discovery-v2 mixed mode with the strict binding, while runtime requires a pending
  reviewed follow-up.

### E) Dependencies And Critical-Path Notes

No dependency. EP-02 and EP-03 must consume these exact types rather than introducing parallel
unversioned evidence structs.

### F) Tasks Checklist

- [ ] Copy the released `schemas/plan-migration-ledger.schema.json` bytes into `schemas/plan-migration-ledger-v2.schema.json` and assert byte identity.
- [ ] Copy released `schemas/kit-config.schema.json` bytes into `schemas/kit-config-v1.schema.json`, promote the unsuffixed schema to version 2, and assert historical byte identity.
- [ ] Define the strict config-v2 `plan_catalog_migration_evidence` object and mixed/discovery conditionals without adding the field to config v1; enforce the cross-ledger pending-deferral invariant at runtime.
- [ ] Replace `schemas/plan-migration-ledger.schema.json` with the strict ledger v3 envelope, per-ref reviews, proof conditionals, overlay binding, and incremental follow-up contract.
- [ ] Add `schemas/plan-inventory.schema.json` for deterministic inventory v3 evidence and public-safe structured blockers.
- [ ] Copy the released `schemas/plan-catalog.schema.json` bytes into `schemas/plan-catalog-v2.schema.json` and assert byte identity.
- [ ] Upgrade `schemas/plan-catalog.schema.json` to v3 compatibility observations and split readiness fields.
- [ ] Add one raw-map-first config v1/v2 dispatcher and replace every hardcoded `ConfigV1Schema` consumer in state observation, lockfile, reconcile, portfolio commands/config/membership, and plan-catalog views; reject missing/unsupported versions before typed decode.
- [ ] Add explicit ledger/catalog v1/v2/v3 dispatch with default rejection in `internal/state/schema.go` and `internal/plancatalog/migrate.go`.
- [ ] Add version-specific ledger decode types and a normalized internal migration model without interpreting v2 free text as v3 clearance.
- [ ] Validate raw cached catalog JSON before version-specific struct decoding in `internal/portfolio/scanner.go`.
- [ ] Extend `tests/test_json_schemas.py` with positive, missing-field, unknown-field, cross-proof, and forbidden-field fixtures.
- [ ] Extend Go compatibility tests with config v1/v2 acceptance, ledger/catalog v1/v2/v3 acceptance, unsupported-version rejection, and old-binary ledger-v3/config-v2 rejection.
- [ ] Run `go test ./internal/state ./internal/plancatalog ./internal/portfolio` and `python -m pytest tests/test_json_schemas.py -q`.

### G) Implementation Notes

Do not mutate v2 semantics. `BuildMigrationPlan` must use an exhaustive version switch; the current
`if version == 2 else v1` shape is unsafe once v3 exists. The v3 evidence arrays must be sorted and
unique at schema and runtime layers. Use full 40-64 lowercase hexadecimal commit/object patterns
where repository object format can vary; retain SHA-256 for content portability.

### H) Open Questions

None. OQ-01 selects the strict config-v2 evidence binding and forbids a config-v1 extension.

## EP-02 - Deterministic Branch Evidence Engine And Inventory UX

### A) Epic ID, Title, And Outcome

`EP-02 - Deterministic Branch Evidence Engine And Inventory UX`

Outcome: local and requested remote inventory produce the same normalized ref/path proof model and
make same-content, incorporated-history, active-owner, and ambiguous candidates reviewable without
granting migration authority.

### B) Scope

- Replace ref-name-only touch facts with immutable evidence snapshots.
- Implement unique merge-base, exact path-state transition, same-content, reachable incorporated
  commit, and stable patch calculations.
- Keep age and PR data non-authoritative.
- Add public-safe text/JSON inventory UX and evidence digests.
- Retain current ref namespace, hard-unowned, dirty, and candidate policy boundaries.

### C) Files Touched

```text
internal/plancatalog/
├── branch_reconciliation.go              # create
├── branch_reconciliation_test.go         # create
├── model.go                              # modify
├── git_index.go                          # modify
├── inventory.go                          # modify
└── classifier_test.go                    # modify
internal/commands/
├── plans.go                              # modify
└── commands_test.go                      # modify
internal/cli/
├── cli.go                                # modify
└── cli_test.go                           # modify
internal/portfolio/
├── source.go                             # modify
├── scanner.go                            # modify
├── mirror.go                             # modify
└── portfolio_test.go                     # modify
tests/fixtures/plans/                     # modify
tests/fixtures/portfolio/                 # modify
```

### D) Acceptance Criteria And Size

- Size: `XL`.
- Every relevant touch has scope/ref/tip/base/target/path-state evidence and one canonical digest.
- Same-content proposals require exact same-path/mode/object/hash identity.
- Incorporated-history proposals require exact transition equality plus stable patch equality to a
  specific reachable commit transition.
- Partial, subset, whitespace-normalized collision, rename ambiguity, deletion ambiguity, multiple
  merge base, missing object, and moved ref cases propose `blocking`.
- Local and remote evidence serialize identically for equivalent Git object graphs.
- Inventory output remains deterministic across Linux, macOS, and Windows.
- No repository content, hook, diff driver, textconv, helper, lazy fetch, or replace object runs.
- Git below 2.43.0, unsupported effective attributes, binary/NUL candidate content, and any
  unproven release-runner Git version block incorporated-history proof deterministically.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 schema/model names. EP-03 cannot start its disposition logic until evidence
digests and recomputation APIs are stable.

### F) Tasks Checklist

- [ ] Implement normalized ref snapshots and ordered path-state transitions in `internal/plancatalog/branch_reconciliation.go`.
- [ ] Replace `ActiveBranchTouch []string` authority with structured touch evidence while retaining a compatibility projection for v1/v2 output.
- [ ] Require one unique merge base and emit a stable blocker for missing and multiple-base histories.
- [ ] Implement exact same-content proposal logic from target and ref-tip path states.
- [ ] Implement incorporated-commit proposal logic with target ancestry, exact transition digest equality, and aggregate stable patch equality.
- [ ] Implement the exact D-04 diff-tree-to-`patch-id --stable` pipeline with disabled renames, canonical path ordering, inert config/environment, bounded stdin, and no shell.
- [ ] Bind effective attribute state at all four proof endpoints; reject custom diff drivers, working-tree encoding, inconsistent text/binary state, invalid UTF-8, NUL, and binary patch output.
- [ ] Enforce the Git 2.43.0 floor, record the Git version, and prove the algorithm on the oldest admitted and every release-runner version.
- [ ] Add bounded size and count limits for refs, paths, diff bytes, and incorporated proof candidates with fail-closed blockers.
- [ ] Compute `branch_evidence_digest` from canonical JSON that binds evidence algorithm, evidence revision, policy, candidates, refs, tips, bases, transitions, and proof results.
- [ ] Extend `plans inventory` text and JSON output with proposed dispositions, proof summaries, blockers, evidence scope, and remote-overlay status.
- [ ] Extend remote mirror/scanner adapters to emit the same ref snapshot and transition model from a prospective discovery-v2 member scan without mutating remote config.
- [ ] Implement D-08 logical-ref normalization and coalesce local tracking/mirror evidence only on exact source, tip, base, and transition equality.
- [ ] Keep pull-request facts under a separate `corroboration` object excluded from proof sufficiency.
- [ ] Add local, tracking-ref, remote-only, rename, deletion, stale, moved, ambiguous-base, partial-patch, and hostile-Git fixtures.
- [ ] Run `go test ./internal/plancatalog ./internal/commands ./internal/cli ./internal/portfolio`.

### G) Implementation Notes

Inventory suggestions are not reviewed dispositions. Do not auto-create a ledger from suggestions.
The incorporated proof scope must include the full ordered transition for every candidate path
affected by that touch, including both sides of rename/delete entries. A target commit with one
matching hunk cannot satisfy a two-path ref transition.

Remote evidence must be filtered to the current repository ID before it enters a migration
artifact. Public output uses sanitized source locators and never exposes credentials, local mirror
paths, private provider URLs, or raw provider errors.

Evidence digesting is candidate-scoped: one ref touching two plans yields two independent review
identities. The branch summary may aggregate those reviews for display but cannot become
clearance authority. The remote result carries both `member_evidence_complete` for the selected
repository and portfolio-wide `complete`; only the former gates that repository's ledger.

### H) Open Questions

- OQ-02 fixes the first implementation at one incorporated target commit.
- OQ-03 defines the separate local and requested-overlay lanes.

## EP-03 - Reviewed Clearance, Mixed Deferral, And Atomic Migration

### A) Epic ID, Title, And Outcome

`EP-03 - Reviewed Clearance, Mixed Deferral, And Atomic Migration`

Outcome: ledger v3 can clear only proof-complete non-owners, can preserve one exact active owner in
mixed mode with mandatory follow-up, and blocks every stale or ambiguous case before metadata or
activation writes, with rollback on authority movement during either transaction.

### B) Scope

- Reconcile every inventory touch to exactly one v3 ledger review.
- Reorder branch enforcement around proof validation.
- Preserve v1/v2 behavior.
- Bind follow-up and mixed cutover evidence.
- Recheck authority on dry-run, apply, idempotent apply, before writes, after writes, and before
  transaction cleanup.
- Keep config/register protected during metadata migration and add a separate guarded config-only
  catalog activation.

### C) Files Touched

```text
internal/plancatalog/
├── migrate.go                         # modify
├── activation.go                      # create
├── migration_v3_test.go               # create
├── migration_v2_test.go               # modify
├── ep02_test.go                        # modify
└── view.go                             # modify: shared mixed proof helper
internal/commands/
├── plans.go                            # modify
└── commands_test.go                    # modify
internal/reconcile/plan.go              # inspect; modify only for required authority interface
internal/reconcile/transaction.go       # modify: zero-action plus pre/post-write authority checks
schemas/kit-config.schema.json          # refine within EP-01 contract
tests/fixtures/plans/                   # modify
```

### D) Acceptance Criteria And Size

- Size: `XL`.
- Every current touch is reviewed once; missing and extra reviews block.
- `same-content-non-owner` and `incorporated-history-non-owner` clear only after runtime proof
  recomputation.
- At most one `deferred-active-owner` exists per candidate and every other touch is cleared.
- Deferred ref-tip plan evidence has lifecycle `active`, valid candidate semantics, exact owner
  path/blob/hash, and complete follow-up.
- Mixed deferral requires unchanged register/cutover/current bytes and no metadata at the default
  path.
- Canonical target rejects deferral and reports a gap.
- Dirty/uncommitted state cannot be cleared or deferred.
- Ref/source/policy/candidate/target/follow-up drift yields zero writes and a stable remediation.
- Dry-run and apply have the same projection; migration never activates config.
- Idempotent apply rechecks branch authority before success.
- V3 migration runs only from clean ledger-only checkpoint `L`; it verifies `L^1=E` and rejects
  any non-ledger `E..L` delta before using evidence from `E`.
- `plans catalog-activate` rechecks the exact ledger, candidate, local ref, requested overlay, and
  projection after migration; any movement since apply yields zero config writes.
- Successful mixed activation writes config v2 plus the evidence binding atomically, and later
  validation blocks if the bound ledger/ref/overlay/candidate authority moves.
- The first clean post-commit validation finds exact activation checkpoint `A`, verifies
  `A^1=L` and the complete `L..A` action/config delta, and permits only later descendants of a
  valid `A`.
- Ref/overlay movement after prewrite, between file actions, after the action batch, or before
  cleanup rolls back all writes.
- Transaction failure restores prior bytes and reports recovery state precisely.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 contracts and EP-02 proof recomputation. This epic is the safety-critical path.
No doctrine or release work can declare the capability usable before EP-03 review passes.

### F) Tasks Checklist

- [ ] Add explicit v3 migration planning that matches every inventory touch to one ledger review and rejects missing, duplicate, and extra reviews.
- [ ] Recompute same-content and incorporated-history proof fields before crediting non-owner dispositions.
- [ ] Move branch blocking after proof reconciliation while retaining all source, dirty, ownership, semantic, and target preconditions.
- [ ] Require exactly one valid `deferred-active-owner` and proof-complete non-owner dispositions for all remaining touches.
- [ ] Bind deferred active lifecycle, owner ref/tip/path/blob/hash, cutover revision/path/hash, and incremental follow-up.
- [ ] Reject deferred owner deletion, multiple active owners, non-active lifecycle, current-path rename, dirty overlap, and canonical target mode.
- [ ] Add branch evidence and follow-up state to the migration authority digest.
- [ ] Invoke the authority check before every dry-run/apply result, including zero-action apply, and implement prewrite, post-action, and pre-cleanup checks for action-bearing transactions.
- [ ] Preserve config and frozen register as protected migration targets and retain `activation_performed: false` in text/JSON output.
- [ ] Implement `plans catalog-activate --ledger ... --dry-run|--yes` as a separate config-only reconcile plan with an exact current-config SHA-256 precondition.
- [ ] Implement the D-06 `E -> L -> A` chronology: verify the ledger-only `E..L` delta, run migration from clean `HEAD=L`, permit only exact projected outputs before activation, and verify the first post-commit `L..A` delta.
- [ ] Make activation require a tracked regular ledger, clean `HEAD=L`, exact applied metadata bytes, complete mixed projection, and matching fresh overlay when bound.
- [ ] Persist config-v2 ledger path/hash, branch evidence digest, evidence revision `E`, activation base `L`, migration-action digest, and evidence scope; require views to revalidate the binding while a follow-up remains.
- [ ] Extend projection output with disposition counts, compatibility rows, pending follow-ups, mixed coverage, canonical readiness, and remote status.
- [ ] Preserve ledger v1/v2 code paths and assert their output and blockers remain unchanged.
- [ ] Add zero-write tests for every proof mismatch, stale source, moved ref, dirty owner, collision, and missing follow-up.
- [ ] Add activation zero-write tests for ledger/ref/overlay/candidate movement after metadata apply and before config write.
- [ ] Add end-to-end clean `E`, ledger-only `L`, migration, guarded activation, exact `A`, and first post-commit validation tests, including later unrelated descendants.
- [ ] Add staged-action rollback tests for ref movement after prewrite, between two actions, after the action batch, before cleanup, write failure, post-check failure, and recovery-marker failure.
- [ ] Run `go test ./internal/plancatalog ./internal/commands ./internal/reconcile`.

### G) Implementation Notes

Non-owner clearance changes only the branch-touch decision. It must not suppress validation
problems, change ownership to unowned, remove the candidate from projection, or relax exact source
and target file preconditions.

For mixed deferral, the migration has no action for the deferred file. It may still apply metadata
to unrelated reviewed candidates in one transaction because projection counts the exact
grandfathered row and no unreviewed gap remains. Any failed clearance makes the whole migration
blocked; there is no partial success.

Catalog activation is not plan lifecycle activation and does not mark this implementation plan
active. It changes only the repository's shared plan-catalog config under the migration runbook.
Manual config editing is unsupported for a deferral because it cannot establish the transaction
authority chronology or the persistent binding.

### H) Open Questions

None. D-05 and D-06 define the single implementation and activation path.

## EP-04 - Mixed Visibility, Remote Overlays, And Portfolio Readiness

### A) Epic ID, Title, And Outcome

`EP-04 - Mixed Visibility, Remote Overlays, And Portfolio Readiness`

Outcome: a deferred active plan remains one visible compatibility row locally and remotely,
complete scans remain distinguishable from canonical readiness, and requested overlay evidence is
freshly rebound at migration time.

### B) Scope

- Add local view coverage/follow-up fields.
- Add remote compatibility observations and catalog v3 readiness fields.
- Add `plans migrate --remote-overlays` composition at the command boundary.
- Feed the same prospective member-scoped overlay to guarded catalog activation.
- Preserve previous complete compatible cache on failure.
- Preserve offline local migration with an explicit visibility limitation.
- Prove later merge/incremental migration behavior.

### C) Files Touched

```text
internal/plancatalog/
├── model.go                            # modify
├── view.go                             # modify
└── migration_v3_test.go                # modify
internal/commands/
├── plans.go                            # modify
└── commands_test.go                    # modify
internal/portfolio/
├── source.go                           # modify
├── scanner.go                          # modify
├── mirror.go                           # modify
├── store.go                            # modify
└── portfolio_test.go                   # modify
schemas/plan-catalog.schema.json        # refine within EP-01 contract
tests/fixtures/portfolio/               # modify
```

### D) Acceptance Criteria And Size

- Size: `L`.
- Mixed `plans list` returns one deferred row with legacy evidence, lifecycle, canonical path,
  `coverage_disposition: mixed-grandfathered`, and `incremental_migration_required: true`.
- Remote catalog v3 emits a compatibility observation instead of omitting the plan.
- `complete` describes the scan, `mixed_coverage_complete` describes reviewed coverage, and
  `canonical_ready` remains false until no grandfathered observation exists.
- Remote-aware ledgers require a fresh complete matching overlay for dry-run and apply.
- Remote-aware activation requires a fresh complete matching target-member overlay; unrelated
  incomplete portfolio members are disclosed but do not block the selected member.
- Post-activation remote scans accept default at exact `A` or a later descendant only after
  verifying remote `E -> L -> A` checkpoints and current bound state; squash/rebase/divergence
  blocks without rebasing proof.
- Offline ledgers work from local objects and state that remote freshness/completeness is unproven.
- A failed/incomplete scan preserves the prior complete compatible cache byte-for-byte.
- A v1/v2 cache remains historical and cannot be presented as current v3 evidence.
- A later branch merge with changed filename-only bytes blocks; owner-supplied metadata and
  reviewed incremental migration both resolve the follow-up.

### E) Dependencies And Critical-Path Notes

Depends on EP-03 disposition and projection semantics. Remote overlay implementation must remain in
`internal/portfolio`; `internal/plancatalog` receives normalized evidence and must not import the
portfolio package.

### F) Tasks Checklist

- [ ] Add mixed compatibility and incremental-follow-up fields to local v3 view rows and text output.
- [ ] Add remote compatibility observations with default ref, commit, path, title, lifecycle, legacy evidence, and content hash.
- [ ] Separate scan completeness, mixed coverage completeness, and canonical readiness in catalog v3.
- [ ] Add `--remote-overlays` to `plans migrate` and compose filtered current-repository evidence at the command layer.
- [ ] Require matching complete remote evidence for ledgers whose evidence scope includes overlays.
- [ ] Add `--remote-overlays` to `plans catalog-activate` and recompute the same prospective member scan before config planning and immediately before write.
- [ ] Define `member_evidence_complete` separately from portfolio `complete`; block on any selected-member default/config/ref/object/proof failure but not an unrelated member failure.
- [ ] Test the exact local/tracking/mirror namespace mapping and reject ambiguous remote-URL source matches.
- [ ] Keep remote proof anchored to `E` after push; verify exact remote `E -> L -> A`, allow unchanged `A` descendants, and reject missing/squashed/rebased/divergent checkpoint histories.
- [ ] Emit explicit offline limitations without consuming stale cache as current authority.
- [ ] Validate raw cache JSON by version and preserve v1/v2 cache as labeled historical evidence.
- [ ] Preserve the last complete compatible v3 cache on auth, access, fetch, ref, proof, and member failures.
- [ ] Add end-to-end mixed activation, deferred visibility, branch merge, validation failure, and incremental migration tests.
- [ ] Add remote same-content, incorporated-history, stale ref, moved ref, missing ref, optional PR, and offline cache tests.
- [ ] Run `go test ./internal/plancatalog ./internal/commands ./internal/portfolio`.

### G) Implementation Notes

Remote pull-request facts may enrich a compatibility observation. The ref's Git proof remains the
only clearance authority. A provider-reported merged PR with missing or mismatched Git evidence is
blocking.

An offline local migration may complete metadata writes and mixed projection. A local-scope ledger
may then activate with `remote_overlay_status: not-requested`, but the runbook and outputs must say
that remote freshness and aggregate portfolio completeness remain unclaimed. A remote-aware
ledger cannot activate from cache. Portfolio analysis cannot call remote evidence current until a
complete target-member refresh succeeds.

### H) Open Questions

OQ-03 supplies the default offline/overlay behavior.

## EP-05 - Managed Doctrine, Routing, And Packaged Resource Parity

### A) Epic ID, Title, And Outcome

`EP-05 - Managed Doctrine, Routing, And Packaged Resource Parity`

Outcome: agent-facing doctrine makes the reviewed evidence, clearance, deferral, activation,
incremental follow-up, compatibility, recovery, and approval paths executable by a fresh agent,
with byte-identical packaged resources and routing evidence.

### B) Scope

- Update catalog, lifecycle, portfolio, migration, register, and routing doctrine.
- Keep the migration runbook agent-facing and L1.
- Keep plan commands L3 and record no L2 promotion.
- Add consumer-impact and migration/release-note inputs.
- Update managed resource identities and mirrors.
- Run a fresh low-context routing probe.

### C) Files Touched

```text
components/planning-workflows/
├── component.yaml
└── managed/
    ├── README.md
    ├── reference/plan-catalog-format.md
    ├── reference/planning-document-lifecycle.md
    ├── reference/portfolio-coordination-format.md
    └── runbooks/
        ├── migrate-plan-catalog.md
        └── maintain-plan-register.md
components/agent-interface/managed/reference/operation-routing-and-dispatch.md
src/codeheart_operating_kit/resources/components/planning-workflows/       # matching mirrors
src/codeheart_operating_kit/resources/components/agent-interface/         # matching mirror
docs/repo/plans/branch-ownership-reconciliation/attachments/
├── consumer-impact-record.md                 # create during execution
└── release-note-input.md                     # create during execution
tests/test_packaging_resources.py              # modify
tests/test_routing.py                          # modify
docs/README.md                                 # inspect/update when command/schema discoverability changes
docs/repo/README.md                            # inspect/update when command/schema discoverability changes
```

### D) Acceptance Criteria And Size

- Size: `L`.
- The migration runbook names source, inputs, preconditions, lane, approval, exact steps, stop
  conditions, evidence, validation, activation separation, incremental follow-up, and recovery.
- The runbook renders `E -> L -> migrate -> catalog-activate -> A -> validate` as one linear
  chronology, names the only permitted delta at each checkpoint, and requests repository-local
  commit authority without treating a CLI flag as Git authority.
- Catalog/portfolio references use the same disposition, proof, readiness, and offline vocabulary
  as the CLI/schema.
- Register doctrine forbids edits and explains deferred-row visibility.
- Routing identifies the owning migration route before Git/provider execution and names the new
  evidence preconditions.
- A fresh low-context probe selects `planning.migrate-catalog`, requests reviewed v3 evidence, and
  does not recommend branch deletion or a tool before routing.
- A missing/older-than-2.43 Git executable routes as local tooling readiness; provider auth,
  source access, and live ref discovery remain portfolio-owned service preflight.
- Source and packaged managed resources are byte-identical.
- Impact and release-note inputs are public-safe and contain only sanitized scenarios.

### E) Dependencies And Critical-Path Notes

Depends on final EP-03/EP-04 behavior and output names. Doctrine must not precede the executable
contract and drift from it.

### F) Tasks Checklist

- [ ] Update `plan-catalog-format.md` with v3 evidence, dispositions, mixed visibility, canonical readiness, compatibility, and invalidation rules.
- [ ] Update `planning-document-lifecycle.md` with the config-v2 ledger binding and mandatory incremental follow-up semantics.
- [ ] Update `portfolio-coordination-format.md` with remote proof, corroboration, readiness, cache-version, and offline boundaries.
- [ ] Rewrite the affected procedure in `migrate-plan-catalog.md` as an agent-facing L1 recipe over `plans inventory`, `plans migrate`, and guarded `plans catalog-activate`.
- [ ] Document separate ledger-only `L` and activation `A` commit approvals, exact first-parent/tree-delta checks, mandatory first post-commit validation, and reinventory after squash/rebase/extra-path drift.
- [ ] Update `maintain-plan-register.md` while excluding register edits, second cutovers, and generated register rows.
- [ ] Update the `planning.migrate-catalog` route card with ledger v3, proof, overlay, approval, blocker, and evidence fields.
- [ ] Route missing/old Git through `handle-tooling-readiness.md`, keep provider auth/access in portfolio preflight, and avoid embedding package-manager instructions in the migration runbook.
- [ ] Record the explicit non-promotion decision that no L2 script wrapper is introduced.
- [ ] At separately authorized implementation start, create the sibling execution log and record epic/task state, commands, review evidence, blockers, and approval boundaries without raw sensitive ref data.
- [ ] Add sanitized consumer-impact and release-note input attachments to the plan bundle.
- [ ] Update component resource identities and byte-identical Python resource mirrors.
- [ ] Update nearest documentation indexes for changed command, schema, and runbook discoverability.
- [ ] Add static routing assertions for clearance authority, mixed deferral, canonical zero-gap, external-action boundaries, and branch-retention safety.
- [ ] Run a fresh low-context migration-routing probe and record its non-secret result in the execution log.
- [ ] Run `python -m pytest tests/test_packaging_resources.py tests/test_routing.py -q` and source-to-packaged byte comparisons.

### G) Implementation Notes

The runbook audience remains `agent-facing`. It already carries a compact intention block and must
retain source/input/precondition/lane/approval/stop/evidence/recovery sections. This is a material
runbook change, not a documentation appendix.

Public examples use placeholder repository IDs, refs, plan paths, commits, and hashes. The two
external acceptance cohorts remain sanitized counts and categories, without consumer names,
private branch names, local paths, or provider identifiers.

### H) Open Questions

None.

## EP-06 - Exhaustive Validation, Security Review, And Acceptance Proof

### A) Epic ID, Title, And Outcome

`EP-06 - Exhaustive Validation, Security Review, And Acceptance Proof`

Outcome: the full required matrix proves correct classification, fail-closed invalidation,
cross-platform reproducibility, visibility, atomic rollback, and the two sanitized acceptance
projections without operating consumer repositories.

### B) Scope

- Add missing unit/integration/CLI/schema/fixture/security tests.
- Retain all existing branch, migration, portfolio, transaction, and platform tests.
- Add deterministic golden evidence across operating systems.
- Run reviewer passes over security and execution readiness.
- Produce sanitized acceptance evidence.

### C) Files Touched

```text
internal/plancatalog/
├── branch_reconciliation_test.go
├── migration_v3_test.go
├── migration_v2_test.go
├── ep02_test.go
└── classifier_test.go
internal/commands/commands_test.go
internal/cli/cli_test.go
internal/portfolio/portfolio_test.go
tests/test_json_schemas.py
tests/test_routing.py
tests/test_packaging_resources.py
tests/fixtures/plans/
tests/fixtures/portfolio/
.github/workflows/validate.yml
docs/repo/plans/branch-ownership-reconciliation/attachments/
└── acceptance-evidence.md                 # create: sanitized synthetic results
```

### D) Acceptance Criteria And Size

- Size: `XL`.
- Every matrix row below has a stable positive/negative assertion and zero-write assertion for
  failed proofs.
- Linux, macOS, and Windows produce identical transition/evidence digests and disposition counts
  from equivalent Git graphs.
- Security review finds no Git content execution, path escape, ref injection, secret capture,
  stale-proof acceptance, partial apply, or cache authority confusion.
- Scenario A retains two deferred owners and all proven non-owners.
- Scenario B yields exactly 12 same-content, 6 incorporated-history, 1 deferred owner, 0
  unreviewed gaps, mixed coverage complete, and canonical ready false.
- Direct canonical projection for both active-owner scenarios fails closed.
- No branch is deleted or modified by the test workflow.

### Required test matrix

| Area | Scenario | Expected result |
| --- | --- | --- |
| Same content | Local branch at different history, exact same path/mode/blob/hash | Reviewable `same-content-non-owner`. |
| Same content | Remote-only branch under a complete overlay, exact same candidate state | Reviewable `same-content-non-owner`; exact remote tip bound. |
| Same content | Equal normalized text but different blob bytes | Blocking; line-ending/whitespace normalization cannot satisfy exact identity. |
| Review identity | One ref touches two unrelated plan candidates and only one has proof | Proven candidate is reviewable; the other remains independently blocking. |
| Incorporated history | Squash commit reachable from target with exact transition and aggregate stable patch | Reviewable `incorporated-history-non-owner`. |
| Incorporated history | Merge-produced reachable commit with exact parent/commit transition | Reviewable when the single-commit contract is satisfied. |
| Partial collision | One matching hunk but additional ref path/blob transition | Blocking `partial_incorporation`. |
| Patch collision | Equal stable patch ID with unequal transition digest | Blocking `transition_digest_mismatch`. |
| Patch pipeline | Global/repository diff config changes, rename threshold changes, and hostile textconv | Identical inert text proof output; no configured command executes. |
| Attributes | Effective attribute stack changes, custom diff driver, working-tree encoding, or binary/NUL content | Exact attributes are bound; unsupported incorporated proof blocks. |
| Git version | Git 2.43.0 and every declared release-runner version | Fixture patch/digest equality; older or unparsable versions block before proof. |
| Ancestry | Ref tip is already an ancestor of target | Not an active touch; ancestry evidence recorded in summary. |
| Stale tracking | Old local tracking ref with exact same-content proof | Reviewable from Git proof; age is non-authoritative. |
| Stale tracking | Old local tracking ref without exact proof | Blocking despite age and merged PR corroboration. |
| Ref drift | Ref tip moves after review and keeps the same touched path set | Blocking `branch_ref_moved` before writes and on idempotent apply. |
| Base drift | Merge base changes after target/ref movement | Blocking and new inventory required. |
| Rename | Exact old/new transitions match incorporated commit | Reviewable incorporated history. |
| Rename | New path differs, is excluded, or collides portably | Blocking. |
| Delete | Exact historical deletion transition is incorporated | Reviewable incorporated history; deletion-only active deferral remains forbidden. |
| Dirty state | Target plan has staged, unstaged, intent-to-add, untracked collision, or filter-modified bytes | Blocking; clearance cannot override dirty authority. |
| Uncommitted owner | Active work exists only in worktree | Blocking and excluded from ref proof. |
| Ledger checkpoint | `L^1=E` and `E..L` is exactly the reviewed ledger path/blob | Migration may run from clean `HEAD=L` while evidence remains bound to `E`. |
| Ledger checkpoint | Extra path, amended ledger, squash/rebase, merge parent, or candidate change in `E..L` | Blocking checkpoint mismatch; new evidence required. |
| Mixed deferral | One active ref-tip plan, exact frozen target bytes, complete follow-up | `deferred-active-owner`, mixed coverage complete, canonical ready false. |
| Mixed deferral | Two active owners | Blocking ownership conflict. |
| Mixed deferral | Owner lifecycle is not active | Blocking invalid owner lifecycle. |
| Mixed deferral | Frozen register, cutover revision, path, or bytes move | Blocking exact mixed proof failure. |
| Canonical | Any deferred active owner | Blocking `active_owner_deferral_requires_mixed` and nonzero gap. |
| Visibility | Mixed activation with one deferred plan | Local and remote compatibility rows retain title/lifecycle/path/hash and follow-up flag. |
| Activation drift | Ref/overlay/candidate/ledger moves after metadata apply but before activation | `plans catalog-activate` performs zero config writes. |
| Activation dirty set | Index equals `L`; worktree contains exactly projected metadata outputs, including expected untracked/removals | Guarded activation may write config; staged or additional/mismatched delta blocks. |
| Activation checkpoint | `A^1=L` and `L..A` is exactly projected actions plus bound config | First clean post-commit `plans validate` succeeds and identifies `A`. |
| Activation checkpoint | Ref moves after config write but before `A` | First post-commit validation fails closed; compatibility evidence is not credited. |
| Activation checkpoint | Extra path, squash/rebase, non-first-parent insertion, or changed config in `L..A` | Binding invalid; reinventory/recovery required. |
| Bound validation | Ref moves after successful mixed activation | Config binding becomes invalid; validate/list/portfolio cannot credit the deferred row as current. |
| Bound validation | Unrelated commits descend from valid `A` without changing bound authority | `A` remains discoverable/ancestral and validation continues normally. |
| Incremental | Owner merge changes legacy bytes without metadata | Validation fails and requires new inventory/ledger. |
| Incremental | Owner merge supplies valid metadata | Follow-up resolves and row becomes canonical. |
| Incremental | Post-merge reviewed metadata migration | Apply succeeds, then canonical projection reaches zero gaps. |
| Offline | Local proof complete, overlay not requested | Migration works; remote freshness and portfolio completeness remain unclaimed. |
| Remote overlay | Ledger binds complete overlay | Dry-run/apply require matching fresh overlay. |
| Remote overlay | Auth/fetch/member/ref failure | Migration blocks and prior complete cache is preserved. |
| Remote overlay | Discovery-v1 default branch with ledger targeting prospective v2 | In-memory target-member scan succeeds without config mutation. |
| Remote overlay | Default is exact pushed `A` or an unchanged later descendant | Scanner verifies remote `E -> L -> A`, keeps proof anchored to `E`, and accepts current evidence. |
| Remote overlay | Default squashes/rebases checkpoint, lacks `E`/`L`/`A`, or diverges from `A` | Member evidence incomplete/blocking; no silent proof rebase and prior cache preserved. |
| Remote overlay | Unrelated portfolio member is incomplete | Target member may proceed when `member_evidence_complete` is true; portfolio `complete` remains false. |
| Remote overlay | Tracking ref and mirror normalize to one source | Coalesce only on exact logical ref, tip, base, and transitions; mismatch blocks as two records. |
| PR corroboration | Merged PR plus missing Git proof | Blocking; PR is not sole evidence. |
| Schema | V3 unknown field/disposition/proof shape | Rejected before planning. |
| Compatibility | New CLI with v1/v2 ledgers | Released behavior preserved. |
| Compatibility | v0.1.25/v0.1.26 CLI with v3 ledger | Unsupported schema; zero writes. |
| Compatibility | v0.1.25/v0.1.26 CLI with config v2 | Config invalid before catalog settings are used; zero writes. |
| Security | Malicious ref/path, `.gitattributes`, diff driver, textconv, replace object, promisor object | No code/network execution; stable blocker. |
| Tamper | Tip/base/blob/hash/transition/patch/incorporated commit/follow-up/overlay or `E`/`L`/action-digest field changes | Proof/checkpoint recomputation blocks. |
| Atomicity | Ref moves after prewrite, between two actions, after the action batch, or before cleanup | Post-write authority drift, rollback, prior bytes restored. |
| Atomicity | Write/post-check/rollback failure | Existing recovery marker and quarantine contract retained. |
| Cross-platform | Same synthetic graph on Linux/macOS/Windows | Identical evidence digest, patch ID, disposition, projection, and JSON ordering. |
| End to end | Clean `E` → ledger-only `L` → migrate → guarded activate → exact `A` → clean validate | No self-reference; each revision/delta is identified and verified linearly. |

### E) Dependencies And Critical-Path Notes

Depends on all implementation and doctrine epics. Existing focused package tests must pass before
the full repository suite. Real Windows proof remains a release-candidate gate, not a cross-compile
substitute.

### F) Tasks Checklist

- [ ] Implement every required matrix row as a deterministic unit, integration, CLI, schema, fixture, security, and transaction test.
- [ ] Preserve and rerun all existing v1/v2 branch touch, mixed baseline, migration, overlay, cache, and rollback tests.
- [ ] Add sanitized Scenario A and Scenario B synthetic repositories with exact expected disposition counts and projections.
- [ ] Add one full `E -> L -> A -> plans validate` fixture plus negative extra-path, squash/rebase, merge-parent, staged-drift, and post-activation-ref-movement variants.
- [ ] Add golden canonical JSON and evidence-digest fixtures shared across Linux, macOS, and Windows jobs.
- [ ] Add command-policy assertions that forbid hooks, diff drivers, textconv, helpers, replace objects, lazy fetch, and interactive prompts, and cover attributes, binaries, renames, and admitted Git versions.
- [ ] Add old-binary compatibility tests using released v0.1.25 and v0.1.26 binaries against a v3 ledger and config-v2 disposable fixture.
- [ ] Run `go test ./...` on Linux, macOS, and real Windows.
- [ ] Run `python -m pytest tests/test_json_schemas.py tests/test_routing.py tests/test_packaging_resources.py tests/test_sync_check.py -q`.
- [ ] Run `python scripts/validate-json-schemas.py`, `python scripts/validate-markdown-headers.py`, `python scripts/validate-public-core.py`, and `python scripts/validate-release-manifest.py`.
- [ ] Record sanitized acceptance counts, commands, platform results, and residual limits in `attachments/acceptance-evidence.md`.
- [ ] Run a security-focused reviewer pass over proof sufficiency, Git inertness, stale evidence, authority digest, output redaction, and rollback.
- [ ] Run the planning-document review route and resolve every High and Medium execution-readiness finding before activation.

### G) Implementation Notes

Synthetic fixtures must use generated placeholder refs and hashes. Cross-version binary tests run
only in disposable fixture repositories and must not download during ordinary offline unit tests;
release binaries become pinned CI inputs through a reviewed test harness.

Stable patch output is not itself a cross-platform success criterion unless the exact transition
digest also matches. Both values must be asserted together.

### H) Open Questions

None. Any newly discovered proof ambiguity is fail-closed and requires a plan revision before the
affected case is supported.

## EP-07 - Versioned Release, Publication Gate, And Consumer Handoff

### A) Epic ID, Title, And Outcome

`EP-07 - Versioned Release, Publication Gate, And Consumer Handoff`

Outcome: after separate authority, the next patch release contains the validated implementation,
managed resources, schemas, migration notes, reproducible macOS/Windows packs, and public proof;
consumer repositories receive a safe staged rollout sequence without automatic upgrades.

### B) Scope

- Recheck baseline and select next patch.
- Update producer version/content/release surfaces.
- Build and validate reproducible release assets.
- Stop before tag/publication without explicit authority.
- Verify public assets after authorized publication.
- Produce a repository-owned consumer rollout handoff; do not operate consumers in this epic.

### C) Files Touched

```text
pyproject.toml
src/codeheart_operating_kit/__init__.py
internal/version/version.go
manifest.yaml
src/codeheart_operating_kit/resources/manifest.yaml
profiles/standard.yaml
components/planning-workflows/component.yaml
bootstrap.md
install.sh
install.ps1
release-notes.md
tests/test_go_cli_parity.py
tests/test_release_assets.py
tests/test_install_metadata.py
.github/workflows/validate.yml
docs/repo/plans/branch-ownership-reconciliation/attachments/
└── release-readiness-evidence.md            # create
```

### D) Acceptance Criteria And Size

- Size: `L`.
- Version is the next available patch selected from live release truth.
- Release notes explain ledger/catalog v3, guarded config v2, old-CLI rejection,
  upgrade-before-use, offline/overlay behavior, mixed follow-up, canonical zero-gap, and rollback.
- Managed content identities and manifest chains match source bytes.
- macOS universal and Windows x64 packs build twice byte-identically.
- External catalog, archive, pack manifest, payload checksums, content identity, binary digest, and
  version chain verify.
- Fresh install and v0.1.25/v0.1.26 upgrade paths preserve existing config v1 and do not activate
  discovery v2, create a ledger, or add a deferral binding. Manifest compatibility advertises
  config versions 1 and 2 only after both dispatch paths pass.
- Failed verification/install/upgrade preserves the previous installation.
- Tag/release/publication occur only under explicit authority and target the validated commit.
- No consumer repository is modified by producer release execution.

### E) Dependencies And Critical-Path Notes

Depends on EP-06 passing evidence. Release versioning and publication are separately gated by the
producer release runbook and current user instruction.

### F) Tasks Checklist

- [ ] Recheck the latest release tag, target commit, release authority, consumer-impact record, and signing boundary.
- [ ] Record explicit version and publication authority before editing version and public release surfaces.
- [ ] Set the next patch version consistently in producer metadata, binary metadata, manifests, profiles, installers, bootstrap, and notes.
- [ ] Add consumer migration notes for ledger v3, config v2, old-CLI rejection, retained v1/v2 behavior, exact `E -> L -> A` checkpoints, guarded activation, mixed follow-up, and canonical cutover.
- [ ] Run public-core, Markdown, schema, content-identity, Go, Python compatibility, installer, and release-contract validation.
- [ ] Build macOS universal and Windows x64 packs twice and require byte-identical output.
- [ ] Verify the external catalog through archive, pack manifest, payload, content identity, binary digest, and version chain.
- [ ] Run isolated macOS and real-Windows fresh-install, v0.1.25 upgrade, v0.1.26 upgrade, failure, and rollback paths.
- [ ] Present the validated commit, version, asset digests, platform evidence, signing state, residual risk, and publication boundary.
- [ ] Stop before tag creation and publication until explicit authority covers the presented release candidate.
- [ ] Publish the authorized tag, packs, catalog, installers, notes, and checksums through the normal non-force workflow.
- [ ] Verify live public URLs, sidecar hashes, install paths, version output, and release-channel behavior.
- [ ] Record release-readiness and publication evidence without private consumer details.

### G) Implementation Notes

Plan activation does not authorize EP-07 publication. A release rollback cannot use an older CLI to
apply ledger v3. Pre-adoption consumers may stay on v0.1.26; consumers that have begun v3 evidence
work require a forward fix or reviewed config-only recovery while plan metadata remains preserved.

### H) Open Questions

OQ-04 resolves version selection at execution time from live release truth.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Multi-commit incorporated-history proof is deferred until an evidence-backed case justifies a
  versioned sequence contract. Trigger: one complete ref candidate transition is provably split
  across more than one target commit and cannot use exact same-content proof.
- Multi-ledger consolidation is deferred. Trigger: a repository has more than one concurrently
  pending deferral ledger and a single config-v2 binding cannot represent the reviewed rollout;
  until then, activation blocks rather than merging authorities.
- Automated ref deletion, pruning, archival, and branch cleanup are excluded. Trigger: a separate
  repository-owned retention policy and explicit destructive-action authority.
- Provider-specific PR enforcement is deferred. Trigger: a provider contract is required for
  corroboration display; it must remain non-authoritative for offline clearance.
- Named consumer upgrades and catalog migrations are separate repository-owned tasks after a public
  release. Trigger: explicit target-repository authority, clean worktree preflight, compatible Kit
  installation, reviewed inventory, and local consumer plan.

## 4.2 Future Considerations

### Consumer rollout sequence

1. Release the producer patch only after EP-07 authority and proof.
2. Upgrade one disposable/pilot repository and prove that existing config, ledger v2, catalog
   mode, discovery version, plan bytes, refs, and cache remain unchanged until commands are
   explicitly invoked.
3. In each authorized consumer, establish clean evidence revision `E`, create a prospective
   inventory with local evidence first, add a complete remote overlay when configured, and retain
   both artifacts outside plan targets.
4. Review every ref touch and build a repository-owned ledger v3 against `E`; never copy raw
   consumer evidence into this public producer. Under separate local commit authority, create the
   ledger-only checkpoint `L` with first parent `E` and validate its single-path delta.
5. From clean `HEAD=L`, run migration dry-run, compare exact
   actions/dispositions/follow-ups, and apply only under local migration authority. Do not commit
   another path before guarded activation.
6. For active owners, run guarded `plans catalog-activate` dry-run/apply only when the final
   projection reports mixed coverage complete and canonical ready false with no unreviewed gap.
   Under separate local commit authority, create activation checkpoint `A` with first parent `L`
   and exactly the projected action paths plus config; immediately require clean post-commit
   `plans validate` to prove the config-v2 binding and `E -> L -> A` history. After separately
   authorized normal push/integration, require a fresh remote member scan to prove `A` ancestry;
   a squash or rebase is a blocker, not a reason to rewrite the evidence.
7. After owner branches merge, require validation immediately. Accept owner-supplied metadata or
   run reviewed incremental migration against the new target bytes.
8. Use guarded zero-gap direct-canonical activation only after a new reviewed projection and
   complete required remote evidence; write config v1 without a deferral binding and validate its
   exact activation checkpoint before treating the member as canonical-ready.
9. Roll out to remaining repositories one at a time, preserving historical refs and recording
   repository-local rollback evidence.

### Rollback and recovery

- Before ledger checkpoint `L`: retain current mode/config/plans and discard only stale unapproved
  projections; regenerate evidence from current refs.
- After `L` but before metadata apply: retain the ledger-only commit as historical review evidence;
  moved authority requires a new `E`/ledger/`L` chain, not amendment, rebase, or hash repair.
- After metadata apply but before activation checkpoint `A`: keep exact reviewed metadata outputs
  and discovery v1 config while correcting evidence through a new ledger. A stale guarded
  activation writes nothing. If config was written but `A` has not been committed, use the guarded
  recovery path to restore its exact pre-activation bytes; do not commit an invalid checkpoint.
- After mixed activation checkpoint `A`: on bound-evidence validation failure, use the guarded config-only
  recovery path to restore the exact pre-activation config-v1 bytes when still semantically valid,
  while preserving ledger, plan metadata, frozen register, refs, and Git history. Never remove the
  binding through an unreviewed hand edit.
- After owner merge: a changed filename-only plan is not rolled back to cutover bytes. Inventory
  the merged truth and complete incremental metadata migration.
- On transaction failure: retain marker, backup, quarantine, and recovery evidence and follow the
  existing `check`/repair path.
- On remote failure: preserve the last complete compatible cache, disclose its timestamp as
  historical, and retry after access/freshness recovers.

### Approval decisions before implementation activation

No unresolved blocker prevents a single-path implementation. An unambiguous future request to
activate this plan authorizes only the bounded plan-status, execution-log, direct metadata, work
branch, plan-only commit, and normal-push checkpoint defined by the planning workflow. It excludes
implementation code, unrelated files, another plan, PR creation, merge, release, force-push,
branch deletion, history rewrite, destructive Git, credential changes, and rejected-push bypass.
After that checkpoint, a separate unambiguous request to execute the active plan may authorize
EP-01 through EP-06 producer changes on that branch. EP-07 versioning/publication and every
consumer rollout remain separately approval-gated by current instruction, repository policy, and
their owning runbooks.

# Revision Notes

- 2026-08-09: Created the draft implementation plan from the verified v0.1.26 runtime/schema/test
  baseline and two sanitized external acceptance scenarios. Selected strict ledger v3, exact
  ref/path-state binding, collision-safe incorporated-history proof, mixed-only active-owner
  deferral, explicit compatibility visibility, and separate release/consumer authority.
- 2026-08-09: Resolved planning-review findings by selecting a strict config-v2 deferral binding
  and guarded catalog activation, adding rollback-capable post-write authority checks, fixing the
  candidate-scoped review identity and normative ledger shape, specifying the inert patch
  pipeline and Git boundary, defining prospective member-scoped remote-overlay completeness, and
  replacing the circular ledger/activation revision model with exact `E -> L -> A` checkpoints and
  mandatory first post-commit validation.
- 2026-08-09: Completed the repository planning-document review route with result `Ready`, no
  remaining High or Medium findings, and residual implementation complexity limited to the
  explicitly specified Git-history and remote-mirror verification matrix.
- 2026-08-09: Activated under explicit user authority for complete producer implementation and
  v0.1.27 release execution. The bounded activation checkpoint contains only this plan, its
  sibling execution log, and directly required documentation indexes; implementation source,
  release, PR, merge, tag, and publication work follows under the separately explicit authority
  in the same request.
