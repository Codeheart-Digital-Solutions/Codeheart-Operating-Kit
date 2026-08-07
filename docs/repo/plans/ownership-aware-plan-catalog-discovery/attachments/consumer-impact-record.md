Last updated: 2026-08-07T01:39:50Z (UTC)

# Repository-Wide Plan Catalog Discovery Consumer Impact Record

Implementation source: `f9d915256cb3dec1143c78a2d4dc79eb002c0ae2`

## Classification

| Impact class | Applicability | Reason |
| --- | --- | --- |
| `instruction-only change` | yes, for the managed planning guidance subset | Managed discovery, authoring, migration, register, portfolio, and review guidance now explains repository-wide discovery-v2 behavior. |
| `validator-only change` | yes, for validation and prospective-read behavior | Plan validation, inventory, list, migration, portfolio completeness, and remote-overlay checks share a stricter versioned classifier. |
| `consumer migration required` | yes, for repositories choosing discovery v2 | Activation is explicit and requires prospective inventory, semantic review, candidate migration, and a reviewed config change. A Kit upgrade alone performs none of those steps. |
| `security or safety policy change` | yes | Authority is restricted to tracked regular Git blobs beneath exact lowercase `docs` segments, with hard-unowned boundaries, inert Markdown parsing, and explicit exclusion evidence. |
| `breaking placement-contract change` | reviewed; not classified as breaking | No managed or consumer-owned path moves. Discovery v2 expands catalog authority only after explicit activation; discovery v1 remains the active contract for existing repositories after upgrade. |
| `component addition` | no | No component or profile-selected component was added. |
| `backwards-compatible scaffold addition` | no | No consumer-owned scaffold or generated path was added. |

## Affected Surfaces

- Plan-catalog models, config loading, Git index/tree enumeration, candidate classification,
  validation, structured views, inventory, migration, and CLI routing under `internal/` and `cmd/`.
- Portfolio remote-source, branch-overlay, completeness, and cache behavior under
  `internal/portfolio/`.
- Versioned plan view, inventory, migration-ledger, and portfolio catalog schemas under
  `schemas/`.
- Managed planning-workflow doctrine and its packaged resource mirror under
  `components/planning-workflows/` and `src/codeheart_operating_kit/resources/`.
- Fresh config defaults, release content identity, tests, fixtures, and the validation workflow.
- Existing consumer config, plans, register, and cached evidence are read and preserved according
  to their active v1 contracts until separately reviewed discovery-v2 activation.

## Compatibility And Consumer Action

Existing repositories with an absent or explicit discovery version `1` retain the fixed
`docs/repo/plans/` catalog after upgrading. They receive prospective discovery-v2 reads but no
automatic candidate expansion, metadata writes, mode change, register rewrite, second cutover
revision, or cache promotion.

A repository that adopts discovery v2 must follow
`discovery-v2-consumer-migration.md`: inventory the tracked candidate set, review ownership and
exclusions, migrate every formal candidate or explicitly choose bounded mixed grandfathering,
validate local and required remote evidence, and activate v2 in one reviewed config change. When
every candidate receives valid metadata, the normal route finishes in canonical mode without a
frozen register or lasting mixed phase.

Freshly initialized repositories default to discovery v2. All ordinary tracked Markdown paths
containing an exact lowercase `docs` segment are owned by default unless a hard-unowned boundary
or explicit `excluded_roots` entry applies.

## Release And Migration Requirements

- Consumer-facing release notes: required.
- Explicit migration/adoption notes: required.
- Versioned machine-contract disclosure: required for view, inventory, ledger, and portfolio
  completeness schema v2; plan metadata remains schema v1.
- Existing v1 cache disclosure: required; v1 evidence remains historical and cannot be called
  v2-complete.
- Signing/audience disclosure: required under the repository's existing unsigned
  internal/prototype boundary unless the release procedure establishes stronger signing.

## Validation Evidence

- Full Go and Python suites, JSON Schema, Markdown, public-core, release-manifest, routing,
  content-identity, package/resource parity, and plan-catalog validation passed on the
  implementation source.
- The synthetic 100,000-path/10,000-Markdown benchmark retained bounded memory, one deterministic
  digest across repeated runs, and no Git subprocess inside the pure classifier.
- GitHub Actions run `31137575877` passed Ubuntu semantic validation, macOS validation, and real
  Windows validation on the exact implementation source.
- Two isolated release-readiness builds were byte-identical; detailed artifact and provenance
  evidence is recorded in `release-readiness-evidence.md`.

## Known Residual Risk

Discovery v2 intentionally makes newly tracked plan-like files beneath any owned `docs` segment
visible. Misleading fixture, vendor, generated, dependency, or build-like locations therefore
become prospective blockers until the content is moved or the non-authoritative root is explicitly
excluded. This is visible safety behavior, not silent omission.

The eventual release that removes discovery v1 and any evidence-based refinement of conventional
ambiguity segments remain separate non-blocking product decisions.
