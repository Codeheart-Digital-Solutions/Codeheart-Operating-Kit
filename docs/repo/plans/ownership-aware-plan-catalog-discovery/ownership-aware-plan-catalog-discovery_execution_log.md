Last updated: 2026-08-06T21:41:02Z (UTC)
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

The user explicitly activated the eight-epic producer implementation on 2026-08-06. Execution uses
`codex/operating-kit-multi-root-plan-catalog` and follows EP-01 through EP-08 linearly. The plan,
this execution log, and the nearest plans index form the bounded activation checkpoint. Runtime,
schema, fixture, managed-doctrine, workflow, and evidence changes begin only after that checkpoint
validates, commits, and normally pushes.

The implementation remains producer-only. Pull-request creation, merge, release, tag, consumer
upgrade, Codeheart-HQ modification, destructive Git, history rewriting, force push, and public
repository governance remain outside current authority.

## Epic Delta Index

| Epic | Status | Meaningful delta | Review gate |
| --- | --- | --- | --- |
| `EP-01` | completed | Discovery models, pure settings, exclusions, v1/v2 schemas, dispatch, and fresh defaults. | Accepted after two rounds |
| `EP-02` | pending | None recorded. | Required |
| `EP-03` | pending | None recorded. | Required |
| `EP-04` | pending | None recorded. | Required |
| `EP-05` | pending | None recorded. | Required |
| `EP-06` | pending | None recorded. | Required, including routing probe |
| `EP-07` | pending | None recorded. | Required, including cross-platform evidence |
| `EP-08` | pending | None recorded. | Required, including release-readiness stop boundary |

## Review Gate Metrics

- Review gate required: yes, after every epic.
- Review gate skipped: no.
- Reviewer mode: fresh read-only subagent when the active environment permits it.
- Reviewer model and reasoning mode: inherited from the implementing agent.
- Review rounds: two for EP-01.
- Material findings: two EP-01 P1 findings in round one; both fixed.
- Files changed because of review: `schemas/kit-config.schema.json`, `internal/plancatalog/view.go`,
  `internal/plancatalog/plancatalog_test.go`, and `tests/test_json_schemas.py`.
- Final accepted result: EP-01 accepted; EP-02 through EP-08 pending.
- Approximate added time: about six minutes for EP-01 review and remediation.
- Token usage: not exposed per review round.
- Worth-it assessment: pending.

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

Status: pending.

No divergence, validation, or review evidence is recorded yet.

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

## Final Validation

Pending completion of EP-01 through EP-08. No release, publication, or consumer evidence may be
inferred from an incomplete execution log.
