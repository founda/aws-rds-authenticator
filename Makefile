# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# Deploy First Mentality

# ==============================================================================
# Install Tooling and Dependencies
#
#	If you are running a mac machine with brew, run these commands:
#	$ make dev-brew  or  make dev-brew-arm64
#	$ make dev-docker
#	$ make dev-gotooling
#
#	If you are running a linux machine with brew, run these commands:
#	$ make dev-brew-common
#	$ make dev-docker
#	$ make dev-gotooling
#   Follow instructions above for Telepresence.
#
#	If you are a windows user with brew, run these commands:
#	$ make dev-brew-common
#	$ make dev-docker
#	$ make dev-gotooling
#   Follow instructions above for Telepresence.

# ==============================================================================
# Define dependencies

GOLANG          := golang:1.20
ALPINE          := alpine:3.18
KIND            := kindest/node:v1.27.1
POSTGRES        := postgres:15.3
TELEPRESENCE    := datawire/tel2:2.13.2

KIND_CLUSTER    := founda-rds-auth-cluster
NAMESPACE       := rds
APP             := aws-rds-authenticator
BASE_IMAGE_NAME := ghcr.io/founda
SERVICE_NAME    := aws-rds-authenticator
VERSION         := 0.0.1
CLI_IMAGE       := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-brew-common:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli

dev-brew: dev-brew-common
	brew list datawire/blackbird/telepresence || brew install datawire/blackbird/telepresence

dev-brew-arm64: dev-brew-common
	brew list datawire/blackbird/telepresence-arm64 || brew install datawire/blackbird/telepresence-arm64

dev-docker:
	docker pull $(GOLANG)
	docker pull $(ALPINE)
	docker pull $(KIND)
	docker pull $(POSTGRES)
	docker pull $(TELEPRESENCE)

# ==============================================================================
# Building containers

all: cli

cli:
	docker build \
		-f zarf/docker/Dockerfile  \
		-t $(CLI_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		.

# ==============================================================================
# Running from within k8s/kind

dev-up-local:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	kind load docker-image $(TELEPRESENCE) --name $(KIND_CLUSTER)
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-up: dev-up-local
	telepresence --context=kind-$(KIND_CLUSTER) helm install
	telepresence --context=kind-$(KIND_CLUSTER) connect

dev-down-local:
	kind delete cluster --name $(KIND_CLUSTER)

dev-down:
	telepresence quit -s
	kind delete cluster --name $(KIND_CLUSTER)

# ------------------------------------------------------------------------------

dev-apply:
	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database
