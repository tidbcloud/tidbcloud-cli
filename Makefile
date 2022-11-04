GOLANGCI_VERSION=v1.49.0

.PHONY: deps
deps:  ## Download go module dependencies
	@echo "==> Installing go.mod dependencies..."
	go mod download
	go mod tidy

.PHONY: devtools
devtools:  ## Install dev tools
	@echo "==> Installing dev tools..."
	go install github.com/google/addlicense@latest
	go install github.com/golang/mock/mockgen@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/google/go-licenses@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

.PHONY: setup
setup: deps devtools ## Set up dev env

.PHONY: addcopy
addcopy: ## Add copyright to all files
	@scripts/add-copy.sh

.PHONY: fmt
fmt: ## Format changed go
	@scripts/fmt.sh

.PHONY: fix-lint
fix-lint: ## Fix linting errors
	golangci-lint run --fix

.PHONY: gen-mocks
gen-mocks: ## Generate mocks
	@echo "==> Generating mocks"
	go generate ./internal...

.PHONY: test
test: ## Run unit-tests
	@go test ./...

.PHONY: build
build: ## Generate a binary in ./bin
	@go build -o ./bin/ticloud ./cmd/

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
