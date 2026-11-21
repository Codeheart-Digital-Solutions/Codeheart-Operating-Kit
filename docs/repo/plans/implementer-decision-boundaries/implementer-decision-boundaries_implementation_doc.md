Last updated: 2026-11-21T10:41:28Z (UTC)
Created: 2026-11-21
Status: active
Execution log: implementer-decision-boundaries_execution_log.md

# Director and Implementer Decision Boundaries

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.implementer-decision-boundaries
  kind: implementation
  purpose: Delegate coherent engineering outcomes and corrections while preserving architectural decisions, review and authority boundaries.
  first_cataloged: 2026-11-21T09:58:31Z
  catalog_metadata_updated: 2026-11-21T09:58:31Z
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.development-feedback-and-validation
```
<!-- END CODEHEART PLAN METADATA -->

Directors need to settle decisions that affect the outcome. Implementers need room to design, build and correct the solution within those decisions. This plan changes the guidance that each role reads.

| Essential context | Why |
| --- | --- |
| [Drafting](../../../../components/planning-workflows/managed/runbooks/draft-implementation-plan.md) and [plan review](../../../../components/planning-workflows/managed/runbooks/review-planning-document.md) | The authoring rules and the criteria a reviewer applies to them must agree. |
| [Execution](../../../../components/planning-workflows/managed/runbooks/execute-implementation-plan.md) and [agent task coordination](../../../../components/agent-interface/managed/reference/agent-task-coordination.md) | Existing authority for whole-plan assignments, review, corrections and report-back. |
| [Cross-tool permission profiles](../../../../components/agent-interface/managed/reference/cross-tool-coordination-contract.md#permission-profiles) and [Claude task operations](../../../../components/agent-interface/managed/reference/claude-code-task-operations.md) | Both describe command selection for explicit profiles. The contract owns the rule and the Claude guide applies it. |
| [Change](../../runbooks/change-operating-kit.md) and [release](../../runbooks/release-operating-kit.md) procedures | Source ownership, resource checks and candidate qualification. |

Contents: [Foundation](#section-1---foundation), [Strategy](#section-2---strategy), [Execution](#section-3---execution-plan), [Future planning](#section-4---future-planning).

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

A director can commission an outcome once, a reviewer can accept a plan that leaves local choices open, and an implementer can act on an in-scope finding in the same assignment. A consequential change in approach or authority reaches the director before dependent work proceeds.

## 1.2 Project And Problem Context

Discovery establishes intent, constraints and accepted capability. The plan translates those into an approach and evidence of success. Repeated command-by-command briefs during corrections duplicate engineering work even when the implementer retains the original session.

## 1.3 Current State Analysis

The Kit already permits whole-plan assignments, safe local decisions and focused correction review. Its mandatory visual file hierarchy, per-epic file list, checklist and corresponding review language pull authors toward predicting implementation. The execution rules identify individual material changes more clearly than cumulative growth in a procedure.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

```text
components/planning-workflows/managed/runbooks/
  draft-implementation-plan.md    # modify planning depth, required format and quality gate
  review-planning-document.md     # modify reviewer criteria to match
  execute-implementation-plan.md  # align progress, review, resume and replanning guidance
components/agent-interface/managed/reference/
  agent-task-coordination.md      # align the existing assignment example and handoff
  cross-tool-coordination-contract.md  # clarify reuse of approved role/repository permissions
  claude-code-task-operations.md       # align the Claude summary and point to the contract
src/codeheart_operating_kit/resources/components/  # synchronize matching packaged resources
```

The coordination home owns its existing `docs/repo/reference/agent-coordination.md`. Its owner records the local policy there. Machine settings and credentials belong to that owner.

## 2.2 Open Questions And Assumptions Requiring Clarification

**OQ-1 · BLOCKER: no · Affects: EP-02 consumer follow-up.** Can the local tools support worktree development with the desired access limits? The owner can qualify a candidate profile independently. Guidance delivery proceeds with existing approved settings.

## 2.3 Architectural Decisions With Reasoning

- The director owns intent, consequential design choices, dependencies, acceptance evidence and operating limits. Ask what would go wrong if the implementer chose differently. Name exact paths or commands when they identify an interface, protected target or necessary constraint. Ordinary local files, functions, test design and debugging belong to the implementer.
- Commission one coherent outcome, which may span related epics. A reviewer states the defect, the behavior to preserve and the proof required. For example, “This path can submit twice after a lost response. Preserve the single-submission guarantee and show the regression test.” The same implementer chooses the repair. Material correctness, security and integrity findings qualify even if the plan missed them. A discretionary style preference alone warrants neither a correction round nor a log entry.
- The plan records each epic's current status and acceptance criteria. Mark the epic complete only after validation and required review. The sibling execution log holds evidence, material deviations and, when useful for resuming or handing off work, the implementer's own working list. Existing plan checkboxes still track validated tasks where present. Update epic status at meaningful checkpoints even when the approved approach stays the same.
- An approved role and repository profile can cover ordinary local work throughout the assignment. At commissioning, check that it covers the intended work and identify any separately authorized external effects. Name individual commands when their identity matters to safety or correctness. Preserve enforced restrictions and the policy owner's authority over profile changes.
- When an assumption fails or individually reasonable fixes together materially change the approach, both planner and implementer reconsider it. Would they choose it from the start with today's facts? Does all remaining work, including recovery, still fit the limits? The director revises the affected part within delegation. Only that work and its dependents wait.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Depends on | Status |
| --- | --- | --- | --- | --- |
| EP-01 | Drafting, review and execution give consistent decision authority. | S | Agreed direction | Complete |
| EP-02 | Released guidance is usable in the designated coordination home. | S | Reviewed EP-01 source for release and adoption | In progress |

### EP-01 · Plan, delegate and correct a coherent outcome

**A) Outcome.** Drafting, review and execution give both roles the same decision boundary.

**B) Scope.** Revise the existing guidance in one batch and preserve existing plans.

**C) Files touched.** The six managed documents in §2.1 and their packaged copies.

**D) Acceptance and size.** S. A fresh planner can leave local details open, a reviewer can still reject an unresolved consequential design choice, and an implementer can distinguish a correction from a replan.

**E) Dependencies.** Agreed planning and delegation decisions in §2.3.

**F · Tasks**

- [x] Make exhaustive file inventories, file trees and planner-written implementation steps conditional. Keep outcome, approach, dependencies, constraints and proof concrete. Align drafting, per-epic progress, review, log and resume instructions with §2.3 for plans with or without task checkboxes. Small changes may use the assignment itself.
- [x] Align assignments, corrections and both permission references with §2.3. Replace “List only the task's ordinary development and delivery commands” in the cross-tool contract and the matching wording in the Claude guide with instructions to check the approved profile once at commissioning and reuse it for covered work without listing those commands per task. Link the owning plan or Program section and relevant decisions in the handoff. Ask for an approach note only to resolve a consequential uncertainty. After a denial, stop the operation and its dependents, report what is blocked and continue independent authorized work. The denial applies across commands, tools, identities and agents. This epic changes wording only.
- [x] Verify packaged copies match source. Use the affected routing/resource checks and one fresh-context walkthrough covering three decisions: repair a duplicate submission under the existing contract, return a recovery design that outgrows its time window, and continue independent local work after a denied publication. Review the coherent diff together.

**G · Implementation notes.** For the duplicate-submission case, use an epic without a predicted file list or task checklist. A fresh reader should be able to locate its review-pending status in the plan, resume from the log, and establish when the epic can be marked complete. The reviewer should still flag a second sample that leaves a consequential data decision unresolved.

**H · Open questions.** None.

### EP-02 · Deliver the guidance to the coordination home

**A) Outcome.** The released rules and the designated consumer's standing instructions work together.

**B) Scope.** Producer release and adoption in one coordination home.

**C) Files touched.** Producer release records, the consumer's installer-managed targets and its existing coordination reference in §2.1.

**D) Acceptance and size.** S. Record the released version, consumer default-branch revision, installation check and owner-reviewed local guidance.

**E) Dependencies.** Release follows source review. The local owner may prepare the reference change while EP-01 is in review.

**F · Tasks**

- [ ] At commissioning, name the producer and consumer owners, targets and authorized Git, release and adoption effects once. Use the normal release with release notes and the eligible checks for the full candidate. Treat this as an instruction change that preserves existing permission policy.
- [ ] Adopt the release through the supported installation workflow. In the consumer reference, link the managed rules, specify which local work the role profile covers, and distinguish covered Git delivery from reserved publication, deployment or shared-resource changes.
- [ ] Verify the installed planning, coordination and permission guidance through the home's existing entry point. EP-02 closes on verified adoption and owner-reviewed local guidance. The separate trial in §4.1 may still be unrun at closure.

**G · Implementation notes.** The commissioning director accepts the source after focused review. The home owner verifies adoption. A separate release solely for this guidance is optional. Existing release and adoption gates apply.

**H · Open questions.** OQ-1 affects selection of a runtime profile only.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

The local owner can trial each proposed runtime profile once in a disposable linked worktree. Exercise edit, representative build/test and local commit. Use harmless fixtures to verify that writes outside the permitted area and reads of a dummy credential are blocked through the worker's shell and file tools. Record any build or Git metadata restriction. Keep live credentials, including cached sessions, outside the trial's reach. A change to runtime defaults requires a local owner decision on the results.

## 4.2 Future Considerations

At the next ordinary handoff, look for a correction completed within the original assignment and for any binding decision the implementer had to guess. Use that evidence if further guidance is needed.

# Revision Notes

- 2026-11-21: Draft from the agreed planning, correction and permission boundaries.
- 2026-11-21: Added the two permission references, defined epic progress with optional checklists, and clarified correction and adoption criteria.
- 2026-11-21: Owner approved the revised plan and commissioned delivery through release and verified adoption. See the [execution log](implementer-decision-boundaries_execution_log.md).
