Last updated: 2026-07-08T14:11:45Z (UTC)

# Repo Documentation

This folder contains public repository governance for Codeheart Operating Kit.

## Contents

- `reference/placement-contract.md`: placement rules for kit docs, generated consumer surfaces,
  managed content, scaffolds, templates, and local user layers.
- `reference/consumer-impact-classification.md`: impact classes for changes that affect
  consumers, release notes, sync, validation, routing, generated surfaces, or safety policy.
- `reference/public-core-hygiene-inventory.md`: inventory of included, generalized, and excluded
  public-core extraction inputs.
- `reference/kit-feedback-label-taxonomy.md`: GitHub label taxonomy for public kit feedback
  intake.
- `runbooks/change-operating-kit.md`: ordered maintainer procedure for changing this repository.
- `runbooks/release-operating-kit.md`: ordered maintainer procedure for public releases.
- `runbooks/promote-consumer-change.md`: ordered maintainer procedure for promoting reusable
  consumer-local guidance into the kit.
- `runbooks/triage-kit-feedback.md`: ordered maintainer procedure for triaging public kit feedback
  issues.
- `plans/README.md`: repository-level discovery and implementation plans.
- `plans/plan-register.md`: lightweight index of registered Operating Kit repository plans.
- `plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md`:
  completed implementation plan for the released self-contained Go CLI, macOS and Windows binary
  release packs, legacy Python-wheel migration, explicit behavior parity tests, and release-stage
  validation for the Operating Kit bootstrap.
- `plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
  execution evidence and per-epic review log for the self-contained bootstrap implementation.
- `plans/operating-kit-script-asset-roles/operating-kit-script-asset-roles_discovery_doc.md`:
  discovery for script asset role doctrine inside the reusable script asset model.
- `plans/operating-kit-script-asset-roles/operating-kit-script-asset-roles_implementation_doc.md`:
  completed implementation plan for script asset role doctrine, generic workflow composition,
  compact review hooks, and managed/cloud portability guidance.
- `plans/operating-kit-script-asset-roles/operating-kit-script-asset-roles_execution_log.md`:
  execution evidence for the script asset role doctrine source implementation.
- `plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_discovery_doc.md`:
  discovery for reusable runbook-to-script promotion triggers, first-script scaffolding, helper
  rules, domain folder triggers, package/CLI promotion, output contracts, and script testing
  expectations.
- `plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_implementation_doc.md`:
  completed instruction-only implementation plan for managed promotion doctrine, a recipe-to-script
  promotion runbook, current doctrine alignment, resource copies, indexes, and validation.
- `plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_execution_log.md`:
  execution evidence for the runbook-to-script promotion standard source implementation.
- `plans/local-runtime-environment-standard/local-runtime-environment-standard_discovery_doc.md`:
  discovery for ignored repo-local machine/runtime state, the default Python venv path, governed
  machine/user-level tooling, and generic local tooling blocker routing.
- `plans/local-runtime-environment-standard/local-runtime-environment-standard_implementation_doc.md`:
  completed implementation plan for `.codeheart/local/`, `.codeheart/local/envs/python/`, additive
  config/schema behavior, gitignore sync behavior, managed readiness doctrine, tests, and
  packaged-resource mirrors.
- `plans/local-runtime-environment-standard/local-runtime-environment-standard_execution_log.md`:
  execution evidence for the local runtime environment standard source implementation.
- `plans/runtime-materialization-hardening/runtime-materialization-hardening_implementation_doc.md`:
  completed implementation plan for consumer runtime materialization and visible-terminal handoff
  hardening, Operating Kit release, AI Execution adopter release, HQ snapshot proof, and named
  Operating Kit installs.
- `plans/runtime-materialization-hardening/runtime-materialization-hardening_execution_log.md`:
  execution evidence for consumer runtime materialization and visible-terminal handoff hardening.
- `plans/business-docs-placement-clarity/business-docs-placement-clarity_implementation_doc.md`:
  completed implementation plan for clarifying the managed `docs/business/` placement rule so it
  means company or organization business-operating records, not software architecture or
  business-logic documentation.
- `plans/business-docs-placement-clarity/business-docs-placement-clarity_execution_log.md`:
  execution evidence for the business docs placement clarity implementation.
- `plans/operational-recipe-maturity-standard`: HQ-owned completed implementation plan pointer for
  the generic Operating Kit operational recipe maturity standard. Canonical plan:
  `Codeheart-HQ:docs/repo/plans/operational-recipe-maturity-standard/operational-recipe-maturity-standard_implementation_doc.md`.
- `plans/discovery-handoff-gate/discovery-handoff-gate_implementation_doc.md`: active
  implementation plan for the discovery capability-scope handoff gate, `v0.1.13` release, and
  named consumer repository sync.
- `plans/discovery-handoff-gate/discovery-handoff-gate_execution_log.md`: execution evidence and
  review-gate log for the discovery handoff gate release.
- `plans/operation-routing-dispatch-standard/operation-routing-dispatch-standard_discovery_doc.md`:
  discovery for generic route-before-execution-surface doctrine, capability advertisements, route
  registries, route cards, authority hierarchy, and routing-bearing validation hooks.
- `plans/operation-routing-dispatch-standard/operation-routing-dispatch-standard_implementation_doc.md`:
  implementation plan for the managed operation-routing standard, compact root route visibility,
  structure-governance placement rules, planning hooks, validation, and release preparation.
- `plans/operation-routing-dispatch-standard/operation-routing-dispatch-standard_execution_log.md`:
  execution evidence for the selected operation-routing implementation epics and review probes.
- `plans/tooling-environment-readiness/tooling-environment-readiness_discovery_doc.md`:
  discovery for shared Operating Kit tooling and environment-readiness doctrine, missing-tool
  routing, module-owned tool declarations, and local readiness evidence boundaries.
- `plans/tooling-environment-readiness/tooling-environment-readiness_implementation_doc.md`:
  completed implementation plan for one central tooling-readiness route, installed route visibility,
  runbook-authoring and planning hooks, packaged resources, and `0.1.12` release preparation.
- `plans/tooling-environment-readiness/tooling-environment-readiness_execution_log.md`:
  execution evidence, validation, release-readiness evidence, and review-gate log for tooling
  environment readiness.
- `plans/module-extension-state-routing/module-extension-state-routing_discovery_doc.md`:
  discovery for committed non-secret module and extension state under `docs/repo/state/<id>/`.
- `plans/module-extension-state-routing/module-extension-state-routing_implementation_doc.md`:
  active implementation plan for managed structure-governance doctrine, generic agent routing,
  packaged resources, tests, and release preparation for committed module and extension state.
- `plans/module-extension-state-routing/module-extension-state-routing_execution_log.md`:
  execution evidence, validation, release-readiness evidence, and per-epic review log for module
  extension state routing.
- `plans/runbook-authoring-standards/runbook-authoring-standards_discovery_doc.md`: discovery
  for reusable human-facing, agent-facing, and hybrid runbook authoring standards.
- `plans/runbook-authoring-standards/runbook-authoring-standards_implementation_doc.md`:
  draft implementation plan for the managed runbook authoring standard, audience classification,
  compact intention block, workflow hooks, packaged resources, and instruction-only release prep.
- `plans/runbook-authoring-standards/runbook-authoring-standards_execution_log.md`: execution
  evidence, validation, release-readiness evidence, and per-epic review log for the runbook
  authoring standards release.
- `plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_implementation_doc.md`:
  draft implementation plan for refined local and coordination-home plan-register reference shapes,
  repository-qualified ID guidance, canonical pointer examples, and relation ownership doctrine.
- `plans/kit-feedback-intake/kit-feedback-intake_discovery_doc.md`: discovery for a safe
  consumer feedback intake and maintainer backlog workflow.
- `plans/repo-feedback-capture/repo-feedback-capture_discovery_doc.md`: discovery for
  Operating Kit-guided repo feedback capture, demand-driven GitHub Issues setup, prompt
  suppression after decline, classification, privacy, and triage promotion.
- `plans/repo-feedback-capture/repo-feedback-capture_implementation_doc.md`: completed
  implementation plan for managed repo feedback capture runbooks, optional config schema support,
  Codeheart organization membership gating, check-first GitHub issue availability, setup guidance,
  validation, release, and approved consumer sync.
- `plans/repo-feedback-capture/repo-feedback-capture_execution_log.md`: execution evidence for
  the repo feedback capture implementation.
- `plans/plan-register-dirty-target-safety/plan-register-dirty-target-safety_implementation_doc.md`:
  completed implementation plan for dirty target-register compatibility rules before falling back
  to coordination-home pending sync.
- `plans/plan-register-dirty-target-safety/plan-register-dirty-target-safety_execution_log.md`:
  execution evidence for dirty target-register compatibility rules.
- `plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_discovery_doc.md`:
  discovery for a reusable plan-register model and optional multi-repository portfolio
  coordination.
- `plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_implementation_doc.md`:
  implementation plan for managed plan-register doctrine, optional portfolio coordination,
  consumer scaffolds, safe sync behavior, and release validation.
- `plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_execution_log.md`:
  execution evidence and per-epic review log for the portfolio coordination and plan-register
  workflow.
- `plans/plan-register-session-lifecycle-hardening/plan-register-session-lifecycle-hardening_implementation_doc.md`:
  implementation plan for self-contained plan-register session-reference resolution and lifecycle
  grouping hardening.
- `plans/plan-register-session-lifecycle-hardening/plan-register-session-lifecycle-hardening_execution_log.md`:
  execution evidence and per-epic review log for the plan-register session-reference and lifecycle
  hardening release.
- `plans/coordination-home-register-id-namespace/coordination-home-register-id-namespace_implementation_doc.md`:
  active implementation plan for coordination-home register ID namespace doctrine.
- `plans/coordination-home-register-id-namespace/coordination-home-register-id-namespace_execution_log.md`:
  execution evidence and per-epic review log for the coordination-home register ID namespace
  release.
- `plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_implementation_doc.md`:
  active implementation plan for managed planning workflow guidance that preserves feature
  capability and rejects policy-only implementation plans.
- `plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_execution_log.md`:
  execution evidence and per-epic review log for the implementation-planning quality release.
- `plans/kit-feedback-intake/kit-feedback-intake_implementation_doc.md`: completed implementation
  plan for the kit feedback intake workflow.
- `plans/kit-feedback-intake/kit-feedback-intake_execution_log.md`: execution evidence and
  per-epic review log for the kit feedback intake workflow.
