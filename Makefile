.PHONY: generate-swagger
generate-swagger:
	export PATH=$PATH:$(go env GOPATH)/bin
	swag init -g ./cmd/main/main.go
