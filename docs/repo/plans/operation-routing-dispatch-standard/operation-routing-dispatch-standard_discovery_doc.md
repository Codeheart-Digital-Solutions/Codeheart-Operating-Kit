Last updated: 2026-06-25T12:18:33Z (UTC)
Created: 2026-06-25
Status: draft

# Operation Routing And Dispatch Standard Discovery

## Overview

This discovery investigates whether Codeheart Operating Kit should define a reusable routing and
dispatch standard for agent-operated work.

The problem is broader than any one external system, connector, module, or runbook. Agents can
enter a repository, see several possible documents, modules, connectors, tools, scripts, browser
surfaces, or APIs, and choose an execution path before resolving which source actually has
authority over the request. That creates inconsistent behavior, especially when a deeper module or
domain can handle an intent that is not obvious from the top-level router.

Recommended current direction: Operating Kit should own a generic routing doctrine and dispatch
sequence. Domain, module, product, and repository owners should expose compact capability
advertisements at their ownership boundaries, then keep detailed route registries and route cards
inside the owning domain. Upper layers should route and advertise capability families; they should
not duplicate every deep operational detail.

This is public Operating Kit discovery. It must not include private tenant details, account
identifiers, credentials, local machine state, raw operational logs, or private business content.

## Goals

- Define a generic pre-execution routing model for agents.
- Clarify how agents choose between user instructions, repository instructions, Operating Kit
  doctrine, installed modules or extensions, repo-owned docs, official external docs, visible
  tools, connectors, APIs, browser paths, and model knowledge.
- Define how upper layers can expose enough capability information for agents to discover deeper
  routes without becoming encyclopedias.
- Define where routers, capability advertisements, route registries, and route cards should live.
- Inventory existing Operating Kit instruction and routing surfaces before proposing new managed
  doctrine.
- Decide which existing surfaces should remain as-is, point to the new standard, be consolidated,
  or be retired during implementation.
- Define which task categories require explicit routing checks and which can stay lightweight.
- Define a simple conflict-resolution rule for competing authorities.
- Define when fresh-agent routing tests should be used.
- Define when capability advertisements must be reviewed or updated.
- Define an implementation-planning hook so routing-bearing runbook or route-surface epics include
  fresh low-context routing probes.
- Define the split between a compact root `AGENTS.md` routing hierarchy and the detailed managed
  routing doctrine.
- Keep module and domain owners responsible for concrete domain routes, execution surfaces,
  preconditions, recipes, and evidence fields.
- Produce a later implementation handoff for managed agent-interface and structure-governance
  guidance.

## Non-Goals

- Do not implement the routing standard in this discovery.
- Do not create a universal plugin, module, connector, MCP, or execution-surface framework.
- Do not define every possible external-service route.
- Do not replace module-owned runbooks, schemas, consent models, or live preflight rules.
- Do not make the root `AGENTS.md` block list every capability of every installed module.
- Do not restructure all existing Operating Kit runbooks and references in this discovery.
- Do not define detailed operational recipe maturity, phase-marker, script-promotion, or test
  shapes here. Route cards may point to canonical recipes, but recipe maturity is a separate
  follow-up topic.

## Context And Evidence

| Evidence | Current Signal | Discovery Implication |
| --- | --- | --- |
| Agent-interface docs | The managed root route is intentionally concise and points agents to task-matched Operating Kit routes. | Root guidance should remain a router, not a catalog of all operational details. |
| Structure-governance docs | README files route, references own durable doctrine, runbooks own procedure, and state belongs under the owner boundary. | Routing doctrine needs placement rules for routers, capability advertisements, route registries, and route cards. |
| Module-extension state routing | Committed module state helps route agents, but live preflight still decides before sensitive reads or changes. | Existing state routing is one instance of the broader routing model, not the whole model. |
| Runbook authoring standard | Agent-facing runbooks already need source of truth, execution lane, preconditions, approval gates, stop conditions, evidence, and validation. | Runbook authoring has pieces of dispatch discipline, but it does not yet define generic route discovery before runbook selection. |
| Tooling-readiness route | Operating Kit separates generic local tooling blockers from module-owned service blockers. | This is a useful ownership split pattern: Operating Kit owns generic behavior; modules own domain-specific preflight. |
| Codex capability surfaces | Codex can operate through instructions, skills, plugins, apps, MCP servers, shell commands, browser surfaces, and local tools. | Visible tools or connectors are execution-surface options, not automatically the highest routing authority. |
| Repeated module use cases | A deep module may handle intents such as workspace, document, mailbox, issue, release, or account work that a top-level router cannot enumerate fully. | Each domain needs a compact capability advertisement so agents know when to descend into its route registry. |
| Existing Operating Kit evolution | Routing-related instructions already appear across root managed routes, agent-interface docs, structure governance, module-extension state, tooling readiness, runbook authoring, and planning workflows. | Implementation should start with an instruction inventory and disposition map so the new standard consolidates or links existing doctrine instead of layering conflicting guidance on top. |

## Terminology

`router`: a concise document or section that tells agents where to look next. `AGENTS.md` files,
README files, and module READMEs commonly act as routers.

`capability advertisement`: a compact boundary summary that states what a domain, module,
product, or repository area can handle. It should list capability families, intent aliases, scope
families, route-registry location, state location when applicable, and ambiguity rules. It should
not duplicate detailed route cards or recipes.

`route registry`: the domain-owned index of request families. A registry may be a reference file,
capability map, table, or set of route-card sections.

`route card`: the dispatch contract for one request family. It tells the agent which authority,
scope, state sources, execution surface, preconditions, stop conditions, and recipe apply before
execution starts.

`execution surface`: the concrete surface used after routing is complete, such as shell, API,
connector, MCP tool, browser, admin portal, PowerShell, CLI, script, document surface, or manual
user action.

`canonical recipe`: the detailed runbook, script, API procedure, portal procedure, or manual path
entered after route selection and precondition checks.

## Candidate Routing Model

The reusable routing sequence should be:

```text
Intent
-> Domain
-> Authority source
-> Scope
-> Capability route
-> Execution surface
-> Preconditions
-> Canonical recipe
-> Evidence or blocker
```

Core rule:

```text
Route through the highest-authority applicable operating source before choosing an execution
surface.
```

Visible tools, connectors, and local commands are important, but they belong after intent,
domain, authority, scope, and route selection. They should not be used as the first source of
truth unless the user explicitly requests that surface or no higher-authority route applies.

## Routing Trigger Categories

The routing standard should scale by task category. Not every action needs a full route card, but
some categories should always trigger a routing check.

`mandatory routing check`

Use a route or owner check before acting when the request involves:

- creating, moving, deleting, or reorganizing durable files, folders, docs, route surfaces, or
  structure;
- changing managed Operating Kit content, generated consumer surfaces, or root routing;
- installed modules, extensions, or committed module state;
- external systems, tenants, accounts, workspaces, permissions, approvals, or live service state;
- sensitive reads, writes, deletes, permission changes, offboarding, diagnostics, or compliance
  surfaces;
- repeated operational workflows that already have or likely need a route card.

`conditional routing check`

Use a quick owner and ambiguity check when the request has multiple plausible domains, providers,
tools, products, or execution surfaces. If the owner is clear and the work is low risk, continue
with the nearest applicable runbook or local convention.

`lightweight direct work`

Tiny local edits, simple command answers, and ordinary code changes may proceed without formal
route-card lookup when the owner and conventions are already clear. Existing repository
instructions, local code patterns, and relevant validation still apply.

For new durable files or folders, the lightweight path should not bypass structure governance. The
agent should at least confirm the owner and placement rule before creating the artifact.

## Authority Conflict Rule

The first implementation should keep conflict handling simple.

When routing authorities conflict:

1. preserve safety and public-core boundaries first;
2. follow explicit user instruction when it is safe and within the user's authority;
3. follow the nearest legitimate owner for the artifact, domain, or operation;
4. use official external truth for current external-system behavior;
5. stop before sensitive reads, writes, deletes, permission changes, or external changes when the
   conflict affects authority, target, approval, or live state.

Do not silently resolve a conflict by choosing the most convenient execution surface.

## Candidate Routing Authority Hierarchy

This hierarchy is an Operating Kit routing heuristic. It does not replace Codex's native
instruction priority or safety rules. Use it after the active instruction stack is already in
force, when an agent needs to decide which repository, Operating Kit, module, product, domain,
external, or tool source should route a request.

The first implementation should treat this as the candidate hierarchy for routing authority:

1. Explicit user instruction, when safe and within the user's authority.
2. Safety, public-core, secrets, destructive-action, and external-state rules.
3. Repository instructions and nearest applicable `AGENTS.md` guidance.
4. Operating Kit generic doctrine for routing, structure, planning, runbook shape, local tooling,
   and managed boundaries.
5. Installed module, extension, product, package, or source-area route owned by the request
   domain.
6. Consumer repo-owned references, route registries, runbooks, and committed non-secret routing
   state.
7. Official external documentation or live external truth for current external-system behavior.
8. Visible execution-surface documentation, such as tool, connector, MCP, CLI, browser, or API
   capability notes.
9. General model knowledge only when no higher-authority route or current source applies.

This hierarchy is not a license to ignore lower layers. It is an order for resolving routing
authority inside the active repository context. Lower layers can still provide required target
detail, execution constraints, preflight results, or current external facts. Live external
preflight, approval gates, and current official documentation may still stop or reshape execution
after a route is selected.

## Layer Model

The routing model should not require every layer to carry every deep detail. Instead:

- upper layers expose routers and capability advertisements;
- owning domains expose route registries and route cards;
- canonical recipes hold concrete execution procedure.

Example generic shape:

```text
AGENTS.md
  -> high-level bootstrap router

.codeheart/kit/
  -> generic routing doctrine and route-card standard

docs/repo/
  -> repo-owned routes, local exceptions, and state pointers

<installed-module-or-extension>/
  -> module capability advertisement, route registry, route cards, and recipes

products/<product>/
  -> product-owned route cards only when repeated product operations justify them
```

Route cards should live at the ownership boundary of the work. A repository does not need one
route card per layer. It needs routers at layer boundaries and route cards only where the owner
has repeated operational routing responsibility.

## Capability Advertisement Concept

Capability advertisements solve the discoverability problem without pushing every detail upward.

A parent router should not need to know every deep command, permission, API endpoint, or fallback
path. It should be able to see enough to decide whether a deeper owner may handle the request.

Minimum V1 advertisement fields:

```text
Owner:
Domain:
Capability families:
Intent aliases:
Route registry:
Fallback or ambiguity rule:
```

Optional V1 advertisement fields:

```text
Scope families:
State location:
```

Example, generalized:

```text
Owner: <workspace module>
Domain: external workspace operations
Capability families:
- users, groups, and access
- workspace files and folders
- named and shared communication resources
- drafts, sends, diagnostics, and validation
Route registry: <module route registry>
State location: docs/repo/state/<module-id>/ when present
Fallback or ambiguity rule: ask when more than one provider or domain can satisfy the intent.
```

This lets an agent match a plain-language request to candidate domains before choosing an
execution surface.

## Decision Inventory

### D-001 - Operating Kit Should Own Generic Routing Doctrine

Status: recommended

Decision: Operating Kit should define the generic pre-execution routing and dispatch doctrine.
Domain, module, product, and repository owners should instantiate that doctrine for their own
operational surfaces.

Why it matters: agents need one reusable mental model for deciding where a request belongs before
selecting tools, scripts, connectors, APIs, or runbooks.

Recommended default: add the doctrine under managed agent-interface guidance, with
structure-governance cross-references for placement.

### D-002 - Route Before Selecting Execution Surface

Status: recommended

Decision: agents should resolve intent, domain, authority, scope, and capability route before
selecting the execution surface.

Why it matters: visible execution surfaces can be incomplete, misleading, unavailable, or less
authoritative than repository or module routes.

Recommended default: execution surfaces should be selected by route cards or canonical runbooks,
unless the user explicitly requests a specific surface.

### D-003 - Use Capability Advertisements At Ownership Boundaries

Status: recommended

Decision: domains, modules, products, and repo areas that own repeated operations should expose a
compact capability advertisement at their boundary.

Why it matters: upper routers cannot enumerate every deep route, but agents still need enough
information to discover that a deeper owner may handle a plain-language intent.

Recommended default: advertisements should live in the owning domain's primary README or
route-registry introduction in V1. Do not require a machine-readable manifest until repeated usage
proves it useful. Advertisements should list capability families, intent aliases, scope families,
route-registry location, state location when applicable, and ambiguity rules. V1 should require
owner, domain, capability families, intent aliases, route registry, and ambiguity rule. Scope
families and state location can be optional until repeated usage proves they are always needed.

### D-004 - Keep Full Route Cards With The Owning Domain

Status: recommended

Decision: full route cards should live where the domain or operation owner lives, not duplicated
in every parent router.

Why it matters: duplicating route details upward makes routers stale and hard to maintain. Keeping
cards with the owner preserves local authority and reduces drift.

Recommended default: parent routers link to route registries and carry only capability
advertisements or route pointers.

### D-005 - Define A Route-Card Standard Without Defining Domain Details

Status: recommended

Decision: Operating Kit should define the route-card shape, but not the domain-specific content
of every card.

Candidate field set:

```text
Route ID:
Intent patterns:
Domain:
Lifecycle:
Scope:
Action class:
Authority source:
State sources:
Live truth source:
Default execution surface:
Fallback surface:
Preconditions:
Approval class:
Canonical runbook or reference:
Stop conditions:
Evidence fields:
```

Why it matters: the shared shape lets agents and reviewers recognize dispatch contracts across
domains, while domain owners keep control of the actual behavior.

Recommended default: treat this as a Markdown authoring standard first. Do not require a
machine-readable route-card schema until repeated usage proves it is needed. Route-card fields may
be marked `not applicable` with a short reason when a field does not fit the route.

### D-006 - Make Ambiguity Explicit

Status: recommended

Decision: when multiple capability advertisements or route cards plausibly match a request, the
agent should use explicit user instruction and local context to choose. If neither decides, the
agent should ask a targeted provider, domain, or target question before execution.

Why it matters: many plain-language requests can belong to multiple domains or providers.
Guessing creates wrong-surface execution and weak approval boundaries.

Recommended default: every route registry should include a short ambiguity rule for common
overlaps.

### D-007 - Preserve Live Preflight And Approval Boundaries

Status: recommended

Decision: routing state, route cards, and capability advertisements help choose the path. They do
not authorize sensitive reads, writes, deletions, permission changes, or external changes.

Why it matters: routing artifacts can become stale or incomplete. Live external systems and
explicit approvals still decide whether execution may proceed.

Recommended default: route cards should identify live truth sources, preconditions, approval
class, and stop conditions.

### D-008 - Treat Recipe Maturity As Adjacent, Not In Scope

Status: recommended

Decision: this routing standard should acknowledge that route cards dispatch to canonical
recipes, but should not define recipe phase models, script promotion, marker taxonomies, or
automation maturity standards.

Why it matters: routing and recipe maturity are related but separable. Routing answers which path
to take; recipe maturity answers how the selected path is structured and hardened.

Recommended default: record operational recipe maturity as a follow-up discovery area after the
routing standard is accepted or stable enough to depend on.

### D-009 - Inventory Existing Routing Doctrine Before Implementation

Status: recommended

Decision: before implementing new managed routing doctrine, the planner should inventory existing
Operating Kit instruction and routing surfaces and assign each one a disposition.

Candidate disposition values:

- `leave-as-is`;
- `link-to-new-standard`;
- `consolidate-into-new-standard`;
- `keep-as-domain-specific-detail`;
- `defer`;
- `retire-or-reword`.

Why it matters: Operating Kit is evolving and already contains partial routing rules. Without an
inventory, the implementation could add a new standard while leaving older guidance ambiguous,
duplicative, or subtly conflicting.

Recommended default: make the inventory an implementation-planning prerequisite. It should cover
managed root `AGENTS.md` wording, `.codeheart/kit/README.md`, agent-interface references and
runbooks, structure-governance references, planning-workflow hooks, module-extension state
guidance, tooling-readiness guidance, and runbook-authoring guidance.

Suggested inventory table:

```text
Surface:
Current routing role:
Overlaps with:
Disposition:
Reason:
Implementation note:
```

### D-010 - Add Fresh-Agent Routing Tests

Status: recommended

Decision: implementation should define fresh-agent routing tests for the new standard.

Why it matters: routing rules are only useful if a fresh agent can discover and apply them without
knowing the history of the discussion.

Recommended default: run fresh-agent routing tests after implementing the standard; after
material changes to managed root routing, structure-governance rules, runbook-authoring rules,
capability advertisements, route registries, or route-card shape; and when a real routing failure
shows that existing guidance is ambiguous. Do not require the test after every ordinary local edit.

Fresh-agent test prompt shape:

```text
Given this repository and a plain-language request, can a fresh agent identify the likely owner,
candidate capability advertisements, applicable route or ambiguity question, and avoid choosing an
execution surface prematurely?
```

### D-011 - Maintain Capability Advertisements With Route Changes

Status: recommended

Decision: capability advertisements should be reviewed when their owner adds, removes, renames,
or materially changes route registries, route cards, scope families, intent aliases, or execution
ownership.

Why it matters: stale advertisements recreate the original problem: the upper router cannot
discover what a deeper owner can handle.

Recommended default: index-maintenance or route-registry maintenance guidance should require an
advertisement review when discoverable capability families change.

### D-012 - Use A Candidate Authority Hierarchy

Status: recommended

Decision: the first implementation should start from the candidate authority hierarchy in this
discovery and refine only where the existing-surface inventory shows a concrete conflict.

Why it matters: authority hierarchy is central to routing. If implementation leaves it implicit,
different agents or module owners may resolve the same request differently.

Recommended default: include the hierarchy in the managed agent-interface routing reference and
cross-link structure-governance and runbook-authoring surfaces to it rather than restating it in
full.

### D-013 - Make Existing-Surface Inventory A Dedicated Implementation Epic

Status: recommended

Decision: the existing Operating Kit routing-surface inventory should be its own implementation
epic, with a detailed inventory artifact as the epic output.

Why it matters: the inventory is substantial enough to affect sequencing and acceptance. Treating
it as a small preflight task would make it too easy to add new doctrine before understanding and
consolidating existing partial rules.

Recommended default: the implementation plan should begin with an inventory epic before writing
the new managed routing reference. The epic output should use this detailed shape:

```text
Surface:
Current routing role:
Instruction type:
Authority level:
Overlaps with:
Conflict risk:
Disposition:
Reason:
Implementation note:
```

### D-014 - Add Fresh Low-Context Routing Probes To Routing-Bearing Epics

Status: recommended

Decision: implementation-planning guidance should require a fresh low-context routing probe for
every epic that creates or materially changes routing-bearing runbooks or routing surfaces.

Routing-bearing surfaces include:

- managed root routing;
- agent-interface routing references;
- structure-governance routing or placement rules;
- runbook-authoring rules that affect route discovery;
- capability advertisements;
- route registries;
- route cards;
- module or extension routing state rules;
- durable runbooks that select routes, execution surfaces, or owner boundaries.

Why it matters: the failure mode being solved is not whether an author understands the route. It
is whether a fresh agent with low historical context can discover and follow the written routing
structure from a vague user-style intent.

Recommended default: this should become a planning and review hook, not a standalone validator in
V1. The implementation-planning runbook should tell planners to add the probe as validation when
an epic changes routing-bearing surfaces. The execution and planning-review runbooks should check
that the probe was run or explicitly marked not applicable.

Probe shape:

```text
Select the deepest nested realistic routing scenario affected by the epic.
Spawn a fresh low-context subagent.
Give it a vague user-style request for that nested intent.
Pass if it identifies the likely owner, discovers the capability advertisement or route registry,
selects the route or asks the required ambiguity question, and avoids choosing an execution
surface prematurely.
```

### D-015 - Put A Compact Routing Hierarchy In Root `AGENTS.md`

Status: recommended

Decision: the managed root `AGENTS.md` block should expose a compact routing hierarchy and link to
the full routing doctrine. The full doctrine should live in a managed agent-interface reference,
not in root `AGENTS.md`.

Why it matters: agents need the routing reflex at the beginning of a session. Root `AGENTS.md` is
read early, but it should remain concise and stable. Detailed doctrine is more likely to grow as
route cards, capability advertisements, validation probes, and examples mature.

Recommended default: root `AGENTS.md` should contain a short hierarchy and the path to the full
reference. The detailed reference should own examples, field definitions, trigger categories,
inventory expectations, and validation guidance.

## Requirements And Evaluation Criteria

### FR-001 - Generic Dispatch Sequence

The first implementation should define a reusable dispatch sequence from intent through evidence
or blocker recording.

### FR-002 - Authority Hierarchy

The first implementation should define how agents choose among explicit user instruction,
repository instructions, Operating Kit doctrine, module or extension routes, repo-owned docs,
official external docs, visible execution-surface documentation, and model knowledge.

### FR-003 - Capability Advertisement Guidance

The first implementation should define when capability advertisements are needed, where they
belong, and what minimum information they should expose.

### FR-004 - Route Registry And Route Card Guidance

The first implementation should define route registries and route cards as domain-owned dispatch
contracts.

### FR-005 - Placement Guidance

The first implementation should update structure-governance guidance so routers, capability
advertisements, route registries, route cards, state, runbooks, and references have clear owners.

### FR-006 - Runbook Authoring Hook

The first implementation should update runbook-authoring guidance so durable agent-facing or
hybrid runbooks expose their routing contract or point to a route card when the runbook handles
repeated operations.

### FR-007 - Ambiguity Handling

The first implementation should define a standard stop-and-ask behavior for ambiguous provider,
domain, scope, or target matches.

### FR-008 - Anti-Catalog Boundary

The first implementation should explicitly avoid turning root routers or Operating Kit docs into
a centralized catalog of every domain route or execution surface.

### FR-009 - Existing Surface Inventory

The first implementation plan should include an inventory of existing Operating Kit routing and
instruction surfaces, with a disposition for each surface before any new doctrine is added.

### FR-010 - Consolidation Strategy

The first implementation should avoid duplicate authority by choosing one owner for each generic
routing rule and changing other surfaces to concise pointers or local exceptions.

### FR-011 - Conflict Handling

The first implementation should define the simple authority-conflict rule and stop behavior for
conflicts that affect external state, sensitive reads, writes, deletions, permissions, approvals,
or target identity.

### FR-012 - Task Category Triggers

The first implementation should define mandatory, conditional, and lightweight routing categories
so agents know when route lookup is required and when ordinary local work can stay low ceremony.

### FR-013 - Fresh-Agent Validation

The first implementation should include fresh-agent routing scenarios that prove the standard can
be discovered and followed without historical context.

Candidate scenarios:

- deep capability exists below a concise root router;
- provider or domain ambiguity requires a targeted question;
- a new durable documentation path routes through structure governance;
- a local tooling blocker routes through tooling readiness;
- committed module state routes the agent, but live preflight would still decide before action;
- a tiny local code edit stays lightweight while still respecting local instructions.

For epics that create or materially change routing-bearing runbooks or routing surfaces, the
implementation plan should include a fresh low-context routing probe as an epic validation task.

### FR-014 - Advertisement Maintenance

The first implementation should define when capability advertisements must be reviewed during
route-registry, route-card, scope, alias, or owner changes.

### FR-015 - Root Compact Hierarchy

The first implementation should add a compact routing hierarchy to the managed root `AGENTS.md`
block and route detailed behavior to the managed agent-interface routing reference.

### FR-016 - Future Implementation Standard-Adoption Hook

The first implementation should update planning guidance so future implementation work that
creates or materially changes routing-bearing runbooks, route surfaces, capability
advertisements, route registries, route cards, module-routing state, or durable operational
workflows implements routing according to the established Operating Kit routing standard.

Discovery guidance should lightly flag routing-bearing scope when a discovery is expected to lead
to such work. Implementation-planning and planning-review guidance should carry the stronger
requirement by making the standard an explicit dependency and validation target for
routing-bearing epics.

## Non-Functional Requirements

- Public-core-safe: examples must be generic or sanitized.
- Low ceremony for simple work: route cards should be used for repeated, external, sensitive, or
  operational work, not every tiny local edit.
- Scalable ownership: Operating Kit defines the generic pattern; owners instantiate it.
- Drift-resistant: upper layers advertise capability families and point to owners instead of
  duplicating deep route details.
- Reviewable: route cards should make authority, scope, preconditions, stop conditions, and
  evidence visible.

## Risks

| Risk | Likelihood | Impact | Mitigation |
| --- | --- | --- | --- |
| Upper routers become stale capability catalogs. | Medium | High | Use compact advertisements and route-registry links, not full route duplication. |
| Route cards become too heavy for ordinary work. | Medium | Medium | Apply route cards to repeated, sensitive, external, or operational workflows first. |
| Agents still choose visible tools before routing. | Medium | High | Add the route-before-surface rule to managed agent-interface guidance and runbook standards. |
| Capability advertisements are too vague to help. | Medium | Medium | Require capability families, intent aliases, scope families, registry path, and ambiguity rule. |
| Domain owners duplicate Operating Kit doctrine locally. | Medium | Medium | Structure governance should prefer local wrappers and pointers to managed doctrine. |

## Open Questions

### OQ-004 - Managed Root `AGENTS.md` Wording

Owner: kit maintainer

BLOCKER: no, but required before release

What is the shortest safe root wording that tells agents to route before selecting execution
surfaces without bloating the managed block? Current recommendation: include a compact hierarchy
and a pointer to the full routing reference.

### OQ-005 - Follow-Up Recipe Maturity Discovery

Owner: kit maintainer and module owners

BLOCKER: no for routing implementation

Should Operating Kit later define a separate operational recipe maturity model? Current
recommendation: yes, but keep it out of this routing discovery except as an adjacent concern.

## Candidate Implementation Surfaces

Likely first implementation surfaces:

- `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`
- `components/agent-interface/managed/reference/runbook-authoring-standard.md`
- `components/agent-interface/managed/README.md`
- `components/planning-workflows/managed/runbooks/discovery-workflow.md`
- `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
- `components/planning-workflows/managed/runbooks/review-planning-document.md`
- `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
- `components/structure-governance/managed/reference/documentation-structure.md`
- `components/structure-governance/managed/reference/module-extension-state.md`
- `components/structure-governance/managed/README.md`
- `templates/agents/AGENTS.managed-block.md`
- packaged mirrors under `src/codeheart_operating_kit/resources/`
- release notes and consumer-impact record for an instruction-only change

Candidate implementation sequence:

1. Inventory existing Operating Kit routing and instruction surfaces.
2. Create the managed operation-routing-and-dispatch reference.
3. Update structure-governance placement and ownership guidance.
4. Update runbook-authoring, implementation-planning, execution, and review hooks, including
   fresh low-context routing probes for routing-bearing epics.
5. Update root managed block, README routes, packaged mirrors, release notes, consumer-impact
   record, and validation.

Possible later surfaces:

- route-card examples;
- validation or linting for route cards;
- module or extension manifest guidance;
- operational recipe maturity discovery.

## Implementation Capability Scope - Operation Routing And Dispatch Standard

Capability:
Agents can route repeated, structural, external, sensitive, module, product, or ambiguous work
through a shared Operating Kit dispatch model before choosing execution surfaces. Operating Kit
managed surfaces expose the compact routing hierarchy, detailed routing doctrine, route-card
standard, capability-advertisement guidance, inventory/consolidation expectations, and validation
hooks needed for future modules and runbooks to adopt the pattern.

Primary workflow:
A future agent starts from root `AGENTS.md`, sees the compact routing hierarchy, opens the managed
agent-interface routing reference when the task category requires it, identifies the likely owner
and capability advertisement, follows the route registry or route card, checks preconditions, and
then enters the canonical recipe or stops with the required ambiguity question or blocker.

Must cover:
- compact root `AGENTS.md` routing hierarchy and pointer to full doctrine;
- managed agent-interface routing reference with dispatch sequence, trigger categories,
  candidate authority hierarchy, conflict handling, capability advertisements, route registries,
  route cards, and fresh low-context routing probes;
- structure-governance placement guidance for routers, capability advertisements, route
  registries, route cards, state, references, and runbooks;
- runbook-authoring and planning-workflow hooks for routing-bearing runbooks and route surfaces;
- future implementation-plan guidance requiring routing-bearing implementation work to use the
  established routing standard instead of inventing local routing rules;
- a lightweight discovery-workflow hook to flag routing-bearing scope before implementation
  planning, without turning discovery into an implementation checklist;
- detailed existing routing-surface inventory with dispositions before adding new managed
  doctrine;
- source-to-packaged resource mirroring, release-note impact classification, and validation for
  an instruction-only Operating Kit change.

Explicitly out of scope:
- domain-specific route cards for Microsoft 365, GitHub, finance, CRM, or product modules;
- module-specific schemas, consent models, live service preflight rules, and recipes;
- operational recipe maturity, phase-marker, script-promotion, or automation-maturity standards;
- machine-readable route-card schemas, route validators, or manifest requirements in V1;
- automatic consumer migration or scaffolding of new consumer-owned route/state files.

Deferred or blocked:
- operational recipe maturity model: deferred to a separate discovery after routing doctrine is
  accepted or stable enough to depend on;
- machine-readable route-card validation: deferred until repeated Markdown route-card usage proves
  a stable shape;
- module or extension manifest guidance: deferred until capability advertisements in README or
  route-registry introductions prove insufficient.

Preserve decisions:
- D-001 through D-015.

Planner must not reinvent:
- Agent Interface owns routing behavior and route-card standard; Structure Governance owns
  placement and ownership boundaries;
- root `AGENTS.md` gets a compact hierarchy and pointer, not the full doctrine;
- capability advertisements require owner, domain, capability families, intent aliases, route
  registry, and ambiguity rule in V1;
- route cards use Markdown in V1, allow `not applicable` with a reason, and do not require a
  schema;
- existing routing-surface inventory is a dedicated first implementation epic;
- fresh low-context routing probes are required validation for routing-bearing epics;
- future routing-bearing implementation plans must cite and apply the established routing
  standard as a dependency;
- discovery should identify routing-bearing scope early, but implementation planning owns the
  detailed standard-adoption checks;
- recipe maturity remains out of this implementation scope.

Feature-level success evidence:
- a fresh agent can use the installed managed root block and routing reference to avoid choosing
  an execution surface prematurely in routing-bearing scenarios;
- implementation inventory shows existing routing surfaces and disposition decisions;
- managed source files and packaged resource mirrors agree;
- focused validation covers public-core hygiene, Markdown headers, packaged resource parity,
  route visibility, and fresh low-context probe evidence for changed routing-bearing surfaces.

## Discovery Status

Current status: implementation-handoff-ready.

The prior blocker questions about primary managed owner, route-card field set, and capability
advertisement placement are resolved as recommended defaults in the decision inventory. The user
approved the implementation capability scope on 2026-06-25. Implementation planning may proceed
from the approved capability scope and preserved decisions in this discovery.

## Revision Notes

- 2026-06-25: Added existing-routing-surface inventory and consolidation requirements so
  implementation does not layer new doctrine over unresolved older routing guidance.
- 2026-06-25: Added routing trigger categories, simple authority-conflict handling,
  fresh-agent validation expectations, and capability-advertisement maintenance requirements.
- 2026-06-25: Clarified that the candidate authority hierarchy is an Operating Kit routing
  heuristic, not a replacement for native Codex instruction priority.
- 2026-06-25: Added compact root `AGENTS.md` hierarchy requirement, route-card
  `not applicable` allowance, and candidate implementation sequence.
- 2026-06-25: Moved `Intent aliases` into the minimum V1 capability advertisement field set.
- 2026-06-25: Resolved primary owner, route-card field set, and capability-advertisement
  placement as accepted defaults and marked the discovery implementation-handoff-ready.
- 2026-06-25: Added formal implementation capability scope and corrected status back to
  manual-review-ready until the capability scope is approved or revised.
- 2026-06-25: Added future implementation standard-adoption scope so routing-bearing
  implementation work must apply the established routing standard, with a light discovery hook
  and stronger implementation-planning/review hooks.
- 2026-06-25: Recorded user approval of the implementation capability scope and marked the
  discovery implementation-handoff-ready.
- 2026-06-25: Added candidate authority hierarchy, made existing-surface inventory a dedicated
  implementation epic, and defined fresh low-context routing probes as a planning/review hook for
  routing-bearing epics.
- 2026-06-25: Created first discovery draft for generic operation routing and dispatch doctrine.
