Last updated: 2026-07-31T22:23:19Z (UTC)

# Planning Document Lifecycle

Planning documents use lifecycle metadata and predictable placement so agents and maintainers can
tell whether a document is draft material, active execution authority, completed traceability,
superseded guidance, or archived history.

## Planning Document Types

- `*_discovery_doc.md`: decision and requirement discovery.
- `*_implementation_doc.md`: execution-ready implementation plan.
- qualifying nested family `README.md`: shared discoverability for independent sibling plans.
- `*_execution_log.md`: goal-style implementation evidence and divergence log.
- Plan-scoped attachments: temporary inventories, reviews, reports, and evidence used by one plan.

Discovery, implementation, and family are separate canonical record kinds. A discovery and its
implementation may share a family and relations, but neither is merely a phase row inside the
other. Execution logs and attachments are supporting evidence, not independent plan records.

## Required Planning Header

Discovery and implementation documents start with:

```text
Last updated: YYYY-MM-DDTHH:MM:SSZ (UTC)
Created: YYYY-MM-DD
Status: draft | active | completed | superseded | archived
```

Use the current UTC clock for `Last updated`. Preserve `Created` after the document is created.

## Canonical Metadata

In mixed and canonical mode, place the bounded Codeheart plan metadata block immediately below the
canonical Markdown title. It owns stable semantic identity, record kind, public-safe purpose,
catalog chronology, optional family/classifications/relations, and legacy aliases.

Use `plan-catalog-format.md` for the exact marker and fields. The lifecycle header continues to own
`Status`; do not add a separate authority field. Metadata-only adoption updates
`catalog_metadata_updated` and must not change the historical `Last updated` value.

## Optional Lifecycle Fields

Use these fields when they help review:

```text
Completed: YYYY-MM-DD
Superseded by: <relative/path>
Execution log: <relative/path>
```

Add `Execution log:` when a sibling execution log exists.

## Status Values

- `draft`: still being written and not approved for execution.
- `active`: current and approved for execution.
- `completed`: executed and retained for current context or traceability.
- `superseded`: replaced by another plan.
- `archived`: retained only as historical context.

Use one status value. Do not execute a `draft`, `completed`, `superseded`, or `archived`
implementation plan without explicit user approval.

The lifecycle value is not a branch, push, pull-request, merge, or release state. An `active` plan
is approved execution authority only within current user instruction, repository policy, and its
own scope.

## Catalog Modes And Authoring

- `legacy`: existing plans may lack metadata and the legacy register may still be maintained.
- `mixed`: every new or materially rewritten formal plan receives metadata; only exact
  pre-cutover register-linked plans may remain legacy temporarily.
- `canonical`: every formal discovery, implementation, and qualifying family record requires
  metadata.

Read mode from `.codeheart/kit.config.yaml`. Do not assume that an absent setting means canonical;
absence defaults to legacy. Use the migration runbook for mode changes.

## Revision Notes

Discovery and implementation documents maintain a bottom `# Revision Notes` section for
meaningful decision, scope, strategy, or execution-plan changes.

Do not add revision notes for:

- timestamp-only edits;
- typos;
- formatting-only edits;
- checklist progress that does not change scope.

## Checkbox State

Implementation documents may use checkbox tasks as execution state. Tick a checkbox only after the
work is completed and validated. Do not tick future work, deferred work, or work that has not been
validated.

Checkbox state is a human-readable execution aid. It does not replace issue state, pull-request
state, release state, or execution logs.

## Execution Logs

Goal-style implementation runs use a sibling execution log:

```text
<feature-slug>_implementation_doc.md
<feature-slug>_execution_log.md
```

For a plan bundle:

```text
<feature-slug>/
  <feature-slug>_implementation_doc.md
  <feature-slug>_execution_log.md
```

The execution log records:

- meaningful divergence from the plan;
- added tasks;
- safe defaults;
- validation substitutions or extensions;
- review-gate rounds and accepted result;
- release, migration, or handoff evidence;
- residual risk and follow-ups.

The execution log is not a command transcript.

## Plan Shapes

Use the smallest shape that preserves ownership, reviewability, and lifecycle clarity.

### Standalone Plan File

Use a standalone plan file for small, low-artifact work unlikely to need an execution log,
attachments, or repeated goal-style execution.

```text
plans/
  <feature-slug>_implementation_doc.md
```

### Plan Bundle

Use a plan bundle when the plan has or may need:

- discovery plus implementation documents;
- execution log;
- attachments;
- repeated goal-style execution;
- multiple epics;
- cross-area review.

```text
plans/
  <feature-slug>/
    <feature-slug>_discovery_doc.md
    <feature-slug>_implementation_doc.md
    <feature-slug>_execution_log.md
    attachments/
      <plan-scoped-artifact>.md
```

### Subplan

Use `subplans/` when a child plan exists mainly to edit, audit, split, repair, or execute a parent
plan, or when the parent defines the child plan's correctness.

```text
<parent-plan-slug>/
  <parent-plan-slug>_discovery_doc.md
  subplans/
    <subplan-slug>/
      <subplan-slug>_implementation_doc.md
      <subplan-slug>_execution_log.md
      attachments/
```

Do not use `subplans/` merely because plans share a topic. Link related plans instead.

### Plan Family

Use a plan family when independently executable sibling plan bundles need shared discoverability
without one plan owning the others.

```text
plans/
  <family-slug>/
    README.md
    <first-plan-slug>/
      <first-plan-slug>_implementation_doc.md
    <second-plan-slug>/
      <second-plan-slug>_implementation_doc.md
```

The second related sibling plan bundle is the normal trigger for a shared family folder.

### Program Folder

Use a program folder only when a parent coordination plan owns multiple workstream plans and has
current authority over scope, sequencing, review state, or cross-workstream decisions.

```text
<program-slug>/
  <program-slug>_implementation_doc.md
  <program-slug>_execution_log.md
  plans/
    <workstream-plan-slug>/
      <workstream-plan-slug>_implementation_doc.md
```

## Plan-Scoped Attachments

Put plan-scoped artifacts in the plan bundle's `attachments/` folder when the material is too
large, too detailed, too volatile, or too audit-specific to keep in the plan.

Examples:

- path migration maps;
- migration inventories;
- stale-path reports;
- one-off review notes;
- temporary decision matrices;
- validation inventories;
- migration ledgers.

Promote an attachment to a durable `reference/` document only when it becomes reusable doctrine.

## Archive Behavior

Use `archive/` for superseded or historical plans retained for traceability. Do not move a plan to
archive during execution unless the user explicitly asks or the active plan includes an approved
archive task.

Archived plans remain historical by default. New active work should link to archived history rather
than adding active plans under `archive/`.

## Plan Views, Register, And Git Visibility

`docs/repo/plans/plan-register.md` remains the stable entry point. In mixed and canonical mode the
current repository view is generated by `codeheart-operating-kit plans list`; the register is not a
manually appended second catalog. A pre-existing register remains byte-frozen legacy evidence from
the mixed cutover revision.

Canonical planning documents remain source of truth for lifecycle, scope, decisions, metadata,
execution evidence, and handoff details. Git provides path, ref, revision, and content identity.
Use `plan-register-format.md` and `../runbooks/maintain-plan-register.md` for mode-specific behavior.

Portfolio coordination sees only pushed remote refs. A plan created or activated on a work branch
becomes globally observable after its authorized normal push; worktree edits, local heads, and
unpushed commits remain local. The coordination home's scanner reports freshness and completeness
and must not present a failed refresh or an older cache as current.

Review-only workflows report missing/invalid metadata or stale views as findings. They do not
silently change lifecycle, register bytes, metadata, config mode, or remote visibility.

## Index Maintenance

Update the nearest README and parent index when discoverability changes:

- new plan, runbook, reference, or attachment that readers need to find;
- moved, renamed, archived, or removed document;
- changed title, purpose, command path, runbook path, or validator path;
- new execution log;
- changed planning entry point.

Index updates are not required for timestamp-only edits or checklist progress inside an already
linked plan.

## Relationship To Structure Governance

Planning workflows own planning lifecycle semantics:

- document types;
- required planning headers;
- status values;
- revision notes;
- execution logs;
- plan bundles, subplans, plan families, and program folders;
- plan-scoped attachments.

Structure governance owns durable placement, naming, index-maintenance rules, and local wrapper
behavior for consumer repositories.

When a consumer keeps a local planning or governance document only for discoverability, that local
document should be a concise wrapper to the managed kit doctrine plus any real local exceptions.
Do not duplicate the managed planning lifecycle rules in a consumer-local wrapper.
