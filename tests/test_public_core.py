import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def run_validator(*paths: Path) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [sys.executable, "scripts/validate-public-core.py", *(str(path) for path in paths)],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )


def test_public_core_validator_passes_repository():
    result = run_validator()
    assert result.returncode == 0, result.stdout + result.stderr


def test_public_core_validator_blocks_private_patterns():
    result = run_validator(ROOT / "tests/fixtures/validator-invalid/public-core-blocked.txt")
    assert result.returncode == 1
    assert "blocked public-core pattern" in result.stdout


def test_public_core_validator_blocks_legacy_onboarding_examples():
    result = run_validator(ROOT / "tests/fixtures/validator-invalid/public-core-legacy-onboarding-examples.txt")
    assert result.returncode == 1
    assert "legacy onboarding example" in result.stdout


def test_public_core_validator_allows_placeholders():
    result = run_validator(ROOT / "tests/fixtures/validator-valid/public-core-placeholder.txt")
    assert result.returncode == 0, result.stdout + result.stderr


def test_public_core_validator_allows_neutral_onboarding_examples():
    result = run_validator(ROOT / "tests/fixtures/validator-valid/public-core-onboarding-placeholders.txt")
    assert result.returncode == 0, result.stdout + result.stderr
