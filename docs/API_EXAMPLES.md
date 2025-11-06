# API Examples

This document provides practical examples for using the microservice API.

## Prerequisites

- Go 1.24+
- Docker and Docker Compose (optional)
- curl or similar HTTP client

## Starting the Service

### Option 1: Run Locally
```bash
# Install dependencies
make install

# Run the service
make run

# Or build and run
make build
./bin/api
```

### Option 2: Run with Docker
```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

### Option 3: Run Full Stack with Docker Compose
```bash
# Start all services (API, Kafka, Jaeger, Prometheus)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down
```

## API Examples

### 1. Health Check

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy"
}
```

### 2. Create an Item

```bash
curl -X POST http://localhost:8080/v1/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "description": "16-inch, M3 Max",
    "price": 2499.99
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "MacBook Pro",
  "description": "16-inch, M3 Max",
  "price": 2499.99,
  "created_at": "2023-11-06T10:00:00Z",
  "updated_at": "2023-11-06T10:00:00Z"
}
```

### 3. Get an Item

```bash
# Replace {id} with actual item ID
curl http://localhost:8080/v1/items/{id}
```

### 4. List Items

```bash
# List with default pagination
curl http://localhost:8080/v1/items

# List with custom pagination
curl "http://localhost:8080/v1/items?limit=20&offset=0"
```

### 5. Update an Item

```bash
curl -X PUT http://localhost:8080/v1/items/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro (Updated)",
    "description": "16-inch, M3 Max, Space Black",
    "price": 2599.99
  }'
```

### 6. Delete an Item

```bash
curl -X DELETE http://localhost:8080/v1/items/{id}
```

## Using JWT Authentication

When authentication is enabled (`AUTH_ENABLED=true`), include a JWT token in requests:

```bash
# Set your JWT token
TOKEN="your-jwt-token-here"

# Make authenticated request
curl http://localhost:8080/v1/items \
  -H "Authorization: Bearer $TOKEN"
```

## Monitoring

### Prometheus Metrics

```bash
# View all metrics
curl http://localhost:9090/metrics
```

### Swagger UI

Open your browser and navigate to:
```
http://localhost:8080/swagger/
```

### Jaeger Tracing (with docker-compose)

Open your browser:
```
http://localhost:16686
```

## Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration
```
