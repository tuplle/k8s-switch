BINARY_NAME=k8s-switch
BUILD_DIR=bin

.PHONY: all install-deps update-deps build clean test snapshot

all: clean install-deps test build

## install-deps: Download dependencies pinned in go.mod/go.sum (does not change versions)
install-deps:
	@echo "Installing dependencies"
	@go mod download

## update-deps: Upgrade dependencies to their latest versions (run deliberately, never as part of build/release)
update-deps:
	@echo "Updating dependencies"
	@go get -u ./...
	@go mod tidy

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

## snapshot: Build a local snapshot release with GoReleaser (no publish, no tag needed)
snapshot:
	@echo "Building snapshot release"
	@goreleaser release --snapshot --clean --skip=publish
