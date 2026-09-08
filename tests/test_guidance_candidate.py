"""Adversarial committed-tree fixtures for the cumulative guidance guard."""
import importlib.util
import json
import subprocess
from pathlib import Path

import pytest
import yaml

ROOT = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location('guidance_guard', ROOT / 'scripts/validate-guidance-candidate.py')
guard = importlib.util.module_from_spec(spec)
spec.loader.exec_module(guard)
MIRROR = guard.MIRROR
SOURCE = 'components/example/managed/reference/example.md'
TARGET = '.codeheart/kit/docs/example/reference/example.md'
COMPONENT = 'components/example/component.yaml'


def run(*args):
    return subprocess.check_output(['git', *args]).decode().strip()


def write(path, text, mirror=False):
    target = Path(path)
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_text(text)
    if mirror:
        write(MIRROR + path, text)


def commit():
    run('add', '.')
    run('commit', '-qm', 'fixture')
    return run('rev-parse', 'HEAD')


@pytest.fixture
def repo(tmp_path, monkeypatch):
    monkeypatch.chdir(tmp_path)
    run('init', '-q')
    run('config', 'user.email', 'fixture@example.invalid')
    run('config', 'user.name', 'Fixture')
    component = {'schema_version': 1, 'component': {'id': 'example', 'version': '1.2.3', 'ownership_modes': ['managed'], 'files': [{'source': SOURCE, 'target': TARGET, 'ownership': 'managed'}]}}
    write(COMPONENT, yaml.safe_dump(component), True)
    write(SOURCE, '# Original\n', True)
    write('manifest.yaml', yaml.safe_dump({'version': '1.2.3', 'components': [], 'profiles': [], 'compatibility': {'commands': ['init']}}), True)
    write('profiles/standard.yaml', yaml.safe_dump({'profile': {'version': '1.2.3', 'selected_components': ['example']}}), True)
    write('pyproject.toml', '[project]\nversion = "1.2.3"\ndependencies = []\n')
    write('internal/version/version.go', 'package version\nvar Version = "1.2.3"\n')
    write('runtime.go', 'package runtime\n')
    return commit()


def test_managed_change_and_known_bookkeeping(repo):
    write(SOURCE, '# Revised\n', True)
    write('docs/repo/plans/example/example_execution_log.md', 'Evidence\n')
    write('release-notes.md', 'Reviewed guidance\n')
    result = guard.validate(repo, commit())
    assert result['changed_resources'] == [{'source': SOURCE, 'target': TARGET}]
    assert result['baseline'] == repo


def test_literal_identity_changes(repo):
    for path in (COMPONENT, 'manifest.yaml', 'profiles/standard.yaml', 'pyproject.toml', 'internal/version/version.go'):
        write(path, Path(path).read_text().replace('1.2.3', '1.2.4'), path.endswith('.yaml'))
    assert guard.validate(repo, commit())['changed_resources'] == []


def test_managed_addition_under_existing_contract(repo):
    data = yaml.safe_load(Path(COMPONENT).read_text())
    source = SOURCE.replace('example.md', 'added.md')
    target = TARGET.replace('example.md', 'added.md')
    data['component']['files'].append({'source': source, 'target': target, 'ownership': 'managed'})
    write(COMPONENT, yaml.safe_dump(data), True)
    write(source, '# Added\n', True)
    assert guard.validate(repo, commit())['changed_resources'] == [{'source': source, 'target': target}]


@pytest.mark.parametrize('path', ['runtime.go', '.github/workflows/validate.yml', 'docs/repo/runbooks/release-operating-kit.md', 'templates/agents/AGENTS.managed-block.md', 'components/new/component.yaml', 'schemas/unknown.md', 'docs/repo/plans/executable.py'])
def test_unknown_or_consequential_paths_reject(repo, path):
    write(path, 'changed\n')
    with pytest.raises(guard.Ineligible):
        guard.validate(repo, commit())


@pytest.mark.parametrize('kind', ['ownership', 'selection', 'compatibility', 'executable_version', 'dependency'])
def test_disguised_identity_changes_reject(repo, kind):
    if kind == 'ownership':
        write(COMPONENT, Path(COMPONENT).read_text().replace('managed', 'scaffold'), True)
    elif kind == 'selection':
        write('profiles/standard.yaml', Path('profiles/standard.yaml').read_text().replace('example', 'other'), True)
    elif kind == 'compatibility':
        write('manifest.yaml', Path('manifest.yaml').read_text().replace('init', 'other'), True)
    elif kind == 'executable_version':
        write('internal/version/version.go', 'package version\nvar Version = calculate()\n')
    else:
        write('pyproject.toml', '[project]\nversion = "1.2.4"\ndependencies = ["surprise"]\n')
    with pytest.raises(guard.Ineligible):
        guard.validate(repo, commit())


@pytest.mark.parametrize('kind', ['missing_mirror', 'mirror_only', 'delete', 'rename', 'executable', 'symlink'])
def test_integrity_and_file_modes_reject(repo, kind):
    if kind == 'missing_mirror':
        write(SOURCE, '# Changed\n')
    elif kind == 'mirror_only':
        write(MIRROR + SOURCE, '# Changed\n')
    elif kind == 'delete':
        Path(SOURCE).unlink()
    elif kind == 'rename':
        Path(SOURCE).rename(Path(SOURCE).with_name('renamed.md'))
    elif kind == 'executable':
        run('update-index', '--chmod=+x', SOURCE)
    else:
        oid = subprocess.check_output(['git', 'hash-object', '-w', '--stdin'], input=b'example.md').decode().strip()
        run('update-index', '--cacheinfo', f'120000,{oid},{SOURCE}')
    if kind in ('executable', 'symlink'):
        # Build Git modes directly so this fixture also works on Windows hosts.
        run('commit', '-qm', 'mode fixture')
        candidate = run('rev-parse', 'HEAD')
    else:
        candidate = commit()
    with pytest.raises(guard.Ineligible):
        guard.validate(repo, candidate)


@pytest.mark.parametrize('kind', ['outside_target', 'outside_source', 'extra_key', 'duplicate', 'undeclared'])
def test_invalid_additions_reject(repo, kind):
    source = SOURCE.replace('example.md', 'added.md')
    target = TARGET.replace('example.md', 'added.md')
    data = yaml.safe_load(Path(COMPONENT).read_text())
    entry = {'source': source, 'target': target, 'ownership': 'managed'}
    if kind == 'outside_target':
        entry['target'] = '.codeheart/kit/new-area/added.md'
    elif kind == 'outside_source':
        entry['source'] = 'templates/added.md'
    elif kind == 'extra_key':
        entry['install_when'] = 'absent'
    elif kind == 'duplicate':
        entry = data['component']['files'][0].copy()
    if kind != 'undeclared':
        data['component']['files'].append(entry)
        write(COMPONENT, yaml.safe_dump(data), True)
    write(source, '# Added\n', True)
    with pytest.raises(guard.Ineligible):
        guard.validate(repo, commit())


def test_cumulative_runtime_change_cannot_hide_behind_recent_guidance(repo):
    write('runtime.go', 'package changed\n')
    commit()
    write(SOURCE, '# Revised\n', True)
    with pytest.raises(guard.Ineligible):
        guard.validate(repo, commit())


def test_exact_ancestor_baseline_required(repo):
    with pytest.raises(guard.Ineligible):
        guard.validate('HEAD', 'HEAD')
    write(SOURCE, '# Revised\n', True)
    later = commit()
    with pytest.raises(guard.Ineligible):
        guard.validate(later, repo)


def test_json_cli_is_clean_machine_output(repo):
    write(SOURCE, '# Revised\n', True)
    candidate = commit()
    import sys
    result = subprocess.run([sys.executable, str(ROOT / 'scripts/validate-guidance-candidate.py'), '--baseline-ref', repo, '--json'], text=True, capture_output=True)
    assert result.returncode == 0, result.stderr
    assert json.loads(result.stdout)['candidate'] == candidate


@pytest.mark.parametrize('path,original', [
    ('src/codeheart_operating_kit/__init__.py', '__version__ = "1.2.3"\n'),
    ('install.sh', 'VERSION="1.2.3"\n  --version VERSION          Release version to install. Default: 1.2.3\necho unchanged\n'),
    ('install.ps1', '    [string]$Version = "1.2.3",\r\n  -Version VERSION       Release version to install. Default: 1.2.3\r\nWrite-Output unchanged\r\n'),
])
def test_known_installer_and_python_literals_only(repo, path, original):
    write(path, original)
    baseline = commit()
    write(path, original.replace('1.2.3', '1.2.4'))
    guard.validate(baseline, commit())
    write(path, original.replace('1.2.3', '1.2.4') + '\n# new behavior\n')
    with pytest.raises(guard.Ineligible):
        guard.validate(baseline, commit())


@pytest.mark.parametrize('extra', ['version: 1.2.4\n', 'unexpected: &anchor foo\n'])
def test_ambiguous_yaml_identity_syntax_rejects(repo, extra):
    write('manifest.yaml', Path('manifest.yaml').read_text() + extra, True)
    with pytest.raises(guard.Ineligible):
        guard.validate(repo, commit())


def test_bootstrap_only_current_release_tokens(repo):
    original = 'Version: v1.2.3\nURL: https://example.invalid/v1.2.3/install.sh\nOld v0.0.9\n'
    write('bootstrap.md', original)
    baseline = commit()
    write('bootstrap.md', original.replace('1.2.3', '1.2.4'))
    write('manifest.yaml', Path('manifest.yaml').read_text().replace('1.2.3', '1.2.4'), True)
    guard.validate(baseline, commit())
    write('bootstrap.md', original.replace('1.2.3', '1.2.4').replace('example.invalid', 'different.invalid'))
    with pytest.raises(guard.Ineligible):
        guard.validate(baseline, commit())
