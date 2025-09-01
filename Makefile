.PHONY: build run test clean

# Binary name
BINARY_NAME=mastercom-service

# Build directory
BUILD_DIR=bin

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run ./cmd/server

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	@swag init -g cmd/server/main.go

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME)

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  docs         - Generate API documentation"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help"

# Start services with Docker Compose (includes Datadog agent)
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up --build

# Stop services
docker-down:
	@echo "Stopping services..."
	@docker-compose down

# Start services in background
docker-up-d:
	@echo "Starting services in background..."
	@docker-compose up -d --build

# View logs
docker-logs:
	@echo "Viewing service logs..."
	@docker-compose logs -f

# View Datadog agent logs
datadog-logs:
	@echo "Viewing Datadog agent logs..."
	@docker-compose logs -f datadog-agent

# Test Datadog agent connectivity
datadog-test:
	@echo "Testing Datadog agent connectivity..."
	@curl -s http://localhost:8126/info || echo "Datadog agent not accessible"
