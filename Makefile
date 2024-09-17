EXECUTABLE=go-link
BUILD_DIR := bin
COVERAGE_DIR := $(BUILD_DIR)/coverage
COVERAGE_FILE := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html
WINDOWS=$(BUILD_DIR)/$(EXECUTABLE)_windows_amd64.exe
LINUX=$(BUILD_DIR)/$(EXECUTABLE)_linux_amd64
VERSION=$(shell git describe --tags --always --long --dirty)

.PHONY: all test clean

all: clean build test ## Build and run tests

test: ## Run unit tests with coverage
	mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated at $(COVERAGE_HTML)"

build: windows linux ## Build binaries
	@echo version: $(VERSION)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -v -o $(WINDOWS) -ldflags="-X main.version=$(VERSION)" ./cmd/go-link

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)" ./cmd/go-link

clean: ## Remove previous build
	rm -rf $(BUILD_DIR)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
