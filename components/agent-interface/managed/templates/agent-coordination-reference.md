Last updated: 2026-10-10T08:30:00Z (UTC)

# Agent Coordination

<!-- Starter from the Operating Kit's agent-coordination-reference template. Copy it to
docs/repo/reference/agent-coordination.md only after the owner agrees a cross-tool arrangement,
replace every placeholder, update the timestamp and remove this comment. The copied file is owned
by this repository; Kit repair, sync and upgrade do not change it. -->

This repository-owned reference records the agreed, non-secret cross-tool coordination
arrangement. Procedure:
`.codeheart/kit/docs/agent-interface/runbooks/coordinate-cross-tool-task.md`. Ordinary work with
the current tool does not need this file.

## Arrangement

| Setting | Value |
| --- | --- |
| Maintenance owner | `<role or team that may change this reference>` |
| Scope | `<this repository / coordination home for listed members>` |
| Supported host | `<tested host and coordinator tool, e.g. Codex desktop on macOS>` |
| Relay | `<ordinary relay chat locator and host>` |
| Report-back route | `<how directors are addressed; the requesting ordinary chat by default>` |
| Preferred substantive model | `<model per role, or "ask">` |
| Relay worker model | `<model and reasoning, or "relay default">` |
| Approved CLI | `<distribution, e.g. app-bundled Claude Code CLI; who approved it>` |
| Permission profile per role | `<profile name, permission mode and permission_prompts value per role; approving record>` |
| Local execution settings | `<ignored file holding exact paths, e.g. .codeheart/user/agent-coordination.local.yaml, including the verified Operating Kit helper the relay runs>` |
| Evidence root | `<ignored local path convention, e.g. .codeheart/local/coordination>` |

Exact host values stay out of this file: the absolute CLI path, its observed version, the
Operating Kit helper path and profile file paths live in the ignored local settings named above
or in an assignment's private record.
Coordinators resolve them from there, verify the executable and its version before each
commission, and never fall back to an unapproved `PATH` lookup. Sign-in lives in each tool's own
credential storage on each machine; nothing here proves that a user is signed in.

## Appointments And Members

- Program or plan appointments: `<links to the records that appoint advisors or reviewers>`
- Member repositories that route here: `<repository URL plus this path, or "none">`
- Coordination home, when this repository is a member: `<home repository URL plus relative path>`

## Change Log

- `<date>`: `<what changed, who agreed and which authority covered it>`
