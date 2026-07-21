#!/bin/bash

CLUSTER_NAME=$1

set -e

# Install Docker first if not present:
if ! systemctl is-active --quiet docker; then
    echo "Installing Docker..."
    # Add Docker's official GPG key:
    sudo apt-get update
    sudo apt-get install -y ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    
    # Actually install Docker
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    # Start and enable Docker service
    sudo systemctl start docker
    sudo systemctl enable docker
else
    echo "Docker already installed and running"
fi

#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#
#                   Install k3s                    #
#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#

if ! command -v k3s &> /dev/null; then
    
    echo "Installing k3s..."
    curl -sfL https://get.k3s.io | sh -

else
    echo "k3s already installed"
fi


#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#
#                   Setup k3d                      #
#~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=#

# Install k3d

if ! command -v k3d &> /dev/null; then
    echo "Installing k3d..."
    wget -q -O - https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
else
    echo "k3d already installed"
fi

# Install kubectl

if command -v kubectl &> /dev/null; then
    echo "kubectl is already installed"

else

    echo "kubectl not found, installing..."
    curl -LO "https://dl.k8s.io/release/$(curl -sL https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
    rm kubectl

fi

# Create the k3d cluster if it doesn't exist
if ! k3d cluster list | grep -q "$CLUSTER_NAME"; then

    echo "Creating k3d cluster $CLUSTER_NAME..."
    k3d cluster create $CLUSTER_NAME --config ./k3d/lgtm.yaml

else

    echo "k3d cluster $CLUSTER_NAME already exists"

fi

# Set the kubeconfig context
export KUBECONFIG=$(k3d kubeconfig write $CLUSTER_NAME)
