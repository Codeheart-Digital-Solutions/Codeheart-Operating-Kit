Last updated: 2026-06-26T16:33:33Z (UTC)
Created: 2026-06-26
Status: completed

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
| EP-03 | completed | Operating Kit `v0.1.15` source commit, tag, GitHub release, assets, installer checks, and published macOS install proof complete. |
| EP-04 | completed | Foundry AI Execution adopter patch, source tests, smoke validation, fix review, and rerun review gate passed. |
| EP-05 | completed | Foundry AI Execution source commit, HQ snapshot, lock update, repo-local venv install, dry-run smoke, auth status, and snapshot hygiene proof complete. |
| EP-06 | completed | Released Operating Kit `0.1.15` synced/initialized into HQ, Foundry, the named private platform repository, and Operating Kit; `check` passed in all four repos. |
| EP-07 | completed | Registers, work board, final validation, review gate, reviewer finding resolution, and completion status updates complete. |

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
- Final cross-repo read-only review agent `019f04c3-27a4-7943-8001-707a8c001e3a`: found two
  material issues, both resolved before completion:
  - P1: HQ, Foundry, and the named private platform repository lockfiles had `kit_version:
    0.1.15` but stale `release.asset_url` metadata pointing at the `0.1.14` macOS release asset.
    Fix: repaired the lockfile release blocks to point at the `0.1.15` macOS asset and checksum,
    then reran `codeheart-operating-kit check` and verified the lock URLs.
  - P2: HQ coordination register and portfolio work board disagreed on
    `CODEHEART-OPERATING-KIT-PR-017` status. Fix: updated HQ register and work board, the
    Operating Kit register, and the Foundry pointer to completed after validation passed.
  Residual risk noted by the reviewer remains: native Windows install was not run in this session.
- Final focused rerun review agent `019f04c9-858d-7660-a053-4f0bd1e6f35c`: verified the prior
  runtime-hardening lockfile/status fixes, then found an adjacent HQ planning-surface inconsistency
  for `CODEHEART-HQ-PR-002`. Fix: aligned the HQ Portfolio Work Board's two `CODEHEART-HQ-PR-002`
  status mentions and the canonical Portfolio Planning Surfaces discovery header to completed,
  matching the HQ plan register.

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
- Release source commit: `3a4332078d1b8d36b5fc50c2886668ae0c22896a`.
- Release tag: `v0.1.15`.
- Release URL:
  `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v0.1.15`.
- Release asset hashes:
  - `bootstrap.md`:
    `1e50c5099427edd0c85c4f386c9e51df755de2cb2d0836ced913e726882a97b1`
  - `install.sh`:
    `14bb968a3b74c29dbd7596b130cdadbcbdc4eb9b5ceb649e7773ba1f55c17bac`
  - `install.ps1`:
    `22796b95ebece6f4ef135eb014081c03bca17fcdf105d03456a9bdeb8fd03ea1`
  - `release-notes.md`:
    `6f7103d1b17c125b51851a6194b3b13246f0baa2ca03e73deb62b6c1ba2f7142`
  - `codeheart-operating-kit-0.1.15-macos.tar.gz`:
    `ceadfb5f81725b0575ff4a883af41e8bb1bffd9e2b77aa6f4f61497283853ea5`
  - `codeheart-operating-kit-0.1.15-macos.tar.gz.sha256`:
    `c4d2473b2d973073e7c48e8573a3c2fd022b0ae197a5ccdcac1733bbb08c6230`
  - `codeheart-operating-kit-0.1.15-windows.zip`:
    `4baaeed6ae7d0e71caa3cc398b952550d807efd2d958329dfbf5505d277a87e6`
  - `codeheart-operating-kit-0.1.15-windows.zip.sha256`:
    `50822de74e51b1a0bc7933d91c0540418671a8327aad875794ed829d4d6ca1fa`
- `install.sh` fail-closed checksum mismatch check passed.
- `install.ps1` fail-closed checksum mismatch check passed under local `pwsh`.
- Temporary macOS install from local `dist/` asset passed and reported
  `codeheart-operating-kit 0.1.15`.
- Temporary PowerShell install from local Windows asset passed under local `pwsh`.
- Published-release macOS install using the default GitHub release URL passed and reported
  `codeheart-operating-kit 0.1.15`.
- Residual Windows release risk: the Windows installer was validated through PowerShell on macOS,
  but not through a native Windows runner in this session.

## EP-05 AI Execution Release And HQ Consumer Proof Evidence

- Final Foundry AI Execution source commit:
  `54072cff5023dd69ce872ee6120ee2e23b7db435`.
- Foundry source tree was clean when the final HQ snapshot source commit was recorded.
- Foundry AI Execution `0.1.1` release source was pushed to
  `Codeheart-Digital-Solutions/Codeheart-Automation-Foundry`.
- During HQ consumer install proof, the first local wheel build failed because ambient pip
  configuration pointed build isolation at a private CodeArtifact index. Fix applied in the
  Foundry source and HQ snapshot: consumer wheel build now uses `PIP_CONFIG_FILE=/dev/null` and
  `--no-build-isolation` for the no-dependency local wheel build.
- HQ managed snapshot was refreshed from the final Foundry source commit with cache/build metadata
  exclusions including `__pycache__/`, `.pytest_cache/`, `*.pyc`, `*.egg-info/`, `.venv/`,
  `uv.lock`, and `tmp/`.
- HQ AI Execution snapshot manifest contains 70 files and has manifest hash
  `988eb2c976320d1ff72f8937ecba35c99b5456d9d5c13aa88db772cade22b1f4`.
- HQ `.codeheart/foundry/foundry.lock.yaml` now records AI Execution `0.1.1`, source plan
  `CODEHEART-AUTOMATION-FOUNDRY-PR-008`, source tree state `clean`, source commit
  `54072cff5023dd69ce872ee6120ee2e23b7db435`, and the new snapshot manifest hash.
- HQ repo-local venv at `.codeheart/local/envs/python/` was reused.
- Consumer-mode non-editable install proof built a wheel from
  `.codeheart/local/build/ai-execution/source/`, installed
  `foundry_ai_execution-0.1.1-py3-none-any.whl`, and `foundry-ai version` reported `0.1.1`.
- `foundry-ai --help` and `python -m foundry_ai_execution --help` passed from the HQ venv.
- `foundry-ai run --dry-run --manifest
  .codeheart/foundry/modules/ai-execution/templates/smoke-input-manifest.yaml` reported
  `sends_provider_request: false`.
- `foundry-ai auth status --json` reported secret-safe status `ready` from `os-keychain`. The
  expected `missing` pre-key state did not apply because the OS keychain already has credentials.
- Auth setup was not started; no secret prompt was opened and no key was entered in chat or an
  agent-hidden terminal.
- Post-install snapshot hygiene check found zero generated artifacts in the managed snapshot.

## EP-06 Named Operating Kit Install Evidence

- User-level `codeheart-operating-kit` was updated from `0.1.14` to `0.1.15` using the published
  release installer.
- `codeheart-operating-kit sync . --json` was run in:
  - Codeheart-HQ;
  - Codeheart-Automation-Foundry;
  - named private platform repository.
- `codeheart-operating-kit init . --project-name Codeheart-Operating-Kit --json` was run in
  Codeheart-Operating-Kit because the self-install state was absent after release-staging cleanup.
- `codeheart-operating-kit check . --json` passed in all four repos with `ok: true`, no drift, no
  missing route targets, no missing routing, and `stale_cli: false`.
- After final review, HQ, Foundry, and the named private platform repository lockfiles were
  repaired so their `release.asset_url` and `checksum_sha256` point at
  `codeheart-operating-kit-0.1.15-macos.tar.gz`:
  `ceadfb5f81725b0575ff4a883af41e8bb1bffd9e2b77aa6f4f61497283853ea5`.
- `codeheart-operating-kit check . --json` was rerun in HQ, Foundry, the named private platform
  repository, and Operating Kit after the lockfile release metadata repair; all four passed.
- Installed `AGENTS.md` managed blocks in synced repos include the missing-local-tooling route to
  the managed tooling-readiness runbook.
- Installed tooling-readiness runbooks in synced repos include the runtime materialization
  hardening and visible-terminal handoff route through the `0.1.15` managed content.
- `.codeheart/local/` ignore/config route was preserved or added by sync/init:
  - HQ: already present; preserved.
  - Foundry: `.gitignore` changed during sync.
  - named private platform repository: already present; preserved.
  - Operating Kit: initialized with `.codeheart/local/` ignored.
- Dirty-work preservation notes:
  - HQ had existing M365, Operating Kit sync, planning/work-board, and module-state changes; they
    were preserved while AI Execution and Operating Kit sync changes were added.
  - Foundry and the named private platform repository had existing Operating Kit sync-state
    changes; sync updated them to `0.1.15` without unrelated source edits.
  - Operating Kit release source commit stayed tagged at `v0.1.15`; subsequent self-install state
    and execution-log updates remain post-release working-tree changes.

## EP-07 Final Register, Work Board, And Validation Evidence

- Operating Kit plan register marks `OK-PR-017` completed and relates `OK-PR-016` to `OK-PR-017`.
- HQ coordination register marks `CODEHEART-OPERATING-KIT-PR-017` completed.
- HQ work board marks `CODEHEART-OPERATING-KIT-PR-017` completed under Operating Kit And Foundry
  System Model.
- Foundry plan register marks `CODEHEART-AUTOMATION-FOUNDRY-PR-008` completed as the local
  AI Execution adopter/release pointer.
- Final `git diff --check` passed in HQ, Foundry, the named private platform repository, and
  Operating Kit.
- Final Operating Kit validations passed:
  - `python3 scripts/validate-markdown-headers.py`
  - `python3 scripts/validate-public-core.py`
  - `python3 scripts/validate-json-schemas.py`
  - `python3 scripts/validate-release-manifest.py`
  - `uv run --with pytest --with pip --with setuptools --with wheel python -m pytest tests`:
    passed, 92 tests.
- Generated `uv.lock` from the final `uv run` validation was removed from the public Operating Kit
  working tree before completion.
