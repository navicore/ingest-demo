.PHONY: build test lint clean run-generate run-process all

# Build variables
BINARY_NAME=ingest
GO=go

all: clean build test

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/ingest

test:
	$(GO) test -v ./...

# Run integration tests
test-integration:
	ENABLE_INTEGRATION_TESTS=true $(GO) test -v ./...

# Run the linter
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf test_output

# Run generator mode
run-generate:
	./$(BINARY_NAME) --generate --count 10

# Run processor mode
run-process:
	./$(BINARY_NAME) --generate --count 10 | ./$(BINARY_NAME)

# Initialize a new git repository
init-git:
	git init
	git add .
	git commit -m "Initial commit"

# Install development dependencies
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest