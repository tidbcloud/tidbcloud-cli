GOLANGCI_VERSION=v1.64.7
COVERAGE=coverage.out

.PHONY: deps
deps:  ## Download go module dependencies
	@echo "==> Installing go.mod dependencies..."
	go mod download
	go mod tidy

.PHONY: devtools
devtools:  ## Install dev tools
	@echo "==> Installing dev tools..."
	go install github.com/google/addlicense@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/google/go-licenses@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

.PHONY: setup
setup: deps devtools ## Set up dev env

.PHONY: generate-mocks
generate-mocks: ## Generate mock objects
	@echo "==> Generating mock objects"
	go install github.com/vektra/mockery/v2@v2.53.5
	mockery --name TiDBCloudClient --recursive --output=internal/mock --outpkg mock --filename api_client.go
	mockery --name EventsSender --recursive --output=internal/mock --outpkg mock --filename sender.go
	mockery --name Uploader --recursive --output=internal/mock --outpkg mock --filename uploader.go

.PHONY: addcopy
addcopy: ## Add copyright to all files
	@scripts/add-copy.sh

.PHONY: generate-v1beta1-client
generate-v1beta1-client: generate-v1beta1-serverless-client generate-v1beta1-iam-client generate-v1beta1-dedicated-client ## Generate v1beta1 client

.PHONY: generate-v1beta1-serverless-client
generate-v1beta1-serverless-client: install-openapi-generator ## Generate serverless client
	@echo "==> Generating serverless branch client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/branch
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/branch.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/branch --package-name branch -c go/config.yaml
	@echo "==> Generating serverless cluster client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/cluster
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/cluster.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/cluster --package-name cluster -c go/config.yaml
	@echo "==> Generating serverless export client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/export
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/export.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/export --package-name export -c go/config.yaml
	@echo "==> Generating serverless br client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/br
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/br.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/br --package-name br -c go/config.yaml
	@echo "==> Generating serverless import client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/imp
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/import.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/imp --package-name imp -c go/config.yaml
	@echo "==> Generating serverless audit log client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/auditlog
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/auditlog.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/auditlog --package-name auditlog -c go/config.yaml
	@echo "==> Generating serverless cdc client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/cdc
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/cdc.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/cdc --package-name cdc -c go/config.yaml
	cd pkg && go fmt ./tidbcloud/v1beta1/serverless/... && goimports -w .
	@echo "==> Generating serverless privatelink client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/privatelink
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/privatelink.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/privatelink --package-name privatelink -c go/config.yaml
	cd pkg && go fmt ./tidbcloud/v1beta1/serverless/... && goimports -w .

.PHONY: generate-v1beta1-dedicated-client
generate-v1beta1-dedicated-client: install-openapi-generator ## Generate dedicated client
	@echo "==> Generating dedicated client"
	rm -rf pkg/tidbcloud/v1beta1/dedicated
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false --skip-validate-spec -i ../../pkg/tidbcloud/v1beta1/dedicated.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/dedicated --package-name dedicated -c go/config.yaml
	cd pkg && go fmt ./tidbcloud/v1beta1/dedicated/... && goimports -w .

.PHONY: generate-v1beta1-iam-client
generate-v1beta1-iam-client: install-openapi-generator ## Generate iam client
	@echo "==> Generating iam client"
	rm -rf pkg/tidbcloud/v1beta1/iam
	cd tools/openapi-generator && npx openapi-generator-cli generate --inline-schema-options RESOLVE_INLINE_ENUMS=true --additional-properties=withGoMod=false,enumClassPrefix=true,disallowAdditionalPropertiesIfNotPresent=false --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/iam.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/iam --package-name iam -c go/config.yaml
	cd pkg && go fmt ./tidbcloud/v1beta1/iam/... && goimports -w .

.PHONY: install-openapi-generator
install-openapi-generator:
	cd tools/openapi-generator && npm install

.PHONY: fmt
fmt: ## Format changed go
	@scripts/fmt.sh

.PHONY: fix-lint
fix-lint: ## Fix linting errors
	golangci-lint run --fix

.PHONY: test
test: ## Run unit-tests
	@go test -race -cover -count=1 -coverprofile $(COVERAGE)  ./...

.PHONY: build
build: ## Generate a binary in ./bin
	@go build -o ./bin/ticloud ./cmd/ticloud

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate-docs
generate-docs: ## Generate docs objects
	@echo "==> Generating docs"
	go run gen_doc.go
