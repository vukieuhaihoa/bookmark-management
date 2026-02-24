# Bookmark Management

A Go-based REST API for managing bookmarks with URL shortening, user authentication, Redis caching, and rate limiting.

## Project Structure

```
bookmark-management/
├── cmd/
│   ├── api/main.go                         # Main API server entry point
│   ├── migrate/main.go                     # Standalone DB migration runner
│   ├── script/                             # Utility scripts (e.g., backfill)
│   └── test/main.go                        # Test helper entry point
│
├── docs/                                   # Auto-generated Swagger documentation
│
├── internal/
│   ├── api/                                # Gin engine setup, routing, middleware wiring
│   │   ├── api.go                          # Route & handler registration
│   │   ├── config.go                       # App config (env vars)
│   │   └── middleware/
│   │       ├── jwtauth.go                  # JWT authentication middleware
│   │       └── ratelimit.go                # IP-based rate limiting middleware
│   │
│   ├── app/
│   │   ├── handler/                        # HTTP request/response layer (one dir per feature)
│   │   │   ├── bookmark/
│   │   │   ├── healthcheck/
│   │   │   ├── link/
│   │   │   ├── random_code_gen/
│   │   │   ├── user/
│   │   │   └── utils/                      # Shared handler utilities (JWT claims, request parsing)
│   │   │
│   │   ├── service/                        # Business logic layer (one dir per feature)
│   │   │   ├── bookmark/                   # Includes cache decorator (NewBookmarkServiceWithCache)
│   │   │   ├── healthcheck/
│   │   │   ├── link/
│   │   │   ├── random_code_gen/
│   │   │   └── user/
│   │   │
│   │   ├── repository/                     # Data access layer (one dir per feature)
│   │   │   ├── bookmark/                   # PostgreSQL via GORM
│   │   │   ├── cache/                      # Generic Redis cache operations
│   │   │   ├── healthcheck/
│   │   │   ├── link/                       # Redis-backed URL storage
│   │   │   ├── ratelimit/                  # Redis-backed rate limit counters
│   │   │   └── user/                       # PostgreSQL via GORM
│   │   │
│   │   └── model/                          # GORM domain entities
│   │       ├── base.go
│   │       ├── bookmark.go
│   │       ├── const.go
│   │       └── user.go
│   │
│   ├── infrastructure/                     # Dependency wiring (composition root)
│   │   ├── api.go                          # CreateAPI() bootstraps all deps → api.New()
│   │   ├── dbs.go                          # CreateSQLDB(), CreateRedisCon()
│   │   └── jwt.go                          # JWT generator/validator init
│   │
│   └── test/
│       ├── fixture/                        # Shared test case helpers
│       └── integration_test/               # Integration tests per feature
│
├── migrations/                             # Sequential SQL migration files (golang-migrate)
│   ├── 000001_add_user.{up,down}.sql
│   ├── 000002_add_bookmark.{up,down}.sql
│   └── 000003_add_code_shorten_and_code_shorten_encoded.{up,down}.sql
│
├── pkg/                                    # Shared, reusable packages
│   ├── common/                             # Standardized error & response helpers
│   ├── dbutils/                            # DB-specific error utilities
│   ├── encoding/                           # Base62 encoding
│   ├── jwtutils/                           # RSA-based JWT generator & validator
│   ├── logger/                             # Log level helpers (zerolog)
│   ├── redis/                              # Redis client & config
│   ├── sqldb/                              # GORM client, config, migrations
│   ├── utils/                              # Code generator, password hashing
│   └── validators/                         # Custom Gin validators (e.g., password_strength)
│
├── Dockerfile
├── Makefile
├── docker-compose.yaml                     # Production stack
├── docker-compose.dev.yaml                 # Local dev stack (PostgreSQL + Redis only)
├── go.mod
└── README.md
```

## Architecture

```
Request → Middleware (JWT / RateLimit) → Handler → Service → Repository → PostgreSQL / Redis
```

### Layer Responsibilities

| Layer | Path | Responsibility |
|---|---|---|
| API | `internal/api/` | Gin engine, route registration, middleware wiring |
| Handler | `internal/app/handler/` | HTTP request parsing, response formatting |
| Service | `internal/app/service/` | Business logic, interface + implementation |
| Repository | `internal/app/repository/` | Data access via GORM (PostgreSQL) or Redis |
| Model | `internal/app/model/` | GORM domain structs (`User`, `Bookmark`, …) |
| Infrastructure | `internal/infrastructure/` | Composition root — wires all dependencies |
| Pkg | `pkg/` | Reusable utilities shared across layers |

### Caching Pattern

The bookmark service uses a **decorator pattern** for transparent Redis caching:

```go
bookmarkSvc := bookmarkService.NewBookmarkService(bookmarkRepo)
bookmarkSvcWithCache := bookmarkService.NewBookmarkServiceWithCache(bookmarkSvc, redisCache)
```

### Rate Limiting

IP-based rate limiting is applied to all `/v1` routes:
- **Max requests**: 20 per 10-second window
- **Storage**: Redis (key format `rate_limit:<client_ip>`)
- Returns `429 Too Many Requests` when exceeded

## Features

- **Bookmark Management** — Full CRUD (create, list, update, delete) with Redis caching
- **URL Shortening** — Shorten URLs with auto-generated codes stored in Redis; redirect by code
- **User Management** — Register, login with bcrypt password hashing
- **JWT Authentication** — RSA-256 signed tokens; middleware protects private routes
- **Rate Limiting** — Redis-backed IP rate limiting on all `/v1` routes
- **Health Check** — Pings both PostgreSQL and Redis; returns service name & instance ID
- **Password Generation** — Cryptographically secure random password generation
- **Swagger UI** — Auto-generated API docs at `/swagger/index.html`
- **Docker Support** — Multi-stage Dockerfile; dev and production compose files
- **Database Migrations** — Sequential SQL migrations via `golang-migrate`
- **80% Coverage Gate** — Test suite enforces minimum coverage threshold

## Getting Started

### Prerequisites

- Go 1.25.5+
- Docker & Docker Compose
- `openssl` (for JWT RSA key generation)
- `swag` CLI: `go install github.com/swaggo/swag/cmd/swag@latest`

### 1. Generate RSA Keys (required for JWT)

```bash
make generate-rsa-key
# Creates private_key.pem and public_key.pem in the project root
```

### 2. Start Local Infrastructure

```bash
make dev-up       # Starts PostgreSQL (5432) and Redis (6379) via Docker
```

### 3. Run the Server

```bash
make dev-run      # Regenerates Swagger docs, then starts the server
# or
go run ./cmd/api/main.go
```

The server starts at `http://localhost:8080`. Swagger UI is at `http://localhost:8080/swagger/index.html`.

## Available Endpoints

### Public (rate-limited: 20 req / 10s per IP)

| Method | Path | Description |
|---|---|---|
| `GET` | `/health-check` | Returns health status, service name, instance ID |
| `GET` | `/swagger/*any` | Swagger UI (not rate-limited) |
| `GET` | `/v1/generate-password` | Generate a secure random password |
| `POST` | `/v1/links/shorten` | Shorten a URL |
| `GET` | `/v1/links/redirect/:code` | Redirect to the original URL by code |
| `POST` | `/v1/users/register` | Register a new user |
| `POST` | `/v1/users/login` | Login and receive a JWT token |

### Protected (requires `Authorization: Bearer <token>`)

| Method | Path | Description |
|---|---|---|
| `GET` | `/v1/self/info` | Get current user profile |
| `PUT` | `/v1/self/info` | Update current user profile |
| `POST` | `/v1/bookmarks` | Create a new bookmark |
| `GET` | `/v1/bookmarks` | List bookmarks (paginated) |
| `PUT` | `/v1/bookmarks/:id` | Update a bookmark |
| `DELETE` | `/v1/bookmarks/:id` | Delete a bookmark |

## Configuration

All configuration is provided via environment variables:

| Variable | Default | Description |
|---|---|---|
| `APP_PORT` | `:8080` | HTTP listen port |
| `SERVICE_NAME` | `bookmark_service` | Identifier shown in health check |
| `INSTANCE_ID` | auto UUID | Per-instance identifier |
| `APP_HOST_NAME` | `localhost:8080` | Swagger host |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `admin` | PostgreSQL user |
| `DB_PASSWORD` | `admin` | PostgreSQL password |
| `DB_NAME` | `bookmark_service` | PostgreSQL database |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | `""` | Redis password |
| `REDIS_DB` | `0` | Redis database index |

## Development

### Make Commands

```bash
# Local dev
make dev-up          # Start PostgreSQL + Redis
make dev-down        # Stop PostgreSQL + Redis
make dev-run         # Regenerate Swagger + run server

# Testing (enforces 80% coverage threshold)
make test            # Run all tests with coverage report
make clean           # Clear test cache + coverage output

# Code generation
make mock-gen        # Regenerate all mocks (go generate ./...)
make swag-gen        # Regenerate Swagger docs

# Database
make migrate                    # Apply all pending migrations
make new-schema name=<name>     # Create a new migration file pair

# Docker
make docker-build    # Build app + migration images
make docker-up       # Start production stack
make docker-down     # Stop production stack
make docker-release  # Build and push images to Docker Hub

# JWT keys
make generate-rsa-key   # Generate private_key.pem + public_key.pem
```

### Running Tests

```bash
# All tests (recommended — enforces coverage gate)
make test

# Specific package
go test ./internal/app/service/bookmark/...
go test ./internal/app/handler/user/...
go test ./internal/api/middleware/...

# Verbose
go test -v ./...
```

### Database Migrations

```bash
# Apply migrations
make migrate

# Create a new migration
make new-schema name=add_some_table
# Creates: migrations/NNNNNN_add_some_table.up.sql
#          migrations/NNNNNN_add_some_table.down.sql
```

### Docker

```bash
# Local dev (PostgreSQL + Redis only)
make dev-up
make dev-down

# Full production stack
make docker-up
make docker-down
```

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/gin-gonic/gin` | HTTP web framework |
| `gorm.io/gorm` + `gorm.io/driver/postgres` | ORM + PostgreSQL driver |
| `github.com/redis/go-redis/v9` | Redis client |
| `github.com/golang-jwt/jwt/v5` | JWT token generation & validation |
| `github.com/golang-migrate/migrate/v4` | SQL schema migrations |
| `github.com/swaggo/swag` + `gin-swagger` | Swagger doc generation & UI |
| `github.com/go-playground/validator/v10` | Request validation |
| `github.com/rs/zerolog` | Structured logging |
| `golang.org/x/crypto` | bcrypt password hashing |
| `github.com/google/uuid` | UUID generation |
| `github.com/kelseyhightower/envconfig` | Environment variable configuration |
| `github.com/stretchr/testify` | Test assertions & mocking |
| `go.uber.org/mock` | Mock generation |
| `github.com/alicebob/miniredis/v2` | In-memory Redis for tests |

## License

[Add your license here]