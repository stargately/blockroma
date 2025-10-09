# Makefile for Stellar RPC Indexer

.PHONY: help build test test-verbose test-coverage clean run docker-build lint fmt vet install-tools

# Default target
.DEFAULT_GOAL := help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary name
BINARY_NAME=indexer
BINARY_PATH=./cmd/indexer

# Build output directory
BUILD_DIR=./build

help: ## Show this help message
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the indexer binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-short: ## Run tests with short flag
	@echo "Running short tests..."
	$(GOTEST) -short -v ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "Coverage report generated: coverage.out"
	@echo "To view coverage in browser, run: go tool cover -html=coverage.out"

test-coverage-html: test-coverage ## Generate and open HTML coverage report
	$(GOCMD) tool cover -html=coverage.out

test-package: ## Run tests for a specific package (usage: make test-package PKG=./pkg/parser)
	@echo "Running tests for $(PKG)..."
	$(GOTEST) -v $(PKG)

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

clean: ## Clean build artifacts and test cache
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out
	@echo "Clean complete"

run: ## Run the indexer
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(BINARY_PATH)

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t stellar-indexer:latest .

docker-run: ## Run indexer in Docker
	@echo "Running indexer in Docker..."
	docker run --rm \
		-e STELLAR_RPC_URL=http://host.docker.internal:8000 \
		-e POSTGRES_DSN="postgresql://stellar:password@host.docker.internal:5432/stellar_indexer?sslmode=disable" \
		stellar-indexer:latest

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

lint: install-tools ## Run golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run ./...

tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	$(GOMOD) tidy

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download

verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	$(GOMOD) verify

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@which golangci-lint > /dev/null || \
		(echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin)

check: fmt vet test ## Run fmt, vet, and tests

ci: deps verify check test-coverage ## Run CI pipeline (deps, verify, fmt, vet, test with coverage)
	@echo "CI checks passed!"

pre-commit: fmt vet ## Run pre-commit checks (fmt, vet)
	@echo "Pre-commit checks passed!"

all: clean deps build test ## Clean, download deps, build, and test

# Development shortcuts
dev-setup: deps install-tools ## Setup development environment
	@echo "Development environment ready!"

dev-test: ## Run tests in watch mode (requires entr)
	@which entr > /dev/null || (echo "Install entr first: brew install entr (macOS) or apt-get install entr (Linux)" && exit 1)
	@echo "Watching for changes... (Ctrl+C to stop)"
	@find . -name '*.go' | entr -c make test-short

# Database tools
db-migrate: ## Run database migrations (requires running indexer)
	@echo "Migrations run automatically on indexer startup"

db-reset: ## Warning: Reset database (requires psql)
	@echo "Resetting database..."
	@echo "This will drop and recreate the database. Press Ctrl+C to cancel, or Enter to continue..."
	@read _
	psql "$(POSTGRES_DSN)" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "Database reset complete. Restart indexer to recreate tables."

# Info targets
version: ## Show Go version
	@$(GOCMD) version

info: ## Show project info
	@echo "Project: Stellar RPC Indexer"
	@echo "Go version: $$($(GOCMD) version)"
	@echo "Binary: $(BUILD_DIR)/$(BINARY_NAME)"
	@echo ""
	@echo "Quick commands:"
	@echo "  make build       - Build the binary"
	@echo "  make test        - Run tests"
	@echo "  make run         - Run the indexer"
	@echo "  make help        - Show all available commands"
