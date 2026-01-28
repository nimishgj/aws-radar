.PHONY: build run test clean docker docker-up docker-down lint fmt help

# Variables
BINARY_NAME=aws-radar
BUILD_DIR=bin
MAIN_PATH=./cmd/aws-radar

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Docker parameters
DOCKER_COMPOSE=docker compose -f docker/docker-compose.yaml

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application locally
run:
	@echo "Running $(BINARY_NAME)..."
	$(GORUN) $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Lint code
lint:
	@echo "Linting code..."
	$(GOVET) ./...

# Build Docker image
docker:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

# Start all services with Docker Compose
docker-up:
	@echo "Starting all services..."
	$(DOCKER_COMPOSE) up -d

# Stop all services
docker-down:
	@echo "Stopping all services..."
	$(DOCKER_COMPOSE) down

# View logs
docker-logs:
	@echo "Showing logs..."
	$(DOCKER_COMPOSE) logs -f

# Rebuild and restart
docker-rebuild:
	@echo "Rebuilding and restarting services..."
	$(DOCKER_COMPOSE) up -d --build

# Remove all volumes
docker-clean:
	@echo "Removing all containers and volumes..."
	$(DOCKER_COMPOSE) down -v

# Show help
help:
	@echo "AWS Radar - Available targets:"
	@echo ""
	@echo "  build         - Build the binary"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  docker        - Build Docker image"
	@echo "  docker-up     - Start all services with Docker Compose"
	@echo "  docker-down   - Stop all services"
	@echo "  docker-logs   - View service logs"
	@echo "  docker-rebuild- Rebuild and restart services"
	@echo "  docker-clean  - Remove all containers and volumes"
	@echo "  help          - Show this help message"

# Default target
.DEFAULT_GOAL := help
