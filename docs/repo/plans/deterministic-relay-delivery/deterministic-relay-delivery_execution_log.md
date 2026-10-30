Last updated: 2026-10-30T19:14:05Z (UTC)
Created: 2026-10-30

# Deterministic Relay Delivery Execution Log

Plan: [Deterministic Relay Delivery](deterministic-relay-delivery_implementation_doc.md).

## Activation And Complete-Delivery Grant

On 2026-10-30 the user approved the plan, normal push/merge, activation and commissioning one
Claude implementer. The approved draft was local commit `87b16b3`; activation changes only the
canonical plan, this execution log and the plans index. The planning checkpoint is published and
merged before source implementation. No implementation or release success is claimed here.
Planning publication is complete: the activation PR merged to `main` as `b443054`, the start
commit of the implementation branch.

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

## EP-01 Source Correction 2

The director and the same independent reviewer examined `76e346a`. The strict classifier resolved
the first finding, and fallback sidecar preservation resolved the ordinary case of the second. One
failure path remained: when a fallback receipt could not be saved, the result still said
`delivered` with the coordinator as next owner, hiding lost evidence. That case is now
`recording_failed`, keeping `derived_status` and `receipt_problem`, and returns the incomplete
handoff to the relay without any further send. Focused fixtures cover sent and rejected fallbacks
whose evidence cannot be saved; `node --test tests/coordination-forwarder.test.cjs` passed 27
tests, and mirror, digest and validator checks were rerun.

The user explicitly approved one retry of the refused PR creation through the prepared body file
under the unchanged permission profile; review, merge and release gates are unchanged. The retry
is made from this correction's head and its result is reported in the handoff.

The retry succeeded: source PR #37 opened at `0ab0503` under the unchanged profile.

## EP-01 Acceptance

On 2026-10-30 the commissioning director accepted EP-01 at `0ab0503`.

- Reviews: primary review and the retained independent reviewer passed with no remaining source
  findings, and the fresh-context recovery routing probe passed. PR feedback run `41106837142`
  passed on that exact head.
- Native host exercise:
  - The actual Codex tool-execution runtime loaded the packaged candidate and confirmed the
    native send and shell tools, loader evaluation and text output.
  - It blocked a malformed message and classified permission-error and conflicting-receipt
    fixtures as `uncertain` without live sends.
  - The actual relay then ran the same candidate once on a clearly synthetic fixture. The native
    send returned the qualified envelope, the receipt was saved, and `record-delivery` recorded
    `sent` with the recipient confirmed. The director received the prepared marker after
    yielding.
- Boundaries: the synthetic proof is candidate qualification only. It does not replace the
  post-adoption smoke with a real Claude consultation. Private locators and receipts stay in the
  coordination home's ignored assignment evidence.
- Merge: the repository rules allowed an ordinary merge without any rule change or bypass. PR #37
  merged normally to `main` as `64288fa`, pinned to the accepted head; the source branch is
  preserved.

## EP-02 Release Candidate

- Version: v0.1.39 is the next unused patch. Live tags and releases end at v0.1.38, and no open
  PR reserves a version.
- Branch: the release candidate is prepared on `codex/deterministic-relay-release` from `64288fa`.
- Identities: all v0.1.39 literals are updated (CLI, Python package, installers, bootstrap,
  release-candidate fixture, agent-interface component and standard profile), and the
  declaration checksums, graph digest and packaged mirrors are refreshed.
- Release notes: they state the executable coordination, external-action and validator impact,
  helper compatibility, the in-flight attempt boundary, the macOS live-host scope and the
  unchanged unsigned internal/prototype audience.
- Gate: the single broad candidate run, artifact verification, publication, HQ adoption and final
  live smoke remain open.

Candidate qualification:

- Release PR #38 feedback passed on candidate `7980b0c`.
- One broad candidate run, `41108029380` (all lanes), passed on that exact head: macOS,
  Windows, Ubuntu semantic and the Git 2.43 proof.
- Retrieved artifacts were verified:
  - Each archive digest matches its sidecar and the catalog. Each pack-manifest digest matches
    the catalog, and every payload checksum matches its file.
  - The binary and content identities agree.
  - The installers, bootstrap and notes are byte-equal to source.
  - Both binaries embed the forwarder bytes qualified on the Codex host.
- The director accepted the candidate for publication.

## EP-02 Publication

- Merge: release PR #38 merged normally to `main` as `5510bec`, pinned to `7980b0c`. The merge
  tree is identical to the qualified candidate tree, so no release input changed after
  qualification.
- Publication: v0.1.39 was confirmed unused and then published from the merge commit:
  - Release: [Operating Kit v0.1.39 — Deterministic relay forwarding](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v0.1.39)
  - Tag `v0.1.39` points to `5510bec`.
  - Audience: unchanged unsigned, unnotarized internal/prototype HTTPS-plus-SHA256.
- Published assets: both packs and their sidecars are the qualified candidate artifacts,
  unchanged. The public catalog differs from the candidate catalog only in its download URLs.
  The 14 published assets and their digests:

| Asset | SHA-256 |
| --- | --- |
| `codeheart-operating-kit-0.1.39-macos-universal.zip` | `2beb0a6e2deb3a09e18959f14031d60493b0e5acfc0c23d45b8a884805a15f42` |
| `codeheart-operating-kit-0.1.39-windows-x64.zip` | `436c584874b0d37054c75cb714c5334d8f15f185ecd7210021ac40d86d0fcb20` |
| `release-catalog-0.1.39.json` | `f29f043c18b3d6e667c13de93d47c36c44dd4ae79e4c2694dc8ba0355556ee29` |
| `bootstrap.md` | `8269a8f4f3dc696639f305323986e97aed339ce6e593e1b885ac702f41585f43` |
| `install.sh` | `6c423066d6e24e579a6f3548420688a3b8dcded1da32156326d68b597151d0ba` |
| `install.ps1` | `488a1f3b2ac4281338b744d1ad5521292081ee288d70d47ecf1ef86acf20659f` |
| `release-notes.md` | `cf1bcc971fd08b7473b9cfe4d9159b2b6fc3efd0ac67b47cc3b7bf365986cc76` |

  Each of these also has a `.sha256` sidecar, for 14 assets in total.
- Public download check: all 14 public downloads are byte-identical to the verified staged
  assets.
- Released-asset smoke: run `41113459170` with explicit `v0.1.39` passed on macOS and Windows.
  Both public installers rejected a wrong checksum, installed v0.1.39 from the public catalog,
  and then onboarded and checked a fresh consumer.
- Separate from publication: HQ adoption, the real Claude relay smoke and closure remain open.

## Current State

- Plan: active; planning publication complete.
- EP-01: accepted; source merged to `main` as `64288fa`.
- EP-02: v0.1.39 is published and its public assets are validated. Consumer adoption in the
  named coordination home, the actual-relay live smoke with a real Claude consultation and the
  final closure are pending; they are not complete.
- Workspace: retained by the commissioning director through implementation/review/release and
  adoption; remove only after the approved preservation and dependency checks, retaining its
  branch and useful evidence.
