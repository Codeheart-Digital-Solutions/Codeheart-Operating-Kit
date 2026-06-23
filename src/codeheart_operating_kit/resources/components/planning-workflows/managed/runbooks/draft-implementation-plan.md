Last updated: 2026-06-23T14:31:23Z (UTC)

# Draft Implementation Plan

Use this runbook to turn accepted discovery, user direction, and targeted repository research into
an execution-ready `*_implementation_doc.md`.

An implementation plan is not a brainstorming note. It is the document a future agent or developer
can execute linearly with low interpretation overhead.

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

## Feature Capability Coverage

Before drafting epics, identify the intended feature capability. Use the accepted discovery
document when one exists, especially any `Implementation Capability Scope` blocks. When no
discovery exists, derive the capability from the user request and targeted repository research.

The execution plan must cover that capability surface. If the plan omits part of the intended
capability, mark the omission explicitly as out of scope, deferred, or blocked with rationale. Do
not let a plan quietly narrow the capability to policy, scaffolding, gates, schemas, stubs, or
validation shells while the intended feature behavior remains unplanned.

Use the fresh-implementer test before writing tasks: if a future implementer can only restate what
must be true, but still has to invent the workflow, command sequence, file edits, data shape,
permission model, or validation method, the epic is not implementation-ready.

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

Runbook-related acceptance criteria must cover the relevant audience checks. Do not let a plan
deliver only routing, policy, or placeholders when the intended runbook still lacks user-facing
flow, agent execution path, approval boundaries, stop conditions, evidence, or validation.

## File And Naming Rules

- Use filename pattern `<feature-slug>_implementation_doc.md`.
- Reuse the discovery slug when a matching `<feature-slug>_discovery_doc.md` exists.
- Use lowercase hyphen-separated slugs.
- Remove special characters and collapse duplicate hyphens.
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

## Required Top-Level Structure

Use exactly these top-level sections in this order:

```text
# Document Header
# Section 1 - Foundation
# Section 2 - Strategy
# Section 3 - Execution Plan
# Section 4 - Future Planning
# Revision Notes
```

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

- `2.1 Implementation Strategy With Visual File/Folder Hierarchy`
- `2.2 Open Questions And Assumptions Requiring Clarification`
- `2.3 Architectural Decisions With Reasoning`

The file tree must show expected paths and use inline comments such as `# create`, `# modify`, and
`# delete` only when deletion is explicitly required and safe.

Record open questions as `OQ-<n>` entries. Each open question must include:

- `BLOCKER: yes` or `BLOCKER: no`;
- `Affects:` with affected epic IDs;
- what decision the question unlocks;
- recommended default when a safe default exists.

Make strategy decisions before task drafting. For each decision include:

1. problem being solved;
2. simplest working solution;
3. what may change in 6-12 months;
4. rationale for the chosen approach;
5. alternatives considered and why not chosen.

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

## Section 3 - Execution Plan

Section 3 must execute from top to bottom. Do not create competing task paths after a strategy
decision is chosen.

Start with `3.0 Epic Map`, a table containing:

- epic ID;
- one-line outcome;
- size `S`, `M`, `L`, or `XL`;
- dependencies.

For each epic include these fields:

- `A) Epic ID, Title, And Outcome`
- `B) Scope`
- `C) Files Touched`
- `D) Acceptance Criteria And Size`
- `E) Dependencies And Critical-Path Notes`
- `F) Tasks Checklist`
- `G) Implementation Notes`
- `H) Open Questions`

The epic outcome owns completion. Checkbox tasks are an execution aid, not the full possible task
universe. Each epic outcome should state what capability exists after the epic, and its tasks
should cover the concrete behavior, artifact changes, and validation needed for that capability.

## Checklist Rules

Every task line in `F) Tasks Checklist` must:

- start with `- [ ]`;
- contain one concrete action;
- name concrete files, commands, components, or validation gates;
- be executable without choosing between branches;
- represent one capability-sized implementation slice rather than one sentence or one handoff;
- include the non-negotiable details the executor must not invent.

Do not use these words in checkbox tasks:

- `either`
- `or`
- `choose`
- `optionally`
- `if needed`
- `depending`
- `TBD`

Reject checklist tasks that only state policy intent, doctrine alignment, readiness, or gate
validation without naming the concrete implementation action. Also reject tasks that deliver only
scaffolding, schemas, stubs, or validation shells when the intended feature capability still lacks
the workflow or behavior that uses them.

Resolve checkable facts during planning when they can be checked safely. If execution-time
variability is legitimate, specify the exact preflight check, expected result, remediation path,
retry validation, and stop condition.

End each epic with validation tasks that prove the epic outcome and feature capability. Use the
smallest validation set that actually covers the changed surface.

## Blocker Handling

Use `BLOCKER: yes` only when the implementation path cannot be safely planned or executed without
the answer.

For an affected blocked epic:

- include blocker-resolution tasks only;
- do not include normal implementation tasks that depend on the answer;
- keep unaffected later epics fully planned when they remain valid regardless of the blocker.

Use `BLOCKER: no` for decisions that can be safely defaulted, deferred, or resolved during
execution without changing the main path.

## Plan Register Hook

When implementation planning creates or materially updates a `*_implementation_doc.md`, maintain
the local plan register if the change affects plan identity, scope, lifecycle status, parent or
child relationships, dependencies, supersession, related plans, implementation path, or review
outcome.

Use `maintain-plan-register.md` for the procedure and `../reference/plan-register-format.md` for
entry shape. The sequence is:

1. Update `docs/repo/plans/plan-register.md` in the local repository.
2. When portfolio coordination is configured and the coordination home is available and safe to
   edit, update the configured coordination-home register.
3. When portfolio coordination is configured but the coordination home is unavailable or unsafe to
   edit, record pending sync in `docs/repo/plans/coordination-sync-pending.md` and continue the
   local planning task.

Record creating or material-update session refs when a session ID is available. Do not block
implementation planning when no session ID is available.

Do not update the register for typos, formatting-only edits, timestamp-only edits, or mechanical
checklist progress that does not change lifecycle, relationships, scope, or implementation path.

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
- every epic has a file tree;
- every epic has verifiable acceptance criteria;
- every epic ends with validation tasks;
- the plan covers intended feature capability or explicitly marks omitted capability areas;
- no epic can be completed by delivering only policy, scaffolding, gates, schemas, stubs, or
  validation shells while intended capability remains unplanned;
- avoidable non-concreteness has been resolved into checked facts or exact execution-time
  preflight/remediation paths;
- every `OQ-*` includes blocker status and affected epic IDs;
- no blocked epic contains normal implementation tasks;
- the plan is self-contained enough for a new implementer;
- the nearest docs indexes are updated when the plan is newly discoverable.
