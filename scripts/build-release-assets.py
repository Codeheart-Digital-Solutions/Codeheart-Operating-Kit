#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import os
import shutil
import subprocess
import sys
import tempfile
import zipfile
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
PYPROJECT = ROOT / "pyproject.toml"
MAIN_PACKAGE = "./cmd/codeheart-operating-kit"


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


def relative_posix(path: Path, root: Path) -> str:
    return path.relative_to(root).as_posix()


def run(command: list[str], *, env: dict[str, str] | None = None) -> None:
    merged_env = os.environ.copy()
    if env:
        merged_env.update(env)
    try:
        subprocess.run(command, cwd=ROOT, env=merged_env, text=True, capture_output=True, check=True)
    except subprocess.CalledProcessError as error:
        message = error.stderr.strip() or error.stdout.strip() or "command failed"
        raise RuntimeError(f"{' '.join(command)} failed: {message}") from error


def build_binary(output: Path, version: str, goos: str, goarch: str) -> None:
    output.parent.mkdir(parents=True, exist_ok=True)
    run(
        [
            "go",
            "build",
            "-trimpath",
            "-ldflags",
            f"-s -w -X github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version.Version={version}",
            "-o",
            str(output),
            MAIN_PACKAGE,
        ],
        env={"CGO_ENABLED": "0", "GOOS": goos, "GOARCH": goarch},
    )


def build_macos_universal(work_dir: Path, version: str) -> Path:
    arm64 = work_dir / "build/macos-arm64/codeheart-operating-kit"
    amd64 = work_dir / "build/macos-amd64/codeheart-operating-kit"
    universal = work_dir / "build/macos-universal/codeheart-operating-kit"
    build_binary(arm64, version, "darwin", "arm64")
    build_binary(amd64, version, "darwin", "amd64")
    universal.parent.mkdir(parents=True, exist_ok=True)
    run(["lipo", "-create", "-output", str(universal), str(arm64), str(amd64)])
    run(["lipo", str(universal), "-verify_arch", "arm64", "x86_64"])
    universal.chmod(0o755)
    return universal


def build_windows_x64(work_dir: Path, version: str) -> Path:
    binary = work_dir / "build/windows-x64/codeheart-operating-kit.exe"
    build_binary(binary, version, "windows", "amd64")
    return binary


def copy_common_payload(payload: Path, version: str, platform: str, binary_name: str) -> None:
    for source in ["bootstrap.md", "install.sh", "install.ps1", "release-notes.md"]:
        shutil.copy2(ROOT / source, payload / source)
    install_doc = payload / "INSTALL.md"
    install_doc.write_text(
        "\n".join(
            [
                f"# Codeheart Operating Kit {version} {platform}",
                "",
                "This staged release pack contains a self-contained Operating Kit CLI binary.",
                "",
                f"Binary: `bin/{binary_name}`",
                "",
                "Installers verify the pack checksum before installing this binary.",
                "",
            ]
        ),
        encoding="utf-8",
    )


def write_payload_checksums(payload: Path) -> None:
    lines = []
    for path in sorted(payload.rglob("*")):
        if path.is_file() and path.name != "checksums.txt":
            lines.append(f"{sha256(path)}  {relative_posix(path, payload)}")
    (payload / "checksums.txt").write_text("\n".join(lines) + "\n", encoding="utf-8")


def create_pack(output_dir: Path, version: str, platform: str, binary: Path, binary_name: str) -> Path:
    pack_name = f"codeheart-operating-kit-{version}-{platform}.zip"
    with tempfile.TemporaryDirectory(prefix=f"codeheart-ok-{platform}-") as tmp:
        payload = Path(tmp) / f"codeheart-operating-kit-{version}-{platform}"
        bin_dir = payload / "bin"
        bin_dir.mkdir(parents=True)
        shutil.copy2(binary, bin_dir / binary_name)
        if not binary_name.endswith(".exe"):
            (bin_dir / binary_name).chmod(0o755)
        copy_common_payload(payload, version, platform, binary_name)
        manifest = {
            "schema_version": 1,
            "version": version,
            "platform": platform,
            "binary": f"bin/{binary_name}",
            "command": "codeheart-operating-kit",
        }
        (payload / "manifest.json").write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8")
        write_payload_checksums(payload)

        pack = output_dir / pack_name
        with zipfile.ZipFile(pack, "w", compression=zipfile.ZIP_DEFLATED) as archive:
            for path in sorted(payload.rglob("*")):
                archive.write(path, path.relative_to(payload.parent))
        return pack


def write_checksum(path: Path) -> Path:
    checksum_path = path.with_name(path.name + ".sha256")
    checksum_path.write_text(f"{sha256(path)}  {path.name}\n", encoding="utf-8")
    return checksum_path


def write_release_candidate_manifest(output_dir: Path, version: str, assets: list[Path]) -> Path:
    manifest = {
        "schema_version": 1,
        "version": version,
        "assets": [
            {
                "name": asset.name,
                "sha256": sha256(asset),
                "platform": platform_for_asset(asset),
            }
            for asset in assets
        ],
    }
    path = output_dir / f"release-candidate-manifest-{version}.json"
    path.write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    return path


def platform_for_asset(asset: Path) -> str:
    if "macos-universal" in asset.name:
        return "macos-universal"
    if "windows-x64" in asset.name:
        return "windows-x64"
    raise RuntimeError(f"cannot infer platform for {asset.name}")


def ensure_no_python_payload(pack: Path) -> None:
    forbidden = (".whl", ".dist-info/", "codeheart_operating_kit/")
    with zipfile.ZipFile(pack) as archive:
        names = archive.namelist()
    for name in names:
        if name.endswith(".whl") or any(part in name for part in forbidden[1:]):
            raise RuntimeError(f"{pack.name} contains forbidden Python payload {name}")


def main() -> int:
    parser = argparse.ArgumentParser(description="Build Codeheart Operating Kit release assets.")
    parser.add_argument("--version", default=None)
    parser.add_argument("--output-dir", default="dist")
    parser.add_argument(
        "--platform",
        choices=["all", "macos-universal", "windows-x64"],
        default="all",
        help="Platform pack to build. Default: all.",
    )
    args = parser.parse_args()
    actual_version = package_version()
    if args.version is None:
        args.version = actual_version
    if args.version != actual_version:
        raise SystemExit(
            f"requested release version {args.version} does not match package version {actual_version}"
        )

    output_dir = Path(args.output_dir).resolve()
    output_dir.mkdir(parents=True, exist_ok=True)
    run(["go", "test", "./..."])

    with tempfile.TemporaryDirectory(prefix="codeheart-ok-assets-") as tmp:
        work_dir = Path(tmp)
        assets = []
        if args.platform in {"all", "macos-universal"}:
            mac_binary = build_macos_universal(work_dir, args.version)
            assets.append(create_pack(output_dir, args.version, "macos-universal", mac_binary, "codeheart-operating-kit"))
        if args.platform in {"all", "windows-x64"}:
            windows_binary = build_windows_x64(work_dir, args.version)
            assets.append(create_pack(output_dir, args.version, "windows-x64", windows_binary, "codeheart-operating-kit.exe"))
        for asset in assets:
            ensure_no_python_payload(asset)
        checksums = [write_checksum(asset) for asset in assets]
        release_candidate = write_release_candidate_manifest(output_dir, args.version, assets)

    result = {
        "version": args.version,
        "assets": [asset.name for asset in assets],
        "checksums": [checksum.name for checksum in checksums],
        "release_candidate_manifest": release_candidate.name,
    }
    print(json.dumps(result, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
