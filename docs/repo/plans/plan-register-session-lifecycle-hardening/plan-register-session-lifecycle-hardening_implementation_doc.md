Last updated: 2026-06-21T19:29:39Z (UTC)
Created: 2026-06-21
Status: completed
Completed: 2026-06-21
Execution log: plan-register-session-lifecycle-hardening_execution_log.md

# Document Header

## Plan Register Session And Lifecycle Hardening Implementation Plan

This implementation plan hardens the Operating Kit plan-register doctrine after the first
`v0.1.5` consumer rollout exposed two gaps:

- session references are required for plan recovery, but agents do not have a direct live-runtime
  way to read their own session ID;
- completed, superseded, and archived planning entries need clearer durable register handling
  without creating a separate archive file by default.

The plan-register workflow must stay self-contained. It must not depend on agent-memory or
session-ledger doctrine because the plan register is intended to become the durable formal
planning and session-recovery surface.

Target release preparation: `v0.1.6`, unless a maintainer reorders release numbers before
execution. Public tagging and publishing remain owned by the release runbook.

## Essential Context Reference Files

| Path | Reason |
| --- | --- |
| `AGENTS.md` | Public-core safety, read order, and maintainer instructions for Operating Kit changes. |
| `README.md` | Public repository purpose and public boundary. |
| `docs/README.md` | Public docs router that must expose this plan. |
| `docs/repo/README.md` | Repository-governance router that must expose this plan. |
| `docs/repo/plans/README.md` | Plan index that must expose this plan. |
| `docs/repo/reference/placement-contract.md` | Ownership model for managed content and kit-initialized consumer state files. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact classes for instruction-only, release-note, and consumer-adoption changes. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer procedure for changing managed docs, resources, tests, and release notes. |
| `docs/repo/runbooks/release-operating-kit.md` | Release procedure for version bump, assets, checksums, validation, tag, and GitHub release. |
| `components/planning-workflows/managed/reference/plan-register-format.md` | Source managed plan-register format to harden. |
| `components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Source managed plan-register maintenance runbook to harden. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md` | Packaged resource copy that must match the component source. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Packaged resource copy that must match the component source. |
| `components/planning-workflows/component.yaml` | Component version and managed-file manifest for planning workflows. |
| `release-notes.md` | Consumer-facing release notes for the hardening release. |
| `scripts/validate-markdown-headers.py` | Markdown timestamp and header validation gate. |
| `scripts/validate-public-core.py` | Public-core hygiene validation gate. |
| `scripts/validate-release-manifest.py` | Release manifest validation gate. |
| `tests/test_packaging_resources.py` | Packaged-resource parity coverage. |
| `tests/test_sync_check.py` | Existing consumer sync and check behavior coverage. |

## Table Of Contents

- [Section 1 - Foundation](#section-1---foundation)
- [Section 2 - Strategy](#section-2---strategy)
- [Section 3 - Execution Plan](#section-3---execution-plan)
- [Section 4 - Future Planning](#section-4---future-planning)
- [Revision Notes](#revision-notes)

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Publish an additive Operating Kit hardening release that makes plan-register maintenance
self-contained for session references and lifecycle organization.

The implementation is complete when:

- `maintain-plan-register.md` contains a bounded read-only session-ID resolution procedure that
  uses Codex local session metadata without reading transcript bodies by default;
- `maintain-plan-register.md` makes session-ID confidence explicit with clear fallback values;
- `maintain-plan-register.md` does not route plan-register maintainers to agent-memory or
  session-ledger docs for session-ID resolution;
- `plan-register-format.md` includes session-ref examples for identified, not recorded, ambiguous,
  and not confidently identified session IDs;
- `plan-register-format.md` includes lifecycle grouping guidance for active/draft, completed, and
  superseded/archived entries;
- the lifecycle guidance keeps `plan-register.md` as the default single durable register and does
  not introduce a separate archive file by default;
- source managed docs and packaged resource copies match;
- consumer-impact and release notes classify the change as instruction-only with no required
  migration;
- validators and focused tests pass;
- release handoff identifies `v0.1.6` as the intended patch release and names required consumer
  sync targets.

## 1.2 Project And Problem Context

Operating Kit `v0.1.5` introduced the plan-register model and optional portfolio coordination.
Real consumer use immediately exposed two practical questions.

First, the plan-register format asks for session refs, but agents cannot reliably introspect their
own session ID directly from visible runtime context. A bounded local metadata scan can usually
identify the current user session without reading transcript bodies, using
`$CODEX_HOME` or `$HOME/.codex`, dated session files, first-record `session_meta` payloads,
`payload.id`, `payload.thread_source`, `payload.source`, `payload.cwd`, timestamps, and file
modification time. When metadata is ambiguous, a narrow `rg -l` search for a unique current-turn
phrase can return filenames only. The current managed plan-register docs do not explain this.

Second, the register needs a better default way to keep completed, superseded, and archived plans
visible without turning the active register into noise. A separate archive file is too early as a
default because it creates split-brain risk. The first hardening step should keep one
`plan-register.md` and add lifecycle grouping guidance.

The user explicitly clarified that plan-register docs should not point to agent-memory or
session-ledger runbooks. The plan register is intended to become the durable formal planning and
session-recovery surface, so any needed session-ref guidance belongs directly in planning-workflow
managed docs.

## 1.3 Current State Analysis

Existing systems:

- `components/planning-workflows/managed/reference/plan-register-format.md` defines the register
  entry fields, lifecycle values, relation vocabulary, session refs, and anti-patterns.
- `components/planning-workflows/managed/runbooks/maintain-plan-register.md` defines local
  register updates, portfolio coordination, pending sync, and high-level session reference
  handling.
- Packaged resource copies exist under
  `src/codeheart_operating_kit/resources/components/planning-workflows/managed/`.
- `components/planning-workflows/component.yaml` declares both managed plan-register files.
- `release-notes.md`, `manifest.yaml`, `pyproject.toml`, installers, packaged resources, tests,
  fixtures, and GitHub Actions currently target `0.1.5`.
- The Operating Kit repository itself does not currently have a local
  `docs/repo/plans/plan-register.md`; discoverability for this planning work is through
  `docs/README.md`, `docs/repo/README.md`, and `docs/repo/plans/README.md`.

Problems to correct:

- Session-ref guidance says "when available" but does not define a public-safe, bounded way for
  an agent to identify its current session ID.
- The current plan-register docs do not define confidence states such as `ambiguous` or
  `not confidently identified`.
- The current lifecycle guidance lists status values but does not explain how to group completed,
  superseded, and archived entries inside the register.
- The current docs say the register does not replace session ledgers, which may be true during the
  transition, but the hardening must avoid making formal plan-register maintenance depend on agent
  memory surfaces.

Target systems:

- Plan-register managed docs become self-contained for formal planning session refs.
- Consumer repositories can use one `docs/repo/plans/plan-register.md` as the durable formal
  planning index with lifecycle grouping.
- Existing consumers can adopt the guidance through normal sync, with no migration required.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Keep the change narrow: update managed planning-workflow doctrine, packaged resource copies,
release notes, consumer-impact metadata where present, and validation surfaces. Do not change CLI
behavior unless tests reveal packaged resource parity requires a mechanical update.

Expected path changes:

```text
Codeheart-Operating-Kit/
  docs/
    README.md                                           # modify
    repo/
      README.md                                         # modify
      plans/
        README.md                                       # modify
        plan-register-session-lifecycle-hardening/
          plan-register-session-lifecycle-hardening_implementation_doc.md  # create
  components/
    planning-workflows/
      component.yaml                                    # modify version metadata when release surfaces update
      managed/
        reference/
          plan-register-format.md                       # modify
        runbooks/
          maintain-plan-register.md                     # modify
  src/
    codeheart_operating_kit/
      resources/
        components/
          planning-workflows/
            component.yaml                              # modify resource copy when version metadata updates
            managed/
              reference/
                plan-register-format.md                 # modify resource copy
              runbooks/
                maintain-plan-register.md               # modify resource copy
  release-notes.md                                      # modify
  manifest.yaml                                         # modify during release preparation
  src/codeheart_operating_kit/resources/manifest.yaml    # modify during release preparation
  pyproject.toml                                        # modify during release preparation
  src/codeheart_operating_kit/__init__.py                # modify during release preparation
  bootstrap.md                                          # modify during release preparation
  install.sh                                            # modify during release preparation
  install.ps1                                           # modify during release preparation
  scripts/build-release-assets.py                       # modify default version during release preparation
  .github/workflows/validate.yml                        # modify release asset version checks during release preparation
  tests/
    test_packaging_resources.py                         # inspect or extend
    test_sync_check.py                                  # inspect or extend
    fixtures/                                           # modify version fixtures during release preparation
```

## 2.2 Open Questions And Assumptions Requiring Clarification

OQ-1 - Exact release version

- BLOCKER: no
- Affects: EPC-04, EPC-05
- Unlocks: The exact version strings used in release surfaces.
- Recommended default: Use `0.1.6` as the next patch release because `0.1.5` is already published
  and this change is additive managed doctrine hardening.

OQ-2 - Plan-register replacement wording

- BLOCKER: no
- Affects: EPC-01
- Unlocks: How strongly the docs state that plan registers are replacing agent-memory/session-ledger
  planning recovery over time.
- Recommended default: State that plan-register guidance is self-contained for formal planning and
  does not depend on agent-memory/session-ledger docs. Avoid declaring immediate removal of
  agent-memory surfaces in this patch release.

OQ-3 - Separate archive file

- BLOCKER: no
- Affects: EPC-01, EPC-02
- Unlocks: Whether to create a default archive register.
- Recommended default: Do not add a separate archive file by default. Add lifecycle grouping in
  `plan-register.md` and defer archive splitting until real registers become too noisy.

OQ-4 - Consumer sync targets

- BLOCKER: no
- Affects: EPC-05
- Unlocks: Which repositories should receive the final published kit update after release.
- Recommended default: After release, sync at least HQ, AWS Platform, and Foundry because they are
  the first configured portfolio coordination set.

## 2.3 Architectural Decisions With Reasoning

AD-1 - Keep plan-register session guidance self-contained

1. Problem being solved: Formal plan-register maintenance needs session IDs, but routing to
   agent-memory/session-ledger docs would make the new formal planning surface depend on an older
   continuity surface.
2. Simplest working solution: Put bounded session-ID resolution directly in
   `maintain-plan-register.md` and examples directly in `plan-register-format.md`.
3. What may change in 6-12 months: Agent runtimes may expose session IDs directly, or the
   Operating Kit may retire parts of agent-memory/session-ledger.
4. Rationale for the chosen approach: The plan register is becoming the formal planning and
   recovery index, so it needs its own complete procedure.
5. Alternatives considered and why not chosen: Pointing to session-ledger maintenance was rejected
   because it preserves the wrong dependency direction.

AD-2 - Use metadata-first session-ID resolution

1. Problem being solved: Reading full local session logs can expose private transcript content and
   waste tokens.
2. Simplest working solution: Resolve `$CODEX_HOME`, fallback to `$HOME/.codex`, inspect dated
   session filenames and first-record metadata, exclude subagents, and use filename-only phrase
   search only when metadata is ambiguous.
3. What may change in 6-12 months: The CLI may expose a direct current-session identifier.
4. Rationale for the chosen approach: Metadata-first scanning is bounded, public-safe as doctrine,
   and practical in current Codex local state.
5. Alternatives considered and why not chosen: Always searching transcript content was rejected
   because it is heavier and less privacy-preserving.

AD-3 - Keep one register with lifecycle grouping

1. Problem being solved: Superseded and archived entries need a place, but a default archive file
   can create split-brain and hide important historical context.
2. Simplest working solution: Add recommended sections inside `plan-register.md` for active/draft,
   completed, and superseded/archived entries.
3. What may change in 6-12 months: Large repositories may need an optional archive register or
   generated filtered views.
4. Rationale for the chosen approach: One durable index is easier for agents and humans to trust
   during early adoption.
5. Alternatives considered and why not chosen: A separate `plan-register-archive.md` was rejected
   as a default because it adds coordination overhead before there is evidence of scale pressure.

AD-4 - Treat this as instruction-only consumer impact

1. Problem being solved: Consumers need better guidance, but no installed paths, sync behavior, or
   schemas need to change for the first hardening release.
2. Simplest working solution: Classify as `instruction-only change`, update release notes, keep
   sync behavior unchanged, and validate resource parity.
3. What may change in 6-12 months: A future release may add validation for register entry shapes
   or a helper CLI for session refs.
4. Rationale for the chosen approach: The requested behavior is doctrine clarification, not a new
   generated surface.
5. Alternatives considered and why not chosen: Adding a CLI command now was rejected because it
   broadens the release beyond the documented gap.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| E1 | Managed plan-register docs are self-contained for session refs and lifecycle grouping. | M | None |
| E2 | Packaged resources and manifests reflect the managed doc changes. | S | E1 |
| E3 | Release notes and consumer-impact handoff describe the hardening release. | S | E1 |
| E4 | Version surfaces and tests are prepared for `v0.1.6`. | M | E2, E3 |
| E5 | Validation passes and release handoff is ready. | M | E4 |

## A) Epic ID, Title, And Outcome

E1 - Harden Managed Plan-Register Doctrine

Outcome: Source managed plan-register docs explain session-ID resolution and lifecycle grouping
without relying on agent-memory or session-ledger docs.

## B) Scope

Modify only source managed planning-workflow docs for the doctrine change.

## C) Files Touched

```text
components/planning-workflows/managed/reference/plan-register-format.md     # modify
components/planning-workflows/managed/runbooks/maintain-plan-register.md    # modify
```

## D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `maintain-plan-register.md` contains a bounded read-only session-ID resolution section.
- The session-ID section reads only metadata by default and names transcript-body search as a
  last narrow fallback that returns filenames only.
- The session-ID section defines confidence states and fallback text.
- `maintain-plan-register.md` does not tell maintainers to use agent-memory or session-ledger docs
  for plan-register session refs.
- `plan-register-format.md` includes session-ref examples for identified, not recorded,
  ambiguous, and not confidently identified states.
- `plan-register-format.md` recommends lifecycle grouping inside one `plan-register.md`.
- `plan-register-format.md` states no separate archive file is created by default.

## E) Dependencies And Critical-Path Notes

No dependency. This epic owns the primary user-requested behavior.

## F) Tasks Checklist

- [x] Update `components/planning-workflows/managed/runbooks/maintain-plan-register.md` with a self-contained session-ID resolution section.
- [x] Add metadata-first scan inputs to `components/planning-workflows/managed/runbooks/maintain-plan-register.md`.
- [x] Add confidence fallback values to `components/planning-workflows/managed/runbooks/maintain-plan-register.md`.
- [x] Confirm `components/planning-workflows/managed/runbooks/maintain-plan-register.md` does not route session-ref work to agent-memory docs.
- [x] Update `components/planning-workflows/managed/reference/plan-register-format.md` session-ref examples.
- [x] Add lifecycle grouping guidance to `components/planning-workflows/managed/reference/plan-register-format.md`.
- [x] Confirm `components/planning-workflows/managed/reference/plan-register-format.md` keeps one default `plan-register.md`.
- [x] Run `python scripts/validate-markdown-headers.py`.

## G) Implementation Notes

Use public-safe placeholders. Do not include private local absolute paths, session IDs, repo names,
or transcript snippets from consumer experiments. Example commands may use `$CODEX_HOME`,
`$HOME/.codex`, `YYYY/MM/DD`, and placeholder unique phrases.

## H) Open Questions

OQ-2 and OQ-3 apply, with safe defaults.

## A) Epic ID, Title, And Outcome

E2 - Sync Packaged Resource Copies

Outcome: Packaged resources match source managed docs, so installed consumers receive the hardened
guidance from the CLI package.

## B) Scope

Copy source managed doc changes into packaged resource mirrors and keep planning-workflows
metadata consistent.

## C) Files Touched

```text
src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-register-format.md   # modify
src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md  # modify
components/planning-workflows/component.yaml                                                                  # modify
src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml                            # modify
```

## D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- Source and packaged `plan-register-format.md` files match.
- Source and packaged `maintain-plan-register.md` files match.
- Source and packaged planning-workflows component manifests use the intended version metadata.
- Packaged-resource tests pass.

## E) Dependencies And Critical-Path Notes

Depends on E1. Do not update resource copies before source docs are stable.

## F) Tasks Checklist

- [x] Copy `components/planning-workflows/managed/reference/plan-register-format.md` to its packaged resource path.
- [x] Copy `components/planning-workflows/managed/runbooks/maintain-plan-register.md` to its packaged resource path.
- [x] Update source planning-workflows component version metadata.
- [x] Update packaged planning-workflows component version metadata.
- [x] Run `python -m pytest tests/test_packaging_resources.py`.

## G) Implementation Notes

Prefer exact file parity for managed docs. Do not hand-edit packaged copies differently from
source component files.

## H) Open Questions

OQ-1 applies, with safe default `0.1.6`.

## A) Epic ID, Title, And Outcome

E3 - Update Release Notes And Consumer Impact

Outcome: Maintainers and consumers can see the hardening change, its impact class, and its
adoption path.

## B) Scope

Update release-facing docs and impact records that already exist. Do not create new feedback or
issue templates.

## C) Files Touched

```text
release-notes.md                                      # modify
manifest.yaml                                         # modify during release preparation
src/codeheart_operating_kit/resources/manifest.yaml    # modify during release preparation
tests/fixtures/release-manifest.json                  # modify during release preparation
```

## D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- `release-notes.md` contains a `v0.1.6` section for plan-register session and lifecycle
  hardening.
- Release notes state consumer impact as `instruction-only change`.
- Release notes state no migration is required and consumers adopt through normal sync.
- Release notes state plan-register docs do not depend on agent-memory/session-ledger docs for
  session refs.
- Release manifest surfaces include the intended version during release preparation.

## E) Dependencies And Critical-Path Notes

Depends on E1 so release notes describe the final behavior.

## F) Tasks Checklist

- [x] Add `v0.1.6` plan-register hardening notes to `release-notes.md`.
- [x] Record consumer impact as `instruction-only change` in `release-notes.md`.
- [x] Record no required migration in `release-notes.md`.
- [x] Update root release manifest version surfaces during release preparation.
- [x] Update packaged release manifest version surfaces during release preparation.
- [x] Update release manifest fixture version surfaces during release preparation.
- [x] Run `python scripts/validate-release-manifest.py manifest.yaml`.

## G) Implementation Notes

If implementation remains docs-only until release preparation, release manifest updates can occur
in the release-prep epic. Keep release notes accurate even before public publishing.

## H) Open Questions

OQ-1 applies, with safe default `0.1.6`.

## A) Epic ID, Title, And Outcome

E4 - Prepare Version Surfaces And Focused Tests

Outcome: The repository is internally consistent for the intended `v0.1.6` patch release.

## B) Scope

Update package, installer, bootstrap, workflow, fixture, and test version surfaces that are
required for a normal release validation run.

## C) Files Touched

```text
pyproject.toml
src/codeheart_operating_kit/__init__.py
bootstrap.md
install.sh
install.ps1
scripts/build-release-assets.py
.github/workflows/validate.yml
profiles/standard.yaml
src/codeheart_operating_kit/resources/profiles/standard.yaml
components/*/component.yaml
src/codeheart_operating_kit/resources/components/*/component.yaml
tests/fixtures/*.yaml
tests/fixtures/*.json
tests/test_sync_check.py                         # inspect or modify
tests/test_install_metadata.py                   # inspect or modify
tests/test_release_assets.py                     # inspect or modify
```

## D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Package version is `0.1.6`.
- CLI `--version` reports `0.1.6`.
- Bootstrap and installers default to `0.1.6`.
- Component/profile/resource metadata uses `0.1.6`.
- Workflow release asset checks expect `0.1.6`.
- Fixtures align with the versioned release surfaces.
- Focused tests that cover install metadata, release assets, schemas, packaging resources, and
  sync/check pass.

## E) Dependencies And Critical-Path Notes

Depends on E2 and E3. This epic prepares release consistency but does not publish.

## F) Tasks Checklist

- [x] Update `pyproject.toml` to version `0.1.6`.
- [x] Update `src/codeheart_operating_kit/__init__.py` to version `0.1.6`.
- [x] Update `bootstrap.md` to version `v0.1.6`.
- [x] Update `install.sh` default version to `0.1.6`.
- [x] Update `install.ps1` default version to `0.1.6`.
- [x] Update `scripts/build-release-assets.py` default version to `0.1.6`.
- [x] Update `.github/workflows/validate.yml` release asset checks to `0.1.6`.
- [x] Update component metadata and packaged component metadata to `0.1.6`.
- [x] Update profile metadata and packaged profile metadata to `0.1.6`.
- [x] Update release and lock fixtures to `0.1.6`.
- [x] Run `python -m pytest tests/test_install_metadata.py tests/test_release_assets.py tests/test_packaging_resources.py tests/test_sync_check.py tests/test_json_schemas.py`.

## G) Implementation Notes

Use the existing `0.1.5` version-surface pattern. Do not publish tags or release assets in this
epic.

## H) Open Questions

OQ-1 applies, with safe default `0.1.6`.

## A) Epic ID, Title, And Outcome

E5 - Validate Release Readiness And Handoff

Outcome: The implementation is validated locally and ready for explicit release execution.

## B) Scope

Run local validations, build release assets if local prerequisites are available, and prepare a
release handoff. Do not tag, push, or publish.

## C) Files Touched

```text
dist/                                          # generated during local asset build
docs/repo/plans/plan-register-session-lifecycle-hardening/
  plan-register-session-lifecycle-hardening_execution_log.md  # create during execution
```

## D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Public-core validation passes.
- Markdown header validation passes.
- JSON schema validation passes.
- Release manifest validation passes.
- Focused pytest suite passes.
- Release asset build succeeds locally when prerequisites are available.
- Release handoff states that public tag, GitHub release, and consumer sync happen only after
  explicit user approval.
- Release handoff names HQ, AWS Platform, and Foundry as first consumer sync targets after
  publication.

## E) Dependencies And Critical-Path Notes

Depends on all prior epics. This is the review gate before public release execution.

## F) Tasks Checklist

- [x] Run `python scripts/validate-public-core.py`.
- [x] Run `python scripts/validate-markdown-headers.py`.
- [x] Run `python scripts/validate-json-schemas.py`.
- [x] Run `python scripts/validate-release-manifest.py manifest.yaml`.
- [x] Run focused pytest suite for metadata, packaging, release assets, sync/check, and schemas.
- [x] Run `python scripts/build-release-assets.py --version 0.1.6 --output-dir dist`.
- [x] Record release asset build result in the execution log.
- [x] Record release handoff for public tag and GitHub release.
- [x] Record post-release consumer sync targets for HQ, AWS Platform, and Foundry.

## G) Implementation Notes

If asset build prerequisites are missing, record the exact failure and keep the release handoff
blocked until the release runbook can complete validation.

## H) Open Questions

OQ-4 applies, with safe default consumer sync targets.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Add a CLI helper for current-session identification only after the metadata procedure proves
  stable across real consumer use.
- Add register entry validation only after the register format stabilizes through organic use.
- Add an optional archive register only after a real consumer register becomes too noisy for one
  file with lifecycle grouping.
- Revisit agent-memory and session-ledger retirement in a separate plan after the plan register
  has carried formal planning recovery for multiple repositories.

## 4.2 Future Considerations

- If future Codex runtimes expose the current session ID directly, replace metadata scanning with
  that direct source and keep the metadata scan as a fallback.
- Consumer repositories should be able to adopt this release through ordinary `sync`; no register
  migration should be required.
- The plan register should keep absorbing formal planning recovery behavior while informal
  continuity state remains outside this release scope.

# Revision Notes

- 2026-06-21: Created draft implementation plan for self-contained plan-register session-ref and
  lifecycle hardening, targeting a likely `v0.1.6` patch release.
- 2026-06-21: Activated plan for execution under user approval and added the sibling execution log.
