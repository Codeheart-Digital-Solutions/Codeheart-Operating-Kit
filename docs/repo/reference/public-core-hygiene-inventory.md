Last updated: 2026-06-13T22:55:57Z (UTC)

# Public-Core Hygiene Inventory

This inventory records the first G1 extraction inputs and their public-core handling.

## Included And Generalized Inputs

- `feature-discovery-doc.SKILL.md.txt`: generalized into planning discovery workflow; no local
  Codex skill authority copied.
- `feature-implementation-doc.SKILL.md.txt`: generalized into implementation-planning workflow and
  planning lifecycle reference.
- `feature-document-review.SKILL.md.txt`: generalized into planning document review workflow.
- `intention-led-discovery.md.txt`: consolidated into the discovery workflow as decision-led and
  goal-style discovery guidance.
- `implementation-plan-execution.md.txt`: generalized into implementation execution workflow and
  execution-log guidance.
- `agent-memory.README.md.txt`: generalized into agent-memory domain routing and scaffold
  boundaries.
- `agent-memory.entry-format.md.txt`: generalized into reusable memory entry formats and status
  vocabulary.
- `agent-memory.session-ledger-maintenance.md.txt`: generalized into memory maintenance runbook.
- `agent-memory-system_discovery_doc.md.txt`: used as historical input for memory boundaries.
- `AGENTS.md.txt`: generalized into root agent-interface contracts and managed block template.
- `agent-guidance-cleanup_discovery_doc.md.txt`: used as historical input for lean root
  `AGENTS.md` routing.
- `documentation-structure.md.txt`: generalized into structure-governance documentation placement
  rules and index maintenance.
- `repository-information-architecture.md.txt`: generalized into naming, ownership, and boundary
  guidance.
- `docs.README.md.txt` and `repo.README.md.txt`: generalized into router and index-maintenance
  examples.

## Architecture Decision Inputs

- `agent-interface-contract.md`: included as managed agent-interface contract.
- `managed-content-ownership-model.md`: included as managed/scaffold/template boundary guidance.
- `native-codex-capability-baseline.md`: included as external capability baseline guidance.
- `onboarding-and-update-policy.md`: included as onboarding and update-check policy guidance.
- `first-run-onboarding-script.md`: used as the canonical source for onboarding copy and prompt
  order.

## Excluded Inputs

- Live `docs/agent-memory/goal-register.md`, `session-ledger.md`, and `untriaged-sessions.md` are
  excluded as consumer-owned state.
- Product-specific API references, cloud deploy rules, tenant details, customer details, business
  records, credentials, secrets, instance records, and restricted strategy content are excluded.
- `docs/repo/reference/public-apis.md` is not copied as universal Operating Kit content.

## Sanitization Result

The G1 content baseline uses public, generic placeholders and reusable operating language. It does
not include product-specific AWS deploy guardrails, account identifiers, tenant identifiers,
customer names, credentials, secrets, or private business strategy.
