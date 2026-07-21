#!/bin/bash

# Colours
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Colour

CLUSTER_NAME=$1

TAG=$(git rev-parse --short HEAD)

set -e

docker build -t ghcr.io/mrlouf/lgtm_back:sha-$TAG ./backend; \
docker build -t ghcr.io/mrlouf/lgtm_front:sha-$TAG ./frontend; \

docker push ghcr.io/mrlouf/lgtm_back:sha-$TAG; \
docker push ghcr.io/mrlouf/lgtm_front:sha-$TAG
