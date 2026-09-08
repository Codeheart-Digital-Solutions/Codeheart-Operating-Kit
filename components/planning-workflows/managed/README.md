Last updated: 2026-09-07T22:40:39Z (UTC)

# Planning Workflows

This component owns managed guidance for discovery, implementation planning, implementation-plan
execution, planning document review, semantic plan catalogs, portfolio coordination, migration,
and the stable plan-register entry point.

## Use

- Use `runbooks/handle-routine-change.md` for known bounded changes with existing-record evidence.
- Use `../agent-interface/reference/agent-task-coordination.md` to commission whole-plan execution,
  name Git/delivery boundaries and delegate epic review with direct report-back.

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
- Use `runbooks/migrate-plan-catalog.md` for reviewed legacy-to-canonical adoption, hash-bound
  branch-touch reconciliation, exact mixed-mode deferral, guarded activation, and mandatory
  incremental follow-up.
- Use `../agent-interface/reference/runbook-authoring-standard.md` when plans create or
  materially change durable runbooks.
- Use `../agent-interface/reference/operation-routing-and-dispatch.md` when plans create or
  materially change routing-bearing surfaces.
- Use `reference/planning-document-lifecycle.md` for planning metadata, statuses, execution logs,
  plan bundles, subplans, plan families, program folders, attachments, archives, and index
  maintenance.
- Use `reference/plan-catalog-format.md` for canonical metadata, semantic IDs, families, modes,
  repository-wide discovery, ownership/exclusions, reviewed ledger-v3 dispositions and proofs,
  derived views, source observations, readiness, and compatibility activation.
- Use `reference/portfolio-coordination-format.md` for roles, configuration, exact membership,
  discovery sources, remote-aware evidence, cache-v3 completeness/readiness, compatibility
  observations, offline limits, and strategic overlays.
- Use `reference/plan-register-format.md` for the stable entry point, legacy evidence, frozen
  mixed-mode baseline, and no-manual-append behavior.

Local commits preserve coherent completed work. Publication and normal PR integration follow the
covered delivery authority; a merged draft remains draft until execution is authorized. Use the
lifecycle reference for these defaults and execution guidance for coherent batches, useful review,
focused correction follow-up and evidence reuse.

## Boundaries

Reusable planning doctrine belongs in this managed component. Consumer repositories own their local
plans, local execution logs, product-specific guidance, release evidence, migration state, and
business-specific planning records.

Under plan-catalog discovery v2, those consumer-owned plans may live in any Git-tracked Markdown
documentation tree with an exact lowercase `docs` segment at arbitrary depth. Do not move an
otherwise truthful product, package, source-area, business, or domain plan into `docs/repo/plans/`
merely to make it catalog-visible. Existing repositories activate v2 only after prospective review;
a Kit upgrade alone preserves their current discovery contract.

Historical branches may be retained. A branch touch blocks unless current reviewed evidence proves
it is a same-content or incorporated-history non-owner, or mixed mode binds one genuine active
owner to an exact deferred-plan snapshot and pending incremental migration. Missing, stale,
ambiguous, or changed evidence always fails closed; direct canonical activation remains zero-gap.

Do not copy consumer-private details into managed planning docs. Use generic placeholder paths and
public-safe examples only.
