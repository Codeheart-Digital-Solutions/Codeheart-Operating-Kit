Last updated: 2026-07-31T22:23:19Z (UTC)

# Migrate Plan Catalog

Audience: agent-facing

Intent:
Inventory every formal plan, semantically review it, apply metadata with optimistic source checks,
and move one repository from legacy through mixed to canonical without freezing unrelated work or
rewriting historical content dates.

Success:
Every in-scope formal record has reviewed canonical metadata, active-branch ownership is
reconciled, migration is idempotent, and canonical local plus requested remote-overlay validation
passes while legacy evidence remains preserved.

Agent judgment boundary:
Interpret plan meaning from content and repository evidence; omit unsupported optional metadata.
Do not convert register rows mechanically, invent classifications, edit another branch's owned
plan, modify another repository, or bypass hashes, blockers, dry-run review, and approval.

Stop boundary:
Stop on ambiguous semantics, invalid inventory/ledger, changed source revision or bytes, dirty
overlap, active-branch ownership without coordination, broken mixed baseline, incomplete remote
scan, recovery-required transaction state, or missing `--yes` approval.

## Source, Inputs, Preconditions, And Lane

Source of truth, in order: current repository plans and sibling evidence; Git revision/ref/history;
the frozen legacy register; current user/branch-owner decisions; `../reference/plan-catalog-format.md`;
and the migration schemas and CLI output.

Inputs: repository path; explicit inventory artifact path outside formal plan authority; reviewed
ledger path; repository ID; target modes; branch ownership/deferral decisions; and approval for
each config or migration write.

Preconditions: initialized Kit with compatible catalog CLI; clean or explicitly understood target
paths; no unresolved transaction marker; committed legacy-mode baseline containing the regular
register before mixed adoption; and remote access when overlay validation is requested.

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

1. Inspect current mode and run `codeheart-operating-kit plans validate <repository>`.
2. In legacy mode, place any frozen-authority notice in `plan-register.md` before the baseline is
   captured. Commit the regular register, current legacy plans, and valid legacy-mode config. Save
   that exact full commit as the future `plan_catalog_cutover_revision`. Do not append to the
   register after this point.
3. Create the explicit artifact directory, then inventory:

   ```sh
   codeheart-operating-kit plans inventory \
     --remote-overlays \
     --output <inventory.json> \
     <repository>
   ```

   Omit `--remote-overlays` only when portfolio overlay evidence is not configured or not in the
   approved scope, and record that limitation.
4. Review every formal discovery, implementation, and qualifying family record semantically.
   Reconcile sibling documents, legacy aliases, lifecycle, purpose, family, relations, taxonomy,
   Git source revision/hash, dirty overlap, and active remote branch ownership.
5. Build a ledger matching `schemas/plan-migration-ledger.schema.json`. Record confidence,
   evidence, ambiguity, conflicts, deferral, and branch owner. Omit unsupported optional
   classification; do not guess.
6. Validate the ledger and dry-run:

   ```sh
   codeheart-operating-kit plans migrate \
     --ledger <reviewed-ledger.yaml> \
     --dry-run \
     <repository>
   ```

7. Review every action, skip, blocker, exact source hash, preserved `Last updated`, and unrelated-
   file exclusion. A material skip is not success.
8. With approval, set shared config to `mixed`, require the stable
   `portfolio.member_repository_id`, and set `plan_catalog_cutover_revision` to the exact committed
   legacy baseline. Validate mixed mode before applying metadata.
9. Apply the exact reviewed ledger:

   ```sh
   codeheart-operating-kit plans migrate \
     --ledger <reviewed-ledger.yaml> \
     --yes \
     <repository>
   ```

10. Reinventory. Files whose revision/hash changed or are owned by active branches remain skipped;
    ask their branch owners to add reviewed metadata in their current work or reevaluate the newest
    bytes after ownership stabilizes. No repository-wide freeze is required.
11. Reconcile skipped/new records with updated ledger evidence and reapply through the same dry-
    run and approval path.
12. Run `plans validate --remote-overlays`. Enter `canonical` only after complete current default
    and accessible remote-overlay coverage passes. Then run local/remote validation and
    `plans list --format json` again.
13. Reapply the final ledger and record zero changes or explicit already-applied outcomes.

## Stop Conditions

Stop on `invalid_ledger`, source revision/hash mismatch, dirty target, active-branch ownership,
semantic ambiguity, incomplete coverage, malformed/unsafe source, broken cutover baseline, changed
frozen register, transaction/recovery blocker, incomplete remote overlays, or any target outside
the enumerated repository scope.

Never update `Last updated` merely because metadata was inserted. Never force a skipped active
plan into migration. Never switch canonical to make validation quieter.

## Evidence And Validation

Retain plan-scoped inventory and reviewed ledger, exact inventory revision, source hashes,
reviewed-at time, reviewer decisions, branch owners/deferrals, dry-run/apply results, preserved
header values, mode/cutover commits, validation problems, remote completeness, and idempotency.

Validation proves: complete formal enumeration; schema-valid metadata; stable IDs and aliases;
byte-identical historical content dates for metadata-only changes; zero writes on changed/dirty/
owned targets; mixed baseline integrity; canonical coverage; no unrelated file changes; public-safe
metadata; and idempotent reapply.

## Recovery

On optimistic mismatch, discard the stale proposed action, inventory the latest bytes, and
reevaluate semantically. On transaction failure, preserve marker, backup, quarantine, and recovery
artifacts and follow CLI `check`/repair guidance. On incomplete remote evidence, remain mixed and
retry via `refresh-portfolio-catalog.md`; the prior complete cache is historical evidence only.
