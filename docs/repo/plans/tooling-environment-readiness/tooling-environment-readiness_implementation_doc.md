Last updated: 2026-06-24T14:16:20Z (UTC)
Created: 2026-06-24
Status: active
Execution log: docs/repo/plans/tooling-environment-readiness/tooling-environment-readiness_execution_log.md

# Document Header

## Tooling Environment Readiness Implementation Plan

Overview: Add an Operating Kit managed tooling-readiness route so agents have one clear way to
handle missing local tools during module onboarding or module operations. The implementation is an
instruction-only Operating Kit release: it creates one central readiness runbook with a small
on-demand baseline catalog, routes it from installed managed surfaces, updates runbook authoring
and planning checks, mirrors packaged resources, prepares `0.1.12` release assets, and keeps
publication behind explicit approval.

This plan does not build an environment manager, install tools by default, require
machine-readable module tool declarations, record durable machine-readiness state, retrofit
Foundry M365 runbooks, or move module-owned domain setup into Operating Kit. Modules remain
responsible for module-specific install commands, authentication, permissions, and live external
preflight.

Essential context:

| Source | Why it matters |
| --- | --- |
| `docs/repo/plans/tooling-environment-readiness/tooling-environment-readiness_discovery_doc.md` | Accepted discovery, durable decisions, trigger model, anti-sprawl boundary, and implementation handoff. |
| `docs/repo/plans/runbook-authoring-standards/runbook-authoring-standards_implementation_doc.md` | Recently implemented runbook-quality pattern that this plan extends for environment blockers. |
| `docs/repo/plans/module-extension-state-routing/module-extension-state-routing_implementation_doc.md` | Recently implemented managed route pattern for root `AGENTS.md`, kit fallback inventory, packaged mirrors, and release prep. |
| `AGENTS.md` | Public-core safety, managed content boundaries, and maintainer routing for this repository. |
| `README.md` | Public repository purpose and consumer-owned boundary. |
| `docs/README.md` | Top-level documentation router that may need discoverability updates if touched. |
| `docs/repo/README.md` | Repository-governance router that must expose this implementation plan. |
| `docs/repo/plans/README.md` | Plan index that must link this implementation plan. |
| `docs/repo/plans/plan-register.md` | Local plan register that must track this implementation plan. |
| `docs/repo/reference/placement-contract.md` | Placement contract that keeps managed kit doctrine separate from consumer-owned and module-owned content. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact-class rules for managed instruction changes, generated surfaces, sync behavior, and release notes. |
| `docs/repo/runbooks/change-operating-kit.md` | Required maintainer procedure before changing managed content, manifests, templates, tests, or release assets. |
| `docs/repo/runbooks/release-operating-kit.md` | Required procedure before publishing a public Operating Kit release. |
| `components/agent-interface/managed/README.md` | Owning route for agent-facing and hybrid managed guidance. |
| `components/agent-interface/managed/kit-readme.md` | Installed kit fallback inventory that must route agents when the root managed block is insufficient. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Existing runbook standard that must route local environment blockers without duplicating install doctrine. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Planning workflow hook for future plans that create or change runbooks with tool readiness concerns. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Execution workflow hook for verifying changed runbooks apply the readiness route. |
| `components/planning-workflows/managed/runbooks/review-planning-document.md` | Review workflow hook for flagging missing route, vague install guidance, or Operating Kit runbook sprawl. |
| `components/structure-governance/managed/README.md` | Structure-governance route that should clarify placement versus runbook-shape ownership. |
| `components/structure-governance/managed/runbooks/change-documentation-placement.md` | Placement runbook that should route durable runbook shape to the authoring standard and central readiness route. |
| `templates/agents/AGENTS.managed-block.md` | Source template for the installed root managed block and the concise missing-tool route. |
| `components/*/component.yaml` | Component manifests that must include changed managed files and version bumps where applicable. |
| `src/codeheart_operating_kit/resources/` | Packaged resource mirror installed consumers receive. |
| `tests/test_packaging_resources.py` | Existing parity test that should cover new and changed managed resources. |
| `release-notes.md` | Consumer-facing release-note surface for the instruction-only release. |
| `pyproject.toml` | Package version surface for the `0.1.12` release preparation. |
| `src/codeheart_operating_kit/__init__.py` | Runtime version surface for the `0.1.12` release preparation. |
| `manifest.yaml` | Public release manifest that must point to `v0.1.12` assets and checksums during release prep. |
| `src/codeheart_operating_kit/resources/manifest.yaml` | Packaged release manifest with current release URLs, component metadata, and established zero-placeholder downloadable asset checksums. |
| `bootstrap.md`, `install.sh`, `install.ps1` | Public bootstrap and installer surfaces that must point to validated `v0.1.12` assets during release prep. |
| `scripts/build-release-assets.py` | Release asset builder that must default to `0.1.12` during release prep. |
| `scripts/validate-release-manifest.py` | Release-manifest validation gate for release prep. |
| `tests/test_install_metadata.py`, `tests/test_release_assets.py`, `tests/test_sync_check.py`, `tests/test_onboard.py` | Focused tests for version metadata, release assets, managed routes, sync behavior, and generated root block content. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Create and ship a reusable Operating Kit route for local tooling and environment readiness so
agents can recover cleanly when a module onboarding or operation is blocked by a missing package
manager, runtime, CLI, PowerShell module, install path, or local tool.

Implementation completion is proven when:

- `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` exists as one central
  hybrid readiness runbook with a compact intention block, trigger model, blocker classification,
  missing-tool behavior contract, approval gates, stop conditions, validation, and return-to-module
  handoff;
- the readiness runbook contains a small on-demand baseline catalog for generic lanes:
  package-manager/bootstrap, PowerShell runtime, PowerShell module, Node.js, Python, browser
  automation, and document/PDF tooling;
- the catalog is explicitly not a default install bundle and not a module-specific install index;
- the readiness runbook gives concrete nontechnical choices for common blockers, such as
  "Install Homebrew", "I will install it another way", and "Stop here";
- managed root `AGENTS.md` and the installed kit fallback inventory expose one concise tooling
  readiness route without listing every module or every tool;
- the runbook authoring standard tells future human-facing, agent-facing, and hybrid runbooks to
  route local environment blockers through the managed readiness route instead of duplicating
  generic install doctrine;
- planning, execution, and planning-review workflows check this route only when plans create or
  materially change runbooks that can hit local tooling blockers;
- no local readiness state, validator, schema, CLI environment manager, consumer scaffold, or
  module-specific M365/AWS/CRM install runbook is created as part of this implementation;
- source managed files and packaged resource mirrors match byte-for-byte;
- root `manifest.yaml` records publishable release-asset hashes while packaged
  `src/codeheart_operating_kit/resources/manifest.yaml` keeps zero-placeholder downloadable asset
  hashes to avoid a self-referential archive checksum;
- package version surfaces, release manifests, installers, release notes, and release assets are
  prepared for `0.1.12`;
- release notes classify the change as an `instruction-only change` with no forced consumer
  migration and no default tool installation;
- local validation covers Markdown headers, public-core hygiene, packaged-resource parity, route
  visibility, release-manifest validation, focused tests, and full pytest;
- any public tag, GitHub release, or consumer sync proof happens only after explicit release
  approval through `docs/repo/runbooks/release-operating-kit.md`.

## 1.2 Project And Problem Context

The immediate product signal came from the Foundry Microsoft 365 module. That module may need
PowerShell, Microsoft Graph PowerShell, PnP PowerShell, and Exchange Online PowerShell. Those
module-specific tools and Microsoft service permissions are legitimate Foundry module
responsibilities. The reusable Operating Kit gap is lower-level: when a consumer machine is
missing common bootstrap tooling, agents need a standard way to explain the blocker, ask for
approval, install or repair through the chosen route, recheck readiness, and then return to the
module runbook.

The current risk is not only technical failure. Without a managed readiness route, an agent may
tell a nontechnical user that it lacks capability, ask the user to solve package-manager setup
alone, or bury the user in local implementation details. The desired user experience is concrete:
explain the outcome, present one useful decision, ask approval before local changes, and stop
cleanly when the user does not want the local machine changed.

The implementation must also protect Operating Kit architecture. Operating Kit should not grow a
separate managed runbook for every module or tool ecosystem. It should own one shared route, a
small reusable lane catalog, and the authoring standard that tells module runbooks when to use that
route. Modules remain allowed to own module-specific install runbooks and commands.

## 1.3 Current State Analysis

Current source state:

- `components/agent-interface/managed/reference/runbook-authoring-standard.md` already defines
  audience classes, compact intention blocks, human-facing flow, agent-facing execution paths,
  hybrid separation, approval gates, stop conditions, and review checks.
- The runbook authoring standard mentions tool readiness checks for agent-facing runbooks, but it
  does not yet name a shared readiness route or anti-sprawl rule for local environment blockers.
- `templates/agents/AGENTS.managed-block.md` exposes planning, feedback, structure governance,
  and module state routes, but it does not yet expose tooling readiness.
- `components/agent-interface/managed/kit-readme.md` is the installed fallback inventory, but it
  does not yet include a missing-tool route.
- Planning workflow runbooks do not yet require plans or reviews to distinguish module-owned
  service blockers from Operating Kit local environment blockers.
- Component manifests and packaged resources mirror managed source docs under
  `src/codeheart_operating_kit/resources/`.
- The package and public release surfaces currently show `0.1.11`; this plan prepares the next
  patch release, `0.1.12`.

Target state:

- Agents can find one managed readiness runbook when module onboarding or module operations hit a
  local tooling blocker.
- Agents classify the blocker into a generic lane, use concrete human-facing choices, ask approval
  before local changes, recheck readiness, and return to the module runbook.
- Agents do not confuse local readiness with live external preflight. Tenant consent, admin roles,
  mailbox access, SharePoint permissions, licenses, app readiness, and service availability remain
  module-owned blockers.
- Future runbooks remain DRY: they may declare required tools and module-specific commands, but
  they should route generic package-manager/runtime/tool readiness through the central managed
  runbook.
- Installed consumers receive the instruction-only route through normal Operating Kit update or
  sync, without any automatic installation or consumer-owned file scaffold.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement one central managed runbook under agent-interface, route it from installed surfaces, add
runbook-authoring and planning workflow hooks, mirror packaged resources, and prepare an
instruction-only `0.1.12` patch release. Keep the first implementation doctrine-first and
human-usable; defer schemas, validators, durable readiness state, and module-specific runbook
retrofits.

Expected source tree:

```text
Codeheart-Operating-Kit/
  docs/
    repo/
      README.md                                                   # modify
      plans/
        README.md                                                 # modify
        plan-register.md                                          # modify
        tooling-environment-readiness/
          tooling-environment-readiness_discovery_doc.md          # existing
          tooling-environment-readiness_implementation_doc.md     # create
  components/
    agent-interface/
      component.yaml                                              # modify
      managed/
        README.md                                                 # modify
        kit-readme.md                                             # modify
        reference/
          runbook-authoring-standard.md                           # modify
        runbooks/
          handle-tooling-readiness.md                             # create
    planning-workflows/
      component.yaml                                              # modify
      managed/
        runbooks/
          draft-implementation-plan.md                            # modify
          execute-implementation-plan.md                          # modify
          review-planning-document.md                             # modify
    structure-governance/
      component.yaml                                              # modify
      managed/
        README.md                                                 # modify
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
              README.md                                           # modify mirror
              kit-readme.md                                       # modify mirror
              reference/
                runbook-authoring-standard.md                     # modify mirror
              runbooks/
                handle-tooling-readiness.md                       # create mirror
          planning-workflows/
            component.yaml                                        # modify mirror
            managed/
              runbooks/
                draft-implementation-plan.md                      # modify mirror
                execute-implementation-plan.md                    # modify mirror
                review-planning-document.md                       # modify mirror
          structure-governance/
            component.yaml                                        # modify mirror
            managed/
              README.md                                           # modify mirror
              runbooks/
                change-documentation-placement.md                 # modify mirror
        templates/
          agents/
            AGENTS.managed-block.md                               # modify mirror
  tests/
    test_install_metadata.py                                      # release metadata validation
    test_onboard.py                                               # managed block route validation
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

No Foundry M365 runbook is modified by this plan. A later Foundry-specific implementation can
update M365 onboarding Step 7 to invoke the managed readiness route once this Operating Kit route
exists.

## 2.2 Resolved Planning Questions

### OQ-1 - Where Should Readiness Observations Live?

Resolution for this implementation: defer durable readiness observations. V1 does not create a
state file, schema, report location, or lockfile extension. The readiness runbook may tell agents
to record a short capability blocker in the current work summary or module run record when the
active task requires evidence, but it must not commit machine-specific readiness under
`docs/repo/state/<id>/`.

### OQ-2 - Which Tool Families Should V1 Cover?

Resolution for this implementation: include a small on-demand baseline catalog covering:

- package-manager/bootstrap;
- PowerShell runtime;
- PowerShell module;
- Node.js;
- Python;
- browser automation;
- document/PDF tooling.

Cloud CLIs and domain-specific service modules stay out of the first Operating Kit catalog unless
future repeated evidence proves they should graduate into generic lanes.

### OQ-3 - Should Tool Declarations Be Machine-Readable In V1?

Resolution for this implementation: no. V1 accepts explicit declaration by module manifest,
reference, or runbook. A module declaration is sufficient when it names the required or optional
tool, the capability unlocked, the readiness check or expected command, the missing-tool route,
and the module-owned live preflight boundary.

### OQ-4 - How Much Install Guidance Should Operating Kit Own?

Resolution for this implementation: Operating Kit owns the generic local readiness route,
concrete user choice patterns, and common baseline lanes. It does not own module-specific install
commands or service-specific setup.

The package-manager/bootstrap lane should be concrete enough for nontechnical users. For example,
on macOS it should be able to offer "Install Homebrew", "I will install it another way", and
"Stop here". The implementation should verify official install sources before adding any
vendor-specific link or command.

### OQ-5 - Should V1 Touch Existing Human-Facing Onboarding Runbooks?

Resolution for this implementation: no broad retrofit. Existing Operating Kit or module runbooks
are changed only where they route to the new shared standard. Foundry M365 onboarding changes are
deferred to a module-owned follow-up plan.

### OQ-6 - Is Release Publication Part Of This Plan?

Resolution for this implementation: release preparation is part of the implementation plan.
Public publication, tag creation, GitHub release creation, and consumer proof are approval-gated
steps and must not run until explicitly approved.

## 2.3 Consumer Impact Classification

Expected impact class: `instruction-only change`.

Reasons:

- managed docs, templates, component manifests, packaged resources, release metadata, and tests
  change;
- no consumer-owned files are scaffolded;
- no local tools are installed;
- no module-owned runbooks are changed;
- no runtime environment manager is introduced;
- no migration is required for existing consumers.

Consumer-visible effect:

- after update or sync, agents see a route for local tooling blockers from root managed guidance
  and kit fallback inventory;
- future runbooks and planning workflows have clearer readiness expectations;
- installed modules may later choose to route environment blockers through this shared runbook.

## 2.4 Boundaries And Non-Goals

Do not implement:

- automatic Homebrew, PowerShell, Node.js, Python, browser, or document tooling installation;
- a default "install all baseline tools" flow;
- a machine-readable `tools:` schema or validator;
- readiness state in `.codeheart/kit.lock.yaml`, `.codeheart/user/`, `docs/repo/state/<id>/`, or
  generated reports;
- a broad environment manager CLI;
- module-specific install runbooks for M365, AWS, CRM, finance, document automation, or other
  module domains;
- Foundry M365 runbook retrofits;
- restricted-device or enterprise endpoint-management flows beyond a simple stop path;
- live service preflight for Microsoft 365, SharePoint, Exchange, Graph, AWS, or other external
  systems.

# Section 3 - Execution Plan

## EP-01 - Add Central Tooling Readiness Runbook

Purpose: create one central managed readiness route that agents can follow whenever module
onboarding or operation is blocked by missing local tooling.

Scope:

- create `components/agent-interface/managed/runbooks/handle-tooling-readiness.md`;
- classify the runbook as `hybrid`;
- include a compact intention block;
- define local environment blockers versus module-owned service blockers;
- define the trigger model from the discovery;
- define the missing-tool behavior contract;
- define the small on-demand baseline catalog;
- define concrete human-facing choice patterns;
- define approval gates before local installs, repairs, path changes, permission prompts, or
  sensitive checks;
- define validation and return-to-module behavior;
- define unresolved blocker reporting;
- explicitly state that V1 records no durable machine-readiness state.

Runbook content requirements:

- Trigger when a module onboarding or operation reports a missing package manager, runtime, CLI,
  PowerShell module, install path, or local tool.
- Do not trigger for module-owned service blockers such as tenant consent, admin roles, mailbox
  access, SharePoint permissions, licenses, app readiness, or service availability.
- Ask at most one user-owned decision per user turn unless the runbook has a clear reason to ask
  two.
- Use concrete choices for common blockers.
- Keep internal file paths, module manifests, command traces, and local diagnostics out of the
  first user-facing explanation.
- Use read-only checks before proposing installation when locally appropriate.
- Ask explicit user approval before local changes.
- Use official vendor sources or module-owned runbooks for actual install commands.
- Recheck readiness after approved install or repair.
- Return control to the calling module runbook when the local blocker is resolved.
- Stop and report a capability blocker when the user declines, install is not possible, or the
  check remains blocked.

Baseline catalog requirements:

| Lane | V1 intent | Ownership boundary |
| --- | --- | --- |
| `package-manager-bootstrap` | Help the agent handle missing package manager or bootstrap route. | Operating Kit owns generic choices; platform-specific official install details must be verified during implementation. |
| `powershell-runtime` | Help the agent handle missing PowerShell runtime. | Operating Kit owns generic readiness route; modules own why PowerShell is needed. |
| `powershell-module` | Help the agent handle missing PowerShell module family. | Operating Kit owns generic module-install pattern; modules own concrete module names, versions, and commands. |
| `node-runtime` | Help the agent handle missing Node.js or package manager runtime. | Operating Kit owns generic lane; repositories or modules own exact version requirements. |
| `python-runtime` | Help the agent handle missing Python runtime or package manager. | Operating Kit owns generic lane; repositories or modules own exact version and virtual environment requirements. |
| `browser-automation` | Help the agent handle missing local browser automation prerequisite. | Operating Kit owns generic lane; specific browser tools or plugins own concrete setup. |
| `document-pdf-tooling` | Help the agent handle missing document conversion or PDF tooling. | Operating Kit owns generic lane; document/PDF skills or modules own exact tool requirements. |

Review checks:

- The runbook does not include private machine details, tenant details, local paths, or raw logs.
- The runbook does not include a default install bundle.
- The runbook does not contain domain-specific M365/AWS/CRM install procedures.
- The runbook gives a nontechnical user an understandable next decision when a tool is missing.

## EP-02 - Expose The Route And Update Authoring Standards

Purpose: make the readiness route visible from installed Operating Kit guidance and make future
runbooks use the route without duplicating generic install doctrine.

Scope:

- update `components/agent-interface/managed/README.md` with the new readiness runbook route;
- update `components/agent-interface/managed/kit-readme.md` with the installed fallback route;
- update `templates/agents/AGENTS.managed-block.md` with one concise immediate rule and one
  managed route entry for tooling readiness;
- update `components/structure-governance/managed/README.md` and
  `components/structure-governance/managed/runbooks/change-documentation-placement.md` with a
  lightweight cross-reference: structure governance owns placement, while agent-interface owns
  durable runbook shape and generic tooling-readiness behavior;
- update `components/agent-interface/managed/reference/runbook-authoring-standard.md` with:
  - local environment blocker routing;
  - blocker-specific human choices;
  - the DRY rule for generic tool guidance;
  - module-owned service blocker boundaries;
  - review checks for changed runbooks that can hit missing local tools;
- avoid adding route entries for individual modules or tools.

Root managed block wording target:

```text
When a module onboarding or operation is blocked by missing local tooling, follow the managed
tooling readiness route before installing, repairing, improvising setup, or declaring the
capability unavailable.
```

Installed route target:

```text
Tooling readiness:
  `.codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md`
```

Runbook authoring standard target:

- Human-facing and hybrid runbooks should use the readiness route for local environment blockers
  and present concrete choices, not abstract package-manager terminology.
- Agent-facing runbooks should name required tools and route missing generic prerequisites through
  the managed readiness runbook.
- Module-specific install commands may stay in module runbooks or references.
- Operating Kit managed runbooks should not duplicate package-manager/runtime setup for every
  module or tool.

Review checks:

- The root block remains concise and generic.
- The route is visible without forcing root `AGENTS.md` to list individual modules.
- The authoring standard improves DRY architecture without becoming a broad environment-manager
  specification.

## EP-03 - Add Planning, Execution, And Review Hooks

Purpose: make future implementation planning, execution, and review detect missing or duplicated
tooling-readiness architecture when runbooks are created or materially changed.

Scope:

- update `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`;
- update `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`;
- update `components/planning-workflows/managed/runbooks/review-planning-document.md`;
- update `components/planning-workflows/component.yaml` and packaged mirrors if files change.

Planning hook requirements:

- When a plan creates or materially changes human-facing, agent-facing, or hybrid runbooks that
  can hit missing local tooling, it must state:
  - which tool blockers can occur;
  - whether they are local environment blockers or module-owned service blockers;
  - which route handles local blockers;
  - whether any module-specific install command remains module-owned;
  - how the plan avoids duplicating generic Operating Kit tooling doctrine.

Execution hook requirements:

- Before marking a runbook-related epic complete, check that changed runbooks:
  - declare their audience class and compact intention block where required;
  - route local environment blockers to the managed readiness runbook;
  - ask approval before installs or repairs;
  - do not leak internal mechanics into first-turn user-facing copy;
  - do not copy broad package-manager/runtime install guidance into module-specific docs unless
    the plan explicitly justifies it.

Review hook requirements:

- Flag plans that:
  - leave agents without an execution path when a required tool is missing;
  - use vague "install tools" wording for a nontechnical user when the blocker can be classified;
  - duplicate generic package-manager or runtime setup in multiple managed runbooks;
  - confuse local tool readiness with live external permissions or service preflight.

Review checks:

- Hooks apply only when plans create or materially change runbooks with tooling blockers.
- Hooks do not force every implementation plan to discuss tooling.
- Hooks do not require module runbooks to adopt machine-readable declarations.

## EP-04 - Mirror Resources, Tests, Docs, And Release Prep

Purpose: keep source and packaged Operating Kit resources aligned, make the new route testable,
and prepare a complete `0.1.12` release candidate.

Scope:

- mirror all changed managed files under `src/codeheart_operating_kit/resources/`;
- update `components/agent-interface/component.yaml`;
- update `components/planning-workflows/component.yaml` if planning workflow files change;
- update `components/structure-governance/component.yaml` if structure-governance files change;
- update packaged component manifests;
- update `tests/test_packaging_resources.py` to cover new and changed managed files;
- update `tests/test_onboard.py` or equivalent generated root block tests so the tooling
  readiness route is protected;
- update `tests/test_sync_check.py` if component manifest or managed route validation requires it;
- update `docs/repo/README.md`, `docs/repo/plans/README.md`, and
  `docs/repo/plans/plan-register.md`;
- prepare `0.1.12` release metadata in `pyproject.toml`,
  `src/codeheart_operating_kit/__init__.py`, `manifest.yaml`,
  `src/codeheart_operating_kit/resources/manifest.yaml`, `bootstrap.md`, `install.sh`,
  `install.ps1`, `scripts/build-release-assets.py`, and `release-notes.md`;
- build release assets locally only after the implementation is otherwise validated;
- keep publication approval-gated.

Required validation:

```bash
python3 scripts/validate-markdown-headers.py
python3 scripts/validate-public-core.py
python3 scripts/validate-json-schemas.py
python3 scripts/validate-release-manifest.py
python3 -m pytest tests/test_packaging_resources.py
python3 -m pytest tests/test_onboard.py tests/test_sync_check.py
python3 -m pytest tests/test_install_metadata.py tests/test_release_assets.py
python3 -m pytest
git diff --check
```

Release-prep review checks:

- Release notes call this an instruction-only change.
- Release notes say the update does not install local tools by default.
- Release notes say modules still own module-specific commands and live service preflight.
- Release metadata consistently targets `0.1.12`.
- Root asset names, manifest URLs, and checksum entries consistently target `v0.1.12` and real
  local hashes.
- Packaged manifest asset names and URLs consistently target `v0.1.12`, and packaged asset hashes
  use the established zero-placeholder pattern.

## EP-05 - Approval-Gated Publication And Consumer Proof

Purpose: publish the release and prove installed consumers can receive the managed route, but only
after explicit release approval.

Approval gate:

- Do not push release commits, create tag `v0.1.12`, create a GitHub release, upload release
  assets, or update configured consumers until the user explicitly approves publication.

Publication scope after approval:

- run the release runbook in `docs/repo/runbooks/release-operating-kit.md`;
- push the release commit to the public repository;
- create and push tag `v0.1.12`;
- publish the GitHub release with macOS and Windows assets, checksums, bootstrap, installers,
  release notes, and manifest;
- verify published asset URLs and checksums;
- run isolated consumer proof that install/sync produces the managed root route and packaged
  readiness runbook without creating consumer-owned tooling state;
- update configured consumer repositories only when explicitly requested or when the release
  runbook requires a separately approved consumer proof;
- record publication and proof evidence in the implementation execution log and plan register.

Completion criteria:

- `v0.1.12` assets are published and verified;
- isolated consumer proof passes;
- no tool installation is performed as part of consumer proof;
- no module-specific runbook is changed during Operating Kit publication;
- plan status and register status are updated according to actual execution outcome.

# Section 4 - Future Planning

## 4.1 Deferred Follow-Ups

Foundry M365 onboarding integration:

- update the M365 onboarding runbook so its local tooling step invokes the managed Operating Kit
  readiness route for missing Homebrew, PowerShell, and PowerShell modules;
- keep Microsoft tenant consent, admin roles, mailbox access, SharePoint permissions, licenses,
  app readiness, and live validation inside the Foundry M365 module;
- preserve the module-specific execution paths for Graph PowerShell, PnP PowerShell, and Exchange
  Online PowerShell.

Machine-readable tool declarations:

- revisit a `tools:` manifest shape only after multiple modules need automated extraction of tool
  requirements;
- decide whether the declaration belongs in module manifests, module references, runbook front
  matter, or a shared schema.

Readiness observation placement:

- revisit whether local readiness observations should remain ephemeral or be written to
  `.codeheart/user/`, generated reports, or another ignored local path;
- do not store machine readiness under committed module state.

Restricted or policy-managed devices:

- add more explicit guidance for devices where the user lacks install rights, Homebrew is not
  allowed, endpoint management blocks installs, or company policy requires IT-admin installation;
- V1 only needs the safe stop path.

Catalog expansion:

- consider cloud CLI lanes, Java/.NET runtimes, database CLIs, or other tool families only after
  repeated modules need them;
- avoid promoting one-off module tools into Operating Kit managed runbooks.

Operating Kit first-run onboarding:

- do not retrofit first-run onboarding in this plan;
- revisit only when first-run onboarding itself is materially changed.

## 4.2 Risks And Mitigations During Implementation

Risk: the central runbook becomes too abstract for nontechnical users.

Mitigation: require concrete user choice examples and verify the first-turn copy does not lead
with package-manager/bootstrap terminology.

Risk: the central runbook becomes a long install manual for every tool.

Mitigation: keep actual commands and version constraints module-owned or official-source routed.
Only generic lanes and behavior belong in Operating Kit V1.

Risk: planning hooks make every future plan discuss tooling even when irrelevant.

Mitigation: apply hooks only when plans create or materially change runbooks that can hit local
tooling blockers.

Risk: agents mistake local readiness for external authorization.

Mitigation: state the boundary repeatedly in the readiness runbook, authoring standard, and review
hooks: local tool checks do not replace module-owned live preflight.

Risk: release prep updates version surfaces before the implementation passes validation.

Mitigation: perform release prep after managed docs, mirrors, and tests pass focused validation;
publication still requires explicit approval.

## 4.3 Handoff To Future Implementer

Start by reading:

1. `docs/repo/plans/tooling-environment-readiness/tooling-environment-readiness_discovery_doc.md`
2. this implementation plan;
3. `docs/repo/runbooks/change-operating-kit.md`;
4. `components/agent-interface/managed/reference/runbook-authoring-standard.md`;
5. `templates/agents/AGENTS.managed-block.md`;
6. `components/agent-interface/managed/kit-readme.md`;
7. the current packaging-resource tests.

Work in this order:

1. Draft the central readiness runbook first.
2. Add route visibility and authoring-standard updates.
3. Add planning workflow hooks.
4. Mirror resources and update manifests/tests.
5. Run focused validation.
6. Prepare release metadata and assets.
7. Stop for explicit publication approval.

Do not use Foundry M365 documents as a source of truth for the Operating Kit implementation. They
are useful as evidence for why the route matters, but the Operating Kit implementation must remain
generic and public-core safe.

# Revision Notes

- 2026-06-24: Created draft implementation plan from tooling environment readiness discovery.
- 2026-06-24: Activated plan for implementation and created sibling execution log.
- 2026-06-24: Added a low-risk structure-governance cross-reference so placement guidance points
  runbook shape and generic tooling-readiness behavior back to agent-interface standards.
- 2026-06-24: Clarified release-manifest strategy: root manifest records publishable hashes, while
  packaged manifest keeps zero-placeholder downloadable asset hashes because archives cannot embed
  their own final checksum.
