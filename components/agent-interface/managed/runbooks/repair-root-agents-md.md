Last updated: 2026-06-13T22:55:57Z (UTC)

# Repair Root AGENTS.md

Use this runbook when a consumer root `AGENTS.md` has missing, stale, or malformed Operating Kit
routing.

## Procedure

1. Detect managed block markers.
2. Report missing or malformed markers.
3. Preserve content outside the managed block.
4. Replace only managed-block content when repair is safe.
5. Verify direct routes and update-check trigger.
6. Report unresolved conflicts instead of overwriting local instructions.

## Stop Conditions

Stop before writing when the file has ambiguous ownership boundaries or when repair would delete
repository-owned instructions.
