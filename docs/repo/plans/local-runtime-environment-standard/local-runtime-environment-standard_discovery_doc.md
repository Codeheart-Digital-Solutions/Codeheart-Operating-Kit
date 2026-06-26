Last updated: 2026-06-26T14:10:39Z (UTC)
Created: 2026-06-26
Status: draft

# Local Runtime Environment Standard Discovery

## Discovery Status

Input state: new Operating Kit doctrine question after AI Execution onboarding exposed that
Homebrew Python blocks direct editable installs under PEP 668, and follow-up discussion clarified
that local tooling readiness should route consistently whenever an agent interaction is blocked by
missing or unsuitable local tools.

Output target: implementation-handoff-ready. This document records the accepted framing,
evidence, defaults, and implementation-shaping decisions for implementation planning.

Approval update: On 2026-06-26, the user accepted the recommended defaults for `.codeheart/local/`,
`.codeheart/local/envs/python/`, `.gitignore` handling, optional
`local_machine_layer_path` config, on-demand first-run behavior, governed machine/user-level
tooling, and generic blocker routing.

Implementation planning is now in scope. Implementation itself still requires a separate
implementation plan and approval before source changes.

## User Intention

Codeheart wants agents to have a standard, reusable way to prepare local tooling when a repo,
module, or agent task needs it. The standard should be simple enough for nontechnical users, safe
for managed consumer repos, and general enough that future modules do not invent separate setup
conventions.

The immediate trigger is `ai-execution`, but the desired standard should belong to the Operating
Kit rather than the AI Execution module. Python virtual environments are the first concrete
runtime convention, not the whole problem.

## Problem Framing

The current Operating Kit already has a generic tooling-readiness runbook in the Agent Interface
component:

```text
components/agent-interface/managed/runbooks/handle-tooling-readiness.md
```

That runbook correctly owns the generic local-tooling conversation: classify the blocker, run
read-only checks, ask approval before local changes, use official or module-owned install routes,
recheck readiness, and return to the calling runbook.

The current route is framed mainly around module onboarding and operation. The desired behavior is
broader: whenever an agent interaction is blocked by missing local tooling, the agent should route
to environment readiness, classify the blocker, use the same approval and recheck pattern, then
return to the original task. The exact tool does not need to be pre-modeled for the route to apply.

The missing piece is a concrete convention for Python virtual environments once the Python tooling
lane is selected. AI Execution revealed this gap because the runbook's direct editable install:

```text
python -m pip install -e .codeheart/foundry/modules/ai-execution
```

is blocked on Homebrew-managed Python by PEP 668. During implementation, the safe path was an
ignored virtual environment. That pattern should become explicit and reusable.

## Goals

- Define where ignored local machine/runtime state belongs in a Codeheart-managed repo.
- Define the default repo-local Python virtual environment location.
- Keep `.codeheart/user/` focused on human local preferences and notes.
- Give agents a safe default when Homebrew Python or another externally managed Python blocks
  direct package installation.
- Clarify when baseline machine/user-level tooling such as package managers, runtimes, browsers,
  and official vendor CLIs may be installed outside the repo.
- Let modules declare package requirements without each module inventing its own venv path.
- Keep generated tool environments out of committed source and managed module snapshots.
- Preserve a simple user-facing approval flow through the existing tooling-readiness runbook.
- Make the environment-readiness route available for any agent task, not only module onboarding.

## Non-Goals

- Do not implement the standard in this discovery.
- Do not patch `.gitignore`, managed kit docs, or AI Execution runbooks yet.
- Do not create a virtual environment during discovery.
- Do not define a full Python dependency lock strategy.
- Do not require every repo to create a Python venv during Operating Kit first-run onboarding.
- Do not define exact installation commands for every OS or package manager in this discovery.
- Do not force every tool into `.codeheart/local/`; baseline package managers, runtimes, browsers,
  and official vendor CLIs may still be machine/user-level tools when appropriate.

## Current Evidence

| Source | Finding | Impact |
| --- | --- | --- |
| `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` | Operating Kit already routes missing local tools through a generic readiness lane and includes `python-runtime`. | The new standard should extend this route rather than duplicate its approval conversation. |
| `components/agent-interface/managed/reference/local-extension-contract.md` | `.codeheart/user/` is for ignored local preferences and personal notes. | A generated venv does not naturally fit the existing user-layer meaning. |
| `components/structure-governance/managed/reference/managed-content-boundaries.md` | `.codeheart/user/` is the ignored local user layer; managed, consumer-owned, committed state, and generated artifacts are separate ownership concepts. | A separate ignored local machine/runtime layer is cleaner than overloading `.codeheart/user/`. |
| Installed consumer `.gitignore` patterns | `.codeheart/user/` is already commonly ignored in consumers. | Ignored local state is accepted, but the current ignored namespace is user-facing. |
| AI Execution implementation evidence | Direct `pip install -e` was blocked by an externally managed Homebrew Python; ignored venv validation succeeded. | The issue is real and already encountered during implementation, not only during onboarding. |
| AI Execution install/onboard runbooks | They route missing Python/pip tooling to Operating Kit tooling readiness but still show direct editable install as the package-validation command. | Module runbooks need to reference the shared venv standard once it exists. |
| Follow-up local-runtime discussion | Homebrew and equivalent package-manager bootstrap should stay in scope, but repo-specific packages should not be installed globally. | The standard needs a machine/user-level tooling category and a repo-local generated-runtime category. |

## Working Definitions

| Term | Draft meaning |
| --- | --- |
| Local user layer | Ignored human-facing preferences, local guidance, and personal notes under `.codeheart/user/`. |
| Local machine layer | Ignored generated runtime/tooling state that is specific to one checkout and can be recreated. |
| Machine/user-level tooling | Tools installed outside a repo because they are naturally machine-level or user-level, such as package managers, shell runtimes, language runtimes, browsers, or official vendor CLIs. |
| Repo-local tooling | Packages, generated runtimes, shims, caches, or temporary files needed for one repo or module and recreated from committed sources or runbooks. |
| Repo-local Python venv | A virtual environment created inside the repo's ignored local machine layer to run Codeheart/Foundry/agent tooling. |
| Managed snapshot | Committed module content under `.codeheart/foundry/modules/<id>/`; must not contain generated venvs or package artifacts. |

## Recommended Defaults

### Default Local Machine Boundary

Recommendation:

```text
.codeheart/local/
```

Meaning: ignored local machine/runtime state. This is separate from `.codeheart/user/`, which
remains the ignored user-facing preference and notes layer.

Why:

- `local` names the real boundary: machine-local and recreatable.
- It leaves `.codeheart/user/` human-readable and preference-oriented.
- It gives future local tooling one parent namespace.

Candidate contents:

```text
.codeheart/local/envs/python/    default Python venv
.codeheart/local/envs/<purpose>/ exception venvs when conflicts require isolation
.codeheart/local/cache/          future non-secret local tool caches
.codeheart/local/tmp/            future scratch files
.codeheart/local/bin/            future generated local shims when justified
```

Do not store secrets, raw provider responses, durable repo state, or managed snapshots under
`.codeheart/local/`.

### Source-Control And Config Defaults

Recommendation: add `.codeheart/local/` to consumer `.gitignore` during init and sync, but do not
create `.codeheart/local/` by default.

Recommendation: add an optional `local_machine_layer_path: .codeheart/local/` field under
`local_consumer_layer` for new installs. Existing configs should remain valid when the field is
absent.

Why:

- The repo is ready for generated local runtime state before the first local tool needs it.
- Empty runtime directories are not scaffolded just to satisfy a future possibility.
- Agents and runbooks get a machine-readable default path without breaking existing consumers.
- The field remains descriptive local setup metadata, not durable runtime state.

### Machine/User-Level Tooling

Recommendation: keep machine/user-level installation in scope, but tightly governed.

Use this lane only for tools that are naturally machine-level or user-level:

- package managers such as Homebrew, Windows package-manager lanes, or Linux distro package
  managers;
- shell runtimes such as PowerShell;
- language runtimes such as Python or Node;
- browsers or browser automation prerequisites;
- official vendor CLIs where the vendor expects a user or machine installation.

Rules:

- ask explicit approval before installing or repairing machine/user-level tooling;
- use official vendor, OS, or package-manager sources;
- do not install repo-specific Python packages globally;
- do not use `sudo pip`, `--break-system-packages`, or ad hoc global package installs as the
  normal path;
- recheck readiness with a read-only command before returning to the calling task.

### Default Python Virtual Environment

Recommendation:

```text
.codeheart/local/envs/python/
```

Meaning: the default repo-local Python venv for Operating Kit, Foundry, and module CLI tooling.

Typical command shape:

```sh
python3 -m venv .codeheart/local/envs/python
.codeheart/local/envs/python/bin/python -m pip install -e .codeheart/foundry/modules/ai-execution
PYTHONDONTWRITEBYTECODE=1 .codeheart/local/envs/python/bin/foundry-ai --help
```

Why:

- It avoids global Python mutation.
- It works with externally managed Python installations such as Homebrew Python.
- It creates one predictable repo-local tool runtime.
- It avoids scattering module-specific venvs for every Python-backed capability.
- It lets modules become more compatible by sharing one standard Python tooling lane when their
  dependencies can coexist.

### Exception Path For Dependency Conflicts

Recommendation:

```text
.codeheart/local/envs/<purpose>/
```

Use only when the default Python venv cannot safely host the needed packages because of dependency,
Python-version, security, lifecycle, or isolation conflicts.

Example:

```text
.codeheart/local/envs/python/
.codeheart/local/envs/ai-execution-smoke/
```

Default behavior should still prefer the shared Python venv.

## Decision Ledger

### D-001 - Add A Separate Local Machine Layer

Question: Should generated local runtime/tooling state live under `.codeheart/user/`, under a new
`.codeheart/local/`, or under a narrower `.codeheart/envs/` path?

Recommended default: introduce `.codeheart/local/`.

Rationale: `.codeheart/user/` already means human local preferences and notes. `.codeheart/envs/`
is too narrow and may conflict conceptually with deployed environments, tenants, staging/prod, or
business environment records. `.codeheart/local/` cleanly covers venvs plus future local caches,
temporary files, or generated shims without implying they are durable repo state.

Decision state: accepted.

BLOCKER: no.

### D-002 - Set The Default Python Venv Path

Question: Should the standard use one general repo-local Python venv or one venv per module?

Recommended default: one general repo-local Python venv at `.codeheart/local/envs/python/`.

Rationale: most Python-backed agent tooling should share one predictable runtime. Module-specific
venvs add cognitive load and duplicate dependency installs before there is evidence of conflicts.
Shared venv first also encourages Foundry modules to stay compatible with the standard tooling lane.

Decision state: accepted.

BLOCKER: no.

### D-003 - Keep Module-Specific Venvs As Exceptions

Question: When should a module get its own Python venv?

Recommended default: use a purpose-specific venv only when there is a concrete conflict or risk.

Valid exception triggers:

- incompatible Python version;
- dependency conflicts with existing shared tooling;
- experimental, heavy, or unstable dependencies that should not affect the shared tool runtime;
- security or sandboxing reasons to isolate dependencies;
- temporary validation fixture rather than normal user path;
- product/application runtime that is not Operating Kit or agent tooling.

Use this path shape:

```text
.codeheart/local/envs/<purpose>/
```

Rationale: a venv is not inherently tied to one module. The default should be shared, and
exceptions should be explicit and explain why the shared venv is not suitable.

Decision state: accepted.

BLOCKER: no.

### D-004 - Keep The Venv Ignored And Recreatable

Question: Should virtual environments ever be committed?

Recommended default: no. Add `.codeheart/local/` to consumer `.gitignore` and treat everything
under it as local generated state.

Rationale: venvs contain platform-specific files, absolute paths, dependency copies, generated
entrypoints, and caches. They should be recreated from committed module/source instructions.

Decision state: accepted.

BLOCKER: no, unless the user wants a different source-control policy.

### D-005 - Keep First-Run Onboarding Lightweight

Question: Should the Operating Kit create `.codeheart/local/envs/python/` during default first-run
onboarding?

Recommended default: no. Create it on demand when a calling runbook needs Python tooling and the
user approves local setup.

Rationale: the existing first-run onboarding explicitly avoids creating a Python virtual
environment by default. On-demand creation keeps first-run setup light and avoids creating unused
runtime state.

Decision state: accepted.

BLOCKER: no.

### D-006 - Divide Generic And Module-Owned Responsibilities

Question: Which layer owns the venv commands?

Recommended default:

- Operating Kit owns the generic `.codeheart/local/` boundary, Python venv default path, approval
  pattern, readiness checks, and recheck/return behavior.
- Modules own package-specific install commands, package names, smoke commands, and module
  readiness evidence.

Rationale: this matches the current tooling-readiness design. Generic doctrine should not hardcode
AI Execution behavior, and AI Execution should not invent the shared venv convention.

Decision state: accepted.

BLOCKER: no.

### D-007 - Broaden The Trigger Beyond Module Onboarding

Question: Should the readiness route apply only to module onboarding/operation, or to any agent
interaction blocked by missing local tooling?

Recommended default: apply it to any agent interaction where missing local tooling blocks progress.

Rationale: users should not need to know whether a blocker came from a module runbook, a repo
task, a validation command, or a generic agent operation. The agent should classify local tooling
blockers consistently and return to the original task after readiness is resolved.

Decision state: accepted.

BLOCKER: no.

### D-008 - Keep Machine/User-Level Tools In Scope But Governed

Question: Should the standard exclude global/user-level tools?

Recommended default: no. Keep them in scope only for baseline tooling that is naturally outside
the repo, such as package managers, runtimes, browsers, shell runtimes, and official vendor CLIs.

Rationale: Homebrew itself cannot live in `.codeheart/local/`, and neither can many system
prerequisites. The standard should govern when and how an agent may propose these installs, while
keeping repo-specific packages inside repo-local environments.

Decision state: accepted.

BLOCKER: no.

### D-009 - Add Gitignore Behavior Without Default Directory Creation

Question: Should init and sync add `.codeheart/local/` to consumer `.gitignore`, and should they
create the directory during default setup?

Recommended default: add `.codeheart/local/` to `.gitignore` during init and sync, but do not
create `.codeheart/local/` by default.

Rationale: generated local machine/runtime state must be ignored before a runbook creates it, but
default onboarding should stay lightweight and should not scaffold unused runtime directories.

Decision state: accepted.

BLOCKER: no.

### D-010 - Add Optional Local Machine Layer Config

Question: Should `.codeheart/kit.config.yaml` expose the local machine layer path?

Recommended default: add optional `local_machine_layer_path: .codeheart/local/` under
`local_consumer_layer` for new installs. Existing configs remain valid when the field is absent.

Rationale: a machine-readable path helps agents and runbooks find the standard local runtime
boundary, while optional schema compatibility avoids forcing existing consumers through a config
migration before they need local runtime tooling.

Decision state: accepted.

BLOCKER: no.

## Resolved Questions

No implementation-planning blockers remain in this discovery.

| ID | Question | Accepted default | BLOCKER |
| --- | --- | --- | --- |
| Q-001 | Should the local machine layer be `.codeheart/local/`? | Yes. | no |
| Q-002 | Should the default Python venv be `.codeheart/local/envs/python/`? | Yes. | no |
| Q-003 | Should module-specific venvs be allowed? | Yes, only as explicit exceptions under `.codeheart/local/envs/<purpose>/`. | no |
| Q-004 | Should first-run onboarding create the Python venv? | No, create on demand. | no |
| Q-005 | Should AI Execution be patched immediately after this standard is accepted? | Yes, update install/onboard runbooks to use the shared venv path. | no |
| Q-006 | Should readiness routing apply to any blocked agent task, not only module onboarding? | Yes. | no |
| Q-007 | Should Homebrew and equivalent package-manager/bootstrap tools remain in scope? | Yes, as governed machine/user-level baseline tooling. | no |
| Q-008 | Should init/sync add `.codeheart/local/` to `.gitignore` without creating the directory? | Yes. | no |
| Q-009 | Should new installs expose an optional `local_machine_layer_path: .codeheart/local/` config field? | Yes, while old configs remain valid without it. | no |

## Likely Implementation Touchpoints

For implementation planning, consider:

- `src/codeheart_operating_kit/components.py`: add `.codeheart/local/` gitignore behavior for
  init/sync without scaffolding the directory, and add the optional new-install config field if
  config generation remains owned there;
- `schemas/kit-config.schema.json`: accept optional
  `local_consumer_layer.local_machine_layer_path` with `.codeheart/local/`;
- `tests/test_init.py`, `tests/test_sync_check.py`, `tests/test_json_schemas.py`, and
  `tests/test_packaging_resources.py`: cover gitignore, schema, generated config, packaging, and
  no-default-directory creation behavior;
- consumer `.gitignore` scaffold or guidance: add `.codeheart/local/` during init/sync;
- `components/structure-governance/managed/reference/managed-content-boundaries.md`: define the
  local machine layer;
- `components/agent-interface/managed/runbooks/handle-tooling-readiness.md`: add the Python venv
  convention to the `python-runtime` lane and broaden the trigger from module onboarding/operation
  to any blocked agent task;
- `components/agent-interface/managed/reference/local-extension-contract.md`: clarify the
  difference between `.codeheart/user/` and `.codeheart/local/`;
- `components/agent-interface/managed/reference/runbook-authoring-standard.md`: keep module-owned
  tool declarations separate from generic readiness routing;
- `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`: point generic
  local tooling blockers to the readiness route regardless of the original task route;
- `components/agent-interface/managed/runbooks/conduct-first-run-onboarding.md`: preserve the rule
  that default first-run onboarding does not create a venv;
- packaged-resource mirrors under `src/codeheart_operating_kit/resources/`;
- AI Execution install/onboard runbooks in the Foundry module, after the Operating Kit standard is
  accepted.

## Draft Implementation Capability Scope - Local Runtime Environment Standard

Capability:
Agents can route any missing local tooling blocker through a standard environment-readiness flow
and, when Python repo/module tooling is needed, create or reuse a standard ignored repo-local
Python venv without modifying global Python or storing generated runtime files in managed
snapshots.

Primary workflow:
A calling runbook or agent task detects missing local tooling, routes through Operating Kit
tooling readiness, classifies the blocker as machine/user-level baseline tooling or repo-local
tooling, receives user approval for local setup, rechecks readiness, then returns to the original
task. For Python repo-local tooling, the default path is to create or reuse
`.codeheart/local/envs/python/`, install the module-owned package, and validate the CLI.

Must cover:

- `.codeheart/local/` as ignored local machine/runtime state.
- `.codeheart/local/envs/python/` as the default Python venv.
- machine/user-level baseline tooling as a governed lane for package managers, runtimes, browsers,
  shell runtimes, and official vendor CLIs.
- `.codeheart/local/` added to consumer `.gitignore` during init/sync without creating the
  directory by default.
- optional `local_consumer_layer.local_machine_layer_path` for new installs while existing configs
  remain valid without it.
- user approval before creating or repairing local tooling.
- no venv creation during default first-run onboarding.
- module-owned package install and smoke validation.
- explicit exception path for dependency conflicts.
- generic fallback behavior when a missing local tool is not explicitly listed yet.

Explicitly out of scope:

- creating the venv during discovery;
- creating a Python dependency lock strategy;
- converting all module runbooks at once;
- ad hoc global Python or package installs;
- committing generated runtime state.

Deferred or blocked:

- implementation plan: ready after this accepted discovery.
- exact OS-specific install commands for package managers and runtimes: deferred to implementation
  and official-source checks.

Preserve decisions:

- D-001: `.codeheart/local/` is the local machine/runtime boundary.
- D-002: `.codeheart/local/envs/python/` is the default Python venv.
- D-003: purpose-specific venvs are exceptions, not defaults.
- D-004: venvs stay ignored and recreatable.
- D-005: first-run onboarding stays lightweight.
- D-006: Operating Kit owns the generic lane; modules own package-specific install and validation.
- D-007: readiness routing applies to any blocked agent task.
- D-008: machine/user-level tools stay in scope but governed.
- D-009: init/sync add `.codeheart/local/` to `.gitignore` without default directory creation.
- D-010: new installs expose optional `local_machine_layer_path` while old configs remain valid.

Planner must not reinvent:

- `.codeheart/user/` remains the local user preference/notes layer.
- `.codeheart/foundry/modules/<id>/` remains managed snapshot content and must not receive venvs,
  caches, or editable-install artifacts.
- live provider auth or external service readiness is not local tooling readiness.
- Homebrew or an equivalent package manager is not repo-local tooling; it is machine/user-level
  baseline tooling and requires explicit approval.

Feature-level success evidence:

- a fresh consumer repo can create or reuse `.codeheart/local/envs/python/`;
- init/sync make `.codeheart/local/` ignored without creating `.codeheart/local/` by default;
- schema validation accepts configs with and without `local_machine_layer_path`, and new install
  config includes the field;
- the AI Execution package can be installed into that venv;
- `foundry-ai --help` and dry-run smoke work from the venv;
- an unknown missing local tool blocker is classified, approval-gated, and either resolved or
  stopped with a clear blocker summary;
- no generated venv, cache, bytecode, or editable-install artifact is committed or placed in the
  managed module snapshot.
