Last updated: 2026-06-13T22:55:57Z (UTC)

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
