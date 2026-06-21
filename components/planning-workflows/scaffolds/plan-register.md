Last updated: 2026-06-21T14:53:02Z (UTC)

# Plan Register

This kit-initialized consumer state file lists important formal plans, plan families, major
workstreams, and portfolio-relevant planning records for this repository.

Operating Kit owns the file contract and format. This repository owns the entries after creation.
Sync may recreate this baseline when the file is absent, but it must not overwrite existing
entries.

Follow `.codeheart/kit/docs/planning-workflows/reference/plan-register-format.md` for entry
fields and `.codeheart/kit/docs/planning-workflows/runbooks/maintain-plan-register.md` for
maintenance.

## Register Coverage

Coverage note: No formal plan entries recorded yet.

## Entries

Add one repeated section per formal plan or workstream:

```md
## PR-001 - Example Plan Title

Type: discovery-plan
Purpose: Short purpose of this plan or workstream.
Status: draft
Owner / repository: Example-Automation
Canonical docs: docs/repo/plans/example/example_discovery_doc.md
Created: 2026-06-21
Last updated: 2026-06-21T14:53:02Z (UTC)
Priority / ordering note: Stable orientation note, not a volatile next action.

Relations:
- related: PR-002 - Example Related Plan

Session refs:
- created: not recorded

Coordination note:
- local-only
```

