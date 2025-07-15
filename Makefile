ifneq (,$(wildcard .env))
  include .env
  export
endif

# Path to PostgreSQL connection string (edit or use from .env)
DB_URL := $(DB_URL)

# Default target
.PHONY: help
help:
	@echo "📘 Available commands:"
	@echo "  make install-tools    - Install goose and sqlc"
	@echo "  make run-services     - Start all services using docker-compose"
	@echo "  make sqlc             - Generate code from SQL files"
	@echo "  make migrate          - Run goose up migration"
	@echo "  make migrate-down     - Rollback one migration"
	@echo "  make reset            - Drop all tables and rerun migrations (⚠️ destructive)"
	@echo "  make setup            - Full project setup"

# Install goose and sqlc
.PHONY: install-tools
install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Run docker-compose to start all services
.PHONY: run-services
run-services:
	docker-compose up -d

# Generate Go code from SQL
.PHONY: sqlc
sqlc-books:
	sqlc generate -f internal/books/infrastructure/persistence/sqlc.yml
	
sqlc-auths:
	sqlc generate -f internal/auth/infrastructure/persistence/sqlc.yml

sqlc-all: sqlc-books sqlc-auths

# Run goose migrations
.PHONY: migrate
migrate:
	goose -dir migrations postgres "$(DB_URL)" up

# Rollback last migration
.PHONY: migrate-down
migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

# Drop all tables and re-run migrations (⚠️ destructive)
.PHONY: reset
reset:
	goose -dir migrations postgres "$(DB_URL)" reset

# Generate Swagger documentation
.PHONY: swag
swag:
	swag init -g cmd/main.go

# Run tests
.PHONY: test
test:
	go test -v ./...

# Full setup
.PHONY: setup
setup: install-tools sqlc migrate
