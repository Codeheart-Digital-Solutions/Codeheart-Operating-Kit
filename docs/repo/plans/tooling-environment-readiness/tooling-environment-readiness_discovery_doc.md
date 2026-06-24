Last updated: 2026-06-24T13:38:01Z (UTC)
Created: 2026-06-24
Status: draft

# Tooling Environment Readiness Discovery

## Overview

This discovery investigates whether Codeheart Operating Kit should define shared tooling and
environment-readiness doctrine for consumer repositories and installed modules.

The immediate trigger is a Foundry Microsoft 365 module that needs local tooling such as
PowerShell and Microsoft 365 PowerShell modules. Those module-specific needs are legitimate module
responsibilities. The reusable gap is one layer lower: some consumers may not have common bootstrap
tooling such as a package manager, shell runtime, or repair path for a missing tool. When a module
runbook reaches a missing-tool condition, the agent should not invent ad hoc install advice or tell
the user it has no capability. It should have a visible Operating Kit route for diagnosing missing
local tooling, explaining why it matters, and guiding the user through an approval-gated install or
repair path.

Recommended direction for discovery: Operating Kit should own an on-demand baseline tooling
catalog, generic environment-readiness routing doctrine, and reusable human-facing install or repair
patterns for common tool families. The baseline catalog is not a default install bundle. It is the
set of tooling lanes every Operating Kit consumer agent should know how to recognize, explain,
check, install after approval, or route around when a module asks for that tool family. Modules
should own their concrete required and optional tool declarations, domain-specific PowerShell
modules or CLIs, service permissions, and live external preflight. Tool readiness should be treated
as local, timestamped, non-secret evidence, not as committed external truth or permission to change
external systems.

The anti-sprawl goal applies to Operating Kit managed surfaces. Operating Kit should not grow a
separate runbook for every module or tool. It should provide one central readiness route and a small
baseline catalog. Modules remain allowed to own module-specific install runbooks and commands.

The route should trigger whenever a module onboarding or operation hits an environment blocker,
not only during initial Operating Kit onboarding. A user should receive concrete choices such as
"Install Homebrew", "I will install it another way", or "Stop here", not abstract package-manager
or bootstrap terminology.

Review of the current managed routes and the Foundry M365 onboarding runbook adds one practical
implementation signal: the route must be visible from installed Operating Kit guidance before a
module needs it. The current M365 onboarding script already has a local tooling step, but it treats
missing Microsoft tools as one broad blocker. The shared Operating Kit route should split local
environment blockers into concrete layers such as package manager, runtime, and module/tool
installation, then return control to the module runbook after the blocker is resolved.

This is public Operating Kit discovery. It must not include tenant details, account identifiers,
customer details, credentials, secrets, raw logs, local absolute paths, or private machine dumps.

## Goals

- Define whether Operating Kit should own generic local tooling and environment-readiness doctrine.
- Give consumer agents a clear route when an installed module is blocked by missing tooling.
- Separate generic bootstrap tooling from module-owned domain setup.
- Decide what kind of readiness state, if any, should be recorded and where it should live.
- Preserve approval gates before installing tools or changing the local machine.
- Keep module onboarding flexible while reducing duplicated package-manager and tool-install
  guidance across modules.
- Define the trigger behavior expected from human-facing and hybrid onboarding runbooks when local
  tooling blocks the flow.
- Produce an implementation-planning handoff if the discovery reaches a stable recommendation.

## Non-Goals

- Do not install Homebrew, PowerShell, package managers, CLIs, or PowerShell modules in this
  discovery.
- Do not define a default "install all baseline tools" bundle.
- Do not build an automated environment manager.
- Do not choose a single universal package manager for every platform.
- Do not make Operating Kit a broad wrapper around PowerShell, Graph, Exchange, SharePoint, cloud
  CLIs, or SDKs.
- Do not add Operating Kit runbooks for individual module tools such as M365 PowerShell modules,
  AWS CLI setup, or product-specific toolchains unless the tool family becomes a proven generic
  baseline lane.
- Do not move module-owned authentication, consent, tenant setup, service permissions, or live
  external validation into the Operating Kit.
- Do not record secrets, token caches, raw command output, full environment dumps, local absolute
  paths, or user account identifiers.
- Do not create a validator, schema, CLI command, or managed runbook before the discovery has been
  reviewed.
- Do not retrofit all existing module runbooks in the first discovery step.
- Do not retrofit Operating Kit first-run onboarding unless the implementation directly changes
  that runbook.

## Context And Evidence

| Evidence | Current Signal | Discovery Implication |
| --- | --- | --- |
| Runbook-authoring standards discovery | Agent-facing runbooks need preconditions, tool readiness checks, execution paths, stop conditions, and user-facing approval gates. | Tool readiness should be explicit enough that a fresh agent can continue when a required tool is missing. |
| Deferred feedback note | `runbook-authoring-standards/attachments/related-operating-kit-feedback.md` records "Shared Environment Readiness And Tooling Register" as `feedback-feature`, `feedback-doctrine`, and `needs-discovery`. | This discovery promotes that deferred feedback into a first-class planning artifact. |
| Native Codex capabilities profile | Operating Kit already has a narrow status-reporting model for capabilities such as documents, spreadsheets, browser, and PDF. | There is a precedent for timestamped local capability status, but it is not a general tooling-readiness system. |
| Module extension state routing | Committed module state belongs under `docs/repo/state/<id>/`, but local machine paths and preferences do not belong there. | Module state and local tool readiness must stay distinct. Module state can route the agent; local tool readiness is machine-specific evidence. |
| Foundry M365 module use case | The module may need PowerShell, Microsoft Graph PowerShell, PnP PowerShell, and Exchange Online PowerShell, while some users may lack prerequisite bootstrap tooling. | The M365 module should not own all generic package-manager and shell installation doctrine by itself. |
| M365 declaration pattern | The current M365 module expresses tooling needs through `module.yaml` lifecycle routes plus `reference/official-tooling.md`, `runbooks/connect-microsoft.md`, and onboarding Step 7. It does not expose one clean machine-readable `tools:` block. | V1 should allow declaration by manifest, reference, or runbook, as long as the module maps its needs to known Operating Kit tooling lanes clearly enough for agents. |
| Current managed route gap | Root `AGENTS.md` and the kit fallback inventory expose module state and runbook authoring routes, but they do not yet expose a tooling-readiness route. | V1 needs a short installed route so agents know where to go when a module onboarding or operation is blocked by local tooling. |
| M365 local tooling step | The current M365 onboarding runbook has a useful Step 7 local tooling check, but its missing-tool choice says "Install the Microsoft tools" rather than distinguishing package-manager, runtime, and PowerShell-module blockers. | The shared route should provide the concrete blocker-specific user choices and then return to the module flow. |
| Operating Kit first-run onboarding | The first-run onboarding runbook is a strong human-facing example, but it predates the compact intention block and is not an environment-readiness runbook. | Do not retrofit it as part of tooling readiness unless directly touched; use it as a quality reference for pacing and user-owned decisions. |
| Consumer-agent visibility need | Installed modules may be the main surface local users interact with. If a module fails due to missing tooling, agents need a standard recovery route. | The route cannot be buried only in one module; it must be visible through managed Operating Kit guidance. |

## Terminology

`tool family`: a reusable class of local tool, such as package manager, shell runtime, PowerShell,
Node.js, Python, browser automation, cloud CLI, document converter, or platform-specific installer.

`baseline tooling catalog`: the Operating Kit list of known on-demand tooling lanes. It tells
agents how to recognize, explain, check, install after approval, repair, or block on a tool family.
It is not a default installation bundle.

`tooling lane`: a reusable check/install/repair route for a tool family, such as `powershell-runtime`
or `powershell-module`.

`module-owned tool`: a tool or extension that is only meaningful for a module's domain, such as a
domain-specific PowerShell module or service CLI.

`declaration by runbook or reference`: a module declaration pattern where tool requirements are not
machine-readable in one manifest field, but are still explicit in required module runbooks or
references.

`readiness check`: a local, read-only check that determines whether a tool appears available,
missing, blocked, stale, or unknown.

`environment blocker`: a local-machine or tooling condition that prevents a module from continuing,
such as a missing package manager, shell runtime, CLI, PowerShell module, path entry, or approved
install route.

`readiness state`: timestamped local evidence about a tool check. It may route the next step, but
it is not proof that a tool still works and does not authorize installs or external changes.

`install or repair runbook`: a human-facing or hybrid runbook that explains the purpose of a tool,
asks for approval before local changes, and guides the user through platform-appropriate setup or
repair.

`live external preflight`: a module-owned check against the external service or tenant before
sensitive reads, writes, consent changes, resource creation, or destructive actions.

## Boundary Model

Tooling readiness has three separate layers:

| Layer | Owner | Examples | Intended Operating Kit Role |
| --- | --- | --- | --- |
| Generic bootstrap environment | Operating Kit, with user approval for changes | package manager availability, shell/runtime prerequisites, permission prompts, local install safety | Define doctrine, checks, user-facing install or repair patterns, and routing. |
| Reusable tool family | Shared; Operating Kit may own common patterns, modules declare need | PowerShell runtime, Node.js, Python, browser automation, cloud CLI family | Provide generic readiness language and optional common runbooks where repeated use proves value. |
| Domain/module tooling and access | Module or product owner | Microsoft Graph PowerShell, PnP PowerShell, Exchange Online PowerShell, service-specific CLIs, tenant consent | Module declares tools, explains capability unlocked, validates versions, handles service auth and live preflight. |

This boundary is meant to prevent two failure modes:

- every module duplicating generic package-manager and tool-install guidance;
- Operating Kit becoming responsible for every domain-specific service module and permission model.

It is acceptable for a module to have module-specific install runbooks. The boundary problem is
Operating Kit runbook sprawl, not module specificity.

## Trigger Model

The tooling readiness route should trigger when:

- module onboarding reaches a missing local package manager, runtime, CLI, PowerShell module, or
  install-path blocker;
- module operation reaches a missing local package manager, runtime, CLI, PowerShell module, or
  install-path blocker;
- the user explicitly asks the agent to install or repair a tool needed by a repository or module;
- an agent-facing runbook declares a required tool and the read-only check reports `missing`,
  `blocked`, or `unknown`.

After the trigger:

1. classify the local blocker by tooling lane;
2. present a concrete nontechnical user choice;
3. ask explicit approval before local installation or repair;
4. run only the approved install or repair path;
5. recheck readiness;
6. return to the module runbook when resolved;
7. record a capability blocker and stop when unresolved.

The tooling route should not trigger for module-owned live external blockers. Microsoft tenant
consent, admin roles, mailbox access, SharePoint permissions, licenses, app readiness, and similar
service-state problems remain module-owned preflight blockers.

## Decision Inventory

### D-001 - Operating Kit Should Own A Generic Readiness Route

Status: recommended

Decision: Operating Kit should provide a visible, generic route for missing local tooling and
environment readiness. When a module is blocked by a missing tool, agents should know to follow the
managed readiness route before inventing ad hoc setup advice.

Why it matters: module users should get a consistent, nontechnical recovery flow when the local
machine is missing a prerequisite. This is especially important for consumers who do not know what a
package manager, shell runtime, or CLI module is.

Recommended default: add this as managed agent-interface or structure-governance doctrine in a
later implementation, with a generic root route rather than one route per tool or module.

### D-002 - Modules Should Declare Concrete Tool Needs

Status: recommended

Decision: modules should declare required and optional tools in their own module documentation,
manifest, references, or runbooks. V1 should not require every module to use the same
machine-readable shape. Each declaration should explain what capability the tool unlocks, whether
it is required for all onboarding or only for specific operations, and which Operating Kit tooling
lane or module-owned install route to use when it is missing.

Why it matters: Operating Kit can route missing-tool handling, but only the module knows whether
PowerShell, a PowerShell module, a cloud CLI, or a browser tool is required for a specific domain
operation.

Recommended default: future module runbooks should include a small tool declaration or preflight
section rather than embedding a full generic install guide. Existing modules may use declaration by
runbook or reference when the route is clear.

### D-003 - Keep Machine Readiness Separate From Committed Module State

Status: recommended

Decision: committed state under `docs/repo/state/<module-or-extension-id>/` is for non-secret
repo-owned module routing context. Machine-specific readiness observations should not be stored
there as durable committed truth.

Why it matters: a repository state file may be shared by multiple agents or machines. Whether a
specific laptop has Homebrew, PowerShell, or a CLI module installed is local and can change without
a repository change.

Recommended default: treat tool availability as timestamped local/generated evidence. The exact
placement for this evidence remains an implementation-shaping open question.

### D-004 - Separate Declarations, Observations, And Instructions

Status: recommended

Decision: the model should separate:

- tool declarations: durable requirements supplied by Operating Kit profiles, modules, or
  repository guidance;
- readiness observations: local, timestamped check results;
- install and repair instructions: managed or module-owned runbooks;
- live external validation: module-owned preflight against the real external service.

Why it matters: conflating these artifacts causes stale records, unsafe automation, and confusing
user-facing flows.

Recommended default: use the model above as the first implementation-planning constraint.

### D-005 - Approval Gates Are Required Before Local Installs

Status: recommended

Decision: Operating Kit readiness guidance must never authorize silent installation or repair.
Agents may run read-only checks when locally appropriate, but package-manager installs, shell
runtime installs, PowerShell module installs, path changes, permission prompts, and similar local
changes require explicit user approval.

Why it matters: local machine setup changes can affect security, compliance, disk state, developer
toolchains, and user trust.

Recommended default: human-facing install or repair runbooks should explain the intended outcome,
what will change locally, how to stop, and what evidence will be checked after approval.

### D-006 - Standardize Tooling Without Over-Standardizing Platforms

Status: recommended

Decision: Operating Kit should define standard concepts and preferred support lanes, but it should
not assume every consumer already has the same package manager, shell runtime, platform, or install
path.

Why it matters: a useful standard must know what to do when a module asks for PowerShell but the
machine does not yet have the expected bootstrap tooling. The first answer should be a clear
tooling lane, not a module-specific guess.

Recommended default: define platform-aware on-demand lanes such as package-manager/bootstrap,
PowerShell runtime, and PowerShell module install. A later implementation can choose which lanes are
officially supported first.

### D-007 - Readiness Does Not Replace Live Preflight

Status: recommended

Decision: a tool readiness check can show that a command exists or a module appears installed, but
it does not prove that service authentication, tenant permissions, mailbox access, SharePoint
access, admin roles, or external resource state are valid.

Why it matters: a successful local tool check is not enough to read or change external systems.

Recommended default: module runbooks must still perform live external preflight before sensitive
reads, writes, or tenant-changing actions.

### D-008 - Baseline Tooling Means On-Demand Catalog, Not Default Install

Status: recommended

Decision: Operating Kit should define a baseline tooling catalog that agents can use on demand. It
should not require every consumer repository to install every baseline tool during kit onboarding.

Why it matters: a flat install list would add local-machine risk and noise before a module actually
needs a tool. A pure ad hoc model would leave agents confused and produce inconsistent install
guidance. The on-demand catalog gives standard behavior without unnecessary setup.

Recommended default: the catalog should include only widely reusable lanes and active module
blockers. Tools are checked or installed only when the user request or module runbook needs them.

### D-009 - Define A Missing-Tool Agent Behavior Contract

Status: recommended

Decision: Operating Kit should define the generic sequence agents follow when a module-required
tool is missing:

1. identify the module-requested tool and capability;
2. map it to a known Operating Kit tooling lane when possible;
3. run read-only local checks;
4. explain the user-facing purpose and local effect;
5. request explicit approval before install or repair;
6. run the approved generic lane or module-owned install command;
7. recheck readiness;
8. return to the module runbook;
9. record a capability blocker when readiness remains unavailable.

Why it matters: this is the reusable behavior that prevents each module from inventing a different
missing-tool conversation.

Recommended default: ship this first as managed doctrine/runbook guidance. Add schema or validator
support only after repeated modules need it.

### D-010 - Environment Blockers Should Trigger The Readiness Route

Status: recommended

Decision: any module onboarding or operation runbook that encounters an environment blocker should
route through the Operating Kit missing-tool behavior contract before giving up, improvising install
steps, or asking the user to solve the problem unaided.

Why it matters: missing tooling will often surface inside a module flow, not during standalone
environment setup. The Operating Kit route must be available at the moment of failure.

Recommended default: module runbooks should state their environment blockers and route them to the
managed readiness contract. The agent should then return to the module runbook after the blocker is
resolved or record a capability blocker when it is not.

### D-011 - Make The Route Visible In Installed Operating Kit Guidance

Status: recommended

Decision: V1 should add a short installed route for tooling readiness in the managed root
`AGENTS.md` block and the kit fallback inventory. The detailed catalog and missing-tool behavior
should live in managed Operating Kit documentation, not in the root block.

Why it matters: environment blockers will often appear while an agent is inside a module runbook.
If the installed root route does not mention tooling readiness, the agent may stay trapped in the
module-local guidance or improvise a fix.

Recommended default: use one concise root rule that says to follow the Operating Kit tooling
readiness route before installing, repairing, or improvising local tools for a module.

### D-012 - Human And Hybrid Runbooks Should Use Blocker-Specific Choices

Status: recommended

Decision: human-facing and hybrid onboarding runbooks should not collapse all missing tools into
one generic "install tools" choice when the blocker can be classified. They should present the
next concrete blocker-specific action, such as installing Homebrew, installing PowerShell,
installing Microsoft 365 PowerShell tools, using the user's own method, using a documented fallback,
or stopping.

Why it matters: nontechnical users can choose concrete outcomes. They should not have to interpret
tooling lanes, package-manager/bootstrap language, or module internals.

Recommended default: the managed readiness route owns the blocker-specific user copy. Module
runbooks may keep a simple local tooling step, but should route environment blockers to the managed
readiness contract instead of duplicating full install UX.

### D-013 - Keep Operating Kit Tooling Guidance Centralized

Status: recommended

Decision: Operating Kit should provide one central tooling-readiness route, one small baseline
catalog, and one missing-tool behavior contract. It should not add separate managed Operating Kit
runbooks for each module-specific tool or ecosystem.

Why it matters: modules may legitimately have very different tool requirements. Forcing all
module-specific install detail into Operating Kit would create managed runbook sprawl and make the
kit harder to maintain.

Recommended default: Operating Kit owns generic lanes and behavior. Modules own module-specific
install commands, version checks, fallback commands, and service-specific tool caveats.

## Requirements And Evaluation Criteria

### FR-001 - Visible Missing-Tool Route

Consumer agents must have a visible managed route for missing local tools before they ask the user
to solve the problem manually or declare the capability unavailable.

### FR-002 - Tool Declaration Shape

The later implementation should define a lightweight declaration expectation that can be satisfied
through a manifest field, reference document, or runbook section. It should express:

- tool ID or family;
- required or optional status;
- capability unlocked;
- Operating Kit tooling lane when one exists;
- read-only readiness check;
- missing, blocked, and unknown handling;
- install or repair runbook route;
- module-owned live preflight route.

### FR-003 - Human-Facing Install And Repair Flow

Install or repair guidance must be suitable for nontechnical users:

- explain the outcome in plain language;
- ask one decision at a time when possible;
- present concrete action choices, not abstract lanes;
- make approval gates explicit;
- avoid leading with internal file paths or implementation mechanics;
- tell users where to look when they need a value or system setting;
- stop safely and record a capability blocker when installation is unavailable.

For example, a macOS package-manager blocker should offer choices in this shape:

```text
This Mac needs Homebrew before I can install the requested tool.

How should we continue?

1. Install Homebrew (Recommended)
2. I will install it another way
3. Stop here
```

Do not ask a nontechnical user to choose a "package-manager/bootstrap lane".

### FR-004 - Local Readiness Evidence Contract

If readiness state is recorded, it must be limited to non-secret summary evidence such as:

- tool ID;
- status;
- checked-at timestamp;
- platform category;
- check method category;
- version when useful and safe;
- caveats and next route.

It must not record secrets, tokens, local absolute paths, raw command logs, full environment dumps,
account identifiers, tenant identifiers, or customer details.

### FR-005 - Module Boundary

The standard must make clear that modules own domain-specific tooling, auth, permissions, service
preflight, and operation recipes even when they reuse generic Operating Kit readiness routes.

### FR-006 - Low-Ceremony First Implementation

The first implementation should be able to ship as instruction-only or mostly instruction-only
unless discovery proves a schema, validator, or CLI command is necessary for agents to use the
route reliably.

### FR-007 - On-Demand Baseline Catalog

The first implementation should define a small baseline tooling catalog. The catalog should not
install tools by default. It should give agents a known route for module-triggered checks,
approval-gated installs, repairs, and blockers.

### FR-008 - Missing-Tool Behavior Contract

The first implementation should define the generic missing-tool sequence that modules can invoke or
reference without copying generic install doctrine.

### FR-009 - Module Environment Block Trigger

Module runbooks should be able to invoke the managed readiness route whenever onboarding or
operation is blocked by the local environment. The route must then either resolve the blocker and
return to the module flow, or record a capability blocker and stop cleanly.

### FR-010 - Installed Route Visibility

The first implementation should make tooling readiness discoverable from installed Operating Kit
routes, including the managed root `AGENTS.md` block and the kit fallback inventory. The root route
should stay concise and generic.

### FR-011 - Blocker Ownership Split

The first implementation should clearly distinguish:

- local environment blockers handled by the Operating Kit readiness route;
- module-owned service blockers handled by module live preflight.

For the M365 example, missing Homebrew, PowerShell, and PowerShell modules are local environment
blockers. Microsoft tenant consent, admin role, mailbox access, SharePoint permission, license, app
readiness, and Microsoft feature availability remain M365 module blockers.

### FR-012 - Operating Kit Anti-Sprawl Boundary

The first implementation should avoid creating multiple managed Operating Kit runbooks for
individual module tools. Tool-specific Operating Kit content should be limited to generic baseline
lanes. Module-specific install details remain in module-owned runbooks and references.

## Open Questions

### OQ-001 - Where Should Readiness Observations Live?

Owner: kit maintainer and implementation planner

BLOCKER: yes, before implementation

Options:

- extend `.codeheart/kit.lock.yaml` for Operating Kit-owned checks;
- use an ignored local user file under `.codeheart/user/` for machine-local readiness notes;
- generate reports under a report/evidence path;
- keep readiness ephemeral and rerun checks every time.

Current recommendation: do not commit machine-specific readiness as module state under
`docs/repo/state/<id>/`. Choose the narrowest placement that supports agent routing without turning
stale local state into false truth.

### OQ-002 - Which Tool Families Should V1 Cover?

Owner: kit maintainer and module owners

BLOCKER: yes, before implementation

Candidate starting set:

- package-manager/bootstrap lane;
- PowerShell runtime;
- PowerShell module lane;
- Node.js and package manager lane;
- Python runtime lane;
- browser automation lane;
- document/PDF conversion lane;
- cloud CLI family as a later candidate.

Current recommendation: v1 should include package-manager/bootstrap, PowerShell runtime,
PowerShell module, Node.js, Python, browser automation, and document/PDF lanes as the initial
catalog. Cloud CLIs should stay later unless an active module needs them. The M365 module maps to
the PowerShell runtime and PowerShell module lanes for Graph, PnP, and Exchange module readiness.

### OQ-003 - Should Tool Declarations Be Machine-Readable In V1?

Owner: kit maintainer and implementation planner

BLOCKER: no, unless automatic checks are planned

Options:

- Markdown-only declaration standard in module runbooks;
- YAML declaration file in module packages;
- Operating Kit schema for tool requirements;
- hybrid Markdown first, schema later.

Current recommendation: do not require machine-readable declarations in v1. Accept declaration by
runbook or reference when it is explicit and maps to a known Operating Kit lane. Revisit a YAML or
schema shape after repeated modules expose the same needs.

### OQ-004 - How Much Install Guidance Should Operating Kit Own?

Owner: kit maintainer

BLOCKER: yes, before implementation

Options:

- only define routing doctrine and let modules own all install steps;
- define generic tool-family install and repair runbooks for common prerequisites;
- define platform package-manager bootstrap guidance;
- define a full managed tool catalog.

Current recommendation: Operating Kit should own generic package-manager/bootstrap and common tool
family patterns where repeated modules need them. For package-manager/bootstrap, v1 should include
thin but concrete human-facing guidance, such as "Install Homebrew", "I will install it another
way", and "Stop here". Modules should still own domain-specific modules, service authentication,
and service permissions.

### OQ-006 - Should V1 Touch Existing Human-Facing Onboarding Runbooks?

Owner: kit maintainer and implementation planner

BLOCKER: no

Question: should the first tooling-readiness implementation retrofit existing Operating Kit
first-run onboarding or module onboarding runbooks?

Current recommendation: do not retrofit Operating Kit first-run onboarding unless it is directly
touched by the implementation. Update generic standards and routes first. Apply the shared
readiness route to module runbooks, such as the M365 onboarding runbook, through a module-specific
follow-up plan after the Operating Kit route exists.

### OQ-005 - How Should Restricted Or Policy-Managed Devices Be Handled Later?

Owner: kit maintainer and consumer operator

BLOCKER: no

Question: how should guidance handle devices where users cannot install package managers or tools
because of endpoint management, administrator rights, or company policy?

Current recommendation: defer detailed restricted-device flows. V1 only needs a simple stop path:
do not attempt workarounds, record a capability blocker, and return to the module's fallback or
stop condition.

## Risks

### R-001 - Operating Kit Becomes A Tool Wrapper

Likelihood: medium

Impact: high

Risk: a broad tooling model could expand into wrappers around every CLI or SDK.

Mitigation: keep the first implementation focused on readiness, user guidance, and routing.
Modules remain responsible for execution recipes.

### R-002 - Stale Local Readiness Causes Bad Decisions

Likelihood: medium

Impact: high

Risk: an old readiness record says a tool exists, but the live environment changed.

Mitigation: readiness records need timestamps and caveats. Agents must rerun read-only checks
before relying on local tools for meaningful work.

### R-003 - Private Machine Details Leak Into Public Repos

Likelihood: medium

Impact: high

Risk: readiness output could accidentally include usernames, local paths, machine details, raw
logs, or account identifiers.

Mitigation: record only summary fields. Exclude raw logs and local paths. Keep public examples
generic.

### R-004 - Modules Still Duplicate Generic Guidance

Likelihood: medium

Impact: medium

Risk: if the route is not visible, module authors will keep embedding their own generic install
instructions.

Mitigation: make the managed route discoverable from root agent guidance, runbook authoring
standards, and module onboarding expectations.

### R-005 - Users Receive Too Much Technical Detail

Likelihood: medium

Impact: medium

Risk: human-facing onboarding may mention package managers, module manifests, or local state before
the user understands the outcome.

Mitigation: install and repair runbooks must follow the human-facing runbook standard: outcome
first, one decision at a time, internal mechanics outside first-turn copy.

## Implementation Handoff Candidate

This discovery will be ready for implementation planning when:

- the readiness observation placement is chosen;
- the first supported tool-family scope is chosen;
- the on-demand baseline catalog shape is accepted;
- the missing-tool behavior contract is accepted;
- the module environment-block trigger is accepted;
- the concrete nontechnical package-manager choice shape is accepted;
- the install-guidance ownership boundary is accepted;
- the consumer-visible route is selected;
- the blocker ownership split between Operating Kit local environment blockers and module-owned
  service blockers is accepted;
- the Operating Kit anti-sprawl boundary is accepted;
- the implementation impact class is confirmed.

Likely implementation surfaces:

- managed agent-interface reference or runbook for tooling readiness;
- managed structure-governance reference update for local/generated readiness state placement;
- runbook-authoring standard cross-reference for tool readiness sections;
- root `AGENTS.md` managed block route;
- kit fallback inventory route;
- optional package resource mirroring;
- optional tests for sync/check packaged resources.

Do not add separate Operating Kit runbooks for M365 PowerShell modules, AWS CLI, or other
module-specific tools in the first implementation. Those belong to their modules unless later
evidence shows that a tool family should graduate into a generic baseline lane.

Deferred follow-up surfaces:

- module-specific update to Foundry M365 onboarding so Step 7 invokes the managed readiness route
  for environment blockers and keeps Microsoft tenant/service blockers module-owned;
- optional retrofit of Operating Kit first-run onboarding only if future implementation touches that
  runbook for another reason.

Likely consumer impact: `instruction-only change` unless the implementation adds generated
readiness state, schemas, validators, CLI behavior, or new scaffolded paths.

## Revision Notes

- 2026-06-24: Created first discovery draft from deferred shared environment-readiness feedback
  and the Foundry Microsoft 365 module tooling-readiness discussion.
- 2026-06-24: Refined the model from a possible flat tooling list into an on-demand baseline
  tooling catalog plus a generic missing-tool agent behavior contract; recorded flexible module
  declaration by manifest, reference, or runbook.
- 2026-06-24: Clarified that module onboarding and operation environment blockers should trigger
  the readiness route, and that package-manager bootstrap guidance needs concrete nontechnical
  choices such as installing Homebrew, choosing another install method, or stopping.
- 2026-06-24: Added trigger-model review results from managed Operating Kit routes and the Foundry
  M365 hybrid onboarding runbook: route visibility is required, local environment blockers must be
  separated from module-owned service blockers, and blocker-specific human choices should replace
  broad missing-tool prompts.
- 2026-06-24: Clarified the anti-sprawl boundary: Operating Kit should keep one central tooling
  readiness route and small baseline catalog, while modules may keep module-specific install
  runbooks and commands.
