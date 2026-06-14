Last updated: 2026-06-13T22:55:57Z (UTC)

# Local Extension Contract

Consumer owners add local rules below the managed Operating Kit block.

## Repository-Owned Instructions

Use repository-owned sections for product rules, local commands, operational guardrails,
validation policy, local doc routing, cloud account rules, and domain-specific safety constraints.

## Local User Guidance

Use `.codeheart/user/` for ignored local preferences and personal notes. Do not commit secrets,
tokens, credentials, or private machine state.

## Conflict Handling

When local guidance conflicts with managed kit rules, preserve immediate safety rules and flag the
conflict. Reusable generic guidance should be proposed as an Operating Kit change, not copied into
consumer `docs/repo/`.
