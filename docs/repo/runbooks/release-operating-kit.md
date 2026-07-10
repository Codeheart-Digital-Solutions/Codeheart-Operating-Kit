Last updated: 2026-07-09T23:45:00Z (UTC)

# Release Operating Kit

Use this runbook before publishing a Codeheart Operating Kit release.

## Procedure

1. Read `AGENTS.md`.
2. Read `README.md`.
3. Read `docs/repo/reference/consumer-impact-classification.md`.
4. Confirm the intended version and consumer-impact classification.
5. Confirm release notes cover consumer-facing behavior and migration.
6. Run public-core, Markdown, JSON Schema, content-identity, Go, Python compatibility, installer,
   and release-contract validation.
7. Build macOS universal and Windows x64 packs twice with `scripts/build-release-assets.py`; require
   byte-identical output from both builds.
8. Confirm each pack contains the expected binary plus `bootstrap.md`, `install.sh`, `install.ps1`,
   `release-notes.md`, `INSTALL.md`, `content-manifest.yaml`, `pack-manifest.json`, and
   `checksums.txt`, with no Python wheel or `*.dist-info` payload.
9. Verify the complete external catalog -> archive digest -> pack manifest -> payload checksums ->
   content identity -> binary digest and version chain.
10. Confirm the external release catalog is generated after the packs. Keep archive URLs and
    digests out of embedded `manifest.yaml`; add final public URLs only to the external catalog.
11. Run isolated macOS and real Windows fresh-install and upgrade dry-run/apply/failure paths.
    Require failed verification, replacement, reconciliation, or post-check to preserve or restore
    the prior installation.
12. Generate and verify sidecar SHA-256 checksums for every public asset.
13. Treat locally generated or unsigned assets as release-candidate evidence only. Before broad
    public distribution, require signing/notarization or explicitly approve and record the
    unsigned internal/prototype boundary.
14. Stop here unless tag creation and publication are separately and explicitly authorized.
15. If publication is authorized, verify the target commit is the validated commit, create the
    tag, publish the packs, catalog, installers, notes, and checksums, then record URLs, digests,
    platform evidence, signing state, and residual risk.

## Stop Conditions

Stop before publishing when public-core hygiene fails; schemas or migrations are incompatible;
release assets are not reproducible; any catalog-to-binary digest is missing or inconsistent;
packs contain a Python payload; macOS or Windows install/upgrade validation fails; rollback does
not preserve the prior installation; release notes omit consumer-impacting changes; staged assets
are confused with live public assets; signing/notarization readiness is unresolved for the
intended audience; publication is not explicitly authorized; or the release target commit differs
from the validated commit.
