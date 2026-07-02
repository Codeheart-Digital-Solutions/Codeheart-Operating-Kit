#!/usr/bin/env python3
from __future__ import annotations

import argparse
import re
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
HEADER = re.compile(r"^Last updated: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \(UTC\)$")
EXCLUDED_PARTS = {
    ".git",
    ".pytest_cache",
    ".venv",
    "__pycache__",
    "build",
    "dist",
    "codeheart_operating_kit.egg-info",
    "validator-invalid",
}


def iter_markdown(paths: list[Path]) -> list[Path]:
    result: list[Path] = []
    for root in paths:
        if root.is_file():
            if root.suffix == ".md":
                result.append(root)
            continue
        for path in root.rglob("*.md"):
            if not any(part in EXCLUDED_PARTS for part in path.parts):
                result.append(path)
    return sorted(result)


def invalid_headers(paths: list[Path]) -> list[str]:
    invalid: list[str] = []
    for path in paths:
        first = path.read_text(encoding="utf-8").splitlines()[0:1]
        if not first or not HEADER.match(first[0]):
            try:
                display = path.relative_to(ROOT)
            except ValueError:
                display = path
            invalid.append(str(display))
    return invalid


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate Markdown first-line timestamps.")
    parser.add_argument("paths", nargs="*", type=Path, default=[ROOT])
    args = parser.parse_args()

    invalid = invalid_headers(iter_markdown([path.resolve() for path in args.paths]))
    if invalid:
        print("Markdown timestamp validation failed.")
        for path in invalid:
            print(f"- {path}")
        return 1
    print("OK: Markdown timestamps validate.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
