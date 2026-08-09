Last updated: 2026-08-09T16:34:26Z (UTC)

# Branch Ownership Reconciliation Acceptance Evidence

Status: implementation validation in progress; final platform and publication evidence pending.

## Sanitized Scenario Proof

The Go acceptance fixture creates only synthetic public-safe repositories and plan identifiers. It
proves the two motivating evidence distributions without copying consumer paths, refs, commit IDs,
or operational records:

| Scenario | Expected classification | Focused result |
| --- | --- | --- |
| Twelve independent same-content refs | 12 `same-content-non-owner` | passed |
| Six divergent refs with exactly incorporated complete transitions | 6 `incorporated-history-non-owner` | passed |
| One divergent active candidate | 1 `active-owner` | passed |
| Second independent active candidate | 2 active owners across two plan candidates | passed |

The fixture is `TestBranchEvidenceSanitizedAcceptanceCohorts` in
`internal/plancatalog/branch_reconciliation_test.go`. It uses real Git commits, refs, blob modes,
content hashes, merge bases, transition digests, and stable patch IDs.

## Focused Evidence And Security Matrix

Passing focused suites cover:

- local and tracking same-content refs, plus independently created remote same-content history;
- exact squash/incorporated-history proof, ambiguous or missing incorporated commits, partial
  transition and normalization collisions, rename/delete transitions, binary and invalid UTF-8
  content, symlink modes, custom attributes, hostile Git config/environment, and digest tampering;
- ambiguous merge bases and proof-scope enumeration failure without candidate-only fallback;
- moved refs, moved merge bases, changed path state, dirty/uncommitted ownership, and proof/config/
  ledger/action tampering;
- genuine mixed active-owner deferral, direct-canonical rejection, separate activation, later exact
  owner integration, canonical-metadata tampering, and post-write rollback;
- matching tracking/mirror coalescing and blocking disagreement;
- remote checkpoint validation at `E`, `L`, and `A`, an unrelated descendant after `A`, and
  fail-closed config/register/candidate byte or mode drift;
- v0.1.25/v0.1.26 old-CLI rejection of schema-v3 ledger and config-v2 activation artifacts with
  zero repository writes.

## Current Command Evidence

The following implementation-focused validation has passed on macOS during the active worktree
phase:

```text
go test ./internal/plancatalog ./internal/portfolio ./internal/commands ./internal/cli -count=1
go test -race ./internal/plancatalog -run '^TestBranchEvidence' -count=1
python -m pytest tests/test_plan_catalog_backward_compatibility.py
python -m pytest tests/test_json_schemas.py tests/test_packaging_resources.py tests/test_routing.py tests/test_sync_check.py
```

The exact release-source full-suite, Ubuntu/macOS/Windows CI, repeated artifact-build digests, and
live publication results will be added to `release-readiness-evidence.md` after those gates finish.

## Residual Limits

- Merged-PR status is optional corroboration and is not an offline clearance dependency.
- Remote-aware evidence is accepted only from a complete fresh overlay; prior cache evidence cannot
  authorize migration.
- The bounded incorporated-commit search stops and blocks when its admitted search set is exceeded
  or ambiguous.
- No historical branch deletion is required or recommended by this feature.
