# Makefile for Kagi MCP Server

# Variables
BINARY_NAME=kagimcp
GO=go
GO_BUILD=$(GO) build
GO_TEST=$(GO) test
GO_CLEAN=$(GO) clean
GO_LINT=golangci-lint
GO_MOD=$(GO) mod
ENV_FILE=.env
GO_FILES=$(wildcard *.go)

# Default target
.PHONY: all
all: help

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GO_BUILD) -o $(BINARY_NAME) .

# Run the server in stdio mode
.PHONY: start
start:
	@if [ ! -f $(BINARY_NAME) ]; then \
		echo "Binary not found. Building..."; \
		make build; \
	fi
	@if [ ! -f $(ENV_FILE) ] && [ -z "$$KAGI_API_KEY" ]; then \
		echo "Warning: No .env file found and KAGI_API_KEY environment variable not set."; \
		echo "You may need to set your API key with -api-key flag."; \
	fi
	@echo "Starting $(BINARY_NAME) in stdio mode..."
	@./$(BINARY_NAME) -t stdio

# Start the server in SSE mode
.PHONY: start-sse
start-sse:
	@if [ ! -f $(BINARY_NAME) ]; then \
		echo "Binary not found. Building..."; \
		make build; \
	fi
	@if [ ! -f $(ENV_FILE) ] && [ -z "$$KAGI_API_KEY" ]; then \
		echo "Warning: No .env file found and KAGI_API_KEY environment variable not set."; \
		echo "You may need to set your API key with -api-key flag."; \
	fi
	@echo "Starting $(BINARY_NAME) in SSE mode..."
	@./$(BINARY_NAME) -t sse

# Initialize the project
.PHONY: init
init:
	@echo "Initializing the project..."
	@$(GO_MOD) download
	@$(GO_MOD) tidy
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Creating .env file template..."; \
		echo "# Add your Kagi API key here" > $(ENV_FILE); \
		echo "KAGI_API_KEY=" >> $(ENV_FILE); \
		echo ".env file created. Please add your Kagi API key."; \
	fi
	@echo "Initialization complete!"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(GO_TEST) -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@$(GO_CLEAN)
	@rm -f $(BINARY_NAME)
	@echo "Cleaned!"

# Show help
.PHONY: help
help:
	@echo "Kagi MCP Server - Makefile Help"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build the application"
	@echo "  start        Build (if needed) and run the server in stdio mode"
	@echo "  start-sse    Build (if needed) and run the server in SSE mode"
	@echo "  init         Initialize the project (download dependencies and create .env template)"
	@echo "  test         Run tests"
	@echo "  clean        Remove binary and build artifacts"
	@echo "  help         Show this help message"
