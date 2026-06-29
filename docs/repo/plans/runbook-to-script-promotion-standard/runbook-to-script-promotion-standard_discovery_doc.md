Last updated: 2026-06-29T14:30:44Z (UTC)
Created: 2026-06-29
Status: draft

# Runbook-To-Script Promotion Standard Discovery

## Discovery Status

Input state: new Operating Kit doctrine request after repeated module and repository operations
showed that some runbook steps are too slow, fragile, or operator-dependent when left as manual
copy/paste instructions, while full automation of judgment-heavy runbooks would be unsafe and
premature.

Output target: manual-review-ready draft. This document captures the proposed generic standard
for deciding when a runbook step should become a script, what scaffolding is required from the
first promoted script, when shared helpers or domain folders are justified, and how scripts
preserve runbook approval and safety boundaries.

Implementation planning is intentionally out of scope for this draft. Downstream module adoption,
including Foundry or Microsoft 365 module updates, should wait until this Operating Kit discovery
is accepted and an Operating Kit implementation plan establishes the managed standard.

## User Intention

Codeheart wants a reusable Operating Kit standard for moving from runbook-only operations toward
deterministic scripts without turning every runbook into automation and without creating
duplicated one-off scripts. The standard should help future repositories, modules, and products
decide:

- what triggers script promotion;
- what should remain a runbook-only conversation or judgment flow;
- what scaffolding is required as soon as the first script exists;
- which shared helpers are worth creating early;
- when domain folders or packages are justified;
- how scripts report evidence and blockers consistently;
- how scripts preserve approvals, safety, and human judgment.

The standard should be generic Operating Kit doctrine. Domain-specific examples may motivate the
work, but the rules should not encode Microsoft 365, Foundry, AWS, AI Execution, or any other
single module as the default case.

## Problem Framing

The current Operating Kit already has related doctrine:

- runbook authoring standards for human-facing, agent-facing, hybrid, and maintainer-facing
  runbooks;
- operational recipe maturity doctrine for shaping executable operational recipes;
- tooling readiness doctrine for local blocker handling;
- operation routing and dispatch doctrine for choosing the right execution surface.

Those standards do not yet answer the next architecture question: when a recurring runbook recipe
should become a durable script, how that script should be structured, and how scripts should grow
from one entrypoint into shared helpers, domain folders, or a package without either accumulating
technical debt or over-abstracting early.

Without a standard, repositories tend to drift into one of two failure modes:

- runbooks remain highly manual, making repeated operations slow, inconsistent, and prone to
  copy/paste or quoting mistakes;
- agents create ad hoc scripts that work once, duplicate fragile logic, hide approval boundaries,
  or become untested mini-wrappers around external systems.

The desired model is a maturity ladder, but not a mandatory prose-first sequence:

```text
runbook guidance
structured operational recipe
reusable script asset
script plus shared infrastructure helpers
domain script group
package or CLI surface
```

Each rung should have clear triggers and stop conditions. A fragile or repeated operation may
start at `runbook + reusable script asset` immediately. The standard should not force agents to
write prose runbooks first when deterministic execution and tests are already clearly needed.

## Goals

- Define generic triggers for promoting a runbook step or operational recipe to a script.
- Define cases that should remain runbook-only.
- Clarify that reusable script assets can be the correct initial implementation surface, not only
  a later cleanup step.
- Inventory current Operating Kit doctrine and runbooks that may still point agents toward the
  older inline-block maturity model.
- Define the minimum scaffolding required when a module or repo adds its first promoted script.
- Define early infrastructure-helper rules that prevent avoidable technical debt without creating
  speculative abstractions.
- Define triggers for adding domain folders under a script area.
- Define triggers for promoting scripts into a package, CLI, or versioned command surface.
- Preserve runbooks as the policy, UX, approval, routing, and judgment layer.
- Require scripts to have explicit inputs, structured outputs, stable markers where useful, and
  proportional tests.
- Keep approval gates, sensitive reads, destructive actions, and external-state changes visible
  and user-approved.
- Provide enough doctrine that downstream modules can adopt the standard later without inventing
  separate promotion rules.

## Non-Goals

- Do not implement the standard in this discovery.
- Do not retrofit existing Operating Kit runbooks or scripts yet.
- Do not retrofit Foundry, Microsoft 365, AI Execution, AWS Platform, or consumer repositories in
  this discovery.
- Do not require every runbook to have a script.
- Do not create a broad automation wrapper or API surface for every domain.
- Do not define one universal language for all scripts.
- Do not force domain folders or packages before the structure is justified.
- Do not make scripts responsible for user conversation, approval decisions, consequence
  explanation, ambiguous target selection, or policy judgment.
- Do not require inline code blocks as an intermediate maturity level before a reusable script
  asset.
- Do not preserve long inline implementations in runbooks after a durable script asset exists.
- Do not authorize silent installs, external changes, sensitive reads, writes, sends, deletes, or
  cleanup actions.
- Do not store secrets, tokens, raw private content, local machine dumps, or customer-specific
  evidence in promoted script output by default.

## Public-Core Safety

This is public Operating Kit discovery. Examples and motivations must stay generic. Do not include
private tenant IDs, account identifiers, customer data, credentials, raw logs, machine-specific
absolute paths, business records, or restricted strategy.

Consumer impact classification for a later implementation is likely `instruction-only change`
unless the implementation also adds validators, managed scaffolds, generated paths, CLI behavior,
or sync behavior.

## Current Evidence

| Source | Finding | Discovery implication |
| --- | --- | --- |
| Runbook authoring standards | Runbooks need audience, intent, success, judgment boundaries, stop boundaries, execution paths, and validation. | Scripts should not replace runbook UX or approval logic; they should execute deterministic steps named by runbooks. |
| Operational recipe maturity doctrine | Durable recipes need explicit execution surfaces, expected markers, structured blockers, and maturity awareness. | The initial execution surface should match risk and repeatability; fragile recipes can start as reusable script assets without a prose-only phase. |
| Tooling readiness doctrine | Missing tools should route through local readiness rather than ad hoc setup advice. | Script scaffolding should include local tooling preconditions and blocker classification instead of assuming the runtime exists. |
| Operation routing doctrine | Agents should route before selecting tools, APIs, connectors, browsers, scripts, or runbooks. | Script entrypoints must be discoverable from runbooks and should not be called randomly outside routing. |
| Consumer runtime materialization work | Runtime tools can need repo-local generated state, visible-terminal handoff, and non-editable consumer installs. | Script standards need to distinguish durable script source from generated local runtime state. |
| Repeated module-operation experience | Deterministic checks such as auth status, package validation, environment checks, inventory, and evidence collection often repeat. | These are strong generic promotion candidates when inputs and outputs are stable. |
| Repeated operator errors | Long inline commands, quoting, request construction, and inconsistent error handling can create avoidable failures. | Promotion should favor reusable script assets for fragile repeatable mechanics and avoid durable long inline blocks. |

## Current Doctrine Alignment Inventory

The implementation plan should include a targeted Operating Kit alignment pass. This is not a broad
inventory of every runbook. It should only update current managed doctrine and runbooks that could
contradict this discovery or keep steering agents toward durable inline-code execution.

Source files to inspect and likely update:

| Source file | Current concern | Expected alignment |
| --- | --- | --- |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Defines `L2 | Tested script block`, which can imply that inline implementations are a durable maturity state. | Remove or reframe that level so reusable script assets are the first durable executable surface. |
| `components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | References executable script blocks and the current maturity labels. | Update planning prompts so implementation plans choose prose, structured recipe, reusable script asset, wrapper, or API without requiring an inline-block phase. |
| `components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | References executable script blocks and promoted recipe assets during execution. | Align execution guidance with reusable script assets, runbook callers, tests, and output contracts. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | Mentions executable script blocks when discovery changes durable operational recipes. | Ensure discovery questions ask whether a reusable script asset is needed when mechanics are fragile or repeated. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Delegates recipe-bearing sections to operational recipe maturity. | Add or adjust cross-reference language after the maturity reference is updated. |

Generated or installed copies under `.codeheart/kit/` should not be treated as the source of truth.
The implementation should modify the `components/.../managed/...` sources and let the normal
managed-content/release path update installed copies.

Historical planning documents may mention the old labels. They should not be rewritten unless they
are active implementation instructions, because old plans are evidence rather than current
doctrine.

## Working Definitions

| Term | Draft meaning |
| --- | --- |
| Runbook | Human, agent, hybrid, or maintainer procedure that owns intent, routing, approval gates, user conversation, and judgment. |
| Operational recipe | A concrete repeatable procedure inside a runbook with named inputs, execution surface, markers, expected outputs, and blockers. |
| Promoted script | A repository- or module-owned executable file that performs one deterministic operation or evidence step with explicit inputs and outputs. |
| Script entrypoint | The script file a runbook names directly for agents to call. |
| Runbook code block | A short invocation, example, or temporary discovery note inside a runbook. It is not a durable execution surface. |
| Infrastructure helper | Shared code for generic mechanics such as output envelopes, redaction, path handling, marker emission, result normalization, or error classification. |
| Domain helper | Shared code for a specific domain's repeated logic, such as a product, module, cloud service, document type, or workflow. |
| Domain folder | A script subfolder used to group cohesive scripts with a shared domain, approval boundary, fixtures, or ownership expectation. |
| Package or CLI surface | A versioned, installed, or reusable command surface that wraps multiple scripts or exposes stable APIs beyond one runbook's local entrypoint. |

## Boundary Model

Runbooks and scripts have different jobs.

| Layer | Owns | Must not own |
| --- | --- | --- |
| Runbook | intent, user-facing flow, routing, target selection, approvals, consequences, fallback choice, stop conditions, evidence requirements | low-level repeated mechanics that are safer as tested deterministic code; full inline script implementations that duplicate durable scripts |
| Script | deterministic execution, local validation, evidence collection, request construction, narrow approved action, structured output | user conversation, approval decisions, hidden scope expansion, broad operation routing |
| Helper | repeated infrastructure or domain logic used by scripts | business or policy decisions that belong in runbooks |
| Package or CLI | stable command surface and reusable internal APIs when script count or consumers justify it | premature abstraction before repeated usage proves the surface |

The runbook remains the source of truth for when a script may be called.

## Trigger Model

### Promotion Triggers

A runbook step or operational recipe becomes a script candidate when one or more are true:

- The same deterministic step appears in multiple runbooks or repeated sessions.
- Manual execution is slow, fragile, or prone to copy/paste mistakes.
- The step has stable inputs and expected outputs.
- The step needs consistent evidence, markers, or JSON output.
- The step repeatedly causes operator confusion.
- The step needs reliable failure classification.
- The step is safe to execute after approvals and inputs are already resolved.
- The step benefits from tests, fixtures, dry-run behavior, or idempotency checks.
- The step contains syntax, quoting, encoding, parsing, or request-construction mechanics that
  are easy to get wrong manually.
- A bug in the step would likely recur across operators or repositories if it stays only as prose.

### Non-Promotion Triggers

Keep a step runbook-only when it is mainly:

- user conversation;
- approval decision;
- consequence explanation;
- ambiguous routing;
- target selection;
- policy judgment;
- fallback choice;
- one-off exploration;
- business decision capture;
- interpretation of incomplete or conflicting user intent.

Runbooks may still reference examples for these steps, but they should not be collapsed into a
script that hides judgment from the user or agent.

### Inline Code Block Rule

Inline code blocks are allowed as short invocations, examples, or temporary discovery notes. They
are not a durable execution surface.

Promote directly to a reusable script asset when a block contains one or more of:

- fragile syntax, quoting, encoding, parsing, or request construction;
- non-trivial branching or error handling;
- structured evidence or marker emission;
- repeated local tooling checks;
- external service request construction;
- sensitive read, write, delete, send, permission, compliance, release, or other state-changing
  mechanics.

After promotion, the runbook should call the script and document inputs, approvals, expected
outputs, blockers, and fallback choices. Do not keep a full duplicate implementation in the
runbook.

### First-Script Scaffolding Triggers

The first promoted script in a repository area, product, or module should trigger basic script
scaffolding immediately. Do not wait until many scripts exist.

Default shape for a module or local source area:

```text
<owner-root>/
  scripts/
    README.md
    <script-entrypoint>
  tests/
    scripts/
      <script-test-or-fixture>
```

The exact owner root may be a module, package, product source tree, or repository governance area.
The important rule is that scripts do not appear without a local README, runbook caller, contract,
and proportional validation.

The first-script README should record:

- script entrypoints and their intended runbook callers;
- required local tooling;
- input contract;
- output contract;
- safety and approval boundary;
- known dry-run or read-only mode;
- where tests and fixtures live.

### Infrastructure Helper Triggers

Create infrastructure helpers early when consistent behavior is known to be required across
promoted scripts, even before the second concrete script exists.

Good early infrastructure-helper candidates:

- output envelope or marker formatting;
- JSON serialization shape;
- secret and private-content redaction;
- command result normalization;
- error or blocker classification;
- dry-run flag handling;
- path resolution;
- evidence record writing;
- local tooling preflight wrappers;
- external request URL or payload construction where errors are common.

These helpers reduce avoidable technical debt because the behavior is infrastructure, not
speculative domain abstraction.

Avoid early helpers for speculative domain behavior. Do not create broad `common utilities` or
domain catch-all helpers without a clear repeated responsibility.

### Domain Helper Triggers

Create domain helpers when repeated concrete domain logic appears or is certain from accepted
script contracts.

Good triggers:

- two scripts need the same domain-specific parsing, mapping, or request-building behavior;
- a domain-specific failure class must be consistent across scripts;
- tests would otherwise duplicate meaningful domain setup;
- a defect would need the same domain fix in multiple scripts;
- a domain-specific helper has a narrow name and clear owner.

Bad triggers:

- a single script might someday need friends;
- the taxonomy looks tidy;
- a helper name is broad enough to absorb unrelated behavior;
- the helper encodes approval, routing, or policy decisions that belong in runbooks.

### Domain Folder Triggers

Start flat unless structure improves navigation, review, or ownership. Add domain folders when
one or more are true:

- three or more scripts form a cohesive domain;
- two domains have different approval or safety boundaries;
- scripts have domain-specific fixtures or test setup;
- filenames become long prefix lists that hurt navigation;
- different maintainers or review expectations apply;
- domain-specific helpers exist and should not be confused with module-wide helpers.

Avoid domain folders when:

- there is only one script in the domain;
- the folder mirrors every runbook section;
- the taxonomy is speculative;
- folder depth makes the runbook call path harder to read.

### Package Or CLI Promotion Triggers

Promote scripts into a package, module, or CLI only when one or more are true:

- scripts need stable reusable internal APIs;
- multiple scripts import the same helpers heavily;
- external consumers need a versioned command surface;
- testing plain scripts becomes awkward or unreliable;
- distribution, install, or runtime materialization behavior matters;
- backwards compatibility becomes a real concern;
- the command surface is used outside one local runbook family.

Do not create a package or CLI only to make a first script look polished.

## Script Contract Requirements

Every promoted script should have:

- one primary purpose;
- explicit required inputs;
- explicit optional inputs;
- read-only, dry-run, or write behavior clearly labeled;
- stable success output;
- stable blocker output;
- nonzero failure behavior when appropriate;
- no secrets in normal output;
- documented local tooling prerequisites;
- a runbook caller;
- proportional tests or fixtures.

For external, sensitive, or state-changing scripts, also require:

- approval packet ID or equivalent runbook-provided approval reference when applicable;
- exact target scope inputs;
- no target broadening inside the script;
- no hidden fallback to broader permissions or different targets;
- idempotency or cleanup expectations;
- structured failure classification.

## Output And Evidence Direction

The standard should prefer a small, mechanical output contract. Scripts should not be asked to
make AI-like or human-like judgments about real-world sensitivity, business meaning, or user
intent. They should report what they are designed to do, what happened at runtime, and what they
chose to emit.

The script author defines the contract at implementation time:

- stable `script_id`;
- supported `mode` values;
- expected `data` shape;
- possible `blocker.class` values;
- output safety behavior;
- whether the script is read-only, dry-run, write-capable, or mixed.

The script fills the runtime result when executed:

- actual `status`;
- short `summary`;
- actual `data`;
- actual `blocker`;
- relevant output safety flags.

Minimum common fields should be:

```json
{
  "status": "success|blocked|failed",
  "script_id": "owner.area.operation",
  "mode": "read_only|dry_run|write|mixed",
  "summary": "Short human-readable result.",
  "blocker": null,
  "data": {},
  "output_safety": {
    "raw_external_content_emitted": false,
    "raw_secret_values_emitted": false,
    "raw_provider_response_emitted": false
  }
}
```

`blocker` should be `null` or an object with stable fields such as `class`, `message`, and
`next_route`.

For external, sensitive, or state-changing scripts, recommended additional fields include
`runbook_caller`, `target_summary`, `action_summary`, and `evidence_summary`, but those should not
be mandatory for every small local helper.

`output_safety` describes emitted output behavior, not a universal guarantee that no sensitive
real-world content exists. A script can reliably say that it did not emit raw provider responses,
tokens, secrets, or raw external content. It cannot reliably decide whether every display name,
site title, subject line, or provider error message is semantically sensitive unless the script
controls and filters that output.

Default output should be stdout. Save output to a local evidence file only when the runbook asks
for durable evidence. Do not commit raw local evidence, secrets, tokens, raw private content, or
provider responses by default.

Additional conventions:

- human-readable marker lines for quick agent scanning when useful;
- machine-readable JSON for durable evidence and tests;
- redacted summaries rather than raw sensitive content;
- clear blocker classes rather than opaque stack traces;
- stable field names for status, mode, summary, data, blocker class, output safety, and next
  suggested route.

## Test Expectations

Tests are part of first-script scaffolding, not an optional late cleanup.

The amount of testing should scale with risk:

- local read-only validation script: fixture or smoke test may be enough;
- parser, serializer, redaction, or blocker classifier: unit tests are expected;
- external request construction: offline construction tests are expected;
- external live operation: non-live tests plus explicit live validation gate are expected;
- state-changing script: dry-run or fixture test plus approval-boundary test is expected.

Tests should prove:

- required input validation;
- output contract;
- blocker classification;
- no raw secret values, raw provider responses, or uncontrolled raw external content in normal
  output;
- output safety flags reflect the script's emitted-output behavior;
- dry-run/read-only behavior where declared;
- helper behavior that would otherwise be duplicated.

## Decision Ledger

### D-001 - Operating Kit Should Own The Generic Promotion Doctrine

Status: recommended

Decision: Operating Kit should own the reusable doctrine for runbook-to-script promotion,
script scaffolding, helper creation, domain folder triggers, package/CLI promotion, and script
output expectations.

Rationale: these questions apply across repositories, modules, and products. Domain modules should
own their own scripts and domain helpers, but not their own generic promotion rules.

BLOCKER: no.

### D-002 - Runbooks Remain The Policy And UX Layer

Status: recommended

Decision: promoted scripts should not replace runbooks. Runbooks remain responsible for user
conversation, target selection, approval gates, consequence explanation, fallback selection,
evidence requirements, and when a script may be called.

Rationale: scripts improve determinism for stable mechanics. They should not hide judgment-heavy
decisions or bypass human approval.

BLOCKER: no.

### D-003 - Promote Deterministic Steps, Not Whole Runbooks

Status: recommended

Decision: the standard should promote individual deterministic steps or recipes rather than whole
runbooks by default.

Rationale: most durable runbooks contain both judgment and mechanics. Promoting the whole runbook
would either under-automate the fragile parts or over-automate the judgment parts.

BLOCKER: no.

### D-004 - Require First-Script Scaffolding

Status: recommended

Decision: the first promoted script in a source area should trigger `scripts/README.md`, a
documented script contract, a runbook caller, and proportional tests or fixtures.

Rationale: waiting for multiple scripts before adding quality scaffolding creates orphan scripts
and unclear ownership. Scaffolding should start early, while domain decomposition can stay light.

BLOCKER: no.

### D-004A - Reusable Script Assets Can Be The First Durable Execution Surface

Status: recommended

Decision: when fragile or repeated execution mechanics are already known, a runbook recipe may
start directly with a reusable script asset and proportional tests. An inline code block is not a
required maturity step.

Rationale: forcing agents through a prose-only or inline-block phase preserves the exact failure
mode the standard is meant to prevent: untested copy/paste mechanics, quoting mistakes, request
construction errors, and inconsistent blocker output.

BLOCKER: no.

### D-005 - Allow Early Infrastructure Helpers

Status: recommended

Decision: infrastructure helpers may be created with the first script when the standard requires
consistent behavior such as output envelopes, redaction, result normalization, path resolution,
dry-run handling, or blocker classification.

Rationale: some helper logic is not speculative abstraction. It is cross-script infrastructure and
avoids preventable technical debt.

BLOCKER: no.

### D-006 - Delay Domain Helpers And Folders Until Cohesion Exists

Status: recommended

Decision: domain helpers and domain folders should be added when repeated concrete domain logic,
review boundaries, fixture boundaries, safety boundaries, or navigation pressure justify them.

Rationale: domain structure created too early often mirrors speculative taxonomy rather than real
maintenance needs.

BLOCKER: no.

### D-007 - Package Or CLI Promotion Is A Later Maturity Step

Status: recommended

Decision: scripts should become a package, PowerShell module, Python package, CLI, or similar
versioned command surface only after repeated usage, helper reuse, distribution needs, or API
stability requirements justify it.

Rationale: a package is useful when there is a real command surface to preserve. It is overhead
when only one or two local script entrypoints exist.

BLOCKER: no.

### D-008 - Scripts Need Structured Outputs And Blockers

Status: recommended

Decision: promoted scripts should emit stable, testable success and blocker output. The output
contract should use a small set of required fields plus domain-shaped `data`, not a rigid universal
envelope. Output safety fields should describe mechanical emitted-output behavior.

Rationale: the main value of a promoted script is deterministic behavior. Unstructured output
reintroduces operator interpretation and makes review harder. Overbroad semantic sensitivity
claims would create false confidence and ask scripts to make judgments they cannot reliably make.

BLOCKER: no.

### D-008A - Script Authors Own Output Contracts

Status: recommended

Decision: the implementer who writes a script must define its `script_id`, allowed modes, data
shape, blocker classes, and output safety behavior. Runtime execution fills the actual status,
summary, data, blocker, and safety flags.

Rationale: output contracts need to be stable enough for agents and tests to rely on them, but
they should remain ordinary software contracts. Agents operating the script should not infer
schema shape from noisy stdout or invent evidence summaries after the fact.

BLOCKER: no.

### D-009 - Tests Are Required From The First Script

Status: recommended

Decision: promoted scripts should include proportional tests or fixtures from the start.

Rationale: a script is executable behavior, not just documentation. Even small scripts need proof
of their input/output contract and safety boundary.

BLOCKER: no.

### D-010 - Inline Code Blocks Are Not Durable Assets

Status: recommended

Decision: inline code blocks in runbooks should be limited to short invocations, examples, or
temporary discovery notes. Long inline implementations should be promoted directly to reusable
script assets when the mechanics are fragile, repeated, or evidence-bearing.

Rationale: keeping full script logic inside runbooks creates drift and weakens testability. The
runbook should call the durable script and preserve the human-facing context around it.

BLOCKER: no.

## Open Questions

### OQ-001 - Should There Be One Generic Script Output Envelope?

Question: Should the first implementation define one reusable JSON envelope for all promoted
scripts, or should it define required fields and let modules use domain-shaped JSON?

Recommendation: define required common fields first, not one rigid envelope. Required fields
should cover `status`, `script_id`, `mode`, `summary`, `blocker`, `data`, and `output_safety`.
External, sensitive, or state-changing scripts should add fields such as `runbook_caller`,
`target_summary`, `action_summary`, and `evidence_summary` when useful. A rigid envelope can come
later if repeated adoption proves value.

BLOCKER: no for discovery draft. The current recommended default is sufficient for first
implementation planning.

### OQ-002 - Where Should Generic Helper Templates Live?

Question: Should Operating Kit implementation include example helper templates, or only doctrine?

Recommendation: start with doctrine and minimal examples. Add templates only after at least one
adopter proves the helper shapes are reusable across domains.

BLOCKER: no.

### OQ-003 - Should The Standard Add Validators Immediately?

Question: Should the first implementation include validators for script README presence, test
presence, or runbook caller references?

Recommendation: not in the first implementation unless the implementation scope is explicitly
validator-bearing. Start with managed doctrine and planning/review hooks. Add validators after
the structure proves stable.

BLOCKER: no.

## Recommended First Implementation Scope

This draft is not yet an implementation handoff, but the likely first implementation should be
instruction-only and managed-doctrine focused:

- add an Operating Kit reference for script promotion and architecture;
- add a runbook for promoting a runbook step or operational recipe to a script;
- update the current managed Operating Kit doctrine and runbooks listed in the alignment
  inventory so they do not contradict the reusable-script-asset model;
- update runbook-authoring and planning workflow guidance to mention script promotion when
  implementation plans create executable script assets;
- add review criteria for orphan scripts, missing tests, hidden approval behavior, and missing
  output contracts;
- avoid validators, templates, CLI behavior, and consumer scaffold changes in the first release
  unless separately approved.

## Downstream Adoption Guidance

After this discovery is accepted and the Operating Kit standard is implemented, downstream
repositories and modules can run their own adoption discovery or implementation plans.

Adoption should start by inventorying existing runbook recipes and scripts, then classifying each
candidate as:

- keep as runbook-only;
- improve as structured runbook recipe;
- promote to script;
- replace long inline implementation with script invocation;
- add infrastructure helper;
- add domain helper or folder;
- defer to package/CLI maturity.

The adoption work should be owned by the repository or module that owns the domain behavior and
validation evidence.

## Readiness For Next Step

This discovery is manual-review-ready as a draft: the main doctrine direction, trigger model,
scaffolding model, and helper/folder/package promotion rules are recorded. It is not
implementation-handoff-ready until the user reviews the recommendations and confirms whether the
first implementation should stay instruction-only or include templates, validators, or planning
workflow hooks.
