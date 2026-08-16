# Simple Makefile for a Go project

# Build the application
all: build test

# The templ CLI must match the templ runtime in go.mod -- a newer CLI emits
# calls the pinned library does not export. Always derive it, never use @latest.
TEMPL_VERSION := $(shell go list -m -f '{{.Version}}' github.com/a-h/templ)

templ-install:
	@if [ "$$(templ version 2>/dev/null | tr -d 'v ')" != "$$(echo $(TEMPL_VERSION) | tr -d 'v')" ]; then \
		echo "Installing templ $(TEMPL_VERSION) to match go.mod..."; \
		go install github.com/a-h/templ/cmd/templ@$(TEMPL_VERSION); \
		if [ ! -x "$$(command -v templ)" ]; then \
			echo "templ installation failed. Exiting..."; \
			exit 1; \
		fi; \
	fi
# Pinned: an unpinned "latest" download means the generated output.css can change
# without any source change, which breaks the CSS drift check in CI.
TAILWIND_VERSION := v4.1.3

tailwind:
	@if [ ! -f tailwindcss ]; then curl -sL https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/tailwindcss-linux-x64 -o tailwindcss; fi
	@chmod +x tailwindcss

build: tailwind templ-install
	@echo "Building..."
	@templ generate
	@./tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/output.css
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch tailwind templ-install
