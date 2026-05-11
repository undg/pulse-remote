#!/usr/bin/env bash
set -euo pipefail

# Colors
BOLD='\033[1m'
BLUE='\033[34m'
GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'

# Get latest available Go version from mise
VERSION=$(mise ls go | tail -1 | awk '{print $2}')

echo -e "${BOLD}${BLUE}==>${RESET} ${BOLD}Updating to Go ${VERSION}${RESET}"

echo -e "${BOLD}${GREEN}==>${RESET} mise install go@${VERSION}"
mise install "go@${VERSION}"

echo -e "${BOLD}${GREEN}==>${RESET} mise use go@${VERSION}"
mise use "go@${VERSION}"

echo -e "${BOLD}${GREEN}==>${RESET} go mod edit -go=${VERSION}"
go mod edit "-go=${VERSION}"

echo -e "${BOLD}${GREEN}==>${RESET} go mod tidy -v"
go mod tidy -v

echo -e "${BOLD}${YELLOW}==>${RESET} ${BOLD}Done. Go updated to ${VERSION}${RESET}"
