Last updated: 2026-08-09T15:24:37Z (UTC)

# Planning Document Lifecycle

Planning documents use lifecycle metadata and predictable placement so agents and maintainers can
tell whether a document is draft material, active execution authority, completed traceability,
superseded guidance, or archived history.

## Planning Document Types

- `*_discovery_doc.md`: decision and requirement discovery.
- `*_implementation_doc.md`: execution-ready implementation plan.
- exact `README.md` with valid family metadata: shared discoverability for independent sibling
  plans.
- `*_execution_log.md`: goal-style implementation evidence and divergence log.
- Plan-scoped attachments: temporary inventories, reviews, reports, and evidence used by one plan.

Discovery, implementation, and family are separate canonical record kinds. A discovery and its
implementation may share a family and relations, but neither is merely a phase row inside the
other. Execution logs and attachments are supporting evidence, not independent plan records.

## Required Planning Header

Discovery and implementation documents start with:

```text
Last updated: YYYY-MM-DDTHH:MM:SSZ (UTC)
Created: YYYY-MM-DD
Status: draft | active | completed | superseded | archived
```

Use the current UTC clock for `Last updated`. Preserve `Created` after the document is created.

## Canonical Metadata

In mixed and canonical mode, place the bounded Codeheart plan metadata block immediately below the
canonical Markdown title. Under discovery v2, canonical discovery and implementation records
always require both their supported filename and valid metadata; an explicit family record always
requires exact `README.md` and valid family metadata. Metadata owns stable semantic identity,
public-safe purpose, catalog chronology, optional family/classifications/relations, and legacy
aliases. Filename or family qualification owns the human-readable kind.

Use `plan-catalog-format.md` for the exact marker and fields. The lifecycle header continues to own
`Status`; do not add a separate authority field. Metadata-only adoption updates
`catalog_metadata_updated` and must not change the historical `Last updated` value.

## Optional Lifecycle Fields

Use these fields when they help review:

```text
Completed: YYYY-MM-DD
Superseded by: <relative/path>
Execution log: <relative/path>
```

Add `Execution log:` when a sibling execution log exists.

## Status Values

- `draft`: still being written and not approved for execution.
- `active`: current and approved for execution.
- `completed`: executed and retained for current context or traceability.
- `superseded`: replaced by another plan.
- `archived`: retained only as historical context.

Use one status value. Do not execute a `draft`, `completed`, `superseded`, or `archived`
implementation plan without explicit user approval.

The lifecycle value is not a branch, push, pull-request, merge, or release state. An `active` plan
is approved execution authority only within current user instruction, repository policy, and its
own scope.

## Discovery, Catalog Modes, And Authoring

Discovery v1 recognizes formal records beneath `docs/repo/plans/`. Discovery v2 recognizes
Git-tracked repository-owned Markdown in any path whose slash-normalized segments contain an exact
lowercase `docs` at arbitrary depth. A supported filename or genuine plan metadata discovers a v2
candidate; canonical validity requires both signals for discovery and implementation records.

Existing repositories stay on v1 until an explicit reviewed v2 activation. A Kit upgrade alone
does not activate v2. Use prospective v2 inventory and validation before changing config. Untracked
authoring preview is optional and never canonical authority.

- `legacy`: existing plans may lack metadata and the legacy register may still be maintained.
- `mixed`: every new or materially rewritten formal plan receives metadata; only exact
  pre-cutover register-linked plans may remain legacy temporarily. A genuine active branch owner
  may remain legacy only through an exact reviewed schema-v3 deferral and pending incremental
  follow-up.
- `canonical`: every formal discovery, implementation, and qualifying family record requires
  metadata; branch-owned deferral is forbidden.

Read mode from `.codeheart/kit.config.yaml`. Do not assume that an absent setting means canonical;
absence defaults to legacy. Mixed is optional: a legacy repository that can migrate all
discovery-v2 candidates may move directly to canonical after the guarded migration. Frozen
register/cutover proof is needed only when filename-only records remain temporarily grandfathered
in mixed mode. Use the migration runbook for mode or discovery-version changes.

### Config-V2 Evidence Binding

When mixed activation retains a `deferred-active-owner`, `.codeheart/kit.config.yaml` uses schema
version 2 and `plan_catalog_migration_evidence` binds the tracked ledger path and SHA-256,
branch-evidence digest, evidence revision `E`, ledger checkpoint `L` as
`activation_base_revision`, migration-action digest, and `local` or `remote-aware` evidence scope.
The frozen cutover revision, exact deferred path/source hash, and owner ref/tip/path/mode/object/
hash snapshot remain in the reviewed ledger.

That binding is required only while a deferred follow-up remains and is forbidden when mixed mode
has no such gap. Direct canonical activation uses the zero-gap contract and must not preserve a
deferral. A consumer CLI that does not understand config schema v2 must reject the config before
use; ignoring unknown evidence fields would silently erase the reason the legacy row remains safe.

### Deferred Owner Follow-Up Lifecycle

Deferral does not change the plan's `Status:` or grant execution authority. The legacy target-row
lifecycle comes from the frozen cutover document; the owner-tip snapshot must independently prove
an `active` candidate. The catalog shows the plan exactly once with
`coverage_disposition: mixed-grandfathered` and `incremental_migration_required: true`.

The reviewed follow-up stays `pending` with action `reinventory-and-incremental-migrate` and trigger
`owner-ref-integrated-or-target-path-changed`. After the owner ref is integrated or the target path
or bytes change, validation must either accept owner-supplied canonical metadata at merge or apply a
new reviewed ledger to the merged target bytes. Until then the member may be
`mixed_coverage_complete` but cannot be `canonical_ready`. A moved ref or any stale bound input is a
blocker, not permission to keep using the old snapshot.

## Revision Notes

Discovery and implementation documents maintain a bottom `# Revision Notes` section for
meaningful decision, scope, strategy, or execution-plan changes.

Do not add revision notes for:

- timestamp-only edits;
- typos;
- formatting-only edits;
- checklist progress that does not change scope.

## Checkbox State

Implementation documents may use checkbox tasks as execution state. Tick a checkbox only after the
work is completed and validated. Do not tick future work, deferred work, or work that has not been
validated.

Checkbox state is a human-readable execution aid. It does not replace issue state, pull-request
state, release state, or execution logs.

## Execution Logs

Goal-style implementation runs use a sibling execution log:

```text
<feature-slug>_implementation_doc.md
<feature-slug>_execution_log.md
```

For a plan bundle:

```text
<feature-slug>/
  <feature-slug>_implementation_doc.md
  <feature-slug>_execution_log.md
```

The execution log records:

- meaningful divergence from the plan;
- added tasks;
- safe defaults;
- validation substitutions or extensions;
- review-gate rounds and accepted result;
- release, migration, or handoff evidence;
- residual risk and follow-ups.

The execution log is not a command transcript.

## Plan Shapes

Use the smallest shape that preserves ownership, reviewability, and lifecycle clarity.

### Standalone Plan File

Use a standalone plan file for small, low-artifact work unlikely to need an execution log,
attachments, or repeated goal-style execution.

```text
<owned-docs-root>/
  plans/
    <feature-slug>_implementation_doc.md
```

`<owned-docs-root>` may be root `docs/` or an arbitrarily nested product, package, source-area,
business, or domain path ending in an exact `docs` segment. Choose placement from the owning
documentation tree; do not centralize plans merely for catalog discovery.

### Plan Bundle

Use a plan bundle when the plan has or may need:

- discovery plus implementation documents;
- execution log;
- attachments;
- repeated goal-style execution;
- multiple epics;
- cross-area review.

```text
<owned-docs-root>/plans/
  <feature-slug>/
    <feature-slug>_discovery_doc.md
    <feature-slug>_implementation_doc.md
    <feature-slug>_execution_log.md
    attachments/
      <plan-scoped-artifact>.md
```

### Subplan

Use `subplans/` when a child plan exists mainly to edit, audit, split, repair, or execute a parent
plan, or when the parent defines the child plan's correctness.

```text
<parent-plan-slug>/
  <parent-plan-slug>_discovery_doc.md
  subplans/
    <subplan-slug>/
      <subplan-slug>_implementation_doc.md
      <subplan-slug>_execution_log.md
      attachments/
```

Do not use `subplans/` merely because plans share a topic. Link related plans instead.

### Plan Family

Use a plan family when independently executable sibling plan bundles need shared discoverability
without one plan owning the others.

```text
plans/
  <family-slug>/
    README.md
    <first-plan-slug>/
      <first-plan-slug>_implementation_doc.md
    <second-plan-slug>/
      <second-plan-slug>_implementation_doc.md
```

The second related sibling plan bundle is the normal trigger for a shared family folder.
Directory shape is only a structural validation aid. It must never qualify the README or create
family authority without genuine valid family metadata.

### Program Folder

Use a program folder only when a parent coordination plan owns multiple workstream plans and has
current authority over scope, sequencing, review state, or cross-workstream decisions.

```text
<program-slug>/
  <program-slug>_implementation_doc.md
  <program-slug>_execution_log.md
  plans/
    <workstream-plan-slug>/
      <workstream-plan-slug>_implementation_doc.md
```

## Plan-Scoped Attachments

Put plan-scoped artifacts in the plan bundle's `attachments/` folder when the material is too
large, too detailed, too volatile, or too audit-specific to keep in the plan.

Examples:

- path migration maps;
- migration inventories;
- stale-path reports;
- one-off review notes;
- temporary decision matrices;
- validation inventories;
- migration ledgers.

Promote an attachment to a durable `reference/` document only when it becomes reusable doctrine.

## Archive Behavior

Use `archive/` for superseded or historical plans retained for traceability. Do not move a plan to
archive during execution unless the user explicitly asks or the active plan includes an approved
archive task.

Archived plans remain historical by default. New active work should link to archived history rather
than adding active plans under `archive/`.

## Plan Views, Register, And Git Visibility

`docs/repo/plans/plan-register.md` remains the stable compatibility entry point even when actual
plans live in multiple owned documentation trees. In mixed and canonical mode the current
repository view is generated by `codeheart-operating-kit plans list`; the register is not a
manually appended second catalog. A pre-existing register remains byte-frozen legacy evidence only
when mixed-mode grandfathering uses its cutover revision. Canonical migrations may retain an old
register as labeled historical evidence without creating a new cutover.

A deferred active-owner row is derived from the frozen legacy entry plus the exact config-v2/
ledger-v3 binding. Do not edit the register, append a deferral notice, synthesize a generated row,
or create a second cutover to make it visible. The compatibility row remains visible locally and
remotely until its incremental follow-up succeeds; silent omission is a validation failure.

Canonical planning documents remain source of truth for lifecycle, scope, decisions, metadata,
execution evidence, and handoff details. Git provides path, ref, revision, and content identity.
Use `plan-register-format.md` and `../runbooks/maintain-plan-register.md` for mode-specific behavior.

Portfolio coordination sees only pushed remote refs and uses the same repository-wide classifier
and candidate set as local list, validate, inventory, and migrate. Default-branch configuration
controls remote discovery; feature branches cannot broaden authority. A plan created or activated
on a work branch becomes globally observable after its authorized normal push; worktree edits,
local heads, and unpushed commits remain local. The coordination home's scanner reports freshness
and completeness and must not present a failed refresh or a v1/older cache as v2-current.

Review-only workflows report missing/invalid metadata or stale views as findings. They do not
silently change lifecycle, register bytes, metadata, config mode, or remote visibility.

## Index Maintenance

Update the nearest README and parent index when discoverability changes:

- new plan, runbook, reference, or attachment that readers need to find;
- moved, renamed, archived, or removed document;
- changed title, purpose, command path, runbook path, or validator path;
- new execution log;
- changed planning entry point.

Index updates are not required for timestamp-only edits or checklist progress inside an already
linked plan.

## Relationship To Structure Governance

Planning workflows own planning lifecycle semantics:

- document types;
- required planning headers;
- status values;
- revision notes;
- execution logs;
- plan bundles, subplans, plan families, and program folders;
- plan-scoped attachments.

Structure governance owns durable placement, naming, index-maintenance rules, and local wrapper
behavior for consumer repositories.

When a consumer keeps a local planning or governance document only for discoverability, that local
document should be a concise wrapper to the managed kit doctrine plus any real local exceptions.
Do not duplicate the managed planning lifecycle rules in a consumer-local wrapper.
