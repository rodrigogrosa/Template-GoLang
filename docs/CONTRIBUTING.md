# Contributing to Golang Microservice Template

Thank you for your interest in contributing! This document provides guidelines for contributing to this project.

## Development Setup

1. **Fork and clone the repository**
```bash
git clone https://github.com/rodrigogrosa/Template-GoLang.git
cd Template-GoLang
```

2. **Install dependencies**
```bash
make install
```

3. **Run tests**
```bash
make test
```

4. **Build the project**
```bash
make build
```

## Code Standards

### Go Code Style

- Follow standard Go formatting (`gofmt`, `goimports`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### Running Linters

```bash
make lint
```

Fix any issues reported by the linter before submitting a PR.

### Testing

- Write unit tests for new functionality
- Maintain or improve test coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate

Example:
```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "result", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MyFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("MyFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Project Structure

```
.
├── cmd/              # Application entry points
├── internal/         # Private application code
│   ├── domain/       # Business entities and logic
│   ├── ports/        # Interface definitions
│   ├── adapters/     # External implementations
│   │   ├── http/     # HTTP handlers
│   │   ├── repository/
│   │   └── kafka/
│   └── config/       # Configuration
├── pkg/              # Public reusable packages
├── api/              # API specs (OpenAPI)
├── docs/             # Documentation
└── test/             # Integration tests
```

## Pull Request Process

1. **Create a feature branch**
```bash
git checkout -b feature/your-feature-name
```

2. **Make your changes**
- Write clean, documented code
- Add tests for new functionality
- Update documentation as needed

3. **Run quality checks**
```bash
make check  # Runs fmt, vet, lint, and test
```

4. **Commit your changes**
```bash
git add .
git commit -m "Brief description of changes"
```

Use conventional commit messages:
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Adding tests
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

5. **Push and create PR**
```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Architecture Guidelines

This project follows **Hexagonal Architecture** (Ports & Adapters):

- **Domain Layer** (`internal/domain`): Pure business logic, no external dependencies
- **Ports** (`internal/ports`): Interfaces that define contracts
- **Adapters** (`internal/adapters`): Implementations of ports (HTTP, DB, Kafka)
- **Services** (`internal/services`): Application services coordinating domain logic

### Adding a New Feature

1. **Define domain entities** in `internal/domain`
2. **Create port interfaces** in `internal/ports`
3. **Implement adapters** in `internal/adapters`
4. **Add service layer** in `internal/services`
5. **Wire up in main** (`cmd/api/main.go`)

## Security Guidelines

- Never commit secrets or credentials
- Validate all user inputs
- Use parameterized queries for database access
- Keep dependencies up to date
- Run security scans: `make security`

## Documentation

- Update README.md for user-facing changes
- Update API_EXAMPLES.md for API changes
- Add inline comments for complex logic
- Update OpenAPI spec for API changes

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about the codebase
- Documentation improvements

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
