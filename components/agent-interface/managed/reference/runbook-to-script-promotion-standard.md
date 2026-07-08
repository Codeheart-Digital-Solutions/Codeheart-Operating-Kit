Last updated: 2026-07-08T14:07:02Z (UTC)

# Runbook-To-Script Promotion Standard

Use this reference when a runbook, operational recipe, or repeated maintainer procedure contains
deterministic mechanics that may need to become a reusable script asset.

Audience: agent-facing

Intent:
Keep runbooks as the policy, UX, routing, approval, and judgment layer while moving fragile,
repeated, evidence-bearing mechanics into tested script assets when that is the safer initial
execution surface.

Success:
Agents can decide whether a step should stay runbook-only, become a structured runbook recipe,
become a primitive or workflow script, extract helpers or folders, or wait for wrapper/API
maturity without creating premature broad command surfaces.

Agent judgment boundary:
The agent may promote deterministic mechanics when triggers are met and an owner, placement,
tests, output contract, and runbook caller are clear. It must not hide user approval, broaden
targets, invent domain policy, commit sensitive output, or create a package, CLI, API, or wrapper
only to make one script look polished.

Stop boundary:
Stop before creating a durable script asset when the owner, placement boundary, approval model,
test path, output contract, or target scope is missing.

## Relationship To Other Standards

Use `operation-routing-and-dispatch.md` before selecting a script, command, API, connector,
browser, portal, or manual execution surface.

Use `operational-recipe-maturity.md` to decide whether content has crossed the operational recipe
threshold and what maturity shape is safe.

Use this reference for the script-promotion decision, first-script scaffolding, helper/folder
growth, script output contract, and review flags.

Use structure-governance references for placement ownership. This reference gives default script
shape; owner-specific modules, products, packages, and repositories may specialize it without
weakening the safety, testing, and output-contract rules.

## Controlled Vocabulary

Preferred terms:

- `runbook`: durable procedure that owns intent, user flow, routing, approvals, consequences,
  fallback choice, stop conditions, evidence expectations, and validation.
- `operational recipe`: repeatable procedure entered after routing is complete, with inputs,
  preconditions, execution, evidence, validation, and blockers.
- `runbook code block`: short invocation, example, or temporary discovery note inside a runbook.
  It is not a durable execution surface.
- `reusable script asset`: repository- or module-owned executable file that performs one
  deterministic operation or evidence step with explicit inputs, outputs, tests, and a runbook
  caller.
- `script asset role`: review label for the job a committed script-area asset performs inside L2,
  such as primitive script, workflow script, or helper. A role is not a new maturity level.
- `primitive script`: reusable script asset that performs one narrow deterministic operation,
  evidence step, transformation, or request-construction task after routing, approvals, and
  prerequisite readiness are already resolved.
- `workflow script`: reusable script asset that deterministically composes stable phases,
  primitive scripts, public script entrypoints, and helpers for one route-selected workflow. It
  owns execution order and phase evidence, not user conversation or broad routing.
- `helper`: script-area asset role for imported infrastructure or domain code shared by scripts.
  A helper may be import-only code rather than an executable entrypoint. It is not a runbook-named
  entrypoint unless it also deliberately exposes a script entrypoint with its own contract.
- `script entrypoint`: script file a runbook names directly for agents or maintainers to call.
- `infrastructure helper`: shared code for generic mechanics such as output formatting,
  redaction, result normalization, path handling, dry-run handling, or blocker classification.
- `domain helper`: shared code for one domain's repeated parsing, mapping, request construction,
  fixtures, or failure classes.
- `domain folder`: script subfolder used when cohesive scripts share a domain, safety boundary,
  fixtures, ownership, or review expectation.
- `thin command wrapper`: command-style surface that validates inputs, runs one or more reusable
  script assets, and emits stable output after repeated usage proves the command shape.
- `mature API/tool surface`: durable productized surface justified by usage, safety, auth,
  observability, scale, or external consumers.
- `local ad hoc script`: temporary, uncommitted script or command used for exploration,
  diagnostics, or unscripted gaps. It is not a durable reusable script asset until deliberately
  promoted with owner, contract, tests, and runbook caller.

Deprecated or avoid:

- Do not use `tested script block` as a maturity state.
- Do not use `executable script block` to mean a durable asset.
- Do not use `command_wrapper` as an L2 script asset role. Use `thin command wrapper` only for the
  L3 maturity state after repeated usage proves a command surface.
- Prefer `reusable script asset` over `promoted script` when the durable asset is specifically a
  script.
- Prefer `reusable script asset` over `promoted recipe asset` when the asset is specifically a
  script. Use `promoted recipe asset` only when speaking generically across scripts, fixtures,
  schemas, wrappers, APIs, or other asset types.

## Boundary Model

Runbooks and scripts have different jobs.

| Layer | Owns | Must not own |
| --- | --- | --- |
| Runbook | intent, user-facing flow, routing, target selection, approvals, consequences, fallback choice, stop conditions, evidence requirements | low-level repeated mechanics that are safer as tested deterministic code; full inline implementations duplicated from scripts |
| Reusable script asset | deterministic execution, local validation, evidence collection, request construction, narrow approved action, structured output | user conversation, approval decisions, hidden scope expansion, broad operation routing |
| Helper | repeated infrastructure or domain logic used by scripts | business or policy decisions that belong in runbooks |
| Package, CLI, or API | stable command surface and reusable internal APIs when script count, consumers, or distribution justify it | premature abstraction before repeated usage proves the surface |

The runbook remains the source of truth for when a reusable script asset may be called.

## Script Asset Roles Inside L2

Primitive scripts, workflow scripts, and helpers are role vocabulary inside the L2 reusable script
asset state. They make review sharper without creating new maturity states. Primitive and
workflow scripts are executable entrypoints. Helpers are script-area assets that support those
entrypoints and may be import-only.

Use `primitive script` when the asset:

- performs one narrow operation, transformation, evidence step, or request-construction task;
- has explicit inputs and preconditions;
- assumes route selection, approvals, and prerequisite readiness were handled before invocation;
- emits stable structured output or a stable blocker.

Use `workflow script` when the asset:

- runs a deterministic multi-phase process after one route has already been selected;
- composes primitive scripts, public script entrypoints, and helpers through documented contracts;
- performs prerequisite readiness checks that are part of the normal deterministic process;
- emits phase-level status, evidence, and blockers;
- reduces repeated agent thinking without hiding approval or routing decisions.

Use `helper` when the asset:

- is imported by scripts rather than named directly by runbooks;
- contains repeated infrastructure or domain logic with a stable local contract;
- is tested in proportion to the logic it owns;
- stays at the narrowest durable owner boundary until real cross-boundary reuse exists.

A local ad hoc script is allowed for exploration and for gaps that do not yet have a durable
script. Keep it local and uncommitted unless promotion triggers are met. When the ad hoc script is
promoted, give it a role, owner, input contract, output contract, tests or fixtures, and runbook
caller.

A workflow script must not own:

- user conversation;
- approval decisions;
- broad route selection;
- ambiguous target selection;
- policy or business judgment;
- hidden fallback to broader permissions or different targets;
- hidden scope expansion or target broadening.

## Workflow Composition And Dependencies

Workflow scripts may depend on other reusable script assets when the dependency is contract-based
and reviewable.

Allowed dependencies include:

- primitive scripts with stable input and output contracts;
- public script entrypoints named by the owning script area;
- infrastructure or domain helpers imported at the narrowest durable owner boundary;
- fixture, schema, or output-contract files owned by the same route, module, package, product, or
  repository boundary.

The workflow script owns phase order, dependency invocation, local retry behavior when safe, and
phase evidence. The called primitive or helper owns its own narrow mechanics. Do not duplicate
primitive internals inside the workflow when a stable primitive already exists.

Document workflow dependencies in the script contract or `scripts/README.md` with enough detail
for a reviewer to see:

- which entrypoints or helpers are called;
- which inputs are passed between phases;
- which phase can return each blocker class;
- whether the workflow is read-only, dry-run, write-capable, or mixed;
- where tests or fixtures prove the composition.

## Prerequisite Readiness Plus Operation Primitive Pattern

When a normal operation always needs the same deterministic readiness or access preflight before a
narrow operation, prefer this shape:

- one workflow script performs the route-selected readiness phases and then calls the operation
  primitive;
- the primitive script remains narrow and declares readiness or access as a precondition;
- the runbook calls the workflow for the normal path and calls the primitive directly only when
  the runbook explicitly states the precondition is already satisfied;
- structured output distinguishes readiness blockers from operation blockers.

Fields such as `readiness_required` or `access_required` may appear in a domain-owned contract
only when the owning domain defines their exact meaning. This generic standard owns the pattern,
not provider-specific readiness semantics.

## Promotion Triggers

A runbook step or operational recipe becomes a reusable script asset candidate when one or more
are true:

- The same deterministic step appears in multiple runbooks or repeated sessions.
- Manual execution is slow, fragile, or prone to copy/paste mistakes.
- The step has stable inputs and expected outputs.
- The step needs consistent evidence, markers, or JSON output.
- The step repeatedly causes operator or agent confusion.
- The step needs reliable failure classification.
- The step is safe to execute after approvals and inputs are already resolved.
- The step benefits from tests, fixtures, dry-run behavior, or idempotency checks.
- The step contains syntax, quoting, encoding, parsing, or request-construction mechanics that
  are easy to get wrong manually.
- A bug in the step would likely recur across operators or repositories if it stays only as prose.

## Non-Promotion Triggers

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

Runbooks may still reference short examples for these steps, but they should not be collapsed
into scripts that hide judgment from the user or agent.

## Inline Code Block Rule

Inline code blocks are allowed as short invocations, examples, or temporary discovery notes. They
are not durable execution surfaces.

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

## First-Script Scaffolding

The first reusable script asset in a repository area, product, package, source area, or module
should trigger basic script scaffolding. Do not wait until many scripts exist.

Default shape:

```text
<owner-root>/
  scripts/
    README.md
    <script-entrypoint>
  tests/
    scripts/
      <script-test>
    fixtures/
      scripts/
        <fixture>
```

The exact owner root may be a module, package, product source tree, repository governance area, or
another locally owned boundary. Start flat. Add subfolders only when they improve navigation,
review, safety, ownership, or fixtures.

The first-script README should record:

- script entrypoints, their script asset role, and their intended runbook callers;
- required local tooling;
- input contract;
- output contract;
- workflow dependencies and phase boundaries when an entrypoint is a workflow script;
- helper placement and importing scripts when helpers are present;
- safety and approval boundary;
- read-only, dry-run, or write behavior;
- where tests and fixtures live.

Keep the README index compact. A small table with entrypoint, role, caller, behavior, and test
path is enough unless the owner has a stronger local convention.

## Helper Rules

Infrastructure helpers may be created with the first script when consistent behavior is known to
be required across promoted scripts.

Good early infrastructure-helper candidates:

- output marker or JSON formatting;
- secret and private-content redaction;
- command result normalization;
- error or blocker classification;
- dry-run flag handling;
- path resolution;
- evidence record writing;
- local tooling preflight wrappers.

Create domain helpers only when repeated concrete domain logic appears or is certain from
accepted script contracts. Avoid broad `utils`, `common`, or manager-style helpers that can absorb
unrelated behavior.

Place helpers at the narrowest durable owner boundary that matches real reuse:

- inside one script file when only that script uses the logic;
- inside one script area when several scripts in that owner area use the logic;
- inside a domain folder when one domain's scripts use the logic;
- inside a repo, package, or product helper area only after real cross-boundary reuse exists;
- inside the Operating Kit only for generic agent or workflow doctrine tooling, not
  product-specific mechanics.

Do not promote helpers upward only because future reuse is imaginable. Promote upward when two or
more durable owners actually need the same helper contract or an approved plan makes that reuse
certain.

## Domain Folder Rules

Add domain folders under a script area when one or more are true:

- three or more scripts form a cohesive domain;
- two domains have different approval or safety boundaries;
- scripts have domain-specific fixtures or test setup;
- filenames become long prefix lists that hurt navigation;
- different maintainers or review expectations apply;
- domain-specific helpers exist and should not be confused with owner-wide helpers.

Do not add domain folders only because a taxonomy looks tidy or a single script might someday
need friends.

Role folders such as `primitives/`, `workflows/`, or `helpers/` are optional. Prefer clear owner
and domain placement first. Add role folders only when they improve review, navigation, safety, or
fixtures for an existing script set. A compact `scripts/README.md` role index is usually enough
before folder separation is justified.

## Package And CLI Promotion

Promote reusable script assets into a package, module, CLI, or similar command surface only when
one or more are true:

- scripts need stable reusable internal APIs;
- multiple scripts import the same helpers heavily;
- external consumers need a versioned command surface;
- testing plain scripts becomes awkward or unreliable;
- distribution, install, or runtime materialization behavior matters;
- backwards compatibility becomes a real concern;
- the command surface is used outside one local runbook family.

Do not create a package or CLI only to make a first script look polished.

Do not promote a workflow script to a thin command wrapper only because it orchestrates phases. A
thin command wrapper is justified by a stable command surface, repeated consumers, distribution
needs, or compatibility expectations, not by composition alone.

## Script Contract

Every reusable script asset should have:

- declared script asset role when committed as a durable asset;
- one primary purpose;
- explicit required inputs;
- explicit optional inputs;
- read-only, dry-run, write, or mixed behavior clearly labeled;
- stable success output;
- stable blocker output;
- nonzero failure behavior when appropriate;
- no secrets in normal output;
- documented local tooling prerequisites;
- a runbook caller;
- proportional tests or fixtures.

For import-only helpers, the "caller" is the importing script or script area, and the output
contract may be a tested function, module, or data contract instead of standalone stdout. The
runbook caller and stdout output contract belong to the executable script entrypoint that imports
the helper.

Workflow scripts should also have:

- stable phase names;
- documented dependency contracts;
- phase-level success and blocker behavior;
- no hidden approval, target broadening, or route fallback;
- tests or fixtures that prove phase ordering and dependency wiring.

For external, sensitive, or state-changing scripts, also require:

- approval packet ID or equivalent runbook-provided approval reference when applicable;
- exact target scope inputs;
- no target broadening inside the script;
- no hidden fallback to broader permissions or different targets;
- idempotency or cleanup expectations;
- structured failure classification.

## Output Contract

Use a small, mechanical output contract. Scripts should not be asked to make AI-like or
human-like judgments about real-world sensitivity, business meaning, or user intent. They should
report what they are designed to do, what happened at runtime, and what they chose to emit.

The script author defines the contract at implementation time:

- stable `script_id`;
- declared script asset role when relevant;
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

Minimum common fields:

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

External, sensitive, or state-changing scripts may add fields such as `runbook_caller`,
`target_summary`, `action_summary`, `phase_summary`, and `evidence_summary`.

`output_safety` describes emitted-output behavior, not a universal guarantee that no sensitive
real-world content exists. A script can reliably say it did not emit raw provider responses,
tokens, secrets, or raw external content. It cannot reliably decide whether every display name,
site title, subject line, or provider error message is semantically sensitive unless the script
controls and filters that output.

Default output should be stdout. Save output to a local evidence file only when the runbook asks
for durable evidence. Do not commit raw local evidence, secrets, tokens, raw private content, or
provider responses by default.

## Test Expectations

Tests are part of first-script scaffolding, not an optional late cleanup.

Testing should scale with risk:

- local read-only validation script: fixture or smoke test may be enough;
- workflow script: phase-order, dependency-wiring, and blocker-routing tests are expected;
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

## Review Flags

Reviewers should flag:

- long inline implementation in a runbook;
- fragile shell, API, portal, or request-construction mechanics without a script;
- script without `scripts/README.md` in its owner area;
- durable script without a declared primitive, workflow, or helper role when the role affects
  review;
- workflow script with undocumented dependencies, phase boundaries, or blocker ownership;
- helper that acts like a runbook-named entrypoint without an entrypoint contract;
- helper promoted above the narrowest real owner boundary without proven reuse;
- script without tests or fixtures;
- script without a runbook caller;
- unclear output contract;
- hidden approval or target-scope expansion;
- raw sensitive output in normal stdout or committed evidence;
- runbook duplicating script internals;
- package, CLI, wrapper, or API created before repeated usage proves the need.
- workflow script that hides approvals, broad routing, ambiguous target selection, policy
  judgment, scope expansion, target broadening, or provider-specific fallback;
- script intended for managed or cloud orchestration without explicit inputs, stable structured
  outputs, non-interactive behavior, artifact/state boundaries, idempotency expectations when
  applicable, and non-secret phase evidence.

## Non-Goals

This standard does not:

- require every runbook to have a script;
- create concrete script assets for any domain;
- define one universal programming language;
- define Foundry, Microsoft 365, AI Execution, AWS, or other domain-specific script layouts;
- add validators or scaffolds;
- authorize silent installs, external changes, sensitive reads, writes, sends, deletes, or
  cleanup actions;
- replace routing, approval, safety, public-core, or tooling-readiness standards.
