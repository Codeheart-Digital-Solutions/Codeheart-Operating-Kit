Last updated: 2026-06-25T13:45:59Z (UTC)
Created: 2026-06-25
Status: completed
Completed: 2026-06-25
Execution log: docs/repo/plans/discovery-handoff-gate/discovery-handoff-gate_execution_log.md

# Document Header

## Discovery Handoff Gate Implementation Plan

Overview: Add a small Operating Kit planning guardrail so implementation-plan drafting cannot
quietly proceed from a discovery document whose implementation capability scope has not been
approved, delegated, or explicitly revised by the user. The change is intentionally narrow: it
updates managed planning workflow guidance, prepares and publishes a `v0.1.13` instruction-only
release, then syncs the released kit into the named consumer repositories `Codeheart-HQ`,
`Codeheart-Automation-Foundry`, and the AWS platform consumer repository identified by the user.

Embedded approval: the user explicitly approved the public release and the three named consumer
repository syncs as part of this plan. During execution, the implementing agent must not ask again
for release publication approval or for approval to run normal Operating Kit sync/check updates in
the named consumer repositories. This approval is limited to this plan's files, the `v0.1.13`
release, and normal managed Operating Kit update paths. Stop before publication or sync when a
release runbook stop condition is hit, a destructive action is required, an unplanned consumer
path would be changed, or the target version changes in a way that needs a new release decision.

Essential context:

| Source | Why it matters |
| --- | --- |
| `docs/repo/plans/operation-routing-dispatch-standard/operation-routing-dispatch-standard_discovery_doc.md` | Triggering discovery whose handoff correction exposed the missing planning guardrail. |
| `AGENTS.md` | Source-repository public-core safety, task routing, and release authority rules. |
| `README.md` | Public repository purpose, consumer-owned boundary, and maintainer entry points. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure before changing managed kit docs, resources, tests, or release assets. |
| `docs/repo/runbooks/release-operating-kit.md` | Required procedure for release validation, assets, checksums, tag, GitHub release, and evidence. |
| `docs/repo/reference/placement-contract.md` | Managed content placement and consumer boundary rules. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact classification and release-note requirements for managed instruction changes. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | Source managed discovery workflow that defines implementation-handoff readiness and capability-scope blocks. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Source managed implementation-planning workflow that must enforce the discovery handoff gate. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | Source managed review workflow that should catch implementation plans drafted from unapproved discovery handoffs. |
| `components/planning-workflows/component.yaml` | Planning-workflows manifest version and managed file list. |
| `src/codeheart_operating_kit/resources/` | Packaged resource mirror installed consumers receive. |
| `tests/test_packaging_resources.py` | Existing parity test for source managed files and packaged resource mirrors. |
| `release-notes.md` | Consumer-facing release-note surface for the instruction-only release. |
| `pyproject.toml`, `src/codeheart_operating_kit/__init__.py`, `manifest.yaml`, `bootstrap.md`, `install.sh`, `install.ps1` | Release metadata and installer surfaces that must move consistently to `0.1.13`. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Update Operating Kit planning workflow guidance so an implementation plan derived from a discovery
document must pass a clear discovery-handoff preflight before normal epics are drafted.

Completion is proven when:

- `draft-implementation-plan.md` tells agents to inspect discovery status, capability-scope
  blocks, and user approval/delegation state before drafting normal implementation epics from a
  discovery document;
- the runbook stops normal implementation-plan drafting when a discovery document is only
  manual-review-ready, lacks required `Implementation Capability Scope` blocks, or carries
  unresolved implementation-shaping blockers;
- simple implementation planning without a discovery document remains allowed when the user
  request and repository research provide enough capability scope;
- `discovery-workflow.md` cross-references the drafting gate without adding implementation-plan
  shapes or runbook-maturity doctrine;
- `review-planning-document.md` catches implementation plans that cite discovery but lack approved
  capability-scope handoff evidence;
- source managed planning-workflow files and packaged resource mirrors match;
- release notes classify the change as an `instruction-only change` with no forced consumer
  migration beyond normal update and sync;
- `v0.1.13` release assets, manifests, installers, checksums, tag, and GitHub release are
  published through `docs/repo/runbooks/release-operating-kit.md`;
- `Codeheart-HQ`, `Codeheart-Automation-Foundry`, and the AWS platform consumer repository are
  synced and checked against the released kit through normal Operating Kit CLI commands.

## 1.2 Project And Problem Context

During operation-routing discovery, the discovery document was briefly treated as implementation
handoff-ready before its formal `Implementation Capability Scope` block had been reviewed and
approved by the user. The user correctly identified the process gap: the discovery workflow says
implementation handoff needs capability scope, but the implementation-plan drafting workflow needs
an explicit guardrail that prevents a future agent from skipping that approval boundary.

The fix belongs in Operating Kit planning workflow guidance because it is generic. It is not part
of the M365 module, not part of the routing standard itself, and not a new runbook-authoring shape.
The routing-standard implementation plan should later be drafted against this improved guardrail.

## 1.3 Current State Analysis

Current source state:

- `discovery-workflow.md` already defines implementation-handoff readiness and the
  `Implementation Capability Scope - <group name>` block.
- `draft-implementation-plan.md` already asks planners to use accepted discovery, identify
  feature capability, and cover the capability surface.
- `review-planning-document.md` already checks whether implementation plans preserve discovery
  goals, decisions, blockers, and capability scope.
- The gap is sequencing: the drafting runbook does not explicitly say to stop when discovery
  capability scope exists but has not been accepted, delegated, or revised by the user.
- The packaged managed resources mirror source files under `src/codeheart_operating_kit/resources/`.
- Current release surfaces are at `0.1.12`, so this plan targets `0.1.13`.

Target state:

- A future agent sees the handoff gate inside the implementation-plan drafting runbook before it
  starts epics.
- Discovery can still stop at manual-review-ready when that is the right output.
- Implementation planning from discovery proceeds only from approved/delegated/revised capability
  scope, a normal implementation-handoff-ready discovery, or a deliberately scoped
  blocker-resolution handoff.
- The operation-routing discovery remains separate and can later receive its own implementation
  plan after the user approves or revises its capability scope.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement the guardrail in the planning workflow source, mirror packaged resources, prepare a
patch release, publish it under the embedded approval, then sync the named consumer repositories.

Expected source tree:

```text
Codeheart-Operating-Kit/
  docs/
    README.md                                                       # modify
    repo/
      README.md                                                     # modify
      plans/
        README.md                                                   # modify
        plan-register.md                                            # modify
        discovery-handoff-gate/
          discovery-handoff-gate_implementation_doc.md              # create
  components/
    planning-workflows/
      component.yaml                                                # modify version only
      managed/
        runbooks/
          discovery-workflow.md                                     # modify cross-reference
          draft-implementation-plan.md                              # modify gate
          review-planning-document.md                               # modify review check
  src/
    codeheart_operating_kit/
      __init__.py                                                   # modify release version
      resources/
        manifest.yaml                                               # modify release metadata
        components/
          planning-workflows/
            component.yaml                                          # modify mirror
            managed/
              runbooks/
                discovery-workflow.md                               # modify mirror
                draft-implementation-plan.md                        # modify mirror
                review-planning-document.md                         # modify mirror
  tests/
    test_packaging_resources.py                                     # inspect existing coverage
    test_install_metadata.py                                        # release validation
    test_release_assets.py                                          # release validation
    fixtures/
      release-manifest.json                                         # modify release fixture
  bootstrap.md                                                      # modify release version
  install.sh                                                        # modify release version
  install.ps1                                                       # modify release version
  manifest.yaml                                                     # modify release metadata
  pyproject.toml                                                    # modify release version
  release-notes.md                                                  # modify release notes
  scripts/build-release-assets.py                                   # modify default version
  dist/                                                             # create release assets
```

Consumer sync targets:

```text
Codeheart-HQ/                                                       # sync managed kit
Codeheart-Automation-Foundry/                                       # sync managed kit
AWS platform consumer repository/                                   # sync managed kit
```

## 2.2 Open Questions And Assumptions Requiring Clarification

`OQ-1` - Should the target release be `v0.1.13`?

- `BLOCKER: no`
- `Affects: EP-03`, `EP-04`
- Unlocks exact release metadata, asset names, tag name, and consumer sync target.
- Recommended default: use `v0.1.13` because the source repository currently exposes `0.1.12`.
  If execution finds a newer released source state, update every release surface to the next patch
  version and record the reason in the execution log.

`OQ-2` - Is release and named consumer sync approval already granted?

- `BLOCKER: no`
- `Affects: EP-03`, `EP-04`
- Unlocks public tag/GitHub release publication and normal sync/check in the three named consumer
  repositories without a second approval prompt.
- Recommended default: yes. The user approved release publication and the three named consumer
  syncs in the request that created this plan. The executor must still stop on release runbook
  stop conditions, unplanned destructive action, unplanned consumer paths, or validation failure.

`OQ-3` - Should this plan add fresh-agent routing probes?

- `BLOCKER: no`
- `Affects: EP-01`, `EP-02`
- Unlocks whether the small guardrail change also implements routing-standard validation doctrine.
- Recommended default: no. Fresh low-context routing probes belong in the later operation-routing
  standard implementation plan. This plan only prevents implementation planning from bypassing
  discovery capability-scope approval.

## 2.3 Architectural Decisions With Reasoning

`AD-1` - Put the gate in `draft-implementation-plan.md`

1. Problem being solved: implementation plans are where a premature handoff causes concrete
   execution work to be drafted from incomplete discovery state.
2. Simplest working solution: add a pre-drafting discovery-handoff gate to the drafting runbook.
3. What may change in 6-12 months: the gate may later become a checklist, validator, or planning
   metadata field when planning documents become more structured.
4. Rationale: the current risk happens before epics are created, so the drafting runbook is the
   direct control point.
5. Alternatives considered: only strengthening discovery workflow was rejected because discovery
   already names the requirement and the failure happened during transition into planning.

`AD-2` - Keep discovery workflow change as a cross-reference

1. Problem being solved: discovery authors should know the drafting runbook will enforce the
   handoff boundary.
2. Simplest working solution: add a short cross-reference in `discovery-workflow.md`.
3. What may change in 6-12 months: discovery status may become machine-readable, making the
   cross-reference less important.
4. Rationale: this avoids duplicating implementation-plan drafting rules inside discovery.
5. Alternatives considered: adding a second capability-scope shape was rejected because it would
   create two standards for the same handoff.

`AD-3` - Review catches the same failure mode

1. Problem being solved: an implementation plan may already exist before a reviewer sees it.
2. Simplest working solution: add one review check for discovery-derived plans.
3. What may change in 6-12 months: review checklists may be split by planning document type.
4. Rationale: review should catch a missing handoff gate even when drafting missed it.
5. Alternatives considered: relying only on the drafting runbook was rejected because review is
   the normal backstop for planning quality.

`AD-4` - Treat the release as instruction-only

1. Problem being solved: the change alters managed guidance but does not alter generated paths,
   validators, schemas, CLI sync behavior, or consumer-owned content.
2. Simplest working solution: classify the change as `instruction-only change`, record release
   notes, and use normal sync adoption.
3. What may change in 6-12 months: a future structured planning validator may require a different
   impact class.
4. Rationale: this release changes agent instructions and packaged managed docs only.
5. Alternatives considered: adding a validator now was rejected because the desired gate is still
   guidance-level and the document structure is Markdown.

`AD-5` - Sync all three named consumer repositories after release

1. Problem being solved: the user wants the improved kit adopted in `Codeheart-HQ`,
   `Codeheart-Automation-Foundry`, and the AWS platform consumer repository as part of the same
   execution.
2. Simplest working solution: install/use the released `0.1.13` CLI, run update-check, sync, and
   check in each named repository.
3. What may change in 6-12 months: a portfolio-wide sync helper may automate repeated consumer
   updates.
4. Rationale: explicit named sync avoids leaving the improved planning guardrail only in the
   source release.
5. Alternatives considered: syncing only one proof consumer was rejected because the user named
   three concrete repositories.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Planning workflow source files enforce and review the discovery handoff gate. | S | none |
| EP-02 | Packaged resources, indexes, register entries, and local validation are updated. | M | EP-01 |
| EP-03 | `v0.1.13` release assets are prepared and published under embedded approval. | M | EP-02 |
| EP-04 | The three named consumer repositories are synced and checked against the released kit. | M | EP-03 |

## EP-01 - Discovery Handoff Gate In Managed Planning Workflows

### A) Epic ID, Title, And Outcome

EP-01 - Discovery Handoff Gate In Managed Planning Workflows

Outcome: Managed planning workflow guidance prevents normal implementation-plan drafting from an
unapproved or incomplete discovery handoff and gives reviewers a matching check.

### B) Scope

Update only the planning workflow runbooks needed for the handoff gate. Do not implement the
operation-routing standard, route cards, fresh-agent routing probes, or new runbook-maturity
shapes in this epic.

### C) Files Touched

```text
components/planning-workflows/managed/runbooks/discovery-workflow.md        # modify
components/planning-workflows/managed/runbooks/draft-implementation-plan.md # modify
components/planning-workflows/managed/runbooks/review-planning-document.md  # modify
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- `draft-implementation-plan.md` has a discovery-handoff preflight before normal epic drafting.
- The preflight distinguishes accepted discovery, manual-review-ready discovery, blocked handoff,
  conditional handoff, and no-discovery straightforward work.
- The preflight requires approved, delegated, or revised `Implementation Capability Scope` blocks
  for discovery-derived normal implementation plans.
- `discovery-workflow.md` briefly points to the drafting runbook as the enforcement point.
- `review-planning-document.md` flags implementation plans that cite discovery without approved
  capability-scope handoff evidence.
- No detailed route-card, routing-probe, or runbook-maturity shapes are added.

### E) Dependencies And Critical-Path Notes

No upstream implementation dependencies. This epic should complete before the later
operation-routing-dispatch-standard implementation plan is drafted.

### F) Tasks Checklist

- [x] Add a `Discovery Handoff Preflight` subsection to `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [x] State that discovery-derived normal implementation plans require capability-scope handoff evidence with approved, delegated, revised status.
- [x] State that manual-review-ready discovery without approved capability scope stops normal epic drafting.
- [x] State that blocked handoff permits blocker-resolution plans and conditional handoff permits explicitly scoped partial plans.
- [x] State that straightforward work without discovery can proceed from user request, repository research, and recorded assumptions.
- [x] Add one cross-reference sentence in `components/planning-workflows/managed/runbooks/discovery-workflow.md`.
- [x] Add one discovery-derived-plan review check in `components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Run `python3 scripts/validate-markdown-headers.py` after editing source runbooks.

### G) Implementation Notes

Keep wording short. The drafting runbook should route agents at the moment they would otherwise
start Section 3 epics. Avoid duplicating the full capability-scope block shape, because
`discovery-workflow.md` already owns it.

### H) Open Questions

OQ-3 applies.

## EP-02 - Packaged Resources, Indexes, Register, And Local Validation

### A) Epic ID, Title, And Outcome

EP-02 - Packaged Resources, Indexes, Register, And Local Validation

Outcome: Source managed files are mirrored into packaged resources, repository planning indexes
and the plan register expose this plan, and local validation proves the instruction-only change.

### B) Scope

Update packaging mirrors and planning discoverability. Validate the changed source tree before
release preparation.

### C) Files Touched

```text
docs/README.md
docs/repo/README.md
docs/repo/plans/README.md
docs/repo/plans/plan-register.md
components/planning-workflows/component.yaml
src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml
src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/discovery-workflow.md
src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/draft-implementation-plan.md
src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/review-planning-document.md
tests/test_packaging_resources.py
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Packaged planning-workflow mirrors match source managed files byte-for-byte.
- Planning-workflows component manifests remain valid and identify the instruction-only impact.
- `docs/README.md`, `docs/repo/README.md`, `docs/repo/plans/README.md`, and
  `docs/repo/plans/plan-register.md` expose this implementation plan.
- `tests/test_packaging_resources.py` already covers the changed source and mirror files.
- Markdown headers, public-core hygiene, packaged-resource parity, and whitespace validation pass.

### E) Dependencies And Critical-Path Notes

Depends on EP-01. Release preparation in EP-03 must not start until the source and packaged
resources are in parity.

### F) Tasks Checklist

- [x] Mirror `components/planning-workflows/managed/runbooks/discovery-workflow.md` into `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/discovery-workflow.md`.
- [x] Mirror `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` into `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/draft-implementation-plan.md`.
- [x] Mirror `components/planning-workflows/managed/runbooks/review-planning-document.md` into `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/review-planning-document.md`.
- [x] Update `components/planning-workflows/component.yaml` version to `0.1.13`.
- [x] Mirror `components/planning-workflows/component.yaml` into `src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml`.
- [x] Confirm `tests/test_packaging_resources.py` covers all changed planning-workflow source files.
- [x] Add this plan to `docs/README.md`.
- [x] Add this plan to `docs/repo/README.md`.
- [x] Add this plan to `docs/repo/plans/README.md`.
- [x] Add `OK-PR-012` for this plan to `docs/repo/plans/plan-register.md`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 scripts/validate-json-schemas.py`.
- [x] Run `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_packaging_resources.py tests/test_sync_check.py tests/test_install_metadata.py tests/test_release_assets.py`.
- [x] Run `git diff --check`.

### G) Implementation Notes

The implementation plan itself is repository governance content and is not a managed consumer
resource. Keep public release notes free of private local paths and raw consumer logs.

### H) Open Questions

None.

## EP-03 - Release Preparation And Publication

### A) Epic ID, Title, And Outcome

EP-03 - Release Preparation And Publication

Outcome: `v0.1.13` is prepared, validated, tagged, and published with release assets and checksums
through the Operating Kit release runbook.

### B) Scope

Update release metadata, release notes, package version surfaces, installers, manifests, tests,
fixtures, and release assets. Publish the public release under the embedded approval after release
runbook stop conditions pass.

### C) Files Touched

```text
release-notes.md
pyproject.toml
src/codeheart_operating_kit/__init__.py
scripts/build-release-assets.py
bootstrap.md
install.sh
install.ps1
manifest.yaml
src/codeheart_operating_kit/resources/manifest.yaml
tests/fixtures/release-manifest.json
dist/
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Version surfaces consistently target `0.1.13`.
- Release notes describe the discovery-handoff gate as an instruction-only managed planning
  workflow change.
- Root release manifest records real publishable asset checksums.
- Packaged resource manifest uses the established zero-placeholder downloadable asset checksums
  inside release assets.
- `dist/` contains macOS and Windows release assets plus checksum files.
- Release runbook stop conditions are checked before publication.
- Public tag `v0.1.13` exists at the validated commit.
- GitHub release `v0.1.13` includes release notes, manifest, installers, assets, and checksums.

### E) Dependencies And Critical-Path Notes

Depends on EP-02. The embedded approval in this plan authorizes EP-03 release publication without
another user prompt. Stop on validation failure, unreproducible assets, missing checksums, dirty
validated commit mismatch, or release runbook stop condition.

### F) Tasks Checklist

- [x] Re-read `docs/repo/runbooks/release-operating-kit.md`.
- [x] Update `release-notes.md` with `v0.1.13` instruction-only release notes.
- [x] Update `pyproject.toml` package version to `0.1.13`.
- [x] Update `src/codeheart_operating_kit/__init__.py` package version to `0.1.13`.
- [x] Update `scripts/build-release-assets.py` default release version to `0.1.13`.
- [x] Update `bootstrap.md`, `install.sh`, and `install.ps1` release references to `v0.1.13`.
- [x] Update `manifest.yaml` release metadata and URLs for `v0.1.13`.
- [x] Update `src/codeheart_operating_kit/resources/manifest.yaml` release metadata and URLs for `v0.1.13`.
- [x] Update `tests/fixtures/release-manifest.json` to match `v0.1.13`.
- [x] Run `python3 scripts/validate-release-manifest.py manifest.yaml`.
- [x] Run `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest`.
- [x] Run `uv run --with pip --with setuptools --with wheel python scripts/build-release-assets.py --version 0.1.13 --output-dir dist`.
- [x] Run `git diff --check`.
- [x] Confirm `git diff --cached --name-status` contains only intended release and plan changes.
- [x] Commit the validated release changes.
- [x] Create public tag `v0.1.13` from the validated commit.
- [x] Publish GitHub release `v0.1.13` with `bootstrap.md`, `install.sh`, `install.ps1`, `release-notes.md`, `manifest.yaml`, release assets, and checksum files.
- [x] Record release URL, asset names, checksums, validation commands, and residual risk in the execution log.

### G) Implementation Notes

The release should stay a patch release. The executor may use established `gh` CLI or GitHub
connector paths already used by prior Operating Kit releases, while preserving the release
runbook stop conditions.

### H) Open Questions

OQ-1 and OQ-2 apply.

## EP-04 - Named Consumer Repository Sync

### A) Epic ID, Title, And Outcome

EP-04 - Named Consumer Repository Sync

Outcome: `Codeheart-HQ`, `Codeheart-Automation-Foundry`, and the AWS platform consumer repository
have their managed Operating Kit snapshots refreshed to the released `0.1.13` kit and pass
managed-content checks.

### B) Scope

Run normal Operating Kit update-check, sync, and check commands in the three named local consumer
repositories after release publication. Do not edit consumer-owned files outside normal managed
Operating Kit sync behavior.

### C) Files Touched

```text
Codeheart-HQ/.codeheart/kit/                                      # modify managed sync output
Codeheart-HQ/.codeheart/kit.lock.yaml                             # modify generated lock state
Codeheart-HQ/AGENTS.md                                            # modify managed block only
Codeheart-Automation-Foundry/.codeheart/kit/                      # modify managed sync output
Codeheart-Automation-Foundry/.codeheart/kit.lock.yaml             # modify generated lock state
Codeheart-Automation-Foundry/AGENTS.md                            # modify managed block only
AWS platform consumer repository/.codeheart/kit/                  # modify managed sync output
AWS platform consumer repository/.codeheart/kit.lock.yaml         # modify generated lock state
AWS platform consumer repository/AGENTS.md                        # modify managed block only
docs/repo/plans/discovery-handoff-gate/discovery-handoff-gate_execution_log.md # create at execution
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- The local CLI used for sync reports `codeheart-operating-kit 0.1.13`.
- Update-check, sync, and check run in `Codeheart-HQ`.
- Update-check, sync, and check run in `Codeheart-Automation-Foundry`.
- Update-check, sync, and check run in the AWS platform consumer repository.
- Each repo reports no managed-content drift after sync.
- Consumer sync changes are limited to normal managed Operating Kit paths and managed root
  `AGENTS.md` block updates.
- Execution log records command outcomes, changed-path summary, validation evidence, and residual
  risk without private local paths or raw operational logs.

### E) Dependencies And Critical-Path Notes

Depends on EP-03. The embedded approval in this plan authorizes normal update-check, sync, and
check in the three named consumer repositories without another user prompt. Stop before any
destructive action, consumer-owned content rewrite, unrelated cleanup, or manual edit outside
managed sync output.

### F) Tasks Checklist

- [x] Install the published `v0.1.13` CLI through the validated macOS installer path.
- [x] Run `codeheart-operating-kit --version` and confirm `codeheart-operating-kit 0.1.13`.
- [x] Run `codeheart-operating-kit update-check <Codeheart-HQ checkout> --agent-notification`.
- [x] Run `codeheart-operating-kit sync <Codeheart-HQ checkout>`.
- [x] Run `codeheart-operating-kit check <Codeheart-HQ checkout> --json`.
- [x] Run `codeheart-operating-kit update-check <Codeheart-Automation-Foundry checkout> --agent-notification`.
- [x] Run `codeheart-operating-kit sync <Codeheart-Automation-Foundry checkout>`.
- [x] Run `codeheart-operating-kit check <Codeheart-Automation-Foundry checkout> --json`.
- [x] Run `codeheart-operating-kit update-check <AWS platform consumer checkout> --agent-notification`.
- [x] Run `codeheart-operating-kit sync <AWS platform consumer checkout>`.
- [x] Run `codeheart-operating-kit check <AWS platform consumer checkout> --json`.
- [x] Inspect `git status --short` in each named consumer repository.
- [x] Confirm each named consumer repository changed only normal managed Operating Kit paths.
- [x] Create `docs/repo/plans/discovery-handoff-gate/discovery-handoff-gate_execution_log.md` with release and consumer sync evidence.

### G) Implementation Notes

Use local checkout paths resolved at execution time, but do not write absolute local paths into
public release notes. Consumer commits are outside this plan unless the user separately requests
repository commit or pull-request publication.

### H) Open Questions

OQ-2 applies.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- The full operation-routing-dispatch standard remains deferred to its own implementation plan
  after the user approves or revises the routing discovery capability scope.
- Fresh low-context routing probes for routing-bearing runbook epics remain deferred to the
  operation-routing standard implementation plan.
- Machine-readable discovery status or capability-scope validation remains deferred until the
  Markdown planning workflow shows repeated failure modes that justify parser work.

## 4.2 Future Considerations

- Consider adding a small planning-review prompt template after several discovery-derived plans
  have used this gate.
- Consider a portfolio sync helper if repeated multi-repository Operating Kit updates become
  common.

# Revision Notes

- 2026-06-25: Created first draft implementation plan for the discovery handoff gate, release,
  and named consumer repository sync.
- 2026-06-25: Activated plan for goal-style execution and created sibling execution log.
- 2026-06-25: Replaced the EP-03 worktree-wide status check with a staged-diff check so unrelated
  untracked files remain protected during release commit preparation.
- 2026-06-25: Completed release publication, consumer repository sync, final review, and
  close-out evidence.
