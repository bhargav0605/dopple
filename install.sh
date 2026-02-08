#!/bin/bash

set -e

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="doppel"

echo "Installing doppel..."

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Map OS names
case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Create install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Check if binary exists locally (for local install)
if [ -f "./$BINARY_NAME" ]; then
    echo "Installing from local binary..."
    cp "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Binary not found. Please build first with: go build -o doppel"
    exit 1
fi

echo ""
echo "✓ doppel installed to $INSTALL_DIR/$BINARY_NAME"
echo ""

# Check if install directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "⚠️  $INSTALL_DIR is not in your PATH"
    echo ""
    echo "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
    echo ""
    echo "    export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
else
    echo "✓ Installation complete! Run 'doppel' from anywhere."
fi
