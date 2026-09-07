Last updated: 2026-09-07T22:33:32Z (UTC)
Created: 2026-09-07
Status: active
Execution log: docs/repo/plans/development-feedback-and-validation/development-feedback-and-validation_execution_log.md

# Development Feedback and Proportionate Validation

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.development-feedback-and-validation
  kind: implementation
  purpose: Deliver coherent discovery, planning and execution guidance with useful review, proportionate validation, minimal machinery and verified consumer adoption.
  first_cataloged: 2026-09-07T21:34:51Z
  catalog_metadata_updated: 2026-09-07T22:33:32Z
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.proportionate-agent-workflows
```
<!-- END CODEHEART PLAN METADATA -->

## Overview

One bounded delivery revises the three core planning workflows together, aligns directly related
instructions and enforcing checks, and ships the result through the normal Kit lifecycle. The
accepted commissioning discovery is retained by the coordinating owner; this public document
contains the complete reusable implementation scope. Private strategy, consumer incident details,
account information and task assignments stay in the coordination home.

The reviewed draft was merged through PR #16. The user subsequently authorized activation and
commissioned one dedicated implementer for the complete guidance, validation, release and verified
adoption outcome on 2026-09-07. This plan is active under the commissioning boundary below.
Consumer product implementations remain paused; only the scoped Kit delivery/adoption is authorized.

Independent planning review completed on 2026-09-07. One material sequencing finding was corrected:
prepublication candidate validation is separate from explicit post-publication smoke. Focused
follow-up closed the finding; no material planning findings remain. This review does not substitute
for implementation, workflow execution, release or adoption evidence.
The subsequently approved Git-checkpoint clarification in Section 2.3 is included in the planned
combined source review. The later whole-plan authorization activates this plan explicitly.

## Essential Context Reference Files

| Source path | Purpose |
| --- | --- |
| `AGENTS.md` | Producer ownership, public-core safety and change authority. |
| `docs/repo/runbooks/change-operating-kit.md` | Change classification and proportionate maintainer checks. |
| `docs/repo/runbooks/release-operating-kit.md` | Current packaging, native platform, signing/audience and publication gates. |
| `docs/repo/reference/consumer-impact-classification.md` | Distinguish instruction, validation and actual safety-policy changes. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | Preserve intention elicitation while making inquiry/review proportional. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Preserve executable technical planning without freezing ordinary implementation choices. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Batch changes, validation timing, review corrections and interruption recovery. |
| `components/planning-workflows/managed/reference/planning-document-lifecycle.md` | Separate document lifecycle from Git publication and execution authority. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Audience, intention, operational detail and authority reuse. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Owning routes and action boundaries; retain actual controls. |
| `.github/workflows/validate.yml` | Current unconditional native, compatibility and release validation. |
| `tests/test_routing.py`, `tests/test_packaging_resources.py` | Relevant route contracts and installed resource equality. |

## Table Of Contents

- [Foundation](#section-1---foundation)
- [Strategy](#section-2---strategy)
- [Execution Plan](#section-3---execution-plan)
- [Future Planning](#section-4---future-planning)
- [Revision Notes](#revision-notes)

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

An agent using the released Kit can elicit the intended outcome, plan a substantive technical
solution, accumulate a coherent implementation batch, receive meaningful independent review and
validate at the right scope and time. Routine corrections and interruptions retain useful work.
Unnecessary mechanisms and checks are removed rather than replaced by another process framework.

Completion requires consistent source guidance and enforcing checks, a coherent independent
review, representative scenarios, appropriate package/install evidence and verified active
instruction routes in the consumer repositories/worktrees designated for subsequent plan revision.
Published artifacts alone do not complete adoption. Consumer product changes and measured savings
remain later work; this plan does not claim their completion.

## 1.2 Project And Problem Context

The approved capability covers discovery, drafting and execution together; planning review,
routine changes, task coordination, relevant templates/tests/triggers, recipe/script promotion,
and the release/adoption handoff. Preserve user authority, substantive reasoning, useful context
and actual safety boundaries. Do not insert private organization-stage assumptions into generic
Kit rules: scale obligations to the consumer's stated commitments and risks.

The approved addition covers predictable local commits, authorized branch publication and PR
integration for planning and implementation work. A finished draft need not remain uncommitted;
publishing it, including on the default branch, does not grant execution authority.

Current support and costly-to-reverse product choices are established in feature discovery. The
Kit explains how to elicit and decide them; it does not prescribe those product answers. This
workstream neither revises another product's support promises nor changes cloud/authentication
policy, caches, runtime identity or acceptance contracts owned elsewhere.

## 1.3 Current State Analysis

Planning base: `6875535fefd137fdc41321128fb82ab65bb166d2`, Kit source version `0.1.31`.
The completed proportionate-agent-workflows delivery already supplies routine routing, approval
reuse and whole-plan coordination. Extend those interfaces without reopening that completed plan.

The execution runbook asks for the smallest validation but repeats material-finding review with a
fresh reviewer. The drafting runbook combines valuable capability coverage with mandatory detail
and checklist vocabulary that can overconstrain ordinary refinement. Discovery already contains
an intention contract and intention ladder; these are strengths to preserve.

Current authoring guidance limits implied publication authority, but does not clearly establish
when to save a local commit. The lifecycle reference already separates status from Git state;
the authoring, activation and execution routes need concrete consistent examples of that rule.

`validate.yml` runs on every push and PR. Native jobs run `go test ./...` and then overlapping Go
package subsets; the oldest-Git job rebuilds its tool on every invocation. Existing routing and
resource tests include prose/structure assertions. These are concrete targets for scoped change,
not evidence that every check is waste. Repository `main` had no protection/ruleset requirements
at planning time; inspect current required-check configuration before modifying trigger behavior.
The existing manual dispatch also runs `macos-public-release` and `windows-public-release`,
which download a published tag and default to the checkout version. Prepublication candidate
validation must not inherit those guards: a newly selected version has no public assets yet.

Managed source is under `components/`; the legacy Python resource mirror is under
`src/codeheart_operating_kit/resources/`. Source and delivered bytes must remain consistent.
No new CLI, schema family, generated evidence graph, authentication helper, context-ingestion
restriction or generic process-testing framework is required.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

```text
components/planning-workflows/managed/
  runbooks/discovery-workflow.md                 # modify: intention and proportionate inquiry
  runbooks/draft-implementation-plan.md          # modify: useful plans and adaptable detail
  runbooks/execute-implementation-plan.md        # modify: batching, validation, review, recovery
  runbooks/review-planning-document.md           # modify: phase-specific material review
  runbooks/handle-routine-change.md              # modify: coherent routine feedback
  reference/planning-document-lifecycle.md       # clarify published drafts and activation
  README.md                                     # align directly affected navigation
components/agent-interface/managed/
  reference/agent-task-coordination.md           # align checkpoint and correction ownership
  reference/runbook-authoring-standard.md        # align proportionate operational detail
  reference/operation-routing-and-dispatch.md    # align affected planning routes
  reference/operational-recipe-maturity.md       # align necessity before promotion
  reference/runbook-to-script-promotion-standard.md # align smallest adequate mechanism
  runbooks/handle-tooling-readiness.md           # align authority reuse and service distinction
src/codeheart_operating_kit/resources/           # update matching declared source mirrors
.github/workflows/validate.yml                   # separate ordinary feedback and full candidates
 tests/test_routing.py                          # update meaningful affected contracts
 tests/test_packaging_resources.py              # update resource and workflow assertions
 docs/repo/runbooks/change-operating-kit.md      # align scoped validation procedure
 docs/repo/runbooks/release-operating-kit.md     # align candidate timing; retain release safeguards
release-notes.md                                # document shipped behavior and adoption
manifest.yaml, pyproject.toml, components/*/component.yaml # normal release identity updates
```

The tree is a checked starting scope, not a demand to change every listed file. Follow direct
references to identify a contradictory template/assertion, record the reason and change it with
its owner guidance. Do not audit unrelated catalog migration, module contracts or all historical
plans. Existing authorization, identity and source-integrity predicates remain binding.

All six planning runbooks are agent-facing with concise user dialogue where inquiry needs it.
Keep compact intention blocks and executable action/stop boundaries; scale examples and detail to
ambiguity. Tooling readiness is hybrid; change/release procedures are maintainer-facing. References
remain short owner guidance rather than new operational recipes. Core workflows stay L1, reviewed
through realistic executability scenarios. No recipe promotion or reusable script is selected.
Apply the recipe-maturity and script-promotion references to strengthen the necessity test, not to
force every operation up a maturity ladder.

Missing Git, Python, Go or a package environment remains local tooling readiness. Authentication,
billing and live provider permissions remain service-owned preflight. Concrete module installation
commands/version facts stay module-owned. Reuse existing sufficient local-setup authority and
installed tools; never infer credentials or introduce an unrelated service login prerequisite.

## 2.2 Open Questions And Assumptions Requiring Clarification

| ID | Blocker and affected epics | Resolution/default |
| --- | --- | --- |
| OQ-01 | BLOCKER: no. Affects: EP-03. | At release preparation, select the next unused patch version from current source/releases and use the existing approved distribution audience. Retain an unsigned prototype candidate until actual signing/audience and publication authority are satisfied. No support expansion is assumed. |
| OQ-02 | BLOCKER: no. Affects: EP-02. | Re-read required check names before trigger edits. Default to preserving a reliable ordinary PR result and moving expensive coverage to an explicit candidate invocation. Required checks cannot silently become skipped/pending; resolve a conflicting enforced setting with the owner before publishing that workflow. |
| OQ-03 | BLOCKER: no. Affects: EP-03. | The commissioning owner supplies exact consumer roots/worktrees and covers adoption in the execution assignment. Use the existing lifecycle route, preview changes and verify installed version/routes at each actual plan-amendment location. An inaccessible consumer stays adoption-pending; no global rollout is implied. |

These are normal action-time inputs with bounded handling, not unanswered architecture choices.
The full-delivery authorization now includes the final effects listed below. Current action-time
checks, actual target bindings and the stated exceptions still apply.

## 2.3 Architectural Decisions With Reasoning

| Decision | Chosen approach and rationale | Alternatives and evolution |
| --- | --- | --- |
| One coherent workflow revision | Revise discovery, drafting and execution in one workstream, retaining existing routes. A caller must receive compatible rules across phases. | Three independent discoveries add coordination and can diverge. Split later only for independently owned capability. |
| Intention before method | Preserve active intention elicitation, trade-offs and meaningful feasibility checks. Confirm material inferences; an explicitly requested method remains a constraint. | Requirement transcription alone loses the desired outcome. Do not replace useful inquiry with a shorter template. |
| Substantive adaptable plans | Record outcomes, scope, approach, consequential constraints, dependencies and validation. Explain requirements versus chosen approach versus assumptions. Let in-scope technical refinement proceed with brief rationale. | Do not prescribe thin plans, exhaustive frozen commands or a new mandatory decision form. Add detail when actual uncertainty warrants it. |
| Purposeful review | Use one independent review at the planned meaningful checkpoint, with the same reviewer following corrections. Escalate review scope for material design/impact change, inadequate independence, a reviewer limitation or unresolved concern. | A new reviewer after every finding repeats work; no review misses real defects. Adjust checkpoints as delivery risk changes. Directors own acceptance without automatically duplicating the technical review. |
| Feedback and candidate acceptance | Cheap affected checks support iteration; broader coverage runs on a coherent candidate. Ordinary progress, documentation and resume do not automatically recreate acceptance. Reuse evidence only while its relevant inputs and environment remain applicable. | Running everything after every edit burns time; universally skipping coverage hides regressions. Unknown dependencies justify investigation/wider checks. No generic proof cache is introduced. |
| Remove before optimizing | Ask what failure a mechanism prevents and whether a simpler existing tool/record suffices. Remove unjustified checks, duplicate artifacts and promotion requirements. Keep justified integrity boundaries. | Downgrading all noise to warnings retains cost. A new evidence framework recreates it. Reconsider complexity when actual consumers/risks justify it. |
| Ordinary recovery and access | Preserve valid output, classify code/test/environment/auth/billing failures, continue safe in-scope repairs, and recheck affected evidence. Establish effect truth before retrying an uncertain external write. Use suitable already-authorized connectors and normal machine sessions. | No universal recovery/auth service, credential copying or session-policy change. Service-specific restrictions remain owner decisions. |
| Predictable Git checkpoints | Save coherent work through local commits, publish within the agreed delivery authority, and merge accepted shared records through the normal PR route. Keep lifecycle and execution approval independent from Git state. | Leaving finished drafts uncommitted makes recovery and handoff fragile; treating every commit as an activation confuses authority. Use existing Git and plan metadata, not a new checkpoint service or approval form. |

### Git checkpoints and planning lifecycle

Implement the following reusable defaults in the existing planning routes and lifecycle reference:

- **Local work and commits:** Keep in-progress edits local while forming a coherent change. At a
  useful completed checkpoint, commit the task's own changes after proportionate inspection and
  checks as part of the authorized authoring or delivery work. An explicit user preference,
  overlapping changes or a concrete blocker can justify leaving work uncommitted; state the
  reason. A local commit does not require another independent review or a new approval prompt.
- **Branch publication:** Push the working branch and create/update its PR at useful checkpoints
  when the user request or delivery assignment covers those effects. Establish the publication
  finish line once and reuse that authority. Drafting alone does not imply external publication;
  a publication-only request can cover a draft without approving its execution.
- **Default-branch integration:** Merge through the repository's normal PR route when the change
  is accepted as shared repository content, the relevant checks pass and merge authority is
  covered. Do not infer direct pushes to the default branch, bypass protection or include
  unrelated files. Acceptance of a documented proposal is distinct from execution approval.
- **Independent lifecycle:** A reviewed, committed, pushed or merged implementation plan can
  remain `Status: draft`. Set it active only under execution/activation authority; publication
  must not dispatch implementation or imply release/adoption approval. Preserve the existing
  bounded activation checkpoint: activation covers its plan-only commit/push, not PR/merge or
  broader effects. An authorized implementer can start on its active pushed branch without
  waiting for a main merge where the owner rules allow it.

Apply the rule to discovery, implementation planning, routine work and whole-plan execution.
Keep checkpoint depth proportionate: no new per-commit test suite, mandatory reviewer, evidence
artifact or approval stage. Explain uncovered publication authority only after preparing the
concrete change; do not ask again for effects already included in the assignment.

### Commissioning and delivery boundary

**Activated 2026-09-07.** The user explicitly requested merging the two documentation PRs, then
activating and commissioning one dedicated Kit task through guidance, proportionate checks,
release and verified adoption. The Organizational Operating Model Improvement Program Director is
the acceptance owner. One implementer owns EP-01 through EP-03 on
`codex/development-feedback-delivery`, starting from the merged planning baseline
`9e73c654aed7de161f0414e37701f8426b1cd44b` and its pushed activation checkpoint.

The grant covers required source changes, normal local development setup, proportionate local and
hosted validation, coherent commits, normal branch pushes, one implementation PR with updates,
source merge after delegated acceptance, next-patch version selection, release preparation and
publication, public-release smoke, and supported Kit adoption at the consumer targets named in the
private assignment. It includes scoped consumer adoption commits and normal pushes/PR integration
where the owning route requires them. Preserve unrelated work and consumer-owned content; never
include product implementation or bypass repository protections, authentication, approval review
or failed acceptance. No force-push, history rewrite, destructive cleanup, credential/security
policy change, billing increase, support expansion or unrelated provider operation is granted.

The finish line is merged source, a validated versioned release under the existing approved
distribution audience and signing boundary, successful published-asset smoke, and verified active
version/routes at every assigned consumer worktree. The private commissioning record identifies
the four current working locations across the coordination home and two product repositories;
confirm their identity and current state immediately before adoption. A missing/moved target stays
explicitly adoption-pending until the acceptance owner resolves its replacement. Never broaden
into all repositories or treat an unverified installation as complete.

The Program Director accepts the combined EP-01/02 source checkpoint and the final release/adoption
handoff. The implementer may prepare independent release inputs while review is pending; after
source acceptance it performs the already-authorized integration, release and adoption steps once
their relevant gates pass. These are delegated acceptance points, not new user-approval turns per
epic. Material scope, authority, cost or unresolved safety departures return to the Director.

Send concise direct reports only to the commissioning Director at the combined source review,
completion, genuine blocker or material exception. The private task assignment supplies its exact
execution locator and authorizes those reports. Include outcome, evidence, Git/release/adoption
state and the actual decision needed. The sibling execution log owns public-safe technical
progress; private consumer paths and task identities remain in the private assignment. No goal,
watcher, automatic monitoring or report to other organizational roles is commissioned.

Apply the approved review direction to this work: EP-01 and EP-02 form one coherent source/review
checkpoint. Independent preparation may overlap; dependent completion waits for that acceptance.
The source review covers their combined effects and the scenarios below. Corrections return to
the same reviewer. EP-03 verifies delivery evidence rather than restarting an unchanged doctrine
review. A material new concern can require a focused additional review. This explicit plan scope
replaces the old automatic fresh-review loop for this delivery once execution is authorized.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Consistent usable discovery, planning, execution and supporting guidance | M | Approved execution assignment |
| EP-02 | Kit checks and CI deliver proportionate feedback and truthful candidate acceptance | M | EP-01 direction; shared source acceptance checkpoint |
| EP-03 | Reviewed change released and adopted for owner-plan revision | M | EP-01/02 accepted; explicit final-effect authority |

## EP-01 — Coherent operating guidance

**A) Epic ID, Title, And Outcome:** EP-01. Users and agents can follow compatible instructions
through inquiry, planning, delivery, review and ordinary recovery without unnecessary ceremony.

**B) Scope:** Three core runbooks, directly related planning/review/routine/coordination routes,
recipe necessity and tooling guidance; public reusable doctrine only.

**C) Files Touched:**
```text
components/planning-workflows/managed/{runbooks,reference/planning-document-lifecycle.md,README.md}
components/agent-interface/managed/{reference,runbooks/handle-tooling-readiness.md}
src/codeheart_operating_kit/resources/          # corresponding mirrors
```

**D) Acceptance Criteria And Size:** M. The guidance preserves intention elicitation and technical
planning, permits coherent batches and ordinary in-scope refinement, defines useful review and
invalidation, and separates necessary blockers, judgment, optional advice and removable noise.
No new control plane, blanket context restriction or private stage policy appears. The routes
explain when to commit, publish and merge; a published draft stays inactive until separately
authorized, and routine checkpoints do not create additional approval or review loops.

**E) Dependencies And Critical-Path Notes:** Follow the accepted capability. EP-02 preparation can
proceed alongside prose corrections; the two epics share one final source review.

**F) Tasks Checklist**

- [ ] Revise `discovery-workflow.md` and `draft-implementation-plan.md` to preserve intention, capability coverage and consequential decisions while scaling detail and removing rigid procedure that adds no decision value.
- [ ] Revise `execute-implementation-plan.md`, `review-planning-document.md` and `handle-routine-change.md` for coherent batches, phase-specific review, focused correction follow-up, applicable evidence reuse and scoped recovery.
- [ ] Align the discovery/drafting/execution/routine runbooks, `planning-document-lifecycle.md`, planning route entries and `agent-task-coordination.md` with Section 2.3's Git checkpoint defaults and independent publication/activation semantics.
- [ ] Align `agent-task-coordination.md`, the affected planning routing entries, `runbook-authoring-standard.md` and `handle-tooling-readiness.md` with those rules, including existing-authority reuse and service-specific authentication boundaries.
- [ ] Reconcile `operational-recipe-maturity.md` and `runbook-to-script-promotion-standard.md` so necessity and simpler existing mechanisms precede promotion; remove directly conflicting requirements and examples.
- [ ] Synchronize the corresponding declared resource mirrors and affected navigation entries.
- [ ] Review the combined EP-01/02 candidate against the representative scenarios in EP-02 and resolve material findings through focused reviewer follow-up.

**G) Implementation Notes:** Keep instructions concrete at consequential boundaries. Do not turn
all prose into a schema, replace every removed requirement with a warning, or let representative
examples become universal support obligations. Read existing local/module instructions when
relevant; this plan does not restrict context ingestion.

**H) Open Questions:** None blocking; EP-02 owns trigger/enforcement alignment.

## EP-02 — Feedback, enforcement and coherent acceptance

**A) Epic ID, Title, And Outcome:** EP-02. Kit development gets useful quick feedback while
meaningful candidates retain their required coverage and release integrity.

**B) Scope:** Kit-owned validation triggers, overlapping invocations, directly affected assertions
and maintainer instructions. No changes to another repository's tests, security policy or CI.

**C) Files Touched:**
```text
.github/workflows/validate.yml
 tests/{test_routing.py,test_packaging_resources.py}
 docs/repo/runbooks/{change-operating-kit.md,release-operating-kit.md}
```

**D) Acceptance Criteria And Size:** M. Small documentation/plan changes do not launch full native
installation, history/compatibility or performance coverage. Ordinary PR feedback is reliable;
full candidate coverage is explicit, traceable and available before release. Overlapping commands
are removed without losing distinct behavior. No required check is evaded. Guidance and tests
agree; a metadata/log edit alone does not invalidate unrelated product evidence.

**E) Dependencies And Critical-Path Notes:** EP-01 defines the behavior. Match candidate inputs to
checks actually covering them. Amend procedure and enforcement together before relying on the
changed requirement. Existing release integrity and supported platform obligations remain.

**F) Tasks Checklist**

- [ ] Inspect current required-check configuration and map existing `validate.yml` jobs to ordinary feedback and full candidate/release acceptance in the execution log.
- [ ] Refactor `validate.yml` to run a small ordinary PR/main feedback lane and retain expensive native/install/compatibility coverage for an explicit coherent-candidate dispatch; prevent duplicate branch-push and PR executions for the same ordinary change.
- [ ] Separate dispatch modes in `validate.yml`: candidate validation uses the checked-out source and staged artifacts; published-release smoke requires an explicit released tag, downloads that release, and runs only the public-asset smoke jobs.
- [ ] Remove overlapping Go/package invocations from each lane where broader coverage already proves the same behavior, preserving separately configured regression and platform cases.
- [ ] Update `change-operating-kit.md` and `release-operating-kit.md` to document the actual lane commands, candidate timing, correction/reuse rules and retained release safeguards.
- [ ] Update directly affected routing/resource/workflow tests to assert meaningful behavior and delivered resource equality; remove wording-only assertions for superseded ceremony without weakening identity, authority and preservation contracts.
- [ ] Run `python -m pytest -q tests/test_routing.py tests/test_packaging_resources.py`, public-core and Markdown checks once on the coherent source candidate; run additional targeted tests for any actual changed executable behavior.
- [ ] Complete one independent source review with representative fresh-context walkthroughs and record findings, corrections and evidence limits in the execution log.

**G) Implementation Notes:** Prefer existing workflow event/branch/path controls and a simple
candidate invocation. Do not build a generic dependency classifier, proof database or orchestration
service. Do not use commit-message bypasses as the new operating model. Recheck current check
configuration before publication; a conflicting protection requires owner resolution.

Walk through five cases in the same review: a routine fix; a larger batch with deferred broad
performance/native validation; a review correction; a plan/log-only update; and an interruption
with one invalidated check and otherwise valid evidence. Include a vague nested module request
that must discover the owner and distinguish missing tooling from service login, and a user's
proposed method that hides an underlying intention. Pass means correct routing, justified check
scope/timing, truthful acceptance and preserved authority. One reviewer can cover routing and
execution cases; no permanent scenario harness or repeated independent review is required.

Extend the plan/log case with a reviewed draft that receives a local commit, then an authorized
publication and merge while remaining draft. Verify that authoring alone does not authorize its
external publication, that no implementation starts from the merge, and that later activation
does not require another main merge or expand its bounded Git authority. Check the directly
affected route tests and resource mirrors; do not add a new publication framework or full suite.

Workflow validation must also demonstrate that unknown/changed executable surfaces cannot be
silently accepted by the cheap lane as release-qualified. Full required release evidence remains
at the candidate boundary. Verify that candidate mode never downloads the unpublished version,
released-smoke mode rejects a missing tag, and post-publication smoke does not restart the full
source suite. Default a manual invocation to candidate mode; never infer public-release intent
from a version string. Cancel only superseded ordinary validation; do not cancel an effectful
release/install operation blindly. Record expensive checks once and rerun only invalidated evidence.

**H) Open Questions:** OQ-02.

## EP-03 — Release, adoption and owner handoff

**A) Epic ID, Title, And Outcome:** EP-03. Consumers preparing subsequent plan revisions actually
use the reviewed guidance from a supported Kit release.

**B) Scope:** Normal release identities/assets and the expressly assigned consumer lifecycle
operations. Technical delivery and observed future efficiency remain distinct facts.

**C) Files Touched:**
```text
release-notes.md, manifest.yaml, pyproject.toml
components/*/component.yaml, src/codeheart_operating_kit/resources/ # normal release identity
 docs/repo/plans/development-feedback-and-validation/*_execution_log.md
<assigned consumers>/.codeheart/kit*           # lifecycle command outputs only
<assigned consumers>/AGENTS.md                # supported managed-block refresh only
```

**D) Acceptance Criteria And Size:** M. Versioned packs contain the reviewed source; applicable
native/install/upgrade and release checks pass on the final candidate. Each assigned consumer
retains authored work and has the intended version and active route contents. The owner receives
a concise handoff enabling later Director onboarding and separately governed plan amendments.

**E) Dependencies And Critical-Path Notes:** Requires EP-01/02 acceptance and explicit assignment of
integration/release/adoption effects. Select the current version and exact targets at action time.
Do not start a second full acceptance run merely because the task resumed.

**F) Tasks Checklist**

- [ ] Select the next unused patch version and prepare `release-notes.md` plus normal version/resource updates, classifying instruction and validation changes and documenting consumer adoption.
- [ ] Run the agreed full candidate lane once and `release-operating-kit.md` packaging gates, including reproducible macOS universal/Windows x64 artifacts and supported native install/upgrade preservation evidence.
- [ ] Integrate and publish the exact validated candidate through normal Git/release routes under the commissioning authority; retain target commit, version, assets and actual acceptance evidence.
- [ ] Run the published-release smoke mode against the explicit released tag, verify public asset retrieval and supported native install behavior, and retain source validation from the unchanged accepted candidate.
- [ ] Preview and apply the released Kit through `maintain-operating-kit-installation.md` in each assigned consumer; verify installed version, changed route contents and preservation of authored work at the actual owner-plan amendment worktrees.
- [ ] Give the acceptance owner the release/adoption results and the changed operating expectations in the existing execution record, explicitly identifying any target still pending and the later Director-onboarding boundary.

**G) Implementation Notes:** This managed-instruction delivery requires a Kit release/adoption,
not Foundry or provider module releases. Preserve supported release checks without replaying
consumer implementations. Existing catalog/archive identity verification is retained. Do not
hand-edit installed snapshots or update metadata; use the supported lifecycle commands. Apply
only the designated consumer boundary; source self-installation is not required for producer
planning, whose authority remains tracked source.

**H) Open Questions:** OQ-01, OQ-03.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

After verified process adoption, the coordinating owner onboards Directors and amends each
existing consumer/product implementation plan separately. That later work selects concrete test,
CI, cache, support, runtime, proof and authentication changes with the relevant owners. It is not
another epic here and may not be silently started from this plan.

Context-ingestion optimization, new monitoring/benchmark frameworks, organization record/profile
changes, platform-support expansion and an overarching implementation plan are excluded.

## 4.2 Future Considerations

Assess practical improvement using existing run/implementation records: unnecessary work removed,
feedback delay and avoidable interruption. Do not promise a measured speedup before observing it.
Revisit assurance/support obligations when real adoption commitments, sensitive operations or
changed risks justify their cost. Preserve effective safeguards while reconsidering machinery.

# Revision Notes

- 2026-09-07: Drafted from the approved shared-process capability and current producer source.
  The draft proposes one coherent source review and a release/adoption finish line; execution
  and final effects require commissioning. HQ strategy remains in its existing private records.
- 2026-09-07: Independent review identified and closed the candidate/public-release dispatch
  sequencing gap. The draft is ready for commissioning with the declared action-time inputs.
- 2026-09-07: Added the user-approved local commit, authorized publication and PR-merge defaults
  to EP-01, with lifecycle clarification and focused EP-02 acceptance coverage. The Operating Kit
  owns this reusable rule; the plan remains draft and no installed guidance changes here.

- 2026-09-07: Activated after the user commissioned complete delivery. The grant now fixes the
  source branch, acceptance owner, integration/release/adoption authority, report-back and finish
  line. Exact private consumer roots are in the commissioning assignment; product work stays held.
