Last updated: 2026-06-21T15:11:52Z (UTC)

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
- Managed section boundaries: `reference/managed-section-boundaries.md`
- Local extension contract: `reference/local-extension-contract.md`
- Onboarding context contract: `reference/onboarding-context-contract.md`
- Update-check policy: `reference/update-check-policy.md`
- Kit feedback item format: `reference/kit-feedback-item-format.md`
- First-run onboarding: `runbooks/conduct-first-run-onboarding.md`
- Root `AGENTS.md` structure: `runbooks/structure-root-agents-md.md`
- Root `AGENTS.md` repair: `runbooks/repair-root-agents-md.md`
- Submit kit feedback: `runbooks/submit-kit-feedback.md`
- Installed fallback inventory: `.codeheart/kit/README.md`
