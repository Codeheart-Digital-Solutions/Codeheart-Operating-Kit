Last updated: 2026-06-17T06:34:53Z (UTC)

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
universe.

## Checklist Rules

Every task line in `F) Tasks Checklist` must:

- start with `- [ ]`;
- contain one concrete action;
- name concrete files, commands, components, or validation gates;
- be executable without choosing between branches.

Do not use these words in checkbox tasks:

- `either`
- `or`
- `choose`
- `optionally`
- `if needed`
- `depending`
- `TBD`

End each epic with validation tasks that prove the epic outcome. Use the smallest validation set
that actually covers the changed surface.

## Blocker Handling

Use `BLOCKER: yes` only when the implementation path cannot be safely planned or executed without
the answer.

For an affected blocked epic:

- include blocker-resolution tasks only;
- do not include normal implementation tasks that depend on the answer;
- keep unaffected later epics fully planned when they remain valid regardless of the blocker.

Use `BLOCKER: no` for decisions that can be safely defaulted, deferred, or resolved during
execution without changing the main path.

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
- every `OQ-*` includes blocker status and affected epic IDs;
- no blocked epic contains normal implementation tasks;
- the plan is self-contained enough for a new implementer;
- the nearest docs indexes are updated when the plan is newly discoverable.
