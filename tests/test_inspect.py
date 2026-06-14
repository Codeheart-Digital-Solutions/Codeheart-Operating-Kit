from pathlib import Path

from codeheart_operating_kit.commands.inspect import inspect_folder


def test_new_folder_setup_for_missing_folder(tmp_path):
    result = inspect_folder(tmp_path / "new")
    assert result["mode"] == "new-folder-setup"


def test_existing_folder_setup(tmp_path):
    (tmp_path / "note.txt").write_text("hello", encoding="utf-8")
    result = inspect_folder(tmp_path)
    assert result["mode"] == "existing-folder-setup"


def test_existing_technical_project_adoption(tmp_path):
    (tmp_path / "pyproject.toml").write_text("[project]\nname = 'x'\n", encoding="utf-8")
    result = inspect_folder(tmp_path)
    assert result["mode"] == "existing-technical-project-adoption"


def test_existing_operating_kit_repair(tmp_path):
    (tmp_path / ".codeheart").mkdir()
    (tmp_path / ".codeheart/kit.lock.yaml").write_text("schema_version: 1\n", encoding="utf-8")
    result = inspect_folder(tmp_path)
    assert result["mode"] == "existing-operating-kit-repair"


def test_ambiguous_folder_stop_for_file(tmp_path):
    file_path = tmp_path / "file.txt"
    file_path.write_text("x", encoding="utf-8")
    result = inspect_folder(file_path)
    assert result["mode"] == "ambiguous-folder-stop"
