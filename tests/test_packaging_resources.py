from pathlib import Path
import os
import subprocess

import pytest
import yaml

import codeheart_operating_kit.components as components
import codeheart_operating_kit.manifest as manifest


ROOT = Path(__file__).resolve().parents[1]
PACKAGED_ROOT = ROOT / "src/codeheart_operating_kit/resources"


def assert_packaged_resource_matches(source: str) -> None:
    source_path = ROOT / source
    packaged_path = PACKAGED_ROOT / source
    assert packaged_path.exists(), source
    assert packaged_path.read_text(encoding="utf-8") == source_path.read_text(encoding="utf-8")


def test_packaged_resource_fallback(monkeypatch, tmp_path):
    monkeypatch.setattr(manifest, "SOURCE_ROOT", Path("/definitely/not/a/checkout"))
    monkeypatch.setattr(components, "kit_root", manifest.kit_root)

    state = components.write_default_state(
        tmp_path,
        project_name="Packaged-Automation",
        purpose="company-automation",
        selected_folder=str(tmp_path),
    )

    assert state["managed_paths"]
    assert (tmp_path / ".codeheart/kit/README.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/README.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/runbooks/submit-kit-feedback.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/runbooks/capture-repo-feedback.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/runbooks/enable-github-issues-feedback-intake.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/runbooks/handle-tooling-readiness.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/reference/kit-feedback-item-format.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/reference/repo-feedback-item-format.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/reference/operational-recipe-maturity.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/reference/runbook-authoring-standard.md").exists()
    assert (tmp_path / ".codeheart/kit/docs/structure-governance/reference/module-extension-state.md").exists()
    for relative in [
        ".codeheart/kit/docs/planning-workflows/runbooks/handle-routine-change.md",
        ".codeheart/kit/docs/agent-interface/reference/agent-task-coordination.md",
        ".codeheart/kit/docs/agent-interface/reference/codex-task-operations.md",
        ".codeheart/kit/docs/planning-workflows/reference/plan-catalog-format.md",
        ".codeheart/kit/docs/planning-workflows/reference/portfolio-coordination-format.md",
        ".codeheart/kit/docs/planning-workflows/runbooks/configure-portfolio-coordination.md",
        ".codeheart/kit/docs/planning-workflows/runbooks/refresh-portfolio-catalog.md",
        ".codeheart/kit/docs/planning-workflows/runbooks/migrate-plan-catalog.md",
    ]:
        assert (tmp_path / relative).exists(), relative
    assert (tmp_path / "AGENTS.md").exists()
    assert (tmp_path / "docs/repo/plans/plan-register.md").exists()
    assert (tmp_path / "docs/repo/portfolio/README.md").exists()
    assert (tmp_path / "docs/repo/portfolio/strategic-overlay.yaml").exists()
    assert not (tmp_path / "docs/repo/plans/coordination-sync-pending.md").exists()
    assert not (tmp_path / "docs/repo/state").exists()
    assert not (tmp_path / ".codeheart/local").exists()
    gitignore = (tmp_path / ".gitignore").read_text(encoding="utf-8")
    assert ".codeheart/user/feedback/" in gitignore
    assert ".codeheart/local/" in gitignore


def test_changed_source_and_packaged_resources_match():
    for source in [
        "components/agent-interface/managed/reference/runbook-to-script-promotion-standard.md",
        "components/agent-interface/managed/runbooks/handle-tooling-readiness.md",
        "components/planning-workflows/managed/runbooks/handle-routine-change.md",
        "components/agent-interface/managed/reference/agent-task-coordination.md",
        "components/agent-interface/managed/reference/codex-task-operations.md",
        "components/planning-workflows/component.yaml",
        "components/planning-workflows/managed/README.md",
        "components/planning-workflows/managed/reference/plan-register-format.md",
        "components/planning-workflows/managed/reference/plan-catalog-format.md",
        "components/planning-workflows/managed/reference/portfolio-coordination-format.md",
        "components/planning-workflows/managed/reference/planning-document-lifecycle.md",
        "components/planning-workflows/managed/runbooks/discovery-workflow.md",
        "components/planning-workflows/managed/runbooks/draft-implementation-plan.md",
        "components/planning-workflows/managed/runbooks/execute-implementation-plan.md",
        "components/planning-workflows/managed/runbooks/maintain-plan-register.md",
        "components/planning-workflows/managed/runbooks/configure-portfolio-coordination.md",
        "components/planning-workflows/managed/runbooks/refresh-portfolio-catalog.md",
        "components/planning-workflows/managed/runbooks/migrate-plan-catalog.md",
        "components/planning-workflows/managed/runbooks/review-planning-document.md",
        "components/planning-workflows/scaffolds/coordination-sync-pending.md",
        "components/planning-workflows/scaffolds/plan-register.md",
        "components/planning-workflows/scaffolds/portfolio-README.md",
        "components/planning-workflows/scaffolds/portfolio-strategic-overlay.yaml",
        "components/agent-interface/component.yaml",
        "components/agent-interface/managed/README.md",
        "components/agent-interface/managed/kit-readme.md",
        "components/agent-interface/managed/reference/local-extension-contract.md",
        "components/agent-interface/managed/reference/operation-routing-and-dispatch.md",
        "components/agent-interface/managed/reference/operational-recipe-maturity.md",
        "components/agent-interface/managed/reference/repo-feedback-item-format.md",
        "components/agent-interface/managed/reference/runbook-authoring-standard.md",
        "components/agent-interface/managed/reference/update-check-policy.md",
        "components/agent-interface/managed/reference/root-agents-md-contract.md",
        "components/agent-interface/managed/runbooks/conduct-first-run-onboarding.md",
        "components/agent-interface/managed/runbooks/capture-repo-feedback.md",
        "components/agent-interface/managed/runbooks/enable-github-issues-feedback-intake.md",
        "components/agent-interface/managed/runbooks/handle-tooling-readiness.md",
        "components/agent-interface/managed/runbooks/maintain-operating-kit-installation.md",
        "components/structure-governance/component.yaml",
        "components/structure-governance/managed/README.md",
        "components/structure-governance/managed/reference/documentation-structure.md",
        "components/structure-governance/managed/reference/managed-content-boundaries.md",
        "components/structure-governance/managed/reference/module-extension-state.md",
        "components/structure-governance/managed/runbooks/change-documentation-placement.md",
        "components/agent-memory/managed/README.md",
        "components/agent-memory/managed/reference/entry-format.md",
        "components/agent-memory/managed/runbooks/session-ledger-maintenance.md",
        "components/agent-memory/scaffolds/goal-register.md",
        "profiles/standard.yaml",
        "templates/agents/AGENTS.managed-block.md",
        "templates/consumer-docs/repo/README.md",
    ]:
        assert_packaged_resource_matches(source)


def test_packaged_migration_resources_expose_reviewed_v3_route_contract():
    manifest_text = (
        PACKAGED_ROOT / "components/planning-workflows/component.yaml"
    ).read_text(encoding="utf-8")
    migrate = (
        PACKAGED_ROOT
        / "components/planning-workflows/managed/runbooks/migrate-plan-catalog.md"
    ).read_text(encoding="utf-8")
    catalog = (
        PACKAGED_ROOT
        / "components/planning-workflows/managed/reference/plan-catalog-format.md"
    ).read_text(encoding="utf-8")
    router = (
        PACKAGED_ROOT
        / "components/agent-interface/managed/reference/operation-routing-and-dispatch.md"
    ).read_text(encoding="utf-8")
    refresh = (
        PACKAGED_ROOT
        / "components/planning-workflows/managed/runbooks/refresh-portfolio-catalog.md"
    ).read_text(encoding="utf-8")

    assert "route_id: planning.migrate-catalog" in manifest_text
    for text in [migrate, catalog, router]:
        assert "schema-v3" in text
        assert "deferred-active-owner" in text
        assert "canonical_ready" in text or "canonical readiness" in text
    assert "E -> L -> migrate -> catalog-activate -> A -> validate" in migrate
    assert "branch deletion is not a migration precondition" in migrate
    assert "branch deletion is neither the primary remedy nor granted authority" in router
    assert "schema-v3/discovery-v2" in refresh
    assert "mixed-grandfathered" in refresh
    assert "mixed_coverage_complete" in refresh and "canonical_ready" in refresh


def test_repo_feedback_runbook_prefers_configured_destination():
    text = (
        ROOT / "components/agent-interface/managed/runbooks/capture-repo-feedback.md"
    ).read_text(encoding="utf-8")
    configured_index = text.index("repo_feedback.destination.owner")
    fallback_index = text.index("resolve the GitHub remote")

    assert configured_index < fallback_index
    assert "repo_feedback.destination.repo" in text


def test_ubuntu_validation_is_semantic_only():
    workflow = (ROOT / ".github/workflows/validate.yml").read_text(encoding="utf-8")
    ubuntu = workflow.split("  ubuntu-semantic-validation:\n", 1)[1].split(
        "\n  windows-validation:", 1
    )[0]

    for required in [
        "runs-on: ubuntu-latest",
        "go test ./...",
        "BenchmarkDiscoveryV2Classifier100kPaths10kMarkdown",
        "tests/test_json_schemas.py",
        "tests/test_routing.py",
        "tests/test_packaging_resources.py",
        "tests/test_sync_check.py",
        "scripts/validate-json-schemas.py",
        "scripts/validate-markdown-headers.py",
        "scripts/validate-public-core.py",
    ]:
        assert required in ubuntu
    assert "build-release-assets.py" not in ubuntu


def validation_workflow():
    # BaseLoader keeps GitHub's "on" key a string (YAML 1.1 treats it as bool).
    return yaml.load((ROOT / ".github/workflows/validate.yml").read_text(), Loader=yaml.BaseLoader)


def test_validation_lanes_preserve_candidate_boundary():
    workflow = validation_workflow()
    assert workflow["on"]["push"] == {"branches": ["main"]}
    assert "pull_request" in workflow["on"]
    inputs = workflow["on"]["workflow_dispatch"]["inputs"]
    assert inputs["mode"]["default"] == "candidate"
    assert inputs["mode"]["options"] == ["candidate", "released-smoke"]
    jobs = workflow["jobs"]
    assert jobs["feedback"]["if"] == "github.event_name != 'workflow_dispatch'"
    feedback = str(jobs["feedback"]["steps"])
    for expensive in ["go test", "build-release-assets", "install.sh", "backward_compatibility", "Benchmark"]:
        assert expensive not in feedback
    for name in ["git-2-43-proof-validation", "macos-validation", "windows-validation", "ubuntu-semantic-validation"]:
        job = jobs[name]
        assert job["needs"] == "dispatch-inputs"
        assert job["if"] == "github.event_name == 'workflow_dispatch' && inputs.mode == 'candidate'"
        assert "/releases/download/" not in str(job["steps"])
    for name in ["macos-validation", "windows-validation", "ubuntu-semantic-validation"]:
        runs = [step.get("run", "") for step in jobs[name]["steps"]]
        broad = [run for run in runs if run.startswith("go test") and "-bench" not in run]
        builders = [run for run in runs if "scripts/build-release-assets.py" in run]
        assert len(broad) + len(builders) == 1
        if builders:
            assert 'run(["go", "test", "-timeout", "30m", "./..."])' in (ROOT / "scripts/build-release-assets.py").read_text()
    for name in ["macos-public-release", "windows-public-release"]:
        job = jobs[name]
        assert job["needs"] == "dispatch-inputs"
        assert job["if"] == "github.event_name == 'workflow_dispatch' && inputs.mode == 'released-smoke'"
        steps = str(job["steps"])
        assert "/releases/download/" in steps
        for source_work in ["go test", "pytest", "build-release-assets", "pyproject", "setup-python"]:
            assert source_work not in steps
    assert workflow["concurrency"]["cancel-in-progress"] == "${{ github.event_name != 'workflow_dispatch' }}"
    assert "github.run_id" in workflow["concurrency"]["group"]


@pytest.mark.parametrize("mode, tag, succeeds", [
    ("candidate", "", True), ("candidate", "v0.1.99", False),
    ("released-smoke", "", False), ("released-smoke", "v0.1.99", True),
    ("released-smoke", "main", False), ("unknown", "", False),
])
def test_dispatch_input_guard(mode, tag, succeeds):
    guard = validation_workflow()["jobs"]["dispatch-inputs"]["steps"][0]["run"]
    result = subprocess.run(["bash", "-c", guard], env={**os.environ, "MODE": mode, "RELEASE_VERSION": tag}, capture_output=True)
    assert (result.returncode == 0) == succeeds
