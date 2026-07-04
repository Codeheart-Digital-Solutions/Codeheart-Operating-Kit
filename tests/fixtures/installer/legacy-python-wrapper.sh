#!/usr/bin/env sh
CODEHEART_OPERATING_KIT_CLI=1 PYTHONPATH="$HOME/.codeheart/operating-kit/lib${PYTHONPATH:+:$PYTHONPATH}" exec python3 -m codeheart_operating_kit.cli "$@"
