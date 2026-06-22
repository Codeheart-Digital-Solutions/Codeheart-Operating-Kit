Last updated: 2026-06-22T21:13:19Z (UTC)
Created: 2026-06-22
Status: active
Execution log: docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_execution_log.md

# Document Header

## Plan Register Portfolio Doctrine Refinement Implementation Plan

Overview: Refine existing Operating Kit plan-register doctrine so local registers and
coordination-home registers have durable reference shapes for portfolio overview use without
turning the register into a board, task tracker, archive system, or consumer-specific planning
surface.

This plan does not introduce the consumer-local Portfolio Work Board pattern into the Operating
Kit. It only updates reusable plan-register reference and maintenance doctrine that is stable
enough for managed public-core guidance.

Essential context:

| Source | Why it matters |
| --- | --- |
| `<consumer-repository>/docs/repo/plans/portfolio-planning-surfaces/portfolio-planning-surfaces_discovery_doc.md` | Accepted consumer-local discovery handoff for register doctrine refinements, sanitized for public-core implementation. |
| `AGENTS.md` | Operating Kit maintainer routing, public-core safety, and release boundaries. |
| `README.md` | Public repository purpose and managed-content ownership boundary. |
| `docs/README.md` | Public docs router that must expose this plan once created. |
| `docs/repo/README.md` | Repository-governance router that must expose this plan once created. |
| `docs/repo/plans/README.md` | Plan index that must link this implementation plan. |
| `docs/repo/runbooks/change-operating-kit.md` | Required procedure before changing kit source docs, managed content, packaged resources, tests, or release assets. |
| `docs/repo/runbooks/promote-consumer-change.md` | Required procedure for promoting consumer-local guidance into public reusable kit doctrine. |
| `docs/repo/runbooks/release-operating-kit.md` | Required procedure before publishing a new Operating Kit release. |
| `docs/repo/reference/placement-contract.md` | Managed content, kit-initialized consumer state, and consumer-owned boundary rules. |
| `docs/repo/reference/consumer-impact-classification.md` | Consumer-impact class and release-note requirement for managed instruction changes. |
| `components/planning-workflows/managed/reference/plan-register-format.md` | Source managed field contract, lifecycle grouping, relation vocabulary, session refs, and coordination-home ID guidance. |
| `components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Source managed procedure for local register updates, coordination-home updates, and pending sync fallback. |
| `components/planning-workflows/managed/reference/planning-document-lifecycle.md` | Planning lifecycle and register relationship reference that may need a compact cross-reference update. |
| `components/planning-workflows/managed/README.md` | Component route that must remain accurate after reference and runbook changes. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md` | Packaged resource mirror installed consumers receive. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Packaged resource mirror installed consumers receive. |
| `tests/test_packaging_resources.py` | Existing parity test for changed source and packaged resources. |
| `release-notes.md` | Consumer-facing release-note surface for the instruction-only patch release. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Update Codeheart Operating Kit source managed plan-register doctrine so future consumers can use
one register model for local repositories and configured coordination homes while preserving clear
ownership boundaries.

Implementation completion is proven when:

- `plan-register-format.md` contains public-safe reference shapes for local entries,
  coordination-home member entries, repository-qualified IDs, canonical document pointers, relation
  fields, and compact repository grouping;
- `plan-register-format.md` clarifies that `docs/repo/plans/plan-register.md` remains the
  default durable register location even when canonical planning documents live in product,
  module, business, or other repository-owned docs roots;
- `maintain-plan-register.md` tells agents how to apply those shapes during local updates,
  configured coordination-home updates, and pending-sync fallback;
- the source managed docs and packaged resource mirrors match byte-for-byte;
- release notes classify the change as an `instruction-only change` with no forced consumer
  migration;
- local validation covers Markdown headers, public-core hygiene, packaged-resource parity, focused
  tests, release manifest validation, release asset build readiness, and whitespace;
- release preparation is complete and any public release publication happens through
  `docs/repo/runbooks/release-operating-kit.md` after explicit release approval.

## 1.2 Project And Problem Context

Operating Kit already defines `docs/repo/plans/plan-register.md` as the default formal plan
register and `docs/repo/plans/coordination-sync-pending.md` as the fallback queue for unavailable
coordination-home updates. Recent plan-register releases added lifecycle grouping, session
reference handling, and coordination-home ID uniqueness.

The remaining gap is not basic location guidance. The remaining gap is the durable reference shape
for portfolio-scale register usage:

- local repositories need a clear entry shape that works whether a canonical plan lives under
  `docs/repo/plans/`, `docs/business/plans/`, product docs, module docs, or another
  repository-owned planning root;
- coordinated portfolios may prefer repository-qualified local IDs so local registers and the
  coordination-home register use the same stable ID for the same entry;
- coordination homes need examples for local portfolio entries, represented member entries,
  relation fields, and repository grouping;
- agents need to know which relation details belong in the canonical plan and which compact
  relation pointers belong in the register;
- the register must remain a durable index, not a current-focus board, issue tracker, or execution
  checklist.

Consumer-local discovery also identified a useful current-focus surface, named a Portfolio Work
Board in that consumer context. That pattern is not mature enough for Operating Kit promotion in
this implementation. It remains consumer-local until real maintenance cycles prove the shape.

## 1.3 Current State Analysis

Current source state:

- `components/planning-workflows/managed/reference/plan-register-format.md` already defines
  required fields, lifecycle values, one default register, lifecycle grouping, relation vocabulary,
  session refs, coordination notes, and anti-patterns.
- The same reference already says local registers and coordination-home registers use the same
  entry shape, with selective coverage for coordination homes.
- `components/planning-workflows/managed/reference/plan-register-format.md` already says
  coordination-home entry IDs must be unique and gives namespace guidance for member entries.
- `components/planning-workflows/managed/runbooks/maintain-plan-register.md` already gives local
  register, portfolio coordination, pending sync, session reference, and safety procedures.
- `components/planning-workflows/managed/reference/planning-document-lifecycle.md` already points
  planning lifecycle changes back to the plan register.
- Source and packaged managed docs are mirrored under
  `src/codeheart_operating_kit/resources/components/planning-workflows/`.
- The source repository is currently at version `0.1.8`.

Target state:

- The format reference contains concrete local and coordination-home Markdown entry examples.
- The format reference explains repository-qualified local IDs as an allowed coordinated-portfolio
  convention without forcing every standalone consumer repository to use them.
- The maintenance runbook tells agents to reuse an already repository-qualified member ID in the
  coordination home instead of double-prefixing it.
- The runbook tells agents to derive `<SOURCE-NAMESPACE>-<LOCAL-ID>` only when the member local ID
  is bare.
- The reference explains that canonical documents own detail, while the register owns compact
  index metadata and relation pointers.
- The change ships as a small instruction-only patch release, likely `v0.1.9` when execution
  starts.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Make a narrow managed-doc refinement in the planning-workflows component, mirror it into packaged
resources, update release notes and version surfaces, then validate and prepare a patch release.

Expected source tree:

```text
Codeheart-Operating-Kit/
  docs/
    README.md                                                            # modify route
    repo/
      README.md                                                          # modify route
      plans/
        README.md                                                        # modify route
        plan-register.md                                                 # modify entry
        plan-register-portfolio-doctrine-refinement/
          plan-register-portfolio-doctrine-refinement_implementation_doc.md  # create
          plan-register-portfolio-doctrine-refinement_execution_log.md       # create at activation
  components/
    planning-workflows/
      component.yaml                                                     # modify during release prep
      managed/
        README.md                                                        # inspect
        reference/
          plan-register-format.md                                        # modify
          planning-document-lifecycle.md                                 # modify compact cross-reference
        runbooks/
          maintain-plan-register.md                                      # modify
  src/
    codeheart_operating_kit/
      resources/
        components/
          planning-workflows/
            component.yaml                                               # modify mirror during release prep
            managed/
              reference/
                plan-register-format.md                                  # modify mirror
                planning-document-lifecycle.md                           # modify mirror
              runbooks/
                maintain-plan-register.md                                # modify mirror
  tests/
    test_packaging_resources.py                                          # inspect existing parity coverage
    test_install_metadata.py                                             # release validation
    test_release_assets.py                                               # release validation
    test_sync_check.py                                                   # focused managed sync validation
    test_json_schemas.py                                                 # manifest and schema validation
  release-notes.md                                                       # modify
  manifest.yaml                                                          # modify during release prep
  src/codeheart_operating_kit/resources/manifest.yaml                     # modify during release prep
  pyproject.toml                                                         # modify during release prep
  src/codeheart_operating_kit/__init__.py                                 # modify during release prep
  scripts/build-release-assets.py                                        # modify default version during release prep
  bootstrap.md                                                           # modify release URLs during release prep
  install.sh                                                             # modify release URLs during release prep
  install.ps1                                                            # modify release URLs during release prep
```

## 2.2 Open Questions And Assumptions Requiring Clarification

OQ-1 - Target release version

- `BLOCKER: no`
- `Affects: EP-04, EP-05`
- Unlocks exact release notes, version surfaces, release asset names, and tag name.
- Recommended default: use `v0.1.9` because source `pyproject.toml` and release notes currently
  show `0.1.8` as the latest release version.

OQ-2 - Public release publication timing

- `BLOCKER: no`
- `Affects: EP-05`
- Unlocks public tag creation, GitHub release publication, and consumer sync proof.
- Recommended default: draft and validate release assets during implementation, then require an
  explicit release-publication instruction before creating the public tag and GitHub release.

OQ-3 - Scope of lifecycle reference edits

- `BLOCKER: no`
- `Affects: EP-01, EP-02`
- Unlocks whether `planning-document-lifecycle.md` receives a compact cross-reference update.
- Recommended default: add only one compact cross-reference when the format and runbook changes
  would otherwise be hard to find from lifecycle docs.

## 2.3 Architectural Decisions With Reasoning

AD-1 - Refine the existing register format reference, not a new runbook

1. Problem being solved: The durable issue is entry shape and relation semantics, not a separate
   maintenance procedure.
2. Simplest working solution: Put reference shapes in `plan-register-format.md` and update
   `maintain-plan-register.md` only where procedure changes.
3. What may change in 6-12 months: A future validator may turn stable entry-shape conventions into
   machine checks.
4. Rationale: The format reference already owns field shape, examples, lifecycle grouping, relation
   vocabulary, and anti-patterns.
5. Alternatives considered and why not chosen: A separate coordination-home register runbook would
   duplicate the existing maintenance procedure and make local versus coordination-home behavior
   drift.

AD-2 - Keep one managed register model for local and coordination-home registers

1. Problem being solved: Coordination homes need broader coverage than a single local repository,
   but the entry shape can remain the same.
2. Simplest working solution: Preserve one format and one maintenance runbook, then document
   different coverage and ID conventions inside them.
3. What may change in 6-12 months: Large coordination homes may need additional scanning aids or
   validation, but not a different register contract.
4. Rationale: One model keeps consumer installs simpler and lets agents transfer local register
   knowledge to coordination-home work.
5. Alternatives considered and why not chosen: A separate coordination-home schema would increase
   instruction load before the current Markdown register proves it needs that complexity.

AD-3 - Allow repository-qualified local IDs for coordinated portfolios

1. Problem being solved: A member local ID such as `PR-001` is easy locally but ambiguous in a
   portfolio conversation.
2. Simplest working solution: Document repository-qualified IDs such as `EXAMPLE-AUTOMATION-PR-001`
   as the preferred convention for coordinated portfolios, while standalone repositories may keep
   short IDs.
3. What may change in 6-12 months: A future CLI may allocate IDs from configuration and warn about
   collisions.
4. Rationale: A repository-qualified local ID can be reused by the coordination home without
   mental mapping and without double-prefixing.
5. Alternatives considered and why not chosen: Forcing all consumers to change local IDs would be
   a migration. Keeping only bare local IDs leaves portfolio-facing work noisier.

AD-4 - Keep canonical documents as detail authority

1. Problem being solved: Registers can become dense and stale when they duplicate discovery,
   implementation, execution, and decision detail.
2. Simplest working solution: State that the register owns compact index metadata and relation
   pointers, while canonical docs own decisions, blockers, evidence, execution state, and detailed
   relation rationale.
3. What may change in 6-12 months: A future summary surface may aggregate selected current work,
   but that should not change source-of-truth ownership.
4. Rationale: The register remains scannable and refreshable from canonical documents.
5. Alternatives considered and why not chosen: Storing relation rationale in the register would
   make the register compete with canonical planning docs.

AD-5 - Exclude the current-focus board pattern from this Operating Kit release

1. Problem being solved: A consumer-local board pattern is useful but not yet proven enough for
   reusable managed doctrine.
2. Simplest working solution: Keep this implementation limited to register doctrine and explicitly
   defer board doctrine.
3. What may change in 6-12 months: The board pattern may be promoted after real consumer-local
   maintenance cycles prove stable admission, exit, and grouping rules.
4. Rationale: Prematurely promoting a board pattern could blur the register versus current-focus
   boundary.
5. Alternatives considered and why not chosen: Adding a managed work-board runbook now would ship
   unproven workflow doctrine.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-00 | Plan activation context, execution log, and public-core implementation context are ready. | S | None |
| EP-01 | Source plan-register format reference has durable portfolio entry shapes. | M | EP-00 |
| EP-02 | Source maintenance runbook applies the refined register shapes procedurally. | M | EP-01 |
| EP-03 | Packaged resources and route indexes mirror the source doctrine changes. | S | EP-02 |
| EP-04 | Release notes, version surfaces, and validation prepare an instruction-only patch release. | M | EP-03 |
| EP-05 | Public release and first-consumer sync proof are completed after explicit approval. | M | EP-04 |

## EP-00 - Activation Context And Execution Log

### A) Epic ID, Title, And Outcome

EP-00 - Activation Context And Execution Log

Outcome: The draft plan is activated after explicit approval, execution evidence has a local home,
public context is fresh, and managed-doc edits can proceed without consumer-private leakage.

### B) Scope

In scope:

- Set this plan to `active` after explicit approval.
- Create the execution log.
- Refresh the local plan-register entry for activation.
- Confirm public-core and change-runbook context before managed-doc edits.

Out of scope:

- Editing managed plan-register doctrine.
- Publishing a release.

### C) Files Touched

```text
docs/repo/plans/plan-register.md                                          # modify
docs/repo/plans/plan-register-portfolio-doctrine-refinement/
  plan-register-portfolio-doctrine-refinement_implementation_doc.md       # modify status
  plan-register-portfolio-doctrine-refinement_execution_log.md            # create
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- The implementation plan status changes from `draft` to `active` only after explicit approval.
- `docs/repo/plans/plan-register.md` refreshes `OK-PR-004` to active state.
- The execution log exists beside this plan.
- The plan uses public-safe placeholder references instead of private consumer repository names.
- `git diff --check` passes.

### E) Dependencies And Critical-Path Notes

No dependencies. This epic is completed by drafting and registering the plan.

### F) Tasks Checklist

- [x] Read `AGENTS.md`, `README.md`, `docs/repo/runbooks/change-operating-kit.md`, `docs/repo/runbooks/promote-consumer-change.md`, `docs/repo/reference/placement-contract.md`, and `docs/repo/reference/consumer-impact-classification.md`.
- [x] Update `docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_implementation_doc.md` status from `draft` to `active` after explicit approval.
- [x] Create `docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_execution_log.md` with activation context and validation sections.
- [x] Refresh `OK-PR-004 - Plan Register Portfolio Doctrine Refinement` in `docs/repo/plans/plan-register.md` with active lifecycle state.
- [x] Run `git status --short` and record unrelated user changes in `docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_execution_log.md`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `git diff --check`.

### G) Implementation Notes

This plan is already the public-core source-repository planning record. Execution starts only after
explicit approval. Do not copy private repository names, local absolute paths, customer names,
tenant names, account IDs, credentials, or restricted business strategy into the execution log.

### H) Open Questions

None.

## EP-01 - Format Reference Shapes

### A) Epic ID, Title, And Outcome

EP-01 - Format Reference Shapes

Outcome: `plan-register-format.md` gives future agents concrete, public-safe reference shapes for
local entries, coordination-home member entries, repository-qualified IDs, canonical pointers,
relations, and grouping.

### B) Scope

In scope:

- Add local and coordination-home entry examples.
- Clarify repository-qualified IDs for coordinated portfolios.
- Clarify already-qualified ID reuse in the coordination home.
- Clarify canonical docs pointer forms for local and member entries.
- Clarify compact relation ownership and repository grouping.
- Add a narrow lifecycle-reference cross-link when useful.

Out of scope:

- Adding a validator.
- Adding a new register file location.
- Adding current-focus board doctrine.

### C) Files Touched

```text
components/planning-workflows/managed/reference/plan-register-format.md        # modify
components/planning-workflows/managed/reference/planning-document-lifecycle.md # modify compact cross-reference
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- The reference keeps `docs/repo/plans/plan-register.md` as the default durable register.
- The reference says canonical planning docs may live outside `docs/repo/plans/`.
- The reference includes a local entry example with a repo-relative canonical docs pointer.
- The reference includes a coordination-home member entry example with an explicit
  repository/path pointer.
- The reference documents repository-qualified local IDs for coordinated portfolios.
- The reference documents reuse of already repository-qualified IDs in coordination homes.
- The reference documents namespace derivation for bare member local IDs.
- The reference clarifies relation ownership between register pointers and canonical-plan detail.
- The reference keeps board, backlog, and task tracker anti-pattern boundaries intact.

### E) Dependencies And Critical-Path Notes

Depends on EP-00. The format reference defines the target shape before the maintenance runbook
procedures are updated.

### F) Tasks Checklist

- [x] Update `components/planning-workflows/managed/reference/plan-register-format.md` `Coverage` text with local repository coverage and coordination-home portfolio overview coverage.
- [x] Update `components/planning-workflows/managed/reference/plan-register-format.md` `Entry Fields` text with canonical pointer guidance for repo-relative paths and explicit repository/path pointers.
- [x] Add repository-qualified local ID guidance to `components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Add already-qualified coordination-home ID reuse guidance to `components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Add bare-member-ID namespace derivation guidance to `components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Add a local entry Markdown example to `components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Add a coordination-home member entry Markdown example to `components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Add compact relation ownership guidance to `components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Add repository grouping guidance to `components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Add a compact plan-register shape cross-reference to `components/planning-workflows/managed/reference/planning-document-lifecycle.md`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `git diff --check`.

### G) Implementation Notes

Use neutral examples such as `EXAMPLE-AUTOMATION-PR-001` and
`Example-Automation:docs/repo/plans/example/example_implementation_doc.md`. Preserve the existing
anti-pattern that the register is not a task backlog, sprint board, transcript index, or source of
truth for lifecycle detail.

### H) Open Questions

OQ-3 applies and has a safe default.

## EP-02 - Maintenance Procedure Refinement

### A) Epic ID, Title, And Outcome

EP-02 - Maintenance Procedure Refinement

Outcome: `maintain-plan-register.md` gives agents an exact procedure for applying the refined
local and coordination-home entry shapes.

### B) Scope

In scope:

- Update local register procedure for repository-qualified ID conventions.
- Update portfolio member procedure for already-qualified ID reuse and bare-ID namespace
  derivation.
- Update coordination-home procedure for selected member updates and relation pointers.
- Clarify pending-sync entry identity expectations.
- Preserve no-sibling-scan safety.

Out of scope:

- Adding automated sync application.
- Adding register validation logic.
- Changing `.codeheart/kit.config.yaml` schema.

### C) Files Touched

```text
components/planning-workflows/managed/runbooks/maintain-plan-register.md  # modify
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Local register procedure tells coordinated portfolios to follow configured local ID conventions.
- Member procedure reuses an already repository-qualified source ID in the coordination home.
- Member procedure derives `<SOURCE-NAMESPACE>-<LOCAL-ID>` when the source local ID is bare.
- Coordination-home procedure applies the same rule for explicitly requested member updates.
- Pending-sync guidance carries enough identity detail for later coordination-home application.
- Safety rules still prohibit private details, implicit sibling scans, and unsafe coordination-home
  writes.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 so the procedure can reference the final format language.

### F) Tasks Checklist

- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` local register procedure with repository-qualified ID convention guidance.
- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` member portfolio procedure with already-qualified ID reuse steps.
- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` member portfolio procedure with bare-ID namespace derivation steps.
- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` coordination-home procedure with the same member ID handling rule.
- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` relation guidance for represented entries and explicit repository/path pointers.
- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` pending-sync shape guidance with source local ID and coordination-home target ID details.
- [x] Verify `components/planning-workflows/managed/runbooks/maintain-plan-register.md` still prohibits implicit sibling repository scans.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `git diff --check`.

### G) Implementation Notes

Keep the runbook procedural. Do not duplicate every Markdown example from the format reference.
Route detailed examples back to `../reference/plan-register-format.md`.

### H) Open Questions

None.

## EP-03 - Packaged Resources And Routing

### A) Epic ID, Title, And Outcome

EP-03 - Packaged Resources And Routing

Outcome: Installed consumers receive the refined doctrine because packaged resources mirror the
changed source docs and public route indexes stay discoverable.

### B) Scope

In scope:

- Mirror changed source managed files into packaged resources.
- Confirm parity coverage.
- Update component routes only when changed wording creates a discoverability gap.

Out of scope:

- Release version bump.
- Release asset build.
- Consumer sync proof.

### C) Files Touched

```text
src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md        # modify mirror
src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/planning-document-lifecycle.md # modify mirror
src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md       # modify mirror
components/planning-workflows/managed/README.md                                                                    # inspect
src/codeheart_operating_kit/resources/components/planning-workflows/managed/README.md                               # modify mirror when source route changes
tests/test_packaging_resources.py                                                                                  # inspect
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- Changed source managed files match packaged resource mirrors byte-for-byte.
- `tests/test_packaging_resources.py` covers all changed mirrored files.
- Planning-workflows README routes remain accurate.
- No generated consumer state files are changed by hand.

### E) Dependencies And Critical-Path Notes

Depends on EP-02 because packaged resources mirror final source wording.

### F) Tasks Checklist

- [x] Copy `components/planning-workflows/managed/reference/plan-register-format.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Copy `components/planning-workflows/managed/reference/planning-document-lifecycle.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/planning-document-lifecycle.md`.
- [x] Copy `components/planning-workflows/managed/runbooks/maintain-plan-register.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md`.
- [x] Inspect `components/planning-workflows/managed/README.md` for route accuracy.
- [x] Mirror `components/planning-workflows/managed/README.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/README.md` when source route text changes.
- [x] Inspect `tests/test_packaging_resources.py` and verify each changed managed file is covered by parity tests.
- [x] Run `uv run --with pytest python -m pytest tests/test_packaging_resources.py`.
- [x] Run `git diff --check`.

### G) Implementation Notes

Prefer exact file copy for mirrors. Packaged resources must not diverge from source managed docs.

### H) Open Questions

None.

## EP-04 - Release Preparation And Validation

### A) Epic ID, Title, And Outcome

EP-04 - Release Preparation And Validation

Outcome: The repository is internally consistent and validated for an instruction-only patch
release carrying the refined plan-register doctrine.

### B) Scope

In scope:

- Update release notes.
- Update version surfaces for the target patch release.
- Run focused and full validation.
- Build release assets locally.
- Create the execution log after plan activation.

Out of scope:

- Public tag creation.
- GitHub release publication.
- Consumer sync proof.

### C) Files Touched

```text
docs/repo/plans/plan-register-portfolio-doctrine-refinement/
  plan-register-portfolio-doctrine-refinement_execution_log.md       # create at activation
release-notes.md                                                     # modify
components/planning-workflows/component.yaml                         # modify
src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml  # modify
manifest.yaml                                                        # modify
src/codeheart_operating_kit/resources/manifest.yaml                  # modify
pyproject.toml                                                       # modify
src/codeheart_operating_kit/__init__.py                              # modify
scripts/build-release-assets.py                                      # modify
bootstrap.md                                                         # modify
install.sh                                                           # modify
install.ps1                                                          # modify
tests/fixtures/release-manifest.json                                 # modify when fixture tracks release version
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `release-notes.md` contains a `v0.1.9` section for plan-register portfolio doctrine refinement.
- Release notes classify the change as `instruction-only change`.
- Release notes state no forced migration and normal update/sync adoption.
- Version surfaces consistently use the target patch version.
- Local validation passes before publication.
- Release assets build locally with the target patch version.
- The execution log records validation evidence and residual risk.

### E) Dependencies And Critical-Path Notes

Depends on EP-03. If release numbering advances before execution, amend the target version before
editing release surfaces.

### F) Tasks Checklist

- [x] Create `docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_execution_log.md` after activation.
- [x] Add `v0.1.9` release notes to `release-notes.md` for plan-register portfolio doctrine refinement.
- [x] Record `instruction-only change` consumer impact in `release-notes.md`.
- [x] Record no forced migration and normal sync adoption in `release-notes.md`.
- [x] Update `pyproject.toml` package version to `0.1.9`.
- [x] Update `src/codeheart_operating_kit/__init__.py` package version to `0.1.9`.
- [x] Update `scripts/build-release-assets.py` default release version to `0.1.9`.
- [x] Update `components/planning-workflows/component.yaml` component version to `0.1.9`.
- [x] Update `src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml` component version to `0.1.9`.
- [x] Update `manifest.yaml` release metadata and URLs for `v0.1.9`.
- [x] Update `src/codeheart_operating_kit/resources/manifest.yaml` release metadata and URLs for `v0.1.9`.
- [x] Update `bootstrap.md`, `install.sh`, and `install.ps1` release URLs for `v0.1.9`.
- [x] Update `tests/fixtures/release-manifest.json` to match the `v0.1.9` release manifest.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 scripts/validate-json-schemas.py`.
- [x] Run `python3 scripts/validate-release-manifest.py manifest.yaml`.
- [x] Run `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_packaging_resources.py tests/test_install_metadata.py tests/test_release_assets.py tests/test_sync_check.py tests/test_json_schemas.py`.
- [x] Run `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest`.
- [x] Run `uv run --with pip --with setuptools --with wheel python scripts/build-release-assets.py --version 0.1.9 --output-dir dist`.
- [x] Run `git diff --check`.
- [x] Record validation commands and residual risk in `docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_execution_log.md`.

### G) Implementation Notes

The target version is `v0.1.9` at draft time. If the source repository has already advanced beyond
`0.1.8` when execution starts, update the plan and all release-surface tasks to the next patch
version before editing files.

### H) Open Questions

OQ-1 applies and has a safe default.

## EP-05 - Release Publication And Consumer Sync Proof

### A) Epic ID, Title, And Outcome

EP-05 - Release Publication And Consumer Sync Proof

Outcome: The patch release is published after explicit approval, and a first consumer can adopt
the refined plan-register doctrine through normal update and sync.

### B) Scope

In scope:

- Follow the release runbook for public release publication.
- Verify release assets and checksums.
- Prove update-check, sync, and check in one consumer repository.
- Record release and adoption evidence in the execution log.

Out of scope:

- Syncing every consumer repository.
- Migrating existing consumer-owned plan registers.
- Populating portfolio coordination-home entries.
- Implementing a Portfolio Work Board.

### C) Files Touched

```text
dist/                                                                    # create release assets
docs/repo/plans/plan-register-portfolio-doctrine-refinement/
  plan-register-portfolio-doctrine-refinement_execution_log.md           # modify
```

Consumer sync proof touches an installed consumer's managed `.codeheart/kit/` snapshot through
normal `codeheart-operating-kit update-check`, `sync`, and `check` commands after publication.

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Public tag `v0.1.9` exists at the validated commit.
- GitHub release `v0.1.9` includes release notes, manifest, installers, assets, and checksums.
- Release runbook stop conditions are checked before publication.
- Consumer update-check detects the published version.
- Consumer sync refreshes managed planning-workflow docs.
- Consumer check reports no managed-content drift after sync.
- Execution log records release URL, asset names, checksums, validation evidence, sync evidence,
  and residual risk.

### E) Dependencies And Critical-Path Notes

Depends on EP-04. Public release publication requires explicit release approval at the time this
epic starts.

### F) Tasks Checklist

- [ ] Re-read `docs/repo/runbooks/release-operating-kit.md` before release publication.
- [ ] Confirm the validated commit with `git status --short` and `git rev-parse HEAD`.
- [ ] Confirm `release-notes.md` covers the plan-register portfolio doctrine consumer impact.
- [ ] Confirm `dist/codeheart-operating-kit-0.1.9-macos.tar.gz` exists from the validated asset build.
- [ ] Confirm `dist/codeheart-operating-kit-0.1.9-windows.zip` exists from the validated asset build.
- [ ] Confirm checksum files exist for both release assets.
- [ ] Create public tag `v0.1.9` from the validated commit after explicit release publication approval.
- [ ] Publish GitHub release `v0.1.9` with `bootstrap.md`, `install.sh`, `install.ps1`, `release-notes.md`, `manifest.yaml`, release assets, and checksum files.
- [ ] Run `codeheart-operating-kit update-check` in one consumer repository after publication.
- [ ] Run `codeheart-operating-kit sync <consumer-repository-path>` in the same consumer repository after update-check sees `v0.1.9`.
- [ ] Run `codeheart-operating-kit check <consumer-repository-path> --json` and confirm managed-content drift is absent.
- [ ] Record release publication evidence in `docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_execution_log.md`.
- [ ] Record consumer update-check, sync, and check evidence in `docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_execution_log.md`.

### G) Implementation Notes

Use the existing release workflow and release asset patterns from `v0.1.8`. Do not publish from a
dirty worktree. Do not include private consumer names or local paths in public release notes or
the public execution log.

### H) Open Questions

OQ-2 applies and has a safe default.

## 3.1 Consumer Impact Record

Maintain this record during source execution.

| Impact class or category | Affected paths | Required validation | Release/adoption note | Known consumer action |
| --- | --- | --- | --- | --- |
| `instruction-only change` | `components/planning-workflows/managed/reference/plan-register-format.md`, `components/planning-workflows/managed/runbooks/maintain-plan-register.md`, possible compact update to `components/planning-workflows/managed/reference/planning-document-lifecycle.md`, packaged mirrors under `src/codeheart_operating_kit/resources/`, and release surfaces for `v0.1.9` | Markdown headers, public-core hygiene, packaged-resource parity, focused tests, full pytest, release manifest validation, asset build, `git diff --check` | Release notes describe refined local and coordination-home register reference shapes, no forced migration, and normal sync adoption | Existing consumers update to the release and run normal sync/check |
| No forced migration | Existing consumer-owned `docs/repo/plans/plan-register.md`, `docs/repo/plans/coordination-sync-pending.md`, canonical plans, execution logs, local runbooks, and board-like local surfaces | No generated path, schema, validator, or CLI behavior changed in this pass | Release notes state this improves future managed guidance and does not rewrite existing consumer planning docs | None beyond normal update/sync |

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Operating Kit Portfolio Work Board doctrine is deferred until a consumer-local board pattern has
  real maintenance evidence and stable admission, exit, grouping, and ownership rules.
- Automated validation for register IDs, relation vocabulary, canonical pointers, and source local
  ID traceability is deferred until enough real registers exist to justify parser and fixture work.
- CLI-assisted pending-sync application is deferred. This implementation only makes the manual
  register procedure clearer.
- Migration of existing consumer-owned register IDs is deferred. Existing registers can be
  refreshed when they are next materially updated or when a consumer executes a local migration
  plan.
- Coordination-home population work is deferred to consumer-owned portfolio setup plans.

## 4.2 Future Considerations

- If coordination-home registers grow large, add a validator that checks ID uniqueness, canonical
  docs pointer shape, relation vocabulary, and source-local-ID traceability.
- If the Portfolio Work Board pattern is promoted later, keep it separate from the durable plan
  register and make its admission rule depend on formal plan-register entries or canonical
  planning documents.
- If several consumers adopt repository-qualified local IDs, add a small naming reference for
  namespace design rather than expanding the maintenance runbook.
- If consumers need archive compaction, define a register review runbook before introducing a
  separate archive file.

# Revision Notes

- 2026-06-22: Created draft implementation plan for Operating Kit plan-register portfolio
  doctrine refinement, targeting source managed reference and runbook updates without promoting
  consumer-local work-board doctrine.
- 2026-06-22: Activated implementation plan for execution after user approval.
- 2026-06-22: Completed managed doctrine, packaged resource mirrors, release preparation, and
  local validation through release asset build. Public release publication and consumer sync proof
  remain gated by explicit release-publication approval.
