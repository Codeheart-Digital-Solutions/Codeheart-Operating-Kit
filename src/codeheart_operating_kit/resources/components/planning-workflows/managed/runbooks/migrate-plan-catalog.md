Last updated: 2026-08-11T14:12:23Z (UTC)

# Migrate Plan Catalog

Audience: agent-facing

Intent:
Route one repository through prospective discovery-v2 inventory, candidate-scoped branch-touch
review, guarded semantic migration, and explicit catalog activation without deleting historical
refs, changing another repository, or hiding a genuine active owner.

Success:
The repository follows `E -> L -> migrate -> catalog-activate -> A -> validate`; every branch touch
has current reviewed schema-v3 evidence; only proved non-owners are cleared; any genuine active
deferral remains visible in mixed mode with an exact pending incremental follow-up; direct
canonical activation has zero gaps; and every checkpoint contains only its permitted delta.

Agent judgment boundary:
Interpret plan meaning, ownership, and optional metadata from repository evidence. Inventory may
propose dispositions but never approves them. Do not invent proof, treat a merged pull request as
sole clearance, weaken a blocker, delete a branch to make migration pass, edit the frozen register,
or treat a CLI `--yes` as Git commit/push authority.

Stop boundary:
Stop on missing or Git older than 2.43, ambiguous semantics/ownership, stale or unsupported proof,
changed ref/candidate/policy/source/config/overlay/ledger/output, dirty or staged overlap, invalid
checkpoint chronology, incomplete required remote preflight, recovery-required transaction state,
or any unapproved local/external write.

## Source, Inputs, Preconditions, And Lane

Source of truth, in order: current user and repository-local authority; current plan and sibling
content; Git index/tree/ref/history evidence; the one frozen legacy cutover/register when mixed is
used; a reviewed schema-v3 ledger; `../reference/plan-catalog-format.md`; current schema and CLI
structured output; and portfolio preflight for remote-aware facts.

Inputs: repository path; public-safe repository ID; target mode (`canonical` or `mixed`); explicit
inventory output outside every candidate, ledger, config, and register path; normalized tracked
ledger path outside all migration targets and protected paths; reviewed exclusions; evidence scope
(`local` or `remote-aware`); branch-owner decisions; exact intended action paths; and separately
scoped approvals for metadata/config writes and the `L` and `A` repository commits.

Preconditions:

- Route the request as `planning.migrate-catalog` before selecting Git, provider, or CLI actions.
- Use Git 2.43 or newer. If Git is absent or older, stop this recipe and route to
  `../../agent-interface/runbooks/handle-tooling-readiness.md`; do not embed package-manager
  instructions here.
- Require an initialized compatible Kit, no unresolved transaction/recovery marker, a regular
  config, and exact understood ownership of every touched repository path.
- Require a clean index and worktree at evidence commit `E`. Uncommitted owner content is not
  evidence and blocks migration.
- For mixed deferral, establish the single reviewed legacy cutover before `E`; its regular register
  and deferred plan bytes must already be frozen. Do not edit the register or create a second
  cutover during this route.
- For `remote-aware` scope, route provider authentication, source access, pagination, mirror/live-
  ref discovery, and complete target-member overlay through the portfolio preflight in
  `refresh-portfolio-catalog.md`. Do not downgrade to local scope after review.

Execution lane: L1 semantic-review orchestration over the L3 `plans inventory`, `plans migrate`,
`plans catalog-activate`, `plans validate`, and `plans list` commands plus repository-owned Git
checkpoint operations. No L2 wrapper is introduced; the commands already own the bounded
deterministic mechanics while this runbook owns judgment, approvals, and chronology.

Approval boundary:

- Inventory, validation, list, and dry-run are evidence reads except for the explicit inventory
  artifact and scanner-owned local portfolio cache/mirrors.
- `plans migrate --yes` authorizes only the enumerated metadata/rename worktree transaction.
- `plans catalog-activate --yes` separately authorizes only the guarded config transaction.
- Creating commit `L` requires explicit repository-local authority for the exact ledger-only path.
- Creating commit `A` requires separate repository-local authority for the exact migration action
  paths plus `.codeheart/kit.config.yaml`.
- None of those approvals authorizes staging unrelated work, push, provider write, PR, merge,
  release, ref deletion, branch deletion, force-push, rebase, history rewrite, or destructive Git.

## Procedure

### 1. Route And Establish Evidence Revision `E`

1. Select `planning.migrate-catalog`. Inspect current config/mode/discovery version, plan/register
   authority, repository status/index, refs, transaction state, portfolio role, and current CLI
   validation before changing anything.
2. Verify Git is 2.43 or newer. Route missing/old tooling as stated above. For remote-aware work,
   complete portfolio service preflight before inventory.
3. Resolve exclusions and genuine plans under ambiguous fixture/vendor/generated/dependency/build-
   like roots. Move a genuine plan through normal repository authority; exclude only proven non-
   authority. Preserve the inventory observation either way.
4. Require a clean index/worktree and record full commit `E`. For mixed deferral, `E` already
   contains the one valid legacy cutover/config/register baseline. Do not create evidence from
   worktree bytes, local-only owner edits, or a moving ref.

### 2. Inventory And Review Every Candidate

Create the explicit output parent, then take one coherent schema-v3 inventory. Direct canonical:

```sh
codeheart-operating-kit plans inventory \
  --target-discovery-version 2 \
  --target-catalog-mode canonical \
  --output <inventory.json> \
  --json \
  <repository>
```

For mixed deferral, the configured mode at `E` is already `mixed`; omit
`--target-catalog-mode canonical`. Add `--remote-overlays` only for the selected remote-aware scope,
and then retain it on every later migrate, activation, and validation call.

Review every owned, excluded, prospective-blocked, hard-unowned, filename-only, metadata-only,
malformed, renamed/deleted-path, duplicate, family, dirty, and branch-touch observation. Require
schema version 3, exact `evidence_revision: E`, policy and candidate-set digests, config
precondition, complete path states, and `git-candidate-proof-v1` branch evidence.

For a direct canonical target, frozen-register `legacy_status_unsupported`,
`legacy_relation_malformed`, and `legacy_relation_unsupported` observations may be warnings only:
the reviewed ledger must supply the current lifecycle and relations, and the register stays
unchanged as historical evidence. The same codes remain errors in legacy and mixed modes. Every
identity, canonical-path, duplicate, date, field, readability, safety, or reconciliation error
remains blocking in canonical mode too. Do not broaden this compatibility rule or edit the register
to clear an observation.

For each branch touch choose exactly the reviewed disposition justified by recomputable evidence:

- `same-content-non-owner` only for exact target/ref mode, object, and SHA-256 identity;
- `incorporated-history-non-owner` only for a complete candidate transition whose digest and stable
  patch ID equal a specific parent/commit transition reachable from `E`;
- `deferred-active-owner` only for one exact active owner in mixed mode with frozen cutover
  authority and the mandatory follow-up;
- `active-owner` when ownership is genuine but deferral is not authorized; or
- `blocking` when evidence is absent, ambiguous, partial, unsupported, conflicting, or stale.

Provider merged-PR state is optional corroboration, never the sole offline requirement. Partial
patch collision, matching a subset of paths, unreachable incorporated commit, unknown attributes,
binary/invalid-UTF-8 proof, missing merge base, stale tracking ref, or a moved/deleted ref blocks.
Keep historical branches; branch deletion is not a migration precondition and must not be
recommended as clearance.

### 3. Author, Review, And Checkpoint Ledger `L`

Author one strict `schemas/plan-migration-ledger.schema.json` schema-v3 ledger. Bind repository ID,
`E`, target mode, inventory revision, policy/candidate/config digests, evidence scope, remote
overlay digest, and `remote_source_identity_sha256` when applicable, every candidate decision, and every candidate/ref review identity,
tip, merge base, path transition, proof, blocker, owner snapshot, and follow-up required by its
disposition. The ledger path is a normalized repository-relative tracked regular file, normally a
plan attachment, outside migration targets, config, and the frozen register.

For a deferral, require all of:

- unchanged current/target plan path and exact cutover source hash;
- exactly one `active` owner-tip path/mode/object/SHA-256 snapshot;
- `action: reinventory-and-incremental-migrate`;
- `trigger: owner-ref-integrated-or-target-path-changed`;
- `state: pending`; and
- both supported success conditions.

After semantic and schema review, request explicit repository-local authority to create ledger-only
commit `L`. The resulting checkpoint must satisfy all of these facts:

```text
L^1 = E
tree-delta(E, L) = {<ledger-path>}
blob(L:<ledger-path>) SHA-256 = reviewed ledger_sha256
HEAD = L; index tree = L; worktree clean
```

No CLI flag grants this commit authority. An amended/squashed/rebased ledger checkpoint, wrong
parent, extra path, non-regular ledger, or changed ledger byte requires reinventory from a new `E`
and renewed review; do not repair hashes in place.

### 4. Dry-Run And Apply Migration From Clean `L`

Run:

```sh
codeheart-operating-kit plans migrate \
  --ledger <reviewed-ledger.yaml> \
  --dry-run \
  --json \
  <repository>
```

Add `--remote-overlays` for remote-aware evidence. Require verified branch reviews, zero material
blockers/skips, exact action list, the correct projected mode, `mixed_coverage_complete: true`,
`canonical_ready: true` for direct canonical or `false` only for exact mixed deferrals, every
pending follow-up, action digest, and `activation_performed: false`.

With separate metadata-write approval, rerun using `--yes` instead of `--dry-run`. The transaction
may change only exact projected metadata/rename action paths, must leave the index at `L`, and must
recompute external authority before and after writes. A ref/overlay/candidate movement during the
transaction rolls back every action. A zero-action/idempotent result still rechecks authority.

### 5. Dry-Run And Apply Guarded Catalog Activation

Without staging migration outputs, run:

```sh
codeheart-operating-kit plans catalog-activate \
  --ledger <reviewed-ledger.yaml> \
  --dry-run \
  --json \
  <repository>
```

Add `--remote-overlays` for remote-aware evidence. The command must require `HEAD=L`, index tree
`L`, and the complete unstaged worktree delta byte-for-byte equal to the migration action digest.
Any extra staged, intent-to-add, untracked, renamed, deleted, or modified path blocks. Review the
one config action, evidence/config preconditions, branch reviews, follow-ups,
`activation_performed: false`, and readiness.

With separate config-write approval, rerun with `--yes`. It writes only the reviewed catalog config
fields and returns `activation_performed: true` plus `activation_checkpoint_required: true`.
Mixed deferral writes config schema v2 and the exact ledger/evidence/action binding. Zero-gap direct
canonical writes config schema v1 without a deferral binding. The command does not stage or commit.

### 6. Checkpoint Activation `A` And Validate Immediately

Request separate repository-local authority for commit `A`. Include exactly the projected
metadata/create/remove paths plus `.codeheart/kit.config.yaml`; do not edit or recommit the ledger,
register, or unrelated work. The checkpoint must satisfy:

```text
A^1 = L
tree-delta(L, A) = {exact migration-action paths, .codeheart/kit.config.yaml}
HEAD = A; index tree = A; worktree clean
```

The first operation after the clean `A` commit is validation:

```sh
codeheart-operating-kit plans validate --json <repository>
codeheart-operating-kit plans list --format json <repository>
```

Add `--remote-overlays` to validation for remote-aware evidence. Require the config binding to
locate and verify `E`, `L`, and `A`, exact first-parent/tree deltas, ledger/ref/proof/overlay/action
hashes, `complete`, `mixed_coverage_complete`, `canonical_ready`, and deferred compatibility rows.
Each deferred plan appears once as `mixed-grandfathered` with
`incremental_migration_required: true`; canonical mode has no such row.

A squash, rebase, non-first-parent insertion, extra checkpoint path, or movement after activation
invalidates the proof. Reinventory from a new clean `E`; do not rewrite the binding manually.

### 7. Incremental Migration And Portfolio Visibility

When an owner ref is integrated or a deferred target path/bytes change, stop claiming canonical
readiness. Reinventory the merged target and current refs. Close the pending follow-up only when
owner-supplied canonical metadata at merge validates or a new reviewed ledger is applied to the
merged bytes through this chronology. Do not update the frozen register or create a second cutover.

Publishing `A`, refreshing a coordination home, or changing another repository requires its own
authority and route. Until pushed/default-branch evidence is refreshed, report remote visibility as
unknown or stale rather than current.

## Stop Conditions And Structured Blockers

Stop on unknown disposition, missing/unreviewed touch, missing/moved ref, missing/ambiguous/moved
merge base, moved path state, blob/hash/transition/patch mismatch, unreachable incorporated commit,
partial incorporation, unsupported Git/attributes/binary proof, multiple active owners, invalid
active lifecycle, dirty owner, missing mixed/follow-up proof, stale overlay, canonical deferral,
checkpoint parent/delta/blob mismatch, activation authority drift, post-write drift, or rollback
failure. Preserve current evidence and do not delete refs or weaken catalog mode.

Record each blocker with: stable code/class; recipe phase; repository and sanitized plan/ref target;
non-secret message; evidence command/result marker; preserved state; safe retry/recovery route; and
the user, branch-owner, tooling, or portfolio decision needed. Never fabricate a cleared review or
success row.

## Evidence And Validation

Retain public-safe inventory and reviewed ledger; `E`, `L`, and `A`; exact first-parent/tree-delta
checks; ledger/config/source/target hashes; policy/candidate/branch/action/overlay digests; review
identities/dispositions/proof kinds; owner snapshots and pending follow-ups; dry-run/apply JSON;
transaction/rollback evidence; readiness values; compatibility rows; first post-commit validation;
and idempotent recheck. Do not retain tokens, credential-bearing URLs, raw private refs, private
topology, or unnecessary machine paths.

Validation proves complete filename-or-metadata enumeration, strict candidate-scoped review,
same-content and incorporated-history proof without false equivalence, fail-closed movement,
metadata/header chronology preservation, exact migration/config atomicity, no register edit,
mixed deferred visibility, canonical zero-gap, remote/offline scope honesty, checkpoint chronology,
and zero unrelated writes.

## Recovery

On any optimistic mismatch, preserve unrelated work, discard the stale proposed action, and
reinventory current evidence; never edit hashes or delete branches to reuse approval. On
transaction failure, preserve marker, backup, quarantine, and recovery artifacts and follow CLI
`check`/repair guidance. Migration apply rolls back all action paths on post-write authority drift;
activation rolls back config to its exact prior bytes on drift.

If `A` validation fails, do not publish or claim readiness. Preserve the reviewed ledger and
transaction evidence, restore only through a separately reviewed exact config recovery, and start
new evidence when chronology or authority moved. An incomplete remote attempt preserves the prior
complete compatible cache only as historical evidence; rerun portfolio preflight rather than
claiming it is current.
