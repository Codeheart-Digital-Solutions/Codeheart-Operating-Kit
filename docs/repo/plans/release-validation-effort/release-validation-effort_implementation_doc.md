Last updated: 2026-09-08T13:59:38Z (UTC)
Created: 2026-09-08
Status: active
Execution log: docs/repo/plans/release-validation-effort/release-validation-effort_execution_log.md

# Proportionate Producer Release Validation

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.release-validation-effort
  kind: implementation
  purpose: Implement and merge guarded guidance candidates with focused identity and native lifecycle validation.
  first_cataloged: 2026-09-08T13:59:38Z
  catalog_metadata_updated: 2026-09-08T13:59:38Z
  relations:
    - kind: related
      target: codeheart-operating-kit.discovery.release-validation-effort
```
<!-- END CODEHEART PLAN METADATA -->

The commissioning owner accepted discovery D-1–3 and its capability scope on 2026-09-08,
including activation and execution. D-4–5 remain separate follow-ups. This plan delivers one
producer-only PR merged to main. No new public version, tag, release or consumer adoption is
included. A controlled guidance candidate is validation evidence only.

| Essential source | Purpose |
| --- | --- |
| `AGENTS.md` | Public producer authority and preserved work. |
| `docs/repo/plans/release-validation-effort/release-validation-effort_discovery_doc.md` | Accepted evidence, capability and exclusions. |
| `docs/repo/runbooks/change-operating-kit.md` | Change gates and local tooling route. |
| `docs/repo/runbooks/release-operating-kit.md` | Existing native, integrity and publication contract. |
| `.github/workflows/validate.yml` | Candidate dispatch and native smoke to adapt. |
| `scripts/build-release-assets.py` | Reproducible packaging without broad source invocation. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Owner and execution-surface separation. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Maintainer recipe requirements. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Minimal justified mechanics. |
| `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` | Small tested script contract. |

Contents: [Foundation](#section-1---foundation), [Strategy](#section-2---strategy),
[Execution](#section-3---execution-plan), [Follow-ups](#section-4---future-planning).

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

An explicit eligible guidance candidate runs fresh content/identity checks, reproducible packs
and real macOS/Windows install/materialization/upgrade/preservation checks without unrelated broad
Go, parity, history, oldest-Git or benchmark suites. Unknown or consequential diffs cannot choose
that route. Default broad and explicit public smoke remain intact. Completion requires one
independent coherent source review, broad machinery acceptance, a real guidance-path exercise,
commissioning final acceptance and normal PR merge. Timing is observed, not promised.

## 1.2 Project And Problem Context

Windows broad Go took 34m39s in the preceding accepted delivery while subsequent packaging and
native smoke took seconds with warmed caches. The root cause is unmeasured. Embedded guidance
still requires fresh binaries; this does not require replaying unchanged branch-history tests.
The accepted design removes repetition before adding machinery and preserves meaningful review.

## 1.3 Current State Analysis

Existing workflow already separates feedback, candidate and released-smoke and supports partial
lane retries. Builders already package twice without source suites. Extend those surfaces rather
than replace them. Current-version hash and manifest test literals cause avoidable patch churn.
The target adds a narrow Git-diff eligibility guard, explicit guidance scope and upgrade tag,
focused identity checks and native byte/preservation assertions. No consumer content is changed
in the final producer branch; controlled demonstration content stays in a separate test branch.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

```text
.github/workflows/validate.yml                    # modify selection and native smoke
scripts/validate-guidance-candidate.py             # create read-only eligibility primitive
scripts/verify-guidance-lifecycle.py                # create shared isolated native proof
scripts/README.md                                 # create compact script contract index
tests/test_guidance_candidate.py                   # create eligibility regression tests
tests/test_packaging_resources.py                  # modify dispatch behavior tests
tests/test_release_assets.py                       # retain independent negatives
internal/hash/hash_test.go                         # modify incidental live-profile vector
internal/manifest/manifest_test.go                 # modify current-version consistency check
docs/repo/runbooks/{change,release}-operating-kit.md # modify governing route together
docs/repo/plans/release-validation-effort/          # plan, log and accepted discovery
```

One coherent source batch covers guard, tests, workflow and instructions. A second epic validates
and integrates it. Read-only script promotion is L2 at existing producer `scripts/`, called by
release runbook/workflow with explicit baseline/candidate inputs, concise output and nonzero
failure. No authority parser, receipt registry, AST/dependency engine or generic output framework.
Use existing Python/PyYAML tooling and Git; missing tooling uses the source readiness runbook.

## 2.2 Open Questions And Assumptions Requiring Clarification

OQ-1 — BLOCKER: no. Affects EP-02: exact accepted broad-source anchor is established after the
broad run and source review. Report that commit to the commissioning owner for evidence acceptance,
then exercise a guidance delta against it. This is the planned acceptance point, not new scope.

OQ-2 — BLOCKER: no. Affects EP-02: hosted duration is unknown; preserve useful lane evidence and
investigate failures without unsupported cancellation or a new timeout. Existing platform promises
remain unchanged. No runner-image byte-identity requirement is introduced.

## 2.3 Architectural Decisions With Reasoning

Add `candidate_scope` (broad default, guidance explicit), exact source `baseline_ref`, and verified
published `upgrade_version` for guidance smoke. The guard compares cumulative baseline-to-candidate
Git trees, rejects mode/removal/rename/unknown consequential changes and constrains known identity
edits to literal fields. Managed additions qualify only through existing target/ownership contracts;
source mirrors must match. Narrow planning/log/index and reviewed release notes are eligible.
Runtime, installers, schemas, workflow/guard and governing runbook changes require broad coverage.
The existing reviewed log owns source acceptance and semantic policy judgment. A green guard proves
only mechanical eligibility. Retain a broad-source anchor across successive guidance releases.

Both scopes keep packaging and native smoke. Guidance runs focused Go manifest/hash checks and
content/schema/routing/resource checks, materializes source bytes and uses a verified matching old
release CLI for current-release upgrade in addition to existing historical smoke. Explicitly test
Windows command exits, not merely later successful output. Preserve independent hash vectors,
manifest corruption and version mismatch negatives while removing incidental release literals.

The producer owns both maintainer-facing runbooks. Preserve intention/authority/evidence/stop and
tooling-readiness boundaries. The independent source reviewer also performs a fresh low-context
route probe: a vague guidance release request must reach the release runbook, distinguish semantic
review from eligibility and select the right evidence route before dispatch.

Authority covers necessary local setup, implementation, focused checks, coherent commits, normal
pushes, one PR with updates and planned hosted validation. Commissioning owner accepts source
anchor and final outcome; after final acceptance and passing applicable gates, this task merges
the PR normally. No force push, deletion, public release, consumer effects or billing/policy change.
Normal branch: `codex/release-validation-effort`; test evidence branch is isolated and retained.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Small enforced guidance candidate and aligned instructions | M | Accepted D-1–3 |
| EP-02 | Proven candidate paths and accepted merged producer change | M | EP-01 coherent source |

## EP-01

### A) Epic ID, Title, And Outcome

EP-01 — Guidance candidates are usable and cannot waive consequential changes mechanically.

### B) Scope

Eligibility guard, workflow scope, identity tests and native lifecycle checks with maintainer docs.

### C) Files Touched

Section 2.1 source/test/runbook paths and plan/log. No unrelated managed content or runtime changes.

### D) Acceptance Criteria And Size

Size M. Actual cumulative diffs and input combinations are tested; native commands explicitly fail;
changed resources are checked as installed bytes; independent integrity negatives remain; default
broad and released-smoke separation survive. One coherent independent review closes findings.

### E) Dependencies And Critical-Path Notes

Guard inputs and workflow selection form one interface. Develop bounded guard work alongside local
workflow integration; synchronize before focused tests and review.

### F) Tasks Checklist

- [ ] Implement narrow guard with positive managed edit/addition and adversarial identity/path tests.
- [ ] Integrate scope/preflight and focused content/native paths; preserve broad/public modes.
- [ ] Remove incidental identity fixture churn while keeping independent golden/negative checks.
- [ ] Align maintainer recipes and script contract; test fresh-agent route selection.
- [ ] Run focused checks and one independent combined source review; correct with same reviewer.

### G) Implementation Notes

No semantic authority automation. Fail closed on unfamiliar fields/paths. Track only relevant
toolchain/dependency/behavior invalidation; incidental host/log updates do not force broad repeats.

### H) Open Questions

OQ-1/2 affect later hosted evidence only; implementation is authorized now.

## EP-02

### A) Epic ID, Title, And Outcome

EP-02 — Both candidate paths are demonstrated and the reviewed producer change is merged.

### B) Scope

One PR, broad candidate, eligible isolated guidance delta on both platforms, negative preflight,
evidence acceptance and integration. No fabricated feature release or version change on main.

### C) Files Touched

Plan/log final state and focused review corrections; controlled test branch content/version only.

### D) Acceptance Criteria And Size

Size M. Applicable broad lanes pass jointly; owner accepts exact anchor; guidance scope really
executes native pack/materialization/upgrade proof without broad suites. Record exact refs/runs,
retained applicability and scoped timing; obtain final owner acceptance and merge current PR head.

### E) Dependencies And Critical-Path Notes

EP-01 source review precedes coherent broad dispatch. Broad acceptance anchors guidance test.
Use ordinary independent work during hosted waits; no duplicate matrix after metadata/integration.

### F) Tasks Checklist

- [ ] Commit/push coherent source and open one PR; dispatch broad candidate on exact checkpoint.
- [ ] Resolve only failed/invalidated lanes and send reviewed checkpoint for source-anchor acceptance.
- [ ] Create controlled eligible delta from accepted anchor; exercise both guidance native lanes and cheap negative preflight.
- [ ] Record evidence, get final source acceptance, verify current checks/ref, merge PR normally.
- [ ] Complete plan/log and report actual merged state and remaining unrelated follow-ups.

### G) Implementation Notes

The demonstration branch never becomes a public release. Use current verified v0.1.32 as upgrade
input. Keep accepted broad anchor separate from this published tag. Do not delete branches or
repeat adoption. A resumed task or documentation update does not invalidate unchanged artifacts.

### H) Open Questions

OQ-1 closes through exact checkpoint acceptance; OQ-2 remains an honest measurement limitation.

# Section 4 - Future Planning

Runtime profiling and older-installation classification are separate bounded follow-ups D-4/D-5.
Content distribution, historical migration retirement and support changes remain out of scope.

# Revision Notes

2026-09-08: Activated directly under accepted discovery and whole-delivery authority. Finish line
is producer PR merged to main; publication/adoption deliberately excluded by the commissioning grant.
