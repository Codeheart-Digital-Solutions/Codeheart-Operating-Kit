Last updated: 2026-06-13T23:08:01Z (UTC)

# Conduct First-Run Onboarding

Use this runbook when guiding a user through `codeheart-operating-kit onboard`. This public
runbook owns the G1 onboarding order and required wording constraints for installed Operating Kit
content.

## Procedure

1. Ask for setup language with the three visible choices: English, Deutsch, and Chinese. Continue
   in the chosen language when localized copy is available.
2. Guide Codex chat setup before folder setup. Tell the user to use the menu in the lower-right
   area of the message box and set Model to `GPT-5.5`, Thinking to `Extra High`, and Speed to
   `Fast`.
3. Guide Codex Settings before folder setup. Tell the user to open Settings from the bottom-left
   of the left sidebar, stay on the General tab, set Work Mode to Coding at the top of the main
   settings screen, and turn on Default permissions, Auto review, and Full access directly beneath
   it. Then ask them to check the chat-box control named Approve for me.
4. Ask whether the setup is for private automation, company automation, or a software product.
5. Explain why the project folder name matters in the Codex sidebar.
6. Ask for a purpose-specific project name.
7. Recommend `Documents > <Project-Name>`.
8. Ask whether to use the recommended folder or a different folder.
9. Inspect the selected folder before writing.
10. Present exactly one setup-mode message.
11. Present the concrete setup, adoption, or repair plan.
12. Ask for write confirmation.
13. Check or record native Codex capability status when the user agrees.
14. Explain quiet weekly update checking.
15. Finish with base Operating Kit setup completion.

## Setup Purpose Wording

Keep the first purpose prompt non-technical:

- Private automation covers personal documents, household tasks, reminders, research, and similar
  private work.
- Company automation covers office work, documents, spreadsheets, email, Microsoft 365
  integration, internal processes, and similar company work.
- Software product covers apps, websites, technical products, and existing code projects.

Do not ask for a company or organization name as a standalone metadata question. Ask for it only
when it helps produce a purpose-specific project name.

## Folder Handling

Default to `Documents > <Project-Name>` because it is understandable to non-technical users and is
easy to reopen in Codex later. If the user chooses a different folder, ask for the folder name or
location immediately and show simple examples.

Always inspect before writing:

- `new-folder-setup`: prepare a new setup plan.
- `existing-folder-setup`: preserve existing files and prepare additive setup.
- `existing-technical-project-adoption`: preserve existing instructions and docs, add managed kit
  files, scaffold only missing memory files, and create an adoption cleanup report.
- `existing-operating-kit-repair`: inspect managed files, routing, config, lockfile, and native
  capability status before repair.
- `ambiguous-folder-stop`: stop and ask for a different folder or more context.

## Write Boundary

Before writing, show the concrete file groups that will be added or repaired. Do not delete
existing files during onboarding. Do not create a Python virtual environment during default G1
onboarding. Do not mention GitHub during first-run onboarding.

## Safety

Do not write files before confirmation. Do not delete existing docs during onboarding. Record
unavailable native capabilities as degraded state, not setup failure.
