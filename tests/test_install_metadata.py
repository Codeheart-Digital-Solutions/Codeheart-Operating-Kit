from pathlib import Path

from codeheart_operating_kit import __version__


ROOT = Path(__file__).resolve().parents[1]


def test_bootstrap_documents_first_run_path():
    text = (ROOT / "bootstrap.md").read_text(encoding="utf-8")
    assert "codeheart-operating-kit onboard" in text
    assert f"https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/tag/v{__version__}" in text
    assert f"releases/download/v{__version__}/install.sh" in text
    assert f"releases/download/v{__version__}/install.ps1" in text
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
    assert "pip install --no-index --no-deps --upgrade --target" in text
    assert "PIP_CONFIG_FILE=/dev/null" in text
    assert "--no-index --no-deps" in text
    assert "CODEHEART_OPERATING_KIT_CLI=1" in text
    assert "codeheart-operating-kit onboard" in text


def test_windows_installer_requires_checksum_and_user_level_path():
    text = (ROOT / "install.ps1").read_text(encoding="utf-8")
    assert r"%LOCALAPPDATA%\Codeheart\OperatingKit" in text
    assert "Get-FileHash -Algorithm SHA256" in text
    assert "Checksum mismatch" in text
    assert "pip install --no-index --no-deps --upgrade --target" in text
    assert '$env:PIP_CONFIG_FILE = "NUL"' in text
    assert "--no-index --no-deps" in text
    assert "CODEHEART_OPERATING_KIT_CLI=1" in text
    assert "codeheart-operating-kit onboard" in text


def test_release_asset_builder_names_expected_assets():
    text = (ROOT / "scripts/build-release-assets.py").read_text(encoding="utf-8")
    assert "macos.tar.gz" in text
    assert "windows.zip" in text
    assert ".sha256" in text
    assert "asset-manifest.json" in text
