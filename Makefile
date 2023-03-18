SHELL := /bin/bash

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==============================================================================
# Development support

## run: run the server
.PHONY: run
run:
	@go run main.go

## build: build the server
.PHONY: build
build:
	@echo "Building the server..."
	@go build -o secret-app main.go
	@echo "Done!"
