import json

from codeheart_operating_kit.cli import main
from codeheart_operating_kit.commands.check import check_repository, cli_available
from codeheart_operating_kit.lockfile import read_lock, write_lock


def test_check_detects_drift_and_sync_repairs(tmp_path):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    managed = tmp_path / ".codeheart/kit/docs/agent-interface/README.md"
    managed.write_text("drift\n", encoding="utf-8")
    report = check_repository(tmp_path)
    assert report["drift"]

    main(["sync", str(tmp_path)])
    report = check_repository(tmp_path)
    assert report["drift"] == []


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
