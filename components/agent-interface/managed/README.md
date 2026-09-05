Last updated: 2026-09-05T21:03:41Z (UTC)

# Agent Interface

This managed domain owns how an installed consumer introduces the Operating Kit to agents, how
local owners extend instructions, and how agents find managed docs without reading the full kit
inventory for every task.

The root `AGENTS.md` managed block routes semantic plan authoring, catalog migration, portfolio
setup, and refresh to planning workflows. It stays generic: repository identities, authorized
source scopes, strategic interpretation, and plan content belong in shared config, ignored local
state, and repository-owned documents, not in the managed bootstrap text.

## Routes

- Generic whole-plan assignments, delegated reviews and report-back: `reference/agent-task-coordination.md`
- Optional Codex task and goal operations: `reference/codex-task-operations.md`
- Root `AGENTS.md` contract: `reference/root-agents-md-contract.md`
- Operation routing and dispatch: `reference/operation-routing-and-dispatch.md`
- Managed section boundaries: `reference/managed-section-boundaries.md`
- Local extension contract: `reference/local-extension-contract.md`
- Onboarding context contract: `reference/onboarding-context-contract.md`
- Operational recipe maturity: `reference/operational-recipe-maturity.md`
- Runbook-to-script promotion standard: `reference/runbook-to-script-promotion-standard.md`
- Runbook authoring standard: `reference/runbook-authoring-standard.md`
- Optional manual update-check policy: `reference/update-check-policy.md`
- Kit feedback item format: `reference/kit-feedback-item-format.md`
- Repo feedback item format: `reference/repo-feedback-item-format.md`
- First-run onboarding: `runbooks/conduct-first-run-onboarding.md`
- Installation lifecycle—init, repair, sync, update-check, upgrade, and check:
  `runbooks/maintain-operating-kit-installation.md`
- Tooling readiness: `runbooks/handle-tooling-readiness.md`
- Promote runbook recipe to script: `runbooks/promote-runbook-recipe-to-script.md`
- Root `AGENTS.md` structure: `runbooks/structure-root-agents-md.md`
- Root `AGENTS.md` repair: `runbooks/repair-root-agents-md.md`
- Submit kit feedback: `runbooks/submit-kit-feedback.md`
- Capture repo feedback: `runbooks/capture-repo-feedback.md`
- Enable GitHub Issues feedback intake: `runbooks/enable-github-issues-feedback-intake.md`
- Installed fallback inventory: `.codeheart/kit/README.md`
- Semantic plan catalog and portfolio coordination:
  `../planning-workflows/README.md`
