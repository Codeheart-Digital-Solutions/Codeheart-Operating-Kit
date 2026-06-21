Last updated: 2026-06-21T14:43:20Z (UTC)
Created: 2026-06-21
Status: completed

# Operating Kit Portfolio Coordination And Plan Register Model Discovery

## Overview

This discovery defines how Codeheart Operating Kit should treat plan registers and optional
multi-repository portfolio coordination.

The immediate problem is conceptual drift. The current scaffolded `goal-register.md` is a very
thin agent-memory file for "ongoing work," while real Codeheart use has started to need a more
durable, reusable register for discovery plans, implementation plans, plan families, and major
workstreams across several repositories. The central artifact should become
`docs/repo/plans/plan-register.md`, while `docs/agent-memory/goal-register.md` remains available
as transitional or informal memory for goals that do not yet have formal plans. Codeheart HQ also
needs a top-level view of cross-repo portfolio work, but the mechanism should not be HQ-specific.
It should be a general Operating Kit standard that works in any repository and optionally connects
repositories to a coordination home.

This discovery is the reviewed decision source for the sibling implementation plan. The document
records the current recommended direction, visible decisions, open questions, and
implementation-shaping requirements for changes to managed agent memory doctrine, planning
runbooks, configuration schemas, installed consumer state files, and sync behavior.

## Essential Context Reference Files

| Path | Reason |
| --- | --- |
| `AGENTS.md` | Public-core safety, maintainer read order, and change safety rules. |
| `README.md` | Repository purpose, public boundary, and maintainer entry points. |
| `docs/README.md` | Public docs router and discoverability surface. |
| `docs/repo/README.md` | Repo-governance router. |
| `docs/repo/plans/README.md` | Repo-plan index receiving this discovery route. |
| `docs/repo/reference/placement-contract.md` | Current ownership model for managed kit content, kit-initialized consumer state files, scaffolds, and consumer-owned docs. |
| `docs/repo/reference/consumer-impact-classification.md` | Impact classes for future managed-doc, consumer state file, scaffold, schema, CLI, or sync behavior changes. |
| `docs/repo/runbooks/change-operating-kit.md` | Maintainer procedure for changing kit source, docs, schemas, templates, validators, installers, or CLI behavior. |
| `src/codeheart_operating_kit/resources/components/agent-memory/managed/README.md` | Current managed doctrine for agent-memory use and consumer state. |
| `src/codeheart_operating_kit/resources/components/agent-memory/scaffolds/goal-register.md` | Current thin scaffold that this discovery preserves as transitional or informal goal memory. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/discovery-workflow.md` | Current managed discovery workflow that may need a plan-register update hook. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/draft-implementation-plan.md` | Current managed implementation-planning workflow that may need a plan-register update hook. |
| `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/execute-implementation-plan.md` | Current managed implementation-execution workflow that may need lifecycle update hooks. |
| `schemas/kit-config.schema.json` | Likely future home for optional portfolio coordination configuration. |

## Related External Or Consumer Evidence

The following evidence motivated this Operating Kit discovery. These paths are not kit source of
truth and must not be copied into public kit doctrine without public-core review.

| Source | Relevance |
| --- | --- |
| Codeheart AWS Platform `docs/repo/plans/codeheart-portfolio-restructure/codeheart-portfolio-restructure_discovery_doc.md` | Accepted portfolio decision that memory format belongs in Operating Kit, memory content belongs in each repo, and a portfolio-level register belongs in Codeheart HQ once that authority exists. |
| Codeheart AWS Platform `docs/repo/plans/agent-memory-system_discovery_doc.md` | Earlier local design for `docs/agent-memory/`, including a broader goal-register concept with goals, workstreams, statuses, related docs, blockers, and next actions. This discovery narrows that responsibility by moving formal planning registration to `docs/repo/plans/plan-register.md`. |
| Codeheart HQ `docs/business/core/portfolio-map.md` | Existing business portfolio map that shows a separate, business-specific portfolio index. It is evidence for coordination needs, not a replacement for a reusable plan/workstream register. |
| Installed consumer `docs/agent-memory/goal-register.md` scaffolds | Current installed state is too thin and too broadly named to function as the canonical portfolio plan register. |

## Table Of Contents

- [Section 1 - Problem Framing](#section-1---problem-framing)
- [Section 2 - Current Model And Gaps](#section-2---current-model-and-gaps)
- [Section 3 - Target Operating Model](#section-3---target-operating-model)
- [Section 4 - Decision Inventory](#section-4---decision-inventory)
- [Section 5 - Requirements And Constraints](#section-5---requirements-and-constraints)
- [Section 6 - Open Questions And Assumptions](#section-6---open-questions-and-assumptions)
- [Section 7 - Risks And Validation](#section-7---risks-and-validation)
- [Section 8 - Review Handoff](#section-8---review-handoff)
- [Revision Notes](#revision-notes)

# Section 1 - Problem Framing

## Problem Statement

Codeheart work now spans multiple repositories: Operating Kit, AWS Platform, Foundry, HQ, and
future consumer repositories. Each repository can contain discovery plans, implementation plans,
execution logs, product or module plans, and agent-memory state. Some work is local to one
repository; other work is a portfolio thread that crosses several repositories.

The current Operating Kit scaffolds a `docs/agent-memory/goal-register.md`, but it does not define
a precise durable model for plan families, discovery documents, implementation documents,
cross-repo relationships, lifecycle metadata, creation or modification session references, or
coordination-home updates. As a result:

- a repository-local goal register can become an ad hoc task tracker or a duplicate plan index;
- formal plans do not have a clear canonical register under `docs/repo/plans/`;
- an HQ-level register can become a different mechanism instead of the same reusable standard;
- agents may not know when plan lifecycle changes should be recorded;
- useful session IDs can be lost even when they created or materially changed a plan;
- multi-repository work may depend on chat memory or individual discipline;
- hardcoded Codeheart HQ assumptions could leak into generic Operating Kit doctrine;
- automatic repository scanning could become brittle or surprising.

## User Intention

Create a simple Operating Kit standard that lets each repository maintain a clear plan register for
important discovery and implementation work, while optionally allowing a coordination home to keep
a portfolio-level view across repositories. The standard should reduce reliance on individual
developer discipline without making every agent session read or update a heavy project-management
system.

## Goals

- Define `docs/repo/plans/plan-register.md` as the central register for formal planning and
  portfolio coordination.
- Keep `docs/agent-memory/goal-register.md` available during transition for informal goals that do
  not yet have formal plans.
- Keep the plan register lightweight and durable enough to be useful across sessions.
- Avoid making the plan register a volatile next-action tracker.
- Use the same register mechanism in normal repositories and coordination-home repositories.
- Let a coordination home hold an aggregate top-level register without creating a separate HQ-only
  runbook or schema.
- Define when planning workflows should update local and coordination-home registers.
- Record lightweight session IDs when sessions create or materially modify plans, without adding
  session summaries to the register.
- Make portfolio coordination explicit, optional, and configuration-driven.
- Avoid default onboarding friction for users who do not need multi-repository coordination.
- Avoid silent GitHub or filesystem scanning.
- Keep repositories usable when the coordination home is unavailable.
- Preserve public-core safety in the Operating Kit repository.

## Non-Goals

- Do not implement the portfolio coordination feature in this discovery.
- Do not decide the exact CLI command names unless the recommendation is accepted during review.
- Do not create Codeheart-HQ-specific instructions in managed kit doctrine.
- Do not make HQ, GitHub, Microsoft 365, or any private Codeheart repository mandatory for generic
  Operating Kit consumers.
- Do not create a full task-management system, kanban board, sprint tracker, or detailed
  next-action ledger.
- Do not require normal first-run onboarding to ask every user about portfolio coordination.
- Do not silently scan sibling folders, GitHub organizations, or remote repositories.
- Do not automatically write across repositories without an explicit user-approved plan.
- Do not duplicate full plan status or execution evidence inside the register.

## Priority Order

1. Clear reusable plan-register semantics.
2. Low-maintenance lifecycle metadata and relations.
3. Lightweight session-reference capture for plan creation and material plan updates.
4. Same mechanism for local repositories and coordination homes.
5. Explicit optional portfolio coordination.
6. No surprising cross-repo writes or scans.
7. Graceful degradation when the coordination home is unavailable.
8. Compatibility with future CLI/schema automation.
9. Minimal changes to normal onboarding.

## Success Criteria For This Discovery

This discovery is ready for manual review when:

- the intended plan-register semantics are explicit enough to reject incompatible designs;
- all visible lifecycle and coordination decisions are captured;
- the relationship between plan-register entries and session IDs is explicit;
- open questions are marked as blocking or non-blocking;
- the relationship between local registers, HQ-style aggregate registers, and Operating Kit
  doctrine is clear;
- a later implementation plan can choose exact files, CLI commands, schema updates, tests, and
  release notes without re-opening the conceptual model.

# Section 2 - Current Model And Gaps

## Current Operating Kit Model

The current managed agent-memory documentation says:

- Operating Kit owns memory format and maintenance procedure.
- Consumer repositories own actual memory state.
- Typical consumer-owned files include `goal-register.md`, `session-ledger.md`,
  `untriaged-sessions.md`, and `archive/`.
- Agent memory should be used when a user asks where prior work left off, asks to continue a
  previous thread, or clearly references an ongoing workstream that may have session history.
- Agent memory should not replace canonical product, business, research, discovery,
  implementation, runbook, or reference docs.

The current goal-register scaffold only says "No active goals recorded yet" and invites curated
goal entries when ongoing work needs continuity across sessions.

## Gap Analysis

The current model is directionally right but under-specified for the way the repositories are now
being used.

| Gap | Practical consequence |
| --- | --- |
| Formal planning has no dedicated plan register under `docs/repo/plans/`. | Agents may overload `goal-register.md`, invent local indexes, or lose parent/child plan relationships. |
| Goal-register purpose is too vague. | Agents may turn it into a task list, session summary, or duplicate plan status table. |
| No durable entry schema exists for plan/workstream registration. | Cross-session and cross-repo work becomes inconsistent. |
| No lifecycle hook exists in planning workflows. | New discovery or implementation docs can be created without being registered anywhere. |
| No session-reference rule exists for sessions that create or materially change plans. | Valuable recovery handles can disappear even when the canonical plan survives. |
| No relation model exists for parent, child, superseded, and related plans. | Big portfolio threads become hard to recover after drilling into sub-discoveries. |
| HQ aggregate register is expected by portfolio planning, but no reusable mechanism exists. | HQ may invent a separate model instead of using Operating Kit doctrine. |
| Coordination home is not represented in kit config. | Agents cannot know where a portfolio-level register lives without chat memory or local assumptions. |
| No unavailable-coordination-home behavior exists. | Agents may fail, skip updates silently, or attempt risky discovery. |
| No enrollment model exists. | Updating several repositories depends on manual discipline or brittle scanning. |
| AGENTS.md cannot hardcode Codeheart HQ. | Generic consumers need a conditional route that depends on local config. |

## Consumer Impact Classification

This discovery document itself is repository governance documentation and has no immediate consumer
runtime impact.

Future implementation is likely to include a mix of:

- `instruction-only change` for managed agent-memory and planning-workflow doctrine;
- `backwards-compatible scaffold addition` if a consumer-owned `docs/repo/plans/plan-register.md`
  or `docs/repo/plans/coordination-sync-pending.md` scaffold is added for new consumers while
  existing consumer files are preserved;
- `instruction-only change` if `goal-register.md` is clarified as transitional or informal
  memory without changing existing consumer files;
- `validator-only change` if register or config validation is added;
- `consumer migration required` only if existing consumers are expected to run a repair, sync, or
  adoption command to add portfolio configuration or register structure.

# Section 3 - Target Operating Model

## Target Principle

A plan register is a lightweight, durable index of important planning and workstream authority. It
helps agents and humans find the right canonical documents, understand their relationships, and
recover the sessions that created or materially changed those plans. It does not replace discovery
docs, implementation plans, execution logs, session ledgers, or task trackers.

## Local Register

Every consumer repository can have its own `docs/repo/plans/plan-register.md`. It should list
important local discovery plans, implementation plans, plan families, major workstreams, and
selected cross-repo dependencies that materially affect local work.

`docs/agent-memory/goal-register.md` should remain available during the transition, but it should
not be the central planning or portfolio register. Its residual role is informal or pre-plan
continuity: goals that are worth remembering but have not yet become formal discovery,
implementation, or workstream plans.

The local register should answer:

- What important planning/workstream records exist here?
- What is each item about?
- Which canonical document owns the current details?
- Is this item active, draft, completed, superseded, or archived according to the canonical doc?
- Which repository owns it?
- Which parent, child, superseding, or related plans matter?
- What priority or ordering note helps orient a future agent?
- Which sessions created or materially updated the canonical plan, when a session ID is available?

It should not try to answer:

- What is the exact current next action?
- What happened in every recent chat session?
- What is the full execution status of every epic?
- What does the implementation plan or discovery doc already say in detail?
- What did each referenced session discuss in full?

## Coordination-Home Register

A coordination home is an optional repository that holds an aggregate top-level register for a
multi-repository portfolio. Codeheart HQ is the expected Codeheart example, but the Operating Kit
standard should not use HQ-specific terminology in generic managed content.

A coordination-home register should use the same entry schema as every other repository. The
difference is coverage:

- a normal repository register covers local plans and important cross-repo links;
- a coordination-home register covers portfolio-level plan families and selected entries from
  member repositories that matter to portfolio steering.

This avoids a special lower-detail HQ mechanism. The coordination home can include less detail by
choice, but it should not require a different schema or runbook.

## Register Entry Shape

Recommended required fields:

| Field | Purpose |
| --- | --- |
| `ID` | Stable local register identifier or canonical plan ID. |
| `Title` | Human-readable title. |
| `Type` | `discovery-plan`, `implementation-plan`, `plan-family`, `workstream`, or `reference-index`. |
| `Purpose` | Short explanation of what the entry is about. |
| `Owner / Repository` | Owning person, role, or repository. |
| `Canonical docs` | Repo-relative path or explicit repository/path pointer to the authoritative docs. |
| `Status` | Lifecycle snapshot copied from the canonical document: `draft`, `active`, `completed`, `superseded`, or `archived`. |
| `Created` | Creation date copied from the canonical document when available. |
| `Last updated` | Last updated timestamp copied from the canonical document when available. |
| `Priority / ordering note` | Stable orientation note, not a volatile next action. |
| `Relations` | Parent, child, supersedes, superseded-by, depends-on, blocks, or related entries. |
| `Session refs` | Lightweight session IDs for creation and material plan updates when available; no session summary required. |
| `Coordination note` | Optional note about portfolio-level relevance or sync state. |

Recommended optional fields:

| Field | Purpose |
| --- | --- |
| `Coverage note` | Explains whether the register is complete for a repository, portfolio, plan family, or intentionally selective. |
| `Last reviewed` | Last explicit register review date, distinct from canonical doc update time. |
| `Sync state` | Local-only marker such as `local-only`, `synced`, or `coordination-sync-pending`. |

Fields intentionally excluded from the required model:

- `Next action`
- detailed status narrative
- per-epic progress
- session summaries
- raw transcript bodies

Those details belong in canonical plans, execution logs, session ledgers, or user discussion.

## Register Entry Markdown Shape

The plan register should use one repeated section per plan or workstream entry, not one wide
table. Each entry should have a stable heading and compact field lines, followed by short lists
for relations, session refs, and coordination notes when present.

Recommended entry shape:

```md
## PR-001 - Operating Kit Portfolio Coordination And Plan Register Model

Type: discovery-plan
Status: draft
Owner / repository: Codeheart-Operating-Kit
Canonical docs:
`docs/repo/plans/portfolio-coordination-plan-register/portfolio-coordination-plan-register_discovery_doc.md`
Created: 2026-06-21
Last updated: 2026-06-21T12:11:05Z (UTC)
Priority / ordering note: Needed before portfolio coordination implementation.

Relations:
- parent: Codeheart portfolio restructure
- related: agent memory model

Session refs:
- created: <session-id-or-not-recorded>
- material update: <session-id-or-not-recorded>

Coordination note:
- Candidate for coordination-home aggregate register after implementation.
```

Rationale: tables become hard to maintain once relations, session refs, and coordination notes
grow. Repeated sections keep each plan readable, diffable, and easy for agents to update without
damaging adjacent entries.

## Session Reference Rule

The plan register should store session IDs as lightweight recovery handles, not as session
summaries.

When a discovery or implementation plan is created, the plan-register entry records the creating
session ID if available. When a plan is materially modified, the plan-register entry records the
modifying session ID if available.

Material modifications include changes to:

- scope;
- decision state;
- lifecycle status;
- parent, child, dependency, supersession, or related-plan links;
- readiness state;
- implementation path;
- review outcome.

Typos, formatting-only edits, timestamp refreshes, and mechanical checklist ticks do not require a
new session reference.

The register should not block work when a session ID is unavailable. The entry may omit the
session row or record `not recorded`. The session ledger remains useful for unplanned work,
abandoned threads, or general session recovery, but formal planned work should be recoverable first
through the plan register.

## Lifecycle Rule

The register may copy lifecycle metadata from canonical documents, but it is not the source of
truth. If the register and canonical plan disagree, the canonical plan wins and the register should
be refreshed.

Planning workflows should update the local register when a discovery or implementation plan is:

- created;
- materially updated in scope or decision state;
- linked as a parent, child, dependency, superseding, or related plan;
- marked active;
- completed;
- superseded;
- archived.

Minor typos, formatting-only edits, timestamp refreshes, and checklist progress that does not
change lifecycle or relationships should not force a register update.

## Portfolio Coordination Configuration

Portfolio coordination should be optional and explicit. A repository should know about a
coordination home through configuration, not through repository-name guesses or GitHub scanning.

Candidate configuration fields:

- `portfolio.role`: `coordination-home` or `member`
- `portfolio.coordination_home_path`
- `portfolio.coordination_home_register_path`
- `portfolio.member_repository_id`

Configuration semantics:

- No `portfolio` block means no configured portfolio coordination.
- `portfolio.role: member` means coordination hooks are active for this member repository.
- `portfolio.role: coordination-home` means this repository is the coordination home.
- The local member register path is inferred from `local_consumer_layer.repo_docs_path` plus
  `plans/plan-register.md`.
- The local pending-sync path is inferred from `local_consumer_layer.repo_docs_path` plus
  `plans/coordination-sync-pending.md`.
- The schema should not include `enabled`, `standalone`, `member_register_path`, or
  `pending_sync_path` unless a later implementation plan explicitly reopens that decision.

The exact JSON Schema encoding is an implementation decision. This discovery establishes that the
configuration should be explicit, non-secret, presence-based, and safe to store in
`.codeheart/kit.config.yaml` or a managed-schema-compatible consumer config location.

## Coordination Hooks In Planning Workflows

Managed planning workflows should include a small conditional hook:

1. Update the local plan register for material lifecycle, relationship, and session-reference
   changes.
2. If portfolio coordination is configured and the coordination home is locally available, update
   the coordination-home register as well.
3. If portfolio coordination is configured but the coordination home is unavailable, record a local
   pending-sync item in `docs/repo/plans/coordination-sync-pending.md` and continue without
   failing the primary task.

The hook should be short enough not to crowd normal discovery or implementation work. It should
route agents to the plan-register format and coordination runbook when needed.

## Managed AGENTS.md Treatment

The managed AGENTS.md block should not hardcode Codeheart HQ, repository names, or local paths. It
may include a generic conditional route such as:

- when portfolio coordination is configured, follow the configured coordination-home route for
  material planning lifecycle changes;
- otherwise maintain only the local repository register when the relevant runbook says to do so.

The exact text should be decided in implementation planning after the register and configuration
model are approved.

## Enrollment And Sync Model

A coordination home may offer an explicit workflow to discover or enroll member repositories, but
it should not silently scan or mutate repositories.

Recommended workflow:

1. User starts portfolio configuration or enrollment explicitly.
2. User provides local repository paths, or explicitly approves scanning a chosen parent folder.
3. Optional GitHub or organization lookup is only used when the user explicitly requests it and
   auth/network access is available.
4. The agent or CLI performs read-only inspection of candidates.
5. A plan is shown with proposed member role, register path, config changes, and any kit updates.
6. The user approves per repository or as a batch.
7. The workflow updates each approved repository, initializes or updates local registers, and
   updates the coordination-home register.

This model reduces individual developer discipline while preserving explicit consent and avoiding
brittle automatic discovery.

# Section 4 - Decision Inventory

## Decision Overview

| ID | Decision | Recommendation | Owner | Status |
| --- | --- | --- | --- | --- |
| `D-1` | One discovery or two | Use one discovery because plan-register semantics and portfolio coordination are coupled. | Operating Kit owner | recommended |
| `D-2` | Plan-register semantics | Define `docs/repo/plans/plan-register.md` as a lightweight plan/workstream authority index, not a task tracker. | Operating Kit owner | recommended |
| `D-3` | Register entry fields and format | Use repeated Markdown sections with title, type, purpose, owner/repo, canonical docs, lifecycle snapshot, priority note, relations, session refs, and coordination note. | Operating Kit owner | recommended |
| `D-4` | Lifecycle metadata authority | Copy created/status/last-updated from canonical docs; canonical docs win on conflict. | Operating Kit owner | recommended |
| `D-5` | HQ aggregate versus local register | Use the same schema and runbook; differ only by coverage and selection. | Operating Kit owner | recommended |
| `D-6` | Portfolio coordination in Operating Kit | Add optional advanced coordination-home support to the kit. | Operating Kit owner | recommended |
| `D-7` | Normal onboarding | Do not ask portfolio coordination questions during normal first-run onboarding. | Operating Kit owner | recommended |
| `D-8` | Configuration source | Use explicit non-secret configuration rather than inference or scanning. | Operating Kit owner | recommended |
| `D-9` | Member discovery and enrollment | Make it explicit, plan-first, and approval-gated. | Operating Kit owner | recommended |
| `D-10` | Planning workflow hooks | Add conditional local and coordination-home register update hooks. | Operating Kit owner | recommended |
| `D-11` | Managed AGENTS.md treatment | Add one lean generic conditional portfolio-coordination hook. | Operating Kit owner | recommended |
| `D-12` | Coordination home unavailable | Record pending sync locally and continue; do not fail the primary task. | Operating Kit owner | recommended |
| `D-13` | Codeheart migration path | Treat HQ as the first coordination-home consumer, but do not hardcode HQ into generic doctrine. | Codeheart owner | recommended |
| `D-14` | First implementation surface | Use runbooks, scaffolds, docs, and managed hooks first; defer CLI commands. | Operating Kit owner | recommended |
| `D-15` | Pending-sync storage | Store human-visible pending coordination sync in `docs/repo/plans/coordination-sync-pending.md`. | Operating Kit owner | recommended |
| `D-16` | Session-reference policy | Record creating and material-update session IDs in the plan register when available, without session summaries. | Operating Kit owner | recommended |
| `D-17` | Release and adoption shape | New installs scaffold the new files; existing consumers receive absent-file scaffolds only; existing files are preserved; release notes state additive no-forced-migration behavior. | Operating Kit owner | recommended |

## `D-1` One Discovery Or Two

Problem: Plan-register semantics and portfolio coordination could be separate discoveries.

Options:

- Split into a plan-register discovery and a portfolio-coordination discovery.
- Keep one discovery that starts with plan-register semantics and then defines optional
  coordination-home behavior.

Recommendation: keep one discovery.

Rationale: the main ambiguity is not just "how should HQ coordinate repositories?" It is "what is
the canonical planning register at all?" The coordination-home design depends directly on the
answer. Splitting the documents now would likely create duplicate decisions and inconsistent
terminology.

Status: recommended.

## `D-2` Plan-Register Semantics

Problem: The current `goal-register.md` name could mean a project-management backlog,
agent-memory summary, session tracker, plan index, or portfolio register. Formal planning needs a
clearer register.

Options:

- Keep `docs/agent-memory/goal-register.md` as the primary register.
- Rename the central concept to `docs/repo/plans/plan-register.md`.
- Make the plan register a lightweight index of important planning and workstream authority.
- Remove it from Operating Kit and leave every repository to decide.

Recommendation: define `docs/repo/plans/plan-register.md` as the primary lightweight index of
important planning and workstream authority.

Rationale: discovery docs, implementation plans, execution logs, and session ledgers already own
detailed state. The plan register should help agents find and relate those canonical records
without duplicating them. `docs/agent-memory/goal-register.md` can remain for informal or pre-plan
continuity during transition.

Status: recommended.

## `D-3` Register Entry Fields And Format

Problem: A useful register needs enough fields to orient agents, but too many fields make it stale.
The Markdown shape also matters because relations, session refs, and coordination notes do not fit
well into a single wide table.

Options:

- Minimal list of names and links.
- Full project-management table with status, next action, owner, blockers, and progress.
- Wide table with one row per plan.
- Repeated Markdown sections with one section per plan or workstream entry.

Recommendation: use repeated Markdown sections with one section per plan or workstream entry.
Each section should include compact fields plus short lists for relations, session refs, and
coordination notes.

Rationale: title, purpose, owner/repo, canonical docs, lifecycle snapshot, priority note, and
relations provide durable orientation. Session IDs provide recovery handles without requiring
session summaries. A repeated-section format is easier to read and update than a wide table.
Required next action and detailed status would require constant updates and duplicate canonical
plans.

Status: recommended.

## `D-4` Lifecycle Metadata Authority

Problem: The register needs status and dates, but the canonical docs already carry created,
status, and last-updated headers.

Options:

- Treat the register as source of truth for lifecycle.
- Copy lifecycle metadata into the register as an index snapshot.
- Omit lifecycle metadata from the register entirely.

Recommendation: copy lifecycle metadata into the register as an index snapshot; canonical docs
win on conflict.

Rationale: seeing status and dates in the register helps prioritize and orient. Treating the
register as authoritative would create conflict with canonical docs and implementation execution
logs.

Status: recommended.

## `D-5` HQ Aggregate Versus Local Register

Problem: A coordination home such as Codeheart HQ needs a portfolio-level view. It could use a
different mechanism or the same plan-register model.

Options:

- Create a special HQ portfolio register with a different schema.
- Use the same plan-register schema and runbook, but let coverage differ.
- Keep all registers local and make HQ link manually to each repository.

Recommendation: use the same schema and runbook, with different coverage.

Rationale: a separate HQ mechanism would create another concept for agents to learn. The same
model can support full local detail in member repositories and selective top-level coverage in a
coordination home.

Status: recommended.

## `D-6` Portfolio Coordination In Operating Kit

Problem: Multi-repository coordination could be a Codeheart-only local pattern or a generic
Operating Kit feature.

Options:

- Keep it entirely local to Codeheart HQ.
- Add optional advanced coordination-home support to Operating Kit.
- Make coordination homes mandatory in the standard profile.

Recommendation: add optional advanced support to Operating Kit.

Rationale: the pattern is reusable for any multi-repository operating context, but it should not
burden simple consumers.

Status: recommended.

## `D-7` Normal Onboarding

Problem: If portfolio coordination is useful, the first-run onboarding could ask about it.

Options:

- Ask every new user during standard onboarding.
- Hide it completely until a future major version.
- Make it an advanced setting or explicit command.

Recommendation: make it an advanced setting or explicit command, not a normal onboarding prompt.

Rationale: most initial users should not have to understand multi-repository coordination before
using the kit. Advanced configuration remains available for consumers who need it.

Status: recommended.

## `D-8` Configuration Source

Problem: Agents need to know whether a repository has a coordination home.

Options:

- Infer from repository names such as `HQ`.
- Scan sibling folders or GitHub organizations.
- Store explicit non-secret configuration.

Recommendation: store explicit non-secret configuration.

Rationale: inference is brittle, and scanning is surprising. Config gives agents a deterministic
route without requiring every prompt to mention the coordination home.

Status: recommended.

## `D-9` Member Discovery And Enrollment

Problem: A coordination-home owner may want to bring several repositories under the same portfolio
without relying on individual developers to configure each repo manually.

Options:

- Do nothing; configure every repository separately.
- Automatically scan and update all nearby or GitHub-visible repositories.
- Offer an explicit read-only discovery, plan, and approval-gated enrollment workflow.

Recommendation: offer explicit read-only discovery, plan, and approval-gated enrollment.

Rationale: this reduces manual discipline while preserving control over cross-repo writes and
dirty worktrees.

Status: recommended.

## `D-10` Planning Workflow Hooks

Problem: Agents need to know when register updates are expected.

Options:

- Leave register updates to user memory.
- Add a broad instruction to update the register after every task.
- Add targeted hooks to discovery, implementation planning, and implementation execution
  workflows.

Recommendation: add targeted hooks for material plan lifecycle, relationship, and session-reference
changes.

Rationale: planning workflows are the natural point where plans are created, activated,
completed, superseded, archived, linked, or materially updated. Updating after every task would be
noisy.

Status: recommended.

## `D-11` Managed AGENTS.md Treatment

Problem: Some routing may need to be visible immediately in `AGENTS.md`, but the managed block
should stay lean.

Options:

- Hardcode the coordination-home route in `AGENTS.md`.
- Add one lean generic conditional route.
- Keep all details in managed runbooks and config references.

Recommendation: add one lean generic conditional portfolio-coordination hook to managed
`AGENTS.md`, and keep details in runbooks and config references.

Rationale: `AGENTS.md` should remain lean. It should tell agents that configured portfolio
coordination exists and that material planning lifecycle changes may require local and
coordination-home register updates. It should not embed HQ-specific paths or full procedures.

Status: recommended.

## `D-12` Coordination Home Unavailable

Problem: A member repository may know its coordination home, but that repository may not be
checked out locally or accessible during the session.

Options:

- Fail the current planning task.
- Skip coordination silently.
- Record pending sync locally and continue.

Recommendation: record pending sync locally and continue.

Rationale: the primary repo must remain usable without the coordination home. Silent skips make
coordination unreliable. A pending-sync marker preserves the need without blocking local work.

Status: recommended.

## `D-13` Codeheart Migration Path

Problem: Codeheart HQ is the expected first coordination home, but Operating Kit must stay generic
and public-safe.

Options:

- Encode Codeheart HQ as the standard.
- Treat Codeheart HQ as an example consumer of a generic standard.
- Keep Codeheart HQ outside the Operating Kit discussion.

Recommendation: treat Codeheart HQ as the first example consumer, not as generic doctrine.

Rationale: this keeps public kit doctrine reusable and avoids leaking private repository shape
into the public core while still supporting Codeheart's immediate need.

Status: recommended.

## `D-14` First Implementation Surface

Problem: Portfolio coordination may be exposed through CLI commands, runbooks, scaffolds, managed
docs, or a mix of these. A first implementation should not overbuild automation before the
operating model has been used.

Candidate options:

- Runbooks, scaffolds, managed docs, and managed AGENTS hook first.
- Add `codeheart-operating-kit portfolio configure` immediately.
- Add `codeheart-operating-kit settings portfolio configure` immediately.
- Extend `inspect` or `sync` before the register model has been piloted.

Recommendation: use runbooks, scaffolds, managed docs, and a lean managed AGENTS hook first. Defer
CLI commands until the model has been used and the command contract is clearer.

Status: recommended.

## `D-15` Pending-Sync Storage

Problem: When coordination home updates cannot be made, the pending work needs a home.

Candidate options:

- a section inside `docs/repo/plans/plan-register.md`;
- a dedicated consumer-owned planning file at `docs/repo/plans/coordination-sync-pending.md`;
- a separate consumer-owned file under `docs/agent-memory/`;
- a generated local report under `.codeheart/`;
- a field in `.codeheart/kit.lock.yaml`.

Recommendation: store human-visible pending coordination sync in
`docs/repo/plans/coordination-sync-pending.md`.

Rationale: pending coordination sync is planning-adjacent state, not managed kit state and not
general agent memory. A dedicated file avoids bloating the plan register while keeping the pending
work discoverable. When pending sync becomes substantial or durable, the pending-sync file itself
can be registered in `docs/repo/plans/plan-register.md`.

Status: recommended.

## `D-16` Session-Reference Policy

Problem: Useful session IDs can be lost even when a session creates or materially changes a
canonical plan. At the same time, a plan register should not become a session ledger or summary
file.

Options:

- Do not record session IDs in the plan register.
- Record only subjectively "notable" sessions.
- Record the creating session and material-update sessions when session IDs are available.
- Merge the session ledger into the plan register.

Recommendation: record the creating session and material-update sessions when session IDs are
available.

Rationale: this avoids subjective judgment about which sessions are notable. A session is relevant
because it created or materially changed a canonical plan. The entry should be lightweight: session
ID, date, and reason such as `created` or `material update`. No session summary is required because
the canonical plan owns the summarized state.

Status: recommended.

## `D-17` Release And Adoption Shape

Problem: The first implementation can affect new installs, existing consumers, or both. It must
make the new plan-register model available without overwriting consumer-owned state or forcing
immediate migration work.

Options:

- New installs only; existing consumers do not receive plan-register scaffolds.
- Existing consumers receive absent-file scaffolds during sync, with no overwrites.
- Existing consumers are forced through a migration that rewrites or moves current goal-register
  content.

Recommendation: new installs scaffold `docs/repo/plans/plan-register.md` and
`docs/repo/plans/coordination-sync-pending.md`. Existing consumers may receive those files through
sync only when the files are absent. Existing `docs/agent-memory/goal-register.md`, existing
`docs/repo/plans/plan-register.md`, and existing `docs/repo/plans/coordination-sync-pending.md`
content must be preserved exactly unless the user explicitly requests a migration or repair.

Release notes should describe this as an additive scaffold and doctrine change, not a forced
migration. Managed docs should clarify the transitional role of `goal-register.md`; sync should not
move, rewrite, or archive existing goal-register content.

Rationale: this gives new and existing consumers the same available planning surface while keeping
consumer-owned state safe. It also avoids making the release feel like a mandatory cleanup project.

Status: recommended.

# Section 5 - Requirements And Constraints

## Functional Requirements

| ID | Requirement | Priority |
| --- | --- | --- |
| `FR-1` | Define `docs/repo/plans/plan-register.md` as a lightweight plan/workstream authority index. | must |
| `FR-2` | Provide a reusable repeated-section entry format for discovery plans, implementation plans, plan families, and major workstreams. | must |
| `FR-3` | Preserve local repository ownership of agent-memory state. | must |
| `FR-4` | Support the same register model for normal repositories and coordination-home repositories. | must |
| `FR-5` | Make portfolio coordination optional and explicit. | must |
| `FR-6` | Store coordination-home configuration in a non-secret, deterministic place. | must |
| `FR-7` | Add planning-workflow hooks for material plan lifecycle, relationship, and session-reference changes. | must |
| `FR-8` | Keep normal onboarding free of mandatory portfolio-coordination questions. | should |
| `FR-9` | Provide an advanced configuration or enrollment path for multi-repository consumers. | should |
| `FR-10` | Do not silently scan or mutate sibling repositories or GitHub repositories. | must |
| `FR-11` | Require plan-first, approval-gated cross-repo enrollment or updates. | must |
| `FR-12` | Record pending sync in `docs/repo/plans/coordination-sync-pending.md` when coordination-home update is due but unavailable. | should |
| `FR-13` | Add one lean managed `AGENTS.md` hook for configured portfolio coordination while keeping details in runbooks. | must |
| `FR-14` | Preserve existing consumer `docs/agent-memory/goal-register.md` files during sync and keep them available for informal or pre-plan continuity. | must |
| `FR-15` | Include release notes or migration/adoption notes when shipped to consumers. | must |
| `FR-16` | Record creating and material-update session IDs in plan-register entries when available, without requiring session summaries. | must |
| `FR-17` | Do not let missing session IDs block planning work. | must |
| `FR-18` | Keep the first implementation runbook/scaffold/docs first and do not add a portfolio CLI command unless a later implementation plan explicitly expands scope. | should |
| `FR-19` | New installs scaffold `docs/repo/plans/plan-register.md` and `docs/repo/plans/coordination-sync-pending.md`. | must |
| `FR-20` | Existing consumers may receive those files through sync only when absent; sync must not overwrite existing consumer-owned plan-register, pending-sync, or goal-register content. | must |
| `FR-21` | Release notes must describe the change as additive scaffold and doctrine behavior with no forced migration. | must |

## Non-Functional Requirements

| ID | Requirement | Priority |
| --- | --- | --- |
| `NFR-1` | Low overhead: register maintenance should not become a mandatory step for every small edit. | must |
| `NFR-2` | Durable orientation: a fresh agent should understand plan relationships without reading every plan first. | must |
| `NFR-3` | Source-of-truth clarity: canonical plans beat register snapshots on conflict. | must |
| `NFR-4` | Public-core safety: generic kit doctrine must not expose private Codeheart repository details. | must |
| `NFR-5` | Cross-repo safety: no cross-repo writes without explicit approval. | must |
| `NFR-6` | Offline/local resilience: member repositories remain usable without the coordination home checkout. | must |
| `NFR-7` | Simplicity: avoid a full project-management system. | must |
| `NFR-8` | Schema compatibility: future config should be machine-readable and validator-friendly. | should |
| `NFR-9` | Discoverability: agents should find the register rules through existing Operating Kit routes. | must |
| `NFR-10` | Session references stay lightweight: IDs, dates, and reasons only, not session summaries. | must |

## Constraints

- The Operating Kit repository is public-core only.
- `.codeheart/kit/` remains managed content and is not a place for consumer-authored portfolio
  state.
- `docs/agent-memory/` remains consumer-owned state and must not be overwritten by sync after
  creation.
- `docs/repo/plans/plan-register.md` and `docs/repo/plans/coordination-sync-pending.md` are
  kit-initialized consumer state files: Operating Kit owns the file contract and presence behavior,
  while the consumer repository owns the contents after creation.
- Future baseline/template changes must preserve existing consumer file contents.
- Existing consumer-owned register contents win over new baseline templates.
- Normal onboarding should remain short and usable for non-technical first-run users.
- Coordination-home paths may be unavailable, moved, or absent from the current checkout.
- GitHub authentication, network access, and organization membership cannot be assumed.
- The feature should work without requiring HQ or any Codeheart-private repository.
- Session IDs are recovery handles, not source-of-truth status records.

# Section 6 - Open Questions And Assumptions

## Open Questions

| ID | Question | Owner | BLOCKER | Current default |
| --- | --- | --- | --- | --- |
| `OQ-1` | What exact future CLI command should configure portfolio coordination if CLI support is added later? | Operating Kit owner | no | Defer CLI in the first implementation; revisit after runbook/scaffold use. |
| `OQ-2` | What exact YAML shape should be added to `.codeheart/kit.config.yaml`? | Operating Kit owner | no | Use optional presence-based `portfolio.role` config: `member` requires member identity plus coordination-home path/register; `coordination-home` requires its register path. Infer local member register and pending-sync paths from `local_consumer_layer.repo_docs_path`. |
| `OQ-3` | Should the first implementation include validation for plan-register shape? | Operating Kit owner | no | Start with docs/scaffold/runbook unless implementation scope justifies a validator. |
| `OQ-4` | Should managed `AGENTS.md` include a new conditional coordination sentence in the first release? | Operating Kit owner | no | Keep `AGENTS.md` lean; add only if runbook routing alone is insufficient. |
| `OQ-5` | Which Codeheart repositories should be enrolled first once the feature exists? | Codeheart owner | no | HQ as coordination home; AWS Platform, Foundry, and possibly Operating Kit as members after review. |
| `OQ-6` | Should the Operating Kit repo itself be a member of the Codeheart portfolio coordination home? | Codeheart owner | no | Likely yes for Codeheart internal visibility, but not required by generic kit doctrine. |
| `OQ-7` | Should `goal-register.md` be removed immediately? | Operating Kit owner | no | No. Keep it during transition for informal or pre-plan continuity. |

## Assumptions

| ID | Assumption | Validation path |
| --- | --- | --- |
| `A-1` | Existing consumers may already have `docs/agent-memory/goal-register.md`, so sync must not overwrite it. | Confirm current sync/scaffold behavior before implementation. |
| `A-2` | The user values a low-maintenance register more than detailed live task tracking. | Manual review of this discovery. |
| `A-3` | A coordination home can be represented as a repository path and register path in non-secret config. | Schema and CLI design review. |
| `A-4` | Planning runbooks are the right place for lifecycle update hooks. | Review managed discovery, implementation-planning, and execution workflows during implementation planning. |
| `A-5` | Codeheart HQ will remain the first concrete coordination-home consumer, but not the generic standard. | Codeheart portfolio review. |
| `A-6` | Public Operating Kit doctrine can describe coordination-home patterns without exposing private HQ details. | Public-core hygiene validation. |
| `A-7` | Session IDs are available often enough to be useful but not guaranteed in every agent environment. | Implementation review and consumer pilots. |
| `A-8` | Existing consumers should benefit from the new files when absent, but should not be forced into a migration. | Sync behavior tests and release-note review. |

# Section 7 - Risks And Validation

## Risks

| ID | Risk | Likelihood | Impact | Mitigation | Detection |
| --- | --- | --- | --- | --- |
| `R-1` | Plan register becomes a stale project-management table. | medium | high | Exclude required next action and detailed status; route details to canonical plans. | Register rows contain long status narratives or frequent task churn. |
| `R-2` | Agents forget to update plan registers. | medium | medium | Add targeted planning-workflow hooks and optional AGENTS routing. | New or completed plans missing from local plan register. |
| `R-3` | HQ develops a special schema that fragments the model. | medium | medium | Use one schema for local and coordination-home registers. | HQ register requires separate runbook for normal operations. |
| `R-4` | Cross-repo enrollment overwrites unrelated user changes. | low | high | Require read-only inspect, plan, dirty-worktree checks, and user approval before writes. | Enrollment review finds unapproved file changes. |
| `R-5` | Coordination-home config hardcodes Codeheart-specific paths in public kit doctrine. | medium | high | Keep examples generic or sanitized; validate public-core hygiene. | Public-core review finds private names, tenant details, or local paths. |
| `R-6` | Coordination-home absence blocks member repository work. | medium | medium | Use `docs/repo/plans/coordination-sync-pending.md` and continue local task. | Planning tasks fail only because HQ is unavailable. |
| `R-7` | Normal onboarding becomes too complex for simple consumers. | low | medium | Keep portfolio configuration advanced and opt-in. | First-run tests or pilots show confusion before setup completes. |
| `R-8` | Lifecycle snapshots drift from canonical docs. | high | low | State canonical docs win; refresh during planning hooks and maintenance. | Plan-register status disagrees with canonical plan header. |
| `R-9` | Silent GitHub or filesystem scanning surprises users. | low | high | Require explicit scan scope and approval. | CLI or runbook can discover repos without a user-provided scope. |
| `R-10` | Session references bloat into session summaries. | medium | medium | Store only ID, date, and reason; keep session summaries in canonical docs or session ledger. | Plan-register entries contain narrative session summaries. |
| `R-11` | Existing consumers experience the change as a forced migration. | low | high | Use absent-file scaffolds only, preserve existing files, and state no-forced-migration behavior in release notes. | Sync rewrites or moves existing goal-register, plan-register, or pending-sync files. |

## Validation Criteria For A Future Implementation Plan

A future implementation plan should include validation that:

- managed planning and agent-memory docs define plan-register semantics and the transition role of
  `goal-register.md`;
- new kit-initialized consumer state baselines can create `docs/repo/plans/plan-register.md` and
  `docs/repo/plans/coordination-sync-pending.md` when absent;
- existing consumer goal-register and plan files are preserved by sync;
- existing consumers receive new plan-register and pending-sync files only when absent;
- plan-register scaffolds use repeated sections, not one wide table;
- planning workflow runbooks include only targeted lifecycle, relationship, and session-reference
  hooks;
- managed `AGENTS.md` includes one lean conditional portfolio-coordination hook;
- portfolio coordination config is explicit and non-secret;
- no default onboarding path asks portfolio questions and no first implementation adds CLI unless
  explicitly re-scoped during implementation planning;
- no workflow silently scans or writes other repositories;
- unavailable coordination home creates visible pending sync state in
  `docs/repo/plans/coordination-sync-pending.md` and does not fail local planning work;
- creating and materially modifying sessions are recorded as lightweight session refs when IDs are
  available;
- README routers and managed inventory routes make the model discoverable;
- public-core hygiene checks pass;
- release notes describe consumer impact and adoption path, including that the change is additive
  and no forced migration is required.

# Section 8 - Review Handoff

## Current Recommendation

Proceed with manual review of this discovery before drafting an implementation plan.

The recommended direction is:

- add `docs/repo/plans/plan-register.md` as the lightweight plan/workstream authority index;
- keep `docs/agent-memory/goal-register.md` during transition for informal or pre-plan continuity;
- use one reusable schema for local repositories and coordination-home repositories;
- use repeated Markdown sections with one section per plan or workstream entry;
- avoid required next-action and detailed-status fields;
- copy lifecycle metadata from canonical docs as an index snapshot only;
- record creating and material-update session IDs as lightweight session refs when available;
- add targeted planning-workflow hooks for material lifecycle, relationship, and session-reference
  changes;
- make portfolio coordination an optional advanced Operating Kit feature;
- configure coordination homes explicitly rather than inferring from names or scanning;
- use an approval-gated enrollment/sync workflow for multi-repository setup;
- add one lean managed `AGENTS.md` hook for configured portfolio coordination;
- defer portfolio CLI commands in the first implementation;
- record pending sync in `docs/repo/plans/coordination-sync-pending.md` if the coordination home is
  unavailable.
- scaffold `docs/repo/plans/plan-register.md` and
  `docs/repo/plans/coordination-sync-pending.md` for new installs and absent existing-consumer
  files only, preserving all existing consumer-owned files.

## Implementation-Planning Readiness

Ready for implementation-plan drafting after manual review of this discovery. The core register
model, repeated-section entry format, runbook-first implementation surface, managed AGENTS hook,
pending-sync storage, session-reference rule, and additive release/adoption shape now have
recommended defaults. Remaining open questions are non-blocking implementation-shaping choices
that can be resolved inside the implementation plan.

## Likely Implementation Areas

A later implementation plan will likely touch:

- `src/codeheart_operating_kit/resources/components/agent-memory/managed/`
- `src/codeheart_operating_kit/resources/components/agent-memory/scaffolds/goal-register.md`
- consumer-owned plan scaffolds for `docs/repo/plans/plan-register.md` and
  `docs/repo/plans/coordination-sync-pending.md`
- `src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/`
- `src/codeheart_operating_kit/resources/components/agent-interface/managed/`
- `schemas/kit-config.schema.json`
- `tests/`
- `docs/README.md`
- `docs/repo/README.md`
- release notes and consumer-impact records

# Revision Notes

- 2026-06-21: Initial draft created from the Codeheart multi-repository coordination and
  goal-register discussion. Captured the first reusable register model, optional coordination-home
  pattern, targeted planning workflow hooks, and open implementation-shaping questions before the
  later plan-register refinement.
- 2026-06-21: Revised the central model from `goal-register.md` to
  `docs/repo/plans/plan-register.md`; kept `goal-register.md` as transitional informal memory;
  selected `docs/repo/plans/coordination-sync-pending.md` for pending coordination sync; and added
  the rule that creating and materially modifying sessions should be recorded as lightweight
  session refs when session IDs are available.
- 2026-06-21: Added final pre-implementation planning decisions: plan-register entries use one
  repeated Markdown section per plan or workstream; the first implementation should use runbooks,
  scaffolds, managed docs, and a lean managed `AGENTS.md` hook before adding CLI commands; and the
  managed `AGENTS.md` hook should remain generic and route details to runbooks.
- 2026-06-21: Added release/adoption decision: new installs receive plan-register and
  coordination-sync-pending scaffolds; existing consumers receive those files only when absent;
  existing consumer-owned register files are preserved; and release notes must describe the change
  as additive with no forced migration.
