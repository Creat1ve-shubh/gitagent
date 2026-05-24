#!/usr/bin/env bash
set -euo pipefail

# Simple Go-based installer for gitclaw

check_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1"
    exit 1
  fi
}

check_cmd git
check_cmd go

ROOT_DIR="$(pwd)"
BIN_DIR="${HOME}/.local/bin"
mkdir -p "${BIN_DIR}"

cd "${ROOT_DIR}"

go build -o "${BIN_DIR}/gitclaw" ./cmd/gitclaw

echo "gitclaw installed to ${BIN_DIR}/gitclaw"

echo "Ensure ${BIN_DIR} is on your PATH."
