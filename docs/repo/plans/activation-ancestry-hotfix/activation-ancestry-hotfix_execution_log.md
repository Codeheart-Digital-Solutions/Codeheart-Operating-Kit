Last updated: 2026-08-10T22:23:09Z (UTC)
Created: 2026-08-10

# Activation Ancestry Resolution Hotfix Execution Log

Plan: `activation-ancestry-hotfix_implementation_doc.md`
Mode: goal-style implementation, authorized producer release, and isolated Kit-only rollout
Status: active
Overall divergence: none at activation

## Summary

A production validation failure was reproduced from externally verified topology evidence: a
reviewed activation checkpoint remained the second parent of a normal merge commit whose tree was
identical to the activation tree, but local persisted-evidence validation searched only the merge
commit's first-parent chain and selected the merge itself as the proposed activation checkpoint.

The correction preserves the exact `E -> L -> A` checkpoint edges and every existing content and
authority check. Only the way `A` is located from a later `HEAD` changes: resolution follows the
complete ancestry path, filters direct first-parent children of `L` by the persisted action digest,
removes later matches that descend from another match, and requires one unique ancestry-minimal
incorporated checkpoint. Independent matching siblings remain ambiguous and block.

## Epic Delta Index

| Epic | Status | Meaningful delta | Review gate |
| --- | --- | --- | --- |
| `EP-01` | completed | Sanitized topology resolver and exact guarded-config enforcement implemented; focused and full affected suites green. | Accepted; no High/Medium blocker. |
| `EP-02` | in progress | v0.1.29 surfaces and reproducible assets complete; full Go/race/vet and 158-test Python suite green; PR and publication pending. | Release review required. |
| `EP-03` | pending | Isolated Kit-only consumer confirmation pending. | Baseline and byte-preservation gates required. |

## Review Gate Metrics

- Review gate required: yes.
- Review gate skipped: no.
- Reviewer mode: fresh read-only reviewer agent when available; otherwise strongest independent
  main-thread review with explicit evidence.
- Review rounds: two completed rounds, including a fresh review after the exact-head Windows
  failure.
- Material findings: four Medium findings resolved: ancestry-minimal selection was documented;
  persisted wrong-ledger/delta/blob negatives were added; exact guarded activation config is now
  reconstructed and byte-checked in both local and remote validation; and the CRLF committed-byte
  boundary gained a direct negative authority matrix.
- Final accepted result: accepted; no remaining High or Medium code or test blocker.

## Activation Delta

- Producer base: released v0.1.28 merge `38e9b85180c94f06dcddf2fd69c14db6bae85996`.
- Branch: `codex/activation-ancestry-v0.1.29`.
- Checkpoint scope: this plan, this execution log, and the nearest plans README only.
- Exclusions: validator source, tests, version surfaces, release assets, consumer repositories,
  and all unrelated local work remain outside the activation checkpoint.
- Consumer impact: validator-only; no migration or semantic adoption is authorized by the Kit
  upgrade itself.

## EP-01 Delta - Topology And Resolver

The live evidence confirms the reported topology, while public tests and committed evidence use
only synthetic repository history. The bounded plan checkpoint was committed and pushed before
production implementation.

Plan activation validation confirmed that the new canonical record is enumerated once with the
expected active lifecycle, semantic ID, and path. Public-core, Markdown-timestamp, and diff checks
pass. Repository-wide plan validation remains red only on pre-existing malformed, unsupported, and
duplicated relation forms in the frozen legacy register; this corrective plan neither introduced
those findings nor rewrites that historical evidence.

The implementation replaces first-parent traversal from current `HEAD` with one shared resolver.
It enumerates the full ancestry path from `L` to the requested target, retains only commits whose
first parent is exactly `L` and whose plan-action digest matches the persisted binding, removes
later matching wrappers that descend from another match, and requires one unique ancestry-minimal
checkpoint. Local persisted validation, command-side target reconstruction, portfolio target
reconstruction, and compatibility projection all use that resolver.

Review identified that the plan-action digest intentionally excludes the guarded config. Without
an additional check, a forged sibling could copy the exact reviewed plan delta while changing a
different schema-valid config field. The corrected implementation now reads the config at `L`,
verifies the ledger's config precondition, deterministically rebuilds the only permitted schema-v2
activation config, and byte-compares it with `A`. The same shared check is required by local and
remote validation.

Sanitized regressions now cover exact `A`, a normal merge with `A` as second parent and an identical
tree, a safe descendant, a merge wrapper whose first parent is `L`, an unrelated sibling,
tree-only reparenting, a wrong action digest, independent ambiguous matching children, a wrong
activation-base binding, a wrong projected blob, and matching-digest sibling config forgery. The
command and portfolio suites also cover merged target reconstruction and remote config forgery.

Focused results after the security repair:

- exact topology/config regression group: pass;
- command-side merged target reconstruction: pass;
- portfolio merged target and forged-config validation: pass;
- full affected packages `internal/plancatalog`, `internal/commands`, and `internal/portfolio`:
  pass after config hardening (`280.340s`, `84.332s`, and `132.718s` respectively in the
  independent review run);
- `go vet` for all three affected packages, gofmt, and `git diff --check`: pass.

A controlled clean-clone reproduction of the reported production merge retained only the exact
ledger-reviewed refs. The released v0.1.28 binary failed solely on
`activation_checkpoint_parent_mismatch`. The current candidate validated the same immutable merge
with 55 rows, 54 canonical records, one reviewed mixed-grandfathered record, complete mixed
coverage, zero errors, and no worktree or tree mutation. Exact private topology identities remain
outside this public execution record.

## EP-02 Delta - Compatibility And Release Preparation

The release identity is now v0.1.29 across the Go and Python versions, bootstrap, macOS and Windows
installers, content manifest, planning-workflows component, standard profile, packaged resource
mirrors, release-candidate fixture, upgrade and parity tests, and release notes. The component and
profile checksums and the compiled standard-profile graph digest were recomputed from the final
retained resources; root and packaged copies are byte-identical.

The release notes classify the correction as validator-only. Upgrade does not inventory, migrate,
reactivate, normalize, repair, or sync a catalog and does not change consumer plan, register,
ledger, catalog, or config bytes. Existing schema and transaction compatibility remains unchanged.

Validation after the version and identity update:

- `go test -timeout 30m ./...`: pass;
- `go test -race -timeout 30m ./...`: pass;
- `go vet ./...`: pass;
- complete Python suite: `158 passed`;
- release-focused Python subset, including repeated platform builds: `94 passed`;
- JSON schemas, release manifests, Markdown timestamps, public-core hygiene, packaged-resource
  mirrors, content checksums, graph identity, backward compatibility, installers, routing, and CLI
  parity: pass;
- `git diff --check`: pass.

The first local Python invocation used the sandbox-blocked default Go cache and failed only where
subprocesses attempted to open that cache. The exact suite passed with a writable isolated
`GOCACHE`; no product assertion changed and no dependency or source workaround was required.

The first exact-head Windows candidate run exposed a platform-specific persisted-evidence gap in
the new merge fixtures. Git for Windows materialized the Git-clean reviewed ledger and config with
CRLF bytes after the branch switch, while local validation compared those checkout bytes directly
with the immutable LF blobs. The focused correction now proves each current path is regular,
stage-zero clean, and equal to the committed blob with only LF/CRLF checkout equivalence admitted;
all ledger/config hashes and security comparisons continue to use the exact committed bytes.
Both affected tests force CRLF checkout materialization on every platform, so the regression is no
longer Windows-only. Direct negative coverage rejects dirty worktrees, staged blob or mode drift,
non-regular and symbolic-link replacements, and a Git-clean checkout transform beyond line-ending
conversion. The final independent review accepted this boundary with no remaining High or Medium
finding. The failed run remains a closed gate until a replacement exact-head run passes.

Final validation after the CRLF correction and negative authority matrix:

- focused CRLF merge, command reconstruction, and authority-boundary regressions: pass;
- full affected packages: pass (`internal/plancatalog` `198.665s`, `internal/commands` `37.038s`,
  and `internal/portfolio` `59.687s`);
- repository-wide `go test -timeout 30m ./...`: pass;
- repository-wide `go test -race -timeout 30m ./...`: pass;
- repository-wide `go vet ./...`: pass;
- complete Python/schema/resource/routing/release suite: `158 passed in 334.29s`;
- gofmt, `git diff --check`, archive sidecars, payload checksums, binary platform/version identity,
  and isolated macOS installer evidence: pass.

Final release-candidate identity:

- macOS universal archive SHA-256:
  `bfbd76a9e8905019abc1b7f571d0526439665395d6d485cfe25694a5eb5b7baf`;
- macOS pack-manifest SHA-256:
  `f673401dac813ea1df052f1d70fef910d3990fca55ee1b014c28afcc3fe03478`;
- macOS binary SHA-256:
  `68fc6e058e4a9ce96813d31bdf332bf734c784dec11bbd01e6bad0c8d4888529`;
- Windows x64 archive SHA-256:
  `99a3472c3870b60626dad1237de88874785aa89d812c072d4e4c1f8cf615290a`;
- Windows pack-manifest SHA-256:
  `a5332d544b5b1df84d93e354d66b19d9a3b764d625d365355a8ab3608274bb44`;
- Windows binary SHA-256:
  `716d3c899700ab74c29be2a359bbd94997495bc4ed2ddb97e430c9d2d778843f`;
- embedded content-manifest SHA-256:
  `71c23246364595e231c211fe71ce32a0e21d18c89cac2b1d86646cc5cd4ce66c`;
- external release-catalog SHA-256:
  `0b96f26c556ccd6e270f39a85773d2ac3c591060cd225cad474f557e6f939b27`.

The final builder reran source validation, built each platform twice, required byte equality, and
emitted the catalog only after the packs. Sidecars verify both archives. Pack-manifest, payload
checksum, content-manifest, binary digest, platform, and version chains all match. The macOS pack
installed through the exact local catalog and installer, reported v0.1.29, and contained x86_64 and
arm64 slices. The Windows binary is PE32+ x86-64. Neither pack contains a Python payload.
