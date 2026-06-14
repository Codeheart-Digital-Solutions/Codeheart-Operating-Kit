from __future__ import annotations

import json
from pathlib import Path

from ..components import copy_managed_files
from ..manifest import validate_release_manifest


def run(args) -> int:
    root = Path(args.path).expanduser()
    release = None
    if args.release_manifest:
        release = validate_release_manifest(Path(args.release_manifest))
    managed = copy_managed_files(root)
    result = {"synced_managed_paths": managed, "release_manifest": release is not None}
    if args.json:
        print(json.dumps(result, indent=2, sort_keys=True))
    else:
        print(f"Synced {len(managed)} managed files under .codeheart/kit/.")
    return 0
