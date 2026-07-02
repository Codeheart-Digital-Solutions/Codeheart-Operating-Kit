Last updated: 2026-07-02T13:16:41Z (UTC)

# Capture Repo Feedback

Use this runbook when an agent notices repository-owned friction while doing another task and
the observation may be worth durable maintainer triage.

Audience: agent-facing

Intent:
Capture sanitized repo-specific feedback only when the current session is a verified Codeheart
GitHub organization member session and the owning repository has an available GitHub Issues
surface.

Success:
The immediate task is not derailed, the feedback is routed to the owning repository when allowed,
and no user-facing feedback prompt appears when the authorization gate fails.

Agent judgment boundary:
The agent may identify likely feedback-worthy friction and draft sanitized issue content after
the authorization gate passes. The agent must not create issues, configure GitHub, install or
repair `gh`, use browser/manual fallback, write local drafts, or route to tooling readiness when
the authorization gate fails.

Stop boundary:
Stop silently when `gh` is missing, GitHub CLI authentication is missing, or active membership in
`Codeheart-Digital-Solutions` cannot be positively verified. Stop and ask before issue creation,
GitHub settings changes, label creation, issue-template creation, or repo config changes.

## When To Consider Capture

Detect repo feedback candidates when one of these happens:

- a runbook step fails, is ambiguous, or needs a workaround;
- a script block fails, produces fragile output, or needs manual repair;
- docs and implementation disagree;
- local tooling readiness blocks or materially delays the task;
- the user says an outcome is confusing, too dense, wrong, unsatisfactory, or not what they meant;
- the agent intentionally overrides a documented default;
- setup or onboarding guidance is missing;
- a recovered detour is likely to recur;
- the same friction appears repeatedly.

Do not offer capture for trivial typos, one-off preference, sensitive incident details, or notes
that belong only in the current task summary.

## Authorization Gate

Run this gate before mentioning repo feedback capture to the user:

1. Check whether `gh` already exists.
2. Check existing GitHub CLI authentication:
   `gh auth status`
3. Verify active Codeheart organization membership:
   `gh api user/memberships/orgs/Codeheart-Digital-Solutions --jq .state`
4. Continue only when the state is `active`.

If any step fails, stop the repo feedback route silently. Do not offer to install `gh`, repair
auth, open a browser, provide manual instructions, create local drafts, route to tooling
readiness, or propose another feedback mechanism. If the user explicitly asks why feedback
capture is unavailable, give a short explanation without prompting setup.

## Destination Resolution

After the authorization gate passes:

1. Decide whether the observation is Operating Kit feedback or repo-owned feedback.
2. For generic Operating Kit doctrine, managed content, install, sync, or routing issues, use
   `runbooks/submit-kit-feedback.md` instead.
3. Inspect optional `repo_feedback` config.
4. If `repo_feedback.mode: disabled`, do not prompt again unless the verified maintainer asks to
   re-enable repo feedback.
5. If `repo_feedback.mode: github_issues` is configured, use
   `repo_feedback.destination.owner` and `repo_feedback.destination.repo` as the target
   repository.
6. If config is missing, treat it as authorization-gated auto-check and resolve the GitHub remote
   for the owning repository.
7. If neither a configured destination nor a GitHub remote can identify one target repository,
   stop the repo feedback route silently unless the verified maintainer explicitly asks for setup.
8. Check issue availability:
   `gh repo view <owner/repo> --json nameWithOwner,hasIssuesEnabled,viewerPermission,isPrivate`
9. If Issues are enabled and reachable, draft sanitized feedback and ask before creating it.
10. If Issues are disabled, unavailable, or the verified maintainer asks for missing label or
   template standardization, route to `enable-github-issues-feedback-intake.md`.

Missing standard labels or issue templates do not block issue capture when Issues works.

## Draft The Issue

Use `reference/repo-feedback-item-format.md`.

The issue body must include:

- summary;
- owning repository;
- affected area;
- feedback kind;
- observed problem;
- expected behavior;
- sanitized evidence;
- workaround used;
- proposed classification;
- privacy confirmation.

Do not include forbidden content from the item-format reference.

## Label Fallback

List existing labels before using labels:

```sh
gh label list --repo <owner/repo> --limit 200 --json name
```

Use only labels that already exist. If recommended labels are missing, put classification in the
issue body and continue.

## Approval Gate

Before creating an issue, show the verified maintainer:

- target repository;
- sanitized title;
- sanitized body summary;
- labels that will be used;
- confirmation that forbidden content was removed.

Create the issue only after explicit approval:

```sh
gh issue create --repo <owner/repo> --title <title> --body-file <body-file>
```

## Decline Suppression

If the verified maintainer declines repo feedback setup for this repository, ask before writing
repo config. When approved, record:

```yaml
repo_feedback:
  mode: disabled
  disabled_reason: verified_maintainer_declined_issue_intake
```

`mode: disabled` is the suppression state. Do not add a separate `suppress_prompts` field.

## Structured Recipe Metadata

Recipe ID: `capture-repo-feedback`
Purpose: capture sanitized repo-owned feedback without burdening unauthorized or non-technical
users.
Inputs: feedback candidate, owning repository, optional configured `repo_feedback.destination`,
GitHub remote when no destination is configured, existing GitHub CLI auth.
Preconditions: existing `gh`, authenticated GitHub CLI, active
`Codeheart-Digital-Solutions` membership, public-safe feedback candidate.
Approval class: read-only checks allowed after route selection; issue creation and config writes
require explicit approval.
Execution surface: existing `gh` only.
Evidence output: authorization gate result without sensitive auth data, issue availability
summary, sanitized draft, labels used, approval text, issue URL or suppression state.
Validation: issue URL after approved creation, or recorded no-op/suppression state.
Stop conditions: failed authorization gate, sensitive disclosure, unavailable issue surface
without setup approval, or missing approval for external changes.
Maturity: L1 structured recipe. Do not promote this workflow to a reusable script, CLI command,
wrapper, API, durable helper, package, or tool surface in v1.
