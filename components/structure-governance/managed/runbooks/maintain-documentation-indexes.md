Last updated: 2026-06-17T06:46:46Z (UTC)

# Maintain Documentation Indexes

Use this runbook when documentation discoverability changes.

README indexes are routers. They should help readers find the owner of a rule, plan, runbook,
reference, product, archive, or generated report without duplicating durable doctrine.

## Trigger

Update indexes when:

- a discoverable doc, runbook, plan, reference, execution log, or attachment is added;
- a doc path, title, or purpose changes;
- a doc is moved, archived, deleted, or removed from current authority;
- a product, source area, package, module, or docs root is created;
- a runbook path, command path, validator input, workflow command, or entry point changes.

Index updates are not required for timestamp-only edits, wording inside an already-linked doc, or
checklist progress inside an existing plan.

## Procedure

1. Identify the nearest README router for the changed path.
2. Identify parent indexes that route readers to that area.
3. Add concise entries for new durable docs and execution logs.
4. Update entries for changed paths, titles, or purposes.
5. Remove current-routing entries for archived or deleted docs.
6. Preserve historical references in completed plans and execution logs when they are traceability
   records.
7. Keep routers concise; point to the owner instead of copying full rules.
8. When a local wrapper replaces duplicated doctrine, route to the wrapper only when local
   exceptions matter. Otherwise route to the managed kit owner.

## Archive And Removal Semantics

- Archived docs may remain linked from an archive README or historical plan.
- Current routers must not present archived docs as current authority.
- Deleted docs should be removed from current routers.
- Superseded docs should link to their replacement when the path remains discoverable.

## Validation

Run:

```sh
git diff --check
```

Then manually verify:

- nearest README includes the changed entry when discoverability changed;
- parent index routes to the nearest README or directly to the doc when appropriate;
- no current router points to an archive as current authority;
- local wrappers point to managed doctrine plus local exceptions;
- moved command, validator, workflow, and runbook paths are updated;
- historical references remain historical rather than current routing.
