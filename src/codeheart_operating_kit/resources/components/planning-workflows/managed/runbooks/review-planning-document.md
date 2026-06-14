Last updated: 2026-06-13T22:55:57Z (UTC)

# Review Planning Document

Use this runbook to review discovery or implementation documents for architecture quality and
execution readiness.

## Review Areas

1. Decision soundness.
2. Completeness.
3. Ambiguity and risk.
4. Maintainability, testability, operability, and safety.
5. Execution readiness.

## Procedure

1. Read the target planning document.
2. Read referenced discovery, implementation, and architecture notes that materially affect it.
3. Check whether decisions map to requirements and constraints.
4. Check whether open questions identify blocker status and ownership.
5. Check whether implementation epics are realistic, linear, and validated.
6. Lead with findings ordered by severity.
7. Use concrete file and section references.
8. Keep summary secondary to findings.

## Severity

- `High`: can cause implementation failure, public safety issue, or major rework.
- `Medium`: likely to create ambiguity, missed work, or weak validation.
- `Low`: useful improvement with limited execution risk.
