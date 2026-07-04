import hashlib
import json
import subprocess
import sys
import zipfile
from pathlib import Path

import pytest

from codeheart_operating_kit import __version__
from codeheart_operating_kit.manifest import load_yaml


ROOT = Path(__file__).resolve().parents[1]


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


def test_release_candidate_fixture_uses_binary_platforms():
    fixture = json.loads((ROOT / "tests/fixtures/release-candidate/release-candidate-manifest.json").read_text(encoding="utf-8"))
    assets = fixture["assets"]
    assert {asset["name"] for asset in assets} == {
        f"codeheart-operating-kit-{__version__}-macos-universal.zip",
        f"codeheart-operating-kit-{__version__}-windows-x64.zip",
    }
    assert {asset["platform"] for asset in assets} == {"macos-universal", "windows-x64"}
    assert all(len(asset["sha256"]) == 64 for asset in assets)


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
        f"codeheart-operating-kit-{__version__}-macos-universal.zip",
        f"codeheart-operating-kit-{__version__}-windows-x64.zip",
    ]:
        assert (tmp_path / name).exists()
        checksum = tmp_path / f"{name}.sha256"
        assert checksum.exists()
        assert name in checksum.read_text(encoding="utf-8")
    assert (tmp_path / f"release-candidate-manifest-{__version__}.json").exists()


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


def test_release_asset_build_can_target_windows_x64_only(tmp_path):
    result = subprocess.run(
        [sys.executable, "scripts/build-release-assets.py", "--platform", "windows-x64", "--output-dir", str(tmp_path)],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert result.returncode == 0, result.stdout + result.stderr
    assert (tmp_path / f"codeheart-operating-kit-{__version__}-windows-x64.zip").exists()
    assert not (tmp_path / f"codeheart-operating-kit-{__version__}-macos-universal.zip").exists()
    candidate = json.loads((tmp_path / f"release-candidate-manifest-{__version__}.json").read_text(encoding="utf-8"))
    assert [asset["platform"] for asset in candidate["assets"]] == ["windows-x64"]


def test_release_assets_contain_binaries_and_no_python_payload(tmp_path):
    result = subprocess.run(
        [sys.executable, "scripts/build-release-assets.py", "--output-dir", str(tmp_path)],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert result.returncode == 0, result.stdout + result.stderr

    expected_binary = {
        f"codeheart-operating-kit-{__version__}-macos-universal.zip": "bin/codeheart-operating-kit",
        f"codeheart-operating-kit-{__version__}-windows-x64.zip": "bin/codeheart-operating-kit.exe",
    }
    for archive_name, binary_suffix in expected_binary.items():
        with zipfile.ZipFile(tmp_path / archive_name) as archive:
            names = archive.namelist()
            assert any(name.endswith(binary_suffix) for name in names)
            assert any(name.endswith("bootstrap.md") for name in names)
            assert any(name.endswith("INSTALL.md") for name in names)
            assert any(name.endswith("manifest.json") for name in names)
            assert any(name.endswith("checksums.txt") for name in names)
            assert not any(name.endswith("manifest.yaml") for name in names)
            assert not any(name.endswith(".whl") for name in names)
            assert not any(".dist-info/" in name for name in names)
            assert not any("codeheart_operating_kit/" in name for name in names)
            manifest_name = next(name for name in names if name.endswith("manifest.json"))
            payload_manifest = json.loads(archive.read(manifest_name).decode("utf-8"))
            assert payload_manifest["version"] == __version__
            assert payload_manifest["binary"] == binary_suffix
            checksums_name = next(name for name in names if name.endswith("checksums.txt"))
            checksums_text = archive.read(checksums_name).decode("utf-8")
            assert binary_suffix in checksums_text
            assert "manifest.json" in checksums_text

    candidate = json.loads((tmp_path / f"release-candidate-manifest-{__version__}.json").read_text(encoding="utf-8"))
    assert {asset["platform"] for asset in candidate["assets"]} == {"macos-universal", "windows-x64"}
    assert {asset["name"] for asset in candidate["assets"]} == set(expected_binary)
    assert all(len(asset["sha256"]) == 64 for asset in candidate["assets"])


def test_current_dist_assets_match_root_manifest_when_present():
    manifest = load_yaml(ROOT / "manifest.yaml")
    source_assets = {"bootstrap.md", "install.sh", "install.ps1", "release-notes.md"}
    local_paths = {
        "bootstrap.md": ROOT / "bootstrap.md",
        "install.sh": ROOT / "install.sh",
        "install.ps1": ROOT / "install.ps1",
        "release-notes.md": ROOT / "release-notes.md",
        f"codeheart-operating-kit-{__version__}-macos-universal.zip": ROOT / "dist" / f"codeheart-operating-kit-{__version__}-macos-universal.zip",
        f"codeheart-operating-kit-{__version__}-macos-universal.zip.sha256": ROOT / "dist" / f"codeheart-operating-kit-{__version__}-macos-universal.zip.sha256",
        f"codeheart-operating-kit-{__version__}-windows-x64.zip": ROOT / "dist" / f"codeheart-operating-kit-{__version__}-windows-x64.zip",
        f"codeheart-operating-kit-{__version__}-windows-x64.zip.sha256": ROOT / "dist" / f"codeheart-operating-kit-{__version__}-windows-x64.zip.sha256",
    }
    missing = [path for path in local_paths.values() if not path.exists()]
    if missing:
        pytest.skip("release dist assets are not present in this checkout")

    unpublished_source_changes = []
    for asset in manifest["assets"]:
        path = local_paths[asset["name"]]
        actual = sha256_file(path)
        if asset["name"] in source_assets and actual != asset["sha256"]:
            unpublished_source_changes.append(asset["name"])
            continue
        assert actual == asset["sha256"]
    assert set(unpublished_source_changes).issubset(source_assets)


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
