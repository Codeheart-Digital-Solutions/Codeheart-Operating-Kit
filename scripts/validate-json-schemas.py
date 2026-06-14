#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
SCHEMA_DRAFT = "https://json-schema.org/draft/2020-12/schema"


def default_schemas() -> list[Path]:
    return sorted((ROOT / "schemas").glob("*.schema.json"))


def local_refs(schema: Any) -> set[str]:
    refs: set[str] = set()
    if isinstance(schema, dict):
        ref = schema.get("$ref")
        if isinstance(ref, str) and ref.startswith("#/"):
            refs.add(ref)
        for value in schema.values():
            refs.update(local_refs(value))
    elif isinstance(schema, list):
        for value in schema:
            refs.update(local_refs(value))
    return refs


def resolve_local_ref(schema: dict[str, Any], ref: str) -> Any:
    current: Any = schema
    for part in ref.removeprefix("#/").split("/"):
        if not isinstance(current, dict) or part not in current:
            raise KeyError(ref)
        current = current[part]
    return current


def validate_required_properties(schema: Any, location: str, errors: list[str]) -> None:
    if isinstance(schema, dict):
        if schema.get("type") == "object":
            properties = schema.get("properties", {})
            required = schema.get("required", [])
            if not isinstance(properties, dict):
                errors.append(f"{location}: object schema properties must be an object")
            if not isinstance(required, list):
                errors.append(f"{location}: object schema required must be an array")
            if isinstance(properties, dict) and isinstance(required, list):
                for property_name in required:
                    if property_name not in properties:
                        errors.append(f"{location}: required property {property_name!r} is not defined")
        for key, value in schema.items():
            validate_required_properties(value, f"{location}/{key}", errors)
    elif isinstance(schema, list):
        for index, value in enumerate(schema):
            validate_required_properties(value, f"{location}[{index}]", errors)


def validate_schema(path: Path) -> list[str]:
    errors: list[str] = []
    try:
        schema = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as error:
        return [f"{path}: invalid JSON: {error.msg} at line {error.lineno}, column {error.colno}"]

    if not isinstance(schema, dict):
        return [f"{path}: schema root must be an object"]
    if schema.get("$schema") != SCHEMA_DRAFT:
        errors.append(f"{path}: $schema must be {SCHEMA_DRAFT}")
    for key in ["$id", "title", "type", "properties"]:
        if key not in schema:
            errors.append(f"{path}: missing required schema key {key}")
    if schema.get("type") != "object":
        errors.append(f"{path}: root type must be object")

    for ref in local_refs(schema):
        try:
            resolve_local_ref(schema, ref)
        except KeyError:
            errors.append(f"{path}: unresolved local $ref {ref}")

    validate_required_properties(schema, str(path), errors)
    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate Operating Kit JSON schema files.")
    parser.add_argument("schemas", nargs="*", type=Path)
    args = parser.parse_args()

    schemas = [path.resolve() for path in args.schemas] if args.schemas else default_schemas()
    errors: list[str] = []
    for schema in schemas:
        errors.extend(validate_schema(schema))
    if errors:
        print("JSON schema validation failed.")
        for error in errors:
            print(f"- {error}")
        return 1
    print("OK: JSON schemas validate.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
