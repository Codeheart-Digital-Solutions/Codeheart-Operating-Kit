import json
from pathlib import Path

from codeheart_operating_kit.capabilities import (
    BASELINE_CAPABILITIES,
    PLUGIN_SELECTORS,
    check_and_install_native_capabilities,
    check_native_capabilities,
    parse_capability_status,
    refresh_native_capability_lock,
)
from codeheart_operating_kit.cli import main
from codeheart_operating_kit.lockfile import read_lock
from codeheart_operating_kit.manifest import load_yaml


FIXTURE_DIR = Path(__file__).parent / "fixtures"


def test_native_capability_plain_text_status_parsing():
    parsed = parse_capability_status("Documents Browser PDF")
    assert parsed["documents"] == "available"
    assert parsed["browser"] == "available"
    assert parsed["pdf"] == "available"
    assert parsed["spreadsheets"] == "unknown"


def test_capability_profile_parsing():
    manifest = load_yaml(Path("components/native-codex-capabilities/component.yaml"))["component"]
    capabilities = {item["id"]: item for item in manifest["capabilities"]}
    assert set(capabilities) == set(BASELINE_CAPABILITIES)
    for capability in capabilities.values():
        assert capability["purpose"]
        assert capability["profile_applicability"] == "standard"
        assert capability["plugin_selector"] == PLUGIN_SELECTORS[capability["id"]]


def test_native_capability_json_status_parsing_installed_and_missing():
    output = (FIXTURE_DIR / "capabilities-installed.json").read_text(encoding="utf-8")
    parsed = parse_capability_status(output)
    assert parsed["documents"] == "installed"
    assert parsed["browser"] == "installed"
    assert parsed["pdf"] == "missing"


def test_native_capability_json_status_parsing_blocked():
    output = (FIXTURE_DIR / "capabilities-blocked.json").read_text(encoding="utf-8")
    parsed = parse_capability_status(output)
    assert parsed["pdf"] == "blocked"


def test_check_native_capabilities_reports_unavailable_when_codex_missing():
    def runner(args):
        return json.loads((FIXTURE_DIR / "capabilities-unavailable.json").read_text(encoding="utf-8"))

    statuses = check_native_capabilities(runner)
    assert statuses["documents"]["status"] == "unavailable"
    assert statuses["documents"]["command_result_category"] == "missing-cli"


def test_check_and_install_native_capabilities_attempts_available_plugins():
    calls = []

    def runner(args):
        calls.append(args)
        if args[:4] == ["plugin", "list", "--available", "--json"]:
            return {
                "ok": True,
                "category": "success",
                "stdout": (FIXTURE_DIR / "capabilities-missing.json").read_text(encoding="utf-8"),
                "stderr": "",
            }
        return {"ok": True, "category": "success", "stdout": "{}", "stderr": ""}

    statuses = check_and_install_native_capabilities(runner)

    assert statuses["documents"]["status"] == "install-attempted"
    assert ["plugin", "add", "documents@openai-primary-runtime", "--json"] in calls


def test_refresh_native_capability_lock_updates_lock(tmp_path):
    main(["init", str(tmp_path), "--project-name", "Example-Automation"])

    def runner(args):
        return {
            "ok": True,
            "category": "success",
            "stdout": (FIXTURE_DIR / "capabilities-installed.json").read_text(encoding="utf-8"),
            "stderr": "",
        }

    refresh_native_capability_lock(tmp_path, codex_runner=runner)

    native = read_lock(tmp_path)["native_capabilities"]
    assert native["documents"]["status"] == "installed"
    assert native["pdf"]["status"] == "missing"
