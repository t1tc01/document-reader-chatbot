# Document Reader Chatbot Backend Makefile

.PHONY: help build run test clean deps migration-generate migration-apply docker-up docker-down docker-clean

# Variables
BINARY_NAME=document-reader-chatbot
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_SERVER=./cmd/server/main.go
MAIN_MIGRATE=./cmd/migrate/main.go

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Dependencies
deps: ## Download and install dependencies
	go mod download
	go mod tidy

# Build commands
build: ## Build the server binary
	@echo "Building server..."
	go build -o $(BINARY_PATH) $(MAIN_SERVER)

build-migrate: ## Build the migration binary
	@echo "Building migration tool..."
	go build -o ./bin/migrate $(MAIN_MIGRATE)

# Run commands
run: ## Run the server (used by docker-compose)
	@echo "Starting server..."
	go run $(MAIN_SERVER)

dev: ## Run the server in development mode with hot reload
	@echo "Starting development server..."
	air || go run $(MAIN_SERVER)

# Database migrations
migration-generate: ## Generate new migration files
	@echo "Generating migration..."
	atlas migrate diff --env local

migration-apply: ## Apply pending migrations
	@echo "Applying migrations..."
	go run $(MAIN_MIGRATE)

migration-status: ## Check migration status
	@echo "Checking migration status..."
	atlas migrate status --env local

# Testing
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -cover ./...

test-race: ## Run tests with race condition detection
	@echo "Running tests with race detection..."
	go test -race -v ./...

# Code quality
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./... || echo "golangci-lint not installed, skipping..."

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Docker commands
docker-up: ## Start all services with docker-compose
	@echo "Starting Docker services..."
	docker-compose up -d

docker-down: ## Stop all services
	@echo "Stopping Docker services..."
	docker-compose down

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker-compose build backend

docker-logs: ## Show Docker logs
	@echo "Showing Docker logs..."
	docker-compose logs -f backend

docker-clean: ## Clean up Docker resources
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f

# Database commands
db-setup: ## Setup database (start postgres and run migrations)
	@echo "Setting up database..."
	docker-compose up -d postgres
	@echo "Waiting for postgres to be ready..."
	sleep 5
	$(MAKE) migration-apply

db-reset: ## Reset database (drop and recreate)
	@echo "Resetting database..."
	docker-compose down postgres
	docker volume rm backend_postgres_data || true
	docker-compose up -d postgres
	@echo "Waiting for postgres to be ready..."
	sleep 5
	$(MAKE) migration-apply

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf ./bin/
	go clean

# Installation
install: build ## Install the binary to GOPATH/bin
	go install $(MAIN_SERVER)

# Development setup
setup: deps docker-up db-setup ## Complete development environment setup
	@echo "Development environment ready!"

# Production build
build-prod: ## Build for production
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_PATH) $(MAIN_SERVER)
