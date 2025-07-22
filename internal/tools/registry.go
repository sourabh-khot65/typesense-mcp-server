// Package tools provides MCP tool registration and management.
package tools

import (
	"fmt"

	"typesense-mcp-server/internal/config"
	"typesense-mcp-server/internal/constants"
	"typesense-mcp-server/internal/handlers"
	"typesense-mcp-server/internal/logger"
	"typesense-mcp-server/internal/services"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Registry manages MCP tool registration
type Registry struct {
	config  *config.Config
	logger  *logger.Logger
	service services.TypesenseService
	handler *handlers.SearchHandler
}

// New creates a new tool registry
func New(cfg *config.Config, log *logger.Logger) (*Registry, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	
	if log == nil {
		log = logger.NewDefault()
	}

	// Create Typesense service
	typesenseService, err := services.NewTypesenseService(&cfg.Typesense, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create typesense service: %w", err)
	}

	// Create handlers
	searchHandler := handlers.NewSearchHandler(typesenseService, log)

	return &Registry{
		config:  cfg,
		logger:  log,
		service: typesenseService,
		handler: searchHandler,
	}, nil
}

// RegisterTools registers all tools with the MCP server
func (r *Registry) RegisterTools(s *server.MCPServer) error {
	r.logger.WithComponent(constants.ComponentMain).
		Info("registering MCP tools")

	// Register collections tool
	if err := r.registerCollectionsTool(s); err != nil {
		return fmt.Errorf("failed to register collections tool: %w", err)
	}

	// Register search tool
	if err := r.registerSearchTool(s); err != nil {
		return fmt.Errorf("failed to register search tool: %w", err)
	}

	r.logger.WithComponent(constants.ComponentMain).
		WithField("tool_count", 2).
		Info("MCP tools registered successfully")

	return nil
}

// registerCollectionsTool registers the collections tool
func (r *Registry) registerCollectionsTool(s *server.MCPServer) error {
	collectionsTool := mcp.NewTool(
		constants.ToolTypesenseCollections,
		mcp.WithDescription("Get all collections with their details such as schema etc. from Typesense. "+
			"This tool retrieves metadata about all available collections in the Typesense instance, "+
			"including collection names, field schemas, and configuration details."),
	)

	s.AddTool(collectionsTool, r.handler.GetTypesenseCollections)
	
	r.logger.WithComponent(constants.ComponentMain).
		WithField("tool_name", constants.ToolTypesenseCollections).
		Debug("collections tool registered")

	return nil
}

// registerSearchTool registers the search tool
func (r *Registry) registerSearchTool(s *server.MCPServer) error {
	searchTool := mcp.NewTool(
		constants.ToolTypesenseSearch,
		mcp.WithDescription("Search documents in a Typesense collection using powerful search capabilities. "+
			"Supports typo-tolerant search, filtering, faceting, pagination, and more. "+
			"Returns matching documents with optional highlights and metadata."),
		mcp.WithString(constants.ParamCollection,
			mcp.Required(),
			mcp.Description("Name of the Typesense collection to search in. Must be a valid collection name."),
		),
		mcp.WithString(constants.ParamQuery,
			mcp.Required(),
			mcp.Description("Search query. Can be keywords, phrases, or natural language queries. "+
				"Supports wildcard (*) and phrase searches."),
		),
		mcp.WithString(constants.ParamQueryBy,
			mcp.Description("Comma-separated list of fields to search in. Use '*' to search all fields."),
			mcp.DefaultString(constants.DefaultQueryBy),
		),
		mcp.WithString(constants.ParamFilterBy,
			mcp.Description("Filter expressions to narrow down results. "+
				"Examples: 'field:value', 'num_field:>100', 'category:=electronics'"),
		),
		mcp.WithNumber(constants.ParamPage,
			mcp.Description(fmt.Sprintf("Page number for pagination (1-based). Must be between 1 and %d.", constants.MaxPage)),
			mcp.DefaultNumber(float64(constants.DefaultPage)),
			mcp.Required(),
		),
		mcp.WithNumber(constants.ParamPerPage,
			mcp.Description(fmt.Sprintf("Number of results per page. Must be between %d and %d.", 
				constants.MinPageSize, constants.MaxPageSize)),
			mcp.DefaultNumber(float64(constants.DefaultPageSize)),
			mcp.Required(),
		),
	)

	s.AddTool(searchTool, r.handler.SearchInTypesenseCollection)
	
	r.logger.WithComponent(constants.ComponentMain).
		WithField("tool_name", constants.ToolTypesenseSearch).
		Debug("search tool registered")

	return nil
}

// Close closes the registry and cleans up resources
func (r *Registry) Close() error {
	r.logger.WithComponent(constants.ComponentMain).
		Info("closing tool registry")

	if r.service != nil {
		if err := r.service.Close(); err != nil {
			r.logger.WithComponent(constants.ComponentMain).
				WithError(err).
				Error("failed to close typesense service")
			return err
		}
	}

	return nil
}

// GetService returns the Typesense service for health checks or other operations
func (r *Registry) GetService() services.TypesenseService {
	return r.service
}
