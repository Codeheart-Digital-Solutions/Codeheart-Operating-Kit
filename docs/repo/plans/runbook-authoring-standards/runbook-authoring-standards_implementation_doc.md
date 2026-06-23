Last updated: 2026-06-23T17:50:15Z (UTC)
Created: 2026-06-23
Status: completed
Completed: 2026-06-23
Execution log: docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_execution_log.md

# Document Header

## Runbook Authoring Standards Implementation Plan

Overview: Add an Operating Kit managed runbook authoring standard that makes new and materially
changed runbooks clearer for their intended audience. The implementation is an instruction-only
Operating Kit release: it creates the standard, routes it from managed docs, integrates it into
planning, execution, and planning-document review workflows, mirrors packaged resources, publishes
the release, and records consumer sync proof.

This plan does not retrofit Foundry M365 runbooks, consumer-local runbooks, or existing Operating
Kit runbooks as a migration project. Future runbooks and materially changed runbooks adopt the
standard through normal Operating Kit guidance.

Essential context:

| Source | Why it matters |
| --- | --- |
| `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_discovery_doc.md` | Accepted discovery, decisions, non-goals, runbook audience model, narrow language-preference rule, and implementation handoff. |
| `docs/repo/plans/runbook-authoring-standards/attachments/runbook-sampling-matrix.md` | Public-safe sample set for testing whether the standard is useful and not overbuilt. |
| `docs/repo/plans/runbook-authoring-standards/attachments/related-operating-kit-feedback.md` | Related local preference and tooling-readiness feedback that must remain mostly out of first implementation scope. |
| `AGENTS.md` | Public-core safety, maintainer routing, and release authority boundaries for this repository. |
| `README.md` | Public repository purpose and consumer-owned boundary. |
| `docs/README.md` | Top-level documentation router that must expose this plan when updated. |
| `docs/repo/README.md` | Repository-governance router that must expose this plan when updated. |
| `docs/repo/plans/README.md` | Plan index that must link this implementation plan. |
| `docs/repo/plans/plan-register.md` | Local plan register that must track this implementation plan. |
| `docs/repo/reference/placement-contract.md` | Managed content placement and installed consumer target rules. |
| `docs/repo/reference/consumer-impact-classification.md` | Consumer-impact class and release-note requirement for managed instruction changes. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure before changing managed content, manifests, templates, tests, or release assets. |
| `docs/repo/runbooks/promote-consumer-change.md` | Public-core promotion rules for reusable guidance derived from consumer feedback. |
| `docs/repo/runbooks/release-operating-kit.md` | Required procedure before publishing a public Operating Kit release. |
| `components/agent-interface/managed/README.md` | Owning route for the new runbook authoring standard. |
| `components/agent-interface/component.yaml` | Source component manifest that must include new managed reference files. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Managed implementation-planning runbook that must apply the standard when plans create or materially change runbooks. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Managed execution runbook that must check the standard before completing runbook-related epics. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | Managed review runbook that must flag vague or user-hostile runbook plans. |
| `components/structure-governance/managed/README.md` | Managed structure route that should point documentation-placement work toward the runbook standard. |
| `src/codeheart_operating_kit/resources/` | Packaged resource mirror installed consumers receive. |
| `tests/test_packaging_resources.py` | Existing parity test that should cover changed source and packaged resources. |
| `release-notes.md` | Consumer-facing release-note surface for the instruction-only release. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Create and ship a reusable Operating Kit standard for authoring durable runbooks so future agents
can distinguish human-facing, agent-facing, hybrid, and maintainer-facing runbooks and apply the
right quality bar.

Implementation completion is proven when:

- `components/agent-interface/managed/reference/runbook-authoring-standard.md` defines audience
  classes, the compact intention block, human-facing flow expectations, agent-facing execution
  expectations, hybrid separation, maintainer-facing scaling, approval and stop boundaries,
  evidence expectations, and review checks;
- the standard includes the narrow language-preference rule: check a visible
  `.codeheart/user/preferences.yaml` for `language`, use it when readable, and ask once only when
  the language value is absent or unreadable;
- managed routes expose the standard from the agent-interface README and installed kit fallback
  inventory, with lightweight cross-links from planning workflows and structure governance;
- implementation planning, execution, and planning-document review runbooks require the standard
  only when runbooks are created or materially changed;
- no broad retrofit of existing Operating Kit, consumer-local, or module runbooks is executed as
  part of this plan;
- source managed files and packaged resource mirrors match byte-for-byte;
- release notes classify the change as an `instruction-only change` with no forced consumer
  migration;
- local validation covers Markdown headers, public-core hygiene, packaged-resource parity,
  focused tests, release manifest checks, and release asset readiness;
- public tag creation, GitHub release publication, release URL verification, and consumer sync
  proof are completed through `docs/repo/runbooks/release-operating-kit.md` after release-readiness
  validation passes.

## 1.2 Project And Problem Context

The immediate product signal came from a Foundry M365 workspace onboarding test. The onboarding
runbook had been hardened technically, but the first live consumer interaction still sounded like
a local technical preflight and asked too many setup questions at once. The user clarified that
human-facing onboarding should explain the intent, guide a nontechnical user step by step, ask
only for information the user can provide, and avoid presenting internal agent mechanics as the
main experience.

The same discussion exposed a broader Operating Kit gap. Some runbooks are intended to guide a
human conversation, some are meant as agent execution recipes, some combine both, and some are
short maintainer procedures. A single vague "runbook" quality standard does not work. Human-facing
runbooks need user-experience flow. Agent-facing runbooks need exact execution paths. Hybrid
runbooks need separation between user copy and operator notes. Maintainer-facing runbooks need
authority, evidence, and validation scaled to blast radius without forcing a verbose template on
every short checklist.

The implementation must stay public-core safe. It may use generalized evidence patterns from the
M365 onboarding feedback and sampled runbooks, but it must not include private tenant identifiers,
customer names, credentials, workspace-specific state, or private business instructions.

## 1.3 Current State Analysis

Current source state:

- `components/agent-interface/managed/reference/onboarding-context-contract.md` already shows a
  strong pattern for separating public prompt text, agent contract, ordered context, and storage
  boundaries.
- `components/agent-interface/managed/runbooks/conduct-first-run-onboarding.md` already shows a
  strong human-facing onboarding flow with ordered questions and explicit approval before writes.
- `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` already contains
  a strong fresh-implementer test for plans, but it does not generalize that standard to
  runbooks.
- `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` already requires
  executing plan intent and review gates, but it does not check runbook audience, intention
  blocks, or user-facing flow.
- `components/planning-workflows/managed/runbooks/review-planning-document.md` already checks
  plan execution readiness, but it does not give reviewers runbook-specific findings to look for.
- `components/structure-governance/managed/README.md` routes documentation placement work, but it
  does not tell agents where to find reusable runbook authoring quality guidance.
- Component manifests and packaged resources mirror managed source docs under
  `src/codeheart_operating_kit/resources/`.
- `tests/test_packaging_resources.py` has a targeted parity list for selected source and packaged
  files; it needs coverage for new or changed managed runbook-standard files.

Target state:

- Agents can find one public-safe runbook authoring standard from the managed agent-interface
  routes.
- A future implementation plan that creates or materially changes a runbook must name the affected
  runbooks, declare audience shape, include the compact intention block, and specify the relevant
  human-facing, agent-facing, hybrid, or maintainer-facing checks.
- A future implementation execution run must refuse to mark a runbook epic complete when the
  standard is missing from changed runbooks.
- A future planning-document review can flag runbook plans that are technical but user-hostile, or
  pleasant but not executable.
- Installed consumers receive the standard through normal Operating Kit sync or update, with no
  forced migration and no automatic edits to consumer-owned runbooks.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement one managed reference under agent-interface, route it from managed inventories, add
planning workflow hooks, mirror packaged resources, and prepare an instruction-only patch release.
Keep the first implementation doctrine-first and review-driven rather than schema-driven.

Expected source tree:

```text
Codeheart-Operating-Kit/
  docs/
    README.md                                                     # modify
    repo/
      README.md                                                   # modify
      plans/
        README.md                                                 # modify
        plan-register.md                                          # modify
        runbook-authoring-standards/
          runbook-authoring-standards_discovery_doc.md            # existing
          runbook-authoring-standards_implementation_doc.md       # create
          attachments/
            runbook-sampling-matrix.md                            # existing
            related-operating-kit-feedback.md                     # existing
  components/
    agent-interface/
      component.yaml                                              # modify
      managed/
        README.md                                                 # modify
        kit-readme.md                                             # modify
        reference/
          runbook-authoring-standard.md                           # create
    planning-workflows/
      component.yaml                                              # modify during release prep
      managed/
        README.md                                                 # modify
        runbooks/
          draft-implementation-plan.md                            # modify
          execute-implementation-plan.md                          # modify
          review-planning-document.md                             # modify
    structure-governance/
      component.yaml                                              # modify during release prep
      managed/
        README.md                                                 # modify
  src/
    codeheart_operating_kit/
      __init__.py                                                 # modify during release prep
      resources/
        manifest.yaml                                             # modify during release prep
        components/
          agent-interface/
            component.yaml                                        # modify mirror
            managed/
              README.md                                           # modify mirror
              kit-readme.md                                       # modify mirror
              reference/
                runbook-authoring-standard.md                     # create mirror
          planning-workflows/
            component.yaml                                        # modify mirror
            managed/
              README.md                                           # modify mirror
              runbooks/
                draft-implementation-plan.md                      # modify mirror
                execute-implementation-plan.md                    # modify mirror
                review-planning-document.md                       # modify mirror
          structure-governance/
            component.yaml                                        # modify mirror
            managed/
              README.md                                           # modify mirror
  tests/
    test_packaging_resources.py                                   # modify
    test_install_metadata.py                                      # release validation
    test_release_assets.py                                        # release validation
    test_sync_check.py                                            # focused managed sync validation
  manifest.yaml                                                   # modify during release prep
  pyproject.toml                                                  # modify during release prep
  release-notes.md                                                # modify
  scripts/
    build-release-assets.py                                       # modify during release prep
    validate-json-schemas.py                                      # validation
    validate-markdown-headers.py                                  # validation
    validate-public-core.py                                       # validation
    validate-release-manifest.py                                  # validation
  bootstrap.md                                                    # modify during release prep
  install.sh                                                      # modify during release prep
  install.ps1                                                     # modify during release prep
  dist/                                                           # generate release assets
```

## 2.2 Open Questions And Assumptions Requiring Clarification

OQ-1 - Target release version

- `BLOCKER: no`
- `Affects: EP-04, EP-05`
- Unlocks exact version surfaces, release notes, release manifest URLs, asset names, and tag name.
- Recommended default: use `v0.1.10` because the repository currently exposes `0.1.9` as the
  package and release-manifest version.

OQ-2 - Root `AGENTS.md` route

- `BLOCKER: no`
- `Affects: EP-02`
- Unlocks whether the root managed block gets a direct route for the runbook authoring standard.
- Recommended default: do not add a root `AGENTS.md` route in the first implementation. Route from
  `.codeheart/kit/docs/agent-interface/README.md`, `.codeheart/kit/README.md`, and planning
  workflow docs so the root block stays concise.

OQ-3 - Sample runbook edits

- `BLOCKER: no`
- `Affects: EP-01, EP-03`
- Unlocks whether any sampled runbooks receive intention blocks during first implementation.
- Recommended default: use the sampling matrix as a quality fixture only. Do not retrofit sampled
  runbooks before the standard proves stable.

## 2.3 Architectural Decisions With Reasoning

AD-1 - Place the standard under agent-interface

1. Problem being solved: Runbook quality spans user conversation, agent execution, and hybrid
   operator guidance, so a placement under only planning or structure would be too narrow.
2. Simplest working solution: Create
   `components/agent-interface/managed/reference/runbook-authoring-standard.md` and link to it
   from planning-workflows and structure-governance docs.
3. What may change in 6-12 months: The standard may split into separate human-facing and
   agent-facing references once repeated usage shows the document is too large.
4. Rationale: Agent-interface already owns first-run onboarding, onboarding context contracts, and
   the relationship between agent behavior and user-facing instructions.
5. Alternatives considered and why not chosen: Structure-governance placement would overemphasize
   document placement. Planning-workflows placement would hide the standard from runbook work that
   is not plan-driven.

AD-2 - Use a compact Markdown intention block before any schema

1. Problem being solved: Runbooks need a visible audience and intent signal, but the shape is not
   mature enough for a validator.
2. Simplest working solution: Require a compact Markdown block with audience, intent, success,
   agent judgment boundary, and stop boundary for new and materially changed durable runbooks.
3. What may change in 6-12 months: A future schema, linter, or validator may check the block once
   repeated examples prove the fields.
4. Rationale: Markdown is visible to humans and agents, works in existing docs, and avoids
   premature automation.
5. Alternatives considered and why not chosen: A YAML front matter schema would be more
   machine-checkable but would add overhead before the standard has practical proof.

AD-3 - Keep the first release instruction-only

1. Problem being solved: Consumers need better runbook guidance without forced rewrites of
   consumer-owned documents.
2. Simplest working solution: Update managed docs, component manifests, packaged resources,
   release notes, tests, release assets, public release metadata, and consumer sync proof without
   changing sync ownership, generated paths, or validators.
3. What may change in 6-12 months: Stable standards may justify a validator and a gradual
   migration guide for selected managed runbooks.
4. Rationale: Instruction-only release minimizes blast radius and matches the discovery decision
   to avoid mass retrofit.
5. Alternatives considered and why not chosen: A retrofit release would mix doctrine creation with
   broad content edits and make feedback on the new standard harder to interpret.

AD-4 - Treat publication as part of active-plan execution

1. Problem being solved: A plan that prepares a release but stops before publication does not
   actually ship the managed standard to consumers.
2. Simplest working solution: Include public tag creation, GitHub release publication, release
   asset verification, and consumer sync proof as the final epic.
3. What may change in 6-12 months: Release automation may combine these steps into one validated
   command.
4. Rationale: Activation of this draft plan is the human release decision. The active plan should
   still stop on release-runbook technical stop conditions, but it should not add a separate human
   approval gate before publication.
5. Alternatives considered and why not chosen: A release-prep-only plan would preserve more manual
   control, but it would not satisfy the goal of shipping the standard.

AD-5 - Integrate the standard into planning, execution, and review

1. Problem being solved: A reference that is not checked during planning and review will not
   reliably change how runbooks are written.
2. Simplest working solution: Add concise hooks to implementation planning, execution, and
   planning-document review docs for runbook creation and material runbook changes.
3. What may change in 6-12 months: A future dedicated runbook-review checklist may become useful
   after several plans use the standard.
4. Rationale: The existing planning workflows are the main path by which durable managed and
   repository runbooks are created.
5. Alternatives considered and why not chosen: A standalone runbook-audit runbook would be useful
   later, but it would not influence normal implementation plans as directly.

AD-6 - Keep local preferences narrow

1. Problem being solved: Human-facing onboarding should not repeatedly ask for language when a
   reliable local preference is already visible.
2. Simplest working solution: The standard only covers `language` in
   `.codeheart/user/preferences.yaml`: use the value when present and readable, ask once when it is
   missing or unreadable, and continue the current flow in the selected language.
3. What may change in 6-12 months: Operating Kit may define a broader local preference contract
   for tone, accessibility, notification, tooling, and workspace defaults.
4. Rationale: Language has immediate user-experience value and a low implementation burden.
5. Alternatives considered and why not chosen: A broad preference system would be useful but
   would distract from runbook authoring standards and needs separate discovery.

AD-7 - Keep environment tooling readiness out of this implementation

1. Problem being solved: Module onboarding often needs tools such as PowerShell, but that is a
   shared environment-readiness concern rather than a runbook-authoring standard.
2. Simplest working solution: Preserve the tooling-readiness topic in the related feedback
   attachment and do not implement new readiness files in this plan.
3. What may change in 6-12 months: Operating Kit may add a tooling register, environment preflight
   runbooks, and reusable install guidance after module usage patterns stabilize.
4. Rationale: Mixing tooling inventory with runbook quality would create a larger release before
   either surface is proven.
5. Alternatives considered and why not chosen: Adding a tooling register now would satisfy a real
   future need but lacks enough discovery and would expand the first implementation beyond the
   accepted handoff.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-00 | Plan is activated and execution evidence starts before implementation work. | S | None |
| EP-01 | Managed runbook authoring standard exists and passes fixture review. | M | EP-00 |
| EP-02 | Standard is discoverable through managed routes and packaged resources. | M | EP-01 |
| EP-03 | Planning, execution, and review workflows enforce the standard for changed runbooks. | M | EP-01 |
| EP-04 | Instruction-only release prep, register updates, and validation are complete. | M | EP-01, EP-02, EP-03 |
| EP-05 | Public release publication and consumer sync proof are complete. | M | EP-04 |

## EP-00 - Activation And Execution Evidence

### A) Epic ID, Title, And Outcome

EP-00 - Activation And Execution Evidence

Outcome: The implementation run starts with active lifecycle metadata and an execution log before
any managed source, package, release, or register implementation work begins.

### B) Scope

In scope:

- Set the plan lifecycle to active when execution is approved.
- Create the sibling execution log before EP-01 starts.
- Record activation state in local and HQ coordination registers.
- Run a metadata validation pass after activation edits.

Out of scope:

- Editing managed component content.
- Building release assets.
- Publishing the release.

### C) Files Touched

- `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_implementation_doc.md`
- `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_execution_log.md`
- `docs/repo/plans/plan-register.md`
- `<coordination-home>/docs/repo/plans/plan-register.md`

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- The implementation plan is `Status: active` before EP-01 work begins.
- The execution log exists beside the implementation plan.
- The local and HQ coordination registers reflect active implementation state.
- Markdown timestamp validation passes after activation edits.

### E) Dependencies And Critical-Path Notes

This epic has no dependency.

Critical path: Execution evidence must exist before managed source work begins so per-epic
validation, review rounds, release evidence, and publication proof are recorded as they happen.

### F) Tasks Checklist

- [x] Set `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_implementation_doc.md` to `Status: active`.
- [x] Create `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_execution_log.md` with header, plan path, mode, status, and activation summary.
- [x] Update `docs/repo/plans/plan-register.md` with active lifecycle state for `OK-PR-007`.
- [x] Update the Codeheart-HQ coordination register with active lifecycle state for `CODEHEART-OPERATING-KIT-PR-007`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Record activation validation evidence in the execution log.

### G) Implementation Notes

Activation of this draft plan is the release decision for the full plan path, including the final
publication epic. The implementation still obeys release-runbook technical stop conditions such as
failed validation, checksum mismatch, unreproducible assets, or tag mismatch.

### H) Open Questions

- None.

## EP-01 - Managed Runbook Authoring Standard

### A) Epic ID, Title, And Outcome

EP-01 - Managed Runbook Authoring Standard

Outcome: The Operating Kit contains one public-safe managed reference that defines the reusable
runbook audience classes, compact intention block, quality requirements, narrow language rule, and
review checklist for new and materially changed durable runbooks.

### B) Scope

In scope:

- Create the new agent-interface managed reference.
- Define audience classes: `human-facing`, `agent-facing`, `hybrid`, and `maintainer-facing`.
- Define the compact intention block and placement near the top of durable runbooks.
- Define human-facing flow quality, agent-facing execution quality, hybrid separation, and
  maintainer-facing scaling.
- Include approval, stop, evidence, and validation expectations.
- Include the narrow language-preference rule.
- Use the sampling matrix as a fixture for standard quality.

Out of scope:

- Retrofitting existing runbooks.
- Creating a validator, schema, or linter.
- Creating a broad preference system.
- Creating a tooling or environment-readiness register.

### C) Files Touched

- `components/agent-interface/managed/reference/runbook-authoring-standard.md`
- `docs/repo/plans/runbook-authoring-standards/attachments/runbook-sampling-matrix.md`

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- The new reference is public-safe and contains no consumer-private identifiers.
- The audience classes are defined clearly enough for agents and reviewers to apply them.
- The compact intention block contains exactly the fields accepted in discovery.
- Human-facing, agent-facing, hybrid, and maintainer-facing sections each state required behavior
  and scaling boundaries.
- The narrow language rule is present and limited to visible `.codeheart/user/preferences.yaml`
  `language` handling.
- The sampling matrix can be reviewed against the standard without requiring source runbook
  retrofits.

### E) Dependencies And Critical-Path Notes

Depends on EP-00. This epic sets the doctrine other epics route and enforce.

Critical path: The standard must land before manifests, packaged resources, and workflow hooks can
reference it.

### F) Tasks Checklist

- [x] Create `components/agent-interface/managed/reference/runbook-authoring-standard.md` with timestamp, title, audience model, compact intention block, and adoption boundary.
- [x] Add human-facing requirements for user intent, step-by-step flow, critical turn wording, question pacing, visible help text, approval wording, and local language preference handling.
- [x] Add agent-facing requirements for source of truth, inputs, preconditions, tool lane, ordered execution path, stop conditions, evidence, validation, and cleanup boundaries.
- [x] Add hybrid requirements for separated user-facing flow, operator notes, execution path, stop conditions, evidence, and validation.
- [x] Add maintainer-facing requirements for authority, scope, evidence, validation, and scaling by blast radius.
- [x] Add a review checklist section that agents can apply to new and materially changed runbooks.
- [x] Review `docs/repo/plans/runbook-authoring-standards/attachments/runbook-sampling-matrix.md` against the reference and adjust the reference wording for practical coverage.
- [x] Scan the new reference for public-core safety excluding private consumer, tenant, account, workspace, credential, and raw operational details.

### G) Implementation Notes

Keep the standard concise enough to be usable. It should be specific about behavior without
turning every runbook into a large template. The compact intention block is required for new and
materially changed durable runbooks; it is not a migration demand for unchanged docs.

The language rule is deliberately narrow. It tells a human-facing or hybrid runbook to check for a
visible `language` value in `.codeheart/user/preferences.yaml` and continue in that language when
available. It does not define the full preference file schema.

### H) Open Questions

- OQ-3 is relevant and defaults to fixture review without source runbook retrofit.

## EP-02 - Managed Routing And Packaged Resources

### A) Epic ID, Title, And Outcome

EP-02 - Managed Routing And Packaged Resources

Outcome: Installed consumers can discover the runbook authoring standard through managed
agent-interface routes and fallback inventory, and the source managed docs match their packaged
resource mirrors.

### B) Scope

In scope:

- Route the standard from the agent-interface README.
- Route the standard from the installed kit fallback inventory.
- Add a lightweight structure-governance cross-link.
- Add the new reference to the agent-interface component manifest.
- Mirror all changed managed source files into packaged resources.
- Expand packaged-resource parity tests for the changed files.

Out of scope:

- Adding a new root `AGENTS.md` managed-block route.
- Adding a new component.
- Changing generated consumer paths.
- Changing scaffold ownership.

### C) Files Touched

- `components/agent-interface/managed/README.md`
- `components/agent-interface/managed/kit-readme.md`
- `components/agent-interface/component.yaml`
- `components/structure-governance/managed/README.md`
- `components/structure-governance/component.yaml`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/README.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/kit-readme.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/component.yaml`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/runbook-authoring-standard.md`
- `src/codeheart_operating_kit/resources/components/structure-governance/managed/README.md`
- `src/codeheart_operating_kit/resources/components/structure-governance/component.yaml`
- `tests/test_packaging_resources.py`

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `components/agent-interface/managed/README.md` lists the new standard in managed routes.
- `components/agent-interface/managed/kit-readme.md` lists the installed path to the standard.
- `components/agent-interface/component.yaml` installs the standard to
  `.codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md`.
- `components/structure-governance/managed/README.md` points runbook documentation placement work
  to the agent-interface standard.
- Packaged resource mirrors exist for every changed managed source file.
- `tests/test_packaging_resources.py` asserts parity for the new and changed managed files.
- A temp consumer install or packaged fallback test proves the installed target exists at
  `.codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md`.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 because routes and component manifests must point to a real file.

Critical path: The component manifest update must happen before sync or packaged fallback testing
can prove installed consumers receive the new standard.

### F) Tasks Checklist

- [x] Add the runbook authoring standard route to `components/agent-interface/managed/README.md`.
- [x] Add the installed fallback inventory route to `components/agent-interface/managed/kit-readme.md`.
- [x] Add `components/agent-interface/managed/reference/runbook-authoring-standard.md` to `components/agent-interface/component.yaml` with managed ownership and installed target path.
- [x] Add a concise cross-link in `components/structure-governance/managed/README.md` from runbook documentation work to the agent-interface standard.
- [x] Mirror changed agent-interface source files under `src/codeheart_operating_kit/resources/components/agent-interface/`.
- [x] Mirror changed structure-governance source files under `src/codeheart_operating_kit/resources/components/structure-governance/`.
- [x] Add parity assertions for the new reference and changed managed files in `tests/test_packaging_resources.py`.
- [x] Add an installed-consumer assertion for `.codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md` in `tests/test_packaging_resources.py`.
- [x] Run a targeted packaged-resource parity check with `python3 -m pytest tests/test_packaging_resources.py`.

### G) Implementation Notes

Do not add a direct root `AGENTS.md` route in this implementation. The root managed block should
stay short and route users to the agent-interface area and kit fallback inventory.

Component versions are handled in EP-04 release prep so source and packaged manifests can be
updated consistently with the release manifest.

### H) Open Questions

- OQ-2 is relevant and defaults to no root `AGENTS.md` route.

## EP-03 - Planning, Execution, And Review Integration

### A) Epic ID, Title, And Outcome

EP-03 - Planning, Execution, And Review Integration

Outcome: Operating Kit planning, execution, and planning-document review workflows require agents
to apply the runbook authoring standard whenever a plan creates or materially changes runbooks.

### B) Scope

In scope:

- Update implementation-planning guidance to identify affected runbooks, audience class, intention
  block, relevant audience checks, and no-retrofit boundary.
- Update implementation-execution guidance so runbook epics cannot complete without the planned
  runbook standard checks.
- Update planning-document review guidance so reviewers can flag weak human-facing, agent-facing,
  hybrid, and maintainer-facing runbook plans.
- Update planning-workflows README route text as needed.
- Mirror changed planning-workflows files into packaged resources.

Out of scope:

- Creating a separate runbook-review runbook.
- Running a repository-wide audit of all existing runbooks.
- Changing implementation plan file format beyond runbook-change checks.

### C) Files Touched

- `components/planning-workflows/managed/README.md`
- `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
- `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
- `components/planning-workflows/managed/runbooks/review-planning-document.md`
- `components/planning-workflows/component.yaml`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/README.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/review-planning-document.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml`
- `tests/test_packaging_resources.py`

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `draft-implementation-plan.md` tells implementation planners to add runbook-specific scope,
  audience class, intention block, and relevant audience checks when plans create or materially
  change runbooks.
- `execute-implementation-plan.md` tells executors to verify changed runbooks against the standard
  before completing runbook-related epics.
- `review-planning-document.md` tells reviewers to flag missing runbook audience class, missing
  intention block, poor human-facing flow, vague agent execution path, missing hybrid separation,
  and overbroad retrofit scope.
- Planning workflow references point to the agent-interface standard without duplicating the full
  standard.
- Packaged resource mirrors match source files.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 because planning workflow hooks need a stable target reference.

Critical path: This epic turns the standard from passive documentation into a checked planning and
review expectation.

### F) Tasks Checklist

- [x] Add runbook-change coverage rules to `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [x] Add runbook-standard completion checks to `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`.
- [x] Add runbook-specific review findings to `components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Add a concise route to the standard in `components/planning-workflows/managed/README.md`.
- [x] Mirror changed planning-workflows files under `src/codeheart_operating_kit/resources/components/planning-workflows/`.
- [x] Add parity assertions for changed planning-workflows files in `tests/test_packaging_resources.py`.
- [x] Run a targeted workflow-doc scan that confirms each changed planning runbook references `runbook-authoring-standard.md`.
- [x] Run `uv run --with pytest python -m pytest tests/test_packaging_resources.py` after planning-workflow mirrors are updated.

### G) Implementation Notes

Keep workflow hooks compact. The planning workflows should route to the standard and state the
required checks, not duplicate the entire standard in three places.

Execution guidance must keep the scope boundary clear: executors apply the standard to planned new
or materially changed runbooks, not every runbook discovered while working.

### H) Open Questions

- None.

## EP-04 - Release Prep, Registers, And Validation

### A) Epic ID, Title, And Outcome

EP-04 - Release Prep, Registers, And Validation

Outcome: The instruction-only Operating Kit release is prepared, indexed, registered, validated,
and ready for immediate publication in EP-05.

### B) Scope

In scope:

- Update release notes and consumer-impact wording.
- Bump modified component versions and package version for the next patch release.
- Regenerate release manifest and packaged release asset metadata.
- Update local plan indexes and plan registers.
- Update the HQ coordination register.
- Run repository validation and focused tests.
- Record release-readiness evidence.

Out of scope:

- Creating the public Git tag before EP-05.
- Publishing the GitHub release before EP-05.
- Running consumer sync proof before EP-05.

### C) Files Touched

- `release-notes.md`
- `components/agent-interface/component.yaml`
- `components/planning-workflows/component.yaml`
- `components/structure-governance/component.yaml`
- `src/codeheart_operating_kit/resources/components/agent-interface/component.yaml`
- `src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml`
- `src/codeheart_operating_kit/resources/components/structure-governance/component.yaml`
- `pyproject.toml`
- `src/codeheart_operating_kit/__init__.py`
- `scripts/build-release-assets.py`
- `manifest.yaml`
- `src/codeheart_operating_kit/resources/manifest.yaml`
- `bootstrap.md`
- `install.sh`
- `install.ps1`
- `docs/README.md`
- `docs/repo/README.md`
- `docs/repo/plans/README.md`
- `docs/repo/plans/plan-register.md`
- `<coordination-home>/docs/repo/plans/plan-register.md`
- `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_implementation_doc.md`

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `release-notes.md` includes the runbook authoring standard under the next patch release with
  `instruction-only change` impact and no forced consumer migration.
- Modified component manifests and package version surfaces use the confirmed release version.
- Release manifest and installed resource manifest agree with package version, component versions,
  asset names, release URLs, and component checksums. The root release manifest records real
  downloadable asset hashes; the installed resource manifest keeps the existing zero-placeholder
  downloadable asset hash pattern.
- Local docs indexes and plan register link the implementation plan and execution log after
  activation.
- HQ coordination register has a matching Codeheart-Operating-Kit entry.
- Validation passes for Markdown headers, public-core hygiene, JSON schemas, release manifest,
  packaged-resource parity, and focused CLI tests.
- Release artifacts are ready for immediate publication in EP-05.

### E) Dependencies And Critical-Path Notes

Depends on EP-01, EP-02, and EP-03.

Critical path: Release prep must occur after all managed source and packaged resources are final
so checksums and manifests represent the validated content.

### F) Tasks Checklist

- [x] Confirm release target `0.1.10` before editing version surfaces.
- [x] Add a `v0.1.10` section to `release-notes.md` with included changes, consumer impact, and validation summary.
- [x] Update agent-interface, planning-workflows, and structure-governance component versions in source component manifests.
- [x] Mirror component version changes into packaged component manifests.
- [x] Update `pyproject.toml`, `src/codeheart_operating_kit/__init__.py`, and `scripts/build-release-assets.py` to release version `0.1.10`.
- [x] Build release assets with `uv run --with pip --with setuptools python scripts/build-release-assets.py --version 0.1.10 --output-dir dist`.
- [x] Update `manifest.yaml` and `src/codeheart_operating_kit/resources/manifest.yaml` with release version, URLs, root release asset checksums, packaged zero-placeholder asset hashes, and component checksums.
- [x] Update `bootstrap.md`, `install.sh`, and `install.ps1` release URLs and defaults from the generated release manifest.
- [x] Update `docs/README.md`, `docs/repo/README.md`, and `docs/repo/plans/README.md` with the implementation plan and execution log routes.
- [x] Update `docs/repo/plans/plan-register.md` with the local implementation-plan entry and lifecycle changes.
- [x] Update the Codeheart-HQ coordination register with the matching Codeheart-Operating-Kit implementation-plan entry.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 scripts/validate-json-schemas.py`.
- [x] Run `python3 scripts/validate-release-manifest.py`.
- [x] Run `uv run --with pytest --with pip --with setuptools python -m pytest tests/test_packaging_resources.py tests/test_install_metadata.py tests/test_release_assets.py tests/test_sync_check.py`.
- [x] Run `git diff --check`.
- [x] Record validation evidence, release-readiness evidence, and residual risk in the execution log.

### G) Implementation Notes

This epic prepares the release artifacts. EP-05 publishes them through
`docs/repo/runbooks/release-operating-kit.md` without adding a separate human approval gate.

The HQ coordination register path is outside this repository. Preserve unrelated HQ changes and
only add the matching coordination entry for this implementation plan.

### H) Open Questions

- OQ-1 is relevant and defaults to `0.1.10`.

## EP-05 - Public Release Publication And Consumer Sync Proof

### A) Epic ID, Title, And Outcome

EP-05 - Public Release Publication And Consumer Sync Proof

Outcome: The validated Operating Kit release is published publicly, release assets are verified
from their published URLs, and at least one consumer sync proof confirms the new runbook authoring
standard is installed through the normal update path.

### B) Scope

In scope:

- Create the public Git tag from the validated commit.
- Publish the GitHub release with release notes, manifest, installers, archives, and checksums.
- Verify published release URLs and checksums.
- Run an isolated consumer install or update proof from the published release.
- Run first configured consumer update-check, sync, and check proof when the configured consumer
  repository is available and safe.
- Update the implementation plan, execution log, local register, and HQ coordination register to
  completed state.

Out of scope:

- Changing release content after the validated release commit.
- Rewriting consumer-owned runbooks during sync proof.
- Publishing a second release from the same implementation plan.

### C) Files Touched

- `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_implementation_doc.md`
- `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_execution_log.md`
- `docs/repo/plans/plan-register.md`
- `<coordination-home>/docs/repo/plans/plan-register.md`

External state touched:

- Git tag `v0.1.10`
- GitHub release `v0.1.10`
- Published release assets and checksums
- Consumer sync proof workspace

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Git tag `v0.1.10` exists at the validated commit.
- GitHub release `v0.1.10` exists with release notes, release manifest, installers, archives, and
  checksums.
- Published release URLs in `manifest.yaml` resolve to assets matching recorded checksums.
- Consumer proof shows the installed target
  `.codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md`.
- The implementation plan is `Status: completed` with `Completed: YYYY-MM-DD`.
- The local and HQ coordination registers reflect completed implementation state.

### E) Dependencies And Critical-Path Notes

Depends on EP-04.

Critical path: Publication must use the exact commit validated in EP-04. If validation fails after
asset generation, fix the source, rebuild release assets, regenerate manifests, repeat EP-04
validation, then publish from the newly validated commit.

### F) Tasks Checklist

- [x] Confirm the current commit matches the validated release commit recorded in the execution log.
- [x] Create Git tag `v0.1.10` from the validated commit.
- [x] Publish GitHub release `v0.1.10` with `release-notes.md`, `manifest.yaml`, installers, release archives, and checksum files.
- [x] Verify published release URLs in `manifest.yaml` return assets matching recorded checksums.
- [x] Run isolated consumer install proof from the published `bootstrap.md` release URL.
- [x] Verify isolated consumer install includes `.codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md`.
- [x] Run first configured consumer `codeheart-operating-kit update-check`, `codeheart-operating-kit sync`, and `codeheart-operating-kit check` proof from the published release.
- [x] Verify first configured consumer install includes `.codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md`.
- [x] Record tag, release URL, asset checksums, isolated consumer proof, configured consumer proof, and residual risk in the execution log.
- [x] Set `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_implementation_doc.md` to `Status: completed` with `Completed: 2026-06-23`.
- [x] Update `docs/repo/plans/plan-register.md` with completed lifecycle state for `OK-PR-007`.
- [x] Update the Codeheart-HQ coordination register with completed lifecycle state for `CODEHEART-OPERATING-KIT-PR-007`.

### G) Implementation Notes

This epic deliberately includes publication. Do not add a separate human approval pause between
EP-04 validation and EP-05 publication during active execution. Technical release-runbook stop
conditions still apply.

The configured consumer proof should preserve unrelated consumer work. If the configured consumer
repository is unavailable or unsafe to edit, record the reason and complete the isolated consumer
proof before closing the release.

### H) Open Questions

- OQ-1 is relevant and defaults to `0.1.10`.

# Section 4 - Future Planning

Future work intentionally left outside this implementation:

- Build a validator or linter for compact runbook intention blocks after the standard has repeated
  usage examples.
- Audit and selectively retrofit high-value managed Operating Kit runbooks once the standard is
  proven practical.
- Apply the standard to Foundry M365 onboarding and operation runbooks through a Foundry-specific
  plan, not through this Operating Kit implementation.
- Define a broader local preference contract beyond language.
- Discover and implement shared environment-readiness or tooling-register guidance for module
  onboarding.
- Continue the separate module-extension state routing discovery into its own implementation
  plan.

# Revision Notes

- 2026-06-23: Initial implementation plan drafted from the accepted runbook authoring standards
  discovery, sampling attachment, and related Operating Kit feedback attachment.
- 2026-06-23: Added activation evidence, public release publication, and consumer sync proof to
  the execution path; removed the separate publication approval gate from the active plan.
