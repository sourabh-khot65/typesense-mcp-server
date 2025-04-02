# Tacitbase MCP Server

This is a Machine Comprehension Protocol (MCP) server for Tacitbase that provides advanced search capabilities using both Tacitbase's native search and Typesense for enhanced search features.

## Project Structure

The project follows a clean, modular architecture:

```
.
├── models/      # Data structures and types
├── services/    # Business logic and external service interactions
├── handlers/    # Request handlers and routing logic
├── tools/       # Tool registration and configuration
├── main.go      # Application entry point
└── README.md    # Project documentation
```

## Packages

- `models`: Contains all data structures and types used across the application
- `services`: Implements the business logic and external service interactions
- `handlers`: Contains the request handlers that process incoming requests
- `tools`: Registers the tools available for searching candidates

## Features

- Basic keyword search with filtering and sorting
- Vector search for similarity-based matching
- Semantic search with natural language understanding
- Fallback to Tacitbase's native search when needed
- Tool registration for searching candidates

## Search Tools

### 1. Basic Search (`search_candidates`)
- Keyword-based search with typo tolerance
- Supports field-specific search
- Filtering and sorting capabilities
- Group by functionality
- Exact matching option

### 2. Vector Search (`vector_search_candidates`)
- Similarity search using vector embeddings
- Ideal for finding candidates with similar profiles
- Supports hybrid search combining vectors with filters

### 3. Semantic Search (`semantic_search_candidates`)
- Natural language understanding
- Automatic embedding generation
- Support for multiple embedding models (OpenAI, SBERT, E5)
- Hybrid search combining semantic understanding with filters

## Configuration

### Environment Variables

- `TACITBASE_AUTH_TOKEN`: Authentication token for Tacitbase API
- `TYPESENSE_API_KEY`: API key for Typesense
- `TYPESENSE_HOST`: Typesense host (default: localhost)
- `TYPESENSE_PORT`: Typesense port (default: 8108)
- `TYPESENSE_PROTOCOL`: Protocol for Typesense (default: http)

## Installation

1. Install dependencies:
```bash
go mod download
```

2. Build the server:
```bash
go build -o tb-mcp-server
```

3. Run the server:
```bash
./tb-mcp-server
```

## Example Usage

### Basic Search
```json
{
  "query": "golang developer",
  "search_fields": "skills,latest_experience",
  "filter_fields": "location:San Francisco",
  "sort_by": "latest_experience:desc",
  "page": 1,
  "per_page": 20
}
```

### Vector Search
```json
{
  "vector_query": "[0.1, 0.2, ..., 0.512]",
  "filter_fields": "years_of_experience:>5",
  "page": 1,
  "per_page": 20
}
```

### Semantic Search
```json
{
  "query": "experienced team lead with cloud architecture background",
  "embedding_model": "openai",
  "filter_fields": "years_of_experience:>8",
  "page": 1,
  "per_page": 20
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Development

To modify or extend the server:

1. Edit `main.go` to add new tools or modify existing ones
2. Run tests (if any)
3. Build and test the server
