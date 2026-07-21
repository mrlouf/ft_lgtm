# ════════════════════════════════════════════════════════════
#
#		Looks Good To Makefile
#
# ════════════════════════════════════════════════════════════

CLUSTER_NAME = lgtm-cluster

APP_URL		= http://lgtm.local
IPFS_URL	= http://ipfs.local
GRAFANA_URL	= http://grafana.local


# ════════════════════════════════════════════════════════════


help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Install the k3d cluster
	@echo 'Installing $(CLUSTER_NAME) on the host'
	@echo ''
	./scripts/setup-cluster.sh $(CLUSTER_NAME)

build: ## Build the docker images
	@echo 'Building docker images'
	@echo ''
	docker build -t lgtm_front:latest ./frontend
	docker build -t lgtm_back:latest ./backend
	k3d image import lgtm_front:latest -c $(CLUSTER_NAME)
	k3d image import lgtm_back:latest -c $(CLUSTER_NAME)

start: ## Start the cluster
	@echo 'Starting $(CLUSTER_NAME) on the host'
	@echo ''
	k3d cluster start $(CLUSTER_NAME)

up: ## Start all services
	@echo ""
	@echo "✅ All services started!"
	@echo "🌐 Access the application at $(APP_URL)"
	@echo "🌐 Access IPFS at $(IPFS_URL)"
	@echo "🌐 Access Grafana at $(GRAFANA_URL)"
	@echo ""

stop: ## Stop cluster
	@echo 'Stopping $(CLUSTER_NAME) on the host'
	@echo ''
	k3d cluster stop $(CLUSTER_NAME)

clean: ## Delete cluster
	@echo "Are you sure you want to delete the cluster $(CLUSTER_NAME)? This action cannot be undone. (y/n)"
	@read answer; \
	if [ "$$answer" != "${answer#[Yy]}" ] ;then \
		k3d cluster delete $(CLUSTER_NAME) \
	else \
		echo "Aborting cluster deletion."; \
	fi

PHONY: build up down logs clean help