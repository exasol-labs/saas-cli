#!/usr/bin/env sh
set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
REPO="exasol-labs/saas-cli"

# ---------------------------------------------------------------------------
# Cleanup on exit
# ---------------------------------------------------------------------------

TMP_FILE=""
cleanup() {
  if [ -n "$TMP_FILE" ] && [ -f "$TMP_FILE" ]; then
    rm -f "$TMP_FILE"
  fi
}
trap cleanup EXIT

# ---------------------------------------------------------------------------
# Platform detection
# ---------------------------------------------------------------------------

printf "Detecting platform... "
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin)  OS_NAME="darwin" ;;
  Linux)   OS_NAME="linux"  ;;
  *)
    echo ""
    echo "Error: unsupported platform: ${OS}/${ARCH}"
    echo "Download manually from: https://github.com/${REPO}/releases"
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64)          ARCH_NAME="amd64" ;;
  arm64|aarch64)   ARCH_NAME="arm64" ;;
  *)
    echo ""
    echo "Error: unsupported platform: ${OS}/${ARCH}"
    echo "Download manually from: https://github.com/${REPO}/releases"
    exit 1
    ;;
esac

BINARY="exasol-saas-${OS_NAME}-${ARCH_NAME}"
echo "${OS_NAME}/${ARCH_NAME}"

# ---------------------------------------------------------------------------
# Version resolution
# ---------------------------------------------------------------------------

printf "Resolving latest version... "
if [ -z "${VERSION:-}" ]; then
  VERSION="$(curl -fsSL -H "User-Agent: exasol-saas-installer" \
    "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' | sed 's/.*"tag_name": *"\(.*\)".*/\1/')"
  if [ -z "$VERSION" ]; then
    echo ""
    echo "Error: could not resolve latest version from GitHub API"
    echo "Try setting VERSION manually: VERSION=v0.1.0 sh install.sh"
    exit 1
  fi
fi
echo "$VERSION"

# ---------------------------------------------------------------------------
# Download
# ---------------------------------------------------------------------------

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"
TMP_FILE="$(mktemp)"

echo "Downloading ${BINARY}..."
if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"; then
  echo "Error: download failed from ${DOWNLOAD_URL}"
  exit 1
fi

# ---------------------------------------------------------------------------
# Install
# ---------------------------------------------------------------------------

INSTALL_PATH="${INSTALL_DIR}/exasol-saas"
echo "Installing to ${INSTALL_PATH}..."
chmod +x "$TMP_FILE"
mv "$TMP_FILE" "$INSTALL_PATH"
TMP_FILE=""  # prevent cleanup trap from trying to remove it

# ---------------------------------------------------------------------------
# Verify
# ---------------------------------------------------------------------------

"$INSTALL_PATH" --version
echo "exasol-saas ${VERSION} installed to ${INSTALL_PATH}"
