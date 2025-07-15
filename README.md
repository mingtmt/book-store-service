[![codecov](https://codecov.io/github/mingtmt/book-store/graph/badge.svg?token=9FLQFFXLLH)](https://codecov.io/github/mingtmt/book-store)
[![CodeQL](https://github.com/mingtmt/book-store/actions/workflows/codeql.yml/badge.svg)](https://github.com/mingtmt/book-store/actions/workflows/codeql.yml)

# Book Store API

A simple, clean, and extensible RESTful API for managing books, built with Go, Gin, and PostgreSQL.

## Features

- CRUD operations for books (Create, Read, Update, Delete)
- Clean architecture: handler, service, repository layers
- PostgreSQL with SQL migrations and type-safe queries (sqlc)
- Structured logging with zerolog
- Request ID middleware for traceability
- Centralized error handling
- JWT authentication for protected endpoints
- Swagger (OpenAPI) documentation

## Project Structure

```
cmd/                # Application entrypoint
internal/
  books/
    controller/     # HTTP handlers (Gin)
    application/    # Business logic (services)
    domain/         # Domain models
    infrastructure/ # Persistence (repositories, db)
  auth/             # Authentication module
  middleware/       # Gin middleware (error, request ID, JWT)
  initialize/       # DI, DB, router setup
migrations/         # SQL migrations
pkg/                # Shared packages (logger, errors, token, response)
test/               # HTTP request samples
```

## Development Setup

### Prerequisites

- Go 1.20+
- PostgreSQL
- [Goose](https://github.com/pressly/goose) (for migrations)
- [sqlc](https://sqlc.dev/) (for type-safe SQL)
- [Swag](https://github.com/swaggo/swag) (for API docs)

### 1. Clone the repository

```sh
git clone https://github.com/mingtmt/book-store-service.git
cd book-store-service
```

### 2. Configure environment variables

Copy `.env.example` to `.env` and set your `DB_URL` and `PORT` as needed.

### 3. Install dependencies

```sh
go mod tidy
```

### 4. Run database migrations

```sh
make migrate
```

### 5. Generate SQL code and Swagger docs

```sh
make sqlc
swag init -g cmd/main.go
```

### 6. Start the development server

```sh
go run cmd/main.go
```

The server will run at `http://localhost:8080` by default.

## How to Use

### API Documentation

- Visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) for interactive Swagger docs.

### Run Tests

```sh
go test ./...
```

## Tech Stack

- Go 1.20+
- Gin Web Framework
- PostgreSQL
- goose (migrations)
- sqlc (type-safe SQL)
- zerolog (logging)
- swaggo/swag (Swagger docs)

## License

MIT
