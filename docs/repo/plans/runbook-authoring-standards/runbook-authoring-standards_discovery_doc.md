Last updated: 2026-06-23T12:02:14Z (UTC)
Created: 2026-06-23
Status: draft

# Runbook Authoring Standards Discovery

## Overview

This discovery investigates whether Codeheart Operating Kit should define reusable standards for
authoring runbooks, especially when runbooks guide humans through onboarding or give agents a
precise execution path.

The immediate trigger came from the Foundry M365 workspace onboarding runbook. The runbook became
technically correct after one hardening pass, but its first live consumer use still produced a
too-technical, batch-question interaction. That exposed a reusable doctrine gap: runbooks need to
declare whether they are human-facing, agent-facing, or hybrid, and the required quality checks
are different for each audience.

This is public Operating Kit discovery. It must not depend on private tenant details, customer
names, credentials, or private business content. Consumer examples are generalized as evidence
patterns.

## Goals

- Discover reusable runbook authoring standards that improve future Operating Kit, repository,
  product, module, and consumer-local runbooks.
- Separate human-facing runbook requirements from agent-facing execution requirements.
- Identify useful existing Operating Kit patterns that should be extracted into a standard.
- Identify weaker existing runbooks that can serve as test applications for the new standard.
- Define where a later implementation should place managed doctrine and which planning/review
  workflows should enforce it.

## Non-Goals

- Do not rewrite existing runbooks in this discovery.
- Do not publish a new Operating Kit release from this discovery.
- Do not create a validator implementation yet.
- Do not move consumer-local runbooks into the Operating Kit.
- Do not make every runbook human-facing; many runbooks are correctly maintainer- or agent-facing.

## Sampled Evidence

Detailed public-core-safe sampling notes are recorded in
`attachments/runbook-sampling-matrix.md`.

| Sample | Audience Shape | Useful Quality Patterns | Gaps Or Test-Case Value |
| --- | --- | --- | --- |
| `components/agent-interface/managed/runbooks/conduct-first-run-onboarding.md` | Human-facing guided onboarding | Ordered procedure, visible choices, direct user-owned decisions, stop before writes, concise wording. | Strong exemplar for human-facing onboarding standards. |
| `components/agent-interface/managed/reference/onboarding-context-contract.md` | Human-facing contract plus agent notes | Separates public prompt, agent contract, ordered context, exact copy, storage boundary, non-interactive rules. | Strong exemplar for separating script from operator notes. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Agent-facing planning workflow | Fresh-implementer test, concrete checklist rules, blocker handling, validation expectations. | Strong exemplar for agent-facing runbook quality, but not a general runbook standard yet. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | Agent-facing review workflow | Review stance, severity, issue order, readiness checks, material-finding criteria. | Strong exemplar for review-gate standards. |
| `components/structure-governance/managed/runbooks/change-documentation-placement.md` | Agent-facing operational procedure | Owner selection order, procedure, archive rules, validation list. | Good agent-facing standard example; could be improved with an explicit input/output/evidence block. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer-facing procedure | Short ordered procedure and stop conditions. | Terse. Useful test case for whether all maintainer runbooks need inputs, outputs, evidence, and validation sections. |
| `docs/repo/runbooks/release-operating-kit.md` | Maintainer-facing release procedure | Strong ordered release sequence and stop conditions. | Could become a test case for explicit evidence fields and authority gates. |
| `docs/repo/runbooks/promote-consumer-change.md` | Maintainer-facing promotion procedure | Captures reusable doctrine promotion and public-safety sanitization. | Terse. Useful test case for promotion-specific inputs/evidence requirements. |
| Consumer business decision runbooks | Human/agent hybrid business operations | Clear folder conventions, canonicalization rules, decision-ready agenda fields, after-meeting updates. | Good domain-runbook pattern, but not a scripted user onboarding flow. |
| Foundry M365 `onboard-workspace.md` pre-UX-hardening | Hybrid onboarding | Good source-of-truth and safety constraints after first hardening. | Test case showing that technically correct operator rules are not enough for nontechnical onboarding. |
| Foundry M365 `operate-workspace.md` | Agent-facing technical execution | Request routing, first response rules, command sequences, approvals, stop conditions. | Good test case for agent-facing runbook standards and recipe-grade execution paths. |

## Condensed Sample Lessons

The sampling detail should stay in the attachment. The durable lessons for the discovery are:

- human-facing runbooks need a clear user-experience intent, not only safety rules;
- agent-facing runbooks need a clear execution intent, not only desired outcomes;
- hybrid runbooks must separate user copy from operator notes;
- terse runbooks are acceptable when risk and ambiguity are low;
- detail should scale with blast radius, repeatability, and how much judgment the agent must apply;
- human-facing runbooks should check for visible local preference context before asking repeated
  setup questions;
- every durable runbook should give the agent enough intent to behave correctly when the exact edge
  case is not scripted.

Related Operating Kit feedback about local user preferences and shared tooling readiness is
recorded in `attachments/related-operating-kit-feedback.md`.

## Decision Inventory

### D-001 - Classify Runbooks By Primary Audience

Status: recommended

Decision: every durable runbook should declare one of these audience shapes:

- `human-facing`;
- `agent-facing`;
- `hybrid`;
- `maintainer-facing`.

Why it matters: a human onboarding script and an agent execution recipe need different structure.
Without audience classification, a runbook can be technically complete but unusable in chat, or
pleasantly worded but too vague for execution.

Recommended default: add the classification near the top of runbooks that are created or materially
changed. Do not require an immediate mass migration of all existing runbooks.

Terminology note: modules do not need object-oriented inheritance from `human-facing` or `hybrid`.
The runbook audience classification is the trigger that tells agents and reviewers which reusable
Operating Kit quality rules apply. A module onboarding runbook that is classified as `hybrid`
inherits the generic human-facing standards through that route rather than by restating them.

Modeling note: in the first implementation, audience classification should be modeled as visible
Markdown metadata in the compact intention block. It does not need a machine-readable schema or
validator until the standard proves stable.

### D-002 - Define Human-Facing Runbook Requirements

Status: recommended

Decision: human-facing runbooks should include a `User-Facing Flow` or equivalent section with:

- plain-language goal and outcome before technical detail;
- exact default wording for key turns;
- one user-owned decision per turn by default;
- no more than two questions in a single user-facing turn;
- visible choices when the option set is small;
- recommended default labeled first;
- internal mechanics kept out of first-turn copy;
- explicit approval-gate wording before writes, sign-ins, installs, external changes, or sensitive
  reads;
- help text for values a nontechnical user may not know how to find;
- a preference check before asking repeated setup questions when an agent-visible local preference
  file exists.

For language specifically, the first standard should stay narrow: check whether an agent-visible
`.codeheart/user/preferences.yaml` exists and contains `language`. If it exists and is readable,
continue in that language. If it is absent, unreadable, or does not contain language, ask the user
once and continue in the answered language for the current flow. Do not build a broader preference
system in the runbook-authoring implementation.

Why it matters: this turns onboarding from a technical questionnaire into a guided experience.

Recommended default: implement first for onboarding runbooks and other recurring user-guided
flows.

### D-003 - Define Agent-Facing Runbook Requirements

Status: recommended

Decision: agent-facing runbooks should include enough execution detail that a fresh agent does not
need to invent the workflow. Required elements should include:

- goal and non-goal;
- source of truth;
- required inputs and accepted formats;
- preconditions and tool readiness checks;
- ordered execution path;
- command, API, portal, or document-surface lane where applicable;
- stop conditions;
- approval gates;
- evidence and run-record requirements;
- validation;
- rollback, retain, offboarding, or cleanup when relevant.

Why it matters: agent-facing runbooks are operational contracts, not policy summaries. Weak
agent-facing runbooks cause inconsistent execution and review churn.

Recommended default: reuse the "fresh implementer test" from implementation planning as the
agent-facing runbook test.

### D-004 - Require Hybrid Runbooks To Separate User Copy From Operator Notes

Status: recommended

Decision: hybrid runbooks should separate:

```text
## User-Facing Flow
## Operator Notes
## Execution Path
## Stop Conditions
## Evidence And Validation
```

Why it matters: mixing internal warnings, file-boundary rules, and technical state into user copy
creates nontechnical onboarding failures.

Recommended default: apply this structure to onboarding, setup, consent, support, and
externally-triggered operational runbooks.

### D-005 - Add The Standard To Planning, Execution, And Review Workflows

Status: recommended

Decision: the runbook authoring standard should not live only as a reference. Later implementation
should update:

- implementation planning, so plans that create or change runbooks include audience classification
  and standard-specific acceptance criteria;
- implementation execution, so review gates check the standard before marking a runbook epic
  complete;
- planning/document review, so reviewers can flag vague agent-facing runbooks or user-facing
  runbooks that expose internal mechanics.

Why it matters: a standard that is not checked during planning, execution, and review will not
reliably change behavior.

Recommended default: update managed planning workflow docs after the reference exists.

### D-006 - Use Existing Runbooks As Standard Test Fixtures

Status: recommended

Decision: do not only write doctrine. Use representative runbooks as test applications:

- strong examples to preserve: first-run onboarding, onboarding context contract, review planning,
  implementation planning;
- terse maintainer examples to evaluate: change Operating Kit, release Operating Kit, promote
  consumer change;
- hybrid UX test case: Foundry M365 onboarding;
- agent-facing recipe test case: Foundry M365 operation runbook.

Why it matters: checking the standard against real runbooks will reveal whether the rule is
practical or too bureaucratic.

Recommended default: use `attachments/runbook-sampling-matrix.md` as the initial audit matrix and
review fixture. Do not actively retrofit consumer or module runbooks in the first implementation.
Consumer repositories should receive the new managed Operating Kit guidance through normal kit
update/sync, and future or materially changed runbooks should follow the updated instructions.

### D-007 - Add A Compact Runbook Intention Block

Status: recommended

Decision: every new or materially changed durable runbook should include a compact intention block
near the top. The block should be short enough to avoid bureaucracy, but explicit enough to guide
agent judgment when the runbook does not cover an edge case.

Recommended shape:

```text
Audience: human-facing | agent-facing | hybrid | maintainer-facing

Intent:
<what this runbook is trying to achieve and what good behavior looks like>

Success:
<observable successful outcome>

Agent judgment boundary:
<what the agent may adapt, and what it must not invent or bypass>

Stop boundary:
<when the agent must stop and ask before continuing>
```

For human-facing runbooks, the intent should describe the desired conversation behavior, such as
guiding a nontechnical user step by step, avoiding internal mechanics in first-turn copy, and
asking one decision at a time.

For agent-facing runbooks, the intent should describe the execution behavior, such as preferring
the named tool lane, inspecting before writes, recording evidence, validating the result, and
stopping rather than inventing high-impact actions.

Why it matters: concrete scripts and command recipes cannot cover every edge case. A short intent
block gives the agent a stable basis for reasonable adaptation without replacing explicit
execution paths where risk or ambiguity requires them.

Recommended default: require the intention block for new or materially changed runbooks. Do not
mass-retrofit existing runbooks before the first standard implementation proves the shape.

## Requirements And Evaluation Criteria

### FR-001 - Audience Classification

Every new or materially changed runbook must identify its primary audience shape.

### FR-002 - Human-Facing Script Standard

Human-facing runbooks must provide exact or near-exact user-facing wording for critical turns and
must avoid exposing internal mechanics in first-turn copy.

### FR-003 - Agent-Facing Execution Standard

Agent-facing runbooks must include enough ordered, concrete execution detail for a fresh agent to
act without inventing the workflow.

### FR-004 - Hybrid Separation

Hybrid runbooks must separate user-facing copy from operator-only notes.

### FR-005 - Planning Workflow Integration

Implementation planning, execution, and review workflows must reference the runbook authoring
standard when runbooks are created or materially changed.

Implementation planning should require runbook-related plans to state:

- which runbooks are created or materially changed;
- audience classification for each affected runbook;
- whether each runbook needs human-facing flow, agent-facing execution path, hybrid separation, or
  maintainer authority/evidence handling;
- whether the compact intention block is present;
- whether the plan is intentionally not retrofitting existing consumer/module runbooks.

Implementation execution should require executors to:

- follow the planned runbook-authoring scope;
- preserve unrelated runbooks and consumer-owned module docs unless the plan explicitly changes
  them;
- update new or materially changed runbooks to the standard;
- avoid marking a runbook epic complete when the audience, intention block, approval/stop
  boundaries, or required execution detail are missing.

Planning-document review should check:

- whether runbook changes are scoped explicitly;
- whether human-facing runbooks avoid internal mechanics and repeated setup questions;
- whether agent-facing runbooks are executable by a fresh agent;
- whether hybrid runbooks separate user copy from operator notes;
- whether maintainer-facing runbooks have enough authority, evidence, and validation handling for
  their blast radius;
- whether the plan accidentally turns a standard rollout into a broad retrofit.

### FR-006 - Compact Intention Block

Every new or materially changed durable runbook must include a compact intention block covering
audience, intent, success, agent judgment boundary, and stop boundary. The block guides behavior;
it does not replace concrete execution steps for risky or ambiguous workflows.

### NFR-001 - Public-Core Safety

The standard must be generic and public-safe. Examples must use placeholders or sanitized
consumer-local patterns.

### NFR-002 - Low Bureaucracy

The standard must not force every short internal maintainer checklist into a verbose template.
It should scale with audience, risk, and repeatability.

### NFR-003 - Testability

The standard must be reviewable with a checklist or audit matrix. Future automation is optional.

## Placement Options

### Option A - Agent-Interface Managed Reference

Path candidate:

`components/agent-interface/managed/reference/runbook-authoring-standard.md`

Pros:

- close to human/agent interface behavior;
- natural owner for user-facing script and operator-note separation;
- existing onboarding context contract lives nearby.

Cons:

- agent-facing technical runbooks also touch planning and structure governance.

### Option B - Structure-Governance Managed Reference

Path candidate:

`components/structure-governance/managed/reference/runbook-authoring-standard.md`

Pros:

- close to documentation artifact kinds and placement;
- naturally governs runbook shape as a documentation class.

Cons:

- human-facing conversation quality is more than document placement.

### Option C - Planning-Workflows Reference Only

Path candidate:

`components/planning-workflows/managed/reference/runbook-quality-standard.md`

Pros:

- close to implementation plan, execution, and review gates.

Cons:

- too narrow; runbooks exist outside planning workflows.

Recommended default: Option A, with links from planning-workflows and structure-governance. The
standard is fundamentally about how agents and users interact with runbooks, while planning and
structure docs can enforce and route to it.

## Open Questions

### OQ-001 - Should The Standard Be One Reference Or Split References?

BLOCKER: no

Owner: implementation planner

Recommended default: start with one reference containing sections for audience classification,
human-facing runbooks, agent-facing runbooks, hybrid runbooks, and review checklist. Split later
only if the document becomes too large.

### OQ-002 - Should Existing Managed Runbooks Be Retrofitted Immediately?

BLOCKER: no

Owner: implementation planner

Recommended default: do not mass-retrofit. Apply the standard to a small sample of high-value
runbooks and require it for new or materially changed runbooks.

### OQ-003 - Should A Validator Be Built In The First Implementation?

BLOCKER: no

Owner: implementation planner

Recommended default: no. Start with doctrine and review checks. Add a validator later only after
the standard stabilizes.

### OQ-004 - Should Consumer Repositories Receive A Managed Route To The Standard?

BLOCKER: no

Owner: implementation planner

Recommended default: yes, if the standard lands under agent-interface managed docs. Consumer
agents should be able to find it from `.codeheart/kit/docs/agent-interface/README.md` and the kit
fallback inventory.

### OQ-005 - Should Preference And Tooling Readiness Be In This Standard?

BLOCKER: no

Owner: implementation planner

Recommended default: include only the immediate runbook-authoring implication: human-facing
runbooks should check visible local language preference before asking again. Keep the broader
preference contract and environment-readiness/tooling register as related Operating Kit feedback
until they receive their own discovery or implementation plan.

## Risks

### R-001 - Over-Templating Runbooks

Likelihood: medium

Impact: medium

Mitigation: scale requirements by audience and risk. Allow short maintainer runbooks to remain
short when they are unambiguous and low-risk.

### R-002 - Human-Facing Scripts Become Too Rigid

Likelihood: medium

Impact: medium

Mitigation: require exact wording for critical onboarding/approval turns but allow domain-specific
adaptation elsewhere.

### R-003 - Standards Are Written But Not Enforced

Likelihood: high

Impact: high

Mitigation: update implementation planning, execution, and review runbooks so the standard is
checked whenever runbooks are created or materially changed.

### R-004 - Private Consumer Examples Leak Into Public Kit

Likelihood: low

Impact: high

Mitigation: use sanitized placeholders and pattern summaries only. Keep consumer-specific details
in consumer repos.

## Implementation Handoff

Next useful step: draft an implementation plan for an Operating Kit instruction-only release that:

1. Adds the managed runbook authoring standard reference, including the compact intention block
   shape.
2. Routes it from the owning component README and kit fallback inventory.
3. Updates implementation planning, execution, and review runbooks to require audience
   classification, compact intention blocks, and standard checks when runbooks are created or
   materially changed.
4. Adds a lightweight audit matrix for sampled existing runbooks.
5. Uses `attachments/runbook-sampling-matrix.md` as the first test set without actively
   retrofitting consumer or module runbooks.
6. Updates release notes and consumer impact classification.

Related future feedback, not part of the first runbook-authoring implementation unless explicitly
accepted later:

- slim local preference contract for language-first human-facing runbooks;
- shared environment-readiness and tooling register model for module onboarding.

Expected impact class: `instruction-only change`.

## Revision Notes

- 2026-06-23: Initial discovery drafted from Operating Kit, HQ, and Foundry runbook sampling.
- 2026-06-23: Added plan-scoped runbook sampling matrix attachment and linked it as the first
  standard test set.
- 2026-06-23: Added condensed sampling lessons and the compact runbook intention block decision.
- 2026-06-23: Added related Operating Kit feedback attachment for slim local preferences and
  shared environment/tooling readiness.
- 2026-06-23: Clarified preference language handling, audience-classification modeling, workflow
  integration scope, and no active consumer/module runbook retrofit for the first implementation.
