Last updated: 2026-08-07T00:26:02Z (UTC)

# Portfolio Coordination Format

Use this reference for coordination-home configuration, exact repository membership, source-
derived portfolio facts, local cache, and the separate strategic overlay. Machine contracts are
shipped as the Kit config, local-source, plan-catalog, and strategic-overlay schemas.

## Concepts And Ownership

`member`
: A Kit-enabled repository that declares one stable repository ID and the coordination home it
  belongs to. A member does not configure discovery sources.

`coordination-home`
: A Kit-enabled repository that declares its own repository ID and home ID, owns authorized
  discovery scopes, includes itself automatically as one member, stores the rebuildable factual
  catalog locally, and owns portfolio strategy in committed repository docs.

`candidate`
: A repository visible in an authorized source that does not satisfy exact membership. Candidate
  visibility never enrolls it and never authorizes plan ingestion.

There is no central repository list and no multiple-role array. Each repository has exactly one
role. The coordination home discovers candidates from configured scopes and enrolls only
repositories whose default-branch declarations match.

## Portfolio V2 Configuration

Shared non-secret configuration lives in `.codeheart/kit.config.yaml`.

Member:

```yaml
portfolio:
  schema_version: 2
  role: member
  member_repository_id: example-service
  coordination_home_id: example-portfolio
```

Coordination home:

```yaml
portfolio:
  schema_version: 2
  role: coordination-home
  member_repository_id: example-coordination-home
  coordination_home_id: example-portfolio
  discovery:
    sources:
      - kind: github-owner
        owner: Example-Organization
```

Both identifiers are lowercase portable identifiers: letters, digits, and interior hyphens. For a
coordination home, `member_repository_id` is its own stable repository identity. It is required so
the home contributes one normal self-member baseline without a second registration entry.

Repeated GitHub owners are committed scope authorization. They are not credentials. Use the
existing authenticated `gh` session at scan time; never copy tokens into config or logs.

## Machine-Local Sources

Local discovery roots live only at `.codeheart/local/portfolio/sources.yaml`:

```yaml
schema_version: 1
sources:
  - kind: local-git-root
    root: /absolute/machine-specific/path
```

This file is ignored, local-machine state. It must not be committed or copied to a coordination
home's strategic docs. A local root authorizes discovery and remote resolution only; developer
worktrees, indexes, local heads, and unpushed commits are not portfolio facts.

## Portfolio V1 Compatibility

Existing path-based portfolio-v1 configuration remains readable compatibility evidence. It may
contain coordination-home paths or register paths and may be paired with
`coordination-sync-pending.md`.

Do not create new v1 configuration or pending-sync scaffolds. Upgrade to v2 only through a
reviewed configuration change. Existing consumer-owned pending-sync files remain preserved and
readable; sync must not delete or overwrite them.

## Exact Membership Predicate

A discovered repository becomes a member only when all applicable evidence succeeds on its
default branch:

1. it is inside an explicitly authorized GitHub-owner or local-root source;
2. either `.codeheart/kit/README.md` proves a Kit installation or the Operating Kit producer's
   schema-valid root `manifest.yaml` proves the equivalent Kit source marker;
3. an installed Kit's `.codeheart/kit.lock.yaml` is readable and valid; a producer source marker
   does not require or authorize committing a consumer installation;
4. `.codeheart/kit.config.yaml` is readable and valid;
5. `portfolio.role` is valid for membership;
6. `member_repository_id` is stable and present;
7. `coordination_home_id` exactly matches the scanning home; and
8. required plan evidence can be read completely; and
9. discovery-v2 completeness is claimed only when the default branch activates discovery v2 and
   every selected-tree candidate has canonical v2 evidence.

The coordination home's own default branch follows the same evidence boundary and is included
once through self-membership. A feature branch cannot enroll a repository. A malformed declaration
is an error, not a candidate success. Missing or mismatched evidence produces a candidate or
exclusion reason without ingesting plan content.

The source-marker alternative exists only for the Operating Kit producer, whose tracked
`manifest.yaml` is validated against the Kit content-manifest schema. It keeps producer source
authority separate from `.codeheart/kit/` consumer installation state; it does not infer membership
from a repository name and does not relax role, repository-ID, home-ID, or authorized-scope checks.

## Discovery Sources And Read Boundary

Supported v1 adapters are:

- local Git roots, which find repositories and resolve their remotes; and
- GitHub owner scopes through authenticated, completely paginated `gh api` reads.

The scanner owns bare mirrors below `.codeheart/local/portfolio/git/`, refreshes and prunes remote
refs, and reads inert Git objects. It does not execute fetched content, hooks, helper transports,
or developer-worktree commands. It does not change developer refs or global Git/`gh` settings.

Provider pagination or truncation gaps, exhausted bounded retries, rate-limit exhaustion, auth or
permission failure, missing comparison evidence, mirror failure, or a required-member read failure
make the attempt incomplete.

## Plan Discovery Authority

The scanner uses the same discovery-v2 classifier and candidate contract as local `plans list`,
`validate`, `inventory`, and `migrate`. Its authoritative universe is regular or executable
Markdown blobs in the selected commit tree. Eligible paths contain an exact lowercase `docs`
directory segment at arbitrary depth; supported filename or genuine metadata discovers a
candidate; canonical remote validity requires both supported filename and valid matching metadata
for discovery/implementation records. Exact `README.md` requires explicit valid family metadata.

Default-branch `.codeheart/kit.config.yaml` controls discovery version and `excluded_roots` for the
repository and all branch observations. A feature branch cannot activate discovery v2, broaden
ownership, or relax exclusions. Excluded, conventionally ambiguous, hard-unowned, malformed, and
filename-only candidates remain structured observations so completeness failures are visible.

## Default Baselines And Branch Overlays

Each discovery-v2 member default branch contributes a complete canonical baseline. Every
accessible unmerged remote branch is evaluated independently of branch name, age, and pull-request
state.

An overlay contains only candidates whose plan evidence is added or materially changed from the
merge base. Semantic same-byte renames and mode-only changes do not add observations. Unchanged
README bytes do not become family authority when a sibling directory appears. Deletions do not
create plan tombstones in v2. Same-ID records and invalid candidate observations coexist with
repository, ref, commit, path, source mode, candidate signal, ownership disposition, and content
hash so conflicts and omissions remain visible.

Optional pull-request facts enrich a selected same-repository ref. They never select or suppress a
branch. Merged or deleted refs disappear after the next complete scan. Accessible observations
unchanged for more than 30 days are marked stale but retained.

## Cache Contract And Completeness

The discovery-v2 rebuildable factual cache uses schema version 2 and is:

```text
.codeheart/local/portfolio/catalog.json
```

It contains scan times, discovery policy/provenance, completeness, enrolled members, valid plan
observations, candidate observations, errors, and metrics including duration, source/member/
candidate/observation/stale/API-call counts and maximum concurrency.

Only a complete schema-v2/discovery-v2 scan atomically replaces the current complete cache. A v1,
inaccessible, invalid, or incomplete required member makes the current v2 attempt incomplete. An
incomplete attempt reports its current errors and preserves the previous complete v2 cache byte-
for-byte. Historical schema-v1 caches remain preserved and labeled but are never reported as
v2-complete. Therefore:

- a command's current result describes the current attempt;
- the cache is only the last complete result; and
- neither proves that unpushed local work is globally visible.

Never call portfolio analysis current after a failed or incomplete refresh. State the last known
complete time when available and disclose the current failure.

## Strategic Overlay

Coordination-owned interpretation lives at:

```text
docs/repo/portfolio/strategic-overlay.yaml
```

It is repo-owned, committed, and created only when absent. Schema v1 supports:

- cross-plan `families`;
- portfolio `themes`;
- cross-plan `relations`, including `overlaps`;
- ranked `priorities`; and
- timestamped narrative `analyses`.

The scanner validates but never rewrites overlay bytes. Scans do not push strategic interpretation
back into member repositories. Member metadata remains member-repository authority.

The sibling `docs/repo/portfolio/README.md` explains local ownership and should contain only
public-safe repository-specific context.

## Configuration And Scan Commands

Preview configuration before writing:

```sh
codeheart-operating-kit portfolio configure \
  --role member \
  --member-repository-id example-service \
  --coordination-home-id example-portfolio \
  --dry-run .
```

For a home, use `--role coordination-home` and repeat `--github-owner` and `--local-root` as
needed. Exactly one of `--dry-run` and `--yes` is required. Matching configuration is idempotent;
conflicting existing role or identities block rather than replace authority.

Refresh from a coordination home:

```sh
codeheart-operating-kit portfolio scan --format text .
codeheart-operating-kit portfolio scan --format json .
```

`--json` is a compatibility alias for JSON. A member cannot run a portfolio scan. A nonzero result
or `complete: false` means current completeness must not be claimed.

Remote-aware plan validation and prospective discovery-v2 inventory use the same scanner:

```sh
codeheart-operating-kit plans validate --target-discovery-version 2 \
  --target-catalog-mode canonical --remote-overlays .
codeheart-operating-kit plans inventory --target-discovery-version 2 \
  --target-catalog-mode canonical --remote-overlays --output <artifact> .
```

These commands may refresh scanner-owned local mirrors and cache only. They do not modify member
repositories.

Use `../runbooks/configure-portfolio-coordination.md` for setup and
`../runbooks/refresh-portfolio-catalog.md` before portfolio analysis.
