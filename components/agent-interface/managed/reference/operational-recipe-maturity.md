Last updated: 2026-06-29T14:49:19Z (UTC)

# Operational Recipe Maturity

Use this reference when a runbook, route-selected procedure, runbook code block, reusable script
asset, or promoted asset starts to behave like repeatable operational machinery.

Audience: agent-facing

Intent:
Keep operational recipes flexible while giving repeated executable work enough structure,
validation, evidence, and promotion discipline to stay reliable.

Success:
Agents can tell ordinary guidance from recipe-bearing work, apply the smallest safe maturity
shape, validate the recipe proportionally, and record promotion or non-promotion decisions without
inventing a wrapper, command, or API too early.

Agent judgment boundary:
The agent may keep low-risk guidance below the recipe threshold, may keep a recipe at the lowest
safe maturity level, and may use owner-specific stricter rules. It must not treat the maturity
levels as a required ladder, bypass approval or routing gates, or freeze module-specific layout
rules in this generic reference.

Stop boundary:
Stop before promoting a recipe into a durable script, command, API, tool, or new folder layout
when no approved owner, placement rule, validation path, output contract, or promotion decision
exists.

## Relationship To Routing

Routing chooses the owner, scope, execution surface, preconditions, approval class, and canonical
recipe. Recipe maturity starts after that route has selected the recipe.

Use `operation-routing-and-dispatch.md` for route-before-surface behavior. Use this reference for
the selected recipe's execution shape, validation, evidence, blocker output, and promotion
boundary.

Route cards should point to canonical recipes. They should not copy long recipe internals.

## Working Terms

`Route card`

Dispatch contract that selects the owner, scope, execution surface, preconditions, approval class,
canonical recipe, stop conditions, and evidence fields.

`Runbook`

Durable operational document that guides user conversation, agent execution, maintainer work, or a
hybrid of those audiences.

`Operational recipe`

Repeatable procedure entered after routing is complete. It performs the work through prose steps,
commands, runbook code blocks, portal steps, APIs, manual actions, reusable script assets,
wrappers, or tool surfaces.

`Recipe asset`

Physical form of a recipe or recipe part, such as a runbook section, short runbook code block,
reusable script asset, test file, fixture, schema, CLI command, wrapper, API procedure, or tool
surface.

`Promotion`

Deliberate movement of a repeated recipe into a more durable, testable, reusable, or productized
form.

`Proto-software`

Recipe content with software-like behavior that still lives in prose or references without enough
lifecycle discipline.

## Threshold Test

Use this question before applying recipe maturity:

```text
If a fresh agent uses this section, is it performing a repeatable operation with inputs,
preconditions, execution, evidence, and validation?
```

If no, the content is below recipe threshold. If yes, it is at least L1. If durable executable
mechanics, markers, structured outputs, blocker classes, fixtures, or tests are needed, review
whether it should start as a reusable script asset.

Do not classify every checklist, explanation, or routing table as a recipe.

## Maturity States

Keep the recipe at the lowest maturity level that is safe, testable, reviewable, and ergonomic for
its actual use. The levels are possible forms, not a required ladder.

| State | Name | Meaning | Typical validation |
| --- | --- | --- | --- |
| Below recipe threshold | Ordinary guidance | Prose, conversation guidance, routing aid, or small checklist that does not itself define a repeatable operation with inputs, execution, evidence, and validation. | Normal runbook or document review. |
| L1 | Structured runbook recipe | Repeatable recipe with inputs, preconditions, approval gates, stop conditions, evidence, and validation. | Fresh-agent executability review. |
| L2 | Reusable script asset | Separate script file invoked by a runbook with explicit inputs, stable output, and tests. This may be the first durable executable surface when mechanics are already fragile, repeated, or evidence-bearing. | Script tests, fixture tests, output-contract checks, and runbook invocation validation. |
| L3 | Thin command wrapper | CLI-style command or wrapper validates inputs, runs the recipe, and emits stable structured output after repeated usage proves the command shape. | Command tests, interface contract tests, and evidence validation. |
| L4 | Mature API or tool surface | Durable tool or API surface exists because usage, safety, auth, observability, or scale justifies productization. | Product-grade tests, auth and safety validation, observability, and release process. |

## Recipe Review Triggers

A recipe review trigger is a signal that ordinary runbook text may have become operational
machinery. The trigger asks whether maturity rules apply; it does not force promotion.

Use a recipe maturity review when work includes:

- repeated workflow;
- external system or live service operation;
- approval-gated action;
- sensitive read, write, delete, permission, compliance, release, or external-state-changing
  action;
- embedded durable executable mechanics;
- expected markers;
- structured summary or blocker output;
- non-live or live validation tests;
- repeated operator or agent mistakes;
- the same logic copied across multiple runbooks, modules, packages, or owner areas.

## Recipe Metadata

L1 recipes need a compact metadata block or equivalent section. Keep this small.

```text
Recipe ID:
Purpose:
Inputs:
Preconditions:
Approval class:
Execution surface:
Evidence output:
Validation:
Stop conditions:
```

Reusable script assets add the execution-contract details needed for tested execution.

```text
Recipe asset level:
Expected markers:
Structured outputs:
Error or blocker taxonomy:
Non-live tests:
Phase model: optional unless phase-specific failure localization is needed
```

When promotion is under consideration, record the review result.

```text
Promotion destination:
Promotion criteria:
Non-promotion decision:
```

## Validation Tiers

Use the smallest validation tier that proves the changed recipe behavior.

| Tier | Purpose |
| --- | --- |
| Fresh-agent executability review | Proves another agent can follow the recipe without inventing workflow. |
| Non-live test | Proves helper logic, parsing, output shape, or blocker classification without touching an external system. |
| Dry-run or preflight | Proves local tooling, auth context, target resolution, or external readiness without the final sensitive or write action. |
| Approval-gated live validation | Proves the recipe outcome after the correct approval and route-specific preconditions. |

L1 recipes normally need fresh-agent executability review. Reusable script assets normally need
non-live tests for helper logic, marker output, blocker shape, summary shape, and output-contract
behavior. Live validation only belongs where the route, approval class, and user approval allow
it.

## Evidence And Blockers

For reusable script assets and above, require stable input and output contracts.

Useful evidence includes:

- expected marker lines for command transport or phase progress;
- compact structured summary output;
- compact structured blocker output;
- non-secret evidence fields that can be copied into a run record;
- validation command names and results.

Generic blocker shape:

```text
Blocker ID or class:
Recipe phase or step:
Target:
Non-secret error text:
Evidence command or marker:
Fallback:
User decision needed:
```

Modules and domains own their specific blocker classes. The Operating Kit owns only the generic
shape and evidence boundary.

Do not commit secrets, tokens, token caches, raw sensitive content, raw mailbox content, raw logs,
local session paths, downloaded sensitive artifacts, or private tenant details as recipe evidence.

## Promotion And Non-Promotion

Promotion must be explicit, justified, and reviewable.

Valid promotion triggers include:

- repeated copy/paste of the same block;
- recurring operator or agent execution errors;
- complex input validation that is risky to do by hand;
- multiple modules or runbooks needing the same recipe;
- structured output consumed by later automation;
- safety depending on consistent execution order;
- inline implementation becoming too long or too hard to review;
- live operation requiring stronger test coverage;
- repeated failure by fresh agents or operators to run the recipe correctly.

Valid non-promotion reasons include:

- low frequency of use;
- experimental or rapidly changing workflow;
- official tooling or external API behavior still evolving;
- domain behavior not yet understood well enough to freeze;
- extra abstraction would hide necessary operator judgment;
- script, CLI, or API packaging would add maintenance cost without reducing risk;
- the current structured runbook recipe is clear, safe, and reliable enough.

`Do not promote yet` is a valid recorded outcome of recipe review.

## Placement Principles

Use structure governance for durable placement decisions. Generic rules:

- promoted recipe assets live under the owning route, recipe, package, module, product,
  repository, or source-area boundary;
- temporary execution scripts stay temporary and are not committed;
- reusable scripts, tests, fixtures, schemas, wrappers, and API surfaces need an obvious owner,
  validation path, and discoverability route;
- owner-specific conventions may specialize this standard but should not silently weaken it.

Use `runbook-to-script-promotion-standard.md` for reusable script asset promotion, first-script
scaffolding, output contracts, helper timing, and review flags.

This reference does not define concrete Foundry, module, script, test, fixture, wrapper, or API
folder paths.

## Runbook Shape After Promotion

Promotion should not make the runbook disappear. The runbook remains the operator-facing entry
point and keeps the operational contract:

- purpose and scope;
- required inputs;
- approval and safety gates;
- invocation path for the promoted asset;
- expected markers or structured output;
- evidence and validation requirements;
- stop conditions and fallback path.

Do not keep a stale copy of a full promoted script in the runbook. If a short excerpt is useful,
make it clearly non-authoritative and link to the canonical asset.

## Planning And Review Expectations

Implementation plans that create or materially change recipe-bearing work should name:

- target maturity state;
- validation tier;
- evidence shape;
- promotion destination or non-promotion decision;
- placement boundary for promoted assets.

Reviewers should confirm that the plan does not promote too early, preserve long inline
implementations as durable assets, hide executable behavior inside prose, duplicate full doctrine
across multiple files, or weaken routing, approval, secrets, or external-state boundaries.

## Non-Goals

This generic standard does not:

- rewrite domain recipes;
- rewrite or adopt M365 mailbox-content recipes;
- create promoted recipe assets;
- define Foundry module packaging conventions;
- define concrete script, test, fixture, command, wrapper, or API paths;
- design broad wrapper or API surfaces;
- retrofit every existing runbook;
- replace routing, approval, safety, or tooling-readiness standards.
