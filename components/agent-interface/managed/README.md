Last updated: 2026-07-09T23:30:00Z (UTC)

# Agent Interface

This managed domain owns how an installed consumer introduces the Operating Kit to agents, how
local owners extend instructions, and how agents find managed docs without reading the full kit
inventory for every task.

The root `AGENTS.md` managed block may route configured portfolio coordination to the planning
workflow register-maintenance runbook. It stays generic: repository-specific coordination paths,
membership decisions, and planning state belong in `.codeheart/kit.config.yaml` and local
repository documents, not in the managed bootstrap text.

## Routes

- Root `AGENTS.md` contract: `reference/root-agents-md-contract.md`
- Operation routing and dispatch: `reference/operation-routing-and-dispatch.md`
- Managed section boundaries: `reference/managed-section-boundaries.md`
- Local extension contract: `reference/local-extension-contract.md`
- Onboarding context contract: `reference/onboarding-context-contract.md`
- Operational recipe maturity: `reference/operational-recipe-maturity.md`
- Runbook-to-script promotion standard: `reference/runbook-to-script-promotion-standard.md`
- Runbook authoring standard: `reference/runbook-authoring-standard.md`
- Update-check policy: `reference/update-check-policy.md`
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
