Last updated: 2026-06-26T15:57:38Z (UTC)

# Local Extension Contract

Consumer owners add local rules below the managed Operating Kit block.

## Repository-Owned Instructions

Use repository-owned sections for product rules, local commands, operational guardrails,
validation policy, local doc routing, cloud account rules, and domain-specific safety constraints.

## Local User Guidance

Use `.codeheart/user/` for ignored local preferences and personal notes. Do not commit secrets,
tokens, credentials, or private machine state.

## Local Machine State

Use `.codeheart/local/` for ignored generated runtime and tooling state that is specific to one
checkout and can be recreated. Examples include repo-local virtual environments, local caches,
temporary files, generated shims, package artifacts, generated install metadata, and
editable-install artifacts.

Do not put human preferences, personal notes, secrets, credentials, durable repo state, managed
snapshots, or live external truth under `.codeheart/local/`.

## Conflict Handling

When local guidance conflicts with managed kit rules, preserve immediate safety rules and flag the
conflict. Reusable generic guidance should be proposed as an Operating Kit change, not copied into
consumer `docs/repo/`.
