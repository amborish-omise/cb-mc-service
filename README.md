# MasterCom Service with Datadog Integration

A Go-based REST API service for MasterCom case management and document handling with comprehensive Datadog logging and tracing.

## Features

- Case filing and management
- Document upload and management
- RESTful API endpoints
- JSON logging with Datadog integration
- Distributed tracing with Datadog APM
- CORS support
- Health check endpoint
- Performance profiling

## Project Structure

```
CB-mc-backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point with Datadog setup
├── internal/
│   ├── config/
│   │   ├── config.go        # Configuration management
│   │   └── datadog.go       # Datadog configuration
│   ├── handlers/
│   │   ├── case_handlers.go     # Case-related HTTP handlers with tracing
│   │   └── document_handlers.go # Document-related HTTP handlers with tracing
│   ├── models/
│   │   ├── case.go          # Case data models
│   │   └── document.go      # Document data models
│   └── services/
│       ├── case_service.go      # Case business logic
│       └── document_service.go  # Document business logic
├── pkg/
│   ├── logger/
│   │   └── datadog.go       # Datadog-integrated logger
│   └── middleware/
│       ├── cors.go          # CORS middleware
│       ├── datadog.go       # Datadog tracing middleware
│       └── logger.go        # Logging middleware
├── tests/                   # Integration tests
├── mastercom-swagger.yaml   # OpenAPI specification
├── docker-compose.yml       # Docker Compose with Datadog agent
├── .env.example            # Environment variables template
├── README.md               # This file
├── Makefile                # Build and run commands
├── Dockerfile              # Container configuration
├── .gitignore              # Git ignore rules
├── go.mod                  # Go module file
└── go.sum                  # Dependencies checksum
```

## Datadog Integration

This service includes comprehensive Datadog integration for:

### Tracing
- Automatic HTTP request tracing
- Custom spans for business operations
- Trace correlation with logs
- Performance monitoring

### Logging
- Structured JSON logging
- Trace ID and Span ID correlation
- Error tracking and monitoring
- Custom log fields for business metrics

### Profiling
- CPU profiling
- Memory profiling
- Block and mutex profiling
- Performance optimization insights

## API Endpoints

### Health Check
- `GET /health` - Service health status with trace information

### Cases
- `POST /api/v6/cases` - Create a new case
- `GET /api/v6/cases` - List all cases (with pagination)
- `GET /api/v6/cases/:id` - Get a specific case
- `PUT /api/v6/cases/:id` - Update a case
- `DELETE /api/v6/cases/:id` - Delete a case

### Documents
- `POST /api/v6/documents` - Upload a document
- `GET /api/v6/documents/:id` - Get a specific document
- `DELETE /api/v6/documents/:id` - Delete a document

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Datadog account and API key

### Local Development

1. Clone the repository
2. Copy environment variables:
   ```bash
   cp .env.example .env
   ```

3. Update `.env` with your Datadog API key:
   ```bash
   DD_API_KEY=your-actual-datadog-api-key
   ```

4. Install dependencies:
   ```bash
   go mod tidy
   ```

5. Run with Docker Compose (includes Datadog agent):
   ```bash
   docker-compose up --build
   ```

6. Or run locally (requires Datadog agent):
   ```bash
   make run
   ```

### Environment Variables

#### Application
- `ENVIRONMENT` - Set to "production" for production mode (default: "development")
- `PORT` - Server port (default: "8080")
- `LOG_LEVEL` - Logging level (default: "info")

#### Datadog
- `DD_ENABLED` - Enable Datadog integration (default: "true")
- `DD_SERVICE` - Service name (default: "mastercom-service")
- `DD_ENV` - Environment (default: "development")
- `DD_VERSION` - Service version (default: "1.0.0")
- `DD_AGENT_HOST` - Datadog agent host (default: "localhost")
- `DD_AGENT_PORT` - Datadog agent port (default: "8126")
- `DD_TRACE_SAMPLE_RATE` - Trace sampling rate (default: "1.0")
- `DD_API_KEY` - Your Datadog API key

### Development

#### Running tests
```bash
make test
```

#### Building for production
```bash
make build
```

#### Cleaning build artifacts
```bash
make clean
```

## Datadog Dashboard

Once the service is running, you can view:

1. **APM Traces** - Distributed tracing for all API requests
2. **Logs** - Structured logs with trace correlation
3. **Profiles** - Performance profiling data
4. **Metrics** - Custom business metrics

### Key Metrics to Monitor

- Request latency and throughput
- Error rates by endpoint
- Case creation and processing times
- Document upload success rates
- Database operation performance

## API Examples

### Create a Case with Tracing
```bash
curl -X POST http://localhost:8080/api/v6/cases \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: $(uuidgen)" \
  -d '{
    "caseType": "PRE_ARBITRATION",
    "primaryAccountNumber": "4111111111111111",
    "transactionAmount": 100.00,
    "transactionCurrency": "USD",
    "transactionDate": "2024-01-01T00:00:00Z",
    "transactionId": "123456789",
    "reasonCode": "10.1",
    "filingAs": "ISSUER",
    "filingIca": "123456",
    "filedAgainstIca": "654321"
  }'
```

### Upload a Document with Tracing
```bash
curl -X POST http://localhost:8080/api/v6/documents \
  -F "file=@document.pdf" \
  -F "caseId=case-id-here" \
  -F "description=Supporting documentation" \
  -F "uploadedBy=test-user" \
  -H "X-Request-ID: $(uuidgen)"
```

### Health Check with Trace Info
```bash
curl http://localhost:8080/health
```

Response includes trace information:
```json
{
  "status": "ok",
  "service": "mastercom-service",
  "version": "v1.0.0",
  "dd_trace_id": "1234567890",
  "dd_span_id": "0987654321"
}
```

## Monitoring and Alerting

### Recommended Alerts

1. **High Error Rate** - Alert when error rate > 5%
2. **High Latency** - Alert when p95 latency > 2s
3. **Service Down** - Alert when health check fails
4. **Document Upload Failures** - Alert when upload success rate < 95%

### Custom Metrics

The service automatically tracks:
- Case creation/update/deletion rates
- Document upload/download rates
- API endpoint usage
- Business transaction volumes

## Troubleshooting

### Datadog Agent Issues

1. Check agent status:
   ```bash
   docker-compose logs datadog-agent
   ```

2. Verify agent connectivity:
   ```bash
   curl http://localhost:8126/info
   ```

3. Check trace ingestion:
   ```bash
   curl http://localhost:8126/health
   ```

### Service Issues

1. Check service logs:
   ```bash
   docker-compose logs mastercom-service
   ```

2. Verify environment variables:
   ```bash
   docker-compose exec mastercom-service env | grep DD_
   ```

## License

This project is licensed under the MIT License.
