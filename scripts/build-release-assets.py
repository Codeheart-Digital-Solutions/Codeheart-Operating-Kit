#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import os
import shutil
import stat
import subprocess
import sys
import tempfile
import zipfile
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
PYPROJECT = ROOT / "pyproject.toml"
MAIN_PACKAGE = "./cmd/codeheart-operating-kit"
FIXED_ZIP_TIME = (1980, 1, 1, 0, 0, 0)


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


def canonical_json(value: dict) -> bytes:
    return (json.dumps(value, indent=2, sort_keys=True) + "\n").encode("utf-8")


def run(command: list[str], *, env: dict[str, str] | None = None) -> None:
    merged = os.environ.copy()
    merged.update({"LC_ALL": "C", "LANG": "C", "TZ": "UTC"})
    if env:
        merged.update(env)
    try:
        subprocess.run(command, cwd=ROOT, env=merged, text=True, capture_output=True, check=True)
    except subprocess.CalledProcessError as error:
        message = error.stderr.strip() or error.stdout.strip() or "command failed"
        raise RuntimeError(f"{' '.join(command)} failed: {message}") from error


def build_binary(output: Path, version: str, goos: str, goarch: str) -> None:
    output.parent.mkdir(parents=True, exist_ok=True)
    run(
        [
            "go", "build", "-trimpath", "-buildvcs=false",
            "-ldflags",
            f"-s -w -buildid= -X github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version.Version={version}",
            "-o", str(output), MAIN_PACKAGE,
        ],
        env={"CGO_ENABLED": "0", "GOOS": goos, "GOARCH": goarch},
    )


def build_platform_binary(work: Path, version: str, platform: str) -> tuple[Path, str]:
    if platform == "windows-x64":
        binary = work / "build/windows-x64/codeheart-operating-kit.exe"
        build_binary(binary, version, "windows", "amd64")
        return binary, "codeheart-operating-kit.exe"
    arm64 = work / "build/macos-arm64/codeheart-operating-kit"
    amd64 = work / "build/macos-amd64/codeheart-operating-kit"
    universal = work / "build/macos-universal/codeheart-operating-kit"
    build_binary(arm64, version, "darwin", "arm64")
    build_binary(amd64, version, "darwin", "amd64")
    universal.parent.mkdir(parents=True, exist_ok=True)
    run(["lipo", "-create", "-output", str(universal), str(arm64), str(amd64)])
    run(["lipo", str(universal), "-verify_arch", "arm64", "x86_64"])
    universal.chmod(0o755)
    return universal, "codeheart-operating-kit"


def assemble_payload(work: Path, version: str, platform: str, binary: Path, binary_name: str) -> Path:
    payload = work / f"codeheart-operating-kit-{version}-{platform}"
    (payload / "bin").mkdir(parents=True)
    shutil.copyfile(binary, payload / "bin" / binary_name)
    if not binary_name.endswith(".exe"):
        (payload / "bin" / binary_name).chmod(0o755)
    for source in ["bootstrap.md", "install.sh", "install.ps1", "release-notes.md"]:
        shutil.copyfile(ROOT / source, payload / source)
    shutil.copyfile(ROOT / "manifest.yaml", payload / "content-manifest.yaml")
    (payload / "INSTALL.md").write_text(
        f"# Codeheart Operating Kit {version} {platform}\n\n"
        "This release pack is verified from the external catalog through the staged binary.\n\n"
        f"Binary: `bin/{binary_name}`\n",
        encoding="utf-8",
    )
    return payload


def write_payload_identity(payload: Path, version: str, platform: str, binary_name: str) -> bytes:
    included = [
        item for item in sorted(payload.rglob("*"))
        if item.is_file() and item.name not in {"checksums.txt", "pack-manifest.json"}
    ]
    checksums = "".join(f"{sha256(item)}  {item.relative_to(payload).as_posix()}\n" for item in included)
    checksums_path = payload / "checksums.txt"
    checksums_path.write_text(checksums, encoding="utf-8")
    manifest = {
        "schema_version": 1,
        "version": version,
        "platform": platform,
        "command": "codeheart-operating-kit",
        "binary_path": f"bin/{binary_name}",
        "binary_sha256": sha256(payload / "bin" / binary_name),
        "content_manifest_path": "content-manifest.yaml",
        "content_manifest_sha256": sha256(payload / "content-manifest.yaml"),
        "payload_checksums_path": "checksums.txt",
        "payload_checksums_sha256": sha256(checksums_path),
    }
    data = canonical_json(manifest)
    (payload / "pack-manifest.json").write_bytes(data)
    return data


def deterministic_zip(payload: Path, output: Path) -> None:
    root_name = payload.name
    with zipfile.ZipFile(output, "w", compression=zipfile.ZIP_DEFLATED, compresslevel=9, strict_timestamps=True) as archive:
        for item in sorted(path for path in payload.rglob("*") if path.is_file()):
            relative = item.relative_to(payload).as_posix()
            info = zipfile.ZipInfo(f"{root_name}/{relative}", FIXED_ZIP_TIME)
            info.create_system = 3
            info.compress_type = zipfile.ZIP_DEFLATED
            mode = 0o755 if relative.startswith("bin/") and not relative.endswith(".exe") else 0o644
            info.external_attr = ((stat.S_IFREG | mode) & 0xFFFF) << 16
            info.flag_bits = 0
            archive.writestr(info, item.read_bytes(), compress_type=zipfile.ZIP_DEFLATED, compresslevel=9)


def build_pack(work: Path, version: str, platform: str) -> tuple[Path, str]:
    binary, binary_name = build_platform_binary(work, version, platform)
    payload = assemble_payload(work / "payload", version, platform, binary, binary_name)
    manifest_data = write_payload_identity(payload, version, platform, binary_name)
    pack = work / f"codeheart-operating-kit-{version}-{platform}.zip"
    deterministic_zip(payload, pack)
    return pack, hashlib.sha256(manifest_data).hexdigest()


def verify_pack_shape(pack: Path, version: str, platform: str) -> None:
    expected_prefix = f"codeheart-operating-kit-{version}-{platform}/"
    with zipfile.ZipFile(pack) as archive:
        names = archive.namelist()
        if names != sorted(names) or not names or any(not name.startswith(expected_prefix) for name in names):
            raise RuntimeError(f"{pack.name} has nondeterministic or unsafe entry order")
        for info in archive.infolist():
            if info.date_time != FIXED_ZIP_TIME or info.create_system != 3:
                raise RuntimeError(f"{pack.name} has non-normalized metadata for {info.filename}")
            if ((info.external_attr >> 16) & 0xF000) != stat.S_IFREG:
                raise RuntimeError(f"{pack.name} contains a non-regular entry: {info.filename}")
        required = {"pack-manifest.json", "checksums.txt", "content-manifest.yaml"}
        suffixes = {name.removeprefix(expected_prefix) for name in names}
        if not required <= suffixes:
            raise RuntimeError(f"{pack.name} is missing required identity files")
        if any(name.endswith(".whl") or ".dist-info/" in name or "codeheart_operating_kit/" in name for name in names):
            raise RuntimeError(f"{pack.name} contains a forbidden Python payload")
        manifest = json.loads(archive.read(expected_prefix + "pack-manifest.json"))
        if manifest["version"] != version or manifest["platform"] != platform or manifest["command"] != "codeheart-operating-kit":
            raise RuntimeError(f"{pack.name} has inconsistent pack identity")
        checksums_data = archive.read(expected_prefix + "checksums.txt")
        if hashlib.sha256(checksums_data).hexdigest() != manifest["payload_checksums_sha256"]:
            raise RuntimeError(f"{pack.name} has inconsistent payload checksum identity")
        expected_checksums = {
            relative: checksum
            for checksum, relative in (
                line.split("  ", 1) for line in checksums_data.decode("utf-8").splitlines()
            )
        }
        for relative, expected in expected_checksums.items():
            actual = hashlib.sha256(archive.read(expected_prefix + relative)).hexdigest()
            if actual != expected:
                raise RuntimeError(f"{pack.name} has a payload checksum mismatch for {relative}")
        if expected_checksums[manifest["binary_path"]] != manifest["binary_sha256"]:
            raise RuntimeError(f"{pack.name} has inconsistent binary identity")
        if expected_checksums[manifest["content_manifest_path"]] != manifest["content_manifest_sha256"]:
            raise RuntimeError(f"{pack.name} has inconsistent content identity")


def write_checksum(path: Path) -> Path:
    checksum_path = path.with_name(path.name + ".sha256")
    checksum_path.write_text(f"{sha256(path)}  {path.name}\n", encoding="utf-8")
    return checksum_path


def build_reproducible(output_dir: Path, version: str, platform: str) -> tuple[Path, str]:
    with tempfile.TemporaryDirectory(prefix=f"codeheart-ok-{platform}-first-") as first_tmp:
        first, first_manifest = build_pack(Path(first_tmp), version, platform)
        verify_pack_shape(first, version, platform)
        first_bytes = first.read_bytes()
    with tempfile.TemporaryDirectory(prefix=f"codeheart-ok-{platform}-second-") as second_tmp:
        second, second_manifest = build_pack(Path(second_tmp), version, platform)
        verify_pack_shape(second, version, platform)
        second_bytes = second.read_bytes()
    if first_bytes != second_bytes or first_manifest != second_manifest:
        raise RuntimeError(f"{platform} pack is not byte-reproducible")
    destination = output_dir / f"codeheart-operating-kit-{version}-{platform}.zip"
    destination.write_bytes(first_bytes)
    return destination, first_manifest


def write_catalog(output_dir: Path, version: str, assets: list[tuple[Path, str]], base_url: str) -> Path:
    catalog_assets = []
    for asset, manifest_sha in assets:
        platform = "macos-universal" if "macos-universal" in asset.name else "windows-x64"
        url = f"{base_url.rstrip('/')}/{asset.name}" if base_url else asset.name
        catalog_assets.append(
            {
                "name": asset.name,
                "version": version,
                "platform": platform,
                "url": url,
                "archive_sha256": sha256(asset),
                "pack_manifest_sha256": manifest_sha,
            }
        )
    catalog = {"schema_version": 1, "version": version, "assets": catalog_assets}
    path = output_dir / f"release-catalog-{version}.json"
    path.write_bytes(canonical_json(catalog))
    return path


def main() -> int:
    parser = argparse.ArgumentParser(description="Build deterministic Codeheart Operating Kit release packs.")
    parser.add_argument("--version")
    parser.add_argument("--output-dir", default="dist")
    parser.add_argument("--platform", choices=["all", "macos-universal", "windows-x64"], default="all")
    parser.add_argument("--base-url", default="", help="Optional publication base URL for catalog assets.")
    args = parser.parse_args()
    actual_version = package_version()
    version = args.version or actual_version
    if version != actual_version:
        raise SystemExit(f"requested release version {version} does not match package version {actual_version}")
    output_dir = Path(args.output_dir).resolve()
    output_dir.mkdir(parents=True, exist_ok=True)

    phases = ["source-validation"]
    run(["go", "test", "./..."])
    platforms = ["macos-universal", "windows-x64"] if args.platform == "all" else [args.platform]
    assets: list[tuple[Path, str]] = []
    for platform in platforms:
        phases.extend([f"{platform}:binary-build", f"{platform}:payload-assembly", f"{platform}:repeat-build-compare"])
        assets.append(build_reproducible(output_dir, version, platform))
    checksums = [write_checksum(asset) for asset, _ in assets]
    catalog = write_catalog(output_dir, version, assets, args.base_url)
    phases.append("catalog-emission")
    print(
        json.dumps(
            {
                "schema_version": 1,
                "version": version,
                "phases": phases,
                "assets": [asset.name for asset, _ in assets],
                "checksums": [path.name for path in checksums],
                "release_catalog": catalog.name,
                "reproducible": True,
            },
            indent=2,
            sort_keys=True,
        )
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
