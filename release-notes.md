Last updated: 2026-06-21T15:40:16Z (UTC)

# Codeheart Operating Kit Release Notes

## v0.1.5 Release Notes

`v0.1.5` prepares the Operating Kit plan-register model and optional portfolio coordination as an
additive, runbook-first release.

### Included

- Managed planning-workflow docs now define `docs/repo/plans/plan-register.md` as the formal
  planning register for plans, plan families, lifecycle state, relationships, and lightweight
  session references.
- Managed planning-workflow docs now define
  `docs/repo/plans/coordination-sync-pending.md` as the local queue for coordination-home updates
  that could not be applied because the coordination home was unavailable.
- New consumers receive plan-register and pending-sync baselines during init/onboarding.
- Existing consumers receive the same baselines through normal sync only when the files are absent.
- `sync` refreshes an existing marked Operating Kit block in `AGENTS.md` so configured portfolio
  coordination routes to the plan-register maintenance runbook.
- The optional `portfolio` config schema is presence-based: no block means no configured portfolio
  coordination. Supported roles are `member` and `coordination-home`.
- Agent-memory guidance now keeps `goal-register.md` available for informal, pre-plan, or
  transitional continuity while formal plans move to `docs/repo/plans/plan-register.md`.

### Consumer Impact

- `instruction-only change`: installed consumers receive managed plan-register doctrine,
  planning-workflow hooks, the agent-memory boundary clarification, and the lean `AGENTS.md` route
  when they sync or update the Operating Kit.
- `backwards-compatible scaffold addition`: `docs/repo/plans/plan-register.md` and
  `docs/repo/plans/coordination-sync-pending.md` are created only when absent and are then
  consumer-owned state files.
- `validator-only change`: `.codeheart/kit.config.yaml` validation accepts optional role-specific
  portfolio configuration and rejects incomplete or contradictory portfolio blocks.
- Existing consumer-owned content is preserved. Sync does not rewrite existing `goal-register.md`,
  `plan-register.md`, or `coordination-sync-pending.md` content.
- No forced migration is required, no automatic goal-register rewrite is performed, and no normal
  onboarding prompt asks users to configure portfolio coordination.
- Optional consumer action: after upgrading, run normal sync/check to receive absent planning
  state files and refreshed managed instructions.

### Validation

- Local validation covers Markdown timestamps, public-core hygiene, JSON schema structure, release
  manifest structure, init/onboarding/sync behavior, packaged fallback behavior, source/package
  resource parity, and the full Python test suite.
- Focused local tests passed under Python 3.12.11: `49 passed`.
- Full local tests passed under Python 3.12.11: `86 passed`.
- Local script validation passed for Markdown timestamps, public-core hygiene, JSON schemas, and
  release manifest structure.
- `git diff --check` passed.
- Local release asset build produced `codeheart-operating-kit-0.1.5-macos.tar.gz` and
  `codeheart-operating-kit-0.1.5-windows.zip`.
- Local macOS installer smoke installed `codeheart-operating-kit 0.1.5` from the generated asset.
- Local checksum mismatch validation failed closed as expected.

## v0.1.4 Release Notes

`v0.1.4` adds a public-safe feedback intake workflow for Operating Kit consumers and maintainers.

### Included

- Public GitHub issue forms now provide sanitized intake for rough feedback, bugs, doctrine gaps,
  install/sync/check issues, docs routing issues, and feature requests.
- Maintainer feedback triage docs and label taxonomy now define feedback lifecycle states,
  discovery and implementation handoff, consumer-specific closure, and accidental disclosure
  response.
- Managed agent-interface docs now route consumers to public-safe feedback submission guidance and
  a reusable feedback item format.
- `init` and `sync` now ensure `.codeheart/user/feedback/` is ignored as optional local draft
  space without scaffolding that directory.

### Consumer Impact

- `instruction-only change`: installed consumers receive managed feedback-submission guidance when
  they sync or update the Operating Kit.
- `security or safety policy change`: public feedback intake now requires sanitization, warns
  against sensitive public submissions, and defines maintainer no-copy disclosure response.
- Repository-governance addition: public GitHub issue forms and feedback lifecycle labels are now
  part of the maintainer intake workflow.
- CLI impact: `init` and `sync` add `.codeheart/user/feedback/` to `.gitignore` so optional local
  feedback drafts remain local-only by default.
- Optional consumer action: run `codeheart-operating-kit sync <path>` after upgrading to refresh
  feedback-intake instructions in an installed consumer folder or repository.
- Migration required: none for already installed consumers. Normal sync/update adoption is enough.

### Validation

- Local validation covers public-core hygiene, Markdown timestamps, JSON schema structure, release
  manifest structure, full CLI tests, and release asset naming for v0.1.4.
- Feedback-intake validation covers issue-form YAML structure, required public-core confirmation
  checkboxes, GitHub label existence, managed-doc packaging parity, packaged fallback install of
  feedback docs, init/sync gitignore behavior, release manifest metadata, and per-epic reviewer
  gates.
- Full local CLI tests passed under Python 3.14.3: `74 passed`.
- Manifest parity checks confirmed source and packaged component manifests and profile manifest are
  byte-identical, and release-manifest component checksums match source component manifests.
- Local release asset build produced `codeheart-operating-kit-0.1.4-macos.tar.gz` and
  `codeheart-operating-kit-0.1.4-windows.zip`.
- Local macOS installer smoke installed `codeheart-operating-kit 0.1.4` from the generated asset.
- Local checksum mismatch validation failed closed as expected.

### Deferred

- CLI-assisted feedback drafting is deferred until real public issue examples prove the intake
  fields and local metadata needs.
- `.codeheart/user/feedback/` scaffold creation remains deferred; the path is optional ignored
  draft space, not a synchronized backlog.
- Private security reporting remains deferred to a focused security intake plan.

## v0.1.3 Release Notes

`v0.1.3` replaces thin managed planning, memory, and structure-governance guidance with the
consolidated public-core doctrine extracted from the mature AWS Platform repository.

### Included

- Managed planning-workflow docs now include fuller implementation-plan drafting, implementation
  execution, planning-document review, and planning-lifecycle doctrine.
- Managed agent-memory docs now distinguish curated continuity state from session transcripts and
  define goal-register, session-ledger, untriaged-session, and entry-format maintenance rules.
- Managed structure-governance docs now define repository documentation structure, durable naming
  and placement decisions, managed-content boundaries, documentation placement changes, and index
  maintenance.
- Consumer `docs/repo/README.md` scaffolding routes generic documentation and structure questions
  to managed Operating Kit governance while leaving repository-specific docs local.

### Consumer Impact

- `instruction-only change`: installed consumers receive stronger managed operating instructions
  when they sync or update the Operating Kit.
- Optional consumer action: run `codeheart-operating-kit sync <path>` after upgrading to refresh
  managed planning, memory, and structure-governance instructions in an installed consumer folder
  or repository.
- Review local duplicates after syncing: local planning, memory, or documentation-governance docs
  should become wrappers, repository-specific extensions, or archived historical references when
  the managed kit now owns the generic rule.

### Validation

- Local validation covers public-core hygiene, Markdown timestamps, JSON schema structure, release
  manifest structure, full CLI tests, and release asset naming for v0.1.3.
- Full local CLI tests passed under Python 3.12: `72 passed`.
- Manifest parity checks confirmed source and packaged component manifests, profile manifest, and
  release manifest are byte-identical, and release-manifest component checksums match source
  component manifests.
- Local release asset build produced `codeheart-operating-kit-0.1.3-macos.tar.gz` and
  `codeheart-operating-kit-0.1.3-windows.zip`.
- Local macOS installer smoke installed `codeheart-operating-kit 0.1.3` from the generated asset.
- Local checksum mismatch validation failed closed as expected.

### Deferred

- Product-specific or Foundry-specific operating profiles. This release keeps the `standard`
  profile general and instruction-only.

## v0.1.2 Release Notes

`v0.1.2` replaces the managed discovery workflow with the consolidated public-core workflow and
refreshes sync metadata when installed consumers update managed kit files.

### Included

- Managed `discovery-workflow.md` now includes the consolidated intention-led, normal discovery,
  and goal-style discovery workflow in one public Operating Kit runbook.
- Deep Discovery is now an explicit autonomous discovery mode that users can trigger with `/goal`
  plus a discovery document path. The runbook defines success criteria, decision-inventory passes,
  helper-agent research and review, reviewer exchange summaries, and manual-review readiness.
- Implementation-handoff readiness is separated from manual-review readiness so deferred or
  unresolved implementation-shaping decisions cannot silently become a normal single-path handoff.
- `sync` refreshes installed lock metadata from the currently installed CLI resources, including
  `kit_version`, selected component metadata, release metadata defaults, and managed-file
  checksums.

### Consumer Impact

- `instruction-only change`: installed consumers receive the expanded managed discovery workflow
  when they sync or update the Operating Kit.
- CLI impact: `sync` updates stale lock metadata to match the installed CLI version and managed
  resource checksums.
- Optional consumer action: run `codeheart-operating-kit sync <path>` after upgrading to refresh
  managed discovery instructions in an installed consumer folder or repository.

### Validation

- Local validation covers the consolidated discovery workflow copy, public-core hygiene, Markdown
  timestamps, JSON schema structure, release manifest structure, sync metadata refresh from a stale
  lock, unsupported-platform release metadata preservation, packaged-resource release metadata
  refresh, and release asset naming for v0.1.2.
- Full local CLI tests passed: `72 passed`.
- Local release asset build produced `codeheart-operating-kit-0.1.2-macos.tar.gz` and
  `codeheart-operating-kit-0.1.2-windows.zip`.
- Local macOS installer smoke installed `codeheart-operating-kit 0.1.2` from the generated asset.
- Local checksum mismatch validation failed closed as expected.

### Deferred

- Splitting the discovery workflow into shorter companion references. This release intentionally
  keeps one managed runbook as the canonical workflow source.

## v0.1.1 Release Notes

`v0.1.1` is an onboarding-contract hardening release for Codeheart Operating Kit.

### Included

- Public `README.md` `Start Here` prompt that keeps the first user prompt short and points to the
  latest public `bootstrap.md`.
- Hardened public `bootstrap.md` agent contract for language-first setup, explicit project-name
  selection, explicit target-folder selection, setup-plan preview, and write confirmation.
- Managed onboarding reference and runbook updates that make purpose optional context instead of a
  required setup branch for the `standard` profile.
- Neutral onboarding examples only: `Yourname-Automation`, `Companyname-Automation`,
  `Productname-Development`, `Teamname-Operations`, and `Existing-Project-Name`.
- Non-interactive onboarding policy that reserves `--yes` for explicit automation or repair flows
  where the user-owned decisions are already supplied.

### Consumer Impact

- `instruction-only change`: public bootstrap, public README, managed onboarding reference, and
  managed onboarding runbook now tell agents not to infer setup decisions.
- `validator-only change`: public-core validation now rejects legacy real-looking onboarding
  examples in public onboarding surfaces.
- `security or safety policy change`: non-interactive onboarding write behavior now fails closed
  when `--yes` is used without explicit target folder and project name decisions.
- CLI impact: `onboard` no longer supplies default target, project name, or purpose values;
  `init` can create a standard setup without purpose metadata.
- Schema impact: `setup_purpose` remains valid when present, but is no longer required for new
  `standard` profile configs.
- Installer impact: macOS and Windows installer defaults point to `v0.1.1` release assets.
- Known consumer action: none for installed consumers; existing `setup_purpose` metadata remains
  supported.

### Validation

- Local validation covers onboarding prompt order, explicit decision requirements, no-write
  behavior without `--yes`, config compatibility with and without `setup_purpose`, public-core
  hygiene, release manifest consistency, release asset build, and local macOS install from built
  assets.

### Deferred

- Purpose-specific installed profiles.
- Foundry module selection during onboarding.
- End-to-end external agent simulation.

## v0.1.0 Release Notes

`v0.1.0` is the first public Codeheart Operating Kit release.

## Included

- Public bootstrap instructions in `bootstrap.md`.
- macOS user-level installer in `install.sh`.
- Windows user-level installer in `install.ps1`.
- macOS release asset `codeheart-operating-kit-0.1.0-macos.tar.gz`.
- Windows release asset `codeheart-operating-kit-0.1.0-windows.zip`.
- SHA-256 checksum files for both platform assets.
- Release manifest `manifest.yaml`.
- Managed Operating Kit docs packaged inside the CLI wheel and described by component manifests.
- `standard` profile for G1 repository and agent-memory scaffolding.
- `codeheart-operating-kit` commands: `onboard`, `inspect`, `init`, `sync`, `check`, and
  `update-check`.
- Native Codex capability status reporting for documents, spreadsheets, presentations, browser,
  and PDF support.

## Consumer Impact

- `component addition`: first public component set, profile, schemas, validators, and managed docs.
- `instruction-only change`: first public operating instructions, runbooks, and references.
- `validator-only change`: first public validation scripts and CI validation matrix.
- `backwards-compatible scaffold addition`: consumer repository docs and agent-memory scaffolds
  are created only when absent.
- `consumer migration required`: existing repositories adopting the kit must run onboarding or
  initialization to add `.codeheart/` metadata and the managed `AGENTS.md` block.
- `breaking placement-contract change`: G1 establishes `.codeheart/kit/` as the managed kit root
  and intentionally does not create a `docs/workspace/` surface.

## Validation

- GitHub Actions validates Linux CLI tests, public-core hygiene, Markdown timestamps, JSON schema
  structure, release-manifest content, native capability tests, release asset build, macOS
  installer smoke, and Windows installer smoke.
- Public release smoke validation is available through the manual `Validate` workflow with
  `release_version` set to `v0.1.0`.
- Installers verify SHA-256 checksums and fail closed on checksum mismatch.

## Deferred

- Signed release artifacts.
- Linux installer asset.
- Homebrew, winget, npm, GitHub Packages, and CodeArtifact publication.
- Specialized profiles beyond `standard`.
- Foundry module selection during onboarding.
