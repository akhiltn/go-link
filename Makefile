BUILD_DIR := bin
EXECUTABLE= quick-url
WINDOWS=$(BUILD_DIR)/$(EXECUTABLE)_windows_amd64.exe
LINUX=$(BUILD_DIR)/$(EXECUTABLE)_linux_amd64
VERSION=$(shell git describe --tags --always --long --dirty)

.PHONY: all test clean

all: clean build ## Build and run tests

test: ## Run unit tests
	go test ./...

build: windows linux ## Build binaries
	@echo version: $(VERSION)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -v -o $(WINDOWS) -ldflags="-X main.version=$(VERSION)"  .

$(LINUX):
	env GOOS=linux GOARCH=amd64  go build -v -o $(LINUX) -ldflags="-s -w -X main.version=$(VERSION)"  .

clean: ## Remove previous build
	rm -rf $(BUILD_DIR)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
