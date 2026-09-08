Last updated: 2026-09-07T22:40:39Z (UTC)

# Review Planning Document

Use this runbook to review discovery or implementation documents for technical quality,
architecture quality, ambiguity, and execution readiness.

The review output should help the owner fix the document before execution, not summarize the
document.

Audience: agent-facing

Intent:
Review discovery and implementation documents for material gaps, ambiguity, unsafe sequencing,
weak validation, missing authority, and incomplete capability coverage before execution.

Success:
The review identifies actionable findings with severity and file or section references, or states
that no material issues remain with residual risk.

Agent judgment boundary:
The agent may choose review depth based on blast radius and evidence. It must not rewrite the
document unless asked, hide material risks behind summary, or treat missing handoff/probe evidence
as ready.

Stop boundary:
Stop at a `Blocked` review result when execution depends on a user decision, missing handoff
authority, unresolved blocker, or external condition.

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

## Phase And Checkpoint

Review discovery for intention, evidence, consequential decisions and the requested handoff.
Review a plan for substantive approach, capability coverage, dependencies, adaptable technical
detail and proportionate validation. Review implementation against actual behavior and evidence;
a planning review cannot prove implementation or release readiness.

Use one independent review at a meaningful coherent checkpoint and the same reviewer for focused
correction follow-up. Broaden or replace the reviewer for material design/impact change,
inadequate independence, reviewer limitation or unresolved concern. Do not restart unchanged
review after every finding. Label optional improvements separately from material blockers; plain
word choice, record duplication and ordinary in-scope refinement are not defects by themselves.

Check local commits, authorized publication and accepted PR integration independently from plan
lifecycle. A merged draft is still draft; publication is not execution approval. No extra review
or approval stage is needed per local commit.

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
- checklist tasks provide a chosen approach and allow safe action-time refinements;
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
- plans that create or materially change routing-bearing surfaces cite and apply
  `../../agent-interface/reference/operation-routing-and-dispatch.md`;
- routing-bearing epics identify affected capability advertisements, route registries, route
  cards, owner boundaries, or routing surfaces;
- routing-bearing epics include fresh low-context routing probes, or explicitly justify why probes
  are not applicable;
- probe pass criteria check that a fresh agent identifies the likely owner, discovers the route or
  ambiguity question, and avoids choosing an execution surface prematurely;
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
- recipe-bearing plans cite
  `.codeheart/kit/docs/agent-interface/reference/operational-recipe-maturity.md`;
- recipe-bearing epics name target maturity state, validation tier, evidence shape, and promotion
  boundary;
- reusable-script-asset epics cite
  `.codeheart/kit/docs/agent-interface/reference/runbook-to-script-promotion-standard.md`;
- reusable-script-asset epics name the runbook caller, script owner, placement boundary, output
  contract, output safety behavior, tests or fixtures, and review flags;
- reusable-script-asset epics name the script role when the role affects implementation or review;
- workflow-script epics document dependencies, phase boundaries, blocker ownership, and
  dependency contracts;
- helper work identifies helper placement, importing scripts, and why the helper is not a hidden
  runbook entrypoint;
- script-bearing plans cover `scripts/README.md` role-index updates where review clarity needs
  them;
- script-bearing plans preserve portability boundaries for managed runners, CI, or cloud
  orchestration when those execution contexts are intended;
- recipe-bearing plans preserve `do not promote yet` as a valid outcome when justified;
- promoted recipe assets have owner, placement boundary, validation path, and discoverability
  route;
- recipe plans do not promote scripts, commands, wrappers, or APIs prematurely;
- recipe plans do not invent an L2 `command_wrapper` role or create premature wrappers where an
  L2 primitive or workflow script is sufficient;
- recipe plans do not hide executable behavior, structured blockers, markers, or validation
  expectations inside vague prose.
- recipe plans do not preserve long inline implementations as durable assets when reusable script
  assets are the safer execution surface;
- recipe plans do not hide approval behavior, target broadening, raw sensitive output, or missing
  tests behind reusable-script language.

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
- repository catalog mode is identified and mixed/canonical metadata uses the exact bounded
  marker, stable semantic ID, matching kind, public-safe purpose, and correct chronology;
- metadata-only migration preserves historical `Last updated`, while a real content/lifecycle
  change updates it normally;
- discovery, implementation, and family records remain distinct and use evidence-backed families,
  relations, classifications, and legacy aliases;
- section order is correct;
- Section 3 is linear;
- every epic has acceptance criteria and validation tasks;
- every epic outcome and checklist set covers the intended feature capability or records explicit
  omissions;
- blocker handling is coherent;
- future planning does not hide required work;
- execution log expectations are present for goal-style runs.
- activation-only wording, when present, grants only plan/log/direct-metadata branch creation/use,
  commit, and normal push; excludes unrelated files, code, PR, merge, release, force-push,
  deletion, destructive Git, ambiguity, auth/policy/rejected-push bypass; and permits execution on
  the active pushed branch without merge;
- whole-plan execution wording separately covers planned implementation, validation, coherent
  commits, normal pushes and PR creation/updates; one implementer can execute all ordered epics
  with delegated director acceptance and explicit report-back, without per-step user prompts;
- the finish line, integration owner, included/delegated/reserved effects and required exact inputs
  are explicit; existing sufficient authority is reused without weakening binding gates;
- work communicates its named plan/epic and actual Git/delivery state; optional goals require an
  explicit request and verified activation, and cannot complete at an intermediate review handoff;
- necessary corrections stay in their epic; genuinely separate routine repairs retain a visible
  relationship and do not evade review;
- a meaningful plan title is used instead of the literal H1 `Document Header`;
- mixed/canonical guidance uses generated views rather than manual register appends, and portfolio
  claims distinguish pushed remote observations from local-only work.

For discovery documents, verify:

- intention, scope, requirements, decisions, assumptions, risks, and open questions align;
- blocker status is accurate;
- recommendations are concrete enough for user review;
- implementation handoff readiness is not claimed when required implementation-shaping decisions
  remain unresolved.
- mode-aware metadata and local-versus-pushed visibility are accurate when the repository uses the
  semantic catalog.

## Final Review Statement

End with one of:

- `Ready`: no material issues remain.
- `Needs improvement`: material issues should be fixed before execution.
- `Blocked`: execution or handoff depends on a user decision or external condition.

State residual risk even when the document is ready.
