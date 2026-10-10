Last updated: 2026-10-10T08:30:00Z (UTC)

# Coordinate Cross-Tool Task

Use this runbook only when a user or assignment selects another tool's agent, currently a Claude
CLI advisor, implementer or independent reviewer commissioned from a Codex coordinator. Ordinary
planning, implementation and delegation continue with the current tool and need nothing here.

Audience: hybrid

Intent:
Let an accountable coordinator commission a Claude CLI session through an independent relay,
remain available to the user, and receive the original question or reply after yielding. Guide
optional first-use setup without making Claude, an account or a relay a Kit prerequisite.

Success:
The selected agent receives a complete brief through the installed helper, its original reply or
blocker reaches the exact coordinator chat or is retained at a known attempt path, and the
coordinator assesses it and continues the same session when needed.

Agent judgment boundary:
The coordinator chooses ordinary brief wording, attempt identifiers and timing within existing
authority. It must not invent preferences, models, permissions, recipients or chat authority,
substitute a provider silently, copy credentials, or let the relay make substantive choices.

Stop boundary:
Stop and return to the user or responsible owner when the selected tool is not installed or not
signed in, a relay or return route is missing, stale or unauthorized, a permission is denied, a
required gate fails, or the requested change exceeds the assignment's authority.

## Routing

| Field | Value |
| --- | --- |
| Intent family | Cross-tool commissioning, continuation and result return. |
| Owner | Managed Operating Kit procedure; consumer operational reference for local choices. |
| Scope | One named assignment with its own attempts; never a standing orchestration service. |
| Authority source | The owning plan or assignment record and the user's actual instructions. |
| Execution surface | `codeheart-operating-kit coordination` helper plus the host's native chat tools. |
| Preconditions | Explicit cross-tool choice, ready CLI and sign-in on this machine, authorized relay and return route. |
| Evidence | Ignored attempt directory, plus durable conclusions in the owning plan or log. |

Read `../reference/cross-tool-coordination-contract.md` for request, attempt-record and brief
shapes, and `../reference/agent-task-coordination.md` for assignment and acceptance rules.

## Tested Support

The executable recipe was qualified with a Codex desktop coordinator, an independent ordinary
Codex relay chat with temporary transport workers, and CLI-owned Claude sessions on one macOS
host. The helper's consultation, continuation and denial paths were exercised with the app-bundled
Claude Code CLI 2.1.286; an older 2.1.153 CLI rejects `--permission-prompts`, which the helper
reports as `unsupported_invocation`. The helper builds for Windows, but Windows process tests,
live Claude and desktop wakeup, cross-host delivery and Claude-led reverse execution are not
qualified. The helper runs the executable directly; a Windows `.cmd` or `.bat` wrapper would be
interpreted by `cmd.exe` with its own quoting rules and is not supported. Host messaging behavior
can change; recheck it when a send is rejected or does not wake the coordinator.

## User-Facing Flow

Do not offer cross-tool delegation at every plan. Explain it when the user asks or when it is
materially useful, for example:

```text
We can continue with Codex. If you want Claude to implement or review, I can guide that setup;
it requires your own working access to Claude.
```

If the user declines, continue natively and do not offer again in this conversation or assignment.
Persist a decline only at a scope the user agrees to. Do not turn one person's choice into a team
default.

When the user selects Claude, reuse an existing clear choice and sufficient setup authority.
Otherwise ask one decision at a time:

1. Confirm the role and any model preference the user or local reference already set.
2. If the CLI is missing, explain what it is and use the tooling-readiness route below to offer
   installation; do not install merely because it is absent.
3. Explain that Claude use requires the user's own account or approved organization access and
   may be chargeable under that account. Link the current official guidance; create no account,
   purchase or usage commitment.
4. If sign-in is missing, ask the user to complete Claude's own sign-in flow in a terminal they
   can see. Passwords, tokens and MFA codes stay in that flow, never in chat or files.
5. If access still is not available, offer the choices: configure it now, use the current tool for
   this assignment, or pause.
6. Ask where to remember the agreed non-secret arrangement: only this assignment, the user's
   personal preferences, this repository, or the coordination home.
7. If no usable relay exists, show the exact proposed relay chat, return recipient and messages,
   and ask for that chat and message authority before creating anything.

End with the result in plain words: what was set up, where it is remembered, and that the
assignment can now be commissioned.

## Operator Notes

- Check CLI presence through `handle-tooling-readiness.md`. Check sign-in with the CLI's supported
  status command (for example `claude auth status`) and recheck after the user signs in. A missing
  login is expected first-use state, not a helper failure, and is handled before writing a request.
- Shared preferences do not prove readiness on another person's machine; authentication is per
  machine in the tool's own credential storage.
- Remember choices at the agreed scope: assignment record; personal ignored user layer
  (`.codeheart/user/`); repository reference `docs/repo/reference/agent-coordination.md`; or the
  coordination home's reference. Ordinary native operation needs no reference file.
- Shared references record policy and role choices: which CLI distribution, model, permission
  profile name and `permission_prompts` value each role uses, and who approved them. Exact host
  values, such as the absolute executable path, its observed version and profile file paths, stay
  in the ignored local layer (for example `.codeheart/user/agent-coordination.local.yaml`) or in
  the assignment's private record. Never commit machine paths to a shared reference.
- Create the repository reference from the installed starter
  `../templates/agent-coordination-reference.md` only when it is absent. When it exists, change
  only the specifically agreed entries. Kit repair, sync and upgrade never overwrite it.
- A member repository links to its coordination home's reference with a repository URL plus
  relative path, or an established local link. Verify home identity against the existing
  portfolio configuration when relying on membership. Never derive a path from a home ID, guess a
  chat by title, or enroll a repository silently.
- A relay locator is configuration, not authority to message or create chats. Missing, stale or
  mismatched routes produce a clear question to the reference owner.
- Keep machine paths, session artifacts and private evidence in ignored local state.

## Execution Path

### 1. Prepare

1. Read the consumer reference for relay, host, model, permission profile and evidence-root
   choices, then resolve them into concrete values from the assignment record or the ignored
   local settings. Reuse values that are already approved there; ask the owner or user only for a
   genuinely missing decision or access. Never let the relay choose them.
2. Use the executable path that those approved settings name. Do not substitute whatever `claude`
   a `PATH` lookup finds unless the approved setting says to. Confirm the file exists and run it
   with `--version`; record the observed version with the attempt. If the version differs from
   the one recorded for the approved setting, note the change; an unsupported flag will surface
   as `unsupported_invocation`, never as a silent fallback.
3. Select the Operating Kit helper the relay will run, using the exact path named in the ignored
   local settings or a verified installation. Confirm it yourself before dispatch: `--version`
   reports a release that includes `coordination`, and `coordination --help` lists
   `invoke-claude` and `record-delivery`. The repository's installed guidance version does not
   prove that the relay's `PATH` resolves the same helper; an older one rejects the subcommand
   before Claude starts. If no suitable helper is available, follow `handle-tooling-readiness.md`.
   Do not silently install or update, change `PATH`, or leave the transport worker to look for
   alternatives.
4. Write the brief from the owning records. For an approved implementation plan, use the native
   goal by default: open the brief with the supported `/goal <condition>` invocation for the
   current phase, aimed at a real handoff such as source review. Honor an explicit opt-out.
5. Write the request JSON with every concrete value, a new `attempt_id`, the exact return chat,
   and `session.mode` `new`, or `resume` with the retained session for the same assignment.
   Set `working_directory` to the repository this attempt changes; when work moves to another
   repository, resume the same session in a new attempt from that repository. With an explicit
   `dontAsk` profile, confirm a new profile once with harmless probes first; see the contract's
   permission profiles section.
6. Note the attempt directory `<state_root>/attempts/<assignment_id>/<attempt_id>/` in the
   assignment record before dispatch.

### 2. Dispatch And Yield

Send the relay one compact envelope and then yield to the user:

```text
Transport request for assignment <assignment_id>, attempt <attempt_id>.
Run: "<verified helper path>" coordination invoke-claude --request "<request path>"
Wait on that same process through ordinary tool output yields; do not restart it.
Then send exactly one notice to send_message_to_thread, loading its arguments unchanged:
- if the helper output names a message_file, load that file and afterwards run
  "<verified helper path>" coordination record-delivery for that attempt (sent, rejected or
  uncertain, with the original receipt or error);
- otherwise, if it contains fallback_message, send that object; do not run record-delivery;
- otherwise, send the helper's original output to the authorized recipient below.
Never load a message file the helper did not name in this output.
Authorized recipient: <coordinator chat> on <host>. Do not read, summarize or act on the reply,
and do not retry the helper. Use only the helper path given above; if it fails to start, send its
original error to the recipient instead of looking for another helper.
```

The helper names a `message_file` only when this invocation wrote it, which covers every outcome
after the attempt directory exists: replies, CLI failures, `session_locked`, `launch_failed` and
`helper_interrupted`. A `fallback_message` covers failures before an attempt record exists, a
message file that could not be written, and `attempt_exists`. For `attempt_exists` the notice is
a new refusal, not the earlier result: the earlier attempt, message and delivery record stay
unchanged and nothing is resent. A fallback send is not recorded in any attempt directory; if it
also fails, the relay reports that in its own reply. Nothing is queued or retried.

With several active assignments, the relay dispatches one temporary transport worker per
assignment and returns control. Workers send directly to the requesting ordinary chat. Do not
route results through a native parent or subagent recipient.

While work runs, the coordinator may discuss other matters. It does not poll transcripts or
duplicate the delegated investigation. It actively waits only when its next decision depends on
the result.

### 3. Receive And Assess

1. Read the original reply (inline, or the referenced `stdout.jsonl` result field) and
   `attempt.json`. Check status, actual session, exit and permission denials.
2. `response_captured` means a reply exists. It is not acceptance. You are the primary reviewer:
   examine the delivered work and evidence against the intended outcome and scope, and verify
   consequential claims; do not re-run the whole investigation. Commission an additional
   independent reviewer only when complexity, risk, your own implementation involvement or a
   binding gate warrants it, and still assess its findings yourself.
3. For a consequential review, report to the user the original locator, core findings and their
   significance, your own assessment including disagreement, recommended disposition and residual
   uncertainty. During discovery, planning or discussion, report before amending and follow
   `../../planning-workflows/runbooks/review-planning-document.md`.
4. For a question, answer from the owning decision or ask the user when it exceeds your mandate,
   then resume the same session with a new attempt. During authorized implementation, routine
   corrections continue under the existing grant; material changes return to the owner.
5. Record accepted conclusions, the substantive session locator and acceptance in the owning plan
   or log. Do not author a separate report per message.

### 4. Recover

- No notice arrived but the result is needed: read the known attempt directory once.
  `delivery.json` shows `pending`, `sent`, `rejected` or `uncertain`. Pending can mean the worker
  stopped before sending. Do not resend or relaunch automatically.
- `session_locked`: another attempt holds the session. Find it from the recorded lock owner.
- `helper_interrupted` or a record left `preparing`, `launching` or `running`: effects are
  uncertain. Inspect the child process, the original output, the worktree and external state
  before anything else. When both recorded processes have ended, run `coordination release-lock`
  for that exact attempt. If the CLI process was never recorded, verify manually that no CLI
  process for the session remains and pass that statement with `--manual-verification`.
- `release-lock` stays refused while a recorded process ID is alive, which includes a reused ID,
  or when the lock is empty or unreadable. No statement overrides those refusals, and there is no
  automatic recovery. The owner either establishes actual ownership and exit (for example, the
  process ID now belongs to an unrelated program and the original CLI has ended), records that
  verification in the assignment record and only then removes that one lock file, or hands the
  work to an explicitly new
  session that keeps the uncertain attempt as evidence. Before resuming overlapping work,
  reconcile the old attempt's effects and make sure no competing writer remains.
- `unsupported_invocation`: the installed CLI rejected the fixed flags. Report the original
  stderr; do not retry with other flags. Check official CLI documentation and update the recipe.
- Never delete a lock silently, replay an uncertain action, or let two writers drive one session.
  A local lock cannot detect a writer on another host or outside the helper.

## Stop Conditions

- The selected provider is unavailable and the user has not chosen native work or a pause.
- A required value, relay, recipient or chat authority is missing, stale or ambiguous.
- A permission denial or failed gate blocks the next dependent step. Return it unchanged; do not
  rephrase, escalate or route it through another agent.
- The work would exceed the assignment, change permissions, or need a decision outside the
  coordinator's mandate.

## Evidence And Validation

- Attempt directories, raw output and message files stay ignored and private. Do not print them
  in full, commit them or send them to other recipients.
- Durable records keep only accepted findings, decisions and the session and attempt locators.
- After acceptance, once no correction, review or recovery needs them, generated duplicates and
  logs may be removed through ordinary authorized cleanup. Keep unfinished attempts and locks.
- Validate a new arrangement with one harmless consultation that returns a reply to the
  coordinator after it yields, before relying on it for implementation.
