Last updated: 2026-06-22T20:58:22Z (UTC)

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
4. Preserve the existing stable ID unless the current task explicitly includes an ID migration.
5. For new coordinated-portfolio entries, follow the repository's local ID convention. A
   repository-qualified ID such as `EXAMPLE-AUTOMATION-PR-001` is acceptable when the entry should
   use the same ID locally and in the coordination home.
6. Add or refresh only the compact index fields defined in `../reference/plan-register-format.md`.
7. Copy lifecycle metadata from the canonical document as a snapshot.
8. Add or update relations using the standard relation vocabulary.
9. Resolve the current session ref with the bounded procedure below when the current task needs a
   new creation or material-update session ref.
10. Record creating or material-update session refs when a session ID is available.
11. Omit the session row or record `not recorded` when no session ID is available.
12. Keep detailed status, blockers, decisions, execution evidence, and next actions in the
   canonical documents or execution logs.

If the register and canonical plan conflict, trust the canonical plan and refresh the register.

## Portfolio Coordination Procedure

Use portfolio coordination only when `.codeheart/kit.config.yaml` contains a `portfolio` block.
No `portfolio` block means no configured portfolio coordination.

For `portfolio.role: member`:

1. Complete the local register update first.
2. Read `portfolio.coordination_home_path` and `portfolio.coordination_home_register_path`.
3. For member entries added to the coordination-home register, derive the source namespace from
   `portfolio.member_repository_id`. Normalize it by uppercasing letters, replacing runs of
   non-alphanumeric characters with one hyphen, and trimming leading or trailing hyphens.
4. If the member local ID already begins with `<SOURCE-NAMESPACE>-`, reuse it as the
   coordination-home ID.
5. If the member local ID is bare, build the coordination-home ID as
   `<SOURCE-NAMESPACE>-<source local ID>`. Do not copy a bare member-local ID such as `PR-001`
   into the coordination-home register as the coordination-home ID.
6. Preserve the member repository's source local ID in `Coordination note` as
   `Source local register ID: <ID>`.
7. If the coordination home is locally available and safe to edit, update the coordination-home
   register using the same entry format and the coordination-home ID.
8. If the coordination home is missing, inaccessible, outside the agent's writable scope, or
   otherwise unsafe to edit, record pending sync locally and continue the local planning task.

For `portfolio.role: coordination-home`:

1. Treat `portfolio.coordination_home_register_path` as the local coordination-home register path.
2. Update it for local portfolio-level entries and for explicitly requested member updates.
3. For explicitly requested member updates, derive the source namespace from the member's
   repository ID when available, otherwise from `Owner / repository`.
4. Reuse an already repository-qualified member ID when it begins with the normalized source
   namespace.
5. Build `<SOURCE-NAMESPACE>-<source local ID>` only when the member source local ID is bare.
6. Preserve the member's source local ID in `Coordination note`.
7. Do not scan sibling repositories, GitHub organizations, or remote repositories unless the user
   explicitly asks for that separate discovery or enrollment work.

Do not infer coordination homes from repository names, sibling folder names, GitHub organizations,
or private conventions.

In coordination-home registers, relations between represented entries should use coordination-home
IDs. Use explicit repository/path pointers when a related member plan is not represented by a
coordination-home entry. Keep detailed dependency rationale and execution evidence in canonical
planning documents and execution logs.

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
Source local register ID: <ID>
Target coordination-home ID: <ID>
Intended change: <add | update | complete | supersede | archive | relation-update>
Reason: <why coordination-home sync is needed>
Date: YYYY-MM-DD
Session ref: session <session-id> | not recorded | ambiguous: <reason> | not confidently identified
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

## Session Reference Resolution

Use this self-contained procedure for plan-register session refs. Do not depend on a separate
continuity runbook for this step.

Prefer a direct runtime-provided session ID when the agent environment exposes one. If no direct
runtime value is available, use a bounded read-only metadata scan:

1. Resolve the local Codex state root:

   ```sh
   CODEX_STATE_ROOT="${CODEX_HOME:-$HOME/.codex}"
   ```

2. Look for dated session JSONL files under:

   ```text
   $CODEX_STATE_ROOT/sessions/YYYY/MM/DD/*.jsonl
   ```

3. Inspect only filenames, modification times, and the first JSON record of likely files. The
   first record should be session metadata. Useful metadata fields include:

   - `payload.id`
   - `payload.timestamp`
   - `payload.thread_source`
   - `payload.source`
   - `payload.cwd`

4. Prefer a single candidate whose metadata indicates the main user thread, matches the current
   repository path, and has plausible start time or recent modification time.
5. Exclude helper, tool, and subagent sessions when metadata identifies them, especially
   `thread_source: subagent`.
6. If multiple user-thread candidates remain and confirmation is necessary, search only for a
   distinctive current-turn phrase with filename-only output:

   ```sh
   rg -l '<unique current-turn phrase>' "$CODEX_STATE_ROOT/sessions/YYYY/MM/DD"
   ```

   Use the matching filename only to select the candidate session metadata. Do not print matching
   transcript lines or paste transcript bodies into the register.
7. Stop at the first confident result. Do not browse historical transcript content to make a
   session ref nicer.

Record confidence explicitly:

- use `session <session-id>` when one main user session is confidently identified;
- use `not recorded` when the environment has no usable session metadata or no scan was performed;
- use `ambiguous: <reason>` when more than one plausible user session remains;
- use `not confidently identified` when metadata exists but does not support a confident current
  session match.

## Session Reference Handling

Record session IDs only when they are available in the current environment or confidently resolved
through the bounded metadata-first procedure. Do not invent session IDs, print private transcript
content, or block the task because the ID is unavailable.

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
- Do not move or rewrite unrelated local state files while maintaining the register.
- Do not write to a coordination home when its local instructions, worktree state, permissions, or
  user direction make the change unsafe.
- Do not add private repository topology, customer names, tenant details, credentials, account
  IDs, or local machine details to managed Operating Kit doctrine.
- Use generic repository examples in managed docs and templates.
