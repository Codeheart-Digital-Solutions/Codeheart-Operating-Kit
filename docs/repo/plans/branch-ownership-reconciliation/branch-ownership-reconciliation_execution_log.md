Last updated: 2026-08-09T20:42:44Z (UTC)
Created: 2026-08-09

# Branch Ownership Reconciliation Execution Log

Plan: `branch-ownership-reconciliation_implementation_doc.md`
Mode: goal-style implementation and authorized producer release
Status: active
Overall divergence: none at activation

## Summary

The user explicitly activated the validated seven-epic implementation on 2026-08-09 and
authorized the complete producer implementation, v0.1.27 version/release surfaces, focused ready
pull request, expected-head guarded merge commit, annotated tag, public release, assets, sidecars,
and live verification. Execution uses
`codex/branch-ownership-reconciliation-v0.1.27` from exact v0.1.26/origin-main commit
`7a6cd54b27b3f0d322d6b65e436542c677fa7e95`.

Consumer repository upgrades and catalog migrations remain excluded. Existing unrelated
worktrees and branches remain untouched. The frozen plan register remains unchanged.

## Epic Delta Index

| Epic | Status | Meaningful delta | Review gate |
| --- | --- | --- | --- |
| `EP-01` | completed | Added strict version-dispatched config v2, ledger v3, inventory v3, and catalog v3 contracts while preserving historical schemas. | Focused schema/config compatibility review passed. |
| `EP-02` | completed | Added deterministic local/tracking/remote proof computation, stable-patch incorporation, candidate-scoped inventory evidence, and fail-closed blockers. | Focused evidence and race suites passed. |
| `EP-03` | completed | Added reviewed clearance, one-owner mixed deferral, E-to-L-to-A checkpoints, guarded activation, post-write checks, and exact incremental resolution. | End-to-end migration/activation/rollback suite passed. |
| `EP-04` | completed | Added remote-aware evidence, tracking/mirror coalescing, compatibility observations, cache v3, and readiness semantics. | Portfolio and command focused suites passed. |
| `EP-05` | completed | Updated managed doctrine, routing, runbooks, component manifest, resource mirrors, and impact/release inputs. | Schema/resource/routing/parity suite passed. |
| `EP-06` | completed | Added adversarial, compatibility, synthetic acceptance, remote descendant, remote-overlay invalidation, and portability coverage; full Go, race, Python, and validator gates pass. | Four iterative fresh review rounds repaired every High/Medium finding; final verdict has no implementation blocker. |
| `EP-07` | in progress | Prepared v0.1.27 version, component, manifest, release-note, installer, compatibility, and CI upgrade/platform surfaces; opened ready PR #10 and repaired the Windows CRLF portability defects it exposed. | Independent repair review passed; final-head reproducible assets, green CI, merge, tag, publication, and live verification remain. |

## Review Gate Metrics

- Review gate required: yes, after each epic and before publication.
- Review gate skipped: no.
- Reviewer mode: fresh read-only subagent.
- Reviewer model and reasoning mode: inherited from the implementing agent.
- Review rounds: nine iterative fresh read-only rounds completed, including five release-platform
  repair rounds.
- Material findings: four High and eight Medium findings across post-activation remote checkpoint
  reconstruction, remote descendant drift, tracking/mirror identity, merged-owner resolution,
  compatibility-row invalidation, complete reviewed-overlay digest binding, checkout conversion
  authority, and Git configuration precedence.
- Files changed because of review: `internal/commands/plans.go`, `internal/portfolio/scanner.go`,
  `internal/plancatalog/inventory.go`, `git_index.go`, `migrate.go`, `migration_binding.go`,
  `migration_branch_review.go`, `migration_checkpoint.go`, and focused tests.
- Final accepted result: no High or Medium findings remain; implementation is ready for normal
  release and platform gates.
- Approximate added time: multiple focused repair-and-rerun cycles during EP-06.
- Token usage: not exposed per reviewer.
- Worth-it assessment: yes; the review prevented false remote clearance after ref/source drift and
  proved the only controlled exception is exact reviewed-owner integration bound back to the
  original complete overlay digest.

## Activation Delta

- Repository identity: `codeheart-operating-kit`.
- Source authority: producer paths under `components/`, `profiles/`, `templates/`, `schemas/`,
  `internal/`, `scripts/`, and `docs/repo/`.
- Baseline: `origin/main` and tag `v0.1.26` both resolve to
  `7a6cd54b27b3f0d322d6b65e436542c677fa7e95` after an authenticated fetch.
- Worktree: dedicated isolated Codex worktree; no transaction, recovery, stage, or quarantine
  marker was present.
- Safe default: create `codex/branch-ownership-reconciliation-v0.1.27` without modifying the two
  unrelated producer worktrees or their branches.
- Activation checkpoint scope: implementation plan, sibling execution log, and the three nearest
  documentation indexes only.
- Explicit exclusions from this checkpoint: product implementation, version/release surfaces,
  consumer repositories, frozen plan register, unrelated worktrees, destructive Git, PR, merge,
  tag, and release. Those producer actions begin only after this checkpoint under the user's
  separate explicit authorization in the same request.

## EP-01 Delta - Versioned Evidence And Compatibility Contracts

Status: completed. Strict raw-map version dispatch rejects unknown/new contracts in old and new
CLIs. Historical config v1, migration-ledger v2, and plan-catalog v2 schemas remain byte-preserved;
new safety-bearing contracts are closed schema versions rather than ignored additive fields.

## EP-02 Delta - Deterministic Branch Evidence Engine And Inventory UX

Status: completed. The proof engine binds ref tip, unique merge base, relevant path states,
attributes, candidate policy, source revision, and target candidate set. Exact same-content and one
complete incorporated transition may clear; ambiguity and all incomplete evidence block.

## EP-03 Delta - Reviewed Clearance, Mixed Deferral, And Atomic Migration

Status: completed. Migration and activation are separate guarded transactions. Mixed deferral
requires frozen cutover authority, one exact active owner, and a mandatory incremental follow-up.
Direct canonical mode remains zero-gap. Integrated owner bytes resolve only with the exact reviewed
metadata; drift triggers rollback or validation failure.

## EP-04 Delta - Mixed Visibility, Remote Overlays, And Portfolio Readiness

Status: completed. Fresh remote overlay evidence participates in review and authority digests.
Tracking and mirror observations coalesce only when proof-equivalent. Mixed compatibility rows stay
visible and force `canonical_ready=false`; cache v3 is raw-validated before decode.

## EP-05 Delta - Managed Doctrine, Routing, And Packaged Resource Parity

Status: completed. Source doctrine and exact packaged mirrors describe the new schema, CLI,
activation, rollback, incremental follow-up, and readiness contracts. Public examples are generic.
Impact, release-note, and preliminary acceptance attachments are present.

## EP-06 Delta - Exhaustive Validation, Security Review, And Acceptance Proof

Status: completed. Focused Go packages, repository-wide Go and race tests, schema/resource/routing
validation, and v0.1.25/v0.1.26 compatibility probes pass. The sanitized acceptance cohort proves
12 same-content refs, six incorporated-history refs, and genuine active owners. Ordinary
remote-aware scans now prove pending deferral, alternate-source and moved-ref invalidation, exact
owner integration with retained or deleted refs, branch-only observation retention, bound-digest
tampering rejection, and compatibility-row suppression on failed authority. Four iterative fresh
review rounds closed every High/Medium finding. The final source gates passed: all Go packages,
repository-wide race coverage plus a post-review portfolio race rerun, Go vet, 158 isolated Python
tests, all four validators, and `git diff --check`.

## EP-07 Delta - Versioned Release, Publication Gate, And Consumer Handoff

Status: in progress. v0.1.27 version and release inputs are prepared. Clean reproducible builds,
pack-shape and macOS installer/upgrade verification passed for the first release-candidate head.
Ready PR #10 targets exact v0.1.26 main. Its oldest-admitted Git 2.43, Ubuntu, and macOS jobs passed;
the real-Windows job exposed clean-checkout CRLF candidate hashing and Git warning/stdout parsing
defects. The repair now classifies Git-clean tracked plans from their stage-zero blobs, retains raw
worktree preconditions for writes, accepts only LF/CRLF checkout equivalence after clean index/
revision proof, and separates Git stderr from machine-parsed stdout. Full Go, focused race, vet,
normal affected-package, and CRLF-simulated migration/remote-overlay suites pass locally. An
independent read-only review found and closed two Medium policy-precedence gaps, then returned no
remaining High or Medium findings. The repaired full Go suite, focused race suite, Go vet, Markdown
timestamp/public-core validators, and diff check pass. A fresh exact-head CI and reproducible build
remain required before merge and publication; prior binary asset digests are superseded by this
code repair. Exact-head CI then proved the Git 2.43, Ubuntu, and macOS lanes while Windows reached
Go's default ten-minute package timeout with completed packages passing and a normal plan-catalog
test still running. The Windows source and later grouped plan-catalog commands now preserve their
complete suites and raise only the Go test deadline to 30 minutes; no test is skipped, narrowed,
or made permissive. The release-asset builder's independent source-validation rerun uses the same
deadline so the real-Windows packaging gate cannot silently inherit the old ten-minute default.
The asynchronous Windows handoff tests now allow the production protocol's 30-second parent-exit
window plus five seconds of scheduling headroom while retaining exact binary/state restoration
assertions; their former six-second polling window caused a hosted-runner-only false failure.

## Final Validation

Implementation validation is complete. EP-07 final-head release-build, Windows CI, merge, tag,
publication, and live-verification evidence remains pending.
