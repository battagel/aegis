GO := go
MOCKERY := mockery
HELM := helm
K3D := k3d
K3D_CONF := k3d-conf.yaml
KUBECTL := kubectl
DOCKER := docker

NAME := aegis
VER := latest
TOPIC := minio-put-events
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
build: vendor
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

## test: Run the tests
.PHONY: test
test:
	@$(GO) test -v ./... --cover

## mock: Generate the mocks for testing
.PHONY: mock
mock:
	@$(MOCKERY) --dir ./internal -r --all --config .mockery.yaml

## docker-build: Build the docker image
.PHONY: docker-build
docker-build:
	@$(DOCKER) build . -t $(NAME):$(VER)

## create-cluster: Create the k3d cluster
.PHONY: create-cluster
create-cluster:
	@$(K3D) cluster create --config $(K3D_CONF)
	@$(K3D) image import $(NAME):$(VER) -c $(NAME)
	@$(HELM) dependency update "$(HELM_DIR)/$(NAME)"
	@$(HELM) install $(NAME) "$(HELM_DIR)/$(NAME)"
	@$(KUBECTL) get pods -w

## delete-cluster: Delete the k3d cluster
.PHONY: delete-cluster
delete-cluster:
	-@$(K3D) cluster delete $(NAME)
	-@$(HELM) uninstall $(NAME)

## rebuild-cluster: Delete and recreate the cluster
.PHONY: rebuild-cluster
rebuild-cluster: delete-cluster docker-build create-cluster

## purge-topic: Deletes and recreates the topic to remove dead messages
.PHONY: purge-topic
purge-topic:
	kafka-topics --delete --topic $(TOPIC) --bootstrap-server localhost:9092
	kafka-topics --create --topic $(TOPIC) --bootstrap-server localhost:9092
