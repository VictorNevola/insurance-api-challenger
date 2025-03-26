This project implements an insurance API using hexagonal architecture (ports and adapters) pattern in Go.

## API Routes

The API exposes the following endpoints (served on port 3000):

- Partner-related routes (managed by handler.go)
  - Partner management
  - Quote creation
  - Policy management

Swagger documentation is available at `/api/v1/docs`

## Dependencies

### External Services

- Insurance Provider API (integrates with challenge-api)

### Infrastructure

- **MongoDB**: Primary data store
- **Redis**: Caching layer

### Main Go Dependencies

- [Fiber](https://github.com/gofiber/fiber): HTTP server framework

## Setup and Running

### Prerequisites

Before running the application, make sure to set the correct `INSURANCE_PROVIDER_TOKEN` environment variable:

```bash
# Set the required token for the insurance provider API or set in docker-compose.yml
export INSURANCE_PROVIDER_TOKEN=your_token_here
```

```bash
# Run with Docker Compose
docker-compose up

```

### Environment Variables

- `MONGO_URL`: MongoDB connection string
- `MONGO_DATABASE`: MongoDB database name
- `REDIS_URL`: Redis connection string
- `INSURANCE_PROVIDER_URL`: URL to the insurance provider API
- `INSURANCE_PROVIDER_TOKEN`: Authentication token for insurance provider

## Project Structure

```
├── api                # HTTP handlers and API documentation
│   ├── docs           # Swagger documentation
│   └── web            # HTTP handlers
├── cmd                # Application entry points
├── configs            # Configuration adapters
├── internal           # Core application code
│   ├── domain         # Business logic and domain entities
│   ├── infra          # Infrastructure adapters
│   └── pkg            # Shared utilities
```
