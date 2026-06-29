Last updated: 2026-06-29T15:06:25Z (UTC)
Created: 2026-06-26
Status: completed
Execution log: runtime-materialization-hardening_execution_log.md

# Document Header

## Consumer Runtime Materialization And Terminal Handoff Hardening Implementation Plan

Overview: Harden the Operating Kit local tooling standard so consumer-mode runtime tooling is
materialized from durable module or package content without editable development links, generated
install metadata in managed snapshots, or mutation of global runtimes. The durable placement for
the generic rule is the managed tooling-readiness runbook, specifically the local-machine layer
and baseline tooling catalog. Python is the first strengthened lane because AI Execution exposed a
real failure, but the plan keeps the standard general enough for future Node, PowerShell, browser,
document/PDF, and binary/archive lanes.

This plan also adds an Operating Kit agent-interface rule for terminal interactions that require
user input. Agent-run terminal tools are execution surfaces, not reliable user-visible terminals.
When a runbook requires the user to type or paste a value into a terminal prompt, especially a
secret, the agent must use a visible-terminal handoff: prepare the command, explain the exact
remaining user action, and avoid asking the user to interact with a hidden tool prompt.

The plan includes the first adopter: the Foundry `ai-execution` module. AI Execution should not
own the generic environment or terminal handoff standard and should not depend explicitly on the
Operating Kit. It should remove contradictory consumer editable-install guidance, declare its
package facts and validation commands, use the terminal handoff rule for keychain auth setup, and
release a hardened module version. The plan then refreshes the HQ managed AI Execution snapshot
and updates the released Operating Kit in Codeheart-HQ, Codeheart Automation Foundry, the named
private platform repository, and Codeheart Operating Kit.

No separate discovery is required. The user direction and the accepted local runtime environment
standard already settle the key decisions: Operating Kit owns the ambient agent tooling standard,
module snapshots remain durable and clean, consumer installs avoid editable/dev links, and module
runbooks declare requirements without becoming environment doctrine.

Essential context:

| Source | Why it matters |
| --- | --- |
| `docs/repo/plans/local-runtime-environment-standard/local-runtime-environment-standard_discovery_doc.md` | Accepted decisions for `.codeheart/local/`, repo-local Python runtime, machine/user-level tooling, and tooling-readiness routing. |
| `docs/repo/plans/local-runtime-environment-standard/local-runtime-environment-standard_implementation_doc.md` | Completed source implementation that added `.codeheart/local/`, broadened tooling readiness, and deferred AI Execution adoption. |
| `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` | Durable home for the generic runtime materialization rule and baseline tooling catalog. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | DRY boundary for local tooling guidance versus module-owned package facts and the durable standard for user-visible terminal handoff in runbooks. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Routing-bearing reference that must keep owner/scope selection before runtime setup. |
| `components/agent-interface/managed/kit-readme.md` | Installed fallback inventory that should expose any new managed runtime-materialization route language. |
| `templates/agents/AGENTS.managed-block.md` | Installed root route that must continue to route missing local tooling through tooling readiness. |
| `components/structure-governance/managed/reference/managed-content-boundaries.md` | Boundary for managed snapshots versus ignored local runtime/tooling state. |
| `docs/repo/reference/placement-contract.md` | Source placement rules for local-machine state, managed content, generated surfaces, and packaged resources. |
| `docs/repo/reference/consumer-impact-classification.md` | Required impact classification for managed docs, generated behavior, release notes, and consumer sync. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer procedure for changing Operating Kit source, docs, templates, schemas, and tests. |
| `docs/repo/runbooks/release-operating-kit.md` | Required procedure for public Operating Kit release publication and release evidence. |
| `tests/test_packaging_resources.py` | Protects source-to-packaged-resource parity for changed managed files. |
| `tests/test_onboard.py` | Protects installed root and fallback route visibility. |
| `../Codeheart-Automation-Foundry/modules/ai-execution/runbooks/install-ai-execution.md` | Current AI Execution install runbook still uses consumer editable-install validation and must be hardened. |
| `../Codeheart-Automation-Foundry/modules/ai-execution/runbooks/onboard-ai-execution.md` | Current onboarding runbook still points agents back to editable install when the package is missing. |
| `../Codeheart-Automation-Foundry/modules/ai-execution/reference/managed-snapshot-install-contract.md` | Snapshot contract already says managed snapshots exclude editable-install metadata; it needs alignment with runtime install validation. |
| `../Codeheart-Automation-Foundry/modules/ai-execution/module.yaml` | Source for AI module version, package name, console script, dependency advertisement, and route cards. |
| `../Codeheart-HQ/.codeheart/foundry/foundry.lock.yaml` | HQ managed snapshot lockfile that must be updated after the AI Execution module release. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Implement and release a hardened consumer runtime materialization standard and terminal handoff
standard so agents do not repair missing local tooling with brittle editable installs and do not
ask users to type secrets or other values into hidden agent tool terminals.

Completion is proven when:

- `handle-tooling-readiness.md` contains a generic runtime materialization rule in the local
  machine layer or adjacent section;
- the baseline tooling catalog is the durable home for future runtime lanes;
- the `python-runtime` lane states consumer-mode package tooling uses non-editable installs and
  editable installs are development-mode only;
- managed runbook-authoring guidance defines the hidden-terminal problem and the required
  visible-terminal handoff pattern;
- runbooks that require user-entered terminal values, including AI Execution auth setup, tell the
  agent to prepare or open a visible terminal surface, run or provide the command there, and tell
  the user exactly what remains to paste/type and when to press Enter;
- the generic rule is not overfit to wheel commands and does not make Foundry modules explicitly
  depend on the Operating Kit;
- AI Execution install and onboarding runbooks no longer instruct consumer agents to use
  `pip install -e` from `.codeheart/foundry/modules/ai-execution`;
- AI Execution runbooks declare package facts, entrypoint, and smoke commands while leaving local
  runtime materialization to the ambient operating/tooling layer;
- AI Execution module version is bumped for a hardened module release;
- Codeheart-HQ receives the hardened AI Execution managed snapshot and lockfile update;
- HQ local runtime validation proves `foundry-ai --help`, dry-run smoke, and
  `foundry-ai auth status --json` work from `.codeheart/local/envs/python/` before auth setup;
- Operating Kit release notes and release assets include the runtime materialization hardening;
- the released Operating Kit is synced/installed into Codeheart-HQ, Codeheart-Automation-Foundry,
  the named private platform repository, and Codeheart-Operating-Kit;
- validation covers managed docs, packaged mirrors, source tests, Foundry AI module tests, HQ
  snapshot manifest, release assets, install checks, and fresh low-context routing behavior.

## 1.2 Project And Problem Context

AI Execution onboarding exposed a runtime fragility after the HQ repo-local venv was created under
the correct Operating Kit path. The consumer installation used an editable package install from the
managed snapshot. Later cleanup removed generated editable-install metadata from the managed
snapshot, which was correct for snapshot cleanliness but broke the console script inside the venv.

The failure shows a general agent-environment rule, not just an AI Execution problem:

- managed snapshots and module source trees are durable content;
- generated runtime state belongs under `.codeheart/local/` or inside repo-local runtimes;
- consumer-mode tooling must not depend on editable links, generated source-tree metadata, or
  mutable development checkouts;
- development-mode installs remain valid only when the task explicitly targets source
  development in the owning source repository.

The Operating Kit already owns local tooling readiness and the baseline tooling catalog. The plan
should strengthen that standard there, then apply it to AI Execution as the first adopter without
making Foundry modules explicitly depend on Operating Kit internals.

## 1.3 Current State Analysis

Current Operating Kit state:

- `.codeheart/local/` is the ignored local machine/runtime layer.
- `.codeheart/local/envs/python/` is the default repo-local Python venv convention.
- `handle-tooling-readiness.md` owns local tooling blocker routing and the baseline tooling
  catalog.
- Runbook authoring currently requires user-facing flow and approval wording, but it does not yet
  state that agent-run terminal tools may be hidden from the user and therefore cannot be used for
  user-entered terminal prompts.
- The current POSIX Python command shape still shows an editable install:
  `python -m pip install -e <module-or-package-path>`.
- The `python-runtime` catalog row names the venv convention but does not yet state
  consumer-mode non-editable package materialization.
- Runbook-authoring guidance says generic tooling guidance is centralized in tooling readiness
  while modules own exact package names, versions, and smoke validation.

Current Foundry AI Execution state:

- `module.yaml` declares package `foundry-ai-execution`, Python module
  `foundry_ai_execution`, and console script `foundry-ai`.
- `install-ai-execution.md` excludes generated caches and editable metadata from managed snapshot
  export.
- The same install runbook still validates consumer packaging with
  `python -m pip install -e .codeheart/foundry/modules/ai-execution`.
- `onboard-ai-execution.md` sends a missing package back to editable install guidance.
- HQ has an AI Execution managed snapshot whose content matches the clean Foundry source commit,
  but local auth setup was blocked when editable-install metadata was removed from the managed
  snapshot.

Target state:

- Operating Kit explains the generic runtime materialization rule once.
- The baseline tooling catalog remains the single entry point for runtime lanes.
- Python lane gives the minimum generic constraint: consumer-mode package tooling is
  non-editable; editable installs are development-mode only.
- AI Execution runbooks no longer contradict the Operating Kit standard.
- AI Execution onboarding uses a visible-terminal handoff for `foundry-ai auth setup` instead of
  implying that a user can paste a key into an agent-hidden terminal prompt.
- AI Execution release and HQ install prove the end-to-end consumer path.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Implement the Operating Kit standard first, then the AI Execution adopter, then releases and
consumer installs. Keep the standard small: define materialization boundaries and mode rules, but
do not create a full package manager, wheel builder, or module dependency framework.

Expected Operating Kit source changes:

```text
Codeheart-Operating-Kit/
  components/
    agent-interface/
      managed/
        kit-readme.md                                      # modify if route inventory needs wording
        reference/
          operation-routing-and-dispatch.md                # modify routing example if needed
          runbook-authoring-standard.md                    # modify DRY/runtime wording
        runbooks/
          handle-tooling-readiness.md                      # modify primary standard
    structure-governance/
      managed/
        reference/
          managed-content-boundaries.md                    # modify if durable/runtime wording needs tightening
  src/codeheart_operating_kit/resources/
    components/...                                         # mirror changed managed files
  templates/
    agents/AGENTS.managed-block.md                         # review, modify only if root route needs wording
  tests/
    test_packaging_resources.py                            # parity coverage
    test_onboard.py                                        # installed route visibility if root wording changes
  docs/
    README.md                                              # add this plan
    repo/
      README.md                                            # add this plan
      plans/
        README.md                                          # add this plan
        plan-register.md                                   # add OK-PR-017
        runtime-materialization-hardening/
          runtime-materialization-hardening_implementation_doc.md
          runtime-materialization-hardening_execution_log.md # during execution
```

Expected Foundry source changes:

```text
Codeheart-Automation-Foundry/
  modules/
    ai-execution/
      README.md                                            # update version/release note if needed
      module.yaml                                          # bump module version, package facts unchanged
      reference/
        managed-snapshot-install-contract.md               # align snapshot/runtime boundary
      runbooks/
        install-ai-execution.md                            # remove consumer editable install
        onboard-ai-execution.md                            # route missing package to consumer-mode install
        troubleshoot-ai-execution.md                       # add setup diagnostic if useful
      tests/                                               # add/adjust validation if useful
  docs/
    repo/
      plans/
        plan-register.md                                   # pointer or module-release entry
```

Expected HQ consumer changes during execution:

```text
Codeheart-HQ/
  .codeheart/
    foundry/
      foundry.lock.yaml                                    # update ai-execution lock entry
      modules/
        ai-execution/                                      # refresh managed snapshot
        ai-execution.manifest.sha256                       # refresh manifest
    local/
      envs/python/                                         # ignored runtime, may be created/updated
  AGENTS.md                                                # root route only if needed
```

Expected named Operating Kit installs after release:

```text
Codeheart-HQ/.codeheart/kit/
Codeheart-Automation-Foundry/.codeheart/kit/
<named-private-platform-repo>/.codeheart/kit/
Codeheart-Operating-Kit/.codeheart/kit/
```

## 2.2 Open Questions And Assumptions Requiring Clarification

### OQ-1 - Should This Be A Discovery First?

BLOCKER: no

Recommended default: no. User direction and prior accepted local-runtime decisions are enough for
implementation planning. This plan records the assumptions explicitly instead of reopening
discovery.

### OQ-2 - What Is The Durable Home For Runtime Lanes?

BLOCKER: no

Recommended default: `components/agent-interface/managed/runbooks/handle-tooling-readiness.md`,
with generic rules in `Local Machine Layer` or a new `Runtime Materialization` section and lane
summaries in `Baseline Tooling Catalog`.

Rationale: this is already the central route for local package managers, runtimes, CLIs, and
repo-local Python tooling. Future lanes should land there first rather than in module runbooks.

### OQ-3 - Should Operating Kit Say "Build A Wheel"?

BLOCKER: no

Recommended default: no. Operating Kit should say consumer-mode Python package tooling uses
non-editable installs and does not depend on generated metadata inside durable managed content.
Wheel build/install can be an implementation tactic, not the doctrine.

### OQ-4 - Should AI Execution Mention Operating Kit?

BLOCKER: no

Recommended default: avoid explicit dependency wording. AI Execution should declare package facts,
entrypoint, and smoke commands. It can say the consumer runtime is selected by the repository's
local operating/tooling standard and that consumer-mode installs are non-editable.

### OQ-5 - Which Versions Should Release?

BLOCKER: no

Recommended default: use the next patch versions unless current repository state shows a newer
unreleased version at execution time.

- Operating Kit target: next patch after `0.1.14`, expected `0.1.15`.
- AI Execution target: next patch after `0.1.0`, expected `0.1.1`.

Execution must confirm tags, version files, and release history before finalizing numbers.

### OQ-6 - What Installs Are Included?

BLOCKER: no

Recommended default: include all user-named Operating Kit installs after the public Operating Kit
release, and include the HQ AI Execution snapshot refresh as the AI module consumer proof.

Named Operating Kit installs:

- Codeheart-HQ;
- Codeheart-Automation-Foundry;
- named private platform repository;
- Codeheart-Operating-Kit.

AI Execution consumer install:

- Codeheart-HQ managed snapshot and repo-local runtime validation.

## 2.3 Architectural Decisions With Reasoning

### AD-1 - Runtime Materialization Is Operating Kit Doctrine

Problem being solved: agents need one ambient rule for preparing local tooling without every
module becoming an environment manual.

Decision: Operating Kit owns the generic consumer runtime materialization standard.

Rationale: Operating Kit already owns agent environment readiness, local-machine state, and
missing-tool routing. This keeps modules environment-agnostic.

### AD-2 - Baseline Tooling Catalog Is The Lane Registry

Problem being solved: Python-specific wording should not become an arbitrary one-off placement.

Decision: place the generic rule in tooling readiness and runtime-specific summaries in the
existing baseline tooling catalog.

Rationale: future lanes have a durable home. Python is simply the first lane strengthened by a
real failure.

### AD-3 - Consumer Mode And Development Mode Are Distinct

Problem being solved: editable installs are useful during development but brittle for consumer
runtime use.

Decision: define two modes:

- `consumer-mode`: installs must not depend on editable source links, generated metadata in
  managed snapshots, or mutable development checkouts;
- `development-mode`: editable/source-linked installs are allowed only when the task explicitly
  targets source development in the owning source repository.

Rationale: this avoids banning editable installs where they are useful while preventing them from
leaking into consumer repo onboarding.

### AD-4 - AI Execution Is An Adopter, Not The Environment Authority

Problem being solved: AI Execution currently contains a consumer editable-install command that
conflicts with the managed snapshot contract.

Decision: patch AI Execution runbooks so they declare package and validation facts, remove
consumer editable-install guidance, and release a hardened module version.

Rationale: the module remains environment-agnostic but no longer contradicts the Operating Kit
standard.

### AD-5 - Release And Install Are In Scope

Problem being solved: doc-only source hardening would not reach the repos that agents operate in.

Decision: the plan includes Operating Kit public release, AI Execution module release, HQ AI
Execution snapshot refresh, and Operating Kit install/sync into four named repos.

Rationale: the failure occurred during real consumer onboarding. The fix is complete only when
the installed surfaces and first adopter are updated.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Operating Kit runtime materialization doctrine is hardened in tooling readiness and the baseline catalog. | M | none |
| EP-02 | Operating Kit runbook-authoring doctrine defines visible-terminal handoff for user-entered terminal prompts. | M | none |
| EP-03 | Operating Kit packaged mirrors, tests, release notes, and release assets are prepared. | M | EP-01, EP-02 |
| EP-04 | Foundry AI Execution runbooks and module metadata adopt the hardened consumer install and terminal handoff model. | M | EP-01, EP-02 decisions |
| EP-05 | AI Execution module release and HQ managed snapshot/runtime proof are completed. | M | EP-04 |
| EP-06 | Released Operating Kit is installed/synced into HQ, Foundry, the named private platform repository, and Operating Kit repos. | M | EP-03 release |
| EP-07 | Cross-repo registers, work board, execution log, validation, and review gates are completed. | S | EP-01 through EP-06 |

## EP-01 - Operating Kit Runtime Materialization Standard

### A) Epic ID, Title, And Outcome

EP-01 - Operating Kit Runtime Materialization Standard

Outcome: tooling readiness clearly tells agents that consumer-mode runtime tooling is
materialized into ignored local runtime state, not editable-linked from durable managed snapshots
or mutable source trees.

### B) Scope

In scope:

- generic runtime materialization rule;
- consumer-mode versus development-mode language;
- Python lane hardening in the baseline tooling catalog;
- DRY/runbook-authoring boundary updates;
- routing example update if needed;
- fresh low-context routing probe for the changed route behavior.

Out of scope:

- full wheel-building recipe;
- package manager implementation;
- module dependency framework;
- exact package versions for modules;
- hosted/cloud runtime policy.

### C) Files Touched

```text
components/agent-interface/managed/runbooks/handle-tooling-readiness.md # modify
components/agent-interface/managed/reference/runbook-authoring-standard.md # modify
components/agent-interface/managed/reference/operation-routing-and-dispatch.md # review/modify
components/agent-interface/managed/kit-readme.md # review/modify if inventory wording changes
components/structure-governance/managed/reference/managed-content-boundaries.md # review/modify
docs/repo/reference/placement-contract.md # review/modify
templates/agents/AGENTS.managed-block.md # review/modify only if root route needs wording
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- `handle-tooling-readiness.md` defines runtime materialization without overfitting to Python
  wheels.
- `Baseline Tooling Catalog` is explicitly named as the durable home for future runtime lanes.
- `python-runtime` lane says consumer-mode Python package installs are non-editable and editable
  installs are development-mode only.
- The previous generic POSIX command shape no longer uses `pip install -e`.
- Operating Kit wording says durable module content declares requirements while local runtime
  materialization happens in ignored local state.
- Runbook authoring guidance says module runbooks declare tool/package facts and smoke commands,
  not generic environment procedure.
- Routing guidance still says select owner/scope first, then use tooling readiness.
- Fresh low-context probe confirms an agent avoids editable install from managed snapshot when
  `foundry-ai` is unavailable in a consumer repo.

### E) Dependencies And Critical-Path Notes

This epic must land before the AI Execution adopter patch, because the module runbook should point
away from environment doctrine and toward a stable generic standard.

### F) Tasks Checklist

- [x] Add a concise runtime materialization rule to `handle-tooling-readiness.md`.
- [x] Replace the editable-install example in the local Python command shape with a
  consumer-safe non-editable package-install shape.
- [x] Update the `python-runtime` lane with consumer-mode and development-mode rules.
- [x] Add wording that future runtime lanes belong in `Baseline Tooling Catalog` before module
  runbooks duplicate environment guidance.
- [x] Update `runbook-authoring-standard.md` so module runbooks own package facts, package names,
  entrypoints, and smoke validation, while tooling readiness owns generic environment procedure.
- [x] Review `operation-routing-and-dispatch.md` and patch the local tooling example if it still
  permits ad hoc install selection before owner/scope.
- [x] Review managed boundaries and placement contract for generated metadata and runtime
  install-state wording.
- [x] Run a fresh low-context routing probe and record evidence in the execution log.

### G) Implementation Notes

Use short, reusable language. Avoid naming AI Execution in the generic Operating Kit standard
except as optional execution-log evidence. The standard should be useful for future runtime lanes.

### H) Open Questions

No blockers.

## EP-02 - Operating Kit Visible Terminal Handoff Standard

### A) Epic ID, Title, And Outcome

EP-02 - Operating Kit Visible Terminal Handoff Standard

Outcome: managed runbook-authoring doctrine tells agents how to handle terminal prompts that need
user input when the agent's terminal/tool execution surface is not visible or interactive for the
user.

### B) Scope

In scope:

- define hidden agent terminal/tool prompts as non-user-interactive surfaces;
- require visible-terminal handoff whenever the user must type or paste a value into a terminal
  prompt;
- require exact remaining-user-action wording, such as what to paste/type and when to press
  Enter;
- require agents to stop before secrets are entered into chat or hidden tool prompts;
- update tooling readiness and/or runbook-authoring guidance so install/auth runbooks know the
  pattern;
- add validation through a fresh low-context prompt involving keychain setup.

Out of scope:

- implementing a GUI terminal controller;
- prescribing one terminal application for every platform;
- handling passwords, MFA, or secrets in agent-visible logs;
- changing native Codex tool behavior.

### C) Files Touched

```text
components/agent-interface/managed/reference/runbook-authoring-standard.md # modify primary standard
components/agent-interface/managed/runbooks/handle-tooling-readiness.md # modify if terminal handoff appears in install/repair flow
components/agent-interface/managed/reference/operation-routing-and-dispatch.md # review/modify if routing examples need terminal handoff
components/agent-interface/managed/kit-readme.md # review/modify if fallback inventory needs wording
templates/agents/AGENTS.managed-block.md # review/modify only if root rule needs visibility
tests/test_onboard.py # modify if root managed block changes
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Runbook-authoring standard states that agent-run terminal tools are not assumed visible or
  interactive for users.
- Any runbook that needs a user to enter terminal input must use a visible-terminal handoff.
- The handoff pattern tells the agent to prepare the command, identify the visible terminal
  surface or provide the command for the user's terminal, and tell the user exactly what remains
  to paste/type and when to press Enter.
- Secret entry is never requested in chat and never delegated to an agent-hidden prompt.
- The standard distinguishes visible-terminal handoff from normal agent-run noninteractive
  commands.
- Fresh low-context probe confirms an agent does not run `foundry-ai auth setup` in a hidden tool
  prompt while asking the user to paste an API key.

### E) Dependencies And Critical-Path Notes

This epic is independent of EP-01, but AI Execution onboarding adoption in EP-04 depends on it.
Keep the rule generic; AI Execution is an example, not the owner of the standard.

### F) Tasks Checklist

- [x] Add visible-terminal handoff doctrine to `runbook-authoring-standard.md`.
- [x] Add or cross-reference the terminal handoff rule in `handle-tooling-readiness.md` where
  install/auth setup can hit interactive terminal prompts.
- [x] Review root route wording and update only if agents need an immediate safety reminder.
- [x] Add a fresh low-context terminal-handoff probe and record evidence in the execution log.
- [x] If root route wording changes, update installed route visibility tests. Not applicable:
  root route wording was reviewed and left unchanged.

### G) Implementation Notes

The wording should make clear that "open terminal and run command" means a user-visible terminal
surface, not an agent-only `exec_command` or hidden PTY. When a visible terminal cannot be opened
or controlled safely, the agent should provide the exact command and working directory for the
user to run in their own terminal and then wait for the user to report completion.

### H) Open Questions

No blockers.

## EP-03 - Operating Kit Packaging, Release Prep, And Public Release

### A) Epic ID, Title, And Outcome

EP-03 - Operating Kit Packaging, Release Prep, And Public Release

Outcome: the hardened runtime materialization standard is mirrored into packaged resources,
validated, released publicly, and ready for named repo sync.

### B) Scope

In scope:

- packaged resource mirrors for changed managed docs/templates;
- parity tests for any changed managed files;
- release notes with consumer-impact classification for runtime materialization and
  visible-terminal handoff;
- package/version bump;
- release asset build and validation;
- public GitHub release after approval gates pass.

Out of scope:

- AI Execution module release;
- named repo sync, which is EP-06;
- changing installer architecture.

### C) Files Touched

```text
src/codeheart_operating_kit/resources/components/... # modify mirrors
src/codeheart_operating_kit/resources/templates/... # modify if template changed
tests/test_packaging_resources.py # modify if needed
tests/test_onboard.py # modify if root route visibility changes
pyproject.toml # version bump
bootstrap.md # release asset content/version
install.sh # release asset content/version if builder updates it
install.ps1 # release asset content/version if builder updates it
release-notes.md # release notes
manifest.yaml # release manifest
src/codeheart_operating_kit/resources/manifest.yaml # packaged release metadata
dist/ # release artifacts
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- changed managed files have matching packaged resources;
- release notes classify impact as backwards-compatible managed-doctrine, tooling-readiness, and
  runbook-authoring change with no forced migration and no default local tool install;
- release version is confirmed at execution time, expected `0.1.15`;
- markdown, public-core, JSON schema, release manifest, focused pytest, full pytest, and
  `git diff --check` pass;
- release assets are built, checksummed, and validated;
- installer fail-closed behavior is validated;
- macOS release install proof passes;
- Windows validation is run through the release workflow or recorded as a residual risk if not
  available;
- GitHub release is published only after validation and explicit execution approval.

### E) Dependencies And Critical-Path Notes

This epic depends on EP-01 and EP-02. Do not publish a release from a dirty or unvalidated source
state.

### F) Tasks Checklist

- [x] Copy changed managed source files into matching packaged resource paths.
- [x] Update packaging parity tests for changed managed files absent from the parity list.
- [x] Update release notes with consumer impact and runtime materialization summary.
- [x] Select and apply the next Operating Kit patch version after checking current tags.
- [x] Build release assets.
- [x] Validate release manifest and asset checksums.
- [x] Validate installer fail-closed behavior.
- [x] Validate a temporary macOS consumer install from release assets.
- [x] Run full source validation.
- [x] Publish GitHub release after approval gates pass.
- [x] Record release URL, tag, commit, asset hashes, and validation evidence in the execution log.

### G) Implementation Notes

Follow `docs/repo/runbooks/release-operating-kit.md`. If the release target commit differs from
the validated commit, stop and revalidate.

### H) Open Questions

No blockers.

## EP-04 - Foundry AI Execution Adopter Hardening

### A) Epic ID, Title, And Outcome

EP-04 - Foundry AI Execution Adopter Hardening

Outcome: AI Execution no longer teaches consumer agents to use editable installs from managed
snapshots, declares package facts and validation commands that the ambient operating layer can
materialize safely, and uses visible-terminal handoff for keychain setup.

### B) Scope

In scope:

- remove consumer `pip install -e` guidance;
- add consumer-mode non-editable install wording;
- distinguish consumer-mode and source-development install modes;
- align auth setup wording with the Operating Kit visible-terminal handoff standard;
- keep module environment-agnostic and avoid explicit Operating Kit dependency language;
- bump AI Execution module version for release;
- update module docs, runbooks, and snapshot contract;
- run AI module tests.

Out of scope:

- changing AI Execution provider behavior;
- changing model profile defaults;
- adding hosted/cloud runtime;
- adding invoice-specific behavior;
- implementing a generic Foundry dependency manager.

### C) Files Touched

```text
../Codeheart-Automation-Foundry/modules/ai-execution/module.yaml # modify version
../Codeheart-Automation-Foundry/modules/ai-execution/README.md # modify version/release summary
../Codeheart-Automation-Foundry/modules/ai-execution/runbooks/install-ai-execution.md # modify
../Codeheart-Automation-Foundry/modules/ai-execution/runbooks/onboard-ai-execution.md # modify
../Codeheart-Automation-Foundry/modules/ai-execution/runbooks/troubleshoot-ai-execution.md # review/modify
../Codeheart-Automation-Foundry/modules/ai-execution/reference/managed-snapshot-install-contract.md # modify
../Codeheart-Automation-Foundry/docs/repo/plans/plan-register.md # pointer/release entry
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- no consumer-facing AI Execution runbook instructs `pip install -e` from
  `.codeheart/foundry/modules/ai-execution`;
- editable installs are described only as source-repository development/test mode;
- install runbook names package `foundry-ai-execution`, module path, console script `foundry-ai`,
  and smoke commands;
- onboarding runbook says do not request an API key until package/CLI validation succeeds;
- onboarding runbook tells the agent not to ask the user to paste a key into a hidden terminal
  prompt and instead uses visible-terminal handoff wording;
- managed snapshot contract states generated install metadata must not be required by consumer
  runtime tooling;
- AI module version is bumped, expected `0.1.1`;
- AI module tests pass;
- source tree is clean for release after generated runtime/cache metadata is removed or ignored.

### E) Dependencies And Critical-Path Notes

This epic depends on EP-01 and EP-02 decisions but can be edited before the Operating Kit release
is published. Keep the module environment-agnostic: it can say "consumer-mode non-editable
install" and use visible-terminal handoff, but should not require a specific Operating Kit file
path by doctrine.

### F) Tasks Checklist

- [x] Update AI Execution `module.yaml` to the selected hardened release version.
- [x] Update AI Execution README with the hardened release note.
- [x] Replace install-runbook editable package validation with consumer-mode non-editable install
  language.
- [x] Add module facts to install validation: package name, source path, console script, smoke
  manifest, auth status command.
- [x] Update onboarding missing-package branch so it returns to the consumer-mode install
  validation path, not editable install.
- [x] Update auth setup wording so the user enters the API key only in a visible terminal or
  self-run terminal command, not chat or a hidden agent terminal prompt.
- [x] Update managed snapshot contract so generated install metadata is not part of durable
  snapshot state and is not required by consumer runtimes.
- [x] Review troubleshooting setup guidance for the same mode boundary.
- [x] Run AI Execution test suite and dry-run smoke from source-development environment.
- [x] Run a source grep proving no consumer `pip install -e` instructions remain.

### G) Implementation Notes

The Foundry source repo can still use editable install for development/test work. The consumer
runbook should not rely on generated metadata inside `.codeheart/foundry/modules/ai-execution/`.

### H) Open Questions

No blockers.

## EP-05 - AI Execution Release And HQ Consumer Proof

### A) Epic ID, Title, And Outcome

EP-05 - AI Execution Release And HQ Consumer Proof

Outcome: the hardened AI Execution module release is installed as a managed snapshot in HQ and
validated from HQ's repo-local runtime before auth setup.

### B) Scope

In scope:

- confirm clean Foundry source commit for AI Execution release;
- install/update HQ managed snapshot;
- update HQ `foundry.lock.yaml`;
- rebuild snapshot manifest;
- create or repair HQ `.codeheart/local/envs/python/`;
- install AI Execution into the venv in consumer mode;
- prove `foundry-ai --help`, dry-run smoke, and `foundry-ai auth status --json`;
- start auth setup only after package validation succeeds and only through visible-terminal
  handoff, if user requests completion.

Out of scope:

- storing or printing an OpenAI key;
- live provider calls;
- downstream M365 writes;
- installing AI Execution into all repos unless separately approved.

### C) Files Touched

```text
../Codeheart-HQ/.codeheart/foundry/foundry.lock.yaml # modify
../Codeheart-HQ/.codeheart/foundry/modules/ai-execution/ # refresh managed snapshot
../Codeheart-HQ/.codeheart/foundry/modules/ai-execution.manifest.sha256 # update
../Codeheart-HQ/AGENTS.md # review/modify only if route missing
../Codeheart-HQ/.gitignore # review, should already ignore .codeheart/local/
../Codeheart-HQ/.codeheart/local/envs/python/ # ignored runtime, create/update
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- Foundry source commit for AI Execution release is recorded and clean;
- HQ snapshot manifest hash validates the installed snapshot;
- HQ lockfile records module version, source repo, source commit, source tree state, installed
  path, snapshot manifest path, and hash;
- managed snapshot contains no `__pycache__/`, `*.pyc`, `*.egg-info/`, `.pytest_cache/`, `.venv/`,
  or local runtime cache;
- HQ venv command resolves `foundry-ai`;
- `foundry-ai --help` passes from HQ venv;
- dry-run smoke reports `sends_provider_request: false`;
- `foundry-ai auth status --json` reports a secret-safe status, expected `missing` before key
  entry;
- auth setup is not started until package validation passes.
- if auth setup is started, the user is instructed to paste the key only into the visible
  terminal prompt and press Enter.

### E) Dependencies And Critical-Path Notes

This epic depends on EP-04. Do not ask for key entry until the consumer-mode install and CLI
validation have passed. Do not run hidden interactive auth prompts that the user cannot see.

### F) Tasks Checklist

- [x] Confirm AI Execution source tree is clean or record approved source-tree state.
- [x] Export the hardened snapshot into HQ using the cache-excluding snapshot command.
- [x] Generate HQ AI Execution snapshot manifest and manifest hash.
- [x] Update HQ `foundry.lock.yaml` for the AI Execution release.
- [x] Confirm HQ root routes discover AI Execution.
- [x] Create or reuse `.codeheart/local/envs/python/` in HQ.
- [x] Install AI Execution into the HQ venv using consumer-mode non-editable install.
- [x] Confirm the install does not leave generated metadata in the managed snapshot.
- [x] Run `foundry-ai --help`.
- [x] Run dry-run smoke against the installed smoke manifest and confirm no provider request.
- [x] Run `foundry-ai auth status --json`.
- [x] If requested, start `foundry-ai auth setup` only after the previous checks pass and only
  with visible-terminal handoff instructions.

### G) Implementation Notes

The auth setup prompt must be handled through a visible terminal handoff, not chat or a hidden
agent tool prompt. The execution log should record only `ready`, `missing`, or `env-fallback`,
never the key.

### H) Open Questions

No blockers.

## EP-06 - Named Operating Kit Installs

### A) Epic ID, Title, And Outcome

EP-06 - Named Operating Kit Installs

Outcome: the released Operating Kit with runtime materialization hardening is installed or synced
into the four named repositories.

### B) Scope

In scope:

- sync/update Codeheart-HQ;
- sync/update Codeheart-Automation-Foundry;
- sync/update the named private platform repository;
- sync/update Codeheart-Operating-Kit;
- verify `.codeheart/local/` ignore/config/route visibility;
- run `codeheart-operating-kit check` in each repo.

Out of scope:

- changing unrelated repo content;
- applying Foundry module snapshots outside HQ;
- public cloud/account changes.

### C) Files Touched

```text
../Codeheart-HQ/.codeheart/kit/** # managed sync
../Codeheart-HQ/.codeheart/kit.lock.yaml # sync metadata
../Codeheart-HQ/AGENTS.md # managed block refresh
../Codeheart-HQ/.gitignore # ensure .codeheart/local/

../Codeheart-Automation-Foundry/.codeheart/kit/** # managed sync
../Codeheart-Automation-Foundry/.codeheart/kit.lock.yaml # sync metadata
../Codeheart-Automation-Foundry/AGENTS.md # managed block refresh
../Codeheart-Automation-Foundry/.gitignore # ensure .codeheart/local/

../<named-private-platform-repo>/.codeheart/kit/** # managed sync
../<named-private-platform-repo>/.codeheart/kit.lock.yaml # sync metadata
../<named-private-platform-repo>/AGENTS.md # managed block refresh
../<named-private-platform-repo>/.gitignore # ensure .codeheart/local/

../Codeheart-Operating-Kit/.codeheart/kit/** # managed sync
../Codeheart-Operating-Kit/.codeheart/kit.lock.yaml # sync metadata
../Codeheart-Operating-Kit/AGENTS.md # managed block refresh
../Codeheart-Operating-Kit/.gitignore # ensure .codeheart/local/
```

### D) Acceptance Criteria And Size

Size: M

Acceptance criteria:

- each repo reports `OK: True` from the released `codeheart-operating-kit check`;
- no repo has managed-kit drift after sync;
- each repo has `.codeheart/local/` ignored;
- each installed managed tooling-readiness runbook contains the runtime materialization hardening;
- each installed root `AGENTS.md` routes missing local tooling to tooling readiness;
- existing unrelated dirty work in each repo is preserved.

### E) Dependencies And Critical-Path Notes

This epic depends on EP-03 release publication. Do not sync from untagged local source for final
release proof unless explicitly recorded as a local pre-release smoke.

### F) Tasks Checklist

- [x] Install/sync released Operating Kit into Codeheart-HQ.
- [x] Install/sync released Operating Kit into Codeheart-Automation-Foundry.
- [x] Install/sync released Operating Kit into the named private platform repository.
- [x] Install/sync released Operating Kit into Codeheart-Operating-Kit.
- [x] Run `codeheart-operating-kit check` in each repo.
- [x] Verify installed tooling-readiness runtime materialization wording in each repo.
- [x] Verify `.codeheart/local/` ignore/config route in each repo.
- [x] Record dirty-work preservation notes for each repo.

### G) Implementation Notes

The Operating Kit source repo now has an installed `.codeheart/kit/` state. Treat it as a named
consumer for sync validation while still respecting public-core safety and source-repo release
rules.

### H) Open Questions

No blockers.

## EP-07 - Registers, Work Board, Review, And Final Validation

### A) Epic ID, Title, And Outcome

EP-07 - Registers, Work Board, Review, And Final Validation

Outcome: cross-repo planning state and validation evidence accurately show the hardening release,
AI module adopter release, HQ proof, and named Operating Kit installs.

### B) Scope

In scope:

- Operating Kit plan register;
- HQ coordination register and work board;
- Foundry plan register pointer or module-release entry;
- execution log;
- review gate;
- validation summary;
- deferred follow-ups.

Out of scope:

- broad portfolio register cleanup;
- unrelated M365 work-board changes;
- changing prior completed plan lifecycle entries except relation pointers if needed.

### C) Files Touched

```text
docs/repo/plans/plan-register.md # Operating Kit register
docs/repo/plans/runtime-materialization-hardening/runtime-materialization-hardening_execution_log.md # during execution
../Codeheart-HQ/docs/repo/plans/plan-register.md # coordination pointer
../Codeheart-HQ/docs/repo/plans/portfolio-work-board/portfolio-work-board.md # board status
../Codeheart-Automation-Foundry/docs/repo/plans/plan-register.md # Foundry module release pointer
```

### D) Acceptance Criteria And Size

Size: S

Acceptance criteria:

- Operating Kit register has `OK-PR-017` for this implementation plan;
- `OK-PR-016` relates to `OK-PR-017`;
- HQ coordination register has `CODEHEART-OPERATING-KIT-PR-017`;
- HQ work board shows the plan under Operating Kit And Foundry System Model;
- Foundry register records the AI Execution hardening release or pointer;
- execution log records all validation, release, install, and residual-risk evidence;
- review gate finds no unresolved material findings before completion;
- final status distinguishes completed source/release/install work from deferred future lanes.

### E) Dependencies And Critical-Path Notes

Keep detailed evidence in the execution log. Registers should summarize lifecycle and pointers,
not duplicate release logs.

### F) Tasks Checklist

- [x] Add or update `OK-PR-017` in the Operating Kit register.
- [x] Add relation from `OK-PR-016` to `OK-PR-017`.
- [x] Add `CODEHEART-OPERATING-KIT-PR-017` to the HQ coordination register.
- [x] Add or refresh the HQ work board entry.
- [x] Add a Foundry plan-register pointer or module-release entry for the AI Execution adopter
  release.
- [x] Create execution log before marking epics complete.
- [x] Run final source, release, module, consumer, and install validation.
- [x] Run `git diff --check` in every touched repo.
- [x] Run a review gate over the cross-repo diff.
- [x] Mark plan/register/work-board statuses complete only after all accepted scope passes.

### G) Implementation Notes

Use `not recorded` for session refs unless a safe session identifier is already available. Do not
scan private transcripts solely to improve the register.

### H) Open Questions

No blockers.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

Additional runtime lanes:

- deferred because Python is the only lane with a concrete current failure;
- trigger: Node, PowerShell, browser automation, document/PDF, binary/archive, or vendor CLI
  runtime materialization creates repeated consumer install ambiguity.

Generic module dependency manager:

- deferred because the current fix needs mode and boundary doctrine, not a dependency resolver;
- trigger: multiple modules need machine-readable runtime requirements and automated setup.

Wheel-builder helper:

- deferred because the standard should not overfit to wheel commands;
- trigger: repeated Python consumer installs need a tested script or command wrapper.

AI Execution live provider smoke:

- deferred unless the user completes keychain onboarding and explicitly approves a live model
  request;
- trigger: auth status reports `ready` and the live-call approval gate is passed.

Operating Kit cloud or hosted runtime policy:

- deferred because this plan only handles local consumer runtime materialization;
- trigger: hosted Foundry runtime or cloud agent work re-enters active planning.

## 4.2 Future Considerations

- If multiple Foundry modules adopt the same runtime declaration shape, promote a lightweight
  module runtime requirement schema into Foundry or Operating Kit doctrine.
- If install commands keep becoming large runbook snippets, consider a reusable script asset under
  the Operating Kit runbook-to-script promotion standard.
- If the same consumer-mode/development-mode distinction appears outside tooling readiness, add a
  short reference page and keep the runbook as the operational route.

# Revision Notes

- 2026-06-26: Initial draft created from user-approved direction after AI Execution auth setup
  exposed brittle editable-install behavior in a managed consumer snapshot. Scope includes
  Operating Kit runtime materialization hardening, Operating Kit release, AI Execution adopter
  release, HQ AI snapshot/runtime proof, named Operating Kit installs, and coordination updates.
- 2026-06-26: Added Operating Kit visible-terminal handoff hardening so runbooks do not ask users
  to type secrets or other values into hidden agent tool terminals.
- 2026-06-26: Activated for execution and created the sibling execution log.
- 2026-06-29: Reworded a forward-looking future consideration to use the current reusable script
  asset vocabulary after the runbook-to-script promotion standard retired tested inline-block
  maturity.
