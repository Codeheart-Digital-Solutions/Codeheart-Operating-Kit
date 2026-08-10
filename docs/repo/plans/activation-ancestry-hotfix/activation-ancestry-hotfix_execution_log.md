Last updated: 2026-08-10T19:08:24Z (UTC)
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
and requires a unique incorporated match.

## Epic Delta Index

| Epic | Status | Meaningful delta | Review gate |
| --- | --- | --- | --- |
| `EP-01` | in progress | Plan activated; sanitized topology regression and resolver implementation pending. | Focused review required. |
| `EP-02` | pending | Full compatibility, release preparation, and publication pending. | Release review required. |
| `EP-03` | pending | Isolated Kit-only consumer confirmation pending. | Baseline and byte-preservation gates required. |

## Review Gate Metrics

- Review gate required: yes.
- Review gate skipped: no.
- Reviewer mode: fresh read-only reviewer agent when available; otherwise strongest independent
  main-thread review with explicit evidence.
- Review rounds: zero at activation.
- Material findings: pending.
- Final accepted result: pending.

## Activation Delta

- Producer base: released v0.1.28 merge `38e9b85180c94f06dcddf2fd69c14db6bae85996`.
- Branch: `codex/activation-ancestry-v0.1.29`.
- Checkpoint scope: this plan, this execution log, and the nearest plans README only.
- Exclusions: validator source, tests, version surfaces, release assets, consumer repositories,
  and all unrelated local work remain outside the activation checkpoint.
- Consumer impact: validator-only; no migration or semantic adoption is authorized by the Kit
  upgrade itself.

## EP-01 Delta - Topology And Resolver

The live evidence confirms the reported topology, but public tests and committed evidence will use
only synthetic repository history. Implementation and focused validation begin after the bounded
plan checkpoint is committed and pushed.

Plan activation validation confirmed that the new canonical record is enumerated once with the
expected active lifecycle, semantic ID, and path. Public-core, Markdown-timestamp, and diff checks
pass. Repository-wide plan validation remains red only on pre-existing malformed, unsupported, and
duplicated relation forms in the frozen legacy register; this corrective plan neither introduced
those findings nor rewrites that historical evidence.
