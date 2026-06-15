Last updated: 2026-06-15T10:43:39Z (UTC)

# Codeheart Operating Kit Release Notes

## v0.1.1 Release Notes

`v0.1.1` is an onboarding-contract hardening release for Codeheart Operating Kit.

### Included

- Public `README.md` `Start Here` prompt that keeps the first user prompt short and points to the
  latest public `bootstrap.md`.
- Hardened public `bootstrap.md` agent contract for language-first setup, explicit project-name
  selection, explicit target-folder selection, setup-plan preview, and write confirmation.
- Managed onboarding reference and runbook updates that make purpose optional context instead of a
  required setup branch for the `standard` profile.
- Neutral onboarding examples only: `Yourname-Automation`, `Companyname-Automation`,
  `Productname-Development`, `Teamname-Operations`, and `Existing-Project-Name`.
- Non-interactive onboarding policy that reserves `--yes` for explicit automation or repair flows
  where the user-owned decisions are already supplied.

### Consumer Impact

- `instruction-only change`: public bootstrap, public README, managed onboarding reference, and
  managed onboarding runbook now tell agents not to infer setup decisions.
- `validator-only change`: public-core validation now rejects legacy real-looking onboarding
  examples in public onboarding surfaces.
- `security or safety policy change`: non-interactive onboarding write behavior now fails closed
  when `--yes` is used without explicit target folder and project name decisions.
- CLI impact: `onboard` no longer supplies default target, project name, or purpose values;
  `init` can create a standard setup without purpose metadata.
- Schema impact: `setup_purpose` remains valid when present, but is no longer required for new
  `standard` profile configs.
- Installer impact: macOS and Windows installer defaults point to `v0.1.1` release assets.
- Known consumer action: none for installed consumers; existing `setup_purpose` metadata remains
  supported.

### Validation

- Local validation covers onboarding prompt order, explicit decision requirements, no-write
  behavior without `--yes`, config compatibility with and without `setup_purpose`, public-core
  hygiene, release manifest consistency, release asset build, and local macOS install from built
  assets.

### Deferred

- Purpose-specific installed profiles.
- Foundry module selection during onboarding.
- End-to-end external agent simulation.

## v0.1.0 Release Notes

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
