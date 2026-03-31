#!/bin/sh
set -e

REPO="orvad/tuibookie"
INSTALL_DIR="/usr/local/bin"
BINARY="tuibookie"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  darwin) OS="darwin" ;;
  linux) OS="linux" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Get latest version
echo "Fetching latest version..."
VERSION=$(curl -sI "https://github.com/${REPO}/releases/latest" | grep -i "^location:" | sed 's/.*tag\///' | tr -d '\r\n')
if [ -z "$VERSION" ]; then
  echo "Failed to determine latest version"
  exit 1
fi
echo "Latest version: $VERSION"

# Download
FILENAME="${BINARY}-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"
echo "Downloading ${FILENAME}..."
curl -sL -o "/tmp/${BINARY}" "$URL"

if [ ! -s "/tmp/${BINARY}" ]; then
  echo "Download failed. Check that a release exists for your platform: ${OS}/${ARCH}"
  exit 1
fi

# Install
chmod +x "/tmp/${BINARY}"
if [ -w "$INSTALL_DIR" ]; then
  mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

echo "Installed ${BINARY} ${VERSION} to ${INSTALL_DIR}/${BINARY}"
echo "Run 'tuibookie' to get started!"
