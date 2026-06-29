Last updated: 2026-06-29T14:49:19Z (UTC)

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
become a reusable script asset, mature into helpers or folders, or wait for wrapper/API maturity
without creating premature broad command surfaces.

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

Deprecated or avoid:

- Do not use `tested script block` as a maturity state.
- Do not use `executable script block` to mean a durable asset.
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

- script entrypoints and their intended runbook callers;
- required local tooling;
- input contract;
- output contract;
- safety and approval boundary;
- read-only, dry-run, or write behavior;
- where tests and fixtures live.

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

## Script Contract

Every reusable script asset should have:

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
`target_summary`, `action_summary`, and `evidence_summary`.

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
- script without tests or fixtures;
- script without a runbook caller;
- unclear output contract;
- hidden approval or target-scope expansion;
- raw sensitive output in normal stdout or committed evidence;
- runbook duplicating script internals;
- package, CLI, wrapper, or API created before repeated usage proves the need.

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
