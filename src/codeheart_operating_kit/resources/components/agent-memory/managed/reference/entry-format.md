Last updated: 2026-06-17T06:40:50Z (UTC)

# Agent Memory Entry Format

This reference owns reusable entry shapes for consumer-owned agent memory files. Live memory files
should start with compact overviews and actual entries, not long templates.

Use grouped, labeled list entries instead of Markdown tables so raw Markdown stays readable and
agents can parse fields reliably.

## Overview Order

Overview sections mirror the order of the detailed entries below them.

- `goal-register.md`: `Goal Overview` follows the same order as `Goals`.
- `session-ledger.md`: `Session Overview` follows the same order as `Active Ledger`.
- `untriaged-sessions.md`: `Inbox Overview` follows the same order as `Inbox`.

Do not maintain a second sort order in the overview. If a maintenance pass reorders detailed
entries, update the overview in the same change.

## Status Values

- `active`: currently moving.
- `blocked`: cannot move until a named blocker is resolved.
- `paused`: intentionally not moving now.
- `resolved`: workstream outcome is complete but still recent.
- `archived`: retained for history outside the active register or ledger.
- `unknown`: needs triage.

## Relation Labels

- `continues`: normal follow-up on the same workstream.
- `spawned-from`: smaller thread discovered inside larger work.
- `blocks`: this workstream prevents progress elsewhere.
- `unblocks`: this workstream resolved a blocker.
- `supersedes`: newer status replaces older status.
- `related`: useful context, but not a dependency.
- `archived`: retained for history, not active.

## Goal Entry

```md
## Goal: <goal name>

Status: active | blocked | paused | resolved | archived | unknown
Last reviewed: YYYY-MM-DDTHH:MM:SSZ (UTC)
Current focus: <one sentence>
Canonical docs:
- <repo-relative path>
Related workstreams:
- <relation>: <workstream id or title>
Ordering notes:
- <timestamp>: <priority or ordering decision, if any>
Current state:
- <short status summary>
Next action:
- <specific next action>
```

## Goal Workstream Entry

```md
### Workstream: <workstream name>

Status: active | blocked | paused | resolved | archived | unknown
Last reviewed: YYYY-MM-DDTHH:MM:SSZ (UTC)
Parent goal: <goal name>
Canonical docs:
- <repo-relative path>
Recent sessions:
- Session: <short title>
  - Started: YYYY-MM-DDTHH:MM:SSZ (UTC)
  - Last observed activity: YYYY-MM-DDTHH:MM:SSZ (UTC)
  - Session ID: <session ID when available>
Related workstreams:
- <relation>: <workstream id or title>
Ordering notes:
- <timestamp>: <priority or ordering decision, if any>
Current state:
- <short status summary>
Next action:
- <specific next action>
Archive references:
- <repo-relative archive path or none>
```

## Session Ledger Entry

```md
### Workstream: <workstream name>

Status: active | blocked | paused | resolved | archived | unknown
Last reviewed: YYYY-MM-DDTHH:MM:SSZ (UTC)
Ordering notes:
- <timestamp>: <priority or ordering decision, if any>

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
- `<repo-relative path>`
```

## Untriaged Entry

```md
### Session: YYYY-MM-DDTHH:MM:SSZ (UTC) - <short title>

Session ID: `<session-id>`
Session started: YYYY-MM-DDTHH:MM:SSZ (UTC)
Last observed activity: YYYY-MM-DDTHH:MM:SSZ (UTC)

Why untriaged:
- <reason>

Possible workstreams:
- <candidate>

Suggested next step:
- <triage action>
```

## Ordering Notes

Use ordering to help the user resume thinking, not as a mechanical status sort.

Ask the user when ordering would imply a meaningful priority decision and the right order is
unclear. Record explicit user ordering or priority decisions as `Ordering notes:` in the relevant
entry.
