# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-w -s -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

.PHONY: build test run clean migration migrate-up migrate-down docker-build docker-run help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application with version information
	@echo "Building version $(VERSION)..."
	@go build $(LDFLAGS) -o bin/ecommerce-golang cmd/main.go

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/ecommerce-golang-linux cmd/main.go

test: ## Run tests
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

run: build ## Build and run the application
	@./bin/ecommerce-golang

clean: ## Clean build artifacts
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Cleaned build artifacts"

migration: ## Create a new migration file (usage: make migration name=migration_name)
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up: ## Run database migrations up
	@go run cmd/migrate/main.go up

migrate-down: ## Run database migrations down
	@go run cmd/migrate/main.go down

docker-build: ## Build Docker image with version information
	@echo "Building Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t ecommerce-api:$(VERSION) \
		-t ecommerce-api:latest \
		.

docker-run: ## Run the application in Docker
	@docker-compose up --build

docker-stop: ## Stop Docker containers
	@docker-compose down

docker-clean: ## Clean Docker resources
	@docker-compose down -v
	@docker rmi ecommerce-api:latest ecommerce-api:$(VERSION) 2>/dev/null || true

lint: ## Run linter
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
	fi

fmt: ## Format code
	@go fmt ./...
	@echo "Code formatted"

vet: ## Run go vet
	@go vet ./...

version: ## Show version information
	@echo "Version:    $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"

# Kubernetes deployment targets
k8s-apply: ## Apply all Kubernetes manifests
	@kubectl apply -f k8s/

k8s-delete: ## Delete all Kubernetes resources
	@kubectl delete -f k8s/

k8s-logs: ## Show logs from Kubernetes pods
	@kubectl logs -f -n ecommerce -l app=ecommerce-api

k8s-status: ## Show Kubernetes deployment status
	@kubectl get pods,svc,hpa -n ecommerce