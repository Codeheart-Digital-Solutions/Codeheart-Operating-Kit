Last updated: 2026-06-21T15:17:48Z (UTC)

# Session Ledger Maintenance

Use this runbook when the user asks to update agent memory, recover where prior work left off, or
reconcile recent agent sessions into curated repository memory.

This runbook maintains consumer-owned memory files, usually:

- `docs/agent-memory/goal-register.md`
- `docs/agent-memory/session-ledger.md`
- `docs/agent-memory/untriaged-sessions.md`
- `docs/agent-memory/archive/`

Raw local agent state is supporting evidence. Curated files under the consumer memory folder are
the durable project memory.

## Scope

In scope:

- Resolve the local agent state root using user-agnostic patterns.
- Inspect local session metadata and session logs read-only.
- Preserve session IDs when available.
- Classify sessions into goals, workstreams, active ledger entries, archive entries, or untriaged
  entries.
- Update freshness fields and coverage windows.
- Compare coverage timestamps against local session activity.
- Maintain grouped list entries and curated ordering.
- Verify current repository state before writing conclusions.

Out of scope:

- Importing raw transcript bodies.
- Editing raw local agent state.
- Creating automation scripts.
- Replacing canonical product, business, research, discovery, implementation, runbook, reference,
  source-control, issue, or release records.

## Formal Plan Boundary

When a session creates or materially changes a formal discovery document, implementation plan, plan
family, or major workstream, record formal lifecycle, canonical plan docs, blockers that belong to
the plan, and plan relationships through `docs/repo/plans/plan-register.md` and the planning
workflow register-maintenance runbook.

Keep `docs/agent-memory/goal-register.md` available for informal, pre-plan, or transitional
continuity that has not yet become a formal plan. Do not duplicate full formal plan status,
per-epic progress, plan relationships, or canonical execution evidence in `goal-register.md`.

## Preflight

1. Confirm the current repository and inspect worktree status.
2. Protect unrelated user work. Do not overwrite, revert, or clean unrelated changes.
3. Read the consumer memory README when present.
4. Read the managed entry-format reference.
5. Read `goal-register.md`, `session-ledger.md`, and `untriaged-sessions.md` when present.
6. Resolve the local agent state root with a user-agnostic pattern:

   ```sh
   CODEX_STATE_ROOT="${CODEX_HOME:-$HOME/.codex}"
   ```

7. If the resolved root does not exist, stop the raw-session scan and update curated memory only
   from repository evidence and user-provided session IDs.

## Freshness Audit

Before deciding ledgers are current, compare their `Coverage window:` timestamps with local
session activity.

1. Read the newest `Coverage window:` timestamp from memory files.
2. Treat date-only coverage as too coarse for same-day continuation checks.
3. Look for session files modified after the coverage timestamp.
4. If a session index is stale or missing recent activity, inspect the dated session folder for
   the relevant day or bounded recent period.
5. Filter candidates to the current repository by inspecting metadata, turn context, or early log
   records for the current repo path or workspace root.
6. If relevant sessions exist after the coverage timestamp, curated memory is incomplete for that
   interval. Summarize and classify those sessions, or add them to `untriaged-sessions.md`.

Prefer timestamp checks over date-only checks.

## Read-Only Raw Session Inspection

Inspect raw state read-only. Do not edit files under the local agent state root.

Use metadata first. Use dated logs only when metadata and curated memory are not enough.

When inspecting session logs:

- summarize outcomes and handoffs;
- preserve session IDs;
- do not copy raw transcript bodies into repository memory;
- do not commit concrete machine-specific absolute paths.

## Classification Rules

Classify each relevant session into one primary destination.

### Goal Register

Update `goal-register.md` when a session changes:

- top-level informal or pre-plan goal status;
- active informal or transitional workstream status;
- current focus;
- informal or pre-plan blocker state;
- next action;
- lightweight references for resumption;
- relationship between informal or transitional workstreams.

Keep goal entries short. Link canonical docs rather than duplicating them. Keep `Goal Overview` in
the same order as the detailed goals.

When the workstream has a formal discovery or implementation plan, update
`docs/repo/plans/plan-register.md` for formal lifecycle state and relationships. Keep only
resumption-oriented memory here.

### Session Ledger

Add or update `session-ledger.md` when a session is:

- active;
- blocked;
- unresolved;
- recent and resolved;
- related to an active workstream;
- needed to explain the current next action.

Each ledger entry should include:

- session heading timestamp;
- session ID when available;
- session started timestamp;
- last observed activity timestamp;
- short title;
- workstream;
- relation label;
- status;
- concise outcome summary;
- handoff or next action;
- repo-relative references.

Keep `Session Overview` in the same order as the detailed ledger.

### Untriaged Sessions

Add an entry to `untriaged-sessions.md` when classification is uncertain.

Use this inbox when:

- the session title is too vague;
- multiple workstreams are plausible;
- the session has no clear outcome;
- inspecting the raw log would require more context than the maintenance pass should spend;
- assigning the session would require guessing.

Preserve the session ID when available. Keep `Inbox Overview` in the same order as the detailed
inbox.

### Archive

Move old resolved entries out of the active ledger only when:

- the related goal or workstream has a current summary;
- no next action is hidden only in the old entry;
- the entry is no longer needed for active context recovery;
- the archived entry preserves the session ID when available.

Archive entries preserve state. A future agent should understand the outcome without reopening raw
session logs.

## Ordering And Priority

Use ordering to help the user resume thinking.

- Order goals and workstreams by curated relevance to current resumption needs.
- Reorder entries when dependencies, current work, or explicit user notes make the order clear.
- Ask the user when ordering would imply a meaningful priority decision and the right order is
  unclear.
- Record explicit ordering or priority decisions as `Ordering notes:` or `Priority notes:`.
- Within a workstream, keep session entries oldest to newest by `Session started:` and then
  `Last observed activity:`.
- In `untriaged-sessions.md`, keep entries newest first by `Last observed activity:`.

## Update Sequence

1. Build a bounded candidate list from curated memory, user-provided IDs, and raw metadata.
2. Verify current repository state for each candidate workstream.
3. Update `goal-register.md` for changed informal, pre-plan, or transitional goal or workstream
   state.
4. Update `session-ledger.md` for active, blocked, unresolved, recent, or active-workstream
   sessions.
5. Apply ordering rules.
6. Move eligible old resolved entries to archive ledgers.
7. Add uncertain sessions to `untriaged-sessions.md`.
8. Update overview sections in the same order as detailed entries.
9. Update `Last updated:` in every changed Markdown file.
10. Update `Last reviewed:` and `Coverage window:` in reviewed memory files. Coverage windows
    should include the latest inspected UTC timestamp.

## Stop Conditions

Stop and ask the user before writing a classification when:

- a session appears to contain sensitive material that should not be summarized into the repo;
- classification would require choosing between active goals with different next actions;
- raw session evidence contradicts current repository state and the correct state is unclear;
- a proposed update would create a new top-level goal outside the user's stated intent;
- the maintenance pass would need broad historical reconstruction instead of bounded update.

Use `untriaged-sessions.md` for non-blocking uncertainty.

## Validation

Run the smallest useful checks after maintenance:

```sh
git diff --check -- docs/agent-memory
rg -n "/U[s]ers/|/h[o]me/[^ ]+" docs/agent-memory
rg -n "\"timestamp\"|\"event_msg\"|\"response_item\"|\"role\"|\"content\"" docs/agent-memory
```

Manual checks:

- no raw transcript bodies were committed;
- changed files start with `Last updated:`;
- `Last reviewed:` and `Coverage window:` are current in reviewed files;
- coverage windows include precise timestamps;
- newer local sessions were checked for the repository;
- ledgers use grouped list entries, not Markdown tables;
- overview sections mirror detailed entry order;
- session IDs are preserved when available;
- references are repo-relative;
- current state was verified against repository evidence before writing next actions.
