Last updated: 2026-08-10T06:32:50Z (UTC)

# v0.1.28 Release Readiness Evidence

## Validated Source

- Release-source commit: `534583827740575186f98f9b230f2be6da17e62c`.
- Branch: `codex/legacy-register-inventory-v0.1.28-hotfix`.
- Base: released v0.1.27 merge `70a8e174fe8767b09eb6e0ccbea64f6fad5e9a3f`.
- Consumer impact: `validator-only change`; no consumer migration, catalog activation, register
  rewrite, sync, or repair is required.
- Frozen producer register: unchanged from base.

The exact release source passed full Go, Go race, Go vet, 158 Python tests, public-core,
Markdown-timestamp, JSON-schema, release-manifest, packaging-resource, routing, installer, and
backward-compatibility validation. Independent review reported no remaining High or Medium
blocker.

## Reproducible Release Assets

`scripts/build-release-assets.py` ran twice in independent build directories with final public
v0.1.28 URLs. Each invocation also built every platform pack twice internally and required byte
equality. Build A and build B matched for both archives, both archive sidecars, and the external
catalog.

| Asset | SHA-256 |
| --- | --- |
| `codeheart-operating-kit-0.1.28-macos-universal.zip` | `49b2fe0f360840d28838bf05018d628037b000e30948e20d9041f8a06a63e8c1` |
| `codeheart-operating-kit-0.1.28-macos-universal.zip.sha256` | `2bfeb2f4390abe247be36608a4729c9143dc7a4f3cf88e0d3a6e050b6cd82070` |
| `codeheart-operating-kit-0.1.28-windows-x64.zip` | `c92708a2869f4ffb905819b1ad8ff1b30da26f1c68151714fe166183bb8c9350` |
| `codeheart-operating-kit-0.1.28-windows-x64.zip.sha256` | `9b7a980290485475f544812e0586fa325642940e2a41af4fba334f3fc013cf84` |
| `release-catalog-0.1.28.json` | `80d133862773138ad9e88e24ac47ebb30ce6ee415f690493cca780c0e7d156d6` |

The catalog binds pack-manifest SHA-256
`0b15c67c276220633fd9900db03363294d020f05a7641a25eee9cc22e03abdfe` for macOS and
`b669139ce6f628f0ec1bec790ba76188c2bda1318aae19bd7749556325346d4f` for Windows. Both packs
bind content-manifest SHA-256
`7085ecec9648ef424fa2d2f0cc3cf8133cac30c9aa9fbe503345fff7f4fee9ab`.

The macOS universal binary is Mach-O x86_64 plus arm64, reports v0.1.28, and has SHA-256
`95f7b3af62ccd5312c6be4cde31b49dd0ee150f300594b1d757933001220ae1c`. The Windows binary is
PE32+ console x86-64 and has SHA-256
`de6857129452a0359371c1e933e2357d5126475b1d003b724a490a7877761e0b`.

Every pack contains exactly the platform binary plus `bootstrap.md`, `install.sh`, `install.ps1`,
`release-notes.md`, `INSTALL.md`, `content-manifest.yaml`, `pack-manifest.json`, and
`checksums.txt`. There is no wheel, `*.dist-info`, or Python package payload. Zip ordering,
timestamps, entry types, modes, payload checksums, and identity paths passed builder verification.

## Catalog, Installer, And Upgrade Proof

The final-public catalog cannot complete its staged live-URL fetch until publication. An otherwise
identical local-URL catalog was therefore used for pre-publication installer and upgrade proof;
its archives are byte-identical to the final-public archives.

- A fresh isolated macOS install from the verified pack succeeded and reported
  `codeheart-operating-kit 0.1.28`.
- Reinstalling the same pack preserved the exact binary digest.
- An explicit checksum mismatch failed before replacement and preserved the prior runnable binary
  with the same digest and version.
- A CLI built from the exact annotated v0.1.27 tag initialized a generic isolated consumer and
  reported a healthy current installation.
- Its v0.1.28 upgrade dry-run passed catalog-to-archive, pack-to-binary, and staged-version checks,
  made no writes, and left the generic plan-register SHA-256 unchanged at
  `609eb01cc595604e7de4a6003c05a3d869ced14745cae597a2ad6732fa904ac7`.
- Applying the upgrade succeeded, replaced the installed binary with the exact released macOS
  digest, reported v0.1.28, left the installation healthy, and retained that register digest.

The pull-request Windows job must provide the final real-Windows install/upgrade/failure-path gate
before merge. After publication, the final-public catalog, sidecars, downloads, installer, binary
versions, and digest chain must be reverified from live GitHub URLs.

## Publication Boundary

The packs are unsigned and unnotarized. Publication is authorized under the repository's existing
HTTPS-plus-SHA-256 internal/prototype boundary; no publisher identity attestation is claimed beyond
GitHub transport, repository control, the annotated tag, and the published digest chain.
