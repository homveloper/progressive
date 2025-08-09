# Go + Templ + Tailwind CSS Build Pipeline

.PHONY: help dev build clean tailwind-build tailwind-watch templ-generate templ-watch run air-dev

# Default target
help:
	@echo "Available targets:"
	@echo "  dev              - Start development server with live reload (Air + Tailwind watch)"
	@echo "  build            - Build the complete application"
	@echo "  run              - Run the application"
	@echo "  clean            - Clean build artifacts"
	@echo "  tailwind-build   - Build Tailwind CSS (production)"
	@echo "  tailwind-watch   - Watch and build Tailwind CSS (development)"
	@echo "  templ-generate   - Generate templ files"
	@echo "  templ-watch      - Watch and generate templ files"
	@echo "  air-dev          - Start Air development server with live reload"

# Development - runs everything with live reload
dev:
	@echo "Starting development servers..."
	@make -j2 tailwind-watch air-dev

# Air development server with live reload
air-dev:
	@echo "Starting Air development server..."
	@air

# Production build
build: clean tailwind-build templ-generate
	@echo "Building Go application..."
	@go build -o ./tmp/progressive ./cmd/web
	@echo "Build complete! Binary: ./tmp/progressive"

# Run the Go application
run:
	@echo "Starting Go application..."
	@go run ./cmd/web

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ./tmp
	@rm -f ./static/css/output.css
	@rm -f ./internal/pages/*_templ.go ./internal/components/*_templ.go

# Tailwind CSS - Production build (minified)
tailwind-build:
	@echo "Building Tailwind CSS (production)..."
	@npx @tailwindcss/cli -c ./tailwind.config.js -i ./static/css/input.css -o ./static/css/output.css

# Tailwind CSS - Development watch (with live reload)
tailwind-watch:
	@echo "Watching Tailwind CSS files..."
	@npx @tailwindcss/cli -c ./tailwind.config.js -i ./static/css/input.css -o ./static/css/output.css --watch

# Templ - Generate template files
templ-generate:
	@echo "Generating templ files..."
	@templ generate

# Templ - Watch and generate template files
templ-watch:
	@echo "Watching templ files..."
	@templ generate --watch

# Create tmp directory if it doesn't exist
./tmp:
	@mkdir -p ./tmp
