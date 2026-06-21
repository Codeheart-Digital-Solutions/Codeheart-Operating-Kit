Last updated: 2026-06-21T15:11:52Z (UTC)

# Root AGENTS.md Contract

Consumer root `AGENTS.md` is the lean entry point for agents.

## Layers

1. Operating Kit managed bootstrap block.
2. Repository-owned instructions.
3. Local user guidance.

The managed block stays short and contains only universal immediate rules plus direct routes into
`.codeheart/kit/docs/`.

## Direct Managed Routes

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
- Structure governance: `.codeheart/kit/docs/structure-governance/README.md`

Use `.codeheart/kit/README.md` as the fallback and full inventory for unclear tasks, kit
management, missing or stale kit state, and explicit user questions about kit structure.

The portfolio-coordination route is conditional. Agents follow it when a planning workflow says a
plan-register update is material and `.codeheart/kit.config.yaml` configures portfolio
coordination. The managed block must not hardcode coordination-home paths, private repository
names, or local machine paths.
