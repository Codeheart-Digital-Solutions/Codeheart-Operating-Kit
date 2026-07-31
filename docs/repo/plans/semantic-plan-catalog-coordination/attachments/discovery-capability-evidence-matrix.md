Last updated: 2026-07-31T23:46:52Z (UTC)

# Discovery Capability Evidence Matrix

This matrix maps every frozen discovery `Must cover` item to its designated source evidence. The
full local Go suite and all 147 Python tests passed on 2026-07-31. Windows packages cross-compile;
real-Windows execution remains an `EP-07` release-candidate gate. The source capability evidence is
separate from this producer's still-blocked default-branch enrollment and canonical cutover.

## Canonical Metadata, Identity, And Families

| Must-cover contract | Designated passing evidence |
| --- | --- |
| Bounded markers and visible YAML below the H1 | `TestParseCanonicalMetadataAndCurrentImplementationTitleShape`, `TestParseRejectsMissingMisplacedMultipleAndUnknownMetadata`, and `plan-catalog-format.md` |
| Existing header and H1 remain content authority | `TestParseCanonicalMetadataAndCurrentImplementationTitleShape`, `TestHistoricalImplementationHandoffStatusMigratesWithoutChangingHeader` |
| Required version, ID, kind, purpose, and catalog timestamps | `TestMetadataValidationFailuresHaveStableCodes` plus positive/negative `plan-metadata` schema fixtures |
| Optional family, classifications, relations, and aliases | `TestParseCanonicalMetadataAndCurrentImplementationTitleShape` plus strict schema validation |
| Portable lowercase ASCII identity grammar | `TestIdentityNormalizationAndOwnership`; Windows cross-compilation of `internal/plancatalog` |
| Repository namespace comes from configured member identity | `TestIdentityNormalizationAndOwnership`, `TestMembershipRequiresDefaultBranchKitLockV2IdentityAndMatchingHome` |
| ID survives title, path, repository-name, status, and classification changes | `TestRenamedPlanKeepsSemanticIdentity` and path-independent semantic-ID fixtures |
| Discovery, implementation, and family remain separate records | `TestListViewKeepsDiscoveryImplementationAndFamilyAsSeparateRows` |
| Family README authority and derived membership | `TestFamilyREADMEQualificationAndDerivedMembership`, `TestRootPlansReadmeNeverQualifiesAsFamilyRecord` |
| Same-ID branch records remain separate observations | `TestDiscoverPreservesDuplicateIDsAndCatalogModes`, `TestPortfolioScanUsesRemoteBaselinesAndIndependentChangedBranchObservations` |
| Schema, identity, relation, family, header, and ownership validation | `TestMetadataValidationFailuresHaveStableCodes`, `TestMalformedMetadataNeverFallsBackToLegacyAndInvalidFamilyPlacementIsDiagnosed`, full schema suite |
| Managed references, templates, hooks, mirrors, and safe examples | `test_low_context_planning_and_portfolio_routes_are_installed`, `test_changed_source_and_packaged_resources_match`, Markdown/public-core validators |

## Local Views And Compatibility Migration

| Must-cover contract | Designated passing evidence |
| --- | --- |
| Generated local listing avoids shared-summary branch edits | `TestPlansListJSONGolden`, `TestLegacyReconciliationAndMixedViewHaveOneRowPerCanonicalPath`; frozen-register contract in `plan-register-format.md` |
| Legacy registers and older branches remain readable | `TestAbsentCatalogModeResolvesToLegacy`, mixed-mode reconciliation fixtures |
| Enumerate plans, then require semantic review | `TestInventoryCapturesRevisionCoverageUnpairedEvidenceBranchTouchesAndDirtyOverlap`; reviewed-ledger requirement in `migrate-plan-catalog.md` |
| Reconcile aliases, sessions, relations, lifecycle, and coordination evidence | Ledger schema round-trip plus the reviewed producer migration ledger and inventory |
| Source revision and content hash guard insertion | `TestMigrationChangedHashAndDirtyOverlapProduceZeroWrites`, `TestMigrationCommitPreconditionPreservesConcurrentEdit` |
| Detect branch touches and assign owners | `TestInventoryCapturesRevisionCoverageUnpairedEvidenceBranchTouchesAndDirtyOverlap`, `TestActiveBranchTouchesIncludeDeletionsAndBothRenamePaths` |
| Changed sources skip and require reevaluation | `TestMigrationChangedHashAndDirtyOverlapProduceZeroWrites`, `TestMigrationActiveBranchSkipAndTransactionalRollback` |
| Durable ledger records confidence, ambiguity, aliases, conflicts, deferrals, and owners | `TestReviewedLedgerRoundTripsThroughDurableSchema` plus positive/negative ledger schema fixtures |
| Separate compatibility, inventory, authoring, and closure phases | `migrate-plan-catalog.md` and `plan-register-format.md`; routing-contract tests |
| Enforce metadata for new plans after cutover | catalog-mode discovery/validation tests and managed discovery/implementation runbooks |
| Stop manual register writes after reliable mixed views | frozen-register byte precondition in `TestMixedModeRequiresPreCutoverPlanAndFrozenRegisterEvidence` and managed register doctrine |
| Reconcile default and active remote-overlay coverage | `TestPlansRemoteOverlayValidationAndInventoryUsePushedRefsOnly` and portfolio multi-branch scan; producer rollout remains blocked until default-branch enrollment |
| Dry-run, idempotency, preservation, rollback, and unknown-consumer compatibility | `TestMigrationDryRunApplyChronologyIdempotencyAndCanonicalCoverage`, transactional rollback/race tests, sync preservation tests |

## Portfolio Discovery, Scanning, And Coordination Home

| Must-cover contract | Designated passing evidence |
| --- | --- |
| Portfolio v2 carries stable home and member identity plus sources | positive/negative config-v2 schema tests and `TestLoadConfigNormalizesV1V2CommittedAndLocalSources` |
| Portfolio v1 paths remain compatible | `test_kit_config_schema_accepts_v1_path_only_coordination_home_fixture`, `TestLoadConfigNormalizesV1V2CommittedAndLocalSources` |
| Machine-specific source scope stays in `.codeheart/local/` | config schema, placement contract, and configure/init/sync tests |
| Membership requires authorized scope, default marker, member role, repository ID, and matching home | `TestMembershipRequiresDefaultBranchKitLockV2IdentityAndMatchingHome` |
| Unconfigured Kit repositories are candidates only | `TestPortfolioScanUsesRemoteBaselinesAndIndependentChangedBranchObservations` and membership fixtures |
| No committed central repository list | config-driven `BuildSources` path exercised by local/GitHub source tests and managed setup doctrine |
| Provider-neutral scanner lives in the Go CLI | `internal/portfolio`, grouped CLI dispatch/help tests, macOS/Windows grouped workflow jobs |
| Local Git and authenticated read-only GitHub-through-`gh` adapters | local multi-repository scan plus `TestGitHubSourcePaginatesRetriesRedactsAndHonorsCancellation` |
| Auth/access preflight stores no credentials and executes no branch code | GitHub redaction/cancellation test; `TestScannerTreatsAdversarialBranchContentAsInertData`; Git protocol, environment, hook, and command-policy tests |
| Scan complete default-branch baselines | `TestPortfolioScanUsesRemoteBaselinesAndIndependentChangedBranchObservations` |
| Enumerate every accessible unmerged remote branch | same multi-branch test plus `TestPlansRemoteOverlayValidationAndInventoryUsePushedRefsOnly` |
| Merge-base filtering ingests only added or materially changed plans | multi-branch scan and `TestBranchOverlaySuppressesPureRenameAndRetainsChangedRename` |
| Preserve conflicts, retire merged/deleted observations, label stale data, and disclose revisions/errors | multi-branch scan, incomplete-history tests, complete-cache preservation tests |
| Refresh on demand before current analysis | `refresh-portfolio-catalog.md` and `test_low_context_planning_and_portfolio_routes_are_installed` |
| Keep member facts separate from coordination-owned strategy | strategic-overlay authority tests and `TestPortfolioScanUsesRemoteBaselinesAndIndependentChangedBranchObservations` byte preservation |
| Install setup/maintenance routes, schemas, validators, mechanics, mirrors, and fresh scaffolds | configure tests, routing tests, schema suite, init/onboard/sync tests, packaged-resource parity |
| Reusable later trigger calls the same scanner authority | grouped `portfolio scan` CLI is the single refresh entry point documented by the refresh runbook |

## Activation Publication And Planning Workflow

| Must-cover contract | Designated passing evidence |
| --- | --- |
| Activation grants plan-only commit and normal-push authority | `test_activation_route_grants_only_plan_checkpoint_normal_push` |
| Confirm or create one unambiguous policy-compliant work branch | managed drafting/execution runbooks plus authority-contract assertions |
| Checkpoint only the plan and directly required planning metadata | `test_activation_checkpoint_publishes_only_planning_paths_to_local_remote` |
| Material active-plan updates use the same bounded authority | managed drafting/execution wording asserted by routing tests |
| Do not ask for a second push approval | `test_activation_route_grants_only_plan_checkpoint_normal_push` |
| Exclude unrelated code, PR, merge, release, force-push, deletion, and destructive Git | static negative assertions in `test_activation_route_grants_only_plan_checkpoint_normal_push` |
| Stop on auth, policy, rejected push, unrelated dirt, or ambiguous scope | managed execution blocker contract asserted by routing tests |
| Disclose global-visibility limits when publication fails | managed execution and refresh runbooks; routing-contract tests |
| Implementation need not wait for merge | explicit execution-runbook assertion in `test_activation_route_grants_only_plan_checkpoint_normal_push` |
| All planning workflows carry proportional hooks | `test_low_context_planning_and_portfolio_routes_are_installed`, `test_new_coordination_runbooks_have_complete_execution_contracts`, packaged parity |

## Explicit Exclusion Audit

The source introduces no UUID or central counter allocation, mandatory taxonomy, central repository
list, branch-code execution, background service, webhook, GitHub App, required schedule, global
plan-body copy, migration-on-sync, repository-wide freeze, automatic PR/merge/release, or standing
arbitrary-push authority. Static routing, schema, sync-preservation, Git-policy, and public-core
tests guard these boundaries.
