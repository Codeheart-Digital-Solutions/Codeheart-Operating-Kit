import io
import hashlib
import subprocess
import sys
import tarfile
import zipfile
from pathlib import Path

import pytest

from codeheart_operating_kit import __version__
from codeheart_operating_kit.manifest import load_yaml


ROOT = Path(__file__).resolve().parents[1]


def read_packaged_manifest_from_wheel(wheel_bytes: bytes) -> str:
    with zipfile.ZipFile(io.BytesIO(wheel_bytes)) as wheel:
        return wheel.read("codeheart_operating_kit/resources/manifest.yaml").decode("utf-8")


def read_wheel_from_tarball(path: Path) -> bytes:
    with tarfile.open(path, "r:gz") as archive:
        member = next(item for item in archive.getmembers() if item.name.endswith(".whl"))
        stream = archive.extractfile(member)
        assert stream is not None
        return stream.read()


def read_wheel_from_zip(path: Path) -> bytes:
    with zipfile.ZipFile(path) as archive:
        name = next(item for item in archive.namelist() if item.endswith(".whl"))
        return archive.read(name)


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


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
        f"codeheart-operating-kit-{__version__}-macos.tar.gz",
        f"codeheart-operating-kit-{__version__}-windows.zip",
    ]:
        assert (tmp_path / name).exists()
        checksum = tmp_path / f"{name}.sha256"
        assert checksum.exists()
        assert name in checksum.read_text(encoding="utf-8")


def test_release_asset_build_ignores_private_pip_index_env(tmp_path, monkeypatch):
    private_index = "https://aws:secret@example-private.invalid/simple/"
    monkeypatch.setenv("PIP_INDEX_URL", private_index)
    result = subprocess.run(
        [sys.executable, "scripts/build-release-assets.py", "--output-dir", str(tmp_path)],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    output = result.stdout + result.stderr
    assert result.returncode == 0, output
    assert "example-private.invalid" not in output
    assert "aws:secret" not in output


def test_release_assets_embed_current_packaged_manifest(tmp_path):
    result = subprocess.run(
        [sys.executable, "scripts/build-release-assets.py", "--output-dir", str(tmp_path)],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert result.returncode == 0, result.stdout + result.stderr

    expected = (ROOT / "src/codeheart_operating_kit/resources/manifest.yaml").read_text(encoding="utf-8")
    expected_data = load_yaml(ROOT / "src/codeheart_operating_kit/resources/manifest.yaml")
    assert {asset["sha256"] for asset in expected_data["assets"]} == {"0" * 64}

    archive_wheels = [
        read_wheel_from_tarball(tmp_path / f"codeheart-operating-kit-{__version__}-macos.tar.gz"),
        read_wheel_from_zip(tmp_path / f"codeheart-operating-kit-{__version__}-windows.zip"),
    ]
    for wheel_bytes in archive_wheels:
        assert read_packaged_manifest_from_wheel(wheel_bytes) == expected


def test_current_dist_assets_match_root_manifest_when_present():
    manifest = load_yaml(ROOT / "manifest.yaml")
    local_paths = {
        "bootstrap.md": ROOT / "bootstrap.md",
        "install.sh": ROOT / "install.sh",
        "install.ps1": ROOT / "install.ps1",
        "release-notes.md": ROOT / "release-notes.md",
        f"codeheart-operating-kit-{__version__}-macos.tar.gz": ROOT / "dist" / f"codeheart-operating-kit-{__version__}-macos.tar.gz",
        f"codeheart-operating-kit-{__version__}-macos.tar.gz.sha256": ROOT / "dist" / f"codeheart-operating-kit-{__version__}-macos.tar.gz.sha256",
        f"codeheart-operating-kit-{__version__}-windows.zip": ROOT / "dist" / f"codeheart-operating-kit-{__version__}-windows.zip",
        f"codeheart-operating-kit-{__version__}-windows.zip.sha256": ROOT / "dist" / f"codeheart-operating-kit-{__version__}-windows.zip.sha256",
    }
    missing = [path for path in local_paths.values() if not path.exists()]
    if missing:
        pytest.skip("release dist assets are not present in this checkout")

    for asset in manifest["assets"]:
        path = local_paths[asset["name"]]
        assert sha256_file(path) == asset["sha256"]


def test_release_manifest_component_impact_matches_component_manifests():
    manifest = load_yaml(ROOT / "manifest.yaml")
    packaged = load_yaml(ROOT / "src/codeheart_operating_kit/resources/manifest.yaml")
    for release_manifest in [manifest, packaged]:
        for component in release_manifest["components"]:
            component_manifest = load_yaml(ROOT / component["manifest_path"])["component"]
            assert component["consumer_impact"] == component_manifest["consumer_impact"]


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
