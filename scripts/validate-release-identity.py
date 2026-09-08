#!/usr/bin/env python3
"""Check source release identities and their packaged declaration mirrors (read-only)."""
from __future__ import annotations

import argparse
import hashlib
import re
import tomllib
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[1]


def validate_identity(root: Path) -> list[str]:
    errors: list[str] = []
    manifest = yaml.safe_load((root / "manifest.yaml").read_text(encoding="utf-8"))
    version = manifest["version"]
    project_version = tomllib.loads((root / "pyproject.toml").read_text(encoding="utf-8"))["project"]["version"]
    go_versions = re.findall(r'^var Version = "([^"]+)"$', (root / "internal/version/version.go").read_text(encoding="utf-8"), re.MULTILINE)
    if project_version != version:
        errors.append("pyproject.toml: project version differs from manifest.yaml")
    if go_versions != [version]:
        errors.append("internal/version/version.go: literal CLI version differs from manifest.yaml")
    for relative, pattern, label in [
        ("src/codeheart_operating_kit/__init__.py", r'^__version__ = "([^"]+)"$', "Python package version"),
        ("install.sh", r'^VERSION="([^"]+)"$', "installer default version"),
        ("install.ps1", r'^    \[string\]\$Version = "([^"]+)",$', "installer default version"),
        ("install.sh", r'^  --version VERSION          Release version to install\. Default: (\S+)$', "installer help version"),
        ("install.ps1", r'^  -Version VERSION       Release version to install\. Default: (\S+)$', "installer help version"),
    ]:
        values = re.findall(pattern, (root / relative).read_text(encoding="utf-8"), re.MULTILINE)
        if values != [version]:
            errors.append(f"{relative}: literal {label} differs from manifest.yaml")
    bootstrap = (root / "bootstrap.md").read_text(encoding="utf-8")
    release_url = re.escape("https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/")
    for pattern in [
        r"^Version: v(\S+)$",
        r"^Release URL: " + release_url + r"tag/v(\S+)$",
        r"^`codeheart-operating-kit-([^`]+)-macos-universal\.zip`; the pack contains$",
        r"^`codeheart-operating-kit-([^`]+)-windows-x64\.zip`; the pack contains$",
        r"^curl -fsSLO " + release_url + r"download/v([^/]+)/install\.sh$",
        r'^Invoke-WebRequest -Uri "' + release_url + r'download/v([^/]+)/install\.ps1" -OutFile install\.ps1$',
        r"^For `v([^`]+)`, `codeheart-operating-kit onboard` is an agent-guided script and setup-plan renderer\.$",
    ]:
        if re.findall(pattern, bootstrap, re.MULTILINE) != [version]:
            errors.append("bootstrap.md: known current-release reference differs or is missing")
    mirrored = ["manifest.yaml"]
    for collection, key, directory, pattern in [
        ("components", "component", "components", "*/component.yaml"),
        ("profiles", "profile", "profiles", "*.yaml"),
    ]:
        entries = manifest[collection]
        paths = [entry["manifest_path"] for entry in entries]
        expected = {path.relative_to(root).as_posix() for path in (root / directory).glob(pattern)}
        if len(paths) != len(set(paths)) or set(paths) != expected:
            errors.append(f"manifest.yaml: {collection} declarations do not match source inventory")
        for entry in entries:
            relative = entry["manifest_path"]
            if relative not in expected:
                errors.append(f"{relative}: missing or unexpected declaration path")
                continue
            data = (root / relative).read_bytes()
            declaration = yaml.safe_load(data)[key]
            for field in ("id", "version"):
                if entry[field] != declaration[field]:
                    errors.append(f"{relative}: {field} differs from manifest.yaml")
            if entry["checksum_sha256"] != hashlib.sha256(data).hexdigest():
                errors.append(f"{relative}: checksum differs from manifest.yaml")
            if key == "profile" and declaration["version"] != version:
                errors.append(f"{relative}: profile version differs from release version")
            mirrored.append(relative)
    for relative in mirrored:
        mirror = root / "src/codeheart_operating_kit/resources" / relative
        if not mirror.is_file() or mirror.read_bytes() != (root / relative).read_bytes():
            errors.append(f"{relative}: packaged declaration mirror differs or is missing")
    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=ROOT)
    args = parser.parse_args()
    try:
        errors = validate_identity(args.root.resolve())
    except (OSError, ValueError, KeyError, TypeError, yaml.YAMLError) as error:
        errors = [f"invalid source identity input: {error}"]
    if errors:
        print("Release identity validation failed:\n" + "\n".join(f"- {error}" for error in errors))
        return 1
    print("OK: release versions, declaration identities/checksums and mirrors agree.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
