Last updated: 2026-06-13T22:55:57Z (UTC)

# Agent Memory Entry Format

Use grouped, labeled list entries instead of Markdown tables so raw Markdown stays readable.

## Status Values

- `active`: currently moving.
- `blocked`: cannot move until a named blocker is resolved.
- `paused`: intentionally not moving now.
- `resolved`: workstream outcome is complete but still recent.
- `archived`: retained for history.
- `unknown`: needs triage.

## Relation Labels

- `continues`: normal follow-up on the same workstream.
- `spawned-from`: smaller thread discovered inside larger work.
- `blocks`: this workstream prevents progress elsewhere.
- `unblocks`: this workstream resolved a blocker.
- `supersedes`: newer status replaces older status.
- `related`: useful context, but not a dependency.
- `archived`: retained for history, not active.

## Goal Entry Shape

```md
## Goal: <goal name>

Status: active | blocked | paused | resolved | archived | unknown
Last reviewed: YYYY-MM-DDTHH:MM:SSZ (UTC)
Current focus: <one sentence>
Canonical docs:
- <relative path>
Related workstreams:
- <relation>: <workstream id or title>
Ordering notes:
- <timestamp>: <priority or ordering decision>
Current state:
- <short status summary>
Next action:
- <specific next action>
```

## Session Entry Shape

```md
#### Session: YYYY-MM-DDTHH:MM:SSZ (UTC) - <short title>

Session ID: `<session-id>`
Session started: YYYY-MM-DDTHH:MM:SSZ (UTC)
Last observed activity: YYYY-MM-DDTHH:MM:SSZ (UTC)
Relation: continues | spawned-from | blocks | unblocks | supersedes | related | archived
Status: active | blocked | paused | resolved | archived | unknown

Summary:
- <outcome summary>

Handoff:
- <next action>

References:
- `<relative path>`
```
