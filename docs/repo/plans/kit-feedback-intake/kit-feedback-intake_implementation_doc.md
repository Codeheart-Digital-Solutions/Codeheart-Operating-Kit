Last updated: 2026-06-17T20:21:29Z (UTC)
Created: 2026-06-17
Status: completed
Completed: 2026-06-17
Execution log: kit-feedback-intake_execution_log.md

# Document Header

## Overview

This implementation plan turns the approved kit feedback intake discovery into a shippable v1 for
Codeheart Operating Kit. It adds a public GitHub issue intake surface, managed consumer guidance,
maintainer triage guidance, label taxonomy, release notes, and validation without adding a CLI
feedback command, project board, private security disclosure channel, or synchronized consumer
draft scaffold.

The first execution step is a repository-governance preflight. Normal implementation must not begin
until the implementer confirms GitHub Issues, issue forms, and label creation are available for the
public `Codeheart-Operating-Kit` repository.

## Essential Context Reference Files

| Path | Reason |
| --- | --- |
| `AGENTS.md` | Maintainer bootstrap, public-core safety, and change-safety rules. |
| `docs/repo/plans/kit-feedback-intake/kit-feedback-intake_discovery_doc.md` | Approved decisions, requirements, risks, and deferred scope. |
| `docs/repo/reference/placement-contract.md` | Ownership rules for managed kit docs, consumer-owned docs, and `.codeheart/user/`. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact classes for managed docs, repository governance, and safety-policy changes. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer procedure for kit source, docs, templates, and validation changes. |
| `docs/repo/runbooks/promote-consumer-change.md` | Existing maintainer path for sanitized consumer-local guidance promotion. |
| `components/agent-interface/component.yaml` | Existing component manifest for managed consumer-facing agent-interface docs. |
| `components/agent-interface/managed/README.md` | Managed agent-interface docs index to route the new consumer feedback runbook. |
| `components/agent-interface/managed/kit-readme.md` | Installed `.codeheart/kit/README.md` fallback route for managed docs. |
| `templates/agents/AGENTS.managed-block.md` | Root `AGENTS.md` managed block route surface for installed consumers. |
| `src/codeheart_operating_kit/resources/manifest.yaml` | Packaged resource manifest requiring updates after component changes. |
| `.github/workflows/validate.yml` | CI command inventory and release asset validation surface. |
| `scripts/validate-public-core.py` | Public-core hygiene validator for issue forms and docs. |
| `scripts/validate-markdown-headers.py` | Timestamp validator for Markdown docs. |
| `scripts/validate-release-manifest.py` | Release manifest validator for packaged metadata. |
| `tests/test_packaging_resources.py` | Packaged resource fallback coverage for installed managed docs. |
| `release-notes.md` | Release-facing consumer-impact and adoption notes. |

## Table Of Contents

- [Section 1 - Foundation](#section-1---foundation)
- [Section 2 - Strategy](#section-2---strategy)
- [Section 3 - Execution Plan](#section-3---execution-plan)
- [Section 4 - Future Planning](#section-4---future-planning)
- [Revision Notes](#revision-notes)

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Ship a v1 feedback intake workflow that lets Operating Kit consumers and agents report sanitized
kit feedback through public GitHub Issues while keeping `.codeheart/kit/` managed, keeping
`.codeheart/user/feedback/` local and ignored, and giving maintainers a repeatable triage path.

Completion is proven when:

- six public GitHub issue forms exist for rough feedback, bugs, doctrine gaps, install/sync/check
  issues, docs routing issues, and feature requests;
- every public issue form requires a public-core sanitization confirmation;
- the approved label taxonomy is documented and the repository label creation route is confirmed;
- installed managed docs route consumers to the feedback workflow without scaffolding local drafts;
- maintainer docs explain triage states, conversion to discovery or implementation plans, and
  accidental public disclosure response;
- authoring files and packaged `src/codeheart_operating_kit/resources/` files are in sync;
- release notes record the instruction-only, repository-governance, and security/safety policy
  impact;
- validation passes locally for Markdown headers, public-core hygiene, package resource coverage,
  release manifests, and the affected pytest suite.

## 1.2 Project And Problem Context

Operating Kit consumers currently have no designed intake path for missing doctrine, confusing
routes, sync issues, stale assumptions, or kit product ideas. Without a designed path, feedback can
land in private chat history, consumer-specific docs, or managed `.codeheart/kit/` files. That
creates drift, privacy risk, and maintainer triage gaps.

The approved discovery chose public GitHub Issues as the canonical shareable backlog. Local
`.codeheart/user/feedback/` may be described as sanitized draft space, but v1 must not scaffold it
or treat it as source of truth. CLI feedback automation, project boards, and private security
reporting are deferred.

## 1.3 Current State Analysis

Existing state:

- `.github/` contains validation workflow configuration but no issue forms.
- `components/agent-interface/` owns managed installed agent-interface docs and the installed
  `.codeheart/kit/README.md` fallback inventory.
- `templates/agents/AGENTS.managed-block.md` owns the managed block inserted into consumer root
  `AGENTS.md` files.
- `src/codeheart_operating_kit/resources/` mirrors component and template resources used when the
  package runs without a source checkout.
- `release-notes.md` exists and must describe consumer-facing impact when the feature ships.
- Validation exists for Markdown headers, public-core hygiene, JSON schemas, release manifests,
  release asset packaging, and packaged resource fallback.

Target state:

- The public repository has issue forms under `.github/ISSUE_TEMPLATE/`.
- Consumers discover feedback guidance from installed managed docs and the managed root `AGENTS.md`
  route surface.
- Maintainers discover triage guidance from `docs/repo/runbooks/`.
- The optional local draft path remains documentation-only and local-user-owned.
- No consumer-owned docs, agent-memory scaffolds, CLI behavior, sync behavior, or local user files
  are overwritten by this v1.

Consumer impact record:

- `instruction-only change` for managed consumer-facing docs and maintainer docs.
- `security or safety policy change` for public-core sanitization and leak-response guidance.
- Repository-governance addition for issue forms and labels.
- No consumer migration required.
- No backwards-compatible scaffold addition in v1.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

```text
Codeheart-Operating-Kit/
  .github/
    ISSUE_TEMPLATE/                                      # create
      config.yml                                        # create
      rough-feedback.yml                                # create
      kit-bug.yml                                       # create
      doctrine-workflow-gap.yml                         # create
      install-sync-check.yml                            # create
      docs-routing.yml                                  # create
      feature-request.yml                               # create
  components/
    agent-interface/
      component.yaml                                    # modify
      managed/
        README.md                                       # modify
        kit-readme.md                                   # modify
        reference/
          kit-feedback-item-format.md                   # create
        runbooks/
          submit-kit-feedback.md                        # create
  templates/
    agents/
      AGENTS.managed-block.md                           # modify
  src/codeheart_operating_kit/resources/
    components/agent-interface/
      component.yaml                                    # modify
      managed/
        README.md                                       # modify
        kit-readme.md                                   # modify
        reference/kit-feedback-item-format.md           # create
        runbooks/submit-kit-feedback.md                 # create
    templates/agents/AGENTS.managed-block.md            # modify
    manifest.yaml                                       # modify
  src/codeheart_operating_kit/
    components.py                                       # modify
    commands/
      sync.py                                           # modify
  docs/
    README.md                                           # modify
    repo/
      README.md                                         # modify
      reference/
        kit-feedback-label-taxonomy.md                  # create
      runbooks/
        triage-kit-feedback.md                          # create
      plans/
        README.md                                       # modify
        kit-feedback-intake/
          kit-feedback-intake_discovery_doc.md          # reference only
          kit-feedback-intake_implementation_doc.md     # this plan
  release-notes.md                                      # modify
  tests/
    test_packaging_resources.py                         # modify
    test_init.py                                        # modify
    test_sync_check.py                                  # modify
```

The implementation should treat `components/` and `templates/` as the authoring source and
`src/codeheart_operating_kit/resources/` as the packaged mirror. Every managed consumer-facing
file added to the authoring tree must be mirrored into packaged resources and listed in the
component manifest.

## 2.2 Open Questions And Assumptions Requiring Clarification

`OQ-1`: GitHub Issues, issue forms, and label governance preflight.

- `BLOCKER: yes`
- Affects: `E1` through `E6`; resolved by `E1`
- Unlocks: issue-form creation, issue-form defaults, and repository label creation.
- Recommended default: proceed only after confirming GitHub Issues are enabled, issue forms may be
  added through repository files, and the approved label set may be created through maintainer
  repository governance.

`OQ-2`: Private security reporting channel.

- `BLOCKER: no`
- Affects: `E2`, `E4`
- Unlocks: future private disclosure design.
- Recommended default: keep private security reporting out of v1, omit public security issue
  forms, and route accidental disclosure handling through maintainer triage.

`OQ-3`: CLI feedback command timing.

- `BLOCKER: no`
- Affects: `E3`, `E5`
- Unlocks: future automation design.
- Recommended default: keep CLI feedback drafting out of v1 and avoid adding CLI tests for this
  feature.

## 2.3 Architectural Decisions With Reasoning

`AD-1`: Use the existing `agent-interface` component for managed consumer feedback guidance.

1. Problem being solved: consumers need a managed route without introducing a new component.
2. Simplest working solution: add the feedback runbook and item-format reference to
   `components/agent-interface/managed/`.
3. What may change in 6-12 months: feedback intake could become a dedicated component after CLI
   support, schemas, or generated draft scaffolds exist.
4. Rationale: this keeps v1 as an instruction-only managed-doc change instead of a component
   addition.
5. Alternatives considered: a new `feedback-intake` component was rejected for v1 because it adds
   profile, manifest, sync, and release complexity before the intake model is proven.

`AD-2`: Use GitHub issue forms, not Markdown issue templates.

1. Problem being solved: public submissions need required sanitization confirmations and structured
   triage fields.
2. Simplest working solution: create YAML issue forms with required fields and checkboxes.
3. What may change in 6-12 months: a CLI command may generate a matching issue body from local
   metadata.
4. Rationale: issue forms give enough structure without creating a separate backlog system.
5. Alternatives considered: one generic Markdown template was rejected because it cannot reliably
   enforce required confirmations.

`AD-3`: Document labels in the repo and create labels through approved repository governance.

1. Problem being solved: triage states must be visible, but GitHub labels are external repository
   state.
2. Simplest working solution: add a label taxonomy reference and confirm the label creation route
   during preflight.
3. What may change in 6-12 months: a label-sync workflow can be added when label drift becomes a
   repeated maintenance problem.
4. Rationale: v1 needs clear labels, not label automation.
5. Alternatives considered: a GitHub project board was rejected because triage volume has not
   justified it.

`AD-4`: Keep `.codeheart/user/feedback/` documentation-only in v1.

1. Problem being solved: local drafts are useful, but scaffolding a local draft folder could imply
   it is synchronized or safe for raw evidence.
2. Simplest working solution: document sanitized local draft expectations without adding scaffold
   files or sync behavior.
3. What may change in 6-12 months: a future CLI command may create sanitized drafts after the
   schema stabilizes.
4. Rationale: ignored local files are not a security boundary and must not become a shadow backlog.
5. Alternatives considered: automatic scaffold creation was rejected because v1 has no draft
   lifecycle automation.

`AD-5`: Put accidental public disclosure response in maintainer triage.

1. Problem being solved: public issue intake can receive sensitive material despite warnings.
2. Simplest working solution: maintainer runbook includes no-copy response steps, public exposure
   minimization, secret rotation guidance, sanitized summary preservation, and private escalation
   to the kit owner.
3. What may change in 6-12 months: a dedicated private security disclosure process may replace the
   temporary triage path.
4. Rationale: prevention alone is not enough for a public intake surface.
5. Alternatives considered: a public security issue form was rejected because it would invite
   sensitive reports into a public repository.

`AD-6`: Treat the release as instruction-only plus security/safety policy and repository
governance.

1. Problem being solved: the change affects managed instructions, public reporting policy, and
   repository issue configuration.
2. Simplest working solution: record the combined impact in release notes and validation evidence.
3. What may change in 6-12 months: CLI or scaffold additions would need new impact records.
4. Rationale: no existing consumer migration is required, but the safety-policy change is material.
5. Alternatives considered: treating this as docs-only was rejected because sanitization and leak
   response change safety guidance.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| `E1` | Repository governance preflight confirms the GitHub issue surface can be implemented. | S | None |
| `E2` | Public issue forms and label taxonomy exist with required sanitization controls. | M | `E1` |
| `E3` | Managed consumer docs route users to safe feedback capture and submission. | M | `E1` |
| `E4` | Maintainer triage docs define lifecycle states, labels, plan conversion, and leak response. | M | `E2` |
| `E5` | Packaged resources and component manifests mirror the authored managed-doc changes. | M | `E3`, `E4` |
| `E6` | Release notes, indexes, and validation evidence are complete. | M | `E2`, `E3`, `E4`, `E5` |

## E1 - Repository Governance Preflight

### A) Epic ID, Title, And Outcome

`E1` - Repository Governance Preflight.

Outcome: GitHub Issues, issue forms, and label creation are confirmed before repository-governance
files or external repository state are changed.

### B) Scope

Confirm the repository can support the approved issue intake model. Record the result in the
execution log for later implementation PR summary use. Stop normal implementation when preflight
fails.

### C) Files Touched

- No required file changes.

```text
(no file changes)
```

### D) Acceptance Criteria And Size

- Size: `S`
- Issues are confirmed enabled for `Codeheart-Operating-Kit`.
- Issue forms under `.github/ISSUE_TEMPLATE/` are confirmed acceptable for this repository.
- The label creation route is confirmed.
- Execution stops before `E2` when any preflight item fails.

### E) Dependencies And Critical-Path Notes

This epic has no dependencies. `E2` through `E6` must not start until this epic passes.

### F) Tasks Checklist

- [x] Confirm Issues are enabled for `Codeheart-Operating-Kit`.
- [x] Confirm issue form files under `.github/ISSUE_TEMPLATE/` are accepted through repository PR review.
- [x] Confirm label creation authority for the approved label set.
- [x] Record the preflight result in the execution log for later implementation PR summary use.
- [x] Stop execution before `E2` when any preflight confirmation fails.
- [x] Run `git status --short --branch` and verify no unrelated file changes were created.

### G) Implementation Notes

Use GitHub repository metadata, maintainer confirmation, or an approved GitHub command path. Do not
change public GitHub settings during preflight unless the user explicitly approves that governance
action.

### H) Open Questions

- `OQ-1` blocks this epic until resolved.

## E2 - Public Issue Forms And Label Taxonomy

### A) Epic ID, Title, And Outcome

`E2` - Public Issue Forms And Label Taxonomy.

Outcome: GitHub issue forms collect structured, sanitized feedback and a public label taxonomy
defines the lifecycle labels maintainers will use.

### B) Scope

Create six issue forms and one label taxonomy reference. Each issue form must include required
public-core warnings and sanitization confirmation. Apply the approved GitHub label set through
the confirmed governance path. No public security report form is added.

### C) Files Touched

- `.github/ISSUE_TEMPLATE/config.yml`
- `.github/ISSUE_TEMPLATE/rough-feedback.yml`
- `.github/ISSUE_TEMPLATE/kit-bug.yml`
- `.github/ISSUE_TEMPLATE/doctrine-workflow-gap.yml`
- `.github/ISSUE_TEMPLATE/install-sync-check.yml`
- `.github/ISSUE_TEMPLATE/docs-routing.yml`
- `.github/ISSUE_TEMPLATE/feature-request.yml`
- `docs/repo/reference/kit-feedback-label-taxonomy.md`
- GitHub repository labels

```text
.github/
  ISSUE_TEMPLATE/
    config.yml
    rough-feedback.yml
    kit-bug.yml
    doctrine-workflow-gap.yml
    install-sync-check.yml
    docs-routing.yml
    feature-request.yml
docs/
  repo/
    reference/
      kit-feedback-label-taxonomy.md
GitHub repository labels                              # external state
```

### D) Acceptance Criteria And Size

- Size: `M`
- Blank issues are disabled or routed to the approved feedback forms.
- Six issue forms exist and use public-safe examples only.
- Every issue form requires confirmation that the issue excludes secrets, credentials, customer or
  tenant details, local machine state, account identifiers, raw logs, and private strategy.
- Rough feedback is routed to `needs-shaping`.
- The label taxonomy defines feedback type labels and lifecycle labels.
- The approved label set exists in GitHub after repository-governance application.
- Issue form YAML parses and local structure validation passes.
- GitHub issue-form preview is recorded as PR-review evidence when the branch is pushed.

### E) Dependencies And Critical-Path Notes

Depends on `E1`. Issue form labels may reference labels before they are created, but the label
creation path must be confirmed by `E1`.

### F) Tasks Checklist

- [x] Create `.github/ISSUE_TEMPLATE/config.yml` for the approved public intake surface.
- [x] Create `.github/ISSUE_TEMPLATE/rough-feedback.yml` with weak-signal fields and `needs-shaping` default labeling.
- [x] Create `.github/ISSUE_TEMPLATE/kit-bug.yml` with kit version, observed behavior, expected behavior, reproduction, and evidence fields.
- [x] Create `.github/ISSUE_TEMPLATE/doctrine-workflow-gap.yml` with affected doctrine, missing guidance, context, and expected guidance fields.
- [x] Create `.github/ISSUE_TEMPLATE/install-sync-check.yml` with command, version, target surface, observed result, and sanitized output summary fields.
- [x] Create `.github/ISSUE_TEMPLATE/docs-routing.yml` with current route, expected route, audience, and confusion summary fields.
- [x] Create `.github/ISSUE_TEMPLATE/feature-request.yml` with user goal, target workflow, value, constraints, and acceptance signal fields.
- [x] Add required public-core sanitization confirmation checkboxes to every issue form.
- [x] Add approved default feedback type labels and lifecycle labels to every issue form.
- [x] Create `docs/repo/reference/kit-feedback-label-taxonomy.md` with label names, meanings, and transition rules.
- [x] Apply the approved GitHub label set through the confirmed governance path.
- [x] Verify the approved GitHub label set exists in GitHub.
- [x] Run a local YAML parse check for every `.github/ISSUE_TEMPLATE/*.yml` file.
- [x] Run local issue-form structure validation for required GitHub issue-form fields.
- [x] Record GitHub issue-form preview as a PR-review validation item for branch review.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Inspect every issue form for required sanitization confirmation.
- [x] Inspect every issue form for public-safe example text.

### G) Implementation Notes

Use GitHub issue form YAML syntax. Keep examples generic and avoid real account IDs, tenant names,
customer names, raw paths, and raw logs. The rough-feedback form should explicitly accept
incomplete friction and ideas without implying that the issue is implementation-ready. Label
application changes external GitHub repository state, so use the governance path confirmed in
`E1` and record the label verification evidence in the PR summary. GitHub renders issue-form
previews from repository branches, not unpushed local files; use local YAML and structure
validation during implementation, then capture GitHub-rendered preview evidence during PR review.

### H) Open Questions

- None after `E1` passes.

## E3 - Managed Consumer Feedback Guidance

### A) Epic ID, Title, And Outcome

`E3` - Managed Consumer Feedback Guidance.

Outcome: installed consumers can find safe guidance for reporting Operating Kit feedback without
editing managed kit files or committing private local drafts.

### B) Scope

Add managed agent-interface docs for submitting kit feedback and formatting a feedback item. Update
managed route surfaces. Do not scaffold `.codeheart/user/feedback/` and do not add CLI behavior.

### C) Files Touched

- `components/agent-interface/managed/runbooks/submit-kit-feedback.md`
- `components/agent-interface/managed/reference/kit-feedback-item-format.md`
- `components/agent-interface/managed/README.md`
- `components/agent-interface/managed/kit-readme.md`
- `templates/agents/AGENTS.managed-block.md`
- `components/agent-interface/component.yaml`
- `src/codeheart_operating_kit/components.py`
- `src/codeheart_operating_kit/commands/sync.py`
- `tests/test_init.py`
- `tests/test_sync_check.py`

```text
components/
  agent-interface/
    component.yaml
    managed/
      README.md
      kit-readme.md
      reference/
        kit-feedback-item-format.md
      runbooks/
        submit-kit-feedback.md
templates/
  agents/
    AGENTS.managed-block.md
src/
  codeheart_operating_kit/
    components.py
    commands/
      sync.py
tests/
  test_init.py
  test_sync_check.py
```

### D) Acceptance Criteria And Size

- Size: `M`
- Managed docs tell consumers not to edit `.codeheart/kit/` for feedback.
- Managed docs define `.codeheart/user/feedback/` as optional ignored sanitized draft space only.
- Managed docs point sanitized shareable feedback to public GitHub Issues.
- The item format captures kit version, component area, problem, expected behavior, sanitized
  evidence summary, privacy confirmation, and proposed classification.
- Component manifest entries target `.codeheart/kit/docs/agent-interface/...`.

### E) Dependencies And Critical-Path Notes

Depends on `E1`. The docs should reference the approved issue intake only after preflight passes.

### F) Tasks Checklist

- [x] Create `components/agent-interface/managed/runbooks/submit-kit-feedback.md` with local draft, sanitization, public submission, and non-editing guidance.
- [x] Create `components/agent-interface/managed/reference/kit-feedback-item-format.md` with field definitions and public-safe examples.
- [x] Update `components/agent-interface/managed/README.md` with routes to the new feedback runbook and item-format reference.
- [x] Update `components/agent-interface/managed/kit-readme.md` with a fallback route to the feedback runbook.
- [x] Update `templates/agents/AGENTS.managed-block.md` with a managed route for Operating Kit feedback.
- [x] Update `components/agent-interface/component.yaml` with managed file entries for the new docs.
- [x] Add `.codeheart/user/feedback/` to local-user gitignore behavior without scaffolding the directory.
- [x] Add `.codeheart/user/feedback/` to sync-time gitignore repair for existing consumers.
- [x] Add tests for local-user feedback draft gitignore behavior.
- [x] Confirm no task creates `.codeheart/user/feedback/` as a scaffolded path.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.

### G) Implementation Notes

Use neutral public-safe placeholders. The runbook should state that private evidence stays outside
the feedback workflow and public issues receive sanitized summaries only.

### H) Open Questions

- None after `E1` passes.

## E4 - Maintainer Triage Workflow

### A) Epic ID, Title, And Outcome

`E4` - Maintainer Triage Workflow.

Outcome: kit maintainers have an ordered workflow for processing feedback issues into closure,
accepted backlog, discovery, implementation planning, release, and synced consumer guidance.

### B) Scope

Create a maintainer runbook and update docs indexes. The runbook must include lifecycle states,
label use, duplicate handling, consumer-specific closure, discovery and implementation handoff,
release-note expectations, and accidental disclosure response.

### C) Files Touched

- `docs/repo/runbooks/triage-kit-feedback.md`
- `docs/repo/README.md`
- `docs/README.md`

```text
docs/
  README.md
  repo/
    README.md
    runbooks/
      triage-kit-feedback.md
```

### D) Acceptance Criteria And Size

- Size: `M`
- Triage runbook covers duplicate, declined, needs information, consumer-specific, accepted
  backlog, needs discovery, implementation planned, released, and superseded outcomes.
- Triage runbook includes `needs-shaping` handling for rough feedback.
- Triage runbook tells maintainers when to create discovery and implementation docs.
- Triage runbook includes no-copy accidental public disclosure response.
- Docs indexes expose the new runbook and label taxonomy.

### E) Dependencies And Critical-Path Notes

Depends on `E2` because runbook labels and issue forms must align.

### F) Tasks Checklist

- [x] Create `docs/repo/runbooks/triage-kit-feedback.md` with ordered intake review steps.
- [x] Add lifecycle state definitions to `docs/repo/runbooks/triage-kit-feedback.md`.
- [x] Add label application guidance to `docs/repo/runbooks/triage-kit-feedback.md`.
- [x] Add conversion rules for discovery docs and implementation docs to `docs/repo/runbooks/triage-kit-feedback.md`.
- [x] Add consumer-specific closure guidance to `docs/repo/runbooks/triage-kit-feedback.md`.
- [x] Add accidental public disclosure response steps to `docs/repo/runbooks/triage-kit-feedback.md`.
- [x] Update `docs/repo/README.md` with routes to the triage runbook and label taxonomy.
- [x] Update `docs/README.md` with routes to the triage runbook and label taxonomy.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.

### G) Implementation Notes

The runbook should explicitly avoid quoting sensitive material from public issues. Preserve only a
sanitized summary when triage must continue after exposure handling.

### H) Open Questions

- None.

## E5 - Packaged Resource And Manifest Sync

### A) Epic ID, Title, And Outcome

`E5` - Packaged Resource And Manifest Sync.

Outcome: source component changes are mirrored into packaged resources so installed CLI behavior
matches the authoring tree in source and packaged fallback modes.

### B) Scope

Mirror changed agent-interface managed docs, component manifest entries, and the managed block
template into `src/codeheart_operating_kit/resources/`. Refresh resource manifest metadata and
package fallback tests.

### C) Files Touched

- `src/codeheart_operating_kit/resources/components/agent-interface/component.yaml`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/README.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/kit-readme.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/kit-feedback-item-format.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/runbooks/submit-kit-feedback.md`
- `src/codeheart_operating_kit/resources/templates/agents/AGENTS.managed-block.md`
- `src/codeheart_operating_kit/resources/manifest.yaml`
- `tests/test_packaging_resources.py`

```text
src/
  codeheart_operating_kit/
    resources/
      manifest.yaml
      components/
        agent-interface/
          component.yaml
          managed/
            README.md
            kit-readme.md
            reference/
              kit-feedback-item-format.md
            runbooks/
              submit-kit-feedback.md
      templates/
        agents/
          AGENTS.managed-block.md
tests/
  test_packaging_resources.py
```

### D) Acceptance Criteria And Size

- Size: `M`
- Packaged resource files match their source counterparts.
- `agent-interface` packaged manifest lists the new managed docs.
- `src/codeheart_operating_kit/resources/manifest.yaml` records updated component metadata and
  consumer impact.
- Packaged fallback tests assert the new managed docs are installed.

### E) Dependencies And Critical-Path Notes

Depends on `E3` and `E4`. Keep source and packaged resource edits in the same commit.

### F) Tasks Checklist

- [x] Mirror `components/agent-interface/managed/runbooks/submit-kit-feedback.md` into packaged resources.
- [x] Mirror `components/agent-interface/managed/reference/kit-feedback-item-format.md` into packaged resources.
- [x] Mirror changed agent-interface managed indexes into packaged resources.
- [x] Mirror changed `components/agent-interface/component.yaml` into packaged resources.
- [x] Mirror changed `templates/agents/AGENTS.managed-block.md` into packaged resources.
- [x] Refresh `src/codeheart_operating_kit/resources/manifest.yaml` component metadata and consumer impact.
- [x] Update `tests/test_packaging_resources.py` to assert packaged fallback installs the new feedback docs.
- [x] Run `python3 -m pytest tests/test_packaging_resources.py -q`.
- [x] Run `python3 scripts/validate-release-manifest.py`.

### G) Implementation Notes

Use Python-based checksum calculation or an existing repo helper for manifest checksum fields so
the plan does not depend on platform-specific shell utilities. The release manifest validator
checks shape; the implementer should still review changed checksums for correctness.

### H) Open Questions

- None.

## E6 - Release Notes, Indexes, And Validation

### A) Epic ID, Title, And Outcome

`E6` - Release Notes, Indexes, And Validation.

Outcome: the implementation is discoverable, release-ready, and backed by validation evidence.

### B) Scope

Update release notes, plan indexes, final docs routes, and validation evidence. Do not tag,
publish, merge, or change repository settings in this epic.

### C) Files Touched

- `release-notes.md`
- `docs/README.md`
- `docs/repo/README.md`
- `docs/repo/plans/README.md`
- `docs/repo/plans/kit-feedback-intake/kit-feedback-intake_implementation_doc.md`

```text
release-notes.md
docs/
  README.md
  repo/
    README.md
    plans/
      README.md
      kit-feedback-intake/
        kit-feedback-intake_implementation_doc.md
```

### D) Acceptance Criteria And Size

- Size: `M`
- Release notes classify the change as instruction-only, repository-governance, and
  security/safety policy.
- Release notes state that no consumer migration is required.
- Plan and docs indexes route to the discovery and implementation documents.
- Full local validation passes for the changed surfaces.

### E) Dependencies And Critical-Path Notes

Depends on `E2`, `E3`, `E4`, and `E5`. This is the final verification epic.

### F) Tasks Checklist

- [x] Update `release-notes.md` with feedback intake summary, consumer impact, validation, and no-migration note.
- [x] Update `docs/repo/plans/README.md` with the implementation plan route.
- [x] Update `docs/repo/README.md` with final feedback intake routes.
- [x] Update `docs/README.md` with final feedback intake routes.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 scripts/validate-json-schemas.py`.
- [x] Run `python3 scripts/validate-release-manifest.py`.
- [x] Run `python3 -m pytest -q`.
- [x] Run `git diff --check`.
- [x] Review `git status --short --branch` and confirm only planned files changed.

### G) Implementation Notes

Do not run release publishing, tag creation, PR merge, or GitHub settings changes as part of this
implementation unless the user explicitly requests those external-state actions after the code and
docs changes are reviewed.

### H) Open Questions

- None.

## 3.1 Release Handoff

Release, publishing, GitHub release creation, and consumer update/sync adoption are intentionally
outside the implementation epics in this plan. After implementation validation is complete, ship
the change through `docs/repo/runbooks/release-operating-kit.md`, then have consumer repositories
update to the newly published Operating Kit version and run normal sync/check adoption.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- CLI-assisted feedback drafting is deferred until real issue examples prove the intake fields and
  local metadata needs.
- `.codeheart/user/feedback/` scaffold creation is deferred because local drafts are not a
  security boundary and should not imply a synchronized backlog.
- A GitHub project board is deferred until issue volume proves labels and maintainer comments are
  insufficient.
- A private security reporting channel is deferred to a focused security intake plan.
- Label automation is deferred until manual label governance creates repeated maintenance cost.
- Issue-form schema validation is deferred until issue form drift becomes a recurring risk.

## 4.2 Future Considerations

- A future CLI command can prefill kit version, component, installed paths, and lockfile metadata
  into a sanitized draft.
- A future private disclosure process should define security contact, response SLA, GitHub
  security advisory use, and release-note boundaries.
- A future feedback-intake component may become warranted if the workflow grows beyond
  agent-interface guidance.
- A future validator can assert that every public issue form contains the required sanitization
  confirmation.

# Revision Notes

- 2026-06-17: Initial draft implementation plan created from the approved kit feedback intake
  discovery.
- 2026-06-17: Tightened execution readiness by widening `OQ-1` impact, adding explicit GitHub
  label application and verification, and adding issue-form YAML and preview validation tasks.
- 2026-06-17: Activated the plan for goal-style execution and linked the sibling execution log.
- 2026-06-17: Adjusted E1 evidence recording to use the execution log before a PR summary exists.
- 2026-06-17: Recorded the E2 validation substitution: local issue-form YAML and structure
  validation during implementation, with GitHub-rendered preview evidence deferred to PR review.
- 2026-06-17: Added E3 code and test tasks to make `.codeheart/user/feedback/` genuinely ignored
  local-user draft space without scaffolding it.
- 2026-06-17: Completed all epics after final validation and reviewer acceptance.
- 2026-06-17: Corrected the packaged resource visual hierarchy for the feedback submission
  runbook path.
- 2026-06-17: Added the release handoff note, keeping tag, publish, and consumer adoption under
  the separate Operating Kit release runbook.
