import json
import os
import re
import subprocess
import sys
from pathlib import Path

import pytest

from codeheart_operating_kit.manifest import dump_yaml, load_yaml


ROOT = Path(__file__).resolve().parents[1]
PYTHON_CLI = [sys.executable, "-m", "codeheart_operating_kit.cli"]
PARITY_FIXTURES = ROOT / "tests/fixtures/parity"
COMMANDS = ["onboard", "inspect", "init", "sync", "check", "update-check"]
ROOT_HELP_TOKENS = ["usage:", "--version", *COMMANDS]
COMMAND_HELP_TOKENS = {
    "onboard": [
        "usage:",
        "--target TARGET",
        "--project-name PROJECT_NAME",
        "--purpose {private-automation,company-automation,software-product}",
        "--yes",
        "--json",
    ],
    "inspect": ["usage:", "path", "--json"],
    "init": [
        "usage:",
        "path",
        "--project-name PROJECT_NAME",
        "--purpose {private-automation,company-automation,software-product}",
        "--selected-folder SELECTED_FOLDER",
        "--json",
    ],
    "sync": ["usage:", "path", "--release-manifest RELEASE_MANIFEST", "--json"],
    "check": ["usage:", "path", "--json"],
    "update-check": [
        "usage:",
        "path",
        "--latest-version LATEST_VERSION",
        "--metadata-url METADATA_URL",
        "--now NOW",
        "--agent-notification",
        "--json",
    ],
}


@pytest.fixture(scope="session")
def go_cli(tmp_path_factory):
    binary = tmp_path_factory.mktemp("go-cli") / "codeheart-operating-kit"
    result = subprocess.run(
        ["go", "build", "-o", str(binary), "./cmd/codeheart-operating-kit"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert result.returncode == 0, result.stdout + result.stderr
    return binary


def run_cli(command, args, *, env=None):
    result = subprocess.run(
        [str(command), *args],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
        env={**os.environ, **(env or {})},
    )
    return result


def run_python(args, *, env=None):
    python_env = {**(env or {})}
    existing = os.environ.get("PYTHONPATH")
    python_env["PYTHONPATH"] = str(ROOT / "src") if not existing else f"{ROOT / 'src'}{os.pathsep}{existing}"
    return run_cli(PYTHON_CLI[0], [*PYTHON_CLI[1:], *args], env=python_env)


def run_go(go_cli, args, *, env=None):
    return run_cli(go_cli, args, env=env)


def json_output(result):
    assert result.stdout, result.stderr
    return json.loads(result.stdout)


def normalize_lock(lock):
    lock = json.loads(json.dumps(lock))
    update = lock.get("update_check", {})
    update["last_update_check_at"] = "<time>"
    update["next_update_check_due"] = "<time>"
    native = lock.get("native_capabilities", {})
    for record in native.values():
        record["checked_at"] = "<time>"
    return lock


def normalize_cli_text(value):
    value = value.replace("\r\n", "\n").replace("\\", "/")
    return re.sub(r"\s+", " ", value).strip()


def normalize_path_text(value):
    return str(value).replace("\\", "/")


def assert_help_surface_equivalent(py_text, go_text, tokens):
    py_normalized = normalize_cli_text(py_text)
    go_normalized = normalize_cli_text(go_text)
    for token in tokens:
        assert normalize_cli_text(token) in py_normalized
        assert normalize_cli_text(token) in go_normalized


def write_yaml(path, value):
    text = dump_yaml(value)
    # The retained Python dumper does not quote digit-only SHA-256 strings. Keep
    # fixture mutations type-preserving so lock-v2 tests exercise the intended drift.
    text = re.sub(r"(:\s+)([0-9]{64})(\n)", r'\1"\2"\3', text)
    path.write_text(text, encoding="utf-8")


def file_tree(root):
    return sorted(path.relative_to(root).as_posix() for path in root.rglob("*") if path.is_file())


def setup_folder_state(base, case):
    target = base / case["id"]
    setup = case["setup"]
    if setup == "missing":
        return target
    target.mkdir()
    if setup == "empty":
        return target
    if setup == "existing":
        (target / "note.txt").write_text("hello\n", encoding="utf-8")
    elif setup == "technical":
        (target / "pyproject.toml").write_text("[project]\nname = 'x'\n", encoding="utf-8")
    elif setup == "existing-kit":
        (target / ".codeheart").mkdir()
        (target / ".codeheart/kit.lock.yaml").write_text("schema_version: 1\n", encoding="utf-8")
    elif setup == "file":
        target.rmdir()
        target.write_text("x\n", encoding="utf-8")
    elif setup == "many-files":
        for index in range(101):
            (target / f"file-{index:03d}.txt").write_text("x\n", encoding="utf-8")
    else:
        raise AssertionError(f"unknown setup {setup}")
    return target


def test_root_surface_parity(go_cli):
    for args, tokens in [
        (["--help"], ROOT_HELP_TOKENS),
        (["onboard", "--help"], COMMAND_HELP_TOKENS["onboard"]),
        (["inspect", "--help"], COMMAND_HELP_TOKENS["inspect"]),
        (["init", "--help"], COMMAND_HELP_TOKENS["init"]),
        (["sync", "--help"], COMMAND_HELP_TOKENS["sync"]),
        (["check", "--help"], COMMAND_HELP_TOKENS["check"]),
        (["update-check", "--help"], COMMAND_HELP_TOKENS["update-check"]),
        (["init", "--json", "--help"], COMMAND_HELP_TOKENS["init"]),
        (["update-check", "--json", "--agent-notification", "--help"], COMMAND_HELP_TOKENS["update-check"]),
    ]:
        py = run_python(args)
        go = run_go(go_cli, args)
        assert py.returncode == go.returncode
        assert py.stderr == go.stderr == ""
        assert_help_surface_equivalent(py.stdout, go.stdout, tokens)

    py = run_python(["--version"])
    go = run_go(go_cli, ["--version"])
    assert py.returncode == go.returncode == 0
    assert py.stdout == go.stdout == "codeheart-operating-kit 0.1.23\n"
    assert py.stderr == go.stderr == ""

    py = run_python([])
    go = run_go(go_cli, [])
    assert py.returncode == go.returncode == 2
    assert "required: command" in py.stderr
    assert "required: command" in go.stderr
    assert_help_surface_equivalent(py.stderr, go.stderr, ROOT_HELP_TOKENS)

    py = run_python(["missing"])
    go = run_go(go_cli, ["missing"])
    assert py.returncode == go.returncode == 2
    assert "invalid choice: 'missing'" in py.stderr
    assert "invalid choice: 'missing'" in go.stderr
    assert_help_surface_equivalent(py.stderr, go.stderr, ROOT_HELP_TOKENS)


def test_inspect_parity(go_cli, tmp_path):
    cases = json.loads((PARITY_FIXTURES / "folder-states.json").read_text(encoding="utf-8"))
    for case in cases:
        target = setup_folder_state(tmp_path, case)
        py = run_python(["inspect", str(target), "--json"])
        go = run_go(go_cli, ["inspect", str(target), "--json"])
        assert py.returncode == go.returncode == 0, case["id"]
        py_data = json_output(py)
        go_data = json_output(go)
        assert normalize_path_text(py_data["path"]) == normalize_path_text(go_data["path"])
        assert py_data["mode"] == go_data["mode"] == case["expected_mode"]
        assert py_data["reason"] == go_data["reason"]
        if "markers" in py_data:
            assert py_data["markers"] == go_data["markers"]


def test_onboard_argument_gating_parity(go_cli, tmp_path):
    py = run_python(["onboard", "--project-name", "Companyname-Automation", "--yes", "--json"])
    go = run_go(go_cli, ["onboard", "--project-name", "Companyname-Automation", "--yes", "--json"])
    assert py.returncode == go.returncode == 2
    assert json_output(py)["required_user_decisions_missing"] == json_output(go)["required_user_decisions_missing"]


def test_argument_parser_parity_for_equals_and_missing_values(go_cli, tmp_path):
    py_equal = tmp_path / "py-equals"
    go_equal = tmp_path / "go-equals"
    py = run_python(["init", str(py_equal), "--project-name=EqualStyle"])
    go = run_go(go_cli, ["init", str(go_equal), "--project-name=EqualStyle"])
    assert py.returncode == go.returncode == 0
    assert load_yaml(py_equal / ".codeheart/kit.config.yaml")["project_display_name"] == "EqualStyle"
    assert load_yaml(go_equal / ".codeheart/kit.config.yaml")["project_display_name"] == "EqualStyle"

    py_bad = tmp_path / "py-bad"
    go_bad = tmp_path / "go-bad"
    py = run_python(["init", str(py_bad), "--project-name", "--json"])
    go = run_go(go_cli, ["init", str(go_bad), "--project-name", "--json"])
    assert py.returncode == go.returncode == 2
    assert not (py_bad / ".codeheart").exists()
    assert not (go_bad / ".codeheart").exists()

    py_update = tmp_path / "py-update"
    go_update = tmp_path / "go-update"
    assert run_python(["init", str(py_update), "--project-name=Update"]).returncode == 0
    assert run_go(go_cli, ["init", str(go_update), "--project-name=Update"]).returncode == 0
    py = run_python(["update-check", str(py_update), "--latest-version", "--json"])
    go = run_go(go_cli, ["update-check", str(go_update), "--latest-version", "--json"])
    assert py.returncode == go.returncode == 2
    assert load_yaml(py_update / ".codeheart/kit.lock.yaml")["update_check"]["latest_seen_version"] != "--json"
    assert load_yaml(go_update / ".codeheart/kit.lock.yaml")["update_check"]["latest_seen_version"] != "--json"

    for args in [
        ["init", "--project-name", "--help"],
        ["init", "--project-name", "--json", "--help"],
        ["init", "--project-name", "-h"],
        ["onboard", "--target", "--help"],
        ["onboard", "--target", "--yes", "--help"],
        ["update-check", "--latest-version", "--help"],
        ["update-check", "--latest-version", "--json", "--help"],
    ]:
        py = run_python(args)
        go = run_go(go_cli, args)
        assert py.returncode == go.returncode == 2

    py = run_python(["onboard", "--target", str(tmp_path), "--yes", "--json"])
    go = run_go(go_cli, ["onboard", "--target", str(tmp_path), "--yes", "--json"])
    assert py.returncode == go.returncode == 2
    assert json_output(py)["required_user_decisions_missing"] == json_output(go)["required_user_decisions_missing"]


def test_onboard_approved_native_capability_difference(go_cli, tmp_path):
    py_target = tmp_path / "python"
    go_target = tmp_path / "go"
    py = run_python(["onboard", "--target", str(py_target), "--project-name", "Companyname-Automation", "--yes", "--json"])
    go = run_go(go_cli, ["onboard", "--target", str(go_target), "--project-name", "Companyname-Automation", "--yes", "--json"])
    assert py.returncode == go.returncode == 0
    py_data = json_output(py)
    go_data = json_output(go)
    assert py_data["written"] is True
    assert go_data["written"] is True
    assert "native_capabilities" in py_data
    assert "native_capabilities" not in go_data
    assert "Should I check and set up these tools now?" in "\n".join(py_data["script"])
    assert "Should I check and set up these tools now?" not in "\n".join(go_data["script"])

    py_help = run_python(["onboard", "--help"])
    go_help = run_go(go_cli, ["onboard", "--help"])
    assert "Base onboarding does not install, offer, or implicitly check optional native capabilities." not in py_help.stdout
    assert "Base onboarding does not install, offer, or implicitly check optional native capabilities." in go_help.stdout

    go_lock = load_yaml(go_target / ".codeheart/kit.lock.yaml")
    for record in go_lock["native_capabilities"].values():
        assert record["status"] == "unknown"
        assert record["command_result_category"] == "not-checked"


def test_init_generated_state_parity(go_cli, tmp_path):
    py_target = tmp_path / "python"
    go_target = tmp_path / "go"
    for target in [py_target, go_target]:
        (target / "docs/repo/plans").mkdir(parents=True)
        (target / "docs/repo/plans/plan-register.md").write_text("custom plan register\n", encoding="utf-8")
        (target / "AGENTS.md").write_text("# Existing Agent Notes\n\nlocal instructions\n", encoding="utf-8")
        (target / ".gitignore").write_text("local.tmp\n", encoding="utf-8")
    py = run_python(["init", str(py_target), "--project-name", "Example-Automation", "--purpose", "company-automation"])
    go = run_go(go_cli, ["init", str(go_target), "--project-name", "Example-Automation", "--purpose", "company-automation"])
    assert py.returncode == go.returncode == 0
    assert file_tree(py_target) == file_tree(go_target)
    for target in [py_target, go_target]:
        assert (target / "docs/repo/plans/plan-register.md").read_text(encoding="utf-8") == "custom plan register\n"
        agents = (target / "AGENTS.md").read_text(encoding="utf-8")
        assert "local instructions" in agents
        assert "<!-- BEGIN CODEHEART OPERATING KIT MANAGED BLOCK -->" in agents
        gitignore = (target / ".gitignore").read_text(encoding="utf-8")
        assert "local.tmp" in gitignore
        assert ".codeheart/user/preferences.yaml" in gitignore
        assert ".codeheart/user/feedback/" in gitignore
        assert ".codeheart/local/" in gitignore

    py_config = load_yaml(py_target / ".codeheart/kit.config.yaml")
    go_config = load_yaml(go_target / ".codeheart/kit.config.yaml")
    py_config["selected_setup_folder"] = "<target>"
    go_config["selected_setup_folder"] = "<target>"
    assert py_config == go_config

    py_lock = normalize_lock(load_yaml(py_target / ".codeheart/kit.lock.yaml"))
    go_lock = normalize_lock(load_yaml(go_target / ".codeheart/kit.lock.yaml"))
    assert py_lock["selected_profile"] == go_lock["selected_profile"]
    assert py_lock["selected_components"] == go_lock["selected_components"]
    assert {item["path"] for item in py_lock["managed_paths"]} == {item["path"] for item in go_lock["managed_paths"]}
    py_surfaces = {item["path"] for item in py_lock["generated_surfaces"]}
    go_surfaces = {item["path"] for item in go_lock["generated_surfaces"]}
    # Go lock v2 records the complete typed graph, including transaction and ignore
    # surfaces; Python lock v1 retains the older directory summary.
    assert py_surfaces - {".codeheart/kit/"} <= go_surfaces
    assert {".codeheart/kit.transaction.json", ".codeheart/local/kit-transactions/", ".gitignore"} <= go_surfaces
    assert go_lock["schema_version"] == 2
    assert go_lock["managed_sections"]
    assert go_lock["release_provenance"]["verification_status"] == "local-source"
    assert py_lock["native_capabilities"] == go_lock["native_capabilities"]


def test_sync_and_check_parity(go_cli, tmp_path):
    py_target = tmp_path / "python"
    go_target = tmp_path / "go"
    assert run_python(["init", str(py_target), "--project-name", "Example-Automation"]).returncode == 0
    assert run_go(go_cli, ["init", str(go_target), "--project-name", "Example-Automation"]).returncode == 0
    release_manifest = tmp_path / "release-manifest.json"
    release_url = "https://example.test/codeheart-operating-kit-0.1.23-universal.zip"
    release_sha = "a" * 64
    release_manifest.write_text(
        json.dumps(
            {
                "assets": [
                    {
                        "name": "codeheart-operating-kit-0.1.23-universal.zip",
                        "platform": "universal",
                        "url": release_url,
                        "sha256": release_sha,
                    }
                ]
            }
        ),
        encoding="utf-8",
    )

    for target in [py_target, go_target]:
        (target / "docs/repo/plans/plan-register.md").write_text("custom plan register\n", encoding="utf-8")
        (target / "docs/repo/plans/coordination-sync-pending.md").unlink()
        managed = target / ".codeheart/kit/docs/agent-interface/README.md"
        managed.write_text("drift\n", encoding="utf-8")
        agents = target / "AGENTS.md"
        text = agents.read_text(encoding="utf-8")
        _before, rest = text.split("<!-- BEGIN CODEHEART OPERATING KIT MANAGED BLOCK -->", 1)
        _managed, after = rest.split("<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->", 1)
        agents.write_text(
            "local prefix\n"
            "<!-- BEGIN CODEHEART OPERATING KIT MANAGED BLOCK -->\n"
            "stale managed block\n"
            "<!-- END CODEHEART OPERATING KIT MANAGED BLOCK -->"
            f"{after}\nlocal suffix\n",
            encoding="utf-8",
        )
        lock = load_yaml(target / ".codeheart/kit.lock.yaml")
        lock["generated_surfaces"] = [
            item
            for item in lock["generated_surfaces"]
            if item["path"] != "docs/repo/plans/coordination-sync-pending.md"
        ]
        lock["generated_surfaces"].append({"path": "docs/custom-local.md", "ownership": "scaffold"})
        write_yaml(target / ".codeheart/kit.lock.yaml", lock)

    env = {"CODEHEART_OPERATING_KIT_CLI": "1"}
    py_check = json_output(run_python(["check", str(py_target), "--json"], env=env))
    go_check = json_output(run_go(go_cli, ["check", str(go_target), "--json"], env=env))
    assert bool(py_check["drift"]) is True
    assert bool(go_check["drift"]) is True

    assert run_python(["sync", str(py_target), "--release-manifest", str(release_manifest)]).returncode == 0
    assert run_go(go_cli, ["sync", str(go_target), "--release-manifest", str(release_manifest)]).returncode == 0
    assert (py_target / "docs/repo/plans/plan-register.md").read_text(encoding="utf-8") == "custom plan register\n"
    assert (go_target / "docs/repo/plans/plan-register.md").read_text(encoding="utf-8") == "custom plan register\n"
    for target in [py_target, go_target]:
        assert (target / "docs/repo/plans/coordination-sync-pending.md").exists()
        agents = (target / "AGENTS.md").read_text(encoding="utf-8")
        assert "local prefix" in agents
        assert "local suffix" in agents
        assert "stale managed block" not in agents
        assert "Operation routing and dispatch" in agents
        lock = load_yaml(target / ".codeheart/kit.lock.yaml")
        generated = {item["path"] for item in lock["generated_surfaces"]}
        assert "docs/repo/plans/coordination-sync-pending.md" in generated
        if target == py_target:
            assert "docs/custom-local.md" in generated
            assert lock["release"] == {"asset_url": release_url, "checksum_sha256": release_sha}
        else:
            # Go v2 resolves lock records from declarations and treats the release
            # manifest flag as validation-only during same-version sync.
            assert "docs/custom-local.md" not in generated
            assert lock["release"]["asset_url"] == "local-source"

    py_after = json_output(run_python(["check", str(py_target), "--json"], env=env))
    go_after = json_output(run_go(go_cli, ["check", str(go_target), "--json"], env=env))
    assert py_after["ok"] == go_after["ok"] == True
    assert py_after["drift"] == go_after["drift"] == []

    for target in [py_target, go_target]:
        (target / ".codeheart/kit/README.md").unlink()
    py_missing = json_output(run_python(["check", str(py_target), "--json"], env=env))
    go_missing = json_output(run_go(go_cli, ["check", str(go_target), "--json"], env=env))
    assert ".codeheart/kit/README.md" in py_missing["missing_route_targets"]
    assert ".codeheart/kit/README.md" in go_missing["missing_route_targets"]


def test_check_stale_cli_and_missing_lock_metadata_parity(go_cli, tmp_path):
    py_target = tmp_path / "python"
    go_target = tmp_path / "go"
    assert run_python(["init", str(py_target), "--project-name", "Example-Automation"]).returncode == 0
    assert run_go(go_cli, ["init", str(go_target), "--project-name", "Example-Automation"]).returncode == 0
    for target in [py_target, go_target]:
        lock = load_yaml(target / ".codeheart/kit.lock.yaml")
        lock["kit_version"] = "0.1.0"
        del lock["release"]["checksum_sha256"]
        del lock["update_check"]["next_update_check_due"]
        write_yaml(target / ".codeheart/kit.lock.yaml", lock)

    env = {"CODEHEART_OPERATING_KIT_CLI": "1"}
    py = json_output(run_python(["check", str(py_target), "--json"], env=env))
    go = json_output(run_go(go_cli, ["check", str(go_target), "--json"], env=env))
    assert py["ok"] == go["ok"] == False
    assert py["stale_cli"] == go["stale_cli"] == True
    assert "release.checksum_sha256" in py["missing_lock_metadata"]
    assert "release.checksum_sha256" in go["missing_lock_metadata"]
    assert "update_check.next_update_check_due" in py["missing_lock_metadata"]
    assert "update_check.next_update_check_due" in go["missing_lock_metadata"]


def test_update_check_parity(go_cli, tmp_path):
    py_target = tmp_path / "python"
    go_target = tmp_path / "go"
    assert run_python(["init", str(py_target), "--project-name", "Example-Automation"]).returncode == 0
    assert run_go(go_cli, ["init", str(go_target), "--project-name", "Example-Automation"]).returncode == 0

    for args, expected in [
        (["--latest-version", "0.1.23", "--now", "2026-06-13T00:00:00Z", "--json"], "current"),
        (["--latest-version", "0.2.0", "--now", "2026-06-13T00:00:00Z", "--json"], "update-available"),
    ]:
        py = json_output(run_python(["update-check", str(py_target), *args]))
        go = json_output(run_go(go_cli, ["update-check", str(go_target), *args]))
        assert py["status"] == go["status"] == expected
        assert py["next_update_check_due"] == go["next_update_check_due"] == "2026-06-20T00:00:00Z"

    py_metadata = py_target / "latest-release.json"
    go_metadata = go_target / "latest-release.json"
    py_metadata.write_text(json.dumps({"tag_name": "v0.3.0"}), encoding="utf-8")
    go_metadata.write_text(json.dumps({"tag_name": "v0.3.0"}), encoding="utf-8")
    py = json_output(run_python(["update-check", str(py_target), "--metadata-url", py_metadata.as_uri(), "--json"]))
    go = json_output(run_go(go_cli, ["update-check", str(go_target), "--metadata-url", go_metadata.as_uri(), "--json"]))
    assert py["status"] == go["status"] == "update-available"
    assert py["latest_seen_version"] == go["latest_seen_version"] == "v0.3.0"

    py_failed = json_output(run_python(["update-check", str(py_target), "--metadata-url", (py_target / "missing.json").as_uri(), "--json"]))
    go_failed = json_output(run_go(go_cli, ["update-check", str(go_target), "--metadata-url", (go_target / "missing.json").as_uri(), "--json"]))
    assert py_failed["status"] == go_failed["status"] == "failed"

    py_current = run_python(["update-check", str(py_target), "--latest-version", "0.1.23", "--agent-notification"])
    go_current = run_go(go_cli, ["update-check", str(go_target), "--latest-version", "0.1.23", "--agent-notification"])
    assert py_current.returncode == go_current.returncode == 0
    assert py_current.stdout == go_current.stdout == ""

    py_available = run_python(["update-check", str(py_target), "--latest-version", "0.2.0", "--agent-notification"])
    go_available = run_go(go_cli, ["update-check", str(go_target), "--latest-version", "0.2.0", "--agent-notification"])
    assert py_available.returncode == go_available.returncode == 0
    assert normalize_cli_text(py_available.stdout) == normalize_cli_text(go_available.stdout)
    assert "update available" in py_available.stdout.lower()
