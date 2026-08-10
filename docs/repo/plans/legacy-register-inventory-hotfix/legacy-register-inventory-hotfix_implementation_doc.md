Last updated: 2026-08-10T06:32:50Z (UTC)
Created: 2026-08-10
Status: active
Execution log: legacy-register-inventory-hotfix_execution_log.md

# Legacy Register Inventory Schema-v3 Projection Hotfix Implementation Plan

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.legacy-register-inventory-hotfix
  kind: implementation
  purpose: Repair schema-v3 inventory projection for recognized frozen legacy-register evidence while preserving strict blockers and compatibility.
  first_cataloged: 2026-08-10T05:25:11Z
  catalog_metadata_updated: 2026-08-10T05:25:11Z
  products:
    - codeheart-operating-kit
  capabilities:
    - semantic-plan-catalog
  strategic_themes:
    - ownership-aware-discovery
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.branch-ownership-reconciliation
```
<!-- END CODEHEART PLAN METADATA -->

## Outcome

Ship v0.1.28 as a focused validator-only compatibility hotfix. `plans inventory` must parse
recognized frozen legacy-register evidence into a raw typed observation and then project a strict,
self-validating schema-v3 wire record. Canonical lifecycle stays closed; recognized historical
status is retained separately; `Last updated` and `Completed` remain separate typed dates; and
unsupported, malformed, duplicated, or ambiguous historical forms produce structured blockers.

The authorized workflow includes implementation, validation, a focused ready pull request,
expected-head guarded merge, annotated v0.1.28 tag, public release assets and sidecars, and live
verification. It excludes every consumer repository and leaves frozen registers byte-for-byte
unchanged.

## Authority And Impact

- Base: exact `origin/main` and released v0.1.27 merge
  `70a8e174fe8767b09eb6e0ccbea64f6fad5e9a3f`.
- Branch: `codex/legacy-register-inventory-v0.1.28-hotfix`.
- Producer authority: `internal/`, `schemas/`, `tests/`, release surfaces, and this plan bundle.
- Consumer impact: `validator-only change`; release notes are required; no migration, sync,
  repair, catalog activation, or consumer write is required by the upgrade itself.
- Compatibility: preserve deterministic digests, inventory v1/v2 behavior, schema-v3 branch and
  activation evidence, local/canonical-target/remote-aware parity, and all v0.1.27 ownership,
  checkpoint, and transaction invariants.
- Public-core boundary: fixtures and evidence are synthetic and generic. No private consumer path,
  ref, identifier, status text, or count is copied into the repository.

## Non-goals

- Do not broaden canonical lifecycle or another closed schema field to arbitrary strings.
- Do not invent a canonical lifecycle mapping for pre-canonical status labels.
- Do not rewrite, normalize, or migrate any frozen consumer register.
- Do not inspect or operate HQ, Foundry, AWS, or another consumer repository.
- Do not change unrelated branch reconciliation, migration, activation, or release authority.

## EP-01 - Typed Legacy Parsing And Strict Projection

Outcome: raw historical evidence and canonical schema-v3 output have distinct representations and
fail closed at their boundary.

- [x] Introduce a typed raw legacy-register observation separate from the canonical inventory wire record.
- [x] Recognize the supported historical status vocabulary without mapping it into canonical lifecycle.
- [x] Parse and retain separate `Last updated` and `Completed` dates.
- [x] Project only unambiguous canonical fields plus closed typed legacy evidence.
- [x] Emit stable structured blockers for unsupported labels, malformed or combined values, duplicates, and ambiguity.
- [x] Preserve deterministic hashing and zero-write inventory behavior.

## EP-02 - Regression And Compatibility Proof

Outcome: sanitized tests prove accepted and rejected evidence across every supported inventory
mode without weakening historical or safety contracts.

- [x] Add generic positive fixtures for recognized legacy statuses, canonical statuses, and separate dates.
- [x] Add generic negative fixtures for combined malformed forms, duplicate fields, unsupported labels, and ambiguous evidence.
- [x] Prove repeated inventory bytes/digests are deterministic.
- [x] Prove local, prospective canonical-target, and remote-aware projections are equivalent.
- [x] Prove inventory v1/v2 compatibility and all v0.1.27 branch/checkpoint/activation invariants remain passing.
- [x] Prove inventory writes only its explicit artifact and never changes frozen register or other consumer bytes.
- [x] Run focused tests, full Go/race/vet, Python/schema/resource/routing/release validation, and independent review; resolve every High and Medium finding.

## EP-03 - v0.1.28 Release And Live Verification

Outcome: the exact reviewed merge is reproducibly packaged, published, and verified as v0.1.28.

- [x] Update every authoritative version, manifest, installer, compatibility, catalog, and release-note surface.
- [x] Build macOS universal and Windows x64 assets twice and require byte equality.
- [x] Generate and verify SHA-256 sidecars and the complete catalog-to-binary/content identity chain.
- [ ] Verify isolated installer/upgrade compatibility and the unsigned internal/prototype boundary.
- [ ] Push a focused ready PR; verify base, head, changed paths, CI, and reviews.
- [ ] Merge with the exact expected-head guard and verify merged `main`.
- [ ] Create an annotated v0.1.28 tag at exact merged `main` and publish every required asset.
- [ ] Live-verify metadata, asset names, sidecars, downloads, installer, CLI version, and evidence chain.

## Acceptance Contract

- Legitimate supported frozen legacy evidence produces schema-v3 inventory that validates against
  `schemas/plan-inventory.schema.json` in all three inventory modes.
- Historical status evidence is retained in a closed typed field and never occupies canonical
  `lifecycle` unless the source lifecycle is itself canonical and unambiguous.
- Separate valid dates remain separate; combined, duplicate, unsupported, malformed, or ambiguous
  evidence blocks with stable structured codes.
- Repeated equivalent observations produce identical wire bytes and digests.
- v1/v2 compatibility, remote overlay parity, no-consumer-write guarantees, and all v0.1.27 safety
  tests pass unchanged.
- v0.1.28 is merged, annotated, published with reproducible verified assets, and live installation
  reports the exact released version.
