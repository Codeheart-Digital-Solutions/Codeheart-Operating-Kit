from __future__ import annotations

import json
from pathlib import Path

from ..components import write_adoption_cleanup_report, write_default_state
from .inspect import inspect_folder


def initialize(root: Path, project_name: str, purpose: str | None = None, selected_folder: str | None = None) -> dict[str, object]:
    root = root.expanduser()
    root.mkdir(parents=True, exist_ok=True)
    inspection = inspect_folder(root)
    preexisting = [
        candidate
        for candidate in ["AGENTS.md", "docs/repo/README.md", "docs/agent-memory/README.md"]
        if (root / candidate).exists()
    ]
    state = write_default_state(root, project_name, purpose, selected_folder or str(root))
    report_path = None
    if inspection["mode"] == "existing-technical-project-adoption":
        if preexisting:
            report_path = write_adoption_cleanup_report(root, preexisting)
    return {"inspection": inspection, "state": state, "adoption_cleanup_report": str(report_path) if report_path else None}


def run(args) -> int:
    result = initialize(Path(args.path), args.project_name, args.purpose, args.selected_folder)
    if args.json:
        print(json.dumps(result, indent=2, sort_keys=True))
    else:
        print("Operating Kit initialized.")
        print(f"Mode: {result['inspection']['mode']}")
        if result["adoption_cleanup_report"]:
            print(f"Adoption cleanup report: {result['adoption_cleanup_report']}")
    return 0
