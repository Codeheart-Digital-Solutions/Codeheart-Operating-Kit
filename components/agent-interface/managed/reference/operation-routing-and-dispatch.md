Last updated: 2026-06-25T13:05:46Z (UTC)

# Operation Routing And Dispatch

Use this reference when a request may need routing before execution. It defines the generic
Operating Kit routing model. Domains, modules, products, packages, and repository areas instantiate
the model with their own capability advertisements, route registries, route cards, runbooks, and
recipes.

This reference does not define every domain route. It defines how an agent should find the right
owner before choosing a tool, connector, API, browser surface, script, runbook, or manual path.

## Core Rule

Route through the highest-authority applicable operating source before choosing an execution
surface.

Visible tools, connectors, MCP servers, browser paths, local commands, APIs, CLIs, and scripts are
execution-surface options. They are not automatically the routing authority. Choose them after the
intent, domain, authority, scope, and capability route are resolved, unless the user explicitly
requests that specific surface and the request is safe.

## Dispatch Sequence

Use this sequence for routing-bearing work:

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

Meanings:

- `Intent`: what the user is trying to accomplish in plain language.
- `Domain`: the operational area that likely owns the request.
- `Authority source`: the instruction, owner, route, state source, official truth, or user
  decision that is allowed to define the path.
- `Scope`: the target boundary, such as local repo, product area, module, workspace, account,
  document set, mailbox, service tenant, browser session, or external resource.
- `Capability route`: the owner route, route registry entry, or route card that matches the
  request family.
- `Execution surface`: the concrete surface selected by the route, such as shell, API, connector,
  MCP tool, browser, portal, CLI, script, document surface, or manual user action.
- `Preconditions`: tooling, auth, consent, role, license, target, approval, live state, and safety
  checks required before execution.
- `Canonical recipe`: the runbook, script, API procedure, portal procedure, or manual path entered
  after route selection and precondition checks.
- `Evidence or blocker`: the non-secret result, validation, run record, stop condition, or
  feedback route.

## Routing Trigger Categories

Use the routing model with proportional effort.

### Mandatory Routing Check

Use a route or owner check before acting when the request involves:

- creating, moving, deleting, or reorganizing durable files, folders, docs, route surfaces, or
  structure;
- changing managed Operating Kit content, generated consumer surfaces, or root routing;
- installed modules, extensions, products, packages, or committed module state;
- external systems, accounts, workspaces, permissions, approvals, or live service state;
- sensitive reads, writes, deletes, permission changes, offboarding, diagnostics, or compliance
  surfaces;
- repeated operational workflows that already have or likely need a route card.

### Conditional Routing Check

Use a quick owner and ambiguity check when a request has multiple plausible domains, providers,
tools, products, or execution surfaces. If the owner is clear and the work is low risk, continue
with the nearest applicable runbook or local convention.

### Lightweight Direct Work

Tiny local edits, simple command answers, and ordinary code changes may proceed without formal
route-card lookup when the owner and conventions are already clear. Existing repository
instructions, local code patterns, and relevant validation still apply.

New durable files, folders, route surfaces, or documentation paths are not lightweight just
because the edit is small. Confirm the owner and placement rule first.

## Routing Authority Hierarchy

This hierarchy is an Operating Kit routing heuristic. It does not replace Codex's native
instruction priority, system safety rules, or the active instruction stack. Use it after those are
already in force, when deciding which repository, Operating Kit, module, product, domain, external,
or tool source should route a request.

Prefer the highest applicable routing authority:

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

Lower layers still matter. They may provide required target detail, execution constraints,
preflight results, or current external facts. Live external preflight, approval gates, and current
official documentation may still stop or reshape execution after a route is selected.

## Ambiguity Handling

When more than one domain, provider, route, target, or execution surface plausibly matches:

1. Apply explicit user instruction when safe.
2. Check nearest repository, product, module, and route-registry context.
3. Use capability advertisements to identify candidate owners.
4. If one owner is still clearly nearest and safe, follow that owner.
5. If ambiguity remains material, ask one targeted provider, domain, scope, or target question
   before execution.

Do not resolve ambiguity by choosing the most visible tool or the surface with the easiest
connector. A visible connector may be the correct execution surface, but the route decides that.

## Authority Conflict Handling

When routing authorities conflict:

1. Preserve safety and public-core boundaries first.
2. Follow explicit user instruction when it is safe and within the user's authority.
3. Follow the nearest legitimate owner for the artifact, domain, or operation.
4. Use official external truth for current external-system behavior.
5. Stop before sensitive reads, writes, deletes, permission changes, or external changes when the
   conflict affects authority, target, approval, or live state.

Record the conflict in the current task summary, execution log, blocker, or feedback route when it
is likely to recur.

## Capability Advertisements

A capability advertisement is a compact boundary summary that lets upper routers discover whether
a deeper owner may handle a plain-language intent. It is not a full route registry and not a recipe.

Use advertisements for domains, modules, products, packages, or repository areas that own repeated
operations or deep routes that are not obvious from the parent router.

Minimum V1 fields:

```text
Owner:
Domain:
Capability families:
Intent aliases:
Route registry:
Fallback or ambiguity rule:
```

Optional V1 fields:

```text
Scope families:
State location:
```

Field meanings:

- `Owner`: the owning domain, module, product, package, repo area, or responsible document.
- `Domain`: the operational domain the owner covers.
- `Capability families`: broad groups the owner can route, not every detailed route.
- `Intent aliases`: plain-language phrases users may use for the capability family.
- `Route registry`: where detailed request families or route cards live.
- `Fallback or ambiguity rule`: what to do when another owner or provider may also match.
- `Scope families`: target boundaries the owner can handle.
- `State location`: committed non-secret state location when applicable.

Generic example:

```text
Owner: <workspace module>
Domain: external workspace operations
Capability families:
- users, groups, and access
- workspace files and folders
- named and shared communication resources
- diagnostics and validation
Intent aliases:
- look up workspace access
- inspect shared communication resources
- validate external workspace setup
Route registry: <module route registry>
Scope families:
- workspace
- named resource
- user or group
State location: docs/repo/state/<module-id>/ when present
Fallback or ambiguity rule: ask when more than one provider or domain can satisfy the intent.
```

## Route Registries

A route registry is the owner-domain index of request families. It may be a reference file,
capability map, table, README section, or set of route-card sections. It should help the agent
choose a route without reading every runbook.

A registry should:

- map intent families and aliases to route cards or canonical recipes;
- identify scope families and required targets;
- identify ambiguity rules for common overlaps;
- point to state sources and live truth sources where applicable;
- keep detailed execution steps in canonical recipes unless the route itself is very small.

Parent routers should link to registries or carry capability advertisements. They should not copy
the full route cards from the owning domain.

## Route Cards

A route card is the dispatch contract for one request family. It tells the agent what must be
resolved before execution starts. It does not have to contain the whole recipe.

Use route cards for repeated, sensitive, external, permissioned, cross-domain, or ambiguous
operations. Do not require route cards for every tiny local edit.

V1 route-card fields:

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

Route-card field meanings:

- `Route ID`: stable identifier inside the owner's registry.
- `Intent patterns`: request phrases, aliases, or patterns that match the route.
- `Domain`: operational domain that owns the route.
- `Lifecycle`: setup, operate, validate, offboard, diagnose, release, maintain, or another owner
  lifecycle.
- `Scope`: target boundary the route applies to.
- `Action class`: read, write, delete, permission, diagnostic, compliance, release, or other
  impact class.
- `Authority source`: instruction, owner, route registry, state source, user approval, or external
  source that authorizes this route selection.
- `State sources`: committed non-secret routing context or local state the route may use.
- `Live truth source`: current external or local system that must be checked before execution when
  current state matters.
- `Default execution surface`: preferred surface after routing is complete.
- `Fallback surface`: alternate surface or manual path when the default cannot be used.
- `Preconditions`: tool readiness, auth, role, consent, license, target resolution, approval,
  safety, or validation requirements.
- `Approval class`: whether the route is read-only, sensitive read, write, destructive,
  permission-changing, release, or external-state-changing.
- `Canonical runbook or reference`: the recipe or durable reference to enter after routing.
- `Stop conditions`: conditions that require asking the user, returning to an owner, or recording a
  blocker.
- `Evidence fields`: non-secret evidence, run-record fields, validation, or feedback needed after
  execution.

Fields may be marked `not applicable` with a short reason. Do not delete a field only because it is
awkward; the awkwardness often reveals a missing owner, scope, precondition, or stop condition.

## Recipe Boundary

Routing selects the path. A canonical recipe performs the work.

Keep recipe details in runbooks, scripts, API procedures, portal procedures, or manual procedure
docs. Use the route card to select the recipe, name preconditions, and define stop conditions and
evidence. Avoid copying long recipe steps into every parent router or capability advertisement.

## Live Preflight And Approval Boundary

Capability advertisements, route registries, route cards, and committed routing state help choose
the path. They do not authorize sensitive reads, writes, deletes, permission changes, releases, or
external changes.

Before those actions, run the route's live preflight and approval gate. If committed state
conflicts with live truth, stop and resolve through the owner route or canonical runbook.

## Advertisement Maintenance

Review a capability advertisement when its owner:

- adds, removes, renames, or materially changes a route registry;
- adds, removes, renames, or materially changes route cards;
- changes capability families, intent aliases, scope families, or execution ownership;
- changes state locations or live truth sources used for routing;
- changes ambiguity handling or fallback ownership.

Advertisement review is part of discoverability maintenance. Update the nearest README or route
registry introduction when discoverable capability families change.

## Fresh Low-Context Routing Probes

Use fresh low-context routing probes when an epic creates or materially changes routing-bearing
surfaces. Routing-bearing surfaces include root routing, agent-interface routing references,
structure-governance routing or placement rules, runbook-authoring routing rules, capability
advertisements, route registries, route cards, module state routing rules, and durable runbooks
that select owners, routes, or execution surfaces.

Probe shape:

```text
Select the deepest nested realistic routing scenario affected by the change.
Spawn a fresh low-context agent.
Give it a vague user-style request for that nested intent.
Pass if it identifies the likely owner, discovers the capability advertisement or route registry,
selects the route or asks the required ambiguity question, and avoids choosing an execution
surface prematurely.
```

Record probe evidence in the implementation plan, execution log, review notes, or final summary.
If no probe is required, record why the changed surface is not routing-bearing.

## Quick Examples

Provider ambiguity:

- Request: "Look into the shared messages for this workspace."
- Correct routing behavior: identify candidate providers and workspace routes first, then ask a
  targeted provider or workspace question if context does not decide.
- Incorrect behavior: use the first visible mail connector because it exists.

Durable documentation path:

- Request: "Add a route card for a repeated release operation."
- Correct routing behavior: check Structure Governance for placement and this reference for route
  card behavior, then place the card with the owning route registry.
- Incorrect behavior: put the card in root `AGENTS.md` because agents read it first.

Local tooling blocker:

- Request: "Run the module operation."
- Correct routing behavior: select the module route first; if the route is blocked by missing
  local tooling, use the managed tooling-readiness runbook.
- Incorrect behavior: silently install a tool before the route or approval class is known.
