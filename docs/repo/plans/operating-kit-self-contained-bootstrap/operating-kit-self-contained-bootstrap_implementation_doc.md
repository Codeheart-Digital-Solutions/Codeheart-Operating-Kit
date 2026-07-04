Last updated: 2026-07-04T22:26:18Z (UTC)
Created: 2026-07-04
Status: completed
Execution log: docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md

# Document Header

## Operating Kit Self-Contained Bootstrap Implementation Plan

Overview: Replace the Python-wheel-first Operating Kit CLI distribution with a self-contained Go
CLI and platform release packs so fresh macOS and Windows machines can install, repair, and run
base Operating Kit onboarding without Python, pip, Homebrew, Git, GitHub CLI, Node.js, or a
package manager.

This plan targets the Codeheart Operating Kit source repository. The discovery and accepted
capability scope were drafted in Codeheart-HQ as the coordination-home artifact:
`Codeheart-HQ:docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_discovery_doc.md`.

This plan does not publish a public release, push a tag, or sync named consumers. It prepares and
validates the source change, local staged release assets, migration behavior, and release gates so
a later release run can publish with explicit release approval. It must not switch live public
release pointers, root release-manifest checksums, or published bootstrap asset references to
unpublished staged assets.

Essential context:

| Source | Why it matters |
| --- | --- |
| `AGENTS.md` | Maintainer bootstrap, public-core safety, source-of-truth order, and required change routes. |
| `README.md` | Public setup entry point and maintainer entry point inventory that will change when the CLI distribution changes. |
| `bootstrap.md` | Public first-run bootstrap document that must remain the user-facing start point. |
| `install.sh` | macOS user-level installer and repair script that currently depends on a Python wheel. |
| `install.ps1` | Windows user-level installer and repair script that currently depends on a Python wheel. |
| `scripts/build-release-assets.py` | Current release asset builder that produces wheel-bearing macOS and Windows assets. |
| `manifest.yaml` | Root release manifest that advertises public assets, checksums, platforms, component impact, and generated surfaces. |
| `src/codeheart_operating_kit/` | Current Python CLI behavior and resource model, retained as the parity oracle during the port. |
| `tests/test_cli.py` | Current command-surface coverage for the Python CLI. |
| `tests/test_init.py` | Current installed file-tree, lock/config, `.gitignore`, and scaffold behavior coverage. |
| `tests/test_sync_check.py` | Current sync, check, drift, managed block, and generated-surface behavior coverage. |
| `tests/test_update_check.py` | Current update-check state and metadata lookup coverage. |
| `tests/test_release_assets.py` | Current release asset builder and manifest validation coverage. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure before changing source, docs, schemas, templates, validators, installers, release assets, or CLI behavior. |
| `docs/repo/runbooks/release-operating-kit.md` | Release procedure and stop conditions for later publication work. |
| `docs/repo/reference/placement-contract.md` | Source and installed placement contract for managed content, generated surfaces, local state, and CLI source. |
| `docs/repo/reference/consumer-impact-classification.md` | Required classification for CLI behavior, generated behavior, release assets, migration, and consumer action. |
| `Codeheart-HQ:docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_discovery_doc.md` | Accepted decisions, scope boundaries, and implementation capability scope. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Implement a self-contained Operating Kit bootstrap capability for macOS and Windows while
preserving the current root command behavior.

Completion is proven when:

- a Go CLI implements `onboard`, `inspect`, `init`, `sync`, `check`, and `update-check`;
- the Go CLI embeds or carries the managed content needed for init, sync, check, and repair
  without a Python runtime;
- macOS and Windows installers download, verify, install, and repair platform release packs
  without using Python or pip;
- staged macOS and Windows release packs contain the Go CLI, bootstrap/install docs, manifest
  data, checksums, and no Python wheel payload;
- legacy Python-wheel installs are detected and repaired without manual deletion as the normal
  path;
- base onboarding does not install, offer to install, or implicitly check optional native
  capabilities during the first-run setup flow;
- behavior parity fixtures compare the current Python CLI and the new Go CLI across command
  surface, exit codes, stable output, generated file trees, normalized lock/config YAML, drift
  reports, update-check state, and explicitly approved intentional differences;
- fresh macOS and Windows install smoke tests prove setup works with Python and pip absent from
  `PATH`;
- checksum mismatch, metadata lookup failure, stale CLI, missing route targets, and drift cases
  fail closed with clear output;
- `bootstrap.md`, `README.md`, install scripts, release-manifest schema or release-candidate
  fixtures, release runbook notes, tests, and plan-register entries are updated without publishing
  or advertising unpublished live assets;
- public-core, Markdown, schema, release-manifest, Python parity, Go unit, release asset, installer
  mismatch, and local staged install validations pass.

## 1.2 Project And Problem Context

The current Operating Kit public bootstrap is intentionally simple for users: give Codex a
`bootstrap.md` URL, let the agent install or repair the CLI, then run agent-guided onboarding.
That model is sound.

The implementation problem is that the current release asset contains a Python wheel, and both
platform installers install that wheel with `python -m pip`. A fresh non-developer machine may not
have Python and pip ready. That makes Python a prerequisite for the very layer that should teach
agents how to handle missing local tooling.

The accepted direction is to keep Python available later for Foundry modules and repo-local
runtimes, but remove Python and pip from the base Operating Kit install path.

## 1.3 Current State Analysis

Current source state:

- `src/codeheart_operating_kit/` is the Python CLI implementation and resource reader.
- `pyproject.toml` declares the Python package and console entry point.
- `scripts/build-release-assets.py` builds a wheel and places that wheel in macOS and Windows
  release assets.
- `install.sh` downloads `codeheart-operating-kit-<version>-macos.tar.gz`, verifies SHA-256,
  extracts a wheel, runs `python3 -m pip install --target`, and writes a shell wrapper.
- `install.ps1` downloads `codeheart-operating-kit-<version>-windows.zip`, verifies SHA-256,
  extracts a wheel, runs `python -m pip install --target`, and writes a `.cmd` wrapper.
- `bootstrap.md` documents the current installer commands and says both installers verify the
  pinned release asset checksum.
- tests cover the Python CLI commands and current wheel-bearing release assets.

Target source state:

- a Go CLI under `cmd/` and `internal/` owns the root command behavior;
- Python CLI source remains temporarily as a parity oracle and compatibility reference, but the
  staged consumer release assets no longer ship a wheel;
- installers install a platform binary release pack and do not invoke Python;
- base `onboard` setup is setup-only for optional tooling: it does not offer native capability
  installation and does not run implicit optional capability checks during first-run onboarding;
- tests prove semantic parity before the Go CLI replaces Python in release assets;
- release-candidate manifests, fixtures, and installer tests understand platform binary assets,
  while live public manifest updates remain release-run work;
- signing and notarization are planned as release gates, with unsigned staged assets allowed only
  for internal/prototype validation.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Port the CLI behavior behind a new Go implementation while keeping the Python CLI available in the
source tree as the parity oracle. Capture Python parity baselines before accepting each command
port. Replace consumer release asset payloads only after parity tests exist. Keep generated
consumer surfaces and managed content behavior unchanged unless a parity test forces an explicitly
approved semantic change. Treat base onboarding optional-tool behavior as one approved intentional
difference: first-run setup must not offer or install optional native capabilities.

Public release publication remains outside this plan. Staged asset validation uses generated
release-candidate manifests under `dist/` and test fixtures. The live root `manifest.yaml` must not
be changed to advertise unpublished asset URLs or staged local checksums during this implementation
plan; switching live public release pointers belongs to the later release run.

Expected file tree:

```text
Codeheart-Operating-Kit/
  README.md                                                     # modify entry point inventory
  bootstrap.md                                                  # modify release-ready bootstrap text without live unpublished asset switch
  install.sh                                                    # modify macOS binary install
  install.ps1                                                   # modify Windows binary install
  manifest.yaml                                                 # leave live release checksums untouched until release run
  release-notes.md                                              # modify release note draft
  go.mod                                                        # create Go module
  go.sum                                                        # create when dependencies require it
  cmd/
    codeheart-operating-kit/
      main.go                                                   # create CLI entry point
  internal/
    cli/                                                        # create command dispatch
    commands/                                                   # create command implementations
    kitfs/                                                      # create embedded managed resource access
    lockfile/                                                   # create lock/config helpers
    manifest/                                                   # create manifest and checksum helpers
    platforms/                                                  # create platform asset helpers
  scripts/
    build-release-assets.py                                     # modify staged binary release packs
    validate-release-manifest.py                                # modify asset schema expectations
  schemas/
    release-manifest.schema.json                                # modify asset platform/name pattern
  src/codeheart_operating_kit/                                  # keep as parity oracle
  tests/
      fixtures/
        parity/                                                   # create normalized expected states
        installer/                                                # create legacy install fixtures
        release-candidate/                                        # create staged manifest fixtures
    test_go_cli_parity.py                                       # create Python-vs-Go parity tests
    test_release_assets.py                                      # modify binary release asset tests
    test_install_metadata.py                                    # modify legacy install metadata tests
  docs/
    repo/
      README.md                                                 # modify plan pointer
      plans/
        README.md                                               # modify plan pointer
        plan-register.md                                        # modify OK-PR-024
        operating-kit-self-contained-bootstrap/
          operating-kit-self-contained-bootstrap_implementation_doc.md # create
          operating-kit-self-contained-bootstrap_execution_log.md       # create during execution
```

Consumer impact classification:

- `consumer migration required`: existing Python-wheel installs must repair or migrate to the
  self-contained CLI.
- `validator-only change`: release asset and manifest validators change.
- `instruction-only change`: bootstrap docs, release runbook notes, README entries, and release
  notes change.
- `security or safety policy change`: installer trust gates, checksum failure behavior, and
  unsigned-versus-signed release staging guidance change.

Runbook change coverage:

- Durable runbooks materially changed by this plan: `docs/repo/runbooks/release-operating-kit.md`
  and source `bootstrap.md`.
- Audience classes: `release-operating-kit.md` is maintainer-facing; `bootstrap.md` is public
  bootstrap/agent contract text, but this plan prepares release-ready source text only and does not
  publish the versioned bootstrap asset.
- Required coverage: binary asset build/verification, legacy migration evidence, checksum
  mismatch validation, macOS staged install, Windows staged install, and signing/notarization
  release gate notes.
- Tooling blockers: Go toolchain is a maintainer source-build prerequisite; it is not a consumer
  install prerequisite. Consumer local tooling readiness remains governed by installed Operating
  Kit routes after base setup.

Routing-standard coverage:

- This plan changes bootstrap and CLI repair behavior, not normal route selection doctrine.
- It does not change capability advertisements, route cards, or owner selection.
- Fresh low-context routing probe is not required for route selection.
- Fresh low-context bootstrap probe is required: a fresh agent receives the public first prompt,
  follows the release-candidate `bootstrap.md` text against staged assets, selects the platform
  installer, verifies asset checksums, and reaches `codeheart-operating-kit onboard` without
  choosing Python, Homebrew, Git, GitHub CLI, Node.js, or a package manager. The probe must not
  use unpublished staged assets as if they were a live public release.

## 2.2 Open Questions And Assumptions Requiring Clarification

OQ-001 - Should the first implementation include public release publication?

- `BLOCKER: no`
- `Affects: EP-09`
- Unlocks GitHub release publication, tag creation, and consumer sync proof.
- Recommended default: keep publication out of this plan. This plan prepares and validates source
  and staged local release assets. Public release requires a separate approval-gated release run.
  Do not switch root `manifest.yaml`, published bootstrap asset references, or live checksums to
  staged local asset values during this implementation plan.

OQ-002 - Should the first release remove the Python package from source?

- `BLOCKER: no`
- `Affects: EP-01, EP-03, EP-04, EP-09`
- Unlocks source cleanup.
- Recommended default: keep the Python package in source as a parity oracle through the first Go
  release. Remove it only after at least one validated Go release and migration proof.

OQ-003 - Should macOS signing/notarization and Windows Authenticode signing be fully automated in
this implementation?

- `BLOCKER: no`
- `Affects: EP-08, EP-09`
- Unlocks broad external release readiness.
- Recommended default: document and gate signing/notarization, but allow unsigned staged internal
  assets. Broad public release stops until signing/notarization inputs exist or a maintainer
  explicitly approves an unsigned public release.

OQ-004 - Should Windows ARM and Linux be included in the first platform matrix?

- `BLOCKER: no`
- `Affects: EP-04, EP-06, EP-09`
- Unlocks additional asset builds and CI jobs.
- Recommended default: first matrix is macOS universal and Windows x64. `macos-universal` is a
  hard target, not a best-effort label. If universal output cannot be produced, stop and revise the
  plan to separate `macos-arm64` and `macos-amd64` assets before changing installer docs. Defer
  Linux and Windows ARM until a real consumer need appears.

OQ-005 - Should offline/local asset install be mandatory in the first implementation?

- `BLOCKER: no`
- `Affects: EP-05, EP-06`
- Unlocks disconnected install proof.
- Recommended default: preserve `--asset-file` and `file://` install support because current
  installers already have that shape. Do not broaden offline behavior beyond asset-file install
  and checksum-file validation.

## 2.3 Architectural Decisions With Reasoning

AD-001 - Use Go for the root CLI

1. Problem being solved: The CLI must run on fresh macOS and Windows machines without Python or
   pip.
2. Simplest working solution: Implement the current root command surface in Go.
3. What may change in 6-12 months: The Python CLI can be removed after Go parity and migration
   prove stable.
4. Rationale for the chosen approach: Go produces self-contained binaries, has a strong standard
   library for filesystem, HTTP, hashing, archive, and JSON behavior, and fits the current CLI's
   operational complexity.
5. Alternatives considered and why not chosen: Rust is viable but heavier for this mostly
   file-orchestration CLI. Bundled Python keeps the dependency hidden inside a larger artifact.
   Shell and PowerShell alone are too weak for robust cross-platform metadata behavior.

AD-002 - Embed managed resources into the Go CLI

1. Problem being solved: `init` and `sync` need managed component files without a Python package
   resource loader.
2. Simplest working solution: Use Go `embed` for `components/`, `profiles/`, `templates/`, and
   required schema/manifest resources.
3. What may change in 6-12 months: Very large resources may justify a separate payload directory.
4. Rationale for the chosen approach: Embedded resources make the installed CLI self-contained and
   reduce installer complexity.
5. Alternatives considered and why not chosen: A separate managed-content payload is workable but
   introduces more path, manifest, and repair edge cases during first bootstrap.

AD-003 - Preserve semantic parity, not byte-for-byte parity

1. Problem being solved: A rewrite can accidentally change user-visible behavior or generated
   state.
2. Simplest working solution: Compare normalized outputs, generated file trees, lock/config YAML,
   checksums, exit codes, and error classes, with every intentional difference recorded as an
   explicit fixture expectation.
3. What may change in 6-12 months: Intentional Go-only behavior can become the new baseline after
   an approved release.
4. Rationale for the chosen approach: Byte-for-byte comparisons are brittle for timestamps, path
   separators, map ordering, and archive metadata.
5. Alternatives considered and why not chosen: Manual smoke testing alone is too weak for a
   bootstrap rewrite.

AD-004 - Keep the legacy Python install path as a migration input

1. Problem being solved: Existing consumers may already have Python-wheel installs under user
   install roots.
2. Simplest working solution: Detect legacy wrapper/lib layouts, install the Go binary into the
   same user-level bin path, preserve old lib content, and report migration state.
3. What may change in 6-12 months: A cleanup command may remove old Python libraries after enough
   releases.
4. Rationale for the chosen approach: Replacing the runnable wrapper while preserving old payloads
   avoids destructive cleanup and gives users a safe rollback trail.
5. Alternatives considered and why not chosen: Deleting old Python libraries during install is too
   risky for a repair path.

AD-005 - Separate staged asset validation from public release publication

1. Problem being solved: The implementation needs release asset validation without silently
   publishing a public release.
2. Simplest working solution: Build and validate local staged assets in `dist/`, generate
   release-candidate manifests and release notes, and stop before tag, live root-manifest switch,
   or GitHub release publication.
3. What may change in 6-12 months: A later release plan can combine implementation and publication
   after signing automation exists.
4. Rationale for the chosen approach: Publication changes public state and needs explicit release
   approval.
5. Alternatives considered and why not chosen: Publishing inside this implementation plan would
   mix source migration risk with public release authority.

AD-006 - Make base onboarding setup-only for optional native capabilities

1. Problem being solved: First-run onboarding must stay simple and must not install further tools
   while proving the base Operating Kit can bootstrap on a fresh machine.
2. Simplest working solution: `onboard --yes` writes base Operating Kit setup only. It must not
   offer optional native capability installation, run plugin installation, or implicitly check
   optional capability state. Optional tooling remains governed by later tooling-readiness routes
   when a real task needs it.
3. What may change in 6-12 months: A separate explicit command or route may check optional
   capabilities when the user asks for that capability, but it remains outside base onboarding.
4. Rationale for the chosen approach: The Operating Kit base layer should establish safe routing
   and repair before any optional local tooling or connector work begins.
5. Alternatives considered and why not chosen: Preserving the current Python `onboard --yes`
   native capability installation behavior would keep parity but violate the accepted bootstrap
   boundary and introduce avoidable external-state changes.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-00 | Plan activation, source preflight, and execution evidence are ready. | S | none |
| EP-01 | Go CLI project skeleton builds and exposes the root command surface. | M | EP-00 |
| EP-02 | Go CLI implements managed resources, YAML, checksums, lock/config, and platform helpers. | L | EP-01 |
| EP-03 | Go CLI implements `inspect`, `onboard`, `init`, `sync`, `check`, and `update-check` with incremental parity baselines. | XL | EP-02 |
| EP-04 | Behavior parity fixtures finalize Go CLI equivalence against the Python CLI. | L | EP-03 |
| EP-05 | macOS and Windows installers install binary release packs and migrate legacy installs. | L | EP-04 |
| EP-06 | Release asset builder, manifest schema, and release asset tests produce binary packs. | L | EP-05 |
| EP-07 | Documentation and managed release guidance describe the self-contained bootstrap. | M | EP-06 |
| EP-08 | Cross-platform validation and staged install proof complete. | L | EP-07 |
| EP-09 | Register, release-readiness record, and handoff are finalized without public publication. | S | EP-08 |

## EP-00 - Plan Activation, Preflight, And Evidence

### A) Epic ID, Title, And Outcome

EP-00 - Plan Activation, Preflight, And Evidence

Outcome: The implementation starts from a known source state, with the plan activated only after
user approval and with execution evidence ready to capture validation and deviations.

### B) Scope

This epic covers activation of this implementation plan, initial source-repo preflight, public-core
classification, and execution-log setup. It does not edit CLI source behavior.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  docs/repo/plans/operating-kit-self-contained-bootstrap/
    operating-kit-self-contained-bootstrap_implementation_doc.md  # modify status during execution
    operating-kit-self-contained-bootstrap_execution_log.md        # create during execution
  docs/repo/plans/plan-register.md                                # modify lifecycle during execution
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- User has explicitly approved execution of this draft implementation plan.
- Source repo status is recorded before edits.
- Consumer impact classes are recorded.
- Execution log exists before implementation work starts.
- Public release publication is explicitly out of scope for this execution.
- Live public release pointers, root `manifest.yaml` asset URLs/checksums, and published
  bootstrap asset references are not switched to staged local assets during this execution.

### E) Dependencies And Critical-Path Notes

Dependencies: none.

Critical path: This epic must run before source changes so execution evidence is not reconstructed
after the fact.

### F) Tasks Checklist

- [x] Record the approved plan activation in `docs/repo/plans/operating-kit-self-contained-bootstrap/operating-kit-self-contained-bootstrap_execution_log.md`.
- [x] Run `git status --short` from `Codeheart-Operating-Kit` and record unrelated dirty paths in the execution log.
- [x] Classify the change as consumer migration required, validator-only change, instruction-only change, and security policy change in the execution log.
- [x] Confirm that public release publication, Git tags, GitHub release creation, and consumer sync remain outside this execution plan.
- [x] Confirm that root `manifest.yaml` and live bootstrap release references will not advertise staged unpublished asset URLs or checksums.
- [x] Update the implementation plan status from `draft` to `active` only after explicit user approval.

### G) Implementation Notes

Use `docs/repo/runbooks/change-operating-kit.md` as the source-change route. Protect existing
untracked `.codeheart/` and `docs/agent-memory/` paths in the source repo.

### H) Open Questions

None.

## EP-01 - Go CLI Skeleton

### A) Epic ID, Title, And Outcome

EP-01 - Go CLI Skeleton

Outcome: The source repo has a buildable Go module and CLI binary that exposes the current command
surface with help output and version output.

### B) Scope

This epic creates the Go module, command entry point, argument parsing, version wiring, and initial
tests for command availability. It does not implement command behavior beyond help and version.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  go.mod                                           # create
  go.sum                                           # create when dependencies require it
  cmd/codeheart-operating-kit/main.go              # create
  internal/cli/cli.go                              # create
  internal/cli/cli_test.go                         # create
  internal/version/version.go                      # create
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `go test ./...` passes for the new skeleton.
- `go run ./cmd/codeheart-operating-kit --help` lists `onboard`, `inspect`, `init`, `sync`,
  `check`, and `update-check`.
- `go run ./cmd/codeheart-operating-kit --version` prints the source version used by current
  release tooling.
- The Go module uses a public-safe module path.

### E) Dependencies And Critical-Path Notes

Dependencies: EP-00.

Critical path: EP-02 through EP-06 depend on a stable binary entry point.

### F) Tasks Checklist

- [x] Create `go.mod` with module path `github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit`.
- [x] Create `cmd/codeheart-operating-kit/main.go` that calls an internal CLI package and exits with the returned code.
- [x] Create `internal/version/version.go` with a version value populated from build flags and a source fallback.
- [x] Create `internal/cli/cli.go` with command parsing for `onboard`, `inspect`, `init`, `sync`, `check`, and `update-check`.
- [x] Create `internal/cli/cli_test.go` covering root help, subcommand help, missing command errors, and version output.
- [x] Run `go test ./...` and record the result in the execution log.

### G) Implementation Notes

Prefer Go standard-library packages for the first implementation. Add external dependencies only
when they materially reduce risk and are recorded in the execution log.

### H) Open Questions

None.

## EP-02 - Go Resource And State Foundation

### A) Epic ID, Title, And Outcome

EP-02 - Go Resource And State Foundation

Outcome: The Go CLI can read embedded Operating Kit resources, parse and write the supported YAML
subset, compute checksums, manage lock/config metadata, and resolve platform-specific release
asset metadata.

### B) Scope

This epic ports shared support behavior from the Python CLI. It does not implement user-facing
commands yet.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  internal/kitfs/kitfs.go                          # create embedded resource access
  internal/manifest/manifest.go                    # create component/profile/release helpers
  internal/yamlmini/yamlmini.go                    # create supported YAML parser/writer
  internal/lockfile/lockfile.go                    # create lock/config helpers
  internal/platforms/platforms.go                  # create platform helpers
  internal/hash/hash.go                            # create sha256 helper
  internal/*/*_test.go                             # create unit coverage
```

### D) Acceptance Criteria And Size

Size: L

Acceptance criteria:

- Go tests load the standard profile and component manifests from embedded resources.
- Go tests round-trip representative lock and config fixtures through the supported YAML subset.
- Go tests compute SHA-256 values matching Python helper outputs for fixture files.
- Go tests select macOS assets on Darwin and Windows assets on Windows with universal fallback
  where accepted by release metadata.

### E) Dependencies And Critical-Path Notes

Dependencies: EP-01.

Critical path: Command behavior in EP-03 relies on these helpers.

### F) Tasks Checklist

- [x] Create `internal/kitfs/kitfs.go` using Go `embed` for `components/`, `profiles/`, `templates/`, and required manifest resources.
- [x] Create `internal/yamlmini/yamlmini.go` with the same supported YAML subset currently used by `src/codeheart_operating_kit/manifest.py`.
- [x] Create `internal/manifest/manifest.go` to load profiles, components, component files, release manifests, and consumer impact records.
- [x] Create `internal/hash/hash.go` with streaming SHA-256 file hashing.
- [x] Create `internal/lockfile/lockfile.go` to read and write `.codeheart/kit.lock.yaml` and `.codeheart/kit.config.yaml`.
- [x] Create `internal/platforms/platforms.go` for macOS, Windows, and unsupported-platform asset resolution.
- [x] Add Go fixture tests for profile loading, component file enumeration, YAML round trip, checksum parity, and platform asset selection.
- [x] Run `go test ./...` and record the result in the execution log.

### G) Implementation Notes

The YAML parser should intentionally preserve the current supported subset rather than becoming a
general YAML implementation without tests. Use normalized fixture comparisons where map order can
differ.

### H) Open Questions

None.

## EP-03 - Root Command Behavior Port

### A) Epic ID, Title, And Outcome

EP-03 - Root Command Behavior Port

Outcome: The Go CLI implements the full current root command behavior for onboarding, inspection,
initialization, sync, check, and update-check, with Python parity baselines captured before each
command group is accepted.

### B) Scope

This epic ports behavior from the Python command modules. It preserves command names, argument
shape, stable user-facing output, JSON output, exit-code expectations, and generated state
semantics except for the approved base-onboarding difference: `onboard --yes` must not offer,
install, or implicitly check optional native capabilities.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  internal/commands/onboard.go                     # create
  internal/commands/inspect.go                     # create
  internal/commands/init.go                        # create
  internal/commands/sync.go                        # create
  internal/commands/check.go                       # create
  internal/commands/update_check.go                # create
  internal/components/components.go                # create managed/scaffold writing helpers
  internal/drift/drift.go                          # create drift detection
  internal/capabilities/capabilities.go            # create baseline capability metadata
  internal/commands/*_test.go                      # create command tests
```

### D) Acceptance Criteria And Size

Size: XL

Acceptance criteria:

- `inspect` classifies missing, empty, existing, technical, existing-kit, many-file, and file-path
  targets like the Python CLI.
- `onboard` enforces explicit `--target` and `--project-name` requirements with `--yes`.
- `onboard --yes` writes base Operating Kit setup only; it does not run native capability
  installation, does not offer optional tool installation, and does not implicitly run optional
  capability checks.
- `init` writes the same managed files, scaffolds, local user files, lock/config fields, and
  `.gitignore` entries with timestamp normalization.
- `sync` refreshes managed files, preserves consumer-owned scaffolds, merges generated surfaces,
  refreshes `AGENTS.md` managed blocks, and updates release metadata like the Python CLI.
- `check` reports missing CLI, stale CLI, missing routing, missing route targets, missing lock
  metadata, drift, and native capability state.
- `update-check` handles latest metadata, injected latest version, injected metadata URL, current
  state, update-available state, and failed metadata lookup state.

### E) Dependencies And Critical-Path Notes

Dependencies: EP-02.

Critical path: Each command group must capture Python baseline fixtures before or during porting.
EP-04 finalizes the full parity suite after all root commands exist.

### F) Tasks Checklist

- [x] Capture Python baseline fixtures for root help, subcommand help, `inspect`, and `onboard` before accepting the corresponding Go ports.
- [x] Port `inspect` behavior from `src/codeheart_operating_kit/commands/inspect.py` into `internal/commands/inspect.go`.
- [x] Port `onboard` prompt-rendering and `--yes` gating behavior into `internal/commands/onboard.go`, excluding optional native capability install/offer/check behavior by approved decision.
- [x] Add Go tests proving `onboard --yes` does not invoke plugin installation, does not offer optional tool installation, and does not run implicit optional capability checks.
- [x] Capture Python baseline fixtures for `init`, `sync`, `check`, and `update-check` before accepting the corresponding Go ports.
- [x] Port `init` file creation, managed resource copy, scaffold copy, local user layer, lock/config writing, and `.gitignore` behavior into `internal/commands/init.go`.
- [x] Port `sync` managed refresh, scaffold preservation, `AGENTS.md` managed block refresh, generated-surface merge, release metadata refresh, and `.gitignore` behavior into `internal/commands/sync.go`.
- [x] Port `check` route-target scanning, lock metadata checks, drift checks, stale CLI checks, and native capability state reporting into `internal/commands/check.go`.
- [x] Port `update-check` latest-release metadata lookup, injected metadata URL behavior, failed lookup state, and agent-notification silence behavior into `internal/commands/update_check.go`.
- [x] Add Go unit tests for every command's success, JSON output, stable text output, and representative error cases.
- [x] Run `go test ./...` and record the result in the execution log.

### G) Implementation Notes

Use the Python tests as behavior documentation, not merely as code examples. Preserve generated
consumer paths and ownership behavior exactly unless the implementation log records an approved
semantic difference. The optional native capability behavior in AD-006 is already approved as an
intentional difference and must be encoded in tests instead of treated as a parity failure.

### H) Open Questions

None.

## EP-04 - Behavior Parity Fixture Suite

### A) Epic ID, Title, And Outcome

EP-04 - Behavior Parity Fixture Suite

Outcome: A parity test suite proves the Go CLI matches the current Python CLI across root command
behavior before release assets switch to the Go binary, with approved intentional differences
documented as explicit fixture expectations.

### B) Scope

This epic finalizes comparison fixtures and test harnesses after incremental baseline capture in
EP-03. It does not change installer behavior.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  tests/
    fixtures/
      parity/                                      # create command fixtures
    test_go_cli_parity.py                          # create
    conftest.py                                    # modify shared helpers when useful
```

### D) Acceptance Criteria And Size

Size: L

Acceptance criteria:

- Python CLI and Go CLI are both executable in the parity harness.
- Python baseline fixtures for each command group were captured before accepting the Go command
  ports.
- Parity tests normalize timestamps, path separators, map ordering, and platform-specific path
  text.
- Parity tests cover `onboard`, `inspect`, `init`, `sync`, `check`, and `update-check`.
- Every intentional difference is recorded in the execution log and implemented as an explicit
  approved fixture expectation.

### E) Dependencies And Critical-Path Notes

Dependencies: EP-03.

Critical path: EP-05 and EP-06 must not replace release assets before this suite passes.

### F) Tasks Checklist

- [x] Create `tests/fixtures/parity/` with folder-state fixtures for missing folder, empty folder, existing folder, technical project, existing kit, ambiguous file path, and high-entry folder.
- [x] Create `tests/test_go_cli_parity.py` with helpers that invoke the Python CLI and compiled Go CLI against isolated temp folders.
- [x] Add parity tests for root help, subcommand help, version output, invalid command output, and `onboard --yes` argument gating.
- [x] Add explicit approved-difference tests proving Go `onboard --yes` does not offer, install, or implicitly check optional native capabilities.
- [x] Add parity tests for `inspect` mode, reason, marker, and JSON output.
- [x] Add parity tests for `init` generated file tree, scaffold preservation, `.gitignore`, `AGENTS.md`, normalized lockfile, and normalized config.
- [x] Add parity tests for `sync` managed-file repair, scaffold preservation, generated-surface merge, release metadata, and managed block refresh.
- [x] Add parity tests for `check` drift, missing route target, missing lock metadata, stale CLI, and JSON output.
- [x] Add parity tests for `update-check` current, update-available, failed metadata lookup, injected latest version, injected metadata URL, and agent-notification behavior.
- [x] Run the parity suite and record the result in the execution log.

### G) Implementation Notes

Semantic parity is the acceptance standard. Normalize timestamps and path separators before
comparison. Do not normalize away actual field presence, ownership mode, generated path, checksum,
status, or exit-code differences. The only accepted non-parity behavior in this plan is the
AD-006 onboarding optional-capability boundary unless the user explicitly approves another
difference during execution.

### H) Open Questions

None.

## EP-05 - Installers And Legacy Migration

### A) Epic ID, Title, And Outcome

EP-05 - Installers And Legacy Migration

Outcome: macOS and Windows installers install binary release packs, verify checksums, repair stale
or broken installs, and migrate legacy Python-wheel installs without requiring Python or pip.

### B) Scope

This epic changes platform installers and installer tests. It does not build final release packs;
that belongs to EP-06.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  install.sh                                       # modify
  install.ps1                                      # modify
  tests/
    fixtures/
      installer/                                   # create legacy install fixtures
    test_install_metadata.py                       # modify
    test_release_assets.py                         # modify installer mismatch coverage
```

### D) Acceptance Criteria And Size

Size: L

Acceptance criteria:

- macOS installer downloads or copies a platform zip, verifies SHA-256, extracts the Go binary,
  installs it under `$HOME/.codeheart/operating-kit/bin/codeheart-operating-kit`, and prints PATH
  guidance.
- Windows installer downloads or copies a platform zip, verifies SHA-256, extracts the Go binary,
  installs it under `%LOCALAPPDATA%\Codeheart\OperatingKit\bin`, writes a `.cmd` shim, and prints
  PATH guidance.
- Both installers preserve `--asset-file`, `file://`, `--checksum`, and `--checksum-file`
  behavior.
- Both installers fail closed on checksum mismatch.
- Both installers extract into a staging location, verify the expected binary path, run a staged
  `--version` or `--help` smoke check, and only then replace the runnable command path.
- If staged binary validation fails, the previous runnable command remains intact.
- Legacy Python wrapper/lib state is detected and migrated by installing the Go binary as the new
  runnable command while preserving legacy files.
- `--python` / `-Python` no longer appears as normal installer help or behavior. If retained for
  one-release compatibility, it is accepted only as a deprecated no-op with a warning, and tests
  prove it is never invoked.

### E) Dependencies And Critical-Path Notes

Dependencies: EP-04.

Critical path: EP-06 release packs depend on installer expectations.

### F) Tasks Checklist

- [x] Update `install.sh` to install `codeheart-operating-kit-<version>-macos-universal.zip` without invoking Python, pip, wheel, or tarball wheel extraction.
- [x] Update `install.sh` to extract `bin/codeheart-operating-kit` into a staging directory, set executable mode, smoke-test the staged binary, and atomically replace the runnable command only after validation passes.
- [x] Update `install.sh` to detect legacy Python wrapper and legacy `lib/codeheart_operating_kit*` paths and print migration state while preserving those paths.
- [x] Remove `--python` from `install.sh` help and behavior, or keep it only as a deprecated no-op compatibility flag with a warning.
- [x] Update `install.ps1` to install `codeheart-operating-kit-<version>-windows-x64.zip` without invoking Python, pip, wheel, or Python module wrappers.
- [x] Update `install.ps1` to extract `bin/codeheart-operating-kit.exe` into a staging directory, smoke-test the staged binary, and atomically replace the runnable command and `.cmd` shim only after validation passes.
- [x] Update `install.ps1` to detect legacy Python wrapper and legacy `lib\codeheart_operating_kit*` paths and print migration state while preserving those paths.
- [x] Remove `-Python` from `install.ps1` help and behavior, or keep it only as a deprecated no-op compatibility parameter with a warning.
- [x] Preserve existing asset-file, checksum-file, checksum override, and `file://` local asset flows in both installers.
- [x] Add installer fixture tests covering no install, legacy Python install, checksum mismatch, local asset install, malformed archive with valid checksum, missing binary, non-executable binary, failed staged validation preserving the previous runnable command, deprecated Python flag behavior when retained, and successful binary install metadata.
- [x] Run installer-focused tests and record the result in the execution log.

### G) Implementation Notes

Do not delete legacy Python libraries during the normal migration path. A later cleanup command can
be planned after the Go release proves stable. Migration is a repair path, so successful checksum
validation is not sufficient by itself; staged binary smoke validation must pass before replacing
any existing runnable command.

### H) Open Questions

None.

## EP-06 - Binary Release Asset Builder And Manifest

### A) Epic ID, Title, And Outcome

EP-06 - Binary Release Asset Builder And Manifest

Outcome: The release asset builder produces macOS and Windows binary release packs with checksums
and manifests, and release validation no longer expects a Python wheel in consumer assets.

### B) Scope

This epic updates staging assets, release-candidate manifest schema/fixtures, and release builder
tests. It does not publish assets and does not switch the live root `manifest.yaml` to unpublished
asset URLs or staged local checksums.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  scripts/build-release-assets.py                  # modify
  scripts/validate-release-manifest.py             # modify when schema expectations change
  schemas/release-manifest.schema.json             # modify
  tests/fixtures/release-manifest.json             # modify
  tests/fixtures/release-candidate/                # create staged manifest fixtures
  tests/fixtures/validator-invalid/                # modify invalid release fixtures
  tests/test_release_assets.py                     # modify
  dist/                                            # generated staged assets and release-candidate manifest
```

### D) Acceptance Criteria And Size

Size: L

Acceptance criteria:

- `scripts/build-release-assets.py` builds the Go CLI for `macos-universal` and `windows-x64`.
  `macos-universal` is mandatory for this asset name.
- Release packs are named `codeheart-operating-kit-<version>-macos-universal.zip` and
  `codeheart-operating-kit-<version>-windows-x64.zip`.
- Release packs include `bootstrap.md`, `INSTALL.md`, platform binary, manifest data, and
  checksum data.
- Release packs do not contain `*.whl`, `codeheart_operating_kit-*.dist-info`, or Python package
  payload directories.
- Release-candidate manifests and fixtures advertise the new asset names, platform values, and
  checksums.
- Root `manifest.yaml` is not updated with staged local checksums or unpublished asset URLs during
  this implementation plan.
- Release asset tests verify binary presence, no-wheel payload, checksum files, manifest values,
  and version mismatch failure.

### E) Dependencies And Critical-Path Notes

Dependencies: EP-05.

Critical path: EP-08 staged install proof depends on these assets.

### F) Tasks Checklist

- [x] Update `scripts/build-release-assets.py` to run `go test ./...` before asset assembly.
- [x] Update `scripts/build-release-assets.py` to build `cmd/codeheart-operating-kit` for mandatory macOS universal and Windows x64 with version build flags.
- [x] Update `scripts/build-release-assets.py` to assemble platform zip packs containing `bootstrap.md`, `INSTALL.md`, `bin/`, `manifest.json`, and `checksums.txt`.
- [x] Update `scripts/build-release-assets.py` to write SHA-256 sidecar files for every platform pack.
- [x] Update `schemas/release-manifest.schema.json` and release-manifest fixtures for `macos-universal` and `windows-x64` asset platform names.
- [x] Update `tests/test_release_assets.py` to reject Python wheel payloads inside staged release packs.
- [x] Update `tests/test_release_assets.py` to validate binary names, checksum sidecars, manifest data, and version mismatch failure.
- [x] Generate staged assets and a release-candidate manifest under `dist/` without changing root `manifest.yaml` to staged local checksum values.
- [x] Run release asset and manifest tests and record the result in the execution log.

### G) Implementation Notes

Mac universal output requires arm64 and amd64 binaries combined with `lipo` or an equivalent
validated universal-binary process. If universal output cannot be produced, stop and revise the
platform matrix before continuing. When the local machine cannot produce the Windows binary
directly, use the CI path in EP-08 as the authoritative Windows build proof.

### H) Open Questions

None.

## EP-07 - Bootstrap, Release, And Documentation Updates

### A) Epic ID, Title, And Outcome

EP-07 - Bootstrap, Release, And Documentation Updates

Outcome: Source and maintainer documentation describe the self-contained bootstrap, migration
behavior, trust gates, and release-stage boundaries without implying that optional tooling is
installed during first-run onboarding or that staged assets are already live public releases.

### B) Scope

This epic updates release-ready bootstrap source text, maintainer docs, release guidance, release
notes, and plan indexes. It does not publish the versioned bootstrap asset, switch live release
URLs, or change command behavior.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  README.md                                        # modify
  bootstrap.md                                     # modify
  release-notes.md                                 # modify
  docs/repo/README.md                              # modify plan pointer
  docs/repo/plans/README.md                        # modify plan pointer
  docs/repo/runbooks/release-operating-kit.md      # modify binary release and signing gates
  docs/repo/reference/consumer-impact-classification.md # modify only when migration wording needs a new note
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `bootstrap.md` source text tells agents to install the self-contained CLI from platform release
  packs and no longer mentions Python/pip installation.
- `bootstrap.md` and release notes do not claim unpublished staged assets are live public release
  assets.
- `README.md` maintainer entry points describe the Go CLI source and legacy Python parity oracle.
- `release-operating-kit.md` includes binary asset validation, checksum mismatch validation,
  staged unsigned asset boundary, live root-manifest/bootstrap switch steps, and signed/notarized
  broad release gate.
- `release-notes.md` includes consumer-impact notes for self-contained CLI migration.
- Docs do not tell users to install Homebrew, Git, GitHub CLI, Python, pip, Node.js, or package
  managers during base onboarding.
- Docs do not offer optional native capability installation during base onboarding.

### E) Dependencies And Critical-Path Notes

Dependencies: EP-06.

Critical path: EP-08 probes use the release-candidate bootstrap source text and staged assets, not
unpublished assets pretending to be a live public release.

### F) Tasks Checklist

- [x] Update `bootstrap.md` to describe macOS and Windows self-contained release pack installation.
- [x] Update `bootstrap.md` to preserve the existing agent-guided onboarding contract and approval gates.
- [x] Ensure `bootstrap.md` does not offer optional native capability installation during base onboarding.
- [x] Ensure `bootstrap.md` and release notes do not advertise staged local asset URLs or checksums as live public release assets.
- [x] Update `README.md` maintainer entry points to include Go CLI source paths and staged binary release asset behavior.
- [x] Update `docs/repo/runbooks/release-operating-kit.md` with binary pack build, checksum mismatch proof, macOS staged install proof, Windows staged install proof, live root-manifest/bootstrap switch steps, unsigned internal prototype boundary, and signed broad release gate.
- [x] Update `release-notes.md` with self-contained CLI, legacy migration, and no-Python bootstrap consumer notes.
- [x] Update `docs/repo/README.md` and `docs/repo/plans/README.md` with this implementation plan path.
- [x] Run Markdown header validation and public-core validation and record the result in the execution log.

### G) Implementation Notes

Keep first-run user copy simple. Do not add broad environment setup guidance to `bootstrap.md`.
Base onboarding should finish after Operating Kit setup; optional tool readiness belongs to later
task-specific routes.

### H) Open Questions

None.

## EP-08 - Cross-Platform Validation And Staged Install Proof

### A) Epic ID, Title, And Outcome

EP-08 - Cross-Platform Validation And Staged Install Proof

Outcome: Local and CI validation prove the Go CLI, release packs, installers, and parity suite work
on macOS and Windows before a public release run is considered.

### B) Scope

This epic updates CI and runs validation. It does not publish a release.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  .github/workflows/validate.yml                   # modify
  tests/test_release_assets.py                     # modify when CI helpers need markers
  docs/repo/plans/operating-kit-self-contained-bootstrap/
    operating-kit-self-contained-bootstrap_execution_log.md # update during execution
```

### D) Acceptance Criteria And Size

Size: L

Acceptance criteria:

- macOS validation runs Go tests, Python parity tests, release asset tests, Markdown validation,
  public-core validation, schema validation, and staged installer proof.
- Windows validation builds the Windows x64 binary, runs Go tests, runs applicable parity tests,
  builds the Windows release pack, and runs `install.ps1` against a staged asset.
- macOS staged install proof runs with `python`, `python3`, `pip`, and `pip3` removed from `PATH`
  in the installer process.
- Windows staged install proof runs with Python removed from `PATH` in the installer process.
- Checksum mismatch tests fail closed on both platforms.
- Validation proves base onboarding does not offer, check, or install optional native capabilities.
- Fresh low-context bootstrap probe evidence is recorded.

### E) Dependencies And Critical-Path Notes

Dependencies: EP-07.

Critical path: EP-09 handoff depends on this validation evidence.

### F) Tasks Checklist

- [x] Update `.github/workflows/validate.yml` to install Go and run `go test ./...` on macOS and Windows.
- [x] Update `.github/workflows/validate.yml` to run Python parity tests after building the Go CLI.
- [x] Update `.github/workflows/validate.yml` to build staged macOS and Windows release packs.
- [x] Add a macOS staged install validation command that runs `install.sh --asset-file <pack> --checksum-file <pack.sha256>` with Python commands absent from `PATH`.
- [x] Add a Windows staged install validation command that runs `install.ps1 -AssetFile <pack> -ChecksumFile <pack.sha256>` with Python commands absent from `PATH`.
- [x] Add checksum mismatch validation for `install.sh` and `install.ps1`.
- [x] Add validation that `onboard --yes` performs base setup without native capability install offers, native capability checks, or plugin installation attempts.
- [x] Run the full local validation set available on macOS and record the result in the execution log.
- [x] Record Windows CI validation evidence in the execution log.
- [x] Record a fresh low-context bootstrap probe in the execution log using the public first prompt, release-candidate bootstrap source text, and staged release asset.

### G) Implementation Notes

Use CI for Windows proof. Do not treat absence of a local Windows shell as a blocker when GitHub
Actions provides the Windows validation lane.

### H) Open Questions

None.

## EP-09 - Register, Release-Readiness Record, And Handoff

### A) Epic ID, Title, And Outcome

EP-09 - Register, Release-Readiness Record, And Handoff

Outcome: The source implementation is review-ready with plan/register state, release-readiness
notes, validation evidence, residual risks, and an explicit stop before public release
publication.

### B) Scope

This epic updates plan/register state and records final handoff. It does not create a Git tag,
publish a GitHub release, or sync consumers.

### C) Files Touched

```text
Codeheart-Operating-Kit/
  docs/repo/plans/plan-register.md                 # modify OK-PR-024 lifecycle
  docs/repo/plans/operating-kit-self-contained-bootstrap/
    operating-kit-self-contained-bootstrap_implementation_doc.md # modify lifecycle during execution
    operating-kit-self-contained-bootstrap_execution_log.md       # update
  release-notes.md                                 # final release-readiness note
Codeheart-HQ/
  docs/repo/plans/plan-register.md                 # coordination pointer update when source plan completes
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- Source `plan-register.md` contains `OK-PR-024` with correct lifecycle and validation summary.
- HQ coordination register points from `CODEHEART-HQ-PR-009` to `OK-PR-024`.
- Execution log records every validation command, staged asset path, checksum, migration proof,
  residual risk, and release publication stop.
- Implementation plan remains `draft` until user approves execution, then moves through `active`
  and completion during execution.
- Public release publication is left as a separate approval-gated action.
- Live root `manifest.yaml` asset URLs/checksums and published bootstrap asset references remain
  release-run work.

### E) Dependencies And Critical-Path Notes

Dependencies: EP-08.

Critical path: This epic closes source implementation readiness.

### F) Tasks Checklist

- [x] Update `docs/repo/plans/plan-register.md` with `OK-PR-024` lifecycle, validation summary, and release publication stop.
- [x] Update `Codeheart-HQ/docs/repo/plans/plan-register.md` so `CODEHEART-HQ-PR-009` references `OK-PR-024`.
- [x] Record staged asset names and SHA-256 values in the execution log.
- [x] Record macOS install proof, Windows install proof, parity proof, and checksum mismatch proof in the execution log.
- [x] Record proof that base onboarding does not offer, check, or install optional native capabilities.
- [x] Record that root `manifest.yaml` live asset URLs/checksums and published bootstrap asset references were not switched to staged unpublished assets.
- [x] Record residual risks for signing/notarization automation and broad public release timing.
- [x] Stop before `git tag`, `gh release`, GitHub release publication, and consumer sync.

### G) Implementation Notes

Public release should run later through `docs/repo/runbooks/release-operating-kit.md` after explicit
release approval and signing/notarization readiness review. That release run owns the live
root-manifest checksum switch, versioned bootstrap publication, Git tag, GitHub release, and
consumer sync.

### H) Open Questions

None.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Public release publication: deferred because publication changes external GitHub state and
  requires explicit release approval.
- Live root `manifest.yaml` and versioned public bootstrap publication: deferred to the release run
  so source implementation cannot advertise unpublished assets as live.
- Broad signed/notarized release automation: deferred until certificate, account, CI secret, and
  policy details are available.
- Linux and Windows ARM support: deferred until a real consumer need appears.
- Removal of Python CLI source: deferred until at least one self-contained Go release proves
  stable and legacy migration has been validated in consumer installs.
- Legacy Python library cleanup command: deferred until migration evidence shows cleanup is safe
  and useful.

## 4.2 Future Considerations

- A future release plan should publish the validated assets, update live root `manifest.yaml` with
  public GitHub release URLs and checksums, publish the versioned bootstrap asset, record release
  evidence, and sync approved consumer repositories.
- A future implementation may move remaining source validation scripts from Python to Go only when
  maintainer runtime reduction becomes a real need.
- A future installer may add signed binary verification metadata once signing automation exists.
- A future consumer proof should include a truly fresh non-developer Mac and a fresh Windows VM.

# Revision Notes

- 2026-07-04: Initial draft from the approved HQ discovery capability scope and targeted
  Operating Kit source reconnaissance.
- 2026-07-04: Revised after planning review and user approval to clarify release-candidate
  boundaries, setup-only onboarding for optional native capabilities, incremental parity,
  atomic installer migration, mandatory macOS universal assets, and Python installer flag removal
  or deprecation behavior.
