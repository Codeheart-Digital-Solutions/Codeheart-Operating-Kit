Last updated: 2026-08-07T01:39:50Z (UTC)

# Repository-Wide Plan Catalog Discovery Release Readiness Evidence

## Bound Source

- Implementation commit: `f9d915256cb3dec1143c78a2d4dc79eb002c0ae2`.
- Branch: `codex/operating-kit-multi-root-plan-catalog`.
- Implementation is active behind explicit discovery-version selection. Existing repository config
  remains discovery v1; no consumer activation occurred.
- The readiness checkpoint is the commit containing this evidence. EP-09 must use that exact
  pushed checkpoint as its source and must stop if source changes invalidate this evidence.

## Platform Validation

GitHub Actions run
[`31137575877`](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/actions/runs/31137575877)
passed on the exact implementation commit:

- Ubuntu semantic validation: passed.
- macOS validation, including staged install/upgrade behavior: passed.
- Real Windows validation, including staged install/upgrade behavior: passed.
- Public-release verification jobs were correctly skipped for the branch push.

The preceding run exposed a Windows-only test-fixture issue caused by a `core.autocrlf` warning
being captured with a Git object ID. Commit `f9d9152` pins `core.autocrlf=false` only for that
fixture's `hash-object` call. The final run above proves the remediation on real Windows without
weakening runtime semantics.

## Reproducible Readiness Builds

Two isolated builds ran from the implementation source:

```sh
env GOCACHE=<isolated-cache> python3 scripts/build-release-assets.py --output-dir <build-a>
env GOCACHE=<isolated-cache> python3 scripts/build-release-assets.py --output-dir <build-b>
diff -rq <build-a> <build-b>
```

Both builds reported `reproducible: true`; the recursive byte comparison produced no differences.
The source currently identifies version `0.1.24`, so these artifacts are readiness candidates only.
EP-09 must select the next doctrine-correct version, update the release surfaces, and rebuild the
final public assets twice.

| Readiness artifact | SHA-256 | Bytes |
| --- | --- | ---: |
| `codeheart-operating-kit-0.1.24-macos-universal.zip` | `363657ef4809e94de2ffce00fe7da88b833b4fbf57b84f2a903dfbf93dedd608` | 7,707,316 |
| macOS archive sidecar file | `00dcfc9692f5f601fcc2097c483a110c3027c20f3f9cc404a03c58980bec5dbb` | 117 |
| `codeheart-operating-kit-0.1.24-windows-x64.zip` | `f3d3d7825af3d2862f8dec7e59c9de0c32816bfbc36fb8b7f57dd15e5e03c844` | 4,066,850 |
| Windows archive sidecar file | `c2858205a6f945561ce629a89ad8aca5d8faa6c8ad60c55de5156a7a7d684731` | 113 |
| `release-catalog-0.1.24.json` | `a1a8bb459b9807d4972cd19323e0476429f3df1443d6123115000dd164ce9f1d` | 857 |

Catalog entries bind the macOS archive to pack-manifest digest
`64eb7f3322b21e5ed37fb46a9a49544ead1038aee14421e26238b9c3a92e91e4` and the Windows archive to
pack-manifest digest `5fd8c0c11df4c68e9fcd6124fc00b55700071212715005dea248371802a8821c`.

Both archives contain only the required release shape: `INSTALL.md`, the platform binary,
`bootstrap.md`, `checksums.txt`, `content-manifest.yaml`, `install.ps1`, `install.sh`,
`pack-manifest.json`, and `release-notes.md`. Entries use deterministic timestamps. No wheel,
`*.dist-info`, or other Python runtime payload is present.

## Catalog-To-Binary Verification

A CLI stamped as the prior version performed a local macOS upgrade dry-run through the generated
external catalog:

```sh
codeheart-operating-kit upgrade \
  --version 0.1.24 \
  --catalog <build-a>/release-catalog-0.1.24.json \
  --dry-run \
  --json \
  <isolated-consumer>
```

Catalog-to-archive, pack-to-binary, and staged-version validation passed. The planned transaction
made no writes. Recorded provenance:

- archive: `363657ef4809e94de2ffce00fe7da88b833b4fbf57b84f2a903dfbf93dedd608`;
- binary: `00e24df46857476c6ac11047aecb2bdbf7340605ce09ce9691c72f1ac29d4d85`;
- catalog: `a1a8bb459b9807d4972cd19323e0476429f3df1443d6123115000dd164ce9f1d`;
- content manifest: `a63cc0d0362e83e402394fa4fb7cfab7760feaae66bb9e27ebda9647ec3160f3`;
- pack manifest: `64eb7f3322b21e5ed37fb46a9a49544ead1038aee14421e26238b9c3a92e91e4`.

The successful macOS and real-Windows Actions jobs provide the corresponding isolated platform
install/upgrade/failure-path evidence at the same implementation source.

## Signing And Audience Boundary

These local assets are unsigned release-candidate evidence and must not be published. Repository
doctrine currently documents HTTPS plus SHA-256 under an unsigned internal/prototype boundary.
EP-09 must reconfirm that boundary for the intended release audience or establish the required
signing/notarization before broad public distribution; unresolved signing/audience ambiguity is a
release blocker, not a reason to substitute these local files.

## EP-09 Publication Inputs

EP-09 must:

1. inspect current main, tags, releases, branch ownership, authentication, identity, and signing;
2. select the doctrine-correct version after `v0.1.24` and update every version, note, installer,
   manifest, content-identity, and release input surface required by the release runbook;
3. run the complete local and platform validation plus two byte-identical final builds;
4. merge the exact validated release commit through the normal protected PR workflow;
5. create the matching `v<version>` tag without force and publish `bootstrap.md`, `install.sh`,
   `install.ps1`, `release-notes.md`, `manifest.yaml`, both platform archives and sidecars, and the
   external release catalog through the normal GitHub release workflow; and
6. verify the live tag, asset URLs, sidecars, catalog chain, binary versions, and consumer channel
   before handing the release to EP-10.

No embedded `manifest.yaml` archive URL or digest may be introduced. Final public URLs belong only
in the external release catalog generated after the final packs.

## Codeheart-HQ Installation Handoff

After publication, EP-10 must use a separate, unambiguously owned HQ worktree/branch based on clean
current main and follow HQ's `AGENTS.md` plus the installed lifecycle runbook:

1. Record pre-upgrade Kit version, `check` output, config/lock/managed state, active discovery
   version and catalog mode, plan inventory/validation result, relevant plan bytes, status,
   transactions, and overlapping branches.
2. Run the published release's verified
   `upgrade --version <released-version> --dry-run`, review the exact managed changes, then run
   `upgrade --version <released-version> --yes`.
3. Verify `codeheart-operating-kit --version`, `codeheart-operating-kit check <HQ>`, config/lock and
   managed-resource parity, and exact installed-version/catalog provenance.
4. Prove installation did not activate discovery v2, change catalog mode, rewrite the plan
   register, alter plan bytes, promote a v2 cache, or modify semantic plan metadata. Repeat the
   smallest discovery-v1 plan list/validate evidence before and after.
5. Stage only explicit upgrade paths, commit/push, and complete HQ's normal protected PR/merge
   workflow; synchronize and repeat final evidence on current HQ main when required.

Prospective v2 inventory, HQ metadata migration, canonical activation, portfolio-v2 completeness,
and the paused semantic cutover are deliberately excluded from this installation checkpoint.

## Residual Risks And Deferred Decisions

- The release that eventually removes discovery v1 remains a later product decision.
- The conventional ambiguity segment list may be refined only from production evidence in a later
  reviewed change.
- New tracked plan-like documents beneath owned `docs` paths become visible after activation and
  may correctly block canonical validation until reviewed.
- Final version-specific assets and their signing/audience state must be re-established by EP-09;
  readiness digests above cannot be reused as publication evidence after version mutation.
