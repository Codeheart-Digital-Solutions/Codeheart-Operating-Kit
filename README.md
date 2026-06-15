Last updated: 2026-06-15T10:16:29Z (UTC)

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

## Maintainer Entry Points

- `AGENTS.md`: agent-facing bootstrap and maintainer routing.
- `docs/README.md`: documentation router.
- `bootstrap.md`: first-run public bootstrap path for users without preinstalled Codeheart skills.
- `install.sh`: macOS user-level installer and repair script.
- `install.ps1`: Windows user-level installer and repair script.
- `components/`: managed Operating Kit component source content.
- `profiles/standard.yaml`: first G1 profile preset.
- `schemas/`: lockfile, config, release-manifest, and consumer-impact contracts.
- `pyproject.toml`: Python package metadata and `codeheart-operating-kit` console entry point.
- `scripts/build-release-assets.py`: macOS and Windows release-asset builder with checksums.
- `scripts/validate-*.py`: public-core, Markdown timestamp, JSON schema, and release-manifest
  validators.
- `src/codeheart_operating_kit/`: CLI source for onboard, inspect, init, sync, check, and
  update-check.
- `tests/`: CLI, manifest, onboarding, sync/check, update-check, and capability tests.
- `templates/`: templates for installed consumer surfaces.
- `docs/repo/reference/placement-contract.md`: repository documentation placement contract.
- `docs/repo/reference/consumer-impact-classification.md`: consumer-impact classes for kit
  changes.
- `docs/repo/reference/public-core-hygiene-inventory.md`: public-core extraction inventory.
- `docs/repo/runbooks/change-operating-kit.md`: procedure for changing this repository.
- `docs/repo/runbooks/release-operating-kit.md`: procedure for public releases.
- `docs/repo/runbooks/promote-consumer-change.md`: procedure for promoting consumer-local
  improvements into the kit.
