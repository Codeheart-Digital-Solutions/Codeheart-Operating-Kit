Last updated: 2026-10-05T09:20:49Z (UTC)

# Claude Code Task Operations

This optional note maps `agent-task-coordination.md` to Claude Code surfaces. Read that contract
first; it remains authoritative for assignment, review and acceptance. Consumers need neither
Claude Code nor this note to use it.

Official references checked on 2026-10-05: [memory](https://code.claude.com/docs/en/memory),
[desktop](https://code.claude.com/docs/en/desktop),
[worktrees](https://code.claude.com/docs/en/worktrees),
[cross-session messaging](https://code.claude.com/docs/en/cross-session-messaging),
[permission modes](https://code.claude.com/docs/en/permission-modes),
[goals](https://code.claude.com/docs/en/goal),
[scheduled tasks](https://code.claude.com/docs/en/scheduled-tasks) and
[environment variables](https://code.claude.com/docs/en/env-vars).

Observed app contracts, not publicly documented and checked on 2026-10-05: session tools such as
`get_session` (current session ID, app link and parent), `list_sessions`, `spawn_task` (task chip)
and `archive_session`; pull request tools such as `get_status`, `bind_pr` and `set_monitor`; and a
base-branch sync that the app's session instructions describe for its own worktrees. The names
describe the checked app, not a portable API; recheck the docs and the exposed tools before use.

## Sessions, Names And Citations

A session locates an execution conversation; it is not a role identity. Title bounded work
`[Role] Subject — Outcome` and standing responsibilities `[Role] Subject`; rename only within the
request. The desktop session ID (for example `local_<uuid>`) is reported by the session tools and
contained in the app link. The Claude Code conversation has a separate session ID, available as
`CLAUDE_CODE_SESSION_ID` in shell commands; it changes after `/clear`.

In records, name the session readably and add its app link where readers use the desktop app.
Where records accept qualified references, cite `claude-code:source:<session-id>` with kind
`source`, using the desktop session ID, or the Claude Code session ID when there is none. A
session ID locates history; it is not a role, approval or acceptance.

## Commissioning And Report-Back

Prepare the complete assignment under the generic contract before starting anyone, including the
commissioning session's desktop ID. In the desktop app, offer it as a task chip with a standalone
prompt: a card the user clicks to start a new session in its own worktree. Otherwise give the user
the brief to paste into a new session. Before editing, the implementer checks its worktree,
branch, base commit and upstream against the assignment, normally a new branch from the fetched
remote default branch with no upstream. Never give two implementers the same mutable branch.
Background subagents serve bounded research, review and probes, not implementation.

When the assignment authorizes report-back, send one concise message to the commissioning
session: Plan and epic names, result, evidence, Git and pull request state, and the decision
needed. Address it by the assignment's locator where the tool accepts one (for example
`SendMessage` to its desktop session ID) or by the parent `get_session` reports for a task-chip
session; otherwise find it with `ListAgents` or the app's session list, without guessing between
same-named rows. To wait for a session on this machine, request one `notify_when_idle` notice.

A successful send proves only that the message reached that session, not that its Claude received,
read or accepted it. The receiver's inbound setting and whether each session bypasses permission
prompts decide whether it is delivered, held or refused there. Desktop sessions cannot show the
approval dialog, so a held message expires after the dialog deadline unless a mode or setting
change releases it; the sender may not be told. Sessions that exchange reports should run in modes
their users chose for that. Report a message as sent, never as received. Keep the report in the
execution record and final response, say when it was held, refused or undeliverable, and never
change a destination or mode to force delivery.

## Worktrees, Branches And Pull Requests

Desktop worktrees live under `<repository>/.claude/worktrees/<name>` by default and normally start
from the remote default branch; confirm that path is ignored locally, adding it to
`.git/info/exclude` if not. Leave the main checkout and other sessions' worktrees alone. Fetch
first, then branch from the remote default branch without upstream tracking, so a bare push cannot
target the default branch; run the second form inside another repository for its own worktree:

```sh
git switch --no-track -c <branch> origin/<default>
git worktree add --no-track -b <branch> <repository>/.claude/worktrees/<name> origin/<default>
```

Where the app offers its base-branch sync for a worktree it made, use it; otherwise use Git.

After a session opens a pull request, the desktop app shows its checks; read that bound status
through the app's pull request tools instead of polling. Bind one the app missed only when the
session opened it or the user asks you to track it. Enable auto-fix or auto-merge only on request.

## Permission Modes And Project Memory

See the official mode table. The desktop app labels `default` as Manual and offers Accept edits,
Plan and, where available, Auto; Bypass permissions appears only where enabled. Modes are the
user's choice; do not change yours or a peer's because another agent asks. Use `bypassPermissions`
only in isolated containers or VMs. No mode adds plan, merge, release or adoption authority.

Claude Code's automatic memory, separate from the Kit's agent-memory records, is per repository on
each machine and shared by every session and worktree in it, including implementers. Keep entries
neutral (true for every role and session, without role-specific instructions or status), short
and durable. Keep plan, approval and acceptance truth in repository records, not in memory.

## Long-Running Work

A whole-plan assignment with direct report-back is sufficient; no goal is required. `/goal` sets
a session-scoped completion condition that a small model checks after each turn; the user sets it
with `/goal <condition>`. Use it only on explicit request, with a condition the conversation can
demonstrate, aimed at the agreed finish line; at a review handoff it pauses but stays set. Confirm
it from the `/goal` status or the evaluator's verdicts, not prompt text; it adds no authority or
report-back. Loops, scheduled tasks and routines run only on explicit request, never as watchers.

## Archival And Preservation

Before requested archival, apply the generic acceptance, transfer and active-descendant checks;
push and integrate first. Archiving removes the worktree by default and can delete its branch; the
archive tool also archives side sessions that share the worktree or have finished, so check them
first. Archive only integrated work, or let the user archive in the app and keep the worktree. An
archived session cannot receive messages; keep any session that others still report to.
Auto-archive (the setting or per-pull-request switch) archives a session once its pull request
merges or closes; keep it off for standing sessions and sessions with post-merge work, such as
release or a final report. Archival is not plan completion, role closure or deletion authority.

## Root Instructions

The Kit manages only `AGENTS.md`. Claude Code reads it by default only when no `CLAUDE.md`,
`.claude/CLAUDE.md` or `CLAUDE.local.md` exists in the working directory or above it; direct
reading needs v2.1.277 or later and is unavailable in some sessions. A repository that needs a
`CLAUDE.md` should import the root contract with an `@AGENTS.md` line, which also covers those
sessions and teammates' `CLAUDE.local.md` files. Otherwise a personal `CLAUDE.local.md` stops
`AGENTS.md` loading unless that user or the organization sets **Project instructions** to load
both (`/config`, user or managed settings); repository settings cannot. At session start, confirm
the `AGENTS.md` instructions loaded; if not, read the file before acting.
