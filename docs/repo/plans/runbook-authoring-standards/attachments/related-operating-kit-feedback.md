Last updated: 2026-06-24T06:50:35Z (UTC)

# Related Operating Kit Feedback

This attachment supports
`../runbook-authoring-standards_discovery_doc.md`.

It records public-core-safe feedback discovered while discussing human-facing module onboarding.
These notes are not implementation plans. They preserve context for later Operating Kit discovery,
feedback triage, or implementation planning.

## Feedback 1 - Slim Local Preference Contract For Human-Facing Runbooks

Summary: Human-facing and hybrid runbooks should reuse reliable interaction context, especially
language, before asking the user to choose it again.

Kit version: `0.1.9`

Affected area: `agent-interface`, local user layer, first-run onboarding, runbook authoring
standard, module onboarding.

Feedback type: doctrine gap / feature request.

Observed problem: The Operating Kit first-run onboarding asks for setup language and continues in
that language, but current repo state treats `.codeheart/user/preferences.yaml` mainly as an
ignored local-user concept with an example file. There is no strong contract that an actual
preferences file exists, no schema for it, and no generic rule requiring later module onboarding
runbooks to read it before asking language again.

Expected behavior: Human-facing and hybrid runbooks should resolve language context before
repeating setup questions:

1. Check whether `.codeheart/user/preferences.yaml` is visible to the agent and contains a readable
   `language` value.
2. If the preference exists, continue in that language.
3. If the preference is absent, unreadable, or does not contain language, ask the user once and
   continue in the answered language for the current flow.
4. If a future workspace-wide setting exists, use it only for preferences that are truly shared and
   do not conflict with local user preferences.

The preference path must be sufficiently visible from managed Operating Kit guidance, for example
from the future runbook authoring standard, the agent-interface README, or the managed root
`AGENTS.md` route. Module runbooks should not need to rediscover or restate the preference
resolution rule.

Recommended v1 preference scope should stay slim:

```yaml
language: English
timezone: Europe/Berlin
```

`language` is the immediate high-value preference because repeated language selection is visible
friction in every human-facing onboarding flow. `timezone` is useful for scheduling, timestamps,
and user-facing date explanations, but should remain optional. Broader preferences such as
explanation depth, tone, portal-instruction style, or approval style should be deferred until real
consumer usage proves they are worth standardizing.

Boundary: `.codeheart/user/preferences.yaml` should remain local-user state, ignored and not
managed. It must not contain secrets, credentials, tenant identifiers, account identifiers, or
private business data. Workspace-wide defaults, if ever needed, should be explicitly designed
separately from personal local preferences.

Proposed classification: `feedback-doctrine`, `feedback-feature`, `needs-discovery`.

Privacy confirmation: This note uses generic examples only and contains no customer, tenant,
account, secret, raw log, or local-machine details.

## Feedback 2 - Shared Environment Readiness And Tooling Register

Summary: Modules should not independently reinvent common tool checks and human-facing
installation guidance when the Operating Kit could own a reusable environment-readiness pattern.

Planning status: promoted to
`docs/repo/plans/tooling-environment-readiness/tooling-environment-readiness_discovery_doc.md`.

Kit version: `0.1.9`

Affected area: `agent-interface`, `native-codex-capabilities`, module onboarding, local readiness
state, tool installation runbooks.

Feedback type: doctrine gap / feature request.

Observed problem: Some modules need specific local tools. A Microsoft 365 module may require
PowerShell plus Microsoft Graph, PnP, and Exchange Online PowerShell modules. Future modules may
need Node.js, Python, browser automation, cloud CLIs, document tooling, or other local capabilities.
If every module owns the full human-facing tool-check and install flow independently, consumers can
see repeated checks, inconsistent explanations, duplicated approval gates, and stale local-state
assumptions.

Expected behavior: Operating Kit should eventually define a generic environment-readiness model:

- modules declare required and optional tools;
- agents check shared readiness before asking the user or repeating install guidance;
- missing tools route to generic human-facing install or repair runbooks for the tool family;
- module runbooks explain why the tool is needed for that module and what capability it unlocks;
- install steps remain approval-gated and user-visible;
- readiness results are recorded locally with timestamp, source, status, and relevant caveats;
- modules remain responsible for domain-specific permissions, tenant/cloud access, and action
  recipes.

Illustrative future state shape:

```yaml
tools:
  powershell:
    status: available
    checked_at: "YYYY-MM-DDTHH:MM:SSZ"
  microsoft_graph_powershell:
    status: missing
    checked_at: "YYYY-MM-DDTHH:MM:SSZ"
    install_runbook: "<managed or module route>"
```

Non-goals for the generic Operating Kit layer:

- do not become a broad wrapper around every external tool;
- do not bypass module-specific setup, consent, credentials, or permissions;
- do not install tools without an explicit approval gate;
- do not treat a stale local readiness record as proof that external access is currently valid.

Proposed classification: `feedback-feature`, `feedback-doctrine`, `needs-discovery`.

Privacy confirmation: This note describes tool families and generic readiness state only. It
contains no customer, tenant, account, secret, raw log, or local-machine details.
