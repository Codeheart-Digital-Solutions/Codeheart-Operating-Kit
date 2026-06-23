Last updated: 2026-06-23T18:09:34Z (UTC)

# Managed Content Boundaries

Operating Kit installed content uses explicit ownership modes so sync can strengthen managed
doctrine without overwriting consumer-owned state.

## Installed Areas

- `.codeheart/kit/`: managed Operating Kit content synchronized from a release.
- `.codeheart/kit/README.md`: managed fallback inventory and route-repair target.
- `.codeheart/kit/docs/`: managed runbooks and references.
- `.codeheart/kit.lock.yaml`: generated installed-state and update-check metadata.
- `.codeheart/kit.config.yaml`: shared non-secret consumer configuration.
- `.codeheart/user/`: ignored local user layer.
- `docs/repo/`: consumer-owned repository-specific documentation scaffold.
- `docs/repo/state/`: consumer-owned committed state for installed modules and extensions.
- `docs/agent-memory/`: consumer-owned agent memory state scaffold.

## Ownership Modes

- `managed`: synchronized from an Operating Kit release and checked for drift.
- `scaffold`: created when absent, then owned by the consumer.
- `template`: starter or example content installed only when a command explicitly uses it.
- `consumer-owned`: durable repository, product, memory, and business content owned by the
  consumer after install.
- `committed module state`: non-secret repository-owned routing state under
  `docs/repo/state/<module-or-extension-id>/`.
- `local-user`: personal local preferences and notes that should stay ignored or local-only.
- `generated`: machine output such as lockfiles, checksums, reports, or release assets.
- `report`: generated or plan-scoped evidence used for review, not managed doctrine.

## Rules

- Do not put user-specific or consumer-authored guidance in `.codeheart/kit/`.
- Do not let sync overwrite consumer-owned docs or memory state.
- Do not promote live memory entries into managed kit docs.
- Keep reusable generic operating doctrine in managed kit docs.
- Keep consumer-specific commands, product details, local exceptions, credentials, and environment
  details in the consumer repository.
- Keep committed module or extension state under `docs/repo/state/<module-or-extension-id>/` only
  when it is non-secret, repo-owned, and useful for routing.
- Treat committed module or extension state as routing context. Validate live external systems
  before sensitive reads, writes, permission changes, or external resource changes.
- Keep local user guidance under `.codeheart/user/` and ignored by source control.
- Treat generated reports as evidence, not as managed source of truth.

## Local Wrappers

Consumer repositories may keep local wrapper docs when existing routes remain useful. A wrapper
should:

- point to the managed kit doc that owns the generic rule;
- keep only local exceptions;
- avoid duplicating long managed doctrine;
- state when local behavior intentionally differs from the kit.

Wrappers are consumer-owned. Sync should not manage or overwrite them.
