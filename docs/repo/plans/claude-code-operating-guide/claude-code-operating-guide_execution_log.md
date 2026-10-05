Last updated: 2026-10-05T09:41:09Z (UTC)
Created: 2026-10-05
Plan: claude-code-operating-guide_implementation_doc.md

# Claude Code Operating Guide Execution Log

## Commissioning And Activation

The accountable owner approved execution of EP-01 and EP-02 on 2026-10-05, in the commissioning
Director's session. One dedicated implementer carries the complete plan on a delivery branch cut
from the default branch after this activation is merged. The Director accepts the exact validated
EP-01 candidate before merge and receives the final report. Section 2.4 of the plan records the
included effects and exclusions. The private assignment holds the consumer targets, their order
and holds, and the report-back locator.

Independent planning review on 2026-10-05 found the plan ready after minor fixes. A simulation of
every planned EP-01 edit was eligible for guidance scope against the accepted anchor. All review
corrections are in the plan. Source implementation, review, release and adoption are pending.

## Impact And Validation Boundary

Expected impact: an `instruction-only change` that adds one managed reference under the existing
agent-interface component, two routing edits, declared mirrors, release identity and notes. No
ownership, runtime, schema, CLI or CI behavior changes. The guidance guard, the Go hash and
manifest tests, both native guidance lanes, released-asset smoke and a fresh-agent walkthrough
with a routing probe prove the delivery. Adoption is verified on each assigned default branch.

## EP-01 Start And Authoring Recheck

The delivery branch `claude/claude-code-operating-guide-delivery` was cut without upstream
tracking from fetched remote main at the merged activation
`09d11eaf2d826a326aadacb58e3c7ead55688f01`. The Kit logs record no newer broad acceptance, so the
accepted broad-source anchor remains `c36061f69f18b8a9aff2e018db2d5e35c6188ec8` (OQ-1); changes
since published `v0.1.33` are planning records only. The next unused patch is `v0.1.34` (no such
tag or release exists); published `v0.1.33` is the old-CLI upgrade input. The unchanged baseline
passed the guard, release identity and the Go hash and manifest tests before any edit.

Official sources were rechecked on 2026-10-05 while authoring (OQ-4). Two findings changed the
companion's wording without changing scope:

- Cross-session messaging, task chips and pull request status are now publicly documented, so the
  companion cites them as official. Only the session ID and link tools and the pull request
  binding, status and monitor tools remain labelled as observed app contracts. Holds depend on
  the receiver's inbound setting and on whether each session bypasses permission prompts, not on
  any mode difference; the companion states that rule.
- Claude Code documents a session-scoped `/goal` completion condition, so the plan's "no goal
  object" was inaccurate. The companion says no goal is required and limits `/goal` to explicit
  requests aimed at the agreed finish line and verified as active; it adds neither authority nor
  report-back. This keeps the plan's intent (assignment plus report-back, no watchers). AD-2
  item 8 is corrected to match, for Director acceptance with the candidate.

The observed contracts were checked against this implementer's own exposed tools: the session
tool reports the current desktop session ID, app link, worktree, branch and permission mode, and
a session started from a task chip records its commissioning parent. Peer listing, messaging with
a one-shot idle notice, pull request status and binding, per-pull-request monitor switches and an
archive tool that removes the worktree by default are exposed. A base-branch sync tool is
described for app-made worktrees but not exposed in this session, so the companion makes it
conditional. That sync is described in the desktop app's session instructions to the agent. The
app-made worktree branch had no upstream, and this repository's local exclude already lists the
desktop worktree folder.

## EP-01 Source Candidate

The companion is 120 lines, every line under 100 characters, and covers all twelve AD-2 topics.
The closing pointer of `agent-task-coordination.md` now names both companions; the Codex wording
and the generic contract are otherwise unchanged. The agent-interface README gains one route, and
the new file is declared after the Codex note with exactly `source`, target and
`ownership: managed`. Component, profile, manifest, packaging, CLI, installer and the seven
bootstrap version tokens moved to `0.1.34`. The profile graph digest was taken from the compiled
value. It covers every managed file's bytes, so each later content correction refreshes it and
its mirror. All six changed managed, declaration, profile and manifest files are mirrored
byte-identically.

A repository-local ignored Python environment, installed from the public index without changing
persistent configuration, supplies the test dependencies. Locally, `go test ./internal/hash
./internal/manifest`, release identity, public-core, Markdown, JSON schema and release-manifest
validation pass, and the feedback test selection passes (119 tests).

## EP-01 Review, Walkthrough And Corrections

Candidate `b9a8cc1` went to PR #22; its feedback check passed. Two background subagents were
commissioned: an independent reviewer, and a separate fresh agent that never saw this plan, the
diff or this log. Both stalled once on a stream watchdog (an environment failure) and were resumed
with their context intact.

The fresh agent received only a temporary installation built from the candidate, with a toy
active plan, and the request "How do I hand this plan to an implementer in Claude Code?". It made
23 calls: 12 file reads and 10 read-only shell commands, all inside that installation, and no
tool with side effects. Its route ran from `AGENTS.md` through the managed coordination route to
`agent-task-coordination.md`, then to the plan facts, then through the closing pointer to the
companion. The routing reference, consulted later, confirmed the same route. It answered all
seven scenarios from the installed text: naming, task-chip commission, report-back, safe archival,
neutral memory, session citation and the `CLAUDE.md` case. Its candid gap list mixed toy-setup
artifacts (no remote, branch, exclude or plan metadata) with six genuine clarity gaps. The
implementer separated the two kinds and the reviewer confirmed the split.

The reviewer confirmed the mechanics, routing, public safety, D-2, D-5 and D-6, and the `/goal`
correction. It judged the probe passed and the scenarios correct, and it confirmed the gap triage.
One material finding: "a successful send proves delivery" was wrong. Arrival is not delivery; a
desktop receiver cannot show the approval dialog, so a held message expires and may go
unreported. The walkthrough had repeated the wrong sentence, which confirmed the fix was needed.
Minor findings and the genuine walkthrough gaps were corrected together:

- report-back addressing by locator or task-chip parent, and observed tool names labelled as such;
- the difference between desktop and Claude Code session IDs, with the environment-variable
  reference;
- the implementer's worktree, base and upstream checks;
- archive cleanup and side-session behavior, and auto-archive off for sessions with post-merge
  work;
- automatic memory distinguished from the Kit's agent-memory records, with "neutral" defined;
- where **Project instructions** is set, and the `@AGENTS.md` import also covering
  `CLAUDE.local.md`;
- pull request binding limited to the session's own or user-requested pull requests;
- AD-2 item 4 and the revision note.

The corrected companion had 123 lines; the final one has 124. Plan Section 4.1 already defers
the routing-reference and agent-memory Codex wording. The reviewer judged that a focused recheck
of the corrected sections is enough, because routing is unchanged; no second walkthrough is
needed.

The same reviewer rechecked correction `831c5bc` and found it ready, with five nits. Three were
accuracy refinements to managed text, folded into the final candidate:

- a commission without a desktop session uses its other locator;
- background subagents are not user-owned implementers;
- only messages held by default expire.

The other two were log wording and the PR description. A personal `@AGENTS.md` line in
`CLAUDE.local.md` would also work. It goes beyond the commissioned root-instruction remedy, so
the Director decides. Native guidance run `37289526658` passed both lanes on `831c5bc`, but that
evidence is superseded for the managed bytes the refinements changed.

## EP-01 Validated Candidate

Final candidate `1a8fef8d7ba48fa5ab1a553a5386f416a5f8033a` passed native guidance run
`37290110971`. The eligibility preflight found 47 paths and 11 managed resources, measured
against the accepted anchor. The macOS lane took 59s and the Windows lane 1m29s. Each proved:

- 51 fresh and 51 upgraded managed resources equal candidate source;
- the published `v0.1.33` to `0.1.34` dry-run, failed verification, apply and check, with
  authored preservation;
- reproducible packs, staged installation and the existing historical upgrade checks.

Broad Go, parity and history, Ubuntu semantic and oldest-Git jobs were skipped, as intended. Their
accepted anchor results remain applicable, because no executable, schema, workflow or toolchain
input changed.

The run's candidate assets were retrieved and verified locally along the full chain: catalog,
archive, pack manifest, payload checksums, content identity, binary digest, and a macOS binary
reporting `0.1.34`. The installers, bootstrap and notes inside the packs equal the candidate's
files. A local staged install from those assets succeeded, and an explicit bad checksum was
rejected. Archive SHA-256:

- macOS universal: `ccaa671f1af2b28f1a4b39217be6df06c5adfdbb332cc695a8ecab17ea50cce4`
- Windows x64: `57327fe63328397a34ee4bfdf9cea98832eee823a99a5cfe825360f09a3b9f81`

The same reviewer found the final refinements ready. Director acceptance of this exact candidate
is requested. This log entry and the checklist ticks are planning records only; they change no
release input.

## EP-01 Acceptance And Director Decisions

The commissioning Director accepted candidate `1a8fef8d7ba48fa5ab1a553a5386f416a5f8033a` on
2026-10-05, to be merged through PR head `4ddd89c` (which adds planning records only). Acceptance
includes the corrections to AD-2 items 4 and 8. The Director also decided:

- **Deferred refinement.** A personal `@AGENTS.md` import inside `CLAUDE.local.md` is an
  alternative root-instruction remedy. It waits for the next guide revision, with no
  re-validation cycle now.
- **Owner merge timing.** For one assigned target, merging to the default branch triggers its
  owner's development deployment. The implementer opens a Kit-only pull request and the owner
  decides when to merge it.
- **Held target.** One target stays local until the Director gives the go-ahead.

## EP-02 Release, Smoke And Adoption

PR #22 merged normally at 2026-10-05T09:33:02Z as `6eb5f1a7152a744d85387a24d859943afbcdd55e`.
The merged tree equals the accepted head. It differs from the validated candidate only in two
planning records, and cumulative guidance eligibility still holds. Annotated tag `v0.1.34`
points to the merge.

The release publishes 14 assets: both verified packs, the external catalog, the installers,
bootstrap and notes, each with a SHA-256 sidecar. The unchanged unsigned internal/prototype
HTTPS-plus-SHA-256 boundary applies. The catalog changes only download URLs from the candidate
catalog; published catalog SHA-256 is
`ef85adc0a578e01bdb93c937181aeaf6442833614b289a1d3d94747b5db75903`. Re-downloaded public assets
are byte-identical to the staged ones. A public install through the published installer reports
`0.1.34`. Released-asset smoke run `37290917668` passed both public native jobs without replaying
source suites.

Each of the four assigned default branches was upgraded in its own isolated worktree, cut from
the current remote default branch. The upgrade used a verified matching published CLI installed
in a temporary location; the shared CLI was not changed. Each ran `check`, a dry-run, apply, a
completed-state wait and a final `check`. The dry-run's catalog-to-archive, pack-to-binary and
staged-version validations passed each time. All four reached `0.1.34` with 51 managed resources
byte-equal to source. Only Kit-managed files and the lock changed, so authored configuration,
plans, module state, instructions and local-user files were preserved. One target upgraded
directly from `v0.1.32`, which confirms A-1; its adoption includes the `v0.1.33` guidance. Using
matching CLIs avoided the partial classification in OQ-2.

Delivery state by target:

- **Two targets merged.** Their pull requests merged normally with no required checks
  triggered. Each fetched remote default branch equals the reviewed adoption head, reports
  `0.1.34`, carries the new route byte-identical to source and checks healthy.
- **One target, pull request open.** Its Kit-only pull request is open and its required source
  check passed. Its owner merges, ideally together with the next planned push so that only one
  deployment runs. The owner has the link and the deployment note.
- **One target, held.** Its adoption is prepared, verified and committed locally, and not
  pushed. Its governance gate's classification of the new Kit-managed path stays unverified
  until the Director's go-ahead.

The plan stays active until both remaining targets are adopted on their default branches. The
Director closes those with their owners. Final completion of the plan, log and index follows
that adoption.
