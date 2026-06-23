Last updated: 2026-06-23T18:09:34Z (UTC)

# Placement Contract

This reference defines the first public placement contract for Codeheart Operating Kit.

## Repository-Owned Areas

- `AGENTS.md`: maintainer bootstrap and task routing.
- `README.md`: repository purpose and public boundary.
- `docs/`: public human-readable kit documentation.
- `docs/repo/`: repository governance, runbooks, and reference docs.
- `components/`: versioned managed content, scaffolds, templates, validators, and component
  metadata.
- `profiles/`: profile presets.
- `schemas/`: machine-readable contracts.
- `src/`: CLI source once implemented.
- `tests/`: validation coverage and fixtures.
- `scripts/`: release, validation, and packaging helpers once implemented.

## Installed Consumer Areas

The Operating Kit may create or manage these consumer paths when the CLI is implemented:

- `.codeheart/kit/`: managed Operating Kit content.
- `.codeheart/kit/README.md`: managed fallback inventory and route repair target.
- `.codeheart/kit.lock.yaml`: installed version, checksum, managed paths, capability status, and
  update-check state.
- `.codeheart/kit.config.yaml`: shared non-secret consumer configuration.
- `.codeheart/user/`: ignored local user layer.
- `docs/repo/`: consumer-owned repository-specific documentation scaffold.
- `docs/repo/state/`: consumer-owned committed state for installed modules and extensions.
- `docs/agent-memory/`: consumer-owned agent memory state scaffold.

## Ownership Modes

- `managed`: synchronized from an Operating Kit release and checked for drift.
- `scaffold`: created when absent, then owned by the consumer.
- `kit-initialized consumer state file`: created or recreated when absent from an Operating Kit
  baseline; the kit owns location, format, and presence behavior, while the consumer owns content
  after creation.
- `template`: available as a starter or example, but installed only when a command explicitly uses
  it.

## Component Target Rules

- Managed documentation component files target `.codeheart/kit/docs/<component>/`.
- `.codeheart/kit/README.md` is the one G1 managed fallback inventory target outside
  `.codeheart/kit/docs/` because root `AGENTS.md` routes unclear kit-structure tasks there.
- Agent-memory state scaffolds target `docs/agent-memory/` and are never overwritten after
  creation.
- Repository documentation starters target `docs/repo/` only as absent-file scaffolds. Reusable
  generic doctrine belongs in the Operating Kit, not in consumer `docs/repo/`.
- Committed module or extension routing state targets `docs/repo/state/<module-or-extension-id>/`
  only when a module or extension has real non-secret repo-owned state to store. The Operating Kit
  defines the placement rule but does not scaffold empty state folders by default.
- Plan-register state files target `docs/repo/plans/plan-register.md` and
  `docs/repo/plans/coordination-sync-pending.md` as kit-initialized consumer state files. Sync may
  create them when absent and must not overwrite them when present.
- Adding these files is additive and does not force migration, movement, rewrite, or archival of
  existing consumer-owned planning or agent-memory content.
- Root `AGENTS.md` receives the Operating Kit managed block from the agent-interface template while
  preserving repository-owned and local-user sections.
- Local user guidance targets `.codeheart/user/` and must stay ignored or local-only.
- G1 does not define or scaffold `docs/workspace/`.

## Consumer-Owned Boundaries

- `docs/repo/`: repository-specific plans, runbooks, references, local commands, validation notes,
  architecture notes, and exceptions to Operating Kit defaults.
- `docs/repo/state/<module-or-extension-id>/`: committed, non-secret routing context for installed
  modules and extensions. This state is not a source of live external truth and does not authorize
  sensitive reads or external changes.
- Product or module docs: local product, package, module, or source-area guidance owned by the
  consumer repository.
- `.codeheart/user/`: personal local preferences and notes.
- `.codeheart/kit.config.yaml`: shared non-secret setup configuration.
- `.codeheart/kit.lock.yaml`: generated installed-state and update-check metadata.

## Placement Rules

- Keep reusable generic operating doctrine in the Operating Kit.
- Keep consumer-specific rules, product details, local commands, credentials, and memory state in
  the consumer repository or workspace.
- Do not make `.codeheart/kit/` a place for user-specific or consumer-authored guidance.
- Do not let sync overwrite consumer-owned docs or memory state.
- Add a consumer-impact classification before changing generated paths, routing, sync behavior, or
  safety policy.
