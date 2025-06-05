.PHONY: help build test fmt lint clean run-example deps mod vet check install-tools

# Default target
help: ## Show help
	@echo "Samsung Wallet Go SDK - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build related
build: ## Build the project
	@echo "Building Samsung Wallet Go SDK..."
	go build ./...

build-example: ## Build example
	@echo "Building example..."
	go build -o bin/example ./examples

run-example: build-example ## Run example
	@echo "Running example..."
	./bin/example

# Test related
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -v -race ./...

benchmark: ## Run benchmark tests
	@echo "Running benchmarks..."
	go test -v -bench=. ./...

# Code quality
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

lint: install-tools ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run

check: fmt vet lint test ## Run all checks

# Dependency management
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

mod: ## Tidy go.mod
	@echo "Tidying go.mod..."
	go mod tidy

mod-verify: ## Verify module integrity
	@echo "Verifying modules..."
	go mod verify

# Tools installation
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean ./...

# Documentation generation
doc: ## Generate documentation
	@echo "Generating documentation..."
	go doc -all ./wallet

doc-server: ## Start documentation server
	@echo "Starting documentation server on :6060..."
	godoc -http=:6060

# Release related
tag: ## Create Git tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required. Usage: make tag VERSION=v1.0.0"; exit 1; fi
	@echo "Creating tag $(VERSION)..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

# Development
dev-setup: deps install-tools ## Setup development environment
	@echo "Development environment setup complete!"

watch: ## Watch file changes and auto-run tests (requires fswatch)
	@echo "Watching for file changes..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} make test; \
	else \
		echo "fswatch is not installed. Install it with: brew install fswatch"; \
	fi

# JWT generation helper (for development/testing)
generate-test-keys: ## Generate test RSA key pair
	@echo "Generating test RSA keys..."
	@mkdir -p test/keys
	openssl genrsa -out test/keys/private.pem 2048
	openssl rsa -in test/keys/private.pem -pubout -out test/keys/public.pem
	@echo "Keys generated in test/keys/"

# Debug
debug-env: ## Show environment information
	@echo "Go version:"
	@go version
	@echo ""
	@echo "Go environment:"
	@go env
	@echo ""
	@echo "Module info:"
	@go list -m

# Integration testing (requires actual API keys)
integration-test: ## Run integration tests (requires SAMSUNG_SERVICE_ID, SAMSUNG_PRIVATE_KEY, SAMSUNG_CERTIFICATE_ID environment variables)
	@if [ -z "$(SAMSUNG_SERVICE_ID)" ] || [ -z "$(SAMSUNG_PRIVATE_KEY)" ] || [ -z "$(SAMSUNG_CERTIFICATE_ID)" ]; then \
		echo "Integration test requires SAMSUNG_SERVICE_ID, SAMSUNG_PRIVATE_KEY, and SAMSUNG_CERTIFICATE_ID environment variables"; \
		exit 1; \
	fi
	@echo "Running integration tests..."
	go test -v -tags=integration ./...

# Profiling
profile-cpu: ## Generate CPU profile
	@echo "Running CPU profiling..."
	go test -cpuprofile=cpu.prof -bench=. ./...
	go tool pprof cpu.prof

profile-mem: ## Generate memory profile
	@echo "Running memory profiling..."
	go test -memprofile=mem.prof -bench=. ./...
	go tool pprof mem.prof 