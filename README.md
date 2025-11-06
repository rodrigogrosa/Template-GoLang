# Template-GoLang

Production-ready Golang microservice template with hexagonal architecture, aligned with Backstage.io standards.

## Features

### Architecture
- **Hexagonal Architecture** (Ports & Adapters)
- Clean separation of concerns
- Domain-driven design principles
- Dependency injection

### API
- **REST API** with OpenAPI 3.0 specification
- `/health` - Health check endpoint
- `/v1/items` - CRUD operations for items
- `/metrics` - Prometheus metrics
- **Swagger UI** - Interactive API documentation at `/swagger/`

### Observability
- **Structured Logging** - Using zerolog with JSON output
- **Prometheus Metrics** - HTTP requests, business metrics
- **OpenTelemetry Tracing** - Distributed tracing with Jaeger

### Configuration
- **Environment Variables** - 12-factor app configuration
- TLS support for secure communications
- Configurable timeouts and settings

### Reliability
- **Graceful Shutdown** - Proper signal handling
- Context propagation
- Connection pooling ready

### Testing
- Unit tests for all layers
- Test coverage reports
- Race condition detection

### DevOps
- **Dockerfile** - Multi-stage build for minimal image size
- **Makefile** - Common development tasks
- **Docker Compose** - Local development environment
- **golangci-lint** - Comprehensive linting

### Security
- Input validation using go-playground/validator
- **JWT Authentication** - Middleware for token validation
- **TLS Support** - HTTPS configuration
- Security scanning ready

### Messaging
- **Kafka Integration** - Producer and consumer examples
- Event-driven architecture support
- Async message processing

### Backstage Integration
- `catalog-info.yaml` - Backstage catalog definition
- API definition integration
- Service metadata

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── core/
│   │   ├── domain/              # Domain entities and business rules
│   │   ├── ports/               # Interfaces (ports)
│   │   └── services/            # Business logic (use cases)
│   ├── adapters/
│   │   ├── http/                # HTTP handlers and server
│   │   ├── repository/          # Data persistence adapters
│   │   └── kafka/               # Kafka producer/consumer
│   └── config/                  # Configuration management
├── pkg/
│   ├── logger/                  # Structured logging
│   ├── metrics/                 # Prometheus metrics
│   ├── tracer/                  # OpenTelemetry tracing
│   ├── validation/              # Input validation
│   └── middleware/              # HTTP middlewares
├── docs/                        # Swagger documentation
├── Dockerfile                   # Container image definition
├── Makefile                     # Build and development commands
├── docker-compose.yml           # Local development stack
├── catalog-info.yaml            # Backstage catalog definition
└── README.md                    # This file
```

## Getting Started

### Prerequisites

- Go 1.21+
- Docker (optional, for containerization)
- Docker Compose (optional, for local stack)
- Make (optional, for using Makefile)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/rodrigogrosa/Template-GoLang.git
cd Template-GoLang
```

2. Install dependencies:
```bash
make deps
```

3. Generate Swagger documentation:
```bash
make swagger
```

### Running the Application

#### Local Development

```bash
# Run directly
make run

# Or with custom environment
export PORT=8080
export LOG_LEVEL=debug
go run ./cmd/api/main.go
```

#### Using Docker

```bash
# Build image
make docker-build

# Run container
make docker-run
```

#### Using Docker Compose (Full Stack)

```bash
# Start all services (app, Kafka, Jaeger, Prometheus)
make docker-compose-up

# Stop all services
make docker-compose-down
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Format code
make fmt
```

## Configuration

Configuration is managed through environment variables. See `.env.example` for all available options.

### Key Configuration Options

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `HOST` | HTTP server host | `0.0.0.0` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `ENVIRONMENT` | Environment (development, production) | `development` |
| `TRACING_ENABLED` | Enable OpenTelemetry tracing | `false` |
| `JAEGER_ENDPOINT` | Jaeger collector endpoint | `""` |
| `KAFKA_ENABLED` | Enable Kafka integration | `false` |
| `KAFKA_BROKERS` | Kafka broker addresses | `localhost:9092` |
| `JWT_SECRET` | JWT signing secret | `""` |
| `TLS_ENABLED` | Enable TLS/HTTPS | `false` |
| `TLS_CERT_FILE` | TLS certificate file path | `""` |
| `TLS_KEY_FILE` | TLS private key file path | `""` |

## API Documentation

Once the service is running, access the Swagger UI at:
- http://localhost:8080/swagger/

### API Endpoints

#### Health Check
```bash
GET /health
```

#### Items API
```bash
# Create item
POST /v1/items
{
  "name": "Item Name",
  "description": "Item Description"
}

# List items
GET /v1/items?limit=10&offset=0

# Get item by ID
GET /v1/items/{id}

# Update item
PUT /v1/items/{id}
{
  "name": "Updated Name",
  "description": "Updated Description"
}

# Delete item
DELETE /v1/items/{id}
```

#### Metrics
```bash
GET /metrics
```

## Observability

### Logging

Logs are output in JSON format (production) or console format (development):

```json
{
  "level": "info",
  "time": "2024-01-01T12:00:00Z",
  "message": "HTTP request",
  "method": "GET",
  "path": "/v1/items",
  "status": 200,
  "duration": 15
}
```

### Metrics

Prometheus metrics are available at `/metrics`:
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request duration histogram
- `items_created_total` - Total items created
- `items_updated_total` - Total items updated
- `items_deleted_total` - Total items deleted
- `kafka_messages_published_total` - Kafka messages published
- `kafka_messages_consumed_total` - Kafka messages consumed

### Tracing

When enabled, traces are sent to Jaeger. View them at:
- http://localhost:16686

## Security

### JWT Authentication

To enable JWT authentication, set the `JWT_SECRET` environment variable:

```bash
export JWT_SECRET=your-secret-key
```

Requests must include a Bearer token:
```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/v1/items
```

### TLS/HTTPS

Enable TLS by setting:
```bash
export TLS_ENABLED=true
export TLS_CERT_FILE=/path/to/cert.pem
export TLS_KEY_FILE=/path/to/key.pem
```

## Kafka Integration

Enable Kafka integration:
```bash
export KAFKA_ENABLED=true
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=items
```

Events published:
- `item.created` - When an item is created
- `item.updated` - When an item is updated
- `item.deleted` - When an item is deleted

## Development

### Make Commands

```bash
make help          # Show all available commands
make deps          # Install dependencies
make swagger       # Generate Swagger docs
make build         # Build the application
make run           # Run the application
make test          # Run tests
make test-coverage # Run tests with coverage
make lint          # Run linter
make fmt           # Format code
make clean         # Clean build artifacts
```

### Adding New Features

1. **Domain Model** - Add entities in `internal/core/domain/`
2. **Ports** - Define interfaces in `internal/core/ports/`
3. **Services** - Implement business logic in `internal/core/services/`
4. **Adapters** - Implement adapters in `internal/adapters/`
5. **Tests** - Add tests for all layers
6. **Documentation** - Update Swagger annotations

## Backstage Integration

This template is ready for Backstage.io integration. The `catalog-info.yaml` file defines:
- Component metadata
- API definitions
- Dependencies
- Links to documentation and monitoring

Import into Backstage:
```yaml
# In your Backstage catalog
catalog:
  locations:
    - type: url
      target: https://github.com/rodrigogrosa/Template-GoLang/blob/main/catalog-info.yaml
```

## License

MIT

## Support

For issues and questions, please open an issue on GitHub.
