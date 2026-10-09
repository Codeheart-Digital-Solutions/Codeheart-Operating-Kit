Last updated: 2026-10-09T15:22:56Z (UTC)
Created: 2026-10-09
Status: draft

# Cross-Tool Agent Coordination Implementation Plan

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.cross-tool-agent-coordination
  kind: implementation
  purpose: Deliver reusable Codex-led Claude CLI commissioning, lean result transport, a bounded invocation helper and discoverable consumer adoption.
  first_cataloged: 2026-10-09T14:54:05Z
  catalog_metadata_updated: 2026-10-09T14:54:05Z
  relations:
    - kind: related
      target: codeheart-operating-kit.discovery.cross-tool-agent-coordination
    - kind: related
      target: codeheart-operating-kit.implementation.claude-code-operating-guide
```
<!-- END CODEHEART PLAN METADATA -->

The discovery has enough evidence for this bounded delivery. This draft turns that direction into
three ordered epics. It does not activate implementation, install permissions, appoint a relay,
create chats, publish a release or upgrade consumers. Planning scope was delegated on 2026-10-09;
implementation and its delivery grant remain for review.

## Essential Context

| File | Why it matters |
| --- | --- |
| [Discovery](cross-tool-agent-coordination_discovery_doc.md) | Accepted requirements, pilot limits, decisions and three capability-scope groups. |
| [Agent task coordination](../../../../components/agent-interface/managed/reference/agent-task-coordination.md) | Whole-plan authority, director review, report-back and lifecycle. |
| [Codex operations](../../../../components/agent-interface/managed/reference/codex-task-operations.md) and [Claude operations](../../../../components/agent-interface/managed/reference/claude-code-task-operations.md) | Current host contracts, commissioning and tool-specific constraints. |
| [Operation routing](../../../../components/agent-interface/managed/reference/operation-routing-and-dispatch.md) | Route before execution surface; owner, target and local-state discovery. |
| [Recipe maturity](../../../../components/agent-interface/managed/reference/operational-recipe-maturity.md) and [script promotion](../../../../components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md) | Keep judgment in a runbook and deterministic mechanics in a justified executable. |
| [Runbook authoring](../../../../components/agent-interface/managed/reference/runbook-authoring-standard.md) | Audience, intent, execution, evidence and stop boundaries. |
| [Placement contract](../../reference/placement-contract.md) and [consumer impact](../../reference/consumer-impact-classification.md) | Producer/consumer ownership and safety-policy review. |
| [Change runbook](../../runbooks/change-operating-kit.md) and [release runbook](../../runbooks/release-operating-kit.md) | Affected checks, broad candidate acceptance, release and native packaging. |
| [Portfolio format](../../../../components/planning-workflows/managed/reference/portfolio-coordination-format.md) | Existing home identity is authoritative; no invented locator fields. |
| internal/cli/cli.go, internal/commands/, components/agent-interface/component.yaml | Existing compiled command and managed-content delivery paths. |

Authoring checkpoint: source placement, relative links, public-core hygiene, Markdown timestamps
and canonical plan catalog were checked. The fresh independent reviewer judged corrected
candidate 7edfde4 Ready (plan only); see [the review record](attachments/independent-plan-review.md).
The final low-severity lock-recovery clarification and the user's later transparency/retention
requirements are recorded below and were not part of that reviewed candidate. Execution approval
remains outstanding. Runtime behavior is planned, not proven by document checks.

## Contents

- Section 1 — Foundation
- Section 2 — Strategy
- Section 3 — Execution Plan
- Section 4 — Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

A newly onboarded ordinary Codex director in the adopted coordination home can follow installed
instructions and consumer routing,
commission a Claude CLI advisor, implementer or reviewer through an independent relay, remain
available for conversation, and receive an original question or response after yielding.
The same Claude session can continue after clarification. Relay tasks do not research, interpret,
approve or extend the assignment.

Completion requires a qualified Operating Kit release, adoption on the named consumer's default
branch, an operational local reference, original evidence from one fresh-director flow and truthful
delivery/host limits. Member-to-home navigation is demonstrated in a configured fixture; actual
member repositories need their owners to adopt a link before their directors gain that route.
A merged source PR or a working-copy pilot alone does not finish the plan.

## 1.2 Project And Problem Context

Ad hoc scripts and remembered chat IDs made working delegation fragile. Native child completion
did not wake an idle parent; independent ordinary chat routing did. A two-director overlap pilot
preserved both assignments, but proves neither a durable queue nor universal exactly-once delivery.
The initial recipe supports Codex desktop coordinating CLI-owned Claude sessions on the tested
host. Generic responsibilities remain direction-neutral; reverse execution is deferred.

Capability coverage: commissioning/transport -> Epic 1 and the real proof in Epic 3; helper ->
Epic 2 and installed proof in Epic 3; discoverability/adoption -> Epic 1 templates and Epic 3 local
application. The plan preserves substantive roles and whole-plan execution, not just notification.

## 1.3 Current State Analysis

Existing: generic coordination and tool references, Go CLI, managed component manifests, release
validation, portfolio identity/configuration and accepted private pilot evidence. The discovery
documents observed merge permission behavior and independent relay topology.

Missing: a reusable installed recipe, compact assignment/reference starters, maintained invocation
and output capture, durable consumer relay route, and fresh-director adoption evidence.

Ownership: Operating Kit owns procedure, templates, command and compatibility. Consumers own
operational relay references, local preferences, appointments, authority and private evidence.
Program records own improvement history and program-specific appointments. A relay locator is
neither a role identity nor authorization. No Organization Home schema change is needed.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

One agent-facing runbook composes existing responsibilities with the tested topology. It calls one
new command in the installed Go CLI; no Python dependency, copied executable, service or separate
CLI package. This narrow L3 command is justified by repeated invocations and installed distribution.
It launches one CLI turn and returns evidence; the runbook and native tools retain communication,
authority, target selection and judgment.

Expected source paths (ordinary internal factoring may change without changing the contract):

    components/agent-interface/
      component.yaml                                      # modify managed file declarations
      managed/
        README.md, kit-readme.md                           # modify discoverability
        reference/
          agent-task-coordination.md                      # modify generic boundary
          codex-task-operations.md                         # modify tested host recipe routing
          claude-code-task-operations.md                   # modify CLI-owned session/permission route
          cross-tool-coordination-contract.md             # create request/result and brief examples
        runbooks/coordinate-cross-tool-task.md             # create agent-facing recipe
        templates/agent-coordination-reference.md          # create opt-in consumer reference starter
    internal/
      cli/cli.go, cli/*tests*                              # modify command dispatch/help
      commands/coordination.go                            # create narrow command entry
      coordination/invoke.go, invoke_test.go               # create mechanics and fake-process proof
    tests/test_routing.py, tests/test_packaging_resources.py # modify where affected
    docs/repo/plans/cross-tool-agent-coordination/
      cross-tool-agent-coordination_execution_log.md       # create at activation
    manifest.yaml, release-notes.md, version/resource mirrors # modify at candidate release boundary

Installed guidance remains under .codeheart/kit/docs/agent-interface/. Template content is a
managed starter, copied into a consumer reference only during explicitly scoped adoption; no new
automatic absent-file scaffold or portfolio schema field.

Named consumer adoption creates docs/repo/reference/agent-coordination.md and updates the repository
AGENTS.md route, docs/repo/README.md, portfolio README link where useful, and the current Program's
working-model cross-reference. Do not copy private operational values into the public Kit plan.
Exact private targets and IDs belong in the companion consumer handover/activation assignment.
The existing managed AGENTS "Whole-plan assignment and agent coordination" link leads to
agent-task-coordination.md, which will link to the new runbook. That managed root template is
unchanged. The consumer-owned root route separately locates local operational values.

## 2.2 Open Questions And Assumptions Requiring Clarification

- OQ-1 — BLOCKER: no; Affects: Epic 3. Release version: choose the next unused patch at the coherent
  candidate boundary, not now. No obsolete preselected release number.
- OQ-2 — BLOCKER: no for planning; Affects: Epics 2–3. Confirm installed tools/authentication,
  exact ordinary recipient and ability to pass a loaded JSON argument without ID transcription.
  Distinguish a tool's output-yield interval from any actual process/worker lifetime limit. Record
  applicable host limits and interruption behavior before qualifying implementer use; unknown
  limits remain disclosed. Test across normal tool yields and a controlled interruption, not an
  arbitrary hours-long idle run. A hard limit incompatible with the assignment blocks that live
  use and returns to the owner; do not silently reduce delivery to consultations or add a daemon.
- OQ-3 — BLOCKER: no; Affects: Epic 3. A known producer portfolio-home mismatch is outside this
  delivery's automatic authority. Use a verified configured home/member fixture; if a selected
  adoption route depends on the mismatch, its owner must reconcile it first through the existing
  configuration route. Do not silently change membership or broaden into portfolio cleanup.
- OQ-4 — BLOCKER: no; Affects: all. This draft needs execution approval. The commissioning section
  below supplies the proposed complete grant so routine steps need not be reapproved individually.

Assumptions: the supported release platforms remain macOS universal and Windows x64; the initial
real app notification recipe is qualified only on the tested macOS host. Native command mechanics
must work on supported release platforms. For this new command Windows coverage is build,
Go unit/fake-process tests and installed command availability within existing native release gates;
live Claude/desktop wakeup on Windows is unqualified. Preserve all existing native Kit release
checks. Consumers do not need Claude installed to use unrelated Kit capabilities.

## 2.3 Architectural Decisions With Reasoning

### Responsibilities and notification

Coordinator writes the complete request.json and brief, resolving executable, model, permission
and tool settings, authority, state root and exact return destination from approved local defaults.
It assigns the attempt ID and knows its deterministic evidence directory before dispatch.
Relay receives only the prepared file path and transport instructions, starts/resumes the specified
session and returns the original response. Missing/invalid fields return a blocker; relay and
worker never choose defaults, permissions, models or substantive context on the coordinator's
behalf. Only the coordinator assesses results or answers substantive questions.
Independent review remains independent of authorship. Session reuse is distinct from role identity.

Reuse a configured independent ordinary relay. For multiple active assignments, it dispatches
temporary transport workers and yields; workers send directly to the requesting ordinary chat
outside their native ancestor hierarchy. Worker/subagent recipients and native-parent app sends
are not valid routes in the observed host. Do not invent native child wakeup guarantees.

Return questions, permission denial and failures as faithfully as completion. Notification failure
is recorded separately from successful CLI execution in the known attempt directory: pending,
sent (tool accepted, not read/accepted by the director), rejected or uncertain, with the original
receipt/error. Persist pending before attempting the send. If a worker disappears before updating
it, pending remains unresolved. Coordinator checks that known directory once on demand when the
next decision depends on the result or the user asks. This recovers discoverability without a
watcher or notification guarantee; it does not automatically resend or relaunch anything.
The coordinator may discuss unrelated matters while work proceeds; it does not busy-wait, continually
inspect transcripts, or duplicate the commissioned research. Active waiting is appropriate when
the next decision genuinely depends on the result.

### Human-visible findings and proportionate records

For consequential reviews, the coordinator presents the material findings in plain language:
what was found, why it matters, the proposed correction and whether it was accepted, challenged,
deferred or remains unresolved. State the review verdict and residual uncertainty. Link the
original response when useful; a terse "review passed" or "findings addressed" is insufficient.
This substantive reporting belongs to the coordinator, never the relay. Ordinary trivial replies
do not need a formal review table or a new report.

The agent replies normally in its CLI conversation; it is not asked to author an extra response
document. CLI JSON contains the original final reply. A generated response.md may expose that
same text conveniently; it is optional, not an additional required report or competing source of
truth. Reading the original conversation is also valid. For deterministic result capture, prefer
the direct CLI result with exact session/attempt identity over scraping a large saved transcript.

Keep one original invocation result and the minimal request/session/process/delivery state needed
to resume or reconcile it. A file plus JSON-field locator can reference the final text without
another Markdown copy. Direct stdout/stderr files remain necessary for interruption evidence;
avoid duplicate summaries, manifests and receipts where one compact attempt record suffices.
The current pilot's collection of scaffolding files is not a production file-count requirement.

Generated evidence remains ignored machine-local data, not committed repository content. Record
meaningful accepted findings/decisions and the implementer/reviewer locator in the existing owning
plan/log; do not create a permanent document for every turn. Once accepted and no longer needed
for active work, correction review or recovery, generated duplicates/logs are eligible for normal
authorized cleanup. Preserve required original evidence and unfinished session/recovery state;
no automatic sweeping, scheduled retention service or silent deletion is introduced.

### Bounded invocation contract

Proposed command: codeheart-operating-kit coordination invoke-claude --request <request.json>.

Required request fields: schema version; assignment/attempt identifiers; exact return thread/host;
working directory; brief file; executable; chosen model; new/resume mode and retained session ID
for resume; explicit permission/tool settings; authority/source references; shared host-local
state root. A new invocation allocates and records its session before execution. Reject unknown
fields or unsupported mode/flag combinations rather than silently falling back.

Settings are concrete paths/values, not shell fragments. Use the approved noninteractive CLI flags,
JSON output and stdin for the brief. No arbitrary extra command string, bypass/full-access default,
implicit installation, model substitution or permission escalation. Preserve default and explicit
deny/ask behavior; unsupported interactive prompts must return a blocker instead of waiting forever.
Use one documented fixed flag set and record the CLI version as metadata. Do not scrape --help
on every run or maintain a version/flag compatibility matrix. Return a recognized unknown-option
or usage error with original stderr as an unsupported-invocation blocker; other nonzero failures
remain execution errors. No fallback flags or replay. Consult official help/docs when maintaining
the recipe, not as a new runtime compatibility engine. Settings validation cannot prove that a
classifier will allow an action.

Use one ignored shared state root for participating launches through a relay on one host.
An exclusive-create session lock holding the attempt ID and unique attempt directories prevent
accidental second launch, overwrite and replay. Retain session/attempt information before spawning;
record launcher and child process IDs with start identities as soon as known. No automatic staleness
heuristics or lock stealing. Recovery documents a deliberate release of the named attempt's lock,
only after verifying the same lock owner and that both recorded processes have ended. Refuse release
while either is alive or process identity/ownership is uncertain, including a launch interrupted
before identity capture. Preserve original outputs and session state when releasing the lock. A participating-process
lock cannot prove there is no independent CLI/app/other-host writer. Fresh CLI ownership and the
runbook remain necessary.

Open attempt stdout/stderr files before spawning and connect the child directly to those files,
not to a helper-owned pipe or an in-memory buffer. Helper/worker death must not erase bytes already
written; it may still terminate the child under host rules. Parse the completed JSON afterwards,
preserving incomplete output and absence of a final record as uncertainty. Keep a compact
result record: requested/actual session, process state/exit, timestamps, model/usage if available,
original response locator (file and optional JSON field), permission-denial metadata and error/status.
Use atomic final record
writes; normal stdout is only this compact metadata. Emit a native-tool message argument file
with the exact return destination and original short reply or response reference. No semantic
digest, fabricated completion state or task acceptance. The relay loads the file directly into
the exposed tool; the command never connects to private app APIs. If the send tool reports a
recipient, compare it with the prepared recipient and retain a mismatch as failed/uncertain
delivery. Never infer recipient confirmation from an omitted field. Cap inline original replies at
2,000 characters; longer responses use an exact file reference without a relay-written summary.

The helper launches once and waits on its child process; it is not a scheduler. It adds no
whole-plan timeout. Normal host output-yield returns are resumed on the same process handle;
they are not timeouts or reasons to restart Claude. An actual limit/interruption records uncertain
effects and observed process state when the writer survives; if it does not, the retained attempt
and raw files remain inspectable. On demand, the coordinator checks writer/child state, original
session and any work effects before resuming. Never infer clean completion from a missing process.
The initial release covers resumable whole-plan work with honest host limits; it does not promise
uninterrupted hours-long execution or delivery while the app/host is shut down. No automatic
child relaunch, retry, lock deletion, cleanup of task work or termination of unrelated processes.
Private artifacts are not blindly printed, committed or sent to another recipient. Protect local
files appropriately to the host; record portable limitations rather than claiming encryption.

### Local routing and lifecycle

Use the consumer operational reference for current relay/host/owner/model preferences and links
to procedure and program appointments. Root instructions point there. Existing portfolio config
owns home identity; an explicit repository URL/path link supplies navigation. Missing, stale or
mismatched routes produce a clear owner question, never a guessed UUID/path or silent enrollment.
Standalone consumers can keep their own reference without configuring a portfolio.

The reference is maintained independently of Program completion. The owning plan/log records who
actually implemented or reviewed work, with session and result/acceptance references. Ignored
runtime files carry raw evidence and machine paths; durable private records retain conclusions
and necessary locators. Replacing a relay must preserve active assignment and substantive-session
references. No global role directory is introduced.

### Permissions, release and evidence

Document the tested invocation-scoped ordinary-merge pattern: tool permission may allow ordinary
task-repository merges while workflow checks delegated authority, review, CI and candidate identity.
Preserve defaults, hard restrictions and explicit ask/deny rules. A real denial returns unchanged;
do not have another agent perform the same action merely to evade the refusal. Existing independent
authority can be considered separately by the accountable owner.

A template/example is not an installed permission change. Live settings require authority covering
their effects; an approved whole-plan grant can cover them once. No global policy installation,
force/admin merge, direct default-branch push or release/deployment permission is implied.

Impact: managed instruction additions plus executable command behavior and security/safety-policy
guidance; no consumer schema migration or automatically owned scaffold. Use broad candidate
acceptance, explicit review of permission wording, native packaging/install proof and release notes.
Run cheap affected checks during edits and one coherent broad candidate gate. Reuse existing live
pilot evidence where unchanged; do not repeat the merge experiment merely for ceremony.

### Proposed whole-plan commissioning and Git boundary

Upon explicit whole-plan execution approval, one implementer can execute all three ordered epics,
with the commissioning director owning acceptance and material decisions. Communicate the plan,
epic and real Git state. Commit coherent work, normally push and open/update one delivery PR;
combine Epics 1 and 2 for independent source review, then correct in the same review cycle.
No user approval per checkbox, routine commit, covered push or epic review.

The proposed finish line includes normal reviewed merge, one patch release under the existing
authorized release audience/signing boundary, and the named consumer default-branch adoption.
The director may delegate integration and renew routine execution within that grant; changed
audience, broader permissions, spending commitments, failed checks and scope expansion return to
the human owner. Exact repositories/branches, release audience and private targets must be in the
activation assignment before effects. Drafting now authorizes none of these external effects.

Include one temporary ordinary onboarding director chat, use of the existing approved independent
relay, its temporary transport workers, the scoped Claude test session and report-back messages
in the activation grant. Archive that temporary director only after its result is accepted;
retain the operational relay. If the host requires explicit human approval for chat creation or
messaging, verify that the actual commissioning grant covers it; agent-authored text cannot
supply human approval. This is a planned single approval boundary, not hidden automatic authority.

Record adoption on the default branch and reconcile the explicitly assigned active worktree
without overwriting unrelated work. A dirty checkout or independent in-flight task is preserved;
use a clean adoption branch and record any remaining reconciliation owner. No repo cleanup,
bulk upgrades, history rewriting, AWS effects or unrelated product work. /goal activation occurs
only if explicitly requested and verified.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| Epic 1 — Reusable coordination contract | Installed routes and compact templates describe faithful delegation, authority, relay-only transport and durable local routing. | S | Execution approval |
| Epic 2 — Maintained invocation helper | One installed command launches/captures assigned Claude CLI work with isolated identities and honest failure evidence. | M | Epic 1 contract |
| Epic 3 — Release, adoption and onboarding | Qualified release is adopted and a fresh director uses it without manual setup coaching. | M | Reviewed Epics 1–2 |

## Epic 1 — Reusable Coordination Contract

**A) Epic ID, Title, And Outcome:** Epic 1. A new director can identify the correct roles, prepare a
complete brief, find its local relay and execute the documented transport workflow.

**B) Scope:** One agent-facing L1 runbook with a compact intent block; update existing generic and
tool references instead of duplicating them. Provide role brief/request examples, relay assignment
and response templates, permission example and consumer reference starter.

**C) Files Touched:** Agent-interface runbook/reference/template paths and component manifest from
Section 2.1; README routers and affected packaging resources/tests. Consumer files wait for Epic 3.

**D) Acceptance Criteria And Size:** S. Procedure states preflight, authority, phases, exact return
route, original evidence, blockers and recovery. Relay cannot research or accept work. The coordinator
yields. Generic versus tested-host claims are explicit. No new portfolio fields or forced scaffold.

**E) Dependencies And Critical-Path Notes:** The discovery is the contract. Required CLI mechanics
are specified for Epic 2, not presented as already available.

**F) Tasks Checklist:**
- [ ] Add runbook and compact role/relay/return templates, including one coherent whole-plan brief,
  session reuse, explicit report-back authority and a question/continuation example.
- [ ] Encode exact topology, strict relay limits, process-versus-task completion, failed-send
  evidence, missing/stale relay handling and coordinator availability in the existing routes.
- [ ] Add the optional local reference starter and explicit member-to-home navigation guidance.
  Preserve Program appointments and existing config authority.
- [ ] Document the scoped ordinary-merge example and denial behavior without changing live settings.
- [ ] Include human-visible review findings and coordinator disposition, optional generated
  response views, minimal ignored runtime evidence and ordinary authorized cleanup guidance.
  Keep durable conclusions in existing plan/log records rather than a new report per turn.
- [ ] Update manifest/resource mirrors and nearest routers; run affected routing/resource and
  public-core/Markdown checks. Review together with Epic 2.

**G) Implementation Notes:** Missing binary/runtime -> existing tooling-readiness route; missing
login, mode capability or app message tools -> service/host preflight, not package installation.
No new runtime is required. The runbook separates technical execution from user-facing decisions,
with plain plan/epic names and no invented per-step approval. Producer owner maintains the route.
Use the final low-context Epic 3 exercise as the fresh routing probe; avoid a duplicate warm test.

**H) Open Questions:** OQ-2 applies at live execution; no unresolved design blocker.

## Epic 2 — Maintained Invocation Helper

**A) Epic ID, Title, And Outcome:** Epic 2. Relay workers can use the installed command instead of
constructing one-off scripts and manually copying session or return identifiers.

**B) Scope:** Narrow Go command and standard-library mechanics from Section 2.3. No daemon,
substantive agent logic, credential manager, direct app messaging or general adapter framework.

**C) Files Touched:** internal/cli/cli.go, internal/commands/coordination.go,
internal/coordination/ implementation/tests; command help/contracts and affected Go test fixtures.
Do not implement the new feature in the legacy Python behavior oracle merely for symmetry.

**D) Acceptance Criteria And Size:** M. One invocation uses the expected brief, worktree, session,
model and settings; originals survive; session mismatch/denial/failure is truthful; notification
payload preserves the exact recipient. Same-session overlapping attempts cannot both launch.

**E) Dependencies And Critical-Path Notes:** Epic 1 fixes the contract. Existing release CLI behavior
must remain unchanged; new command parity is specified by its Go tests, not a Python copy.

**F) Tasks Checklist:**
- [ ] Implement strict request/preflight, new/resume identity persistence, explicit settings and
  structured argument/stdin launch; reject unsupported flags/modes without fallback.
- [ ] Implement shared-root session exclusivity, isolated attempts, retained original output,
  atomic result metadata and exact native-message argument output.
- [ ] Implement failure/denial/interruption reporting and coordinator-directed recovery. Do not
  auto-retry or infer semantic completion from exit zero.
- [ ] Add focused fake-CLI tests for quoting/spaces and stdin fidelity; wrong/empty session and
  malformed/error/denial outputs; nonzero exit; two distinct concurrent assignments; same-session
  collision/replayed attempt; helper termination while a fake child is writing; retained partial
  files and unresolved pending state; interrupted/stale attempt; deliberate lock release refused
  while the recorded process is alive or identity is unknown, permitted only after verified exit
  and matching lock ownership; no raw normal output or implicit escalation. Use sanitized real CLI response/denial and message-rejection shapes to anchor fake
  fixtures. Exercise actual process behavior, not just mocked return values.
- [ ] Verify on the supported live host that a controlled harmless process continues through
  normal tool output-yield/resume cycles and retains evidence on interruption. Record actual
  process-limit knowledge separately from observed duration; do not burn hours testing idleness.
- [ ] Run affected Go tests on supported platforms and one independent coherent source review
  across Epics 1–2, emphasizing authority, result fidelity and process/concurrency failure behavior.
  Address findings without restarting unchanged reviews.

**G) Implementation Notes:** Target maturity L3 narrow command, called by the L1 runbook. Request
and result versioning belong to this command, not the global Kit config schema. Use per-host
filesystem/process primitives appropriately; record platform limits. Check current CLI help/docs
when implementing, and test unsupported versions as structured blockers. Do not consume raw
reasoning traces or copy private pilot transcripts into fixtures.

**H) Open Questions:** OQ-2 is action-time compatibility. Coordinator decides any material CLI/API
scope change; ordinary internal factoring remains implementer-owned.

## Epic 3 — Release, Adoption And Fresh-Director Onboarding

**A) Epic ID, Title, And Outcome:** Epic 3. A qualified installed release and consumer-owned routing
let a fresh ordinary director use the arrangement without conversation-history coaching.

**B) Scope:** One patch release, one named coordination-home consumer default-branch adoption,
reconciliation of the assigned active checkout and one combined real onboarding/routing proof.
Use a temporary member fixture for member-to-home navigation rather than changing another live
program repository merely for a test.

**C) Files Touched:** Release identity/notes and affected mirrors under the release runbook;
execution log; consumer kit config/lock/managed content only through upgrade, local operational
reference and its existing root/index/Program routes. Private identities/evidence stay consumer-owned.

**D) Acceptance Criteria And Size:** M. Broad release checks pass; assets match accepted source;
consumer default branch adopts that release; original question and final response reach the fresh
director after yielding; it reads originals and assesses them. Stale/missing routes fail clearly.
Record host support and any outstanding working-copy reconciliation separately from adoption.

**E) Dependencies And Critical-Path Notes:** Source review and authority precede release. Release
precedes live installed-use acceptance. Use candidate fixtures for cheap routing checks before
publication, then the same scenario once on installed release. Do not repeatedly live-test merges.

**F) Tasks Checklist:**
- [ ] Confirm exact execution grant, named consumer/default branch, host/relay owner and existing
  membership. Preserve conflicts/unrelated work; use a configured fixture for routing proof.
- [ ] Select release version, record impact/safety review, update identity and notes, run the broad
  candidate/native release gates and integrate the reviewed candidate normally.
- [ ] Publish and verify the permitted patch release; upgrade the named consumer through the
  managed lifecycle route. Commit/push/merge the adoption checkpoint under the whole-plan grant.
- [ ] Create its local operational reference and links; move ongoing relay routing authority out
  of the improvement Program while retaining its history and program-specific appointments.
  Configure exact model/host/session preferences privately and reconcile the assigned checkout.
- [ ] Give a fresh ordinary director a small user-style request. Require it to discover the
  installed route/reference from a configured member fixture/home, launch one harmless consultation
  through the maintained helper/relay, yield, receive a question, resume the same session with an
  authorized answer, yield again and read the original final response. Do not mention the relay
  or its path in the initial user-style request. Observe relay behavior: no research/digest/
  acceptance, exact destination and session, coordinator remains available.
- [ ] Within that exercise, use one harmless deliberately denied tool request in a scoped
  read-only test profile; verify the installed helper preserves the real denial and returns it.
  Anchor the failure fixtures to this actual CLI shape. Exercise rejected-notification persistence
  using an injected send rejection with the already observed native rejection shape; verify the
  coordinator can find the retained result by its known attempt path on demand. No arbitrary
  real recipient or repeat live negative messaging is required unless actual host behavior changed.
- [ ] Cover missing/stale route handling in the same onboarding exercise or non-live fixture.
  Retain existing accepted overlap evidence; repeat only if the new mechanics/topology invalidate it.
  Record live usage/time if exposed and distinguish cached input from new output; no invented budget threshold.
- [ ] Director presents material findings, accepted/challenged corrections and residual limits
  to the human, with an original-response reference. Verify generated artifacts are ignored and
  sufficient for recovery without requiring a duplicate authored report.
- [ ] Director accepts evidence and honest support limits; archive only the temporary approved
  test chat, retain relay and original results, publish final plan/log/adoption state.

**G) Implementation Notes:** Installation tests prove packaged routing and command availability;
the real pilot proves user experience and host wakeup. Neither substitutes for the other. Raw
private evidence remains ignored, with a short durable consumer acceptance record and public-safe
summary in the producer log. A failure after publication is a disclosed adoption defect to correct,
not permission to rewrite a tag or claim completion. Preserve earlier release evidence where valid.

**H) Open Questions:** OQ-1–3 are named action-time preflights. Broader organization rollout and
Windows desktop notification qualification are outside this finish line.

# Section 4 - Future Planning

Deferred deliberately: reverse Claude-led executable recipe, cross-host delivery, Windows app
notification qualification, durable queues/retry services, automated advisor selection, global
role registry and general agent orchestration. Promote only when a real workflow requires them.

A future host change can invalidate notification proof without invalidating generic responsibility
rules. Recheck the affected boundary, retain original failures and adapt one recipe. Do not turn
every transient failure into framework work. No AWS or product feature implementation is included.

# Revision Notes

- 2026-10-09: recorded independent Ready verdict on 7edfde4; added the reviewer's remaining
  low-severity deliberate lock-release clarification and the user's requirement for transparent
  findings and minimal ignored evidence. These final clarifications are author-checked; no new
  implementation or release claim is implied.

- 2026-10-09: addressed independent review with coordinator-owned requests, explicit durable
  failure lookup, direct child output files, bounded lifetime/interruption proof, narrower member
  adoption claims, a fixed CLI flag contract and a real harmless denial within onboarding.
- 2026-10-09: drafted from the completed discovery and accepted bounded pilots under delegated
  planning scope; selected narrow compiled invocation helper, three epics, strict relay boundary,
  stable consumer routing and one combined fresh-director acceptance. Execution remains inactive.
