Last updated: 2026-08-06T23:46:20Z (UTC)
Created: 2026-08-06

# Ownership-Aware Repository-Wide Plan Catalog Discovery Execution Log

Plan: `ownership-aware-plan-catalog-discovery_implementation_doc.md`
Mode: goal-style implementation
Status: active
Overall divergence: reviewed pre-activation defaults were incorporated before activation. EP-01
also updates the shared Python YAML emitter so fresh empty exclusion arrays remain arrays, promotes
the directly imported Unicode package in `go.mod`, and dispatches existing portfolio-v1 cache
validation through the historical schema; these are bounded support changes required by the
planned fresh-default, portability, and historical-evidence contracts.

## Summary

The user explicitly activated the producer implementation on 2026-08-06. Execution uses
`codex/operating-kit-multi-root-plan-catalog` and follows EP-01 through EP-10 linearly. The plan,
this execution log, and the nearest plans index form the bounded activation checkpoint. Runtime,
schema, fixture, managed-doctrine, workflow, and evidence changes begin only after that checkpoint
validates, commits, and normally pushes.

The user subsequently expanded completion authority through the normal protected producer
PR/merge, doctrine-selected version bump, release validation, tag/assets/catalog/publication, and a
protected Codeheart-HQ installation-validation checkpoint. EP-09 and EP-10 now own those steps; no
second approval is required. The separate HQ semantic plan migration, discovery-v2 activation, and
portfolio-v2 cutover remain excluded, as do destructive Git, protection bypass, history rewriting,
force push, secret exposure, and silent conflict resolution.

## Epic Delta Index

| Epic | Status | Meaningful delta | Review gate |
| --- | --- | --- | --- |
| `EP-01` | completed | Discovery models, pure settings, exclusions, v1/v2 schemas, dispatch, and fresh defaults. | Accepted after two rounds |
| `EP-02` | completed | Shared Git-backed candidate classifier, local index adapter, exact path and marker semantics, ownership boundaries, and family qualification. | Accepted after six rounds |
| `EP-03` | completed | Versioned local views, prospective v2/canonical reads, preview evidence, repository-wide branch touches, and output protection. | Accepted after two rounds |
| `EP-04` | completed | Guarded schema-v2 migration, direct canonical readiness, optional exact mixed proof, atomic renames, and transaction-time authority binding. | Accepted after four rounds |
| `EP-05` | completed | Shared remote tree classification, default-policy branch overlays, canonical-readiness completeness, structured candidate evidence, and v2-only cache replacement. | Accepted after two rounds |
| `EP-06` | pending | None recorded. | Required, including routing probe |
| `EP-07` | pending | None recorded. | Required, including cross-platform evidence |
| `EP-08` | pending | None recorded. | Required, including release-readiness source binding |
| `EP-09` | pending | Doctrine-selected producer release and publication. | Required, including release runbook gates |
| `EP-10` | pending | Verified Codeheart-HQ Kit installation with inactive-v2 proof. | Required, including HQ lifecycle gates |

## Review Gate Metrics

- Review gate required: yes, after every epic.
- Review gate skipped: no.
- Reviewer mode: fresh read-only subagent when the active environment permits it.
- Reviewer model and reasoning mode: inherited from the implementing agent.
- Review rounds: two for EP-01, six for EP-02, two for EP-03, four for EP-04, and two for EP-05.
- Material findings: two EP-01 P1 findings plus EP-02 findings covering unsafe-source ordering,
  path portability, family authority, marker parsing, v1 compatibility, exact object provenance,
  and deterministic candidate evidence; EP-03 findings covered preview context isolation,
  incomplete target-canonical evidence, and filename-owned kind semantics; all fixed.
- EP-04 findings covered v1/v2 authority separation, pre-commit evidence drift, protected register
  handling, unstaged rename reapplication, executable-mode preservation, deterministic final-state
  rebinding, and cross-platform Git mode semantics; all fixed.
- EP-05 findings covered ownership-changing branch candidate visibility and canonical-readiness
  completeness for discovery-v2 members retaining filename-only legacy records; both fixed.
- Files changed because of review include `schemas/kit-config.schema.json`, the classifier, local
  view, migration, commit-tree, and reconcile transaction implementations, their focused tests,
  and JSON-schema coverage.
- Final accepted result: EP-01 through EP-05 accepted; EP-06 through EP-10 pending.
- Approximate added time: about six minutes for EP-01, about thirty-five minutes for EP-02, and
  about twelve minutes for EP-03 review and remediation.
- EP-04 review and remediation added about thirty-five minutes across four rounds.
- EP-05 review and remediation added about twenty minutes across two rounds.
- Token usage: not exposed per review round.
- Worth-it assessment: yes. EP-02 review prevented unsafe local fallback evidence from becoming
  authority, preserved v1 activation compatibility, and bound fallback reads to exact offline Git
  objects despite replacement refs. EP-03 review prevented preview conflicts from being reported
  safe, ensured canonical readiness always carries v2 evidence, and preserved filename-owned kind.
- EP-04 review prevented historical v1 evidence from mutating v2 authority, bound policy and
  candidate evidence inside the reconcile transaction, made protected paths unconditional, and
  prevented false final-state readiness after source/body/mode drift.
- EP-05 review prevented an excluded ownership-changing branch rename from disappearing and
  prevented filename-only legacy records from being omitted by a falsely complete v2 portfolio.

## Activation Delta

- Repository: `Codeheart-Operating-Kit` producer worktree.
- Branch: `codex/operating-kit-multi-root-plan-catalog`, tracking the same-named `origin` branch.
- Activation paths: the canonical implementation plan, this sibling execution log, and
  `docs/repo/plans/README.md` only.
- Safe defaults: keep discovery v1 compatibility, omit child-shape family validity, distinguish
  local nested-repository evidence from remote gitlinks, freeze the first conventional ambiguity
  list, and add Linux semantic CI without Linux distribution assets.
- Exclusions: runtime/source code, schemas, fixtures, managed doctrine/resources, workflows,
  release metadata, consumer repositories, and HQ are absent from the activation checkpoint.
- Activation commit and normal-push result: recorded in the activation report because a commit
  cannot contain its own final object ID; later epic evidence will retain the pushed activation
  revision as its base.

## EP-01 Delta - Discovery Settings And Versioned Contracts

Status: completed and accepted for checkpoint publication.

- Added discovery-version, ownership, signal, Git provenance, deterministic policy, and digest
  models. Non-authoritative untracked preview does not change the authoritative policy digest.
- Added one pure config-byte decoder used by local loading and later remote policy reads. Absence
  remains discovery v1; v2 exclusions are normalized, sorted, and rejected for globs, escapes,
  absolute/drive paths, Unicode/case collisions, and overlaps.
- Retained the planning-workflow extension point while explicitly rejecting `owned_roots`; this
  resolved the first review P1 without invalidating existing schema-v1 extension settings.
- Preserved v1 cache/ledger schemas and made the current schema paths v2 with explicit discovery,
  policy, candidate-set, target, and completeness evidence. Added schema-version dispatch so
  existing v1 portfolio bytes remain historical evidence.
- Fresh Go and Python config writers now emit discovery v2 plus an empty exclusions array. Existing
  config remains byte-preserved during sync.
- Validation: `GOCACHE=/tmp/codeheart-go-cache go test ./internal/plancatalog ./internal/state
  ./internal/commands` passed; `uv run --with pytest pytest tests/test_json_schemas.py
  tests/test_sync_check.py` passed 55 tests; `python3 scripts/validate-json-schemas.py` passed; and
  `git diff --check` passed.
- Review gate: round one found two P1 issues—over-closing the config extension point and host-native
  Windows-drive detection. Both received portable regression coverage. Round two accepted EP-01
  with no remaining actionable findings.

## EP-02 Delta - Shared Git-Backed Candidate Classifier

Status: completed and accepted for checkpoint publication.

- Added one deterministic classifier over neutral Git blob descriptors. It recognizes exact
  lowercase `docs` segments at arbitrary depth, strict discovery/implementation filenames, and
  genuine metadata markers while keeping fenced, quoted, and indented examples inert.
- Added NUL-safe Git-index enumeration plus one bounded streaming `cat-file --batch` reader.
  Fallback reads disable replacement objects and lazy fetch, retain the exact index object ID, and
  never follow an unsafe worktree path. Missing, replaced, symlinked, nested-repository, and other
  unsafe candidate paths remain visible evidence but cannot create local record authority.
- Added hard-unowned managed/local/user, Git mode, containment, path-case, UTF-8, Unicode NFC,
  portable-collision, excluded-root, and conventional-ambiguity behavior with deterministic codes
  and ordering. Ordinary unsafe Markdown without a plan signal does not pollute completeness.
- Added explicit v1/v2 discovery selection for enumerate, discover, snapshot, and formal-path
  classification. V1 retains historical empty-stem filename, source-size, fixed-root, and family
  compatibility; v2 applies strict filename, bounded-content, and exact family rules.
- Added exact `README.md` family qualification requiring a structurally valid header/metadata
  block, metadata kind `family`, and semantic ID kind `family`. Partial metadata preserves identity
  collision evidence only and invalid families cannot satisfy child references.
- Used deterministic temporary-repository fixtures rather than adding a persistent fixture tree;
  this covers root, deep, repeated-docs, malformed, excluded, ambiguous, case, UTF-8, symlink,
  nested-repository, missing-worktree, family, duplicate-ID, and replacement-ref cases without
  shipping test-only catalog documents.
- Validation: `GOCACHE=/tmp/codeheart-plan-catalog-go-cache go test ./internal/plancatalog -count=2`
  passed; `go test ./internal/portfolio ./internal/commands` passed; `go vet` for those three
  packages passed; the repository `plans validate` command reported 35 records and `valid=true`;
  two generated JSON list reads were deterministic in the unit fixture; and `git diff --check`
  passed.
- Review gate: six rounds. Review-driven corrections covered indexed metadata retention behind
  unsafe worktrees, semantic and structural family qualification, candidate-only unsafe errors,
  pre-validation authority filtering, invalid UTF-8 rejection, CommonMark fence details,
  duplicate-ID preservation, v1 compatibility, bounded streaming memory, batch failure handling,
  and exact offline object reads. Round six accepted EP-02 with no actionable findings.

## EP-03 Delta - Local Views, Prospective Reads, And Authoring Preview

Status: completed and accepted for checkpoint publication.

- Added schema-v2 list, validate, and inventory evidence over the same authoritative classifier,
  while keeping the exact schema-v1 list golden unchanged for active discovery v1.
- Added read-only `--target-discovery-version 2` and canonical-readiness lenses. A canonical target
  implies discovery v2 so it cannot emit an incomplete v1 candidate/provenance contract.
- Added explicit non-ignored untracked preview for validate and inventory. Preview records remain
  outside authoritative counts, records, digests, migration inputs, remote evidence, and
  completeness, but are checked against tracked IDs, portable paths, and family context.
- Replaced fixed-root branch ownership for v2 with NUL-safe repository-wide Git diff/tree
  classification, including nested metadata-only additions and both sides of eligibility-changing
  renames. Discovery v1 retains the historical branch-touch route.
- Versioned inventory provenance and coverage, preserved filename/family-owned human kind, and
  protected existing metadata-only plus prospective reserved output targets even when ignored or
  omitted from preview.
- Validation: focused `go test ./internal/plancatalog ./internal/commands ./internal/portfolio`
  passed; remediation regressions repeated two times passed; focused `go vet` passed; active-v1
  JSON golden remained byte-exact; deterministic prospective-v2 list reads matched; and `git diff
  --check` passed. The repository-wide Go suite has only the expected managed content-graph digest
  drift caused by producer changes, reserved for EP-06/EP-08 parity regeneration.
- Review gate: round one found two P1 issues and one P2 issue—preview isolation from tracked
  identity/path/family context, incomplete canonical-target evidence over discovery v1, and
  metadata overriding filename-owned inventory kind. All received regression coverage; round two
  accepted EP-03 with no remaining actionable findings.

## EP-04 Delta - Guarded Migration And Mode Compatibility

Status: completed and accepted for checkpoint publication.

- Kept schema-v1 migration behavior on its historical route and added a hard incompatibility
  blocker when a v1 ledger is presented to an active discovery-v2 repository.
- Implemented schema-v2 migration over the exact prospective inventory revision, policy digest,
  candidate digest, source hashes, owned disposition, branch evidence, dirty state, and target
  preconditions. Every authoritative candidate requires one reviewed disposition.
- Added direct legacy-to-canonical preparation without config writes, plus optional mixed deferral
  only for exact register-proven cutover blobs. Grandfathered plan comparison is byte-exact and no
  longer tolerates checkout line-ending conversion.
- Added deterministic in-place metadata insertion and atomic filename correction as no-replace
  create plus exact-hash remove. Created/Last updated/lifecycle headers, Git executable class, config,
  frozen register, and original cutover revision are preserved.
- Added an authority digest over configured settings/problems, config/register bytes, HEAD,
  policy/candidates/completeness, sources, dirty state, and branch touches. Reconcile revalidates it
  after staging and immediately before the first repository mutation; drift rolls back without
  applying plan actions.
- Added unconditional config/register path protection, including a hostile valid-metadata marker
  in the frozen register, and exact final-state reapplication proof derived from each reviewed
  original blob, semantic decision, target precondition, deterministic inserted bytes, and Git
  executable class. Immediate unstaged rename reapply is a safe no-op.
- Candidate-set hashing now excludes the surrounding selected revision so identical local-index
  and exact commit-tree evidence has the same digest; revision remains separately bound inventory
  evidence. Added one reusable inert exact commit-tree adapter for this proof and later EP-05 use.
- CLI JSON/text now reports versioned projected candidate/canonical/grandfathered/gap readiness and
  explicitly states that migration does not activate discovery or catalog mode.
- Validation: `go test ./internal/plancatalog ./internal/reconcile ./internal/commands -count=1`
  passed; focused v2 migration/CLI regressions repeated two and three times passed; focused `go
  vet` passed; JSON-schema tests passed 40 cases; schema validation and `git diff --check` passed.
  The full Go suite has only the expected managed content-graph digest drift reserved for
  EP-06/EP-08 parity regeneration.
- Review gate: four rounds. The accepted remediation covers v1/v2 separation, transaction-time
  authority revalidation, protected paths, exact mixed proof, deterministic source-to-target
  final-state binding, immediate rename reapply, source mode preservation, and portable
  executable-class comparison with an explicit Windows rule. No P1/P2 findings remain.

## EP-05 Delta - Remote Discovery, Branch Overlays, And Cache Completeness

Status: completed and accepted for checkpoint publication.

- Replaced the fixed `docs/repo/plans` remote scan and duplicated family heuristics with one
  repository-wide selected-tree adapter to the shared classifier. Exact root and arbitrarily deep
  lowercase `docs` segments, filename-or-metadata candidates, exclusions, ambiguity boundaries,
  family metadata, path portability, and canonical errors now match local commit-tree evidence.
- Added NUL-safe repository-wide raw Git change evidence with old/new paths, modes, and object IDs.
  Same-byte renames are suppressed only when candidate eligibility, ownership, signal, and kind
  semantics match; eligibility-, ownership-, and kind-changing renames remain visible.
- Parsed discovery version and exclusions once from default-branch config and applied that immutable
  policy to every branch. Branch config cannot broaden ownership. Discovery-v1 required members are
  retained as explicit incomplete evidence; discovery-v2 members are evaluated through a canonical
  readiness lens even while their local catalog mode remains legacy or mixed.
- Added bounded retained-directory `git cat-file --batch` reads with a shared 8 MiB source limit,
  sanitized inert Git environment, exact object evidence, safe fallback for injected runners, and
  post-stream mirror-authority verification. Symlinks and mode `160000` gitlinks never contribute
  content authority.
- Added schema-v2 plan-candidate observations carrying repository, ref, commit, path, signal,
  ownership, filename kind, Git mode/object, content and policy hashes, exclusion/boundary evidence,
  visibility, and verification. This preserves excluded and otherwise non-record candidate evidence
  without granting stable-ID plan authority or making a reviewed exclusion incomplete.
- Portfolio catalogs now publish discovery-v2 member policy/candidate digests and member completeness.
  Only complete schema-v2/discovery-v2 evidence can replace the cache. Incomplete attempts label and
  preserve the last complete v2 cache; historical v1 cache bytes remain historical and cannot supply
  a current `last_complete_scan_at`.
- Removed the temporary CLI conflict between remote overlays and prospective/active discovery v2.
  Added parity, nested docs, metadata-only, family, default-policy, excluded rename, eligibility/kind
  rename, filename-only legacy, v1 member, gitlink, hostile/oversized content, cache, inert transport,
  deterministic sorting, and public-evidence tests.
- Validation: `go test ./internal/portfolio ./internal/commands ./internal/plancatalog` passed;
  `go test -race ./internal/portfolio` passed; JSON-schema tests passed 40 cases; focused `go vet`
  passed; Windows amd64 portfolio and plan-catalog test binaries compiled; `git diff --check` passed.
  The full Go suite has only the expected managed content-graph digest drift reserved for
  EP-06/EP-08 parity regeneration.
- Review gate: accepted after two rounds. The two blocking findings added structured evidence for
  ownership-changing excluded renames and required canonical-valid evidence before a discovery-v2
  portfolio member/cache can be complete. No P1/P2 findings remain.

## EP-06 Delta - Managed Doctrine, Routing, And Resource Parity

Status: pending.

No divergence, validation, routing-probe, or review evidence is recorded yet.

## EP-07 Delta - Comprehensive Validation And Cross-Platform Proof

Status: pending.

No divergence, validation, platform, or review evidence is recorded yet.

## EP-08 Delta - Release Readiness And Consumer Handoff Evidence

Status: pending.

No divergence, readiness, handoff, or review evidence is recorded yet.

## EP-09 Delta - Versioned Producer Release And Publication

Status: pending.

User authority is durable; execution waits for accepted EP-08 readiness evidence.

## EP-10 Delta - Codeheart-HQ Installation And Validation

Status: pending.

User authority is durable; execution waits for a verified EP-09 published release. HQ semantic
catalog migration and discovery-v2 activation remain outside this epic.

## Final Validation

Pending completion of EP-01 through EP-10. No release, publication, or consumer-installation
evidence may be inferred before its owning epic records the exact validated result.
