Last updated: 2026-07-10T11:29:07Z (UTC)

# Maintain An Operating Kit Installation

Use this runbook to choose the smallest lifecycle command that matches installed state. The command
name carries the intent; routine repair and same-version sync do not require an extra confirmation
prompt.

## State-To-Command Route

| Observed state or intent | Command | Boundary |
| --- | --- | --- |
| No installation; folder is absent or adoptable | `init --dry-run`, then `init` | Creates lock v2 and preserves existing consumer files. |
| Diagnose any state | `check` | Read-only; returns one primary state and actionable blockers. |
| Compatible drift, partial installation, or lock-v1 migration | `repair --dry-run`, then `repair` | Restores the running version; never authorizes a version change. |
| Refresh embedded files for the already installed version | `sync --dry-run`, then `sync` | Uses only the running binary and preserves the installed kit version. |
| User explicitly asks whether a newer release exists | `update-check` | Optional manual lookup for a valid lock-v2 installation; due metadata never triggers it. |
| User approved a newer release | `upgrade --version <version> --dry-run`, then `upgrade --version <version> --yes` | Only command that may change kit version; verifies catalog-to-binary provenance first. |
| Active transaction | `check` and wait | Do not start concurrent lifecycle work. |
| Dead verified pre-commit transaction | `repair --dry-run`, then `repair` | Stale takeover requires process-liveness and transaction-identity proof. |
| Recovery required, schema-invalid, or unsupported future state | `check` | Stop; preserve evidence and follow the returned blocker. |

## Normal Procedure

1. Run `codeheart-operating-kit check <repository>`.
2. Select one command from the table. Do not substitute `init` for repair or `sync` for upgrade.
3. Use `--dry-run` before a material repair, sync, or upgrade when the change set is not already
   obvious.
4. Run the selected command. Direct invocation authorizes `init`, `repair`, and same-version
   `sync`; only cross-version `upgrade` also requires `--yes`.
5. Run `check` after manual recovery or when the command reports a blocker. Successful lifecycle
   transactions already perform their own post-check and remove successful transaction data.

## Common Blockers

| Blocker | Action |
| --- | --- |
| `managed_path_modified` | Preserve or move the modified retired file, then retry the named command. |
| `version_change_requires_upgrade` | Preview the named upgrade; do not force repair or sync. |
| `transaction_in_progress` | Wait for the owning process and run `check`. |
| `recovery_required` | Preserve `.codeheart/kit.transaction.json` and its recovery directory; diagnose before retry. |
| `schema-invalid` or `unsupported-future-version` | Do not coerce or rewrite state; use a compatible CLI or repair path. |
| `cli_unavailable` | Follow the tooling-readiness route before installing or improvising tooling. |

## Preservation And Stop Conditions

Config, repository instructions outside the managed block, plans, memory, local-user state, and
create-once scaffolds are consumer-owned and must remain unchanged. Stop before hand-editing
managed kit files, deleting recovery evidence, forcing a version change, trusting an unverified
catalog or pack, or publishing a release.

Tooling readiness route:
`.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`.
