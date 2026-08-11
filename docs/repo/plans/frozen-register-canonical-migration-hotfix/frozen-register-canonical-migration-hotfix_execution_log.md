Last updated: 2026-08-11T14:45:48Z (UTC)
Created: 2026-08-11

# Frozen Register Canonical Migration Hotfix Execution Log

Plan: `frozen-register-canonical-migration-hotfix_implementation_doc.md`
Mode: focused v0.1.30 producer correction, release, Foundry migration, and HQ refresh
Status: active
Overall divergence: none at activation

## Summary

Foundry's exact reviewed schema-v3 migration dry-run contains 90 owned records, 78 metadata
actions, 12 already-canonical skips, complete branch dispositions, and no identity, alias, hash, or
candidate-set ambiguity. It is blocked only by 29 frozen-register parser errors: 13 malformed
relations, eight unsupported relation kinds, and eight unsupported historical status labels.

The register cannot be normalized after mixed cutover. The correction will preserve those
observations but make them warnings only when the effective catalog mode is canonical and the
problem path is the frozen register. Legacy and mixed authority and all other errors remain strict.

## Epic Delta Index

| Epic | Status | Meaningful delta | Review gate |
| --- | --- | --- | --- |
| `EP-01` | completed | Exact blocker and narrow behavioral contract established; this commit is the explicit three-path plan checkpoint. | Accepted with exact baseline parity. |
| `EP-02` | completed | Fixed allowlist projection, strict negative matrix, managed doctrine, and packaged mirrors complete. | Accepted; no High or Medium finding. |
| `EP-03` | in progress | v0.1.30 surfaces, full local matrix, reproducible packs, and isolated macOS install complete; PR/CI/publication pending. | Exact-head CI and live publication required. |
| `EP-04` | pending | Foundry remains paused at its reviewed v0.1.29 ledger checkpoint; HQ refresh not started. | Fresh post-upgrade evidence and consumer byte-preservation gates. |

## Activation Evidence

- Producer base/tag/main: `996036893faa251e7f38c3391f94efe7f4d1273d` / `v0.1.29`.
- Activation branch: `codex/frozen-register-canonical-migration-v0.1.30`.
- Candidate next tag: `v0.1.30`; remote tag inventory has no collision at activation.
- Affected public subsystem: `internal/plancatalog` plus managed migration/reference doctrine,
  packaged mirrors, version/release surfaces, and tests.
- Consumer impact: validator-only correction plus instruction-only clarification; Kit upgrades do
  not themselves migrate plan metadata or alter frozen registers.
- Review mode: strongest independent main-thread review because team delegation is unavailable in
  the current execution environment; no review requirement is waived.

## Foundry Blocker Evidence

The paused Foundry dry-run at its reviewed ledger checkpoint reports:

- 94 discovered candidates;
- 90 owned/formal candidates;
- 78 reviewed metadata changes;
- 12 exact already-applied skips;
- zero branch observations or unresolved ownership blockers;
- zero projected canonical gaps;
- 29 errors, all at `docs/repo/plans/plan-register.md` and all within the three accepted
  non-authoritative historical-field codes.

No Foundry metadata/config write or activation occurred. The existing ledger will be discarded as
stale after Kit upgrade and rebuilt from a fresh evidence revision before any consumer migration
write.

## EP-01 Delta - Plan Activation

The new canonical plan is enumerated once as active with its expected semantic ID and path.
Markdown timestamp, public-core, and diff checks pass. Repository plan validation contains the same
pre-existing frozen-register errors as exact v0.1.29 plus only the expected compatibility-layout
warning for this newly authored structured plan; it introduces no new error. The checkpoint stages
only this plan, this log, and `docs/repo/plans/README.md`.

## Review Gate Metrics

- Review gate required: yes.
- Review gate skipped: no.
- Review rounds completed: one fresh main-thread review after the complete local matrix.
- High findings: zero.
- Medium findings: zero.
- Final accepted result: accepted for PR and exact-head CI.

## EP-02 Delta - Contextual Severity Correction

`LoadRepositorySnapshotWithOptions` now projects a copied problem slice after the effective target
mode is selected. It changes an error to a warning only when the mode is canonical, the path is the
exact frozen register, and the code is one of `legacy_status_unsupported`,
`legacy_relation_malformed`, or `legacy_relation_unsupported`. The parser itself remains strict,
so legacy and mixed snapshots still expose the original errors. Duplicate relations, malformed or
duplicate canonical paths and fields, dates, identity, readability, safety, reconciliation, and
all non-register paths remain errors.

The integration regression creates a valid frozen mixed baseline containing exactly the three
historical-field forms. It proves mixed errors, target-canonical warnings, a ready reviewed
canonical migration, configured canonical warnings, input immutability, and exact register-byte
preservation. Unit coverage proves wrong-path and authority-critical negatives. Managed migration
and register references state the same bounded rule, and source/resource copies are byte-identical.

Fresh review traced every caller through repository snapshot, inventory, migration, validation,
command-side, and portfolio paths. The change does not extend the migration remediable-error
allowlist, touch ledger parsing, alter candidate or branch authority, change transaction behavior,
or bypass checkpoint/config/digest guards. No High or Medium finding remains.

## EP-03 Delta - Release Candidate

The release identity is v0.1.30 across Go/Python versions, bootstrap, macOS and Windows installers,
planning-workflows component, standard profile, root/resource manifests, release-candidate fixture,
upgrade/parity tests, and release notes. Component/profile checksums and the standard graph digest
were recomputed, and every packaged mirror is identical to its source.

Validation results:

- focused frozen-register regression group: pass;
- complete `internal/plancatalog` package: pass (`148.637s` before version surfaces and `157.462s`
  in the full repository run);
- `go test -timeout 30m ./...`: pass;
- `go test -race -timeout 30m ./...`: pass;
- `go vet ./...`: pass;
- Python/schema/resource/routing/release suite: 156 passed plus two historical-tag setup errors in
  the shallow clone; after fetching immutable v0.1.25/v0.1.26 tags, both affected cases passed,
  yielding 158/158 effective passes;
- JSON schemas, release manifests, Markdown timestamps, public-core hygiene, packaged-resource
  parity, gofmt, and `git diff --check`: pass;
- two separate full asset builds: byte-identical; each builder also repeated both platform builds
  internally and reported `reproducible: true`;
- isolated local-catalog macOS install: pass; CLI reports 0.1.30 and the binary contains x86_64 and
  arm64 slices;
- Windows artifact identity: PE32+ x86-64 with the exact pack-manifest binary digest.

Final retained release-candidate identity:

- macOS universal archive SHA-256:
  `d23df9463bc50f1e8749e893839d03bd5283d5b20a9c0a0f0069e3456ee0ab9d`;
- macOS pack-manifest SHA-256:
  `8cf501bb504bf454e7098f84403954c8e6b5c100c05938b97ad403ba5cf08492`;
- macOS binary SHA-256:
  `4c5e7f047430438104a091df914be8732431d43d5d144f716e44f36d3cc90384`;
- Windows x64 archive SHA-256:
  `e4df78f6b942ffd8f6cbd7e9ec2c969ef201c2da3da895342487b7d14a288827`;
- Windows pack-manifest SHA-256:
  `a620c58a13233c592f8f96c862fa8e1863f19ced7056cc390e5d3ba497e35d3c`;
- Windows binary SHA-256:
  `342ba605b89579927e940735b36b8e0af7cbe9c3c3dfa93961c17c08c1343ee1`;
- embedded content-manifest SHA-256:
  `8534f13155ff47d8a4a04a9c960b5c2aeaf4f0ab9de2363ca029db9ae20344a6`;
- external release-catalog SHA-256:
  `c8a3a9faff903767b11b0a0dce2af07d16c1d472dce9c55bf1c49086cade2e03`.

The candidate v0.1.30 binary was also run against the exact paused Foundry ledger checkpoint. Its
dry-run reports 90/90 canonical coverage, 78 replacements, 12 already-applied skips, zero blockers,
`ready: true`, `canonical_ready: true`, and the unchanged migration action digest
`c30f0a52ffaad4c74d374e1360c7fddb210d6183de04bc6fc8ec29f7eade631f`.
The Foundry worktree remained clean and no migration or activation write occurred.

## Divergence Log

- The first release-asset test invocation exposed two incomplete version-bump fixtures: the
  negative wrong-target value had become the real v0.1.30 target, and the profile hash still
  expected v0.1.29 bytes. Both were corrected and their focused and full suites passed.
- The first full Python run used a narrow clone without v0.1.25/v0.1.26 tags. The 156 executable
  tests passed; the two archive fixtures passed unchanged after fetching those immutable tags.
- A test-only `uv.lock` generated by the environment runner was removed before staging. It was not
  product or user work and is absent from the candidate diff.
