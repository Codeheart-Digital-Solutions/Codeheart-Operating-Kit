Last updated: 2026-06-22T18:05:35Z (UTC)

# Discovery Workflow

Use this runbook for substantial discovery, requirements clarification, architecture or operating
model framing, hard-to-reverse decisions, and discovery document drafting.

Discovery is not a task checklist. Discovery creates enough shared understanding, evidence,
decision state, and handoff clarity for the requested next step.

## Core Model

Use one discovery lifecycle across domains.

The domain can change. The mechanics stay the same:

1. Clarify intention.
2. Detect candidate domains, boundaries, and decision clusters when the request is ambiguous.
3. Gather evidence before broad questions.
4. Build or update a decision ledger.
5. Work decisions in dependency order.
6. Produce concrete recommendations or justified `No safe default` packets.
7. Use reviewer scrutiny for important decisions.
8. Record the decision state.
9. Stop at the requested readiness gate.

Avoid special-case workflow branches. Treat "feature discovery", "business discovery",
"repository discovery", "goal-style discovery", and "document drafting" as different input states,
domains, or output targets inside the same lifecycle.

## When To Use This Runbook

Use this runbook when any of these are true:

- the user asks for discovery, architecture exploration, requirements clarification, or decision
  framing;
- a `*_discovery_doc.md` is being created or materially updated;
- an existing discovery document has `BLOCKER: yes` open questions;
- the work involves hard-to-reverse decisions about ownership, boundaries, source of truth, naming,
  security, identity, lifecycle, compliance, business model, operating model, user behavior, or
  external interfaces;
- the user asks to progress an existing discovery document with autonomy;
- the requested next step needs a decision ledger rather than a direct answer.

Do not use this runbook for small wording edits, straightforward implementation-plan execution, or
single-answer questions that do not need decision traceability.

## Discovery Inputs

Start from the available input state:

- `unknown-domain-request`: the user describes a broad goal, concern, or desired outcome but may not
  know which domain or workstream it belongs to.
- `new-request`: the user has described a problem, goal, concern, or proposed solution.
- `existing-discovery-doc`: a discovery document already exists and should be progressed.
- `existing-decision-ledger`: decisions, requirements, risks, or open questions already exist in
  another artifact.
- `review-or-cleanup`: the user wants gaps, ambiguity, or decision quality reviewed.
- `handoff-preparation`: discovery content exists and needs to be made ready for implementation
  planning or another next step.

For an existing discovery document, treat the document as the state ledger. Refresh it before
resolving new decisions. Do not assume previous content is correct without checking current
evidence and user intent.

## Output Targets

Choose the output target before working. If the user did not name one, infer the smallest useful
target and state it briefly.

### Exploration Ready

Use this target when the user mainly needs orientation.

Discovery is exploration-ready when:

- the starting problem and likely scope are restated clearly;
- candidate domains and boundaries are named when the request may span more than one domain;
- key uncertainties are named;
- the next useful questions, evidence checks, or focused discovery cluster are known.

### Draft Ready

Use this target when the user wants a discovery document draft.

Discovery is draft-ready when:

- the intention contract is recorded or explicitly marked incomplete;
- candidate domains and decision clusters are recorded when domain fit is still uncertain;
- meaningful requirements, decisions, open questions, assumptions, and risks are captured;
- blockers are labeled with `BLOCKER: yes/no`;
- the document clearly states what still needs user judgment or evidence.

A draft may carry blockers when the user asks for an early draft. Do not present a blocker-carrying
draft as ready for implementation planning.

### Manual-Review Ready

Use this target when progressing an existing discovery document with autonomy, when decisions need
user review before implementation planning, or when the user asks for decision-by-decision
progress.

Discovery is manual-review-ready when:

- the decision inventory is coherent;
- visible blocker and implementation-shaping decisions have concrete recommendations, delegated or
  approved outcomes, or justified `No safe default` packets;
- reviewer gates are complete where required;
- reviewer exchange summaries are recorded where used;
- remaining open questions are non-blocking follow-ups or clear manual-choice points;
- the exact user action needed is visible.

Manual-review-ready does not mean decisions are closed. It means the user can now approve, reject,
delegate, defer, or correct the recommendations with enough context.

### Implementation-Handoff Ready

Use this target only when the user asks for implementation planning input, implementation handoff,
or a next-step handoff that needs a single coherent path.

Discovery is implementation-handoff-ready when:

- user intention, non-goals, and priority order are stable enough for planning;
- MVP or first-step requirements are clear;
- critical non-functional requirements are clear;
- implementation-shaping decisions are approved, delegated, superseded, rejected as out of scope,
  or explicitly non-blocking;
- no required implementation-shaping decision remains deferred;
- no unresolved `BLOCKER: yes` prevents a single-path plan;
- decision clusters have been promoted to candidate workstreams, implementation groups, or explicit
  follow-ups only where dependencies and ownership boundaries are clear;
- risks, assumptions, and validation needs are recorded;
- each implementation-relevant decision group included in the normal handoff has an
  `Implementation Capability Scope - <group name>` block;
- the handoff lists frozen inputs and remaining follow-ups.

For each implementation-relevant decision group included in a normal handoff, add a compact block
using this shape:

```text
## Implementation Capability Scope - <group name>

Capability:
<what should be possible after implementation>

Primary workflow:
<who or what uses it and through which workflow>

Must cover:
- <feature behavior the implementation plan must cover>

Explicitly out of scope:
- <exclusion>

Deferred or blocked:
- <capability piece>: <deferred reason or blocker ID>

Preserve decisions:
- <D-* decision>

Planner must not reinvent:
- <name, value, permission, constraint, or validation expectation>

Feature-level success evidence:
- <what proves the capability works at feature level>
```

The capability-scope block is a planning handoff target, not an implementation plan. Do not add
epics, execution checklists, or sentence-level implementation tasks to discovery handoff. If the
capability cannot be summarized concretely enough for a single-path plan, label the output as
`blocked handoff`, `conditional handoff`, or `blocker-resolution handoff` instead of normal
implementation-handoff-ready.

If the user explicitly asks to carry a blocker forward, label the handoff as `blocked handoff`,
`conditional handoff`, or `blocker-resolution handoff`. Do not label it as a normal single-path
implementation handoff.

## Discovery Status Checkpoints

Users do not need to know the readiness gates. For substantial discovery, surface the agent's
current interpretation in one or two sentences.

At the start of substantial discovery, state:

- input state;
- current output target;
- why that target is the smallest useful next target;
- whether implementation planning is intentionally out of scope for now.

When the stage changes, say so briefly.

Examples:

- "This looks like an unknown-domain request. I will aim for exploration-ready first so we can
  identify candidate domains before turning this into a plan."
- "This is now draft-ready: the main decisions and blockers are captured, but it is not ready for
  implementation planning yet."
- "This is manual-review-ready, not implementation-handoff-ready: the recommendations are prepared,
  but you still need to approve, reject, delegate, or correct them."
- "This blocker prevents a single-path implementation handoff. The next useful output is a
  blocker-resolution handoff."

Before creating an implementation handoff, explicitly state why the discovery is or is not
implementation-handoff-ready.

Do not turn status checkpoints into ceremony for small tasks. Infer the right target, state the
inference, and let the user correct it.

## Deep Discovery Trigger

The user can trigger the most autonomous discovery style with a short phrase.

Recognize these as Deep Discovery triggers:

- "Run deep discovery."
- "Deep discovery on this."
- "Run deep discovery on this discovery doc."
- "Take this discovery to manual-review-ready."
- "Progress this discovery autonomously."
- `/goal <path-to-discovery-doc>`
- `/goal Run deep discovery on <path-to-discovery-doc>`

Do not require the user to mention helper agents, decision inventory, reviewer gates, or readiness
gates.

Deep Discovery means:

- Use the current request, referenced artifact, or existing `*_discovery_doc.md` as the input.
- Default the output target to `Manual-Review Ready`.
- State a brief discovery status checkpoint before deep work.
- Build or refresh the decision inventory.
- Detect candidate domains and decision clusters when the request is broad or ambiguous.
- Work blocker and implementation-shaping decisions in dependency order.
- Use research helpers by default when bounded evidence gathering, file or document inspection,
  option comparison, or external constraint checks would improve the decision.
- Use review helpers by default for blocker, hard-to-reverse, security, identity, architecture,
  ownership, compliance, cross-domain, or broad operating-model decisions.
- Keep the main agent responsible for the intention contract, decision ledger, final
  recommendations, decision state, user-facing summary, and document updates.
- Stop at `Manual-Review Ready` unless the user explicitly asks for implementation handoff.

Goal-style execution:

- Treat Deep Discovery as a user-facing trigger for goal-style autonomous discovery.
- When the user invokes `/goal` with a `*_discovery_doc.md` path, treat that as enough information
  to run Deep Discovery on the referenced document unless the user explicitly scopes the goal
  differently.
- The default goal objective is: take the referenced request, artifact, or discovery document to
  `Manual-Review Ready` through Deep Discovery.
- If the formal goal mechanism is available, activate or continue it with that objective and the
  success criteria below.
- If the formal goal mechanism is unavailable, execute the same contract in the current thread and
  state that no formal goal mechanism was used.

Deep Discovery goal success criteria:

- A discovery status checkpoint records input state, output target, goal success criterion, and
  whether implementation handoff is out of scope.
- Candidate domains and decision clusters are recorded when the request or document may span more
  than one domain.
- The decision inventory is refreshed. Every visible blocker and implementation-shaping decision
  has an ID, class, state, owner, dependencies, affected IDs, and review or closure criteria.
- Every visible blocker and implementation-shaping decision has a decision packet, evidence
  summary, option comparison, concrete recommendation or justified `No safe default` packet, and
  updated decision state.
- Review helper critique is recorded for every blocker or implementation-shaping decision unless a
  review helper is unavailable or explicitly waived with a reason.
- For every reviewed decision, the main agent response to helper critique is recorded before the
  decision is marked `review-ready recommendation`, `manual choice required`, or closed.
- Helper/main-agent exchange summaries are recorded. Each summary includes helper role, decision
  ID, packet or recommendation sent to the helper, helper objections or missing evidence,
  main-agent response, recommendation changes, and remaining dissent or uncertainty.
- Newly discovered blocker or implementation-shaping decisions are added to the inventory and
  processed through the same criteria, unless explicitly classified as non-blocking with rationale.
- A manual review packet lists approved or delegated decisions, review-ready recommendations,
  `No safe default` choices, remaining non-blocking follow-ups, reviewer dissent or residual
  uncertainty, and exact user actions needed.

Do not mark a Deep Discovery goal complete after only an inventory, broad summary, partial
recommendation set, or unreviewed blocker decision unless the user explicitly scoped the goal that
narrowly.

Reasoning guidance:

- Use high reasoning by default for Deep Discovery.
- Use the strongest available reasoning level when the environment permits it for hard-to-reverse,
  security-critical, identity, architecture, ownership, compliance, or broad cross-domain
  decisions.
- Do not hardcode a specific model version in managed workflow doctrine.

## Autonomy And Decision Authority

Separate recommendation authority from closure authority.

By default, the agent may:

- infer and propose intention;
- research;
- identify decisions;
- recommend options;
- prepare `No safe default` packets;
- mark decisions as review-ready.

The agent may close a decision only when the user approves, rejects, delegates, or defers it, or
when the user has explicitly granted decision authority for that class of decision.

Recommendations are not approvals by themselves.

## Universal Discovery Lifecycle

### 1. Select Target And Input State

Before deep work, state the effective input state and output target when useful.

Examples:

- "I will treat this as a new-request discovery and aim for draft-ready."
- "I will treat the existing discovery document as the ledger and aim for manual-review-ready."
- "I will check whether the document is implementation-handoff-ready."

Do not over-explain this to the user for small turns. The point is to guide the work, not create
ceremony.

### 2. Restate The Starting Point

Restate:

- the problem or opportunity;
- the intended outcome;
- known constraints;
- assumed scope;
- what would make discovery useful.

If the starting point is too vague for useful research, ask a small number of orientation questions
first.

### 3. Detect Candidate Domains And Decision Clusters

When the user may not know which domain the request belongs to, do not force the request into one
domain from the first noun or example.

First identify candidate domains and boundaries. Candidate domains may include, depending on the
request:

- software or product work;
- business operations;
- documentation or knowledge management;
- automation or local runner workflows;
- Microsoft 365, SaaS, or external-service integration;
- data, privacy, security, compliance, or identity;
- repository, workspace, governance, naming, or operating model;
- sourcing, vendors, finance, market, customer, or research.

For each plausible domain, record:

- why it might be involved;
- what owner, source of truth, system, repository, workspace, or team may govern it;
- what evidence would confirm or reject it;
- which decisions it likely creates;
- whether it is a first-step candidate or a later dependency.

Use `decision clusters` before implementation groups. A decision cluster is a provisional grouping
of related decisions or uncertainties. It helps orient discovery, but it is not a promise that a
future implementation epic or repository exists.

Only promote clusters into candidate workstreams, implementation groups, or epics when:

- the intention contract is stable enough;
- ownership and source-of-truth boundaries are known;
- dependencies between clusters are visible;
- the selected output target needs sequencing.

A discovery may validly end with a domain map and a recommended next focused discovery target. It
does not need to force implementation groups when the correct domain boundary is still unclear.

### 4. Establish The Intention Contract

Before resolving design or operating decisions, establish the highest useful level of user
intention.

Record or confirm:

- `User Intention`: what the user is trying to accomplish.
- `Non-Goals`: what the user explicitly does not want.
- `Priority Order`: the tradeoff order, such as clarity, safety, future-proofing, operability,
  auditability, minimal churn, speed, compatibility, cost, or risk.
- `Durable Principles`: what should remain true across future products, users, clients, operators,
  workflows, or implementations.
- `Success Criteria`: what must be true for discovery to be considered ready for the selected
  output target.

The intention contract is clear enough when it can reject at least one plausible option. If every
option still seems equally compatible with the intention, clarify before treating recommendations
as review-ready.

### 5. Extract Intention From User Wording

Treat user wording as evidence, not automatically as final intention. Users may describe examples,
symptoms, concerns, proposed tools, partial solutions, or mechanisms.

Separate the input into:

- `Outcome`: what should become possible.
- `Pain`: what mistake, friction, ambiguity, or failure should be avoided.
- `Constraint`: what must remain true.
- `Proposed Mechanism`: a tool, path, framework, service, model, package, name, or process the
  user mentions.
- `Future Pressure`: what may be added later and should not be made harder by today's design.

Then abstract upward:

- Do not treat a proposed mechanism as the intention unless the user makes it a hard requirement.
- Convert mechanism-level language into durable intent and decision-policy candidates.
- Offer two to four candidate intention statements when the user is struggling to formulate intent.
- Ask for confirmation or correction before using inferred intention as a decisive basis.
- Record confirmed intent as `User Intention`, `Durable Principles`, `Priority Order`,
  `Non-Goals`, or `Decision Policy`.

### 6. Use The Intention Ladder

When a decision exposes unclear or conflicting intention, move up before answering:

1. Outcome: what should become possible?
2. Durable principle: what should remain true across future work?
3. Decision policy: what rule should guide this class of decisions?
4. Design choice: which architecture, workflow, or model follows from that policy?
5. Concrete change: which file, package, path, name, field, process, or artifact changes?

If blocked at level 4 or 5, ask or reason at level 1, 2, or 3.

Useful intention questions:

- What future mistake are we trying to avoid?
- What should a new participant infer without private context?
- Which mistake is worse for this work?
- Is the priority reuse, ownership clarity, minimal churn, auditability, operability, speed,
  compatibility, cost, or risk?
- Should this decision establish a general policy or solve only the current case?

### 7. Research Before Broad Questions

Run targeted research before asking broad follow-up questions.

List essential context files or sources with one-line reasons. Use the evidence type that fits the
domain:

- code, tests, configuration, validators, docs, runbooks, and existing plans;
- operating records, process docs, policy sources, meeting notes, or internal decisions;
- supplier, vendor, market, customer, regulatory, financial, or external platform evidence.

Summarize findings briefly and state which questions or assumptions the research closed.

Do not ask the user questions that targeted research can answer safely.

### 8. Create Or Refresh The Ledger

Use stable traceability IDs:

- `FR-*`: functional requirements.
- `NFR-*`: non-functional requirements such as security, privacy, performance, reliability,
  accessibility, observability, auditability, maintainability, cost, or operability.
- `D-*`: decisions.
- `OQ-*`: open questions.
- `A-*`: assumptions.
- `R-*`: risks.

Every `OQ-*` must include:

- question;
- owner;
- `BLOCKER: yes` or `BLOCKER: no`;
- affected IDs;
- required information or resolution path.

`BLOCKER: yes` means no safe single-path next-step handoff exists until the question is resolved,
explicitly accepted as a blocked or conditional handoff, delegated with a concrete resolution path
that is sufficient for the selected output target, or reclassified as non-blocking with rationale.

### 9. Build The Decision Inventory

For substantial discovery, maintain a decision inventory.

Each decision should have:

- stable ID;
- purpose;
- exact question;
- owner;
- state;
- class;
- depends on;
- blocks;
- affected `FR-*`, `NFR-*`, `OQ-*`, `A-*`, and `R-*`;
- affected docs, plans, code, workflows, policies, external interfaces, operations, or handoff
  sections;
- closure criteria.

Decision classes:

- `blocking`: must be resolved before the selected output target can be coherent.
- `implementation-shaping`: changes architecture, paths, APIs, packages, security model, naming,
  data model, lifecycle, external interface, operating model, or validation strategy.
- `non-blocking`: can remain a follow-up without changing the selected next-step path.
- `duplicate/subsumed`: covered by another decision.
- `not-a-decision`: factual research, implementation task, validation task, or wording cleanup.

Default decision order:

1. intention, non-goals, durable principles, and priority order;
2. ownership, boundary, and placement;
3. source of truth and authority;
4. security, identity, lifecycle, policy, and human-facing behavior;
5. profile, data model, API, integration, and scope;
6. validation, documentation, rollout, and handoff.

Work on one active decision at a time unless the decisions are independent and the user accepts
parallel discovery.

### 10. Work The Active Decision

For each blocker or implementation-shaping decision, prepare a decision packet.

```text
Decision:
<D-* or OQ-* and title>

User Intention:
<highest-level intention relevant to this decision>

Non-Goals:
<what this decision must not optimize for or change>

Priority Order:
<tradeoff order for this decision>

Decision Relevance:
<why this decision matters to the intention>

Exact Question:
<the concrete decision to answer>

Scope:
<in scope>

Out Of Scope:
<not covered by this decision>

Options:
<viable choices>

Evaluation Criteria:
<how options will be judged>

Required Evidence:
<workspace files, docs, external references, examples, tests, constraints, or other evidence>

Recommendation Standard:
<what a valid recommendation must include, or what would make No safe default valid>

Risks And Tradeoffs:
<known risks and tradeoffs>

Downstream Effects:
<affected requirements, risks, docs, plans, code, validators, runbooks, operations, or handoff>

Closure Criteria:
<what must be true before the decision is closed>
```

Gather the smallest evidence set that can answer the decision safely.

Produce:

- concise evidence summary;
- option comparison against evaluation criteria;
- concrete recommendation or `No safe default` packet;
- alternatives considered;
- risks and mitigations;
- residual uncertainty;
- exact user decision request, delegated-decision statement, or review-ready state.

### 11. Use Helper Agents For Bounded Research And Review

Helper agents are an execution mechanism inside the same discovery lifecycle. They do not create a
separate workflow or change the decision authority model.

Use helper agents when bounded parallel work improves evidence quality, review quality, or speed
without closing decisions from stale assumptions.

Useful helper roles:

- `research helper`: gathers evidence, inspects files or documents, compares bounded options,
  checks external constraints, or summarizes domain facts for the active decision.
- `review helper`: critiques the active decision packet, recommendation, evidence sufficiency,
  dependency order, risks, and `No safe default` claims.

The main agent owns:

- the intention contract;
- the decision inventory;
- the active decision packet;
- final recommendation or `No safe default` packet;
- user-facing summary;
- decision state;
- discovery document or ledger updates.

Helper agents must not independently close decisions, change decision state, update the canonical
discovery document, or present their output as the final recommendation.

Parallel research is allowed when:

- tasks are evidence-gathering tasks, not closure tasks;
- the active decision packet defines the scope;
- the work does not depend on an unresolved parent decision;
- outputs can be reconciled by the main agent before recommendation;
- conflicting helper findings are resolved before a decision becomes review-ready.

Parallel review is allowed when:

- decisions are independent according to the dependency metadata; or
- the user explicitly asks for bulk review.

Dependent decisions must receive final review in dependency order. Do not mark a dependent decision
as `review-ready recommendation` while its parent boundary, authority, security, lifecycle, or
validation decision is unresolved, except as an explicit manual-choice blocker.

When the user asks for per-decision scrutiny:

- send one active decision packet to the reviewer helper at a time;
- include the intention contract, decision inventory, dependency map, approved and delegated
  decisions, unresolved blockers, and relevant evidence;
- ask the reviewer to critique only the active decision;
- the main agent must respond to the critique before moving to the next decision;
- if the reviewer identifies a new blocker or implementation-shaping decision, add it to the
  inventory and process it through the same lifecycle;
- one broad review of many decisions counts only as inventory review or closure-sweep review unless
  the user explicitly accepts bulk review.

Record helper outputs as concise summaries. Do not paste full helper transcripts by default.

### 12. Handle Broad Option Landscapes

When a decision depends on choosing between external tools, platforms, vendors, services,
technologies, operating models, suppliers, or comparable alternatives, consider a bounded landscape
scan.

Use the scan only when it improves the decision. Keep it bounded:

- compare three to five viable options;
- map options to the active `D-*`;
- include pros and cons;
- evaluate fit to relevant `NFR-*`;
- include integration or adoption cost;
- include lock-in, migration, or reversibility risk;
- name official or authoritative sources when external evidence is used;
- end with a recommended default plus one or two serious alternatives, or a justified
  `No safe default` packet.

Do not let broad scans expand discovery into open-ended research. The scan exists to support a
specific decision.

### 13. Ask High-Impact Questions

Ask only high-impact or hard-to-reverse questions. Ask in batches of no more than five.

For each question, include:

- question;
- recommended default when one is safe;
- one or two alternatives;
- downstream blast radius;
- `BLOCKER: yes/no`;
- owner.

If no safe recommendation exists, say why and use the `No safe default` standard.

## Decision States

Use decision states consistently.

- `candidate`: possible decision, not yet fully framed.
- `researched`: evidence exists but no recommendation is ready.
- `in-review`: recommendation is being reviewed or challenged.
- `review-ready recommendation`: recommendation is ready for user review but not approved.
- `manual choice required`: no safe default exists or user judgment is explicitly required.
- `approved`: user approved the decision.
- `delegated`: user delegated the decision to the agent, a rule, or an owner.
- `rejected`: user rejected the recommendation or option.
- `deferred`: user chose to postpone the decision.
- `non-blocking follow-up`: recorded for later and not blocking the selected output target.
- `superseded`: replaced by another decision or changed context.

Closed decision states are `approved`, `delegated`, `rejected`, `deferred`, and `superseded`.

`review-ready recommendation` and `manual choice required` are not closed states. They can satisfy
manual-review readiness, but they do not authorize implementation handoff unless the user
explicitly requests or delegates that use.

## Concrete Recommendation Standard

A valid recommendation must choose one viable option, recommend a specific hybrid of options, or
explicitly state `No safe default`.

"The user must decide", "decide this before implementation", and "review-ready recommendation" are
not recommendations. Review-ready is a decision state. It does not replace the obligation to
recommend the option that best fits the intention contract when a safe recommendation exists.

For a concrete recommendation, include:

- recommended option;
- why this option best satisfies the intention, priority order, durable principles, and evaluation
  criteria;
- evidence that supports the recommendation;
- why serious alternatives were rejected or deferred;
- downstream effects on architecture, naming, docs, validators, runbooks, operations, external
  interfaces, or handoff;
- risks, mitigations, and residual uncertainty;
- exact approval wording the user can accept, change, delegate, reject, or defer.

The recommendation should be detailed enough that a later planner can tell why the option won, what
assumptions it depends on, and what would cause the decision to reopen.

## No Safe Default Standard

Use `No safe default` only when the best option genuinely depends on unresolved user intention,
authority, long-term scope appetite, cost or risk tolerance, compliance posture, business posture,
or another unresolved decision that cannot be inferred from the confirmed intention contract.

Do not use `No safe default` because the decision is hard, politically sensitive, or requires
manual approval.

A valid `No safe default` packet must include:

- explicit statement that no safe default exists;
- why the intention contract and evidence cannot prefer one option without inventing user intent;
- option-by-option comparison showing which intention or tradeoff would make each option correct;
- the minimum missing user judgment or upstream decision needed to choose safely;
- a provisional default if the user delegates authority under the stated priority order, or an
  explicit statement that delegation would still be unsafe;
- consequences of deferring the decision;
- exact user-choice prompt with the smallest complete set of options.

If one option is preferred by the current intention contract but still needs user approval, record
it as the concrete recommendation and mark it `review-ready recommendation` or
`manual choice required` when the remaining issue is an explicit user choice. Do not mark it
`No safe default`.

## Reviewer Gate

Use a reviewer gate for:

- `BLOCKER: yes` decisions;
- architecture, security, identity, naming, repository structure, product boundary, operating
  model, compliance, or cross-product decisions;
- hard-to-reverse choices;
- decisions that shape multiple future implementation plans or workflows.

Reviewer execution guidance:

- Use a separate reviewer agent when the environment permits it or the user explicitly asks for
  reviewer-agent execution.
- If a separate reviewer agent is not available, perform a documented reviewer pass in the main
  thread and state that no separate reviewer was used.
- Use high reasoning by default for important decisions.
- Use extra-high reasoning only for hard-to-reverse, security-critical, or broad cross-plan
  decisions.
- Do not hardcode a specific model version in managed workflow doctrine.

Give the reviewer a complete context packet:

- intention contract;
- decision inventory and dependency map;
- already approved, delegated, deferred, rejected, or manual-choice decisions;
- decisions blocked by the active decision;
- active decision packet;
- current options and recommendation;
- relevant evidence;
- adjacent domains or future growth pressure;
- what is out of scope.

Ask the reviewer to challenge:

- whether the intention was understood correctly;
- whether the decision is being answered at the right abstraction level;
- whether the decision question is framed too narrowly;
- whether important options are missing;
- whether the evidence is sufficient;
- whether the recommendation chooses a concrete option and follows from the criteria;
- whether a `No safe default` claim is justified by missing user intent rather than indecision;
- whether risks, tradeoffs, or downstream effects are underplayed.

The main agent must respond to reviewer critique by accepting and integrating it or rejecting it
with a reason.

For each decision that uses a reviewer gate, record a concise reviewer exchange summary before
treating the decision as `review-ready recommendation`, `manual choice required`, or closed.

The summary should include:

- decision ID and title;
- reviewer mode;
- recommendation sent to reviewer;
- reviewer objections, missing evidence, or reframing requests;
- main-agent response and recommendation changes;
- convergence result, remaining dissent, or user-choice point;
- residual risks and uncertainty after review.

Do not paste full reviewer transcripts by default.

## Convergence And Closure

Do not loop indefinitely.

Default convergence loop:

1. Main agent prepares the decision packet and recommendation.
2. Reviewer critiques when a reviewer gate is required.
3. Main agent revises or responds with rationale.
4. Reviewer performs one final check when needed.
5. Main agent records convergence or bounded disagreement.

A decision is ready for user action when:

- the recommendation is concrete and evidence-backed;
- `No safe default` is justified; or
- remaining disagreement is clear, bounded, and presented to the user.

A decision is closed only when:

- the intention contract is explicit enough for this decision;
- required evidence has been gathered and summarized;
- options have been compared against evaluation criteria;
- the recommendation chooses a concrete option and ties it to intention and priority order, or a
  `No safe default` packet meets the required standard;
- reviewer gate is complete or explicitly waived;
- reviewer exchange summary is recorded when a reviewer gate was used;
- dissent is resolved or escalated;
- the user approves, rejects, delegates, or defers the decision, unless prior explicit delegation
  applies;
- the discovery document records final status, rationale, alternatives, risks, affected IDs, and
  remaining blockers.

When a decision changes state, update the discovery document or ledger before moving to the next
decision.

## Document Update Requirements

After each material decision-state change, update:

- decision state;
- open question blocker state;
- rationale;
- alternatives considered;
- risks and mitigations;
- assumptions added or changed;
- affected requirements and non-functional requirements;
- next-step or handoff implications;
- reviewer exchange summary when a reviewer gate was used;
- revision notes for meaningful decision, scope, strategy, or handoff changes.

Do not add revision notes for typo fixes, formatting-only edits, or timestamp-only updates.

## Discovery Document Shape

Use the repository or workspace documentation placement rules when they exist.

Default discovery filename:

```text
<feature-or-topic-slug>_discovery_doc.md
```

Every Markdown discovery document starts with:

```text
Last updated: YYYY-MM-DDTHH:MM:SSZ (UTC)
Created: YYYY-MM-DD
Status: draft | active | completed | superseded | archived
```

Use the current UTC clock for `Last updated`. Preserve `Created` after the document is created.

Use sections that fit the domain, but ensure the document contains:

1. Problem Framing:
   - problem statement;
   - user intention;
   - goals;
   - non-goals;
   - priority order;
   - durable principles;
   - success criteria.
2. Context And Evidence:
   - essential context files or sources;
   - external references when used;
   - candidate domains and decision clusters when domain fit was uncertain;
   - relevant existing patterns, constraints, or operating facts.
3. Requirements And Constraints:
   - `FR-*`;
   - `NFR-*`;
   - workflows or operating scenarios when relevant;
   - data, integration, policy, platform, financial, regulatory, or organizational constraints.
4. Decision Ledger:
   - `D-*` overview table or list;
   - decision state, class, owner, dependencies, affected IDs, and closure criteria;
   - decision log per important `D-*` with options, recommendation, rationale, tradeoffs, state,
     reviewer summary when used, and next user action.
5. Open Questions And Assumptions:
   - `OQ-*` with `BLOCKER: yes/no`;
   - `A-*`;
   - owner and resolution path.
6. Risks And Validation:
   - `R-*`;
   - mitigation and detection;
   - targeted validation plan.
7. Manual Review Or Handoff:
   - include the section matching the selected output target.
8. Revision Notes:
   - meaningful decision, scope, strategy, or handoff changes.

## Manual Review Packet

When the output target is manual-review-ready, prepare a concise packet for the user.

Include:

- decisions already approved, delegated, rejected, deferred, or superseded;
- review-ready recommendations awaiting approval or correction;
- `No safe default` manual-choice packets;
- remaining non-blocking follow-ups;
- reviewer dissent or residual uncertainty that matters;
- exact user actions needed.

Do not call the document implementation-ready merely because the manual review packet is complete.

## Handoff Rules

Only create an implementation or next-step handoff when the selected output target requires it or
the user explicitly asks for it.

The handoff should include:

- approved or delegated decisions;
- MVP or first-step requirements;
- critical constraints;
- hard non-goals;
- remaining blockers and resolution paths;
- high-level epic outline only, if implementation planning is expected next;
- minimum validation gates;
- follow-ups that are explicitly non-blocking.

If unresolved blockers remain by user instruction, label the handoff as one of:

- `blocked handoff`: implementation planning cannot choose a single path yet.
- `conditional handoff`: implementation planning can proceed only under named assumptions.
- `blocker-resolution handoff`: the next step is to resolve the blocker, not implement the feature.

Implementation epics and checkbox tasks belong in a `*_implementation_doc.md`, not in the discovery
document.

If a closed decision changes after handoff, treat the handoff as stale and regenerate it before
drafting or updating an implementation plan.

## Plan Register Hook

When discovery work creates or materially updates a `*_discovery_doc.md`, maintain the local plan
register if the change affects plan identity, scope, decision state, lifecycle state, readiness,
parent or child relationships, dependencies, supersession, related plans, or review outcome.

Use `maintain-plan-register.md` for the procedure and `../reference/plan-register-format.md` for
entry shape. The sequence is:

1. Update `docs/repo/plans/plan-register.md` in the local repository.
2. When portfolio coordination is configured and the coordination home is available and safe to
   edit, update the configured coordination-home register.
3. When portfolio coordination is configured but the coordination home is unavailable or unsafe to
   edit, record pending sync in `docs/repo/plans/coordination-sync-pending.md` and continue the
   local discovery task.

Record creating or material-update session refs when a session ID is available. Do not block
discovery when no session ID is available.

Do not update the register for typos, formatting-only edits, timestamp-only edits, or mechanical
checklist progress that does not change lifecycle, relationships, readiness, or decision state.

## Domain Adaptation

Use the same discovery mechanics across domains, but translate domain-specific terms instead of
forcing every discovery into software language.

Examples:

- implementation handoff may mean validation handoff, execution handoff, sourcing handoff,
  operational handoff, or implementation handoff depending on the domain;
- technical decisions may be business, product, sourcing, regulatory, financial, operational,
  partner, market, research, repository, or policy decisions;
- repo evidence may be replaced or supplemented by supplier evidence, customer evidence, market
  evidence, regulatory evidence, financial assumptions, operational constraints, or external
  references;
- architecture options may be replaced by business-model, operating-model, channel, supplier,
  compliance, validation, or governance options.

Keep the invariant: decisions trace back to intention, evidence, alternatives, risk, downstream
effects, reviewer scrutiny when required, and explicit approval, rejection, delegation, or deferral
before they are closed.

## Checklist

- [ ] Select input state and output target.
- [ ] State the discovery status checkpoint for substantial discovery.
- [ ] If the user triggered Deep Discovery, default to manual-review-ready and use helper agents by
      default where useful.
- [ ] Do not mark Deep Discovery complete until its goal success criteria are satisfied or the user
      explicitly narrows the goal.
- [ ] Restate problem, outcome, constraints, and scope.
- [ ] Detect candidate domains and decision clusters when the request may span more than one
      domain.
- [ ] Establish or refresh the intention contract.
- [ ] Extract durable intention from user wording.
- [ ] Use the intention ladder when decisions expose unclear intent.
- [ ] Run targeted research before broad questions.
- [ ] List essential context files or sources.
- [ ] Create or update `FR-*`, `NFR-*`, `D-*`, `OQ-*`, `A-*`, and `R-*`.
- [ ] Mark every open question with `BLOCKER: yes/no`.
- [ ] Build or refresh the decision inventory.
- [ ] Classify decisions and add dependency metadata.
- [ ] Select one active decision.
- [ ] Prepare the decision packet.
- [ ] Gather required evidence.
- [ ] Compare options against criteria.
- [ ] Use helper agents for bounded research or review when useful.
- [ ] Run a bounded landscape scan when broad external options shape the active decision.
- [ ] Produce a concrete recommendation or justified `No safe default` packet.
- [ ] Run reviewer gate for important or blocker decisions.
- [ ] Record reviewer exchange summary when a reviewer gate was used.
- [ ] Mark the decision as review-ready, manual-choice-required, closed, or non-blocking.
- [ ] Update the discovery document or ledger before moving to the next decision.
- [ ] Continue until the selected output target's readiness gate is satisfied.
- [ ] Prepare a manual review packet or handoff only when the selected target requires it.
