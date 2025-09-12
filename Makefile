# Makefile (минимальный)
APP := tui
PKG := ./cmd/tui
DIST := dist
SERVER_PKG ?= ./cmd/server
GRPC_ADDR ?= 127.0.0.1:8080
DB_DSN ?= mem://



.PHONY: build run clean linux mac windows

build:
	mkdir -p $(DIST)
	go mod tidy
	CGO_ENABLED=0 go build -o $(DIST)/$(APP) $(PKG)

run: build
	$(DIST)/$(APP)

clean:
	rm -rf $(DIST)

# Кросс-сборки, если нужно
linux:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(APP)_linux_amd64 $(PKG)

mac:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(DIST)/$(APP)_darwin_arm64 $(PKG)

windows:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(DIST)/$(APP)_windows_amd64.exe $(PKG)

help:
	@echo "make run  - run server via 'go run'"
	@echo "make build  -build server binary into ./bin/server"
	@echo "env overrides: GRPC_ADDR (default: $(GRPC_ADDR)), DB_DSN (default: $(DB_DSN))"
	@echo "example: make run GRPC_ADDR=:9090 DB_DSN='postgres://...@localhost:5432/vault?sslmode=disable"
buildServer:
	@mkdir -p bin
	@echo "===> go build ./cmd/server -> ./bin/server"
	@go build -o bin/server $(SERVER_PKG)

tidy:
	@go mod tidy
