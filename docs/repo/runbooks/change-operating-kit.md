Last updated: 2026-07-09T23:30:00Z (UTC)

# Change Operating Kit

Use this runbook before changing Codeheart Operating Kit source, docs, schemas, templates,
validators, installers, release assets, or CLI behavior.

## Procedure

1. Read `AGENTS.md`.
2. Read `README.md`.
3. Read `docs/repo/reference/placement-contract.md`.
4. Read `docs/repo/reference/consumer-impact-classification.md`.
5. Classify the change by consumer impact.
6. Check whether the change belongs in the Operating Kit or in a consumer repository.
7. Keep public-core material only.
8. Update the nearest README when discoverability changes.
9. Update tests, schemas, manifests, or fixtures when behavior changes.
10. Record release-note or migration-note needs for consumer-affecting changes.
11. Run the smallest validation set that proves the changed surface.
12. Summarize validation and residual risk in the PR.

For state, lifecycle, installer, or release changes, also run the matching gates:

- schema and migration tests for declaration, config, lock, content, catalog, or pack contracts;
- transaction failure tests for planning, staging, commit, post-check, rollback, and recovery;
- build-twice byte comparison plus catalog-to-binary verification for release-pack changes;
- isolated installer and upgrade success/failure paths for each affected platform;
- consumer materialization and routing checks when managed guidance changes.

Do not use the source repository's ignored `.codeheart/kit/` tree as producer authority. Source
components, profiles, templates, schemas, Go packages, and maintainer runbooks are canonical.

## Stop Conditions

Stop before editing when the change would expose private content, change release authority, alter
consumer ownership boundaries, or require a new public repository setting that is not already
approved.
