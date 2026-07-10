#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "src"))

from codeheart_operating_kit.capabilities import BASELINE_CAPABILITIES  # noqa: E402
from codeheart_operating_kit.manifest import load_yaml  # noqa: E402


STATUS_VALUES = {
    "available",
    "installed",
    "missing",
    "install-attempted",
    "unavailable",
    "blocked",
    "unknown",
}
SHA256 = re.compile(r"^[a-fA-F0-9]{64}$")


def load_manifest(path: Path) -> dict[str, Any]:
    if path.suffix == ".json":
        return json.loads(path.read_text(encoding="utf-8"))
    return load_yaml(path)


def validate_manifest(path: Path) -> list[str]:
    manifest = load_manifest(path)
    if "compatibility" in manifest:
        return validate_content_manifest(path, manifest)
    return validate_legacy_release_manifest(path, manifest)


def validate_content_manifest(path: Path, manifest: dict[str, Any]) -> list[str]:
    errors: list[str] = []
    for key in ["schema_version", "version", "compatibility", "components", "profiles", "consumer_impact"]:
        if key not in manifest:
            errors.append(f"{path}: missing required key {key}")
    for forbidden in ["assets", "released_at"]:
        if forbidden in manifest:
            errors.append(f"{path}: embedded content identity contains forbidden key {forbidden}")
    compatibility = manifest.get("compatibility", {})
    if set(compatibility.get("platforms", [])) != {"macos-universal", "windows-x64"}:
        errors.append(f"{path}: compatibility platforms are incomplete")
    if set(compatibility.get("commands", [])) != {"init", "repair", "sync", "update-check", "upgrade", "check"}:
        errors.append(f"{path}: compatibility commands are incomplete")
    for item in [*manifest.get("components", []), *manifest.get("profiles", [])]:
        if not SHA256.match(str(item.get("checksum_sha256", ""))):
            errors.append(f"{path}: {item.get('id', '<unnamed>')} has invalid checksum_sha256")
            continue
        source = ROOT / str(item.get("manifest_path", ""))
        if source.is_file():
            actual = hashlib.sha256(source.read_bytes()).hexdigest()
            if actual != item["checksum_sha256"]:
                errors.append(f"{path}: {item.get('id', '<unnamed>')} checksum does not match {item.get('manifest_path')}")
    for profile in manifest.get("profiles", []):
        if not SHA256.match(str(profile.get("graph_sha256", ""))):
            errors.append(f"{path}: profile {profile.get('id', '<unnamed>')} has invalid graph_sha256")
    return errors


def validate_legacy_release_manifest(path: Path, manifest: dict[str, Any]) -> list[str]:
    errors: list[str] = []
    for key in [
        "schema_version",
        "version",
        "assets",
        "components",
        "profiles",
        "native_baseline_capability_checks",
        "consumer_impact",
    ]:
        if key not in manifest:
            errors.append(f"{path}: missing required key {key}")

    for asset in manifest.get("assets", []):
        checksum = str(asset.get("sha256", ""))
        if not SHA256.match(checksum):
            errors.append(f"{path}: asset {asset.get('name', '<unnamed>')} has invalid sha256")
        if asset.get("platform") not in {"macos", "macos-universal", "windows", "windows-x64", "universal"}:
            errors.append(f"{path}: asset {asset.get('name', '<unnamed>')} has invalid platform")

    component_ids = {component.get("id") for component in manifest.get("components", [])}
    expected_components = {
        "planning-workflows",
        "agent-memory",
        "agent-interface",
        "structure-governance",
        "native-codex-capabilities",
        "validators",
    }
    if not expected_components.issubset(component_ids):
        missing = sorted(expected_components - component_ids)
        errors.append(f"{path}: missing components {', '.join(missing)}")

    capability_checks = {
        item.get("id"): set(item.get("status_values", []))
        for item in manifest.get("native_baseline_capability_checks", [])
    }
    for capability in BASELINE_CAPABILITIES:
        if capability not in capability_checks:
            errors.append(f"{path}: missing native capability check {capability}")
        elif capability_checks[capability] != STATUS_VALUES:
            errors.append(f"{path}: native capability check {capability} has wrong status values")
    return errors


def default_manifests() -> list[Path]:
    paths = [ROOT / "tests/fixtures/release-manifest.json"]
    release_manifest = ROOT / "manifest.yaml"
    if release_manifest.exists():
        paths.append(release_manifest)
    packaged_manifest = ROOT / "src/codeheart_operating_kit/resources/manifest.yaml"
    if packaged_manifest.exists():
        paths.append(packaged_manifest)
    return paths


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate Operating Kit release manifests.")
    parser.add_argument("manifests", nargs="*", type=Path)
    args = parser.parse_args()

    manifests = [path.resolve() for path in args.manifests] if args.manifests else default_manifests()
    errors: list[str] = []
    for manifest in manifests:
        errors.extend(validate_manifest(manifest))
    if errors:
        print("Release manifest validation failed.")
        for error in errors:
            print(f"- {error}")
        return 1
    print("OK: release manifests validate.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
