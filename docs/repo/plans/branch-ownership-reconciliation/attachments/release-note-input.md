Last updated: 2026-08-09T16:34:26Z (UTC)

# Branch Ownership Reconciliation Release Note Input

## Included

- Adds strict schema-v3 reviewed branch-touch ledgers with candidate-scoped dispositions for exact
  same-content non-owners, incorporated-history non-owners, active owners, deferred active owners,
  and blocking evidence.
- Adds deterministic Git-object proofs that bind evidence revision, policy, candidate set, logical
  ref, ref tip, unique merge base, complete relevant path transitions, blob identities, content
  hashes, and proof-specific attributes.
- Adds bounded incorporated-history matching against commits reachable from the reviewed target.
  Pull-request merge facts may corroborate the result but never replace offline Git proof.
- Adds a guarded `plans catalog-activate` step so reviewed evidence, migration actions, and config
  activation form explicit evidence, ledger, and activation checkpoints.
- Adds strict config-v2 bindings for mixed catalogs that retain one genuine active branch owner,
  plus structured pending incremental follow-up and compatibility observations.
- Adds remote-aware migration and activation checks, source-stable mirror identities, tracking-ref
  coalescing, cache/catalog v3 readiness, and structured JSON blocker output.
- Adds post-write authority validation and rollback so drift that occurs during a transaction does
  not survive a reported failure.

## Compatibility

Existing config-v1 and ledger-v1/v2 repositories keep their released behavior. Historical schemas
remain addressable. Operating Kit v0.1.25 and v0.1.26 reject the new ledger/config contracts before
writing, while the new CLI continues to accept their supported historical contracts.

A Kit upgrade alone does not activate discovery v2, change catalog mode, migrate plans, delete
branches, or publish a new portfolio cache. Direct canonical mode remains zero-gap and cannot use
active-owner deferral.

## Consumer Impact

Repositories previously blocked by same-byte or already-incorporated historical refs can retain
those refs and migrate after an explicit evidence review. Repositories with genuine branch-owned
active work may use one exact mixed-mode deferral per candidate when frozen cutover authority and
the mandatory incremental follow-up are present. Canonical readiness remains false until every
deferral is resolved.

## Security And Failure Behavior

Evidence is fail-closed. Ref movement, stale tracking data, ambiguous merge bases, partial patch
matches, path rename/delete changes, dirty owner bytes, unsupported Git, malformed or tampered
evidence, target changes, and overlay drift produce actionable blockers. Proof digests omit the
diagnostic Git version but retain the versioned algorithm and minimum admitted Git floor, allowing
identical object evidence to reproduce across supported operating systems.

## Rollout

Upgrade each consumer separately after the producer release is verified. Run prospective local or
remote-aware inventory, review every branch candidate, checkpoint the schema-v3 ledger directly
after its evidence revision, dry-run and apply the migration, run guarded catalog activation, then
validate mixed coverage and canonical readiness. When a deferred owner lands, reinventory and run
the reviewed incremental migration against the integrated bytes.

Consumer installation and plan-catalog activation are intentionally not part of this producer
release task.
