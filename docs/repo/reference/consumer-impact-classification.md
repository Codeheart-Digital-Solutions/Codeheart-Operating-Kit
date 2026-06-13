Last updated: 2026-06-13T22:47:44Z (UTC)

# Consumer Impact Classification

Use this reference before changing Codeheart Operating Kit behavior that may affect installed
consumers.

## Impact Classes

`managed-content`

Changes synchronized files under `.codeheart/kit/`. Requires drift-aware validation and release
notes.

`scaffold`

Changes files created only when absent, such as agent-memory state scaffolds. Requires clear
consumer ownership language and tests proving existing files are preserved.

`template`

Changes examples or optional starting points. Requires discoverability updates when template names
or locations change.

`consumer-owned-guidance`

Changes documentation that tells consumers where local rules, product guidance, plans, or memory
state belong. Requires placement-contract review.

`generated-surface`

Changes generated paths, root `AGENTS.md` managed sections, config files, lockfiles, reports, or
`.gitignore` entries. Requires migration notes and validation fixtures.

`release-surface`

Changes installers, release assets, manifests, checksums, release notes, tags, or bootstrap
instructions. Requires release-runbook execution.

`safety-policy`

Changes public-core hygiene, destructive-action rules, secrets handling, or external-action rules.
Requires explicit review and release notes.

## Required Record

For each non-trivial change, record:

- impact class;
- affected paths;
- validation used;
- release-note requirement;
- migration or adoption-note requirement;
- known consumer action.
