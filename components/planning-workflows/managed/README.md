Last updated: 2026-06-21T14:53:02Z (UTC)

# Planning Workflows

This component owns managed guidance for discovery, implementation planning, implementation-plan
execution, planning document review, planning document lifecycle, and plan-register maintenance.

## Use

- Use `runbooks/discovery-workflow.md` for unclear, early, cross-domain, or decision-heavy
  discovery work.
- Use `runbooks/draft-implementation-plan.md` to create execution-ready
  `*_implementation_doc.md` files.
- Use `runbooks/execute-implementation-plan.md` to execute active implementation plans,
  including goal-style runs.
- Use `runbooks/review-planning-document.md` to review discovery and implementation documents for
  quality and execution readiness.
- Use `runbooks/maintain-plan-register.md` to update local and configured coordination-home plan
  registers for material planning lifecycle and relationship changes.
- Use `reference/planning-document-lifecycle.md` for planning metadata, statuses, execution logs,
  plan bundles, subplans, plan families, program folders, attachments, archives, and index
  maintenance.
- Use `reference/plan-register-format.md` for `docs/repo/plans/plan-register.md` entry shape,
  lifecycle snapshots, relation vocabulary, session refs, and coordination notes.

## Boundaries

Reusable planning doctrine belongs in this managed component. Consumer repositories own their local
plans, local execution logs, product-specific guidance, release evidence, migration state, and
business-specific planning records.

Do not copy consumer-private details into managed planning docs. Use generic placeholder paths and
public-safe examples only.
