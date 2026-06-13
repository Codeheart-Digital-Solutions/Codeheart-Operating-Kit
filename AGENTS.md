Last updated: 2026-06-13T22:47:44Z (UTC)

# Repository Guidelines

## Purpose And Read Order

`AGENTS.md` is the first agent-facing bootstrap contract for Codeheart Operating Kit maintainers.
Keep it short, public-safe, and routing-focused.

Use task-matched documentation reading:

- Read `README.md` when the task involves repository purpose, public scope, or unsure placement.
- Read `docs/README.md` before creating, moving, or reorganizing documentation.
- Read `docs/repo/reference/placement-contract.md` before adding a new docs area, generated
  surface, component target, or externally visible path.
- Read `docs/repo/reference/consumer-impact-classification.md` before changing managed content,
  scaffolds, templates, schemas, validators, installers, release assets, or sync/check behavior.
- Read `docs/repo/runbooks/change-operating-kit.md` before changing kit source, docs, schemas,
  templates, validators, installers, or CLI behavior.
- Read `docs/repo/runbooks/release-operating-kit.md` before tagging, packaging, publishing, or
  updating release notes.
- Read `docs/repo/runbooks/promote-consumer-change.md` before moving consumer-local guidance into
  this kit.

Do not read every document by default. Read the smallest matching route and then follow links when
the task requires deeper context.

## Public-Core Safety

This repository is public. Do not add private Codeheart business records, customer names, tenant
details, credentials, secrets, local machine state, account identifiers, restricted strategy, or raw
operational logs.

Before committing copied or extracted material, classify it as reusable public doctrine,
public-safe placeholder content, plan-scoped evidence, or excluded private material. Sanitize
source-specific names and paths unless the public identifier is intentionally part of the kit.

## Source Of Truth

- `README.md` owns repository purpose and public boundary.
- `docs/repo/reference/placement-contract.md` owns placement rules for kit docs and generated
  consumer surfaces.
- `docs/repo/reference/consumer-impact-classification.md` owns impact classes for kit changes.
- Runbooks own ordered maintainer procedures.
- Component manifests, schemas, and release manifests own machine-readable behavior after they
  exist.

When instructions conflict, preserve public-core safety, avoid destructive or public-state-changing
actions, and record the inconsistency before proceeding.

## Change Safety

- Protect existing work. Do not overwrite, delete, or revert unrelated changes.
- Do not run destructive cleanup unless explicitly requested.
- Do not publish a release, push a tag, or change public GitHub settings unless the task explicitly
  asks for release or repository governance work.
- Do not print secrets or tokens. Mask sensitive values in logs.
- Keep docs and implementation changes in the same PR unless the design is still uncertain or the
  change is high-risk.

## Validation

Run the smallest validation that proves the changed surface. For repository-wide changes, run the
available markdown, schema, public-core, CLI, and release validation commands once they exist. If a
validator does not exist yet, record the residual risk in the PR or execution log.
