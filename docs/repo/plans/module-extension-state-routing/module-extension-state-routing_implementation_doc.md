Last updated: 2026-06-23T18:17:47Z (UTC)
Created: 2026-06-23
Status: active

# Document Header

## Module Extension State Routing Implementation Plan

Overview: Add Operating Kit managed doctrine and generic routing for committed, non-secret
module or extension state under `docs/repo/state/<module-or-extension-id>/`. The implementation
is an instruction-only Operating Kit release: it creates a structure-governance reference, routes
the convention from managed docs and the root managed agent block, mirrors packaged resources,
updates focused tests, and prepares the `0.1.11` release assets and release evidence.

This plan does not create state folders in consumer repositories, define Foundry internals,
create M365-specific state files, implement a plugin framework, or validate module-specific state
schemas. It only makes the durable placement rule and agent routing visible.

Essential context:

| Source | Why it matters |
| --- | --- |
| `docs/repo/plans/module-extension-state-routing/module-extension-state-routing_discovery_doc.md` | Accepted discovery, durable decisions, non-goals, risks, and implementation scope handoff. |
| `AGENTS.md` | Public-core safety, maintainer routing, managed content boundaries, and release authority boundaries for this repository. |
| `README.md` | Public repository purpose and consumer-owned boundary. |
| `docs/README.md` | Top-level documentation router that must expose this plan when updated. |
| `docs/repo/README.md` | Repository-governance router that must expose this plan when updated. |
| `docs/repo/plans/README.md` | Plan index that must link this implementation plan. |
| `docs/repo/plans/plan-register.md` | Local plan register that must track this implementation plan. |
| `docs/repo/reference/placement-contract.md` | Durable source and consumer placement rules; this change adds a consumer state path category. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact-class rules for managed instruction changes, generated surfaces, sync behavior, and release notes. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure before changing managed content, manifests, templates, tests, or release assets. |
| `docs/repo/runbooks/release-operating-kit.md` | Required procedure before publishing a public Operating Kit release. |
| `components/structure-governance/managed/README.md` | Primary managed route for structure-governance doctrine. |
| `components/structure-governance/managed/reference/documentation-structure.md` | Repository documentation structure model that must include `docs/repo/state/`. |
| `components/structure-governance/managed/reference/managed-content-boundaries.md` | Managed versus consumer-owned boundary reference that must classify committed module state. |
| `components/structure-governance/managed/runbooks/change-documentation-placement.md` | Placement runbook that must route committed module or extension state to the new path. |
| `templates/agents/AGENTS.managed-block.md` | Source template for the installed root managed block and the generic state-routing rule. |
| `components/agent-interface/managed/kit-readme.md` | Installed kit fallback inventory that must route agents to the state convention. |
| `components/structure-governance/component.yaml` | Structure-governance component manifest that must include the new managed reference. |
| `src/codeheart_operating_kit/resources/` | Packaged resource mirror installed consumers receive. |
| `tests/test_packaging_resources.py` | Existing parity test that should cover changed source and packaged resources. |
| `release-notes.md` | Consumer-facing release-note surface for the instruction-only release. |
| `pyproject.toml` | Package version surface for the `0.1.11` release. |
| `src/codeheart_operating_kit/__init__.py` | Runtime version surface for the `0.1.11` release. |
| `manifest.yaml` | Public release manifest that must point to `v0.1.11` assets and checksums. |
| `src/codeheart_operating_kit/resources/manifest.yaml` | Packaged release manifest mirror. |
| `bootstrap.md`, `install.sh`, `install.ps1` | Public bootstrap and installer surfaces that must point to validated `v0.1.11` assets. |
| `scripts/build-release-assets.py` | Release asset builder that must default to `0.1.11` after release prep. |
| `scripts/validate-release-manifest.py` | Release-manifest validation gate for release prep. |
| `tests/test_install_metadata.py`, `tests/test_release_assets.py` | Focused tests for release version and asset metadata. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Create and ship a reusable Operating Kit convention for committed module or extension state so
agents can find repo-owned routing context before asking repeated setup questions or touching
external systems.

Implementation completion is proven when:

- `components/structure-governance/managed/reference/module-extension-state.md` defines the
  `docs/repo/state/<module-or-extension-id>/` convention, ownership model, allowed content,
  prohibited content, live-preflight boundary, and examples;
- managed structure-governance docs and placement runbooks route committed, non-secret module or
  extension state to `docs/repo/state/<id>/`;
- the managed root `AGENTS.md` template contains one generic rule for module or extension state,
  without listing individual modules;
- the installed kit fallback inventory points agents to the state convention;
- no state folder, M365 state file, Foundry manifest, plugin framework, or state validator is
  created as part of this implementation;
- source managed files and packaged resource mirrors match byte-for-byte;
- package version surfaces, release manifests, installers, release notes, and release assets are
  prepared for `0.1.11`;
- release notes classify the change as an `instruction-only change` with no forced consumer
  migration and no state-folder scaffold;
- local validation covers Markdown headers, public-core hygiene, JSON schemas, packaged-resource
  parity, release-manifest validation, focused tests, and full pytest;
- any public tag, GitHub release, or consumer sync proof happens only after explicit release
  approval through `docs/repo/runbooks/release-operating-kit.md`.

## 1.2 Project And Problem Context

The immediate product signal came from Foundry M365 workspace onboarding. A consumer repository
needs committed, non-secret routing files such as `workspace-registry.yaml` or
`tenant-onboarding.yaml` so agents can avoid repeatedly asking for known tenant and workspace
context. The concrete first-use path is `docs/repo/state/m365-workspace/`, but the convention is
not M365-specific.

The broader Operating Kit problem is placement and discoverability. Root `AGENTS.md` should not
grow a line for every installed module, but agents still need one obvious generic rule: discover
installed modules through the module system that is present, then look for committed repo-owned
state under `docs/repo/state/<id>/` before asking setup questions. That local state can route the
agent, but it does not replace live external preflight or authorize tenant-changing actions.

The implementation must stay public-core safe. It may define the generic convention and examples,
but it must not include private tenant IDs, customer names, credentials, live service state, raw
logs, or module-specific secrets.

## 1.3 Current State Analysis

Current source state:

- Structure-governance docs explain documentation placement and managed boundaries, but they do
  not define a durable location for committed module or extension state.
- The source root managed block in `templates/agents/AGENTS.managed-block.md` routes agents to
  core Operating Kit docs, but it does not tell agents where committed module state belongs.
- The installed kit fallback inventory in `components/agent-interface/managed/kit-readme.md`
  routes agents to major kit surfaces, but it does not include module or extension state routing.
- `docs/repo/reference/placement-contract.md` defines durable source and consumer path
  ownership, but it does not classify `docs/repo/state/<id>/`.
- Component manifests and packaged resources mirror managed source docs under
  `src/codeheart_operating_kit/resources/`.
- `tests/test_packaging_resources.py` has parity coverage for selected managed files; it needs
  coverage for the new state reference and changed routing files.
- Current package and release surfaces are `0.1.10`; this plan prepares the next patch release,
  `0.1.11`.
- Current agent-interface and structure-governance component manifests are `0.1.7`; this plan
  bumps those modified component manifests to `0.1.8`.

Target state:

- Agents understand that committed module or extension state belongs under
  `docs/repo/state/<module-or-extension-id>/`.
- Agents discover the module or extension ID through whatever module system is installed.
- Agents check visible local committed state before asking repeated setup questions.
- Agents treat local state as routing context and still run live preflight before sensitive reads
  or external changes.
- Operating Kit owns placement doctrine only; module systems own module discovery, modules own
  concrete schemas, and consumers own committed state contents.
- Installed consumers receive the instruction-only route through normal Operating Kit sync or
  update, with no forced migration and no automatic scaffolding.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement one managed structure-governance reference, route it from managed structure docs and
agent-facing fallback surfaces, mirror packaged resources, prepare an instruction-only `0.1.11`
patch release, and keep publication behind explicit release approval. Keep the first
implementation doctrine-first and module-agnostic.

Expected source tree:

```text
Codeheart-Operating-Kit/
  docs/
    README.md                                                     # modify
    repo/
      README.md                                                   # modify
      reference/
        placement-contract.md                                     # modify
      plans/
        README.md                                                 # modify
        plan-register.md                                          # modify
        module-extension-state-routing/
          module-extension-state-routing_discovery_doc.md         # existing
          module-extension-state-routing_implementation_doc.md    # create
  components/
    agent-interface/
      component.yaml                                              # modify during release prep
      managed/
        kit-readme.md                                             # modify
    structure-governance/
      component.yaml                                              # modify
      managed/
        README.md                                                 # modify
        reference/
          documentation-structure.md                              # modify
          managed-content-boundaries.md                           # modify
          module-extension-state.md                               # create
        runbooks/
          change-documentation-placement.md                       # modify
  templates/
    agents/
      AGENTS.managed-block.md                                     # modify
  src/
    codeheart_operating_kit/
      __init__.py                                                 # modify during release prep
      resources/
        manifest.yaml                                             # modify during release prep
        components/
          agent-interface/
            component.yaml                                        # modify mirror
            managed/
              kit-readme.md                                       # modify mirror
          structure-governance/
            component.yaml                                        # modify mirror
            managed/
              README.md                                           # modify mirror
              reference/
                documentation-structure.md                        # modify mirror
                managed-content-boundaries.md                     # modify mirror
                module-extension-state.md                         # create mirror
              runbooks/
                change-documentation-placement.md                 # modify mirror
        templates/
          agents/
            AGENTS.managed-block.md                               # modify mirror
  tests/
    test_install_metadata.py                                      # release metadata validation
    test_onboard.py                                               # managed block validation
    test_packaging_resources.py                                   # source/resource parity
    test_release_assets.py                                        # release asset validation
    test_sync_check.py                                            # managed component validation
  bootstrap.md                                                    # modify during release prep
  install.sh                                                      # modify during release prep
  install.ps1                                                     # modify during release prep
  manifest.yaml                                                   # modify during release prep
  pyproject.toml                                                  # modify during release prep
  release-notes.md                                                # modify during release prep
  scripts/
    build-release-assets.py                                       # modify during release prep
    validate-release-manifest.py                                  # validation
```

No consumer module state or consumer sync target is modified during EP-01 through EP-04. The
Codeheart-HQ coordination register may be updated because it is the configured coordination home.
Consumer sync proof runs only in EP-05 after explicit release publication approval and should
validate that the updated managed instructions install cleanly without scaffolding a
`docs/repo/state/` directory.

## 2.2 Open Questions And Assumptions Requiring Clarification

OQ-1 - Target release version

- `BLOCKER: no`
- `Affects: EP-04`
- Unlocks exact version surfaces, release notes, release manifest URLs, asset names, and tag name.
- Recommended default: use `0.1.11` because the repository currently exposes `0.1.10` as the
  package and release-manifest version.

OQ-2 - Consumer impact class

- `BLOCKER: no`
- `Affects: EP-04`
- Unlocks release-note wording and whether consumer sync proof is required before publication.
- Recommended default: classify as `instruction-only change` because the plan changes managed
  docs, routes, templates, manifests, packaged resources, and tests, but does not scaffold files
  or change runtime behavior.

OQ-3 - `docs/repo/state/README.md` scaffold

- `BLOCKER: no`
- `Affects: EP-01, EP-04`
- Unlocks whether the implementation creates a generic README scaffold in consumer repositories.
- Recommended default: do not scaffold `docs/repo/state/README.md` in the first implementation.
  Let a module or extension create `docs/repo/state/<id>/` when it has real committed state to
  store.

OQ-4 - Root managed block wording

- `BLOCKER: no`
- `Affects: EP-02`
- Unlocks the exact managed `AGENTS.md` language.
- Recommended default: add one concise generic rule that tells agents to use installed module or
  extension discovery, then check `docs/repo/state/<id>/` for committed non-secret routing state.
  Do not name Foundry, M365, finance, CRM, or any specific module in the root managed block.

OQ-5 - Public release publication timing

- `BLOCKER: no`
- `Affects: EP-05`
- Unlocks public tag creation, GitHub release publication, release asset upload, and consumer
  sync proof.
- Recommended default: prepare and validate `0.1.11` release assets in EP-04, then stop before
  public tag creation and GitHub release publication until the user explicitly approves release
  publication.

## 2.3 Architectural Decisions With Reasoning

AD-1 - Place the doctrine under structure-governance

1. Problem being solved: The main decision is a durable placement and ownership rule for
   committed repo-owned state.
2. Simplest working solution: Create
   `components/structure-governance/managed/reference/module-extension-state.md` and route it
   from structure-governance docs.
3. What may change in 6-12 months: If module systems converge, a future module reference may
   explain discovery mechanics, manifests, and capability metadata.
4. Rationale: Structure-governance already owns durable path decisions and boundaries between
   managed kit content, generated surfaces, consumer-owned docs, and local-only state.
5. Alternatives considered and why not chosen: Agent-interface placement would make routing
   visible but would hide the path doctrine from placement work. A Foundry-specific placement
   would solve only the first use case.

AD-2 - Use `module or extension`, not `plugin`

1. Problem being solved: The user wants a convention broad enough for Foundry modules, shared
   services, and future extensions without committing to a plugin architecture.
2. Simplest working solution: Use `module or extension` in durable doctrine and reserve `plugin`
   for future explicit architecture work.
3. What may change in 6-12 months: A real plugin model may emerge after repeated module systems
   prove common metadata and lifecycle needs.
4. Rationale: The current need is path placement and state routing, not lifecycle design.
5. Alternatives considered and why not chosen: Introducing a plugin framework now would add
   vocabulary and obligations without proof of use.

AD-3 - Add one generic root managed-block rule

1. Problem being solved: Agents need to know where module state lives without root `AGENTS.md`
   listing every installed module.
2. Simplest working solution: Add one generic managed-block rule that points agents to installed
   module or extension discovery and then to `docs/repo/state/<id>/`.
3. What may change in 6-12 months: A future module registry could expose richer installed-module
   metadata, but root routing should remain generic.
4. Rationale: Root instructions are the highest-signal routing surface and must stay concise.
5. Alternatives considered and why not chosen: Per-module root routes would not scale and would
   make every module install compete for root instruction space.

AD-4 - Do not scaffold empty state directories

1. Problem being solved: A visible route is needed, but empty state folders in every repository
   would create noise and imply module state exists.
2. Simplest working solution: Define the convention and let modules create
   `docs/repo/state/<id>/` when committed state exists.
3. What may change in 6-12 months: A future installer may add a `docs/repo/state/README.md` when
   multiple modules use the convention and consumer confusion proves real.
4. Rationale: No scaffolding keeps the first release instruction-only and avoids inventing state
   for repositories that do not use modules.
5. Alternatives considered and why not chosen: Creating a generic README now would be useful for
   discoverability, but it would make this release alter consumer-owned tree shape.

AD-5 - Local state routes; live preflight decides

1. Problem being solved: Committed state can be stale, incomplete, or wrong.
2. Simplest working solution: Doctrine must say local state is routing context only and live
   external systems must be checked before sensitive reads, writes, deletes, permission changes,
   or tenant-changing actions.
3. What may change in 6-12 months: Repeated module patterns may justify state freshness metadata
   or module-specific validators.
4. Rationale: The convention should reduce repeated setup questions without turning committed
   files into authority over external systems.
5. Alternatives considered and why not chosen: Treating committed state as truth would be faster
   for happy paths but unsafe for external systems.

AD-6 - Module systems and modules own specifics

1. Problem being solved: Operating Kit cannot define every module registry, module ID, or state
   schema without overreaching.
2. Simplest working solution: Operating Kit owns only the generic placement convention. Module
   systems own how installed modules are discovered. Modules own concrete state files and
   schemas. Consumers own committed state values.
3. What may change in 6-12 months: Stable repeated module patterns may become a separate
   Operating Kit module-extension contract.
4. Rationale: This keeps the doctrine useful across Foundry and future systems without freezing
   their internals.
5. Alternatives considered and why not chosen: Defining a global module manifest now would
   prematurely couple independent module systems.

AD-7 - Prepare `0.1.11` release assets before publication

1. Problem being solved: Installed consumers only receive the new managed routing convention
   through a published Operating Kit release, but publishing is an external-state-changing action.
2. Simplest working solution: EP-04 prepares and validates `0.1.11` package, manifest, installer,
   release-note, and asset surfaces; EP-05 publishes only after explicit release approval.
3. What may change in 6-12 months: A future release workflow may separate release preparation and
   publication into different maintained plans.
4. Rationale: This keeps the implementation useful for consumers without hiding release authority.
5. Alternatives considered and why not chosen: Excluding release prep would make the plan a local
   docs change only. Publishing inside the same final validation epic would blur the approval
   boundary for external state.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Managed state doctrine reference and structure routes exist. | M | None |
| EP-02 | Root and fallback agent routing surfaces expose the generic convention. | S | EP-01 |
| EP-03 | Component manifests, packaged resources, and focused tests include the new route. | M | EP-01, EP-02 |
| EP-04 | `0.1.11` release prep, planning records, release notes, and validation are complete. | M | EP-01, EP-02, EP-03 |
| EP-05 | Public `v0.1.11` publication and consumer proof are complete after explicit approval. | M | EP-04 |

## EP-01 - Managed State Doctrine Reference And Structure Routes

### A) Epic ID, Title, And Outcome

EP-01 - Managed State Doctrine Reference And Structure Routes

Outcome: Structure governance contains the durable `docs/repo/state/<module-or-extension-id>/`
doctrine and routes committed module or extension state to that path.

### B) Scope

In scope:

- Create the new structure-governance managed reference.
- Define state path shape, ownership boundaries, allowed content, prohibited content, live
  preflight boundary, examples, and anti-examples.
- Update structure-governance route docs and placement runbook.
- Update the source repository placement contract with the durable consumer path category.

Out of scope:

- Creating consumer state folders.
- Creating module-specific state files.
- Defining Foundry lockfiles, module manifests, or module schemas.
- Creating generic validators for module state files.

### C) Files Touched

- `components/structure-governance/managed/reference/module-extension-state.md`
- `components/structure-governance/managed/README.md`
- `components/structure-governance/managed/reference/documentation-structure.md`
- `components/structure-governance/managed/reference/managed-content-boundaries.md`
- `components/structure-governance/managed/runbooks/change-documentation-placement.md`
- `docs/repo/reference/placement-contract.md`

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- The new reference defines `docs/repo/state/<module-or-extension-id>/` as committed,
  non-secret, repo-owned routing state.
- The new reference states that local state does not replace live preflight or authorize external
  changes.
- The new reference assigns ownership clearly: Operating Kit owns placement doctrine, module
  systems own discovery, modules own schemas, consumers own committed values.
- Structure-governance README exposes the reference.
- Documentation structure reference includes `docs/repo/state/` in the consumer repository model.
- Managed-content boundaries reference classifies `docs/repo/state/<id>/` as consumer-owned
  committed state, not managed kit content.
- Placement runbook routes committed module or extension state to `docs/repo/state/<id>/`.
- Placement contract records the durable consumer path category.

### E) Dependencies And Critical-Path Notes

This epic has no implementation dependency. It creates the doctrine other epics route and package.

Critical path: The state reference must land before the root managed block and kit fallback
inventory can point to it.

### F) Tasks Checklist

- [x] Create `components/structure-governance/managed/reference/module-extension-state.md` with timestamp, purpose, path convention, ownership model, allowed content, prohibited content, examples, and anti-examples.
- [x] Add the live-preflight rule to `components/structure-governance/managed/reference/module-extension-state.md`.
- [x] Add the module/extension state route to `components/structure-governance/managed/README.md`.
- [x] Add `docs/repo/state/` to `components/structure-governance/managed/reference/documentation-structure.md`.
- [x] Add the consumer-owned committed-state boundary to `components/structure-governance/managed/reference/managed-content-boundaries.md`.
- [x] Add placement instructions for committed module/extension state to `components/structure-governance/managed/runbooks/change-documentation-placement.md`.
- [x] Add `docs/repo/state/<id>/` to `docs/repo/reference/placement-contract.md` as a durable consumer path category.
- [x] Review changed doctrine files for `module/extension` wording and remove plugin-framework language.
- [x] Run `rg -n "docs/repo/state|module-extension-state" components/structure-governance/managed docs/repo/reference/placement-contract.md`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.

### G) Implementation Notes

Keep examples generic and public-safe. `m365-workspace`, `finance-ops`, `crm-sync`, and
`document-automation` are acceptable placeholder module IDs, but the reference must not contain
real tenant, customer, account, credential, mailbox, or workspace state.

The reference should explicitly separate committed routing context from live truth. A module can
read local state to decide what to check next, but it must still validate live systems before
sensitive reads or external mutations.

### H) Open Questions

- OQ-2 is relevant and defaults to `instruction-only change`.
- OQ-3 is relevant and defaults to no scaffold.

## EP-02 - Generic Agent Routing Surfaces

### A) Epic ID, Title, And Outcome

EP-02 - Generic Agent Routing Surfaces

Outcome: Agents see a concise generic route to committed module or extension state from the root
managed block and installed kit fallback inventory.

### B) Scope

In scope:

- Add one generic state-routing rule to the source root managed block template.
- Add the installed fallback inventory route to the kit README.
- Keep root routing module-agnostic and compact.

Out of scope:

- Adding a root route for every module.
- Naming Foundry, M365, finance, CRM, document automation, or any other concrete module in the
  root managed block.
- Changing consumer-owned repository sections.

### C) Files Touched

- `templates/agents/AGENTS.managed-block.md`
- `components/agent-interface/managed/kit-readme.md`

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- The managed root block tells agents to use installed module or extension discovery and then
  check `docs/repo/state/<id>/` for committed non-secret routing state.
- The managed root block states or clearly implies that local state does not replace live
  preflight.
- The managed root block does not list module-specific routes.
- The installed kit fallback inventory links to the structure-governance state reference.
- The route is visible enough that agents do not need to rediscover the convention from unrelated
  documents.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 because the route must point to an existing managed reference.

Critical path: This epic must complete before consumer sync proof can validate the updated agent
routing behavior.

### F) Tasks Checklist

- [x] Add a concise module/extension state rule to `templates/agents/AGENTS.managed-block.md`.
- [x] Add the installed state reference route to `components/agent-interface/managed/kit-readme.md`.
- [x] Review `templates/agents/AGENTS.managed-block.md` to confirm the new route contains no per-module names.
- [x] Review `components/agent-interface/managed/kit-readme.md` to confirm the fallback route points to the structure-governance reference.
- [x] Run `rg -n "docs/repo/state|module-extension-state" templates/agents/AGENTS.managed-block.md components/agent-interface/managed/kit-readme.md`.
- [x] Run `rg -n "Foundry|M365|finance|crm|document-automation" templates/agents/AGENTS.managed-block.md` and record no route-specific matches from the new rule.

### G) Implementation Notes

Root `AGENTS.md` should remain a router, not a module registry. The preferred shape is one short
rule under managed routes or immediate rules that says agents should check committed
module/extension state under `docs/repo/state/<id>/` after discovering installed modules through
the module system present in the repository.

### H) Open Questions

- OQ-4 is relevant and defaults to one generic root managed-block rule.

## EP-03 - Component Manifest, Packaged Resources, And Tests

### A) Epic ID, Title, And Outcome

EP-03 - Component Manifest, Packaged Resources, And Tests

Outcome: Installed consumers receive the new doctrine and routes through packaged resources, and
focused tests prove source/resource parity and managed routing expectations.

### B) Scope

In scope:

- Add the new managed reference to the structure-governance component manifest.
- Mirror changed managed files into `src/codeheart_operating_kit/resources/`.
- Mirror the changed root managed block template into packaged resources.
- Update focused tests for packaging, onboarding managed block expectations, and sync checks.

Out of scope:

- Adding validators for module state file schemas.
- Adding runtime code that interprets `docs/repo/state/<id>/`.
- Changing install behavior to scaffold state directories.

### C) Files Touched

- `components/structure-governance/component.yaml`
- `src/codeheart_operating_kit/resources/components/structure-governance/component.yaml`
- `src/codeheart_operating_kit/resources/components/structure-governance/managed/README.md`
- `src/codeheart_operating_kit/resources/components/structure-governance/managed/reference/documentation-structure.md`
- `src/codeheart_operating_kit/resources/components/structure-governance/managed/reference/managed-content-boundaries.md`
- `src/codeheart_operating_kit/resources/components/structure-governance/managed/reference/module-extension-state.md`
- `src/codeheart_operating_kit/resources/components/structure-governance/managed/runbooks/change-documentation-placement.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/kit-readme.md`
- `src/codeheart_operating_kit/resources/templates/agents/AGENTS.managed-block.md`
- `tests/test_onboard.py`
- `tests/test_packaging_resources.py`
- `tests/test_sync_check.py`

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- The structure-governance component manifest installs the new state reference to
  `.codeheart/kit/docs/structure-governance/reference/module-extension-state.md`.
- Every changed managed source file has a matching packaged resource mirror.
- Packaging tests assert parity for the new state reference and changed routing files.
- Onboarding tests account for the updated root managed block.
- Sync-check tests account for the new structure-governance managed file.
- Tests prove no state directory is scaffolded as part of normal install or sync.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 and EP-02 because packaging and tests must mirror final source content.

Critical path: The packaged resources must match source before release asset validation can pass.

### F) Tasks Checklist

- [x] Add `components/structure-governance/managed/reference/module-extension-state.md` to `components/structure-governance/component.yaml`.
- [x] Mirror changed structure-governance source files under `src/codeheart_operating_kit/resources/components/structure-governance/`.
- [x] Mirror changed agent-interface source files under `src/codeheart_operating_kit/resources/components/agent-interface/`.
- [x] Mirror changed root managed block template under `src/codeheart_operating_kit/resources/templates/agents/`.
- [x] Update `tests/test_packaging_resources.py` to assert source/resource parity for the new state reference and changed route files.
- [x] Review `tests/test_onboard.py` for managed block expectations.
- [x] Update `tests/test_sync_check.py` expected managed component content for the new structure-governance reference.
- [x] Add focused test coverage proving install and sync do not scaffold `docs/repo/state/`.

### G) Implementation Notes

Use the existing resource mirroring pattern in the repository. The implementation should not
invent a new packaging path or install target. If tests already cover managed block text through
fixtures, update the fixture expectations instead of adding broad brittle assertions.

The no-scaffold proof can be a focused assertion in an existing install or sync test. It should
verify repository tree behavior only, not module-specific state semantics.

### H) Open Questions

- OQ-3 is relevant and defaults to no scaffold.

## EP-04 - Release Prep, Planning Records, Release Notes, And Validation

### A) Epic ID, Title, And Outcome

EP-04 - Release Prep, Planning Records, Release Notes, And Validation

Outcome: The `0.1.11` instruction-only release is prepared, registered, indexed, release-noted,
asset-ready, and locally validated, with publication blocked on explicit release approval.

### B) Scope

In scope:

- Keep local and coordination-home plan registers accurate.
- Keep documentation indexes accurate.
- Prepare release notes and consumer-impact wording.
- Bump package version surfaces to `0.1.11`.
- Bump modified agent-interface and structure-governance component manifests from `0.1.7` to
  `0.1.8`.
- Build and validate `0.1.11` release assets.
- Update public release manifests, bootstrap, and installer references for `v0.1.11`.
- Run validation commands for the instruction-only release.
- Record validation results in the implementation handoff or execution log created during
  execution.

Out of scope:

- Publishing a public release without explicit approval.
- Running consumer sync proof without explicit approval.
- Marking this implementation plan complete before execution occurs.

### C) Files Touched

- `docs/README.md`
- `docs/repo/README.md`
- `docs/repo/plans/README.md`
- `docs/repo/plans/plan-register.md`
- `release-notes.md`
- `components/agent-interface/component.yaml`
- `components/structure-governance/component.yaml`
- `src/codeheart_operating_kit/resources/components/agent-interface/component.yaml`
- `src/codeheart_operating_kit/resources/components/structure-governance/component.yaml`
- `pyproject.toml`
- `src/codeheart_operating_kit/__init__.py`
- `scripts/build-release-assets.py`
- `manifest.yaml`
- `src/codeheart_operating_kit/resources/manifest.yaml`
- `bootstrap.md`
- `install.sh`
- `install.ps1`
- `tests/test_install_metadata.py`
- `tests/test_release_assets.py`
- `docs/repo/plans/module-extension-state-routing/module-extension-state-routing_implementation_doc.md`
- `Codeheart-HQ:docs/repo/plans/plan-register.md`

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- Local plan register has `OK-PR-008` for this implementation plan.
- Codeheart-HQ coordination register has `CODEHEART-OPERATING-KIT-PR-008`.
- Documentation indexes link the implementation plan.
- Release notes include a `v0.1.11` section that describes the change as instruction-only and
  states that no consumer state folders are scaffolded.
- Package version surfaces use `0.1.11`.
- Modified agent-interface and structure-governance component manifests use `0.1.8`.
- Release manifest and packaged release manifest agree with `0.1.11` asset names, URLs,
  checksums, and component checksums.
- `bootstrap.md`, `install.sh`, and `install.ps1` point to validated `v0.1.11` assets.
- Validation commands complete successfully or the execution handoff records exact failures and
  remaining risk.
- Public release publication is not performed without explicit user approval.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 through EP-03 because release notes, release assets, and validation must reflect
final changed content.

Critical path: Release prep and validation are the final local gates before release publication
approval can be requested.

### F) Tasks Checklist

- [x] Add this implementation plan to `docs/README.md`.
- [x] Add this implementation plan to `docs/repo/README.md`.
- [x] Add this implementation plan to `docs/repo/plans/README.md`.
- [x] Register `OK-PR-008` in `docs/repo/plans/plan-register.md` as the local implementation-plan entry with dependency on `OK-PR-006`.
- [x] Register `CODEHEART-OPERATING-KIT-PR-008` in `Codeheart-HQ:docs/repo/plans/plan-register.md` as the coordination-home entry.
- [x] Confirm release target `0.1.11` before editing version surfaces.
- [x] Update `release-notes.md` with a `v0.1.11` instruction-only summary, no-scaffold note, and consumer impact classification.
- [x] Update `pyproject.toml`, `src/codeheart_operating_kit/__init__.py`, and `scripts/build-release-assets.py` to release version `0.1.11`.
- [x] Update `components/agent-interface/component.yaml` and `components/structure-governance/component.yaml` component versions to `0.1.8`.
- [x] Mirror component version changes into packaged component manifests.
- [x] Build release assets with `python3 scripts/build-release-assets.py --version 0.1.11`.
- [x] Update `manifest.yaml` and `src/codeheart_operating_kit/resources/manifest.yaml` with `0.1.11` release version, checksums, URLs, and component checksums.
- [x] Update `bootstrap.md`, `install.sh`, and `install.ps1` release URLs and checksums from the generated `0.1.11` release manifest.
- [x] Run `git diff --check`.
- [x] Run `python3 scripts/validate-markdown-headers.py`.
- [x] Run `python3 scripts/validate-public-core.py`.
- [x] Run `python3 scripts/validate-json-schemas.py`.
- [x] Run `python3 scripts/validate-release-manifest.py`.
- [x] Run `pytest tests/test_packaging_resources.py tests/test_onboard.py tests/test_sync_check.py tests/test_install_metadata.py tests/test_release_assets.py`.
- [x] Run `pytest`.
- [x] Record validation evidence, release-readiness evidence, and any residual release risks in the execution handoff.

### G) Implementation Notes

During plan drafting, the plan and registers can be created before the implementation epics are
executed. During implementation execution, update this plan only for material scope corrections
and create an execution log when the execution runbook requires one.

Release-note wording should be explicit that consumers receive clearer routing instructions, not
new state files or automatic module behavior.

### H) Open Questions

- OQ-1 is relevant and defaults to `0.1.11`.
- OQ-2 is relevant and defaults to `instruction-only change`.
- OQ-5 is relevant and defaults to stopping before public release publication.

## EP-05 - Public Release Publication And Consumer Proof

### A) Epic ID, Title, And Outcome

EP-05 - Public Release Publication And Consumer Proof

Outcome: After explicit release publication approval, the validated `v0.1.11` Operating Kit
release is published publicly, release assets are verified from their published URLs, and isolated
consumer proof confirms the module/extension state routing convention installs through the normal
release path.

### B) Scope

In scope:

- Confirm explicit approval for public release publication.
- Create the public Git tag from the validated commit.
- Publish the GitHub release with release notes, manifest, installers, archives, and checksums.
- Verify published release URLs and checksums.
- Run isolated temporary consumer install proof from the published release.
- Run isolated temporary consumer sync/check proof from the published release.
- Update the implementation plan, execution evidence, local register, and HQ coordination register
  to completed state.

Out of scope:

- Publishing before explicit release approval.
- Changing release content after the validated release commit.
- Modifying consumer-owned module state during sync proof.
- Publishing a second release from this implementation plan.

### C) Files Touched

- `docs/repo/plans/module-extension-state-routing/module-extension-state-routing_implementation_doc.md`
- `docs/repo/plans/module-extension-state-routing/module-extension-state-routing_execution_log.md`
- `docs/repo/plans/plan-register.md`
- `Codeheart-HQ:docs/repo/plans/plan-register.md`

External state touched:

- Git tag `v0.1.11`
- GitHub release `v0.1.11`
- Published release assets and checksums
- Isolated temporary consumer proof workspace

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Explicit release publication approval is recorded before creating tag `v0.1.11`.
- Git tag `v0.1.11` points to the validated commit from EP-04.
- GitHub release `v0.1.11` includes release notes, manifest, installers, archives, and checksums.
- Published release asset checksums match `manifest.yaml`.
- Isolated consumer install proof installs the module/extension state reference under
  `.codeheart/kit/docs/structure-governance/reference/module-extension-state.md`.
- Isolated consumer proof shows the root managed block contains the generic module/extension
  state route and does not scaffold `docs/repo/state/`.
- Local plan register and HQ coordination register mark the implementation complete.

### E) Dependencies And Critical-Path Notes

Depends on EP-04 because publication must use validated release assets from the validated commit.

Critical path: Do not create the public tag until approval is explicit and EP-04 validation
evidence is complete.

### F) Tasks Checklist

- [ ] Confirm explicit release publication approval for `v0.1.11`.
- [ ] Record release publication approval in `docs/repo/plans/module-extension-state-routing/module-extension-state-routing_execution_log.md`.
- [ ] Create Git tag `v0.1.11` from the EP-04 validated commit.
- [ ] Publish GitHub release `v0.1.11` with release notes, manifest, installers, archives, and checksums.
- [ ] Verify published release asset URLs and checksums against `manifest.yaml`.
- [ ] Run isolated temporary consumer install proof from published `v0.1.11` assets.
- [ ] Run isolated temporary consumer sync/check proof from published `v0.1.11` assets.
- [ ] Verify installed `.codeheart/kit/docs/structure-governance/reference/module-extension-state.md` exists in the proof workspace.
- [ ] Verify installed root `AGENTS.md` contains the generic module/extension state route and no `docs/repo/state/` scaffold exists in the proof workspace.
- [ ] Mark `docs/repo/plans/module-extension-state-routing/module-extension-state-routing_implementation_doc.md` complete after release proof passes.
- [ ] Mark `OK-PR-008` complete in `docs/repo/plans/plan-register.md`.
- [ ] Mark `CODEHEART-OPERATING-KIT-PR-008` complete in `Codeheart-HQ:docs/repo/plans/plan-register.md`.
- [ ] Record release URLs, checksums, validation evidence, consumer proof evidence, and residual risk in the execution log.

### G) Implementation Notes

This epic changes external state. It must not run from a general implementation approval alone
unless the user explicitly includes public release publication approval for `v0.1.11`.

Keep the consumer proof isolated. It should prove installed managed content and absence of
`docs/repo/state/` scaffolding, not modify a real consumer repository.

### H) Open Questions

- OQ-5 is relevant and defaults to stopping before public release publication.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

Future work intentionally deferred from this implementation:

- A Foundry-specific module registry or lockfile contract.
- M365 workspace state files such as `workspace-registry.yaml` and `tenant-onboarding.yaml`.
- A broader module or extension architecture.
- A plugin framework.
- Validators for module state schemas.
- A generic `docs/repo/state/README.md` scaffold.
- A migration of existing consumer or module state.
- A tooling/environment readiness register for module onboarding.
- A state freshness convention or schema for external preflight evidence.

## 4.2 Future Considerations

Likely follow-up candidates:

- Foundry M365 can use this convention in its own implementation plan after Operating Kit
  structure-governance routing is released.
- Runbook authoring standards can reference the convention later as an example of visible local
  context agents should check before asking repeated setup questions.
- Operating Kit can discover a broader module-extension contract after at least two module
  systems need shared lifecycle, registry, or capability metadata.

# Revision Notes

- 2026-06-23: Drafted implementation plan from
  `module-extension-state-routing_discovery_doc.md` implementation scope handoff.
- 2026-06-23: Updated plan to target `0.1.11` release prep, add approval-gated publication proof,
  and address implementation-plan review findings.
- 2026-06-23: Activated implementation plan for execution.
