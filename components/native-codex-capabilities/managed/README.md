Last updated: 2026-06-14T00:23:12Z (UTC)

# Native Codex Capabilities

This managed domain describes external native Codex capabilities expected by the Operating Kit
baseline. These capabilities are not Codeheart-owned workflow doctrine.

## Routes

- Baseline capability profile: `reference/baseline-capability-profile.md`

Onboarding attempts installation only when Codex exposes a supported plugin install command and the
capability is available but not installed. `check` reports the lockfile status. Missing, blocked,
unavailable, and unknown states are degraded states, not G1 setup failures.
