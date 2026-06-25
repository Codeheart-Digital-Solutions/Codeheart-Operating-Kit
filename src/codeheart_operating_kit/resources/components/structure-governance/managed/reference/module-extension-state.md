Last updated: 2026-06-25T13:05:46Z (UTC)

# Module Extension State

Use this reference when a repository has installed modules or extensions that need committed,
non-secret routing context.

## Standard Path

Committed module or extension state belongs under:

```text
docs/repo/state/<module-or-extension-id>/
```

Use the stable module or extension ID exposed by the installed module system as the folder name.
The Operating Kit defines this placement convention, but it does not define every module registry,
module manifest, state file, or schema.

## Ownership

- Operating Kit owns the generic placement and routing doctrine.
- Module systems own how installed modules and extension IDs are discovered.
- Modules own concrete state file names, schemas, lifecycle, route registries, route cards, and
  validation rules.
- Consumer repositories own committed state contents after creation.
- Live external systems own current truth for external resources.

Use `../../agent-interface/reference/operation-routing-and-dispatch.md` for the generic routing
sequence, capability advertisement fields, and route-card field semantics. This reference only
defines where committed module or extension routing state belongs and how it relates to live
preflight.

## Allowed Content

Use `docs/repo/state/<module-or-extension-id>/` for committed, repository-owned context that helps
agents route work before asking repeated setup questions or touching external systems.

Appropriate examples include:

- selected non-secret workspace identifiers;
- non-secret target names, slugs, or URLs;
- routing registries for repo-owned module setup;
- onboarding facts that are safe to share with repository contributors;
- stable lifecycle notes that tell an agent what live preflight to run next.

Example shape:

```text
docs/repo/state/
  m365-workspace/
    workspace-registry.yaml
    tenant-onboarding.yaml
  finance-ops/
    account-registry.yaml
  crm-sync/
    sync-targets.yaml
  document-automation/
    publishing-registry.yaml
```

These examples are placeholder module IDs and file names. A real module owns its own schema and
validation rules.

## Prohibited Content

Do not store any of the following in committed module or extension state:

- secrets, credentials, tokens, certificate material, or token caches;
- sensitive personal data;
- private tenant details, customer records, mailbox contents, or account identifiers;
- raw operational logs or message traces;
- generated lockfiles, checksums, or transient run records;
- local machine paths, local user preferences, or ignored personal notes;
- live external truth that must be read from the external system before action.

Use `.codeheart/user/` for ignored local-only user notes. Use module-owned run records or report
locations for execution evidence. Use live service preflight for external truth.

## Live Preflight Boundary

Committed module state routes the agent; it does not authorize action.

Before sensitive reads, writes, deletes, permission changes, tenant changes, or external resource
changes, the agent must validate the relevant live system through the module's approved tool lane
or runbook. If local state conflicts with live preflight, stop and resolve the conflict through
the module's runbook instead of trusting the committed file.

Route cards may point to committed module state as a state source. They must still name the live
truth source, preconditions, approval class, stop conditions, and canonical runbook before any
sensitive or external action proceeds.

## Agent Routing Rule

When operating installed modules or extensions:

1. Discover the installed module or extension ID through the module system present in the
   repository.
2. Check `docs/repo/state/<module-or-extension-id>/` for committed non-secret routing context.
3. Use the module's runbooks and schemas for concrete file interpretation.
4. Run live preflight before sensitive reads or external changes.
5. Ask the user only for missing information that is not already visible in committed safe state.

## Anti-Examples

- Adding a root `AGENTS.md` route for each installed module.
- Treating `docs/repo/state/<id>/` as a token cache or generated run-output folder.
- Making Operating Kit define a Foundry-specific lockfile or module manifest shape.
- Creating empty state folders in every consumer repository before a module needs committed state.
- Proceeding with an external change only because a committed state file names a target.
