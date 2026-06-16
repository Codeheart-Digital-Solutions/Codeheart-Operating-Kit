from __future__ import annotations

import json
import sys
from datetime import timedelta
from pathlib import Path
from typing import Any

from .. import __version__
from ..components import copy_managed_files
from ..lockfile import format_time, read_lock, utc_now, write_lock
from ..manifest import kit_root, load_profile, load_yaml, validate_release_manifest


def _valid_sha256(value: object) -> bool:
    text = str(value or "")
    return len(text) == 64 and all(char in "0123456789abcdefABCDEF" for char in text)


def _current_platform() -> str:
    if sys.platform.startswith("win"):
        return "windows"
    if sys.platform == "darwin":
        return "macos"
    return "universal"


def _release_asset_from_manifest(manifest: dict[str, Any] | None) -> dict[str, str] | None:
    if not manifest:
        return None
    assets = manifest.get("assets") or []
    if not isinstance(assets, list):
        return None
    preferred_platform = _current_platform()

    def is_cli_asset(asset: dict[str, Any]) -> bool:
        name = str(asset.get("name", ""))
        return name.startswith(f"codeheart-operating-kit-{__version__}") and not name.endswith(".sha256")

    candidates = [
        asset
        for asset in assets
        if isinstance(asset, dict)
        and asset.get("platform") in {preferred_platform, "universal"}
        and is_cli_asset(asset)
    ]
    if not candidates:
        candidates = [
            asset
            for asset in assets
            if isinstance(asset, dict)
            and is_cli_asset(asset)
            and asset.get("platform") in {preferred_platform, "universal"}
            and str(asset.get("url", ""))
            and _valid_sha256(asset.get("sha256"))
        ]
    if not candidates:
        return None
    asset = candidates[0]
    return {
        "asset_url": str(asset["url"]),
        "checksum_sha256": str(asset["sha256"]),
    }


def _bundled_release_manifest() -> dict[str, Any] | None:
    manifest_path = kit_root() / "manifest.yaml"
    if not manifest_path.exists():
        return None
    return load_yaml(manifest_path)


def _release_metadata(existing: dict[str, Any], release_manifest: dict[str, Any] | None) -> dict[str, str]:
    from_manifest = _release_asset_from_manifest(release_manifest) or _release_asset_from_manifest(_bundled_release_manifest())
    if from_manifest:
        return from_manifest

    existing_release = existing.get("release") if isinstance(existing.get("release"), dict) else {}
    asset_url = str(existing_release.get("asset_url") or "local-source")
    checksum = str(existing_release.get("checksum_sha256") or "")
    if not _valid_sha256(checksum):
        checksum = "0" * 64
    return {"asset_url": asset_url, "checksum_sha256": checksum}


def _refresh_lock(root: Path, managed: list[dict[str, Any]], release_manifest: dict[str, Any] | None) -> dict[str, Any]:
    now = utc_now()
    lock = read_lock(root)
    profile_id = str(lock.get("selected_profile") or "standard")
    profile = load_profile(profile_id)

    existing_update = lock.get("update_check") if isinstance(lock.get("update_check"), dict) else {}
    existing_cli = lock.get("cli_repair") if isinstance(lock.get("cli_repair"), dict) else {}
    existing_native = lock.get("native_capabilities") if isinstance(lock.get("native_capabilities"), dict) else {}

    refreshed = {
        "schema_version": 1,
        "kit_version": __version__,
        "selected_profile": profile_id,
        "selected_components": profile["selected_components"],
        "release": _release_metadata(lock, release_manifest),
        "managed_paths": managed,
        "generated_surfaces": lock.get("generated_surfaces") or [],
        "cli_repair": {
            "installed_cli_path": existing_cli.get("installed_cli_path", "codeheart-operating-kit"),
            "repair_source_url": existing_cli.get("repair_source_url", "local-source"),
            "repair_checksum_sha256": existing_cli.get("repair_checksum_sha256", "0" * 64),
        },
        "update_check": {
            "last_update_check_at": existing_update.get("last_update_check_at", format_time(now)),
            "next_update_check_due": existing_update.get("next_update_check_due", format_time(now + timedelta(days=7))),
            "latest_seen_version": __version__,
            "update_status": "current",
        },
        "native_capabilities": existing_native,
    }
    if not _valid_sha256(refreshed["cli_repair"].get("repair_checksum_sha256")):
        refreshed["cli_repair"]["repair_checksum_sha256"] = "0" * 64
    write_lock(root, refreshed)
    return refreshed


def run(args) -> int:
    root = Path(args.path).expanduser()
    release = None
    if args.release_manifest:
        release = validate_release_manifest(Path(args.release_manifest))
    managed = copy_managed_files(root)
    refreshed_lock = _refresh_lock(root, managed, release)
    result = {
        "synced_managed_paths": managed,
        "release_manifest": release is not None,
        "kit_version": refreshed_lock["kit_version"],
    }
    if args.json:
        print(json.dumps(result, indent=2, sort_keys=True))
    else:
        print(f"Synced {len(managed)} managed files under .codeheart/kit/.")
    return 0
