Last updated: 2026-06-25T12:28:30Z (UTC)

# Codeheart Operating Kit Release Notes

## v0.1.13 Release Notes

`v0.1.13` adds a managed discovery-handoff gate for implementation planning.

### Included

- `draft-implementation-plan.md` now requires agents to check discovery handoff state before
  drafting normal implementation epics from a discovery document.
- Discovery-derived normal implementation plans now require implementation-handoff-ready
  discovery or recorded user approval, delegation, or revision of relevant
  `Implementation Capability Scope` blocks.
- `discovery-workflow.md` now points to the implementation-plan drafting runbook as the handoff
  enforcement point.
- `review-planning-document.md` now checks that discovery-derived plans do not bypass capability
  scope approval.
- Packaged planning-workflow resources include the updated managed runbooks.

### Consumer Impact

- `instruction-only change`: installed consumers receive clearer managed planning guidance when
  they sync or update the Operating Kit.
- No forced migration is required. Existing consumer-owned plans, execution logs, runbooks, and
  local state files are not rewritten.
- Consumers adopt the new guidance through normal sync or update.

### Validation

- Release-readiness validation is recorded in the discovery handoff gate execution log before
  public tagging, release publication, and named consumer sync.

## v0.1.12 Release Notes

`v0.1.12` adds managed tooling-readiness guidance for local environment blockers encountered
during module onboarding or operation.

### Included

- Agent-interface guidance now includes `handle-tooling-readiness.md`, one central hybrid runbook
  for missing local tools such as package managers, runtimes, CLIs, PowerShell modules, browser
  automation prerequisites, and document/PDF tooling.
- The managed root `AGENTS.md` block and installed fallback inventory now route agents to the
  tooling-readiness runbook before they install, repair, improvise setup, or declare a module
  capability unavailable.
- The runbook authoring standard now tells future runbooks to route generic local environment
  blockers through the managed readiness route while keeping module-specific install commands and
  live service preflight module-owned.
- Planning, execution, and planning-review runbooks now check tooling-readiness routing only when
  plans create or materially change runbooks that can encounter missing local tools.
- Structure-governance guidance now clarifies that it owns runbook placement while
  agent-interface owns durable runbook shape and generic tooling-readiness behavior.
- Packaged resources include the new readiness runbook and updated route files.

### Consumer Impact

- `instruction-only change`: installed consumers receive clearer managed guidance for local
  tooling blockers when they sync or update the Operating Kit.
- No default tool installation is performed. The baseline tooling catalog is on-demand, not an
  install bundle.
- No forced migration is required. Existing consumer-owned runbooks, module-owned runbooks, plans,
  execution logs, and local state files are not rewritten.
- Modules still own concrete module-specific commands, version requirements, authentication,
  permissions, and live external preflight.
- Consumers adopt the new guidance through normal sync or update.

### Validation

- Release-readiness validation is recorded in the tooling environment readiness execution log
  before public tagging, release publication, and consumer proof.

## v0.1.11 Release Notes

`v0.1.11` adds managed doctrine and generic routing for committed, non-secret module and extension
state under `docs/repo/state/<module-or-extension-id>/`.

### Included

- Structure-governance guidance now defines `docs/repo/state/<module-or-extension-id>/` as the
  standard location for committed repo-owned module/extension routing state.
- The managed root `AGENTS.md` block now tells agents to discover installed modules through the
  module system present in the repository and check `docs/repo/state/<id>/` before asking repeated
  setup questions.
- The installed kit fallback inventory now routes agents to the module/extension state reference.
- Placement and managed-boundary guidance now clarify that local committed state routes work but
  does not replace live external preflight.
- Packaged resources include the new structure-governance reference and updated route files.

### Consumer Impact

- `instruction-only change`: installed consumers receive clearer managed placement and routing
  guidance for module/extension state when they sync or update the Operating Kit.
- No forced migration is required. Existing consumer-owned module files, runbooks, plans,
  execution logs, and local state files are not rewritten.
- No `docs/repo/state/` folder is scaffolded by this release. Modules or extensions create
  `docs/repo/state/<id>/` only when they have real committed non-secret state to store.
- Consumers adopt the new guidance through normal sync or update.

### Validation

- Release-readiness validation is recorded in the module extension state routing execution log
  before public tagging, release publication, and consumer proof.

## v0.1.10 Release Notes

`v0.1.10` adds a managed runbook authoring standard and connects it to planning,
execution, review, and installed fallback routes.

### Included

- Agent-interface guidance now includes `runbook-authoring-standard.md` for human-facing,
  agent-facing, hybrid, and maintainer-facing durable runbooks.
- The standard defines compact intention blocks, user-visible flow requirements, agent execution
  path requirements, explicit approval and stop boundaries, and the narrow local language
  preference rule.
- Planning workflow runbooks now require runbook-authoring scope, audience class, intention-block
  coverage, execution checks, and review checks when plans create or materially change durable
  runbooks.
- Structure-governance guidance now routes durable runbook creation or material changes to the
  agent-interface standard.
- Packaged resources and fallback inventory include the new standard so installed consumers can
  discover it after normal sync or update.

### Consumer Impact

- `instruction-only change`: installed consumers receive stronger managed runbook authoring,
  planning, execution, and review guidance when they sync or update the Operating Kit.
- No forced migration is required. Existing consumer-owned runbooks, module-owned runbooks,
  plans, execution logs, and local state files are not rewritten.
- Consumers adopt the new guidance through normal sync or update.

### Validation

- Release-readiness validation is recorded in the runbook authoring standards execution log before
  public tagging, release publication, and consumer sync proof.

## v0.1.9 Release Notes

`v0.1.9` refines managed plan-register doctrine for local repositories and configured
coordination homes. It keeps the existing `docs/repo/plans/plan-register.md` location and adds
clearer reference shapes for portfolio overview usage.

### Included

- Plan-register format guidance now clarifies that canonical planning documents may live outside
  `docs/repo/plans/` while the durable register still lives at
  `docs/repo/plans/plan-register.md`.
- Local and coordination-home entry examples now show repository-qualified IDs, canonical document
  pointer shapes, compact relation fields, and coordination notes.
- Coordinated portfolios may use repository-qualified local IDs so a member register and the
  coordination-home register can refer to the same plan with the same stable ID.
- Coordination-home guidance now reuses already repository-qualified member IDs and derives
  `<SOURCE-NAMESPACE>-<LOCAL-ID>` only for bare member-local IDs.
- Plan-register guidance now clarifies that registers own compact index metadata and relation
  pointers while canonical planning documents own detailed decisions, blockers, execution
  evidence, and rationale.
- Repository grouping guidance now allows coordination-home registers to group entries by owning
  repository when that improves scanning.

### Consumer Impact

- `instruction-only change`: installed consumers receive clearer managed plan-register doctrine
  when they sync or update the Operating Kit.
- No forced migration is required. Existing consumer-owned plan registers, pending-sync files,
  canonical plans, local runbooks, and local board-like surfaces are not rewritten.
- Consumers adopt the new guidance through normal sync or update.

### Validation

- Release-readiness validation is recorded in the plan-register portfolio doctrine refinement
  execution log before public tagging or publishing.

## v0.1.8 Release Notes

`v0.1.8` clarifies coordination-home plan-register identity rules so multiple member repositories
can safely use local IDs such as `PR-001` while the coordination home keeps a unique portfolio
index.

### Included

- Plan-register format guidance now states that coordination-home entry IDs must be unique inside
  the coordination-home register.
- Member-repository entries promoted to a coordination home now use a stable source namespace plus
  the source local ID, such as `EXAMPLE-AUTOMATION-PR-001`.
- Coordination notes now preserve source local register IDs with
  `Source local register ID: <ID>`.
- Plan-register maintenance guidance now tells agents how to derive member namespaces from
  `portfolio.member_repository_id` and how to avoid copying bare member-local IDs into a
  coordination-home register.
- Coordination-home relation guidance now uses coordination-home IDs when represented entries are
  related inside the coordination-home register.

### Consumer Impact

- `instruction-only change`: installed consumers receive clearer managed plan-register doctrine
  when they sync or update the Operating Kit.
- No forced migration is required. Existing consumer-owned plan registers and pending-sync files
  are not rewritten.
- Consumers adopt the new guidance through normal sync or update.

### Validation

- Release-readiness validation is recorded in the coordination-home register ID namespace
  execution log before public tagging or publishing.

## v0.1.7 Release Notes

`v0.1.7` improves managed planning workflow guidance so discovery handoff,
implementation planning, planning-document review, and implementation execution preserve intended
feature capability instead of accepting policy-only or under-covered plans.

### Included

- Discovery handoffs now require implementation capability-scope blocks for implementation-relevant
  decision groups included in normal handoff.
- Implementation planning now requires intended feature capability coverage from accepted
  discovery when available, or from user request and targeted research when no discovery exists.
- Checklist guidance now emphasizes capability-sized concrete tasks, non-negotiable details, and
  exact preflight/remediation paths for legitimate execution-time variability.
- Planning-document review now checks for quiet capability narrowing, support-structure-only plans,
  avoidable non-concreteness, and the lazy implementer failure mode.
- Implementation execution review now checks delivered feature capability against epic outcome and
  discovery capability scope when available.

### Consumer Impact

- `instruction-only change`: installed consumers receive stronger managed planning workflow
  guidance when they sync or update the Operating Kit.
- No forced migration is required. Existing consumer-owned plans, discovery docs, execution logs,
  and local runbooks are not rewritten.
- Consumers adopt the new guidance through normal sync or update.

### Validation

- Release validation is recorded in the implementation-planning quality execution log before
  public tagging or publishing.

## v0.1.6 Release Notes

`v0.1.6` hardens plan-register session references and lifecycle organization as an additive
planning-workflow doctrine release.

### Included

- Plan-register maintenance now includes a self-contained bounded session-reference resolution
  procedure for agents that cannot read their own session ID directly from runtime context.
- The session-reference procedure is metadata-first: it uses local Codex session filenames,
  modification times, and first-record session metadata before any filename-only phrase check.
- Session-reference fallbacks are explicit: `not recorded`, `ambiguous: <reason>`, and
  `not confidently identified`.
- Plan-register format guidance now shows session-reference examples for identified, missing,
  ambiguous, and not-confident states.
- Plan-register format guidance now recommends lifecycle grouping inside one default
  `docs/repo/plans/plan-register.md` for active/draft, completed, and superseded/archived
  entries.
- No separate archive register is created by default.

### Consumer Impact

- `instruction-only change`: installed consumers receive clearer managed plan-register doctrine
  when they sync or update the Operating Kit.
- No migration is required. Existing consumer-owned plan registers and pending-sync files are not
  rewritten.
- Consumers adopt the new guidance through normal sync or update.
- Plan-register session-reference guidance does not depend on agent-memory or session-ledger docs.

### Validation

- Release-readiness validation is recorded in the implementation execution log before public
  tagging or publishing.

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
