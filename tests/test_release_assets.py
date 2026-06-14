import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def test_release_manifest_validator_passes_fixture():
    result = subprocess.run(
        [sys.executable, "scripts/validate-release-manifest.py"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert result.returncode == 0, result.stdout + result.stderr


def test_release_manifest_validator_fails_fixture():
    result = subprocess.run(
        [
            sys.executable,
            "scripts/validate-release-manifest.py",
            "tests/fixtures/validator-invalid/release-manifest-invalid.json",
        ],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert result.returncode == 1
    assert "Release manifest validation failed" in result.stdout


def test_release_asset_build_check(tmp_path):
    result = subprocess.run(
        [sys.executable, "scripts/build-release-assets.py", "--output-dir", str(tmp_path)],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert result.returncode == 0, result.stdout + result.stderr
    for name in [
        "codeheart-operating-kit-0.1.0-macos.tar.gz",
        "codeheart-operating-kit-0.1.0-windows.zip",
    ]:
        assert (tmp_path / name).exists()
        checksum = tmp_path / f"{name}.sha256"
        assert checksum.exists()
        assert name in checksum.read_text(encoding="utf-8")


def test_release_asset_build_rejects_version_mismatch(tmp_path):
    result = subprocess.run(
        [sys.executable, "scripts/build-release-assets.py", "--version", "9.9.9", "--output-dir", str(tmp_path)],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert result.returncode != 0
    assert "does not match package version" in result.stderr
