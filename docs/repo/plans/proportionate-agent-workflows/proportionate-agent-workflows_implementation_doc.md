Last updated: 2026-09-05T21:06:54Z (UTC)
Created: 2026-09-05
Status: active
Execution log: docs/repo/plans/proportionate-agent-workflows/proportionate-agent-workflows_execution_log.md

# Proportionate Changes, Approval Reuse and Agent Task Coordination

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.proportionate-agent-workflows
  kind: implementation
  purpose: Ship proportionate routine work and whole-plan execution with predictable Git checkpoints, delegated review, simple agent report-back and clear communication.
  first_cataloged: 2026-09-05T16:59:13Z
  catalog_metadata_updated: 2026-09-05T18:02:06Z
```
<!-- END CODEHEART PLAN METADATA -->

## Overview

Deliver a shipped Kit workflow through which an implementer executes an entire approved multi-epic
plan under one execution authorization, with predictable Git checkpoints and delegated reviews.
Routine work uses short existing records. Assignments include a simple report-back instruction;
user interaction is reserved for material exceptions and decisions outside the agreed authority.

This active plan translates the accepted capability into a self-contained public implementation
scope. It contains no private consumer facts and has no private-document dependency. The delivery
authorization and delegated acceptance boundary are recorded in 2.4; consumer-specific assignment
details remain in the commissioning task.

## Essential Context Reference Files

All source paths below are relative to this repository.

| File | Why it matters |
| --- | --- |
| `AGENTS.md` | Source ownership and public-core boundary. |
| `docs/repo/runbooks/change-operating-kit.md` | Consumer-impact and smallest applicable validation gates. |
| `docs/repo/runbooks/release-operating-kit.md` | Current release, platform, signing and publication requirements. |
| `docs/repo/reference/consumer-impact-classification.md` | Classify the final change honestly. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Owner-first routing and controlled-operation authority. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Existing broad approval wording to align. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Formal planning, activation and conflicting title instructions. |
| `components/planning-workflows/managed/reference/plan-catalog-format.md` | Canonical semantic identity and meaningful title. |
| `tests/test_packaging_resources.py` | Source/mirror equality and consumer materialization. |

## Table Of Contents

- [Foundation](#section-1---foundation)
- [Strategy](#section-2---strategy)
- [Execution Plan](#section-3---execution-plan)
- [Future Planning](#section-4---future-planning)
- [Revision Notes](#revision-notes)

# Section 1 - Foundation

## 1.1 Goal

A consumer can commission an executable whole plan once, have an implementer progress through its
ordered epics with evidence and delegated acceptance, and understand its actual position without
repeated prompts for ordinary implementation steps. Routine fixes retain a proportionate path.
The generic result must work without Organization Home and without Codex.

## 1.2 Context and accepted boundaries

**Act through the plan.** Every implementation action belongs to a named epic and its intended
outcome. Necessary low-risk corrections within that outcome are executed and recorded there;
they do not become separately commissioned micro-deliveries. A changed outcome, significant scope,
dependency or authority requires an explicit plan amendment. A genuinely separate routine change
uses the routine route and a visible relationship to affected work. Never silently relabel plan
work as a routine exception to escape its review, controls or acceptance.

**Authorize execution once.** The reviewed plan identifies the whole delivery outcome, ordered
epics, owning repository/branch, acceptance evidence, delegated decisions and planned Git effects.
A request to execute that whole plan authorizes the covered implementation, validation, commits,
normal branch pushes and PR creation/updates at its planned delivery checkpoints. Planning must
surface these ordinary steps together, rather than leave permission blanks for the implementer to
rediscover. Mere draft review or activation-only bookkeeping is distinct from whole-plan execution.
Record this distinction in the existing plan/assignment; introduce no new approval form or registry.

State the final delivery boundary up front: reviewed branch, merged main, released product or
adopted consumer. Name who integrates and whether merge/release/install/provider effects are
included, delegated or explicitly reserved. Covered effects proceed after their predicates pass;
CI success alone does not grant an unspecified effect. Reuse valid existing authorization across
steps while target, scope and limits remain applicable. A tool-enforced gate remains binding.
Do not invent approval evidence or treat an approval for one rollout as authority for another.

**Use epics as agent checkpoints.** One implementer can execute all ordered epics under that
assignment. Where director acceptance is required, send the result to the designated director;
the director reviews within its delegated mandate and authorizes continuation without another
user decision. Do not require a user approval, new task, PR or product release for every epic.
Use coherent validated commits and ordinary pushes at the plan's agreed review, recovery and
handover boundaries; update an existing PR rather than manufacture one per checkbox. Report
uncommitted, committed, pushed, merged and released states accurately. Avoid the ambiguous word
"published" when the action is only pushing a branch and opening a PR.

**Simple report-back.** The assignment explicitly instructs the implementer to send a concise
message to the identified commissioning task at required epic review, completion, a real blocker
or material scope issue. Include the plan/epic by name, result/evidence, Git state and requested
next action. Return to the director for delegated decisions; escalate to the user only outside
that mandate. The director stays available for discussion while implementation proceeds. No
continuous waiting, scheduled automation, cron, watcher, retry service or callback framework is
required. Use the supported task-message tool under the commissioned authority. If delivery is
blocked, retain the result and state that the message was not delivered; do not claim receipt.
Best-effort messaging and ordinary follow-up are sufficient for this first delivery.

**Communicate the same structure.** At starts, meaningful transitions and completion, explain
whether work is discovery, plan amendment, implementation within a named epic, or a separate
routine repair. Name the outcome and next responsible actor. Use readable names first and stable
identifiers second; explain new abbreviations on first use. These are short orientation updates,
not permission requests or repeated ceremonial headers.

**Use goals for sustained execution.** The optional Codex note shows one user-authorized goal for
an entire executable multi-epic plan, referring to its canonical file, finish line, constraints,
review/report-back instructions and Git boundary. Verify goal activation through the supported
mechanism; a slash-command string in another task's prompt does not itself prove activation. Do
not infer a goal or token budget from discussion. A goal supports continuation; the plan remains
authoritative and goal state creates no new authority or automatic cross-task notification.

Routine eligibility requires a known owner and method, bounded understood impact, reversibility
and sufficient current authority. Lines/files are not the classifier. Use an existing PR, issue
or owning record for purpose, scope, authority, result and evidence; create a short record only
when none fits. Promote uncertain or consequential work to formal planning. Preserve specific
module, managed-content, external-effect and security contracts.

Generic assignments retain the owning project, evidence, authority and acceptance owner. Recommend
`[Role] Subject — Outcome` for bounded tasks and `[Role] Subject` for standing responsibilities.
Task IDs locate execution; durable roles and work records retain identity. Before archival, accept
the result or transfer unfinished work, preserve needed context and check active descendants.
UI archival, organizational lifecycle and deletion/worktree cleanup remain distinct. Model
preferences remain consumer-owned; no role schema or model-routing engine is introduced.

## 1.3 Current state

Planning base: `be4b92e5f997abb3ba7b7de03f85937c83ef8193`, version `0.1.30`.
Managed guidance is declared by component manifests and mirrored under
`src/codeheart_operating_kit/resources/`. The Kit already supports proportionate judgment in parts
of its guidance, but lacks an explicit end-to-end routine route. The runbook-authoring standard
contains broad ask-before-write language; formal activation already recognizes covered plan-only
publication. Align these instructions and whole-plan commissioning without changing CLI permission semantics.
The existing activation rule covers only a plan-checkpoint push, leaving later implementation Git
effects easy to fragment into fresh questions. The new whole-execution route must explicitly
resolve that distinction across drafting, dispatch, execution and review. Because it changes the
interpretation of external-action authority, classify that slice as a `security or safety policy
change` under the existing impact reference and include explicit review and release notes; do not
label the entire delivery instruction-only merely because its implementation is Markdown.

The drafting runbook still requires a generic `Document Header` H1 while the catalog expects a
meaningful title. This plan uses the meaningful title and retains the required planning sections;
EP-01 corrects the producer inconsistency. Installed managed content is not authoring authority.

Authoring validation accepts this new Plan. The configured repository-wide validator reports an
existing `legacy_relation_duplicate` for frozen register entry `OK-PR-002`; the same result was
confirmed with the current source CLI. The register is unchanged from the planning base. The user has explicitly authorized a bounded catalog prerequisite correction as part of this
complete delivery. Activate the plan checkpoint with the unchanged baseline error recorded, then
resolve that prerequisite before the remaining EP-01 source work. Require zero new Plan errors and
a passing configured catalog after the correction; this is not a general validator waiver. Do not
follow the generic error suggestion by editing frozen history. No catalog migration is
part of this delivery. The prospective discovery-v2 preview also encounters existing test fixtures;
only its new-Plan preview is used here, without changing repository configuration.

# Section 2 - Strategy

## 2.1 File tree and ownership

```text
components/planning-workflows/
  component.yaml
  managed/README.md
  managed/runbooks/handle-routine-change.md                 [new]
  managed/runbooks/draft-implementation-plan.md
  managed/runbooks/execute-implementation-plan.md
  managed/runbooks/review-planning-document.md
components/agent-interface/
  component.yaml
  managed/README.md
  managed/reference/operation-routing-and-dispatch.md
  managed/reference/runbook-authoring-standard.md
  managed/reference/agent-task-coordination.md              [new, generic]
  managed/reference/codex-task-operations.md                [new, optional app note]
templates/agents/AGENTS.managed-block.md
src/codeheart_operating_kit/resources/                     [matching declared mirrors]
tests/test_packaging_resources.py
internal/plancatalog/legacy.go                            [bounded prerequisite]
internal/plancatalog/legacy_projection_test.go            [focused regression evidence]
release-notes.md
```

Keep root guidance short: one routine route and pointers to owner references. Verify that the
consumer entry route can discover its installed module inventory; do not duplicate each module's
policy in a new root managed block. Shared installer machinery remains outside this delivery.
 The optional Codex
note describes supported task operations and links current official documentation; app permissions
and available tools decide execution. No task creation, archival or cleanup is performed merely
because this plan documents a convention.

## 2.2 Open questions

| ID | BLOCKER | Affected epic | Default and resolution |
| --- | --- | --- | --- |
| OQ-01 | no | EP-02 | Select the next available patch version and exact release audience/commit after the source diff passes. Recheck current release state and the release runbook; retain a candidate until applicable platform, signing and publication requirements are satisfied. This does not block EP-01. |
| OQ-02 | no | EP-02 | One combined fresh-context walkthrough can be shared with the downstream module delivery. Record who actually performed it and which installed/candidate versions were read. Do not claim independent review from a main-task read-through; missing fresh-context access stays a visible pilot-readiness gap. |

## 2.3 Decisions

| Decision | Rationale and future consequence |
| --- | --- |
| Add a routine route with existing-record evidence. | Makes the small-work path discoverable without a new schema, ledger or mandatory mini-plan. Formal plans remain available for consequential work. |
| Execute the whole authorized plan with delegated reviews. | Plans declare ordinary Git effects and final acceptance up front; epics guide actual execution and agent review without fresh user prompts. Explicitly review the changed external-action doctrine while preserving enforced permissions. |
| Instruct implementers to report back. | A direct message at completion, blockage or required director review keeps the director available. Best-effort reporting is accepted; no continuous watcher, scheduler or notification framework. |
| Use readable plan/epic names and optional whole-plan goals. | Communication follows actual execution state. Goal mode supports sustained execution of the canonical plan; it neither replaces that plan nor creates authority. |
| Separate generic coordination from the optional Codex note. | Consumers can use the Kit without a particular app. App behavior must be rechecked when operating it. |
| Use focused checks plus practical experience. | Reuse materialization and packaging tests for breakable contracts; use one concise walkthrough for instruction clarity and consumer feedback for everyday usability. No prompt benchmark, UI suite or lifecycle simulator. |
| Ship through the existing release route. | A repository edit does not make a reusable improvement available to consumers. Required release checks remain; this delivery creates no release machinery. |

## 2.4 Approved whole-delivery authority and execution ownership

On 2026-09-05, the user approved the proposed sequence and one delivery authorization covering both
complete implementation plans. This Plan is active. One implementer owns all its ordered epics;
the commissioning CEO/Program Director owns acceptance, cross-plan coordination and integration.
The task assignment supplies execution locators without making task IDs durable role identities.

The grant covers in-scope implementation and necessary corrections; work branches, coherent
commits, normal pushes and PR creation/updates; Director review and continuation; merge to the
source repository's `main` after required checks and acceptance; version selection and publication
of a validated release using existing procedures; and adoption into the named pilot consumer with
preservation and workflow checks. The Director selects and approves the exact candidate, target,
starting state and other action-time inputs within this grant. Resolve those inputs and retain
their evidence; do not substitute fabricated approval references or ask the user again merely
because the approved plan has reached its release or adoption epic.

The implementer commits and pushes at coherent review/recovery boundaries and opens or updates
one delivery PR as appropriate. It reports the first epic result to the Director for acceptance,
then continues through the second epic after that delegated review. Release preparation may
proceed while awaiting review when it does not depend on acceptance. Merge, final release and
consumer adoption are performed by the Director or explicitly delegated to the implementer after
review of the exact candidate. There is no routine user gate per epic or Git operation.

Report back directly to the commissioning task at required review, completion, a genuine blocker
or material scope issue. Include Plan/epic names, evidence, Git state and the needed decision.
Use a simple supported task message; preserve the result and state any delivery failure honestly.
Do not add continuous waiting, scheduled automation or a reporting framework. A whole-plan goal
may be used when explicitly requested in the commissioning task and verified active; no token
budget is inferred. The plan remains the execution authority and cannot be marked complete merely
because source edits are done or the implementer awaits review.

The finish line is a released product, adoption in the specifically authorized pilot consumer,
technical verification and a clear user-pilot handoff. The coordinated pilot consumer is supplied
in the assignment. User experience remains pilot-pending until actually observed. Escalate only
material outcome, scope, cost or authority changes and enforced gates requiring the user's own
decision. Existing required checks and repository protection remain binding. No unrelated provider
publication, live infrastructure operation, public audience expansion or signing-policy waiver is
granted by this commissioning.

**Bounded catalog prerequisite.** The configured canonical catalog has an unchanged baseline
`legacy_relation_duplicate` in frozen entry `OK-PR-002`. Inspection shows two different prose
relations beginning with the same word and followed by different wrapped document paths; the
parser currently projects their first word as the target. Establish the correct classification
with focused evidence before choosing a fix. This grant permits the smallest catalog parser or
projection correction needed to remove the false blocker, with regression tests and release notes.
Preserve the frozen register, existing schemas/configuration, schema-v1 compatibility, strict
legacy/mixed authority, actual duplicate relation guards and ambiguous identity/path guards.
Do not broadly downgrade duplicate errors or migrate discovery modes. Record this as validator-only
impact in addition to the main instruction and external-action doctrine impacts. If investigation
shows a materially broader catalog redesign is required, report that concrete exception.

The activation checkpoint records this existing failure rather than claiming global validation
passes. It must introduce no new Plan error. The implementer's first EP-01 task resolves the
prerequisite and demonstrates a passing configured source catalog before proceeding with the
remaining workflow implementation. This explicitly approved bootstrap exception is bounded to
the observed baseline and does not alter the reusable activation rule.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Coherent generic routine, approval and coordination guidance is materialized correctly. | M | Plan activation through the existing route. |
| EP-02 | Validated release is available and an authorized consumer can adopt it. | M | EP-01; current release/platform readiness and covered publication/adoption authority. |

## EP-01 — Deliver the reusable workflow

**A) Epic ID, Title, And Outcome:** EP-01. A fresh consumer can commission whole-plan execution
once, delegate epic review, receive simple report-back and follow clear actual plan/Git state;
routine changes retain their proportionate route.

**B) Scope:** Managed instructions, route declarations, matching packaged resources and existing
resource tests, plus the explicitly authorized catalog prerequisite in 2.4. No other CLI change,
schema, permission engine or consumer-owned record migration.

**C) Files Touched:** The paths in 2.1, excluding release-only version/asset changes.

**D) Acceptance Criteria And Size:** M. All installed entry links resolve; the routine route carries
the eligibility and evidence requirements in 1.2; generic and Codex guidance are distinct; covered
approval is not requested again by conflicting producer instructions; a smooth multi-epic plan
continues through planned Git work and delegated review without user reauthorization;
source/mirror bytes agree. Notification perfection and an automated coordination system are not
acceptance requirements.

**E) Dependencies And Critical-Path Notes:** Activate this plan using the existing planning route.
Publish its bounded plan checkpoint under that route before source implementation. Protect existing
work and recheck current source changes against the planning base. Rebase only within this branch.

**F) Tasks Checklist**

- [x] Resolve the catalog prerequisite in 2.4 with focused parser/projection evidence, preserve frozen register bytes and genuine duplicate/identity guards, and require a passing configured catalog before the remaining source tasks.
- [x] Add `handle-routine-change.md` with eligibility, existing-record evidence, escalation triggers, authority reuse and proportionate validation from 1.2.
- [x] Add `agent-task-coordination.md` with whole-plan assignments, explicit report-back instructions, director review, readable communication, naming, handover, descendant checks and consumer model policy from 1.2.
- [x] Add `codex-task-operations.md` with supported direct task reporting, an explicitly authorized whole-plan goal example, goal-activation verification and current official task/archival references; preserve app permissions.
- [x] Align `operation-routing-and-dispatch.md`, `runbook-authoring-standard.md` and the draft/execution/review planning runbooks with whole-plan authorization, actual epic sequencing, planned commits/pushes/PRs, delegated review and material-exception handling; correct the formal-plan H1 instruction.
- [x] Add the new routes to the two component manifests, their managed READMEs and `AGENTS.managed-block.md` with short owner pointers.
- [x] Synchronize the declared changed resources under `src/codeheart_operating_kit/resources/` and extend `tests/test_packaging_resources.py` to check their materialization and exact source equality.
- [x] Run `python -m pytest -q tests/test_packaging_resources.py tests/test_routing.py` in the repository's ready Python environment and inspect the materialized routine and coordination routes.
- [x] Run `scripts/validate-markdown-headers.py` and `scripts/validate-public-core.py` against the changed Markdown/resource files and review the final diff for conflicting approval instructions.

**G) Implementation Notes:** Do not assert exact prose in new tests. Existing tests that enforce
materialization and declarations remain relevant. Tooling gaps follow the tooling-readiness route.
A workflow may retain a specific mandatory confirmation; name that requirement and recognize
already sufficient authorization to the extent its actual contract allows.

**H) Open Questions:** None for the source scope. Any demonstrated permission-code need is a scope
amendment, not an excuse to weaken a binding contract through prose.

## EP-02 — Release and consumer handoff

**A) Epic ID, Title, And Outcome:** EP-02. A versioned Kit release carries the guidance, and the
adoption handoff distinguishes tested behavior from pending real-use feedback.

**B) Scope:** Current Kit release process, installed compatibility evidence and bounded adoption.
**C) Files Touched:** `release-notes.md`, version/release files selected by the existing release
runbook, and this plan's execution evidence. No installer redesign.
**D) Acceptance Criteria And Size:** M. Required release gates pass, the published commit/assets
match the validated candidate, and an authorized consumer adoption confirms discoverable guidance
with consumer-owned content preserved. Technical delivery can complete with usability pilot pending.
**E) Dependencies And Critical-Path Notes:** EP-01 and OQ-01. Generic Kit source work does not wait
for a module change. Coordinate the single walkthrough with the downstream delivery to avoid repeats.

**F) Tasks Checklist**

- [ ] Record explicit review of the whole-plan external-action doctrine and its impact classification from 1.3 in the existing execution evidence, with workflow and adoption notes in `release-notes.md`.
- [ ] Prepare the exact version and release candidate using `docs/repo/runbooks/release-operating-kit.md`, retaining its required Go, Python, schema, installer and reproducibility evidence.
- [ ] Perform the shared fresh-context walkthrough using a neutral authorized two-epic plan: proceed through implementation, testing, agreed commits/pushes/PR work and delegated epic review without extra user approval; check readable state and explicit result-report instructions.
- [ ] Contrast a material scope exception, separate routine repair, blocked report-back and active-descendant archival in that same walkthrough; retain a module-owner routing example without live provider work.
- [ ] Publish the validated version through `release-operating-kit.md` under applicable explicit publication authority and record exact commit, asset digests, URLs and platform/signing evidence.
- [ ] Apply the published Kit through the installed `maintain-operating-kit-installation.md` route to the specifically authorized pilot consumer; verify installed route discovery and preservation of its local guidance and records.
- [ ] Record source completion, release availability, consumer adoption and pilot-pending status separately in the execution evidence, with a concise handoff to the consumer's coordination owner.

**G) Implementation Notes:** Resolve this Plan's intended integration, release and adoption
authority together at execution commissioning. Routine Git work should follow the declared
whole-plan boundary. Specific release/install runbooks may retain binding action-time inputs;
prepare those effects fully and identify actual uncovered decisions rather than introducing
per-epic user approvals. OQ-01 publication/adoption gates are checked against current authority
before asking. An unapproved final effect waits with the concrete candidate ready for
review. Do not waive release gates because the code change is small. Do not repeat a clean walkthrough
without a new material change. Actual user experience is assessed after delivery, with concrete
examples routed to the source owner; no fixed pilot duration, quota or telemetry requirement.

**H) Open Questions:** OQ-01 and OQ-02. Missing release/consumer authority blocks that effect, not
already authorized source work. Release readiness alone does not satisfy the epic's availability claim.

# Section 4 - Future Planning

Defer broader role/task schemas, assignment registries, automatic model selection, automated UI
cleanup, continuous monitoring, cron/heartbeat scheduling and general migration machinery. Capture a specific recurring problem at its owner before
expanding this scope. Guidance about supported module records belongs to that module, not the Kit.

# Revision Notes

- 2026-09-05: Drafted the bounded reusable workflow, approval recognition and coordination delivery
  from accepted requirements. Source, release, adoption and real-use evidence remain distinct.
- 2026-09-05: Revised the accepted operating direction to whole-plan execution under one
  commissioning authorization, delegated epic review, predictable Git effects, simple instructed
  report-back and optional whole-plan goals. Communication must match actual plan structure.
  The Plan remains a draft; no product implementation, task dispatch or Git publication was activated.

- 2026-09-05: Activated under the user-approved whole-delivery authority in 2.4, including delegated
  acceptance, routine Git work, source merge, validated release and named-consumer adoption.
  Implementation and actual delivery evidence are tracked in the execution log.
