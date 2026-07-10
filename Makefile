.PHONY: fmt lint vuln test build clean up

BINARY_NAME := rbac-system
BIN_DIR := ./bin

fmt:
	gofumpt -l -w .
	goimports -w .

lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

test:
	go test -race -cover ./...

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...

clean:
	rm -rf $(BIN_DIR)
	go clean -cache -testcache

up: fmt lint vuln test build