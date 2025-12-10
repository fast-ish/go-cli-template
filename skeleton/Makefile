# ${{values.name}} makefile
.PHONY: help build install test lint format clean run docker-build docker-run release

# Build variables
BINARY_NAME=${{values.name}}
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Colors
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "$(GREEN)${{values.name}} - ${{values.description}}$(NC)"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'

## build: Build the CLI binary
build:
	@echo "$(YELLOW)Building ${{values.name}}...$(NC)"
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)
	@echo "$(GREEN)✓ Binary built: bin/$(BINARY_NAME)$(NC)"

## install: Install dependencies
install:
	@echo "$(YELLOW)Installing dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

## test: Run unit tests
test:
	@echo "$(YELLOW)Running tests...$(NC)"
	$(GOTEST) -v -race -cover ./...
	@echo "$(GREEN)✓ Tests passed$(NC)"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

## lint: Run linting
lint:
	@echo "$(YELLOW)Running linters...$(NC)"
{%- if values.testFramework == "standard" or values.testFramework == "testify" %}
	golangci-lint run ./...
{%- elif values.testFramework == "ginkgo" %}
	golangci-lint run --exclude-dirs=vendor ./...
{%- endif %}
	@echo "$(GREEN)✓ Linting passed$(NC)"

## format: Format code
format:
	@echo "$(YELLOW)Formatting code...$(NC)"
	gofmt -s -w .
	goimports -w .
	@echo "$(GREEN)✓ Code formatted$(NC)"

## clean: Remove build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Cleaned$(NC)"

## run: Run the CLI locally
run: build
	@./bin/$(BINARY_NAME) $(ARGS)

{%- if values.dockerize %}

## docker-build: Build Docker image
docker-build:
	@echo "$(YELLOW)Building Docker image...$(NC)"
	docker build -t ${{values.name}}:$(VERSION) .
	docker tag ${{values.name}}:$(VERSION) ${{values.name}}:latest
	@echo "$(GREEN)✓ Docker image built$(NC)"

## docker-run: Run Docker container
docker-run:
	@echo "$(YELLOW)Running Docker container...$(NC)"
	docker run --rm -it ${{values.name}}:latest $(ARGS)
{%- endif %}

{%- if values.releaseType == "goreleaser" %}

## release: Create release with GoReleaser
release:
	@echo "$(YELLOW)Creating release...$(NC)"
	goreleaser release --clean
	@echo "$(GREEN)✓ Release created$(NC)"

## release-snapshot: Create snapshot release
release-snapshot:
	@echo "$(YELLOW)Creating snapshot release...$(NC)"
	goreleaser release --snapshot --clean
	@echo "$(GREEN)✓ Snapshot created$(NC)"
{%- endif %}

## dev: Run in development mode with hot reload
dev:
	@echo "$(YELLOW)Starting development mode...$(NC)"
	air

## validate: Run all validation checks
validate: lint test
	@echo "$(GREEN)✓ All validation checks passed$(NC)"

## deps-update: Update dependencies
deps-update:
	@echo "$(YELLOW)Updating dependencies...$(NC)"
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

{%- if values.includeMocks %}

## mocks: Generate mocks
mocks:
	@echo "$(YELLOW)Generating mocks...$(NC)"
	mockery --all --dir internal --output internal/mocks
	@echo "$(GREEN)✓ Mocks generated$(NC)"
{%- endif %}
