Last updated: 2026-06-13T23:16:03Z (UTC)

# Consumer Impact Classification

Use this reference before changing Codeheart Operating Kit behavior that may affect installed
consumers.

Machine-readable consumer-impact records use `schemas/consumer-impact.schema.json`.

## Impact Classes

`instruction-only change`

Changes managed instructions, references, runbooks, templates, or docs without changing generated
paths, validators, safety policy, or required consumer action. Requires discoverability review and
release notes when shipped.

`validator-only change`

Changes validation logic, validator manifests, or validation inputs without changing installed
content placement. Requires validator evidence and release notes when shipped.

`component addition`

Adds a new component, component manifest, or profile-selected component. Requires manifest,
profile, release-manifest, and sync/check validation.

`backwards-compatible scaffold addition`

Adds a file or directory that is created only when absent and then owned by the consumer. Requires
tests or fixtures proving existing consumer files are preserved.

`consumer migration required`

Changes behavior in a way that requires an installed consumer to run a migration, repair, sync, or
manual adoption step. Requires explicit migration notes.

`breaking placement-contract change`

Changes where managed, scaffold, template, generated, consumer-owned, or local-user content lives
in a way that can break existing consumers. Requires explicit approval during planning or release
review.

`security or safety policy change`

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

The impact classes in this reference are the public G1 classes used by component manifests,
release manifests, and consumer-impact records.
