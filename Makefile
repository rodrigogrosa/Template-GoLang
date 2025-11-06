.PHONY: help build test lint run docker-build docker-run clean install

# Variables
APP_NAME=golang-microservice
BINARY_NAME=api
DOCKER_IMAGE=$(APP_NAME):latest
GO_FILES=$(shell find . -name '*.go' -type f)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	@echo "Building application..."
	CGO_ENABLED=0 go build -o bin/$(BINARY_NAME) ./cmd/api

run: ## Run the application
	@echo "Running application..."
	go run ./cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test -v -tags=integration ./...

lint: ## Run linter
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 -p 9090:9090 --rm $(DOCKER_IMAGE)

docker-compose-up: ## Start services with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	docker-compose down

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

swagger: ## Generate swagger documentation
	@echo "Generating Swagger docs..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/api/main.go

mod-update: ## Update Go modules
	go get -u ./...
	go mod tidy

security: ## Run security scan
	@echo "Running security scan..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	gosec ./...

check: fmt vet lint test ## Run all checks

all: clean install lint test build ## Run all tasks
