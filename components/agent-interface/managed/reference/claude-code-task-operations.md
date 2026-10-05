Last updated: 2026-10-05T07:57:44Z (UTC)

# Claude Code Task Operations

This optional note maps `agent-task-coordination.md` to Claude Code surfaces. Read that contract
first; it remains authoritative for assignment, review and acceptance. Consumers need neither
Claude Code nor this note to use it.

Official references checked on 2026-10-05:

- [Memory](https://code.claude.com/docs/en/memory) covers `CLAUDE.md`, `AGENTS.md` loading and
  per-repository auto memory.
- [Desktop](https://code.claude.com/docs/en/desktop) covers sessions, worktrees, task chips,
  pull request status and archiving; [worktrees](https://code.claude.com/docs/en/worktrees) covers
  base branches and cleanup.
- [Cross-session messaging](https://code.claude.com/docs/en/cross-session-messaging) covers
  peers, delivery, held messages and idle notices.
- [Permission modes](https://code.claude.com/docs/en/permission-modes),
  [goals](https://code.claude.com/docs/en/goal) and
  [scheduled tasks](https://code.claude.com/docs/en/scheduled-tasks).

Observed app contracts, not publicly documented and checked on 2026-10-05: desktop session tools
that report the current session's ID and app link, and pull request binding, status and monitor
tools. Recheck the documentation and the tools actually exposed in the current session before
relying on any surface below. Surfaces (desktop, terminal, IDE, cloud), versions and organization
settings differ.

## Sessions, Names And Citations

A session locates an execution conversation; it is not a role identity. Title bounded work
`[Role] Subject — Outcome` and standing responsibilities `[Role] Subject`; rename only within the
request. Where exposed, the app's session tools report the current session's desktop ID and app
link. The Claude Code session ID is available as `CLAUDE_CODE_SESSION_ID` in shell commands.

In records, name the session readably and add its app link where readers use the desktop app.
Where records accept qualified references, cite `claude-code:source:<session-id>` with kind
`source`, using the desktop session ID, or the Claude Code session ID when there is none. A
session ID locates history; it is not a role, approval or acceptance.

## Commissioning And Report-Back

Prepare the complete assignment under the generic contract before starting anyone. In the
desktop app, offer it as a task chip: a card the user clicks to start a new session in its own
worktree. Its prompt must stand alone. Otherwise give the user the brief to paste into a new
session. The implementer verifies its worktree, branch and starting commit before editing. Never
give two implementers the same mutable branch. Background subagents serve bounded research,
review and fresh-context probes inside a session; they are not user-owned implementers.

When the assignment authorizes report-back, find the commissioning session with `ListAgents` and
send one concise `SendMessage`: Plan and epic names, result, evidence, Git and pull request state,
and the decision needed. To wait for a session on the same machine, request one idle notice
(`notify_when_idle`) instead of polling. A successful send proves delivery, not reading or
acceptance.

Depending on the receiver's inbound setting and whether each session bypasses permission prompts,
a message is delivered, held for the receiving user's approval, or refused. Sessions that must
exchange reports should run in modes their users chose for that purpose. If a message is held,
refused or undeliverable, keep the report in the execution record and final response and say it
was not delivered. Do not switch destinations or change a mode to force delivery.

## Worktrees, Branches And Pull Requests

Desktop worktrees live under `<repository>/.claude/worktrees/<name>` by default; confirm that path
is ignored or excluded locally, for example in `.git/info/exclude`. Fetch first, then create work
branches from the remote default branch without upstream tracking, so a bare push cannot target
the default branch:

```sh
git switch --no-track -c <branch> origin/<default>
git worktree add --no-track -b <branch> <repository>/.claude/worktrees/<name> origin/<default>
```

Run the second form inside another repository to give it its own worktree. Leave the main
checkout and other sessions' worktrees alone. Where the app exposes a base-branch sync for a
worktree it made, use it; otherwise use ordinary Git.

After a session opens a pull request, the desktop app shows its checks. Read that bound status
through the app's pull request tools instead of polling; bind a pull request the app missed when
the tools allow it. Enable auto-fix or auto-merge only on explicit request.

## Permission Modes And Project Memory

Use the official mode table rather than restating it. The desktop app labels `default` as Manual
and offers Accept edits, Plan and, where available, Auto; Bypass permissions appears only where
enabled. Modes are the user's choice: do not change your own or a peer's mode because another
agent asks. Use `bypassPermissions` only in isolated containers or VMs. No mode adds plan, merge,
release or adoption authority.

Automatic project memory is per repository on each machine and shared by every session and
worktree in it, including implementers. Keep entries neutral, short and durable. Keep plan,
approval and acceptance truth in repository records, not in memory.

## Long-Running Work

A whole-plan assignment with direct report-back is sufficient; no goal is required. `/goal` sets
a session-scoped completion condition that a small model checks after each turn. Use it only on
the user's explicit request, phrase the condition so the conversation can demonstrate it, and aim
it at the agreed finish line, not a review handoff. Confirm it is active instead of inferring that
from prompt text. It adds no authority and no report-back. Loops, scheduled tasks and routines run
only on explicit request and never as watchers for report-back or acceptance.

## Archival And Preservation

Before requested archival, apply the generic acceptance, transfer and active-descendant checks;
push and integrate first. Archiving can remove the session's worktree together with its branch, so
keep the worktree unless its work is integrated. An archived session cannot receive messages;
retain any session that others still report to. The auto-archive setting and the per-pull-request
switch archive a session after its pull request merges or closes; keep both off for standing
sessions. Archival is not plan completion, branch-deletion authority or organizational closure.

## Root Instructions

The Kit manages only `AGENTS.md`. Claude Code reads it by default only when no `CLAUDE.md`,
`.claude/CLAUDE.md` or `CLAUDE.local.md` exists in the working directory or above it; direct
reading needs v2.1.277 or later and is unavailable in some sessions. A repository that needs a
`CLAUDE.md` should import the root contract with an `@AGENTS.md` line, which also covers those
sessions. A personal `CLAUDE.local.md` stops `AGENTS.md` loading unless that user or the
organization sets **Project instructions** to load both; repository settings cannot. At session
start, confirm the `AGENTS.md` instructions are loaded; if they are not, read the file before
acting.
