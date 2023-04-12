GO := go
MOCKERY := mockery
HELM := helm

NAME := aegis
CMD_DIR := $(CURDIR)/cmd
BIN_DIR := $(CURDIR)/bin
MAIN_LOCATION := $(CMD_DIR)/$(NAME)/main.go

## help: Print this message
.PHONY: help
help:
	@fgrep -h '##' $(MAKEFILE_LIST) | fgrep -v fgrep | column -t -s ':' | sed -e 's/## //'

## build: Create the binary
.PHONY: build
build:
	@$(GO) build -o $(BIN_DIR)/$(NAME) -mod=vendor $(MAIN_LOCATION)

## run: Run the binary
.PHONY: run
run:
	@$(BIN_DIR)/$(NAME)

## vendor: Download the vendored dependencies
.PHONY: vendor
vendor:
	@$(GO) mod tidy
	@$(GO) mod vendor

.PHONY: test
test:
	@$(GO) test -v ./... --cover

.PHONY: mock
mock:
	@$(MOCKERY) --dir ./internal -r --all --config .mockery.yaml


.PHONY: helm-preview
helm-preview:
	@$(HELM) template aegis helm -f ./helm/values.yaml

.PHONY: helm-install
helm-install:
	@$(HELM) install aegis helm -f ./helm/values.yaml
