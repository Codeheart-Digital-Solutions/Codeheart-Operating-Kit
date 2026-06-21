Last updated: 2026-06-21T15:32:08Z (UTC)
Created: 2026-06-21
Status: completed
Completed: 2026-06-21
Execution log: docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_execution_log.md

# Document Header

## Portfolio Coordination And Plan Register Implementation Plan

This implementation plan turns the reviewed discovery for portfolio coordination and plan-register
doctrine into concrete Codeheart Operating Kit changes.

The plan deliberately starts with runbooks, managed doctrine, kit-initialized consumer state files,
sync behavior, schema support, and tests. It does not add a portfolio CLI command in the first
implementation. That keeps the first release small enough to validate through real use while still
making the model available to new and existing consumers.

Target release preparation: `v0.1.5`, unless a maintainer reorders release numbers before
execution. Tagging, publishing, and PR merging remain out of scope for this implementation plan and
belong to the release runbook after the implementation is reviewed.

## Essential Context Reference Files

| Path | Reason |
| --- | --- |
| `AGENTS.md` | Public-core safety, read order, and maintainer instructions for Operating Kit changes. |
| `README.md` | Public repository purpose, release asset description, and bootstrap entry point. |
| `docs/README.md` | Public docs router that must expose the new plan once created. |
| `docs/repo/README.md` | Repository-governance router that must expose the new plan once created. |
| `docs/repo/plans/README.md` | Plan index that must link this implementation plan. |
| `docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_discovery_doc.md` | Source discovery and decision authority for this implementation plan. |
| `docs/repo/reference/placement-contract.md` | Ownership model for managed content, kit-initialized consumer state files, scaffolded consumer docs, templates, and generated surfaces. |
| `docs/repo/reference/consumer-impact-classification.md` | Consumer-impact classes for consumer state file, scaffold, sync, schema, managed-doc, and release-note decisions. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure for changes to kit source, docs, schemas, templates, validators, installers, or CLI behavior. |
| `components/planning-workflows/component.yaml` | Source manifest for planning-workflow managed files and new scaffold entries. |
| `components/planning-workflows/managed/` | Source managed planning workflow doctrine that will receive plan-register routes and hooks. |
| `components/agent-memory/component.yaml` | Source manifest for agent-memory managed docs and scaffolds. |
| `components/agent-memory/managed/` | Source managed agent-memory doctrine that must clarify the transitional role of `goal-register.md`. |
| `components/agent-memory/scaffolds/goal-register.md` | Existing consumer state baseline that must be preserved and clarified without forced migration. |
| `templates/agents/AGENTS.managed-block.md` | Managed `AGENTS.md` block template that may receive one lean conditional portfolio hook. |
| `templates/consumer-docs/repo/README.md` | Consumer `docs/repo/README.md` scaffold that must route the new plan-register files when appropriate. |
| `profiles/standard.yaml` | Source standard profile generated-surface declaration. |
| `schemas/kit-config.schema.json` | Consumer config schema that must accept optional portfolio coordination settings. |
| `src/codeheart_operating_kit/components.py` | Absent-file creation and generated-surface behavior for init/onboard and sync. |
| `src/codeheart_operating_kit/commands/sync.py` | Existing consumer update path that must create new absent state files without overwriting consumer-owned files. |
| `src/codeheart_operating_kit/resources/` | Packaged resource mirror that must stay in sync with source components, profiles, schemas, and templates. |
| `tests/test_init.py` | Init/onboard state-file creation coverage. |
| `tests/test_onboard.py` | Agent-guided onboarding and `onboard --yes` write behavior coverage. |
| `tests/test_sync_check.py` | Sync, check, drift, and update behavior coverage. |
| `tests/test_packaging_resources.py` | Packaged-resource fallback coverage. |
| `tests/test_json_schemas.py` | JSON schema validation coverage for config changes. |
| `release-notes.md` | Release-note surface for consumer-impact and adoption-path notes. |

## Table Of Contents

- [Section 1 - Foundation](#section-1---foundation)
- [Section 2 - Strategy](#section-2---strategy)
- [Section 3 - Execution Plan](#section-3---execution-plan)
- [Section 4 - Future Planning](#section-4---future-planning)
- [Revision Notes](#revision-notes)

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Implement the Operating Kit plan-register and optional portfolio-coordination model as an
additive, public-core-safe release.

The implementation must create a durable managed doctrine path for:

- `docs/repo/plans/plan-register.md` as the central kit-initialized consumer state file for formal
  plans, plan families, major workstreams, and portfolio-relevant planning records;
- `docs/repo/plans/coordination-sync-pending.md` as the kit-initialized consumer state file for
  coordination-home sync work that could not be applied because the coordination home is
  unavailable;
- `docs/agent-memory/goal-register.md` as transitional or informal pre-plan continuity memory,
  not the canonical formal planning register;
- optional portfolio coordination through explicit non-secret configuration rather than repository
  name guesses, GitHub scanning, or hidden assumptions;
- planning-workflow hooks that update the local register, optionally update a configured
  coordination-home register, and record pending sync when the coordination home is unavailable;
- lightweight session references for sessions that create or materially modify a plan when a
  session ID is available.

The implementation must make the new files available to:

- new consumers during init/onboarding;
- existing consumers through sync only when the state-file targets are absent.

The implementation must not:

- overwrite existing consumer state in `goal-register.md`, `plan-register.md`, or
  `coordination-sync-pending.md`;
- force a migration of existing goal-register contents;
- add a portfolio CLI command in this release;
- ask portfolio-coordination questions during normal first-run onboarding;
- silently scan or mutate sibling repositories, GitHub repositories, or coordination homes;
- hardcode Codeheart HQ, private repository names, local absolute paths, or tenant-specific details
  into public managed doctrine.

## 1.2 Project And Problem Context

The Operating Kit currently has a thin `docs/agent-memory/goal-register.md` scaffold and managed
agent-memory instructions that say memory state belongs to consumers. That is correct for agent
memory, but it is too vague for formal planning across multiple repositories.

Codeheart now has several repositories and workstreams that need a higher-level planning index:
Operating Kit, AWS Platform, Foundry, HQ, and future consumer repositories. Discovery plans and
implementation plans exist in each repository, while some work spans several repositories. Without
a reusable plan-register model, agents can lose the big picture, invent local indexes, overload
agent-memory files, or depend on chat memory.

The discovery chose a small model:

- local repositories keep their own `docs/repo/plans/plan-register.md`;
- a coordination home can use the same register model for a portfolio-level view;
- coordination is optional and explicit;
- `goal-register.md` remains available for informal or pre-plan continuity during transition;
- the first implementation should be runbooks, scaffolds, managed docs, a lean `AGENTS.md` hook,
  schema support, and tests, not a new CLI command.

The resulting release is mostly instruction and state-surface baseline work, with one important
behavior change: sync must be able to create newly introduced absent kit-initialized consumer state
files for existing consumers without touching any existing consumer-owned content.

Terminology for this plan:

- Operating Kit owns the file contract: location, baseline template, format, runbooks,
  creation/repair behavior, and validation expectations.
- The consumer repository owns the file content after creation: plan entries, local lifecycle
  snapshots, relations, session refs, and pending-sync items.
- Sync may create or recreate the baseline file when it is absent.
- Sync must never overwrite, normalize, or migrate the file when it already exists.

## 1.3 Current State Analysis

### Current Managed Planning Workflows

The planning-workflows component currently owns managed runbooks for discovery, implementation-plan
drafting, implementation execution, and planning-document review. Those runbooks do not yet define
a reusable plan-register format or instruct agents when plan lifecycle, relationships, or session
references should be recorded.

Current gap: a plan can be created, materially updated, completed, superseded, or linked to another
plan without any durable index update.

### Current Agent Memory

The agent-memory component currently scaffolds:

- `docs/agent-memory/README.md`;
- `docs/agent-memory/goal-register.md`;
- `docs/agent-memory/session-ledger.md`;
- `docs/agent-memory/untriaged-sessions.md`.

Those files are consumer state files created from Operating Kit baselines when absent. Their
consumer-owned contents must be preserved. The gap is semantic: `goal-register.md` is too broad and
too thin to be the formal plan/workstream register. It must be clarified as transitional or
informal continuity memory.

### Current Consumer Scaffolding And Sync

Init/onboarding creates managed kit content and absent consumer state baselines. Sync refreshes
managed content, but it does not currently add newly introduced state files to an existing consumer
when they are absent.

Current gap: if the kit adds `docs/repo/plans/plan-register.md` as a new kit-initialized consumer
state file, new installs can receive it, but existing consumers will not get it through normal sync
unless sync is extended.

### Current Configuration Schema

`.codeheart/kit.config.yaml` is schema-validated and currently has no optional portfolio block.
The discovery selected explicit non-secret configuration over inference or scanning.

Current gap: agents and future CLI commands need a deterministic place to read portfolio role and
coordination-home paths. The first implementation only needs schema support and doctrine; it does
not need an interactive command.

### Current Agent Interface

The managed `AGENTS.md` block is intentionally short and routing-focused. It does not mention
portfolio coordination.

Current gap: a fresh agent needs a lean reminder that configured portfolio coordination may require
local and coordination-home register maintenance for material planning lifecycle changes. The block
must stay generic and must route details to managed runbooks.

### Current Public-Core Risk

The coordination-home concept is motivated by Codeheart HQ, but the public kit must remain generic.
The implementation must avoid private Codeheart examples in managed doctrine, schemas, templates,
release notes, and test fixtures unless sanitized as generic examples.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement the change in six epics:

1. Add managed plan-register doctrine and a maintenance runbook.
2. Add kit-initialized consumer state files and teach sync to create absent baselines safely.
3. Add targeted hooks to the managed planning workflows.
4. Add the lean agent-interface hook and optional portfolio config schema.
5. Clarify agent-memory transition and consumer documentation routes.
6. Update packaging, release notes, indexes, and validation.

The target source tree additions and updates are:

```text
Codeheart-Operating-Kit/
  components/
    planning-workflows/
      component.yaml
      managed/
        README.md
        reference/
          plan-register-format.md                    # new managed reference
          planning-document-lifecycle.md             # update lifecycle routes
        runbooks/
          discovery-workflow.md                      # update targeted hook
          draft-implementation-plan.md               # update targeted hook
          execute-implementation-plan.md             # update targeted hook
          maintain-plan-register.md                  # new managed runbook
      scaffolds/
        plan-register.md                             # new kit-initialized consumer state baseline
        coordination-sync-pending.md                 # new kit-initialized consumer state baseline
    agent-memory/
      managed/
        README.md                                    # clarify role split
        reference/
          entry-format.md                            # clarify goal-register boundaries
        runbooks/
          session-ledger-maintenance.md              # clarify plan-register versus memory boundary
      scaffolds/
        goal-register.md                             # clarify transitional/informal purpose
    agent-interface/
      managed/
        README.md                                    # route conditional portfolio hook
        reference/
          root-agents-md-contract.md                 # update immediate-contract expectations
  templates/
    agents/
      AGENTS.managed-block.md                        # lean conditional hook
    consumer-docs/
      repo/
        README.md                                    # route plan register and pending sync
  docs/
    README.md                                        # route this implementation plan
    repo/
      README.md                                      # route this implementation plan
      reference/
        placement-contract.md                        # classify new kit-initialized consumer state files
      plans/
        README.md                                    # route this implementation plan
        portfolio-coordination-plan-register/
          portfolio-coordination-plan-register_discovery_doc.md
          portfolio-coordination-plan-register_implementation_doc.md
  profiles/
    standard.yaml                                    # add generated surfaces for new scaffolds
  schemas/
    kit-config.schema.json                           # optional portfolio block
  src/
    codeheart_operating_kit/
      components.py                                  # return/merge scaffold records for sync
      commands/
        sync.py                                      # create absent new scaffolds, preserve existing
      resources/
        manifest.yaml                                # include new packaged resources
        components/                                  # mirror component source files
        profiles/
          standard.yaml                              # mirror profile changes
        templates/
          agents/
            AGENTS.managed-block.md                  # packaged managed-block template mirror
          consumer-docs/
            repo/
              README.md                              # packaged consumer repo README scaffold mirror
  tests/
    test_init.py                                     # new install scaffold expectations
    test_onboard.py                                  # onboarding scaffold and default-config expectations
    test_sync_check.py                               # existing consumer sync/preservation tests
    test_packaging_resources.py                      # packaged-resource coverage
    test_json_schemas.py                             # optional portfolio config coverage
    test_public_core.py                              # existing public-core guard coverage
  release-notes.md                                   # additive scaffold/doctrine release note
```

Generated or packaged mirrors must match source files. Source component, profile, and template
files must not diverge from their packaged counterparts under
`src/codeheart_operating_kit/resources/`.

## 2.2 Open Questions And Assumptions Requiring Clarification

No blocking questions remain from the discovery. The following implementation defaults are
accepted for this plan:

| ID | Question | Default For This Plan |
| --- | --- | --- |
| `IQ-1` | Should the first release add a portfolio CLI command? | No. Add runbooks, scaffolds, schema support, hooks, and tests first. |
| `IQ-2` | Should normal onboarding ask portfolio-coordination questions? | No. Keep portfolio coordination out of normal first-run onboarding. |
| `IQ-3` | Should existing consumers receive the new planning files? | Yes, through sync when each target file is absent, while preserving any existing content. |
| `IQ-4` | Should sync migrate or rewrite `goal-register.md`? | No. Existing consumer-owned files are preserved exactly. |
| `IQ-5` | Should the managed `AGENTS.md` block mention portfolio coordination? | Yes, but only as one lean generic conditional hook that routes to runbooks. |
| `IQ-6` | Should portfolio config be generated by default? | No. Schema allows it; absent config means no configured portfolio coordination. |
| `IQ-7` | Should plan-register validation be implemented now? | No dedicated validator in this plan. Tests and docs are sufficient for first adoption. |
| `IQ-8` | Should release notes call this a migration? | No. It is additive scaffold and doctrine behavior with no forced migration. |

Assumptions to validate during implementation:

- `sync` can call existing absent-file creation logic or a small shared helper without changing
  managed file semantics or overwriting consumer-owned state.
- Lock `generated_surfaces` can merge newly created scaffold records without rewriting unrelated
  lock state.
- Adding optional `portfolio` to the config schema is backwards-compatible because existing config
  files remain valid.
- Managed planning docs can add targeted register hooks without becoming too long or duplicative.
- Tests can use generic repository names and paths to preserve public-core hygiene.

## 2.3 Architectural Decisions With Reasoning

### `AD-1` Plan Register Belongs Under `docs/repo/plans/`

Decision: initialize and document `docs/repo/plans/plan-register.md` as the central formal planning
register.

Reasoning: formal plans live under repository planning docs, not agent memory. This placement makes
the register part of repository documentation governance while keeping it consumer-owned.

### `AD-2` Goal Register Remains But Is Not The Formal Register

Decision: keep `docs/agent-memory/goal-register.md` as a scaffold and clarify that it is for
informal or pre-plan continuity during transition.

Reasoning: existing consumers may already have useful goal-register content. Rewriting or removing
it would be risky and unnecessary. The new register provides a clearer home for formal plans
without breaking existing memory state.

### `AD-3` Register Entries Use Repeated Markdown Sections

Decision: the plan-register scaffold and reference use one repeated section per plan/workstream
entry, not a wide table.

Reasoning: relations, session refs, and coordination notes become unreadable in a wide table.
Repeated sections are easier for agents to update without damaging adjacent entries.

### `AD-4` Session References Are Lightweight Recovery Handles

Decision: record creating and material-update session IDs when available, but do not require
session summaries.

Reasoning: canonical plans should summarize state. Session IDs help recovery without turning the
register into a transcript index or session ledger.

### `AD-5` Sync Creates Absent New Scaffolds For Existing Consumers

Decision: extend sync so new kit-initialized consumer state definitions can create absent files in
existing installs, while preserving every existing file.

Reasoning: limiting scaffolds to new installs leaves existing consumers without the new planning
surface. Overwriting consumer files would violate the placement contract. Absent-file sync is the
safe middle path.

### `AD-6` Portfolio Config Is Optional And Explicit

Decision: add an optional presence-based `portfolio` object to the config schema, but do not create
it in default init/onboard output.

Reasoning: explicit config avoids brittle inference and silent scans. Leaving it absent by default
keeps simple consumers simple and preserves the current onboarding shape.

### `AD-7` First Release Is Runbook-First, Not CLI-First

Decision: defer `portfolio configure` or equivalent CLI work.

Reasoning: the operating model should be used through docs and scaffolds before the command
contract is frozen. The first release still gives agents enough instructions to act correctly.

### `AD-8` Coordination Home Unavailable Means Pending Sync, Not Failure

Decision: when coordination is configured but the coordination home is unavailable, agents record
the needed update in `docs/repo/plans/coordination-sync-pending.md` and continue local work.

Reasoning: local repository work must remain possible offline or without sibling checkouts. Silent
skips make coordination unreliable; pending sync keeps the gap visible.

### `AD-9` Managed `AGENTS.md` Stays Lean

Decision: add only one generic conditional hook to the managed block.

Reasoning: `AGENTS.md` is a bootstrap contract, not the full procedure. The detailed rules belong
in managed planning runbooks and plan-register references.

### `AD-10` Public-Core Hygiene Controls Examples

Decision: use generic examples such as `Companyname-Automation`, `Yourname-Automation`, or
sanitized repository names only when examples are necessary.

Reasoning: the public kit must not encode private Codeheart repository topology or tenant details.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Title | Outcome | Size | Depends On |
| --- | --- | --- | --- | --- |
| `EP1` | Managed Plan-Register Doctrine | Full managed format reference and maintenance runbook exist and are packaged. | M | None |
| `EP2` | Consumer Scaffolds And Safe Sync | New installs and existing sync create absent plan-register files without overwriting consumer state. | L | `EP1` |
| `EP3` | Planning Workflow Lifecycle Hooks | Discovery, implementation planning, and execution workflows route material changes to the register. | M | `EP1` |
| `EP4` | Agent Interface And Portfolio Config Schema | Managed `AGENTS.md` gets one lean hook and config schema accepts optional portfolio settings. | M | `EP1` |
| `EP5` | Agent Memory And Consumer Documentation Transition | Goal-register boundaries and consumer repo routes are clear and non-migratory. | M | `EP1`, `EP2` |
| `EP6` | Packaging, Release Notes, And Validation | Packaged resources, indexes, release notes, and tests prove the release is additive and safe. | M | `EP1`-`EP5` |

## EP1 - Managed Plan-Register Doctrine

### A) Epic ID, Title, And Outcome

`EP1` - Managed Plan-Register Doctrine

Outcome: the Operating Kit has a full managed plan-register reference and maintenance runbook that
agents can follow without needing the discovery document or private Codeheart examples.

### B) Scope

In scope:

- Create `components/planning-workflows/managed/reference/plan-register-format.md`.
- Create `components/planning-workflows/managed/runbooks/maintain-plan-register.md`.
- Update `components/planning-workflows/managed/README.md` to route both files.
- Update `components/planning-workflows/component.yaml` to include the new managed files.
- Mirror all new and changed files under
  `src/codeheart_operating_kit/resources/components/planning-workflows/`.
- Include full behavior, not stub text.

Out of scope:

- Adding a portfolio CLI command.
- Adding a plan-register validator.
- Writing Codeheart-HQ-specific doctrine.
- Changing consumer scaffolding behavior; that is `EP2`.

### C) Files Touched

Expected files:

- `components/planning-workflows/component.yaml`
- `components/planning-workflows/managed/README.md`
- `components/planning-workflows/managed/reference/plan-register-format.md`
- `components/planning-workflows/managed/runbooks/maintain-plan-register.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/README.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md`

### D) Acceptance Criteria And Size

Size: `M`

Acceptance criteria:

- `plan-register-format.md` defines the register purpose, source-of-truth rule, entry fields,
  repeated-section Markdown shape, lifecycle values, relation vocabulary, session-ref shape,
  coordination note, and anti-patterns.
- `maintain-plan-register.md` gives a concrete procedure for local register updates, coordination
  home updates, pending-sync fallback, session-ref handling, and canonical-doc conflict handling.
- `maintain-plan-register.md` defines a concrete pending-sync entry shape covering source
  repository, target coordination register, affected plan entry, intended change, reason, date,
  session ref when available, status, and clear/complete handling.
- The runbook states that typos, formatting-only edits, timestamp refreshes, and mechanical
  checklist progress do not require register updates.
- The runbook states that creating or materially modifying sessions should be recorded when a
  session ID is available and should not block work when no ID is available.
- The runbook states that missing coordination-home access records pending sync and does not fail
  the local planning task.
- The managed docs contain no Codeheart-private repository names, local absolute paths, tenant
  details, or HQ-only assumptions.
- Source and packaged-resource mirrors match.

### E) Dependencies And Critical-Path Notes

This epic is first because later runbook hooks and scaffolds need an authoritative managed
destination. Implement this before changing workflow instructions so every hook can route to a real
document.

### F) Tasks Checklist

- [x] Read the discovery decisions and current planning-workflows managed files before drafting.
- [x] Add `components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Add `components/planning-workflows/managed/runbooks/maintain-plan-register.md`.
- [x] Define the pending-sync entry format inside `maintain-plan-register.md`.
- [x] Update `components/planning-workflows/managed/README.md` with routes to the new reference and runbook.
- [x] Update `components/planning-workflows/component.yaml` with the new managed files.
- [x] Mirror new and changed planning-workflow files under `src/codeheart_operating_kit/resources/components/planning-workflows/`.
- [x] Verify the reference and runbook include complete instructions rather than placeholder summaries.
- [x] Verify examples are generic and public-core safe.

### G) Implementation Notes

The reference should be declarative and stable. The runbook should be procedural. Avoid duplicating
large sections of one into the other. The runbook can link to the reference for entry shape and
field semantics.

Use exact lifecycle terms from the discovery: `draft`, `active`, `completed`, `superseded`, and
`archived`.

Use relation terms from the discovery: `parent`, `child`, `supersedes`, `superseded-by`,
`depends-on`, `blocks`, and `related`.

### H) Open Questions

None.

## EP2 - Consumer State Files And Safe Sync

### A) Epic ID, Title, And Outcome

`EP2` - Consumer State Files And Safe Sync

Outcome: new consumers receive plan-register and pending-sync baseline files, and existing
consumers receive those files through sync when absent. Existing consumer-owned file content is
preserved.

### B) Scope

In scope:

- Add source baseline files for:
  - `docs/repo/plans/plan-register.md`;
  - `docs/repo/plans/coordination-sync-pending.md`.
- Register the files in `components/planning-workflows/component.yaml` as kit-initialized consumer
  state files or the nearest supported ownership metadata updated by this plan.
- Update planning-workflow component metadata and ownership doctrine to include
  kit-initialized consumer state file behavior and the relevant backwards-compatible consumer
  impact.
- Mirror the baseline files and manifest metadata under packaged resources.
- Update `profiles/standard.yaml` and packaged profile generated surfaces.
- Extend sync to create newly defined absent scaffold files.
- Ensure sync merges created scaffold records into lock `generated_surfaces` without dropping or
  rewriting unrelated lock data.
- Add init and sync tests for scaffold creation and preservation.
- Add onboarding tests for scaffold creation and no portfolio-coordination default.

Out of scope:

- Migrating any existing `goal-register.md` content.
- Rewriting any existing consumer `plan-register.md` or `coordination-sync-pending.md` content.
- Adding a dedicated register validator.
- Adding a CLI command for portfolio configuration.

### C) Files Touched

Expected files:

- `components/planning-workflows/component.yaml`
- `components/planning-workflows/scaffolds/plan-register.md`
- `components/planning-workflows/scaffolds/coordination-sync-pending.md`
- `profiles/standard.yaml`
- `src/codeheart_operating_kit/components.py`
- `src/codeheart_operating_kit/commands/sync.py`
- `src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml`
- `src/codeheart_operating_kit/resources/components/planning-workflows/scaffolds/plan-register.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/scaffolds/coordination-sync-pending.md`
- `src/codeheart_operating_kit/resources/profiles/standard.yaml`
- `tests/test_init.py`
- `tests/test_onboard.py`
- `tests/test_sync_check.py`

### D) Acceptance Criteria And Size

Size: `L`

Acceptance criteria:

- `codeheart-operating-kit init` and onboarding create `docs/repo/plans/plan-register.md` when the
  file is absent.
- `codeheart-operating-kit init` and onboarding create `docs/repo/plans/coordination-sync-pending.md`
  when the file is absent.
- Sync creates both files for an existing consumer only when each file is absent.
- Sync preserves existing consumer content in:
  - `docs/agent-memory/goal-register.md`;
  - `docs/repo/plans/plan-register.md`;
  - `docs/repo/plans/coordination-sync-pending.md`.
- Sync does not mark preserved consumer-owned files as drift.
- Lock `generated_surfaces` includes newly created scaffold records after sync and preserves
  existing records.
- Source and packaged planning-workflow component metadata identifies the new files as
  kit-initialized consumer state files, or as the nearest supported ownership mode plus explicit
  presence-only/no-overwrite semantics, and records the relevant backwards-compatible consumer
  impact.
- The plan-register scaffold uses repeated Markdown sections, not one wide table.
- The pending-sync scaffold clearly states that local planning work can continue while coordination
  sync is pending.
- `onboard --yes` creates both planning scaffolds when absent.
- Onboarding output does not ask normal first-run users to configure portfolio coordination.
- Default init/onboard config output contains no `portfolio` block.
- No normal onboarding prompt asks about portfolio coordination.

### E) Dependencies And Critical-Path Notes

Depends on `EP1` so scaffolds can reference the managed format and runbook. This is the main code
epic because it changes sync behavior for existing consumers.

### F) Tasks Checklist

- [x] Add `components/planning-workflows/scaffolds/plan-register.md` with a repeated-section starter.
- [x] Add `components/planning-workflows/scaffolds/coordination-sync-pending.md`.
- [x] Update `components/planning-workflows/component.yaml` to declare kit-initialized consumer state behavior and absent-file install/repair behavior for both files.
- [x] Update source and packaged planning-workflow component metadata with the chosen ownership/presence semantics and backwards-compatible consumer-impact classification.
- [x] Mirror scaffold files and component metadata under `src/codeheart_operating_kit/resources/components/planning-workflows/`.
- [x] Update `profiles/standard.yaml` with the new generated surfaces.
- [x] Mirror the standard profile under `src/codeheart_operating_kit/resources/profiles/standard.yaml`.
- [x] Refactor absent-file creation as needed so sync can create missing state baselines safely.
- [x] Update `src/codeheart_operating_kit/commands/sync.py` to create newly introduced absent scaffolds.
- [x] Update lock refresh behavior so newly created scaffold records are merged into `generated_surfaces`.
- [x] Add init tests that assert both new planning files are created on new installs.
- [x] Add onboarding tests that assert `onboard --yes` creates both new planning files when absent.
- [x] Add onboarding tests that assert normal onboarding output contains no portfolio-coordination setup prompt.
- [x] Add onboarding or init tests that assert default config contains no `portfolio` block.
- [x] Add sync tests that assert missing planning files are created for existing installs.
- [x] Add sync tests that assert existing `goal-register.md`, `plan-register.md`, and `coordination-sync-pending.md` contents are preserved byte-for-byte.
- [x] Add sync tests that assert generated-surface lock records are preserved and extended rather than replaced.

### G) Implementation Notes

The existing `scaffold_consumer_files` helper already has absent-file behavior for init/onboard.
Prefer reusing or carefully extracting that logic instead of creating a second file-creation
engine. Keep the content rule explicit: create when absent, preserve exactly when present.

When updating lock `generated_surfaces`, merge by path or equivalent stable identity. Do not
replace the whole list with only this sync run's newly created files.

If sync reports changed files to the user, distinguish managed refreshes from newly created
kit-initialized consumer state files.

### H) Open Questions

None.

## EP3 - Planning Workflow Lifecycle Hooks

### A) Epic ID, Title, And Outcome

`EP3` - Planning Workflow Lifecycle Hooks

Outcome: managed discovery, implementation planning, and implementation execution workflows tell
agents when to maintain plan registers and how to avoid noisy updates.

### B) Scope

In scope:

- Update discovery workflow instructions for plan creation, material decision changes, relation
  changes, and discovery completion/supersession.
- Update implementation-plan drafting instructions for plan creation, scope changes, status
  changes, relation links, and session refs.
- Update implementation execution instructions for activation, completion, supersession, archived
  state, material implementation-path changes, and execution handoff.
- Each workflow hook must route to the full maintenance sequence: update the local plan register;
  update the configured coordination-home register when configured and available; write
  `docs/repo/plans/coordination-sync-pending.md` when coordination is configured but unavailable.
- Update planning lifecycle reference to mention plan-register synchronization.
- Mirror all changed managed files under packaged resources.

Out of scope:

- Requiring register updates for every typo, timestamp refresh, formatting-only change, or
  mechanical checklist tick.
- Duplicating implementation-plan or execution-log details inside the plan register.
- Adding non-planning task hooks.
- Editing `review-planning-document.md` for register updates. Review-only work remains read-only;
  stale register observations are reported as review findings, not applied as automatic mutations.

### C) Files Touched

Expected files:

- `components/planning-workflows/managed/runbooks/discovery-workflow.md`
- `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
- `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
- `components/planning-workflows/managed/reference/planning-document-lifecycle.md`
- Matching files under `src/codeheart_operating_kit/resources/components/planning-workflows/managed/`

### D) Acceptance Criteria And Size

Size: `M`

Acceptance criteria:

- Discovery workflow instructs agents to maintain the local plan register when discovery docs are
  created or materially updated.
- Discovery workflow instructs goal-style discovery agents to record material decision and
  relationship changes without bloating the register.
- Implementation-planning workflow instructs agents to register new implementation plans and
  material plan changes.
- Implementation execution workflow instructs agents to update the register when a plan becomes
  active, completed, superseded, archived, or materially changes path.
- Discovery, implementation-planning, and implementation-execution hooks all require the same
  conditional sequence: local register update first, configured coordination-home update second,
  pending-sync fallback when the coordination home is unavailable.
- All hooks route to `maintain-plan-register.md` and `plan-register-format.md` instead of
  embedding full procedures in every runbook.
- Hooks explicitly exclude typo-only, formatting-only, timestamp-only, and mechanical checklist
  progress changes.
- Review-only workflows do not update registers. When a user authorizes edits after a review, those
  edits route through the discovery or implementation planning workflow and use its register hook
  only if the canonical plan is materially changed.
- Source and packaged-resource mirrors match.

### E) Dependencies And Critical-Path Notes

Depends on `EP1`. This epic is mostly managed-doc work, but it is important because the
plan-register model will fail if agents only see the scaffold and not the lifecycle triggers.

### F) Tasks Checklist

- [x] Update `discovery-workflow.md` with targeted plan-register lifecycle and session-ref hooks.
- [x] Update `draft-implementation-plan.md` with targeted plan-register lifecycle and session-ref hooks.
- [x] Update `execute-implementation-plan.md` with targeted plan-register lifecycle and session-ref hooks.
- [x] Ensure each planning workflow hook routes to local register, coordination-home update, and pending-sync fallback behavior through `maintain-plan-register.md`.
- [x] Update `planning-document-lifecycle.md` with the register relationship.
- [x] Mirror every changed planning-workflow managed file under packaged resources.
- [x] Confirm `review-planning-document.md` remains outside the mutation path for register updates.
- [x] Review all hooks for noise control and source-of-truth clarity.

### G) Implementation Notes

Do not make every workflow repeat the full register entry schema. Each workflow should say when to
invoke the maintenance runbook and what kind of lifecycle change is material.

The implementation execution hook should distinguish canonical execution state in the
implementation plan or execution log from register metadata. The register is an index snapshot.

Do not turn planning-document review into a side-effecting workflow. A review may report that a
plan register appears stale, but the register update itself belongs to the later user-authorized
discovery or implementation planning edit that changes the canonical plan.

### H) Open Questions

None.

## EP4 - Agent Interface And Portfolio Config Schema

### A) Epic ID, Title, And Outcome

`EP4` - Agent Interface And Portfolio Config Schema

Outcome: the managed bootstrap interface includes one generic conditional portfolio hook, and
consumer config schema can represent optional coordination-home settings without changing default
onboarding output.

### B) Scope

In scope:

- Add one lean conditional portfolio-coordination hook to
  `templates/agents/AGENTS.managed-block.md`.
- Update agent-interface managed docs to describe the hook and route details to planning-workflow
  runbooks.
- Update sync behavior so existing consumers receive managed `AGENTS.md` block template changes
  while preserving repository-owned content outside the managed block.
- Add optional presence-based `portfolio` config schema support.
- Add schema tests for omitted portfolio config and valid portfolio config.
- Confirm default init/onboard config does not add a portfolio block.
- Mirror changed agent-interface files and schema resources as needed.

Out of scope:

- Adding a portfolio CLI command.
- Asking about portfolio coordination during normal onboarding.
- Automatically setting portfolio config based on repository names.
- Auto-discovering or scanning coordination homes.

### C) Files Touched

Expected files:

- `templates/agents/AGENTS.managed-block.md`
- `components/agent-interface/managed/README.md`
- `components/agent-interface/managed/reference/root-agents-md-contract.md`
- `schemas/kit-config.schema.json`
- `src/codeheart_operating_kit/components.py`
- `src/codeheart_operating_kit/commands/sync.py`
- `src/codeheart_operating_kit/resources/templates/agents/AGENTS.managed-block.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/README.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/root-agents-md-contract.md`
- `tests/test_json_schemas.py`
- `tests/test_init.py`
- `tests/test_onboard.py`
- `tests/test_sync_check.py`

### D) Acceptance Criteria And Size

Size: `M`

Acceptance criteria:

- Managed `AGENTS.md` contains one lean, generic conditional hook for configured portfolio
  coordination.
- The hook does not hardcode Codeheart HQ, local absolute paths, or private repository names.
- The hook directs agents to managed planning-workflow runbooks for details.
- Sync refreshes the managed block of existing consumer `AGENTS.md` files when markers are present
  and preserves repository-owned content outside the managed block.
- `kit-config.schema.json` accepts config files with no `portfolio` block.
- `kit-config.schema.json` accepts a valid optional portfolio block with explicit role and
  coordination-home fields.
- `kit-config.schema.json` pins portfolio role values to `coordination-home` and `member`.
- When a `portfolio` block is present, `role` is required and no unknown portfolio fields are
  accepted.
- `role: member` requires `member_repository_id`, `coordination_home_path`, and
  `coordination_home_register_path`.
- `role: coordination-home` requires `coordination_home_register_path`.
- The schema rejects `enabled`, `standalone`, `member_register_path`, and `pending_sync_path`.
- `kit-config.schema.json` rejects invalid portfolio role values.
- Default init/onboard config output remains free of portfolio settings.
- Tests prove schema compatibility, instance-level config behavior, sync-managed-block behavior,
  and default config behavior.

### E) Dependencies And Critical-Path Notes

Depends on `EP1` for the runbook route. Can be implemented in parallel with `EP3` after the runbook
path is final.

### F) Tasks Checklist

- [x] Draft the managed `AGENTS.md` conditional portfolio hook in the source template.
- [x] Update agent-interface managed docs to explain the hook boundary.
- [x] Mirror source template and managed docs under packaged resources.
- [x] Update sync behavior to refresh only the managed `AGENTS.md` block when markers are present.
- [x] Add sync tests proving managed-block refresh and preservation of repository-owned content outside the block.
- [x] Extend `schemas/kit-config.schema.json` with optional presence-based `portfolio`.
- [x] Add fixture-based config instance tests for omitted portfolio config.
- [x] Add schema tests for a valid member portfolio config.
- [x] Add schema tests for a valid coordination-home portfolio config.
- [x] Add schema tests for invalid role rejection.
- [x] Add schema tests rejecting `portfolio: {}`.
- [x] Add schema tests rejecting `role: standalone`.
- [x] Add schema tests rejecting `enabled`.
- [x] Add schema tests rejecting `member_register_path`.
- [x] Add schema tests rejecting `pending_sync_path`.
- [x] Add schema tests rejecting member configs missing `member_repository_id`, `coordination_home_path`, or `coordination_home_register_path`.
- [x] Add schema tests rejecting coordination-home configs missing `coordination_home_register_path`.
- [x] Confirm default init/onboard test fixtures or assertions still show no default portfolio block.

### G) Implementation Notes

Recommended schema shape:

```yaml
portfolio:
  role: member
  coordination_home_path: ../Coordination-Home
  coordination_home_register_path: docs/repo/plans/plan-register.md
  member_repository_id: Example-Repository
```

The exact requiredness within `portfolio` can be conservative:

- no `portfolio` block means no configured portfolio coordination.
- `role` should be required when `portfolio` exists.
- `role` must be one of `coordination-home` or `member`.
- `member` requires `member_repository_id`, `coordination_home_path`, and
  `coordination_home_register_path`.
- `coordination-home` requires `coordination_home_register_path`.
- the local member register path is inferred from `local_consumer_layer.repo_docs_path` plus
  `plans/plan-register.md`.
- the local pending-sync path is inferred from `local_consumer_layer.repo_docs_path` plus
  `plans/coordination-sync-pending.md`.
- the schema should use strings and enums only; no secrets.
- the schema should reject `enabled`, `standalone`, `member_register_path`, `pending_sync_path`,
  and incomplete role-specific configs.
- tests should validate real config fixtures against the touched `portfolio` schema behavior
  without adding a runtime dependency.

### H) Open Questions

None.

## EP5 - Agent Memory And Consumer Documentation Transition

### A) Epic ID, Title, And Outcome

`EP5` - Agent Memory And Consumer Documentation Transition

Outcome: consumer-facing docs clearly explain that the formal plan register lives under
`docs/repo/plans/`, while `docs/agent-memory/goal-register.md` remains available for informal or
pre-plan continuity.

### B) Scope

In scope:

- Update managed agent-memory docs to clarify the split between formal planning registration and
  agent memory.
- Update `session-ledger-maintenance.md` so it no longer routes formal plan lifecycle, formal
  workstream status, blockers, canonical plan docs, or plan relationships into `goal-register.md`.
- Update the `goal-register.md` scaffold text to avoid implying it is the canonical formal plan
  register.
- Update consumer `docs/repo/README.md` scaffold to route plan-register and pending-sync files.
- Update placement contract with the new kit-initialized consumer state files and ownership
  expectations.
- Mirror changed files under packaged resources.

Out of scope:

- Moving existing goal-register content.
- Removing goal-register scaffolding.
- Creating a separate HQ-specific register mechanism.
- Adding detailed project-management or next-action fields.

### C) Files Touched

Expected files:

- `components/agent-memory/managed/README.md`
- `components/agent-memory/managed/reference/entry-format.md`
- `components/agent-memory/managed/runbooks/session-ledger-maintenance.md`
- `components/agent-memory/scaffolds/goal-register.md`
- `templates/consumer-docs/repo/README.md`
- `docs/repo/reference/placement-contract.md`
- Matching packaged-resource files under `src/codeheart_operating_kit/resources/components/`
- `src/codeheart_operating_kit/resources/templates/consumer-docs/repo/README.md`
- `tests/test_packaging_resources.py`

### D) Acceptance Criteria And Size

Size: `M`

Acceptance criteria:

- Managed agent-memory docs state that formal discovery and implementation plans are registered in
  `docs/repo/plans/plan-register.md`.
- Managed agent-memory docs state that `goal-register.md` remains for informal or pre-plan
  continuity and transition.
- `session-ledger-maintenance.md` states that formal plan lifecycle, formal workstream status,
  blockers, canonical plan docs, and plan relationships belong in `docs/repo/plans/plan-register.md`
  when a formal plan exists or is created.
- `session-ledger-maintenance.md` keeps `goal-register.md` available for informal, pre-plan, or
  transitional continuity that has not yet become a formal discovery or implementation plan.
- The goal-register scaffold does not ask agents to duplicate full formal plan status.
- The consumer `docs/repo/README.md` scaffold routes to `plans/plan-register.md` and
  `plans/coordination-sync-pending.md`.
- The placement contract identifies both new files as kit-initialized consumer state files:
  Operating Kit owns the file contract and presence behavior, while the consumer repository owns
  the contents after creation.
- The docs preserve existing consumer ownership and no-forced-migration behavior.
- Source and packaged-resource mirrors match.

### E) Dependencies And Critical-Path Notes

Depends on `EP1` and `EP2` so the docs can refer to real managed references and scaffold paths.

### F) Tasks Checklist

- [x] Update agent-memory managed README with the formal-plan versus memory-state boundary.
- [x] Update agent-memory entry-format reference with the same boundary.
- [x] Update `session-ledger-maintenance.md` with the formal-plan versus memory-state boundary.
- [x] Update `components/agent-memory/scaffolds/goal-register.md` to describe informal or pre-plan continuity.
- [x] Update consumer `docs/repo/README.md` scaffold with plan-register and pending-sync routes.
- [x] Update placement contract with ownership, install-when-absent, repair-when-missing, and
  preserve-when-present rules for the new files.
- [x] Mirror changed agent-memory and structure-governance files under packaged resources.
- [x] Review wording to ensure it does not imply a required migration.

### G) Implementation Notes

The goal-register scaffold should stay useful. It should not become a tombstone. Recommended
language: use it for informal, pre-plan, or transitional continuity; promote durable formal plans
to `docs/repo/plans/plan-register.md`.

The session-ledger runbook should not depend on the plan register for all session memory. It should
only remove the overlapping authority: when a session created or materially changed a formal plan,
record formal lifecycle and relationships through the plan-register workflow. Unplanned session
recovery and non-plan continuity can remain in agent memory.

### H) Open Questions

None.

## EP6 - Packaging, Release Notes, And Validation

### A) Epic ID, Title, And Outcome

`EP6` - Packaging, Release Notes, And Validation

Outcome: all source and packaged resources are synchronized, public docs route the new plan, release
notes explain consumer impact, and tests validate the behavior.

### B) Scope

In scope:

- Update packaged resource manifest and any generated inventories.
- Update public docs routers and plan indexes.
- Add release notes for the next release target.
- Add and maintain the Consumer Impact Record in this implementation plan.
- Add packaged-resource parity tests for changed source and mirror files.
- Run focused and full validation.
- Record validation results in this implementation plan if execution happens under this plan.

Out of scope:

- Creating a GitHub release.
- Tagging a version.
- Publishing release assets.
- Merging a PR.

### C) Files Touched

Expected files:

- `src/codeheart_operating_kit/resources/manifest.yaml`
- `manifest.yaml` for review only unless release execution is explicitly in scope
- `docs/README.md`
- `docs/repo/README.md`
- `docs/repo/plans/README.md`
- `release-notes.md`
- `tests/test_public_core.py`
- `tests/test_packaging_resources.py`
- `tests/test_markdown_headers.py`
- Additional tests touched by previous epics

### D) Acceptance Criteria And Size

Size: `M`

Acceptance criteria:

- Public docs indexes link this implementation plan.
- The resource manifest includes every new managed and scaffold file needed by package fallback.
- Packaged-resource parity tests compare every changed source/mirror pair touched by this plan.
- Packaged fallback tests prove the new plan-register and coordination-sync-pending scaffolds are
  available when running from packaged resources.
- A Consumer Impact Record is present in this implementation plan and is reflected in release
  notes.
- Consumer impact distinguishes additive absent-file sync adoption from data migration: existing
  consumers run normal update/sync to receive absent scaffolds and managed-block refresh, and sync
  must not rewrite existing consumer-owned register files.
- Release notes describe:
  - plan-register doctrine;
  - new absent-file scaffolds;
  - existing consumer preservation;
  - optional portfolio configuration;
  - no forced migration;
  - no normal onboarding changes.
- Release handoff explicitly states that final version bump, root release manifest regeneration,
  `bootstrap.md`, `install.sh`, `install.ps1`, asset checksums, tags, and publication are release
  runbook tasks unless the user explicitly expands execution scope.
- Public-core hygiene validation passes.
- Markdown header validation passes.
- JSON schema validation passes.
- Release manifest validation passes.
- Focused tests pass for init, onboarding, sync, packaged resources, schemas, public-core hygiene,
  and Markdown headers.
- Full test suite passes or residual risk is documented if a full suite cannot be run.

### E) Dependencies And Critical-Path Notes

Depends on all previous epics. This epic should be last because the manifest, release notes, and
indexes need the final file list.

### F) Tasks Checklist

- [x] Update `src/codeheart_operating_kit/resources/manifest.yaml`.
- [x] Review root `manifest.yaml` and record whether it remains release-runbook-owned or is updated in this implementation.
- [x] Update `docs/README.md` with the implementation plan route.
- [x] Update `docs/repo/README.md` with the implementation plan route.
- [x] Update `docs/repo/plans/README.md` with the implementation plan route.
- [x] Update the Consumer Impact Record in this implementation plan.
- [x] Update `release-notes.md` with the additive consumer-impact summary for the target release.
- [x] Add packaged-resource parity tests for changed source and mirror files.
- [x] Add packaged fallback assertions for plan-register and coordination-sync-pending scaffolds.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 scripts/validate-json-schemas.py`.
- [x] Run `python3 scripts/validate-release-manifest.py`.
- [x] Run `pytest tests/test_init.py tests/test_onboard.py tests/test_sync_check.py tests/test_packaging_resources.py tests/test_json_schemas.py tests/test_public_core.py tests/test_markdown_headers.py`.
- [x] Run the full `pytest` suite when the focused suite passes.
- [x] Run `git diff --check`.
- [x] Record any validation failures and fixes in the execution log if this plan is executed.

### G) Implementation Notes

Keep release publication separate. This implementation plan can prepare release notes and a clean
tree, but `docs/repo/runbooks/release-operating-kit.md` owns final version tagging, asset build,
GitHub release publication, and post-release verification.

Do not leave release ownership implicit. If implementation does not update root `manifest.yaml`,
`bootstrap.md`, installers, package version fields, or release checksums, record that they are
release-runbook-owned in the handoff.

### H) Open Questions

None.

## 3.1 Release Handoff

When all epics are complete and validated, hand off to the release runbook with:

- final changed file inventory;
- consumer-impact classification;
- test and validation output summary;
- explicit statement that the release is additive and no forced migration is required;
- explicit statement that existing consumer-owned register files are preserved;
- explicit statement that normal onboarding does not ask portfolio-coordination questions;
- explicit statement that the release does not add a portfolio CLI command.
- explicit statement whether root `manifest.yaml`, package version fields, bootstrap, installers,
  and release checksums were updated or deferred to the release runbook.

Do not tag, publish, or merge as part of this implementation plan unless the user explicitly asks
for release execution after reviewing the implementation result.

Implementation execution note: the packaged resource manifest's component checksum, consumer
impact, and generated-surface metadata were updated for fallback installs. The root `manifest.yaml`
and the bundled resource manifest's release metadata fields, including `version`, `released_at`,
asset URLs, and asset checksums, remain release-runbook-owned because this implementation does not
execute a public release. Package version fields, `bootstrap.md`, `install.sh`, `install.ps1`,
release asset checksums, tag creation, and publication are also deferred to the release runbook.

Current validation evidence:

- `python3 scripts/validate-markdown-headers.py` passed.
- `python3 scripts/validate-public-core.py` passed.
- `python3 scripts/validate-json-schemas.py` passed.
- `python3 scripts/validate-release-manifest.py` passed.
- Focused suite passed: `pytest tests/test_init.py tests/test_onboard.py tests/test_sync_check.py tests/test_packaging_resources.py tests/test_json_schemas.py tests/test_public_core.py tests/test_markdown_headers.py`
  reported `49 passed`.
- Full suite passed: `pytest` reported `86 passed`.
- `git diff --check` passed.

## 3.2 Consumer Impact Record

Maintain this record during execution. Update affected paths and validation evidence as files are
implemented.

| Impact class or category | Affected paths | Required validation | Release/adoption note | Known consumer action |
| --- | --- | --- | --- | --- |
| `instruction-only change` | Managed planning-workflow docs, agent-memory docs, agent-interface docs, managed `AGENTS.md` template | Markdown headers, public-core hygiene, packaged-resource parity, focused docs tests | Release notes explain plan-register doctrine, coordination hook, and agent-memory boundary | Update to the release and run normal sync/check |
| Backwards-compatible kit-initialized consumer state file addition | `docs/repo/plans/plan-register.md`, `docs/repo/plans/coordination-sync-pending.md` | Init, onboarding, sync, packaged fallback, and preservation tests | New files are created or recreated only when absent; existing consumer file contents are preserved | Existing consumers run normal sync; absent files are initialized |
| Additive config/schema compatibility | Optional presence-based `portfolio` block in `.codeheart/kit.config.yaml` schema | Fixture-based omitted/member/coordination-home/invalid-role/incomplete-config config tests | No default portfolio config is written; no normal onboarding prompt is added; no `enabled`, `standalone`, `member_register_path`, or `pending_sync_path` fields are accepted | No action unless the consumer opts into portfolio coordination |
| Managed-block template refresh | `AGENTS.md` managed block generated from `templates/agents/AGENTS.managed-block.md` | Sync test proves managed block refresh and outside-content preservation | Existing consumers receive the lean hook through sync | Run normal sync/check after updating |
| Existing-consumer adoption note | Normal update/sync/check flow for installed consumers | Sync preservation tests, managed-block refresh tests, and generated-surface merge tests | This is additive sync adoption, not a data migration; no consumer-owned register content is moved, rewritten, or archived | Existing consumers update to the release and run normal sync/check |
| No forced migration | Existing `goal-register.md`, `plan-register.md`, and `coordination-sync-pending.md` content | Byte-for-byte preservation tests | Release notes state no forced migration and no automatic goal-register rewrite | None beyond normal update/sync |

# Section 4 - Future Planning

## 4.1 Deferred Tasks

The following work is intentionally deferred:

- Add a CLI command such as `codeheart-operating-kit portfolio configure` after the runbook-first
  model has been piloted.
- Add explicit multi-repository enrollment workflows that inspect user-approved local paths and
  produce approval-gated change plans.
- Add a plan-register validator after real register entries reveal what should be machine-checked.
- Add coordination-home aggregation automation after HQ-style use proves the update shape.
- Add richer portfolio status reports after the register model stabilizes.
- Add migration guidance for existing `goal-register.md` contents only if consumers ask for help
  promoting informal goals into formal plan-register entries.
- Decide whether `goal-register.md` should remain permanently or be further narrowed in a later
  major release.

## 4.2 Future Considerations

Future implementation plans should revisit:

- whether a coordination-home repository should offer a guided enrollment command;
- whether member repositories should expose a read-only status command for coordination homes;
- whether plan-register entries should use generated IDs or user-authored IDs;
- whether stale lifecycle snapshots should be detected by a validator;
- whether release notes should include a reusable adoption guide for existing consumers;
- whether the Operating Kit should provide an advanced settings flow that includes portfolio
  coordination without adding friction to normal onboarding.

# Revision Notes

- 2026-06-21: Initial implementation plan drafted from
  `portfolio-coordination-plan-register_discovery_doc.md`, covering managed doctrine, consumer
  scaffolds, safe sync behavior, planning-workflow hooks, optional portfolio config schema,
  agent-memory transition docs, packaging, release notes, and validation.
- 2026-06-21: Narrowed EP3 so review-only planning work remains read-only and
  `review-planning-document.md` is not edited for register updates; register maintenance belongs to
  user-authorized discovery or implementation planning edits that materially change canonical plans.
- 2026-06-21: Added remaining review hardening: clean agent-memory versus plan-register boundary,
  full coordination-home and pending-sync hook behavior, concrete pending-sync entry format,
  correct root template paths and packaged mirrors, existing-consumer `AGENTS.md` sync behavior,
  planning-workflow scaffold metadata, onboarding and schema fixture tests, release-boundary
  handoff, Consumer Impact Record, and packaged-resource parity requirements.
- 2026-06-21: Clarified existing-consumer sync as additive adoption rather than data migration and
  removed the remaining conditional wording around the known planning lifecycle reference.
- 2026-06-21: Marked the implementation plan active for later execution; no implementation tasks
  were started.
- 2026-06-21: Finalized review-driven schema and ownership recommendations: replaced ambiguous
  scaffold language with kit-initialized consumer state file terminology, removed `enabled`,
  `standalone`, `member_register_path`, and `pending_sync_path` from the proposed portfolio schema,
  and added role-specific required-field and negative-test expectations.
- 2026-06-21: Executed EP1 through EP6, recorded validation evidence, kept release publication
  tasks under the release runbook, fixed final review findings, and completed the plan.
