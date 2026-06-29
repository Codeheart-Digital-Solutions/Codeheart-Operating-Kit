Last updated: 2026-06-29T20:06:45Z (UTC)
Created: 2026-06-29
Status: completed
Completed: 2026-06-29
Execution log:
`docs/repo/plans/plan-register-dirty-target-safety/plan-register-dirty-target-safety_execution_log.md`

# Document Header

## Plan Register Dirty Target Safety Implementation Plan

Overview: Clarify managed Operating Kit plan-register doctrine so agents do not treat an
unrelated dirty coordination-home repository as an automatic blocker. Agents should inspect the
target register and distinguish compatible, non-overlapping changes from ambiguous target-entry
conflicts before falling back to `coordination-sync-pending`.

This plan is intentionally narrow. It does not promote the HQ Portfolio Work Board, add a generic
work-board hook, change sync/check behavior, add validators, add a new scaffold path, or change
scaffold ownership. It does align the existing pending-sync scaffold wording so new absent-file
baselines do not preserve the stale broad fallback.

Essential context:

| Source | Why it matters |
| --- | --- |
| `AGENTS.md` | Maintainer bootstrap, public-core safety, and change-safety rules. |
| `docs/repo/reference/consumer-impact-classification.md` | Consumer-impact class and release-note expectations. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure for managed Operating Kit changes. |
| `components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Primary doctrine that defines local register, portfolio coordination, and pending-sync behavior. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | Planning hook that repeats the current pending-sync fallback during discovery. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Planning hook that repeats the current pending-sync fallback during implementation-plan drafting. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Execution hook that repeats the current pending-sync fallback during implementation execution. |
| `components/planning-workflows/managed/reference/plan-register-format.md` | Reference vocabulary for register entries, coordination notes, and pending-sync state. |
| `components/planning-workflows/scaffolds/coordination-sync-pending.md` | Existing absent-file baseline that must not preserve stale broad pending-sync wording. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/` | Packaged managed-resource mirrors that must match source managed docs. |
| `scripts/validate-markdown-headers.py` | Markdown timestamp/header validation. |
| `scripts/validate-public-core.py` | Public-core hygiene validation. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Implement a small managed doctrine refinement for dirty plan-register handling.

Completion is proven when:

- `maintain-plan-register.md` says a dirty repository alone is not enough to skip a
  coordination-home register update;
- agents must inspect the target register state before choosing direct update versus pending
  sync;
- compatible cases are defined as unique, non-overlapping inserts or clearly compatible field
  refreshes that preserve existing dirty content;
- ambiguous cases are defined as dirty target entries, lifecycle/status/relation conflicts,
  duplicate IDs or canonical paths, ownership ambiguity, or updates that require interpreting
  another unfinished edit;
- pending sync is framed as a fallback for unsafe target-register updates, not a default response
  to unrelated dirty worktree state;
- discovery, implementation-planning, and implementation-execution hooks point back to this
  compatibility test instead of repeating the broad unsafe-edit rule;
- `plan-register-format.md` keeps `coordination-sync-pending` as valid vocabulary while clarifying
  its boundary;
- source and packaged managed-resource copies match;
- validation proves headers, public-core hygiene, resource parity, and no whitespace errors.

## 1.2 Project And Problem Context

The current plan-register doctrine correctly protects unrelated work and prevents agents from
writing into an unsafe coordination-home register. However, its wording is broad enough for agents
to avoid a direct coordination-home update whenever the coordination-home repository is dirty,
even when the target register update would be independent and safe.

That behavior is too conservative. In portfolio coordination, unrelated dirty files should not
force pending sync. The safer distinction is whether the target register edit itself overlaps
uncommitted register content or requires the agent to interpret unfinished work from another
actor.

## 1.3 Current State Analysis

Current managed state:

- `maintain-plan-register.md` says to update the coordination-home register when it is locally
  available and safe to edit.
- The same runbook says to record pending sync when the coordination home is missing,
  inaccessible, outside writable scope, or otherwise unsafe to edit.
- The safety rules include not writing to a coordination home when its worktree state makes the
  change unsafe.
- Discovery, draft-implementation, and execute-implementation runbooks repeat the broad
  "available and safe to edit" versus "unavailable or unsafe" wording.
- `plan-register-format.md` defines `coordination-sync-pending` as valid sync-state vocabulary.

Problem:

- "Dirty worktree" and "unsafe target edit" are not separated.
- A target register can be dirty but still accept a clearly compatible non-overlapping insert.
- A target register can be dirty in a way that makes the intended update unsafe, especially when
  the same entry, lifecycle state, relation set, ID, or canonical path is already modified.

Target state:

- Agents use a target-register compatibility test before pending sync.
- Agents preserve unrelated dirty work and do not rewrite or normalize register content outside
  the intended entry.
- Pending sync remains the conservative fallback when target-entry intent is unclear.

Consumer impact classification: `instruction-only change`.

Release-note requirement: required when shipped.

Migration requirement: none.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Patch the managed planning-workflows source docs and existing pending-sync scaffold wording, mirror
those changes into packaged resources, and update repository planning indexes/registers for this
plan. Do not add new managed files, components, scaffold paths, validators, or CLI behavior.

```text
Codeheart-Operating-Kit/
  docs/
    README.md                                                        # modify index
    repo/
      README.md                                                      # modify index
      plans/
        README.md                                                    # modify index
        plan-register.md                                             # modify entry
        plan-register-dirty-target-safety/
          plan-register-dirty-target-safety_implementation_doc.md    # this plan
  components/
    planning-workflows/
      managed/
        reference/
          plan-register-format.md                                    # modify
        runbooks/
          discovery-workflow.md                                      # modify
          draft-implementation-plan.md                               # modify
          execute-implementation-plan.md                             # modify
          maintain-plan-register.md                                  # modify
      scaffolds/
        coordination-sync-pending.md                                  # modify wording only
  src/
    codeheart_operating_kit/
      resources/
        components/planning-workflows/managed/...                    # mirror modified docs
        components/planning-workflows/scaffolds/coordination-sync-pending.md
                                                                        # mirror wording only
```

## 2.2 Open Questions And Assumptions Requiring Clarification

### OQ-001 - Should Direct Updates Be Allowed When The Target Register Is Dirty?

BLOCKER: no

Decision: yes. Direct updates should be allowed when the intended change is clearly compatible
with the current target-register diff and preserves existing dirty content.

### OQ-002 - Should Agents Update A Dependent Work Board During This Plan?

BLOCKER: no

Decision: no. Work-board behavior remains repository-owned and is explicitly out of scope for
this implementation.

## 2.3 Architectural Decisions With Reasoning

### AD-001 - Use A Target-Register Compatibility Test

Decision: define "safe to edit" by inspecting the target register, not by the whole repository's
dirty state.

Reasoning: unrelated dirty files do not affect register correctness. Dirty target entries can
affect correctness and must be handled more carefully.

### AD-002 - Preserve Pending Sync As The Conservative Fallback

Decision: keep `coordination-sync-pending` for missing access, unwritable coordination homes, and
ambiguous target-register conflicts.

Reasoning: the fallback is still useful. The refinement only prevents overuse when a direct
register update is clearly safe.

### AD-003 - Keep This Instruction-Only

Decision: do not add validators, CLI support, or new scaffolds in this plan. Update only the
existing pending-sync scaffold wording to avoid shipping stale fallback language to absent-file
consumers.

Reasoning: the behavior is an agent decision rule inside existing managed planning guidance.
Validation can be done through docs checks, resource parity checks, and review.

# Section 3 - Execution Plan

## EP1 - Clarify Dirty Target Register Doctrine

Objective: Update `maintain-plan-register.md` with the concrete compatibility test.

Tasks:

- Add a section that distinguishes dirty repository, dirty target register, and dirty target
  entry.
- Define clearly compatible cases:
  - target register is clean;
  - target register is dirty, but the intended entry ID and canonical path are absent and the
    update is a unique non-overlapping insert;
  - target register is dirty, but the intended field refresh does not touch lines already changed
    by another uncommitted edit and is directly supported by the canonical plan;
  - unrelated dirty files exist outside the target register.
- Define unsafe or ambiguous cases:
  - intended entry, ID, title, canonical path, lifecycle state, relation set, or coordination note
    is already modified in the dirty target-register diff;
  - the intended update would require reordering, renumbering, normalizing, or reconciling
    existing dirty content;
  - duplicate IDs or canonical paths appear;
  - canonical plan and dirty register content disagree in a way that changes lifecycle,
    relationship, or implementation-path meaning;
  - local instructions, permissions, or user direction make the coordination-home edit unsafe.
- Require agents to mention existing dirty target-register state when they update through it.
- Require pending sync or a targeted user question when compatibility cannot be established.

Validation:

- `rg -n "dirty|target register|pending sync|coordination-home" components/planning-workflows/managed/runbooks/maintain-plan-register.md`
- Manual review that pending sync remains available but is not the default for unrelated dirty
  files.

## EP2 - Align Planning Workflow Register Hooks

Objective: Update repeated plan-register hook wording in discovery, drafting, and execution
runbooks.

Tasks:

- In `discovery-workflow.md`, route portfolio coordination updates through the new
  target-register compatibility test before pending sync.
- In `draft-implementation-plan.md`, use the same wording.
- In `execute-implementation-plan.md`, use the same wording.
- Keep the local register first, coordination-home second sequence intact.
- Avoid duplicating the full compatibility test in every runbook; point to
  `maintain-plan-register.md` as the authority.

Validation:

- `rg -n "coordination home is available|unsafe to edit|pending sync|compatib" components/planning-workflows/managed/runbooks`
- Manual review that the hook text is consistent and not broader than EP1.

## EP3 - Clarify Register Format Boundary

Objective: Keep `coordination-sync-pending` vocabulary while clarifying when it belongs.

Tasks:

- Add one concise note in `plan-register-format.md` explaining that pending sync is for unsafe or
  unavailable coordination-home register updates, not for unrelated dirty worktree state.
- Align the existing `coordination-sync-pending.md` scaffold wording with the target-register
  compatibility boundary.
- Preserve existing entry shapes and examples.
- Do not add new lifecycle values or sync-state values.

Validation:

- `rg -n "coordination-sync-pending|dirty|unsafe" components/planning-workflows/managed/reference/plan-register-format.md`
- Manual review that the format reference does not become a procedure duplicate.

## EP4 - Mirror Resources, Indexes, And Validate

Objective: Make the source implementation shippable after EP1 through EP3.

Tasks:

- Copy changed managed planning-workflow docs and scaffold wording into matching
  `src/codeheart_operating_kit/resources/components/planning-workflows/` paths.
- Update `docs/README.md`, `docs/repo/README.md`, and `docs/repo/plans/README.md` with this
  implementation plan.
- Update `docs/repo/plans/plan-register.md` with `OK-PR-023`.
- Run source/resource parity checks for changed managed docs.
- Run Markdown header validation.
- Run public-core hygiene validation.
- Run `git diff --check`.
- Record release-note requirement for the next Operating Kit release, but do not release in this
  plan unless explicitly requested.

Validation:

- `python3 scripts/validate-markdown-headers.py`
- `python3 scripts/validate-public-core.py`
- `diff -u components/planning-workflows/managed/runbooks/maintain-plan-register.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md`
- `diff -u components/planning-workflows/managed/runbooks/discovery-workflow.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/discovery-workflow.md`
- `diff -u components/planning-workflows/managed/runbooks/draft-implementation-plan.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
- `diff -u components/planning-workflows/managed/runbooks/execute-implementation-plan.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
- `diff -u components/planning-workflows/managed/reference/plan-register-format.md src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md`
- `diff -u components/planning-workflows/scaffolds/coordination-sync-pending.md src/codeheart_operating_kit/resources/components/planning-workflows/scaffolds/coordination-sync-pending.md`
- `git diff --check`

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Generic dependent planning surfaces or work boards are deferred. HQ Portfolio Work Board remains
  HQ-owned local doctrine.
- Automated validation for register patch compatibility is deferred until repeated failures prove
  a script or validator is worth adding.
- Release, tag, and consumer sync are deferred until the implementation is completed and approved
  for shipping.

## 4.2 Future Considerations

- If multiple repositories repeatedly need the same dirty-register decision logic, consider a
  reusable script or validator that reports target-entry overlap without modifying files.
- If coordination-home syncing becomes frequent, consider a source-owned coordination register
  update recipe with explicit diff evidence and conflict classification.

## Revision Notes

- 2026-06-29: Drafted narrow implementation plan after user approved fixing dirty plan-register
  rules and explicitly set work-board behavior aside.
- 2026-06-29: Activated implementation and created the execution log.
- 2026-06-29: Accepted review-gate finding that the existing pending-sync scaffold also needed
  wording alignment; no new scaffold path or ownership behavior was added.
- 2026-06-29: Recorded validation and review evidence, resolved close-out metadata findings, and
  marked the implementation complete.
