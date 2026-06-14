Last updated: 2026-06-13T22:55:57Z (UTC)

# Documentation Structure

Use ownership first, then domain and artifact kind.

## Common Type Folders

- `plans/`: discovery docs, implementation docs, migration plans, and long-form change plans.
- `runbooks/`: operational procedures.
- `reference/`: stable contracts, naming rules, lifecycle definitions, schemas, and explanatory
  source-of-truth docs.
- `archive/`: superseded or historical docs retained for traceability.

## Consumer Repository Docs

Consumer `docs/repo/` is for local repository plans, runbooks, references, build/test/release
details, repository-specific architecture notes, and local exceptions to Operating Kit defaults.

Generic operating doctrine belongs in the Operating Kit. Agents should flag reusable generic
guidance proposed under consumer `docs/repo/` and recommend changing the Operating Kit instead.

## Index Rules

Update the nearest README and parent index when discoverability changes, including added, moved,
renamed, removed, or archived docs; changed runbook paths; changed command paths; changed validator
inputs; or changed repo entry points.
