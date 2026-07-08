Last updated: 2026-07-08T12:50:24Z (UTC)
Created: 2026-07-08
Status: implementation-handoff-ready

# Operating Kit Script Asset Roles Discovery

## Discovery Status

Input state: focused Operating Kit doctrine discovery after Foundry M365 authentication and invoice
intake work exposed a recurring distinction between small deterministic script assets, larger
workflow-style scripts, shared helpers, and local ad hoc scripts.

Output target: implementation-handoff-ready. This document records the approved vocabulary,
placement guidance, review expectations, non-goals, and candidate implementation scope for the
next Operating Kit implementation planning step. It is not an implementation plan and does not
authorize editing installed consumer managed `.codeheart/kit/` copies directly.

The user approved the recommendations on 2026-07-08: workflow script naming stays recommended but
not mandatory, `scripts/README.md` indexing stays compact, `thin command wrapper` remains the L3
maturity state rather than an L2 role, workflow scripts may compose prerequisite readiness and
operation primitives, and portability guidance stays generic rather than M365- or AWS-specific.

This discovery is recipe-bearing because the likely Operating Kit source change would materially
update reusable script asset vocabulary, script placement guidance, first-script README guidance,
and recipe review criteria. It is routing-bearing only in a narrow sense: it clarifies where script
assets should live by owner boundary, but it should not create domain-specific Foundry paths as
generic doctrine.

## User Intention

The user wants a generic Operating Kit update that keeps script doctrine simple while making it
clearer for agents and reviewers. The main concern is that "reusable script asset" is currently too
broad when a module has both small operation scripts and larger deterministic workflow scripts.

Key intentions:

- anticipate workflow scripts as a normal asset role when a process is stable enough to avoid
  repeated agent thinking at every step;
- keep primitive scripts, workflow scripts, and helpers conceptually distinct without creating a
  heavy bureaucracy;
- keep `thin command wrapper` as the existing L3 maturity state, not as another L2 script role;
- make script placement rules generic and portable, not hard-coded to Foundry or Microsoft 365;
- anticipate shared helpers so helpers are not duplicated across modules unnecessarily;
- keep runbooks responsible for user conversation, approvals, routing, target selection, and
  policy judgment;
- let a stable user-facing operation use one workflow entrypoint that composes prerequisite
  access/readiness checks and the narrow operation primitive;
- allow agents to improvise local ad hoc scripts when no deterministic script exists, while making
  clear that ad hoc scripts are not reusable script assets;
- keep workflow phase boundaries and input/output contracts portable enough that the same phases
  can later move to a managed runner, CI worker, or cloud orchestration surface;
- improve human review by making script roles visible through naming and compact README indexes,
  not mandatory heavyweight metadata.

## Scope

In scope for this discovery:

- Operating Kit vocabulary for script asset roles inside the existing reusable script asset layer;
- guidance for primitive scripts, workflow scripts, and helpers;
- guidance for shared helper promotion across owner boundaries;
- guidance for workflow scripts that compose prerequisite readiness and operation primitives;
- guidance on script placement by logical owner boundary;
- guidance on role visibility through filenames and `scripts/README.md`;
- clarification that role folders such as `primitives/` and `workflows/` are optional, not the
  default starting structure;
- clarification of local ad hoc scripts versus committed reusable script assets.

Out of scope for this discovery:

- changing Foundry M365 implementation code;
- adding scripts for Microsoft user reads or other new M365 operations;
- changing the Foundry managed runtime wrapper;
- defining Foundry-specific module paths as universal Operating Kit paths;
- creating a full CLI, package, or mature API doctrine beyond the existing L3/L4 model;
- editing installed consumer managed `.codeheart/kit/` copies directly;
- release, consumer install, live Microsoft validation, or invoice batch processing.

## Current Evidence

| Source | Current fact | Impact |
| --- | --- | --- |
| `.codeheart/kit/docs/agent-interface/reference/runbook-to-script-promotion-standard.md` | Defines runbook, operational recipe, runbook code block, reusable script asset, script entrypoint, infrastructure helper, domain helper, domain folder, thin command wrapper, and mature API/tool surface. | The core taxonomy is strong, but `reusable script asset` does not distinguish primitive-style scripts from workflow-style scripts. |
| Same reference | Says a reusable script asset performs one deterministic operation or evidence step, and that runbooks retain conversation, approvals, fallback choice, target selection, and routing. | A larger deterministic workflow script can be valid, but the current wording makes that role less visible. |
| Same reference | First-script scaffolding starts with `<owner-root>/scripts/README.md` and a flat script entrypoint, adding subfolders only when they improve navigation, review, safety, ownership, or fixtures. | This supports a lightweight extension rather than a new mandatory folder taxonomy. |
| Same reference | Helper rules distinguish infrastructure helpers and domain helpers and warn against broad `utils`, `common`, or manager-style helpers. | Existing helper doctrine can be reused, but shared-helper promotion across boundaries needs clearer language. |
| `.codeheart/kit/docs/agent-interface/reference/operational-recipe-maturity.md` | L2 is reusable script asset, L3 is thin command wrapper, and maturity states are not a required ladder. | Primitive and workflow script roles should be roles within L2, not new maturity states. |
| `.codeheart/kit/docs/planning-workflows/runbooks/draft-implementation-plan.md` | Plans that create reusable script assets must state target maturity, validation tier, promotion destination, placement boundary, runbook caller, output contract, tests, and review criteria. | Implementation-plan guidance can be tightened to ask for script role and placement rationale. |
| `modules/m365-workspace/scripts/README.md` | Existing module script index lists PowerShell access scripts and Python invoice-intake scripts in one table with callers and behavior. | Real Foundry usage already benefits from a compact index; adding a light `Role` column would improve review without a new registry. |
| Recent M365 evidence threads | Stable auth onboarding benefits from a script that composes phases in one process, while user-list reads can still be ad hoc until repeated need justifies a script. | Workflow scripts are useful when the process is stable; ad hoc work should remain allowed for unscripted operations. |

## Problem Framing

The Operating Kit currently has the right maturity ladder, but one layer needs sharper internal
language. "Reusable script asset" covers both a narrow operation and a deterministic multi-phase
workflow. That is technically workable, but it makes reviews harder because a reader cannot quickly
tell whether a script is intended to be a primitive building block, a workflow entrypoint, or a
helper.

The missing distinction has practical effects:

- agents may keep reasoning through an already stable process instead of using a workflow script;
- reviewers may not see whether a large script is acceptable composition or hidden policy logic;
- agents may treat local ad hoc scripts and committed reusable script assets as the same thing;
- helpers may be duplicated locally because the doctrine does not clearly say when to promote them
  upward;
- user-facing operations that always need the same prerequisites may still be executed as a series
  of manual agent decisions instead of one deterministic workflow entrypoint;
- teams may create premature `primitives/` and `workflows/` folders when a flat or domain-first
  layout would be simpler.

The clean fix is not a new maturity ladder. The clean fix is a small script-role vocabulary inside
the existing L2 reusable script asset concept.

## Recommended Decision Ledger

### D-001 - Add Script Asset Role Vocabulary

State: review-ready recommendation.

Recommendation: Add `script asset role` as a sub-classification for reusable script assets.

Recommended roles:

- `primitive script`: a narrow deterministic operation or evidence step with explicit inputs,
  outputs, blocker behavior, tests, and a runbook caller.
- `workflow script`: a deterministic composition of stable phases, primitives, or tool calls. It
  may perform deterministic branching and phase orchestration, but must not own approval decisions,
  broad target selection, policy judgment, or hidden scope expansion.
- `helper`: imported support code used by scripts. A helper is not a runbook-facing entrypoint.
  Existing infrastructure-helper and domain-helper concepts remain the normal helper subtypes.

Do not add `command_wrapper` as a script asset role. `Thin command wrapper` remains the existing L3
maturity state for a stable command-style surface after repeated use proves that a wrapper is
justified.

Rationale: This preserves the existing L2 maturity state while giving reviewers and agents the
words they need for different script shapes.

### D-002 - Workflow Scripts Compose Primitives, They Do Not Replace Them

State: review-ready recommendation.

Recommendation: Clarify that primitive scripts are not automatically "promoted into" workflow
scripts. A workflow script may import or call primitives, import shared helpers, or contain
deterministic phase logic directly when splitting would add noise. A primitive can remain a
reusable building block after a workflow script appears.

Script dependencies are normal when they are deliberate, documented, and contract-based. Workflow
scripts should depend on stable public entrypoints or shared helpers, not copied internals or
unrecorded process state.

Rationale: Promotion should be based on risk, repetition, reviewability, and stable contracts, not
on an artificial ladder from primitive to workflow.

### D-003 - Keep Workflow Scripts Inside L2 Unless A Wrapper Is Justified

State: review-ready recommendation.

Recommendation: A workflow script can still be an L2 reusable script asset. It should become an L3
thin command wrapper only when repeated use proves the need for a stable command-style interface,
distribution boundary, compatibility promise, or broader consumer surface.

Rationale: This keeps "workflow script" from becoming a premature CLI requirement and keeps
`command_wrapper` out of the L2 role vocabulary.

### D-004 - Place Scripts By Owner Boundary

State: review-ready recommendation.

Recommendation: Operating Kit doctrine should say that reusable script assets live at the
narrowest durable owner boundary that owns the behavior, safety boundary, tests, fixtures, and
review responsibility.

Generic owner-boundary levels:

- `domain-local`: one domain, slice, feature, integration, or workflow family;
- `owner-area`: multiple domains inside one module, package, product, or source area;
- `repository-level`: multiple modules, packages, products, or source areas in one repository;
- `organization/toolkit-level`: multiple repositories or generic agent/workflow doctrine tooling.

Rationale: This gives clear placement logic without encoding Foundry-specific path rules in the
Operating Kit.

### D-005 - Prefer Domain Folders Before Role Folders

State: review-ready recommendation.

Recommendation: Do not require `primitives/`, `workflows/`, or `helpers/` folders by default.
Start flat or domain-first under the owner `scripts/` area. Add role folders only when the owner
area documents that they improve navigation, review, safety boundaries, fixture layout, or
maintainer ownership.

Recommended default:

```text
<owner-root>/
  scripts/
    README.md
    <operation-script>
    <workflow-script>
    _lib/
      <helper>
```

Domain-first example:

```text
<owner-root>/
  scripts/
    README.md
    <domain-or-slice>/
      <operation-script>
      <workflow-script>
      _lib/
        <domain-helper>
```

Rationale: The boundary between primitive and workflow can be relative. A script can be a workflow
internally while still acting as a primitive for a larger runbook. Domain and ownership usually age
better than role-only folders.

### D-006 - Make Roles Visible With Names And A Compact README Index

State: review-ready recommendation.

Recommendation: Use filenames and `scripts/README.md` for lightweight role visibility.

Filename guidance:

- workflow script filenames should generally include `workflow` when they compose multiple phases
  and are intended as a normal runbook-facing workflow entrypoint;
- primitive script filenames should name the operation; no `primitive` prefix or suffix is
  required;
- helpers should either live under `_lib/` or have clearly non-entrypoint helper names;
- role folders are optional and should not replace the README index.

Recommended compact index shape:

```text
| Script | Role | Runbook caller | Behavior |
| --- | --- | --- | --- |
```

Rationale: Human reviewers can see the intended role without maintaining a heavyweight metadata
system.

### D-007 - Keep `scripts/README.md` Lightweight

State: review-ready recommendation.

Recommendation: Preserve the existing first-script README idea, but clarify that it is a compact
index and contract summary, not a large registry. Add a `Role` column where an owner area has more
than one script role or where review clarity benefits.

Do not require stale-index validation everywhere. Add validation only where the owner area already
has static checks or the script surface is risky enough to justify it.

Rationale: The README is useful for agents and reviewers, but mandatory heavy metadata would
become stale and would discourage script promotion.

### D-008 - Clarify Shared Helper Promotion

State: review-ready recommendation.

Recommendation: Anticipate shared helpers, but promote them upward only after real cross-boundary
reuse or a clearly accepted shared contract.

Helper placement rules:

- keep helpers at the narrowest boundary that currently owns the behavior;
- use clear inputs and outputs so a helper can move upward later;
- do not copy a helper into multiple modules when a second real consumer appears and the behavior
  is generic enough to share;
- promote to repository-level or organization/toolkit-level only when ownership, tests, and
  compatibility expectations are clear.

Rationale: This avoids both speculative abstraction and helper duplication.

### D-009 - Clarify Ad Hoc Scripts

State: review-ready recommendation.

Recommendation: Clarify that local ad hoc scripts are allowed for exploration, diagnostics, or
unscripted one-off operations when no deterministic asset exists. They are not reusable script
assets unless committed under an owner boundary with a runbook caller, input/output contract, and
tests or fixtures.

Recommended language:

- ad hoc scripts must stay local or temporary unless promoted deliberately;
- ad hoc scripts must not be used to bypass existing runbook approvals, target constraints, secret
  handling, or safety boundaries;
- repeated ad hoc scripts are a review trigger for promotion to a primitive or workflow script;
- domain-specific managed runners, such as Foundry runtime wrappers, belong to the owning product
  or module. The Operating Kit should define the generic boundary, not one product's invocation
  command.

Rationale: Agents need freedom to solve unscripted gaps, but the path from improvisation to
durable asset must be explicit.

### D-010 - Compose Prerequisites And Operations Through Workflow Entrypoints

State: review-ready recommendation.

Recommendation: For a stable repeated user-facing operation that always needs the same prerequisite
readiness or access checks, prefer one workflow entrypoint over asking the agent to manually run
multiple commands and decide the sequence each time.

Generic pattern:

```text
operation workflow
  resolves approved target or profile from runbook-provided inputs or committed state
  runs access, auth, or readiness workflow/helper in the required mode
  runs the operation primitive
  emits one structured result

operation primitive
  performs the narrow operation
  assumes preconditions are satisfied
  returns access_required, auth_required, or readiness_required when not satisfied
```

Early in discovery or while a process is still unstable, a runbook may call two commands in order.
Once the sequence is stable and repeated, the workflow script should own the deterministic
composition so the agent does not re-think the same prerequisite chain for each normal request.

If the operation must share the same process, session, or runtime context as the prerequisite
phase, the workflow should import shared helper code rather than shelling out to a separate process
only to recreate state.

Rationale: This keeps primitives narrow, avoids duplicated prerequisite logic, and gives agents a
single normal entrypoint for common requests.

### D-011 - Preserve Managed Runner And Cloud Portability

State: review-ready recommendation.

Recommendation: Workflow scripts should expose clear phase boundaries, explicit inputs, stable
structured outputs, and explicit artifact/state references so the same phases can later move into a
managed runner, CI worker, queue worker, or cloud orchestration surface without changing the domain
contract.

Portability guidance:

- approvals, ambiguous target selection, account/profile choice, and user conversation happen
  before script execution or remain in the runbook;
- scripts avoid interactive prompts, GUI/browser assumptions, hidden local caches as core logic,
  and hard-coded user-machine paths unless the script is explicitly host-local;
- write-capable scripts use explicit target IDs, run IDs, or idempotency keys where retries are
  possible;
- local file artifacts are treated as explicit artifacts that can later map to another artifact
  store;
- logs and output include non-secret run IDs, phase names, blocker classes, and evidence fields
  instead of relying on chat memory.

Rationale: The Operating Kit should stay generic. Cloud portability should be stated as a managed
execution and orchestration principle, not as an AWS-specific rule.

## Proposed Operating Kit Patch Targets

Likely source-doc targets for a future Operating Kit implementation:

- `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md`
  - add script asset role vocabulary;
  - clarify primitive, workflow, helper, and ad hoc boundaries;
  - add workflow-script constraints and README role guidance;
  - add generic prerequisite-plus-operation workflow composition guidance;
  - clarify script dependency and helper-import expectations;
  - clarify managed-runner and cloud-orchestration portability expectations;
  - clarify owner-boundary placement and optional role folders.
- `components/agent-interface/managed/reference/operational-recipe-maturity.md`
  - clarify that primitive/workflow/helper are roles inside L2 reusable script asset, not new
    maturity states;
  - keep thin command wrapper as L3.
- `components/planning-workflows/managed/runbooks/draft-implementation-plan.md`
  - ask recipe-bearing plans to name script asset role when creating or materially changing
    reusable script assets;
  - ask for placement boundary and role-folder rationale when applicable.
- `components/planning-workflows/managed/runbooks/execute-implementation-plan.md`
  - verify that implemented scripts match the planned role, owner boundary, README index, and
    helper placement.
- `components/structure-governance/managed/reference/documentation-structure.md`
  - only if needed, cross-reference script asset role and placement guidance without duplicating
    the agent-interface standard.

Do not patch installed managed files directly inside consumer repositories unless the task is
explicitly Operating Kit sync/drift repair. The actual implementation should happen in the
Operating Kit source repository and then be consumed through the normal release/update path.

## Candidate Implementation Capability Scope - Script Asset Role Doctrine

Capability:
Operating Kit users can distinguish primitive scripts, workflow scripts, helpers, and ad hoc
scripts inside the existing reusable script asset model.

Primary workflow:
Agents and maintainers use runbook-to-script promotion, implementation planning, and
implementation execution review to decide whether a repeated recipe stays in a runbook, becomes a
primitive script, becomes a workflow script, uses helpers, or later graduates to a thin command
wrapper.

Must cover:

- controlled vocabulary for primitive script, workflow script, helper, and local ad hoc script;
- explicit statement that `thin command wrapper` remains L3, not an L2 script role;
- workflow-script boundary that allows deterministic orchestration but forbids hidden approvals,
  target broadening, policy judgment, and broad routing;
- contract-based script dependency guidance;
- generic composition pattern for prerequisite readiness plus operation primitive workflows;
- placement by narrowest durable owner boundary;
- domain-first and flat-first layout guidance;
- optional role folders only when justified;
- lightweight `scripts/README.md` role index guidance;
- shared-helper promotion guidance;
- managed-runner and cloud-orchestration portability guidance;
- planning and execution review checks for role, placement, README, and helper boundaries.

Explicitly out of scope:

- Foundry-specific script paths as generic Operating Kit rules;
- new universal validators or generated registries;
- new script assets in Foundry modules;
- runtime wrapper changes;
- live external validation.

Preserve decisions:

- `D-001` through `D-011`.

Planner must not reinvent:

- L2/L3/L4 maturity states;
- runbook ownership of approvals, user conversation, routing, target selection, and fallback
  choice;
- existing helper distinction between infrastructure helper and domain helper;
- the generic nature of Operating Kit doctrine. M365, AWS, Azure, GCP, CI workers, and local
  runners may be examples or consumers, but not hard-coded doctrine.

Feature-level success evidence:

- Operating Kit docs define script asset roles without changing maturity states;
- Operating Kit docs explicitly keep thin command wrappers as L3;
- draft-plan and execute-plan runbooks ask for role and placement checks when scripts are
  created or materially changed;
- example README table shape includes `Role`;
- workflow composition guidance explains how a normal operation can run prerequisite readiness and
  an operation primitive through one workflow entrypoint;
- portability guidance explains how workflows remain liftable to managed runners or cloud
  orchestration surfaces;
- the docs clearly say ad hoc scripts remain allowed but are not reusable assets until promoted.

## Risks And Mitigations

| Risk | Impact | Mitigation |
| --- | --- | --- |
| The update becomes too bureaucratic. | Teams may avoid promoting scripts. | Keep roles as vocabulary and README guidance, not mandatory heavy metadata. |
| Role folders become default cargo-cult structure. | Small script surfaces get noisy. | Prefer flat or domain-first layouts; role folders require owner-area rationale. |
| Workflow scripts hide judgment. | Scripts could bypass approvals or scope decisions. | State that runbooks own approvals, routing, target selection, and policy judgment. |
| Workflow scripts duplicate prerequisite logic. | Auth, readiness, or setup behavior can drift across operations. | Prefer shared helpers or stable public entrypoints; document dependencies in the script index. |
| Scripts depend on local interactive behavior. | Later managed-runner or cloud execution becomes hard. | Require explicit inputs, non-interactive execution, stable output, and clear artifact/state boundaries. |
| Shared helpers are promoted too early. | Generic helpers become vague and hard to maintain. | Promote upward only after real reuse or an accepted shared contract. |
| Ad hoc scripts are mistaken for approved assets. | Agents may bypass deterministic paths or safety boundaries. | Clarify that ad hoc scripts are temporary/local unless deliberately promoted. |

## Open Questions

### OQ-001 - Should `workflow` In Filenames Be Required?

Recommendation: no. It should be the default naming signal when a script is clearly a
runbook-facing workflow entrypoint, but not an absolute rule.

Outcome: approved on 2026-07-08.

BLOCKER: no.

### OQ-002 - Should Every Script README Require A `Role` Column?

Recommendation: require it only when the owner area has multiple script roles or the role would not
be obvious to a reviewer. Strongly recommend it for new or changed script indexes.

Outcome: approved on 2026-07-08.

BLOCKER: no.

### OQ-003 - Should `_lib/` Be The Required Helper Folder Name?

Recommendation: no. Use `_lib/` as the recommended default helper folder name or allow a
documented local equivalent. The underscore should mean "internal support code, not a runbook-facing
entrypoint."

Outcome: approved on 2026-07-08.

BLOCKER: no.

### OQ-004 - Should Managed Ad Hoc Execution Be Generic Operating Kit Doctrine?

Recommendation: no. The Operating Kit should define the generic distinction between local ad hoc
scripts and reusable script assets. Product-specific managed runners and command names belong to
the owning repo, module, or runtime product.

Outcome: approved on 2026-07-08.

BLOCKER: no.

## Readiness Assessment

This discovery is implementation-handoff-ready. The proposed role vocabulary, naming guidance,
helper promotion rules, workflow composition model, portability guidance, and ad hoc script
boundary have user approval and no blocker-marked open questions remain.

The next meaningful step is to draft an Operating Kit implementation plan that patches the source
doctrine and then lets this Foundry repo consume the updated Operating Kit through the normal
managed update path.
