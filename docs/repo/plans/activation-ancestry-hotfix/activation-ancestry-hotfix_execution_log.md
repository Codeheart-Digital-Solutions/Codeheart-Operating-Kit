Last updated: 2026-08-10T20:22:47Z (UTC)
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
| `EP-02` | pending | Full compatibility, release preparation, and publication pending. | Release review required. |
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
