Last updated: 2026-06-25T12:14:15Z (UTC)

# Review Planning Document

Use this runbook to review discovery or implementation documents for technical quality,
architecture quality, ambiguity, and execution readiness.

The review output should help the owner fix the document before execution, not summarize the
document.

## Trigger

Use this runbook when the user asks to:

- review a discovery document;
- review an implementation document;
- find gaps, ambiguities, or risks in a plan;
- check whether a planning document is ready for execution.

Do not rewrite the source document unless the user explicitly asks for edits.

## Review Stance

Lead with findings. Put summaries and praise after issues or omit them.

Prioritize:

- bugs in the plan;
- missing decisions;
- hidden blockers;
- weak acceptance criteria;
- ambiguous tasks;
- unsafe sequencing;
- insufficient validation;
- ownership or boundary mistakes.

If no material issues are found, say that clearly and mention residual test gaps or remaining
risk.

## Required Inputs

Read:

- target planning document;
- referenced discovery or implementation document;
- referenced architecture notes, runbooks, and source inventories that materially affect the plan;
- local routing docs when placement or lifecycle is under review.

Use targeted context. Do not perform a broad repository sweep unless the document scope requires
it.

## Review Areas

### Decision Soundness

Check whether:

- problem statements are clear and correctly scoped;
- trade-offs are explicit;
- chosen decisions follow constraints;
- alternatives were considered where they matter;
- hidden assumptions should be explicit.

### Completeness

Check whether:

- required workstreams are missing;
- dependencies and sequencing are realistic;
- migration, release, rollback, and handoff concerns are covered when applicable;
- security, reliability, observability, or public-safety concerns are addressed to the needed
  level.

### Ambiguity And Risk

Check whether:

- requirements are testable;
- tasks allow conflicting interpretations;
- independent implementers could produce incompatible outcomes;
- open questions have correct blocker status;
- non-blocker defaults are safe.

### Best-Practice Alignment

Check whether the approach fits maintainability, testability, operability, and safety best
practices. Call out weak legacy patterns when they are copied without justification.

### Execution Readiness

For implementation docs, check whether:

- epics are linear;
- each epic has clear outcome, scope, files, acceptance criteria, dependencies, tasks, notes, and
  open questions;
- validation gates cover the changed surface;
- checklist tasks are concrete and non-branching;
- release or migration authority is explicit;
- the plan concretely implements the intended feature capability, not only surrounding policy,
  scaffolding, gates, schemas, stubs, or validation shells;
- the plan preserves accepted discovery goals, non-goals, decisions, blockers, and capability scope
  when a discovery document exists;
- a discovery-derived normal implementation plan has implementation-handoff-ready discovery or
  recorded user approval, delegation, or revision of the relevant `Implementation Capability
  Scope` blocks;
- omitted capability areas are explicitly out of scope, deferred, or blocked;
- a lazy implementer cannot complete every checklist item while delivering only scaffolding,
  policy, stubs, validation shells, or a narrow slice that does not fulfill the intended feature
  capability;
- checkable facts are resolved or represented by exact preflight checks, expected results,
  remediation paths, retry validation, and stop conditions.
- plans that create or materially change durable runbooks route to
  `../../agent-interface/reference/runbook-authoring-standard.md`;
- affected runbooks have clear audience classes and compact intention-block requirements;
- human-facing runbook work includes user-visible flow, question pacing, explicit action wording,
  and non-technical language boundaries;
- agent-facing runbook work includes concrete execution paths, evidence, validation, and stop
  conditions;
- hybrid runbook work separates user dialogue from agent-only execution;
- maintainer-facing runbook work preserves authority, evidence, rollback, and validation
  boundaries;
- runbook work that can hit missing local tooling names the blocker, routes generic local
  readiness to `../../agent-interface/runbooks/handle-tooling-readiness.md`, and keeps
  module-owned service preflight separate;
- human-facing or hybrid runbook work presents blocker-specific choices instead of vague
  "install tools" prompts when the blocker can be classified;
- runbook work does not duplicate generic package-manager, runtime, or local-tool setup guidance
  across multiple managed runbooks;
- runbook plans do not hide missing user-facing flow or agent execution detail behind routing,
  policy, placeholders, or broad future audits.

For discovery docs, check whether:

- requirements, decisions, open questions, assumptions, and risks connect coherently;
- blockers are separated from implementation-shaping decisions;
- recommendations have evidence and trade-off reasoning;
- the handoff target is clear.

## Severity

Use:

- `High`: can cause implementation failure, public safety issue, security issue, or major rework.
- `Medium`: likely to create ambiguity, missed work, weak validation, or review churn.
- `Low`: useful improvement with limited execution risk.

Every finding should include:

- severity;
- concrete file and line or section reference;
- why it matters;
- actionable recommendation.

## Output Order

Use this order by default:

1. Critical risks.
2. Gaps in scope or plan.
3. Ambiguities to resolve.
4. Decision-quality concerns.
5. Best-practice improvements.
6. Targeted questions.

Ask at most five targeted questions. Include a recommended default and blast radius for each.

## Planning-Document-Specific Checks

For implementation documents, verify:

- required lifecycle header exists;
- section order is correct;
- Section 3 is linear;
- every epic has acceptance criteria and validation tasks;
- every epic outcome and checklist set covers the intended feature capability or records explicit
  omissions;
- blocker handling is coherent;
- future planning does not hide required work;
- execution log expectations are present for goal-style runs.

For discovery documents, verify:

- intention, scope, requirements, decisions, assumptions, risks, and open questions align;
- blocker status is accurate;
- recommendations are concrete enough for user review;
- implementation handoff readiness is not claimed when required implementation-shaping decisions
  remain unresolved.

## Final Review Statement

End with one of:

- `Ready`: no material issues remain.
- `Needs improvement`: material issues should be fixed before execution.
- `Blocked`: execution or handoff depends on a user decision or external condition.

State residual risk even when the document is ready.
