Last updated: 2026-10-09T19:01:13Z (UTC)
Created: 2026-10-09
Status: active

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

The user commissioned this three-epic implementation on 2026-10-09 after the approved planning
amendments merged. The [execution log](cross-tool-agent-coordination_execution_log.md) records
activation, delivery boundaries and actual progress. One implementer owns the full plan;
the director accepts independent source review after Epics 1 and 2 before release and adoption.
Activation is not evidence that an epic, release, adoption or native-goal invocation has completed.

## Essential Context

| File | Why it matters |
| --- | --- |
| [Discovery](cross-tool-agent-coordination_discovery_doc.md) | Accepted requirements, pilot limits, decisions and three capability-scope groups. |
| [Plan execution](../../../../components/planning-workflows/managed/runbooks/execute-implementation-plan.md) | Owns the native-goal default during implementation commissioning. |
| [Planning review](../../../../components/planning-workflows/managed/runbooks/review-planning-document.md) | Owns the review-versus-amendment boundary; coordination guidance points here. |
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
and canonical plan catalog were checked. The independent Claude reviewer judged corrected
candidate b6132dd **Ready (plan only)** after closing all material and low-severity amendment
findings; see [the review record](attachments/independent-plan-review.md). The coordinator read
the original response and accepted that assessment. The user subsequently authorized starting
implementation; the activation assignment now covers the whole plan as recorded in the execution log.
Runtime behavior is planned, not proven by document checks or this planning verdict. The subsequent
user-approved report-before-amendment and native-goal default rules below are author-checked;
they were not part of b6132dd's
independent verdict. Recording these approved rules does not commission another review cycle.

## Contents

- Section 1 — Foundation
- Section 2 — Strategy
- Section 3 — Execution Plan
- Section 4 — Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

A fresh Kit consumer can implement with the current tool without a Claude account, a relay or
a coordination reference. When the user selects cross-tool delegation, an ordinary Codex director
in the adopted coordination home can follow guided setup and installed consumer routing,
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
and output capture, durable consumer relay route, and fresh-director adoption evidence. Current
first-run onboarding prescribes a specific Codex model and Full access in the runbook, context
contract, bootstrap guide and compiled onboard output. Align those conflicting lines and their
existing oracle/test expectations with preserving the user's tool/model/permissions.

Ownership: Operating Kit owns procedure, templates, command and compatibility. Consumers own
operational relay references, local preferences, appointments, authority and private evidence.
Program records own improvement history and program-specific appointments. A relay locator is
neither a role identity nor authorization. No Organization Home schema change is needed.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

One hybrid runbook separates a short optional setup dialogue from agent execution of the tested
topology. Existing generic guidance supplies the native-tool default. The cross-tool runbook calls one
new command in the installed Go CLI; no Python dependency, copied executable, service or separate
CLI package. This narrow L3 command is justified by repeated invocations and installed distribution.
It launches one CLI invocation and returns evidence; the runbook and native tools retain communication,
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
          onboarding-context-contract.md                  # align conflicting model/access lines only
          cross-tool-coordination-contract.md             # create request/result and brief examples
        runbooks/coordinate-cross-tool-task.md             # create hybrid opt-in/setup/execution recipe
        runbooks/conduct-first-run-onboarding.md            # modify only conflicting model/access defaults
        templates/agent-coordination-reference.md          # create opt-in consumer reference starter
    components/planning-workflows/managed/runbooks/
      review-planning-document.md                          # clarify coordinator review/amendment authority
      execute-implementation-plan.md                       # default to verified native goal commissioning
    bootstrap.md                                          # align conflicting model/access lines only
    internal/
      commands/onboard.go                                 # align existing onboard output only
      cli/cli.go, cli/*tests*                              # modify command dispatch/help
      commands/coordination.go                            # create narrow command entry
      coordination/invoke.go, invoke_test.go               # create mechanics and fake-process proof
    src/codeheart_operating_kit/commands/onboard.py         # align existing behavior oracle output
    src/codeheart_operating_kit/resources/components/agent-interface/ # sync touched managed mirrors
    src/codeheart_operating_kit/resources/components/planning-workflows/ # sync touched workflow mirrors
    tests/test_onboard.py, tests/test_install_metadata.py   # replace conflicting output assertions
    tests/test_routing.py, tests/test_packaging_resources.py # modify where affected
    docs/repo/plans/cross-tool-agent-coordination/
      cross-tool-agent-coordination_execution_log.md       # create at activation
    manifest.yaml, release-notes.md, version/resource mirrors # modify at candidate release boundary

Installed guidance remains under .codeheart/kit/docs/agent-interface/. Template content is a
managed starter, copied into a consumer reference only during explicitly scoped adoption; no new
automatic absent-file scaffold or portfolio schema field. Every installation gets the managed
native-tool default through generic guidance; no empty consumer reference is required for that.
After the user configures a custom coordination arrangement, the runbook creates the consumer-owned
reference from the managed starter only if absent, or makes the specifically authorized change to
an existing one. Kit repair/upgrade must preserve its contents.

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
- OQ-4 — BLOCKER: no; Affects: all. Implementation start is authorized. Record the exact targets,
  goal authority and delivery grant in the activation assignment before dispatch; routine covered
  steps do not need repeated approval.

Assumptions: the supported release platforms remain macOS universal and Windows x64; the initial
real app notification recipe is qualified only on the tested macOS host. Native command mechanics
must work on supported release platforms. For this new command Windows coverage is build,
Go unit/fake-process tests and installed command availability within existing native release gates;
live Claude/desktop wakeup on Windows is unqualified. Preserve all existing native Kit release
checks. Consumers do not need Claude installed for normal planning, implementation or native
agent delegation. Choosing a provider does not create an account or prove that another developer
on another machine is signed in.

## 2.3 Architectural Decisions With Reasoning

### Native default and optional first use

The generic installed rule is: honor the explicit assignment choice and applicable local
preferences; otherwise continue with the current tool. A user working in Codex can use normal
Codex implementation, and a user working in Claude can use its normal native workflow. This does
not qualify reverse cross-tool execution. Native work uses its own supported delegation/reporting
surface, not the Claude invocation helper or a relay. Usual chat-creation and execution authority
still applies; installing the Kit never creates an implementation chat by itself.

A fresh user approving implementation should not have to answer a provider questionnaire or
configure an optional service. Explain cross-tool delegation briefly when requested or materially
useful; do not advertise it at every plan. If declined, continue the native workflow and avoid
repeated offers. Remember a choice only at its agreed scope: this assignment, personal preference,
repository or coordination home. An unscoped decline applies to the current conversation or
assignment; do not persist a team preference without agreement. The general rule against repeated
offers still applies. Do not turn one developer's decision into a team-wide default.

Example opt-in explanation: "We can continue with Codex. If you want Claude to implement or review,
I can guide that setup; it requires your own working access to Claude." Ask for the choice before
starting optional installation, authentication or relay creation. Existing clear choice and
sufficient setup authority are reused; no new approval per routine setup step.

For the selected Claude route, the coordinator performs this order:
1. Resolve existing local/home preferences and the requested role/model. Check local CLI presence
   through the existing tooling-readiness route; do not install or upgrade merely because it is absent.
2. Explain the applicable account/access and potentially chargeable use requirements from current
   official service guidance. Guide an authorized installation if needed. Existing base Kit setup
   is not rerun and no new account, purchase or usage commitment is created implicitly.
3. Check authentication through the tool's supported status/preflight. If missing, guide the user
   through the tool's own sign-in flow. Passwords, tokens and MFA stay there, never in chat prompts,
   the consumer reference, invocation request or Git. Recheck success before commissioning.
4. If unavailable, explain the concrete blocker and offer native execution for this assignment,
   completing the requested setup, or pausing. Never silently replace an explicitly selected
   provider/model. Do not interpret a shared preference as proof of personal machine readiness.
5. Record the agreed non-secret arrangement at the agreed scope. Prefer the existing home
   reference when shared; use the repository reference for repo-local choices; personal preferences
   use the existing ignored user layer. Ordinary native operation needs no new reference file.
6. Reuse an authorized existing relay. If none is usable, include the exact proposed ordinary relay
   chat and report-back recipient/effects in the user's setup choice; create only with sufficient
   explicit chat/message authorization, then record its actual locator. No user needs to know or
   transcribe a chat ID. A launcher/request cannot run before this route and authority are ready.

Managed defaults and setup procedure live under .codeheart/kit/; the generated local reference is
consumer-owned even though the Kit-guided workflow creates it. Preserve edits on reconfiguration
unless the exact change is authorized; never sync a managed template over consumer choices.
Authentication stays in the tool's credential storage and is checked on the current machine.
Machine paths/session artifacts stay ignored. No new global Kit-config fields or role registry.

Align only the conflicting model/reasoning/speed/access lines across the first-run runbook,
onboarding-context-contract.md and its packaged mirror, bootstrap.md, Go onboard output and the
existing Python onboard behavior oracle. Replace the old output assertions in test_onboard.py
and test_install_metadata.py with the accepted preserve-user-choice wording; retain coverage of
the existing onboarding sequence. Use existing parity/resource checks for mirrors, not a new
Python implementation of the coordination helper. Add only a short optional-route pointer where
appropriate; no mandatory provider setup or unrelated onboarding redesign. This bounded alignment
prevents the CLI itself from retaining Full access as a hidden prerequisite.

Generic native implementation works with the current tool after Kit installation. Converting
existing Codex-specific base installation/onboarding UI wording to a complete Claude-host guide
is deferred; release notes must not imply that broader onboarding path was delivered here.

### Native goals for implementation-plan execution

Executing an approved implementation plan defaults to the implementing tool's native goal mode,
for Codex and Claude alike. This applies to complete implementation assignments, not routine edits,
discovery, consultation or review. Preserve the user's explicitly selected exception.

Commission the concrete objective, agreed finish line, constraints, review handoffs and reporting
route together with the existing execution grant. Carry actual user authority, including an
accepted standing goal preference; never treat installing the Kit as consent that overrides an
explicit-authorization requirement in the host. Where a further explicit goal choice is required,
resolve it in the same commissioning decision. Do not ask again per epic or routine continuation.

Use each tool's supported activation and status mechanisms and retain observable activation
proof. A prompt mentioning /goal is not proof. For Claude CLI, pass the supported native goal
invocation in the prepared request/brief; the existing helper transports it rather than becoming
a goal engine. Codex uses its own exposed goal mechanism. Do not add a scheduler, automatic
permission escalation, invented token budget or separate continuation framework.

Required review or genuine blockers must return control and original evidence to the coordinator,
without spinning or advancing dependent work. Verify the selected host's goal/handoff behavior
before claiming unattended continuation; state any limitation. Keep overall plan completion
separate from process exit, evaluator verdict or an intermediate checkpoint. Use supported native
continuation after delegated acceptance. If activation cannot be verified or is unsupported,
disclose that fact and continue only under an authorized alternative; do not silently claim goal
mode or change tools. The review/reporting boundary remains unchanged.

Update the existing execution runbook, generic coordination reference and both per-tool guides
together, replacing contradictory optional-only wording. Apply matching packaged mirrors and
examples. Include the change in the already-planned patch release and named consumer adoption;
there is no separate prerequisite release. The current assignment can explicitly use this
user-approved default before the reusable guidance ships.

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

For consequential reviews, the coordinator reads the original and reports its accessible link
or precise source locator, core findings and significance, an independent assessment including
any disagreement, recommended disposition, verdict and residual uncertainty. A terse "review
passed" or "findings addressed" is insufficient. The relay never supplies this judgment.

During discovery, planning and discussion, report before incorporating findings or commissioning
corrective edits. Wait for the user's decision unless explicit existing authority covers the
correction cycle. "Apply this agreed change and get a review" covers that change and review, not
new changes arising from it. A reviewer recommendation, Ready verdict or the coordinator's
agreement is not a substitute for that decision. Distinguish recommended, user-approved and
implemented changes in the report. Once the user approves specific corrections, carry them out
without another per-edit approval. A wider correction-cycle delegation is reusable only within
its stated limits.

During authorized implementation, routine review fixes remain covered by the existing execution
grant; preserve delegated director acceptance and report outcomes. Material scope, outcome or
authority changes still return to the responsible owner. This boundary adds no user checkpoint
per implementation finding or epic. Ordinary trivial replies need no formal review table.

Put this reusable boundary in review-planning-document.md, clarifying its existing no-rewrite
rule for both the reviewing agent and the commissioning coordinator. Generic task coordination
and the cross-tool recipe link to that rule and carry the assignment's actual correction authority;
do not create a second approval system or tool-specific policy copy.

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
accidental second launch, overwrite and replay. Capture launcher identity before creating the
lock and persist it with session/attempt information before spawning; record child ID and start
identity as soon as known. No automatic staleness heuristics or lock stealing. Recovery documents
a deliberate release of the named attempt's lock only after verifying the same lock owner and
that both processes have ended. Refuse automated release while either is alive or identity/ownership
is uncertain, including interruption around child identity capture. For unresolved identity, the
coordinator records the blocker and chooses either manual verification sufficient to establish
exit and ownership before release, or an explicitly handed-over new session retaining the old
uncertain attempt. Do not replay uncertain actions or permit a competing writer: resumed effects
require reconciliation of prior process/workspace/external state, or isolated non-conflicting work.
No indefinite wait, silent lock deletion or automatic new-session retry. Preserve outputs and
session state. A participating-process lock cannot prove there is no independent CLI/app/other-host
writer. Fresh CLI ownership and the runbook remain necessary.

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

Impact: managed instruction additions, the new invocation command, changed existing onboard
command output, security/safety-policy guidance and optional first-use routing; no consumer schema migration or automatically owned
scaffold. The consumer reference is created by guided opt-in, not by mandatory base installation.
Use broad candidate acceptance, explicit review of permission wording, native packaging/install
proof and release notes.
Run cheap affected checks during edits and one coherent broad candidate gate. Reuse existing live
pilot evidence where unchanged; do not repeat the merge experiment merely for ceremony.

### Approved whole-plan commissioning and Git boundary

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
activation assignment before effects. The 2026-10-09 activation records that grant; effects still
require their action-time checks and delegated review acceptance.

Include one temporary ordinary onboarding director chat, use of the existing approved independent
relay, its temporary transport workers, the scoped Claude test session and report-back messages
in the activation grant. Archive that temporary director only after its result is accepted;
retain the operational relay. If the host requires explicit human approval for chat creation or
messaging, verify that the actual commissioning grant covers it; agent-authored text cannot
supply human approval. This is a planned single approval boundary, not hidden automatic authority.

Record adoption on the default branch and reconcile the explicitly assigned active worktree
without overwriting unrelated work. A dirty checkout or independent in-flight task is preserved;
use a clean adoption branch and record any remaining reconciliation owner. No repo cleanup,
bulk upgrades, history rewriting, AWS effects or unrelated product work. Native goal execution is
the approved default; the assignment explicitly carries its authority and objective, and verifies
activation under the selected host contract before claiming that it is active.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| Epic 1 — Reusable coordination contract | Native defaults, guided optional setup and compact templates make coordination usable with or without Claude. | M | Recorded activation assignment |
| Epic 2 — Maintained invocation helper | One installed command launches/captures assigned Claude CLI work with isolated identities and honest failure evidence. | M | Epic 1 contract |
| Epic 3 — Release, adoption and onboarding | Qualified release is adopted and a fresh director uses it without manual setup coaching. | M | Reviewed Epics 1–2 |

## Epic 1 — Reusable Coordination Contract

**A) Epic ID, Title, And Outcome:** Epic 1. A new director can identify the correct roles, prepare a
complete brief and use the native default immediately, or select guided optional setup, find its
local relay and execute the documented cross-tool workflow.

**B) Scope:** One hybrid L1 runbook with a compact intent block, paced opt-in dialogue and separate
agent execution path. Update existing generic/tool references and conflicting first-run defaults.
Provide role brief/request examples, relay/response templates, permission example and opt-in
consumer reference starter.

**C) Files Touched:** Agent-interface runbook/reference/template paths and component manifest from
Section 2.1; planning-workflows review-planning-document.md, execute-implementation-plan.md and
their existing packaged mirrors;
README routers and affected packaging resources/tests. The bounded onboarding
alignment also touches onboarding-context-contract.md and its packaged mirror, bootstrap.md,
internal/commands/onboard.go, src/codeheart_operating_kit/commands/onboard.py and existing
tests/test_onboard.py and tests/test_install_metadata.py. Consumer files wait for Epic 3.

**D) Acceptance Criteria And Size:** M. Native work proceeds with no Claude install/account or
reference file. Optional setup handles missing CLI/login, user choice, permission preservation,
exact relay authority and consumer-owned reference creation/update. Procedure states preflight,
authority, phases, exact return route, original evidence, blockers and recovery. Relay cannot
research or accept work. The coordinator yields. Generic versus tested-host claims are explicit.
No new portfolio fields or forced scaffold. Review-only/discussion assignments report before
amending; explicitly delegated correction cycles and routine authorized implementation retain
their autonomy. Native implementation goals are the default in both tool guides and the generic
execution route, with explicit authority/activation and honest fallback behavior.

**E) Dependencies And Critical-Path Notes:** The discovery is the contract. Required CLI mechanics
are specified for Epic 2, not presented as already available.

**F) Tasks Checklist:**
- [ ] Add runbook and compact role/relay/return templates, including one coherent whole-plan brief,
  session reuse, explicit report-back authority and a question/continuation example.
- [ ] Encode exact topology, strict relay limits, process-versus-task completion, failed-send
  evidence, missing/stale relay handling and coordinator availability in the existing routes.
- [ ] Add native-current-tool selection defaults to generic coordination, the paced optional setup
  flow to the hybrid runbook and a bounded first-run pointer. Replace conflicting pinned-model/
  Full-access prescriptions across every listed onboarding surface and its output assertions.
  Preserve the existing sequence and check Go/oracle output and packaged-resource parity.
  No mandatory provider setup, repeated offers or unrelated onboarding rewrite.
- [ ] Add the opt-in local reference starter and explicit member-to-home navigation. Guide creation
  only after an agreed arrangement; preserve existing contents, personal/team scope, Program
  appointments and config authority. Authentication is per machine in the tool's storage.
- [ ] Document the scoped ordinary-merge example and denial behavior without changing live settings.
- [ ] Align the execution runbook, generic coordination and Codex/Claude guides on default native
  goal commissioning for approved implementation plans. Include actual user authority, activation
  proof, review/blocker handoff and explicit exceptions. Replace conflicting optional-only wording;
  native current-tool execution still needs no second tool/account/relay.
- [ ] Clarify the existing planning-review runbook's coordinator boundary: original response link,
  findings, independent assessment and recommended disposition before discussion-stage amendments;
  honor explicit correction delegation and existing implementation authority. Link from generic
  coordination and cross-tool guidance. Include the agreed-edit-plus-review example.
- [ ] Include human-visible review findings and coordinator disposition, optional generated
  response views, minimal ignored runtime evidence and ordinary authorized cleanup guidance.
  Keep durable conclusions in existing plan/log records rather than a new report per turn.
- [ ] Update manifest/resource mirrors and nearest routers; run affected routing/resource and
  public-core/Markdown checks. Review together with Epic 2.

**G) Implementation Notes:** Missing binary/runtime -> existing tooling-readiness route; missing
login, mode capability or app message tools -> service/host preflight, not package installation.
No new runtime is required. Initial missing authentication is an expected guided setup state,
not a helper capability failure. The coordinator handles it before constructing a launch request.
The hybrid runbook separates technical execution from user-facing decisions,
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
Do not implement the new helper in the legacy Python behavior oracle merely for symmetry.
Epic 1 deliberately updates the existing onboard output and its oracle/tests together.

**D) Acceptance Criteria And Size:** M. One invocation uses the expected brief, worktree, session,
model and settings; originals survive; session mismatch/denial/failure is truthful; notification
payload preserves the exact recipient. Same-session overlapping attempts cannot both launch.

**E) Dependencies And Critical-Path Notes:** Epic 1 fixes the contract and explicitly changes the
existing onboard wording. Other existing CLI behavior remains unchanged; new helper behavior is
specified by its Go tests, not a Python copy.

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
  and matching lock ownership; interruption before child identity capture retains a blocker and
  documented owner-directed recovery. No raw normal output or implicit escalation. Use sanitized
  real CLI response/denial and message-rejection shapes to anchor fake fixtures. Exercise actual process behavior, not just mocked return values.
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
Add small absent-service/opt-in/preservation fixtures around that same recipe; do not install or
uninstall tools or sign out a real account to manufacture these states. Use a temporary member
fixture for member-to-home navigation rather than changing another live program repository merely
for a test.

**C) Files Touched:** Release identity/notes and affected mirrors under the release runbook;
execution log; consumer kit config/lock/managed content only through upgrade, local operational
reference and its existing root/index/Program routes. Private identities/evidence stay consumer-owned.

**D) Acceptance Criteria And Size:** M. Broad release checks pass; assets match accepted source;
consumer default branch adopts that release; original question and final response reach the fresh
director after yielding; it reads originals and assesses them. Stale/missing routes fail clearly.
Native-only and selected-but-unconfigured cases pass the scoped first-use checks; an existing
custom reference survives setup repetition and Kit upgrade. Record host support and any
outstanding working-copy reconciliation separately from adoption.

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
- [ ] Use bounded fresh-agent probes in fixture repositories for behavior, not document review
  alone: (a) no reference/Claude/account/preference -> ordinary native work proceeds with no
  Claude call or setup prompt; (b) selected Claude with missing CLI/login -> guides the concrete
  setup/choice without fallback, credential copying, bypass or unapproved installation; (c) feed
  an accepted native/decline choice, then another routine request -> no renewed offer or team-wide
  preference write. Cases (a)/(c) can share one small probe; (b) is a separate selected-provider
  probe, with scoped readiness fixtures and no real account changes. Check intended preference
  scope in those probes. Use an automated existing-reference fixture for (d): repeat setup and
  Kit upgrade preserve custom reference contents and leave personal choices in the ignored layer.
  Use the real installed cross-tool exercise below for positive login/preflight.
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
  Record live usage/time if exposed and distinguish cached input from new output; no invented
  budget threshold. Disclose that initial real adoption reuses an existing relay: creating one
  from scratch is fixture/procedure coverage unless independently observed, not claimed live proof.
- [ ] In a bounded fixture continuation of the existing agent exercise, request an agreed planning
  edit plus review and provide a reviewer finding suggesting another change. Pass: coordinator
  links the original, summarizes findings, gives its own assessment/recommendation and yields;
  only the agreed edit exists, with no new amendment or corrective dispatch. Then supply explicit
  correction authority: it applies covered fixes without per-finding approval. Verify the recipe
  preserves routine corrections under an existing whole-plan execution grant. No extra live
  external review or new production task is required solely for this probe.
- [ ] Extend the bounded commissioning fixture to cover default-goal selection, an explicit opt-out
  and unavailable/unverifiable activation without false claims or host-authority bypass. Qualify
  the selected Claude native-goal invocation and review/blocker handoff within the installed live
  exercise; the director reads the activation/evaluator evidence. Check Codex's native instruction
  path against exposed goal tools, without claiming the deferred reverse cross-tool recipe.
- [ ] Director presents material findings, recommended versus approved/implemented corrections and
  residual limits to the human, with the original-response reference. Verify generated artifacts
  are ignored and sufficient for recovery without requiring a duplicate authored report.
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

- 2026-10-09: activated the whole-plan assignment on the implementation branch after explicit
  commissioning. Source review follows Epics 1 and 2; release, named adoption and onboarding
  remain part of the finish line. See the execution log for actual state.

- 2026-10-09: user approved native goals by default for Codex and Claude implementation-plan
  execution; added source routes, activation/handoff requirements and proportionate proof to
  the existing delivery. This author-checked amendment ships in its planned release.

- 2026-10-09: user approved the discussion-stage report-before-amendment boundary, its canonical
  planning-review route and a bounded authority probe. This accepted rule is author-checked and
  is not retrospectively covered by the earlier independent review. No implementation is activated.

- 2026-10-09: recorded independent Ready verdict on b6132dd and closure of the amendment findings.
  This update records review status only; implementation remains draft and inactive.

- 2026-10-09: addressed the amendment review: align all existing onboarding output/mirror/test
  surfaces, specify the oracle decision and impact, clarify uncertain-lock recovery, name behavioral
  versus automated first-use proof, and bound Claude-host onboarding/relay-creation claims.
  Same-reviewer correction check is pending; the three epics remain inactive.

- 2026-10-09: added the user-approved native-default/optional-Claude first-use flow, guided tooling
  and sign-in, scoped preferences, managed defaults versus consumer-owned reference creation,
  bounded conflicting onboarding-default alignment and first-use/preservation acceptance cases.
  Same-reviewer amendment review is pending; execution remains inactive.

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
