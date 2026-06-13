Last updated: 2026-06-13T22:55:57Z (UTC)

# Planning Document Lifecycle

Planning documents use lifecycle metadata so agents and maintainers can tell whether a document is
draft material, active execution authority, completed traceability, superseded guidance, or archive
history.

## Required Header

Every Markdown planning document starts with:

```text
Last updated: YYYY-MM-DDTHH:MM:SSZ (UTC)
Created: YYYY-MM-DD
Status: draft | active | completed | superseded | archived
```

Use the current UTC clock for `Last updated`. Preserve `Created` after the document is created.

## Optional Lifecycle Fields

Use these fields when useful:

```text
Completed: YYYY-MM-DD
Superseded by: <relative/path>
Execution log: <relative/path>
```

## Status Values

- `draft`: still being written and not approved for execution.
- `active`: current and approved for execution.
- `completed`: executed and retained for current context or traceability.
- `superseded`: replaced by another plan.
- `archived`: retained only as historical context.

## Revision Notes

Planning documents should maintain a bottom `# Revision Notes` section for meaningful decisions,
scope changes, strategy changes, and execution-plan changes. Do not add revision notes for
timestamp-only edits, typos, or checklist progress without scope change.

## Execution Logs

Goal-style implementation runs use a sibling `*_execution_log.md`. The log records meaningful
divergence, validation evidence, review-gate results, and follow-ups. It is not a command
transcript.
