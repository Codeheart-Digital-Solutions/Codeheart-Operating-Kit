Last updated: 2026-10-10T16:27:45Z (UTC)

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
commissioning session's locator (its desktop session ID where it has one). In the desktop app,
offer it as a task chip with a standalone prompt: a card the user clicks to start a new session in
its own worktree. Otherwise give the user the brief to paste into a new session. Before editing,
the implementer checks its worktree, branch, base commit and upstream against the assignment,
normally a new branch from the fetched remote default branch with no upstream. Never give two
implementers the same mutable branch. Background subagents serve bounded research, review and
probes, not user-owned implementation.

When the assignment authorizes report-back, send one concise message to the commissioning
session: Plan and epic names, result, evidence, Git and pull request state, and the decision
needed. Address it by the assignment's locator where the tool accepts one (for example
`SendMessage` to its desktop session ID) or by the parent `get_session` reports for a task-chip
session; otherwise find it with `ListAgents` or the app's session list, without guessing between
same-named rows. To wait for a session on this machine, request one `notify_when_idle` notice.

A successful send proves only that the message reached that session, not that its Claude received,
read or accepted it. The receiver's inbound setting and whether each session bypasses permission
prompts decide whether it is delivered, held or refused there. Desktop sessions cannot show the
approval dialog, so a message held by default expires after the dialog deadline unless a mode or
setting change releases it; the sender may not be told. Sessions that exchange reports should run
in modes their users chose for that. Report a message as sent, never as received. Keep the report
in the execution record and final response, say when it was held, refused or undeliverable, and
never change a destination or mode to force delivery.

## Worktrees, Branches And Pull Requests

Apply "Assignment Workspaces" in `agent-task-coordination.md` for selection, closure ownership
and preservation. Standing coordination normally uses the usual checkout. App-managed locations
and manual workspace preferences are distinct; do not relocate an app-managed workspace to match
a manual folder convention.

Desktop worktrees live under `<repository>/.claude/worktrees/<name>` by default and normally start
from the remote default branch; confirm that path is ignored locally, adding it to
`.git/info/exclude` if not. Leave other sessions' worktrees alone; use the usual checkout only as
"Assignment Workspaces" permits. Fetch first, then branch from the remote default branch without
upstream tracking, so a bare push cannot target the default branch. Use the first form only inside
your own app-created assignment worktree when a new branch is needed. For a manually created
worktree, resolve `<workspace-path>` under the chosen manual workspace root and run the second
form from the owning repository:

```sh
git switch --no-track -c <branch> origin/<default>
git worktree add --no-track -b <branch> <workspace-path> origin/<default>
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

Executing an approved implementation plan uses `/goal` by default, under the generic rule in
`agent-task-coordination.md`; reviews, discovery and routine changes do not. `/goal <condition>`
sets a session-scoped completion condition that a small model checks after each turn. Commission
it with the user's actual authority, with a condition the conversation can demonstrate, aimed at
a real handoff such as source review or the agreed finish line. Confirm activation from the
`/goal` confirmation, status or evaluator verdicts, not from prompt text. A goal ending at a
handoff does not complete the plan; continue the next phase in the same session with a new goal
after delegated acceptance. Honor an explicit opt-out, and if a goal cannot be activated or
verified, say so and continue only authorized ordinary work. It adds no authority or report-back.
Loops, scheduled tasks and routines run only on explicit request, never as watchers.

## CLI-Owned Sessions For Cross-Tool Work

When a Codex coordinator commissions Claude, follow `../runbooks/coordinate-cross-tool-task.md`.
The installed helper starts or resumes one CLI-owned session in print mode with a fixed flag set,
the brief on stdin and stream-json output kept in the attempt directory. One driver per session:
do not resume a CLI-owned session from the desktop app or a second terminal while an attempt is
running, and do not drive app-owned sessions through the helper.

A brief that begins with `/goal <condition>` sets the goal in print mode; the coordinator reads the
confirmation in the original output as activation evidence. Each phase is a new attempt that
resumes the retained session with its own goal condition.

Permissions come from the request: an explicit mode (never `bypassPermissions`), tool lists and an
optional authorized settings profile. Headless runs cannot show approval dialogs; denials appear
in the result's permission-denial metadata even when the process exits successfully, and the
helper returns them unchanged. A profile that allows ordinary pull-request merges in the task
repository grants tool permission only; the workflow still checks authority, review, CI and
candidate identity, and no profile adds release or deployment authority.

For delivery work, an explicit `dontAsk` profile listing the task's ordinary development and
GitHub commands is simpler than an Auto-mode classifier allowance, which refused a self-authored
merge and a public release. `dontAsk` denies actions that would otherwise prompt; allowed actions
still run, and a denial returns to the agent as a refusal. Choosing or changing a profile is the
policy owner's explicit, recorded decision before work resumes; an agent never switches profiles,
modes or tools to get past its own denial, and relayed approval does not override a refusal.
Confirm a new profile once with harmless allowed, ask and unlisted probes. Start or resume the
session in the repository it changes, because changing directory and then running Git in one
command is evaluated separately. See the permission profiles section of
`cross-tool-coordination-contract.md`.

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
