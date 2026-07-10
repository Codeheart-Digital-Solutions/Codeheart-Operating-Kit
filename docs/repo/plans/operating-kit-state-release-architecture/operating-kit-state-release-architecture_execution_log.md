Last updated: 2026-07-10T00:54:31Z (UTC)
Created: 2026-07-09

# Operating Kit State And Release Architecture Execution Log

Plan: `operating-kit-state-release-architecture_implementation_doc.md`

Mode: goal-style source implementation

Status: completed; source and validation-only macOS/Windows CI passed

## Overall Divergence

The simplified five-epic plan remains the execution authority. Three bounded execution
differences are recorded:

- focused routing probes were implemented as fresh deterministic producer and materialized-
  consumer scenarios because subagent execution was not permitted;
- the retained Python sync path preserves existing external release provenance because embedded
  content identity no longer contains archive assets; Go lifecycle commands are authoritative;
- source and local validation completed before the scoped branch/commit/push route was authorized;
  after authorization, the configured real Windows workflow ran and passed the completion gate.

## Summary

The plan was activated after explicit user approval. Execution is limited to source implementation,
local validation, validation-only CI, planning lifecycle updates, and evidence. No release, tag,
named consumer sync, local Operating Kit update, signing change, Python retirement, or new platform
work is authorized.

Consumer impact:

- consumer migration required;
- validator-only change;
- instruction-only change;
- security or safety policy change;
- no breaking placement-contract change.

## Epic Delta Index

| Epic | State | Meaningful divergence |
| --- | --- | --- |
| `EP-01` | completed | Shared profile-level state defaults replaced repetitive component-level strategy fields; `resources.go` already embedded the schema glob and required no edit. |
| `EP-02` | completed | Existing Python parity tests remain compatibility evidence, but lock-v2 graph records and same-version sync provenance are intentionally Go-authoritative. |
| `EP-03` | completed | Real Windows install/upgrade execution passed through the EP-05 completion gate. |
| `EP-04` | completed | Deterministic focused probes replaced subagent probes because delegation was not permitted. |
| `EP-05` | completed | Local evidence and validation-only macOS/Windows CI passed. |

## Review Gate Metrics

Review gate required: yes

Review gate skipped: no

Reviewer mode: strongest practical main-thread review; reviewer-agent execution is not permitted
unless the user or applicable repository instructions explicitly request subagent work.

Review rounds: 8 during implementation and CI

Material findings: EP-01 state precedence, EP-02 pre-commit cleanup, EP-03 pack-verifier checksum
mapping, final Windows expected-failure messaging, and missing explicit Windows junction coverage
plus bare-filename catalog inference, catalog/archive containment, and stale final graph identity
plus no-Python installer dependency, Windows drive-path parsing, and symlink-scanner pipefail
handling, host-native archive ordering, and expected Windows native-exit handling found and fixed

Final accepted result: EP-01 through EP-05 accepted; validation-only macOS and Windows execution
passed

Worth-it assessment: positive. The capability is concentrated in three internal packages and one
new durable lifecycle runbook. It adds no service, registry, successful-operation journal,
evidence attachment set, or extra approval step for routine repair and same-version sync.

## EP-01 Delta - Validated State Model And Migration

State: completed

Safe defaults and divergence:

- Kept per-file strategy overrides available in the component schema while defining shared defaults
  once in `profiles/standard.yaml`; this avoids repeating the same policies across every component
  file without narrowing graph behavior.
- Left `resources.go` unchanged because its existing `schemas` embed pattern automatically includes
  the new schemas.
- Reused the existing ignored `.venv` after the shell `python3` lacked pytest; no tooling was
  installed or repaired.

Validation:

- `go test ./...`: passed after updating the intentional profile digest fixture.
- State, component, lockfile, and manifest package tests: passed.
- State fixtures cover numeric-looking strings, optional config, legacy placeholders, unrelated
  invalid v1 data, future versions, graph semantics, all primary classifications, and migration.
- `.venv/bin/python -m pytest tests/test_json_schemas.py tests/test_packaging_resources.py -q`:
  `31 passed`.
- `python3 scripts/validate-json-schemas.py`: passed.
- `go mod verify`: passed.

Review gate:

- Reviewer mode: strongest practical main-thread review.
- Round 1 finding: unsupported future lock state could be classified as `partial` before the lock
  version was inspected when another required surface was missing.
- Fix: inspect and reject unsupported or invalid lock versions before partial-state evaluation when
  a lock exists; added complete and incomplete future-lock regression fixtures.
- Material findings after fix: none.
- Files changed because of review: `internal/state/observed.go` and `internal/state/state_test.go`.
- EP-01 result: accepted.

## EP-02 Delta - Transactional Lifecycle Commands

State: completed

Safe defaults and divergence:

- Kept legacy JSON fields as top-level compatibility projections and placed the shared concise
  operation result under `result`, avoiding a second public schema for an internal result type.
- Kept `--release-manifest` on `sync` as validation-only compatibility input. Same-version sync now
  preserves installed release provenance and uses only embedded running-binary content.
- Retained the Python CLI as parity evidence. Its lock-v1 surface remains intentionally narrower;
  Go lock v2 records the complete typed graph.
- Real Windows reparse behavior was retained as an EP-05 platform gate and passed in the final
  validation-only workflow.

Validation:

- `go test ./...`: passed.
- `go test -race ./internal/reconcile ./internal/commands ./internal/cli ./internal/drift`:
  passed.
- `GOOS=windows GOARCH=amd64 go build ./...`: passed.
- `.venv/bin/python -m pytest tests/test_cli.py tests/test_init.py tests/test_onboard.py
  tests/test_sync_check.py tests/test_update_check.py tests/test_go_cli_parity.py -q`:
  `51 passed`.
- State-transition tests cover absent, adoptable, current, drifted, stale CLI, partial, invalid,
  legacy, active transaction, recovery, and future-version inputs.
- Preservation and failure tests cover config, root instructions, plans, memory, local-user data,
  safe retirement, managed-section repair, concurrency, parent replacement, staged validation,
  interrupted commit, successful rollback, failed rollback, stale recovery, no-op, and cleanup.

Review gate:

- Reviewer mode: strongest practical main-thread review.
- Round 2 finding: several filesystem errors after marker acquisition but before commit returned
  directly, which could remove the marker while retaining partial stage data.
- Fix: route every pre-commit setup, staging, marker-write, and validation failure through the same
  cleanup rollback path; marker acquisition now also removes a partially written marker.
- Additional tightening: require plan/marker identity agreement before stale staged takeover,
  revalidate parent identity before each replacement, and compare resolved prefixes to reject
  symbolic links and Windows reparse targets.
- Material findings after fix: none.
- EP-02 result: accepted.

## EP-03 Delta - Reproducible Release And Safe Upgrade

State: completed; real Windows execution passed under EP-05

Implementation:

- Split embedded content identity from external release catalog and per-platform pack identity;
  added strict schemas and catalog-to-binary verification.
- Added deterministic macOS universal and Windows x64 pack builds, normalized archive metadata,
  repeat-build byte comparison, and final catalog emission.
- Added safe extraction, traversal and link rejection, payload and identity validation, forward-
  only upgrade selection, hash-bound handoff, new-binary reconciliation, final check, and binary/
  state restoration after failure.
- Added public `upgrade --dry-run` and `upgrade --yes`; kept internal handoff and verification
  commands absent from help.
- Refactored both installers to validate catalog, archive, pack, payload, content, binary, and
  staged version before replacement.

Validation:

- `go test -race ./...`: passed, including release, handoff, rollback, traversal, schema,
  transaction, and command tests.
- `.venv/bin/python -m pytest tests/test_install_metadata.py tests/test_release_assets.py
  tests/test_packaging_resources.py -q`: `27 passed` after the final catalog/archive hardening.
- Final release builder: macOS universal and Windows x64 packs each built twice, final bytes
  matched, and the catalog/pack/content/binary chain verified.
- Isolated macOS install: passed. A conflicting checksum failed closed and the installed binary
  hash remained unchanged.
- Isolated upgrade: dry-run preserved the old binary and lock hashes; approved handoff upgraded
  `0.1.20` to `0.1.21`; the resulting installation reported `current` with no drift.
- `GOOS=windows GOARCH=amd64 go build ./...`: passed. The configured Windows workflow subsequently
  passed actual installation, deferred handoff, reparse containment, and failure preservation.

Review gate:

- The first EP-03 review found that the Python pack verifier interpreted checksum/path columns in
  reverse order. The mapping was corrected and release tests were rerun.
- Handoff validation was tightened to require a sibling handoff directory, staged binary
  containment, and a hash-bound handoff file; Windows cleanup uses the staged binary after parent
  exit.
- No remaining material source defect was found. EP-03 source result: accepted.

## EP-04 Delta - Compact Lifecycle Guidance And Producer Routing

State: completed

Implementation:

- Added one compact managed lifecycle runbook with state-to-command selection, dry-runs, upgrade
  approval, blockers, recovery, tooling readiness, and stop conditions.
- Added producer authority to root instructions and aligned onboarding, update-check, managed
  routes, the root managed template, maintainer runbooks, and retained Python resource mirrors.
- Kept routine authorization minimal: invoking init, repair, or same-version sync authorizes that
  operation; only a version-changing upgrade requires `--yes`.

Focused routing probes:

| Prompt | Selected route | Approval class | Result |
| --- | --- | --- | --- |
| Change Operating Kit implementation source. | Tracked `components/`, `internal/`, schemas, templates, and `docs/repo/`; never ignored `.codeheart/kit/`. | Existing source-implementation authorization. | pass |
| Restore an existing compatible drifted materialized consumer without changing version. | `repair --dry-run`, then `repair`; `sync` remains same-version and `upgrade` is version-only. | Named repair invocation authorizes repair; no upgrade approval. | pass |

Validation:

- `.venv/bin/python -m pytest tests/test_onboard.py tests/test_packaging_resources.py
  tests/test_routing.py -q`: `15 passed`.
- Lifecycle preservation tests prove consumer config, root instructions, plans, memory, local-user
  content, and managed-section-local content remain intact.
- Markdown and public-core validators passed; source and retained resource mirrors match.
- EP-04 result: accepted.

## EP-05 Delta - Integrated Validation And Source Handoff

State: completed; local source and validation-only macOS/Windows CI passed

Local integrated evidence:

- `go test -race ./...`: passed.
- `.venv/bin/python -m pytest -q`: `130 passed`.
- `go vet ./...`: passed.
- JSON Schema, release-contract, public-core, and Markdown validators: passed.
- Shell syntax and PowerShell help parsing: passed. `actionlint` and PyYAML are not present, so no
  new local tooling was installed; the workflow is retained for GitHub-native parsing and runtime.
- Final reproducibility, macOS install, upgrade dry-run/apply, preservation, and post-check
  evidence passed in isolated temporary roots; repository `.codeheart/` state and the developer
  CLI were not used or changed.

Final review:

- Capability, consumer preservation, transaction safety, release integrity, platform symmetry,
  routing, and publication stops were reviewed on the main thread.
- Finding: Windows checksum-failure lanes expected the phrase `Checksum mismatch`, while the new
  fail-closed catalog conflict used a different phrase. Both installers now use one matching
  message; installer/release tests (`25 passed`), Bash parsing, PowerShell parsing, and repeat-build
  verification passed after the fix.
- Finding: the general reparse test skips on Windows but the workflow had no replacement fixture.
  The Windows install lane now creates a real directory junction, requires repair to fail, and
  verifies the outside sentinel is preserved. The final Windows workflow passed this runtime
  evidence.
- Finding: PowerShell catalog inference passed an empty parent to `Join-Path` when `-AssetFile`
  was a bare filename, as used by the public-release smoke lane. It now resolves that case against
  the current directory; the focused PowerShell probe and affected installer/release suite passed.
- Finding: catalog asset names were not bound to the requested version/platform filename, and ZIP
  extraction relied too heavily on declared sizes and incomplete filesystem-type checks. Catalog
  schema/runtime selection and both installers now reject unsafe names; Go extraction enforces
  entry, declared-size, actual-size, symlink, and special-file limits; deterministic packs encode
  regular-file types. Negative catalog and symbolic-link fixtures pass.
- The handoff rollback fixture now uses the native Go test executable as the staged candidate and
  forces only its child reconciliation process to fail. On Windows CI this exercises native
  executable replacement, child-process failure, file-lock release, and byte-for-byte restoration
  rather than relying on a Unix script fixture.
- No remaining critical or high-severity source finding is known. The required real Windows
  runtime evidence passed before plan completion.

CI round 1:

- Push-triggered Validate run `29059953879` tested commit
  `2200e960c4f14e736ade188aaa1e572e5bc62cbb`.
- Both platform jobs stopped at `go test ./...`: the final removal of an extra EOF blank line from
  the managed lifecycle runbook changed the compiled graph digest after content identity had been
  generated. Recorded `e534aba6...` did not match compiled `8acfeed3...`.
- Root and retained-resource content manifests now record the final graph digest. Local Go,
  packaged-resource, release-manifest, and schema validation pass after the fix. A new pushed CI
  run remains required.

CI round 2:

- Push-triggered Validate run `29060094718` tested commit
  `a0b98f018e3dc9fd1c6d90133accf26ce27403d5`.
- macOS passed Go, Python parity, installer/release tests, public-core, Markdown, schemas, release
  validation, and reproducible builds, then the minimal no-Python installer lane found unqualified
  `dirname` calls. Those calls now use `/usr/bin/dirname`; the exact minimal-PATH offline install
  passes locally.
- Windows passed setup and reached Go release tests, where absolute `D:\...` local asset paths
  were reparsed as URL scheme `d`. Windows absolute filesystem paths are now read before URL
  parsing; Windows-target tests compile and focused release tests pass locally.
- The affected full installer suite then exposed `grep -q`/`pipefail` SIGPIPE ambiguity in the
  symlink scanner. The scanner now consumes the complete ZIP listing; all 14 installer tests pass.
- A new pushed CI run remains required for these fixes.

CI round 3:

- Push-triggered Validate run `29060371292` tested commit
  `7883939cb87ba09fe466415c57ad346b0e105f53`.
- The full macOS job passed, including reproducible builds, minimal-PATH installation, and approved
  upgrade. Windows passed Go and Python parity, then its release builder found that native `Path`
  ordering differed from POSIX archive-name ordering.
- Payload checksum and ZIP entry order now sort explicitly by archive-relative POSIX names. All 10
  release-asset tests, including repeat-build comparison, pass locally. A new pushed Windows run
  remains required.

CI round 4:

- Push-triggered Validate run `29060578007` tested commit
  `e2b57488637acb79c027b73098e3a50d0baf2b1a`.
- macOS passed again. Windows passed Go, parity, and deterministic pack construction; fresh install
  and onboarding also succeeded. The intentional junction-containment rejection then exited 1 as
  designed, but GitHub's PowerShell wrapper promoted that native exit before the assertion ran.
- The junction fixture now temporarily disables native-error promotion, captures exit and JSON,
  requires a nonzero exit plus `unsafe_target`, restores the shell preference, and verifies the
  outside sentinel. PowerShell syntax and native-failure capture probes pass locally. A new pushed
  Windows run remains required.

CI round 5:

- Push-triggered Validate run `29060790532` tested commit
  `b71a40a59403821df06d9a8a8dba9fc771b3d75d`.
- macOS passed all validation. Windows passed Go, parity, deterministic pack construction, fresh
  install, onboarding, check, and the junction fixture's expected `unsafe_target` result plus
  outside-sentinel assertion. The step still ended nonzero because PowerShell retained the
  intentionally failed native repair as `$LASTEXITCODE` after the assertions.
- The workflow now clears only that captured, asserted expected exit before leaving the step. The
  product's fail-closed exit and containment assertions remain unchanged.

CI round 6 - completion evidence:

- Push-triggered Validate run `29060996193` tested commit
  `f0b31bd925ba19798ad28ff54c8ba19d2b77af75` and concluded `success`.
- macOS validation job `86262615364` passed Go, Python parity, installer/release tests, public-core,
  Markdown, schemas, release contracts, reproducible builds, no-Python staged install, and approved
  handoff.
- Windows validation job `86262615369` passed Go, Python parity, deterministic Windows pack build,
  no-Python staged install, checksum rejection, onboarding/check, real junction containment,
  upgrade dry-run, deferred replacement, native failed-reconciliation rollback, and final check.
- Both public-release jobs were skipped. No assets were published.

Negative evidence:

- no release or tag created;
- no pull request created; branch, commit, push, and validation-only workflow actions were limited
  to the user's explicit approval;
- no named consumer synchronized and no local Operating Kit update applied;
- no signing/notarization policy changed;
- no Python implementation removed;
- no new platform added.

## Final Validation

Local validation and the configured macOS/Windows validation-only workflow are complete. Run
`29060996193` passed for commit `f0b31bd925ba19798ad28ff54c8ba19d2b77af75`, including fresh
installation, checksum failure, dry-run hash preservation, approved upgrade, deferred Windows
replacement, native failed-reconciliation rollback, junction containment, and final health checks.
This closes the implementation gate; publication and consumer rollout remain separate.

Blocked audit: the same missing authorization for branch creation, commit, push, and validation-
only workflow execution recurred for three consecutive goal turns. No further local evidence can
substitute for the plan's mandatory real Windows runtime gate.

Resumed: the user explicitly authorized creation of a `codex/` branch, scoped commit, push, and
validation-only workflow execution. PR creation, release, tag, and consumer sync remain excluded.

Final local regression after review fixes: `go test -race ./...` passed; full Python validation
reported `130 passed`; JSON Schema, release-contract, public-core, Markdown, Windows cross-build,
Bash syntax, PowerShell parsing, and diff checks passed. The only diff-check output is Git's
existing LF-to-CRLF working-copy warning for `install.ps1`.
The final catalog/archive fixes were followed by the affected installer, release, packaging, and
repeat-build suite: `27 passed`; fresh isolated macOS install and `0.1.20` -> `0.1.21` upgrade plus
post-check also passed. Windows-target test binaries for reconcile, release, commands, and CLI
compiled successfully, and all Windows workflow PowerShell blocks passed parser validation after
GitHub-expression substitution.
