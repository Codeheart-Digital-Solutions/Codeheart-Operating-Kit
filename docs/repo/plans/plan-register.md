Last updated: 2026-06-25T13:45:59Z (UTC)

# Plan Register

This kit-initialized consumer state file lists important formal plans, plan families, major
workstreams, and portfolio-relevant planning records for this repository.

Operating Kit owns the file contract and format. This repository owns the entries after creation.
Sync may recreate this baseline when the file is absent, but it must not overwrite existing
entries.

Follow `.codeheart/kit/docs/planning-workflows/reference/plan-register-format.md` for entry
fields and `.codeheart/kit/docs/planning-workflows/runbooks/maintain-plan-register.md` for
maintenance.

## Register Coverage

Coverage note: This register currently lists public Operating Kit repository plans that have been
entered during plan-register adoption. Earlier repository plans may be added during later register
maintenance.

## Entries

## OK-PR-013 - Operation Routing And Dispatch Standard Implementation

Type: implementation-plan
Purpose: Implement managed Operating Kit operation-routing and dispatch doctrine, including a
two-pass routing-surface inventory and consolidation flow, Agent Interface routing reference,
compact root route, Structure Governance placement rules, runbook and planning hooks, fresh
low-context routing probes, packaged-resource mirrors, validation, and release preparation.
Status: draft
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/operation-routing-dispatch-standard/operation-routing-dispatch-standard_implementation_doc.md
Created: 2026-06-25
Last updated: 2026-06-25T13:14:40Z (UTC)
Priority / ordering note: Should execute after OK-PR-012 completes because both plans touch
planning-workflow runbooks and release surfaces.

Relations:
- depends-on: OK-PR-011 - Operation Routing And Dispatch Standard Discovery
- depends-on: OK-PR-012 - Discovery Handoff Gate Implementation
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-006 - Module Extension State Routing Discovery
- related: OK-PR-009 - Tooling Environment Readiness Discovery
- related: future module route-registry adoption work

Session refs:
- created: 2026-06-25, not recorded, drafted implementation plan from the approved routing
  discovery capability scope with inventory-first sequencing, routing-bearing workflow hooks, and
  fresh low-context routing probe validation.
- material update: 2026-06-25, not recorded, tightened the implementation plan after review by
  adding source-governance inventory tasks, a fresh low-context probe matrix attachment
  requirement, and explicit release checksum sequencing.
- material update: 2026-06-25, not recorded, changed the routing-surface inventory model to a
  two-pass flow with provisional inventory before doctrine writing and final consolidation after
  the core standard, root routing, placement rules, and planning hooks exist.
- material update: 2026-06-25, not recorded, directly implemented selected EP-02 through EP-05
  managed-source and packaged-resource changes with validation, fresh low-context routing probes,
  and a new execution log while leaving full plan activation, release prep, and consumer sync
  pending.

Coordination note:
- local-only

## OK-PR-012 - Discovery Handoff Gate Implementation

Type: implementation-plan
Purpose: Implement a managed planning-workflow guardrail that stops normal implementation-plan
drafting from discovery documents whose implementation capability scope has not been approved,
delegated, or explicitly revised; publish the instruction-only `v0.1.13` release; and sync the
released kit into the named consumer repositories.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/discovery-handoff-gate/discovery-handoff-gate_implementation_doc.md
docs/repo/plans/discovery-handoff-gate/discovery-handoff-gate_execution_log.md
Created: 2026-06-25
Last updated: 2026-06-25T13:45:59Z (UTC)
Priority / ordering note: Should execute before drafting the Operation Routing And Dispatch
Standard implementation plan so that discovery capability-scope approval is enforced by the
planning workflow itself.

Relations:
- related: OK-PR-011 - Operation Routing And Dispatch Standard Discovery
- related: OK-PR-002 - Codeheart Operating Kit Implementation-Planning Quality
- related: managed planning-workflows release path

Session refs:
- created: 2026-06-25, not recorded, drafted a narrow implementation plan for the discovery
  handoff gate, embedded release and named consumer sync approval, and planned `v0.1.13`
  instruction-only release execution.
- material update: 2026-06-25, not recorded, activated the implementation plan for goal-style
  execution and created the sibling execution log.
- completed: 2026-06-25, published Operating Kit `v0.1.13`, verified the public release
  workflows, synced the released kit into the three named consumer repositories, and recorded final
  review evidence in the execution log.

Coordination note:
- local-only

## OK-PR-011 - Operation Routing And Dispatch Standard Discovery

Type: discovery-plan
Purpose: Discover shared Operating Kit doctrine for pre-execution routing and dispatch, including
authority hierarchy, capability advertisements, route registries, route cards, ambiguity handling,
and the split between generic Operating Kit routing behavior and domain-owned route details.
Status: draft
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/operation-routing-dispatch-standard/operation-routing-dispatch-standard_discovery_doc.md
Created: 2026-06-25
Last updated: 2026-06-25T12:18:33Z (UTC)
Priority / ordering note: Should precede module-specific adoption work that depends on a shared
route-before-execution-surface standard.

Relations:
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-006 - Module Extension State Routing Discovery
- related: OK-PR-009 - Tooling Environment Readiness Discovery
- related: future module route-registry adoption work

Session refs:
- created: 2026-06-25, not recorded, drafted first discovery for generic operation routing and
  dispatch doctrine after module-operation routing discussions identified the need for a reusable
  Kit-owned standard.
- material update: 2026-06-25, not recorded, added an existing routing-surface inventory and
  consolidation requirement before implementation planning.
- material update: 2026-06-25, not recorded, added routing trigger categories, simple
  authority-conflict handling, fresh-agent validation expectations, and capability-advertisement
  maintenance requirements.
- material update: 2026-06-25, not recorded, added a candidate authority hierarchy, resolved the
  inventory granularity question as a dedicated implementation epic, and defined fresh
  low-context routing probes as a planning/review hook for routing-bearing epics.
- material update: 2026-06-25, not recorded, clarified that the candidate authority hierarchy is
  an Operating Kit routing heuristic and does not replace native Codex instruction priority.
- material update: 2026-06-25, not recorded, added compact root `AGENTS.md` hierarchy
  requirement, minimum capability-advertisement fields, route-card `not applicable` allowance, and
  candidate implementation sequence.
- material update: 2026-06-25, not recorded, moved `Intent aliases` into the minimum V1
  capability advertisement field set.
- material update: 2026-06-25, not recorded, resolved primary owner, route-card field set, and
  capability-advertisement placement as accepted defaults and marked the discovery
  implementation-handoff-ready.
- material update: 2026-06-25, not recorded, added formal implementation capability scope and
  corrected discovery status back to manual-review-ready until scope approval or revision.
- material update: 2026-06-25, not recorded, added future implementation standard-adoption scope
  so routing-bearing implementation work must apply the established routing standard, with a
  light discovery hook and stronger implementation-planning/review hooks.
- material update: 2026-06-25, not recorded, recorded user approval of the implementation
  capability scope and marked the discovery implementation-handoff-ready.

Coordination note:
- local-only

## OK-PR-010 - Tooling Environment Readiness Implementation

Type: implementation-plan
Purpose: Implement one central Operating Kit tooling-readiness route, installed route visibility,
runbook-authoring and planning hooks, packaged-resource mirroring, tests, and `0.1.12` release
preparation for local environment blockers encountered during module onboarding or operation.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/tooling-environment-readiness/tooling-environment-readiness_implementation_doc.md
Created: 2026-06-24
Last updated: 2026-06-24T14:37:00Z (UTC)
Priority / ordering note: Should execute after OK-PR-009 is accepted and before modules depend on
a managed Operating Kit route for missing package managers, runtimes, CLIs, PowerShell modules, or
other local tooling blockers.

Relations:
- depends-on: OK-PR-009 - Tooling Environment Readiness Discovery
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: OK-PR-008 - Module Extension State Routing Implementation
- related: Foundry Microsoft 365 module onboarding tooling-readiness discussion

Session refs:
- created: not recorded
- material update: 2026-06-24, not recorded, activated implementation plan and created sibling
  execution log.
- material update: 2026-06-24, not recorded, added a low-risk structure-governance
  cross-reference to keep placement guidance separate from runbook-shape and tooling-readiness
  standards.
- material update: 2026-06-24, not recorded, clarified EP-04 release-manifest strategy after
  review: root manifest records publishable asset hashes while packaged manifest keeps
  zero-placeholder downloadable asset hashes to avoid self-referential archive checksums.
- completed: 2026-06-24, published Operating Kit `v0.1.12`, verified release assets, completed
  isolated consumer proof, and passed GitHub Actions workflow run
  `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28106384737`.

Coordination note:
- local-only

## OK-PR-009 - Tooling Environment Readiness Discovery

Type: discovery-plan
Purpose: Discover shared Operating Kit doctrine for local tooling and environment readiness,
including missing-tool routing, package-manager/bootstrap guidance boundaries, module-owned tool
declarations, local readiness evidence placement, and approval-gated install or repair flows.
Status: draft
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/tooling-environment-readiness/tooling-environment-readiness_discovery_doc.md
Created: 2026-06-24
Last updated: 2026-06-24T13:38:01Z (UTC)
Priority / ordering note: Should precede any Operating Kit implementation that adds shared
tooling-readiness routes, package-manager/bootstrap guidance, readiness state placement, or
module-facing tool declaration standards.

Relations:
- promoted-from: runbook-authoring-standards related feedback item "Shared Environment Readiness
  And Tooling Register"
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-006 - Module Extension State Routing Discovery
- related: Foundry Microsoft 365 module onboarding tooling-readiness discussion

Session refs:
- created: not recorded
- material update: 2026-06-24, not recorded, refined discovery around an on-demand baseline
  tooling catalog, flexible module declarations by manifest/reference/runbook, and a generic
  missing-tool behavior contract.
- material update: 2026-06-24, not recorded, clarified that module onboarding and operation
  environment blockers should trigger the readiness route and that package-manager guidance should
  use concrete nontechnical choices.
- material update: 2026-06-24, not recorded, added trigger-model review results from managed
  route surfaces and the Foundry M365 hybrid onboarding runbook, including route visibility,
  blocker-specific human choices, and the split between local environment blockers and
  module-owned service blockers.
- material update: 2026-06-24, not recorded, clarified the Operating Kit anti-sprawl boundary:
  keep one central readiness route and small baseline catalog while leaving module-specific install
  runbooks and commands in modules.

Coordination note:
- local-only

## OK-PR-008 - Module Extension State Routing Implementation

Type: implementation-plan
Purpose: Implement managed Operating Kit structure-governance doctrine, generic agent routing,
packaged-resource mirroring, tests, and release preparation for committed non-secret module and
extension state under `docs/repo/state/<id>/`.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/module-extension-state-routing/module-extension-state-routing_implementation_doc.md
Created: 2026-06-23
Last updated: 2026-06-23T18:26:27Z (UTC)
Priority / ordering note: Should execute after OK-PR-006 is accepted and before Foundry modules or
other extensions depend on a shared committed-state placement route.

Relations:
- depends-on: OK-PR-006 - Module Extension State Routing Discovery
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: managed structure-governance placement doctrine
- related: Foundry M365 workspace registry discussion

Session refs:
- created: not recorded
- material update: 2026-06-23, not recorded, updated implementation plan to target `0.1.11`
  release prep with explicit approval-gated publication proof.
- material update: 2026-06-23, not recorded, activated implementation plan and created sibling
  execution log.
- material update: 2026-06-23, not recorded, completed EP-01 through EP-04 local implementation,
  packaged-resource mirroring, `v0.1.11` release preparation, and validation; public release
  publication and consumer proof remain pending explicit approval.
- material update: 2026-06-23, not recorded, published `v0.1.11`, verified published assets,
  completed isolated consumer proof, and marked the implementation complete.

Coordination note:
- promoted into the Codeheart-HQ coordination register as CODEHEART-OPERATING-KIT-PR-008.

## OK-PR-007 - Runbook Authoring Standards Implementation

Type: implementation-plan
Purpose: Implement the managed Operating Kit runbook authoring standard with audience
classification, compact intention blocks, human-facing and agent-facing quality rules, planning
workflow hooks, packaged resource mirroring, public release publication, and consumer sync proof.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_implementation_doc.md
Created: 2026-06-23
Last updated: 2026-06-23T17:44:23Z (UTC)
Priority / ordering note: Should execute after the runbook authoring standards discovery is
accepted and before future Operating Kit runbook-quality guidance is promoted into managed
consumer routes.

Relations:
- depends-on: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-006 - Module Extension State Routing Discovery
- related: OK-PR-002 - Codeheart Operating Kit Implementation-Planning Quality
- related: managed first-run onboarding and onboarding context contract

Session refs:
- created: not recorded
- material update: 2026-06-23, not recorded, added activation evidence, public release
  publication, and consumer sync proof to the planned execution path.
- material update: 2026-06-23, not recorded, activated implementation plan and created sibling
  execution log.
- material update: 2026-06-23, not recorded, prepared v0.1.10 release assets, manifests, and
  validation evidence for publication.
- material update: 2026-06-23, not recorded, published v0.1.10 and completed isolated plus
  configured consumer sync proof.

Coordination note:
- promoted into the Codeheart-HQ coordination register as CODEHEART-OPERATING-KIT-PR-007.

## OK-PR-006 - Module Extension State Routing Discovery

Type: discovery-plan
Purpose: Discover the generic Operating Kit placement and routing convention for committed,
non-secret, consumer-owned module and extension state under `docs/repo/state/<id>/`.
Status: draft
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/module-extension-state-routing/module-extension-state-routing_discovery_doc.md
Created: 2026-06-23
Last updated: 2026-06-23T13:41:25Z (UTC)
Priority / ordering note: Should precede any structure-governance implementation that makes
module or extension state routing visible in managed `AGENTS.md`, kit fallback inventory, and
placement doctrine.

Relations:
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: Foundry M365 workspace registry discussion
- related: managed structure-governance placement doctrine

Session refs:
- created: not recorded
- material update: 2026-06-23, not recorded, added implementation scope handoff for a
  structure-governance instruction-only release.

Coordination note:
- promoted into the Codeheart-HQ coordination register as CODEHEART-OPERATING-KIT-PR-006.

## OK-PR-005 - Runbook Authoring Standards Discovery

Type: discovery-plan
Purpose: Discover reusable Operating Kit standards for human-facing, agent-facing, hybrid, and
maintainer runbooks so future runbooks have clear user scripts, explicit execution paths, source
of truth, stop conditions, evidence, and validation.
Status: draft
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_discovery_doc.md
Created: 2026-06-23
Last updated: 2026-06-23T12:02:14Z (UTC)
Priority / ordering note: Should precede the next Operating Kit runbook-quality implementation
release and inform the Foundry M365 onboarding UX hardening work.

Relations:
- related: OK-PR-002 - Codeheart Operating Kit Implementation-Planning Quality
- related: consumer-discovered Foundry M365 onboarding UX feedback
- related: managed first-run onboarding and onboarding context contract

Session refs:
- created: not recorded
- material update: 2026-06-23, not recorded, added the runbook sampling matrix attachment as the
  first standard test set.
- material update: 2026-06-23, not recorded, added condensed sampling lessons and the compact
  runbook intention block decision.
- material update: 2026-06-23, not recorded, added related Operating Kit feedback notes for slim
  local preferences and shared environment/tooling readiness.
- material update: 2026-06-23, not recorded, clarified language preference handling, audience
  modeling, planning/execution/review integration, and no active consumer/module runbook retrofit.

Coordination note:
- promoted into the Codeheart-HQ coordination register as CODEHEART-OPERATING-KIT-PR-005.

## OK-PR-004 - Plan Register Portfolio Doctrine Refinement

Type: implementation-plan
Purpose: Refine managed plan-register doctrine with local and coordination-home reference shapes,
repository-qualified ID guidance, canonical pointer examples, relation ownership, and public-safe
portfolio overview conventions.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_implementation_doc.md
Created: 2026-06-22
Last updated: 2026-06-22T21:18:13Z (UTC)
Priority / ordering note: Should execute before consumers rely on refined portfolio register
shapes for coordination-home population or repository-qualified local register IDs.

Relations:
- related: OK-PR-003 - Coordination Home Register ID Namespace
- related: OK-PR-001 - Plan Register Session And Lifecycle Hardening
- related: docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_implementation_doc.md
- related: consumer-local portfolio planning surfaces discovery handoff

Session refs:
- created: not recorded
- material update: 2026-06-22, not recorded, activated implementation plan.
- material update: 2026-06-22, not recorded, completed managed doctrine edits, packaged resource
  mirroring, `v0.1.9` release preparation, and local validation; public release publication and
  consumer sync proof remain pending explicit approval.
- material update: 2026-06-22, not recorded, published `v0.1.9`, completed first-consumer sync
  proof, and marked implementation complete.

Coordination note:
- local-only
- consumer-local discovery source sanitized for public-core doctrine planning

## OK-PR-003 - Coordination Home Register ID Namespace

Type: implementation-plan
Purpose: Add managed plan-register doctrine so coordination-home entries use unique IDs for
member-repository plans while preserving source local register IDs.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/coordination-home-register-id-namespace/coordination-home-register-id-namespace_implementation_doc.md
Created: 2026-06-22
Last updated: 2026-06-22T19:47:33Z (UTC)
Priority / ordering note: Completed `v0.1.8` instruction-only release and isolated consumer sync
proof.

Relations:
- related: OK-PR-001 - Plan Register Session And Lifecycle Hardening
- related: docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_implementation_doc.md

Session refs:
- created: not recorded
- material update: 2026-06-22, not recorded, activated implementation plan and created sibling
  execution log.
- material update: 2026-06-22, not recorded, completed release-preparation epics through local
  asset build and validation; public release publication remains pending explicit approval.
- material update: 2026-06-22, not recorded, completed EP-04 local pre-publication checks; public
  tag, GitHub release, and consumer sync proof remain pending.
- material update: 2026-06-22, not recorded, completed `v0.1.8` public release publication and
  isolated consumer update-check, sync, and check proof.

Coordination note:
- local-only

## OK-PR-002 - Codeheart Operating Kit Implementation-Planning Quality

Type: implementation-plan
Purpose: Update managed discovery, implementation-planning, planning-review, and execution
workflows so plans preserve intended feature capability and reject policy-only or under-covered
implementation.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_implementation_doc.md
Created: 2026-06-22
Last updated: 2026-06-22T18:50:19Z (UTC)
Priority / ordering note: Completed `v0.1.7` planning workflow quality release and first-consumer
sync proof.

Relations:
- related: first consumer repository discovery handoff -
  <first-consumer-repository>/docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_discovery_doc.md
- related: first consumer repository handoff implementation plan -
  <first-consumer-repository>/docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_implementation_doc.md

Session refs:
- created: 2026-06-22, session 019eef87-f252-7b91-aa50-ecf54b357c6c
- material update: 2026-06-22, session 019eef87-f252-7b91-aa50-ecf54b357c6c,
  completed `v0.1.7` release publication and first-consumer sync proof.

Coordination note:
- local-only

## OK-PR-001 - Plan Register Session And Lifecycle Hardening

Type: implementation-plan
Purpose: Harden Operating Kit plan-register doctrine for self-contained session-reference
resolution and lifecycle grouping.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/plan-register-session-lifecycle-hardening/plan-register-session-lifecycle-hardening_implementation_doc.md
Created: 2026-06-21
Last updated: 2026-06-21T19:30:29Z (UTC)
Priority / ordering note: Prepares the next additive hardening release after `v0.1.5`.

Relations:
- related: Portfolio coordination and plan-register implementation plan -
  docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_implementation_doc.md

Session refs:
- material update: 2026-06-21, not recorded, activated implementation plan.
- material update: 2026-06-21, not recorded, completed implementation plan.

Coordination note:
- local-only
