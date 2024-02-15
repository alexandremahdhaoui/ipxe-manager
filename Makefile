PROJECT := ipxe-api

GO_GEN := go generate
CONTROLLER_GEN := go run sigs.k8s.io/controller-tools/cmd/controller-gen@latest

OAPI_SPEC := ./api/ipxe.v1.yaml
OAPI_SERVER_PKG := server
OAPI_CLIENT_PKG := ipxeclient
OAPI_SERVER_FILE := ./internal/interface/server/zz_generated.server.go
OAPI_CLIENT_FILE := ./pkg/ipxeclient/zz_generated.ipxeclient.go

OAPI_CODEGEN := go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

.PHONY: generate
generate: ## Generate REST API server/client code, CRDs and other go generators.
	mkdir -p $$(dirname $(OAPI_SERVER_FILE)) $$(dirname $(OAPI_CLIENT_FILE))
	$(OAPI_CODEGEN) -generate types,server,spec -package $(OAPI_SERVER_PKG) -o $(OAPI_SERVER_FILE) $(OAPI_SPEC) || exit 1
	$(OAPI_CODEGEN) -generate types,client -package $(OAPI_CLIENT_PKG) -o $(OAPI_CLIENT_FILE) $(OAPI_SPEC) || exit 1

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
