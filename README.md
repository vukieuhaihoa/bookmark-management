# Bookmark Management

A Go-based bookmark management application with a clean architecture pattern, supporting password generation and health check functionality.

## Project Structure

```
bookmark-management/
├── cmd/                                    # Application entry points
│   └── api/                                # Main API server entry point
│       └── main.go                         # Application bootstrap and startup
│
├── internal/                               # Private application code (not importable by other projects)
│   ├── api/                                # API layer - HTTP server setup and routing
│   │   ├── api.go                          # Gin framework implementation with Engine interface
│   │   └── config.go                       # Configuration management with environment variables
│   │
│   ├── handler/                            # HTTP request handlers - translates HTTP requests to service calls
│   │   ├── pass_handler.go                 # Password generation HTTP handlers
│   │   ├── pass_handler_test.go            # Unit tests for password handler
│   │   ├── health_check_handler.go         # Health check HTTP handlers
│   │   └── health_check_handler_test.go    # Unit tests for health check handler
│   │
│   ├── service/                            # Business logic layer - core application logic
│   │   ├── pass_service.go                 # Password generation service implementation
│   │   ├── pass_service_test.go            # Unit tests for password service
│   │   ├── health_check_service.go         # Health check service implementation
│   │   ├── health_check_service_test.go    # Unit tests for health check service
│   │   └── mocks/                          # Mock implementations for testing
│   │       ├── pass_service.go             # Mock password service
│   │       └── health_check_service.go     # Mock health check service
│   │
│   ├── test/                               # Integration and endpoint tests
│   │   └── endpoint/                       # Endpoint-level integration tests
│   │       └── password_endpoint_test.go   # Password endpoint integration tests
│   │
│   ├── model/                              # Data models - domain entities and data structures
│   │                                       # (Reserved for future bookmark models)
│   │
│   └── repository/                         # Data access layer - database/storage abstractions
│                                           # (Reserved for future data persistence)
│
├── .gitignore                              # Git ignore patterns
├── go.mod                                  # Go module dependencies
├── go.sum                                  # Go module checksums
└── README.md                               # Project documentation
```

## Folder Descriptions

### `cmd/`
**Purpose**: Contains the main entry points for the application.

- **`cmd/api/`**: The main application entry point that initializes and starts the HTTP server. This follows Go's standard project layout where each executable has its own directory under `cmd/`.

### `internal/`
**Purpose**: Contains private application code that is not meant to be imported by other projects. This enforces encapsulation and prevents external dependencies on internal implementation details.

#### `internal/api/`
**Purpose**: API layer responsible for HTTP server setup, routing, and framework initialization.

- Implements the `Engine` interface for HTTP server abstraction
- Handles the creation of HTTP servers using the Gin framework
- Registers routes and connects handlers to endpoints
- Manages server lifecycle (startup and shutdown)
- Provides configuration management through environment variables
- Supports the `http.Handler` interface via `ServeHTTP` method

#### `internal/handler/`
**Purpose**: HTTP request handlers that translate HTTP requests into service calls and format responses.

- Acts as the presentation layer, handling HTTP-specific concerns
- Validates and parses incoming requests
- Calls appropriate services to execute business logic
- Formats and returns HTTP responses
- Includes comprehensive unit tests with mocked dependencies
- Currently implements:
  - Password generation handlers
  - Health check handlers

#### `internal/service/`
**Purpose**: Business logic layer containing the core application logic.

- Implements domain-specific business rules and operations
- Independent of HTTP concerns (can be used by handlers, CLI tools, etc.)
- Contains the password generation service with cryptographic security
- Contains the health check service for monitoring application status
- Includes comprehensive unit tests to ensure correctness
- Provides mock implementations in the `mocks/` subdirectory for testing

#### `internal/model/`
**Purpose**: Data models and domain entities.

- Defines the structure of domain objects (bookmarks, users, etc.)
- Currently empty, reserved for future bookmark-related models
- Will contain structs representing the core business entities

#### `internal/repository/`
**Purpose**: Data access layer providing abstractions for data persistence.

- Defines interfaces and implementations for database/storage operations
- Abstracts away database-specific details from the service layer
- Reserved for future data persistence implementations
- Will handle CRUD operations for bookmarks and other entities

#### `internal/test/`
**Purpose**: Integration and endpoint tests for the application.

- Contains endpoint-level integration tests
- Tests the full request/response cycle
- Validates API behavior and contract
- Uses real HTTP requests to test handlers

## Architecture

This project follows a **layered architecture** pattern:

1. **API Layer** (`internal/api/`) - HTTP server and routing
2. **Handler Layer** (`internal/handler/`) - Request/response handling
3. **Service Layer** (`internal/service/`) - Business logic
4. **Repository Layer** (`internal/repository/`) - Data access (to be implemented)
5. **Model Layer** (`internal/model/`) - Domain entities (to be implemented)

## Features

- **Password Generation**: Secure random password generation with cryptographic security
- **Health Checks**: Application health monitoring with service name and instance ID tracking
- **Configuration Management**: Environment-based configuration with sensible defaults
- **Clean Architecture**: Separation of concerns with layered architecture pattern
- **Comprehensive Testing**: Unit tests, integration tests, and mock implementations
- **HTTP Server Abstraction**: Engine interface for flexible HTTP server implementations
- **Gin Framework**: High-performance HTTP routing with middleware support

## Getting Started

### Prerequisites

- Go 1.25.5 or higher

### Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Server port | `:8080` |
| `SERVICE_NAME` | Service identifier | `bookmark_service` |
| `INSTANCE_ID` | Unique instance identifier | Auto-generated UUID |

### Running the Application

```bash
# Run with default configuration
go run cmd/api/main.go

# Run with custom port
APP_PORT=:3000 go run cmd/api/main.go

# Run with custom service name and instance ID
SERVICE_NAME=my_service INSTANCE_ID=instance-1 go run cmd/api/main.go
```

The server will start on `http://localhost:8080` (or your configured port)

### Available Endpoints

- `GET /health-check` - Returns application health status, service name, and instance ID
- `GET /generate-password` - Generates a secure random password (12 characters, alphanumeric)

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests in a specific package
go test ./internal/service/...
go test ./internal/handler/...
go test ./internal/test/endpoint/...

# Run tests with verbose output
go test -v ./...
```

### Project Dependencies

Main dependencies include:
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/google/uuid` - UUID generation for instance IDs
- `github.com/kelseyhightower/envconfig` - Environment variable configuration
- `github.com/stretchr/testify` - Testing toolkit with assertions and mocks

## License

[Add your license here]
