Last updated: 2026-06-13T22:47:44Z (UTC)

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

## Stop Conditions

Stop before editing when the change would expose private content, change release authority, alter
consumer ownership boundaries, or require a new public repository setting that is not already
approved.
