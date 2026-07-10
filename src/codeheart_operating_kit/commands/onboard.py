from __future__ import annotations

import json
from pathlib import Path
from typing import Any

from ..capabilities import refresh_native_capability_lock
from .init import initialize
from .inspect import inspect_folder


PURPOSE_LABELS = {
    "private-automation": "Personal automation",
    "company-automation": "Company operations",
    "software-product": "Software or product development",
}


def required_user_decisions_missing(args) -> list[str]:
    missing: list[str] = []
    if not getattr(args, "target", None):
        missing.append("target_folder")
    if not getattr(args, "project_name", None):
        missing.append("project_name")
    return missing


def onboarding_script(project_name: str | None, purpose: str | None, target: Path | None, mode: str | None) -> list[str]:
    display_name = project_name or "<Project-Name>"
    selected_folder = str(target) if target else "<Selected-Folder>"
    recommended = f"Documents > {display_name}"
    lines = [
        "Choose setup language / Sprache waehlen / 选择设置语言:",
        "1. English",
        "2. Deutsch",
        "3. 中文",
        "",
        "Before we set up your project folder, please adjust Codex in this chat.",
        "Look at the message box on the right. In the lower-right area, open the menu for model, thinking, and speed.",
        "Set:",
        "- Model: GPT-5.5",
        "- Thinking: Extra High",
        "- Speed: Fast",
        "Tell me when this is done.",
        "",
        "Now open Codex Settings.",
        "Look at the left sidebar. At the bottom-left, click Settings.",
        "The General tab should open automatically.",
        "At the very top of the main settings screen, find Work Mode and select Coding.",
        "Directly beneath that, find Permissions and turn on all three setup options:",
        "- Default permissions",
        "- Auto review",
        "- Full access",
        "Then return to this chat. In the chat box area, check the lower-left control named Approve for me.",
        "Turn it on when it is not already selected.",
        "Tell me when this is done.",
        "",
        "Do you already know what this Codex project should be called, or should I suggest a name?",
        "",
        "If the user wants help, ask:",
        "What is this mainly for?",
        "1. Personal automation",
        "2. Company operations",
        "3. Software or product development",
        "4. I am not sure yet",
    ]
    if purpose:
        lines.extend([
            "",
            f"Explicit setup context metadata supplied: {purpose} ({PURPOSE_LABELS[purpose]}).",
            "Use this only for naming help and next-step guidance; it does not change the standard profile.",
        ])
    lines.extend([
        "",
        "Codex needs one project folder.",
        "The folder name is important because Codex will show it in the left sidebar. Chat threads for this work will be grouped under that project name, so choose a name you will recognize later.",
        f"Selected project name: {display_name}",
        "Neutral examples: Yourname-Automation; Companyname-Automation; Productname-Development; Teamname-Operations",
        "",
        "Do you already know where this project folder should be, or should I suggest a simple location?",
        f"Recommended folder: {recommended}",
        "What should I use?",
        f"1. Yes, use {recommended}",
        "2. Use a different folder",
        "Please tell me the folder name or location if you choose a different folder.",
        "Examples: Documents > Companyname-Automation; Documents > Yourname-Automation; Desktop > Productname-Development; Documents > Existing-Project-Name",
        "",
        f"Thank you. I will check this folder now so I can prepare the setup plan: {selected_folder}",
        "This check only looks at the folder. I will show you the setup plan before changing files.",
    ])
    if mode:
        lines.extend([mode_message(mode), plan_preview(mode)])
    else:
        lines.extend([
            "Folder inspection is waiting for the selected target folder.",
            "Do not write setup files until the user supplies the target folder, sees the setup plan, and approves setup.",
        ])
    lines.extend([
        "Should I continue with this setup?",
        "1. Yes, set it up",
        "2. No, stop here",
        "",
        "After setup writes complete, ask whether to check native Codex capabilities.",
        "I can check whether Codex has the tools for documents, spreadsheets, presentations, browser work, and PDFs.",
        "Should I check and set up these tools now?",
        "1. Yes, check these tools",
        "2. No, skip this for now",
        "",
        "After successful setup, finish with: Base Operating Kit setup is complete.",
        "Foundry module setup will become available later, after the first Foundry module is released.",
    ])
    return lines


def mode_message(mode: str) -> str:
    messages = {
        "new-folder-setup": "This folder is ready for a new setup. I can add Codex working instructions and a small memory area.",
        "existing-folder-setup": "This folder already contains files. I can set up Operating Kit here without replacing your existing files. I will show the exact additions before changing anything.",
        "existing-technical-project-adoption": "This is an existing technical project. I will not overwrite existing docs or instructions. I will add the managed Operating Kit area, preserve local instructions, scaffold only missing memory files, and create an adoption cleanup report for overlapping docs.",
        "existing-operating-kit-repair": "This folder already has Operating Kit. I will check whether the managed kit files, routing, and lifecycle state need repair. I will not apply a version update unless you ask for it.",
        "ambiguous-folder-stop": "I cannot tell whether this folder is safe to set up. Please use a different folder, or tell me more about what this folder is for.",
    }
    return messages[mode]


def plan_preview(mode: str) -> str:
    if mode == "existing-technical-project-adoption":
        return "Adoption plan: add .codeheart/kit/, config, lock, managed AGENTS block, missing memory files, and .codeheart/reports/adoption-cleanup-report.md without deleting overlapping docs."
    if mode == "existing-operating-kit-repair":
        return "Repair plan: check managed Operating Kit files, agent routing, setup information, lifecycle metadata, and native Codex capability status."
    if mode == "ambiguous-folder-stop":
        return "Stop plan: choose a different folder or provide more context before writing files."
    return "Setup plan: add .codeheart/kit/, kit config, kit lock, local user notes, AGENTS.md, docs/repo/, and missing docs/agent-memory/ files. I will not delete existing files."


def emit_result(result: dict[str, Any], script: list[str], json_output: bool) -> None:
    if json_output:
        payload = {**result, "script": script}
        print(json.dumps(payload, indent=2, sort_keys=True))
    else:
        print("\n".join(script))
        if result.get("error"):
            print(f"Error: {result['error']}")
        elif result.get("written"):
            print("Base Operating Kit setup is complete.")


def run(args) -> int:
    target = Path(args.target).expanduser() if args.target else None
    inspection = inspect_folder(target) if target else {
        "path": None,
        "mode": "target-folder-pending",
        "reason": "target folder must be chosen by the user before inspection",
    }
    mode = str(inspection["mode"]) if target else None
    missing_decisions = required_user_decisions_missing(args)
    script = onboarding_script(args.project_name, args.purpose, target, mode)
    result: dict[str, Any] = {
        "inspection": inspection,
        "written": False,
        "write_approved": bool(args.yes),
        "required_user_decisions_missing": missing_decisions,
    }

    if args.yes and missing_decisions:
        result["error"] = "Cannot write setup files because required user decisions are missing: " + ", ".join(missing_decisions)
        emit_result(result, script, args.json)
        return 2

    if args.yes and inspection["mode"] == "ambiguous-folder-stop":
        result["error"] = "Cannot write setup files because folder inspection is ambiguous."
        emit_result(result, script, args.json)
        return 1

    if args.yes:
        assert target is not None
        assert args.project_name is not None
        result = initialize(target, args.project_name, args.purpose, str(target))
        result["native_capabilities"] = refresh_native_capability_lock(target, attempt_install=True)
        result["written"] = True
        result["write_approved"] = True
        result["required_user_decisions_missing"] = []

    emit_result(result, script, args.json)
    return 0
