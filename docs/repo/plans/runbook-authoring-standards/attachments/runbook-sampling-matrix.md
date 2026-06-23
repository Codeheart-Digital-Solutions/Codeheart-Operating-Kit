Last updated: 2026-06-23T12:02:14Z (UTC)

# Runbook Sampling Matrix

This attachment supports
`../runbook-authoring-standards_discovery_doc.md`.

It records public-core-safe sampling notes for existing runbooks. Consumer and module examples are
generalized so private tenant, customer, account, and local-environment details do not enter the
public Operating Kit repository.

## Sampling Method

The sampling focused on whether existing runbooks give either humans or agents enough structure to
complete the intended workflow without inventing missing decisions.

Evaluation lenses:

- audience clarity;
- user-facing script quality;
- agent-facing execution detail;
- required inputs and source of truth;
- approval gates and stop conditions;
- evidence and validation expectations;
- public-core safety;
- suitability as a standard test case.

`Strong` means the runbook should influence the future standard. `Test case` means the runbook is
useful for proving the standard because it exposes a gap or boundary question. A test case is not
necessarily a poor runbook; it may simply predate the proposed standard or intentionally stay terse.

## Sample Matrix

| Sample | Candidate Classification | Strong Signals | Gaps Or Test-Case Value | Standard Implication |
| --- | --- | --- | --- | --- |
| `components/agent-interface/managed/runbooks/conduct-first-run-onboarding.md` | Strong human-facing exemplar | Ordered onboarding procedure, visible choices, user-owned decisions, neutral examples, write boundary, native-capability and update-check handoff. | Does not yet declare an explicit audience metadata field. Evidence/validation is implicit in the setup flow rather than a named section. | Human-facing runbooks should include exact or near-exact chat flow, decision ownership, neutral examples, and write approval boundaries. |
| `components/agent-interface/managed/reference/onboarding-context-contract.md` | Strong hybrid contract exemplar | Separates first prompt, agent contract, ordered context, copy blocks, inspection modes, storage boundaries, and non-interactive rules. | It is a reference contract rather than a runbook, so the future standard should avoid forcing every runbook into this much detail. | Hybrid runbooks need a clean split between user-visible copy and operator-only contract details. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Strong agent-facing planning exemplar | Trigger, inputs, file rules, required structure, decision-quality rules, blocker handling, checklist rules, quality gate, and fresh-implementer test. | Highly specialized to implementation planning; not a universal runbook template. | Reuse the fresh-implementer test for agent-facing runbooks: a future agent should not need to invent workflow, commands, evidence, or validation. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | Strong agent-facing review exemplar | Review stance, required inputs, review areas, severity model, output order, final readiness statement. | Could gain audience metadata and an explicit evidence/output contract. | Review runbooks should define issue order, severity, and output shape, not just ask the agent to "review." |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Strong agent-facing execution exemplar | Required read order, lifecycle state rules, safe defaults, per-epic flow, review gate, execution log shape, register hook, and final summary contract. | Long and dense; future standards should not make small runbooks this heavy by default. | Complex execution runbooks need lifecycle rules, authority gates, evidence expectations, and review gates. |
| `components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Strong technical operations exemplar | Trigger conditions, required references, local and portfolio procedures, pending-sync fallback, session-reference safety, and safety rules. | Complexity may be high for casual maintainers; could benefit from a short input/output summary at the top. | High-risk coordination runbooks should include fallback behavior and precise stop/safety rules. |
| `components/structure-governance/managed/runbooks/change-documentation-placement.md` | Good agent-facing operational baseline | Preflight, owner-selection order, procedure, archive handling, validation, and the plan-scoped `attachments/` rule. | Missing explicit audience classification and input/output/evidence blocks. | This is a good minimal operational runbook shape: preflight, owner choice, procedure, validation. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer-facing test case | Very concise ordered procedure, public-core safety, placement and impact-classification routing, stop conditions. | Too terse to show expected evidence, release-note recording, or validation output. | The standard should allow short maintainer runbooks, but require extra evidence fields when blast radius rises. |
| `docs/repo/runbooks/release-operating-kit.md` | Maintainer-facing high-risk test case | Linear release procedure, validation sequence, asset/checksum steps, publication stop conditions. | It names evidence needs but does not structure the expected evidence record. Authority gates are embedded in prose. | High-risk maintainer runbooks should have explicit authority, evidence, rollback/stop, and artifact-record sections. |
| `docs/repo/runbooks/promote-consumer-change.md` | Maintainer-facing promotion test case | Clear promotion flow, sanitization requirement, placement/impact routing, consumer-specific stop condition. | Could use an explicit sanitization checklist and examples of what must stay in the consumer repo. | Promotion runbooks need a reusable public-safety and ownership-boundary checklist. |
| `docs/repo/runbooks/triage-kit-feedback.md` | Maintainer-facing lifecycle exemplar | Scope boundary, lifecycle-state table, discovery/implementation handoff, public-disclosure response, release/sync notes. | Not a human-facing issue-response script; triage comments are described but not scripted. | State-machine runbooks benefit from lifecycle tables and handoff rules. |
| Sanitized consumer/module M365 onboarding runbook | Hybrid onboarding test case | After hardening, includes conversation rules, user-facing script, help finding tenant details, required input record, operator notes, execution flow, and stop conditions. | The live failure showed that technical correctness alone did not ensure a good nontechnical onboarding experience. The sample should be used to test user-copy separation. | Hybrid setup runbooks should begin with plain-language intent, ask one decision per turn, hide internals in first-turn copy, and provide help for hard-to-find inputs. |
| Sanitized consumer/module M365 operations runbook | Agent-facing technical recipe test case | Routing table, first-response rules, common execution discipline, required fields, tool lanes, command sequences, approval packets, stop conditions, evidence section. | Large recipe surface may eventually need sub-runbooks or a generated action catalog if reuse patterns stabilize. | Agent-facing technical runbooks should name the lane, required fields, preflight, approval boundary, commands, stop conditions, and evidence. |
| Sanitized consumer business decision runbooks | Domain operations test case | Clear folder convention, agenda/notes/decision/todo sections, canonicalization rules, source separation. | Strong domain workflow but not a chat script. It should not be judged by human-facing onboarding standards. | Audience classification prevents over-applying user-facing script requirements to domain workflow runbooks. |

## Cross-Sample Findings

### Finding 1 - Audience Metadata Is Missing Even In Strong Runbooks

Most sampled runbooks communicate their audience through structure and wording rather than an
explicit field. The future standard should add a lightweight audience declaration for new or
materially changed runbooks instead of retrofitting every existing file immediately.

### Finding 2 - Human-Facing Quality Depends On Scripted Conversation, Not Just Safety Rules

Strong human-facing examples give the agent the words, order, visible choices, and stop points.
Weak or pre-standard onboarding can still be technically safe while producing a poor interaction
because it exposes local files, internal state, or too many questions before the user understands
the goal.

### Finding 3 - Agent-Facing Quality Depends On Execution Specificity

Strong agent-facing examples name the trigger, required inputs, read order, execution path, stop
conditions, evidence, and validation. The planning workflows show the clearest version of this
through the fresh-implementer test.

### Finding 4 - Maintainer Runbooks Need A Low-Bureaucracy Path

Short maintainer runbooks are useful when they govern experienced maintainers and have low
ambiguity. The future standard should not force them into full onboarding-style templates. The
threshold for extra structure should rise with blast radius, release authority, security exposure,
or repeatability.

### Finding 5 - Hybrid Runbooks Need Hard Separation

The M365 onboarding example shows why hybrid runbooks should separate:

- user-facing flow;
- operator notes;
- execution path;
- stop conditions;
- evidence and validation.

Without that separation, agents tend to mix internal caveats, local setup details, and technical
questionnaires into the user's first experience.

## Candidate Standard Test Set

Use this compact set when drafting and validating the future runbook authoring standard:

| Test Purpose | Candidate Sample | Expected Use |
| --- | --- | --- |
| Human-facing exemplar preservation | `conduct-first-run-onboarding.md` | Ensure the standard preserves plain-language guided setup quality. |
| Hybrid contract preservation | `onboarding-context-contract.md` | Ensure user copy and agent contract can coexist without mixing. |
| Agent-facing execution specificity | `draft-implementation-plan.md` and `execute-implementation-plan.md` | Ensure the standard captures fresh-agent execution readiness. |
| Minimal maintainer runbook threshold | `change-operating-kit.md` | Ensure the standard does not over-template low-ambiguity maintainer procedures. |
| High-risk maintainer evidence | `release-operating-kit.md` | Ensure release authority, evidence, and stop conditions are explicit enough. |
| Promotion and sanitization boundary | `promote-consumer-change.md` | Ensure public-core promotion standards are tested. |
| Hybrid UX failure/recovery | Sanitized M365 onboarding sample | Ensure nontechnical onboarding flow is tested against a real failure mode. |
| Technical recipe execution | Sanitized M365 operations sample | Ensure broad agent-facing operation recipes stay executable. |

## Implementation Handoff Notes

The later implementation plan should:

1. Turn this matrix into acceptance criteria for the managed runbook authoring standard.
2. Avoid actively retrofitting consumer or module runbooks in the first implementation.
3. Use the sampled runbooks as review fixtures to check whether the standard is practical.
4. Add a review checklist that can catch missing audience classification, mixed user/operator
   content, vague execution paths, missing approval gates, and missing evidence expectations.
