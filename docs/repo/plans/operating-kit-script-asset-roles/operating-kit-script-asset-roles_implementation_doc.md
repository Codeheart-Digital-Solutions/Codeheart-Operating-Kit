Last updated: 2026-07-08T14:10:04Z (UTC)
Created: 2026-07-08
Status: completed
Completed: 2026-07-08
Execution log: docs/repo/plans/operating-kit-script-asset-roles/operating-kit-script-asset-roles_execution_log.md

# Document Header

## Operating Kit Script Asset Roles Implementation Plan

Overview: update Operating Kit source doctrine so reusable script assets can be classified as
`primitive script`, `workflow script`, or `helper` without changing the existing L2/L3/L4 maturity
model. The implementation keeps `thin command wrapper` as L3, adds generic workflow-composition
guidance, clarifies contract-based script dependencies, keeps helper promotion conservative, and
adds compact conditional planning, execution, and review checks.

This plan is derived from the approved discovery:
`docs/repo/plans/operating-kit-script-asset-roles/operating-kit-script-asset-roles_discovery_doc.md`.
It is the canonical Operating Kit source plan. It must be executed in this source repository and
later shipped through the normal Operating Kit release/update path, not by hand-editing any
consumer repository's managed `.codeheart/kit/` copy.

Essential context:

| Source | Why It Matters |
| --- | --- |
| `AGENTS.md` | Maintainer routing, public-core safety, and required change routes. |
| `docs/repo/runbooks/change-operating-kit.md` | Ordered maintainer procedure for source docs, managed content, validation, and release-note needs. |
| `docs/repo/reference/placement-contract.md` | Confirms reusable generic doctrine belongs in the Operating Kit and consumer `.codeheart/kit/` copies are managed. |
| `docs/repo/reference/consumer-impact-classification.md` | Classifies this as an instruction-only managed-doc change with release notes required when shipped. |
| `docs/repo/plans/operating-kit-script-asset-roles/operating-kit-script-asset-roles_discovery_doc.md` | Approved decisions `D-001` through `D-011`, patch targets, and implementation capability scope. |
| `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` | Source file for reusable script asset vocabulary, helper rules, first-script scaffolding, output contract, and review flags. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Source file for the L1-L4 maturity model; L2 is reusable script asset and L3 is thin command wrapper. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Planning source runbook that must gain compact role, dependency, and placement checks for script-bearing plans. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Execution source runbook that must verify script roles and boundaries during implementation. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | Planning review source runbook that must flag missing role, dependency, helper, and portability coverage in script-bearing plans. |
| `components/structure-governance/managed/reference/documentation-structure.md` | Structure source reference that currently points to the script-promotion standard for promoted recipe asset placement. |
| `scripts/validate-markdown-headers.py` | Markdown timestamp/header validation for changed docs. |
| `scripts/validate-public-core.py` | Public-core hygiene validation for changed public docs. |
| `tests/test_packaging_resources.py` | Packaged-resource parity validation for managed content changes. Run through `uv` when global pytest is unavailable. |

Table of contents:

- [Section 1 - Foundation](#section-1---foundation)
- [Section 2 - Strategy](#section-2---strategy)
- [Section 3 - Execution Plan](#section-3---execution-plan)
- [Section 4 - Future Planning](#section-4---future-planning)
- [Revision Notes](#revision-notes)

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Implement the approved Operating Kit script asset role doctrine so future agents and maintainers
can distinguish narrow primitive scripts, deterministic workflow scripts, imported helpers, local
ad hoc scripts, and L3 thin command wrappers without adding bureaucracy or domain-specific paths.

Completion is proven when:

- Operating Kit source docs define `primitive script`, `workflow script`, `helper`, and local
  ad hoc script boundaries inside the existing reusable script asset model;
- `thin command wrapper` remains the L3 maturity state and is explicitly not an L2
  `command_wrapper` role;
- workflow scripts are documented as deterministic composition surfaces that can use primitives,
  public script entrypoints, and shared helpers through stable contracts;
- workflow scripts are forbidden from owning approvals, broad routing, ambiguous target
  selection, policy judgment, hidden scope expansion, and hidden target broadening;
- the generic prerequisite-readiness plus operation-primitive pattern is documented without
  hard-coding Microsoft 365, AWS, Azure, GCP, or Foundry examples as doctrine;
- helper placement and shared-helper promotion guidance keeps helpers at the narrowest durable
  owner boundary until real cross-boundary reuse exists;
- `scripts/README.md` guidance remains compact and includes role visibility where it improves
  review;
- planning, execution, and planning-review runbooks contain compact conditional checks for script
  role, dependencies, README index, helper placement, and managed/cloud portability;
- structure-governance guidance points to the enhanced script-promotion standard without
  duplicating its rules;
- validation shows the source docs are internally consistent, no managed consumer copy was edited
  directly, and no new validator, generated registry, Foundry path rule, or domain-specific
  product mechanics were introduced.
- consumer impact is recorded as `instruction-only change`: release notes are required when
  shipped and no consumer migration or manual adoption step is required.

This plan does not authorize:

- hand-editing `.codeheart/kit/` in any consumer repository;
- creating or changing Foundry module scripts;
- changing Foundry runtime wrappers;
- adding Microsoft 365 user-read, email, invoice, or auth scripts;
- creating validators, generated registries, scaffolding generators, or CLI wrappers;
- releasing the Operating Kit or syncing consumer repositories after release.

## 1.2 Project And Problem Context

The current Operating Kit has a good maturity model: L1 structured runbook recipe, L2 reusable
script asset, L3 thin command wrapper, and L4 mature API/tool surface. The gap is inside L2.
`Reusable script asset` currently covers both a narrow deterministic operation and a larger
deterministic multi-phase workflow. That is not wrong, but it makes reviews less precise.

Recent Foundry M365 and invoice-intake discussions exposed the recurring practical cases:

- a stable normal request can benefit from one workflow entrypoint that runs prerequisite
  readiness/access checks and then a narrow operation primitive;
- primitives should not duplicate prerequisite workflow logic;
- workflow scripts can depend on stable helpers and public entrypoints;
- local ad hoc scripts must remain allowed for exploration and unscripted gaps, but they are not
  durable reusable assets;
- future managed-runner and cloud execution needs explicit inputs, structured outputs,
  non-interactive behavior, and clear artifact/state boundaries.

The Operating Kit update should capture these generic rules once, so future Foundry, M365, AI,
document-processing, cloud, and other domain plans do not reinvent the vocabulary.

## 1.3 Current State Analysis

Existing systems, constraints, and problems:

- `runbook-to-script-promotion-standard.md` defines reusable script asset, script entrypoint,
  infrastructure helper, domain helper, domain folder, thin command wrapper, and mature API/tool
  surface, but it does not define primitive or workflow script roles.
- `operational-recipe-maturity.md` correctly keeps L2 and L3 separate, but it does not say that
  primitive, workflow, and helper are roles inside L2.
- `draft-implementation-plan.md` requires script owner, placement boundary, output contract,
  output safety behavior, tests, fixtures, and review criteria, but does not ask for a script
  role or workflow dependency review.
- `execute-implementation-plan.md` verifies reusable-script basics, but does not yet verify
  role, dependency, helper, and portability boundaries.
- `review-planning-document.md` reviews recipe-bearing plans, but does not yet flag missing
  script role, workflow dependency, helper-entrypoint, and managed/cloud portability coverage.
- `documentation-structure.md` points to the existing script-promotion standard for promoted
  recipe assets and can remain a light cross-reference.
- the source files to change already exist under `components/agent-interface/managed/`,
  `components/planning-workflows/managed/`, and `components/structure-governance/managed/`.

Target systems and ownership boundaries:

- Operating Kit source components own generic doctrine and managed runbooks.
- Consumer repositories receive these docs only through the normal Operating Kit update/sync path.
- Domain modules such as M365, AI Execution, and document-processing own concrete script names,
  helper modules, command invocations, blocker classes, and domain validation.
- This implementation changes doctrine and review behavior only. It does not create operational
  scripts.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement in the Operating Kit source tree that owns the component paths mapped by
the managed component source layout.

Expected source file tree:

```text
<operating-kit-source-root>/
  components/
    agent-interface/
      managed/
        reference/
          runbook-to-script-promotion-standard.md      # modify
          operational-recipe-maturity.md               # modify
    planning-workflows/
      managed/
        runbooks/
          draft-implementation-plan.md                 # modify
          execute-implementation-plan.md               # modify
          review-planning-document.md                  # modify
    structure-governance/
      managed/
        reference/
          documentation-structure.md                   # modify
  src/
    codeheart_operating_kit/
      resources/
        components/
          ...                                          # mirror changed managed files
```

Packaged resource mirrors under `src/codeheart_operating_kit/resources/components/` must be kept
in sync with the changed `components/` source docs because consumer installs are built from those
resources.

Installed consumer paths that must not be edited directly:

```text
<consumer-repo>/
  .codeheart/
    kit.lock.yaml                                      # inspect
    kit/
      docs/
        agent-interface/
          reference/
            runbook-to-script-promotion-standard.md    # do not hand-edit
            operational-recipe-maturity.md             # do not hand-edit
        planning-workflows/
          runbooks/
            draft-implementation-plan.md               # do not hand-edit
            execute-implementation-plan.md             # do not hand-edit
            review-planning-document.md                # do not hand-edit
        structure-governance/
          reference/
            documentation-structure.md                 # do not hand-edit
```

## 2.2 Open Questions And Assumptions Requiring Clarification

### OQ-001 - Exact Operating Kit Source Checkout

BLOCKER: no.

Affects: none.

Decision unlocked: none. The canonical plan now lives in the Operating Kit source repository.

Recommended default: execute the plan from this repository root and stop if the component source
paths in Section 2.1 are absent.

### OQ-002 - Release And Consumer Sync Timing

BLOCKER: no.

Affects: none in this plan.

Decision unlocked: determines when the doctrine change reaches consumer repositories.

Recommended default: keep release and consumer sync outside this plan. Execute release and sync
through a separate Operating Kit release/update task after the source-doc implementation is
reviewed and approved.

### OQ-003 - Consumer Impact Classification

BLOCKER: no.

Affects: all epics.

Decision unlocked: release-note and migration-note expectations.

Recommended default: classify this as `instruction-only change`. Release notes are required when
shipped; no consumer migration, repair, sync beyond normal update, or manual adoption note is
required.

## 2.3 Architectural Decisions With Reasoning

### AD-001 - Add Script Roles Inside L2, Not New Maturity Levels

1. Problem being solved: the current L2 layer is too broad for review precision.
2. Simplest working solution: add role vocabulary inside reusable script asset doctrine.
3. What may change in 6-12 months: repeated workflow scripts may justify more examples or a
   helper extraction guide.
4. Rationale: the existing L2/L3/L4 maturity model is correct and should not be disrupted.
5. Alternatives considered and why not chosen: create new maturity levels between L2 and L3.
   Rejected because it would make the model heavier without improving execution.

### AD-002 - Keep Review Guidance Conditional And Compact

1. Problem being solved: review guidance must catch script role and workflow risks without turning
   every plan review into a long checklist.
2. Simplest working solution: add a few conditional bullets under existing reusable-script review
   sections.
3. What may change in 6-12 months: high-risk modules may add local review gates in their own docs.
4. Rationale: most plans do not touch reusable scripts, so the added guidance should activate only
   for script-bearing plans.
5. Alternatives considered and why not chosen: add a separate script architecture review gate.
   Rejected as unnecessary process overhead.

### AD-003 - Keep Examples Generic In Operating Kit Doctrine

1. Problem being solved: M365 motivated the issue, but the Operating Kit must stay provider-neutral.
2. Simplest working solution: use generic terms such as prerequisite readiness, operation
   primitive, managed runner, and cloud orchestration surface.
3. What may change in 6-12 months: domain modules may add domain-specific examples in their own
   module docs.
4. Rationale: generic doctrine stays reusable across local, CI, cloud, Microsoft, and non-Microsoft
   work.
5. Alternatives considered and why not chosen: use concrete M365 auth and email examples in the
   Operating Kit. Rejected because it would make a generic standard look product-specific.

### AD-004 - Use Structure Governance As A Pointer, Not A Duplicate Rulebook

1. Problem being solved: placement guidance spans agent-interface recipe doctrine and structure
   governance.
2. Simplest working solution: put the detailed role and placement rules in the script-promotion
   standard and add only a short cross-reference in structure governance.
3. What may change in 6-12 months: structure governance may gain a broader recipe-asset index if
   multiple standards need one.
4. Rationale: one source of truth reduces drift.
5. Alternatives considered and why not chosen: duplicate the role rules in structure governance.
   Rejected because duplicated doctrine will age poorly.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Source boundary and baseline inventory are confirmed before source edits. | S | None |
| EP-02 | Agent-interface recipe references define script roles, dependencies, workflow composition, helpers, ad hoc scripts, and portability. | M | EP-01 |
| EP-03 | Planning, execution, and review runbooks add compact conditional checks for script-bearing work. | M | EP-02 |
| EP-04 | Structure-governance pointer and validation prove the doctrine is coherent and not domain-specific. | S | EP-02, EP-03 |

## EP-01 - Source Boundary And Baseline Inventory

### A) Epic ID, Title, And Outcome

EP-01 - Source Boundary And Baseline Inventory.

Outcome: the executor is in the Operating Kit source boundary, has mapped every managed consumer
source file in scope, and has recorded the current baseline before edits.

### B) Scope

In scope:

- source-boundary confirmation;
- source path confirmation for all managed component files in scope;
- baseline search for existing terms and affected sections.

Out of scope:

- source doctrine edits;
- installed consumer `.codeheart/kit/` edits;
- release and sync.

Recipe maturity coverage: below recipe threshold. This epic inventories doctrine only.

Routing coverage: applicable because the plan changes generic placement and review guidance.
Fresh low-context routing probe is deferred to EP-04 after the changed docs exist.

### C) Files Touched

- Operating Kit source files are inspected only in this epic.
- `docs/repo/plans/operating-kit-script-asset-roles/operating-kit-script-asset-roles_implementation_doc.md`
  may receive execution notes only when the executing repo uses this plan as the active log
  source.

### D) Acceptance Criteria And Size

Size: S.

Acceptance criteria:

- The executor confirms the active repository is the Operating Kit source repository.
- All six source-component files listed in Section 2.1 exist.
- Baseline searches identify current vocabulary and target insertion sections.
- No installed consumer `.codeheart/kit/` file is modified.
- Consumer impact is recorded as `instruction-only change`.

### E) Dependencies And Critical-Path Notes

No dependencies. This epic blocks all source edits because it prevents accidental consumer managed
copy edits.

### F) Tasks Checklist

- [x] Confirm the active repository root contains `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`.
- [x] Confirm the active repository root contains `components/agent-interface/managed/reference/operational-recipe-maturity.md`.
- [x] Confirm the active repository root contains `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [x] Confirm the active repository root contains `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`.
- [x] Confirm the active repository root contains `components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Confirm the active repository root contains `components/structure-governance/managed/reference/documentation-structure.md`.
- [x] Run `rg -n "reusable script asset|thin command wrapper|infrastructure helper|domain helper|temporary execution scripts|scripts/README.md" components/agent-interface components/planning-workflows components/structure-governance`.
- [x] Record the baseline command output summary in the execution notes.
- [x] Record consumer impact as `instruction-only change` with release notes required when shipped and no consumer migration required.
- [x] Run `git status --short components docs/repo/plans` and verify edits are limited to Operating Kit source docs and planning records.

### G) Implementation Notes

Do not use an installed consumer `.codeheart/kit/` copy as an edit source. Source changes belong
under `components/`.

### H) Open Questions

None.

## EP-02 - Agent-Interface Script Role Doctrine

### A) Epic ID, Title, And Outcome

EP-02 - Agent-Interface Script Role Doctrine.

Outcome: the agent-interface reference docs define the approved script asset roles and boundaries
inside the existing L2 reusable script asset model.

### B) Scope

In scope:

- controlled vocabulary updates;
- primitive, workflow, helper, ad hoc, and thin-command-wrapper boundaries;
- script dependency guidance;
- generic prerequisite-readiness plus operation-primitive workflow pattern;
- owner-boundary, role-folder, README role, helper-promotion, and portability guidance.

Out of scope:

- concrete Foundry, M365, AWS, Azure, GCP, or local-runtime command examples;
- new validators, scaffolds, registries, scripts, wrappers, APIs, or generated metadata.

Recipe maturity coverage: L2 reusable script asset doctrine and L3 thin command wrapper boundary.

Routing coverage: placement guidance changes generic owner-boundary routing language. Fresh
low-context probe is in EP-04.

### C) Files Touched

- `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`
- `components/agent-interface/managed/reference/operational-recipe-maturity.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/operational-recipe-maturity.md`

### D) Acceptance Criteria And Size

Size: M.

Acceptance criteria:

- `runbook-to-script-promotion-standard.md` defines `primitive script`, `workflow script`, and
  `helper` as script asset roles.
- The standard explicitly says `command_wrapper` is not an L2 script role and `thin command
  wrapper` remains L3.
- Workflow scripts are allowed to compose stable phases, public script entrypoints, primitives,
  and helpers through documented contracts.
- Workflow scripts are forbidden from hidden approvals, broad routing, ambiguous target
  selection, policy judgment, hidden scope expansion, and hidden target broadening.
- The generic prerequisite-readiness plus operation-primitive pattern is present without
  product-specific examples.
- Helper placement and shared-helper promotion remain narrow and testable.
- Role folders are optional and domain-first layout remains the default recommendation.
- `scripts/README.md` guidance stays compact and includes a role column where review clarity
  benefits.
- Ad hoc scripts remain allowed for local exploration and are separated from committed reusable
  script assets.
- Portability guidance covers explicit inputs, stable structured outputs, non-interactive
  behavior, idempotency keys, artifact/state boundaries, and non-secret phase evidence.
- `operational-recipe-maturity.md` states that primitive, workflow, and helper are L2 role
  vocabulary, not new maturity states.

### E) Dependencies And Critical-Path Notes

Depends on EP-01. EP-03 should not start until the main doctrine wording is stable enough to point
to.

### F) Tasks Checklist

- [x] Update `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` controlled vocabulary with `script asset role`, `primitive script`, `workflow script`, and `helper`.
- [x] Add explicit wording in `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` that `command_wrapper` is not an L2 role and `thin command wrapper` remains L3.
- [x] Add workflow-script boundary language in `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` covering deterministic composition and forbidden hidden judgment.
- [x] Add contract-based dependency guidance in `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` covering public entrypoints and shared helpers.
- [x] Add generic prerequisite-readiness plus operation-primitive workflow pattern in `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`.
- [x] Add role visibility guidance in `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` covering filenames and compact `scripts/README.md` role index.
- [x] Add helper promotion guidance in `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` covering narrowest durable owner boundary and real cross-boundary reuse.
- [x] Add ad hoc script boundary guidance in `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` covering temporary local use and deliberate promotion triggers.
- [x] Add managed-runner and cloud-orchestration portability guidance in `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`.
- [x] Update `components/agent-interface/managed/reference/operational-recipe-maturity.md` to state script roles live inside L2 and `thin command wrapper` remains L3.
- [x] Sync the matching packaged resource mirror files under `src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/`.
- [x] Run `rg -n "primitive script|workflow script|command_wrapper|thin command wrapper|access_required|readiness_required|cloud orchestration" components/agent-interface/managed/reference`.
- [x] Review the matched sections and confirm no product-specific path is used as generic doctrine.

### G) Implementation Notes

Keep the detailed doctrine in `runbook-to-script-promotion-standard.md`. Use
`operational-recipe-maturity.md` only to preserve the maturity model boundary.

### H) Open Questions

None.

## EP-03 - Planning, Execution, And Review Hooks

### A) Epic ID, Title, And Outcome

EP-03 - Planning, Execution, And Review Hooks.

Outcome: planning, execution, and planning-review runbooks contain compact conditional checks that
activate only when work creates or materially changes reusable script assets.

### B) Scope

In scope:

- implementation-plan drafting guidance;
- implementation execution verification guidance;
- planning document review guidance;
- concise pointers to the enhanced script-promotion standard.

Out of scope:

- a standalone script architecture review gate;
- heavy metadata registry requirements;
- mandatory validation for every doc-only plan.

Recipe maturity coverage: review criteria for L2 reusable script asset planning and execution.

Routing coverage: this epic changes planning/review route behavior for script-bearing plans.
Fresh low-context probe is in EP-04.

### C) Files Touched

- `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
- `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
- `components/planning-workflows/managed/runbooks/review-planning-document.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/review-planning-document.md`

### D) Acceptance Criteria And Size

Size: M.

Acceptance criteria:

- `draft-implementation-plan.md` asks script-bearing plans to name script role, placement
  boundary, dependency shape, helper placement, and portability constraints.
- `execute-implementation-plan.md` verifies implemented scripts match planned role, dependencies,
  README role index, helper placement, and no hidden approval or target broadening.
- `review-planning-document.md` flags missing script role, undocumented workflow dependencies,
  helpers that act like entrypoints, premature command wrappers, and weak portability boundaries.
- New review wording is conditional and compact.
- New wording cites the script-promotion standard rather than duplicating detailed doctrine.

### E) Dependencies And Critical-Path Notes

Depends on EP-02. This epic should reference role names after EP-02 defines them.

### F) Tasks Checklist

- [x] Update `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` reusable-script section with compact script role, dependency, helper, README, and portability planning bullets.
- [x] Update `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` reusable-script verification section with compact role, dependency, helper, README, and portability checks.
- [x] Update `components/planning-workflows/managed/runbooks/review-planning-document.md` recipe-plan review section with compact findings criteria for missing role, undocumented dependencies, helper-entrypoint confusion, premature wrappers, and weak portability.
- [x] Sync the matching packaged resource mirror files under `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/`.
- [x] Run `rg -n "script role|workflow script|helper placement|README|portability|premature wrappers" components/planning-workflows/managed/runbooks`.
- [x] Review the matched sections and confirm the detailed doctrine remains in `runbook-to-script-promotion-standard.md`.

### G) Implementation Notes

Keep each runbook change short. The goal is to route agents to the reference standard, not to copy
the whole standard into every planning workflow.

### H) Open Questions

None.

## EP-04 - Structure Governance Pointer And Validation

### A) Epic ID, Title, And Outcome

EP-04 - Structure Governance Pointer And Validation.

Outcome: structure governance points to the enhanced script role doctrine and validation confirms
the change is generic, coherent, and safe to release later.

### B) Scope

In scope:

- structure-governance cross-reference update;
- source-doc consistency validation;
- fresh low-context routing probe using a generic script-bearing scenario.

Out of scope:

- Operating Kit release, version bump, publication, and consumer sync;
- generated validator additions;
- live external validation.

Recipe maturity coverage: documentation-level validation for L2/L3 doctrine.

Routing coverage: this epic performs the fresh low-context routing probe required because
placement and review-routing guidance changed.

### C) Files Touched

- `components/structure-governance/managed/reference/documentation-structure.md`
- `src/codeheart_operating_kit/resources/components/structure-governance/managed/reference/documentation-structure.md`

### D) Acceptance Criteria And Size

Size: S.

Acceptance criteria:

- `documentation-structure.md` points readers to the enhanced script-promotion standard for script
  asset roles, placement, helper, README, and promotion guidance.
- Structure governance does not duplicate detailed role rules.
- Source validation commands pass or produce documented findings that are fixed before plan
  completion.
- A fresh low-context routing probe confirms a generic request involving repeated script mechanics
  routes to reusable script asset review without selecting a domain-specific path prematurely, or
  the execution log records reviewer-agent unavailability and the strongest practical main-thread
  proxy probe.
- The final source diff contains no installed consumer `.codeheart/kit/` edits.

### E) Dependencies And Critical-Path Notes

Depends on EP-02 and EP-03. This is the final validation epic.

### F) Tasks Checklist

- [x] Update `components/structure-governance/managed/reference/documentation-structure.md` promoted recipe asset section with a short pointer to script asset role and placement guidance.
- [x] Sync the matching packaged resource mirror file under `src/codeheart_operating_kit/resources/components/structure-governance/managed/reference/`.
- [x] Run `rg -n "primitive script|workflow script|script asset role|thin command wrapper|command_wrapper|cloud orchestration|scripts/README.md" components/agent-interface components/planning-workflows components/structure-governance`.
- [x] Run `git diff --check`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `uv run --with pytest python -m pytest tests/test_packaging_resources.py`.
- [x] Run a fresh low-context routing probe by giving a fresh low-context reviewer only this prompt: "A repeated external-service check needs the same readiness step before a narrow resource operation; where should the durable mechanics live?". If reviewer-agent execution is unavailable, record the strongest practical main-thread proxy probe.
- [x] Record routing-probe evidence showing the reviewer or proxy identifies the owner boundary, routes to reusable script asset review, asks any needed ambiguity question, and avoids selecting a product-specific execution surface.
- [x] Run `git status --short components src docs/repo/plans` and verify the source diff contains no installed consumer `.codeheart/kit` edits.
- [x] Record final validation evidence in the execution notes.

### G) Implementation Notes

The routing probe should prove generic owner-boundary thinking. It should not mention M365, AWS,
invoice intake, or Foundry runtime.

### H) Open Questions

None.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Operating Kit release and consumer sync: deferred because this plan is source-doctrine
  implementation only. Trigger: source changes pass review and the user requests release.
- Concrete Foundry module script refactors: deferred because this plan is generic Operating Kit
  doctrine. Trigger: a Foundry module implementation plan creates or changes reusable script
  assets.
- Validators or generated script registries: deferred because the approved approach is lightweight
  README/index guidance. Trigger: repeated stale indexes or broken script role declarations create
  measurable review failures.
- Domain-specific examples: deferred to module docs. Trigger: a domain module needs examples for
  its own script architecture.
- Consumer release notes: deferred to release work. Trigger: the instruction-only source change is
  shipped in an Operating Kit release.

## 4.2 Future Considerations

- A later Operating Kit example appendix could show generic examples for several domains after the
  doctrine has been used in at least two independent modules.
- A later cloud-execution plan may translate workflow script phase boundaries into worker,
  queue, function, or orchestrator contracts.
- A later structure-governance pass may add broader promoted-recipe asset discovery rules if
  multiple standards start to overlap.

# Revision Notes

- 2026-07-08: Initial draft implementation plan created from the approved
  `operating-kit-script-asset-roles` discovery.
- 2026-07-08: Moved canonical plan into the Operating Kit source repository, replaced
  consumer-repo handoff wording, classified impact as instruction-only, and tightened validation
  and routing-probe tasks after planning review.
- 2026-07-08: Activated implementation after explicit user request to activate and implement with
  per-epic subagent review gates.
