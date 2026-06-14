from __future__ import annotations

import json
import os
import re
import shutil
import sys
from pathlib import Path

from .. import __version__
from ..capabilities import BASELINE_CAPABILITIES
from ..drift import drift_report
from ..lockfile import missing_required_lock_metadata, read_lock
from ..components import BEGIN_MARKER, END_MARKER


ROUTE_TARGET_PATTERN = re.compile(r"\.codeheart/kit/[A-Za-z0-9._/\-]+\.md")


def missing_route_targets(root: Path, agents_text: str) -> list[str]:
    targets = sorted(set(ROUTE_TARGET_PATTERN.findall(agents_text)))
    return [target for target in targets if not (root / target).exists()]


def cli_available() -> bool:
    if os.environ.get("CODEHEART_OPERATING_KIT_CLI") == "1":
        return True
    if shutil.which("codeheart-operating-kit"):
        return True
    executable = Path(sys.argv[0])
    return executable.name == "codeheart-operating-kit" and executable.exists()


def check_repository(root: Path) -> dict[str, object]:
    root = root.expanduser()
    lock = read_lock(root)
    missing_lock_metadata = missing_required_lock_metadata(lock)
    missing_cli = not cli_available()
    agents_text = (root / "AGENTS.md").read_text(encoding="utf-8") if (root / "AGENTS.md").exists() else ""
    missing_route_target_list = missing_route_targets(root, agents_text)
    missing_routing = (
        BEGIN_MARKER not in agents_text
        or END_MARKER not in agents_text
        or bool(missing_route_target_list)
    )
    stale_cli = bool(lock.get("kit_version") and lock.get("kit_version") != __version__)
    drift = drift_report(root, lock)
    native = lock.get("native_capabilities") or {
        capability: {"status": "unknown", "profile_applicability": "standard"}
        for capability in BASELINE_CAPABILITIES
    }
    return {
        "ok": not missing_lock_metadata and not missing_routing and not drift and not missing_cli and not stale_cli,
        "missing_cli": missing_cli,
        "stale_cli": stale_cli,
        "missing_routing": missing_routing,
        "missing_route_targets": missing_route_target_list,
        "missing_lock_metadata": missing_lock_metadata,
        "drift": drift,
        "native_capabilities": native,
    }


def run(args) -> int:
    result = check_repository(Path(args.path))
    if args.json:
        print(json.dumps(result, indent=2, sort_keys=True))
    else:
        print("Operating Kit check")
        print(f"OK: {result['ok']}")
        print(f"Missing CLI: {result['missing_cli']}")
        print(f"Stale CLI: {result['stale_cli']}")
        print(f"Missing routing: {result['missing_routing']}")
        print(f"Missing route targets: {len(result['missing_route_targets'])}")
        print(f"Drift findings: {len(result['drift'])}")
    return 0 if result["ok"] else 1
