Last updated: 2026-10-05T16:32:21Z (UTC)
Created: 2026-10-05
Status: draft

# Cross-Tool Agent Coordination Discovery

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.discovery.cross-tool-agent-coordination
  kind: discovery
  purpose: Establish proportionate commissioning, context, permission and result-return procedures for agents working across Codex and Claude CLI sessions.
  first_cataloged: 2026-10-05T16:32:21Z
  catalog_metadata_updated: 2026-10-05T16:32:21Z
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.claude-code-operating-guide
    - kind: related
      target: codeheart-operating-kit.implementation.proportionate-agent-workflows
```
<!-- END CODEHEART PLAN METADATA -->

## Problem And Intended Outcome

A coordinating agent should be able to commission another tool's agent, remain available for
discussion, receive its questions and results, and resume the same assignment. Consumers should
not have to invent scripts, restate authority at every routine step or move messages manually.

This discovery consolidates accepted commissioning requirements and a bounded consultation. It
does not activate implementation, change permissions or reopen the earlier Claude Code guide.
Target: draft-ready, with unresolved implementation-shaping evidence explicitly retained.

Priorities are faithful delegation and results, continued user interaction, few interruptions,
simple reusable operation, and proportionate cost. Keep one accountable coordinator and existing
plan/epic acceptance boundaries. Whole-plan execution may proceed through routine steps and
delegated reviews under sufficient authority; a tool's successful exit never constitutes acceptance.

Out of scope: a general orchestration platform, polling watchers, scheduled monitoring, a new
role registry, automatic permission escalation, desktop transcript injection, and universal
cross-tool compatibility claims. Existing assignments are not automatically restarted or archived.

## Owners And Existing Routes

The Kit owns reusable procedure, commissioning templates, tool-specific instructions and any
subsequently justified helper. Consumer governance owns appointments, preferred models, local
authority and private strategy. The owning plan/log records the actual implementer or reviewer,
session locator, scoped authority, work/evidence references and accepted disposition. Session IDs
locate execution; they are not durable role identities. Root agent instructions remain short routers.

Existing source routes:

- [Generic coordination](../../../../components/agent-interface/managed/reference/agent-task-coordination.md)
- [Codex operations](../../../../components/agent-interface/managed/reference/codex-task-operations.md)
- [Claude Code operations](../../../../components/agent-interface/managed/reference/claude-code-task-operations.md)
- [Recipe maturity](../../../../components/agent-interface/managed/reference/operational-recipe-maturity.md)
- [Script promotion](../../../../components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md)

Future delivery is routing-bearing and recipe-bearing. This draft changes only producer planning
documents: no managed instructions, safety policy, installation behavior or consumer paths change.
Consumer-impact classification and release/adoption requirements must be revisited for actual
implementation, particularly if permission guidance or executable assets are introduced.

## Evidence And Its Limits

The commissioning pilot supplied these observations; private transcripts, identities, resource
names, local paths and raw logs stay with its owner and are not public Kit content.

- A Codex relay launched a fresh Claude CLI consultation, captured its response and resumed the
  same session for follow-up. A separate native Codex message-delivery test succeeded. It waited
  for process output rather than repeatedly querying saved chats.
- CLI resume of an app-owned Codex conversation failed with an active-writer conflict. Saved
  conversation access is not a supported live desktop message channel. Use a fresh CLI-only
  session for this pattern, with one active driver; do not override session ownership.
- A real review and correction round trip worked. Whole-plan implementation, fresh director
  onboarding, reverse Claude-led delivery and unattended merge under a selected profile remain
  unproven. An available CLI in both directions does not establish symmetric notification.
- The relay also researched evidence and interpreted results, while the coordinator repeated some
  investigation. Usage showed substantial repeated context processing. This supports a narrower
  transport role and compact context; it does not establish future cost or model reliability.
- Recorded merge attempts were denied as `Merge Without Review`. The inspected Claude 2.1.286
  default is a soft rule expecting human approval; independent agent review does not itself meet
  that default. Teammate text does not meet soft-rule consent requirements. A direct human merge
  instruction was also present in one denied case. These findings identify a concrete mismatch,
  not a universal hard ban on delegation or proof of every historical policy input.
- Current documented `autoMode.allow` prose can provide scoped exceptions to soft rules. Its
  settings come from user, managed or invocation configuration, not repository project settings.
  Actual acceptance of the intended workflow is untested. Normal command allow patterns do not
  establish plan scope or independent review. A separate agent identity is not human approval.
  Protected auto-merge is an alternative where suitable infrastructure already exists, but
  enqueue-time commit matching alone does not prove the final merged candidate after updates.

Official references checked on 2026-10-05:
[Claude permission modes](https://code.claude.com/docs/en/permission-modes),
[auto-mode configuration](https://code.claude.com/docs/en/auto-mode-config),
[permission rules](https://code.claude.com/docs/en/permissions), and
[Codex subagents](https://learn.chatgpt.com/docs/agent-configuration/subagents).
Recheck relevant installed versions when selecting the pilot. Settings observations cover only
the inspected files and current configuration extracts, not every desktop/server/historical policy.

## Requirements And Working Scenarios

**FR-1 — Commissioning.** Supply the question/outcome, role, repository/worktree, relevant records,
settled decisions, authority and limits, review owner and finish line. Faithfully attribute existing
human approval; do not represent a relay or advisor's new wording as a new human decision.

**FR-2 — Context by role.** Use a current brief plus changes and exact source references; relevant
saved conversations supply rationale and intent, not automatic live context or replacement truth.

| Function | Relevant context and responsibility |
| --- | --- |
| Organization advisor | Company direction and relevant coordinating discussion; challenge substantive strategy and cross-program choices. |
| Program advisor | Program strategy, related plans, constraints and relevant director discussion; advise on consequential program/plan decisions. |
| Implementer | Approved whole plan, dependencies, authority and acceptance criteria; deliver the assignment and its corrections. |
| Independent reviewer | Candidate, requirements, authority limits and acceptance criteria; reach an independent judgment. An advisor who co-authored should not be the sole independent reviewer. |
| Temporary relay | Ready message, session/start instructions, paths and return destination; transport faithfully without substantive research or judgment. |

These are functions, not a requirement for five standing agents or a review on every commit.
Reuse an advisor's session while useful; provide a handover to a fresh session when context becomes
stale or unwieldy. Reuse an implementer's session for its coherent assignment and corrections.

**FR-3 — Proportionate preparation and review.** Prior research by the coordinator resolves a
commissioning uncertainty, such as the correct plan or current work state. It should not pre-solve
the delegated question. Investigation and interpretation belong to the commissioned substantive
role. The coordinator assesses the result and verifies consequential gaps or questionable claims;
routine duplication of the investigation is not the default. Small specified extraction can be
mechanical preparation, but open-ended searching is a substantive assignment.

**FR-4 — Lean return contract.** Return assignment identity, response/process status, unchanged
short reply or exact full-response reference, session locator, and relevant failure details. Forward
questions unchanged. A response-available notice is not a claim of implementation completion.
Check process/session/output integrity without adding a second substantive digest. Preserve the
original result. A maintained helper should perform deterministic mechanics where proportionate.

**FR-5 — Conversation availability.** Launch the delegated work and return control to the user
instead of holding the main conversation in a repeated wait/research loop. The coordinator can
discuss independent matters while the relay waits. Review results when received; actively wait
only when the next decision depends on them. Verify actual background/notification behavior in
the host before promising unattended report-back after a turn ends.

**FR-6 — Lifecycle and failures.** Keep one writer per session. Return questions, permission
denials, interruption and failure with preserved session/work/evidence. Verify uncertain effects
before any retry. Do not infer plan completion from process exit, auto-retry ambiguous mutations,
start follow-on scope, or escalate model/permissions without the responsible owner's decision.

**FR-7 — Reuse and direction.** Generic doctrine must permit either tool to coordinate. Tool
recipes must name tested directions and limitations. Fresh agents should discover the appropriate
procedure and appointed advisor from normal repository/Program routes without manual coaching.

**NFR-1 — Economy.** Give relays compact, fresh context and configurable model/reasoning selection.
Model preference belongs in consumer configuration, not a hardcoded public model version. Record
actual model, input/cache/output usage and elapsed time where available; do not equate aggregate
cached processing with new content or assume quoted list cost equals subscription billing.

**NFR-2 — Boundaries.** Preserve real repository checks/reviews, credential scope, hard permission
boundaries and human-controlled decisions. A denied outcome is not rerouted through another agent.
No blanket full-access default is selected. Keep secrets and raw private evidence out of the Kit.

## Decision Ledger

These are accepted discovery inputs from the commissioning owner. They are not a completed
independent review of this consolidated document or authorization to ship implementation.

| Decision | State, rationale and consequence |
| --- | --- |
| D-1 Reusable procedure in the Kit | Approved direction. Consumer-specific appointments/preferences stay in consumer governance; execution stays with owning plans. Extends existing guides rather than inventing a registry. |
| D-2 Distinct substantive and relay responsibilities | Approved direction. Consultant investigates, director assesses, relay transports. Some mechanical preparation is allowed; no default parallel duplicate research. |
| D-3 Role-matched context and durable continuity | Approved direction. Brief plus relevant records/conversation evidence, stable role appointments with replaceable session locators, whole-plan implementer assignments. |
| D-4 Main conversation remains available | Approved direction. Dispatch and yield by default; verify the host's completion/attention route. No polling watcher is selected. |
| D-5 Permission approach | Approved investigation direction; profile unresolved. First evaluate a small supported invocation-scoped policy for approved normal merges after independent review. Existing protected auto-merge may be used if already suitable; extra identities/protection machinery are not prerequisites. |
| D-6 Model and tooling | Approved intent for a capable but narrowly instructed relay with compact context. Exact consumer model and compatible reasoning are resolved at commissioning. A small maintained launcher is a candidate, not yet selected or implemented. |
| D-7 Proof proportionality | Approved direction. One meaningful upcoming delivery and onboarding exercise, with relevant failure/question return. No arbitrary repeated-delivery gate or exhaustive framework. Expand evidence only for an actual unresolved risk. |

The permission consultation was reviewed by the coordinator. It corrected overclaims about agent
approval being human approval, queued-merge candidate binding, settings coverage, and invented
requirements for verbatim approval copying, repeated trials and re-approval of existing delegation.
Residual uncertainty is captured below. No separate review of the complete new discovery has yet
been commissioned; this remains a draft, not implementation-handoff-ready.

## Open Questions, Assumptions And Risks

| Open question | Blocker and resolving evidence |
| --- | --- |
| OQ-1 Which exact permission profile supports the authorized workflow? | BLOCKER: yes for unattended merge support. Coordinator selects a concrete scoped policy and permitted pilot effects; verify it with the installed CLI. Preserve defaults, deny/ask rules and real server-required checks. Merge findings do not qualify releases, deployments or other actions. |
| OQ-2 What should ship as a launcher? | BLOCKER: yes for executable scope. Assess the repeated invocation/result mechanics, supported hosts, ownership, minimal inputs/output and interruption behavior. Prefer a small maintained helper if it removes recurring scripts; do not design an orchestration platform. |
| OQ-3 Does the lean relay preserve delivery while the coordinator yields? | BLOCKER: yes for a background-delivery claim. Pilot response, question, failed invocation and continuation with compact context and the selected model. Capture integrity and usage. Stop on unreliable delivery or a genuine permission boundary. |
| OQ-4 What reverse workflow is actually supported? | BLOCKER: no for initially scoped Codex-led support; yes for advertising bidirectional support. Verify Claude-led commissioning and return into an available Codex execution surface without an app-writer conflict. |
| OQ-5 Does a fresh director discover and use the procedure? | BLOCKER: yes for onboarding acceptance. Exercise installed instructions and a Program appointment/absence case; account for tool-specific instruction-loading precedence. Do not assume an existing warm session proves onboarding. |

Assumptions: A-1 selected CLIs are installed/authenticated and available; check at launch without
silently installing or changing versions. A-2 the coordinator can receive a native relay result;
verify for the host and idle/active states used. A-3 generic direction can be symmetric while
initial tested transport is asymmetric. A-4 existing Program/Plan records are sufficient unless a
concrete pilot shows otherwise. An unavailable requested model is a commissioning gap, not license
to silently substitute another model or reasoning level.

Risks: R-1 stale or oversized context loses intent or wastes processing; use scoped briefs and
measured context. R-2 policy prose is classifier-interpreted; distinguish supported configuration
from demonstrated reliability. R-3 duplicate drivers or retries damage work; preserve ownership and
inspect state. R-4 a relay summary alters advice; retain unchanged results for coordinator review.
R-5 completion notices fail when the parent yields; verify the actual lifecycle rather than add a
watcher by assumption. R-6 broad allow patterns or inaccurate human attribution weaken authority;
keep scope explicit and do not apply a profile as part of this documentation update.

## Next Evidence And Handoff State

Next work is bounded discovery: prepare the concrete permission/profile candidate and lean relay
pilot, identify the smallest reusable invocation mechanics, and resolve the applicable questions
above. The actual pilot target/effects and exact consumer model remain to be selected. Its output
should support a coherent implementation capability scope before epics are drafted.

This is not an implementation plan. No permission policy, launcher, agent appointment, release or
consumer adoption is activated. Public source review and normal release/adoption remain required
when a delivery is commissioned. The earlier guide and its remaining obligations retain their own
authority. Draft publication alone would not authorize those effects.

## Revision Notes

- 2026-10-05: consolidated accepted cross-tool commissioning, role context, lean relay, asynchronous
  interaction, durability and proportionate-evidence direction; recorded permission consultation
  findings and the remaining implementation-shaping questions. Private model choices and pilot
  evidence remain with the commissioning owner.
