#!/usr/bin/env bash
set -euo pipefail
export LC_ALL=C
export LANG=C

VERSION="0.1.21"
INSTALL_DIR="${CODEHEART_OPERATING_KIT_HOME:-$HOME/.codeheart/operating-kit}"
ASSET_URL=""
ASSET_FILE=""
CHECKSUM=""
CHECKSUM_FILE=""
DEPRECATED_PYTHON=""
STAGING_DIR=""

usage() {
  cat <<'USAGE'
Install or repair codeheart-operating-kit for the current macOS user.

Options:
  --version VERSION          Release version to install. Default: 0.1.21
  --install-dir PATH         User-level install root. Default: $HOME/.codeheart/operating-kit
  --asset-url URL            Release asset URL. Defaults to the GitHub release asset.
  --asset-file PATH          Local release asset path for validation or offline repair.
  --checksum SHA256          Expected asset SHA-256.
  --checksum-file PATH       File containing the expected SHA-256.
  -h, --help                 Show this help.
USAGE
}

need_value() {
  if [[ $# -lt 2 || "$2" == -* ]]; then
    echo "Option $1 requires a value." >&2
    exit 2
  fi
}

file_url_path() {
  local path="${1#file://}"
  if [[ "$path" == localhost/* ]]; then
    path="/${path#localhost/}"
  fi
  printf '%b' "${path//%/\\x}"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      need_value "$@"
      VERSION="$2"
      shift 2
      ;;
    --install-dir)
      need_value "$@"
      INSTALL_DIR="$2"
      shift 2
      ;;
    --asset-url)
      need_value "$@"
      ASSET_URL="$2"
      shift 2
      ;;
    --asset-file)
      need_value "$@"
      ASSET_FILE="$2"
      shift 2
      ;;
    --checksum)
      need_value "$@"
      CHECKSUM="$2"
      shift 2
      ;;
    --checksum-file)
      need_value "$@"
      CHECKSUM_FILE="$2"
      shift 2
      ;;
    --python)
      need_value "$@"
      DEPRECATED_PYTHON="$2"
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

if [[ -n "$DEPRECATED_PYTHON" ]]; then
  echo "Warning: --python is deprecated and ignored; the Operating Kit installer uses a self-contained binary." >&2
fi

RELEASE_BASE_URL="https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v${VERSION}"
ASSET_NAME="codeheart-operating-kit-${VERSION}-macos-universal.zip"
if [[ -z "$ASSET_URL" ]]; then
  ASSET_URL="${RELEASE_BASE_URL}/${ASSET_NAME}"
fi

TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
  if [[ -n "$STAGING_DIR" ]]; then
    rm -rf "$STAGING_DIR"
  fi
}
trap cleanup EXIT

ASSET_PATH="$TMP_DIR/$ASSET_NAME"
if [[ -n "$ASSET_FILE" ]]; then
  cp "$ASSET_FILE" "$ASSET_PATH"
elif [[ "$ASSET_URL" == file://* ]]; then
  cp "$(file_url_path "$ASSET_URL")" "$ASSET_PATH"
else
  curl -fsSL "$ASSET_URL" -o "$ASSET_PATH"
fi

if [[ -z "$CHECKSUM" ]]; then
  if [[ -n "$CHECKSUM_FILE" ]]; then
    CHECKSUM="$(awk '{print $1}' "$CHECKSUM_FILE")"
  elif [[ "$ASSET_URL" == file://* ]]; then
    CHECKSUM_PATH="$(file_url_path "${ASSET_URL}.sha256")"
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

if ! command -v unzip >/dev/null 2>&1; then
  echo "unzip is required to extract the Operating Kit release pack." >&2
  exit 1
fi

EXTRACT_DIR="$TMP_DIR/extract"
mkdir -p "$EXTRACT_DIR"
if ! unzip -q "$ASSET_PATH" -d "$EXTRACT_DIR"; then
  echo "Release asset could not be extracted; installation stopped." >&2
  exit 1
fi

BINARY_PATH="$(find "$EXTRACT_DIR" -path "*/bin/codeheart-operating-kit" -type f | head -n 1)"
if [[ -z "$BINARY_PATH" ]]; then
  echo "Release asset did not contain bin/codeheart-operating-kit." >&2
  exit 1
fi

BIN_DIR="$INSTALL_DIR/bin"
LIB_DIR="$INSTALL_DIR/lib"
TARGET_BINARY="$BIN_DIR/codeheart-operating-kit"
mkdir -p "$BIN_DIR"

LEGACY_FOUND=0
if [[ -f "$TARGET_BINARY" ]] && head -n 20 "$TARGET_BINARY" | grep -Eq 'python|codeheart_operating_kit'; then
  LEGACY_FOUND=1
fi
if compgen -G "$LIB_DIR/codeheart_operating_kit*" >/dev/null 2>&1; then
  LEGACY_FOUND=1
fi
if [[ "$LEGACY_FOUND" == "1" ]]; then
  echo "Legacy Python install detected; installing the self-contained binary and preserving legacy files."
fi

STAGING_DIR="$(mktemp -d "$INSTALL_DIR/.staging.XXXXXX")"
STAGED_BINARY="$STAGING_DIR/codeheart-operating-kit"
cp "$BINARY_PATH" "$STAGED_BINARY"
chmod 0755 "$STAGED_BINARY"

SMOKE_LOG="$TMP_DIR/staged-version.log"
if ! "$STAGED_BINARY" --version >"$SMOKE_LOG" 2>&1; then
  cat "$SMOKE_LOG" >&2
  echo "Staged binary validation failed; previous runnable command preserved." >&2
  exit 1
fi

mv -f "$STAGED_BINARY" "$TARGET_BINARY"
chmod 0755 "$TARGET_BINARY"
rm -rf "$STAGING_DIR"
STAGING_DIR=""

echo "codeheart-operating-kit installed at $TARGET_BINARY"
if [[ "$LEGACY_FOUND" == "1" ]]; then
  echo "Legacy Python files were preserved under $INSTALL_DIR for later cleanup."
fi
case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *) echo "Add this folder to PATH to run it by name: $BIN_DIR" ;;
esac
echo "Next: codeheart-operating-kit onboard"
