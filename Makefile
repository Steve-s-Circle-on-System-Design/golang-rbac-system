.PHONY: fmt lint test build clean

BINARY_NAME := rbac-system
BIN_DIR := ./bin

fmt:
	gofumpt -l -w .
	goimports -w .

lint:
	golangci-lint run ./...

test:
	go test -race -cover ./...

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/...

clean:
	rm -rf $(BIN_DIR)
	go clean -cache -testcache