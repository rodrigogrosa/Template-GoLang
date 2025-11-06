# Contributing to Template-GoLang

Thank you for your interest in contributing to this project!

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Make (optional but recommended)
- Docker & Docker Compose (for full stack testing)
- golangci-lint (for linting)
- swag (for generating API documentation)

### Getting Started

1. Fork and clone the repository:
```bash
git clone https://github.com/rodrigogrosa/Template-GoLang.git
cd Template-GoLang
```

2. Install dependencies:
```bash
make deps
```

3. Run the application locally:
```bash
make run
```

## Development Workflow

### Making Changes

1. Create a new branch for your feature or bugfix:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes following the project structure:
   - Domain logic goes in `internal/core/domain/`
   - Business logic goes in `internal/core/services/`
   - Adapters (HTTP, Kafka, DB) go in `internal/adapters/`
   - Shared utilities go in `pkg/`

3. Write tests for your changes:
```bash
make test
```

4. Format and lint your code:
```bash
make fmt
make lint
```

5. Update Swagger documentation if you modified API endpoints:
```bash
make swagger
```

## Code Standards

### Hexagonal Architecture

This project follows hexagonal (ports and adapters) architecture:

- **Domain Layer** (`internal/core/domain/`): Business entities and rules
- **Ports Layer** (`internal/core/ports/`): Interfaces defining contracts
- **Services Layer** (`internal/core/services/`): Application business logic
- **Adapters Layer** (`internal/adapters/`): External interface implementations

### Code Style

- Follow standard Go conventions
- Use `gofmt` and `goimports` for formatting
- Pass all `golangci-lint` checks
- Write meaningful variable and function names
- Add comments for exported functions and types

### Testing

- Write unit tests for all business logic
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

Example test:
```go
func TestItemService_CreateItem(t *testing.T) {
    tests := []struct {
        name    string
        item    *domain.Item
        wantErr bool
    }{
        {
            name: "valid item",
            item: &domain.Item{Name: "Test"},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Documentation

- Update README.md for significant changes
- Add Swagger annotations for API endpoints
- Document complex logic with inline comments
- Update catalog-info.yaml if changing service metadata

## API Documentation

This project uses Swagger/OpenAPI for API documentation.

### Adding Swagger Annotations

For new endpoints, add annotations above the handler:

```go
// CreateItem godoc
// @Summary Create a new item
// @Description Create a new item with the provided details
// @Tags items
// @Accept json
// @Produce json
// @Param item body CreateItemRequest true "Item details"
// @Success 201 {object} domain.Item
// @Failure 400 {object} ErrorResponse
// @Router /v1/items [post]
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
    // implementation
}
```

Regenerate docs:
```bash
make swagger
```

## Pull Request Process

1. Ensure all tests pass: `make test`
2. Ensure linting passes: `make lint`
3. Update documentation as needed
4. Create a pull request with:
   - Clear description of changes
   - Reference to any related issues
   - Screenshots for UI changes (if applicable)
5. Wait for review and address feedback

## Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Example:
```
feat(items): add pagination support to list endpoint

- Add limit and offset query parameters
- Update tests for pagination
- Update API documentation

Closes #123
```

## Running the Full Stack

Use Docker Compose to run the complete stack:

```bash
make docker-compose-up
```

This starts:
- The microservice
- Kafka (with Zookeeper)
- Jaeger (for tracing)
- Prometheus (for metrics)

Access:
- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/
- Jaeger UI: http://localhost:16686
- Prometheus: http://localhost:9090

## Observability

### Logging

Use the structured logger:

```go
import "github.com/rodrigogrosa/Template-GoLang/pkg/logger"

logger.Logger.Info().
    Str("user_id", userID).
    Msg("User created")
```

### Metrics

Add custom metrics:

```go
import "github.com/rodrigogrosa/Template-GoLang/pkg/metrics"

myCounter := promauto.NewCounter(
    prometheus.CounterOpts{
        Name: "my_events_total",
        Help: "Total number of my events",
    },
)

myCounter.Inc()
```

### Tracing

Use context for distributed tracing:

```go
ctx, span := tracer.GetTracer().Start(ctx, "operation-name")
defer span.End()

// operation code
```

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about the architecture
- Suggestions for improvements

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
