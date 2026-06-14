Last updated: 2026-06-14T00:23:12Z (UTC)

# Baseline Capability Profile

The G1 `standard` profile expects status reporting for these external Codex capabilities:

- `documents`: Word and `.docx` creation, editing, review, and rendering. Applies to the
  `standard` profile.
- `spreadsheets`: Excel, CSV, workbook analysis, formula work, and modification. Applies to the
  `standard` profile.
- `presentations`: PowerPoint and slide deck creation, editing, rendering, and export. Applies to
  the `standard` profile.
- `browser`: UI inspection, screenshots, interaction, and frontend verification. Applies to the
  `standard` profile.
- `pdf`: PDF reading, extraction, generation, rendering checks, and review. Applies to the
  `standard` profile.

## Status Values

- `available`
- `installed`
- `missing`
- `install-attempted`
- `unavailable`
- `blocked`
- `unknown`

Record capability status, checked-at timestamp, profile applicability, and command result category
in the lockfile. Do not record secrets, auth tokens, local absolute paths, or raw command logs.

Release manifests record the native baseline checks expected by the release. Lockfiles record the
status observed in the consumer folder. Missing, unavailable, blocked, or unknown capabilities are
degraded states for reporting; they are not automatic Operating Kit setup failures.
