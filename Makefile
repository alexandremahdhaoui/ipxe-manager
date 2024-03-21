PROJECT := ipxer

OAPI_IPXER_SPEC        := ./api/ipxer.v1.yaml
OAPI_IPXER_SERVER_PKG  := server
OAPI_IPXER_SERVER_FILE := ./internal/drivers/server/zz_generated.server.go
OAPI_IPXER_CLIENT_PKG  := ipxerclient
OAPI_IPXER_CLIENT_FILE := ./pkg/ipxerclient/zz_generated.ipxerclient.go

OAPI_WEBHOOK_RESOLVER_SPEC        := ./api/ipxer-webhook-resolver.v1.yaml
OAPI_WEBHOOK_RESOLVER_CLIENT_PKG  := resolverclient
OAPI_WEBHOOK_RESOLVER_CLIENT_FILE := ./pkg/resolverclient/zz_generated.resolverclient.go
OAPI_WEBHOOK_RESOLVER_SERVER_PKG  := resolverserver
OAPI_WEBHOOK_RESOLVER_SERVER_FILE := ./pkg/resolverserver/zz_generated.resolverserver.go

OAPI_WEBHOOK_TRANSFORMER_SPEC        := ./api/ipxer-webhook-transformer.v1.yaml
OAPI_WEBHOOK_TRANSFORMER_CLIENT_PKG  := transformerclient
OAPI_WEBHOOK_TRANSFORMER_CLIENT_FILE := ./pkg/transformerclient/zz_generated.transformerclient.go
OAPI_WEBHOOK_TRANSFORMER_SERVER_PKG  := transformerserver
OAPI_WEBHOOK_TRANSFORMER_SERVER_FILE := ./pkg/transformerserver/zz_generated.transformerserver.go

GO_GEN         := go generate
CONTROLLER_GEN := go run sigs.k8s.io/controller-tools/cmd/controller-gen@latest
OAPI_CODEGEN   := go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

MOCKS_CLEAN   := rm -rf ./internal/util/mocks
MOCKS_INSTALL := go install github.com/vektra/mockery/v2@v2.42.0
MOCKS_GEN     := mockery

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
		rbac:roleName=$(PROJCET) \
		webhook \
		output:rbac:dir=charts/$(PROJECT)/templates/rbac \
		output:webhook:dir=charts/$(PROJECT)/templates/webhook

	$(MOCKS_CLEAN)
	$(MOCKS_INSTALL)
	$(MOCKS_GEN)
