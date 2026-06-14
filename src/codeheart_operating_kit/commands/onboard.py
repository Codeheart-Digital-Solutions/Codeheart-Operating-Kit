from __future__ import annotations

import json
from pathlib import Path

from ..capabilities import refresh_native_capability_lock
from .init import initialize
from .inspect import inspect_folder


PURPOSES = {
    "private-automation": "Private automation: personal documents, household tasks, reminders, research, or similar private work.",
    "company-automation": "Company automation: office work, documents, spreadsheets, email, Microsoft 365 integration, internal processes, or similar company work.",
    "software-product": "Software product: building or maintaining an app, website, technical product, or existing code project.",
}


def onboarding_script(project_name: str, purpose: str, target: Path, mode: str) -> list[str]:
    recommended = f"Documents > {project_name}"
    return [
        "Choose setup language / Sprache waehlen / 选择设置语言:",
        "1. English",
        "2. Deutsch",
        "3. 中文",
        "",
        "Before we set up your project folder, please adjust Codex in this chat.",
        "Look at the message box on the right. In the lower-right area, open the menu for model, thinking, and speed.",
        "Set Model: GPT-5.5",
        "Set Thinking: Extra High",
        "Set Speed: Fast",
        "",
        "Now open Codex Settings from the bottom-left of the left sidebar.",
        "The General tab should open automatically.",
        "At the top of the main settings screen, set Work Mode to Coding.",
        "Directly beneath that, turn on Default permissions, Auto review, and Full access.",
        "In the chat box area, check the lower-left control named Approve for me.",
        "",
        "What are you setting this up for?",
        f"1. {PURPOSES['private-automation']}",
        f"2. {PURPOSES['company-automation']}",
        f"3. {PURPOSES['software-product']}",
        "",
        "Codex needs a project folder. The folder name appears in the Codex left sidebar and groups chat threads for this setup.",
        purpose_name_prompt(purpose),
        f"I recommend this setup folder: {recommended}",
        "What should I use?",
        f"1. Yes, use {recommended}",
        "2. Use a different folder",
        "Please tell me the folder name or location if you choose a different folder.",
        "Examples: Documents > Maria-Automation; Documents > Finance-Team-Automation; Desktop > Booking-App; Documents > Existing-Project-Name",
        "",
        f"I will check this folder now so I can prepare the setup plan: {target}",
        mode_message(mode),
        plan_preview(mode),
        "Should I continue with this setup?",
        "1. Yes, set it up",
        "2. No, stop here",
        "",
        "I can check whether Codex has the tools for documents, spreadsheets, presentations, browser work, and PDFs.",
        "Operating Kit setup includes quiet weekly update checking. If everything is current, Codex will not mention it. If an update is available, Codex will ask before applying anything.",
    ]


def purpose_name_prompt(purpose: str) -> str:
    if purpose == "private-automation":
        return "What should this personal automation project be called? Examples: Maria-Automation, Family-Planning, Home-Research"
    if purpose == "company-automation":
        return "What is the company, team, or office area this automation is for? Examples: Bluebird-Automation, Finance-Team-Automation, M365-Office-Automation"
    return "What is the product or project name? Examples: Client-Portal, Booking-App, Storefront-Relaunch"


def mode_message(mode: str) -> str:
    messages = {
        "new-folder-setup": "This folder is ready for a new setup.",
        "existing-folder-setup": "This folder already contains files. I can set up Operating Kit here without replacing your existing files.",
        "existing-technical-project-adoption": "This is an existing technical project. I will preserve existing project instructions and docs.",
        "existing-operating-kit-repair": "This folder already has Operating Kit. I will check managed kit files, routing, and update state.",
        "ambiguous-folder-stop": "I cannot tell whether this folder is safe to set up.",
    }
    return messages[mode]


def plan_preview(mode: str) -> str:
    if mode == "existing-technical-project-adoption":
        return "Adoption plan: add .codeheart/kit/, config, lock, managed AGENTS block, missing memory files, and .codeheart/reports/adoption-cleanup-report.md without deleting overlapping docs."
    if mode == "existing-operating-kit-repair":
        return "Repair plan: check managed Operating Kit files, agent routing, setup information, update-check information, and native Codex capability status."
    if mode == "ambiguous-folder-stop":
        return "Stop plan: choose a different folder or provide more context before writing files."
    return "Setup plan: add .codeheart/kit/, kit config, kit lock, local user notes, AGENTS.md, docs/repo/, and missing docs/agent-memory/ files. I will not delete existing files."


def run(args) -> int:
    target = Path(args.target).expanduser()
    inspection = inspect_folder(target)
    script = onboarding_script(args.project_name, args.purpose, target, str(inspection["mode"]))
    result = {"inspection": inspection, "written": False}
    if args.yes and inspection["mode"] != "ambiguous-folder-stop":
        result = initialize(target, args.project_name, args.purpose, str(target))
        result["native_capabilities"] = refresh_native_capability_lock(target, attempt_install=True)
        result["written"] = True
    if args.json:
        result["script"] = script
        print(json.dumps(result, indent=2, sort_keys=True))
    else:
        print("\n".join(script))
        if args.yes:
            print("Base Operating Kit setup is complete.")
    return 0
