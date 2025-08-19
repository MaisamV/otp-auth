# OTP Authentication Service Makefile

# Variables
APP_NAME := otp-auth
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO_VERSION := $(shell go version | awk '{print $$3}')

# Build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT) -X main.goVersion=$(GO_VERSION)"

# Docker
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
DOCKER_LATEST := $(APP_NAME):latest

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

.PHONY: help build run test clean docker-build docker-run docker-stop deps lint fmt vet security audit migrate-up migrate-down dev prod logs

# Default target
all: clean deps test build

## Help
help: ## Show this help message
	@echo "$(BLUE)$(APP_NAME) - Available commands:$(RESET)"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

deps: ## Install dependencies
	@echo "$(BLUE)Installing dependencies...$(RESET)"
	go mod download
	go mod tidy

run: ## Run the application locally
	@echo "$(BLUE)Running application...$(RESET)"
	go run cmd/server/main.go

dev: ## Run in development mode with hot reload (requires air)
	@echo "$(BLUE)Starting development server with hot reload...$(RESET)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)Air not found. Install with: go install github.com/cosmtrek/air@latest$(RESET)"; \
		echo "$(YELLOW)Falling back to regular run...$(RESET)"; \
		make run; \
	fi

fmt: ## Format Go code
	@echo "$(BLUE)Formatting code...$(RESET)"
	go fmt ./...

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(RESET)"
	go vet ./...

lint: ## Run golangci-lint
	@echo "$(BLUE)Running linter...$(RESET)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(RED)golangci-lint not found. Install from https://golangci-lint.run/usage/install/$(RESET)"; \
	fi

security: ## Run security checks with gosec
	@echo "$(BLUE)Running security checks...$(RESET)"
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "$(RED)gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest$(RESET)"; \
	fi

audit: ## Run dependency audit
	@echo "$(BLUE)Running dependency audit...$(RESET)"
	@if command -v nancy > /dev/null; then \
		go list -json -deps ./... | nancy sleuth; \
	else \
		echo "$(YELLOW)nancy not found. Install with: go install github.com/sonatypecommunity/nancy@latest$(RESET)"; \
		echo "$(BLUE)Running go mod audit instead...$(RESET)"; \
		go mod verify; \
	fi

##@ Testing

test: ## Run tests
	@echo "$(BLUE)Running tests...$(RESET)"
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(RESET)"

test-integration: ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(RESET)"
	go test -v -tags=integration ./...

benchmark: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	go test -bench=. -benchmem ./...

##@ Building

build: ## Build the application
	@echo "$(BLUE)Building $(APP_NAME)...$(RESET)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME) cmd/server/main.go
	@echo "$(GREEN)Build complete: bin/$(APP_NAME)$(RESET)"

build-windows: ## Build for Windows
	@echo "$(BLUE)Building $(APP_NAME) for Windows...$(RESET)"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME).exe cmd/server/main.go
	@echo "$(GREEN)Build complete: bin/$(APP_NAME).exe$(RESET)"

build-mac: ## Build for macOS
	@echo "$(BLUE)Building $(APP_NAME) for macOS...$(RESET)"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME)-darwin cmd/server/main.go
	@echo "$(GREEN)Build complete: bin/$(APP_NAME)-darwin$(RESET)"

build-all: build build-windows build-mac ## Build for all platforms

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning...$(RESET)"
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache
	@echo "$(GREEN)Clean complete$(RESET)"

##@ Docker

docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(RESET)"
	docker build -t $(DOCKER_IMAGE) -t $(DOCKER_LATEST) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE)$(RESET)"

docker-run: ## Run application with Docker Compose
	@echo "$(BLUE)Starting services with Docker Compose...$(RESET)"
	docker-compose up -d
	@echo "$(GREEN)Services started. Check status with: make docker-status$(RESET)"

docker-stop: ## Stop Docker Compose services
	@echo "$(BLUE)Stopping Docker Compose services...$(RESET)"
	docker-compose down

docker-restart: docker-stop docker-run ## Restart Docker services

docker-logs: ## Show Docker Compose logs
	docker-compose logs -f

docker-status: ## Show Docker Compose status
	docker-compose ps

docker-clean: ## Clean Docker images and containers
	@echo "$(BLUE)Cleaning Docker resources...$(RESET)"
	docker-compose down -v --remove-orphans
	docker system prune -f
	@echo "$(GREEN)Docker cleanup complete$(RESET)"

##@ Database

migrate-up: ## Run database migrations
	@echo "$(BLUE)Running database migrations...$(RESET)"
	@if command -v migrate > /dev/null; then \
		migrate -path internal/infrastructure/persistence/postgres/migrations -database "postgres://postgres:postgres@localhost:5432/otp_auth?sslmode=disable" up; \
	else \
		echo "$(RED)migrate not found. Install from https://github.com/golang-migrate/migrate$(RESET)"; \
	fi

migrate-down: ## Rollback database migrations
	@echo "$(BLUE)Rolling back database migrations...$(RESET)"
	@if command -v migrate > /dev/null; then \
		migrate -path internal/infrastructure/persistence/postgres/migrations -database "postgres://postgres:postgres@localhost:5432/otp_auth?sslmode=disable" down; \
	else \
		echo "$(RED)migrate not found. Install from https://github.com/golang-migrate/migrate$(RESET)"; \
	fi

migrate-force: ## Force migration version
	@echo "$(BLUE)Forcing migration version...$(RESET)"
	@read -p "Enter version number: " version; \
	migrate -path internal/infrastructure/persistence/postgres/migrations -database "postgres://postgres:postgres@localhost:5432/otp_auth?sslmode=disable" force $$version

##@ Production

prod: ## Run production environment
	@echo "$(BLUE)Starting production environment...$(RESET)"
	docker-compose -f docker-compose.yml --profile production up -d
	@echo "$(GREEN)Production environment started$(RESET)"

prod-logs: ## Show production logs
	docker-compose -f docker-compose.yml --profile production logs -f

prod-stop: ## Stop production environment
	@echo "$(BLUE)Stopping production environment...$(RESET)"
	docker-compose -f docker-compose.yml --profile production down

##@ Monitoring

health: ## Check application health
	@echo "$(BLUE)Checking application health...$(RESET)"
	@curl -s http://localhost:8080/health | jq . || echo "$(RED)Health check failed$(RESET)"

ready: ## Check application readiness
	@echo "$(BLUE)Checking application readiness...$(RESET)"
	@curl -s http://localhost:8080/ready | jq . || echo "$(RED)Readiness check failed$(RESET)"

live: ## Check application liveness
	@echo "$(BLUE)Checking application liveness...$(RESET)"
	@curl -s http://localhost:8080/live | jq . || echo "$(RED)Liveness check failed$(RESET)"

status: health ready live ## Check all health endpoints

logs: ## Show application logs (Docker)
	docker-compose logs -f otp-auth

##@ Utilities

generate: ## Generate code (mocks, etc.)
	@echo "$(BLUE)Generating code...$(RESET)"
	go generate ./...

mod-update: ## Update Go modules
	@echo "$(BLUE)Updating Go modules...$(RESET)"
	go get -u ./...
	go mod tidy

install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(RESET)"
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/sonatypecommunity/nancy@latest
	@echo "$(GREEN)Development tools installed$(RESET)"

version: ## Show version information
	@echo "$(BLUE)Version Information:$(RESET)"
	@echo "App Name: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

##@ Quick Start

quick-start: deps docker-run ## Quick start for development
	@echo "$(GREEN)Quick start complete!$(RESET)"
	@echo "$(BLUE)Application should be running at: http://localhost:8080$(RESET)"
	@echo "$(BLUE)Check health: make health$(RESET)"

full-setup: install-tools deps test build docker-build ## Full development setup
	@echo "$(GREEN)Full setup complete!$(RESET)"
	@echo "$(BLUE)Ready for development. Run 'make dev' to start with hot reload$(RESET)"