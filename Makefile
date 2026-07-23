.PHONY: build run test test-race clean docker docker-up docker-down docker-logs docker-logs-follow lint fmt fmt-check mod-check ci help

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
DOCKER_COMPOSE=docker-compose -f docker/docker-compose.yaml

# golangci-lint (installed into $(BUILD_DIR) on demand so `make ci` is
# self-contained and runs the same locally as in CI)
GOLANGCI_LINT_VERSION=v2.12.2
GOLANGCI_LINT=$(BUILD_DIR)/golangci-lint

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

# Install golangci-lint into $(BUILD_DIR) if it is not already present
$(GOLANGCI_LINT):
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(BUILD_DIR) $(GOLANGCI_LINT_VERSION)

# Lint code (golangci-lint's standard set includes go vet)
lint: $(GOLANGCI_LINT)
	@echo "Linting code..."
	@$(GOLANGCI_LINT) run ./...

# Check formatting (fails if not formatted)
fmt-check:
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Code not formatted. Run 'make fmt'" && gofmt -l . && exit 1)

# Check go.mod is tidy
mod-check:
	@echo "Checking go.mod..."
	@go mod tidy
	@git diff --exit-code go.mod go.sum || (echo "go.mod/go.sum not tidy. Run 'go mod tidy'" && exit 1)

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -race -v ./...

# CI recipe - runs all checks
ci: fmt-check mod-check lint build test
	@echo "CI checks passed!"

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

# View logs (alias)
docker-logs-follow:
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
	@echo "  fmt-check     - Check code formatting (CI)"
	@echo "  mod-check     - Check go.mod is tidy (CI)"
	@echo "  lint          - Lint code"
	@echo "  test-race     - Run tests with race detection"
	@echo "  ci            - Run all CI checks (fmt, mod, lint, build, test)"
	@echo "  docker        - Build Docker image"
	@echo "  docker-up     - Start all services with Docker Compose"
	@echo "  docker-down   - Stop all services"
	@echo "  docker-logs   - View service logs"
	@echo "  docker-rebuild- Rebuild and restart services"
	@echo "  docker-clean  - Remove all containers and volumes"
	@echo "  help          - Show this help message"

# Default target
.DEFAULT_GOAL := build
