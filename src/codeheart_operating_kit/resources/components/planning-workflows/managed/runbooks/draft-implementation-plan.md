Last updated: 2026-11-21T10:34:09Z (UTC)

# Draft Implementation Plan

Use this runbook to turn accepted discovery, user direction, and targeted repository research into
an execution-ready `*_implementation_doc.md`.

An implementation plan gives a future agent or developer the decisions and evidence needed
to deliver the intended outcome in dependency order.

Audience: agent-facing

Intent:
Record the accepted capability scope, technical approach, dependencies, validation and
review gates. Leave local design and execution choices to the implementer.

Success:
The plan has a valid lifecycle header, complete sections, concrete epics, covered feature
capability, explicit dependencies, validation, and no hidden blocker that prevents execution.

Agent judgment boundary:
The agent may derive safe defaults from accepted discovery, user direction, and repository
research. It must not narrow the intended capability silently, turn discovery recommendations into
authority without approval, or replace feature behavior with policy-only scaffolding.

Stop boundary:
Stop before normal epic drafting when discovery handoff authority is missing, capability scope is
absent or unapproved, or a blocker prevents a single-path implementation plan.

For known bounded changes, first consider `handle-routine-change.md`. Do not require a formal plan
solely because a change touches several files. For a small change that fits the routine route, the
assignment itself can carry the outcome, limits and evidence. Work already necessary to an active
epic stays there.

## Trigger

Use this runbook when the user asks for any of the following:

- an implementation plan;
- conversion from discovery to implementation planning;
- a plan for executing a feature, migration, release, or repo change;
- a review-ready `*_implementation_doc.md`.

If the user only asks to discuss possibilities, do not create execution checklists until the user
asks for a plan or the path is clearly accepted.

## Inputs

Prefer inputs in this order:

1. Accepted discovery document.
2. User requirements and constraints from the current thread.
3. Targeted repository reconnaissance.
4. Relevant managed kit runbooks, local runbooks, references, and product docs.

Discovery is recommended but not mandatory for straightforward work. When no discovery document
exists, record the baseline problem, constraints, and assumptions in Sections 1 and 2 before
drafting execution tasks.

## Discovery Handoff Preflight

Before drafting normal implementation epics from a discovery document, verify the discovery
handoff state.

Use this preflight when the user asks to convert discovery into an implementation plan, cites a
`*_discovery_doc.md`, or the intended plan depends on decisions recorded in discovery.

Normal implementation-plan drafting may proceed only when one of these is true:

- the discovery document is explicitly implementation-handoff-ready;
- the user has approved the implementation capability scope after manual review;
- the user has delegated or revised the implementation capability scope and the revision is
  recorded in the current planning artifact;
- the user requests a blocker-resolution or conditional handoff plan and the plan scope is limited
  to that handoff type.

For a normal implementation plan derived from discovery, the discovery must include
`Implementation Capability Scope - <group name>` blocks for implementation-relevant decision
groups. Use those blocks as the primary capability source for Section 1 goals, Section 2 scope
decisions, and Section 3 epic outcomes.

Stop before normal epic drafting when:

- the discovery is only draft-ready or manual-review-ready and no approved, delegated, or revised
  capability-scope handoff is recorded;
- required `Implementation Capability Scope` blocks are absent for discovery-owned capability;
- an unresolved `BLOCKER: yes` prevents a single-path implementation plan;
- the discovery labels the next step as blocked handoff, conditional handoff, or
  blocker-resolution handoff and the requested plan would exceed that limited scope.

When this gate stops the plan, update or review the discovery document first. Do not silently
convert review-ready recommendations into implementation authority.

Discovery is still optional for straightforward work. When no discovery document exists and the
work is simple enough to plan directly, derive capability from the user request, targeted
repository research, and recorded assumptions.

## Planning Depth And Technical Judgment

A useful plan explains the problem, intended behavior, existing architecture, chosen technical
approach, consequential constraints, dependencies and evidence that will prove completion. Keep
requirements distinct from the chosen approach and assumptions. Preserve enough context for a
fresh implementer to reason about the work; brevity alone is not quality.

The director sets intent, consequential design choices, dependencies, acceptance evidence and
operating limits before dependent execution. Ask what would go wrong if the implementer chose
differently. Name exact paths or commands when they identify an interface, a protected target or
a constraint that matters to correctness or safety. Include concrete examples or data flows when
they resolve material ambiguity. Within the approved outcome the implementer chooses ordinary
files, functions, commands, test design and debugging.

Plan cheap affected checks during iteration and broader validation on a coherent candidate.
Plan the commissioning agent's primary review at meaningful checkpoints, combining tightly
related epics when useful. Add an independent reviewer where complexity, risk, the commissioning
agent's own authorship or a binding gate warrants it, and keep the same reviewers for focused
correction follow-up. Required owner and release gates remain.

## Feature Capability Coverage

Before drafting epics, identify the intended feature capability. Use the accepted discovery
document when one exists, especially any `Implementation Capability Scope` blocks. When no
discovery exists, derive the capability from the user request and targeted repository research.

The execution plan must cover that capability surface. If the plan omits part of the intended
capability, mark the omission explicitly as out of scope, deferred, or blocked with rationale. Do
not let a plan quietly narrow the capability to policy, scaffolding, gates, schemas, stubs, or
validation shells while the intended feature behavior remains unplanned.

Apply the fresh-implementer test while drafting epics. A fresh implementer must be able to identify
the chosen workflow, consequential data and permission decisions, and how to verify acceptance.
An unresolved choice that can materially affect the outcome belongs in the plan as a decision or
an explicit blocker for the affected work.

## Runbook Change Coverage

When a plan creates or materially changes durable runbooks, use
`../../agent-interface/reference/runbook-authoring-standard.md` and make the runbook-authoring scope
explicit.

The plan must state:

- each runbook created or materially changed;
- audience class for each affected runbook;
- whether each affected runbook needs human-facing flow, agent-facing execution path, hybrid
  separation, or maintainer authority and evidence handling;
- whether each affected runbook needs a compact intention block;
- whether existing consumer-owned, module-owned, or unrelated runbooks are intentionally outside
  scope.

When affected human-facing, agent-facing, or hybrid runbooks can hit missing local tooling, the
plan must also state:

- which local tool blockers can occur;
- whether each blocker is local environment readiness or module-owned service preflight;
- whether local blockers route through
  `../../agent-interface/runbooks/handle-tooling-readiness.md`;
- which concrete module-specific install commands or version requirements remain module-owned;
- how the plan avoids duplicating generic package-manager, runtime, or local-tool guidance across
  multiple runbooks.

Runbook-related acceptance criteria must cover the relevant audience checks. Do not let a plan
deliver only routing, policy, or placeholders when the intended runbook still lacks user-facing
flow, agent execution path, approval boundaries, stop conditions, evidence, or validation.

## Recipe-Maturity Coverage

When a plan creates or materially changes durable operational recipes, route-selected recipes,
durable executable mechanics, expected markers, structured summary or blocker output, reusable
script assets, promoted recipe assets, or recipe validation expectations, cite and apply
`.codeheart/kit/docs/agent-interface/reference/operational-recipe-maturity.md`.

When a plan creates or materially changes reusable script assets, first-script scaffolding,
script output contracts, script tests, script helper rules, or script-promotion review criteria,
also cite and apply
`.codeheart/kit/docs/agent-interface/reference/runbook-to-script-promotion-standard.md`.

For recipe-bearing epics, the plan must state:

- target maturity state, such as below recipe threshold, L1 structured runbook recipe, L2
  reusable script asset, L3 thin command wrapper, or L4 mature API/tool surface;
- validation tier, such as fresh-agent executability review, non-live test, dry-run or preflight,
  or approval-gated live validation;
- evidence shape and blocker shape when structured output is expected;
- promotion destination or explicit non-promotion decision;
- placement boundary for any promoted recipe asset;
- whether module-owned blocker classes, concrete package layouts, or domain-specific live
  validation are intentionally outside scope.

For reusable script asset epics, the plan must state:

- script asset role, such as primitive script, workflow script, or helper, when the role affects
  implementation or review;
- runbook caller;
- script owner and placement boundary;
- first-script scaffolding path when applicable;
- workflow dependencies, phase boundaries, and dependency contracts when the asset is a workflow
  script;
- helper placement, importing scripts, and the narrowest durable owner boundary when helpers are
  created or moved;
- whether the owner area's `scripts/README.md` needs a compact role index update;
- output contract appropriate to its caller, with structured fields only when needed;
- output safety behavior;
- managed-runner, CI, or cloud portability constraints when the script is expected to run outside
  a local interactive shell;
- proportional tests or fixtures;
- review criteria for hidden approvals, target broadening, raw output, and premature wrappers or
  APIs.

Do not treat the maturity states as a required ladder. `Do not promote yet` is a valid planned
outcome when the recipe is safe, testable, reviewable, and ergonomic at its current level.

## Routing-Standard Coverage

When a plan creates or materially changes routing-bearing surfaces, cite and apply
`../../agent-interface/reference/operation-routing-and-dispatch.md`.

Routing-bearing surfaces include:

- managed root routing;
- agent-interface routing references;
- structure-governance routing or placement rules;
- runbook-authoring rules that affect route discovery;
- capability advertisements;
- route registries;
- route cards;
- module or extension routing state rules;
- durable runbooks that select routes, owners, scopes, execution surfaces, or approval classes.

For routing-bearing epics, the plan must state:

- which routing standard sections apply;
- the capability advertisement, route registry, route card, or routing surface being created or
  changed;
- which owner maintains the route after implementation;
- whether the epic needs a fresh low-context routing probe;
- how probe evidence will be recorded.

Use a fresh low-context routing probe when the epic changes a route surface or route-selection
behavior. Select the deepest nested realistic scenario affected by the epic, give a fresh agent a
vague user-style request, and require it to identify the owner, discover the route or ambiguity
question, and avoid choosing an execution surface prematurely.

Mark the probe `not applicable` only when the changed work is not routing-bearing and explain why.

## File And Naming Rules

- Use filename pattern `<feature-slug>_implementation_doc.md`.
- Reuse the discovery slug when a matching `<feature-slug>_discovery_doc.md` exists.
- Use lowercase hyphen-separated slugs.
- Remove special characters and collapse duplicate hyphens.
- When implementation work spans multiple repositories, identify the repository that owns the work
  boundary before creating the canonical implementation plan. Place the canonical plan in that
  owning repository's planning root. Coordination derives observations from its pushed canonical
  document; it does not become the plan's canonical home.
- Put the plan in the owning `plans/` folder or plan bundle according to the planning lifecycle
  reference.

## Required Header

Every implementation document starts with:

```text
Last updated: YYYY-MM-DDTHH:MM:SSZ (UTC)
Created: YYYY-MM-DD
Status: draft
```

Use the current UTC clock for `Last updated`. Preserve `Created` after initial creation. Keep the
plan as `Status: draft` until the user explicitly approves execution.

Read the repository catalog mode before authoring. In mixed and canonical mode, add the exact
bounded metadata block from `../reference/plan-catalog-format.md` immediately below the canonical
plan title. Use kind `implementation`, choose a stable semantic ID from meaning rather than a
sequential register number, relate the discovery/family when evidence supports it, and omit
unsupported optional classifications. In legacy mode, retain compatibility behavior until an
explicit migration.

## Required Top-Level Structure

Use exactly these top-level sections in this order:

```text
# <Meaningful Plan Title>
# Section 1 - Foundation
# Section 2 - Strategy
# Section 3 - Execution Plan
# Section 4 - Future Planning
# Revision Notes
```

The first H1 is the meaningful plan title, matching semantic catalog identity. Do not use the
literal generic title `Document Header`; retain the required sections below it.

## Document Header Content

Include:

- concise overview;
- essential context reference files with reasons;
- table of contents.

The essential-context table should name concrete files, not vague areas. Include only context a
future implementer must read to execute safely.

## Section 1 - Foundation

Include these subsections:

- `1.1 Goal Of The Implementation`
- `1.2 Project And Problem Context`
- `1.3 Current State Analysis`

The goal must be measurable. State what proves completion, including user-visible behavior,
managed content changes, validation, release or migration evidence, and downstream handoff when
applicable.

Current state analysis must distinguish:

- existing systems, constraints, and problems;
- new or target systems, requirements, and ownership boundaries.

## Section 2 - Strategy

Include:

- `2.1 Implementation Strategy`
- `2.2 Open Questions And Assumptions Requiring Clarification`
- `2.3 Architectural Decisions With Reasoning`

Include a file tree when placement or layout matters to the chosen strategy. Show the paths that
carry those decisions. If helpful, annotate them with `# create` or `# modify`. Use `# delete`
only for explicitly required and safe deletions. Existing plans may keep their current headings
and trees.

Record open questions as `OQ-<n>` entries. Each open question must include:

- `BLOCKER: yes` or `BLOCKER: no`;
- `Affects:` with affected epic IDs;
- what decision the question unlocks;
- recommended default when a safe default exists.

Make consequential strategy decisions before dependent task drafting. Explain the problem,
simplest adequate solution, rationale and serious alternatives. Discuss future change only when
credible pressure affects the decision; do not invent a fixed forecast for every choice.

If no safe choice exists, keep the question as `BLOCKER: yes` and draft only blocker-resolution
tasks for affected epics.

## Decision Quality Rules

Use this decision order when choosing the implementation strategy:

1. Hard constraints.
2. Maintainability.
3. Divergence cost and risk.
4. Existing patterns when quality remains acceptable or divergence cost is too high now.

Prefer the simplest robust path that can work as an MVP without creating avoidable long-term
traps. When best practice differs from current repository patterns, contain the difference behind
a clean boundary or add explicit deferred standardization in Section 4.

Do not use existing patterns as a reason to preserve weak design when the plan is the right place
to set a better boundary.

## Whole-Plan Commissioning And Delivery Boundary

Before execution, settle the intended outcome, owning repository/branch, ordered epics, acceptance
evidence, director and delegated review points in the existing plan or assignment. Link the owning
plan or relevant Program section and the decisions the implementer needs. Include coherent
commits, normal pushes and PR creation/updates at agreed review, recovery and handover checkpoints.
Identify independently usable integration outcomes, their dependencies and relevant acceptance.
One coherent outcome may use one PR; a longer plan may use several without a PR per checkbox.
Do not leave ordinary covered effects as permission blanks to rediscover in each epic. Plan
coherent commits/publication inside implementation, not only after epic review, and apply the
Git action table in `../reference/planning-document-lifecycle.md`.

State the finish line (reviewed branch, merged main, released product or adopted consumer), who
integrates, and whether merge/release/adoption/provider effects are included, delegated or reserved.

At commissioning, check that the approved role and repository profile covers the planned local
work. Identify any separately authorized external effects in the assignment. Follow
[agent task coordination](../../agent-interface/reference/agent-task-coordination.md) for profile
reuse, corrections and handoffs.

For repository adoption, name the intended default/integration branches and active worktrees;
a working-copy pilot is partial unless that narrower finish line is explicit. Include final
record publication in the grant if closure must reach the shared baseline. Name required checks
and exact action-time inputs. Reuse sufficient applicable authority; success
of checks alone does not authorize an unspecified effect. A draft review is not execution approval.

Use `../../agent-interface/reference/agent-task-coordination.md` to commission the complete plan
with explicit report-back to the director at required epic review, completion, genuine blocker or
material scope issue. Plan delegated acceptance rather than a fresh user gate for every epic. State
how material exceptions return to the owner and how independently useful preparation may proceed
while a review is pending. Execution uses the implementing tool's native goal by default unless the
user opts out; name the phase handoff each goal ends at and require verified activation.

## Section 3 - Execution Plan

Section 3 must execute from top to bottom. Do not create competing task paths after a strategy
decision is chosen.

Start with `3.0 Epic Map`, a table containing:

- epic ID;
- one-line outcome;
- size `S`, `M`, `L`, or `XL`;
- dependencies;
- current epic status.

For each epic state the outcome, scope, acceptance criteria, size and dependencies. Include
critical-path notes, implementation notes and open questions where they matter to execution.
A file list is useful when specific paths identify an interface, protected area or starting point.
Include a task checklist when the plan needs to communicate required sequencing or work that
must be tracked individually.

The epic outcome and acceptance criteria govern completion. The plan tracks epic status, and
its sibling execution log records evidence, material deviations and any working list the
implementer needs for resume or handoff. Existing plans may keep their headings, file lists and
checkboxes. Record epic status once in the epic map or the epic itself.

## Checklist Rules

When a checklist helps, use `- [ ]` for capability-sized actions or required checkpoints. Name
the constraints and chosen approach and leave safe implementation choices open. During
execution, check existing tasks only after they are done and validated.

Do not let policy, scaffolds, schemas or tests substitute for the intended usable capability.
Resolve safely checkable consequential facts during planning. For action-time uncertainty, name
the check and how its result affects execution. Specify exact remediation and stops when risk
requires them. Give every epic validation that proves its outcome and identify any shared review
checkpoint.

## Blocker Handling

Use `BLOCKER: yes` only when the implementation path cannot be safely planned or executed without
the answer.

For an affected blocked epic:

- limit planned work on the affected epic to resolving the blocker;
- keep unaffected later epics fully planned when they remain valid regardless of the blocker.

Use `BLOCKER: no` for decisions that can be safely defaulted, deferred, or resolved during
execution without changing the main path.

## Authoring Git Checkpoints

Save the plan's own coherent completed work in a local commit after proportionate inspection and
checks. Follow `../reference/planning-document-lifecycle.md` for exceptions, authorized branch
publication and normal PR integration. Drafting alone does not authorize push or PR creation.
An accepted published or merged draft remains inactive until execution/activation is authorized.
Local commits do not add a review or permission layer. Establish covered publication once for a
shared handoff and use it without waiting for execution approval. State the canonical path,
branch/revision, open questions and actual visibility, or the reason/owner/resolving event for a
local-only exception. Main placement does not approve recommendations or activate the plan.

## Activation And Plan-Checkpoint Publication

An activation-only request for a specific implementation plan is approval to:

- set that plan to `Status: active` and make directly required plan metadata/log changes;
- create or use its unambiguous work branch;
- commit only that canonical planning checkpoint and directly required planning metadata; and
- normally push that branch so the authority is recoverable and visible to coordination.

Do not ask for a second push approval for that bounded checkpoint. The same rule applies when the
user explicitly requests a material update to an already active plan.

This activation-only authority excludes unrelated dirty files, implementation code not already included in the
requested checkpoint, another plan, ambiguous targets or branches, PR creation, merge, release,
force-push, branch deletion, history rewrite, destructive Git, credential changes, rejected-push
bypass, and any broader external action. Activation remains an L1 managed workflow; do not invent
an activation/commit/push CLI. A whole-plan execution request additionally covers the implementation
and planned Git effects declared above; do not use the activation-only boundary to ask again for
those covered steps.

Before publication:

1. resolve the exact plan, repository, and work branch;
2. inspect status and exclude unrelated changes;
3. validate mode-aware metadata and planning docs;
4. show or record the exact included planning paths;
5. create an intentional plan-only commit; and
6. perform a normal push without force.

Stop on an ambiguous target/branch, overlapping unrelated dirty content, missing auth, policy
rejection, non-fast-forward or other rejected normal push, or a repository instruction requiring
narrower authority. Report the local commit and visibility limitation; do not escalate the Git
operation. Successful normal push makes the plan observable on the next complete coordination-home
refresh. Execution may begin on that active pushed work branch without waiting for merge.

Evidence is a concise non-secret record containing repository,
branch, included paths, commit identity, normal-push result, and exclusions. Review must verify the positive plan-only authority and negative boundaries. Use meaningful
route/resource checks and a realistic walkthrough, without a wording-only test for every phrase.

## Catalog And Register Hook

When implementation planning creates or materially updates a formal plan:

1. validate lifecycle and mode-aware metadata with `plans validate`;
2. use `plans list` as the current local view;
3. use `maintain-plan-register.md` only for legacy compatibility or the stable entry point;
4. do not append numbered entries in mixed/canonical mode or create new pending-sync files; and
5. use the activation publication contract above only when the user's request activates the plan
   or materially updates an active one.

## Section 4 - Future Planning

Include:

- `4.1 Deferred Tasks`
- `4.2 Future Considerations`

Deferred tasks should name why the work is deferred and what would trigger it later. Do not hide
required implementation work in future planning.

## Quality Gate

Before finalizing the plan, verify:

- top-level section order is exact;
- every required subsection exists;
- Section 3 is linear and non-branching;
- consequential decisions are settled or explicitly blocked, and any binding paths or commands
  are concrete enough to apply;
- every epic has a status, dependencies and verifiable acceptance criteria with a validation method;
- the plan covers intended feature capability or explicitly marks omitted capability areas;
- no epic can be completed by delivering only policy, scaffolding, gates, schemas, stubs, or
  validation shells while intended capability remains unplanned;
- consequential facts are checked, or the plan gives the execution-time check and how it
  determines the next action;
- every `OQ-*` includes blocker status and affected epic IDs;
- no blocked epic contains normal implementation tasks;
- the plan is self-contained enough for a new implementer;
- the nearest docs indexes are updated when the plan is newly discoverable.
