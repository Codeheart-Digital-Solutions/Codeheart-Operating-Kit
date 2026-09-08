#!/usr/bin/env python3
"""Candidate-only native content and current-release upgrade proof in isolated directories."""
from __future__ import annotations

import argparse
import os
import json
from pathlib import Path
import re
import subprocess
import time
import tomllib

import yaml

ROOT = Path(__file__).resolve().parents[1]


def run(*args: str | Path, succeeds: bool = True) -> str:
    result = subprocess.run([str(arg) for arg in args], capture_output=True, text=True)
    if (result.returncode == 0) != succeeds:
        raise RuntimeError(f"command {Path(args[0]).name} returned unexpected exit {result.returncode}: {result.stderr}")
    return result.stdout.strip()


def verify_materialized(root: Path, source: Path = ROOT) -> int:
    count = 0
    for declaration in (source / 'components').glob('*/component.yaml'):
        component = yaml.safe_load(declaration.read_text())['component']
        for entry in component.get('files', []):
            if entry['ownership'] != 'managed':
                continue
            target = root / entry['target']
            if target.is_symlink() or not target.is_file() or target.read_bytes() != (source / entry['source']).read_bytes():
                raise RuntimeError(f"materialized resource mismatch: {entry['target']}")
            count += 1
    if not count:
        raise RuntimeError('no managed resources verified')
    return count


def snapshot(root: Path) -> dict[str, bytes]:
    return {str(path.relative_to(root)): path.read_bytes() for path in root.rglob('*') if path.is_file()}


def unchanged(root: Path, before: dict[str, bytes], label: str) -> None:
    if snapshot(root) != before:
        raise RuntimeError(f'{label} changed installation state')


def wait_ready(binary: Path, target: Path, version: str, timeout: float = 60) -> None:
    """A replaced binary can precede transaction completion; retry only known pending state."""
    deadline = time.monotonic() + timeout
    while True:
        try:
            reported = run(binary, '--version')
        except (OSError, RuntimeError):
            reported = ''  # Executable may be temporarily unavailable during replacement.
        if reported == f'codeheart-operating-kit {version}':
            result = subprocess.run([str(binary), 'check', str(target), '--json'], capture_output=True,
                                    text=True, env={**os.environ, 'CODEHEART_OPERATING_KIT_CLI': '1'})
            state = json.loads(result.stdout)
            if result.returncode == 0 and state.get('ok') is True:
                return
            if state.get('state') != 'transaction-in-progress':
                raise RuntimeError(f"upgraded installation failed check: {state.get('state', 'unknown')}")
        if time.monotonic() >= deadline:
            raise RuntimeError('native upgrade transaction did not finish')
        time.sleep(.25)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    modes = parser.add_subparsers(dest='mode', required=True)
    wait = modes.add_parser('wait-ready', help='Wait for expected version and a completed native check')
    wait.add_argument('--binary', type=Path, required=True)
    wait.add_argument('--target', type=Path, required=True)
    wait.add_argument('--version', required=True)
    smoke = modes.add_parser('smoke', help='Run isolated fresh/current-release upgrade proof')
    smoke.add_argument('--fresh-target', type=Path, required=True)
    smoke.add_argument('--upgrade-version', required=True, help='Verified published vX.Y.Z input, distinct from source baseline')
    smoke.add_argument('--catalog', type=Path, required=True)
    smoke.add_argument('--work-dir', type=Path, required=True, help='New isolated directory; must not exist')
    args = parser.parse_args()
    if args.mode == 'wait-ready':
        wait_ready(args.binary, args.target, args.version)
        print('OK: expected version and completed installed state verified')
        return 0
    if not re.fullmatch(r'v\d+\.\d+\.\d+', args.upgrade_version):
        parser.error('upgrade-version must be an explicit published tag')
    version = tomllib.loads((ROOT / 'pyproject.toml').read_text())['project']['version']
    old_version = args.upgrade_version[1:]
    if tuple(map(int, version.split('.'))) <= tuple(map(int, old_version.split('.'))):
        parser.error('guidance candidate must advance the published upgrade version')
    fresh_count = verify_materialized(args.fresh_target)
    args.work_dir.mkdir(parents=True, exist_ok=False)
    home = args.work_dir / 'old-cli'
    consumer = args.work_dir / 'consumer'
    # Existing native installer verifies the public catalog -> archive -> payload -> binary chain.
    if os.name == 'nt':
        run('pwsh', '-NoProfile', '-File', ROOT / 'install.ps1', '-Version', old_version, '-InstallDir', home)
        binary = home / 'bin/codeheart-operating-kit.exe'
    else:
        run('bash', ROOT / 'install.sh', '--version', old_version, '--install-dir', home)
        binary = home / 'bin/codeheart-operating-kit'
    if run(binary, '--version') != f'codeheart-operating-kit {old_version}':
        raise RuntimeError('published old binary version mismatch')
    run(binary, 'init', consumer, '--project-name', 'Guidance-Smoke', '--json')
    authored = {
        'docs/repo/plans/authored.md': b'Preserve authored planning content.\n',
        '.codeheart/user/sentinel.txt': b'Preserve local user content.\n',
    }
    for relative, data in authored.items():
        target = consumer / relative
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_bytes(data)
    agents = consumer / 'AGENTS.md'
    with agents.open('ab') as stream:
        stream.write(b'\n# Consumer instructions\nPreserve this authored instruction.\n')
    config_before = (consumer / '.codeheart/kit.config.yaml').read_bytes()
    before = snapshot(consumer)
    binary_before = binary.read_bytes()
    common = [binary, 'upgrade', consumer, '--version', version, '--installed-binary', binary]
    run(*common, '--catalog', args.catalog.resolve(), '--dry-run', '--json')
    unchanged(consumer, before, 'dry-run')
    if binary.read_bytes() != binary_before:
        raise RuntimeError('dry-run replaced binary')
    run(*common, '--catalog', args.work_dir / 'missing-catalog.json', '--yes', '--json', succeeds=False)
    unchanged(consumer, before, 'failed verification')
    if binary.read_bytes() != binary_before:
        raise RuntimeError('failed verification replaced binary')
    run(*common, '--catalog', args.catalog.resolve(), '--yes', '--json')
    wait_ready(binary, consumer, version)
    upgraded_count = verify_materialized(consumer)
    for relative, data in authored.items():
        if (consumer / relative).read_bytes() != data:
            raise RuntimeError(f'upgrade changed authored file: {relative}')
    if (consumer / '.codeheart/kit.config.yaml').read_bytes() != config_before:
        raise RuntimeError('upgrade changed authored configuration')
    if not agents.read_bytes().endswith(b'\n# Consumer instructions\nPreserve this authored instruction.\n'):
        raise RuntimeError('upgrade changed authored AGENTS section')
    print(f'OK: {fresh_count} fresh and {upgraded_count} upgraded managed resources equal source; '
          f'{args.upgrade_version} -> {version} dry-run/failure/apply and authored preservation pass')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
