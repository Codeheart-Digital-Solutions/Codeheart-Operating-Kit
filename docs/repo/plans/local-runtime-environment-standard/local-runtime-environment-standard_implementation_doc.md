Last updated: 2026-06-26T14:38:50Z (UTC)
Created: 2026-06-26
Status: completed

# Document Header

## Local Runtime Environment Standard Implementation Plan

Overview: Implement the accepted Operating Kit standard for ignored local machine/runtime state.
The implementation adds `.codeheart/local/` as the repo-local machine boundary, makes
`.codeheart/local/envs/python/` the default Python virtual environment path, adds safe init/sync
and config/schema behavior, broadens managed tooling-readiness guidance beyond module-only
blockers, and keeps first-run onboarding lightweight.
It also adds a narrow planning-workflow clarification learned from this plan: when work spans
multiple repositories, the canonical discovery or implementation plan belongs in the repository
that owns the work boundary, while coordination-home registers point to that canonical plan.

This plan does not create a virtual environment by default, install Python packages, introduce a
dependency lock strategy, convert every module runbook, create durable readiness state, or publish
a release. It prepares the Operating Kit source so a later release can deliver the standard through
normal sync/update behavior.

Essential context:

| Source | Why it matters |
| --- | --- |
| `docs/repo/plans/local-runtime-environment-standard/local-runtime-environment-standard_discovery_doc.md` | Accepted discovery, decisions D-001 through D-010, and implementation capability scope. |
| `AGENTS.md` | Public-core safety, change routing, and required maintainer references. |
| `README.md` | Public repository purpose and consumer-owned boundary. |
| `docs/README.md` | Top-level docs router that must expose the new implementation plan. |
| `docs/repo/README.md` | Repository-governance router that must expose the plan bundle. |
| `docs/repo/plans/README.md` | Plan index that must link the new implementation plan. |
| `docs/repo/plans/plan-register.md` | Local register that must record the implementation plan. |
| `docs/repo/reference/placement-contract.md` | Source placement contract for installed consumer areas and ownership rules. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact classification for managed docs, schema, generated config, and sync behavior. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure before changing source, docs, schemas, templates, tests, or CLI behavior. |
| `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` | Managed readiness route that must gain the local machine layer and Python venv workflow. |
| `components/agent-interface/managed/reference/local-extension-contract.md` | Defines local user guidance and must distinguish user notes from local machine/runtime state. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Runbook authoring standard that must keep generic local tooling guidance centralized. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Routing standard that must route local tooling blockers after owner and scope are selected. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Recipe-maturity standard for repeatable readiness and venv setup procedure text. |
| `components/agent-interface/managed/runbooks/conduct-first-run-onboarding.md` | Must preserve the no-default-venv first-run rule. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | Discovery drafting route that should name the canonical owner-repository placement check. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Implementation-plan drafting route that should name the canonical owner-repository placement check. |
| `components/planning-workflows/managed/reference/plan-register-format.md` | Existing register doctrine for canonical owner-repo plans and coordination pointers. |
| `components/structure-governance/managed/reference/managed-content-boundaries.md` | Managed boundary doctrine that must add the local machine layer. |
| `components/agent-interface/managed/kit-readme.md` | Installed fallback inventory that should expose local machine/runtime state. |
| `templates/agents/AGENTS.managed-block.md` | Source template for installed root rules, including the tooling-readiness route. |
| `src/codeheart_operating_kit/components.py` | Init/sync source for `.gitignore`, generated config, scaffolds, and first-run behavior. |
| `schemas/kit-config.schema.json` | Consumer config schema that must accept optional `local_machine_layer_path`. |
| `tests/test_init.py` | Init tests for config, gitignore, and absence of default runtime directory creation. |
| `tests/test_sync_check.py` | Sync tests for repairing missing gitignore lines in existing consumers. |
| `tests/test_json_schemas.py` | Schema tests for old and new kit config shapes. |
| `tests/test_packaging_resources.py` | Source and packaged-resource parity checks for managed content. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Implement the local runtime environment standard so installed Operating Kit consumers have one
ignored, documented, schema-visible local machine boundary and one default Python venv convention
that agents can use when local tooling blocks a task.

Implementation completion is proven when:

- `.codeheart/local/` is documented as ignored local machine/runtime state, distinct from
  `.codeheart/user/`;
- `.codeheart/local/envs/python/` is documented as the default repo-local Python venv for
  Operating Kit, Foundry, and module CLI tooling;
- purpose-specific venvs are documented as exceptions under `.codeheart/local/envs/<purpose>/`;
- init and sync add `.codeheart/local/` to consumer `.gitignore`;
- init and sync do not create `.codeheart/local/` or `.codeheart/local/envs/python/` by default;
- new `.codeheart/kit.config.yaml` files include
  `local_consumer_layer.local_machine_layer_path: .codeheart/local/`;
- the config schema accepts existing configs that omit `local_machine_layer_path`;
- the managed tooling-readiness route applies to any agent task blocked by missing local tooling,
  not only module onboarding and operations;
- discovery and implementation drafting guidance tells agents to place canonical plans in the
  repository that owns the work boundary when multiple repositories are involved;
- the `python-runtime` lane gives the shared venv command shape without installing repo-specific
  packages globally;
- managed docs, source mirrors, packaged resources, tests, docs indexes, and plan register
  pointers are consistent;
- validation proves gitignore behavior, schema compatibility, generated config shape, packaged
  resource parity, and no generated runtime state in source control.

## 1.2 Project And Problem Context

AI Execution onboarding exposed a concrete local setup failure: direct editable install into a
Homebrew-managed Python can be blocked by PEP 668. The successful local path was an ignored venv,
but putting that venv under an AI-specific path would make every future Python-backed module
invent its own environment convention.

The accepted Operating Kit direction is broader. The immediate pain came from AI Execution, but
the standard belongs in the Operating Kit because any repo, Foundry module, or agent task can hit
missing local tooling. The shared boundary should be simple enough for nontechnical users,
predictable for agents, and safe for managed repos.

The implementation must keep three boundaries intact:

- `.codeheart/user/` remains human local preferences and notes;
- `.codeheart/local/` becomes generated, ignored, recreatable machine/runtime state;
- module runbooks remain responsible for module-specific package names, package versions, smoke
  commands, auth, and live service preflight.

## 1.3 Current State Analysis

Current source state:

- `src/codeheart_operating_kit/components.py` writes generated config and ensures local-user
  gitignore lines, but it does not add `.codeheart/local/`.
- `schemas/kit-config.schema.json` allows `local_consumer_layer` fields for repo docs, agent
  memory, and user layer only.
- `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` already handles local
  tooling blockers, but its route text is still mainly framed around module onboarding and
  operations.
- The readiness runbook has a `python-runtime` lane, but it does not define the default repo-local
  Python venv path.
- `components/agent-interface/managed/runbooks/conduct-first-run-onboarding.md` already says not
  to create a Python virtual environment during default onboarding.
- `components/planning-workflows/managed/reference/plan-register-format.md` already says
  canonical planning documents may live outside the local register and that coordination-home
  entries can point to member-repository canonical docs.
- `components/planning-workflows/managed/runbooks/discovery-workflow.md` and
  `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` do not yet make
  the multi-repository owner-boundary placement check prominent at plan creation time.
- `components/structure-governance/managed/reference/managed-content-boundaries.md`,
  `docs/repo/reference/placement-contract.md`, and
  `components/agent-interface/managed/reference/local-extension-contract.md` define
  `.codeheart/user/`, but not `.codeheart/local/`.
- Source managed docs are mirrored under `src/codeheart_operating_kit/resources/`, and parity is
  protected by `tests/test_packaging_resources.py`.

Target state:

- Agents can find `.codeheart/local/` from managed docs and generated config when a repo-local
  runtime is needed.
- Existing consumers can sync and gain the `.gitignore` protection without a forced config
  migration.
- New installs expose the local machine layer path in config without creating empty runtime
  directories.
- Tooling readiness remains approval-gated, official-source oriented, and distinct from live
  service preflight.
- Future Foundry modules can depend on the generic standard instead of copying venv placement
  rules.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement the source behavior first, then managed doctrine, then tests and packaged mirrors. Keep
the implementation intentionally small: add the boundary, schema, and runbook guidance, but do not
build an environment manager or create the venv automatically.

Expected source tree:

```text
Codeheart-Operating-Kit/
  docs/
    README.md                                                           # modify
    repo/
      README.md                                                         # modify
      reference/
        placement-contract.md                                           # modify
      plans/
        README.md                                                       # modify
        plan-register.md                                                # modify
        local-runtime-environment-standard/
          local-runtime-environment-standard_discovery_doc.md            # existing
          local-runtime-environment-standard_implementation_doc.md       # create
  components/
    agent-interface/
      managed/
        kit-readme.md                                                   # modify
        reference/
          local-extension-contract.md                                   # modify
          operation-routing-and-dispatch.md                             # modify
          runbook-authoring-standard.md                                 # modify
        runbooks/
          conduct-first-run-onboarding.md                               # review
          handle-tooling-readiness.md                                   # modify
    structure-governance/
      managed/
        reference/
          managed-content-boundaries.md                                 # modify
    planning-workflows/
      managed/
        reference/
          plan-register-format.md                                       # modify
        runbooks/
          discovery-workflow.md                                         # modify
          draft-implementation-plan.md                                  # modify
  schemas/
    kit-config.schema.json                                              # modify
  src/
    codeheart_operating_kit/
      components.py                                                     # modify
      resources/
        components/
          agent-interface/
            managed/
              kit-readme.md                                             # modify mirror
              reference/
                local-extension-contract.md                             # modify mirror
                operation-routing-and-dispatch.md                       # modify mirror
                runbook-authoring-standard.md                           # modify mirror
              runbooks/
                handle-tooling-readiness.md                             # modify mirror
          structure-governance/
            managed/
              reference/
                managed-content-boundaries.md                           # modify mirror
          planning-workflows/
            managed/
              reference/
                plan-register-format.md                                 # modify mirror
              runbooks/
                discovery-workflow.md                                   # modify mirror
                draft-implementation-plan.md                            # modify mirror
        templates/
          agents/
            AGENTS.managed-block.md                                     # modify mirror
  templates/
    agents/
      AGENTS.managed-block.md                                           # modify
  tests/
    fixtures/
      kit-config.yaml                                                   # modify
      validator-valid/
        kit-config-without-purpose.yaml                                 # preserve old-shape fixture
    test_init.py                                                        # modify
    test_sync_check.py                                                  # modify
    test_json_schemas.py                                                # modify
    test_packaging_resources.py                                         # modify
```

No deletion is planned.

## 2.2 Open Questions And Assumptions Requiring Clarification

### OQ-1 - Which Consumer-Impact Class Applies?

BLOCKER: no

Affects: EP-01, EP-02, EP-03, EP-04

Decision unlocked: release-note and validation language for consumers.

Recommended default: record a mixed backwards-compatible impact. The managed-doc updates are
`instruction-only change`; the config schema update is schema/validator-affecting; the init/sync
gitignore update is a backwards-compatible generated-surface behavior change. No forced migration
is required.

### OQ-2 - Should `.codeheart/local/` Be Listed In Profile Generated Surfaces?

BLOCKER: no

Affects: EP-01, EP-03

Decision unlocked: whether profile metadata should list a path that init does not create.

Recommended default: do not list `.codeheart/local/` in `profiles/standard.yaml` generated
surfaces during this implementation. The path is ignored and documented, but it is not generated
during first-run onboarding.

### OQ-3 - Should First-Run Onboarding Text Change?

BLOCKER: no

Affects: EP-02

Decision unlocked: whether the implementation edits `conduct-first-run-onboarding.md`.

Recommended default: review and leave unchanged when its existing "Do not create a Python virtual
environment during default onboarding" rule remains clear. The implementation can validate that
rule without adding duplicate local-runtime wording.

### OQ-4 - Should AI Execution Runbooks Be Patched In This Plan?

BLOCKER: no

Affects: EP-04

Decision unlocked: whether this Operating Kit source plan edits Foundry module content.

Recommended default: no. Record AI Execution onboarding as a downstream Foundry follow-up that
depends on this Operating Kit standard.

### OQ-5 - Should This Plan Clarify Multi-Repo Plan Placement?

BLOCKER: no

Affects: EP-02, EP-03

Decision unlocked: whether this implementation should adjust managed planning workflow guidance
while it already touches cross-repo coordination and placement boundaries.

Recommended default: yes, narrowly. Discovery and implementation drafting runbooks should tell
agents to identify the repository that owns the work boundary before creating the canonical plan.
Coordination-home or HQ registers should point to the canonical plan rather than becoming the
canonical home by default.

## 2.3 Architectural Decisions With Reasoning

### AD-1 - Add `.codeheart/local/` As A New Local Machine Layer

Problem being solved: generated runtime files need a predictable ignored home that is not confused
with human preferences under `.codeheart/user/`.

Simplest working solution: document `.codeheart/local/` and add it to consumer `.gitignore`
during init/sync.

What may change in 6-12 months: this folder may host additional local caches, generated shims, or
purpose-specific runtime environments.

Rationale: the name is broad enough for runtime state without implying durable repo or deployed
environment state.

Alternatives considered: `.codeheart/user/` was rejected because it is human-facing; `.codeheart/envs/`
was rejected because it is too narrow and can conflict with business environment terminology.

### AD-2 - Use One Default Repo-Local Python Venv

Problem being solved: Python-backed agent tooling needs a reliable path that avoids global Python
mutation and works with externally managed Python distributions.

Simplest working solution: use `.codeheart/local/envs/python/` as the default venv path and use
`.codeheart/local/envs/<purpose>/` only for concrete conflicts.

What may change in 6-12 months: high-conflict modules may declare purpose-specific venvs, or a
future tool may automate venv creation.

Rationale: one shared path minimizes cognitive load and keeps modules compatible by default.

Alternatives considered: one venv per module was rejected as premature fragmentation; global
package installs were rejected as unsafe and incompatible with PEP 668.

### AD-3 - Make Config Additive And Backwards Compatible

Problem being solved: agents benefit from a machine-readable local machine layer path, but
existing consumers should not need config migrations.

Simplest working solution: write `local_machine_layer_path: .codeheart/local/` for new installs
and make it optional in `schemas/kit-config.schema.json`.

What may change in 6-12 months: config may add more local runtime metadata after real usage
patterns emerge.

Rationale: additive optional schema preserves existing configs while giving new installs a clear
path.

Alternatives considered: requiring the new field was rejected because it would force unnecessary
consumer migration.

### AD-4 - Keep Readiness Runbook At L1 Structured Recipe Maturity

Problem being solved: venv creation is repeatable operational work, but the first standard should
not freeze a wrapper or CLI before usage patterns are proven.

Simplest working solution: keep command shapes and evidence rules in the managed runbook, with
fresh-agent executability review and focused tests around generated behavior.

What may change in 6-12 months: repeated errors or duplicated command blocks may justify a script
asset or CLI helper.

Rationale: L1 gives enough structure for agents without overbuilding an environment manager.

Alternatives considered: an L2 tested script block or L4 CLI wrapper was rejected for V1 because
the current need is a standard path and approval flow, not automation.

### AD-5 - Treat This As Routing-Bearing Managed Doctrine

Problem being solved: missing local tools should not cause agents to pick an install surface
before selecting the owner, route, and approval class.

Simplest working solution: update root managed guidance, operation-routing examples, and the
tooling-readiness route so agents route local blockers consistently after owner and scope are
known.

What may change in 6-12 months: domains may add route cards that point to the readiness route for
specific module workflows.

Rationale: this keeps route-before-surface behavior intact while making local blockers easier to
resolve.

Alternatives considered: placing Python venv guidance only in AI Execution was rejected because
the standard is cross-module.

### AD-6 - Clarify Plan Placement At Draft Time

Problem being solved: when work spans HQ, Operating Kit, Foundry, and consumer repositories,
agents can create the canonical discovery or implementation plan in the most visible coordination
repo instead of the repository that owns the change.

Simplest working solution: update discovery and implementation drafting guidance to require an
owner-boundary placement check before creating the canonical planning document, and use
plan-register entries for cross-repo coordination pointers.

What may change in 6-12 months: product or module repositories may introduce local planning roots
or route registries that further specialize where module plans live.

Rationale: the canonical plan should live where implementation authority, source changes, and
validation evidence are owned. Coordination homes should coordinate; they should not silently
become source of truth for another repository's work.

Alternatives considered: relying only on plan-register-format was rejected because the rule exists
there but is easy to miss during initial drafting.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Init, sync, config schema, and tests implement the local machine layer behavior. | M | none |
| EP-02 | Managed docs and runbooks define `.codeheart/local/`, Python venv usage, and routing behavior. | M | EP-01 strategy |
| EP-03 | Packaged mirrors, docs indexes, and validation protect the changed source surfaces. | M | EP-01, EP-02 |
| EP-04 | Plan/register coordination, release-note requirement, and downstream Foundry handoff are recorded. | S | EP-03 |

## EP-01 - Generated Behavior And Schema

### A) Epic ID, Title, And Outcome

EP-01 - Generated Behavior And Schema

Outcome: new installs and syncs ignore `.codeheart/local/`, new configs expose the local machine
layer path, old configs remain valid, and no default runtime directory is created.

### B) Scope

In scope:

- `.gitignore` behavior for `.codeheart/local/`;
- generated config field for new installs;
- schema compatibility for old and new configs;
- focused init, sync, and schema tests.

Out of scope:

- creating `.codeheart/local/`;
- creating `.codeheart/local/envs/python/`;
- installing Python packages;
- changing profile generated surfaces.

### C) Files Touched

```text
src/codeheart_operating_kit/components.py                 # modify
schemas/kit-config.schema.json                            # modify
tests/test_init.py                                        # modify
tests/test_sync_check.py                                  # modify
tests/test_json_schemas.py                                # modify
tests/fixtures/kit-config.yaml                            # modify
tests/fixtures/validator-valid/kit-config-without-purpose.yaml # review
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `codeheart-operating-kit init` writes `.codeheart/local/` to `.gitignore`;
- `codeheart-operating-kit sync` repairs a missing `.codeheart/local/` ignore line;
- init does not create `.codeheart/local/`;
- new config includes `local_consumer_layer.local_machine_layer_path: .codeheart/local/`;
- schema accepts configs with the field;
- schema accepts old configs without the field;
- schema rejects an incorrect `local_machine_layer_path` value;
- profile generated surfaces do not claim `.codeheart/local/` as a created path.

### E) Dependencies And Critical-Path Notes

This epic is the critical path for implementation because docs should describe behavior that the
CLI actually supports. Keep `.codeheart/local/` out of generated surfaces because init must not
create the directory.

### F) Tasks Checklist

- [x] Update `src/codeheart_operating_kit/components.py` to include `.codeheart/local/` in `ensure_gitignore` output with a local machine layer heading.
- [x] Update `src/codeheart_operating_kit/components.py` so `write_default_state` writes `local_consumer_layer.local_machine_layer_path: .codeheart/local/` for new configs.
- [x] Confirm `scaffold_consumer_files` never creates `.codeheart/local/` during default onboarding.
- [x] Extend `schemas/kit-config.schema.json` with optional `local_machine_layer_path` using const `.codeheart/local/`.
- [x] Update `tests/fixtures/kit-config.yaml` to represent the new generated config shape.
- [x] Preserve `tests/fixtures/validator-valid/kit-config-without-purpose.yaml` as an old-shape compatibility fixture.
- [x] Add `tests/test_init.py` assertions for `.codeheart/local/` in `.gitignore`, generated config field presence, and absent `.codeheart/local/` directory.
- [x] Add `tests/test_sync_check.py` coverage for sync repairing a missing `.codeheart/local/` ignore line.
- [x] Add `tests/test_json_schemas.py` coverage for accepted new config, accepted old config, and rejected wrong local machine layer path.
- [x] Run `python3 -m pytest tests/test_init.py tests/test_sync_check.py tests/test_json_schemas.py`.

### G) Implementation Notes

Reuse the existing `ensure_gitignore` pattern rather than adding a second gitignore writer. The
heading can be separate from the local user layer heading so future agents can understand the
different boundary.

### H) Open Questions

No blockers. OQ-2 is resolved by keeping `.codeheart/local/` out of profile generated surfaces.

## EP-02 - Managed Doctrine And Readiness Runbook

### A) Epic ID, Title, And Outcome

EP-02 - Managed Doctrine And Readiness Runbook

Outcome: managed Operating Kit guidance distinguishes local user and local machine layers,
routes any missing local tooling blocker through tooling readiness, and gives the Python venv
lane a concrete default path.

### B) Scope

In scope:

- placement and managed-boundary docs for `.codeheart/local/`;
- planning workflow guidance for choosing the canonical owning repository before drafting
  discovery or implementation plans;
- local extension contract distinction between `.codeheart/user/` and `.codeheart/local/`;
- tooling-readiness trigger broadening;
- Python venv convention in the `python-runtime` lane;
- root managed block wording and installed fallback inventory;
- routing-standard and runbook-authoring references where the behavior is discoverable.

Out of scope:

- rewriting all planning workflow runbooks again;
- adding route cards for modules;
- adding exact OS install commands beyond current official-source references;
- changing AI Execution runbooks in this repository.

### C) Files Touched

```text
docs/repo/reference/placement-contract.md                           # modify
components/structure-governance/managed/reference/managed-content-boundaries.md # modify
components/agent-interface/managed/reference/local-extension-contract.md        # modify
components/agent-interface/managed/runbooks/handle-tooling-readiness.md         # modify
components/agent-interface/managed/reference/runbook-authoring-standard.md      # modify
components/agent-interface/managed/reference/operation-routing-and-dispatch.md  # modify
components/agent-interface/managed/runbooks/conduct-first-run-onboarding.md     # review
components/agent-interface/managed/kit-readme.md                               # modify
templates/agents/AGENTS.managed-block.md                                       # modify
components/planning-workflows/managed/runbooks/discovery-workflow.md           # modify
components/planning-workflows/managed/runbooks/draft-implementation-plan.md    # modify
components/planning-workflows/managed/reference/plan-register-format.md        # modify
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- placement docs list `.codeheart/local/` as ignored local machine/runtime state;
- managed boundaries define a local-machine/runtime ownership category;
- local extension contract keeps `.codeheart/user/` human-facing and `.codeheart/local/`
  generated/recreatable;
- root managed block says missing local tooling blockers route through tooling readiness for any
  applicable agent task;
- tooling-readiness runbook trigger includes repository tasks, module tasks, extension tasks, and
  agent-facing runbook blockers;
- `python-runtime` lane names `.codeheart/local/envs/python/` as the default venv path;
- venv setup text uses `.codeheart/local/envs/<purpose>/` only as an exception path;
- the runbook states no `sudo pip`, no `--break-system-packages`, no ad hoc global repo package
  installs, and no default first-run venv creation;
- first-run onboarding still says no Python venv is created during default onboarding;
- discovery and implementation drafting runbooks tell agents to place canonical planning documents
  in the repository that owns the work boundary when multiple repositories are involved;
- plan-register-format remains the coordination pointer source of truth and does not require
  moving canonical plans into `docs/repo/plans/`;
- routing example and authoring checks point local blockers to tooling readiness without turning
  live service preflight into local readiness.

### E) Dependencies And Critical-Path Notes

This epic depends on EP-01 decisions but can be edited in parallel with tests. Keep generic
Operating Kit doctrine public-safe and avoid source-specific local paths from the development
machine.

### F) Tasks Checklist

- [x] Update `docs/repo/reference/placement-contract.md` to add `.codeheart/local/` as ignored local machine/runtime state.
- [x] Update `components/structure-governance/managed/reference/managed-content-boundaries.md` to define the local machine layer and generated runtime boundary.
- [x] Update `components/agent-interface/managed/reference/local-extension-contract.md` to distinguish `.codeheart/user/` from `.codeheart/local/`.
- [x] Update `templates/agents/AGENTS.managed-block.md` to route missing local tooling blockers from any applicable agent task through tooling readiness.
- [x] Update `components/agent-interface/managed/kit-readme.md` to list `.codeheart/local/` under generated/local state.
- [x] Update `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` trigger text so it covers repository tasks, module tasks, extension tasks, and agent-facing runbook blockers.
- [x] Update `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` with the default Python venv path `.codeheart/local/envs/python/`.
- [x] Update `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` with exception venv path `.codeheart/local/envs/<purpose>/`.
- [x] Update `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` with global Python mutation guardrails.
- [x] Update `components/agent-interface/managed/reference/runbook-authoring-standard.md` to keep module package commands module-owned while routing generic Python venv placement through tooling readiness.
- [x] Update `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` local tooling example so the readiness route applies after owner and scope are selected.
- [x] Review `components/agent-interface/managed/runbooks/conduct-first-run-onboarding.md` and confirm its no-default-venv rule remains explicit.
- [x] Update `components/planning-workflows/managed/runbooks/discovery-workflow.md` to require a canonical owner-repository placement check before creating a discovery document in multi-repo work.
- [x] Update `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` to require a canonical owner-repository placement check before creating an implementation plan in multi-repo work.
- [x] Update `components/planning-workflows/managed/reference/plan-register-format.md` to make coordination-home register entries explicit pointers to canonical owner-repo plans.
- [x] Run a fresh low-context routing probe for a vague missing Python package blocker in an installed-module task and record the expected route in the execution log.

### G) Implementation Notes

Runbook-authoring coverage: `handle-tooling-readiness.md` is a hybrid runbook. It already has a
compact intention block, user-facing flow, execution path, stop conditions, and evidence. Preserve
that shape and add only the local runtime standard details needed for this capability.

Recipe-maturity coverage: the venv setup guidance should remain L1 structured runbook recipe
content. Do not promote it to a script, command wrapper, or API. Validation tier is fresh-agent
executability review plus focused non-live tests around generated repo behavior.

Routing-standard coverage: the changed root managed block, routing reference, and readiness
runbook are routing-bearing surfaces. A fresh low-context routing probe is required.

### H) Open Questions

No blockers. OQ-3 is resolved by preserving first-run onboarding unless review finds the existing
rule unclear.

## EP-03 - Packaged Mirrors, Docs Indexes, And Validation

### A) Epic ID, Title, And Outcome

EP-03 - Packaged Mirrors, Docs Indexes, And Validation

Outcome: all changed source docs are discoverable, packaged mirrors match source, and validation
proves the changed behavior.

### B) Scope

In scope:

- packaged resource mirrors for every changed managed file and template;
- packaging parity tests for newly covered managed files;
- docs indexes for the new implementation plan;
- focused validation commands and full test run.

Out of scope:

- publishing release artifacts;
- changing package version;
- changing public installer scripts;
- syncing named consumer repositories.

### C) Files Touched

```text
docs/README.md                                                       # modify
docs/repo/README.md                                                  # modify
docs/repo/plans/README.md                                            # modify
src/codeheart_operating_kit/resources/components/agent-interface/...  # modify mirrors
src/codeheart_operating_kit/resources/components/structure-governance/... # modify mirrors
src/codeheart_operating_kit/resources/components/planning-workflows/... # modify mirrors
src/codeheart_operating_kit/resources/templates/agents/AGENTS.managed-block.md # modify mirror
tests/test_packaging_resources.py                                    # modify
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- docs indexes link both the discovery and implementation plan;
- every changed managed source file has a matching packaged resource file;
- `tests/test_packaging_resources.py` covers changed managed files that were not previously in
  the parity list;
- generated `AGENTS.md` route visibility is covered so missing local tooling blockers remain
  discoverable in installed consumers;
- markdown header validation passes;
- public-core validation passes;
- JSON schema validation passes;
- focused pytest and full pytest pass;
- `git diff --check` passes.

### E) Dependencies And Critical-Path Notes

Run this epic after EP-01 and EP-02 edits so packaged mirrors are copied once from final source
docs. Do not hand-diverge mirrors from their source files.

### F) Tasks Checklist

- [x] Copy changed managed source files into matching `src/codeheart_operating_kit/resources/` paths.
- [x] Copy changed `templates/agents/AGENTS.managed-block.md` into the packaged template mirror.
- [x] Update `tests/test_packaging_resources.py` to cover changed managed files absent from the parity list.
- [x] Update `tests/test_onboard.py` or equivalent generated-root tests to protect installed
  `AGENTS.md` tooling-readiness route visibility.
- [x] Update `docs/README.md` with the implementation plan route.
- [x] Update `docs/repo/README.md` with the implementation plan route.
- [x] Update `docs/repo/plans/README.md` with the implementation plan route.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 scripts/validate-json-schemas.py`.
- [x] Run `python3 -m pytest tests/test_init.py tests/test_sync_check.py tests/test_json_schemas.py tests/test_packaging_resources.py`.
- [x] Run `python3 -m pytest`.
- [x] Run `git diff --check`.

### G) Implementation Notes

Use existing mirror conventions. Packaged resources should be byte-for-byte identical for managed
content already protected by parity tests.

### H) Open Questions

No blockers.

## EP-04 - Coordination, Release Requirement, And Foundry Handoff

### A) Epic ID, Title, And Outcome

EP-04 - Coordination, Release Requirement, And Foundry Handoff

Outcome: the plan/register state reflects implementation readiness, release impact is recorded,
and Foundry AI Execution has a clear downstream dependency on this Operating Kit standard.

### B) Scope

In scope:

- local plan register entry for this implementation plan;
- HQ coordination register and work-board pointer;
- release-note requirement and consumer-impact summary;
- downstream handoff note for AI Execution onboarding updates;
- execution-log placeholder during execution.

Out of scope:

- public release publication;
- tag creation;
- consumer sync proof;
- direct Foundry module patching.

### C) Files Touched

```text
docs/repo/plans/plan-register.md                                    # modify
docs/repo/plans/local-runtime-environment-standard/                  # execution log during execution
Codeheart-HQ:docs/repo/plans/plan-register.md                       # modify coordination pointer
Codeheart-HQ:docs/repo/plans/portfolio-work-board/portfolio-work-board.md # modify coordination pointer
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- Operating Kit register has an implementation-plan entry for `OK-PR-016`;
- discovery register entry relates to the implementation plan;
- HQ coordination register has a matching portfolio pointer for `CODEHEART-OPERATING-KIT-PR-016`;
- HQ work board shows the implementation plan as active and dependent on accepted discovery;
- release-note requirement states mixed backwards-compatible impact and no forced migration;
- downstream handoff states AI Execution install/onboard runbooks should adopt
  `.codeheart/local/envs/python/` after this standard is implemented.

### E) Dependencies And Critical-Path Notes

Register and HQ coordination updates happen during planning and after implementation completion.
Keep detailed execution evidence in the implementation execution log, not in the register.

### F) Tasks Checklist

- [x] Verify `OK-PR-016 - Local Runtime Environment Standard Implementation` exists in
  `docs/repo/plans/plan-register.md` and refresh lifecycle fields as implementation state changes.
- [x] Verify the child relation from `OK-PR-015` to `OK-PR-016` remains present.
- [x] Verify `CODEHEART-OPERATING-KIT-PR-016` exists in the HQ coordination register and refresh
  lifecycle fields as implementation state changes.
- [x] Verify `CODEHEART-OPERATING-KIT-PR-016` exists in the HQ portfolio work board under
  Operating Kit And Foundry System Model and refresh the board snapshot status.
- [x] Record release-note requirement as mixed backwards-compatible impact with no forced migration and no default local tool installation.
- [x] Record downstream Foundry handoff for AI Execution onboarding adoption of `.codeheart/local/envs/python/`.
- [x] During execution, create `local-runtime-environment-standard_execution_log.md` before marking any epic complete.
- [x] Run `rg -n "OK-PR-016|CODEHEART-OPERATING-KIT-PR-016"` against the Operating Kit and HQ
  plan registers and the HQ portfolio work board to verify coordination pointers.

### G) Implementation Notes

Use `not recorded` for session refs when no safe session ID is available. Do not scan private
transcripts solely to improve the register.

### H) Open Questions

No blockers. OQ-4 is resolved by keeping AI Execution changes outside this Operating Kit source
implementation.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

AI Execution onboarding adoption:

- deferred because the canonical local runtime standard should exist in Operating Kit source
  before module runbooks depend on it;
- trigger: EP-01 through EP-03 complete and validation passes;
- expected follow-up: update AI Execution install/onboard runbooks to create or reuse
  `.codeheart/local/envs/python/`, install the module package there, validate `foundry-ai`, and
  keep output routing module/runbook-owned.

Dependency lock strategy:

- deferred because the V1 requirement is a standard runtime boundary, not reproducible package
  locking;
- trigger: recurring dependency drift or module conflicts in the shared venv.

Script or CLI helper:

- deferred because L1 runbook guidance is sufficient for the first implementation;
- trigger: repeated agent mistakes, duplicated venv command blocks, or a need for stable
  structured output.

Durable readiness evidence:

- deferred because V1 should not create repo state, lockfile extensions, or reports for local
  machine readiness;
- trigger: repeated operations need non-secret readiness records across sessions.

Release publication:

- deferred because this plan prepares source implementation only;
- trigger: user approves release publication through the release runbook after implementation
  validation.

## 4.2 Future Considerations

- A future local runtime standard may add Node, browser, document/PDF, or vendor CLI runtime
  conventions under `.codeheart/local/` after real module usage proves common needs.
- Module dependency declarations may later become machine-readable, but the first implementation
  should keep module package details module-owned.
- A future Operating Kit release train should decide the release number and update public release
  surfaces once this plan is executed and reviewed.

# Revision Notes

- 2026-06-26: Initial draft created from the accepted local runtime environment standard
  discovery. Scope covers Operating Kit source behavior, managed doctrine, tests, packaged
  mirrors, and HQ coordination pointers while deferring AI Execution module adoption and release
  publication.
- 2026-06-26: Activated and completed the source implementation. Validation passed for focused
  init/sync/schema/onboarding/packaging tests, full pytest, markdown headers, public-core hygiene,
  JSON schemas, coordination pointers, packaged-resource parity, and whitespace checks. Public
  release publication, named consumer sync, and AI Execution module adoption remain deferred.
