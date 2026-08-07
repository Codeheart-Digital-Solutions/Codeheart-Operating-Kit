Last updated: 2026-08-07T00:26:02Z (UTC)

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

## Discovery Versions And Path Authority

Discovery v1 is the compatibility contract that recognizes formal records only beneath
`docs/repo/plans/`. Discovery v2 is the repository-wide contract. In v2, the authoritative local
path universe is regular or executable Markdown blobs in the Git index. The authoritative remote
universe is regular or executable Markdown blobs in the selected commit tree. A path is eligible
when its slash-normalized repository-relative segments contain an exact lowercase segment named
`docs` at any depth.

Examples of eligible v2 paths include:

```text
docs/repo/plans/example/example_discovery_doc.md
docs/business/procurement/vendor-review_implementation_doc.md
products/example/packages/api/docs/plans/migration/migration_discovery_doc.md
source/areas/payments/docs/initiatives/settlement/settlement_implementation_doc.md
```

`Docs/`, `mydocs/`, and `documentation/` are not `docs` segments. Classification is semantic
path-segment matching, not `*/docs/**`, and therefore includes root `docs/**` and arbitrarily
nested `**/docs/**` without depending on glob semantics.

The supported filename kinds are exact, case-sensitive suffixes:

- `*_discovery_doc.md` -> `discovery`
- `*_implementation_doc.md` -> `implementation`
- exact `README.md` with valid `kind: family` metadata -> `family`

Discovery and implementation candidates are selected by supported filename **or** one genuine
Codeheart plan metadata block. Candidate selection is intentionally broader than validity so a
misnamed or metadata-missing plan cannot disappear. Marker-like content inside Markdown fences,
examples, or malformed marker layouts is not genuine metadata authority and is reported when it
also has a filename signal.

In discovery v2 canonical mode, discovery and implementation records are valid only when both the
supported filename and valid matching metadata are present. Stable identity comes from metadata;
the human-readable kind comes from the filename. A metadata-bearing discovery or implementation
under an eligible `docs` tree but without the supported filename is a visible invalid candidate,
not a record. Lifecycle continues to come only from the document header.

A family record requires exact `README.md`, genuine valid family metadata, and meaningful child
plan grouping. Directory shape alone never creates family authority and an ordinary domain README
does not become a candidate merely because sibling plan directories appear. Existing v1 family
records remain v1 compatibility evidence until prospective v2 migration gives the README explicit
family metadata. Branch overlays must not infer new family authority from unchanged README bytes.

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

## Repository Ownership And Exclusions

In discovery v2, every ordinary tracked Markdown path with an exact `docs` segment is repository-
owned by default. There is no central docs-root allowlist and no `owned_roots` setting.

Hard-unowned boundaries never provide consumer plan authority:

- managed Kit content under `.codeheart/kit/` and local/user state under `.codeheart/local/` or
  `.codeheart/user/`;
- symlinks, gitlinks/submodules, and nested Git repositories;
- paths outside the repository; and
- ignored or untracked files as catalog authority.

Conventional fixture, vendor, generated, dependency, example, build, output, or virtual-
environment segments are ambiguous, not silently ignored. A plan signal beneath one becomes a
visible `prospective-blocked` candidate until the document is moved to a truthful owned docs path
or the misleading root is explicitly excluded. The implementation-defined ambiguity segment set
is emitted in the discovery policy and may be refined only from evidence.

Configure only reviewed repository-relative directory exclusions:

```yaml
component_settings:
  planning-workflows:
    plan_catalog_ownership:
      excluded_roots:
        - vendor/copied-product/docs/
```

Exclusions remain inventory evidence with their path, signal, disposition, and policy digest;
they do not silently disappear. Roots must be slash-normalized, unique, non-overlapping, portable,
and end in `/`. Globs, absolute paths, drive-qualified paths, `..`, case-fold collisions, and
repository-wide exclusions are invalid. If a genuine plan sits beneath a misleading conventional
boundary, move it normally; do not add `owned_roots` without a separate reviewed requirement.

`plans validate` and `plans inventory` accept `--include-untracked` only with a prospective or
active discovery-v2 read. This opt-in authoring preview is non-authoritative. It is excluded from
candidate-set and policy hashes, canonical completeness, migration writes, and remote evidence;
ignored files are not previewed.

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

Mixed mode is optional. When every prospective v2 candidate can receive valid metadata in one
reviewed migration, the normal route is legacy -> migrate every candidate -> canonical. Use mixed
only when proven pre-cutover filename-only records must be deferred. A frozen register and cutover
revision are required only for that grandfathering route. Do not create a second cutover revision
merely because discovery v2 expands the recognized plan set.

## Discovery V2 Activation

Existing repositories remain on discovery v1 until an explicit reviewed activation. Installing or
upgrading the Kit must not expand their active candidate set or invalidate their current catalog.
Freshly initialized repositories may default to v2.

Activate v2 only after prospective inventory, exclusions/ambiguity review, guarded migration, and
complete local plus required remote evidence. Store activation in shared config:

```yaml
component_settings:
  planning-workflows:
    plan_catalog_discovery_version: 2
    plan_catalog_mode: canonical
    plan_catalog_ownership:
      excluded_roots: []
```

Metadata schema remains version 1 because identity semantics do not change. Structured views,
inventories, migration ledgers, portfolio caches, and completeness evidence use version 2 where
their discovery/provenance contract changes. Preserve v1 artifacts as labeled historical evidence;
never report them as v2-complete.

## Derived Local Views

Use:

```sh
codeheart-operating-kit plans validate .
codeheart-operating-kit plans list --format text .
codeheart-operating-kit plans list --format json .
codeheart-operating-kit plans validate --target-discovery-version 2 \
  --target-catalog-mode canonical --json .
codeheart-operating-kit plans inventory --target-discovery-version 2 \
  --target-catalog-mode canonical --output <inventory.json> .
```

Text is for orientation. JSON is the stable structured view. A nonzero validation/list result
means the reported problems must be resolved or disclosed; do not present the view as clean.

The target flags perform prospective reads; they do not activate v2 or change repository config.
`plans inventory --output <path> .` records Git-backed migration evidence. It is not an approved
ledger and does not authorize plan writes. `plans migrate` accepts only a semantically reviewed
ledger and requires exactly one of `--dry-run` or `--yes`. Migration does not activate discovery
v2; config activation is a separately reviewed write after the projected catalog is complete.

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

Default-branch configuration controls discovery version and exclusions for remote scanning; a
feature branch cannot broaden authority. Default-branch records form the baseline. An accessible
unmerged branch adds only candidates whose semantic evidence is added or materially changed from
its merge base. Pure same-byte renames and inherited plans are not duplicated, and unchanged
family README bytes never gain authority from sibling-directory changes. Pull-request facts enrich
observations but never select them.

For a discovery-v2 member, remote readiness is always evaluated against canonical v2 validity,
even when its local catalog mode temporarily remains legacy or mixed. Excluded, ambiguous,
hard-unowned, malformed, and filename-only candidate observations remain structured evidence. A
required member that is v1, inaccessible, invalid, or lacks canonical v2 candidate coverage makes
the current portfolio attempt incomplete. Only a complete schema-v2/discovery-v2 scan may replace
the current complete cache; older v1 caches remain historical and byte-preserved.

Only pushed remote refs are portfolio facts. Worktree edits, local heads, and unpushed commits
remain local evidence until normally pushed.

## Compatibility And Validation

Legacy IDs and register rows are evidence to reconcile, not strings to convert mechanically. A
migration reviewer must inspect each plan's content, sibling records, relations, lifecycle, Git
ownership, and ambiguity. Preserve unclear optional classifications as omissions.

Validation blocks malformed or misplaced metadata (`metadata_*`), missing metadata
(`metadata_missing`), missing canonical filenames (`canonical_filename_missing`), path/metadata
kind disagreement (`record_kind_path_mismatch` or `plan_kind_mismatch`), duplicate IDs
(`duplicate_plan_id`), invalid family placement, portable path collisions, ambiguous conventional
roots (`plan_documentation_root_ambiguous`), unsafe sources (`plan_source_unsafe`), paths outside
owned docs (`plan_path_unowned` or `plan_metadata_misplaced`), invalid/overlapping
exclusions, broken mixed baselines, and incomplete canonical coverage.
Excluded candidates are retained as `plan_candidate_excluded` evidence rather than counted as
owned records. Never repair a problem by inventing metadata, weakening mode, silently extending an
exclusion, or treating preview content as authority.

Use `../runbooks/migrate-plan-catalog.md` for adoption and
`../runbooks/maintain-plan-register.md` for legacy and generated-view behavior.
