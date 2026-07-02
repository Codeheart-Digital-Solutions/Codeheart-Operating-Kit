Last updated: 2026-07-02T13:16:41Z (UTC)

# Codeheart Operating Kit Inventory

This file is the installed fallback inventory for managed Operating Kit content.

Agents should prefer the direct routes in the root `AGENTS.md` managed block. Use this inventory
when a task is unclear, the direct route does not match, the user asks about kit structure, or kit
state looks missing, stale, or damaged.

## Managed Domains

- Planning workflows: `.codeheart/kit/docs/planning-workflows/README.md`
- Agent memory: `.codeheart/kit/docs/agent-memory/README.md`
- Agent interface: `.codeheart/kit/docs/agent-interface/README.md`
- Structure governance: `.codeheart/kit/docs/structure-governance/README.md`
- Native Codex capabilities: `.codeheart/kit/docs/native-codex-capabilities/README.md`

## Operation Routing

- Route-before-surface standard:
  `.codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md`

## Generated State

- Lock metadata: `.codeheart/kit.lock.yaml`
- Shared non-secret config: `.codeheart/kit.config.yaml`
- Local user layer: `.codeheart/user/`
- Local machine/runtime layer: `.codeheart/local/`

## Module And Extension State

- Committed non-secret module/extension state convention:
  `.codeheart/kit/docs/structure-governance/reference/module-extension-state.md`
- Consumer-owned module/extension state root:
  `docs/repo/state/<id>/`

## Feedback

- Capture sanitized repo-specific feedback:
  `.codeheart/kit/docs/agent-interface/runbooks/capture-repo-feedback.md`
- Enable GitHub Issues feedback intake:
  `.codeheart/kit/docs/agent-interface/runbooks/enable-github-issues-feedback-intake.md`
- Repo feedback item format:
  `.codeheart/kit/docs/agent-interface/reference/repo-feedback-item-format.md`
- Submit sanitized Operating Kit feedback:
  `.codeheart/kit/docs/agent-interface/runbooks/submit-kit-feedback.md`
- Operating Kit feedback item format:
  `.codeheart/kit/docs/agent-interface/reference/kit-feedback-item-format.md`

## Runbook Authoring

- Runbook authoring standard:
  `.codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md`
- Runbook-to-script promotion standard:
  `.codeheart/kit/docs/agent-interface/reference/runbook-to-script-promotion-standard.md`
- Promote runbook recipe to script:
  `.codeheart/kit/docs/agent-interface/runbooks/promote-runbook-recipe-to-script.md`
- User-entered terminal prompt standard:
  `.codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md`

## Tooling Readiness

- Missing local tooling route:
  `.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`
- Consumer runtime materialization route:
  `.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`

## Operating Rule

Treat files under `.codeheart/kit/` as managed Operating Kit content. Repair or refresh them with
`codeheart-operating-kit sync` instead of hand-editing them.
