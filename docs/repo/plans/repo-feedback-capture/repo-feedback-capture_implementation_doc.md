Last updated: 2026-06-29T14:19:10Z (UTC)
Created: 2026-06-29
Status: draft

# Document Header

## Overview

This implementation plan turns the approved repo feedback capture discovery into a shippable
Operating Kit v1. It adds managed agent guidance for capturing repo-specific feedback, a
check-first GitHub Issues flow, demand-driven setup guidance, optional repo-feedback config
schema support, route visibility for fresh installs, packaged-resource mirrors, validation, and
release notes.

The first release deliberately avoids a new CLI command and avoids automatic GitHub changes.
Agents may draft issues and run read-only preflight checks, but issue creation, label/template
changes, and repository setting changes remain explicit user-approved actions.

## Essential Context Reference Files

| Path | Reason |
| --- | --- |
| `AGENTS.md` | Maintainer bootstrap, public-core safety, and change-safety rules. |
| `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_discovery_doc.md` | Approved capability scope, decisions, non-goals, and handoff requirements. |
| `docs/repo/plans/kit-feedback-intake/kit-feedback-intake_implementation_doc.md` | Precedent for Operating Kit feedback intake, issue forms, labels, managed guidance, and triage. |
| `docs/repo/plans/kit-feedback-intake/kit-feedback-intake_execution_log.md` | Evidence that GitHub Issues were checked first, already enabled, and labels/forms were handled as governance. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Required quality bar for the new durable runbooks. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Route-before-surface doctrine for the new feedback route and fresh-agent probe. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | Required maturity reference for the repeated capture/setup recipes and non-promotion boundary. |
| `components/agent-interface/managed/runbooks/submit-kit-feedback.md` | Existing route for feedback about the Operating Kit itself; repo feedback must not replace it. |
| `components/agent-interface/managed/reference/kit-feedback-item-format.md` | Existing feedback item format precedent. |
| `docs/repo/reference/consumer-impact-classification.md` | Canonical release and manifest impact class vocabulary. |
| `components/agent-interface/component.yaml` | Managed agent-interface component manifest that must include new managed docs. |
| `templates/agents/AGENTS.managed-block.md` | Installed root route surface for fresh repos. |
| `components/agent-interface/managed/README.md` | Managed agent-interface route index. |
| `components/agent-interface/managed/kit-readme.md` | Installed fallback inventory for managed routes. |
| `schemas/kit-config.schema.json` | Optional `repo_feedback` config contract. |
| `src/codeheart_operating_kit/components.py` | Fresh install config writer; v1 must keep `repo_feedback` absent by default. |
| `tests/test_json_schemas.py` | Focused schema tests for optional `repo_feedback` states. |
| `tests/test_init.py` | Fresh install tests for no default `repo_feedback` block and route target presence. |
| `tests/test_sync_check.py` | Sync/check route target validation. |
| `tests/test_packaging_resources.py` | Packaged resource fallback coverage. |
| `release-notes.md` | Release-facing consumer-impact notes. |

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
- the capture runbook tells agents to check config, GitHub remote, and live issue availability
  before offering setup;
- missing `repo_feedback` config means unconfigured auto-check, not disabled;
- GitHub Issues capture works without standard labels/templates by putting classification in the
  issue body and using only labels that exist;
- sensitive, private, and security-disclosure material stops before GitHub issue creation and is
  not stored in feedback drafts;
- setup is offered only when Issues are unavailable, disabled, or the user wants missing standard
  labels/templates;
- decline can be recorded as disabled/suppressed state after explicit user approval;
- `schemas/kit-config.schema.json` accepts optional configured, disabled, and local-draft states
  while preserving old config validity;
- fresh install tests prove no `repo_feedback` block is written by default;
- packaged-resource fallback installs the new runbooks and reference files;
- validation covers Markdown headers, public-core hygiene, schema behavior, route targets,
  packaged resources, release manifests, focused tests, full tests, and a fresh low-context
  routing probe;
- repeated capture/setup recipes record L1 structured-recipe maturity and no-promotion evidence;
- release notes record instruction, schema, and safety-policy impact.

## 1.2 Project And Problem Context

The Operating Kit already has a public feedback intake path for feedback about the kit itself.
That path sends sanitized feedback to the public `Codeheart-Operating-Kit` GitHub Issues backlog
and provides a maintainer triage workflow.

Repo-specific feedback is different. A runbook gap, script failure, conservative undocumented
override, docs conflict, or user dissatisfaction usually belongs in the repository that owns the
affected artifact. Without a managed route, agents can solve the immediate issue in chat and lose
the reusable maintenance signal.

The approved direction is to make the Operating Kit own the generic capture route while each
repository owns its feedback destination and issue triage. A fresh repo should not need feedback
setup during onboarding. When a real feedback item appears, the agent checks whether GitHub
Issues already works, drafts an issue when it does, and offers setup only when needed.

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

- Missing `repo_feedback` config remains valid and means auto-check on first feedback trigger.
- Optional `repo_feedback` config records durable preferences after first use, setup, local draft
  opt-in, or decline.
- Managed runbooks describe the exact check-first, draft, approval, setup, and suppression flow.
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
          repo-feedback-item-format.md                    # create
        runbooks/
          capture-repo-feedback.md                        # create
          enable-github-issues-feedback-intake.md          # create
  templates/
    agents/
      AGENTS.managed-block.md                             # modify
  src/
    codeheart_operating_kit/
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
```

## 2.2 Open Questions And Assumptions Requiring Clarification

`OQ-1`: Exact release version.

- `BLOCKER: no`
- Affects: `E6`
- Unlocks: release notes and release manifest wording.
- Recommended default: use the next patch release available at execution time and record the
  exact version in the execution log.

`OQ-2`: GitHub CLI availability in consumer repos.

- `BLOCKER: no`
- Affects: `E3`, `E4`, `E6`
- Unlocks: preflight and setup command lane.
- Recommended default: make `gh` the default command lane because it supports `gh repo view
  --json hasIssuesEnabled`, `gh repo edit --enable-issues`, `gh label list`, `gh label create`,
  and `gh issue create`; route missing local `gh` through tooling readiness and provide browser
  issue creation as the non-CLI fallback.

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

`AD-2`: Missing `repo_feedback` config means auto-check.

1. Problem being solved: fresh repos must work without onboarding-time feedback setup.
2. Simplest working solution: keep `repo_feedback` optional; absence means the capture runbook
   checks the GitHub remote and live issue availability when a real feedback item appears.
3. What may change in 6-12 months: onboarding may offer feedback setup after usage data proves it
   is worth front-loading.
4. Rationale: this preserves lightweight onboarding and matches the Operating Kit feedback
   precedent where Issues were already enabled.
5. Alternatives considered: writing `repo_feedback.mode: auto_check` during init was rejected
   because it creates config churn without adding capability.

`AD-3`: Schema records preferences, not behavior automation.

1. Problem being solved: agents need durable state for configured destinations and declined setup.
2. Simplest working solution: add an optional `repo_feedback` schema with `github_issues`,
   `disabled`, and `local_draft_only` states.
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

`AD-5`: Use `gh` as the default command lane and keep browser fallback.

1. Problem being solved: agents need exact, testable preflight and setup commands.
2. Simplest working solution: document `gh repo view`, `gh issue create`, `gh repo edit`, and
   `gh label` commands, with approval gates before write commands.
3. What may change in 6-12 months: a connector or CLI wrapper may provide a cleaner interface.
4. Rationale: `gh` is already available in the maintainer environment and is the direct GitHub
   governance lane.
5. Alternatives considered: a new Operating Kit CLI command was rejected for v1 because the manual
   route needs real usage first.

`AD-6`: Treat setup as repository governance.

1. Problem being solved: enabling Issues, labels, and templates changes shared repository state.
2. Simplest working solution: the setup runbook preflights, explains the change, asks approval,
   executes the approved change, then verifies it.
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
| `E6` | Validation, fresh-repo proof, routing probe, review, and closeout evidence are complete. | M | `E5` |

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

Outcome: `.codeheart/kit.config.yaml` accepts optional repo feedback preferences, while fresh
installs still omit `repo_feedback` and therefore use auto-check behavior on first feedback
trigger.

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
- Valid `github_issues` config includes destination owner and repo.
- Valid `disabled` config records decline or owner-policy reason.
- Valid `local_draft_only` config is explicit and local-user-owned.
- Invalid partial and mode-incompatible states fail schema tests.

### E) Dependencies And Critical-Path Notes

Depends on `E1`. This epic establishes config semantics that runbooks reference.

### F) Tasks Checklist

- [ ] Add optional top-level `repo_feedback` to `schemas/kit-config.schema.json`.
- [ ] Add `repo_feedback.mode` enum values `github_issues`, `disabled`, and `local_draft_only` to `schemas/kit-config.schema.json`.
- [ ] Add `repo_feedback.destination.type` const `github_issues` plus `owner` and `repo` string fields for `github_issues` mode.
- [ ] Add `repo_feedback.github_standardization.labels` and `repo_feedback.github_standardization.issue_templates` enum values `configured`, `not_configured`, and `declined`.
- [ ] Add `repo_feedback.disabled_reason` enum values `user_declined_issue_intake`, `repo_owner_policy`, `issues_unavailable`, and `other` for `disabled` mode.
- [ ] Add `repo_feedback.draft_path` const `.codeheart/user/feedback/` for `local_draft_only` mode.
- [ ] Keep `src/codeheart_operating_kit/components.py` fresh install config output unchanged for default `repo_feedback` absence.
- [ ] Add `tests/test_json_schemas.py` coverage for missing `repo_feedback` config validity.
- [ ] Add `tests/test_json_schemas.py` coverage for valid `github_issues`, `disabled`, and `local_draft_only` config blocks.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `github_issues` without `destination.type`.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `github_issues` without destination owner and repo.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `disabled` without `disabled_reason`.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting `local_draft_only` without `draft_path`.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting invalid `github_standardization.labels` and `github_standardization.issue_templates` enum values.
- [ ] Add `tests/test_json_schemas.py` coverage rejecting mode-incompatible fields across `github_issues`, `disabled`, and `local_draft_only`.
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
- Capture runbook has the runbook-authoring compact intention block.
- Capture runbook states triggers, prompt timing, and stop conditions.
- Capture runbook states missing config means auto-check.
- Capture runbook tells agents to check config, GitHub remote, and live issue availability before
  offering setup.
- Capture runbook states missing labels/templates do not block issue capture.
- Capture runbook requires explicit approval before issue creation.
- Capture runbook stops before GitHub issue creation for security vulnerabilities, sensitive
  disclosures, secrets, credentials, customer or tenant details, account identifiers, raw logs,
  local machine state, private strategy, and raw private evidence.
- Capture runbook routes feedback about the Operating Kit itself to `submit-kit-feedback.md`.
- Item-format reference includes title shape, required body fields, privacy confirmation,
  classification, label fallback, and triage promotion.
- Capture runbook records L1 structured-recipe maturity and no-promotion boundary.

### E) Dependencies And Critical-Path Notes

Depends on `E2`. The runbook must use the schema semantics from `E2`.

### F) Tasks Checklist

- [ ] Create `components/agent-interface/managed/runbooks/capture-repo-feedback.md` with audience `agent-facing`.
- [ ] Add compact intention block to `capture-repo-feedback.md` covering intent, success, agent judgment boundary, and stop boundary.
- [ ] Add trigger section to `capture-repo-feedback.md` for blockers, runbook failures, script failures, workarounds, user dissatisfaction, docs conflicts, undocumented overrides, missing guidance, and repeated friction.
- [ ] Add prompt-timing section to `capture-repo-feedback.md` for blocker report, direct dissatisfaction, checkpoint, and final summary.
- [ ] Add destination-resolution procedure to `capture-repo-feedback.md` covering `repo_feedback.mode`, missing config auto-check, GitHub remote detection, and `gh repo view <owner/repo> --json nameWithOwner,hasIssuesEnabled,viewerPermission,isPrivate`.
- [ ] Add issue-draft procedure to `capture-repo-feedback.md` requiring sanitized title, sanitized body, classification fields, and no raw sensitive evidence.
- [ ] Add sensitive-disclosure stop boundary to `capture-repo-feedback.md` covering security vulnerabilities, sensitive disclosures, secrets, credentials, customer and tenant details, account identifiers, raw logs, local machine state, private strategy, and raw private evidence.
- [ ] Add label fallback procedure to `capture-repo-feedback.md` requiring `gh label list --repo <owner/repo> --limit 200 --json name` and use of existing labels only.
- [ ] Add issue-creation approval gate to `capture-repo-feedback.md` before `gh issue create --repo <owner/repo> --title <title> --body-file <body-file>`.
- [ ] Add decline-suppression procedure to `capture-repo-feedback.md` requiring user approval before writing `repo_feedback.mode: disabled`.
- [ ] Add Operating Kit boundary procedure to `capture-repo-feedback.md` routing generic kit feedback to `submit-kit-feedback.md`.
- [ ] Add L1 structured-recipe maturity and no-promotion note to `capture-repo-feedback.md`.
- [ ] Create `components/agent-interface/managed/reference/repo-feedback-item-format.md` with required issue body fields, title prefix guidance, privacy confirmation, and sensitive-disclosure stop boundary.
- [ ] Update `components/agent-interface/managed/README.md` with routes for capture runbook and item-format reference.
- [ ] Update `components/agent-interface/managed/kit-readme.md` with installed fallback route for repo feedback capture.
- [ ] Update `templates/agents/AGENTS.managed-block.md` with a concise `Repo feedback capture` managed route.
- [ ] Update `tests/test_sync_check.py` to assert refreshed root `AGENTS.md` contains the repo feedback route and the route target exists.
- [ ] Run `python3 scripts/validate-markdown-headers.py components/agent-interface/managed/runbooks/capture-repo-feedback.md components/agent-interface/managed/reference/repo-feedback-item-format.md`.
- [ ] Run `python3 scripts/validate-public-core.py components/agent-interface/managed/runbooks/capture-repo-feedback.md components/agent-interface/managed/reference/repo-feedback-item-format.md templates/agents/AGENTS.managed-block.md`.
- [ ] Run `PYTHONPATH=src python3 -m pytest tests/test_sync_check.py -q`.

### G) Implementation Notes

This runbook is routing-bearing. Apply `operation-routing-and-dispatch.md` and keep route
selection separate from issue creation. Apply `operational-recipe-maturity.md`, keep the capture
workflow at L1 structured-recipe maturity, and record that v1 does not promote the workflow into a
script, CLI command, wrapper, or API.

The setup runbook route becomes public in E4 after the setup runbook exists, so E3 does not
temporarily advertise a missing installed target.

### H) Open Questions

- `OQ-2` does not block because the runbook has a `gh` lane and browser fallback.

## E4 - GitHub Issue-Intake Setup Runbook

### A) Epic ID, Title, And Outcome

`E4` - GitHub Issue-Intake Setup Runbook.

Outcome: agents have a managed, approval-gated setup path for repositories where Issues are
disabled, unavailable, or missing user-approved standard labels/templates.

### B) Scope

Add the setup runbook. It is maintainer-facing and external-state-changing. It defines read-only
preflight, approval gates, Issues enablement, standard label verification/creation, issue-template
file guidance, config recording, validation, and stop conditions.

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
- Setup runbook has a compact intention block and maintainer-facing audience.
- Setup runbook checks GitHub remote and `gh` auth before writes.
- Setup runbook requires explicit approval before `gh repo edit --enable-issues`.
- Setup runbook requires explicit approval before label creation.
- Setup runbook treats issue-template file creation as repo-owned source change requiring
  approval.
- Setup runbook records configured, declined, and disabled config states.
- Setup runbook routes missing `gh` through tooling readiness and provides browser setup fallback.
- Capture runbook points to the setup runbook only after the setup runbook exists.
- Managed indexes and installed fallback docs expose the setup route after the target exists.

### E) Dependencies And Critical-Path Notes

Depends on `E2` and `E3`. Setup config state must match the schema and capture handoff must point
to this runbook.

### F) Tasks Checklist

- [ ] Create `components/agent-interface/managed/runbooks/enable-github-issues-feedback-intake.md` with audience `maintainer-facing`.
- [ ] Add compact intention block to `enable-github-issues-feedback-intake.md` covering setup scope, approval gates, and stop boundary.
- [ ] Add required inputs section for target repository, feedback trigger summary, desired setup level, and privacy confirmation.
- [ ] Add read-only preflight procedure using `git remote get-url origin` and `gh repo view <owner/repo> --json nameWithOwner,hasIssuesEnabled,viewerPermission,isPrivate`.
- [ ] Add missing-tool route for absent `gh` through `handle-tooling-readiness.md`.
- [ ] Add browser fallback path for users who prefer GitHub web setup.
- [ ] Add approval packet for enabling Issues with command `gh repo edit <owner/repo> --enable-issues`.
- [ ] Add standard label inventory with `feedback-intake`, `kind-runbook-gap`, `kind-tooling-gap`, `kind-docs-conflict`, `kind-bug`, `kind-feature`, `source-agent-runbook`, `source-user-feedback`, `privacy-sanitized`, and `needs-triage`.
- [ ] Add label verification command `gh label list --repo <owner/repo> --limit 200 --json name`.
- [ ] Add approval packet for creating missing labels with `gh label create <name> --repo <owner/repo> --description <description> --color <hex>`.
- [ ] Add repo-owned issue-template section for `.github/ISSUE_TEMPLATE/repo-feedback.yml` creation after explicit approval.
- [ ] Add config recording section for `repo_feedback.mode: github_issues` and `repo_feedback.github_standardization`.
- [ ] Add decline recording section for `repo_feedback.mode: disabled` with `disabled_reason: user_declined_issue_intake`.
- [ ] Add validation section requiring a second `gh repo view` check, label list check, and `git diff --check` for any repo-owned template file change.
- [ ] Add setup handoff procedure to `capture-repo-feedback.md` pointing to `enable-github-issues-feedback-intake.md` for disabled Issues, unavailable Issues, and missing user-requested standardization.
- [ ] Update `components/agent-interface/managed/README.md` with the setup runbook route.
- [ ] Update `components/agent-interface/managed/kit-readme.md` with the setup runbook route.
- [ ] Run `python3 scripts/validate-markdown-headers.py components/agent-interface/managed/runbooks/enable-github-issues-feedback-intake.md`.
- [ ] Run `python3 scripts/validate-public-core.py components/agent-interface/managed/runbooks/enable-github-issues-feedback-intake.md`.

### G) Implementation Notes

The setup runbook documents future external changes; executing this Operating Kit implementation
must not enable Issues, create labels, create templates, or create test issues in consumer repos.
The only live GitHub command expected during this implementation is read-only command-shape proof.

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
refresh package manifest checksums and consumer-impact records, update docs indexes, update
release notes, and add packaged-resource tests.

### C) Files Touched

```text
components/
  agent-interface/
    component.yaml                                      # modify
src/
  codeheart_operating_kit/
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
tests/
  test_packaging_resources.py                           # modify
```

### D) Acceptance Criteria And Size

- Size: `M`
- Component manifest includes new capture/setup runbooks and item-format reference.
- Packaged resources match authored source files.
- Installed fallback test covers new files.
- Root and packaged manifests include the correct component checksum and impact classes.
- Docs indexes include discovery and implementation plan routes.
- Release notes explain consumer impact and adoption behavior.

### E) Dependencies And Critical-Path Notes

Depends on `E3` and `E4`.

### F) Tasks Checklist

- [ ] Add `capture-repo-feedback.md`, `enable-github-issues-feedback-intake.md`, and `repo-feedback-item-format.md` entries to `components/agent-interface/component.yaml`.
- [ ] Bump `components/agent-interface/component.yaml` version from `0.1.15` to the next patch version recorded in the execution log.
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

- `OQ-1` is resolved during execution by recording the release version chosen for this change.

## E6 - Validation, Fresh-Repo Proof, Routing Probe, Review, And Closeout

### A) Epic ID, Title, And Outcome

`E6` - Validation, Fresh-Repo Proof, Routing Probe, Review, And Closeout.

Outcome: the implementation is validated locally, route behavior is probed, review evidence is
recorded, and the plan is ready for PR/release handoff.

### B) Scope

Run focused validation, full validation, read-only GitHub command proof, fresh install proof,
fresh low-context routing probe, review gate, execution-log closeout, and final register/index
refresh.

### C) Files Touched

```text
docs/repo/plans/repo-feedback-capture/
  repo-feedback-capture_implementation_doc.md       # modify status at completion
  repo-feedback-capture_execution_log.md            # modify
docs/repo/plans/plan-register.md                    # modify lifecycle snapshot
```

### D) Acceptance Criteria And Size

- Size: `M`
- Focused tests and full tests pass.
- Markdown, public-core, schema, release-manifest, and diff checks pass.
- Fresh install proof confirms route targets exist and no default `repo_feedback` block is
  written.
- Read-only GitHub command proof confirms command shape for issue availability checks.
- Fresh low-context routing probe shows correct route selection.
- Review gate has no unresolved material findings.
- Execution log records validation, residual risk, and release handoff.

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
- [ ] Run `gh repo view Codeheart-Digital-Solutions/Codeheart-Operating-Kit --json nameWithOwner,hasIssuesEnabled,viewerPermission,isPrivate` as read-only command proof.
- [ ] Record fresh install proof from `tests/test_init.py` and packaged fallback proof from `tests/test_packaging_resources.py` in the execution log.
- [ ] Run a fresh low-context routing probe with a prompt about an unclear repo runbook default in a newly installed Operating Kit repo.
- [ ] Record probe outcome showing repo feedback capture route, check-first GitHub behavior, and no confusion with Operating Kit feedback.
- [ ] Run a read-only review of the implementation against the discovery capability scope and this implementation plan.
- [ ] Resolve material review findings in source files and record the resolution in the execution log.
- [ ] Update `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_execution_log.md` with validation summaries and residual risks.
- [ ] Update `docs/repo/plans/repo-feedback-capture/repo-feedback-capture_implementation_doc.md` status to `completed` after all acceptance criteria pass.
- [ ] Update `docs/repo/plans/plan-register.md` with completed lifecycle snapshot after completion.
- [ ] Run final `git status --short` and record changed files in the execution log.

### G) Implementation Notes

Do not create a live GitHub issue during validation. Use read-only GitHub checks and local tests.
If full pytest requires an isolated local environment, create it outside tracked source or under
ignored local state and record that in the execution log.

### H) Open Questions

- None.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- CLI-assisted issue drafting is deferred until manual repo feedback capture produces real issue
  examples and stable fields.
- Label/template sync automation is deferred until repeated repo drift makes manual setup costly.
- Cross-repo batch setup is deferred until Codeheart identifies a concrete repo cohort and
  approval model.
- GitHub Projects, dashboards, and milestone automation are deferred until feedback volume
  justifies them.
- Private security disclosure handling is deferred to a separate security disclosure discovery.
- Local draft workflow automation is deferred because local drafts are explicit opt-in and not
  the preferred durable inbox.

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
