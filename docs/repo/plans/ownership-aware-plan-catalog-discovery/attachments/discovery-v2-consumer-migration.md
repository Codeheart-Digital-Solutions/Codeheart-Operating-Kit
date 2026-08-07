Last updated: 2026-08-07T01:39:50Z (UTC)

# Discovery V2 Consumer Migration

This note describes adoption after installing a release that contains repository-wide plan
discovery. Installation and semantic adoption are separate checkpoints.

## Compatibility Default

An existing repository stays on discovery v1 when
`plan_catalog_discovery_version` is absent or set to `1`. Upgrading the Kit does not scan new roots
as active authority, modify plans, change catalog mode, rewrite a register, create a cutover
revision, promote a cache, or activate discovery v2.

Discovery v2 considers every Git-tracked regular Markdown blob whose normalized repository-relative
path contains an exact lowercase directory segment named `docs`, at any depth. Candidate discovery
uses supported filename or genuine metadata. Canonical discovery/implementation validity requires
both a supported filename and valid metadata. An exact `README.md` becomes a family record only
with genuine valid family metadata.

## Reviewed Adoption Sequence

1. Record the installed Kit version, active discovery version and catalog mode, current Git
   revision/status, plan bytes, register bytes when present, config, lock, and portfolio role.
2. Keep discovery v1 active. Run active-contract validation before changing anything:

   ```sh
   codeheart-operating-kit plans validate <repository>
   ```

3. Create a plan-scoped evidence destination outside every actual or prospective plan target, then
   inventory prospective canonical v2 authority:

   ```sh
   codeheart-operating-kit plans inventory \
     --target-discovery-version 2 \
     --target-catalog-mode canonical \
     --remote-overlays \
     --output <inventory.json> \
     <repository>
   ```

   Omit `--remote-overlays` only when the repository has no configured portfolio overlay or that
   evidence is explicitly out of scope, and record the limitation. `--include-untracked` may be
   used for an authoring preview, but its results are non-authoritative and cannot enter hashes,
   migration writes, completeness, or remote evidence.
4. Review every owned, excluded, prospective-blocked, hard-unowned, filename-only, metadata-only,
   malformed, mismatch, duplicate, and family observation. Move genuine plans out of misleading
   fixture/vendor/generated/dependency/build-like locations. Use normalized `excluded_roots` only
   for tracked non-authority; exclusions stay visible in inventory. Do not add `owned_roots`.
5. Resolve active-branch ownership and dirty overlap. Build a schema-v2 semantic ledger bound to
   the inventory revision, repository identity, discovery policy, complete candidate set, source
   hashes, target paths, and reviewed metadata decisions.
6. Dry-run and review the exact mutation:

   ```sh
   codeheart-operating-kit plans migrate \
     --ledger <reviewed-ledger.yaml> \
     --dry-run \
     <repository>
   ```

7. Apply the approved ledger while active discovery remains v1, then re-run it until it is
   idempotent:

   ```sh
   codeheart-operating-kit plans migrate \
     --ledger <reviewed-ledger.yaml> \
     --yes \
     <repository>
   ```

8. Reinventory and require complete prospective evidence. Changed bytes, policy, candidates,
   source revision, active ownership, dirty overlap, or target collision invalidate stale ledger
   decisions and require review of the current evidence.
9. Run prospective canonical validation:

   ```sh
   codeheart-operating-kit plans validate \
     --target-discovery-version 2 \
     --target-catalog-mode canonical \
     --remote-overlays \
     <repository>
   ```

10. Only after complete evidence, make a separate reviewed config change that sets
    `plan_catalog_discovery_version: 2`, the reviewed
    `plan_catalog_ownership.excluded_roots`, and the selected catalog mode. Validate, commit, and
    publish through the repository's normal protected workflow.

## Normal Complete Route

When every formal candidate can receive a supported filename and valid metadata, keep config
unchanged during migration and activate discovery v2 directly in canonical mode after the final
prospective validation. The old register may remain historical evidence. Do not create a frozen
register, new cutover revision, or grandfathered gap merely to pass through mixed mode.

In shorthand, the route is:

`legacy v1 -> inventory/review -> migrate every candidate while v1 stays active -> validate -> canonical v2`

## Optional Deferred Route

Mixed mode is available only when specific unchanged filename-only records must be grandfathered.
Before migration, preserve the original legacy baseline and exact frozen-register/cutover proof,
set mixed mode while discovery remains v1, and validate that baseline. The reviewed ledger must
enumerate every deferred record and prove its pre-cutover bytes. Retain the original cutover
revision; do not create a second one for discovery v2.

Mixed is therefore a compatibility tool, not a required middle phase. New or materially changed
records still require valid metadata, and canonical completeness is unavailable while a required
portfolio member remains v1 or incomplete.

## Remote Evidence And Caches

Default-branch config controls remote discovery version and exclusions; a feature branch cannot
broaden authority. A v2-complete portfolio result requires compatible canonical v2 evidence from
every required member. Inaccessible, v1, malformed, or incomplete required members make the
current scan incomplete. Preserve and label the last complete cache; never present a v1 cache as
v2-complete.

## Recovery

Before activation, leave or restore the reviewed config on discovery v1 and preserve all inventory,
ledger, register, and migration evidence. Do not destructively remove valid metadata merely to
return to v1 views. After activation, a config-only return to v1 requires a reviewed repository
change and cannot be used to hide unresolved v2 findings. Transaction failure requires preserving
the marker and recovery evidence and following `check`/repair guidance.

Stop on ambiguous ownership or meaning, unsafe/non-regular authority, changed evidence, incomplete
coverage, broken mixed proof, failed validation, incomplete required remote evidence, or any write
outside the reviewed candidate set.

## Codeheart-HQ Sequence

The producer release is first installed into HQ with discovery v1 unchanged and plan bytes
preserved. HQ's separate semantic adoption then runs prospective local and remote inventory,
reviews every newly found/excluded/ambiguous candidate, migrates all genuine plans, validates
complete evidence, and activates canonical v2 directly if no filename-only record is deferred.
Installation of the release does not authorize or perform that semantic adoption.
