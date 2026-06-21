Last updated: 2026-06-21T14:53:02Z (UTC)

# Plan Register Format

Use this reference for `docs/repo/plans/plan-register.md` in consumer repositories and
coordination-home repositories.

A plan register is a lightweight, durable index of important planning and workstream authority. It
helps agents and humans find canonical planning documents, understand relationships between plans,
and recover session IDs for plan creation or material plan changes.

The register does not replace discovery documents, implementation plans, execution logs, session
ledgers, runbooks, or task trackers.

## Source Of Truth

Canonical planning documents own their own details.

The register may copy lifecycle metadata such as status, creation date, and last-updated timestamp
as an index snapshot. When the register and the canonical document disagree, the canonical
document wins and the register should be refreshed.

## Coverage

A local repository register should list important local discovery plans, implementation plans,
plan families, major workstreams, and selected cross-repository dependencies that materially
affect local work.

A coordination-home register uses the same entry shape but may be selective. It should cover
portfolio-level plan families and selected member-repository entries that matter to portfolio
coordination.

## Entry Fields

Use one repeated Markdown section per entry.

Required fields:

- `ID`: stable local register identifier or canonical plan ID.
- `Title`: human-readable title.
- `Type`: one of `discovery-plan`, `implementation-plan`, `plan-family`, `workstream`, or
  `reference-index`.
- `Purpose`: short explanation of what the entry covers.
- `Owner / repository`: owning person, role, team, or repository.
- `Canonical docs`: repo-relative path or explicit repository/path pointer to the authoritative
  document or bundle.
- `Status`: lifecycle snapshot copied from the canonical document when available.
- `Created`: creation date copied from the canonical document when available.
- `Last updated`: last-updated timestamp copied from the canonical document when available.
- `Priority / ordering note`: stable orientation note, not a volatile next action.
- `Relations`: parent, child, supersession, dependency, blocking, or related links.
- `Session refs`: lightweight session IDs for creation and material plan updates when available.
- `Coordination note`: optional note about portfolio relevance or sync state.

Optional fields:

- `Coverage note`: explains whether the register is complete for a repository, portfolio, plan
  family, or intentionally selective.
- `Last reviewed`: last explicit register review date.
- `Sync state`: local-only marker such as `local-only`, `synced`, or
  `coordination-sync-pending`.

## Lifecycle Values

Use these lifecycle values when copying status from canonical planning documents:

- `draft`
- `active`
- `completed`
- `superseded`
- `archived`

Do not invent additional status values in the register. If a canonical document uses a different
local status, preserve the canonical value in the document and use the closest standard lifecycle
snapshot in the register.

## Relation Vocabulary

Use these relation terms:

- `parent`: broader plan, program, family, or workstream that owns this entry.
- `child`: subordinate plan or workstream.
- `supersedes`: older entry replaced by this entry.
- `superseded-by`: newer entry replacing this entry.
- `depends-on`: entry that must happen or remain true before this one can proceed.
- `blocks`: entry materially blocked by this one.
- `related`: relevant but non-blocking relationship.

Prefer stable IDs and relative paths over prose-only references.

## Session References

Session refs are recovery handles, not summaries.

Record the creating session and material-update sessions when a session ID is available. A session
reference should contain only:

- reason such as `created` or `material update`;
- session ID;
- date;
- optional short note naming the material change.

Do not add session summaries, transcript excerpts, or detailed status narratives to the register.
When no session ID is available, omit the row or record `not recorded`; missing session IDs do not
block planning work.

## Coordination Notes

Use `Coordination note` for portfolio-level relevance and sync state that helps a future agent
choose the right register updates.

Examples:

- `local-only`
- `candidate for coordination-home register`
- `synced to coordination home`
- `coordination-sync-pending`

Detailed pending sync belongs in `docs/repo/plans/coordination-sync-pending.md`, not in the main
register entry.

## Markdown Shape

Use this repeated-section shape:

```md
## PR-001 - Example Portfolio Coordination Model

Type: discovery-plan
Purpose: Define a reusable coordination model for related repositories.
Status: active
Owner / repository: Example-Automation
Canonical docs: docs/repo/plans/example-portfolio/example-portfolio_discovery_doc.md
Created: 2026-06-21
Last updated: 2026-06-21T14:53:02Z (UTC)
Priority / ordering note: Needed before implementation planning.

Relations:
- parent: PF-001 - Example Portfolio Foundation
- related: PR-002 - Example Agent Memory Model

Session refs:
- created: not recorded
- material update: 2026-06-21, session not recorded, selected repeated-section format.

Coordination note:
- candidate for coordination-home register
```

Keep entries compact. If a field would need multiple paragraphs, move the detail to the canonical
plan and link to it.

## Anti-Patterns

Do not use the register as:

- a task backlog;
- a sprint board;
- a per-epic progress table;
- a transcript or session-summary index;
- the source of truth for lifecycle state;
- a duplicate copy of discovery or implementation plan details;
- a place for private tenant, customer, credential, or local-machine information;
- a mechanism for silently discovering or writing sibling repositories.

