Last updated: 2026-10-05T07:42:54Z (UTC)
Created: 2026-10-05
Status: active
Execution log: claude-code-operating-guide_execution_log.md

# Claude Code Operating Guide

<!-- BEGIN CODEHEART PLAN METADATA -->
```yaml
plan:
  schema_version: 1
  id: codeheart-operating-kit.implementation.claude-code-operating-guide
  kind: implementation
  purpose: Deliver an optional Claude Code companion to the generic agent task coordination contract, release it through the guidance route and verify adoption on assigned consumer default branches.
  first_cataloged: 2026-10-05T07:19:31Z
  catalog_metadata_updated: 2026-10-05T07:42:54Z
  relations:
    - kind: related
      target: codeheart-operating-kit.implementation.proportionate-agent-workflows
```
<!-- END CODEHEART PLAN METADATA -->

## Overview

The Kit's coordination doctrine is tool-neutral, but its only tool-specific companion,
`codex-task-operations.md`, maps it to the Codex app. Agents running in Claude Code have no
supported mapping for sessions, starting implementers, messages between sessions, worktrees,
permission modes, archiving, project memory or citing sessions, so they improvise. This plan adds
one optional Claude Code companion beside the Codex note, routes to both, and ships the result as
an instruction-only patch release through the existing guidance candidate route, followed by
adoption on the assigned consumer default branches.

The accepted commissioning discovery is retained by the coordinating owner. This public plan
contains the complete reusable scope; private organization facts, consumer names, paths and task
locators stay in the commissioning assignment. Onboarding and capability wording that recognize
the agent tool are a later, separate plan because they change the CLI and need broad validation.

Independent planning review on 2026-10-05 found the plan ready after minor fixes. A simulation of
every EP-01 edit passed the guidance guard, release-identity validation and the guidance lanes'
Python tests. Three material points were corrected in this revision: explicit release-identity
and graph-digest steps; correct root-instruction guidance; and a genuinely fresh walkthrough
agent. Minor points were also corrected.

The accountable owner approved execution on 2026-10-05. This plan is active under the delivery
grant in Section 2.4, with one dedicated implementer.

## Essential Context Reference Files

| Source path | Purpose |
| --- | --- |
| `AGENTS.md` | Producer authority, public-core safety and protected existing work. |
| `docs/repo/runbooks/change-operating-kit.md` | Producer source route and focused feedback. |
| `docs/repo/reference/consumer-impact-classification.md` | Adding a managed file under an existing component target is an `instruction-only change`. |
| `docs/repo/runbooks/release-operating-kit.md` | Guidance candidate scope, packaging, publication and released-asset smoke. |
| `docs/repo/plans/release-validation-effort/release-validation-effort_execution_log.md` | Accepted broad-source anchor and guidance-lane evidence. |
| `docs/repo/plans/git-delivery-discipline/git-delivery-discipline_execution_log.md` | Most recent guidance release, identity edits and main-branch adoption pattern. |
| `components/agent-interface/component.yaml` | Managed file declarations; the Codex note is declared at the end of the list. |
| `components/agent-interface/managed/README.md` | Agent-interface route list. |
| `components/agent-interface/managed/reference/agent-task-coordination.md` | Generic contract the companion maps; its closing paragraph points to the Codex note. |
| `components/agent-interface/managed/reference/codex-task-operations.md` | Structure and tone model (84 lines). |
| `components/agent-interface/managed/reference/root-agents-md-contract.md` | Root `AGENTS.md` contract that the companion's instruction-file section must respect. |
| `components/agent-interface/managed/reference/operation-routing-and-dispatch.md` | Routing standard for the routing probe. |
| `components/agent-interface/managed/runbooks/maintain-operating-kit-installation.md` | Supported consumer upgrade, preservation and verification. |
| `scripts/validate-guidance-candidate.py` | Guidance eligibility rules this delivery must stay within. |

## Table Of Contents

- [Section 1 - Foundation](#section-1---foundation)
- [Section 2 - Strategy](#section-2---strategy)
- [Section 3 - Execution Plan](#section-3---execution-plan)
- [Section 4 - Future Planning](#section-4---future-planning)
- [Revision Notes](#revision-notes)

# Section 1 - Foundation

## 1.1 Goal Of The Implementation

An agent working in Claude Code, reading only the installed Kit, can name its session, commission
an implementer, report back between sessions, work in the right worktree, choose an appropriate
permission mode, keep shared project memory neutral, archive safely, cite its session in records
and keep the Kit's root instructions loading, all without Codex-only instructions.

Completion means all of the following are true:

1. `claude-code-task-operations.md` exists as a managed reference in the agent-interface
   component, with a byte-identical mirror, a declaration in `component.yaml`, and routes from the
   agent-interface README and the closing paragraph of `agent-task-coordination.md`.
2. The Codex companion and the generic contract keep their meaning; Codex users see no change
   beyond the second pointer.
3. Release identity is consistent: versions, checksums and the profile graph digest are refreshed,
   and `go test ./internal/hash ./internal/manifest` passes.
4. Cumulative guidance eligibility passes against the accepted broad-source anchor, and both native
   guidance lanes pass for the exact candidate the Director accepts.
5. The next unused patch (expected `v0.1.34`) is published with its normal assets, and
   released-asset smoke passes for that tag.
6. Each assigned consumer default branch runs the new version, contains the new route
   byte-identical to source, and keeps its authored configuration, plans, instructions and
   local-user files.
7. A fresh-context walkthrough and a routing probe by an agent that has not seen the plan pass,
   judged by the independent reviewer.
8. This plan, its execution log and the plans index are completed on producer main.

## 1.2 Project And Problem Context

Claude Code reads a repository's `AGENTS.md` by default only when there is no `CLAUDE.md`,
`.claude/CLAUDE.md` or `CLAUDE.local.md` in the working directory or any directory above it.
Direct reading needs Claude Code v2.1.277 or later, and some sessions cannot read `AGENTS.md`.
A `CLAUDE.md` that imports `@AGENTS.md` works everywhere. The setting that loads both files is a
user or organization setting (`/config`, user settings, `--settings` or managed settings); Claude
Code ignores it in project and local settings files. So the Kit's root routing usually works
unchanged, but one personal `CLAUDE.local.md` silently replaces it for that user.

What is missing is the tool-specific layer. Claude Code differs from Codex in ways that matter to
the generic contract:

- work happens in sessions; the desktop app gives each a session ID and an app link;
- a session can offer the user a task card that starts a new session in its own worktree, and
  background subagents handle bounded research and review inside a session;
- sessions can list peers and send them messages; a receiving session in a different permission
  mode can hold messages for its user's approval;
- desktop worktrees live under `<repository>/.claude/worktrees/<name>`, and archiving offers to
  remove a session's worktree together with its branch;
- automatic project memory is per repository and shared across its worktrees, so every session
  in the repository reads it;
- permission modes decide what runs without asking;
- there is no goal object; loops, scheduled tasks and routines exist but are not report-back
  mechanisms.

Official sources, checked 2026-10-05 and to be rechecked when authoring:

- <https://code.claude.com/docs/en/memory.md> (instruction files, `AGENTS.md`, auto memory)
- <https://code.claude.com/docs/en/desktop.md> and <https://code.claude.com/docs/en/worktrees.md>
- <https://code.claude.com/docs/en/permission-modes.md>
- <https://code.claude.com/docs/en/env-vars.md>
- <https://code.claude.com/docs/en/scheduled-tasks.md>
- <https://code.claude.com/docs/en/skills.md> and <https://code.claude.com/docs/en/plugins.md>

Task cards, messages between sessions, session links and pull-request binding were observed in
the desktop app's exposed tools rather than in public documentation. The guide must label them as
observed app contracts with their check date and tell agents to verify the tools they actually
have.

## 1.3 Current State Analysis

- `agent-task-coordination.md` is tool-neutral; only its last sentence points to the Codex note.
- `codex-task-operations.md` (84 lines) has four parts: project and task operations, direct
  report-back, one explicit whole-plan goal, and archival and preservation. It is declared as the
  last managed entry in `components/agent-interface/component.yaml` and routed from the
  agent-interface README.
- All managed sources, component declarations, profiles and `manifest.yaml` are mirrored
  byte-identically under `src/codeheart_operating_kit/resources/` for the legacy Python package.
  The guidance guard requires an identical mirror for every changed managed or declaration file.
  `tests/test_packaging_resources.py` lists only some paths; editing it would force broad scope,
  so this plan does not touch tests.
- `resources.go` embeds whole directories, so a new managed file needs no Go change.
- Release identity lives in: the agent-interface `component.yaml` version; `profiles/standard.yaml`
  version; `manifest.yaml` (release version, agent-interface version and `component.yaml` sha256,
  profile version, sha256 and `graph_sha256`); `pyproject.toml`; `internal/version/version.go`;
  `src/codeheart_operating_kit/__init__.py`; `install.sh`; `install.ps1`; and seven version tokens
  in `bootstrap.md`. The guard admits these literal and masked fields only.
- The profile graph digest is checked only by `go test ./internal/hash ./internal/manifest`, which
  the guidance lanes run. A stale digest failed a `v0.1.33` candidate run.
- The Kit contains no mention of Claude Code.
- Guidance releases reuse the accepted broad-source anchor recorded in the release-validation-effort
  log; `v0.1.33` used that anchor with `v0.1.32` as the old-CLI upgrade input.
- Known lifecycle follow-up: a newer CLI can classify a pristine older installation as partial when
  the release adds managed paths; the matching published CLI initiates the normal upgrade.

# Section 2 - Strategy

## 2.1 Implementation Strategy With Visual File/Folder Hierarchy

Add one managed reference under the existing agent-interface component, update two routing
surfaces, mirror them, and make the release-identity edits the guidance guard admits.

```text
components/agent-interface/
  component.yaml                                   # modify: version; declare the new managed file
  managed/README.md                                # modify: route to the Claude Code companion
  managed/reference/
    agent-task-coordination.md                     # modify: closing pointer names both companions
    claude-code-task-operations.md                 # create
profiles/standard.yaml                             # modify: version
manifest.yaml                                      # modify: versions, sha256 values, graph_sha256
src/codeheart_operating_kit/resources/
  (mirrors of every changed file above)            # modify/create byte-identically
pyproject.toml                                     # modify: version
internal/version/version.go                        # modify: version literal
src/codeheart_operating_kit/__init__.py            # modify: version literal
install.sh, install.ps1                            # modify: default version literals
bootstrap.md                                       # modify: version tokens only
release-notes.md                                   # modify: v0.1.34 notes
docs/repo/plans/
  README.md                                        # modify: index entry
  claude-code-operating-guide/
    claude-code-operating-guide_implementation_doc.md   # this plan
    claude-code-operating-guide_execution_log.md        # create at activation
```

## 2.2 Open Questions And Assumptions Requiring Clarification

- OQ-1 — Accepted anchor. BLOCKER: no. Affects: EP-01 guard and candidate dispatch. Use the latest
  accepted broad-source anchor recorded in the Kit execution logs; at planning time that is the
  release-validation-effort anchor used for `v0.1.33`. If a newer broad acceptance exists, use it
  and record why.
- OQ-2 — Upgrade classification with a new managed path. BLOCKER: no. Affects: EP-02 adoption.
  If the newer CLI classifies a consumer as partial, use the documented matching-CLI route; never
  edit locks or bypass checks. Record any occurrence as lifecycle evidence.
- OQ-3 — Consumer merge gates. BLOCKER: no. Affects: EP-02 adoption. Inspect each consumer's
  triggers and required checks before opening its adoption PR. If a consumer gate rejects the new
  Kit-managed path, stop that consumer's adoption and report to the Director; changing consumer
  policy or CI is out of scope.
- OQ-4 — Observed app contracts. BLOCKER: no. Affects: EP-01 content. Verify each observed behavior
  against the implementer's own exposed tools while authoring and describe only what can be
  confirmed, with its date. Live behavior between real sessions (messaging, archiving) is covered
  by this check and the commissioning organization's ongoing pilot, not by the walkthrough.
- A-1 — A consumer running `v0.1.32` can upgrade directly to the new patch through the supported
  route; verify rather than assume.

## 2.3 Architectural Decisions With Reasoning

**AD-1 — One optional companion per agent tool.** Keep the generic contract tool-neutral and add
`claude-code-task-operations.md` beside the Codex note. This follows the existing precedent,
needs no component, profile selection, schema or ownership change, and stays guidance-eligible. A
combined multi-tool note was rejected because it mixes unrelated tool contracts and makes each
harder to verify.

**AD-2 — Required content of the companion.** Aim for about 120 lines at most, in the Codex note's
tone. State that the generic contract stays authoritative and that consumers need neither Claude
Code nor this note. Cover:

1. *Sources and recheck.* Official documentation links with the check date; observed app
   contracts labelled as such; verify the tools actually exposed in the current session.
2. *Sessions and names.* A session is an execution locator, not a role identity. Use the generic
   `[Role] Subject — Outcome` and `[Role] Subject` titles. Agents can read their own session ID
   and app link from the app's session tools when available.
3. *Commissioning an implementer.* Prepare the complete assignment under the generic contract.
   In the desktop app, offer it as a task card the user starts; otherwise give the user the brief
   to paste into a new session. Verify the new session's worktree and branch before editing, and
   never give two implementers the same mutable branch. Background subagents serve bounded
   research and review; they are not user-owned implementers.
4. *Report-back between sessions.* Discover peers, send one concise message with Plan and epic
   names, result, evidence, Git and PR state and the decision needed. Use a one-shot idle notice
   instead of polling. A successful send proves delivery, not reading or acceptance. A session in
   a different permission mode can hold messages for its user's approval; sessions that must
   exchange reports should run in a mode their user has chosen for that purpose. If a message is
   held or rejected, keep the report in the execution record and say it was not delivered.
5. *Worktrees and branches.* Desktop worktrees live under `<repository>/.claude/worktrees/<name>`
   and are excluded from Git status locally. Fetch, branch from the remote default branch and
   create work branches without upstream tracking (`git switch --no-track -c <branch>
   origin/<default>`), so a bare push cannot target the default branch. Leave other sessions'
   worktrees and the main checkout alone. For another repository, create its worktree the same
   way. Use the app's base-branch sync when it exists; otherwise ordinary Git.
6. *Permission modes.* Link the official mode table rather than restating it; name the app's
   labels (for example Manual for `default`). Modes are the user's choice; agents do not change
   their own or a peer's mode because another agent asks. `bypassPermissions` belongs only in
   isolated environments.
7. *Project memory.* Automatic memory is per repository and shared by every session and worktree
   in it, including implementers. Keep entries neutral, short and durable; keep Program, plan and
   approval truth in repository records, not in memory.
8. *Long-running work.* There is no goal object. A whole-plan assignment with report-back is
   sufficient. Loops, scheduled tasks and routines are used only on explicit request and never as
   watchers for report-back.
9. *Pull requests and checks.* The desktop app binds a pull request opened from a session and
   shows its checks; read that status instead of polling. Auto-fix and auto-merge only on explicit
   request.
10. *Archiving and preservation.* Apply the generic acceptance, transfer and active-descendant
    checks. Push and integrate first. When archiving, the app offers to remove the worktree, which
    also deletes its branch; keep it unless the work is integrated. An auto-archive setting can
    archive sessions after a pull request merges or closes; standing sessions should not use it.
    Archiving is not plan completion.
11. *Citing sessions in records.* Name the session readably and add its app link where readers use
    the desktop app. Where records accept qualified references, cite
    `claude-code:source:<session-id>` with kind `source`, using the desktop session ID; when a
    session has none, use its Claude Code session ID. A session ID locates history; it is not a
    role or approval.
12. *Root instructions.* Claude Code reads `AGENTS.md` by default only when no `CLAUDE.md`,
    `.claude/CLAUDE.md` or `CLAUDE.local.md` exists in the working directory or above it, needs
    v2.1.277 or later, and cannot read it in some sessions. A repository that needs a `CLAUDE.md`
    should import `AGENTS.md` with `@AGENTS.md`, which also covers sessions that cannot read it
    directly. A personal `CLAUDE.local.md` stops `AGENTS.md` loading unless that user or the
    organization sets **Project instructions** to load both; repositories cannot set this.
    Confirm at session start that `AGENTS.md` is reported as loaded. The Kit manages only
    `AGENTS.md`.

**AD-3 — Routing.** Change the closing sentence of `agent-task-coordination.md` to point to
`codex-task-operations.md` or `claude-code-task-operations.md` for tool-specific surfaces, and add
"Optional Claude Code session operations" to the agent-interface README routes. No other managed
document changes. The route owner is the agent-interface component; the routing probe's evidence
goes in the execution log.

**AD-4 — Release.** Ship the next unused patch as an `instruction-only change` through guidance
candidate scope, with the retained accepted anchor and the latest published tag as the old-CLI
upgrade input, under the existing unsigned internal/prototype HTTPS-plus-SHA-256 boundary. Bump
the agent-interface component version and the profile version as earlier guidance releases did.

**AD-5 — Adoption.** Adopt on each assigned consumer default branch through the supported upgrade
route and one Kit-only adoption PR per consumer, merged after that consumer's required checks.
Sessions working elsewhere pick up the release when they next sync with their default branch; no
separate worktree reconciliation is planned unless the assignment names one.

**AD-6 — Validation.** The independent reviewer checks content and public safety. A separate,
fresh subagent that has not seen this plan, the diff or the log receives only a temporary
installation built from the candidate, the vague request "How do I hand this plan to an
implementer in Claude Code?" and these scenarios: naming a session, preparing a task-card
commission, reporting back, archiving safely, keeping project memory neutral, citing a session
and the `CLAUDE.md` case. The probe must reach the generic contract and then the companion
without choosing a tool prematurely. The reviewer judges that transcript and follows any
corrections. The guidance candidate and released-asset smoke prove packaging and installation.

## 2.4 Commissioning And Delivery Boundary

The accountable owner approved execution on 2026-10-05, in the commissioning Director's session.
One dedicated implementer carries EP-01 and EP-02 on a delivery branch cut from the default branch
after this activation is merged. The grant covers:

- implementation and focused checks, coherent commits, normal pushes and one delivery PR with
  updates;
- the independent review, the fresh walkthrough and the routing probe;
- the guidance candidate dispatch at the end of EP-01;
- source merge after Director acceptance of that exact validated candidate;
- publication of the next unused patch under the existing boundary, and released-asset smoke;
- Kit-only adoption PRs merged on the assigned consumer default branches after their required
  checks, in the order and with the holds stated in the private assignment;
- the final plan, log and index closure on producer main through one small PR.

Excluded: Go, Python, test, workflow or CLI behavior changes, beyond the literal version edits in
AD-4; onboarding prose; the root `README.md`; `bootstrap.md` beyond its version tokens, which
count as release identity; consumer policy or CI changes; force pushes or history rewrites;
cleanup; product work.

The Director is the acceptance owner. The implementer reports to the Director's session at the
EP-01 review, at completion, at a genuine blocker or at a material scope issue. Material scope,
authority, preservation or cost changes return to the Director. Exact consumer targets, their
order and holds, and the report-back locator remain in the private assignment.

# Section 3 - Execution Plan

## 3.0 Epic Map

| Epic | Outcome | Size | Depends on | Review point |
| --- | --- | --- | --- | --- |
| EP-01 — Claude Code companion and routing | One exact, validated candidate with the companion, routes, mirrors, release identity and notes | M | Activation | Director acceptance after independent review, fresh walkthrough, routing probe and both guidance lanes |
| EP-02 — Release and adoption | Published patch, released-asset smoke, adoption on every assigned consumer default branch, closure on main | M | EP-01 accepted | Final report to the Director |

## EP-01 — Claude Code companion and routing

### A) Epic ID, Title, And Outcome

EP-01. One exact candidate contains the Claude Code companion, its routes, mirrors, release
identity and release notes, passes local checks and both native guidance lanes, and is accepted by
the Director.

### B) Scope

AD-1 to AD-4 source work and AD-6 validation. No behavior changes in code, tests or workflows.

### C) Files Touched

The files in Section 2.1, except the execution log and closure edits made in EP-02.

### D) Acceptance Criteria And Size

- The companion covers every AD-2 topic in about 120 lines or fewer, is public-safe, and labels
  observed app contracts with a check date.
- The generic contract and the Codex note keep their meaning.
- `go test ./internal/hash ./internal/manifest` passes, then
  `python scripts/validate-guidance-candidate.py --baseline-ref <anchor> --candidate-ref HEAD`,
  `python scripts/validate-release-identity.py` and the ordinary feedback validators pass.
- The fresh walkthrough and routing probe pass, judged by the reviewer.
- Both native guidance lanes pass for the exact candidate presented for acceptance.
- Size: M, about one focused day including review and lanes.

### E) Dependencies And Critical-Path Notes

Activation and the Section 2.4 grant, with this plan merged on the default branch. Confirm the
anchor (OQ-1) before running the guard.

### F) Tasks Checklist

- [ ] Fetch and create the delivery branch from the remote default branch (which contains the
      activated plan) without upstream tracking; create the execution log; record the anchor and
      the next unused patch.
- [ ] Author `components/agent-interface/managed/reference/claude-code-task-operations.md` per
      AD-2, verifying each observed app contract against the implementer's own tools (OQ-4).
- [ ] Update the closing paragraph of `agent-task-coordination.md` and the agent-interface README
      per AD-3.
- [ ] Declare the new file in `components/agent-interface/component.yaml` after the Codex note,
      with exactly `source`, target
      `.codeheart/kit/docs/agent-interface/reference/claude-code-task-operations.md` and
      `ownership: managed`.
- [ ] Release identity:
      - bump the agent-interface `component.yaml` version and `profiles/standard.yaml` version;
      - in `manifest.yaml`, set the release version, the agent-interface version and the sha256 of
        its `component.yaml`, and the profile version, sha256 and `graph_sha256`, taking the graph
        value from the `want` line printed by
        `go test ./internal/manifest -run TestLoadEmbeddedContentManifest`;
      - make the literal version edits in `pyproject.toml`, `internal/version/version.go`,
        `src/codeheart_operating_kit/__init__.py`, `install.sh`, `install.ps1` and the seven tokens
        in `bootstrap.md`.
- [ ] Mirror every changed managed, declaration, profile and manifest file byte-identically under
      `src/codeheart_operating_kit/resources/`.
- [ ] Add `v0.1.34` release notes as an `instruction-only change` with adoption guidance.
- [ ] Run `go test ./internal/hash ./internal/manifest`, the guidance guard, the release-identity
      check and the ordinary feedback validators locally.
- [ ] Commit coherent progress, push, and open one delivery PR.
- [ ] Commission the independent reviewer and the separate fresh walkthrough agent per AD-6; apply
      corrections with the same reviewer following them; record the probe evidence in the log.
- [ ] Dispatch the guidance candidate on the exact candidate:
      `gh workflow run validate.yml --ref <candidate> -f mode=candidate -f candidate_scope=guidance -f baseline_ref=<anchor> -f upgrade_version=v0.1.33`
      and record both lane results.
- [ ] Report the validated candidate to the Director for acceptance.

### G) Implementation Notes

Write for both audiences: short imperative guidance an agent can execute, readable by a person.
Keep consumer-specific names, paths and session IDs out of the public guide. A correction after
the lanes run invalidates only the affected evidence; rerun the guard and the affected lane.

### H) Open Questions

OQ-1 and OQ-4.

## EP-02 — Release and adoption

### A) Epic ID, Title, And Outcome

EP-02. The patch is published, released-asset smoke passes, every assigned consumer default branch
runs it with authored content preserved, and the plan is closed on producer main.

### B) Scope

AD-4 publication and AD-5 adoption.

### C) Files Touched

Release assets through the release runbook; Kit-managed paths in each consumer; this plan, its
log and the plans index at closure.

### D) Acceptance Criteria And Size

- The merged tree equals the accepted, validated candidate.
- The release has its normal assets and checksums; released-asset smoke passes for its tag.
- Each assigned consumer default branch reports the new version, contains the new route
  byte-identical to source, passes managed checksums, and preserves authored files.
- Plan, log and index show completion on producer main.
- Size: M, about one day, mostly waiting on checks.

### E) Dependencies And Critical-Path Notes

EP-01 acceptance. Inspect consumer checks before promising cheap adoption (OQ-3). Handle upgrade
classification through the matching-CLI route if needed (OQ-2).

### F) Tasks Checklist

- [ ] Merge the delivery PR after Director acceptance; verify the merged tree equals the accepted
      candidate.
- [ ] Build, verify and publish the release per `release-operating-kit.md`, then run
      `gh workflow run validate.yml -f mode=released-smoke -f release_version=v0.1.34`.
- [ ] For each assigned consumer, in the order the assignment gives: create a fresh worktree from
      its remote default branch, preview and apply the supported upgrade, verify version, the new
      route, managed checksums and preservation, open a Kit-only adoption PR, merge it after the
      consumer's required checks, and verify the remote default branch.
- [ ] Complete this plan, its log and the plans index on producer main, and send the final report
      to the Director.

### G) Implementation Notes

Keep private consumer roots, names and locators in the assignment and the private report, not in
this public log.

### H) Open Questions

OQ-2, OQ-3 and A-1.

# Section 4 - Future Planning

## 4.1 Deferred Tasks

- Tool-aware onboarding and capability wording: recognize Claude Code through the documented
  `CLAUDECODE=1`, describe the baseline capabilities per tool, correct stale installation wording,
  and remove or isolate the legacy Python Codex installer and its test hazard. This needs CLI and
  test changes and broad validation, so it is a separate later plan.
- Codex-only examples elsewhere (`/goal` triggers in discovery and execution runbooks, the routing
  reference's instruction-priority example, agent-memory session paths): generalize with the
  onboarding plan.

## 4.2 Future Considerations

If more agent tools appear, keep one optional companion per tool rather than growing the generic
contract.

# Revision Notes

- 2026-10-05: Drafted from the commissioning organization's accepted discovery and its
  two-step decision (guide first).
- 2026-10-05: Activated after the accountable owner approved execution; Section 2.4 records
  the delivery grant.
- 2026-10-05: Independent planning review corrections:
  - explicit release-identity and graph-digest steps, with the Go hash and manifest tests;
  - corrected root-instruction guidance (`CLAUDE.local.md`, `@AGENTS.md` import, user-level
    setting);
  - a separate fresh walkthrough agent with the probe evidence in the log;
  - candidate validation before acceptance;
  - project-memory guidance;
  - the session-ID fallback;
  - a length target;
  - clearer exclusions, sizes and open-question effects.
