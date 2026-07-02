Last updated: 2026-07-02T13:16:41Z (UTC)
Created: 2026-06-29
Status: release-candidate

# Document Header

## Overview

This implementation plan turns the approved repo feedback capture discovery into a shippable
Operating Kit v1. It adds managed agent guidance for capturing repo-specific feedback as a
Codeheart GitHub organization member feature, a check-first GitHub Issues flow behind an existing
authenticated `gh` gate, demand-driven maintainer setup guidance, optional repo-feedback config
schema support, route visibility for fresh installs, packaged-resource mirrors, validation, and
release notes.

The first release deliberately avoids a new CLI command and avoids automatic GitHub changes.
Agents may run read-only `gh` authorization and issue preflight checks, but only existing
authenticated `gh` is allowed. Missing `gh`, missing authentication, or unverifiable
`Codeheart-Digital-Solutions` organization membership makes repo feedback capture unavailable
without user-facing fallback. Issue creation, label/template changes, and repository setting
changes remain explicit user-approved actions for verified Codeheart maintainers.

## Essential Context Reference Files

| Path | Reason |
| --- | --- |
| `AGENTS.md` | Maintainer bootstrap, public-core safety, and change-safety rules. |
| `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_discovery_doc.md` | Approved capability scope, decisions, non-goals, and handoff requirements. |
| `docs/repo/plans/kit-feedback-intake/kit-feedback-intake_implementation_doc.md` | Precedent for Operating Kit feedback intake, issue forms, labels, managed guidance, and triage. |
| `docs/repo/plans/kit-feedback-intake/kit-feedback-intake_execution_log.md` | Evidence that GitHub Issues were checked first, already enabled, and labels/forms were handled as governance. |
| `docs/repo/runbooks/release-operating-kit.md` | Public release procedure, asset, checksum, tag, GitHub release, and evidence requirements. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Required quality bar for the new durable runbooks. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Route-before-surface doctrine for the new feedback route and fresh-agent probe. |
| `components/agent-interface/managed/reference/root-agents-md-contract.md` | Contract for direct managed routes that may appear in the installed root `AGENTS.md` managed block. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Authoring source for installed `.codeheart/kit/docs/agent-interface/reference/operational-recipe-maturity.md`; required maturity reference for repeated capture/setup recipes and non-promotion boundary. |
| `components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md` | Authoring source for installed `.codeheart/kit/docs/agent-interface/reference/runbook-to-script-promotion-standard.md`; required reference for consciously keeping v1 runbook-only and not creating reusable script assets. |
| `components/agent-interface/managed/runbooks/submit-kit-feedback.md` | Existing route for feedback about the Operating Kit itself; repo feedback must not replace it. |
| `components/agent-interface/managed/reference/kit-feedback-item-format.md` | Existing feedback item format precedent. |
| `docs/repo/reference/consumer-impact-classification.md` | Canonical release and manifest impact class vocabulary. |
| `components/agent-interface/component.yaml` | Managed agent-interface component manifest that must include new managed docs. |
| `templates/agents/AGENTS.managed-block.md` | Installed root route surface for fresh repos. |
| `components/agent-interface/managed/README.md` | Managed agent-interface route index. |
| `components/agent-interface/managed/kit-readme.md` | Installed fallback inventory for managed routes. |
| `schemas/kit-config.schema.json` | Optional `repo_feedback` config contract. |
| `src/codeheart_operating_kit/components.py` | Fresh install config writer; v1 must keep `repo_feedback` absent by default. |
| `src/codeheart_operating_kit/__init__.py` | Package runtime version used by the CLI and release asset validation. |
| `pyproject.toml` | Python package version consumed by the release asset builder. |
| `tests/test_json_schemas.py` | Focused schema tests for optional `repo_feedback` states. |
| `tests/test_init.py` | Fresh install tests for no default `repo_feedback` block and route target presence. |
| `tests/test_sync_check.py` | Sync/check route target validation. |
| `tests/test_packaging_resources.py` | Packaged resource fallback coverage. |
| `release-notes.md` | Release-facing consumer-impact notes. |
| `scripts/build-release-assets.py` | Release asset builder and checksum generator. |

## Table Of Contents

- [Section 1 - Foundation](#section-1---foundation)
- [Section 2 - Strategy](#section-2---strategy)
- [Section 3 - Execution Plan](#section-3---execution-plan)
- [Section 4 - Future Planning](#section-4---future-planning)
- [Revision Notes](#revision-notes)

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Ship a managed Operating Kit repo feedback capture v1 that works in fresh and existing consumer
repositories.

Completion is proven when:

- root `AGENTS.md` managed routes expose repo feedback capture separately from Operating Kit
  feedback;
- installed managed docs include a repo feedback capture runbook, GitHub issue-intake setup
  runbook, and repo feedback item-format reference;
- the capture runbook tells agents to silently verify existing `gh`, GitHub CLI authentication,
  and active `Codeheart-Digital-Solutions` organization membership before any feedback prompt;
- missing `repo_feedback` config means authorization-gated auto-check, not disabled;
- GitHub Issues capture works without standard labels/templates by putting classification in the
  issue body and using only labels that exist;
- sensitive, private, and security-disclosure material stops before GitHub issue creation and is
  not stored in feedback drafts;
- missing `gh`, missing authentication, or unverifiable organization membership stops silently and
  does not route to tooling readiness, browser fallback, manual fallback, or local drafts;
- setup is offered only to verified Codeheart maintainers when Issues are unavailable, disabled,
  or the maintainer wants missing standard labels/templates;
- decline can be recorded as disabled/suppressed state after explicit user approval;
- `schemas/kit-config.schema.json` accepts optional configured and disabled states, including the
  Codeheart authorization policy, while preserving old config validity;
- fresh install tests prove no `repo_feedback` block is written by default;
- packaged-resource fallback installs the new runbooks and reference files;
- validation covers Markdown headers, public-core hygiene, schema behavior, route targets,
  packaged resources, release manifests, focused tests, full tests, and a fresh low-context
  routing probe;
- repeated capture/setup recipes record L1 structured-recipe maturity, validation tier, evidence
  shape, blocker shape, and no-promotion evidence;
- v1 creates no reusable script asset, CLI command, wrapper, API, or durable executable helper;
- release notes record instruction, schema, and safety-policy impact;
- the Operating Kit release is published from the validated commit;
- approved consumer repositories are synced to the released version or have explicit blockers
  recorded.

## 1.2 Project And Problem Context

The Operating Kit already has a public feedback intake path for feedback about the kit itself.
That path sends sanitized feedback to the public `Codeheart-Operating-Kit` GitHub Issues backlog
and provides a maintainer triage workflow.

Repo-specific feedback is different. A runbook gap, script failure, conservative undocumented
override, docs conflict, or user dissatisfaction usually belongs in the repository that owns the
affected artifact. Without a managed route, agents can solve the immediate issue in chat and lose
the reusable maintenance signal.

The approved direction is to make the Operating Kit own the generic capture route while each
repository owns its feedback destination and issue triage. The feature is internal to Codeheart
maintainers and must not burden non-technical users. A fresh repo should not need feedback setup
during onboarding. When a real feedback item appears, the agent first verifies existing `gh`,
GitHub authentication, and Codeheart organization membership. Only then does it check whether
GitHub Issues already works, draft an issue when it does, and offer setup only when needed.

## 1.3 Current State Analysis

Existing state:

- `components/agent-interface/managed/runbooks/submit-kit-feedback.md` handles feedback about the
  Operating Kit itself.
- There is no managed route for repo-specific feedback capture.
- `templates/agents/AGENTS.managed-block.md` exposes Operating Kit feedback but not repo
  feedback.
- `schemas/kit-config.schema.json` allows `portfolio` and local layer config, but no
  `repo_feedback` block.
- `src/codeheart_operating_kit/components.py` writes fresh `.codeheart/kit.config.yaml` without
  `repo_feedback`.
- GitHub issue forms and labels exist for the Operating Kit repository, but consumer repos may
  have no standard labels/templates.

Target state:

- Missing `repo_feedback` config remains valid and means authorization-gated auto-check on first
  feedback trigger.
- Optional `repo_feedback` config records durable preferences after first use, setup, or verified
  maintainer decline.
- Managed runbooks describe the exact `gh`/organization gate, check-first, draft, approval, setup,
  no-fallback, and suppression flow.
- Fresh installs receive the route and docs through normal managed agent-interface sync.
- Standard labels/templates are useful setup, not prerequisites for creating a feedback issue.

Consumer impact record:

- `instruction-only change` for managed routes, runbooks, references, and docs.
- `validator-only change` for optional kit-config schema coverage.
- `security or safety policy change` for external issue creation, GitHub settings, labels,
  templates, and sanitization approval rules.
- No consumer migration required.
- No new scaffold path required.
- No new CLI behavior in v1.
- No reusable script asset in v1.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

```text
Codeheart-Operating-Kit/
  components/
    agent-interface/
      component.yaml                                      # modify
      managed/
        README.md                                         # modify
        kit-readme.md                                     # modify
        reference/
          root-agents-md-contract.md                      # modify
          repo-feedback-item-format.md                    # create
        runbooks/
          capture-repo-feedback.md                        # create
          enable-github-issues-feedback-intake.md          # create
  templates/
    agents/
      AGENTS.managed-block.md                             # modify
  pyproject.toml                                          # modify package version
  src/
    codeheart_operating_kit/
      __init__.py                                         # modify runtime version
      resources/
        components/agent-interface/...                    # mirror created and modified files
        templates/agents/AGENTS.managed-block.md          # mirror modified template
        manifest.yaml                                     # modify checksum and impact
  schemas/
    kit-config.schema.json                                # modify
  tests/
    test_json_schemas.py                                  # modify
    test_init.py                                          # modify
    test_sync_check.py                                    # modify
    test_packaging_resources.py                           # modify
  docs/
    README.md                                             # modify
    repo/
      README.md                                           # modify
      plans/
        README.md                                         # modify
        plan-register.md                                  # modify
        repo-feedback-capture/
          repo-feedback-capture_discovery_doc.md          # modify status and handoff
          repo-feedback-capture_implementation_doc.md     # this plan
          repo-feedback-capture_execution_log.md          # create at activation
manifest.yaml                                           # modify release impact
release-notes.md                                        # modify
pyproject.toml                                          # modify package version
```

## 2.2 Open Questions And Assumptions Requiring Clarification

`OQ-1`: Exact release version.

- `BLOCKER: no`
- Affects: `E5`, `E6`, `E7`
- Unlocks: package version bump, release notes, release manifest wording, release asset names,
  tag, and consumer sync target.
- Recommended default: use the next patch release available at execution time and record the
  exact version in the execution log.

`OQ-2`: GitHub CLI availability in consumer repos.

- `BLOCKER: no`
- Affects: `E3`, `E4`, `E6`
- Unlocks: preflight and setup command lane.
- Recommended default: use existing authenticated `gh` as the only command lane because it
  supports organization membership verification, `gh repo view --json hasIssuesEnabled`,
  `gh repo edit --enable-issues`, `gh label list`, `gh label create`, and `gh issue create`.
  Missing local `gh`, missing auth, or unverifiable `Codeheart-Digital-Solutions` membership is
  silent unavailability. Do not route missing `gh` through tooling readiness and do not provide
  browser, manual, or local-draft fallback.

`OQ-3`: Standard issue template installation in consumer repos.

- `BLOCKER: no`
- Affects: `E4`
- Unlocks: setup runbook wording.
- Recommended default: setup runbook may add a repo-owned issue template only after explicit user
  approval; missing templates never block issue capture.

## 2.3 Architectural Decisions With Reasoning

`AD-1`: Use the existing `agent-interface` component.

1. Problem being solved: installed agents need the route without adding another component.
2. Simplest working solution: add the capture/setup runbooks and item-format reference under
   `components/agent-interface/managed/`.
3. What may change in 6-12 months: a dedicated feedback component may be justified after CLI
   support, schemas, or wider repo onboarding exists.
4. Rationale: this keeps v1 as managed agent guidance and schema support, not a component
   addition.
5. Alternatives considered: a new `repo-feedback` component was rejected because it adds profile
   and release complexity before behavior is proven.

`AD-2`: Missing `repo_feedback` config means authorization-gated auto-check.

1. Problem being solved: fresh repos must work without onboarding-time feedback setup.
2. Simplest working solution: keep `repo_feedback` optional; absence means the capture runbook
   first verifies existing `gh`, GitHub authentication, and active Codeheart organization
   membership, then checks the GitHub remote and live issue availability when a real feedback item
   appears.
3. What may change in 6-12 months: onboarding may offer feedback setup after usage data proves it
   is worth front-loading.
4. Rationale: this preserves lightweight onboarding and matches the Operating Kit feedback
   precedent where Issues were already enabled.
5. Alternatives considered: writing `repo_feedback.mode: auto_check` during init was rejected
   because it creates config churn without adding capability.

`AD-3`: Schema records preferences and authorization policy, not behavior automation.

1. Problem being solved: agents need durable state for configured destinations, Codeheart-only
   authorization policy, and declined setup.
2. Simplest working solution: add an optional `repo_feedback` schema with `github_issues`,
   `disabled`, and authorization-policy fields requiring existing `gh`, verified
   `Codeheart-Digital-Solutions` membership, and silent unavailability.
3. What may change in 6-12 months: a CLI command may write or validate the block directly.
4. Rationale: v1 can rely on agent-runbook edits with approval while still validating the shape.
5. Alternatives considered: no schema was rejected because disabled/suppressed state would be
   inconsistent across repos.

`AD-4`: Missing labels/templates do not block capture.

1. Problem being solved: many repos will already have Issues enabled but not Codeheart feedback
   labels/templates.
2. Simplest working solution: draft body-classified issues and pass only labels proven to exist.
3. What may change in 6-12 months: label/template setup may become standardized across selected
   repos.
4. Rationale: feedback capture should work immediately in a fresh repo when Issues are enabled.
5. Alternatives considered: requiring standard labels first was rejected because it turns a
   working issue surface into a setup blocker.

`AD-5`: Use existing authenticated `gh` as the only command lane.

1. Problem being solved: agents need exact, testable preflight and setup commands without
   prompting non-technical users into GitHub tooling setup.
2. Simplest working solution: document `gh repo view`, `gh issue create`, `gh repo edit`, and
   `gh label` commands, plus `gh auth status` and `gh api user/memberships/orgs/Codeheart-Digital-Solutions`,
   with approval gates before write commands.
3. What may change in 6-12 months: a connector or CLI wrapper may provide a cleaner interface.
4. Rationale: `gh` is already available in the maintainer environment and is the direct GitHub
   governance lane. Missing `gh` means repo feedback capture is unavailable; it is not a tooling
   readiness blocker.
5. Alternatives considered: a new Operating Kit CLI command was rejected for v1 because the
   maintainer-run route needs real usage first. Browser/manual fallback was rejected because it
   would surface feedback administration to non-technical users.

`AD-6`: Treat setup as repository governance.

1. Problem being solved: enabling Issues, labels, and templates changes shared repository state.
2. Simplest working solution: the setup runbook first verifies Codeheart organization membership,
   preflights, explains the change, asks approval, executes the approved change, then verifies it.
3. What may change in 6-12 months: repo-wide setup could be batched by a separate governance
   plan.
4. Rationale: this matches the kit-feedback-intake precedent and avoids silent external changes.
5. Alternatives considered: silent setup was rejected because repository settings and labels are
   shared governance surfaces.

`AD-7`: Keep repo feedback separate from Operating Kit feedback.

1. Problem being solved: agents may otherwise route all feedback to the public Operating Kit repo.
2. Simplest working solution: managed docs state the split and route kit feedback to
   `submit-kit-feedback.md`.
3. What may change in 6-12 months: cross-repo triage may promote repeated repo issues into
   Operating Kit feedback.
4. Rationale: the owning repo is the right first inbox for repo-specific runbook and script
   friction.
5. Alternatives considered: one universal feedback inbox was rejected because it would create
   noise and privacy risk.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic ID | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| `E1` | Activation, baseline evidence, and handoff state are recorded. | S | None |
| `E2` | Optional `repo_feedback` config schema and fresh-install semantics are implemented and tested. | M | `E1` |
| `E3` | Managed capture route, runbook, and item-format reference are installed for fresh repos. | M | `E2` |
| `E4` | GitHub issue-intake setup runbook covers Issues, labels, templates, approval, and suppression. | M | `E2`, `E3` |
| `E5` | Packaged resources, component manifests, release surfaces, and indexes are synchronized. | M | `E3`, `E4` |
| `E6` | Validation, fresh-repo proof, routing probe, review, and release-candidate evidence are complete. | M | `E5` |
| `E7` | Operating Kit release is published and approved consumer repositories are synced or blocked with evidence. | M | `E6` |

## E1 - Activation, Baseline Evidence, And Handoff State

### A) Epic ID, Title, And Outcome

`E1` - Activation, Baseline Evidence, And Handoff State.

Outcome: the approved discovery handoff, implementation activation, execution log, and baseline
repository state are recorded before source implementation begins.

### B) Scope

Activate this plan after explicit execution approval. Create the sibling execution log, record the
dirty worktree baseline, preserve unrelated local changes, and refresh discovery/register
lifecycle state for the implementation path.

### C) Files Touched

```text
docs/repo/plans/repo-feedback-capture/
  repo-feedback-capture_discovery_doc.md          # modify lifecycle/handoff evidence
  repo-feedback-capture_implementation_doc.md     # modify status at activation
  repo-feedback-capture_execution_log.md          # create
docs/repo/plans/plan-register.md                  # modify lifecycle and relations
```

### D) Acceptance Criteria And Size

- Size: `S`
- Implementation plan status is `active` only after explicit execution approval.
- Execution log exists before source changes begin.
- Execution log records `git status --short` and notes unrelated local changes.
- Discovery handoff points to the approved implementation capability scope.
- Plan register links `OK-PR-020` and the implementation plan entry.

### E) Dependencies And Critical-Path Notes

No dependencies. `E2` through `E6` wait for activation evidence.

### F) Tasks Checklist

- [ ] Update `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_implementation_doc.md` status from `draft` to `active` after explicit execution approval.
- [ ] Create `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_execution_log.md` with activation context, epic table, validation sections, and review-gate placeholders.
- [ ] Record `git status --short` in the execution log with unrelated local changes preserved.
- [ ] Update `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_discovery_doc.md` with activation reference and implementation-plan pointer.
- [ ] Update `docs/repo/plans/plan-register.md` with active implementation lifecycle and relation from `OK-PR-020`.
- [ ] Run `git diff --check -- docs/repo/plans/repo-feedback-capture docs/repo/plans/plan-register.md`.

### G) Implementation Notes

Do not clean, revert, or restage unrelated existing changes. Treat existing dirty files as user
work unless the execution plan explicitly owns them.

### H) Open Questions

- None.

## E2 - Config Schema And Fresh-Install Semantics

### A) Epic ID, Title, And Outcome

`E2` - Config Schema And Fresh-Install Semantics.

Outcome: `.codeheart/kit.config.yaml` accepts optional repo feedback preferences and Codeheart
authorization policy, while fresh installs still omit `repo_feedback` and therefore use
authorization-gated auto-check behavior on first feedback trigger.

### B) Scope

Add optional `repo_feedback` schema support and tests. Do not change fresh install output to write
`repo_feedback` by default.

### C) Files Touched

```text
schemas/
  kit-config.schema.json                         # modify
tests/
  test_json_schemas.py                           # modify
  test_init.py                                   # modify
  fixtures/
    kit-config.yaml                              # preserve unless schema fixture coverage needs update
```

### D) Acceptance Criteria And Size

- Size: `M`
- Existing config files without `repo_feedback` remain valid.
- Fresh installs do not write a `repo_feedback` block.
- Valid `github_issues` config includes destination owner, repo, and Codeheart authorization
  policy.
- Valid `disabled` config records decline or owner-policy reason, and `mode: disabled` itself is
  the suppression state.
- Invalid partial and mode-incompatible states fail schema tests.

### E) Dependencies And Critical-Path Notes

Depends on `E1`. This epic establishes config semantics that runbooks reference.

### F) Tasks Checklist

- [ ] Add optional top-level `repo_feedback` to `schemas/kit-config.schema.json`.
- [ ] Add `repo_feedback.mode` enum values `github_issues` and `disabled` to `schemas/kit-config.schema.json`.
- [ ] Add `repo_feedback.destination.type` const `github_issues` plus `owner` and `repo` string fields for `github_issues` mode.
- [ ] Add `repo_feedback.authorization.organization` const `Codeheart-Digital-Solutions` for `github_issues` mode.
- [ ] Add `repo_feedback.authorization.require_verified_membership` const `true` for `github_issues` mode.
- [ ] Add `repo_feedback.authorization.require_gh_cli` const `true` for `github_issues` mode.
- [ ] Add `repo_feedback.authorization.unavailable_behavior` const `silent` for `github_issues` mode.
- [ ] Add `repo_feedback.github_standardization.labels` and `repo_feedback.github_standardization.issue_templates` enum values `configured`, `not_configured`, and `declined`.
- [ ] Add `repo_feedback.disabled_reason` enum values `verified_maintainer_declined_issue_intake`, `repo_owner_policy`, `issues_unavailable`, and `other` for `disabled` mode.
- [ ] Do not add a separate `suppress_prompts` field; `repo_feedback.mode: disabled` implies suppression.
- [ ] Keep `src/codeheart_operating_kit/components.py` fresh install config output unchanged for default `repo_feedback` absence.
- [ ] Add `tests/test_json_schemas.py` coverage for missing `repo_feedback` config validity.
- [ ] Add `tests/test_json_schemas.py` coverage for valid `github_issues` and `disabled` config blocks.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `github_issues` without `destination.type`.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `github_issues` without destination owner and repo.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `github_issues` without the Codeheart authorization policy block.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `github_issues` with any organization other than `Codeheart-Digital-Solutions`.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `github_issues` when `require_verified_membership`, `require_gh_cli`, or `unavailable_behavior: silent` are missing or false.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `disabled` without `disabled_reason`.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting invalid `github_standardization.labels` and `github_standardization.issue_templates` enum values.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting mode-incompatible fields across `github_issues` and `disabled`.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting unknown `repo_feedback.mode` values.
- [ ] Add `tests/test_init.py` assertion that new installs do not write a `repo_feedback` block.
- [ ] Run `python3 scripts/validate-json-schemas.py schemas/kit-config.schema.json`.
- [ ] Run `PYTHONPATH=src python3 -m pytest tests/test_json_schemas.py tests/test_init.py -q`.

### G) Implementation Notes

The existing lightweight schema test helper supports object, string, integer, const, enum,
required, additionalProperties, and allOf patterns. Keep the schema shape compatible with that
helper unless the implementation also extends the helper and tests.

Mode-incompatible field rejection is part of the intended schema behavior. If the current helper
cannot express or test the final schema shape cleanly, extend the helper in the same epic and
cover the helper extension with failing and passing examples.

### H) Open Questions

- `OQ-1` affects release wording only and does not block this epic.

## E3 - Managed Capture Route, Runbook, And Item Format

### A) Epic ID, Title, And Outcome

`E3` - Managed Capture Route, Runbook, And Item Format.

Outcome: fresh and synced repositories expose a managed route for repo-specific feedback capture,
and the capture runbook gives agents the full check-first workflow without confusing repo
feedback with Operating Kit feedback.

### B) Scope

Add the managed repo feedback capture runbook and item-format reference. Update installed route
indexes and root `AGENTS.md` managed block. The capture runbook is agent-facing and
routing-bearing.

### C) Files Touched

```text
components/
  agent-interface/
    managed/
      README.md                                      # modify
      kit-readme.md                                  # modify
      reference/
        root-agents-md-contract.md                    # modify
        repo-feedback-item-format.md                 # create
      runbooks/
        capture-repo-feedback.md                     # create
templates/
  agents/
    AGENTS.managed-block.md                          # modify
tests/
  test_sync_check.py                                 # modify
  test_packaging_resources.py                        # modify in E5
```

### D) Acceptance Criteria And Size

- Size: `M`
- Root managed block routes repo feedback capture separately from Operating Kit feedback.
- Root `AGENTS.md` contract explicitly allows the repo feedback route as a direct managed route.
- Capture runbook has the runbook-authoring compact intention block.
- Capture runbook states triggers, prompt timing, and stop conditions.
- Capture runbook states missing config means authorization-gated auto-check.
- Capture runbook tells agents to silently check existing `gh`, GitHub authentication, and active
  `Codeheart-Digital-Solutions` organization membership before any user-facing feedback prompt.
- Capture runbook states missing `gh`, missing auth, or unverifiable Codeheart organization
  membership means silent unavailability, not tooling readiness.
- Capture runbook tells agents to check config, GitHub remote, and live issue availability only
  after the authorization gate passes.
- Capture runbook states missing labels/templates do not block issue capture.
- Capture runbook requires explicit approval before issue creation.
- Capture runbook provides no browser fallback, manual fallback, local draft fallback, `gh`
  install prompt, or auth-repair prompt when the authorization gate fails.
- Capture runbook stops before GitHub issue creation for security vulnerabilities, sensitive
  disclosures, secrets, credentials, customer or tenant details, account identifiers, raw logs,
  local machine state, private strategy, and raw private evidence.
- Capture runbook routes feedback about the Operating Kit itself to `submit-kit-feedback.md`.
- Item-format reference includes title shape, required body fields, privacy confirmation,
  classification, label fallback, and triage promotion.
- Capture runbook records L1 structured-recipe maturity, fresh-agent executability validation
  tier, evidence shape, blocker shape, and no-promotion boundary.

### E) Dependencies And Critical-Path Notes

Depends on `E2`. The runbook must use the schema semantics from `E2`.

### F) Tasks Checklist

- [ ] Create `components/agent-interface/managed/runbooks/capture-repo-feedback.md` with audience `agent-facing`.
- [ ] Add compact intention block to `capture-repo-feedback.md` covering intent, success, agent judgment boundary, and stop boundary.
- [ ] Add trigger section to `capture-repo-feedback.md` for blockers, runbook failures, script failures, workarounds, user dissatisfaction, docs conflicts, undocumented overrides, missing guidance, and repeated friction.
- [ ] Add prompt-timing section to `capture-repo-feedback.md` for blocker report, direct dissatisfaction, checkpoint, and final summary after the authorization gate passes.
- [ ] Add authorization gate requiring existing `gh`, successful `gh auth status`, and active membership from `gh api user/memberships/orgs/Codeheart-Digital-Solutions --jq .state`.
- [ ] Add explicit rule that failed authorization gate stops silently and must not route to tooling readiness, `gh` install, auth repair, browser fallback, manual fallback, local draft, or any other feedback mechanism.
- [ ] Add destination-resolution procedure to `capture-repo-feedback.md` covering `repo_feedback.mode`, missing config authorization-gated auto-check, GitHub remote detection, and `gh repo view <owner/repo> --json nameWithOwner,hasIssuesEnabled,viewerPermission,isPrivate`.
- [ ] Add issue-draft procedure to `capture-repo-feedback.md` requiring sanitized title, sanitized body, classification fields, and no raw sensitive evidence.
- [ ] Add sensitive-disclosure stop boundary to `capture-repo-feedback.md` covering security vulnerabilities, sensitive disclosures, secrets, credentials, customer and tenant details, account identifiers, raw logs, local machine state, private strategy, and raw private evidence.
- [ ] Add label fallback procedure to `capture-repo-feedback.md` requiring `gh label list --repo <owner/repo> --limit 200 --json name` and use of existing labels only.
- [ ] Add issue-creation approval gate to `capture-repo-feedback.md` before `gh issue create --repo <owner/repo> --title <title> --body-file <body-file>`.
- [ ] Add decline-suppression procedure to `capture-repo-feedback.md` requiring verified maintainer approval before writing `repo_feedback.mode: disabled`.
- [ ] Add Operating Kit boundary procedure to `capture-repo-feedback.md` routing generic kit feedback to `submit-kit-feedback.md`.
- [ ] Add L1 structured-recipe metadata to `capture-repo-feedback.md` covering recipe ID, purpose, inputs, preconditions, approval class, execution surface, evidence output, validation, and stop conditions.
- [ ] Add no-promotion note to `capture-repo-feedback.md` stating v1 creates no reusable script asset, CLI command, wrapper, API, durable helper, package, and tool surface.
- [ ] Create `components/agent-interface/managed/reference/repo-feedback-item-format.md` with required issue body fields, title prefix guidance, privacy confirmation, and sensitive-disclosure stop boundary.
- [ ] Update `components/agent-interface/managed/reference/root-agents-md-contract.md` so direct managed routes include repo feedback capture separately from Operating Kit feedback.
- [ ] Update `components/agent-interface/managed/README.md` with routes for capture runbook and item-format reference.
- [ ] Update `components/agent-interface/managed/kit-readme.md` with installed fallback route for repo feedback capture.
- [ ] Update `templates/agents/AGENTS.managed-block.md` with a concise `Repo feedback capture` managed route.
- [ ] Update `tests/test_sync_check.py` to assert refreshed root `AGENTS.md` contains the repo feedback route and the route target exists.
- [ ] Run `python3 scripts/validate-markdown-headers.py components/agent-interface/managed/runbooks/capture-repo-feedback.md components/agent-interface/managed/reference/repo-feedback-item-format.md`.
- [ ] Run `python3 scripts/validate-public-core.py components/agent-interface/managed/runbooks/capture-repo-feedback.md components/agent-interface/managed/reference/repo-feedback-item-format.md components/agent-interface/managed/reference/root-agents-md-contract.md templates/agents/AGENTS.managed-block.md`.
- [ ] Run `PYTHONPATH=src python3 -m pytest tests/test_sync_check.py -q`.

### G) Implementation Notes

This runbook is routing-bearing. Apply `operation-routing-and-dispatch.md` and keep route
selection separate from issue creation. Apply `operational-recipe-maturity.md`, keep the capture
workflow at L1 structured-recipe maturity, and record that v1 does not promote the workflow into a
script, CLI command, wrapper, or API.

Capture recipe validation tier: fresh-agent executability review plus non-live command-shape
review for documented `gh` examples.

Capture evidence shape: authorization gate result without exposing sensitive auth data, sanitized
issue draft, privacy confirmation, selected existing labels, read-only issue-availability
preflight summary, approval prompt text, and recorded destination or suppression state.

Capture blocker shape: non-secret blocker class, recipe phase, target repository, preflight
command or route attempted, sanitized reason, and user decision needed. Failed authorization gate
is not a blocker to present to the user during normal work; it is silent unavailability.

Runbook `gh` command lines are invocation examples inside the L1 runbook recipe. They are not
reusable script assets, promoted helpers, wrappers, CLI behavior, or APIs.

The setup runbook route becomes public in E4 after the setup runbook exists, so E3 does not
temporarily advertise a missing installed target.

### H) Open Questions

- `OQ-2` does not block because the runbook uses existing authenticated `gh` only and has an
  explicit no-fallback rule.

## E4 - GitHub Issue-Intake Setup Runbook

### A) Epic ID, Title, And Outcome

`E4` - GitHub Issue-Intake Setup Runbook.

Outcome: verified Codeheart organization members have a managed, approval-gated setup path for
repositories where Issues are disabled, unavailable, or missing maintainer-approved standard
labels/templates.

### B) Scope

Add the setup runbook. It is maintainer-facing and external-state-changing. It defines
authorization preflight for verified Codeheart GitHub organization members, read-only repo
preflight, approval gates, Issues enablement, standard label verification/creation,
issue-template file guidance, config recording, validation, and stop conditions.

### C) Files Touched

```text
components/
  agent-interface/
    managed/
      README.md                                      # modify
      kit-readme.md                                  # modify
      runbooks/
        capture-repo-feedback.md                     # modify setup handoff after setup runbook exists
        enable-github-issues-feedback-intake.md      # create
```

### D) Acceptance Criteria And Size

- Size: `M`
- Setup runbook has a compact intention block and `maintainer-facing` audience.
- Setup runbook checks existing `gh`, GitHub authentication, active
  `Codeheart-Digital-Solutions` organization membership, and GitHub remote before writes.
- Setup runbook states missing `gh`, auth, or Codeheart organization membership stops silently in
  normal capture and stops setup without install/browser/manual fallback.
- Setup runbook requires explicit approval before `gh repo edit --enable-issues`.
- Setup runbook requires explicit approval before label creation.
- Setup runbook treats issue-template file creation as repo-owned source change requiring
  approval.
- Setup runbook records configured, declined, and disabled config states.
- Setup runbook does not route missing `gh` through tooling readiness and provides no browser or
  manual setup fallback.
- Capture runbook points to the setup runbook only after the setup runbook exists.
- Managed indexes and installed fallback docs expose the setup route after the target exists.
- Setup runbook records L1 structured-recipe maturity, fresh-agent executability validation tier,
  evidence shape, blocker shape, and no-promotion boundary.

### E) Dependencies And Critical-Path Notes

Depends on `E2` and `E3`. Setup config state must match the schema and capture handoff must point
to this runbook.

### F) Tasks Checklist

- [ ] Create `components/agent-interface/managed/runbooks/enable-github-issues-feedback-intake.md` with audience `maintainer-facing`.
- [ ] Add compact intention block to `enable-github-issues-feedback-intake.md` covering setup scope, approval gates, and stop boundary.
- [ ] Add required inputs section for target repository, feedback trigger summary, desired setup level, authorization-gate evidence, and privacy confirmation.
- [ ] Add authorization preflight requiring existing `gh`, successful `gh auth status`, and active membership from `gh api user/memberships/orgs/Codeheart-Digital-Solutions --jq .state`.
- [ ] Add read-only repo preflight procedure using `git remote get-url origin` and `gh repo view <owner/repo> --json nameWithOwner,hasIssuesEnabled,viewerPermission,isPrivate`.
- [ ] Add explicit no-fallback rule for absent `gh`, failed auth, or unverifiable Codeheart membership: do not use tooling readiness, do not offer install/repair, do not use browser/manual fallback, and stop setup.
- [ ] Add approval packet for enabling Issues with command `gh repo edit <owner/repo> --enable-issues`.
- [ ] Add standard label inventory with `feedback-intake`, `kind-runbook-gap`, `kind-tooling-gap`, `kind-docs-conflict`, `kind-bug`, `kind-feature`, `source-agent-runbook`, `source-user-feedback`, `privacy-sanitized`, and `needs-triage`.
- [ ] Add label verification command `gh label list --repo <owner/repo> --limit 200 --json name`.
- [ ] Add approval packet for creating missing labels with `gh label create <name> --repo <owner/repo> --description <description> --color <hex>`.
- [ ] Add repo-owned issue-template section for `.github/ISSUE_TEMPLATE/repo-feedback.yml` creation after explicit approval.
- [ ] Add config recording section for `repo_feedback.mode: github_issues` and `repo_feedback.github_standardization`.
- [ ] Add decline recording section for `repo_feedback.mode: disabled` with `disabled_reason: verified_maintainer_declined_issue_intake`.
- [ ] Add validation section requiring a second `gh repo view` check, label list check, and `git diff --check` for any repo-owned template file change.
- [ ] Add L1 structured-recipe metadata to `enable-github-issues-feedback-intake.md` covering recipe ID, purpose, inputs, preconditions, approval class, execution surface, evidence output, validation, and stop conditions.
- [ ] Add no-promotion note to `enable-github-issues-feedback-intake.md` stating v1 creates no reusable script asset, CLI command, wrapper, API, durable helper, package, and tool surface.
- [ ] Add setup handoff procedure to `capture-repo-feedback.md` pointing to `enable-github-issues-feedback-intake.md` for disabled Issues, unavailable Issues, and missing user-requested standardization.
- [ ] Update `components/agent-interface/managed/README.md` with the setup runbook route.
- [ ] Update `components/agent-interface/managed/kit-readme.md` with the setup runbook route.
- [ ] Run `python3 scripts/validate-markdown-headers.py components/agent-interface/managed/runbooks/enable-github-issues-feedback-intake.md`.
- [ ] Run `python3 scripts/validate-public-core.py components/agent-interface/managed/runbooks/enable-github-issues-feedback-intake.md`.

### G) Implementation Notes

The setup runbook documents future external changes; executing this Operating Kit implementation
must not enable Issues, create labels, create templates, or create test issues in consumer repos.
The only live GitHub command expected during this implementation is read-only command-shape proof.

Setup recipe validation tier: fresh-agent executability review plus read-only command-shape proof.

Setup evidence shape: target repository, authorization gate result without exposing sensitive auth
data, read-only issue availability preflight, approval packet, approved action summary,
config-state diff summary, and validation summary.

Setup blocker shape: non-secret blocker class, recipe phase, target repository, missing
permission, unavailable Issues state, and user decision needed. Missing `gh`, failed auth, or
unverifiable Codeheart membership is setup unavailability and must not become an install or
browser/manual fallback prompt.

The setup runbook may include short `gh` invocation examples. It must not promote those examples
into reusable script assets, helpers, wrappers, CLI behavior, packages, or APIs in v1.

### H) Open Questions

- `OQ-3` does not block because the issue-template setup path is approval-gated and missing
  templates do not block capture.

## E5 - Packaged Resources, Manifests, Release Surfaces, And Indexes

### A) Epic ID, Title, And Outcome

`E5` - Packaged Resources, Manifests, Release Surfaces, And Indexes.

Outcome: authored managed docs, templates, component manifests, packaged-resource mirrors,
indexes, and release surfaces are synchronized.

### B) Scope

Add the new managed files to component manifests, mirror authoring files into packaged resources,
refresh package and release manifest versions/checksums and consumer-impact records, update docs
indexes, update release notes, and add packaged-resource tests.

### C) Files Touched

```text
components/
  agent-interface/
    component.yaml                                      # modify
src/
  codeheart_operating_kit/
    __init__.py                                         # modify runtime version
    resources/
      components/agent-interface/component.yaml         # mirror
      components/agent-interface/managed/...            # mirror
      templates/agents/AGENTS.managed-block.md          # mirror
      manifest.yaml                                     # modify
docs/
  README.md                                             # modify
  repo/
    README.md                                           # modify
    plans/
      README.md                                         # modify
      plan-register.md                                  # modify
manifest.yaml                                           # modify
release-notes.md                                        # modify
pyproject.toml                                          # modify package version
tests/
  test_packaging_resources.py                           # modify
```

### D) Acceptance Criteria And Size

- Size: `M`
- Component manifest includes new capture/setup runbooks and item-format reference.
- Package metadata version, runtime `__version__`, component version, and release manifests align
  with the selected release version.
- Packaged resources match authored source files.
- Installed fallback test covers new files.
- Root and packaged manifests include the correct component checksum and impact classes.
- Docs indexes include discovery and implementation plan routes.
- Release notes explain consumer impact and adoption behavior.

### E) Dependencies And Critical-Path Notes

Depends on `E3` and `E4`.

### F) Tasks Checklist

- [ ] Add `capture-repo-feedback.md`, `enable-github-issues-feedback-intake.md`, and `repo-feedback-item-format.md` entries to `components/agent-interface/component.yaml`.
- [ ] Bump `components/agent-interface/component.yaml` from the current source version at execution time to the next patch version recorded in the execution log.
- [ ] Bump `pyproject.toml` project version to the selected release version recorded in the execution log.
- [ ] Bump `src/codeheart_operating_kit/__init__.py` `__version__` to the selected release version recorded in the execution log.
- [ ] Copy updated agent-interface managed files into `src/codeheart_operating_kit/resources/components/agent-interface/managed/`.
- [ ] Copy updated `components/agent-interface/component.yaml` into `src/codeheart_operating_kit/resources/components/agent-interface/component.yaml`.
- [ ] Copy updated `templates/agents/AGENTS.managed-block.md` into `src/codeheart_operating_kit/resources/templates/agents/AGENTS.managed-block.md`.
- [ ] Update `manifest.yaml` and `src/codeheart_operating_kit/resources/manifest.yaml` with the new agent-interface checksum.
- [ ] Update root and packaged manifest consumer-impact lists with `instruction-only change`, `validator-only change`, and the canonical security/safety policy impact class.
- [ ] Update `tests/test_packaging_resources.py` to assert packaged fallback installs the capture runbook, setup runbook, item-format reference, and root route target.
- [ ] Update `docs/README.md`, `docs/repo/README.md`, and `docs/repo/plans/README.md` with the implementation plan route.
- [ ] Update `docs/repo/plans/plan-register.md` with the implementation plan entry and relation to `OK-PR-020`.
- [ ] Update `release-notes.md` with repo feedback capture summary, consumer impact, fresh-install behavior, and no-migration note.
- [ ] Run source-to-packaged `diff -q` checks for each mirrored agent-interface file and the managed block template.
- [ ] Run `PYTHONPATH=src python3 -m pytest tests/test_packaging_resources.py -q`.
- [ ] Run `python3 scripts/validate-release-manifest.py`.

### G) Implementation Notes

Use existing packaged-resource mirroring patterns. Do not add a new component. Do not scaffold a
repo feedback state file.

### H) Open Questions

- `OQ-1` is resolved in this epic by recording the selected release version; E7 confirms that
  the package metadata, tag, assets, manifests, and GitHub release all use that same version.

## E6 - Validation, Fresh-Repo Proof, Routing Probe, Review, And Release-Candidate Handoff

### A) Epic ID, Title, And Outcome

`E6` - Validation, Fresh-Repo Proof, Routing Probe, Review, And Release-Candidate Handoff.

Outcome: the implementation is validated locally, route behavior is probed, review evidence is
recorded, and the validated commit is ready for release execution.

### B) Scope

Run focused validation, full validation, read-only GitHub command proof, fresh install proof,
fresh low-context routing probe, review gate, execution-log validation evidence, and
release-candidate register/index refresh.

### C) Files Touched

```text
docs/repo/plans/repo-feedback-capture/
  repo-feedback-capture_implementation_doc.md       # modify status at release-candidate checkpoint
  repo-feedback-capture_execution_log.md            # modify
docs/repo/plans/plan-register.md                    # modify release-candidate lifecycle snapshot
```

### D) Acceptance Criteria And Size

- Size: `M`
- Focused tests and full tests pass.
- Markdown, public-core, schema, release-manifest, and diff checks pass.
- Fresh install proof confirms route targets exist and no default `repo_feedback` block is
  written.
- Read-only GitHub command proof confirms command shape for issue availability checks.
- Fresh low-context routing probe shows correct route selection.
- Negative safety probe confirms missing `gh`, missing auth, or unverifiable Codeheart
  organization membership produces no feedback prompt and no fallback path.
- Review gate has no unresolved material findings.
- Execution log records validation, residual risk, and release-candidate handoff.

### E) Dependencies And Critical-Path Notes

Depends on `E5`.

### F) Tasks Checklist

- [ ] Run `python3 scripts/validate-markdown-headers.py`.
- [ ] Run `python3 scripts/validate-public-core.py`.
- [ ] Run `python3 scripts/validate-json-schemas.py`.
- [ ] Run `python3 scripts/validate-release-manifest.py`.
- [ ] Run `git diff --check`.
- [ ] Run `PYTHONPATH=src python3 -m pytest tests/test_json_schemas.py tests/test_init.py tests/test_sync_check.py tests/test_packaging_resources.py -q`.
- [ ] Run `PYTHONPATH=src python3 -m pytest -q`.
- [ ] Run `gh auth status` as read-only command proof when `gh` already exists.
- [ ] Run `gh api user/memberships/orgs/Codeheart-Digital-Solutions --jq .state` as read-only membership-gate proof when `gh` is already authenticated.
- [ ] Run `gh repo view Codeheart-Digital-Solutions/Codeheart-Operating-Kit --json nameWithOwner,hasIssuesEnabled,viewerPermission,isPrivate` as read-only issue-surface proof after the membership gate passes.
- [ ] Record fresh install proof from `tests/test_init.py` and packaged fallback proof from `tests/test_packaging_resources.py` in the execution log.
- [ ] Run a fresh low-context routing probe with a prompt about an unclear repo runbook default in a newly installed Operating Kit repo.
- [ ] Record probe outcome showing repo feedback capture route, Codeheart membership gate, check-first GitHub behavior, no fallback prompt when the gate fails, and no confusion with Operating Kit feedback.
- [ ] Run a negative fresh-agent probe or equivalent isolated command-path review where `gh` is unavailable, unauthenticated, or organization membership is unverifiable.
- [ ] Record negative-probe outcome showing no repo feedback prompt, no tooling-readiness route, no `gh` install or auth-repair suggestion, no browser/manual fallback, and no local-draft fallback.
- [ ] Run a read-only review of the implementation against the discovery capability scope and this implementation plan.
- [ ] Resolve material review findings in source files and record the resolution in the execution log.
- [ ] Update `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_execution_log.md` with validation summaries, residual risks, and release-candidate handoff.
- [ ] Update `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_implementation_doc.md` status to `release-candidate` after E6 acceptance criteria pass.
- [ ] Update `docs/repo/plans/plan-register.md` with release-candidate lifecycle snapshot after validation.
- [ ] Run final `git status --short` and record changed files in the execution log.

### G) Implementation Notes

Do not create a live GitHub issue during validation. Use read-only GitHub checks and local tests.
If full pytest requires an isolated local environment, create it outside tracked source or under
ignored local state and record that in the execution log.

### H) Open Questions

- None.

## E7 - Release And Approved Consumer Sync

### A) Epic ID, Title, And Outcome

`E7` - Release And Approved Consumer Sync.

Outcome: the validated Operating Kit release is published, release evidence is recorded, and
approved consumer repositories are synced to the released version or have explicit blockers
recorded.

### B) Scope

Follow the Operating Kit release runbook after E6 succeeds. Build and validate release assets,
publish the Git tag and GitHub release from the validated commit, then sync approved consumer
repositories through their local Operating Kit update/sync route. Consumer repository names,
private paths, and repo-specific evidence stay out of this public plan; record them in the
execution log only in public-safe form.

### C) Files Touched

```text
docs/repo/plans/repo-feedback-capture/
  repo-feedback-capture_implementation_doc.md       # modify status at release completion
  repo-feedback-capture_execution_log.md            # modify release and sync evidence
docs/repo/plans/plan-register.md                    # modify completed lifecycle snapshot

external release state:
  git tag v<version>                                # create after validated commit
  GitHub release v<version>                         # publish after release validation
  release assets and checksums                      # build and attach

approved consumer repositories:
  AGENTS.md                                         # sync managed block when changed
  .codeheart/kit/**                                 # sync managed kit content
  .codeheart/kit.lock.yaml                          # update installed version/checksums
```

### D) Acceptance Criteria And Size

- Size: `M`
- Release runbook has been followed from the E6 validated commit.
- Release version is recorded and matches package metadata, release assets, manifests, tag, and
  GitHub release.
- Release notes cover instruction, validator, and security/safety policy impact.
- Release assets and checksum files are built and recorded.
- Installer checksum-mismatch failure and macOS release-asset install validation are recorded.
- Windows validation is recorded through GitHub Actions or a release blocker is recorded.
- Git tag and GitHub release point to the validated commit.
- Approved consumer repositories are synced to the released version, committed/pushed according
  to each repository's policy, or recorded with explicit blockers.
- Consumer sync evidence confirms installed kit version, route target presence, and no unwanted
  consumer-owned overwrite.
- Implementation document, execution log, and plan register reflect completed release and sync.

### E) Dependencies And Critical-Path Notes

Depends on `E6`. Stop before release publication if the source tree changes after E6 validation
without rerunning the affected validation. Stop before consumer sync when a consumer repository
has unrelated dirty state that makes managed sync unsafe under that repository's local
instructions.

### F) Tasks Checklist

- [ ] Read `docs/repo/runbooks/release-operating-kit.md` before release work.
- [ ] Confirm release authority, target version, target commit, and consumer-sync scope in the execution log.
- [ ] Re-run release-critical validation if any source file changed after E6.
- [ ] Run `python3 scripts/build-release-assets.py --version <version> --output-dir dist`.
- [ ] Record release asset filenames and checksum filenames.
- [ ] Verify checksum-mismatch failure behavior for the installer path covered by the release runbook.
- [ ] Validate macOS install from the built release asset in an isolated temporary repository.
- [ ] Validate Windows install through GitHub Actions, or stop and record a release blocker.
- [ ] Confirm release notes and manifests reference the final version and asset checksums.
- [ ] Create the Git tag `v<version>` from the validated commit.
- [ ] Publish the GitHub release with release notes, manifests, assets, installers, and checksums.
- [ ] Record release URL, tag, asset URLs, checksums, validation evidence, and residual risk in the execution log.
- [ ] For each approved consumer repository, read its local `AGENTS.md` before syncing.
- [ ] For each approved consumer repository, record `git status --short` and preserve unrelated changes.
- [ ] For each approved consumer repository, run the Operating Kit update/sync path to the released version.
- [ ] For each approved consumer repository, validate `.codeheart/kit.lock.yaml` records the released version and the repo feedback route targets exist.
- [ ] For each approved consumer repository, run `git diff --check` for touched files.
- [ ] For each approved consumer repository, commit and push the sync changes according to local repository policy, or record the branch/PR/blocker handoff.
- [ ] Update `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_execution_log.md` with release and consumer-sync evidence.
- [ ] Update `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_implementation_doc.md` status to `completed`.
- [ ] Update `docs/repo/plans/plan-register.md` with completed lifecycle snapshot.
- [ ] Run final `git status --short` in the Operating Kit repo and record changed files in the execution log.

### G) Implementation Notes

Release and GitHub publication are external-state-changing actions. Do not create tags, publish
GitHub releases, or mutate consumer repositories unless release authority and consumer-sync scope
are explicit in the execution log.

Consumer sync must use each repository's own instructions and preserve unrelated work. If managed
sync is compatible with dirty unrelated files, update carefully and record the evidence. If the
target managed files are dirty or intent is unclear, stop for that repository and record the
blocker rather than overwriting.

Keep public-core safety: do not add private consumer repo names, local machine paths, secrets,
tokens, raw local logs, or private business context to public release notes, manifests, or this
plan. The execution log may record sanitized consumer-sync evidence.

Release URLs and consumer-sync outcomes are post-release evidence. If recording them changes docs
after the release tag is published, treat that follow-up as documentation-only closeout evidence;
do not change release behavior after tagging without rerunning validation and publishing a new
release.

### H) Open Questions

- `OQ-1` release version is resolved in this epic by the confirmed target version.
- Consumer-sync repository list is resolved at execution time from explicit release authority.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- CLI-assisted issue drafting is deferred until maintainer-run repo feedback capture produces real
  issue examples and stable fields.
- Reusable script assets for feedback drafting and setup are deferred until repeated usage proves
  deterministic mechanics that are safer as tested scripts than as L1 runbook recipes.
- Label/template sync automation is deferred until repeated repo drift makes manual setup costly.
- Cross-repo batch setup is deferred until Codeheart identifies a concrete repo cohort and
  approval model.
- GitHub Projects, dashboards, and milestone automation are deferred until feedback volume
  justifies them.
- Private security disclosure handling is deferred to a separate security disclosure discovery.
- Local draft workflow automation is deferred because local drafts are not part of v1 and should
  not become a fallback for missing `gh`, missing authentication, or unverifiable Codeheart
  organization membership.

## 4.2 Future Considerations

- A future `codeheart-operating-kit feedback draft` command could read `repo_feedback`, gather
  sanitized fields, and print an issue body without creating an issue.
- A future validator could warn when `repo_feedback.mode: github_issues` points at a malformed
  destination.
- A future maintainer runbook could batch-standardize labels/templates across selected
  Codeheart-owned repositories.
- Repeated repo feedback issue patterns may inform Operating Kit doctrine, runbook authoring,
  script promotion, module state routing, or planning workflow updates.

# Revision Notes

- 2026-06-29: Created draft implementation plan from the approved repo feedback capture discovery
  capability scope and user-approved defaults.
- 2026-06-29: Patched review findings for schema partial states, sensitive-disclosure stop
  boundary, setup-route sequencing, canonical consumer-impact wording, and recipe-maturity
  evidence.
- 2026-06-29: Patched Operating Kit 0.1.17 alignment for current component-version baseline,
  recipe validation tier, evidence shape, blocker shape, and explicit no-script-promotion scope.
- 2026-07-02: Patched review findings for disabled-mode suppression semantics, root
  `AGENTS.md` contract coverage, canonical maintainer-facing audience, authorization-gate
  evidence wording, and negative no-fallback validation.
- 2026-07-02: Added explicit release and approved consumer-sync epic, including package-version
  alignment, release asset validation, GitHub publication, and consumer sync evidence.
- 2026-07-02: Activated implementation after explicit user goal to implement the plan.
