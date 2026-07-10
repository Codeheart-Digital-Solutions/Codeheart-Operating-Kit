Last updated: 2026-07-09T23:30:00Z (UTC)

# Codeheart Operating Kit

Codeheart Operating Kit is the public Codeheart foundation for agent-first operating standards,
bootstrap instructions, managed documentation, reusable components, validators, and native Codex
capability checks.

This repository contains public-core-safe material only. Do not add private Codeheart business
records, customer or tenant details, secrets, credentials, instance state, or restricted strategy
content.

## Start Here

Give Codex this prompt to set up the latest public Operating Kit:

```text
Set up Codeheart Operating Kit:

https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/latest/download/bootstrap.md
```

## Public Boundary

The Operating Kit owns reusable operating doctrine and release artifacts. Consumer repositories own
their local product guidance, plans, runbooks, memory state, credentials, and environment-specific
configuration.

## Reliable Installed State And Releases

The Go CLI compiles component and profile declarations into one desired-state graph and classifies
installed state before any lifecycle command acts. Mutating commands use one target-filesystem
transaction with dry-run, lock-last commit, post-check, rollback, and recovery markers only when
rollback cannot finish. `init`, `repair`, `sync`, `upgrade`, and `check` have separate
responsibilities; `update-check` is an optional, user-invoked lookup, and only `upgrade --yes` may
change kit version.

Embedded `manifest.yaml` identifies content, compatibility, components, profiles, and graph digest
without circular archive URLs or checksums. Deterministic release packs carry a pack manifest and
payload checksums; an external catalog generated after the packs binds archive and pack-manifest
digests. Install and upgrade verify that catalog-to-binary chain. Lock v1 has a bounded migration
for the two released zero-checksum placeholders; unrelated invalid or future state fails closed.

Supported release platforms remain macOS universal and Windows x64. Catalog verification is
currently HTTPS plus SHA-256 under the unsigned internal/prototype boundary.

## Maintainer Entry Points

- `AGENTS.md`: agent-facing bootstrap and maintainer routing.
- `docs/README.md`: documentation router.
- `bootstrap.md`: first-run public bootstrap path for users without preinstalled Codeheart skills.
- `install.sh`: macOS user-level installer and repair script.
- `install.ps1`: Windows user-level installer and repair script.
- `cmd/codeheart-operating-kit/`: self-contained Go CLI entry point.
- `internal/`: Go CLI implementation, embedded resource handling, command behavior, and release
  helpers.
- `resources.go`: Go embedded Operating Kit resource bundle.
- `components/`: managed Operating Kit component source content.
- `profiles/standard.yaml`: first G1 profile preset.
- `schemas/`: state, declarations, embedded content identity, external release catalog, pack,
  compatibility release-manifest, and consumer-impact contracts.
- `pyproject.toml`: legacy Python package metadata and behavior-oracle entry point retained during
  the Go CLI migration.
- `src/codeheart_operating_kit/`: legacy Python CLI oracle for parity tests during migration.
- `scripts/build-release-assets.py`: deterministic macOS universal and Windows x64 pack builder,
  repeat-build verifier, and external catalog emitter.
- `scripts/validate-*.py`: public-core, Markdown timestamp, JSON schema, and release-manifest
  validators.
- `tests/`: Go/Python parity, CLI, manifest, onboarding, sync/check, update-check, installer, and
  release-asset tests.
- `templates/`: templates for installed consumer surfaces.
- `docs/repo/reference/placement-contract.md`: repository documentation placement contract.
- `docs/repo/reference/consumer-impact-classification.md`: consumer-impact classes for kit
  changes.
- `docs/repo/reference/public-core-hygiene-inventory.md`: public-core extraction inventory.
- `docs/repo/runbooks/change-operating-kit.md`: procedure for changing this repository.
- `docs/repo/runbooks/release-operating-kit.md`: procedure for public releases.
- `docs/repo/runbooks/promote-consumer-change.md`: procedure for promoting consumer-local
  improvements into the kit.
