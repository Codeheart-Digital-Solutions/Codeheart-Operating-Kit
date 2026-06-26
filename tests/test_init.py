from codeheart_operating_kit.cli import main
from codeheart_operating_kit.components import ensure_gitignore
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
        "docs/repo/plans/plan-register.md",
        "docs/repo/plans/coordination-sync-pending.md",
        "docs/agent-memory/README.md",
        "docs/agent-memory/goal-register.md",
        "docs/agent-memory/session-ledger.md",
        "docs/agent-memory/untriaged-sessions.md",
        ".gitignore",
    ]
    for relative in expected:
        assert (tmp_path / relative).exists(), relative
    gitignore = (tmp_path / ".gitignore").read_text(encoding="utf-8")
    assert ".codeheart/user/preferences.yaml" in gitignore
    assert ".codeheart/user/*.local.yaml" in gitignore
    assert ".codeheart/user/feedback/" in gitignore
    assert ".codeheart/local/" in gitignore
    assert not (tmp_path / ".codeheart/local").exists()
    lock = read_lock(tmp_path)
    assert lock["selected_profile"] == "standard"
    assert ".codeheart/kit/README.md" in {item["path"] for item in lock["managed_paths"]}
    generated = {item["path"] for item in lock["generated_surfaces"]}
    assert "docs/repo/plans/plan-register.md" in generated
    assert "docs/repo/plans/coordination-sync-pending.md" in generated
    assert ".codeheart/local/" not in generated
    assert set(lock["native_capabilities"]) == {"documents", "spreadsheets", "presentations", "browser", "pdf"}
    config = load_yaml(tmp_path / ".codeheart/kit.config.yaml")
    assert config["setup_purpose"] == "company-automation"
    assert config["local_consumer_layer"]["local_machine_layer_path"] == ".codeheart/local/"
    assert "portfolio" not in config


def test_init_can_omit_purpose_metadata(tmp_path):
    assert main(["init", str(tmp_path), "--project-name", "Companyname-Automation"]) == 0
    config = load_yaml(tmp_path / ".codeheart/kit.config.yaml")
    assert "setup_purpose" not in config
    assert config["project_display_name"] == "Companyname-Automation"


def test_gitignore_adds_feedback_draft_path_to_existing_local_user_block(tmp_path):
    (tmp_path / ".gitignore").write_text(
        "# Codeheart Operating Kit local user layer\n"
        ".codeheart/user/preferences.yaml\n"
        ".codeheart/user/*.local.yaml\n",
        encoding="utf-8",
    )

    assert ensure_gitignore(tmp_path) is True

    gitignore = (tmp_path / ".gitignore").read_text(encoding="utf-8")
    assert ".codeheart/user/feedback/" in gitignore
    assert ".codeheart/local/" in gitignore
    assert gitignore.count("# Codeheart Operating Kit local user layer") == 1
    assert gitignore.count("# Codeheart Operating Kit local machine layer") == 1
