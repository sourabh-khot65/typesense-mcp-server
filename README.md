# Typesense MCP Server

A Model Control Protocol (MCP) server for interacting with Typesense, a fast, typo-tolerant search engine. This server provides a standardized interface for performing searches across any Typesense collection.

## Features

- Generic search interface for any Typesense collection
- Support for all Typesense search parameters
- Typo-tolerant search
- Filtering and faceting support
- Pagination

## Configuration

The server can be configured using the following environment variables:

- `TYPESENSE_HOST`: Typesense server host (default: "localhost")
- `TYPESENSE_PORT`: Typesense server port (default: 8108)
- `TYPESENSE_PROTOCOL`: Protocol to use (http/https) (default: "http")
- `TYPESENSE_API_KEY`: Typesense API key (default: "xyz")

## Available Tools

### typesense_search

Search documents in any Typesense collection.

Parameters:
- `collection` (required): Name of the Typesense collection to search in
- `q` (required): Search query to find documents
- `query_by` (optional): Comma-separated list of fields to search in (default: "*")
- `filter_by` (optional): Filter expressions (e.g., "field:value", "num_field:>100")
- `page` (required): Page number for pagination (1-based)
- `per_page` (required): Number of results per page (default: 10, max: 100)

## Development

### Prerequisites

- Go 1.23 or later
- Access to a Typesense server

### Building

```bash
go build -o typesense-mcp-server
```

### Running

```bash
./typesense-mcp-server
```
