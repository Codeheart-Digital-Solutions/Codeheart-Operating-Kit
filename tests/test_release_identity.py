"""Identity drift must fail before a release reaches expensive native tests."""
import importlib.util
from pathlib import Path
import shutil

import pytest

ROOT = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location("release_identity", ROOT / "scripts/validate-release-identity.py")
identity = importlib.util.module_from_spec(spec)
spec.loader.exec_module(identity)


def test_current_release_identity():
    assert identity.validate_identity(ROOT) == []


@pytest.fixture
def source(tmp_path):
    paths = ["manifest.yaml", "pyproject.toml", "internal/version/version.go",
             "src/codeheart_operating_kit/__init__.py", "install.sh", "install.ps1", "bootstrap.md"]
    paths += [p.relative_to(ROOT).as_posix() for p in (ROOT / "components").glob("*/component.yaml")]
    paths += [p.relative_to(ROOT).as_posix() for p in (ROOT / "profiles").glob("*.yaml")]
    paths += ["src/codeheart_operating_kit/resources/" + p for p in paths if p == "manifest.yaml" or p.startswith(("components/", "profiles/"))]
    for relative in paths:
        destination = tmp_path / relative
        destination.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(ROOT / relative, destination)
    return tmp_path


@pytest.mark.parametrize("relative,old,new,message", [
    ("pyproject.toml", 'version = "', 'version = "9', "project version"),
    ("internal/version/version.go", 'Version = "', 'Version = "9', "CLI version"),
    ("src/codeheart_operating_kit/__init__.py", '__version__ = "', '__version__ = "9', "Python package version"),
    ("install.sh", 'VERSION="', 'VERSION="9', "installer default version"),
    ("install.ps1", '$Version = "', '$Version = "9', "installer default version"),
    ("install.sh", "Default: ", "Default: 9", "installer help version"),
    ("install.ps1", "Default: ", "Default: 9", "installer help version"),
    ("profiles/standard.yaml", 'version: "', 'version: "9', "profile version"),
    ("components/agent-interface/component.yaml", 'version: "', 'version: "9', "version differs"),
    ("profiles/standard.yaml", "name: Standard", "name: Changed", "checksum differs"),
    ("manifest.yaml", "checksum_sha256: ", "checksum_sha256: bad", "checksum differs"),
    ("src/codeheart_operating_kit/resources/manifest.yaml", "schema_version: 1", "schema_version: 2", "mirror differs"),
])
def test_identity_drift_is_rejected(source, relative, old, new, message):
    path = source / relative
    original = path.read_text()
    assert old in original
    path.write_text(original.replace(old, new, 1))
    assert any(message in error for error in identity.validate_identity(source))


def test_missing_declaration_is_rejected(source):
    (source / "profiles/standard.yaml").unlink()
    assert any("source inventory" in error for error in identity.validate_identity(source))


def test_missing_mirror_is_rejected(source):
    (source / "src/codeheart_operating_kit/resources/profiles/standard.yaml").unlink()
    assert any("mirror differs" in error for error in identity.validate_identity(source))


@pytest.mark.parametrize("marker", [
    "Version: v", "releases/tag/v", "-macos-universal.zip", "-windows-x64.zip",
    "curl -fsSLO ", "Invoke-WebRequest -Uri ", "For `v",
])
@pytest.mark.parametrize("omit", [False, True])
def test_bootstrap_release_references_are_checked(source, marker, omit):
    path = source / "bootstrap.md"
    lines = path.read_text().splitlines(keepends=True)
    matches = [i for i, line in enumerate(lines) if marker in line]
    assert len(matches) == 1
    index = matches[0]
    version = identity.yaml.safe_load((source / "manifest.yaml").read_text())["version"]
    lines[index] = "" if omit else lines[index].replace(version, "9.9.9")
    path.write_text("".join(lines))
    assert any("bootstrap.md:" in error for error in identity.validate_identity(source))
