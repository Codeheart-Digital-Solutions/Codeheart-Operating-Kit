Last updated: 2026-08-06T21:13:00Z (UTC)
Created: 2026-08-06
Status: active
Execution log: ownership-aware-plan-catalog-discovery_execution_log.md

# Document Header

## Ownership-Aware Repository-Wide Plan Catalog Discovery Implementation Plan

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.ownership-aware-plan-catalog-discovery
  kind: implementation
  purpose: Implement discovery-v2 plan catalog coverage across repository-owned documentation trees at arbitrary depth.
  first_cataloged: 2026-08-06T19:58:42Z
  catalog_metadata_updated: 2026-08-06T19:58:42Z
  products:
    - codeheart-operating-kit
  capabilities:
    - semantic-plan-catalog
    - portfolio-coordination
  strategic_themes:
    - ownership-aware-discovery
  relations:
    - kind: depends-on
      target: codeheart-operating-kit.discovery.ownership-aware-plan-catalog-discovery
```
<!-- END CODEHEART PLAN METADATA -->

Overview: Implement the approved discovery-v2 contract from
`ownership-aware-plan-catalog-discovery_discovery_doc.md`. One shared classifier will consume Git
index regular-blob paths locally and selected commit-tree regular-blob paths remotely, recognize
eligible Markdown under an exact lowercase `docs` segment at arbitrary depth, discover candidates
from supported filenames or genuine metadata, and apply the same candidate set to local commands,
migration, branch ownership, remote overlays, and portfolio completeness.

The implementation preserves existing repositories on discovery v1 until explicit activation.
Fresh repositories may write discovery v2. Ordinary tracked docs paths are owned by default;
hard-unowned boundaries and explicit `excluded_roots` are the only ownership controls in the first
implementation. There is no `owned_roots` setting, docs-root registry, automatic activation,
forced mixed-mode transition, plan-register rewrite, or second cutover revision.

The user activated this plan for goal-style implementation on 2026-08-06. Execution is authorized
on `codex/operating-kit-multi-root-plan-catalog` after the bounded plan, execution-log, and nearest-
index activation checkpoint validates and is normally pushed. This authority covers the producer
implementation and planned review checkpoints. It does not authorize a pull request, merge,
release, tag, consumer upgrade, Codeheart-HQ change, destructive Git action, history rewrite, or
force push.

Planning-checkpoint validation note: the current implementation-planning runbook requires the
generic `# Document Header` top-level section with the canonical title beneath it, while the current
catalog validator labels that layout `legacy_title_layout`. The warning is non-blocking and the
catalog remains valid. Reconciliation of that pre-existing runbook/validator mismatch is outside
this feature plan.

Essential context:

| Source | Why it matters |
| --- | --- |
| `ownership-aware-plan-catalog-discovery_discovery_doc.md` | Approved D-01 through D-11, exact discovery semantics, compatibility model, risks, acceptance criteria, and implementation capability scopes. |
| `AGENTS.md` | Producer authority, public-core boundary, managed-resource rules, and repository safety constraints. |
| `docs/repo/runbooks/change-operating-kit.md` | Required producer change route, consumer-impact classification, parity, validation, and release-note gates. |
| `docs/repo/runbooks/release-operating-kit.md` | Later release-candidate, reproducibility, platform, and publication boundaries; publication is outside this plan's authority. |
| `docs/repo/reference/placement-contract.md` | Managed, scaffold, consumer-owned, generated, and local-machine placement boundaries. |
| `docs/repo/reference/consumer-impact-classification.md` | Required instruction, validator, migration, placement, and safety impact evidence. |
| `components/planning-workflows/managed/reference/plan-catalog-format.md` | Current fixed-root catalog contract, metadata, modes, and derived-view behavior to revise. |
| `components/planning-workflows/managed/reference/planning-document-lifecycle.md` | Lifecycle header, record kinds, bundles, family behavior, and register compatibility. |
| `components/planning-workflows/managed/reference/portfolio-coordination-format.md` | Current remote membership, observation, branch overlay, completeness, and cache contract. |
| `components/planning-workflows/managed/runbooks/migrate-plan-catalog.md` | Current mandatory legacy-to-mixed route, frozen baseline, ledger, write, recovery, and validation procedure. |
| `components/planning-workflows/managed/runbooks/refresh-portfolio-catalog.md` | Current completeness, stale-cache disclosure, external-read, and recovery procedure. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Routing cards for catalog migration and portfolio refresh; changed routes require a low-context probe. |
| `components/agent-interface/managed/reference/operational-recipe-maturity.md` | L1 runbook and L3 Go-command maturity, evidence, blocker, and recovery expectations. |
| `components/agent-interface/managed/reference/runbook-authoring-standard.md` | Agent-facing intention, lane, authority, stop, evidence, and recovery shape for changed runbooks. |
| `internal/plancatalog/discover.go` | Current `Enumerate`, `Discover`, `FormalPathKind`, safe local reads, and filesystem-derived family qualification. |
| `internal/plancatalog/metadata.go` and `validate.go` | Current header/metadata parser, fence handling, expected-kind validation, duplicates, relations, and stable problem shape. |
| `internal/plancatalog/view.go`, `inventory.go`, and `migrate.go` | Current config mode parsing, structured views, fixed-root branch evidence, v1 ledger, and guarded migration. |
| `internal/commands/plans.go` | Current CLI flags and output for list, validate, inventory, migrate, and remote overlays. |
| `internal/portfolio/scanner.go` and `mirror.go` | Current fixed-prefix remote tree listing, branch diff, duplicated family classifier, inert Git transport, and observation construction. |
| `internal/portfolio/membership.go`, `source.go`, and `store.go` | Default-branch membership, catalog/cache models, completeness, and atomic prior-cache preservation. |
| `schemas/kit-config.schema.json` | Config-v1 contract containing catalog mode and cutover revision but no discovery version or exclusions. |
| `schemas/plan-metadata.schema.json` | Metadata schema v1, which remains unchanged because identity semantics do not change. |
| `schemas/plan-migration-ledger.schema.json` | Ledger v1 without policy, candidate-set, target-path, or discovery-version binding. |
| `schemas/plan-catalog.schema.json` | Cache/catalog v1 without discovery-version completeness evidence. |
| `components/planning-workflows/component.yaml`, `resources.go`, and `tests/test_packaging_resources.py` | Managed source declarations, Go embedding, legacy packaged-resource mirrors, and byte-parity validation. |
| `.github/workflows/validate.yml` | Existing macOS and real-Windows Go, CLI, route, schema, packaging, and release-asset validation surface. |

Table of contents:

- Section 1 - Foundation
- Section 2 - Strategy
- Section 3 - Execution Plan
- Section 4 - Future Planning
- Revision Notes

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

Deliver discovery v2 as a compatible, explicitly activated catalog capability whose local and
remote results are complete, deterministic, ownership-safe, and migration-ready across any
repository-owned Markdown path containing an exact lowercase `docs` segment.

Completion is proven when:

- one `internal/plancatalog` classifier receives neutral Git blob descriptors and produces the
  same path eligibility, ownership, candidate signal, expected kind, and stable problems for local
  index paths, committed trees, migration projection, remote defaults, and branch overlays;
- local authority comes from Git index regular-blob membership and remote authority from selected
  commit-tree regular blobs, while opt-in untracked authoring preview is labeled non-authoritative
  and excluded from completeness, hashes, migration writes, and remote evidence;
- root `docs/**` plus arbitrarily nested paths such as `products/a/packages/b/docs/**` are found by
  segment comparison without `*/docs/**`, recursive-glob, or fixed-depth semantics;
- supported discovery and implementation filenames or genuine metadata discover candidates, and
  canonical validity requires supported filename plus valid metadata;
- exact `README.md` plus genuine valid family metadata establishes v2 family authority, directory
  structure only validates an existing declaration, and unchanged README bytes cannot gain branch
  authority from a new sibling directory;
- `.codeheart/kit/`, `.codeheart/local/`, `.codeheart/user/`, symlinks, gitlinks, nested
  repositories, paths outside the repository, and ignored/untracked authority remain hard-unowned;
- conventional fixture, vendor, generated, dependency, and build-like candidate paths appear as
  prospective blockers until moved or listed in `excluded_roots`, and exclusions remain visible
  inventory evidence;
- existing repositories stay on v1 after Kit upgrade, prospective v2 reads require no config
  mutation, explicit activation controls adoption, and fresh repository config writes v2;
- structured views, inventories, migration ledgers, catalog/cache evidence, and member completeness
  carry versions and policy/candidate provenance wherever their meaning changes, while plan
  metadata stays schema v1;
- legacy, mixed, and canonical repositories retain compatibility, including unchanged historical
  cutover revisions and frozen registers, optional mixed grandfathering, and a direct
  legacy-to-canonical route when every formal record receives valid metadata;
- complete portfolio scans require compatible v2 evidence from every required member, inaccessible
  and v1 members make the attempt incomplete, and the last complete cache remains preserved and
  clearly historical;
- managed references, runbooks, routing, templates, source mirrors, consumer-impact evidence,
  migration notes, and release-readiness evidence describe the implemented behavior exactly;
- the full unit, integration, CLI, schema, fixture, remote, portfolio, packaging, routing,
  performance, macOS, Linux, and real-Windows matrix passes; and
- no Codeheart-HQ file changes occur in this producer plan, while a release-bounded HQ validation
  handoff is explicit and executable after an approved Kit release.

## 1.2 Project And Problem Context

The Kit already allows plans beneath product, package, business, and source-area documentation
trees, yet every executable catalog path currently begins at `docs/repo/plans/`. The fixed prefix
is repeated across local discovery, prospective path checks, branch ownership, migration,
inventory protection, remote baselines, remote branch diffs, and family inference. A partial change
would create inconsistent completeness and security claims.

Discovery v2 is a consumer-affecting capability change with these expected impact classes:

- `instruction-only change` for managed catalog, lifecycle, migration, refresh, routing, and
  consumer-template guidance;
- `validator-only change` for path, ownership, candidate, family, mode, version, and completeness
  validation;
- `consumer migration required` because existing repositories adopt v2 only after prospective
  review and explicit activation;
- `security or safety policy change` because Git modes, symlink/gitlink/nested-repository
  boundaries, inert remote reads, and exclusions become catalog authority controls; and
- placement review against `breaking placement-contract change`, with the expected conclusion that
  existing ownership-first placement is honored rather than moved automatically.

No central docs-root allowlist is introduced. No plan register becomes current authority again.
No upgrade writes consumer activation. No remote branch controls discovery scope. No v1 cache is
relabelled as v2-complete.

## 1.3 Current State Analysis

### Existing Systems, Constraints, And Problems

| Surface | Current symbol/evidence | Required change |
| --- | --- | --- |
| Local enumeration | `internal/plancatalog/discover.go: Enumerate` uses `filepath.WalkDir` below `docs/repo/plans`. | Replace path authority with Git index regular-blob enumeration and shared v1/v2 policy selection. |
| Filename and family | `kindForFilename`, `FormalPathKind`, and `qualifiesAsFamily` use suffixes plus filesystem sibling shape. | Keep exact suffixes, make family authority metadata-qualified, and remove shape-created authority. |
| Metadata parsing | `ParseDocument`, `lineIndexesOutsideFences`, and `markdownFence` parse only files already chosen by filename. | Scan eligible tracked Markdown for genuine markers, ignore fenced and indented examples, and retain malformed marker candidates. |
| Mode validation | `Discover` and `ValidateRecords` apply legacy/mixed/canonical behavior only to fixed-root filename candidates. | Apply mode validity to the shared v2 candidate set and preserve stable problem codes. |
| Settings | `RepositorySettings` and `LoadRepositorySettings` read mode, repository ID, and cutover revision. | Add discovery version, exclusions, canonical policy digest, and pure default-branch config decoding. |
| Local view | `BuildView`, `WriteViewJSON`, and `plansValidationOutput` emit schema v1 without discovery provenance. | Version changed structured output and expose active/target discovery evidence. |
| Inventory | `BuildInventory`, `activeBranchTouches`, and output protection reuse fixed-root candidates and path-limited Git diff. | Bind policy/candidate digests, all classifications, repository-wide branch touches, and preview evidence. |
| Migration | `BuildMigrationPlan` rejects legacy mode, accepts ledger v1, inserts metadata in place, and validates fixed-root coverage. | Accept reviewed v2 migration in legacy, support guarded rename plus insert, bind policy/candidate hashes, and prove projected canonical coverage. |
| CLI | `runPlansValidate`, `runPlansList`, `runPlansInventory`, and `runPlansMigrate` have no target-version or preview flags. | Add read-only prospective flags, preview provenance, v2 ledger application, and stable text/JSON output. |
| Remote tree | `GitRepository.ListFiles(ref, prefix)` and `collectObservations` enumerate only `docs/repo/plans`. | Enumerate one selected tree, adapt regular blobs to the shared classifier, and stream only required bytes. |
| Remote branch diff | `GitRepository.ChangedPaths` passes the fixed prefix and discards same-byte rename semantics without ownership/kind context. | Produce repository-wide old/new change evidence and classify semantic eligibility with default-branch policy. |
| Remote family | `remoteFormalKind` and `remoteFamilyQualified` duplicate local filename and sibling logic. | Remove the duplicate classifier and consume metadata-qualified shared family results. |
| Membership/completeness | `EvaluateMembership`, `Catalog`, and schema v1 do not bind member discovery version. | Carry default-branch discovery version/policy evidence and fail v2 completeness for incompatible required members. |
| Cache | `StoreCompleteCatalog` validates and atomically stores catalog schema v1. | Preserve v1 cache as historical evidence and replace only with a complete compatible v2 result. |
| Config creation | `initialConfig`, `components.WriteDefaultState`, and the Python compatibility writer emit empty `component_settings`. | Make new repositories explicitly v2 while leaving existing absent settings as v1. |
| Managed doctrine | Catalog and migration docs require the fixed root and mixed route. | Publish v2 semantics, direct complete migration, optional mixed route, exclusions-only ownership, and remote completeness. |
| Tests | Current Go/Python suites cover the fixed root, fences, mixed baseline, transactions, remote branches, cache, and platform containment. | Retain those baselines and add repository-wide, metadata-only, exclusion, family, version, preview, scale, and cross-platform coverage. |

### Target Systems And Ownership Boundaries

- `internal/plancatalog` owns neutral path, mode, ownership, marker, filename, family, and problem
  semantics. It must not import `internal/portfolio`.
- Local Git-index and safe-worktree adapters live in `internal/plancatalog`; remote Git-tree and
  branch-change adapters remain in `internal/portfolio` and feed the shared classifier.
- `schemas/kit-config.schema.json` owns the public discovery-version and exclusions shape.
- Structured output/cache/ledger contracts own their own version numbers; metadata schema v1 is
  unchanged.
- `internal/reconcile` remains the guarded transaction engine. Migration composes existing
  create/replace/remove actions with exact byte preconditions; it does not add an unreviewed file
  mutation path.
- `components/planning-workflows/` owns managed catalog and adoption doctrine.
- `components/agent-interface/` owns routing cards and low-context route expectations.
- Producer source under `components/` and `templates/` is authoritative; matching
  `src/codeheart_operating_kit/resources/` files are parity mirrors only.
- Codeheart-HQ owns its eventual inventory, semantic decisions, activation, remote refresh, and
  canonical cutover after release.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Expected producer surface:

```text
internal/
├── plancatalog/
│   ├── classifier.go                                  # create: neutral v1/v2 classifier and policy
│   ├── git_index.go                                   # create: index path universe and batch Git evidence
│   ├── family.go                                      # create: metadata-qualified family validation
│   ├── model.go                                       # modify: versions, ownership, signals, provenance
│   ├── discover.go                                    # modify: shared classification and safe local reads
│   ├── metadata.go                                    # modify: genuine marker scan and indented-code handling
│   ├── validate.go                                    # modify: stable v2 diagnostics and mode rules
│   ├── view.go                                        # modify: settings, prospective target, structured v2 view
│   ├── inventory.go                                   # modify: policy/candidate digests and branch evidence
│   ├── migrate.go                                     # modify: v2 ledger, direct migration, guarded rename
│   ├── legacy.go                                      # modify: v1 family/register compatibility only
│   ├── plancatalog_test.go                            # modify: shared classifier and schema cases
│   ├── ep02_test.go                                   # modify: modes, inventory, migration, safety
│   └── classifier_benchmark_test.go                   # create: deterministic scale and memory evidence
├── portfolio/
│   ├── mirror.go                                      # modify: whole-tree listing and old/new change evidence
│   ├── scanner.go                                     # modify: shared classifier adapter and v2 observations
│   ├── membership.go                                  # modify: default-branch discovery-policy evidence
│   ├── config.go                                      # modify: home discovery settings and version checks
│   ├── source.go                                      # modify: v2 member/cache/completeness fields
│   ├── store.go                                       # modify: historical v1 cache and v2 atomic replacement
│   └── portfolio_test.go                              # modify: remote depth, branch, member, cache, security
├── commands/
│   ├── plans.go                                       # modify: prospective/preview CLI and output versions
│   ├── lifecycle.go                                   # modify: fresh config defaults to discovery v2
│   └── commands_test.go                               # modify: flags, goldens, config defaults, redaction
├── components/components.go                           # modify: fresh compatibility config default
└── state/schema.go                                    # modify: explicit v1/v2 schema dispatch

schemas/
├── kit-config.schema.json                             # modify: discovery version and excluded_roots
├── plan-metadata.schema.json                          # unchanged: metadata remains v1
├── plan-migration-ledger-v1.schema.json               # create: retained historical v1 validation
├── plan-migration-ledger.schema.json                  # modify: v2 policy/candidate/target binding
├── plan-catalog-v1.schema.json                        # create: retained historical cache validation
└── plan-catalog.schema.json                           # modify: v2 completeness and member provenance

components/
├── planning-workflows/
│   └── managed/
│       ├── README.md                                  # modify: discovery-v2 route
│       ├── reference/
│       │   ├── plan-catalog-format.md                 # modify: canonical v2 contract
│       │   ├── planning-document-lifecycle.md         # modify: metadata-qualified families
│       │   └── portfolio-coordination-format.md       # modify: version-complete remote evidence
│       └── runbooks/
│           ├── maintain-plan-register.md              # modify: historical register and optional mixed
│           ├── migrate-plan-catalog.md                # modify: prospective/direct migration workflow
│           └── refresh-portfolio-catalog.md           # modify: v2 member/cache completeness
└── agent-interface/
    └── managed/reference/
        └── operation-routing-and-dispatch.md          # modify: v2 migration/refresh routing cards

templates/consumer-docs/repo/README.md                 # modify: explicit activation and nested docs visibility
src/codeheart_operating_kit/components.py              # modify: fresh compatibility config default
src/codeheart_operating_kit/resources/                 # modify: byte-identical mirrors of changed managed sources

tests/
├── fixtures/plans/multi-root-repository/              # create: root/deep/excluded/ambiguous/family cases
├── fixtures/portfolio/                                # modify: member discovery-version configurations
├── test_json_schemas.py                               # modify: v1/v2 config, ledger, and catalog contracts
├── test_packaging_resources.py                        # modify: source/mirror parity coverage
├── test_routing.py                                    # modify: v2 migration and refresh routing probes
└── test_sync_check.py                                 # modify: upgrade preserves v1 and fresh init writes v2

docs/repo/plans/ownership-aware-plan-catalog-discovery/attachments/
├── consumer-impact-record.md                          # create: reviewed impact classification
├── discovery-v2-consumer-migration.md                 # create: release-bounded adoption sequence
├── release-note-input.md                              # create: version-neutral release narrative
└── release-readiness-evidence.md                      # create: bounded producer validation evidence
```

No file is deleted by the implementation strategy. The v1 schema copies preserve historical
interpretation. The release-prep attachment supplies version-neutral input; `manifest.yaml`,
`internal/version/version.go`, tags, packs, public release notes, and published artifacts remain
owned by a separately authorized release task.

## 2.2 Open Questions And Assumptions Requiring Clarification

### OQ-01 - Eventual Discovery-V1 Removal Release

- `BLOCKER: no`.
- Affects: EP-01, EP-06, EP-08.
- Unlocks: the later release that may remove v1 readers, v1 schemas, compatibility notices, and
  upgrade fixtures.
- Recommended default: retain discovery v1 for at least one released compatibility cycle and keep
  removal outside this implementation plan. Removal requires separate discovery/release approval.

### OQ-02 - Later Conventional Ambiguity Segment Refinement

- `BLOCKER: no`.
- First-implementation decision: freeze the exact discovery-v2 ambiguity segments as `vendor`,
  `third_party`, `third-party`, `external`, `deps`, `dependencies`, `node_modules`, `fixtures`,
  `testdata`, `examples`, `generated`, `.generated`, `build`, `dist`, `out`, `target`, `.venv`, and
  `venv`.
- EP-02 through EP-08 must not add or remove segments from this set. Matching a plan signal beneath
  one of these segments produces prospective-blocker evidence until the path is moved or its root is
  explicitly listed in `excluded_roots`.
- Deferred question: later additions or removals require public-safe adoption evidence, explicit
  positive and negative fixtures, prospective inventory impact, release notes, and separate review.

### Assumptions

- `A-01`: Git remains a prerequisite for inventory, migration, and portfolio behavior.
- `A-02`: Config schema version 1 can add optional planning-workflow properties without changing
  root config identity.
- `A-03`: Existing reconcile create/replace/remove actions can express reviewed rename plus content
  insertion atomically with exact source and target preconditions.
- `A-04`: The existing maximum portfolio concurrency of four remains sufficient.
- `A-05`: The current GitHub Actions macOS and Windows jobs remain authoritative supported-platform
  environments; this plan adds an Ubuntu semantic-validation job without adding a Linux release
  asset or Linux distribution support.
- `A-06`: HQ remains paused and available only as a later consumer-owned validation target after an
  approved release.

## 2.3 Architectural Decisions With Reasoning

### AD-01 - Neutral Shared Classifier In `internal/plancatalog`

1. Problem being solved: local filesystem traversal, formal-path checks, and remote scanners
   independently define plan authority.
2. Simplest working solution: add a classifier that accepts normalized Git blob descriptors,
   content readers, discovery settings, and an optional selected-path set, then emits candidates,
   records, exclusions, blockers, and stable problems.
3. Six-to-twelve-month change: additional Git providers can adapt tree evidence without adding a
   new catalog policy implementation.
4. Rationale: `internal/portfolio` already imports `internal/plancatalog`; keeping policy in
   `plancatalog` avoids a cycle and makes local/remote parity testable.
5. Alternatives rejected: a new generic package would split metadata/validation ownership, while
   duplicating local and remote classifiers preserves the current inconsistency.

### AD-02 - Git Path Universe With Safe Content Adapters

1. Problem being solved: filesystem walks include untracked/ignored content and cannot represent a
   remote commit tree consistently.
2. Simplest working solution: enumerate local index stages with NUL-safe Git output, reject
   non-regular/ambiguous stages, bind tracked worktree bytes safely for local validation, and adapt
   remote `ls-tree` regular blobs to the same descriptor shape.
3. Six-to-twelve-month change: batch-object APIs can be optimized without changing classifier
   semantics.
4. Rationale: Git spelling, modes, and revisions are portable authority evidence; safe worktree
   reads preserve current local-edit validation.
5. Alternatives rejected: recursive globs, host filesystem walks, and per-file `git show` calls are
   inconsistent, unsafe, or unbounded.

Local nested-repository evidence and remote tree evidence have different representable boundaries.
The local adapter must detect a nested `.git` directory or file below the outer repository root,
without following it, and mark indexed descendants unsafe/incomplete. The outer repository's own
Git directory is not a nested boundary. A selected remote commit tree has no worktree-only `.git`
state: mode `160000` gitlinks are hard-unowned, while `100644` and `100755` blobs belong to the
selected outer commit. The remote adapter must not claim to detect an unrepresented nested worktree.

### AD-03 - Owned By Default With Exclusions Only

1. Problem being solved: a root allowlist silently omits future domains, while tracking alone can
   expose plan-shaped embedded content.
2. Simplest working solution: classify ordinary tracked exact-`docs` Markdown as owned, hard-block
   managed/local/user state and unsafe Git boundaries, surface conventional signals as prospective
   blockers, and accept only explicit `excluded_roots` resolutions.
3. Six-to-twelve-month change: an inclusion override can be discovered separately after a concrete
   legitimate use case exists.
4. Rationale: the model finds future owned roots automatically and keeps misleading locations
   visible without normalizing poor placement.
5. Alternatives rejected: `owned_roots`, mandatory root markers, and a complete docs-root registry
   add silent omission and overlap risk without present evidence.

### AD-04 - Signal Discovery Separated From Canonical Validity

1. Problem being solved: current filename-only enumeration cannot diagnose metadata-bearing wrong
   filenames and structure-created families can appear without explicit authority.
2. Simplest working solution: scan eligible Markdown once for genuine markers, combine filename
   and metadata signals, validate expected kind after discovery, and require exact `README.md` plus
   valid family metadata for v2 families.
3. Six-to-twelve-month change: a fuller CommonMark parser may replace the bounded scanner while
   preserving marker/position semantics.
4. Rationale: invalid candidates remain visible, identity stays metadata-owned, kind stays
   filename/qualification-owned, and lifecycle stays header-owned.
5. Alternatives rejected: filename-only scanning hides malformed placement; metadata-only
   authority removes human-readable kind; directory-only family rules misclassify routers.

Discovery v2 intentionally has no child-count or directory-shape validity requirement for a family.
Family authority and validity come from the exact `README.md`, valid family metadata, matching
family ID/kind, eligible owned path, and normal header rules. Child appearance or disappearance is
visible only through independently classified records; it cannot create, remove, validate, or
invalidate family authority. Relationships between records remain metadata-owned.

### AD-05 - Explicit Activation And Version-Honest Evidence

1. Problem being solved: upgrade-only semantic expansion can invalidate canonical repositories and
   old caches can falsely claim completeness.
2. Simplest working solution: absence means discovery v1 for existing repositories, fresh config
   writes v2, read commands accept a target-version preview, and each changed structured contract
   carries its own v2 evidence.
3. Six-to-twelve-month change: discovery v1 can be removed under OQ-01 after adoption evidence.
4. Rationale: compatibility is explicit, prospective review is write-free, and historical evidence
   remains interpretable.
5. Alternatives rejected: automatic upgrade activation and relabelling schema-v1 cache bytes make
   completeness claims unverifiable.

### AD-06 - Reviewed Migration Without Mandatory Mixed Mode

1. Problem being solved: current migration rejects legacy mode and requires mixed before metadata
   writes even when all candidates can migrate coherently.
2. Simplest working solution: ledger v2 binds discovery policy, candidate digest, revision, hashes,
   target paths, and target mode; migration may apply reviewed metadata/rename actions in legacy;
   activation remains a separate explicit config change.
3. Six-to-twelve-month change: mixed stays available for unusually large deferred migrations and
   can be retired only with separate compatibility evidence.
4. Rationale: HQ and complete repositories use the shorter safe route while frozen registers and
   cutover proofs remain meaningful only for intentional grandfathering.
5. Alternatives rejected: a second cutover revision corrupts historical meaning; automatic config
   activation combines authority with file mutation; forced mixed adds ceremony without safety.

### AD-07 - Default-Branch Remote Policy And Complete-Member Gate

1. Problem being solved: branch config could broaden authority and v1 members could be counted as
   v2-complete.
2. Simplest working solution: parse discovery version/exclusions from each member default branch,
   apply that policy to its branch overlays, and require compatible v2 evidence from every required
   member before cache replacement.
3. Six-to-twelve-month change: later catalog versions can negotiate explicit compatible ranges.
4. Rationale: unmerged branches contribute plan observations but never redefine membership or
   authority; incomplete current scans preserve historical cache honestly.
5. Alternatives rejected: branch-controlled policy and best-effort mixed-version completeness are
   vulnerable to false authority and incomparable candidate sets.

### AD-08 - L3 Go Mechanics With Managed L1 Routes

1. Problem being solved: the capability spans deterministic Git mechanics, semantic review,
   migration authorization, and external-read failure handling.
2. Simplest working solution: keep reusable mechanics in existing Go commands, update agent-facing
   L1 runbooks and routing cards, and add no standalone conversion script.
3. Six-to-twelve-month change: command internals may expose narrower APIs after operational use,
   without changing the route or evidence model.
4. Rationale: the existing CLI already owns parsing, transactions, remotes, cache, and structured
   output; a new script would duplicate safety boundaries.
5. Alternatives rejected: policy-only docs leave the feature unimplemented, and an L2 script would
   bypass the reviewed ledger and transaction path.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Dependencies |
| --- | --- | --- | --- |
| EP-01 | Versioned discovery settings, neutral models, config defaults, and historical schema dispatch exist. | L | none |
| EP-02 | One shared classifier discovers and validates arbitrary-depth local Git-index candidates. | XL | EP-01 |
| EP-03 | List, validate, and inventory expose authoritative v2 plus prospective and preview evidence consistently. | L | EP-02 |
| EP-04 | Guarded migration supports complete direct legacy-to-canonical preparation and optional mixed grandfathering. | XL | EP-03 |
| EP-05 | Remote defaults and branch overlays use the shared classifier with version-honest completeness and cache preservation. | XL | EP-04 |
| EP-06 | Managed doctrine, routing, templates, and packaged mirrors describe the implemented contract exactly. | L | EP-04, EP-05 |
| EP-07 | The complete fixture, unit, CLI, remote, performance, packaging, and cross-platform matrix passes. | XL | EP-01 through EP-06 |
| EP-08 | Consumer-impact, migration, release-note, release-readiness, and HQ handoff evidence is complete without publishing. | M | EP-07 |

Planned implementation checkpoints after later plan activation:

| Epic | Intended commit message | Included workstream |
| --- | --- | --- |
| EP-01 | `feat: add plan catalog discovery v2 contracts` | Models, config, schema versions, fresh defaults |
| EP-02 | `feat: classify repository-wide plan candidates` | Shared classifier, Git index, metadata, families |
| EP-03 | `feat: expose prospective multi-root plan views` | CLI, view, inventory, preview, branch ownership |
| EP-04 | `feat: migrate multi-root plan catalogs safely` | Ledger v2, guarded rename/insert, modes, recovery |
| EP-05 | `feat: scan multi-root portfolio plans` | Remote tree/branch adapters, member completeness, cache |
| EP-06 | `docs: document plan catalog discovery v2` | Managed doctrine, routing, templates, resource parity |
| EP-07 | `test: cover repository-wide plan discovery` | Full fixtures, integration, scale, platform evidence |
| EP-08 | `docs: record plan catalog v2 release readiness` | Impact, migration, release-note, HQ handoff evidence |

Each checkpoint must contain only its reviewed paths, pass its epic gate, and be normally pushed
only under execution authority. Failed gates stop forward execution. Recovery uses the existing
transaction marker and cache-preservation paths; Git history rewriting is never a recovery step.

## 3.1 EP-01 - Discovery Settings And Versioned Contracts

### A) Epic ID, Title, And Outcome

`EP-01 - Discovery Settings And Versioned Contracts`

Outcome: config, Go models, schema dispatch, and fresh-repository defaults can represent discovery
v1/v2, exclusions, candidate provenance, and changed completeness contracts without altering plan
metadata schema v1 or silently activating existing repositories.

### B) Scope

- Add `DiscoveryVersion`, discovery policy, ownership class, candidate signal, Git mode, and policy
  digest types to `internal/plancatalog/model.go`.
- Extend `RepositorySettings` and split config byte decoding from filesystem loading in
  `internal/plancatalog/view.go` so remote membership can reuse it.
- Add optional `plan_catalog_discovery_version` and exclusions-only ownership settings to config
  schema v1.
- Preserve old catalog and ledger schema files for historical interpretation, then make current
  schema paths v2.
- Make new Go and Python compatibility config writers emit discovery v2; preserve absent settings
  as v1 in existing repositories.

### C) Files Touched

- `internal/plancatalog/model.go`
- `internal/plancatalog/view.go`
- `internal/state/schema.go`
- `internal/commands/lifecycle.go`
- `internal/components/components.go`
- `src/codeheart_operating_kit/components.py`
- `schemas/kit-config.schema.json`
- `schemas/plan-migration-ledger-v1.schema.json` (new)
- `schemas/plan-migration-ledger.schema.json`
- `schemas/plan-catalog-v1.schema.json` (new)
- `schemas/plan-catalog.schema.json`
- `tests/test_json_schemas.py`
- `internal/commands/commands_test.go`
- `tests/test_sync_check.py`

### D) Acceptance Criteria And Size

- Size: `L`.
- Existing config with absent discovery settings decodes as v1 without byte mutation.
- Fresh Go and Python compatibility config generation writes explicit v2.
- Config accepts normalized unique `excluded_roots` and rejects `owned_roots`, globs, absolute
  paths, escapes, portable collisions, and overlapping exclusions.
- Metadata schema remains byte-identical.
- V1 cache and ledger instances validate only through historical dispatch and cannot be labelled
  v2-complete or applied as v2 migration authority.
- V2 structures require discovery version, policy digest, candidate-set digest, and relevant
  completeness/provenance fields.

### E) Dependencies And Critical-Path Notes

No dependency. Freeze public model and schema names before classifier work. A schema disagreement
blocks EP-02 through EP-05.

### F) Tasks Checklist

- [ ] Add `DiscoveryVersion`, `DiscoveryPolicy`, `OwnershipClass`, `CandidateSignal`, and Git provenance types to `internal/plancatalog/model.go`.
- [ ] Add pure config-byte decoding plus filesystem loading for discovery version and exclusions in `internal/plancatalog/view.go`.
- [ ] Extend `schemas/kit-config.schema.json` with discovery versions `1` and `2` plus normalized `excluded_roots`.
- [ ] Copy current ledger schema bytes to `schemas/plan-migration-ledger-v1.schema.json` and define the v2 contract in `schemas/plan-migration-ledger.schema.json`.
- [ ] Copy current catalog schema bytes to `schemas/plan-catalog-v1.schema.json` and define v2 member/completeness evidence in `schemas/plan-catalog.schema.json`.
- [ ] Add explicit v1/v2 schema constants and dispatch in `internal/state/schema.go`.
- [ ] Make `initialConfig` in `internal/commands/lifecycle.go` write fresh discovery-v2 settings.
- [ ] Make `WriteDefaultState` in `internal/components/components.go` write the same fresh settings.
- [ ] Make `write_default_state` in `src/codeheart_operating_kit/components.py` preserve fresh-config parity.
- [ ] Add config, schema-version, historical-evidence, and fresh-default cases to `tests/test_json_schemas.py`, `internal/commands/commands_test.go`, and `tests/test_sync_check.py`.
- [ ] Run `go test ./internal/plancatalog ./internal/state ./internal/commands`.
- [ ] Run `python3 -m pytest tests/test_json_schemas.py tests/test_sync_check.py`.
- [ ] Run `python3 scripts/validate-json-schemas.py`.

### G) Implementation Notes

Use deterministic slash-normalized repository-relative exclusions. Schema parsing must reject
attempted `owned_roots` even though planning-workflow settings previously allowed extra keys.
Keep root config schema version 1. Historical schemas are read-only evidence contracts.

Rollback boundary: before a release, revert the EP-01 checkpoint as one unit. After any consumer
creates v2 config, do not downgrade bytes automatically; preserve evidence and follow the later
release recovery route.

### H) Open Questions

- OQ-01 is non-blocking; retain v1 dispatch in this epic.

## 3.2 EP-02 - Shared Git-Backed Candidate Classifier

### A) Epic ID, Title, And Outcome

`EP-02 - Shared Git-Backed Candidate Classifier`

Outcome: one neutral classifier enumerates arbitrary-depth docs candidates from local Git-index
regular blobs, applies hard-unowned/exclusion/ambiguity policy, parses filename and genuine metadata
signals, and establishes metadata-qualified family authority with deterministic problems.

### B) Scope

- Add neutral descriptor/classification APIs in `classifier.go`.
- Add one NUL-safe local Git-index listing and bounded batch Git evidence in `git_index.go`.
- Refactor `Discover`, `Enumerate`, `FormalPathKind`, and safe local reads to use active v1 or v2
  policy without a filesystem authority walk.
- Expand marker detection to eligible Markdown and misplaced-metadata anomaly paths while ignoring
  backtick fences, tilde fences, indented code, and non-structural examples.
- Move family validation into `family.go`; v2 authority requires exact `README.md` plus valid
  family metadata.
- Retain legacy v1 fixed-root behavior behind the discovery-version selector.

### C) Files Touched

- `internal/plancatalog/classifier.go` (new)
- `internal/plancatalog/git_index.go` (new)
- `internal/plancatalog/family.go` (new)
- `internal/plancatalog/model.go`
- `internal/plancatalog/discover.go`
- `internal/plancatalog/metadata.go`
- `internal/plancatalog/validate.go`
- `internal/plancatalog/legacy.go`
- `internal/plancatalog/plancatalog_test.go`
- `internal/plancatalog/ep02_test.go`
- `tests/fixtures/plans/multi-root-repository/` (new)

### D) Acceptance Criteria And Size

- Size: `XL`.
- Exact lowercase `docs` directory-segment comparison finds root and arbitrarily deep paths once.
- Exact case-sensitive `<non-empty-name>_discovery_doc.md` and
  `<non-empty-name>_implementation_doc.md` filenames supply discovery/implementation kind.
- Git modes `100644` and `100755` are the only authoritative blob modes.
- Filename and metadata signals discover candidates; canonical validity requires both.
- Genuine markers inside malformed blocks remain candidates; fenced, indented, and quoted examples
  do not become authority.
- Duplicate semantic IDs across every owned root produce stable errors preserving all paths.
- Hard-unowned content is never read through symlinks, gitlinks, nested repositories, escapes, or
  managed/local/user boundaries.
- Conventional plan signals block prospective v2 until excluded or moved; excluded candidates
  stay visible with reason evidence.
- Exact `README.md` plus valid family metadata creates v2 family authority; directory shape alone
  never does.
- Zero, one, or many descendant records do not affect v2 family authority or validity.
- Output ordering is bytewise deterministic across host platforms.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 public types. The shared classifier API freezes before local CLI and remote adapter
work. Keep `internal/plancatalog` independent of `internal/portfolio`.

### F) Tasks Checklist

- [ ] Implement neutral blob descriptor, content reader, policy compiler, and classification result APIs in `internal/plancatalog/classifier.go`.
- [ ] Implement NUL-safe index-stage enumeration and batch blob metadata reads in `internal/plancatalog/git_index.go`.
- [ ] Refactor `Enumerate`, `Discover`, and `FormalPathKind` in `internal/plancatalog/discover.go` to use version-selected v1 compatibility and v2 shared classification paths.
- [ ] Add exact lowercase path-segment, extension, suffix, portable-collision, and normalization validation.
- [ ] Add local hard-unowned checks for managed state, local state, user state, symlinks, nested `.git` file/directory boundaries below the outer repository root, gitlinks, escaping paths, and non-regular sources without following unsafe entries.
- [ ] Add remote hard-unowned handling for mode `160000` gitlinks while treating only selected-tree modes `100644` and `100755` as authoritative outer-commit blobs.
- [ ] Add conventional ambiguity classification for the frozen Section 2.2 segment set plus exclusions-only resolution using the EP-01 policy.
- [ ] Extend marker scanning in `internal/plancatalog/metadata.go` for fenced code, indented code, structural placement, malformed markers, and metadata-only candidates.
- [ ] Implement exact-`README.md`, metadata-kind, semantic-ID-kind, owned-path, and header family validation plus v1-family prospective migration evidence in `internal/plancatalog/family.go` without child-count and directory-shape validity checks.
- [ ] Update `ValidateRecords` with stable metadata-missing, filename-missing, kind-mismatch, duplicate-ID, misplaced, overlap, unowned, excluded, case, portability, family, and unsafe-source behavior.
- [ ] Create multi-root fixtures covering root docs, deeply nested docs, repeated docs segments, metadata-only paths, malformed metadata, exclusions, conventional blockers, and family routers.
- [ ] Add local classifier parity, deterministic ordering, case, Unicode, symlink, gitlink, nested-repository, and hostile-content tests.
- [ ] Run `go test ./internal/plancatalog`.
- [ ] Run the local classifier fixture twice and compare byte-identical JSON output.

### G) Implementation Notes

Index membership decides authority; bound worktree bytes remain the local validation content for
tracked paths and must retain identity checks. Missing, replaced, intent-to-add, conflicted-stage,
racy, and ancestor-nested-repository paths produce explicit incomplete/unsafe evidence. Remote
commit-tree classification rejects gitlinks and does not infer worktree-only `.git` boundaries from
regular outer-commit blobs. The classifier reads inert bytes only.

The misplaced-metadata anomaly pass may inspect tracked Markdown outside eligible docs but can
only emit diagnostics; it cannot create authority. Case-variant extensions are diagnostic-only.

Rollback boundary: this epic changes no repository plan bytes. Revert the classifier checkpoint
before dependent API changes. Never recover by falling back silently from v2 to filesystem walk.

### H) Open Questions

- OQ-02 is deferred; implement the exact frozen Section 2.2 list and preserve evidence for a later
  separately reviewed refinement.

## 3.3 EP-03 - Local Views, Prospective Reads, And Authoring Preview

### A) Epic ID, Title, And Outcome

`EP-03 - Local Views, Prospective Reads, And Authoring Preview`

Outcome: list, validate, inventory, branch ownership, and inventory-target protection consume the
same classifier; users can preview discovery v2 and canonical readiness without config writes and
can opt into clearly non-authoritative untracked authoring feedback.

### B) Scope

- Thread active and target discovery settings through repository snapshots and views.
- Add `--target-discovery-version 2` to read-only list, validate, and inventory flows.
- Add `--target-catalog-mode canonical` to prospective validate/inventory readiness checks.
- Add `--include-untracked` to validate/inventory only, label preview paths, and exclude them from
  authoritative counts/digests/migration/remote evidence.
- Replace path-limited branch touch logic with repository-wide Git changes classified by the
  shared policy.
- Version view, validation, and inventory JSON where completeness/provenance meaning changes.
- Protect all actual and prospective formal v2 targets from inventory-output collisions.

### C) Files Touched

- `internal/plancatalog/view.go`
- `internal/plancatalog/inventory.go`
- `internal/plancatalog/discover.go`
- `internal/plancatalog/legacy.go`
- `internal/commands/plans.go`
- `internal/commands/commands_test.go`
- `internal/commands/testdata/plans-list.json`
- `internal/plancatalog/ep02_test.go`

### D) Acceptance Criteria And Size

- Size: `L`.
- Active v1 output remains compatible after upgrade.
- Prospective v2 commands do not edit config, plans, register, cache, or ledger.
- View/inventory output includes discovery version, target mode, policy digest, candidate digest,
  signal, ownership, Git mode/blob, source hash, problems, and authoritative/preview provenance.
- Untracked preview never contributes to valid canonical completeness or migration coverage.
- Branch touches include eligibility-changing renames and nested docs paths across the repository.
- Every read command reports the same authoritative candidate set for the same settings and bytes.

### E) Dependencies And Critical-Path Notes

Depends on EP-02. Freeze structured local outputs before migration ledger v2 is implemented.

### F) Tasks Checklist

- [ ] Extend `RepositorySnapshot`, `View`, `Inventory`, and coverage models with discovery and policy provenance.
- [ ] Add target discovery-version parsing to list, validate, and inventory in `internal/commands/plans.go`.
- [ ] Add target canonical-mode readiness parsing to validate and inventory without config writes.
- [ ] Add non-authoritative non-ignored untracked preview collection to validate and inventory with explicit provenance.
- [ ] Replace fixed-root `activeBranchTouches` with repository-wide change classification from the shared policy.
- [ ] Bind inventory candidate digests to sorted path, signal, expected kind, ownership, Git mode/blob, and source hash.
- [ ] Extend inventory-output protection to every classifier-visible present and prospective formal target.
- [ ] Update text output, JSON output, and `internal/commands/testdata/plans-list.json` with versioned deterministic fields.
- [ ] Add CLI rejection tests for unsupported target versions, unsafe exclusions, preview misuse, and conflicting flags.
- [ ] Add no-write tests proving prospective and preview commands preserve config, plans, register, cache, and transaction paths.
- [ ] Run `go test ./internal/plancatalog ./internal/commands`.
- [ ] Run prospective v1 and v2 command goldens on the multi-root fixture.

### G) Implementation Notes

Use one immutable compiled policy per command. Text output should summarize included, excluded,
prospective-blocked, misplaced, unsafe, and preview counts while JSON retains full evidence.
Prospective canonical mode is a read-only readiness lens, not activation.

Rollback boundary: command output versions must change atomically with model changes. A failed
inventory write preserves the previous artifact via existing bound-output handling. No preview
result becomes migration input.

### H) Open Questions

None. Later OQ-02 refinement remains outside this implementation and does not change CLI shape.

## 3.4 EP-04 - Guarded Migration And Mode Compatibility

### A) Epic ID, Title, And Outcome

`EP-04 - Guarded Migration And Mode Compatibility`

Outcome: ledger v2 can safely migrate the complete discovery-v2 candidate set while a repository
is still legacy, support reviewed wrong-filename renames, prove projected canonical readiness, and
retain mixed mode only for explicit frozen-baseline grandfathering.

### B) Scope

- Extend ledger records with discovery version, target mode, policy/candidate digests, current and
  target paths, target absence/hash preconditions, and reviewed ownership disposition.
- Permit v2 metadata/rename application in legacy mode when the reviewed ledger covers every
  authoritative candidate and projected canonical validation passes.
- Compose exact create/replace/remove reconcile actions for reviewed rename plus metadata insertion.
- Preserve `Created`, `Last updated`, semantic ID, frozen register bytes, and original cutover
  revision.
- Strengthen mixed unchanged-byte proof beyond line-ending conversion and reject new/renamed/
  changed filename-only plans.
- Keep activation as a separate reviewed config write after prospective validation.

### C) Files Touched

- `schemas/plan-migration-ledger.schema.json`
- `internal/plancatalog/migrate.go`
- `internal/plancatalog/inventory.go`
- `internal/plancatalog/view.go`
- `internal/plancatalog/validate.go`
- `internal/plancatalog/ep02_test.go`
- `internal/commands/plans.go`
- `internal/commands/commands_test.go`

### D) Acceptance Criteria And Size

- Size: `XL`.
- V1 ledgers are readable historical evidence but cannot apply as v2.
- Any revision, policy, candidate set, source hash, target state, branch owner, dirty state, or
  semantic decision drift blocks mutation.
- Reviewed filename-only candidates gain metadata without changing historical content timestamps.
- Reviewed metadata-only candidates can move to supported filenames with no-replace target safety
  and one transaction rollback boundary.
- Legacy repositories can reach zero projected canonical gaps without entering mixed.
- Mixed grandfathering requires unchanged register-proven cutover blobs; new or materially changed
  filename-only records fail.
- Migration never writes config, register, unrelated files, excluded candidates, preview files, or
  another repository.
- Dry-run/apply/reapply is deterministic, guarded, recoverable, and idempotent.

### E) Dependencies And Critical-Path Notes

Depends on EP-03 inventory/output freeze. Remote completeness lands in EP-05 before doctrine and
release evidence claim end-to-end readiness.

### F) Tasks Checklist

- [ ] Add discovery version, target mode, policy digest, candidate digest, target path, and target preconditions to ledger v2 decoding.
- [ ] Bind `BuildMigrationPlan` to the exact EP-03 inventory revision, policy, candidate set, source hashes, and branch ownership.
- [ ] Permit reviewed v2 metadata writes while active config remains legacy and reject incomplete projected canonical coverage.
- [ ] Generate exact replace actions for in-place metadata insertion while preserving lifecycle header timestamps.
- [ ] Generate paired no-replace create and exact-hash remove actions for reviewed filename correction in one reconcile transaction.
- [ ] Reject excluded, prospective-blocked, unowned, preview, malformed, duplicate-ID, dirty, branch-owned, and stale-ledger actions.
- [ ] Strengthen mixed proof to compare current bytes with exact cutover blobs beyond checkout line endings.
- [ ] Preserve frozen register bytes and original cutover revision in every legacy, mixed, and canonical migration path.
- [ ] Add direct legacy-to-canonical readiness, optional mixed deferral, rename collision, rollback, recovery-marker, and idempotent reapply tests.
- [ ] Add CLI text/JSON evidence for actions, skips, blockers, projected coverage, and explicit no-activation behavior.
- [ ] Run `go test ./internal/plancatalog ./internal/reconcile ./internal/commands`.
- [ ] Run dry-run, apply, recovery, and reapply fixtures with byte comparisons for headers and frozen registers.

### G) Implementation Notes

Reuse the existing reconcile transaction rather than adding a second writer. A rename is a reviewed
target create plus source remove with exact preconditions; transaction rollback restores both.
Activation edits `plan_catalog_discovery_version` and `plan_catalog_mode` only after a separate
prospective canonical validation and user-approved repository checkpoint.

Recovery boundary: preserve `.codeheart/kit.transaction.json` plus
`.codeheart/local/kit-transactions/` evidence on incomplete rollback. Follow existing check/repair
recovery; never delete markers, retry stale actions, or reconstruct a ledger from partial output.

### H) Open Questions

None. Mixed is an available compatibility path, not an unresolved strategy choice.

## 3.5 EP-05 - Remote Discovery, Branch Overlays, And Cache Completeness

### A) Epic ID, Title, And Outcome

`EP-05 - Remote Discovery, Branch Overlays, And Cache Completeness`

Outcome: portfolio defaults and unmerged branches adapt selected commit-tree regular blobs to the
shared classifier, default-branch config controls discovery authority, and only compatible complete
v2 evidence replaces the last complete cache.

### B) Scope

- Generalize `ListFiles` to one selected commit tree and return path plus Git mode.
- Generalize `ChangedPaths` to repository-wide NUL-safe old/new rename evidence with modes and byte
  comparison deferred until classifier eligibility is known.
- Parse member discovery settings from default-branch config and attach version/policy evidence to
  membership and catalog members.
- Remove `remoteFormalKind` and `remoteFamilyQualified` in favor of the shared classifier.
- Apply default-branch exclusions to branch overlays; branch config cannot broaden authority.
- Treat v1, inaccessible, unsafe, truncated, or incompatible required members as incomplete.
- Preserve v1 caches as historical and replace cache atomically only after complete v2 validation.

### C) Files Touched

- `internal/portfolio/mirror.go`
- `internal/portfolio/scanner.go`
- `internal/portfolio/membership.go`
- `internal/portfolio/config.go`
- `internal/portfolio/source.go`
- `internal/portfolio/store.go`
- `internal/portfolio/portfolio_test.go`
- `internal/plancatalog/model.go`
- `schemas/plan-catalog.schema.json`
- `tests/fixtures/portfolio/kit-config-v2-member.yaml`
- `tests/fixtures/portfolio/kit-config-v2-home.yaml`
- `tests/fixtures/portfolio/kit-config-v1-member.yaml`
- `internal/commands/plans.go`
- `internal/commands/commands_test.go`

### D) Acceptance Criteria And Size

- Size: `XL`.
- Baseline scans find root and deeply nested docs candidates with the same paths, kinds, problems,
  and IDs as local committed-tree classification.
- Remote Git reads remain inert and never checkout content, follow symlinks/gitlinks, run hooks,
  inherit unsafe Git environment, or execute repository bytes.
- Remote mode `160000` gitlinks are hard-unowned; regular selected-tree blobs are outer-commit
  authority, and remote evidence never claims detection of unrepresented worktree-only `.git`
  boundaries.
- Added and materially changed branch candidates produce observations with repository, ref,
  commit, path, content hash, verification, visibility, and source-change evidence.
- Same-byte rename suppression occurs only when old/new eligibility, ownership, and kind match.
- Rename into docs, kind-changing rename, ownership-changing rename, and metadata changes remain
  visible.
- An unchanged README never gains family authority because a sibling directory appears.
- V1 and inaccessible required members make the current scan incomplete.
- Incomplete v2 scans preserve and label the previous complete cache; a v1 cache is never current
  v2 evidence.

### E) Dependencies And Critical-Path Notes

Depends on EP-04. Reuse the frozen EP-02 classifier and EP-03 structured provenance after guarded
migration behavior is complete; execution remains linear before EP-06.

### F) Tasks Checklist

- [ ] Change `GitRepository.ListFiles` to enumerate one selected tree and return regular-blob mode evidence without a plans prefix.
- [ ] Change `GitRepository.ChangedPaths` to emit NUL-safe old/new path, status, and mode evidence across the repository.
- [ ] Parse default-branch discovery version and exclusions through the pure EP-01 settings decoder in membership evaluation.
- [ ] Adapt remote tree descriptors and bounded blob readers to the shared `internal/plancatalog` classifier.
- [ ] Remove `remoteFormalKind` and `remoteFamilyQualified` from `internal/portfolio/scanner.go`.
- [ ] Apply default-branch policy to every unmerged branch observation and ignore branch attempts to broaden authority.
- [ ] Preserve eligibility-changing and kind-changing rename observations with exact ref, commit, path, and hash evidence.
- [ ] Add discovery version and policy digest to catalog member and completeness structures.
- [ ] Mark v1, inaccessible, truncated, unsafe, and incompatible required-member evidence incomplete with stable scan codes.
- [ ] Dispatch historical v1 cache validation and prevent it from satisfying current v2 completeness.
- [ ] Preserve the prior complete cache on every incomplete attempt and atomically store only a complete validated v2 catalog.
- [ ] Add remote multi-root, exclusion, metadata-only, family, rename, hostile-content, incompatible-member, and cache-preservation tests.
- [ ] Run `go test ./internal/portfolio ./internal/commands`.
- [ ] Run remote default/branch fixture parity against the local committed-tree classifier.

### G) Implementation Notes

Keep the portfolio maximum concurrency at four. Within a repository, use bounded batch blob reads
and deterministic result ordering. Default-branch membership/config evidence is required before
branch bytes contribute observations. Pull-request enrichment remains non-authoritative.

Rollback boundary: mirror/cache writes retain existing bound-directory and atomic no-replace
safety. On incomplete scan, return current errors and preserve the previous cache without copying
its observations into the current attempt.

### H) Open Questions

None.

## 3.6 EP-06 - Managed Doctrine, Routing, And Resource Parity

### A) Epic ID, Title, And Outcome

`EP-06 - Managed Doctrine, Routing, And Resource Parity`

Outcome: installed agents and maintainers receive exact discovery-v2 authoring, migration,
portfolio, ownership, activation, stop, evidence, and recovery guidance, with producer source and
legacy packaged mirrors byte-identical.

### B) Scope

- Update catalog, lifecycle, portfolio, register, migration, and refresh guidance.
- Update planning-workflows README and consumer repo template.
- Update `planning.migrate-catalog` and portfolio-refresh route cards.
- Classify runbooks as agent-facing L1 recipes over L3 commands; preserve tooling and external-read
  preflights.
- Copy every changed managed/template source into its matching legacy packaged-resource path.
- Add a low-context nested-domain migration/refresh route probe.

### C) Files Touched

- `components/planning-workflows/managed/README.md`
- `components/planning-workflows/managed/reference/plan-catalog-format.md`
- `components/planning-workflows/managed/reference/planning-document-lifecycle.md`
- `components/planning-workflows/managed/reference/portfolio-coordination-format.md`
- `components/planning-workflows/managed/runbooks/maintain-plan-register.md`
- `components/planning-workflows/managed/runbooks/migrate-plan-catalog.md`
- `components/planning-workflows/managed/runbooks/refresh-portfolio-catalog.md`
- `components/agent-interface/managed/reference/operation-routing-and-dispatch.md`
- `templates/consumer-docs/repo/README.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/README.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/plan-catalog-format.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/planning-document-lifecycle.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/reference/portfolio-coordination-format.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/maintain-plan-register.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/migrate-plan-catalog.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/refresh-portfolio-catalog.md`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/reference/operation-routing-and-dispatch.md`
- `src/codeheart_operating_kit/resources/templates/consumer-docs/repo/README.md`
- `tests/test_packaging_resources.py`
- `tests/test_routing.py`
- `tests/test_sync_check.py`

### D) Acceptance Criteria And Size

- Size: `L`.
- Managed doctrine states exact segment, signals, validity, ownership, exclusions, family, modes,
  prospective, migration, remote completeness, and cache behavior.
- Migration guidance makes direct complete legacy-to-canonical adoption normal and mixed optional.
- Register guidance preserves historical evidence and forbids a second cutover revision.
- Refresh guidance stops current conclusions for v1/inaccessible required members and labels old
  cache historical.
- A low-context agent finds the correct migration/refresh route for a deeply nested domain plan,
  asks the ownership question, and does not select writes before evidence/approval.
- Source and packaged mirrors are byte-identical.

### E) Dependencies And Critical-Path Notes

Depends on EP-04 and EP-05 behavior freeze. Doctrine must not lead runtime semantics.

### F) Tasks Checklist

- [ ] Update `plan-catalog-format.md` with discovery-v2 configuration, exact path segments, signals, validity, errors, versions, and exclusions-only authority.
- [ ] Update `planning-document-lifecycle.md` with exact README plus valid metadata family authority and v1-family migration.
- [ ] Update `portfolio-coordination-format.md` with default-policy branch overlays, compatible-member completeness, and historical cache labels.
- [ ] Update `maintain-plan-register.md` with historical register scope, optional mixed grandfathering, and no second cutover revision.
- [ ] Rewrite `migrate-plan-catalog.md` for prospective v2 inventory, reviewed ledger v2, direct complete migration, guarded activation, and recovery.
- [ ] Update `refresh-portfolio-catalog.md` with discovery-version checks, exclusions, inaccessible-member blockers, and prior-cache disclosure.
- [ ] Update the planning-workflows README and consumer repo template with v2 authoring and activation routes.
- [ ] Update `operation-routing-and-dispatch.md` route cards for multi-root migration and version-complete refresh.
- [ ] Copy each changed producer source to its exact `src/codeheart_operating_kit/resources/` mirror.
- [ ] Add source/mirror parity assertions to `tests/test_packaging_resources.py` for every changed path.
- [ ] Add a low-context nested-domain route probe to `tests/test_routing.py` with owner, route, ambiguity, approval, and stop assertions.
- [ ] Run `python3 -m pytest tests/test_packaging_resources.py tests/test_routing.py tests/test_sync_check.py`.
- [ ] Run `python3 scripts/validate-markdown-headers.py` and `python3 scripts/validate-public-core.py`.

### G) Implementation Notes

Affected runbooks remain agent-facing. Each needs compact intention, success, judgment, stop,
source/input/precondition, lane, approval, evidence, validation, and recovery sections. Generic Git
availability routes through tooling readiness; provider auth stays portfolio-route-owned.

Recipe maturity remains L1 orchestration over L3 Go commands. Do not add an L2 migration script.
Producer component files are authority; mirrors are copied only after source review.

Rollback boundary: revert source and mirror bytes together. A parity mismatch blocks the epic
checkpoint and release-readiness evidence.

### H) Open Questions

- OQ-01 affects only future v1-removal wording; this epic documents retained compatibility.

## 3.7 EP-07 - Comprehensive Validation And Cross-Platform Proof

### A) Epic ID, Title, And Outcome

`EP-07 - Comprehensive Validation And Cross-Platform Proof`

Outcome: unit, integration, CLI, schema, fixture, remote, portfolio, performance, packaging,
routing, macOS, Linux, and real-Windows evidence proves the complete discovery-v2 capability and
its compatibility/safety boundaries.

### B) Scope

- Consolidate the discovery acceptance/test matrix into automated fixtures.
- Add scale benchmark and bounded-process/memory assertions.
- Exercise platform-specific case, separator, Unicode, symlink/reparse, gitlink, nested-repo, and
  rename behavior.
- Run full Go/Python/schema/Markdown/public-core/package/resource suites.
- Add one Ubuntu semantic-validation job while retaining the existing macOS and Windows validation
  jobs and their supported release assets.

### C) Files Touched

- `internal/plancatalog/plancatalog_test.go`
- `internal/plancatalog/ep02_test.go`
- `internal/plancatalog/classifier_benchmark_test.go` (new)
- `internal/portfolio/portfolio_test.go`
- `internal/commands/commands_test.go`
- `tests/fixtures/plans/multi-root-repository/` (new)
- `tests/fixtures/portfolio/`
- `tests/test_json_schemas.py`
- `tests/test_packaging_resources.py`
- `tests/test_routing.py`
- `tests/test_sync_check.py`
- `.github/workflows/validate.yml`

### D) Acceptance Criteria And Size

- Size: `XL`.
- Every discovery acceptance-criteria and test-matrix row maps to a named automated test.
- A synthetic 100,000-path/10,000-Markdown benchmark uses a constant bounded Git-process count,
  deterministic output, and memory bounded independently of aggregate blob bytes.
- Real Windows and macOS jobs agree on reserved case, slash normalization, portability collisions,
  symlink/reparse/gitlink containment, and rename semantics.
- Existing v1, mixed baseline, transaction, remote transport, redaction, cache, and installer tests
  remain green.
- No test relies on private HQ content or network mutation.

### E) Dependencies And Critical-Path Notes

Depends on EP-01 through EP-06. This is the final implementation validation gate; failures return
to the owning epic and require its checkpoint to be corrected before EP-08.

### F) Tasks Checklist

- [ ] Map each discovery acceptance criterion to named Go, Python, schema, routing, packaging, benchmark, macOS, Linux, and Windows tests.
- [ ] Add filename-only, metadata-only, malformed, mismatch, duplicate, misplaced, excluded, prospective-blocked, unowned, unsafe, and preview fixtures.
- [ ] Add root docs, deeply nested docs, repeated docs segments, case variants, Unicode collisions, symlinks, gitlinks, nested repositories, and escape fixtures.
- [ ] Add metadata-qualified zero-child, one-child, many-child, router README, wrong-filename family, malformed family, v1-family migration, and unchanged-README branch fixtures.
- [ ] Add direct legacy-to-canonical, optional mixed, frozen-register, original-cutover, source-hash, candidate-digest, dirty, branch-owned, rollback, and recovery tests.
- [ ] Add remote v2 member, v1 member, inaccessible member, branch policy, rename, cache preservation, inert content, and redaction tests.
- [ ] Implement the 100,000-path and 10,000-Markdown benchmark in `classifier_benchmark_test.go` with deterministic process and memory assertions.
- [ ] Add an `ubuntu-latest` semantic-validation job to `.github/workflows/validate.yml` that runs `go test ./...`, the focused plan-catalog/portfolio/command suites, Python schema/routing/resource tests, and JSON-schema/Markdown/public-core validators without Linux release-asset generation and distribution advertising.
- [ ] Run `go test ./...`.
- [ ] Run `go test ./internal/plancatalog ./internal/portfolio ./internal/commands ./internal/cli`.
- [ ] Run `python3 -m pytest`.
- [ ] Run `python3 scripts/validate-json-schemas.py`, `python3 scripts/validate-markdown-headers.py`, and `python3 scripts/validate-public-core.py`.
- [ ] Run the focused benchmark twice and compare deterministic candidate/result digests.
- [ ] Require the Linux, macOS, and real-Windows GitHub validation jobs to pass on the same checkpoint.

### G) Implementation Notes

Benchmark wall-clock numbers are recorded but not frozen to one developer machine. Enforce
algorithmic evidence: one index/tree listing per ref, bounded batch readers, bounded workers,
constant Git-process count, and sorted output. Hostile Markdown remains inert bytes.

Rollback boundary: tests and fixtures may be corrected with their owning behavior checkpoint.
Never weaken expected errors, completeness, containment, or cache preservation merely to pass a
platform job.

### H) Open Questions

- OQ-02 refinement is deferred beyond this implementation; EP-07 tests the frozen Section 2.2 list.

## 3.8 EP-08 - Release Readiness And Consumer Handoff Evidence

### A) Epic ID, Title, And Outcome

`EP-08 - Release Readiness And Consumer Handoff Evidence`

Outcome: public-safe impact, migration, release-note, reproducibility, and consumer-validation
evidence is ready for a separately authorized Kit release, while no tag, release, consumer upgrade,
or HQ change occurs.

### B) Scope

- Create a plan-scoped consumer-impact record with affected paths, validation, notes, and consumer
  action.
- Create version-neutral migration and release-note input from implemented behavior.
- Build release assets twice in isolated temporary directories and record byte/digest comparison as
  bounded readiness evidence.
- Record the exact release stop boundary and remaining unsigned distribution risk from the release
  runbook without publishing.
- Create an HQ handoff sequence that starts only after an approved release and keeps all HQ writes
  outside this producer plan.

### C) Files Touched

- `docs/repo/plans/ownership-aware-plan-catalog-discovery/attachments/consumer-impact-record.md` (new)
- `docs/repo/plans/ownership-aware-plan-catalog-discovery/attachments/discovery-v2-consumer-migration.md` (new)
- `docs/repo/plans/ownership-aware-plan-catalog-discovery/attachments/release-note-input.md` (new)
- `docs/repo/plans/ownership-aware-plan-catalog-discovery/attachments/release-readiness-evidence.md` (new)
- sibling execution log created only after later activation

### D) Acceptance Criteria And Size

- Size: `M`.
- Impact evidence covers instruction, validator, consumer migration, safety, and placement review.
- Migration notes contain existing-repository v1 hold, prospective inventory, exclusions review,
  semantic migration, direct complete canonical route, optional mixed route, activation, remote
  completeness, rollback, and cache limits.
- Release-note input describes feature, compatibility, security, validation, and consumer action
  without selecting a release version.
- Two isolated release builds are byte-identical and catalog-to-binary verification passes on the
  implementation checkpoint.
- Release evidence states that version selection, manifest/version changes, tag, push, release,
  assets, and publication require separate authority.
- HQ handoff keeps portfolio-v2 migration paused until release and requires prospective local plus
  remote review before direct legacy-to-canonical activation.

### E) Dependencies And Critical-Path Notes

Depends on EP-07. This epic ends implementation readiness and hands off to release review. It does
not run the public release procedure beyond non-publishing reproducibility checks.

### F) Tasks Checklist

- [ ] Create `attachments/consumer-impact-record.md` with impact classes, affected paths, validation, release-note need, migration need, and consumer action.
- [ ] Create `attachments/discovery-v2-consumer-migration.md` with v1 hold, prospective inventory, exclusions review, semantic migration, direct canonical route, optional mixed route, activation, remote completeness, and recovery.
- [ ] Create `attachments/release-note-input.md` with version-neutral included, compatibility, impact, security, validation, and rollout sections.
- [ ] Build supported release assets twice in separate temporary directories with `scripts/build-release-assets.py`.
- [ ] Compare release packs, manifests, catalogs, installers, checksums, and embedded discovery-v2 resources byte-for-byte.
- [ ] Create `attachments/release-readiness-evidence.md` with source commit, commands, platform results, digests, residual risks, and explicit publication stop.
- [ ] Record the Codeheart-HQ post-release validation sequence without editing HQ files and without resuming its paused migration.
- [ ] Run `go test ./...` and `python3 -m pytest` on the exact readiness commit.
- [ ] Run every schema, Markdown, public-core, manifest, packaging, routing, and release-candidate validator required by `change-operating-kit.md`.
- [ ] Verify Git status contains only reviewed producer implementation/evidence paths and no consumer repository changes.
- [ ] Stop before version selection, manifest/version bump, tag, release publication, consumer upgrade, and HQ activation.

### G) Implementation Notes

Use temporary output directories and preserve non-secret digest summaries only. Do not commit build
artifacts. Release-note input is intentionally version-neutral so OQ-01 and later release timing do
not block implementation.

HQ sequence after release: install without activation; run prospective v2 inventory with remote
overlays; review every included/excluded/prospective-blocked/ambiguous candidate; migrate all formal
records; prove canonical local and compatible remote evidence; activate discovery v2 and canonical
directly when zero filename-only records remain; preserve the old register as historical evidence.

Rollback boundary: readiness evidence can be regenerated from the exact implementation commit.
Failed reproducibility, platform, digest, or completeness evidence blocks release handoff and does
not authorize cache deletion, version changes, publication, or consumer work.

### H) Open Questions

- OQ-01 remains non-blocking and belongs to a later v1-removal release.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Discovery-v1 removal is deferred until at least one compatibility release plus adoption evidence
  and separate approval resolve OQ-01.
- An `owned_roots` inclusion override is deferred indefinitely; it requires a concrete genuine-plan
  placement that cannot be corrected and a separate discovery/review.
- Release version selection, component/profile/manifest version bumps, public `release-notes.md`
  integration, tag creation, asset publication, and post-publication verification are deferred to a
  separately authorized release task using EP-08 evidence.
- Codeheart-HQ installation, prospective scan, semantic migration, discovery-v2 activation,
  canonical cutover, portfolio refresh, commit, and push remain consumer-owned work after release.
- Deletion tombstones and proactive portfolio services remain outside discovery v2.

## 4.2 Future Considerations

- Reassess the frozen conventional ambiguity list only in a later reviewed contract change backed by
  public-safe adoption findings, positive and negative fixtures, prospective inventory impact, and
  release notes without weakening the exclusions-only contract.
- Evaluate v1 removal only after portfolio homes can prove every required member supplies v2
  evidence.
- Consider a fuller CommonMark parser only when real documents exceed the bounded structural
  scanner and fixtures prove compatibility.
- Consider compatible discovery-version ranges for a future v3 only after more than one active
  catalog version exists.
- Keep plan identity in metadata, kind in filename/family qualification, lifecycle in the header,
  and path/ref/hash in observations through future versions.

# Revision Notes

- 2026-08-06: Created the draft, manual-review-ready implementation plan from committed discovery
  `codeheart-operating-kit.discovery.ownership-aware-plan-catalog-discovery` and current producer
  source evidence. Planned shared Git-backed classification, arbitrary-depth docs eligibility,
  exclusions-only ownership, metadata-qualified families, prospective activation, versioned
  evidence, guarded migration, optional mixed compatibility, remote completeness, managed parity,
  comprehensive validation, bounded release readiness, and post-release HQ handoff. No execution,
  release, or consumer authority was granted.
- 2026-08-06: Completed a main-thread planning review against the implementation-planning and
  planning-document-review runbooks; checked capability coverage, file/symbol surfaces,
  non-blocking OQ-01/OQ-02 handling, per-epic validation/recovery/checkpoints, draft-only authority,
  and the release/HQ stop boundaries.
- 2026-08-06: Applied review defaults: distinguished local nested-repository evidence from remote
  gitlinks/outer-commit blobs, removed directory-shape validity from v2 families, made EP-05 linear
  after EP-04, required an Ubuntu semantic-validation job without a Linux release asset, and froze
  the first discovery-v2 conventional ambiguity segment set while deferring later refinement.
- 2026-08-06: Activated the plan for goal-style producer implementation on
  `codex/operating-kit-multi-root-plan-catalog`; added the sibling execution log and retained the
  pull-request, merge, release, tag, consumer, HQ, destructive-Git, history-rewrite, and force-push
  stop boundaries.
