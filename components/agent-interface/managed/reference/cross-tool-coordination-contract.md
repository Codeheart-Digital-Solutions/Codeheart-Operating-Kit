Last updated: 2026-10-10T07:07:56Z (UTC)

# Cross-Tool Coordination Contract

This reference defines the request, records and brief shapes used by
`../runbooks/coordinate-cross-tool-task.md`. The generic assignment, review and acceptance rules
remain in `agent-task-coordination.md`. Nothing here is needed for ordinary current-tool work.

## Roles

| Function | Owns | Must not |
| --- | --- | --- |
| Coordinator | Complete brief and request, defaults from approved local choices, attempt identity, primary review of delivered work and original replies, answers to substantive questions, acceptance within its mandate. | Delegate those choices to the relay, forward a reviewer's verdict as its own review, or duplicate the commissioned investigation while waiting. |
| Substantive agent (advisor, implementer, independent reviewer) | The commissioned question or work, its own questions, blockers and final reply. | Treat relayed text as new human authority or widen its assignment. |
| Relay or temporary transport worker | Passing the prepared request path to the helper, waiting on that process, loading the message file into the native send tool and recording the send result. | Research, summarize, interpret, answer, accept, choose defaults, retry or start follow-on work. |

A separate independent reviewer is commissioned only when complexity, risk, the coordinator's
own implementation involvement or a binding gate warrants it; the coordinator still reviews the
work itself and decides acceptance. Independence comes from contribution and assessment, not
from session freshness. Reuse an
advisor or implementer session while it remains useful; a relay is replaceable without replacing
the substantive session. Session IDs locate execution; they are not role identities.

## Invocation Request

The coordinator writes one JSON request per attempt. The helper rejects unknown fields, missing
values, relative paths, unsupported modes and trailing data before creating any attempt state.

```json
{
  "schema_version": 1,
  "assignment_id": "example-plan-implementation",
  "attempt_id": "source-phase-1",
  "return_destination": {"thread_id": "<coordinator chat id>", "host_id": "<host id>"},
  "working_directory": "/absolute/path/to/worktree",
  "brief_file": "/absolute/path/to/ignored/brief.md",
  "executable": "/absolute/path/to/claude",
  "model": "<approved model>",
  "effort": "high",
  "session": {"mode": "new"},
  "permissions": {
    "mode": "default",
    "permission_prompts": "none",
    "tools": ["Read", "Glob", "Grep", "Edit", "Write", "Bash"],
    "allowed_tools": ["Read", "Glob", "Grep"],
    "disallowed_tools": [],
    "settings_file": "/absolute/path/to/ignored/profile.json",
    "add_dirs": [],
    "disable_mcp_servers": true,
    "disable_chrome": true
  },
  "authority_refs": ["docs/repo/plans/<plan>_implementation_doc.md", "<assignment record locator>"],
  "state_root": "/absolute/path/to/repo/.codeheart/local/coordination"
}
```

- `session.mode` is `new` (the helper allocates and records a UUID before launch; omit `id`) or
  `resume` (supply the retained CLI-owned session UUID).
- `permissions.mode` is `default` (named `manual` in newer CLIs; both are passed through),
  `acceptEdits`, `auto`, `plan` or `dontAsk`.
  `bypassPermissions` is refused. Optional fields are omitted when not approved; the helper never
  fills them in. `effort` accepts `low`, `medium`, `high`, `xhigh` or `max`.
- `permission_prompts` is passed only when set. A CLI version without that option rejects it and
  the attempt returns `unsupported_invocation`; the helper does not retry without it.
- Use one shared ignored `state_root` on a host for every participating launch, so the session
  lock can refuse overlapping writers.
- `executable` and `settings_file` come from approved local settings, never from a `PATH` lookup
  the coordinator did not approve. The executable is run directly; Windows `.cmd` or `.bat`
  wrappers are not supported.

## Fixed Invocation

The helper runs the executable directly, without a shell, in `working_directory`:

```text
<executable> --print --output-format stream-json --verbose --model <model>
  (--session-id <new uuid> | --resume <retained uuid>) --permission-mode <mode>
  [--permission-prompts <value>] [--effort <level>] [--settings <file>]
  [--strict-mcp-config --mcp-config {"mcpServers":{}}] [--no-chrome]
  [--tools <a,b>] [--allowedTools <rule>...] [--disallowedTools <rule>...] [--add-dir <dir>...]
```

The brief file is connected to stdin unchanged. Child stdout and stderr are connected directly to
files in the attempt directory, so bytes already written survive helper or worker death. The CLI
version is recorded from the CLI's own init record; the helper does not scrape help output or keep
a compatibility matrix. Check official CLI documentation when maintaining this recipe.

## Attempt Directory

The coordinator knows the directory before dispatch:
`<state_root>/attempts/<assignment_id>/<attempt_id>/`. Session locks live at
`<state_root>/sessions/<session_id>.lock`. All files are private, ignored machine-local evidence.

| File | Content |
| --- | --- |
| `attempt.json` | Compact record: status, request path and digest, requested and actual session, model, CLI version, argv, launcher and child identity, exit, response locator, `earlier_results`, denial metadata, usage, lock release. Usage, cost, duration and turn count come from the last result record only; they are not summed across records. |
| `stdout.jsonl` | Original CLI output. The final reply is the `result` field of the last `type: result` line and must be a JSON string (an empty string counts); a missing, null or non-string value is `incomplete_output`, never a captured reply. A run can emit more than one result record, for example a substantive reply or a failure followed by a background-task completion. Each earlier record is listed in `earlier_results` with its line, `subtype`, `is_error` and, when it has a text reply, a `reply` locator. |
| `stderr.txt` | Original CLI error output. |
| `message.json` | Native send-message arguments: `threadId`, `hostId` and `prompt`. `message_file` is the path to this file, not a send-tool argument: read and parse it, then pass these three values unchanged as the tool's arguments. Use it only when this invocation's output names it as `message_file`. |
| `delivery.json` | `pending`, then `sent`, `rejected` or `uncertain`, with the original receipt or error. |

Statuses: `response_captured` (a final reply was captured; not task acceptance), `cli_error`,
`session_mismatch`, `incomplete_output`, `unsupported_invocation`, `execution_failed`,
`interrupted` (child ended by a signal), `helper_interrupted` (helper stopped while the child may
still run; lock retained), `launch_failed`, `session_locked` (no second writer launched) and
`attempt_exists` (replay refused; the original attempt is untouched). Permission denials are
reported separately with their count and tool names; the CLI can exit successfully after a denial.
Denials are gathered from every result record of the invocation, so a later record without
denials cannot hide an earlier one; a denial repeated with the same tool use ID counts once.

Helper output (stdout) is compact JSON: `status`, `detail`, identifiers, `attempt_dir`,
`attempt_file`, and `message_file` and `delivery_file` only when this invocation wrote them. Every
outcome after the attempt directory exists writes all three records, including `session_locked`,
`launch_failed` and `helper_interrupted`; the lock is released only when this launcher owns it and
no child can still run. `fallback_message` holds send arguments for the request's own return
destination when no message file could be written: `invalid_request` (only when the destination
itself is valid), `helper_error` before an attempt record exists, a failed message write, and
`attempt_exists`. A refused replay changes nothing in the earlier attempt and its notice says so.
`persistence_errors` lists records that could not be written. Exit status is 0 only for
`response_captured` with no persistence error, 2 for `invalid_request` and 1 otherwise.

The message prompt states the assignment, attempt, status, session, denial count and attempt
record path. It carries the original reply unchanged when it has at most 2,000 characters, and
otherwise an exact file, line and JSON-field reference. When earlier result records exist, it lists
each one's line, `subtype`, `is_error` and whether it has a text reply, without quoting it; the
coordinator reads every one. An earlier record marked as an error or not `success` is also named
in `detail`. The status can still be `response_captured` when the last record holds a reply, but
the earlier failure stays explicit. The message contains no summary or judgment.

Helpers released before this capture change, including v0.1.35, keep only the last result
record, so an earlier reply, failure or denial can be missing from `attempt.json` and the
message. With those helpers, read `stdout.jsonl` directly whenever it holds more than one
`type: result` line. The release that ships this change must say so in its notes.

## Delivery And Recovery Commands

```text
codeheart-operating-kit coordination record-delivery --attempt-dir <dir>
  --status {sent,rejected,uncertain} [--receipt-file <original receipt or error>]
  [--reported-thread-id <id>] [--reported-host-id <id>]
codeheart-operating-kit coordination release-lock --state-root <root> --session-id <uuid>
  --assignment-id <id> --attempt-id <id> [--manual-verification "<statement>"]
```

Run both commands, like `invoke-claude`, with the exact helper path the coordinator verified
and supplied in the transport request; a `PATH` lookup can resolve an older helper without the
`coordination` subcommand.

`record-delivery` changes a pending record once. A reported recipient that differs from the
prepared one is retained as `uncertain`; an omitted reported recipient never counts as
confirmation. `sent` means the tool accepted the send, not that the coordinator read or accepted it.

`release-lock` removes only the named attempt's lock, after checking the lock owner, the host and
that the recorded launcher and CLI processes have ended. A reused process ID counts as running,
which errs toward retaining the lock. When the CLI process was never recorded, it refuses unless
the coordinator supplies a manual verification statement, which is kept in `attempt.json`. A
statement never overrides a live or unknown process or an empty or unreadable lock; those
cases follow the owner-directed recovery in the runbook.

## Brief Shapes

A brief is ordinary prompt text written by the coordinator. It names the role, outcome, sources,
settled decisions, authority and limits, review owner, finish line and return behavior. Attribute
authority to the owning record or actual user instruction; transport in the user-message role is
not new human authorship.

Whole-plan implementer, first phase (native goal by default; see `agent-task-coordination.md`):

```text
/goal Implement <phase> of <Plan name> at <canonical path> on <branch>. Complete its code,
documentation and validation, commit, push and open/update the delivery PR. Return the
source-review handoff with head, tests, deviations and remaining work, or the original blocker.
Stop before <reserved effects>; <Director role> accepts first. This ends the phase goal, not the plan.

Context: <assignment record>, <plan>, <discovery>, <runbooks>. Authority: <grant and its limits>.
Reply normally in this conversation; do not write a separate response file.
```

Advisor or reviewer consultation (no goal):

```text
You are the <independent reviewer | program advisor> for <subject>. Read <sources>. Question:
<question>. Settled decisions: <list>. Do not edit files. Reply with findings, severity, evidence
and your verdict in this conversation.
```

Continuation after a question (resume the same session):

```text
Answer from <Coordinator role>, based on <owning decision or user instruction>: <answer>.
Continue the same assignment within the original authority.
```

For a later implementation phase after director acceptance, resume the same session with a new
`/goal` condition for that phase. A goal ending at a checkpoint never completes the whole plan.

## Permission Profiles

A profile is a separate, explicitly authorized settings file passed with the request; the brief
cannot authorize itself. Nothing here is installed globally, and a consumer that never selects
cross-tool delegation needs no profile. Choosing and approving one is the consumer's decision;
its exact command list and paths stay in the consumer's ignored local layer.

Choosing or changing a profile is the policy owner's explicit, recorded decision, made before
work resumes. An agent never switches profiles, modes or tools in response to its own denial.
Once approved, a profile covers the actions it lists; the owner does not need to re-approve
each covered action.

### Explicit command permissions (`dontAsk`)

For delivery work, an explicit profile is simpler and more predictable than an AI classifier.
`dontAsk` denies any action that would otherwise ask for permission. Actions the rules allow, and
actions that never need permission, still run. A denial comes back to the agent as a refusal; the
session does not end, so the agent reports it instead of being judged case by case (see the
official permissions documentation). List only the task's ordinary development and delivery
commands, keep explicit ask and deny rules for risky options, and grant file edits only inside
the task's directories:

```json
{
  "permissions": {
    "defaultMode": "dontAsk",
    "allow": [
      "Read", "Glob", "Grep",
      "Edit(//<task worktree>/**)", "Write(//<task worktree>/**)",
      "Bash(git status *)", "Bash(git diff *)", "Bash(git log *)", "Bash(git fetch *)",
      "Bash(git add *)", "Bash(git commit *)", "Bash(git push *)", "Bash(git switch *)",
      "Bash(gh pr view *)", "Bash(gh pr create *)", "Bash(gh pr edit *)", "Bash(gh pr checks *)",
      "Bash(go test *)", "Bash(<repository test command> *)",
      "Bash(printf <allowed readiness probe>)"
    ],
    "ask": ["Bash(printf <ask readiness probe>)"],
    "deny": ["Bash(git push *--force*)", "Bash(git push -f*)", "Bash(git push * -f*)", "Bash(git push *--delete*)", "Bash(git branch -D *)", "Bash(gh pr merge *--admin*)"]
  }
}
```

Add release, merge or other public-surface commands only when the assignment includes those
effects. Command patterns are workflow control for a trusted agent, not operating-system
isolation: an allowed test or build command runs repository code, and `gh` and `git push` rules
do not check repository identity, review readiness or CI. The deny entries are illustrative
examples, not exhaustive enforcement: other spellings, such as a `+<refspec>` force push or
`:<branch>` deletion, are not matched. Branch protection, the reviewed scope and the delivery
gates stay authoritative. The workflow still checks authority, review, CI and the exact target.

When setting up or changing a profile, confirm it once with harmless probes before real work:
the allowed probe runs, and the ask probe and one unlisted command that would normally prompt
(for example an interpreter one-liner) are refused. Any other refusal during real work is a
blocker to report unchanged.

Launch or resume the session with its working directory set to the repository it must change,
and add other task directories as working directories. A command that changes directory and then
runs Git is evaluated separately, because Git can run hooks from the new directory, and was
refused under such a profile even with `cd` allowed. Use one session start per repository
instead. In one run, listed file commands on staging paths outside the working directories were
refused and the same kind of work succeeded once those paths were working directories; keep task
files inside the session's working directories.

Environment failures are not permission problems. Example: `shasum` on macOS aborted under an
inherited `C.UTF-8` locale; setting `LC_ALL`, `LANG` and `LC_CTYPE` to a supported locale such as
`C` for that invocation fixed it. Diagnose the actual platform error rather than adding
permissions or switching tools.

### Auto-mode allowance

Auto mode keeps the default protections and lets an AI classifier decide each action. An
`autoMode.allow` entry can describe an allowance, for example ordinary pull-request merges in
the task repository:

```json
{
  "permissions": {"ask": ["<explicit ask rule>"], "deny": ["<explicit deny rule>"]},
  "autoMode": {
    "allow": [
      "$defaults",
      "Ordinary pull-request merges in the current task repository are allowed. The task workflow owns authorization, review, validation and readiness. This grants tool permission, not an instruction to merge. It does not permit bypassing required reviews or checks, --admin or --force, direct default-branch pushes, releases, deployments or unrelated operations. All other default protections and explicit ask/deny rules remain."
    ]
  }
}
```

In practice the classifier refused a self-authored merge and later a public tag and release.
Approval relayed in a brief is not new human authority and does not override a tool refusal.
For delivery work that includes merge or release effects, the policy owner may prefer an
explicit profile, chosen as a recorded decision before work resumes; agents do not make that
switch themselves.

Neither example is an installed setting. Applying one requires authority covering its effects.
A real denial returns unchanged; do not route the same action through another agent, tool or
command spelling to evade it.
