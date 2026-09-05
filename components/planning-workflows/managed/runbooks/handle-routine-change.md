Last updated: 2026-09-05T19:15:20Z (UTC)

# Handle Routine Change

Audience: agent-facing

Intent:
Complete a bounded, understood change through its owning route with proportionate evidence and
reuse of sufficient existing authority.

Success:
The requested outcome works, relevant checks pass, and the owning record states the actual result
and delivery state.

Agent judgment boundary:
Choose ordinary implementation details within the authorized outcome. Do not use this route to
escape a plan's epic, review, module contract, managed boundary, or controlled-operation gate.

Stop boundary:
Return to planning or the acceptance owner when the outcome, impact, method, ownership, scope, or
authority is materially uncertain or changes. Preserve work when a tool-enforced gate rejects an
action; do not bypass it.

## Route And Inputs

Route: `planning.routine-change` in
`../../agent-interface/reference/operation-routing-and-dispatch.md`.

Read the request, root instructions and nearest owner guidance. Inputs are the intended outcome,
known target and method, current state, applicable authority, and an existing PR, issue, or owning
record. Inspect that record before creating another. Use the owner's existing editor, command or
runbook after routing; this runbook adds no CLI command or approval registry.

A change is routine when all of these hold:

- the owner and method are known;
- impact is bounded and understood;
- the change can be reversed within the understood boundary;
- sufficient current authority covers the intended effects.

Line count, file count, and speed are not classifiers. A one-line permission change can require
formal planning. A known correction across several files can remain routine. Creating a durable
path still requires the placement check. Discover installed module IDs through the repository's
module system and read `docs/repo/state/<id>/` plus the module's own route before operating it.

## Execution

1. State the routine outcome and owner. If the action is necessary to an active plan's epic,
   continue in that named epic using `execute-implementation-plan.md`; do not relabel it routine.
   For a genuinely separate repair, record its relationship to affected plan work.
2. Inspect current state and protect overlapping user changes. Confirm eligibility and the owning
   method. If uncertain, use `discovery-workflow.md` or `draft-implementation-plan.md` at the scope
   needed to settle the uncertainty.
3. Record purpose, scope and applicable authority in the existing PR, issue or owning record.
   A short paragraph is enough. Create a short owner-placed record only when none fits; no mini-plan
   or new schema is required.
4. Check the actual authorization against target, effects and limits. A request for the covered
   change can supply it. Reuse sufficient earlier approval; do not ask again for ordinary covered
   steps. A routine label grants no push, PR, release, install, sensitive-read or provider effect
   that the request and owner contract do not cover. Resolve required exact inputs before action.
5. Make the change through the owner route. Use
   `../../agent-interface/runbooks/handle-tooling-readiness.md` for missing local tools before
   improvising setup. Keep service preflight with the module or service owner.
6. Run the smallest meaningful checks that prove the outcome, including mandatory owner gates.
   Inspect a document correction; exercise changed behavior for a defect; verify materialization
   and preservation for managed content. Add a regression when it proves a breakable contract,
   rather than mechanically testing prose or every reversible edit.
7. Review the result and complete any covered Git checkpoint. Record result, validation, remaining
   limits, commit/push/PR state and the next owner. If validation fails, correct within scope or
   preserve the current state with the concrete blocker and recovery route.

## Evidence And Handoff

Example existing-record update:

```text
Purpose/scope: correct the known link in the owning component's guide.
Authority: requested correction; normal PR update covered by the assignment.
Result/evidence: installed link resolves; declared source and packaged bytes agree.
Delivery: committed and pushed to the existing PR; reviewer owns acceptance.
```

For a new effect outside authority, prepare the concrete reviewable result first, then ask only
for the missing decision. Explain the exact applicable requirement when a fresh approval is needed.
Do not fabricate authority, weaken required gates, or perform unrelated cleanup. Retain unrelated
changes and undo only this change through its owner's recovery method when appropriate.

This is an L1 structured runbook recipe, validated by fresh-context walkthrough and owner-specific
technical checks. Do not promote it to a script, ledger or workflow engine without demonstrated need.
