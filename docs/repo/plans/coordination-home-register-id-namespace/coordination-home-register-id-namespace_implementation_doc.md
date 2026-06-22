Last updated: 2026-06-22T19:47:33Z (UTC)
Created: 2026-06-22
Status: completed
Execution log: docs/repo/plans/coordination-home-register-id-namespace/coordination-home-register-id-namespace_execution_log.md

# Document Header

## Coordination Home Register ID Namespace Implementation Plan

Overview: Add explicit Operating Kit doctrine for coordination-home plan-register IDs so member
repositories can safely have local entries such as `PR-001` without colliding in the coordination
home. The change is instruction-only: it updates managed plan-register format and maintenance
guidance, mirrors the packaged resources, records release notes, and prepares a patch release.

Essential context:

| Source | Why it matters |
| --- | --- |
| `AGENTS.md` | Operating Kit maintainer routing, public-core safety, and release boundaries. |
| `README.md` | Public repository purpose and managed-content ownership boundary. |
| `docs/repo/runbooks/change-operating-kit.md` | Required procedure before changing kit source docs and packaged resources. |
| `docs/repo/runbooks/release-operating-kit.md` | Required procedure before publishing a new Operating Kit release. |
| `docs/repo/reference/placement-contract.md` | Managed content and consumer-owned state placement rules. |
| `docs/repo/reference/consumer-impact-classification.md` | Consumer impact class and release-note requirement for managed instruction changes. |
| `components/planning-workflows/managed/reference/plan-register-format.md` | Source managed field contract for local and coordination-home plan registers. |
| `components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Source managed maintenance procedure for local and coordination-home register updates. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md` | Packaged resource mirror that installed consumers receive. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Packaged resource mirror that installed consumers receive. |
| `tests/test_packaging_resources.py` | Existing parity test for changed source and packaged resources. |
| `release-notes.md` | Consumer-facing release-note surface for the instruction-only patch release. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Update Codeheart Operating Kit plan-register doctrine so coordination-home registers have a clear,
repeatable identity rule for member-repository entries.

Completion is proven when:

- `plan-register-format.md` states that coordination-home entry IDs must be unique within the
  coordination-home register;
- `plan-register-format.md` defines how to preserve a member repository's local register ID in
  `Coordination note`;
- `maintain-plan-register.md` instructs member and coordination-home updates to use a stable
  source namespace when adding member entries to the coordination home;
- coordination-home relations use coordination-home IDs when referring to coordination-home
  entries;
- source managed docs and packaged resource mirrors match byte-for-byte;
- release notes classify the change as `instruction-only change` with no forced consumer
  migration;
- local validation covers Markdown headers, public-core hygiene, packaged-resource parity,
  focused tests, and release asset build readiness;
- a patch release can be published through the existing release runbook after explicit release
  execution approval.

## 1.2 Project And Problem Context

Operating Kit `v0.1.5` introduced local plan registers and optional portfolio coordination.
Coordination-home registers use the same repeated-entry shape as local registers and may include
selected member-repository entries that matter to portfolio coordination.

The current format requires an `ID`, `Owner / repository`, and `Canonical docs`, but it does not
state how to handle common local ID collisions across repositories. A local repository can use
`PR-001` for its first registered plan. Another member repository can also use `PR-001`. Those
local IDs are valid in their local registers, but they collide when copied directly into one
coordination-home register.

The intended reusable doctrine is:

- local register IDs remain local to the owning repository;
- coordination-home register IDs are unique inside the coordination-home register;
- member entries in the coordination home use a stable repository namespace plus the source local
  ID when that is the simplest clear identity;
- the source local ID remains traceable in `Coordination note`;
- `Owner / repository` and `Canonical docs` remain required disambiguators and source-of-truth
  pointers;
- managed examples stay generic and public-core-safe.

## 1.3 Current State Analysis

Current source state:

- `components/planning-workflows/managed/reference/plan-register-format.md` says `ID` is a stable
  local register identifier or canonical plan ID.
- The same reference says a coordination-home register uses the same entry shape and may include
  selected member-repository entries.
- The required fields already include `Owner / repository` and `Canonical docs`, which can
  disambiguate source ownership and canonical authority.
- `Coordination note` already exists for portfolio relevance and sync state, but it does not name
  the source local register ID as a standard use.
- `components/planning-workflows/managed/runbooks/maintain-plan-register.md` tells member repos
  to update the coordination-home register using the same entry format, but it does not define a
  collision-safe ID procedure.
- Source and packaged managed docs are mirrored under
  `src/codeheart_operating_kit/resources/components/planning-workflows/`.
- `tests/test_packaging_resources.py` already verifies source and packaged resource parity for the
  two files this change needs.

Target state:

- Coordination-home ID uniqueness is explicit doctrine.
- Member local IDs are preserved without becoming the coordination-home IDs when collision risk
  exists.
- A future agent can apply pending sync items from several repositories without inventing an ID
  convention.
- The change ships as a small instruction-only patch release, likely `v0.1.8` if `v0.1.7` remains
  the current published release when execution starts.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Make the smallest managed-doc change that resolves the ambiguity, then mirror, validate, and
release through the normal Operating Kit path.

Expected source tree:

```text
Codeheart-Operating-Kit/
  components/
    planning-workflows/
      component.yaml                                             # modify version metadata
      managed/
        reference/
          plan-register-format.md                                # modify ID and coordination note doctrine
        runbooks/
          maintain-plan-register.md                              # modify coordination-home update procedure
  src/
    codeheart_operating_kit/
      resources/
        components/
          planning-workflows/
            component.yaml                                       # modify version metadata mirror
            managed/
              reference/
                plan-register-format.md                          # modify mirror
              runbooks/
                maintain-plan-register.md                        # modify mirror
  tests/
    test_packaging_resources.py                                  # inspect existing parity coverage
    test_install_metadata.py                                     # release validation
    test_release_assets.py                                       # release validation
    test_sync_check.py                                           # focused managed sync validation
    test_json_schemas.py                                         # manifest and schema validation
  release-notes.md                                               # modify with v0.1.8 notes
  manifest.yaml                                                  # modify during release prep
  src/codeheart_operating_kit/resources/manifest.yaml             # modify during release prep
  pyproject.toml                                                 # modify during release prep
  src/codeheart_operating_kit/__init__.py                         # modify during release prep
  scripts/build-release-assets.py                                # modify default version during release prep
  bootstrap.md                                                   # modify release URLs during release prep
  install.sh                                                     # modify release URLs during release prep
  install.ps1                                                    # modify release URLs during release prep
  docs/
    repo/
      plans/
        coordination-home-register-id-namespace/
          coordination-home-register-id-namespace_implementation_doc.md  # created by this draft
          coordination-home-register-id-namespace_execution_log.md        # create at activation
```

## 2.2 Open Questions And Assumptions Requiring Clarification

OQ-1 - Target release version

- BLOCKER: no
- Affects: EP-03, EP-04
- Unlocks: Version surfaces and release asset names.
- Recommended default: Use `v0.1.8` when `v0.1.7` remains the current published release at
  execution start.

OQ-2 - Public release execution

- BLOCKER: no
- Affects: EP-04
- Unlocks: Tag creation, GitHub release publication, and consumer adoption proof.
- Recommended default: Treat approval to execute this plan as approval to prepare release assets,
  then require an explicit release-publish instruction before creating the public tag and GitHub
  release.

OQ-3 - Repository namespace spelling

- BLOCKER: no
- Affects: EP-01, EP-02
- Unlocks: The exact generic example text.
- Recommended default: Use generic examples such as `MEMBER-A-PR-001` and derive the namespace
  from `portfolio.member_repository_id` when present, with `Owner / repository` as the fallback.

## 2.3 Architectural Decisions With Reasoning

AD-1 - Use coordination-home-unique IDs, not global local IDs

1. Problem being solved: Local register IDs like `PR-001` collide when several member repositories
   are represented in one coordination-home register.
2. Simplest working solution: Require coordination-home IDs to be unique inside the coordination
   register and use a stable source namespace for member entries.
3. What may change in 6-12 months: A future CLI could generate IDs and validate collisions
   automatically.
4. Rationale: The register is a Markdown index today. A simple namespace rule is enough for
   humans and agents without adding schema or CLI behavior.
5. Alternatives considered and why not chosen: Requiring globally unique local IDs would create
   unnecessary coordination overhead in every member repository. Using canonical paths only would
   make relations harder to scan.

AD-2 - Preserve source local IDs in `Coordination note`

1. Problem being solved: Changing the ID in the coordination-home register can obscure the source
   member register entry.
2. Simplest working solution: Add `Source local register ID: <ID>` to `Coordination note` for
   member entries promoted into the coordination home.
3. What may change in 6-12 months: A future register schema could add a dedicated optional
   `Source local ID` field.
4. Rationale: `Coordination note` already exists for portfolio relevance and sync state, so the
   current format can carry this traceability without a schema change.
5. Alternatives considered and why not chosen: Adding a required field now would make existing
   registers feel stale and would be more than an instruction-only patch.

AD-3 - Keep the change instruction-only

1. Problem being solved: The ambiguity is procedural and format-level, not a CLI behavior defect.
2. Simplest working solution: Update managed docs, packaged mirrors, component version metadata,
   release notes, and validation evidence.
3. What may change in 6-12 months: A validator may later check coordination-home ID uniqueness if
   register volume justifies it.
4. Rationale: Consumers can adopt the doctrine through normal sync. Existing consumer-owned plan
   registers do not need migration.
5. Alternatives considered and why not chosen: Adding validation now would broaden a small
   doctrine fix into behavior work and require more fixture design.

AD-4 - Keep examples generic and public-core-safe

1. Problem being solved: The motivating use case comes from private Codeheart repositories, but
   Operating Kit managed doctrine is public.
2. Simplest working solution: Use neutral member names such as `Example-Automation` and
   `MEMBER-A-PR-001`.
3. What may change in 6-12 months: Public docs may add richer examples after real public consumer
   feedback.
4. Rationale: Generic examples explain the rule without exposing private topology.
5. Alternatives considered and why not chosen: Using real internal repository names would make the
   example clearer locally but would violate the public-core boundary.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Source managed docs define coordination-home ID namespace doctrine. | S | None |
| EP-02 | Packaged resources and component metadata mirror the source doctrine. | S | EP-01 |
| EP-03 | Release notes, version surfaces, and local validation are ready for `v0.1.8`. | M | EP-02 |
| EP-04 | The patch release is published and a first consumer can sync the new doctrine. | M | EP-03 |

## EP-01 - Source Managed Doctrine

### A) Epic ID, Title, And Outcome

EP-01 - Source Managed Doctrine

Outcome: Source managed plan-register docs define collision-safe coordination-home IDs and source
local ID traceability.

### B) Scope

Modify only the source managed plan-register format reference and maintenance runbook. Keep the
rule generic and instruction-only.

### C) Files Touched

```text
components/planning-workflows/managed/reference/plan-register-format.md     # modify
components/planning-workflows/managed/runbooks/maintain-plan-register.md    # modify
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- The format reference says coordination-home entry IDs must be unique within the coordination
  home.
- The format reference explains namespaced member IDs with a generic example.
- The format reference standardizes `Coordination note` traceability for source local register
  IDs.
- The maintenance runbook tells agents how to derive and apply a member namespace before adding
  member entries to the coordination-home register.
- The maintenance runbook tells agents to use coordination-home IDs in coordination-home relations.

### E) Dependencies And Critical-Path Notes

No dependencies. This epic defines the intended doctrine before resource mirroring and release
work.

### F) Tasks Checklist

- [x] Update `components/planning-workflows/managed/reference/plan-register-format.md` so the `ID` field distinguishes local-register IDs from coordination-home IDs.
- [x] Add a `Coordination-Home ID Uniqueness` section to `components/planning-workflows/managed/reference/plan-register-format.md` with a generic member entry example.
- [x] Update the `Coordination Notes` section in `components/planning-workflows/managed/reference/plan-register-format.md` with `Source local register ID: <ID>` guidance.
- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` member procedure steps to derive a stable namespace from `portfolio.member_repository_id`.
- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` coordination-home procedure steps to reject bare member-local IDs when adding member entries.
- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` relation guidance so coordination-home relations use coordination-home IDs.
- [x] Run `python3 scripts/validate-markdown-headers.py` and fix header failures in edited Markdown files.
- [x] Run `python3 scripts/validate-public-core.py` and fix public-core hygiene failures in edited managed docs.

### G) Implementation Notes

Use generic example repository names. Do not mention private repository names, local absolute
paths, tenant names, customer names, account IDs, tokens, or private business context.

### H) Open Questions

None.

## EP-02 - Packaged Resource Mirrors

### A) Epic ID, Title, And Outcome

EP-02 - Packaged Resource Mirrors

Outcome: Installed consumers receive the same managed doctrine because packaged resources mirror
the changed source docs.

### B) Scope

Mirror the two source managed docs into packaged resources and update planning-workflows component
version metadata for the patch release.

### C) Files Touched

```text
src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md   # modify
src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md  # modify
components/planning-workflows/component.yaml                                                                  # modify
src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml                            # modify
tests/test_packaging_resources.py                                                                            # inspect
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- Source and packaged `plan-register-format.md` are byte-identical.
- Source and packaged `maintain-plan-register.md` are byte-identical.
- Source and packaged planning-workflows component manifests use the target patch version.
- Existing parity test coverage remains sufficient for the changed files.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 so packaged resources mirror the final source wording.

### F) Tasks Checklist

- [x] Copy `components/planning-workflows/managed/reference/plan-register-format.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Copy `components/planning-workflows/managed/runbooks/maintain-plan-register.md` to `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md`.
- [x] Update `components/planning-workflows/component.yaml` component version to the target patch version.
- [x] Update `src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml` component version to the target patch version.
- [x] Inspect `tests/test_packaging_resources.py` and confirm the two changed managed docs remain covered by `test_changed_source_and_packaged_resources_match`.
- [x] Run `uv run --with pytest python -m pytest tests/test_packaging_resources.py` and fix source-package parity failures.

### G) Implementation Notes

Prefer exact file copy for mirrors. Do not hand-edit packaged resources differently from source
managed docs.

### H) Open Questions

None.

## EP-03 - Release Notes, Version Surfaces, And Validation

### A) Epic ID, Title, And Outcome

EP-03 - Release Notes, Version Surfaces, And Validation

Outcome: The repository is internally consistent and validated for a patch release that carries
the instruction-only doctrine change.

### B) Scope

Update release notes and release version surfaces. Run focused and full local validation before
release publication.

### C) Files Touched

```text
release-notes.md                                           # modify
manifest.yaml                                              # modify
src/codeheart_operating_kit/resources/manifest.yaml         # modify
pyproject.toml                                             # modify
src/codeheart_operating_kit/__init__.py                     # modify
scripts/build-release-assets.py                            # modify
bootstrap.md                                               # modify
install.sh                                                 # modify
install.ps1                                                # modify
tests/fixtures/release-manifest.json                       # modify when version fixture exists
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `release-notes.md` has a `v0.1.8` section describing the coordination-home ID namespace change.
- Consumer impact is recorded as `instruction-only change`.
- Release notes say no forced consumer migration is required.
- Version surfaces point at the target patch version.
- Local validation passes before tag publication.
- Release assets build locally with the target version.

### E) Dependencies And Critical-Path Notes

Depends on EP-02 so release notes and version surfaces describe the final changed files.

### F) Tasks Checklist

- [x] Add `v0.1.8` release notes to `release-notes.md` for coordination-home register ID namespace doctrine.
- [x] Record `instruction-only change` consumer impact in `release-notes.md`.
- [x] Record no forced migration and normal sync adoption in `release-notes.md`.
- [x] Update `pyproject.toml` package version to `0.1.8`.
- [x] Update `src/codeheart_operating_kit/__init__.py` package version to `0.1.8`.
- [x] Update `scripts/build-release-assets.py` default release version to `0.1.8`.
- [x] Update `.github/workflows/validate.yml` release smoke-test asset names and default release version to `0.1.8`.
- [x] Update `manifest.yaml` release metadata and URLs for `v0.1.8`.
- [x] Update `src/codeheart_operating_kit/resources/manifest.yaml` release metadata and URLs for `v0.1.8`.
- [x] Update `bootstrap.md`, `install.sh`, and `install.ps1` release URLs for `v0.1.8`.
- [x] Update `tests/fixtures/release-manifest.json` to match the `v0.1.8` release manifest.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 scripts/validate-json-schemas.py`.
- [x] Run `python3 scripts/validate-release-manifest.py manifest.yaml`.
- [x] Run `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_packaging_resources.py tests/test_install_metadata.py tests/test_release_assets.py tests/test_sync_check.py tests/test_json_schemas.py`.
- [x] Run `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest`.
- [x] Run `uv run --with pip --with setuptools --with wheel python scripts/build-release-assets.py --version 0.1.8 --output-dir dist`.
- [x] Run `git diff --check`.

### G) Implementation Notes

If the latest release version changes before execution, update the target version consistently
before editing release surfaces. Keep this as a patch release unless a maintainer intentionally
bundles more consumer-impacting work.

### H) Open Questions

OQ-1 applies.

## EP-04 - Release Publication And Consumer Sync Proof

### A) Epic ID, Title, And Outcome

EP-04 - Release Publication And Consumer Sync Proof

Outcome: The patch release is published, and a first consumer can adopt the new coordination-home
ID doctrine through normal update and sync.

### B) Scope

Follow `docs/repo/runbooks/release-operating-kit.md` for public release publication and record
post-release evidence. Run one consumer sync proof in a locally available consumer repository.

### C) Files Touched

```text
dist/                                                              # create release assets
docs/repo/plans/coordination-home-register-id-namespace/
  coordination-home-register-id-namespace_execution_log.md          # create at activation
```

Consumer sync proof touches an installed consumer's managed `.codeheart/kit/` snapshot through
the normal `codeheart-operating-kit sync` command after release publication.

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Public tag `v0.1.8` exists at the validated commit.
- GitHub release `v0.1.8` includes release notes, manifest, installers, assets, and checksums.
- Release runbook stop conditions are checked before publication.
- Consumer update-check detects the published version.
- Consumer sync refreshes managed planning-workflow docs.
- Consumer check reports no managed-content drift after sync.
- Execution log records release URL, asset names, checksums, validation evidence, and residual
  risk.

### E) Dependencies And Critical-Path Notes

Depends on EP-03. Public release publication requires explicit release execution approval at the
time this epic starts.

### F) Tasks Checklist

- [x] Re-read `docs/repo/runbooks/release-operating-kit.md` before release publication.
- [x] Confirm the validated commit matches the intended release commit with `git status --short` and `git rev-parse HEAD`.
- [x] Confirm `release-notes.md` covers the coordination-home ID namespace consumer impact.
- [x] Confirm `dist/codeheart-operating-kit-0.1.8-macos.tar.gz` exists from the validated asset build.
- [x] Confirm `dist/codeheart-operating-kit-0.1.8-windows.zip` exists from the validated asset build.
- [x] Confirm checksum files exist for both release assets.
- [x] Create public tag `v0.1.8` from the validated commit after explicit release publication approval.
- [x] Publish GitHub release `v0.1.8` with `bootstrap.md`, `install.sh`, `install.ps1`, `release-notes.md`, `manifest.yaml`, release assets, and checksum files.
- [x] Run `codeheart-operating-kit update-check` in one consumer repository after publication.
- [x] Run `codeheart-operating-kit sync <consumer-repository-path>` in the same consumer repository after update-check sees `v0.1.8`.
- [x] Run `codeheart-operating-kit check <consumer-repository-path> --json` and confirm managed-content drift is absent.
- [x] Create `docs/repo/plans/coordination-home-register-id-namespace/coordination-home-register-id-namespace_execution_log.md` with validation and release evidence.

### G) Implementation Notes

Use the existing release workflow and release asset patterns from `v0.1.7`. Do not publish from a
dirty worktree. Do not include private consumer names or local paths in public release notes.

### H) Open Questions

OQ-2 applies.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Automated validation for coordination-home ID uniqueness is deferred until real coordination
  registers show enough volume to justify parser and fixture work.
- A dedicated `Source local ID` field is deferred because the current repeated-section format can
  carry source local ID traceability in `Coordination note`.
- CLI-assisted pending-sync application is deferred; this plan only makes the manual register
  procedure unambiguous.
- Migration of existing coordination-home register entries is deferred. Existing consumer-owned
  registers can be refreshed when they are next materially updated.

## 4.2 Future Considerations

- If coordination-home registers become large, add a validator that checks ID uniqueness,
  `Owner / repository`, `Canonical docs`, and source-local-ID traceability.
- If multiple coordination homes need shared conventions, add a short reference on namespace
  design rather than overloading the maintenance runbook.
- If pending-sync files accumulate repeated updates for the same member entry, consider a compaction
  procedure that preserves history while applying only the latest coordination-home snapshot.

# Revision Notes

- 2026-06-22: Created draft implementation plan for coordination-home register ID namespace
  doctrine and likely `v0.1.8` instruction-only patch release.
- 2026-06-22: Activated plan for implementation after user approval.
- 2026-06-22: Updated local validation commands to use isolated `uv` environments because the
  default local Python lacked the required pytest and packaging tooling.
- 2026-06-22: Completed EP-04 local pre-publication checks for release notes, assets, checksums,
  checksum mismatch fail-closed behavior, and local macOS install from the built asset. Public tag,
  GitHub release, and consumer sync proof remain pending explicit release-publication approval.
- 2026-06-22: Completed EP-04 after explicit release-publication approval by publishing
  `v0.1.8`, verifying release assets, and proving update-check, sync, and check adoption in an
  isolated consumer repository.
