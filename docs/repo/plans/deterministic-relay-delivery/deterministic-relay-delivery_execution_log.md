Last updated: 2026-10-30T18:20:26Z (UTC)
Created: 2026-10-30

# Deterministic Relay Delivery Execution Log

Plan: [Deterministic Relay Delivery](deterministic-relay-delivery_implementation_doc.md).

## Activation And Complete-Delivery Grant

On 2026-10-30 the user approved the plan, normal push/merge, activation and commissioning one
Claude implementer. The approved draft was local commit `87b16b3`; activation changes only the
canonical plan, this execution log and the plans index. The planning checkpoint is published and
merged before source implementation. No implementation or release success is claimed here.

One fresh Claude Opus 5.5 CLI implementer owns both epics. The commissioning Codex director owns
acceptance, actual Codex-host exercises and workspace closure. The existing ordinary relay
transports requests and original results. Private model/profile, session, recipient, host,
consumer and evidence locators are fixed in the coordination home's ignored assignment record.
The already-created isolated planning workspace is reused for the same delivery; the shared
producer checkout remains on main with its local files preserved.

The grant covers the two epics through reviewed source, normal commits/pushes/PRs and merge,
one next-unused patch release under the unchanged unsigned internal/prototype audience,
supported adoption in the named coordination home's main branch, its relay-instruction update,
the bounded harmless live proof, and publication of final closure records. No production
feature changes, cloud-resource effects, credential changes, permission widening, branch
deletion, force/admin merge or new messaging service is included.

## Review And Phase Boundaries

EP-01 starts with implementation and affected local tests, then returns a pushed reviewable
source candidate to the director. The director supplies the actual host exercise and independent
review, assesses the original evidence, and authorizes correction or continuation. The native
Claude phase goal ends at that handoff; its ending is not EP-01 acceptance or plan completion.
The same session continues EP-02 after acceptance, with its next goal aimed at the release/adoption
finish line. Confirm each goal from observable CLI output.

Run cheap affected checks during iteration. Executable changes retain the current broad candidate
gate at the coherent release boundary, followed by required public-asset smoke. Reuse applicable
evidence and rerun only invalidated checks. A host/tool or permission failure is reported
unchanged, not treated as authority to select another transport or broaden a profile.

Before the new function is qualified, the relay loads helper-prepared message arguments as one
parsed object into the native send tool. It preserves the exact receipt and records delivery.
The director checks the known attempt once at a dependent checkpoint if a notice is missing.
There is no watcher or automatic re-execution of Claude.

## EP-01 Source Candidate

This section records the first candidate, `3f4b477`, as it was handed off. The source
correction below supersedes its receipt and fallback behavior.

Implementation started on 2026-10-30 from the merged activation checkpoint on the isolated
branch `codex/deterministic-relay-delivery`. The implementer's native phase goal was confirmed
active from the CLI's goal confirmation. Its permission profile was probed once: the allowed
probe ran and the explicit ask probe was refused headlessly. The unlisted harmless `printf` probe
unexpectedly ran, apparently because the CLI treats it as a read-only command needing no
permission; an unlisted non-read-only command was refused. The profile was not changed. The
deviation is reported to the director as qualification evidence.

Changed behavior:

- New managed workflow `components/agent-interface/managed/scripts/forward-prepared-message.js`,
  installed at `.codeheart/kit/docs/agent-interface/scripts/`. It runs in the Codex
  tool-execution runtime with injected native-send and shell callbacks, binds the helper outcome
  to the authorized assignment, attempt, state root and destination, and checks that the helper
  provides `record-delivery`. It reads `delivery.json` and `message.json` through framed,
  byte-counted shell reads and claims the attempt by atomic `delivery-claim/` creation. It then
  passes the parsed message object unchanged to the send tool once, saves the native result,
  derives the status and records it with the existing `record-delivery`. After a definite
  argument rejection only, it sends one compact failure notice. Fallback outcomes are sent once
  per invocation without records. The compact result names the outcome, locators and next owner,
  never the message or raw receipt.
- Receipt interpretation: `sent` only for a success envelope (or a bare `{threadId}`) naming the
  expected recipient; `rejected` only for the known invalid-arguments envelope. Any other error
  envelope, throw, mismatch or unfamiliar shape is `uncertain`, with no notice or resend.
- Guidance: the contract gains "Deterministic Forwarding" with the loader, inputs, sequence,
  receipt table, result fields and recovery. The hybrid runbook's dispatch envelope runs the
  loader instead of manual send and record steps, and its recovery covers claims. The Codex
  operations reference, managed README script index, installed inventory and placement contract
  name the asset and its caller.
- Distribution: declared in the agent-interface component, mirrored as a packaged resource, with
  the component checksum and content-graph digest refreshed. Go embedding already covered the
  file; the helper's Go source is unchanged.
- Feedback: `node --test tests/coordination-forwarder.test.cjs` added to the cheap `feedback`
  job; public-core hygiene now also scans `.js` and `.cjs` files.

Deliberate choices for review: helper output that is not one JSON outcome (for example a helper
that failed to start) is a blocker that the relay reports in its own reply; the forwarder
composes no message for it. An outcome naming a message file without its delivery record is a
binding blocker. A failure notice follows a definite rejection even when its recording failed.
The bare `{threadId}` receipt shape is accepted, but only the tool envelope has been observed
from the native send tool so far.

Consumer impact: additive managed file under an existing component target, with a reviewed
external-action and recovery-policy change for relays (`security or safety policy change`), and
a producer validator change. No consumer scaffold, migration or new prerequisite: consumers need
no Node runtime. In-flight attempts keep their commissioned route; adoption is deliberate for new
attempts. The script is not guidance-scope eligible, so release uses the broad candidate scope.

Focused validation on macOS (Node 26, Go 1.26, Python 3.11 in the repo-local environment):

- `node --test tests/coordination-forwarder.test.cjs`: 24 passed. Recording uses the real
  helper built from this source; `exec_command` is emulated by running each generated command
  string through `/bin/sh`; the native send tool is always a fake. It covers success,
  error-shaped and unknown/mismatched receipts, throws, the incident's `thread` key, unsafe
  paths and shell-sensitive text, truncated reads, failed persistence/recording, interruption,
  concurrent and repeated entry, long-reply locators, fallbacks, and the documented loader.
- `pytest` feedback set (routing, packaging, guidance candidate/lifecycle, release identity):
  122 passed; sync/drift/init/install metadata/public-core/markdown: 43 passed; JSON schemas,
  CLI and Go parity: 59 passed.
- `go test` for manifest, state, kitfs, hash, coordination, CLI and commands packages: passed.
- Public-core, Markdown header, JSON schema, release manifest and release identity validators:
  passed.

Limits: none of this exercises the Codex host. The actual `exec_command` parameters (`login`,
`tty`, `yield_time_ms`, `max_output_tokens`), its result shape, `new Function` loading, the
text-output function name and the live native receipt remain to be confirmed by the director's
candidate exercise, together with the fresh low-context route walkthrough, primary review and
one independent source review. No merge, release or adoption has occurred.

Candidate handoff outcome: the commit was pushed, but the implementer's inline `gh pr create`
was refused by the headless permission check, although the profile lists that command; its
Markdown body contained backquotes and dollar signs. Per the assignment, the refusal was reported
unchanged and not retried in another form. The phase goal kept requiring the PR, the same blocker
was repeated until the CLI's stop-hook limit ended the turn, and the attempt's final result was
empty. The candidate, branch and evidence were retained, and the substantive handoff had already
been given earlier in that session. The source-review PR remains absent pending the owner's
decision on a body-file retry; the prepared body is kept for that route.

## EP-01 Source Correction 1

The director and the independent reviewer reviewed `3f4b477` and withheld acceptance pending two
receipt-evidence corrections. They found no other material issue; the fresh low-context recovery
walkthrough passed and changed packaged resources matched source.

- Receipt classifier: a success envelope with the expected `threadId` beside a nested
  `isError: true`, a `status` of `rejected` or `failed`, or a top-level `errors` field was
  classified `sent`. Success now requires the observed minimal shape exactly: an envelope of only
  `{content, isError}` with one `{type, text}` text item whose JSON is exactly `{threadId}` for the
  destination, or that exact payload. `rejected` requires the exact known invalid-arguments text
  in such an envelope, so a permission or security error, or any extra or conflicting field, is
  `uncertain` and cannot trigger a failure notice.
- Fallback receipt preservation: a fallback send kept no native evidence. It now writes the native
  result or thrown error to a new, uniquely named private file below
  `<state_root>/forwarder-receipts/<assignment_id>/`, returns its locator, and reports
  `receipt_problem` when it cannot be saved. It never touches `attempts/` (including the earlier
  attempt named by `attempt_exists`), records no delivery state, and still sends at most once per
  invocation.
- The contract's receipt table, fallback paragraph and result fields describe both changes.
  Markdown timestamps that were set ahead of the actual time in the first candidate are
  corrected to the original commit time or the correction time.

Focused validation: `node --test tests/coordination-forwarder.test.cjs` passed 27 tests. New
cases cover the reviewers' conflicting wrappers, permission-like errors near the rejection text,
qualified success, fallback success, rejection, unknown and thrown receipts, failed sidecar
persistence, and an unchanged `attempt_exists` attempt. Source-to-packaged equality, routing and
packaging tests, the manifest graph digest, and the public-core, Markdown, release-manifest and
release-identity validators were rerun. Unaffected Go and broad Python evidence is reused.

## Current State

- Plan: active.
- EP-01: corrected source candidate pushed for director and independent re-review and the
  Codex-host candidate exercise; unaccepted. The source PR is absent pending the owner's
  permission decision.
- EP-02 release and consumer adoption: not started; authorized after EP-01 acceptance.
- Workspace: retained by the commissioning director through implementation/review/release and
  adoption; remove only after the approved preservation and dependency checks, retaining its
  branch and useful evidence.
