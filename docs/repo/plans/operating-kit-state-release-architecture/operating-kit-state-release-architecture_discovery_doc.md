Last updated: 2026-07-09T20:40:04Z (UTC)
Created: 2026-07-09
Status: completed
Completed: 2026-07-09

# Operating Kit State And Release Architecture Discovery

## Discovery Status

Input state: `review-or-cleanup` followed by a user-requested focused architecture discovery.

Output target: `implementation-handoff-ready`. The user approved `D-001` through `D-006` and the
recommended command, migration, and delivery defaults on 2026-07-09. This document now authorizes
implementation planning for the capability scopes below. It does not authorize source execution,
release publication, consumer sync, migration execution, cleanup, or deletion.

This discovery is:

- `routing-bearing` because the likely implementation may change generated root routing,
  managed-section ownership, source-repository routing, and route validation;
- `recipe-bearing` because initialization, repair, synchronization, upgrade, checking, installer,
  packaging, and release mechanics are repeatable operational recipes with failure, evidence, and
  rollback requirements.

The applicable managed references are:

- `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`;
- `components/agent-interface/managed/reference/operational-recipe-maturity.md`;
- `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`;
- `components/agent-interface/managed/reference/runbook-authoring-standard.md`.

## Problem Framing

### Problem Statement

The Operating Kit has a strong managed-versus-consumer ownership model and a well-tested CLI, but
several high-authority contracts are not enforced by the same runtime model:

- lifecycle commands overlap and do not fail closed on incompatible repository states;
- profiles, component manifests, schemas, release manifests, lock state, and runtime hardcoding
  each describe parts of installation behavior;
- schema files exist, but consumer `check` does not validate complete lock and config instances;
- synchronization writes directly and cannot preview, stage, roll back, or reconcile removed
  managed paths safely;
- release metadata embeds the checksum of the archive that contains the binary carrying that
  metadata, creating a circular provenance boundary;
- the Operating Kit source repository routes maintainers through an older installed consumer
  layer instead of an explicit producer mode.

The result is not general instability: the current test and validator suite passes. The problem is
that important invalid, partial, stale, or contradictory states are outside the suite's present
contract.

### User Intention

Create a coherent architecture for Operating Kit state, lifecycle commands, updates, release
provenance, and source-repository self-hosting before drafting implementation work.

### Goals

- Preserve consumer-owned content and valid consumer configuration across every repair, sync, and
  upgrade path.
- Give every lifecycle command a clear state transition, precondition, mutation boundary, and
  failure result.
- Establish one typed desired-state model that runtime code and generated release metadata share.
- Make declared schemas executable for source manifests and installed consumer state.
- Make managed-state mutation previewable, transactional, recoverable, and auditable.
- Remove the release archive checksum self-reference and make installed provenance trustworthy.
- Give this source repository an explicit producer-mode route that cannot silently use stale
  installed doctrine.
- Preserve the existing ownership model, public-core boundary, and current consumer paths unless
  a migration is explicitly justified.

### Non-Goals

- Replacing the Go CLI or changing its public name.
- Redesigning the full planning doctrine, plan register, or documentation hierarchy in this work.
- Introducing new domain modules or provider integrations.
- Adding Linux or Windows ARM support.
- Requiring signed or notarized broad-public releases in the first implementation, although the
  design must leave a clean trust boundary for them.
- Removing consumer-owned plans, memory, configuration, or repository instructions.
- Removing the legacy Python oracle before a separately approved decommission gate.
- Publishing a release, applying the locally available kit update, or syncing named consumers.

### Priority Order

1. Consumer-state safety and recoverability.
2. Contract integrity and trustworthy provenance.
3. Backward compatibility and migration clarity.
4. Maintainability and one source of truth.
5. Operator and agent ergonomics.
6. Implementation speed.

This order is inferred from the repository's existing safety rules and the user's request for
fundamental improvement. The user may correct it during manual review.

### Durable Principles

- Managed content may be reconciled by the kit; consumer-owned content must not be silently
  rewritten.
- A command name must describe one lifecycle responsibility.
- Desired state, observed state, and historical evidence are different data classes.
- The same declared contract should drive mutation, validation, reporting, and release snapshots.
- A failed operation must leave either the previous valid state or a machine-readable recoverable
  state.
- Archive integrity metadata must not depend circularly on the archive that contains it.
- Producer repositories and consumer repositories may share doctrine while using different state
  authority paths.
- Security-sensitive and state-changing behavior remains approval-gated.

### Manual-Review Success Criteria

This discovery is manual-review-ready when:

- the current failure modes and contract gaps are recorded as evidence;
- lifecycle, desired-state, serialization, provenance, self-hosting, and scope decisions each
  have concrete recommendations;
- architecture reviewer objections are incorporated or bounded;
- the exact approval needed before implementation planning is visible.

## Candidate Domains And Decision Clusters

| Domain | Why it is involved | Authority source | Decision cluster |
| --- | --- | --- | --- |
| Consumer lifecycle | Commands create, repair, synchronize, validate, and update installed state. | CLI behavior, lock/config schemas, ownership contract | Lifecycle state machine and transactional mutation |
| Desired-state and contracts | Manifests, profiles, schemas, lock records, and runtime code describe overlapping surfaces. | Component/profile manifests and schemas | Typed installation graph and validation |
| Release engineering | Binaries, installers, archives, manifests, checksums, and tags establish installed provenance. | Release runbook, build script, root release manifest | External provenance and reproducibility |
| Agent routing | Root `AGENTS.md` and managed docs choose the operating authority. | Agent Interface and source-repository instructions | Producer-mode routing and managed-section verification |
| Migration governance | Existing consumers and the Python oracle constrain safe rollout. | Consumer-impact classification and completed bootstrap plan | First implementation boundary and deferred migration |

## Context And Evidence

### Essential Sources

| Source | Current fact | Discovery impact |
| --- | --- | --- |
| `internal/commands/onboard.go` | Approved onboarding calls `Initialize` for every non-ambiguous inspected state, including existing-kit repair. | Repair behavior currently re-enters initialization. |
| `internal/components/components.go` | Initialization writes managed files, scaffolds, root instructions, ignore rules, a new lock, and a new config without merging existing kit state. | Existing kit configuration and generated state can be replaced. |
| `internal/commands/sync.go` | Sync copies files and changes root/ignore state before writing the refreshed lock. | Failure can leave a partially updated installation. |
| `internal/commands/check.go` | Check evaluates lock metadata presence, managed-file digests, root route existence, and CLI version, but not the complete config or declared JSON Schemas. | A partial or schema-invalid install can report healthy. |
| `profiles/standard.yaml` and `components/*/component.yaml` | Declarative ownership and surface metadata exists. | These are the natural inputs to a compiled desired-state graph. |
| `internal/components/components.go` and `internal/capabilities/capabilities.go` | Generated surfaces, ignore behavior, default config, and capability IDs are also hardcoded. | Declarations are not yet runtime-authoritative. |
| `schemas/kit-lock.schema.json` | Repair and release checksums are typed as 64-character strings. | Runtime must preserve scalar type and validate instances. |
| `internal/yamlmini/yamlmini.go` | Digit-only scalars are parsed as integers, while digit-only strings are not always quoted when emitted. | All-zero checksum placeholders do not round-trip as schema-valid strings. |
| `.codeheart/kit.lock.yaml` | The current local lock contains `repair_checksum_sha256: 0`. | A concrete schema-invalid value is not detected by `check`. |
| `resources.go`, `manifest.yaml`, and `scripts/build-release-assets.py` | The binary embeds the release manifest, while the release manifest records the digest of the archive containing that binary. | Final archive provenance cannot be self-contained without circularity. |
| `scripts/build-release-assets.py` | ZIP entries inherit source timestamps and metadata. | Consecutive builds need not be byte-reproducible. |
| `install.sh` and `install.ps1` | Installers verify the outer archive checksum and smoke-test a located binary; pack manifest and internal checksum metadata are not authoritative inputs. | The pack contract and installed provenance are weaker than the emitted metadata suggests. |
| `AGENTS.md`, `.codeheart/kit.lock.yaml`, and `manifest.yaml` | The source repository routes into an untracked v0.1.17 installed layer while source is v0.1.21. | Maintainers can follow stale installed doctrine. |
| `src/codeheart_operating_kit/` and `tests/test_go_cli_parity.py` | A complete Python oracle remains and approved parity differences already exist. | The oracle is migration debt, not the future source of behavior authority. |

### Validation Baseline

The architecture review established this baseline on 2026-07-09:

- `go test ./...`: passed;
- full Python suite: `124 passed`;
- public-core validator: passed;
- Markdown timestamp validator: passed;
- JSON-schema-structure validator: passed;
- release-manifest validator: passed;
- source and packaged Python resource mirrors: matched;
- current-source `check` on this repository: no managed-file drift or missing route targets, but
  `stale_cli: true` because the local installed state is older than source.

This baseline proves the repository is testable. It also proves the new work needs state-transition,
instance-schema, provenance, reproducibility, and semantic-routing coverage rather than only more
of the existing checks.

## Confirmed Failure And Gap Inventory

- `F-001`: existing-kit onboarding can rewrite generated kit configuration while presenting the
  operation as repair.
- `F-002`: `init`, `sync`, and `update-check` do not consistently reject absent, partial, or
  incompatible state.
- `F-003`: sync has no dry-run, staged reconciliation, rollback journal, or safe removal rule for
  managed paths retired by the desired state.
- `F-004`: `check` does not validate `.codeheart/kit.config.yaml` or full schema conformance.
- `F-005`: the current checksum placeholder can serialize as an integer even though the schema
  requires a string.
- `F-006`: component/profile declarations and runtime hardcoding can diverge.
- `F-007`: the binary/archive release manifest relationship is circular.
- `F-008`: release packs are not proven byte-reproducible by CI.
- `F-009`: source-repository routing can use stale consumer-installed doctrine.
- `F-010`: Python parity can preserve obsolete behavior or require exceptions after Go becomes the
  behavioral authority.

## Requirements And Constraints

### Functional Requirements

- `FR-001`: `init` must create a new installation only when the target is absent or explicitly
  eligible for adoption.
- `FR-002`: `repair` must preserve valid consumer configuration and consumer-owned content while
  restoring managed and generated invariants.
- `FR-003`: `sync` must reconcile managed content from the currently running trusted kit without
  pretending to install a newer kit version.
- `FR-004`: `upgrade` must be the explicit approval-gated path that resolves trusted release
  metadata, stages and verifies a new binary, invokes its reconciliation behavior, validates the
  result, and rolls back on failure.
- `FR-005`: `check` must validate installation completeness, schemas, managed digests, managed root
  sections, route targets, CLI/source version coherence, and recoverable transaction state.
- `FR-006`: every mutating lifecycle command must expose a dry-run or change-plan mode before
  applying changes.
- `FR-007`: one resolved desired-state graph must describe managed files, managed sections,
  generated state, scaffolds, local-user/local-machine surfaces, and conditional reports.
- `FR-008`: lock and config writes must validate against their schemas before replacement.
- `FR-009`: synchronization must remove a retired managed path only when the prior lock proves kit
  ownership and the current file still matches the prior managed digest; otherwise it must stop or
  report a conflict.
- `FR-010`: archive digest and release provenance must be supplied by a release catalog outside the
  binary/archive self-reference.
- `FR-011`: pack creation must be reproducible for a pinned source and toolchain, and installers
  must validate the declared platform, version, binary path, and resulting binary version.
- `FR-012`: the source repository must route maintainers to current canonical source or a
  consistency-checked current materialization.

### Non-Functional Requirements

- `NFR-001` Safety: no lifecycle command silently overwrites consumer-owned state or valid optional
  config.
- `NFR-002` Atomicity: a failed mutation preserves the last valid state or a bounded recoverable
  transaction record.
- `NFR-003` Idempotency: repeating a successful repair, sync, or check produces no additional
  changes.
- `NFR-004` Compatibility: existing G1 paths and YAML files remain readable through the first
  implementation unless an explicit migration is approved.
- `NFR-005` Determinism: desired-state resolution, serialization, lock ordering, archive contents,
  and validation output are stable.
- `NFR-006` Auditability: machine-readable results identify state before, intended changes,
  applied changes, conflicts, validation, provenance, and rollback outcome without exposing
  secrets.
- `NFR-007` Portability: macOS and Windows receive equivalent lifecycle, installer, failure-path,
  and rollback coverage.
- `NFR-008` Maintainability: version, surface, capability, route, and release metadata are generated
  or validated from authoritative inputs instead of synchronized manually.
- `NFR-009` Public-core safety: logs, fixtures, and release evidence remain public-safe.
- `NFR-010` Context economy: the implementation should add compact machine contracts rather than
  duplicating long prose across routers and runbooks.

### Consumer Impact Boundary

The likely implementation will require at least these impact classifications:

- `validator-only change` for executable schema and semantic validation;
- `consumer migration required` if lock/config representation or update mechanics require adoption;
- `security or safety policy change` for trust, rollback, and mutation gates;
- `instruction-only change` for managed routing and runbook alignment.

A `breaking placement-contract change` is not recommended. Preserve current paths unless later
evidence proves a path migration necessary.

## Decision Inventory

| ID | Decision | Class | State | Depends on | Affected IDs | Blocks | Closure criteria |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `D-001` | Lifecycle command state model | blocking | approved | None | `FR-001` through `FR-006`, `NFR-001` through `NFR-003`, `R-001`, `R-002` | `D-002`, `D-004`, `D-006` | Closed by user approval on 2026-07-09. |
| `D-002` | Desired-state authority | implementation-shaping | approved | `D-001` | `FR-007`, `FR-009`, `FR-012`, `NFR-005`, `NFR-008`, `R-003`, `R-005`, `R-008` | `D-003`, `D-005`, `D-006` | Closed by user approval on 2026-07-09. |
| `D-003` | Machine-state format and validation | implementation-shaping | approved | `D-002` | `FR-005`, `FR-008`, `NFR-004`, `NFR-005`, `R-004` | `D-006` | Closed with lock schema v2, config schema v1, mature YAML parsing, and versioned migration approved on 2026-07-09. |
| `D-004` | Release provenance and reproducibility | blocking | approved | `D-001`, `D-002` | `FR-004`, `FR-010`, `FR-011`, `NFR-002`, `NFR-005` through `NFR-007`, `R-002`, `R-006` | `D-006` | Closed by user approval on 2026-07-09. |
| `D-005` | Source-repository producer mode | implementation-shaping | approved | `D-002` | `FR-012`, `NFR-008`, `NFR-010`, `R-008` | `D-006` | Closed by user approval on 2026-07-09. |
| `D-006` | First implementation boundary | blocking | approved | `D-001` through `D-005` | All requirements and risks | Implementation handoff | Closed with one source implementation plan and all recorded deferrals approved on 2026-07-09. |

Owners: the user owns approval or correction of every decision. The Operating Kit source
repository owns future implementation after approval.

## Decision Ledger

### D-001 - Lifecycle Command State Model

Exact question: Should the current overlapping command behavior remain, collapse into one generic
reconcile command, or become explicit lifecycle commands with state-specific preconditions?

Options:

1. Preserve current commands and add isolated guards.
2. Replace them with one general `reconcile` command.
3. Use explicit `init`, `repair`, `sync`, `upgrade`, and `check` responsibilities backed by one
   shared state engine.

Recommendation: option 3.

Rationale:

- It makes user intent and approval class visible at the command boundary.
- It prevents an operation presented as repair from performing fresh initialization.
- It allows shared transactional mechanics without hiding materially different authority levels.
- It gives future agents and runbooks stable route cards: create, repair, reconcile current
  content, install a new version, and validate.

Rejected alternatives:

- Guard-only changes leave overlapping semantics and make future failures command-specific.
- One reconcile command obscures whether a new installation, local repair, or external release
  update is authorized.

Required behavior:

- `init`: absent/adoption state only; stop on an existing or partial installation unless the caller
  explicitly chooses the appropriate route.
- `repair`: preserve and validate consumer configuration; restore current-version invariants.
- `sync`: current trusted binary only; no network version change.
- `upgrade`: explicit external-state-changing release transition with approval, staging, rollback,
  and post-upgrade check.
- `update-check`: valid installed state only; update update-check metadata without creating or
  repairing an installation.
- `check`: read-only full-state diagnosis.

Risks: adding commands increases surface area. Mitigation: all commands share one typed state
classifier and transaction engine rather than separate implementations.

State: `approved` by the user on 2026-07-09.

### D-002 - Desired-State Authority

Exact question: What should be the authoritative runtime representation of installed artifacts and
ownership behavior?

Options:

1. Keep runtime hardcoding and use tests to compare it with manifests.
2. Treat raw component/profile YAML as runtime logic everywhere.
3. Compile validated component and profile declarations into one typed desired-state graph used by
   lifecycle commands, check, lock generation, routing generation, and release snapshots.

Recommendation: option 3.

The graph should distinguish:

- `content_owner`: kit, consumer, local user, or local machine;
- `presence_policy`: required, create-when-absent, optional, conditional, or historical-only;
- `update_strategy`: replace, managed-section merge, preserve, append-only, or report-only;
- `removal_strategy`: reconcile, preserve, or explicit migration;
- `source`, `target`, component, profile, condition, and expected digest;
- route identity and validation requirements where the artifact is routing-bearing.

Rationale:

- The existing component/profile model becomes executable rather than decorative.
- Desired state can be compared with observed state before mutation.
- Lock records can report observed and historical facts without being reused as desired-state
  declarations.
- Root routes, release manifests, and docs indexes can be generated or semantically validated from
  the same owner data.

Rejected alternatives:

- Hardcoding remains a competing source of truth.
- Using untyped raw YAML directly spreads interpretation across commands and weakens validation.

State: `approved` by the user on 2026-07-09.

### D-003 - Machine-State Format And Validation

Exact question: Should the first implementation preserve YAML, switch installed state to JSON, or
retain the custom YAML subset?

Options:

1. Retain the custom YAML parser and patch known scalar cases.
2. Migrate lock/config files to deterministic JSON.
3. Preserve existing YAML paths and syntax, replace the custom parser with a mature YAML library,
   and validate parsed instances with the declared schemas.

Recommendation: option 3 for the first implementation.

Rationale:

- It fixes scalar typing and parsing completeness without a placement or format migration.
- It preserves existing user-facing config files and release compatibility.
- It allows schema validation to become authoritative before any future format discussion.

Validation requirements:

- validate component, profile, release, lock, and config instances;
- validate before write and after read;
- reject unknown or incompatible fields according to the schema version;
- preserve supported optional consumer config during repair and upgrade;
- dispatch reads and migrations by `schema_version`, fail closed on unsupported future versions,
  and never normalize a file through a schema version the running CLI does not understand;
- require an explicit dry-run migration plan before a versioned state rewrite;
- include fixtures for all-zero hashes, numeric-looking strings, timestamps, nested optional config,
  unknown fields, and schema-version migration.

Residual concern: a mature YAML dependency increases binary dependency surface. This is preferable
to maintaining a parser for public machine contracts; dependency version and license validation
belong in implementation planning.

Approved default: introduce lock schema v2 for transaction and provenance changes; retain config
schema v1 unless implementation evidence requires a structural config change. The new CLI must read
lock v1, produce a dry-run migration plan, migrate to lock v2 only through an authorized lifecycle
operation, and never downgrade or normalize unsupported future versions.

State: `approved` by the user on 2026-07-09.

### D-004 - Release Provenance And Reproducibility

Exact question: Where should the authoritative digest and provenance of a release archive live?

Options:

1. Continue embedding the enclosing archive checksum in the binary's release manifest.
2. Remove release provenance from installed state and rely only on a successful binary smoke test.
3. Keep archive digest and publication provenance in an external release catalog; keep only
   non-circular kit/component identity inside the binary; have the installer and upgrade flow
   record the verified external provenance in installed state.

Recommendation: option 3.

Rationale:

- It removes the binary/archive self-reference.
- It lets the release catalog be signed or attested later without changing consumer state layout.
- It separates content identity from distribution-container identity.
- It gives installation and upgrade a precise evidence chain: requested release, catalog, asset,
  outer digest, pack metadata, binary identity, reconciliation result, and final check.

The external catalog is the only authority for the enclosing archive digest. The binary may embed
its version, component identities, content digests, and compatible schema versions, but it must not
embed the digest of its enclosing release archive.

Required first-version behavior:

- deterministic ZIP metadata and pinned build inputs;
- build-twice reproducibility check;
- outer archive checksum verification;
- pack manifest validation for version, platform, binary path, and command identity;
- staged binary version verification;
- macOS and Windows malformed archive, missing binary, wrong platform/version, failed smoke test,
  rollback, and prior-install preservation tests;
- explicit unsigned internal/prototype boundary until signing policy is separately approved.

Rejected alternatives:

- Continuing the embedded checksum cannot produce a stable fixed point.
- Removing provenance would weaken repair and upgrade evidence.

State: `approved` by the user on 2026-07-09.

### D-005 - Source-Repository Producer Mode

Exact question: Should this source repository operate through a normal installed consumer layer or
through canonical source with an explicit self-host consistency test?

Options:

1. Continue using an untracked local installed layer as the maintainer route.
2. Commit the full installed consumer layer into the source repository.
3. Route maintainer work to canonical source components and treat a generated consumer install as
   an ignored test fixture checked for parity in CI.

Recommendation: option 3.

Rationale:

- Maintainers cannot be routed through stale installed doctrine while changing newer source.
- The repository avoids committing generated consumer state and duplicate installed docs.
- Consumer behavior remains dogfooded through an isolated materialization test.
- The tracked root managed section can be generated or checked against current canonical source.

Required behavior:

- producer-specific repository guidance points maintainers to source owners;
- CI materializes a fresh consumer installation and runs full check/route probes;
- CI compares the tracked root managed section with the canonical template or an explicitly
  producer-owned equivalent;
- local `.codeheart/` remains ignored or clearly classified as a non-authoritative fixture;
- adoption normalizes existing instructions into the defined repository-owned layer instead of
  appending an extra instruction hierarchy.

State: `approved` by the user on 2026-07-09.

### D-006 - First Implementation Boundary

Exact question: Which recommendations belong in the first implementation plan, and which should
remain separate follow-ups?

Options:

1. Fix only the observed defects and defer architecture.
2. Include every audit recommendation, including Python retirement, documentation simplification,
   profile redesign, signing, and new platforms.
3. Implement the coherent state-and-release foundation from `D-001` through `D-005`, while keeping
   unrelated simplification and distribution expansion separate.

Recommendation: option 3.

The first implementation plan should cover these capability groups:

- immediate state-preservation and fail-closed safeguards;
- lifecycle state classification, change planning, transaction, rollback, and structured results;
- typed desired-state graph and schema-backed parsing/validation;
- explicit versioned state migration with unsupported-future-version stop behavior;
- safe managed-path and managed-section reconciliation;
- approval-gated upgrade using external provenance;
- reproducible packs and symmetric installer failure-path validation;
- source-repository producer mode and self-host proof;
- documentation and release-runbook alignment required by those behaviors.

Explicitly defer:

- deletion of the Python oracle and packaged Python resource tree;
- broad planning-document and index redesign;
- profile modularization beyond changes required for the typed graph;
- Linux, Windows ARM, signing/notarization automation, and new external providers;
- named consumer sync and public release publication.

Rationale: option 3 fixes the unsafe and circular foundations without turning the plan into a full
Operating Kit redesign. Deferring the desired-state architecture would cause guard patches to be
rewritten. Including every audit recommendation would dilute validation and migration review.

State: `approved` by the user on 2026-07-09.

## Open Questions And Assumptions

### Open Questions

#### OQ-001 - Broad Distribution Trust Target

Question: What signing, notarization, or attestation target should govern a future broad-public
release?

Owner: user and future release-governance work.

`BLOCKER: no`

Affected IDs: `NFR-006`, `D-004`, `R-006`.

Resolution path: keep the current unsigned internal/prototype boundary explicit; resolve trust
policy before a later broad-public release plan.

#### OQ-002 - Python Oracle Decommission Evidence

Question: What minimum consumer migration evidence should close the legacy Python oracle and
resource-mirror retirement gate?

Owner: future decommission discovery or implementation planning.

`BLOCKER: no`

Affected IDs: `NFR-004`, `R-007`.

Resolution path: define a bounded compatibility fixture set, named migration evidence, and rollback
window after the Go state/release architecture becomes authoritative.

### Assumptions

- `A-001`: GitHub Releases remains the initial external distribution catalog, but the architecture
  should not hardcode catalog semantics into consumer state.
- `A-002`: existing `.codeheart/kit.lock.yaml` and `.codeheart/kit.config.yaml` paths should remain
  stable in the first implementation.
- `A-003`: existing valid optional config such as portfolio and repo-feedback state is
  consumer-owned and must survive repair and upgrade.
- `A-004`: no current requirement justifies switching machine state from YAML to JSON.
- `A-005`: current macOS universal and Windows x64 targets remain the first-class release targets.
- `A-006`: the source repository may use producer-specific routing as a repository-owned exception
  without changing consumer root doctrine.

## Risks And Mitigations

| ID | Risk | Impact | Mitigation and detection |
| --- | --- | --- | --- |
| `R-001` | Lifecycle refactor changes existing CLI behavior unexpectedly. | Consumers may receive different exit codes or repair behavior. | Preserve public command compatibility where safe; add a state-transition matrix and golden CLI fixtures. |
| `R-002` | Transaction rollback fails after partial filesystem mutation. | Consumer state may become inconsistent. | Stage on the target filesystem, use atomic replacement where supported, write a bounded transaction record, inject failures at every phase, and keep lock commit last. |
| `R-003` | Desired-state graph becomes a new abstraction that still duplicates raw manifests. | Complexity increases without removing drift. | Define graph compilation as the only runtime input and generate snapshots from it; prohibit command-specific surface registries. |
| `R-004` | YAML library or schema enforcement rejects previously tolerated files. | Existing consumers may need repair or migration. | Add compatibility fixtures from released versions, explicit schema migration, dry-run diagnostics, and no silent normalization. |
| `R-005` | Safe retirement of removed managed files deletes user-modified content. | User work may be lost. | Delete only prior managed paths whose digest still matches the prior lock; report conflicts and preserve modified files. |
| `R-006` | External release catalog is available but not authentic. | Attackers could substitute both asset and checksum if the publication authority is compromised. | Keep HTTPS/checksum baseline, record unsigned boundary, and leave a stable catalog-signature field for later trust policy. |
| `R-007` | Python oracle remains indefinitely and continues to shape new Go behavior. | Duplicate implementation and resource drift persist. | Make Go plus black-box compatibility fixtures authoritative after this implementation; schedule a separate decommission gate. |
| `R-008` | Producer mode drifts from real consumer behavior. | Maintainer tests pass while consumers fail. | Materialize a fresh consumer in CI and run install, sync, check, route, and upgrade probes against it. |
| `R-009` | Scope expands into documentation, profiles, signing, or platform work. | The first plan becomes hard to validate and complete. | Preserve `D-006` deferrals and require separate user approval for scope expansion. |

## Targeted Validation Strategy

Future implementation evidence should include:

- a complete state-transition matrix covering absent, adoptable, current, drifted, stale CLI,
  partial, schema-invalid, transaction-in-progress, and rollback states;
- preservation fixtures for repository instructions, config extensions, consumer plans, memory,
  local-user content, and modified formerly-managed files;
- full JSON Schema instance tests for component, profile, release, lock, and config contracts;
- YAML round-trip tests for numeric-looking strings, all-zero hashes, timestamps, nested maps, and
  unknown fields;
- deterministic desired-state graph snapshots and conflict diagnostics;
- injected failure tests before and after every mutation phase;
- managed-section checksum and semantic route validation;
- build-twice reproducibility on macOS and Windows release lanes;
- pack manifest, outer checksum, wrong-platform, wrong-version, missing-binary, malformed-archive,
  failed-smoke, rollback, and existing-install-preservation tests on both platforms;
- tag/source/binary/catalog version and provenance coherence checks without archive self-reference;
- a fresh producer-mode routing probe and a fresh materialized-consumer low-context routing probe;
- public-core, Markdown, schema, manifest, Go, migration, installer, and release validations;
- explicit evidence that no public release or named consumer sync occurred without separate
  approval.

## Reviewer Exchange Summary

Reviewer mode: documented main-thread architecture review after two separate reviewer-agent
attempts did not return a result. No separate reviewer output was treated as evidence. The review
challenged each decision against lifecycle safety, compatibility, source-of-truth boundaries,
release circularity, migration scope, and handoff rules.

| Decision | Reviewer critique | Main-agent response | Convergence |
| --- | --- | --- | --- |
| `D-001` | The lifecycle model omitted the mutating precondition of `update-check`, leaving a path that could still create partial state. | Added `update-check` as valid-install-only metadata maintenance that cannot initialize or repair. | Accept with revision; review-ready. |
| `D-002` | A new graph could become another source rather than the compiled runtime form of component/profile authority. | Clarified that validated declarations compile into the only runtime graph and generated snapshots derive from it. | Accept; review-ready. |
| `D-003` | Mature YAML parsing alone does not define schema-version compatibility or protect newer files from older CLIs. | Added schema-version dispatch, unsupported-future-version failure, no lossy normalization, and explicit dry-run migration. | Accept with revision; review-ready. |
| `D-004` | External provenance remained ambiguous unless the archive catalog and binary content identity were separated explicitly. | Declared the external catalog the sole archive-digest authority and limited embedded identity to non-circular content facts. | Accept with revision; review-ready. |
| `D-005` | Producer routing could weaken consumer dogfooding if canonical source bypassed installed behavior entirely. | Retained a fresh materialized-consumer CI path with full lifecycle and routing probes. | Accept; review-ready. |
| `D-006` | The first plan could omit migration safety or accidentally absorb every audit recommendation. | Added versioned migration to the required foundation and retained explicit Python, docs/profile, signing, platform, release, and sync deferrals. | Accept with revision; review-ready. |

Remaining dissent: none. Residual uncertainty is limited to the non-blocking signing target and
Python decommission evidence in `OQ-001` and `OQ-002`.

## Candidate Future Workstreams

These are discovery clusters, not implementation epics or task lists:

1. Consumer state safety and lifecycle semantics.
2. Desired-state compilation and contract validation.
3. Transactional reconciliation and managed ownership.
4. Release provenance, reproducibility, and upgrade flow.
5. Producer-mode routing and consumer materialization proof.
6. Deferred migration and decommission follow-ups.

## Manual Review Packet

### Approved, Delegated, Rejected, Deferred, Or Superseded Decisions

- Approved on 2026-07-09: `D-001` through `D-006` as recommended.
- Approved command default: retain current public commands, add `repair` and `upgrade`, make `init`
  reject existing installations, and make `sync` and `update-check` reject absent or invalid
  installations.
- Approved migration default: introduce lock schema v2, retain config schema v1 unless its
  structure changes, support dry-run lock-v1 migration, and never rewrite unsupported future
  versions.
- Approved delivery default: use one source implementation plan and stop before public release or
  consumer sync.

### No-Safe-Default Decisions

None. The current intention and evidence support a recommended default for every visible
implementation-shaping decision.

### Non-Blocking Follow-Ups

- `OQ-001`: future signing/notarization/attestation target.
- `OQ-002`: Python oracle decommission evidence.
- planning-document/index simplification and real profile modularity from the broader audit.

### Reviewer State

Architecture reviewer gate: complete through a documented main-thread pass after two separate
reviewer-agent attempts stalled. All six recommendations were revised where needed and are
approved for implementation planning.

### Exact User Action Needed

No further discovery decision is required before implementation planning. A later user action is
still required to approve execution of the draft implementation plan.

## Implementation Capability Scope - Consumer Lifecycle And State Safety

Capability:
The CLI can classify installed state and perform create, repair, current-version sync,
update-metadata check, trusted-version upgrade, and full validation through distinct safe
lifecycle operations.

Primary workflow:
A user or agent inspects the target, previews the requested lifecycle transition, approves the
operation when it changes state, receives a structured result, and can validate or recover the
installation afterward.

Must cover:

- `init`, `repair`, `sync`, `update-check`, `upgrade`, and `check` responsibilities from `D-001`;
- absent, adoptable, current, drifted, partial, schema-invalid, stale, transaction-in-progress,
  rollback, and unsupported-future-version states;
- preservation of repository instructions, optional consumer config, plans, memory, and local-user
  state;
- dry-run change plans, approval boundaries, stable blocker classes, and idempotent success.

Explicitly out of scope:

- public release publication and named consumer execution;
- broad onboarding UX redesign unrelated to lifecycle correctness.

Deferred or blocked:

- none inside the approved source implementation scope.

Preserve decisions:

- `D-001`, `D-006`.

Planner must not reinvent:

- command responsibilities, state preconditions, or the distinction between `sync` and `upgrade`;
- the rule that `update-check` cannot initialize or repair state.

Feature-level success evidence:

- a state-transition test matrix proves safe success, refusal, preservation, and recovery behavior
  for every supported starting state.

## Implementation Capability Scope - Desired State, Schemas, And Migration

Capability:
Validated component and profile declarations compile into the one typed desired-state graph used
by mutation, checking, lock generation, routing verification, and release snapshots.

Primary workflow:
The running CLI loads versioned declarations, validates them, resolves one profile-specific graph,
compares desired and observed state, and produces a deterministic change or validation result.

Must cover:

- typed ownership, presence, update, removal, condition, route, source, target, and digest fields;
- mature YAML parsing and full instance-schema validation;
- lock schema v1 reading and dry-run migration to lock schema v2;
- config schema v1 preservation unless implementation evidence requires an approved structural
  change;
- unsupported future-version stop behavior without normalization or downgrade;
- component and profile schema coverage.

Explicitly out of scope:

- new consumer paths, component families, and broad profile modularization;
- a JSON migration for lock or config files.

Deferred or blocked:

- config schema v2: deferred until a real structural config change is approved.

Preserve decisions:

- `D-002`, `D-003`, `D-006`.

Planner must not reinvent:

- YAML versus JSON, lock-v2/config-v1 defaults, or raw manifests versus a compiled typed graph.

Feature-level success evidence:

- schema, round-trip, migration, graph-snapshot, and unsupported-version tests prove deterministic
  and non-lossy state handling.

## Implementation Capability Scope - Transactional Reconciliation

Capability:
Every mutating lifecycle operation can preview, stage, validate, commit, and roll back managed and
generated state without silently overwriting consumer-owned content.

Primary workflow:
The lifecycle command resolves a deterministic change plan, reports conflicts, stages changes on
the target filesystem, validates the staged state, commits the lock last, and restores the prior
valid state on failure.

Must cover:

- transaction phases, durable or bounded recovery markers, structured evidence, and blocker shape;
- managed file and managed-section replacement;
- safe retirement of previously managed paths only when prior ownership and digest still match;
- conflict preservation for user-modified formerly managed files;
- failure injection before and after each phase;
- macOS and Windows filesystem behavior.

Explicitly out of scope:

- deletion of consumer-owned scaffolds, plans, memory, config, or local-user state.

Deferred or blocked:

- automated cleanup of legacy Python install libraries: deferred to decommission work.

Preserve decisions:

- `D-001`, `D-002`, `D-003`.

Planner must not reinvent:

- lock-last commit, safe managed-path removal criteria, or failure recovery expectations.

Feature-level success evidence:

- deterministic dry-run output and injected-failure tests prove atomic success, conflict handling,
  rollback, and idempotent retry.

## Implementation Capability Scope - Release Provenance And Upgrade

Capability:
Release packs are reproducible and externally verifiable, while approved upgrades stage a trusted
new binary, reconcile through that binary, validate the result, and roll back on failure.

Primary workflow:
An approved upgrade resolves the GitHub release catalog, verifies the archive digest and pack
identity, stages and validates the platform binary, invokes the new lifecycle engine, records
non-circular provenance in lock v2, and completes only after full check passes.

Must cover:

- external catalog as sole archive-digest authority;
- embedded non-circular version, component, content, and schema compatibility identity;
- deterministic ZIP creation and build-twice proof;
- pack version, platform, binary path, command identity, outer digest, and staged binary version
  validation;
- equivalent macOS and Windows failure, rollback, and prior-install-preservation tests;
- release and installer runbook alignment.

Explicitly out of scope:

- public release publication, tag creation, named consumer sync, signing/notarization automation,
  Linux, and Windows ARM.

Deferred or blocked:

- broad-public trust policy: deferred under `OQ-001` and remains a release-time gate.

Preserve decisions:

- `D-001`, `D-004`, `D-006`.

Planner must not reinvent:

- archive digest location, unsigned internal/prototype boundary, supported release platforms, or
  the release-publication stop.

Feature-level success evidence:

- byte-identical repeat builds, catalog/pack/binary coherence tests, and cross-platform staged
  upgrade failure tests prove reproducibility, provenance, rollback, and preservation.

## Implementation Capability Scope - Producer Mode And Routing Proof

Capability:
Maintainers use current canonical source authority while CI continuously proves that a freshly
materialized consumer receives correct managed routing and lifecycle behavior.

Primary workflow:
A maintainer follows producer-owned root routes into canonical source; CI builds the CLI,
materializes a clean consumer, runs lifecycle checks, and executes low-context routing probes.

Must cover:

- producer-specific source routing and tracked root managed-section consistency;
- local `.codeheart/` classification as non-authoritative generated state;
- consumer materialization, check, managed-section, and route-target validation;
- adoption normalization into the defined managed, repository-owned, and local-user layers;
- fresh low-context producer and consumer routing probes.

Explicitly out of scope:

- broad planning-index redesign and unrelated managed doctrine changes.

Deferred or blocked:

- Python oracle removal: deferred under `OQ-002`.

Preserve decisions:

- `D-002`, `D-005`, `D-006`.

Planner must not reinvent:

- canonical-source producer authority, generated consumer-fixture boundary, or the need for both
  producer and consumer routing proof.

Feature-level success evidence:

- source consistency checks and fresh materialized-consumer routing probes pass without relying on
  a stale local installed layer.

## Implementation Handoff State

This discovery is implementation-handoff-ready.

Frozen inputs:

- approved `D-001` through `D-006`;
- approved command compatibility default;
- approved lock-v2/config-v1 migration default;
- approved one-plan source delivery boundary;
- `FR-001` through `FR-012` and `NFR-001` through `NFR-010`;
- the five implementation capability scopes above;
- public release and named consumer sync remain separate approval-gated work.

Remaining follow-ups are non-blocking: `OQ-001`, `OQ-002`, planning/index simplification, broad
profile modularity, new platforms, and signing automation.

## Revision Notes

- 2026-07-09: Created the manual-review-ready draft from the repository architecture, automation,
  validation, documentation, and consumer-state audit. Recorded verified defects as requirements,
  separated six implementation-shaping decisions, and kept implementation planning, release, and
  consumer sync out of scope pending review and approval.
- 2026-07-09: Completed the architecture reviewer gate in the main thread after two reviewer-agent
  attempts stalled. Added `update-check` lifecycle boundaries, schema-version migration rules,
  non-circular catalog authority, affected-ID traceability, and explicit migration scope; advanced
  `D-001` through `D-006` to review-ready recommendations.
- 2026-07-09: Recorded user approval of `D-001` through `D-006`, the strict command compatibility
  default, lock-schema-v2/config-schema-v1 migration default, and one-plan source delivery boundary.
  Added five implementation capability scopes and advanced the discovery to
  implementation-handoff-ready while preserving release, consumer sync, Python retirement,
  signing, profile redesign, and new platform deferrals.
