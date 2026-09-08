Last updated: 2026-09-08T15:12:57Z (UTC)
Created: 2026-09-08

# Release Validation Effort Execution Log

[Active plan](release-validation-effort_implementation_doc.md) owns scope and acceptance.
Commissioning owner accepted discovery D-1–3 at `8c99688` and authorized this complete producer-only
implementation, planned hosted validation, normal pushes/one PR and merge after final acceptance.
No public release, version/tag, consumer adoption or unrelated runtime optimization is included.

## Activation

Branch `codex/release-validation-effort` is isolated from `origin/main` at `0e78015`. Current source
was fetched again; repository has no rulesets and main protection endpoint reports unprotected.
The activation checkpoint contains only discovery acceptance, implementation plan, log and index.
Native source implementation has not started. Canonical plan validation precedes normal push.

## Progress

EP-01 implementation pending; EP-02 hosted proof and integration pending. One coherent independent
source review with same-reviewer corrections is planned. Exact broad anchor requires commissioning
acceptance before the real guidance-path demonstration. Final outcome acceptance precedes merge.

## Coherent source checkpoint

Activation was committed at `9f1bb0b` and normally pushed before source implementation. EP-01
uses three small scripts with distinct responsibilities: cumulative eligibility, cheap release
identity and shared isolated native lifecycle proof. The latter avoids duplicating new current
published-release upgrade mechanics across shell dialects. Existing historical native smoke stays.
No product source version or embedded managed guidance is changed on the delivery branch.

Eligibility accepts only known release identity literals, existing managed targets and narrow
planning/release-note surfaces; semantic policy review and accepted broad-source evidence remain
owner judgment. Installer behavior and governing workflow changes remain broad. New Windows
assertions check intermediate native exits. Review exposed an inherited checksum-negative false
positive: its unexpected-success message matched the expected error catch. Both workflow copies
are corrected, with actual PowerShell try/catch regressions for rejecting and wrongly accepting
installer stubs.

Local setup used the supported ignored development environment. An inherited package-index
authentication failure was isolated to retrieval; public producer dependencies were installed
from the public package source without shared auth or machine configuration changes.

One independent coherent source review found the checksum issue and an identity gap (installer
and Python defaults could disagree while explicit-version smoke passed). Both corrections remain
inside EP-01; identity now checks Python, installer defaults/help and known bootstrap references.
Independent hash vectors, historical fixture digests and version/corruption negatives remain.
The reviewer also passed a fresh low-context guidance-release routing scenario through producer
AGENTS, change and release routes. No semantic acceptance service or new public-policy surface.

Initial combined focused suite: 87 passed in 11.96s before final review corrections. Public-core,
Markdown, schema and release-manifest checks passed. Focused Go hash/manifest passed. Final focused
correction tests and same-reviewer follow-up precede the coherent broad candidate. Hosted native
and actual guidance branch proof remain pending; no performance gain is yet claimed.

Same-reviewer focused follow-up closed both material findings. Final correction suite: 61 tests
passed in 3.53s (workflow, actual PowerShell negative gate and identity); guard's 35 tests and
lifecycle's 3 tests remain applicable. Source review is clear; hosted acceptance is pending.
Consumer impact: validator-only producer enforcement and maintainer instructions. No shipped
managed content, runtime behavior, current release identity or consumer migration is changed.

## Hosted checkpoint and route evidence

Source checkpoint `80e179bd9ac9b4335ea097d6b49845e5257fa570` is pushed in PR #18. Broad candidate
run `34236648411` targets that exact source. PR feedback passed in 47s. Ubuntu passed in 2m02s
(Go 81s; new focused guard/lifecycle/identity tests 4s); exact Git 2.43 passed in 1m03s; macOS
passed in 5m57s. Windows broad source remains pending. These are scoped observations, not a
controlled cold/warm comparison or a diagnosed Windows cause.

Negative preflight run `34236680552` deliberately used the pre-change main anchor and rejected
`.github/workflows/validate.yml` as consequential. Every native/public job was skipped. This is
expected negative evidence, not a failed candidate acceptance lane.

Retrieved macOS-job artifacts verify through external catalog, pack manifest, payload checksums
and sidecars. Both pack digests equal the already published v0.1.32 assets (macOS `3e44ffeb...c415`,
Windows `bb7db54e...6468`); no new distributed payload was introduced. A local corrupted copy of
release-note payload was rejected by the checksum verifier. The real eligible test delta is
prepared in an isolated branch from the source checkpoint; no guidance dispatch precedes accepted
broad-source evidence.

A separate read-only agent with no conversation history performed the low-context route probe.
It found AGENTS -> producer change runbook -> classification/release runbook, distinguished
semantic policy review from the guard, required cumulative accepted broad-source anchor and
separate verified upgrade tag, retained both native proof lanes, and escalated unknown/consequential
changes. It selected no execution surface prematurely and performed no validation or external action.

Commissioning owner accepted source anchor `80e179bd9ac9b4335ea097d6b49845e5257fa570`, conditional
on the complete remaining Windows lane in `34236648411` passing, including packaging/native smoke.
No repeated anchor request is needed after that condition. Final integration acceptance remains
after the isolated real guidance result. This acceptance/log update does not invalidate source,
artifacts or successful broad lanes and does not change publication/adoption exclusions.

## Windows native completion correction

Windows broad Go passed in 38m40s; packaging passed in 6s and staged installation in 3s. The
historical upgrade smoke then failed its newly explicit exit assertion: the new binary version
was visible while `check` still returned `transaction-in-progress`. This is a smoke sequencing
race, not a failed broad suite or a diagnosed runtime slowdown. The test must wait for completed
installed state before checking preservation; it must not suppress an actual failed state.

The shared native helper now supports a narrow `wait-ready` operation used by Windows historical
smoke and by guidance upgrade. It requires expected version plus successful native check, retries
only known pending state, and fails on drift/corruption. Tests prove version alone is insufficient
and permanent failure is not retried. A `windows-smoke` partial retry lane rebuilds packages and
reruns affected native checks while retaining the successful Windows source/parity and other
platform evidence. It does not qualify acceptance by itself or parse authority receipts. This
keeps the correction within the existing retained-evidence contract without replaying 38m of Go.
The same reviewer checks the correction before its scoped hosted run. The original conditional
source-anchor gate is not claimed satisfied by the failed job; corrected joint evidence will be
reported with its exact source checkpoint.

The correction passes 39 focused workflow/lifecycle tests in 3.14s. Same-reviewer follow-up is
clear and confirms retained unchanged source evidence is sufficient alongside fresh native proof.
The new exact source anchor requires commissioning acceptance; the original failed job alone
does not meet its condition.

## Pre-reconciliation observation correction

Scoped Windows run `34241685510` at `db171fcee901a046cb074a0cd3d718eb6a866943` passed
resource/identity/build/staged install and skipped retained broad suites, but historical upgrade
reported `partial` before the wait completed. Source inspection of `ApplyHandoff` shows the new
executable is copied before starting `__upgrade-reconcile`, so there is also a pre-marker interval
in which the new graph observes an old-version repository. The previous helper handled the
transaction-marker interval only. This does not establish a permanent historical runtime defect.

The helper now observes that explicit version-mismatched old tree as pending, requires successful
expected-version installed-state check, and never treats partial state as success. Current-version
partial fails immediately; persistent old-version state times out with full last check diagnostics.
Focused tests cover old-tree -> transaction -> current and both persistent partial variants,
including retained missing-path evidence. All 44 affected lifecycle/workflow tests pass in 2.50s.
No runtime source, package bytes or broad suite inputs changed; prior passing evidence remains
applicable. The affected Windows smoke is rerun after same-reviewer correction review.

Same reviewer confirmed failure gates and requested the release runbook match the expanded
pre-reconciliation observation boundary. That instruction is now aligned; no other findings.

Corrected checkpoint `c36061f69f18b8a9aff2e018db2d5e35c6188ec8` passed scoped Windows run
`34242213815` (2m34s total, 15:03:00–15:05:34 UTC). Both historical upgrades reached expected
version and successful final check, including preservation assertions. PR feedback passed.
This confirms the synchronization correction without runtime edits. Joint evidence retains the
unchanged broad source suites from `34236648411` and replaces failed native completion evidence
with this scoped run. Corrected exact-anchor acceptance was requested before guidance dispatch.

The commissioning owner accepted `c36061f69f18b8a9aff2e018db2d5e35c6188ec8` as the corrected
broad-source anchor, superseding the unmet conditional decision for `80e179bd`. Acceptance
explicitly combines retained broad results with corrected Windows smoke; it does not convert the
original failed workflow into a passing run. No new broad matrix is required.

Isolated demonstration `721ebd7c62ecaa3e4bfa8ea12f4676d4743c4065` passes the actual cumulative
Git guard and release identity check against that accepted anchor. Its only managed content
delta is a routine-change guidance paragraph with matching resource mirror and consistent test
identity/checksum changes. Both native guidance lanes were dispatched with published `v0.1.32`
as the separate old-CLI upgrade input. Test identity `0.1.33` remains unpublished and outside the
delivery branch.

Initial guidance run `34242697847` passed eligibility and focused Go identity, then both native
resource/schema test steps rejected the demonstration manifest's serializer-generated indentless
YAML lists. This is a test-candidate preparation error; the existing parser contract was correctly
enforced. The isolated candidate now preserves supported indentation and its source mirror.
Corrected candidate `115e53d` passes all 86 resource/schema/routing/workflow tests locally (6.54s),
release identity and cumulative eligibility against the same accepted source anchor. Both affected
guidance lanes were redispatched; no producer source or accepted broad evidence changed.

## Successful guidance demonstration

Run `34242965649` at isolated candidate `115e53db3836216e8738a395b6dd2c2b52c4bdc8` passed
against accepted broad-source anchor `c36061f69f18b8a9aff2e018db2d5e35c6188ec8`. Eligibility
preflight took 8s. macOS native job took 1m11s and Windows 1m27s; both restored Go caches. These
are observed hosted runs, not controlled speedup benchmarks. The original broad Windows Go step
alone took 38m40s; that unchanged work was not rerun for the guidance candidate.

Both native logs prove 50 fresh and 50 upgraded managed resources equal candidate source bytes,
including the changed runbook. Matching published `v0.1.32` CLI initiated upgrade to unpublished
test `0.1.33`; dry-run, failed verification, successful apply/check and authored preservation pass.
Existing historical upgrade, staged install, containment and checksum-negative checks also pass.
macOS package build took 21s and Windows 7s; fresh content/lifecycle proof took 4s and 5s.
Broad Go, Python parity/history, Ubuntu semantic/benchmark, exact Git 2.43 and public-release jobs
were skipped exactly as intended. Focused Go identity remained active (3s macOS, 9s Windows).

The guard's earlier negative run remains valid. This positive exercise and its correctly detected
fixture error show that eligibility permits fresh focused validation, not automatic acceptance.
No public version/tag/release or consumer adoption occurred; delivery still carries v0.1.32.
Final integration acceptance is the remaining planned boundary before normal PR merge.
