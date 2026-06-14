from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def test_bootstrap_documents_first_run_path():
    text = (ROOT / "bootstrap.md").read_text(encoding="utf-8")
    assert "codeheart-operating-kit onboard" in text
    assert "GPT-5.5" in text
    assert "Extra High" in text
    assert "Fast" in text
    assert "Documents > <Project-Name>" in text
    assert "stays silent" in text


def test_macos_installer_requires_checksum_and_user_level_path():
    text = (ROOT / "install.sh").read_text(encoding="utf-8")
    assert "$HOME/.codeheart/operating-kit" in text
    assert "shasum -a 256" in text
    assert "Checksum mismatch" in text
    assert "pip install --upgrade --target" in text
    assert "codeheart-operating-kit onboard" in text


def test_windows_installer_requires_checksum_and_user_level_path():
    text = (ROOT / "install.ps1").read_text(encoding="utf-8")
    assert "%LOCALAPPDATA%\\Codeheart\\OperatingKit" in text
    assert "Get-FileHash -Algorithm SHA256" in text
    assert "Checksum mismatch" in text
    assert "pip install --upgrade --target" in text
    assert "codeheart-operating-kit onboard" in text


def test_release_asset_builder_names_expected_assets():
    text = (ROOT / "scripts/build-release-assets.py").read_text(encoding="utf-8")
    assert "macos.tar.gz" in text
    assert "windows.zip" in text
    assert ".sha256" in text
    assert "asset-manifest.json" in text
