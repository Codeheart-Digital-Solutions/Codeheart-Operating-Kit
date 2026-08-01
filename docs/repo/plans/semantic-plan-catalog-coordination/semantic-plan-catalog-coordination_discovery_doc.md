Last updated: 2026-07-31T09:22:30Z (UTC)
Created: 2026-07-30
Status: draft

# Semantic Plan Catalog And Branch-Aware Coordination Discovery

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
    capabilities:
        - semantic-plan-catalog
    catalog_metadata_updated: "2026-07-31T23:16:16Z"
    first_cataloged: "2026-07-31T23:16:16Z"
    id: codeheart-operating-kit.discovery.semantic-plan-catalog-coordination
    kind: discovery
    legacy_aliases:
        - OK-PR-027
    products:
        - codeheart-operating-kit
    purpose: Define semantic plan identity, generated views, branch-aware portfolio facts, and resilient compatibility migration.
    relations:
        - kind: related
          target: codeheart-operating-kit.implementation.semantic-plan-catalog-coordination
    schema_version: 1
    strategic_themes:
        - cohesive-product-planning
```
<!-- END CODEHEART PLAN METADATA -->

## Overview

This discovery re-evaluates the Operating Kit plan-register model after practical use exposed
identity, merge, comprehension, portfolio-coverage, branch-visibility, and synchronization
problems.

The current model uses manually maintained repeated sections in
`docs/repo/plans/plan-register.md`, sequential identifiers such as `OK-PR-026`, and optional
agent-driven synchronization to a coordination home. That model established useful local and
portfolio planning discipline, but it assumes that agents can safely allocate identifiers and
edit shared register files in sequence. In practice, separate branches can allocate the same next
number, edit the same register location, and begin authorized implementation before their
canonical plan is visible from the default branch. Bare references such as `PR45` also provide too
little context for humans or agents.

The intended evolution is a semantic plan catalog whose canonical facts come from rich structured
metadata in discovery and implementation documents. Local and global overviews are derived rather
than manually duplicated, and plan lifecycle remains separate from Git visibility and catalog
freshness. A coordination home should maintain a complete strategic overview of all formal
planning records in its configured portfolio, not only records already known to cross repository
boundaries. Portfolio-owned strategic families, themes, and relationships should enrich that
factual catalog without taking authority away from repositories that own the canonical plans.

The first version should be deliberately operationally small. Existing portfolio configuration
provides self-description and discovery scope; it does not require a manually maintained central
repository list. The coordinator refreshes its view on demand before portfolio analysis. It scans
the default branch as the baseline and scans only plan documents added or materially changed on
accessible unmerged remote branches. Reliable pre-push announcements, local outboxes, webhooks,
GitHub Apps, and background schedules remain possible later additions rather than prerequisites.

This discovery is `implementation-handoff-ready`, not implementation authority. It records the
discussion, evidence, approved target model, requirements, compatibility migration, and frozen
implementation capability. The user approved the metadata contract, identifier normalization,
active-branch scope, activation publication authority, first scanner adapters, family authority,
portfolio configuration evolution, and branch-aware migration approach on 2026-07-31. A future
implementation plan must preserve those decisions and still requires separate approval before
execution.

This discovery is both:

- `routing-bearing`, because future implementation would change managed planning routes,
  coordination ownership, publication checkpoints, and external Git-host selection; and
- `recipe-bearing`, because reliable discovery, validation, branch-diff ingestion, migration, and
  catalog generation require reusable executable mechanics rather than prose-only instructions.

## Essential Context

| Path | Reason |
| --- | --- |
| `AGENTS.md` | Public-core safety, producer authority, planning routes, external-action boundaries, and repository change rules. |
| `README.md` | Public repository purpose and validation entry points. |
| `docs/repo/reference/placement-contract.md` | Current ownership and placement contract for plan-register state files and consumer-owned planning content. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact classes for future metadata, schema, generated-path, sync, validator, migration, and safety changes. |
| `components/planning-workflows/managed/reference/planning-document-lifecycle.md` | Current canonical discovery and implementation document types, lifecycle headers, plan bundles, families, and register relationship. |
| `components/planning-workflows/managed/reference/plan-register-format.md` | Current register fields, sequential-ID examples, coordination namespace rules, entry format, relations, and coverage doctrine. |
| `components/planning-workflows/managed/runbooks/maintain-plan-register.md` | Current local-first, agent-maintained coordination procedure and pending-sync fallback. |
| `components/planning-workflows/scaffolds/plan-register.md` | Current kit-initialized consumer register baseline. |
| `schemas/kit-config.schema.json` | Current portfolio enrollment configuration and likely authority for future coordination settings. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Route-before-execution doctrine for selecting Git-host, coordinator, CLI, service, or local mechanics. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Maturity guidance for future reliable coordination mechanics. |
| `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` | Promotion guidance when repeated coordination and reconciliation should become reusable scripts. |
| `docs/repo/plans/plan-register.md` | Current producer-repository register and concrete evidence of the large shared-file model. |
| `docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_discovery_doc.md` | Original reusable plan-register and optional coordination-home discovery. |
| `docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_implementation_doc.md` | Implemented local-register, coordination-home, configuration, scaffold, and pending-sync model. |
| `docs/repo/plans/coordination-home-register-id-namespace/coordination-home-register-id-namespace_implementation_doc.md` | Existing repository-qualified coordination-home ID doctrine and collision analysis. |
| `docs/repo/plans/plan-register-portfolio-doctrine-refinement/plan-register-portfolio-doctrine-refinement_implementation_doc.md` | Existing local versus coordination-home coverage and relation refinements. |
| `docs/repo/plans/plan-register-dirty-target-safety/plan-register-dirty-target-safety_implementation_doc.md` | Existing compatibility rules for agent edits to dirty coordination-home registers. |
| `docs/repo/plans/plan-register-session-lifecycle-hardening/plan-register-session-lifecycle-hardening_implementation_doc.md` | Existing session-reference and lifecycle maintenance doctrine. |

## Table Of Contents

- [Section 1 - Problem Framing](#section-1---problem-framing)
- [Section 2 - Current Evidence And Gaps](#section-2---current-evidence-and-gaps)
- [Section 3 - Candidate Target Operating Model](#section-3---candidate-target-operating-model)
- [Section 4 - Requirements And Constraints](#section-4---requirements-and-constraints)
- [Section 5 - Decision Inventory And Recommendations](#section-5---decision-inventory-and-recommendations)
- [Section 6 - Open Questions And Assumptions](#section-6---open-questions-and-assumptions)
- [Section 7 - Risks And Validation](#section-7---risks-and-validation)
- [Section 8 - Implementation Handoff](#section-8---implementation-handoff)
- [Revision Notes](#revision-notes)

# Section 1 - Problem Framing

## Problem Statement

The current register model combines three concerns in one manually edited Markdown document:

1. stable identity for plans;
2. a local overview of canonical planning authority; and
3. a synchronization surface for portfolio coordination.

That coupling creates several practical failures:

- Sequential identifiers require shared allocation state. Two branches can independently choose
  the same next number without either branch being wrong from its local view.
- A bare identifier such as `PR45` carries no title, document kind, product, repository, or family
  context and can also be confused with a Git pull request.
- Every new or materially changed plan edits the same register file, often touching the same
  timestamp and insertion location even when the plans are unrelated.
- Discovery and implementation documents can be combined into one register entry even though they
  have different purposes, lifecycle states, and canonical paths.
- The current coordination-home model permits selective portfolio coverage, while the intended
  strategic coordinator needs a complete view of formal plans in enrolled repositories.
- Manual agent synchronization can provide early visibility, but omissions, unavailable
  coordination homes, and unsafe target-register edits can leave the global view incomplete.
- Git-host scanning can verify pushed work, but it cannot observe an uncommitted local plan.
- A plan can be active and approved for execution on a work branch while remaining unmerged.
  Calling that plan merely `proposed` confuses plan authority with Git integration state.
- Copying lifecycle and relationship metadata into separate local entry files would create
  another source of drift when the canonical discovery or implementation document already owns
  those facts.
- A mechanical migration of current register rows would reproduce incomplete coverage, combined
  entries, weak identifiers, missing families, and stale relationships instead of using semantic
  understanding to improve the catalog.

The problem is therefore not only duplicate numbering. It is the lack of a branch-aware,
source-derived, semantically useful catalog and reconciliation model.

## User Intention

Create an Operating Kit planning and coordination model in which:

- agents and humans can understand a plan reference without private conversational context;
- discovery documents, implementation documents, and plan families are distinguishable;
- parallel branches do not need to coordinate a shared next number or edit the same local register
  section;
- each enrolled repository's formal plans contribute to a complete coordination-home overview;
- the coordination home can perform strategic analysis across products, capabilities, plan
  families, and seemingly isolated features;
- authorized work on non-default branches is visible and correctly classified;
- the coordinator obtains a fresh view automatically whenever strategic analysis begins, without
  depending on a developer to update a register;
- repositories join a portfolio through existing Operating Kit markers and portfolio
  configuration rather than a separate repository registry;
- a fresh Kit-enabled repository can be configured as a coordination home using only managed
  concepts, setup instructions, schemas, and commands;
- requesting activation makes the plan checkpoint remotely visible without a second push-approval
  exchange while keeping broader external actions outside that authority;
- canonical planning documents remain the source of truth; and
- existing plans are semantically inventoried and migrated without rewriting their historical
  content dates.

## Goals

- Define a human-meaningful identity model that does not depend on a global sequential counter.
- Define separate catalog identities for discovery, implementation, and plan-family authority.
- Define the structured metadata contract inside canonical planning documents.
- Remove the need for manually maintained local sidecar entry files when canonical metadata is
  sufficient.
- Define local and global register views as derived summaries or stable entry points.
- Give the coordination home complete coverage of formal plans in explicitly enrolled
  repositories.
- Preserve repository ownership of canonical plan facts while enabling global strategic
  interpretation.
- Represent plan lifecycle, Git visibility, and catalog verification as separate concepts without
  adding a redundant execution-authority state.
- Retain rich metadata without making every classification mandatory or creating a centrally
  controlled ontology in the first version.
- Reuse portfolio configuration for automatic member discovery within an authorized scope,
  without adding a separate central repository list.
- Define a small first-version synchronization contract based on an on-demand source scan before
  coordination analysis.
- Treat the default branch as the catalog baseline and non-default branches as change overlays so
  inherited plans are not duplicated.
- Put reusable coordination-home concepts, setup, scanning, validation, and migration mechanics in
  the Operating Kit while leaving each coordination home's portfolio data and analysis local.
- Define branch publication checkpoints that make active work remotely visible without treating
  merge as a prerequisite for execution authority.
- Define a semantic migration process that evaluates each existing plan rather than blindly
  converting legacy register rows.
- Preserve existing IDs as aliases and preserve content-update dates during metadata migration.
- Define generic provider-neutral coordination contracts with a practical GitHub adapter when
  GitHub owns the enrolled repositories.
- Keep private repository topology, source content, credentials, and business-specific strategic
  analysis outside generic managed doctrine.

## Non-Goals

- Do not implement the new catalog, migration, schema, validator, coordinator, workflow, service,
  or Git-host integration in this discovery.
- Do not choose exact implementation epics or command names.
- Do not perform branch creation, commits, pushes, pull requests, merges, releases, repository
  settings, or external service installation merely by approving this discovery; the activation
  policy governs future plan-workflow authority.
- Do not make merge to the default branch a prerequisite for executing an explicitly activated
  implementation plan.
- Do not make a central coordinator the source of truth for canonical plan content.
- Do not require a full copy of every discovery or implementation body in the coordination home.
- Do not treat execution logs, task checkboxes, or session transcripts as independent strategic
  plans by default.
- Do not turn the plan catalog into a sprint tracker, work board, next-action list, or live
  progress database.
- Do not infer portfolio enrollment from repository names, sibling folders, organizations, or
  private conventions.
- Do not add a manually maintained coordination-home repository registry.
- Do not promise global visibility of uncommitted or unpushed local work in the first version.
- Do not require pre-push announcements, local outboxes, webhooks, GitHub Apps, background
  services, or scheduled scans in the first version.
- Do not scan and duplicate every inherited plan on every remote branch.
- Do not require a global controlled vocabulary for every optional strategic field before the
  catalog can operate.
- Do not force unknown external Operating Kit consumers through an automatic destructive
  migration.

## Priority Order

1. Complete and trustworthy strategic visibility.
2. Canonical source-of-truth clarity.
3. Correct representation of active branch work.
4. Human comprehension and semantic identity.
5. Collision and merge resilience under parallel branches.
6. Automatic portfolio discovery without a second repository registry.
7. Explicit authority and safe external-action boundaries.
8. Historical accuracy and reversible migration.
9. Provider-neutral reusable doctrine.
10. Low ongoing operational and maintenance overhead.

## Durable Principles

- One fact should have one canonical owner.
- Plan lifecycle is not Git integration state.
- `active` remains the lifecycle signal that a plan is current and approved for execution.
- Discovery, implementation, and family records remain distinguishable even when they share a
  topic or bundle.
- A coordination home catalogs repository-owned facts and owns portfolio interpretation; it does
  not silently rewrite member plan authority.
- Global completeness means complete coverage of formal registered planning authority in enrolled
  repositories, not indiscriminate collection of every note or task.
- A coordination analysis must refresh source observations before claiming a current overview.
- The default branch is the baseline; non-default branch observations are overlays containing only
  added or materially changed canonical plans.
- Uncommitted and unpushed local state is outside the first-version global visibility boundary.
- External writes and branch publication remain approval- and repository-policy-bound; an explicit
  activation request is the approval for its plan-only normal commit and push.
- Migration must improve semantic quality rather than reproduce known register defects.
- Content history and catalog-adoption history require separate timestamps.
- Generated summaries must not become feature-branch merge hotspots.
- Portfolio membership is configuration-driven and discovered automatically inside authorized
  provider or local scope; it is not maintained in a second repository list.
- Generic coordination mechanics belong to the Operating Kit; portfolio configuration, catalog
  state, and strategic interpretation belong to the coordination-home repository.

## Candidate Domains And Decision Clusters

| Cluster | Domain | Why it is involved | Likely authority |
| --- | --- | --- | --- |
| Identity and canonical metadata | Planning lifecycle and schema | IDs, document kinds, families, aliases, dates, and relations need a durable contract. | Planning-workflows managed references and schemas |
| Local catalog and generated views | Repository structure and documentation | The shared register file and possible sidecars affect placement, ownership, and merge behavior. | Structure governance and planning workflows |
| Portfolio catalog and strategic overlay | Coordination operating model | Complete portfolio awareness and cross-plan interpretation need separate fact and strategy ownership. | Coordination-home configuration and repo-owned portfolio doctrine |
| Branch visibility and publication | Git workflow and external actions | Active work may be authorized before merge and invisible before push. | Repository policy, planning runbooks, Git-host route |
| Portfolio discovery and source refresh | Operational recipe and automation | Config-driven repository discovery, baseline scans, branch overlays, and freshness are repeated mechanics. | Agent-interface recipe doctrine and coordinator implementation |
| Coordination-home bootstrap | Installation and role configuration | A fresh Kit-enabled repository must be able to establish the same coordination model without prior conversational context. | Managed setup runbook, config schema, commands, and repo-owned portfolio surfaces |
| Semantic migration and compatibility | Consumer lifecycle and adoption | Existing plans require interpretation, aliases, timestamp preservation, and public consumer safety. | Producer migration plan and consumer-impact contract |

## Success Criteria For This Discovery

This discovery is implementation-handoff-ready when:

- the intended global coverage and strategic purpose are explicit;
- the distinction between plan lifecycle, Git visibility, and verification is explicit;
- the role of canonical structured metadata is framed;
- automatic portfolio discovery and on-demand source refresh are described;
- default-branch baselines and branch-diff overlays are described;
- the reusable coordination-home setup boundary is explicit;
- the semantic migration direction and historical date rules are recorded;
- functional and non-functional requirements are traceable;
- implementation-shaping decisions are approved and remaining questions are non-blocking;
- parallel migration and cutover behavior are explicit;
- implementation capability-scope blocks freeze the required workflows, boundaries, and evidence;
  and
- a future planner can draft one coherent implementation path without reinventing the decisions.

# Section 2 - Current Evidence And Gaps

## Current Repository Evidence

At discovery creation:

- `docs/repo/plans/plan-register.md` contains 26 sequential `OK-PR-001` through `OK-PR-026`
  entries in one repeated-section Markdown file.
- The register contains 7 `discovery-plan` entries, 17 `implementation-plan` entries, and 2
  combined `discovery-plan + implementation-plan` entries.
- The repository contains 11 `*_discovery_doc.md` files, 20 `*_implementation_doc.md` files, and
  20 `*_execution_log.md` files before this discovery is added.
- The current register coverage note explicitly states that earlier repository plans may be added
  later, so the register is not a complete inventory of canonical planning documents.
- No current register entry uses `Type: plan-family`, although managed lifecycle doctrine supports
  plan-family structure and the register format supports `plan-family` records.
- The required planning header contains only `Last updated`, `Created`, and `Status`, with optional
  completion, supersession, and execution-log fields.
- The register repeats title, type, purpose, owner, canonical paths, lifecycle dates, relations,
  session references, and coordination notes outside the canonical documents.
- The current checkout has no configured `portfolio` block, so no coordination-home update is
  required for this discovery under the current runbook.

## Existing Model Strengths To Preserve

The redesign should preserve what the current model established successfully:

- canonical planning documents beat register snapshots on conflict;
- plan registers are indexes rather than task trackers;
- local planning work continues when coordination is unavailable;
- coordination is explicit and configuration-driven;
- member repositories own their canonical planning documents;
- coordination homes use namespaced source identity;
- session references are recovery handles rather than transcript summaries;
- dirty unrelated repository state does not automatically block a compatible coordination update;
- existing consumer-owned files are preserved rather than silently overwritten; and
- generic managed doctrine does not hardcode one private coordination repository.

The current `portfolio` configuration should also be preserved. Automatic onboarding means that
the coordinator discovers configured member repositories inside its authorized scope; it does not
mean that portfolio affiliation, coordination-home identity, or discovery scope becomes implicit.

## Gap 1 - Sequential Identity

The current short sequential identifiers are easy to sort but require a shared allocator. Git
branches do not share mutable state after they diverge. Two correct branch-local choices can
therefore collide at merge or coordination time.

Repository qualification such as `EXAMPLE-AUTOMATION-PR-001` prevents different repositories from
colliding in one coordination home, but it does not prevent two branches in the same repository
from allocating the same next local number.

The identifier also fails as conversational context. `OK-PR-023` does not reveal that the record is
an implementation plan about dirty coordination-target safety.

## Gap 2 - Shared Register Merge Hotspot

Changing the ID scheme alone does not solve the shared-file problem. Parallel branches may still
update the register timestamp and insert entries under the same heading. The conflict is caused by
the central mutation surface as well as by duplicate identifiers.

## Gap 3 - Canonical Metadata Duplication

Separate local entry files would reduce shared-file conflicts, but they would duplicate lifecycle,
relations, title, and classification data already associated with the canonical plan. If canonical
planning documents contain a complete structured metadata block, local sidecar entries become
unnecessary by default.

Generated or materialized global catalog records may still be useful as coordinator-owned cache,
but those records should be replaceable from canonical source and should not become a second
manual authority.

## Gap 4 - Incomplete Global Coverage

The original coordination discovery recommended selective coordination-home coverage. The current
user intention is broader: the coordination home should have a total overview of formal planning
authority in enrolled repositories so that it can recognize strategic relationships between work
that appears locally isolated.

A complete factual catalog does not mean copying all plan bodies. It means every formal discovery,
implementation, and first-class family record is represented with sufficient structured metadata
and a canonical pointer.

## Gap 5 - Strategic Interpretation Ownership

Local repositories can state their product, capabilities, family, and known relations. Some
strategic relationships become visible only when plans from several repositories are considered
together.

If global strategic interpretation is written into generated member snapshots, later
reconciliation could overwrite it. The coordination home therefore needs a separate
coordination-owned layer for portfolio families, themes, relationships, and analysis.

## Gap 6 - Branch Visibility And Authority

Agents commonly create a work branch, draft and activate an implementation plan there, and begin
authorized execution on the same branch. Default-branch-only scanning sees that work too late.

However, an unmerged plan is not necessarily unapproved. Current lifecycle doctrine already gives
`active` the meaning “current and approved for execution.” A separate `Authority: approved`
state would duplicate lifecycle semantics and permit contradictory combinations.

The catalog needs orthogonal observations:

- plan lifecycle: `draft`, `active`, `completed`, `superseded`, or `archived`;
- Git visibility: pushed branch, open pull request when available, or default branch; and
- catalog freshness: source revision, scan time, and any stale or unavailable condition.

## Gap 7 - Coordination Reliability

Runbooks can instruct agents to update coordination state, but any design that requires a second
manual register edit or an announcement to a separate repository still depends on perfect agent
discipline. Event infrastructure can reduce latency, but webhooks, apps, outboxes, credentials,
retries, hosting, and cleanup create substantial operational scope.

The simplest reliable first version is source-derived and pull-based: before the coordination home
performs strategic analysis, one command discovers configured members and refreshes their
catalogs. Pushed source is the visibility boundary. Early publication of an activated plan makes
the work visible on the next scan; unpushed work remains local.

A naïve scan of every plan file on every branch would be both expensive and incorrect because each
branch inherits the default branch's plan history. The scanner therefore needs a default-branch
baseline plus branch overlays containing only canonical plans added or materially changed relative
to the merge base. Fully merged branches contribute no overlay. Multiple changed observations of
the same plan remain separate and are flagged instead of silently reconciled.

Optional scheduled or event-triggered execution can invoke the same scanner later without changing
the catalog contract.

## Gap 8 - Historical Migration Quality

The existing register is valuable migration evidence, but it is known to have incomplete coverage,
combined record types, no family records, and identifiers that communicate little meaning.

A mechanical row conversion would reproduce these defects. The repository set is small enough for
an implementation agent to enumerate and semantically evaluate every canonical plan, use the
legacy register as reconciliation evidence, and produce a migration ledger for ambiguous cases.

Adding catalog metadata must not change the historical content-update date merely because a
machine-readable block was inserted. Content changes and catalog metadata changes require distinct
timestamps.

## Gap 9 - Portfolio Discovery And Coordination-Home Bootstrap

Automatic member onboarding should not require a manually maintained central repository list.
However, automatic discovery without explicit scope or affiliation would be unsafe and could
silently collect unrelated repositories.

The existing portfolio configuration can solve both concerns:

- a member repository self-describes its portfolio affiliation and member identity;
- a coordination home declares its role, home identity, and authorized local or provider scope;
- the scanner finds repositories with an Operating Kit marker on their default branch, reads their
  portfolio configuration, and includes only members matching the coordination home; and
- Kit-enabled but unconfigured repositories may be reported as candidates, but do not silently
  become full members.

Default-branch configuration determines membership. A feature branch cannot enroll its repository
by itself.

The Operating Kit must also explain and implement this model completely. A fresh repository with a
standard Kit installation and no knowledge of this discussion should be able to become a
coordination home through a low-context request such as “Set this repository up as a coordination
home.” Managed instructions should route the agent through read-only inspection, scope and
provider selection, required authorization, configuration, creation of repository-owned portfolio
surfaces, an initial scan, and validation.

# Section 3 - Candidate Target Operating Model

## 3.1 Canonical Planning Documents

Each formal discovery and implementation document should carry a versioned, machine-readable
metadata block immediately below its first heading. The existing three-line planning header and
first heading remain in place:

- the heading owns the human-readable title;
- `Created` owns the original creation date;
- `Status` owns plan lifecycle; and
- `Last updated` owns the last meaningful plan-content or authority update.

The metadata block uses the YAML subset already supported by the self-contained Go CLI and is
bounded by explicit `BEGIN CODEHEART PLAN METADATA` and `END CODEHEART PLAN METADATA` HTML comment
markers. The markers are machine anchors; the fenced YAML between them remains visible and
reviewable in rendered Markdown.

Approved placement:

````markdown
Last updated: YYYY-MM-DDTHH:MM:SSZ (UTC)
Created: YYYY-MM-DD
Status: draft | active | completed | superseded | archived

# Human-Readable Plan Title

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: repository.kind.stable-slug
  kind: discovery
  purpose: Concise strategic purpose
  first_cataloged: YYYY-MM-DDTHH:MM:SSZ
  catalog_metadata_updated: YYYY-MM-DDTHH:MM:SSZ
```
<!-- END CODEHEART PLAN METADATA -->
````

The required block shape is:

```yaml
plan:
  schema_version: 1
  id: <repository>.<kind>.<stable-slug>
  kind: discovery | implementation | family
  purpose: <concise strategic purpose>
  first_cataloged: <catalog adoption timestamp>
  catalog_metadata_updated: <last classification update>
```

Optional rich fields use the same block:

```yaml
plan:
  family: <stable family identifier when applicable>
  products:
    - <public-safe product or domain identifier>
  capabilities:
    - <public-safe capability identifier>
  strategic_themes:
    - <optional source-owned theme>
  relations:
    - kind: <relation>
      target: <plan identifier>
  legacy_aliases:
    - <legacy register identifier>
```

The first version validates the required core, optional field shapes, and reserved relation kinds
without imposing a globally rigid product or capability ontology. Single-line purpose values and
portable scalar/list/map forms keep the block within the existing YAML subset. Catalog timestamps
use UTC RFC 3339 values ending in `Z`.

The canonical plan owns its ID, kind, title, purpose, lifecycle, family, products, capabilities,
known relations, aliases, and content chronology. The scanner derives repository, canonical path,
branch, commit, pull request, checksum, and scan time. The coordination home owns cross-repository
relationships, portfolio-level families or themes, priorities, and strategic analysis.

Migration-only metadata insertion preserves the historical `Last updated` value. It records the
adoption time in `first_cataloged` and `catalog_metadata_updated`. A concurrent change to plan
meaning, scope, lifecycle, relations, or authority updates `Last updated` normally. This
catalog-only exception must be added explicitly to lifecycle doctrine so historical chronology
does not depend on migration timing.

## 3.2 Record Granularity

Discovery and implementation documents should have separate catalog identities even when they live
in the same plan bundle and share a family. Their purposes, lifecycle, authority, timestamps, and
relations may differ.

Execution logs should normally remain evidence associated with an implementation plan rather than
independent strategic records.

A plan family should become first-class when independently executable sibling plan bundles need
shared discoverability. Preserve the current second-sibling trigger and place canonical family
metadata in the family folder's `README.md` with `kind: family`.

The family record owns family title, purpose, strategic scope, lifecycle when applicable, and
source-owned relations. It does not manually enumerate every member. The scanner derives
membership from plan records whose `family` field points to the family ID. A single plan may carry
a provisional family reference before a family record exists; validation reports that condition
without silently inventing a family.

## 3.3 Identity

The identifier should:

- include a stable repository or source namespace;
- include human-meaningful topic context;
- distinguish discovery, implementation, and family kinds;
- avoid dependence on a shared next number;
- remain stable when title wording or lifecycle changes;
- work safely in Markdown, filenames, Git, and supported operating systems;
- preserve old `PR-*` identifiers as aliases; and
- make a same-slug choice on two branches visible as a possible overlap rather than hiding it
  behind unrelated generated numbers.

The recommended initial shape is:

```text
<repository>.<kind>.<stable-slug>
```

Example:

```text
operating-kit.discovery.semantic-plan-catalog
```

Family, owner, product hierarchy, priority, and lifecycle do not belong in the ID because they may
change. A generated collision token is not required initially. If two branches independently
choose the same ID, the coordinator preserves both branch observations and flags the overlap. That
may reveal duplicated intent. If they are genuinely distinct plans, an author adds a meaningful
qualifier to one slug.

Normalization is fixed as follows:

- use lowercase ASCII;
- use dots only between repository, kind, and stable-slug segments;
- use hyphens between words inside a segment;
- collapse runs of unsupported characters to one hyphen and trim leading or trailing hyphens;
- allow only `discovery`, `implementation`, and `family` as kind segments;
- derive the repository segment from the configured stable member repository ID rather than a
  folder name;
- require a meaningful explicit ASCII transliteration when normalization would otherwise erase a
  non-ASCII name; and
- reject invalid IDs instead of silently generating an opaque replacement.

Canonical filenames keep the existing lowercase hyphenated plan-bundle convention and remain
independent from IDs. Moving a file, changing a title, or renaming a repository does not
automatically change an established plan ID. A deliberate identity migration preserves the prior
ID as an alias.

This discovery uses `OK-PR-027` only as a legacy register identity under the currently active
format. Adoption of a new scheme should retain it as an alias rather than pretending the new
scheme was already approved.

## 3.4 Local Views

Canonical planning documents should be the local catalog source. A local overview can be:

- generated on demand by an Operating Kit command;
- generated centrally after merge;
- represented by a stable `plan-register.md` entry point that explains how to inspect the
  distributed records.

Feature branches should not manually edit or regenerate a committed shared summary. Otherwise the
same merge hotspot returns. After migration, manual local register maintenance should stop.
Whether the legacy `plan-register.md` becomes a stable explanatory entry point or is retired is an
implementation choice; it should not remain a second catalog authority.

## 3.5 Complete Coordination-Home Catalog

The coordination home should ingest every formal discovery, implementation, and first-class family
record from each configured member repository discovered inside its authorized scope.

The coordination home should not maintain a second repository registry. Discovery should use:

1. the coordination home's portfolio role, stable home identity, and authorized provider or local
   scope;
2. valid Operating Kit lock and config markers on each candidate repository's default branch;
3. the candidate's configured member role and stable repository ID; and
4. an exact match between the member's configured coordination-home identity and the scanning
   home.

The portfolio block receives its own nested schema version while preserving the root Kit config
contract. The approved member shape is:

```yaml
portfolio:
  schema_version: 2
  role: member
  member_repository_id: operating-kit
  coordination_home_id: codeheart-development
```

The approved coordination-home shape is:

```yaml
portfolio:
  schema_version: 2
  role: coordination-home
  coordination_home_id: codeheart-development
  discovery:
    sources:
      - kind: github-owner
        owner: example-organization
```

Provider or owner scope that is safe and intentionally shared may remain in repository
configuration. Machine-specific local roots belong in `.codeheart/local/` rather than committed
configuration.

Existing v1 fields such as `coordination_home_path` and
`coordination_home_register_path` remain readable compatibility inputs. New provider-based member
configuration does not require them. The migration must accept both v1 and v2 portfolio blocks and
must not discard working legacy paths merely to rename them.

A repository becomes a catalog member only when all membership predicates match on its default
branch. A feature branch cannot self-enroll. An unconfigured Kit-enabled repository may appear in
a non-authoritative candidate report, but its plans are not ingested and its configuration is not
modified silently.

The factual catalog should contain source-derived snapshots such as:

- source repository;
- plan ID and legacy aliases;
- title and kind;
- family;
- lifecycle status;
- products and capabilities;
- known source-owned relations;
- canonical path;
- source branch or pull request;
- source commit;
- verification and freshness state.

The coordination home should not require full plan bodies for routine strategic orientation. It
may fetch canonical documents when deeper analysis is needed and permitted.

## 3.6 Portfolio Strategic Layer

The coordination home should own separate durable records for:

- portfolio plan families;
- strategic themes;
- cross-repository relations;
- conflicts and overlaps;
- strategic analysis; and
- coordination decisions.

Generated member snapshots must not overwrite these coordination-owned interpretations.

## 3.7 Branch-Aware Catalog States And Overlays

The catalog should not replace plan lifecycle with one overloaded proposal state.

The scanner should build one default-branch baseline per member repository, then enumerate every
accessible unmerged remote branch. For each branch it compares against the default-branch merge
base and ingests only canonical plans that were added or materially changed. Unchanged inherited
plans do not become branch observations, branches with no canonical plan changes are ignored, and
fully merged branches do not remain overlays.

The first version does not require a branch-name prefix, open pull request, or recency cutoff.
Those filters could silently hide legitimate agent work. Old unmerged plan observations are marked
stale for operator attention rather than omitted. Open pull requests enrich branch observations
with review context when available but are not an enrollment or visibility prerequisite.

An active plan on a pushed work branch is authorized work that has not yet been integrated. It
should participate in strategic analysis with its source state and freshness visible.

If several branches modify the same plan ID, each source-ref observation remains visible and the
catalog reports a conflict or overlap. It must not pick a winner by scan order.

## 3.8 On-Demand Coordination Refresh

The first-version flow is:

```text
Agent drafts canonical plan
-> user requests activation
-> activation supplies plan-only commit-and-push authority
-> coordination analysis invokes one portfolio scan
-> scanner discovers configured members
-> scanner refreshes each default-branch baseline
-> scanner applies plan-changing unmerged-branch diffs
-> coordinator performs analysis on the refreshed catalog
```

The scan should be idempotent, record source revisions and failures, and use the repository or
provider credentials already available to the authorized operator. The provider-neutral scanner
belongs in the existing self-contained Go CLI and has two first-version source adapters:

- a local Git adapter for configured nearby repositories, fixtures, deterministic testing, and
  offline use; and
- a GitHub adapter that uses the operator's existing authenticated `gh` CLI session for repository
  discovery and read-only source access.

The GitHub adapter performs an authentication and scope preflight, never stores tokens in Kit
configuration or catalog output, and distinguishes inaccessible repositories from repositories
with no formal plans. Neither adapter executes code from scanned branches.

The scan is automatic when the coordination home is used, not a background watcher. A later
scheduled job, webhook, or GitHub App may call the same scan contract if practical evidence
justifies proactive alerts or tighter latency.

## 3.9 Publication Checkpoints

For coordinated repositories, activation of a formal implementation plan should normally create a
publication checkpoint:

1. confirm or create the dedicated work branch under repository policy;
2. set the canonical implementation plan to `active`;
3. preserve approval or activation evidence;
4. commit the activated plan as a planning checkpoint;
5. push the work branch under the activation authority;
6. optionally open a draft pull request against the configured integration branch; and
7. continue implementation on the same branch.

An explicit user request to activate a plan is sufficient scoped authority to create or use the
work branch, commit the canonical plan and directly required planning metadata, and push that
checkpoint. The agent does not ask for a second push approval.

A user-requested material update to an active plan carries the same plan-only commit-and-push
authority when the update changes plan scope, authority, dependencies, or expected outcomes.
Typographical or formatting edits do not create an independent publication requirement.

The scoped authority does not include unrelated files, implementation code not otherwise
authorized, draft pull-request creation, merge, release, force-push, destructive branch changes, or
a different repository or branch. The agent must stop or ask when the target is ambiguous,
unrelated changes would be included, normal push is rejected, or repository policy requires a
broader or destructive action.

If the user explicitly limits or revokes publication authority, or normal push is unavailable,
authorized local work may continue according to user direction, but the coordinator cannot
represent that local-only work in the first version. The agent must state this visibility
limitation clearly.

## 3.10 Semantic Migration

Migration is branch-aware, optimistic, and compatible with ongoing work. It uses four phases.

### Phase 1 - Compatibility Tooling

- Implement metadata parsing, legacy fallback, validators, scanner, and generated views.
- Keep current authoring behavior until the mixed-format reader is proven.
- Treat canonical plan content as authority and the legacy register as fallback index evidence.

### Phase 2 - Semantic Inventory And Safe Migration

The implementation agent should evaluate, for each canonical planning document:

- actual kind and purpose;
- lifecycle evidence;
- owning repository and product;
- stable semantic identity;
- family membership;
- capabilities and strategic relevance;
- discovery-to-implementation relationship;
- dependencies, supersession, and related plans;
- legacy IDs and session references;
- whether the record is current authority or historical evidence; and
- confidence and unresolved ambiguity.

Legacy registers remain evidence for aliases, session refs, coordination notes, and relations.
They do not become the migration authority.

The inventory records each plan's source commit and content hash and detects every active branch
that adds or changes the same plan. Before inserting metadata, migration verifies that the source
commit and hash are still current:

- unchanged and uncontested files may be migrated;
- changed files are skipped and semantically reevaluated from the latest source;
- plans touched by an active branch are deferred or assigned to the agent already owning that
  branch;
- unchanged inherited copies are never migrated separately on every branch; and
- every skip, conflict, ambiguity, alias, and assignment is recorded in the migration ledger.

### Phase 3 - Authoring Cutover

- Require metadata for newly created formal plans.
- Stop manual register maintenance once mixed-format generated views are reliable.
- Continue reading the legacy register as fallback evidence for unmigrated plans.
- Let older active branches remain compatibility-readable; when they materially change a plan,
  their owning agent adds metadata or the plan remains a ledgered deferral.

This deliberately avoids a long dual-write period in which agents must maintain both canonical
metadata and the shared register.

### Phase 4 - Coverage Closure

- Reach complete default-branch metadata coverage for Codeheart-controlled repositories.
- Reconcile every active branch overlay.
- Resolve or explicitly retain every migration-ledger deferral.
- Run final local and portfolio catalog scans.
- Retain legacy aliases and compatibility reading for historical branches and unknown consumers.

No repository-wide work freeze is required. At most, migration coordinates briefly around an
individual plan file whose source and active branch are changing concurrently.

## 3.11 Consumer Compatibility

Codeheart-controlled repositories may receive a bounded semantic migration under an approved
implementation plan. Unknown Operating Kit consumers require:

- a compatibility reader for the legacy register shape during transition;
- an explicit migration command or runbook;
- preservation of consumer-owned planning content;
- dry-run and ambiguity reporting;
- no forced rewrite during ordinary sync, repair, or upgrade; and
- clear adoption and rollback notes.

## 3.12 Operating Kit And Coordination-Home Ownership

All reusable concepts and mechanics needed to understand, establish, and operate coordination
belong in managed Operating Kit source:

- the canonical plan metadata and catalog model;
- member and coordination-home roles;
- portfolio configuration schema and compatibility behavior;
- coordination-home setup and maintenance runbooks;
- plan and portfolio scan commands;
- branch-overlay logic;
- validators;
- migration mechanics; and
- generic templates or scaffolds.

The Operating Kit should make the setup route discoverable but keep its heavier runtime dormant
unless a repository is configured for the coordination-home role. Becoming a coordination home is
an explicit user action; automatic member discovery begins afterward.

Each repository retains its own instance data:

- member repositories own canonical plans and source metadata;
- the coordination-home repository owns its actual portfolio catalog or rebuildable cache,
  strategic relationships, themes, priorities, local discovery scope, and analyses; and
- repository-owned wrapper documentation should contain only local provider, scope, command, or
  policy details rather than duplicating the generic runbook.

A decisive acceptance test is that a fresh repository with the standard Operating Kit and a
low-context agent can establish a coordination home, discover matching member repositories without
a repository list, run an initial scan, and validate the result.

# Section 4 - Requirements And Constraints

## Functional Requirements

| ID | Requirement | Priority |
| --- | --- | --- |
| `FR-1` | Represent every formal discovery document with its own stable catalog identity. | must |
| `FR-2` | Represent every formal implementation document with its own stable catalog identity. | must |
| `FR-3` | Represent first-class plan families independently from their member plans. | must |
| `FR-4` | Avoid shared sequential allocation as the only source of plan identity. | must |
| `FR-5` | Preserve legacy register identifiers as searchable aliases. | must |
| `FR-6` | Store versioned rich catalog metadata in a bounded YAML block below the canonical title, with required identity, kind, purpose, and catalog timestamps plus optional strategic fields. | must |
| `FR-7` | Preserve the existing header and title as canonical lifecycle, creation, content-update, and title sources; migration-only metadata insertion must not rewrite their chronology. | must |
| `FR-8` | Generate or derive local plan overviews without requiring feature branches to edit one shared summary. | must |
| `FR-9` | Give the coordination home complete coverage of formal plans in configured member repositories. | must |
| `FR-10` | Keep global strategic interpretation separate from source-derived member facts. | must |
| `FR-11` | Preserve plan lifecycle independently from Git visibility and catalog verification. | must |
| `FR-12` | Treat `active` as the existing execution-authority-bearing lifecycle state without a redundant authority enum. | must |
| `FR-13` | Discover candidate repositories automatically inside configured local or provider scope using valid default-branch Operating Kit lock and config markers. | must |
| `FR-14` | Treat only matching configured repositories as portfolio members; unconfigured Kit repositories may be reported as candidates but must not silently enroll. | must |
| `FR-15` | Refresh the portfolio catalog on demand before coordination-home strategic analysis. | must |
| `FR-16` | Scan each member's default branch as the complete baseline. | must |
| `FR-17` | Enumerate every accessible unmerged remote branch and ingest only canonical plans added or materially changed relative to the default-branch merge base. | must |
| `FR-18` | Ignore unchanged inherited plans and branches already fully merged into the baseline. | must |
| `FR-19` | Keep multiple changed branch observations of the same plan distinguishable and flag conflicts or overlaps. | must |
| `FR-20` | Treat a user request to activate a plan, or materially update an active plan, as scoped authority to commit and push the plan-only checkpoint to its work branch. | must |
| `FR-21` | Keep unrelated files, unauthorized implementation code, pull-request creation, merge, release, force-push, destructive branch actions, and ambiguous targets outside activation authority. | must |
| `FR-22` | Enumerate and semantically evaluate all existing canonical plans during Codeheart migration. | must |
| `FR-23` | Reconcile every migrated document against legacy register and repository evidence. | must |
| `FR-24` | Produce a migration ledger with confidence, ambiguity, legacy aliases, and unresolved findings. | must |
| `FR-25` | Support non-destructive compatibility for unknown consumer-owned legacy registers and both v1 path-based and v2 identity-and-scope portfolio configuration. | must |
| `FR-26` | Validate metadata markers, required fields, normalized IDs, kinds, lifecycle, family references, relation shapes, supported schema versions, and ownership violations. | must |
| `FR-27` | Record derived repository, ref, commit, canonical path, version evidence, and scan time in global observations. | must |
| `FR-28` | Expose last successful refresh, freshness, inaccessible members, and ingestion errors to the coordination home. | must |
| `FR-29` | Allow deeper strategic analysis to fetch canonical source only through the applicable repository and permission route. | should |
| `FR-30` | Keep agent-facing references human-readable by requiring title and kind on first conversational reference. | should |
| `FR-31` | Provide managed concepts, schemas, setup and maintenance runbooks, scanner, validators, and migration mechanics sufficient for a fresh Kit-enabled repository to become a coordination home. | must |
| `FR-32` | Keep actual portfolio scope, catalog or cache, strategic relationships, priorities, and analysis repository-owned in the coordination home. | must |
| `FR-33` | Allow a later scheduler or event trigger to invoke the same scan contract without making background infrastructure part of the first version. | should |
| `FR-34` | Provide a provider-neutral Go scanner with local Git and authenticated read-only GitHub-through-`gh` adapters in the first version. | must |
| `FR-35` | Use commit-and-hash preconditions and active-branch touch detection before semantically migrating an existing plan. | must |
| `FR-36` | Support a mixed-format authoring cutover that stops manual register writes before every historical plan must be migrated. | must |
| `FR-37` | Put first-class family authority in the family `README.md` at the current second-sibling trigger and derive membership from member plan metadata. | must |

## Non-Functional Requirements

| ID | Requirement | Priority |
| --- | --- | --- |
| `NFR-1` | Collision resilience: independently created branch records must not require a shared next-number lock. | must |
| `NFR-2` | Merge resilience: unrelated plans should normally change unrelated source files. | must |
| `NFR-3` | Idempotency: repeated source scans must converge without duplicate records. | must |
| `NFR-4` | Use-time freshness: coordination analysis must refresh or explicitly disclose that it is using a stale catalog. | must |
| `NFR-5` | Historical integrity: metadata migration must not rewrite content chronology. | must |
| `NFR-6` | Source clarity: canonical plan content must beat generated or cached catalog observations. | must |
| `NFR-7` | Security: scanning must use authorized read access and must not execute untrusted branch code. | must |
| `NFR-8` | Privacy: generic managed doctrine and global metadata must not expose secrets, customer data, private topology, or unapproved plan bodies. | must |
| `NFR-9` | Availability: local planning may continue when the coordinator or Git host is unavailable. | must |
| `NFR-10` | Observability: operators can see last successful scan, source revision, errors, inaccessible members, and stale observations. | must |
| `NFR-11` | Provider neutrality: generic doctrine must not require GitHub even when GitHub is the first practical adapter. | must |
| `NFR-12` | Portability: IDs and generated filenames must work on supported macOS and Windows filesystems. | must |
| `NFR-13` | Bounded maintenance: routine plan edits should not require hand-editing multiple index records. | must |
| `NFR-14` | Reversibility: generated catalog state can be rebuilt, and consumer migration has an explicit rollback or compatibility path. | must |
| `NFR-15` | Strategic usefulness: classifications support portfolio reasoning without becoming a high-churn project-management taxonomy. | must |
| `NFR-16` | Bounded scan scope: inherited default-branch plans must not be multiplied across branch overlays. | must |
| `NFR-17` | Automatic onboarding: adding a configured member inside authorized scope must not require updating a central repository list. | must |
| `NFR-18` | Dormancy: coordination-home runtime should remain inactive in ordinary member repositories unless the role is configured or a matching command is invoked. | must |
| `NFR-19` | Authority containment: activation-triggered publication must be plan-only, non-destructive, and limited to the unambiguous work branch and repository. | must |
| `NFR-20` | Concurrency safety: migration must never overwrite a plan whose recorded source version changed or whose active branch ownership is unresolved. | must |

## Operating Scenarios

### Scenario A - Parallel New Plans

Two agents create unrelated plans on separate branches. Each receives a stable semantic identity
without consulting a global sequence. Each branch changes its own canonical plan
document. Once pushed, the next on-demand scan sees each as a branch overlay without duplicating
the default-branch catalog.

### Scenario B - Competing Changes To One Plan

Two branches modify the same plan ID. The coordinator preserves both source-ref observations,
compares their changes with the baseline, and flags a conflict rather than letting scan order
overwrite either branch.

### Scenario C - Activation Publication

A user requests activation of an implementation plan. The agent commits the canonical plan and
directly required planning metadata and pushes the unambiguous work branch without asking for a
second push approval. It does not include unrelated files, open a pull request, merge, release, or
force-push. If normal push is unavailable or the target is ambiguous, the agent stops or reports
the visibility limitation.

### Scenario D - Published Active Branch

An activated plan is committed and pushed to its dedicated work branch before implementation
continues. Before the next strategic analysis, the scanner discovers the changed canonical plan
relative to the merge base and includes it as active, branch-visible, and unmerged.

### Scenario E - Automatic Member Discovery

A repository within authorized scope installs the Operating Kit and commits portfolio
configuration that names the coordination home. The next on-demand scan discovers it from its
default branch and includes it without an edit to a central repository list.

### Scenario F - Semantic Migration

The migration agent enumerates every canonical plan, reads each plan, proposes identity, family,
capability, lifecycle, and relations, reconciles the proposal with the legacy register, and records
its source commit and hash. A plan changed since inventory or touched on an active branch is
deferred or assigned to the branch owner rather than overwritten.

### Scenario G - Unknown Consumer

An external consumer upgrades the Operating Kit while retaining a consumer-owned legacy register.
Ordinary upgrade or sync does not rewrite the register or canonical plans. The consumer may run an
explicit dry-run migration later.

### Scenario H - Fresh Coordination Home

A user asks a low-context agent to set up a fresh Kit-enabled repository as a coordination home.
The managed route explains the role, obtains scope and authorization, writes repository-owned
configuration and portfolio surfaces, performs the initial member scan, and validates the result
without relying on this discovery conversation.

# Section 5 - Decision Inventory And Recommendations

## Decision Inventory

| ID | Question | Class | State | Depends on | Affects | Closure criteria |
| --- | --- | --- | --- | --- | --- | --- |
| `D-1` | What should be the primary source and syntax of catalog metadata? | implementation-shaping | approved | None | `FR-6`, `FR-7`, `FR-26`, `NFR-6`, `NFR-13` | Closed: bounded YAML metadata below the title, existing header and H1 authority preserved, rich optional fields, and derived scanner facts. |
| `D-2` | What identifier model replaces sequential allocation? | implementation-shaping | approved | `D-1` | `FR-1`-`FR-5`, `NFR-1`, `NFR-12` | Closed: lowercase ASCII `<repo>.<kind>.<stable-slug>`, independent filenames, no initial token, and preserved aliases. |
| `D-3` | What is the record granularity for discovery, implementation, execution logs, and families? | implementation-shaping | approved | `D-1` | `FR-1`-`FR-3`, `FR-11`, `FR-37` | Closed: discovery and implementation are separate, logs are evidence, and family README authority begins at the second sibling. |
| `D-4` | How complete should coordination-home coverage be? | implementation-shaping | approved | None | `FR-9`, `NFR-15` | Closed by explicit user direction: all formal plans in enrolled repositories contribute to the overview. |
| `D-5` | Who owns strategic families and cross-repository interpretation? | implementation-shaping | approved | `D-4` | `FR-10`, `FR-29`, `NFR-15` | Closed: member plans own source facts; the coordination home owns additional portfolio interpretation. |
| `D-6` | How should lifecycle, authority, Git state, and verification be represented? | implementation-shaping | approved | None | `FR-11`, `FR-12`, `FR-17`-`FR-19` | Closed by user correction: lifecycle carries execution authority; Git visibility and verification remain separate. |
| `D-7` | Which first-version synchronization model is required? | implementation-shaping | approved | `D-1`, `D-4` | `FR-15`-`FR-19`, `NFR-3`, `NFR-4`, `NFR-10`, `NFR-16` | Closed: on-demand source scan before analysis; background announcement and event infrastructure deferred. |
| `D-8` | When should coordinated plans be committed and pushed? | implementation-shaping | approved | `D-6`, `D-7` | `FR-20`, `FR-21`, `NFR-8`, `NFR-9`, `NFR-19` | Closed: activation or requested material update authorizes the plan-only commit and normal push; broader external actions remain separate. |
| `D-9` | How should existing plans be migrated? | implementation-shaping | approved | `D-1`-`D-3` | `FR-22`-`FR-25`, `FR-35`, `FR-36`, `NFR-5`, `NFR-14`, `NFR-20` | Closed: semantic inventory, optimistic hash checks, branch ownership, mixed-format cutover, ledger, timestamp preservation, and compatibility. |
| `D-10` | How should human-readable local and global summaries be materialized? | implementation-shaping | approved | `D-1`, `D-4`, `D-7` | `FR-8`, `FR-9`, `NFR-2`, `NFR-13` | Closed for V1: generate on demand; no feature-branch edits to a committed summary. |
| `D-11` | Which runtime owns first-version source refresh? | implementation-shaping | approved | `D-7` | `FR-13`-`FR-19`, `FR-27`-`FR-29`, `FR-34`, `NFR-7`, `NFR-11` | Closed: provider-neutral Go scanner with local Git and authenticated read-only GitHub-through-`gh` adapters. |
| `D-12` | How should pre-push intent visibility work? | non-blocking | approved | `D-7`, `D-8` | `NFR-9` | Closed for V1: no pre-push global visibility; reconsider only with evidence. |
| `D-13` | How are portfolio members enrolled and discovered? | implementation-shaping | approved | `D-4` | `FR-13`, `FR-14`, `FR-25`, `NFR-17` | Closed: nested v2 home/member identity and discovery sources, default-branch predicates, v1 compatibility, and no central repository list. |
| `D-14` | Which coordination capabilities belong in the Kit versus the coordination home? | implementation-shaping | approved | `D-11`, `D-13` | `FR-31`, `FR-32`, `NFR-18` | Closed: reusable concepts and mechanics are managed; portfolio instance state and strategy are repository-owned. |

## D-1 - Canonical Metadata Source

Options:

1. Keep the manually maintained central register as metadata authority.
2. Add one manually maintained sidecar entry beside each canonical plan.
3. Add versioned structured metadata to each canonical discovery and implementation document.
4. Maintain metadata only in a central database or coordination repository.

Decision: retain the existing planning header and H1, then place one explicitly bounded YAML
metadata block below the title. Require schema version, ID, kind, purpose, and catalog timestamps;
keep family, products, capabilities, themes, relations, and aliases optional. Derive repository,
path, ref, commit, checksum, and scan facts automatically. Do not require local sidecar entry
files.

Rationale: the canonical document already owns kind, lifecycle, dates, relations, and meaning.
Embedding a structured block removes local duplication while preserving independent files for
parallel branches. A central database alone would create dual authority, while sidecars would
reduce merge conflicts but retain drift risk.

Tradeoff: machine-readable metadata introduces a schema into human documents and requires
migration and validation. Retaining the established header means the parser reads the document
contract rather than treating the YAML block as a duplicate self-contained record.

State: approved by the user on 2026-07-31.

## D-2 - Stable Semantic Identity

Options:

1. Keep repository-local sequential numbers.
2. Allocate numbers centrally through an online coordinator.
3. Use repository-qualified semantic IDs.
4. Use opaque UUID or ULID values.
5. Use semantic IDs plus a short generated collision token.

Recommendation: start with `<repository>.<kind>.<stable-slug>` and preserve old sequential IDs as
aliases. Do not add a token until evidence shows that meaningful qualification and branch-conflict
detection are insufficient.

Rationale: semantic content helps conversation and strategic scanning. A same-ID choice on two
branches is often valuable evidence that the work overlaps. The coordinator should surface both
observations. If they are distinct plans, one receives a meaningful slug qualifier. Opaque or
random suffixes solve accidental collision at the cost of comprehension and can conceal duplicate
intent.

Normalization: lowercase ASCII; dots between structural segments; hyphens within segments;
repository namespace from configured stable member identity; explicit transliteration when
needed; filenames remain independent.

State: approved by the user on 2026-07-31.

## D-3 - Record Granularity

Options:

1. Keep one record per plan bundle, even when it combines discovery and implementation.
2. Create separate records for discovery and implementation, associate execution logs, and create
   first-class family records when warranted.
3. Treat every planning and evidence file as an independent catalog record.

Recommendation: choose option 2.

Rationale: discovery and implementation have distinct authority and lifecycle; execution logs are
evidence of implementation; families represent shared discoverability across independent sibling
plans. The model remains expressive without cataloging every attachment and log as strategy.

Family authority: preserve the current second-sibling trigger, use the family `README.md` as
canonical authority, and derive member lists from plan metadata.

State: approved, including the family-authority placement, by the user on 2026-07-31.

## D-4 - Complete Coordination Coverage

Decision: the coordination home should represent all formal planning records from every explicitly
enrolled repository, not only known cross-repository work.

Rationale: locally isolated capabilities may combine into one product or strategic opportunity.
Selective ingestion prevents the coordinator from discovering those relationships.

Boundary: complete coverage applies to formal catalog records, not every note, task, attachment, or
full plan body.

State: approved by explicit user direction.

## D-5 - Strategic Ownership Split

Options:

1. Member repositories own both factual metadata and all strategic interpretation.
2. The coordination home owns all catalog facts and may rewrite member classifications.
3. Member repositories own source facts; the coordination home owns additional portfolio
   families, themes, cross-repository relations, and analysis.

Recommendation: choose option 3.

Rationale: local owners know canonical scope and lifecycle; the coordinator sees portfolio
relationships unavailable locally. Separate ownership prevents source reconciliation from
overwriting strategic interpretation.

State: approved.

## D-6 - Lifecycle And Git State

Decision:

- keep the existing plan lifecycle values;
- retain `active` as current and approved for execution;
- do not add a redundant authority enum;
- represent Git visibility separately; and
- represent catalog verification separately.

Rationale: `active + unapproved` would be contradictory under existing lifecycle doctrine. An
active plan on a feature branch remains authorized even when unmerged.

State: approved through user correction during discovery.

## D-7 - First-Version Synchronization

Options:

1. Agent runbooks only.
2. Git-host events only.
3. Scheduled global polling only.
4. On-demand source refresh before every coordination analysis.
5. Immediate agent announcement plus durable retry, event-driven verification, and scheduled
   reconciliation.

Decision: choose option 4 for the first version.

Rationale: it removes dependence on agents remembering a second register update while avoiding the
hosting, authentication, retry, expiry, and cleanup burden of a proactive coordination system. It
is timely at the point where freshness matters: before strategic analysis. The scanner remains one
reusable authority that optional schedules or events may invoke later.

Boundary: pushed source is visible on the next scan. Uncommitted and unpushed work is intentionally
not visible globally in the first version.

State: approved.

## D-8 - Branch Publication Contract

Options:

1. Allow active implementation to remain local without a publication checkpoint by default.
2. Require every plan to merge before implementation.
3. Treat the user's activation request as scoped authority to commit and push the activated plan
   checkpoint to its work branch; keep broader external actions separate.

Decision: choose option 3.

Rationale: it makes canonical content remotely verifiable without confusing merge with execution
authority or requiring a redundant approval prompt. The authority includes the canonical plan and
directly required planning metadata on the unambiguous work branch. It excludes unrelated files,
unauthorized implementation code, pull-request creation, merge, release, force-push, destructive
branch changes, and ambiguous repositories or branches.

State: approved with the user-specified activation authority on 2026-07-31.

## D-9 - Semantic Migration

Options:

1. Mechanically translate legacy register entries.
2. Leave all legacy plans untouched and apply the new model only to future plans.
3. Enumerate all canonical plans, semantically evaluate each, reconcile legacy evidence, review
   ambiguities, then insert structured metadata while preserving content chronology.

Recommendation: choose option 3 for Codeheart-controlled repositories. Retain explicit,
non-destructive migration for unknown consumers.

Rationale: the existing set is bounded and valuable enough to inventory properly. A semantic agent
can recognize families, relationships, duplicates, status evidence, and strategic capabilities
that a row converter cannot.

Migration mechanics: use commit-and-hash preconditions, detect active-branch ownership, skip and
reevaluate changed plans, let branch owners migrate contested plans, avoid migrating inherited
copies, and use compatibility-tooling, semantic-inventory, authoring-cutover, and coverage-closure
phases.

State: approved by the user on 2026-07-31.

## D-10 - Generated Views

Options:

1. Continue manual shared Markdown summaries.
2. Commit generated summaries from every feature branch.
3. Generate summaries on demand or centrally after source integration; keep a stable entry point
   when a committed overview is not safe.
4. Replace all human-readable views with a database UI.

Recommendation: choose option 3.

Rationale: the summary remains useful without making unrelated feature branches edit the same
artifact. A database UI may become an additional view but should not be required for repository
portability.

State: approved for the first version. A stable explanatory `plan-register.md` may remain for
compatibility, but it is not manually updated as catalog authority.

## D-11 - Source Refresh Runtime

Candidate options:

- an Operating Kit CLI scan using local repositories;
- the same scan command with a provider adapter and the operator's existing authorized
  credentials;
- a later schedule that invokes the command; or
- a later event-triggered service that invokes the same contract.

Decision: make a provider-neutral Operating Kit Go command the first runtime. The
coordination-home runbook invokes it on demand before analysis. V1 includes a local Git adapter and
a GitHub adapter that uses the operator's existing authenticated `gh` session for read-only
discovery and source access. It stores no tokens and executes no branch code.

Branch scope: enumerate every accessible unmerged remote branch, ignore branches without canonical
plan changes, and use pull-request state only as enrichment.

State: approved by the user on 2026-07-31.

## D-12 - Pre-Push Intent Layer

Options:

1. Require an idempotent agent announcement for local work.
2. Add a distinct short-lived intent service.
3. Accept no global visibility until first push and use early publication when visibility matters.

Decision: choose option 3 for the first version.

Rationale: an announcement route still depends on agent discipline and adds a second state surface
whose authority, retry, expiry, and reconciliation must be maintained. Early authorized push
creates durable, reviewable source evidence and uses the normal repository workflow.

State: approved as a first-version deferral. Reconsider announcements only if observed work
requires proactive awareness before publication.

## D-13 - Automatic Portfolio Membership

Options:

1. Maintain a central list of member repositories.
2. Treat every repository visible to the coordinator as a member.
3. Preserve portfolio configuration and automatically discover matching Kit-enabled repositories
   within an explicitly authorized scope.

Decision: choose option 3.

Rationale: the existing configuration already expresses affiliation. Automatic discovery removes
the duplicate registry and its maintenance burden while explicit scope and matching configuration
prevent silent collection of unrelated repositories. Default-branch state controls membership, so
a feature branch cannot self-enroll.

Configuration: add a nested portfolio schema version, stable `coordination_home_id`, explicit
discovery sources for the home, and default-branch membership predicates; retain v1 path fields as
compatibility inputs. Unconfigured Kit repositories may appear only in a non-authoritative
candidate report.

State: approved by the user on 2026-07-31.

## D-14 - Kit Versus Coordination-Home Ownership

Decision:

- the Operating Kit owns the generic catalog model, roles, setup and maintenance routes, config
  schema, scanner, branch-overlay mechanics, validators, migration, and reusable scaffolds;
- member repositories own canonical plans and their source metadata; and
- coordination-home repositories own actual portfolio scope, rebuildable catalog state,
  strategic relationships, themes, priorities, and analysis.

Rationale: a fresh Kit-enabled repository can recreate the operating model without prior context,
while private or portfolio-specific state remains outside the public reusable core. A local wrapper
may add repository details but should not fork generic doctrine.

State: approved.

# Section 6 - Open Questions And Assumptions

## Resolved Questions

| ID | Resolution | Closure |
| --- | --- | --- |
| `OQ-1` | Lowercase ASCII `<repository>.<kind>.<stable-slug>` with explicit transliteration, stable configured repository namespace, and independent filenames. | Approved with `D-2` on 2026-07-31. |
| `OQ-2` | Existing header and H1 plus a marker-bounded YAML metadata block below the title using schema version 1 and the existing supported YAML subset. | Approved with `D-1` on 2026-07-31. |
| `OQ-3` | Enumerate every accessible unmerged remote branch; ingest only canonical plans added or materially changed relative to the merge base; PR data is optional enrichment. | Approved with `D-11` on 2026-07-31. |
| `OQ-4` | User-requested activation or material active-plan update authorizes the plan-only commit and normal push to the unambiguous work branch; broader external actions remain separate. | Approved with `D-8` on 2026-07-31. |
| `OQ-5` | Provider-neutral Go scanner with local Git and authenticated read-only GitHub-through-`gh` adapters; no token storage or branch-code execution. | Approved with `D-11` on 2026-07-31. |
| `OQ-7` | Required metadata core is schema version, ID, kind, purpose, and catalog timestamps; family and strategic classifications remain optional. | Approved with `D-1` on 2026-07-31. |
| `OQ-10` | Canonical family authority lives in the family `README.md` at the current second-sibling trigger; membership is derived from plan metadata. | Approved with `D-3` on 2026-07-31. |
| `OQ-12` | Nested portfolio schema v2 adds stable home identity and discovery sources while v1 paths remain readable compatibility inputs. | Approved with `D-13` on 2026-07-31. |

## Remaining Non-Blocking Questions

| ID | Question | Owner | BLOCKER | Affects | Resolution path |
| --- | --- | --- | --- | --- | --- |
| `OQ-6` | Should `plan-register.md` remain a stable explanatory entry point after migration, or be retired once the generated local view exists? | Operating Kit owner | no | `D-10`, `FR-8`, `NFR-2`, `NFR-13`, `NFR-14` | Decide during compatibility implementation after testing discoverability; either outcome preserves canonical metadata authority. |
| `OQ-8` | What retention presentation should stale, deleted, inaccessible, or out-of-scope branch observations use? | Portfolio operator | no | `FR-18`, `FR-19`, `NFR-10`, `NFR-14` | Implement deterministic retirement for merged or deleted refs and use fixtures to select a human-facing stale-retention window. |
| `OQ-9` | What scan freshness and duration targets are necessary for interactive strategic analysis? | Portfolio operator | no | `FR-15`, `FR-28`, `NFR-4`, `NFR-10` | Benchmark representative portfolios during implementation and set explicit limits without changing the scan model. |
| `OQ-11` | Which current register-only session refs and coordination notes remain valuable enough to migrate into canonical metadata or separate evidence? | Migration owner | no | `D-9`, `FR-23`, `FR-24`, `NFR-5` | Classify each legacy-only field in the semantic migration ledger; ambiguity does not block tooling or authoring cutover. |
| `OQ-13` | Which generated view should expose unconfigured Kit-enabled candidate repositories? | Portfolio owner | no | `D-13`, `FR-14` | Select a non-authoritative setup or scan report during UX implementation; candidates never become catalog members automatically. |

## Assumptions

| ID | Assumption | Validation |
| --- | --- | --- |
| `A-1` | The number of Codeheart-controlled canonical plans is small enough for a complete semantic review during one bounded implementation. | Produce a cross-repository inventory and estimate document volume before implementation planning closes. |
| `A-2` | Canonical planning filenames are consistent enough to enumerate most discovery and implementation documents automatically. | Scan enrolled repositories and report exceptions. |
| `A-3` | A complete metadata catalog is sufficient for routine global strategic orientation without copying full plan bodies. | Test strategic analysis against catalog-only fixtures, then fetch source for unresolved cases. |
| `A-4` | Repository-qualified semantic identifiers remain understandable enough for human conversation, and same-ID branch overlaps are uncommon but useful signals. | Review representative names, parallel-branch fixtures, and agent handoff examples. |
| `A-5` | Repositories can provide or explicitly configure a default or integration branch and portfolio affiliation. | Validate current and target config shapes. |
| `A-6` | An on-demand scanner can access relevant member refs with existing authorized local or provider credentials without executing untrusted branch code. | Complete permission and threat-model review for the selected adapter. |
| `A-7` | Current lifecycle values remain sufficient; a separate `paused` state is not required by this redesign. | Review semantic migration findings for real paused-but-valid plans. |
| `A-8` | Product and capability classifications can be public-safe or appropriately scoped in the coordination home. | Run public-core and repository-privacy classification during migration. |
| `A-9` | Portfolio size permits an on-demand refresh quickly enough for routine coordination analysis. | Benchmark default-baseline and branch-overlay scans across the intended scope. |
| `A-10` | Operating Kit installation markers and portfolio configuration are available on member default branches and can support discovery without a repository list. | Inventory representative member repositories and test the onboarding predicate. |

# Section 7 - Risks And Validation

## Risks

| ID | Risk | Likelihood | Impact | Mitigation | Detection |
| --- | --- | --- | --- | --- | --- |
| `R-1` | Semantic IDs remain too long or unstable for conversation. | medium | medium | Separate stable ID from required title-and-kind first reference; prohibit mutable status or ownership in IDs. | Human and agent reference review. |
| `R-2` | Two branches create the same semantic ID for genuinely distinct plans. | medium | medium | Preserve both observations, flag the overlap, and require a meaningful qualifier when review confirms distinct intent. | Concurrent-branch fixture. |
| `R-3` | Structured metadata duplicates or contradicts document content. | medium | high | Define canonical fields, derive what can be derived, and validate lifecycle/title/path consistency. | Schema and semantic validation. |
| `R-4` | Metadata migration falsely updates historical content dates. | medium | high | Preserve the existing `Last updated` value for catalog-only insertion and record separate catalog timestamps. | Before/after migration ledger. |
| `R-5` | Semantic migration invents families or relationships not supported by evidence. | medium | high | Require evidence, confidence, ambiguity reporting, and manual review of strategic classifications. | Migration review gate and relation audit. |
| `R-6` | Global completeness collects private or sensitive metadata. | low | high | Explicit enrollment, public-safe schemas, repository-scoped permissions, field allowlists, and public-core review. | Privacy and public-core validators. |
| `R-7` | Branch ingestion creates noisy or stale strategic records. | high | medium | Ingest only plan-changing merge-base diffs, retire merged or deleted refs, and expose old unmerged observations as stale. | Branch-overlay lifecycle tests. |
| `R-8` | A naïve branch scan duplicates every inherited default-branch plan. | high | high | Full-scan only the baseline and ingest added or materially changed canonical plans on relevant branches. | Multi-branch fixture with a large inherited plan set. |
| `R-9` | Scanner credentials expose more repository scope than intended or branch content executes during collection. | low | high | Use explicit read scope, preflight access, parse files as data, and never execute branch code. | Permission and adversarial-branch review. |
| `R-10` | On-demand synchronization is mistaken for proactive monitoring. | medium | medium | State the use-time freshness contract and disclose last scan; add scheduling later only if justified. | Low-context operator probe. |
| `R-11` | Activation authority is interpreted as permission to push unrelated or destructive changes. | medium | high | Define activation as plan-only normal-push authority; exclude unrelated files, broader code, PRs, merge, release, force-push, destructive actions, and ambiguous targets. | Publication-authority boundary tests. |
| `R-12` | Generated summaries return as merge hotspots. | medium | medium | Do not regenerate committed summaries on feature branches; centralize or generate on demand. | Concurrent branch merge test. |
| `R-13` | Global strategic overlays are overwritten by member reconciliation. | medium | high | Separate source-derived records from coordination-owned analysis paths and ownership. | Reconciliation preservation test. |
| `R-14` | Unknown consumers experience a forced migration. | low | high | Compatibility reader, explicit dry-run migration, no sync overwrite, and adoption notes. | Existing-consumer fixture. |
| `R-15` | Provider-specific GitHub mechanics leak into generic doctrine. | medium | medium | Define provider-neutral capability and route contracts with adapter-specific implementation docs. | Managed-doc public-core and portability review. |
| `R-16` | The coordinator cannot see uncommitted or unpushed work and is assumed to be complete. | high | medium | State the visibility boundary and make early authorized publication part of the activation route. | Local-unpublished scenario test. |
| `R-17` | Automatic repository discovery silently enrolls unrelated repositories. | low | high | Require authorized scope, default-branch Kit marker, matching portfolio configuration, and explicit candidate-only treatment for unconfigured repositories. | Scope and false-positive discovery fixtures. |
| `R-18` | Rich metadata becomes high-maintenance taxonomy governance. | medium | medium | Require only a stable core, keep strategic fields optional, and defer controlled vocabularies until migration evidence supports them. | Schema usability and migration review. |
| `R-19` | A fresh repository cannot reconstruct the intended coordination-home model from the Kit alone. | medium | high | Ship discoverable managed concepts, setup runbook, schemas, commands, validation, and a low-context acceptance test. | Fresh-repository bootstrap probe. |
| `R-20` | Portfolio scans become too slow for on-demand analysis. | medium | medium | Use safe baseline caching, diff branch overlays before parsing content, and expose stale fallback explicitly. | Representative performance benchmark. |
| `R-21` | Migration collides with active plan edits or later creates header conflicts when branches merge. | high | high | Record commit/hash and active-branch touch sets, skip changed or contested files, and assign migration to branch owners where appropriate. | Concurrent-edit and later-merge fixtures. |
| `R-22` | Portfolio schema v2 breaks working v1 path-based configurations or older consumers. | medium | high | Version the nested portfolio contract, retain v1 reading and fields, provide dry-run migration, and do not rewrite unknown consumer config during ordinary sync. | v1/v2 compatibility fixtures. |

## Targeted Validation Plan

Future implementation planning should cover at least:

- schema validation for discovery, implementation, and family metadata;
- marker-bounded YAML extraction without replacing the existing lifecycle header or H1;
- required-core and optional-rich-field validation using the supported YAML subset;
- portable ID generation and normalization tests on macOS and Windows;
- simultaneous creation of different plans on parallel branches;
- simultaneous modification of the same plan on parallel branches;
- duplicate semantic ID overlap and meaningful-qualifier handling;
- discovery-to-implementation and family relation validation;
- lifecycle, Git visibility, and verification invariant tests;
- automatic member discovery from default-branch Kit markers and matching portfolio configuration;
- v1 path-based and v2 identity-and-scope portfolio configuration compatibility;
- proof that no central repository list is required;
- exclusion or candidate-only reporting for unconfigured and mismatched repositories;
- idempotent repeated on-demand scan behavior;
- full default-branch baseline ingestion;
- merge-base branch overlays containing only added or materially changed plan documents;
- enumeration of plan-changing unmerged remote branches without requiring prefixes, pull requests,
  or recency;
- proof that inherited plans are not duplicated across branches;
- pushed-branch ingestion before merge, with or without a pull request;
- branch merge, deletion, and inaccessible-ref cleanup;
- multiple source-ref observations for conflicting changes to one plan;
- source commit and freshness evidence checks;
- preservation of coordination-owned strategic overlays during source reconciliation;
- scan preflight and failure disclosure before strategic analysis;
- local Git and authenticated read-only GitHub-through-`gh` adapter contract tests;
- no execution of untrusted branch code;
- migration dry-run inventory of every canonical plan;
- semantic migration ledger completeness and legacy alias reconciliation;
- preservation of original content-update dates;
- commit/hash mismatch, active-branch ownership, skip, reevaluation, and branch-owner migration
  tests;
- mixed-format authoring cutover without a long manual register and metadata dual-write period;
- existing consumer-owned legacy register preservation;
- compatibility for existing portfolio path fields;
- public-core and privacy validation of catalog fields;
- failure-closed permission tests for private or unavailable repositories;
- generated-summary merge-conflict proof; and
- activation-authorized plan-only normal push plus refusal of unrelated, ambiguous, destructive,
  PR, merge, release, and force-push actions;
- low-context agent probes that reference plans by title, kind, family, and stable ID;
- fresh-repository coordination-home setup using only managed Operating Kit guidance; and
- optional scheduling proof that invokes the same scanner rather than creating a second catalog
  authority.

## Preliminary Consumer Impact

Future implementation is likely to span:

- `instruction-only change` for managed planning, coordination, publication, and migration
  doctrine;
- `validator-only change` for metadata, identity, relation, lifecycle, and catalog validation;
- `consumer migration required` for consumers that opt into canonical structured metadata and
  generated catalogs;
- `backwards-compatible scaffold addition` if new coordination-home config or repository-owned
  portfolio surfaces are created only when absent;
- `breaking placement-contract change` if the current plan-register state-file contract is
  replaced or moved rather than preserved through compatibility;
- `security or safety policy change` for configured coordination read scope and activation-based
  plan-only commit-and-push authority; and
- possible generated-path or component changes that require explicit placement and manifest
  review.

No implementation plan should collapse these into one impact class.

# Section 8 - Implementation Handoff

## Approval Record

The user approved the complete recommendation set on 2026-07-31, with one correction that is now
part of `D-8`: requesting plan activation is itself sufficient scoped authority to commit and push
the plan checkpoint. No second push-approval prompt is required.

The approved set includes:

- marker-bounded YAML metadata below the existing title and lifecycle header;
- the required metadata core and optional rich strategic fields;
- exact semantic ID normalization and independent filename behavior;
- separate discovery, implementation, and family records with family README authority;
- on-demand generated local and portfolio views;
- automatic config-driven membership without a repository list;
- default-branch baselines and all accessible plan-changing unmerged remote branch overlays;
- provider-neutral Go scanning with local Git and GitHub-through-`gh` adapters;
- activation-based plan-only publication authority;
- branch-aware optimistic semantic migration without a repository-wide freeze; and
- managed coordination-home setup with repository-owned portfolio instance data.

## Reviewer State

The implementation-shaping decisions received a documented main-thread evidence and option pass
against current lifecycle doctrine, portfolio schema, family placement, YAML support, and the
self-contained Go CLI architecture. The user reviewed the concrete recommendations, corrected the
publication authority boundary, and approved the revised set. No unresolved reviewer dissent
remains.

## Frozen Inputs

- Canonical plan metadata, not generated catalog output, is the source of truth.
- Existing header and H1 authority are preserved.
- IDs use lowercase ASCII `<repository>.<kind>.<stable-slug>` without an initial random token.
- Portfolio configuration evolves through a nested v2 contract while v1 remains readable.
- Membership comes only from authorized scope plus matching default-branch configuration.
- Coordination analysis refreshes on demand before claiming a current overview.
- Branch overlays include only canonical plans added or materially changed from the merge base.
- V1 scans all accessible unmerged remote branches rather than depending on naming, PR, or recency
  discipline.
- Activation authorizes a plan-only normal commit and push on the unambiguous work branch.
- Migration uses semantic review, optimistic version checks, active-branch ownership, a ledger,
  mixed-format authoring cutover, and final coverage reconciliation.
- Announcements, outboxes, webhooks, GitHub Apps, background services, and required schedules are
  outside V1.

## Remaining Non-Blocking Follow-Ups

- `OQ-6`: final stable-entry-point treatment for `plan-register.md`.
- `OQ-8`: human-facing stale branch-observation retention.
- `OQ-9`: measured scan duration and freshness targets.
- `OQ-11`: legacy-only session-ref and coordination-note classification.
- `OQ-13`: candidate-repository report placement.

These questions may be resolved from implementation evidence without changing the approved
capability or architecture.

## Readiness

Current status: `implementation-handoff-ready`.

All implementation-shaping decisions are approved, no `BLOCKER: yes` question remains, and the
capability scopes below provide a single coherent planning path. No further discovery decision is
required before drafting the implementation plan. A later user action is still required to approve
execution of that draft plan.

## Implementation Capability Scope - Canonical Plan Metadata, Identity, And Families

Capability:
Every formal discovery, implementation, and first-class family record can identify and describe
itself through a stable machine-readable contract that remains understandable to humans and safe
under parallel branches.

Primary workflow:
An agent creates or updates a canonical plan using the existing lifecycle header and title, writes
or validates the bounded metadata block, assigns the stable semantic ID, relates discovery,
implementation, and family records, and receives deterministic validation before publication.

Must cover:

- explicit metadata begin/end markers and a visible YAML block below the H1;
- existing header and H1 parsing as canonical title, created, status, and content-update sources;
- required schema version, ID, kind, purpose, and catalog timestamp fields;
- optional family, products, capabilities, strategic themes, relations, and legacy aliases;
- lowercase ASCII repository, kind, and stable-slug normalization;
- repository namespace from stable configured member identity;
- ID stability across title, path, repository-name, status, and classification changes;
- separate discovery, implementation, and family records;
- family README authority at the second-sibling trigger with derived membership;
- duplicate same-ID branch observation rather than scan-order overwrite;
- schema, identity, relation, family, header, and ownership validation; and
- managed references, templates, runbook hooks, packaged mirrors, and migration-safe examples.

Explicitly out of scope:

- mandatory global product, capability, or strategic-theme vocabularies;
- opaque UUIDs, central sequential allocation, or random ID suffixes in V1;
- independent strategic records for execution logs, checkboxes, transcripts, or attachments; and
- duplicating scanner-derived repository, path, ref, commit, checksum, pull-request, or scan facts
  in canonical metadata.

Deferred or blocked:

- none inside the approved metadata and identity capability.

Preserve decisions:

- `D-1`, `D-2`, `D-3`, `D-5`, and `D-6`.

Planner must not reinvent:

- metadata placement, marker contract, required core, optional rich fields, header authority, ID
  grammar, filename independence, family authority, or legacy alias preservation;
- the rule that catalog-only migration does not rewrite historical `Last updated`; and
- the rule that same-ID branches remain separate observations until human or source integration
  resolves them.

Feature-level success evidence:

- fixtures prove valid and invalid metadata across discovery, implementation, and family records;
- portable ID tests pass on supported macOS and Windows paths;
- parallel-branch fixtures preserve separate same-ID and different-ID observations; and
- a low-context agent can identify a plan by title, kind, ID, family, and canonical path without
  relying on a register number or private conversation.

## Implementation Capability Scope - Local Views And Compatibility Migration

Capability:
Repositories can derive useful local plan views from canonical documents, adopt the new metadata
contract while ordinary development continues, and retire manual register maintenance without
losing legacy evidence or historical chronology.

Primary workflow:
The tooling reads new metadata when present, falls back to canonical content and the legacy
register when absent, inventories every plan semantically, migrates only source-stable files,
records contested work in a ledger, begins metadata-only authoring, and closes coverage through a
final reconciliation.

Must cover:

- generated local listing without feature-branch edits to a shared summary;
- compatibility reading for legacy registers and older branches;
- automated plan enumeration followed by semantic agent review of each canonical plan;
- legacy alias, session-reference, relation, lifecycle, and coordination-note reconciliation;
- source commit and content-hash preconditions before metadata insertion;
- active-branch touch detection and branch-owner assignment for contested plans;
- skip and reevaluation behavior when a source version changes;
- a migration ledger with confidence, ambiguity, aliases, conflicts, deferrals, and ownership;
- separate compatibility-tooling, semantic-inventory, authoring-cutover, and coverage-closure
  phases;
- new-plan metadata enforcement after authoring cutover;
- cessation of manual register writes once mixed-format generated views are reliable;
- final default-branch and active-overlay coverage reconciliation; and
- dry-run, idempotency, preservation, rollback, and unknown-consumer compatibility behavior.

Explicitly out of scope:

- a repository-wide development freeze;
- blind row conversion from legacy registers;
- automatic overwrite of a changed or branch-contested plan;
- migrating unchanged inherited copies separately on every branch; and
- forced migration during ordinary Operating Kit sync, repair, or upgrade.

Deferred or blocked:

- final stable-entry-point treatment for `plan-register.md`: non-blocking `OQ-6`;
- legacy-only evidence classification details: non-blocking `OQ-11`.

Preserve decisions:

- `D-1`, `D-9`, and `D-10`.

Planner must not reinvent:

- semantic rather than mechanical migration;
- commit/hash optimistic concurrency, branch ownership, and ledger requirements;
- the mixed-format authoring cutover that avoids a long dual-write period; or
- content chronology and legacy alias preservation.

Feature-level success evidence:

- a representative migration dry-run inventories every canonical plan and reconciles the legacy
  register;
- concurrent-edit and active-branch fixtures prove that migration skips rather than overwrites;
- generated local views remain correct across mixed new and legacy records; and
- authoring cutover stops manual shared-register edits while unresolved historical records remain
  discoverable through compatibility reading.

## Implementation Capability Scope - Portfolio Discovery, Scanning, And Coordination Home

Capability:
A freshly configured coordination home can automatically discover matching Kit-enabled member
repositories, refresh a complete factual plan catalog including active remote-branch work, and
support separate durable strategic analysis without a central repository list.

Primary workflow:
A user asks a low-context agent to configure a coordination home; the managed route inspects the
repository, selects authorized local or provider scope, writes repo-owned portfolio configuration
and surfaces, performs preflight, scans default baselines and changed-plan branch overlays, reports
freshness and access failures, and makes the refreshed catalog available for strategic analysis.

Must cover:

- nested portfolio schema v2 for stable home identity, stable member identity, and discovery
  sources;
- v1 `coordination_home_path` and `coordination_home_register_path` compatibility;
- machine-specific local scope in `.codeheart/local/`;
- authorized-scope, default-branch Kit marker, member-role, repository-ID, and matching-home-ID
  membership predicates;
- candidate-only treatment for unconfigured Kit-enabled repositories;
- no manually maintained repository list;
- a provider-neutral scanner implemented in the self-contained Go CLI;
- local Git and authenticated read-only GitHub-through-`gh` adapters;
- authentication and access preflight without credential storage or branch-code execution;
- complete default-branch baseline scans;
- enumeration of every accessible unmerged remote branch;
- merge-base path filtering and ingestion of only added or materially changed canonical plans;
- independent conflicting branch observations, merged/deleted retirement, stale labeling, source
  revisions, and scan errors;
- on-demand refresh before current strategic analysis;
- separate member-derived factual observations and coordination-owned strategic overlays;
- managed setup and maintenance routes, config schemas, validators, scanner mechanics, routing
  hooks, packaged mirrors, and fresh-repository scaffolding; and
- a later trigger interface that can invoke the same scanner without creating another authority.

Explicitly out of scope:

- pre-push announcements, local outboxes, webhooks, GitHub Apps, background services, or mandatory
  schedules in V1;
- implicit enrollment based on repository names, folder adjacency, or provider visibility alone;
- executing scanned branch code;
- copying complete plan bodies into the portfolio catalog by default; and
- storing private portfolio strategy in generic public managed doctrine.

Deferred or blocked:

- stale-observation presentation and retention: non-blocking `OQ-8`;
- measured scan duration and freshness thresholds: non-blocking `OQ-9`;
- candidate-report presentation: non-blocking `OQ-13`.

Preserve decisions:

- `D-4`, `D-5`, `D-7`, `D-11`, `D-13`, and `D-14`.

Planner must not reinvent:

- the v2 identity-and-scope configuration direction or v1 compatibility;
- the membership predicate;
- the default-baseline plus plan-changing unmerged-branch overlay model;
- the V1 local Git and GitHub-through-`gh` adapters;
- the on-demand rather than proactive synchronization boundary; or
- the ownership split between managed mechanics, member source facts, and coordination-home
  strategic instance data.

Feature-level success evidence:

- a fresh low-context agent can configure and validate a coordination home from installed managed
  guidance alone;
- adding a matching configured member requires no central-list edit;
- fixtures exclude mismatched and unconfigured repositories from catalog membership;
- baseline and multi-branch fixtures prove completeness without inherited-plan duplication;
- local and GitHub adapter tests prove idempotency, read-only credential handling, and failure
  disclosure; and
- reconciliation never overwrites coordination-owned strategic relationships or analyses.

## Implementation Capability Scope - Activation Publication And Planning Workflow

Capability:
Activating or materially updating a formal plan makes its canonical planning checkpoint remotely
visible without a redundant push-approval prompt or accidental authority for broader external
actions.

Primary workflow:
The user requests activation or a material update, the agent confirms the unambiguous work branch,
updates the canonical plan, commits the plan and directly required planning metadata, performs a
normal push, reports any blocker or visibility limitation, and continues only within the
separately authorized implementation scope.

Must cover:

- planning lifecycle and runbook doctrine that defines activation as plan-only commit-and-push
  authority;
- dedicated work-branch confirmation or creation under repository policy;
- the canonical plan and directly required planning metadata checkpoint;
- the same scoped authority for a user-requested material active-plan update;
- no second push-approval prompt;
- explicit exclusions for unrelated files, unauthorized implementation code, pull-request
  creation, merge, release, force-push, destructive branch actions, and ambiguous targets;
- stop or blocker behavior for authentication, repository policy, rejected normal push, dirty
  unrelated changes, or ambiguous branch/repository scope;
- a clear global-visibility limitation when the checkpoint cannot be published;
- no requirement to merge before authorized implementation begins; and
- managed discovery, implementation-planning, execution, review, and register-transition hooks
  with proportional safety tests.

Explicitly out of scope:

- standing authority to push arbitrary repository changes;
- automatic pull-request creation, merge, release, branch deletion, or force-push;
- pre-push global intent records; and
- treating Git integration state as plan lifecycle or execution authority.

Deferred or blocked:

- none inside the approved activation-publication capability.

Preserve decisions:

- `D-6`, `D-8`, and `D-12`.

Planner must not reinvent:

- activation as sufficient plan-only normal-push authority;
- the exact exclusions and ambiguity stop conditions;
- the separation between lifecycle, Git visibility, and catalog freshness; or
- the rule that merge is not a prerequisite for authorized execution.

Feature-level success evidence:

- scenario tests prove activation commits and pushes only the planning checkpoint without a second
  approval prompt;
- negative tests refuse unrelated, ambiguous, destructive, PR, merge, release, and force-push
  actions;
- an active pushed plan appears in the next branch-aware coordination scan; and
- low-context agent probes explain both the authority granted by activation and its boundaries.

# Revision Notes

- 2026-07-30: Created the draft discovery from the plan-register identity, merge-conflict,
  document-kind, plan-family, complete portfolio coverage, structured metadata, branch visibility,
  synchronization reliability, publication checkpoint, and semantic migration discussion.
- 2026-07-31: Simplified the recommended operating model after further review: retained rich
  canonical metadata and existing portfolio configuration; selected automatic config-driven member
  discovery without a repository list; changed first-version synchronization to an on-demand scan
  before analysis; defined default-branch baselines plus changed-plan branch overlays; deferred
  announcements, outboxes, webhooks, apps, and scheduled infrastructure; selected semantic IDs
  without generated tokens initially; and required the Operating Kit to support complete
  low-context coordination-home setup.
- 2026-07-31: Recorded user approval of the complete implementation-shaping recommendation set;
  fixed the bounded YAML metadata, exact ID, family authority, all-unmerged-branch diff, Go
  scanner, local Git and GitHub-through-`gh`, nested portfolio v2, optimistic parallel migration,
  and mixed-format cutover contracts; defined activation as sufficient plan-only normal
  commit-and-push authority; closed every blocking question; and added four implementation
  capability scopes to make the discovery implementation-handoff-ready.
