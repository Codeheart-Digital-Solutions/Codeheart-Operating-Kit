Last updated: 2026-06-17T19:42:42Z (UTC)

# Triage Kit Feedback

Use this runbook when maintaining public GitHub Issues for Operating Kit feedback intake.

## Scope

This runbook applies to public, sanitized feedback submitted through the Operating Kit issue forms.
It does not define a private security disclosure process.

## Triage Procedure

1. Confirm the issue is public-core-safe before quoting, labeling, or summarizing it.
2. Apply or verify the feedback type label from
   `docs/repo/reference/kit-feedback-label-taxonomy.md`.
3. Remove `needs-triage` after the first maintainer decision.
4. Choose one lifecycle route:
   - `needs-shaping`;
   - `needs-information`;
   - `consumer-specific`;
   - `duplicate`;
   - `declined`;
   - `accepted-backlog`;
   - `needs-discovery`;
   - `implementation-planned`;
   - `released`;
   - `superseded`.
5. Leave a concise triage comment that names the route and next action.
6. Link accepted work to a discovery document, implementation plan, pull request, or release note
   when one exists.

## Lifecycle States

| State | Use When | Maintainer Action |
| --- | --- | --- |
| `needs-triage` | New issue has not been reviewed. | Keep until first triage decision. |
| `needs-shaping` | Rough feedback is not implementation-ready. | Ask clarifying questions or convert to a clearer issue. |
| `needs-information` | Reporter can provide more public-safe detail. | Ask for specific sanitized information. |
| `consumer-specific` | The issue belongs in a consumer repository or workspace. | Close with local-ownership guidance. |
| `duplicate` | Existing issue or plan already covers the feedback. | Link the canonical issue or plan, then close. |
| `declined` | Feedback is out of scope or not actionable for the kit. | Explain the reusable-kit boundary and close. |
| `accepted-backlog` | Feedback is generic kit work but not yet planned. | Keep open or link to backlog tracking. |
| `needs-discovery` | Decisions are needed before implementation planning. | Create or update a discovery document. |
| `implementation-planned` | Accepted work has an implementation plan. | Link the plan and keep state current. |
| `released` | The change shipped in release notes or a public release. | Link the release evidence and close when appropriate. |
| `superseded` | Newer issue, plan, or release path replaces this item. | Link the replacement and close when appropriate. |

## Discovery And Implementation Handoff

Create or update a discovery document when feedback changes reusable doctrine, safety policy,
ownership boundaries, managed routes, or implementation direction and the right answer is not
obvious.

Create or update an implementation plan when the feedback is accepted, the intended solution is
clear, and execution touches managed docs, templates, validators, schemas, CLI behavior, release
assets, or repository governance.

Use `docs/repo/runbooks/change-operating-kit.md` before changing the kit. Use
`docs/repo/reference/consumer-impact-classification.md` to record consumer impact before accepted
work ships.

## Consumer-Specific Closure

Close as `consumer-specific` when the feedback is valid but belongs to a consumer repository,
workspace, local user layer, product docs, local command inventory, or private operating context.

Do not promote consumer-specific details into the public kit. When reusable guidance might exist,
ask the reporter for a sanitized generic summary or route maintainers through
`docs/repo/runbooks/promote-consumer-change.md`.

## Accidental Public Disclosure Response

If an issue or comment includes sensitive material:

1. Do not copy or quote the sensitive material.
2. Use available GitHub moderation tools to hide, edit, delete, or minimize public exposure.
3. Ask the reporter to rotate exposed secrets or credentials when relevant.
4. Preserve only a sanitized summary for triage.
5. Escalate privately to the kit owner when the exposure cannot be handled through normal issue
   moderation.

After exposure handling, continue triage only with public-safe summaries.

## Release And Sync Notes

Accepted feedback that changes managed content, templates, validators, safety policy, generated
surfaces, or repository governance needs release-note review before shipping. Record whether
consumers need no action, a normal sync, a repair, or a migration.
