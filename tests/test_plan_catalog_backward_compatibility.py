from __future__ import annotations

import hashlib
import io
import os
import shutil
import subprocess
import tarfile
from pathlib import Path

import pytest


ROOT = Path(__file__).resolve().parents[1]
FIXTURE = ROOT / "tests" / "fixtures" / "plans" / "migration-repository"
RELEASED_TAGS = ("v0.1.25", "v0.1.26")


def run(command: list[str], *, cwd: Path, check: bool = True) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        command,
        cwd=cwd,
        check=check,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        env={**os.environ, "GIT_TERMINAL_PROMPT": "0"},
    )


def repository_bytes(root: Path) -> dict[str, str]:
    result: dict[str, str] = {}
    for path in sorted(root.rglob("*")):
        if not path.is_file() or ".git" in path.parts:
            continue
        result[path.relative_to(root).as_posix()] = hashlib.sha256(path.read_bytes()).hexdigest()
    return result


@pytest.fixture(scope="module", params=RELEASED_TAGS)
def released_cli(request: pytest.FixtureRequest, tmp_path_factory: pytest.TempPathFactory) -> tuple[str, Path]:
    tag = str(request.param)
    archive = subprocess.run(
        ["git", "archive", "--format=tar", tag],
        cwd=ROOT,
        check=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    ).stdout
    source = tmp_path_factory.mktemp(tag.removeprefix("v") + "-source")
    with tarfile.open(fileobj=io.BytesIO(archive), mode="r:") as bundle:
        bundle.extractall(source, filter="data")
    binary_name = "codeheart-operating-kit.exe" if os.name == "nt" else "codeheart-operating-kit"
    binary = tmp_path_factory.mktemp(tag.removeprefix("v") + "-bin") / binary_name
    run(["go", "build", "-trimpath", "-buildvcs=false", "-o", str(binary), "./cmd/codeheart-operating-kit"], cwd=source)
    return tag, binary


def test_released_clis_fail_closed_on_ledger_v3_and_config_v2(
    released_cli: tuple[str, Path], tmp_path: Path
) -> None:
    tag, binary = released_cli
    repository = tmp_path / "consumer"
    shutil.copytree(FIXTURE, repository)
    run(["git", "init", "-b", "main"], cwd=repository)
    run(["git", "config", "user.name", "Compatibility Test"], cwd=repository)
    run(["git", "config", "user.email", "compatibility@example.invalid"], cwd=repository)
    run(["git", "add", "."], cwd=repository)
    run(["git", "commit", "-m", "fixture"], cwd=repository)

    compatible = run([str(binary), "plans", "validate", "--json", str(repository)], cwd=repository, check=False)
    assert compatible.returncode == 0, (tag, compatible.stdout, compatible.stderr)

    ledger = tmp_path / "reviewed-v3.yaml"
    ledger.write_text("schema_version: 3\n", encoding="utf-8")
    before_ledger_check = repository_bytes(repository)
    rejected_ledger = run(
        [str(binary), "plans", "migrate", "--ledger", str(ledger), "--dry-run", str(repository)],
        cwd=repository,
        check=False,
    )
    assert rejected_ledger.returncode != 0
    assert "unsupported plan migration schema version 3" in rejected_ledger.stderr.lower()
    assert repository_bytes(repository) == before_ledger_check

    config = repository / ".codeheart" / "kit.config.yaml"
    config.write_text(config.read_text(encoding="utf-8").replace("schema_version: 1", "schema_version: 2", 1), encoding="utf-8")
    before_config_check = repository_bytes(repository)
    rejected_config = run([str(binary), "plans", "validate", "--json", str(repository)], cwd=repository, check=False)
    assert rejected_config.returncode != 0
    assert '"code": "catalog_config_invalid"' in rejected_config.stdout
    assert repository_bytes(repository) == before_config_check
