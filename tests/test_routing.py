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
