Last updated: 2026-09-05T23:32:12Z (UTC)
Created: 2026-09-05

# Proportionate Agent Workflows Execution Log

Plan: `proportionate-agent-workflows_implementation_doc.md`
Status: completed
Mode: whole-plan delivery with delegated Director acceptance

Current outcome: technical delivery accepted, v0.1.31 released and adopted by the authorized
pilot consumer. Real-use experience is pilot-pending. Earlier checkpoints below are historical.

## Commissioning and acceptance

The user approved the proposed sequence and full delivery authority on 2026-09-05. The canonical
Plan section 2.4 records covered Git, merge, release and adoption effects, ownership and exceptions.
The commissioning Director is the acceptance and integration owner. One implementer will own both
epics and report directly at review, completion or a genuine blocker; no continuous watcher is used.
The private assignment records the commissioning and implementer task locators and exact consumer.

## Activation evidence

- Source main refreshed on 2026-09-05 and unchanged from planning base `be4b92e5f997abb3ba7b7de03f85937c83ef8193`.
- Canonical execution branch: `codex/proportionate-agent-workflows`.
- This bounded activation checkpoint contains only this Plan, its execution log and the plans index.
  The commissioning handoff records its actual pushed commit before source work begins.
- Product implementation, release and adoption: not yet performed.
- Catalog baseline: the unchanged frozen register reports `legacy_relation_duplicate` for
  `OK-PR-002`; this Plan introduces no new catalog error. The user-approved prerequisite in 2.4
  permits bounded correction after this checkpoint. Do not claim global validation passes yet.

## Epic progress

| Epic | State | Acceptance |
| --- | --- | --- |
| EP-01 — Deliver the reusable workflow | accepted by Director at source commit 1ef4fe7 | Director reviews the coherent source result and evidence. |
| EP-02 — Release and consumer handoff | candidate preparation in progress | Exact release, consumer preservation and pilot readiness required. |

## Validation and delivery evidence

Record commands/results, source review, commit/PR state, release identity and adopted version here
as they occur. Source completion, branch push, PR, merge, release and adoption are separate facts.
No source, release, adoption or user-experience result is claimed by activation alone.

Activation checks: new Plan metadata has no error; configured catalog retains exactly its one
recorded baseline error. Changed Markdown timestamp, public-core and whitespace checks passed.


## EP-01 — Deliver the reusable workflow

The bounded prerequisite was completed before workflow editing. A focused regression first failed
on the two distinct wrapped prose relations, reproducing invented duplicate target `first`.
Projection now accepts one target with an optional ` - title` suffix; other multiword prose remains
malformed evidence instead of becoming a fabricated identity. This corrects classification rather
than changing severity policy. Canonical frozen-history handling retains malformed evidence as
warnings; legacy and mixed remain errors. Genuine duplicate IDs (including different display
titles), duplicate paths, ambiguous identity/path guards and schema-v1 projection remain unchanged.
The frozen register and catalog configuration have no diff.

Source CLI `plans validate` passes with 40 canonical records. Full
`go test ./internal/plancatalog -count=1` passes (189.699 seconds). The focused regression covers
wrapped/unwrapped prose, actual ID/path duplicates, valid targets, legacy/mixed strictness,
canonical compatibility and schema-v1 continuation behavior.

The reusable workflow now includes the routine change runbook, generic agent-task coordination
reference and optional Codex operations note. Drafting, execution, review, routing and runbook
authoring agree on whole-plan commissioning, authority reuse, coherent Git effects, direct report-back,
material exceptions and distinct final delivery boundaries. Formal plan titles use meaningful H1s.
Existing specific approval contracts, exact-effect bindings and enforced gates remain binding.
No CLI permission semantics, schema, scheduler, role registry or migration framework changed.

The routine recipe remains L1, with existing-record evidence and owner-specific validation; no
script promotion is justified. All changed declared managed/template sources are mirrored exactly.
The new routes materialize from packaged resources, and root guidance discovers the routine and
coordination routes. The existing installed-module discovery instruction is preserved.

Validation:

- `python -m pytest -q tests/test_packaging_resources.py tests/test_routing.py`: 16 passed.
- JSON schema, sync/preservation and plan-catalog backward-compatibility tests: 59 passed.
- Materialized a neutral consumer and inspected all three new routes and root pointers.
- Changed-file Markdown timestamps, public-core validation and whitespace checks pass.
- Component checksums and standard graph identity are refreshed in source/resource manifests.
- The Python test environment is repository-local development tooling, excluded from source.

Impact review: instruction-only guidance plus validator-only classification correction, and
explicit security/safety-policy impact for the changed interpretation of whole-plan external-action
authority. The release notes disclose all three. The latter is not classified as merely prose.
Independent review found no material issue; Director acceptance remains pending. Source readiness is not release/adoption.

Execution continuity: an initial enforced filesystem/reporting boundary delayed checkout. After
user-controlled access changed, assigned-branch checkout and direct Director report-back succeeded.
No history rewrite or alternate checkout was used. The whole-plan goal had been verified active,
then marked blocked at the enforced gate. The current API cannot resume that status; authorized
implementation continues without claiming reactivation. This does not narrow the whole finish line.


## EP-01 Independent Review

One fresh read-only reviewer examined the complete source diff against Plan 2.4 and EP-01. No
material findings. It independently checked parser classification/compatibility boundaries,
authority reuse, the installed-path link graph, all 13 changed managed/template source mirrors,
root/resource identity equality, component hashes and `go test ./internal/manifest`. Full catalog
and resource/routing suite results were supplied by the implementer. This was the required source
review, not the single shared consumer walkthrough. No source changes were required by review.
Review was useful for the changed safety doctrine and classifier boundary. Director acceptance
and the shared walkthrough remain distinct gates.


## EP-02 — Release and consumer handoff

Director accepted EP-01 at `1ef4fe7` after review of the classifier, workflow doctrine, declared
routes, evidence and independent review. That commit is pushed on
`codex/proportionate-agent-workflows`; delivery PR is
[PR 14](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/pull/14).
EP-02 continuation is authorized under Plan 2.4. Integration/publication/adoption remain with the
Director until an exact candidate is explicitly delegated.

Release inventory on 2026-09-05 showed latest `v0.1.30`; next candidate is `v0.1.31`. Current-version
Go/Python, installer/bootstrap, component/profile/content identities and version-dependent fixtures
are prepared using the existing release route. Historical release notes and compatibility baselines
are preserved. Required macOS/real-Windows, reproducibility, catalog-to-binary and signing/audience
gates remain mandatory. Candidate preparation is not release availability.

The single shared fresh-context walkthrough will use the actual isolated Kit candidate and the
companion module's installed candidate. Source review is not substituted for installed release
proof. Consumer locators remain outside this public log.


Release correction: CI on candidate `f043070` exposed the standard-profile golden hash still bound
to v0.1.30. Regenerated its exact expectation using the existing Python SHA-256 helper and verified
all three helper fixtures plus `go test ./internal/hash -count=1`. The assertion remains exact.
The corrected head, not the failing candidate, must supply final release/platform evidence.


## EP-02 Local Candidate Evidence

Validated product source: `be3936b28dbb6190d9d5e993d346499f82f31511`, version `0.1.31`, on
`codex/proportionate-agent-workflows` / PR 14. No product source changed after these builds.

- Two independent `scripts/build-release-assets.py` invocations completed. Each internally repeated
  both macOS universal and Windows x64 builds; the two complete output directories also match byte
  for byte. Their source-validation phase runs `go test -timeout 30m ./...` and passes.
- Full Python suite: 156 passed initially, with a local test-environment Markdown scan conflict and
  the stale profile-hash candidate build as its two failures. After correction, both failed cases
  pass (29.03 seconds), yielding 158 effective passes. No assertions or validation gates were waived.
- The default environment path exposed dependency license Markdown to the blanket repository
  validator. The task-owned environment was moved under the existing `.venv` exclusion inside the
  ignored local tooling layer. Full default Markdown/public-core validators now pass. This is a
  tooling-layout exception, not a source-validator change.
- JSON schemas, content identity, `go vet ./...`, whitespace and configured source catalog pass;
  the catalog contains 40 records. Frozen register/config remain byte-identical.
- Verified external catalog -> archive -> pack manifest -> every payload checksum -> content
  manifest and binary/version chain. Packs have the required assets and no Python payload. All
  staged public asset sidecars match. Public URLs exist only in the staged external catalog;
  nothing has been uploaded by this preparation.
- Isolated macOS fresh install and `0.1.30 -> 0.1.31` upgrade preview/failure/apply/check pass.
  Dry-run and missing-catalog failure preserve old binary and lock hashes. Four consumer-owned
  config/record/local-note/memory hashes and the root local-instruction sentinel survive upgrade.
- A separate macOS PATH containing no Python/pip commands rejects a bad checksum and passes
  fresh install, init and check. Universal binary has arm64 and x86_64 slices; Windows binary is
  PE32+ x86-64. Native Windows execution remains a CI gate, not inferred from cross-compilation.
- macOS binary has linker ad-hoc signature, no Developer ID/team identity and no notarization.
  Broad public distribution is not authorized by candidate readiness. Director must preserve the
  existing internal/prototype audience boundary and action-time release conditions.

Candidate identity:

| Artifact | SHA-256 |
| --- | --- |
| macos-universal archive | `02d469078b8896b790e9e0a0b1cf33d1e483ed5215229dae8530aff5db48db64` |
| macos-universal pack manifest | `972f6824b36d729fcf6004cc1f4819e86b6a04f418a312f54b2912f6955fb3ac` |
| macos-universal binary | `7388e1eee4e986b0070ee5e765d9e20b4b627373819c4874152ad157d6855d89` |
| windows-x64 archive | `78d487fa7ed48e6aa97f217b3e22e70e9589ef790534753f5616636920d75e2f` |
| windows-x64 pack manifest | `9d72073ed6e78da85bc281d4891a86932fd927632d2b86fa1b8cc9cee805d0df` |
| windows-x64 binary | `bedec2b9c52952847ad1f475c99001ec8a740d6b221fabaaf809784805a913ea` |
| Shared content manifest | `53cb77f53eb0b991c17c58661807d146173c17bbe079ee8a963c6711175785bd` |
| Staged public-URL catalog | `b25e38cc10e4f64598f81c5e49060e794a56db240b8e8f726cb9060d67eae988` |


## Single Shared Installed Walkthrough

One fresh read-only reviewer started from a synthetic consumer's installed root instructions and
discovered both owners without using the source Plan as an answer key. It read Operating Kit
`0.1.31` and Organization Home `0.1.2` bundle schema 3. It verified all 50 Kit managed checksums,
38 module-manifest members and nine shared components; all 52 pre-existing module/tool files were
preserved by Kit init. The combined Kit check passed.

Kit source association is retained execution provenance at `be3936b`; the installed lock honestly
records `embedded-running-binary` / `local-source`, not an embedded source commit. The companion's
build evidence names `523f4041fe2fb2460990d99a74a3b9ad6136683e`; its separate release declaration
names checkpoint `1d88f66323fc2b06340b70bd8230826f4702c38e`. Bundle digest is
`a8b9c2d769107f7dc4ad9d2b8340c9e8a8220ea7c2bcbd448000192be5775859`; runtime wheel digest is
`597eeec3f4c11d07592a358aa459a44646587661e283bf6c753dc9d0ee67f71a`.

Walkthrough outcomes:

- Root -> operation router -> `planning.execute-whole-plan` -> coordination and execution routes
  resolves the authorized two-epic assignment. Covered implementation, tests, coherent commits,
  normal pushes and one PR proceed; required review returns to the Director. A sent message is not
  acceptance; independent release preparation may continue while dependent work waits.
- Material outcome/scope changes require a Plan amendment with the right owner; genuinely separate
  routine repairs retain an existing record and relationship. Necessary epic corrections stay in
  the epic. Rejected reporting preserves the result and discloses nondelivery without a watcher.
- Active-descendant archival requires reachable transferred review ownership or retention of the
  parent, plus worktree preservation. Goals are explicit-only and require observed activation;
  no goal creates authority or completes an intermediate review handoff.
- The nested Program Governance route discovers `organization-home` through the installed module
  system, then its guidance, operation contract and reconciliation route before choosing a tool.
  Governance prose is material: preserve typed payload/owner/history and exact before/after and
  approval bindings. Prior sufficient authority is reused; downstream effects retain their own
  owner gates. Proposal intake does not silently become an approved Plan or a second epic tracker.

No blocking workflow contradiction was found. One nonblocking companion wording ambiguity about
Program/Plan endpoints in the Strategy/Capability section was reported to its owner for assessment.
The walkthrough itself made no source, provider, task, archival or consumer-record changes. It is
not live release/adoption or observed user experience. Do not repeat this clean shared walkthrough
without a material change.

## Integration Handoff State

EP-01 accepted. EP-02 local candidate and shared walkthrough evidence are ready. At this checkpoint,
corrected-head Git-2.43 and Ubuntu semantic CI pass; macOS and real Windows jobs remain in progress.
Director acceptance, required platform results, merge, release publication and authorized pilot
consumer adoption are still required. Final release assets must match the reviewed candidate and
applicable audience; release availability and adoption cannot be inferred from this source evidence.
User experience remains pilot-pending until observed. The Plan stays active and the whole finish
line remains unchanged.

## Published Release And Final Native CI Evidence

Source PR 14 merged as `29be2378af80e58af212bbc4e652bb2a8dedeb93`. Published release
[v0.1.31](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v0.1.31)
has a lightweight tag at exact validated source `9ba3ce561660f296d22e15c6e42661e402ec6d31`,
which is contained in merged main. Product build inputs remain those validated at `be3936b`;
the later source checkpoint only added Plan/execution evidence.

Both existing current-head runs finished successfully:
[push validation](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/33992880223)
and [PR validation](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/33992883932).
All 25 macOS and all 20 Windows job steps succeeded in each run, including release builds,
no-Python installation, checksum rejection, old-version upgrade preview, failure preservation,
apply and check. Git 2.43 and Ubuntu semantic checks also passed. No duplicate workflow was
triggered to obtain this final evidence.

After publication, an authenticated fresh download read back all 14 advertised assets. Their
names, sizes, bytes and digests match the frozen reviewed candidate; all seven sidecars pass.
Both catalog/archive/pack-manifest/payload/content-manifest/binary chains are intact. Downloaded
macOS binary reports 0.1.31. Windows declares 0.1.31, confirms the Windows amd64 build target,
and matches the frozen reviewed binary digest. Native Windows behavior uses the existing CI
source-release evidence; no fresh execution of the downloaded Windows binary on macOS is claimed.
No new build, full platform matrix or second walkthrough was needed for unchanged publication.

The release body retains the explicitly accepted HTTPS-plus-SHA-256 internal/prototype boundary
with ad-hoc macOS linker signature, no Developer ID/notarization or equivalent publisher signing.

| Published asset | Bytes | SHA-256 |
| --- | --- | --- |
| [bootstrap.md](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/bootstrap.md) | 12041 | `37d4650067b0f06a23d9f9e634784aefbb01fcf8e67c41ec2d0db41939234eed` |
| [bootstrap.md.sha256](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/bootstrap.md.sha256) | 79 | `9cbf813a968f94f85096d65f4e817b81a1ea14e55dd14d338fb3b2bf7b99936f` |
| [codeheart-operating-kit-0.1.31-macos-universal.zip](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/codeheart-operating-kit-0.1.31-macos-universal.zip) | 8186244 | `02d469078b8896b790e9e0a0b1cf33d1e483ed5215229dae8530aff5db48db64` |
| [codeheart-operating-kit-0.1.31-macos-universal.zip.sha256](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/codeheart-operating-kit-0.1.31-macos-universal.zip.sha256) | 117 | `00b47cf96b67d0109da45ea446a1e176a22a47367970b90453ff3e6f21d74d8e` |
| [codeheart-operating-kit-0.1.31-windows-x64.zip](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/codeheart-operating-kit-0.1.31-windows-x64.zip) | 4321392 | `78d487fa7ed48e6aa97f217b3e22e70e9589ef790534753f5616636920d75e2f` |
| [codeheart-operating-kit-0.1.31-windows-x64.zip.sha256](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/codeheart-operating-kit-0.1.31-windows-x64.zip.sha256) | 113 | `e68b00b3ff8b80c73462b38c93376dd1957d72166f48a15d1ac3c331e15b7550` |
| [install.ps1](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/install.ps1) | 14052 | `9930ad2541acd642f2ce3bcc7d3a3adcac8a7e60c63a6b86eb82d742aee27f68` |
| [install.ps1.sha256](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/install.ps1.sha256) | 78 | `403698b947f4d54c013e756aa622f5adbb39d63a73c77ad66c609fea02843b90` |
| [install.sh](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/install.sh) | 11501 | `eef59ccd60c8152819e8a7398f8877a451eda33211c2e240e714fe0356f56416` |
| [install.sh.sha256](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/install.sh.sha256) | 77 | `724d0696d1b38c0f506887d322a804990035336a7a2ce19b6ae0f265a483ced0` |
| [release-catalog-0.1.31.json](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/release-catalog-0.1.31.json) | 1051 | `b25e38cc10e4f64598f81c5e49060e794a56db240b8e8f726cb9060d67eae988` |
| [release-catalog-0.1.31.json.sha256](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/release-catalog-0.1.31.json.sha256) | 94 | `1ce0d726d77d335e132ee9fb815469c772f71bb7e089b6b5e1e22f93cc6d2fb6` |
| [release-notes.md](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/release-notes.md) | 78384 | `20c68b2ad1fa883e6263689db4956dd15d5255a12daecf57bea1b1c162c28b35` |
| [release-notes.md.sha256](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.31/release-notes.md.sha256) | 83 | `0799e959ec0178159f301926c840f889625b2994c7d61cab08f91c73b0e7b401` |


## Shared Walkthrough AWS-Owner Follow-Up

The same reviewer completed the omitted AWS-owner question as a bounded follow-up to the single
shared installed walkthrough. The synthetic installed inventory exposes only Organization Home,
excludes live providers, and has no AWS owner advertisement or committed provider routing state.
Before choosing an execution tool, resolve the owning product/repository and approved change or
runbook, then the intended effect, authoritative account/environment/region and target, sufficient
prior authority, and applicable preflight/gates. The Kit/module delivery grant does not authorize AWS effects.
The reviewed routes establish missing owner/context; they do not supply an AWS execution recipe
or authorize inferred target facts. No AWS access, provider action, configuration inference or
mutation occurred. The module owner received this evidence mapping.


## Authorized Consumer Adoption And Technical Acceptance

The Director accepted the post-publication proof and completed the named pilot-consumer upgrade
through the published external catalog using the managed lifecycle route. Preview and approved
apply succeeded with the reviewed catalog/archive/binary identities; final Kit check reports
`current` / `ok`, with no drift or missing routes. Installed Kit version is 0.1.31.

Consumer-owned records and local guidance, the repository-owned root-instruction section, frozen
legacy files, unrelated installed modules and shared module-tool payload were preserved. A
separate read of the committed consumer checkpoint verified that all three new routine/coordination
routes exactly match released source and that the repository-owned root-instruction section is
byte-identical to its parent. Consumer-specific locations and operational proof remain with that
owner rather than in this public repository.

The coordinated Organization Home 0.1.2 adoption also passed its owner checks. The authorized
Program Governance clarification reconciled across 12 records and four foundations while
preserving payload, owner, assertions and history; both existing older module bindings still
resolve. These are companion-owner acceptance facts, not new Kit runtime features. The earlier
module build identities in the shared walkthrough remain historical clarity evidence; they are
not asserted as the final adopted module payload. Final module release/provenance is retained by
its source owner. No second whole walkthrough or unrelated live provider operation occurred.

The Director gave final technical acceptance after release and actual consumer adoption, then
integrated the adoption checkpoint into the consumer main branch while preserving pre-existing
unrelated work. Both ordered epics are complete. No product input changed in this closure: only
this execution log and the canonical Plan status/checklist changed. Markdown timestamps, public-core
hygiene and whitespace checks passed. The released CLI validates the configured source catalog:
40 records, `valid=true`, with existing canonical compatibility warnings and no new Plan error.
A complete changed-path comparison against released source proves that all product/build/release
inputs and frozen register/config remain byte-identical. Existing native/source release proof
remains valid; no product rebuild or deliberate full platform rerun is required.

## Pilot Handoff And Retained Evidence

Source completion, validated release availability and technical consumer adoption are complete.
The consumer's coordination owner owns the next practical Program Director pilot: use the
installed workflows on real commissioned work and route concrete usability observations to the
appropriate source owner. Actual user experience remains explicitly pilot-pending; no successful
real-use trial, fixed pilot duration, quota or telemetry collection is claimed or required here.

Published URLs and immutable source/check identities above are the public evidence. Detailed
non-secret download/digest proofs, original/repeat builds, installer preservation results and
private consumer-owner proof remain retained for audit and the pilot handoff. No needed artifact,
consumer content, worktree or active task was removed. The applicable release route has no
mandatory disposable-root cleanup step, so no speculative deletion was introduced for closure.
No product delivery obligation remains; the Director retains normal integration review of this
single documentation closure and the handed-off practical pilot.
