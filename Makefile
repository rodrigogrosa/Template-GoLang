.PHONY: help build run test lint clean docker-build docker-run swagger deps

# Variables
APP_NAME=template-service
DOCKER_IMAGE=$(APP_NAME):latest
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod download
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

swagger: ## Generate swagger documentation
	swag init -g cmd/api/main.go -o docs

build: swagger ## Build the application
	go build -o bin/$(APP_NAME) ./cmd/api

run: ## Run the application
	go run ./cmd/api/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	golangci-lint run --timeout=5m

fmt: ## Format code
	go fmt ./...
	gofmt -s -w $(GO_FILES)

vet: ## Run go vet
	go vet ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf docs/swagger.*
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env.example $(DOCKER_IMAGE)

docker-compose-up: ## Start services with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	docker-compose down

.DEFAULT_GOAL := help
