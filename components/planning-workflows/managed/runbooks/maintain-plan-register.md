Last updated: 2026-08-09T15:24:37Z (UTC)

# Maintain Plan Register

Audience: agent-facing

Intent:
Use the stable register path correctly across legacy, mixed, and canonical repositories without
reintroducing manual numbering conflicts or duplicating canonical plan authority.

Success:
Legacy repositories receive only authorized compatible maintenance; mixed repositories keep an
actually used frozen register byte-identical and expose deferred rows through the derived view;
canonical repositories use the derived repository-wide view; migration never edits the register;
v1 coordination evidence remains preserved without creating new pending-sync work.

Agent judgment boundary:
Inspect mode and route to the right behavior. Do not infer a mode change, allocate global numbers,
rewrite frozen evidence, invent metadata, or write to another repository.

Stop boundary:
Stop on invalid config, broken mixed baseline, ambiguous legacy insertion, overlapping dirty
register changes, migration semantics, or any request that would modify a frozen register.

## Source, Inputs, Preconditions, And Lane

Source of truth: `.codeheart/kit.config.yaml`, the canonical plan, Git baseline evidence,
`../reference/plan-register-format.md`, `../reference/plan-catalog-format.md`, and current CLI
validation.

Inputs: repository path, requested planning change, current catalog mode, canonical plan path, and
explicit write authority when legacy maintenance or migration preparation is required.

Preconditions: initialized compatible Kit; readable config; understood dirty overlap; regular
contained register when present; no unresolved transaction/recovery state.

Execution lane: local plan CLI and managed planning docs. Legacy Markdown edit is allowed only in
legacy mode. Mixed/canonical current views are read-only CLI output.

Approval boundary: listing and validation are local reads. Creating a missing scaffold, updating a
legacy entry, or changing mode requires user-approved repository writes. Plan-catalog migration,
clearance, deferral, and activation never authorize a register edit or generated register row.
This runbook does not authorize coordination-home/member writes, commits, pushes, PRs, merges,
releases, force-push, deletion, or destructive Git.

## Procedure

1. Read catalog mode and discovery version from config, then run:

   ```sh
   codeheart-operating-kit plans validate <repository>
   ```

2. Select behavior from the reported mode.
3. If the task concerns discovery-v2 adoption, do not infer the expanded plan set from the active
   v1 view. Route to `migrate-plan-catalog.md` for prospective inventory and explicit activation.

### Legacy Mode

1. If the register is absent and the task includes creating its entry point, use the Kit scaffold
   without overwriting any existing path.
2. If this repository still intentionally uses legacy entries, locate the entry by canonical path,
   existing ID, and title. Preserve its ID. Add or refresh only the compact legacy fields the
   existing file already uses.
3. Preserve unrelated dirty content and reject duplicate/ambiguous IDs or canonical paths. Do not
   renumber parallel entries during ordinary authoring.
4. Only when filename-only records will be intentionally deferred through mixed mode, validate and
   record the existing regular register and legacy plans at the one exact legacy cutover. Do not
   add a deferral notice, generated row, or migration-only register edit; if the current register
   is not already adequate frozen authority, stop for a separate explicit legacy-maintenance
   decision before taking migration evidence.
5. When all candidates can receive metadata, do not manufacture a mixed cutover. Route the direct
   legacy-to-canonical migration to `migrate-plan-catalog.md`.

### Mixed Mode

1. Do not edit `plan-register.md`; mixed validation requires it to equal the cutover revision.
2. Author every new or materially rewritten formal plan with its supported filename and canonical
   metadata, wherever it lives beneath an owned tracked `docs` tree.
3. Use `plans list` for the current view. A schema-v3 `deferred-active-owner` remains visible once
   with `coverage_disposition: mixed-grandfathered` and
   `incremental_migration_required: true`. It may fill only its exact pre-cutover compatibility
   gap; it is never a generated register row.
4. If a legacy plan needs metadata, use the reviewed migration ledger; do not append a register row.
5. When the owner ref is integrated or its target path/bytes change, route to
   `migrate-plan-catalog.md` for `reinventory-and-incremental-migrate`. Owner-supplied canonical
   metadata at merge or a new reviewed ledger against merged bytes may close the follow-up.
6. Treat `legacy_register_modified_after_cutover`, stale config-v2 evidence, or other baseline
   errors as blockers and restore only through reviewed recovery from the exact baseline.

### Canonical Mode

1. Do not manually append or regenerate register entries.
2. Use `plans validate` and `plans list --format text|json`. With discovery v2 active, these views
   cover every owned tracked Markdown documentation tree at arbitrary depth, not only
   `docs/repo/plans/`.
3. Update current facts only in the canonical document metadata/header/content under the relevant
   authoring or execution runbook.
4. Retain the old register as historical evidence unless an explicit archival plan says otherwise.

### Portfolio Coordination

1. For portfolio v2, do not mirror member plans into a coordination Markdown register. Refresh the
   home catalog through `refresh-portfolio-catalog.md` and maintain strategy separately.
2. Preserve existing v1 register and `coordination-sync-pending.md` files. Apply an already-recorded
   v1 pending item only under that repository's explicit compatibility authority.
3. Do not create a new pending-sync file. Missing global visibility normally means the plan has not
   been pushed; publication follows the plan-checkpoint authority in the planning/execution
   runbooks.

## Stop Conditions

Stop on malformed metadata/config, non-regular or symlinked authority, duplicate semantic IDs,
ambiguous legacy matches, dirty overlap, missing baseline commit, changed frozen register,
stale/missing ledger-v3 evidence, incomplete mixed coverage, pending incremental follow-up when
canonical readiness is required, incomplete canonical coverage, or another repository boundary.

Do not use a register edit to hide validation errors or make an unpushed plan globally visible.

## Evidence And Validation

Record mode, validation result, canonical path/ID, changed paths, baseline commit when applicable,
problems/blockers, and confirmation that unrelated/frozen files were preserved.

Validation succeeds when legacy edits are unique and canonical-document-backed, mixed register
bytes equal the baseline when mixed is used, every deferred row remains visible with its exact
follow-up, `mixed_coverage_complete` and `canonical_ready` retain their distinct meanings,
canonical discovery-v2 views contain every owned candidate and disclose excluded/ambiguous/unowned
observations, and no migration-time register edit, portfolio member write, or new pending-sync file
was introduced.

## Recovery

If a mixed register changed, stop new authoring and compare it to the exact cutover commit. Preserve
concurrent bytes and use a reviewed recovery change; never blindly overwrite. For migration
blockers use `migrate-plan-catalog.md`. For incomplete coordination evidence use
`refresh-portfolio-catalog.md` and disclose freshness.
