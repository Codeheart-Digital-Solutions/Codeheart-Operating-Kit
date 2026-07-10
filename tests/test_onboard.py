import json

from codeheart_operating_kit.cli import main
from codeheart_operating_kit.manifest import load_yaml


NEUTRAL_EXAMPLES = [
    "Yourname-Automation",
    "Companyname-Automation",
    "Productname-Development",
    "Teamname-Operations",
    "Existing-Project-Name",
]


def test_onboard_prompt_order_and_copy(tmp_path, capsys):
    main([
        "onboard",
        "--target",
        str(tmp_path),
        "--project-name",
        "Companyname-Automation",
    ])
    output = capsys.readouterr().out
    assert output.index("Choose setup language") < output.index("GPT-5.5")
    assert output.index("Do you already know what this Codex project should be called") < output.index("What is this mainly for?")
    assert output.index("What is this mainly for?") < output.index("Selected project name")
    assert "Extra High" in output
    assert "Fast" in output
    assert "Settings" in output
    assert "Work Mode" in output
    assert "Default permissions" in output
    assert "Auto review" in output
    assert "Full access" in output
    assert "Approve for me" in output
    assert "Documents > Companyname-Automation" in output
    assert "Use a different folder" in output
    assert "Desktop > Productname-Development" in output
    assert not (tmp_path / ".codeheart").exists()
    assert "portfolio" not in output.lower()
    assert "coordination" not in output.lower()
    assert "weekly update" not in output.lower()
    assert "update-check" not in output.lower()


def test_onboard_project_name_and_target_folder_help_prompts(tmp_path, capsys):
    main(["onboard", "--target", str(tmp_path), "--project-name", "Yourname-Automation"])
    output = capsys.readouterr().out
    assert "Do you already know what this Codex project should be called, or should I suggest a name?" in output
    assert "Do you already know where this project folder should be, or should I suggest a simple location?" in output
    assert "What is this mainly for?" in output
    assert "Documents > <Project-Name>" not in output


def test_onboard_setup_plan_is_shown_before_write_confirmation(tmp_path, capsys):
    main(["onboard", "--target", str(tmp_path), "--project-name", "Companyname-Automation"])
    output = capsys.readouterr().out
    assert output.index("Setup plan:") < output.index("Should I continue with this setup?")
    assert "written" not in output.lower()
    assert not (tmp_path / ".codeheart").exists()


def test_onboard_uses_neutral_example_names(tmp_path, capsys):
    main(["onboard", "--target", str(tmp_path), "--project-name", "Productname-Development"])
    output = capsys.readouterr().out
    for example in NEUTRAL_EXAMPLES:
        assert example in output


def test_onboard_yes_writes_and_creates_adoption_report(tmp_path):
    (tmp_path / "AGENTS.md").write_text("local instructions\n", encoding="utf-8")
    assert main([
        "onboard",
        "--target",
        str(tmp_path),
        "--project-name",
        "Existing-Project-Name",
        "--purpose",
        "software-product",
        "--yes",
    ]) == 0
    assert (tmp_path / ".codeheart/kit.lock.yaml").exists()
    assert (tmp_path / ".codeheart/reports/adoption-cleanup-report.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md").exists()
    assert (tmp_path / "docs/repo/plans/plan-register.md").exists()
    assert (tmp_path / "docs/repo/plans/coordination-sync-pending.md").exists()
    agents_text = (tmp_path / "AGENTS.md").read_text(encoding="utf-8")
    kit_readme_text = (tmp_path / ".codeheart/kit/README.md").read_text(encoding="utf-8")
    assert "local instructions" in agents_text
    assert "route before selecting" in agents_text
    assert ".codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md" in agents_text
    assert "repository, module, extension, or agent task is blocked by missing local tooling" in agents_text
    assert ".codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md" in agents_text
    assert "Do not check for a new Operating Kit version at session start" in agents_text
    assert "At the start of each agent session" not in agents_text
    assert "Route-before-surface standard" in kit_readme_text
    assert ".codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md" in kit_readme_text
    assert "Local machine/runtime layer" in kit_readme_text
    config = load_yaml(tmp_path / ".codeheart/kit.config.yaml")
    assert config["setup_purpose"] == "software-product"
    assert "portfolio" not in config


def test_onboard_yes_requires_target_folder(capsys):
    assert main(["onboard", "--project-name", "Companyname-Automation", "--yes", "--json"]) == 2
    data = json.loads(capsys.readouterr().out)
    assert data["written"] is False
    assert data["required_user_decisions_missing"] == ["target_folder"]


def test_onboard_yes_requires_project_name(tmp_path, capsys):
    assert main(["onboard", "--target", str(tmp_path), "--yes", "--json"]) == 2
    data = json.loads(capsys.readouterr().out)
    assert data["written"] is False
    assert data["required_user_decisions_missing"] == ["project_name"]
    assert not (tmp_path / ".codeheart").exists()


def test_onboard_without_yes_does_not_write(tmp_path, capsys):
    assert main(["onboard", "--target", str(tmp_path), "--project-name", "Companyname-Automation", "--json"]) == 0
    data = json.loads(capsys.readouterr().out)
    assert data["written"] is False
    assert data["write_approved"] is False
    assert data["required_user_decisions_missing"] == []
    assert not (tmp_path / ".codeheart").exists()


def test_onboard_omits_purpose_metadata_when_absent(tmp_path):
    assert main(["onboard", "--target", str(tmp_path), "--project-name", "Companyname-Automation", "--yes"]) == 0
    config = load_yaml(tmp_path / ".codeheart/kit.config.yaml")
    assert "setup_purpose" not in config


def test_onboard_preserves_existing_purpose_metadata_values(tmp_path):
    for purpose in ["private-automation", "company-automation", "software-product"]:
        target = tmp_path / purpose
        assert main([
            "onboard",
            "--target",
            str(target),
            "--project-name",
            "Companyname-Automation",
            "--purpose",
            purpose,
            "--yes",
        ]) == 0
        config = load_yaml(target / ".codeheart/kit.config.yaml")
        assert config["setup_purpose"] == purpose
