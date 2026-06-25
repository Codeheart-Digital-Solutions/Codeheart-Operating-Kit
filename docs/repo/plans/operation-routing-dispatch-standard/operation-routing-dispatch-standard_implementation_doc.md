Last updated: 2026-06-25T12:47:17Z (UTC)
Created: 2026-06-25
Status: draft

# Document Header

## Operation Routing And Dispatch Standard Implementation Plan

Overview: Implement the approved Operating Kit routing and dispatch standard so future agents
route repeated, structural, external, sensitive, module, product, or ambiguous work before
choosing tools, connectors, APIs, browser surfaces, scripts, or runbooks. The implementation
starts with a detailed provisional inventory of existing routing surfaces, then adds one managed
Agent Interface routing reference, compact root visibility, structure-governance placement rules,
runbook/planning hooks, a late final consolidation pass, fresh low-context routing probes,
packaged-resource mirrors, validation, and release preparation.

Execution dependency: this plan should execute after `OK-PR-012 - Discovery Handoff Gate
Implementation` completes. `OK-PR-012` is currently active and targets `v0.1.13`; this plan
therefore targets the next patch release, expected `v0.1.14`, unless execution finds that the
source release state has moved and records a consistent next patch target.

Essential context:

| Source | Why it matters |
| --- | --- |
| `docs/repo/plans/operation-routing-dispatch-standard/operation-routing-dispatch-standard_discovery_doc.md` | Approved discovery, implementation capability scope, decisions D-001 through D-015, requirements FR-001 through FR-016, and feature-level success evidence. |
| `docs/repo/plans/discovery-handoff-gate/discovery-handoff-gate_implementation_doc.md` | Active prerequisite that hardens discovery-to-implementation planning and currently owns the `v0.1.13` release lane. |
| `AGENTS.md` | Source-repository public-core safety, task routing, and release authority rules. |
| `README.md` | Public repository purpose, managed/consumer boundary, and maintainer entry points. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure before changing managed kit source, docs, templates, tests, or release assets. |
| `docs/repo/runbooks/release-operating-kit.md` | Required procedure for release validation, release assets, checksums, tag, GitHub release, and release evidence. |
| `docs/repo/reference/placement-contract.md` | Public placement contract for managed content, generated surfaces, and consumer-owned boundaries. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact classification and release-note requirements for managed instruction changes. |
| `components/agent-interface/managed/README.md` | Managed Agent Interface router that must expose the new routing reference. |
| `components/agent-interface/managed/kit-readme.md` | Installed fallback inventory that must route unclear tasks to the new standard. |
| `components/agent-interface/managed/reference/root-agents-md-contract.md` | Contract for keeping root `AGENTS.md` compact while exposing the routing hierarchy. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Existing durable-runbook standard that must receive route-card and routing-bearing runbook guidance. |
| `templates/agents/AGENTS.managed-block.md` | Source template for installed root `AGENTS.md` managed block and compact routing hierarchy. |
| `components/structure-governance/managed/README.md` | Managed Structure Governance router for placement and ownership guidance. |
| `components/structure-governance/managed/reference/documentation-structure.md` | Durable placement reference for routers, references, runbooks, plans, state, and indexes. |
| `components/structure-governance/managed/reference/module-extension-state.md` | Existing module/extension routing-state boundary that the new standard must preserve. |
| `components/structure-governance/managed/reference/managed-content-boundaries.md` | Managed/consumer ownership reference that should distinguish route artifacts from consumer state. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | Planning runbook that should lightly flag routing-bearing discovery scope. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Planning runbook that should require routing-bearing implementation plans to cite and apply the routing standard. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Execution runbook that should verify routing-bearing epic probes and review-gate evidence. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | Review runbook that should catch missing routing-standard adoption in plans. |
| `components/*/component.yaml` | Component manifest versions and managed file lists for changed components. |
| `src/codeheart_operating_kit/resources/` | Packaged resource mirror installed consumers receive. |
| `tests/test_packaging_resources.py`, `tests/test_onboard.py`, `tests/test_sync_check.py`, `tests/test_install_metadata.py`, `tests/test_release_assets.py` | Focused validation for packaged-resource parity, root block visibility, sync behavior, release metadata, and release assets. |
| `release-notes.md`, `pyproject.toml`, `src/codeheart_operating_kit/__init__.py`, `manifest.yaml`, `bootstrap.md`, `install.sh`, `install.ps1` | Release surfaces that must move consistently to the selected patch version during release preparation. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Implement the approved Operation Routing And Dispatch Standard capability in managed Operating
Kit source and packaged resources.

Completion is proven when:

- a detailed provisional existing-routing-surface inventory exists at
  `docs/repo/plans/operation-routing-dispatch-standard/attachments/routing-surface-inventory.md`;
- the inventory gives each relevant routing surface a provisional disposition before new doctrine
  is added;
- the finalized inventory records final dispositions and consolidation edits after the core
  doctrine, placement rules, and planning hooks exist;
- a fresh low-context routing probe matrix exists at
  `docs/repo/plans/operation-routing-dispatch-standard/attachments/fresh-low-context-routing-probes.md`;
- `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` exists as the
  managed owner of routing behavior, authority hierarchy, route-before-surface rules,
  trigger categories, ambiguity handling, conflict handling, capability advertisements, route
  registries, route cards, route-card field semantics, fresh low-context routing probes, and
  advertisement maintenance;
- root `AGENTS.md` managed block exposes only a compact routing hierarchy and pointer to the full
  routing reference;
- installed fallback inventory and Agent Interface routes expose the new reference;
- Structure Governance explains where routers, capability advertisements, route registries, route
  cards, references, runbooks, state, and local wrappers belong;
- durable runbook-authoring and planning workflows require routing-bearing work to apply the
  established standard instead of inventing local routing rules;
- materially changed planning runbooks have audience/intention coverage required by the runbook
  authoring standard;
- source managed files and packaged resource mirrors agree;
- focused tests and public-core validation pass;
- fresh low-context routing probes show that a new agent can discover the route, avoid choosing a
  visible execution surface prematurely, ask required ambiguity questions, and keep tiny local
  work lightweight;
- release notes classify the change as an `instruction-only change` with no automatic consumer
  migration or new consumer-owned route scaffolding.

## 1.2 Project And Problem Context

The motivating failure was structural. An agent can enter a repository, see available tools,
connectors, docs, scripts, browser surfaces, and modules, then choose an execution surface before
resolving the authoritative route. The Microsoft 365 discussion exposed this because a visible
mail connector and a module-owned workspace route could both appear relevant to a plain-language
request. The reusable issue is not email-specific and not connector-discoverability-specific. It
is a missing pre-execution routing doctrine.

The approved discovery assigns generic routing doctrine to Operating Kit. Domains, modules,
products, packages, and repo areas instantiate the doctrine through capability advertisements,
route registries, route cards, and recipes at their ownership boundary. Upper layers should route
and advertise capability families without becoming stale catalogs of every deep route.

The implementation must also respect the evolving state of Operating Kit. Existing managed docs
already contain partial routing guidance: root routes, module state routing, tooling readiness,
runbook authoring, planning workflows, structure governance, and managed-content boundaries. The
first epic therefore inventories those surfaces and assigns provisional dispositions before adding
the new reference. A later consolidation epic revisits that inventory after the new standard is
concrete and applies final linking, consolidation, and rewording.

## 1.3 Current State Analysis

Current source state:

- Root `AGENTS.md` managed block has immediate rules and direct managed routes, including module
  state routing and tooling readiness, but it does not yet expose a general route-before-surface
  hierarchy.
- Agent Interface has no dedicated operation-routing-and-dispatch reference.
- Runbook authoring already names source of truth, execution lane, preconditions, approval gates,
  stop conditions, evidence, and validation for agent-facing runbooks, but it does not yet define
  route-card adoption for repeated operational work.
- Structure Governance already distinguishes README routers, references, runbooks, plans, state,
  managed content, and local wrappers, but it does not yet place capability advertisements, route
  registries, or route cards explicitly.
- Module Extension State already says committed module state routes the agent while live
  preflight still decides before sensitive or external action.
- Tooling Readiness already separates generic local tooling blockers from module-owned service
  blockers.
- Planning workflow changes from `OK-PR-012` are active in the worktree and should complete before
  this plan executes.
- Package and release surfaces currently expose `0.1.12`, while the active handoff-gate plan
  targets `0.1.13`. This routing plan should target the next patch release after that work.

Target state:

- A fresh agent sees the routing hierarchy early, opens the full reference only when the task
  category requires routing, and can decide whether to use a route card, ask an ambiguity
  question, follow module state, use tooling readiness, or proceed with lightweight local work.
- Repeated or high-risk operations have a clear route artifact owner.
- Modules and domain owners can adopt the standard without copying Operating Kit doctrine into
  each module.
- Planning and review workflows prevent future routing-bearing implementation plans from adding
  local routing rules that bypass the standard.

## 1.4 Runbook Change Coverage

This implementation materially changes durable planning runbooks. The plan therefore applies
`components/agent-interface/managed/reference/runbook-authoring-standard.md` to those files.

Affected runbooks:

| Runbook | Audience class | Required change |
| --- | --- | --- |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | agent-facing | Add a lightweight routing-bearing scope flag and compact intention block when absent. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | agent-facing | Require routing-bearing plans to cite and apply the routing standard, and add compact intention block when absent. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | agent-facing | Require execution evidence for routing probes and routing-standard adoption checks, and add compact intention block when absent. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | agent-facing | Add review checks for missing routing-standard adoption and routing probe evidence, and add compact intention block when absent. |

These runbooks do not need human-facing flow, hybrid separation, or module-specific local tooling
guidance. They may use existing shell, Python, `uv`, pytest, and fresh-agent/subagent execution
capabilities during implementation. If those local tools are missing during execution, route that
blocker through `components/agent-interface/managed/runbooks/handle-tooling-readiness.md`.

Existing consumer-owned, module-owned, and unrelated runbooks are outside scope. Module-specific
operational recipes remain module-owned.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Sequence the implementation from provisional inventory to doctrine to routing visibility to
placement and workflow hooks, then run a final disposition and consolidation pass. Package and
validate after source behavior is coherent.

Expected source tree:

```text
Codeheart-Operating-Kit/
  docs/
    README.md                                                            # modify
    repo/
      README.md                                                          # modify
      plans/
        README.md                                                        # modify
        plan-register.md                                                 # modify
        operation-routing-dispatch-standard/
          operation-routing-dispatch-standard_discovery_doc.md           # existing
          operation-routing-dispatch-standard_implementation_doc.md       # create
          operation-routing-dispatch-standard_execution_log.md            # create at execution
          attachments/
            routing-surface-inventory.md                                 # create in EP-01
            fresh-low-context-routing-probes.md                           # create in EP-07
  components/
    agent-interface/
      component.yaml                                                     # modify
      managed/
        README.md                                                        # modify
        kit-readme.md                                                    # modify
        reference/
          operation-routing-and-dispatch.md                              # create
          root-agents-md-contract.md                                     # modify
          runbook-authoring-standard.md                                  # modify
    planning-workflows/
      component.yaml                                                     # modify
      managed/
        runbooks/
          discovery-workflow.md                                          # modify
          draft-implementation-plan.md                                   # modify
          execute-implementation-plan.md                                 # modify
          review-planning-document.md                                    # modify
    structure-governance/
      component.yaml                                                     # modify
      managed/
        README.md                                                        # modify
        reference/
          documentation-structure.md                                     # modify
          managed-content-boundaries.md                                  # modify
          module-extension-state.md                                      # modify
  templates/
    agents/
      AGENTS.managed-block.md                                            # modify
  src/
    codeheart_operating_kit/
      __init__.py                                                        # modify during release prep
      resources/
        manifest.yaml                                                    # modify during release prep
        components/
          agent-interface/                                               # mirror source changes
          planning-workflows/                                            # mirror source changes
          structure-governance/                                          # mirror source changes
        templates/
          agents/
            AGENTS.managed-block.md                                      # mirror source changes
  tests/
    test_packaging_resources.py                                          # modify
    test_onboard.py                                                      # modify
    test_sync_check.py                                                   # modify
    test_install_metadata.py                                             # release validation
    test_release_assets.py                                               # release validation
    fixtures/
      release-manifest.json                                              # modify during release prep
  bootstrap.md                                                           # modify during release prep
  install.sh                                                             # modify during release prep
  install.ps1                                                            # modify during release prep
  manifest.yaml                                                          # modify during release prep
  pyproject.toml                                                         # modify during release prep
  release-notes.md                                                       # modify during release prep
  scripts/build-release-assets.py                                        # modify during release prep
  dist/                                                                  # create release assets
```

## 2.2 Open Questions And Assumptions Requiring Clarification

`OQ-1` - What exact patch version should this plan target?

- `BLOCKER: no`
- `Affects: EP-00`, `EP-08`, `EP-09`
- Unlocks exact release metadata, release asset names, tag name, and consumer proof target.
- Recommended default: execute after `OK-PR-012` completes `v0.1.13`, then target `v0.1.14`.
  If the current package version differs at execution time, use the next patch version and record
  the version decision in the execution log before editing release surfaces.

`OQ-2` - Does plan activation include public release publication approval?

- `BLOCKER: no`
- `Affects: EP-09`
- Unlocks whether EP-09 can create the public tag and publish the GitHub release during the same
  execution run.
- Recommended default: do not assume release approval from this draft. When the user activates
  the plan, require the activation request or a separate user message to explicitly approve public
  release publication before creating the tag. Source implementation and release asset preparation
  may proceed without publication.

`OQ-3` - Should V1 include route-card examples beyond generic examples inside the reference?

- `BLOCKER: no`
- `Affects: EP-02`
- Unlocks the amount of example material in the new routing reference.
- Recommended default: include only public-safe generic examples in the reference. Defer
  domain-specific examples to domain owners.

`OQ-4` - Should route-card schemas be added now?

- `BLOCKER: no`
- `Affects: Section 4`
- Unlocks whether Markdown route cards become machine-readable in this release.
- Recommended default: no. V1 uses Markdown guidance only, with schema work deferred until
  repeated route-card usage proves a stable shape.

## 2.3 Architectural Decisions With Reasoning

`AD-1` - Start with a detailed provisional inventory artifact

1. Problem being solved: Operating Kit already has partial routing guidance, and adding a new
   standard without disposition could create conflicting authority.
2. Simplest working solution: make EP-01 produce a detailed inventory attachment with current
   routing role, authority level, overlap, risk, provisional disposition, reason, and
   implementation note.
3. What may change in 6-12 months: a future validator or route-surface index may replace manual
   inventory.
4. Rationale: this gives the doctrine-writing epics enough existing-surface context without
   pretending final wording can be known before the standard exists.
5. Alternatives considered: writing the new reference first was rejected because the discovery
   explicitly made inventory the first implementation epic.

`AD-2` - Put the durable routing doctrine in Agent Interface

1. Problem being solved: agents need behavioral guidance before execution surface selection.
2. Simplest working solution: create
   `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` as the single
   durable reference for generic routing behavior and route-card standard.
3. What may change in 6-12 months: route-card examples or schemas may become separate files when
   real usage grows.
4. Rationale: Agent Interface already owns how installed consumers introduce Operating Kit to
   agents and how agents find managed docs.
5. Alternatives considered: placing the full doctrine in Structure Governance was rejected
   because structure owns placement, not the whole pre-execution dispatch behavior.

`AD-3` - Keep root `AGENTS.md` compact

1. Problem being solved: agents need the route-before-surface reflex early, but root managed
   blocks must stay short.
2. Simplest working solution: add a compact routing hierarchy and one pointer to the full routing
   reference.
3. What may change in 6-12 months: root wording may be shortened after the hierarchy stabilizes.
4. Rationale: root `AGENTS.md` is the early agent interface, while detailed doctrine belongs in a
   managed reference.
5. Alternatives considered: adding full route-card guidance to root `AGENTS.md` was rejected
   because it would bloat the managed block and become stale.

`AD-4` - Use Markdown route cards in V1

1. Problem being solved: route cards need a recognizable standard without forcing schema work
   before actual usage proves stable.
2. Simplest working solution: define the field set in Markdown and allow `not applicable` with a
   reason.
3. What may change in 6-12 months: repeated route-card usage may justify schemas, linting, or
   manifest guidance.
4. Rationale: Markdown aligns with current managed docs and keeps V1 lightweight.
5. Alternatives considered: adding a machine-readable route-card schema now was rejected because
   it would overfit before adoption.

`AD-5` - Planning workflows enforce future routing-standard adoption

1. Problem being solved: future routing-bearing implementations could invent local routing rules
   instead of using the standard.
2. Simplest working solution: add a light discovery hook and stronger implementation-planning,
   execution, and review hooks.
3. What may change in 6-12 months: a structured planning validator may replace some checklist
   review.
4. Rationale: implementation planning is where routing-bearing changes become concrete, so it is
   the right enforcement point.
5. Alternatives considered: making discovery own detailed standard-adoption checks was rejected
   because discovery should flag scope, not become an implementation checklist.

`AD-6` - Treat this as an instruction-only release

1. Problem being solved: the change updates managed guidance and packaging, but does not alter
   CLI sync behavior, generated paths, validators, schemas, or consumer-owned content placement.
2. Simplest working solution: classify as `instruction-only change`, update release notes, and
   use normal sync adoption.
3. What may change in 6-12 months: route-card schemas or scaffolds could introduce a different
   consumer-impact class.
4. Rationale: V1 changes instructions and managed docs only.
5. Alternatives considered: adding automatic route-card scaffolding was rejected because the
   discovery explicitly excludes automatic consumer migration and new consumer-owned route/state
   scaffolding.

`AD-7` - Finalize dispositions after the standard is written

1. Problem being solved: exact linking, consolidation, and rewording decisions depend on the
   concrete routing reference, root wording, placement rules, and planning hooks.
2. Simplest working solution: add a late consolidation epic after EP-02 through EP-05 and before
   packaging.
3. What may change in 6-12 months: a future route-surface index may make the consolidation pass
   more mechanical.
4. Rationale: early inventory prevents blind doctrine writing, while late consolidation prevents
   premature final dispositions.
5. Alternatives considered: moving the entire inventory to the end was rejected because the core
   standard should be written with known existing overlaps and authority risks in view.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-00 | Execution starts from a clean dependency baseline after `OK-PR-012`. | S | none |
| EP-01 | Existing routing surfaces are inventoried with provisional dispositions. | M | EP-00 |
| EP-02 | Agent Interface owns the full routing and dispatch reference. | L | EP-01 |
| EP-03 | Root and installed fallback routing surfaces expose compact route-before-surface guidance. | M | EP-02 |
| EP-04 | Structure Governance owns placement for routing artifacts and wrappers. | M | EP-02 |
| EP-05 | Runbook authoring and planning workflows enforce routing-standard adoption. | L | EP-02, EP-04 |
| EP-06 | Final dispositions are applied across existing routing surfaces. | M | EP-03, EP-04, EP-05 |
| EP-07 | Packaged resources, tests, and fresh low-context probe evidence prove the managed standard. | L | EP-06 |
| EP-08 | Release metadata and assets are prepared for the selected patch version. | M | EP-07 |
| EP-09 | Public release publication and first-consumer proof complete after explicit approval. | M | EP-08 |

## EP-00 - Dependency Baseline And Execution Setup

### A) Epic ID, Title, And Outcome

EP-00 - Dependency Baseline And Execution Setup

Outcome: Execution begins only after `OK-PR-012` has completed, the patch target is known, and the
execution log is ready.

### B) Scope

Establish the source baseline. Do not edit routing surfaces in this epic.

### C) Files Touched

```text
docs/repo/plans/operation-routing-dispatch-standard/
  operation-routing-dispatch-standard_execution_log.md              # create at execution
docs/repo/plans/plan-register.md                                    # modify lifecycle at activation
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- `OK-PR-012` is completed before routing implementation starts.
- Current package version is known.
- Target patch version is recorded in the execution log.
- Plan lifecycle and plan register reflect activation before source changes begin.
- No routing managed surface changes occur in this epic.

### E) Dependencies And Critical-Path Notes

Execution should stop when `OK-PR-012` is still active. This prevents concurrent edits to the same
planning workflow files and release surfaces.

### F) Tasks Checklist

- [ ] Confirm `OK-PR-012 - Discovery Handoff Gate Implementation` is completed in `docs/repo/plans/plan-register.md`.
- [ ] Confirm source package version with `PYTHONPATH=src python -c "import codeheart_operating_kit as c; print(c.__version__)"`.
- [ ] Record target patch version in `docs/repo/plans/operation-routing-dispatch-standard/operation-routing-dispatch-standard_execution_log.md`.
- [ ] Create `docs/repo/plans/operation-routing-dispatch-standard/operation-routing-dispatch-standard_execution_log.md`.
- [ ] Mark this implementation plan active.
- [ ] Update `docs/repo/plans/plan-register.md` lifecycle state for `OK-PR-013`.
- [ ] Re-read `docs/repo/runbooks/change-operating-kit.md`.
- [ ] Re-read `docs/repo/reference/consumer-impact-classification.md`.
- [ ] Re-read `components/agent-interface/managed/reference/runbook-authoring-standard.md`.

### G) Implementation Notes

The expected target after `OK-PR-012` is `0.1.14`. If the source package version has moved, record
the next patch target before EP-08 and use that value consistently.

### H) Open Questions

OQ-1 applies.

## EP-01 - Provisional Existing Routing Surface Inventory

### A) Epic ID, Title, And Outcome

EP-01 - Provisional Existing Routing Surface Inventory

Outcome: The implementation has a detailed inventory and provisional disposition map for existing
Operating Kit routing and instruction surfaces before adding new routing doctrine.

### B) Scope

Create a plan-scoped inventory artifact. Include managed source files and repository governance
files that already route agents, define authority, place docs, manage state, select tools, define
runbook behavior, or guide implementation/review. This epic proposes dispositions and candidate
source edits; EP-06 finalizes and applies dispositions after the core standard exists.

### C) Files Touched

```text
docs/repo/plans/operation-routing-dispatch-standard/attachments/
  routing-surface-inventory.md                                      # create
docs/repo/plans/operation-routing-dispatch-standard/
  operation-routing-dispatch-standard_execution_log.md              # update at execution
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Inventory artifact uses this row shape: `Surface`, `Current routing role`, `Instruction type`,
  `Authority level`, `Overlaps with`, `Conflict risk`, `Provisional disposition`, `Reason`,
  `Implementation note`.
- Inventory covers every direct managed route in `templates/agents/AGENTS.managed-block.md`.
- Inventory covers Agent Interface, Structure Governance, Planning Workflows, Module Extension
  State, Tooling Readiness, Runbook Authoring, Native Codex Capabilities, and managed fallback
  inventory surfaces.
- Inventory covers source-repository governance surfaces that route maintainers and agents before
  managed content is installed.
- Every inventory row has one provisional disposition from the discovery-approved set.
- The inventory produces the candidate source-edit list used by EP-02 through EP-05.
- EP-02 through EP-05 do not begin until every inventoried surface has a provisional disposition
  and the inventory has a `Candidate Implementation Edit List`.
- Inventory contains no private consumer paths, tenant details, customer details, or raw local
  logs.

### E) Dependencies And Critical-Path Notes

Depends on EP-00. Later epics should follow the inventory dispositions rather than re-deciding
ownership.

### F) Tasks Checklist

- [ ] Create `docs/repo/plans/operation-routing-dispatch-standard/attachments/routing-surface-inventory.md`.
- [ ] Add the approved provisional inventory row shape to `docs/repo/plans/operation-routing-dispatch-standard/attachments/routing-surface-inventory.md`.
- [ ] Inventory `templates/agents/AGENTS.managed-block.md`.
- [ ] Inventory `components/agent-interface/managed/kit-readme.md`.
- [ ] Inventory `components/agent-interface/managed/README.md`.
- [ ] Inventory `components/agent-interface/managed/reference/root-agents-md-contract.md`.
- [ ] Inventory `components/agent-interface/managed/reference/runbook-authoring-standard.md`.
- [ ] Inventory `components/agent-interface/managed/runbooks/handle-tooling-readiness.md`.
- [ ] Inventory `components/structure-governance/managed/README.md`.
- [ ] Inventory `components/structure-governance/managed/reference/documentation-structure.md`.
- [ ] Inventory `components/structure-governance/managed/reference/managed-content-boundaries.md`.
- [ ] Inventory `components/structure-governance/managed/reference/module-extension-state.md`.
- [ ] Inventory `components/planning-workflows/managed/README.md`.
- [ ] Inventory `components/planning-workflows/managed/runbooks/discovery-workflow.md`.
- [ ] Inventory `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [ ] Inventory `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`.
- [ ] Inventory `components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [ ] Inventory `components/native-codex-capabilities/managed/README.md`.
- [ ] Inventory `AGENTS.md`.
- [ ] Inventory `README.md`.
- [ ] Inventory `docs/README.md`.
- [ ] Inventory `docs/repo/README.md`.
- [ ] Inventory `docs/repo/reference/placement-contract.md`.
- [ ] Inventory `docs/repo/reference/consumer-impact-classification.md`.
- [ ] Inventory `docs/repo/runbooks/change-operating-kit.md`.
- [ ] Inventory `docs/repo/runbooks/release-operating-kit.md`.
- [ ] Inventory `docs/repo/runbooks/promote-consumer-change.md`.
- [ ] Inventory `docs/repo/runbooks/triage-kit-feedback.md`.
- [ ] Add provisional disposition values for every inventory row.
- [ ] Add a `Candidate Implementation Edit List` section to `docs/repo/plans/operation-routing-dispatch-standard/attachments/routing-surface-inventory.md`.
- [ ] Confirm `docs/repo/plans/operation-routing-dispatch-standard/attachments/routing-surface-inventory.md` maps every provisional disposition to candidate EP-02 through EP-05 source edits.
- [ ] Run `python3 scripts/validate-markdown-headers.py`.
- [ ] Run `python3 scripts/validate-public-core.py`.
- [ ] Run `git diff --check`.

### G) Implementation Notes

Use `rg -n "route|routing|authority|surface|tool|connector|module|extension|state|preflight|owner|capability|runbook|approval" components templates docs/repo` to find additional candidate rows, then include only reusable public-core-safe routing surfaces.

### H) Open Questions

None.

## EP-02 - Agent Interface Routing Reference

### A) Epic ID, Title, And Outcome

EP-02 - Agent Interface Routing Reference

Outcome: Agent Interface owns the full generic operation-routing and dispatch doctrine in one
managed reference.

### B) Scope

Create the new reference and update Agent Interface routing surfaces. Do not add domain-specific
route cards, schemas, validators, or module manifest rules.

### C) Files Touched

```text
components/agent-interface/component.yaml                           # modify
components/agent-interface/managed/README.md                        # modify
components/agent-interface/managed/reference/
  operation-routing-and-dispatch.md                                 # create
src/codeheart_operating_kit/resources/components/agent-interface/
  component.yaml                                                    # modify mirror in EP-07
  managed/README.md                                                 # modify mirror in EP-07
  managed/reference/operation-routing-and-dispatch.md               # create mirror in EP-07
```

### D) Acceptance Criteria And Size

Size: L

Acceptance criteria:

- New reference defines the route-before-surface rule.
- New reference defines the dispatch sequence: intent, domain, authority source, scope,
  capability route, execution surface, preconditions, canonical recipe, evidence or blocker.
- New reference defines mandatory, conditional, and lightweight routing trigger categories.
- New reference contains the approved routing authority hierarchy.
- New reference defines ambiguity handling and the approved conflict rule.
- New reference defines capability advertisements with required fields: owner, domain,
  capability families, intent aliases, route registry, ambiguity rule.
- New reference defines optional advertisement fields: scope families and state location.
- New reference defines route registries and route cards as owner-domain dispatch contracts.
- New reference defines the route-card field set and `not applicable` usage.
- New reference defines advertisement maintenance requirements.
- New reference defines fresh low-context routing probe shape and pass criteria.
- New reference includes public-safe generic examples only.
- Agent Interface README links to the new reference.
- Agent Interface component manifest includes the new managed file and version target.

### E) Dependencies And Critical-Path Notes

Depends on EP-01. Follow inventory dispositions when deciding which Agent Interface surfaces link
to the new reference.

### F) Tasks Checklist

- [ ] Create `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add the route-before-surface rule to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add the dispatch sequence to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add routing trigger categories to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add the approved authority hierarchy to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add ambiguity handling to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add authority conflict handling to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add capability advertisement guidance to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add route registry guidance to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add route-card guidance to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add advertisement maintenance guidance to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add fresh low-context routing probe guidance to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add public-safe generic examples to `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Add the new reference route to `components/agent-interface/managed/README.md`.
- [ ] Add the new managed file entry to `components/agent-interface/component.yaml`.
- [ ] Run fresh low-context routing probe `P-02-reference` against the new reference and record evidence in the execution log.
- [ ] Run `python3 scripts/validate-markdown-headers.py`.
- [ ] Run `python3 scripts/validate-public-core.py`.
- [ ] Run `git diff --check`.

### G) Implementation Notes

Use the discovery language as the source of truth, but refine wording where the EP-01 inventory
shows concrete overlap. Probe `P-02-reference` prompt: "A repository has a visible connector and a
deeper workspace module that may both handle a plain-language communication-resource request.
What should route the work before a tool is selected?" Pass when the fresh agent identifies the
root hierarchy, the routing reference, capability advertisement or ambiguity rule, and avoids
selecting the connector first.

### H) Open Questions

OQ-3 and OQ-4 apply.

## EP-03 - Root And Fallback Routing Visibility

### A) Epic ID, Title, And Outcome

EP-03 - Root And Fallback Routing Visibility

Outcome: Agents see compact route-before-surface guidance early through root `AGENTS.md` and can
find the full reference through installed fallback routes.

### B) Scope

Update root managed block wording, root contract, installed fallback inventory, and focused tests
that verify generated root block route visibility.

### C) Files Touched

```text
templates/agents/AGENTS.managed-block.md                            # modify
components/agent-interface/managed/kit-readme.md                    # modify
components/agent-interface/managed/reference/root-agents-md-contract.md # modify
tests/test_onboard.py                                               # modify
tests/test_sync_check.py                                            # modify
src/codeheart_operating_kit/resources/templates/agents/
  AGENTS.managed-block.md                                           # modify mirror in EP-07
src/codeheart_operating_kit/resources/components/agent-interface/
  managed/kit-readme.md                                             # modify mirror in EP-07
  managed/reference/root-agents-md-contract.md                      # modify mirror in EP-07
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Root managed block contains a compact routing hierarchy and pointer to the full routing
  reference.
- Root managed block stays concise and does not list every domain route or execution surface.
- Root contract explains the compact hierarchy and anti-catalog boundary.
- Installed fallback inventory routes unclear routing tasks to the new reference.
- Tests cover the new managed block route visibility after onboarding or sync.

### E) Dependencies And Critical-Path Notes

Depends on EP-02. Root wording must use the shortest safe form from the discovery and avoid
domain-specific examples.

### F) Tasks Checklist

- [ ] Add compact routing hierarchy wording to `templates/agents/AGENTS.managed-block.md`.
- [ ] Add the full routing reference path to `templates/agents/AGENTS.managed-block.md`.
- [ ] Update `components/agent-interface/managed/reference/root-agents-md-contract.md` with the compact hierarchy contract.
- [ ] Update `components/agent-interface/managed/kit-readme.md` with the operation-routing reference route.
- [ ] Update `tests/test_onboard.py` to assert root managed block operation-routing visibility.
- [ ] Update `tests/test_sync_check.py` to assert managed sync preserves operation-routing visibility.
- [ ] Run fresh low-context routing probe `P-03-root` against the root managed block and record evidence in the execution log.
- [ ] Run `python3 scripts/validate-markdown-headers.py`.
- [ ] Run `python3 scripts/validate-public-core.py`.
- [ ] Run `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_onboard.py tests/test_sync_check.py`.
- [ ] Run `git diff --check`.

### G) Implementation Notes

Probe `P-03-root` prompt: "Starting only from root `AGENTS.md`, route a vague request that may
belong to a deeper installed module while a visible execution tool also appears relevant." Pass
when the fresh agent finds the compact hierarchy, opens the full reference, and avoids choosing a
tool first.

### H) Open Questions

None.

## EP-04 - Structure Governance Placement Rules

### A) Epic ID, Title, And Outcome

EP-04 - Structure Governance Placement Rules

Outcome: Structure Governance defines where routing artifacts live and how parent routers avoid
duplicating deep route details.

### B) Scope

Update managed placement and boundary references. Do not scaffold consumer route folders or
define module-specific route-card content.

### C) Files Touched

```text
components/structure-governance/component.yaml                      # modify
components/structure-governance/managed/README.md                   # modify
components/structure-governance/managed/reference/documentation-structure.md # modify
components/structure-governance/managed/reference/managed-content-boundaries.md # modify
components/structure-governance/managed/reference/module-extension-state.md # modify
src/codeheart_operating_kit/resources/components/structure-governance/
  component.yaml                                                    # modify mirror in EP-07
  managed/README.md                                                 # modify mirror in EP-07
  managed/reference/documentation-structure.md                      # modify mirror in EP-07
  managed/reference/managed-content-boundaries.md                   # modify mirror in EP-07
  managed/reference/module-extension-state.md                       # modify mirror in EP-07
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Structure Governance distinguishes routers, capability advertisements, route registries, route
  cards, canonical recipes, references, runbooks, state, and local wrappers.
- Documentation Structure states that README files and root `AGENTS.md` route without duplicating
  full route cards.
- Module Extension State remains a routing-context rule and not a live-truth authorization rule.
- Managed Content Boundaries keeps route artifacts separate from consumer-owned state and
  generated evidence.
- Structure Governance README routes placement questions for routing artifacts.
- No new consumer-owned route scaffold is added.

### E) Dependencies And Critical-Path Notes

Depends on EP-02. Keep Agent Interface as routing behavior owner and Structure Governance as
placement owner.

### F) Tasks Checklist

- [ ] Update `components/structure-governance/managed/reference/documentation-structure.md` with routing artifact placement rules.
- [ ] Update `components/structure-governance/managed/reference/documentation-structure.md` with parent-router anti-catalog guidance.
- [ ] Update `components/structure-governance/managed/reference/managed-content-boundaries.md` with route artifact ownership boundaries.
- [ ] Update `components/structure-governance/managed/reference/module-extension-state.md` with the relationship between module state and route cards.
- [ ] Update `components/structure-governance/managed/README.md` with a route to the Agent Interface routing reference.
- [ ] Update `components/structure-governance/component.yaml` version target.
- [ ] Run fresh low-context routing probe `P-04-placement` and record evidence in the execution log.
- [ ] Run `python3 scripts/validate-markdown-headers.py`.
- [ ] Run `python3 scripts/validate-public-core.py`.
- [ ] Run `git diff --check`.

### G) Implementation Notes

Probe `P-04-placement` prompt: "A team wants to add a route card and a capability advertisement
for a repeated product operation. Where should the artifacts live, and what should the parent
README contain?" Pass when the fresh agent routes to Structure Governance for placement and Agent
Interface for route-card behavior.

### H) Open Questions

None.

## EP-05 - Runbook Authoring And Planning Workflow Hooks

### A) Epic ID, Title, And Outcome

EP-05 - Runbook Authoring And Planning Workflow Hooks

Outcome: Future routing-bearing runbooks and implementation plans are required to apply the
established routing standard and prove adoption through fresh low-context routing probes.

### B) Scope

Update the runbook authoring standard and the four planning workflow runbooks named in Section
1.4. Do not retrofit unrelated runbooks or add route-card schemas.

### C) Files Touched

```text
components/agent-interface/managed/reference/runbook-authoring-standard.md # modify
components/planning-workflows/component.yaml                      # modify
components/planning-workflows/managed/runbooks/discovery-workflow.md # modify
components/planning-workflows/managed/runbooks/draft-implementation-plan.md # modify
components/planning-workflows/managed/runbooks/execute-implementation-plan.md # modify
components/planning-workflows/managed/runbooks/review-planning-document.md # modify
src/codeheart_operating_kit/resources/components/agent-interface/
  managed/reference/runbook-authoring-standard.md                 # modify mirror in EP-07
src/codeheart_operating_kit/resources/components/planning-workflows/
  component.yaml                                                   # modify mirror in EP-07
  managed/runbooks/discovery-workflow.md                           # modify mirror in EP-07
  managed/runbooks/draft-implementation-plan.md                    # modify mirror in EP-07
  managed/runbooks/execute-implementation-plan.md                  # modify mirror in EP-07
  managed/runbooks/review-planning-document.md                     # modify mirror in EP-07
```

### D) Acceptance Criteria And Size

Size: L

Acceptance criteria:

- Runbook authoring standard tells new or materially changed agent-facing and hybrid runbooks to
  expose their routing contract or point to a route card when they handle repeated routing-bearing
  operations.
- Runbook authoring standard keeps recipe execution detail separate from route selection.
- Discovery workflow lightly flags routing-bearing scope.
- Implementation planning requires routing-bearing epics to cite and apply the established
  routing standard.
- Implementation planning requires fresh low-context routing probes for routing-bearing epics.
- Execution workflow verifies probe evidence before marking routing-bearing epics complete.
- Planning review checks for missing routing-standard adoption and missing probe evidence.
- Affected planning runbooks have compact intention blocks and audience class coverage.
- Generic local tooling blockers still route through Tooling Readiness.

### E) Dependencies And Critical-Path Notes

Depends on EP-02 and EP-04. This epic changes files currently touched by `OK-PR-012`, so EP-00
must verify that work is complete first.

### F) Tasks Checklist

- [ ] Update `components/agent-interface/managed/reference/runbook-authoring-standard.md` with routing-bearing runbook guidance.
- [ ] Update `components/agent-interface/managed/reference/runbook-authoring-standard.md` with route-card pointer guidance.
- [ ] Add compact intention block to `components/planning-workflows/managed/runbooks/discovery-workflow.md`.
- [ ] Add routing-bearing scope flag guidance to `components/planning-workflows/managed/runbooks/discovery-workflow.md`.
- [ ] Add compact intention block to `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [ ] Add routing-standard dependency guidance to `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [ ] Add fresh low-context routing probe planning guidance to `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [ ] Add compact intention block to `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`.
- [ ] Add routing probe evidence guidance to `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`.
- [ ] Add compact intention block to `components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [ ] Add routing-standard adoption review checks to `components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [ ] Update `components/planning-workflows/component.yaml` version target.
- [ ] Run fresh low-context routing probe `P-05-future-plan` and record evidence in the execution log.
- [ ] Run `python3 scripts/validate-markdown-headers.py`.
- [ ] Run `python3 scripts/validate-public-core.py`.
- [ ] Run `git diff --check`.

### G) Implementation Notes

Probe `P-05-future-plan` prompt: "Draft the validation approach for an implementation epic that
adds a new route registry and route cards for a product-owned operational workflow." Pass when the
fresh agent names the routing standard as a dependency and includes a fresh low-context routing
probe.

### H) Open Questions

None.

## EP-06 - Final Disposition And Consolidation

### A) Epic ID, Title, And Outcome

EP-06 - Final Disposition And Consolidation

Outcome: Existing routing surfaces are revisited against the written standard, final dispositions
are recorded, and necessary linking, consolidation, and rewording edits are applied.

### B) Scope

Finalize the provisional inventory from EP-01 after EP-02 through EP-05 have made the standard
concrete. Apply source edits that remove duplicate generic routing authority, convert redundant
local text into pointers, preserve legitimate domain-specific detail, and keep root surfaces
compact.

### C) Files Touched

```text
docs/repo/plans/operation-routing-dispatch-standard/attachments/
  routing-surface-inventory.md                                      # update
docs/repo/plans/operation-routing-dispatch-standard/
  operation-routing-dispatch-standard_execution_log.md              # update at execution
<inventoried routing surfaces from EP-01>                           # modify as listed by inventory
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Every inventory row has a final disposition.
- Inventory includes a `Final Consolidation Map` section.
- Final consolidation edits are applied to every surface marked `link-to-new-standard`,
  `consolidate-into-new-standard`, `keep-as-domain-specific-detail`, `retire-or-reword`, or
  `defer`.
- Generic routing behavior has one durable owner in the Agent Interface routing reference.
- Structure Governance owns placement wording without restating full routing doctrine.
- Root and README router surfaces stay compact and avoid full route-card duplication.
- The final inventory explains any surface left unchanged.
- Fresh low-context routing probe `P-06-consolidation` passes.

### E) Dependencies And Critical-Path Notes

Depends on EP-03, EP-04, and EP-05. This epic is the last source-behavior pass before packaging
and broader probe validation.

### F) Tasks Checklist

- [ ] Re-read `docs/repo/plans/operation-routing-dispatch-standard/attachments/routing-surface-inventory.md`.
- [ ] Add final disposition values for every inventory row.
- [ ] Add `Final Consolidation Map` to `docs/repo/plans/operation-routing-dispatch-standard/attachments/routing-surface-inventory.md`.
- [ ] Apply final `link-to-new-standard` edits from `Final Consolidation Map`.
- [ ] Apply final `consolidate-into-new-standard` edits from `Final Consolidation Map`.
- [ ] Apply final `keep-as-domain-specific-detail` edits from `Final Consolidation Map`.
- [ ] Apply final retirement and rewording edits from `Final Consolidation Map`.
- [ ] Record final `defer` reasons in `Final Consolidation Map`.
- [ ] Confirm root router surfaces stay compact after consolidation.
- [ ] Confirm generic routing doctrine has one durable owner in `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Run fresh low-context routing probe `P-06-consolidation` and record evidence in the execution log.
- [ ] Run `python3 scripts/validate-markdown-headers.py`.
- [ ] Run `python3 scripts/validate-public-core.py`.
- [ ] Run `git diff --check`.

### G) Implementation Notes

Probe `P-06-consolidation` prompt: "Find the durable owner of generic routing doctrine and decide
whether a nearby README should restate the rule, link to the standard, or keep only local
exceptions." Pass when the fresh agent finds the Agent Interface routing reference as behavior
owner, Structure Governance as placement owner, and avoids treating parent routers as full
route-card catalogs.

### H) Open Questions

None.

## EP-07 - Packaged Mirrors, Tests, And Probe Suite

### A) Epic ID, Title, And Outcome

EP-07 - Packaged Mirrors, Tests, And Probe Suite

Outcome: Source managed docs are mirrored into packaged resources, tests cover the changed
managed surfaces, and fresh low-context probes prove the feature-level routing behavior.

### B) Scope

Mirror changed managed files under `src/codeheart_operating_kit/resources/`, update focused
tests, run validation, and record probe evidence. Do not prepare release metadata in this epic.

### C) Files Touched

```text
src/codeheart_operating_kit/resources/components/agent-interface/   # modify mirrors
src/codeheart_operating_kit/resources/components/planning-workflows/ # modify mirrors
src/codeheart_operating_kit/resources/components/structure-governance/ # modify mirrors
src/codeheart_operating_kit/resources/templates/agents/AGENTS.managed-block.md # modify mirror
tests/test_packaging_resources.py                                  # modify
tests/test_onboard.py                                              # modify from EP-03
tests/test_sync_check.py                                           # modify from EP-03
docs/repo/plans/operation-routing-dispatch-standard/
  attachments/fresh-low-context-routing-probes.md                    # create
  operation-routing-dispatch-standard_execution_log.md              # update at execution
```

### D) Acceptance Criteria And Size

Size: L

Acceptance criteria:

- Every changed managed source file has a matching packaged resource mirror.
- `tests/test_packaging_resources.py` covers the new operation-routing reference.
- Focused tests cover root managed block routing visibility.
- Probe matrix uses this row shape: `Probe ID`, `Scenario`, `Fresh-agent prompt`,
  `Expected owner`, `Expected route artifacts`, `Required ambiguity question`,
  `Execution-surface anti-pattern to avoid`, `Pass criteria`, `Evidence fields`, `Result`.
- Fresh low-context probe suite covers deep capability, provider ambiguity, structure placement,
  tooling readiness, module state with live preflight, and lightweight local work.
- EP-07 is not complete until fresh low-context subagent execution is available and every probe
  result is recorded.
- Validation passes for Markdown headers, public-core hygiene, JSON schemas, packaged-resource
  parity, focused tests, and whitespace.

### E) Dependencies And Critical-Path Notes

Depends on EP-03, EP-04, and EP-05. This is the consolidation and proof epic before release
preparation.

### F) Tasks Checklist

- [ ] Mirror changed Agent Interface source files into `src/codeheart_operating_kit/resources/components/agent-interface/`.
- [ ] Mirror changed Planning Workflows source files into `src/codeheart_operating_kit/resources/components/planning-workflows/`.
- [ ] Mirror changed Structure Governance source files into `src/codeheart_operating_kit/resources/components/structure-governance/`.
- [ ] Mirror changed root managed block template into `src/codeheart_operating_kit/resources/templates/agents/AGENTS.managed-block.md`.
- [ ] Update `tests/test_packaging_resources.py` to include `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- [ ] Create `docs/repo/plans/operation-routing-dispatch-standard/attachments/fresh-low-context-routing-probes.md`.
- [ ] Add the approved probe row shape to `docs/repo/plans/operation-routing-dispatch-standard/attachments/fresh-low-context-routing-probes.md`.
- [ ] Add probe matrix row `P-07-deep-capability` to `docs/repo/plans/operation-routing-dispatch-standard/attachments/fresh-low-context-routing-probes.md`.
- [ ] Add probe matrix row `P-07-provider-ambiguity` to `docs/repo/plans/operation-routing-dispatch-standard/attachments/fresh-low-context-routing-probes.md`.
- [ ] Add probe matrix row `P-07-structure-placement` to `docs/repo/plans/operation-routing-dispatch-standard/attachments/fresh-low-context-routing-probes.md`.
- [ ] Add probe matrix row `P-07-tooling-readiness` to `docs/repo/plans/operation-routing-dispatch-standard/attachments/fresh-low-context-routing-probes.md`.
- [ ] Add probe matrix row `P-07-module-state-live-preflight` to `docs/repo/plans/operation-routing-dispatch-standard/attachments/fresh-low-context-routing-probes.md`.
- [ ] Add probe matrix row `P-07-lightweight-local-work` to `docs/repo/plans/operation-routing-dispatch-standard/attachments/fresh-low-context-routing-probes.md`.
- [ ] Run fresh low-context routing probe `P-07-deep-capability` and record evidence in the execution log.
- [ ] Run fresh low-context routing probe `P-07-provider-ambiguity` and record evidence in the execution log.
- [ ] Run fresh low-context routing probe `P-07-structure-placement` and record evidence in the execution log.
- [ ] Run fresh low-context routing probe `P-07-tooling-readiness` and record evidence in the execution log.
- [ ] Run fresh low-context routing probe `P-07-module-state-live-preflight` and record evidence in the execution log.
- [ ] Run fresh low-context routing probe `P-07-lightweight-local-work` and record evidence in the execution log.
- [ ] Run `python3 scripts/validate-markdown-headers.py`.
- [ ] Run `python3 scripts/validate-public-core.py`.
- [ ] Run `python3 scripts/validate-json-schemas.py`.
- [ ] Run `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_packaging_resources.py tests/test_onboard.py tests/test_sync_check.py tests/test_install_metadata.py tests/test_release_assets.py`.
- [ ] Run `git diff --check`.

### G) Implementation Notes

Use fresh subagents with low historical context and the lowest reasoning profile available. Record
prompt, expected route, observed route, pass/fail result, and fixes made. Keep probe prompts
public-safe and generic.

If fresh low-context subagent execution is unavailable, record the blocker and leave EP-07
incomplete. The routing capability depends on fresh-agent evidence.

### H) Open Questions

None.

## EP-08 - Release Preparation

### A) Epic ID, Title, And Outcome

EP-08 - Release Preparation

Outcome: Release metadata, release notes, fixtures, and local assets are prepared for the selected
patch version.

### B) Scope

Update versioned release surfaces and build release assets. Do not create a public tag in this
epic.

### C) Files Touched

```text
release-notes.md
pyproject.toml
src/codeheart_operating_kit/__init__.py
scripts/build-release-assets.py
bootstrap.md
install.sh
install.ps1
manifest.yaml
src/codeheart_operating_kit/resources/manifest.yaml
tests/fixtures/release-manifest.json
dist/
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Release notes describe the routing standard as an instruction-only managed-doc release.
- Version surfaces consistently target the selected patch version.
- Root release manifest records publishable asset checksums.
- Packaged release manifest follows the established placeholder checksum strategy for bundled
  downloadable asset metadata.
- `dist/` contains macOS and Windows release assets plus checksum files.
- Full pytest passes after release asset preparation.
- Release runbook stop conditions are checked before publication in EP-09.

### E) Dependencies And Critical-Path Notes

Depends on EP-07. Target version should be the next patch after the completed `OK-PR-012` release.

### F) Tasks Checklist

- [ ] Re-read `docs/repo/runbooks/release-operating-kit.md`.
- [ ] Update `release-notes.md` with the selected patch version and routing-standard release notes.
- [ ] Update `pyproject.toml` package version to the selected patch version.
- [ ] Update `src/codeheart_operating_kit/__init__.py` package version to the selected patch version.
- [ ] Update `scripts/build-release-assets.py` default release version to the selected patch version.
- [ ] Update `bootstrap.md` release references to the selected patch version.
- [ ] Update `install.sh` release references to the selected patch version.
- [ ] Update `install.ps1` release references to the selected patch version.
- [ ] Update `manifest.yaml` release metadata and URLs for the selected patch version.
- [ ] Update `src/codeheart_operating_kit/resources/manifest.yaml` release metadata and URLs for the selected patch version.
- [ ] Update `tests/fixtures/release-manifest.json` to match the selected patch version.
- [ ] Run `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest`.
- [ ] Run `uv run --with pip --with setuptools --with wheel python scripts/build-release-assets.py --version <selected-patch-version> --output-dir dist`.
- [ ] Update `manifest.yaml` with built release asset checksums from `dist/`.
- [ ] Confirm `src/codeheart_operating_kit/resources/manifest.yaml` keeps placeholder downloadable asset checksums.
- [ ] Run `python3 scripts/validate-release-manifest.py manifest.yaml`.
- [ ] Run `git diff --check`.
- [ ] Record release-preparation evidence in the execution log.

### G) Implementation Notes

Use the selected patch version consistently. Replace `<selected-patch-version>` with the actual
version recorded in EP-00 before running release asset commands.

### H) Open Questions

OQ-1 applies.

## EP-09 - Release Publication And First-Consumer Proof

### A) Epic ID, Title, And Outcome

EP-09 - Release Publication And First-Consumer Proof

Outcome: After explicit public release approval, the patch release is published and one consumer
can adopt the routing standard through normal update, sync, and check.

### B) Scope

Publish the GitHub release from the validated commit and run one consumer proof. Do not sync named
private consumers unless the user explicitly adds them during activation.

### C) Files Touched

```text
dist/                                                                 # validated release assets
docs/repo/plans/operation-routing-dispatch-standard/
  operation-routing-dispatch-standard_execution_log.md                 # update at execution
```

The consumer proof touches an installed consumer's managed `.codeheart/kit/` snapshot through
normal `codeheart-operating-kit sync` behavior.

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Explicit public release approval is recorded before tag creation.
- Public tag for the selected patch version exists at the validated commit.
- GitHub release includes release notes, manifest, installers, assets, and checksums.
- Consumer update-check detects the published version.
- Consumer sync refreshes managed routing standard docs.
- Consumer check reports no managed-content drift after sync.
- Execution log records release URL, asset names, checksums, validation evidence, consumer proof,
  and residual risk.

### E) Dependencies And Critical-Path Notes

Depends on EP-08. Stop before tag creation when explicit public release approval is absent.

### F) Tasks Checklist

- [ ] Confirm explicit public release approval is recorded in the execution log.
- [ ] Confirm the validated commit matches the release commit with `git status --short` and `git rev-parse HEAD`.
- [ ] Confirm release assets and checksum files exist in `dist/`.
- [ ] Confirm `release-notes.md` covers routing-standard consumer impact.
- [ ] Create public tag for the selected patch version from the validated commit.
- [ ] Publish GitHub release for the selected patch version with `bootstrap.md`, `install.sh`, `install.ps1`, `release-notes.md`, `manifest.yaml`, release assets, and checksum files.
- [ ] Run `codeheart-operating-kit update-check <consumer-repository-path> --agent-notification`.
- [ ] Run `codeheart-operating-kit sync <consumer-repository-path>`.
- [ ] Run `codeheart-operating-kit check <consumer-repository-path> --json`.
- [ ] Record release URL, asset names, checksums, consumer proof, and residual risk in the execution log.

### G) Implementation Notes

Do not include private consumer names, local machine paths, tenant details, or raw operational logs
in public release notes. Use the existing release workflow and release asset patterns from recent
Operating Kit releases.

### H) Open Questions

OQ-2 applies.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Domain-specific route cards for Microsoft 365, GitHub, finance, CRM, product modules, or
  customer-specific workflows are deferred to the owning domains.
- Operational recipe maturity, phase markers, script promotion, and automation maturity standards
  are deferred to a separate discovery after the routing standard is stable enough to depend on.
- Machine-readable route-card schemas, route validators, and route-card linting are deferred
  until repeated Markdown route-card usage proves a stable shape.
- Module or extension manifest guidance is deferred until README or route-registry capability
  advertisements prove insufficient.
- Automatic consumer scaffolding for route registries, route cards, or consumer-owned state is
  deferred because V1 is instruction-only.

## 4.2 Future Considerations

- After the release, draft a module-adoption implementation plan for the M365 workspace module
  that adapts its capability map and operation runbooks to the established routing standard.
- After two or more modules adopt route cards, review whether the route-card field set needs
  schema support.
- After repeated fresh low-context probe usage, consider a lightweight probe template or review
  checklist.

# Revision Notes

- 2026-06-25: Created draft implementation plan from the approved Operation Routing And Dispatch
  Standard discovery capability scope.
- 2026-06-25: Tightened EP-01 inventory scope, added EP-07 probe matrix requirements, and fixed
  EP-08 release checksum sequencing after plan review.
