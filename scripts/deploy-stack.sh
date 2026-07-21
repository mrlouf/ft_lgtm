#!/bin/bash

# Colours
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Colour

CLUSTER_NAME=$1

set -e

#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#
#               Install ArgoCD via Helm            #
#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#

if kubectl -n argocd get deployment argocd-server &> /dev/null; then

    echo -e "${GREEN}ArgoCD already installed${NC}"

else
    echo -e "${BLUE}Installing ArgoCD...${NC}"

    helm repo add argo https://argoproj.github.io/argo-helm
    helm repo update

    helm install argocd argo/argo-cd \
    --namespace argocd \
    --create-namespace \
    --values argocd/values.yaml 1>/dev/null

    kubectl wait --namespace argocd --for=condition=available deployment/argocd-server --timeout=120s

fi

#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#
#               Deploy the app via ArgoCD          #
#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#