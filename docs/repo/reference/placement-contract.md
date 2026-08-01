Last updated: 2026-07-31T22:23:19Z (UTC)

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
- `.codeheart/local/`: ignored local machine/runtime state.
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
- `local-machine`: ignored generated runtime/tooling state that is specific to one checkout and
  can be recreated.

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
- The stable plan-catalog entry point targets `docs/repo/plans/plan-register.md` as an absent-file
  scaffold. Sync may create it when absent and must not overwrite it when present. Existing
  `coordination-sync-pending.md` files are preserved portfolio-v1 compatibility evidence, but
  fresh installs do not scaffold them.
- Coordination-home repository guidance targets `docs/repo/portfolio/README.md` and
  `docs/repo/portfolio/strategic-overlay.yaml` as absent-file scaffolds. They are repo-owned after
  creation; configuration may replace only the exact untouched generic overlay placeholder with
  the approved home ID.
- Adding these files is additive and does not force migration, movement, rewrite, or archival of
  existing consumer-owned planning or agent-memory content.
- Root `AGENTS.md` receives the Operating Kit managed block from the agent-interface template while
  preserving repository-owned and local-user sections.
- Local user guidance targets `.codeheart/user/` and must stay ignored or local-only.
- Local machine/runtime state targets `.codeheart/local/` and must stay ignored or local-only.
  Init and sync may add the ignore rule without creating the directory by default.
- Portfolio local sources, bare mirrors, and complete factual cache target
  `.codeheart/local/portfolio/`. They are rebuildable local-machine state and never plan or
  strategic authority.
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
- `.codeheart/local/`: generated local runtime/tooling state, such as repo-local virtual
  environments, caches, temporary files, generated shims, package artifacts, or generated install
  metadata. It is machine-local, ignored, and recreatable.
- `docs/repo/plans/`: canonical repository plans, plan-scoped evidence, and the stable catalog
  entry point. Current mixed/canonical views are derived rather than committed.
- `docs/repo/portfolio/`: coordination-home strategic interpretation and local coordination docs;
  source-derived factual cache does not belong here.
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
