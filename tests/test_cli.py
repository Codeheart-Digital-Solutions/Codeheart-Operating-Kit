import pytest

from codeheart_operating_kit.cli import main


def test_help_lists_commands(capsys):
    with pytest.raises(SystemExit) as exc:
        main(["--help"])
    assert exc.value.code == 0
    output = capsys.readouterr().out
    for command in ["onboard", "inspect", "init", "sync", "check", "update-check"]:
        assert command in output


@pytest.mark.parametrize("command", ["onboard", "inspect", "init", "sync", "check", "update-check"])
def test_subcommand_help(command, capsys):
    with pytest.raises(SystemExit) as exc:
        main([command, "--help"])
    assert exc.value.code == 0
    assert "usage:" in capsys.readouterr().out


def test_onboard_help_describes_non_interactive_requirements(capsys):
    with pytest.raises(SystemExit) as exc:
        main(["onboard", "--help"])
    assert exc.value.code == 0
    output = capsys.readouterr().out
    normalized = " ".join(output.split())
    assert "Non-interactive --yes writes require explicit --target and --project-name values" in normalized
    assert "Required with --yes" in normalized
    assert "does not select a different profile" in normalized
