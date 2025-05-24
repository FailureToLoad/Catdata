#!/bin/bash

set -e

INSTALL_DIR="$HOME/.local/bin"

mkdir -p "$INSTALL_DIR"

echo "Downloading templ binary..."

curl -L -o "/tmp/templ.tar.gz" "https://github.com/a-h/templ/releases/latest/download/templ_Linux_x86_64.tar.gz"

echo "Extracting..."

tar -xzf "/tmp/templ.tar.gz" -C "/tmp"

mv "/tmp/templ" "$INSTALL_DIR/templ"

rm "/tmp/templ.tar.gz"

chmod +x "$INSTALL_DIR/templ"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "Adding $INSTALL_DIR to PATH..."
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    echo "Please run 'source ~/.bashrc' or restart your terminal to update PATH"
else
    echo "$INSTALL_DIR is already in PATH"
fi

echo "templ installed successfully to $INSTALL_DIR/templ"
echo "Run 'templ --help' to get started"