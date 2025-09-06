SERVER_PKG ?= ./cmd/server

GRPC_ADDR ?= 127.0.0.1:8080
DB_DSN ?= mem://

.PHONY: run build tidy help

help:
	@echo "make run  - run server via 'go run'"
	@echo "make build  -build server binary into ./bin/server"
	@echo "env overrides: GRPC_ADDR (default: $(GRPC_ADDR)), DB_DSN (default: $(DB_DSN))"
	@echo "example: make run GRPC_ADDR=:9090 DB_DSN='postgres://...@localhost:5432/vault?sslmode=disable"

run:
	@echo "===> go run $(SERVER_PKG) (DB_DSN=$(GRPC_ADDR), DB_DSN=$(DB_DSN))"
	@GRPC_ADDR=$(GRPC_ADDR) DB_DSN=$(DB_DSN) go run $(SERVER_PKG)

build:
	@mkdir -p bin
	@echo "===> go build ./cmd/server -> ./bin/server"
	@go build -o bin/server $(SERVER_PKG)

tidy:
	@go mod tidy
