.PHONY: build lint run clean dep
.DEFAULT_GOAL := run

BINARY_NAME = chimes
GOLANGCI_LINT_VERSION = v1.42.1
BUILD_DIR ?= ./bin
DEFAULT_OS ?= linux

# Build the project
build:
		@echo "Building Project 🔨"
		GOARCH=amd64 GOOS=darwin go build -o ${BUILD_DIR}/${BINARY_NAME}-darwin ./cmd/chimes 
		GOARCH=amd64 GOOS=linux go build -o ${BUILD_DIR}/${BINARY_NAME}-linux ./cmd/chimes 
		GOARCH=amd64 GOOS=windows go build -o ${BUILD_DIR}/${BINARY_NAME}-windows ./cmd/chimes 

# Run the project
run:	build
				@echo "Running Project 🚀"
				./${BUILD_DIR}/${BINARY_NAME}-${DEFAULT_OS}

# Link the project
lint: 
		@echo "Linting Project 🔍️"
		golangci-lint run

# Clean up
clean: 
		@echo "Cleaning up Project 🧹"
		go clean
		rm ./${BUILD_DIR}/${BINARY_NAME}-darwin
		rm ./${BUILD_DIR}/${BINARY_NAME}-linux
		rm ./${BUILD_DIR}/${BINARY_NAME}-windows

# Download all dependencies
dep:
	@echo "Downloading Project Dependencies 📦️"
	go mod download

# Installing golangci-lint
install-linter:
		@echo "Installling golangci-lint ⏬"
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)