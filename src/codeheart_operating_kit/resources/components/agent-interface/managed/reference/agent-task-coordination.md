Last updated: 2026-10-10T16:09:55Z (UTC)

# Agent Task Coordination

Use this generic reference to commission, review, report and hand over agent work. It applies
without a particular app or an Organization Home installation. Durable roles and work records own
identity; a task ID only locates an execution conversation.

## Tool Selection

Honor the assignment's explicit tool choice and applicable accepted preferences; otherwise
continue with the current tool. Work in Codex continues with Codex and work in Claude Code with
Claude Code, using that tool's own delegation and reporting surfaces. No second tool, account,
relay or coordination reference is needed, and installing the Kit never creates a chat by itself.

Mention cross-tool delegation only when the user asks or it is materially useful. A decline holds
for the current conversation or assignment unless the user agrees a wider scope. When the user
selects another tool's agent, follow `../runbooks/coordinate-cross-tool-task.md`; its tested
direction is a Codex coordinator commissioning Claude CLI sessions. Reverse execution is not yet
qualified.

## Whole-Plan Assignment

Commission one implementer for the complete executable plan and its ordered epics. Before starting,
record these facts in the existing plan and assignment, using private locators only in the proper
audience's records:

- canonical plan name/path and owning project, repository and branch;
- selected workspace, closure owner and any covered workspace-removal authority, following
  "Assignment Workspaces" below;
- intended outcome, ordered epics, constraints, non-goals and acceptance evidence;
- execution authority, coherent commit/normal-push/PR checkpoints during work, independently
  usable integration outcomes and their dependencies, and final delivery boundary;
- for repository adoption, named default/integration branches plus worktrees to reconcile, any
  held-ancestry publication limits, and covered final record publication;
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

## Assignment Workspaces

A temporary workspace belongs to an assignment; deciding its disposition is part of finishing
that assignment. These rules cover linked Git worktrees and independent clones, for whole-plan
execution and routine changes alike. They introduce no registry, cleanup service or new settings
schema.

### Select And Locate

- Standing coordinators and directors normally use the usual checkout for discussion, reading,
  review and coordinated drafting. A standing role does not need a permanent worktree.
- Use an isolated workspace for an implementation assignment, concurrent editing, or a different
  branch that would disturb other chats. A small sequential change may use the usual checkout
  when its files and branch can be changed without disrupting anyone else. Coordinate one writer
  per checkout; do not switch a shared checkout's branch underneath dependent chats or processes.
- Inspect existing workspaces first and reuse a suitable one for the same assignment. A new epic
  or review alone does not require another workspace. Prefer linked worktrees; record a concrete
  reason for an independent clone, such as an isolated comparison experiment.
- For manually created workspaces, use the repository's or user's chosen workspace root, outside
  the usual checkout and other repository checkouts. A readable convention is
  `Workspaces/<repository>/<assignment>/`; a date prefix is optional. Reuse an existing local
  convention, or select a suitable sibling location and state it when none exists. Keep actual
  machine paths in ignored local settings or local assignment evidence. This is a guidance
  preference, not a new Kit configuration key or a requirement to create a coordination reference.
- App-managed workspaces stay in their supported locations and use the app's lifecycle tools.
  Do not move them solely to match the manual convention. Existing workspaces do not need a
  relocation merely to adopt this guidance.

Record the workspace kind, repository/branch and closure owner in the existing assignment or
execution record; keep its exact local path in the appropriate local evidence. The commissioning
agent owns closure by default; for undelegated work, the agent doing the work owns that step.
Establish whether the actual grant covers removing this assignment's exact temporary workspace
after the conditions below hold. Reuse sufficient authority without asking again; the record
does not create authority. Branch deletion, other workspaces, force operations and unexpected
unique material are outside that bounded removal grant. Tool-enforced gates still apply.

While working, place deliverables intended to outlast the assignment in their lasting home: a
committed record, an agreed document folder or the shared evidence location. Standing chats,
shared helpers and their sole configuration must not depend on a temporary assignment workspace.

### Close Or Retain

At final handoff, the implementer reports the workspace's branch/commit state, useful untracked
or ignored material and any dependent sessions or processes. The closure owner verifies:

1. The outcome is accepted and integrated, or explicitly superseded with useful work preserved.
   Required release, adoption, review and follow-up work is finished or transferred to a named
   owner with a usable workspace; a merged PR alone does not establish this.
2. Unique commits, unfinished changes, documents and useful untracked or ignored evidence have
   a verified durable or recoverable home outside the workspace. Inspect ignored files as well
   as Git status. Regenerable build outputs need no routine backup.
3. No running process, continuing chat, active report-back destination or shared tool still
   depends on the folder. Preserve conversation continuity through the tool's supported handoff
   or archival route before removing it; an idle chat alone is not proof of independence.
4. Removal is covered by applicable authority. Remove only the identified temporary workspace
   through its supported app or Git operation, without force. Inspect any coupled branch or chat
   effects first. Preserve the branch unless its deletion is separately covered; never remove
   the usual checkout as assignment cleanup.

Record the result as removed or retained. For retention, one line in the existing record gives
the reason, responsible owner and next resolving event; no expiry timer or watcher is required.
If preservation or removal fails, retain the workspace and report the specific blocker. Report
delivery completion and workspace disposition separately.

After integration, refresh the responsible agent's usual checkout by fast-forward only when its
branch and dependent sessions permit and existing local work is preserved. If this is not safe or
possible, report the divergence or collision and next owner; do not reset, switch a shared branch,
or stash unrelated changes merely to complete closure.

## Native Goal By Default

Executing an approved implementation plan uses the implementing tool's native goal mode by
default, for Codex and Claude Code alike. This applies to complete implementation assignments,
not routine changes, discovery, consultation or review. Honor an explicit user opt-out.

Commission the goal together with the execution grant: concrete objective, agreed finish line or
phase handoff, constraints, review points and report-back. Carry the actual user authority,
including an accepted standing goal preference. Installing the Kit is not consent that overrides
a host's explicit-authorization requirement; resolve any missing goal authority in the same
commissioning decision, not again at every epic or routine continuation.

Activate through the tool's supported mechanism and confirm it from observable status or
evaluator evidence; a prompt mentioning `/goal` is not proof. Aim each goal at a real handoff,
such as source review, so required review or a genuine blocker returns control with original
evidence instead of spinning. A goal ending at a checkpoint, a process exit or an evaluator
verdict is not plan completion. After delegated acceptance, continue with the tool's supported
continuation for the next phase. If activation is unavailable or cannot be verified, say so and
continue only under authorized ordinary execution; never claim a goal is active. Goals add no
permissions, notification channel or scheduler. Tool details: `codex-task-operations.md` and
`claude-code-task-operations.md`.

## Execute And Review

Follow `../../planning-workflows/runbooks/execute-implementation-plan.md`. Every implementation
action belongs to a named epic and its outcome. Add necessary low-risk corrections there. Amend
the plan for a material outcome, scope, dependency, cost or authority change. A genuinely separate
routine repair follows `../../planning-workflows/runbooks/handle-routine-change.md` and retains a
visible relationship to affected work. Do not use it to escape epic review.

Use local commits for coherent recoverable progress before formal acceptance, after proportionate
inspection and checks. Preserve incomplete progress before pauses and inspect relevant target
changes on substantial starts/resumes. Follow
`../../planning-workflows/reference/planning-document-lifecycle.md` for exceptions, authorized
publication and accepted PR integration; a published draft remains independent of execution
authority. Do not create a review or approval layer per commit.

The commissioning agent is the primary reviewer and owns acceptance. At the planned meaningful
checkpoint, combining related epics when declared, it examines the delivered work against the
intended outcome, scope and validation evidence, in proportion to risk. It may rely on competent
existing evidence and targeted checks. Forwarding another agent's verdict without examining the
work is not a review. It adds an independent reviewer when complexity or risk warrants it, when
it took a substantive part in the implementation, or when the user, assignment or a binding gate
requires it; another full review is not automatic. Checking one's own implementation is
self-review, never independent review. The commissioning agent assesses any additional findings
itself. Review is tool-neutral: a Codex director can review a Claude implementer's work directly,
and same-tool work is equally valid, with no extra account, model, relay or chat. Keep reviewers
the user selected and binding gates. The same reviewers follow corrections; broader review needs
a material change, independence issue, limitation or unresolved concern. Retain applicable validation through corrections and interruptions and rerun
invalidated checks, with full coverage at the required coherent candidate boundary. Choose PR boundaries around usable accepted outcomes: update one coherent PR or integrate several
independent outcomes during a longer plan. Do not create a task, release or PR per checklist item.
Integrate ready work under covered authority and passing gates, or record a concrete reason,
owner and next resolving event in the existing record. When the
assignment requires director acceptance, report the coherent result to that director and hand off
for review. Do independent preparation that does not depend on the decision. Continue dependent
work when acceptance arrives. The director can request in-scope corrections and authorize the next
epic within its delegated mandate. Escalate only a decision outside that mandate or a binding gate
that genuinely requires the user's intervention.

When an additional reviewer reports on consequential work, the coordinator reads the original
response and reports its locator,
core findings and significance, its own assessment including disagreement, recommended
disposition, verdict and residual uncertainty; a bare "review passed" is not enough. During
discovery, planning or discussion, report before incorporating findings and wait for the user
unless an explicit delegation covers the correction cycle; see
`../../planning-workflows/runbooks/review-planning-document.md`. Routine review corrections during
authorized implementation stay within the existing grant.

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
acceptance. A native implementation goal supports the same plan but adds neither authority nor
automatic cross-task notifications.

Describe actual delivery facts separately: uncommitted, committed, pushed, PR opened, merged,
released, adopted on each named branch, worktree-reconciled, and pilot-pending. Local installed
health and a published adoption branch do not establish default-branch rollout. An adoption commit
on held product ancestry may remain local with its owner/resumption exception recorded. Do not call a branch push a released product or claim that
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
cleanup require their own applicable authority and preservation checks; the bounded authority
established under "Assignment Workspaces" can cover that exact removal without another approval.
For optional tool-specific
surfaces and worktree consequences, read `codex-task-operations.md` for Codex or
`claude-code-task-operations.md` for Claude Code. For commissioning another tool's agent, use
`../runbooks/coordinate-cross-tool-task.md` and `cross-tool-coordination-contract.md`.
