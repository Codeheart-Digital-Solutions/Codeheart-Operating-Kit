Last updated: 2026-06-13T22:47:44Z (UTC)

# Placement Contract

This reference defines the first public placement contract for Codeheart Operating Kit.

## Repository-Owned Areas

- `AGENTS.md`: maintainer bootstrap and task routing.
- `README.md`: repository purpose and public boundary.
- `docs/`: public human-readable kit documentation.
- `docs/repo/`: repository governance, runbooks, and reference docs.
- `components/`: versioned managed content, scaffolds, templates, validators, and component
  metadata once implemented.
- `profiles/`: profile presets once implemented.
- `schemas/`: machine-readable contracts once implemented.
- `src/`: CLI source once implemented.
- `tests/`: validation coverage once implemented.
- `scripts/`: release, validation, and packaging helpers once implemented.

## Installed Consumer Areas

The Operating Kit may create or manage these consumer paths when the CLI is implemented:

- `.codeheart/kit/`: managed Operating Kit content.
- `.codeheart/kit.lock.yaml`: installed version, checksum, managed paths, capability status, and
  update-check state.
- `.codeheart/kit.config.yaml`: shared non-secret consumer configuration.
- `.codeheart/user/`: ignored local user layer.
- `docs/repo/`: consumer-owned repository-specific documentation scaffold.
- `docs/agent-memory/`: consumer-owned agent memory state scaffold.

## Ownership Modes

- `managed`: synchronized from an Operating Kit release and checked for drift.
- `scaffold`: created when absent, then owned by the consumer.
- `template`: available as a starter or example, but installed only when a command explicitly uses
  it.

## Placement Rules

- Keep reusable generic operating doctrine in the Operating Kit.
- Keep consumer-specific rules, product details, local commands, credentials, and memory state in
  the consumer repository or workspace.
- Do not make `.codeheart/kit/` a place for user-specific or consumer-authored guidance.
- Do not let sync overwrite consumer-owned docs or memory state.
- Add a consumer-impact classification before changing generated paths, routing, sync behavior, or
  safety policy.
