# GoPlace - Guest Service

The GoPlace Guest Service manages guests and their reservations within the GoPlace application.

## Table of Contents

- [Getting Started](#getting-started)
- [Directory Structure](#directory-structure)
- [Endpoints](#endpoints)
- [Configuration](#configuration)
- [Running Tests](#running-tests)

## Getting Started

### Prerequisites

- Golang 1.22 or higher
- Docker

### Installation

1. Clone the repository:

```bash
git clone git@github.com:goplaceapp/goplace-guest.git
cd goplace-guest
```

2. Install dependencies:

```bash
go mod tidy
```

### Running the Service

To run the service:

```bash
go run cmd/main.go
```

Or using Docker:

```bash
docker build -t goplace-guest .
docker run -p 8080:8080 goplace-guest
```

## Directory Structure

- **api/v1**: Contains the API endpoints for version 1 of the microservice.
  - **guest_service.proto**: Protobuf definitions for RPCs and structs.
  - Generate protobuf files using `make gen` and clear them with `make clean`.
- **cmd**: Main entry point of the application.
  - Run the application via `docker` or `go run`.
- **config**: Configuration files and settings.
- **database**: Handles database connections and queries.
  - **database.go**: Gorm configuration.
  - **migrator.go**: Defines and runs shared or tenant migrations.
  - **postgres.go**: Manages tenant connection pool and runs migrations.
  - **seed.go**: Seeds data in all environments.
  - **triggers.go**: Deprecated trigger functions for database actions.
- **internal**: Internal packages not exposed externally.
  - **clients/**: Client connections to other microservices for method access.
  - **services/**: Services, including protobuf, with three layers:
    - **adapters/**:
      - **converters/**: Methods to convert data types.
      - **grpc/**: gRPC handlers for RPCs, called first by the gateway.
    - **application/**: Handles requests from gRPC to the infrastructure layer, applying security or validation checks.
    - **domain/**: Contains domain schemas.
    - **infrastructure/repository/**: Repositories interacting with the database and processing data before response.
- **migrations**: Database migration files.
  - **shared/**: Migrations for the shared database.
  - **tenant/**: Migrations for tenant databases.
- **pkg**: Common packages used across the application.
- **scripts**: Utility scripts.
- **server**: Server-related files and settings.
- **tests**: Test files and configurations.
- **utils**: Utility functions and helpers.

## Running Tests

To execute the tests, run:

```bash
go test ./...
```
