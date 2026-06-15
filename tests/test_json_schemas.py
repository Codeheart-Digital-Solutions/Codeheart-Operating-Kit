import json
import subprocess
import sys
from pathlib import Path

from codeheart_operating_kit.manifest import load_yaml


ROOT = Path(__file__).resolve().parents[1]


def run_validator(*paths: Path) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [sys.executable, "scripts/validate-json-schemas.py", *(str(path) for path in paths)],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )


def test_json_schema_validator_passes_repository():
    result = run_validator()
    assert result.returncode == 0, result.stdout + result.stderr


def test_json_schema_validator_passes_fixture():
    result = run_validator(ROOT / "tests/fixtures/validator-valid/schema-valid.schema.json")
    assert result.returncode == 0, result.stdout + result.stderr


def test_json_schema_validator_fails_fixture():
    result = run_validator(ROOT / "tests/fixtures/validator-invalid/schema-invalid.schema.json")
    assert result.returncode == 1
    assert "required property 'missing' is not defined" in result.stdout


def test_kit_config_schema_allows_missing_setup_purpose_fixture():
    schema = json.loads((ROOT / "schemas/kit-config.schema.json").read_text(encoding="utf-8"))
    fixture = load_yaml(ROOT / "tests/fixtures/validator-valid/kit-config-without-purpose.yaml")
    assert "setup_purpose" not in schema["required"]
    assert set(schema["required"]).issubset(fixture)
    assert "setup_purpose" not in fixture


def test_kit_config_schema_preserves_existing_setup_purpose_values():
    schema = json.loads((ROOT / "schemas/kit-config.schema.json").read_text(encoding="utf-8"))
    assert schema["properties"]["setup_purpose"]["enum"] == [
        "private-automation",
        "company-automation",
        "software-product",
    ]
