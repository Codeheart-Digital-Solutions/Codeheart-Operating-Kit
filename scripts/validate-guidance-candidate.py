#!/usr/bin/env python3
"""Read-only cumulative guidance eligibility check; semantic acceptance remains with owners."""
import argparse
import copy
import json
import re
import subprocess
import sys
from pathlib import PurePosixPath

import yaml

MIRROR = 'src/codeheart_operating_kit/resources/'
VERSION = r'[0-9]+\.[0-9]+\.[0-9]+'
LITERAL_IDENTITIES = {
    'pyproject.toml': [r'^version = "' + VERSION + r'"$'],
    'internal/version/version.go': [r'^var Version = "' + VERSION + r'"$'],
    'src/codeheart_operating_kit/__init__.py': [r'^__version__ = "' + VERSION + r'"$'],
    'install.sh': [r'^VERSION="' + VERSION + r'"$', r'^  --version VERSION          Release version to install\. Default: ' + VERSION + r'$'],
    'install.ps1': [r'^    \[string\]\$Version = "' + VERSION + r'",\r?$', r'^  -Version VERSION       Release version to install\. Default: ' + VERSION + r'\r?$'],
}


class Ineligible(ValueError):
    pass


def git(*args):
    result = subprocess.run(['git', *args], capture_output=True)
    if result.returncode:
        raise Ineligible('Git input is unavailable or invalid')
    return result.stdout


def content(ref, path):
    return git('show', f'{ref}:{path}')


def tree(ref):
    result = {}
    for row in git('ls-tree', '-r', '-z', ref).split(b'\0'):
        if row:
            metadata, path = row.split(b'\t', 1)
            mode, kind, oid = metadata.decode().split()
            result[path.decode()] = (mode, kind, oid)
    return result


def safe_path(path):
    return isinstance(path, str) and not path.startswith('/') and all(p not in ('', '.', '..') for p in path.split('/')) and '\\' not in path


class IdentityLoader(yaml.SafeLoader):
    pass


def unique_mapping(loader, node):
    result = {}
    for key_node, value_node in node.value:
        key = loader.construct_object(key_node)
        if key in result:
            raise Ineligible('Duplicate YAML identity key')
        result[key] = loader.construct_object(value_node)
    return result


IdentityLoader.add_constructor(yaml.resolver.BaseResolver.DEFAULT_MAPPING_TAG, unique_mapping)


def parsed(ref, path):
    raw = content(ref, path)
    if any(isinstance(token, (yaml.tokens.AnchorToken, yaml.tokens.AliasToken, yaml.tokens.TagToken)) for token in yaml.scan(raw)):
        raise Ineligible(f'Non-literal YAML identity syntax: {path}')
    value = yaml.load(raw, Loader=IdentityLoader)
    if not isinstance(value, dict):
        raise Ineligible(f'Invalid identity document: {path}')
    return value


def mask_identity(doc, kind):
    value = copy.deepcopy(doc)
    def mask(obj, key, pattern):
        if key not in obj or not isinstance(obj[key], str) or not re.fullmatch(pattern, obj[key]):
            raise Ineligible(f'Unrecognized {key} identity value')
        obj[key] = '<identity>'
    if kind == 'manifest':
        mask(value, 'version', VERSION)
        for group in ('components', 'profiles'):
            for entry in value[group]:
                mask(entry, 'version', VERSION)
                mask(entry, 'checksum_sha256', r'[a-f0-9]{64}')
                if group == 'profiles':
                    mask(entry, 'graph_sha256', r'[a-f0-9]{64}')
    else:
        mask(value[kind], 'version', VERSION)
    return value


def validate(baseline, candidate):
    if not re.fullmatch(r'[a-f0-9]{40}', baseline):
        raise Ineligible('baseline-ref must be an exact full commit SHA')
    baseline = git('rev-parse', '--verify', baseline + '^{commit}').decode().strip()
    candidate = git('rev-parse', '--verify', candidate + '^{commit}').decode().strip()
    git('merge-base', '--is-ancestor', baseline, candidate)
    old, new = tree(baseline), tree(candidate)
    changed = sorted(p for p in old.keys() | new.keys() if old.get(p) != new.get(p))
    for p in changed:
        if p not in new or new[p][:2] != ('100644', 'blob') or (p in old and old[p][:2] != new[p][:2]):
            raise Ineligible(f'Removal or unsupported file mode: {p}')
    # Rename detection is explicit: an eligible declaration cannot hide a moved path.
    statuses = git('diff', '--name-status', '-M', baseline, candidate).decode().splitlines()
    if any(line.startswith(('R', 'C')) for line in statuses):
        raise Ineligible('Renames/copies require broad validation')
    components = sorted(p for p in old if re.fullmatch(r'components/[^/]+/component.yaml', p))
    managed = {}
    identities = {'manifest.yaml': 'manifest'}
    identities.update({p: 'component' for p in components})
    identities.update({p: 'profile' for p in old if re.fullmatch(r'profiles/[^/]+\.yaml', p)})
    for path in components:
        before, after = parsed(baseline, path), parsed(candidate, path)
        bfiles, afiles = before['component']['files'], after['component']['files']
        if [entry for entry in afiles if entry in bfiles] != bfiles:
            raise Ineligible(f'Changed existing file declarations or order: {path}')
        known = [e for e in bfiles if e.get('ownership') == 'managed']
        for entry in afiles:
            if entry not in bfiles:
                if set(entry) != {'source', 'target', 'ownership'} or entry['ownership'] != 'managed':
                    raise Ineligible(f'Unknown managed declaration: {path}')
                if not all(safe_path(entry[k]) for k in ('source', 'target')):
                    raise Ineligible(f'Unsafe declaration: {path}')
                if not entry['source'].endswith('.md') or not entry['target'].endswith('.md'):
                    raise Ineligible(f'Non-Markdown addition: {path}')
                if not any(PurePosixPath(entry['source']).parent == PurePosixPath(e['source']).parent and PurePosixPath(entry['target']).parent == PurePosixPath(e['target']).parent for e in known):
                    raise Ineligible(f'Addition outside existing target contract: {path}')
                if entry['source'] in old or entry['source'] not in new:
                    raise Ineligible(f'Addition must declare a new source: {path}')
            if entry.get('ownership') == 'managed' and entry['source'].endswith('.md'):
                source, target = entry['source'], entry['target']
                if not safe_path(source) or not safe_path(target) or not source.startswith(str(PurePosixPath(path).parent) + '/managed/') or not target.startswith('.codeheart/kit/'):
                    raise Ineligible(f'Unrecognized managed placement: {path}')
                if source in managed or target in managed.values():
                    raise Ineligible(f'Duplicate managed declaration: {path}')
                managed[source] = target
        bmasked, amasked = mask_identity(before, 'component'), mask_identity(after, 'component')
        amasked['component']['files'] = bmasked['component']['files']
        if amasked != bmasked:
            raise Ineligible(f'Consequential component change: {path}')
    resources = []
    for p in changed:
        source = p.removeprefix(MIRROR)
        if p.startswith(MIRROR):
            if source not in new or content(candidate, p) != content(candidate, source):
                raise Ineligible(f'Resource mirror mismatch: {p}')
        elif p in managed or p in identities:
            if MIRROR + p not in new or content(candidate, MIRROR + p) != content(candidate, p):
                raise Ineligible(f'Missing or stale resource mirror: {p}')
        if source in managed:
            if not p.startswith(MIRROR):
                resources.append({'source': p, 'target': managed[p]})
            continue
        if source in identities:
            if source not in old:
                raise Ineligible(f'New identity surface: {source}')
            kind = identities[source]
            if kind != 'component' and mask_identity(parsed(baseline, source), kind) != mask_identity(parsed(candidate, source), kind):
                raise Ineligible(f'Consequential identity change: {source}')
            continue
        if p in LITERAL_IDENTITIES:
            a, b = content(baseline, p).decode(), content(candidate, p).decode()
            for pattern in LITERAL_IDENTITIES[p]:
                pattern = '(?m)' + pattern
                if len(re.findall(pattern, a)) != 1 or len(re.findall(pattern, b)) != 1:
                    raise Ineligible(f'Unrecognized literal version syntax: {p}')
                # Retain every byte other than the numeric release identity, including CRLF.
                replace = lambda match: re.sub(VERSION, '<version>', match.group())
                a, b = re.sub(pattern, replace, a), re.sub(pattern, replace, b)
            if a != b:
                raise Ineligible(f'Non-literal version edit: {p}')
            continue
        if p == 'bootstrap.md':
            old_version = parsed(baseline, 'manifest.yaml')['version']
            new_version = parsed(candidate, 'manifest.yaml')['version']
            a, b = content(baseline, p).decode(), content(candidate, p).decode()
            token = r'(?<![0-9.])' + re.escape(old_version) + r'(?![0-9.])'
            if not re.search(token, a) or re.sub(token, new_version, a) != b:
                raise Ineligible('Bootstrap permits only current release version token replacement')
            continue
        if p == 'release-notes.md' or (p.startswith('docs/repo/plans/') and p.endswith('.md')):
            continue
        raise Ineligible(f'Unknown or consequential path: {p}')
    return {'baseline': baseline, 'candidate': candidate, 'changed_paths': changed, 'changed_resources': resources}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--baseline-ref', required=True)
    parser.add_argument('--candidate-ref', default='HEAD')
    parser.add_argument('--json', action='store_true')
    args = parser.parse_args()
    try:
        result = validate(args.baseline_ref, args.candidate_ref)
    except (Ineligible, KeyError, TypeError, yaml.YAMLError, UnicodeError) as exc:
        print(f'Guidance candidate rejected: {exc}', file=sys.stderr)
        return 1
    if args.json:
        print(json.dumps(result, indent=2))
    else:
        print(f'Eligible cumulative diff: {len(result["changed_paths"])} paths, {len(result["changed_resources"])} managed resources.')
        print('Mechanical eligibility only; owner acceptance and semantic review remain required.')
    return 0


if __name__ == '__main__':
    sys.exit(main())
