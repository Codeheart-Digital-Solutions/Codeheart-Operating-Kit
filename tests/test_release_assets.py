import hashlib
import json
import stat
import subprocess
import sys
import zipfile
from pathlib import Path

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
    assert all(len(asset["archive_sha256"]) == 64 for asset in assets)
    assert all(len(asset["pack_manifest_sha256"]) == 64 for asset in assets)


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
    catalog = tmp_path / f"release-catalog-{__version__}.json"
    assert catalog.exists()
    payload = json.loads(catalog.read_text(encoding="utf-8"))
    assert payload["version"] == __version__
    assert all("archive_sha256" in asset and "pack_manifest_sha256" in asset for asset in payload["assets"])


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
    outputs = [tmp_path / "first", tmp_path / "second"]
    for output in outputs:
        result = subprocess.run(
            [sys.executable, "scripts/build-release-assets.py", "--platform", "windows-x64", "--output-dir", str(output)],
            cwd=ROOT,
            text=True,
            capture_output=True,
            check=False,
        )
        assert result.returncode == 0, result.stdout + result.stderr
    assert (outputs[0] / f"codeheart-operating-kit-{__version__}-windows-x64.zip").exists()
    assert not (outputs[0] / f"codeheart-operating-kit-{__version__}-macos-universal.zip").exists()
    for name in [
        f"codeheart-operating-kit-{__version__}-windows-x64.zip",
        f"codeheart-operating-kit-{__version__}-windows-x64.zip.sha256",
        f"release-catalog-{__version__}.json",
    ]:
        assert (outputs[0] / name).read_bytes() == (outputs[1] / name).read_bytes()
    candidate = json.loads((outputs[0] / f"release-catalog-{__version__}.json").read_text(encoding="utf-8"))
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
            assert any(name.endswith("pack-manifest.json") for name in names)
            assert any(name.endswith("checksums.txt") for name in names)
            assert any(name.endswith("content-manifest.yaml") for name in names)
            assert not any(name.endswith(".whl") for name in names)
            assert not any(".dist-info/" in name for name in names)
            assert not any("codeheart_operating_kit/" in name for name in names)
            manifest_name = next(name for name in names if name.endswith("pack-manifest.json"))
            payload_manifest = json.loads(archive.read(manifest_name).decode("utf-8"))
            assert payload_manifest["version"] == __version__
            assert payload_manifest["binary_path"] == binary_suffix
            assert len(payload_manifest["binary_sha256"]) == 64
            assert len(payload_manifest["content_manifest_sha256"]) == 64
            assert len(payload_manifest["payload_checksums_sha256"]) == 64
            checksums_name = next(name for name in names if name.endswith("checksums.txt"))
            checksums_text = archive.read(checksums_name).decode("utf-8")
            assert binary_suffix in checksums_text
            assert "content-manifest.yaml" in checksums_text
            assert "pack-manifest.json" not in checksums_text

            for info in archive.infolist():
                assert info.date_time == (1980, 1, 1, 0, 0, 0)
                assert info.create_system == 3
                assert stat.S_IFMT(info.external_attr >> 16) == stat.S_IFREG

    candidate = json.loads((tmp_path / f"release-catalog-{__version__}.json").read_text(encoding="utf-8"))
    assert {asset["platform"] for asset in candidate["assets"]} == {"macos-universal", "windows-x64"}
    assert {asset["name"] for asset in candidate["assets"]} == set(expected_binary)
    assert all(len(asset["archive_sha256"]) == 64 for asset in candidate["assets"])
    assert all(len(asset["pack_manifest_sha256"]) == 64 for asset in candidate["assets"])


def test_current_dist_assets_match_root_manifest_when_present():
    manifest = load_yaml(ROOT / "manifest.yaml")
    assert "assets" not in manifest
    assert "released_at" not in manifest
    assert set(manifest["compatibility"]["platforms"]) == {"macos-universal", "windows-x64"}


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
