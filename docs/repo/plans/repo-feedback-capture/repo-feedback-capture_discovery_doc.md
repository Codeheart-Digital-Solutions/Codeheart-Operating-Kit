Last updated: 2026-06-29T13:58:49Z (UTC)
Created: 2026-06-29
Status: draft

# Repo Feedback Capture And Issue Intake Discovery

## Discovery Status

Input state: new Operating Kit doctrine request after repo and module operations showed that
small runbook gaps, script failures, user dissatisfaction, and agent detours can be solved in chat
but then disappear before they become durable maintenance work.

Output target: manual-review-ready draft. This document captures the proposed generic standard
for how Operating Kit-guided agents should recognize repo-specific feedback, route it to the
owning repository's issue intake when already available, and avoid repeatedly prompting users when
issue intake is unavailable or declined.

Implementation planning is intentionally out of scope for this draft. This discovery does not
change managed docs, root instructions, repo configuration schemas, GitHub repository settings,
labels, templates, CLI behavior, or release assets.

## User Intention

Codeheart wants the Operating Kit to help agents capture smaller operational issues over time
across configured repositories, so maintainers can later batch triage them into direct patches,
discovery documents, or implementation plans.

The mechanism should apply broadly to Codeheart-owned or otherwise configured maintainer-operated
repositories. It should not be limited to feedback about the Operating Kit itself. When an agent
encounters a feedback-worthy issue while running another runbook, it should know when to offer
durable capture, where to route it, how to classify it, and when to stop asking.

## Problem Framing

The current Operating Kit has a public-safe feedback workflow for issues about the Operating Kit
itself. That workflow sends sanitized, shareable kit feedback to the public
`Codeheart-Operating-Kit` GitHub Issues backlog.

That does not fully cover repo-specific friction. Examples include:

- a repository runbook has an unclear default;
- a script block fails or needs a workaround;
- a user is dissatisfied with an agent outcome;
- docs and implementation disagree;
- an agent overrides a documented default for unstated caution;
- a module or repo has missing setup guidance.

These observations may belong in the repository that owns the runbook, module, product, or
script. Sending every such observation to the public Operating Kit backlog would create noise and
may expose repo-specific context. Leaving them in chat loses maintenance evidence.

The desired model is:

```text
agent notices feedback-worthy friction
-> resolve or report the immediate work first
-> attempt repo feedback route
-> check repo feedback config, GitHub remote, and live issue availability
-> if GitHub Issues already works, draft the issue and ask before creating it
-> if Issues are unavailable or disabled, ask once whether to configure issue intake
-> if declined, record suppression and stop prompting repeatedly
-> triage accepted issues into patch, batch, discovery, or implementation planning
```

## Goals

- Define a generic Operating Kit route for repo-specific feedback capture.
- Keep the existing Operating Kit feedback workflow focused on feedback about the kit itself.
- Make GitHub Issues the default durable inbox for configured repositories.
- Check whether GitHub Issues already works before offering setup.
- Use demand-driven setup when feedback is first identified, not mandatory setup during kit
  onboarding.
- Give agents reliable triggers for when to offer feedback capture during other runbooks.
- Avoid repeated prompts when a user declines GitHub issue setup for a repository.
- Require explicit user approval before creating issues or changing GitHub repository settings.
- Define a small classification scheme that supports later batch triage.
- Preserve public-core and repo privacy even when issues are private.
- Define how issue intake promotes into direct patches, discovery, or implementation plans.

## Non-Goals

- Do not implement the route in this discovery.
- Do not change GitHub settings, create labels, or add issue templates in this discovery.
- Do not make GitHub Issues mandatory for Operating Kit installation or onboarding.
- Do not create a local draft feedback system as the default fallback.
- Do not auto-create issues without user approval.
- Do not make every minor detour interrupt the current runbook.
- Do not store secrets, API keys, tokens, OAuth codes, account identifiers, customer or tenant
  details, raw logs, raw document content, local machine dumps, private file paths, or private
  strategy in feedback issues.
- Do not replace plan registers, discovery documents, implementation plans, execution logs, PR
  review, or release notes with issue intake.

## Public-Core Safety

This is public Operating Kit discovery. Examples and motivations must stay generic. Do not include
private repository names, tenant IDs, account identifiers, credentials, raw logs, machine-specific
absolute paths, customer data, mailbox or document content, business records, or restricted
strategy.

Future implementation may include a safety-policy dimension because it changes agent behavior for
external issue creation and repository-setting setup prompts. At minimum it is an
`instruction-only change`; if implementation adds repo config schema behavior, issue templates,
labels, validators, or CLI support, the impact class must be broadened accordingly.

## Current Evidence

| Source | Finding | Discovery implication |
| --- | --- | --- |
| Existing kit feedback intake | Public GitHub Issues are already the canonical backlog for sanitized Operating Kit feedback. | Reuse the pattern, but do not overload it with repo-specific feedback. |
| Kit feedback triage runbook | Existing triage can close issues as consumer-specific and route reusable kit work to discovery or implementation planning. | Repo feedback needs a parallel ownership rule: the owning repo is the first durable inbox. |
| Operation routing doctrine | Agents route through the highest-authority applicable operating source before choosing execution surfaces. | Feedback capture needs destination resolution before issue creation. |
| Runbook authoring standards | Runbooks own procedure, stop conditions, evidence, and validation. | Feedback-worthy runbook gaps should become issues against the owning repo or module. |
| Tooling readiness route | Missing tooling should be handled through a managed route instead of ad hoc setup. | Issue intake setup should have a runbook and explicit approval, not improvised GitHub changes. |
| Plan-register doctrine | Planning docs own accepted discovery and implementation state, while registers own compact index metadata. | Issues are intake objects; accepted larger changes should promote to formal plans and registers. |
| Repeated repo/module operation | Small confusions, conservative overrides, and recovered detours can disappear after the immediate task succeeds. | Agents need a trigger and prompt discipline that captures useful lessons without derailing work. |
| Kit feedback intake implementation | The Operating Kit feedback implementation preflight found GitHub Issues already enabled and then added issue forms plus labels through planned governance. | Repo feedback should check whether Issues already works before asking to set anything up. |

## Decision Ledger

| ID | Decision | Class | Recommendation | Blocks |
| --- | --- | --- | --- | --- |
| `D-1` | Ownership boundary | blocking | Operating Kit owns the generic feedback-capture route; each repository owns its configured feedback destination and triage. | route design |
| `D-2` | Eligible repositories | implementation-shaping | All Operating Kit consumer repositories are eligible, but feedback capture activates through repo config or first attempted feedback route. | config design |
| `D-3` | Default durable inbox | blocking | Use GitHub Issues as the default durable inbox when configured or live preflight confirms it already works. | setup runbook |
| `D-4` | Setup timing | implementation-shaping | Check first, then configure on demand only when Issues is unavailable or repo feedback preferences are absent and needed. | onboarding and prompt policy |
| `D-5` | Unavailable or declined issue intake | blocking | Offer setup only after the check fails; if declined, record disabled or suppressed state and do not keep asking. | config schema |
| `D-6` | Agent triggers | implementation-shaping | Offer capture for blockers, runbook/script failures, workarounds, user dissatisfaction, docs conflicts, undocumented overrides, and repeated friction. | managed runbook |
| `D-7` | Prompt timing | implementation-shaping | Prompt at blocker report, immediate user dissatisfaction, checkpoint, or final summary; do not interrupt every minor recovered detour. | managed runbook |
| `D-8` | Approval boundary | blocking | Agents may draft issues but must ask before issue creation, label/template creation, or GitHub repository setting changes. | external-action safety |
| `D-9` | Classification | implementation-shaping | Use labels as canonical classification and title prefixes as readable summaries. | labels/templates |
| `D-10` | Privacy posture | blocking | Sanitize all issue content; private repos do not allow secrets or raw sensitive evidence. | issue format |
| `D-11` | Triage and promotion | implementation-shaping | Issues are intake; triage promotes to direct patch, batch milestone, discovery, implementation plan, defer, duplicate, or close. | maintainer runbook |

### D-1 - Ownership Boundary

Question: Who owns the generic feedback-capture behavior and who owns individual feedback items?

Recommendation: Operating Kit owns the generic agent route, trigger model, privacy rules,
destination-resolution procedure, and prompt discipline. The repository that owns the affected
runbook, script, module, product, or docs owns the feedback issue and triage.

Rationale: The behavior must be reusable across repos, so the Operating Kit should teach agents
how to recognize and route feedback. The actual issue belongs where maintainers can act on it.
Operating Kit feedback remains reserved for reusable kit doctrine, tooling, sync, or managed
content issues.

State: `draft recommendation`.

### D-2 - Eligible Repositories

Question: Which repositories should use the route?

Recommendation: Treat all Operating Kit consumer repositories as eligible for repo feedback
capture. Activation should depend on repo-local configuration or a demand-driven first-use setup
prompt.

Rationale: The user intention is broad collection across Codeheart-owned and configured
maintainer-operated repos. A universal mechanism is simpler than hand-authoring separate guidance
per repo. The route still needs repo-local destination state because not every repo has GitHub
Issues enabled or wants issue intake.

State: `draft recommendation`.

### D-3 - Default Durable Inbox

Question: What is the default durable inbox for repo-specific feedback?

Recommendation: Use GitHub Issues when the repository has a configured GitHub issue destination
or a live preflight confirms that the current GitHub repository already has Issues enabled and the
agent can reach the issue surface. The user must still approve creating the issue. Repo
configuration should record durable preferences after first use or after setup/decline, but a
missing config entry should not by itself block a check-first GitHub Issues route.

Rationale: GitHub Issues are searchable, linkable from PRs, and suitable for small triage items.
They also let maintainers batch related issues into milestones or formal planning work. Local
draft files are not durable enough as the default and are easy to forget.

State: `draft recommendation`.

### D-4 - Setup Timing

Question: Should GitHub Issues setup happen during Operating Kit onboarding?

Recommendation: Do not make GitHub Issues setup part of default Operating Kit onboarding. When an
agent identifies a real feedback-worthy issue, it should first check whether GitHub Issues already
works for the current repo. Run feedback setup on demand only when the check shows that the issue
surface is unavailable, disabled, missing required repo preferences, or missing agreed labels or
templates that the user wants standardized.

Rationale: Onboarding should stay light. Issue setup is easier to understand when there is a
concrete issue to capture. This also avoids forcing GitHub configuration on local-only repos,
external repos, or repos that intentionally do not use Issues.

State: `draft recommendation`.

### D-5 - Unavailable Or Declined Issue Intake

Question: What happens when GitHub Issues are not configured or the user declines setup?

Recommendation: When a feedback item appears and no configured issue destination exists, the
agent should first inspect the GitHub remote and live issue availability. If Issues already works,
draft the issue and ask before creating it; optionally record the destination in repo config after
approval. If Issues are disabled, unavailable, or the user wants standardized labels/templates
that are missing, offer to run an issue-intake setup runbook. If the user declines setup, record a
repo-local preference that suppresses repeated GitHub setup prompts. Default declined behavior
should be disabled/suppressed, not local-draft-only.

Candidate config shape:

```yaml
repo_feedback:
  mode: disabled
  suppress_prompts: true
  disabled_reason: user_declined_issue_intake
```

Local draft mode may exist, but only when the user explicitly chooses it.

Rationale: If GitHub Issues are the intended durable inbox, local drafts are a weak default.
Repeated prompts would be annoying and would teach agents to interrupt the user. Disabled/suppressed
state respects the user's decision while preserving a future manual path such as "enable repo
feedback." The Operating Kit feedback intake precedent also shows that Issues may already be
enabled, so setup should not be assumed.

State: `draft recommendation`.

### D-6 - Agent Triggers

Question: When should agents offer feedback capture?

Recommendation: Agents should offer capture when one of these happens during repo work:

- a runbook step fails, is ambiguous, or needs a workaround;
- a script block fails, produces fragile output, or needs manual repair;
- docs and implementation disagree;
- local tooling readiness blocks or materially delays the task;
- the user says an outcome is confusing, too dense, wrong, unsatisfactory, or not what they meant;
- the agent intentionally overrides a documented default;
- the agent discovers missing setup/onboarding guidance;
- a recovered detour is likely to recur for future agents or users;
- the same friction appears repeatedly across sessions or repos.

Do not offer issue capture for trivial typos, one-off user preference, sensitive incident details,
or observations that belong only in the current final summary.

State: `draft recommendation`.

### D-7 - Prompt Timing

Question: When should agents prompt the user?

Recommendation: Prompt immediately only when the user expresses dissatisfaction or the task is
blocked. For recovered detours, mention feedback capture at the next natural checkpoint or final
summary. Do not interrupt every small detour while another runbook is in progress.

Rationale: Feedback capture should make small issues durable without derailing the primary task.
The agent should finish the operational work when possible, then capture the lesson.

State: `draft recommendation`.

### D-8 - Approval Boundary

Question: What may agents do automatically?

Recommendation: Agents may identify feedback-worthy friction, inspect repo-local feedback config,
draft a sanitized issue body, and ask for approval. Agents must not create issues, enable GitHub
Issues, create labels, add issue templates, or change repository settings without explicit user
approval.

Rationale: Issue creation and GitHub settings are external-state-changing actions. Labels and
templates are repository governance changes. Approval keeps the route safe and predictable.

State: `draft recommendation`.

### D-9 - Classification

Question: How should feedback be named and classified?

Recommendation: Use labels as canonical classification and optional title prefixes for human
scanning.

Recommended title shape:

```text
[feedback][<kind>][<area>] <short issue title>
```

Recommended standard labels:

- `feedback-intake`;
- `kind-runbook-gap`;
- `kind-tooling-gap`;
- `kind-docs-conflict`;
- `kind-bug`;
- `kind-feature`;
- `source-agent-runbook`;
- `source-user-feedback`;
- `privacy-sanitized`;
- `needs-triage`.

Area should be an issue field by default. Area labels may be configured per repo when useful, but
the Operating Kit should not require every repo to pre-create dynamic `area-*` labels.

State: `draft recommendation`.

### D-10 - Privacy Posture

Question: What must never enter feedback issues?

Recommendation: All feedback issues must be sanitized, even in private repositories. Do not put
these in issues:

- secrets, credentials, API keys, passwords, tokens, OAuth codes, or MFA material;
- customer, tenant, account, mailbox, or personal identifiers;
- raw logs, command transcripts, local environment dumps, or local machine state;
- raw mailbox, document, invoice, image, PDF, or source data content;
- private file paths, session-cache paths, or absolute local paths unless intentionally public and
  necessary;
- private business strategy, commercial details, or sensitive operational context;
- security vulnerability details that need a private disclosure route.

Use short sanitized summaries and placeholders instead.

State: `draft recommendation`.

### D-11 - Triage And Promotion

Question: How do feedback issues become actual work?

Recommendation: Treat GitHub Issues as intake objects. During triage, choose one route:

- direct patch for clear small fixes;
- batch with related small issues;
- discovery when doctrine, ownership, behavior, or source of truth is unclear;
- implementation plan when accepted work changes managed kit content, schemas, validators,
  tooling, release surfaces, or cross-repo behavior;
- defer when valid but not currently important;
- duplicate, superseded, declined, or not actionable;
- route to Operating Kit feedback when the issue is actually generic kit work.

Accepted larger work should link the issue to the canonical discovery document, implementation
plan, PR, release note, or execution evidence. The plan register should record formal plans, not
every feedback issue.

State: `draft recommendation`.

## Requirements

## Functional Requirements

- `FR-1`: The Operating Kit provides a managed route for repo-specific feedback capture.
- `FR-2`: Agents can distinguish Operating Kit feedback from repo-owned feedback.
- `FR-3`: Agents can discover repo feedback mode from repo-local non-secret configuration.
- `FR-4`: Agents can offer feedback capture while executing another runbook without derailing it.
- `FR-5`: Agents can detect when GitHub Issues are already available, unavailable, disabled, or
  unconfigured.
- `FR-6`: Agents can offer an explicit issue-intake setup runbook on first real need only after
  checking whether GitHub Issues already works.
- `FR-7`: A user can decline setup and suppress repeated prompts.
- `FR-8`: Agents can draft sanitized issue content and ask before creating it.
- `FR-9`: The issue format captures affected repo, area, runbook or script, observed friction,
  expected behavior, sanitized evidence, workaround, proposed classification, and privacy
  confirmation.
- `FR-10`: Maintainers can triage issues into patch, batch, discovery, implementation planning,
  defer, close, or route-to-kit-feedback.

## Non-Functional Requirements

- `NFR-1`: Feedback capture must not expose secrets or sensitive operational evidence.
- `NFR-2`: The route must be low-interruption during active operational work.
- `NFR-3`: Issue creation and GitHub setup must require explicit approval.
- `NFR-4`: Declined setup must be durable enough that future agents do not keep asking.
- `NFR-5`: The mechanism must not require GitHub Issues for basic Operating Kit installation.
- `NFR-6`: Classification should be simple enough to reuse across repositories without large label
  setup overhead.
- `NFR-7`: The design must preserve plan-register boundaries by not registering every issue as a
  formal plan.

## Candidate Implementation Surfaces

Future implementation likely touches:

```text
components/agent-interface/managed/runbooks/capture-repo-feedback.md
components/agent-interface/managed/runbooks/enable-github-issues-feedback-intake.md
components/agent-interface/managed/reference/repo-feedback-capture.md
components/agent-interface/managed/reference/repo-feedback-item-format.md
components/agent-interface/managed/reference/root-agents-md-contract.md
components/agent-interface/managed/reference/operation-routing-and-dispatch.md
components/agent-interface/component.yaml
profiles/standard.yaml
schemas/kit-config.schema.json
src/codeheart_operating_kit/resources/...
tests/...
release-notes.md
```

Exact paths are implementation-planning candidates, not approved execution scope.

## Open Questions

| ID | Question | Blocking? | Recommendation |
| --- | --- | --- | --- |
| `OQ-1` | What exact `.codeheart/kit.config.yaml` schema should represent `repo_feedback` mode, destination, and suppression? | yes for implementation | Use `mode`, `destination`, `suppress_prompts`, and optional `disabled_reason`; finalize in implementation planning. |
| `OQ-2` | Should setup use `gh`, GitHub API through a connector, browser, or manual instructions? | yes for implementation | Prefer a runbook that detects available surfaces and asks before any external change. |
| `OQ-3` | Should labels and issue templates be standardized across all repos or only recommended? | no for discovery | Standardize a minimal label set, but allow repo-local extensions. |
| `OQ-4` | Should local draft mode exist? | no for discovery | Yes as an explicit opt-in mode, not as the default after GitHub setup is declined. |
| `OQ-5` | Should issue availability be detected live or trusted from config? | yes for implementation | Use config when present, otherwise check the GitHub remote and live issue availability before offering setup; use live checks as preflight before creating issues or setup changes. |
| `OQ-6` | Should security-sensitive feedback get a separate private disclosure route? | no for this scope | Keep out of v1 repo feedback unless a separate security disclosure discovery accepts it. |

## Risks

| Risk | Likelihood | Impact | Mitigation |
| --- | --- | --- | --- |
| Agents create too many low-value issues. | medium | medium | Use trigger and prompt timing rules; require user approval. |
| Agents keep asking to configure GitHub Issues after user declines. | high without config | medium | Record disabled/suppressed feedback state. |
| Sensitive evidence enters issues. | medium | high | Require sanitization fields and strict forbidden-content rules. |
| Operating Kit feedback and repo feedback get mixed. | medium | medium | Make destination ownership explicit in the route. |
| Label setup becomes too heavy for every repo. | medium | low | Keep required labels small and put area in issue body by default. |
| Onboarding becomes too complex. | low with recommendation | medium | Use demand-driven setup instead of onboarding setup. |

## Manual Review Packet

Recommended approval packet:

- Accept Operating Kit ownership of the generic repo feedback capture route.
- Accept demand-driven issue-intake setup instead of onboarding-time setup.
- Accept check-first behavior: verify whether GitHub Issues already works before offering setup.
- Accept GitHub Issues as the default durable inbox when configured or live preflight confirms it
  already works.
- Accept disabled/suppressed state when the user declines setup.
- Accept local draft mode only as explicit opt-in.
- Accept labels as canonical classification and title prefixes as readability helpers.
- Accept strict sanitization for private and public repos.
- Delegate exact config schema, setup execution surface, and validation details to
  implementation planning.

## Implementation Handoff Status

This discovery is implementation-handoff-ready for a draft implementation plan. The user approved
the recommended defaults on 2026-06-29, including check-first GitHub Issues behavior,
demand-driven setup, disabled/suppressed state after decline, labels as canonical classification,
local draft mode as explicit opt-in only, and strict sanitization.

## Implementation Capability Scope - Repo Feedback Capture And Issue Intake

Capability:
Agents in freshly installed and existing Operating Kit repositories can recognize repo-specific
feedback-worthy friction, check whether GitHub Issues already works, draft sanitized repo feedback
issues, ask before creating external issues or changing repo settings, offer setup only when the
issue surface is unavailable or incomplete, and suppress repeated setup prompts after user decline.

Primary workflow:
An agent encounters feedback-worthy friction while running another task or runbook. At a natural
checkpoint, blocker report, or final summary, the agent routes to the managed repo-feedback
runbook. The agent checks `.codeheart/kit.config.yaml` for `repo_feedback`. If feedback is
disabled, it stops prompting. If a GitHub Issues destination is configured, it preflights that
destination. If no destination is configured, it checks the current GitHub remote and live issue
availability. When Issues already works, the agent drafts a sanitized issue and asks before
creation. Missing standard labels or templates do not block capture; classification goes in the
body and only existing labels are used. When Issues are disabled, unavailable, or the user wants
standard labels/templates that are missing, the agent offers the setup runbook. If the user
declines setup, the agent asks to record disabled/suppressed state so future agents do not keep
asking.

Must cover:

- installed root route visibility for repo feedback capture;
- managed capture runbook with trigger model, prompt timing, destination resolution, check-first
  GitHub preflight, approval gates, issue drafting, label fallback, and stop conditions;
- managed setup runbook for enabling Issues, creating/verifying standard labels, and optionally
  adding a repo-owned issue template only after explicit approval;
- managed item-format or reference guidance for title prefixes, body fields, labels, privacy, and
  triage promotion;
- optional `repo_feedback` config schema that keeps missing config valid and interprets absence as
  auto-check behavior;
- disabled/suppressed config state after user decline;
- tests proving fresh installs do not need `repo_feedback` to validate and installed route targets
  exist;
- packaged-resource mirrors and component manifest updates;
- release notes and consumer-impact record for instruction, schema, and safety-policy changes;
- fresh low-context routing probe for the new route.

Explicitly out of scope:

- automatic GitHub issue creation;
- automatic GitHub Issues setup during Operating Kit onboarding;
- a new Operating Kit CLI feedback command;
- background telemetry or upload;
- mandatory local draft mode;
- forcing standard labels or issue templates before a repo feedback issue can be created;
- creating labels/templates across all Codeheart repositories in this implementation;
- a private security disclosure route.

Deferred or blocked:

- CLI-assisted issue drafting: deferred until manual issue capture proves the format;
- label/template sync automation: deferred until drift becomes repeated;
- private security disclosure: deferred to separate security disclosure discovery;
- GitHub Projects or dashboards: deferred until issue volume proves the need.

Preserve decisions:

- `D-1`: Operating Kit owns generic capture route; owning repo owns feedback issue and triage.
- `D-2`: All Operating Kit consumer repos are eligible, activated by config or first route use.
- `D-3`: GitHub Issues is default durable inbox when configured or live preflight confirms it.
- `D-4`: Setup is demand-driven and check-first, not onboarding-time.
- `D-5`: Setup is offered only after the check fails; decline records disabled/suppressed state.
- `D-6`: Agent triggers include blockers, runbook/script failures, workarounds, dissatisfaction,
  docs conflicts, undocumented overrides, and repeated friction.
- `D-7`: Prompt timing is blocker, dissatisfaction, checkpoint, or final summary.
- `D-8`: Issue creation, labels/templates, and GitHub settings require explicit approval.
- `D-9`: Labels are canonical classification; title prefixes aid readability.
- `D-10`: Sanitization is required even for private repos.
- `D-11`: Issues are intake; triage promotes to patch, batch, discovery, implementation plan,
  defer, duplicate, or close.

Planner must not reinvent:

- Missing `repo_feedback` config means unconfigured auto-check, not disabled.
- Missing standard labels or templates must not block issue capture when Issues works.
- The capture path may use only existing labels; otherwise classification belongs in the issue
  body.
- The setup runbook may use `gh repo view --json hasIssuesEnabled` and
  `gh repo edit --enable-issues` only after approval and tool readiness checks.
- Setup decline should default to disabled/suppressed state, not local draft mode.
- Local draft mode is explicit opt-in and remains ignored/local-user-owned.
- Feedback about the Operating Kit itself stays routed to `submit-kit-feedback.md`.

Feature-level success evidence:

- Fresh install test shows no `repo_feedback` block is written by default, config remains valid,
  and managed route targets for capture/setup exist.
- Config schema tests cover no block, configured GitHub Issues, disabled/suppressed, local draft
  opt-in, and invalid partial states.
- Packaged-resource fallback installs the capture/setup runbooks and reference files.
- Static grep or focused tests confirm the capture runbook contains check-first flow, label
  fallback, explicit approval gates, decline suppression, and Operating Kit feedback boundary.
- Read-only `gh repo view ... --json hasIssuesEnabled` proof is recorded for command shape.
- Fresh low-context routing probe shows an agent routes repo-specific runbook friction to repo
  feedback capture and does not confuse it with Operating Kit feedback.

## Revision Notes

- 2026-06-29: Created draft discovery from user discussion about repo-wide agent feedback capture,
  GitHub Issues as durable inbox, demand-driven setup, disabled suppression after decline,
  classification, privacy, and promotion to planning.
- 2026-06-29: Clarified check-first behavior from the Operating Kit feedback intake precedent:
  agents should verify whether GitHub Issues already works before offering setup, and should only
  offer setup when the issue surface is unavailable, disabled, or missing user-approved standard
  labels/templates.
- 2026-06-29: Added approved implementation capability scope after user accepted the recommended
  defaults and requested implementation planning.
