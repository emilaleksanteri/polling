# Simple Makefile for a Go project
#

.PHONY: help
help:
	@echo 'Usage: '
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


## all: Build the application
all: build

build:
	@echo "Building..."
	
	@go build -o main cmd/api/main.go

## run: Run the application
run:
	@go run cmd/api/main.go

## test: Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

## clean: Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

## watch: Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
