GOLANGCI_VERSION=v1.55.2
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
	go install github.com/vektra/mockery/v2@latest
	mockery --name TiDBCloudClient --recursive --output=internal/mock --outpkg mock --filename api_client.go
	mockery --name EventsSender --recursive --output=internal/mock --outpkg mock --filename sender.go
	mockery --name MySQLHelper --recursive --output=internal/mock --outpkg mock --filename mysql_helper.go

# Required to install go-swagger https://goswagger.io/install.html
.PHONY: generate-import-client
generate-import-client: ## Generate import client
	@echo "==> Generating import client"
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	swagger generate client -f pkg/tidbcloud/import/import-api.json -A tidbcloud-import -t pkg/tidbcloud/import

.PHONY: generate-connect-info-client
 generate-connect-info-client:
	@echo "==> Generating connect info client"
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	swagger generate client -f pkg/tidbcloud/connect_info/connect-info-api.json -A tidbcloud-connect-info -t pkg/tidbcloud/connect_info

.PHONY: generate-pingchat-client
generate-pingchat-client: ## Generate PingChat client
	@echo "==> Generating PingChat client"
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	swagger generate client -f pkg/tidbcloud/pingchat/pingchat_swagger.json -A tidbcloud-pingchat -t pkg/tidbcloud/pingchat

.PHONY: addcopy
addcopy: ## Add copyright to all files
	@scripts/add-copy.sh

.PHONY: generate-branch-client
 generate-branch-client:
	@echo "==> Generating branch client"
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	swagger generate client -f pkg/tidbcloud/branch/branch.swagger.json -A tidbcloud-branch -t pkg/tidbcloud/branch

.PHONY: generate-serverless-client
 generate-serverless-client:
	@echo "==> Generating serverless client"
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	swagger generate client -f pkg/tidbcloud/serverless/serverless.swagger.json -A tidbcloud-serverless -t pkg/tidbcloud/serverless

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
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
