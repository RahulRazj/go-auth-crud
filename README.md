# go-auth-crud

A simple Go authentication CRUD service built with a clean project structure.

## Features

- Config management with environment variables
- Structured logging using zerolog
- Modular internal packages for config, logger, and API
- Support for CORS origin configuration

## Getting Started

### Requirements

- Go 1.21+
- PostgreSQL
- `task` CLI (optional, if using Taskfile)

### Setup

1. Copy `.env.example` to `.env`:

    ```bash
    cp .env.example .env
    ```

2. Update `.env` with your database connection values.

3. Install dependencies:

    ```bash
    go mod tidy
    ```

### Run

```bash
go run ./cmd/api
```

Or with Taskfile:

```bash
task run
```

### Development

```bash
task run:dev
```

### Helpful tasks

- `task fmt` - format Go code
- `task lint` - lint code
- `task test` - run tests
- `task build` - build the application
- `task clean` - clean generated artifacts

## Environment

Use `.env.example` as a template for environment variables.

## Notes

This repository is intended as a starting point for building a Go-based auth CRUD service.
