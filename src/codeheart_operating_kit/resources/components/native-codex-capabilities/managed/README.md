Last updated: 2026-06-13T22:55:57Z (UTC)

# Native Codex Capabilities

This managed domain describes external native Codex capabilities expected by the Operating Kit
baseline. These capabilities are not Codeheart-owned workflow doctrine.

## Routes

- Baseline capability profile: `reference/baseline-capability-profile.md`

Bootstrap and `check` should verify or attempt installation when Codex exposes supported commands,
record status, and report degraded state without hard-failing G1 setup solely because a capability
is unavailable.
