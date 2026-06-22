Last updated: 2026-06-22T19:47:33Z (UTC)

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

## OK-PR-003 - Coordination Home Register ID Namespace

Type: implementation-plan
Purpose: Add managed plan-register doctrine so coordination-home entries use unique IDs for
member-repository plans while preserving source local register IDs.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/coordination-home-register-id-namespace/coordination-home-register-id-namespace_implementation_doc.md
Created: 2026-06-22
Last updated: 2026-06-22T19:47:33Z (UTC)
Priority / ordering note: Completed `v0.1.8` instruction-only release and isolated consumer sync
proof.

Relations:
- related: OK-PR-001 - Plan Register Session And Lifecycle Hardening
- related: docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_implementation_doc.md

Session refs:
- created: not recorded
- material update: 2026-06-22, not recorded, activated implementation plan and created sibling
  execution log.
- material update: 2026-06-22, not recorded, completed release-preparation epics through local
  asset build and validation; public release publication remains pending explicit approval.
- material update: 2026-06-22, not recorded, completed EP-04 local pre-publication checks; public
  tag, GitHub release, and consumer sync proof remain pending.
- material update: 2026-06-22, not recorded, completed `v0.1.8` public release publication and
  isolated consumer update-check, sync, and check proof.

Coordination note:
- local-only

## OK-PR-002 - Codeheart Operating Kit Implementation-Planning Quality

Type: implementation-plan
Purpose: Update managed discovery, implementation-planning, planning-review, and execution
workflows so plans preserve intended feature capability and reject policy-only or under-covered
implementation.
Status: completed
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_implementation_doc.md
Created: 2026-06-22
Last updated: 2026-06-22T18:50:19Z (UTC)
Priority / ordering note: Completed `v0.1.7` planning workflow quality release and first-consumer
sync proof.

Relations:
- related: first consumer repository discovery handoff -
  <first-consumer-repository>/docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_discovery_doc.md
- related: first consumer repository handoff implementation plan -
  <first-consumer-repository>/docs/repo/plans/codeheart-operating-kit-implementation-planning-quality/codeheart-operating-kit-implementation-planning-quality_implementation_doc.md

Session refs:
- created: 2026-06-22, session 019eef87-f252-7b91-aa50-ecf54b357c6c
- material update: 2026-06-22, session 019eef87-f252-7b91-aa50-ecf54b357c6c,
  completed `v0.1.7` release publication and first-consumer sync proof.

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
