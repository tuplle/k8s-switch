BINARY_NAME=k8s-switch
BUILD_DIR=bin

.PHONY: all install-deps build clean test

all: clean install-deps test build

## install-deps: Get dependencies
install-deps:
	@echo "Installing dependencies"
	@go get -u .

## test: Run the test suite
test:
	@echo "Testing"
	@go test ./...

## build: Build the binary
build:
	@echo "Building"
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) github.com/tuplle/k8s-switch

## clean: Remove build artifacts
clean:
	@echo "Cleaning"
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)

## install: Install the binary as go command
install:
	@echo "Installing"
	@go install github.com/tuplle/k8s-switch

## lint: Lint and format project
lint:
	@echo "Linting"
	@go vet ./...
	@echo "Formatting"
	@go fmt ./...
