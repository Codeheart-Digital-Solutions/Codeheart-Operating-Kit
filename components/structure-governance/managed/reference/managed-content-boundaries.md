Last updated: 2026-06-13T22:55:57Z (UTC)

# Managed Content Boundaries

Operating Kit installed content uses clear ownership modes.

## Boundaries

- `.codeheart/kit/`: managed Operating Kit content.
- `.codeheart/user/`: ignored local user layer.
- `docs/`: consumer-owned docs, plans, product docs, and memory state.

## Modes

- `managed`: synchronized from a release and checked for drift.
- `scaffold`: created when absent, then consumer-owned.
- `template`: example or starting point, installed only by explicit command/profile behavior.

Sync must not treat consumer docs as managed Operating Kit authority.
