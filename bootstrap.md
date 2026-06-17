Last updated: 2026-06-16T21:22:03Z (UTC)

# Bootstrap Codeheart Operating Kit

Use this public bootstrap when Codeheart Operating Kit is not installed yet. It does not require
preinstalled Codeheart skills.

Pinned corrective release:

```text
Version: v0.1.3
Release URL: https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v0.1.3
```

## Install The CLI

macOS installs into a user-level Operating Kit folder:

```sh
curl -fsSLO https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.3/install.sh
bash install.sh
```

The default command installs the CLI under:

```text
$HOME/.codeheart/operating-kit/bin/codeheart-operating-kit
```

Windows installs into the current user's local application data folder:

```powershell
Invoke-WebRequest -Uri "https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.3/install.ps1" -OutFile install.ps1
.\install.ps1
```

The default command installs the CLI under:

```text
%LOCALAPPDATA%\Codeheart\OperatingKit\bin\codeheart-operating-kit.cmd
```

Both installers download the pinned release asset for the selected version and verify its SHA-256
checksum before installing or repairing the CLI. A checksum mismatch stops installation.

## Agent Contract

Follow this contract exactly during first-run setup:

- Read this public bootstrap before taking setup action.
- Install or repair `codeheart-operating-kit` before running kit commands.
- Verify release checksums before trusting downloaded release assets.
- Start onboarding after the CLI is installed.
- Ask the user for language before continuing with user-facing setup prompts.
- Ask the user before deciding the Codex project name.
- Ask the user before deciding the target folder.
- Ask the user before writing setup files.
- Show the setup or adoption plan before writing files.
- Do not infer private automation, company automation, or software-product purpose from folder
  names, user names, repository names, or surrounding files.
- Do not use non-interactive flags to fill missing user decisions.
- Do not use `--yes` unless the user already supplied every required setup decision and approved
  writing files.
- Keep current update-check results silent and mention only available updates.

For `v0.1.3`, `codeheart-operating-kit onboard` is an agent-guided script and setup-plan renderer.
It is not a terminal stdin prompt loop. Show rendered prompts in Codex chat, collect user decisions
in chat, and rerun the command with explicit values only when applying setup.

## Start Onboarding After CLI Install

After the CLI is installed, run:

```sh
codeheart-operating-kit onboard
```

The first run must show the language prompt before any setup-choice prompt:

```text
Choose setup language / Sprache waehlen / 选择设置语言:
1. English
2. Deutsch
3. 中文
```

After the user selects a language, continue in that language in chat. English is the source copy for
this release.

## Codex Chat Setup

Show this before folder setup:

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

Then show this settings step:

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

Do not add extra cautionary copy to this user-facing step. The Operating Kit safety rules still
apply, and setup files are not written before user approval.

## Project Name

Ask whether the user already knows the project name:

```text
Do you already know what this Codex project should be called, or should I suggest a name?
```

When the user already knows the name, ask for the name and continue to folder selection.

When the user wants help, ask this lightweight context question:

```text
What is this mainly for?

1. Personal automation
2. Company operations
3. Software or product development
4. I am not sure yet
```

Use the answer only to help with naming and next-step guidance. It is not a required setup branch
and must not control the installed `standard` profile.

Then explain why the name matters:

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

## Target Folder

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

## Folder Inspection

Before running inspection, show:

```text
Thank you. I will check this folder now so I can prepare the setup plan:

<Selected-Folder>

This check only looks at the folder. I will show you the setup plan before changing files.
```

Then run read-only inspection:

```sh
codeheart-operating-kit inspect <Selected-Folder>
```

Use exactly one setup mode message after inspection.

### `new-folder-setup`

```text
This folder is ready for a new setup.

I can add Codex working instructions, a small memory area, and quiet weekly update checking.
```

### `existing-folder-setup`

```text
This folder already contains files.

I can set up Operating Kit here without replacing your existing files. I will show the exact
additions before changing anything.
```

### `existing-technical-project-adoption`

```text
This is an existing technical project.

I will not overwrite existing docs or instructions. I will add the managed Operating Kit area,
preserve local instructions, scaffold only missing memory files, and create an adoption cleanup
report for overlapping docs.
```

### `existing-operating-kit-repair`

```text
This folder already has Operating Kit.

I will check whether the managed kit files, routing, and update state need repair. I will not apply
a version update unless you ask for it.
```

### `ambiguous-folder-stop`

```text
I cannot tell whether this folder is safe to set up.

Please use a different folder, or tell me more about what this folder is for.
```

Stop before writing files in `ambiguous-folder-stop`.

## Setup Plan Preview

For `new-folder-setup` and `existing-folder-setup`, show:

```text
Here is the setup plan.

I will add:
- Operating Kit files in .codeheart/kit/
- Setup information in .codeheart/kit.config.yaml
- Update-check information in .codeheart/kit.lock.yaml
- Local user notes in .codeheart/user/
- Agent instructions in AGENTS.md
- Repository notes in docs/repo/
- Agent memory files in docs/agent-memory/ when they are missing

I will not delete existing files.
```

For `existing-technical-project-adoption`, show:

```text
Here is the adoption plan.

I will add:
- Operating Kit files in .codeheart/kit/
- Setup information in .codeheart/kit.config.yaml
- Update-check information in .codeheart/kit.lock.yaml
- Local user notes in .codeheart/user/
- A managed Operating Kit block in AGENTS.md
- Missing agent memory files when needed
- An adoption cleanup report in .codeheart/reports/adoption-cleanup-report.md

I will preserve existing project instructions and docs.
I will not delete overlapping docs during onboarding.
```

For `existing-operating-kit-repair`, show:

```text
Here is the repair plan.

I will check:
- Managed Operating Kit files
- Agent instruction routing
- Setup information
- Update-check information
- Native Codex capability status

I will show repair findings before changing files.
```

## Write Confirmation

Before any setup write, ask:

```text
Should I continue with this setup?

1. Yes, set it up
2. No, stop here
```

If the user stops, do not write setup files.

## Native Capabilities

After setup writes complete, ask:

```text
I can check whether Codex has the tools for documents, spreadsheets, presentations, browser work,
and PDFs.

Should I check and set up these tools now?

1. Yes, check these tools
2. No, skip this for now
```

Unavailable tools are recorded as degraded state, not treated as setup failure.

## Update Checks

After setup and capability checks, show:

```text
Operating Kit setup includes quiet weekly update checking.

Once a week, Codex checks whether an Operating Kit update is available.

If everything is current, Codex will not mention it.

If an update is available, Codex will ask before applying anything.
```

## Completion Message

Finish with:

```text
Base Operating Kit setup is complete.

Foundry module setup will become available later, after the first Foundry module is released.
```

## Non-Interactive Onboarding

Non-interactive onboarding is allowed only for automation and repair contexts where every
user-owned decision is already explicit. `codeheart-operating-kit onboard --yes` must fail when the
target folder or project name is missing. `codeheart-operating-kit onboard` must not write setup
files unless `--yes` is present.

The `--purpose` flag remains supported as backward-compatible metadata. It is not required for a
normal setup and must not select a different G1 profile.
