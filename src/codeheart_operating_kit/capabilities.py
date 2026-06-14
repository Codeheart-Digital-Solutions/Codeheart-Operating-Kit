from __future__ import annotations

from datetime import datetime, timezone

from .codex_cli import run_codex


BASELINE_CAPABILITIES = ["documents", "spreadsheets", "presentations", "browser", "pdf"]


def parse_capability_status(output: str) -> dict[str, str]:
    statuses = {capability: "unknown" for capability in BASELINE_CAPABILITIES}
    lower = output.lower()
    for capability in BASELINE_CAPABILITIES:
        if capability in lower:
            statuses[capability] = "available"
    return statuses


def check_native_capabilities() -> dict[str, dict[str, str]]:
    result = run_codex(["plugins", "list"])
    now = datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")
    parsed = parse_capability_status(str(result.get("stdout", ""))) if result["ok"] else {}
    statuses: dict[str, dict[str, str]] = {}
    for capability in BASELINE_CAPABILITIES:
        statuses[capability] = {
            "status": parsed.get(capability, "unknown" if result["ok"] else "unavailable"),
            "checked_at": now,
            "profile_applicability": "standard",
            "command_result_category": str(result["category"]),
        }
    return statuses


def install_native_capability(capability: str) -> dict[str, object]:
    return run_codex(["plugins", "install", capability])
