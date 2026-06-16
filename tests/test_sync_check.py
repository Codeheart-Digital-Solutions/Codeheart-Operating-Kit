import json
from pathlib import Path

import codeheart_operating_kit.commands.sync as sync_command
import codeheart_operating_kit.components as components
import codeheart_operating_kit.manifest as manifest
from codeheart_operating_kit import __version__
from codeheart_operating_kit.cli import main
from codeheart_operating_kit.commands.check import check_repository, cli_available
from codeheart_operating_kit.lockfile import read_lock, write_lock
from codeheart_operating_kit.manifest import sha256_file


def test_check_detects_drift_and_sync_repairs(tmp_path):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    managed = tmp_path / ".codeheart/kit/docs/agent-interface/README.md"
    managed.write_text("drift\n", encoding="utf-8")
    report = check_repository(tmp_path)
    assert report["drift"]

    main(["sync", str(tmp_path)])
    report = check_repository(tmp_path)
    assert report["drift"] == []


def test_sync_refreshes_v011_lock_to_installed_cli_metadata(tmp_path, monkeypatch):
    monkeypatch.setenv("CODEHEART_OPERATING_KIT_CLI", "1")
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    lock = read_lock(tmp_path)
    lock["kit_version"] = "0.1.1"
    lock["release"] = {
        "asset_url": "https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.1/codeheart-operating-kit-0.1.1-macos.tar.gz",
        "checksum_sha256": "1" * 64,
    }
    discovery_path = ".codeheart/kit/docs/planning-workflows/runbooks/discovery-workflow.md"
    discovery_record = next(item for item in lock["managed_paths"] if item["path"] == discovery_path)
    discovery_record["checksum_sha256"] = "0" * 64
    write_lock(tmp_path, lock)

    assert check_repository(tmp_path)["ok"] is False

    main(["sync", str(tmp_path)])
    refreshed = read_lock(tmp_path)
    refreshed_record = next(item for item in refreshed["managed_paths"] if item["path"] == discovery_path)

    assert refreshed["kit_version"] == __version__
    assert refreshed["selected_profile"] == "standard"
    assert "planning-workflows" in refreshed["selected_components"]
    assert "v0.1.2" in refreshed["release"]["asset_url"]
    assert refreshed_record["checksum_sha256"] == sha256_file(tmp_path / discovery_path)
    assert check_repository(tmp_path)["ok"] is True


def test_sync_preserves_release_metadata_on_platform_without_cli_asset(tmp_path, monkeypatch):
    monkeypatch.setenv("CODEHEART_OPERATING_KIT_CLI", "1")
    monkeypatch.setattr("codeheart_operating_kit.commands.sync.sys.platform", "linux")
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    lock = read_lock(tmp_path)
    lock["release"] = {
        "asset_url": "local-source",
        "checksum_sha256": "0" * 64,
    }
    write_lock(tmp_path, lock)

    main(["sync", str(tmp_path)])
    release = read_lock(tmp_path)["release"]

    assert release["asset_url"] == "local-source"
    assert not release["asset_url"].endswith("bootstrap.md")


def test_sync_refreshes_release_metadata_from_packaged_resources(tmp_path, monkeypatch):
    monkeypatch.setenv("CODEHEART_OPERATING_KIT_CLI", "1")
    monkeypatch.setattr(manifest, "SOURCE_ROOT", Path("/definitely/not/a/checkout"))
    monkeypatch.setattr(components, "kit_root", manifest.kit_root)
    monkeypatch.setattr(sync_command, "kit_root", manifest.kit_root)
    monkeypatch.setattr("codeheart_operating_kit.commands.sync.sys.platform", "darwin")
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    lock = read_lock(tmp_path)
    lock["release"] = {
        "asset_url": "https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v0.1.1/codeheart-operating-kit-0.1.1-macos.tar.gz",
        "checksum_sha256": "1" * 64,
    }
    write_lock(tmp_path, lock)

    main(["sync", str(tmp_path)])
    release = read_lock(tmp_path)["release"]

    assert release["asset_url"].endswith("codeheart-operating-kit-0.1.2-macos.tar.gz")
    assert "v0.1.2" in release["asset_url"]


def test_check_json_reports_missing_cli_and_routing(tmp_path, capsys):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    capsys.readouterr()
    (tmp_path / "AGENTS.md").write_text("broken\n", encoding="utf-8")
    main(["check", str(tmp_path), "--json"])
    data = json.loads(capsys.readouterr().out)
    assert "missing_cli" in data
    assert data["missing_routing"] is True
    assert "native_capabilities" in data


def test_check_accepts_config_without_setup_purpose(tmp_path, monkeypatch):
    monkeypatch.setenv("CODEHEART_OPERATING_KIT_CLI", "1")
    assert main(["init", str(tmp_path), "--project-name", "Companyname-Automation"]) == 0
    report = check_repository(tmp_path)
    assert report["ok"] is True
    assert report["missing_lock_metadata"] == []


def test_check_recognizes_installer_wrapper_marker(monkeypatch):
    monkeypatch.setenv("CODEHEART_OPERATING_KIT_CLI", "1")

    assert cli_available() is True


def test_check_reports_missing_managed_route_target(tmp_path):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    (tmp_path / ".codeheart/kit/README.md").unlink()

    report = check_repository(tmp_path)

    assert report["missing_routing"] is True
    assert ".codeheart/kit/README.md" in report["missing_route_targets"]


def test_check_reports_missing_nested_lock_metadata(tmp_path):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    lock = read_lock(tmp_path)
    del lock["release"]["checksum_sha256"]
    del lock["update_check"]["next_update_check_due"]
    write_lock(tmp_path, lock)

    report = check_repository(tmp_path)

    assert "release.checksum_sha256" in report["missing_lock_metadata"]
    assert "update_check.next_update_check_due" in report["missing_lock_metadata"]
