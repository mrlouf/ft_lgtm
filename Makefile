# ════════════════════════════════════════════════════════════
#
#		Looks Good To Makefile
#
# ════════════════════════════════════════════════════════════


RED 		= \033[0;91m
GREEN 		= \033[0;92m
YELLOW		= \033[0;93m
BLUE		= \033[0;94m
NC			= \033[0m


# ════════════════════════════════════════════════════════════


CLUSTER_NAME = lgtm

APP_URL		= http://lgtm.local
IPFS_URL	= http://ipfs.local
GRAFANA_URL	= http://grafana.local
ARGOCD_URL	= http://argocd.local


# ════════════════════════════════════════════════════════════


help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: cluster build deploy ## Setup, build and deploy all services

cluster: ## Install the k3d cluster
	@printf "\n$(BLUE)Installing $(CLUSTER_NAME) on the host$(NC)\n"
	@echo ''
	@./scripts/setup-cluster.sh $(CLUSTER_NAME)

build: ## Build the docker images and push them to the GHCR registry
	@printf "\n$(BLUE)Building docker images$(NC)\n"
	@echo ''
	@./scripts/build-images.sh $(CLUSTER_NAME)


deploy: ## Deploy all services
	@printf "\n$(BLUE) Deploying the stack... this may take a moment $(NC)\n\n"
	@./scripts/deploy-stack.sh $(CLUSTER_NAME)
	@echo ''
	@echo -e "$(BLUE)🌐 Access the application at $(APP_URL)$(NC)"
	@echo -e "$(BLUE)🌐 Access IPFS at $(IPFS_URL)$(NC)"
	@echo -e "$(BLUE)🌐 Access Grafana at $(GRAFANA_URL)$(NC)"
	@echo -e "$(BLUE)🌐 Access ArgoCD at $(ARGOCD_URL)$(NC)"
	@echo ""

stop: ## Stop cluster
	@printf "\n$(RED)Stopping $(CLUSTER_NAME) on the host$(NC)\n"
	@echo ''
	k3d cluster stop $(CLUSTER_NAME)

start: ## Start the cluster
	@printf "\n$(YELLOW)Starting $(CLUSTER_NAME) on the host$(NC)\n"
	@echo ''
	k3d cluster start $(CLUSTER_NAME)

clean: ## Delete cluster
	@printf "\n$(RED)Are you sure you want to delete the cluster $(CLUSTER_NAME)? This action cannot be undone. (y/n)$(NC)\n"
	@read answer; \
	if [ "$$answer" != "$(answer#[Yy])" ] ;then \
		k3d cluster delete $(CLUSTER_NAME) \
	else \
		echo "Aborting cluster deletion."; \
	fi

PHONY: help all cluster build start deploy stop clean
