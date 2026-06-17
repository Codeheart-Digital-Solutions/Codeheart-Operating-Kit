Last updated: 2026-06-17T06:40:50Z (UTC)

# Agent Memory

This managed domain owns reusable rules for curated agent memory.

The Operating Kit owns memory format and maintenance procedure. Consumer repositories own actual
memory state.

## Use

Use agent memory when the user asks where prior work left off, asks to continue a previous thread,
or clearly references an ongoing workstream that may have session history.

Do not read agent memory by default for every task. Verify the current repository state and the
relevant canonical docs before acting on a memory summary.

## Source-Of-Truth Order

Use this order when recovering context:

1. Current repository state: files, diffs, plans, tests, validation outputs, PRs, and current docs.
2. Curated agent memory: goals, workstreams, session summaries, handoffs, and untriaged sessions.
3. Raw local agent state: supporting evidence when curated memory is incomplete or needs
   verification.

Session IDs are resumability handles when available. Preserve them, but do not rely on them as the
only record of status because local session state may move, archive, delete, or change format.

## Freshness Check

Before using memory to resume work, compare current UTC time with the `Coverage window:` in the
consumer memory files.

If local sessions continued after the coverage timestamp, treat curated memory as incomplete until
the maintenance runbook has inspected the newer activity.

Coverage windows should include timestamp precision, not only a date, so same-day continuation can
be detected.

## Routes

- Entry format: `reference/entry-format.md`
- Session-ledger maintenance: `runbooks/session-ledger-maintenance.md`

## Consumer State

Consumer memory state is scaffolded outside managed kit content, normally under
`docs/agent-memory/`. Sync may create missing scaffold files but must not overwrite existing
consumer memory state.

Typical consumer-owned files:

- `goal-register.md`
- `session-ledger.md`
- `untriaged-sessions.md`
- `archive/`

## Durable Reference Rules

- Use repo-relative paths for repository files.
- Use `$CODEX_HOME` or `$HOME/.codex` patterns for local agent state.
- Use session IDs when available.
- Do not commit raw transcript dumps.
- Do not commit concrete user home paths, editor profile paths, terminal-specific paths, or
  machine-specific absolute checkout paths.
- Do not treat memory as a replacement for canonical product, business, research, discovery,
  implementation, runbook, or reference docs.
