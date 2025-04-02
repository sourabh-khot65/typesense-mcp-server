package tools

import (
	"tb-mcp-server/config"
	"tb-mcp-server/handlers"
	"tb-mcp-server/services"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all tools with the server
func RegisterTools(s *server.MCPServer) {
	// Initialize configuration
	typesenseConfig := config.NewTypesenseConfig()

	// Initialize services
	tacitbaseService := services.NewTacitbaseService()
	typesenseService := services.NewTypesenseService(typesenseConfig)

	// Initialize handlers
	searchHandler := handlers.NewSearchHandler(tacitbaseService, typesenseService)

	// Add basic search tool for candidates
	searchTool := mcp.NewTool("search_candidates",
		mcp.WithDescription("Search for candidates in Tacitbase"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query to find candidates"),
		),
		mcp.WithString("search_fields",
			mcp.Description("Comma-separated list of fields to search in (e.g., 'first_name,last_name,skills')"),
		),
		mcp.WithString("filter_fields",
			mcp.Description("Comma-separated list of fields to filter on"),
		),
		mcp.WithString("sort_by",
			mcp.Description("Comma-separated list of fields to sort by (e.g., 'latest_experience:desc')"),
		),
		mcp.WithString("group_by",
			mcp.Description("Comma-separated list of fields to group by"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of results per page"),
		),
	)
	s.AddTool(searchTool, searchHandler.HandleSearch)

	// Add vector search tool for candidates
	vectorSearchTool := mcp.NewTool("vector_search_candidates",
		mcp.WithDescription("Search for candidates using vector embeddings"),
		mcp.WithString("vector_query",
			mcp.Required(),
			mcp.Description("Vector query for similarity search"),
		),
		mcp.WithString("search_fields",
			mcp.Description("Comma-separated list of fields to search in"),
		),
		mcp.WithString("filter_fields",
			mcp.Description("Comma-separated list of fields to filter on"),
		),
		mcp.WithString("sort_by",
			mcp.Description("Comma-separated list of fields to sort by"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of results per page"),
		),
	)
	s.AddTool(vectorSearchTool, searchHandler.HandleVectorSearch)

	// Add semantic search tool for candidates
	semanticSearchTool := mcp.NewTool("semantic_search_candidates",
		mcp.WithDescription("Search for candidates using semantic search"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Natural language query for semantic search"),
		),
		mcp.WithArray("search_fields",
			mcp.Description("Array list of fields to search in"),
		),
		mcp.WithString("embedding_field",
			mcp.Description("Field containing pre-computed embeddings"),
		),
		mcp.WithString("embedding_model",
			mcp.Description("Model to use for generating embeddings (e.g., 'openai', 'sbert', 'e5')"),
		),
		mcp.WithString("filter_fields",
			mcp.Description("Comma-separated list of fields to filter on"),
		),
		mcp.WithString("sort_by",
			mcp.Description("Comma-separated list of fields to sort by"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of results per page"),
		),
	)
	s.AddTool(semanticSearchTool, searchHandler.HandleSemanticSearch)

	// Add search tool for attachments
	attachmentsSearchTool := mcp.NewTool("search_attachments",
		mcp.WithDescription("Search for candidate attachments"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query to find attachments"),
		),
		mcp.WithString("search_fields",
			mcp.Description("Comma-separated list of fields to search in (e.g., 'name,content')"),
		),
		mcp.WithString("filter_fields",
			mcp.Description("Comma-separated list of fields to filter on"),
		),
		mcp.WithString("sort_by",
			mcp.Description("Comma-separated list of fields to sort by (e.g., 'created_at:desc')"),
		),
		mcp.WithString("record_id",
			mcp.Description("Filter by record ID"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of results per page"),
		),
	)
	s.AddTool(attachmentsSearchTool, searchHandler.HandleAttachmentsSearch)
}
