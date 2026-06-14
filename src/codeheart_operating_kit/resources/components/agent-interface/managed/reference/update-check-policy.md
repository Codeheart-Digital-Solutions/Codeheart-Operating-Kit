Last updated: 2026-06-13T23:55:46Z (UTC)

# Update Check Policy

The Operating Kit uses a static root `AGENTS.md` trigger plus lockfile state for weekly update
checks.

G1 resolves latest-version metadata from the public GitHub latest-release endpoint unless a
specific metadata URL is provided for tests, mirrors, or controlled environments.

## Session Trigger

At the start of each agent session, inspect `.codeheart/kit.lock.yaml`. If
`next_update_check_due` is in the past, run `codeheart-operating-kit update-check`.

Do not apply a kit version update unless the user asks for update, sync, or kit maintenance.

The trigger is static so root `AGENTS.md` does not churn weekly. The changing due date belongs in
the lockfile only.

## Weekly Cadence

When `codeheart-operating-kit update-check` succeeds, it writes:

```yaml
last_update_check_at: <now>
next_update_check_due: <now + 7 days>
latest_seen_version: <latest release version>
update_status: current | update-available
```

The command should update `next_update_check_due` after a successful latest-version check whether
the installed kit is current or an update is available. Failed checks may record attempt metadata,
but they should preserve the previous successful cadence state unless the CLI has a stronger
diagnostic model.

## User Notification

Stay silent when the installed kit is current. Mention an available update briefly and ask whether
the user wants to apply it. Mention failed checks only during kit-related work or repeated failure
attention.
