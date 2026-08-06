#!/bin/bash

# Colours
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Colour

ADD_ARGOCD=$1

set -e

if [[ $ADD_ARGOCD == "TRUE" ]]; then
    HOSTS='172.18.0.2	lgtm.local argocd.lgtm.local ipfs.lgtm.local grafana.lgtm.local'
else
    HOSTS='172.18.0.2	lgtm.local ipfs.lgtm.local grafana.lgtm.local'
fi

if [[ $(grep -c "$HOSTS" /etc/hosts) -gt 0 ]]; then
    echo -e "${YELLOW}Hosts already exist in /etc/hosts${NC}"
    exit 0

else
    echo -e "${BLUE}Adding hosts to /etc/hosts${NC}"
    sudo -- sh -c "echo '$HOSTS' >> /etc/hosts"
    echo -e "${GREEN}/etc/hosts configured${NC}"

fi

