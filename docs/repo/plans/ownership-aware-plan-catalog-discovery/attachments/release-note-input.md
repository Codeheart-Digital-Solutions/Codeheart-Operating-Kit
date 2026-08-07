Last updated: 2026-08-07T01:39:50Z (UTC)

# Repository-Wide Plan Catalog Discovery Release Note Input

This input is intentionally version-neutral. The release checkpoint selects the version from
current repository doctrine.

## Included

- Adds opt-in discovery v2 for formal plan records in Git-tracked, repository-owned Markdown
  documentation trees at arbitrary depth. Eligibility uses an exact lowercase `docs` path segment,
  so root `docs/**`, nested `product/docs/**`, and more deeply repeated documentation paths are all
  handled without fixed-depth glob semantics.
- Uses supported filename or genuine plan metadata to discover candidates. Canonical
  discovery/implementation validity requires both a supported filename and valid metadata;
  metadata owns stable identity, the filename or qualified family record owns human-readable kind,
  and the document header owns lifecycle.
- Uses one local-index/remote-commit-tree classifier across list, validate, inventory, migration,
  changed-path ownership, remote baselines, branch overlays, and portfolio completeness.
- Adds explicit ownership evidence: hard-unowned managed/local/user state, symlinks, submodules,
  nested repositories, external paths, and untracked authority; tracked conventional ambiguity
  paths remain visible blockers until moved or explicitly listed in `excluded_roots`.
- Requires exact `README.md` plus genuine valid family metadata for new discovery-v2 family
  authority. Directory shape alone cannot promote an ordinary domain router.
- Adds versioned structured plan views, inventory/ledger provenance, and portfolio-v2 completeness
  while preserving plan metadata schema v1.

## Compatibility

Existing repositories remain on discovery v1 after upgrade unless they explicitly activate v2.
Their fixed-root plan view, catalog mode, register, plan bytes, config, and historical caches are
not automatically migrated or rewritten. Fresh repositories default to discovery v2.

Mixed mode remains supported for intentionally grandfathered filename-only records, but it is not
mandatory. A repository that migrates every discovered formal record can move from legacy to
canonical v2 without a new cutover revision or frozen-register gap.

## Consumer Impact

The change includes managed instruction, validator, consumer-migration, and safety-policy impact.
No component, scaffold, managed path, or consumer-owned path moves. Activating discovery v2 may
surface plan-like tracked documents that were previously outside the active catalog; those findings
must be reviewed rather than silently ignored.

## Security And Ownership

Only regular blobs from the local Git index or selected remote commit tree can become catalog
authority. Markdown is parsed as inert bytes, and marker-like text in code fences, quotations, or
examples does not become genuine metadata. Default-branch policy controls remote scans, branch
content cannot broaden authority, and excluded content remains inventory evidence.

## Validation

The implementation includes a persistent multi-root fixture, focused classifier/metadata/mode/
family/ownership/migration/remote/cache/portability/security tests, a 100,000-path scale benchmark,
and Ubuntu, macOS, and real-Windows validation. Release readiness additionally requires two
byte-identical supported-platform builds and the full catalog-to-binary provenance chain.

## Rollout

Upgrade the Kit first; existing repositories stay on discovery v1. Before activation, run
prospective v2 inventory and validation, review ownership and exclusions, migrate every formal
candidate or explicitly select bounded mixed grandfathering, validate all required remote members,
then activate v2 through a reviewed config change. Preserve old caches and registers as labeled
historical evidence.

Codeheart-HQ receives the released Kit and proves that discovery v1 semantics remain unchanged.
Its plan metadata migration and canonical-v2 activation remain a separate consumer workstream.

## Deferred Product Decisions

- The later release, if any, that removes discovery v1.
- Evidence-based refinement of the frozen conventional ambiguity segment list.
