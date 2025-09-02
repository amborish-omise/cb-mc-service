# MasterCom Service

A Go service following the OmisePayments Go template structure for handling MasterCom case filing and document management.

## Project Structure

This project follows the OmisePayments Go template structure:

```
├── cmd/
│   ├── http/           # HTTP-only server
│   ├── grpc/           # gRPC-only server  
│   └── http-grpc/      # Combined HTTP + gRPC server
├── internal/            # Internal application code
├── pkg/                 # Public packages
├── specs/               # API specifications (gRPC, OpenAPI)
├── .buildkite/          # CI/CD configuration
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yml   # Local development setup
└── Makefile            # Build and development commands
```

## Quick Start

### Prerequisites

- Go 1.25+
- Docker and Docker Compose
- Make

### Environment Setup

1. Copy the environment template:
   ```bash
   cp .env.example .env
   ```

2. Update the `.env` file with your configuration:
   ```bash
   APP_NAME=mastercom-service
   DEBUGGER_PORT=2345
   APP_GRPC_SERVER_PORT=50000
   APP_HTTP_SERVER_PORT=8080
   ```

### Running the Service

#### HTTP Server Only
```bash
make run-http
```

#### gRPC Server Only
```bash
make run-grpc
```

#### Combined HTTP + gRPC Server
```bash
make run-http-grpc
```

### Development

#### Debug Mode
```bash
make debug-http      # Debug HTTP server
make debug-grpc      # Debug gRPC server
make debug-http-grpc # Debug combined server
```

#### Testing
```bash
make test            # Run tests
make check           # Run linting and checks
```

#### Code Generation
```bash
make gen-proto      # Generate gRPC code from protobuf
make gen-openapi    # Generate OpenAPI specs from protobuf
```

### Docker Commands

```bash
make build-dev       # Build development image
make build-test      # Build testing image
make build           # Build production image
make release-http    # Build release image for HTTP
make release-grpc    # Build release image for gRPC
make release-http-grpc # Build release image for combined
```

## API Endpoints

### Health Checks
- `GET /__ops/ping` - Template-style health check (returns "pong")
- `GET /health` - Detailed health check with Datadog tracing

### Case Management
- `POST /api/v6/cases` - Create a new case
- `GET /api/v6/cases` - List all cases
- `GET /api/v6/cases/:id` - Get a specific case
- `PUT /api/v6/cases/:id` - Update a case
- `DELETE /api/v6/cases/:id` - Delete a case

### Document Management
- `POST /api/v6/documents` - Upload a document
- `GET /api/v6/documents/:id` - Get a specific document
- `DELETE /api/v6/documents/:id` - Delete a document

## Configuration

The service uses Viper for configuration management. Configuration can be set via:

- Environment variables
- Configuration files
- Command line flags

## Observability

- **Logging**: Structured logging with Datadog integration
- **Tracing**: Distributed tracing with Datadog APM
- **Metrics**: Prometheus metrics (planned)
- **Health Checks**: Built-in health check endpoints

## Development Workflow

1. **Local Development**: Use `make run-*` for local development
2. **Testing**: Use `make test` for running tests
3. **Linting**: Use `make check` for code quality checks
4. **Building**: Use `make build-*` for Docker builds
5. **Deployment**: Use `make release-*` for production builds

## Webhook Development

### Testing Webhooks

Use the provided sample payload for testing:

```bash
# Test the webhook endpoint
curl -X POST http://localhost:8080/api/v6/webhooks/ethoca \
  -H "Content-Type: application/json" \
  -d @docs/sample-webhook-payload.json

# Check webhook health
curl http://localhost:8080/api/v6/webhooks/ethoca/health

# Get webhook statistics
curl http://localhost:8080/api/v6/webhooks/ethoca/stats
```

### Webhook Documentation

For detailed webhook integration information, see:
- [Ethoca Webhook Documentation](docs/ETHOCA_WEBHOOK.md)
- [Sample Webhook Payload](docs/sample-webhook-payload.json)
- [System Diagrams](docs/system-diagrams/) - Comprehensive visual documentation of architecture and workflows

## Contributing

1. Follow the existing code structure
2. Add tests for new functionality
3. Run `make check` before committing
4. Follow the Go template patterns established
5. Include webhook tests when adding new webhook functionality

## License

[Add your license information here]
