Last updated: 2026-06-25T13:05:46Z (UTC)

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
- When a module onboarding or operation is blocked by missing local tooling, follow the managed
  tooling readiness route before installing, repairing, improvising setup, or declaring the
  capability unavailable.

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
- Operating Kit feedback:
  `.codeheart/kit/docs/agent-interface/runbooks/submit-kit-feedback.md`
- Tooling readiness:
  `.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`
- Structure governance: `.codeheart/kit/docs/structure-governance/README.md`
- Module and extension state:
  `.codeheart/kit/docs/structure-governance/reference/module-extension-state.md`
- Full kit inventory and fallback: `.codeheart/kit/README.md`

## Weekly Update Check

At the start of each agent session, inspect `.codeheart/kit.lock.yaml`. If
`next_update_check_due` is in the past, run `codeheart-operating-kit update-check`.

Stay silent when the installed kit is current. When an update is available, mention it briefly and
ask before applying anything.

Do not edit the due date manually. `codeheart-operating-kit update-check` owns
`last_update_check_at`, `next_update_check_due`, `latest_seen_version`, and `update_status`.

<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->

# Repository-Owned Instructions

Add repository-specific rules below this heading. Keep local safety rules, build and test commands,
product documentation routes, release procedures, and repository-specific exceptions here.

# Local User Guidance

Add local user guidance below this heading or link to `.codeheart/user/` when present. Do not place
private local preferences inside the managed block.
