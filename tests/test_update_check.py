import json

from codeheart_operating_kit.cli import main
from codeheart_operating_kit.lockfile import read_lock


def test_update_check_writes_weekly_cadence(tmp_path):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    main([
        "update-check",
        str(tmp_path),
        "--latest-version",
        "0.1.0",
        "--now",
        "2026-06-13T00:00:00Z",
    ])
    update = read_lock(tmp_path)["update_check"]
    assert update["last_update_check_at"] == "2026-06-13T00:00:00Z"
    assert update["next_update_check_due"] == "2026-06-20T00:00:00Z"
    assert update["update_status"] == "current"


def test_update_check_silent_current_for_agent_notification(tmp_path, capsys):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    capsys.readouterr()
    main(["update-check", str(tmp_path), "--latest-version", "0.1.0", "--agent-notification"])
    assert capsys.readouterr().out == ""


def test_update_check_prompts_for_available_update(tmp_path, capsys):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    capsys.readouterr()
    main(["update-check", str(tmp_path), "--latest-version", "0.2.0", "--agent-notification"])
    assert "update available" in capsys.readouterr().out.lower()


def test_update_check_reads_latest_release_metadata_url(tmp_path, capsys):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    metadata = tmp_path / "latest-release.json"
    metadata.write_text(json.dumps({"tag_name": "v0.2.0"}), encoding="utf-8")
    capsys.readouterr()

    main([
        "update-check",
        str(tmp_path),
        "--metadata-url",
        metadata.as_uri(),
        "--now",
        "2026-06-13T00:00:00Z",
        "--json",
    ])

    output = json.loads(capsys.readouterr().out)
    assert output["status"] == "update-available"
    update = read_lock(tmp_path)["update_check"]
    assert update["latest_seen_version"] == "v0.2.0"
    assert update["next_update_check_due"] == "2026-06-20T00:00:00Z"


def test_update_check_reports_metadata_lookup_failure_without_advancing_due_date(tmp_path, capsys):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    before = read_lock(tmp_path)["update_check"]["next_update_check_due"]
    capsys.readouterr()

    main([
        "update-check",
        str(tmp_path),
        "--metadata-url",
        (tmp_path / "missing-release.json").as_uri(),
        "--now",
        "2026-06-13T00:00:00Z",
        "--json",
    ])

    output = json.loads(capsys.readouterr().out)
    assert output["status"] == "failed"
    update = read_lock(tmp_path)["update_check"]
    assert update["update_status"] == "failed"
    assert update["next_update_check_due"] == before


def test_update_check_json(tmp_path, capsys):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])
    capsys.readouterr()
    main(["update-check", str(tmp_path), "--latest-version", "0.2.0", "--json"])
    assert json.loads(capsys.readouterr().out)["status"] == "update-available"
