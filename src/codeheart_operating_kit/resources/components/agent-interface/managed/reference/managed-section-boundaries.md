Last updated: 2026-06-13T22:55:57Z (UTC)

# Managed Section Boundaries

The Operating Kit-managed root `AGENTS.md` block must use stable markers so sync can repair it
without rewriting local instructions.

## Required Markers

```md
<!-- BEGIN CODEHEART OPERATING KIT MANAGED BLOCK -->
...
<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->
```

Repository-owned instructions start after the managed block. Local user guidance may be linked or
included below repository-owned instructions when the consumer chooses.

## Repair Rules

- Sync may replace content inside the managed markers.
- Sync must preserve content outside the managed markers.
- Repair must not delete repository-owned instructions.
- Repair must report malformed or missing boundaries before writing.
