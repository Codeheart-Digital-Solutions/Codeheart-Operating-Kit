Last updated: 2026-06-14T01:00:37Z (UTC)

# Bootstrap Codeheart Operating Kit

Use this public bootstrap when Codeheart Operating Kit is not installed yet. It does not require
preinstalled Codeheart skills.

Pinned G1 release:

```text
Version: v0.1.0
Release URL: https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v0.1.0
```

## Install The CLI

macOS installs into a user-level Operating Kit folder:

```sh
curl -fsSLO https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.0/install.sh
bash install.sh
```

The default command installs the CLI under:

```text
$HOME/.codeheart/operating-kit/bin/codeheart-operating-kit
```

Windows installs into the current user's local application data folder:

```powershell
Invoke-WebRequest -Uri "https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.0/install.ps1" -OutFile install.ps1
.\install.ps1
```

The default command installs the CLI under:

```text
%LOCALAPPDATA%\Codeheart\OperatingKit\bin\codeheart-operating-kit.cmd
```

Both installers download the pinned release asset for the selected version and verify its SHA-256
checksum before installing or repairing the CLI. A checksum mismatch stops installation.

## Start Onboarding

After the CLI is installed, run:

```sh
codeheart-operating-kit onboard
```

The onboarding flow starts with language selection, guides Codex setup, inspects the selected
folder, shows the setup plan, and asks before writing files.

## Codex Chat Setup

In the chat window, use the menu in the lower-right area of the message box:

- Set Model to `GPT-5.5`.
- Set Thinking to `Extra High`.
- Set Speed to `Fast`.

Open Settings from the bottom-left of the left sidebar. The General tab should open by default.
At the top of the main settings screen, set Work Mode to Coding. Directly beneath that, turn on
Default permissions, Auto review, and Full access. In the chat box area, check the lower-left
control named Approve for me.

## Project Setup

The onboarding flow asks whether the setup is for private automation, company automation, or a
software product before it asks for a project name.

The project folder name appears in the Codex left sidebar and groups chat threads for this setup.
After purpose-specific naming, the default recommendation is:

```text
Documents > <Project-Name>
```

Different-folder examples:

- `Documents > Maria-Automation`
- `Documents > Finance-Team-Automation`
- `Desktop > Booking-App`
- `Documents > Existing-Project-Name`

## Native Capabilities

After CLI bootstrap, onboarding can check whether Codex has baseline support for documents,
spreadsheets, presentations, browser work, and PDFs. Unavailable capabilities are reported as
degraded state, not setup failure.

## Update Checks

The installed root `AGENTS.md` uses `.codeheart/kit.lock.yaml` for quiet weekly update checks.
When the installed kit is current, Codex stays silent. When an update is available, Codex mentions
it briefly and asks before applying changes.
