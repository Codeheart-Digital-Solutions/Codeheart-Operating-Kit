Last updated: 2026-06-23T18:09:34Z (UTC)

# Structure Governance

This managed domain owns reusable documentation placement, durable naming, ownership boundaries,
managed-content boundaries, and documentation index-maintenance guidance.

## Use

- Use `reference/documentation-structure.md` before creating, moving, archiving, or reorganizing
  docs.
- Use `reference/repository-information-architecture.md` before introducing durable names, new
  folder boundaries, or externally visible identifiers.
- Use `reference/managed-content-boundaries.md` to distinguish managed, scaffold, template,
  consumer-owned, local-user, generated, and report content.
- Use `reference/module-extension-state.md` before placing committed state for installed modules
  or extensions.
- Use `../agent-interface/reference/runbook-authoring-standard.md` when creating or materially
  changing durable runbooks.
- Use `runbooks/change-documentation-placement.md` before changing documentation placement.
- Use `runbooks/maintain-documentation-indexes.md` when discoverability changes.

## Boundary

Reusable generic structure doctrine belongs in the Operating Kit. Consumer repositories own local
product docs, local plans, local runbooks, local references, memory state, credentials,
environment details, and local exceptions.

When a consumer-local doc duplicates managed doctrine, convert it to a concise wrapper that points
to managed kit doctrine and keeps only real local exceptions.
