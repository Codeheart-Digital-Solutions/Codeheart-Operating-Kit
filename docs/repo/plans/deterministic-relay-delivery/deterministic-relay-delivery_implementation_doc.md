Last updated: 2026-10-30T17:38:12Z (UTC)
Created: 2026-10-30
Status: draft

# Deterministic Relay Delivery

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.deterministic-relay-delivery
  kind: implementation
  purpose: Forward prepared cross-tool notifications through Codex without agent reconstruction, retain truthful receipts, and verify released adoption in the named coordination home.
  first_cataloged: 2026-10-30T17:31:22Z
  catalog_metadata_updated: 2026-10-30T17:38:12Z
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.cross-tool-agent-coordination
```
<!-- END CODEHEART PLAN METADATA -->

## Overview

Replace the relay's manual reconstruction of native message arguments and delivery receipts with
one Kit-owned JavaScript workflow, executed through the Codex host's existing tool-execution
environment. Preserve the current CLI invocation, original replies, recipient authority and
coordinator review. Finish with a published patch, adoption on the named coordination home's
`main`, and a verified handoff through its actual relay.

The user requested this draft after a bounded incident investigation and read-only feasibility
probes on 2026-10-30. A separate discovery document is unnecessary. This document authorizes no
implementation, messages, push, merge, release or adoption until execution is commissioned.

Authoring checkpoint: the configured canonical catalog accepts this draft; relative links,
required structure, Markdown timestamps and public-core checks pass. The local branch
`codex/deterministic-relay-delivery-plan` and its isolated workspace are retained by the drafting
director for review and authorized publication. The shared producer checkout remains on `main`.

## Essential Context Reference Files

| File | Purpose |
| --- | --- |
| [Existing coordination plan](../cross-tool-agent-coordination/cross-tool-agent-coordination_implementation_doc.md) | Completed predecessor; preserve its established boundaries. |
| [Coordination runbook](../../../../components/agent-interface/managed/runbooks/coordinate-cross-tool-task.md) | Hybrid commissioning, original-result return and recovery. |
| [Coordination contract](../../../../components/agent-interface/managed/reference/cross-tool-coordination-contract.md) | Helper outcomes, message files, private evidence and delivery states. |
| [Codex operations](../../../../components/agent-interface/managed/reference/codex-task-operations.md) | Ordinary-chat recipients and temporary relay workers. |
| [Invocation](../../../../internal/coordination/invoke.go), [records](../../../../internal/coordination/records.go), [delivery recorder](../../../../internal/coordination/recovery.go) | Existing mechanics to reuse; the recorder currently trusts the caller's status and reported recipient. |
| [Recipe maturity](../../../../components/agent-interface/managed/reference/operational-recipe-maturity.md), [script promotion](../../../../components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md) | Justify an L2 workflow script with a narrow owner, contract and tests. |
| [Runbook authoring](../../../../components/agent-interface/managed/reference/runbook-authoring-standard.md), [routing](../../../../components/agent-interface/managed/reference/operation-routing-and-dispatch.md) | Preserve audience, authority, readiness and error routes. |
| [Placement](../../reference/placement-contract.md), [consumer impact](../../reference/consumer-impact-classification.md), [change](../../runbooks/change-operating-kit.md), [release](../../runbooks/release-operating-kit.md) | Managed asset distribution, focused feedback and release gates. |

## Table Of Contents

- Section 1: outcome, incident and verified baseline.
- Section 2: architecture, boundaries, authority and validation choices.
- Section 3: two ordered delivery epics.
- Section 4: deferred capabilities and support limits.

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

A relay invokes one released function with a helper outcome and the already authorized return
destination. The function forwards the helper-prepared arguments unchanged, preserves the native
receipt, and records a truthful result. A malformed recipient key is rejected before an external
send. Re-entry or interruption cannot silently trigger another send of the same attempt.

Completion requires the installed function to deliver an original reply or its exact locator to
the intended ordinary coordinator after it yields. Verify both the accepted send and actual
coordinator receipt; neither implies review or acceptance of the substantive work. Controlled
failure must remain visible and recoverable without rerunning the implementer.

## 1.2 Project And Problem Context

An observed relay worker read a valid `message.json` with `threadId`, then manually called the
native tool with `thread`. The tool rejected the arguments. The implementer's reply was retained,
but the director received no notice; the relay reported the failure only in its own chat. An
earlier correctly formed send to the same destination had succeeded. This was a forwarding defect,
not evidence of lost implementation output or a general messaging outage.

The accepted direction is deterministic forwarding and receipt handling, with proportionate
validation. Private conversations, paths and raw receipts remain in the consumer's existing
ignored assignment evidence. Public tests use synthetic identifiers and sanitized shapes.

## 1.3 Current State Analysis

- The v0.1.38 helper emits `message_file` or `fallback_message`. Long answers already use exact
  file/line/field locators; shortening or summarizing the implementer's answer is unnecessary.
- `record-delivery` moves an existing pending record to a terminal state once. It does not call
  the host messaging tool, interpret its receipt, or prevent two senders from racing before recording.
- The actual relay worker used Codex's JavaScript tool-execution surface. That surface exposes
  native send and shell-execution tools; its V8 runtime has no direct filesystem or Node imports.
  The separate Node REPL has filesystem access but no native chat-tool bridge.
- Read-only probes confirmed callable native tool bindings, loading a function from source,
  unchanged object forwarding to a fake sender, original simulated receipt preservation, and
  rejection of `thread` before sending. No live message or receipt record was changed by the probes.
- These facts establish feasibility on the observed macOS host. Loading the installed asset in
  the actual relay, real receipt interpretation, interruption and coordinator wakeup still need
  the explicit qualification below. Public [Programmatic Tool Calling documentation](https://developers.openai.com/api/docs/guides/tools-programmatic-tool-calling)
  explains the general runtime; it does not guarantee desktop delivery or persistent execution.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

```text
components/agent-interface/
  component.yaml                                           # add managed script declaration
  managed/scripts/forward-prepared-message.js               # create one host-executed workflow
  managed/README.md                                        # add compact script role/caller index
  managed/runbooks/coordinate-cross-tool-task.md            # replace manual send/receipt steps
  managed/reference/cross-tool-coordination-contract.md      # document adapter input/results/recovery
  managed/reference/codex-task-operations.md                 # link the installed invocation route
src/codeheart_operating_kit/resources/components/agent-interface/
                                                           # mirror changed managed resources
tests/coordination-forwarder.test.cjs                        # create focused built-in Node tests
tests/test_packaging_resources.py                           # cover declared installed asset
.github/workflows/validate.yml                              # add only focused script feedback hook
docs/repo/reference/placement-contract.md                    # clarify host-script placement
docs/repo/plans/deterministic-relay-delivery/
  deterministic-relay-delivery_implementation_doc.md         # this plan
  deterministic-relay-delivery_execution_log.md              # create on activation, not at drafting
manifest.yaml, release-notes.md, existing version mirrors    # update at release boundary
```

Install the script under the existing agent-interface managed area at
`.codeheart/kit/docs/agent-interface/scripts/forward-prepared-message.js`. It is an optional,
explicitly invoked host asset, never an install hook. Document that placement; add no consumer
scaffold or global installation. Use the existing managed resource and packaged-mirror mechanisms.
The companion consumer changes its existing coordination reference and ignored helper selection
through supported adoption; its exact private files and chat locators belong in the activation
assignment. Routine test organization may be refined without changing this architecture.

## 2.2 Open Questions And Assumptions Requiring Clarification

No unresolved design blocker remains. The following are action-time checks with chosen defaults:

- OQ-1 — **BLOCKER: no. Affects: EP-01, EP-02.** Verify that the selected relay still exposes
  JavaScript execution, the native send tool and the approved shell tool. Use those supplied
  capabilities. If absent or changed incompatibly, stop at a concrete host-support blocker; do
  not substitute a private app endpoint, new server, account connection or another recipient.
- OQ-2 — **BLOCKER: no. Affects: EP-02.** Select the next unused patch from live source/tags/releases
  at candidate preparation; do not reserve a version in this draft. Default to broad candidate
  scope because executable behavior is changing. Preserve applicable evidence after corrections.
- Assumptions: the named consumer's existing macOS relay remains the first supported host; its
  current helper understands `record-delivery`. The existing message authority is checked at
  commissioning, not inferred from a file's existence. A whole-plan grant must include the named
  live probe and consumer adoption before those steps run.

## 2.3 Architectural Decisions With Reasoning

**One L2 workflow, executed by the host.** Package the reviewed function as a managed file. A
small documented loader reads that trusted file through the supplied shell tool and evaluates it
in Codex, injecting native-send and shell callbacks. It passes outcome data as data, not source
code. Await the complete operation in that execution cell. The agent selects the authorized
assignment and invokes the loader; it does not transcribe recipient fields, message content,
receipt contents or delivery status. No new CLI command, daemon or MCP server is needed.

**Inputs and phases.** Accept the just-finished helper outcome, verified helper executable,
assignment/attempt identity, approved destination and expected state root. Load only the message
file named by that invocation; validate its binding to those inputs and its exact allowed keys
`threadId`, `hostId`, `prompt`, with non-empty strings. Pass the parsed object itself to native
send. Support the existing inline fallback outcome without inventing an attempt or touching an
earlier attempt's records. Never resume Claude or scan for an alternative result from this function.

**Small local claim before sending.** For a recorded attempt, check the existing delivery state,
then exclusively create a private `delivery-claim` directory beneath that exact attempt before
sending. Use the approved shell's atomic directory creation; concurrent calls and re-entry must
refuse an existing claim. Retain the claim after completion or interruption. Existing sent,
rejected or uncertain records are not reset. A pending claimed attempt means sending may have
started and requires owner reconciliation; age or a dead process alone cannot authorize replay.
Fallback outcomes without an attempt have no durable once-only guarantee: send once within that
invocation, never automatically replay, and report that limitation explicitly.

**Interpret the actual receipt.** Preserve the returned native tool object in a private receipt
file using safe argument encoding, then derive status from it. A recognized success must include
the expected recipient identifier and no error. Known pre-delivery argument rejection is
`rejected`; exceptions, unfamiliar receipts, recipient mismatch and possible post-send failures
are `uncertain`. An error envelope must never count as success. Reuse `record-delivery` with the
derived status and recipient; check its exit and structured result. If receipt persistence or
recording fails, expose `recording_failed` and possible delivery rather than claiming completion.
This is not an atomic transaction with the chat service and makes no exactly-once guarantee.

**Bounded failure visibility.** Return a compact result identifying assignment/attempt, delivery
state, receipt/attempt locator and next owner; keep raw output private. Before any send, malformed
input or missing capabilities produce a blocker. After a definite mechanical rejection, the
workflow may send one compact failure notice through the same authorized native tool to the
validated coordinator, retaining its separate receipt and preserving the original rejection.
It does not resend the original reply. Never do this after a permission/security refusal, an
uncertain send, or a missing/invalid destination; never switch tools or recipients. If that notice
also fails, return an explicit blocked handoff to the relay with no further attempts. At its next
dependent decision, the coordinator reads the known attempt once. No watcher or claim of reliable
wakeup when the host channel is unavailable is introduced.

**Compatibility and support.** Keep existing helper requests, original output and old delivery
records readable. The new adapter owns its claim and receipt files; older in-flight attempts
continue under their commissioned route. MacOS Codex-host execution is the qualified live scope.
Pure JavaScript fixture tests may run in Node using its built-in test runner; Node is a development
tool, not a new consumer prerequisite. Shell paths and arbitrary receipt text must be safely
quoted, never interpolated as executable code. Local missing-tool blockers follow the existing
tooling-readiness runbook; module service setup and permission-profile changes are outside scope.

**Guidance and impact.** The hybrid coordination runbook retains its compact intention block,
user/setup flow and owner judgment while naming the script for the mechanical phase. Agent-facing
Codex and coordination references describe its invocation, output, failures and recovery. The
existing coordination route remains the owner. One fresh low-context walkthrough checks that a
transport failure routes here without guessing another tool. This is an executable coordination
change with an additive managed asset and reviewed external-action/recovery-policy implications;
record those impacts and adoption notes, not an instruction-only classification. Consumer-owned
runbooks beyond the named coordination arrangement are outside this delivery.

## 2.4 Proposed Commissioning And Delivery Boundary

The commissioning director owns acceptance; one implementer owns both epics. At activation, bind
the named coordination home (`main`), actual relay/coordinator probe chats and host, supported CLI
and private evidence root in the existing assignment record. The intended consumer is the home
whose relay exposed this incident, not an organization-wide rollout. Adopt its managed Kit,
existing coordination reference and relay dispatch instructions together.

Proposed whole-delivery authority covers implementation, focused tests, coherent commits, normal
push/PR updates, ordinary reviewed merge, next-patch publication under the unchanged
unsigned internal/prototype audience, the named consumer adoption and harmless live messages for
verification. It includes final plan/log publication and safe cleanup of this assignment's clean,
integrated worktrees while retaining useful evidence and branches. These effects require an
explicit whole-delivery commission; today's drafting request supplies none of them.

Use isolated implementation/publication workspaces and preserve shared checkouts and live
assignments. Establish report-back and review checkpoints once. The director performs primary
review plus one independent review of the coherent source/failure semantics, retaining that
reviewer for corrections. A native phase goal ends at the source-review handoff, then at complete
release/adoption. No per-checkbox approval. Material permission, host-support, scope or release
policy changes return to the owner. Until corrected delivery is adopted, a failed report-back is
visible in the assignment and relay; the known attempt remains the recovery source.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 — Deterministic delivery | Packaged workflow and instructions pass focused contract tests, actual-host candidate exercise and source review. | M | Execution commission and current host capability check. |
| EP-02 — Release and verified adoption | Accepted patch is published, named home main and relay use it, real handoff is verified, records and workspaces are closed. | M | EP-01 acceptance and named release/adoption/probe authority. |

## EP-01 — Deterministic Delivery

### A) Epic ID, Title, And Outcome

EP-01 provides an executable host-side send-and-record workflow with guarded re-entry and truthful
failure reporting. Completion requires behavior, distribution and host-execution evidence.

### B) Scope

The script, its managed declaration/mirror, focused tests, minimal feedback hook and affected
coordination guidance from Section 2.1. Reuse existing CLI behavior; do not build a scheduler or
expand messaging permissions. Preserve existing attempt/output contracts.

### C) Files Touched

The script, component manifest, managed README/runbook/references, resource mirror,
`tests/coordination-forwarder.test.cjs`, affected packaging tests, `validate.yml` focused hook,
placement clarification and this plan's execution log.

### D) Acceptance Criteria And Size

Size M. Prove unchanged arguments and prompt, receipt-derived status, no send on invalid input,
one claimant per recorded attempt, original evidence preservation and explicit incomplete handoff.
Required fixtures cover success, error-shaped response, unknown/mismatched receipt, unsafe path or
shell-sensitive text, failed persistence/recording, concurrent/repeated entry, interrupted delivery,
long-reply locator, and no-record fallback. A permission refusal cannot cause an alternative send.
An actual Codex-host candidate exercise loads the packaged function and records the native result.

### E) Dependencies And Critical-Path Notes

Check the host bridge before substantive implementation. The existing incident and feasibility
evidence are retained; do not rerun the feature implementation. The candidate host exercise uses
only the named harmless-message authority and precedes release acceptance.

### F) Tasks Checklist

- [ ] Implement the workflow, deterministic loader and narrowly injected host callbacks; validate
      destination/outcome binding, safely claim a recorded attempt, forward once and retain receipts.
- [ ] Implement receipt classification, existing-helper recording and bounded failure-notice
      behavior, with conservative interruption and re-entry outcomes from Section 2.3.
- [ ] Declare and mirror the installed asset; update the compact role index, placement guidance,
      hybrid runbook and agent references with the exact invocation and recovery boundaries.
- [ ] Add and run focused fixtures and packaging checks; integrate only those script checks into
      cheap feedback. Run applicable public-core, Markdown and source/mirror equality checks.
- [ ] Exercise the candidate through the actual host, including successful native receipt and
      controlled non-live failure cases; complete the fresh-context route walkthrough and primary
      plus independent source review. Record evidence, limitations and the reviewed candidate.

### G) Implementation Notes

Use fake native-send callbacks for negative cases; do not send malformed requests or duplicates to
real chats just to test rejection. For interruption, simulate failure after a send but before local
record completion and show that re-entry refuses. Claim cleanup/recovery is a conscious owner
decision, not a retry loop. Keep process output and receipts out of public commits.

### H) Open Questions

OQ-1 only; its capability check is a stop condition, not permission to invent another transport.

## EP-02 — Release And Verified Adoption

### A) Epic ID, Title, And Outcome

EP-02 completes delivery when the released workflow is installed on the named home's main branch,
the actual relay uses it, a yielded coordinator receives the original result, and closure is shared.

### B) Scope

Producer release identities/assets, required candidate and public-asset checks, the one named
consumer adoption, relay instructions, live smoke and final plan/log/index/workspace closure.

### C) Files Touched

Existing release identity and notes files, this plan/log and `docs/repo/plans/README.md`; in the
named consumer, generated Kit state and the existing coordination reference/local helper settings.
Exact private paths are fixed in the activation assignment, not copied into public Kit records.

### D) Acceptance Criteria And Size

Size M. Ordinary source merge and next-patch release have their required evidence. Published
assets contain the reviewed script; upgrade preserves consumer content. Named home remote `main`
has adopted, selected helpers and installed resources agree, and an actual relay invocation sends
an original harmless Claude result to the coordinator after it yields. Verify recipient history
and original reply/locator, not just `sent`. Failure evidence and all pre-existing attempts remain
available. State explicitly if host wakeup requires a user turn; such a result does not meet the
automatic-handoff finish line and needs owner disposition before completion.

### E) Dependencies And Critical-Path Notes

EP-01 accepted. Select and freeze the release candidate once, using current release authority and
required gates. Broad candidate validation is currently required for executable changes; this plan
does not alter that policy. Run it at the coherent candidate boundary and rerun only invalidated
lanes/checks after scoped corrections. Routine commits and document closure do not recreate it.

### F) Tasks Checklist

- [ ] Select the unused patch, prepare release/impact/adoption notes and identities; complete
      normal PR feedback, reviewed merge and required candidate qualification on applicable inputs.
- [ ] Publish under the covered audience, verify downloaded assets and required public smoke,
      and preserve the exact candidate/source/release evidence relationship.
- [ ] Upgrade the named home through the supported lifecycle, select verified matching helper and
      script, update the existing relay invocation instructions, and merge the isolated adoption
      PR. Reconcile the usual checkout without disturbing active assignments or local files.
- [ ] Run the bounded original-reply handoff through the actual relay after the coordinator yields;
      verify native receipt, private delivery record and destination history. Confirm the relay
      calls the released function rather than reconstructing arguments. Assess with the director.
- [ ] Publish final plan/log/index status with source, release and adoption distinguished. Remove
      only clean integrated assignment workspaces with no remaining dependency; otherwise record
      the exact retained workspace, owner and next resolving event. Preserve branches and evidence.

### G) Implementation Notes

Use one harmless consultation for the released end-to-end smoke; negative cases remain local
fixtures. Do not involve production feature work, AWS resources or credentials. Adoption is
deliberate for new attempts; no mid-flight migration or replay of the historical failed handoff.
The director's acceptance concerns delivery mechanics, not the implementer's product decisions.

### H) Open Questions

OQ-1 and OQ-2 are rechecked at their stated boundaries. A changed enforced release gate or missing
publication authority returns a concrete owner decision, preserving successful evidence.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Additional hosts, Windows live Claude/Codex coordination and reverse-tool operation: qualify
  only when there is a named consumer need. Distribution on a platform is not live-host support.
- Automatic redelivery, queues, watchers, a new messaging service and exactly-once delivery: not
  selected. Revisit only with repeated evidence the bounded workflow and owner recovery cannot
  serve actual use; no hidden delivery commitment remains in this plan.
- Other consumer rollouts and general release-cost optimization remain separately owned.

## 4.2 Future Considerations

Host tool availability and receipt schemas can change. An unrecognized schema must produce an
explicit uncertain outcome and focused compatibility work, rather than an optimistic success.
Keep the function portable through injected capabilities while qualifying only demonstrated hosts.

# Revision Notes

- 2026-10-30: Drafted from the accepted deterministic-forwarding direction, observed recipient-key
  failure and read-only host feasibility probes. Defined two complete-delivery epics, local claim
  and receipt semantics, bounded failure visibility, candidate validation and named-home adoption.
  Execution remains inactive; drafting authorizes only the local planning checkpoint.
