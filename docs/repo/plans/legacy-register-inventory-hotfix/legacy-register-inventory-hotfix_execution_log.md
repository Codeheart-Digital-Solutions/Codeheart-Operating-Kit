Last updated: 2026-08-10T06:04:50Z (UTC)
Created: 2026-08-10

# Legacy Register Inventory Hotfix Execution Log

Plan: `legacy-register-inventory-hotfix_implementation_doc.md`
Mode: goal-style implementation and authorized producer release
Status: active
Overall divergence: none at activation

## Summary

The user authorized the focused implementation and full v0.1.28 producer release workflow. The
work starts from exact live `origin/main` and v0.1.27 merge
`70a8e174fe8767b09eb6e0ccbea64f6fad5e9a3f` on the distinct branch
`codex/legacy-register-inventory-v0.1.28-hotfix`.

The worktree was clean and detached before branch creation. Git identity is configured. No
transaction, recovery, staging, or quarantine state was present. The separate stalled design
branch/worktree remains untouched. The frozen register and all consumer repositories are excluded.

## Epic Delta Index

| Epic | Status | Meaningful delta | Review gate |
| --- | --- | --- | --- |
| `EP-01` | completed | Split raw observations from strict wire projection; canonical lifecycle, typed legacy status, and separate dates now fail closed. | Focused parser/schema review and final independent gate passed. |
| `EP-02` | completed | Added generic positive/negative fixtures plus v1 compatibility, schema-v3 mode parity, deterministic digest/bytes, ambiguity, and zero-write coverage. | Focused and full validation passed; every High/Medium review finding is resolved. |
| `EP-03` | pending | None yet. | Required before publication. |

## Review Gate Metrics

- Review gate required: yes.
- Review gate skipped: no.
- Reviewer mode: fresh read-only reviewer agent when implementation is ready.
- Review rounds: three read-only passes across two fresh reviewer agents.
- Material findings: zero High; seven Medium and one Low found and resolved.
- Final accepted result: no remaining High or Medium release blocker.
- Worth-it assessment: accepted; review closed schema-invalid relation/path edges, ambiguous
  reconciliation authority, and schema-v1 compatibility drift before release preparation.

## Activation Delta

- Repository: `Codeheart-Digital-Solutions/Codeheart-Operating-Kit`.
- Live default branch: `main` at `70a8e174fe8767b09eb6e0ccbea64f6fad5e9a3f`.
- Released base: annotated `v0.1.27` peels to the same commit.
- Local identity: configured public maintainer name and repository-authorized email.
- GitHub connector identity: repository administrator with push and maintain permission.
- Activation checkpoint scope: this implementation plan, sibling execution log, and the nearest
  plan index only.
- External note: local `gh` credentials were invalid at preflight, while authenticated `git` fetch
  and the GitHub connector succeeded. Publication will recheck both routes before the live gate.

## EP-01 Delta - Typed Parsing And Strict Projection

Raw legacy-register observations now retain repeated historical fields without JSON/YAML wire
tags. A separate projector emits only canonical lifecycle values or the closed
`implementation-handoff-ready` evidence value, never both. `Completed` is a field boundary and a
separate calendar date; `Last updated` remains a second-precision UTC timestamp. Invalid or
ambiguous values are omitted from typed wire fields and retained as stable error problems.

## EP-02 Delta - Focused Regression Evidence

Synthetic fixtures cover all canonical lifecycles, the recognized pre-canonical status, separate
completion/update dates, combined date text, duplicate fields, invalid dates and paths,
unsupported status labels, and two entries claiming one canonical path. Focused tests prove
schema-v3 validation in configured local, prospective canonical-target, and attached remote-aware
modes, repeat-build byte/digest determinism, and byte preservation for the frozen register and
formal plans. Affected `plancatalog`, `commands`, and `portfolio` packages plus the JSON-schema
validator and 42 Python schema tests pass.

The compatibility adapter reconstructs the exact v0.1.27 schema-v1 legacy wire projection from
the raw observation while schema-v3 remains strict. Repeated or malformed relation blocks,
canonical-doc fields, canonical-doc values, backslash paths, and rooted paths now emit structured
blockers and project no ambiguous reconciliation authority. Focused cross-platform assertions
retain historical schema-v1 path semantics.

Independent review required three repair passes. The first pass found malformed relation wire,
duplicate canonical-doc field authority, and schema-v1 projection drift. The second found
platform-dependent backslash acceptance, duplicate path-value authority, and missing duplicate
relation proof. The final follow-up found and closed a Windows-specific assertion and a rooted-path
strictness edge. The concluding read-only review reported no remaining High or Medium blocker.

Validation completed before release versioning:

- focused `plancatalog`, `commands`, and `portfolio` Go suites;
- repository-wide Go tests, race tests, and `go vet`;
- 158 Python tests, including 42 JSON-schema tests;
- public-core, Markdown, JSON-schema, and release-manifest validators;
- producer inventory validation and deterministic local/canonical-target/remote-aware evidence;
- `git diff --check` after the final reviewer repairs.

## Final Validation

EP-01 and EP-02 are complete. Exact v0.1.28 release-source validation, reproducible assets,
installer/upgrade proof, publication, and live verification remain under EP-03.
