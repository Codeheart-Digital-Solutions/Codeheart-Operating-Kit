Last updated: 2026-09-07T22:40:39Z (UTC)

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
evidence is recorded, and lifecycle/catalog visibility is updated where required.

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
- confirm mixed/canonical mode metadata validates and identifies this implementation record;
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

## Whole-Plan Authority And Assignment

Read the whole delivery grant and the existing execution log before selecting effects. The
assignment should name the owning branch, ordered epics, acceptance owner, coherent commits,
normal pushes, PR creation/updates and final boundary. Follow
`../../agent-interface/reference/agent-task-coordination.md` for the complete assignment and
explicit direct report-back. Distinguish draft review, activation-only bookkeeping and whole-plan
execution. A whole-plan request covers the specified implementation, validation and Git checkpoints;
it is not a new user permission decision at each epic.

Merge, release, adoption and provider effects require included or delegated authority, their named
owner, exact action-time inputs and passing gates. Reuse valid authority while target, scope and
limits still apply. A specific mandatory fresh confirmation or binding tool gate remains binding.
Never invent approval evidence. At starts and transitions, name the phase, Plan/epic outcome and
next actor; use readable names before IDs.

## Activation Publication Preflight

When the current user request activates the plan, publish the bounded planning checkpoint before
source implementation:

1. resolve the exact repository, canonical plan, and unambiguous work branch;
2. include only the plan, execution log, and directly required planning metadata;
3. validate the planning checkpoint;
4. create or use the work branch, create an intentional plan-only commit, and normally push it;
5. record repository, branch, paths, commit, normal-push result, and exclusions without secrets.

The activation request itself is approval for this one plan-only branch/use/commit/normal-push
checkpoint; do not ask for a duplicate push prompt. Execution may continue on that active pushed
branch without waiting for a PR or merge.

Activation alone is not authority for unrelated dirty files, unauthorized implementation code, another plan,
PR creation, merge, release, force-push, branch deletion, history rewrite, destructive Git, or a
broader external action. Stop on ambiguity, overlapping dirty paths, auth failure, policy rejection,
or rejected normal push. Preserve and report the local checkpoint; do not bypass the rejection.

A whole-plan assignment may separately cover the later implementation and Git effects. Apply that
existing grant rather than asking again under this activation-only rule.

For a user-requested material update to an active plan, apply the same bounded checkpoint after the
plan/log change validates. Routine checkbox progress does not trigger a publication checkpoint.

## Execution Contract

- Execute epics sequentially.
- Treat each epic outcome as the authority for completion.
- Treat checkbox tasks as planned execution aids.
- Add missing tasks when they are required for the epic outcome and are low-risk.
- Return a material outcome, scope, dependency, cost or authority change to the delegated director
  for plan amendment; involve the user only outside that mandate or at a genuinely enforced user gate.
- Keep necessary low-risk corrections inside the named epic. A separate routine repair uses
  `handle-routine-change.md` with a visible relationship, never as an escape from epic acceptance.
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
- approval gates recognize existing sufficient authority and use explicit wording for any
  uncovered effect or specific mandatory fresh confirmation;
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

- declared script asset role matches the planned role and the script-promotion standard;
- runbook caller exists and does not duplicate full script internals;
- script placement follows the owning area's convention;
- existing owner docs/tests are reused; add only missing first-script contract and validation;
- workflow dependencies, phase boundaries, and blocker ownership are documented when the asset is
  a workflow script;
- helpers are placed at the narrowest durable owner boundary and do not act as hidden runbook
  entrypoints;
- the owner area's `scripts/README.md` has a compact role index when that improves review
  clarity;
- output matches its actual caller; use common structured fields only when needed by that caller;
- output safety describes emitted-output behavior;
- managed-runner, CI, or cloud portability constraints are satisfied when the plan requires those
  execution contexts;
- proportional tests or fixtures prove the contract;
- script does not broaden approval or target scope;
- wrapper, package, CLI, or API promotion has repeated-use rationale.

Record recipe validation evidence in the execution log. Include fresh-agent executability review,
non-live tests, dry-run or preflight, and approval-gated live validation only where the plan and
approval class require them.

## Coherent Batches, Validation And Recovery

Accumulate a coherent implementation batch before broad validation. During iteration, inspect
changes and use cheap affected checks that can expose useful defects. Run required broader
native, integration, compatibility or performance coverage at the planned candidate boundary.
Do not run a full product suite merely because a document changed, a local commit was made or
work resumed. Unknown dependencies warrant investigation and wider checks when justified.

Keep evidence while its relevant source, artifacts, configuration, environment and target remain
applicable. Rerun checks invalidated by corrections; a plan/log-only update does not invalidate
unrelated behavior. Explain retained evidence and its limits in the existing log. A cheap pass
never proves an executable surface it did not exercise or substitutes for required release gates.

After interruption, inspect current files, Git state, outputs and pending operations. Preserve
usable work. Classify a failure as code, test, environment, authentication or billing before repair.
Use tooling readiness for local prerequisites and the service owner for login/access/cost issues.
For an uncertain external write, establish whether it happened before retrying. Continue safe
in-scope fixes; do not blindly retry broad suites, increase budgets or invent an auth service.

At a useful completed checkpoint, locally commit the task's own inspected and proportionately
validated changes. Use `../reference/planning-document-lifecycle.md` for explicit exceptions,
authorized pushes/PR updates and normal accepted merge. Establish publication once and reuse its
authority. Commit, publication, lifecycle, release and adoption remain separate states.

## Safe Defaults

Use a safe default without stopping only when all of these are true:

- the decision is local to the current epic;
- the decision is reversible or low blast radius;
- the decision follows existing repository, product, or runbook rules;
- the decision does not change ownership, security, product boundaries, release authority, or
  external governance;
- the decision is needed to satisfy the epic outcome.

Resolve new decisions outside the approved plan with the acceptance owner before changing:

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
7. Run the planned meaningful review checkpoint; related epics may share one coherent review.
8. Fix material findings and return affected corrections to the same reviewer.
9. Update checklist state only for completed and validated tasks.
10. Update the execution log with meaningful divergence and review evidence.
11. Make the agreed coherent commit/normal-push/PR checkpoint, reporting each actual state accurately.
12. Send the named epic result, evidence, Git state and requested decision directly to the assigned
    director at required review, completion or genuine blocker. Retain the result and disclose any
    message rejection; best-effort reporting needs no watcher, heartbeat or busy polling.
13. At a required acceptance point, complete independent preparation then hand off for review.
    Continue dependent epic work once delegated acceptance arrives. Do not require a new task,
    release, PR or user authorization per epic.
14. Recap whether the epic intention and acceptance are achieved; keep the plan and any whole-plan
    goal incomplete while required release/adoption outcomes remain.

## Meaningful Review Checkpoints

Before accepting the planned coherent source checkpoint, use one independent read-only reviewer
when the active environment and user request permit reviewer-agent execution. Closely related
epics may share that checkpoint when declared in the plan. Keep those epics acceptance-pending
until review passes. Director acceptance need not duplicate the technical review.

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

Fix material findings and use the same reviewer for focused follow-up on corrections and their
effects. Preserve valid review evidence. Broaden review or use a different reviewer only for a
material design/impact change, inadequate independence, a reviewer limitation or unresolved concern.
Continue until no material issues remain or a clear blocker is recorded. A material issue is anything that makes the epic incomplete,
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

Keep the timestamp, creation date and plan link, then summarize outcomes, meaningful divergence,
validation and actual delivery state in prose or a compact table. Record the review scope,
material findings, corrections, reviewer continuity and accepted result. Name residual limits and
remaining owner decisions. Do not require separate metrics, duplicate evidence artifacts or a
command transcript; include timing or measured cost only when useful and actually known.

## Relationship To The Plan

The implementation plan is the canonical execution state. The execution log is the quick review
surface for divergence and evidence.

When the plan and log conflict:

- trust the implementation plan for current scope, checklist state, lifecycle state, and
  acceptance criteria;
- use the execution log to understand what changed during execution and why;
- correct the conflicting document before closing the epic or plan.

## Catalog And Visibility Hook

The implementation plan and execution log remain canonical execution state. On activation,
completion, supersession, archive, or user-requested material active-plan change:

1. update the canonical header/content and metadata chronology as applicable;
2. run `codeheart-operating-kit plans validate` and inspect `plans list`;
3. keep mixed/canonical `plan-register.md` frozen and avoid new pending-sync files;
4. follow `maintain-plan-register.md` only for legacy compatibility; and
5. use the bounded publication preflight above when the user requested activation or a material
   update to an active plan.

After a successful push, coordination visibility still depends on the next complete home refresh.
After push failure or before push, state that the checkpoint is local-only. Completion itself
does not authorize additional effects; execute covered PR, merge or release work under the existing
whole-plan grant and its required gates. Keep committed, pushed, PR opened, merged, released,
adopted and pilot-pending states distinct.

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
