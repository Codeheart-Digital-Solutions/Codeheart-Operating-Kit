Last updated: 2026-07-31T23:43:02Z (UTC)

# Semantic Plan Catalog Consumer Impact Record

This public-safe record classifies the implemented source surfaces. The final version and public
release evidence remain owned by `EP-07`.

| Impact | Affected surfaces | Validation | Adoption or action |
| --- | --- | --- | --- |
| `instruction-only change` | Managed planning and Agent Interface references/runbooks, root routes, templates, and indexes | Routing assertions, fresh low-context probes, Markdown validation, source/package parity | Adopt through the released Kit; use the managed authoring, setup, refresh, and migration routes. |
| `validator-only change` | Plan/config/catalog/ledger/source/overlay schemas and Go parsing, identity, relation, compatibility, and diagnostic validation | JSON-schema fixtures, Go unit/integration tests, stable-code assertions | No action unless adopting semantic catalog modes. |
| `backwards-compatible scaffold addition` | Fresh plan-register entry point plus repo-owned portfolio README and strategic overlay | Init/onboard/sync fixtures; byte-preservation, symlink, non-regular, and pristine-placeholder tests | Fresh installs create missing files; existing repository-owned bytes remain protected. |
| `consumer migration required` | Repository plan metadata, tracked shared identity/config, frozen register baseline, mixed/canonical cutover | Inventory/ledger schema, hashes, dry-run/apply, branch-owner reconciliation, remote-overlay gate, idempotency | Upgrade first, inventory current refs, review and activate a repository-owned migration, remain mixed until complete. |
| `breaking placement-contract change` (not triggered) and generated/local-path effect | Formal plans and strategic overlay under `docs/repo/`; rebuildable catalog, mirrors, and local sources under `.codeheart/local/` | Placement-contract review, profile/component declarations, init/sync/check tests | Placement is additive; no existing authority moves and no generated catalog is committed. |
| `security or safety policy change` | Plan-only activation publication; inert branch reads; Git environment/protocol/hook policy; atomic cache/mirror/transaction operations; public-safe evidence | Authority negative tests, adversarial fixtures, race/containment tests, public-core validation, secret-capture assertions | Do not broaden activation into code/PR/merge/release authority; disclose incomplete scans and preserve prior complete evidence. |

Release notes: required. Migration/adoption notes: required. Known consumer action: none merely to
upgrade; each repository opts into semantic migration under its own authority.
