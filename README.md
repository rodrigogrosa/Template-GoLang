# Golang Microservice Template

A production-ready Golang microservice template with hexagonal architecture, comprehensive observability, and modern DevOps practices.

## 🚀 Features

- **Clean Architecture**: Hexagonal (ports & adapters) architecture for maintainability
- **RESTful API**: HTTP REST API with `/health` and `/v1/items` endpoints
- **OpenAPI 3.0**: Full API specification with Swagger UI at `/swagger/`
- **Observability**:
  - Structured logging with `zerolog`
  - Prometheus metrics exposed on port 9090
  - OpenTelemetry distributed tracing
- **Event Streaming**: Kafka producer/consumer for domain events
- **Security**:
  - Input validation on all endpoints
  - JWT authentication middleware (configurable)
  - TLS/HTTPS support
  - Non-root Docker container
- **Configuration**: Environment variable-based configuration
- **Graceful Shutdown**: Proper handling of OS signals
- **Testing**: Unit and integration tests with good coverage
- **DevOps Ready**: Dockerfile, Makefile, docker-compose, and linting

## 📁 Project Structure

```
.
├── cmd/
│   └── api/              # Application entry point
│       └── main.go
├── internal/             # Private application code
│   ├── domain/           # Business logic and entities
│   │   ├── item.go
│   │   ├── service.go
│   │   └── errors.go
│   ├── ports/            # Interfaces (hexagonal architecture)
│   │   ├── ports.go
│   │   └── service.go
│   ├── adapters/         # External implementations
│   │   ├── http/         # HTTP handlers and server
│   │   ├── repository/   # Data persistence
│   │   └── kafka/        # Event streaming
│   └── config/           # Configuration management
├── pkg/                  # Reusable packages
│   ├── logger/           # Structured logging
│   ├── middleware/       # HTTP middleware
│   └── validator/        # Input validation
├── api/                  # API specifications
│   └── openapi.yaml      # OpenAPI 3.0 spec
├── test/                 # Integration tests
│   └── integration/
├── Dockerfile            # Container image definition
├── docker-compose.yml    # Local development stack
├── Makefile              # Build automation
├── .golangci.yml         # Linter configuration
└── README.md

```

## 🛠️ Prerequisites

- Go 1.21 or later
- Docker and Docker Compose (for containerized deployment)
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/rodrigogrosa/Template-GoLang.git
cd Template-GoLang
```

### 2. Install dependencies

```bash
make install
```

### 3. Run locally

```bash
make run
```

The API will be available at:
- API: http://localhost:8080
- Health: http://localhost:8080/health
- Swagger UI: http://localhost:8080/swagger/
- Metrics: http://localhost:9090/metrics

### 4. Run with Docker Compose

```bash
docker-compose up -d
```

This starts:
- API service on port 8080
- Kafka for event streaming
- Jaeger for distributed tracing (UI on http://localhost:16686)
- Prometheus for metrics collection

## 📝 Configuration

Configuration is managed through environment variables. See `.env.example` for all available options:

```bash
cp .env.example .env
# Edit .env with your settings
```

### Key Configuration Options

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 8080 | HTTP server port |
| `LOG_LEVEL` | info | Logging level (debug, info, warn, error) |
| `LOG_FORMAT` | json | Log format (json, console) |
| `KAFKA_ENABLED` | false | Enable Kafka integration |
| `AUTH_ENABLED` | false | Enable JWT authentication |
| `TLS_ENABLED` | false | Enable HTTPS |
| `METRICS_ENABLED` | true | Enable Prometheus metrics |
| `TRACING_ENABLED` | true | Enable OpenTelemetry tracing |

## 🔌 API Endpoints

### Health Check
```bash
GET /health
```

### Items API

#### Create Item
```bash
POST /v1/items
Content-Type: application/json

{
  "name": "Example Item",
  "description": "Item description",
  "price": 19.99
}
```

#### Get Item
```bash
GET /v1/items/{id}
```

#### List Items
```bash
GET /v1/items?limit=10&offset=0
```

#### Update Item
```bash
PUT /v1/items/{id}
Content-Type: application/json

{
  "name": "Updated Item",
  "description": "Updated description",
  "price": 29.99
}
```

#### Delete Item
```bash
DELETE /v1/items/{id}
```

## 🔐 Authentication

JWT authentication can be enabled by setting `AUTH_ENABLED=true` and providing a `JWT_SECRET`.

To make authenticated requests, include the JWT token in the Authorization header:

```bash
curl -H "Authorization: Bearer <your-token>" http://localhost:8080/v1/items
```

## 🧪 Testing

### Run unit tests
```bash
make test
```

### Run integration tests
```bash
make test-integration
```

### Run with coverage
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔍 Code Quality

### Linting
```bash
make lint
```

### Formatting
```bash
make fmt
```

### Security scanning
```bash
make security
```

## 🐳 Docker

### Build image
```bash
make docker-build
```

### Run container
```bash
make docker-run
```

## 📊 Observability

### Metrics
Prometheus metrics are exposed at `http://localhost:9090/metrics`

Key metrics include:
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - HTTP request latency

### Tracing
When running with docker-compose, Jaeger UI is available at `http://localhost:16686`

### Logging
Structured JSON logs are written to stdout and include:
- Request/response details
- Error traces
- Performance metrics

## 🔄 Kafka Events

When Kafka is enabled, the service publishes domain events:

- `item.created` - When an item is created
- `item.updated` - When an item is updated
- `item.deleted` - When an item is deleted

Example event consumer is included in the main application.

## 🏗️ Architecture

This template follows **Hexagonal Architecture** (Ports & Adapters):

- **Domain Layer**: Core business logic, independent of external concerns
- **Ports**: Interfaces defining contracts
- **Adapters**: Concrete implementations (HTTP, Kafka, Repository)

Benefits:
- ✅ Testable: Easy to mock dependencies
- ✅ Flexible: Swap implementations without changing business logic
- ✅ Maintainable: Clear separation of concerns

## 📦 Building for Production

```bash
# Build optimized binary
make build

# Build Docker image
docker build -t golang-microservice:latest .

# Run in production mode
docker run -p 8080:8080 -p 9090:9090 \
  -e LOG_LEVEL=warn \
  -e AUTH_ENABLED=true \
  -e JWT_SECRET=your-secret \
  golang-microservice:latest
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- Built with best practices from the Go community
- Inspired by Domain-Driven Design and Clean Architecture principles
- Uses industry-standard libraries and tools
