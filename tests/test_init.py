from codeheart_operating_kit.cli import main
from codeheart_operating_kit.lockfile import read_lock
from codeheart_operating_kit.manifest import load_yaml


def test_init_writes_standard_surfaces(tmp_path):
    assert main([
        "init",
        str(tmp_path),
        "--project-name",
        "Example-Automation",
        "--purpose",
        "company-automation",
    ]) == 0
    expected = [
        ".codeheart/kit",
        ".codeheart/kit/README.md",
        ".codeheart/kit.lock.yaml",
        ".codeheart/kit.config.yaml",
        ".codeheart/user/README.md",
        ".codeheart/user/examples/preferences.yaml",
        "AGENTS.md",
        "docs/repo/README.md",
        "docs/agent-memory/README.md",
        "docs/agent-memory/goal-register.md",
        "docs/agent-memory/session-ledger.md",
        "docs/agent-memory/untriaged-sessions.md",
        ".gitignore",
    ]
    for relative in expected:
        assert (tmp_path / relative).exists(), relative
    assert ".codeheart/user/preferences.yaml" in (tmp_path / ".gitignore").read_text(encoding="utf-8")
    lock = read_lock(tmp_path)
    assert lock["selected_profile"] == "standard"
    assert ".codeheart/kit/README.md" in {item["path"] for item in lock["managed_paths"]}
    assert set(lock["native_capabilities"]) == {"documents", "spreadsheets", "presentations", "browser", "pdf"}
    config = load_yaml(tmp_path / ".codeheart/kit.config.yaml")
    assert config["setup_purpose"] == "company-automation"


def test_init_can_omit_purpose_metadata(tmp_path):
    assert main(["init", str(tmp_path), "--project-name", "Companyname-Automation"]) == 0
    config = load_yaml(tmp_path / ".codeheart/kit.config.yaml")
    assert "setup_purpose" not in config
    assert config["project_display_name"] == "Companyname-Automation"
