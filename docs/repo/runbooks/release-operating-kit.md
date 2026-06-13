Last updated: 2026-06-13T22:47:44Z (UTC)

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
9. Run CLI tests when the CLI exists.
10. Build release assets.
11. Generate checksums.
12. Verify installers fail closed on checksum mismatch.
13. Validate macOS install from release assets.
14. Validate Windows install through GitHub Actions.
15. Create the Git tag from the validated commit.
16. Publish the GitHub release with release notes, manifests, assets, installers, and checksums.
17. Record release URLs, checksums, validation evidence, and residual risk.

## Stop Conditions

Stop before publishing when public-core hygiene fails, release assets are not reproducible,
checksums are missing, installer validation fails, release notes omit consumer-impacting changes,
or the release target commit differs from the validated commit.
