Last updated: 2026-06-26T14:30:34Z (UTC)

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

Only `upgrade --yes` may change the installed kit version. `repair` and `sync` restore or refresh
the currently installed version without an additional confirmation prompt.

Do not edit the due date manually. `codeheart-operating-kit update-check` owns
`last_update_check_at`, `next_update_check_due`, `latest_seen_version`, and `update_status`.

<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->

# Repository-Owned Instructions

Add repository-specific rules below this heading. Keep local safety rules, build and test commands,
product documentation routes, release procedures, and repository-specific exceptions here.

## Producer Authority

This repository is the Operating Kit producer. For source, release, state, lifecycle, installer, or
routing changes, use tracked source under `components/`, `profiles/`, `templates/`, `schemas/`,
`internal/`, `scripts/`, and `docs/repo/` as authority. Never use the ignored consumer installation
under `.codeheart/kit/` as source truth or copy changes from it back into producer files.

Route implementation through `docs/repo/runbooks/change-operating-kit.md` and release work through
`docs/repo/runbooks/release-operating-kit.md`. Consumer lifecycle guidance is authored at
`components/agent-interface/managed/runbooks/maintain-operating-kit-installation.md`.

# Local User Guidance

Add local user guidance below this heading or link to `.codeheart/user/` when present. Do not place
private local preferences inside the managed block.


# Existing Instructions Preserved

Last updated: 2026-06-13T22:47:44Z (UTC)

# Repository Guidelines

## Purpose And Read Order

`AGENTS.md` is the first agent-facing bootstrap contract for Codeheart Operating Kit maintainers.
Keep it short, public-safe, and routing-focused.

Use task-matched documentation reading:

- Read `README.md` when the task involves repository purpose, public scope, or unsure placement.
- Read `docs/README.md` before creating, moving, or reorganizing documentation.
- Read `docs/repo/reference/placement-contract.md` before adding a new docs area, generated
  surface, component target, or externally visible path.
- Read `docs/repo/reference/consumer-impact-classification.md` before changing managed content,
  scaffolds, templates, schemas, validators, installers, release assets, or sync/check behavior.
- Read `docs/repo/runbooks/change-operating-kit.md` before changing kit source, docs, schemas,
  templates, validators, installers, or CLI behavior.
- Read `docs/repo/runbooks/release-operating-kit.md` before tagging, packaging, publishing, or
  updating release notes.
- Read `docs/repo/runbooks/promote-consumer-change.md` before moving consumer-local guidance into
  this kit.

Do not read every document by default. Read the smallest matching route and then follow links when
the task requires deeper context.

## Public-Core Safety

This repository is public. Do not add private Codeheart business records, customer names, tenant
details, credentials, secrets, local machine state, account identifiers, restricted strategy, or raw
operational logs.

Before committing copied or extracted material, classify it as reusable public doctrine,
public-safe placeholder content, plan-scoped evidence, or excluded private material. Sanitize
source-specific names and paths unless the public identifier is intentionally part of the kit.

## Source Of Truth

- `README.md` owns repository purpose and public boundary.
- `docs/repo/reference/placement-contract.md` owns placement rules for kit docs and generated
  consumer surfaces.
- `docs/repo/reference/consumer-impact-classification.md` owns impact classes for kit changes.
- Runbooks own ordered maintainer procedures.
- Component manifests, schemas, and release manifests own machine-readable behavior after they
  exist.

When instructions conflict, preserve public-core safety, avoid destructive or public-state-changing
actions, and record the inconsistency before proceeding.

## Change Safety

- Protect existing work. Do not overwrite, delete, or revert unrelated changes.
- Do not run destructive cleanup unless explicitly requested.
- Do not publish a release, push a tag, or change public GitHub settings unless the task explicitly
  asks for release or repository governance work.
- Do not print secrets or tokens. Mask sensitive values in logs.
- Keep docs and implementation changes in the same PR unless the design is still uncertain or the
  change is high-risk.

## Validation

Run the smallest validation that proves the changed surface. For repository-wide changes, run the
available markdown, schema, public-core, CLI, and release validation commands once they exist. If a
validator does not exist yet, record the residual risk in the PR or execution log.
