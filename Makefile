GO := go
MOCKERY := mockery
HELM := helm
K3D := k3d
KUBECTL := kubectl

NAME := aegis
VER := 1.0.0
CMD_DIR := $(CURDIR)/cmd
BIN_DIR := $(CURDIR)/bin
HELM_DIR := $(CURDIR)/helm
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

.PHONY: create-cluster
create-cluster:
	@$(K3D) cluster create --config k3d-conf.yaml
	@$(K3D) image import $(NAME):$(VER) -c $(NAME)
	@$(HELM) dependency update "$(HELM_DIR)/$(NAME)"
	@$(HELM) install $(NAME) "$(HELM_DIR)/$(NAME)"
	@$(KUBECTL) get pods

.PHONY: delete-cluster
delete-cluster:
	@$(K3D) cluster delete $(NAME)
	@$(HELM) uninstall $(NAME)

.PHONY: cluster-ports
cluster-ports:
	@$(KUBECTL) port-forward svc/aegis-minio 9000:9000
	@$(KUBECTL) port-forward svc/aegis-postgresql 5432:5432
