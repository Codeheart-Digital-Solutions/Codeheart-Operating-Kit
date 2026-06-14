Last updated: 2026-06-13T22:55:57Z (UTC)

# Execute Implementation Plan

Use this runbook when executing an active `*_implementation_doc.md`.

## Procedure

1. Confirm the plan header is `Status: active`.
2. Create or update the sibling execution log.
3. Execute epics sequentially.
4. Treat each epic outcome as the authority, not only its checkbox list.
5. Add low-risk missing tasks when required to achieve the epic outcome.
6. Run the smallest validation set that proves the changed surface.
7. Run the per-epic review gate before marking the epic complete.
8. Update the plan checklist only after work is complete and validated.
9. Record meaningful divergence, validation substitutions, review findings, fixes, and evidence in
   the execution log.
10. Keep the plan active until all epics pass validation.

## Review Gate

Before marking an epic complete, use a fresh read-only reviewer agent when available. The reviewer
checks the implementation against the epic outcome, acceptance criteria, completed checklist,
validation evidence, execution-log state, scope boundaries, and accidental future-epic work.

Fix material findings and repeat with a fresh reviewer until no material issues remain or a blocker
is recorded.
