Last updated: 2026-08-10T21:10:34Z (UTC)
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
- Review rounds: one completed round.
- Material findings: three Medium findings resolved: ancestry-minimal selection was documented;
  persisted wrong-ledger/delta/blob negatives were added; exact guarded activation config is now
  reconstructed and byte-checked in both local and remote validation.
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

Final release-candidate identity:

- macOS universal archive SHA-256:
  `14a64e1ae25088291046353858a153fef0f3938340431cacb3937c0d50850449`;
- macOS pack-manifest SHA-256:
  `7857dd4653e5963bdd2290ad9a3899b34144e685d3e4c15213567b3358b5aa39`;
- macOS binary SHA-256:
  `a3b7e91615c2c4e20004c4392c3648a8ceb79509e65ab8f3496990240e206d7a`;
- Windows x64 archive SHA-256:
  `c4e9cf0c06939dbceddb26ee350de8c0fe8fa76762a8b770dc7968d1a78fa856`;
- Windows pack-manifest SHA-256:
  `4d94a7f8209a68e33b46e4347a05b8ad590679ad7f46faf1ed82da92213ab95c`;
- Windows binary SHA-256:
  `79972ec886c69a8336f96b5cad7e11cbfa32690fbe07bccc495a2b00363632e5`;
- embedded content-manifest SHA-256:
  `71c23246364595e231c211fe71ce32a0e21d18c89cac2b1d86646cc5cd4ce66c`;
- external release-catalog SHA-256:
  `2b03bf5b8ad5196e0699b8bd51ef2e4fb2415a63610cb5d1b2b3c23f6bc10bf6`.

The final builder reran source validation, built each platform twice, required byte equality, and
emitted the catalog only after the packs. Sidecars verify both archives. Pack-manifest, payload
checksum, content-manifest, binary digest, platform, and version chains all match. The macOS pack
installed through the exact local catalog and installer, reported v0.1.29, and contained x86_64 and
arm64 slices. The Windows binary is PE32+ x86-64. Neither pack contains a Python payload.
