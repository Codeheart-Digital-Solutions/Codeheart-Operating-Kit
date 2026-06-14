from __future__ import annotations

import shutil
import subprocess


def run_codex(args: list[str], timeout: int = 20) -> dict[str, object]:
    executable = shutil.which("codex")
    if not executable:
        return {"ok": False, "category": "missing-cli", "stdout": "", "stderr": "codex not found"}
    try:
        completed = subprocess.run(
            [executable, *args],
            text=True,
            capture_output=True,
            timeout=timeout,
            check=False,
        )
    except subprocess.TimeoutExpired:
        return {"ok": False, "category": "timeout", "stdout": "", "stderr": "codex timed out"}
    return {
        "ok": completed.returncode == 0,
        "category": "success" if completed.returncode == 0 else "failed",
        "stdout": completed.stdout,
        "stderr": completed.stderr,
        "returncode": completed.returncode,
    }
