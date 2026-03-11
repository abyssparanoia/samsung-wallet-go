.PHONY: build test test.coverage test.race lint.go lint.go.fix fmt vet mod.tidy check clean

build: ## Build the project
	go build ./...

test: ## Run tests
	go test ./...

test.coverage: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test.race: ## Run tests with race detector
	go test -race ./...

lint.go: ## Run golangci-lint
	go tool golangci-lint run

lint.go.fix: ## Run golangci-lint with auto-fix
	go tool golangci-lint run --fix

fmt: ## Format code (goimports + gofumpt via golangci-lint)
	go tool golangci-lint fmt

vet: ## Run go vet
	go vet ./...

mod.tidy: ## Tidy go.mod
	go mod tidy

check: build vet lint.go test ## Run all checks (build, vet, lint, test)

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean ./...
