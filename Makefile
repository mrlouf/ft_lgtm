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
IPFS_URL	= http://ipfs.lgtm.local
GRAFANA_URL	= http://grafana.lgtm.local
ARGOCD_URL	= http://argocd.lgtm.local



APP_URL_DEV = http://localhost:5173
APP_IPFS_DEV = http://localhost:5001/webui
APP_PROMETHEUS_DEV = http://localhost:9090
APP_GRAFANA_DEV = http://localhost:3000
APP_OTEL_DEV = http://localhost:55679/debug/servicez


# ════════════════════════════════════════════════════════════


help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: cluster deploy ## Setup, build and deploy all services

cluster: ## Install the k3d cluster
	@printf "\n$(YELLOW)Installing $(CLUSTER_NAME) on the host$(NC)\n"
	@echo ''
	@./scripts/setup-cluster.sh $(CLUSTER_NAME)

build: ## Build the docker images and push them to the GHCR registry
	@printf "\n$(YELLOW)Building docker images$(NC)\n"
	@echo ''
	@./scripts/build-images.sh $(CLUSTER_NAME)


deploy: ## Deploy all services
	@printf "\n$(YELLOW)Deploying the stack... this may take a moment $(NC)\n\n"
	@./scripts/deploy-stack.sh $(CLUSTER_NAME)
	@./scripts/append-hosts.sh $(CLUSTER_NAME)
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
		k3d cluster delete $(CLUSTER_NAME); \
	else \
		echo "Aborting cluster deletion."; \
	fi
	
develop: ## Start the development environment
	@printf "\n$(YELLOW)Starting the development environment...$(NC)\n"
	@echo ''
	@docker compose -f dev/docker-compose.yaml up --build -d
	@echo ''
	@echo -e "$(BLUE)🌐 Access the application at $(APP_URL_DEV)$(NC)"
	@echo -e "$(BLUE)🌐 Access IPFS at $(APP_IPFS_DEV)$(NC)"
	@echo ""
	@echo -e "$(BLUE)🌐 Access the OpenTelemetry Collector at $(APP_OTEL_DEV)$(NC)"
	@echo -e "$(BLUE)🌐 Access Prometheus at $(APP_PROMETHEUS_DEV)$(NC)"
	@echo -e "$(BLUE)🌐 Access Grafana at $(APP_GRAFANA_DEV)$(NC)"
	@echo ''

develop-stop: ## Stop the development environment
	@printf "\n$(RED)Stopping the development environment...$(NC)\n"
	@echo ''
	@docker compose -f dev/docker-compose.yaml down

PHONY: help all cluster build start deploy stop clean develop develop-stop
