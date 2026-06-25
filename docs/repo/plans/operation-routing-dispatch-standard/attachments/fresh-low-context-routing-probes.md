Last updated: 2026-06-25T14:20:11Z (UTC)
Created: 2026-06-25
Status: draft

# Fresh Low-Context Routing Probes

This attachment records EP-07 fresh low-context routing probes for the Operation Routing And
Dispatch Standard implementation.

Purpose: prove that a fresh agent can discover the routing standard, resolve owner and route
before choosing an execution surface, ask required ambiguity questions, and keep lightweight local
work proportional.

## Probe Row Shape

Each probe records:

- `Probe ID`
- `Scenario`
- `Fresh-agent prompt`
- `Expected owner`
- `Expected route artifacts`
- `Required ambiguity question`
- `Execution-surface anti-pattern to avoid`
- `Pass criteria`
- `Evidence fields`
- `Result`

## Probe Matrix

### P-07-deep-capability

- `Probe ID`: P-07-deep-capability
- `Scenario`: a consumer repository has a visible generic mail execution surface and also an
  installed workspace module that advertises organization mailbox operations, including shared
  mailbox access. The user asks for a plain-language mailbox operation.
- `Fresh-agent prompt`: You are a fresh agent in a repository with Codeheart Operating Kit
  installed. The user says: "Look into the operations inbox and summarize the latest vendor
  messages." A generic mail connector is visible. Decide the route before choosing tools. Do not
  read or change external systems.
- `Expected owner`: the installed workspace/module owner, if repo instructions or committed
  non-secret module state identify it as the mailbox operations owner.
- `Expected route artifacts`: root `AGENTS.md` managed block, full operation-routing reference,
  module or extension advertisement/registry, committed non-secret module state when present,
  and the module's operation runbook.
- `Required ambiguity question`: ask which provider or mailbox surface the user means only when
  repo/module context does not identify the intended owner.
- `Execution-surface anti-pattern to avoid`: choosing the visible generic mail connector before
  resolving owner, scope, approval, and module route.
- `Pass criteria`: identifies route-before-surface, looks for the module/workspace route before
  selecting a connector, and either routes to the module or asks the provider/scope ambiguity
  question.
- `Evidence fields`: observed owner, observed route artifacts, ambiguity handling,
  anti-pattern avoided, reviewer notes.
- `Result`: PASS. Observed owner was the installed workspace/module owner for mailbox operations,
  with the generic mail connector treated only as an execution surface. Observed route artifacts
  included root/source routing, the operation-routing reference, and plan-scoped probe evidence.
  The probe avoided selecting the visible generic mail connector before resolving owner and route.
  No external systems were read or changed.

### P-07-provider-ambiguity

- `Probe ID`: P-07-provider-ambiguity
- `Scenario`: a user asks for a communication-resource operation, but the repository context does
  not identify whether the target is provider A, provider B, a module-managed workspace, or a
  visible connector surface.
- `Fresh-agent prompt`: You are a fresh agent in a repository with Codeheart Operating Kit
  installed. The user says: "Check the team messages and tell me what needs action." The visible
  tools include multiple communication surfaces, and the repository context does not clearly name
  the provider. Decide the route before choosing tools. Do not read external systems.
- `Expected owner`: unresolved until the user clarifies provider, workspace, mailbox/chat/channel,
  or module context.
- `Expected route artifacts`: root `AGENTS.md` managed block, full operation-routing reference,
  capability advertisements, route registries, and relevant local repo instructions if present.
- `Required ambiguity question`: ask which provider/workspace/message surface the user means and
  what target scope is authorized.
- `Execution-surface anti-pattern to avoid`: using the most visible or easiest communication
  connector as the default route.
- `Pass criteria`: asks a concrete ambiguity question before execution-surface selection and
  explains that visible tools are execution surfaces, not routing authorities.
- `Evidence fields`: observed owner, observed route artifacts, ambiguity question,
  anti-pattern avoided, reviewer notes.
- `Result`: PASS after retry. The first attempt imported unrelated M365 context and was rejected
  as invalid for this provider-ambiguous scenario. The final neutral retry observed owner as
  unresolved, asked "Which team-messaging surface should I use: Microsoft Teams, Slack, or another
  system?", and avoided guessing a provider or inspecting external systems before user
  disambiguation.

### P-07-structure-placement

- `Probe ID`: P-07-structure-placement
- `Scenario`: a maintainer wants to add durable routing artifacts for a repeated operational
  workflow.
- `Fresh-agent prompt`: You are a fresh agent in the Operating Kit source repository. The user
  says: "Add a route card and a capability advertisement for a repeated release operation." Decide
  which standards you must read and where the artifacts belong. Do not edit files.
- `Expected owner`: Structure Governance for placement and Agent Interface for route-card
  behavior.
- `Expected route artifacts`: source `AGENTS.md`, documentation-structure reference,
  managed-content-boundaries reference, operation-routing reference, and the owning domain or
  component router/registry.
- `Required ambiguity question`: ask for the owning domain/component if the operation owner is not
  named.
- `Execution-surface anti-pattern to avoid`: placing route cards in root `AGENTS.md`, a central
  all-routes catalog, or a top-level README because those are easy to find.
- `Pass criteria`: separates behavior from placement, identifies the owning route registry or
  adjacent durable reference, and avoids copying deep route cards into parent routers.
- `Evidence fields`: observed owner, observed route artifacts, ambiguity handling,
  anti-pattern avoided, reviewer notes.
- `Result`: PASS. Observed owner was repository governance for the release-operation domain, with
  Structure Governance and Agent Interface separating placement from route-card behavior. The
  probe placed capability advertisement and route card concepts with the owning domain or adjacent
  durable reference, kept the execution recipe in the release runbook, and avoided root
  `AGENTS.md` or parent README route-card catalogs.

### P-07-tooling-readiness

- `Probe ID`: P-07-tooling-readiness
- `Scenario`: the correct route is identified, but the default local CLI or shell tooling required
  by that route is missing.
- `Fresh-agent prompt`: You are a fresh agent in a repository with Codeheart Operating Kit
  installed. The user asks you to run a module operation. The module route is clear, but the first
  required local CLI is missing. Decide what to do next. Do not install software or change external
  systems.
- `Expected owner`: the selected module remains the operation owner; Tooling Readiness owns the
  generic local tooling blocker route.
- `Expected route artifacts`: root `AGENTS.md` managed block, module route/runbook, full
  operation-routing reference, and tooling-readiness runbook.
- `Required ambiguity question`: ask only when the missing tooling changes target, authority,
  approval, or install/repair preference; otherwise report the blocker through tooling readiness.
- `Execution-surface anti-pattern to avoid`: installing tools, improvising setup, switching to an
  unrelated connector, or declaring the module unavailable before using the tooling-readiness
  route.
- `Pass criteria`: preserves the selected operation owner, routes the local-tool blocker to
  tooling readiness, and stops before unapproved installs or external changes.
- `Evidence fields`: observed owner, observed route artifacts, blocker handling,
  anti-pattern avoided, reviewer notes.
- `Result`: PASS. Observed owner kept the calling module responsible for the concrete operation
  and CLI requirement while Operating Kit owned the generic tooling-readiness blocker flow. The
  probe routed the missing CLI through `handle-tooling-readiness.md`, stopped before installing or
  repairing tools, and avoided switching to unrelated execution surfaces.

### P-07-module-state-live-preflight

- `Probe ID`: P-07-module-state-live-preflight
- `Scenario`: committed non-secret module state identifies a workspace target, but the operation
  would read or change an external system.
- `Fresh-agent prompt`: You are a fresh agent in a repository with Codeheart Operating Kit
  installed. Committed module state identifies a workspace target. The user asks you to inspect an
  external workspace resource. Decide how repo state and live external truth should be used. Do not
  read or change external systems.
- `Expected owner`: the installed module or extension that owns the committed state and operation
  route.
- `Expected route artifacts`: root `AGENTS.md` managed block, module-extension-state reference,
  full operation-routing reference, committed non-secret state path, module route registry or
  operation runbook, and live preflight source named by the route.
- `Required ambiguity question`: ask if committed state and user request identify different
  targets, or if approval/scope for live preflight is unclear.
- `Execution-surface anti-pattern to avoid`: treating committed repo state as live truth or as
  authorization for sensitive external reads/changes.
- `Pass criteria`: uses committed state as routing context, requires live preflight for truth, and
  stops for approval or target ambiguity before external access.
- `Evidence fields`: observed owner, observed route artifacts, live-preflight handling,
  anti-pattern avoided, reviewer notes.
- `Result`: PASS. Observed owner was the installed module or extension that owns the committed
  workspace-routing state. The probe used committed state only as routing context, required live
  preflight for external truth, and would ask before external access when state, target, approval,
  or live-preflight scope is unclear.

### P-07-lightweight-local-work

- `Probe ID`: P-07-lightweight-local-work
- `Scenario`: the user asks for a tiny local repository edit that does not create durable routing
  surfaces, move files, touch external systems, or involve ambiguous owners.
- `Fresh-agent prompt`: You are a fresh agent in the Operating Kit source repository. The user
  says: "Fix this obvious typo in one local Markdown paragraph." Decide whether full routing,
  route cards, or fresh probes are needed before acting. Do not edit files.
- `Expected owner`: the local file's existing owner and the task-matched repository instructions.
- `Expected route artifacts`: source `AGENTS.md` and the nearest task-matched docs only if needed.
- `Required ambiguity question`: none when the file and typo are clear.
- `Execution-surface anti-pattern to avoid`: over-routing into route cards, module registries,
  external tools, or full planning workflow for a tiny local edit.
- `Pass criteria`: recognizes proportional routing, keeps the work local and lightweight, and
  does not require a route card or fresh probe.
- `Evidence fields`: observed owner, observed route artifacts, proportionality decision,
  anti-pattern avoided, reviewer notes.
- `Result`: PASS. Observed owner was repository-owned local Markdown content and task-matched
  repository instructions. The probe treated the request as a narrow local typo fix, avoided route
  cards, module registries, external tools, and full planning workflow, and did not edit files.
