Last updated: 2026-08-06T19:10:19Z (UTC)
Created: 2026-08-06
Status: draft

# Ownership-Aware Repository-Wide Plan Catalog Discovery

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.discovery.ownership-aware-plan-catalog-discovery
  kind: discovery
  purpose: Define safe repository-wide plan discovery across owned documentation trees at arbitrary depth.
  first_cataloged: 2026-08-06T16:22:57Z
  catalog_metadata_updated: 2026-08-06T16:22:57Z
  products:
    - codeheart-operating-kit
  capabilities:
    - semantic-plan-catalog
    - portfolio-coordination
  strategic_themes:
    - ownership-aware-discovery
  relations:
    - kind: related
      target: codeheart-operating-kit.discovery.semantic-plan-catalog-coordination
```
<!-- END CODEHEART PLAN METADATA -->

## Discovery Status

Input state: new request with an agreed semantic direction and a current implementation to audit.

Output target: manual-review-ready. This document records one approved coherent contract and an
unconditional, bounded implementation-planning handoff. It is discovery authority for planning,
not authority to implement, release, activate a consumer repository, or publish Git changes.

This discovery is:

- `routing-bearing`, because the next implementation changes local plan commands, portfolio scan
  routing, migration routes, and ownership-error handling; and
- `recipe-bearing`, because deterministic Git enumeration, Markdown candidate detection, guarded
  migration, and remote overlay comparison are executable mechanics.

A separate reviewer agent was not used because this task did not authorize sub-agent delegation.
The main thread performed and recorded a reviewer pass over the implementation-shaping decisions.

## Problem Framing

### Problem Statement

The plan catalog recognizes formal records only below `docs/repo/plans/`, but the managed
structure-governance contract places plans with the owner: repository, business/domain, product,
source area, package, module, or another repository-owned documentation boundary. The catalog can
therefore report a canonical repository as complete while omitting valid plan authority stored in
other owned documentation trees.

The inconsistency is observable in both doctrine and implementation:

- `components/structure-governance/managed/reference/documentation-structure.md` permits plans
  under repository, product, source-area, package, business, and other owned docs roots;
- `components/planning-workflows/managed/reference/plan-catalog-format.md` restricts formal
  discovery to `docs/repo/plans/`; and
- local discovery, inventory, migration placement, branch-touch detection, remote baseline scans,
  and remote branch overlays all repeat the fixed `docs/repo/plans` boundary.

The intended invariant is broader: a real plan in a repository-owned documentation tree should be
catalog-visible regardless of directory depth, while managed, generated, dependency, vendored,
fixture, submodule, nested-repository, and external content must not become plan authority.

### User Intention

Make the semantic plan catalog complete across repository-owned documentation trees without
turning arbitrary embedded Markdown into consumer-owned authority.

### Resolved User Decisions

The following are treated as approved input, not new recommendations:

- actual plan records inside repository-owned documentation trees are catalog-visible at arbitrary
  depth;
- a documentation match means a normalized repository-relative path containing a directory
  segment exactly named `docs`, rather than a single-level glob such as `*/docs/**`;
- filename or genuine metadata discovers a candidate;
- canonical validity requires a supported filename and valid metadata;
- stable identity comes from metadata;
- human-readable record kind comes from filename and family qualification;
- lifecycle remains the document header `Status:` value; and
- managed, vendored, generated, dependency, submodule, nested-repository, and external content is
  not consumer-owned plan authority;
- authoritative local paths come from Git index regular blobs and authoritative remote paths from
  selected commit-tree regular blobs, with only a non-authoritative opt-in untracked preview;
- automatic tracked-docs discovery uses hard-unowned boundaries and `excluded_roots`, without an
  `owned_roots` registry or mandatory hybrid ownership phase;
- discovery v2 requires explicit reviewed activation in existing repositories, while fresh
  repositories may default to v2;
- plan metadata remains schema v1 while machine contracts whose completeness or provenance changes
  are versioned independently;
- family authority requires exact `README.md` plus genuine valid family metadata and is never
  created by directory shape alone;
- existing cutover revisions and frozen registers retain their historical meaning, and mixed mode
  remains optional rather than a mandatory adoption phase; and
- default-branch configuration controls remote discovery authority, with v2 portfolio completeness
  requiring compatible evidence from every required member.

### Goals

- Define one portable path classifier for local commands, migration, remote baselines, branch
  overlays, inventory protection, and family qualification.
- Define exact discovery/implementation filename rules and metadata-qualified family rules.
- Detect metadata-bearing plans even when the filename is wrong, without treating examples inside
  code blocks as metadata.
- Preserve legacy, mixed, and canonical mode meaning while making direct complete migration from
  legacy to canonical the normal route when no filename-only record is deferred.
- Make ownership exclusions explicit, reviewable, deterministic, and safe across platforms.
- Prevent a Kit upgrade from silently invalidating existing canonical repositories.
- Preserve frozen registers, semantic IDs, lifecycle headers, source hashes, branch ownership, and
  remote completeness evidence through migration.
- Give Codeheart-HQ and other consumers a bounded sequence for adopting the expanded set.

### Non-Goals

- Do not implement CLI, runtime, schema, fixture, test, managed-resource, or doctrine changes in
  this discovery.
- Do not modify a consumer repository or inspect private consumer content.
- Do not upgrade or resume Codeheart-HQ migration or catalog cutover.
- Do not change plan metadata schema fields, semantic ID grammar, lifecycle values, or relation
  kinds merely because paths expand.
- Do not treat execution logs, attachments, general project READMEs, or arbitrary Markdown as plan
  records.
- Do not add a central repository allowlist that becomes the only way future docs roots are found.
- Do not authorize commits, pushes, pull requests, merges, releases, or destructive Git actions.

### Priority Order

1. Ownership and authority safety.
2. Complete plan discovery.
3. Compatibility and reversible adoption.
4. Deterministic cross-platform behavior.
5. Clear diagnostics and migration evidence.
6. Remote-scan security and completeness.
7. Performance and operational simplicity.

### Durable Principles

- Ownership is decided before artifact kind.
- A path convention may discover authority but cannot create ownership by itself.
- Invalid candidates stay visible as findings; they are not silently dropped or repaired.
- Git paths and blobs are inert evidence. Scanning never executes repository content.
- One classifier owns local and remote semantics.
- A compatibility baseline is evidence, not permission to keep materially changing filename-only
  plans.
- An old complete cache is historical evidence, not current completeness under a broader contract.

### Manual-Review Success Criteria

This discovery is manual-review-ready when it records current evidence, the approved contract,
alternatives and tradeoffs, mode and error semantics, compatibility and migration, acceptance
criteria, a test matrix, risks, only genuinely non-blocking release questions, and a bounded
implementation-planning handoff.

## Context And Evidence

### Essential Sources

| Source | Evidence used |
| --- | --- |
| `AGENTS.md` | Producer authority, public-core safety, required change and release routes, and preservation rules. |
| `components/planning-workflows/managed/runbooks/discovery-workflow.md` | Manual-review target, decision ledger, reviewer gate, and handoff boundary. |
| `components/planning-workflows/managed/reference/plan-catalog-format.md` | Current fixed root, metadata, modes, identities, observations, and validation contract. |
| `components/planning-workflows/managed/reference/planning-document-lifecycle.md` | Record kinds, lifecycle header, plan shapes, families, attachments, and register behavior. |
| `components/structure-governance/managed/reference/documentation-structure.md` | Ownership-first placement across repo, business, product, source-area, and package docs. |
| `components/structure-governance/managed/reference/managed-content-boundaries.md` | Managed, consumer-owned, local, generated, and report boundaries. |
| `docs/repo/reference/placement-contract.md` | Producer and consumer ownership surfaces. |
| `docs/repo/reference/consumer-impact-classification.md` | Migration, validator, instruction, and placement impact classes. |
| `components/planning-workflows/managed/reference/portfolio-coordination-format.md` | Git-object scanning, completeness, cache, branch overlay, and strategic-overlay boundaries. |
| `components/planning-workflows/managed/runbooks/migrate-plan-catalog.md` | Inventory, ledger, hash, branch-owner, mixed, canonical, and remote-overlay gates. |
| `internal/plancatalog/*.go` | Local enumeration, parsing, validation, inventory, mixed baseline, view, and migration behavior. |
| `internal/portfolio/scanner.go` and `internal/portfolio/mirror.go` | Remote Git-tree enumeration, branch comparison, observation generation, and cache behavior. |
| `internal/commands/plans.go` | `list`, `validate`, `inventory`, `migrate`, and remote-overlay command behavior. |
| `schemas/plan-metadata.schema.json` | Path-independent metadata schema v1. |
| `schemas/plan-catalog.schema.json` | Portfolio catalog/cache schema v1 and path-bearing observations. |
| `schemas/plan-migration-ledger.schema.json` | Per-record revision/hash decisions without a discovery-scope fingerprint. |
| `schemas/kit-config.schema.json` | Catalog mode and cutover revision, with no discovery-contract or ownership settings. |

### Repository And Worktree Preflight

Evidence was gathered from clean detached `origin/main` commit
`7ef63d4c6b8eb090c17691b2a67c3f6b699ca25b`. The semantic-plan-catalog development, cutover,
handoff, and release-evidence branches are ancestors of this commit. No current
`.codeheart/kit.transaction.json` marker or transaction recovery directory was present. No
overlapping tracked changes existed before this discovery bundle was created.

### Evidence-Backed Current Behavior

| Surface | Current behavior | Consequence |
| --- | --- | --- |
| Local enumeration | `internal/plancatalog/discover.go` calls `filepath.WalkDir` only at `docs/repo/plans`. | Other docs roots are invisible; untracked and ignored worktree files below the fixed root can still be candidates. |
| Discovery filename | A case-sensitive suffix match recognizes `*_discovery_doc.md` and `*_implementation_doc.md`. | Metadata in a differently named Markdown file is never inspected. |
| Family filename | A case-insensitive nested `README.md` qualifies when two immediate child directories contain filename-formal descendants. | Qualification is fixed to one plans root and can traverse fixture or nested-repository content inside it. |
| Metadata markers | Exact begin/end markers outside backtick or tilde fences are parsed; marker examples in fenced code are ignored. | Fenced examples are covered, but discovery never scans a non-filename Markdown file. |
| Legacy mode | Missing metadata is a warning. | Filename-only records can be valid. |
| Mixed mode | Missing metadata is a warning in discovery, then an error unless path and frozen-register evidence exist at the cutover revision. | The proof is path-based and does not detect material post-cutover edits to a grandfathered file. |
| Canonical mode | Missing metadata is an error for enumerated filenames. | A metadata-only plan with the wrong filename is omitted instead of diagnosed. |
| Inventory | Reuses fixed enumeration, HEAD, per-file hashes, dirty checks, and branch touches limited to `docs/repo/plans`. | Coverage and branch-ownership evidence omit other docs roots. |
| Migration | `FormalPathKind` accepts targets only beneath `docs/repo/plans`; projected canonical coverage uses the same set. | Migration cannot repair or adopt repository-wide plan candidates. |
| Inventory output protection | Protects the root README, register, and fixed-root formal paths. | An output artifact can collide with a future out-of-root plan path without current recognition. |
| Remote baseline | The scanner lists regular Git blobs only below `docs/repo/plans` and validates them canonically. | Remote scans safely exclude symlinks/gitlinks but omit owned nested docs elsewhere. |
| Remote branch overlay | Git diff is path-limited to `docs/repo/plans`; identical-byte renames are suppressed and deletions produce no tombstone. | Branch plans under other docs roots and eligibility-changing renames are invisible. |
| Cache schema | Catalog schema v1 has a free canonical path but no discovery-contract version. | A v1 cache can appear complete after discovery semantics expand. |
| Config | Catalog mode and original cutover revision exist, but no path-ownership or discovery-version gate exists. | An immediate semantic change would alter completeness on upgrade with no explicit adoption boundary. |

The source CLI at the inspected commit reported canonical mode, `valid: true`, and 33 records. That
result proves current consistency only for the fixed root. The focused current suites passed:

```text
go test ./internal/plancatalog ./internal/portfolio ./internal/commands
```

They cover fixed-root modes, metadata markers, family qualification, symlink containment, frozen
mixed baselines, migration hashes and transactions, remote branches, cache safety, and command
behavior. They do not cover repository-wide docs roots or ownership classification.

### Contract Inconsistency

The catalog format currently narrows the authority surface below the placement doctrine. This is
not merely a documentation wording issue: the same fixed prefix is embedded independently in
local enumeration, formal-path classification, active-branch diffing, migration coverage,
inventory protection, remote tree listing, and remote changed-path selection. Changing only the
reference or one glob would leave inconsistent authority and security boundaries.

### Candidate Domains And Decision Clusters

| Cluster | Why involved | Future owner |
| --- | --- | --- |
| Plan semantics | Candidate signals, kind, metadata, lifecycle, modes, and families. | Planning workflows |
| Documentation ownership | Owned docs trees, fixtures, managed content, generated output, and placement. | Structure governance plus repo config |
| Git evidence | Tracked paths, modes, hashes, refs, branch changes, path normalization, and nested boundaries. | Plan catalog and portfolio scanner |
| Migration compatibility | Discovery-version gate, inventory, ledger, cutover, frozen register, and HQ adoption. | Planning migration route |
| Portfolio completeness | Default baselines, branch overlays, cache version, source observations, and stale evidence. | Portfolio coordination |
| Security and performance | Inert reads, no symlink traversal, bounded concurrency, streaming, and deterministic output. | CLI runtime and tests |

## Requirements And Constraints

### Functional Requirements

- `FR-01`: Find plan candidates in every repository-owned path whose normalized segments contain
  exact segment `docs`, including root `docs/**` and arbitrarily nested paths.
- `FR-02`: Use one path classifier rather than filesystem glob expansion.
- `FR-03`: Discover candidates by supported filename or structurally genuine plan metadata.
- `FR-04`: Keep metadata-like examples inside fenced or indented code from becoming candidates.
- `FR-05`: Require filename plus valid metadata for a canonical record.
- `FR-06`: Derive discovery or implementation kind from filename. Recognize family authority only
  from exact `README.md` plus genuine valid family metadata, then require metadata kind and ID kind
  to agree; directory shape may validate but never create family authority.
- `FR-07`: Preserve lifecycle exclusively in the document header.
- `FR-08`: Apply exact legacy, mixed, and canonical behavior to the expanded candidate set.
- `FR-09`: Report invalid and ownership-ambiguous candidates with stable problem codes.
- `FR-10`: Detect duplicate semantic IDs across all owned docs roots.
- `FR-11`: Apply the same semantics to list, validate, inventory, migrate, remote baselines,
  changed-path/branch ownership, remote overlays, family qualification, inventory-target
  protection, and portfolio completeness.
- `FR-12`: Preserve old discovery behavior until an existing repository explicitly activates v2.
- `FR-13`: Provide prospective v2 inventory and validation before activation.
- `FR-14`: Preserve frozen registers and the original catalog cutover revision.
- `FR-15`: Bind migration review to the discovery contract, ownership policy, source revision,
  candidate set, and source hashes.
- `FR-16`: Treat incomplete or older-version remote evidence as incomplete, never current.
- `FR-17`: Provide an authoring preview for untracked, non-ignored candidates without granting
  them canonical authority.

### Non-Functional Requirements

- `NFR-01 Determinism`: Normalize Git paths to `/`, preserve Git spelling, sort bytewise, and do
  not depend on host glob, filesystem traversal order, or case-folding behavior.
- `NFR-02 Cross-platform`: Require exact lowercase `docs`, exact lowercase suffixes, and exact
  `README.md`; detect case/Unicode portability collisions rather than choosing by host behavior.
- `NFR-03 Security`: Read Git blobs and bound regular files as inert bytes; do not checkout remote
  content, follow symlinks, enter gitlinks/nested repositories, execute filters/hooks, or run
  repository code.
- `NFR-04 Ownership safety`: Hard-unowned boundaries cannot be overridden. Candidate signals below
  conventional fixture, vendor, generated, dependency, or build-like boundaries remain visible
  prospective blockers until an explicit `excluded_roots` decision or truthful plan relocation.
- `NFR-05 Performance`: Enumerate the Git index/tree once per ref, stream Markdown detection and
  hashes, avoid a Git subprocess per file, retain the portfolio scanner's maximum concurrency of
  four repositories, and keep output ordering deterministic.
- `NFR-06 Memory`: Do not require full repository trees or all Markdown blobs in memory at once.
- `NFR-07 Compatibility`: Upgrade alone produces no new canonical errors in an existing repo.
- `NFR-08 Auditability`: Inventory records the included, excluded, ambiguous, invalid, untracked-
  preview, and remote candidate sets with reasons.
- `NFR-09 Public-core safety`: Config and artifacts contain repository-relative paths and hashes,
  never credentials, raw private logs, or unnecessary local paths.
- `NFR-10 Failure honesty`: Any missing tree, truncated stream, unsafe path, scope mismatch, or
  ownership ambiguity makes the affected validation or remote scan incomplete.

## Options Considered

### Option A - Replace The Fixed Root With Globs

Examples: `docs/**`, `*/docs/**`, or `**/docs/**`.

Advantages:

- smallest apparent code change;
- easy to describe for one filesystem.

Rejected because:

- `*/docs/**` is not arbitrary-depth;
- `**` behavior varies by glob engine and pathspec dialect;
- independent local, Git, and remote glob implementations can disagree;
- globs do not establish ownership; and
- they do not solve metadata-only candidates, fixtures, submodules, case, cache versioning, or
  migration compatibility.

### Option B - Explicit Allowlist Of Documentation Roots

Advantages:

- strong explicit ownership;
- simple exclusion of embedded content.

Rejected as the primary model because every future business, venture, product, source-area, or
package docs root can be silently omitted until configuration is updated. It reverses the agreed
default that real owned docs should be catalog-visible.

### Option C - Required Root Marker In Every Docs Tree

Advantages:

- ownership is local and visible;
- nested repositories can declare themselves separately.

Rejected as a required model because it adds marker churn to every docs root, silently omits new
roots, complicates remote branch behavior, and duplicates Git/config authority. A marker may be
considered later as an optional ownership signal, not the v2 prerequisite.

### Option D - Automatic Git Discovery With Hard Boundaries And Exclusions

Advantages:

- future docs roots are found automatically;
- ordinary tracked docs paths require no registration;
- hard-unowned content remains outside authority; and
- explicit exclusions are reviewable inventory evidence.

Recommendation: choose Option D with an explicit discovery-contract v2 activation gate. Tracked
fixtures, examples, vendored content, generated docs, dependencies, and build output with plan
signals are prospective blockers until explicitly excluded. A genuine plan below such a misleading
boundary should normally move to a truthful owned docs path instead of gaining an inclusion
override.

### Option E - Hybrid Automatic Discovery With Owned And Excluded Roots

Advantages:

- future owned docs roots are included automatically;
- unusual conventional locations can be force-classified as owned.

Rejected for the first implementation because no concrete genuine-plan case requires an
`owned_roots` override, it complicates authority and overlap rules, and it encourages keeping plans
under misleading fixture/vendor/generated-style boundaries. Add no `owned_roots` capability
without a demonstrated requirement and separate review.

## Recommended Contract

### 1. Configuration And Activation

Add these planning-workflow settings:

```yaml
component_settings:
  planning-workflows:
    plan_catalog_mode: legacy | mixed | canonical
    plan_catalog_discovery_version: 1 | 2
    plan_catalog_ownership:
      excluded_roots:
        - <repo-relative-directory>/
```

The field names above are the proposed public shape; the implementation plan must verify exact
schema placement and upgrade mechanics without changing the approved semantics.

Semantics:

- absent `plan_catalog_discovery_version` means v1 for existing repositories;
- fresh repositories created by the new release write v2;
- `excluded_roots` explicitly acknowledges a repository-relative external, fixture, vendored, or
  generated, dependency, or build-like boundary;
- entries are slash-normalized relative directories with no glob syntax, absolute paths, empty
  segments, `.`/`..`, symlink traversal, or Git-link traversal;
- duplicate, case-fold/Unicode-equivalent, or ancestor/descendant overlaps between exclusions are
  invalid rather than silently normalized; and
- hard-unowned boundaries cannot be reclassified by configuration.

There is no `owned_roots` setting in v2. Every ordinary tracked Markdown path containing an exact
`docs` segment is owned by default unless it is hard-unowned, conventionally ambiguous pending
review, or explicitly excluded.

The original `plan_catalog_cutover_revision` retains its current meaning. Do not repoint it or use
it as the v2 activation revision.

### 2. Path Universe

Authoritative local candidate paths come from the Git index, using NUL-delimited stage records and
Git modes. Remote candidate paths come from the selected commit tree. The classifier does not walk
the repository filesystem to decide authority.

For each Git path:

1. require a regular blob mode (`100644` or `100755`);
2. preserve the Git path spelling and normalize separators to `/`;
3. reject absolute, empty, invalidly encoded, or dot-segment paths;
4. require an exact lowercase `.md` extension for canonical Markdown discovery, while an anomaly
   pass may inspect case-variant Markdown extensions only to report non-portable metadata claims;
5. classify ownership; and
6. treat the path as an eligible documentation path when at least one directory segment is exactly
   lowercase `docs`.

This includes, for example:

```text
docs/repo/plans/feature/feature_discovery_doc.md
docs/business/procurement/supplier-onboarding_implementation_doc.md
docs/business/ventures/venture-a/discovery/market-fit_discovery_doc.md
products/widget/docs/plans/widget-roadmap_implementation_doc.md
products/widget/api/docs/architecture/api-shape_discovery_doc.md
products/widget/packages/client/docs/client-migration_implementation_doc.md
```

It does not rely on a `plans` segment. A `plans/` directory remains preferred placement where it
fits, but procurement, business, and other owned domain trees are not excluded merely because they
use a different type folder.

When multiple `docs` segments occur, eligibility remains true and the path is enumerated once.
No deepest-anchor or fixed-depth calculation changes authority.

### 3. Ownership Classification

Classify each potential path as `owned`, `excluded`, `prospective-blocked`, or `hard-unowned`
before cataloging it.

Hard-unowned boundaries:

- `.codeheart/kit/`, `.codeheart/local/`, and `.codeheart/user/`;
- Git symlinks, gitlinks/submodules, and paths outside the superproject index/tree;
- nested Git repositories not represented as superproject regular blobs;
- ignored and untracked files as authoritative catalog inputs; and
- an explicit path that escapes the repository or traverses a non-directory/symlink boundary.

Conventional ambiguity signals include directory segments such as:

```text
vendor
third_party
third-party
external
deps
dependencies
node_modules
fixtures
testdata
examples
generated
.generated
build
dist
out
target
.venv
venv
```

These names are conventional ambiguity signals, not silent exclusions. A plan signal below one is
retained as a prospective blocker and produces `plan_documentation_root_ambiguous` until the
repository lists the relevant boundary in `excluded_roots` or moves the genuine plan to a truthful
repository-owned docs location. This makes the Operating Kit producer's plan-shaped test fixtures
visible during prospective v2 adoption without treating them as real plans.

An explicit excluded root is not silently discarded: inventory retains the candidate path, signal,
and exclusion reason as non-authoritative evidence. Validation does not fail merely because an
intentionally excluded fixture demonstrates plan syntax.

There is no inclusion override. A genuine plan beneath a fixture, vendor, generated, dependency,
or build-like boundary is misplaced by default; the normal resolution is relocation. A future
`owned_roots` mechanism requires concrete evidence, a separate discovery/review, and explicit
contract revision.

Tracked external content in an unusual, unmarked path cannot always be inferred automatically.
Repository owners remain responsible for declaring such a root excluded. This residual limitation
is preferable to a permanent allowlist that silently omits all future owned roots.

### 4. Candidate Signals And Markdown Structure

For owned, eligible Markdown, a discovery or implementation candidate exists when either signal is
present:

- `F`: a supported discovery or implementation filename; or
- `M`: exact Codeheart plan metadata markers in structural Markdown, outside fenced and indented
  code.

Family candidate handling is narrower: an exact `README.md` becomes v2 family authority only when
it also contains genuine valid metadata declaring kind `family`. A metadata-bearing non-README
family claim remains a candidate so the filename error is visible. An ordinary `README.md` with no
family metadata is an index/router, not a plan candidate merely because child directories exist.

Use one shared Markdown structure scanner for candidate detection and metadata parsing. It must:

- recognize backtick and tilde fences of length three or greater using CommonMark indentation and
  close rules;
- ignore indented code blocks;
- require exact marker comment lines, not substring matches;
- retain unmatched, reversed, duplicate, or misplaced real markers as invalid candidates;
- treat the canonical bounded block immediately below the canonical H1 as the only valid metadata
  position; and
- never interpret marker-like content inside examples as plan metadata.

Scanning all tracked Markdown is necessary to diagnose a metadata-bearing plan outside an eligible
docs path. The anomaly pass also checks case-variant Markdown extensions for exact markers so
`plan_path_case_mismatch` is observable. It is a read-only pass; content outside owned docs never
becomes authority.

### 5. Exact Filename And Kind Rules

Discovery filename:

```text
<non-empty-name>_discovery_doc.md
```

Implementation filename:

```text
<non-empty-name>_implementation_doc.md
```

The suffix and `.md` are exact and case-sensitive. This discovery intentionally does not add a new
stem-slug grammar; existing valid names keep working, and stable identity does not come from the
stem. A later naming-lint proposal may narrow stems separately.

Family filename:

```text
README.md
```

V2 family qualification requires all of the following:

- the exact case-sensitive filename `README.md`;
- an owned eligible docs path;
- genuine, structurally positioned, schema-valid plan metadata declaring kind `family`; and
- a semantic ID whose kind segment is `family`.

Directory shape never independently creates family authority. Structural validation may require a
declared family to meaningfully group child plan bundles inside the same owned boundary, but the
later implementation plan must define the precise grouping invariant and diagnostics. Whatever
invariant it selects cannot turn root `docs/README.md`, `<owner>/docs/README.md`,
`plans/README.md`, or another ordinary domain index/router into a candidate without valid family
metadata.

Already recognized v1 family records keep their v1 compatibility before activation. Prospective
v2 inventory must identify each one and require valid family metadata before v2 activation unless
the repository intentionally remains on v1. Directory shape and a frozen register cannot
grandfather a metadata-free family into v2 authority.

For discovery and implementation, filename determines expected kind. For a qualified family,
exact `README.md` and valid metadata jointly establish family kind. Metadata `kind` and the middle
segment of metadata `id` must both match the expected kind.

### 6. Mode Semantics

| Signals and evidence | Legacy | Mixed | Canonical |
| --- | --- | --- | --- |
| `F`, no `M` | Valid legacy record when lifecycle/header is valid; `metadata_missing` warning. | Valid only when frozen register references the exact path, the cutover revision contains a regular blob, and current bytes are unchanged from that blob beyond line-ending conversion; otherwise `mixed_new_or_changed_plan_metadata_missing` error. | Invalid: `metadata_missing` error. |
| `F` plus valid `M` | Valid canonical record. | Valid canonical record. | Valid canonical record. |
| `F` plus malformed `M` | Invalid; never fall back to legacy. | Invalid; never grandfather malformed metadata. | Invalid. |
| Valid/malformed `M`, no `F` | Invalid: `canonical_filename_missing`; retain metadata diagnostics. | Invalid: `canonical_filename_missing`; not grandfatherable. | Invalid: `canonical_filename_missing`. |
| Neither `F` nor `M` | Not a candidate. | Not a candidate. | Not a candidate. |

`materially changed` in mixed mode is content evidence, not a subjective authoring label. Compare
the current Git blob or bound worktree bytes with the exact cutover blob, accepting only CRLF/LF
checkout conversion. A rename is a new path and requires metadata. A legacy register row alone does
not grandfather later content changes.

Mode controls validity, not discovery. All modes enumerate the same v2 candidate universe so
invalid records cannot disappear by changing mode. The table's filename-only legacy allowance
applies to discovery and implementation records. V2 family authority always requires exact
`README.md` plus valid family metadata; only already-recognized v1 families retain v1 behavior
while the repository remains on discovery v1.

### 7. Stable Validation And Error Behavior

Required problem codes:

| Code | Condition | Severity/behavior |
| --- | --- | --- |
| `metadata_missing` | Filename candidate lacks metadata. | Warning in legacy, conditional in mixed, error in canonical. |
| `mixed_new_or_changed_plan_metadata_missing` | Mixed filename-only candidate is new, renamed, not register-proven, or materially changed after cutover. | Error. |
| `canonical_filename_missing` | Genuine metadata candidate lacks a supported discovery/implementation filename or exact family `README.md`. | Error in every mode. |
| `record_kind_path_mismatch` | Metadata kind or ID kind disagrees with filename-derived kind. | Error. |
| `duplicate_plan_id` | Two valid or partially valid owned candidates use one semantic ID. | Error; preserve every observation. |
| `plan_metadata_misplaced` | Metadata signal is in a repository-owned tracked Markdown path with no exact `docs` segment. | Error; never a valid record. |
| `plan_documentation_root_ambiguous` | A candidate intersects a conventional or conflicting ownership boundary without an explicit resolution. | Error; inventory retains evidence. |
| `plan_discovery_root_overlap` | Excluded roots overlap, collide by portable path comparison, or cause more than one exclusion resolution. | Config error before scan. |
| `plan_path_unowned` | A candidate crosses a hard-unowned managed, symlink, gitlink, nested-repository, ignored/untracked-authority, or escaping boundary. | Error for the candidate; do not read external target bytes. |
| `plan_candidate_excluded` | A candidate signal exists in an explicitly excluded regular Git blob. | Information or warning in inventory; not catalog authority. |
| `plan_path_case_mismatch` | A metadata candidate relies on `Docs`, `.MD`, `readme.md`, or case-varied reserved suffix. | Error with canonical spelling. |
| `plan_path_portability_collision` | Candidate paths collide after Unicode normalization and case-fold comparison. | Error on all colliding paths. |
| `family_placement_invalid` | Metadata claims family but exact `README.md`, valid family metadata, or the selected structural grouping invariant fails. | Error. |
| `plan_source_unsafe` | Candidate Git mode or bound local path is non-regular or changes identity while read. | Error. |
| `remote_discovery_version_mismatch` | A v2 portfolio scan cannot prove a member's v2 discovery contract. | Incomplete scan; preserve prior cache. |

Existing metadata/header schema errors remain stable. Implementations should avoid redundant
cascades: `canonical_filename_missing` may coexist with metadata syntax/schema findings, while a
single ownership classification should not also emit several synonymous root errors.

### 8. Local Git State

- Tracked, unstaged modifications: validate the bound worktree bytes and report dirty provenance.
- Staged additions: part of the authoritative local candidate universe because they are in the
  index; inventory/migration still requires a coherent committed source revision before writes.
- Untracked, non-ignored files: excluded from authority by default. Add
  `plans validate --include-untracked` and the matching inventory preview so authors can validate
  them as `untracked-preview`; never include them in canonical completeness, migration writes,
  remote facts, or frozen baselines.
- Ignored, untracked files: do not open or preview by default. A regular blob already present in
  the Git index remains tracked even if a later ignore rule matches it; classify it by ownership
  rather than silently dropping it.
- Symlinks: never follow. A filename-signaled symlink reports `plan_source_unsafe`.
- Submodules/gitlinks: never descend. Any attempted plan claim beneath one reports
  `plan_path_unowned`.
- Nested Git repositories: do not descend from the superproject. Their plans belong to their own
  repository catalog.

### 9. Local Commands And Migration

`plans list` and `plans validate` must use the same candidate set and ownership classifier.
Structured output should include discovery version and candidate validity/provenance. A nonzero
result remains mandatory when errors exist.

`plans inventory` v2 must record:

- discovery contract version;
- ownership config and a canonical policy digest;
- source revision;
- included, excluded, ambiguous, misplaced, unsafe, and optional untracked-preview candidates;
- candidate-set digest over sorted path, signal, expected kind, Git mode/blob, and ownership
  classification;
- per-record source hash, mode compatibility, branch touches, dirty state, and problem codes; and
- remote discovery version/completeness when requested.

`plans migrate` must accept only a reviewed v2 ledger bound to the same source revision, policy
digest, and candidate-set digest. Enhance it to support:

- metadata insertion for filename-only candidates;
- reviewed rename plus metadata insertion for metadata-only candidates with a wrong filename;
- exact target-path and expected-byte preconditions;
- the existing branch-owner, dirty-overlap, transaction, rollback, and recovery safeguards; and
- projected completeness over the entire v2 candidate set.

Migration must not invent a filename, kind, ID, purpose, ownership decision, family, or target
directory. A prospective-blocked or unowned source remains a blocker until it is explicitly
excluded as non-authority or a genuine plan is moved through a separately reviewed placement
decision. There is no owned-root override.

Inventory-output protection must use the same classifier and protect every present or prospective
formal v2 path, not a second hard-coded root list.

### 10. Remote Baselines And Branch Overlays

The remote scanner should enumerate one commit tree, filter regular `.md` blobs in code, and stream
candidate detection. It must not use `*/docs/**`, `**/docs/**`, or a separate filename classifier.

For branch overlays:

- obtain the repository-wide NUL-delimited changed-path set once;
- classify old and new paths with the same v2 classifier;
- include added or materially changed owned plan candidates;
- suppress an identical-byte rename only when old and new paths have the same ownership,
  eligibility, and filename-derived kind;
- treat a rename into an owned docs path or across kind/ownership boundaries as semantic change;
- retain the existing no-tombstone rule for deletions in the first v2 version; and
- never infer new family authority from directory shape or an unchanged `README.md` merely because
  a second sibling directory appears; family authority comes from exact filename plus valid family
  metadata, while structural checks only validate an already-declared family.

Default-branch membership and discovery configuration controls remote scanning. An unmerged branch
cannot remove exclusions, broaden authority, or activate v2. Pull-request facts remain enrichment
only.

The scanner continues to read inert Git objects, use disabled hooks and restricted transports,
redact local source locators from public migration artifacts, and omit worktrees, local heads, and
unpushed commits.

### 11. Performance And Determinism

The implementation should use one Git index/tree listing per revision and a batch blob reader
rather than spawning `git show` per Markdown file. Candidate detection and hashing should stream.
Repository scans remain bounded to four concurrent repositories; within one repository, bounded
workers may parse blobs but must emit bytewise path-sorted results.

Performance validation should prove linear enumeration over a synthetic repository with at least
100,000 tracked paths and 10,000 Markdown files, bounded memory independent of total blob bytes,
and a constant number of Git processes per ref. Wall-clock thresholds should be recorded from CI
baselines rather than frozen to one developer machine.

## Compatibility And Migration Strategy

### Versioning Recommendation

- Keep plan metadata schema at version 1; identity and fields do not change.
- Add `plan_catalog_discovery_version` to config without changing the global config schema version.
- Bump local view/inventory structured output to schema version 2 because completeness and
  provenance semantics change.
- Bump migration ledger to schema version 2 to require discovery/policy/candidate-set binding and
  optional reviewed target paths.
- Bump portfolio plan-catalog/cache schema to version 2 and add the discovery-contract version to
  completeness evidence.
- Do not treat a v1 portfolio cache as complete under v2. Preserve it as historical evidence until
  a complete v2 scan atomically replaces it.

### Compatibility Window And Upgrade Gate

The release that introduces v2 should keep absent/explicit v1 behavior for existing repositories
and report a non-error `plan_discovery_upgrade_available` notice. Fresh repositories use v2.

Provide prospective reads:

```sh
codeheart-operating-kit plans inventory \
  --target-discovery-version 2 \
  --remote-overlays \
  --output <inventory.json> \
  .

codeheart-operating-kit plans validate \
  --target-discovery-version 2 \
  --remote-overlays \
  .
```

Only an explicit reviewed config change activates v2. Keep v1 for at least one released
compatibility cycle. The exact release in which v1 support may be removed remains a manual product
decision and must be announced in release and migration notes.

### Existing Catalog Modes

`legacy`:

- v1 remains unchanged until activation;
- prospective v2 reports all new candidates and ownership findings;
- when every formal candidate can receive valid metadata in one reviewed migration, the normal
  route is `legacy -> migrate all candidates -> canonical` with v2 activation in the reviewed
  cutover change;
- remaining in legacy v2 may retain valid discovery/implementation filename-only records, but it
  is not the recommended HQ destination; and
- metadata-only, malformed-metadata, and metadata-free v2 family candidates remain errors.

`mixed`:

- preserve the original cutover revision and frozen register bytes;
- do not grandfather newly discovered v2 paths merely because their content predates the Kit
  upgrade;
- migrate new v2 candidates to both filename and metadata before activating v2; and
- continue grandfathering only unchanged, register-proven original cutover records.

Mixed mode remains supported for large or deferred migrations but is not a mandatory intermediate
phase. Frozen-register and cutover proofs are required only when filename-only discovery or
implementation records are intentionally grandfathered through mixed mode.

`canonical`:

- remain on v1 while prospective inventory and reviewed migration resolve every new v2 candidate;
- activate v2 only when local and required remote validation passes with both signals for every
  owned record; and
- never lower mode to hide expanded coverage gaps.

This gate avoids a second cutover revision. The activation commit, v2 inventory/ledger, policy
digest, and candidate-set digest provide the expansion evidence while the original cutover remains
historically accurate.

### Frozen Registers And Source Hashes

- Keep `docs/repo/plans/plan-register.md` byte-identical in mixed/canonical repositories.
- Do not append newly discovered business, product, procurement, or nested venture paths.
- Reconcile old aliases only into metadata.
- Preserve `Created` and `Last updated` during metadata-only adoption.
- A reviewed rename may change path and Git history evidence but not semantic ID.
- Reinventory after every skipped, branch-owned, dirty, renamed, or changed source.
- A ledger is stale when revision, ownership policy, candidate-set digest, or any expected source
  hash changes.

### Portfolio And Remote Cutover

A v2-complete home scan requires every required member to expose compatible v2 default-branch
configuration and canonically valid plan evidence. A v1 member produces
`remote_discovery_version_mismatch`; the current attempt is incomplete and the previous cache is
preserved.

Branch overlays are rescanned with v2 after each member's default branch activates v2. Same-ID
default/branch observations continue to coexist with repository/ref/commit/path/hash evidence.
The strategic overlay remains repo-owned and unchanged by scans.

### Consumer Migration Sequence

Codeheart-HQ is currently legacy and intends to catalog every genuine plan candidate under any
nested docs tree. Its portfolio-v2 semantic migration remains paused until an approved Kit release
implements this contract. For HQ and other repositories:

1. Wait for an approved Operating Kit release containing discovery v2; upgrade through the normal
   consumer lifecycle route without changing catalog discovery version.
2. Verify clean/understood target paths, no transaction marker, current legacy catalog mode,
   historical register integrity, active branches, and portfolio role. Installing the release must
   not activate v2 automatically.
3. Run prospective v2 local inventory; for a coordination home such as HQ, include remote overlays
   and retain incomplete/version-mismatched member evidence rather than claiming completeness.
4. Review every newly included, excluded, prospective-blocked, metadata-only, filename-only,
   duplicate, and family candidate. Add explicit exclusions for fixture/vendor/generated/
   dependency/build-like non-authority; move any genuine plan found under a misleading conventional
   boundary to a truthful repository-owned docs location.
5. Build a reviewed v2 semantic migration ledger bound to policy, candidate set, source revision,
   and hashes. For Codeheart-HQ, explicitly review the known
   `docs/business/plans/`, `docs/business/procurement/`, and nested venture discovery paths supplied
   by the repository owner.
6. Let active branch owners migrate their current plan bytes. Do not overwrite or freeze unrelated
   work.
7. Dry-run and apply only the reviewed metadata/rename actions; preserve lifecycle headers,
   semantic IDs, frozen register bytes, and unrelated files.
8. Reinventory and require zero deferred filename-only formal records, zero material skips, and
   complete local evidence for the direct route.
9. Activate discovery v2 and move the repository directly from legacy to canonical in the same
   reviewed migration sequence when every formal record has valid filename and metadata. Use mixed
   only when a large or deferred migration intentionally preserves proven filename-only records.
10. For portfolio members, land v2 on default branches through their separately authorized
    workflows and rerun remote scans. A feature branch cannot broaden authority.
11. Activate v2 in the HQ coordination home only after every required member and branch overlay has
    compatible v2 evidence. Preserve and label the last complete cache while the current scan is
    incomplete.
12. Resume HQ's paused portfolio-v2 semantic migration from its existing approved consumer plan
    only after the Kit release is available. The old plan register may remain historical evidence;
    HQ needs neither mixed mode, a second cutover revision, nor grandfathered gaps when every plan
    receives metadata.

Upgrade, migration apply, commit, push, and HQ cutover each require their own repository authority.

## Decision Ledger

| ID | Decision | Class | State | Depends on | Affects | Closure criteria |
| --- | --- | --- | --- | --- | --- | --- |
| `D-01` | Exact `docs` segment at arbitrary depth defines documentation eligibility. | implementation-shaping | approved user direction | none | `FR-01`, `NFR-01` | Preserved exactly in planning. |
| `D-02` | Filename or genuine metadata discovers candidates; canonical validity requires both. | implementation-shaping | approved user direction | `D-01` | `FR-03` to `FR-06` | Preserve signal/validity separation. |
| `D-03` | Identity is metadata, kind is filename/qualification, lifecycle is header. | implementation-shaping | approved user direction | `D-02` | `FR-05` to `FR-07` | No competing authority field. |
| `D-04` | Use the specified legacy/mixed/canonical mode behavior. | implementation-shaping | approved user direction | `D-02`, `D-03` | `FR-08` | Mixed material-change proof is implemented. |
| `D-05` | Use Git index regular blobs locally and selected commit-tree regular blobs remotely as the authoritative universe; retain only an opt-in non-authoritative untracked preview. | implementation-shaping | approved user decision | `D-01` | `FR-11`, `FR-17`, `NFR-01`, `NFR-03` | Preserve authority/preview separation in planning. |
| `D-06` | Use automatic tracked-docs discovery plus hard-unowned boundaries and `excluded_roots`; omit `owned_roots` and any mandatory hybrid root-registry phase. | implementation-shaping | approved user decision | `D-05` | `FR-01`, `FR-09`, `NFR-04` | Preserve owned-by-default ordinary docs and visible prospective blockers. |
| `D-07` | Require explicit discovery-v2 activation for existing repositories; upgrade alone leaves them on v1, while fresh repositories may default to v2. | blocking | approved user decision | `D-05`, `D-06` | `FR-12`, `FR-13`, `NFR-07` | Preserve prospective review and explicit activation. |
| `D-08` | Keep plan metadata schema v1; version view/inventory, ledger, portfolio cache/completeness, and other machine contracts wherever provenance/completeness meaning changes. | implementation-shaping | approved user decision | `D-07` | `FR-15`, `FR-16` | Never present v1 evidence as v2-complete. |
| `D-09` | Require exact `README.md` plus genuine valid family metadata; directory shape only validates and never creates family authority, with explicit v1-family migration. | implementation-shaping | approved user decision | `D-01`, `D-02`, `D-06` | `FR-06`, `FR-11` | Preserve metadata-qualified authority and overlay stability. |
| `D-10` | Preserve existing cutover revisions/frozen registers, bind prospective v2 migration evidence, make direct legacy-to-canonical migration normal when complete, and keep mixed optional for grandfathered deferrals. | implementation-shaping | approved user decision | `D-07`, `D-08` | `FR-14`, `FR-15` | No second cutover merely for v2; proofs only for mixed grandfathering. |
| `D-11` | Let default-branch config control remote authority and require compatible v2 evidence from every required member for v2 completeness, preserving the last complete cache otherwise. | implementation-shaping | approved user decision | `D-07`, `D-08` | `FR-16`, `NFR-10` | Preserve inert evidence and honest incomplete scans. |

### Reviewer Pass

Reviewer mode: documented main-thread architecture and ownership pass.

`D-05` and `D-06`:

- Challenge: Git tracking does not prove semantic ownership, while an allowlist would omit future
  docs roots and a filesystem walk includes ignored/nested content.
- Response: keep Git index/tree regular blobs as authority, apply non-overridable hard boundaries,
  and make ordinary exact-`docs` paths owned by default. Conventional embedded roots with plan
  signals become visible prospective blockers until excluded or truthful placement is restored.
  Add only `excluded_roots`; retain untracked preview strictly outside authority.
- Result: approved exclusions-only model; `owned_roots` and a hybrid root-registry phase are out of
  first-implementation scope.

`D-07`, `D-08`, and `D-10`:

- Challenge: immediate expansion breaks canonical consumers; a second cutover revision could
  corrupt the historical meaning of the first; an old cache could falsely claim completeness.
- Response: add a separate discovery-version activation gate, prospective reads, schema-v2
  completeness evidence, and policy/candidate-set hashes. Preserve the original cutover and frozen
  register. Allow a complete legacy repository to migrate every candidate and enter canonical
  directly; use mixed only for intentional grandfathering.
- Result: approved; no automatic config rewrite, no mandatory mixed phase, and no second cutover
  revision.

`D-09`:

- Challenge: arbitrary nested README qualification can classify collection routers or fixtures as
  families, and branch creation of a second sibling can change qualification without README bytes.
- Response: exact `README.md` plus genuine valid family metadata jointly establish v2 family
  authority. Directory shape may validate a declared family but cannot create one, so unchanged
  README bytes never gain new authority merely because a sibling appears. Prospective inventory
  explicitly migrates already-recognized v1 families.
- Result: approved metadata-qualified family model; the implementation plan must define only the
  non-authoritative structural grouping check.

`D-11`:

- Challenge: allowing branch config to broaden discovery would let an unmerged ref change scan
  authority; accepting v1 members would make v2 completeness incomparable.
- Response: default-branch config controls ownership/version and incompatible required members make
  the scan incomplete; preserve and label the last complete cache rather than replacing it with a
  false-completeness result.
- Result: approved.

Residual reviewer uncertainty is non-blocking: the conventional ambiguity-segment list should be
refined only from evidence and the exact v1 removal release remains a later release decision.
Neither prevents unconditional implementation planning from this approved architecture.

## Open Questions And Assumptions

### Open Questions

`OQ-01` - In which later release may discovery v1 support be removed?

- Owner: release owner.
- `BLOCKER: no` for implementation planning; blocker only for the later removal release.
- Affects: release notes and deprecation tests.
- Recommended answer: at least one released compatibility cycle, with removal separately approved.

`OQ-02` - Which conventional ambiguity segments should the first implementation include based on
repository evidence?

- Owner: implementation plan reviewer.
- `BLOCKER: no`; the approved behavior is stable even if the initial evidence-backed list changes.
- Affects: ownership config schema and tests.
- Recommended answer: begin with the documented list, validate it against producer/consumer
  fixtures, and require explicit evidence and tests for additions or removals.

### Assumptions

- `A-01`: The Codeheart-HQ paths named in the request are repository-owned, but their individual
  plan semantics still require HQ-local review.
- `A-02`: Existing metadata schema v1 remains sufficient because path is observation evidence, not
  metadata identity.
- `A-03`: A complete v2 portfolio may require coordinated member adoption; temporary incomplete
  scans are acceptable when disclosed.
- `A-04`: Git is available wherever current plan inventory and portfolio scanning are supported.
- `A-05`: Exact lowercase path conventions are preferable to host-dependent case acceptance.

## Risks And Mitigations

| ID | Risk | Impact | Mitigation/detection |
| --- | --- | --- | --- |
| `R-01` | Tracked fixtures or vendored docs expose plan signals. | Activation blocker or accidental false authority. | Visible prospective blocker, explicit exclusion, inventory evidence, and producer fixture tests; no owned override. |
| `R-02` | A new owned docs root is omitted by config. | Incomplete catalog. | Automatic exact-segment discovery; config is exceptions, not allowlist. |
| `R-03` | Metadata examples trigger candidates. | False errors. | Shared CommonMark-aware fence/indented-code scanner and positive/negative fixtures. |
| `R-04` | Existing canonical repos fail immediately after upgrade. | Adoption outage. | Default v1 for existing config; prospective v2 gate. |
| `R-05` | Mixed grandfathering permits post-cutover edits without metadata. | New authority remains legacy. | Compare current bytes to exact cutover blob beyond line endings. |
| `R-06` | Old cache is presented as complete under v2. | False portfolio conclusions. | Catalog schema v2 and remote version mismatch blocker. |
| `R-07` | Remote scan executes or follows embedded content. | Security boundary breach. | Git blobs only, safe transport policy, no checkout/hooks/filters/symlinks/gitlinks. |
| `R-08` | Repository-wide Markdown scanning is slow. | Poor command/scan operability. | One tree listing, batch streaming, bounded concurrency, synthetic benchmark. |
| `R-09` | Case variants work on Linux but fail on Windows/macOS. | Non-portable catalogs. | Exact reserved casing and portability-collision errors. |
| `R-10` | Migration renames overwrite branch-owned or dirty paths. | Data loss. | Ledger target hashes, active branch ownership, no-replace transaction, rollback/recovery. |
| `R-11` | Directory shape creates family authority without branch-visible bytes. | Missing/stale overlay or ordinary router cataloged as a plan. | Require exact README plus valid family metadata; structural checks never create authority. |
| `R-12` | Unusual unmarked external content is treated as owned. | False authority. | Document residual limit, retain candidate evidence, add explicit exclusions, and review prospective inventory before activation. |

## Acceptance Criteria

- One shared classifier produces the same candidate paths and kinds for local worktree/index,
  committed tree, migration, inventory protection, remote baseline, and remote overlay tests.
- Root `docs/**` and paths with `docs` at depth greater than one are discovered without glob
  semantics.
- Filename-only, metadata-only, malformed, mismatch, duplicate, family, ambiguous, excluded,
  unowned, case, symlink, gitlink, nested-repository, ignored, and untracked-preview cases produce
  the specified results.
- Ordinary tracked Markdown beneath any exact lowercase `docs` segment is owned by default;
  `excluded_roots` and hard-unowned boundaries are honored, conventional plan signals block
  activation until resolved, and no `owned_roots` setting exists.
- Legacy, mixed, and canonical behavior matches the mode table, including byte-based mixed
  material-change detection.
- A fully migrated legacy fixture activates v2 and moves directly to canonical without mixed mode;
  a deferred fixture proves that mixed remains optional and requires frozen-register/cutover proof.
- V2 family authority requires exact `README.md` plus valid family metadata; directory changes
  alone never create a family candidate or new branch observation, and recognized v1 families have
  explicit prospective migration results.
- Existing repositories remain v1 after upgrade; fresh repositories use v2; prospective commands
  do not change config.
- V2 migration is source-, policy-, candidate-set-, and target-hash-bound; dry-run/apply is
  idempotent and preserves `Last updated` for metadata-only adoption.
- Frozen register bytes and the original cutover revision remain unchanged.
- Remote baselines and branch overlays find arbitrary-depth docs paths, reject incompatible member
  versions, and preserve the previous complete cache on any incomplete attempt.
- Cache/view/inventory/ledger schemas and command goldens expose v2 semantics explicitly.
- Source and packaged managed resources remain byte-identical after the future implementation.
- Cross-platform tests pass on macOS, Linux, and real Windows; case/path behavior is identical.
- The synthetic enumeration benchmark proves one listing/batch-read design and bounded memory.
- Release and migration notes classify the change as at least `instruction-only change`,
  `validator-only change`, and `consumer migration required`; any ownership-placement effect is
  reviewed against the breaking placement class.

## Test Matrix

| Area | Positive cases | Negative/boundary cases |
| --- | --- | --- |
| Path depth | `docs/business/...`; `products/a/docs/...`; deeply nested `a/b/docs/...` | `mydocs/`, `Docs/`, `.MD`, no docs segment |
| Filename | discovery and implementation exact suffixes | empty stem, case variants, metadata-only ordinary filename, kind mismatch |
| Metadata detection | exact block after H1 | 3/4-backtick, tilde, indented code, blockquote example, unmatched/reversed/duplicate/misplaced markers |
| Modes | legacy filename-only; direct legacy-to-canonical complete migration; optional mixed unchanged baseline; canonical both | mixed new/renamed/changed/unproven; canonical missing metadata; malformed never falls back |
| Identity | unique IDs across roots | duplicate IDs across repo/business/product roots and default/branch refs |
| Families | exact `README.md` with valid family metadata and meaningful child bundles; recognized v1 family migration | ordinary docs/domain/`plans` router without metadata, family metadata on wrong filename, malformed family metadata, directory-only sibling trigger, excluded/cross-boundary child |
| Ownership config | automatic ordinary owned docs path, explicit exclusion retained as evidence | attempted `owned_roots`, overlapping/case-colliding exclusions, absolute/glob/escape path, hard-boundary override |
| Git modes | regular tracked blob, staged addition, dirty tracked bytes | ignored, untracked default, symlink, gitlink, nested repo, non-regular path, identity race |
| Inventory | v2 policy/candidate digest and all classifications | stale policy, changed candidate set, changed hash, missing HEAD evidence |
| Migration | insert; reviewed rename+insert; idempotent reapply | dirty/branch-owned, target collision, ambiguous/unowned, rollback/recovery, partial coverage |
| Remote baseline | arbitrary-depth regular blobs | tree/read failure, incompatible discovery version, unsafe modes, invalid ownership config |
| Branch overlay | added/changed nested plan; rename into docs; metadata-qualified family change | pure same-class rename suppression, kind-changing rename, branch config broadening, unchanged README plus new sibling, deletion no tombstone |
| Cache | complete v2 atomic replacement | v1 cache under v2, incomplete scan preservation, schema mismatch, concurrent readers |
| Portability | slash normalization and exact case on all platforms | Unicode/case-fold path collisions, Windows reserved/path behavior, CRLF/LF mixed proof |
| Scale/security | 100k paths, 10k Markdown, batch streaming | cancellation, truncated NUL stream, oversized blobs without unbounded memory, hostile blob contents stay inert |

## Manual Review Packet

### Approved Contract Decisions

- `D-01` exact arbitrary-depth `docs` segment.
- `D-02` candidate OR signals and canonical AND validity.
- `D-03` metadata identity, filename kind, header lifecycle.
- `D-04` legacy/mixed/canonical semantics, with mixed optional rather than mandatory.
- `D-05` Git index/tree authority plus non-authoritative untracked preview.
- `D-06` automatic discovery, hard-unowned boundaries, exclusions only, and no `owned_roots`.
- `D-07` explicit discovery-v2 activation and existing-repository compatibility.
- `D-08` metadata v1 plus versioned changed-semantics machine evidence.
- `D-09` exact README plus valid metadata for family authority.
- `D-10` preserved historical cutover/register evidence and normal direct legacy-to-canonical route.
- `D-11` default-config-controlled remote authority and version-honest completeness.

### Remaining Non-Blocking Follow-Ups

- Select the later v1 removal release (`OQ-01`).
- Refine the conventional ambiguity-segment set from evidence during planning (`OQ-02`).

### Review Result

All implementation-shaping decisions `D-01` through `D-11` are resolved. The two remaining
questions are non-blocking release/evidence refinements. This discovery is unconditional input to
implementation planning, but it grants no implementation, release, publication, activation, or HQ
migration authority.

## Implementation-Planning Handoff

Handoff state: unconditional for drafting the bounded producer implementation plan below.

### Implementation Capability Scope - Shared Discovery And Ownership

Capability:
Discover and validate plan candidates across owned documentation trees at arbitrary depth with one
portable ownership-aware classifier.

Primary workflow:
Local plan authors and validators use list/validate/inventory; all consume the same Git-backed path
and Markdown classifier.

Must cover:

- config discovery version, hard-unowned boundaries, and exclusions-only policy;
- Git index/tree enumeration;
- exact filename, marker, mode, family, error, path, case, symlink, gitlink, nested-repo, ignored,
  and untracked-preview semantics; and
- deterministic structured/text output.

Explicitly out of scope:

- a central docs-root allowlist;
- an `owned_roots` inclusion override or hybrid root registry;
- new plan metadata fields or lifecycle states;
- execution log/attachment records.

Deferred but non-blocking:

- v1 removal release: `OQ-01`;
- evidence-based ambiguity list refinement: `OQ-02` during planning.

Preserve decisions:

- `D-01` through `D-09`.

Planner must not reinvent:

- exact `docs` segment semantics;
- OR candidate/AND validity split;
- identity/kind/lifecycle authority; or
- upgrade gate, exclusions-only ownership role, or metadata-qualified family authority.

Feature-level success evidence:

- local classifier/mode/ownership/family test matrix and cross-platform command goldens pass.

### Implementation Capability Scope - Migration And Compatibility

Capability:
Preview, review, and safely adopt discovery v2 without changing frozen history or overwriting
concurrent work.

Primary workflow:
Repository owner inventories prospective v2, reviews a ledger, dry-runs, applies approved
metadata/rename changes, validates, then activates v2.

Must cover:

- prospective v2 commands;
- output/ledger versioning and policy/candidate hashes;
- direct complete legacy-to-canonical migration and optional mixed material-change proof;
- reviewed rename plus metadata insertion;
- full candidate protection and projected coverage; and
- release/migration guidance for all existing modes.

Explicitly out of scope:

- automatic migration during sync/upgrade;
- rewriting frozen registers;
- changing semantic IDs on rename.

Preserve decisions:

- `D-04`, `D-05`, `D-07`, `D-08`, `D-10`.

Planner must not reinvent:

- original cutover meaning, content chronology, branch ownership, dry-run/apply approval, or
  transaction safety.

Feature-level success evidence:

- legacy/mixed/canonical upgrade fixtures prove no upgrade-only break, complete prospective
  inventory, direct complete legacy-to-canonical adoption, optional mixed grandfathering proof,
  zero-write blockers, idempotent migration, and explicit v2 activation.

### Implementation Capability Scope - Portfolio V2 Completeness

Capability:
Produce complete default and branch observations across arbitrary-depth owned docs under the v2
contract without weakening remote security or cache honesty.

Primary workflow:
Coordination home scans configured member remotes, validates v2 membership/discovery evidence,
derives observations, and atomically publishes only complete v2 cache data.

Must cover:

- commit-tree candidate scanning;
- repository-wide changed-path classification;
- eligibility-aware rename behavior;
- member discovery-version compatibility;
- cache/catalog schema v2; and
- existing redaction, transport, concurrency, stale, conflict, and incomplete-scan behavior.

Explicitly out of scope:

- branch-controlled ownership expansion;
- checkout or execution of remote content;
- deletion tombstones in the first v2 version;
- strategic-overlay rewrites.

Preserve decisions:

- `D-01` through `D-03`, `D-06` through `D-11`.

Planner must not reinvent:

- pushed-refs-only visibility, default baseline, branch overlay, exact membership, or cache
  preservation semantics.

Feature-level success evidence:

- remote multi-root, branch, cache-version, incomplete-member, inert-content, and cross-platform
  portfolio suites pass.

## Recommended Implementation-Plan Scope

Draft one producer implementation plan covering the three capability scopes above plus managed
doctrine/resource parity, consumer-impact record, migration/release notes, focused low-context
route probes, and release-candidate validation. The plan must distinguish approved external
semantics from internal implementation choices still to design, including exact schema field
placement, parser architecture, batch Git plumbing, the non-authoritative family grouping check,
and the evidence-backed conventional ambiguity list. Keep HQ execution in its own existing
consumer plan after a Kit release. Do not combine Kit implementation with HQ migration or release
publication authority.

## Revision Notes

- 2026-08-06: Created the manual-review-ready discovery from current producer doctrine, Go CLI,
  schemas, fixtures, tests, current canonical validation, and focused package-test evidence;
  recommended ownership-aware discovery v2, compatibility gating, migration semantics, remote
  completeness, HQ adoption sequence, acceptance criteria, and bounded capability scopes;
  clarified case-variant anomaly detection and ignored-untracked versus already-tracked behavior
  during the main-thread review pass.
- 2026-08-06: Incorporated approved `D-05` through `D-11`: Git index/tree authority with preview-
  only untracked files, automatic exact-`docs` discovery with hard boundaries and exclusions only,
  explicit v2 activation, version-honest machine evidence, metadata-qualified families, optional
  mixed mode with a normal direct legacy-to-canonical route, default-config-controlled portfolio
  completeness, and the paused HQ adoption sequence. Removed `owned_roots`, mandatory hybrid/
  mixed phases, conditional-review language, and directory-created family authority.
