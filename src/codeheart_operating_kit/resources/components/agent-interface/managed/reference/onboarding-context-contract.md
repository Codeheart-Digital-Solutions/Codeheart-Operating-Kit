Last updated: 2026-06-13T23:02:10Z (UTC)

# Onboarding Context Contract

First-run onboarding is purpose-first and non-technical until folder inspection detects an existing
technical project.

## Goal

Onboarding produces a usable Operating Kit consumer folder while keeping the early conversation
plain-language. The agent should collect only the context needed to select a folder, explain the
setup, inspect before writing, and persist non-secret setup metadata outside managed kit content.

## Ordered Context

1. Ask for setup language: English, Deutsch, or Chinese.
2. Guide Codex chat setup in the lower-right message-box menu:
   - Model: `GPT-5.5`
   - Thinking: `Extra High`
   - Speed: `Fast`
3. Guide Codex Settings from the left sidebar:
   - General tab.
   - Work Mode: Coding.
   - Permissions: Default permissions, Auto review, Full access.
   - Chat-box control: Approve for me.
4. Ask setup purpose:
   - Private automation.
   - Company automation.
   - Software product.
5. Explain that the project folder name appears in the Codex left sidebar and groups chat threads.
6. Ask for a purpose-specific project name.
7. Recommend `Documents > <Project-Name>`.
8. Offer two folder choices: use recommended folder or use a different folder.
9. Inspect the selected folder before writing.
10. Show the setup, adoption, repair, or stop plan.
11. Ask before writing files.
12. Ask whether to check native Codex capabilities.
13. Explain quiet weekly update checking.

## Purpose Flow

Ask for one purpose before asking for a project name:

- Private automation: personal documents, household tasks, reminders, research, or similar private
  work.
- Company automation: office work, documents, spreadsheets, email, Microsoft 365 integration,
  internal processes, or similar company work.
- Software product: building or maintaining an app, website, technical product, or existing code
  project.

Use the purpose to make the naming prompt concrete. Do not ask for organization metadata as a
separate abstract onboarding question.

## Project Naming Flow

Before asking for the name, explain that Codex uses the project folder name in the left sidebar and
groups chat threads under that project. A recognizable name helps the user reopen the same setup
later.

Recommended example categories:

- Private automation: `Maria-Automation`, `Family-Planning`, `Home-Research`
- Company automation: `Bluebird-Automation`, `Finance-Team-Automation`, `M365-Office-Automation`
- Software product: `Client-Portal`, `Booking-App`, `Storefront-Relaunch`

After the name is selected, recommend `Documents > <Project-Name>`. Offer only:

1. Use the recommended folder.
2. Use a different folder.

When the user chooses a different folder, immediately ask for the folder name or location and show
plain examples such as `Documents > Finance-Team-Automation` or `Desktop > Booking-App`.

## Inspection Modes

- `new-folder-setup`
- `existing-folder-setup`
- `existing-technical-project-adoption`
- `existing-operating-kit-repair`
- `ambiguous-folder-stop`

Existing technical project adoption must preserve existing docs and instructions, add managed kit
files, scaffold only missing memory files, and create
`.codeheart/reports/adoption-cleanup-report.md` for overlapping docs. It must not delete
overlapping docs during onboarding.

Inspection is read-only. It may classify folder contents, detect existing `AGENTS.md`, `docs/`,
`.git/`, `.codeheart/`, package manifests, or build files, and prepare a setup plan. It must not
write files until the user has seen the plan and confirmed setup.

The adoption cleanup report is consumer-owned evidence. It lists docs, instructions, or memory
surfaces that may now overlap with managed Operating Kit doctrine. It does not move, delete, or
rewrite those files during onboarding.

## Storage

- Shared non-secret setup context: `.codeheart/kit.config.yaml`
- Installed version and update-check state: `.codeheart/kit.lock.yaml`
- Local personal preferences: `.codeheart/user/preferences.yaml`

Do not write user-specific answers inside `.codeheart/kit/`. The managed kit directory is replaced
or repaired by sync operations, while config, lock, reports, memory state, and local preferences
belong to the consumer folder.
