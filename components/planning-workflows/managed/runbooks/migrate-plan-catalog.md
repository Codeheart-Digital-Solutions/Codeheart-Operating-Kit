Last updated: 2026-08-07T00:26:02Z (UTC)

# Migrate Plan Catalog

Audience: agent-facing

Intent:
Prospectively inventory every repository-wide discovery-v2 candidate, semantically review its
ownership and meaning, apply guarded metadata/rename changes, and explicitly activate v2 without
freezing unrelated work or rewriting historical content dates.

Success:
Every owned candidate has a supported canonical filename and reviewed valid metadata, exclusions
and prospective blockers are resolved, active-branch ownership is reconciled, migration is
idempotent, and explicit activation occurs only after complete local and required remote evidence
while legacy evidence remains preserved.

Agent judgment boundary:
Interpret plan meaning from content and repository evidence; omit unsupported optional metadata.
Do not convert register rows mechanically, invent classifications, edit another branch's owned
plan, modify another repository, or bypass hashes, blockers, dry-run review, and approval.

Stop boundary:
Stop on ambiguous ownership/semantics, invalid inventory/ledger, changed policy/candidate set/
source revision/bytes, dirty overlap, active-branch ownership without coordination, broken mixed
baseline, incomplete required remote scan, recovery-required transaction state, or missing
`--yes` approval.

## Source, Inputs, Preconditions, And Lane

Source of truth, in order: current repository plans and sibling evidence; Git revision/ref/history;
the frozen legacy register; current user/branch-owner decisions; `../reference/plan-catalog-format.md`;
and the migration schemas and CLI output.

Inputs: repository path; explicit inventory artifact path outside every actual or prospective plan
target; reviewed schema-v2 ledger path; repository ID; `excluded_roots`; target mode; branch
ownership/deferral decisions; and approval for each config or migration write.

Preconditions: initialized Kit with discovery-v2-capable CLI; clean or explicitly understood target
paths; no unresolved transaction marker; current Git index/revision evidence; and remote access
when overlay validation is requested. A committed legacy-mode register baseline is required only
when mixed mode will grandfather filename-only records.

Execution lane: Operating Kit `plans inventory`, semantic agent review, `plans migrate`, guarded
config edit, and `plans validate`. Inventory writes only its explicit artifact. Migration requires
exactly one of `--dry-run` and `--yes`.

Recipe maturity: L1 semantic-review orchestration over L3 inventory/migrate/validate commands.
Fresh-agent review proves the judgment and approval flow; schema, command, transaction, and fixture
tests prove deterministic mechanics. Do not create an L2 conversion script that bypasses semantic
review or duplicates the L3 commands.

Approval boundary: inventory and dry-run are evidence operations. User approval is required for
metadata/config/register-notice writes. Approval applies only to enumerated plan and directly
required migration evidence; it does not authorize unrelated code, other repositories, PRs,
merges, releases, force-push, deletion, or destructive Git.

## Procedure

1. Inspect current catalog mode, discovery version, config, Git index/status/branches, transaction
   state, historical register, and portfolio role. Run active-contract validation without changing
   config:

   ```sh
   codeheart-operating-kit plans validate <repository>
   ```

2. If filename-only records will be intentionally deferred through mixed mode, preserve the
   existing original cutover revision or create the one legacy baseline required by mixed doctrine.
   With explicit reviewed config authority, set `plan_catalog_mode: mixed` and that exact
   `plan_catalog_cutover_revision` before migration dry-run/apply, but keep discovery v1 active by
   leaving `plan_catalog_discovery_version` absent or `1`. Validate the v1 mixed baseline. Otherwise
   do not add a register notice, create a cutover revision, or enter mixed mode; the direct
   canonical route keeps config unchanged through migration apply.
3. Review proposed `excluded_roots`. Move genuine plans out of misleading fixture/vendor/generated/
   dependency/build-like boundaries; use exclusions only for non-authority and retain their
   inventory evidence.
4. Create the explicit artifact directory, then run a prospective canonical v2 inventory:

   ```sh
   codeheart-operating-kit plans inventory \
     --target-discovery-version 2 \
     --target-catalog-mode canonical \
     --remote-overlays \
     --output <inventory.json> \
     <repository>
   ```

   Omit `--remote-overlays` only when portfolio overlay evidence is not configured or not in the
   approved scope, and record that limitation.
5. Review every owned, excluded, prospective-blocked, hard-unowned, filename-only, metadata-only,
   malformed, mismatch, duplicate, and family candidate. Reconcile sibling documents, aliases,
   lifecycle, purpose, family, relations, taxonomy, ownership, source revision/hash, target-path
   precondition, dirty overlap, and active remote branch ownership. Exact `README.md` requires
   genuine valid family metadata; directory shape is never authority.
6. Build a schema-v2 ledger matching `schemas/plan-migration-ledger.schema.json`. Bind
   `discovery_version: 2`, repository ID, policy digest, candidate-set digest,
   inventory revision, source hashes, current/target paths, target preconditions, ownership,
   confidence, evidence, ambiguity, conflicts, deferral, and branch owner. Include every owned
   candidate exactly once. Set `target_catalog_mode: canonical` for the normal complete route, or
   `mixed` only when reviewed unchanged register-proven filename-only records are deferred. Omit
   unsupported optional classification; do not guess.
7. Validate the ledger and dry-run:

   ```sh
   codeheart-operating-kit plans migrate \
     --ledger <reviewed-ledger.yaml> \
     --dry-run \
     <repository>
   ```

8. Review every replace/create/remove action, rename, skip, blocker, projected coverage, exact
   policy/candidate/source hash, preserved header, and unrelated-file exclusion. Require
   `activation_performed: false`; a material skip or coverage gap is not success.
9. Apply the exact reviewed ledger while active discovery remains v1. The direct canonical route
   keeps config unchanged; the optional deferral route is already catalog mode mixed with its
   validated original cutover proof:

   ```sh
   codeheart-operating-kit plans migrate \
     --ledger <reviewed-ledger.yaml> \
     --yes \
     <repository>
   ```

10. Reinventory prospectively. Files whose revision/hash changed or are owned by active branches remain skipped;
    ask their branch owners to add reviewed metadata in their current work or reevaluate the newest
    bytes after ownership stabilizes. No repository-wide freeze is required.
11. Reconcile skipped/new records with updated ledger evidence and reapply through the same dry-
    run and approval path.
12. Reapply the final ledger and record zero changes or explicit `already_applied` outcomes.
13. For the direct canonical route, run prospective canonical validation and require zero owned
    gaps:

    ```sh
    codeheart-operating-kit plans validate \
      --target-discovery-version 2 \
      --target-catalog-mode canonical \
      --remote-overlays \
      <repository>
    ```

    For the optional mixed route, the final dry-run/apply projection must instead show that every
    remaining filename-only candidate is exact register-proven grandfathering and that there are
    zero unreviewed gaps. A coordination home may retain version-mismatched member evidence until
    those member default branches activate v2; do not call the portfolio v2-complete meanwhile.
14. With explicit reviewed config authority, set `plan_catalog_discovery_version: 2` and the
    reviewed `plan_catalog_ownership.excluded_roots`. For the direct route, also set
    `plan_catalog_mode: canonical`. For the optional deferral route, retain the already-validated
    mixed mode and original `plan_catalog_cutover_revision`. Migration itself never activates v2.
15. Validate, commit, and normally publish activation through the owning repository workflow.
    Rerun local validation/list. For portfolio members, land v2 on required member default branches
    before claiming v2 completeness; then refresh the coordination home and require complete
    remote/branch evidence.

## Stop Conditions

Stop on `invalid_ledger`, policy/candidate/revision/hash mismatch, dirty target, active-branch
ownership, semantic or ownership ambiguity, incomplete coverage, malformed/unsafe source,
unresolved conventional root, broken cutover baseline, changed frozen register, transaction/
recovery blocker, incomplete required remote overlays, or any target outside the enumerated
repository scope.

Never update `Last updated` merely because metadata was inserted. Never force a skipped active
plan into migration. Never switch canonical to make validation quieter.

## Evidence And Validation

Retain plan-scoped inventory and reviewed ledger, discovery policy/candidate-set digests, exact
inventory revision, source and target hashes, reviewed-at time, ownership/exclusion decisions,
branch owners/deferrals, dry-run/apply results, preserved header values, mode/discovery activation
commit, any original cutover commit, validation problems, remote completeness, and idempotency.

Validation proves: complete filename-or-metadata candidate enumeration at arbitrary docs depth;
schema-valid metadata plus supported canonical filename; stable IDs and aliases; byte-identical
historical content dates for metadata-only changes; zero writes on excluded, unowned, preview,
changed, dirty, or branch-owned targets; optional mixed baseline integrity; canonical coverage; no
unrelated file changes; public-safe metadata; separate explicit activation; and idempotent reapply.

## Recovery

On optimistic mismatch, discard the stale proposed action, inventory the latest bytes and policy,
and reevaluate semantically. On transaction failure, preserve marker, backup, quarantine, and
recovery artifacts and follow CLI `check`/repair guidance. If activation validation fails, preserve
the migration evidence and keep or restore v1 through a reviewed config-only recovery; never undo
reviewed plan metadata destructively. On incomplete remote evidence, do not claim v2 completeness;
retry via `refresh-portfolio-catalog.md` after required members activate. The prior complete cache
is historical evidence only.
