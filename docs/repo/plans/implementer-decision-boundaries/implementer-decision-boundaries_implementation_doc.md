Last updated: 2026-11-21T09:58:31Z (UTC)
Created: 2026-11-21
Status: draft

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
  execute-implementation-plan.md  # modify correction and replanning guidance
components/agent-interface/managed/reference/
  agent-task-coordination.md      # align the existing assignment example and handoff
src/codeheart_operating_kit/resources/components/  # synchronize matching packaged resources
```

The coordination home owns its existing `docs/repo/reference/agent-coordination.md`. Its owner records the local policy there. Machine settings and credentials belong to that owner.

## 2.2 Open Questions And Assumptions Requiring Clarification

**OQ-1 · BLOCKER: no · Affects: EP-02 consumer follow-up.** Can the local tools support worktree development with the desired access limits? The owner can qualify a candidate profile independently. Guidance delivery proceeds with existing approved settings.

## 2.3 Architectural Decisions With Reasoning

- The director owns intent, consequential design choices, dependencies, acceptance evidence and operating limits. Ask what would go wrong if the implementer chose differently. Name exact paths or commands when they identify an interface, protected target or necessary constraint. Ordinary local files, functions, test design and debugging belong to the implementer.
- Commission one coherent outcome, which may span related epics. A reviewer states the defect, the behavior to preserve and the proof required. The same implementer chooses the repair within the agreed approach. Material correctness, security and integrity findings qualify even if the original checklist missed them.
- When an assumption fails or individually reasonable fixes together materially change the approach, both planner and implementer reconsider it. Would they choose it from the start with today's facts? Does all remaining work, including recovery, still fit the limits? The director revises the affected part within delegation. Only that work and its dependents wait.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Depends on |
| --- | --- | --- | --- |
| EP-01 | Drafting, review and execution give consistent decision authority. | S | Agreed direction |
| EP-02 | Released guidance is usable in the designated coordination home. | S | Reviewed EP-01 source for release and adoption |

### EP-01 · Plan, delegate and correct a coherent outcome

**A) Outcome.** Drafting, review and execution give both roles the same decision boundary.

**B) Scope.** Revise the existing guidance in one batch and preserve existing plans.

**C) Files touched.** The four managed references in §2.1 and their packaged copies.

**D) Acceptance and size.** S. A fresh planner can leave local details open, a reviewer can still reject an unresolved consequential design choice, and an implementer can distinguish a correction from a replan.

**E) Dependencies.** Agreed planning and delegation decisions in §2.3.

**F · Tasks**

- [ ] Make exhaustive file inventories, file trees and planner-written implementation steps conditional. Keep outcome, chosen approach, dependencies, binding constraints and proof concrete. Align the review criteria and quality gates. Small changes may use the assignment itself.
- [ ] Align execution and assignment guidance with the decisions above. Link the owning plan or Program section and relevant decisions in the handoff. Provide a short approach note only when it resolves a real uncertainty. Reuse covered local permissions through corrections. Respect a denial across commands, tools, identities and agents, report the blocked effect and continue independent authorized work. External effects require matching assignment authority.
- [ ] Verify packaged copies match source. Use the affected routing/resource checks and one fresh-context walkthrough covering three decisions: repair a duplicate submission under the existing contract, return a recovery design that outgrows its time window, and continue independent local work after a denied publication. Review the coherent diff together.

**G · Implementation notes.** Review a sample plan without a predicted file inventory and one that omits a material data decision. The reviewer should accept the first if executable and flag the second. Retain existing evidence for unaffected behavior.

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
- [ ] Verify the installed drafting, review and execution guidance and its entry from the home. Record the local owner’s disposition of the separate permission trial.

**G · Implementation notes.** The commissioning director accepts the source after focused review. The home owner verifies adoption. A separate release solely for this guidance is optional. Existing release and adoption gates apply.

**H · Open questions.** OQ-1 affects selection of a runtime profile only.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

The local owner can trial each proposed runtime profile once in a disposable linked worktree. Exercise edit, representative build/test and local commit. Use harmless fixtures to verify that writes outside the permitted area and reads of a dummy credential are blocked through the worker's shell and file tools. Record any build or Git metadata restriction. Keep live credentials, including cached sessions, outside the trial's reach. A change to runtime defaults requires a local owner decision on the results.

## 4.2 Future Considerations

At the next ordinary handoff, look for a correction completed within the original assignment and for any binding decision the implementer had to guess. Use that evidence if further guidance is needed.

# Revision Notes

- 2026-11-21: Draft from the agreed planning, correction and permission boundaries.
