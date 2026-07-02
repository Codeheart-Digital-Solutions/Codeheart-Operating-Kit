from pathlib import Path

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
    assert (tmp_path / "AGENTS.md").exists()
    assert (tmp_path / "docs/repo/plans/plan-register.md").exists()
    assert (tmp_path / "docs/repo/plans/coordination-sync-pending.md").exists()
    assert not (tmp_path / "docs/repo/state").exists()
    assert not (tmp_path / ".codeheart/local").exists()
    gitignore = (tmp_path / ".gitignore").read_text(encoding="utf-8")
    assert ".codeheart/user/feedback/" in gitignore
    assert ".codeheart/local/" in gitignore


def test_changed_source_and_packaged_resources_match():
    for source in [
        "components/planning-workflows/component.yaml",
        "components/planning-workflows/managed/README.md",
        "components/planning-workflows/managed/reference/plan-register-format.md",
        "components/planning-workflows/managed/reference/planning-document-lifecycle.md",
        "components/planning-workflows/managed/runbooks/discovery-workflow.md",
        "components/planning-workflows/managed/runbooks/draft-implementation-plan.md",
        "components/planning-workflows/managed/runbooks/execute-implementation-plan.md",
        "components/planning-workflows/managed/runbooks/maintain-plan-register.md",
        "components/planning-workflows/managed/runbooks/review-planning-document.md",
        "components/planning-workflows/scaffolds/coordination-sync-pending.md",
        "components/planning-workflows/scaffolds/plan-register.md",
        "components/agent-interface/component.yaml",
        "components/agent-interface/managed/README.md",
        "components/agent-interface/managed/kit-readme.md",
        "components/agent-interface/managed/reference/local-extension-contract.md",
        "components/agent-interface/managed/reference/operation-routing-and-dispatch.md",
        "components/agent-interface/managed/reference/operational-recipe-maturity.md",
        "components/agent-interface/managed/reference/repo-feedback-item-format.md",
        "components/agent-interface/managed/reference/runbook-authoring-standard.md",
        "components/agent-interface/managed/reference/root-agents-md-contract.md",
        "components/agent-interface/managed/runbooks/capture-repo-feedback.md",
        "components/agent-interface/managed/runbooks/enable-github-issues-feedback-intake.md",
        "components/agent-interface/managed/runbooks/handle-tooling-readiness.md",
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


def test_repo_feedback_runbook_prefers_configured_destination():
    text = (
        ROOT / "components/agent-interface/managed/runbooks/capture-repo-feedback.md"
    ).read_text(encoding="utf-8")
    configured_index = text.index("repo_feedback.destination.owner")
    fallback_index = text.index("resolve the GitHub remote")

    assert configured_index < fallback_index
    assert "repo_feedback.destination.repo" in text
