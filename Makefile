.PHONY: build run test clean docker-build docker-run docker-stop docs

# Build the application
build:
	go build -o otp-auth-service .

# Run the application locally
run:
	go run main.go

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Run tests for specific package
test-package:
	@echo "Usage: make test-package PKG=package_name"
	@echo "Example: make test-package PKG=./internal/validation"
	@if [ -z "$(PKG)" ]; then echo "Please specify PKG parameter"; exit 1; fi
	go test $(PKG) -v

# Run tests with race detection
test-race:
	go test -race ./...

# Clean build artifacts
clean:
	rm -f otp-auth-service
	rm -rf docs/

# Generate Swagger documentation
docs:
	swag init -g main.go

# Build Docker image
docker-build:
	docker build -t otp-auth-service .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose services
docker-stop:
	docker-compose down

# View logs
logs:
	docker-compose logs -f app

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# All-in-one setup
setup: deps docs build

# Development mode (with hot reload)
dev:
	air

# Show help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-package  - Run tests for specific package"
	@echo "  test-race     - Run tests with race detection"
	@echo "  clean        - Clean build artifacts"
	@echo "  docs         - Generate Swagger documentation"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker Compose services"
	@echo "  logs         - View application logs"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  setup        - Complete setup (deps + docs + build)"
	@echo "  help         - Show this help"
