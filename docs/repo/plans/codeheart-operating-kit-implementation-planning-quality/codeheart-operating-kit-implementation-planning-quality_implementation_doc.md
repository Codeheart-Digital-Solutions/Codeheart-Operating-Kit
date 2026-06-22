Last updated: 2026-06-22T18:50:19Z (UTC)
Created: 2026-06-22
Status: completed
Completed: 2026-06-22
Execution log: docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_execution_log.md

# Document Header

This draft implementation plan turns the completed implementation-planning quality discovery into
a source-repository change plan for Codeheart Operating Kit. The plan targets the reusable managed
planning workflow runbooks in this `Codeheart-Operating-Kit` source repository. It does not
hand-edit the installed `.codeheart/kit/` snapshot in this consumer repository; consumer adoption
happens through the normal Operating Kit sync command after release.

Essential context:

| Source | Why it matters |
| --- | --- |
| `<first-consumer-repository>/docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_discovery_doc.md` | Accepted discovery, decisions, capability-scope handoff, and feature-level success evidence. |
| `AGENTS.md` | Source-repository bootstrap rules for public-core safety, task routing, and release boundaries. |
| `README.md` | Source-repository purpose, public boundary, and maintainer entry points. |
| `docs/repo/runbooks/change-operating-kit.md` | Required source-repository procedure before changing kit source, docs, resources, tests, or release assets. |
| `docs/repo/runbooks/release-operating-kit.md` | Required source-repository procedure for release assets, checksums, tag, GitHub release, and release evidence. |
| `docs/repo/reference/placement-contract.md` | Managed content placement and source/resource ownership rules. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact classification for managed instruction changes and release-note expectations. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | Source managed discovery workflow to receive the capability-scope handoff requirement. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Source managed implementation-planning workflow to receive feature-capability and concreteness drafting rules. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | Source managed planning review workflow to receive capability coverage and hollow-plan checks. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Source managed execution workflow to receive per-epic delivered-capability review checks. |
| `tests/test_packaging_resources.py` | Packaged-resource parity test covering changed managed planning workflow files. |
| `release-notes.md` | Consumer-facing release-note surface for instruction-only managed workflow changes. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Update Codeheart Operating Kit source managed planning workflow guidance so downstream consumer
repositories receive implementation-planning, review, and execution runbooks that:

- preserve intended feature capability from discovery into implementation plans;
- require concrete implementation planning rather than policy, doctrine, gates, stubs,
  scaffolding, or validation shells alone;
- keep atomic checklist tasks at capability-sized granularity;
- reject avoidable non-concreteness when facts can be checked during planning;
- require exact preflight, remediation, retry validation, and stop conditions when legitimate
  execution-time variability remains;
- make planning-document reviewers catch hollow, under-covered, or quietly narrowed plans;
- make per-epic execution reviewers check delivered feature capability, not only checklist
  completion.

Completion is proven when:

- the source runbooks under `components/planning-workflows/managed/runbooks/`
  contain the new discovery handoff, implementation drafting, planning review, and execution
  review guidance;
- the mirrored packaged resources under
  `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/`
  match the source managed runbooks;
- source release notes record this as an instruction-only managed workflow change with no required
  consumer migration beyond normal update and sync;
- source validation passes for Markdown headers, public-core hygiene, packaged-resource parity,
  focused tests, and whitespace;
- the source release is published through `docs/repo/runbooks/release-operating-kit.md`;
- this consumer repository runs normal update, sync, and check against the released kit with zero
  managed-content drift.

## 1.2 Project And Problem Context

The Operating Kit is consumed by other repositories through managed snapshots under
`.codeheart/kit/`. The reusable source-of-truth change belongs in the public
`Codeheart-Operating-Kit` source repository, not in this consumer repository's installed snapshot.

The completed discovery found that the main problem is not actor labeling. The problem is an
implementation-planning mindset: plans can satisfy current document shape, gates, and checklist
rules while leaving future implementers to invent the feature workflow, command sequence, data
shape, permission model, validation method, or meaningful behavior. Plans can also be concrete
about scaffolding, policy, schemas, evidence, and validators while still under-covering the
intended feature capability.

The implementation must keep the atomic checkbox style the user prefers. The fix should raise the
quality bar for what each checkbox represents, without forcing mandatory `Agent:` / `User:`
prefixes, approval-only task taxonomy, sentence-level tasks, or handoff-sized granularity.

## 1.3 Current State Analysis

Current source state:

- `discovery-workflow.md` has implementation-handoff readiness criteria, but it does not require
  a structured feature-capability scope block before normal implementation handoff.
- `draft-implementation-plan.md` requires concrete, non-branching checklist tasks and linear
  epics, but it does not strongly require capability coverage against accepted discovery or
  explicitly reject plans that are concrete only around support structures.
- `review-planning-document.md` checks execution readiness, concrete checklist tasks, and
  validation coverage, but it does not ask whether the plan concretely implements the intended
  feature capability instead of only surrounding policy, scaffolding, gates, schemas, stubs, or
  validation shells.
- `execute-implementation-plan.md` already says epic outcomes outrank literal checklist items and
  has a per-epic review gate, but the reviewer checklist does not explicitly verify delivered
  feature capability against discovery capability scope when one exists.
- The source repository has mirrored packaged resources under `src/codeheart_operating_kit/resources/`
  and `tests/test_packaging_resources.py` asserts parity for the planning workflow files.
- The source repository has validators for Markdown headers, public-core hygiene, JSON schemas,
  and release manifests. This change is instruction-only and should not require schema or CLI
  behavior changes.

Target state:

- Discovery handoffs include compact structured capability-scope blocks for implementation-
  relevant decision groups included in normal implementation handoff.
- Implementation planning starts from the intended feature capability and makes omitted areas
  explicit as out of scope, deferred, or blocked.
- Checklist rules teach capability-sized concrete tasks and reject avoidable vague tasks.
- Planning review has a clear hollow-plan test: can a lazy implementer complete every checklist
  item while leaving the intended feature capability narrow, stubbed, policy-only, or unusable?
- Execution review checks that each completed epic delivered the intended capability described by
  the epic outcome and, when available, the discovery capability scope.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement the change in the source Operating Kit repository, then mirror the changed managed
runbooks into package resources and validate the public managed surface.

Expected source tree:

```text
./
  components/
    planning-workflows/
      managed/
        runbooks/
          discovery-workflow.md                 # modify
          draft-implementation-plan.md          # modify
          review-planning-document.md           # modify
          execute-implementation-plan.md        # modify
  src/
    codeheart_operating_kit/
      resources/
        components/
          planning-workflows/
            managed/
              runbooks/
                discovery-workflow.md           # modify, mirror source
                draft-implementation-plan.md    # modify, mirror source
                review-planning-document.md     # modify, mirror source
                execute-implementation-plan.md  # modify, mirror source
  tests/
    test_packaging_resources.py                 # inspect, existing parity coverage
    test_install_metadata.py                    # release validation
    test_release_assets.py                      # release validation
  release-notes.md                              # modify
  pyproject.toml                                # modify during release
  src/codeheart_operating_kit/__init__.py        # modify during release
  manifest.yaml                                 # modify during release
  src/codeheart_operating_kit/resources/manifest.yaml  # modify during release
  bootstrap.md                                  # modify during release
  install.sh                                    # modify during release
  install.ps1                                   # modify during release
  dist/
    codeheart-operating-kit-0.1.7-macos.tar.gz         # create during release
    codeheart-operating-kit-0.1.7-macos.tar.gz.sha256  # create during release
    codeheart-operating-kit-0.1.7-windows.zip          # create during release
    codeheart-operating-kit-0.1.7-windows.zip.sha256   # create during release
  docs/
    repo/
      plans/
        codeheart-operating-kit-implementation-planning-quality/
          codeheart-operating-kit-implementation-planning-quality_implementation_doc.md  # create in source execution
          codeheart-operating-kit-implementation-planning-quality_execution_log.md        # create in source execution
```

This source copy is the canonical source-repository implementation plan for execution. Source
execution uses source-repo relative paths.

## 2.2 Open Questions And Assumptions Requiring Clarification

`OQ-1` - Should the source execution create a source-repository copy of this implementation plan?

- `BLOCKER: no`
- `Affects: EP0`
- Unlocks whether the source repo has its own canonical implementation plan before editing source
  runbooks.
- Recommended default: create the source plan copy under
  `docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/` before implementation
  execution so source register, execution log, and review gates live beside the changed kit files.

`OQ-2` - What target release version should EP6 use?

- `BLOCKER: no`
- `Affects: EP5`, `EP6`
- Unlocks exact release surface updates, asset names, tag name, and consumer sync target version.
- Recommended default: use `v0.1.7` because the source repository currently exposes `v0.1.6` in
  `pyproject.toml`, `src/codeheart_operating_kit/__init__.py`, release manifests, installers, and
  release notes. If source HEAD already contains `v0.1.7` before EP6 starts, stop and amend the
  source plan to the next patch version.

`OQ-3` - Should source execution add new tests beyond the existing packaged-resource parity test?

- `BLOCKER: no`
- `Affects: EP5`
- Unlocks whether validation stays docs-only or adds text-contract assertions.
- Recommended default: do not add new tests unless source review finds a stable text contract worth
  asserting. The smallest sufficient validation is Markdown headers, public-core hygiene,
  packaged-resource parity, focused tests, and whitespace.

## 2.3 Architectural Decisions With Reasoning

AD-1 - Implement in Operating Kit source managed runbooks, not consumer snapshot

1. Problem being solved: consumer repositories need reusable guidance changes, and `.codeheart/kit/`
   is managed installed content.
2. Simplest working solution: update source `components/planning-workflows/managed/runbooks/` files
   and their packaged resource mirrors.
3. What may change in 6-12 months: the kit may gain a docs-generation or resource-sync command.
4. Rationale: source-managed content is the only durable place for reusable doctrine consumed by
   other repos.
5. Alternatives considered: hand-editing the consumer snapshot was rejected because it would create
   drift and would not help other repositories.

AD-2 - Patch existing workflow gates instead of adding a new taxonomy

1. Problem being solved: the current failure mode is weak implementation mindset, not missing actor
   labels.
2. Simplest working solution: add capability coverage, capability-sized task, fresh implementer,
   and avoidable non-concreteness checks inside the current discovery, planning, review, and
   execution runbooks.
3. What may change in 6-12 months: the kit may add planning-lint checks or reviewer heuristics
   after repeated plan reviews show stable failure patterns.
4. Rationale: existing runbooks already define document shape and review flow; targeted amendments
   reduce ambiguity without forcing task granularity that the user rejected.
5. Alternatives considered: mandatory actor prefixes, approval-only taxonomy, and deliverable-type
   labels were rejected because they can be gamed while the feature capability remains unplanned.

AD-3 - Require discovery capability-scope blocks only for normal implementation handoff candidates

1. Problem being solved: discovery prose can be too easy to narrow accidentally during planning.
2. Simplest working solution: add a compact structured block to normal implementation handoffs for
   implementation-relevant decision groups.
3. What may change in 6-12 months: broader capability-scope usage may be considered after the
   workflow is piloted.
4. Rationale: this gives implementation planning and review a concrete coverage target without
   turning discovery into an implementation plan.
5. Alternatives considered: requiring the block for every discovery decision group was rejected as
   too heavy for early discovery.

AD-4 - Treat this as instruction-only consumer impact

1. Problem being solved: source release planning needs impact classification.
2. Simplest working solution: classify as `instruction-only change`, update release notes, publish
   the patch release through the release runbook, and require normal consumer update and sync for
   adoption.
3. What may change in 6-12 months: future validators could turn some guidance into machine-checked
   planning lint.
4. Rationale: the planned change updates managed docs and packaged resources without changing CLI
   behavior, generated path placement, schemas, or required migration behavior. Release and sync are
   needed so the managed guidance actually reaches consumers.
5. Alternatives considered: validator-only change was rejected because no validator behavior is
   planned in this pass.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP0 | Source execution context is prepared and the source repo has a canonical plan copy. | S | none |
| EP1 | Discovery workflow requires capability-scope blocks before normal implementation handoff. | M | EP0 |
| EP2 | Draft implementation plan workflow requires concrete feature-capability coverage. | M | EP1 |
| EP3 | Planning-document review workflow catches hollow, under-covered, and avoidably vague plans. | M | EP2 |
| EP4 | Execution workflow per-epic review checks delivered feature capability. | S | EP3 |
| EP5 | Packaged resources, release notes, validation, and consumer-impact handoff are complete. | M | EP1-EP4 |
| EP6 | The Operating Kit release is published and first consumer repository sync/check proves consumer adoption. | L | EP5 |

## EP0 - Source Context And Canonical Plan Setup

### A) Epic ID, Title, And Outcome

EP0 - Source Context And Canonical Plan Setup

Outcome: the source repository has a canonical implementation plan copy, local routing context is
fresh, and execution can proceed against source paths with no consumer-snapshot edits.

### B) Scope

In scope:

- Re-read source-repository routing and change runbooks immediately before source execution.
- Create or refresh the source-repository implementation plan copy.
- Add a source plan-register entry.
- Confirm source worktree state before edits.

Out of scope:

- Editing managed runbooks.
- Editing release notes.
- Running public release steps.

### C) Files Touched

```text
./
  docs/
    repo/
      plans/
        codeheart-operating-kit-implementation-planning-quality/
          codeheart-operating-kit-implementation-planning-quality_implementation_doc.md  # create
          codeheart-operating-kit-implementation-planning-quality_execution_log.md        # create
      plans/
        plan-register.md                                                                # modify
```

### D) Acceptance Criteria And Size

Size: S.

Acceptance criteria:

- Source execution context is documented in the source plan copy.
- Source execution log exists beside the source plan copy.
- Source plan-register entry points to the source plan copy.
- Source `git status --short` is reviewed before runbook edits begin.
- No `.codeheart/kit/` consumer snapshot files are edited.

### E) Dependencies And Critical-Path Notes

No upstream epic dependency. This epic should run first so later edits and review evidence live in
the source repository.

### F) Tasks Checklist

- [x] Read `AGENTS.md`, `README.md`, `docs/repo/runbooks/change-operating-kit.md`, `docs/repo/reference/placement-contract.md`, and `docs/repo/reference/consumer-impact-classification.md` in this source repository.
- [x] Create `docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_implementation_doc.md` in this source repository as an adapted source-repository plan with source-relative paths.
- [x] Create `docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_execution_log.md` in this source repository.
- [x] Add the source plan entry to `docs/repo/plans/plan-register.md` in this source repository.
- [x] Run `git status --short` and record unrelated user changes in `docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_execution_log.md`.
- [x] Verify no planned target path starts with `.codeheart/kit/`.
- [x] Run `python3 scripts/validate-markdown-headers.py` in this source repository.
- [x] Run `git diff --check`.

### G) Implementation Notes

Use the source repository's existing `docs/repo/plans/plan-register.md` format for the entry. Adapt
this handoff plan before execution so source tasks use source-relative paths instead of sibling-repository paths. Do not create additional planning infrastructure in this epic.

### H) Open Questions

- `OQ-1` applies and has a safe default.

## EP1 - Discovery Capability-Scope Handoff

### A) Epic ID, Title, And Outcome

EP1 - Discovery Capability-Scope Handoff

Outcome: the discovery workflow requires a structured implementation capability-scope block before
a discovery can claim normal implementation-handoff readiness for implementation-relevant decision
groups.

### B) Scope

In scope:

- Update implementation-handoff readiness criteria.
- Add a compact capability-scope block template.
- Define when the block is required.
- Define blocked, conditional, and blocker-resolution handoff behavior when the capability scope is
  not concrete enough.
- Preserve the rule that discovery handoff is not an implementation plan with epics or checkboxes.

Out of scope:

- Adding implementation epics to the discovery workflow.
- Requiring the capability-scope block for every early discovery decision group.

### C) Files Touched

```text
./
  components/
    planning-workflows/
      managed/
        runbooks/
          discovery-workflow.md                 # modify
  src/
    codeheart_operating_kit/
      resources/
        components/
          planning-workflows/
            managed/
              runbooks/
                discovery-workflow.md           # mirror source
```

### D) Acceptance Criteria And Size

Size: M.

Acceptance criteria:

- `discovery-workflow.md` says normal implementation handoff for implementation-relevant decision
  groups requires structured capability scope.
- The block fields are exactly: capability, primary workflow, must cover, explicitly out of scope,
  deferred or blocked, preserve decisions, planner must not reinvent, and feature-level success
  evidence.
- The guidance says discovery handoff must not contain implementation epics, execution checklists,
  or sentence-level implementation tasks.
- The guidance says missing capability scope produces blocked, conditional, or blocker-resolution
  handoff instead of normal implementation-handoff readiness.
- Source and packaged resource copies match byte-for-byte.

### E) Dependencies And Critical-Path Notes

Depends on EP0. EP2 will reference the capability-scope block as a preferred source for
implementation-plan coverage.

### F) Tasks Checklist

- [x] Add capability-scope readiness language to `components/planning-workflows/managed/runbooks/discovery-workflow.md` under `Implementation-Handoff Ready`.
- [x] Add the `Implementation Capability Scope - <group name>` template to `components/planning-workflows/managed/runbooks/discovery-workflow.md`.
- [x] Add guidance that the capability-scope block is required for implementation-relevant decision groups included in normal handoff.
- [x] Add guidance that missing capability scope produces blocked, conditional, and blocker-resolution handoff labels.
- [x] Add guidance that discovery handoff must not include implementation epics, execution checklists, and sentence-level implementation tasks.
- [x] Copy `components/planning-workflows/managed/runbooks/discovery-workflow.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/discovery-workflow.md`.
- [x] Run `diff -q components/planning-workflows/managed/runbooks/discovery-workflow.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/discovery-workflow.md`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 -m pytest tests/test_packaging_resources.py -q`.
- [x] Run `git diff --check`.

### G) Implementation Notes

Use the discovery document's accepted block shape as the content source. Keep examples generic and
public-core safe. Avoid tenant names, customer names, private repo paths, raw chat excerpts, and
local machine details.

### H) Open Questions

- `OQ-3` applies and has a safe default.

## EP2 - Draft Implementation Plan Capability And Concreteness Rules

### A) Epic ID, Title, And Outcome

EP2 - Draft Implementation Plan Capability And Concreteness Rules

Outcome: the draft implementation plan workflow makes feature-capability coverage and concrete
implementation tasks first-class drafting requirements while preserving atomic capability-sized
checkboxes.

### B) Scope

In scope:

- Add feature-capability coverage guidance before epic drafting.
- Require explicit treatment of omitted capability as out of scope, deferred, or blocked.
- Strengthen checklist rules around capability-sized tasks.
- Add avoidable non-concreteness guidance.
- Add the fresh implementer test.
- Reject policy-only, doctrine-only, scaffolding-only, stub-only, validation-shell-only, and
  avoidably vague tasks as implementation-ready work.

Out of scope:

- Mandatory actor prefixes.
- Approval-only task taxonomy.
- Deliverable-type taxonomy.
- Sentence-level task decomposition.
- Weak-plan sample artifacts or anti-pattern catalogs.

### C) Files Touched

```text
./
  components/
    planning-workflows/
      managed/
        runbooks/
          draft-implementation-plan.md          # modify
  src/
    codeheart_operating_kit/
      resources/
        components/
          planning-workflows/
            managed/
              runbooks/
                draft-implementation-plan.md    # mirror source
```

### D) Acceptance Criteria And Size

Size: M.

Acceptance criteria:

- The runbook tells planners to derive intended feature capability from accepted discovery when
  available, using the capability-scope block when present.
- The runbook tells planners to use user request and targeted research as capability sources when
  no accepted discovery exists.
- The runbook requires each epic outcome and task set to cover intended capability or explicitly
  mark omitted capability areas as out of scope, deferred, or blocked.
- The checklist rules define capability-sized tasks as concrete implementation slices.
- The checklist rules reject policy-only, doctrine-only, scaffolding-only, stub-only,
  validation-shell-only, and avoidably vague tasks.
- The avoidable non-concreteness rule requires safe planning-time checks for checkable facts.
- Legitimate execution-time variability requires exact preflight, expected result, remediation
  path, retry validation, and stop condition.
- Source and packaged resource copies match byte-for-byte.

### E) Dependencies And Critical-Path Notes

Depends on EP1 so the implementation plan runbook can refer to the discovery capability-scope
block.

### F) Tasks Checklist

- [x] Add feature-capability coverage guidance to `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` near `Inputs` and `Section 2 - Strategy`.
- [x] Add omitted-capability handling guidance to `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` near `Section 3 - Execution Plan`.
- [x] Add capability-sized task guidance to `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` under `Checklist Rules`.
- [x] Add invalid task-shape guidance for policy-only, doctrine-only, scaffolding-only, stub-only, validation-shell-only, and avoidably vague checklist items.
- [x] Add the fresh implementer test to `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [x] Add avoidable non-concreteness guidance with planning-time fact checks plus exact preflight, expected result, remediation path, retry validation, and stop condition.
- [x] Preserve existing atomic checkbox, linear epic, and banned branching-word rules in `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [x] Copy `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [x] Run `diff -q components/planning-workflows/managed/runbooks/draft-implementation-plan.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 -m pytest tests/test_packaging_resources.py -q`.
- [x] Run `git diff --check`.

### G) Implementation Notes

Avoid making the runbook say every task needs a delivery-type label. That was rejected as a weak
solution because a drafter can label policy, evidence, schema, or validation work while still
under-planning the feature capability.

Do not add weak-plan samples or anti-pattern catalogs in this pass. Express the desired task shape
as drafting rules, review checks, and execution review criteria. Use future real plan reviews to
decide whether examples or lint are worth adding later.

### H) Open Questions

- `OQ-3` applies and has a safe default.

## EP3 - Planning Review Capability Coverage Checks

### A) Epic ID, Title, And Outcome

EP3 - Planning Review Capability Coverage Checks

Outcome: the planning-document review workflow flags implementation plans that can pass structure
checks while failing to plan the intended feature capability.

### B) Scope

In scope:

- Add review checks comparing implementation docs against accepted discovery when available.
- Add review checks using the user request and targeted research when no discovery exists.
- Add checks for quiet capability narrowing.
- Add checks for support-structure-only plans.
- Add checks for avoidable non-concreteness.
- Add the lazy implementer risk question.

Out of scope:

- Full architectural review replacement.
- Review automation.
- Mandatory reviewer checklist expansion beyond planning-document review scope.

### C) Files Touched

```text
./
  components/
    planning-workflows/
      managed/
        runbooks/
          review-planning-document.md           # modify
  src/
    codeheart_operating_kit/
      resources/
        components/
          planning-workflows/
            managed/
              runbooks/
                review-planning-document.md     # mirror source
```

### D) Acceptance Criteria And Size

Size: M.

Acceptance criteria:

- Review guidance asks whether the implementation plan concretely implements intended feature
  capability, not only surrounding policy, scaffolding, gates, schemas, stubs, or validation
  shells.
- Review guidance compares against accepted discovery, including goals, non-goals, accepted
  decisions, blockers, and capability scope when available.
- Review guidance flags quiet capability narrowing unless omitted areas are explicitly out of
  scope, deferred, or blocked.
- Review guidance includes the lazy implementer question: could a lazy implementer complete every
  checklist item while delivering only scaffolding, policy, stubs, validation shells, or a narrow
  slice that does not fulfill the intended feature capability?
- Review guidance flags avoidable non-concreteness when a checkable fact lacks a planning-time
  result, exact preflight, expected result, remediation path, retry validation, and stop condition.
- Source and packaged resource copies match byte-for-byte.

### E) Dependencies And Critical-Path Notes

Depends on EP2 because review wording should enforce the same drafting contract.

### F) Tasks Checklist

- [x] Add feature-capability coverage checks to `components/planning-workflows/managed/runbooks/review-planning-document.md` under `Execution Readiness`.
- [x] Add accepted-discovery comparison checks to `components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Add quiet-capability-narrowing checks to `components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Add support-structure-only checks for policy, scaffolding, gates, schemas, stubs, and validation shells.
- [x] Add the lazy implementer risk question to `components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Add avoidable non-concreteness checks with planning-time fact checks plus exact preflight, expected result, remediation path, retry validation, and stop condition.
- [x] Copy `components/planning-workflows/managed/runbooks/review-planning-document.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Run `diff -q components/planning-workflows/managed/runbooks/review-planning-document.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 -m pytest tests/test_packaging_resources.py -q`.
- [x] Run `git diff --check`.

### G) Implementation Notes

Keep reviewer guidance concise enough that architectural review remains the main review workload.
The point is not a second implementation-planning runbook. The point is a clear test that catches
plans which are structurally compliant but not feature-capability complete.

### H) Open Questions

- `OQ-3` applies and has a safe default.

## EP4 - Execution Per-Epic Delivered-Capability Review

### A) Epic ID, Title, And Outcome

EP4 - Execution Per-Epic Delivered-Capability Review

Outcome: the implementation execution workflow's per-epic review gate checks delivered feature
capability against the epic outcome and discovery capability scope when available.

### B) Scope

In scope:

- Extend the per-epic review gate checks.
- Add feature-capability finding handling.
- Preserve existing execution contract that epic outcome owns completion.
- Avoid a broad new stop rule unless feature-capability gaps require plan amendment.

Out of scope:

- Changing goal lifecycle mechanics.
- Changing execution-log schema.
- Adding a standalone execution reviewer workflow.

### C) Files Touched

```text
./
  components/
    planning-workflows/
      managed/
        runbooks/
          execute-implementation-plan.md        # modify
  src/
    codeheart_operating_kit/
      resources/
        components/
          planning-workflows/
            managed/
              runbooks/
                execute-implementation-plan.md  # mirror source
```

### D) Acceptance Criteria And Size

Size: S.

Acceptance criteria:

- Per-epic review checks include delivered feature capability.
- Reviewer compares implementation against epic outcome, acceptance criteria, and discovery
  capability scope when referenced by the plan.
- Reviewer flags completed checklists that leave the intended capability incomplete, narrow,
  policy-only, stubbed, unusable, or unvalidated.
- Material feature-capability gaps are treated as material review findings.
- The runbook says to fix in scope, and to amend the plan when the gap requires new high-impact
  decisions or scope expansion.
- Source and packaged resource copies match byte-for-byte.

### E) Dependencies And Critical-Path Notes

Depends on EP3 so execution review enforces the same capability model as planning review.

### F) Tasks Checklist

- [x] Add delivered feature-capability checks to `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` under `Per-Epic Review Gate`.
- [x] Add discovery capability-scope comparison guidance to `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`.
- [x] Add material finding handling for incomplete, narrow, policy-only, stubbed, unusable, and unvalidated feature capability.
- [x] Preserve the existing execution contract that epic outcome owns completion above literal checkbox completion.
- [x] Copy `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/execute-implementation-plan.md`.
- [x] Run `diff -q components/planning-workflows/managed/runbooks/execute-implementation-plan.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/execute-implementation-plan.md`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 -m pytest tests/test_packaging_resources.py -q`.
- [x] Run `git diff --check`.

### G) Implementation Notes

This epic should be a narrow wording change. Do not add a separate execution stop rule unless the
source wording cannot make feature-capability gaps material findings cleanly.

### H) Open Questions

- `OQ-3` applies and has a safe default.

## EP5 - Packaging, Release Notes, Validation, And Consumer Impact

### A) Epic ID, Title, And Outcome

EP5 - Packaging, Release Notes, Validation, And Consumer Impact

Outcome: the changed managed workflow files are mirrored into packaged resources, documented as an
instruction-only consumer-impact change, and validated with the source repository's smallest
sufficient validation set.

### B) Scope

In scope:

- Verify all changed source managed runbooks match packaged resource mirrors.
- Update release notes with consumer impact and no-migration note.
- Run source validators and focused tests.
- Record consumer-impact evidence in the source plan.
- Prepare release notes and manifests for EP6 release execution.

Out of scope:

- Version bump, release asset build, installer updates, checksum updates, Git tag creation, and
  GitHub release publication for EP6.

### C) Files Touched

```text
./
  components/
    planning-workflows/
      managed/
        runbooks/
          discovery-workflow.md                 # already modified
          draft-implementation-plan.md          # already modified
          review-planning-document.md           # already modified
          execute-implementation-plan.md        # already modified
  src/
    codeheart_operating_kit/
      resources/
        components/
          planning-workflows/
            managed/
              runbooks/
                discovery-workflow.md           # mirror verified
                draft-implementation-plan.md    # mirror verified
                review-planning-document.md     # mirror verified
                execute-implementation-plan.md  # mirror verified
  tests/
    test_packaging_resources.py                 # inspect, existing parity coverage
  release-notes.md                              # modify
```

### D) Acceptance Criteria And Size

Size: M.

Acceptance criteria:

- Each changed managed runbook has a byte-identical packaged resource mirror.
- `release-notes.md` describes the planning workflow quality change.
- Release notes classify the change as `instruction-only change`.
- Release notes state no forced migration and normal update/sync adoption.
- Release manifests are ready for the target `v0.1.7` release.
- Consumer Impact Record in this plan is accurate at completion.
- `python3 scripts/validate-markdown-headers.py` passes.
- `python3 scripts/validate-public-core.py` passes.
- `python3 -m pytest tests/test_packaging_resources.py tests/test_public_core.py tests/test_markdown_headers.py -q` passes.
- Full `python3 -m pytest -q` passes, or the execution log records why it could not run and the
  residual risk.
- `git diff --check` passes.

### E) Dependencies And Critical-Path Notes

Depends on EP1 through EP4. This epic should run after all managed runbook wording is final.

### F) Tasks Checklist

- [x] Run `diff -q components/planning-workflows/managed/runbooks/discovery-workflow.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/discovery-workflow.md`.
- [x] Run `diff -q components/planning-workflows/managed/runbooks/draft-implementation-plan.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [x] Run `diff -q components/planning-workflows/managed/runbooks/review-planning-document.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Run `diff -q components/planning-workflows/managed/runbooks/execute-implementation-plan.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/execute-implementation-plan.md`.
- [x] Inspect `tests/test_packaging_resources.py` and confirm changed planning workflow runbooks remain in `test_changed_source_and_packaged_resources_match`.
- [x] Update `release-notes.md` with the implementation-planning quality summary, `instruction-only change` impact, no forced migration, and normal update/sync adoption path.
- [x] Update root and packaged release manifest surfaces for the target `v0.1.7` release.
- [x] Update this plan's Consumer Impact Record with final affected paths and validation evidence.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 -m pytest tests/test_packaging_resources.py tests/test_public_core.py tests/test_markdown_headers.py -q`.
- [x] Run `python3 -m pytest -q`.
- [x] Run `git diff --check`.
- [x] Record final validation evidence and residual risk in `docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_execution_log.md`.

### G) Implementation Notes

Keep release publication in EP6. This epic prepares managed source docs, packaged resources,
release notes, and manifest surfaces. `docs/repo/runbooks/release-operating-kit.md` owns the
release procedure used by EP6.

### H) Open Questions

- `OQ-2` applies and has a safe default.
- `OQ-3` applies and has a safe default.

## EP6 - Release And Consumer Sync Proof

### A) Epic ID, Title, And Outcome

EP6 - Release And Consumer Sync Proof

Outcome: Codeheart Operating Kit `v0.1.7` is published through the source release runbook, and a first consumer repository syncs the released kit and verifies zero managed-content drift.

### B) Scope

In scope:

- Release version surfaces for `v0.1.7`.
- Release asset build, checksums, tag, and GitHub release through the source release runbook.
- First-consumer update-check, sync, and check proof.
- Release URL, checksum, sync, and drift evidence in the source execution log.

Out of scope:

- Syncing additional consumer repositories beyond the first-consumer proof.
- Changing consumer-owned plans, runbooks, or local guidance outside Operating Kit managed sync.
- Adding planning-lint validators.

### C) Files Touched

```text
./
  pyproject.toml                                      # modify
  src/codeheart_operating_kit/__init__.py             # modify
  manifest.yaml                                       # modify
  src/codeheart_operating_kit/resources/manifest.yaml  # modify
  bootstrap.md                                        # modify
  install.sh                                          # modify
  install.ps1                                         # modify
  release-notes.md                                    # verify
  dist/
    codeheart-operating-kit-0.1.7-macos.tar.gz         # create
    codeheart-operating-kit-0.1.7-macos.tar.gz.sha256  # create
    codeheart-operating-kit-0.1.7-windows.zip          # create
    codeheart-operating-kit-0.1.7-windows.zip.sha256   # create
<first-consumer-repository>/
  .codeheart/kit.lock.yaml                             # modify through sync
  .codeheart/kit/docs/planning-workflows/runbooks/     # modify through sync
```

### D) Acceptance Criteria And Size

Size: L.

Acceptance criteria:

- Source release surfaces consistently target `0.1.7` and `v0.1.7`.
- Release assets and checksum files exist in `dist/` for macOS and Windows.
- `docs/repo/runbooks/release-operating-kit.md` stop conditions are satisfied.
- Git tag `v0.1.7` exists on the validated source commit.
- GitHub release `v0.1.7` is published with release notes, manifests, installers, assets, and
  checksums.
- First consumer repository `codeheart-operating-kit update-check` sees `v0.1.7`.
- First consumer repository `codeheart-operating-kit sync .` updates the managed kit snapshot.
- First consumer repository `codeheart-operating-kit check` reports `OK: True` and `Drift findings: 0`.

### E) Dependencies And Critical-Path Notes

Depends on EP5. This is the public release and first-consumer adoption proof, so it should run only
after source managed wording, packaged resources, release notes, and local validation are complete.

### F) Tasks Checklist

- [x] Read `docs/repo/runbooks/release-operating-kit.md` in this source repository.
- [x] Confirm target release `v0.1.7` against `pyproject.toml`, `src/codeheart_operating_kit/__init__.py`, `manifest.yaml`, `src/codeheart_operating_kit/resources/manifest.yaml`, `bootstrap.md`, `install.sh`, `install.ps1`, and `release-notes.md`.
- [x] Update source version surfaces to `0.1.7` and `v0.1.7` in `pyproject.toml`, `src/codeheart_operating_kit/__init__.py`, `manifest.yaml`, `src/codeheart_operating_kit/resources/manifest.yaml`, `bootstrap.md`, `install.sh`, and `install.ps1`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-json-schemas.py`.
- [x] Run `python3 scripts/validate-release-manifest.py`.
- [x] Run `python3 -m pytest -q`.
- [x] Run `python3 scripts/build-release-assets.py --version 0.1.7 --output-dir dist`.
- [x] Run `python3 -m pytest tests/test_install_metadata.py tests/test_release_assets.py tests/test_packaging_resources.py tests/test_sync_check.py tests/test_json_schemas.py -q`.
- [x] Create Git tag `v0.1.7` from the validated source commit.
- [x] Publish GitHub release `v0.1.7` with `bootstrap.md`, `install.sh`, `install.ps1`, `release-notes.md`, `manifest.yaml`, `dist/codeheart-operating-kit-0.1.7-macos.tar.gz`, `dist/codeheart-operating-kit-0.1.7-macos.tar.gz.sha256`, `dist/codeheart-operating-kit-0.1.7-windows.zip`, and `dist/codeheart-operating-kit-0.1.7-windows.zip.sha256`.
- [x] Record release URL, Git tag, asset filenames, checksum values, and residual risk in `docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_execution_log.md`.
- [x] Run `codeheart-operating-kit update-check` in `../<first-consumer-repository>`.
- [x] Run `codeheart-operating-kit sync .` in `../<first-consumer-repository>`.
- [x] Run `codeheart-operating-kit check` in `../<first-consumer-repository>`.
- [x] Run `git -C ../<first-consumer-repository> diff --check`.
- [x] Record first consumer repository update-check, sync, check, and diff evidence in `docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_execution_log.md`.

### G) Implementation Notes

The current source repository version at planning time is `0.1.6`; `v0.1.7` is the intended patch
release for this plan. If source release numbering advances before EP6 starts, update the source
plan, release surfaces, asset names, and consumer sync evidence to the next patch version before
publishing.

First-consumer sync should be performed through `codeheart-operating-kit sync .`, not by hand-editing
`.codeheart/kit/`. The consumer worktree will show managed snapshot changes after sync; that is
expected adoption evidence, not drift.

### H) Open Questions

- `OQ-2` applies and has a safe default.

## 3.1 Consumer Impact Record

Maintain this record during source execution.

| Impact class or category | Affected paths | Required validation | Release/adoption note | Known consumer action |
| --- | --- | --- | --- | --- |
| `instruction-only change` | `components/planning-workflows/managed/runbooks/discovery-workflow.md`, `draft-implementation-plan.md`, `review-planning-document.md`, `execute-implementation-plan.md`, packaged mirrors under `src/codeheart_operating_kit/resources/`, and release surfaces for `v0.1.7` | Markdown headers, public-core hygiene, packaged-resource parity, focused tests, full pytest, release manifest validation, asset build, `git diff --check` | Release notes describe capability-scope discovery handoff, concrete feature-capability planning, planning-review checks, per-epic delivered-capability review, and no forced migration | Existing consumers update to the release and run normal sync/check |
| First-consumer adoption proof | `.codeheart/kit.lock.yaml` and managed `.codeheart/kit/docs/planning-workflows/runbooks/` files in `../<first-consumer-repository>` | `codeheart-operating-kit update-check`, `codeheart-operating-kit sync .`, `codeheart-operating-kit check`, `git diff --check` | This repository proves released managed guidance can be consumed without managed-content drift | The first consumer repository adopts the released kit through normal sync |
| No forced migration | Existing consumer plans, discovery docs, execution logs, local runbooks, and consumer-owned docs | No generated path, schema, validator, or CLI behavior changed in this pass | Release notes state this improves future managed guidance and does not rewrite existing consumer planning docs | None beyond normal update/sync |

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Add dedicated planning-lint validators after repeated real-world plan reviews make stable machine
  checks possible.
- Extend capability-scope blocks beyond normal implementation handoff candidates only after the
  updated discovery workflow has been piloted.
- Sync additional consumer repositories after first-consumer adoption proves the release path.

## 4.2 Future Considerations

- A later kit version may turn some concreteness checks into automated planning-doc lint.
- The review runbook may eventually split architecture review and execution-readiness review if
  the combined reviewer workload becomes too broad.
- Existing consumer implementation plans are not automatically invalidated by this change. Apply
  the new guidance when they are materially revised or reviewed for execution.

# Revision Notes

- 2026-06-22: Created draft implementation plan from the completed implementation-planning quality
  discovery, targeting Codeheart Operating Kit source managed planning workflow runbooks.
- 2026-06-22: Added source release and first-consumer sync proof, made source-plan copy
  path adaptation explicit, and named the source execution log path.
- 2026-06-22: Activated plan for implementation.
- 2026-06-22: Removed weak-plan sample work from the implementation scope while preserving release
  and consumer sync planning.
- 2026-06-22: Published `v0.1.7` and recorded first-consumer sync proof.
