Last updated: 2026-07-31T22:23:19Z (UTC)

# Plan Register Format

Use this reference for `docs/repo/plans/plan-register.md`. The path remains a stable entry point so
old links and low-context agents keep working while plan facts move to canonical document metadata
and generated views.

## Current Authority

Canonical discovery, implementation, and qualifying family documents own current plan facts. Use:

```sh
codeheart-operating-kit plans validate .
codeheart-operating-kit plans list --format text .
codeheart-operating-kit plans list --format json .
```

The CLI discovers formal files, validates their lifecycle and bounded metadata, reconciles allowed
legacy evidence, and generates deterministic repository-scoped output. It does not write a
committed catalog or renumber plans.

The stable `plan-register.md` may contain:

- the compact command entry point;
- stable repository-specific orientation or exceptions;
- a frozen historical register retained as migration evidence; and
- a clear frozen-authority notice in repositories that adopted mixed/canonical mode.

Do not copy generated rows back into this file.

## Behavior By Catalog Mode

`legacy`
: Existing numbered register entries remain authoring/index authority under the legacy contract.
  New repositories may instead keep only the compact entry point. Use legacy maintenance only
  until the repository deliberately captures its mixed-mode baseline.

`mixed`
: The exact pre-cutover register bytes are frozen. Current new plans require metadata. Existing
  legacy plans are grandfathered only when both their path and register evidence exist in the full
  `plan_catalog_cutover_revision`. Do not append, reorder, normalize, or add a notice after the
  baseline; prepare any frozen notice before committing the baseline.

`canonical`
: Every formal plan has canonical metadata. The old register body is historical evidence only.
  Keep it stable unless a later explicit archival plan changes the retention policy.

See `plan-catalog-format.md` for mode and baseline requirements.

## Compact Entry-Point Shape

Fresh repositories use a short document like:

````md
# Plan Catalog Entry Point

## Current View

```sh
codeheart-operating-kit plans validate .
codeheart-operating-kit plans list --format text .
```

## Repository Notes

Add only stable repository-specific orientation or exceptions here.
````

This shape is repository-owned after creation. Sync creates it only when absent and never replaces
an existing register.

## Legacy Entry Evidence

Older registers may contain sequential IDs such as `PR-001`, types, titles, purposes, canonical
paths, lifecycle snapshots, relations, and session/coordination notes. Treat these fields as
semantic-review evidence:

- preserve a linked historical ID in canonical metadata `legacy_aliases`;
- split combined discovery/implementation rows into the actual separate records;
- prefer the current canonical document for title, lifecycle, purpose, and relationships;
- record ambiguous or unpaired evidence in the migration ledger; and
- never allocate a new sequential ID merely to add metadata.

Parallel legacy branches may contain duplicate numbers. Repository/path, content, Git revision,
and semantic review disambiguate them. The generated view uses canonical semantic IDs rather than
attempting a global renumber.

## Frozen-Authority Notice

When an existing legacy register will become the mixed baseline, add a notice before committing
the last legacy-mode revision. Example:

```md
> Frozen legacy evidence: current plan authority is derived from canonical plan metadata and
> `codeheart-operating-kit plans list`. Do not append or renumber entries after mixed cutover.
```

Commit the notice, regular register, legacy plans, and valid legacy-mode config together. Record
that commit as `plan_catalog_cutover_revision` only in the later mixed-mode config change. Mixed
validation requires current register bytes to remain identical to that committed baseline.

## Generated Views And Stable Links

Text output is for human orientation. JSON output owns exact structured fields. Both are derived
at request time and may report problems. A nonzero result must not be presented as complete.

Keep durable links pointed to canonical plan paths or semantic IDs. Keep old links to
`plan-register.md` valid as an entry point. Do not create a second generated Markdown register on
every branch; separate plan files and derived views avoid the old shared insertion point.

## Coordination-Home Behavior

Portfolio coordination does not maintain a second global Markdown register. The coordination home
refreshes pushed default and unmerged-branch observations into its ignored local cache and keeps
strategic interpretation in `docs/repo/portfolio/strategic-overlay.yaml`.

Existing portfolio-v1 coordination registers and `coordination-sync-pending.md` remain readable
compatibility evidence. Do not create new pending-sync files or append v2 source observations to a
Markdown register. Use `portfolio-coordination-format.md` and the refresh runbook.

## Session References And Progress

Session references already present in a legacy register remain evidence. New canonical metadata
does not include session IDs. Keep detailed status, decision history, review outcomes, progress,
and validation in the canonical plan or execution log, not in this entry point.

## Anti-Patterns

Do not use the register as:

- a manual current catalog in mixed/canonical mode;
- a sequential ID allocator;
- a branch-visible activity feed;
- a task backlog, sprint board, or per-epic progress table;
- a duplicate plan body or transcript index;
- a generated cache committed after every scan; or
- a place for credentials, private topology, customer data, or local-machine paths.
