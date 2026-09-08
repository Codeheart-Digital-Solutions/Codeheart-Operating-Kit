Last updated: 2026-09-08T14:09:59Z (UTC)

# Change Operating Kit

Use this runbook before changing Codeheart Operating Kit source, docs, schemas, templates,
validators, installers, release assets, or CLI behavior.

Audience: maintainer-facing

Intent:
Change producer source with proportionate checks and preserved consumer boundaries.

Success:
The intended change is verified with truthful source, validation and delivery evidence.

Agent judgment boundary:
Choose ordinary in-scope implementation details and reuse sufficient authority. Preserve required
integrity, platform, ownership and signing/audience gates.

Stop boundary:
Stop on a failed required gate or a material scope, authority or preservation conflict.

## Procedure

1. Read `AGENTS.md`.
2. Read `README.md`.
3. Read `docs/repo/reference/placement-contract.md`.
4. Read `docs/repo/reference/consumer-impact-classification.md`.
5. Classify the change by consumer impact.
6. Check whether the change belongs in the Operating Kit or in a consumer repository.
7. Keep public-core material only.
8. Update the nearest README when discoverability changes.
9. Update tests, schemas, manifests, or fixtures when behavior changes.
10. Record release-note or migration-note needs for consumer-affecting changes.
11. Run the smallest validation set that proves the changed surface.
12. Summarize validation and residual risk in the PR.

For state, lifecycle, installer, or release changes, also run the matching gates:

- schema and migration tests for declaration, config, lock, content, catalog, or pack contracts;
- transaction failure tests for planning, staging, commit, post-check, rollback, and recovery;
- build-twice byte comparison plus catalog-to-binary verification for release-pack changes;
- isolated installer and upgrade success/failure paths for each affected platform;
- consumer materialization and routing checks when managed guidance changes.
- low-context routing probes when plan authoring, activation publication, portfolio setup/refresh,
  or migration routes change;
- absent-file and preservation tests when scaffold declarations change, including proof that
  removed fresh surfaces remain preserved in existing consumers;
- source-to-packaged byte equality for every changed managed or scaffold resource.

Do not use the source repository's ignored `.codeheart/kit/` tree as producer authority. Source
components, profiles, templates, schemas, Go packages, and maintainer runbooks are canonical.

## Stop Conditions

Stop before editing when the change would expose private content, change release authority, alter
consumer ownership boundaries, or require a new public repository setting that is not already
approved.

## Feedback And Candidate Timing

Ordinary PRs and main pushes run the `feedback` job in `validate.yml`: routing/resource tests,
public-core, Markdown, schema and release-manifest checks. Feature-branch pushes do not duplicate
PR runs. No commit-message bypass is needed. This cheap result does not qualify executable changes
or an unknown changed surface for release. Inspect dependencies and run additional affected tests
for changed behavior; select the broad or eligible guidance acceptance route at the coherent candidate boundary.

Local instruction feedback (use the supported repo-local Python environment when needed):

```sh
python -m pytest -q tests/test_routing.py tests/test_packaging_resources.py
python scripts/validate-public-core.py
python scripts/validate-markdown-headers.py
```

The tests use pytest and PyYAML in the development environment. Missing local tooling follows
`components/agent-interface/managed/runbooks/handle-tooling-readiness.md`; this is producer
source development, where an editable installation is permitted.

Use one meaningful independent source review for a coherent batch with the same reviewer for
focused correction follow-up. Run broader checks after the batch is ready. Reuse evidence only
while relevant source, artifacts, configuration, target and environment remain applicable; rerun
invalidated checks. A metadata/log edit or a resumed task does not recreate native acceptance.
Classify code, test, environment, authentication and billing failures separately before retrying.
Do not increase budgets or blindly restart broad suites.

Candidate scope and explicit public-smoke invocation are owned by `release-operating-kit.md`.
Guidance scope retains fresh package/native lifecycle proof and applicable broad-source evidence;
it is never an ordinary-feedback waiver. The cumulative diff guard proves mechanical eligibility
only. Review semantic policy, routing and dependency impact separately. Runtime, ownership,
installer behavior, schema, toolchain, workflow/guard or unknown changes default to broad scope.
Known literal release-version updates are identity changes, not installer behavior changes.
Use `python scripts/validate-release-identity.py` before native dispatch to catch identity drift
without patch-specific golden-fixture edits. Preserve independent hash/corruption negatives.
Inspect current branch protection and rulesets before changing check names/triggers. Required
checks must remain reliable; resolve an enforced conflict with the owner before publication.
