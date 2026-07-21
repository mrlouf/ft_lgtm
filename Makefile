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


# ════════════════════════════════════════════════════════════


help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: setup build start up ## Setup, build, start and run all services

setup: ## Install the k3d cluster
	@printf "\n$(BLUE)Installing $(CLUSTER_NAME) on the host$(NC)\n"
	@echo ''
	./scripts/setup-cluster.sh $(CLUSTER_NAME)

build: ## Build the docker images
	@printf "\n$(GREEN)Building docker images$(NC)\n"
	@echo ''
	docker build -t lgtm_front:latest ./frontend
	docker build -t lgtm_back:latest ./backend
	k3d image import lgtm_front:latest -c $(CLUSTER_NAME)
	k3d image import lgtm_back:latest -c $(CLUSTER_NAME)

start: ## Start the cluster
	@printf "\n$(YELLOW)Starting $(CLUSTER_NAME) on the host$(NC)\n"
	@echo ''
	k3d cluster start $(CLUSTER_NAME)

up: ## Start all services
	@printf "\n$(BLUE) All services started!$(NC)\n\n"
	@echo "🌐 Access the application at $(APP_URL)"
	@echo "🌐 Access IPFS at $(IPFS_URL)"
	@echo "🌐 Access Grafana at $(GRAFANA_URL)"
	@echo ""

stop: ## Stop cluster
	@printf "\n$(RED)Stopping $(CLUSTER_NAME) on the host$(NC)\n"
	@echo ''
	k3d cluster stop $(CLUSTER_NAME)

clean: ## Delete cluster
	@printf "\n$(RED)Are you sure you want to delete the cluster $(CLUSTER_NAME)? This action cannot be undone. (y/n)$(NC)\n"
	@read answer; \
	if [ "$$answer" != "$(answer#[Yy])" ] ;then \
		k3d cluster delete $(CLUSTER_NAME) \
	else \
		echo "Aborting cluster deletion."; \
	fi

PHONY: build up down logs clean help