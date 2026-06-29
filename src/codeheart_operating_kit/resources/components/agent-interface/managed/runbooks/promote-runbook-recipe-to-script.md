Last updated: 2026-06-29T14:49:19Z (UTC)

# Promote Runbook Recipe To Script

Use this runbook when an existing or planned runbook recipe contains deterministic mechanics that
may be safer as a reusable script asset.

Audience: agent-facing

Intent:
Convert one fragile, repeated, or evidence-bearing recipe step into a reusable script asset while
keeping the runbook responsible for routing, approvals, target selection, consequences, and
fallback choice.

Success:
The selected recipe has an explicit promotion decision. When promoted, it has an owner, script
entrypoint, runbook caller, output contract, tests or fixtures, and review evidence.

Agent judgment boundary:
The agent may promote deterministic mechanics when the owning repo, module, package, or source
area is clear. It must not promote user conversation, approval decisions, policy judgment, or
ambiguous target selection into a script.

Stop boundary:
Stop before creating or changing a script when the owner, placement boundary, inputs, approval
model, target scope, output contract, or validation path is unclear.

## Required References

Read:

- `../reference/runbook-to-script-promotion-standard.md`;
- `../reference/operational-recipe-maturity.md`;
- `../reference/runbook-authoring-standard.md` when changing a runbook;
- `../reference/operation-routing-and-dispatch.md` when owner, scope, route, or execution surface
  is not already resolved;
- the owning repository, module, package, or source-area guidance for concrete script placement.

## Trigger

Use this runbook when:

- a runbook contains long inline implementation logic;
- command syntax, quoting, encoding, parsing, or request construction is fragile;
- a deterministic step repeats across runbooks or sessions;
- consistent output, markers, blockers, or run records are needed;
- a fresh agent repeatedly has to invent execution details;
- an implementation plan asks for reusable script assets or script-promotion review.

Do not use this runbook only because a runbook contains a short command invocation or example.

## Inputs

Collect:

- owning area;
- current runbook path;
- recipe or section being considered;
- intended script purpose;
- input values and accepted formats;
- read-only, dry-run, write, or mixed mode;
- approval requirements;
- target scope;
- expected output and blockers;
- test or fixture path.

## Procedure

1. Confirm owner and route.

   Identify the owning repository, module, product, package, or source area. Use routing guidance
   before choosing a script, CLI, API, portal, browser, connector, or manual surface.

2. Classify the recipe.

   Use the promotion standard to classify the selected work as one of:

   - runbook-only judgment;
   - structured runbook recipe;
   - reusable script asset;
   - thin command wrapper;
   - mature API or tool surface.

   Record the promotion or non-promotion decision in the plan, runbook, or execution log that owns
   the change.

3. Keep runbook-only work out of scripts.

   Do not promote user conversation, approval decision, consequence explanation, ambiguous routing,
   target selection, policy judgment, fallback choice, one-off exploration, or business decision
   capture.

4. Define the script contract.

   For a reusable script asset, define:

   - `script_id`;
   - required inputs;
   - optional inputs;
   - mode;
   - data shape;
   - blocker classes;
   - output safety behavior;
   - local tooling prerequisites;
   - approval reference handling;
   - target-scope handling.

5. Select placement.

   Use the owning area's placement convention. When no stricter owner convention exists, use:

   ```text
   <owner-root>/
     scripts/
       README.md
       <script-entrypoint>
     tests/
       scripts/
         <script-test>
       fixtures/
         scripts/
           <fixture>
   ```

6. Create or update `scripts/README.md`.

   Record script entrypoints, runbook callers, local tooling, input contract, output contract,
   safety boundary, mode behavior, tests, and fixtures.

7. Implement the script.

   Keep the script narrow. It should perform one deterministic operation or evidence step. It
   should not ask the user questions, broaden targets, hide approvals, or choose a different route.

8. Implement tests or fixtures.

   Prove required input validation, output contract, blocker classification, output safety flags,
   dry-run/read-only behavior when declared, and helper behavior that would otherwise be duplicated.

9. Update the runbook caller.

   Replace full inline implementation with:

   - script invocation path;
   - required inputs;
   - approval preconditions;
   - expected success output;
   - expected blocker output;
   - fallback and stop conditions.

   Keep short invocations or examples when useful. Do not keep a full duplicate implementation in
   the runbook.

10. Review for safety and maturity.

    Confirm:

    - runbook still owns approvals, target selection, consequences, fallback choice, and stop
      conditions;
    - script has a clear owner and placement boundary;
    - output contract is stable;
    - output safety describes emitted-output behavior;
    - tests or fixtures prove the contract;
    - package, CLI, wrapper, or API promotion is not premature.

11. Record evidence.

    Record changed files, validation commands, test results, stale-inline cleanup, and residual
    risk in the owning plan, execution log, or final task summary.

## Stop Conditions

Stop and return to planning or owner clarification when:

- owner or route is ambiguous;
- script would need to decide user approval or business policy;
- target scope cannot be expressed exactly;
- tests cannot be scoped;
- output would include uncontrolled raw secrets, raw provider responses, or raw external content;
- promotion would create a broad package, CLI, wrapper, or API before repeated use proves the
  surface.

## Review Checklist

- The step promoted is deterministic.
- The runbook caller remains the operator-facing entry point.
- Long inline implementation has been replaced by script invocation.
- `scripts/README.md` exists for the owner area.
- Script inputs and mode are explicit.
- Script output includes `status`, `script_id`, `mode`, `summary`, `blocker`, `data`, and
  `output_safety`.
- Tests or fixtures prove the output contract.
- Approval and target scope are not broadened inside the script.
- No raw secrets, raw provider responses, or uncontrolled raw external content are emitted in
  normal output.
- Wrapper, package, CLI, or API promotion has a repeated-use rationale.
