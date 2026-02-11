.PHONY: run build clean test help

# Variables
BINARY_NAME=habit-server
GO=go
GOFLAGS=-v

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run: ## Run the application
	$(GO) run main.go

build: ## Build the application
	$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) main.go

clean: ## Remove build artifacts
	rm -rf bin/
	rm -rf logs/

test: ## Run tests
	$(GO) test -v ./...

deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

fmt: ## Format code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

lint: ## Run golangci-lint
	golangci-lint run

dev: ## Run in development mode with hot reload (requires air)
	air

docker-build: ## Build docker image
	docker build -t $(BINARY_NAME) .

docker-run: ## Run docker container
	docker run -p 8080:8080 $(BINARY_NAME)
