Last updated: 2026-06-13T22:55:57Z (UTC)

# Discovery Workflow

Use this runbook for substantial discovery, architecture decision framing, requirements
clarification, and discovery document drafting.

## Procedure

1. Restate the problem, intended outcome, and current constraints.
2. Run targeted repo or workspace research before asking broad follow-up questions.
3. List essential context files and why each matters.
4. Identify decision axes as `D-*` items with owner, options, criteria, and status.
5. Track functional requirements as `FR-*` and non-functional requirements as `NFR-*`.
6. Track open questions as `OQ-*` with `BLOCKER: yes` or `BLOCKER: no`.
7. Track assumptions as `A-*` and risks as `R-*`.
8. Ask only high-impact questions, in batches of no more than five.
9. Converge decisions before drafting implementation work.
10. Draft or update a discovery document when blocker decisions are resolved or when the user asks
    for a draft with blockers carried forward.

## Goal-Style Discovery

For goal-style discovery, use the discovery document as the decision ledger. Refresh the decision
inventory, work through blocker and implementation-shaping decisions one at a time, add new
decisions when found, and stop when the document is ready for manual review.

## Handoff

End discovery with an implementation handoff that names approved decisions, MVP requirements,
critical constraints, remaining blockers, high-level epic outline, and minimum validation gates.
Do not include a granular implementation checklist in discovery output.
