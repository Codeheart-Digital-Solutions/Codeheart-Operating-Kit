Last updated: 2026-07-04T22:26:18Z (UTC)

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
maintenance. `OK-PR-024` is the completed source implementation plan for removing Python/pip from
the base Operating Kit bootstrap by porting root commands to a self-contained Go CLI, adding macOS
and Windows binary release packs, legacy Python-wheel migration, and explicit behavior parity
tests.
`OK-PR-022` is the completed implementation plan for generic runbook-to-script promotion doctrine,
including reusable script asset guidance, current doctrine alignment, and instruction-only release
readiness. `OK-PR-021` is the completed implementation plan for Operating Kit-guided repo feedback
capture with Codeheart organization membership gating, check-first GitHub Issues availability,
demand-driven issue-intake setup, release, and approved consumer sync scope.

## Entries

## OK-PR-024 - Operating Kit Self-Contained Bootstrap Implementation

Type: implementation-plan
Purpose: Implement a self-contained Operating Kit bootstrap for fresh macOS and Windows machines
by replacing the Python-wheel-first root CLI distribution with a Go CLI, binary release packs,
legacy Python-wheel migration, and behavior parity tests against the current Python CLI.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md
docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md
Created: 2026-07-04
Last updated: 2026-07-04T22:26:18Z (UTC)
Priority / ordering note: Source implementation is complete and validated. Public release
publication, live manifest pointer switch, signing/notarization decision, Git tag, GitHub release,
and consumer sync remain separate approval-gated release-run work.

Relations:
- depends-on: Codeheart-HQ:CODEHEART-HQ-PR-009 - Operating Kit Self-Contained Bootstrap Discovery
- related: OK-PR-017 - Consumer Runtime Materialization Hardening Implementation
- related: OK-PR-010 - Tooling Environment Readiness Implementation
- related: docs/repo/runbooks/change-operating-kit.md
- related: docs/repo/runbooks/release-operating-kit.md

Session refs:
- created: 2026-07-04, not recorded, drafted after the user approved the discovery capability
  scope, including Go as the root CLI implementation language, macOS and Windows first-class
  support, legacy Python-wheel migration, release-stage trust gates, and explicit behavior parity
  tests.
- material update: 2026-07-04, not recorded, revised the draft after planning review and user
  approval to clarify release-candidate boundaries, setup-only onboarding for optional native
  capabilities, incremental parity baselines, atomic installer migration, mandatory macOS
  universal assets, and Python installer flag removal or deprecation behavior.
- material update: 2026-07-04, active goal, activated the implementation plan, created the
  execution log, recorded source preflight evidence, and confirmed publication remains out of
  scope.
- material update: 2026-07-04, active goal, completed source implementation, generated staged
  macOS universal and Windows x64 release-candidate packs, proved local macOS no-Python install,
  proved Windows no-Python install through GitHub Actions, and recorded the low-context bootstrap
  probe.
- completed: 2026-07-04, active goal, closed the source implementation after Go tests, Python
  parity tests, installer tests, release asset tests, schema/public-core/Markdown validation,
  GitHub Actions macOS and Windows validation, and fresh review gates passed.

Coordination note:
- Canonical source implementation plan for the HQ-owned discovery.
- Public release publication, tagging, and named consumer sync are not in this implementation plan
  unless separately approved through the release runbook.
- Consumer impact classification: migration required, validator-only change, instruction-only
  change, and security or safety policy change.
- Validation summary: local `go test`, Python-vs-Go parity, installer, release asset,
  release-manifest, public-core, Markdown, staged macOS install, and low-context bootstrap proof
  passed; GitHub Actions Validate run `28721439492` passed macOS and Windows validation.
- Staged source assets remain local release candidates. Root `manifest.yaml` live asset URLs and
  checksums were not switched to unpublished staged assets.

## OK-PR-023 - Plan Register Dirty Target Safety Implementation

Type: implementation-plan
Purpose: Clarify managed plan-register doctrine so agents distinguish unrelated dirty repository
state from target-register conflicts before deciding whether to update a coordination-home
register directly or record pending sync.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/plan-register-dirty-target-safety/plan-register-dirty-target-safety_implementation_doc.md
docs/repo/plans/plan-register-dirty-target-safety/plan-register-dirty-target-safety_execution_log.md
Created: 2026-06-29
Last updated: 2026-06-29T20:23:54Z (UTC)
Priority / ordering note: Completed and shipped in Operating Kit `v0.1.18`. Consumers receive the
refined dirty target-register rule after normal install, update, or sync.

Relations:
- related: OK-PR-013 - Operation Routing And Dispatch Standard Implementation
- related: OK-PR-004 - Plan Register Portfolio Doctrine Refinement
- related: OK-PR-003 - Coordination Home Register ID Namespace
- related: OK-PR-022 - Runbook-To-Script Promotion Standard Implementation

Session refs:
- created: 2026-06-29, not recorded, drafted after user discussion clarified that unrelated dirty
  worktree state should not block compatible coordination-home register updates.
- material update: 2026-06-29, not recorded, activated implementation and created the execution
  log.
- material update: 2026-06-29, not recorded, accepted review-gate finding and aligned the
  existing pending-sync scaffold wording without adding a new scaffold path.
- completed: 2026-06-29, not recorded, completed managed doctrine, hook, scaffold wording,
  packaged-resource mirror, validation, and review-gate closeout.
- material update: 2026-06-29, not recorded, released Operating Kit `v0.1.18`, published the
  GitHub release, and recorded release validation evidence in the execution log.

Coordination note:
- Generic Operating Kit implementation candidate.
- Consumer impact classification: instruction-only change.
- Work-board behavior is explicitly out of scope.
- Released in Operating Kit `v0.1.18`.

## OK-PR-022 - Runbook-To-Script Promotion Standard Implementation

Type: implementation-plan
Purpose: Implement managed Operating Kit doctrine for promoting fragile, repeated, or
evidence-bearing runbook mechanics into reusable script assets without turning whole runbooks into
premature CLIs, APIs, or broad automation wrappers.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_implementation_doc.md
docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_execution_log.md
Created: 2026-06-29
Last updated: 2026-06-29T15:01:31Z (UTC)
Priority / ordering note: Should execute after `OK-PR-019` recommendation approval and before
downstream Foundry, Microsoft 365, AI Execution, or consumer repositories adopt script promotion
rules.

Relations:
- depends-on: OK-PR-019 - Runbook-To-Script Promotion Standard Discovery
- related: OK-PR-014 - Operational Recipe Maturity Standard Implementation Pointer
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: OK-PR-011 - Operation Routing And Dispatch Standard Discovery
- related: OK-PR-013 - Operation Routing And Dispatch Standard Implementation
- related: OK-PR-017 - Consumer Runtime Materialization Hardening Implementation

Session refs:
- created: 2026-06-29, not recorded, drafted after the user approved the recommended defaults
  for output contract, script placement, helper policy, validator scope, adoption scope, migration
  rule, and review gates.
- material update: 2026-06-29, not recorded, activated for implementation and began source
  managed doctrine changes.
- completed: 2026-06-29, not recorded, completed the instruction-only source implementation,
  resource mirrors, indexes, stale-vocabulary review, validation, and low-context routing probe.
- material update: 2026-06-29, not recorded, released Operating Kit `v0.1.17` with the
  runbook-to-script promotion standard and recorded release validation evidence in the execution
  log.

Coordination note:
- Generic Operating Kit implementation candidate.
- Consumer impact classification: instruction-only change.
- Release notes are required when shipped.
- No consumer migration is required in this implementation scope.

## OK-PR-021 - Repo Feedback Capture And Issue Intake Implementation

Type: implementation-plan
Purpose: Implement managed Operating Kit repo feedback capture by adding optional `repo_feedback`
config schema support, capture and setup runbooks, item-format guidance, installed root route
visibility, packaged-resource mirrors, validation, release surfaces, and fresh-repo proof while
keeping GitHub issue creation and repository setup approval-gated.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/repo-feedback-capture/repo-feedback-capture_implementation_doc.md
docs/repo/plans/repo-feedback-capture/repo-feedback-capture_execution_log.md
Created: 2026-06-29
Last updated: 2026-07-02T13:16:41Z (UTC)
Priority / ordering note: Should execute after `OK-PR-020` capability scope approval and before
agents are expected to use repo feedback capture in consumer repositories.

Relations:
- depends-on: OK-PR-020 - Repo Feedback Capture And Issue Intake Discovery
- related: docs/repo/plans/kit-feedback-intake/kit-feedback-intake_implementation_doc.md
- related: OK-PR-011 - Operation Routing And Dispatch Standard Discovery
- related: OK-PR-013 - Operation Routing And Dispatch Standard Implementation
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: OK-PR-010 - Tooling Environment Readiness Implementation

Session refs:
- created: 2026-06-29, not recorded, drafted after user approved the repo feedback capture
  recommendations, including check-first GitHub Issues behavior, demand-driven setup,
  disabled/suppressed state after decline, and missing-label fallback.
- material update: 2026-07-02, active goal, activated implementation and created the execution
  log after plan refinements for Codeheart organization membership gating, no-fallback behavior,
  release, and approved consumer sync.
- material update: 2026-07-02, active goal, completed source implementation, validation,
  packaged-resource proof, GitHub authorization proof, routing probe, and review-gate fix before
  release execution.
- completed: 2026-07-02, active goal, released Operating Kit `v0.1.19`, validated public release
  install paths, and synced approved consumer repositories while preserving unrelated worktree
  changes.

Coordination note:
- Generic Operating Kit implementation candidate.
- Consumer impact classification: instruction-only change, validator-only change, and security or
  safety policy change.
- Execution should not create live repo feedback issues. Release and consumer sync are in scope
  after validation and explicit release authority.
- Released in Operating Kit `v0.1.19`.

## OK-PR-020 - Repo Feedback Capture And Issue Intake Discovery

Type: discovery-plan
Purpose: Discover a reusable Operating Kit standard for how agents recognize repo-specific
feedback during runbook execution, check whether the owning repository's GitHub Issues intake
already works, route to it when available, handle demand-driven setup only when issue intake is
unavailable or incomplete, suppress repeated prompts after user decline, classify issues, protect
privacy, and promote accepted feedback into direct patches, batches, discovery, or implementation
planning.
Status: implementation-handoff-ready
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/repo-feedback-capture/repo-feedback-capture_discovery_doc.md
Created: 2026-06-29
Last updated: 2026-07-02T13:16:41Z (UTC)
Priority / ordering note: Should be reviewed before implementing managed repo-feedback capture
routes, check-first GitHub issue availability, issue-intake setup runbooks, config schema support,
or prompt-suppression behavior.

Relations:
- child: OK-PR-021 - Repo Feedback Capture And Issue Intake Implementation
- related: docs/repo/plans/kit-feedback-intake/kit-feedback-intake_discovery_doc.md
- related: docs/repo/plans/kit-feedback-intake/kit-feedback-intake_implementation_doc.md
- related: OK-PR-011 - Operation Routing And Dispatch Standard Discovery
- related: OK-PR-013 - Operation Routing And Dispatch Standard Implementation
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: OK-PR-009 - Tooling Environment Readiness Discovery
- related: OK-PR-010 - Tooling Environment Readiness Implementation

Session refs:
- created: 2026-06-29, not recorded, drafted from user discussion about repo-specific feedback
  capture, GitHub Issues as durable inbox, demand-driven setup, suppression after decline,
  classification, privacy, and triage promotion.
- material update: 2026-06-29, not recorded, clarified that agents should check whether GitHub
  Issues already works before offering setup, following the Operating Kit feedback intake
  precedent where Issues were already enabled.
- material update: 2026-07-02, active goal, handed off to active implementation in `OK-PR-021`
  after Codeheart organization membership gating and no-fallback behavior were accepted.

Coordination note:
- Generic Operating Kit doctrine candidate. Repo-specific issue destinations and GitHub settings
  should remain owned by the target repository and require explicit approval before external
  changes.

## OK-PR-019 - Runbook-To-Script Promotion Standard Discovery

Type: discovery-plan
Purpose: Discover a reusable Operating Kit standard for deciding when runbook steps or
operational recipes should be promoted to scripts, what first-script scaffolding is required, when
infrastructure helpers or domain folders are justified, when package/CLI promotion is warranted,
and how promoted scripts preserve runbook approval, safety, output, and testing expectations.
Status: draft
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_discovery_doc.md
Created: 2026-06-29
Last updated: 2026-06-29T14:30:44Z (UTC)
Priority / ordering note: Should be accepted and implemented in the Operating Kit before
downstream repositories or modules use it to promote deterministic runbook recipes into scripts.

Relations:
- child: OK-PR-022 - Runbook-To-Script Promotion Standard Implementation
- related: OK-PR-014 - Operational Recipe Maturity Standard Implementation Pointer
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: OK-PR-011 - Operation Routing And Dispatch Standard Discovery
- related: OK-PR-013 - Operation Routing And Dispatch Standard Implementation
- related: OK-PR-009 - Tooling Environment Readiness Discovery
- related: OK-PR-017 - Consumer Runtime Materialization Hardening Implementation

Session refs:
- created: 2026-06-29, not recorded, drafted from user discussion about generic script
  promotion triggers, first-script quality scaffolding, helper abstraction timing, domain folder
  triggers, and package/CLI maturity.
- material update: 2026-06-29, not recorded, refined the discovery to remove tested inline block
  maturity, add current doctrine alignment inventory, and define the mechanical script output
  contract.

Coordination note:
- Generic Operating Kit doctrine candidate. Foundry and Microsoft 365 module adoption should be
  handled by separate downstream work only after this standard is accepted and implemented.

## OK-PR-018 - Business Docs Placement Clarity Implementation

Type: implementation-plan
Purpose: Clarify the managed Structure Governance `docs/business/` placement rule so it means
company or organization business-operating records when a consumer repository intentionally stores
them, and explicitly does not mean software product architecture, module design, platform solution
design, application business logic, or implementation planning.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/business-docs-placement-clarity/business-docs-placement-clarity_implementation_doc.md
docs/repo/plans/business-docs-placement-clarity/business-docs-placement-clarity_execution_log.md
Created: 2026-06-26
Last updated: 2026-06-26T21:09:26Z (UTC)
Priority / ordering note: Instruction-only managed doctrine clarification. No scaffold, sync,
schema, validator, CLI, or consumer migration behavior should change.

Relations:
- related: OK-PR-013 - Operation Routing And Dispatch Standard Implementation
- related: OK-PR-017 - Consumer Runtime Materialization Hardening Implementation
- related: CODEHEART-AUTOMATION-FOUNDRY-PR-009 - Relational Workspace View Module Discovery

Session refs:
- created: 2026-06-26, not recorded, drafted after a Codeheart-HQ placement correction moved a
  reusable relational workspace view module discovery from HQ business docs into the Foundry repo.
- completed: 2026-06-26, not recorded, activated and completed the instruction-only managed
  wording clarification with release-note, mirror, index, register, execution-log, and validation
  evidence.
- material update: 2026-06-26, not recorded, prepared Operating Kit `v0.1.16` release surfaces,
  assets, installer proof, and release-readiness validation.
- material update: 2026-06-26, not recorded, published Operating Kit `v0.1.16` and recorded
  GitHub release plus workflow validation evidence in the execution log.

Coordination note:
- Canonical plan is Operating Kit-owned.
- Triggering consumer clarification was applied ad hoc in Codeheart-HQ local docs.

## OK-PR-017 - Consumer Runtime Materialization Hardening Implementation

Type: implementation-plan
Purpose: Harden the Operating Kit local tooling and agent-interface standards so consumer-mode
runtime tooling is materialized from durable module/package content into ignored local runtime
state without editable development links, generated install metadata in managed snapshots, or
global runtime mutation, and so runbooks use visible-terminal handoff for user-entered terminal
prompts instead of hidden agent tool prompts; release the Operating Kit update; apply the first
adopter through an AI Execution module release; refresh the HQ AI Execution snapshot; and sync the
released kit into HQ, Foundry, the named private platform repository, and Operating Kit repos.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/runtime-materialization-hardening/runtime-materialization-hardening_implementation_doc.md
docs/repo/plans/runtime-materialization-hardening/runtime-materialization-hardening_execution_log.md
Created: 2026-06-26
Last updated: 2026-06-26T16:33:33Z (UTC)
Priority / ordering note: Should execute after OK-PR-016 because it hardens the newly established
`.codeheart/local/` and `python-runtime` lane. It should execute before retrying AI Execution
auth setup in HQ so `foundry-ai` does not depend on editable-install metadata in the managed
snapshot.

Relations:
- depends-on: OK-PR-016 - Local Runtime Environment Standard Implementation
- related: OK-PR-009 - Tooling Environment Readiness Discovery
- related: OK-PR-010 - Tooling Environment Readiness Implementation
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: CODEHEART-AUTOMATION-FOUNDRY-PR-007 - Foundry AI Execution Module Implementation Plan

Session refs:
- created: 2026-06-26, not recorded, drafted implementation plan from user-approved direction
  after AI Execution auth setup exposed brittle editable-install behavior in a managed consumer
  snapshot.
- material update: 2026-06-26, not recorded, expanded scope to include visible-terminal handoff
  doctrine for runbooks that require user-entered terminal input.
- material update: 2026-06-26, not recorded, activated the implementation plan and created the
  execution log.
- completed: 2026-06-26, not recorded, released Operating Kit `v0.1.15`, released and proved AI
  Execution `0.1.1` through HQ's repo-local runtime, synced the released kit into all named repos,
  resolved final review findings, and marked the execution log complete.

Coordination note:
- promoted into the Codeheart-HQ coordination register as CODEHEART-OPERATING-KIT-PR-017.

## OK-PR-016 - Local Runtime Environment Standard Implementation

Type: implementation-plan
Purpose: Implement the accepted Operating Kit local runtime standard, including `.codeheart/local/`
as ignored local machine/runtime state, `.codeheart/local/envs/python/` as the default Python venv,
additive kit config/schema behavior, init/sync `.gitignore` readiness, managed readiness doctrine,
tests, packaged-resource mirrors, and downstream Foundry AI Execution handoff.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/local-runtime-environment-standard/local-runtime-environment-standard_implementation_doc.md
docs/repo/plans/local-runtime-environment-standard/local-runtime-environment-standard_execution_log.md
Created: 2026-06-26
Last updated: 2026-06-26T14:38:50Z (UTC)
Priority / ordering note: Should execute after the accepted local runtime discovery and before
patching Foundry AI Execution onboarding against the shared `.codeheart/local/envs/python/`
convention.

Relations:
- child: OK-PR-017 - Consumer Runtime Materialization Hardening Implementation
- depends-on: OK-PR-015 - Local Runtime Environment Standard Discovery
- related: OK-PR-009 - Tooling Environment Readiness Discovery
- related: OK-PR-010 - Tooling Environment Readiness Implementation
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: Foundry AI Execution Module Implementation Plan

Session refs:
- created: 2026-06-26, not recorded, drafted implementation plan from the accepted discovery
  capability scope with generated behavior, managed doctrine, validation, packaging, coordination,
  and downstream Foundry handoff sequencing.
- material update: 2026-06-26, not recorded, added a narrow managed planning-workflow
  clarification so discovery and implementation plan drafting selects the canonical owning
  repository before using coordination-home register pointers.
- material update: 2026-06-26, not recorded, activated the implementation plan, created the
  sibling execution log, and began source implementation for generated behavior, managed doctrine,
  packaging parity, validation, and HQ coordination pointers.
- completed: 2026-06-26, not recorded, implemented the source standard, validated focused and
  full test suites, recorded the fresh low-context routing probe, and left public release
  publication, named consumer sync, and downstream AI Execution module adoption deferred.

Coordination note:
- promoted into the Codeheart-HQ coordination register as CODEHEART-OPERATING-KIT-PR-016.

## OK-PR-015 - Local Runtime Environment Standard Discovery

Type: discovery-plan
Purpose: Discover a reusable Operating Kit convention for ignored repo-local machine/runtime state,
including the default Python virtual environment path for Foundry and Operating Kit tooling such as
`foundry-ai`, and clarify how generic environment-readiness routing handles missing local tools in
any blocked agent interaction.
Status: draft
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/local-runtime-environment-standard/local-runtime-environment-standard_discovery_doc.md
Created: 2026-06-26
Last updated: 2026-06-26T14:10:39Z (UTC)
Priority / ordering note: Accepted discovery should precede implementation-plan drafting for the
generic local runtime standard, then AI Execution onboarding can be patched against the shared
`.codeheart/local/envs/python/` convention.

Relations:
- child: OK-PR-016 - Local Runtime Environment Standard Implementation
- related: OK-PR-009 - Tooling Environment Readiness Discovery
- related: OK-PR-010 - Tooling Environment Readiness Implementation
- related: OK-PR-005 - Runbook Authoring Standards Discovery
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: Foundry AI Execution Module Implementation Plan

Session refs:
- created: 2026-06-26, not recorded, moved the discovery from HQ into the Operating Kit source
  repo after confirming the canonical scope is reusable managed Operating Kit doctrine.
- material update: 2026-06-26, not recorded, broadened scope from Python venv placement to
  generic local tooling blocker routing, machine/user-level baseline tooling, and
  shared-vs-purpose-specific Python venv defaults.
- material update: 2026-06-26, not recorded, accepted recommended defaults for `.codeheart/local/`,
  `.codeheart/local/envs/python/`, `.gitignore` handling, optional
  `local_machine_layer_path`, on-demand first-run behavior, governed machine/user-level tooling,
  and generic blocker routing.

Coordination note:
- promoted into the Codeheart-HQ coordination register as CODEHEART-OPERATING-KIT-PR-015.

## OK-PR-014 - Operational Recipe Maturity Standard Implementation Pointer

Type: implementation-plan
Purpose: Compact source-repo pointer to the HQ-owned implementation plan that adds a generic
Operating Kit standard for identifying, structuring, validating, reviewing, and deliberately
promoting operational recipes inside runbooks without forcing premature scripts, commands,
wrappers, APIs, Foundry packaging conventions, or M365 module adoption.
Status: completed
Owner / repository: Codeheart-HQ
Canonical docs:
Codeheart-HQ:docs/repo/plans/operational-recipe-maturity-standard/operational-recipe-maturity-standard_discovery_doc.md
Codeheart-HQ:docs/repo/plans/operational-recipe-maturity-standard/operational-recipe-maturity-standard_implementation_doc.md
Codeheart-HQ:docs/repo/plans/operational-recipe-maturity-standard/operational-recipe-maturity-standard_execution_log.md
Created: 2026-06-25
Last updated: 2026-06-25T20:42:30Z (UTC)
Priority / ordering note: Source implementation touches this repository, but the accepted
discovery and canonical execution log live in Codeheart-HQ. This entry is a pointer only and does
not duplicate the full plan.

Relations:
- related: OK-PR-013 - Operation Routing And Dispatch Standard Implementation
- related: OK-PR-007 - Runbook Authoring Standards Implementation
- related: OK-PR-008 - Module Extension State Routing Implementation

Session refs:
- created: 2026-06-25, not recorded, added compact source-repo pointer during HQ-owned plan
  activation.
- completed: 2026-06-25, not recorded, source implementation completed under HQ-owned plan.
  Release publication and consumer sync remain deferred.

Coordination note:
- canonical plan is HQ-owned
- local-only source-repo pointer

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
