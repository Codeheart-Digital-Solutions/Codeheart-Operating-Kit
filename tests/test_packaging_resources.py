from pathlib import Path

import codeheart_operating_kit.components as components
import codeheart_operating_kit.manifest as manifest


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
    assert (tmp_path / ".codeheart/kit/docs/agent-interface/reference/kit-feedback-item-format.md").exists()
    assert (tmp_path / "AGENTS.md").exists()
    gitignore = (tmp_path / ".gitignore").read_text(encoding="utf-8")
    assert ".codeheart/user/feedback/" in gitignore
