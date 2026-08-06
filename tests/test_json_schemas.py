import json
import re
import subprocess
import sys
from copy import deepcopy
from pathlib import Path

from codeheart_operating_kit.manifest import load_yaml


ROOT = Path(__file__).resolve().parents[1]


def validate_instance(schema, instance):
    errors = []

    def resolve_ref(ref):
        current = schema
        for part in ref.removeprefix("#/").split("/"):
            current = current[part]
        return current

    def validate(subschema, value, location):
        if not isinstance(subschema, dict):
            return
        ref = subschema.get("$ref")
        if isinstance(ref, str) and ref.startswith("#/"):
            validate(resolve_ref(ref), value, location)
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
            elif isinstance(subschema.get("additionalProperties"), dict):
                for key in value:
                    if key not in properties:
                        validate(subschema["additionalProperties"], value[key], f"{location}.{key}")
            for key, property_schema in properties.items():
                if key in value:
                    validate(property_schema, value[key], f"{location}.{key}")
        elif expected_type == "integer" and not isinstance(value, int):
            errors.append(f"{location}: expected integer")
        elif expected_type == "boolean" and not isinstance(value, bool):
            errors.append(f"{location}: expected boolean")
        elif expected_type == "array":
            if not isinstance(value, list):
                errors.append(f"{location}: expected array")
                return
            if len(value) < subschema.get("minItems", 0):
                errors.append(f"{location}: too few items")
            if subschema.get("uniqueItems") and len({json.dumps(item, sort_keys=True) for item in value}) != len(value):
                errors.append(f"{location}: duplicate item")
            for index, item in enumerate(value):
                validate(subschema.get("items", {}), item, f"{location}[{index}]")
        elif expected_type == "string":
            if not isinstance(value, str):
                errors.append(f"{location}: expected string")
            elif len(value) < subschema.get("minLength", 0):
                errors.append(f"{location}: too short")

        if isinstance(value, int) and value < subschema.get("minimum", value):
            errors.append(f"{location}: below minimum")
        if isinstance(value, str) and "pattern" in subschema and re.search(subschema["pattern"], value) is None:
            errors.append(f"{location}: pattern mismatch")

        if "const" in subschema and value != subschema["const"]:
            errors.append(f"{location}: expected const {subschema['const']}")
        if "enum" in subschema and value not in subschema["enum"]:
            errors.append(f"{location}: invalid enum {value}")

        for clause in subschema.get("allOf", []):
            condition = clause.get("if")
            then = clause.get("then")
            otherwise = clause.get("else")
            if condition:
                if condition_matches(condition, value):
                    if then:
                        validate(then, value, location)
                elif otherwise:
                    validate(otherwise, value, location)
            else:
                validate(clause, value, location)

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


def durable_schema(name):
    return json.loads((ROOT / "schemas" / name).read_text(encoding="utf-8"))


def test_state_foundation_schemas_are_versioned_and_declared():
    component = json.loads((ROOT / "schemas/component.schema.json").read_text(encoding="utf-8"))
    profile = json.loads((ROOT / "schemas/profile.schema.json").read_text(encoding="utf-8"))
    lock_v1 = json.loads((ROOT / "schemas/kit-lock-v1.schema.json").read_text(encoding="utf-8"))
    lock_v2 = json.loads((ROOT / "schemas/kit-lock.schema.json").read_text(encoding="utf-8"))

    file_properties = component["$defs"]["file"]["properties"]
    assert {"presence_policy", "update_strategy", "removal_strategy", "route_id"}.issubset(
        file_properties
    )
    state_defaults = profile["properties"]["profile"]["properties"]["state_defaults"]
    assert set(state_defaults["required"]) == {
        "managed",
        "scaffold",
        "template",
        "generated-surface",
        "local-user",
        "local-machine",
    }
    assert lock_v1["properties"]["schema_version"]["const"] == 1
    assert lock_v2["properties"]["schema_version"]["const"] == 2
    assert {"state_generation", "release_provenance", "last_operation"}.issubset(
        lock_v2["required"]
    )


def test_profile_schema_keeps_legacy_update_check_metadata_optional_and_compatible():
    schema = json.loads((ROOT / "schemas/profile.schema.json").read_text(encoding="utf-8"))
    profile = load_yaml(ROOT / "profiles/standard.yaml")

    assert "update_check" not in schema["properties"]["profile"]["required"]
    assert "update_check" not in profile["profile"]

    profile["profile"]["update_check"] = {
        "cadence_days": 7,
        "current_result_agent_message": "silent",
        "update_available_agent_message": "prompt-user-before-apply",
    }
    assert validate_instance(schema, profile) == []


def test_release_identity_schemas_are_acyclic_and_exact():
    content = json.loads((ROOT / "schemas/content-manifest.schema.json").read_text(encoding="utf-8"))
    catalog = json.loads((ROOT / "schemas/release-catalog.schema.json").read_text(encoding="utf-8"))
    pack = json.loads((ROOT / "schemas/pack-manifest.schema.json").read_text(encoding="utf-8"))
    embedded = load_yaml(ROOT / "manifest.yaml")

    assert "assets" not in embedded
    assert "released_at" not in embedded
    assert set(content["required"]) == {
        "schema_version", "version", "compatibility", "components", "profiles", "consumer_impact"
    }
    asset_required = set(catalog["properties"]["assets"]["items"]["required"])
    assert {"archive_sha256", "pack_manifest_sha256", "url", "platform", "version"} <= asset_required
    assert "codeheart-operating-kit" in catalog["properties"]["assets"]["items"]["properties"]["name"]["pattern"]
    assert {"binary_sha256", "content_manifest_sha256", "payload_checksums_sha256"} <= set(pack["required"])


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


def test_kit_config_schema_reserves_discovery_sources_for_coordination_homes():
    portfolio_schema = kit_config_schema()["properties"]["portfolio"]
    member_rule = portfolio_schema["allOf"][0]["then"]["allOf"][0]["then"]

    assert member_rule["properties"]["discovery"] is False


def test_plan_catalog_schemas_define_strict_versioned_contracts():
    expected = {
        "plan-metadata.schema.json": {"plan"},
        "plan-catalog.schema.json": {
            "schema_version", "discovery_version", "policy_digest", "candidate_set_digest",
            "coordination_home_id", "started_at", "completed_at", "complete", "members",
            "observations", "candidates", "errors", "metrics",
        },
        "plan-catalog-v1.schema.json": {
            "schema_version", "coordination_home_id", "started_at", "completed_at", "complete",
            "members", "observations", "candidates", "errors", "metrics",
        },
        "plan-migration-ledger.schema.json": {
            "schema_version", "repository_id", "discovery_version", "target_catalog_mode",
            "policy_digest", "candidate_set_digest", "inventory_revision", "reviewed_at", "records",
        },
        "plan-migration-ledger-v1.schema.json": {
            "schema_version", "repository_id", "inventory_revision", "reviewed_at", "records",
        },
        "portfolio-local-sources.schema.json": {"schema_version", "sources"},
        "portfolio-strategic-overlay.schema.json": {
            "schema_version", "coordination_home_id", "families", "themes", "relations",
            "priorities", "analyses",
        },
    }
    for filename, required in expected.items():
        schema = json.loads((ROOT / "schemas" / filename).read_text(encoding="utf-8"))
        assert schema["$schema"] == "https://json-schema.org/draft/2020-12/schema"
        assert schema["additionalProperties"] is False
        if filename == "plan-metadata.schema.json":
            assert set(schema["required"]) == required
            assert set(schema["properties"]["plan"]["required"]) == {
                "schema_version", "id", "kind", "purpose", "first_cataloged",
                "catalog_metadata_updated",
            }
        else:
            assert set(schema["required"]) == required


def test_every_new_durable_schema_accepts_a_positive_instance_and_rejects_a_negative_one():
    commit = "a" * 40
    digest = "b" * 64
    timestamp = "2026-07-31T10:00:00Z"
    cases = {
        "plan-metadata.schema.json": (
            {
                "plan": {
                    "schema_version": 1,
                    "id": "example.discovery.catalog",
                    "kind": "discovery",
                    "purpose": "Define the catalog",
                    "first_cataloged": timestamp,
                    "catalog_metadata_updated": timestamp,
                    "products": ["operating-kit"],
                }
            },
            lambda value: value["plan"].update({"unknown": True}),
            "unknown unknown",
        ),
        "plan-catalog.schema.json": (
            {
                "schema_version": 2,
                "discovery_version": 2,
                "policy_digest": digest,
                "candidate_set_digest": "c" * 64,
                "coordination_home_id": "example-home",
                "started_at": timestamp,
                "completed_at": timestamp,
                "complete": True,
                "members": [],
                "observations": [],
                "candidates": [],
                "errors": [],
                "metrics": {
                    "duration_ms": 0,
                    "source_count": 0,
                    "member_count": 0,
                    "candidate_count": 0,
                    "observation_count": 0,
                    "stale_count": 0,
                    "api_call_count": 0,
                    "max_concurrency": 1,
                },
            },
            lambda value: value["metrics"].update({"max_concurrency": 0}),
            "below minimum",
        ),
        "plan-migration-ledger.schema.json": (
            {
                "schema_version": 2,
                "repository_id": "example",
                "discovery_version": 2,
                "target_catalog_mode": "canonical",
                "policy_digest": digest,
                "candidate_set_digest": "c" * 64,
                "inventory_revision": commit,
                "reviewed_at": timestamp,
                "records": [
                    {
                        "current_path": "docs/repo/plans/catalog.md",
                        "target_path": "docs/repo/plans/catalog_discovery_doc.md",
                        "source_revision": commit,
                        "source_sha256": digest,
                        "target_precondition": {"state": "absent"},
                        "ownership_disposition": "owned",
                        "decision": {
                            "id": "example.discovery.catalog",
                            "kind": "discovery",
                            "purpose": "Define the catalog",
                            "first_cataloged": timestamp,
                            "catalog_metadata_updated": timestamp,
                        },
                        "confidence": "high",
                        "ambiguity": [],
                        "evidence": ["Reviewed the canonical plan"],
                        "legacy_aliases": ["PR01"],
                        "conflicts": [],
                        "deferred": False,
                        "branch_owner": "none",
                    }
                ],
            },
            lambda value: value.pop("inventory_revision"),
            "missing inventory_revision",
        ),
        "portfolio-local-sources.schema.json": (
            load_yaml(ROOT / "tests/fixtures/portfolio/local-sources.yaml"),
            lambda value: value["sources"][0].update({"kind": "worktree"}),
            "expected const local-git-root",
        ),
        "portfolio-strategic-overlay.schema.json": (
            load_yaml(ROOT / "tests/fixtures/portfolio/strategic-overlay.yaml"),
            lambda value: value.pop("analyses"),
            "missing analyses",
        ),
    }

    for name, (positive, make_negative, expected_error) in cases.items():
        schema = durable_schema(name)
        assert validate_instance(schema, positive) == [], name
        negative = deepcopy(positive)
        make_negative(negative)
        errors = validate_instance(schema, negative)
        assert any(expected_error in error for error in errors), (name, errors)


def test_historical_plan_catalog_and_ledger_schemas_remain_v1_only():
    catalog_v1 = durable_schema("plan-catalog-v1.schema.json")
    ledger_v1 = durable_schema("plan-migration-ledger-v1.schema.json")
    assert catalog_v1["properties"]["schema_version"]["const"] == 1
    assert ledger_v1["properties"]["schema_version"]["const"] == 1
    assert durable_schema("plan-catalog.schema.json")["properties"]["schema_version"]["const"] == 2
    assert durable_schema("plan-migration-ledger.schema.json")["properties"]["schema_version"]["const"] == 2


def test_kit_config_schema_exposes_exclusions_only_discovery_v2_settings():
    config = base_config()
    config["component_settings"]["planning-workflows"] = {
        "plan_catalog_discovery_version": 2,
        "plan_catalog_ownership": {"excluded_roots": ["vendor/docs/"]},
    }
    assert_config_valid(config)
    config["component_settings"]["planning-workflows"]["plan_catalog_ownership"][
        "owned_roots"
    ] = ["docs/"]
    assert_config_invalid(config, "unknown owned_roots")

    compatible = base_config()
    compatible["component_settings"]["planning-workflows"] = {
        "unrelated_extension_setting": {"enabled": True}
    }
    assert_config_valid(compatible)


def test_kit_config_schema_accepts_v2_member_and_home_fixtures():
    for name in ["kit-config-v2-member.yaml", "kit-config-v2-home.yaml"]:
        config = load_yaml(ROOT / "tests" / "fixtures" / "portfolio" / name)
        assert_config_valid(config)


def test_kit_config_schema_accepts_v1_path_only_coordination_home_fixture():
    config = load_yaml(
        ROOT / "tests" / "fixtures" / "portfolio" / "kit-config-v1-home.yaml"
    )

    assert_config_valid(config)


def test_kit_config_schema_validates_mixed_cutover_revision_shape():
    config = base_config()
    config["component_settings"]["planning-workflows"] = {
        "plan_catalog_mode": "mixed",
        "plan_catalog_cutover_revision": "a" * 40,
    }
    assert_config_valid(config)
    config["component_settings"]["planning-workflows"][
        "plan_catalog_cutover_revision"
    ] = "not-a-commit"
    assert_config_invalid(config, "pattern mismatch")


def test_kit_config_schema_rejects_v2_home_without_member_repository_id():
    config = load_yaml(
        ROOT / "tests" / "fixtures" / "portfolio" /
        "kit-config-v2-home-missing-repository-id.yaml"
    )
    assert_config_invalid(config, "missing member_repository_id")


def test_kit_config_schema_preserves_portfolio_v1_and_defaults_catalog_mode_to_legacy():
    member = base_config()
    member["portfolio"] = {
        "role": "member",
        "member_repository_id": "legacy-member",
        "coordination_home_path": "../home",
        "coordination_home_register_path": "docs/repo/plans/plan-register.md",
    }
    home = base_config()
    home["portfolio"] = {
        "role": "coordination-home",
        "coordination_home_register_path": "docs/repo/plans/plan-register.md",
    }
    assert_config_valid(member)
    assert_config_valid(home)
    assert "planning-workflows" not in member["component_settings"]


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
