.PHONY: build build-embed run generate test test-unit test-integration lint clean dev migrate-up migrate-new

APP_NAME := server
BUILD_DIR := bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
MIGRATIONS_DIR := migrations

build:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$(VERSION)" -trimpath -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

build-embed:
	CGO_ENABLED=0 go build -tags embed -ldflags="-s -w -X main.Version=$(VERSION)" -trimpath -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

dev:
	go run ./cmd/server

generate:
	go generate ./ent
	go generate ./cmd/server

# Generate a new empty migration file
migrate-new:
	@read -p "Enter migration name: " name; \
	touch $(MIGRATIONS_DIR)/$(shell date +%Y%m%d%H%M%S)_$$name.sql; \
	echo "-- $(MIGRATIONS_DIR)/$(shell date +%Y%m%d%H%M%S)_$$name.sql created"

# Run migrations (handled by Ent auto-migration at startup)
# To generate a migration from Ent schema changes:
migrate-diff:
	go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema

test:
	go test ./...

test-unit:
	go test -tags=unit ./...

test-integration:
	go test -tags=integration ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)
