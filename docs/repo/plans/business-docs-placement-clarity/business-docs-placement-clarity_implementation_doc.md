Last updated: 2026-06-26T21:05:06Z (UTC)
Created: 2026-06-26
Status: completed
Completed: 2026-06-26
Execution log: business-docs-placement-clarity_execution_log.md

# Document Header

## Overview

This implementation plan clarifies the managed Operating Kit placement doctrine for
`docs/business/`. The intended change is to make clear that `docs/business/` means company or
organization business-operating records when a consumer repository intentionally stores them. It
does not mean software product architecture, module design, platform solution design, application
business logic, or implementation planning.

This is an instruction-only managed documentation change. It must not change scaffold placement,
sync behavior, CLI behavior, schemas, validators, or consumer ownership boundaries.

## Essential Context Reference Files

| Path | Reason |
| --- | --- |
| `AGENTS.md` | Public-core safety, task routing, and managed content constraints. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer procedure for changing kit docs/source and recording impact. |
| `docs/repo/reference/placement-contract.md` | Source-of-truth for managed docs, consumer-owned docs, and installed consumer paths. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact class definitions and required record for consumer-affecting changes. |
| `components/structure-governance/managed/reference/documentation-structure.md` | Authoring source for the managed placement wording to clarify. |
| `src/codeheart_operating_kit/resources/components/structure-governance/managed/reference/documentation-structure.md` | Packaged resource mirror that must stay in sync with the authoring source. |
| `Codeheart-Automation-Foundry:docs/repo/plans/relational-workspace-view-module/relational-workspace-view-module_discovery_doc.md` | Consumer-originating evidence for the ambiguity: a software/module architecture discovery was initially placed under HQ `docs/business/`. |
| `release-notes.md` | Release-facing note for the instruction-only managed doctrine clarification. |
| `scripts/validate-markdown-headers.py` | Timestamp/header validation for changed Markdown files. |
| `scripts/validate-public-core.py` | Public-core hygiene validation for public managed docs. |
| `tests/test_packaging_resources.py` | Packaged resource mirror coverage. |

## Table Of Contents

- [Section 1 - Foundation](#section-1---foundation)
- [Section 2 - Strategy](#section-2---strategy)
- [Section 3 - Execution Plan](#section-3---execution-plan)
- [Section 4 - Future Planning](#section-4---future-planning)
- [Revision Notes](#revision-notes)

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Clarify the reusable Operating Kit `docs/business/` placement rule so future agents and consumers
do not confuse business-operating records with software product/module architecture or application
business-logic documentation.

Completion is proven when:

- the managed Structure Governance documentation-structure source describes `docs/business/` as
  company/organization business-operating records;
- the same wording explicitly excludes software product architecture, module design, platform
  solution design, application business logic, and implementation planning;
- the packaged resource mirror matches the authoring source;
- release notes record an instruction-only managed doctrine clarification;
- repo indexes and plan register route this implementation plan;
- focused validation passes or residual risk is recorded.

## 1.2 Project And Problem Context

The current managed wording says `docs/business/` is for "business operating docs when the
consumer repository intentionally stores business records." That wording is broadly correct but
can be ambiguous to an agent when the word "business" is used in software contexts, such as
business logic, business apps, business-product architecture, or module/platform design.

A Codeheart-HQ discussion exposed the ambiguity. A discovery about a reusable relational
workspace/view module was initially placed in HQ `docs/business/plans/` because the conversation
used business-operating examples. After clarification, the artifact was moved into the Foundry
repository because its owner is reusable module/platform architecture rather than HQ company
business-operating records.

## 1.3 Current State Analysis

Existing state:

- `components/structure-governance/managed/reference/documentation-structure.md` defines
  top-level consumer documentation areas.
- The packaged mirror in `src/codeheart_operating_kit/resources/` carries the same managed file.
- The current `docs/business/` line does not explicitly exclude software product/module
  architecture or business-logic documentation.
- Consumer repositories may add local clarification, but the reusable placement doctrine lives in
  the Operating Kit.

Target state:

- Managed docs clarify that `docs/business/` is for company/organization operating records such
  as ventures, procurement, meetings, decisions, and business research.
- Managed docs explicitly route software product architecture, module design, platform solution
  design, application business logic, and implementation plans to the owning repo/product/module,
  package, or source-area docs.
- No consumer-owned files are moved or changed by the kit.

Consumer impact record:

- Impact class: `instruction-only change`.
- Affected paths: managed Structure Governance reference and packaged mirror.
- Release notes: required when shipped.
- Consumer action: none.
- Migration/adoption note: none.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

```text
Codeheart-Operating-Kit/
  components/
    structure-governance/
      managed/
        reference/
          documentation-structure.md                         # modify
  src/codeheart_operating_kit/
    resources/
      components/
        structure-governance/
          managed/
            reference/
              documentation-structure.md                     # modify mirror
  docs/
    README.md                                                # modify index
    repo/
      README.md                                              # modify index
      plans/
        README.md                                            # modify index
        plan-register.md                                     # modify index
        business-docs-placement-clarity/
          business-docs-placement-clarity_implementation_doc.md # this plan
  release-notes.md                                           # modify when implemented
```

The implementation should update the authoring source first, copy the same semantic wording into
the packaged resource mirror, then run focused validation.

## 2.2 Open Questions And Assumptions Requiring Clarification

`OQ-1`: Exact wording for the exclusion list.

- `BLOCKER: no`
- Affects: `E1`, `E2`
- Unlocks: final managed wording.
- Recommended default: use public-safe generic terms only: software product architecture, module
  design, platform solution design, application business logic, and implementation planning.

`OQ-2`: Whether to add examples beyond the top-level `docs/business/` bullet.

- `BLOCKER: no`
- Affects: `E1`
- Unlocks: wording depth.
- Recommended default: keep the managed reference concise and avoid consumer-specific examples.

`OQ-3`: Whether planning workflow docs need a parallel clarification.

- `BLOCKER: no`
- Affects: `E3`
- Unlocks: future follow-up.
- Recommended default: do not change planning workflows in this implementation unless validation
  finds another ambiguous managed sentence.

## 2.3 Architectural Decisions With Reasoning

`AD-1`: Treat this as an instruction-only managed docs change.

1. Problem being solved: agents need clearer placement language.
2. Simplest working solution: clarify the existing `docs/business/` bullet.
3. What may change in 6-12 months: additional placement examples may be added if more ambiguous
   consumer cases appear.
4. Rationale: no behavior, scaffold, sync, or schema change is needed.
5. Alternatives considered and why not chosen: a new component or validator is unnecessary.

`AD-2`: Do not change consumer-owned HQ docs from this plan.

1. Problem being solved: Operating Kit should own only reusable generic doctrine.
2. Simplest working solution: update managed source; consumer repos can keep local specializations.
3. What may change in 6-12 months: a future release may include a local wrapper example if
   repeated consumer confusion appears.
4. Rationale: consumer-owned docs should not be overwritten by kit implementation.
5. Alternatives considered and why not chosen: migrating consumer docs would exceed an
   instruction-only doctrine clarification.

`AD-3`: Preserve product/module/source-area ownership language.

1. Problem being solved: software artifacts need an obvious placement route.
2. Simplest working solution: explicitly route software architecture and module design to the
   owning repo/product/module/package/source-area docs.
3. What may change in 6-12 months: product documentation doctrine may become more detailed.
4. Rationale: this matches the existing ownership-first structure doctrine.
5. Alternatives considered and why not chosen: forcing all cross-cutting product planning into
   `docs/repo/` would be too broad.

# Section 3 - Execution Plan

## Epic E1 - Clarify Managed Structure Governance Wording

Outcome: The managed `docs/business/` placement sentence is clearer and public-safe.

Tasks:

- [x] Read the current `docs/business/` bullet in
  `components/structure-governance/managed/reference/documentation-structure.md`.
- [x] Replace the ambiguous sentence with wording that describes company/organization
  business-operating records.
- [x] Add a concise exclusion sentence for software product architecture, module design, platform
  solution design, application business logic, and implementation planning.
- [x] Keep the change generic and public-safe; do not mention private Codeheart examples.

Acceptance criteria:

- [x] The managed wording can reject placing a software/module architecture discovery under
  `docs/business/` solely because examples are business-domain examples.
- [x] The wording still allows consumer repositories that intentionally store company business
  records to use `docs/business/`.

## Epic E2 - Synchronize Packaged Resource Mirror

Outcome: Packaged resources match the authoring source.

Tasks:

- [x] Apply the same managed wording to
  `src/codeheart_operating_kit/resources/components/structure-governance/managed/reference/documentation-structure.md`.
- [x] Confirm no source-only or package-only wording drift remains.

Acceptance criteria:

- [x] `rg "docs/business" components/structure-governance src/codeheart_operating_kit/resources`
  shows equivalent wording in source and packaged mirror.

## Epic E3 - Update Release And Planning Metadata

Outcome: The planned docs change is discoverable and release-facing impact is recorded.

Tasks:

- [x] Add a release-note entry classifying the change as instruction-only.
- [x] Update repo docs/plans indexes if the plan path is not already discoverable.
- [x] Update the plan register entry lifecycle if implementation is activated or completed.

Acceptance criteria:

- [x] Release notes mention that the `docs/business/` placement rule was clarified.
- [x] The implementation plan remains discoverable from `docs/repo/plans/README.md`.

## Epic E4 - Validate

Outcome: Focused validation proves the documentation and packaged-resource change.

Tasks:

- [x] Run Markdown header validation for changed Markdown files.
- [x] Run public-core hygiene validation for changed public docs.
- [x] Run packaged-resource tests that cover managed documentation mirrors.
- [x] Run `git diff --check`.

Suggested commands:

```sh
python scripts/validate-markdown-headers.py
python scripts/validate-public-core.py
pytest tests/test_packaging_resources.py
git diff --check
```

Acceptance criteria:

- [x] Validation passes, or every failed command has a recorded reason and residual risk.

## Epic E5 - Review Gate

Outcome: The final wording is checked against the original ambiguity before release.

Tasks:

- [x] Review the changed wording against the Codeheart-HQ scenario that triggered this plan.
- [x] Confirm the wording would route company business records to `docs/business/` and reusable
  software/module architecture to the owning product/module/source-area docs.
- [x] Record any release-note or follow-up impact.

Acceptance criteria:

- [x] Review confirms the clarification solves the ambiguity without weakening business-docs
  support for consumer repositories that intentionally store company business records.

# Section 4 - Future Planning

- Consider adding a short placement example table if repeated consumer confusion continues.
- Consider a future product-docs placement discovery if module/product/source-area ownership
  remains hard for agents to apply.
- Keep consumer-local specializations in consumer repositories; do not make Operating Kit carry
  private company examples.

# Revision Notes

- 2026-06-26: Drafted implementation plan from Codeheart-HQ placement clarification and Foundry
  discovery relocation.
- 2026-06-26: Activated and completed implementation after explicit user approval; per-epic
  reviewer-agent gates were waived by user request, with main-thread validation and final review
  recorded in the execution log.
- 2026-06-26: Prepared Operating Kit `v0.1.16` release after explicit user request to bump and
  release the implemented clarification.
