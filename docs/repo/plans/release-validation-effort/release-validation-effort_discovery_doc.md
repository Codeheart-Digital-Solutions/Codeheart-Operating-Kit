Last updated: 2026-09-08T13:53:53Z (UTC)
Created: 2026-09-08
Status: completed
Completed: 2026-09-08

# Release Validation Effort Discovery and Proposal

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.discovery.release-validation-effort
  kind: discovery
  purpose: Reduce recurring release validation effort through an enforced guidance candidate route and focused release identity checks without weakening integrity or platform preservation.
  first_cataloged: 2026-09-08T13:48:34Z
  catalog_metadata_updated: 2026-09-08T13:53:53Z
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.development-feedback-and-validation
```
<!-- END CODEHEART PLAN METADATA -->

## Readiness and authority

New-request discovery; readiness: manual-review-ready proposal for the commissioning acceptance
owner. The proposed capability and acceptance below are concrete inputs to implementation
planning, conditional on that owner's acceptance. They are not active implementation epics.
Discovery, read-only evidence gathering, document authoring and a coherent local commit are
authorized. Workflow execution, source changes, publication and adoption await the chosen
delivery grant. No hosted diagnostic or full candidate was dispatched for discovery.

Producer baseline: `origin/main` at `0e7801598c9b753e5c39df36f56f4da4c112a50a`, fetched and
verified on 2026-09-08; isolated branch `codex/release-validation-effort` began clean.
Repository catalog mode is canonical. Tracked producer source is authority; the ignored
consumer installation is not used as doctrine or installed to perform this work.

This is routing-bearing and recipe-bearing work: it proposes changes to the maintainer release
route and its executable workflow. Apply source references
`components/agent-interface/managed/reference/operation-routing-and-dispatch.md`,
`operational-recipe-maturity.md`, `runbook-to-script-promotion-standard.md` and
`runbook-authoring-standard.md` in that same directory.

## Problem and intention

Make a coherent Kit release substantially less effortful by removing checks that do not test
the changed behavior, catching cheap release-identity defects early, and avoiding speculative
cancel/retry cycles. Preserve technical understanding, intention elicitation, substantive plans,
one combined independent review with same-reviewer corrections, and director acceptance of
outcomes without a second technical review.

Priority: truthful acceptance and preservation first; reduced recurring elapsed time and human
coordination second; minimal machinery and maintenance third. Meaningful guidance and safety
policy review still matters even when no Go algorithm changes. A file suffix is not impact.

Success is an enforceable guidance release route with fresh package and native lifecycle proof,
clear escalation for consequential or unknown changes, and measured candidate-step outcomes.
No fixed speedup is promised from historical timings. Do not restrict context ingestion, replace
intention questions with a form, or create a new review/task/release for each correction.

Non-goals: content distribution redesign, generic dependency classifier or proof cache, auth or
billing changes, support/audience changes, company-wide process inventory, migration retirement,
consumer adoption, or unmeasured plancatalog production optimization in this first delivery.

## Evidence and obligations

The preceding delivery is complete. [PR #17](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/pull/17)
merged at the baseline and [v0.1.32](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v0.1.32)
was published. The final accepted plan/log record is at
`f29172b3b591e3329dc8f3664f4a7a989d575a3a` on `codex/development-feedback-delivery`, under
`docs/repo/plans/development-feedback-and-validation/`. Read it with `git show`; main's earlier
status does not reopen that delivery. Public smoke run
[34172479978](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/34172479978)
and the closure record establish completed release/adoption. Private target records stay private.

### Measured observations

| Observation | Evidence and interpretation |
| --- | --- |
| Ordinary feedback already became small | Prior delivery reported roughly 30–38 seconds; branch pushes no longer duplicate PR feedback. Preserve this work. |
| Windows source suite dominated | [Job 101889250541](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/34170316992/job/101889250541) succeeded. API step timestamps independently rechecked: Go 23:32:12–00:06:51 UTC, 34m39s; builder 6s; staged install 3s; upgrade 5s. These warmed-job timings are not a cold lightweight-lane forecast. |
| Large package disparity | Same job log independently rechecked: plancatalog 1079.607s, portfolio 288.422s, commands 179.521s. Commissioning evidence gives Ubuntu 52.571s, 29.393s, 12.764s respectively, full Go about 78s; macOS candidate about 7m. Package elapsed times do not add directly to wall time. |
| First failures were machinery defects | Closure record and `bec52c4`: Windows test harness chose unavailable WSL Bash; patch identity change left `internal/hash/hash_test.go` profile digest stale. Focused correction retained the same reviewer. Neither was a new managed-guidance defect. |
| A cancellation was premature | Closure record corrects run 34168210016 cancellation after 35m of silent builder output. Earlier successful Windows run 33998404410/job 101392834844 took 38m42s for Go, plancatalog 1055.871s. Silence and elapsed time did not prove a stall. |
| Source/package separation already exists | `.github/workflows/validate.yml` runs visible native Go steps. `scripts/build-release-assets.py` packages twice and verifies; it no longer runs source suites internally. `candidate_lane` retains applicable lane results; explicit released-smoke does not replay source suites. |

The commissioning investigation found no billing-limit error and established a public repository
using standard hosted runners. This proposal does not attribute failures to budget or claim an
inherent Windows cause. No billing investigation or account change is warranted by this evidence.

### Source-backed hypotheses, not measured causes

`resources.go` embeds components, profiles, templates, schemas and manifest identity in Go
executables. Guidance changes need fresh binaries and packages under the current design; they
do not inherently change branch-proof algorithms.

`internal/plancatalog/inventory.go` loads snapshots, inspects branch candidates and launches Git
processes. `proposeBranchCandidateEvidence` computes initial evidence, searches up to 512 history
candidates and calls `RecomputeBranchCandidateEvidence` for each. `branch_reconciliation.go`
rechecks repository/ref/commit/merge-base/path/attribute facts. Migration binding/checkpoint paths
verify persisted historical evidence. `ep02_test.go:migrationRepository` copies a fixture, creates
real Git history and adopts mixed mode; additional tests create more history and inventories.
The commissioning source count (119 tests, 9 files, no `t.Parallel`, 77 fixture-helper call sites,
326 Git-helper call sites) is static, not runtime profiling.

Repeated serial Git/process/filesystem work is a plausible cost center. It does not prove which
test, process-start mechanism or filesystem behavior caused Windows disparity. The explicit
Ubuntu 100k-path classifier benchmark is unrelated evidence, not a Windows workload count.

### Which checks still earn their cost

| Obligation | Guidance release recommendation | Why |
| --- | --- | --- |
| Content/routing, public-core, schema and release identity | Fresh focused checks and meaningful review | Guidance can change operational behavior; stale mirrors and identity are real defects. |
| Both distributable packs, twice-built equality and catalog-to-binary verification | Fresh for every new payload | Embedded bytes and version change; old artifacts cannot qualify new ones. |
| Native fresh install, materialization, check, upgrade and preservation | Fresh on macOS and real Windows | Detect actual new-content integration and native handoff failures. |
| Existing historical 0.1.25/0.1.26 upgrade smoke and Windows reparse smoke | Retain initially | Existing inexpensive steps already exercise lock compatibility and containment; dropping them saves little and changes assurances. |
| Unrelated broad Go, full Python parity/plan-history compatibility, exact oldest-Git proof, scale benchmark | Retain exact applicable prior evidence for guidance route | These cover unchanged mechanisms, not newly edited instructions. Consequential/unknown changes still require broader evidence. |
| Source review and director acceptance | One coherent technical review, then outcome acceptance | Neither repeated reviewers nor duplicated technical review improves this boundary. |
| Public release download smoke, signing/audience and authority | Unchanged, only for authorized publication | Candidate acceptance and publication remain distinct. |

## Requirements

- FR-1: A maintainer can explicitly choose a guidance candidate against an exact accepted broad
  source baseline, while broad candidate remains the default and fallback. A separate verified
  published tag supplies the matching old CLI for current-release upgrade smoke.
- FR-2: The workflow rejects an ineligible guidance request before expensive jobs. Missing or
  unresolvable baseline, runtime/dependency/toolchain/installer/schema/ownership changes, unknown
  paths, and unrecognized release-identity edits must not silently qualify.
- FR-3: Candidate evidence records exact baseline, resolved source commit, reviewed semantic
  impact, fresh lanes, retained run/commit inputs, artifact identities and remaining gates in
  the existing execution log and workflow summary. No external cache or new evidence database.
- FR-4: All new payloads get full content/resource/version/integrity verification, reproducible
  packs, native materialization and upgrade preservation. Prove changed source bytes at their
  installed targets, including added managed files; verify declarations and mirrors.
- FR-5: Cheap release-identity checks catch normal patch inconsistencies before native dispatch.
  Remove current-version fixture churn only where it is incidental; retain independent golden
  vectors and corrupted/mismatched identity rejection.
- NFR-1: Preserve supported platforms, ownership, rollback, signing/audience and action authority.
- NFR-2: Unknown semantics escalate. A reviewed safety-policy Markdown change needs explicit
  policy review and affected scenario probes; narrow mechanical eligibility cannot waive it.
- NFR-3: Keep ordinary feedback small; retain visible package output and applicable evidence.
  Do not infer stalls from uncalibrated duration or impose a new whole-command timeout.

## Decision ledger

All proposals below await commissioning acceptance; recommendation is not execution authority.

| ID | Decision / class | Dependencies | Recommended outcome and closure |
| --- | --- | --- | --- |
| D-1 | Release assurance route / operating model | FR-1–4, NFR-1–3 | Add explicit guarded guidance route; owner accepts matrix and retention boundary. |
| D-2 | Enforcement shape / recipe | D-1 | Small producer-owned eligibility guard plus existing workflow; owner accepts scope and default-broad behavior. |
| D-3 | Patch preparation / validation | D-1 | Focused identity preflight and incidental-fixture cleanup; preserve independent integrity assertions. |
| D-4 | Runtime bottleneck / scope | D-1 | Defer optimization; preserve a bounded diagnostic follow-up and no hosted run at discovery. |
| D-5 | Newer CLI on older installation / product boundary | None | Separate bounded lifecycle follow-up; do not broaden validation implementation. |

### D-1 — Two practical candidate routes

Recommend `candidate` with a default broad validation scope and an explicit guidance scope.
Keep `released-smoke` separate. Do not infer guidance eligibility automatically from `.md`.
Compared with keeping every historical test mandatory, this removes the observed dominant
unrelated cost. Compared with unbundling content distribution, it preserves existing release
contracts and avoids a new update/integrity architecture. Separate distribution is only worth
reopening if fresh packaging and smoke later dominate measured cost.

Meaningful change classes:

| Change | Minimum adequate evidence |
| --- | --- |
| Unshipped planning/log documentation only | Proportionate document/catalog checks; no release solely for bookkeeping. |
| Managed instruction/content change within existing ownership and placement | Reviewed semantic impact; guidance route obligations above and applicable prior broad evidence. |
| Safety or consequential operational policy in guidance | Explicit policy/owner review and focused scenario/low-context route probes in addition; broaden if changed obligations invoke uncertain executable dependencies. |
| Scaffold, component/profile selection, placement/ownership, schema, installer/runtime, dependencies or toolchain | Broad route initially, plus affected native success/failure/rollback, schema/migration or integrity tests under `change-operating-kit.md`. |
| Workflow, eligibility guard, build/release-contract or test infrastructure | Focused mechanism tests during iteration and one coherent broad acceptance for this delivery; new gate cannot certify itself solely through its lighter path. |
| Mixed or unknown changes | Broad route until dependencies are resolved; manual assertion alone cannot override mechanical ineligibility. |

Adding a managed file through a recognized declaration under an existing component target and
unchanged ownership is eligible content addition; it is not a new component, scaffold, ownership
mode or placement contract. Any other declaration/placement change follows broad scope.

Guidance candidate native smoke retains current historical upgrades and adds the named current
published release as the current-content upgrade case (deduplicate when already present). Use its verified
matching CLI to initiate normal upgrade. Verify dry-run and failed verification preserve state,
successful apply/check, all changed materialized resources and authored/config/AGENTS preservation.
Retain applicable transaction fault-injection evidence for unchanged runtime; rerun it when
mechanisms or dependencies change. Do not claim a missing-catalog smoke covers every rollback phase.

### D-2 — Small guard, explicit judgment

Use one small testable producer script in existing `scripts/`, called by the workflow's input
preflight and documented in the release runbook. It is a read-only L2 primitive: explicit source
and baseline revisions; concise eligibility/blocker output and nonzero rejection. This warrants
a script because structured identity diffs and path handling are fragile in workflow shell.
No generic classifier, CLI extension, service, nested framework or generic JSON envelope.

Compare the complete baseline-to-candidate diff, including additions, removals, renames and mode
changes, against a deliberately narrow known surface. Permit managed content at declared owned
targets and byte-equal mirrors, plus specifically parsed release bookkeeping. Explicitly admit
unshipped planning/log documents under `docs/repo/plans/` and their index, subject to existing
document/catalog checks; those records must not force a broad run for their own evidence.
Release-note prose is a known packaged input requiring fresh pack integrity and content review.
Changes to governing runbooks, guard/workflow, bootstrap or installers still use broad scope.
Do not allow whole
`components/`, `templates/`, `profiles/`, `manifest.yaml`, tests or Go files by directory/suffix.
Constrain component/profile/manifest changes to reviewed identity/content additions permitted by
existing target contracts; selection, ownership, dependencies and compatibility stay equal.
Constrain version defaults and package metadata to literal identity-field changes. Unrecognized
syntax or fields reject the guidance route. Whole arbitrary test edits do not qualify.

The named source baseline must be an exact ancestor commit with owner-accepted source and
applicable successful broad evidence, recorded in the existing log. A passing run alone is not
source acceptance. It may be unpublished; the separate upgrade tag must actually be published
and verified. Preserve the same broad-source anchor across successive guidance releases until a
new coherent broad acceptance replaces it, so cumulative executable drift cannot disappear.
Record both identities and retained run/commit evidence; no new registry is required. An operator's
semantic review establishes policy and operational dependency impact; the mechanical guard
only proves a bounded diff. Candidate preflight reports that distinction. Prior environment,
configuration or dependency changes invalidate affected results even if source paths match.

Reuse existing `candidate_lane` correction handling. A partial run cannot become complete just
because other jobs were skipped; record fresh and retained lane applicability to the exact
candidate. Existing summaries and execution log suffice. No universal proof-reuse system.

### D-3 — Catch identity churn cheaply

The profile hash fixture in `internal/hash/hash_test.go` and hardcoded release version in
`internal/manifest/manifest_test.go` are concrete starting points. Inspect whether each assertion
is an algorithm golden vector or incidental use of current release files before changing it.
For hashing, prefer a stable independent fixture/vector plus a separate current-profile
integration assertion; do not replace a known-good digest with a call to the function under test.
For versioning, compare authoritative release inputs while retaining an intentional mismatch
negative case. Existing `tests/test_release_assets.py` already has a cheap version mismatch test.
Keep fixed historical released fixtures fixed. Add only focused current identity checks to early
feedback if measured cheap; avoid pulling plancatalog integration into ordinary feedback.

### D-4 — Measure before optimizing runtime

The first delivery removes unrelated runtime repetition from eligible releases; it does not need
a measured Windows root cause to justify that scope. Keep the broad runtime suite for actual
runtime changes. Do not delete old migration guarantees merely because their tests are costly.

If runtime optimization is commissioned later, first gather per-test elapsed output from one
bounded plancatalog scenario family with `go test -json -count=1 -run <exact-test-expression>`.
Compare cold compilation separately from a repeat with warm build cache; count actual Git calls
and fixture setup cost for those tests before proposing fixture reuse, bounded parallelism or
request-scoped immutable Git-fact reuse. Each has different isolation or stale-proof risks.
No package-wide hosted profiling is requested here. Before any hosted diagnostic, return its
exact test names, expected duration supported by a local/historical sample, and the question it
will answer to the commissioning owner. Long silence alone is not a cancellation criterion.

### D-5 — Related lifecycle usability observation

The closure record reports a newer shared CLI classified pristine older installations as partial
because newer managed paths were absent. A verified official matching old CLI initiated normal
upgrade successfully. Recommend separate lifecycle discovery with a minimal old/new reproduction
and no repair/lock edit bypass. This matters to upgrade ergonomics but changes state classification
and ownership-sensitive runtime, so it should not quietly enter this guidance validation delivery.

## Implementation Capability Scope - Proportionate release validation

Capability: release eligible managed-guidance changes with fresh package/native lifecycle proof
without replaying unrelated runtime history suites.

Primary workflow: maintainer reviews semantic impact, names an accepted broad source baseline
and a verified published upgrade tag, dispatches
candidate with explicit scope, reviews fresh/retained evidence, then uses unchanged publication
and smoke procedures only under a delivery grant.

Must cover: D-1 matrix, D-2 fail-closed guard, D-3 identity checks, current baseline upgrade,
source-to-installed equality, safe partial reruns and exact artifact/source evidence.
Windows smoke must explicitly check each relevant native command's exit status as well as
resulting bytes and state. Existing PowerShell command sequencing alone is insufficient when
intermediate failures can be hidden by later successful commands.

Explicitly out of scope: D-4 production optimization, D-5 lifecycle repair, distribution redesign,
new platform/support/audience promises, proof caches and consumer changes.

Deferred or blocked: owner acceptance of D-1–3 and execution/final-effects grant are pending.
Preserve decisions: D-4 and D-5 stay bounded follow-ups unless separately commissioned.
Planner must not reinvent: broad default; meaningful semantic review; both native platforms;
fresh reproducible packs; integrity/ownership/rollback gates; unchanged publication authority.

Feature-level success evidence: positive eligible-content and negative runtime/schema/unknown
fixtures, dispatch boundary tests, both native guidance jobs with installed-byte and preservation
proof, broad fallback wiring, independent combined review, truthful recorded timing and evidence.

## Bounded implementation proposal

After acceptance, use one implementation plan and one coherent implementation/review batch.
First settle the guidance eligibility contract and matrix in producer
`docs/repo/runbooks/change-operating-kit.md`, `release-operating-kit.md` and, where required,
`docs/repo/reference/consumer-impact-classification.md`. These are maintainer-facing recipes;
retain compact intention, exact input/evidence/stop boundaries and existing tooling-readiness
route. Do not copy generic doctrine or alter unrelated managed runbooks.

Then implement the narrow guard, focused identity corrections and workflow selection together.
Keep existing native pack and smoke steps; condition broad Go/parity/history/oldest-Git/benchmark
work on broad scope and add focused content checks and named-baseline upgrade to guidance scope.
Scope includes `scripts/`, `.github/workflows/validate.yml`, their existing tests and only the
directly implicated identity fixtures. Update nearest script/index documentation if needed.
Broader architecture or newly discovered runtime defects return to the owner.

Validate cheap behavior during iteration: eligible managed edit and new file, wrong mirror,
manifest compatibility/ownership change, renamed/deleted/unknown path, executable/mode change,
version-only field versus hidden code edit, invalid/nonancestor baseline, missing retained
evidence, policy-review escalation, scope/lane combinations and explicit public tag handling.
Use actual dispatch-condition tests and guard fixtures, not string-presence tests alone.
Corrupt a pack/content identity and prove rejection; preserve independent hash/mismatch negatives.

At the coherent candidate boundary, run one broad candidate for the changed validation machinery
and obtain source acceptance using the combined review and result. Then run the new guidance
path on a controlled eligible content/version delta against that exact accepted broad-source
baseline, with the existing verified release as the upgrade input, on both native platforms,
to prove the branch actually executes without the
excluded source suites. This second run is focused package/smoke validation, not another broad
matrix. Keep an ineligible guidance dispatch as a cheap preflight-negative proof. Exact hosted
refs and effects belong in the accepted execution plan; none is dispatched during discovery.
This exercises the real eligible path without publishing a release to bootstrap the test.

Acceptance requires the combined review to close material findings; all required fresh or
retained evidence applicable; documented broad fallback; no published-tag confusion; recorded
step timing and residual risks; runbook and enforcement agreement. Publication/adoption are
optional later effects only if expressly included by the commissioning owner. An unreleased
source/workflow improvement must be reported as such.

## Risks, questions and handoff

- R-1: A too-wide allowlist misclassifies behavior. Mitigate with field-level identity checks,
  unknown rejection, semantic review and adversarial guard fixtures.
- R-2: Source-only diff misses environment drift. Record tooling/runner inputs and rerun affected
  lanes when applicability is uncertain; do not claim byte-identical code proves environment identity.
- R-3: Dynamic fixture expectations hide corruption. Preserve independent vectors and negatives.
- R-4: New guidance branch is only inspected, never executed. Require focused native branch proof.
- R-5: Process improvements become more machinery than saved effort. Keep one guard, existing
  workflow/log and one review; measure actual lane timings before further optimization.
- OQ-1 — BLOCKER: yes for execution, not proposal review. Commissioning owner must accept the
  capability/acceptance boundary and grant chosen effects once. Recommendation: accept D-1–3 as
  one bounded delivery; D-4–5 remain separate follow-ups.
- OQ-2 — BLOCKER: no. Exact slowest test and Windows mechanism remain unknown; not needed to
  remove unrelated regression suites from proven guidance-only releases.
- A-1: Existing platform/compatibility promises remain in force. This proposal changes validation
  selection, not support policy.

Manual review packet: approve or revise D-1–3 and the capability scope, confirm D-4–5 exclusions,
and define implementation plus any publication boundary. No extra intention question is needed:
the commissioning constraints already resolve the consequential scope and priority choices.
Normal implementation-plan drafting follows that acceptance. Coordination visibility currently
comes from direct handoff; no branch push or portfolio refresh is claimed.

## Review and revision notes

2026-09-08: Initial evidence-backed proposal authored from current producer source, historical
closure and live successful-job timing. One independent read-only reviewer assessed D-1–5.
The reviewer identified ambiguous eligibility for delivery documentation and unnecessary coupling
of source assurance to a published baseline. Both were accepted and corrected: narrow planning
documents are eligible, release notes remain reviewed package inputs, and an accepted broad-source
anchor is separate from the verified old-release upgrade input. The same anchor persists across
successive guidance releases. Additional clarifications cover managed-file additions under an
existing target and explicit Windows native exit-status assertions. D-3 integrity safeguards and
D-4–5 deferrals were supported. No public-safety or authority findings; timing/root-cause
uncertainties remain explicit. Same-reviewer focused follow-up closed every material finding;
the proposal is ready for commissioning acceptance. Review does not authorize implementation.

Local Markdown and public-core checks pass. Canonical `plans validate` passes (42 records) with
existing legacy-layout warnings outside this proposal; `plans list` includes the new draft.
Validation used the current catalog mode without changing config or discovery version. No source
or native suite is warranted for this proposal-only edit.

Commissioning acceptance on 2026-09-08: D-1–3 and capability scope approved for the sibling active
implementation plan. D-4–5 deferred. Finish line is producer workflow merged to main, with no
public version, tag/release or consumer adoption. Authority and evidence applicability remain
meaningful owner judgment in the existing log; no receipt registry or runner-image equality gate.
