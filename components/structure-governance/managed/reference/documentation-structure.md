Last updated: 2026-06-25T13:05:46Z (UTC)

# Documentation Structure

This reference defines the reusable documentation organization model for Operating Kit consumers.

Documentation placement follows ownership first, then domain and artifact kind inside that
boundary.

## Principles

- Keep reusable generic operating doctrine in the Operating Kit.
- Keep repository-specific guidance in consumer repository docs.
- Keep product-owned guidance under the owning product, package, module, or source-area boundary.
- Keep curated agent memory state under the consumer memory area.
- Keep local user preferences in the local user layer.
- Long-form planning docs live under `plans/`.
- Operational procedures live under `runbooks/`.
- Stable contracts, naming rules, lifecycle definitions, schemas, and source-of-truth guidance
  live under `reference/`.
- Superseded or historical docs retained for traceability live under `archive/`.
- README files are routers. They should route, not duplicate durable doctrine.
- Parent routers should advertise capability families and route-registry pointers, not copy deep
  route cards or execution recipes.
- Do not create new durable folders until there is enough material to justify a stable owner and
  boundary.

## Documentation Authority Boundaries

Use clear ownership:

- Reference docs own durable rules, stable contracts, naming rules, lifecycle definitions,
  placement rules, and explanatory source-of-truth guidance.
- Runbooks own operational procedure: ordered steps, preflight checks, execution checks, stop
  conditions, validation, rollback, and handoff.
- README routers own discoverability.
- Plans own scoped decisions, strategy, tasks, and execution state for a specific change.
- Execution logs own goal-style implementation divergence, validation evidence, review-gate
  results, and residual risk.
- Capability advertisements own compact boundary discoverability for repeated operations.
- Route registries and route cards own dispatch contracts inside the owning domain.

Before adding durable guidance, check the nearest README and relevant `reference/` and `runbooks/`
folders for an existing owner. If a durable rule already has an owner, link to that owner instead
of repeating the full rule.

## Routing Artifact Placement

Use `../../agent-interface/reference/operation-routing-and-dispatch.md` for routing behavior,
authority hierarchy, capability advertisement fields, route registry expectations, route-card
fields, ambiguity handling, and fresh low-context routing probes.

Use this structure reference for placement:

- Root `AGENTS.md`: compact bootstrap router and pointer to the full routing reference. It should
  not carry full capability catalogs, route registries, or route cards.
- README routers: local discoverability, capability advertisements, and route-registry pointers
  for the owner boundary. They should not duplicate deep route cards or canonical recipes.
- Capability advertisements: in the owning domain's primary README or route-registry introduction
  when the owner handles repeated operations or deep routes that parent routers cannot enumerate.
- Route registries: in the owning domain's `reference/`, primary README, or equivalent durable
  docs area. Use a local convention that keeps the registry close to the owner.
- Route cards: in the owning route registry or adjacent durable reference. Put cards where the
  operational owner can maintain them.
- Canonical recipes: in `runbooks/`, scripts, API procedure docs, portal procedure docs, or other
  owner-approved procedure surfaces after routing is complete.
- Stable route behavior and field definitions: in managed Operating Kit references unless the
  rule is domain-specific.
- Committed non-secret routing state: under `docs/repo/state/<module-or-extension-id>/` only when
  the installed module or extension owns real repo state.
- Execution evidence and blockers: in execution logs, run records, reports, or final summaries,
  not in route-card definitions.
- Local wrappers: in the consumer-owned repository, product, module, package, or source-area docs
  where the local exception lives.

Do not create a central repository catalog of every route card. Upper layers should expose enough
capability information for agents to find the owner and then descend to the owner's registry.

## Generic Target Shape

```text
/
  AGENTS.md
  .codeheart/
    kit/
    kit.config.yaml
    kit.lock.yaml
    user/
  docs/
    README.md
    agent-memory/
      README.md
      goal-register.md
      session-ledger.md
      untriaged-sessions.md
      archive/
    repo/
      README.md
      plans/
      runbooks/
      reference/
      state/
    research/
      README.md
    business/
      README.md
  products/
    README.md
    <product-slug>/
      README.md
      docs/
        README.md
        plans/
        runbooks/
        reference/
        archive/
      <source-area>/
        README.md
        docs/
          README.md
          plans/
          runbooks/
          reference/
          archive/
      packages/
        README.md
        <package>/
          README.md
          docs/
```

Consumers may omit areas they do not need. Do not create business, research, product, package, or
source-area folders until the repository has a real owner and use for them.

## Top-Level Documentation Areas

- `docs/repo/`: consumer repository organization, local development conventions,
  repo-maintenance plans, local runbooks, local references, and local exceptions to Operating Kit
  defaults.
- `docs/repo/state/`: committed, non-secret, consumer-owned routing state for installed modules
  and extensions. Use `docs/repo/state/<module-or-extension-id>/` after a module or extension ID
  is known.
- `docs/agent-memory/`: curated agent memory state. This is demand-driven context recovery, not a
  default read for every task.
- `docs/business/`: business operating docs when the consumer repository intentionally stores
  business records.
- `docs/research/`: non-product research that is intentionally separate from product
  implementation authority and business operating docs.
- `products/`: product roots when the repository contains product-owned source or product-owned
  docs.

## Product Documentation Model

Use product roots when the repository contains product-owned source or product-owned docs.

- `products/<product-slug>/README.md`: first product-specific router.
- `products/<product-slug>/docs/README.md`: product-wide documentation router.
- `products/<product-slug>/<source-area>/README.md`: source-area router.
- `products/<product-slug>/<source-area>/docs/README.md`: source-area documentation router.
- `products/<product-slug>/packages/README.md`: package collection router.
- `products/<product-slug>/packages/<package>/README.md`: package router.

Generic docs should use placeholder paths such as `products/<product-slug>/docs/`. Product-specific
domain names belong in the consumer repository's owning product README tree.

## Type Folders

Use these consistently inside repo docs, product docs, source-area docs, package docs, and other
owned docs roots:

- `plans/`: discovery docs, implementation docs, migration plans, refactor plans, and other
  long-form change plans.
- `runbooks/`: step-by-step operational procedures.
- `reference/`: stable contracts, naming rules, lifecycle definitions, schemas, and durable
  architecture references.
- `state/`: committed non-secret routing state for installed modules or extensions under
  `docs/repo/state/<module-or-extension-id>/`.
- `archive/`: superseded plans and historical docs retained for traceability.

Package-local docs may use package-native conventions when an ecosystem has a strong standard.
Prefer a root `README.md` plus `docs/` for long-form package maintainer, API, schema, or generated
reference docs.

## Planning Doc Placement

Use the planning lifecycle reference for planning document shape. Placement still follows
ownership:

- repo-maintenance plans under `docs/repo/plans/`;
- product-wide plans under `products/<product-slug>/docs/plans/`;
- source-area plans under `products/<product-slug>/<source-area>/docs/plans/`;
- package-local plans under `products/<product-slug>/packages/<package>/docs/`;
- plan-scoped attachments under the plan bundle's `attachments/` folder.

## Planning Doc Conventions And Shapes

Discovery and implementation docs use lifecycle headers:

```text
Last updated: YYYY-MM-DDTHH:MM:SSZ (UTC)
Created: YYYY-MM-DD
Status: draft | active | completed | superseded | archived
```

Use repo-relative links and maintain a bottom `# Revision Notes` section for meaningful decision,
scope, strategy, or execution-plan changes. Goal-style implementation runs use a sibling
`*_execution_log.md`.

Use the smallest planning shape that preserves ownership, reviewability, and lifecycle clarity:

- Standalone plan file: one low-artifact plan directly under `plans/`.
- Plan bundle: a folder for one plan plus execution log, discovery doc, or attachments.
- Subplan: a child plan under a parent plan's `subplans/` folder when the parent defines the
  child's correctness or review context.
- Plan family: related sibling plan bundles with shared discoverability but independent execution.
- Program folder: a parent coordination plan that owns multiple workstream plans.

Use the planning lifecycle reference for the full lifecycle and shape rules. This structure
reference owns where those shapes belong.

## Consumer Repository Docs

Consumer `docs/repo/` is for local repository plans, runbooks, references, build/test/release
details, committed non-secret module or extension state, local architecture notes, and local
exceptions to Operating Kit defaults.

Generic operating doctrine belongs in the Operating Kit. If an agent or contributor proposes
reusable generic rules under consumer `docs/repo/`, flag that they are kit rules and recommend
changing the Operating Kit instead.

Committed state under `docs/repo/state/<module-or-extension-id>/` is routing context, not live
external truth. Agents must still run live preflight before sensitive reads, writes, permission
changes, or external system changes.

Consumer-owned route cards and route registries belong with the repository, product, module,
package, or source-area owner that maintains the operation. Keep generic route-card field
semantics in the managed Agent Interface reference and keep local route details with the local
owner.

## Index Maintenance

Update the nearest README and parent index when discoverability changes:

- a new discoverable doc, runbook, plan, reference, or execution log is added;
- a doc path, title, or purpose changes;
- a doc is moved, archived, or removed;
- a product, product area, package, or module folder is created;
- a committed module or extension state namespace is created under `docs/repo/state/`;
- a capability advertisement, route registry, or route card is added, removed, renamed, or
  materially changed;
- a runbook path, command path, validator input, workflow command, or entry point changes.

Index updates are not required for timestamp-only edits, local wording inside an already-linked
document, or checklist progress inside an existing plan.

## Retired And Local Paths

Do not preserve obsolete routers as current authority. Historical plans may mention old paths as
traceability records, but current docs should route to the current owner.

When a consumer-local document becomes redundant with managed doctrine, prefer a concise wrapper
with local exceptions over a duplicated copy.
