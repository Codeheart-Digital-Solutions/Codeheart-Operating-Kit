Last updated: 2026-07-31T22:23:19Z (UTC)

# Plan Catalog Format

Use this reference for canonical discovery, implementation, and family records and for the
repository-scoped views derived from them. The machine contract is
`schemas/plan-metadata.schema.json`; installed consumers use the matching schema shipped with
their Operating Kit version.

## Authority Model

- A canonical plan document owns its title, lifecycle header, content, and bounded metadata.
- Git owns revision, ref, path, and content-hash evidence.
- `codeheart-operating-kit plans list` derives a current repository view; it does not create a
  committed catalog.
- `plan-register.md` is frozen legacy evidence after mixed cutover, not current authoring
  authority.
- A coordination home derives portfolio facts from pushed remote refs and keeps strategic
  interpretation in its separate repo-owned overlay.

Never copy source-derived cache content into a plan as authority. Never use an execution log as an
independent plan record unless a later reviewed contract explicitly adds that kind.

## Canonical Record Kinds And Paths

Formal records are discovered under `docs/repo/plans/`:

- `*_discovery_doc.md` -> `discovery`
- `*_implementation_doc.md` -> `implementation`
- a qualifying nested family `README.md` -> `family`

The root `docs/repo/plans/README.md` is an index, not a family. A nested family README qualifies
when its directory contains at least two child directories that contain formal discovery or
implementation documents.

Discovery and implementation are separate records even when they belong to the same feature or
family. A family is a discoverability and relationship record; it is not an extra lifecycle role
on every member.

## Canonical Header And Metadata Position

Keep the lifecycle header first:

```text
Last updated: YYYY-MM-DDTHH:MM:SSZ (UTC)
Created: YYYY-MM-DD
Status: draft | active | completed | superseded | archived
```

Place exactly one metadata block immediately below the canonical Markdown title and outside all
other code fences:

````md
# Example Discovery

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: example.discovery.requirements
  kind: discovery
  purpose: Resolve the public-safe requirements for the example capability.
  first_cataloged: 2026-07-31T12:00:00Z
  catalog_metadata_updated: 2026-07-31T12:00:00Z
  family: example.family.example-capability
  products:
    - example-product
  capabilities:
    - example-capability
  strategic_themes:
    - reliable-coordination
  relations:
    - kind: related
      target: example.implementation.delivery
  legacy_aliases:
    - PR-012
```
<!-- END CODEHEART PLAN METADATA -->
````

The visible YAML fence is intentional. Marker-like text inside examples or other Markdown fences
does not count as metadata.

## Required Metadata

`schema_version`
: Integer `1`.

`id`
: Stable semantic ID in `<namespace>.<kind>.<slug>` form. The middle segment must match
  `discovery`, `implementation`, or `family`.

`kind`
: The record kind and filename/family placement must agree.

`purpose`
: One non-empty public-safe line explaining why the record exists. Do not copy a title as a
  substitute for semantic review.

`first_cataloged`
: UTC RFC3339 time ending in `Z`. Set once when canonical metadata is first adopted and preserve
  it through renames and later metadata edits.

`catalog_metadata_updated`
: UTC RFC3339 time ending in `Z`. Change only when metadata meaning changes.

## Optional Metadata

`family`
: Canonical family plan ID. Family membership is derived from member records that name it; do not
  maintain a second member array in the family metadata.

`products`, `capabilities`, `strategic_themes`
: Unique public identifiers. Use lowercase letters, digits, dots, underscores, and hyphens. Omit
  classifications that evidence does not support.

`relations`
: Unique `{kind, target}` objects. Kinds are `parent`, `child`, `supersedes`, `superseded-by`,
  `depends-on`, `blocks`, and `related`. Targets are canonical semantic IDs.

`legacy_aliases`
: Stable historical identifiers such as `PR-012`. Preserve aliases; do not use them as the new
  canonical ID and do not renumber them to remove collisions.

## Stable Identity And Families

Choose IDs from meaning, not allocation order, branch name, path, repository display name, or
current status. A rename does not change the ID. Parallel branches may independently create plans
without a shared number allocator because the semantic ID is part of each document.

When two branches choose the same semantic ID, keep both observations with their repository, ref,
commit, path, and content hash. Report the conflict; do not silently choose one.

Use `family` when independently useful discovery or implementation records should be found and
analyzed together. Use relations for directional meaning. Merely sharing a folder or topic does
not require a family record.

## Lifecycle Versus Authority

Lifecycle is the single `Status:` value in the document header:

- `draft`: not execution authority
- `active`: approved current execution authority
- `completed`: executed traceability
- `superseded`: replaced by another record
- `archived`: historical context only

Do not add a separate authority field. User approval activates an implementation plan by changing
its lifecycle to `active`; revocation or replacement uses the existing lifecycle and relation
model. Repository policy and current user instruction still bound what an active plan may do.

## Content Chronology

`Last updated` describes meaningful plan-content or lifecycle change. Metadata-only migration must
leave that historical value byte-identical. `catalog_metadata_updated` records metadata adoption
or metadata-semantic change instead.

Mechanical migration must also preserve `Created`. Later content edits update `Last updated` in
the normal planning workflow.

## Repository Catalog Modes

Mode is stored at:

```yaml
component_settings:
  planning-workflows:
    plan_catalog_mode: legacy | mixed | canonical
```

`legacy`
: Default when the setting is absent. Existing formal plans may lack metadata and the register may
  still be maintained under legacy doctrine.

`mixed`
: New or materially rewritten formal plans require canonical metadata. Only pre-cutover plans
  proven by the frozen register and exact baseline commit may remain legacy temporarily. Mixed
  mode also requires `portfolio.member_repository_id` and:

```yaml
    plan_catalog_cutover_revision: <full pre-cutover commit>
```

That commit must be a reachable ancestor, must still be legacy mode, and must contain the regular
frozen register and grandfathered plan files. Current register bytes must remain identical to the
baseline.

`canonical`
: Every formal record requires valid metadata. The legacy register may remain as historical
  evidence but no longer fills canonical gaps.

Do not enter mixed mode until the frozen legacy baseline is committed. Do not enter canonical
mode until current default-branch and accessible remote-overlay validation is complete.

## Derived Local Views

Use:

```sh
codeheart-operating-kit plans validate .
codeheart-operating-kit plans list --format text .
codeheart-operating-kit plans list --format json .
```

Text is for orientation. JSON is the stable structured view. A nonzero validation/list result
means the reported problems must be resolved or disclosed; do not present the view as clean.

`plans inventory --output <path> .` records Git-backed migration evidence. It is not an approved
ledger and does not authorize plan writes. `plans migrate` accepts only a semantically reviewed
ledger and requires exactly one of `--dry-run` or `--yes`.

## Recipe Maturity, Evidence, And Blockers

Plan activation and plan-checkpoint publication remain an L1 structured runbook recipe. Validate
that route with a fresh low-context agent and static positive/negative authority checks. Do not add
an activation, commit, push, PR, or merge command.

`plans validate`, `plans list`, `plans inventory`, `plans migrate`, `portfolio configure`, and
`portfolio scan` are L3 thin command wrappers: they validate inputs, execute bounded mechanics,
return stable text/JSON or transaction output, and require command/interface/fixture evidence.
Their calling runbooks remain the operator entry points.

No L2 reusable script is introduced between those runbooks and commands. Do not add one merely to
wrap an existing command. Consider a separate L2 promotion only when new repeated mechanics cannot
be expressed safely by the existing L3 interface and a reviewed owner, contract, tests, and
placement exist.

Non-secret evidence fields include repository ID, plan ID/path, mode, ref/branch, source revision,
content hash, inventory/ledger path, action/skip/blocker code, included planning paths, commit,
normal-push result, validation result, scan times/completeness/metrics, and cache update or
preservation. Do not record tokens, credential-bearing URLs, raw sensitive logs, unnecessary local
paths, or private topology.

Use this blocker shape in run records: blocker code/class, recipe step, repository/plan target,
non-secret message, evidence command/result, preserved state, safe retry/recovery path, and user or
branch-owner decision needed. A blocker must not be converted into a fabricated success row.

## Source Observations

A coordination catalog observation identifies at least:

- repository ID, plan ID, title, kind, purpose, and lifecycle;
- canonical path, ref, commit, and SHA-256 content identity;
- visibility (`default` or `unmerged-branch`);
- verification, observed time, stale state, and conflict state; and
- optional family, classifications, relations, aliases, source-change time, and pull-request fact.

Default-branch records form the baseline. An accessible unmerged branch adds only canonical plan
files whose bytes are added or materially changed from its merge base. Inherited plans are not
duplicated. Pull-request facts enrich observations but never select them.

Only pushed remote refs are portfolio facts. Worktree edits, local heads, and unpushed commits
remain local evidence until normally pushed.

## Compatibility And Validation

Legacy IDs and register rows are evidence to reconcile, not strings to convert mechanically. A
migration reviewer must inspect each plan's content, sibling records, relations, lifecycle, Git
ownership, and ambiguity. Preserve unclear optional classifications as omissions.

Validation blocks malformed metadata, kind/ID disagreement, duplicate IDs, unsafe or unreadable
sources, missing required identity, broken mixed baselines, incomplete canonical coverage, and
other stable problem codes. Never repair these by inventing metadata or weakening the configured
mode.

Use `../runbooks/migrate-plan-catalog.md` for adoption and
`../runbooks/maintain-plan-register.md` for legacy and generated-view behavior.
