Last updated: 2026-06-15T10:16:29Z (UTC)

# Onboarding Context Contract

First-run onboarding is user-decision-first and non-technical until folder inspection detects an
existing technical project.

## Goal

Onboarding produces a usable Operating Kit consumer folder while keeping the early conversation
plain-language. The agent collects only the context needed to select a project name, select a
folder, explain the setup, inspect before writing, and persist non-secret setup metadata outside
managed kit content.

The agent must not infer setup purpose, project name, target folder, or write approval from folder
names, user names, repository names, or surrounding files.

## First Prompt

The public recommended first prompt is:

```text
Set up Codeheart Operating Kit:

https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/latest/download/bootstrap.md
```

Do not put onboarding explanations, setup options, project naming advice, permissions guidance, or
folder strategy in the first prompt. `bootstrap.md` and this managed contract own those details.

## Agent Contract

- Read the public `bootstrap.md`.
- Install or repair `codeheart-operating-kit` before running kit commands.
- Verify release checksums before trusting downloaded release assets.
- Start onboarding after the CLI is installed.
- Ask the user for language before continuing with user-facing setup prompts.
- Ask the user before deciding the Codex project name.
- Ask the user before deciding the target folder.
- Ask the user before writing setup files.
- Show the setup or adoption plan before writing files.
- Do not infer private automation, company automation, or software-product purpose from context.
- Do not use non-interactive flags to fill missing user decisions.
- Do not use `--yes` unless the user already supplied every required setup decision and approved
  writing files.
- Keep current update-check results silent and mention only available updates.

`codeheart-operating-kit onboard` is an agent-guided script and setup-plan renderer. It is not a
terminal stdin prompt loop. The agent shows rendered prompts in Codex chat, collects user decisions
in chat, and reruns the command with explicit values only when applying setup.

## Ordered Context

1. Ask for setup language: English, Deutsch, or Chinese.
2. Guide Codex chat setup in the lower-right message-box menu.
3. Guide Codex Settings from the left sidebar.
4. Ask whether the user already knows the Codex project name or wants a suggestion.
5. Ask lightweight purpose/context only when the user wants naming help.
6. Explain that the project folder name appears in the Codex left sidebar and groups chat threads.
7. Ask for the project name.
8. Ask whether the user already knows the target folder or wants a simple location.
9. Recommend `Documents > <Project-Name>` when the user wants a suggestion.
10. Inspect the selected folder before writing.
11. Show the setup, adoption, repair, or stop plan.
12. Ask before writing files.
13. Ask whether to check native Codex capabilities.
14. Explain quiet weekly update checking.

## Codex Setup Copy

Use the chat setup prompt before folder setup:

```text
Before we set up your project folder, please adjust Codex in this chat.

Look at the message box on the right. In the lower-right area, open the menu for model, thinking,
and speed.

Set:
- Model: GPT-5.5
- Thinking: Extra High
- Speed: Fast

Tell me when this is done.
```

Use the settings prompt after chat setup:

```text
Now open Codex Settings.

Look at the left sidebar. At the bottom-left, click Settings.

The General tab should open automatically.

At the very top of the main settings screen, find Work Mode and select Coding.

Directly beneath that, find Permissions and turn on all three setup options:
- Default permissions
- Auto review
- Full access

Then return to this chat. In the chat box area, check the lower-left control named Approve for me.
Turn it on when it is not already selected.

Tell me when this is done.
```

Do not add extra cautionary copy to this user-facing step.

## Project Naming Flow

Ask whether the user already knows the project name:

```text
Do you already know what this Codex project should be called, or should I suggest a name?
```

When the user wants help, ask:

```text
What is this mainly for?

1. Personal automation
2. Company operations
3. Software or product development
4. I am not sure yet
```

Use this answer only to help with naming and next-step guidance. It is not a required setup branch
and must not control the installed `standard` profile. Do not ask for organization metadata as a
separate abstract onboarding question.

Before asking for the name, explain:

```text
Codex needs one project folder.

The folder name is important because Codex will show it in the left sidebar. Chat threads for this
work will be grouped under that project name, so choose a name you will recognize later.
```

Use neutral examples only:

- Personal automation: `Yourname-Automation`
- Company operations: `Companyname-Automation`
- Software or product development: `Productname-Development`
- Team operations: `Teamname-Operations`

Do not use real-looking person, family, company, customer, tenant, or product names as examples.

## Target Folder Flow

Ask whether the user already knows the folder:

```text
Do you already know where this project folder should be, or should I suggest a simple location?
```

When the user wants a suggestion, recommend:

```text
Documents > <Project-Name>
```

Then ask:

```text
What should I use?

1. Yes, use Documents > <Project-Name>
2. Use a different folder
```

When the user selects a different folder, ask:

```text
Please tell me the folder name or location.

Examples:
- Documents > Companyname-Automation
- Documents > Yourname-Automation
- Desktop > Productname-Development
- Documents > Existing-Project-Name
```

## Inspection Modes

- `new-folder-setup`
- `existing-folder-setup`
- `existing-technical-project-adoption`
- `existing-operating-kit-repair`
- `ambiguous-folder-stop`

Inspection is read-only. It may classify folder contents, detect existing `AGENTS.md`, `docs/`,
`.git/`, `.codeheart/`, package manifests, or build files, and prepare a setup plan. It must not
write files until the user has seen the plan and confirmed setup.

Existing technical project adoption must preserve existing docs and instructions, add managed kit
files, scaffold only missing memory files, and create
`.codeheart/reports/adoption-cleanup-report.md` for overlapping docs. It must not delete
overlapping docs during onboarding.

The adoption cleanup report is consumer-owned evidence. It lists docs, instructions, or memory
surfaces that may now overlap with managed Operating Kit doctrine. It does not move, delete, or
rewrite those files during onboarding.

## Setup Plan Preview

Show the setup, adoption, or repair plan before asking for write approval. The plan must describe
the exact file groups that will be added or repaired and must say that existing files will not be
deleted.

## Native Capabilities

After setup writes complete, ask whether to check baseline support for documents, spreadsheets,
presentations, browser work, and PDFs. Unavailable tools are recorded as degraded state, not setup
failure.

## Update Checks

After setup and capability checks, explain quiet weekly update checking:

```text
Operating Kit setup includes quiet weekly update checking.

Once a week, Codex checks whether an Operating Kit update is available.

If everything is current, Codex will not mention it.

If an update is available, Codex will ask before applying anything.
```

## Non-Interactive Onboarding

Non-interactive onboarding is allowed only for automation and repair contexts where all user-owned
decisions are already explicit.

`codeheart-operating-kit onboard --yes` must fail when either required setup decision is missing:

- target folder;
- project name.

`codeheart-operating-kit onboard` must not write setup files unless `--yes` is present.

The `--purpose` flag remains supported as backward-compatible metadata. It is not required for a
normal setup and must not select a different G1 profile.

## Configuration Compatibility

Existing configs may contain one of these purpose values:

```yaml
setup_purpose: private-automation
setup_purpose: company-automation
setup_purpose: software-product
```

The kit must continue to read those values. New interactive setups should store purpose only when
the user explicitly answered the context question. Missing purpose metadata must pass validation.

## Storage

- Shared non-secret setup context: `.codeheart/kit.config.yaml`
- Installed version and update-check state: `.codeheart/kit.lock.yaml`
- Local personal preferences: `.codeheart/user/preferences.yaml`

Do not write user-specific answers inside `.codeheart/kit/`. The managed kit directory is replaced
or repaired by sync operations, while config, lock, reports, memory state, and local preferences
belong to the consumer folder.
