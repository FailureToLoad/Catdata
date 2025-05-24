#!/bin/bash

set -e

INSTALL_DIR="$HOME/.local/bin"

echo "Installing Tailwind CSS standalone CLI..."

mkdir -p "$INSTALL_DIR"

echo "Downloading tailwindcss binary..."

curl -L -o "$INSTALL_DIR/tailwindcss" "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64"

chmod +x "$INSTALL_DIR/tailwindcss"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "Adding $INSTALL_DIR to PATH..."
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    echo "Please run 'source ~/.bashrc' or restart your terminal to update PATH"
else
    echo "$INSTALL_DIR is already in PATH"
fi

echo "Tailwind CSS installed successfully to $INSTALL_DIR/tailwindcss"
echo "Run 'tailwindcss --help' to get started"