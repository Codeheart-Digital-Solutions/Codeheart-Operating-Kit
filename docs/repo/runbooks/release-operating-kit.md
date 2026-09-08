Last updated: 2026-09-08T14:57:57Z (UTC)

# Release Operating Kit

Use this runbook before publishing a Codeheart Operating Kit release.

Audience: maintainer-facing

Intent:
Publish the validated candidate under the approved audience and verify public assets.

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
3. Read `docs/repo/reference/consumer-impact-classification.md`.
4. Confirm the intended version and consumer-impact classification.
5. Confirm release notes cover consumer-facing behavior and migration.
6. Select broad (default) or eligible guidance candidate scope below on the intended source ref.
   Require applicable public-core, Markdown, schema, content identity, native installer and
   release-contract evidence. Broad source/compatibility results may be retained for eligible
   guidance only under the explicit applicability boundary below. Ordinary feedback alone never
   qualifies a release.
7. Build macOS universal and Windows x64 packs twice with `scripts/build-release-assets.py`; require
   byte-identical output from both builds.
8. Confirm each pack contains the expected binary plus `bootstrap.md`, `install.sh`, `install.ps1`,
   `release-notes.md`, `INSTALL.md`, `content-manifest.yaml`, `pack-manifest.json`, and
   `checksums.txt`, with no Python wheel or `*.dist-info` payload.
9. Verify the complete external catalog -> archive digest -> pack manifest -> payload checksums ->
   content identity -> binary digest and version chain.
10. Confirm the external release catalog is generated after the packs. Keep archive URLs and
    digests out of embedded `manifest.yaml`; add final public URLs only to the external catalog.
11. Run isolated macOS and real Windows fresh-install and upgrade dry-run/apply/failure paths.
    Require failed verification, replacement, reconciliation, or post-check to preserve or restore
    the prior installation.
12. Generate and verify sidecar SHA-256 checksums for every public asset.
13. Treat locally generated or unsigned assets as release-candidate evidence only. Before broad
    public distribution, require signing/notarization or explicitly approve and record the
    unsigned internal/prototype boundary.
14. Verify tag creation, publication and the signing/audience boundary against the existing
    grant. Stop only if authority or a required action-time condition is missing; do not request
    duplicate approval for already-covered effects.
15. If publication is authorized, verify the target commit is the validated commit, create the
    tag, publish the packs, catalog, installers, notes, and checksums, then record URLs, digests,
    platform evidence, signing state, and residual risk.

## Stop Conditions

Stop before publishing when public-core hygiene fails; schemas or migrations are incompatible;
release assets are not reproducible; any catalog-to-binary digest is missing or inconsistent;
packs contain a Python payload; macOS or Windows install/upgrade validation fails; rollback does
not preserve the prior installation; release notes omit consumer-impacting changes; staged assets
are confused with live public assets; signing/notarization readiness is unresolved for the
intended audience; publication is not explicitly authorized; or the release target commit differs
from the validated commit.

## Candidate And Published-Asset Modes

Select the next unused patch from current source, tags and releases before freezing the candidate.
Update release identity, notes and declared resource mirrors through source. Recheck current
required-check configuration before publishing workflow changes.

Default manual dispatch is `candidate` with `candidate_scope=broad` and validates checked-out source and staged packs. It never
downloads that candidate as a public release. Invoke on the exact intended branch/ref:

```sh
gh workflow run validate.yml --ref <candidate-ref> -f mode=candidate
```

`candidate_lane` defaults to `all`. To rerun one invalidated lane after a scoped correction, pass
`-f candidate_lane=windows` (also `macos`, `ubuntu`, `git-2-43`).
For a Windows packaging/native-smoke-only correction with applicable successful source evidence,
use `candidate_scope=broad` and `candidate_lane=windows-smoke`; it omits Go/parity/history suites
and rebuilds packages before native smoke. Record the retained source run and unchanged inputs.
This option is partial evidence, not a new acceptance scope or an authority bypass. A selected-lane pass is partial
evidence, never full release acceptance by itself. Record the previous run, unchanged relevant
inputs and retained results alongside the correction run. If applicability is uncertain, use all.

Record the run URL and resolved commit. Broad candidate jobs cover macOS and Windows native suites,
staged installers and old-version upgrade preservation, Ubuntu semantic/performance coverage and
exact Git 2.43 proof regression. Each broad native lane runs the Go suite once as a visible CI step with its retained 30-minute
per-package Go timeout; Ubuntu runs it directly. Do not mistake this for a whole-command
wall-time limit: compilation and multiple package executions can take longer. The builder only packages and
verifies artifacts; invoking it alone does not prove source acceptance.
Release test cases that repeat entire builds are covered by that builder in CI; it runs under an
unused package-index URL to prove no Python package retrieval is needed. Separately
configured benchmark and oldest-Git cases remain distinct. The macOS builder performs two builds
and byte comparison for both supported packs and uploads `candidate-release-assets` for reuse.
No new platform support is implied by the Ubuntu semantic lane.

## Guidance Candidate Scope

Use guidance scope for reviewed managed instruction changes under existing component ownership
and target contracts, including declared new managed files. A new component, scaffold, ownership
mode or placement contract is not such an addition. Review policy and operational implications:
a Markdown suffix does not waive safety-policy review or affected scenario/routing probes.
Unknown dependencies and consequential runtime changes require broad scope.

Select an exact ancestor commit whose coherent source and broad candidate evidence the owner has
accepted. Retain that broad-source anchor across successive guidance releases until a new broad
acceptance replaces it; do not hide cumulative executable changes behind the last guidance tag.
The existing execution log records accepted source, run/commit evidence, semantic impact and why
retained lanes remain applicable. The guard cannot manufacture acceptance from prose. There is
no receipt registry or automatic authority check.

The source anchor may be unpublished. Separately select an explicit verified published tag for
the matching old CLI upgrade smoke; this release must precede the candidate version. This lets
workflow improvements be tested before publication. A controlled demonstration candidate is
validation evidence only and need not be published or merged as a feature release.

```sh
python scripts/validate-guidance-candidate.py --baseline-ref <accepted-broad-source-sha> --candidate-ref <candidate-sha>
python scripts/validate-release-identity.py
gh workflow run validate.yml --ref <candidate-ref> -f mode=candidate -f candidate_scope=guidance -f baseline_ref=<accepted-broad-source-sha> -f upgrade_version=v<published-old-version>
```

The workflow enforces cumulative committed-tree eligibility before native jobs. It admits
known managed content and identical declared mirrors, narrowly recognized declaration additions
and literal release identity edits, planning/log/index Markdown under `docs/repo/plans/`, and
reviewed release-note prose. It rejects unknown paths, removals, renames, mode changes, unrelated
code, schema/ownership/selection changes and changed validation machinery. See the script's
explicit rules/tests for the known literal fields; an unrecognized edit uses broad scope.
No suffix-wide exemption for executable mechanics or governing instructions exists.

Both native guidance lanes run focused hash/embedded-manifest and routing/schema/resource checks,
rebuild reproducible packages and perform existing staged installation, historical upgrade and
containment smoke. They also compare all declared managed source bytes with fresh and upgraded
installed targets. A verified matching published CLI initiates the current-release upgrade;
dry-run and failed catalog verification preserve consumer and binary state, successful apply/check
preserves authored config, instructions, plans and local-user files. Native command exit codes
are asserted, including Windows intermediate calls. Wait for both expected binary version and a
completed installed-state check before preservation assertions; version replacement can precede
transaction completion. Retry only the known transaction-in-progress state, not drift or corruption. Missing-catalog smoke does not claim to
retest every transaction rollback phase; retain applicable unchanged runtime fault-injection
results and rerun them when affected.

Guidance omits the broad Go, full Python parity/plan-history compatibility, oldest-Git and scale
benchmark runs. These results must remain applicable to the unchanged mechanisms. Relevant
source, toolchain, dependency, target or behavior changes invalidate affected evidence; incidental
host/image/log differences do not automatically demand a matrix rerun. Investigate uncertainty
and use broad scope when applicability cannot be established.

`candidate_lane=macos` or `windows` can rerun an invalidated guidance lane with the same explicit
scope/anchor/tag inputs. Ubuntu and oldest-Git are not guidance lanes. A selected-lane pass is
partial evidence; list fresh and retained results for the exact candidate. Record scoped step
timing and cache conditions honestly, without treating earlier warmed timings as a guarantee.
Broad remains the default, including for changes to this route, workflow, guard or build contract.

## Artifact Acceptance and Publication

Retrieve the successful run's artifacts and verify their catalog/checksum/version identity before
publication. Required gates apply to the intended source and release payload; retain valid
evidence across review corrections only if their relevant inputs are unchanged. Rerun invalidated
checks, not every suite after a log update. If integration changes release inputs, rebuild and
validate those changes before tagging. Record the final target and evidence applicability.

After publication, use only the explicit released tag:

```sh
gh workflow run validate.yml --ref <reviewed-workflow-ref> -f mode=released-smoke -f release_version=v<released-version>
```

`released-smoke` rejects a missing or malformed tag and downloads its public assets; a nonexistent
release fails retrieval. It runs only the two public native smoke jobs after input validation,
without restarting source suites. Never infer public-release intent from the checkout version.
Candidate mode rejects a release tag input to avoid misleading evidence. Concurrency cancels only
superseded ordinary feedback, leaving candidate/install runs undisturbed.

Publication is followed by the assigned lifecycle adoption, installed route verification and
consumer-authored content preservation. Public smoke alone is not consumer adoption. Use
`components/agent-interface/managed/runbooks/maintain-operating-kit-installation.md`; exact private
targets and operational records stay outside public producer evidence.
