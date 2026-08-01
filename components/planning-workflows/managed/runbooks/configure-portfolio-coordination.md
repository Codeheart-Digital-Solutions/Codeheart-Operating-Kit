Last updated: 2026-07-31T23:03:05Z (UTC)

# Configure Portfolio Coordination

Audience: hybrid

Intent:
Guide a user one decision at a time, preview the exact member or coordination-home configuration,
and write only the reviewed non-secret scope and identity, missing home scaffolds, and the approved
home ID binding for an exact untouched overlay placeholder.

Success:
The repository has valid matching portfolio-v2 configuration; a coordination home also has its
approved committed and machine-local sources, preserved repo-owned scaffolds, and a disclosed first
scan result.

Agent judgment boundary:
Explain and normalize valid values, inspect existing state, and assemble the exact command. Do not
invent repository IDs, home IDs, organizations, local roots, roles, or replacement authority.

Stop boundary:
Stop on ambiguous repository authority, conflicting existing identity, private/sensitive source
material, missing write approval, unsafe overlapping changes, failed authentication, or an
incomplete first scan that the user expects to treat as current.

## Source, Inputs, And Lane

Source of truth, in order:

1. current user choices and repository instructions;
2. `.codeheart/kit.config.yaml` and `.codeheart/kit.lock.yaml`;
3. `../reference/portfolio-coordination-format.md`;
4. `codeheart-operating-kit portfolio configure --help`;
5. live `gh auth status` only when a GitHub-backed first scan is requested.

Inputs:

- target repository path;
- role: `member` or `coordination-home`;
- this repository's lowercase portable `member_repository_id`;
- lowercase portable `coordination_home_id`;
- for a home only, zero or more GitHub owner names and local Git roots; and
- explicit approval of the previewed write.

Execution lane: managed hybrid runbook plus the Operating Kit CLI. Missing local CLI or Git/`gh`
tooling routes to the managed tooling-readiness runbook; do not improvise installation.

Recipe maturity: L1 hybrid orchestration over the L3 `portfolio configure`/`portfolio scan`
commands. Fresh-agent review proves the conversation and authority flow; command and fixture tests
prove mechanics and structured blockers. Do not add an L2 script that merely wraps these commands.

Preconditions:

- the target is an initialized Kit repository;
- role and stable IDs are repository-owner decisions;
- repeated GitHub owners contain no credentials and are safe to commit;
- repeated local roots are machine-local and will remain below `.codeheart/local/`; and
- existing config and repo-owned portfolio files are preserved unless the exact reviewed action
  changes them.

Approval boundary: preview and inspection are local reads. `--yes` is a local repository write and
requires the user's approval of the displayed change plan. A first GitHub scan is an authenticated
read and requires the relevant live preflight, but no additional prompt when it is already part of
the user's approved setup request and local policy permits it.

## User-Facing Flow

Ask one decision per turn by default. Reuse a readable local language preference; otherwise ask
for language once before setup decisions.

1. Ask: “Should this repository be a normal portfolio member, or the coordination home that builds
   the cross-repository overview?” Explain that a home is also included as one member automatically.
2. Ask: “What stable lowercase repository ID should this repository use?” Offer a public-safe
   repository slug as a suggestion only; explain that renames should not change the ID.
3. Ask: “What stable lowercase coordination-home ID should this portfolio use?” For a member, this
   must be supplied by the portfolio owner. For a home, it identifies the portfolio rather than a
   filesystem path.
4. For a coordination home, ask: “Which GitHub owner scope may this home inspect?” Record zero or
   more explicit owner names. Do not infer an organization from remotes.
5. For a coordination home, ask separately: “Should this machine add a local folder containing Git
   repositories as a discovery source?” Explain that absolute paths remain ignored local state.
6. Present the exact command, files affected, committed-versus-local distinction, which portfolio
   files will be created because they are absent, and whether the byte-exact untouched overlay
   placeholder will be bound to the approved home ID. State that no other existing scaffold bytes
   will change.
7. Ask: “Do you approve this configuration write?” Do not combine this with unrelated upgrade,
   migration, commit, push, or release approval.

## Operator Notes

- `member_repository_id` is required for both roles. For a coordination home it is the home's own
  repository identity, not another member's ID.
- Members cannot carry discovery flags.
- Matching role and identity are idempotent. A conflict produces `portfolio_identity_conflict` and
  must be resolved by the repository owner; do not add a multiple-role field or overwrite it.
- GitHub owners are merged, deduplicated, and committed. Local roots are normalized, deduplicated,
  and written only to `.codeheart/local/portfolio/sources.yaml`.
- The existing README is preserved. Edited, non-regular, and symlink overlays are also preserved.
  Configuration replaces only the byte-exact untouched overlay placeholder with the approved home
  ID; any other existing overlay bytes remain repository-owned.

## Execution Path

1. Inspect the repository and run CLI help if command availability is uncertain.
2. Assemble one preview command.

   Member:

   ```sh
   codeheart-operating-kit portfolio configure \
     --role member \
     --member-repository-id <repository-id> \
     --coordination-home-id <home-id> \
     --dry-run <repository-path>
   ```

   Coordination home (repeat source flags as approved):

   ```sh
   codeheart-operating-kit portfolio configure \
     --role coordination-home \
     --member-repository-id <home-repository-id> \
     --coordination-home-id <home-id> \
     --github-owner <owner> \
     --local-root <root> \
     --dry-run <repository-path>
   ```

3. Show the deterministic action plan and blockers. Confirm only intended paths are included.
4. After approval, repeat the exact command with `--yes` instead of `--dry-run`.
5. Run the same `--dry-run` command again; matching configuration should produce no unintended
   change.
6. Run `codeheart-operating-kit plans validate <repository-path>`.
7. For a coordination home, preflight the approved sources and run
   `codeheart-operating-kit portfolio scan --format text <repository-path>`.

## Stop Conditions

Stop without writing when:

- the role, repository ID, home ID, target path, or owner scope is ambiguous;
- existing role or identity conflicts;
- a member request includes discovery sources;
- the config is malformed or the Kit installation is not initialized;
- the preview includes unrelated or overlapping dirty files;
- a local root would be committed or a credential would enter config;
- `--yes` approval is missing; or
- repository instructions require a narrower authority.

After a scan failure, retain the valid configuration but do not claim a current complete catalog.
Use the refresh runbook for disclosure and retry.

## Evidence And Validation

Record only non-secret evidence:

- target repository and role;
- stable repository and home IDs;
- approved GitHub owner names and count of local roots, without unnecessary absolute-path copying;
- preview result and applied action paths;
- idempotent re-preview result;
- plan-validation result; and
- first-scan completeness, timestamps, metrics, cache-updated/preserved state, and error codes.

Validation succeeds when config reloads, matching reconfiguration is idempotent, member/home
boundaries are correct, local sources remain ignored, README and repository-owned overlay bytes
are preserved, the byte-exact untouched overlay placeholder alone is bound to the approved home
ID, and a home scan is honestly reported as complete or incomplete.

## Recovery

If apply fails, preserve the CLI transaction evidence and follow its `check`/repair guidance; do
not hand-delete transaction files. If identity conflicts, leave existing bytes unchanged and ask
the repository owner whether a separately reviewed migration is intended. If a first scan is
incomplete, fix only the reported auth, scope, evidence, or access blocker and retry through
`refresh-portfolio-catalog.md`.
