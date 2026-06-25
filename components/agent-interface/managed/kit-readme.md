Last updated: 2026-06-25T13:05:46Z (UTC)

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

## Module And Extension State

- Committed non-secret module/extension state convention:
  `.codeheart/kit/docs/structure-governance/reference/module-extension-state.md`
- Consumer-owned module/extension state root:
  `docs/repo/state/<id>/`

## Feedback

- Submit sanitized Operating Kit feedback:
  `.codeheart/kit/docs/agent-interface/runbooks/submit-kit-feedback.md`
- Feedback item format:
  `.codeheart/kit/docs/agent-interface/reference/kit-feedback-item-format.md`

## Runbook Authoring

- Runbook authoring standard:
  `.codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md`

## Tooling Readiness

- Missing local tooling route:
  `.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`

## Operating Rule

Treat files under `.codeheart/kit/` as managed Operating Kit content. Repair or refresh them with
`codeheart-operating-kit sync` instead of hand-editing them.
