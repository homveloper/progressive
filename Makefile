# Progressive Spreadsheet Makefile

.PHONY: dev build-wasm build clean deps

# Development server with WASM build
dev: build-wasm
	@echo "Starting development server..."
	go run *.go

# Build WebAssembly
build-wasm:
	@echo "Building WebAssembly..."
	@mkdir -p web
	GOOS=js GOARCH=wasm go build -o web/app.wasm *.go
	@echo "WebAssembly built: web/app.wasm"

# Build server binary
build:
	@echo "Building server binary..."
	go build -o progressive *.go
	@echo "Server binary built: progressive"

# Build both
build-all: build-wasm build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f progressive
	rm -f web/app.wasm

# Copy wasm_exec.js from Go installation
wasm-exec:
	@echo "Copying wasm_exec.js..."
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" web/

# Setup project (run once)
setup: deps wasm-exec build-wasm
	@echo "Project setup complete!"
	@echo "Run 'make dev' to start development server"

# Help
help:
	@echo "Available commands:"
	@echo "  dev        - Start development server"
	@echo "  build-wasm - Build WebAssembly binary"
	@echo "  build      - Build server binary"
	@echo "  build-all  - Build both WebAssembly and server"
	@echo "  deps       - Install dependencies"
	@echo "  clean      - Clean build artifacts"
	@echo "  setup      - Setup project (run once)"
	@echo "  help       - Show this help"