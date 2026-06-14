from __future__ import annotations

import json
from datetime import timedelta
from pathlib import Path
from urllib.error import URLError
from urllib.request import Request, urlopen

from .. import __version__
from ..lockfile import format_time, read_lock, utc_now, write_lock


DEFAULT_LATEST_RELEASE_URL = (
    "https://api.github.com/repos/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/latest"
)


class LatestVersionLookupError(RuntimeError):
    pass


def _version_tuple(value: str) -> tuple[int, ...]:
    result = []
    for part in value.lstrip("v").split("."):
        try:
            result.append(int(part))
        except ValueError:
            result.append(0)
    return tuple(result)


def _latest_version_from_metadata(metadata_url: str) -> str:
    request = Request(
        metadata_url,
        headers={
            "Accept": "application/vnd.github+json, application/json",
            "User-Agent": f"codeheart-operating-kit/{__version__}",
        },
    )
    try:
        with urlopen(request, timeout=10) as response:
            payload = json.loads(response.read().decode("utf-8"))
    except (OSError, URLError, json.JSONDecodeError) as exc:
        raise LatestVersionLookupError(str(exc)) from exc

    if not isinstance(payload, dict):
        raise LatestVersionLookupError("latest-version metadata must be a JSON object")

    latest = payload.get("tag_name") or payload.get("latest_version") or payload.get("version") or payload.get("name")
    if not latest:
        raise LatestVersionLookupError("latest-version metadata did not contain tag_name, latest_version, version, or name")
    return str(latest)


def _ensure_lock_defaults(lock: dict[str, object], current: str) -> None:
    lock.setdefault("schema_version", 1)
    lock.setdefault("kit_version", current)
    lock.setdefault("selected_profile", "standard")
    lock.setdefault("selected_components", [])
    lock.setdefault("release", {"asset_url": "unknown", "checksum_sha256": "0" * 64})
    lock.setdefault("managed_paths", [])
    lock.setdefault("generated_surfaces", [])
    lock.setdefault(
        "cli_repair",
        {
            "installed_cli_path": "codeheart-operating-kit",
            "repair_source_url": "unknown",
            "repair_checksum_sha256": "0" * 64,
        },
    )


def update_check(
    root: Path,
    latest_version: str | None = None,
    now_text: str | None = None,
    metadata_url: str | None = None,
) -> dict[str, object]:
    root = root.expanduser()
    lock = read_lock(root)
    now = utc_now()
    if now_text:
        from ..lockfile import parse_time

        now = parse_time(now_text)
    current = str(lock.get("kit_version") or __version__)
    _ensure_lock_defaults(lock, current)

    try:
        latest = latest_version or _latest_version_from_metadata(metadata_url or DEFAULT_LATEST_RELEASE_URL)
    except LatestVersionLookupError as exc:
        update_check_state = dict(lock.get("update_check") or {})
        update_check_state.setdefault("last_update_check_at", format_time(now))
        update_check_state.setdefault("next_update_check_due", format_time(now))
        update_check_state.setdefault("latest_seen_version", current)
        update_check_state["update_status"] = "failed"
        lock["update_check"] = update_check_state
        write_lock(root, lock)
        return {
            "status": "failed",
            "latest_seen_version": update_check_state["latest_seen_version"],
            "next_update_check_due": update_check_state["next_update_check_due"],
            "error": str(exc),
        }

    status = "update-available" if _version_tuple(str(latest)) > _version_tuple(current) else "current"
    lock["update_check"] = {
        "last_update_check_at": format_time(now),
        "next_update_check_due": format_time(now + timedelta(days=7)),
        "latest_seen_version": str(latest),
        "update_status": status,
    }
    write_lock(root, lock)
    return {"status": status, "latest_seen_version": str(latest), "next_update_check_due": lock["update_check"]["next_update_check_due"]}


def run(args) -> int:
    result = update_check(Path(args.path), args.latest_version, args.now, args.metadata_url)
    if args.json:
        print(json.dumps(result, indent=2, sort_keys=True))
    elif args.agent_notification and result["status"] == "current":
        return 0
    elif result["status"] == "update-available":
        print(f"Operating Kit update available: {result['latest_seen_version']}. Apply it only if the user asks.")
    elif result["status"] == "failed":
        print(f"Operating Kit update check failed: {result['error']}")
    else:
        print("Operating Kit is current.")
    return 0
