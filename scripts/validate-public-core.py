#!/usr/bin/env python3
from __future__ import annotations

import argparse
import re
from dataclasses import dataclass
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
EXCLUDED_PARTS = {
    ".git",
    ".pytest_cache",
    "__pycache__",
    "build",
    "dist",
    "node_modules",
    "codeheart_operating_kit.egg-info",
    "validator-invalid",
}
TEXT_SUFFIXES = {
    ".md",
    ".py",
    ".ps1",
    ".sh",
    ".json",
    ".yaml",
    ".yml",
    ".toml",
    ".txt",
}

LEGACY_ONBOARDING_EXAMPLE_PARTS = [
    ("Maria", "Automation"),
    ("Family", "Planning"),
    ("Home", "Research"),
    ("Bluebird", "Automation"),
    ("Finance", "Team", "Automation"),
    ("M365", "Office", "Automation"),
    ("Client", "Portal"),
    ("Booking", "App"),
    ("Storefront", "Relaunch"),
]
LEGACY_ONBOARDING_EXAMPLE_PATTERN = re.compile(
    r"\b(?:" + "|".join("-".join(parts) for parts in LEGACY_ONBOARDING_EXAMPLE_PARTS) + r")\b"
)


@dataclass(frozen=True)
class Rule:
    name: str
    pattern: re.Pattern[str]


RULES = [
    Rule("private platform repo name", re.compile("Codeheart-" + "AWS-Platform")),
    Rule("retired standards repo name", re.compile("Codeheart-" + "Standards")),
    Rule("AWS ARN", re.compile(r"\b" + "arn:" + r"aws:[A-Za-z0-9_:/+=,.@-]+")),
    Rule("AWS access key", re.compile(r"\b(?:AKIA|ASIA)[0-9A-Z]{16}\b")),
    Rule("private CodeArtifact account marker", re.compile(r"codeheart-[0-9]{12}")),
    Rule("private CodeArtifact host marker", re.compile(r"d\.codeartifact")),
    Rule("private key marker", re.compile(r"BEGIN (?:RSA |EC |OPENSSH |)PRIVATE KEY")),
    Rule("legacy onboarding example", LEGACY_ONBOARDING_EXAMPLE_PATTERN),
]


def should_scan(path: Path) -> bool:
    if any(part in EXCLUDED_PARTS for part in path.parts):
        return False
    return path.is_file() and path.suffix in TEXT_SUFFIXES


def iter_scan_paths(roots: list[Path]) -> list[Path]:
    paths: list[Path] = []
    for root in roots:
        if root.is_file():
            if root.suffix in TEXT_SUFFIXES:
                paths.append(root)
            continue
        for path in root.rglob("*"):
            if should_scan(path):
                paths.append(path)
    return sorted(paths)


def scan_paths(paths: list[Path]) -> list[str]:
    findings: list[str] = []
    for path in paths:
        text = path.read_text(encoding="utf-8", errors="ignore")
        for rule in RULES:
            if rule.pattern.search(text):
                try:
                    display = path.relative_to(ROOT)
                except ValueError:
                    display = path
                findings.append(f"{display}: blocked public-core pattern: {rule.name}")
    return findings


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate public-core hygiene.")
    parser.add_argument("paths", nargs="*", type=Path, default=[ROOT])
    args = parser.parse_args()

    findings = scan_paths(iter_scan_paths([path.resolve() for path in args.paths]))
    if findings:
        print("Public-core hygiene failed.")
        for finding in findings:
            print(f"- {finding}")
        return 1
    print("OK: public-core hygiene validates.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
