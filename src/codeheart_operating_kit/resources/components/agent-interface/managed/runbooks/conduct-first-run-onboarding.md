Last updated: 2026-07-09T23:30:00Z (UTC)

# Conduct First-Run Onboarding

Use this runbook when guiding a user through `codeheart-operating-kit onboard`. This public
runbook owns the onboarding order and required stop points for installed Operating Kit content.

The agent must not infer setup purpose, project name, target folder, or write approval from local
context. Ask the user directly and keep the conversation visible in Codex chat.

## Procedure

1. Ask for setup language with the three visible choices: English, Deutsch, and Chinese. Continue
   in the chosen language in chat after the user answers.
2. Guide Codex chat setup before folder setup. Tell the user to use the menu in the lower-right
   area of the message box and set Model to `GPT-5.5`, Thinking to `Extra High`, and Speed to
   `Fast`.
3. Guide Codex Settings before folder setup. Tell the user to open Settings from the bottom-left
   of the left sidebar, stay on the General tab, set Work Mode to Coding at the top of the main
   settings screen, and turn on Default permissions, Auto review, and Full access directly beneath
   it. Then ask them to check the chat-box control named Approve for me.
4. Ask whether the user already knows the Codex project name or wants a suggestion.
5. Ask lightweight purpose/context only when the user wants naming help.
6. Explain why the project folder name matters in the Codex sidebar.
7. Ask for the project name.
8. Ask whether the user already knows the target folder or wants a simple location.
9. Recommend `Documents > <Project-Name>` when the user wants a suggestion.
10. Ask whether to use the recommended folder or a different folder.
11. Inspect the selected folder before writing.
12. Present exactly one setup-mode message.
13. Present the concrete setup, adoption, or repair plan.
14. Ask for write confirmation.
15. Execute the state-matched lifecycle command after approval.
16. Explain quiet weekly update checking.
17. Finish with base Operating Kit setup completion.

## Required User-Owned Decisions

Ask the user before deciding:

- setup language;
- Codex project name;
- target folder;
- setup writes.

Do not use non-interactive flags to fill missing user decisions. Do not use `--yes` unless the user
already supplied the target folder, supplied the project name, and approved writing files.

## Purpose And Context Wording

Purpose is optional context, not a required setup branch. Ask for it only when the user wants help
choosing a name or needs next-step guidance:

```text
What is this mainly for?

1. Personal automation
2. Company operations
3. Software or product development
4. I am not sure yet
```

Use the answer only to make naming help concrete. Do not use it to select a different profile;
`standard` remains the installed profile.

## Project Naming Flow

Before asking for the name, explain that Codex uses the project folder name in the left sidebar and
groups chat threads under that project. A recognizable name helps the user reopen the same setup
later.

Use neutral examples only:

- Personal automation: `Yourname-Automation`
- Company operations: `Companyname-Automation`
- Software or product development: `Productname-Development`
- Team operations: `Teamname-Operations`

Do not use real-looking person, family, company, customer, tenant, or product names as examples.

## Folder Handling

Default to `Documents > <Project-Name>` when the user asks for a simple recommendation. If the
user chooses a different folder, ask for the folder name or location immediately and show only
neutral examples:

- `Documents > Companyname-Automation`
- `Documents > Yourname-Automation`
- `Desktop > Productname-Development`
- `Documents > Existing-Project-Name`

Always inspect before writing:

- `new-folder-setup`: prepare a new setup plan.
- `existing-folder-setup`: preserve existing files and prepare additive setup.
- `existing-technical-project-adoption`: preserve existing instructions and docs, add managed kit
  files, scaffold only missing memory files, and create an adoption cleanup report.
- `existing-operating-kit-repair`: inspect managed files, routing, config, lockfile, and native
  capability status before repair.
- `ambiguous-folder-stop`: stop and ask for a different folder or more context.

## Agent Execution Route

Keep the visible user dialogue above separate from state classification and command execution:

- absent or adoptable: show the setup plan; after approval run `init`;
- compatible existing lock-v2 state: show the repair plan; after approval run `repair`;
- compatible lock v1: show `repair --dry-run` so bounded migration is visible, then repair after
  approval;
- active transaction: wait and run `check`; do not start concurrent setup;
- schema-invalid, recovery-required, stale-CLI, or unsupported-future state: stop onboarding and
  return the `check` blocker;
- version change: stop onboarding and use the separate upgrade route only after the user asks for
  that version change.

Do not re-enter `init` for an existing installation. Onboarding confirmation authorizes only the
state-matched init or repair operation shown in its plan.

## Setup Plan And Write Boundary

Before writing, show the concrete file groups that will be added or repaired. Do not write files
before confirmation. Do not delete existing files during onboarding. Do not create a Python virtual
environment during default onboarding. Do not mention GitHub during first-run onboarding.

## Native Capabilities And Updates

Base onboarding does not offer, install, or implicitly check optional native capabilities. Their
existing lock state remains unchanged unless the user separately requests capability work.

Explain that Operating Kit includes quiet weekly update checking. If everything is current, Codex
does not mention it. If an update is available, Codex asks before applying anything.
