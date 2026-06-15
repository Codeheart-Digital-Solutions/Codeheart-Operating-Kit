from __future__ import annotations

import argparse
from pathlib import Path

from . import __version__
from .commands import check, init, inspect, onboard, sync, update_check


PURPOSE_CHOICES = ["private-automation", "company-automation", "software-product"]


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(prog="codeheart-operating-kit")
    parser.add_argument("--version", action="version", version=f"%(prog)s {__version__}")
    subparsers = parser.add_subparsers(dest="command", required=True)

    onboard_parser = subparsers.add_parser(
        "onboard",
        help="Guide first-run Operating Kit onboarding",
        description="Render the agent-guided first-run onboarding script and setup plan.",
        epilog=(
            "First-run user setup should run without --yes so Codex can ask for language, "
            "project name, target folder, setup-plan review, and write approval in chat. "
            "Non-interactive --yes writes require explicit --target and --project-name values."
        ),
    )
    onboard_parser.add_argument(
        "--target",
        help="Explicit folder to inspect or set up. Required with --yes.",
    )
    onboard_parser.add_argument(
        "--project-name",
        help="Explicit Codex project folder name. Required with --yes.",
    )
    onboard_parser.add_argument(
        "--purpose",
        choices=PURPOSE_CHOICES,
        help="Optional backward-compatible setup metadata; it does not select a different profile.",
    )
    onboard_parser.add_argument("--yes", action="store_true", help="Write setup files after explicit decisions are supplied")
    onboard_parser.add_argument("--json", action="store_true")
    onboard_parser.set_defaults(func=onboard.run)

    inspect_parser = subparsers.add_parser("inspect", help="Inspect a folder before setup")
    inspect_parser.add_argument("path", nargs="?", default=".")
    inspect_parser.add_argument("--json", action="store_true")
    inspect_parser.set_defaults(func=inspect.run)

    init_parser = subparsers.add_parser("init", help="Initialize Operating Kit in a folder")
    init_parser.add_argument("path", nargs="?", default=".")
    init_parser.add_argument("--project-name", default="Example-Automation")
    init_parser.add_argument(
        "--purpose",
        choices=PURPOSE_CHOICES,
        help="Optional backward-compatible setup metadata; omitted for new generic setups.",
    )
    init_parser.add_argument("--selected-folder")
    init_parser.add_argument("--json", action="store_true")
    init_parser.set_defaults(func=init.run)

    sync_parser = subparsers.add_parser("sync", help="Refresh managed Operating Kit files")
    sync_parser.add_argument("path", nargs="?", default=".")
    sync_parser.add_argument("--release-manifest", help="Optional release manifest fixture to validate before sync")
    sync_parser.add_argument("--json", action="store_true")
    sync_parser.set_defaults(func=sync.run)

    check_parser = subparsers.add_parser("check", help="Check installed Operating Kit state")
    check_parser.add_argument("path", nargs="?", default=".")
    check_parser.add_argument("--json", action="store_true")
    check_parser.set_defaults(func=check.run)

    update_parser = subparsers.add_parser("update-check", help="Check latest version metadata without applying updates")
    update_parser.add_argument("path", nargs="?", default=".")
    update_parser.add_argument("--latest-version")
    update_parser.add_argument(
        "--metadata-url",
        help="Latest release metadata URL; defaults to the public GitHub latest-release endpoint",
    )
    update_parser.add_argument("--now")
    update_parser.add_argument("--agent-notification", action="store_true")
    update_parser.add_argument("--json", action="store_true")
    update_parser.set_defaults(func=update_check.run)

    return parser


def main(argv: list[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)
    Path.cwd()
    return int(args.func(args))


if __name__ == "__main__":
    raise SystemExit(main())
