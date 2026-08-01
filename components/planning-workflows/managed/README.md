Last updated: 2026-07-31T22:23:19Z (UTC)

# Planning Workflows

This component owns managed guidance for discovery, implementation planning, implementation-plan
execution, planning document review, semantic plan catalogs, portfolio coordination, migration,
and the stable plan-register entry point.

## Use

- Use `runbooks/discovery-workflow.md` for unclear, early, cross-domain, or decision-heavy
  discovery work.
- Use `runbooks/draft-implementation-plan.md` to create execution-ready
  `*_implementation_doc.md` files.
- Use `runbooks/execute-implementation-plan.md` to execute active implementation plans,
  including goal-style runs.
- Use `runbooks/review-planning-document.md` to review discovery and implementation documents for
  quality and execution readiness.
- Use `runbooks/maintain-plan-register.md` to select legacy maintenance versus mixed/canonical
  derived-view behavior.
- Use `runbooks/configure-portfolio-coordination.md` to preview and configure a member or
  coordination home.
- Use `runbooks/refresh-portfolio-catalog.md` before current portfolio analysis.
- Use `runbooks/migrate-plan-catalog.md` for reviewed legacy-to-canonical adoption.
- Use `../agent-interface/reference/runbook-authoring-standard.md` when plans create or
  materially change durable runbooks.
- Use `../agent-interface/reference/operation-routing-and-dispatch.md` when plans create or
  materially change routing-bearing surfaces.
- Use `reference/planning-document-lifecycle.md` for planning metadata, statuses, execution logs,
  plan bundles, subplans, plan families, program folders, attachments, archives, and index
  maintenance.
- Use `reference/plan-catalog-format.md` for canonical metadata, semantic IDs, families, modes,
  derived views, source observations, and compatibility.
- Use `reference/portfolio-coordination-format.md` for roles, configuration, exact membership,
  discovery sources, cache completeness, and strategic overlays.
- Use `reference/plan-register-format.md` for the stable entry point, legacy evidence, frozen
  mixed-mode baseline, and no-manual-append behavior.

## Boundaries

Reusable planning doctrine belongs in this managed component. Consumer repositories own their local
plans, local execution logs, product-specific guidance, release evidence, migration state, and
business-specific planning records.

Do not copy consumer-private details into managed planning docs. Use generic placeholder paths and
public-safe examples only.
