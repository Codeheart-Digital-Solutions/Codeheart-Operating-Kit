from __future__ import annotations

from pathlib import Path
from typing import Any

from .manifest import sha256_file


def drift_report(root: Path, lock: dict[str, Any]) -> list[dict[str, str]]:
    findings: list[dict[str, str]] = []
    for record in lock.get("managed_paths", []):
        path = root / record["path"]
        expected = record.get("checksum_sha256")
        if not path.exists():
            findings.append({"path": record["path"], "status": "missing"})
            continue
        actual = sha256_file(path)
        if expected and actual != expected:
            findings.append({"path": record["path"], "status": "drift", "expected": expected, "actual": actual})
    return findings
