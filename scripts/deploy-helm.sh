#!/bin/bash

# Colours
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Colour

CLUSTER_NAME=$1

set -e

# Set the kubeconfig context
export KUBECONFIG=$(k3d kubeconfig write $CLUSTER_NAME)

#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#
#               Deploy LGTM with Helm              #
#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#

if kubectl get ns lgtm &> /dev/null; then

    echo -e "${GREEN}LGTM already deployed${NC}"

else
    echo -e "${BLUE}Deploying LGTM...${NC}"

    helm upgrade --install lgtm ./helm/lgtm \
    --namespace lgtm \
    --create-namespace \
    --values ./helm/lgtm/values.yaml 1>/dev/null

fi
