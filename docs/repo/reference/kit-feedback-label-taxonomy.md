Last updated: 2026-06-17T19:11:03Z (UTC)

# Kit Feedback Label Taxonomy

This reference defines the v1 GitHub label set for public Operating Kit feedback intake.

## Feedback Type Labels

| Label | Meaning | Default form |
| --- | --- | --- |
| `feedback-rough` | Early friction, weak signal, or idea that needs shaping before implementation. | Rough feedback or idea |
| `feedback-bug` | Public-safe defect or regression in Operating Kit behavior or docs. | Kit bug or regression |
| `feedback-doctrine` | Missing, confusing, stale, or conflicting reusable doctrine or workflow guidance. | Doctrine or workflow gap |
| `feedback-install-sync-check` | Install, sync, check, update-check, inspect, or release-asset friction. | Install, sync, or check issue |
| `feedback-docs-routing` | Documentation route, index, wording, or placement confusion. | Docs clarity or routing issue |
| `feedback-feature` | Reusable Operating Kit feature or capability request. | Feature or capability request |

## Lifecycle Labels

| Label | Meaning | Applied by |
| --- | --- | --- |
| `needs-triage` | New feedback awaiting maintainer review. | Issue forms |
| `needs-shaping` | Rough feedback that is not yet implementation-ready. | Rough feedback form or maintainer |
| `needs-information` | Maintainer needs more public-safe information from the reporter. | Maintainer |
| `accepted-backlog` | Maintainer accepted the issue as generic kit backlog. | Maintainer |
| `needs-discovery` | Accepted feedback requires discovery before implementation planning. | Maintainer |
| `implementation-planned` | Accepted feedback has or is linked to an implementation plan. | Maintainer |
| `released` | Feedback has shipped in a release. | Maintainer |
| `consumer-specific` | Feedback is valid for a consumer context but not generic kit doctrine. | Maintainer |
| `superseded` | Feedback was replaced by a newer issue, plan, or release path. | Maintainer |
| `declined` | Maintainer declined the feedback as out of scope or not actionable for the kit. | Maintainer |
| `duplicate` | Feedback duplicates an existing issue or plan. | Maintainer |

## Transition Rules

- New issue forms apply one `feedback-*` type label and `needs-triage`.
- Rough feedback also starts with `needs-shaping`.
- Maintainers remove `needs-triage` after the first triage decision.
- Maintainers use `needs-information` only when the reporter can provide more public-safe detail.
- Maintainers use `accepted-backlog` when the feedback is reusable kit work but not yet planned.
- Maintainers use `needs-discovery` when a decision ledger is needed before implementation.
- Maintainers use `implementation-planned` when an accepted issue links to an active or draft
  implementation plan.
- Maintainers use `released` only after the feedback is included in release notes or a published
  release.
- Maintainers use `consumer-specific`, `declined`, `duplicate`, or `superseded` to close feedback
  that should not become new generic kit work.

## Public-Core Rule

Labels do not relax public-core hygiene. Every issue body and maintainer comment must exclude
secrets, credentials, customer or tenant details, local machine state, account identifiers, raw
logs, and private strategy.
