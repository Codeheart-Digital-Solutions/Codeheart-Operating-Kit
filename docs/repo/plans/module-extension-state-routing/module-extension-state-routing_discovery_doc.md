Last updated: 2026-06-23T13:41:25Z (UTC)
Created: 2026-06-23
Status: draft

# Module Extension State Routing Discovery

## Overview

This discovery defines whether Codeheart Operating Kit should own a generic placement and routing
rule for committed, non-secret state used by installed modules or extensions.

The immediate trigger came from a Foundry Microsoft 365 workspace module discussion. The module
needs repo-specific routing state such as a workspace registry and tenant onboarding facts. That
state is specific to the consumer repository, should be visible before Microsoft sign-in, and
should be committed when it is non-secret and shared by the repo's agents. The reusable question is
broader than M365: future module or extension systems may need similar committed routing state for
finance operations, CRM sync, document automation, shared services, AWS automation, or other
consumer-local operating surfaces.

Recommended direction: use `docs/repo/state/<module-or-extension-id>/` for committed,
non-secret, consumer-owned module or extension routing state. Operating Kit should own the generic
placement and routing doctrine. Module systems should own how installed modules are discovered.
Individual modules should own their concrete state files and schemas.

This is public Operating Kit discovery. It must not include private tenant details, account
identifiers, credentials, local machine state, or raw operational logs.

## Goals

- Define a generic committed-state location for installed modules and extensions.
- Keep the rule visible enough that agents check module state before asking repeated setup
  questions or touching external systems.
- Separate Operating Kit structure doctrine from Foundry-specific module metadata.
- Clarify what belongs in committed repo state versus local ignored state, managed snapshots,
  generated lockfiles, run records, and live external systems.
- Avoid over-designing a full plugin framework before repeated usage proves it is needed.
- Produce an implementation handoff for a later structure-governance update.

## Non-Goals

- Do not implement the rule in this discovery.
- Do not create or move any consumer state files.
- Do not define M365-specific schemas such as workspace registry contents.
- Do not define a universal plugin manifest, marketplace, dependency model, or compatibility
  framework.
- Do not make Operating Kit own Foundry lockfile or module manifest formats.
- Do not scaffold empty state folders in every consumer repository.
- Do not treat committed state as live external truth or authorization for external changes.
- Do not store secrets, tokens, credentials, sensitive personal data, account identifiers, or raw
  operational logs in committed module state.

## Context And Evidence

| Evidence | Current Signal | Discovery Implication |
| --- | --- | --- |
| Operating Kit placement contract | `docs/repo/` is consumer-owned repository documentation and local rules; `.codeheart/kit/` is managed; `.codeheart/user/` is local ignored state. | A committed repo-owned state path belongs under `docs/repo/`, not managed kit content. |
| Structure governance docs | Existing references distinguish `plans/`, `runbooks/`, `reference/`, `archive/`, managed content, scaffolds, templates, local-user, generated, and reports. | `state/` is a missing artifact category for committed routing context. |
| Foundry installed snapshot | Foundry exposes installed modules through `.codeheart/foundry/foundry.lock.yaml` and module metadata through `.codeheart/foundry/modules/<module-id>/module.yaml`. | Operating Kit should not define Foundry internals; it can use the discovered module ID as a namespace. |
| M365 workspace module need | Workspace routing needs local, committed, non-secret files such as `workspace-registry.yaml` and `tenant-onboarding.yaml`. | `docs/repo/state/m365-workspace/` is a concrete first use case. |
| Runbook-authoring discovery | Human-facing and hybrid runbooks should check visible local context before asking repeated setup questions. | Module state must be routed visibly enough for agents to find it. |

Example target shape:

```text
docs/repo/state/
  m365-workspace/
    workspace-registry.yaml
    tenant-onboarding.yaml
  finance-ops/
    account-registry.yaml
  crm-sync/
    sync-targets.yaml
  document-automation/
    publishing-registry.yaml
```

## Terminology

`module or extension`: a reusable installed capability surface that may provide runbooks,
references, templates, schemas, or operational behavior. Foundry modules are the current concrete
example. Future shared services or product modules may also fit.

`module system`: the owning system that installs or exposes modules and defines its registry,
lockfile, manifest, and lifecycle format. Foundry is a module system.

`module or extension ID`: the stable identifier exposed by the module system or module manifest.
It becomes the default namespace under `docs/repo/state/`.

`committed module state`: non-secret, consumer-owned routing state intended to be committed and
shared with repo agents.

`live external truth`: the actual state of an external system such as Microsoft 365, AWS, a CRM,
or a finance platform. Committed module state helps route the agent, but live preflight validates
before sensitive reads or changes.

## Decision Inventory

### D-001 - Use `docs/repo/state/<module-or-extension-id>/`

Status: recommended

Decision: committed, non-secret, consumer-owned routing state for installed modules and extensions
should live under:

```text
docs/repo/state/<module-or-extension-id>/
```

Why it matters: agents need a visible, repository-owned place to find module state before asking
setup questions, signing in to external systems, or touching external state.

Recommended default: use the module or extension ID as the folder name. For the current M365
module, the state root would be:

```text
docs/repo/state/m365-workspace/
```

### D-002 - Split Ownership Between Operating Kit, Module Systems, And Modules

Status: recommended

Decision:

- Operating Kit owns the generic placement and routing doctrine.
- Module systems own how installed modules are discovered.
- Modules own concrete state file names, schemas, lifecycle, and validation rules.
- Consumer repositories own the committed state contents after creation.

Why it matters: this lets Operating Kit provide a stable convention without coupling itself to
Foundry internals or future module systems.

Recommended default: Operating Kit documentation should describe "when a module or extension ID is
known" rather than hardcoding Foundry-specific discovery as the only mechanism.

### D-003 - Keep State Committed, Non-Secret, And Routing-Oriented

Status: recommended

Decision: `docs/repo/state/<id>/` is for committed routing context, selected targets, registries,
and lifecycle facts. It is not for secrets, token caches, generated lockfiles, raw run records, or
external-system truth.

Why it matters: the path must be safe to commit and useful for agents. If the path starts
collecting sensitive or transient data, it becomes unsafe and unreliable.

Recommended default: module state files should include only non-secret fields that help route an
agent. External changes still require live preflight and approval.

### D-004 - Make Routing Visible Without Listing Every Module

Status: recommended

Decision: the managed root `AGENTS.md` block and kit fallback inventory should include one generic
module/extension state route, not a line per module.

Recommended managed route shape:

```text
Before operating installed modules or extensions, inspect the module or extension registry exposed
by that system. When a module or extension ID is known, check
docs/repo/state/<module-or-extension-id>/ for committed non-secret routing state before asking
repeated setup questions or touching external systems.
```

Why it matters: Foundry modules may become a main local interaction surface, but root routing must
stay stable as new modules are installed.

### D-005 - Add Structure-Governance Doctrine, Not A Full Plugin Model

Status: recommended

Decision: the next implementation should add structure-governance doctrine for module/extension
state. It should not introduce a broad Operating Kit plugin architecture.

Why it matters: a plugin model would imply manifests, dependency resolution, compatibility
contracts, distribution semantics, and capability negotiation. The current evidence supports a
placement/routing convention, not a full plugin framework.

Recommended default: use neutral wording such as `module or extension`, and defer `plugin`
terminology until there is repeated evidence across module systems.

### D-006 - Do Not Scaffold Empty State Folders By Default

Status: recommended

Decision: Operating Kit should not create empty `docs/repo/state/` folders in every consumer
repository by default.

Why it matters: empty state folders create noise and may imply module state exists when it does
not.

Recommended default: create `docs/repo/state/<id>/` when a module or extension actually needs
committed state. A `docs/repo/state/README.md` may be useful when the first state namespace is
created, but should not be mandatory in every repo.

### D-007 - Require Live Preflight Before Sensitive Reads Or Changes

Status: recommended

Decision: committed module state routes the agent. It does not authorize changes or prove live
external state.

Why it matters: a local registry can become stale. Sensitive reads, writes, tenant changes,
account changes, and external resource changes must still validate against the live external
system and collect approval when required.

Recommended default: module runbooks should say how their state files are validated against live
systems before action.

## Requirements And Evaluation Criteria

### FR-001 - Standard Path

Operating Kit structure governance must define `docs/repo/state/<module-or-extension-id>/` as the
standard location for committed, non-secret, consumer-owned module or extension routing state.

### FR-002 - Visible Agent Route

Managed agent routes must make the state convention visible enough that agents check it before
repeating setup questions or touching external systems.

### FR-003 - Ownership Boundary

The standard must clearly separate:

- Operating Kit generic placement and routing doctrine;
- module-system discovery and lifecycle metadata;
- module-owned state file schemas;
- consumer-owned committed state contents;
- live external truth.

### FR-004 - Placement Procedure

The structure-governance placement runbook must route committed module/extension state to
`docs/repo/state/<id>/` when the artifact is shared, non-secret, repo-owned routing context.

### FR-005 - No Per-Module Root Route

The root managed route must stay generic. It must not require one `AGENTS.md` line per installed
module.

### NFR-001 - Public-Core Safety

Generic doctrine and examples must use placeholders and must not include real tenant, customer,
account, credential, local machine, or raw operational details.

### NFR-002 - Low Coupling

Operating Kit must not depend on Foundry-specific lockfile fields as the only module discovery
model.

### NFR-003 - Low Ceremony

The first implementation should be instruction-only unless scaffolding is explicitly approved.

## Placement And Implementation Surface

Recommended implementation files:

- `components/structure-governance/managed/reference/module-extension-state.md`: new durable
  reference for the state convention.
- `components/structure-governance/managed/README.md`: route to the new reference.
- `components/structure-governance/managed/reference/documentation-structure.md`: add
  `docs/repo/state/` to the reusable documentation organization model.
- `components/structure-governance/managed/reference/managed-content-boundaries.md`: classify
  `docs/repo/state/` as consumer-owned committed state.
- `components/structure-governance/managed/runbooks/change-documentation-placement.md`: route
  committed module/extension state to `docs/repo/state/<id>/`.
- `components/agent-interface/managed/kit-readme.md`: add fallback route to the state convention.
- `templates/agents/AGENTS.managed-block.md`: add one generic route/rule for module and extension
  state.
- `docs/repo/reference/placement-contract.md`: add the new consumer-owned path category if the
  implementation changes the durable placement contract.

If packaged resources mirror managed docs, the implementation must update those mirrors and
packaging tests.

Expected consumer impact class for a no-scaffold implementation: `instruction-only change`.

If the implementation scaffolds `docs/repo/state/README.md`, impact class becomes
`backwards-compatible scaffold addition` and needs preservation tests.

## Open Questions

### OQ-001 - Should The Term Be Module, Extension, Or Plugin?

BLOCKER: no

Owner: implementation planner

Recommended default: use `module or extension`. Avoid `plugin` until repeated usage proves a
formal plugin system is warranted.

### OQ-002 - Should Operating Kit Scaffold `docs/repo/state/README.md`?

BLOCKER: no

Owner: implementation planner

Recommended default: no for the first implementation. Define the convention and let module
installers create the state root when needed.

### OQ-003 - Should Operating Kit Validate Module State Files?

BLOCKER: no

Owner: module owners and later implementation planners

Recommended default: no generic validator in the first implementation. Modules own their concrete
schemas and validation.

### OQ-004 - Should This Discovery Merge With Runbook Authoring Standards?

BLOCKER: no

Owner: implementation planner

Recommended default: no. Keep it related but separate. Runbook authoring cares that agents check
visible local context; structure governance owns where committed module state lives.

## Risks

### R-001 - Sensitive Data Is Committed

Likelihood: medium

Impact: high

Mitigation: doctrine must explicitly ban secrets, credentials, account identifiers, raw logs, and
sensitive external state in `docs/repo/state/`. Future validators can inspect common anti-patterns
if repeated use warrants automation.

### R-002 - Stale State Is Treated As Live Truth

Likelihood: medium

Impact: high

Mitigation: every module using committed state must validate live systems before sensitive reads or
changes.

### R-003 - Routing Remains Hidden

Likelihood: medium

Impact: medium

Mitigation: add a generic route to the managed root `AGENTS.md` block and kit fallback inventory.

### R-004 - Operating Kit Overbuilds A Plugin Framework

Likelihood: medium

Impact: medium

Mitigation: keep the first implementation to placement, ownership, routing, and safety. Defer
plugin architecture until there is repeated evidence across multiple module systems.

## Implementation Scope Handoff

Next useful step: draft an implementation plan for an Operating Kit structure-governance
instruction-only release.

### Goal

Add Operating Kit doctrine so agents know where to find committed, non-secret module or extension
state before asking repeated setup questions or touching external systems.

Core convention:

```text
docs/repo/state/<module-or-extension-id>/
```

Example:

```text
docs/repo/state/m365-workspace/workspace-registry.yaml
```

### Primary Outcome

After the Operating Kit update is released and synced, agents should understand:

- installed modules/extensions may expose their own module IDs through their own registry or
  manifest;
- once a module ID is known, committed repo-owned state belongs under `docs/repo/state/<id>/`;
- this state is routing/context state, not live external truth;
- live external preflight is still required before sensitive reads or external changes;
- root `AGENTS.md` should route generally, not list every module.

### In Scope

Create a managed structure-governance reference:

```text
components/structure-governance/managed/reference/module-extension-state.md
```

Update managed structure governance:

```text
components/structure-governance/managed/README.md
components/structure-governance/managed/reference/documentation-structure.md
components/structure-governance/managed/reference/managed-content-boundaries.md
components/structure-governance/managed/runbooks/change-documentation-placement.md
```

Update generic agent routing:

```text
templates/agents/AGENTS.managed-block.md
components/agent-interface/managed/kit-readme.md
```

Update source repo placement contract if the implementation treats this as a durable consumer path
category:

```text
docs/repo/reference/placement-contract.md
```

Update component manifests and packaged resource mirrors as required:

```text
components/structure-governance/component.yaml
src/codeheart_operating_kit/resources/...
tests/...
```

Update release notes and consumer impact records.

### Out Of Scope

- Do not create `docs/repo/state/` in every consumer repo.
- Do not create M365-specific files.
- Do not modify Foundry installer behavior.
- Do not define Foundry lockfile or module manifest formats.
- Do not implement a plugin framework.
- Do not build validators for module state files.
- Do not migrate existing consumer/module state.

### Key Decisions To Encode

- Use `module or extension`, not `plugin`.
- Operating Kit owns generic placement/routing doctrine.
- Module systems own module discovery.
- Modules own concrete state files and schemas.
- Consumers own committed state contents.
- `docs/repo/state/<id>/` is for committed, non-secret routing state only.
- Local state helps route the agent; it does not authorize action.

### Acceptance Criteria

- Root managed `AGENTS.md` includes a generic module/extension state rule.
- Kit fallback inventory routes to the state convention.
- Structure governance explains what belongs in `docs/repo/state/<id>/`.
- Placement runbook routes committed module/extension state to that path.
- No per-module route is added to root `AGENTS.md`.
- No scaffolding is added unless explicitly approved.
- Packaging/resource tests pass after adding the managed reference.

### Validation

Run at minimum:

```sh
git diff --check
python3 scripts/validate-markdown-headers.py
python3 scripts/validate-public-core.py
python3 scripts/validate-json-schemas.py
pytest
```

Expected impact class: `instruction-only change` unless scaffolding is explicitly added.

## Revision Notes

- 2026-06-23: Initial discovery drafted from module-state routing discussion and current Operating
  Kit structure-governance boundaries.
- 2026-06-23: Added implementation scope handoff for a structure-governance instruction-only
  release.
