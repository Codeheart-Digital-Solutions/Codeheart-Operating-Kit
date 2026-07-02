Last updated: 2026-07-02T13:16:41Z (UTC)

# Repo Feedback Item Format

Use this reference to structure sanitized feedback about the repository that owns the affected
runbook, script, module, product area, or document.

Repo feedback is not Operating Kit feedback by default. Use
`runbooks/submit-kit-feedback.md` only when the issue is reusable Operating Kit doctrine,
managed content, install, sync, or routing behavior.

## Title

Recommended title shape:

```text
[feedback][<kind>][<area>] <short issue title>
```

Use a short plain-language title. Do not include account identifiers, tenant names, local paths,
raw error payloads, customer details, or private business context.

## Body Fields

| Field | Purpose |
| --- | --- |
| Summary | One sentence describing the repository-owned friction. |
| Owning repository | Repository name or public-safe placeholder when the issue is drafted elsewhere. |
| Affected area | Runbook, script, module, product area, docs area, or workflow. |
| Feedback kind | Runbook gap, tooling gap, docs conflict, bug, feature, or repeated friction. |
| Observed problem | Sanitized description of what happened or what was unclear. |
| Expected behavior | Sanitized description of the desired repository-owned behavior. |
| Sanitized evidence | Public-safe summaries, command names, or short generic observations. |
| Workaround used | Sanitized workaround or `none`. |
| Proposed classification | Suggested labels or triage route when known. |
| Privacy confirmation | Confirmation that forbidden content was removed. |

## Recommended Labels

Use existing labels only. Do not block issue creation when labels or templates are missing.

Recommended standard labels when already present:

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

Put classification in the issue body when labels are unavailable.

## Forbidden Content

Do not include:

- secrets, credentials, API keys, passwords, tokens, OAuth codes, or MFA material;
- customer, tenant, mailbox, account, or personal identifiers;
- raw logs, command transcripts, provider responses, environment dumps, or local machine state;
- raw mailbox, document, invoice, image, PDF, source data, or operational content;
- local absolute paths or private repository paths;
- private strategy, commercial details, or sensitive operational context;
- security vulnerability details that need a private disclosure route.

Use short sanitized summaries and placeholders instead.

## Triage Promotion

Repo feedback issues are intake. During triage, maintainers choose one route:

- direct patch for clear small fixes;
- batch with related feedback;
- discovery when doctrine, ownership, or behavior is unclear;
- implementation plan when accepted work changes managed content, schemas, validators, release
  surfaces, or cross-repo behavior;
- defer, duplicate, supersede, decline, or close;
- route to Operating Kit feedback when the issue is generic kit work.
