from codeheart_operating_kit.cli import main


def test_onboard_prompt_order_and_copy(tmp_path, capsys):
    main([
        "onboard",
        "--target",
        str(tmp_path),
        "--project-name",
        "Bluebird-Automation",
        "--purpose",
        "company-automation",
    ])
    output = capsys.readouterr().out
    assert output.index("Choose setup language") < output.index("GPT-5.5")
    assert output.index("What are you setting this up for?") < output.index("What is the company")
    assert "Extra High" in output
    assert "Fast" in output
    assert "Settings" in output
    assert "Work Mode" in output
    assert "Default permissions" in output
    assert "Auto review" in output
    assert "Full access" in output
    assert "Approve for me" in output
    assert "Documents > Bluebird-Automation" in output
    assert "Use a different folder" in output
    assert "Desktop > Booking-App" in output
    assert not (tmp_path / ".codeheart").exists()


def test_onboard_yes_writes_and_creates_adoption_report(tmp_path):
    (tmp_path / "AGENTS.md").write_text("local instructions\n", encoding="utf-8")
    assert main([
        "onboard",
        "--target",
        str(tmp_path),
        "--project-name",
        "Existing-App",
        "--purpose",
        "software-product",
        "--yes",
    ]) == 0
    assert (tmp_path / ".codeheart/kit.lock.yaml").exists()
    assert (tmp_path / ".codeheart/reports/adoption-cleanup-report.md").exists()
    assert "local instructions" in (tmp_path / "AGENTS.md").read_text(encoding="utf-8")
