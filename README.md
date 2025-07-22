# Typesense MCP Server

A Model Control Protocol (MCP) server for interacting with Typesense, a fast, typo-tolerant search engine. This server provides a standardized interface for performing searches across any Typesense collection with enterprise-grade reliability, comprehensive error handling, and observability features.

## 🌟 Features

### Core Functionality
- **Generic search interface** for any Typesense collection
- **Complete Typesense API support** with all search parameters
- **Typo-tolerant search** with intelligent ranking
- **Advanced filtering and faceting** capabilities
- **Intelligent pagination** with validation
- **Collection metadata retrieval** with schema information

### Enterprise Features
- **Structured logging** with configurable formats (JSON/text)
- **Comprehensive error handling** with detailed error types
- **Input validation** with security-focused parameter checking
- **Health checks** with service monitoring
- **Graceful shutdown** with proper resource cleanup
- **Configuration validation** with environment-based setup
- **Request tracing** with duration tracking

## 🏗️ Architecture

This project follows Go best practices and patterns from well-structured open source projects:

```
├── cmd/server/          # Application entry points
├── internal/            # Private application code
│   ├── constants/       # Application constants
│   ├── errors/          # Custom error types
│   ├── logger/          # Structured logging
│   ├── config/          # Configuration management
│   ├── validation/      # Input validation
│   ├── services/        # Business logic
│   ├── handlers/        # Request handlers
│   ├── health/          # Health checks
│   └── tools/           # MCP tool registry
└── pkg/                 # Public library code (future)
```

## 📋 Configuration

The server supports comprehensive configuration through environment variables:

### Typesense Configuration
- `TYPESENSE_HOST`: Typesense server host (default: "localhost")
- `TYPESENSE_PORT`: Typesense server port (default: 8108)
- `TYPESENSE_PROTOCOL`: Protocol to use - http/https (default: "http")
- `TYPESENSE_API_KEY`: Typesense API key (default: "xyz")

### Logging Configuration
- `LOG_LEVEL`: Log level - debug/info/warn/error/fatal/panic (default: "info")
- `LOG_FORMAT`: Log format - json/text (default: "text")

## 🛠️ Available Tools

### typesense_search

Search documents in any Typesense collection with advanced capabilities.

**Parameters:**
- `collection` (required): Name of the Typesense collection to search in
- `q` (required): Search query - supports keywords, phrases, and natural language
- `query_by` (optional): Comma-separated list of fields to search in (default: "*")
- `filter_by` (optional): Filter expressions (e.g., "field:value", "num_field:>100")
- `page` (required): Page number for pagination (1-based, max: 10000)
- `per_page` (required): Number of results per page (1-100, default: 10)

**Response includes:**
- Found documents with metadata
- Text match scores and highlights
- Query performance metrics
- Facet counts (if requested)

### typesense_collections

Retrieve all collections with their metadata and schema information.

**Response includes:**
- Collection names and configurations
- Field schemas and types
- Index statistics
- Creation timestamps

## 🚀 Development

### Prerequisites

- **Go 1.23 or later**
- **Access to a Typesense server**

### Building

```bash
# Build the server
go build -o typesense-mcp-server ./cmd/server

# Or build with version info
go build -ldflags "-X main.version=$(git describe --tags)" -o typesense-mcp-server ./cmd/server
```

### Running

```bash
# Run with default configuration
./typesense-mcp-server

# Run with custom configuration
TYPESENSE_HOST=production-host \
TYPESENSE_API_KEY=your-key \
LOG_LEVEL=debug \
LOG_FORMAT=json \
./typesense-mcp-server
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

## 🔧 Error Handling

The server implements comprehensive error handling with structured error types:

- **Validation Errors**: Invalid input parameters
- **Configuration Errors**: Invalid configuration settings
- **Service Errors**: Typesense service issues
- **Client Errors**: Network or authentication problems
- **Internal Errors**: Unexpected server issues

All errors include:
- Error type and code
- Descriptive messages
- Contextual details
- Request tracing information

## 📊 Observability

### Logging
- **Structured logging** with consistent fields
- **Request tracing** with unique request IDs
- **Performance metrics** with operation durations
- **Error tracking** with stack traces

### Health Checks
- **Service health monitoring** with detailed status
- **Dependency checks** for Typesense connectivity
- **Performance metrics** for response times

## 🔒 Security

- **Input validation** with parameter sanitization
- **Configuration validation** with secure defaults
- **Error message sanitization** to prevent information leakage
- **Resource limits** with pagination controls

## 🤝 Contributing

1. Follow Go best practices and conventions
2. Add comprehensive tests for new features
3. Update documentation for API changes
4. Ensure all health checks pass
5. Use structured logging for observability

## 📄 License

This project follows standard open source licensing practices.
