Last updated: 2026-07-10T11:29:07Z (UTC)

# Update Check Policy

Operating Kit version checks are manual and optional. There is no session-start, background, or
cadence-based requirement to look for a newer release.

G1 resolves latest-version metadata from the public GitHub latest-release endpoint unless a
specific metadata URL is provided for tests, mirrors, or controlled environments.

`update-check` is valid-install-only. It requires a schema-valid lock-v2 installation, changes only
update metadata plus the normal operation generation record, and never repairs content or changes
the installed kit version. Route compatible legacy or damaged state through `repair` first.

## Invocation Boundary

Run `codeheart-operating-kit update-check` only when the user asks whether a newer release exists
or the current authorized task explicitly requires latest-release information. Do not run it at
session start, on a background schedule, or merely because `next_update_check_due` is in the past.

Do not apply a kit version update unless the user explicitly approves `upgrade --yes`. `sync` is a
same-version refresh and is not version-change authorization.

## Compatibility Metadata

When `codeheart-operating-kit update-check` succeeds, it writes:

```yaml
last_update_check_at: <now>
next_update_check_due: <now + 7 days>
latest_seen_version: <latest release version>
update_status: current | update-available
```

These fields remain for lockfile and CLI compatibility. `next_update_check_due` is informational;
it does not instruct an agent or process to perform a future check. Failed checks may record
attempt metadata, but they should preserve the previous successful metadata unless the CLI has a
stronger diagnostic model.

## User Notification

Respond in the context that requested the check. Agent-notification mode may remain silent when the
installed kit is current. Mention an available update briefly and ask whether the user wants to
apply it; never apply it automatically.
