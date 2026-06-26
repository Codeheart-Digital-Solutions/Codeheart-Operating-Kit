Last updated: 2026-06-26T16:18:05Z (UTC)
Created: 2026-06-26
Status: active

# Consumer Runtime Materialization And Terminal Handoff Hardening Execution Log

Plan path:
`docs/repo/plans/runtime-materialization-hardening/runtime-materialization-hardening_implementation_doc.md`

Mode: direct implementation after user-approved implementation planning.

## Summary

Implementation was activated from the consumer runtime materialization and terminal handoff
hardening plan. Scope includes Operating Kit doctrine and release work, Foundry AI Execution
adopter hardening, HQ managed snapshot proof, named Operating Kit installs, and per-epic
read-only review gates through subagents.

## Epic Delta Index

| Epic | Status | Notes |
| --- | --- | --- |
| EP-01 | completed | Doctrine, packaged mirrors, focused validation, fresh probe, and rerun review gate passed. |
| EP-02 | completed | Doctrine, packaged mirrors, focused validation, fresh probe, and rerun review gate passed. |
| EP-03 | active | Version bump, release notes, asset build, release manifest validation, and full source validation complete; release publication pending. |
| EP-04 | completed | Foundry AI Execution adopter patch, source tests, smoke validation, fix review, and rerun review gate passed. |
| EP-05 | pending | HQ AI Execution snapshot and runtime proof pending. |
| EP-06 | pending | Named Operating Kit installs pending. |
| EP-07 | active | Lifecycle activation and execution log created. |

## Activation Notes

- Activated canonical plan status from `draft` to `active` after the user requested implementation.
- Created this sibling execution log before source edits, following the managed implementation
  execution runbook.
- Review gates will use fresh read-only subagents after each epic where the environment permits
  subagent execution.

## Validation Log

- `python3 scripts/validate-markdown-headers.py`: passed.
- `python3 scripts/validate-public-core.py`: passed.
- `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest tests/test_packaging_resources.py tests/test_onboard.py tests/test_sync_check.py`:
  passed, 26 tests.
- `git diff --check` in Codeheart-Operating-Kit: passed.
- Managed source-to-packaged-resource parity was checked by `tests/test_packaging_resources.py`
  and direct `cmp` for changed managed files.

## Fresh Low-Context Probes

### EP-01 Runtime Materialization Probe

Status: passed.

Fresh low-context probe agent: `019f04ab-a96f-7e11-95c7-dc2780cbcc82`.

Scenario: consumer repo has a managed `ai-execution` snapshot and `foundry-ai` is unavailable in
the repo-local environment.

Result: the agent should not use editable install from the managed snapshot. It should route the
missing `foundry-ai` blocker through Operating Kit tooling readiness, materialize runtime tooling
into ignored repo-local state, normally `.codeheart/local/envs/python/`, and use a non-editable
package install from the snapshot path before smoke validation. AI Execution owns concrete package
facts such as package name, Python requirement, console script, and smoke commands.

Evidence cited by the probe:

- `docs/repo/reference/placement-contract.md` for `.codeheart/local/` as generated local
  runtime/tooling state;
- `components/agent-interface/managed/runbooks/handle-tooling-readiness.md` for consumer-mode
  runtime materialization, default Python venv, and non-editable package install shape;
- `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` for selecting
  the owner route before tooling readiness and using consumer-mode non-editable materialization.

### EP-02 Visible-Terminal Handoff Probe

Status: passed.

Fresh low-context probe agent: `019f04ab-c8cc-7420-9169-09677d853ebd`.

Scenario: user asks to start AI Execution auth setup so they can paste an OpenAI API key, while the
agent can run hidden shell tools that the user may not see or control.

Result: the agent should not run `foundry-ai auth setup` in an agent-hidden prompt and tell the
user to paste the key there. It should run read-only/noninteractive preflight first, then use
visible-terminal handoff: prepare command and working directory, use a user-visible terminal or
provide the exact command for the user to run, tell the user to paste the key only into that
terminal prompt and press Enter, and wait for completion before validating. Secret evidence must be
status-only and token-free.

Evidence cited by the probe:

- `components/agent-interface/managed/reference/runbook-authoring-standard.md` for visible
  terminal handoff and secret handling;
- the implementation plan's AI Execution onboarding criteria for package validation before auth
  setup and no hidden interactive prompts;
- `AGENTS.md` for the immediate secret-safety rule.

## Review Gate Log

- EP-01 initial read-only review agent `019f04aa-1a07-77c2-82b6-8b527424f3a8`: doctrine aligned;
  material finding was missing fresh low-context probe evidence. Probe evidence has since been
  recorded above.
- EP-02 initial read-only review agent `019f04aa-3e21-7cf0-89c1-566f3978d0e2`: doctrine aligned;
  material finding was missing fresh low-context probe evidence. Probe evidence has since been
  recorded above.
- EP-01 rerun read-only review agent `019f04ad-37c3-7f73-ab6f-05ab2248ff2c`: no material
  findings. Prior probe-evidence finding is resolved.
- EP-02 rerun read-only review agent `019f04ad-4f53-7b13-9a85-23dd13d881b7`: no material
  findings. Prior probe-evidence finding is resolved.

## EP-04 AI Execution Adopter Evidence

- Updated AI Execution version surfaces from `0.1.0` to `0.1.1` in `module.yaml`,
  `pyproject.toml`, package `__version__`, and CLI version test.
- Updated AI Execution README, install, onboarding, troubleshooting, and managed snapshot contract
  so consumer package validation uses consumer-mode non-editable install into an approved local
  runtime, builds wheels from an ignored local copy under `.codeheart/local/`, and keeps editable
  installs source-development mode only.
- Added package data for AI Execution schemas and default templates so non-editable wheel installs
  can run dry-run validation without relying on source-tree layout.
- Updated onboarding so API-key entry happens only after package/CLI validation and through
  visible-terminal handoff, not chat or an agent-hidden prompt.
- `uv run --no-project --with pytest --with pip --with setuptools --with wheel --with jsonschema --with keyring --with openai --with PyYAML python -m pytest tests`
  from `modules/ai-execution`: passed, 61 tests.
- `python3 -m pytest modules/ai-execution/tests`: failed before test execution because the visible
  Homebrew Python environment has no `pytest` module installed.
- Temporary consumer-style validation copied the managed snapshot to a temporary consumer repo,
  built a non-editable wheel from `.codeheart/local/build/ai-execution/source`, installed it into
  `.codeheart/local/envs/python/`, and ran module help, CLI help, version, dry-run smoke, and auth
  status.
- Temporary consumer-style validation passed and reported version `0.1.1`,
  `sends_provider_request: false`, and secret-safe auth status `ready` from `os-keychain`.
- The same temporary validation checked both the managed snapshot copy and source module tree for
  `uv.lock`, `*.egg-info`, `build`, `dist`, `__pycache__`, and `*.pyc`; no generated metadata was
  present in the managed snapshot copy, and the source tree remained free of `uv.lock`, `*.egg-info`,
  `build`, and `dist`.
- Source grep for `pip install -e`, approved editable-install phrasing, managed-snapshot editable
  install phrasing, and `0.1.0`: found only the explicit guardrail `Do not use pip install -e from
  a managed consumer snapshot`; no consumer instruction to use editable install remains.
- `git diff --check` in Codeheart-Automation-Foundry and Codeheart-HQ: passed.
- EP-04 initial read-only review agent `019f04b0-7c54-70a3-9d37-8de56cde1260`: found that
  `modules/ai-execution/uv.lock` generated by validation could leak editable source metadata into
  the managed snapshot because snapshot export did not exclude `uv.lock`. Fix applied: snapshot
  export excludes `uv.lock`, managed contract names source-development lock metadata, validation
  uses lock-free commands, and consumer wheel builds from an ignored local copy under
  `.codeheart/local/`.
- EP-04 rerun read-only review agent `019f04b5-e40b-7740-b0a0-d0581afb9144`: no material
  findings. Prior `uv.lock` and editable metadata finding is resolved.

## EP-03 Operating Kit Release Prep Evidence

- Version surfaces were bumped to Operating Kit `0.1.15` in package metadata, installers,
  bootstrap guidance, release asset builder defaults, source manifests, packaged release metadata,
  and changed component manifests.
- Release notes include a `v0.1.15` entry for consumer runtime materialization, the Python
  runtime lane, visible-terminal handoff, generated metadata boundaries, and planning/resource
  packaging impact.
- Built release assets:
  - `codeheart-operating-kit-0.1.15-macos.tar.gz`
  - `codeheart-operating-kit-0.1.15-macos.tar.gz.sha256`
  - `codeheart-operating-kit-0.1.15-windows.zip`
  - `codeheart-operating-kit-0.1.15-windows.zip.sha256`
- Source release manifest asset hashes were updated after build and `python3
  scripts/validate-release-manifest.py` passed.
- Stale release-version scan found `0.1.14` only in historical release notes.
- Operating Kit validation passed:
  - `python3 scripts/validate-json-schemas.py`
  - `python3 scripts/validate-markdown-headers.py`
  - `python3 scripts/validate-public-core.py`
  - `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest tests`:
    passed, 92 tests.
  - `git diff --check`: passed.
