Last updated: 2026-06-29T15:06:25Z (UTC)
Created: 2026-06-29
Status: completed
Completed: 2026-06-29
Execution log:
`docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_execution_log.md`

# Document Header

## Runbook-To-Script Promotion Standard Implementation Plan

Overview: Implement the accepted Operating Kit doctrine for promoting fragile, repeated, or
evidence-bearing runbook mechanics into reusable script assets without turning whole runbooks into
premature CLIs, APIs, or broad automation wrappers. The implementation adds a managed reference,
adds an agent-facing promotion runbook, aligns current recipe-maturity and planning guidance away
from the old durable inline-block model, and updates discoverability plus validation coverage.

This plan does not retrofit Foundry, Microsoft 365, AI Execution, AWS Platform, or consumer
repositories. It does not add validators, generated scaffolds, new sync behavior, CLI behavior, or
script templates. It prepares Operating Kit doctrine so downstream repositories can adopt the
standard through their own plans.

Discovery handoff basis: the discovery remains `Status: draft`, but the user approved the
recommendations in the current planning discussion and explicitly requested this implementation
plan. This plan records that approved capability scope instead of silently changing discovery
status.

Essential context:

| Source | Why it matters |
| --- | --- |
| `docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_discovery_doc.md` | Accepted recommendations, boundary model, output contract direction, alignment inventory, and downstream adoption limits. |
| `AGENTS.md` | Public-core safety, managed-content boundary, and maintainer routing. |
| `README.md` | Public repository purpose and consumer-owned boundary. |
| `docs/README.md` | Top-level docs router that must expose this implementation plan. |
| `docs/repo/README.md` | Repository-governance router that must expose this plan bundle. |
| `docs/repo/plans/README.md` | Plan index that must link this implementation plan. |
| `docs/repo/plans/plan-register.md` | Local register that must record this implementation plan and parent discovery relation. |
| `docs/repo/reference/placement-contract.md` | Source placement and managed-content ownership rules. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact class for instruction-only managed doctrine changes. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure for Operating Kit source, docs, and release surfaces. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Current recipe maturity doctrine with the conflicting `Tested script block` state. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Runbook-authoring guidance that must point recipe-bearing sections to the updated promotion model. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Route-before-surface doctrine that should point selected recipes with promotion pressure to the new standard. |
| `components/structure-governance/managed/reference/documentation-structure.md` | Placement doctrine for promoted assets that should cross-reference the new script promotion/scaffolding standard. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | Discovery workflow guidance that must stop steering durable mechanics toward inline blocks. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Implementation-planning workflow that must use the updated maturity vocabulary and coverage requirements. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Execution workflow that must use reusable script assets, tests, output contracts, and review gates. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | Planning review runbook that must detect stale inline-block semantics, missing tests, missing output contracts, hidden approvals, and premature wrappers/APIs. |
| `components/agent-interface/managed/README.md` | Managed agent-interface inventory that must expose the new reference and runbook. |
| `components/agent-interface/managed/kit-readme.md` | Installed fallback inventory that should expose the new standard. |
| `components/agent-interface/component.yaml` | Component manifest that must materialize the new managed reference and runbook during install or sync. |
| `src/codeheart_operating_kit/resources/` | Packaged managed-resource copies that must match component source content. |
| `tests/test_packaging_resources.py` | Source and packaged-resource parity validation. |
| `scripts/validate-markdown-headers.py` | Timestamp/header validator for changed markdown. |
| `scripts/validate-public-core.py` | Public-core hygiene validator. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Implement a managed Operating Kit standard that tells agents and maintainers when to keep work in
runbooks, when to create reusable script assets immediately, how scripts must be placed and tested,
how script output contracts work, and how current managed doctrine should stop presenting durable
inline code as a maturity target.

Implementation completion is proven when:

- a managed reference defines runbook/script/helper/package boundaries;
- reusable script assets are documented as the first durable executable surface when mechanics are
  already fragile, repeated, or evidence-bearing;
- inline code blocks are limited to short invocations, examples, and temporary discovery notes;
- first-script scaffolding is documented with `scripts/README.md`, script entrypoints, tests, and
  fixtures;
- infrastructure-helper, domain-helper, domain-folder, package, and CLI promotion triggers are
  documented;
- the mechanical script output contract is documented with required common fields and
  domain-shaped `data`;
- output-safety fields are framed as emitted-output behavior, not semantic sensitivity judgment;
- current managed recipe-maturity and planning runbooks no longer preserve the old
  `Tested script block` maturity model;
- route-selection and structure-governance references point recipe promotion pressure to the new
  standard rather than stopping at recipe maturity alone;
- review criteria flag long inline implementations, orphan scripts, missing tests, hidden
  approval behavior, unclear output contracts, and premature package/CLI surfaces;
- packaged managed-resource copies, docs indexes, plan register, and validation evidence are
  consistent.

## 1.2 Project And Problem Context

The original Operating Kit direction deliberately avoided premature API surfaces and broad CLIs
because real operational requirements were still emerging. That was the right constraint, but
live module work exposed the downside of prose-heavy runbooks: long inline commands, quoting
rules, request construction, and ad hoc error handling were less tested than scripts and caused
avoidable detours during agent execution.

The accepted resolution is not to turn every operation into a CLI. Runbooks should remain the
policy, UX, approval, routing, and judgment layer. Reusable script assets should own deterministic
mechanics when the risk and repeatability already justify tests, stable outputs, and clear blocker
classification.

This implementation belongs in the Operating Kit because the rule is generic. Foundry and M365
will adopt it later, but the standard should not encode those domains as the default case.

## 1.3 Current State Analysis

Current source state:

- `components/agent-interface/managed/reference/operational-recipe-maturity.md` defines `L2 |
  Tested script block`, which can lead agents to treat inline script bodies as durable assets.
- `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` asks plans to name
  maturity states including `L2 tested script block`.
- `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` and
  `components/planning-workflows/managed/runbooks/discovery-workflow.md` still use executable
  script-block vocabulary.
- `components/agent-interface/managed/reference/runbook-authoring-standard.md` delegates
  recipe-bearing sections to the recipe-maturity reference, but it does not yet point to a
  reusable script asset standard.
- `components/planning-workflows/managed/runbooks/review-planning-document.md` reviews
  recipe-bearing plans, but it does not yet require the reusable-script-asset output, test,
  runbook-caller, and approval-boundary checks from this standard.
- `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` and
  `components/structure-governance/managed/reference/documentation-structure.md` route promoted
  recipe questions through recipe maturity and placement, but they do not yet point promotion
  pressure to the new runbook-to-script promotion standard.
- There is no managed reference dedicated to runbook-to-script promotion, first-script
  scaffolding, helper timing, output contracts, or review flags.
- There is no agent-facing runbook for promoting one recipe into a script without creating a broad
  CLI/API surface.

Target state:

- Agents see one generic promotion standard under the managed agent-interface reference set.
- Agents can follow one managed promotion runbook when a recipe should become a script asset.
- Planning and execution runbooks use the same vocabulary: ordinary guidance, structured runbook
  recipe, reusable script asset, thin command wrapper, and mature API/tool surface.
- Inline code blocks are allowed as short examples, but not as durable execution surfaces.
- Managed-resource copies and package parity remain valid.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement the source managed doctrine first, then align existing managed runbooks, then update
packaged resource copies and indexes. Keep the first implementation instruction-only. Do not add
validators, scaffolds, CLI behavior, or template files in this plan.

Expected source tree:

```text
Codeheart-Operating-Kit/
  docs/
    README.md                                                           # modify
    repo/
      README.md                                                         # modify
      plans/
        README.md                                                       # modify
        plan-register.md                                                # modify
        runbook-to-script-promotion-standard/
          runbook-to-script-promotion-standard_discovery_doc.md          # existing
          runbook-to-script-promotion-standard_implementation_doc.md     # create
  components/
    agent-interface/
      component.yaml                                                   # modify
      managed/
        README.md                                                       # modify
        kit-readme.md                                                   # modify
        reference/
          operation-routing-and-dispatch.md                             # modify
          operational-recipe-maturity.md                                # modify
          runbook-authoring-standard.md                                 # modify
          runbook-to-script-promotion-standard.md                       # create
        runbooks/
          promote-runbook-recipe-to-script.md                           # create
    planning-workflows/
      managed/
        runbooks/
          discovery-workflow.md                                         # modify
          draft-implementation-plan.md                                  # modify
          execute-implementation-plan.md                                # modify
          review-planning-document.md                                   # modify
    structure-governance/
      managed/
        reference/
          documentation-structure.md                                    # modify
  src/
    codeheart_operating_kit/
      resources/
        components/
          agent-interface/
            component.yaml                                             # modify copy
            managed/
              README.md                                                 # modify copy
              kit-readme.md                                             # modify copy
              reference/
                operation-routing-and-dispatch.md                       # modify copy
                operational-recipe-maturity.md                          # modify copy
                runbook-authoring-standard.md                           # modify copy
                runbook-to-script-promotion-standard.md                 # create copy
              runbooks/
                promote-runbook-recipe-to-script.md                     # create copy
          planning-workflows/
            managed/
              runbooks/
                discovery-workflow.md                                   # modify copy
                draft-implementation-plan.md                            # modify copy
                execute-implementation-plan.md                          # modify copy
                review-planning-document.md                             # modify copy
          structure-governance/
            managed/
              reference/
                documentation-structure.md                              # modify copy
  tests/
    test_packaging_resources.py                                         # validate only
  scripts/
    validate-markdown-headers.py                                        # validate only
    validate-public-core.py                                             # validate only
```

Consumer impact classification: `instruction-only change`. The implementation changes managed
instructions, references, and runbooks. It does not add generated paths, consumer scaffolds, sync
behavior, validators, templates, or CLI behavior.

## 2.2 Open Questions And Assumptions Requiring Clarification

### OQ-001 - Should The Output Contract Be One Rigid Envelope?

BLOCKER: no

Affects: EP1, EP2, EP3

Decision unlocked: exact script output guidance.

Recommended default: use required common fields plus domain-shaped `data`. Required fields are
`status`, `script_id`, `mode`, `summary`, `blocker`, `data`, and `output_safety`. External,
sensitive, or state-changing scripts may add `runbook_caller`, `target_summary`,
`action_summary`, and `evidence_summary`.

### OQ-002 - Should The First Implementation Add Validators?

BLOCKER: no

Affects: EP5

Decision unlocked: validation scope.

Recommended default: do not add validators in this implementation. Use managed doctrine, planning
hooks, review criteria, and existing markdown/public-core/resource-parity tests. Add validators
after one or more real adopters prove the structure.

### OQ-003 - Should Historical Plans Be Rewritten?

BLOCKER: no

Affects: EP3

Decision unlocked: alignment scope.

Recommended default: update active managed doctrine and runbooks only. Historical plan documents
may keep old terms as evidence unless an active plan uses them as current implementation
instructions.

### OQ-004 - Should Script Placement Be More Structured From The Start?

BLOCKER: no

Affects: EP1, EP2

Decision unlocked: default script folder guidance.

Recommended default: start flat with `scripts/README.md`, script entrypoints,
`tests/scripts/`, and `tests/fixtures/scripts/`. Add `scripts/lib/` when infrastructure helpers
exist. Add domain folders only when cohesion, safety boundaries, fixtures, ownership, or
navigation pressure justify them.

## 2.3 Architectural Decisions With Reasoning

### AD-001 - Runbooks Remain The Policy And UX Layer

Problem being solved: scripts can make mechanics reliable, but they can also hide approval
boundaries, consequences, and target selection.

Simplest working solution: keep runbooks as the caller and decision layer; scripts execute named
deterministic steps after inputs and approvals are resolved.

What may change in 6-12 months: frequently repeated script families may become thin command
wrappers or mature tool surfaces.

Rationale: this preserves human judgment and approval visibility while reducing fragile manual
mechanics.

Alternatives considered and why not chosen: whole-runbook automation was rejected because most
runbooks mix judgment and mechanics.

### AD-002 - Reusable Script Asset Is The First Durable Executable Surface

Problem being solved: the previous `Tested script block` vocabulary can make long inline code
look durable.

Simplest working solution: document inline code blocks as examples or temporary discovery notes,
and promote fragile or repeated mechanics directly to reusable script assets.

What may change in 6-12 months: some script assets may be grouped into domain folders or command
wrappers when repeated use proves the surface.

Rationale: this removes the confusing middle state while preserving flexibility for short
examples.

Alternatives considered and why not chosen: keeping `Tested script block` as a level was rejected
because it invites drift between runbooks and scripts.

### AD-003 - Output Contracts Are Mechanical Software Contracts

Problem being solved: scripts cannot reliably judge semantic sensitivity or business meaning, but
agents need stable outputs.

Simplest working solution: require small common fields, domain-shaped `data`, stable blocker
classes, and `output_safety` flags that describe emitted-output behavior.

What may change in 6-12 months: a rigid envelope may be added if repeated adoption proves value.

Rationale: this gives agents reliable parsing without asking scripts to make AI-like judgments.

Alternatives considered and why not chosen: a universal JSON envelope was rejected for first
implementation because it would freeze structure too early.

### AD-004 - First Implementation Is Doctrine-Only

Problem being solved: validators and templates can overconstrain an immature standard.

Simplest working solution: ship managed references, runbooks, planning hooks, and review gates
first.

What may change in 6-12 months: validators, helper templates, or script scaffolds may be added
after downstream adoption proves common shapes.

Rationale: the standard is still early, and downstream modules need room to apply it.

Alternatives considered and why not chosen: adding validators now was rejected because it would
make untested assumptions enforceable.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP1 | Managed reference defines the runbook-to-script promotion standard. | M | none |
| EP2 | Managed promotion runbook gives agents a concrete recipe-to-script workflow. | M | EP1 |
| EP3 | Existing managed doctrine and planning runbooks align with reusable script assets. | M | EP1, EP2 |
| EP4 | Packaged resources, indexes, register, and release-note surfaces are discoverable. | M | EP1, EP2, EP3 |
| EP5 | Validation and review gates prove the instruction-only implementation. | S | EP1, EP2, EP3, EP4 |

## EP1

### A) Epic ID, Title, And Outcome

EP1 - Managed Reference For Runbook-To-Script Promotion

Outcome: Operating Kit has one managed reference that defines when deterministic runbook mechanics
become reusable script assets, how first-script scaffolding works, how helpers and folders mature,
and how script outputs stay stable without becoming a broad CLI/API surface.

### B) Scope

In scope:

- Add the managed reference under agent-interface.
- Define runbook, script, helper, domain folder, package, and CLI boundaries.
- Define preferred and deprecated vocabulary.
- Define promotion and non-promotion triggers.
- Define the inline-code-block rule.
- Define first-script scaffolding.
- Define output contract and output-safety behavior.
- Define test expectations and review flags.

Out of scope:

- No validators.
- No script templates.
- No consumer scaffold changes.
- No downstream module adoption.

### C) Files Touched

```text
components/agent-interface/managed/reference/
  runbook-to-script-promotion-standard.md                               # create
src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/
  runbook-to-script-promotion-standard.md                               # create copy
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- The reference contains boundary, trigger, placement, helper, folder, package, output, test, and
  review sections.
- The reference contains a controlled vocabulary section that normalizes preferred terms and
  deprecated terms.
- The reference states that reusable script assets can be the first durable executable surface.
- The reference states that inline code blocks are not durable assets.
- The reference uses required common fields plus domain-shaped `data`.
- The reference treats `output_safety` as emitted-output behavior.
- The packaged resource copy matches the source file.

### E) Dependencies And Critical-Path Notes

No implementation dependency. EP1 is the foundation for EP2 and EP3.

### F) Tasks Checklist

- [x] Create `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` with the boundary model, promotion triggers, non-promotion triggers, inline-code-block rule, first-script scaffolding, helper rules, folder rules, package promotion rules, output contract, test expectations, and review flags.
- [x] Add a controlled vocabulary section to `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` with preferred terms `runbook`, `operational recipe`, `runbook code block`, `reusable script asset`, `script entrypoint`, `infrastructure helper`, `domain helper`, `domain folder`, `thin command wrapper`, and `mature API/tool surface`.
- [x] Add deprecated terms to `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`, including `tested script block` as a maturity state, `executable script block` as a durable asset, `promoted script` when `reusable script asset` is clearer, and `promoted recipe asset` when the asset is specifically a script.
- [x] Add the required common output fields `status`, `script_id`, `mode`, `summary`, `blocker`, `data`, and `output_safety` to `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`.
- [x] Add emitted-output safety guidance to `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`.
- [x] Add first-script placement guidance for `scripts/README.md`, script entrypoints, `tests/scripts/`, and `tests/fixtures/scripts/`.
- [x] Add package and CLI promotion stop conditions to `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`.
- [x] Create the packaged resource copy at `src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`.
- [x] Run `diff -u components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`.

### G) Implementation Notes

Use generic examples only. Do not include tenant names, private paths, provider-specific logs, or
domain-specific command bodies. Keep examples short and public-core safe.

### H) Open Questions

OQ-001 and OQ-004 apply with recommended defaults.

## EP2

### A) Epic ID, Title, And Outcome

EP2 - Managed Runbook For Promoting A Recipe To A Script

Outcome: Agents have an ordered runbook for converting one deterministic runbook recipe into a
reusable script asset while preserving runbook ownership of approvals, target selection, fallback
choice, and consequences.

### B) Scope

In scope:

- Add one managed runbook under agent-interface.
- Include preflight, classification, placement, output-contract, test, runbook-caller, and review
  steps.
- Tell agents to replace long inline implementations with script invocations after the script
  exists.
- Keep non-promotion paths clear for judgment-heavy work.

Out of scope:

- No concrete module script generation.
- No automatic code generation.
- No live external operations.

### C) Files Touched

```text
components/agent-interface/managed/runbooks/
  promote-runbook-recipe-to-script.md                                   # create
src/codeheart_operating_kit/resources/components/agent-interface/managed/runbooks/
  promote-runbook-recipe-to-script.md                                   # create copy
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- The runbook has audience, intent, success, judgment boundary, and stop boundary.
- The runbook has a linear workflow from recipe identification through script validation.
- The runbook requires a runbook caller for each script asset.
- The runbook requires proportional tests from the first script.
- The runbook prohibits hidden approval expansion and hidden target broadening.
- The packaged resource copy matches the source file.

### E) Dependencies And Critical-Path Notes

Depends on EP1. The runbook should cite the new reference created in EP1.

### F) Tasks Checklist

- [x] Create `components/agent-interface/managed/runbooks/promote-runbook-recipe-to-script.md` with a compact intention block, trigger list, inputs, stop conditions, and ordered procedure.
- [x] Add classification steps for runbook-only judgment, structured runbook recipe, reusable script asset, thin command wrapper, and mature API surface.
- [x] Add placement steps for `scripts/README.md`, script entrypoint, `tests/scripts/`, `tests/fixtures/scripts/`, and runbook caller updates.
- [x] Add output-contract steps requiring `status`, `script_id`, `mode`, `summary`, `blocker`, `data`, and `output_safety`.
- [x] Add review steps for long inline implementations, missing tests, hidden approval behavior, target broadening, raw sensitive output, and premature package surfaces.
- [x] Create the packaged resource copy at `src/codeheart_operating_kit/resources/components/agent-interface/managed/runbooks/promote-runbook-recipe-to-script.md`.
- [x] Run `diff -u components/agent-interface/managed/runbooks/promote-runbook-recipe-to-script.md src/codeheart_operating_kit/resources/components/agent-interface/managed/runbooks/promote-runbook-recipe-to-script.md`.

### G) Implementation Notes

Keep the runbook agent-facing. It should guide implementers, not end users. It should call out
that state-changing scripts still need runbook-provided approval references and exact target
scope.

### H) Open Questions

OQ-001 and OQ-002 apply with recommended defaults.

## EP3

### A) Epic ID, Title, And Outcome

EP3 - Align Existing Managed Doctrine And Planning Workflows

Outcome: Current managed Operating Kit doctrine and planning runbooks no longer preserve the
durable inline-code maturity model and instead point to reusable script assets for fragile,
repeated, or evidence-bearing mechanics.

### B) Scope

In scope:

- Update operational recipe maturity vocabulary.
- Update implementation planning guidance.
- Update implementation execution guidance.
- Update discovery workflow guidance.
- Update planning review guidance.
- Update runbook authoring cross-references.
- Update route-before-surface and structure-governance cross-references.

Out of scope:

- Do not rewrite historical plan documents.
- Do not modify installed `.codeheart/kit/` copies by hand.
- Do not apply the standard to downstream modules.

### C) Files Touched

```text
components/agent-interface/managed/reference/
  operation-routing-and-dispatch.md                                     # modify
  operational-recipe-maturity.md                                        # modify
  runbook-authoring-standard.md                                         # modify
components/planning-workflows/managed/runbooks/
  discovery-workflow.md                                                 # modify
  draft-implementation-plan.md                                          # modify
  execute-implementation-plan.md                                        # modify
  review-planning-document.md                                           # modify
components/structure-governance/managed/reference/
  documentation-structure.md                                            # modify
src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/
  operation-routing-and-dispatch.md                                     # modify copy
  operational-recipe-maturity.md                                        # modify copy
  runbook-authoring-standard.md                                         # modify copy
src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/
  discovery-workflow.md                                                 # modify copy
  draft-implementation-plan.md                                          # modify copy
  execute-implementation-plan.md                                        # modify copy
  review-planning-document.md                                           # modify copy
src/codeheart_operating_kit/resources/components/structure-governance/managed/reference/
  documentation-structure.md                                            # modify copy
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- No active managed source file recommends `Tested script block` as a durable maturity state.
- Recipe maturity states align with ordinary guidance, structured runbook recipe, reusable script
  asset, thin command wrapper, and mature API/tool surface.
- Planning guidance tells implementers to select the correct initial execution surface.
- Execution guidance requires tests, output contracts, and runbook callers for reusable script
  assets.
- Planning review guidance can detect non-compliant recipe/script promotion plans, including long
  inline implementations, missing tests, missing output contracts, missing runbook callers, hidden
  approval behavior, target broadening, and premature wrappers or APIs.
- Routing and structure-governance references point from repeatable recipe promotion pressure to
  the new promotion standard without replacing their existing routing and placement roles.
- Runbook authoring guidance points recipe-bearing sections to the updated maturity and promotion
  references.
- Resource copies match source content.

### E) Dependencies And Critical-Path Notes

Depends on EP1 and EP2 because existing docs should cite the new reference and runbook.

### F) Tasks Checklist

- [x] Update `components/agent-interface/managed/reference/operational-recipe-maturity.md` to remove `Tested script block` as a durable maturity state.
- [x] Update `components/agent-interface/managed/reference/operational-recipe-maturity.md` to state that inline code blocks are short examples and temporary discovery notes.
- [x] Update `components/agent-interface/managed/reference/operational-recipe-maturity.md` to point durable executable mechanics to reusable script assets.
- [x] Update `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` to replace old maturity-state examples with the updated model.
- [x] Update `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` to require recipe-bearing epics to state placement, output contract, tests, and runbook caller for reusable script assets.
- [x] Update `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` to align execution checks with reusable script assets, script tests, and output contract proof.
- [x] Update `components/planning-workflows/managed/runbooks/discovery-workflow.md` to ask when fragile mechanics should start as reusable script assets.
- [x] Update `components/planning-workflows/managed/runbooks/review-planning-document.md` to flag long inline implementations, missing reusable script tests, missing output contract, missing runbook caller, hidden approval behavior, target broadening, premature wrappers, and premature APIs.
- [x] Update `components/agent-interface/managed/reference/runbook-authoring-standard.md` to cite the new runbook-to-script promotion reference and promotion runbook.
- [x] Update `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` to point selected recipes with promotion pressure to `runbook-to-script-promotion-standard.md` after route selection.
- [x] Update `components/structure-governance/managed/reference/documentation-structure.md` to point concrete script promotion and scaffolding questions to `runbook-to-script-promotion-standard.md` while preserving placement ownership.
- [x] Update matching resource copies under `src/codeheart_operating_kit/resources/components/`.
- [x] Run `rg -n "Tested script block|tested script block|executable script block|inline script block|script block|embedded executable block|L2 tested|L2\\+ executable|L2 recipes|promoted recipe asset|promoted asset" components src/codeheart_operating_kit/resources`.
- [x] Review each stale-vocabulary search hit and classify it with one label from this set: removed; reworded; short-example reference; historical reference; generic-asset reference.

### G) Implementation Notes

Preserve the idea that maturity states are possible forms, not a mandatory ladder. Do not replace
all references to code examples; the target is durable inline implementations, not short
invocations. `Promoted recipe asset` may remain when the text intentionally speaks generically
across scripts, fixtures, schemas, wrappers, APIs, or other assets. Prefer `reusable script asset`
when the durable asset is specifically a script.

### H) Open Questions

OQ-003 applies with the recommended default.

## EP4

### A) Epic ID, Title, And Outcome

EP4 - Discoverability, Register, And Release-Readiness Surfaces

Outcome: The new managed standard and runbook are discoverable from component indexes, repository
docs, the local plan register, and release-readiness notes without publishing a release.

### B) Scope

In scope:

- Update managed component indexes.
- Update the agent-interface component manifest so future install or sync materializes the new
  managed docs.
- Update repository docs indexes.
- Update local plan register.
- Record release-note and consumer-impact expectations.

Out of scope:

- No public release.
- No tag.
- No GitHub release edits.
- No consumer sync.

### C) Files Touched

```text
components/agent-interface/managed/
  README.md                                                            # modify
  kit-readme.md                                                        # modify
components/agent-interface/
  component.yaml                                                       # modify
src/codeheart_operating_kit/resources/components/agent-interface/
  component.yaml                                                       # modify copy
src/codeheart_operating_kit/resources/components/agent-interface/managed/
  README.md                                                            # modify copy
  kit-readme.md                                                        # modify copy
docs/
  README.md                                                            # modify
  repo/
    README.md                                                          # modify
    plans/
      README.md                                                        # modify
      plan-register.md                                                 # modify
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Agent-interface indexes link the new reference and runbook.
- Agent-interface component manifest lists the new reference and runbook so installed
  `.codeheart/kit/` material can contain advertised files after release or sync.
- Repository docs indexes link the implementation plan.
- Plan register includes `OK-PR-022` as a child of `OK-PR-019`.
- `OK-PR-019` register entry reflects the latest discovery timestamp and child relation.
- Release-readiness note states instruction-only change, release notes required when shipped, and
  no migration required.
- Resource copies match source content.

### E) Dependencies And Critical-Path Notes

Depends on EP1 through EP3 because discoverability must point to final paths.

### F) Tasks Checklist

- [x] Update `components/agent-interface/managed/README.md` to list `reference/runbook-to-script-promotion-standard.md` and `runbooks/promote-runbook-recipe-to-script.md`.
- [x] Update `components/agent-interface/managed/kit-readme.md` to list the new standard and promotion runbook.
- [x] Update `components/agent-interface/component.yaml` and the packaged manifest mirror to list the new reference and runbook.
- [x] Update matching resource copies under `src/codeheart_operating_kit/resources/components/agent-interface/managed/`.
- [x] Update `docs/README.md` to list `docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_implementation_doc.md`.
- [x] Update `docs/repo/README.md` to list the implementation plan and summarize its instruction-only scope.
- [x] Update `docs/repo/plans/README.md` to list the implementation plan under current plans.
- [x] Add `OK-PR-022 - Runbook-To-Script Promotion Standard Implementation` to `docs/repo/plans/plan-register.md`.
- [x] Update `OK-PR-019` in `docs/repo/plans/plan-register.md` with the current discovery timestamp and child relation to `OK-PR-022`.
- [x] Record `instruction-only change`, release-note requirement, and no migration requirement in the `OK-PR-022` coordination note.

### G) Implementation Notes

This plan draft already creates the implementation document and local register entry. During
execution, refresh timestamps and register lifecycle only when material scope changes or the plan
status changes.

### H) Open Questions

None.

## EP5

### A) Epic ID, Title, And Outcome

EP5 - Validation And Review Gate

Outcome: The instruction-only implementation is validated for markdown hygiene, public-core
safety, packaged-resource parity, stale vocabulary, and low-context usability before execution is
marked complete.

### B) Scope

In scope:

- Run existing validators.
- Run stale-vocabulary search.
- Run resource parity tests.
- Perform a fresh low-context review probe for the new promotion route.
- Record release and migration residual risk.

Out of scope:

- No new validator implementation.
- No live consumer install proof.
- No downstream adoption proof.

### C) Files Touched

```text
docs/repo/plans/runbook-to-script-promotion-standard/
  runbook-to-script-promotion-standard_execution_log.md                 # create during execution
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- Markdown timestamp validation passes.
- Public-core hygiene validation passes.
- Packaged-resource parity validation passes.
- Stale vocabulary search has no active managed source hits for the old durable inline-block
  model.
- Stale vocabulary review has an explicit allowlist or remediation note for any remaining generic
  `promoted asset`, short-example, or historical-reference hits.
- A fresh low-context routing probe identifies the promotion reference/runbook for a vague
  request about fragile repeated runbook mechanics.
- Execution log records impact class, release-note requirement, migration requirement, validation
  commands, and residual risk.

### E) Dependencies And Critical-Path Notes

Depends on EP1 through EP4.

### F) Tasks Checklist

- [x] Create `docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_execution_log.md` with initial execution metadata.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 -m pytest tests/test_packaging_resources.py`.
- [x] Run `rg -n "Tested script block|tested script block|executable script block|inline script block|script block|embedded executable block|L2 tested|L2\\+ executable|L2 recipes|promoted recipe asset|promoted asset" components src/codeheart_operating_kit/resources`.
- [x] Record the stale-vocabulary remediation and allowlist result in `docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_execution_log.md`.
- [x] Run a fresh low-context routing probe with a vague request about repeated fragile runbook commands and record whether the probe finds the new reference and runbook before selecting a script surface.
- [x] Record validation results in `docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_execution_log.md`.
- [x] Record release-note requirement and no migration requirement in `docs/repo/plans/runbook-to-script-promotion-standard/runbook-to-script-promotion-standard_execution_log.md`.

### G) Implementation Notes

The stale-vocabulary search may find historical plan documents outside `components/` and packaged
resources. Those are not failures unless they are active managed doctrine or active execution
instructions.

### H) Open Questions

OQ-002 applies with the recommended default.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Add validators for script README presence, runbook caller references, output-contract fields,
  and tests after one or more downstream adopters prove the structure.
- Add helper templates only after repeated script assets show a stable reusable helper shape.
- Apply the standard to Foundry M365 through a separate Foundry-owned adoption plan.
- Apply the standard to AI Execution through a separate Foundry-owned adoption plan when a
  concrete repeated recipe needs promotion.
- Consider a rigid script output envelope only after multiple domains use the required common
  fields successfully.
- Consider command-wrapper promotion only after script assets have repeated stable consumers.

## 4.2 Future Considerations

Downstream adoption should start with an inventory of existing runbook recipes and scripts, then
classify each candidate as runbook-only, structured recipe, reusable script asset,
infrastructure-helper candidate, domain-helper candidate, domain-folder candidate, or package/CLI
candidate. Adoption plans should not rewrite every runbook by default.

Future release planning should include release notes because installed consumers will receive new
managed guidance. No migration note is required unless a later implementation adds validators,
sync behavior, scaffold behavior, generated paths, or consumer action requirements.

# Revision Notes

- 2026-06-29: Initial draft created from the accepted runbook-to-script promotion discovery and
  user-approved defaults for output contract, placement, helper policy, validator scope, adoption
  scope, migration rule, and review gates.
- 2026-06-29: Activated and completed the instruction-only source implementation, including the
  managed promotion reference, promotion runbook, doctrine alignment, resource mirrors, indexes,
  register updates, validation evidence, and low-context routing probe.
- 2026-06-29: Added the agent-interface component manifest update to EP4 after fresh review found
  that source indexes alone would not materialize the new docs into installed `.codeheart/kit/`
  content.
