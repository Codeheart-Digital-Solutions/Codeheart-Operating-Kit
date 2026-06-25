Last updated: 2026-06-25T14:04:53Z (UTC)
Created: 2026-06-25
Status: draft

# Routing Surface Inventory

This is the EP-01 provisional inventory for the Operation Routing And Dispatch Standard
implementation plan.

Purpose: identify existing Operating Kit routing and instruction surfaces before new routing
doctrine is written. This inventory is plan-scoped evidence. It proposes source-edit targets for
the implementation epics, but it does not itself finalize wording, activate the implementation
plan, or change managed route behavior.

Consumer impact class: `instruction-only change` when implemented. This attachment does not
change generated paths, sync behavior, schemas, validators, release authority, consumer-owned
state, or external-system behavior.

## Inventory Row Shape

Each entry uses the approved row fields:

- `Surface`
- `Current routing role`
- `Instruction type`
- `Authority level`
- `Overlaps with`
- `Conflict risk`
- `Provisional disposition`
- `Reason`
- `Implementation note`

## Provisional Disposition Values

- `leave-as-is`: no source edit is currently expected beyond final recheck.
- `link-to-new-standard`: keep the existing surface and add a concise route or reference to the
  new routing standard.
- `consolidate-into-new-standard`: move generic routing doctrine into the new standard and leave
  only a short local pointer or local-specific rule.
- `keep-as-domain-specific-detail`: preserve the surface as a concrete instance after the generic
  standard exists.
- `defer`: wait for the final consolidation pass before deciding whether this source needs a
  direct edit.
- `retire-or-reword`: remove, narrow, or rephrase wording that would conflict with the new
  standard.

## Provisional Inventory

### INV-001 - Root Managed Block Template

- `Surface`: `templates/agents/AGENTS.managed-block.md`
- `Current routing role`: first installed agent-facing bootstrap; contains immediate rules, direct
  managed routes, module/extension state routing, and tooling-readiness routing.
- `Instruction type`: template, root router, managed bootstrap.
- `Authority level`: early managed agent interface for installed consumers.
- `Overlaps with`: Agent Interface routing reference, root contract, module-extension state,
  tooling readiness, fallback inventory.
- `Conflict risk`: high if it remains silent about route-before-surface behavior, because agents
  see this before deeper doctrine.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: root should expose the compact hierarchy and pointer, not duplicate the full doctrine.
- `Implementation note`: EP-03 should add only compact routing hierarchy language and a route to
  the full Agent Interface reference.

### INV-002 - Installed Kit Fallback Inventory

- `Surface`: `components/agent-interface/managed/kit-readme.md`
- `Current routing role`: installed fallback inventory when root routes do not match or kit state
  looks unclear, missing, stale, or damaged.
- `Instruction type`: managed fallback router.
- `Authority level`: managed installed kit inventory.
- `Overlaps with`: root managed block, Agent Interface README, Structure Governance README,
  module-extension state, tooling readiness.
- `Conflict risk`: medium if unclear requests route through inventory but the routing standard is
  absent from inventory.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: fallback inventory should tell agents where to find generic routing doctrine when the
  direct route is unclear.
- `Implementation note`: EP-03 should add a short Operation Routing route without turning this
  file into the full route registry.

### INV-003 - Agent Interface Router

- `Surface`: `components/agent-interface/managed/README.md`
- `Current routing role`: managed domain router for root agent introduction, local extension
  guidance, runbook authoring, tooling readiness, and fallback inventory.
- `Instruction type`: managed component README router.
- `Authority level`: managed domain router.
- `Overlaps with`: root contract, runbook authoring standard, tooling readiness, fallback
  inventory.
- `Conflict risk`: high until the new standard is routable from the domain that owns agent
  behavior.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: Agent Interface is the planned owner of the full routing doctrine.
- `Implementation note`: EP-02 should add the new
  `reference/operation-routing-and-dispatch.md` route.

### INV-004 - Root AGENTS.md Contract

- `Surface`: `components/agent-interface/managed/reference/root-agents-md-contract.md`
- `Current routing role`: defines the root `AGENTS.md` layer model, direct managed routes, fallback
  inventory, and compactness boundary.
- `Instruction type`: managed reference contract.
- `Authority level`: durable managed contract for installed root routing.
- `Overlaps with`: root managed block template, fallback inventory, Agent Interface README.
- `Conflict risk`: medium if root compactness is preserved but route-before-surface visibility is
  not defined.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: the contract should specify how much routing hierarchy belongs in root and where the
  detailed doctrine lives.
- `Implementation note`: EP-03 should update this contract with the compact hierarchy and pointer
  rule.

### INV-005 - Runbook Authoring Standard

- `Surface`: `components/agent-interface/managed/reference/runbook-authoring-standard.md`
- `Current routing role`: defines durable runbook audience classes, source of truth, execution
  lane, preconditions, approvals, stop conditions, evidence, validation, and tooling-readiness
  separation.
- `Instruction type`: managed reference standard.
- `Authority level`: durable managed runbook-shape standard.
- `Overlaps with`: new route-card standard, planning workflow hooks, tooling readiness.
- `Conflict risk`: medium because it already has dispatch-like fields but does not distinguish
  route selection from recipe execution.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: keep runbook maturity guidance minimal here while linking routing-bearing runbooks to
  the new standard.
- `Implementation note`: EP-05 should add lightweight route-card and routing-bearing runbook
  guidance without defining new runbook maturity shapes.

### INV-006 - Tooling Readiness Runbook

- `Surface`: `components/agent-interface/managed/runbooks/handle-tooling-readiness.md`
- `Current routing role`: shared route for missing local tools encountered by repositories,
  modules, or extensions, with service blockers returned to the calling module.
- `Instruction type`: managed hybrid runbook.
- `Authority level`: managed procedure for generic local tooling blockers.
- `Overlaps with`: module-owned runbooks, runbook authoring standard, root managed block tooling
  route.
- `Conflict risk`: low if preserved as a downstream blocker route after operation routing selects
  the owner and task.
- `Provisional disposition`: `keep-as-domain-specific-detail`
- `Reason`: this is a concrete generic blocker route, not the full routing doctrine.
- `Implementation note`: EP-06 should recheck whether a one-line pointer is useful; EP-02 through
  EP-05 do not need to reshape this runbook.

### INV-007 - Structure Governance Router

- `Surface`: `components/structure-governance/managed/README.md`
- `Current routing role`: routes agents to placement, naming, managed boundaries, module state,
  runbook authoring, placement-change, and index-maintenance guidance.
- `Instruction type`: managed component README router.
- `Authority level`: managed structure-governance router.
- `Overlaps with`: documentation structure, managed-content boundaries, runbook authoring,
  module-extension state.
- `Conflict risk`: medium if routing artifacts have no explicit placement owner.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: structure governance should route placement questions for route artifacts to the
  correct placement reference while leaving behavior doctrine to Agent Interface.
- `Implementation note`: EP-04 should add concise routes for routing artifact placement and owner
  boundaries.

### INV-008 - Documentation Structure Reference

- `Surface`: `components/structure-governance/managed/reference/documentation-structure.md`
- `Current routing role`: durable placement reference for docs, routers, references, runbooks,
  plans, state, product docs, and index maintenance.
- `Instruction type`: managed reference contract.
- `Authority level`: durable managed placement doctrine.
- `Overlaps with`: new route-artifact placement rules, module-extension state, runbook authoring,
  planning lifecycle.
- `Conflict risk`: high until route registries, route cards, and capability advertisements have
  explicit placement guidance.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: this is the right owner for where routing artifacts live, not for their behavioral
  semantics.
- `Implementation note`: EP-04 should define placement for capability advertisements, route
  registries, route cards, local wrappers, and plan-scoped routing evidence.

### INV-009 - Managed Content Boundaries Reference

- `Surface`: `components/structure-governance/managed/reference/managed-content-boundaries.md`
- `Current routing role`: distinguishes managed, scaffold, template, consumer-owned, local-user,
  generated, report, and committed module-state ownership modes.
- `Instruction type`: managed reference contract.
- `Authority level`: durable managed ownership-boundary doctrine.
- `Overlaps with`: route-card placement, capability advertisements, local wrappers, module-state
  routing.
- `Conflict risk`: medium if route artifacts are introduced without ownership-mode boundaries.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: route artifacts need managed versus consumer-owned placement boundaries.
- `Implementation note`: EP-04 should add only ownership-boundary language needed for routing
  artifacts and wrappers.

### INV-010 - Module Extension State Reference

- `Surface`: `components/structure-governance/managed/reference/module-extension-state.md`
- `Current routing role`: defines committed non-secret module and extension routing state under
  `docs/repo/state/<module-or-extension-id>/` and the live-preflight boundary.
- `Instruction type`: managed reference for a specific routing-state pattern.
- `Authority level`: durable managed placement and state-routing doctrine for modules/extensions.
- `Overlaps with`: new generic routing standard, root managed block, managed boundaries,
  documentation structure.
- `Conflict risk`: medium if this specific state-routing rule is mistaken for the whole routing
  doctrine.
- `Provisional disposition`: `keep-as-domain-specific-detail`
- `Reason`: preserve it as the module/extension state instance of the broader routing model.
- `Implementation note`: EP-04 should add a concise cross-reference to the generic routing
  standard while keeping module-state details here.

### INV-011 - Planning Workflows Router

- `Surface`: `components/planning-workflows/managed/README.md`
- `Current routing role`: managed router for discovery, implementation planning, execution,
  review, lifecycle, and plan-register maintenance.
- `Instruction type`: managed component README router.
- `Authority level`: managed planning-domain router.
- `Overlaps with`: discovery workflow, implementation-plan drafting, execution, review, runbook
  authoring standard.
- `Conflict risk`: low to medium; planning workflows should enforce routing-standard adoption for
  routing-bearing plans but should not own routing doctrine.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: a light route or note can make planning hooks discoverable without bloating the
  planning router.
- `Implementation note`: EP-05 should add a concise pointer only if needed after workflow runbooks
  are updated.

### INV-012 - Discovery Workflow Runbook

- `Surface`: `components/planning-workflows/managed/runbooks/discovery-workflow.md`
- `Current routing role`: governs discovery lifecycle, decision state, output targets, and
  implementation capability-scope handoff.
- `Instruction type`: managed discovery runbook.
- `Authority level`: managed planning runbook for discovery.
- `Overlaps with`: routing-trigger categories, capability-advertisement decisions,
  implementation handoff, review-ready versus handoff-ready states.
- `Conflict risk`: low if it only flags routing-bearing scope; medium if it starts defining
  routing mechanics.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: discovery should identify routing-bearing scope and handoff expectations, not define
  the full routing doctrine.
- `Implementation note`: EP-05 should add a light scope flag and route to the standard.

### INV-013 - Draft Implementation Plan Runbook

- `Surface`: `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
- `Current routing role`: turns accepted discovery, user direction, and targeted research into
  execution-ready implementation plans; already enforces discovery capability-scope handoff.
- `Instruction type`: managed implementation-planning runbook.
- `Authority level`: managed planning runbook for implementation-plan creation.
- `Overlaps with`: routing-standard adoption, runbook-change coverage, feature capability
  coverage, review readiness.
- `Conflict risk`: high for future routing-bearing work if this runbook does not require plans to
  cite and apply the new standard.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: implementation planning is the correct enforcement point for routing-bearing changes.
- `Implementation note`: EP-05 should require routing-bearing plans to apply the new standard and
  include fresh low-context routing probes when relevant.

### INV-014 - Execute Implementation Plan Runbook

- `Surface`: `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
- `Current routing role`: governs active implementation execution, epic review gates, execution
  logs, validation, and plan-register hooks.
- `Instruction type`: managed execution runbook.
- `Authority level`: managed implementation execution procedure.
- `Overlaps with`: routing-probe evidence, final consolidation evidence, runbook change execution.
- `Conflict risk`: medium if routing-bearing epic execution can complete without probe evidence
  or standard-adoption checks.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: execution should verify routing-bearing epics against the routing standard.
- `Implementation note`: EP-05 should add execution evidence expectations for routing probes and
  standard adoption.

### INV-015 - Review Planning Document Runbook

- `Surface`: `components/planning-workflows/managed/runbooks/review-planning-document.md`
- `Current routing role`: reviews discovery and implementation documents for ambiguity, scope,
  decision quality, execution readiness, capability coverage, and runbook coverage.
- `Instruction type`: managed review runbook.
- `Authority level`: managed planning-document review procedure.
- `Overlaps with`: routing-standard adoption checks, fresh low-context routing probe expectations,
  discovery capability-scope preservation.
- `Conflict risk`: high if reviewers do not catch routing-bearing plans that invent local routing
  rules or skip probes.
- `Provisional disposition`: `link-to-new-standard`
- `Reason`: review must catch missing routing-standard adoption before execution.
- `Implementation note`: EP-05 should add targeted review checks for routing-bearing discovery,
  implementation plans, and probe evidence.

### INV-016 - Native Codex Capabilities Router

- `Surface`: `components/native-codex-capabilities/managed/README.md`
- `Current routing role`: describes baseline native Codex capability expectations and routes to
  the baseline capability profile.
- `Instruction type`: managed component README router.
- `Authority level`: managed capability-awareness router, not workflow doctrine.
- `Overlaps with`: execution-surface awareness, connector/tool visibility, installation checks.
- `Conflict risk`: medium if visible native capabilities are mistaken for routing authority.
- `Provisional disposition`: `keep-as-domain-specific-detail`
- `Reason`: the component should remain capability awareness, while the new routing standard
  should explain that visible surfaces are not automatically the highest routing authority.
- `Implementation note`: EP-02 should cite native capabilities as one class of execution surface;
  no direct edit is expected unless EP-06 finds ambiguity.

### INV-017 - Source Repo AGENTS.md

- `Surface`: `AGENTS.md`
- `Current routing role`: maintainer bootstrap for the public Operating Kit source repository,
  including public-core safety, task-matched read order, source of truth, change safety, and
  validation.
- `Instruction type`: source-repository maintainer router.
- `Authority level`: nearest source-repo instruction surface for maintainers.
- `Overlaps with`: README, docs router, placement contract, change runbook, managed root template.
- `Conflict risk`: low; this file routes maintainers of the kit source, not installed consumer
  operations.
- `Provisional disposition`: `leave-as-is`
- `Reason`: adding the consumer operation-routing doctrine here would blur source maintainer
  routing with installed agent behavior.
- `Implementation note`: EP-06 should recheck only if the final standard changes maintainer read
  order or public-core safety.

### INV-018 - Source Repo README

- `Surface`: `README.md`
- `Current routing role`: public repository purpose, boundary, and maintainer entry points.
- `Instruction type`: source-repository overview router.
- `Authority level`: public repository entry point.
- `Overlaps with`: docs index, placement contract, change runbook, component routers.
- `Conflict risk`: low; it should remain an entry point, not detailed routing doctrine.
- `Provisional disposition`: `leave-as-is`
- `Reason`: managed Agent Interface and Structure Governance are better owners for the new
  standard and placement rules.
- `Implementation note`: no EP-02 through EP-05 edit expected.

### INV-019 - Docs Index

- `Surface`: `docs/README.md`
- `Current routing role`: public documentation router for repository governance, runbooks, plans,
  managed components, schemas, scripts, CLI source, and tests.
- `Instruction type`: source documentation index.
- `Authority level`: public docs router.
- `Overlaps with`: repo documentation index, component READMEs, plan bundle indexes.
- `Conflict risk`: low to medium; direct links may be useful after the new standard exists, but
  over-indexing could make the docs router noisy.
- `Provisional disposition`: `defer`
- `Reason`: decide in EP-06 whether the new routing reference needs a top-level docs index entry
  or whether the Agent Interface route is enough.
- `Implementation note`: EP-06 final consolidation should decide the index behavior after EP-02
  through EP-05 source surfaces exist.

### INV-020 - Repo Documentation Index

- `Surface`: `docs/repo/README.md`
- `Current routing role`: public repository-governance index for references, runbooks, and formal
  plans.
- `Instruction type`: source repository documentation index.
- `Authority level`: public repo governance router.
- `Overlaps with`: docs index, plan register, plan bundle documents, component routers.
- `Conflict risk`: low to medium; it should route plan evidence but not own generic routing
  doctrine.
- `Provisional disposition`: `defer`
- `Reason`: the implementation plan and final standard may warrant index updates, but exact
  links should be decided after the standard exists.
- `Implementation note`: EP-06 should decide whether to add the implementation plan and final
  reference links.

### INV-021 - Source Placement Contract

- `Surface`: `docs/repo/reference/placement-contract.md`
- `Current routing role`: public placement contract for source repository areas, installed
  consumer areas, ownership modes, component targets, consumer-owned boundaries, and placement
  rules.
- `Instruction type`: source repository reference.
- `Authority level`: public source placement contract.
- `Overlaps with`: managed documentation-structure reference, managed-content boundaries,
  consumer-impact classification.
- `Conflict risk`: low; existing component target rules already cover managed references and
  consumer-owned state.
- `Provisional disposition`: `leave-as-is`
- `Reason`: V1 does not add generated paths, scaffolds, schemas, or placement-contract behavior.
- `Implementation note`: no EP-02 through EP-05 source edit expected unless EP-04 discovers a
  placement-contract gap.

### INV-022 - Consumer Impact Classification

- `Surface`: `docs/repo/reference/consumer-impact-classification.md`
- `Current routing role`: classifies consumer impact for instruction-only, validator-only,
  component addition, scaffold addition, migration, placement, and safety-policy changes.
- `Instruction type`: source repository reference.
- `Authority level`: source governance classification reference.
- `Overlaps with`: release notes, change runbook, release runbook, implementation plan.
- `Conflict risk`: low; routing-standard V1 is already representable as instruction-only.
- `Provisional disposition`: `leave-as-is`
- `Reason`: no new impact class is needed for route-before-surface managed guidance.
- `Implementation note`: no EP-02 through EP-05 source edit expected.

### INV-023 - Change Operating Kit Runbook

- `Surface`: `docs/repo/runbooks/change-operating-kit.md`
- `Current routing role`: maintainer procedure before changing kit source, docs, schemas,
  templates, validators, installers, release assets, or CLI behavior.
- `Instruction type`: source maintainer runbook.
- `Authority level`: source repository change procedure.
- `Overlaps with`: source AGENTS, consumer-impact classification, placement contract, release
  runbook.
- `Conflict risk`: low; it already routes maintainers to the relevant source governance checks.
- `Provisional disposition`: `leave-as-is`
- `Reason`: EP-01 does not change the source change procedure.
- `Implementation note`: no EP-02 through EP-05 source edit expected.

### INV-024 - Release Operating Kit Runbook

- `Surface`: `docs/repo/runbooks/release-operating-kit.md`
- `Current routing role`: maintainer procedure for release validation, assets, checksums, tags,
  GitHub release, and evidence.
- `Instruction type`: source maintainer release runbook.
- `Authority level`: source repository release procedure.
- `Overlaps with`: release notes, consumer-impact classification, release asset tests, plan
  release epics.
- `Conflict risk`: low; release routing is a separate maintainer process.
- `Provisional disposition`: `leave-as-is`
- `Reason`: the routing standard does not alter release authority or release procedure.
- `Implementation note`: no EP-02 through EP-05 source edit expected.

### INV-025 - Promote Consumer Change Runbook

- `Surface`: `docs/repo/runbooks/promote-consumer-change.md`
- `Current routing role`: maintainer procedure for promoting reusable consumer-local rules,
  runbooks, templates, or workflows into the public kit.
- `Instruction type`: source maintainer runbook.
- `Authority level`: source repository promotion procedure.
- `Overlaps with`: public-core safety, placement contract, consumer-impact classification,
  structure governance.
- `Conflict risk`: low; it can already route reusable consumer guidance through placement and
  impact checks.
- `Provisional disposition`: `leave-as-is`
- `Reason`: no new promotion procedure is required for V1.
- `Implementation note`: no EP-02 through EP-05 source edit expected.

### INV-026 - Triage Kit Feedback Runbook

- `Surface`: `docs/repo/runbooks/triage-kit-feedback.md`
- `Current routing role`: maintainer procedure for public feedback issue triage, lifecycle route
  selection, public-safe handling, and release/sync notes.
- `Instruction type`: source maintainer runbook.
- `Authority level`: source repository feedback triage procedure.
- `Overlaps with`: discovery workflow, implementation planning, promote-consumer-change runbook,
  release notes.
- `Conflict risk`: low; feedback triage can continue routing accepted reusable work to discovery
  or implementation planning.
- `Provisional disposition`: `leave-as-is`
- `Reason`: feedback lifecycle routing is distinct from operation-routing doctrine.
- `Implementation note`: no EP-02 through EP-05 source edit expected.

## Candidate Implementation Edit List

### EP-02 - Agent Interface Routing Reference

Candidate source edits:

- Create `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`.
- Update `components/agent-interface/managed/README.md` with a route to the new reference.
- Update component metadata and packaged-resource mirrors for Agent Interface when EP-02 is
  implemented.
- Use `components/native-codex-capabilities/managed/README.md` as an input example for execution
  surfaces, but do not make native capabilities an authority source.

Rows informing EP-02: INV-003, INV-005, INV-016.

### EP-03 - Root And Fallback Routing Visibility

Candidate source edits:

- Update `templates/agents/AGENTS.managed-block.md` with compact route-before-surface hierarchy
  language and a pointer to the new routing reference.
- Update `components/agent-interface/managed/kit-readme.md` with an installed fallback route to
  operation routing.
- Update `components/agent-interface/managed/reference/root-agents-md-contract.md` to define the
  compact root boundary for routing hierarchy.
- Update packaged-resource mirrors and focused tests that assert installed root/fallback route
  visibility.

Rows informing EP-03: INV-001, INV-002, INV-004.

### EP-04 - Structure Governance Placement Rules

Candidate source edits:

- Update `components/structure-governance/managed/README.md` with routing-artifact placement
  routes.
- Update `components/structure-governance/managed/reference/documentation-structure.md` with
  placement rules for capability advertisements, route registries, route cards, and local
  wrappers.
- Update `components/structure-governance/managed/reference/managed-content-boundaries.md` with
  ownership-mode guidance for route artifacts.
- Update `components/structure-governance/managed/reference/module-extension-state.md` with a
  concise cross-reference that preserves module state as a specific routing-context instance.
- Update component metadata and packaged-resource mirrors for Structure Governance when EP-04 is
  implemented.

Rows informing EP-04: INV-007, INV-008, INV-009, INV-010, INV-021.

### EP-05 - Runbook And Planning Workflow Hooks

Candidate source edits:

- Update `components/agent-interface/managed/reference/runbook-authoring-standard.md` with
  minimal routing-bearing runbook guidance and route-card reference language.
- Update `components/planning-workflows/managed/README.md` only if the updated runbooks need a
  discoverability pointer.
- Update `components/planning-workflows/managed/runbooks/discovery-workflow.md` with a light
  routing-bearing scope flag.
- Update `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` so
  routing-bearing implementation plans cite and apply the routing standard.
- Update `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` so
  routing-bearing epics record standard-adoption and fresh low-context probe evidence.
- Update `components/planning-workflows/managed/runbooks/review-planning-document.md` so reviewers
  catch missing routing-standard adoption and missing probe evidence.
- Update component metadata and packaged-resource mirrors for Planning Workflows when EP-05 is
  implemented.

Rows informing EP-05: INV-005, INV-011, INV-012, INV-013, INV-014, INV-015.

### EP-06 - Final Disposition And Consolidation

Final consolidation should re-read this inventory after EP-02 through EP-05 have concrete text.
It should then add a `Final Consolidation Map` that records each row as kept, linked,
consolidated, reworded, deferred, or intentionally unchanged.

Candidate final-pass surfaces:

- Recheck source repo routers and governance docs for over-linking or missing discoverability:
  INV-017 through INV-026.
- Recheck tooling readiness and native capability surfaces for possible one-line cross-reference:
  INV-006 and INV-016.
- Decide whether `docs/README.md` and `docs/repo/README.md` should link the final standard
  directly or only through component routers: INV-019 and INV-020.

## Final Consolidation Map

Decision-only status: draft for user review. These decisions record recommended final
dispositions after EP-02 through EP-05. They do not apply the remaining consolidation edits.

### D-001 - Root Managed Block Template / INV-001

- Final disposition: linked to new standard.
- Final recommendation: keep the EP-03 compact route-before-surface rule and managed route pointer
  as the root exposure. Do not add route-card fields, capability catalogs, or domain examples to
  the managed block.
- Rationale: the root block must be visible early, but it should stay a bootstrap router.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Copernicus; no follow-up before user review.

### D-002 - Installed Kit Fallback Inventory / INV-002

- Final disposition: linked to new standard.
- Final recommendation: keep the `Operation Routing` entry in the installed kit fallback
  inventory. Do not duplicate the full hierarchy in the fallback list.
- Rationale: `.codeheart/kit/README.md` is a fallback inventory, not the doctrine owner.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Cicero; no follow-up before user review.

### D-003 - Agent Interface Router / INV-003

- Final disposition: linked to new standard.
- Final recommendation: keep the Agent Interface README as the primary component-level route to
  `reference/operation-routing-and-dispatch.md`.
- Rationale: Agent Interface owns agent-facing behavior, so its router should expose the standard
  directly.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Hypatia; no follow-up before user review.

### D-004 - Root AGENTS.md Contract / INV-004

- Final disposition: linked to new standard.
- Final recommendation: keep the compact-routing-hierarchy contract and anti-catalog boundary.
  Do not promote deep module capability details into root `AGENTS.md`.
- Rationale: the contract explains why the root sees the route early without becoming a stale
  upper-layer catalog.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Aquinas; no follow-up before user review.

### D-005 - Runbook Authoring Standard / INV-005

- Final disposition: linked to new standard.
- Final recommendation: keep the routing-bearing runbook hook and pointer to the operation routing
  reference. Do not define additional runbook maturity shapes in this routing epic.
- Rationale: the runbook authoring standard should tell runbooks when to expose a routing contract,
  while route-card semantics remain in the routing reference.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Bohr; no follow-up before user review.

### D-006 - Tooling Readiness Runbook / INV-006

- Final disposition: kept as downstream specialty route.
- Final recommendation: leave the tooling readiness runbook unchanged in this consolidation pass.
  Use it after a route identifies a local tooling blocker; do not add routing doctrine to the
  tooling runbook.
- Rationale: the routing standard already points blocked tool availability to tooling readiness,
  and duplicating the hierarchy there would blur route selection with blocker handling.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Helmholtz; no follow-up before user review.

### D-007 - Structure Governance Router / INV-007

- Final disposition: linked to new standard.
- Final recommendation: keep the split where Structure Governance owns placement and boundary
  rules while Agent Interface owns routing behavior and field semantics.
- Rationale: this preserves a clean behavior-vs-placement boundary.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Anscombe; no follow-up before user review.

### D-008 - Documentation Structure Reference / INV-008

- Final disposition: linked to new standard.
- Final recommendation: keep the routing artifact placement section and the relative link to the
  routing behavior reference. Do not centralize every route card in this structure reference.
- Rationale: the reference should decide where artifacts belong, not copy owner-domain routing
  registries.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Dewey; no follow-up before user review.

### D-009 - Managed Content Boundaries Reference / INV-009

- Final disposition: linked to new standard.
- Final recommendation: keep the ownership split for generic routing behavior, consumer
  route-state, and domain-specific route artifacts.
- Rationale: managed-content boundaries are the right place to prevent doctrine from drifting into
  consumer-owned local state.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Dalton; no follow-up before user review.

### D-010 - Module Extension State Reference / INV-010

- Final disposition: kept as domain-specific detail.
- Final recommendation: keep the existing module-state routing rule and cross-reference to the
  generic routing standard. Do not fold module state mechanics into the generic routing reference.
- Rationale: the generic standard owns behavior; this reference owns committed module/extension
  state placement and state-vs-live-truth boundaries.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Peirce; no follow-up before user review.

### D-011 - Planning Workflows Router / INV-011

- Final disposition: linked to new standard.
- Final recommendation: keep the planning-workflows README hook for plans that create or
  materially change routing-bearing surfaces.
- Rationale: the router should expose when planning work must apply the routing standard, without
  duplicating the standard itself.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with McClintock; no follow-up before user review.

### D-012 - Discovery Workflow Runbook / INV-012

- Final disposition: linked to new standard.
- Final recommendation: keep the lightweight `routing-bearing` discovery scope flag and pointer to
  the operation-routing standard.
- Rationale: discovery should identify routing-bearing scope, but implementation planning owns the
  detailed adoption checklist.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Boole; no follow-up before user review.

### D-013 - Draft Implementation Plan Runbook / INV-013

- Final disposition: linked to new standard.
- Final recommendation: keep the routing-standard coverage section and probe requirement for
  routing-bearing implementation plans.
- Rationale: implementation planning is where approved capability scope becomes concrete routing
  work and validation.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Mencius; no follow-up before user review.

### D-014 - Execute Implementation Plan Runbook / INV-014

- Final disposition: linked to new standard.
- Final recommendation: keep routing-standard adoption and fresh low-context probe evidence in the
  execution checklist/log expectations.
- Rationale: routing-bearing implementation needs execution evidence that the route works for a
  fresh agent, not only source text.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Confucius; no follow-up before user review.

### D-015 - Review Planning Document Runbook / INV-015

- Final disposition: linked to new standard.
- Final recommendation: keep the review checks for missing routing-standard adoption, affected
  routing surfaces, and fresh low-context probe evidence.
- Rationale: planning review should catch routing gaps before execution starts.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Wegener; no follow-up before user review.

### D-016 - Native Codex Capabilities Router / INV-016

- Final disposition: kept as execution-surface detail.
- Final recommendation: leave the native-capability surface unchanged in this consolidation pass.
  The routing standard already classifies visible tools, connectors, MCP tools, CLIs, browser
  surfaces, and APIs as execution-surface documentation below higher routing authorities.
- Rationale: editing native capability docs would make the routing epic reach outside its owner
  boundary and duplicate the core rule in the wrong surface.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Beauvoir; no follow-up before user review.

### D-017 - Source Repo AGENTS.md / INV-017

- Final disposition: intentionally unchanged.
- Final recommendation: leave the source repository `AGENTS.md` unchanged for this routing
  standard. Keep it focused on public-core maintainer instructions and repository-specific change
  authority.
- Rationale: installed consumer routing belongs in managed kit content; source-repo maintainer
  instructions should not become a second doctrine entry point.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Descartes; no follow-up before user review.

### D-018 - Source Repo README / INV-018

- Final disposition: intentionally unchanged.
- Final recommendation: leave the public source `README.md` unchanged for this routing standard.
- Rationale: the README is a high-level product and maintainer entry point. The standard is
  discoverable through component routers and plan artifacts without turning the public README into
  a detailed doctrine index.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Hooke; no follow-up before user review.

### D-019 - Docs Index / INV-019

- Final disposition: intentionally unchanged.
- Final recommendation: do not add a direct `operation-routing-and-dispatch.md` link to
  `docs/README.md` in this pass. Keep discovery through the Agent Interface component and plan
  indexes.
- Rationale: the top-level docs index should remain a broad navigation surface. Directly linking
  every new standard from it would create an upper-layer catalog problem.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Ptolemy; no follow-up before user review.

### D-020 - Repo Documentation Index / INV-020

- Final disposition: intentionally unchanged.
- Final recommendation: do not add a direct `operation-routing-and-dispatch.md` link to
  `docs/repo/README.md` in this pass. Keep repository documentation focused on maintainer
  governance, plans, references, and release/change processes.
- Rationale: `docs/repo/README.md` should route source-repo governance work, not become the
  installed-agent behavior index.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Volta; no follow-up before user review.

### D-021 - Source Placement Contract / INV-021

- Final disposition: intentionally unchanged.
- Final recommendation: leave `docs/repo/reference/placement-contract.md` unchanged for this
  routing standard.
- Rationale: V1 is instruction-only and does not introduce generated consumer route-card
  scaffolding or new public placement-contract obligations.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Sartre; no follow-up before user review.

### D-022 - Consumer Impact Classification / INV-022

- Final disposition: reworded.
- Final recommendation: keep the clarified `instruction-only change` definition in
  `docs/repo/reference/consumer-impact-classification.md`. It now explicitly covers additive
  managed instruction, reference, runbook, template, or doc files under an existing component
  target when they do not add a component, create or move consumer-owned scaffolds or generated
  paths, change validators, sync or ownership behavior, safety policy, or require consumer action.
- Rationale: the routing standard is instruction-only in behavior, but this implementation also
  adds a new managed installed reference target under an existing component. The prior
  classification text said instruction-only changes do not change generated paths, which could be
  read too narrowly for additive managed documentation targets.
- Source edits in this pass: applied to
  `docs/repo/reference/consumer-impact-classification.md`.
- Subagent review: CONVERGED with Ohm on the revised recommendation before implementation;
  post-edit review READY with Banach, Faraday, and Socrates. No findings.

### D-023 - Change Operating Kit Runbook / INV-023

- Final disposition: intentionally unchanged.
- Final recommendation: leave `docs/repo/runbooks/change-operating-kit.md` unchanged for this
  routing standard.
- Rationale: maintainer changes already pass through placement, impact, tests, and release
  governance; routing-bearing implementation planning now carries the specific standard hook.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Locke; no follow-up before user review.

### D-024 - Release Operating Kit Runbook / INV-024

- Final disposition: intentionally unchanged.
- Final recommendation: leave `docs/repo/runbooks/release-operating-kit.md` unchanged for this
  routing standard.
- Rationale: release publication mechanics do not change because the routing standard is a managed
  instruction-only update.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Arendt; no follow-up before user review.

### D-025 - Promote Consumer Change Runbook / INV-025

- Final disposition: intentionally unchanged.
- Final recommendation: leave `docs/repo/runbooks/promote-consumer-change.md` unchanged for this
  routing standard.
- Rationale: consumer promotion mechanics are not changed by the new route-before-surface
  doctrine.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Boyle; no follow-up before user review.

### D-026 - Triage Kit Feedback Runbook / INV-026

- Final disposition: intentionally unchanged.
- Final recommendation: leave `docs/repo/runbooks/triage-kit-feedback.md` unchanged for this
  routing standard.
- Rationale: feedback lifecycle stays separate from routing doctrine; future feedback can still
  propose route-standard changes through the existing feedback route.
- Source edits in this pass: not applied; no additional source edit recommended.
- Subagent review: CONVERGED with Maxwell; no follow-up before user review.

## Final Disposition Summary

- Linked to new standard: D-001, D-002, D-003, D-004, D-005, D-007, D-008, D-009,
  D-011, D-012, D-013, D-014, D-015.
- Kept as downstream specialty, domain-specific, or execution-surface detail: D-006, D-010,
  D-016.
- Intentionally unchanged: D-017, D-018, D-019, D-020, D-021, D-023, D-024, D-025, D-026.
- Deferred: none.
- Consolidate into new standard: none.
- Reworded: D-022.

## Provisional Disposition Summary

- `link-to-new-standard`: INV-001, INV-002, INV-003, INV-004, INV-005, INV-007, INV-008,
  INV-009, INV-011, INV-012, INV-013, INV-014, INV-015.
- `keep-as-domain-specific-detail`: INV-006, INV-010, INV-016.
- `leave-as-is`: INV-017, INV-018, INV-021, INV-022, INV-023, INV-024, INV-025, INV-026.
- `defer`: INV-019, INV-020.
- `consolidate-into-new-standard`: none provisionally.
- `retire-or-reword`: none provisionally.

No final source wording should be inferred from this summary. It is an implementation input for
EP-02 through EP-06.
