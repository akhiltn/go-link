CONTAINER_ENGINE ?= podman
IMAGE ?= go-link
DEV_IMAGE ?= $(IMAGE):dev
DATA_VOLUME ?= golink-data
DOMAIN ?= localhost
PORT ?= 3000
AIR_BIN ?= $(shell go env GOPATH)/bin/air
VERSION := $(shell git describe --tags --always --long --dirty)

.PHONY: all test container-test image run air dev-image dev clean help

all: test image ## Run local tests and build the container image

test: ## Run local Go tests
	go test ./...

container-test: ## Run tests in the Dockerfile test stage
	$(CONTAINER_ENGINE) build --target test .

image: ## Build the production container image
	$(CONTAINER_ENGINE) build --build-arg VERSION=$(VERSION) -t $(IMAGE) .

run: image ## Run the production container image
	$(CONTAINER_ENGINE) run --rm -p $(PORT):3000 -e DOMAIN=$(DOMAIN) \
		-v $(DATA_VOLUME):/data $(IMAGE)

air: ## Build a static Air binary for the development container
	CGO_ENABLED=0 go install github.com/air-verse/air@latest

dev-image: ## Build the hot-reload development container image
	$(CONTAINER_ENGINE) build --target dev -t $(DEV_IMAGE) .

dev: air dev-image ## Run the hot-reload development container
	$(CONTAINER_ENGINE) run --rm -p $(PORT):3000 -e DOMAIN=$(DOMAIN) \
		-v "$(CURDIR)":/src -v $(DATA_VOLUME):/data \
		-v "$(AIR_BIN)":/usr/local/bin/air:ro $(DEV_IMAGE)

clean: ## Remove local build artifacts
	rm -rf bin

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
