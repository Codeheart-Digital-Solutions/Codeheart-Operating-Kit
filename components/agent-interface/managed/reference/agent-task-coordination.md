Last updated: 2026-09-07T22:40:39Z (UTC)

# Agent Task Coordination

Use this generic reference to commission, review, report and hand over agent work. It applies
without a particular app or an Organization Home installation. Durable roles and work records own
identity; a task ID only locates an execution conversation.

## Whole-Plan Assignment

Commission one implementer for the complete executable plan and its ordered epics. Before starting,
record these facts in the existing plan and assignment, using private locators only in the proper
audience's records:

- canonical plan name/path and owning project, repository and branch;
- intended outcome, ordered epics, constraints, non-goals and acceptance evidence;
- execution authority, coherent commit/normal-push/PR checkpoints and final delivery boundary;
- acceptance owner (the director), delegated decisions and required review points;
- who merges, releases or adopts, which effects are included, delegated or reserved, and the
  action-time inputs and gates still required;
- explicit report-back instruction and the commissioning task's execution locator.

A reviewed draft is not an execution request. An activation-only request grants the bounded
plan-checkpoint effects documented in planning workflows. A request to execute the whole plan
covers the implementation, validation, coherent commits, normal pushes and PR creation/updates
specified by that plan. Name the finish line: reviewed branch, merged main, released product or
adopted consumer. Never infer unspecified final effects from CI success or from an epic finishing.

Reuse authority while its target, scope and limits still apply. A specific mandatory confirmation,
exact-effect binding, signing/audience condition, security boundary or tool-enforced gate remains
binding. Determine whether the existing grant satisfies the actual contract before asking again.
Do not fabricate approval references or treat another rollout's approval as current authority.

Use an assignment such as:

```text
Execute the complete <Plan name> at <canonical path> in <owning project/repository>,
on <branch>, through its ordered epics. The approved scope includes implementation,
validation, coherent commits, normal pushes and one delivery PR with updates.
The finish line is <reviewed branch / merged main / released and adopted product>.
<Director role> accepts the first epic and final candidate; <owner> performs the
explicitly covered integration effects after <required checks and exact inputs>.
Send a concise direct message to <commissioning task locator> at required epic review,
completion, genuine blocker or material scope issue. Include Plan/epic names, outcome,
evidence, commit/branch/PR state and the decision or next action needed from the director.
Preserve the result and disclose failed message delivery. Do not add a watcher.
```

Fill the placeholders with real authorized scope before commissioning. Do not copy the example as
approval evidence. Keep consumer facts out of reusable public guidance.

## Execute And Review

Follow `../../planning-workflows/runbooks/execute-implementation-plan.md`. Every implementation
action belongs to a named epic and its outcome. Add necessary low-risk corrections there. Amend
the plan for a material outcome, scope, dependency, cost or authority change. A genuinely separate
routine repair follows `../../planning-workflows/runbooks/handle-routine-change.md` and retains a
visible relationship to affected work. Do not use it to escape epic review.

Use local commits for coherent completed work after proportionate inspection and checks. Follow
`../../planning-workflows/reference/planning-document-lifecycle.md` for exceptions, authorized
publication and accepted PR integration; a published draft remains independent of execution
authority. Do not create a review or approval layer per commit.

Use one independent source review at the planned meaningful checkpoint, combining related epics
when declared. The same reviewer follows corrections; broader review needs a material change,
independence issue, limitation or unresolved concern. Director acceptance need not repeat that
technical review. Retain applicable validation through corrections and interruptions and rerun
invalidated checks, with full coverage at the required coherent candidate boundary. Normally keep
one delivery PR and update it; do not create a task, release or PR per checklist item. When the
assignment requires director acceptance, report the coherent result to that director and hand off
for review. Do independent preparation that does not depend on the decision. Continue dependent
work when acceptance arrives. The director can request in-scope corrections and authorize the next
epic within its delegated mandate. Escalate only a decision outside that mandate or a binding gate
that genuinely requires the user's intervention.

## Report And Orient

At starts, meaningful transitions and completion, state the actual phase: discovery, plan amendment,
implementation within a named epic, or a separate routine repair. Name the outcome and next actor.
Use readable plan, epic, role and task names first; add stable IDs when needed to locate them.
Explain a new abbreviation on first use.

Report directly through the supported task-message surface when the assignment authorizes the
recipient and payload. Include result/evidence, Git state and requested director action. The
director can remain available for discussion while the implementer works. Best-effort report-back
is sufficient: if the tool rejects delivery, retain the result in the existing record and final
response and state that it was not delivered. Do not claim receipt or bypass the rejection.

No continuous watching, busy polling, scheduler, cron, heartbeat, retry service or callback framework
is required. Ordinary follow-up and a review handoff are enough. A sent review request is not
acceptance. An explicitly requested sustained-execution goal can support the same plan but adds
neither authority nor automatic cross-task notifications.

Describe actual delivery facts separately: uncommitted, committed, pushed, PR opened, merged,
released, adopted, and pilot-pending. Do not call a branch push a released product or claim that
technical validation proves a user's real-use experience.

## Naming, Handover And Archival

Recommend `[Role] Subject — Outcome` for bounded tasks and `[Role] Subject` for standing
responsibilities. Follow consumer naming conventions where present. Consumer-owned policy decides
model preferences; this reference introduces no role schema or model-routing engine.

Before archiving, accept the result or transfer unfinished work to a named owner with its current
state, evidence, remaining actions and required authority. Preserve durable context in the owning
record. Check active descendants and dependent assignments: transfer their report/review destination
or retain the parent until they have a reachable owner. Do not leave active work reporting to an
unavailable acceptance owner.

UI archival, organizational lifecycle and deletion/worktree cleanup are separate decisions, even
when an app couples some effects. Inspect actual app behavior and preserve required work before
archival. Archiving a conversation does not itself close a role or plan. Deletion and destructive
cleanup require their own applicable authority and preservation checks. For optional Codex-specific
surfaces and worktree consequences, read `codex-task-operations.md`.
