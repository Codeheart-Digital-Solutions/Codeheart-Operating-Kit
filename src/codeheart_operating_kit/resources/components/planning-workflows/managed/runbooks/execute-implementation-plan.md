Last updated: 2026-06-29T19:51:37Z (UTC)

# Execute Implementation Plan

Use this runbook when executing an active `*_implementation_doc.md`, including goal-style Codex
runs.

Implementation plans are authoritative, but their checklists may not be exhaustive. Complete the
stated intention and outcome of each epic, not only the literal checkbox list.

Audience: agent-facing

Intent:
Execute an active implementation plan linearly, complete each epic's intended capability, record
meaningful divergence and evidence, and stop before unapproved scope or authority changes.

Success:
The implemented plan satisfies its stated outcomes, validation and review gates pass, execution
evidence is recorded, and lifecycle/register state is updated where required.

Agent judgment boundary:
The agent may add low-risk tasks required for the epic outcome and choose safe local defaults. It
must not expand scope, skip validation, mark incomplete epics complete, or bypass approval,
release, safety, or routing-standard gates.

Stop boundary:
Stop before executing draft, completed, superseded, or archived plans; before high-impact
unplanned decisions; and before closing an epic whose evidence does not prove the intended
capability.

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

## Runbook Change Execution

When an epic creates or materially changes durable runbooks, use
`../../agent-interface/reference/runbook-authoring-standard.md` before marking that epic complete.

Verify the changed runbook surface against the plan and standard:

- audience classification is explicit;
- required compact intention blocks are present;
- human-facing runbooks provide the user-visible flow, question pacing, and approval wording;
- agent-facing runbooks provide the concrete execution path, evidence, validation, and stop
  conditions;
- hybrid runbooks clearly separate user dialogue from agent-only execution;
- maintainer-facing runbooks preserve authority, evidence, rollback, and validation boundaries;
- approval gates use explicit user-facing wording before external-state-changing actions;
- runbooks that can hit missing local tooling route generic environment blockers to
  `../../agent-interface/runbooks/handle-tooling-readiness.md`;
- module-specific install commands, version requirements, service authentication, and live
  preflight remain in module-owned guidance;
- broad package-manager, runtime, or local-tool setup guidance is not copied into multiple
  managed runbooks unless the active plan explicitly changes the shared tooling-readiness route;
- planned scope is honored without accidental retrofit of unrelated runbooks.

Fix gaps within the approved epic scope before completing the epic. If the fix requires a broader
runbook audit, durable format change, ownership change, or new authority decision, stop and amend
the plan instead of expanding execution ad hoc.

## Recipe-Maturity Execution

When an epic creates or materially changes durable operational recipes, route-selected recipes,
durable executable mechanics, expected markers, structured summary or blocker output, reusable
script assets, promoted recipe assets, or recipe validation expectations, use
`.codeheart/kit/docs/agent-interface/reference/operational-recipe-maturity.md` before marking
that epic complete.

When an epic creates or materially changes reusable script assets, first-script scaffolding,
script output contracts, script tests, script helper rules, or script-promotion review criteria,
use `.codeheart/kit/docs/agent-interface/reference/runbook-to-script-promotion-standard.md`
before marking that epic complete.

Verify the changed recipe surface against the plan and reference:

- target maturity state is named and reflected in the artifact;
- recipe metadata, validation tier, evidence shape, and blocker shape match the planned level;
- promotion destination or non-promotion decision is recorded;
- promoted assets have an owner, placement boundary, validation path, and discoverability route;
- runbooks remain the operator-facing entry point after promotion;
- domain blocker classes and concrete package layouts remain with the owning domain when they are
  outside the generic plan;
- approval, secrets, public-core, and external-state boundaries are preserved.

For reusable script assets, also verify:

- runbook caller exists and does not duplicate full script internals;
- script placement follows the owning area's convention;
- first-script scaffolding exists when this is the first script in the owner area;
- output contract includes required common fields;
- output safety describes emitted-output behavior;
- proportional tests or fixtures prove the contract;
- script does not broaden approval or target scope;
- wrapper, package, CLI, or API promotion has repeated-use rationale.

Record recipe validation evidence in the execution log. Include fresh-agent executability review,
non-live tests, dry-run or preflight, and approval-gated live validation only where the plan and
approval class require them.

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
6. For routing-bearing epics, run or verify the planned fresh low-context routing probe, or record
   why the probe is not applicable.
7. Run the per-epic review gate.
8. Fix material findings and repeat the review gate.
9. Update checklist state only for completed and validated tasks.
10. Update the execution log with meaningful divergence and review evidence.
11. Recap whether the epic intention is achieved.

## Per-Epic Review Gate

Before marking an epic complete, spawn a fresh read-only reviewer agent when the active environment
and user request permit reviewer-agent execution.

The reviewer checks the implemented epic against:

- epic outcome and acceptance criteria;
- completed and incomplete checklist items;
- validation evidence;
- execution-log state;
- scope boundaries and out-of-scope guardrails;
- delivered feature capability, not only completed checklist lines;
- discovery capability scope when the implementation plan references one;
- routing-standard adoption and probe evidence for routing-bearing epics;
- accidental future-epic work.

Use the same default model and reasoning mode as the implementing agent unless the user requests a
different reviewer setup or the epic is unusually high-risk.

Fix material findings and repeat with a fresh read-only reviewer until no material issues remain
or a clear blocker is recorded. A material issue is anything that makes the epic incomplete,
misleading, out of scope, unvalidated, not reproducible, narrow, policy-only, stubbed, unusable,
or incomplete against the intended feature capability. If the gap is within the approved epic
scope, fix it. If fixing it requires a new high-impact decision or scope expansion, stop and
amend the plan instead of improvising.

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
- routing-standard adoption evidence and fresh low-context probe results for routing-bearing
  epics;
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
2. When portfolio coordination is configured, use the target-register compatibility test in
   `maintain-plan-register.md` before choosing direct coordination-home update versus pending
   sync.
3. When the coordination-home register update is compatible, update the configured
   coordination-home register.
4. When the coordination home is unavailable, unwritable, or unsafe under that compatibility test,
   record pending sync in `docs/repo/plans/coordination-sync-pending.md` and continue the local
   execution task.

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
