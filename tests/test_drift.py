from codeheart_operating_kit.cli import main
from codeheart_operating_kit.commands.check import check_repository


def test_drift_validator_detects_and_repairs_managed_file(tmp_path):
    main(["init", str(tmp_path), "--project-name", "Drift-Automation"])
    managed = tmp_path / ".codeheart/kit/docs/planning-workflows/README.md"
    managed.write_text("changed\n", encoding="utf-8")

    assert check_repository(tmp_path)["drift"]

    main(["sync", str(tmp_path)])

    assert check_repository(tmp_path)["drift"] == []
