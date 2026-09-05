Last updated: 2026-09-05T19:15:20Z (UTC)

# Codex Task Operations

This optional note maps `agent-task-coordination.md` to Codex surfaces. The generic assignment and
acceptance contract remains authoritative. Consumers need neither Codex nor these tools to use it.

Official references checked on 2026-09-05:

- [Projects and chats](https://learn.chatgpt.com/docs/projects) describes project context, task
  naming, pinning, archival and restoring archived chats.
- [Long-running work](https://learn.chatgpt.com/docs/long-running-work) documents Goal mode and its
  progress controls; a goal keeps the same sandbox and approval policy.
- [Git worktrees](https://learn.chatgpt.com/docs/environments/git-worktrees) explains separate
  checkouts, handoff and managed-worktree cleanup. Archiving a managed-worktree chat can cause
  worktree deletion; permanent worktrees have different retention behavior.

Recheck official documentation and the currently exposed tool contracts when operating these
features. Tool availability, permissions and host context can differ between tasks. The surface
names below describe available desktop contracts at the checked date, not a portable API promise.

## Project And Task Operations

When the user explicitly asks for a new task, inspect available projects with `list_projects` and
select the owning source project. Use `create_thread` according to its current local/worktree
contract, including an explicitly assigned existing branch when applicable. Verify the resulting
worktree and branch before editing; a created task may start detached or still be initializing.
Protect existing worktrees and avoid assigning two implementers the same mutable branch.

A task created for the entire approved plan should retain that plan's full outcome, ordered epics,
Git boundary, constraints, acceptance owner and report-back instruction. A slash-command string
inside its prompt does not prove that a goal was activated. A short-lived reviewer is distinct
from a new user-owned implementation task; use it only under the applicable request and tools.

Use the supported `list_threads`/`read_thread` surfaces to resolve a requested task by its actual
name and project before acting. Use the current tools for rename, pin, handoff or archival only
within the requested scope. Do not perform these actions merely because this note lists them.

## Direct Report-Back

When the assignment explicitly authorizes reporting to the commissioning task, use the current
`send_message_to_thread` contract with that destination and a concise human-readable prompt:
Plan and epic names, result/evidence, Git/PR state, and the director decision needed. Task IDs are
execution locators, not durable organizational role IDs. Include only information appropriate to
the destination's audience.

Verify the tool's result. A successful send proves dispatch, not acceptance. If unavailable or
rejected, keep the report in the execution log/final result and say it was not delivered. Do not
switch to another destination, bypass permission review, or install a reporting framework. No
heartbeat, cron or continuous polling is needed; ordinary task follow-up carries the decision.

## One Explicit Whole-Plan Goal

The user can request a goal for the entire plan, for example:

```text
Start one goal to execute <Plan name> at <canonical path> in <owning project>,
through all ordered epics to <the agreed finish line>. Preserve <constraints>.
Use the plan's approved commits, normal pushes and PR checkpoints. Send required
reviews and result reports to <commissioning task> for <Director role> acceptance.
Keep reserved merge/release/adoption effects with their named owner until delegated.
Verify goal activation; do not infer a token budget or mark it complete at review handoff.
```

Use the supported goal creation mechanism only after that explicit request. Where exposed,
`create_goal` followed by `get_goal` can verify the actual objective and active state; otherwise
inspect the supported app goal control. Do not report activation from printing `/goal` or merely
putting it in another task's prompt. If activation cannot be verified, disclose that limit and
continue already authorized useful work under the plan.

Keep the goal incomplete at required review handoff. It creates no new permissions and does not
require repeated polling. Match completion to the whole agreed finish line, including release or
adoption when covered. Report technical completion and real-use pilot status separately. Follow the
current goal tool's status rules; do not invent a resume API when only user progress controls can
resume a blocked goal.

## Archival And Preservation

Before requested archival, apply the generic acceptance/transfer and active-descendant checks.
Preserve committed work, unfinished work context and reachable review ownership. Inspect whether
the task uses a managed or permanent worktree and account for the documented cleanup consequences.
Use the supported archival tool, such as `set_thread_archived`, only after those checks and within
current authority. UI archival is not plan completion or an instruction to delete branches, erase
organizational records or clean up other worktrees.
