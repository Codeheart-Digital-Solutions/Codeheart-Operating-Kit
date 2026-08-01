Last updated: 2026-08-01T06:46:12Z (UTC)
Created: 2026-07-31
Status: active
Execution log: semantic-plan-catalog-coordination_execution_log.md

# Document Header

## Semantic Plan Catalog And Branch-Aware Coordination Implementation Plan

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
    capabilities:
        - semantic-plan-catalog
    catalog_metadata_updated: "2026-07-31T23:16:16Z"
    first_cataloged: "2026-07-31T23:16:16Z"
    id: codeheart-operating-kit.implementation.semantic-plan-catalog-coordination
    kind: implementation
    legacy_aliases:
        - OK-PR-028
    products:
        - codeheart-operating-kit
    purpose: Implement semantic plan metadata, guarded migration, branch-aware coordination, managed UX, and an approval-gated release.
    relations:
        - kind: depends-on
          target: codeheart-operating-kit.discovery.semantic-plan-catalog-coordination
    schema_version: 1
    strategic_themes:
        - cohesive-product-planning
```
<!-- END CODEHEART PLAN METADATA -->

Overview: Implement the approved semantic plan-catalog model as a compatibility migration inside
the self-contained Operating Kit. Canonical discovery, implementation, and family documents gain
stable metadata and semantic identifiers. Local plan views become derived CLI output. A configured
coordination home can discover enrolled repositories, refresh default-branch facts plus changed
plan observations from unmerged remote branches, and preserve its own strategic overlay separately
from source-derived data.

The plan deliberately avoids a second repository registry, committed generated catalogs, proactive
services, and a forced all-repository cutover. It builds compatibility readers first, migrates this
producer repository with optimistic source checks while ordinary branch work continues, and leaves
each additional repository migration under that repository's authority. Existing portfolio-v1
configuration remains readable. Existing `plan-register.md` files remain stable compatibility
entry points and legacy evidence; new authoring stops appending manual entries after a repository
enters mixed catalog mode.

After source completion and producer migration, the plan now carries the change through a patch
version bump, reproducible release-candidate construction, macOS and real-Windows release proof,
an explicit publication gate, public tag and release publication, and post-publication download
verification. It does not upgrade a coordination home or any member repository; that portfolio
rollout remains repository-owned follow-up work performed from the released Kit.

This plan also implements the approved publication boundary. A user's request to activate a plan,
or materially update an active plan, authorizes the canonical planning checkpoint and directly
required planning metadata to be committed and normally pushed on the unambiguous work branch.
That request does not authorize unrelated files, implementation code outside the approved scope,
pull-request creation, merge, release, force-push, branch deletion, or destructive Git actions.

This is a draft implementation plan, not execution authority. It does not authorize source edits,
semantic migration, commits, pushes, pull requests, releases, consumer upgrades, or changes in
other repositories until the user activates the plan. Activation will include the bounded plan
publication checkpoint defined above. Release preparation is part of the activated execution
scope, but tag creation and public publication pause at the `EP-07` release gate until the user
explicitly authorizes the target-version rule, presented release-candidate source tree, remote
validation actions, distribution boundary, and conditional publication after every gate passes. A
future activation request may include that explicit release authority. Cross-repository rollout
remains a separate authority gate.

Essential context:

| Source | Why it matters |
| --- | --- |
| `semantic-plan-catalog-coordination_discovery_doc.md` | Approved decisions, exact metadata and identity contracts, risks, and four frozen implementation capability scopes. |
| `AGENTS.md` | Producer source authority, public-core boundary, managed routes, and external-action safety. |
| `docs/repo/runbooks/change-operating-kit.md` | Required source-change workflow, impact classification, mirrors, indexes, and validation evidence. |
| `docs/repo/runbooks/release-operating-kit.md` | Required version, reproducibility, platform, digest-chain, distribution-boundary, publication, and release-evidence gates. |
| `docs/repo/reference/placement-contract.md` | Ownership rules for managed doctrine, repo-owned plans, scaffolds, generated views, and local-machine cache state. |
| `docs/repo/reference/consumer-impact-classification.md` | Separate instruction, validator, scaffold, migration, placement, and safety impact declarations. |
| `components/planning-workflows/managed/reference/planning-document-lifecycle.md` | Existing plan headers, lifecycle, bundles, family trigger, and register relationship to preserve during migration. |
| `components/planning-workflows/managed/reference/plan-register-format.md` | Current sequential identity and register fields that become compatibility evidence. |
| `components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Current agent-maintained local and coordination-home procedure that the derived catalog replaces. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Required route ownership, preflight, low-context selection, and external Git-host boundaries. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Maturity, evidence, structured output, failure, and recovery expectations for reusable coordination mechanics. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Audience, intention, input, approval, stop, evidence, and validation structure for new managed runbooks. |
| `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` | Basis for keeping durable mechanics in Go command packages instead of adding ad hoc scripts. |
| `schemas/kit-config.schema.json` | Current root config-v1 and portfolio-v1 compatibility contract. |
| `internal/cli/`, `internal/commands/`, `internal/state/`, and `internal/reconcile/` | Existing self-contained Go command dispatch, schema validation, deterministic result, and safe mutation foundations. |
| `components/planning-workflows/component.yaml`, `profiles/standard.yaml`, and `resources.go` | Managed/scaffold inventory, generated consumer surfaces, and embedded resource boundary. |
| `scripts/validate-json-schemas.py`, `scripts/validate-markdown-headers.py`, and `scripts/validate-public-core.py` | Existing schema, chronology-header, and public-safety validators to extend or reuse. |
| `manifest.yaml`, `internal/version/version.go`, `release-notes.md`, and `scripts/build-release-assets.py` | Current version authority, release narrative, and deterministic macOS/Windows release-pack path. |
| `.github/workflows/validate.yml` | Supported macOS and real-Windows validation surface. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Replace shared manual plan-register authority with a source-derived semantic catalog while
preserving existing plans, chronology, aliases, ordinary parallel development, and explicit
repository ownership.

Completion is proven when:

- discovery, implementation, and family records expose a marker-bounded YAML metadata contract
  below the H1 while the existing three-line header and H1 retain their current authority;
- stable lowercase ASCII IDs use `<repository>.<kind>.<stable-slug>`, remain independent of file
  names and mutable classifications, and preserve previous register IDs as aliases;
- metadata, header, identity, relation, family, and repository-ownership validators produce stable
  actionable results;
- `codeheart-operating-kit plans validate`, `plans list`, `plans inventory`, and `plans migrate`
  provide one self-contained local workflow with mixed-format legacy compatibility;
- absent catalog-adoption configuration behaves as `legacy`, `mixed` mode uses metadata for new
  plans without manual register appends, and `canonical` mode requires complete metadata coverage;
- `plan-register.md` remains a stable compatibility entry point, existing legacy entries remain
  readable evidence, and feature branches do not commit regenerated summaries;
- migration inventory records source commit and content hash, semantic decisions, ambiguity,
  confidence, aliases, branch ownership, conflicts, and deferrals before a guarded write occurs;
- concurrent edits and active-branch ownership cause skips and reevaluation instead of overwrite;
- portfolio config-v2 expresses stable member and home identities plus discovery sources while
  current portfolio-v1 paths remain readable;
- ordinary member enrollment requires authorized scope, default-branch Kit markers, member role,
  stable repository identity, and matching coordination-home identity;
- a coordination home also has a required stable member-repository identity and contributes its
  own formal plans to its factual catalog without a role array or self-registration entry;
- `codeheart-operating-kit portfolio configure` accepts explicit role, repository identity, home
  identity, provider-scope, and local-root inputs and creates only reviewed config and absent
  repo-owned portfolio surfaces;
- `codeheart-operating-kit portfolio scan` supports local Git and authenticated read-only
  GitHub-through-`gh`, executes no branch code, and stores no credential;
- each complete scan uses default branches as baselines and ingests only canonical plans added or
  materially changed on every accessible unmerged remote branch;
- scanner-owned bare mirrors refresh and prune remote refs without reading uncommitted work,
  including unpushed local branches, or changing developer worktrees;
- provider pagination, truncation, retry exhaustion, and permission failures are detected and
  prevent an incomplete scan from replacing the last complete cache;
- duplicate semantic IDs remain distinct source observations, branch disappearance retires the
  observation on the next complete scan, older accessible branch observations receive a visible
  stale label, and access failure does not replace the previous complete cache;
- coordination-owned strategic themes, relationships, families, priorities, and analysis remain
  separate from member-derived factual observations;
- current strategic analysis routes through an on-demand refresh and discloses freshness,
  duration, candidates, inaccessible sources, and completeness;
- activation and material active-plan-update instructions grant plan-only normal commit-and-push
  authority without a second push prompt and refuse every broader Git or release action;
- this producer repository is semantically migrated without a repository-wide freeze, content
  dates are preserved for metadata-only changes, and a branch-aware validation command proves its
  default-branch plus remote-overlay coverage before canonical mode;
- managed files, component manifests, profile surfaces, packaged mirrors, schemas, tests, and
  public documentation agree; and
- the next unused patch version is applied consistently after source completion, deterministic
  release packs and their complete digest chain pass macOS and real-Windows validation, the
  explicitly authorized release is published from the validated commit, and public assets pass
  post-publication verification without opening a pull request, merging, upgrading a consumer, or
  modifying another repository.

## 1.2 Project And Problem Context

The current plan-register model solved an earlier coordination problem by giving every repository a
local register and optionally synchronizing relevant entries to a coordination home. Practical use
has exposed three kinds of coupling:

- sequential IDs require branches to share allocation state;
- every plan touches the same Markdown register heading, timestamp, and insertion point; and
- coordination freshness depends on an agent remembering a second write outside the canonical
  plan.

The approved discovery separates those responsibilities. Canonical documents own plan facts.
Generated views supply local comprehension. A coordination-home scanner builds a complete factual
portfolio view at the time of analysis. Portfolio-owned overlays add strategy without changing
member-owned facts. Git visibility and scan freshness become observable provenance rather than
being folded into plan lifecycle.

The repository already has useful implementation foundations:

- a self-contained Go CLI with explicit command dispatch;
- maintained YAML and JSON Schema validation in `internal/state`;
- deterministic result and blocker types in `internal/reconcile`;
- managed component manifests and embedded packaged resources;
- scaffold preservation and local-machine placement rules;
- public-core and cross-platform validation; and
- established planning, routing, recipe, tooling-readiness, and producer-change runbooks.

The implementation should extend those foundations instead of creating a second runtime or
coordination service. The scanner is an L3 command surface over reusable Go packages. The managed
setup, maintenance, migration, and authoring instructions are L1 recipes. No L2 script asset is
added because the distributed Go CLI already owns stable execution, structured results,
cross-platform behavior, and compatibility.

## 1.3 Current State Analysis

Current behavior:

- `docs/repo/plans/plan-register.md` is repo-owned, manually updated, and treated as plan-list
  authority;
- register records may combine discovery and implementation documents under one sequential ID;
- no schema identifies plan metadata blocks, stable semantic IDs, source relations, or family
  records;
- planning documents have consistent human headers but no portable machine-readable core;
- the root config schema has portfolio-v1 path fields but no nested portfolio schema version,
  stable home identity, or discovery sources;
- `component_settings` accepts arbitrary component objects but no catalog-adoption state is
  defined;
- the CLI has flat commands only and no plan or portfolio command group;
- no package enumerates plans, reconciles legacy register evidence, computes branch overlays, or
  scans provider scope;
- no coordination-home scaffold separates source-derived local cache from committed strategic
  overlay;
- existing runbooks require manual register maintenance and do not grant activation-scoped plan
  publication authority; and
- current tests cover installation and routing but not plan metadata, portfolio membership,
  migration concurrency, or branch-aware source scanning.

Target boundaries:

- `internal/plancatalog` owns metadata extraction, canonical record discovery, semantic identity,
  validation, legacy fallback, local view models, inventory, and guarded metadata application;
- `internal/portfolio` owns nested config compatibility, enrollment predicates, source adapters,
  default baselines, branch overlays, factual catalog assembly, strategic-overlay validation, and
  atomic local cache replacement;
- `internal/commands` owns thin `plans` and `portfolio` command adapters;
- schemas own durable public metadata, catalog, migration-ledger, local-source, and strategic
  overlay shapes;
- the planning-workflows component owns generic doctrine, runbooks, references, templates, and
  coordination-home scaffolds;
- member repositories own their canonical planning documents and metadata;
- coordination-home repositories own committed portfolio configuration and strategic overlays;
- a coordination home uses one `coordination-home` role, has its own `member_repository_id`, and is
  automatically included as a factual catalog member;
- `.codeheart/local/portfolio/` owns rebuildable cache and machine-specific source configuration;
- `plan-register.md` remains a repo-owned compatibility entry point and is never regenerated by
  feature-branch authoring;
- each repository controls its own catalog-mode transition and migration writes;
- this plan migrates the Operating Kit producer repository only; and
- `EP-07` owns the producer version bump and release under the release runbook's explicit
  publication gate, while consumer upgrade, other-repository writes, and persistent scheduling
  stay outside this plan.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement seven dependency-ordered epics. Establish durable contracts first, then local catalog and
migration mechanics, portfolio scanning, managed UX, this repository's semantic migration, and
integrated proof before versioning and releasing the completed capability. The local catalog
remains useful without a coordination home, so portfolio features do not become a prerequisite for
ordinary plan authoring.

Expected source hierarchy:

```text
Codeheart-Operating-Kit/
  internal/
    plancatalog/                                      # create: plan facts and migration
      model.go
      metadata.go
      identity.go
      discover.go
      legacy.go
      validate.go
      view.go
      inventory.go
      migrate.go
      plancatalog_test.go
    portfolio/                                        # create: enrollment and scanning
      config.go
      membership.go
      source.go
      localgit.go
      mirror.go
      github.go
      scanner.go
      overlay.go
      store.go
      portfolio_test.go
    commands/
      plans.go                                         # create grouped local commands
      portfolio.go                                     # create grouped coordination commands
      util.go                                          # modify repeatable value options
      commands_test.go                                 # modify command contract coverage
    cli/
      cli.go                                           # modify grouped dispatch and help
      cli_test.go                                      # modify CLI UX coverage
    state/schema.go                                    # modify embedded schema constants
    version/version.go                                 # modify selected release version
  schemas/
    plan-metadata.schema.json                          # create canonical metadata contract
    plan-catalog.schema.json                           # create source-observation/cache contract
    plan-migration-ledger.schema.json                  # create guarded semantic decisions
    portfolio-local-sources.schema.json                # create machine-specific source contract
    portfolio-strategic-overlay.schema.json            # create coordinator-owned facts
    kit-config.schema.json                             # modify nested portfolio-v2 compatibility
  components/planning-workflows/
    component.yaml                                     # modify managed/scaffold inventory and impact
    managed/
      README.md                                        # modify route index
      reference/
        planning-document-lifecycle.md                 # modify metadata and catalog modes
        plan-register-format.md                        # modify compatibility-entry contract
        plan-catalog-format.md                         # create exact metadata/catalog contract
        portfolio-coordination-format.md               # create config, membership, and ownership
      runbooks/
        discovery-workflow.md                          # modify new-plan metadata and publication
        draft-implementation-plan.md                   # modify new-plan metadata and activation handoff
        execute-implementation-plan.md                 # modify activation publication boundary
        review-planning-document.md                    # modify metadata and authority review
        maintain-plan-register.md                      # modify into compatibility/local-view route
        configure-portfolio-coordination.md            # create member/home hybrid setup recipe
        refresh-portfolio-catalog.md                   # create read-only coordinator recipe
        migrate-plan-catalog.md                        # create guarded migration recipe
    scaffolds/
      plan-register.md                                 # modify stable entry-point baseline
      coordination-sync-pending.md                     # retain source compatibility; stop fresh install
      portfolio-README.md                              # create repo-owned coordination entry point
      portfolio-strategic-overlay.yaml                 # create repo-owned overlay baseline
  components/agent-interface/
    component.yaml                                     # modify routing inventory version
    managed/
      README.md                                        # modify coordination route
      reference/operation-routing-and-dispatch.md      # modify plan/portfolio route selection
  profiles/standard.yaml                               # modify generated surfaces
  manifest.yaml                                        # modify release content identity
  release-notes.md                                     # modify impact, migration, and adoption notes
  bootstrap.md                                         # modify public bootstrap release reference
  install.sh                                           # modify macOS release reference
  install.ps1                                          # modify Windows release reference
  scripts/build-release-assets.py                      # reuse deterministic pack/catalog builder
  templates/
    agents/AGENTS.managed-block.md                     # modify immediate managed routes
    consumer-docs/repo/README.md                       # modify plan and portfolio entry points
  tests/
    fixtures/
      plans/                                           # create valid, legacy, conflict, family cases
      portfolio/                                       # create config, membership, branch, gh cases
    test_json_schemas.py                               # modify schema-instance coverage
    test_packaging_resources.py                        # modify source/mirror inventory coverage
    test_routing.py                                    # modify low-context route probes
    test_release_assets.py                             # modify release-contract coverage
    test_install_metadata.py                           # modify install/upgrade release coverage
  .github/workflows/validate.yml                       # modify cross-platform catalog tests
  docs/
    README.md                                          # modify producer docs index where required
    repo/
      README.md                                        # modify current plan and impact summary
      reference/
        placement-contract.md                          # modify new owned/local paths
        consumer-impact-classification.md              # retain authority; no new class required
      runbooks/change-operating-kit.md                 # modify validation commands
      plans/
        README.md                                      # modify plan index
        plan-register.md                               # freeze entries after producer cutover
        semantic-plan-catalog-coordination/
          semantic-plan-catalog-coordination_discovery_doc.md
          semantic-plan-catalog-coordination_implementation_doc.md
          semantic-plan-catalog-coordination_execution_log.md       # create on activation
          attachments/
            producer-plan-migration-ledger.yaml                     # create during EP-05
            producer-plan-migration-inventory.json                  # create during EP-05
```

Expected consumer-owned and local-machine surfaces:

```text
<member-repository>/
  .codeheart/kit.config.yaml
    portfolio.schema_version: 2
    portfolio.role: member
    portfolio.member_repository_id: <stable-id>
    portfolio.coordination_home_id: <stable-home-id>
    component_settings.planning-workflows.plan_catalog_mode: <legacy|mixed|canonical>
    component_settings.planning-workflows.plan_catalog_cutover_revision: <pre-cutover-git-commit> # required in mixed mode
  docs/repo/plans/
    <plan-family>/
      <slug>_discovery_doc.md                         # canonical discovery metadata
      <slug>_implementation_doc.md                    # canonical implementation metadata
      README.md                                       # family record at second-sibling trigger
    plan-register.md                                  # stable entry point and legacy evidence

<coordination-home>/
  .codeheart/kit.config.yaml
    portfolio.schema_version: 2
    portfolio.role: coordination-home
    portfolio.member_repository_id: development-hq    # identity of this repository
    portfolio.coordination_home_id: <stable-home-id>
    portfolio.discovery.sources: [...]                # committed provider scope
  .codeheart/local/
    portfolio/
      sources.yaml                                    # machine-specific local Git roots
      catalog.json                                    # last complete rebuildable source scan
      git/<member-repository-id>.git/                  # scanner-owned bare remote mirrors
  docs/repo/portfolio/
    README.md                                         # repo-owned analysis entry point
    strategic-overlay.yaml                            # repo-owned families/themes/relations/priorities
```

Command surface:

| Command | Mode | Durable effect | Primary output |
| --- | --- | --- | --- |
| `plans validate [--remote-overlays] [path]` | read-only | scanner-owned local mirror refresh only with `--remote-overlays` | Validation results, source scope, and blockers. |
| `plans list [--format text|json] [path]` | read-only | none | Generated local view from metadata plus legacy fallback. |
| `plans inventory [--remote-overlays] --output <path> [path]` | local write | Explicit inventory artifact plus scanner-owned mirror refresh with `--remote-overlays`. | Source revisions, hashes, legacy evidence, ownership signals, and gaps. |
| `plans migrate --ledger <path> --dry-run [path]` | preview | none | Guarded change plan and skips. |
| `plans migrate --ledger <path> --yes [path]` | repository write | Metadata blocks for hash-matching ledger entries. | Applied changes, skips, validations, and rollback evidence. |
| `portfolio configure --role member --member-repository-id <id> --coordination-home-id <id> (--dry-run|--yes) [path]` | preview or repository write | Approved shared config fields only. | Change plan or applied member configuration and validation. |
| `portfolio configure --role coordination-home --member-repository-id <id> --coordination-home-id <id> [--github-owner <owner>]... [--local-root <path>]... (--dry-run|--yes) [path]` | preview or repository/local write | Approved shared config, machine-local sources, and absent repo-owned scaffolds. | Change plan or applied home configuration and validation. |
| `portfolio scan [--format text|json] [path]` | read-only external plus local cache | Atomic replacement of last-complete local cache. | Freshness, duration, members, candidates, branch observations, errors, and completeness. |

All grouped commands use the existing binary and error conventions. Mutation commands reuse
`internal/reconcile` change planning, blockers, containment, atomic staging, rollback, and shared
result fields. Read commands use the same blocker shape and add command-specific data without
pretending that local cache materialization is repository mutation.

Recipe and routing treatment:

| Surface | Maturity and owner |
| --- | --- |
| Metadata/catalog/migration packages | Internal implementation behind L3 commands, owned by `internal/plancatalog`; strict structs, embedded schemas, fixtures, and deterministic results. They are not independent recipe assets or L2 scripts. |
| Portfolio scanner and adapters | Internal implementation behind L3 commands, owned by `internal/portfolio`; provider-neutral interfaces, scanner-owned bare mirrors, local Git and `gh` adapters, bounded processes, read-only remote access, and no branch execution. They are not independent recipe assets or L2 scripts. |
| `plans` and `portfolio` commands | L3 thin command wrappers owned by `internal/commands` and dispatched by `internal/cli`. |
| Setup, refresh, and migration runbooks | L1 managed recipes owned by planning-workflows with explicit audience, preflight, approval, stop, evidence, and validation contracts. |
| Planning lifecycle hooks | Managed doctrine owned by planning-workflows; no independent workflow runtime. |
| Future schedules or events | Deferred triggers that may call `portfolio scan`; never a second scanner or catalog authority. |

Validation tiers and evidence:

- L1 runbooks require fresh-agent executability probes with selected route, resolved target,
  approval class, command, non-secret blocker or result, and recovery evidence;
- L3 commands require non-live unit and fixture tests, interface-contract tests, structured-result
  evidence, dry-run or preflight validation for writes and live reads, and approval-gated live
  evidence only when a configured authorized scope exists;
- internal Go packages are proven through the L3 command and package tests rather than assigned a
  separate recipe maturity level; and
- no L2 script asset or publication command is promoted in V1. Reconsider promotion only when an
  independently invoked cross-runtime recipe appears or repeated command use justifies an L4 tool
  surface.

Routing ownership:

- ordinary repository plan work routes through planning-workflows;
- coordination analysis routes through `refresh-portfolio-catalog.md` before a freshness claim;
- GitHub access routes to the existing authenticated `gh` session after tooling and auth preflight;
- missing `git` or `gh` routes through tooling readiness rather than implicit installation;
- config/scaffold writes route through `portfolio configure` and require its named approval;
- migration writes route through reviewed ledger plus `plans migrate --yes`;
- plan activation itself supplies the narrowly defined plan-checkpoint commit-and-push authority;
- source implementation in this producer routes through `docs/repo/runbooks/change-operating-kit.md`;
  and
- release preparation and publication route through `docs/repo/runbooks/release-operating-kit.md`;
- `EP-07` stops for explicit release authority before tag creation and public publication; and
- PR, merge, consumer sync, and other-repository changes route separately and remain unauthorized
  by this plan.

## 2.2 Open Questions And Assumptions Requiring Clarification

### OQ-1 - Portfolio Performance Targets

`BLOCKER: no`

Affects: `EP-03`, `EP-06`.

Decision unlocked: hard duration budgets, adapter concurrency, and future incremental caching.

Recommended default: use four concurrent read operations per provider, record per-source and total
duration, apply explicit process cancellation, and set no product-level freshness SLO until the
representative benchmark in `EP-06` exists. A current overview requires a complete scan from the
present coordination request, not merely a young cache timestamp.

### OQ-2 - Live GitHub Portfolio Probe Scope

`BLOCKER: no`

Affects: `EP-03`, `EP-06`.

Decision unlocked: real-provider evidence beyond deterministic fake-`gh` fixtures.

Recommended default: make fixture and command-contract tests mandatory. Run one live read-only
probe only when an explicitly configured coordination-home scope and authenticated `gh` session
are available during execution. Absence of that environment is disclosed residual evidence, not
permission to infer or create portfolio scope.

### OQ-3 - Additional Repository Rollout Order

`BLOCKER: no`

Affects: `EP-06` and later repository-owned migration plans.

Decision unlocked: which Codeheart-controlled member repository adopts mixed and canonical mode
after this producer.

Recommended default: complete the reusable source capability and this producer migration first,
then use a coordination-home dry-run inventory to prepare per-repository migration handoffs.
Repository writes occur only from an activated plan in the owning repository.

### OQ-4 - Legacy Register Evidence Retention

`BLOCKER: no`

Affects: `EP-02`, `EP-05`.

Decision unlocked: eventual archival or deletion of historic manual register entries.

Recommended default: keep existing register bodies in place as read-only compatibility evidence,
add a clear frozen-authority notice at producer cutover, and stop appending entries. Fresh
repositories receive a compact explanatory entry point. Historical removal is a later explicit
cleanup decision, not part of semantic migration.

### OQ-5 - Target Release Version

`BLOCKER: no`

Affects: `EP-07`.

Decision unlocked: the exact version written to producer, packaged-resource, installer, test, and
release surfaces.

Recommended default: select the next unused patch version after the latest validated public tag at
`EP-07` preflight. That is `v0.1.24` while `v0.1.23` remains the latest release. Stop and recompute
the target when another release lands before version selection; do not preserve a stale version
merely because this plan mentioned it. This exact preflight makes the question non-blocking during
planning.

### OQ-6 - Release Distribution Boundary

`BLOCKER: no`

Affects: `EP-07`.

Decision unlocked: whether the built assets may be publicly distributed in their recorded signing
state.

Recommended default: retain the repository's current HTTPS-plus-SHA-256 unsigned
internal/prototype boundary only through explicit approval recorded at the publication gate. A
different intended audience requires the signing and notarization evidence demanded by the release
runbook before publication. Release-candidate work can proceed before that gate, so the question
does not block earlier implementation.

Assumptions:

- the approved discovery remains the source for implementation-shaping decisions;
- root config schema version remains `1`; only the nested `portfolio.schema_version` becomes `2`;
- both portfolio roles require `member_repository_id`; `coordination-home` additionally requires
  `coordination_home_id` and automatically catalogs its own repository;
- absent `component_settings.planning-workflows.plan_catalog_mode` means `legacy`;
- `mixed` accepts metadata plus legacy fallback and requires metadata on newly authored plans;
- `mixed` records the exact pre-cutover Git revision so grandfathered plan paths and the frozen
  register are provenance-backed rather than inferred from mutable current entries;
- `canonical` rejects missing metadata for every formal discovery, implementation, and family
  record in scope;
- family authority remains `README.md` only after the current second-sibling trigger is met;
- Git default-branch, fetch, prune, and merge-base operations are available through the installed
  Git executable;
- the GitHub adapter uses `gh api` and `gh auth status` without token capture or persistence;
- adapter processes receive bounded arguments, use scanner-owned bare mirrors, exclude local heads
  and working-tree bytes, and parse fetched content as inert bytes;
- repeated `--github-owner` values become committed provider discovery sources, repeated
  `--local-root` values become `.codeheart/local/portfolio/sources.yaml` entries, duplicate inputs
  are idempotent, and conflicting existing identities block instead of being replaced;
- inaccessible member source makes a scan incomplete and leaves the last complete cache untouched;
- branch observations remain present while their refs are accessible, receive `stale: true` after
  30 days without a source-commit change, and disappear on the next complete scan after merge or
  deletion;
- pull-request data is optional enrichment and never selects the branch set;
- provider candidate repositories remain non-authoritative records in command output and cache;
- coordination-owned overlay data is never written by source reconciliation;
- this public producer migration includes public-safe metadata only;
- current work may continue on other branches during migration; and
- no release version is hard-coded before `EP-07`; its preflight resolves the next unused patch
  against the then-current public tags.

## 2.3 Architectural Decisions With Reasoning

### AD-1 - Marker-Bounded Metadata With Existing Header Authority

Problem: plans need machine-readable facts without invalidating human-readable chronology and
lifecycle conventions.

Simplest working solution: parse exactly one YAML block between explicit Codeheart plan-metadata
markers immediately below the H1. Keep the three-line header authoritative for created date,
lifecycle status, and last meaningful content update; keep H1 authoritative for title.

Six-to-twelve-month change: metadata schema v2 may add optional fields while preserving v1 parsing
and markers.

Rationale: the document remains understandable without tooling, migration is localized, and
catalog-only updates cannot silently rewrite historical content chronology.

Rejected alternative: YAML front matter would replace the current header contract; unbounded
embedded YAML makes deterministic extraction and conflict diagnosis harder.

### AD-2 - Semantic IDs Without Central Allocation

Problem: sequential register IDs collide on parallel branches and communicate little context.

Simplest working solution: validate lowercase ASCII `<repository>.<kind>.<stable-slug>` IDs, derive
the repository segment from configured `member_repository_id`, normalize words with hyphens, use
dots only between structural segments, and require explicit transliteration for non-ASCII source
words.

Six-to-twelve-month change: an optional disambiguator can be added only after evidence shows real
same-repository, same-kind, same-slug collisions that cannot be resolved semantically.

Rationale: IDs are understandable, independently assignable, stable across path and title changes,
and namespace-scoped without a coordination service.

Rejected alternative: UUIDs solve collision but harm conversation; random suffixes add noise
without present evidence; continued counters retain branch coordination.

### AD-3 - Separate Records And Derived Families

Problem: combined register entries obscure discovery versus implementation authority and make
family membership drift.

Simplest working solution: catalog discovery and implementation documents separately, associate
execution logs as evidence, and treat a second-sibling family `README.md` as a family record whose
member set is derived from each plan's `family` metadata.

Six-to-twelve-month change: additional formal record kinds require schema and lifecycle review.

Rationale: each record has one purpose and lifecycle while families remain first-class and
nonduplicative.

Rejected alternative: cataloging every log, attachment, and task creates strategy noise.

### AD-4 - Explicit Repository Catalog Modes

Problem: inferring cutover from the presence of one metadata block makes authoring behavior
ambiguous and risks a long dual-write period.

Simplest working solution: define
`component_settings.planning-workflows.plan_catalog_mode` with `legacy`, `mixed`, and `canonical`.
Absence means `legacy`; `mixed` stops manual appends and enforces metadata for new plans while
reading old evidence proven by `plan_catalog_cutover_revision`; `canonical` requires complete
metadata coverage.

Six-to-twelve-month change: a future mode may require stronger controlled vocabularies without
altering current repository state.

Rationale: one explicit adoption flag makes validation and agent behavior deterministic while
allowing ordinary work during migration.

Rejected alternative: a separate migration-state file duplicates config authority; metadata
presence inference cannot express coverage closure.

### AD-5 - Stable Register Entry Point, Uncommitted Generated Views

Problem: users need a predictable place to learn how plans are listed, but committed regenerated
summaries recreate the merge hotspot.

Simplest working solution: preserve `docs/repo/plans/plan-register.md` as a repo-owned explanatory
and legacy-evidence file. Generate current local views through `plans list` and never update the
register during ordinary plan authoring after mixed-mode cutover.

Six-to-twelve-month change: a UI or centrally generated report can consume the same catalog model.

Rationale: stable discovery and backwards compatibility remain while feature branches create only
their canonical plan files.

Rejected alternative: deleting the path breaks learned entry points; committing generated views
returns the original conflict pattern.

### AD-6 - Guarded Semantic Migration Ledgers

Problem: automatic register-row conversion cannot infer correct record boundaries, family,
relations, status evidence, and aliases, while blind plan edits can overwrite active work.

Simplest working solution: `plans inventory` records source revision, byte hash, legacy evidence,
and branch touch signals. An agent semantically completes a schema-validated ledger. `plans
migrate` writes only approved entries whose preconditions still match, skips contested entries,
and reports reevaluation ownership.

Six-to-twelve-month change: repositories may automate high-confidence classifications after
multiple reviewed migrations establish safe rules.

Rationale: semantic understanding is used where it adds value, and deterministic tooling owns
write safety and idempotency.

Rejected alternative: freezing all work is disproportionate; blind mechanical migration preserves
bad data and creates conflicts.

### AD-7 - Nested Portfolio V2 With V1 Read Compatibility

Problem: automatic membership requires stable identities and discovery scope, but a second central
repository list would duplicate configuration.

Simplest working solution: keep root config schema v1, add `portfolio.schema_version: 2`, stable
member/home IDs, and coordination-home discovery sources. Require `member_repository_id` for both
roles. A `coordination-home` repository automatically contributes its own plans to the factual
catalog and does not need a second role or self-registration entry. Interpret current path-only
portfolio objects as v1 compatibility inputs.

Six-to-twelve-month change: provider source kinds can grow through versioned adapter contracts.

Rationale: repositories remain self-describing, the home has the same plan-identity namespace as
every other repository, and a fresh Kit installation can configure either role without hidden
institutional knowledge or a multiple-role state machine.

Rejected alternative: a central member list creates another source of truth; provider visibility
alone silently enrolls unrelated repositories.

### AD-8 - Provider-Neutral Scanner With Two Initial Adapters

Problem: coordination must work now for local Git and GitHub while generic managed doctrine must
not become GitHub-only.

Simplest working solution: define a Go `Source` interface that returns repository/ref/file facts.
Resolve configured local roots to their remotes, refresh and prune scanner-owned bare mirrors below
`.codeheart/local/portfolio/git/`, and scan remote refs without reading local heads or worktrees.
Use authenticated `gh api` calls for GitHub discovery and enrichment with complete pagination,
bounded retry, rate-limit, and truncation detection. Run no branch code, change no developer
worktree, change no global authentication setting, and store no credential.

Six-to-twelve-month change: GitLab, Bitbucket, scheduled, and event-driven adapters may call the
same scanner contract.

Rationale: one factual catalog pipeline is portable, testable, and ready for additional triggers.

Rejected alternative: a GitHub App or webhook service adds hosting, credentials, events, retry,
and lifecycle complexity before it is needed.

### AD-9 - Baseline Plus Changed-Plan Branch Observations

Problem: full scanning of every branch duplicates inherited plans, while PR-only scanning misses
valid pushed work.

Simplest working solution: refresh each source, scan each member default branch fully, enumerate
every accessible unmerged remote branch, compute the merge base, inspect only added or materially
changed canonical plan paths, and retain each ref/commit observation independently. Provider
pagination exhaustion, response truncation, retry exhaustion, missing merge-base evidence, and
permission failure make the scan incomplete rather than silently reducing coverage.

Six-to-twelve-month change: content-addressed incremental reads may optimize large portfolios
without changing catalog semantics.

Rationale: the coordinator sees active pushed plans without relying on branch names, PR creation,
or agent announcements and without multiplying baseline content.

Rejected alternative: recency filters can hide long-running active work; PR filters confuse review
state with planning authority.

### AD-10 - Complete-Scan Atomic Cache And Separate Strategic Overlay

Problem: partial scans can falsely claim current completeness, and source refresh can overwrite
coordinator-owned strategic interpretation.

Simplest working solution: write `.codeheart/local/portfolio/catalog.json` only after all required
member reads complete. Report candidates and failures in the scan result. Validate but never mutate
`docs/repo/portfolio/strategic-overlay.yaml` during source reconciliation.

Six-to-twelve-month change: a durable remote catalog may replace the local cache only through the
same factual/strategic ownership split.

Rationale: last-known complete evidence remains distinguishable from a failed refresh, and
portfolio analysis survives factual reconciliation.

Rejected alternative: best-effort partial cache replacement makes absence ambiguous; storing
strategic fields in member observations creates overwrite conflicts.

### AD-11 - Grouped CLI Commands In The Existing Binary

Problem: six new flat top-level commands would crowd the current lifecycle CLI and obscure related
workflows.

Simplest working solution: add `plans` and `portfolio` top-level groups with explicit subcommand
dispatch and help in the existing small parser. Keep each handler thin over internal packages.

Six-to-twelve-month change: a command-tree helper can replace repeated dispatch code when a third
group demonstrates need.

Rationale: the UX remains discoverable without adding a CLI framework or another executable.

Rejected alternative: standalone scripts weaken release portability and result consistency; a new
CLI framework is unnecessary for two groups.

### AD-12 - Activation Is Plan-Checkpoint Publication Authority

Problem: active branch work is invisible to the coordinator until pushed, but a second approval
prompt weakens timely publication and broad implied Git authority is unsafe.

Simplest working solution: managed planning routes state that the user's activation request, and a
requested material update to an active plan, authorize an unambiguous work branch plus a commit and
normal push containing only the canonical plan and directly required planning metadata.

Execution surface: this remains L1 agent workflow governed by managed runbooks. V1 adds no plan
publication CLI. Static contract tests plus fresh-agent probes in isolated repositories record the
selected files, Git commands, local bare-remote result, refusal cases, and visibility outcome.

Six-to-twelve-month change: a pre-push intent channel may be evaluated only after measured need for
unpublished awareness.

Rationale: durable remote visibility follows the user's plan decision while PR, merge, release,
force-push, deletion, unrelated files, and unauthorized code retain separate gates.

Rejected alternative: requiring merge blocks legitimate branch execution; agent announcements
depend on discipline and add another transient authority.

### AD-13 - Repository-Owned Rollout, Not Global Write Authority

Problem: a complete portfolio design must eventually cover all controlled repositories, but this
producer plan cannot safely authorize arbitrary cross-repository edits.

Simplest working solution: implement reusable tooling and managed doctrine here, semantically
migrate this producer, generate dry-run portfolio migration handoffs later, and require an
activated plan in each target repository before writes.

Six-to-twelve-month change: a portfolio operator may approve a bounded multi-repository rollout
plan with explicit targets and rollback evidence.

Rationale: strategic overview becomes possible without weakening repository ownership or user
authority.

Rejected alternative: silent global migration turns a source implementation plan into broad
external mutation authority.

### AD-14 - Release Completes The Producer Plan Behind A Publication Gate

Problem: stopping at release readiness leaves the reusable capability unavailable to coordination
homes and members, while treating plan activation as blanket release authority would bypass the
producer's version, signing-boundary, validated-commit, and publication checks.

Simplest working solution: add one final epic that selects the next unused patch version only after
source completion, updates every authoritative and packaged version surface, builds each supported
pack twice, proves the complete catalog-to-binary digest chain, validates isolated macOS and real
Windows install/upgrade paths, and prepares the release. Before remote candidate push and Windows
validation, the epic pauses for explicit user authorization covering the resolved version,
presented source tree, remote validation, distribution boundary, and conditional tag/publication
after every release gate passes.

Six-to-twelve-month change: signed and notarized assets may replace the explicitly recorded
unsigned internal/prototype boundary without changing the plan-catalog architecture or release
authority model.

Rationale: one implementation record now reaches an installable public outcome, but irreversible
external publication still follows the existing producer release contract. The release remains
independent from coordination-home and member-repository rollout, so publication failure cannot
silently widen cross-repository authority.

Rejected alternative: hard-coding `v0.1.24` during planning can collide with parallel releases;
publishing automatically after source tests pass ignores platform, digest, signing-boundary, and
explicit-authority gates.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| `EP-01` | Versioned plan metadata, semantic identity, family, catalog, ledger, local-source, overlay, and config contracts compile and validate. | L | None |
| `EP-02` | Local grouped plan commands produce derived views and guarded compatibility migration. | XL | `EP-01` |
| `EP-03` | Config-driven portfolio discovery and branch-aware local/GitHub scanning produce one complete atomic catalog. | XXL | `EP-01`, local read model from `EP-02` |
| `EP-04` | Managed planning and coordination UX, scaffolds, routing, activation publication, and packaged surfaces agree with the runtime. | XL | `EP-01`, `EP-02`, `EP-03` command contracts |
| `EP-05` | This producer repository completes semantic inventory, optimistic migration, mixed-mode cutover, and canonical coverage closure. | XXL | `EP-02`, `EP-04` |
| `EP-06` | Integrated, cross-platform, low-context, privacy, migration, and release-candidate evidence proves the complete source capability. | XL | `EP-01` through `EP-05` |
| `EP-07` | The next unused patch version is reproducibly built, platform-validated, explicitly authorized, published, and verified from public assets. | XL | `EP-06` |

Critical path: `EP-01` -> `EP-02` -> `EP-03` -> `EP-04` -> `EP-05` -> `EP-06` -> `EP-07`.

`EP-03` design can begin after the `EP-01` record schema stabilizes, but integration depends on the
local parser and observation model from `EP-02`. Managed prose in `EP-04` must describe actual
command output and blocker codes rather than anticipated behavior. Producer migration in `EP-05`
must use the completed tool rather than one-off edits. `EP-07` must release the validated
post-migration source state and recompute its target version when the public tag set has advanced.

## EP-01 - Canonical Contracts, Identity, And Validation Foundation

### A) Epic ID, Title, And Outcome

`EP-01 - Canonical Contracts, Identity, And Validation Foundation`

Outcome: durable schemas and one Go read model define canonical plan facts, semantic IDs, family
authority, catalog observations, migration decisions, portfolio-local sources, strategic overlays,
and nested portfolio-v2 configuration without changing repository content.

### B) Scope

In scope:

- exact metadata marker extraction and YAML schema v1;
- existing header and H1 parsing;
- plan kind, required fields, optional strategic fields, relations, aliases, and timestamps;
- semantic ID grammar, explicit transliteration validation, and configured repository ownership;
- family record placement and derived membership rules;
- source observation, duplicate-ID, visibility, verification, freshness, and completeness fields;
- migration ledger preconditions and semantic-review evidence;
- local source and strategic overlay schemas;
- nested portfolio-v2 config, required repository identity for both roles, automatic home
  self-membership, plus v1 read compatibility;
- explicit `legacy`, `mixed`, and `canonical` catalog modes; and
- embedded schema access plus focused fixtures.

Out of scope:

- CLI mutation;
- repository migration;
- provider processes; and
- controlled vocabularies for optional strategic fields.

### C) Files Touched

- `internal/plancatalog/model.go`
- `internal/plancatalog/metadata.go`
- `internal/plancatalog/identity.go`
- `internal/plancatalog/discover.go`
- `internal/plancatalog/validate.go`
- `internal/plancatalog/plancatalog_test.go`
- `internal/state/schema.go`
- `schemas/plan-metadata.schema.json`
- `schemas/plan-catalog.schema.json`
- `schemas/plan-migration-ledger.schema.json`
- `schemas/portfolio-local-sources.schema.json`
- `schemas/portfolio-strategic-overlay.schema.json`
- `schemas/kit-config.schema.json`
- `tests/fixtures/plans/`
- `tests/fixtures/portfolio/`
- `tests/test_json_schemas.py`

### D) Acceptance Criteria And Size

Size: `L`.

Acceptance criteria:

- valid discovery, implementation, and family fixtures parse into deterministic typed records;
- missing markers, multiple blocks, misplaced blocks, unknown required-shape fields, invalid UTC
  timestamps, duplicate relation targets, malformed aliases, and unsupported schema versions fail
  with stable blocker codes;
- header/H1 mismatches are reported without copying their canonical values into metadata;
- metadata-only catalog timestamps are distinct from the existing `Last updated` header;
- semantic IDs enforce configured repository ownership and permitted kind segments;
- path, title, lifecycle, and optional taxonomy changes do not require ID changes;
- family membership is derived and second-sibling placement is validated;
- same-ID source observations remain a list and are never map-overwritten;
- root config-v1 fixtures remain valid;
- portfolio-v1 path objects remain readable;
- portfolio-v2 member and coordination-home fixtures require `member_repository_id`, home fixtures
  require `coordination_home_id`, and the home produces one automatic self-member record;
- absent catalog mode resolves to `legacy`; and
- every new schema compiles through both Go and repository schema validation.

### E) Dependencies And Critical-Path Notes

No implementation dependency. The approved discovery supplies the field and normalization
contract. Schema names and typed Go models must stabilize before downstream commands begin.

### F) Tasks Checklist

- [x] Create `schemas/plan-metadata.schema.json` with the approved required core, optional strategic fields, exact kinds, relation shape, alias shape, and UTC catalog timestamp format.
- [x] Create `schemas/plan-catalog.schema.json` with member facts, independent source observations, Git visibility, verification, freshness, candidates, scan errors, and completeness fields.
- [x] Create `schemas/plan-migration-ledger.schema.json` with source revision, byte hash, semantic decision, evidence, confidence, ambiguity, aliases, conflicts, deferral, and branch-owner fields.
- [x] Create `schemas/portfolio-local-sources.schema.json` and `schemas/portfolio-strategic-overlay.schema.json` with strict public-safe field boundaries.
- [x] Extend `schemas/kit-config.schema.json` with portfolio-v2 role contracts, v1 compatibility, discovery sources, required repository identities for both roles, home identity, and planning-workflows catalog modes.
- [x] Add embedded schema constants and loaders in `internal/state/schema.go`.
- [x] Implement canonical record structs and deterministic ordering in `internal/plancatalog/model.go`.
- [x] Implement exact marker, YAML, lifecycle-header, and H1 extraction in `internal/plancatalog/metadata.go`.
- [x] Implement semantic ID parsing, ASCII normalization checks, transliteration diagnostics, and repository-identity ownership in `internal/plancatalog/identity.go`.
- [x] Implement canonical plan and family enumeration in `internal/plancatalog/discover.go`.
- [x] Implement schema, identity, relation, family, header, alias, and catalog-mode validation in `internal/plancatalog/validate.go`.
- [x] Add valid, invalid, renamed, same-ID, family, legacy, and mixed-mode fixtures under `tests/fixtures/plans/`.
- [x] Add v1, v2-member, v2-home-self-member, missing-home-repository-ID, mismatched-home, local-source, and strategic-overlay fixtures under `tests/fixtures/portfolio/`.
- [x] Extend `tests/test_json_schemas.py` with positive and negative instance tests for every new durable schema.
- [x] Run `go test ./internal/plancatalog ./internal/state` and record the result in the execution log.
- [x] Run `python3 scripts/validate-json-schemas.py` and record the result in the execution log.

### G) Implementation Notes

- Use `go.yaml.in/yaml/v3` and the existing `internal/state` schema compiler. Do not expand
  `yamlmini`.
- Metadata extraction must operate on bytes and line boundaries; it must not render and rewrite a
  document during validation.
- Required metadata fields are `schema_version`, `id`, `kind`, `purpose`, `first_cataloged`, and
  `catalog_metadata_updated`.
- Optional fields are `family`, `products`, `capabilities`, `strategic_themes`, `relations`, and
  `legacy_aliases`.
- Scanner-derived repository, path, ref, commit, checksum, PR, and scan facts belong only to source
  observations.
- `role: coordination-home` is the single home role. It does not create a role array. Its required
  `member_repository_id` supplies the namespace for home-owned plans, and the scanner derives the
  self-member observation without a discovery-source entry.
- Lifecycle values stay exactly those defined by the existing planning lifecycle reference.
- Treat unrecognized metadata schema versions as unsupported, not as legacy absence.

### H) Open Questions

None. The discovery froze the contract fields, placement, identity grammar, record kinds, family
authority, and lifecycle separation.

### Validation Tasks

- [x] Prove metadata fixtures round-trip through typed parsing without changing canonical document bytes.
- [x] Prove parallel same-ID fixtures preserve separate ref and commit observations.
- [x] Prove existing config-v1 fixtures remain valid and portfolio-v1 objects remain readable.
- [x] Prove a coordination home validates and catalogs its own plans through its required repository identity.
- [x] Prove malformed metadata never falls back silently to legacy inference.

## EP-02 - Local Catalog Views And Guarded Migration Mechanics

### A) Epic ID, Title, And Outcome

`EP-02 - Local Catalog Views And Guarded Migration Mechanics`

Outcome: repositories can validate, list, inventory, and semantically migrate plans through grouped
commands while mixed-format legacy records remain discoverable and concurrent edits remain safe.

### B) Scope

In scope:

- legacy register parser and evidence reconciliation;
- generated local list in text and JSON;
- grouped `plans` CLI dispatch and stable help;
- inventory generation with Git revision, byte hash, branch-touch, metadata coverage, and legacy
  evidence;
- schema-validated semantic ledger application;
- dry-run, explicit `--yes`, path containment, atomic writes, rollback, idempotency, skip, and
  reevaluation results;
- catalog-mode enforcement; and
- frozen register-entry behavior after mixed cutover.

Out of scope:

- automatic semantic classification;
- migration writes outside the selected repository;
- removal of legacy register entries; and
- portfolio discovery.

### C) Files Touched

- `internal/plancatalog/legacy.go`
- `internal/plancatalog/view.go`
- `internal/plancatalog/inventory.go`
- `internal/plancatalog/migrate.go`
- `internal/plancatalog/plancatalog_test.go`
- `internal/commands/plans.go`
- `internal/commands/commands_test.go`
- `internal/cli/cli.go`
- `internal/cli/cli_test.go`
- `internal/reconcile/plan.go`
- `internal/reconcile/transaction.go`
- `internal/reconcile/reconcile_test.go`
- `go.mod`
- `tests/fixtures/plans/`

### D) Acceptance Criteria And Size

Size: `XL`.

Acceptance criteria:

- `plans validate` reports canonical metadata failures, legacy gaps, duplicate IDs, and mode
  violations with deterministic exit codes;
- `plans list` returns one separately typed discovery, implementation, and family row with title,
  kind, semantic ID or explicit legacy identity, lifecycle, family, and canonical path;
- JSON output is deterministic and usable as an agent data source;
- legacy register evidence fills only missing compatibility facts and never overrides valid
  canonical metadata;
- `plans inventory` includes every canonical document plus unpaired legacy evidence and records the
  exact source revision and byte hash;
- `plans migrate --dry-run` produces no plan writes;
- `plans migrate --yes` requires a valid reviewed ledger and applies only exact precondition
  matches;
- concurrent change, active-branch ownership, ambiguous evidence, invalid target placement, and
  unrelated dirty-file conditions produce skips or blockers without overwrite;
- a failed multi-file apply rolls back every change made by the command;
- transaction staging, backup, quarantine, and rollback mutations remain bound to the opened
  repository and target-parent handles even when path components are concurrently replaced;
- repeated apply is idempotent;
- a catalog-only insertion preserves the original `Last updated` value;
- `mixed` mode enforces metadata for new plans and stops manual register append instructions while
  retaining legacy fallback; and
- `canonical` mode fails until every formal record has valid metadata.

### E) Dependencies And Critical-Path Notes

Depends on all `EP-01` schemas and typed records. The migration engine must be complete before the
managed migration runbook and producer migration are authored as executable instructions.

### F) Tasks Checklist

- [x] Implement legacy register entry parsing and canonical-path reconciliation in `internal/plancatalog/legacy.go`.
- [x] Implement deterministic text and JSON local views in `internal/plancatalog/view.go`.
- [x] Implement Git revision, byte-hash, branch-touch, coverage, and legacy-gap inventory in `internal/plancatalog/inventory.go`.
- [x] Implement schema-validated ledger loading and optimistic precondition evaluation in `internal/plancatalog/migrate.go`.
- [x] Implement marker-bounded metadata insertion that preserves existing header bytes and meaningful content dates.
- [x] Implement idempotent multi-file migration through `internal/reconcile` staging, containment, rollback, and result types.
- [x] Implement grouped `plans validate`, `plans list`, `plans inventory`, and `plans migrate` dispatch in `internal/commands/plans.go`.
- [x] Extend `internal/cli/cli.go` with `plans` group help, subcommand help, argument errors, and exit-code routing.
- [x] Add mode enforcement for `legacy`, `mixed`, and `canonical` repository states.
- [x] Add stable blocker codes for malformed metadata, ambiguous legacy evidence, source mismatch, active-branch ownership, dirty overlap, invalid ledger, incomplete coverage, and rollback failure.
- [x] Add golden command-output tests in `internal/commands/commands_test.go`.
- [x] Add grouped help and invalid-subcommand tests in `internal/cli/cli_test.go`.
- [x] Add dry-run, idempotency, concurrent-edit, active-branch, dirty-file, rollback, chronology, and canonical-coverage fixtures under `tests/fixtures/plans/`.
- [x] Raise the source-build floor to Go 1.25 and use `os.Root` handle-relative operations for transaction containment and rollback.
- [x] Run `go test ./internal/plancatalog ./internal/commands ./internal/cli ./internal/reconcile` and record the result in the execution log.
- [x] Run a temporary-repository migration simulation and record the before-hash, after-hash, preserved content date, and second-run result.

### G) Implementation Notes

- Inventory output is evidence, not an approved migration ledger. The semantic agent must fill
  classifications and confidence explicitly.
- Active-branch ownership should be detected from reachable commits and changed paths, then
  surfaced for a branch owner to resolve. Do not attempt to edit another worktree.
- Dirty state is path-scoped. Unrelated dirty files must not be staged or included; an overlapping
  dirty plan is a blocker.
- `plans inventory --output` is the only command in this group that writes an explicit artifact
  without `--yes`; the destination is explicit and the write must not modify canonical plans.
  Remote-overlay validation and inventory may refresh only scanner-owned local-machine mirrors.
- Preserve old register identifiers as metadata aliases. Retain session references and
  coordination notes as ledger evidence until semantically classified.
- Text output should favor comprehension; JSON output owns exact automation fields.

### H) Open Questions

- `OQ-4` affects future legacy-register cleanup only. The implementation default is retention and
  a frozen-authority notice.

### Validation Tasks

- [x] Prove mixed views contain metadata records and legacy-only records without duplicate canonical paths.
- [x] Prove changed source hashes and overlapping dirty plans produce zero canonical writes.
- [x] Prove rollback restores byte-identical plan files after an injected migration failure.
- [x] Prove concurrent parent replacement cannot redirect transaction artifacts, backup,
  quarantine, restoration, or cleanup outside the opened repository and parent handles.
- [x] Prove canonical mode rejects one missing formal record and passes after its reviewed ledger entry applies.

## EP-03 - Config-Driven Portfolio Discovery And Branch-Aware Scanning

### A) Epic ID, Title, And Outcome

`EP-03 - Config-Driven Portfolio Discovery And Branch-Aware Scanning`

Outcome: a configured coordination home can discover matching Kit-enabled repositories through
authorized local and GitHub scopes, build complete default baselines plus plan-changing branch
overlays, preserve independent conflicts, and atomically publish a fresh local factual catalog.

### B) Scope

In scope:

- portfolio-v1/v2 config loading;
- explicit member and coordination-home configuration inputs with required repository identities;
- committed provider sources and machine-specific local sources;
- exact default-branch membership predicate plus automatic coordination-home self-membership;
- candidate-only unconfigured Kit repository reporting;
- provider-neutral source interface;
- bounded local Git adapter;
- authenticated read-only GitHub-through-`gh` adapter;
- scanner-owned bare-mirror refresh and prune without developer-worktree reads;
- complete provider pagination, truncation detection, bounded transient retry, and rate-limit
  blockers;
- default-branch full plan reads;
- accessible unmerged remote-branch enumeration;
- merge-base changed-path filtering;
- source observations, optional PR enrichment, conflicts, staleness, retirement, failures, metrics,
  and completeness;
- strategic-overlay preservation and validation;
- grouped `portfolio configure` and `portfolio scan`;
- branch-aware `plans validate --remote-overlays` and `plans inventory --remote-overlays`; and
- atomic complete-cache replacement.

Out of scope:

- webhook, GitHub App, daemon, scheduler, outbox, and announcement infrastructure;
- executing repository content;
- credential storage;
- implicit enrollment;
- full plan-body storage by default; and
- source reconciliation writes to member repositories.

### C) Files Touched

- `internal/portfolio/config.go`
- `internal/portfolio/membership.go`
- `internal/portfolio/source.go`
- `internal/portfolio/localgit.go`
- `internal/portfolio/mirror.go`
- `internal/portfolio/github.go`
- `internal/portfolio/scanner.go`
- `internal/portfolio/overlay.go`
- `internal/portfolio/store.go`
- `internal/portfolio/portfolio_test.go`
- `internal/commands/portfolio.go`
- `internal/commands/plans.go`
- `internal/commands/util.go`
- `internal/commands/commands_test.go`
- `internal/cli/cli.go`
- `internal/cli/cli_test.go`
- `internal/reconcile/plan.go`
- `tests/fixtures/portfolio/`

### D) Acceptance Criteria And Size

Size: `XXL`.

Acceptance criteria:

- a member is enrolled only after scope authorization, valid default-branch Kit evidence, valid
  member role, stable repository ID, and matching home ID all pass; ordinary consumers require the
  installed marker and lock, while the Operating Kit producer may use its schema-valid tracked
  content manifest without committing consumer installation state;
- a coordination home requires its own stable repository ID and contributes one self-member
  baseline without a second role or discovery-source entry;
- a feature branch cannot self-enroll a repository;
- unconfigured and mismatched Kit repositories appear only as candidates or exclusions;
- no central repository list exists;
- local Git and GitHub adapters return the same provider-neutral source model;
- member configuration requires explicit role, member repository ID, and home ID inputs;
- home configuration requires explicit role, member repository ID, and home ID inputs, commits
  repeated GitHub-owner scopes, and stores repeated local roots only under `.codeheart/local/`;
- configuration is idempotent for matching values and blocks conflicting existing identities
  without replacement;
- GitHub preflight uses the existing `gh` session and never reads, prints, logs, or persists a
  token;
- subprocess argument construction prevents shell interpolation and every process is cancellable;
- local Git scanning refreshes and prunes scanner-owned bare mirrors, scans remote refs only,
  ignores local heads and worktree bytes, and changes no developer repository refs;
- no fetched branch content executes and no global Git or `gh` setting changes;
- GitHub repository and branch enumeration consumes every page; response truncation, exhausted
  three-attempt transient retry, rate-limit exhaustion, and incomplete comparison evidence make the
  scan incomplete;
- each member default branch contributes a complete canonical baseline;
- every accessible unmerged remote branch is evaluated independently of name, age, and PR state;
- branch overlays contain only added or materially changed canonical plan files from the merge
  base;
- inherited baseline plans are not duplicated across branches;
- same-ID branch observations coexist with ref, commit, content hash, and conflict diagnostics;
- optional PR facts enrich but never select observations;
- accessible observations older than 30 unchanged days are visibly stale and remain retained;
- merged or deleted refs disappear after the next complete scan;
- any required-member read failure makes the new scan incomplete and leaves the prior complete
  cache byte-identical;
- a successful scan atomically replaces `.codeheart/local/portfolio/catalog.json`;
- coordinator strategic overlay bytes remain unchanged across every scan;
- scan results disclose start/end time, duration, source counts, member counts, candidates,
  observations, stale count, errors, and completeness; and
- `plans validate --remote-overlays` and `plans inventory --remote-overlays` reuse the same source
  engine, distinguish current-worktree evidence from pushed remote observations, and disclose that
  local-only evidence is not globally visible; and
- local fixture scans pass on macOS and Windows.

### E) Dependencies And Critical-Path Notes

Depends on `EP-01` portfolio and catalog contracts and `EP-02` plan reading. Managed coordination
instructions in `EP-04` wait for command names, preflight, blockers, and output fields to stabilize.

### F) Tasks Checklist

- [x] Implement portfolio-v1 and portfolio-v2 loading with normalized role, identity, source, and compatibility records in `internal/portfolio/config.go`.
- [x] Implement exact default-branch membership, home self-membership, candidate, mismatch, and exclusion decisions in `internal/portfolio/membership.go`.
- [x] Define provider-neutral repository, ref, file, merge-base, change, and enrichment interfaces in `internal/portfolio/source.go`.
- [x] Implement scanner-owned bare-mirror creation, authenticated fetch, ref prune, remote-only ref namespace, and atomic refresh in `internal/portfolio/mirror.go`.
- [x] Implement bounded direct-argument remote discovery and bare-object content reads in `internal/portfolio/localgit.go`.
- [x] Implement `gh auth status` preflight, complete paginated read-only discovery, truncation detection, three-attempt transient retry, rate-limit blockers, and redacted diagnostics in `internal/portfolio/github.go`.
- [x] Implement default-baseline collection, unmerged-branch enumeration, merge-base filtering, observation assembly, conflict detection, staleness, and retirement in `internal/portfolio/scanner.go`.
- [x] Implement strategic-overlay validation and immutable preservation in `internal/portfolio/overlay.go`.
- [x] Implement atomic complete-cache replacement and incomplete-scan preservation in `internal/portfolio/store.go`.
- [x] Implement repeatable `--github-owner` and `--local-root` argument parsing without changing existing single-value command behavior in `internal/commands/util.go`.
- [x] Implement role-specific `portfolio configure` inputs, deterministic placement, matching-value idempotency, and conflicting-identity blockers in `internal/commands/portfolio.go`.
- [x] Implement `portfolio scan` text and JSON results in `internal/commands/portfolio.go`.
- [x] Extend `plans validate` and `plans inventory` with `--remote-overlays` through the shared portfolio source engine in `internal/commands/plans.go`.
- [x] Extend `internal/cli/cli.go` with `portfolio` group help, subcommand help, argument errors, and exit-code routing.
- [x] Add fake-Git and fake-`gh` fixtures for bare-mirror refresh, stale local refs, local-only branches, default branches, multiple unmerged branches, pagination, truncation, retry exhaustion, rate limiting, PR enrichment, merged refs, deleted refs, auth failure, access failure, and command cancellation.
- [x] Add configuration and false-enrollment fixtures for repeatable sources, matching reruns, conflicting identities, home self-membership, branch-only config, missing Kit markers, missing stable identity, mismatched home identity, and provider-visible unrelated repositories.
- [x] Add cache-integrity tests for success, incomplete scan, interrupted write, strategic-overlay preservation, and deterministic ordering.
- [x] Run `go test ./internal/portfolio ./internal/plancatalog ./internal/commands ./internal/cli` and record the result in the execution log.
- [x] Run local multi-repository fixture scans on macOS and Windows through `.github/workflows/validate.yml`.

### G) Implementation Notes

- Use `exec.CommandContext` with argument slices. Do not invoke a shell.
- Create mirrors only below `.codeheart/local/portfolio/git/`, disable repository hooks for every
  Git subprocess, redact credential-bearing remote material, and use object reads such as
  `git show <ref>:<path>`. A developer checkout must never become the scan source or execution
  working directory.
- A configured local root supplies repository discovery and remote resolution only. Its worktree,
  index, local heads, and unpushed commits are excluded from the factual portfolio catalog.
- Determine unmerged status relative to the repository's configured default branch. Treat missing
  merge-base evidence as an explicit branch error.
- A materially changed canonical plan means added content or changed canonical plan bytes. Pure
  path deletion retires a branch-local observation and does not invent a plan tombstone in V1.
- Preserve last complete cache on auth, permission, network, parse, or required-member failure.
  The current command still reports the failed attempt distinctly.
- Treat pagination, provider truncation, retry exhaustion, rate-limit exhaustion, mirror-refresh
  failure, and missing remote provenance as incomplete-scan blockers.
- Candidate output may include repository identity and reason only; do not ingest plan content.
- Four provider operations is the initial concurrency ceiling. Capture metrics before tuning.

### H) Open Questions

- `OQ-1` defers a hard performance SLO; bounded concurrency and duration reporting are fixed.
- `OQ-2` controls optional live-provider evidence; deterministic adapter tests remain mandatory.

### Validation Tasks

- [x] Prove a large inherited baseline appears once while each changed branch plan appears as one independent observation.
- [x] Prove provider visibility alone never creates catalog membership.
- [x] Prove the coordination home appears once through self-membership and needs no self-registration source.
- [x] Prove stale local refs, local heads, worktree changes, and unpushed commits never appear as current remote observations.
- [x] Prove pagination, truncation, retry exhaustion, and rate-limit exhaustion prevent complete-cache replacement.
- [x] Prove a failed required-member scan preserves the prior complete cache byte-for-byte.
- [x] Prove adapter output and artifacts contain no secret-bearing environment value.
- [x] Prove strategic overlay bytes remain unchanged across successful and failed scans.

## EP-04 - Managed Planning, Coordination, And Publication UX

### A) Epic ID, Title, And Outcome

`EP-04 - Managed Planning, Coordination, And Publication UX`

Outcome: fresh Kit-enabled repositories and low-context agents can author semantic plans,
configure a coordination home, refresh before analysis, migrate safely, and publish activated plan
checkpoints using managed instructions that exactly match the implemented commands and authority
boundaries.

### B) Scope

In scope:

- metadata and portfolio reference docs;
- lifecycle and register compatibility updates;
- discovery, implementation planning, execution, review, and register-transition hooks;
- coordination-home setup, refresh, and migration runbooks;
- plan-only activation publication authority and stop conditions;
- route selection and tooling readiness hooks;
- stable register, portfolio README, and strategic-overlay scaffolds;
- AGENTS and repo README routes;
- component/profile inventories and packaged mirrors;
- placement and consumer-impact updates; and
- low-context route tests.

Out of scope:

- private portfolio configuration;
- repository-specific strategy content;
- automatic external mutations during setup;
- release publication; and
- mandatory schedule guidance.

### C) Files Touched

- `components/planning-workflows/component.yaml`
- `components/planning-workflows/managed/README.md`
- `components/planning-workflows/managed/reference/planning-document-lifecycle.md`
- `components/planning-workflows/managed/reference/plan-register-format.md`
- `components/planning-workflows/managed/reference/plan-catalog-format.md`
- `components/planning-workflows/managed/reference/portfolio-coordination-format.md`
- `components/planning-workflows/managed/runbooks/discovery-workflow.md`
- `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
- `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
- `components/planning-workflows/managed/runbooks/review-planning-document.md`
- `components/planning-workflows/managed/runbooks/maintain-plan-register.md`
- `components/planning-workflows/managed/runbooks/configure-portfolio-coordination.md`
- `components/planning-workflows/managed/runbooks/refresh-portfolio-catalog.md`
- `components/planning-workflows/managed/runbooks/migrate-plan-catalog.md`
- `components/planning-workflows/scaffolds/plan-register.md`
- `components/planning-workflows/scaffolds/portfolio-README.md`
- `components/planning-workflows/scaffolds/portfolio-strategic-overlay.yaml`
- `components/agent-interface/component.yaml`
- `components/agent-interface/managed/README.md`
- `components/agent-interface/managed/kit-readme.md`
- `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`
- `profiles/standard.yaml`
- `templates/agents/AGENTS.managed-block.md`
- `templates/consumer-docs/repo/README.md`
- `docs/repo/reference/placement-contract.md`
- `docs/repo/runbooks/change-operating-kit.md`
- `tests/test_packaging_resources.py`
- `tests/test_routing.py`
- `src/codeheart_operating_kit/resources/`

### D) Acceptance Criteria And Size

Size: `XL`.

Acceptance criteria:

- installed guidance alone explains member and coordination-home concepts, semantic plans,
  catalog modes, source-derived facts, strategic overlays, and freshness limits;
- new references publish the exact metadata, ID, family, config, membership, observation,
  ownership, and compatibility contracts;
- every new runbook declares audience, intention, source, inputs, preconditions, execution lane,
  approval boundary, stop conditions, evidence, validation, and recovery;
- the hybrid configuration runbook asks one decision at a time for role, repository ID, home ID,
  committed GitHub scope, and machine-local roots, then shows the exact command and change plan;
- ordinary discovery and implementation drafting create metadata in mixed and canonical mode;
- legacy mode retains current register behavior;
- activation guidance treats the request as plan-only branch creation/use, commit, and normal push
  authority without a second prompt;
- activation remains an L1 managed workflow with no publication CLI, and its evidence contract is a
  fresh-agent Git transcript plus static authority-contract checks;
- negative boundaries explicitly exclude unrelated files, unauthorized code, PR, merge, release,
  force-push, deletion, destructive actions, ambiguous targets, and rejected pushes;
- execution may begin on an active published work branch without merge;
- coordination analysis refreshes first and cannot claim current completeness after a failed scan;
- a fresh repository can preview and configure a coordination home without hidden documentation;
- fresh coordination-home scaffolds create only when absent and remain repo-owned;
- machine-specific sources and catalog cache stay below `.codeheart/local/`;
- existing consumer register and coordination-pending files remain preserved by sync while fresh
  installations stop creating `coordination-sync-pending.md`;
- no managed file contains real private repository topology, credentials, customer data, or
  Codeheart strategic records;
- component and profile declarations include every new managed and scaffold surface; and
- packaged Python fallback resources exactly mirror source resources.

### E) Dependencies And Critical-Path Notes

Depends on stabilized contracts from `EP-01`, command and blocker behavior from `EP-02`, and scan
behavior from `EP-03`. Doctrine must not speculate beyond implemented behavior.

### F) Tasks Checklist

- [x] Create `plan-catalog-format.md` with exact metadata markers, fields, IDs, kinds, relations, families, aliases, modes, observations, and compatibility examples.
- [x] Create `portfolio-coordination-format.md` with v1/v2 config, authorized scope, membership predicates, source ownership, candidates, cache, overlay, and freshness rules.
- [x] Update `planning-document-lifecycle.md` with separate record kinds, metadata authority, catalog modes, Git visibility, verification, and content chronology.
- [x] Update `plan-register-format.md` with stable-entry, legacy-evidence, generated-view, and no-manual-append behavior.
- [x] Update the discovery, implementation-planning, execution, and review runbooks with mode-aware metadata authoring and validation gates.
- [x] Add the activation and material-active-update plan-checkpoint publication contract to implementation-planning and execution runbooks.
- [x] Add explicit publication exclusions and stop conditions for ambiguous scope, auth failure, policy rejection, normal-push rejection, overlapping dirty work, and broader external actions.
- [x] Define L1 fresh-agent and L3 command validation tiers, non-secret evidence fields, blocker fields, and the no-L2-script promotion boundary in the new managed references and runbooks.
- [x] Replace manual catalog authority in `maintain-plan-register.md` with local view, legacy compatibility, cutover, and coordination routing.
- [x] Create `configure-portfolio-coordination.md` as one hybrid member/home setup recipe with paced identity and scope questions, exact flags, preview, user review, approved write, validation, and first scan.
- [x] Create `refresh-portfolio-catalog.md` as an agent-facing read recipe with preflight, freshness, completeness, failure disclosure, and strategic-overlay preservation.
- [x] Create `migrate-plan-catalog.md` as an agent-facing inventory, semantic review, dry-run, guarded apply, reevaluation, cutover, and closure recipe.
- [x] Update planning-workflows, agent-interface, fallback kit inventory routes, and `operation-routing-and-dispatch.md` with local-plan and portfolio selections.
- [x] Modify the fresh `plan-register.md` scaffold into a compact stable entry point.
- [x] Remove `coordination-sync-pending.md` from fresh component and profile generated surfaces while preserving existing consumer-owned copies and compatibility reads.
- [x] Create repo-owned `portfolio-README.md` and `portfolio-strategic-overlay.yaml` scaffolds with public-safe placeholders.
- [x] Update `components/planning-workflows/component.yaml`, `components/agent-interface/component.yaml`, `profiles/standard.yaml`, `templates/agents/AGENTS.managed-block.md`, and `templates/consumer-docs/repo/README.md` with exact new surfaces.
- [x] Update `docs/repo/reference/placement-contract.md` and `docs/repo/runbooks/change-operating-kit.md` with ownership and validation gates.
- [x] Synchronize every changed embedded source file into `src/codeheart_operating_kit/resources/`.
- [x] Extend `tests/test_packaging_resources.py` with every new managed and scaffold resource.
- [x] Extend `tests/test_routing.py` with low-context member authoring, coordination setup, refresh, migration, activation, and failed-publication scenarios.
- [x] Run focused managed-doc, routing, packaging, markdown-header, and public-core tests and record the result in the execution log.

### G) Implementation Notes

- Keep runbook prose compact. Put field-level detail in references and command behavior in the CLI.
- `configure-portfolio-coordination.md` is hybrid because it collects user-facing member/home
  identity and scope decisions before agent-only config, local-source, and scaffold writes.
  `refresh-portfolio-catalog.md` is read-only external access plus rebuildable local cache.
  `migrate-plan-catalog.md` is a write recipe requiring reviewed evidence and explicit apply.
- The activation rule is not generic push authority. It applies only to the planning checkpoint
  caused by activation or a user-requested material active-plan update.
- Activation publication remains doctrine-driven. Do not add an activation, commit, push, PR, or
  merge command to the Operating Kit CLI in V1.
- Existing portfolio-v1 pending-sync files remain readable compatibility evidence, but the scaffold
  is not installed in fresh repositories once portfolio-v2 scanning is available.
- Scaffolds provide generic structure. Actual source scopes, repository identities, strategy,
  themes, priorities, and analysis remain repo-owned.
- Update component versions and profile version according to normal producer practice during
  implementation; do not choose a release version.

### H) Open Questions

None that block managed UX. Live scope and later repository rollout remain operational follow-ups,
not doctrine gaps.

### Validation Tasks

- [x] Prove a fresh initialized fixture contains every new managed route and repo-owned scaffold without `coordination-sync-pending.md`.
- [x] Prove an existing consumer-owned `coordination-sync-pending.md` survives sync and remains readable as compatibility evidence.
- [x] Prove a low-context agent selects refresh before current portfolio analysis and reports incomplete scans.
- [x] Prove activation wording grants one plan-only normal push and refuses every named broader action.
- [x] Prove source and packaged resource bytes match for every changed embedded file.
- [x] Prove public-core validation rejects private topology and credential examples in managed content.

## EP-05 - Producer Semantic Migration And Catalog Cutover

### A) Epic ID, Title, And Outcome

`EP-05 - Producer Semantic Migration And Catalog Cutover`

Outcome: every formal plan in this Operating Kit producer receives reviewed canonical metadata,
legacy aliases and evidence are reconciled, actively edited plans are handled by their branch
owners, manual register appends stop, and the default branch plus active overlays pass canonical
coverage without pausing unrelated work.

### B) Scope

In scope:

- full producer plan inventory;
- semantic review of discovery, implementation, and family records;
- public-safe purpose, family, product, capability, theme, relation, and alias decisions;
- source revision and hash preconditions;
- active-branch detection and owner assignment;
- explicit deferral and reevaluation;
- metadata-only chronology preservation;
- legacy register evidence reconciliation;
- `legacy` to `mixed` to `canonical` mode transitions;
- frozen-register notice and cessation of manual appends;
- final default and active-overlay coverage; and
- committed plan-scoped migration evidence.

Out of scope:

- other repository writes;
- deletion of historical register entries;
- forced migration of unknown consumers;
- migration of execution logs as plan records; and
- rewriting plan content merely to modernize prose.

### C) Files Touched

- `docs/repo/plans/**/*_discovery_doc.md`
- `docs/repo/plans/**/*_implementation_doc.md`
- qualifying `docs/repo/plans/**/README.md` family records
- `docs/repo/plans/plan-register.md`
- `.codeheart/kit.config.yaml` as the tracked shared producer catalog-mode and portfolio-identity authority
- `docs/repo/plans/semantic-plan-catalog-coordination/attachments/producer-plan-migration-inventory.json`
- `docs/repo/plans/semantic-plan-catalog-coordination/attachments/producer-plan-migration-ledger.yaml`
- `docs/repo/plans/semantic-plan-catalog-coordination/semantic-plan-catalog-coordination_execution_log.md`

### D) Acceptance Criteria And Size

Size: `XXL`.

Acceptance criteria:

- inventory enumerates every canonical discovery and implementation document and every qualifying
  family record in the producer;
- each record receives a semantic human review rather than a converted register row;
- every sequential register ID linked to a canonical record is preserved as a legacy alias;
- combined discovery-plus-implementation register entries are separated into their actual records;
- family, relation, taxonomy, lifecycle, session, and coordination evidence decisions are recorded
  with confidence and ambiguity;
- current branch/worktree ownership is checked before any file is assigned to migration;
- changed hashes, overlapping dirty files, and active ownership cause a recorded skip;
- branch owners can supply metadata through their ongoing work without a global freeze;
- metadata-only insertion leaves the historical `Last updated` header byte-identical;
- ambiguous optional classifications remain omitted rather than invented;
- producer mode changes to `mixed` only after compatibility tooling and new authoring guidance are
  installed in source;
- register appends cease at mixed cutover and the existing register receives a frozen-authority
  notice without deleting its history;
- a final inventory catches plans created or changed during migration;
- producer mode changes to `canonical` only after default-branch and accessible active-overlay
  coverage passes through `plans validate --remote-overlays`;
- the migration apply is idempotent;
- public-core validation passes on every migrated metadata value; and
- no other repository is modified.

### E) Dependencies And Critical-Path Notes

Depends on executable migration mechanics from `EP-02` and final authoring/cutover doctrine from
`EP-04`. The active implementation plan itself is branch-owned and must be migrated by its owning
branch after source hashes stabilize.

### F) Tasks Checklist

- [x] Create the plan-scoped execution log before migration evidence is collected.
- [x] Run `codeheart-operating-kit plans inventory --remote-overlays --output docs/repo/plans/semantic-plan-catalog-coordination/attachments/producer-plan-migration-inventory.json .` from the producer work branch.
- [x] Reconcile every inventory record against `docs/repo/plans/plan-register.md`, sibling documents, execution evidence, and current lifecycle headers.
- [x] Assign one stable semantic ID, kind, purpose, first-cataloged time, catalog-update time, and legacy alias set to every formal record.
- [x] Record evidence-backed family, product, capability, theme, and relation values in `producer-plan-migration-ledger.yaml`.
- [x] Record omitted optional classifications, ambiguity, confidence, branch owner, conflicts, and deferrals in `producer-plan-migration-ledger.yaml`.
- [x] Validate the completed ledger against `schemas/plan-migration-ledger.schema.json`.
- [x] Recheck source revisions, byte hashes, dirty overlap, and active-branch ownership immediately before the dry run.
- [x] Run `codeheart-operating-kit plans migrate --ledger docs/repo/plans/semantic-plan-catalog-coordination/attachments/producer-plan-migration-ledger.yaml --dry-run .`.
- [x] Review every planned write, skip, blocker, preserved header value, and unrelated-file exclusion from the dry-run result.
- [x] Inspect the existing `.codeheart/kit.config.yaml` path, preserve untracked and dirty user content, and block on conflicting portfolio identity before adoption.
- [x] Adopt `.codeheart/kit.config.yaml` as shared tracked producer configuration with the approved stable repository identity and `plan_catalog_mode: mixed`.
- [x] Apply the reviewed ledger with `codeheart-operating-kit plans migrate --ledger docs/repo/plans/semantic-plan-catalog-coordination/attachments/producer-plan-migration-ledger.yaml --yes .`.
- [x] Add the frozen-authority and generated-listing notice to `docs/repo/plans/plan-register.md` without deleting historical entries.
- [x] Re-run remote-overlay inventory and semantically review every record created, changed, skipped, and deferred during the migration window.
- [x] Assign contested plan migration to each current branch owner and record the resulting source evidence in the ledger.
- [x] Apply reviewed reconciliation entries with the same guarded migration command.
- [x] Run `codeheart-operating-kit plans validate --remote-overlays .` in mixed mode and resolve every canonical-record blocker.
- [x] Change producer catalog mode from `mixed` to `canonical` after default and active-overlay coverage passes.
- [x] Run `codeheart-operating-kit plans validate --remote-overlays .` and `codeheart-operating-kit plans list --format json .` in canonical mode.
- [x] Re-run the migration apply and record the zero-change idempotency result.
- [x] Run `python3 scripts/validate-public-core.py` across the migrated repository and record the result.

### G) Implementation Notes

- This epic is intentionally semantic and human/agent-reviewed. The CLI owns enumeration and safe
  application, not strategic inference.
- Do not normalize plan prose, whitespace, lifecycle, or historical dates as a side effect of
  inserting metadata.
- Use the stable repository ID in tracked shared `.codeheart/kit.config.yaml` for every semantic
  namespace. A missing, untracked, dirty-conflicting, or ambiguous config blocks migration rather
  than deriving an ID from the current folder name.
- Treat register-only notes as evidence. Place only approved canonical fields in metadata.
- An active branch may merge during migration. Re-run inventory and evaluate the latest default
  branch rather than applying a stale ledger entry.
- The producer's ignored `.codeheart/kit/` remains non-authoritative. The sibling tracked shared
  `.codeheart/kit.config.yaml` is repository configuration, not managed-kit source. Existing user
  content must be inspected and preserved before it becomes tracked authority.

### H) Open Questions

- `OQ-4` does not block migration; existing register history remains in place.
- Semantic ambiguity on optional metadata does not block a record's required core. Ambiguity in
  record identity, kind, canonical path, or source ownership blocks that record's migration.

### Validation Tasks

- [x] Prove the inventory count equals canonical discovery, implementation, and qualifying family record counts.
- [x] Prove every migrated legacy register ID resolves to exactly one canonical record alias, with each exception recorded as a documented ambiguity.
- [x] Prove all metadata-only changes preserve pre-migration content-update header values.
- [x] Prove `plans validate --remote-overlays` covers the default branch and every accessible active plan-changing branch before canonical cutover.
- [x] Prove a second migration apply produces zero canonical plan changes.

## EP-06 - Integrated Validation And Release-Candidate Handoff

### A) Epic ID, Title, And Outcome

`EP-06 - Integrated Validation And Release-Candidate Handoff`

Outcome: automated and low-context evidence proves the end-to-end capability on supported
platforms, consumer impact is explicit, the validated source is ready for `EP-07` versioning and
release, and additional repositories receive safe migration handoffs without being modified.

### B) Scope

In scope:

- repository-wide Go and Python validation;
- macOS and Windows fixture execution;
- low-context plan reference, setup, refresh, migration, and doctrine-driven activation probes;
- privacy, token, branch-code, path, and authority negative tests;
- mixed/canonical and portfolio-v1/v2 compatibility;
- performance measurement without a premature SLO;
- consumer-impact record and migration/release notes;
- source/packaged mirror and manifest proof;
- optional live read-only GitHub probe under configured scope; and
- dry-run rollout handoff design for additional repositories.

Out of scope:

- release tag and publication;
- release-version selection and version-surface mutation;
- pull request creation and merge;
- consumer upgrade or sync;
- persistent scheduling;
- cross-repository migration writes; and
- claims of live-provider proof when no authorized environment exists.

### C) Files Touched

- `.github/workflows/validate.yml`
- `tests/fixtures/plans/`
- `tests/fixtures/portfolio/`
- `tests/test_json_schemas.py`
- `tests/test_markdown_headers.py`
- `tests/test_packaging_resources.py`
- `tests/test_public_core.py`
- `tests/test_routing.py`
- `internal/plancatalog/*_test.go`
- `internal/portfolio/*_test.go`
- `internal/commands/commands_test.go`
- `internal/cli/cli_test.go`
- `docs/repo/README.md`
- `docs/repo/plans/README.md`
- `docs/repo/plans/plan-register.md`
- `docs/repo/plans/semantic-plan-catalog-coordination/semantic-plan-catalog-coordination_execution_log.md`
- release-note and impact surfaces selected by `docs/repo/runbooks/change-operating-kit.md`

### D) Acceptance Criteria And Size

Size: `XL`.

Acceptance criteria:

- all Go tests pass;
- all repository Python tests pass;
- JSON schemas, Markdown headers, public-core rules, component declarations, profile surfaces, and
  packaged resources validate;
- macOS and Windows execute metadata, ID, local Git branch-overlay, migration, and CLI tests;
- low-context agents identify plans by title, kind, semantic ID, family, and path without a legacy
  register number;
- low-context coordination setup succeeds from installed managed guidance alone;
- low-context strategic analysis refreshes first, distinguishes current from last-known complete,
  and discloses inaccessible members;
- a fresh-agent activation probe follows managed instructions, commits and normally pushes only the
  canonical planning checkpoint to an isolated local bare remote, and asks no second approval;
- static authority-contract tests plus fresh-agent negative probes refuse unrelated files,
  ambiguous branch/repository, unauthorized code, PR, merge, release, force-push, deletion,
  destructive action, and failed auth;
- adversarial branch content is never executed;
- logs and artifacts contain no captured credential or secret-bearing environment value;
- representative local and fake-GitHub benchmarks record repository, branch, plan, API-call,
  duration, and concurrency measurements;
- consumer impact separately declares instruction-only, validator-only, backwards-compatible
  scaffold, consumer migration, generated/local path, placement, and security/safety effects;
- migration guidance states that existing work continues, compatibility precedes cutover, and
  repository-owned activation is required for writes;
- release-candidate readiness identifies the exact validated source state, required version
  surfaces, remaining platform evidence, and `EP-07` publication gate; and
- no additional repository receives a write.

### E) Dependencies And Critical-Path Notes

Depends on all earlier epics. This epic closes source implementation and hands one validated
candidate to `EP-07`. Public publication remains subject to the producer release runbook and its
explicit user-authorization gate. Additional repository migration remains a separate
owner-repository plan or explicitly scoped portfolio rollout.

### F) Tasks Checklist

- [x] Extend `.github/workflows/validate.yml` with supported-platform plan metadata, migration, local Git overlay, and grouped CLI tests.
- [x] Add an isolated repository and local bare-remote fixture for recorded fresh-agent activation publication probes.
- [x] Add static authority-contract assertions and fresh-agent refusal scenarios for unrelated staged files, ambiguous branch scope, rejected normal push, PR creation, merge, release, force-push, branch deletion, and destructive commands.
- [x] Add adversarial branch fixtures containing executable files, hooks, malformed metadata, path traversal names, oversized metadata, and secret-like text.
- [x] Add low-context routing probes for semantic plan reference, member authoring, coordination-home setup, current analysis, failed refresh, migration skip, and activation publication.
- [x] Benchmark representative local and fake-GitHub portfolios and record repository, branch, plan, API-call, duration, and concurrency evidence.
- [x] Record configured live GitHub probe availability and redacted evidence status in the execution log.
- [x] Classify every changed surface under `docs/repo/reference/consumer-impact-classification.md`.
- [x] Draft compatibility migration, adoption, consumer-impact, and release notes for the final implemented behavior.
- [x] Update `docs/repo/README.md` and `docs/repo/plans/README.md`, and preserve the frozen current register compatibility entry byte-for-byte.
- [x] Run `gofmt -w` on changed Go files and verify the resulting diff.
- [x] Run `go test ./...` and record the result in the execution log.
- [x] Run the repository Python suite through the available isolated pytest runtime and record the result in the execution log.
- [x] Run `python3 scripts/validate-json-schemas.py` and record the result in the execution log.
- [x] Run `python3 scripts/validate-markdown-headers.py` and record the result in the execution log.
- [x] Run `python3 scripts/validate-public-core.py` and record the result in the execution log.
- [x] Run `git diff --check` and record the result in the execution log.
- [x] Verify source and packaged resource identity for every changed embedded file.
- [x] Review the implementation against every frozen discovery decision, capability must-cover item, explicit exclusion, and success-evidence statement.
- [x] Produce repository-scoped migration handoff summaries without changing any additional repository.
- [x] Record the validated source revision and hand it to EP-07 with release version, tag, public assets, consumer installations, and other repositories unchanged.

### G) Implementation Notes

- The live GitHub task is evidence-sensitive, not capability-sensitive. It must not cause the suite
  to fake success when no authorized scope exists.
- Use a fake `gh` executable for deterministic provider behavior and secret-capture assertions.
- The isolated activation fixture is executed by a fresh-agent probe following managed runbooks; it
  proves selected-file, Git-transcript, approval, refusal, and visibility behavior without adding a
  publication CLI or creating external state.
- Cross-repository handoffs may include repository ID, observed coverage, missing metadata, branch
  ownership, and recommended mode transition. They are not migration ledgers until the owning
  repository semantically reviews and activates them.
- Do not update component, profile, packaged-resource, installer, or release versions merely to
  make source tests pass. `EP-07` selects and applies one coherent version after source completion.

### H) Open Questions

- `OQ-1`: benchmark evidence may recommend a later performance plan without blocking V1.
- `OQ-2`: live provider proof is recorded only under real configured authorization.
- `OQ-3`: additional repository order remains a portfolio decision after source completion.

### Validation Tasks

- [x] Prove all mandatory local and cross-platform gates pass from a clean fixture state.
- [x] Prove activation publication through static managed-contract tests and a recorded fresh-agent local-remote probe without a publication CLI.
- [x] Prove source completion leaves release tags, public assets, consumer sync, and other-repository writes untouched before EP-07.
- [x] Prove every discovery must-cover requirement maps to one designated passing evidence item.
- [x] Prove migration and rollout documentation preserves repository ownership and ordinary parallel work.

## EP-07 - Version Bump, Reproducible Release, And Public Verification

### A) Epic ID, Title, And Outcome

`EP-07 - Version Bump, Reproducible Release, And Public Verification`

Outcome: the completed semantic plan-catalog capability receives the next unused patch version,
passes the full producer release contract on macOS and real Windows, is published from the exact
validated commit after explicit authorization, and is verified through its public assets.

### B) Scope

In scope:

- target-version resolution against the latest public tag at epic entry;
- one coherent version bump across source, packaged resources, installers, fixtures, and binary
  identity;
- versioned release notes covering behavior, consumer impact, compatibility, migration, adoption,
  and rollout prerequisites;
- full producer validation after the version bump;
- two independent deterministic macOS universal and Windows x64 pack builds;
- external catalog, archive, pack-manifest, payload, content, binary, and sidecar digest proof;
- isolated macOS fresh-install, upgrade, and failure-path validation;
- real-Windows fresh-install, upgrade, and failure-path validation;
- signing/notarization state and intended-distribution-boundary decision;
- explicit release-execution authorization for the resolved version, presented candidate source
  tree, remote validation, distribution boundary, and conditional publication;
- release-candidate commit and branch push for real-Windows validation;
- validated-commit annotated tag creation, public release publication, and release evidence;
  and
- post-publication download, update-check, install, and upgrade verification against public assets.

Out of scope:

- pull request creation and merge;
- upgrading a real coordination home or member repository;
- portfolio-wide migration execution;
- writes in any additional repository;
- persistent coordination scheduling; and
- new signing infrastructure beyond evidence required for the approved distribution boundary.

### C) Files Touched

- `manifest.yaml`
- `internal/version/version.go`
- `pyproject.toml`
- `src/codeheart_operating_kit/__init__.py`
- `src/codeheart_operating_kit/resources/manifest.yaml`
- `components/planning-workflows/component.yaml`
- `src/codeheart_operating_kit/resources/components/planning-workflows/component.yaml`
- `components/agent-interface/component.yaml`
- `src/codeheart_operating_kit/resources/components/agent-interface/component.yaml`
- `profiles/standard.yaml`
- `src/codeheart_operating_kit/resources/profiles/standard.yaml`
- `release-notes.md`
- `bootstrap.md`
- `install.sh`
- `install.ps1`
- `.github/workflows/validate.yml`
- `tests/fixtures/release-manifest.json`
- `tests/fixtures/release-candidate/release-candidate-manifest.json`
- `tests/test_release_assets.py`
- `tests/test_install_metadata.py`
- `internal/manifest/manifest_test.go`
- `internal/release/release_test.go`
- `internal/commands/commands_test.go`
- `tests/test_go_cli_parity.py`
- `docs/repo/README.md`
- `docs/repo/plans/README.md`
- `docs/repo/plans/plan-register.md`
- `docs/repo/plans/semantic-plan-catalog-coordination/semantic-plan-catalog-coordination_execution_log.md`
- `dist/` as ignored generated release-candidate output

### D) Acceptance Criteria And Size

Size: `XL`.

Acceptance criteria:

- the selected target is the next unused patch after the latest validated public tag at epic
  entry; the current recommended target is `v0.1.24` only while `v0.1.23` remains latest;
- every authoritative and packaged version surface agrees with the selected version, with no stale
  previous-version value outside intentional historical release notes and compatibility fixtures;
- release notes describe the new `plans` and `portfolio` command groups, metadata and config
  contracts, compatibility modes, producer migration, pending-sync scaffold retirement, activation
  publication boundary, consumer impact, and repository-owned rollout prerequisites;
- all public-core, Markdown, JSON Schema, content-identity, Go, Python compatibility, installer,
  and release-contract validation passes after the version bump;
- two isolated builds produce byte-identical macOS universal and Windows x64 packs;
- every pack contains the expected binary, bootstrap and installer files, release notes,
  `INSTALL.md`, content manifest, pack manifest, checksums, and no Python package payload;
- the external catalog is generated after the packs and proves the complete
  catalog-to-archive-to-pack-to-payload-to-content-to-binary identity chain;
- isolated macOS and real-Windows fresh-install, upgrade dry-run, upgrade apply, and injected
  failure paths pass while failed replacement, reconciliation, and post-check preserve the prior
  installation;
- sidecar SHA-256 checksums exist and verify for every public asset;
- signing/notarization state and the approved audience boundary are recorded without presenting an
  unsigned candidate as broadly trusted distribution;
- explicit user release-execution authority covers the selected version, presented candidate
  source tree, remote validation actions, distribution boundary, and conditional publication
  before the candidate commit and branch are pushed;
- the tag target is rechecked against the exact commit that passed all local and real-Windows
  validation before public publication;
- the public tag and release point to the validated commit and contain packs, external catalog,
  installers, bootstrap, notes, and checksums;
- public downloads reproduce the recorded digests and pass isolated update-check, fresh-install,
  and upgrade smoke tests; and
- the execution log records commands, CI evidence, URLs, digests, signing state, authorization,
  residual risk, and the fact that no consumer repository was upgraded.

### E) Dependencies And Critical-Path Notes

Depends on `EP-06` and the exact source revision it validates. Version selection must occur after a
fresh public-tag preflight because parallel work may have consumed the currently recommended patch
number. Any source change after release validation invalidates the candidate and returns execution
to the relevant validation gate. Remote candidate push, workflow dispatch, tag creation, and
publication remain blocked until the `OQ-5` preflight, `OQ-6` distribution decision, and explicit
release-runbook authorization are resolved.

### F) Tasks Checklist

- [x] Re-read `docs/repo/runbooks/release-operating-kit.md` and record the final consumer-impact classification, migration requirement, and release scope.
- [x] Resolve the next unused patch from current public tags and record `v0.1.24` only while `v0.1.23` remains latest.
- [x] Update root, binary, legacy-compatibility, component, profile, packaged-resource, installer, fixture, and workflow version surfaces to the selected version.
- [x] Add versioned release notes covering command UX, contracts, compatibility, producer migration, consumer adoption, rollout prerequisites, safety boundaries, and known residual risk.
- [x] Regenerate content identity and verify every source-to-packaged-resource mirror after the version bump.
- [x] Run public-core, Markdown, JSON Schema, content-identity, Go, Python compatibility, installer, and release-contract validation from the versioned source tree.
- [x] Build macOS universal and Windows x64 release packs twice in separate output directories and prove byte-identical archives.
- [x] Verify required pack members, Python-payload exclusion, external-catalog sequencing, and the complete catalog-to-binary digest chain.
- [x] Run isolated macOS fresh-install, upgrade dry-run, upgrade apply, replacement failure, reconciliation failure, and post-check failure paths.
- [x] Generate and verify sidecar SHA-256 checksum files for every proposed public asset.
- [x] Record signing and notarization evidence plus the intended distribution boundary in the release-candidate summary.
- [x] Present the resolved version, candidate source tree, local evidence, remote validation actions, distribution boundary, conditional publication action, and rollback boundary for explicit user authorization.
- [ ] Create and push the release-candidate commit on the unambiguous branch under the recorded release-execution authority.
- [ ] Dispatch the real-Windows fresh-install, upgrade dry-run, upgrade apply, replacement failure, reconciliation failure, and post-check failure matrix and retain CI evidence.
- [ ] Verify the release target still equals the validated commit and rerun invalidated gates after every intervening source change.
- [ ] Create and push the annotated version tag at the exact locally and remotely validated commit under the confirmed publication authority.
- [ ] Publish the packs, external catalog, bootstrap, installers, release notes, and checksum sidecars as one public release.
- [ ] Verify public asset URLs and digests, then run isolated public update-check, fresh-install, and upgrade smoke tests.
- [ ] Record the release URL, tag, commit, assets, digests, macOS evidence, Windows evidence, signing state, authority, and residual risk in the execution log.
- [ ] Update the producer plan indexes and compatibility register to the released status while leaving coordination homes, member repositories, and additional checkouts unchanged.

### G) Implementation Notes

- Treat `v0.1.24` as a planning recommendation, not reserved global state. The tag preflight owns
  the final value.
- Apply one version consistently; do not bump individual components early merely because their
  source changes land in an earlier epic.
- Build outputs remain ignored release artifacts. The external catalog is generated only after
  final pack bytes exist and contains the public archive URLs and digests that must not enter the
  embedded manifest.
- A locally successful build is release-candidate evidence only. Real-Windows validation and the
  distribution-boundary decision remain publication gates.
- An activation request authorizes release execution only when it explicitly covers the dynamic
  target-version rule, candidate commit and branch push, remote validation workflow, distribution
  boundary, and conditional tag/publication after every gate passes. Otherwise execution pauses
  before remote release-candidate actions.
- Public release makes rollout possible; it does not prove that a coordination home or member has
  upgraded. The later rollout begins with a fresh coordination-home portfolio scan, then upgrades
  and inventories each selected repository before migration.

### H) Open Questions

- `OQ-5`: resolve the target version from live public tags at epic entry.
- `OQ-6`: resolve and record the signing/notarization and intended-audience boundary before
  publication.

### Validation Tasks

- [x] Prove one selected version across source identity, packaged resources, binaries, installers, fixtures, workflows, archives, and external catalog.
- [x] Prove two builds produce byte-identical supported-platform packs and a coherent external-catalog digest chain.
- [ ] Prove isolated macOS and real-Windows install, upgrade, failure, rollback, and recovery behavior from the release candidate.
- [ ] Prove tag and public assets originate from the exact validated commit under recorded publication authority.
- [ ] Prove public downloads match recorded sidecars and succeed in isolated update-check, fresh-install, and upgrade smoke tests.
- [ ] Prove release completion leaves every coordination-home installation, member installation, portfolio config, and member plan unchanged.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Add optional scheduled invocation of `portfolio scan` after measured demand, reusing the same
  command and catalog contract.
- Add webhook, GitHub App, or provider-event triggers only after an explicit reliability and
  credential-lifecycle discovery.
- Reconsider a pre-push intent or announcement layer only when coordination needs visibility before
  the first authorized push.
- Add provider adapters beyond local Git and GitHub-through-`gh` through the existing `Source`
  interface.
- Add content-addressed incremental scan caching after benchmark evidence shows that full on-demand
  refresh is too slow.
- Define a hard scan freshness and duration SLO after representative portfolio evidence exists.
- Decide long-term archival or deletion policy for frozen legacy register bodies.
- Add controlled product, capability, or strategic-theme vocabularies only after actual drift makes
  validation more valuable than flexibility.
- Add catalog UI, database, or remote service views without moving canonical authority away from
  repository plans.
- Add a bounded multi-repository migration plan after repository-specific handoffs are reviewed.
- Upgrade the coordination home to the published Kit, run a fresh complete portfolio scan, select
  a repository-owned adoption batch, upgrade each member before migration, inventory current
  default and remote branches, execute reviewed ledgers, and rescan after each bounded rollout.
- Add signing and notarization infrastructure when the intended distribution audience exceeds the
  explicitly approved unsigned internal/prototype boundary.

## 4.2 Future Considerations

- Semantic ID collisions should be measured before adding tokens. Independent conflicting source
  observations already prevent data loss in V1.
- A durable remote catalog may improve multi-user coordination, but it must preserve source-derived
  facts, coordination-owned overlays, complete-scan semantics, and visible freshness.
- Strategic overlay structure may grow to include decisions, roadmap horizons, and portfolio risks.
  Those fields remain coordinator-owned and should not be pushed back into member metadata without
  repository approval.
- Large portfolios may need incremental branch/ref checkpoints. Such optimization must not change
  membership rules or silently hide unmerged branches.
- Provider APIs may expose incomplete merge relationships for unusual histories. Future adapters
  should retain explicit unknown state rather than infer merged status.
- Family records may eventually span repositories. Cross-repository families belong in the
  coordination-home overlay; member family README files remain repository-local authority.
- Execution logs could expose derived progress summaries later, but they should not become
  independent strategic plans by default.
- A future formal deprecation can remove portfolio-v1 path compatibility after observed consumer
  adoption and explicit migration notice.

# Revision Notes

- 2026-07-31: Created the draft implementation plan from the approved semantic plan catalog and
  branch-aware coordination discovery. Defined six dependency-ordered epics for canonical
  contracts, local compatibility and migration, portfolio scanning, managed UX and activation
  publication, producer semantic cutover, and integrated release-readiness evidence. Resolved the
  first-version register entry point, candidate presentation, stale observation, cache atomicity,
  and catalog-mode defaults while preserving all frozen discovery decisions and repository-owned
  rollout authority.
- 2026-07-31: Applied the planning review recommendations without expanding the six-epic
  architecture: specified role-aware configuration flags and conflict behavior; required a stable
  repository ID and automatic self-membership for coordination homes; added scanner-owned bare
  mirrors, remote-only ref semantics, provider pagination and truncation failure handling, and
  reusable remote-overlay validation; fixed producer config authority; corrected L1/L3 recipe
  maturity and evidence tiers; retired fresh pending-sync scaffolding while preserving existing
  files; and kept activation publication doctrine-driven with static contracts and fresh-agent
  probes rather than a new publication command.
- 2026-07-31: Expanded the draft to seven epics so the producer change ends in an installable
  release instead of a release-readiness handoff. Added dynamic next-patch selection, coherent
  version surfaces, build-twice reproducibility, full catalog-to-binary digest proof, isolated
  macOS and real-Windows release validation, signing-boundary review, an explicit publication gate,
  public release verification, and a strict boundary that leaves coordination-home and member
  rollout for later repository-owned execution.
- 2026-07-31: Activated the implementation plan by explicit user request, created its sibling
  execution log, and selected `codex/semantic-plan-catalog-coordination` as the unambiguous work
  branch. Public release execution remains subject to the `EP-07` release gate.
- 2026-07-31: Completed and independently accepted `EP-01` after two review rounds. The accepted
  foundation includes versioned plan, catalog, migration, local-source, overlay, and nested config
  contracts; semantic identities and families; exact metadata/header parsing; deterministic
  validation; compatibility-mode defaults; automatic home self-membership; and focused fixtures
  proving rename stability, mixed and legacy behavior, same-ID observations, and malformed-input
  blockers.
- 2026-07-31: Completed and independently accepted `EP-02` after twenty-three review rounds. The
  accepted compatibility and migration layer includes deterministic legacy reconciliation and
  views, Git-backed inventory, guarded ledger application, Go 1.25 handle-relative transaction
  containment, rollback and stale recovery, safe explicit inventory replacement, grouped plan
  commands, stable blockers, and race regressions for parent exchange and concurrent-byte
  preservation.
- 2026-07-31: Completed and independently accepted `EP-03` after fourteen review rounds. The
  accepted portfolio layer includes exact configuration and membership, provider-neutral local
  and GitHub discovery, complete default and unmerged-branch observations, immutable strategic
  overlays, atomic cache publication, policy-confined descriptor-bound Git execution, and rooted
  atomic no-replace mirror capture, installation, restoration, and cleanup across supported
  release targets.
- 2026-07-31: Completed and independently accepted `EP-04` after three review rounds. The accepted
  managed UX includes exact semantic-catalog and portfolio references, mode-aware authoring,
  stable register and strategic-overlay scaffolds, refresh-first analysis, semantic migration,
  bounded activation publication, fresh low-context routing, and byte-identical packaged
  resources. Review tightened the generic portfolio scaffold to remain role-neutral and made the
  pristine-placeholder replacement exception consistent and regression-tested throughout setup
  guidance.
