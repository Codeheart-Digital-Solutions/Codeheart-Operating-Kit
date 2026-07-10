import hashlib
import json
import os
import stat
import subprocess
import zipfile
from pathlib import Path

from codeheart_operating_kit import __version__


ROOT = Path(__file__).resolve().parents[1]
INSTALLER_FIXTURES = ROOT / "tests/fixtures/installer"


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def fake_cli_script(version: str = __version__) -> str:
    return "\n".join(
        [
            "#!/usr/bin/env sh",
            'if [ "$1" = "--version" ]; then',
            f'  echo "codeheart-operating-kit {version}"',
            "  exit 0",
            "fi",
            'if [ "$1" = "--help" ]; then',
            '  echo "usage: codeheart-operating-kit"',
            "  exit 0",
            "fi",
            'if [ "$1" = "__verify-content-identity" ]; then',
            "  exit 0",
            "fi",
            'if [ "$1" = "__verify-release-evidence" ]; then',
            "  exit 0",
            "fi",
            "exit 0",
            "",
        ]
    )


def write_pack(tmp_path: Path, *, content: str | None = None, include_binary: bool = True) -> tuple[Path, Path]:
    payload = tmp_path / "payload" / f"codeheart-operating-kit-{__version__}-macos-universal"
    files: list[Path] = []
    if include_binary:
        binary = payload / "bin/codeheart-operating-kit"
        binary.parent.mkdir(parents=True)
        binary.write_text(content if content is not None else fake_cli_script(), encoding="utf-8")
        binary.chmod(0o755)
        files.append(binary)
    else:
        payload.mkdir(parents=True)
        binary = payload / "bin/codeheart-operating-kit"
    content_manifest = payload / "content-manifest.yaml"
    content_manifest.write_text((ROOT / "manifest.yaml").read_text(encoding="utf-8"), encoding="utf-8")
    files.append(content_manifest)
    readme = payload / "README.md"
    readme.write_text("fixture\n", encoding="utf-8")
    files.append(readme)
    checksums = payload / "checksums.txt"
    checksums.write_text(
        "".join(f"{sha256_file(path)}  {path.relative_to(payload).as_posix()}\n" for path in sorted(files)),
        encoding="utf-8",
    )
    pack_manifest = payload / "pack-manifest.json"
    pack_manifest.write_text(
        json.dumps(
            {
                "schema_version": 1,
                "version": __version__,
                "platform": "macos-universal",
                "command": "codeheart-operating-kit",
                "binary_path": "bin/codeheart-operating-kit",
                "binary_sha256": sha256_file(binary) if include_binary else "0" * 64,
                "content_manifest_path": "content-manifest.yaml",
                "content_manifest_sha256": sha256_file(content_manifest),
                "payload_checksums_path": "checksums.txt",
                "payload_checksums_sha256": sha256_file(checksums),
            },
            indent=2,
            sort_keys=True,
        )
        + "\n",
        encoding="utf-8",
    )
    pack = tmp_path / f"codeheart-operating-kit-{__version__}-macos-universal.zip"
    with zipfile.ZipFile(pack, "w", compression=zipfile.ZIP_DEFLATED) as archive:
        if payload.exists():
            for path in sorted(payload.rglob("*")):
                archive.write(path, path.relative_to(payload.parent))
    checksum = pack.with_name(pack.name + ".sha256")
    checksum.write_text(f"{sha256_file(pack)}  {pack.name}\n", encoding="utf-8")
    write_catalog(tmp_path, pack, sha256_file(pack_manifest))
    return pack, checksum


def write_catalog(tmp_path: Path, pack: Path, pack_manifest_sha256: str = "0" * 64) -> Path:
    catalog = tmp_path / f"release-catalog-{__version__}.json"
    catalog.write_text(
        json.dumps(
            {
                "schema_version": 1,
                "version": __version__,
                "assets": [
                    {
                        "name": pack.name,
                        "version": __version__,
                        "platform": "macos-universal",
                        "url": pack.name,
                        "archive_sha256": sha256_file(pack),
                        "pack_manifest_sha256": pack_manifest_sha256,
                    }
                ],
            },
            indent=2,
            sort_keys=True,
        )
        + "\n",
        encoding="utf-8",
    )
    return catalog


def run_macos_installer(tmp_path: Path, args: list[str], *, install_dir: Path | None = None) -> subprocess.CompletedProcess[str]:
    install_dir = install_dir or tmp_path / "install"
    env = os.environ.copy()
    env["PATH"] = "/usr/bin:/bin"
    return subprocess.run(
        [
            "bash",
            str(ROOT / "install.sh"),
            "--version",
            __version__,
            "--install-dir",
            str(install_dir),
            *args,
        ],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
        env=env,
    )


def test_bootstrap_documents_first_run_path():
    text = (ROOT / "bootstrap.md").read_text(encoding="utf-8")
    assert "codeheart-operating-kit onboard" in text
    assert f"https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v{__version__}" in text
    assert f"releases/download/v{__version__}/install.sh" in text
    assert f"releases/download/v{__version__}/install.ps1" in text
    assert f"codeheart-operating-kit-{__version__}-macos-universal.zip" in text
    assert f"codeheart-operating-kit-{__version__}-windows-x64.zip" in text
    assert "Should I check and set up these tools now?" not in text
    assert "Do not offer optional native capability installation during base onboarding." in text
    assert "pip install" not in text
    assert "GPT-5.5" in text
    assert "Extra High" in text
    assert "Fast" in text
    assert "Documents > <Project-Name>" in text
    assert "If everything is current, Codex will not mention it" in text


def test_macos_installer_requires_checksum_and_user_level_path():
    text = (ROOT / "install.sh").read_text(encoding="utf-8")
    assert "$HOME/.codeheart/operating-kit" in text
    assert "shasum -a 256" in text
    assert "Checksum mismatch" in text
    assert "macos-universal" in text
    assert "release-catalog" in text
    assert "pack-manifest.json" in text
    assert "unzip -q" in text
    assert "bin/codeheart-operating-kit" in text
    assert "Staged binary validation failed; previous runnable command preserved." in text
    assert "symbolic link or unsupported filesystem entry" in text
    assert "pip install" not in text
    assert "PIP_CONFIG_FILE" not in text
    assert "codeheart-operating-kit onboard" in text


def test_windows_installer_requires_checksum_and_user_level_path():
    text = (ROOT / "install.ps1").read_text(encoding="utf-8")
    assert r"%LOCALAPPDATA%\Codeheart\OperatingKit" in text
    assert "Get-FileHash -Algorithm SHA256" in text
    assert "Checksum mismatch" in text
    assert "windows-x64" in text
    assert "release-catalog" in text
    assert "pack-manifest.json" in text
    assert "Expand-Archive" in text
    assert "bin/codeheart-operating-kit.exe" in text
    assert "Staged binary validation failed; previous runnable command preserved." in text
    assert "$AssetParent = (Get-Location).Path" in text
    assert "Release catalog asset name does not match" in text
    assert "pip install" not in text
    assert "PIP_CONFIG_FILE" not in text
    assert "CODEHEART_OPERATING_KIT_CLI=1" in text
    assert "codeheart-operating-kit onboard" in text
    assert "-Python PATH" not in text


def test_release_asset_builder_names_expected_assets():
    text = (ROOT / "scripts/build-release-assets.py").read_text(encoding="utf-8")
    assert "macos-universal" in text
    assert "windows-x64" in text
    assert ".sha256" in text
    assert "release-catalog" in text


def test_macos_installer_help_hides_deprecated_python_flag():
    result = subprocess.run(
        ["bash", str(ROOT / "install.sh"), "--help"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert result.returncode == 0
    assert "--python" not in result.stdout


def test_macos_installer_installs_local_binary_pack_without_python(tmp_path):
    pack, checksum = write_pack(tmp_path)
    install_dir = tmp_path / "install"

    result = run_macos_installer(
        tmp_path,
        [
            "--asset-file",
            str(pack),
            "--checksum-file",
            str(checksum),
            "--python",
            "/does/not/exist/python",
        ],
        install_dir=install_dir,
    )

    assert result.returncode == 0, result.stdout + result.stderr
    assert "--python is deprecated and ignored" in result.stderr
    target = install_dir / "bin/codeheart-operating-kit"
    assert target.exists()
    version = subprocess.run([str(target), "--version"], text=True, capture_output=True, check=False)
    assert version.stdout == f"codeheart-operating-kit {__version__}\n"
    assert not (install_dir / "lib").exists()
    assert "Add this folder to PATH" in result.stdout


def test_macos_installer_fails_closed_on_checksum_mismatch(tmp_path):
    pack, _checksum = write_pack(tmp_path)
    result = run_macos_installer(tmp_path, ["--asset-file", str(pack), "--checksum", "0" * 64])

    assert result.returncode != 0
    assert "checksum" in result.stderr.lower()
    assert not (tmp_path / "install/bin/codeheart-operating-kit").exists()


def test_macos_installer_uses_file_url_checksum_sidecar(tmp_path):
    asset_dir = tmp_path / "asset space"
    asset_dir.mkdir()
    pack, _checksum = write_pack(asset_dir)
    result = run_macos_installer(tmp_path, ["--asset-url", pack.as_uri(), "--catalog-file", str(asset_dir / f"release-catalog-{__version__}.json")])

    assert result.returncode == 0, result.stdout + result.stderr
    assert (tmp_path / "install/bin/codeheart-operating-kit").exists()


def test_macos_installer_rejects_malformed_archive_with_valid_checksum(tmp_path):
    pack = tmp_path / f"codeheart-operating-kit-{__version__}-macos-universal.zip"
    pack.write_text("not a zip\n", encoding="utf-8")
    checksum = sha256_file(pack)
    write_catalog(tmp_path, pack)

    result = run_macos_installer(tmp_path, ["--asset-file", str(pack), "--checksum", checksum])

    assert result.returncode != 0
    assert "could not be extracted" in result.stderr
    assert not (tmp_path / "install/bin/codeheart-operating-kit").exists()


def test_macos_installer_rejects_unsafe_catalog_asset_name(tmp_path):
    pack, _checksum = write_pack(tmp_path)
    catalog_path = tmp_path / f"release-catalog-{__version__}.json"
    catalog = json.loads(catalog_path.read_text(encoding="utf-8"))
    catalog["assets"][0]["name"] = f"../codeheart-operating-kit-{__version__}-macos-universal.zip"
    catalog_path.write_text(json.dumps(catalog, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    result = run_macos_installer(tmp_path, ["--asset-file", str(pack), "--catalog-file", str(catalog_path)])

    assert result.returncode != 0
    assert "asset name does not match" in result.stderr
    assert not (tmp_path / "install/bin/codeheart-operating-kit").exists()


def test_macos_installer_rejects_symbolic_link_archive_entry(tmp_path):
    pack, _checksum = write_pack(tmp_path)
    entry = zipfile.ZipInfo(f"codeheart-operating-kit-{__version__}-macos-universal/linked")
    entry.create_system = 3
    entry.external_attr = (stat.S_IFLNK | 0o777) << 16
    with zipfile.ZipFile(pack, "a", compression=zipfile.ZIP_DEFLATED) as archive:
        archive.writestr(entry, "outside")
    manifest = tmp_path / "payload" / f"codeheart-operating-kit-{__version__}-macos-universal" / "pack-manifest.json"
    write_catalog(tmp_path, pack, sha256_file(manifest))

    result = run_macos_installer(tmp_path, ["--asset-file", str(pack)])

    assert result.returncode != 0
    assert "symbolic link or unsupported filesystem entry" in result.stderr
    assert not (tmp_path / "install/bin/codeheart-operating-kit").exists()


def test_macos_installer_rejects_archive_missing_binary(tmp_path):
    pack, checksum = write_pack(tmp_path, include_binary=False)
    result = run_macos_installer(tmp_path, ["--asset-file", str(pack), "--checksum-file", str(checksum)])

    assert result.returncode != 0
    assert "Binary checksum mismatch" in result.stderr
    assert not (tmp_path / "install/bin/codeheart-operating-kit").exists()


def test_macos_installer_preserves_previous_command_when_staged_validation_fails(tmp_path):
    pack, checksum = write_pack(tmp_path, content="not an executable binary\n")
    install_dir = tmp_path / "install"
    target = install_dir / "bin/codeheart-operating-kit"
    target.parent.mkdir(parents=True)
    target.write_text("previous runnable\n", encoding="utf-8")

    result = run_macos_installer(tmp_path, ["--asset-file", str(pack), "--checksum-file", str(checksum)], install_dir=install_dir)

    assert result.returncode != 0
    assert "previous runnable command preserved" in result.stderr
    assert target.read_text(encoding="utf-8") == "previous runnable\n"


def test_macos_installer_detects_legacy_python_install_and_preserves_files(tmp_path):
    pack, checksum = write_pack(tmp_path)
    install_dir = tmp_path / "install"
    target = install_dir / "bin/codeheart-operating-kit"
    legacy_lib = install_dir / "lib/codeheart_operating_kit"
    target.parent.mkdir(parents=True)
    legacy_lib.mkdir(parents=True)
    target.write_text((INSTALLER_FIXTURES / "legacy-python-wrapper.sh").read_text(encoding="utf-8"), encoding="utf-8")
    legacy_marker = legacy_lib / "__init__.py"
    legacy_marker.write_text("# preserved\n", encoding="utf-8")

    result = run_macos_installer(tmp_path, ["--asset-file", str(pack), "--checksum-file", str(checksum)], install_dir=install_dir)

    assert result.returncode == 0, result.stdout + result.stderr
    assert "Legacy Python install detected" in result.stdout
    assert "Legacy Python files were preserved" in result.stdout
    assert legacy_marker.read_text(encoding="utf-8") == "# preserved\n"
    assert "codeheart_operating_kit.cli" not in target.read_text(encoding="utf-8")
