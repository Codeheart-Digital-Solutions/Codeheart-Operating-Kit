from __future__ import annotations

import shutil
from datetime import timedelta
from pathlib import Path
from typing import Any

from . import __version__
from .lockfile import format_time, utc_now, write_config, write_lock
from .manifest import iter_component_files, kit_root, load_profile, sha256_file


BEGIN_MARKER = "<!-- BEGIN CODEHEART OPERATING KIT MANAGED BLOCK -->"
END_MARKER = "<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->"


def copy_managed_files(root: Path, profile_id: str = "standard") -> list[dict[str, Any]]:
    records: list[dict[str, Any]] = []
    for entry in iter_component_files(profile_id):
        if entry.get("ownership") != "managed":
            continue
        source = kit_root() / entry["source"]
        target = root / entry["target"]
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(source, target)
        records.append(
            {
                "path": entry["target"],
                "component": entry["component"],
                "source": entry["source"],
                "ownership": "managed",
                "checksum_sha256": sha256_file(target),
            }
        )
    return records


def render_agents(root: Path) -> str:
    template = (kit_root() / "templates/agents/AGENTS.managed-block.md").read_text(encoding="utf-8")
    existing_path = root / "AGENTS.md"
    if not existing_path.exists():
        existing_path.write_text(template, encoding="utf-8")
        return "created"

    existing = existing_path.read_text(encoding="utf-8")
    if BEGIN_MARKER in existing and END_MARKER in existing:
        before, rest = existing.split(BEGIN_MARKER, 1)
        _old, after = rest.split(END_MARKER, 1)
        managed = template.split(BEGIN_MARKER, 1)[1].split(END_MARKER, 1)[0]
        existing_path.write_text(f"{before}{BEGIN_MARKER}{managed}{END_MARKER}{after}", encoding="utf-8")
        return "repaired-managed-block"

    existing_path.write_text(template + "\n\n# Existing Instructions Preserved\n\n" + existing, encoding="utf-8")
    return "added-managed-block"


def scaffold_consumer_files(root: Path) -> list[dict[str, str]]:
    created: list[dict[str, str]] = []
    for entry in iter_component_files("standard"):
        if entry.get("ownership") != "scaffold":
            continue
        source = kit_root() / entry["source"]
        target = root / entry["target"]
        if target.exists():
            continue
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(source, target)
        created.append({"path": entry["target"], "ownership": "scaffold"})

    user_dir = root / ".codeheart/user"
    user_dir.mkdir(parents=True, exist_ok=True)
    user_readme = user_dir / "README.md"
    if not user_readme.exists():
        shutil.copyfile(kit_root() / "templates/user-layer/README.md", user_readme)
        created.append({"path": ".codeheart/user/README.md", "ownership": "local-user"})
    examples = user_dir / "examples"
    examples.mkdir(exist_ok=True)
    example_pref = examples / "preferences.yaml"
    if not example_pref.exists():
        shutil.copyfile(kit_root() / "templates/user-layer/example.preferences.yaml", example_pref)
        created.append({"path": ".codeheart/user/examples/preferences.yaml", "ownership": "local-user"})
    return created


def ensure_gitignore(root: Path) -> bool:
    path = root / ".gitignore"
    marker = ".codeheart/user/preferences.yaml"
    text = path.read_text(encoding="utf-8") if path.exists() else ""
    lines = text.splitlines()
    additions = [
        "",
        "# Codeheart Operating Kit local user layer",
        marker,
        ".codeheart/user/*.local.yaml",
    ]
    if marker in lines:
        return False
    path.write_text("\n".join(lines + additions).lstrip("\n") + "\n", encoding="utf-8")
    return True


def write_default_state(root: Path, project_name: str, purpose: str, selected_folder: str) -> dict[str, Any]:
    now = utc_now()
    profile = load_profile("standard")
    managed_records = copy_managed_files(root)
    scaffold_records = scaffold_consumer_files(root)
    agents_status = render_agents(root)
    gitignore_changed = ensure_gitignore(root)

    generated_surfaces = [
        {"path": ".codeheart/kit/", "ownership": "managed"},
        {"path": ".codeheart/kit.lock.yaml", "ownership": "generated-surface"},
        {"path": ".codeheart/kit.config.yaml", "ownership": "generated-surface"},
        {"path": "AGENTS.md", "ownership": "template"},
    ] + scaffold_records

    lock = {
        "schema_version": 1,
        "kit_version": __version__,
        "selected_profile": "standard",
        "selected_components": profile["selected_components"],
        "release": {
            "asset_url": "local-source",
            "checksum_sha256": "0" * 64,
        },
        "managed_paths": managed_records,
        "generated_surfaces": generated_surfaces,
        "cli_repair": {
            "installed_cli_path": "codeheart-operating-kit",
            "repair_source_url": "local-source",
            "repair_checksum_sha256": "0" * 64,
        },
        "update_check": {
            "last_update_check_at": format_time(now),
            "next_update_check_due": format_time(now + timedelta(days=7)),
            "latest_seen_version": __version__,
            "update_status": "current",
        },
        "native_capabilities": {
            capability: {
                "status": "unknown",
                "checked_at": format_time(now),
                "profile_applicability": "standard",
                "command_result_category": "not-checked",
            }
            for capability in ["documents", "spreadsheets", "presentations", "browser", "pdf"]
        },
    }
    write_lock(root, lock)
    write_config(
        root,
        {
            "schema_version": 1,
            "selected_profile": "standard",
            "setup_purpose": purpose,
            "project_display_name": project_name,
            "selected_setup_folder": selected_folder,
            "local_consumer_layer": {
                "repo_docs_path": "docs/repo/",
                "agent_memory_path": "docs/agent-memory/",
                "user_layer_path": ".codeheart/user/",
            },
            "component_settings": {},
        },
    )
    return {
        "managed_paths": managed_records,
        "generated_surfaces": generated_surfaces,
        "agents_status": agents_status,
        "gitignore_changed": gitignore_changed,
    }


def write_adoption_cleanup_report(root: Path, findings: list[str]) -> Path:
    report = root / ".codeheart/reports/adoption-cleanup-report.md"
    report.parent.mkdir(parents=True, exist_ok=True)
    body = [
        "Last updated: 2026-06-13T00:00:00Z (UTC)",
        "",
        "# Adoption Cleanup Report",
        "",
        "Existing project guidance was preserved. Review these overlapping surfaces before cleanup:",
        "",
    ]
    body.extend(f"- {item}" for item in findings)
    report.write_text("\n".join(body) + "\n", encoding="utf-8")
    return report
