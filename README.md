# Bookmark Management

A Go-based bookmark management application with a clean architecture pattern, supporting password generation functionality.

## Project Structure

```
bookmark-management/
├── cmd/                          # Application entry points
│   └── api/                      # Main API server entry point
│       └── main.go               # Application bootstrap and startup
│
├── internal/                     # Private application code (not importable by other projects)
│   ├── api/                      # API layer - HTTP server setup and routing
│   │   ├── api.go                # Gin framework implementation
│   │   └── api_mux.go            # Standard library mux implementation
│   │
│   ├── handler/                  # HTTP request handlers - translates HTTP requests to service calls
│   │   └── pass_handler.go       # Password generation HTTP handlers (Gin & Mux)
│   │
│   ├── service/                  # Business logic layer - core application logic
│   │   ├── pass_service.go       # Password generation service implementation
│   │   └── pass_service_test.go  # Unit tests for password service
│   │
│   ├── model/                    # Data models - domain entities and data structures
│   │                             # (Currently empty - for future bookmark models)
│   │
│   └── repository/               # Data access layer - database/storage abstractions
│                                 # (Currently empty - for future data persistence)
│
├── go.mod                        # Go module dependencies
├── go.sum                        # Go module checksums
└── README.md                     # Project documentation
```

## Folder Descriptions

### `cmd/`
**Purpose**: Contains the main entry points for the application.

- **`cmd/api/`**: The main application entry point that initializes and starts the HTTP server. This follows Go's standard project layout where each executable has its own directory under `cmd/`.

### `internal/`
**Purpose**: Contains private application code that is not meant to be imported by other projects. This enforces encapsulation and prevents external dependencies on internal implementation details.

#### `internal/api/`
**Purpose**: API layer responsible for HTTP server setup, routing, and framework initialization.

- Handles the creation of HTTP servers using different frameworks (Gin and standard library mux)
- Registers routes and connects handlers to endpoints
- Manages server lifecycle (startup and shutdown)

#### `internal/handler/`
**Purpose**: HTTP request handlers that translate HTTP requests into service calls and format responses.

- Acts as the presentation layer, handling HTTP-specific concerns
- Validates and parses incoming requests
- Calls appropriate services to execute business logic
- Formats and returns HTTP responses
- Currently implements password generation handlers for both Gin and standard library mux

#### `internal/service/`
**Purpose**: Business logic layer containing the core application logic.

- Implements domain-specific business rules and operations
- Independent of HTTP concerns (can be used by handlers, CLI tools, etc.)
- Contains the password generation service with cryptographic security
- Includes unit tests to ensure correctness

#### `internal/model/`
**Purpose**: Data models and domain entities.

- Defines the structure of domain objects (bookmarks, users, etc.)
- Currently empty, reserved for future bookmark-related models
- Will contain structs representing the core business entities

#### `internal/repository/`
**Purpose**: Data access layer providing abstractions for data persistence.

- Defines interfaces and implementations for database/storage operations
- Abstracts away database-specific details from the service layer
- Currently empty, reserved for future data persistence implementations
- Will handle CRUD operations for bookmarks and other entities

## Architecture

This project follows a **layered architecture** pattern:

1. **API Layer** (`internal/api/`) - HTTP server and routing
2. **Handler Layer** (`internal/handler/`) - Request/response handling
3. **Service Layer** (`internal/service/`) - Business logic
4. **Repository Layer** (`internal/repository/`) - Data access (to be implemented)
5. **Model Layer** (`internal/model/`) - Domain entities (to be implemented)

## Features

- Password generation with cryptographic security
- Support for multiple HTTP frameworks (Gin and standard library mux)
- Clean architecture with separation of concerns
- Unit tests for business logic

## Getting Started

### Prerequisites

- Go 1.25.5 or higher

### Running the Application

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

### Available Endpoints

- `GET /generate-password` - Generates a secure random password (12 characters, alphanumeric)

## Development

### Running Tests

```bash
go test ./internal/service/...
```

## License

[Add your license here]
