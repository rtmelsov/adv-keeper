APP := tui
PKG := ./cmd/tui
DIST := dist
LDFLAGS := -s -w
CGO := 0

.PHONY: build run clean linux-amd64 linux-arm64 windows-amd64 darwin-arm64 release tidy

build:
	mkdir -p $(DIST)
	go mod tidy
	CGO_ENABLED=$(CGO) go build -ldflags '$(LDFLAGS)' -o $(DIST)/$(APP) $(PKG)

run: build
	$(DIST)/$(APP)

clean:
	rm -rf $(DIST)

# Универсальный шаблон: build-<goos>_<goarch>
build-%:
	mkdir -p $(DIST)
	@GOOS=$(word 1,$(subst _, ,$*)) GOARCH=$(word 2,$(subst _, ,$*)) \
	CGO_ENABLED=$(CGO) go build -ldflags '$(LDFLAGS)' \
	-o $(DIST)/$(APP)_$*$(if $(findstring windows,$(word 1,$(subst _, ,$*))),.exe,) $(PKG)

linux-amd64:  build-linux_amd64
linux-arm64:  build-linux_arm64
windows-amd64: build-windows_amd64
darwin-arm64: build-darwin_arm64

# Сборка всего и упаковка
release: linux-amd64 linux-arm64 windows-amd64 darwin-arm64
	cd $(DIST) && \
	tar -czf $(APP)_linux_amd64.tgz  $(APP)_linux_amd64 && \
	tar -czf $(APP)_linux_arm64.tgz  $(APP)_linux_arm64 && \
	tar -czf $(APP)_darwin_arm64.tgz $(APP)_darwin_arm64 && \
	zip -9    $(APP)_windows_amd64.zip $(APP)_windows_amd64.exe

tidy:
	go mod tidy

