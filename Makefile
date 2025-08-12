.PHONY: build run-server run-client test clean benchmark help all

build:
	go build -o gcache ./cmd/server

#Run the server with default settings
run-server: build
	./gcache -mode=server -addr=localhost:8080 -capacity=1000


# Run interactive client
run-client: build
	./gcache -mode=client -addr=localhost:8080 -interactive

# Run client demo
demo: build
	./gcache -mode=client -addr=localhost:8080

# Clean build artifacts
clean:
	rm -f gcache
# Install dependencies (if any)
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Show help
help:
	@echo "Available commands:"
	@echo "  build       - Build the gcache binary"
	@echo "  run-server  - Run the cache server"
	@echo "  run-client  - Run interactive client"
	@echo "  demo        - Run client demo"
	@echo "  test        - Run unit tests"
	@echo "  test-race   - Run tests with race detection"
	@echo "  test-server - Run server integration tests"
	@echo "  benchmark   - Run performance benchmarks"
	@echo "  clean       - Remove build artifacts"
	@echo "  fmt         - Format Go code"
	@echo "  help        - Show this help"

# Default target
all: build
