# ------------------------------------------------------- ENVS ------------------------------------------------------- #

PROJECT    := ipxer

COMMIT_SHA := $(shell git rev-parse --short HEAD)
TIMESTAMP  := $(shell date --utc --iso-8601=seconds)
VERSION    ?= $(shell git describe --tags --always --dirty)

CHARTS     := $(shell ./hack/list-subprojects.sh charts)
CONTAINERS := $(shell ./hack/list-subprojects.sh containers)
CMDS       := $(shell ./hack/list-subprojects.sh cmd)

GO_BUILD_LDFLAGS ?= "-X main.BuildTimestamp=$(TIMESTAMP) -X main.CommitSHA=$(COMMIT_SHA) -X main.Version=$(VERSION)"

# ------------------------------------------------------- VERSIONS --------------------------------------------------- #

# renovate: datasource=github-release depName=kubernetes-sigs/controller-tools
CONTROLLER_GEN_VERSION := v0.14.0
# renovate: datasource=github-release depName=mvdan/gofumpt
GOFUMPT_VERSION        := v0.6.0
# renovate: datasource=github-release depName=golangci/golangci-lint
GOLANGCI_LINT_VERSION  := v1.59.1
# renovate: datasource=github-release depName=gotestyourself/gotestsum
GOTESTSUM_VERSION      := v1.12.0
# renovate: datasource=github-release depName=vektra/mockery
MOCKERY_VERSION        := v2.42.0
# renovate: datasource=github-release depName=deepmap/oapi-codegen
OAPI_CODEGEN_VERSION   := v2.1.0

# ------------------------------------------------------- TOOLS ------------------------------------------------------ #

CONTAINER_ENGINE ?= podman

CONTROLLER_GEN := go run sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)
GO_GEN         := go generate
GOFUMPT        := go run mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
GOLANGCI_LINT  := go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
GOTESTSUM      := go run gotest.tools/gotestsum@$(GOTESTSUM_VERSION) --format pkgname
MOCKERY        := go run github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
OAPI_CODEGEN   := go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION)

CLEAN_MOCKS    := rm -rf ./internal/util/mocks

# ------------------------------------------------------- GENERATE --------------------------------------------------- #

# -- ipxer
# spec
OAPI_IPXER_SPEC        := ./api/ipxer.v1.yaml
# client
OAPI_IPXER_CLIENT_PKG  := ipxerclient
OAPI_IPXER_CLIENT_FILE := ./pkg/generated/ipxerclient/zz_generated.ipxerclient.go
# server
OAPI_IPXER_SERVER_PKG  := server
OAPI_IPXER_SERVER_FILE := ./internal/driver/server/zz_generated.server.go

# -- ipxer-webhook-resolver
# spec
OAPI_WEBHOOK_RESOLVER_SPEC        := ./api/ipxer-webhook-resolver.v1.yaml
# client
OAPI_WEBHOOK_RESOLVER_CLIENT_PKG  := resolverclient
OAPI_WEBHOOK_RESOLVER_CLIENT_FILE := ./pkg/generated/resolverclient/zz_generated.resolverclient.go
# server
OAPI_WEBHOOK_RESOLVER_SERVER_PKG  := resolverserver
OAPI_WEBHOOK_RESOLVER_SERVER_FILE := ./pkg/generated/resolverserver/zz_generated.resolverserver.go

# -- ipxer-webhook-transformer
# spec
OAPI_WEBHOOK_TRANSFORMER_SPEC        := ./api/ipxer-webhook-transformer.v1.yaml
# client
OAPI_WEBHOOK_TRANSFORMER_CLIENT_PKG  := transformerclient
OAPI_WEBHOOK_TRANSFORMER_CLIENT_FILE := ./pkg/generated/transformerclient/zz_generated.transformerclient.go
# server
OAPI_WEBHOOK_TRANSFORMER_SERVER_PKG  := transformerserver
OAPI_WEBHOOK_TRANSFORMER_SERVER_FILE := ./pkg/generated/transformerserver/zz_generated.transformerserver.go

# TODO: simplify the OAPI code generation. Use a yaml config file and create a script.

.PHONY: generate
generate: ## Generate REST API server/client code, CRDs and other go generators.
	mkdir -p $$(dirname $(OAPI_IPXER_SERVER_FILE)) $$(dirname $(OAPI_IPXER_CLIENT_FILE))
	$(OAPI_CODEGEN) -generate types,server,spec -package $(OAPI_IPXER_SERVER_PKG) -o $(OAPI_IPXER_SERVER_FILE) $(OAPI_IPXER_SPEC) || exit 1
	$(OAPI_CODEGEN) -generate types,client -package $(OAPI_IPXER_CLIENT_PKG) -o $(OAPI_IPXER_CLIENT_FILE) $(OAPI_IPXER_SPEC) || exit 1

	mkdir -p $$(dirname $(OAPI_WEBHOOK_RESOLVER_SERVER_FILE)) $$(dirname $(OAPI_WEBHOOK_RESOLVER_CLIENT_FILE))
	$(OAPI_CODEGEN) -generate types,server,spec -package $(OAPI_WEBHOOK_RESOLVER_SERVER_PKG) -o $(OAPI_WEBHOOK_RESOLVER_SERVER_FILE) $(OAPI_WEBHOOK_RESOLVER_SPEC) || exit 1
	$(OAPI_CODEGEN) -generate types,client -package $(OAPI_WEBHOOK_RESOLVER_CLIENT_PKG) -o $(OAPI_WEBHOOK_RESOLVER_CLIENT_FILE) $(OAPI_WEBHOOK_RESOLVER_SPEC) || exit 1

	mkdir -p $$(dirname $(OAPI_WEBHOOK_TRANSFORMER_SERVER_FILE)) $$(dirname $(OAPI_WEBHOOK_TRANSFORMER_CLIENT_FILE))
	$(OAPI_CODEGEN) -generate types,server,spec -package $(OAPI_WEBHOOK_TRANSFORMER_SERVER_PKG) -o $(OAPI_WEBHOOK_TRANSFORMER_SERVER_FILE) $(OAPI_WEBHOOK_TRANSFORMER_SPEC) || exit 1
	$(OAPI_CODEGEN) -generate types,client -package $(OAPI_WEBHOOK_TRANSFORMER_CLIENT_PKG) -o $(OAPI_WEBHOOK_TRANSFORMER_CLIENT_FILE) $(OAPI_WEBHOOK_TRANSFORMER_SPEC) || exit 1

	$(GO_GEN) "./..."

	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	$(CONTROLLER_GEN) paths="./..." \
		crd:generateEmbeddedObjectMeta=true \
		output:crd:artifacts:config=charts/$(PROJECT)/templates/crds

	$(CONTROLLER_GEN) paths="./..." \
		rbac:roleName=$(PROJECT) \
		webhook \
		output:rbac:dir=charts/$(PROJECT)/templates/rbac \
		output:webhook:dir=charts/$(PROJECT)/templates/webhook

	$(CLEAN_MOCKS)
	$(MOCKERY)

# ------------------------------------------------------- BUILD BINARIES --------------------------------------------- #

.PHONY: build-binary
build-binary: generate
	GO_BUILD_LDFLAGS=$(GO_BUILD_LDFLAGS) ./hack/build-binary.sh "${BINARY_NAME}"

.PHONY: build-binaries
build-binaries: generate ## Build the binaries.
	echo $(CMDS) | \
		GO_BUILD_LDFLAGS=$(GO_BUILD_LDFLAGS) \
		xargs -n1 ./hack/build-binary.sh

# ------------------------------------------------------- BUILD CONTAINERS -------------------------------------------- #

.PHONY: build-container
build-container: generate
	CONTAINER_ENGINE=$(CONTAINER_ENGINE) GO_BUILD_LDFLAGS=$(GO_BUILD_LDFLAGS) VERSION=$(VERSION) \
		./hack/build-container.sh "${CONTAINER_NAME}"

.PHONY: build-containers
build-containers: generate
	echo $(CONTAINERS) | \
		CONTAINER_ENGINE=$(CONTAINER_ENGINE) \
		GO_BUILD_LDFLAGS=$(GO_BUILD_LDFLAGS) \
		VERSION=$(VERSION) \
		xargs -n1 ./hack/build-container.sh

# ------------------------------------------------------- FMT -------------------------------------------------------- #

.PHONY: fmt
fmt:
	$(GOFUMPT) -w .

# ------------------------------------------------------- LINT ------------------------------------------------------- #

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run --fix

# ------------------------------------------------------- TEST ------------------------------------------------------- #

.PHONY: test-unit
test-unit:
	$(GOTESTSUM) --junitfile .ignore.test-unit.xml -- -tags unit -race ./... -count=1 -short -cover -coverprofile .ignore.test-unit-coverage.out ./...

.PHONY: test-integration
test-integration:
	$(GOTESTSUM) --junitfile .ignore.test-integration.xml -- -tags integration -race ./... -count=1 -short -cover -coverprofile .ignore.test-integration-coverage.out ./...

.PHONY: test-functional
test-functional:
	$(GOTESTSUM) --junitfile .ignore.test-functional.xml -- -tags functional -race ./... -count=1 -short -cover -coverprofile .ignore.test-functional-coverage.out ./...

.PHONY: test-e2e
test-e2e:
	echo TODO: test-e2e

.PHONY: test
test: test-unit test-integration test-functional

# ------------------------------------------------------- PRE-PUSH --------------------------------------------------- #

.PHONY: githooks
githooks: ## Initializes Git hooks to run before a push.
	git config core.hooksPath .githooks

.PHONY: pre-push
pre-push: generate fmt lint test
	git status --porcelain