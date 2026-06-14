from __future__ import annotations

import json
from pathlib import Path


TECHNICAL_MARKERS = {
    ".git",
    "pyproject.toml",
    "package.json",
    "Cargo.toml",
    "go.mod",
    "Makefile",
    "AGENTS.md",
    "src",
    "tests",
}


def inspect_folder(path: Path) -> dict[str, object]:
    path = path.expanduser()
    if path.exists() and not path.is_dir():
        return {"path": str(path), "mode": "ambiguous-folder-stop", "reason": "target is not a folder"}
    if not path.exists():
        return {"path": str(path), "mode": "new-folder-setup", "reason": "folder does not exist"}
    entries = {entry.name for entry in path.iterdir()}
    if ".codeheart" in entries and ((path / ".codeheart/kit.lock.yaml").exists() or (path / ".codeheart/kit").exists()):
        return {"path": str(path), "mode": "existing-operating-kit-repair", "reason": "Operating Kit metadata exists"}
    if not entries:
        return {"path": str(path), "mode": "new-folder-setup", "reason": "folder is empty"}
    if entries & TECHNICAL_MARKERS:
        return {
            "path": str(path),
            "mode": "existing-technical-project-adoption",
            "reason": "technical project markers found",
            "markers": sorted(entries & TECHNICAL_MARKERS),
        }
    if len(entries) > 100:
        return {"path": str(path), "mode": "ambiguous-folder-stop", "reason": "folder has many existing files"}
    return {"path": str(path), "mode": "existing-folder-setup", "reason": "folder contains existing files"}


def run(args) -> int:
    result = inspect_folder(Path(args.path))
    if args.json:
        print(json.dumps(result, indent=2, sort_keys=True))
    else:
        print(f"{result['mode']}: {result['reason']}")
    return 0
