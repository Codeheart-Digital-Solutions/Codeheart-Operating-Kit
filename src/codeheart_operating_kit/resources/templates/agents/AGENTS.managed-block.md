Last updated: 2026-07-10T11:29:07Z (UTC)

<!-- BEGIN CODEHEART OPERATING KIT MANAGED BLOCK -->

# Codeheart Operating Kit

## Immediate Rules

- Protect existing work. Do not overwrite, revert, clean up, or delete unrelated user changes.
- Do not run destructive actions or external-state-changing actions without an explicit user
  request or matching local runbook.
- Never expose, commit, or print secrets or tokens.
- Read task-matched instructions instead of sweeping every document.
- Treat `.codeheart/kit/` as managed Operating Kit content. Do not hand-edit it outside explicit
  drift, repair, or sync work.
- Keep reusable generic doctrine in the Operating Kit. Repository-owned sections should contain
  only local product, project, command, operational, and exception guidance.
- For structural, external, sensitive, module, product, or ambiguous work, route before selecting
  a tool, connector, API, browser, script, or runbook. Full routing reference:
  `.codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md`.
- Before operating installed modules or extensions, discover the module or extension ID through
  the module system present in the repo, then check `docs/repo/state/<id>/` for committed
  non-secret routing state. Local state routes the work; live external preflight still decides.
- When a repository, module, extension, or agent task is blocked by missing local tooling, follow
  the managed tooling readiness route before installing, repairing, improvising setup, or
  declaring the capability unavailable.

## Managed Routes

- Discovery: `.codeheart/kit/docs/planning-workflows/runbooks/discovery-workflow.md`
- Implementation planning:
  `.codeheart/kit/docs/planning-workflows/runbooks/draft-implementation-plan.md`
- Implementation execution:
  `.codeheart/kit/docs/planning-workflows/runbooks/execute-implementation-plan.md`
- Planning document review:
  `.codeheart/kit/docs/planning-workflows/runbooks/review-planning-document.md`
- Plan registers and configured portfolio coordination:
  `.codeheart/kit/docs/planning-workflows/runbooks/maintain-plan-register.md`
- Agent memory: `.codeheart/kit/docs/agent-memory/README.md`
- Agent interface: `.codeheart/kit/docs/agent-interface/README.md`
- Operation routing and dispatch:
  `.codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md`
- Operating Kit installation lifecycle:
  `.codeheart/kit/docs/agent-interface/runbooks/maintain-operating-kit-installation.md`
- Repo feedback capture:
  `.codeheart/kit/docs/agent-interface/runbooks/capture-repo-feedback.md`
- Operating Kit feedback:
  `.codeheart/kit/docs/agent-interface/runbooks/submit-kit-feedback.md`
- Tooling readiness:
  `.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`
- Structure governance: `.codeheart/kit/docs/structure-governance/README.md`
- Module and extension state:
  `.codeheart/kit/docs/structure-governance/reference/module-extension-state.md`
- Full kit inventory and fallback: `.codeheart/kit/README.md`

## Optional Update Checks

Do not check for a new Operating Kit version at session start or merely because
`next_update_check_due` is in the past. Run `codeheart-operating-kit update-check` only when the
user asks or the current task explicitly requires latest-release information.

`update-check` changes metadata only. Its due-date fields are informational and never create a
mandatory or background check.

Only `upgrade --yes` may change the installed kit version. `repair` and `sync` restore or refresh
the currently installed version without an additional confirmation prompt.

Do not edit update metadata manually. When invoked, `codeheart-operating-kit update-check` owns
`last_update_check_at`, `next_update_check_due`, `latest_seen_version`, and `update_status`.

<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->

# Repository-Owned Instructions

Add repository-specific rules below this heading. Keep local safety rules, build and test commands,
product documentation routes, release procedures, and repository-specific exceptions here.

# Local User Guidance

Add local user guidance below this heading or link to `.codeheart/user/` when present. Do not place
private local preferences inside the managed block.
