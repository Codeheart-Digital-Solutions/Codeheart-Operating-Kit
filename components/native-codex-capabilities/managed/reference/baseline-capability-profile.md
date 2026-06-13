Last updated: 2026-06-13T23:16:03Z (UTC)

# Baseline Capability Profile

The G1 `standard` profile expects status reporting for these external Codex capabilities:

- `documents`: Word and `.docx` creation, editing, review, and rendering.
- `spreadsheets`: Excel, CSV, workbook analysis, formula work, and modification.
- `presentations`: PowerPoint and slide deck creation, editing, rendering, and export.
- `browser`: UI inspection, screenshots, interaction, and frontend verification.
- `pdf`: PDF reading, extraction, generation, rendering checks, and review.

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
