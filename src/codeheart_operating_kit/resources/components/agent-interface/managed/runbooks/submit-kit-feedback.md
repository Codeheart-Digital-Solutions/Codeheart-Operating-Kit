Last updated: 2026-06-17T19:34:19Z (UTC)

# Submit Kit Feedback

Use this runbook when a consumer, user, or agent notices a reusable Operating Kit issue, doctrine
gap, sync problem, routing confusion, or product idea.

## Boundary

Do not edit `.codeheart/kit/` to record feedback. Files under `.codeheart/kit/` are managed
Operating Kit content and should be refreshed through kit sync or repair workflows.

Use `.codeheart/user/feedback/` only as optional ignored local draft space for sanitized notes.
Local drafts are not source of truth and are not a safe place for secrets, credentials, customer or
tenant details, local machine state, account identifiers, raw logs, or private strategy.

Public, shareable kit feedback belongs in GitHub Issues for the public
`Codeheart-Operating-Kit` repository after it is sanitized.

## Stop Conditions

Do not submit security vulnerabilities, sensitive disclosures, secrets, credentials, customer or
tenant details, local machine state, account identifiers, raw logs, or private strategy through
public GitHub Issues. Keep that material out of the feedback workflow.

## Submit Feedback

1. Decide whether the observation is generic kit feedback or a consumer-specific local exception.
2. Keep consumer-specific guidance in the consumer repository or workspace.
3. Sanitize the feedback before writing it into a draft or public issue.
4. Use `.codeheart/user/feedback/` only for optional sanitized drafts or pointers to private
   evidence that remains outside the kit workflow.
5. Open the issue chooser:
   `https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/issues/new/choose`
6. Use the matching public GitHub issue form:
   - rough feedback or idea: `rough-feedback.yml`;
   - kit bug or regression: `kit-bug.yml`;
   - doctrine or workflow gap: `doctrine-workflow-gap.yml`;
   - install, sync, or check issue: `install-sync-check.yml`;
   - docs clarity or routing issue: `docs-routing.yml`;
   - feature or capability request: `feature-request.yml`.
7. Include the kit version, affected area, observed problem, expected behavior, sanitized evidence
   summary, privacy confirmation, and proposed classification when known.

## Sanitization Rules

Before submitting or drafting feedback, remove:

- secrets and credentials;
- customer or tenant details;
- account identifiers;
- raw local paths, machine state, or environment dumps;
- raw logs;
- private business strategy;
- consumer-specific product details that are not necessary for generic kit triage.

Use generic placeholders and concise summaries instead of copying raw evidence.

## Maintainer Handoff

After submission, maintainers triage the issue. Accepted reusable kit changes may become discovery
documents, implementation plans, release notes, or managed doc updates. Consumer-specific feedback
may be closed as local to the consumer.
