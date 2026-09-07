Last updated: 2026-09-07T22:48:11Z (UTC)
Created: 2026-09-07

# Development Feedback and Proportionate Validation — Execution Log

## Authority and delivery boundary

The user authorized the complete plan on 2026-09-07, including source integration, release and
verified adoption. The [active plan](development-feedback-and-validation_implementation_doc.md)
owns scope, acceptance and exclusions. The Program Director is the acceptance owner; the dedicated
implementer receives its exact report-back locator and private adoption targets in the assignment.

Implementation branch: `codex/development-feedback-delivery`. Planning baseline: merged
[PR #16](https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/pull/16), commit
`9e73c654aed7de161f0414e37701f8426b1cd44b`.

## Activation checkpoint

- The reviewed draft and repository index were merged before activation.
- Independent planning review closed the candidate/public-release dispatch sequencing issue.
  The later user-approved publication/lifecycle clarification is in the combined source review.
- The activation checkpoint includes only the plan, this log and the index lifecycle label.
  Product source, workflows and installed consumers are unchanged at activation.
- Documentation-only checkpoints use `[skip ci]`; their relevant local document/catalog checks
  supply evidence. This does not waive any implementation-candidate or release acceptance.
- The active source branch is to be pushed before dispatch; the commissioning task records the
  resulting checkpoint and actual task creation. No goal or monitoring automation is requested.

## Progress

| Epic | State | Completion evidence |
| --- | --- | --- |
| EP-01 — Coherent operating guidance | Implemented; combined acceptance pending | Guidance and declared mirrors updated; independent review correction addressed |
| EP-02 — Feedback, enforcement and coherent acceptance | Implemented; combined acceptance pending | 27 focused tests pass; workflow mode guards exercised locally |
| EP-03 — Release, adoption and owner handoff | Release inputs prepared; delivery pending | Next unused patch selected: 0.1.32; native candidate and adoption pending |

EP-01/02 share one coherent independent source review with focused correction follow-up. After
Director acceptance, EP-03 proceeds through the already-authorized final effects and reports
actual release and per-target adoption results. Record measured outcomes and relevant deviations;
do not claim savings or completed consumer product changes from a process release.

## Combined guidance and validation checkpoint

Impact: instruction-only managed guidance and validator-only producer enforcement. No placement,
consumer ownership, authentication, support or release-audience changes. Release notes describe the
normal upgrade and route verification requirement. Source authority is tracked producer content;
all affected declared Python resources match it byte for byte.

The implementation preserves active intention elicitation and substantive technical planning,
removes checklist-word prohibitions, adds coherent batch/evidence/recovery guidance, uses one
independent review with same-reviewer corrections and establishes local commit/publication/merge
rules independent from lifecycle. Recipe guidance tests necessity before adding machinery.

### CI mapping and current enforcement

GitHub live preflight found main unprotected (branch protection endpoint: `Branch not protected`)
and no repository rulesets. No required check names are being silently skipped or renamed.

| Lane | Trigger and evidence |
| --- | --- |
| feedback | PR and main push; routing/resources, public-core, Markdown, schemas, release manifest; feature branch push does not duplicate PR |
| git-2-43-proof-validation | Explicit candidate only; exact pinned Git build and branch-proof regression |
| macos-validation | Explicit candidate; compatibility/parity, twice-built native packs, install and historical upgrade preservation; builder owns broad Go invocation |
| ubuntu-semantic-validation | Explicit candidate; broad Go, separately configured scale benchmark, schema/routing/resource/sync and compatibility |
| windows-validation | Explicit candidate; compatibility/parity, Windows-only twice-build, native install/reparse and deferred-upgrade preservation; builder owns broad Go invocation |
| macos/windows-public-release | Explicit released-smoke with required versioned tag; public download/install/check only |

Mode guard is exercised with valid candidate, rejected candidate tag, rejected empty/malformed
smoke tag, valid smoke tag and unknown mode. Candidate jobs do not download the unpublished
version. Public smoke does not set up Python or run source suites. Only superseded ordinary
feedback is cancelled; each manual candidate/install dispatch has its own concurrency group.

Inspection found additional overlapping full-suite calls inside the existing pack builder and
release tests. Native workflow jobs now rely on the builder's one Go invocation and twice-built
pack verification; CI omits equivalent full-builder Python cases. The distinct Windows-only
builder, oldest Git regression, benchmark and native preservation checks remain. Builder execution
under an unused package-index URL retains evidence that packaging needs no Python package access.

A supported local tooling environment exposed the Markdown validator traversing generated vendor
files. The validator now excludes only `.codeheart/local/` during directory traversal. A regression
proves similarly named consumer-owned `docs/local/` content is still checked. This directly related
EP-02 enforcement correction adds `scripts/validate-markdown-headers.py` and its existing tests to
the starting scope. No new evidence framework or classifier is introduced.

### Validation and review

- Focused routing, resource, workflow-boundary and Markdown regression suite: 27 passed.
- Public-core, Markdown, JSON schemas and release-manifest validators pass.
- Release identity updated to next unused patch 0.1.32; `internal/manifest` passes, including
  embedded graph identity. Full native candidate evidence is still pending.
- One independent read-only reviewer walked all requested scenarios. Routine repair, deferred
  broad checks, focused correction, published draft and later bounded activation, interruption,
  nested module/tooling/service routing and proposed-method intention elicitation pass at source
  walkthrough level. This is not live service, hosted CI or adoption evidence.
- One material review finding identified stale mandatory script-scaffolding/common-output bullets
  in planning/execution guidance. These were aligned with the owner standard and mirrored; focused
  same-reviewer follow-up closed it. The corrected selector retains the cheap version-mismatch
  test, which passes (1 test); no material source findings remain. An apparent discovery batching conflict was retracted after
  verifying it applied only to an explicitly requested per-decision review mode.

Tooling setup used the supported ignored local development environment. An inherited package-index
access failure was isolated to package retrieval and resolved using the public package source for
this public producer development environment, without changing authentication or machine policy.
No account budgets were changed and no broad suite was blindly retried.

The focused reviewer also identified a broad test selector excluding the cheap negative version
mismatch guard; it is now explicitly retained. Native builders cover required pack integrity,
not every historical Python assertion: separate-invocation Windows catalog comparison and captured
private-index-marker omission are not replayed. Those tests remain available for builder changes.

Source acceptance, hosted candidate results, integration, publication, public smoke and all assigned
consumer adoptions remain open. Exact target paths and preservation records stay private.
