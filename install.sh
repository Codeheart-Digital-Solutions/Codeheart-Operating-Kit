#!/usr/bin/env bash
set -euo pipefail
export LC_ALL=C
export LANG=C

VERSION="0.1.21"
INSTALL_DIR="${CODEHEART_OPERATING_KIT_HOME:-$HOME/.codeheart/operating-kit}"
CATALOG_URL=""
CATALOG_FILE=""
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
  --catalog-url URL          External release catalog URL.
  --catalog-file PATH        Local external release catalog.
  --asset-url URL            Optional release asset URL override.
  --asset-file PATH          Local release asset for offline validation.
  --checksum SHA256          Optional checksum that must agree with the catalog.
  --checksum-file PATH       Optional checksum file that must agree with the catalog.
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
  local value="${1#file://}"
  if [[ "$value" == localhost/* ]]; then
    value="/${value#localhost/}"
  fi
  printf '%b' "${value//%/\\x}"
}

json_get() {
  /usr/bin/plutil -extract "$2" raw -o - "$1" 2>/dev/null
}

sha256_file() {
  /usr/bin/shasum -a 256 "$1" | /usr/bin/awk '{print $1}'
}

lower() {
  printf '%s' "$1" | /usr/bin/tr '[:upper:]' '[:lower:]'
}

fail() {
  echo "$1" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version) need_value "$@"; VERSION="$2"; shift 2 ;;
    --install-dir) need_value "$@"; INSTALL_DIR="$2"; shift 2 ;;
    --catalog-url) need_value "$@"; CATALOG_URL="$2"; shift 2 ;;
    --catalog-file) need_value "$@"; CATALOG_FILE="$2"; shift 2 ;;
    --asset-url) need_value "$@"; ASSET_URL="$2"; shift 2 ;;
    --asset-file) need_value "$@"; ASSET_FILE="$2"; shift 2 ;;
    --checksum) need_value "$@"; CHECKSUM="$2"; shift 2 ;;
    --checksum-file) need_value "$@"; CHECKSUM_FILE="$2"; shift 2 ;;
    --python) need_value "$@"; DEPRECATED_PYTHON="$2"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

if [[ -n "$DEPRECATED_PYTHON" ]]; then
  echo "Warning: --python is deprecated and ignored; the Operating Kit installer uses a self-contained binary." >&2
fi
[[ -x /usr/bin/plutil ]] || fail "plutil is required to validate the external release catalog."

RELEASE_BASE_URL="https://github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/releases/download/v${VERSION}"
CATALOG_NAME="release-catalog-${VERSION}.json"
if [[ -z "$CATALOG_URL" && -z "$CATALOG_FILE" ]]; then
  if [[ -n "$ASSET_FILE" && -f "$(/usr/bin/dirname "$ASSET_FILE")/$CATALOG_NAME" ]]; then
    CATALOG_FILE="$(/usr/bin/dirname "$ASSET_FILE")/$CATALOG_NAME"
  else
    CATALOG_URL="${RELEASE_BASE_URL}/${CATALOG_NAME}"
  fi
fi
if [[ -n "$CATALOG_URL" && "$CATALOG_URL" != https://* && "$CATALOG_URL" != file://* ]]; then
  fail "Release catalogs must use HTTPS or a local file URL."
fi

TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
  if [[ -n "$STAGING_DIR" ]]; then
    rm -rf "$STAGING_DIR"
  fi
}
trap cleanup EXIT

CATALOG_PATH="$TMP_DIR/$CATALOG_NAME"
if [[ -n "$CATALOG_FILE" ]]; then
  cp "$CATALOG_FILE" "$CATALOG_PATH"
elif [[ "$CATALOG_URL" == file://* ]]; then
  cp "$(file_url_path "$CATALOG_URL")" "$CATALOG_PATH"
else
  curl -fsSL "$CATALOG_URL" -o "$CATALOG_PATH"
fi

[[ "$(json_get "$CATALOG_PATH" schema_version)" == "1" ]] || fail "Release catalog schema is unsupported."
[[ "$(json_get "$CATALOG_PATH" version)" == "$VERSION" ]] || fail "Release catalog version does not match the requested version."

MATCHES=0
INDEX=0
CATALOG_ASSET_URL=""
CATALOG_ARCHIVE_SHA=""
CATALOG_PACK_SHA=""
while ASSET_VERSION="$(json_get "$CATALOG_PATH" "assets.$INDEX.version")"; do
  ASSET_PLATFORM="$(json_get "$CATALOG_PATH" "assets.$INDEX.platform")"
  if [[ "$ASSET_VERSION" == "$VERSION" && "$ASSET_PLATFORM" == "macos-universal" ]]; then
    MATCHES=$((MATCHES + 1))
    CATALOG_ASSET_URL="$(json_get "$CATALOG_PATH" "assets.$INDEX.url")"
    CATALOG_ARCHIVE_SHA="$(json_get "$CATALOG_PATH" "assets.$INDEX.archive_sha256")"
    CATALOG_PACK_SHA="$(json_get "$CATALOG_PATH" "assets.$INDEX.pack_manifest_sha256")"
    ASSET_NAME="$(json_get "$CATALOG_PATH" "assets.$INDEX.name")"
  fi
  INDEX=$((INDEX + 1))
done
[[ "$MATCHES" == "1" ]] || fail "Release catalog must contain exactly one macos-universal asset for $VERSION."
[[ "$ASSET_NAME" == "codeheart-operating-kit-${VERSION}-macos-universal.zip" ]] || fail "Release catalog asset name does not match the requested version and platform."
[[ "$CATALOG_ARCHIVE_SHA" =~ ^[A-Fa-f0-9]{64}$ && "$CATALOG_PACK_SHA" =~ ^[A-Fa-f0-9]{64}$ ]] || fail "Release catalog contains an invalid digest."

if [[ -n "$CHECKSUM_FILE" ]]; then
  CHECKSUM="$(/usr/bin/awk '{print $1}' "$CHECKSUM_FILE")"
fi
if [[ -n "$CHECKSUM" && "$(lower "$CHECKSUM")" != "$(lower "$CATALOG_ARCHIVE_SHA")" ]]; then
  fail "Checksum mismatch: explicit checksum disagrees with the external release catalog."
fi

if [[ -z "$ASSET_URL" ]]; then
  ASSET_URL="$CATALOG_ASSET_URL"
  if [[ "$ASSET_URL" != *://* ]]; then
    if [[ -n "$CATALOG_FILE" ]]; then
      ASSET_URL="$(/usr/bin/dirname "$CATALOG_FILE")/$ASSET_URL"
    else
      ASSET_URL="${CATALOG_URL%/*}/$ASSET_URL"
    fi
  fi
fi
ASSET_PATH="$TMP_DIR/$ASSET_NAME"
if [[ -n "$ASSET_FILE" ]]; then
  cp "$ASSET_FILE" "$ASSET_PATH"
elif [[ "$ASSET_URL" == file://* ]]; then
  cp "$(file_url_path "$ASSET_URL")" "$ASSET_PATH"
elif [[ "$ASSET_URL" == https://* ]]; then
  curl -fsSL "$ASSET_URL" -o "$ASSET_PATH"
elif [[ "$ASSET_URL" == *://* ]]; then
  fail "Release assets must use HTTPS or a local file URL."
else
  cp "$ASSET_URL" "$ASSET_PATH"
fi

ACTUAL_ARCHIVE_SHA="$(sha256_file "$ASSET_PATH")"
[[ "$(lower "$ACTUAL_ARCHIVE_SHA")" == "$(lower "$CATALOG_ARCHIVE_SHA")" ]] || fail "Checksum mismatch for $ASSET_NAME; installation stopped."

EXTRACT_DIR="$TMP_DIR/extract"
mkdir -p "$EXTRACT_DIR"
while IFS= read -r ENTRY; do
  case "$ENTRY" in
    /*|*\\*|../*|*/../*|*/..) fail "Release asset contains an unsafe archive path." ;;
  esac
done < <(/usr/bin/unzip -Z1 "$ASSET_PATH")
if /usr/bin/unzip -Z -v "$ASSET_PATH" | /usr/bin/grep -E 'Unix file attributes \((010|020|060|120|140)[0-7]{3} octal\)' >/dev/null; then
  fail "Release asset contains a symbolic link or unsupported filesystem entry."
fi
if ! /usr/bin/unzip -q "$ASSET_PATH" -d "$EXTRACT_DIR"; then
  fail "Release asset could not be extracted; installation stopped."
fi

PACK_MANIFESTS=()
while IFS= read -r FILE; do PACK_MANIFESTS+=("$FILE"); done < <(find "$EXTRACT_DIR" -name pack-manifest.json -type f)
[[ "${#PACK_MANIFESTS[@]}" == "1" ]] || fail "Release asset must contain exactly one pack-manifest.json."
PACK_MANIFEST="${PACK_MANIFESTS[0]}"
PAYLOAD_ROOT="$(/usr/bin/dirname "$PACK_MANIFEST")"
[[ "$(sha256_file "$PACK_MANIFEST")" == "$CATALOG_PACK_SHA" ]] || fail "Pack manifest checksum mismatch."
[[ "$(json_get "$PACK_MANIFEST" schema_version)" == "1" ]] || fail "Pack manifest schema is unsupported."
[[ "$(json_get "$PACK_MANIFEST" version)" == "$VERSION" ]] || fail "Pack version does not match."
[[ "$(json_get "$PACK_MANIFEST" platform)" == "macos-universal" ]] || fail "Pack platform does not match."
[[ "$(json_get "$PACK_MANIFEST" command)" == "codeheart-operating-kit" ]] || fail "Pack command does not match."
BINARY_REL="$(json_get "$PACK_MANIFEST" binary_path)"
BINARY_SHA="$(json_get "$PACK_MANIFEST" binary_sha256)"
CONTENT_REL="$(json_get "$PACK_MANIFEST" content_manifest_path)"
CONTENT_SHA="$(json_get "$PACK_MANIFEST" content_manifest_sha256)"
CHECKSUMS_REL="$(json_get "$PACK_MANIFEST" payload_checksums_path)"
CHECKSUMS_SHA="$(json_get "$PACK_MANIFEST" payload_checksums_sha256)"
[[ "$BINARY_REL" == "bin/codeheart-operating-kit" && "$CONTENT_REL" == "content-manifest.yaml" && "$CHECKSUMS_REL" == "checksums.txt" ]] || fail "Pack identity paths are invalid."
BINARY_PATH="$PAYLOAD_ROOT/$BINARY_REL"
CONTENT_PATH="$PAYLOAD_ROOT/$CONTENT_REL"
CHECKSUMS_PATH="$PAYLOAD_ROOT/$CHECKSUMS_REL"
[[ -f "$BINARY_PATH" && "$(sha256_file "$BINARY_PATH")" == "$BINARY_SHA" ]] || fail "Binary checksum mismatch."
[[ -f "$CONTENT_PATH" && "$(sha256_file "$CONTENT_PATH")" == "$CONTENT_SHA" ]] || fail "Content manifest checksum mismatch."
[[ -f "$CHECKSUMS_PATH" && "$(sha256_file "$CHECKSUMS_PATH")" == "$CHECKSUMS_SHA" ]] || fail "Payload checksum identity mismatch."

BINARY_LISTED=0
while IFS= read -r LINE; do
  EXPECTED="${LINE%%  *}"
  RELATIVE="${LINE#*  }"
  case "$RELATIVE" in ""|/*|*\\*|../*|*/../*|pack-manifest.json|checksums.txt) fail "Payload checksum path is invalid." ;; esac
  [[ -f "$PAYLOAD_ROOT/$RELATIVE" ]] || fail "Checksummed payload file is missing: $RELATIVE"
  [[ "$(sha256_file "$PAYLOAD_ROOT/$RELATIVE")" == "$EXPECTED" ]] || fail "Payload checksum mismatch for $RELATIVE"
  [[ "$RELATIVE" == "$BINARY_REL" ]] && BINARY_LISTED=1
done < "$CHECKSUMS_PATH"
[[ "$BINARY_LISTED" == "1" ]] || fail "Binary is absent from payload checksums."

BIN_DIR="$INSTALL_DIR/bin"
LIB_DIR="$INSTALL_DIR/lib"
TARGET_BINARY="$BIN_DIR/codeheart-operating-kit"
mkdir -p "$BIN_DIR"

LEGACY_FOUND=0
if [[ -f "$TARGET_BINARY" ]] && head -n 20 "$TARGET_BINARY" | grep -Eq 'python|codeheart_operating_kit'; then LEGACY_FOUND=1; fi
if compgen -G "$LIB_DIR/codeheart_operating_kit*" >/dev/null 2>&1; then LEGACY_FOUND=1; fi
if [[ "$LEGACY_FOUND" == "1" ]]; then
  echo "Legacy Python install detected; installing the self-contained binary and preserving legacy files."
fi

STAGING_DIR="$(mktemp -d "$INSTALL_DIR/.staging.XXXXXX")"
STAGED_BINARY="$STAGING_DIR/codeheart-operating-kit"
cp "$BINARY_PATH" "$STAGED_BINARY"
chmod 0755 "$STAGED_BINARY"
EXPECTED_VERSION="codeheart-operating-kit $VERSION"
[[ "$("$STAGED_BINARY" --version 2>&1)" == "$EXPECTED_VERSION" ]] || fail "Staged binary validation failed; previous runnable command preserved."
"$STAGED_BINARY" __verify-content-identity --path "$CONTENT_PATH" --version "$VERSION" >/dev/null 2>&1 || fail "Staged content identity validation failed; previous runnable command preserved."
"$STAGED_BINARY" __verify-release-evidence --catalog "$CATALOG_PATH" --version "$VERSION" >/dev/null 2>&1 || fail "Staged release evidence validation failed; previous runnable command preserved."

BACKUP_BINARY="$STAGING_DIR/previous-binary"
HAD_PREVIOUS=0
if [[ -f "$TARGET_BINARY" ]]; then
  mv "$TARGET_BINARY" "$BACKUP_BINARY"
  HAD_PREVIOUS=1
fi
restore_previous() {
  rm -f "$TARGET_BINARY"
  if [[ "$HAD_PREVIOUS" == "1" ]]; then mv "$BACKUP_BINARY" "$TARGET_BINARY"; fi
}
if ! mv "$STAGED_BINARY" "$TARGET_BINARY"; then
  restore_previous
  fail "Binary replacement failed; previous runnable command restored."
fi
chmod 0755 "$TARGET_BINARY"
if [[ "$("$TARGET_BINARY" --version 2>&1)" != "$EXPECTED_VERSION" ]]; then
  restore_previous
  fail "Installed binary validation failed; previous runnable command restored."
fi
rm -f "$BACKUP_BINARY"
rm -rf "$STAGING_DIR"
STAGING_DIR=""

echo "codeheart-operating-kit installed at $TARGET_BINARY"
if [[ "$LEGACY_FOUND" == "1" ]]; then echo "Legacy Python files were preserved under $INSTALL_DIR for later cleanup."; fi
case ":$PATH:" in *":$BIN_DIR:"*) ;; *) echo "Add this folder to PATH to run it by name: $BIN_DIR" ;; esac
echo "Next: codeheart-operating-kit onboard"
