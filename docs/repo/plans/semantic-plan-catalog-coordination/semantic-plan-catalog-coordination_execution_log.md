Last updated: 2026-08-01T05:12:42Z (UTC)
Created: 2026-07-31

# Semantic Plan Catalog And Branch-Aware Coordination Execution Log

Plan: `semantic-plan-catalog-coordination_implementation_doc.md`
Mode: goal-style implementation
Status: active
Overall divergence: none at activation

## Summary

The user explicitly activated the seven-epic implementation and authorized the bounded planning
checkpoint on 2026-07-31. Execution uses
`codex/semantic-plan-catalog-coordination` and proceeds in epic order. Existing unrelated
`.codeheart/`, agent-memory, and pending-sync work remains outside the activation checkpoint.
Public release execution remains subject to the `EP-07` release gate recorded in the plan.

## Epic Delta Index

| Epic | Status | Meaningful delta | Review gate |
| --- | --- | --- | --- |
| `EP-01` | completed | Added strict schemas, typed metadata parsing, identity validation, canonical discovery, fixtures, and config-v2 compatibility. | Accepted in round two |
| `EP-02` | completed | Added compatibility reconciliation, deterministic views, Git inventory, grouped commands, and optimistic transactional migration. | Accepted in round twenty-three |
| `EP-03` | completed | Added exact portfolio configuration and membership, local/GitHub discovery with optional PR enrichment, scanner-owned policy-confined bare mirrors, byte-aware default and unmerged-branch observations, atomic complete cache, and remote-overlay commands; hardened compatibility, evidence completeness, source authority, remote normalization, local namespace containment, atomic publication, and Git command-policy confinement after review. | Accepted in round fourteen |
| `EP-04` | completed | Added exact semantic-catalog and portfolio references, mode-aware planning and activation publication doctrine, setup/refresh/migration recipes, stable entry-point and strategy scaffolds, routing, declarations, packaged mirrors, and materialization tests. | Accepted in round three |
| `EP-05` | active | Began the producer-only compatibility migration with an explicit legacy-mode identity/config adoption and frozen-register baseline sequence. | Pending |
| `EP-06` | active | Added cross-platform grouped validation, low-context and activation probes, adversarial fixtures, public-safe remote evidence, impact/release/adoption documentation, discovery evidence mapping, benchmarks, and read-only migration handoffs; immutable candidate handoff remains pending. | Pending |
| `EP-07` | pending | None yet. | Pending explicit release gate |

## Review Gate Metrics

- Review gate required: yes, after every epic.
- Review gate skipped: no.
- Reviewer mode: fresh read-only subagent.
- Reviewer model and reasoning mode: inherited from the implementing agent.
- Review rounds: two for `EP-01` (round one rejected; round two accepted).
- Material findings: five in round one (two high, three medium): incomplete family evidence,
  absent-mode defaulting, incomplete positive/negative schema fixtures, unreachable unsupported
  version diagnostics, and silent compatibility-title authority.
- Files changed because of review: `internal/plancatalog/model.go`, `metadata.go`, `validate.go`,
  `plancatalog_test.go`, plan and portfolio fixtures, catalog/migration taxonomy patterns, and
  `tests/test_json_schemas.py`.
- Final accepted result: `EP-01` accepted with no material findings remaining.
- Approximate added time: eleven minutes from recorded first-round remediation start to validation;
  token usage is not exposed per subagent review.
- Worth-it assessment: yes. The first round found five contract/evidence gaps before downstream
  command work could depend on them; the second round confirmed the remediation without widening
  epic scope.

## Activation Delta

- Safe default: create the dedicated `codex/semantic-plan-catalog-coordination` branch before
  source implementation.
- Authority: the activation request authorizes the canonical plan, execution log, and directly
  required planning metadata checkpoint to be committed and normally pushed without another push
  prompt.
- Scope boundary: no unrelated local files, pull request, merge, consumer upgrade, other-repository
  write, tag, or release is included in the activation checkpoint.

## EP-01 Delta - Canonical Contracts, Identity, And Validation Foundation

Status: completed and accepted in review round two.

- Safe compatibility default: existing implementation plans use `# Document Header` followed by
  an H2 human title, while the approved new contract uses the canonical title heading directly.
  The parser accepts the current generated shape and treats its immediate H2 as title authority;
  new direct-H1 documents follow the approved contract unchanged.
- Added task: update the embedded content graph digest after the managed Kit config schema changed;
  no release version was changed.
- Validation: `go test ./...` passed with a sandbox-local Go cache; all 34 focused JSON-schema
  tests passed; `scripts/validate-json-schemas.py` passed; `gofmt` and `git diff --check` passed.
- Round-one review: rejected with two high and three medium evidence/diagnostic gaps. Remediation
  added qualifying family fixtures and derived membership, explicit empty-mode-to-legacy behavior,
  actual positive/negative instances for every durable schema, early unsupported-version and raw
  UTC diagnostics, an observable `legacy_title_layout` warning, direct-H1 rejection for a generic
  header without a title, static same-ID and mixed-mode fixtures, rename-stability proof, separate
  ref/commit observations, and portfolio-v1 fixture validation.
- Post-remediation validation: `go test ./...` passed; all 35 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py` and `git diff --check` passed.
- Round-two review: accepted with no material findings. It confirmed family qualification and
  derived membership, absent-mode defaulting, every durable schema's positive/negative coverage,
  stable unsupported-version handling, narrow and observable implementation-title compatibility,
  static rename/same-ID/mixed evidence, and config-v1/v2/home-self-member behavior. No accidental
  future-epic implementation was found.

## EP-02 Delta - Local Catalog Views And Guarded Migration Mechanics

Status: completed and accepted in review round twenty-three.

- Compatibility findings from the real producer corpus: quoted marker examples inside Markdown
  fences must not count as canonical metadata; one historical discovery uses
  `implementation-handoff-ready`; some implementation plans have only the generic
  `# Document Header`; and the existing untracked v1 config encodes empty `component_settings` as
  YAML null. The reader now ignores fenced marker examples, inventories the historical status as
  draft with an explicit warning, exposes generic implementation titles as observable legacy
  compatibility, and normalizes null component settings in memory without changing config bytes.
- Added deterministic legacy-register reconciliation and local text/JSON views. Canonical metadata
  wins; register evidence fills only missing compatibility fields. Mixed mode grandfathers only
  paths with register evidence, while canonical mode rejects every metadata gap.
- Added Git-backed inventory with exact revision and byte hashes, metadata coverage, unpaired
  register evidence, path-scoped dirty overlap, and changed paths from unmerged local/remote refs.
  A read-only producer inventory found 33 formal legacy records, one unpaired legacy entry, zero
  active-branch touches, one dirty formal plan, and no inventory errors.
- Added schema-validated reviewed-ledger loading, semantic ambiguity/deferral checks, path and
  repository containment, inventory/source-revision/hash checks, active-branch and dirty-target
  skips, marker-bounded insertion, metadata revalidation, and exact precondition hashes at commit.
  Reused `internal/reconcile` staging with a new reviewed-file plan constructor and optimistic
  action preconditions for atomic multi-file rollback and recovery-required results.
- Added grouped `plans validate`, `plans list`, `plans inventory`, and `plans migrate` command and
  help dispatch. Inventory is the only non-`--yes` artifact write and cannot overwrite an
  enumerated plan; migration requires exactly one of `--dry-run` or `--yes`.
- Simulation evidence: alpha changed from
  `b04922785df534dbdb7ee25fe0e8deb278a74b9de01ea22796c65e2d20a57fad` to
  `fe1ff37216c1c56a5ff2de1244fea851c0f3c28f7591292f3d83aed5fc7b7633` while preserving
  `Last updated: 2026-07-31T09:00:00Z (UTC)`; beta changed from
  `4746448a59cf7c941c51c7e511313a9325b3d9a4820230e3242a0aa8ee814a77` to
  `802b134cad6cf074ec96777a0841204cc008ed2479803437e555e94583070ec9` while preserving
  `Last updated: 2026-07-31T09:05:00Z (UTC)`. The second apply succeeded with zero changes and two
  explicit `already_applied` records.
- Validation: the focused catalog/commands/CLI/reconcile suites passed; `go test ./...` passed;
  all 35 focused JSON-schema tests passed; `scripts/validate-json-schemas.py` and
  `git diff --check` passed. Tests cover mixed/canonical coverage, stable output, dry-run,
  idempotency, source mismatch, dirty overlap, unrelated dirty files, active-branch ownership,
  semantic ambiguity, invalid placement, concurrent commit-time edits, successful rollback, and
  recovery-required rollback failure.
- Round-one review: rejected with five high and one medium finding. The reviewer found that the
  known producer-only `implementation-handoff-ready` status could be inventoried but not migrated;
  fenced metadata-marker examples blocked insertion; a partially committed current action was
  omitted from rollback; rollback could overwrite a concurrent edit to an earlier committed
  action; branch deletion and rename-old-path evidence was excluded; and `plans list` lacked the
  approved `--format text|json` contract.
- Remediation: the strict header reader now maps only the known historical producer status to
  draft while retaining and warning on its original value; insertion shares the fenced-marker
  semantics of parsing; transaction records distinguish mutation from completed commit and
  rollback verifies both backups and current target bytes before restoration; Git branch scanning
  includes deletions and both rename paths; and list output accepts `--format text|json` while
  retaining `--json` as a compatibility alias. Regression tests exercise the historical status
  through the complete migration transaction, four-backtick marker examples, current-action
  partial commit, concurrent earlier-target mutation, branch deletion, both rename paths, and all
  format flag combinations.
- Post-remediation validation: focused catalog/reconcile/commands/CLI tests passed; `go test ./...`
  passed; all 35 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-two review: rejected with one additional high finding after confirming that all six
  round-one findings were resolved. The remaining check/rename race allowed a target to change
  after the first expected-hash check but before its bytes were moved to backup.
- Additional remediation: each action now verifies the exact bytes actually moved to backup
  against `ExpectedSHA256` before installing staged content. A disappeared target also fails the
  post-validation precondition. The transaction returns the mutated current-action record so
  rollback restores the concurrent bytes. A deterministic `pre-commit` regression writes new
  bytes after the first check and proves that the migration rolls back without overwriting them.
- Post-round-two validation: focused catalog/reconcile/commands/CLI tests passed; `go test ./...`
  passed; all 35 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-three review: rejected with two high and two medium findings. The reviewer found that a
  target recreated after backup validation could still be replaced by the staged rename; inventory
  output protected only successfully parsed records and could overwrite an omitted malformed plan
  or the retained legacy register; canonical migration did not block projected incomplete
  coverage; and the approved inventory command contract required `--output` while help and dispatch
  treated it as optional. The reviewer reconfirmed all prior findings as resolved.
- Additional remediation: staged installation now uses an atomic hard-link no-replace operation,
  so a recreated target remains untouched and the transaction retains the original backup as a
  `recovery-required` case. Inventory output protection enumerates every formal plan independently
  of parse success and also protects `plan-register.md`. Canonical migration projects valid
  existing metadata plus planned actions across the full formal enumeration and adds an
  `incomplete_coverage` blocker before any write. `plans inventory` and its help now require an
  explicit `--output` destination. Regression tests prove all four behaviors, including exact byte
  preservation and zero-write incomplete canonical migration.
- Post-round-three validation: focused catalog/reconcile/commands/CLI tests passed; `go test ./...`
  passed; all 35 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-four review: rejected with two high findings after confirming the no-replace transaction,
  exact-path protection, canonical coverage, required-output behavior, and all earlier remediations.
  A repo-local symlinked parent could bypass lexical inventory-target protection, and the inventory
  record set still came from successfully parsed discovery records rather than the complete formal
  enumeration, causing malformed readable plans to appear only as problems and not counted records.
- Additional remediation: inventory destinations now reject repo-local symlink traversal, resolve
  missing-path identities through their nearest existing ancestor, compare resolved identities for
  every protected plan and the legacy register, and write through the resolved path. Inventory is
  now enumeration-driven: every formal candidate receives a record with parse status, problem
  codes, exact current hash when readable, dirty/branch evidence, and explicit canonical, legacy,
  invalid, or unreadable coverage. Coverage separately counts invalid and unreadable records. Tests
  reproduce both plan and register symlink aliases and a committed malformed formal plan.
- Post-round-four validation: focused catalog/reconcile/commands/CLI tests passed; `go test ./...`
  passed; all 35 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-five review: rejected with two high and one medium finding. On a case-insensitive macOS
  volume, a case-variant destination named the same existing plan but bypassed string comparison;
  a nonexistent destination with a formal discovery, implementation, or qualifying family name
  could be created as inventory YAML without `--yes`; and malformed enumerated plans lost matching
  register evidence because reconciliation still used only successfully parsed records. All prior
  migration, rollback, branch, coverage, format, and inventory remediations remained effective.
- Additional remediation: protected existing destinations now compare `os.SameFile` identities
  before the missing-target path fallback. Repository-relative prospective destinations are
  classified with the same formal-path rules as discovery and migration, preventing creation of
  discovery, implementation, qualifying-family README, or register authority. Legacy register
  reconciliation now uses the complete candidate set while views remain based on parsed records.
  Tests cover case-variant macOS identity when the current volume is case-insensitive, all three
  prospective formal kinds, and matching register evidence on a malformed record.
- Post-round-five validation: focused catalog/reconcile/commands/CLI tests passed; `go test ./...`
  passed; all 35 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-six review: rejected with two high and one medium finding after confirming every earlier
  fix. Formal source discovery and register reading could dereference symlinks or block on other
  non-regular files; mixed-mode grandfathering trusted mutable current register entries without
  proof that the plan and evidence predated cutover; and the shared formal-path classifier allowed
  the root plans `README.md` as a family even though enumeration explicitly excludes it.
- Additional remediation: plan and register reads now require a regular `Lstat` identity, open the
  verified inode, compare `os.SameFile`, and read only from that descriptor. Unsafe sources receive
  stable errors and inventory records with explicit unsafe coverage and no external hash/content.
  Mixed mode now requires `plan_catalog_cutover_revision`; validation reads the frozen register
  from that exact Git commit, requires current register bytes to remain identical, and grandfathers
  only register-linked plan paths that also existed at the baseline commit. The config schema,
  fixture adoption flow, implementation contract, and manifest graph digest were updated. Root
  plans `README.md` is excluded consistently from family classification and separately protected as
  planning authority. Tests cover external plan/register symlinks, post-cutover committed legacy
  authoring plus register append, cutover-revision shape, and root-index classification.
- Post-round-six validation: focused catalog/reconcile/commands/CLI tests passed; `go test ./...`
  passed; all 36 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-seven review: rejected with four high and one medium finding. Migration planning did not
  consume the same mixed-cutover/register invariants as validation; the safe reader verified only
  the final path component while migration retained a separate raw read; missing register,
  non-commit, and non-ancestor cutover cases were not all frozen-authority failures; rollback still
  checked a target before separately deleting/restoring it; and root README exclusion remained
  case-sensitive in several paths.
- Additional remediation: migration now plans from the full repository snapshot, promotes every
  non-repairable snapshot error to a transaction blocker, and permits a metadata gap only when a
  resulting action repairs that exact path. One root-relative reader validates containment, every
  parent component, final regular-file identity, and the opened descriptor; discovery, inventory,
  register validation, and migration all use it. Cutover validation requires an exact commit that
  is an ancestor of HEAD, predates mixed mode, and contains regular register/plan blobs; missing or
  unsafe current register state is a frozen-authority error. Rollback atomically quarantines the
  actual current target, verifies those moved bytes, and restores backups through hard-link
  no-replace operations, preserving concurrent targets plus recovery evidence. Root index
  recognition and prospective protection are case-insensitive. Tests cover migration with a
  modified frozen register, symlinked parents, FIFO plans, identity replacement, missing/non-commit/
  non-ancestor/non-legacy baselines, rollback recreation after quarantine, and lowercase root
  index behavior.
- Post-round-seven validation: focused catalog/reconcile/commands/CLI tests passed; `go test ./...`
  passed; all 36 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-eight review: rejected with one high data-safety finding after confirming every round-seven
  remediation. Rollback retained only raw target paths, so replacing an already committed target's
  parent with a symlink before a later failure could redirect quarantine and restoration outside
  the repository.
- Additional remediation: every mutated action now retains the exact parent identity recorded at
  commit, and every installed file retains its inode identity. Rollback refuses a changed parent
  before touching the target, verifies that the current and quarantined entries are the exact
  installed inode, and revalidates the parent again before no-replace backup restoration. An
  adversarial regression replaces the parent with an external symlink after the first action,
  triggers the second action's failure, and proves external bytes remain untouched while the
  committed file, backup, marker, and transaction evidence remain for recovery.
- Post-round-eight validation: focused reconcile/catalog/commands race tests passed; `go test ./...`
  passed; all 36 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-nine review: rejected with two high containment findings. Parent identity checks remained
  separate from path-based rename/link/remove operations, leaving a narrow check/use race, and the
  transaction staging/recovery tree itself could follow a pre-existing `.codeheart/local` symlink
  outside the repository. The reviewer required directory-handle-relative mutation rather than
  another timing check.
- Additional remediation: the transaction layer now opens an `os.Root` for the canonical
  repository, a bound transaction root, and a bound target-parent root per action. Staging,
  marker, backup, quarantine, restore, cleanup, and stale-recovery operations use root-relative
  APIs; rollback moves and restores only inside the opened parent handle and rechecks both the
  path identity and bound handle identity. Backup restoration is copy-with-`O_EXCL`, and exact
  original/installed inode identities remain recorded. Tests now replace a parent during the
  rollback quarantine hook, prove the external target remains unchanged and recovery evidence is
  retained, and prove an external `.codeheart/local` symlink receives no transaction bytes.
- Necessary implementation divergence: the source-build floor changed from Go 1.23 to Go 1.25,
  the first standard-library version containing the required cross-platform `os.Root` mutation
  APIs. Released binary consumers are unaffected; source builders and release notes must disclose
  the new floor. This avoids a platform-specific syscall layer or another security dependency.
- Post-round-nine validation: focused and race-enabled reconcile/catalog/commands tests passed;
  `go test ./...` passed; Windows amd64 reconcile tests cross-compiled; all 36 focused JSON-schema
  tests passed; `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and
  `git diff --check` passed.
- Round-ten review: rejected with three high concurrency findings. Backup restoration reopened a
  mutable backup pathname rather than a bound descriptor; quarantined installed bytes were not
  rechecked after the rollback hook; and the newly opened target-parent handle was not compared to
  the pre-transaction parent identity, allowing an in-repository symlink redirection.
- Additional remediation: backup creation now returns and retains the exact open descriptor,
  inode identity, mode, and SHA-256 of copied original bytes. Rollback requires the transaction
  path to retain that identity, reads from the bound descriptor, revalidates its hash, and restores
  only verified in-memory bytes with `O_EXCL`. The quarantine inode and bytes are rechecked after
  every hook; a changed quarantine is restored as the concurrent current target and retained as a
  recovery-required outcome. Parent binding now receives the pre-transaction nearest-parent
  identity, revalidates it before and after rooted creation/open, rejects any remaining symlink,
  and compares the bound directory handle with the intended current path before mutation. Tests
  deterministically substitute backup bytes, modify the same quarantined inode after the hook, and
  redirect a parent to an in-repository symlink; all preserve concurrent evidence without applying
  substituted bytes or writing to the redirected directory.
- Post-round-ten validation: `go test -race ./...` passed; Windows amd64 reconcile tests
  cross-compiled; all 36 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-eleven review: rejected with three high data-safety findings. Backup-path substitution
  prevented descriptor restore and allowed the unlinked original-byte descriptor to disappear on
  close; an in-repository `.codeheart/local` symlink could redirect recursive cleanup into
  unrelated repository data; and lifecycle `create` actions could replace a file that appeared
  after planning because create/replace/remove preconditions were not encoded consistently.
- Additional remediation: a substituted backup path is now recorded as an anomaly rather than
  making the bound descriptor unusable. Rollback verifies the descriptor identity/hash, restores
  original bytes from verified memory, preserves both the substituted backup path and installed
  quarantine, and returns recovery-required. Transaction setup rejects internal and external
  symlink components. Cleanup enumerates and removes only entries through the bound transaction
  root, removes the transaction directory non-recursively only when its recorded identity still
  matches, and never recursively follows a mutable repository path. Lifecycle planning now hashes
  every replace and remove source; Apply rejects missing SHA-256 preconditions and a create action
  fails if its target appears after planning. Tests prove original-byte survival after backup
  substitution, sentinel preservation behind an internal transaction symlink, and a real
  BuildPlan/Apply init preserving user config bytes created at the pre-commit boundary.
- Post-round-eleven validation: `go test -race ./...` passed; Windows amd64 reconcile tests
  cross-compiled; all 36 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-twelve review: rejected with one high finding after every round-eleven regression passed.
  Apply reused a pre-existing deterministic transaction-ID directory and bound cleanup then
  recursively removed an unrelated sentinel already inside that directory.
- Additional remediation: Apply now creates only the shared transaction parent and atomically
  acquires the exact transaction-ID directory with exclusive `Mkdir`. Any existing directory
  produces a stable `transaction_path_exists` blocker, retains its complete contents, removes only
  the separately acquired marker, and performs no transaction-directory cleanup. Only after
  exclusive acquisition does Apply bind the transaction root and create fresh `stage` and `backup`
  children. A regression precreates the exact directory and sentinel, then proves byte-identical
  preservation, zero target writes, and marker cleanup.
- Post-round-twelve validation: `go test -race ./...` passed; Windows amd64 reconcile tests
  cross-compiled; all 36 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed.
- Round-thirteen review: rejected with three high transaction-authority findings and one medium
  cutover-validation finding. Apply could implicitly stale-recover an unowned exact directory from
  a dead planning marker; pre-commit failure discarded the bound transaction handle and reopened a
  substituted path; phase updates reopened a substitutable marker path; and historical baseline
  config was decoded but not schema-validated.
- Additional remediation: Apply no longer performs implicit stale recovery; any existing marker
  blocks until explicit check/repair recovery. Pre-commit rollback carries the already-bound
  transaction root and original identity through every staging, validation, and marker failure,
  while unbound cleanup remains limited to a newly and safely resolved path. Marker acquisition
  retains its exact descriptor and inode; every phase update verifies path and descriptor identity,
  writes through that descriptor, rechecks the path, and removal refuses substitutions. Tests
  preserve a sentinel behind a forged dead marker, substitute an unrelated directory during the
  staged hook without cleanup damage, and replace the marker with an internal symlink without
  changing its target. Mixed-cutover loading now requires any historical config entry to be a
  regular Git blob and validates it against the Kit config schema before accepting legacy mode; a
  malformed legacy-looking baseline regression now fails with
  `mixed_cutover_baseline_config_invalid`.
- Post-round-thirteen validation: `go test -race ./...` passed; Windows amd64 reconcile and plan
  catalog suites cross-compiled; all 36 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-fourteen review was requested.
- Round-fourteen review: rejected with three high combined-authority findings. A marker-hook
  failure before exclusive acquisition could clean a newly appearing unowned transaction
  directory; committed rollback ignored marker/cleanup authority failures and could report false
  success; and simultaneous target-parent plus backup-path substitution could leave original bytes
  only in an unlinked descriptor that disappeared on close.
- Additional remediation: unbound pre-acquisition rollback never opens or cleans a transaction
  directory. Once acquired, every failure uses only the retained transaction handle. Committed
  rollback aggregates marker-update, rollback, cleanup, and marker-removal failures; any authority
  failure becomes recovery-required and skips destructive cleanup. When an unsafe parent prevents
  target restoration, rollback reads the verified bound backup descriptor; if its normal path was
  removed or substituted, it writes original bytes to a new exclusive
  `recovery-original/<target>.original` file through the bound transaction root before the
  descriptor may close. Regressions prove a marker-hook-created sentinel survives, a post-check
  marker substitution cannot produce false rollback success or alter its target, and combined
  parent/backup substitution retains original bytes durably without touching the external parent.
- Post-round-fourteen validation: `go test -race ./...` passed; Windows amd64 reconcile and plan
  catalog suites cross-compiled; all 36 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-fifteen review was requested.
- Round-fifteen review: rejected with two high evidence-preservation findings after all existing
  validation remained green. Committed rollback could empty its bound transaction directory
  before detecting a marker or transaction-path authority substitution made during rollback, and
  the deterministic `recovery-original/<target>.original` path could be pre-created so the only
  surviving descriptor for original bytes was later closed without durable recovery evidence.
- Additional remediation: committed and pre-commit rollback now validate both marker and bound
  transaction identities before any evidence cleanup. Committed rollback removes the verified
  marker before transaction deletion and refuses transaction cleanup when marker authority is
  lost; bound cleanup validates transaction authority before enumeration, before every deletion,
  and immediately before closing the handle. Original-byte preservation now reserves an
  unpredictable 128-bit-suffixed path with exclusive create and bounded collision retries, verifies
  its opened/path identity and exact bytes before and after close, and synchronizes its containing
  directory. Regressions substitute the marker and transaction path from the rollback quarantine
  hook while returning success, prove both forms retain transaction evidence, and pre-create the
  former deterministic recovery path while proving the verified original survives in a separate
  random file.
- Post-round-fifteen validation: focused reconcile tests passed; `go test -race ./...` passed;
  Windows amd64 reconcile and plan catalog suites cross-compiled; all 36 focused JSON-schema tests
  passed; `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check`
  passed. A fresh round-sixteen review was requested.
- Round-sixteen review: rejected with two high residual concurrency findings. Marker removal and
  bound-transaction cleanup still separated identity validation from pathname deletion, leaving a
  substitution window after validation. Recovery-copy verification checked exact bytes only
  through the original descriptor before close and did not reopen and recheck the inode and bytes
  after synchronizing the containing directory.
- Additional remediation: marker and transaction cleanup now rename the current entry atomically
  to an unpredictable 128-bit-suffixed sibling before deletion, verify that the captured entry is
  the exact bound inode, and retain rather than delete any mismatched capture. Transaction cleanup
  uses the verified quarantine name for every subsequent rooted check and deletion. Deterministic
  hook phases substitute the marker and transaction path after the preliminary validation and
  prove unrelated bytes plus transaction evidence survive. Recovery copies now close, synchronize
  their containing directory, reopen through the bound root, recheck the exact inode and bytes,
  and perform a final path-identity check. Authority-change failures trigger a fresh unpredictable
  exclusive-name attempt while the verified source backup descriptor remains open.
- Post-round-sixteen validation: `go test -race ./...` passed; Windows amd64 reconcile and plan
  catalog suites cross-compiled; all 36 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-seventeen review was requested.
- Round-seventeen review: rejected with one high authority finding in the separate explicit stale
  recovery route. That route still read the marker and `plan.json` through mutable paths, later
  reopened the transaction directory for cleanup, and removed the marker by its current pathname,
  bypassing the bound-descriptor and atomic-quarantine guarantees used by normal Apply cleanup.
- Additional remediation: explicit stale recovery now opens and retains a verified regular marker
  descriptor and exact marker bytes, binds the exact transaction root before semantic validation,
  reads and revalidates `plan.json` through that root, jointly rechecks marker bytes plus marker and
  transaction identities, atomically removes the captured verified marker first, and cleans only
  the retained transaction root through the quarantine protocol. The former path-reopening cleanup
  helper was removed. Deterministic stale-recovery regressions substitute either the marker or the
  transaction directory after preliminary validation and prove mismatched captures, unrelated
  bytes, and original transaction evidence are retained.
- Additional finalization hardening: a successful committed operation now also jointly validates
  marker and transaction authority and removes the verified marker before transaction deletion;
  any cleanup-authority loss is reported as recovery-required rather than a success warning. Tests
  cover both post-validation substitutions after a successful commit. Recovery-copy testing now
  mutates the first random copy after directory sync and proves a fresh retry preserves one exact
  original while retaining the concurrent evidence.
- Post-round-seventeen validation: `go test -race ./...` passed; Windows amd64 reconcile and plan
  catalog suites cross-compiled; all 36 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-eighteen review was requested.
- Round-eighteen review: rejected with two high residual findings. Normal source backup and
  installed-target rollback still selected deterministic stash names, so a concurrent file could
  be replaced by Unix rename between the existence check and capture. Explicit stale recovery
  verified marker and plan bytes before, but not after, the cleanup hooks, allowing a same-inode
  post-validation edit to be deleted as if it were the reviewed evidence.
- Additional remediation: original-source and installed-target captures now use unpredictable
  128-bit-suffixed sibling names and verify the captured inode before proceeding; the former
  deterministic check-then-rename names were removed. A regression pre-creates both former names
  during commit and proves their bytes remain untouched through rollback. Stale marker removal now
  revalidates exact expected marker bytes through its retained descriptor after atomic quarantine,
  and stale transaction cleanup revalidates the retained `plan.json` descriptor, inode, and bytes
  after transaction quarantine but before any entry deletion. Regressions mutate each same inode
  at the post-validation hook and prove the complete mismatched capture is retained.
- Ownership hardening: any stale transaction directory now requires a bound matching `plan.json`,
  including `planning` and `failed` phases. A marker without a transaction directory may still be
  removed when its dead pre-commit identity is verified; an existing unproved directory is never
  inferred as owned from a marker alone. Marker acquisition error cleanup also uses verified
  atomic quarantine rather than pathname removal.
- Post-round-eighteen validation: `go test -race ./...` passed; Windows amd64 reconcile and plan
  catalog suites cross-compiled; all 36 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-nineteen review was requested.
- Round-nineteen review: rejected with one critical containment finding. Explicit stale recovery
  accepted any non-empty transaction ID; a traversal-bearing value could normalize the derived
  recovery path outside `.codeheart/local/kit-transactions` while remaining inside the repository,
  and planted matching marker/plan bytes could then misclassify another repository directory as a
  transaction owned by Apply.
- Additional remediation: explicit recovery now requires the exact current Apply transaction-ID
  grammar of 32 lowercase hexadecimal characters before deriving or opening a path. The recovery
  path must independently equal exactly one direct child of
  `.codeheart/local/kit-transactions`, and the bound plan must declare schema version 1 plus the
  same valid transaction identity, canonical root, and command. A traversal regression plants a
  matching malicious marker and `plan.json` under `docs/repo/plans`, then proves the marker, plan,
  and sentinel remain byte-identical when recovery rejects the identity before opening the target
  directory.
- Post-round-nineteen validation: `go test -race ./...` passed; Windows amd64 reconcile and plan
  catalog suites cross-compiled; all 36 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-twenty review was requested.
- Round-twenty review: the critical stale-recovery remediation passed, but the only separate
  non-transactional write path was rejected with one high containment finding and one medium
  platform finding. Inventory output validated mutable parent pathnames before later `MkdirAll`,
  temporary creation, and rename, allowing a parent substitution to redirect the write; final
  replacement via `os.Rename` also lacked reliable repeat-refresh semantics on Windows.
- Additional remediation: inventory output now resolves, opens, retains, and verifies the exact
  destination-parent handle; creates an unpredictable exclusive temporary file inside that bound
  parent; revalidates the requested parent path before commit; atomically quarantines and verifies
  any prior regular output; installs the new inode with an exclusive same-directory hard link;
  verifies installed identity and exact bytes; and synchronizes the directory. Concurrent targets
  are never overwritten, mismatched captures are retained, and a post-commit parent-authority loss
  is reported. A deterministic hook replaces the requested parent with a symlink to a formal plan
  directory after binding and proves the plan remains byte-identical. The normal command test now
  refreshes the same output twice, and the complete commands suite cross-compiles for Windows.
- Post-round-twenty validation: `go test -race ./...` passed; Windows amd64 reconcile, plan catalog,
  and commands suites cross-compiled; all 36 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-twenty-one review was requested.
- Round-twenty-one review: rejected with one high residual parent-creation race. The new writer
  bound an existing destination parent, but still ran `MkdirAll` first when the parent was absent;
  a missing ancestor could be concurrently replaced before binding and receive unauthorized
  directory creation. Existing-parent redirection, repeated refresh, concurrent target, and
  same-inode output-change tests otherwise passed.
- Simplifying remediation: inventory output now requires its explicit destination directory to
  exist and never creates output ancestors. It verifies that parent as a non-symlink directory
  before resolving and binding it. This removes the directory-creation authority problem without
  adding a second handle-relative directory-construction subsystem. A regression targets a nested
  missing output directory, requires `inventory_target_parent_missing`, and proves no ancestor or
  file appears. Existing explicit-file creation and repeat refresh UX remain unchanged.
- Additional concurrency hardening: when refreshing an existing inventory, the writer now retains
  the old output descriptor and exact bytes from binding through quarantine. A same-inode edit is
  retained and reported rather than overwritten; a concurrently appearing previously missing
  target makes the exclusive install fail and remains byte-identical.
- Post-round-twenty-one validation: `go test -race ./...` passed; Windows amd64 reconcile, plan
  catalog, and commands suites cross-compiled; all 36 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-twenty-two review was requested.
- Round-twenty-two review: rejected with two high inventory concurrency findings. Protection
  checks still ran before the eventual destination-parent descriptor was bound, allowing a
  regular-directory exchange to change which path received the protected-path decision. The
  writer also verified the quarantined old output only immediately after quarantine and later
  removed it after an identity-only check, so a same-inode edit made during that interval could be
  deleted.
- Additional remediation: the writer now requires, resolves, opens, and validates the explicit
  destination parent before any protected-path decision, invokes the deterministic race hook only
  after that binding, and revalidates the requested path against the bound directory before
  continuing. A regular-directory exchange regression moves the real formal-plan directory into
  the requested pathname during that hook and proves validation rejects the changed authority
  before any output file is created or plan byte is changed.
- Quarantine remediation: the original output descriptor and exact captured bytes now remain live
  through the final quarantine cleanup. Immediately before deletion, the writer rechecks the
  quarantine pathname identity, descriptor identity, and exact bytes; any mismatch retains the
  quarantine as evidence and reports a stable error. A same-inode post-quarantine mutation
  regression proves the concurrently changed bytes survive.
- Post-round-twenty-two validation: `go test -race ./...` passed; Windows amd64 reconcile, plan
  catalog, and commands suites cross-compiled; all 36 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-twenty-three review was requested.
- Round-twenty-three review: accepted with no material findings. It confirmed both round-twenty-two
  race fixes, passed the focused EP-02 race suites and Windows amd64 cross-compilation, and found no
  remaining correctness, containment, concurrency, schema, command, migration, or scope blocker.

## EP-03 Delta - Config-Driven Portfolio Discovery And Branch-Aware Scanning

Status: completed and accepted in review round fourteen.

- Added normalized portfolio-v1/v2 config loading, committed GitHub-owner sources, machine-local
  Git-root sources, exact default-branch membership, and automatic coordination-home
  self-membership. Default-branch Kit marker, schema-valid lock, v2 role, stable repository ID, and
  matching home ID are all required; branch-only enrollment remains a candidate.
- Added provider-neutral discovery contracts, direct-argument cancellable Git and `gh` runners,
  bounded local-root discovery, authenticated paginated GitHub discovery, three-attempt transient
  retry, distinct permission and rate-limit blockers, truncation detection, redacted diagnostics,
  and credential-bearing remote rejection.
- Added scanner-owned fresh bare-mirror replacement under `.codeheart/local/portfolio/git/`.
  Mirrors fetch and prune only `refs/remotes/origin/*`, disable repository hooks, read inert Git
  objects, and never fetch into or read bytes from a developer checkout. A single identity-bound
  scan lock serializes mirror and cache mutation.
- Added complete default-branch canonical baselines plus every accessible unmerged remote branch.
  Branch overlays use merge-base changed paths, omit inherited plans and deletions, retain
  independent same-ID observations, expose conflicts, use file-specific Git history for the
  30-day stale label, and retire merged or deleted refs on the next complete scan. Pull-request
  data remains optional and does not select branches.
- Added schema-valid deterministic factual catalogs with members, candidates, observations,
  errors, metrics, completeness, last-complete freshness, and explicit local-only evidence
  omission. Successful complete scans atomically replace the local cache; provider, mirror,
  required-member, comparison, or cancellation failures preserve the prior complete bytes.
  Strategic-overlay schema, identity, inode, and exact bytes remain unchanged across scans.
- Added `portfolio configure` and `portfolio scan`, repeatable `--github-owner` and `--local-root`
  parsing, matching-value idempotency, conflicting-identity refusal, role-specific placement,
  absent-only portfolio scaffolds, and grouped CLI help. Added shared-engine
  `plans validate --remote-overlays` and `plans inventory --remote-overlays` with pushed-only
  provenance disclosure.
- Real Git fixtures prove one home self-member, an exactly enrolled member, provider-visible
  exclusion, branch-only enrollment refusal, baseline de-duplication, independent changed-plan
  overlays, conflict retention, staleness, deleted-ref retirement, untouched developer refs and
  worktree, interrupted-cache preservation, failed-member preservation, strategic-overlay
  immutability, and exclusion of a committed unpushed local branch. Fake `gh` fixtures prove full
  pagination, retry, rate-limit, permission, truncation, cancellation, and secret redaction.
- Validation: `go test -race ./...` passed. The portfolio, commands, and CLI suites cross-compiled
  for Windows amd64, and the Windows CLI binary built successfully. Existing macOS and Windows
  workflow jobs both execute `go test ./...`, so the multi-repository fixture scan is included on
  both real platforms without a workflow change. All 36 focused JSON-schema tests,
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-one EP-03 review was requested.
- Round-one review: rejected after focused race tests and vet passed. The reviewer confirmed that
  unreadable default-branch config, lock, or marker objects were silently treated as missing and
  could falsely retire a required member in a complete cache; mirror-root creation could follow a
  `.codeheart/local` symlink outside repository authority; equivalent GitHub SSH and HTTPS remotes
  were not deduplicated, which could duplicate the coordination home's self-member; member role
  configuration could still construct coordination discovery sources; and missing file-specific
  history silently weakened staleness evidence. The extended audit was interrupted after those
  material findings were established so remediation could begin.
- Round-one remediation: membership object failures now emit
  `membership_evidence_unavailable` and make the scan incomplete while genuine missing files
  remain candidate evidence. Plan-specific history failures emit
  `observation_freshness_unavailable` instead of silently claiming complete freshness. Remote
  identities canonicalize GitHub HTTPS, SSH-URL, and SCP forms before source deduplication, with
  self-source precedence. The config schema forbids committed discovery on v2 members, member
  loading ignores machine-local source files, and runtime source construction is home-only.
- Containment remediation: mirror roots and cache parents are created one component at a time
  beneath retained `os.Root` handles with non-symlink directory checks. Refresh and installed
  mirror directories remain identity-bound; every later Git object read verifies the retained root
  and mirror handles before and after the subprocess. Regression tests cover local symlink escape,
  equivalent-remote self deduplication, member source isolation, unreadable membership evidence,
  and unavailable plan-specific freshness history.
- Post-round-one validation: `go test -race ./...` passed. The portfolio, commands, and CLI suites
  and the CLI binary cross-compiled for Windows amd64 to explicit temporary outputs. All 37
  focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed. The embedded profile graph digest was
  refreshed after the role-specific config-schema change. A fresh round-two review was requested.
- Round-two review: rejected with two high findings after confirming the named round-one
  remediations. First, readable but malformed config or lock data still became candidate output,
  and every failed home self-membership could still leave a scan falsely complete. Second, cache
  publication retained a parent handle without revalidating its canonical namespace, removed the
  old target before linking the new target, and allowed unbound symlink-following cache reads.
- Round-two membership remediation: membership decisions now distinguish absent/unconfigured
  evidence from malformed evidence. Malformed config or lock data is an incomplete-scan blocker;
  any failed coordination-home self-membership is also a blocker. Regressions seed a complete
  cache, corrupt readable member config and lock evidence, and change the home remote to a member
  role; every attempt remains incomplete and preserves the prior cache byte-for-byte.
- Round-two publication remediation: cache writing binds and repeatedly verifies both the
  repository-relative canonical parent and its retained directory identity. Parent-authority loss
  preserves the staged evidence and stops publication. The old remove/quarantine/link sequence was
  replaced by one same-directory `os.Root.Rename` from the fully synced staged file to the
  canonical target, so concurrent readers see either complete old bytes or complete new bytes.
  Existing and installed target identities, retained descriptors, and exact bytes are verified;
  cache reads now use a contained regular-file handle and reject symlinks. Regressions cover parent
  substitution, symlink reads, and concurrent readers across publication.
- Post-round-two validation: `go test -race ./...` passed. The portfolio, commands, and CLI suites
  and the CLI binary cross-compiled for Windows amd64 to explicit temporary outputs. All 37
  focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed. A fresh round-three review was requested.
- Round-three review: rejected with two high contract gaps after confirming the round-two
  membership and atomic-replacement changes. Intermediate in-repository symlink components could
  still redirect cache reads, and the scan lock lived inside the same replaceable cache subtree
  without canonical namespace revalidation. Separately, the approved optional pull-request
  enrichment capability and fake-`gh` coverage had not been implemented.
- Round-three namespace remediation: config, local-source, cache, and lock reads now share a
  retained bound-directory abstraction. It validates the repository-root identity, rejects every
  intermediate symlink component, binds the exact parent directory, and revalidates canonical and
  retained identities before and after reads or writes. Exact read bytes are rechecked through the
  retained descriptor. The scan lock moved above the replaceable `.codeheart/` subtree to the
  repository root and is revalidated after repository scans and around cache publication. Tests
  replace both the cache parent and the entire `.codeheart/` namespace while the lock is held and
  prove a second scan cannot acquire authority or mutate replacement bytes. An in-repository
  intermediate-symlink cache read is rejected.
- Round-three enrichment remediation: `RepositorySource` now carries an optional provider-neutral
  ref enricher. The GitHub adapter performs bounded, completely paginated open-PR lookup through
  direct `gh api` arguments, validates public GitHub PR URLs, correlates facts to refs the scanner
  already selected, and contributes API-call metrics. Enrichment failures remain optional and do
  not select or suppress branches. Duplicate SSH/HTTPS home sources retain the provider enricher
  while preserving the home self-source. Tests prove two-page PR lookup, matching-ref enrichment,
  and inclusion of a changed remote branch with no PR fact.
- Additional authority hardening: an overlay that appears as a dangling symlink during a scan is
  now detected with `Lstat` rather than treated as absent, and cache/config read descriptors verify
  exact bytes after namespace checks.
- Post-round-three validation: `go test -race ./...` passed. The portfolio, commands, and CLI
  suites and the CLI binary cross-compiled for Windows amd64 to explicit temporary outputs. All 37
  focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed. A fresh round-four review was requested.
- Round-four review: rejected with two high compatibility/normalization findings and one medium PR
  correlation finding. The compatibility schema permits a path-only v1 coordination home, but the
  loader still required the new stable repository identity. Home self-discovery also passed a
  relative origin URL directly to a fresh mirror, resolving it relative to the wrong directory.
  Finally, matching PR facts by branch name alone could attach a fork PR to an origin branch with
  the same head-ref name.
- Round-four remediation: v1 path-only coordination homes now load as compatibility records with
  empty v2 identities, while remote scanning continues to require the v2 home identity. A static
  v1-home fixture and Go/Python regressions cover that boundary. Self-discovery normalizes relative
  origins against the coordination-home repository exactly like local-source discovery; an
  end-to-end relative-origin fixture produces one complete self-member baseline. GitHub PR
  enrichment now captures `head.repo.full_name` and requires it to match the scanned repository;
  a colliding fork head ref remains unenriched while the same-repository PR is retained.
- Post-round-four validation: `go test -race ./...` passed. The portfolio, commands, and CLI suites
  and the CLI binary cross-compiled for Windows amd64 to explicit temporary outputs. All 38
  focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed. A fresh round-five review was requested.
- Round-five review: rejected with one critical transport-policy finding and one high branch-
  overlay finding. Arbitrary Git transport helpers such as `ext::` were not rejected, so remote
  helper commands or unsafe URL rewrites could escape the inert-content boundary. Separately,
  rename detection treated every renamed canonical path as materially changed without comparing
  its merge-base and branch bytes.
- Round-five transport remediation: accepted remote syntax is now limited to credential-free
  HTTPS, SSH with the optional `git` user, SCP-style `git@host:path`, `file://`, and absolute file
  paths. Every scanner-owned Git command receives an explicit deny-by-default protocol policy with
  only HTTPS, SSH, and file enabled and `ext` denied; remote operations additionally pin the exact
  configured URL against global `insteadOf` rewrites. Command-scope protocol/config injection and
  SSH/proxy override environment variables are removed before Git starts, while hooks remain
  disabled. Regressions reject `ext::`, custom helpers, Git and plain HTTP transports; confirm all
  supported forms and policy arguments; and install a malicious global URL rewrite whose helper
  sentinel must remain absent.
- Round-five overlay remediation: changed-path enumeration now consumes rename-aware name/status
  evidence. For every rename it reads the merge-base source and branch destination objects and
  emits the new path only when canonical bytes differ. A real multi-branch Git regression proves a
  pure rename contributes no overlay while a rename with changed canonical bytes contributes one.
- Post-round-five validation: `go test -race ./...` passed. The portfolio, commands, and CLI suites
  and the CLI binary cross-compiled for Windows amd64 to explicit temporary outputs. All 38
  focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed. A fresh round-six review is pending.
- Round-six review: rejected with one critical command-scope injection finding and one high
  command-policy consistency finding. The environment sanitizer removed a fixed list of Git
  variables but still admitted `GIT_CONFIG_PARAMETERS`, allowing a caller to inject an alternate
  `remote.origin.uploadpack`. Local and self-origin discovery also invoked Git outside the central
  protocol and hook policy, so the scanner did not yet have one auditable safety boundary for all
  of its Git subprocesses.
- Round-six remediation: the scanner now removes every inherited `GIT_*` variable and restores
  only non-interactive prompting. All scanner-owned Git calls, including local and self discovery
  plus inert-object reads, pass through the central deny-by-default protocol and disabled-hook
  policy. Remote operations additionally pin the exact URL identity and the built-in
  `git-upload-pack` service. A thread-safe recording runner proves every scanner-owned invocation
  carries the central policy; a real command-scope upload-pack injection fixture proves the raw
  attack is active before the safe mirror fetch and that its helper sentinel remains absent.
  Mode-only branch changes are also byte-compared and suppressed alongside pure renames.
- Post-round-six validation: the targeted transport, policy-coverage, and rename tests and the
  complete `go test -race ./...` suite passed. The portfolio, commands, and CLI suites and the CLI
  binary cross-compiled for Windows amd64 to explicit temporary outputs. All 38 focused
  JSON-schema tests passed; `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and
  `git diff --check` passed. A fresh round-seven review was requested.
- Round-seven review: rejected with one high containment finding. Refresh bound the random
  staging parent but repeatedly addressed its nested `mirror.git` directory by pathname without
  binding that child. A concurrent child substitution between initialization and later remote or
  fetch commands could therefore redirect scanner-owned writes outside `.codeheart/local/`.
- Round-seven remediation: refresh now creates and binds the empty `mirror.git` directory before
  initialization, verifies both the staging parent and child identities before and after every
  refresh-time Git subprocess, rejects a symlink or replacement child, and releases the child
  handle only before the atomic install rename. An adversarial runner substitutes the initialized
  child with an external bare-repository symlink; the next command is refused, only initialization
  runs, and the external repository remains byte-identical.
- Post-round-seven validation: the focused child-substitution, transport, central-policy, and
  rename tests and the complete `go test -race ./...` suite passed. The portfolio, commands, and
  CLI suites and the CLI binary cross-compiled for Windows amd64 to explicit temporary outputs.
  All 38 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed. A fresh round-eight review was requested.
- Round-eight review: rejected with one critical Windows environment finding and two high
  descriptor-lifetime findings. Case-sensitive filtering admitted mixed-case `GIT_*` variables on
  Windows. Refresh checks still surrounded a subprocess that resolved a mutable pathname, leaving
  a check/use window, and the refreshed child identity was released before installation without
  comparing the installed directory to the refreshed object.
- Round-eight remediation: Git environment keys are now case-folded before every inherited
  `GIT_*` variable is removed; a mixed-case sanitizer fixture and Windows-only live upload-pack
  sentinel cover the platform behavior. On Unix, bound Git execution passes the retained
  directory descriptor into a helper process, performs `fchdir`, and then replaces that process
  with Git, so Git never selects the repository through the mutable path. On Windows, each command
  opens and identity-matches a non-delete-sharing directory handle before setting the command
  directory, preventing rename or deletion until Git exits. Both refresh mutations and installed
  object reads use the bound command path. The exact refreshed root and file descriptors remain
  open across the install rename; the installed namespace is compared with both, the original
  file descriptor is transferred into the installed repository object, mismatches are rejected,
  and the previous mirror is restored. Tests substitute the path immediately before delegation
  and replace the staged identity before installation; the external target remains unchanged and
  a prior installed mirror is restored byte-for-byte. An initial round-nine review was paused
  without a decision before this stronger descriptor-lifetime amendment, so the fresh gate reviews
  the final behavior rather than a stale snapshot.
- Post-round-eight validation: the complete portfolio race suite and `go test -race ./...` passed.
  The portfolio, commands, and CLI suites and CLI binary cross-compiled for Windows amd64 to
  explicit temporary outputs. All 38 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-nine review was requested.
- Round-nine review: rejected with one high late-rollback finding after confirming every named
  round-eight remediation. Two retained-authority failure exits did not remove the suspect
  installation or restore a quarantined prior mirror, while earlier rollback paths ignored their
  own restoration failures.
- Round-nine remediation: every failure after prior-mirror quarantine now passes through one
  checked rollback helper. It removes an installed suspect, restores the quarantined prior mirror,
  and promotes any cleanup or restoration failure into an explicit `mirror_rollback_failed`
  diagnostic rather than silently returning. Fault injection at the late retained-authority check
  proves that path restores the prior mirror byte-for-byte; the earlier staged-identity
  substitution regression continues to prove the same behavior.
- Post-round-nine validation: the focused substitution, installation rollback, transport, and
  central-policy race tests and complete `go test -race ./...` suite passed. The portfolio,
  commands, and CLI suites and CLI binary cross-compiled for Windows amd64 to explicit temporary
  outputs. All 38 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed. A fresh round-ten review was requested.
- Round-ten review: rejected with two high name-authority/unchecked-error findings and one medium
  proof gap. Rollback deleted and restored mutable sibling names without retaining the prior
  mirror identity, several handle-close and quarantine-cleanup errors were ignored, and the
  byte-identical restoration regression compared only the Git config file.
- Round-ten remediation: refresh now retains both root and file descriptors for the prior mirror
  before quarantine and snapshots its complete tree deterministically, including paths, modes,
  symlink targets, file lengths, and SHA-256 content digests. Rollback verifies the canonical entry
  against retained new-mirror descriptors before deletion, verifies the quarantine against
  retained prior descriptors before restoration, revalidates the restored namespace and complete
  tree, and refuses to delete or restore substituted names. Every final-handle close, prior-handle
  close, and quarantine cleanup result is checked and surfaced. Canonical namespace authority is
  revalidated after the late retained-authority hook as well as before cleanup. Tests prove
  full-tree restoration after retained-authority and cleanup-precondition failures; a hostile late
  canonical substitution remains byte-identical and produces `mirror_rollback_failed` instead of
  being deleted as unrelated data. An initial round-eleven pass was paused without a decision so
  this final late-hook ordering amendment receives a fresh review.
- Post-round-ten validation: the complete portfolio race suite and `go test -race ./...` passed.
  The portfolio, commands, and CLI suites and CLI binary cross-compiled for Windows amd64 to
  explicit temporary outputs. All 38 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-eleven review was requested.
- Round-eleven review: rejected with two high late-namespace/TOCTOU findings after confirming the
  round-ten identity, snapshot, and error-checking remediation. The last cleanup hook could replace
  the canonical entry after its latest check, and rollback/cleanup still verified a mutable name
  before performing a later path-based delete or restore.
- Round-eleven remediation: canonical authority is rechecked after the cleanup hook and immediately
  before success. Destructive and restorative operations now atomically rename the current name to
  a fresh scanner-owned capture before verifying the captured entry against the retained root and
  file descriptors. Only a matching capture is eligible for cleanup or restoration. A mismatching
  object is preserved in its capture and reported; it is never passed to deletion. The hostile late
  replacement regression proves its sentinel remains byte-identical in rollback quarantine while
  the retained prior mirror is restored and full-tree verified.
- Post-round-eleven validation: the focused containment tests and complete `go test -race ./...`
  suite passed. The portfolio, commands, and CLI suites and CLI binary cross-compiled for Windows
  amd64 to explicit temporary outputs. All 38 focused JSON-schema tests passed;
  `scripts/validate-json-schemas.py`, `go vet ./internal/...`, and `git diff --check` passed. A
  fresh round-twelve review was requested.
- Round-twelve review: rejected with one high staging-cleanup finding after confirming the new
  canonical/prior capture protocol and all named Git-policy defenses. The random refresh parent
  still used an unchecked deferred path-based `RemoveAll`, so a replacement at that name could be
  deleted and cleanup errors were invisible.
- Round-twelve remediation: refresh now retains both root and file descriptors for the staging
  parent. Deferred cleanup atomically captures the current staging name, verifies the capture
  against both retained descriptors, removes only a match, preserves mismatches in quarantine,
  checks both cleanup and close results, and converts an otherwise successful refresh into an
  explicit failure if safe cleanup cannot be proven. A deterministic hook replaces the bound
  staging parent; its sentinel remains byte-identical in capture and no repository is returned.
- Post-round-twelve validation: the focused staging, child, installation, transport, and central-
  policy race tests and complete `go test -race ./...` suite passed. The portfolio, commands, and
  CLI suites and CLI binary cross-compiled for Windows amd64 to explicit temporary outputs. All 38
  focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed. A fresh round-thirteen review was
  requested.
- Round-thirteen review: rejected with one high atomic-publication finding after confirming the
  round-twelve staging cleanup and all named containment and Git-policy defenses. Both the final
  staged-mirror installation and prior-mirror restoration still used ordinary renames after their
  preceding absence or quarantine decisions. An unexpected entry appearing in either interval
  could therefore be overwritten instead of preserved as unrelated data.
- Round-thirteen remediation: all scanner-owned mirror capture, publication, and restoration
  transitions now use a rooted atomic no-replace primitive. Darwin uses
  `renameatx_np(RENAME_EXCL)`, supported Linux architectures use
  `renameat2(RENAME_NOREPLACE)`, and Windows retains and identity-matches the mirror-root handle
  before calling the destination-preserving `MoveFile` operation. Unsupported platform or Linux
  architecture combinations fail explicitly instead of falling back to a replacing rename.
  Deterministic hooks create unexpected occupied directories in the final install and rollback
  restore windows; both regressions preserve their sentinel bytes and return
  `mirror_rollback_failed` without replacing the occupants.
- Post-round-thirteen validation: the focused staging, child, installation, no-replace rollback,
  transport, and central-policy race tests and complete `go test -race ./...` suite passed. The
  portfolio, commands, and CLI suites and CLI binary cross-compiled for Windows amd64. The
  portfolio package also cross-compiled with `CGO_ENABLED=0` for Linux amd64/arm64 and Darwin
  amd64/arm64. All 38 focused JSON-schema tests passed; `scripts/validate-json-schemas.py`,
  `go vet ./internal/...`, and `git diff --check` passed. A fresh round-fourteen review was
  requested.
- Round-fourteen review: accepted with no material finding. The reviewer independently traced
  every capture, quarantine, install, and restoration through the rooted atomic no-replace
  primitive; confirmed Darwin, Linux amd64/arm64, Windows amd64, and explicit unsupported-target
  behavior; reran focused race, schema, vet, diff, and cross-compilation checks; reconfirmed prior
  membership, cache, namespace, transport, environment, descriptor-bound execution, branch,
  enrichment, and compatibility defenses; and found no accidental `EP-04` work.

## EP-04 Delta - Managed Planning, Coordination, And Publication UX

Status: completed and accepted in review round three.

- Added exact managed references for semantic plan metadata, semantic IDs, record kinds, families,
  relations, aliases, catalog modes, chronology, derived views, observations, portfolio roles,
  exact membership, source ownership, cache completeness, and strategic overlays.
- Replaced manual mixed/canonical register maintenance with a stable entry point and generated CLI
  views while retaining legacy-mode maintenance and frozen v1 evidence. Fresh declarations no
  longer create `coordination-sync-pending.md`; existing consumer-owned copies remain untouched and
  compatibility guidance remains installed.
- Added hybrid coordination setup plus agent-facing refresh and semantic migration recipes. Each
  declares source, inputs, preconditions, lane, approval, stop, evidence, validation, recovery,
  L1/L3 maturity, and the no-L2-wrapper boundary. Setup asks role, repository ID, home ID, GitHub
  scope, and local roots one decision at a time.
- Added the bounded activation contract: an activation or user-requested material active-plan
  update grants one plan/log/direct-metadata branch creation/use, intentional commit, and normal
  push without a duplicate prompt. Static routing checks cover unrelated files/code, ambiguity,
  auth and policy failure, rejected normal push, PR, merge, release, force-push, deletion, history
  rewrite, and destructive-Git exclusions. Execution may continue on the active pushed work branch
  without merge.
- Added root/fallback/direct routes for local plan authoring, portfolio configuration, refresh-
  before-analysis, migration, and activation. Route evidence distinguishes pushed remote facts
  from worktree, local-head, and unpushed evidence and forbids current completeness claims after a
  failed scan.
- Safe added task: because the fresh generic strategic-overlay scaffold cannot know a repository's
  approved home ID, `portfolio configure` now loads the declared scaffold resource and replaces
  only its exact untouched placeholder under optimistic transaction checks. Edited overlays and
  unsafe symlink/non-regular occupants remain preserved. The README is also created from the
  declared packaged scaffold instead of duplicated command text.
- Updated placement and consumer-impact contracts, planning/agent-interface component versions,
  the standard profile, root and repository templates, content-manifest checksums/digest, and every
  changed packaged Python fallback resource.
- Validation: complete `go test ./...` passed after refreshing the intentional profile-hash
  fixture; the focused EP-04 Python set passed 59 tests including release-asset construction,
  packaged fallback, source/package byte equality, routing, init/onboarding, sync preservation,
  Go/Python parity, and public-core rejection. All 38 focused JSON-schema tests and the schema,
  Markdown-header, public-core, release-manifest, and diff validators passed. The
  `go vet ./internal/...` gate and Windows amd64 commands/CLI test-binary cross-compilation passed.
- Round-one review: rejected with two medium documentation/scaffold inconsistencies after every
  executable, routing-probe, manifest, packaging, and public-core gate passed. The setup runbook
  incorrectly said configuration never overwrites an existing overlay even though the command
  intentionally replaces the byte-exact untouched placeholder. The generic portfolio README also
  claimed every fresh repository was already a coordination home despite fresh configuration
  declaring no portfolio role.
- Round-one remediation: the setup runbook now states the exact preservation boundary: README,
  edited overlay, symlink, and non-regular occupants remain unchanged; only the byte-exact pristine
  overlay placeholder is replaced with the approved home ID. The generic README is now role-neutral
  and points to `.codeheart/kit.config.yaml` as role and identity authority. Source and packaged
  resources match byte-for-byte and the content graph digest was refreshed.
- Post-round-one validation: complete `go test ./...` passed; the focused EP-04 Python suite passed
  100 tests; JSON-schema, Markdown-header, public-core, release-manifest, and diff validators passed.
  A fresh round-two review was requested.
- Round-two review: rejected with one medium residual wording inconsistency after confirming the
  two round-one fixes and all five low-context routes. The setup runbook's preview and validation
  language still implied that an overlay changes only when absent and that every existing scaffold
  remains byte-preserved, contradicting its precise pristine-placeholder exception.
- Round-two remediation: the runbook's intent, preview, operator notes, and validation now use one
  exact rule: create missing files, preserve README and repository-owned overlay bytes, and bind
  only the byte-exact untouched overlay placeholder to the approved home ID. A static regression
  asserts this wording, rejects the superseded broad claim, and confirms source/package identity.
- Post-round-two validation: focused manifest, command, and embedded-resource Go tests passed; 12
  focused Python routing, packaging, and init tests passed; Markdown-header, public-core,
  release-manifest, and diff validators passed. A fresh round-three review was requested.
- Round-three review: accepted with no material findings. The reviewer independently confirmed the
  absent-create and pristine-placeholder exception across intent, preview, operator, validation,
  and recovery guidance; role-neutral scaffolding; source/package parity; declarations, manifests,
  content graph, all five low-context routes, the static wording regression, focused Go/Python
  tests, release-manifest, Markdown-header, public-core, vet, and diff gates. No `EP-05` or later
  implementation was reviewed.

## EP-05 Delta - Producer Semantic Migration And Catalog Cutover

Status: active; all 33 local records are canonical in mixed mode, with remote/default-branch
coverage and canonical cutover pending default-branch enrollment.

- Safe sequence correction: the task list originally placed the frozen-register notice after mixed
  adoption. The implemented runtime requires current register bytes to match the exact legacy
  baseline, so the managed migration runbook's corrected order is authoritative: add the notice in
  legacy mode, commit the baseline, then record that commit while entering mixed mode. This changes
  ordering only and preserves the approved capability.
- Repository identity decision: adopted the public stable repository ID
  `codeheart-operating-kit`. Existing V1 coordination-home configuration in the sibling
  `Codeheart-HQ` repository provides the evidence for public-safe home ID `codeheart-hq`; this
  producer remains a member and is not reclassified as the coordination home.
- Preserved-state boundary: the pre-existing untracked `.codeheart/kit.config.yaml` was inspected
  and its consumer-layer values retained. Only explicit legacy catalog settings and portfolio-v2
  identity were added. The untracked lock, installed Kit, local memory, pending-sync work, `uv.lock`,
  and Windows binary remain excluded.
- Pre-baseline validation: source CLI validation found 33 formal records, all valid in legacy mode,
  with expected metadata/title-layout warnings and the documented unpaired `OK-PR-014` evidence.
  JSON schemas, Markdown headers, public-core hygiene, and planning-surface diff checks passed.
- Baseline checkpoint: committed and normally pushed planning/configuration-only commit
  `3b134a07b806d1522e7cf76fe52e2028d79155c2`. It contains the frozen register, valid legacy config,
  active plan, and execution log; source implementation and unrelated untracked files remain
  outside the checkpoint.
- Remote-overlay preflight: the source CLI enumerated all 33 local records but correctly returned
  incomplete remote coverage because `origin`'s default branch does not yet contain portfolio-v2
  enrollment. The failed artifact contained a local source locator, so it was retained only under
  `/private/tmp` and was not committed to this public repository. No membership rule was weakened
  and no current-completeness claim was made.
- Local migration inventory: the committed plan-scoped inventory records exact baseline revision
  `3b134a07b806d1522e7cf76fe52e2028d79155c2`, 21 implementation records, 12 discovery records, no
  qualifying family README, 33 legacy records, and unpaired historical evidence `OK-PR-014`.
- Semantic review: reviewed every canonical document against its own content, sibling records,
  legacy register evidence, lifecycle, source hash, and branch ownership. The ledger assigns 33
  stable semantic IDs, public product/capability/theme classifications, evidence-backed sibling
  relations, 27 unique aliases without duplication, and no unsupported family classification.
  Combined `OK-PR-025` and `OK-PR-026` rows retain their alias only on the implementation record;
  the discoveries remain separately addressable. Plans without a register row receive no invented
  alias.
- Mixed adoption: set `plan_catalog_cutover_revision` to the exact baseline and entered mixed mode.
  Pre-apply validation reported the four expected unregistered legacy records as transient mixed
  blockers; the reviewed migration plan contained guarded writes for all four and no blocker. The
  dry run planned exactly 32 replacements and deferred only the active branch-owned implementation
  plan.
- Guarded apply: `plans migrate --yes --json` transaction
  `117ec086a8a4436a9208a7ed709c65ba` applied all 32 unchanged records with staged-state,
  post-check, transaction-cleanup, and marker-cleanup validation. The current branch owner then
  inserted the reviewed metadata for the active implementation plan and reconciled its ledger
  evidence without changing the metadata decision.
- Local closure evidence: mixed validation is valid for all 33 canonical records. A second dry run
  and apply returned zero changes and 33 `already_applied` outcomes. All 33 IDs are unique; the 27
  assigned aliases are unique. Metadata-only changes preserved every historical `Last updated`
  line byte-for-byte; the active implementation plan's later checklist progress is a separate
  meaningful content change. JSON-schema, Markdown-header, public-core, and diff validation passed.
- Remaining gate: remote overlay validation cannot succeed until the repository's portfolio-v2
  enrollment exists on the actual default branch. Therefore the producer remains mixed and no
  canonical-cutover checkbox is claimed. Default-branch integration is a repository-governance
  action outside the current plan-only push authority.
- Isolated rollout simulation: an exact local clone and local bare remote reproduced the intended
  two-integration sequence without changing GitHub. A merge commit integrating source plus mixed
  enrollment produced 33 canonical default-branch observations, one self-member, no candidates or
  scan errors, and complete remote validation. A synthetic unmerged plan-changing branch added
  exactly one verified branch observation while retaining all 33 baseline observations. Switching
  only the simulated producer config to canonical passed remote validation and JSON listing before
  and after a second merge commit. The final simulated default branch remained canonical with 33
  canonical local records, 34 remote observations, zero invalid/unreadable/unsafe records, no
  dirty overlap, and the frozen register byte-identical to baseline `3b134a0`. This proves the
  mechanics and merge-commit topology but does not substitute for actual GitHub integration or
  current live-remote reconciliation.

## EP-06 Delta - Integrated Validation And Release-Candidate Handoff

Status: active; source-level validation is clean, but immutable candidate handoff cannot close
before the `EP-05` default-branch gate and the later release-candidate authority boundary.

- Cross-platform validation surface: macOS and Windows workflow jobs now run grouped plan catalog,
  migration, portfolio-overlay, command, CLI, routing, schema, and packaged-resource tests. The
  affected Go packages also cross-compiled successfully for Windows amd64 in this worktree. Real
  Windows execution remains release-candidate evidence and is not claimed locally.
- Low-context and activation fixtures: installed-route probes cover semantic reference, member
  authoring, coordination-home setup, refresh-before-analysis, incomplete-refresh disclosure,
  migration skips, and activation boundaries. An automated isolated repository and local bare
  remote fixture proves
  that activation publishes only the selected plan and execution-log checkpoint while unrelated
  code remains local; no PR is created and no publication CLI was added. This fixture is not by
  itself fresh-agent evidence; the independent positive and negative probe remains an `EP-06`
  review-gate condition.
- Negative and adversarial evidence: static managed-contract assertions cover unrelated staged
  files, ambiguous scope, rejected normal push, PR, merge, release, force-push, branch deletion,
  destructive Git, auth failure, and policy rejection. The adversarial descriptor includes an
  executable plan-adjacent file, Git hook, malformed metadata, traversal name, oversized metadata,
  and secret-like placeholder. Scanner Git-policy tests disable hooks, sanitize inherited Git
  configuration, reject unsafe transports, and read branch content only as inert Git objects. A
  real local-remote adversarial branch test confirms neither its executable nor hook runs and its
  ignored sentinel/placeholder bytes do not enter scan output.
- Public artifact hardening: remote-overlay plan validation and inventory now redact absolute and
  `file:` machine locators to stable `local-git` evidence labels and replace internal scan-error
  messages with fixed public descriptions. A regression covers members, candidates, scan errors,
  preservation of ordinary provider locators, and an end-to-end failing local-source inventory
  artifact. Exact locators and internal error detail remain only in ignored scanner-local cache for
  operational recovery.
- Benchmarks: a representative local inventory measured one repository, one local/default view,
  33 plans, zero provider API calls, concurrency one, and 1.07 seconds wall time. The deterministic
  fake-GitHub portfolio measured two repositories, two branches, four plans, two API calls,
  concurrency four, and sub-millisecond scanner duration; its complete test took 2.42 seconds.
  No V1 SLO or incremental cache is justified by this evidence.
- Live-provider boundary: authenticated `gh` read access was available for a redacted repository
  visibility preflight. No live portfolio-completeness claim was made because the producer is not
  yet enrolled on its default branch. No credential value, private source locator, or failed raw
  artifact was committed.
- Consumer and discovery evidence: release/adoption notes and a dedicated impact record separately
  classify instruction, validator, scaffold, migration, additive placement/local-path, and safety
  effects. `discovery-capability-evidence-matrix.md` maps every frozen must-cover item and explicit
  exclusion to designated passing source evidence while retaining the real-Windows and producer
  rollout qualifications.
- Additional-repository handoffs: current read-only inventories and owner-scoped migration
  recommendations were produced for three explicitly scoped private repositories. They remain
  under `/private/tmp`; no private names, counts, paths, or topology were added to this public repo,
  and no file, branch, remote, plan, register, or configuration in another repository changed.
- Register compatibility correction: the approved architecture freezes the legacy register at the
  mixed-mode baseline, so the original task to update its current entry was unsafe. The two README
  entry points were updated instead, and a byte comparison against baseline
  `3b134a07b806d1522e7cf76fe52e2028d79155c2` proves the register remains unchanged.
- Local validation: `gofmt` and diff review passed; `go test ./...` passed; the focused catalog,
  portfolio, commands, and CLI race suite passed; `go vet ./internal/... ./cmd/...` passed; all 147
  Python tests passed through `uv run --with pytest python -m pytest -q`; JSON-schema, Markdown
  header, public-core, and release-manifest validators passed; `git diff --check` passed; and the
  packaging tests proved source/package byte identity for every changed embedded resource. A fresh
  source binary also validated and listed all 33 local canonical records in mixed mode; its 25
  diagnostics are the recorded legacy-title/status and unpaired `OK-PR-014` compatibility warnings,
  not errors.
- Remaining source handoff: no release tag, public asset, consumer sync, external-repository write,
  version mutation, or release-candidate commit/push occurred. Because the validated worktree is
  not yet an immutable commit and `EP-05` remote/canonical closure is unresolved, no exact source
  revision is handed to `EP-07` yet.
- Coordination visibility: bounded plan/evidence checkpoint
  `a85bc2e0151f1d3d63d4999190c7a20810d24f22` was committed and normally pushed to
  `codex/semantic-plan-catalog-coordination`. It contains only the active plan/log and plan-scoped
  inventory, impact, and discovery-evidence attachments; implementation source remains local.
- Source gate round one: immutable local candidate
  `bdb386618b927151c64122a812bfb7b033060fc1` was rejected with two high and one medium finding.
  The reviewer found that default-branch enrollment still depended on the producer's untracked
  consumer installation, public error messages could retain absolute local paths, and the recorded
  activation evidence was an automated fixture rather than a fresh-agent probe.
- Source gate remediation: normal consumers retain installed-marker-plus-lock enrollment, while a
  Kit source repository may use a schema-valid tracked root content manifest without committing
  `.codeheart/kit/` or a lock. Remote-overlay public artifacts now use fixed error descriptions and
  have an end-to-end failing-local-source redaction test. The full Go suite, 147 Python tests,
  focused race suite, vet, schema, Markdown, public-core, release-manifest, diff, Windows-amd64
  cross-compilation, and source-CLI local catalog checks pass. Fresh-agent evidence and a new
  independent source review remain pending.
- Source gate round two: accepted exact remediation commit
  `deb39dc5b851c056471a2510c579a254e04f9a81` with no material findings. From copied exact managed
  routing and activation instructions, the fresh reviewer created an isolated repository and bare
  remote, activated branch `codex/probe-activation`, committed only the canonical implementation
  plan and sibling execution log as `0801958e3e8d7493488fbf04f4eaa34bd4aef405`, and normally pushed
  it without a second approval. Unrelated `src/app.txt` remained dirty locally with baseline bytes
  on the remote; the remote branch diff contained only the two planning files and no PR ref.
  Negative probes stopped or refused unrelated staged files, ambiguous repository or branch,
  unauthorized code, PR, merge, release, force-push, deletion, history rewrite, destructive Git,
  missing authentication, policy rejection, and rejected normal push. The reviewer independently
  passed the full Go and 147-test Python suites, race, vet, schema, Markdown, public-core,
  release-manifest, diff, frozen-register, source-marker enrollment, and failing-local-source
  redaction checks from a clean archive. No shared-workspace or unauthorized repository write,
  release version, tag, or public asset changed.
- Source-and-mixed integration CI: ready PR `#3` was opened from exact source head
  `43d11826adadf70c95a5a1493ec6edb47e52c0a1`. Both macOS validation runs passed, while both
  Windows runs exposed three platform-compatibility defects before merge: checkout line-ending
  conversion made worktree hashes differ from committed blob bytes, long temporary fixture paths
  exceeded Git for Windows defaults, and stale-transaction cleanup retained a child evidence handle
  while renaming its parent directory. The remediation retains exact worktree SHA write
  preconditions, proves the reviewed revision through clean index/blob identity, enables
  scanner-owned Git long-path handling, and closes/reopens only the Windows evidence handle around
  quarantine while revalidating identity and bytes. A new CRLF regression, all focused packages,
  the complete race-enabled Go suite, Windows amd64 test-binary cross-compilation, vet, all 147
  Python tests, JSON-schema, Markdown, public-core, release-manifest, and diff gates pass locally.
  Real Windows execution remains pending the updated PR checks; no failed check was bypassed.
- Source-and-mixed integration CI round two: updated head
  `6b9ef22a019623d1eaf2245911660a0b93823a83` passed both macOS runs. Both Windows runs confirmed
  that long-path mirror refresh and stale-transaction cleanup were fixed, then exposed two residual
  boundaries: a later clean branch checkout could convert line endings after inventory, and cache
  replacement/read verification could encounter Windows sharing or identity-change races. The
  second remediation defines a clean inventory SHA as the committed blob SHA, permits migration
  only when differing clean checkout bytes normalize to that blob through line endings alone, and
  derives the transaction precondition from the exact current bytes. The publisher closes its own
  Windows target handle before replacement and retries only bounded sharing failures while
  revalidating parent and target identity; `ReadCachedCatalog` retries only bounded sharing or
  atomic identity-change reads. The concurrency regression now exercises that supported catalog
  API. Focused packages, repeated CRLF/branch tests, 20 race-enabled cache-publication repetitions,
  the complete race-enabled Go suite, Windows amd64 test-binary cross-compilation, vet, all 147
  Python tests, JSON-schema, Markdown, public-core, release-manifest, and diff gates pass locally.
  The second Windows rerun remains pending; no failing gate was bypassed.
- Source-and-mixed integration CI round three: both Windows jobs at exact head
  `5b627b9238b0b9f6861f1002980e438e2ae9825a` passed the full `go test ./...` step, proving the
  catalog, cache, and transaction remediations on real Windows. They then reached the previously
  unexecuted Python routing group and exposed a fixture-only portability defect: the activation
  proof passed a Windows backslash path to Git's platform-neutral `revision:path` object syntax.
  The fixture now uses `Path.as_posix()` for all three remote object assertions. The exact focused
  activation test and all 50 grouped schema, routing, and packaged-resource tests pass locally. A
  fresh real-Windows rerun remains required; no failed gate was bypassed.

## EP-07 Delta - Version Bump, Reproducible Release, And Public Verification

Status: pending. The explicit release gate remains unresolved until the release candidate exists.

## Final Validation

Pending completion of all seven epics.
