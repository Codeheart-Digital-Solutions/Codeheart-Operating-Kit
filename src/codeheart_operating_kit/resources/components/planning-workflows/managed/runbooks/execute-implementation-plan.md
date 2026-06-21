Last updated: 2026-06-21T15:07:29Z (UTC)

# Execute Implementation Plan

Use this runbook when executing an active `*_implementation_doc.md`, including goal-style Codex
runs.

Implementation plans are authoritative, but their checklists may not be exhaustive. Complete the
stated intention and outcome of each epic, not only the literal checkbox list.

## Trigger

Use this runbook when the user asks to:

- execute an implementation plan;
- implement a plan with goal-style autonomy;
- continue or resume an active `*_implementation_doc.md`;
- run the epics from an implementation document.

If the user asks only to review, explain, or refine a plan, do not execute it.

## Required Read Order

1. Read the root `AGENTS.md`.
2. Read the repository or product docs router that owns the plan path.
3. Read the planning lifecycle reference.
4. Read the referenced implementation plan.
5. Read the nearest product, component, package, or repo README for the plan scope.
6. Read matching runbooks for operational, release, migration, publishing, or controlled work.

Do not read every document by default. Follow the plan's essential context first and expand only
when the changed surface requires it.

## Plan Shape Check

Before execution, confirm the plan is in a valid planning shape:

- standalone plan file;
- plan bundle;
- subplan;
- plan family;
- program folder.

If the plan shape would require broad path churn before execution, stop and ask before moving
files.

## Lifecycle State

Before execution:

- confirm the plan header has `Status: active`;
- if the plan is `draft`, `completed`, `superseded`, or `archived`, stop and ask before executing.

During execution:

- keep the plan `Status: active`;
- tick checkboxes only after the work is completed and validated;
- update `Last updated` when the plan changes;
- add or preserve revision notes for meaningful scope, decision, strategy, or execution-plan
  changes.

After successful execution:

- set `Status: completed`;
- add `Completed: YYYY-MM-DD`;
- add or preserve `Execution log: <relative/path>`;
- do not archive the plan unless the user explicitly asks or the plan includes an approved archive
  task.

## Execution Contract

- Execute epics sequentially.
- Treat each epic outcome as the authority for completion.
- Treat checkbox tasks as planned execution aids.
- Add missing tasks when they are required for the epic outcome and are low-risk.
- Stop and ask when a new high-impact decision has no clear safe default.
- Preserve unrelated user work.
- Do not mark an epic complete until validation and the review gate pass.

## Safe Defaults

Use a safe default without stopping only when all of these are true:

- the decision is local to the current epic;
- the decision is reversible or low blast radius;
- the decision follows existing repository, product, or runbook rules;
- the decision does not change ownership, security, product boundaries, release authority, or
  external governance;
- the decision is needed to satisfy the epic outcome.

Stop before decisions that affect:

- architecture or product boundaries;
- durable docs or code path conventions;
- security, secrets, customer data, or tenant data;
- cloud accounts, deployment targets, or external systems;
- public repository settings, tags, releases, or permissions;
- destructive cleanup;
- scope that later epics depend on.

## Per-Epic Flow

For each epic:

1. Restate the epic outcome in practical terms.
2. Identify affected files, commands, and validation gates.
3. Implement planned tasks.
4. Add required low-risk tasks omitted by the checklist.
5. Run the smallest validation set that proves the outcome.
6. Run the per-epic review gate.
7. Fix material findings and repeat the review gate.
8. Update checklist state only for completed and validated tasks.
9. Update the execution log with meaningful divergence and review evidence.
10. Recap whether the epic intention is achieved.

## Per-Epic Review Gate

Before marking an epic complete, spawn a fresh read-only reviewer agent when the active environment
and user request permit reviewer-agent execution.

The reviewer checks the implemented epic against:

- epic outcome and acceptance criteria;
- completed and incomplete checklist items;
- validation evidence;
- execution-log state;
- scope boundaries and out-of-scope guardrails;
- accidental future-epic work.

Use the same default model and reasoning mode as the implementing agent unless the user requests a
different reviewer setup or the epic is unusually high-risk.

Fix material findings and repeat with a fresh read-only reviewer until no material issues remain
or a clear blocker is recorded. A material issue is anything that makes the epic incomplete,
misleading, out of scope, unvalidated, or not reproducible.

When reviewer-agent execution is unavailable, record why and run the strongest practical
main-thread review.

## Execution Log

Create or update a sibling execution log beside the implementation plan:

```text
<feature-slug>_implementation_doc.md
<feature-slug>_execution_log.md
```

For a plan bundle:

```text
<feature-slug>/
  <feature-slug>_implementation_doc.md
  <feature-slug>_execution_log.md
```

Create an execution log for goal-style implementation runs even when divergence is low.

## Execution Log Content

The execution log is a review surface, not a command transcript.

Log:

- extra tasks added because the checklist was not exhaustive;
- safe defaults chosen during execution;
- validation substitutions or extensions;
- review-gate rounds, findings, fixes, metrics, and accepted result;
- changed sequencing;
- corrected assumptions;
- meaningful plan wording changes caused by implementation reality;
- unresolved follow-ups.

Do not log:

- every command;
- routine file edits already covered by the checklist;
- timestamp-only edits;
- formatting-only edits;
- checklist progress without meaningful divergence.

## Execution Log Shape

Use prose and lists. Include:

- header with timestamp and created date;
- plan path;
- mode;
- status;
- overall divergence;
- summary;
- epic delta index;
- review gate metrics;
- per-epic delta sections;
- final validation when the plan is complete.

Review gate metrics should include:

- review gate required;
- review gate skipped status and reason;
- reviewer mode;
- reviewer model or reasoning mode when known;
- review rounds;
- material findings status;
- concise findings by round;
- whether files changed because of review;
- final accepted result;
- approximate added time when known;
- token usage when known;
- worth-it assessment.

## Relationship To The Plan

The implementation plan is the canonical execution state. The execution log is the quick review
surface for divergence and evidence.

When the plan and log conflict:

- trust the implementation plan for current scope, checklist state, lifecycle state, and
  acceptance criteria;
- use the execution log to understand what changed during execution and why;
- correct the conflicting document before closing the epic or plan.

## Plan Register Hook

When implementation execution changes a plan's lifecycle or material path, maintain the local plan
register. Material execution changes include activation, completion, supersession, archive state,
major implementation-path changes, new parent or child links, changed dependencies, changed
related-plan links, and execution handoff changes.

Use `maintain-plan-register.md` for the procedure and `../reference/plan-register-format.md` for
entry shape. The sequence is:

1. Update `docs/repo/plans/plan-register.md` in the local repository.
2. When portfolio coordination is configured and the coordination home is available and safe to
   edit, update the configured coordination-home register.
3. When portfolio coordination is configured but the coordination home is unavailable or unsafe to
   edit, record pending sync in `docs/repo/plans/coordination-sync-pending.md` and continue the
   local execution task.

The implementation plan and execution log remain the canonical execution state. The register is an
index snapshot and should not duplicate epic progress, validation details, review logs, or
execution evidence.

Record material-update session refs when a session ID is available. Do not block execution when no
session ID is available.

Do not update the register for typos, formatting-only edits, timestamp-only edits, or mechanical
checklist progress that does not change lifecycle, relationships, or implementation path.

## Final User Summary

After the full plan is achieved, summarize:

- overall divergence;
- epics completed;
- meaningful deltas by epic;
- safe defaults chosen;
- user decisions required;
- validation summary;
- open follow-ups.

Keep the final chat response concise. Do not paste the full execution log unless the user asks for
it.
