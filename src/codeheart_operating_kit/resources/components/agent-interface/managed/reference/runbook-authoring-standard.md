Last updated: 2026-06-25T13:05:46Z (UTC)

# Runbook Authoring Standard

Use this reference when creating or materially changing a durable runbook.

This standard improves how agents use runbooks and how humans experience runbook-guided work. It
does not require a mass retrofit of existing runbooks, and it does not turn every short maintainer
checklist into a large template. Apply the parts that match the runbook audience, blast radius,
repeatability, and ambiguity.

## Audience Classes

Declare one primary audience near the top of every new or materially changed durable runbook.

`human-facing`

The runbook guides a conversation with a user. It must make the user's goal clear, ask for
decisions in a manageable order, provide help for values the user may not know, and keep internal
agent mechanics out of user-facing copy.

`agent-facing`

The runbook gives agents an execution path. It must name source of truth, inputs, preconditions,
tool or command lane, ordered procedure, approval gates, stop conditions, evidence, and
validation. A fresh agent should not need to invent the workflow.

`hybrid`

The runbook has both user-facing conversation and operator execution. It must separate user copy
from operator notes so internal mechanics do not leak into the user's experience.

`maintainer-facing`

The runbook guides maintainers through repository, release, governance, or operational work.
Short maintainer runbooks may stay short when ambiguity and blast radius are low. Higher-risk
maintainer runbooks need explicit authority, evidence, validation, and stop conditions.

## Compact Intention Block

Use this compact block near the top of every new or materially changed durable runbook:

```text
Audience: human-facing | agent-facing | hybrid | maintainer-facing

Intent:
<what this runbook is trying to achieve and what good behavior looks like>

Success:
<observable successful outcome>

Agent judgment boundary:
<what the agent may adapt, and what it must not invent or bypass>

Stop boundary:
<when the agent must stop and ask before continuing>
```

The intention block helps the agent handle edge cases. It does not replace concrete execution
steps for risky, ambiguous, or external-state-changing work.

For human-facing runbooks, the intent should describe the desired conversation behavior: guide the
user step by step, avoid internal mechanics in the first turn, and ask only the next useful
decision.

For agent-facing runbooks, the intent should describe the desired execution behavior: use the
named lane, inspect before writes, record evidence, validate results, and stop before inventing
high-impact actions.

## Human-Facing Runbooks

Human-facing runbooks must include a `User-Facing Flow` section or equivalent.

Required quality bar:

- State the user's plain-language goal and expected outcome before technical detail.
- Provide exact or near-exact wording for critical turns.
- Ask one user-owned decision per turn by default.
- Ask no more than two questions in one user-facing turn.
- Present visible choices when the option set is small.
- Put the recommended default first when a safe default exists.
- Keep internal file paths, local state, logs, and implementation mechanics out of first-turn copy.
- Explain how the user can find values they may not know.
- Ask explicit approval before writes, sign-ins, installs, external changes, destructive actions,
  or sensitive reads.
- Route missing local tooling through `../runbooks/handle-tooling-readiness.md` and offer
  blocker-specific choices instead of broad "install tools" prompts.
- End with a clear result and next step.

### Language Preference

Before asking for setup language in a human-facing or hybrid flow, check whether an agent-visible
`.codeheart/user/preferences.yaml` file exists and contains a readable `language` value.

When a readable `language` value exists, continue the current flow in that language.

When the file is absent, unreadable, or has no `language` value, ask once for language and continue
the current flow in the selected language.

Do not expand this rule into a broad preference system inside a runbook. Broader preference
handling belongs in separate Operating Kit guidance.

## Agent-Facing Runbooks

Agent-facing runbooks must be executable by a fresh agent without inventing the workflow.

Required quality bar:

- State the goal and non-goal.
- Name the source of truth and required read order.
- List required inputs and accepted formats.
- State preconditions and tool readiness checks.
- Name the execution lane, such as CLI, API, portal, document surface, or managed runbook.
- Route missing generic local prerequisites through `../runbooks/handle-tooling-readiness.md`
  before improvising install or repair guidance.
- Provide an ordered procedure with concrete commands, API calls, document edits, or portal steps
  where applicable.
- State approval gates before external-state-changing, destructive, sensitive, release, or
  security-relevant actions.
- State stop conditions.
- Define evidence and run-record requirements.
- Define validation that proves the requested outcome, not only that commands ran.
- Include rollback, retain, offboarding, or cleanup guidance when relevant.

Use the fresh-agent test: if another agent can only restate what must be true but still has to
invent commands, workflow, evidence, or validation, the runbook is not specific enough.

### Routing-Bearing Runbooks

When an agent-facing or hybrid runbook handles a repeated routing-bearing operation, expose its
routing contract before recipe execution begins.

A runbook is routing-bearing when it selects an owner, route, execution surface, target scope,
approval class, or external/service path for repeated work. In that case, either:

- point to the owning route card or route registry; or
- include a compact routing section that names the intent family, domain, scope, authority source,
  execution surface, preconditions, approval class, stop conditions, and evidence expectation.

Use `operation-routing-and-dispatch.md` for route-before-surface behavior, authority hierarchy,
capability advertisements, route registries, route-card fields, ambiguity handling, and fresh
low-context routing probes.

Keep route selection separate from recipe execution. The route or route card chooses the lane and
preconditions. The runbook recipe performs the work after routing is complete.

## Hybrid Runbooks

Hybrid runbooks must separate user copy from operator-only material. Use these sections unless a
local runbook shape gives the same separation more clearly:

```text
## User-Facing Flow
## Operator Notes
## Execution Path
## Stop Conditions
## Evidence And Validation
```

User-facing flow owns what the user sees. Operator notes own internal state, file paths, local
configuration, tool selection, and implementation cautions. Execution path owns the work. Stop
conditions and evidence must be visible enough for reviewers to confirm the runbook is safe.

## Maintainer-Facing Runbooks

Maintainer-facing runbooks may be concise when the procedure is low risk, familiar, and
unambiguous.

Add more structure when the runbook controls:

- public releases;
- repository governance;
- public-core safety;
- generated or managed content;
- consumer sync behavior;
- external services;
- destructive cleanup;
- security-sensitive changes.

Higher-risk maintainer runbooks should include:

- trigger and authority;
- required inputs;
- required references;
- ordered procedure;
- evidence record;
- validation;
- stop conditions;
- release, migration, rollback, or handoff notes when relevant.

## Tooling Readiness And DRY Architecture

When a durable runbook can encounter missing local tooling, keep the responsibilities separate:

- the runbook names required tools, capability unlocked, readiness checks, and module-owned
  service preflight;
- `../runbooks/handle-tooling-readiness.md` owns the generic missing-tool conversation, approval
  gate, blocker-specific user choices, and return-to-runbook behavior;
- module-owned runbooks and references own concrete module-specific install commands, versions,
  and service caveats;
- structure governance owns where the runbook belongs, not the internal runbook shape.

Do not duplicate generic package-manager, runtime, or local-tool setup guidance across multiple
managed Operating Kit runbooks. Use a concise route to the tooling-readiness runbook unless the
implementation plan explicitly creates or changes the shared readiness route itself.

Local environment blockers include missing package managers, shell runtimes, CLIs, PowerShell
runtime, PowerShell modules, PATH discovery, browser automation prerequisites, and document/PDF
tools. Service blockers such as external sign-in, tenant consent, admin roles, mailbox access,
SharePoint permissions, licenses, app readiness, API authorization, or live external preflight
remain module-owned.

## Review Checklist

Use this checklist when reviewing a new or materially changed runbook.

Audience and intention:

- The runbook declares exactly one primary audience class.
- The compact intention block is present and understandable.
- The success condition is observable.
- The judgment boundary says what the agent may adapt.
- The stop boundary says when the agent must ask before continuing.

Human-facing checks:

- The user goal is clear before technical detail.
- Critical user turns have exact or near-exact wording.
- The first turn avoids internal mechanics.
- Questions are paced for a nontechnical user.
- Help text exists for hard-to-find values.
- Approval wording is explicit before writes, sign-ins, installs, external changes, or sensitive
  reads.
- Language preference is reused when a readable local `language` preference exists.
- Missing local tooling routes to the managed tooling-readiness runbook with concrete user
  choices.

Agent-facing checks:

- Source of truth and inputs are clear.
- Preconditions and tool readiness checks are clear.
- The execution lane is named.
- Routing-bearing runbooks expose their routing contract or point to the owning route card or
  route registry.
- Route selection is separate from recipe execution.
- Missing generic local tools route to the managed tooling-readiness runbook.
- The ordered procedure is concrete enough for a fresh agent.
- Approval gates and stop conditions are explicit.
- Evidence and validation prove the outcome.
- Cleanup, retain, rollback, or offboarding is covered when relevant.

Hybrid checks:

- User-facing copy is separate from operator notes.
- Internal state and local mechanics do not leak into the user's first experience.
- Execution path, stop conditions, evidence, and validation are not hidden inside user copy.

Maintainer-facing checks:

- The runbook is as short as its risk allows.
- Higher-risk maintainer work has authority, evidence, validation, and stop conditions.
- Public-core, release, migration, and consumer-impact boundaries are explicit when applicable.

Scope checks:

- The change does not accidentally retrofit unrelated runbooks.
- Consumer-owned and module-owned runbooks are preserved unless the plan explicitly changes them.
- Examples are public-safe and use placeholders or sanitized patterns.
- Generic package-manager, runtime, and local-tool readiness guidance is centralized instead of
  copied into every module or managed runbook.
