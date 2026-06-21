Last updated: 2026-06-21T14:53:02Z (UTC)

# Maintain Plan Register

Use this runbook when creating or materially updating discovery documents, implementation plans,
plan families, major workstreams, or portfolio-relevant planning records.

The plan register is a durable index. Keep canonical details in the canonical planning documents.

## Trigger

Maintain the local register when a discovery or implementation plan is:

- created;
- materially changed in scope, decision state, readiness, implementation path, or review outcome;
- linked as a parent, child, dependency, superseding, or related plan;
- marked `active`, `completed`, `superseded`, or `archived`.

Do not update the register for typos, formatting-only edits, timestamp-only edits, or mechanical
checklist progress that does not change lifecycle, relationships, or implementation path.

## Required References

Read:

- the canonical plan or workstream document being registered;
- `../reference/plan-register-format.md`;
- `.codeheart/kit.config.yaml` when present and portfolio coordination may apply.

## Local Register Procedure

1. Locate the local register at `docs/repo/plans/plan-register.md`.
2. Create the register from the kit baseline when it is absent and the current task includes
   register maintenance.
3. Find an existing entry by stable ID, title, canonical path, or relation.
4. Add or refresh only the compact index fields defined in `../reference/plan-register-format.md`.
5. Copy lifecycle metadata from the canonical document as a snapshot.
6. Add or update relations using the standard relation vocabulary.
7. Record creating or material-update session refs when a session ID is available.
8. Omit the session row or record `not recorded` when no session ID is available.
9. Keep detailed status, blockers, decisions, execution evidence, and next actions in the
   canonical documents or execution logs.

If the register and canonical plan conflict, trust the canonical plan and refresh the register.

## Portfolio Coordination Procedure

Use portfolio coordination only when `.codeheart/kit.config.yaml` contains a `portfolio` block.
No `portfolio` block means no configured portfolio coordination.

For `portfolio.role: member`:

1. Complete the local register update first.
2. Read `portfolio.coordination_home_path` and `portfolio.coordination_home_register_path`.
3. If the coordination home is locally available and safe to edit, update the coordination-home
   register using the same entry format.
4. If the coordination home is missing, inaccessible, outside the agent's writable scope, or
   otherwise unsafe to edit, record pending sync locally and continue the local planning task.

For `portfolio.role: coordination-home`:

1. Treat `portfolio.coordination_home_register_path` as the local coordination-home register path.
2. Update it for local portfolio-level entries and for explicitly requested member updates.
3. Do not scan sibling repositories, GitHub organizations, or remote repositories unless the user
   explicitly asks for that separate discovery or enrollment work.

Do not infer coordination homes from repository names, sibling folder names, GitHub organizations,
or private conventions.

## Pending Sync Fallback

When a member repository has configured portfolio coordination but the coordination-home register
cannot be updated, write a pending item to
`docs/repo/plans/coordination-sync-pending.md`.

Use this entry shape:

```md
## Pending Sync - YYYY-MM-DD - <affected plan ID or title>

Source repository: <member repository ID or repository name>
Target coordination register: <coordination_home_path>/<coordination_home_register_path>
Affected plan entry: <ID, title, and canonical path>
Intended change: <add | update | complete | supersede | archive | relation-update>
Reason: <why coordination-home sync is needed>
Date: YYYY-MM-DD
Session ref: <session ID, not recorded, or unavailable>
Status: pending

Notes:
- <brief note about why the coordination home was unavailable or unsafe to edit>
```

When the pending sync is later applied:

1. Update the coordination-home register.
2. Change `Status: pending` to `Status: completed`.
3. Add a completion date note.
4. Keep the completed item for traceability unless the user explicitly asks to archive or clean up
   old pending-sync entries.

Missing coordination-home access does not fail the local planning task.

## Session Reference Handling

Record session IDs only when they are available in the current environment. Do not invent session
IDs, search private transcripts, or block the task because the ID is unavailable.

Material update reasons include:

- scope changed;
- decision state changed;
- lifecycle status changed;
- parent, child, dependency, supersession, or related-plan links changed;
- readiness state changed;
- implementation path changed;
- review outcome changed.

Use short notes. Do not add session summaries.

## Safety Rules

- Never overwrite existing register or pending-sync content.
- Do not move or rewrite `docs/agent-memory/goal-register.md`.
- Do not write to a coordination home when its local instructions, worktree state, permissions, or
  user direction make the change unsafe.
- Do not add private repository topology, customer names, tenant details, credentials, account
  IDs, or local machine details to managed Operating Kit doctrine.
- Use generic repository examples in managed docs and templates.
