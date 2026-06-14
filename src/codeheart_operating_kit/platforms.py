from __future__ import annotations

import platform


def current_platform() -> str:
    system = platform.system().lower()
    if system == "darwin":
        return "macos"
    if system == "windows":
        return "windows"
    return system


def is_supported_g1_platform() -> bool:
    return current_platform() in {"macos", "windows"}
