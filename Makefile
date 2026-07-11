ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: fmt lint vuln test build clean up

BINARY_NAME := rbac-system
BIN_DIR := ./bin

DB_USER ?= $(DATABASE_USER)
DB_PASS ?= $(DATABASE_PASS)
DB_HOST ?= $(DATABASE_HOST)
DB_NAME ?= $(DATABASE_NAME)
DB_PORT ?= $(DATABASE_PORT)

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

migrate-up:
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
        -path migrations up

migrate-down:
	migrate -database "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		-path migrations down