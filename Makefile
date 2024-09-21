# Makefile for Docker Dogger

# Variables
BINARY_NAME = colker
CONFIG_FILE = config.yaml

# Go commands
GO_BUILD = go build
GO_CLEAN = go clean
GO_FMT = go fmt
GO_RUN = go run

# Default target: build
all: build

# Build the Go binary
build:
	@echo "Building the Go binary..."
	$(GO_BUILD) -o $(BINARY_NAME) main.go

# Run the application
run: build
	@echo "Running the application..."
	./$(BINARY_NAME)

# Format the code
fmt:
	@echo "Formatting the Go code..."
	$(GO_FMT) ./...

# Clean up the build
clean:
	@echo "Cleaning up..."
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)

# Install dependencies (if any, not necessary for stdlib dependencies)
deps:
	@echo "Installing Go dependencies..."
	go mod tidy

# Rebuild the application
rebuild: clean build

.PHONY: all build run clean fmt deps rebuild

