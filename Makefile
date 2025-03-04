.PHONY: build clean test run lint deps

# Binary name
BINARY_NAME=scrutiny

# Go build and run parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOLINT=golangci-lint

# Build flags
BUILD_DIR=build
LDFLAGS=-ldflags "-w -s"

# Sources and packages
SRC_DIRS=cmd internal pkg
PACKAGES=$(shell go list ./... | grep -v /vendor/)

# Set the build target
build:
	@echo "Building..."
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cnapp

# Run the application
run:
	@echo "Running..."
	$(GORUN) ./cmd/cnapp/main.go

# Clean the build directory
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Testing..."
	$(GOTEST) -v $(PACKAGES)

# Run tests with coverage
test-coverage:
	@echo "Testing with coverage..."
	$(GOTEST) -coverprofile=coverage.out $(PACKAGES)
	go tool cover -html=coverage.out

# Lint the code
lint:
	@echo "Linting..."
	$(GOLINT) run

# Get dependencies
deps:
	@echo "Getting dependencies..."
	$(GOGET) -v -t -d ./...
	go mod tidy

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/scrutiny
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/scrutiny
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/scrutiny