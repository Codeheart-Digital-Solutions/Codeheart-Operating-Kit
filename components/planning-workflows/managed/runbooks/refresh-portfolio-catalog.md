Last updated: 2026-08-09T15:24:37Z (UTC)

# Refresh Portfolio Catalog

Audience: agent-facing

Intent:
Refresh source-derived portfolio facts before strategic analysis and disclose freshness and
completeness without treating local-only work or an older cache as current remote truth.

Success:
A complete compatible schema-v3/discovery-v2 current attempt atomically updates the local cache
with separate scan, mixed-coverage, and canonical-readiness facts. An incomplete attempt preserves
the prior complete compatible cache and all historical schema-v1/v2 evidence and is reported with
its current blockers and freshness limit.

Agent judgment boundary:
Choose text or JSON output and summarize non-secret facts. Do not broaden discovery scope, modify
member repositories, rewrite the strategic overlay, weaken membership, or infer completeness.

Stop boundary:
Stop current-state analysis on wrong role, malformed config, auth/scope ambiguity, required-member
failure, incomplete scan, or evidence that the result omits a source the user considers required.

## Source, Inputs, Preconditions, And Lane

Source of truth: current user scope, local repository instructions,
`../reference/portfolio-coordination-format.md`, `.codeheart/kit.config.yaml`, machine-local
sources, the current scan result, and the last complete cache only as labeled historical evidence.

Inputs: coordination-home path, requested output format, and the analysis scope that will consume
the result.

Preconditions: configured `coordination-home` role; readable Kit config/lock/marker; Git available;
authenticated `gh` session when committed GitHub scopes exist; network and repository access; no
expectation that unpushed work is globally visible. Route missing tools to the tooling-readiness
runbook.

Execution lane: `codeheart-operating-kit portfolio scan`. This is authenticated external read plus
scanner-owned local cache/mirror refresh. It is not a member-repository write.

Recipe maturity: L1 read recipe over the L3 scan command. Fresh-agent review proves refresh-first
selection and failure disclosure; command/adapter tests prove pagination, completeness, cache, and
blocker behavior. Do not add an L2 polling or wrapper script in this version.

Approval boundary: a user request for current portfolio analysis or refresh authorizes the
configured read-only scopes and rebuildable local cache when repository policy allows them. It
does not authorize adding scopes, sign-in, member writes, strategic-overlay edits, commits, pushes,
PRs, merges, or releases.

## Procedure

1. Confirm the target is the configured coordination home and inspect source counts without
   printing credentials or unnecessary machine paths.
2. If GitHub sources exist, run the route's normal `gh auth status` preflight. If sign-in is
   required, use visible-terminal handoff and wait for the user; never request a token in chat.
3. Run:

   ```sh
   codeheart-operating-kit portfolio scan --format json <coordination-home-path>
   ```

   Use text for direct human orientation; use JSON when analysis or evidence needs exact fields.
4. Check schema version 3, discovery version 2, `catalog.complete`,
   `mixed_coverage_complete`, `canonical_ready`, start/completion time, `last_complete_scan_at`,
   errors, source/member/candidate/plan-observation/compatibility-observation/stale counts, API
   calls, duration, and cache update/preservation.
5. Confirm every required member's default branch activates discovery v2 and every candidate is
   either canonical remote evidence or one exact config-v2/ledger-v3 compatibility observation.
   Review structured excluded, prospective-blocked, hard-unowned, malformed, and filename-only
   observations rather than counting them as valid plans.
6. For each compatibility observation, require `coverage_disposition: mixed-grandfathered`,
   `incremental_migration_required: true`, one visible legacy row, the bound migration evidence,
   current cutover/register/owner/ref proof, and `canonical_ready: false`. Missing, stale, or moved
   evidence makes the member incomplete; it does not silently remove the row.
7. Confirm the result says default-branch config controls branch classification and local heads,
   worktree changes, untracked files, and unpushed commits are omitted.
8. If complete, use the current result for analysis and consult the repo-owned strategic overlay
   separately. Distinguish source facts from coordination-owned interpretation. A complete mixed
   member may be `mixed_coverage_complete` while not `canonical_ready`; disclose the pending
   incremental migration rather than calling it canonical.
9. If incomplete, report current error codes including incompatible schema/config members and that
   the previous complete compatible cache was preserved when applicable. Keep schema-v1/v2 caches
   labeled historical. Do not silently substitute any cache as current.

## Stop Conditions

Stop before current schema-v3 portfolio conclusions when the command exits nonzero, `complete` is
false, a required member is v1/inaccessible/invalid, claimed mixed coverage is incomplete,
canonical readiness is falsely claimed while a deferred row remains, a config-v2/ledger-v3 binding
or remote overlay is missing/stale, membership evidence is unavailable, pagination/truncation/
retry/rate-limit evidence is incomplete, auth or permission fails, merge-base or freshness evidence
is unavailable, or a required source is missing from configured scope.

Do not edit config to make a scan pass. Return to `configure-portfolio-coordination.md` for a
user-approved scope/identity change.

## Evidence And Validation

Record: home ID, attempt start/completion times, `complete`, `mixed_coverage_complete`,
`canonical_ready`, last-complete time, compatibility/factual observation counts, pending
incremental follow-ups, metrics, cache updated or previous preserved, error codes with retryability,
overlay status/digest/source-identity SHA-256, stale/conflict counts, and the local-only omission notice. Redact credentials
and avoid copying absolute local paths unless needed to resolve a blocker.

Validation succeeds when a complete schema-v3 result validates against discovery v2, enrolled
members satisfy exact default-branch evidence, branch overlays use that same policy and candidate
set, every deferred compatibility row stays visible with its verified binding, readiness states are
not conflated, the cache is updated atomically, and strategic-overlay bytes are unchanged. For a
failed attempt, validation succeeds only as failure handling: the previous complete compatible
cache and historical schema-v1/v2 caches are unchanged and the limitation is disclosed.

## Recovery

Retry transient or rate-limit failures only within the command's bounded policy, then stop. Repair
auth through visible-terminal handoff, permissions through the owning repository/provider route,
and malformed member declarations in the member repository under separate authority. Preserve
cache and mirror evidence; do not delete it merely to make a scan appear fresh.
