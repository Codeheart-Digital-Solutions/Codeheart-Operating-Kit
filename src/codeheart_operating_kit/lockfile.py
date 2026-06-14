from __future__ import annotations

from datetime import datetime, timezone
from pathlib import Path
from typing import Any

from .manifest import dump_yaml, load_yaml


LOCK_PATH = Path(".codeheart/kit.lock.yaml")
CONFIG_PATH = Path(".codeheart/kit.config.yaml")


def utc_now() -> datetime:
    return datetime.now(timezone.utc).replace(microsecond=0)


def parse_time(value: str) -> datetime:
    return datetime.fromisoformat(value.replace("Z", "+00:00"))


def format_time(value: datetime) -> str:
    return value.astimezone(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def read_lock(root: Path) -> dict[str, Any]:
    path = root / LOCK_PATH
    if not path.exists():
        return {}
    return load_yaml(path)


def write_lock(root: Path, lock: dict[str, Any]) -> None:
    path = root / LOCK_PATH
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(dump_yaml(lock), encoding="utf-8")


def write_config(root: Path, config: dict[str, Any]) -> None:
    path = root / CONFIG_PATH
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(dump_yaml(config), encoding="utf-8")


def required_lock_keys() -> set[str]:
    return {
        "schema_version",
        "kit_version",
        "selected_profile",
        "selected_components",
        "release",
        "managed_paths",
        "generated_surfaces",
        "cli_repair",
        "update_check",
    }


def required_lock_metadata_paths() -> list[str]:
    return [
        "schema_version",
        "kit_version",
        "selected_profile",
        "selected_components",
        "release",
        "release.asset_url",
        "release.checksum_sha256",
        "managed_paths",
        "generated_surfaces",
        "cli_repair",
        "cli_repair.installed_cli_path",
        "cli_repair.repair_source_url",
        "update_check",
        "update_check.last_update_check_at",
        "update_check.next_update_check_due",
        "update_check.latest_seen_version",
        "update_check.update_status",
    ]


def missing_required_lock_metadata(lock: dict[str, Any]) -> list[str]:
    missing: list[str] = []
    for path in required_lock_metadata_paths():
        current: Any = lock
        found = True
        for part in path.split("."):
            if not isinstance(current, dict) or part not in current:
                found = False
                break
            current = current[part]
        if not found:
            missing.append(path)
    return missing
