GOLANGCI_VERSION=v1.56.2
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
	go install github.com/vektra/mockery/v2@v2.43.0
	mockery --name TiDBCloudClient --recursive --output=internal/mock --outpkg mock --filename api_client.go
	mockery --name EventsSender --recursive --output=internal/mock --outpkg mock --filename sender.go
	mockery --name Uploader --recursive --output=internal/mock --outpkg mock --filename uploader.go

# Required to install go-swagger https://goswagger.io/install.html
.PHONY: generate-import-client
generate-import-client: ## Generate import client
	@echo "==> Generating import client"
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	swagger generate client -f pkg/tidbcloud/import/import-api.json -A tidbcloud-import -t pkg/tidbcloud/import

.PHONY: generate-pingchat-client
generate-pingchat-client: ## Generate PingChat client
	@echo "==> Generating PingChat client"
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	swagger generate client -f pkg/tidbcloud/pingchat/pingchat_swagger.json -A tidbcloud-pingchat -t pkg/tidbcloud/pingchat

.PHONY: addcopy
addcopy: ## Add copyright to all files
	@scripts/add-copy.sh

.PHONY: generate-v1beta1-client
generate-v1beta1-client: install-openapi-generator ## Generate v1beta1 client
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	@echo "==> Generating serverless branch client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/branch
	cd tools/openapi-generator && npx openapi-generator-cli generate --additional-properties=withGoMod=false,enumClassPrefix=true --global-property=apiTests=false,apiDocs=false,modelDocs=fasle,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/branch.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/branch --package-name branch
	@echo "==> Generating serverless client"
	swagger generate client -f pkg/tidbcloud/v1beta1/serverless/serverless.swagger.json -A tidbcloud-serverless -t pkg/tidbcloud/v1beta1/serverless
	@echo "==> Generating serverless br client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/br
	cd tools/openapi-generator && npx openapi-generator-cli generate --additional-properties=withGoMod=false,enumClassPrefix=true --global-property=apiTests=false,apiDocs=false,modelDocs=false,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/br.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/br --package-name br
	@echo "==> Generating serverless import client"
	swagger generate client -f pkg/tidbcloud/v1beta1/serverless_import/import.swagger.json -A tidbcloud-serverless -t pkg/tidbcloud/v1beta1/serverless_import
	@echo "==> Generating serverless export client"
	rm -rf pkg/tidbcloud/v1beta1/serverless/export
	cd tools/openapi-generator && npx openapi-generator-cli generate --additional-properties=withGoMod=false,enumClassPrefix=true --global-property=apiTests=false,apiDocs=false,modelDocs=fasle,modelTests=false -i ../../pkg/tidbcloud/v1beta1/serverless/export.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/serverless/export --package-name export
	@echo "==> Generating iam client"
	rm -rf pkg/tidbcloud/v1beta1/iam
	cd tools/openapi-generator && npx openapi-generator-cli generate --additional-properties=withGoMod=false,enumClassPrefix=true --global-property=apiTests=false,apiDocs=false,modelDocs=fasle,modelTests=false -i ../../pkg/tidbcloud/v1beta1/iam.swagger.json -g go -o ../../pkg/tidbcloud/v1beta1/iam --package-name iam
	go fmt ./pkg/...

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
generate-docs: ## Generate mock objects
	@echo "==> Generating docs"
	go run gen_doc.go
