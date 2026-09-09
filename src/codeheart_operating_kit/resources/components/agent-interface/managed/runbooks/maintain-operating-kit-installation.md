Last updated: 2026-09-09T14:08:21Z (UTC)

# Maintain An Operating Kit Installation

Use this runbook to choose the smallest lifecycle command that matches installed state. The command
name carries the intent; routine repair and same-version sync do not require an extra confirmation
prompt.

Audience: agent-facing

Intent:
Choose the supported lifecycle command and complete the assigned installation or repository rollout.

Success:
The intended version is healthy at each named target, authored content is preserved, and actual
local, published-branch and default-branch adoption states are distinguished.

Agent judgment boundary:
Choose the smallest supported lifecycle operation and preservation-safe adoption branch. Reuse
sufficient authority; do not expand it to held product history or unspecified rollout targets.

Stop boundary:
Preserve state on real drift/recovery failure, overlapping authored work, missing authority or a
failed required gate. Never hand-edit installed state or bypass target protections.

## State-To-Command Route

| Observed state or intent | Command | Boundary |
| --- | --- | --- |
| No installation; folder is absent or adoptable | `init --dry-run`, then `init` | Creates lock v2 and preserves existing consumer files. |
| Diagnose any state | `check` | Read-only; returns one primary state and actionable blockers. |
| Compatible drift, partial installation, or same-version lock-v1 migration | `repair --dry-run`, then `repair` | Restores the running version; never authorizes a version change. |
| Refresh embedded files for the already installed version | `sync --dry-run`, then `sync` | Uses only the running binary and preserves the installed kit version. |
| User explicitly asks whether a newer release exists | `update-check` | Optional manual lookup for a valid lock-v2 installation; due metadata never triggers it. |
| User approved a newer release | `upgrade --version <version> --dry-run`, then `upgrade --version <version> --yes` | Only command that may change kit version; verifies catalog-to-binary provenance first and may atomically migrate a pristine compatible lock v1 to v2. |
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

## Repository Adoption And Worktree Reconciliation

For a repository rollout, inputs are the verified released version, named default/integration
branches, active worktrees, publication/merge/reconciliation authority and acceptance owner.
Route before choosing Git or provider operations; use ordinary Git plus the lifecycle commands
above. A local pilot is a narrower outcome and must be reported as such.

1. Inspect current target refs, instructions, installation, dirty/overlapping work and unpublished
   ancestry. Recheck actual target CI triggers, required checks and merge policy. Producer release
   evidence does not rewrite consumer CI. Return a disproportionate enforced gate or preservation
   conflict to its owner before the affected merge; continue independent authorized preparation.
2. Prepare an isolated Kit-only adoption branch from the current intended main/integration baseline.
   Do not merge an unfinished product branch to distribute guidance. An existing adoption commit
   can inform a compatible change, but the target's own baseline controls lifecycle and preservation.
3. Run `check`, preview the supported released upgrade, then apply under covered authority. If a
   newer shared CLI calls a pristine older tree partial, verify it with the matching verified
   published old CLI and initiate the supported forward upgrade there. A genuine drift, partial or
   recovery failure remains a blocker; never relabel it healthy or hand-edit a lock.
4. Verify intended installed version, managed routes and preservation of config, module state,
   product plans, local-user files and instructions outside the managed block. Inspect the exact
   adoption diff and all ancestors that a push would publish. Commit only the adoption changes.
5. Publish the scoped branch/PR and merge under covered authority when applicable review and
   required checks pass. Recheck current target changes before integration; resolve relevant
   overlap and rerun invalidated evidence. Verify the actual remote target lock/content and a
   healthy checkout after integration before claiming repository adoption.
6. After main adoption, reconcile that released guidance into each assigned active worktree using
   supported lifecycle operations or a preservation-safe compatible Git change. Inspect current
   owner work again; do not overwrite it or merge held product history merely for adoption.
   Verify health/preservation and commit the bounded change locally. Push only if the entire
   unpublished ancestry is authorized. A held-ancestry local commit is legitimate: record the
   concrete exception, product owner and resumption event without claiming branch publication.
7. Report named main adoption, worktree reconciliation and any remaining exceptions separately.
   Publish final plan/log closure through its covered route only after the actual required effects.
   A source merge, public release, branch push or healthy pilot is not full repository rollout.

Use `../../planning-workflows/reference/planning-document-lifecycle.md` for Git timing and
exceptions. Retain applicable source review for mechanical adoption; add target-specific checks
and required owner gates. Git events alone do not require unrelated broad product suites. This is
L1 structured guidance using existing tools; no new Git wrapper or evidence registry is needed.

## Common Blockers

| Blocker | Action |
| --- | --- |
| `managed_path_modified` | Preserve or move the modified retired file, then retry the named command. |
| `version_change_requires_upgrade` | Preview the named upgrade; do not force repair or sync. |
| `upgrade_requires_clean_legacy_v1_installation` | Preserve the reported state; resolve invalid, missing, modified, unsafe, or ambiguous legacy authority before retrying. Do not edit the lock manually. |
| `transaction_in_progress` | Wait for the owning process and run `check`. |
| `recovery_required` | Preserve `.codeheart/kit.transaction.json` and its recovery directory; diagnose before retry. |
| `schema-invalid` or `unsupported-future-version` | Do not coerce or rewrite state; use a compatible CLI or repair path. |
| `cli_unavailable` | Follow the tooling-readiness route before installing or improvising tooling. |

## Preservation And Stop Conditions

Config, repository instructions outside the managed block, plans, memory, local-user state, and
create-once scaffolds are consumer-owned and must remain unchanged. A legacy forward upgrade binds
the exact source lock and rechecks its declared managed-file checksums before target-side
reconciliation. Stop before hand-editing the lock or managed kit files, deleting recovery evidence,
forcing a version change, trusting an unverified catalog or pack, or publishing a release.

Tooling readiness route:
`.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`.
