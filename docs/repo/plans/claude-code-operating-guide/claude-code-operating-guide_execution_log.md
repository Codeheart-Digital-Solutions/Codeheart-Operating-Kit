Last updated: 2026-10-05T08:01:28Z (UTC)
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
conditional. The app-made worktree branch had no upstream, and this repository's local exclude
already lists the desktop worktree folder.

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
