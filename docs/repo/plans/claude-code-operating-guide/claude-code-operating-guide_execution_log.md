Last updated: 2026-10-05T09:20:49Z (UTC)
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
neutral memory, session citation and the `CLAUDE.md` case. Its candid gap list separated toy-setup
artifacts (no remote, branch, exclude or plan metadata) from six genuine clarity gaps.

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

The companion is now 123 lines. Plan Section 4.1 already defers routing-reference and agent-memory
Codex wording. The reviewer judged that a focused recheck of the corrected sections is enough,
because routing is unchanged; no second walkthrough is needed.
