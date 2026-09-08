import importlib.util
from pathlib import Path
import sys

import pytest

ROOT = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location('guidance_lifecycle', ROOT / 'scripts/verify-guidance-lifecycle.py')
lifecycle = importlib.util.module_from_spec(spec)
spec.loader.exec_module(lifecycle)


def test_native_failure_is_not_hidden():
    with pytest.raises(RuntimeError, match='unexpected exit 7'):
        lifecycle.run(sys.executable, '-c', 'raise SystemExit(7)')
    lifecycle.run(sys.executable, '-c', 'raise SystemExit(7)', succeeds=False)
    with pytest.raises(RuntimeError, match='unexpected exit 0'):
        lifecycle.run(sys.executable, '-c', 'pass', succeeds=False)


def test_materialization_detects_corrupt_and_missing_bytes(tmp_path):
    source = tmp_path / 'source'
    component = source / 'components/example'
    component.mkdir(parents=True)
    (component / 'guide.md').write_bytes(b'expected\n')
    (component / 'component.yaml').write_text('component:\n  files:\n    - source: components/example/guide.md\n      target: installed.md\n      ownership: managed\n')
    target = tmp_path / 'target'
    target.mkdir()
    with pytest.raises(RuntimeError, match='mismatch'):
        lifecycle.verify_materialized(target, source)
    (target / 'installed.md').write_bytes(b'expected\n')
    assert lifecycle.verify_materialized(target, source) == 1
    (target / 'installed.md').write_bytes(b'corrupted\n')
    with pytest.raises(RuntimeError, match='mismatch'):
        lifecycle.verify_materialized(target, source)


def test_preservation_detects_added_and_changed_files(tmp_path):
    (tmp_path / 'owned').write_bytes(b'original')
    before = lifecycle.snapshot(tmp_path)
    lifecycle.unchanged(tmp_path, before, 'dry-run')
    (tmp_path / 'extra').write_bytes(b'unexpected')
    with pytest.raises(RuntimeError, match='changed installation'):
        lifecycle.unchanged(tmp_path, before, 'dry-run')


@pytest.mark.parametrize("state,exit_code,ok", [("current", 0, True), ("drifted", 1, False)])
def test_wait_ready_retries_only_transaction_pending(monkeypatch, tmp_path, state, exit_code, ok):
    import json
    from types import SimpleNamespace
    monkeypatch.setattr(lifecycle, "run", lambda *args: "codeheart-operating-kit 0.1.33")
    results = iter([
        SimpleNamespace(returncode=1, stdout=json.dumps({"ok": False, "state": "transaction-in-progress"})),
        SimpleNamespace(returncode=exit_code, stdout=json.dumps({"ok": ok, "state": state})),
    ])
    monkeypatch.setattr(lifecycle.subprocess, "run", lambda *args, **kwargs: next(results))
    sleeps = []
    monkeypatch.setattr(lifecycle.time, "sleep", lambda duration: sleeps.append(duration))
    if ok:
        lifecycle.wait_ready(tmp_path / "binary", tmp_path, "0.1.33")
    else:
        with pytest.raises(RuntimeError, match='failed check: .*drifted'):
            lifecycle.wait_ready(tmp_path / "binary", tmp_path, "0.1.33")
    assert sleeps == [.25]


def test_wait_ready_does_not_accept_version_alone(monkeypatch, tmp_path):
    import json
    from types import SimpleNamespace
    monkeypatch.setattr(lifecycle, "run", lambda *args: "codeheart-operating-kit 0.1.33")
    monkeypatch.setattr(lifecycle.subprocess, "run", lambda *args, **kwargs:
        SimpleNamespace(returncode=1, stdout=json.dumps({"ok": False, "state": "transaction-in-progress"})))
    with pytest.raises(RuntimeError, match="transaction did not finish"):
        lifecycle.wait_ready(tmp_path / "binary", tmp_path, "0.1.33", timeout=0)


@pytest.mark.parametrize("initial", ["partial", "drifted", "stale-cli"])
def test_wait_ready_observes_old_tree_before_reconcile(monkeypatch, tmp_path, initial):
    import json
    from types import SimpleNamespace
    monkeypatch.setattr(lifecycle, "run", lambda *args: "codeheart-operating-kit 0.1.33")
    states = iter([
        {"ok": False, "state": initial, "stale_cli": True},
        {"ok": False, "state": "transaction-in-progress"},
        {"ok": True, "state": "current", "stale_cli": False},
    ])
    def check(*args, **kwargs):
        state = next(states)
        return SimpleNamespace(returncode=0 if state["ok"] else 1, stdout=json.dumps(state))
    monkeypatch.setattr(lifecycle.subprocess, "run", check)
    monkeypatch.setattr(lifecycle.time, "sleep", lambda _: None)
    lifecycle.wait_ready(tmp_path / "binary", tmp_path, "0.1.33")


@pytest.mark.parametrize("stale", [True, False])
def test_wait_ready_partial_never_passes_and_preserves_diagnostics(monkeypatch, tmp_path, stale):
    import json
    from types import SimpleNamespace
    monkeypatch.setattr(lifecycle, "run", lambda *args: "codeheart-operating-kit 0.1.33")
    diagnostic = {"ok": False, "state": "partial", "stale_cli": stale,
                  "drift": [{"path": "missing-resource.md", "status": "missing"}]}
    monkeypatch.setattr(lifecycle.subprocess, "run", lambda *args, **kwargs:
        SimpleNamespace(returncode=1, stdout=json.dumps(diagnostic)))
    with pytest.raises(RuntimeError, match="missing-resource.md"):
        lifecycle.wait_ready(tmp_path / "binary", tmp_path, "0.1.33", timeout=0)
