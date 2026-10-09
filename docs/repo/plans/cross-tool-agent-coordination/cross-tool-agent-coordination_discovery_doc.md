Last updated: 2026-10-09T15:53:13Z (UTC)
Created: 2026-10-05
Status: completed

# Cross-Tool Agent Coordination Discovery

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.discovery.cross-tool-agent-coordination
  kind: discovery
  purpose: Establish proportionate commissioning, context, permission and result-return procedures for agents working across Codex and Claude CLI sessions.
  first_cataloged: 2026-10-05T16:32:21Z
  catalog_metadata_updated: 2026-10-09T14:56:20Z
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.cross-tool-agent-coordination
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

This discovery consolidates accepted commissioning requirements, consultation and bounded pilots.
On 2026-10-09 the user approved consolidating the findings, settling the small helper's scope and
preparing the implementation/onboarding work. That delegates the bounded capability detail below
for planning; it does not activate implementation, change live permissions or reopen the earlier
Claude Code guide. Readiness: implementation-handoff-ready. The related implementation plan remains
draft for review. Completion here means discovery handoff, not shipped coordination capability.

Priorities are faithful delegation and results, continued user interaction, few interruptions,
simple reusable operation, and proportionate cost. Keep one accountable coordinator and existing
plan/epic acceptance boundaries. Whole-plan execution may proceed through routine steps and
delegated reviews under sufficient authority; a tool's successful exit never constitutes acceptance.

Out of scope: a general orchestration platform, polling watchers, scheduled monitoring, a new
role registry, automatic permission escalation, desktop transcript injection, and universal
cross-tool compatibility claims. Existing assignments are not automatically restarted or archived.

## Owners And Existing Routes

The Kit owns reusable procedure, commissioning templates, tool-specific instructions and any
selected invocation helper. Consumer operational references own shared relay routing, preferred
models and local operating choices; Program governance owns its appointments and private strategy.
The owning plan/log records the actual implementer or reviewer,
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
Implementation includes executable mechanics and permission guidance. It therefore requires broad
candidate acceptance, explicit safety-policy review and release/adoption notes; it is not eligible
for instruction-only release qualification.

## Evidence And Its Limits

Original responses, exact sessions, permissions and process/delivery records were inspected by the
coordinator. They remain private consumer evidence. The following is the public-safe conclusion,
not a copy of transcripts, local paths, account details or live chat identities.

| Observation | What it establishes | Limit |
| --- | --- | --- |
| Fresh Claude CLI consultation and same-session follow-up | Prepared briefs and original results can be transported without the relay doing research. | Use CLI-owned sessions with one driver; app-owned sessions are not part of this recipe. |
| Documentation PR 24 merged normally on 2026-10-06 | One reviewed candidate passed validation and merged under an invocation-scoped ordinary-merge permission profile with defaults and explicit restrictions retained. | No release, deployment or blanket bypass was qualified; one success is not classifier reliability proof. |
| Earlier narrow merge exception was denied | Tool permission and workflow authorization are separate. A rule requiring the classifier to prove review/CI from tool output was unsuitable in the tested setup. | The replacement allows ordinary task-repository merges; workflow still checks authority, review and readiness. |
| Native child completion after the parent ended its turn | Output was retained, but native completion did not start a new idle-parent turn in that test. | Active-parent collection is not idle notification. |
| Native app-message restrictions | Sending to a native ancestor and sending direct app input to a native-subagent recipient were rejected. | Do not bypass these boundaries or propose native mock directors as ordinary-chat substitutes. |
| Independent ordinary relay, tested 2026-10-09 | A completion notice started a new turn in an idle ordinary coordinator without user prompting. | Observed host/tool behavior, not a portable notification guarantee. |
| Clarification round trip | Question and final response reached the coordinator; its answer resumed the same Claude session. | Relay forwarded; coordinator made the choice. |
| Concurrent relay worker assignments | Separate sessions, output and ordinary return destinations stayed isolated. | Worker assignments overlapped; tiny model processes happened sequentially. This does not prove simultaneous model capacity. |
| Two ordinary requesting chats, same active relay turn | Both incoming requests survived the overlap; each received its correct result once through a distinct worker. | One bounded interleaving, not durable queueing or exactly-once delivery. |
| Mistyped chat identifier in an earlier fixture | Tool rejection exposed a transport-only copy error. Corrected structured identifiers worked. | Supports loading destinations from records instead of retyping them, not automatic retry of uncertain sends. |

The original merge pilot's delayed notification was a real failure of timely reporting, not a
missing result. The independent-relay tests resolve that observed topology problem for the selected
host. Failed sends retain their response and error; no scheduler or continuous watcher was added.

One-off launch scripts, repeated request/result assembly and the identifier error justify a small
maintained invocation helper. No helper, installed reusable procedure or fresh-director adoption
has yet shipped. Those are implementation outcomes; reverse-direction support remains unproven.

## Requirements And Working Scenarios

**FR-1 — Commissioning.** Supply the question/outcome, role, repository/worktree, relevant records,
settled decisions, authority and limits, review owner and finish line. Faithfully attribute existing
human approval; do not represent a relay or advisor's new wording as a new human decision.
Some CLIs receive a relayed brief in the user-message role. That transport role does not establish
human authorship or grant additional authority, even if it influences the tool's interpretation.
Attribute authority to the owning record or actual user instruction and preserve its limits.
Permission exceptions come from an authorized profile decision, not self-authorizing brief text.

**FR-2 — Context by role.** Use a current brief plus changes and exact source references; relevant
saved conversations supply rationale and intent, not automatic live context or replacement truth.

| Function | Relevant context and responsibility |
| --- | --- |
| Organization advisor | Company direction and relevant coordinating discussion; challenge substantive strategy and cross-program choices. |
| Program advisor | Program strategy, related plans, constraints and relevant director discussion; advise on consequential program/plan decisions. |
| Implementer | Approved whole plan, dependencies, authority and acceptance criteria; deliver the assignment and its corrections. |
| Independent reviewer | Candidate, requirements, relevant design rationale/evidence, authority limits and acceptance criteria; reach an independent judgment. An advisor who co-authored should not be the sole independent reviewer. |
| Relay or temporary transport worker | Ready message, session/start instructions, paths and exact return destination; transport faithfully without substantive research or judgment. |

These are functions, not a requirement for five standing agents or a review on every commit.
Reuse an advisor's session while useful; provide a handover to a fresh session when context becomes
stale or unwieldy. Reuse an implementer's session for its coherent assignment and corrections.
Independence depends on contribution and an independent assessment, not session freshness alone.
Do not withhold relevant rationale to simulate independence. A durable advisor can challenge a
discovery without becoming its independent acceptance reviewer; replacing a relay does not require
replacing the substantive agent's session.

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
Where supported, allocate and retain a new session's locator before launch, outside the relay's
conversation. Retain the existing locator on resume. Recovery still requires checking for an
active process/writer and uncertain effects; knowing the locator alone does not make resume safe.
Determine how approval prompts, ask rules and repeated denials behave in the selected headless
invocation. A denied/unanswerable prompt or failed invocation must produce a returned blocker or
execution failure with the original evidence, rather than a silent stop or a success claim. The
substantive agent owns its question/blocker text; the relay forwards it and execution metadata.

**FR-6a — Strict relay boundary.** A relay may check that the expected process/session ran and
that output exists. It must not investigate the subject, prepare alternatives, answer substantive
questions, interpret findings, summarize recommendations, accept work or start follow-on scope.
Forward original questions and failures to the designated coordinator. Handling concurrent
assignments does not expand this role. A relay receives a compact prepared envelope, not the whole
coordinator conversation. The coordinator dispatches and yields, then assesses the original reply;
it does not routinely duplicate the delegated investigation while waiting.

**FR-7 — Reuse and direction.** Generic doctrine must permit either tool to coordinate. Tool
recipes must name tested directions and limitations. Fresh agents should discover the appropriate
procedure and appointed advisor from normal repository/Program routes without manual coaching.

**FR-8 — Human transparency.** For consequential reviews, the coordinator exposes the core
findings, their significance and accepted/challenged/deferred corrections in readable language,
with an original-response reference when useful. Verdict alone is not enough. This does not turn
the relay into a summarizer or require a new report for every message.

**FR-9 — Optional first use.** Normal planning/implementation continues with the current tool
unless the assignment or applicable accepted preference selects another arrangement. Claude is
optional, not a Kit prerequisite. The coordinator guides requested CLI installation/sign-in and
explicit relay setup before commissioning, explains account/usage requirements, and offers an
explicit native alternative or pause when unavailable. No silent substitution, repeated offers,
blanket-access requirement or credentials in repository records. Generic defaults are managed;
agreed consumer choices are created from an opt-in starter and preserved on Kit upgrades.

**NFR-0 — Minimal records.** The CLI's original final reply is sufficient review material, whether
read from its conversation or captured invocation result. A generated response.md is an optional
convenience view, not an extra authored deliverable. Keep minimum request/session/process/delivery
evidence ignored locally; preserve partial output and active recovery data. Durable accepted
conclusions belong in existing plans/logs. Eligible generated duplicates/logs may be removed through
ordinary authorized cleanup after acceptance and recovery needs end; no retention service or
automatic deletion is required.

**NFR-1 — Economy.** Give relays compact, fresh context and configurable model/reasoning selection.
Model preference belongs in consumer configuration, not a hardcoded public model version. Record
actual model, input/cache/output usage and elapsed time where available; do not equate aggregate
cached processing with new content or assume quoted list cost equals subscription billing.

**NFR-2 — Boundaries.** Preserve real repository checks/reviews, credential scope, hard permission
boundaries and human-controlled decisions. A denied outcome is not rerouted through another agent.
No blanket full-access default is selected. Keep secrets and raw private evidence out of the Kit.
The human policy owner, or a coordinator with explicit delegation covering the configuration and
its effects, authorizes the concrete permission profile. The coordinator records the exact profile,
scope and authority in the owning assignment/plan. Approval to investigate an approach is not
approval to apply an unspecified exception. This names the existing decision owner, not a new
per-action approval gate. Preserve applicable defaults and enforced deny/ask rules.

## Decision Ledger

| Decision | State, rationale and consequence |
| --- | --- |
| D-1 Reusable procedure in the Kit | Approved. Extend existing task coordination and tool guides; preserve one accountable coordinator and plan/epic acceptance. |
| D-2 Distinct substantive and relay responsibilities | Approved and reaffirmed 2026-10-09. Substantive agent investigates, coordinator assesses, relay transports without parallel thinking. |
| D-3 Role-matched context and durable continuity | Approved. Brief plus relevant records; advisors can persist, implementers retain coherent whole-plan sessions, relay identity is replaceable. |
| D-4 Main conversation remains available | Approved. Dispatch and yield; use the tested independent ordinary relay and supported native messages. No watcher, cron or polling service. |
| D-5 Permission approach | Approved direction with observed pilot. Invocation-scoped ordinary-merge allowance may support an authorized task; preserve defaults, explicit restrictions, actual checks and review. Never use blanket bypass, auto-escalate, or turn denial into a handoff to another agent to evade it. |
| D-6 Bounded maintained helper | Selected under the user's 2026-10-09 delegation to settle precise scope. Add one invocation command to the existing Go CLI, using standard libraries. It launches/captures one assigned Claude CLI turn; it does not message Codex, schedule work or interpret results. |
| D-7 Proportionate evidence | Approved. Reuse the bounded pilot conclusions; validate changed mechanics and one fresh-director onboarding flow. No repeated experiments merely to increase confidence. |
| D-8 Durable local placement | Accepted direction. Consumer-owned docs/repo/reference/agent-coordination.md holds operational routing; root instructions link there. Owning Programs retain appointments and improvement history, not the sole location of a shared relay service. |
| D-9 Reuse existing home identity | Preserve portfolio identity in .codeheart/kit.config.yaml. Home ID is neither a filesystem path nor a chat address. Add ordinary reference links for navigation; do not add unsupported config fields or a second registry. |
| D-10 Explicit initial support | Generic responsibilities are direction-neutral. Initial executable recipe is Codex desktop coordinating CLI-only Claude on a verified host. Reverse direction and other hosts are deferred until qualified. |
| D-11 Native default and guided opt-in | Approved 2026-10-09. Current-tool operation works without optional account, relay or local reference. Requested cross-tool setup handles CLI/login, non-secret scoped preferences and explicitly authorized relay creation. Managed guidance owns defaults; optional consumer-owned records retain local choices without sync overwrite. |

## Selected Helper And Communication Shape

Target command: codeheart-operating-kit coordination invoke-claude --request <request.json>.
This is a proposed new command, not an available command. Reuse the existing compiled Go CLI and
its release/install path so users do not need copied scripts or another runtime. Repeated
commissions, installed distribution and compatibility justify this narrow L3 command under the
script-promotion standard; do not build a separate CLI, framework or general agent API.

The coordinator writes the complete request and brief, including approved local defaults; the
relay passes their path unchanged. Missing fields are blockers, not relay choices. Before dispatch,
the coordinator knows the attempt directory. The request identifies the assignment and attempt,
exact ordinary return chat/host,
working directory, brief file, CLI executable, model, session mode/ID, approved permission and
tool configuration, and ignored local evidence/state root. A new session ID is allocated and
persisted before launch; a resume requires the retained CLI-owned session. Arguments are structured
and prompts go through stdin; no shell interpolation, implicit mode changes or arbitrary retry.
Optional execution limits come from the assignment; the helper adds no short whole-plan deadline.

The command performs deterministic validation, prevents a second participating writer to that
session in the shared host-local state root, launches once with child stdout/stderr connected
directly to private files (not helper-owned pipes), preserves partial output on interruption and
the original reply, and writes a small result and native-message argument file. That file contains
the exact approved recipient and original short reply or full-response reference; the relay
passes it to the exposed app tool without retyping identifiers. The helper cannot authenticate to
or call a private Codex app API. Its successful exit means response capture, not task acceptance.

Artifacts are private ignored runtime evidence. Normal output contains only status and locators,
not raw transcripts, prompts, secrets or reasoning. Record model/usage when supplied by the CLI;
unknown usage stays unknown. Do not scrape private internal reasoning from session storage.
On denial, mismatched session, timeout, interrupted execution or malformed output, retain evidence
and report uncertainty. Inspect before resuming; no automatic replay, force-unlock or mode fallback.
A local lock does not prove that an external driver or another host is absent.

One independent ordinary relay chat can serve multiple ordinary directors. It dispatches compact
assignments to temporary transport workers and returns control. Workers address the original
requesting ordinary chat directly, outside the worker's ancestor hierarchy. One driver per Claude
session still applies. Without an authorized valid return target, report a blocker rather than
guessing a similarly named chat. A failed or uncertain send preserves its result and delivery
state in the known attempt directory; the coordinator checks it on demand when a result is needed
or the user asks. No exactly-once or eventual-delivery guarantee is claimed.

The helper uses a fixed documented flag set and records CLI version, returning invocation errors
without compatibility-matrix machinery or fallback. Whole-plan implementer use remains required.
Validate normal tool yield/resume and interruption behavior on the host; disclose actual limits
and unknowns without promising uninterrupted hours or silently narrowing to consultations.

## Durable Routing And Onboarding

Base installation provides reusable native-current-tool defaults through managed guidance. It
does not create a Claude account, install an optional CLI, change permissions or create a relay.
The coordinator consults explicit user choice and accepted local preferences before defaulting
to native operation. Offer cross-tool setup only when requested or materially relevant, not at
every plan; a declined offer does not block native work.

Requested setup is a hybrid workflow: explain service/account/usage implications, route missing
CLI through tooling readiness, guide the user through native sign-in, verify readiness and
record the agreed non-secret preference at its intended scope. User credentials stay in the
tool's credential storage. Shared preference is not proof of authentication on another machine.
Reuse sufficient scoped setup authority; honor explicit chat-creation/message requirements.

When needed, create docs/repo/reference/agent-coordination.md from the installed starter only
if absent, preserving existing consumer content on repeat setup and Kit upgrades. No empty
record is needed for the native default. A shared home reference is reused by member routing;
personal preferences use the existing ignored user layer. Missing selected Claude access prompts
a clear choice to configure it, explicitly use the native tool for that assignment or pause.

Consumer operational reference: current relay locator/host, maintenance owner, model preferences,
supported local host, shared local evidence-root convention and links to the managed procedure.
Private machine paths stay in ignored local setup; examples in the Kit use public-safe placeholders.
Local reference data is discoverable configuration, not authority to message or create chats.

An ordinary member's existing root/repository reference links to its coordination home's reference,
using a repository URL plus relative path or an established local link. Verify home identity against
the existing portfolio declaration when that membership is used. Do not infer a path from the home
ID, scan every chat by title, use an old Program as the routing database, or rewrite portfolio
identity silently. A repo without portfolio membership can use its own local operational reference.

Initial real adoption covers the named coordination home. The member-navigation pattern is
fixture-proved; each actual member needs its owner's link adoption before it gains this route.
A missing or stale relay produces a bounded owner decision/replacement procedure. Preserve advisor
sessions and unfinished assignment references when replacing a relay. No automatic chat creation,
standing role registry or program-lifecycle dependency is introduced.

## Implementation Capability Scope - Commissioning And Transport

Capability:
An accountable coordinator commissions substantive work, remains available and receives original
questions/results through an existing independent relay.

Primary workflow:
Read the installed route and consumer reference; prepare role-matched brief and authority; relay
dispatches; substantive agent works; relay returns; coordinator reads and accepts or responds.

Must cover:
- Advisors, implementers and independent reviewers; durable substantive sessions and role-specific context.
- Strict relay-only mechanics, original responses, exact return destinations and one driver per session.
- Question, denial, execution failure and failed-notification reporting; no automatic authority escalation.
- Whole-plan assignments with delegated checkpoints and clear Git/release/adoption finish lines.
- Native operation without optional tooling/account and guided user-selected cross-tool setup;
  no provider questionnaire for normal execution, no forced fallback or blanket-access default.

Explicitly out of scope:
- Schedulers, durable queues, automatic retries, app-owned CLI resumes and internal app API injection.

Deferred or blocked:
- Reverse Claude-led executable recipe: deferred pending separate evidence; not a first-release blocker.

Preserve decisions:
- D-1 through D-5, D-7, D-10 and D-11.

Planner must not reinvent:
- Relay judgment boundary, human-authorization attribution, coordinator availability or tested topology.

Feature-level success evidence:
- Installed procedure lets a fresh ordinary director receive a question and final result after yielding;
  it assesses the original response and the relay performs no substantive research.

## Implementation Capability Scope - Invocation Helper

Capability:
One maintained, installed command replaces ad hoc invocation/result scripts.

Primary workflow:
A relay worker passes the approved structured request, awaits the process result and delivers the
unchanged notice through the native tool.

Must cover:
- Explicit request, non-shell launch, retained session identity, one participating writer and isolated attempts.
- Original response capture, exact message arguments, honest process/denial/error and usage metadata.
- Safe interruption/recovery evidence, private output handling, current CLI/tooling preflight.

Explicitly out of scope:
- Tool installation, permission changes, worktree creation, background daemons, direct app messaging,
  business-task interpretation, automatic retry or credential management.

Deferred or blocked:
- No blocker to implementing the selected path; unsupported CLI versions fail with a useful blocker.

Preserve decisions:
- D-2, D-5, D-6 and D-10.

Planner must not reinvent:
- Existing Go CLI delivery, no Python prerequisite, optional assignment-owned time limit and no bypass default.

Feature-level success evidence:
- Fake-process tests prove identity/failure/concurrency mechanics; installed command passes the bounded
  real onboarding flow without hand-written launch scripts.

## Implementation Capability Scope - Discoverability And Adoption

Capability:
A fresh director finds the operational relay without the coordinator explaining the setup manually.

Primary workflow:
Member/root routing leads to the consumer operational reference, then managed procedure and helper;
the named consumer adopts the qualified release and records the observed result.

Must cover:
- Stable local reference independent of Program lifecycle, owning Program appointment links and private evidence.
- Existing portfolio identity, explicit navigation link, missing/stale/mismatched route handling.
- Coherent release and named consumer default-branch adoption; expose working-copy or host limitations.
- Managed defaults plus opt-in consumer reference creation, update/upgrade preservation, personal
  versus shared preference scope and selected-but-unconfigured service handling.
- Bounded alignment of older first-run model/permission prescriptions across runbook, context
  contract, bootstrap, compiled onboard output and existing oracle/tests. Broader Claude-host
  base-onboarding UI wording remains deferred.

Explicitly out of scope:
- New portfolio schema/registry, Organization Home record families, bulk unrelated adoption or identity migration.

Deferred or blocked:
- Any local membership conflict is resolved by its owner before using that route; it does not block
  building the Kit procedure/helper or a correctly configured consumer's adoption.

Preserve decisions:
- D-7 through D-11.

Planner must not reinvent:
- Local reference placement or the distinction between home ID, chat identity and machine-local state.

Feature-level success evidence:
- One low-context director follows installed routes in a configured home/member fixture and completes
  the real response/clarification flow; no guessed identifiers, new policy registry or manual briefing.
- Small fixtures prove absent-Claude native operation, guided opt-in with missing login/tooling,
  no repeated declined offers and preserved existing references; no real account churn for tests.

## Remaining Questions And Risks

- OQ-1 — BLOCKER: no. Exact release version is chosen from current tags at the candidate boundary.
- OQ-2 — BLOCKER: no for planning. Action-time tool/version/host availability and delegated authority
  are preflight checks; a missing capability blocks the affected live step, not unrelated preparation.
- OQ-3 — BLOCKER: no. Fresh-director onboarding is delivery acceptance, not missing architectural research.
- OQ-4 — BLOCKER: no. Reverse direction and cross-host operation remain explicitly deferred.

Residual risks: host messaging behavior can change; ordinary-merge classifier behavior is not
guaranteed; local locks only coordinate participating launchers; tool outputs can contain private
material; uncertain sends may require manual reconciliation. Preserve evidence, name the actual
failure and return to the coordinator instead of adding unattended retry infrastructure.

## Handoff State

The user delegated consolidation and bounded helper design for planning on 2026-10-09. The capability
blocks freeze the chosen first delivery while leaving safe command/file details to implementation.
The [implementation plan](cross-tool-agent-coordination_implementation_doc.md) is draft and inactive.
Next review is of that coherent plan. No live relay appointment, policy installation, implementation,
release or consumer upgrade is authorized merely by this discovery checkpoint.

## Revision Notes

- 2026-10-09: user approved native-current-tool defaults and guided optional Claude first use;
  added setup/persistence boundaries and acceptance cases without creating a mandatory provider
  questionnaire, new config schema or consumer record on every install.

- 2026-10-09: included the user's human-transparency and minimal-record requirements: material
  findings plus coordinator disposition, optional response view, ignored recovery evidence and
  ordinary cleanup instead of a permanent report per message.

- 2026-10-09: clarified mechanics after independent plan review without reducing whole-plan scope:
  coordinator owns complete requests, child output survives helper interruption, known attempt
  paths support on-demand recovery, and real-member adoption is distinct from a routing fixture.

- 2026-10-09: consolidated the accepted merge/notification/overlap pilots with explicit limits;
  selected a bounded installed CLI helper under delegated planning scope; fixed durable routing,
  relay-only responsibility and fresh-director acceptance; completed discovery handoff without
  activating implementation.

- 2026-10-05: recorded the authorized, bounded documentation delivery pilot and its review,
  permission, headless-blocker and evidence boundaries; results remain pending execution.
- 2026-10-05: incorporated the approved advisory review: relayed-input authority, headless blockers,
  permission-profile ownership, session-locator recovery, reviewer context, compaction resilience
  and review-to-merge evidence. Recorded the bounded lean-relay observation without closing the
  remaining permission, failure/attention, onboarding or reverse-direction evidence questions.
- 2026-10-05: consolidated accepted cross-tool commissioning, role context, lean relay, asynchronous
  interaction, durability and proportionate-evidence direction; recorded permission consultation
  findings and the remaining implementation-shaping questions. Private model choices and pilot
  evidence remain with the commissioning owner.
