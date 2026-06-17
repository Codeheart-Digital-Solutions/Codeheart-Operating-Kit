Last updated: 2026-06-17T19:21:22Z (UTC)

# Kit Feedback Item Format

Use this reference to structure sanitized Operating Kit feedback before opening a public issue.

## Fields

| Field | Purpose |
| --- | --- |
| Summary | One sentence describing the reusable kit concern or idea. |
| Kit version | Installed version, release tag, or `unknown`. |
| Affected area | Component, runbook, document, command, or workflow. |
| Feedback type | Rough feedback, bug, doctrine gap, install/sync/check issue, docs routing issue, or feature request. |
| Observed problem | Public-safe description of what happened or what is missing. |
| Expected behavior | Public-safe description of the desired reusable kit behavior. |
| Sanitized evidence | Generic reproduction notes, public-safe links, or summarized output. |
| Privacy confirmation | Confirmation that sensitive material was removed. |
| Proposed classification | Suggested maintainer route when known. |

## Public-Safe Example

```text
Summary: The managed documentation route for a recurring planning case is unclear.
Kit version: 0.1.4
Affected area: planning-workflows
Feedback type: doctrine or workflow gap
Observed problem: Agents may create local guidance instead of using a managed route.
Expected behavior: The managed route should explain when to group related plans.
Sanitized evidence: The issue appeared during a generic repository cleanup discussion.
Privacy confirmation: No secrets, credentials, customer or tenant details, local machine state,
account identifiers, raw logs, or private strategy are included.
Proposed classification: needs-triage, feedback-doctrine
```

## Privacy Rule

Do not use this format to store raw private evidence. Keep private evidence outside the feedback
workflow and submit only a sanitized summary or pointer.
