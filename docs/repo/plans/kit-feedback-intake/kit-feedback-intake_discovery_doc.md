Last updated: 2026-06-17T18:31:42Z (UTC)
Created: 2026-06-17
Status: active

# Document Header

## Overview

This discovery defines how Codeheart Operating Kit should collect consumer feedback, kit doctrine
gaps, sync issues, and product ideas without letting users or agents hand-edit managed kit content
inside consumer repositories.

The target is implementation-plan-ready. It records approved recommendations for the kit owner and
keeps only repository-governance preflight as a blocker before implementation planning. It does not
create an implementation plan yet and does not change issue forms, managed docs, CLI behavior, or
release assets.

Approved direction: make public GitHub Issues in `Codeheart-Operating-Kit` the canonical
shareable backlog, allow `.codeheart/user/feedback/` only as ignored sanitized draft space, add a
managed consumer-facing feedback runbook, add a maintainer triage runbook, add GitHub issue forms
with required public-core confirmations, include a rough-feedback issue form for not-yet-crystallized
ideas or friction, add minimal lifecycle labels, and defer a CLI `feedback` command until the
intake format has stabilized.

## Essential Context Reference Files

| Path | Reason |
| --- | --- |
| `AGENTS.md` | Public-core safety, maintainer read order, and change safety rules. |
| `README.md` | Repository purpose and public boundary. |
| `docs/README.md` | Current docs router; needs an index route for repo plans. |
| `docs/repo/README.md` | Repo docs router; needs a route to the new plan area. |
| `docs/repo/reference/placement-contract.md` | Defines `.codeheart/kit/`, `.codeheart/user/`, consumer-owned docs, and ownership modes. |
| `docs/repo/reference/consumer-impact-classification.md` | Defines impact classes for managed docs, scaffolds, templates, validators, and safety policy. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer procedure and stop conditions for kit changes. |
| `docs/repo/runbooks/promote-consumer-change.md` | Existing path for sanitized promotion of reusable consumer-local guidance. |
| `components/agent-interface/managed/reference/local-extension-contract.md` | Managed consumer guidance for repository-owned and local user sections. |
| `components/structure-governance/managed/reference/managed-content-boundaries.md` | Managed, scaffold, consumer-owned, local-user, generated, and report boundary doctrine. |
| `templates/user-layer/README.md` | Current ignored local user layer purpose and privacy constraints. |
| `src/codeheart_operating_kit/cli.py` | Current CLI command set; no feedback command exists. |
| `.github/workflows/validate.yml` | Existing CI validation surface. |

## Table Of Contents

- [Section 1 - Problem Framing](#section-1---problem-framing)
- [Section 2 - Context And Evidence](#section-2---context-and-evidence)
- [Section 3 - Requirements And Constraints](#section-3---requirements-and-constraints)
- [Section 4 - Decision Ledger](#section-4---decision-ledger)
- [Section 5 - Open Questions And Assumptions](#section-5---open-questions-and-assumptions)
- [Section 6 - Risks And Validation](#section-6---risks-and-validation)
- [Section 7 - Manual Review Packet](#section-7---manual-review-packet)
- [Revision Notes](#revision-notes)

# Section 1 - Problem Framing

## Problem Statement

Operating Kit consumers will notice missing doctrine, confusing routes, sync issues, stale
assumptions, and useful workflow ideas while working in consumer repositories. Today there is no
first-class intake path for those observations. Without a designed intake, users and agents may:

- edit managed `.codeheart/kit/` content locally and create drift;
- leave useful kit ideas in private chat context where maintainers cannot triage them;
- commit private or consumer-specific details into a public kit repository;
- create duplicate local backlog documents that are not visible to kit owners;
- skip the existing promotion path and implement consumer-specific fixes as generic doctrine.

## User Intention

Create a clean way for different users and consumers to record Operating Kit ideas and gaps without
messing up managed kit content, while giving kit owners a systematic backlog and triage path.

## Goals

- Give consumers and agents a safe intake route for kit feedback.
- Preserve `.codeheart/kit/` as managed content only.
- Keep private or rough feedback local until it is sanitized and intentionally submitted.
- Make shareable kit feedback visible to maintainers in one canonical backlog.
- Give maintainers a repeatable triage path from intake to discovery, implementation, release, and
  consumer sync.
- Keep v1 small enough to ship without overdesigning a second planning or issue system.

## Non-Goals

- Do not implement the feedback system in this discovery.
- Do not create a private security disclosure process in v1 unless the kit owner explicitly chooses
  that scope.
- Do not make `.codeheart/user/` a synchronized or canonical backlog.
- Do not let feedback drafts include secrets, credentials, tenant details, customer details, or
  raw local machine state in public outputs.
- Do not replace planning documents, release notes, or GitHub PR review with issue intake.
- Do not require a CLI command for the first shippable intake workflow.

## Priority Order

1. Public-core safety and privacy.
2. Managed-content boundary clarity.
3. Maintainer triage usefulness.
4. Low-friction consumer reporting.
5. Auditability from issue to release.
6. Minimal v1 implementation scope.
7. Future CLI automation compatibility.

## Durable Principles

- `.codeheart/kit/` is never a user feedback scratchpad.
- Public backlog records must be sanitized and public-core-safe.
- Local draft notes can exist, but they are not source of truth.
- GitHub Issues should hold shareable backlog state; planning docs should hold accepted discovery
  and implementation state.
- GitHub account friction is acceptable for canonical upstream feedback because it reduces
  low-quality submissions and keeps the first intake path maintainer-oriented.
- Maintainers, not consumers, decide when feedback becomes kit doctrine.
- Any new consumer-facing surface must have a consumer-impact classification.

## Success Criteria

This discovery is implementation-plan-ready when:

- candidate feedback domains are identified;
- all visible implementation-shaping decisions have concrete recommendations;
- privacy and managed-boundary risks are explicit;
- open questions are classified as blocking or non-blocking;
- a kit owner can verify GitHub issue governance during implementation preflight;
- the next step can be an implementation plan.

# Section 2 - Context And Evidence

## Candidate Domains And Decision Clusters

- Consumer intake workflow: how users and agents capture an observation safely.
- Public backlog authority: where shareable feedback becomes visible to kit maintainers.
- Privacy and public-core safety: what must be stripped, warned against, or kept local.
- Sensitive evidence and leak response: what local drafts may contain and how maintainers respond
  when private material is posted publicly.
- Managed-content boundary: what consumers may edit and what sync owns.
- Maintainer triage workflow: how issues become discovery, implementation, release, or closure.
- Implementation scope: docs and issue forms now versus CLI automation later.
- Validation and release impact: tests, public-core checks, release notes, and consumer migration
  notes.

## Evidence Summary

- The Operating Kit repo is public and forbids private business records, customer or tenant
  details, credentials, secrets, local machine state, account identifiers, restricted strategy, and
  raw operational logs.
- The placement contract defines `.codeheart/kit/` as managed content and `.codeheart/user/` as an
  ignored local user layer.
- Managed content boundaries state that user-specific or consumer-authored guidance must not go in
  `.codeheart/kit/`.
- The local user layer currently says it is for ignored local preferences and notes, not managed
  content, secrets, credentials, tokens, or public release material.
- The current CLI commands are `onboard`, `inspect`, `init`, `sync`, `check`, and `update-check`.
  No `feedback` command exists.
- `.github/` currently contains only CI workflow configuration. There are no issue templates.
- The existing `promote-consumer-change` runbook covers sanitized promotion of reusable
  consumer-local guidance, but it is maintainer-facing and not a consumer feedback intake route.
- Existing validation covers tests, public-core hygiene, Markdown timestamps, JSON schemas, and
  release manifests. There is no issue-template-specific validator.
- Issue availability and repository-governance permissions must be verified before implementation;
  this discovery does not assume repo settings changes are already approved.

# Section 3 - Requirements And Constraints

## Functional Requirements

- `FR-1`: Consumers and agents can report Operating Kit feedback without editing
  `.codeheart/kit/`.
- `FR-2`: The workflow supports local sanitized draft capture and public sanitized submission as
  separate steps.
- `FR-3`: Shareable kit feedback has one canonical backlog visible to kit maintainers.
- `FR-4`: Feedback records capture kit version, installed component or workflow area, consumer
  context, observed problem, expected behavior, evidence, privacy status, and proposed
  classification.
- `FR-5`: Maintainers can triage feedback into declined, duplicate, needs information, accepted
  backlog, needs discovery, implementation planned, released, or consumer-specific.
- `FR-6`: Consumers can distinguish generic kit feedback from consumer-specific local exceptions.
- `FR-7`: The v1 workflow is discoverable from managed consumer docs.
- `FR-8`: The maintainer workflow is discoverable from Operating Kit repo docs.
- `FR-9`: The workflow can later support a CLI-assisted feedback draft without changing the
  canonical authority model.
- `FR-10`: Maintainers have an immediate response path when sensitive material is accidentally
  posted publicly.

## Non-Functional Requirements

- `NFR-1`: Public-core safety: public issue bodies and examples must not encourage disclosure of
  secrets, customer details, tenant details, account identifiers, private strategy, or raw local
  logs.
- `NFR-2`: Boundary safety: sync must not overwrite consumer-owned feedback records or local user
  drafts.
- `NFR-3`: Auditability: accepted feedback should trace to discovery, implementation, release
  notes, and synced consumer behavior.
- `NFR-4`: Low friction: a user with a concern should not need to understand every kit component to
  report it.
- `NFR-5`: Maintainer usefulness: issue forms should collect enough information to reproduce,
  classify, and route feedback.
- `NFR-6`: Backwards compatibility: v1 should not require existing consumers to migrate.
- `NFR-7`: Extensibility: a later CLI command should be able to read the same intake format.
- `NFR-8`: Local draft safety: ignored draft space is not a safe place for secrets, raw logs,
  customer or tenant details, account identifiers, token values, or real local machine state.

## Operating Scenarios

### OS-1 - Consumer Notices A Managed Doctrine Gap

An agent working in a consumer repo sees that a managed runbook is missing guidance. It records a
local sanitized note or opens a GitHub issue, then the kit maintainer triages it.

### OS-2 - Consumer Has Private Evidence

A user has evidence that includes local paths, client context, or logs. The workflow tells them not
to store raw sensitive evidence in feedback drafts, to keep private evidence outside the kit
workflow, to write a sanitized summary or pointer, and to include only public-safe reproduction
details in public issues.

### OS-3 - Maintainer Converts Feedback Into Work

A maintainer reviews an issue, labels it, requests missing info, closes it as consumer-specific, or
opens a discovery or implementation plan when the feedback changes kit doctrine or behavior.

### OS-4 - Future CLI Draft

A later `codeheart-operating-kit feedback new` command can read lockfile metadata, create a local
draft, warn about secrets, and optionally print a GitHub issue body. It should not auto-submit raw
local drafts in v1.

# Section 4 - Decision Ledger

## Decision Inventory

| ID | Decision | Class | State | Owner | Depends On | Blocks |
| --- | --- | --- | --- | --- | --- | --- |
| `D-1` | Canonical backlog authority | blocking | approved with implementation preflight | Kit owner | `D-5`, `D-8`, GitHub Issues preflight | `FR-3`, `FR-5`, implementation plan |
| `D-2` | Consumer-local capture boundary | implementation-shaping | approved | Kit owner | `D-1` | `FR-1`, `FR-2`, `NFR-2` |
| `D-3` | First shippable scope | implementation-shaping | approved | Kit owner | `D-1`, `D-2`, `D-5`, `D-8` | implementation plan shape |
| `D-4` | Feedback taxonomy and issue form set | implementation-shaping | approved | Kit owner | `D-1`, `D-3` | `FR-4`, `FR-5` |
| `D-5` | Privacy and sanitization posture | blocking | approved | Kit owner | None | public issue design |
| `D-6` | Maintainer triage lifecycle | implementation-shaping | approved | Kit owner | `D-1`, `D-4`, `D-5`, `D-8` | maintainer runbook |
| `D-7` | Consumer impact and validation class | implementation-shaping | approved | Kit owner | `D-3`, `D-5`, `D-8` | release planning |
| `D-8` | Sensitive evidence and accidental disclosure response | blocking | approved | Kit owner | `D-5` | public issue forms, maintainer triage |

## D-1 - Canonical Backlog Authority

Question: Where should shareable Operating Kit feedback become canonical?

Options:

- Public GitHub Issues in `Codeheart-Operating-Kit`.
- A repository-owned Markdown backlog under `docs/repo/`.
- Consumer-local `.codeheart/user/feedback/` records.
- External private task system.

Recommendation: use public GitHub Issues in `Codeheart-Operating-Kit` as the canonical shareable
backlog after an implementation preflight verifies that Issues are enabled and issue forms can be
added without an unapproved public repository settings change. Use planning docs only after
maintainers accept work that needs discovery or an implementation plan. Use
`.codeheart/user/feedback/` only as optional ignored sanitized draft space.

Rationale: GitHub Issues are visible to maintainers, are designed for backlog state, avoid
overloading planning documents, and are already adjacent to PR/release workflow. Markdown backlog
files would create a second manually curated issue system. Local drafts are useful before
sanitization, but they are invisible to maintainers and must not become source of truth.

State: `approved with implementation preflight`.

Reviewer summary: reviewer agreed with GitHub Issues as the right authority, but challenged the
missing dependencies on privacy posture and issue availability. The recommendation now includes
preflight verification and depends on the tightened sensitive-evidence decisions.

Next user action: none before implementation planning. Implementation planning must verify GitHub
Issues, issue forms, and label governance before editing repository governance files.

## D-2 - Consumer-Local Capture Boundary

Question: What may consumers write locally before submitting upstream feedback?

Options:

- No local draft surface; users open issues directly.
- Optional ignored `.codeheart/user/feedback/` drafts.
- Consumer-owned `docs/repo/backlog/`.
- Managed `.codeheart/kit/feedback/`.

Recommendation: allow optional ignored `.codeheart/user/feedback/` drafts for sanitized summaries,
public-safe reproduction notes, and pointers to private evidence that remains outside the feedback
workflow. Explicitly forbid managed `.codeheart/kit/` feedback edits. Do not create a committed
consumer `docs/repo/backlog/` as a generic kit feature.

Rationale: `.codeheart/user/` already exists for local user notes and is ignored/local-only, but
existing local-user rules still forbid secrets and private machine state. Managed kit content must
stay synchronized from releases. A committed consumer backlog would mix consumer-specific work with
generic kit governance and would be harder to triage upstream.

State: `approved`.

Reviewer summary: reviewer agreed with the boundary but flagged that "ignored local draft" could be
misread as safe raw-evidence storage. The recommendation now limits drafts to sanitized summaries
and pointers, not secrets, raw logs, tenant/customer details, account identifiers, or machine
state.

Next user action: none before implementation planning.

## D-3 - First Shippable Scope

Question: What should v1 implement?

Options:

- Docs and issue forms only.
- Docs, issue forms, and a local draft scaffold.
- Docs, issue forms, maintainer triage runbook, consumer managed runbook, feedback item
  format, and no CLI.
- Full CLI-assisted feedback command.

Recommendation: v1 should include GitHub issue forms, a managed consumer-facing feedback runbook, a
feedback item format reference, a maintainer triage runbook, docs index updates, minimal lifecycle
labels, release notes, and no CLI command. Defer CLI until real intake examples prove the format.

Rationale: the main risk is boundary and privacy confusion, not missing automation. Docs and
templates can establish the authority model with lower blast radius. A CLI command too early may
freeze the wrong schema and require tests, packaging, and UX design before the intake model is
validated.

State: `approved`.

Reviewer summary: reviewer agreed with deferring CLI and asked to make sanitization a dependency.
The recommendation now depends on `D-5` and `D-8` before issue forms or runbooks ship.

Next user action: none before implementation planning.

## D-4 - Feedback Taxonomy And Issue Form Set

Question: Which feedback types should v1 support?

Options:

- One generic feedback issue form.
- Separate forms for rough feedback, bug, doctrine gap, sync/install issue, docs clarity, and
  feature request.
- A larger taxonomy including security reports, consumer exceptions, and release requests.

Recommendation: use a small set of public GitHub issue forms with required fields and required
public-core confirmation checkboxes:

- rough feedback or idea;
- kit bug or regression;
- doctrine or workflow gap;
- install, sync, or check issue;
- docs clarity or routing issue;
- feature or capability request.

Do not add a public security-report template in v1 unless a private disclosure channel is also
defined. Each public issue form should contain strong public-core warnings and required
confirmation that the submission excludes sensitive content.

Rationale: one generic form is too weak for triage. A rough-feedback form lets users submit
not-yet-crystallized friction, ideas, and weak signals without pretending they are bugs or
implementation-ready requests. Plain Markdown templates cannot enforce required confirmations as
well as issue forms. A larger taxonomy is likely premature and may invite sensitive reports into a
public repository.

State: `approved`.

Reviewer summary: reviewer agreed the taxonomy is plausible and identified issue forms as a
missing implementation-shaping choice. The recommendation now chooses issue forms with required
public-core confirmations. Follow-up review added a rough-feedback form and explicitly accepted
GitHub account friction for canonical upstream feedback.

Next user action: none before implementation planning.

## D-5 - Privacy And Sanitization Posture

Question: How should feedback intake prevent private data leakage?

Options:

- Trust users to sanitize.
- Put warnings in templates only.
- Make sanitization a first-class field and workflow step.
- Avoid public issues entirely.

Recommendation: make sanitization a first-class field and workflow step. Public issue forms must
require reporters to confirm that the issue excludes secrets, credentials, customer or tenant
details, local machine state, account identifiers, raw logs, and private strategy. Local draft
space must use sanitized summaries or pointers only; raw sensitive evidence should stay outside the
feedback workflow and should not be copied into `.codeheart/user/feedback/`.

Rationale: the kit repository is public and already has a strict public-core boundary. Warnings
alone are easy to miss. Avoiding public issues entirely would reduce maintainer visibility and is
not required if templates and runbooks are explicit.

State: `approved`.

Reviewer summary: reviewer agreed that the decision is correctly blocking, but found public-body
guidance insufficient without local-draft and leak-response rules. The recommendation now covers
local draft content and is paired with `D-8`.

Next user action: none before implementation planning. Private disclosure remains deferred to a
later security-specific plan unless the implementation scope is explicitly expanded.

## D-6 - Maintainer Triage Lifecycle

Question: How should kit owners process incoming feedback?

Options:

- Handle issues ad hoc.
- Use labels only.
- Add a maintainer runbook with statuses and conversion rules.

Recommendation: add a maintainer triage runbook that classifies issues as duplicate, declined,
needs information, consumer-specific, accepted backlog, needs discovery, implementation planned,
released, or superseded. Represent lifecycle state with minimal GitHub labels and maintainer
triage comments in v1; include a `needs-shaping` label for rough-feedback or idea submissions
that are not implementation-ready; do not require a project board. Accepted doctrine or behavior
changes should create or update discovery and implementation plans when they are non-trivial.

Rationale: labels alone do not explain when feedback becomes doctrine, when a plan is required, or
how release notes and consumer impact are recorded. Comments alone are too hard to scan. Minimal
labels plus a runbook keep maintainer behavior consistent without requiring a new project-management
tool in v1.

State: `approved`.

Reviewer summary: reviewer agreed that lifecycle states are useful and asked for a machine-readable
or scan-friendly representation plus sensitive-data response dependency. The recommendation now
uses minimal labels plus comments and depends on `D-8`.

Next user action: none before implementation planning.

## D-7 - Consumer Impact And Validation Class

Question: How should v1 be classified and validated?

Options:

- Treat v1 as documentation-only and skip consumer impact.
- Classify managed docs as instruction-only and issue forms as repository governance.
- Add a scaffold or CLI and classify as backwards-compatible scaffold or behavior change.

Recommendation: classify v1 as both:

- `instruction-only change` for managed consumer-facing docs and maintainer docs; and
- `security or safety policy change` for public-core sanitization, sensitive-evidence handling,
  and leak-response rules.

Also record a repository-governance addition for GitHub issue forms and labels. If
`.codeheart/user/feedback/` is scaffolded or CLI behavior is added later, classify that later work
separately.

Validation should include Markdown timestamp validation, public-core validation, tests affected by
resource packaging if managed docs are added, release manifest validation, explicit safety-policy
review, and issue-form review for required public-core confirmations unless an issue-form validator
is added.

State: `approved`.

Reviewer summary: reviewer found the initial classification underpowered because sanitization and
secrets handling affect safety policy. The recommendation now includes both instruction-only and
security/safety policy classifications with explicit release-note and review expectations.

Next user action: none before implementation planning.

## D-8 - Sensitive Evidence And Accidental Disclosure Response

Question: What may local drafts contain, and how should maintainers respond if sensitive material
is posted publicly?

Options:

- Allow raw evidence in `.codeheart/user/feedback/` because it is ignored.
- Allow sanitized summaries and private-evidence pointers only.
- Ban local feedback drafts entirely.
- Define a full private disclosure program in v1.

Recommendation: allow only sanitized summaries and pointers in `.codeheart/user/feedback/`. Do not
store secrets, credentials, token values, raw logs, customer or tenant details, account
identifiers, private strategy, or real local machine state in feedback drafts. Public issue forms
must require a sanitization confirmation.

Maintainer triage must include an accidental public disclosure response:

1. Do not copy or quote the sensitive material.
2. Use available GitHub moderation tools to hide, edit, delete, or minimize the public exposure.
3. Ask the reporter to rotate exposed secrets or credentials when relevant.
4. Preserve only a sanitized summary for triage.
5. Escalate privately to the kit owner when the exposure cannot be handled through normal issue
   moderation.

Rationale: ignored local files are not a security boundary. A public intake feature needs both
reporter-side prevention and maintainer-side response. A full private disclosure program is
important but larger than the focused v1 feedback intake.

State: `approved`.

Reviewer summary: reviewer identified this as a missing blocking decision. The decision was added
and now controls local draft content, public issue confirmation, and maintainer leak response.

Next user action: none before implementation planning.

# Section 5 - Open Questions And Assumptions

## Open Questions

### OQ-1 - GitHub Issues As Canonical Shareable Backlog (Resolved)

- Owner: Kit owner
- `BLOCKER: no`
- Affected IDs: `D-1`, `FR-3`, `FR-5`
- Question: Should public GitHub Issues in `Codeheart-Operating-Kit` be the canonical shareable
  backlog for sanitized kit feedback?
- Recommended default: Yes.
- Resolution path: Resolved by user approval on 2026-06-17. Implementation preflight must still
  verify that GitHub Issues are enabled and issue forms can be added without an unapproved public
  repository settings change.

### OQ-2 - Private Disclosure Channel

- Owner: Kit owner
- `BLOCKER: no`
- Affected IDs: `D-4`, `D-5`, `R-1`
- Question: Should this feature define a private security or sensitive-feedback channel, or should
  v1 only warn users not to post sensitive reports publicly?
- Recommended default: Defer private disclosure to a focused security reporting plan; v1 should
  warn users not to post sensitive reports publicly.
- Resolution path: Decide before adding any public security issue template.

### OQ-3 - CLI Timing

- Owner: Kit owner
- `BLOCKER: no`
- Affected IDs: `D-3`, `FR-9`
- Question: Should `codeheart-operating-kit feedback new` be included in v1?
- Recommended default: No. Defer until issue templates and maintainer triage have produced real
  examples.
- Resolution path: Approve v1 docs/forms first or explicitly expand implementation scope.

### OQ-4 - Label And Project Board Governance

- Owner: Kit owner
- `BLOCKER: no`
- Affected IDs: `D-4`, `D-6`
- Question: Should v1 create GitHub labels or a project board, or only issue forms and a
  triage runbook?
- Recommended default: Create minimal lifecycle labels and a triage runbook. Do not create a
  project board in v1.
- Resolution path: Minimal lifecycle labels are approved. Project board remains deferred unless
  future triage volume proves it is needed.

### OQ-5 - GitHub Issues Availability And Governance

- Owner: Kit maintainer
- `BLOCKER: yes`
- Affected IDs: `D-1`, `D-4`, `D-6`
- Question: Are GitHub Issues enabled for `Codeheart-Operating-Kit`, and can issue forms and labels
  be added without an unapproved public repository settings change?
- Recommended default: Verify in implementation preflight before editing `.github/ISSUE_TEMPLATE/`.
- Resolution path: Use GitHub repository metadata or maintainer confirmation during implementation
  planning.

## Assumptions

- `A-1`: The Operating Kit repository remains public.
- `A-2`: GitHub Issues are expected to be available for the public repository, but implementation
  must verify this before adding issue forms.
- `A-3`: Consumers may have private context, so public templates must be defensive by default.
- `A-4`: v1 should minimize release and migration risk.
- `A-5`: The existing `promote-consumer-change` runbook remains useful after feedback is accepted
  and sanitized.

# Section 6 - Risks And Validation

## Risks

### R-1 - Sensitive Data Leakage

- Likelihood: medium
- Impact: high
- Mitigation: first-class sanitization fields, public-core warnings, required issue-form
  confirmations, sanitized local draft rules, no public security template without a private
  channel, no auto-submit from local drafts, and maintainer leak-response steps.
- Detection: public-core review, issue-form review, safety-policy review, maintainer triage.

### R-2 - Local Drafts Become A Shadow Backlog

- Likelihood: medium
- Impact: medium
- Mitigation: state that `.codeheart/user/feedback/` is optional ignored draft space only and that
  GitHub Issues are canonical for shareable feedback.
- Detection: managed runbook wording and maintainer triage checks.

### R-3 - Overbuilt CLI Freezes The Wrong Schema

- Likelihood: medium
- Impact: medium
- Mitigation: defer CLI until templates produce real examples.
- Detection: implementation plan keeps CLI out of v1.

### R-4 - Feedback Bypasses Consumer Impact Classification

- Likelihood: medium
- Impact: medium
- Mitigation: maintainer triage runbook requires impact classification before accepted changes.
- Detection: PR and release review.

### R-5 - Issue Templates Are Too Heavy

- Likelihood: low
- Impact: medium
- Mitigation: keep template fields minimal and route complex items to maintainer follow-up.
- Detection: first few issue triage cycles.

## Targeted Validation Plan

- Run Markdown timestamp validation.
- Run public-core validation.
- Run affected tests for managed resource packaging when adding managed docs.
- Run release-manifest validation when managed content changes.
- Inspect issue forms for required public-core confirmations and useful triage fields.
- Verify issue availability and issue-form support before implementation edits.
- Verify maintainer triage includes accidental public disclosure response.
- Verify release notes flag the `security or safety policy change` impact.
- Verify no new consumer-owned path is overwritten by sync.
- Verify release notes classify the consumer impact.

# Section 7 - Manual Review Packet

## Approved Recommendations

- `D-1`: public GitHub Issues in `Codeheart-Operating-Kit` are the canonical shareable backlog for
  sanitized feedback, pending implementation preflight for issue availability and governance.
- `D-2`: `.codeheart/user/feedback/` may be optional ignored local sanitized draft space, while
  `.codeheart/kit/` must never be edited for feedback.
- `D-3`: v1 ships docs, issue forms, feedback item format, maintainer triage runbook, and minimal
  lifecycle labels; CLI feedback command is deferred.
- `D-4`: v1 uses six public issue forms: rough feedback or idea, bug/regression, doctrine/workflow
  gap, install/sync/check issue, docs clarity/routing issue, and feature/capability request.
- `D-5`: sanitization is a first-class required workflow step for public submission.
- `D-6`: maintainer triage uses explicit lifecycle states represented by minimal labels and
  comments, includes a `needs-shaping` label for rough feedback, and promotes non-trivial accepted
  changes into discovery or implementation plans.
- `D-7`: v1 is both an instruction-only managed-doc change and a security/safety policy change,
  plus a repository-governance issue-form and label addition; CLI or scaffold changes are deferred
  and classified later.
- `D-8`: local feedback drafts may contain sanitized summaries and pointers only, and maintainer
  triage must include no-copy leak-response steps.

## Manual Choice Points

- `OQ-5`: Verify GitHub Issues availability and issue-form governance before implementation.

## Non-Blocking Follow-Ups

- Decide whether a private security reporting path belongs in a separate security intake plan.
- Decide whether a GitHub project board should be added after issue forms and labels exist.
- Revisit a CLI-assisted draft command after real issue examples prove the data model.

## Exact User Actions Needed

1. Verify GitHub Issues, issue forms, and label governance during implementation planning.
2. Keep private security reporting out of v1 unless implementation scope is explicitly expanded.
3. Request an implementation plan for `kit-feedback-intake`.

# Revision Notes

- 2026-06-17: Initial deep discovery created for Operating Kit feedback intake, backlog authority,
  privacy posture, consumer-local draft boundaries, maintainer triage, and v1 scope.
- 2026-06-17: Integrated reviewer critique by adding sensitive-evidence and leak-response decision,
  issue-form recommendation, minimal lifecycle labels, GitHub Issues preflight, and corrected
  security/safety impact classification.
- 2026-06-17: Recorded user approval for the GitHub Issues intake model, v1 scope, sanitized draft
  boundary, privacy posture, minimal lifecycle labels, security/safety impact classification, and
  rough-feedback issue form.
