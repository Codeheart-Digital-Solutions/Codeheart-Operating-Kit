Last updated: 2026-07-04T23:31:11Z (UTC)
Created: 2026-07-04
Status: completed

# Operating Kit Self-Contained Bootstrap Execution Log

## Execution Scope

Plan:
`docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md`

Consumer impact classification:

- `consumer migration required`: existing Python-wheel installs must repair or migrate to the
  self-contained CLI.
- `validator-only change`: release asset and manifest validators change.
- `instruction-only change`: bootstrap docs, release runbook notes, README entries, release notes,
  and plan/register surfaces change.
- `security or safety policy change`: installer trust gates, checksum failure behavior, staged
  unsigned asset boundaries, and signed broad release gates change.

Release-note requirement: required when this change ships in an Operating Kit release.

Migration requirement: required for consumers with existing Python-wheel installs.

Excluded:

- Public release publication, Git tags, GitHub release creation, and named consumer sync.
- Switching live root `manifest.yaml` asset URLs/checksums or published bootstrap asset references
  to unpublished staged assets.
- Removing Python from Foundry AI Execution or module-owned runtimes.
- Linux or Windows ARM support in this milestone.
- Destructive cleanup of legacy Python install libraries.
- Rewriting unrelated dirty repository work.

## Preflight State

Activation time: 2026-07-04T20:27:00Z.

User approval: the active goal explicitly requests activation and implementation of this plan, with
fresh subagents for epic review gates.

Plan status before activation: `draft`.

Plan status after activation: `active`.

Operating Kit installed update state in this checkout:

- `kit_version`: `0.1.17`
- `last_update_check_at`: `2026-07-04T19:24:33Z`
- `next_update_check_due`: `2026-07-11T19:24:33Z`
- `latest_seen_version`: `v0.1.19`
- `update_status`: `update-available`

Source preflight `git status --short` before activation:

```text
 M docs/repo/README.md
 M docs/repo/plans/README.md
 M docs/repo/plans/plan-register.md
?? .codeheart/
?? docs/agent-memory/
?? docs/repo/plans/coordination-sync-pending.md
?? docs/repo/plans/operating-kit-self-contained-bootstrap/
```

Dirty-worktree handling:

- The modified plan indexes/register entries and new
  `docs/repo/plans/operating-kit-self-contained-bootstrap/` files are in scope.
- Existing untracked `.codeheart/`, `docs/agent-memory/`, and
  `docs/repo/plans/coordination-sync-pending.md` surfaces are preserved and must not be swept into
  implementation changes unless a later plan step explicitly owns them.

## Execution Events

- 2026-07-04T20:27:00Z: Activated the implementation plan, created this execution log, recorded
  source preflight state, recorded consumer-impact classes, confirmed release publication is out
  of scope, and confirmed live root release pointers must not switch to staged unpublished assets.
- 2026-07-04T20:30:44Z: Accepted and resolved the EP-00 fresh-review finding by recording the
  review result before treating the epic as complete. Began EP-01 locally by adding the Go module
  skeleton files, then hit a missing Go toolchain blocker during validation.
- 2026-07-04T20:36:28Z: Rechecked EP-01 tooling readiness. `go` and `gofmt` are still missing,
  Homebrew reports `go` is available but not installed, and no local install was performed because
  the tooling-readiness route requires explicit user approval before local toolchain changes.
- 2026-07-04T20:40:06Z: User explicitly approved installing Go for implementation. Installed Go
  with Homebrew, verified `go` and `gofmt`, aligned the skeleton no-argument behavior to the
  Python CLI usage-error surface, formatted Go files, and completed EP-01 local validation.
- 2026-07-04T20:43:03Z: EP-01 fresh review gate passed with no findings. Recorded residual
  version-sync guard risk for later implementation, closed EP-01, and started EP-02.
- 2026-07-04T20:47:29Z: Implemented EP-02 Go resource/state foundation with embedded resource
  access, YAML subset parsing/writing, manifest/component helpers, lock/config helpers, SHA-256
  helpers, platform asset selection, and unit tests. Local Go validation and whitespace checks
  passed; fresh review gate is next.
- 2026-07-04T20:52:55Z: EP-02 fresh review found incomplete representative config
  round-trip coverage and misleading `git diff --check` evidence for untracked Go files. Replaced
  the inline lock/config round-trip fixtures with existing schema-representative YAML fixtures,
  reran validation, and replaced untracked-file whitespace evidence with `gofmt -l` over Go files.
- 2026-07-04T20:55:54Z: EP-02 re-review found SHA-256 parity evidence was still too weak because
  tests used only an inline string and a hardcoded digest. Added Go fixture-file checksum tests
  using expected values produced by the Python `sha256_file` helper, then reran validation.
- 2026-07-04T20:58:59Z: EP-02 second fresh re-review passed with no findings. Closed EP-02 and
  started EP-03.
- 2026-07-04T21:09:04Z: Implemented EP-03 Go root command behavior for `inspect`, `onboard`,
  `init`, `sync`, `check`, and `update-check`; added component writing, drift, capability baseline,
  and command tests; validated Python command oracle tests, Go tests, CLI smokes, and a
  Python-vs-Go parity probe with the approved AD-006 onboarding difference.
- 2026-07-04T21:15:29Z: EP-03 fresh review found the custom parser was not argparse-compatible
  for `--flag=value` and could consume the next flag as a malformed value, risking writes where
  Python would exit `2`. Fixed parser semantics, added command tests and formal Python-vs-Go
  parity tests for equals-style and missing-value flags, and reran validation.
- 2026-07-04T21:19:44Z: EP-03 re-review found subcommand help was not accepted after other
  options like Python argparse accepts. Updated CLI dispatch to detect help anywhere after the
  subcommand and added Go plus Python-vs-Go parity coverage for delayed help.
- 2026-07-04T21:24:29Z: EP-03 second re-review found delayed help was over-accepted when `--help`
  appeared where a value-taking option required a value. Refined command-aware help detection so
  `--help` is callable after complete options but not consumed as a missing option value; added Go
  and Python-vs-Go parity coverage for those value-position help cases.
- 2026-07-04T21:28:24Z: EP-03 third re-review found delayed help was still over-accepted when a
  value-taking option was missing a value and another option appeared before `--help`, such as
  `init --project-name --json --help`. Tightened command-aware help scanning and command parsing
  so any option-looking token in a value slot is treated as missing value; added Go and parity
  coverage for these variants.
- 2026-07-04T21:32:35Z: EP-03 fourth fresh re-review passed with no findings. Closed EP-03 and
  started EP-04 using the formal parity suite created during EP-03.
- 2026-07-04T21:40:34Z: EP-04 fresh review found parity coverage still too shallow for help/error
  surfaces, init/sync/check/update-check checklist cases, and path/platform normalization. Expanded
  `tests/test_go_cli_parity.py` to normalize help text and path separators while checking shared
  help/error tokens, explicitly record the approved AD-006 onboarding help/behavior difference,
  cover init preservation, sync generated-surface merge/release metadata/managed-block refresh,
  stale CLI and missing lock metadata, and update-check agent-notification behavior. Updated Go
  unknown-command output to the Python argparse-compatible invalid-choice shape and reran
  validation.
- 2026-07-04T21:42:59Z: EP-04 fresh re-review passed with no findings. Closed EP-04 and started
  EP-05 installer and legacy migration work.
- 2026-07-04T21:48:21Z: Implemented EP-05 installer migration locally. `install.sh` now installs
  `macos-universal` zip packs without Python or pip, verifies checksums, extracts with `unzip`,
  stages the binary under the install root, smoke-tests `--version`, and replaces the runnable
  command only after validation. `install.ps1` now installs `windows-x64` zip packs without Python
  or pip, verifies checksums, stages and smoke-tests `codeheart-operating-kit.exe`, then writes the
  `.exe` and `.cmd` shim. Both installers preserve local asset/checksum flows, detect legacy
  Python wrapper/lib state without deleting it, and accept the deprecated Python flag only as an
  ignored warning outside normal help. Added installer fixtures and tests for local asset install,
  checksum mismatch, file URL sidecar checksums, malformed archive, missing binary, failed staged
  validation preserving the previous runnable command, and legacy migration preservation.
- 2026-07-04T21:51:51Z: EP-05 fresh review found macOS `file://` URLs with percent-encoded
  spaces were not decoded before `cp`/checksum sidecar reads. Added a shell `file_url_path`
  decoder that handles `file://localhost/...` and percent-encoded paths without invoking Python,
  updated the file-URL installer test to use a space-containing asset path, and reran validation.
- 2026-07-04T21:54:05Z: EP-05 fresh re-review passed with no findings. Closed EP-05 and started
  EP-06 binary release asset builder and manifest work.
- 2026-07-04T21:57:36Z: Completed EP-06 local implementation. The release asset builder now
  builds Go binaries for mandatory `macos-universal` and `windows-x64`, runs `go test ./...`
  before assembly, packages platform zips with binary, bootstrap/install/release-note docs,
  `INSTALL.md`, pack-specific `manifest.json`, and payload `checksums.txt`, writes sidecar
  SHA-256 files, and emits a release-candidate manifest. Release schema, validator fixtures, and
  tests accept `macos-universal` and `windows-x64`; release tests reject Python wheel payloads and
  stale live `manifest.yaml` inside staged packs. Generated ignored staged assets under `dist/`
  and confirmed root `manifest.yaml` and packaged manifest were not changed to staged local
  checksums or unpublished URLs.
- 2026-07-04T22:00:01Z: EP-06 fresh review passed with no findings. Closed EP-06 and started
  EP-07 bootstrap, release, and documentation updates. Updated bootstrap source copy for
  self-contained platform packs and removed the optional native capability setup prompt from base
  onboarding. Updated maintainer README, release runbook, unreleased source release notes, and
  plan indexes for the Go CLI, binary release packs, migration boundary, and staged-release
  guardrails.
- 2026-07-04T22:02:49Z: EP-07 local validation passed for Markdown timestamps, public-core
  hygiene, installer/release doc guard tests, and whitespace. Scans confirmed the edited bootstrap
  no longer contains the optional native capability setup prompt, Python/pip installer guidance, or
  legacy platform asset names.
- 2026-07-04T22:05:44Z: EP-07 fresh review passed with no findings. Closed EP-07 and started
  EP-08 cross-platform validation and staged install proof.
- 2026-07-04T22:10:12Z: Implemented EP-08 workflow and portability updates. Replaced the GitHub
  Actions validation workflow with macOS and Windows validation lanes using Go, Python parity
  tests, binary release packs, staged installer proof, checksum mismatch proof, and base onboarding
  no-native-offer checks. Added `--platform windows-x64` to the release asset builder for Windows
  CI, made Go file-URL handling portable for Windows drive-letter and UNC paths, and made the
  parity harness use the active Python executable instead of assuming `python3`.
- 2026-07-04T22:13:37Z: EP-08 fresh review found two closure blockers: missing actual Windows
  runner evidence and missing fresh low-context bootstrap probe evidence. Completed the
  low-context probe locally using the public first prompt text, current release-candidate
  bootstrap source, and staged macOS release-candidate pack with Python and pip absent from the
  installer PATH. Probe installed `codeheart-operating-kit 0.1.19`, rendered the language prompt
  first, completed non-interactive setup without reporting `native_capabilities`, did not offer
  optional native setup, and `check --json` reported `ok: true`. Windows CI evidence remains
  pending.
- 2026-07-04T22:16:10Z: First pushed Validate run
  `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28721342285`
  failed in `windows-validation` during `go test ./...`. The failure exposed Windows checkout
  line-ending sensitivity in byte-level hash fixtures and CRLF parsing in the mini YAML parser.
- 2026-07-04T22:18:55Z: Added repository line-ending policy through `.gitattributes`, made the Go
  mini YAML parser trim CRLF line endings, added CRLF parser coverage, reran local Go and parity
  validation, rebuilt staged assets, committed `a6196c5`, and pushed the fix branch.
- 2026-07-04T22:20:33Z: Follow-up Validate run
  `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28721439492`
  passed at head `a6196c5202f0890f9ea64680dc237f4c2a1d194b`. `windows-validation` passed
  `go test ./...`, Python-vs-Go parity tests, Windows x64 release pack build, and staged
  `install.ps1` proof without Python on `PATH`. `macos-validation` also passed Go tests, parity,
  installer/release tests, public-core, Markdown, JSON schema, release-manifest validation, full
  release asset build, and staged `install.sh` proof without Python on `PATH`.
- 2026-07-04T22:22:12Z: Reran the low-context bootstrap probe against the rebuilt final staged
  macOS universal pack after the line-ending repair. The probe again installed
  `codeheart-operating-kit 0.1.19`, rendered the language prompt first, completed non-interactive
  setup without reporting `native_capabilities`, did not offer optional native setup, and
  `check --json` reported `ok: true`.
- 2026-07-04T22:26:18Z: EP-08 fresh re-review passed with no blockers. Updated the implementation
  plan lifecycle and checklists, Operating Kit plan register, HQ coordination register, and
  unreleased source release notes for source-complete handoff while preserving the stop before
  public release publication, Git tag creation, GitHub release creation, live manifest pointer
  switch, signing/notarization release decision, and consumer sync.
- 2026-07-04T22:27:47Z: EP-09 local validation passed for Markdown timestamps, public-core
  hygiene, Operating Kit whitespace diff checks, HQ register whitespace diff checks, and no
  remaining unchecked implementation-plan checklist items. Final fresh review gate is next.
- 2026-07-04T22:31:35Z: Final EP-09 fresh review passed with no blockers. Accepted the
  non-blocking caution that the HQ coordination register file is broadly dirty and must not be
  staged wholesale. Marked this execution log complete while keeping public release publication,
  Git tag creation, GitHub release creation, live manifest pointer switch, signing/notarization
  release decision, and consumer sync out of scope.

## Epic Status

| Epic | Status | Evidence |
| --- | --- | --- |
| EP-00 - Plan Activation, Preflight, And Evidence | completed | Plan is active, source preflight is recorded, scope guardrails are recorded, validation passed, and fresh review gate passed after the bookkeeping fix below. |
| EP-01 - Go CLI Skeleton | completed | Go skeleton files created, local validation passed, and fresh review gate passed with no findings. |
| EP-02 - Go Resource And State Foundation | completed | Go foundation files created, local validation passed, review findings fixed, and second fresh re-review passed with no findings. |
| EP-03 - Root Command Behavior Port | completed | Go command ports created, parser/help parity findings fixed, local validation passed, and fourth fresh re-review passed with no findings. |
| EP-04 - Behavior Parity Fixture Suite | completed | Formal Python-vs-Go parity suite expanded after review findings, local validation passed, and fresh re-review passed with no findings. |
| EP-05 - Installers And Legacy Migration | completed | Installer migration implemented locally, review finding fixed, validation passed, and fresh re-review passed with no findings. |
| EP-06 - Binary Release Asset Builder And Manifest | completed | Binary release builder and manifest updates implemented locally, staged ignored dist assets generated, validation passed, and fresh review passed with no findings. |
| EP-07 - Bootstrap, Release, And Documentation Updates | completed | Bootstrap, release notes, README, release runbook, and plan indexes updated locally; validation passed and fresh review passed with no findings. |
| EP-08 - Cross-Platform Validation And Staged Install Proof | completed | Workflow and portability updates implemented; local macOS staged install proof passed; Windows CI runner evidence passed in Validate run `28721439492`; fresh re-review passed with no blockers. |
| EP-09 - Register, Release-Readiness Record, And Handoff | completed | Register, release-readiness, HQ coordination, checklist, and release-stop updates completed locally; validation passed and final fresh review passed with no blockers. |

## Validation Results

EP-00:

- Source status preflight recorded.
- Consumer impact classes recorded.
- Release publication exclusion recorded.
- Live root release pointer exclusion recorded.
- `python3 scripts/validate-markdown-headers.py`: passed.
- `python3 scripts/validate-public-core.py`: passed.
- `git diff --check -- docs/repo/plans/plan-register.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
  passed.

EP-01:

- Added:
  - `go.mod`
  - `cmd/codeheart-operating-kit/main.go`
  - `internal/version/version.go`
  - `internal/cli/cli.go`
  - `internal/cli/cli_test.go`
- `gofmt -w ... && go test ./... && go run ./cmd/codeheart-operating-kit --help && go run ./cmd/codeheart-operating-kit --version`:
  blocked because `gofmt` is missing.
- `go version`: blocked because `go` is missing.
- Recheck:
  - `brew info go`: reports Go stable `1.26.4`, not installed.
  - broad read-only search for existing `go` or `gofmt` under common local paths found no usable
    result before being stopped as no longer useful.
- Read-only checks:
  - `command -v go`: missing.
  - `command -v gofmt`: missing.
  - `command -v brew`: `/opt/homebrew/bin/brew`.
  - platform: macOS arm64.
- Python-side control validation:
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_cli.py -q`:
    passed, 8 tests.
  - `python3 scripts/validate-markdown-headers.py`: passed.
  - `python3 scripts/validate-public-core.py`: passed.
  - `git diff --check -- docs/repo/plans/plan-register.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md go.mod cmd/codeheart-operating-kit/main.go internal/version/version.go internal/cli/cli.go internal/cli/cli_test.go`:
    passed.
- Tooling-readiness route applied:
  `.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`.
- User approval for Go install:
  "Then please resume the goal of implementing everything, and you may install Go. You can do it.
  I give you approval."
- `brew install go`: passed; installed `/opt/homebrew/Cellar/go/1.26.4`.
- `command -v go && command -v gofmt && go version`: passed.
  - `go`: `/opt/homebrew/bin/go`
  - `gofmt`: `/opt/homebrew/bin/gofmt`
  - version: `go version go1.26.4 darwin/arm64`
- Python parity spot checks:
  - `PYTHONPATH=src python3 -m codeheart_operating_kit.cli`: exited `2` with required-command
    usage error.
  - `PYTHONPATH=src python3 -m codeheart_operating_kit.cli --help`: exited `0` and listed root
    commands.
  - `PYTHONPATH=src python3 -m codeheart_operating_kit.cli --version`: exited `0` and printed
    `codeheart-operating-kit 0.1.19`.
- Go skeleton was adjusted so no-argument execution exits `2` with a required-command usage error
  instead of succeeding with root help.
- `gofmt -w cmd/codeheart-operating-kit/main.go internal/version/version.go internal/cli/cli.go internal/cli/cli_test.go`:
  passed.
- `go test ./...`: passed.
- `go run ./cmd/codeheart-operating-kit --help`: passed and listed `onboard`, `inspect`, `init`,
  `sync`, `check`, and `update-check`.
- `go run ./cmd/codeheart-operating-kit --version`: passed and printed
  `codeheart-operating-kit 0.1.19`.
- `go run ./cmd/codeheart-operating-kit onboard --help`: passed and documented explicit
  `--target`/`--project-name` requirements plus the approved setup-only optional capability
  boundary.

EP-02:

- Added:
  - `resources.go`
  - `internal/kitfs/kitfs.go`
  - `internal/kitfs/kitfs_test.go`
  - `internal/hash/hash.go`
  - `internal/hash/hash_test.go`
  - `internal/yamlmini/yamlmini.go`
  - `internal/yamlmini/yamlmini_test.go`
  - `internal/manifest/manifest.go`
  - `internal/manifest/manifest_test.go`
  - `internal/lockfile/lockfile.go`
  - `internal/lockfile/lockfile_test.go`
  - `internal/platforms/platforms.go`
  - `internal/platforms/platforms_test.go`
- `resources.go` exists at the repository root because Go `embed` patterns cannot reference files
  above the package directory. It exposes embedded `components/`, `profiles/`, `templates/`,
  `schemas/`, and `manifest.yaml`; `internal/kitfs` remains the internal access layer.
- YAML compatibility note: the parser intentionally retains the current Python subset behavior,
  including numeric-looking scalar parsing. The lock/config round-trip fixture uses a hex checksum
  with letters so it proves checksum string handling without changing parser semantics.
- `gofmt -w resources.go internal/kitfs/kitfs.go internal/kitfs/kitfs_test.go internal/hash/hash.go internal/hash/hash_test.go internal/yamlmini/yamlmini.go internal/yamlmini/yamlmini_test.go internal/manifest/manifest.go internal/manifest/manifest_test.go internal/lockfile/lockfile.go internal/lockfile/lockfile_test.go internal/platforms/platforms.go internal/platforms/platforms_test.go cmd/codeheart-operating-kit/main.go internal/version/version.go internal/cli/cli.go internal/cli/cli_test.go`:
  passed.
- `go test ./...`: passed.
- `git diff --check -- resources.go internal/kitfs internal/hash internal/yamlmini internal/manifest internal/lockfile internal/platforms docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
  passed, but this evidence was superseded because the EP-02 Go files were still untracked and
  therefore not actually covered by `git diff --check`.
- Review-finding fix:
  - `internal/lockfile/lockfile_test.go` now round-trips existing schema-representative fixtures
    `tests/fixtures/kit-lock.yaml` and `tests/fixtures/kit-config.yaml`.
  - `gofmt -w internal/lockfile/lockfile_test.go`: passed.
  - `go test -count=1 ./...`: passed.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output, covering tracked and untracked Go files.
  - `git diff --check -- docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
    passed for the tracked plan/log edits.
- Second review-finding fix:
  - Python helper command:
    `PYTHONPATH=src python3 - <<'PY' ... from codeheart_operating_kit.manifest import sha256_file ...`
    produced these fixture digests:
    - `tests/fixtures/kit-lock.yaml`:
      `25c20277b11464ec8083bc8459f701f023229a57fc6ed17a77e1f148013ddfee`
    - `tests/fixtures/kit-config.yaml`:
      `454b734b25583549364ceb055fe70785b0116619f9277adee8fc4e19ad8cb00c`
    - `profiles/standard.yaml`:
      `badcf4f34b4cbe7d2eb3f29163b3b3bf65dfa426437793d579ce116e1516f77d`
  - `internal/hash/hash_test.go` now verifies Go file hashing against those Python-helper fixture
    outputs.
  - `gofmt -w internal/hash/hash_test.go`: passed.
  - `go test -count=1 ./...`: passed.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output.
  - `git diff --check -- docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
    passed for the tracked plan/log edits.

## Review Gates

EP-00 review gate:

- Fresh subagent: `019f2ed2-20ad-7833-8e58-1ed4980b7026`.
- Finding: EP-00 was marked completed while the review gate was still recorded as pending.
- Action: recorded the review finding and updated this execution log so EP-00 completion rests on
  validation evidence plus accepted review evidence.
- Verdict after action: accepted; no source behavior issue found in EP-00.

EP-01 review gate:

- Fresh subagent: `019f2edc-81b4-7e93-92e4-77477f3ff871`.
- Findings: none.
- Additional residual risk noted by reviewer: no automated guard yet keeps the Go fallback version
  in `internal/version/version.go` synchronized with release metadata on future version bumps.
- Verdict: pass.

EP-02 review gate:

- Fresh subagent: `019f2ee3-372c-7961-af27-59685fde77fc`.
- Findings:
  - P2: config round-trip fixture was not schema-representative.
  - P3: `git diff --check` evidence did not cover untracked EP-02 Go files.
- Action: fixed both findings as recorded above.
- Verdict after initial review: needs-fix.
- Fresh re-review: pending.
- Fresh re-review subagent: `019f2ee8-2a68-7860-821e-1f59b5e40ab3`.
- Re-review finding:
  - P2: SHA-256 parity evidence did not yet use repository fixture files with Python-helper output
    values.
- Action: fixed as recorded above.
- Verdict after first re-review: needs-fix.
- Second fresh re-review: pending.
- Second fresh re-review subagent: `019f2eea-eaeb-7c80-952d-fe257812e329`.
- Second re-review findings: none.
- Second re-review validation: `go test -count=1 ./...` passed; Python helper digests matched
  recorded values; `gofmt -l` over Go files produced no output; tracked plan/log `git diff --check`
  passed.
- Verdict: pass.

EP-03 review gate: pending; command ports are in progress.

EP-03:

- Added:
  - `internal/capabilities/capabilities.go`
  - `internal/drift/drift.go`
  - `internal/components/components.go`
  - `internal/commands/util.go`
  - `internal/commands/inspect.go`
  - `internal/commands/init.go`
  - `internal/commands/onboard.go`
  - `internal/commands/sync.go`
  - `internal/commands/check.go`
  - `internal/commands/update_check.go`
  - `internal/commands/commands_test.go`
  - `tests/fixtures/parity/folder-states.json`
  - `tests/test_go_cli_parity.py`
- Modified:
  - `internal/cli/cli.go`
  - `internal/cli/cli_test.go`
- Ported command behavior:
  - `inspect`: missing, empty, existing, technical, existing-kit, file-path, and many-entry mode
    logic.
  - `init`: managed content writing, scaffold writing, local user layer, `AGENTS.md` managed
    block rendering, `.gitignore`, lock/config writing, and adoption cleanup report behavior.
  - `onboard`: prompt rendering, required `--target`/`--project-name` gating, write approval,
    ambiguous-folder stop, and base setup write behavior.
  - `sync`: managed refresh, scaffold preservation, `AGENTS.md` managed-block refresh, generated
    surface merge, release metadata refresh, and `.gitignore` repair.
  - `check`: missing CLI, stale CLI, missing routing, missing route targets, missing lock
    metadata, drift, and native capability state reporting.
  - `update-check`: injected latest version, injected metadata URL, current/update-available
    status, failed metadata lookup state, and agent-notification silence behavior.
- Approved intentional difference:
  - Go `onboard --yes` does not offer optional native capability installation/checks, does not
    report `native_capabilities` from onboarding, and leaves lock native capability records at
    `status: unknown` / `command_result_category: not-checked`.
- Python command oracle validation:
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_cli.py tests/test_inspect.py tests/test_onboard.py tests/test_init.py tests/test_sync_check.py tests/test_update_check.py -q`:
    passed, 47 tests.
- Go command validation:
  - `go test -count=1 ./...`: passed.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output.
  - `git diff --check -- docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
    passed for tracked plan/log edits.
- CLI smoke validation:
  - `go run ./cmd/codeheart-operating-kit inspect --json .`: passed.
  - `go run ./cmd/codeheart-operating-kit onboard --target <temp> --project-name Companyname-Automation --yes --json`:
    passed and created `.codeheart/kit.lock.yaml`.
  - `CODEHEART_OPERATING_KIT_CLI=1 go run ./cmd/codeheart-operating-kit check <temp> --json`:
    passed after Go `init` and reported `ok: true`.
  - `go run ./cmd/codeheart-operating-kit update-check <temp> --latest-version 0.2.0 --now 2026-06-13T00:00:00Z --json`:
    passed and reported `update-available`.
- Incremental Python-vs-Go parity probe:
  - Built a temporary Go binary with `go build -o <temp> ./cmd/codeheart-operating-kit`.
  - Compared Python and Go root/subcommand help and version output.
  - Compared Python and Go `inspect --json` modes/reasons for missing, empty, technical,
    existing-kit, and file targets.
  - Compared Python and Go `init` generated file trees, normalized config, selected profile,
    selected components, managed path sets, generated surface sets, and native capability baseline
    state.
  - Compared Python and Go `onboard --yes --json`, verifying the approved difference: Python
    still asks/checks optional native capabilities, Go does not.
  - Verified Go drift/check/sync repair and `update-check` update-available behavior.
  - Result: passed.
- Fresh review gate:
  - Fresh subagent: `019f2ef7-5116-7980-8418-2f6602e83e87`.
  - Findings:
    - High: parser rejected valid `--flag=value` forms and consumed following flags as missing
      values, allowing malformed write commands to proceed.
    - Medium: formal parity files and parser cases were missing from the EP-03 evidence.
  - Verdict after initial review: needs-fix.
- Review-finding fix:
  - `internal/commands/util.go` now supports `--flag=value`, rejects values for boolean flags, and
    treats a following `--flag` as a missing value for value-taking options.
  - `internal/commands/commands_test.go` now covers equals-style flags and missing-value safety
    for `init` and `update-check`.
  - `tests/test_go_cli_parity.py` now covers equals-style flag parity and missing-value parity,
    proving malformed `init` does not create `.codeheart/` and malformed `update-check` does not
    write `latest_seen_version: --json`.
  - `tests/fixtures/parity/folder-states.json` records the inspect folder-state matrix used by
    the formal parity suite.
  - `go test -count=1 ./...`: passed.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_go_cli_parity.py -q`:
    passed, 8 tests.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_cli.py tests/test_inspect.py tests/test_onboard.py tests/test_init.py tests/test_sync_check.py tests/test_update_check.py -q`:
    passed, 47 tests.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output.
  - `git diff --check -- docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md tests/test_go_cli_parity.py tests/fixtures/parity/folder-states.json`:
    passed.
- EP-03 fresh re-review: pending.
- Fresh re-review subagent: `019f2efc-ed23-70d2-a062-ae15da1f06cd`.
- Re-review finding:
  - Medium: subcommand help was only handled as the first subcommand argument, while Python accepts
    `--help` after other options.
- Action:
  - `internal/cli/cli.go` now detects `-h`/`--help` anywhere after the subcommand before command
    parsing.
  - `internal/cli/cli_test.go` covers `init --json --help`.
  - `tests/test_go_cli_parity.py` covers delayed help forms including `init --json --help` and
    `update-check --json --agent-notification --help`.
- Validation after action:
  - `go test -count=1 ./...`: passed.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_go_cli_parity.py -q`:
    passed, 8 tests.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_cli.py tests/test_inspect.py tests/test_onboard.py tests/test_init.py tests/test_sync_check.py tests/test_update_check.py -q`:
    passed, 47 tests.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output.
  - `git diff --check -- docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md tests/test_go_cli_parity.py tests/fixtures/parity/folder-states.json`:
    passed.
- EP-03 second fresh re-review: pending.
- Second fresh re-review subagent: `019f2f00-bbcc-72a0-bcb3-3891bcba91ac`.
- Second re-review finding:
  - Medium: delayed help was over-accepted in value position, e.g. `init --project-name --help`,
    where Python exits `2` for a missing value.
- Action:
  - `internal/cli/cli.go` now derives value-taking options from command help metadata and skips
    only completed value options while scanning for callable help.
  - `internal/cli/cli_test.go` covers `init --project-name --help`,
    `onboard --target --help`, and `update-check --latest-version --help` as missing-value
    errors.
  - `tests/test_go_cli_parity.py` covers the same Python-vs-Go value-position help cases.
- Validation after action:
  - `go test -count=1 ./...`: passed.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_go_cli_parity.py -q`:
    passed, 8 tests.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_cli.py tests/test_inspect.py tests/test_onboard.py tests/test_init.py tests/test_sync_check.py tests/test_update_check.py -q`:
    passed, 47 tests.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output.
  - `git diff --check -- docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md tests/test_go_cli_parity.py tests/fixtures/parity/folder-states.json`:
    passed.
- EP-03 third fresh re-review: pending.
- Third fresh re-review subagent: `019f2f05-227f-70c2-aa8a-c4014b1a5d72`.
- Third re-review finding:
  - High: delayed help was still over-accepted after a missing value option when another option
    appeared before `--help`, e.g. `init --project-name --json --help`.
- Action:
  - `internal/cli/cli.go` now refuses callable-help interception when the token after a
    value-taking option starts with `-`.
  - `internal/commands/util.go` now rejects any option-looking token, including short `-h`, as a
    missing value for value-taking options.
  - `internal/cli/cli_test.go` covers `init --project-name --json --help`,
    `onboard --target --yes --help`, `update-check --latest-version --json --help`, and
    `init --project-name -h` as missing-value errors.
  - `tests/test_go_cli_parity.py` covers the same Python-vs-Go cases.
- Validation after action:
  - `go test -count=1 ./...`: passed.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_go_cli_parity.py -q`:
    passed, 8 tests.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_cli.py tests/test_inspect.py tests/test_onboard.py tests/test_init.py tests/test_sync_check.py tests/test_update_check.py -q`:
    passed, 47 tests.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output.
  - `git diff --check -- docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md tests/test_go_cli_parity.py tests/fixtures/parity/folder-states.json`:
    passed.
- EP-03 fourth fresh re-review: pending.
- Fourth fresh re-review subagent: `019f2f08-b834-7110-a13a-0f9d55ebeb89`.
- Fourth re-review findings: none.
- Fourth re-review validation: `go test -count=1 ./...` passed; `tests/test_go_cli_parity.py`
  passed with 8 tests; Python oracle command suite passed with 47 tests; extra parser checks
  confirmed delayed help and missing-value variants match Python.
- Verdict: pass.

EP-04:

- Added formal parity suite:
  - `tests/fixtures/parity/folder-states.json`
  - `tests/test_go_cli_parity.py`
- Covered parity areas:
  - root help, subcommand help, delayed help, missing-value help, version output, and invalid
    command exit behavior;
  - `inspect` missing, empty, existing, technical, existing-kit, file-path, and many-file folder
    states;
  - `onboard --yes` argument gating and approved AD-006 optional-native-capability difference;
  - `init` generated file tree, normalized config, managed-path set, generated-surface set, and
    normalized native capability baseline;
  - `sync` managed-file repair and scaffold preservation;
  - `check` drift and missing route target reporting;
  - `update-check` current, update-available, metadata URL, and failed metadata lookup behavior;
  - parser parity for equals-style flags and malformed value flags.
- EP-04 validation:
  - `go test -count=1 ./...`: passed.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_go_cli_parity.py -q`:
    passed, 8 tests.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output.
  - `git diff --check -- docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md tests/test_go_cli_parity.py tests/fixtures/parity/folder-states.json`:
    passed.
- EP-04 fresh review:
  - Fresh subagent: `019f2f0d-0cd2-7500-8425-6167072dc016`.
  - Verdict after initial review: needs-fix.
  - Findings:
    - Help/error parity was too shallow because help assertions checked only broad usage tokens and
      invalid-command stderr did not match Python's argparse-compatible surface.
    - Checklist coverage was incomplete for init scaffold preservation, sync generated-surface
      merge/release metadata/managed-block refresh, check stale CLI/missing lock metadata, and
      update-check agent-notification behavior.
    - Path/platform normalization was incomplete because file-tree paths were raw platform strings
      and inspect path text was skipped rather than compared.
- Review-finding fix:
  - `internal/cli/cli.go` now reports unknown commands with the same invalid-choice shape as the
    Python CLI.
  - `internal/cli/cli_test.go` now expects the invalid-choice diagnostic.
  - `tests/test_go_cli_parity.py` now normalizes help text and path separators, compares inspect
    paths, uses POSIX-style relative file-tree paths, checks shared help/error tokens, records the
    approved AD-006 onboarding help and behavior difference, and covers the missing init, sync,
    check, and update-check cases.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_go_cli_parity.py -q`:
    passed, 9 tests.
  - `go test -count=1 ./...`: passed.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_cli.py tests/test_inspect.py tests/test_onboard.py tests/test_init.py tests/test_sync_check.py tests/test_update_check.py -q`:
    passed, 47 tests.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output.
- EP-04 fresh re-review:
  - Fresh subagent: `019f2f14-32d0-7a80-9dd8-c12c816be9ee`.
  - Findings: none.
  - Re-review validation: `go test -count=1 ./...` passed; `tests/test_go_cli_parity.py`
    passed with 9 tests.
  - Residual risk noted by reviewer: help parity is semantic token parity, not byte-for-byte
    output parity. This is accepted by AD-003 and is not blocking.
  - Verdict: pass.

EP-05:

- Started installer and legacy migration work after EP-04 closure.
- Implemented binary installer migration:
  - `install.sh` now targets `codeheart-operating-kit-<version>-macos-universal.zip`, verifies
    SHA-256, extracts with `unzip`, finds `bin/codeheart-operating-kit`, stages under the install
    root, runs staged `--version`, and only then replaces
    `$HOME/.codeheart/operating-kit/bin/codeheart-operating-kit`.
  - `install.ps1` now targets `codeheart-operating-kit-<version>-windows-x64.zip`, verifies
    SHA-256, expands with `Expand-Archive`, finds `bin/codeheart-operating-kit.exe`, stages under
    the install root, runs staged `--version`, and only then replaces the `.exe` and `.cmd` shim
    under `%LOCALAPPDATA%\Codeheart\OperatingKit\bin`.
  - Both installers preserve `--asset-file` / `-AssetFile`, `file://`, `--checksum` /
    `-Checksum`, and `--checksum-file` / `-ChecksumFile` flows.
  - Both installers detect legacy Python wrapper/lib state and preserve legacy files.
  - Deprecated `--python` / `-Python` is no longer shown in normal help and is accepted only as an
    ignored compatibility warning.
  - `internal/commands/check.go` now treats `codeheart-operating-kit.exe` as an installed CLI
    executable name for direct Windows binary execution.
- Installer validation:
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_install_metadata.py tests/test_release_assets.py -q`:
    passed, 20 tests.
  - `go test -count=1 ./...`: passed.
  - `bash -n install.sh`: passed.
  - `pwsh -NoProfile -Command <PowerShell parser check for install.ps1>`: passed.
  - `git diff --check -- install.sh install.ps1 internal/commands/check.go internal/commands/commands_test.go tests/test_install_metadata.py tests/test_release_assets.py tests/fixtures/installer/legacy-python-wrapper.sh tests/fixtures/installer/legacy-python-wrapper.cmd docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
    passed.
- EP-05 fresh review:
  - Fresh subagent: `019f2f1b-5a95-7881-892a-f70d1c4b1a29`.
  - Verdict after initial review: needs-fix.
  - Finding:
    - Medium: macOS `file://` install failed for valid URLs containing encoded spaces because
      `install.sh` stripped `file://` but did not percent-decode the local path before `cp` or
      sidecar checksum reads.
- Review-finding fix:
  - `install.sh` now decodes `file://` paths with `file_url_path`, including `file://localhost/`
    and percent-encoded spaces, without invoking Python.
  - `tests/test_install_metadata.py` now runs the `file://` checksum-sidecar installer path from a
    space-containing directory so `Path.as_uri()` produces `%20`.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_install_metadata.py tests/test_release_assets.py -q`:
    passed, 20 tests.
  - `go test -count=1 ./...`: passed.
  - `bash -n install.sh`: passed.
  - `pwsh -NoProfile -Command <PowerShell parser check for install.ps1>`: passed.
- EP-05 fresh re-review:
  - Fresh subagent: `019f2f1e-3fcc-7053-9615-6fce33d8663b`.
  - Findings: none.
  - Re-review validation: installer metadata and release asset tests passed with 20 tests; Go test
    suite passed; `bash -n install.sh` passed; PowerShell parser check passed; `git diff --check`
    over EP-05 touched files passed.
  - Residual risk noted by reviewer: macOS `file://` decoding covers standard local
    `Path.as_uri()` paths with spaces and `file://localhost`, but not every possible encoded byte
    or non-local file URI host. This is accepted for the EP-05 local installer scope.
  - Verdict: pass.

EP-06:

- Started binary release asset builder and manifest work after EP-05 closure.
- Implemented binary release asset builder and manifest updates:
  - `scripts/build-release-assets.py` runs `go test ./...` before assembly.
  - The builder produces `codeheart-operating-kit-<version>-macos-universal.zip` and
    `codeheart-operating-kit-<version>-windows-x64.zip`.
  - The macOS pack is built from `darwin/arm64` and `darwin/amd64` binaries combined with `lipo`
    and verified for `arm64` and `x86_64`.
  - Platform packs include `bin/`, `bootstrap.md`, `install.sh`, `install.ps1`,
    `release-notes.md`, `INSTALL.md`, pack-specific `manifest.json`, and `checksums.txt`.
  - Platform packs intentionally do not include the live root `manifest.yaml`, because that file
    still points at published assets until the later release run updates live release pointers.
  - Platform pack sidecar `.sha256` files and `release-candidate-manifest-0.1.19.json` were
    generated under ignored `dist/`.
  - `schemas/release-manifest.schema.json`, `scripts/validate-release-manifest.py`, and
    `tests/fixtures/release-manifest.json` accept `macos-universal` and `windows-x64`.
  - `tests/fixtures/release-candidate/release-candidate-manifest.json` records the staged
    release-candidate manifest shape.
- Generated staged assets:
  - `dist/codeheart-operating-kit-0.1.19-macos-universal.zip`
  - `dist/codeheart-operating-kit-0.1.19-macos-universal.zip.sha256`
  - `dist/codeheart-operating-kit-0.1.19-windows-x64.zip`
  - `dist/codeheart-operating-kit-0.1.19-windows-x64.zip.sha256`
  - `dist/release-candidate-manifest-0.1.19.json`
- EP-06 validation:
  - `python3 scripts/build-release-assets.py --output-dir dist`: passed.
  - `go test -count=1 ./...`: passed.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_release_assets.py tests/test_install_metadata.py -q`:
    passed, 21 tests.
  - `python3 scripts/validate-release-manifest.py`: passed.
  - Zip inspection confirmed each staged pack contains the expected binary, pack-specific
    `manifest.json`, and `checksums.txt`; contains no `manifest.yaml`; and contains no `.whl`,
    `.dist-info/`, or `codeheart_operating_kit/` Python payload.
  - `LC_ALL=C LANG=C shasum -a 256 <staged packs>` matched the generated `.sha256` sidecars.
  - `git diff --name-only -- manifest.yaml src/codeheart_operating_kit/resources/manifest.yaml`:
    passed with no output.
  - `git check-ignore -v <staged dist assets>` confirmed the generated staged assets remain
    ignored by the repository `dist/` rule.
  - `git diff --check -- scripts/build-release-assets.py scripts/validate-release-manifest.py schemas/release-manifest.schema.json tests/fixtures/release-manifest.json tests/fixtures/release-candidate/release-candidate-manifest.json tests/test_release_assets.py tests/test_install_metadata.py docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
    passed.
- EP-06 fresh review:
  - Fresh subagent: `019f2f23-f791-7e13-92f6-284730465df1`.
  - Findings: none.
  - Review validation: staged zips contain expected payloads and no Python payloads; sidecar
    checksums match; staged `dist/` assets remain ignored; root `manifest.yaml` and packaged
    manifest have no diff.
  - Residual risk noted by reviewer: the asset builder depends on maintainer-side Go and macOS
    `lipo` availability for the universal macOS pack. This matches the approved maintainer-build
    model and does not affect consumers.
  - Verdict: pass.

EP-07:

- Started bootstrap, release, and documentation updates after EP-06 closure.
- Updated documentation:
  - `bootstrap.md` now names the self-contained `macos-universal` and `windows-x64` platform
    release packs, says base bootstrap installs only the Operating Kit CLI, and removes the
    optional native capability setup prompt from the base onboarding flow.
  - `README.md` now lists the Go CLI entry point, `internal/` implementation, embedded resource
    bundle, binary release-pack builder, and legacy Python behavior oracle.
  - `docs/repo/runbooks/release-operating-kit.md` now includes binary pack validation, checksum
    mismatch validation, staged install proof, live manifest/bootstrap switch timing, staged
    unsigned asset boundary, and signed/notarized broad release gates.
  - `release-notes.md` now has an unreleased source section for the self-contained bootstrap
    migration and explicitly states staged local `dist/` assets are not public release assets.
  - `docs/repo/README.md` and `docs/repo/plans/README.md` now include the implementation plan and
    execution log path for this work.
- EP-07 validation:
  - `python3 scripts/validate-markdown-headers.py`: passed.
  - `python3 scripts/validate-public-core.py`: passed.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_install_metadata.py tests/test_release_assets.py -q`:
    passed, 21 tests.
  - `git diff --check -- bootstrap.md README.md release-notes.md docs/repo/runbooks/release-operating-kit.md docs/repo/README.md docs/repo/plans/README.md tests/test_install_metadata.py docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
    passed.
  - `rg` scans confirmed edited bootstrap surfaces no longer contain the optional native
    capability setup prompt, Python/pip installer guidance, or legacy platform asset names.
- EP-07 fresh review:
  - Fresh subagent: `019f2f28-5dd1-7510-b56b-3f042e74adc2`.
  - Findings: none.
  - Residual risk noted by reviewer: `bootstrap.md` is release-ready source text, not proof that
    public `v0.1.19` assets are already switched. EP-08 and EP-09 still need staged install proof
    and final public manifest/bootstrap switch gating.
  - Verdict: pass.

EP-08:

- Started cross-platform validation and staged install proof after EP-07 closure.
- Implemented validation workflow updates:
  - `.github/workflows/validate.yml` now has a `macos-validation` lane that installs Go and
    Python, runs Go tests, Python-vs-Go parity tests, installer/release asset tests, public-core
    validation, Markdown timestamp validation, JSON schema validation, release-manifest
    validation, binary release asset build, macOS staged install without Python or pip on PATH,
    checksum mismatch proof, and base onboarding no-native-offer proof.
  - `.github/workflows/validate.yml` now has a `windows-validation` lane that installs Go and
    Python, runs Go tests, Python-vs-Go parity tests, builds the Windows x64 release pack, runs
    `install.ps1` without Python or pip on PATH, validates checksum mismatch failure, validates
    staged install, and checks base onboarding does not report or offer native capability setup.
  - Manual public release smoke lanes now use `macos-universal` and `windows-x64` asset names and
    assert base onboarding does not report `native_capabilities`.
- Implemented source portability updates for validation:
  - `scripts/build-release-assets.py` now accepts `--platform all`, `--platform macos-universal`,
    and `--platform windows-x64`; default remains `all`.
  - `internal/commands/util.go` now normalizes `file://` paths for Windows drive-letter and UNC
    file URLs, used by update-check metadata URLs and sync release-manifest file URLs.
  - `tests/test_go_cli_parity.py` now invokes the Python CLI through `sys.executable`, avoiding a
    Windows-only `python3` assumption.
  - `.gitattributes` now forces LF checkout for repository text fixtures and scripts, with
    Windows command and PowerShell scripts checked out as CRLF where appropriate.
  - `internal/yamlmini` now trims CRLF endings while parsing the supported YAML subset.
- EP-08 local validation:
  - `gofmt -w internal/commands/util.go internal/commands/update_check.go internal/commands/sync.go internal/commands/commands_test.go`: passed.
  - `go test -count=1 ./...`: passed.
  - `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_go_cli_parity.py tests/test_release_assets.py tests/test_install_metadata.py -q`:
    passed, 31 tests.
  - `python3 scripts/build-release-assets.py --platform windows-x64 --output-dir /tmp/codeheart-ok-windows-pack-test`:
    passed and produced only the Windows x64 pack, checksum, and release-candidate manifest.
  - `ruby -e 'require "yaml"; YAML.load_file(".github/workflows/validate.yml")'`: passed.
  - `python3 scripts/build-release-assets.py --output-dir dist`: passed and regenerated staged
    ignored macOS universal and Windows x64 assets.
  - Local macOS staged install proof ran `install.sh` with `python`, `python3`, `pip`, and `pip3`
    absent from installer `PATH`; checksum mismatch failed closed; valid staged install succeeded;
    installed binary reported `codeheart-operating-kit 0.1.19`; `onboard --yes --json` did not
    report `native_capabilities` and did not offer optional native setup; `check --json` reported
    `ok: true`.
  - `python3 scripts/validate-json-schemas.py`: passed.
  - `python3 scripts/validate-markdown-headers.py`: passed.
  - `python3 scripts/validate-public-core.py`: passed.
  - `python3 scripts/validate-release-manifest.py`: passed.
  - `find . -path './.git' -prune -o -name '*.go' -print0 | xargs -0 gofmt -l`: passed with no
    output.
  - `git diff --check -- .github/workflows/validate.yml scripts/build-release-assets.py internal/commands/util.go internal/commands/update_check.go internal/commands/sync.go internal/commands/commands_test.go tests/test_go_cli_parity.py tests/test_release_assets.py docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`:
    passed.
  - After Windows CI exposed CRLF sensitivity, reran `gofmt -w
    internal/yamlmini/yamlmini.go internal/yamlmini/yamlmini_test.go`: passed.
  - After Windows CI exposed CRLF sensitivity, reran `go test -count=1 ./...`: passed.
  - After Windows CI exposed CRLF sensitivity, reran `uv run --no-project --with pytest --with pip
    --with setuptools --with wheel python -m pytest tests/test_go_cli_parity.py
    tests/test_release_assets.py tests/test_install_metadata.py -q`: passed, 31 tests.
  - After Windows CI exposed CRLF sensitivity, reran
    `python3 scripts/build-release-assets.py --platform windows-x64 --output-dir /tmp/codeheart-ok-windows-pack-test`:
    passed.
  - After Windows CI exposed CRLF sensitivity, reran
    `python3 scripts/build-release-assets.py --output-dir dist`: passed and regenerated staged
    ignored macOS universal and Windows x64 assets.
  - `env LC_ALL=C LANG=C shasum -a 256 <staged packs>` matched the regenerated sidecar checksums.
- EP-08 external validation:
  - Initial pushed Validate run
    `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28721342285`
    failed in `windows-validation` during `go test ./...`, proving the earlier local-only Windows
    build simulation was insufficient as closure evidence.
  - Follow-up pushed Validate run
    `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28721439492`
    passed at head `a6196c5202f0890f9ea64680dc237f4c2a1d194b`.
  - `windows-validation` job
    `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28721439492/job/85171600104`
    passed, including `go test ./...`, Python-vs-Go parity tests,
    `python scripts/build-release-assets.py --platform windows-x64 --output-dir dist`, and
    `Validate Windows staged install without Python on PATH`.
  - `macos-validation` job
    `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28721439492/job/85171600115`
    passed, including Go tests, parity tests, installer/release tests, public-core validation,
    Markdown validation, JSON schema validation, release-manifest validation, full release asset
    build, and staged `install.sh` proof without Python on `PATH`.
  - `windows-public-release` and `macos-public-release` jobs were skipped as expected because this
    implementation run does not publish public release assets.
- EP-08 fresh review:
  - Fresh subagent: `019f2f2f-7995-7531-9c62-d17c063973bc`.
  - Verdict after initial review: needs-fix.
  - Findings:
    - High: EP-08 cannot close without actual Windows runner evidence for the x64 build, Go tests,
      parity tests, Windows pack build, and staged `install.ps1` proof.
    - High: fresh low-context bootstrap probe evidence was not recorded.
- Review-finding fix completed:
  - Fresh low-context bootstrap probe was rerun after the Windows line-ending repair and passed
    using:
    - public first prompt text:
      `Set up Codeheart Operating Kit: https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/latest/download/bootstrap.md`
    - current release-candidate `bootstrap.md` source;
    - staged release-candidate pack
      `dist/codeheart-operating-kit-0.1.19-macos-universal.zip`;
    - installer `PATH` with `python`, `python3`, `pip`, and `pip3` absent.
  - Probe evidence:
    - first prompt SHA-256 for the exact prompt string without a trailing newline:
      `651026e955b450f13564dda427b4cd1022a9a96a4dc04600dc70cda35e6ca33c`
    - bootstrap source SHA-256:
      `8efd5891b87417481fc4fb1b001e73a95926c738bf9e829b247cc5553309f014`
    - staged macOS pack SHA-256:
      `845a783bea2dfd57cf5a4072e4dfbab30cd319053faf33a40bbc8104fc128e0a`
    - installed binary version: `codeheart-operating-kit 0.1.19`
    - first onboarding line:
      `Choose setup language / Sprache waehlen / 选择设置语言:`
    - non-interactive onboarding did not report `native_capabilities`;
    - non-interactive onboarding did not offer `Should I check and set up these tools now?`;
    - `check --json` reported `ok: true`.
  - Updated staged pack checksums:
    - `dist/codeheart-operating-kit-0.1.19-macos-universal.zip`:
      `845a783bea2dfd57cf5a4072e4dfbab30cd319053faf33a40bbc8104fc128e0a`
    - `dist/codeheart-operating-kit-0.1.19-windows-x64.zip`:
      `247ec2d3815ba326a875a7d256fff3c086554c09d9f8e5c69d2690205f9e9a20`
- EP-08 fresh re-review:
  - Fresh subagent: `019f2f3a-be6c-7ee0-bcf0-efdabfd213ca`.
  - Findings: none.
  - Review validation: verified EP-08 acceptance coverage, GitHub Actions Validate run
    `28721439492`, Windows job `85171600104`, macOS job `85171600115`, line-ending policy, CRLF
    YAML parser coverage, and release-pack checksum behavior.
  - Residual risks noted by reviewer: Windows proof is CI-backed rather than a physical fresh
    Windows machine; the Windows no-Python proof removes `python`, `python3`, `pip`, and `pip3`
    from `PATH` but does not enumerate every possible Python launcher; public release smoke jobs,
    tag/release publication, signing, notarization, and public asset switching remain outside
    EP-08.
  - Verdict: pass.

EP-09:

- Started register, release-readiness, and handoff updates after EP-08 closure.
- Updated lifecycle surfaces:
  - `docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md`
    now has status `completed` and all implementation checklist items checked.
  - `docs/repo/plans/plan-register.md` now marks `OK-PR-024` completed, adds the execution log as
    canonical evidence, records the validation summary, records staged-source asset boundaries,
    and preserves the stop before public release publication.
  - `Codeheart-HQ/docs/repo/plans/plan-register.md` now records that the HQ discovery's canonical
    Operating Kit child implementation `OK-PR-024` completed source readiness while public release
    publication remains separate approval-gated release-run work.
  - `release-notes.md` now records the passing GitHub Actions macOS and Windows validation as
    unreleased source-readiness evidence.
- EP-09 validation:
  - `python3 scripts/validate-markdown-headers.py`: passed.
  - `python3 scripts/validate-public-core.py`: passed.
  - `git diff --check -- docs/repo/plans/plan-register.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md release-notes.md .gitattributes internal/yamlmini/yamlmini.go internal/yamlmini/yamlmini_test.go`:
    passed.
  - In Codeheart-HQ, `git diff --check -- docs/repo/plans/plan-register.md`: passed.
  - `rg -n "^- \\[ \\]" docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_implementation_doc.md`:
    passed with no unchecked implementation checklist items.
- EP-09 fresh review:
  - Fresh subagent: `019f2f3f-6fef-7fc1-a849-44ecedf6e3e4`.
  - Findings: no blockers.
  - Non-blocking caution: the HQ coordination register contains the intended `CODEHEART-HQ-PR-009`
    update, but the file also contains unrelated dirty register changes and must not be staged
    wholesale.
  - Closeout note: execution log completion status was expected final bookkeeping after the review
    returned.
  - Review validation: verified GitHub Actions Validate run `28721439492`, macOS and Windows
    validation, skipped public-release jobs, release-stop language, Markdown timestamp validation,
    public-core validation, diff whitespace checks, no unchecked implementation checklist items,
    and staged pack hashes.
  - Verdict: pass.

## Release Evidence

Release run:

- User approval: explicit chat request to make the new release.
- Release version: `v0.1.20`.
- Release source commit: `e3acfb2717d707ff9a5f523db9440e9da8d34834`.
- Release source commit message: `Release Operating Kit 0.1.20 self-contained bootstrap`.
- Release tag: `v0.1.20`, pushed to GitHub.
- Release URL:
  `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v0.1.20`.
- Published at: `2026-07-04T23:28:52Z`.
- `main` was fast-forwarded to the validated release source commit before tagging and publishing.

Published assets:

- `manifest.yaml`
- `bootstrap.md`
- `install.sh`
- `install.ps1`
- `release-notes.md`
- `codeheart-operating-kit-0.1.20-macos-universal.zip`
- `codeheart-operating-kit-0.1.20-macos-universal.zip.sha256`
- `codeheart-operating-kit-0.1.20-windows-x64.zip`
- `codeheart-operating-kit-0.1.20-windows-x64.zip.sha256`

Live manifest checksums:

- `bootstrap.md`:
  `33337c85754ddf8e838af8d5cf4fced19ed7e3c9dcb68d49e75398d18e750ad8`
- `install.sh`:
  `ee8a975f81454ee8c18cc4bcea21c984e85029696bc862b7ce8249b32e56104c`
- `install.ps1`:
  `1db4ba0a2c7c3a1073672a4b6ff634e682e8795336009e7b94c6b2f8ffbb4758`
- `release-notes.md`:
  `fc0930d36058e0e560b6d99bea71d56df58370a1b44facd1345f45b0fb6a9e2a`
- `codeheart-operating-kit-0.1.20-macos-universal.zip`:
  `19eafcfbcfffe7daa1baf256d6b05aee356417d1939c8f302e7015105061036e`
- `codeheart-operating-kit-0.1.20-macos-universal.zip.sha256`:
  `97eaaf2b13c999852509345114a336c36018b05502693a4c2cb60fe76d6f3cec`
- `codeheart-operating-kit-0.1.20-windows-x64.zip`:
  `b3d9c32d1b06fc60ced465eefc07d27078b69f09dc3902260defccd7a928c799`
- `codeheart-operating-kit-0.1.20-windows-x64.zip.sha256`:
  `4188a17bc28e200b81166a121950156dd72401f6d62ca2215339dab08204c382`

Release-source validation:

- Local `python3 scripts/build-release-assets.py --output-dir dist`: passed and produced the
  macOS universal and Windows x64 packs.
- Local `go test -count=1 ./...`: passed.
- Local `uv run --no-project --with pytest --with pip --with setuptools --with wheel python -m
  pytest tests/test_go_cli_parity.py tests/test_release_assets.py tests/test_install_metadata.py
  -q`: passed, 31 tests.
- Local `python3 scripts/validate-public-core.py`: passed.
- Local `python3 scripts/validate-markdown-headers.py`: passed.
- Local `python3 scripts/validate-json-schemas.py`: passed.
- Local `python3 scripts/validate-release-manifest.py`: passed.
- Local `git diff --check`: passed.
- Local macOS restricted-path staged install: passed with `python`, `python3`, `pip`, and `pip3`
  absent from installer `PATH`; bad checksum failed closed; valid install succeeded; installed
  binary reported `codeheart-operating-kit 0.1.20`.
- Branch Validate run before `main` fast-forward:
  `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28723014116`
  passed for commit `e3acfb2717d707ff9a5f523db9440e9da8d34834`.
- `main` push Validate run:
  `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28723062205`
  passed for commit `e3acfb2717d707ff9a5f523db9440e9da8d34834`.
- Public release smoke workflow-dispatch run:
  `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/28723075077`
  passed with `release_version=v0.1.20`, including macOS and Windows public release download
  checks plus normal macOS and Windows validation lanes.

Signing and distribution boundary:

- `v0.1.20` is published as an unsigned internal/prototype Operating Kit release.
- The release is not Apple-notarized and does not include Windows Authenticode signing.
- The installer trust gate for this release is the published SHA-256 manifest, checksum sidecars,
  and staged binary smoke tests.
- Broad external distribution remains gated on a later signing/notarization decision.

Consumer sync:

- Named consumer repository sync was not performed in this release run.
- Consumers can install or repair through the published `v0.1.20` bootstrap and platform packs.

## Residual Risk

- Go was installed locally with explicit user approval as a maintainer source-build prerequisite.
  This remains a maintainer dependency only, not a consumer install prerequisite.
- Go fallback version synchronization is not yet automatically guarded against future release
  metadata drift.
- `v0.1.20` is unsigned and not notarized; broad external distribution should wait for a signing
  and notarization decision.
- Consumer repository sync remains separate work and was not performed during this release run.
