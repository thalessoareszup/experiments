#!/bin/sh
set -e

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

cd cli
mkdir -p build "$INSTALL_DIR"
go build -o build/plan main.go
cp build/plan "$INSTALL_DIR/plan"

echo "Installed plan to $INSTALL_DIR/plan"
