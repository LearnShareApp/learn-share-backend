.PHONY: generate-swagger

generate-swagger:
	@echo "Getting GOPATH..."
	@GOPATH=$$(go env GOPATH); \
	echo "GOPATH is: $$GOPATH"; \
	echo "Adding $$GOPATH/bin to PATH"; \
	export PATH=$$PATH:$$GOPATH/bin; \
	swag init -g ./cmd/main/main.go
