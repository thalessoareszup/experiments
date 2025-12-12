#!/bin/sh
set -e

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

mkdir -p build "$INSTALL_DIR"
go build -o build/wk ./cmd/wk
cp build/wk "$INSTALL_DIR/wk"

echo "Installed wk to $INSTALL_DIR/wk"
