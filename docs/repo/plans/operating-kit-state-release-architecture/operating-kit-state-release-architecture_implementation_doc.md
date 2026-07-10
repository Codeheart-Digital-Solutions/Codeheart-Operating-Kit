Last updated: 2026-07-10T00:45:48Z (UTC)
Created: 2026-07-09
Status: active
Source state: complete; validation-only Windows CI authorized and pending
Execution log: operating-kit-state-release-architecture_execution_log.md

# Document Header

## Operating Kit State And Release Architecture - Minimum Reliable Foundation Plan

Overview: Implement the approved state-and-release capability as a minimum reliable foundation,
not as a new governance framework. The plan improves how the Operating Kit installs, checks,
repairs, synchronizes, upgrades, and releases itself. It does not attempt a broad redesign of
Operating Kit guidance or add ceremony to ordinary agent work.

The implementation uses three cohesive internal packages: state, reconciliation, and release. It
adds schemas only for durable public data, keeps successful transaction evidence transient, uses
automated tests and CI as the primary proof, and introduces one compact lifecycle runbook. Explicit
approval is required for cross-version upgrade; invoking `init`, `repair`, or `sync` is itself the
authorization for that named operation. Healthy routine results remain concise and update checks
remain silent when current.

This plan preserves the approved capability scope: lock schema v1 read compatibility, lock schema
v2 writes through an authorized lifecycle command, config schema v1 preservation, transactional
mutation, rollback, non-circular release provenance, reproducible packs, safe upgrade, producer
routing, and real macOS/Windows validation. It stops before release publication, tag creation,
named consumer sync, Python retirement, signing, new platforms, and the available local kit update.

Essential context:

| Source | Why it matters |
| --- | --- |
| `operating-kit-state-release-architecture_discovery_doc.md` | Approved decisions, requirements, risks, and implementation capability scopes. |
| `AGENTS.md` | Repository safety, managed-content, routing, and producer-work rules. |
| `docs/repo/reference/placement-contract.md` | Authority for source, generated, local-user, local-machine, and packaged-resource placement. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact classes for state migration, validation, instructions, and safety behavior. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer procedure for CLI, schema, template, installer, and release changes. |
| `docs/repo/runbooks/release-operating-kit.md` | Release authority and publication stop conditions. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Route ownership and low-context routing proof. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Recipe maturity, evidence, blocker, and validation expectations. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Audience and execution requirements for the compact lifecycle runbook. |
| `internal/commands/`, `internal/components/`, `internal/lockfile/`, and `internal/manifest/` | Current command, installation, state, and declaration behavior to consolidate. |
| `internal/yamlmini/yamlmini.go` | Custom YAML subset with the known scalar-typing defect. |
| `schemas/kit-lock.schema.json` and `schemas/kit-config.schema.json` | Current installed-state contracts. |
| `resources.go`, `manifest.yaml`, and `scripts/build-release-assets.py` | Current circular content/archive provenance and nondeterministic release builder. |
| `install.sh`, `install.ps1`, and `.github/workflows/validate.yml` | macOS and Windows installation and failure-validation surfaces. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Make the Operating Kit reliable at managing its own installed state without materially increasing
day-to-day user or agent ceremony.

Completion is proven when:

- component, profile, lock, config, embedded content identity, release catalog, and pack manifest
  data are parsed with a maintained YAML library and validated against versioned schemas;
- validated declarations resolve into one deterministic desired-state graph used by checking and
  mutation;
- one read-only classifier identifies absent, adoptable, current, drifted, stale CLI, partial,
  schema-invalid, legacy-v1-compatible, transaction-in-progress, recovery-required, and
  unsupported-future-version state;
- `init`, `repair`, `sync`, `update-check`, `upgrade`, and `check` enforce distinct responsibilities;
- mutating commands support dry-run, preserve consumer-owned content, stage changes on the target
  filesystem, commit lock v2 last, and roll back on failure;
- known released lock-v1 checksum-placeholder anomalies have a narrow migration path while
  unrelated invalid state still fails closed;
- transaction staging and recovery evidence are removed after success and retained only after a
  failed or interrupted operation;
- `check` validates the complete installed state and returns concise, actionable blockers;
- embedded content identity no longer contains its enclosing archive digest;
- release candidates are reproducible and installers verify catalog, archive, pack, platform,
  version, command, payload, and binary identity;
- upgrade requires explicit approval, runs reconciliation through the verified new binary, and
  restores the prior binary plus repository state after failure;
- producer work routes to canonical source and two focused low-context probes prove producer and
  consumer lifecycle routing;
- local and cross-platform validation passes without creating a tag, publishing a release, syncing
  a named consumer, or applying the local kit update.

## 1.2 Project And Problem Context

The Operating Kit already has a tested Go CLI, component and profile declarations, JSON Schemas,
managed-content ownership, binary release packs, and macOS/Windows installers. The reliability
problem is concentrated in its self-management path:

- existing-kit onboarding can re-enter initialization;
- sync writes directly and can leave a partial installation;
- check validates only part of the state;
- the custom YAML subset changes some scalar types;
- declarations and runtime hardcoding can disagree;
- lock-v1 consumers need a controlled migration;
- embedded release metadata depends on the archive containing it;
- release packs are not proven byte-reproducible;
- installers do not verify the complete pack identity;
- this producer repository can route through stale ignored consumer materialization.

These are real runtime and release defects. The implementation does not need a large new process
layer to fix them. It needs a small number of durable contracts, cohesive code boundaries, strong
filesystem safety, and automated validation.

## 1.3 Current State Analysis

Current systems and constraints:

- command implementations combine preconditions, state interpretation, mutation, and rendering;
- component/profile YAML describes useful ownership data but does not fully drive runtime behavior;
- lock/config reads use `yamlmini`, and full schema instance validation is absent from the Go CLI;
- lock v1 has no completed state generation or verified external provenance;
- sync has no deterministic change plan, target-filesystem staging, exclusive operation marker, or
  rollback boundary;
- successful command JSON shapes are not consistently derived from one shared model;
- `manifest.yaml` combines embedded content facts with external archive facts;
- ZIP output inherits build-host metadata;
- installer and Windows failure coverage is incomplete;
- managed routing docs exist, but producer authority is not explicit.

Target systems and boundaries:

- `internal/state/` owns schema-backed parsing, deterministic encoding, desired-state resolution,
  observed-state classification, and lock migration;
- `internal/reconcile/` owns change planning, target containment, exclusive operation state,
  staging, commit, rollback, recovery, and shared operation results;
- `internal/release/` owns content identity, external catalog, pack verification, reproducible build
  contracts, and staged upgrade;
- public commands remain thin adapters over those packages;
- only durable externally consumed data receives JSON Schema files;
- `.codeheart/kit.transaction.json` and `.codeheart/local/kit-transactions/` exist only during an
  active or failed operation and are deleted after successful completion;
- lock v2 stores the current generation, last successful operation summary, and verified release
  provenance, not a verbose historical journal;
- tests and CI produce detailed evidence; users and agents receive short success output and
  actionable blockers;
- one compact managed lifecycle runbook owns command selection and recovery guidance.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement the safety-critical path in five linear epics: state, lifecycle transactions, release
and upgrade, compact routing, then integrated validation. Keep existing package boundaries as
compatibility adapters only where removing them immediately would create avoidable churn.

Expected source hierarchy:

```text
Codeheart-Operating-Kit/
  AGENTS.md                                               # modify producer route
  README.md                                               # modify lifecycle overview
  bootstrap.md                                            # modify verified install/upgrade flow
  go.mod                                                  # modify two direct parsing dependencies
  go.sum                                                  # create
  manifest.yaml                                           # modify embedded content identity
  resources.go                                            # modify embedded-resource boundary
  install.sh                                              # modify verification and rollback
  install.ps1                                             # modify verification and rollback
  internal/
    state/                                                # create: contracts, graph, classifier, migration
      schema.go
      yaml.go
      graph.go
      observed.go
      migrate.go
      state_test.go
    reconcile/                                            # create: plan, filesystem transaction, result
      plan.go
      transaction.go
      result.go
      reconcile_test.go
    release/                                              # create: catalog, pack, reproducibility, upgrade
      catalog.go
      pack.go
      upgrade.go
      handoff.go
      release_test.go
    cli/cli.go                                             # modify command surface
    commands/                                              # modify adapters; create repair.go and upgrade.go
    components/components.go                               # reduce to state-backed adapter
    lockfile/lockfile.go                                   # reduce to state-backed adapter
    manifest/manifest.go                                   # reduce to state-backed adapter
  schemas/
    component.schema.json                                 # create
    profile.schema.json                                   # create
    kit-lock-v1.schema.json                               # create frozen compatibility schema
    kit-lock.schema.json                                  # modify to v2
    kit-config.schema.json                                # retain v1 structure
    content-manifest.schema.json                          # create
    release-catalog.schema.json                           # create
    pack-manifest.schema.json                             # create
  components/agent-interface/
    component.yaml                                        # modify managed inventory
    managed/
      kit-readme.md                                       # modify compact command routes
      reference/update-check-policy.md                    # modify valid-install behavior
      runbooks/
        conduct-first-run-onboarding.md                   # modify init/repair selection
        maintain-operating-kit-installation.md            # create single lifecycle runbook
  templates/agents/AGENTS.managed-block.md                # modify consumer routes
  src/codeheart_operating_kit/resources/                  # mirror changed retained resources
  scripts/
    build-release-assets.py                               # modify deterministic builder
    validate-json-schemas.py                              # modify durable-contract validation
    validate-release-manifest.py                          # modify content/catalog split
  tests/
    fixtures/state/                                       # create focused state and recovery fixtures
    fixtures/release/                                     # create focused pack/catalog fixtures
    test_routing.py                                       # create two routing scenarios
    test_json_schemas.py                                  # modify
    test_sync_check.py                                    # modify
    test_release_assets.py                                # modify
    test_install_metadata.py                              # modify
  .github/workflows/validate.yml                          # modify macOS/Windows validation
  docs/repo/
    runbooks/change-operating-kit.md                      # modify validation gates
    runbooks/release-operating-kit.md                     # modify release gates
    plans/operating-kit-state-release-architecture/
      operating-kit-state-release-architecture_discovery_doc.md
      operating-kit-state-release-architecture_implementation_doc.md
      operating-kit-state-release-architecture_execution_log.md # create only on activation
```

Generated consumer transaction state:

```text
<consumer>/.codeheart/
  kit.transaction.json                                   # active/failed operation marker only
  local/kit-transactions/<transaction-id>/
    plan.json                                             # intended changes
    stage/                                                # staged replacement bytes
    backup/                                               # affected prior bytes
```

The operation marker is atomically created as the per-target lock. Success removes the marker,
plan, stage, and backup. Failed rollback retains the minimum files required for `repair` to recover.
No completed journal archive is added.

Recipe and routing treatment:

| Surface | Treatment |
| --- | --- |
| Lifecycle CLI | L3 command surface over reusable Go packages; unit, fixture, dry-run, and filesystem failure tests; concise structured result with stable Go/golden-test contract. |
| `maintain-operating-kit-installation.md` | One L1 agent-facing recipe calling the CLI; compact intention, command table, approval boundary, blocker table, and recovery path. |
| Release builder | Existing L2 workflow script; deterministic phases and JSON summary remain in the script and tests; no new wrapper, script index, or service. |
| Producer/consumer routes | Two focused low-context probes recorded in the execution log; no separate probe attachment or broad routing audit. |

The applicable installed references are
`.codeheart/kit/docs/agent-interface/reference/operational-recipe-maturity.md` and
`.codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md`. Producer work
uses their canonical source under `components/agent-interface/managed/reference/`. Missing Go,
Python, `lipo`, shell, or PowerShell remains a local tooling-readiness blocker. GitHub availability
and publication authority remain release-runbook preflight.

## 2.2 Open Questions And Assumptions Requiring Clarification

### OQ-1 - Future Public Signing Policy

`BLOCKER: no`

Affects: `EP-03`, `EP-05`.

Decision unlocked: catalog signatures, attestations, macOS notarization, and Windows signing.

Recommended default: keep the approved unsigned internal/prototype boundary, reserve compatible
catalog fields, reject unsupported signature claims, and stop before public release.

### OQ-2 - Python Decommission

`BLOCKER: no`

Affects: `EP-05` and future work.

Decision unlocked: deletion of the Python CLI, resource mirror, and parity suite.

Recommended default: retain Python compatibility evidence, add no new Python lifecycle behavior,
and make new capability authoritative in Go tests.

### OQ-3 - Publication Version

`BLOCKER: no`

Affects: future release execution after `EP-05`.

Decision unlocked: version bump, final catalog URLs, tag, and release notes.

Recommended default: do not select a publication version in source implementation. Build isolated
release candidates for validation and leave versioning to separately approved release work.

Assumptions:

- existing lock, config, root instruction, and managed kit paths remain stable;
- config schema v1 remains sufficient; a discovered structural config need stops work and returns
  to planning rather than creating config v2 implicitly;
- macOS universal and Windows x64 remain the supported platforms;
- GitHub Releases remains the initial catalog host;
- staging below the target repository provides same-filesystem replacement semantics;
- the source repository's current `.codeheart/` tree remains untouched and non-authoritative;
- no named consumer repository is modified during this plan.

## 2.3 Architectural Decisions With Reasoning

### AD-1 - Three Cohesive Internal Packages

Problem: the first draft introduced too many domain packages and boundaries.

Simplest working solution: use `internal/state`, `internal/reconcile`, and `internal/release`, with
public commands as adapters and existing packages retained only where compatibility requires them.

Six-to-twelve-month change: a package may split after measured growth or an independent caller
appears.

Rationale: the three packages follow the real lifecycle—understand state, change state, change
version—without building a framework around each substep.

Rejected alternative: separate contracts, graph, installation, lifecycle, catalog, and upgrade
packages create more coordination than the current scope needs.

### AD-2 - Schemas Only For Durable Public Data

Problem: public YAML needs reliable validation, but schemas for every internal record add ceremony.

Simplest working solution: use `go.yaml.in/yaml/v3 v3.0.4` and
`github.com/santhosh-tekuri/jsonschema/v6 v6.0.2`; schema-validate component, profile, lock, config,
content manifest, release catalog, and pack manifest data. Validate transaction markers and
operation results with strict Go structs plus golden fixtures.

Six-to-twelve-month change: publish an operation-result schema only after a real external consumer
needs one.

Rationale: durable interoperability gets executable contracts while internal mechanics remain
easy to change.

Rejected alternative: retaining `yamlmini` preserves known defects; schemas for transient markers
and every command result create maintenance without current consumer value.

### AD-3 - One Desired-State Graph And One Classifier

Problem: declarations, runtime lists, checking, and lock generation can disagree.

Simplest working solution: validated component/profile declarations compile into a sorted typed
graph; a read-only classifier compares that graph with observed state and returns one primary state
plus diagnostic traits.

Six-to-twelve-month change: additional profiles can add graph node types without changing command
preconditions.

Rationale: one runtime model improves correctness and reduces duplicated logic.

Rejected alternative: command-specific registries plus drift tests retain multiple authorities.

### AD-4 - Bounded Lock-V1 Migration

Problem: lock v2 is needed, but the kit's released v1 checksum placeholders can be schema-invalid.

Simplest working solution: freeze lock-v1 schema; before strict validation recognize only scalar
`0`, scalar `"0"`, and 64 zero digits at `release.checksum_sha256` and
`cli_repair.repair_checksum_sha256`; migrate them to explicit `unverified-legacy` provenance during
an authorized dry-run-visible operation. All unrelated invalid state fails closed.

Six-to-twelve-month change: lock v3 can use the same explicit version dispatch after evidence
justifies it.

Rationale: released consumers are recoverable without creating a general coercion mechanism.

Rejected alternative: strict rejection strands valid historical installs; broad normalization can
hide corruption.

### AD-5 - Small Transaction, No Successful Journal Archive

Problem: direct writes are unsafe, while a durable multi-phase evidence system would be excessive.

Simplest working solution: atomically create one active marker, stage affected bytes on the target
filesystem, back up affected prior bytes, validate, commit lock last, post-check, then delete all
transaction state after success. Retain recovery state only when rollback cannot finish.

Use a canonical root, no-follow platform primitives, symbolic-link and Windows reparse-point
refusal, identity revalidation before replacement, and live-process checks before stale-marker
takeover.

Six-to-twelve-month change: persistent audit storage can be added only after a real operational or
compliance requirement.

Rationale: this provides atomicity and recovery without turning every sync into a retained case
file.

Rejected alternative: direct writes remain unsafe; a permanent journal for every success adds
storage, code, and review burden.

### AD-6 - Explicit Commands, Minimal Approval Ceremony

Problem: overlapping command behavior is unsafe, but extra prompts can make routine maintenance
annoying.

Simplest working solution: keep explicit `init`, `repair`, `sync`, `update-check`, `upgrade`, and
`check`; support `--dry-run` on mutations; require `upgrade --yes` for a version change. Direct
invocation authorizes init, repair, and current-version sync. Existing onboarding confirmation
remains the approval for its selected operation.

Six-to-twelve-month change: an API can call the same packages without changing CLI semantics.

Rationale: command names communicate intent; only external version change needs an additional
approval flag.

Rejected alternative: one generic reconcile command hides authority; confirmation on every repair
or sync adds ceremony without improving informed consent.

### AD-7 - External Archive Catalog And Reproducible Packs

Problem: embedded archive digests are circular and current ZIP output is nondeterministic.

Simplest working solution: embed content identity only; generate the archive catalog after packs
exist; normalize entry metadata, build inputs, and Go build IDs; verify the full catalog-to-binary
chain in installers and upgrade.

Six-to-twelve-month change: signatures and attestations can extend the external catalog.

Rationale: provenance becomes acyclic and release candidates become comparable byte-for-byte.

Rejected alternative: removing provenance weakens safety; keeping the embedded enclosing digest
cannot produce a stable artifact.

### AD-8 - Compact Routing Proof

Problem: source maintainers can use stale installed doctrine, but a broad routing program would not
improve the core lifecycle implementation.

Simplest working solution: add producer authority to root instructions, update consumer command
routes, maintain one lifecycle runbook, and run one fresh producer probe plus one fresh consumer
probe.

Six-to-twelve-month change: broader guidance-effectiveness evaluation remains a separate workstream.

Rationale: the changed route is proven without expanding into a general audit of all Operating Kit
guidance.

Rejected alternative: relying on the source repository's ignored installation preserves hidden
drift; testing every runbook is outside this capability.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| `EP-01` | One validated state model, desired-state graph, classifier, and bounded lock migration exist. | XL | None |
| `EP-02` | Lifecycle commands use one small transactional reconciler with dry-run, rollback, recovery, and concise results. | XL | `EP-01` |
| `EP-03` | Release packs are reproducible and catalog-backed install/upgrade is verifiable and reversible. | XL | `EP-01`, `EP-02` |
| `EP-04` | Producer and consumer routing plus one compact lifecycle runbook accurately expose the implemented behavior. | M | `EP-02`, `EP-03` |
| `EP-05` | Automated local and cross-platform evidence proves the capability and stops before release or consumer sync. | L | `EP-01` through `EP-04` |

## Epic EP-01 - Validated State Model And Migration

### A) Epic ID, Title, And Outcome

`EP-01` - Validated State Model And Migration.

Outcome: Public state and declarations parse deterministically, validated declarations resolve into
one desired-state graph, observed installations receive one stable classification, and lock v1 can
migrate safely to lock v2 without changing config v1.

### B) Scope

Create `internal/state`, add only the durable public schemas, replace runtime `yamlmini` use, compile
the graph, classify state, and implement the bounded lock migration. Do not mutate repositories in
this epic.

### C) Files Touched

```text
go.mod, go.sum                                         # modify/create
resources.go                                           # modify schema embedding
internal/state/                                        # create
internal/components/components.go                      # modify adapter
internal/lockfile/lockfile.go                          # modify adapter
internal/manifest/manifest.go                          # modify adapter
components/*/component.yaml                            # validate; accept optional per-file overrides
profiles/standard.yaml                                 # modify generated-surface semantics
schemas/component.schema.json                          # create
schemas/profile.schema.json                            # create
schemas/kit-lock-v1.schema.json                        # create
schemas/kit-lock.schema.json                           # modify v2
schemas/kit-config.schema.json                         # retain v1 structure
tests/fixtures/state/                                  # create
tests/test_json_schemas.py                             # modify
```

### D) Acceptance Criteria And Size

Size: `XL`.

Acceptance criteria:

- maintained YAML and embedded JSON Schema validation handle all durable contracts;
- graph nodes contain owner, presence, update, removal, source, target, route, and digest behavior;
- duplicate targets, invalid paths, missing sources, and incompatible owners fail compilation;
- classifier precedence is deterministic and read-only;
- absent/adoptable state is distinct from partial/invalid state;
- lock v1 and v2 dispatch by exact version;
- known released checksum anomalies receive only the bounded compatibility treatment;
- config v1 optional fields survive round-trip unchanged;
- future state versions remain byte-identical and blocked.

### E) Dependencies And Critical-Path Notes

No epic dependency. This is the shared foundation. New graph fields must describe current paths and
ownership rather than redesign profiles.

### F) Tasks Checklist

- [x] Add pinned `go.yaml.in/yaml/v3 v3.0.4` and `github.com/santhosh-tekuri/jsonschema/v6 v6.0.2` to `go.mod`, resolve `go.sum`, and run `go mod verify`.
- [x] Create `internal/state/schema.go` with embedded schema compilation, local reference resolution, strict instance validation, and public-safe errors.
- [x] Create `internal/state/yaml.go` with YAML-node decoding, JSON-compatible normalization, deterministic encoding, and no filesystem mutation helpers.
- [x] Add component and profile schemas covering current ownership, presence, update, removal, source, target, route, and digest semantics.
- [x] Freeze the current lock contract as `kit-lock-v1.schema.json` and change `kit-lock.schema.json` to v2 with state generation, last successful operation, managed-section digest, and release provenance.
- [x] Add component-schema support for per-file strategy overrides and extend the standard profile with shared state defaults required by the typed desired-state graph.
- [x] Create `internal/state/graph.go` to compile, validate, conflict-check, and sort one immutable graph from selected declarations.
- [x] Create `internal/state/observed.go` to read lock, config, managed files, root managed section, routes, generated surfaces, CLI identity, and active transaction state without writing.
- [x] Implement deterministic primary-state precedence plus diagnostic traits for every approved installation state.
- [x] Create `internal/state/migrate.go` with exact version dispatch and the bounded released-v1 checksum compatibility decoder.
- [x] Preserve valid lock-v1 fields as historical facts and mark placeholder provenance `unverified-legacy` without inventing a digest.
- [x] Refactor existing components, lockfile, and manifest packages into thin compatibility adapters over `internal/state`.
- [x] Add fixtures for numeric-looking strings, zero placeholders, unrelated invalid v1 fields, nested optional config, graph conflicts, state classifications, lock migration, and future versions.
- [x] Run `go test ./internal/state ./internal/components ./internal/lockfile ./internal/manifest` and `python3 -m pytest tests/test_json_schemas.py tests/test_packaging_resources.py -q`.
- [x] Run `python3 scripts/validate-json-schemas.py` and confirm every durable public schema plus representative instance validates.

### G) Implementation Notes

Schemas are public contracts; internal Go types are implementation contracts. No operation-result
or transaction-marker schema is added. The compatibility decoder runs before strict lock-v1
validation and recognizes only the two named paths and three named placeholder forms.

`internal/yamlmini` remains temporarily for retained compatibility tests but is not used by new Go
state paths. Removing it belongs to the Python/legacy decommission gate.

### H) Open Questions

None. Lock-v2/config-v1 and bounded legacy compatibility are settled.

## Epic EP-02 - Transactional Lifecycle Commands

### A) Epic ID, Title, And Outcome

`EP-02` - Transactional Lifecycle Commands.

Outcome: Every lifecycle command uses the same classifier and transactional reconciler, preserves
consumer-owned state, supports concise dry-run output, and leaves no transaction artifacts after
success.

### B) Scope

Create `internal/reconcile`, move init/repair/sync/check/update-check/onboarding onto it, and define
shared Go result/blocker types with golden JSON compatibility tests. Upgrade applies the same engine
in `EP-03`.

### C) Files Touched

```text
internal/reconcile/                                    # create
internal/cli/cli.go                                    # modify
internal/commands/init.go                              # modify
internal/commands/onboard.go                           # modify
internal/commands/repair.go                            # create
internal/commands/sync.go                              # modify
internal/commands/check.go                             # modify
internal/commands/update_check.go                      # modify
internal/commands/commands_test.go                     # modify
internal/drift/drift.go                                # modify graph-backed comparison
tests/test_cli.py                                      # modify
tests/test_init.py                                     # modify
tests/test_onboard.py                                  # modify
tests/test_sync_check.py                               # modify
tests/test_update_check.py                             # modify
tests/test_go_cli_parity.py                            # modify approved differences
```

### D) Acceptance Criteria And Size

Size: `XL`.

Acceptance criteria:

- `init` accepts only absent/adoptable state;
- `repair` restores current-version invariants and can apply the bounded v1 migration;
- `sync` uses only the running binary's content and never changes kit version;
- `update-check` requires a valid installation and changes only update metadata;
- `check` validates schemas, completeness, graph drift, managed sections, routes, version,
  provenance, and recovery state without writing;
- mutating commands support `--dry-run`;
- only upgrade requires an additional `--yes` approval flag;
- change plans are deterministic and consumer-owned paths are never replacement targets;
- the active marker is acquired atomically and blocks concurrent work;
- no-follow containment prevents target escape and replacement races;
- staged state validates before commit and lock v2 commits last;
- rollback restores prior bytes; unrecoverable rollback retains minimal repair state;
- success deletes the marker, stage, backup, and plan;
- text and JSON output use shared Go types and concise golden-tested shapes.

### E) Dependencies And Critical-Path Notes

Depends on `EP-01`. Filesystem safety is the highest-risk path. Real Windows reparse and executable
behavior remains a completion gate in `EP-05`.

### F) Tasks Checklist

- [x] Create `internal/reconcile/result.go` with stable Go types for command, status, state, changes, blockers, validation, provenance, and rollback plus text/JSON renderers.
- [x] Create `internal/reconcile/plan.go` with deterministic create, replace, managed-section merge, preserve, report, migrate, and safe-remove actions.
- [x] Create `internal/reconcile/transaction.go` with canonical-root resolution, atomic create-new marker ownership, target-filesystem staging, bounded backup, staged validation, lock-last commit, post-check, rollback, and success cleanup.
- [x] Implement no-follow traversal, symbolic-link refusal, Windows reparse-point refusal, parent identity revalidation, and stale-marker takeover only after process-liveness plus transaction-identity checks.
- [x] Remove a retired managed path only when prior lock ownership exists and the observed digest still matches the prior managed digest.
- [x] Preserve modified retired paths and return a concise `managed_path_modified` blocker with the exact repair action.
- [x] Refactor `init` to reject existing/partial/invalid state and create lock v2 through the reconciler.
- [x] Add `repair` to restore current-version state, preserve config and consumer content, and expose bounded lock migration in dry-run.
- [x] Refactor `sync` to reconcile current embedded content through the same plan and transaction path.
- [x] Refactor `update-check` to refuse absent/invalid state, preserve managed files, and remain silent when current in agent-notification mode.
- [x] Refactor `check` to execute complete read-only classification and return only concise health facts plus actionable blockers.
- [x] Refactor onboarding so approved absent/adoptable setup calls init and approved existing-kit setup calls repair.
- [x] Keep existing JSON compatibility fields as projections and add golden fixtures for the concise shared result shape.
- [x] Add lifecycle fixtures covering every starting state, no-op idempotency, consumer config, repository instructions, plans, memory, local-user state, safe removal, conflict, concurrent command, interrupted commit, failed validation, successful rollback, failed rollback, and recovery retry.
- [x] Run `go test -race ./internal/reconcile ./internal/commands ./internal/cli ./internal/drift`.
- [x] Run `python3 -m pytest tests/test_cli.py tests/test_init.py tests/test_onboard.py tests/test_sync_check.py tests/test_update_check.py tests/test_go_cli_parity.py -q`.

### G) Implementation Notes

The transaction marker is a strict versioned Go struct written as JSON. It contains only
transaction ID, command, phase, process identity, timestamps, target identity, and recovery
location. Successful operations keep only the short last-operation summary in lock v2.

Invocation of init, repair, and sync is the operation authorization. Onboarding retains its current
user write confirmation. Healthy `check` output should be one compact summary in text mode; JSON
remains available for agents and automation.

### H) Open Questions

None. The user explicitly requested the minimum reliable foundation and minimal approval ceremony.

## Epic EP-03 - Reproducible Release And Safe Upgrade

### A) Epic ID, Title, And Outcome

`EP-03` - Reproducible Release And Safe Upgrade.

Outcome: Embedded identity is non-circular, release packs are byte-reproducible, installers verify
the entire pack, and approved upgrade stages the verified new binary, applies its reconciler, and
rolls back on failure.

### B) Scope

Create `internal/release`, split content identity from archive catalog, harden the existing release
builder and installers, add upgrade, and validate both platforms. Do not publish, sign, tag, or
modify named consumers.

### C) Files Touched

```text
manifest.yaml                                          # modify embedded content identity
resources.go                                           # modify catalog exclusion
schemas/content-manifest.schema.json                  # create
schemas/release-catalog.schema.json                   # create
schemas/pack-manifest.schema.json                     # create
schemas/release-manifest.schema.json                  # retain compatibility fixture support
internal/release/                                     # create
internal/cli/cli.go                                   # modify upgrade command
internal/commands/upgrade.go                          # create
install.sh                                             # modify
install.ps1                                           # modify
bootstrap.md                                          # modify
scripts/build-release-assets.py                       # modify
scripts/validate-release-manifest.py                  # modify
tests/fixtures/release/                               # create
tests/fixtures/installer/                             # modify
tests/test_release_assets.py                          # modify
tests/test_install_metadata.py                        # modify
.github/workflows/validate.yml                        # modify staged platform lanes
```

### D) Acceptance Criteria And Size

Size: `XL`.

Acceptance criteria:

- embedded content identity contains version, graph/component/profile digests, compatibility, and
  consumer impact without enclosing archive URL or digest;
- external catalog is generated after packs and is never embedded or committed as pre-publication
  authority;
- pack manifest declares platform, version, command, binary path/digest, content digest, and
  payload checksum identity;
- ZIP order, timestamps, modes, compression, and Go build IDs are deterministic;
- two identical builds produce byte-identical macOS universal and Windows x64 packs;
- installers validate catalog, outer digest, pack manifest, payload checksums, platform, version,
  command, binary digest, and staged binary version;
- `upgrade --dry-run` verifies candidate evidence in temporary storage and leaves target state
  untouched;
- `upgrade --yes` is required to change version;
- the staged new binary revalidates its handoff, replaces the installed binary, reconciles state,
  runs full check, records provenance, and restores the prior binary/state after failure;
- malformed catalog/archive, missing binary, wrong platform/version, failed smoke, failed
  reconciliation, and failed post-check preserve the prior installation.

### E) Dependencies And Critical-Path Notes

Depends on `EP-01` and `EP-02`. The release builder remains one L2 workflow script. No new script
framework, registry, or service is created.

### F) Tasks Checklist

- [x] Change `manifest.yaml` to the content-manifest contract and remove enclosing archive URLs plus digests from embedded identity.
- [x] Create content-manifest, release-catalog, and pack-manifest schemas with exact compatibility, identity, platform, command, path, and digest requirements.
- [x] Change `resources.go` so content identity and schemas are embedded while external catalogs remain release artifacts.
- [x] Create `internal/release/catalog.go` and `pack.go` to validate catalog selection, safe extraction, payload checksums, pack identity, and binary identity.
- [x] Create `internal/release/upgrade.go` and `handoff.go` for version-direction checks, temporary dry-run verification, approved target staging, hash-bound handoff, new-binary reconciliation, final check, and rollback.
- [x] Add public `upgrade --version <version> --dry-run` and `upgrade --version <version> --yes` adapters while keeping internal handoff mode absent from help.
- [x] Refactor the release builder into deterministic source validation, binary build, payload assembly, normalized archive, pack verification, repeat-build comparison, and catalog emission phases.
- [x] Normalize ZIP entry order, timestamps, modes, creator flags, compression parameters, and Go build IDs.
- [x] Keep release-builder structured output inside the existing script and tests without adding a script registry and without adding a separate wrapper.
- [x] Refactor `install.sh` to verify the full official catalog-to-binary chain, stage replacement, smoke-test, and restore the prior binary after failure.
- [x] Refactor `install.ps1` to perform the equivalent Windows x64 verification, staging, smoke, replacement, and restoration path.
- [x] Update `bootstrap.md` with verified fresh install, explicit upgrade approval, dry-run, rollback, and unsigned internal/prototype wording.
- [x] Add fixtures for invalid catalog, checksum mismatch, traversal, malformed archive, missing binary, wrong platform, wrong version, wrong command, failed smoke, failed reconciliation, failed post-check, handoff tampering, and preserved prior install.
- [x] Add build-twice tests comparing final bytes, digests, entry metadata, pack manifests, payload checksums, and candidate catalogs.
- [x] Run `go test ./internal/release ./internal/reconcile ./internal/commands`.
- [x] Run `python3 -m pytest tests/test_release_assets.py tests/test_install_metadata.py tests/test_packaging_resources.py -q`.
- [x] Run staged macOS install and upgrade success plus failure paths in an isolated temporary home while leaving both the developer CLI and source-repository `.codeheart/` state unchanged.
- [ ] Run the Windows install and upgrade matrix through the validation workflow and retain CI results as the completion evidence.

### G) Implementation Notes

Dry-run may use bounded operating-system temporary storage for network verification, but it must
clean up and leave the target repository, transaction path, and installed binary unchanged. Apply
mode acquires the target transaction marker before durable staging.

The staged binary controls executable replacement after the parent exits, which avoids Windows
self-locking. The external catalog uses HTTPS plus SHA-256 under the approved unsigned prototype
boundary. Signing remains separate.

### H) Open Questions

`OQ-1` and `OQ-3` remain non-blocking and outside source implementation.

## Epic EP-04 - Compact Lifecycle Guidance And Producer Routing

### A) Epic ID, Title, And Outcome

`EP-04` - Compact Lifecycle Guidance And Producer Routing.

Outcome: Maintainers use canonical source, consumers can select the correct lifecycle command from
one compact runbook, and two focused probes prove the changed routing without expanding the broader
guidance system.

### B) Scope

Add producer instructions, one managed lifecycle runbook, compact consumer routes, onboarding and
update-policy alignment, retained resource mirrors, and two low-context probes.

Affected runbooks:

| Runbook | Audience and treatment |
| --- | --- |
| `conduct-first-run-onboarding.md` | Hybrid: preserve visible question pacing; separate agent classification; route setup to init and existing kit to repair. |
| `maintain-operating-kit-installation.md` | Agent-facing, new: one short command/precondition/approval/blocker/recovery reference. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer: add the smallest state, transaction, release, and platform validation gates. |
| `docs/repo/runbooks/release-operating-kit.md` | Maintainer: add reproducibility, catalog, pack, upgrade, unsigned-boundary, and publication stops. |

### C) Files Touched

```text
AGENTS.md                                              # modify repository-owned producer section
README.md                                              # modify concise architecture summary
components/agent-interface/component.yaml             # modify inventory
components/agent-interface/managed/kit-readme.md      # modify command routes
components/agent-interface/managed/reference/
  update-check-policy.md                              # modify
components/agent-interface/managed/runbooks/
  conduct-first-run-onboarding.md                     # modify
  maintain-operating-kit-installation.md              # create
templates/agents/AGENTS.managed-block.md              # modify
src/codeheart_operating_kit/resources/                # mirror changed resources
docs/repo/runbooks/change-operating-kit.md            # modify
docs/repo/runbooks/release-operating-kit.md           # modify
tests/test_onboard.py                                 # modify
tests/test_packaging_resources.py                     # modify
tests/test_routing.py                                 # create
```

### D) Acceptance Criteria And Size

Size: `M`.

Acceptance criteria:

- producer instructions name canonical source and mark local installed materialization
  non-authoritative;
- consumer routes distinguish init, repair, sync, update-check, upgrade, and check;
- the lifecycle runbook fits one command-selection and recovery workflow without duplicating
  implementation architecture or generic tooling setup;
- onboarding preserves user dialogue and routes existing installations to repair;
- update-check remains valid-install-only and silent when current;
- source and retained resource mirrors match;
- one fresh producer probe and one materialized-consumer probe identify owner, route, scope, and
  approval before tool selection.

### E) Dependencies And Critical-Path Notes

Depends on `EP-02` and `EP-03`. This epic documents implemented behavior only. Broader guidance
quality, route scoring, and full runbook evaluation remain outside scope.

### F) Tasks Checklist

- [x] Add a repository-owned producer section to root `AGENTS.md` naming canonical source authorities and classifying local `.codeheart/` materialization as non-authoritative.
- [x] Create `maintain-operating-kit-installation.md` with a compact intention, state-to-command table, dry-run examples, upgrade approval, blocker table, recovery path, tooling-readiness route, and stop conditions.
- [x] Update first-run onboarding with separate user-dialogue and agent-execution paths that select init, repair, diagnosis, then stop for incompatible state.
- [x] Update the managed kit readme, update-check policy, and root managed template with the six lifecycle command routes and minimal approval wording.
- [x] Update the Agent Interface component manifest with the new runbook and required state semantics.
- [x] Update root `README.md` with a concise state, transaction, provenance, platform, and migration overview.
- [x] Update change and release runbooks with the smallest required schema, migration, transaction, reproducibility, catalog, installer, upgrade, platform, and publication gates.
- [x] Mirror changed managed files, templates, schemas, profiles, and content identity into the retained Python resource tree.
- [x] Add adoption tests proving one managed root section plus preserved repository-owned and local-user instructions.
- [x] Add two focused routing test scenarios covering producer source authority and consumer repair selection.
- [x] Run one fresh low-context producer probe and one fresh materialized-consumer probe, then record only prompt, selected route, approval class, and pass/fail in the execution log.
- [x] Run `python3 -m pytest tests/test_onboard.py tests/test_packaging_resources.py tests/test_routing.py -q`.
- [x] Run `python3 scripts/validate-markdown-headers.py` and `python3 scripts/validate-public-core.py`.

### G) Implementation Notes

The lifecycle runbook is the only new durable guidance artifact. It calls the L3 CLI and does not
restate schemas, transaction phases, package-manager setup, or release-builder internals.

Probe evidence stays in the execution log. No probe attachment, route registry expansion, or broad
guidance audit is created.

### H) Open Questions

None. General operating-guidance effectiveness is explicitly deferred.

## Epic EP-05 - Integrated Validation And Source Handoff

### A) Epic ID, Title, And Outcome

`EP-05` - Integrated Validation And Source Handoff.

Outcome: Automated local and real-platform evidence proves state safety, lifecycle behavior,
release integrity, upgrade rollback, and routing; one final review confirms capability coverage;
the work stops before publication and consumer rollout.

### B) Scope

Run the complete regression and failure matrices, update unreleased migration notes, perform one
final implementation review, update planning lifecycle records, and record negative evidence. Do
not require per-epic manual review documents or separate evidence attachments.

### C) Files Touched

```text
.github/workflows/validate.yml                         # modify final automated matrix
release-notes.md                                       # modify unreleased notes
docs/README.md                                         # update final discoverability
docs/repo/README.md                                    # update final status
docs/repo/plans/README.md                              # update final status
docs/repo/plans/plan-register.md                       # update lifecycle/evidence
docs/repo/plans/operating-kit-state-release-architecture/
  operating-kit-state-release-architecture_implementation_doc.md # update lifecycle
  operating-kit-state-release-architecture_execution_log.md      # create on activation
```

### D) Acceptance Criteria And Size

Size: `L`.

Acceptance criteria:

- all starting-state and lifecycle-command combinations have automated expected results;
- consumer config, instructions, plans, memory, local-user content, and modified retired managed
  files remain unchanged across applicable operations;
- failure injection proves rollback, recovery, and idempotent retry;
- macOS and Windows containment tests cover concurrency, links/reparse points, parent replacement,
  and stale markers;
- schema, migration, graph, lifecycle, release, installer, upgrade, routing, public-core, and
  Markdown validation passes;
- release candidates are byte-reproducible;
- real Windows install/upgrade success and failure paths pass;
- the final reviewer finds no material capability or safety gap;
- execution log contains concise command/CI evidence and no copied raw operational logs;
- negative evidence confirms no release, tag, named consumer sync, local update, signing change,
  Python deletion, or new platform work.

### E) Dependencies And Critical-Path Notes

Depends on `EP-01` through `EP-04`. Real Windows validation is a completion gate. Inability to run
it does not authorize repository-setting changes or publication.

### F) Tasks Checklist

- [x] Extend `.github/workflows/validate.yml` with state migration, transaction failure, reproducibility, consumer materialization, routing, macOS install/upgrade, and Windows install/upgrade jobs.
- [x] Run the automated state-transition matrix and confirm every command precondition, ending state, blocker, and preservation digest.
- [ ] Run transaction failure, concurrent-operation, symbolic-link, Windows reparse-point, parent-replacement, stale-marker, rollback, and recovery-retry fixtures. Local cases passed; Windows reparse execution remains in the pending Windows workflow.
- [x] Build macOS universal and Windows x64 release candidates twice and confirm byte-identical packs plus coherent catalog, pack, content, and binary digests.
- [x] Run isolated macOS fresh-install and upgrade success plus failure paths.
- [ ] Run the validation-only Windows workflow and record run IDs, commit SHA, and pass/fail summary without publishing assets.
- [x] Run `go test -race ./...`.
- [x] Run `python3 -m pytest -q` and record approved Go-only compatibility differences.
- [x] Run public-core, Markdown, JSON Schema, and release-contract validators.
- [x] Update `release-notes.md` with unreleased lock migration, command, validation, rollback, provenance, and unsigned-boundary notes without selecting a version.
- [x] Classify impact as consumer migration required, validator-only change, instruction-only change, and security/safety policy change without a placement break.
- [x] Run one final implementation review covering capability, consumer preservation, transaction safety, release integrity, platform symmetry, routing, and publication stop.
- [x] Resolve all critical and high-severity review findings and rerun the affected automated validation.
- [x] Update documentation indexes and the plan register with final source status plus execution-log and CI references.
- [x] Record explicit negative evidence for release, tag, consumer sync, local kit update, signing, Python removal, and new platform actions.
- [x] Re-run the full automated validation set and record one concise final evidence summary in the execution log.

### G) Implementation Notes

CI artifacts and command outputs are the detailed evidence. The execution log records commands,
results, run IDs, and residual risk; it does not duplicate test matrices or retain raw logs.

There are no per-epic manual review files. The final review occurs after integrated validation and
before source handoff. Publication and named consumer sync require separate explicit approval.

### H) Open Questions

`OQ-1`, `OQ-2`, and `OQ-3` remain non-blocking deferred work.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Define signing, notarization, catalog-signature, and attestation policy before broad-public
  distribution. Trigger: user approval of a trust target.
- Retire the Python CLI, resource mirror, legacy parser, and parity suite. Trigger: Go fixtures are
  authoritative and named migration evidence plus rollback window are approved.
- Evaluate general routing accuracy, runbook executability, guidance ambiguity, tool selection, and
  operational outcomes. Trigger: the user requests improvement of guidance effectiveness beyond
  the kit's self-management path.
- Modularize profiles or redesign documentation indexes. Trigger: a concrete second profile or
  information-architecture problem is approved.
- Add Linux, Windows ARM, package-manager channels, and enterprise catalogs. Trigger: platform
  demand and CI support are approved.
- Select a release version, update final catalog URLs, tag, publish, and sync named consumers.
  Trigger: source validation passes and the user separately approves release execution.
- Apply the available local Operating Kit update. Trigger: separate user approval.

## 4.2 Future Considerations

- Publish an operation-result schema only when an external integration needs a stable independent
  contract.
- Split an internal package only after measured growth, independent ownership, or a second caller
  makes the boundary useful.
- Add persistent transaction audit history only after a real operational or compliance requirement.
- Extend the catalog with detached signatures without changing consumer state paths.
- Monitor real consumer rollout for Windows file locking and antivirus interference before adding
  a platform helper.
- Treat general guidance effectiveness as a product-quality workstream with outcome benchmarks,
  not as more lifecycle ceremony.

# Revision Notes

- 2026-07-09: Created the initial eight-epic implementation draft from the approved discovery.
- 2026-07-09: Resolved the initial planning review findings for bounded lock-v1 migration,
  filesystem containment, concurrent operation safety, target-write-free upgrade dry-run,
  contracts ownership, generated catalog placement, and immediate discoverability; re-review
  returned `Ready`.
- 2026-07-09: Simplified the plan at the user's request while preserving approved capability.
  Consolidated seven proposed internal domains into three packages, reduced eight epics to five,
  removed schemas for transient transaction/results, removed persistent successful journals,
  removed the script-role index and evidence attachments, consolidated lifecycle guidance into one
  runbook, limited routing proof to two focused probes, replaced per-epic manual evidence with
  automated validation, and kept additional approval only for upgrade.
- 2026-07-09: Reviewed the simplified plan against `FR-001` through `FR-012`, `NFR-001` through
  `NFR-010`, the five approved capability scopes, the implementation-plan quality gate, routing
  coverage, and recipe-maturity coverage. No approved capability was lost and no planning blocker
  remains. Real Windows validation remains an execution completion gate.
- 2026-07-09: Activated for source implementation after explicit user approval. Release
  publication, tag creation, named consumer sync, local kit update, signing, Python retirement, and
  new platforms remain outside execution authority.
