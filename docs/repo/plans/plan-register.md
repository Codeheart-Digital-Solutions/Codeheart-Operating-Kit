Last updated: 2026-06-22T18:05:35Z (UTC)

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

Coverage note: This register currently lists public Operating Kit repository plans that have been
entered during plan-register adoption. Earlier repository plans may be added during later register
maintenance.

## Entries

## OK-PR-002 - Codeheart Operating Kit Implementation-Planning Quality

Type: implementation-plan
Purpose: Update managed discovery, implementation-planning, planning-review, and execution
workflows so plans preserve intended feature capability and reject policy-only or under-covered
implementation.
Status: active
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_implementation_doc.md
Created: 2026-06-22
Last updated: 2026-06-22T18:05:35Z (UTC)
Priority / ordering note: Targets `v0.1.7` planning workflow quality release and first consumer repository
first-consumer sync proof.

Relations:
- related: first consumer repository discovery handoff -
  <first-consumer-repository>/docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_discovery_doc.md
- related: first consumer repository handoff implementation plan -
  <first-consumer-repository>/docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_implementation_doc.md

Session refs:
- created: 2026-06-22, session 019eef87-f252-7b91-aa50-ecf54b357c6c

Coordination note:
- local-only

## OK-PR-001 - Plan Register Session And Lifecycle Hardening

Type: implementation-plan
Purpose: Harden Operating Kit plan-register doctrine for self-contained session-reference
resolution and lifecycle grouping.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/plan-register-session-lifecycle-hardening/plan-register-session-lifecycle-hardening_implementation_doc.md
Created: 2026-06-21
Last updated: 2026-06-21T19:30:29Z (UTC)
Priority / ordering note: Prepares the next additive hardening release after `v0.1.5`.

Relations:
- related: Portfolio coordination and plan-register implementation plan -
  docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_implementation_doc.md

Session refs:
- material update: 2026-06-21, not recorded, activated implementation plan.
- material update: 2026-06-21, not recorded, completed implementation plan.

Coordination note:
- local-only
