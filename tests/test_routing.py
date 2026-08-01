import json
import subprocess
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def test_producer_route_prefers_tracked_source_over_ignored_installation():
    agents = (ROOT / "AGENTS.md").read_text(encoding="utf-8")
    change_runbook = (ROOT / "docs/repo/runbooks/change-operating-kit.md").read_text(encoding="utf-8")

    assert "This repository is the Operating Kit producer" in agents
    assert "Never use the ignored consumer installation" in agents
    assert "components/" in agents and "internal/" in agents and "docs/repo/" in agents
    assert "Source\ncomponents, profiles, templates, schemas, Go packages" in change_runbook


def test_consumer_route_selects_repair_without_turning_sync_into_upgrade(tmp_path):
    binary = tmp_path / "codeheart-operating-kit"
    build = subprocess.run(
        ["go", "build", "-o", str(binary), "./cmd/codeheart-operating-kit"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert build.returncode == 0, build.stdout + build.stderr
    consumer = tmp_path / "consumer"
    init = subprocess.run(
        [str(binary), "init", str(consumer), "--project-name", "Routing-Probe"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert init.returncode == 0, init.stdout + init.stderr

    runbook = (
        consumer
        / ".codeheart/kit/docs/agent-interface/runbooks/maintain-operating-kit-installation.md"
    ).read_text(encoding="utf-8")
    agents = (consumer / "AGENTS.md").read_text(encoding="utf-8")

    assert "Compatible drift, partial installation, or lock-v1 migration" in runbook
    assert "repair --dry-run" in runbook
    assert "sync" in runbook and "preserves the installed kit version" in runbook
    assert "Only command that may change kit version" in runbook
    assert ".codeheart/kit/docs/agent-interface/runbooks/maintain-operating-kit-installation.md" in agents


def test_low_context_planning_and_portfolio_routes_are_installed(tmp_path):
    binary = tmp_path / "codeheart-operating-kit"
    build = subprocess.run(
        ["go", "build", "-o", str(binary), "./cmd/codeheart-operating-kit"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert build.returncode == 0, build.stdout + build.stderr
    consumer = tmp_path / "consumer"
    init = subprocess.run(
        [str(binary), "init", str(consumer), "--project-name", "Routing-Probe"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=False,
    )
    assert init.returncode == 0, init.stdout + init.stderr

    agents = (consumer / "AGENTS.md").read_text(encoding="utf-8")
    router = (
        consumer
        / ".codeheart/kit/docs/agent-interface/reference/operation-routing-and-dispatch.md"
    ).read_text(encoding="utf-8")
    planning = (
        consumer / ".codeheart/kit/docs/planning-workflows/README.md"
    ).read_text(encoding="utf-8")
    runbooks = consumer / ".codeheart/kit/docs/planning-workflows/runbooks"
    configure = (runbooks / "configure-portfolio-coordination.md").read_text(
        encoding="utf-8"
    )
    refresh = (runbooks / "refresh-portfolio-catalog.md").read_text(
        encoding="utf-8"
    )
    migrate = (runbooks / "migrate-plan-catalog.md").read_text(encoding="utf-8")
    draft = (runbooks / "draft-implementation-plan.md").read_text(encoding="utf-8")

    for route in [
        "configure-portfolio-coordination.md",
        "refresh-portfolio-catalog.md",
        "migrate-plan-catalog.md",
    ]:
        assert route in agents
        assert route in planning
    for route_id in [
        "planning.author-or-update",
        "portfolio.configure",
        "portfolio.refresh-and-analyze",
        "planning.migrate-catalog",
        "planning.activate-and-publish-checkpoint",
    ]:
        assert route_id in router

    assert "current cross-repo work" in router
    assert "portfolio scan --format json" in router
    assert "incomplete scan" in router
    assert "older cache is historical only" in router
    assert "member_repository_id` is required for both roles" in configure
    assert "previous complete cache was" in refresh and "historical" in refresh
    assert "changed or are owned by active branches remain skipped" in migrate
    assert "ask their branch owners" in migrate
    assert "do not ask for a second push approval" in draft.lower()


def test_activation_route_grants_only_plan_checkpoint_normal_push():
    draft = (
        ROOT / "components/planning-workflows/managed/runbooks/draft-implementation-plan.md"
    ).read_text(encoding="utf-8")
    execute = (
        ROOT / "components/planning-workflows/managed/runbooks/execute-implementation-plan.md"
    ).read_text(encoding="utf-8")
    router = (
        ROOT
        / "components/agent-interface/managed/reference/operation-routing-and-dispatch.md"
    ).read_text(encoding="utf-8")
    combined = "\n".join([draft, execute, router])

    assert "do not ask for a second push approval" in draft.lower()
    assert "plan-only branch/use/commit/normal-push" in router
    assert "without waiting for a PR or merge" in execute
    for excluded in [
        "unrelated dirty files",
        "PR creation",
        "merge",
        "release",
        "force-push",
        "branch deletion",
        "destructive Git",
        "ambiguous",
        "auth failure",
        "policy rejection",
        "rejected normal push",
    ]:
        assert excluded.lower() in combined.lower()
    assert "do not bypass the rejection" in execute.lower()


def test_new_coordination_runbooks_have_complete_execution_contracts():
    runbook_root = ROOT / "components/planning-workflows/managed/runbooks"
    for name in [
        "configure-portfolio-coordination.md",
        "refresh-portfolio-catalog.md",
        "migrate-plan-catalog.md",
    ]:
        text = (runbook_root / name).read_text(encoding="utf-8")
        for field in [
            "Audience:",
            "Intent:",
            "Success:",
            "Agent judgment boundary:",
            "Stop boundary:",
            "Source",
            "Inputs",
            "Preconditions",
            "Execution lane",
            "Approval boundary",
            "Stop Conditions",
            "Evidence And Validation",
            "Recovery",
        ]:
            assert field in text, (name, field)


def test_portfolio_configure_runbook_states_exact_scaffold_replacement_boundary():
    source = (
        ROOT
        / "components/planning-workflows/managed/runbooks/configure-portfolio-coordination.md"
    ).read_text(encoding="utf-8")
    packaged = (
        ROOT
        / "src/codeheart_operating_kit/resources/components/planning-workflows/managed/runbooks/configure-portfolio-coordination.md"
    ).read_text(encoding="utf-8")

    assert source == packaged
    assert "byte-exact untouched overlay placeholder" in source
    assert "README and repository-owned overlay bytes\nare preserved" in source
    assert "no other existing scaffold bytes\n   will change" in source
    assert "existing scaffolds are byte-preserved" not in source


def test_semantic_plan_reference_does_not_require_a_legacy_register_number():
    reference = (
        ROOT / "components/planning-workflows/managed/reference/plan-catalog-format.md"
    ).read_text(encoding="utf-8")
    router = (
        ROOT / "components/agent-interface/managed/reference/operation-routing-and-dispatch.md"
    ).read_text(encoding="utf-8")

    for field in ["title", "kind", "semantic ID", "family", "canonical path"]:
        assert field.lower() in reference.lower()
    assert "legacy_aliases" in reference
    assert "planning.author-or-update" in router
    assert "planning.migrate-catalog" in router


def test_activation_checkpoint_publishes_only_planning_paths_to_local_remote(tmp_path):
    repository = tmp_path / "repository"
    remote = tmp_path / "remote.git"
    repository.mkdir()

    def git(*args, cwd=repository):
        result = subprocess.run(
            ["git", *args], cwd=cwd, text=True, capture_output=True, check=False
        )
        assert result.returncode == 0, result.stdout + result.stderr
        return result.stdout.strip()

    git("init", "--bare", str(remote), cwd=tmp_path)
    git("init", "-b", "main")
    git("config", "user.email", "activation-probe@example.invalid")
    git("config", "user.name", "Activation Probe")
    plan = repository / "docs/repo/plans/example/example_implementation_doc.md"
    log = repository / "docs/repo/plans/example/example_execution_log.md"
    code = repository / "src/application.txt"
    plan.parent.mkdir(parents=True)
    code.parent.mkdir(parents=True)
    plan.write_text("Status: draft\n", encoding="utf-8")
    log.write_text("Status: pending\n", encoding="utf-8")
    code.write_text("baseline\n", encoding="utf-8")
    git("add", ".")
    git("commit", "-m", "Baseline")
    git("remote", "add", "origin", str(remote))
    git("push", "-u", "origin", "main")
    git("checkout", "-b", "codex/example-plan")

    plan.write_text("Status: active\n", encoding="utf-8")
    log.write_text("Status: active\n", encoding="utf-8")
    code.write_text("unrelated worktree change\n", encoding="utf-8")
    git("add", str(plan.relative_to(repository)), str(log.relative_to(repository)))
    git("commit", "-m", "Activate example plan")
    git("push", "-u", "origin", "codex/example-plan")

    assert git("show", "origin/codex/example-plan:" + str(plan.relative_to(repository))) == "Status: active"
    assert git("show", "origin/codex/example-plan:" + str(log.relative_to(repository))) == "Status: active"
    assert git("show", "origin/codex/example-plan:" + str(code.relative_to(repository))) == "baseline"
    assert "src/application.txt" in git("status", "--short")
    assert not git("for-each-ref", "--format=%(refname)", "refs/remotes/origin").endswith("refs/pull")


def test_adversarial_branch_fixture_covers_non_execution_and_safety_cases():
    fixture = json.loads(
        (
            ROOT
            / "tests/fixtures/portfolio/adversarial-branch/branch-tree.json"
        ).read_text(encoding="utf-8")
    )
    kinds = {entry["kind"] for entry in fixture["entries"]}
    assert kinds == {
        "executable-plan-adjacent-file",
        "git-hook",
        "malformed-metadata",
        "path-traversal-name",
        "oversized-metadata",
        "secret-like-placeholder",
    }
    assert all(entry["must_execute"] is False for entry in fixture["entries"])
    assert fixture["secret_policy"] == "placeholder-only-never-capture"
