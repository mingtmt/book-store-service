[![codecov](https://codecov.io/github/mingtmt/book-store/graph/badge.svg?token=9FLQFFXLLH)](https://codecov.io/github/mingtmt/book-store)

# Book Store API

A simple, clean, and extensible RESTful API for managing books, built with Go, Gin, and PostgreSQL.

## Features

- CRUD operations for books (Create, Read, Update, Delete)
- Clean architecture: handler, service, repository layers
- PostgreSQL with SQL migrations and type-safe queries (sqlc)
- Structured logging with zerolog
- Request ID middleware for traceability
- Centralized error handling

## Project Structure

```
cmd/                # Application entrypoint
internal/
  books/
    controller/     # HTTP handlers (Gin)
    application/    # Business logic (services)
    domain/         # Domain models
    infrastructure/ # Persistence (repositories, db)
  middleware/       # Gin middleware (error, request ID)
  initialize/       # DI, DB setup
migrations/         # SQL migrations
pkg/                # Shared packages (logger, errors)
test/               # HTTP request samples
```

## Quick Start

1. **Clone the repo**
2. **Configure your DB**: copy `.env.example` to `.env` and set `DB_URL`
3. **Run migrations**: `make migrate`
4. **Start the server**: `go run cmd/main.go`
5. **Test the API**: see `test/post.http` for sample requests

## Tech Stack

- Go 1.20+
- Gin Web Framework
- PostgreSQL
- goose
- sqlc (type-safe SQL)
- zerolog (logging)

## License

MIT
