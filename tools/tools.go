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
	searchTool := mcp.NewTool("mcp_tacitbase_search_candidates",
		mcp.WithDescription("Search for candidates in Tacitbase using Typesense's powerful search capabilities. "+
			"Supports keyword search with typo tolerance, field-specific search, filtering, sorting, and grouping."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query to find candidates. Can be keywords, phrases, or natural language queries."),
		),
		mcp.WithArray("search_fields",
			mcp.Description("Fields to search in. Available fields: first_name, last_name, email, phone, skills, latest_experience, highest_education, description"),
		),
		mcp.WithArray("filter_fields",
			mcp.Description("Filter expressions in format: field:value. Example: years_of_experience:>5, location:San Francisco"),
		),
		mcp.WithArray("sort_by",
			mcp.Description("Fields to sort by with optional direction. Format: field:direction. Example: latest_experience:desc, skills:asc"),
		),
		mcp.WithArray("group_by",
			mcp.Description("Fields to group results by. Example: skills, highest_education"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (1-based)"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of results per page (default: 10, max: 100)"),
		),
	)
	s.AddTool(searchTool, searchHandler.HandleSearch)

	// Add search tool for attachments
	attachmentsSearchTool := mcp.NewTool("mcp_tacitbase_search_attachments",
		mcp.WithDescription("Search for candidate attachments (resumes, portfolios, etc.) in the candidates_candidate-attachments collection. "+
			"Supports full-text search within document content and metadata."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query to find attachments. Can search in file names and content."),
		),
		mcp.WithArray("search_fields",
			mcp.Description("Fields to search in. Available fields: name, content, record_id"),
		),
		mcp.WithArray("filter_fields",
			mcp.Description("Filter expressions for attachments. Example: record_id:123, created_at:>2024-01-01"),
		),
		mcp.WithArray("sort_by",
			mcp.Description("Fields to sort attachments by. Example: created_at:desc"),
		),
		mcp.WithString("record_id",
			mcp.Description("Filter attachments by specific candidate record ID"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (1-based)"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of results per page (default: 10, max: 50)"),
		),
	)
	s.AddTool(attachmentsSearchTool, searchHandler.HandleAttachmentsSearch)

	// Add staging search tool for candidates
	stagingSearchTool := mcp.NewTool("mcp_tacitbase_staging_search_candidates",
		mcp.WithDescription("Search for candidates directly in Tacitbase staging environment, bypassing Typesense."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query to find candidates. Can be keywords, phrases, or natural language queries."),
		),
		mcp.WithArray("search_fields",
			mcp.Description("Fields to search in. Available fields: first_name, last_name, email, phone, skills, latest_experience, highest_education, description"),
		),
		mcp.WithArray("filter_fields",
			mcp.Description("Filter expressions in format: field:value. Example: years_of_experience:>5, location:San Francisco"),
		),
		mcp.WithArray("sort_by",
			mcp.Description("Fields to sort by with optional direction. Format: field:direction. Example: latest_experience:desc, skills:asc"),
		),
		mcp.WithArray("group_by",
			mcp.Description("Fields to group results by. Example: skills, highest_education"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (1-based)"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of results per page (default: 10, max: 100)"),
		),
	)
	s.AddTool(stagingSearchTool, searchHandler.HandleStagingSearch)
}
