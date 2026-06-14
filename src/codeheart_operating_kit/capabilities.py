from __future__ import annotations

import json
from datetime import datetime, timezone
from pathlib import Path
from typing import Callable

from .codex_cli import run_codex
from .lockfile import read_lock, write_lock


BASELINE_CAPABILITIES = ["documents", "spreadsheets", "presentations", "browser", "pdf"]
PLUGIN_SELECTORS = {
    "documents": "documents@openai-primary-runtime",
    "spreadsheets": "spreadsheets@openai-primary-runtime",
    "presentations": "presentations@openai-primary-runtime",
    "browser": "browser@openai-bundled",
    "pdf": "pdf@openai-primary-runtime",
}
CodexRunner = Callable[[list[str]], dict[str, object]]


def _now_text() -> str:
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")


def _blank_status(status: str, category: str, checked_at: str) -> dict[str, str]:
    return {
        "status": status,
        "checked_at": checked_at,
        "profile_applicability": "standard",
        "command_result_category": category,
    }


def parse_capability_status(output: str) -> dict[str, str]:
    statuses = {capability: "unknown" for capability in BASELINE_CAPABILITIES}
    try:
        payload = json.loads(output)
    except json.JSONDecodeError:
        payload = None

    if isinstance(payload, dict):
        entries = []
        for key in ("installed", "available"):
            value = payload.get(key)
            if isinstance(value, list):
                entries.extend(item for item in value if isinstance(item, dict))
        by_name: dict[str, dict[str, object]] = {}
        for entry in entries:
            name = str(entry.get("name") or entry.get("pluginId", "").split("@", 1)[0])
            if name in BASELINE_CAPABILITIES and name not in by_name:
                by_name[name] = entry
        for capability in BASELINE_CAPABILITIES:
            entry = by_name.get(capability)
            if not entry:
                statuses[capability] = "missing"
                continue
            install_policy = str(entry.get("installPolicy", "")).lower()
            if entry.get("installed") is True:
                statuses[capability] = "installed"
            elif install_policy and install_policy != "available":
                statuses[capability] = "blocked"
            elif entry.get("installed") is False:
                statuses[capability] = "available"
            else:
                statuses[capability] = "unknown"
        return statuses

    lower = output.lower()
    for capability in BASELINE_CAPABILITIES:
        if capability in lower:
            statuses[capability] = "available"
    return statuses


def _discovery_failure_status(result: dict[str, object], checked_at: str) -> str:
    category = str(result.get("category", "unknown"))
    if category == "missing-cli":
        return "unavailable"
    if category in {"timeout", "permission-denied"}:
        return "blocked"
    return "unknown"


def check_native_capabilities(codex_runner: CodexRunner = run_codex) -> dict[str, dict[str, str]]:
    result = codex_runner(["plugin", "list", "--available", "--json"])
    now = _now_text()
    parsed = parse_capability_status(str(result.get("stdout", ""))) if result["ok"] else {}
    statuses: dict[str, dict[str, str]] = {}
    for capability in BASELINE_CAPABILITIES:
        status = parsed.get(capability, "unknown" if result["ok"] else _discovery_failure_status(result, now))
        statuses[capability] = _blank_status(status, str(result["category"]), now)
    return statuses


def install_native_capability(capability: str, codex_runner: CodexRunner = run_codex) -> dict[str, object]:
    selector = PLUGIN_SELECTORS[capability]
    return codex_runner(["plugin", "add", selector, "--json"])


def check_and_install_native_capabilities(codex_runner: CodexRunner = run_codex) -> dict[str, dict[str, str]]:
    statuses = check_native_capabilities(codex_runner)
    now = _now_text()
    for capability, record in statuses.items():
        if record["status"] != "available":
            continue
        result = install_native_capability(capability, codex_runner)
        if result["ok"]:
            record["status"] = "install-attempted"
        elif result["category"] in {"timeout", "permission-denied"}:
            record["status"] = "blocked"
        else:
            record["status"] = "unavailable"
        record["checked_at"] = now
        record["command_result_category"] = str(result["category"])
    return statuses


def refresh_native_capability_lock(
    root: Path,
    *,
    attempt_install: bool = False,
    codex_runner: CodexRunner = run_codex,
) -> dict[str, dict[str, str]]:
    statuses = (
        check_and_install_native_capabilities(codex_runner)
        if attempt_install
        else check_native_capabilities(codex_runner)
    )
    lock = read_lock(root)
    lock["native_capabilities"] = statuses
    write_lock(root, lock)
    return statuses
