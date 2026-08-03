#!/bin/bash

# Colours
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Colour

CLUSTER_NAME=$1

set -e

install_k9s() {
    if command -v brew &> /dev/null; then
        echo -e "${BLUE}Installing k9s via Homebrew...${NC}"
        brew install derailed/k9s/k9s
    elif command -v apt &> /dev/null; then
        echo -e "${BLUE}Installing k9s via apt...${NC}"
        wget https://github.com/derailed/k9s/releases/latest/download/k9s_linux_amd64.deb && sudo apt install ./k9s_linux_amd64.deb && rm k9s_linux_amd64.deb
    elif command -v pacman &> /dev/null; then
        echo -e "${BLUE}Installing k9s via pacman...${NC}"
        sudo pacman -S k9s
    elif command -v dnf &> /dev/null; then
        echo -e "${BLUE}Installing k9s via dnf...${NC}"
        sudo dnf install k9s
    else
        echo -e "${BLUE}No supported package manager found. Installing via webinstall...${NC}"
        curl -sS https://webinstall.dev/k9s | bash || echo -e "${RED}Failed to install k9s via webinstall, consider installing from source.${NC}"
    fi
}

if ! command -v k9s &> /dev/null; then
    echo -e "${BLUE}k9s could not be found, installing...${NC}"

    install_k9s
    
else
    echo -e "${GREEN}k9s is already installed.${NC}"
    k9s version
fi
