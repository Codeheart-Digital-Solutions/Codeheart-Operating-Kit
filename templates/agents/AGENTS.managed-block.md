Last updated: 2026-06-13T23:02:10Z (UTC)

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

## Managed Routes

- Discovery: `.codeheart/kit/docs/planning-workflows/runbooks/discovery-workflow.md`
- Implementation planning:
  `.codeheart/kit/docs/planning-workflows/runbooks/draft-implementation-plan.md`
- Implementation execution:
  `.codeheart/kit/docs/planning-workflows/runbooks/execute-implementation-plan.md`
- Planning document review:
  `.codeheart/kit/docs/planning-workflows/runbooks/review-planning-document.md`
- Agent memory: `.codeheart/kit/docs/agent-memory/README.md`
- Agent interface: `.codeheart/kit/docs/agent-interface/README.md`
- Structure governance: `.codeheart/kit/docs/structure-governance/README.md`
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
