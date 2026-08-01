Last updated: 2026-07-31T22:23:19Z (UTC)

# Plan Catalog Entry Point

This stable repository-owned file points agents and humans to the current derived plan view.
Canonical discovery, implementation, and family documents own plan facts; this file is not a
second plan database.

## Current View

From the repository root, run:

```sh
codeheart-operating-kit plans validate .
codeheart-operating-kit plans list --format text .
```

Use `plans list --format json .` when structured output is required. A nonzero result means the
reported problems must be resolved or disclosed.

## Authoring And Migration

- New formal plans follow
  `.codeheart/kit/docs/planning-workflows/reference/plan-catalog-format.md`.
- Repository mode and legacy compatibility follow
  `.codeheart/kit/docs/planning-workflows/runbooks/maintain-plan-register.md`.
- Existing repositories adopt canonical metadata through
  `.codeheart/kit/docs/planning-workflows/runbooks/migrate-plan-catalog.md`.

Do not manually append numbered entries after this repository enters mixed or canonical mode.

## Repository Notes

Add only stable repository-specific orientation or exceptions here. Do not copy generated plan
rows, volatile progress, private topology, credentials, or plan-body content into this entry point.
