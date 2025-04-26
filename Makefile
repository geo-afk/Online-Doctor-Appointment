# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@@go build cmd/api/main.go
	
	
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@docker compose up --build

# Shutdown DB container
docker-down:
	@docker compose down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"


swag:
	swag init -g cmd/api/main.go


# APP_NAME = apiserver
# BUILD_DIR = $(PWD)/build
# MIGRATIONS_FOLDER = $(PWD)/platform/migrations
# DATABASE_URL = postgres://postgres:password@localhost/postgres?sslmode=disable


# migrate.up:
# 	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

# migrate.down:
# 	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down


.PHONY: all build run test clean watch docker-run docker-down itest
