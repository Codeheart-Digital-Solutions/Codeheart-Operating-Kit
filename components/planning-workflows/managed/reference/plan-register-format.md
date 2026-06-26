Last updated: 2026-06-26T14:30:34Z (UTC)

# Plan Register Format

Use this reference for `docs/repo/plans/plan-register.md` in consumer repositories and
coordination-home repositories.

A plan register is a lightweight, durable index of important planning and workstream authority. It
helps agents and humans find canonical planning documents, understand relationships between plans,
and recover session IDs for plan creation or material plan changes.

The register does not replace discovery documents, implementation plans, execution logs, runbooks,
or task trackers. Session-reference guidance for formal plan registration lives in this planning
workflow reference and its maintenance runbook.

## Source Of Truth

Canonical planning documents own their own details.

The register may copy lifecycle metadata such as status, creation date, and last-updated timestamp
as an index snapshot. When the register and the canonical document disagree, the canonical
document wins and the register should be refreshed.

The register owns compact index metadata and stable relationship pointers. Keep decision
rationale, blockers, execution evidence, review findings, detailed dependencies, and active work
state in the canonical planning document or execution log.

## Register Location And Canonical Documents

Keep one `docs/repo/plans/plan-register.md` as the default durable register for the repository.
The register location does not require every canonical planning document to live under
`docs/repo/plans/`.

Canonical planning documents may live in any repository-owned planning root, including
`docs/repo/plans/`, `docs/business/plans/`, product docs, module docs, or source-area planning
folders. The register points to those canonical homes instead of moving or duplicating them.

When work spans multiple repositories, the canonical planning document should live in the
repository that owns the work boundary, implementation authority, and validation evidence.
Coordination-home registers represent that work with explicit pointers; they do not become the
canonical home for member-repository plans by default.

Use repo-relative paths for canonical documents owned by the same repository. Use explicit
repository/path pointers for member-repository entries represented in a coordination-home
register, such as `Example-Automation:docs/repo/plans/example/example_implementation_doc.md`.

## Coverage

A local repository register should list important local discovery plans, implementation plans,
plan families, major workstreams, and selected cross-repository dependencies that materially
affect local work.

A coordination-home register uses the same entry shape but may be selective. It should cover
portfolio-level plan families, local coordination-home plans, and selected member-repository
entries that matter to portfolio coordination. It is not limited to cross-repository entries; it
may also give a compact overview of what each represented repository is planning, implementing,
and completing.

## Entry Fields

Use one repeated Markdown section per entry.

Required fields:

- `ID`: stable local register identifier, coordination-home identifier, or canonical plan ID.
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

## ID Conventions

Standalone repository registers may use short local IDs such as `PR-001` when those IDs are clear
inside that one repository.

Coordinated portfolios may use repository-qualified local IDs such as
`EXAMPLE-AUTOMATION-PR-001` when the same entry should be easy to recognize in both the local
register and the coordination-home register. Repository-qualified local IDs are preferred when a
consumer repository expects frequent portfolio coordination and wants to avoid mental mapping
between local and coordination-home IDs.

Do not force existing consumer-owned registers to migrate only because this convention exists.
Apply it when creating new coordinated entries or when a consumer intentionally refreshes register
IDs through a local plan.

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

## Coordination-Home ID Uniqueness

Register IDs must be unique inside the register that contains them.

Local register IDs are local to the owning repository. A standalone member repository may use
`PR-001` in its own local register even when another member repository also has a local `PR-001`.

Coordination-home register IDs must be unique inside the coordination-home register. When adding a
member-repository entry to a coordination-home register, do not copy a bare member-local ID such
as `PR-001` as the coordination-home ID. Use a coordination-home ID that includes a stable source
namespace plus the source local ID.

Derive the namespace from `portfolio.member_repository_id` when present. Normalize it for register
IDs by uppercasing letters, replacing runs of non-alphanumeric characters with one hyphen, and
trimming leading or trailing hyphens. If `portfolio.member_repository_id` is unavailable, use the
`Owner / repository` value with the same normalization.

If the member local ID is already repository-qualified with that normalized source namespace, the
coordination home may reuse that ID. Do not double-prefix an already-qualified ID.

Examples:

- local member ID: `PR-001`
- member repository ID: `Example-Automation`
- coordination-home ID: `EXAMPLE-AUTOMATION-PR-001`
- already-qualified local member ID: `EXAMPLE-AUTOMATION-PR-001`
- coordination-home ID when already qualified: `EXAMPLE-AUTOMATION-PR-001`

Preserve the source local ID in `Coordination note`:

```md
Coordination note:
- Source local register ID: PR-001
- synced to coordination home
```

Use coordination-home IDs in coordination-home relations when the related entry is represented in
the coordination-home register. Use explicit repository/path pointers for related plans that are
not represented as coordination-home entries.

## Lifecycle Grouping

Keep one `docs/repo/plans/plan-register.md` as the default durable register. Do not create a
separate archive register by default.

Recommended grouping:

```md
## Active And Draft Entries

### PR-001 - Current Implementation Plan

## Completed Entries

### PR-002 - Completed Discovery Plan

## Superseded And Archived Entries

### PR-003 - Superseded Planning Model
```

Use grouping when the register has enough entries that lifecycle sections make scanning easier.
For very small registers, a single `## Entries` section is acceptable. When lifecycle grouping is
used, entry headings should sit below the grouping heading. The entry fields do not change.

Move entries between lifecycle groups when the canonical document status changes. The canonical
planning document remains the source of truth; grouping is only an index convenience.

## Repository Grouping

Coordination-home registers may group entries by owning repository when that improves scanning.
Repository grouping is an index convenience, not a separate ownership model.

Example:

```md
## Codeheart Operating Kit Entries

### OPERATING-KIT-PR-001 - Example Managed Workflow Plan

## Example Automation Entries

### EXAMPLE-AUTOMATION-PR-001 - Example Member Implementation Plan
```

Use repository grouping when ownership is the stable scanning boundary. Keep plan families,
workstreams, and lifecycle status in fields and relations rather than making every fuzzy grouping
a durable section.

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

Relations in the register are compact pointers. Put detailed dependency rationale, sequencing
notes, decision history, and implementation evidence in canonical planning documents and execution
logs.

## Session References

Session refs are recovery handles, not summaries.

Record the creating session and material-update sessions when a session ID is available. A session
reference should contain only:

- reason such as `created` or `material update`;
- session ID;
- date;
- optional short note naming the material change.

Do not add session summaries, transcript excerpts, or detailed status narratives to the register.
When no session ID is available, omit the row or record one of the explicit fallback values below;
missing session IDs do not block planning work.

Use these compact forms:

```md
Session refs:
- created: 2026-06-21, session <session-id>
- material update: 2026-06-21, session <session-id>, activated implementation plan.
- created: not recorded
- material update: 2026-06-21, ambiguous: multiple matching user sessions in this repository.
- material update: 2026-06-21, not confidently identified, metadata did not isolate one user session.
```

Fallback meanings:

- `not recorded`: no session ID was available or no session scan was performed.
- `ambiguous: <reason>`: bounded metadata or filename-only checks found more than one plausible
  user session.
- `not confidently identified`: metadata existed, but it did not support a confident match to the
  current user session.

## Coordination Notes

Use `Coordination note` for portfolio-level relevance and sync state that helps a future agent
choose the right register updates.

Examples:

- `local-only`
- `candidate for coordination-home register`
- `synced to coordination home`
- `coordination-sync-pending`
- `Source local register ID: PR-001`

Detailed pending sync belongs in `docs/repo/plans/coordination-sync-pending.md`, not in the main
register entry.

## Local Entry Shape

Use this repeated-section shape for a local repository entry:

```md
## Active And Draft Entries

### EXAMPLE-AUTOMATION-PR-001 - Example Local Implementation Plan

Type: implementation-plan
Purpose: Implement a repository-local capability that is relevant to portfolio coordination.
Status: active
Owner / repository: Example-Automation
Canonical docs: docs/repo/plans/example-local/example-local_implementation_doc.md
Created: 2026-06-21
Last updated: 2026-06-21T14:53:02Z (UTC)
Priority / ordering note: Needed before implementation planning.

Relations:
- parent: EXAMPLE-AUTOMATION-PF-001 - Example Local Plan Family
- related: EXAMPLE-AUTOMATION-PR-002 - Example Dependency Model

Session refs:
- created: not recorded
- material update: 2026-06-21, not recorded, selected repeated-section format.

Coordination note:
- candidate for coordination-home register
```

## Coordination-Home Member Entry Shape

Use this repeated-section shape when a coordination-home register represents a member-repository
entry:

```md
## Example Automation Entries

### EXAMPLE-AUTOMATION-PR-001 - Example Local Implementation Plan

Type: implementation-plan
Purpose: Implement a member-repository capability that affects portfolio coordination.
Status: active
Owner / repository: Example-Automation
Canonical docs: Example-Automation:docs/repo/plans/example-local/example-local_implementation_doc.md
Created: 2026-06-21
Last updated: 2026-06-21T14:53:02Z (UTC)
Priority / ordering note: Needed before the portfolio can depend on the member capability.

Relations:
- depends-on: OPERATING-KIT-PR-001 - Example Managed Workflow Plan
- related: Example-Automation:docs/repo/plans/example-dependency/example-dependency_discovery_doc.md

Session refs:
- created: not recorded

Coordination note:
- Source local register ID: PR-001
- synced to coordination home
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
