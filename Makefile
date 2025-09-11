# Makefile (минимальный)
APP := tui
PKG := ./cmd/tui
DIST := dist

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

