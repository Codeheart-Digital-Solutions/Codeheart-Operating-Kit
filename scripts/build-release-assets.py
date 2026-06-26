#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import os
import shutil
import subprocess
import sys
import tarfile
import tempfile
import zipfile
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
PYPROJECT = ROOT / "pyproject.toml"


def package_version() -> str:
    for line in PYPROJECT.read_text(encoding="utf-8").splitlines():
        if line.startswith("version = "):
            return line.split("=", 1)[1].strip().strip('"')
    raise RuntimeError("pyproject.toml does not define project version")


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def sanitized_pip_env() -> dict[str, str]:
    env = {key: value for key, value in os.environ.items() if not key.startswith("PIP_")}
    env["PIP_CONFIG_FILE"] = os.devnull
    env["PIP_DISABLE_PIP_VERSION_CHECK"] = "1"
    env["PIP_NO_CACHE_DIR"] = "1"
    env["PIP_NO_INDEX"] = "1"
    env["PIP_NO_INPUT"] = "1"
    env["PYTHONNOUSERSITE"] = "1"
    return env


def build_wheel(version: str, work_dir: Path) -> Path:
    wheel_dir = work_dir / "wheelhouse"
    egg_base = work_dir / "egg-base"
    build_base = work_dir / "build-base"
    wheel_dir.mkdir()
    egg_base.mkdir()
    build_base.mkdir()
    command = [
        sys.executable,
        "-m",
        "pip",
        "wheel",
        "--no-build-isolation",
        "--no-deps",
        "--no-index",
        "--wheel-dir",
        str(wheel_dir),
        "--config-settings=--build-option=egg_info",
        f"--config-settings=--build-option=--egg-base={egg_base}",
        "--config-settings=--build-option=build",
        f"--config-settings=--build-option=--build-base={build_base}",
        str(ROOT),
    ]
    try:
        subprocess.run(
            command,
            check=True,
            env=sanitized_pip_env(),
            text=True,
            capture_output=True,
        )
    except subprocess.CalledProcessError as error:
        raise RuntimeError("wheel build failed while building local release asset payload") from error
    wheels = sorted(wheel_dir.glob("codeheart_operating_kit-*.whl"))
    if not wheels:
        raise RuntimeError("wheel build did not produce a codeheart-operating-kit wheel")
    expected = f"codeheart_operating_kit-{version}-"
    matching = [wheel for wheel in wheels if wheel.name.startswith(expected)]
    if not matching:
        raise RuntimeError(f"built wheel version does not match requested release version {version}")
    return matching[0]


def create_payload(wheel: Path, version: str, work_dir: Path) -> Path:
    payload = work_dir / f"codeheart-operating-kit-{version}"
    payload.mkdir()
    shutil.copy2(wheel, payload / wheel.name)
    (payload / "asset-manifest.json").write_text(
        json.dumps(
            {
                "schema_version": 1,
                "version": version,
                "wheel": wheel.name,
                "command": "codeheart-operating-kit",
            },
            indent=2,
            sort_keys=True,
        )
        + "\n",
        encoding="utf-8",
    )
    return payload


def create_tarball(payload: Path, output_dir: Path, version: str) -> Path:
    target = output_dir / f"codeheart-operating-kit-{version}-macos.tar.gz"
    with tarfile.open(target, "w:gz") as archive:
        archive.add(payload, arcname=payload.name)
    return target


def create_zip(payload: Path, output_dir: Path, version: str) -> Path:
    target = output_dir / f"codeheart-operating-kit-{version}-windows.zip"
    with zipfile.ZipFile(target, "w", compression=zipfile.ZIP_DEFLATED) as archive:
        for path in sorted(payload.rglob("*")):
            archive.write(path, path.relative_to(payload.parent))
    return target


def write_checksum(path: Path) -> Path:
    checksum_path = path.with_name(path.name + ".sha256")
    checksum_path.write_text(f"{sha256(path)}  {path.name}\n", encoding="utf-8")
    return checksum_path


def main() -> int:
    parser = argparse.ArgumentParser(description="Build Codeheart Operating Kit release assets.")
    parser.add_argument("--version", default="0.1.15")
    parser.add_argument("--output-dir", default="dist")
    args = parser.parse_args()
    actual_version = package_version()
    if args.version != actual_version:
        raise SystemExit(
            f"requested release version {args.version} does not match package version {actual_version}"
        )

    output_dir = Path(args.output_dir).resolve()
    output_dir.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory(prefix="codeheart-ok-assets-") as tmp:
        work_dir = Path(tmp)
        wheel = build_wheel(args.version, work_dir)
        payload = create_payload(wheel, args.version, work_dir)
        assets = [
            create_tarball(payload, output_dir, args.version),
            create_zip(payload, output_dir, args.version),
        ]
        checksums = [write_checksum(asset) for asset in assets]

    result = {
        "version": args.version,
        "assets": [asset.name for asset in assets],
        "checksums": [checksum.name for checksum in checksums],
    }
    print(json.dumps(result, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
