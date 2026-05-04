.PHONY: help build build-all build-linux build-macos build-windows clean test vet fmt

# Default target
help:
	@echo "Available targets:"
	@echo "  make build          - Build for current platform"
	@echo "  make build-all      - Build for all platforms (Linux, macOS, Windows)"
	@echo "  make build-linux    - Build for Linux x86_64"
	@echo "  make build-macos    - Build for macOS (Intel + ARM)"
	@echo "  make build-windows  - Build for Windows"
	@echo "  make test           - Run tests"
	@echo "  make vet            - Run go vet"
	@echo "  make fmt            - Format code"
	@echo "  make clean          - Remove build artifacts"

# Detect OS and set defaults
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# Build for current platform
build:
	go build -o bin/amnezia-vless-builder ./cmd/amnezia-vless-builder
	@echo "✅ Build complete: bin/amnezia-vless-builder"

# Build all binaries
build-all: build-linux build-macos build-windows
	@echo "✅ All builds complete!"

# Linux build
build-linux:
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
		-o bin/amnezia-vless-builder-linux-amd64 \
		./cmd/amnezia-vless-builder
	@echo "✅ Linux build complete"

# macOS builds
build-macos:
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build \
		-o bin/amnezia-vless-builder-macos-amd64 \
		./cmd/amnezia-vless-builder
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build \
		-o bin/amnezia-vless-builder-macos-arm64 \
		./cmd/amnezia-vless-builder
	@echo "✅ macOS builds complete"

# Windows build
build-windows:
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
		-o bin/amnezia-vless-builder-windows-amd64.exe \
		./cmd/amnezia-vless-builder
	@echo "✅ Windows build complete"

# Run tests
test:
	go test -v ./...

# Run go vet
vet:
	go vet ./...

# Format code
fmt:
	go fmt ./...

# Clean build artifacts
clean:
	@rm -rf bin/
	@go clean
	@echo "✅ Cleaned build artifacts"
