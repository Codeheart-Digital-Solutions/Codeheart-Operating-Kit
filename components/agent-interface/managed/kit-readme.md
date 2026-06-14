Last updated: 2026-06-13T23:55:46Z (UTC)

# Codeheart Operating Kit Inventory

This file is the installed fallback inventory for managed Operating Kit content.

Agents should prefer the direct routes in the root `AGENTS.md` managed block. Use this inventory
when a task is unclear, the direct route does not match, the user asks about kit structure, or kit
state looks missing, stale, or damaged.

## Managed Domains

- Planning workflows: `.codeheart/kit/docs/planning-workflows/README.md`
- Agent memory: `.codeheart/kit/docs/agent-memory/README.md`
- Agent interface: `.codeheart/kit/docs/agent-interface/README.md`
- Structure governance: `.codeheart/kit/docs/structure-governance/README.md`
- Native Codex capabilities: `.codeheart/kit/docs/native-codex-capabilities/README.md`

## Generated State

- Lock metadata: `.codeheart/kit.lock.yaml`
- Shared non-secret config: `.codeheart/kit.config.yaml`
- Local user layer: `.codeheart/user/`

## Operating Rule

Treat files under `.codeheart/kit/` as managed Operating Kit content. Repair or refresh them with
`codeheart-operating-kit sync` instead of hand-editing them.
