Last updated: 2026-08-10T07:19:50Z (UTC)

# v0.1.28 Release Readiness Evidence

## Validated Source

- Release-source commit: `bff24c8644f78c2490dbb723aff8ce2df6a94ba0`.
- Branch: `codex/legacy-register-inventory-v0.1.28-hotfix`.
- Base: released v0.1.27 merge `70a8e174fe8767b09eb6e0ccbea64f6fad5e9a3f`.
- Consumer impact: `validator-only change`; no consumer migration, catalog activation, register
  rewrite, sync, or repair is required.
- Frozen producer register: unchanged from base.

The exact release source passed full Go, Go race, Go vet, 158 Python tests, public-core,
Markdown-timestamp, JSON-schema, release-manifest, packaging-resource, routing, installer, and
backward-compatibility validation. Independent review reported no remaining High or Medium
blocker. The source also includes the bounded Windows previous-binary restore retry and preserved
backup evidence required after the first pull-request run exposed a transient executable-release
race in the existing legacy lock-v1 rollback gate.

## Reproducible Release Assets

`scripts/build-release-assets.py` ran twice in independent build directories with final public
v0.1.28 URLs. Each invocation also built every platform pack twice internally and required byte
equality. Build A and build B matched for both archives, both archive sidecars, and the external
catalog.

| Asset | SHA-256 |
| --- | --- |
| `codeheart-operating-kit-0.1.28-macos-universal.zip` | `ffc3cbea73241099fc6820f97cc13a4294f65b096421f7e7aa88779514a63701` |
| `codeheart-operating-kit-0.1.28-macos-universal.zip.sha256` | `9a6e4197402d6f3835c6f29d6139cbef04ae7e0ee2db687b5b43f87f9dcb4507` |
| `codeheart-operating-kit-0.1.28-windows-x64.zip` | `46b8d44285e917310b2c751dd5dab72c0ed0b1805034268ef72b7da201d14e61` |
| `codeheart-operating-kit-0.1.28-windows-x64.zip.sha256` | `8947b345fa7a325169c39ce2bc6e64e371bef1c37eaf089e1444f0477226cd18` |
| `release-catalog-0.1.28.json` | `124422c57df4733a23b0fb87cb786b7089b38b157d5b1ab737c3906f5d9768ee` |

The catalog binds pack-manifest SHA-256
`9ac99c034f55056407730e94bd89b95f1920700e5b473565173b049b0ea73f90` for macOS and
`c6ae8fa27a6f7a9e81c6f92b42d73ca13d5f81487c1a81d63029eb0717ce14e8` for Windows. Both packs
bind content-manifest SHA-256
`7085ecec9648ef424fa2d2f0cc3cf8133cac30c9aa9fbe503345fff7f4fee9ab`.

The macOS universal binary is Mach-O x86_64 plus arm64, reports v0.1.28, and has SHA-256
`85e824433157e39e67e9ccc4cce811c9aff64ed98647cddba412a31bfcd7346a`. The Windows binary is
PE32+ console x86-64 and has SHA-256
`b31272945df2a32e9f070f198542284e5547d6485fc8a87ec066310e713ed8fd`.

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
- Applying the upgrade succeeded against the corrected release-source asset, replaced the
  installed binary with exact SHA-256
  `85e824433157e39e67e9ccc4cce811c9aff64ed98647cddba412a31bfcd7346a`, reported v0.1.28, left
  the installation healthy, and retained that register digest. The local-URL catalog used for
  this proof has SHA-256
  `7c35130cd905131112c58dc2f6ed1ca8834dd126fd1f7b22c5dfb7550696d1e4`; its four archives and
  sidecars are byte-identical to the final-public set.

The pull-request Windows job must provide the final real-Windows install/upgrade/failure-path gate
before merge. After publication, the final-public catalog, sidecars, downloads, installer, binary
versions, and digest chain must be reverified from live GitHub URLs.

## Publication Boundary

The packs are unsigned and unnotarized. Publication is authorized under the repository's existing
HTTPS-plus-SHA-256 internal/prototype boundary; no publisher identity attestation is claimed beyond
GitHub transport, repository control, the annotated tag, and the published digest chain.
