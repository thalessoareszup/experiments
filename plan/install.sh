#!/bin/sh
set -e

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

mkdir -p build "$INSTALL_DIR"
go build -o build/plan cmd/cli/main.go
cp build/plan "$INSTALL_DIR/plan"

echo "Installed plan to $INSTALL_DIR/plan"
