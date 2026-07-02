Last updated: 2026-07-02T13:16:41Z (UTC)

# Enable GitHub Issues Feedback Intake

Use this runbook when a verified Codeheart maintainer wants to enable or standardize GitHub
Issues intake for repo feedback.

Audience: maintainer-facing

Intent:
Guide a verified Codeheart maintainer through approval-gated GitHub Issues, label, template, and
repo-feedback config setup for one repository.

Success:
The repository has an approved feedback intake path, or setup is declined and recorded without
repeated prompts.

Agent judgment boundary:
The agent may run read-only authorization and repository preflight checks. The agent must not
enable Issues, create labels, create templates, write config, install `gh`, repair auth, or use a
browser/manual fallback without the required approval and preconditions.

Stop boundary:
Stop setup when `gh` is missing, auth fails, Codeheart organization membership cannot be verified,
the target repository is unclear, or the maintainer does not approve an external or repo-owned
change.

## Required Inputs

- Target repository.
- Feedback trigger summary.
- Desired setup level: Issues only, labels, issue template, config recording, or decline record.
- Authorization-gate evidence from `capture-repo-feedback.md`.
- Privacy confirmation for any feedback example or template text.

## Authorization Preflight

Use existing `gh` only:

```sh
gh auth status
gh api user/memberships/orgs/Codeheart-Digital-Solutions --jq .state
```

Continue only when organization membership state is `active`.

If `gh` is absent, authentication fails, or membership is unverifiable, stop setup. Do not route
to tooling readiness, offer install or repair, open a browser, provide manual fallback, create a
local draft, or propose another feedback mechanism.

## Repository Preflight

Resolve the target repository and inspect issue state:

```sh
git remote get-url origin
gh repo view <owner/repo> --json nameWithOwner,hasIssuesEnabled,viewerPermission,isPrivate
```

Record the sanitized result. Do not change repository settings during preflight.

## Enable Issues Approval

If Issues are disabled and the maintainer wants GitHub Issues feedback intake, present an
approval packet:

- target repository;
- current issue state;
- command to run;
- effect: enables GitHub Issues for the repository;
- validation command.

Run only after explicit approval:

```sh
gh repo edit <owner/repo> --enable-issues
```

Validate with `gh repo view <owner/repo> --json hasIssuesEnabled`.

## Standard Labels

Recommended labels:

- `feedback-intake`
- `kind-runbook-gap`
- `kind-tooling-gap`
- `kind-docs-conflict`
- `kind-bug`
- `kind-feature`
- `source-agent-runbook`
- `source-user-feedback`
- `privacy-sanitized`
- `needs-triage`

Inventory existing labels:

```sh
gh label list --repo <owner/repo> --limit 200 --json name
```

Create missing labels only after explicit approval. Use concise descriptions and stable colors:

```sh
gh label create <name> --repo <owner/repo> --description <description> --color <hex>
```

Missing labels never block feedback issue capture.

## Issue Template

Creating `.github/ISSUE_TEMPLATE/repo-feedback.yml` is a repo-owned source change. Do it only
after explicit approval in that repository. The template must enforce the fields and forbidden
content from `reference/repo-feedback-item-format.md`.

Validate template edits with `git diff --check` and the repository's normal review path.

## Config Recording

When setup is approved, record repo feedback config only after approval:

```yaml
repo_feedback:
  mode: github_issues
  destination:
    type: github_issues
    owner: <owner>
    repo: <repo>
  authorization:
    organization: Codeheart-Digital-Solutions
    require_verified_membership: true
    require_gh_cli: true
    unavailable_behavior: silent
  github_standardization:
    labels: configured | not_configured | declined
    issue_templates: configured | not_configured | declined
```

When setup is declined and the maintainer wants prompts suppressed, record:

```yaml
repo_feedback:
  mode: disabled
  disabled_reason: verified_maintainer_declined_issue_intake
```

`mode: disabled` is the suppression state. Do not add `suppress_prompts`.

## Validation

After approved setup:

1. Re-run `gh repo view <owner/repo> --json nameWithOwner,hasIssuesEnabled`.
2. Re-run label inventory if labels changed.
3. Run `git diff --check` for repo-owned template or config changes.
4. Record the approved action summary and validation result.

## Structured Recipe Metadata

Recipe ID: `enable-github-issues-feedback-intake`
Purpose: approval-gated setup for repo feedback GitHub Issues intake.
Inputs: target repository, setup level, authorization-gate evidence, privacy confirmation.
Preconditions: existing authenticated `gh`, active `Codeheart-Digital-Solutions` membership,
clear target repository, explicit approval for writes.
Approval class: read-only preflight allowed; repository settings, labels, templates, and config
writes require explicit approval.
Execution surface: existing `gh` and repo-owned file edits only.
Evidence output: authorization gate result without sensitive auth data, issue availability
preflight, approval packet, approved action summary, config diff summary, validation summary.
Validation: second repo view, label list when applicable, `git diff --check` when files changed.
Stop conditions: failed authorization gate, unclear target, missing approval, unavailable
permission, or sensitive disclosure.
Maturity: L1 structured recipe. Do not promote this workflow to a reusable script, CLI command,
wrapper, API, durable helper, package, or tool surface in v1.
