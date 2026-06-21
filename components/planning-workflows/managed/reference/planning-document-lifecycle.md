Last updated: 2026-06-21T15:07:29Z (UTC)

# Planning Document Lifecycle

Planning documents use lifecycle metadata and predictable placement so agents and maintainers can
tell whether a document is draft material, active execution authority, completed traceability,
superseded guidance, or archived history.

## Planning Document Types

- `*_discovery_doc.md`: decision and requirement discovery.
- `*_implementation_doc.md`: execution-ready implementation plan.
- `*_execution_log.md`: goal-style implementation evidence and divergence log.
- Plan-scoped attachments: temporary inventories, reviews, reports, and evidence used by one plan.

## Required Planning Header

Discovery and implementation documents start with:

```text
Last updated: YYYY-MM-DDTHH:MM:SSZ (UTC)
Created: YYYY-MM-DD
Status: draft | active | completed | superseded | archived
```

Use the current UTC clock for `Last updated`. Preserve `Created` after the document is created.

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

## Plan Register Relationship

`docs/repo/plans/plan-register.md` is the lightweight index for formal plans, plan families, major
workstreams, and portfolio-relevant planning records.

Planning documents remain the source of truth for their own lifecycle state, scope, decisions,
execution evidence, and handoff details. The register may copy lifecycle metadata as a snapshot.
When a register entry and canonical planning document disagree, refresh the register from the
canonical document.

Use `../runbooks/maintain-plan-register.md` when a discovery document or implementation plan is
created, materially changed, activated, completed, superseded, archived, or linked to parent,
child, dependency, supersession, or related planning records.

Do not treat review-only workflows as side-effecting register maintenance. A planning-document
review may report a stale or missing register entry as a finding. The register update belongs to a
later user-authorized discovery, implementation-planning, or execution edit that materially
changes or refreshes the canonical plan.

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
