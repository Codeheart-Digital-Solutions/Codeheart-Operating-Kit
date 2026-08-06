Last updated: 2026-08-06T22:34:04Z (UTC)
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
| `EP-03` | pending | None recorded. | Required |
| `EP-04` | pending | None recorded. | Required |
| `EP-05` | pending | None recorded. | Required |
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
- Review rounds: two for EP-01 and six for EP-02.
- Material findings: two EP-01 P1 findings plus EP-02 findings covering unsafe-source ordering,
  path portability, family authority, marker parsing, v1 compatibility, exact object provenance,
  and deterministic candidate evidence; all fixed.
- Files changed because of review: `schemas/kit-config.schema.json`, `internal/plancatalog/view.go`,
  `internal/plancatalog/plancatalog_test.go`, and `tests/test_json_schemas.py`.
- Final accepted result: EP-01 and EP-02 accepted; EP-03 through EP-10 pending.
- Approximate added time: about six minutes for EP-01 and about thirty-five minutes for EP-02
  review and remediation.
- Token usage: not exposed per review round.
- Worth-it assessment: yes. EP-02 review prevented unsafe local fallback evidence from becoming
  authority, preserved v1 activation compatibility, and bound fallback reads to exact offline Git
  objects despite replacement refs.

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

Status: pending.

No divergence, validation, or review evidence is recorded yet.

## EP-04 Delta - Guarded Migration And Mode Compatibility

Status: pending.

No divergence, validation, or review evidence is recorded yet.

## EP-05 Delta - Remote Discovery, Branch Overlays, And Cache Completeness

Status: pending.

No divergence, validation, or review evidence is recorded yet.

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
