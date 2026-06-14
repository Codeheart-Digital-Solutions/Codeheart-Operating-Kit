import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def run_validator(*paths: Path) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [sys.executable, "scripts/validate-markdown-headers.py", *(str(path) for path in paths)],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )


def test_markdown_header_validator_passes_repository():
    result = run_validator()
    assert result.returncode == 0, result.stdout + result.stderr


def test_markdown_header_validator_passes_fixture():
    result = run_validator(ROOT / "tests/fixtures/validator-valid/markdown-pass.md")
    assert result.returncode == 0, result.stdout + result.stderr


def test_markdown_header_validator_fails_fixture():
    result = run_validator(ROOT / "tests/fixtures/validator-invalid/markdown-fail.md")
    assert result.returncode == 1
    assert "Markdown timestamp validation failed" in result.stdout
