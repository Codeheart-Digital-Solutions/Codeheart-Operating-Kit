import json
import subprocess
import sys
from pathlib import Path

from codeheart_operating_kit.manifest import load_yaml


ROOT = Path(__file__).resolve().parents[1]


def validate_instance(schema, instance):
    errors = []

    def validate(subschema, value, location):
        if not isinstance(subschema, dict):
            return
        expected_type = subschema.get("type")
        if expected_type == "object" or "properties" in subschema or "required" in subschema:
            if not isinstance(value, dict):
                errors.append(f"{location}: expected object")
                return
            properties = subschema.get("properties", {})
            for required in subschema.get("required", []):
                if required not in value:
                    errors.append(f"{location}: missing {required}")
            if subschema.get("additionalProperties") is False:
                for key in value:
                    if key not in properties:
                        errors.append(f"{location}: unknown {key}")
            for key, property_schema in properties.items():
                if key in value:
                    validate(property_schema, value[key], f"{location}.{key}")
        elif expected_type == "integer" and not isinstance(value, int):
            errors.append(f"{location}: expected integer")
        elif expected_type == "string":
            if not isinstance(value, str):
                errors.append(f"{location}: expected string")
            elif len(value) < subschema.get("minLength", 0):
                errors.append(f"{location}: too short")

        if "const" in subschema and value != subschema["const"]:
            errors.append(f"{location}: expected const {subschema['const']}")
        if "enum" in subschema and value not in subschema["enum"]:
            errors.append(f"{location}: invalid enum {value}")

        for clause in subschema.get("allOf", []):
            condition = clause.get("if")
            then = clause.get("then")
            if condition and then and condition_matches(condition, value):
                validate(then, value, location)

    def condition_matches(condition, value):
        before = list(errors)
        validate(condition, value, "$condition")
        matched = errors == before
        del errors[len(before):]
        return matched

    validate(schema, instance, "$")
    return errors


def base_config():
    return {
        "schema_version": 1,
        "selected_profile": "standard",
        "project_display_name": "Example-Automation",
        "selected_setup_folder": "/tmp/Example-Automation",
        "local_consumer_layer": {
            "repo_docs_path": "docs/repo/",
            "agent_memory_path": "docs/agent-memory/",
            "user_layer_path": ".codeheart/user/",
            "local_machine_layer_path": ".codeheart/local/",
        },
        "component_settings": {},
    }


def github_repo_feedback():
    return {
        "mode": "github_issues",
        "destination": {
            "type": "github_issues",
            "owner": "Codeheart-Digital-Solutions",
            "repo": "Codeheart-Operating-Kit",
        },
        "authorization": {
            "organization": "Codeheart-Digital-Solutions",
            "require_verified_membership": True,
            "require_gh_cli": True,
            "unavailable_behavior": "silent",
        },
        "github_standardization": {
            "labels": "not_configured",
            "issue_templates": "not_configured",
        },
    }


def kit_config_schema():
    return json.loads((ROOT / "schemas/kit-config.schema.json").read_text(encoding="utf-8"))


def assert_config_valid(config):
    errors = validate_instance(kit_config_schema(), config)
    assert errors == []


def assert_config_invalid(config, expected):
    errors = validate_instance(kit_config_schema(), config)
    assert any(expected in error for error in errors), errors


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


def test_kit_config_schema_accepts_no_portfolio_block():
    config = base_config()

    assert_config_valid(config)


def test_kit_config_schema_accepts_no_repo_feedback_block():
    config = base_config()

    assert_config_valid(config)


def test_kit_config_schema_accepts_valid_repo_feedback_github_issues_config():
    config = base_config()
    config["repo_feedback"] = github_repo_feedback()

    assert_config_valid(config)


def test_kit_config_schema_accepts_valid_repo_feedback_disabled_config():
    config = base_config()
    config["repo_feedback"] = {
        "mode": "disabled",
        "disabled_reason": "verified_maintainer_declined_issue_intake",
    }

    assert_config_valid(config)


def test_kit_config_schema_rejects_repo_feedback_github_issues_without_destination_type():
    config = base_config()
    config["repo_feedback"] = github_repo_feedback()
    del config["repo_feedback"]["destination"]["type"]

    assert_config_invalid(config, "missing type")


def test_kit_config_schema_rejects_repo_feedback_github_issues_without_destination_owner_or_repo():
    for missing_field in ["owner", "repo"]:
        config = base_config()
        config["repo_feedback"] = github_repo_feedback()
        del config["repo_feedback"]["destination"][missing_field]

        assert_config_invalid(config, f"missing {missing_field}")


def test_kit_config_schema_rejects_repo_feedback_github_issues_without_authorization_policy():
    config = base_config()
    config["repo_feedback"] = github_repo_feedback()
    del config["repo_feedback"]["authorization"]

    assert_config_invalid(config, "missing authorization")


def test_kit_config_schema_rejects_repo_feedback_github_issues_with_wrong_organization():
    config = base_config()
    config["repo_feedback"] = github_repo_feedback()
    config["repo_feedback"]["authorization"]["organization"] = "Other-Organization"

    assert_config_invalid(config, "expected const Codeheart-Digital-Solutions")


def test_kit_config_schema_rejects_repo_feedback_github_issues_with_missing_or_false_authorization_policy():
    for field, value in [
        ("require_verified_membership", False),
        ("require_gh_cli", False),
        ("unavailable_behavior", "prompt"),
    ]:
        config = base_config()
        config["repo_feedback"] = github_repo_feedback()
        config["repo_feedback"]["authorization"][field] = value

        assert_config_invalid(config, "expected const")

        config = base_config()
        config["repo_feedback"] = github_repo_feedback()
        del config["repo_feedback"]["authorization"][field]

        assert_config_invalid(config, f"missing {field}")


def test_kit_config_schema_rejects_repo_feedback_disabled_without_reason():
    config = base_config()
    config["repo_feedback"] = {"mode": "disabled"}

    assert_config_invalid(config, "missing disabled_reason")


def test_kit_config_schema_rejects_invalid_repo_feedback_github_standardization_values():
    for field in ["labels", "issue_templates"]:
        config = base_config()
        config["repo_feedback"] = github_repo_feedback()
        config["repo_feedback"]["github_standardization"][field] = "required"

        assert_config_invalid(config, "invalid enum required")


def test_kit_config_schema_rejects_repo_feedback_mode_incompatible_fields():
    config = base_config()
    config["repo_feedback"] = github_repo_feedback()
    config["repo_feedback"]["disabled_reason"] = "other"

    assert_config_invalid(config, "unknown disabled_reason")

    config = base_config()
    config["repo_feedback"] = {
        "mode": "disabled",
        "disabled_reason": "other",
        "destination": {
            "type": "github_issues",
            "owner": "Codeheart-Digital-Solutions",
            "repo": "Codeheart-Operating-Kit",
        },
    }

    assert_config_invalid(config, "unknown destination")


def test_kit_config_schema_rejects_unknown_repo_feedback_mode():
    config = base_config()
    config["repo_feedback"] = {"mode": "local_draft_only"}

    assert_config_invalid(config, "invalid enum local_draft_only")


def test_kit_config_schema_accepts_old_config_without_local_machine_layer_path():
    config = base_config()
    del config["local_consumer_layer"]["local_machine_layer_path"]

    assert_config_valid(config)


def test_kit_config_schema_rejects_wrong_local_machine_layer_path():
    config = base_config()
    config["local_consumer_layer"]["local_machine_layer_path"] = ".codeheart/envs/"

    assert_config_invalid(config, "expected const .codeheart/local/")


def test_kit_config_schema_accepts_valid_member_portfolio_config():
    config = base_config()
    config["portfolio"] = {
        "role": "member",
        "member_repository_id": "Example-Automation",
        "coordination_home_path": "../Coordination-Home",
        "coordination_home_register_path": "docs/repo/plans/plan-register.md",
    }

    assert_config_valid(config)


def test_kit_config_schema_accepts_valid_coordination_home_portfolio_config():
    config = base_config()
    config["portfolio"] = {
        "role": "coordination-home",
        "coordination_home_register_path": "docs/repo/plans/plan-register.md",
    }

    assert_config_valid(config)


def test_kit_config_schema_rejects_empty_portfolio_config():
    config = base_config()
    config["portfolio"] = {}

    assert_config_invalid(config, "missing role")


def test_kit_config_schema_rejects_invalid_portfolio_role():
    config = base_config()
    config["portfolio"] = {"role": "standalone"}

    assert_config_invalid(config, "invalid enum standalone")


def test_kit_config_schema_rejects_removed_portfolio_fields():
    for removed_field in ["enabled", "member_register_path", "pending_sync_path"]:
        config = base_config()
        config["portfolio"] = {
            "role": "member",
            "member_repository_id": "Example-Automation",
            "coordination_home_path": "../Coordination-Home",
            "coordination_home_register_path": "docs/repo/plans/plan-register.md",
            removed_field: True,
        }

        assert_config_invalid(config, f"unknown {removed_field}")


def test_kit_config_schema_rejects_incomplete_member_portfolio_config():
    for missing_field in [
        "member_repository_id",
        "coordination_home_path",
        "coordination_home_register_path",
    ]:
        portfolio = {
            "role": "member",
            "member_repository_id": "Example-Automation",
            "coordination_home_path": "../Coordination-Home",
            "coordination_home_register_path": "docs/repo/plans/plan-register.md",
        }
        del portfolio[missing_field]
        config = base_config()
        config["portfolio"] = portfolio

        assert_config_invalid(config, f"missing {missing_field}")


def test_kit_config_schema_rejects_incomplete_coordination_home_portfolio_config():
    config = base_config()
    config["portfolio"] = {"role": "coordination-home"}

    assert_config_invalid(config, "missing coordination_home_register_path")
