# Makefile for Go project

# Variables
APP_NAME := /tmp/web
SRC_DIR := ./cmd/web
TEST_DIR := ./...
TLS_SRC := ./tls
TLS_DEST := /tmp/tls

# Build the application
build:
	go build -o $(APP_NAME) $(SRC_DIR)

# Run the application
run: build copy-tls
	$(APP_NAME)

# Run tests
test:
	go test $(TEST_DIR)

# Copy TLS directory
copy-tls:
	cp -r $(TLS_SRC) $(TLS_DEST)

# Clean up build artifacts
clean:
	rm -f $(APP_NAME)

.PHONY: build run test clean copy-tls