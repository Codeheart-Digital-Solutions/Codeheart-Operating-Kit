from __future__ import annotations

import hashlib
import json
from importlib import resources
from pathlib import Path
from typing import Any

from . import __version__


SOURCE_ROOT = Path(__file__).resolve().parents[2]
PACKAGE_RESOURCE_ROOT = Path(resources.files("codeheart_operating_kit") / "resources")


def kit_root() -> Path:
    if (SOURCE_ROOT / "components").exists() and (SOURCE_ROOT / "profiles").exists():
        return SOURCE_ROOT
    return PACKAGE_RESOURCE_ROOT


def component_dir() -> Path:
    return kit_root() / "components"


def profile_dir() -> Path:
    return kit_root() / "profiles"


def parse_scalar(value: str) -> Any:
    value = value.strip()
    if value == "":
        return ""
    if value == "[]":
        return []
    if value == "{}":
        return {}
    if value.startswith('"') and value.endswith('"'):
        return value[1:-1]
    if value.startswith("'") and value.endswith("'"):
        return value[1:-1]
    if value == "true":
        return True
    if value == "false":
        return False
    if value.isdigit():
        return int(value)
    return value


def _strip_comments(line: str) -> str:
    in_quote = False
    quote = ""
    for index, char in enumerate(line):
        if char in {"'", '"'} and (index == 0 or line[index - 1] != "\\"):
            if in_quote and char == quote:
                in_quote = False
                quote = ""
            elif not in_quote:
                in_quote = True
                quote = char
        if char == "#" and not in_quote:
            return line[:index].rstrip()
    return line.rstrip()


def _lines(text: str) -> list[tuple[int, str]]:
    parsed: list[tuple[int, str]] = []
    for raw in text.splitlines():
        line = _strip_comments(raw)
        if not line.strip():
            continue
        parsed.append((len(line) - len(line.lstrip(" ")), line.lstrip(" ")))
    return parsed


def parse_yaml(text: str) -> Any:
    lines = _lines(text)

    def parse_block(index: int, indent: int) -> tuple[Any, int]:
        if index >= len(lines):
            return {}, index
        current_indent, content = lines[index]
        if current_indent < indent:
            return {}, index
        if content == "-" or content.startswith("- "):
            return parse_list(index, current_indent)
        return parse_map(index, current_indent)

    def parse_list(index: int, indent: int) -> tuple[list[Any], int]:
        result: list[Any] = []
        while index < len(lines):
            current_indent, content = lines[index]
            if current_indent != indent or not (content == "-" or content.startswith("- ")):
                break
            item = content[1:].strip()
            index += 1
            if item == "":
                child, index = parse_block(index, indent + 2)
                result.append(child)
            elif ":" in item:
                key, value = item.split(":", 1)
                mapping: dict[str, Any] = {}
                if value.strip():
                    mapping[key.strip()] = parse_scalar(value)
                else:
                    child, index = parse_block(index, indent + 2)
                    mapping[key.strip()] = child
                while index < len(lines):
                    next_indent, next_content = lines[index]
                    if next_indent <= indent:
                        break
                    if next_indent == indent + 2 and not next_content.startswith("- "):
                        subkey, subvalue = next_content.split(":", 1)
                        index += 1
                        if subvalue.strip():
                            mapping[subkey.strip()] = parse_scalar(subvalue)
                        else:
                            child, index = parse_block(index, next_indent + 2)
                            mapping[subkey.strip()] = child
                    else:
                        break
                result.append(mapping)
            else:
                result.append(parse_scalar(item))
        return result, index

    def parse_map(index: int, indent: int) -> tuple[dict[str, Any], int]:
        result: dict[str, Any] = {}
        while index < len(lines):
            current_indent, content = lines[index]
            if current_indent != indent or content.startswith("- "):
                break
            key, value = content.split(":", 1)
            index += 1
            if value.strip():
                result[key.strip()] = parse_scalar(value)
            else:
                child, index = parse_block(index, indent + 2)
                result[key.strip()] = child
        return result, index

    parsed, next_index = parse_block(0, 0)
    if next_index != len(lines):
        raise ValueError("Unsupported YAML structure")
    return parsed


def _quote(value: str) -> str:
    if value == "" or any(char in value for char in ":#[]{}") or value.strip() != value:
        return json.dumps(value)
    if value in {"true", "false", "null"}:
        return json.dumps(value)
    if "\n" in value:
        return json.dumps(value)
    return value


def dump_yaml(value: Any, indent: int = 0) -> str:
    lines: list[str] = []

    def emit(item: Any, level: int, key_prefix: str | None = None) -> None:
        pad = " " * level
        if isinstance(item, dict):
            if key_prefix is not None:
                lines.append(f"{pad}{key_prefix}:")
                pad = " " * (level + 2)
                level += 2
            for key, nested in item.items():
                if isinstance(nested, (dict, list)):
                    lines.append(f"{pad}{key}:")
                    emit(nested, level + 2)
                else:
                    lines.append(f"{pad}{key}: {format_scalar(nested)}")
        elif isinstance(item, list):
            if key_prefix is not None:
                lines.append(f"{pad}{key_prefix}:")
                level += 2
                pad = " " * level
            for nested in item:
                if isinstance(nested, dict):
                    lines.append(f"{pad}-")
                    emit(nested, level + 2)
                elif isinstance(nested, list):
                    lines.append(f"{pad}-")
                    emit(nested, level + 2)
                else:
                    lines.append(f"{pad}- {format_scalar(nested)}")
        elif key_prefix is not None:
            lines.append(f"{pad}{key_prefix}: {format_scalar(item)}")
        else:
            lines.append(f"{pad}{format_scalar(item)}")

    def format_scalar(item: Any) -> str:
        if isinstance(item, bool):
            return "true" if item else "false"
        if isinstance(item, int):
            return str(item)
        return _quote(str(item))

    emit(value, indent)
    return "\n".join(lines) + "\n"


def load_yaml(path: Path) -> Any:
    return parse_yaml(path.read_text(encoding="utf-8"))


def load_profile(profile_id: str = "standard") -> dict[str, Any]:
    profile_path = profile_dir() / f"{profile_id}.yaml"
    return load_yaml(profile_path)["profile"]


def component_manifest_paths() -> list[Path]:
    return sorted(component_dir().glob("*/component.yaml"))


def load_components(profile_id: str = "standard") -> list[dict[str, Any]]:
    selected = set(load_profile(profile_id)["selected_components"])
    components: list[dict[str, Any]] = []
    for path in component_manifest_paths():
        component = load_yaml(path)["component"]
        if component["id"] in selected:
            component["_manifest_path"] = str(path.relative_to(kit_root()))
            components.append(component)
    return components


def iter_component_files(profile_id: str = "standard") -> list[dict[str, Any]]:
    files: list[dict[str, Any]] = []
    for component in load_components(profile_id):
        for entry in component.get("files", []):
            item = dict(entry)
            item["component"] = component["id"]
            files.append(item)
    return files


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def validate_release_manifest(path: Path) -> dict[str, Any]:
    manifest = json.loads(path.read_text(encoding="utf-8"))
    for asset in manifest.get("assets", []):
        checksum = asset.get("sha256", "")
        if len(checksum) != 64 or any(char not in "0123456789abcdefABCDEF" for char in checksum):
            raise ValueError(f"Invalid asset checksum for {asset.get('name', '<unnamed>')}")
    return manifest


def kit_version() -> str:
    return __version__
