#!/usr/bin/env bash

set -e

VERSION="1.0.0"
INSTALL_DIR="/usr/local/bin"

echo "Installing Kyra SDK $VERSION..."

# Detect platform
OS=$(uname -s)
ARCH=$(uname -m)

if [[ "$OS" == "Linux" ]]; then
    PLATFORM="linux"
elif [[ "$OS" == "Darwin" ]]; then
    PLATFORM="darwin"
else
    echo "Unsupported OS: $OS"
    exit 1
fi

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

KYRA="kyra-$VERSION-$PLATFORM-$ARCH"
KYRAC="kyrac-$VERSION-$PLATFORM-$ARCH"

echo "Detected platform: $PLATFORM-$ARCH"
echo "Installing binaries..."

sudo cp "$KYRA" "$INSTALL_DIR/kyra"
sudo cp "$KYRAC" "$INSTALL_DIR/kyrac"

sudo chmod +x "$INSTALL_DIR/kyra"
sudo chmod +x "$INSTALL_DIR/kyrac"

echo "Kyra SDK installed successfully."
echo "You can now run:"
echo "  kyra -kbc file.kbc"
echo "  kyrac -kbc file.kyra"
