Last updated: 2026-07-31T22:23:19Z (UTC)

# Portfolio Coordination

This repository can participate in portfolio coordination. Its configured role and stable
identity, when present, are authoritative in `.codeheart/kit.config.yaml`; this generic scaffold
does not assign a role. In a coordination home, source-derived plan facts are rebuilt into ignored
local state and strategic interpretation remains repository-owned here.

## Current Facts

Refresh before current analysis:

```sh
codeheart-operating-kit portfolio scan --format text .
```

The last complete cache lives at `.codeheart/local/portfolio/catalog.json`. It excludes worktree
changes, local heads, and unpushed commits. An incomplete refresh preserves the previous complete
cache but does not make it current.

## Strategic Interpretation

Keep public-safe cross-repository families, themes, relations, priorities, and analyses in
`strategic-overlay.yaml`. The scanner validates but never rewrites that file or member plans.

Managed format and operating guidance:

- `.codeheart/kit/docs/planning-workflows/reference/portfolio-coordination-format.md`
- `.codeheart/kit/docs/planning-workflows/runbooks/refresh-portfolio-catalog.md`

Add repository-specific coordination context below without copying credentials, private machine
paths, customer data, or a generated repository registry.
