Last updated: 2026-07-04T22:00:01Z (UTC)

# Release Operating Kit

Use this runbook before publishing a Codeheart Operating Kit release.

## Procedure

1. Read `AGENTS.md`.
2. Read `README.md`.
3. Read `docs/repo/reference/consumer-impact-classification.md`.
4. Confirm the intended version.
5. Confirm release notes cover consumer-impacting changes.
6. Run public-core hygiene validation.
7. Run markdown timestamp validation.
8. Run schema and manifest validation when schemas and manifests exist.
9. Run Go CLI tests.
10. Run Python parity tests while the Python CLI remains the behavior oracle.
11. Run installer and release-asset tests.
12. Build binary release packs with `scripts/build-release-assets.py`.
13. Confirm the macOS pack is `codeheart-operating-kit-<version>-macos-universal.zip` and contains
    a universal `bin/codeheart-operating-kit` binary.
14. Confirm the Windows pack is `codeheart-operating-kit-<version>-windows-x64.zip` and contains
    `bin/codeheart-operating-kit.exe`.
15. Confirm release packs include `bootstrap.md`, `install.sh`, `install.ps1`, `release-notes.md`,
    `INSTALL.md`, `manifest.json`, and `checksums.txt`.
16. Confirm release packs do not include Python wheels, Python package payloads, or
    `*.dist-info` directories.
17. Generate and verify sidecar SHA-256 checksums for every public asset.
18. Verify installers fail closed on checksum mismatch.
19. Validate macOS install from the staged release pack.
20. Validate Windows install from the staged release pack through GitHub Actions.
21. Validate macOS and Windows installer failure paths for malformed archive, missing binary, and
    failed staged binary smoke validation.
22. Confirm staged unsigned or locally generated assets are treated as release-candidate evidence
    only and are not described as live public release assets.
23. Confirm the release is either signed/notarized for broad public distribution or explicitly
    recorded as an unsigned internal prototype boundary.
24. Update root `manifest.yaml` and packaged manifest release asset entries only after final public
    asset URLs and checksums are known.
25. Update `bootstrap.md` only with public release URLs that will exist in the GitHub release.
26. Create the Git tag from the validated commit.
27. Publish the GitHub release with release notes, manifests, installers, binary packs, and
    checksums.
28. Record release URLs, checksums, validation evidence, signing/notarization state, staged asset
    boundary, and residual risk.

## Stop Conditions

Stop before publishing when public-core hygiene fails, release assets are not reproducible,
checksums are missing, binary packs contain a Python wheel payload, installer validation fails,
release notes omit consumer-impacting changes, staged assets are confused with live public assets,
signing/notarization readiness is unresolved for the intended audience, or the release target
commit differs from the validated commit.
