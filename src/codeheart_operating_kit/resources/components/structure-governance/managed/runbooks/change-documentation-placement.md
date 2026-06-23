Last updated: 2026-06-23T18:09:34Z (UTC)

# Change Documentation Placement

Use this runbook before adding, moving, renaming, archiving, deleting, or reclassifying
documentation.

## Preflight

1. Identify the durable owner.
2. Identify artifact kind: router, plan, runbook, reference, archive, scaffold, template, managed
   doc, committed module state, generated report, or local wrapper.
3. Read the documentation structure reference.
4. Read the repository information architecture reference when a durable name or new folder
   boundary is involved.
5. Check for an existing owner before adding a new document.
6. Protect unrelated user work.

## Owner Selection

Use this order:

1. Managed Operating Kit docs for reusable generic operating doctrine.
2. Consumer repo docs for local repository rules, plans, commands, and exceptions.
3. Product docs for product-owned rules, plans, runbooks, and references.
4. Package, module, or source-area docs for local source-owned guidance.
5. Consumer repo state under `docs/repo/state/<module-or-extension-id>/` for committed
   non-secret routing context owned by installed modules or extensions.
6. Agent memory for curated continuity state.
7. Business or research docs only when the repository intentionally stores that type of material.

If proposed content is generic and reusable, recommend changing the Operating Kit instead of
adding duplicate consumer-local doctrine.

## Procedure

1. State the current path and proposed target path.
2. State owner, lifecycle, artifact kind, and discoverability reason.
3. Check whether the change creates, removes, or changes a README route.
4. Check whether the change affects commands, validators, workflow paths, generated paths, or
   public links.
5. Move or create the document in the owning path.
6. Convert duplicate local doctrine to a wrapper when managed kit doctrine owns the generic rule.
7. Place committed non-secret module or extension routing state under
   `docs/repo/state/<module-or-extension-id>/`.
8. Place plan-scoped temporary artifacts under the plan bundle's `attachments/` folder.
9. Place historical material under `archive/` only when it no longer governs current behavior.
10. Update nearest README and parent indexes.
11. Update links, commands, validator inputs, and stale-path guards affected by the change.
12. Record consumer impact when generated surfaces or managed routes change.

## Archive Handling

Archive only when material is superseded or retained for traceability. Do not archive a current
source of truth just because the path is old.

Historical plans may retain old paths as traceability records. Current routers should point to
current owners.

## Validation

Run checks appropriate to the consumer repository, at minimum:

```sh
git diff --check
```

Also validate:

- changed Markdown files start with the timestamp header;
- nearest README and parent index routes are current;
- moved paths have updated links;
- archived docs are not still presented as current authority;
- generated or managed paths were not changed without impact classification;
- committed module or extension state is non-secret and remains under `docs/repo/state/<id>/`;
- local wrappers point to the managed owner and keep only real local exceptions.
