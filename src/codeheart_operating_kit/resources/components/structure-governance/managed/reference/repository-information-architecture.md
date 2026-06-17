Last updated: 2026-06-17T06:46:46Z (UTC)

# Repository Information Architecture

Use this reference before introducing durable names, new folder boundaries, or externally visible
identifiers.

This is a decision framework, not a repository map.

## Core Rule

Choose the ownership boundary first, then organize by domain inside that boundary.

Do not hide repo-level concepts inside the first product, feature, or task that needs them. Do not
place product-owned source in a repo-level folder just because it is widely used.

## Decision Sequence

Before choosing a path or name, identify:

1. Durable domain entities: products, clients, environments, workflows, packages, templates,
   policies, releases, instances, registries, or evidence.
2. Cardinality: singular, repeatable, or expected to become repeatable.
3. Ownership boundary: repo-owned, product-owned, instance-owned, generated, archived, managed kit,
   consumer-owned, local-user, or external platform-owned.
4. Artifact kind: source, template, instantiated record, generated output, runbook, reference,
   registry state, evidence, plan, or archive.
5. Lifecycle: source of truth, derived output, temporary staging, historical context, or
   operational state.
6. Visibility: where a new contributor would reasonably look first.
7. External constraints: official platform naming limits and character rules.

## Reasoning Tests

- Repeatable-entity test: if this concept can multiply, represent it as a collection or scoped
  member.
- Future-sibling test: if a likely future sibling would make the current name ambiguous, use the
  more specific name now.
- New-contributor test: a reader should infer owner, lifecycle, and artifact kind from the path.
- First-use trap test: do not place a broader concept under the first product or feature that
  needs it.
- Lifecycle-mixing test: do not mix templates, instantiated records, generated output, tracking
  state, and operational instructions only because one workflow touches them.
- External-boundary test: if content is meant to live in a separate repository or platform, say
  that explicitly.

## Placement Guidance

- Repo-level conventions, local development rules, repo governance, cross-cutting references, and
  repo-maintenance plans belong in consumer repo docs.
- Product-owned implementation, product contracts, product plans, product references, and product
  runbooks belong under the owning product, source area, package, or module boundary.
- Product-wide documentation belongs under `products/<product-slug>/docs/`.
- Source-area docs belong under `products/<product-slug>/<source-area>/docs/`.
- Package-local docs belong under `products/<product-slug>/packages/<package>/docs/`.
- Concrete instance, environment, estate, or deployed-state records should be modeled as visible
  collections when they are repeatable.
- Cross-cutting repository automation should remain visible as repo automation.
- Historical migration maps and transitional inventories should be labeled as migration,
  previous-generation, or archive context.
- Generated output should be excluded, placed under an explicit generated-output location, or
  documented as derived.

## Abstraction Scope

When planning a domain-specific implementation, identify which mechanics may apply beyond that
domain.

Candidates for shared primitives:

- external platform rules;
- file formats;
- protocols;
- parsers;
- validation constraints;
- reusable transformations;
- deterministic shortening;
- date and time normalization.

Domain-owned semantics:

- product names;
- business vocabulary;
- lifecycle states specific to one workflow;
- ownership policies;
- allowed repository lists;
- product configuration schema.

Centralize only when doing so reduces likely drift without coupling domains that may evolve
independently.

## Naming Guidance

Prefer names that communicate domain, responsibility, and scope without requiring private planning
context.

- Use the shortest name that remains specific enough for likely future siblings.
- Scope generic nouns with the domain when a repository may later contain multiple meanings.
- Avoid overloaded names such as `audit`, `ops`, `reader`, `bootstrap`, `shared`, `common`, or
  `core` unless the surrounding path makes the meaning unambiguous.
- Name collections as collections and members as members.
- Distinguish templates from instantiated records, plans from runbooks, source from generated
  output, and logical identifiers from concrete external names.
- Follow local separator and casing conventions for the file type or platform.

## External Platform Names

External names are durable interfaces. Treat source-control repositories, workflow names, cloud
resources, package names, release assets, and registry names as public contracts.

For every durable external name:

- distinguish the logical token used in source from the concrete platform name;
- include enough domain and responsibility to be understandable in the external console or audit
  trail;
- check current official platform limits and allowed characters before freezing the name;
- leave room for environment, region, tenant, account, or generated uniqueness tokens;
- prefer a shorter descriptive name over prose-like names near platform limits;
- use a naming helper, validator, or platform-specific naming convention for exact budgets.

## Deterministic Uniqueness Suffixes

Use deterministic uniqueness suffixes only when a durable external name needs reproducible
uniqueness and the readable owner, domain, and purpose prefix is not enough.

Normal triggers:

- shared or global platform namespace;
- replacement-sensitive resource;
- too little readable discriminator budget;
- tooling must reproduce the same name from reviewed source inputs;
- suffix preimage inputs are stable authority data.

Do not use suffixes when parent-local uniqueness is safe and readable, readability is the primary
contract, or the suffix would replace owner/domain/purpose semantics.

The default pattern is a readable prefix plus a reviewed deterministic suffix. The owning product
or platform contract defines exact algorithm, suffix length, preimage fields, and tests.

## Review Expectations

Plans or code changes that introduce durable names or folders should state intended owner, domain,
lifecycle, cardinality, and visibility when not obvious.

Reviewers should challenge:

- names that depend on private context;
- folders that mix artifact kinds or lifecycles;
- product-local placement for central repo concepts;
- repo-local placement for product-owned implementation;
- ambiguous external-repository staging paths;
- narrowly scoped helpers that reimplement shared mechanics without explanation;
- generic helpers that absorb product or domain semantics;
- platform-visible names that omit domain or responsibility context;
- ad hoc abbreviations introduced only after a name hits a platform limit.
