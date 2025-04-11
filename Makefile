GOPATH := $(shell go env GOPATH)
export PATH := $(PATH):$(GOPATH)/bin

LINT_VERSION = 1.64.5
GOLANGCI_LINT_PATH = $(GOPATH)/bin/golangci-lint


.PHONY: generate-swagger
generate-swagger:
	@echo "Generating Swagger docs..."
	@swag init -g ./cmd/main/main.go

.PHONY: start-linters
start-linters: ensure-golangci-lint
	golangci-lint run

.PHONY: ensure-golangci-lint
ensure-golangci-lint:
	@if  ! test -f $(GOLANGCI_LINT_PATH); then \
		echo "golangci-lint doesn't installed. installing version $(LINT_VERSION)..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(LINT_VERSION); \
	elif ! golangci-lint --version | grep -q $(LINT_VERSION); then \
		echo "golangci-lint installed, but not version $(LINT_VERSION). installing version $(LINT_VERSION)..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(LINT_VERSION); \
	else \
		echo "golangci-lint version $(LINT_VERSION) installed"; \
	fi

# migrate create -ext sql -dir migrations -seq create_teachers_table