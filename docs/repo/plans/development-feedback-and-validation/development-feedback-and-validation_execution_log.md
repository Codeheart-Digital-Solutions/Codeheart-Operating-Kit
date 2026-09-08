Last updated: 2026-09-08T00:15:42Z (UTC)
Created: 2026-09-07

# Development Feedback and Proportionate Validation — Execution Log

## Authority and delivery boundary

The user authorized the complete plan on 2026-09-07, including source integration, release and
verified adoption. The [active plan](development-feedback-and-validation_implementation_doc.md)
owns scope, acceptance and exclusions. The Program Director is the acceptance owner; the dedicated
implementer receives its exact report-back locator and private adoption targets in the assignment.

Implementation branch: `codex/development-feedback-delivery`. Planning baseline: merged
[PR #16](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/pull/16), commit
`9e73c654aed7de161f0414e37701f8426b1cd44b`.

## Activation checkpoint

- The reviewed draft and repository index were merged before activation.
- Independent planning review closed the candidate/public-release dispatch sequencing issue.
  The later user-approved publication/lifecycle clarification is in the combined source review.
- The activation checkpoint includes only the plan, this log and the index lifecycle label.
  Product source, workflows and installed consumers are unchanged at activation.
- Documentation-only checkpoints use `[skip ci]`; their relevant local document/catalog checks
  supply evidence. This does not waive any implementation-candidate or release acceptance.
- The active source branch is to be pushed before dispatch; the commissioning task records the
  resulting checkpoint and actual task creation. No goal or monitoring automation is requested.

## Progress

| Epic | State | Completion evidence |
| --- | --- | --- |
| EP-01 — Coherent operating guidance | Accepted at 861ed11 | Guidance and declared mirrors updated; independent review correction addressed |
| EP-02 — Feedback, enforcement and coherent acceptance | Accepted at 861ed11 | 27 focused tests pass; workflow mode guards exercised locally |
| EP-03 — Release, adoption and owner handoff | Release inputs prepared; delivery pending | Next unused patch selected: 0.1.32; native candidate and adoption pending |

EP-01/02 share one coherent independent source review with focused correction follow-up. After
Director acceptance, EP-03 proceeds through the already-authorized final effects and reports
actual release and per-target adoption results. Record measured outcomes and relevant deviations;
do not claim savings or completed consumer product changes from a process release.

## Combined guidance and validation checkpoint

Impact: instruction-only managed guidance and validator-only producer enforcement. No placement,
consumer ownership, authentication, support or release-audience changes. Release notes describe the
normal upgrade and route verification requirement. Source authority is tracked producer content;
all affected declared Python resources match it byte for byte.

The implementation preserves active intention elicitation and substantive technical planning,
removes checklist-word prohibitions, adds coherent batch/evidence/recovery guidance, uses one
independent review with same-reviewer corrections and establishes local commit/publication/merge
rules independent from lifecycle. Recipe guidance tests necessity before adding machinery.

### CI mapping and current enforcement

GitHub live preflight found main unprotected (branch protection endpoint: `Branch not protected`)
and no repository rulesets. No required check names are being silently skipped or renamed.

| Lane | Trigger and evidence |
| --- | --- |
| feedback | PR and main push; routing/resources, public-core, Markdown, schemas, release manifest; feature branch push does not duplicate PR |
| git-2-43-proof-validation | Explicit candidate only; exact pinned Git build and branch-proof regression |
| macos-validation | Explicit candidate; compatibility/parity, twice-built native packs, install and historical upgrade preservation; visible CI step owns broad Go invocation |
| ubuntu-semantic-validation | Explicit candidate; broad Go, separately configured scale benchmark, schema/routing/resource/sync and compatibility |
| windows-validation | Explicit candidate; compatibility/parity, Windows-only twice-build, native install/reparse and deferred-upgrade preservation; visible CI step owns broad Go invocation |
| macos/windows-public-release | Explicit released-smoke with required versioned tag; public download/install/check only |

Mode guard is exercised with valid candidate, rejected candidate tag, rejected empty/malformed
smoke tag, valid smoke tag and unknown mode. Candidate jobs do not download the unpublished
version. Public smoke does not set up Python or run source suites. Only superseded ordinary
feedback is cancelled; each manual candidate/install dispatch has its own concurrency group.

Inspection found additional overlapping full-suite calls inside the existing pack builder and
release tests. The initial native workflow revision relied on the builder's one Go invocation and twice-built
pack verification; the later diagnostic correction separates source tests from packaging; CI omits equivalent full-builder Python cases. The distinct Windows-only
builder, oldest Git regression, benchmark and native preservation checks remain. Builder execution
under an unused package-index URL retains evidence that packaging needs no Python package access.

A supported local tooling environment exposed the Markdown validator traversing generated vendor
files. The validator now excludes only `.codeheart/local/` during directory traversal. A regression
proves similarly named consumer-owned `docs/local/` content is still checked. This directly related
EP-02 enforcement correction adds `scripts/validate-markdown-headers.py` and its existing tests to
the starting scope. No new evidence framework or classifier is introduced.

### Validation and review

- Focused routing, resource, workflow-boundary and Markdown regression suite: 27 passed.
- Public-core, Markdown, JSON schemas and release-manifest validators pass.
- Release identity updated to next unused patch 0.1.32; `internal/manifest` passes, including
  embedded graph identity. Full native candidate evidence is still pending.
- One independent read-only reviewer walked all requested scenarios. Routine repair, deferred
  broad checks, focused correction, published draft and later bounded activation, interruption,
  nested module/tooling/service routing and proposed-method intention elicitation pass at source
  walkthrough level. This is not live service, hosted CI or adoption evidence.
- One material review finding identified stale mandatory script-scaffolding/common-output bullets
  in planning/execution guidance. These were aligned with the owner standard and mirrored; focused
  same-reviewer follow-up closed it. The corrected selector retains the cheap version-mismatch
  test, which passes (1 test); no material source findings remain. An apparent discovery batching conflict was retracted after
  verifying it applied only to an explicitly requested per-decision review mode.

Tooling setup used the supported ignored local development environment. An inherited package-index
access failure was isolated to package retrieval and resolved using the public package source for
this public producer development environment, without changing authentication or machine policy.
No account budgets were changed and no broad suite was blindly retried.

The focused reviewer also identified a broad test selector excluding the cheap negative version
mismatch guard; it is now explicitly retained. Native builders cover required pack integrity,
not every historical Python assertion: separate-invocation Windows catalog comparison and captured
private-index-marker omission are not replayed. Those tests remain available for builder changes.

Source acceptance is complete. Hosted candidate results, integration, publication, public smoke
and all assigned consumer adoptions remain open. Exact target paths and preservation records stay private.

### Hosted candidate correction

Full candidate run `34167972472` targets source `861ed11`. Windows stopped in the new dispatch-guard
Python test: `bash` resolves to a WSL launcher with no distribution. This is a test-environment
mismatch, not Kit native behavior. The guard actually runs on Ubuntu; its executable cases now
run on POSIX while Windows still checks workflow structure and all native behavior. The product's
Windows runtime remains untouched.

Added a bounded `candidate_lane` choice (default all) to permit a focused hosted correction run
without replaying valid macOS, Ubuntu or oldest-Git evidence. Partial runs never qualify a release
alone. This uses existing workflow conditions; no proof cache or orchestration service is added.
Source/packaged payload is unchanged by this test/dispatch correction. Retained evidence must name
the original full run and final commit; any failed or invalidated lane remains open.

Director accepted the combined EP-01/02 source checkpoint at `861ed11` in PR #17. The acceptance
covers source outcome; hosted, publication and adoption gates remain binding. Guidance and release
inputs are unchanged by the subsequent test-harness correction. No new user permission is needed
for the covered final effects.

Ubuntu's broad suite found a stale fixed Python-helper checksum for `profiles/standard.yaml` after
the normal patch-version update. The expected fixture digest was recomputed from the unchanged
hash algorithm and current profile. This is release-fixture maintenance, not a product change.
Its affected native/semantic lanes must complete; the successful exact-Git proof remains valid.

### Native source diagnostics boundary

The corrected Windows run remained in the builder beyond 35 minutes without exposed output.
The builder captured its whole source-suite subprocess, preventing useful live diagnosis. Source
validation is now an explicit visible native CI step, retaining the 30-minute per-package Go
timeout. There is still one broad source suite per lane; no validation is waived and builder
invocation alone is not release acceptance. This changes orchestration, not embedded release
inputs or the pack algorithm. Successful macOS/Ubuntu/Git evidence and verified pack bytes remain
applicable. Windows is pending.

Same-reviewer focused follow-up found no material issue in the visible-source-step correction.
Twenty focused workflow/release tests pass. The earlier run was confirmed still in source/build
validation before cancellation; no staged install, release, consumer or provider operation was in
progress. However, subsequent historical evidence corrects the diagnosis: successful Windows run
`33998404410`, job `101392834844`, took 38m42s for the same Go command. Its per-package timeout does
not bound compilation plus total execution time. Duration alone did not establish a stall; the
cancellation was premature. The unsupported new 35-minute step bound is removed. The replacement
run is allowed to finish; no further cancellation or rerun follows merely from elapsed time.
This correction preserves the original timeout policy and does not change account budgets.

### Candidate acceptance

All required candidate lanes now pass jointly:

- macOS run `34168213234` at `bec52c4`: native validation, both reproducible packs,
  staged install and historical upgrade preservation.
- Ubuntu run `34168211657` at `bec52c4`: semantic, compatibility and scale validation.
- Windows run `34170316992` at `62fba97`: native Go (34m39s), compatibility, reproducible
  Windows pack, staged install and deferred upgrade preservation.
- Exact Git 2.43 proof job `101882711480` in original run `34167972472` at `861ed11`.

These are joint results, not acceptance from a selected lane alone. Final source `f555742`
removes only an unsupported whole-step timeout from the Windows-tested orchestration and corrects
its diagnostic record. The source suite, embedded payload and pack algorithm are unchanged.
Ordinary PR feedback passes on this final source. Same-reviewer follow-up closed the timeout
correction. No further candidate run is justified by the record update or unchanged integration.

Both public packs were retrieved from the successful macOS job. Catalog/archive, pack manifest,
payload/content manifest, binary version/digest and every public sidecar were verified. The
verified native binary materialized all 13 changed managed resources byte for byte. Pack SHA-256:

- macOS universal: `3e44ffebc3f566a83ade6b90785f4f24ad71a646c9f3ff4702e8ce26fc84c415`
- Windows x64: `bb7db54e5507a5a864a64c7d451a0c630e1c15bc431363de0f774294f6c16468`

Main protection and repository rulesets were rechecked before integration: no enforced branch
protection or rulesets. Director source acceptance and the existing unsigned/unnotarized
HTTPS-plus-SHA-256 internal/prototype grant cover the release. Publication and adoption remain
the next actions, not claims inferred from candidate success.

### Release and adoption completion

PR #17 merged at `0e7801598c9b753e5c39df36f56f4da4c112a50a`; its tree exactly matches
accepted checkpoint `dbb581b`. Tag and release
[v0.1.32](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v0.1.32)
contain the verified fourteen public assets. The established unsigned/unnotarized
HTTPS-plus-SHA-256 internal/prototype boundary remains unchanged.

[Published-asset smoke run 34172479978](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/34172479978)
passed both native public jobs against explicit tag `v0.1.32`, without replaying source suites.
All four assigned consumer worktrees were previewed and upgraded through the supported lifecycle.
Each now checks current at 0.1.32, with all 13 changed routes matching producer source byte for byte.
All lock-declared managed checksums pass. Consumer configuration, authored plans and other tracked
owned files, instructions outside the managed block and local-user content match their baselines.
Scoped adoption commits leave all four worktrees clean. Exact roots, branches, commits and raw
operational evidence remain in the private handoff, outside this public repository.

One adoption readiness correction was necessary: the shared newer CLI classified three pristine
0.1.30 installations as partial because newer managed paths did not exist in their old graph.
The official 0.1.30 pack was retrieved and its catalog-to-binary chain verified. That matching CLI
confirmed each installation current and initiated the normal explicit upgrade to 0.1.32, with the
released target performing reconciliation and post-check. No repair, lock edit, bypass, source
patch or provider operation was used. The shared installed CLI remains 0.1.32.

The final completion record is committed and published on the delivery branch after the release;
it changes no release input and requires no second candidate suite or implementation PR. Source
is merged, the release and public smoke are complete, and every assigned active instruction route
is verified. Later Director onboarding, consumer-plan amendment and product implementation remain
separately governed work. No measured efficiency gain is claimed.
