#!/usr/bin/env bash
set -euo pipefail
export LC_ALL=C
export LANG=C

VERSION="0.1.11"
INSTALL_DIR="${CODEHEART_OPERATING_KIT_HOME:-$HOME/.codeheart/operating-kit}"
ASSET_URL=""
ASSET_FILE=""
CHECKSUM=""
CHECKSUM_FILE=""
PYTHON_BIN="${PYTHON:-python3}"

usage() {
  cat <<'USAGE'
Install or repair codeheart-operating-kit for the current macOS user.

Options:
  --version VERSION          Release version to install. Default: 0.1.11
  --install-dir PATH         User-level install root. Default: $HOME/.codeheart/operating-kit
  --asset-url URL            Release asset URL. Defaults to the GitHub release asset.
  --asset-file PATH          Local release asset path for validation or offline repair.
  --checksum SHA256          Expected asset SHA-256.
  --checksum-file PATH       File containing the expected SHA-256.
  --python PATH              Python executable to use. Default: python3
  -h, --help                 Show this help.
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      VERSION="$2"
      shift 2
      ;;
    --install-dir)
      INSTALL_DIR="$2"
      shift 2
      ;;
    --asset-url)
      ASSET_URL="$2"
      shift 2
      ;;
    --asset-file)
      ASSET_FILE="$2"
      shift 2
      ;;
    --checksum)
      CHECKSUM="$2"
      shift 2
      ;;
    --checksum-file)
      CHECKSUM_FILE="$2"
      shift 2
      ;;
    --python)
      PYTHON_BIN="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

RELEASE_BASE_URL="https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v${VERSION}"
ASSET_NAME="codeheart-operating-kit-${VERSION}-macos.tar.gz"
if [[ -z "$ASSET_URL" ]]; then
  ASSET_URL="${RELEASE_BASE_URL}/${ASSET_NAME}"
fi

TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

ASSET_PATH="$TMP_DIR/$ASSET_NAME"
if [[ -n "$ASSET_FILE" ]]; then
  cp "$ASSET_FILE" "$ASSET_PATH"
elif [[ "$ASSET_URL" == file://* ]]; then
  cp "${ASSET_URL#file://}" "$ASSET_PATH"
else
  curl -fsSL "$ASSET_URL" -o "$ASSET_PATH"
fi

if [[ -z "$CHECKSUM" ]]; then
  if [[ -n "$CHECKSUM_FILE" ]]; then
    CHECKSUM="$(awk '{print $1}' "$CHECKSUM_FILE")"
  elif [[ "$ASSET_URL" == file://* ]]; then
    CHECKSUM_URL="${ASSET_URL}.sha256"
    CHECKSUM_PATH="${CHECKSUM_URL#file://}"
    CHECKSUM="$(awk '{print $1}' "$CHECKSUM_PATH")"
  else
    curl -fsSL "${ASSET_URL}.sha256" -o "$TMP_DIR/$ASSET_NAME.sha256"
    CHECKSUM="$(awk '{print $1}' "$TMP_DIR/$ASSET_NAME.sha256")"
  fi
fi

if [[ -z "$CHECKSUM" ]]; then
  echo "Expected checksum is required; installation stopped." >&2
  exit 1
fi

ACTUAL_CHECKSUM="$(shasum -a 256 "$ASSET_PATH" | awk '{print $1}')"
ACTUAL_LOWER="$(printf '%s' "$ACTUAL_CHECKSUM" | tr '[:upper:]' '[:lower:]')"
EXPECTED_LOWER="$(printf '%s' "$CHECKSUM" | tr '[:upper:]' '[:lower:]')"
if [[ "$ACTUAL_LOWER" != "$EXPECTED_LOWER" ]]; then
  echo "Checksum mismatch for $ASSET_NAME; installation stopped." >&2
  exit 1
fi

EXTRACT_DIR="$TMP_DIR/extract"
mkdir -p "$EXTRACT_DIR"
tar -xzf "$ASSET_PATH" -C "$EXTRACT_DIR"
WHEEL_PATH="$(find "$EXTRACT_DIR" -name 'codeheart_operating_kit-*.whl' -type f | head -n 1)"
if [[ -z "$WHEEL_PATH" ]]; then
  echo "Release asset did not contain a codeheart-operating-kit wheel." >&2
  exit 1
fi

BIN_DIR="$INSTALL_DIR/bin"
LIB_DIR="$INSTALL_DIR/lib"
mkdir -p "$BIN_DIR" "$LIB_DIR"
PIP_LOG="$TMP_DIR/pip-install.log"
if ! PIP_CONFIG_FILE=/dev/null \
  PIP_DISABLE_PIP_VERSION_CHECK=1 \
  PIP_NO_CACHE_DIR=1 \
  PIP_NO_INDEX=1 \
  PIP_NO_INPUT=1 \
  PYTHONNOUSERSITE=1 \
    "$PYTHON_BIN" -m pip install --no-index --no-deps --upgrade --target "$LIB_DIR" "$WHEEL_PATH" > "$PIP_LOG" 2>&1; then
  cat "$PIP_LOG" >&2
  exit 1
fi

WRAPPER="$BIN_DIR/codeheart-operating-kit"
cat > "$WRAPPER" <<EOF
#!/usr/bin/env sh
CODEHEART_OPERATING_KIT_CLI=1 PYTHONPATH="$LIB_DIR\${PYTHONPATH:+:\$PYTHONPATH}" exec "$PYTHON_BIN" -m codeheart_operating_kit.cli "\$@"
EOF
chmod +x "$WRAPPER"

echo "codeheart-operating-kit installed at $WRAPPER"
case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *) echo "Add this folder to PATH to run it by name: $BIN_DIR" ;;
esac
echo "Next: codeheart-operating-kit onboard"
