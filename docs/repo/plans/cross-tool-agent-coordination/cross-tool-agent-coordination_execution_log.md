Last updated: 2026-10-09T20:48:19Z (UTC)
Created: 2026-10-09

# Cross-Tool Agent Coordination Execution Log

Plan: [Cross-Tool Agent Coordination](cross-tool-agent-coordination_implementation_doc.md).

## Approved Review-Responsibility Amendment

At the first correction handoff, the director recorded the user's approved bounded amendment:
the commissioning agent performs the primary review against outcome, scope and evidence, uses
additional independent scrutiny when complexity, risk, authorship or a binding gate warrants it,
and remains accountable for acceptance. This replaces an automatic extra-reviewer default in
generic guidance; no tool/model or second chat is inherently required. Existing independent
source findings and correction proof for this delivery remain useful and are retained.

The amendment is in the active plan and pending source-guidance implementation. It adds no helper
mechanism or release. The corrected candidate 685515578ca8c2ba7df4d77e312d320e108c695b has passed
PR feedback. The director read the original correction report and inspected failure finalization,
replay refusal, local-settings guidance and original selected-host new/resume/denial results;
observed CLI version is 2.1.286. Source acceptance still awaits the guidance amendment, primary
review and focused follow-up on the earlier independent findings. No merge or release is accepted.

## Activation And Assignment

The user commissioned the complete three-epic plan on 2026-10-09 after planning PRs 25 and 26
merged. Source starts at main commit 62c350606e13b5de0966ab44a9f5ecc9fd752c16 on
`codex/cross-tool-coordination-implementation`. The activation checkpoint changes only the
plan, this log and its index; source implementation follows successful normal publication.

A fresh Claude Opus 5.5 CLI implementer owns the whole plan, with the Codex commissioning
director accountable for acceptance. The existing ordinary relay and temporary transport workers
only deliver prepared requests and original results. Private session/recipient/evidence locators
are recorded in the commissioning home's existing handover record, never in this public source.
Independent source review after Epics 1 and 2 remains separate from implementation authorship.

The approved native-goal default applies. The first execution goal ends at the source-review
handoff (or a faithfully reported blocker); after director acceptance the same session continues
the remaining work under another native invocation. This bounds the evaluator at a real
dependency rather than letting it spin or cross the review gate. Such a goal ending does not
complete the plan. Actual activation and handoff evidence remain to be observed.

## Delivery Boundary

Authorized delivery includes coherent commits, normal feature-branch pushes and a delivery PR
to this repository's main; independent review and passing applicable checks precede accepted
normal merge. One next-unused patch release uses existing macOS universal/Windows x64 validation
and the unchanged unsigned, unnotarized internal/prototype HTTPS-plus-SHA-256 audience.
No version is reserved at activation. Broad candidate qualification is required for executable
changes. No force/admin merge, direct main push or repository policy change is authorized.

The named private coordination-home consumer receives the release on its main branch and its
explicitly assigned clean active worktree; other production consumers and unrelated dirty work
are excluded. The private assignment identifies that repository, branch and worktree. The planned
temporary onboarding director, scoped test session and report-back are within the commission;
the director handles native chat operations using the actual user authority.
No AWS effects, additional spending commitments, global permission changes or unrelated features.

Impact includes managed instructions, additive CLI behavior, existing onboarding output and
security/safety-policy guidance. Adoption creates a consumer-owned reference by guided opt-in;
there is no schema migration or automatic scaffold ownership change. Release notes and explicit
permission/source review are required.

## Source Phase: Epics 1 And 2

The implementer's readiness phase returned READY with the native goal confirmed from the
host's goal confirmation and evaluator hook. The source phase ran under a second native goal
that ends at this source-review handoff.

Epic 1 delivered the hybrid `coordinate-cross-tool-task.md` runbook, the
`cross-tool-coordination-contract.md` reference (request, attempt records, delivery, recovery,
brief and permission-profile examples) and the opt-in `templates/agent-coordination-reference.md`
starter. Generic coordination, both tool guides, the execution and planning-review runbooks now
carry the native-tool default, native implementation goals by default and the report-before-
amendment boundary. Onboarding preserves the user's model, reasoning, speed and permission
choices across the runbook, context contract, bootstrap, Go output and Python oracle.

Epic 2 delivered `codeheart-operating-kit coordination` in the Go CLI with standard-library
mechanics only; no Python helper was added.

Meaningful divergence and safe defaults:

| Item | Decision |
| --- | --- |
| Command shape | One `coordination` group with `invoke-claude`, plus `record-delivery` and `release-lock`. The plan's delivery-state and deliberate-release requirements need deterministic writers; hand-edited JSON would be fragile. The contract and scope are otherwise unchanged. |
| Output format | Fixed `--output-format stream-json --verbose`, so partial progress survives interruption; the reply is the final `result` record. The CLI version comes from its init record, not help scraping. |
| Optional flags | `permission_prompts`, `effort`, settings, tool lists, extra directories, MCP and Chrome isolation are passed only when the request sets them. An older CLI rejects `--permission-prompts`; that returns `unsupported_invocation` with original stderr and is not retried. |
| Message file | Uses the tested host's send-message argument names (`threadId`, `hostId`, `prompt`) so relays load it unchanged. |
| Extra wording fix | `draft-implementation-plan.md` still required an explicit request for goals; aligned with the approved default. |
| Pre-existing test drift | Upgrade, release-pack and Go/Python parity tests pinned the 0.1.32 literal and failed at HEAD because guidance-scope releases skip the native suite. They now follow the built version. The release-candidate asset-name fixture remains a release-identity literal for the release step. |
| Reviewer | No reviewer agent was started by the implementer; the director obtains the independent source review. |

Validation on macOS: `go vet ./...` and `go test ./...` pass; `go test -race` and a Windows
`go vet` cover the new package. Python routing, packaging, onboarding, install metadata,
Go/Python parity, schema, release identity and manifest checks pass; the full suite has one
release-candidate fixture failure that the release step owns. Public-core and Markdown pass.

Fake-process tests cover quoting, spaces and stdin fidelity; exact argv; new and resumed
sessions; wrong and empty sessions; malformed, empty, error, nonzero and unsupported-flag
outputs; denial metadata; long-reply references; replay refusal; two concurrent assignments;
same-session collision without a second launch; helper interruption; a killed helper while the
child keeps writing; guarded release refused while processes live or ownership differs and
permitted after verified exit; interruption before child identity requiring a manual statement;
and pending/sent/rejected/uncertain delivery with recipient comparison.

Earlier live checks on the macOS host with the built helper and an older Claude CLI (2.1.153):
a new consultation and same-session continuation returned the original replies; a harmless
denied write returned its denial metadata unchanged; the older CLI rejected the optional prompt
flag as `unsupported_invocation`. A controlled harmless child kept running across separate tool
calls; a real helper termination left `helper_interrupted`, retained the lock and complete child
output, refused early release and released after verified exit.

Impact: managed instruction additions, additive CLI command, changed onboard output and
security/safety-policy guidance (permission profile and denial handling). No schema migration
or automatic scaffold; the consumer reference is created by guided opt-in. Release notes and
explicit permission-wording review are required at the release step.

## Source Review And First Correction Cycle

An independent read-only reviewer judged candidate `e0433f2` **Needs improvement** for source
acceptance. The director read the original, reported its own assessment and authorized one
bounded correction cycle. Director source acceptance is not yet given.

| Finding | Director disposition | Correction |
| --- | --- | --- |
| M-1: failures after the attempt existed wrote no message or delivery record; the relay envelope loaded a fixed path and could resend an earlier result on `attempt_exists` | Fix through shared finalization; never resend on replay; narrow fallback only | One finish step writes the record, message arguments and pending delivery for every post-creation outcome, including `session_locked`, `launch_failed` and `helper_interrupted`. It releases only a lock this launcher owns with no possible child. Output names `message_file` and `delivery_file` only when written. `persistence_errors` reports failed writes, and the exit status is nonzero. `fallback_message` uses the request's own return destination for failures before an attempt exists, a failed message write and `attempt_exists`. A refused replay changes nothing in the earlier attempt and its notice says it is not the earlier result. The relay envelope loads only the `message_file` named in this output. A lock whose content write failed is removed by its creator. |
| M-2: no durable place for the approved executable and profile | Agree on discoverability; keep machine paths out of committed references | The template records the approved CLI distribution, per-role profile, mode and `permission_prompts`, plus the ignored local settings route. The runbook resolves exact values from approved local or assignment records, verifies the executable and its observed version, and forbids unapproved `PATH` substitution. It asks only for genuinely missing decisions. |
| M-3: no proof on the selected host's CLI | Run the built helper with the assigned 2.1.286 executable and `claude-opus-5-5` | See the evidence below. Init and result fixtures are now sanitized 2.1.286 shapes. The unsupported-flag fixture remains labeled 2.1.153 evidence. |
| L-1: `release-lock` dead ends for a reused process ID or an empty lock | No free-text bypass; document owner-directed options | Refusals stay fail-closed, and tests prove a statement cannot override a live launcher, a live child or an unreadable lock. The runbook names the owner options: record actual ownership and exit verification, then remove that one lock deliberately; or hand over to an explicitly new session that keeps the uncertain attempt and reconciles effects first. No automatic recovery is promised. An atomic link-based lock claim was not added, because hard-link support is not portable across all consumer filesystems; the remaining window is the instant between exclusive create and write. |
| L-2: `preparing` missing from recovery | Add | Added with `launching` and `running`. |

Also during this cycle, the selected 2.1.286 CLI lists permission modes as `manual`,
`acceptEdits`, `auto`, `dontAsk`, `plan` and `bypassPermissions`. It still parses `default`.
The helper now passes `manual` through as well, and still refuses `bypassPermissions`.

Additional documentation:
- Windows `.cmd` and `.bat` wrappers are documented as unsupported.
- Optional effort support has an argument-fidelity test covering every optional field. The
  2.1.286 help lists `--effort` with low, medium, high, xhigh and max. No effort override is used
  by the current commission.

Not adopted:
- A brief digest in `attempt.json` (optional suggestion).
- Any change to reviewer-selection policy.

Selected-host evidence, retained privately in ignored state:
- **Setup:** the built helper with the assigned app-bundled CLI, observed version 2.1.286, and
  model `claude-opus-5-5`. Permission mode `default` with `permission_prompts` `none`, tools Read
  and Write, Read allowed, MCP and Chrome disabled, no profile file and no effort override.
- **New consultation:** captured the exact reply, with requested and actual session equal.
- **Continuation:** resuming that session returned a reply that showed memory of the earlier turn.
- **Denied write:** a deliberately denied Write to a disposable directory returned
  `permission_denials` [Write] with exit 0 and terminal reason `completed`. The directory stayed
  empty. The denial was not retried.
- **Locks:** none remained after the three attempts.

The commissioning session's own 2.1.286 print-mode stream showed `Goal set` for a brief opening
with `/goal`. That was the native-goal activation evidence; it was not rerun.

Validation after corrections:
- `go vet ./...` and `go test ./...` pass on macOS, and the new package passes under `-race`.
- Windows `go vet` and test compilation succeed for the new package.
- Python: 262 pass. The remaining failure is the release-candidate asset-name fixture.
- Release identity and manifest validators pass after refreshing the standard-profile digest.
- Public-core and Markdown checks pass.

## Primary-Review Guidance Implemented

The approved rule is in the managed guidance:
- The commissioning agent is the primary reviewer. It examines delivered work against outcome,
  scope and evidence, in proportion to risk, and owns acceptance.
- Forwarding another agent's verdict is not a review.
- Additional independent review is selected for complexity, risk or the commissioner's own
  implementation involvement, or when a user, assignment or binding gate requires it.
- Self-review is never called independent.
- The rule is tool-neutral and needs no extra account, model, relay or chat.
- A separate reviewer's model follows approved choices and never dictates the commissioner's own.

Changed routes:
- `execute-implementation-plan.md`: the ordered per-epic steps, review checkpoints and log shape.
- `agent-task-coordination.md`.
- The cross-tool runbook and role contract.
- One relay-review sentence in `codex-task-operations.md`.
- Contradictory reviewer wording in `review-planning-document.md`,
  `draft-implementation-plan.md`, `operation-routing-and-dispatch.md` and the
  `handle-routine-change.md` example.
- Mirrors and the standard-profile digest.
- The discovery received a dated pointer note; its findings are unchanged.

Unchanged:
- The discovery reviewer gate for high-risk decisions.
- The Claude guide, which had no contradiction.
- Producer release and source-review gates.

Proof:
- A routing assertion over the packaged guidance covers ordinary commissioner review, justified
  additional review assessed by the commissioner, and authorship-conflict self-review labeling.
  It also checks that the old automatic-reviewer and same-model wording is gone.
- Routing, packaging, release-identity, sync, guidance-candidate, manifest, public-core, Markdown
  and plan validation pass.
- Executable code is unchanged, so earlier helper and live proof are retained.

For this delivery, the retained independent reviewer performs the focused follow-up. The
director's primary source review and acceptance are pending.

## Current State And Remaining Evidence

Epics 1 and 2 are implemented, with the first correction cycle applied. They await the same
reviewer's focused follow-up and director acceptance.

Required Epic 3 gates, none waived:
- The release-candidate asset-name fixture test must be green at the release candidate.
- Native Windows Go tests and release gates must pass.
- Codex relay output-yield/resume must be observed with the installed helper.
- Release notes and permission-wording review.
- Named consumer adoption, first-use probes and the fresh-director exercise.

The plan stays active.
