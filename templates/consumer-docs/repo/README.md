Last updated: 2026-07-31T22:23:19Z (UTC)

# Repository Documentation

This folder is for consumer-repository-specific documentation.

Use it for:

- local repository plans;
- local repository runbooks;
- local repository references;
- stable local plan-catalog entry point: `plans/plan-register.md`;
- coordination-home strategy and repository-owned interpretation: `portfolio/`;
- build, test, release, and validation details;
- repository-specific architecture notes;
- local exceptions to Operating Kit defaults.

Reusable generic operating doctrine belongs in the Operating Kit, not in this folder. When a local
rule is broadly reusable, propose an Operating Kit change instead of duplicating generic doctrine
here.

Managed structure governance lives under `.codeheart/kit/docs/structure-governance/`. Use those
managed references and runbooks for generic documentation placement, durable names, managed
content boundaries, placement changes, and index maintenance.

Managed semantic plan and portfolio guidance lives under
`.codeheart/kit/docs/planning-workflows/`. Current plan views are derived from canonical plan
metadata; do not manually append numbered register entries after mixed or canonical cutover.
