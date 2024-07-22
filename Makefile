BINARY_NAME=squeeze_box_server

all: build

build:
	@echo "Building..."
	go build -o ./bin/server/$(BINARY_NAME) ./cmd/server

run: build
	@echo "Starting server..."
	./bin/server/$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	rm -rf ./bin/server/*

dev:
	@echo "[Dev Mode] Starting server..."
	go run ./cmd/server/main.go

help:
	@echo "Makefile commands:"
	@echo "all   - Build the application"
	@echo "build - Build the binary"
	@echo "run   - Build and run the application"
	@echo "clean - Clean cache and Remove binaries"
	@echo "help  - Show this help menu"
