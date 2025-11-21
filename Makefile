.PHONY: help build test clean install-xk6 examples

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-xk6: ## Install xk6 tool
	go install go.k6.io/xk6/cmd/xk6@latest

build: ## Build k6 with xk6-parquet extension
	xk6 build --with github.com/mmga-lab/xk6-parquet=.

test: ## Run tests
	go test -v -race -coverprofile=coverage.txt ./...

test-verbose: ## Run tests with verbose output
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

coverage: test ## Generate and view coverage report
	go tool cover -html=coverage.txt

lint: ## Run linters
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Skipping."; \
	fi

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

mod-tidy: ## Tidy go modules
	go mod tidy

mod-download: ## Download go modules
	go mod download

clean: ## Clean build artifacts
	rm -f k6 coverage.txt
	rm -rf dist/

examples: build ## Build and run example
	@echo "Building k6 with xk6-parquet..."
	@./k6 version

docker-build: ## Build Docker image
	docker build -t k6-parquet:latest .

docker-run: docker-build ## Run k6 in Docker
	docker run --rm k6-parquet:latest version

generate-data: ## Generate sample Parquet data files
	@echo "Generating sample data..."
	cd examples/data && go run generate_sample.go

all: clean mod-download build test ## Run all: clean, download, build, test
