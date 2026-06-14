Last updated: 2026-06-14T01:01:25Z (UTC)

# Codeheart Operating Kit v0.1.0 Release Notes

`v0.1.0` is the first public Codeheart Operating Kit release.

## Included

- Public bootstrap instructions in `bootstrap.md`.
- macOS user-level installer in `install.sh`.
- Windows user-level installer in `install.ps1`.
- macOS release asset `codeheart-operating-kit-0.1.0-macos.tar.gz`.
- Windows release asset `codeheart-operating-kit-0.1.0-windows.zip`.
- SHA-256 checksum files for both platform assets.
- Release manifest `manifest.yaml`.
- Managed Operating Kit docs packaged inside the CLI wheel and described by component manifests.
- `standard` profile for G1 repository and agent-memory scaffolding.
- `codeheart-operating-kit` commands: `onboard`, `inspect`, `init`, `sync`, `check`, and
  `update-check`.
- Native Codex capability status reporting for documents, spreadsheets, presentations, browser,
  and PDF support.

## Consumer Impact

- `component addition`: first public component set, profile, schemas, validators, and managed docs.
- `instruction-only change`: first public operating instructions, runbooks, and references.
- `validator-only change`: first public validation scripts and CI validation matrix.
- `backwards-compatible scaffold addition`: consumer repository docs and agent-memory scaffolds
  are created only when absent.
- `consumer migration required`: existing repositories adopting the kit must run onboarding or
  initialization to add `.codeheart/` metadata and the managed `AGENTS.md` block.
- `breaking placement-contract change`: G1 establishes `.codeheart/kit/` as the managed kit root
  and intentionally does not create a `docs/workspace/` surface.

## Validation

- GitHub Actions validates Linux CLI tests, public-core hygiene, Markdown timestamps, JSON schema
  structure, release-manifest content, native capability tests, release asset build, macOS
  installer smoke, and Windows installer smoke.
- Public release smoke validation is available through the manual `Validate` workflow with
  `release_version` set to `v0.1.0`.
- Installers verify SHA-256 checksums and fail closed on checksum mismatch.

## Deferred

- Signed release artifacts.
- Linux installer asset.
- Homebrew, winget, npm, GitHub Packages, and CodeArtifact publication.
- Specialized profiles beyond `standard`.
- Foundry module selection during onboarding.
