Last updated: 2026-07-02T13:16:41Z (UTC)

# Root AGENTS.md Contract

Consumer root `AGENTS.md` is the lean entry point for agents.

## Layers

1. Operating Kit managed bootstrap block.
2. Repository-owned instructions.
3. Local user guidance.

The managed block stays short and contains only universal immediate rules plus direct routes into
`.codeheart/kit/docs/`.

## Compact Routing Hierarchy

The managed block should expose the route-before-surface reflex without becoming a catalog of
deep routes.

Root wording should tell agents to route structural, external, sensitive, module, product, or
ambiguous work before selecting a tool, connector, API, browser, script, or runbook. It should
link to:

```text
.codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md
```

The detailed hierarchy, trigger categories, capability advertisements, route registries, route
cards, ambiguity handling, and probes belong in that managed Agent Interface reference. Root
`AGENTS.md` should not list every domain route, module capability, provider, execution surface, or
route-card field.

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
- Operation routing and dispatch:
  `.codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md`
- Repo feedback capture:
  `.codeheart/kit/docs/agent-interface/runbooks/capture-repo-feedback.md`
- Operating Kit feedback:
  `.codeheart/kit/docs/agent-interface/runbooks/submit-kit-feedback.md`
- Tooling readiness:
  `.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`
- Structure governance: `.codeheart/kit/docs/structure-governance/README.md`

Use `.codeheart/kit/README.md` as the fallback and full inventory for unclear tasks, kit
management, missing or stale kit state, and explicit user questions about kit structure.

The portfolio-coordination route is conditional. Agents follow it when a planning workflow says a
plan-register update is material and `.codeheart/kit.config.yaml` configures portfolio
coordination. The managed block must not hardcode coordination-home paths, private repository
names, or local machine paths.
