Last updated: 2026-09-09T14:25:47Z (UTC)
Created: 2026-09-09
Status: active
Execution log: git-delivery-discipline_execution_log.md

# Predictable Git Delivery and Main-Branch Adoption

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.git-delivery-discipline
  kind: implementation
  purpose: Deliver lightweight Git discipline in existing agent workflows and verify the released guidance on designated consumer default branches and active working copies.
  first_cataloged: 2026-09-09T13:49:45Z
  catalog_metadata_updated: 2026-09-09T13:58:03Z
  relations:
    - kind: related
      target: codeheart-operating-kit.discovery.git-delivery-discipline
```
<!-- END CODEHEART PLAN METADATA -->

The accepted planning checkpoint is `e2bdc238ff8b0076f28b547d9d7c65ba3e6c34ff`.
Independent planning review found no material issues. On 2026-09-09 the user commissioned the
complete EP-01 through EP-03 delivery, with one dedicated Operating Kit implementer on
`codex/git-delivery-discipline-delivery`. The commissioning Organizational Operating Program
Director accepts the coherent source/candidate after independent review and applicable evidence.

The grant covers necessary supported tooling, implementation and focused checks, coherent commits,
normal pushes and PR work, applicable candidate/public smoke, normal source merge and next unused
patch publication under the existing unsigned internal/prototype HTTPS-plus-SHA256 audience,
Kit-only adoption PR merges on three named consumer main branches, and preservation-safe local
reconciliation in four named owner worktrees. Exact private targets and the direct report-back
locator remain in the commissioning task. Held product history is not authorized for publication.
Material scope, authority, preservation, safety or CI-cost changes return to that director.

Completion requires actual released guidance, all three main adoptions, all four reconciliations,
and truthful completed plan/log/index on producer main. After those effects, the grant permits a
normal documentation-only fast-forward closure push limited to this plan, its sibling log and
`docs/repo/plans/README.md` if policy permits; otherwise one normal closure PR. Force pushes,
provider/product work, consumer CI changes, security/billing changes and cleanup remain excluded.
This activation is execution authority, not a second planning approval gate.

| Essential context | Why it matters |
| --- | --- |
| `AGENTS.md` | Public producer authority and protected existing work. |
| `docs/repo/plans/git-delivery-discipline/git-delivery-discipline_discovery_doc.md` | Accepted intent, D-1 through D-6, requirements, two capability scopes and scenarios. |
| `docs/repo/runbooks/change-operating-kit.md` | Producer source route, impact classification and focused feedback. |
| `docs/repo/runbooks/release-operating-kit.md` | Existing guidance-candidate eligibility, native evidence, publication and retained-evidence rules. |
| `docs/repo/plans/release-validation-effort/release-validation-effort_execution_log.md` | Accepted source anchor and completed producer validation improvement. |
| `components/planning-workflows/managed/reference/planning-document-lifecycle.md` | Canonical existing Git/lifecycle policy to refine. |
| `components/agent-interface/managed/reference/agent-task-coordination.md` | Whole-plan assignment, report-back and covered final effects. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Audience, executable guidance and authority preservation. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Task/owner route before selecting tools and controlled effects. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Keep this at structured guidance with ordinary Git, without unnecessary wrappers. |
| `components/agent-interface/managed/runbooks/maintain-operating-kit-installation.md` | Supported lifecycle route and consumer preservation. |

Contents: [Foundation](#section-1---foundation), [Strategy](#section-2---strategy),
[Execution](#section-3---execution-plan), [Follow-ups](#section-4---future-planning).

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

An agent using the delivered Kit commits coherent recoverable work before formal epic acceptance,
publishes useful progress under established authority, and integrates independently usable
accepted outcomes without waiting unnecessarily for the whole plan. It checks relevant divergence
and identifies a concrete exception when ready work remains unmerged. Planning publication makes
the canonical document discoverable without changing execution authority.

Completion requires verified managed guidance, one coherent source review, applicable release
evidence, the released version merged on every assigned consumer default branch, healthy guidance
in the designated active worktrees, and truthful completed records on producer main. A source
merge, public release or working-copy-only upgrade is an intermediate state. No improvement in
consumer CI latency or general merge-conflict rate is claimed without evidence.

## 1.2 Project And Problem Context

The prior release established useful commit/publication rules and proportionate validation. Its
adoption boundary covered active worktrees while main remained older. Timing phrases such as
"useful completed checkpoint" and end-of-epic Git steps also leave room for delayed recovery and
integration. This plan refines those existing workflows; it does not create a new process program.

The accepted solution preserves authorization boundaries, substantive reasoning, current support
and meaningful review. Early publication improves recovery and visibility; integration still
requires usable behavior and applicable checks. One plan may contain several integration outcomes,
but neither every edit nor every epic mandates a new PR, test suite or approval round.

## 1.3 Current State Analysis

Planning branch: `codex/git-delivery-discipline-discovery`, based on producer main
`bbfce9b54b8e54174a58518f155be2874f310390`. Current source identity is v0.1.32. The accepted broad
source anchor `c36061f69f18b8a9aff2e018db2d5e35c6188ec8` is its ancestor. Source review and combined
broad/corrected-native evidence are recorded in the release-validation-effort log.

Existing producer CI separates cheap feedback, broad/guidance candidates and explicit public
smoke. The proposed managed edits use existing declarations, paths and ownership. No runtime,
schema, validator, workflow or installer behavior change is expected. Existing routing/resource
tests plus realistic agent walkthroughs can validate this instruction-only behavior without a
new sentence-presence test suite. If necessary executable changes emerge, classify them and
reassess candidate eligibility before scheduling broader work; do not weaken the guard to fit.

Consumer CI is independently owned and was not changed by the producer improvement. Current
targets have different triggers; inspect the actual workflow/rules at their action-time refs.
A long configured timeout is not a measured duration. Preliminary read-only inspection does not
authorize a run, a CI rewrite or a promise that every target merge will be inexpensive.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

```text
components/planning-workflows/managed/
  reference/planning-document-lifecycle.md          # modify canonical Git policy
  runbooks/discovery-workflow.md                    # modify publication/handoff hook
  runbooks/draft-implementation-plan.md             # modify integration/authority planning
  runbooks/execute-implementation-plan.md           # modify inner execution loop
  runbooks/handle-routine-change.md                 # modify equivalent routine behavior
  runbooks/review-planning-document.md              # modify readiness/coverage checks
components/agent-interface/managed/
  reference/agent-task-coordination.md              # modify assignment and delivery states
  reference/operation-routing-and-dispatch.md       # align affected route evidence/limits
  runbooks/maintain-operating-kit-installation.md    # modify repository rollout boundary
src/codeheart_operating_kit/resources/components/   # synchronize declared source mirrors
manifest.yaml, components/*/component.yaml,
profiles/standard.yaml and their resource mirrors  # regenerate normal release identities
pyproject.toml, internal/version/version.go,
src/codeheart_operating_kit/__init__.py,
install.sh, install.ps1, bootstrap.md               # existing literal version updates only
release-notes.md                                   # describe workflow change and adoption
docs/repo/plans/git-delivery-discipline/
  git-delivery-discipline_discovery_doc.md           # accepted discovery and capability
  git-delivery-discipline_implementation_doc.md      # this plan and final lifecycle
  git-delivery-discipline_execution_log.md           # create on activation; existing evidence home
docs/repo/plans/README.md                           # maintain discoverability/lifecycle
```

No new managed file, template, component, Git helper, evidence schema or CI service is planned.
Use one canonical policy with concise hooks in the named runbooks. All changed runbooks remain
agent-facing and retain compact intent/success/judgment/stop blocks. Coordination and installation
routes must distinguish scope, authority, preflight, actual effect and handoff. Reusable policy
stays in the Kit; concrete consumer commands, protected-branch requirements and private target
identities stay with their owners.

Recipe maturity is L1 structured guidance, using ordinary Git and supported lifecycle commands.
Validation is focused existing tests plus fresh-agent scenario review and actual release/adoption
evidence. Record material results/blockers briefly in the existing log/PR; no new output protocol.
No script promotion is warranted: the decision requires contextual judgment and existing tools
already implement the mechanics.

## 2.2 Open Questions And Assumptions Requiring Clarification

- OQ-1 — BLOCKER: no. Affects EP-01 through EP-03: publication authority stays as currently defined;
  the one assignment must explicitly include its bounded Git and final effects. Commissioning is recorded above; the plan is active. It does not infer authority from a broad user goal or a passing check.
- OQ-2 — BLOCKER: no. Affects EP-02: choose the next unused patch and verify the published old-CLI
  upgrade input at action time. Do not reuse the isolated demonstration's unpublished test identity
  by assumption. Preferred source anchor is the accepted reachable `c36061f` while applicable.
- OQ-3 — BLOCKER: no. Affects EP-03: the private commissioning record must bind the three consumer
  repositories/default branches, four worktrees, acceptance owner and authorized ancestry before
  effects. Verify current repository state, policies and target CI before publication/merge.
- OQ-4 — BLOCKER: no known current blocker; live target checks can create one. Affects EP-03: if
  an enforced check or preservation conflict prevents proportionate adoption, prepare the exact
  issue and return it to the owning director before that target's merge. Other independent target
  preparation may continue. Consumer CI/product changes are not silently added to this plan.
- A-1: Existing declarations and source runtime remain unchanged in behavior. Unknown or
  consequential changes invalidate the guidance-only assumption and go to the acceptance owner.

## 2.3 Architectural Decisions With Reasoning

**Git actions inside work.** Add a compact action table and execution hooks: commit coherent
progress; push for useful recovery/coordination under existing authority; inspect relevant main
divergence at meaningful starts/resumes/integration; merge usable accepted outcomes; give a short
reason/owner/resolution event for an exception. Preserve incomplete work safely before a planned
pause without falsely marking acceptance. No per-message trigger, age quota, forced rebase or
automatic merge is introduced. Update the per-epic flow so its end checkpoint cannot be read as
the first permitted commit or push.

**Plans and shared visibility.** Draft authoring includes local commits. A shared planning
assignment names its publication/merge authority once; authorized publication and handoff need
not wait for implementation approval. A pushed canonical branch becomes visible to configured
remote discovery after a successful refresh; main merge establishes shared placement, not
automatic activation. A local-only exception states its visibility limit. Revise one-PR wording
to permit several independently useful outcomes without requiring fragmentation.

**Main-branch adoption.** Prepare Kit-only changes from each target's current main, use the
supported upgrade route, verify and merge them, then reconcile the same released guidance into
active worktrees. Never merge a paused feature branch just to distribute the Kit. Existing
working-copy adoption commits may inform or supply a clean compatible change, but the actual
target baseline controls the supported upgrade and preservation checks.

Publishing a worktree branch publishes its ancestors too. A held branch containing unapproved
product history may therefore retain a local adoption commit even after main adoption succeeds.
Record that explicit exception and the owning product resumption trigger. Main rollout is still
complete; the local exception must not become an accidental authorization to publish product work.

**Validation and review.** Run existing route/resource checks during one coherent source batch;
retain content/manifest/integrity checks. One independent source reviewer challenges all six
discovery scenarios; use the same reviewer for corrections. A separate fresh agent receives the
deepest realistic vague request with no conversation context and must discover the owner route
before choosing Git/provider operations. Record decisions, not a transcript. No per-commit review
or live consumer push is needed for the probe. Producer publication uses the existing eligible
guidance route and native checks; consumer publication uses its own applicable gates.

**Authority and closure.** At commissioning, name one implementer and director for the whole
delivery. Proposed covered effects: necessary supported local tooling; source edits and focused
checks; coherent commits, normal branch pushes and PR work; applicable candidate/public smoke;
normal source merge, tag/release under the existing unsigned internal/prototype audience; scoped
consumer adoption PRs/merges; and bounded reconciliation/local commits in assigned worktrees.
Source acceptance follows the coherent review/candidate evidence. The implementer then completes
covered release/adoption without a new user gate per repository and reports final evidence.
Material scope, authority, safety or CI-cost changes return to the director.

Bind completion bookkeeping in the same grant: after actual release and target adoption, permit
a documentation-only fast-forward closure commit to producer main limited to this plan, its log
and `docs/repo/plans/README.md`, if repository policy allows it. Otherwise use one ordinary closure
PR. This is an explicit effect, not a generic permission to push directly to main. Verify final
remote state; never pre-mark unperformed merges/adoption complete or leave required closure local.

Excluded effects: provider/cloud operations, held product implementation, billing or security
setting changes, consumer CI rewrites, force pushes, deletion and wholesale branch publication.
Missing local Git/CLI/Python uses the existing tooling-readiness route; GitHub access uses the
owning service preflight. No new package-manager/auth recipe or runtime repair is added.

# Section 3 - Execution Plan

Current delivery: EP-01 and EP-02 complete; v0.1.33 is released and publicly verified. EP-03 is
partial: two of three target main branches and three of four worktrees are verified. The remaining
target is held by an enforced owner restriction outside this grant; its isolated adoption draft
is prepared. Do not close the plan before that adoption, reconciliation and remote closure.

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Coherent Git workflow guidance is implemented and reviewed | M | Accepted discovery; explicit execution grant |
| EP-02 | Accepted guidance is merged, released and publicly verified | M | EP-01 source and relevant review evidence |
| EP-03 | Target main branches and working copies use the release; delivery is closed | M | EP-02 published release; target preflight |

## EP-01

### A) Epic ID, Title, And Outcome

EP-01 — Git recovery, visibility and usable integration are explicit inside existing workflows.

### B) Scope

The nine named managed source documents, their mirrors, necessary normal release identity/notes,
and plan/log evidence. Preserve D-1 through D-5 and all discovery FR/NFRs.

### C) Files Touched

Section 2.1 managed source/mirror paths, existing release-identity surfaces and plan/log/index.
Existing tests are run unchanged unless a demonstrated contract gap requires a reviewed scope
correction; executable changes are not hidden merely to keep a lighter candidate classification.

### D) Acceptance Criteria And Size

Size M. Each discovery requirement maps to an actual decision/action in the updated workflow.
Independent review and the low-context probe show coherent commits before epic review, bounded
publication, accurate draft visibility, ready-outcome integration, protection of held ancestry,
and default-branch adoption without new ceremony. Existing route/resource and identity checks pass.

### E) Dependencies And Critical-Path Notes

Activation creates the sibling log and publishes its bounded planning checkpoint under the grant.
The implementer resolves substantive wording within accepted scope. Prepare the intended release
identity before final coherent source review so review and EP-02 candidate share relevant inputs.

### F) Tasks Checklist

- [x] Record the one grant, targets, branch and finish line; activate/commit/push the plan checkpoint.
- [x] Refine canonical Git policy and the existing authoring/execution/review/coordination hooks.
- [x] Add default-branch adoption and active-worktree reconciliation to the installation route,
      including the distinction between local health, published branches and repository rollout.
- [x] Synchronize declared resources and prepare current unused release identity/notes; inspect
      source/resource equality and run focused routing/resource, document and identity checks.
- [x] Complete one independent coherent source review plus a fresh low-context route probe across
      the discovery scenarios; resolve findings with the same reviewer and record retained evidence.
- [x] Commit and publish the coherent source PR and report its exact candidate and impact to the
      director for the planned source/candidate acceptance sequence.

### G) Implementation Notes

Use one source PR for this coherent guidance change. That bounded choice does not become a
universal one-PR-per-plan rule. Current source tests include real activation publication fixtures;
retain their authority boundary. Do not add wording-only assertions or a new Git command.

### H) Open Questions

OQ-1 and A-1. Material changes to authority, test/runtime behavior or new rule machinery return to
the owner before dependent execution.

## EP-02

### A) Epic ID, Title, And Outcome

EP-02 — The exact accepted guidance is on producer main and in a verified published Kit release.

### B) Scope

Existing candidate, artifact, merge and public-release routes. Preferred path is a mechanically
eligible guidance candidate plus semantic review, with applicable accepted broad evidence retained.

### C) Files Touched

Release identity/notes already prepared in EP-01, source PR and existing log. Packaging outputs
remain generated artifacts. No new CI workflow or guard is introduced.

### D) Acceptance Criteria And Size

Size M. Exact cumulative eligibility and identity checks pass; both native guidance lanes prove
fresh and upgraded managed content, integrity, reproducibility and authored preservation. Required
review/checks pass before source merge/tag. Published-tag public smoke passes on both supported
platforms. Exact source/tag/assets/audience and retained evidence are recorded.

### E) Dependencies And Critical-Path Notes

Use accepted reachable anchor `c36061f69f18b8a9aff2e018db2d5e35c6188ec8` only while its evidence
remains applicable. Use a separately verified published older tag for upgrade input. Candidate
validation precedes release publication; public smoke follows it. This is an actual release,
unlike the previous isolated demonstration. Do not create an extra demonstration release.

### F) Tasks Checklist

- [x] Inspect the exact cumulative diff and semantic impact, run existing guidance eligibility and
      release-identity preflight, and select the applicable candidate route with retained evidence.
- [x] Run the two native guidance lanes once for the coherent candidate; retrieve and verify the
      accepted packs/catalog. Retry only invalidated lanes and preserve meaningful failure output.
- [x] Obtain the planned director source acceptance, verify current head/checks and merge normally;
      compare relevant integrated inputs before tagging, rebuilding only if they changed.
- [x] Publish the exact accepted patch and artifacts under the covered audience, then run only the
      explicit released-tag public smoke and verify retrieval/version/identity.

### G) Implementation Notes

A guard rejection or material applicability concern is not permission to bypass or rewrite
eligibility. Resolve the concrete diff/classification; return an actual broad-validation scope or
cost change to the director. Do not replay broad suites because of a commit, log edit or resume.
There is no fixed timing promise and no automatic stall cancellation based on historical duration.

### H) Open Questions

OQ-2 and A-1. Tag/version/toolchain selection is action-time input, not a new architectural choice.

## EP-03

### A) Epic ID, Title, And Outcome

EP-03 — The release is the verified shared baseline on all assigned consumer main branches.

### B) Scope

Three designated repositories, four designated active working copies, and final producer records.
Integrate only the Kit adoption changes; do not resume the held owner plans or publish their work.

### C) Files Touched

Only supported lifecycle-generated Kit managed files/lock/managed AGENTS sections in each target,
with an adoption-only commit/PR and preservation evidence; producer plan/log/index for closure.
Repository-owned plans/configuration, module state and product changes are preserved.

### D) Acceptance Criteria And Size

Size M. Each target main contains the intended released version and correct managed routes, with
its actual required gates satisfied and authored content preserved. All four assigned working
copies are healthy on that version. No held ancestry is accidentally published. Main rollout and
any legitimate local-only worktree exception are distinguished. Final producer records are
completed and remotely visible after the actual effects, not just prepared for them.

### E) Dependencies And Critical-Path Notes

Bind target identities at commissioning and inspect live main/policies/CI before each adoption
publication. Prepare independent targets concurrently where safe. An expensive enforced target
check or conflicting owned change goes to that owner; do not silently alter its CI. Newly
released guidance is applied directly to the intended baseline, avoiding redundant intermediate
upgrades merely to reproduce the earlier worktree-only adoption.

### F) Tasks Checklist

- [x] Confirm current main refs, target access, relevant required checks and actual CI triggers;
      record justified checks and any material cost/overlap issue before the corresponding effect.
- [x] Prepare isolated Kit-only adoption branches from target main. Check installed state and use
      the supported released upgrade with preview/apply; verify changed routes, installed health
      and preservation against each target's own baseline.
- [ ] Open/update and merge the scoped target PRs under the one grant once their applicable gates
      pass. Verify the actual remote-main lock/content and healthy checkout after integration.
- [ ] Reconcile the released main guidance into each designated active worktree while preserving
      product/plan work. Commit locally; publish only when the whole branch ancestry is authorized.
      Record any held-ancestry local exception and its later owner trigger.
- [ ] Report source/release/target-main/worktree evidence to the director, close the plan/log/index
      after actual completion through the covered closure route, and verify final remote state.

### G) Implementation Notes

A newer CLI can classify a pristine older installation using its newer graph. If encountered,
verify the old tree with the matching published old CLI and use that supported forward-upgrade
route, as in the prior delivery; do not hand-edit locks or apply unrequested runtime repairs.
Any real drift/recovery blocker preserves state and returns through the lifecycle owner route.

Existing valid producer source review is reused for mechanical target adoption; only required
repository review and target-specific preservation/health evidence are added. No product-wide
suite is justified solely by a Kit change, but existing enforced checks remain binding until
their owner changes them. Completion does not imply that consumer CI has been redesigned.

### H) Open Questions

OQ-3 and OQ-4. Missing exact target authority or a real failed gate blocks only its dependent
effects; keep independent authorized preparation moving and report the concrete boundary.

# Section 4 - Future Planning

After the completed main-branch rollout, amend the held consumer product plans separately under
the new process and establish any actual dependencies before resuming implementation. Runtime
profiling, older-installation classification improvements and consumer CI redesign remain distinct
owner work. This plan creates no fleet audit, goal, watcher, approval service or reporting program.

# Revision Notes

- 2026-09-09: Drafted from the user's accepted discovery direction and capability handoff. Planned
  source guidance, proportionate release and default-branch adoption as one bounded delivery,
  preserving inactive lifecycle and the separate commissioning authority for external effects.

- 2026-09-09: Activated under the complete delivery grant; bound acceptance, rollout and closure
  authority while retaining private target locators only in the commissioning task.
