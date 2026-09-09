Last updated: 2026-09-09T13:39:07Z (UTC)
Created: 2026-09-09
Status: draft

# Git Delivery Discipline Discovery

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.discovery.git-delivery-discipline
  kind: discovery
  purpose: Make commits, publication, timely integration and repository adoption predictable throughout agent work without adding procedural overhead or broad validation churn.
  first_cataloged: 2026-09-09T13:39:07Z
  catalog_metadata_updated: 2026-09-09T13:39:07Z
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.development-feedback-and-validation
    - kind: related
      target: codeheart-operating-kit.implementation.release-validation-effort
```
<!-- END CODEHEART PLAN METADATA -->

## Readiness and scope

This draft consolidates the user's discussion and the inspected current rules. Its target is
draft-ready, not an approved implementation handoff. The user requested predictable Git discipline
inside ongoing agent work and explicitly challenged the risk of reintroducing ceremony and costly
CI. Recommendations below remain proposals where indicated. No implementation is activated.

Producer baseline: `bbfce9b54b8e54174a58518f155be2874f310390`, fetched from `origin/main` on
2026-09-09. The work is isolated on `codex/git-delivery-discipline-discovery`. This authoring scope
covers the discovery, its nearest index, proportionate local validation and a local commit.
External publication, managed-rule changes, a release and consumer rollout are separate effects.

The Operating Kit owns the reusable rules. Consumer repositories own their concrete targets,
branch conventions, required checks and product plans. Private task locators and target details
are excluded from this public document. The earlier completed deliveries remain completed within
their recorded scopes; this discovery addresses an uncovered integration boundary.

## Problem, intention and success

The concern is costly, late integration caused by work accumulating in uncommitted files,
unpublished branches or long-lived branches containing otherwise finished changes. More commits
and pushes alone cannot prevent integration difficulties: finished outcomes must reach the shared
baseline at appropriate times, and work must account for relevant changes to that baseline.

The immediate example was a shared guidance upgrade verified in several active working copies
while the corresponding repository default branches still carried older guidance. Existing tasks
could use the new process, but fresh tasks based on main could inherit older instructions. The
commissioning finish line admitted working-copy adoption; it did not establish a repository-wide
default-branch rollout. This is a gap in both the chosen finish line and the precision of guidance,
not evidence that all existing commit/publication rules were ignored.

User intention: keep work recoverable, shared plans discoverable, main usable and integration
manageable while preserving fast iteration and low coordination cost. Commit, push, PR, merge,
release and adoption are means to those outcomes, not activity quotas.

Priorities are preserved work and accurate authority; timely shared integration; proportionate
validation; and minimal ongoing ceremony. Context ingestion, technical reasoning and substantive
planning must not be weakened to save tokens.

Success means agents can identify the next appropriate Git action during authoring, iteration,
review, planned pauses and handoffs. They do not wait for an entire epic merely to save progress,
leave independently usable accepted work unmerged without a concrete reason, or confuse a local
installation with completed repository adoption. These behaviors require no new Git service,
receipt registry, periodic watchdog or mandatory reporting artifact.

Non-goals: automatic merging of incomplete behavior; universal branch-age or commit-frequency
limits; a PR per file or checkbox; default force-push/rebase/branch deletion; bypassing required
checks; changing execution approval through document publication; a fleet-wide CI inventory;
resuming held product implementations; or optimizing unrelated runtime suites.

## Evidence and current rules

Paths in this document are relative to the producer root.

| Source | Current behavior and implication |
| --- | --- |
| `components/planning-workflows/managed/reference/planning-document-lifecycle.md`, Git Checkpoints And Independent Lifecycle | Requires a local commit at a useful completed checkpoint, authorized branch publication and accepted PR integration. Publication and lifecycle are separate. Timing remains broad. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`, Coherent Batches and Per-Epic Flow | Preserves evidence across resumes and rejects full suites solely for a commit. Its numbered epic flow places the Git checkpoint after implementation/review, which can encourage delaying smaller Git actions. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md`, Catalog And Visibility Hook | Discovery drafting includes a local commit, but does not itself authorize push or PR publication. Remote coordination requires publication and a successful refresh. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`, Authoring and Activation | Drafting and activation have different authority. Activation authorizes its bounded plan-only commit/push; it does not independently authorize PR merge or product execution. |
| `components/agent-interface/managed/reference/agent-task-coordination.md` | Requires a named delivery boundary and accurate delivery states. A whole-plan assignment can cover Git effects once. The one-delivery-PR example must not become a requirement to keep a long plan unintegrated. |
| `components/planning-workflows/managed/runbooks/handle-routine-change.md` | Uses existing records, proportionate checks and covered authority; suitable for carrying the same lightweight Git behavior into routine work. |
| `components/agent-interface/managed/runbooks/maintain-operating-kit-installation.md` | Owns supported install/check/upgrade and preservation. A successful local transaction alone does not establish which Git branches received the upgrade. |
| `docs/repo/plans/development-feedback-and-validation/` | The final closure at `f29172b3b591e3329dc8f3664f4a7a989d575a3a` records the released v0.1.32 guidance and adoption in assigned working copies. That bounded outcome did not establish adoption on all intended default branches. |
| `docs/repo/plans/release-validation-effort/` and `docs/repo/runbooks/release-operating-kit.md` | Completed producer improvement supports guarded guidance candidates and applicable evidence reuse. It did not change consumer CI workflows. Guidance demonstration jobs took 1m11s and 1m27s with restored caches; those observations do not predict consumer merge latency. |

Local plan discovery and remote portfolio discovery are different. Under the configured catalog
and source scope, pushed branches can supply canonical plans after a successful refresh. Main
merge is not intrinsically required merely for remote plan visibility. Local-only edits or commits
are not remotely observable. Main supplies the default shared baseline for fresh work.

No evidence here quantifies merge-conflict frequency or establishes the current cost of each
consumer's CI. The observed inconsistency is enough to justify clearer workflow defaults, not a
claim that every long-lived branch or every test run is wasteful.

## Requirements and operating scenarios

- FR-1: Distinguish local recovery, remote visibility, accepted integration, released artifacts and
  adoption at named targets. A report must not imply a broader state than was verified.
- FR-2: Put coherent local commits inside the execution loop, independently of formal review or
  epic acceptance. Before a planned pause/handoff, preserve recoverable work; incomplete work is
  identified as incomplete, and unrelated or sensitive material is excluded.
- FR-3: Under established publication authority, publish useful progress for coordination/recovery
  and update the existing PR where appropriate. Avoid a new review or permission round per push.
- FR-4: Identify independently usable integration outcomes in planning. A plan can span multiple
  coherent PRs; neither an entire plan nor each microscopic edit defines a mandatory PR boundary.
- FR-5: At the start or resumption of substantial work and before integration, inspect relevant
  target-branch changes. Resolve actual overlap/dependency drift before accumulating more dependent
  work. Preserve shared history and other agents' changes; no unconditional rebase is prescribed.
- FR-6: Integrate ready outcomes under covered authority and passing relevant gates. If ready work
  remains unmerged, record its concrete blocker, owner and next resolving event in the existing
  record. No new tracker or elapsed-time escalation service is required.
- FR-7: Make planning publication predictable. A coherent draft is locally committed; a shared
  handoff identifies its canonical location and actual Git visibility. An accepted shared plan can
  merge while retaining draft status and unresolved questions. Publication does not approve its
  recommendations, activate execution or dispatch implementation.
- FR-8: For a repository adoption assignment, reach the explicitly named default/integration
  branch and verify supported installation plus preservation. A working-copy pilot is explicitly
  partial. Shared guidance changes should reach main through isolated adoption changes before
  dependent work treats the baseline as generally available; existing branches then reconcile it.
- FR-9: Couple cadence to relevant validation, not arbitrary Git events. Ordinary commits, pushes,
  resumes or unchanged integrations do not by themselves justify broad runtime validation.

- NFR-1: Use one canonical rule source with concise links and execution hooks in the existing
  runbooks. Do not copy a large policy into each route or create another overarching process.
- NFR-2: Preserve explicit local-only requests, public/private boundaries, repository protections,
  installed-content ownership, current authorization and meaningful source review.
- NFR-3: Keep main usable. Split separable ready work or document a real dependency instead of
  integrating unsafe unfinished behavior to satisfy cadence. Each meaningful integration outcome
  receives its appropriate review; each commit does not require a separate review.
- NFR-4: Use affected checks during iteration and coherent acceptance at integration/release
  boundaries. Retain evidence when relevant inputs are unchanged. Unknown or consequential changes
  need investigation and appropriate wider checks; a documentation extension is not a risk waiver.

Representative scenarios for later validation:

1. A multi-step epic has useful progress before its formal review: the agent saves/publishes it
   under covered authority without marking the epic complete or replaying the full suite.
2. A draft plan with open questions is needed by another owner: it becomes discoverable through
   authorized publication; accepted shared placement on main does not activate implementation.
3. A shared operating update sits on an unfinished product branch: the agent isolates adoption
   onto the intended main baseline without merging the held product work.
4. A ready integration unit exists inside a longer plan: the agent integrates it after its own
   relevant review/checks and retains the remaining plan as active.
5. Consumer CI runs an unrelated expensive suite on every push: the owner identifies the mismatch
   and proposes a scoped correction or coherent batching while required checks remain enforced.
6. Main changes during a pause: the agent inspects relevant divergence, preserves concurrent work,
   resolves overlap and reruns only invalidated evidence plus required integration checks.

## Decision ledger

Decision owner: the commissioning user, with routine implementation details delegable later.
The agreed user intention is distinguished from recommendations not yet explicitly accepted.

| ID | Decision, state and class | Recommendation, alternatives and closure |
| --- | --- | --- |
| D-1 | Predictable Git discipline with low ceremony — user requirement; operating outcome | Make recovery, visibility and usable integration the outcomes. Reject activity quotas and per-operation approval. Preserve FR-1 and NFR-1 through NFR-4. |
| D-2 | Timing inside agent work — recommended workflow clarification | Use coherent progress and meaningful transitions, rather than waiting only for epic completion or imposing timers. Covers FR-2 through FR-6. Close through acceptance of this behavior and representative scenario review. |
| D-3 | Planning publication authority — recommended clarification; external-action boundary | Establish bounded publication/PR/merge authority once for shared repository planning assignments; retain the existing rule that drafting alone does not grant external publication. A broader standing authoring grant is outside this recommendation and would require an explicit policy choice and impact review. Covers FR-3 and FR-7; do not invent authority. |
| D-4 | Merge timing and main adoption — recommended integration boundary | Merge independently usable accepted outcomes and distinguish pilot installation from repository-main adoption. Avoid both whole-plan delay and artificial tiny PRs. Covers FR-4, FR-6 through FR-8; exact rollout targets are chosen in commissioning. |
| D-5 | Validation and CI interaction — recommended preservation of approved proportionate-validation direction | Reuse applicable evidence and inspect actual affected repository triggers before promising cheap publication. The Kit producer's faster route does not update consumer CI. CI fixes remain owner-scoped; no global audit or check bypass. Covers FR-9 and OQ-3. |
| D-6 | Delivery of the clarification — recommended bounded Kit follow-up | Update existing canonical guidance and workflow hooks, release through the applicable producer route, and adopt on designated repository main branches before claiming general rollout. A source-only or pilot-only boundary is an explicit alternative, not implied full adoption. Depends on accepted D-2 through D-5 and action-time authority. |

The tradeoff is intentional: early coherent publication makes work recoverable and visible, while
integration waits for usable behavior and meaningful acceptance. Repeated tiny PRs or blanket
branch-refresh rules would create the overhead the user wants removed. A short reason for an
exception is enough; this proposal does not add a Git event log or proof-of-compliance framework.

## Open questions and assumptions

- OQ-1 — BLOCKER: no. The proposed route preserves current authoring authority and makes bounded
  publication explicit once at commissioning. A broader standing grant is not required to solve
  the observed problem and is outside this proposal unless the user requests it. Owner: the
  commissioning owner resolves actual publication scope once before those effects; the absence of
  a broader policy is not a reason to restart discovery or request permission per Git operation.
- OQ-2 — BLOCKER: no for this discovery or a bounded draft implementation plan. Exact release
  version, consumer default branches, rollout targets and granted final effects must be selected
  before those actions. Recommendation: include source merge, release, default-branch adoption
  and reconciliation of the designated active worktrees in one eventual delivery grant. Keep
  private target locators in the private commissioning record. Do not perform successive redundant
  upgrades merely to repair a label; use the intended accepted version at action time.
- OQ-3 — BLOCKER: no for generic guidance; unresolved before promising a specific consumer's merge
  latency. The relevant owners must inspect their actual CI triggers and required checks before
  increasing publication cadence or proposing scoped changes. This discovery does not authorize
  consumer CI changes or a general repository inventory.
- A-1: Main is the ordinary shared baseline in the triggering case. Generic wording must support a
  repository's explicitly chosen default/integration branch.
- A-2: Existing branch, PR, plan and execution records are sufficient; no new runtime capability,
  persistent registry, scheduling service or approval parser is expected.

## Risks, validation and source ownership

- R-1: More Git events can amplify expensive CI. Detect by inspecting actual triggers for affected
  delivery targets; mitigate with coherent updates and scoped owner changes, preserving enforced
  checks. Do not claim a universal speedup from producer timings.
- R-2: Pressure to merge can introduce unfinished behavior. Require independently usable changes,
  meaningful review, relevant checks and explicit treatment of integration dependencies.
- R-3: Publication can be confused with approval or expose private context. Keep lifecycle and
  authority distinct, verify the destination/audience, and preserve explicit local-only work.
- R-4: Default-branch or adoption reconciliation can damage concurrent work. Use isolated scoped
  changes, supported lifecycle commands, and actual branch-state inspection; no product-branch
  merge or destructive cleanup is implied by a Kit update.
- R-5: The clarification itself can become ceremony. Validate realistic agent choices in the six
  scenarios above, not a new wording-only test for every sentence. Keep one coherent source review
  with same-reviewer corrections at the meaningful implementation boundary.

This discovery changes producer planning documentation only. Its local checks are Markdown,
public-core, diff hygiene and canonical plan validation/listing. No hosted suite is needed to
consolidate the discussion, and no independent implementation review is claimed for this draft.

The later managed change is routing-bearing and recipe-bearing. Candidate owner surfaces are the
existing planning lifecycle reference; discovery, drafting, execution and routine-change runbooks;
agent-task coordination; and installation/adoption guidance. Use the producer's operation-routing,
runbook-authoring, operational-recipe-maturity and runbook-to-script-promotion references. This is
guidance refinement using ordinary Git, not a reason to introduce a script or new Git abstraction.

The expected managed impact is instruction-only if external-action authority stays unchanged.
Expanding that authority instead requires explicit security/safety-policy classification and
review; a Markdown-only diff does not decide this. The eventual implementation must synchronize
declared resource mirrors, verify useful route behavior and use the applicable release checks.

## Next decision and handoff boundary

Review this consolidated draft and accept or revise the proposed scope in D-2 through D-6. There
is no unresolved technical blocker requiring another broad discovery. The next artifact is a
bounded implementation capability handoff and implementation plan using this document. It should
define coherent delivery outcomes,
affected guidance/tests, applicable validation and a truthful main-branch adoption finish line.

This is not yet implementation-handoff-ready, a new overarching program, or authorization to
resume consumer product implementations. Those owner-plan amendments remain subsequent work.

# Revision Notes

- 2026-09-09: Consolidated the discussion of commits, pushes, timely main integration, planning
  visibility, low-ceremony operation and the working-copy/default-branch adoption gap. Separated
  existing rules, user requirements, proposed defaults and action-time authority boundaries.
