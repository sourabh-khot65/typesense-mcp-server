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
	// searchTool := mcp.NewTool("mcp_tacitbase_search_candidates",
	// 	mcp.WithDescription("Search for candidates in Tacitbase using Typesense's powerful search capabilities. "+
	// 		"Supports keyword search with typo tolerance, field-specific search, filtering, sorting, and grouping."),
	// 	mcp.WithString("q",
	// 		mcp.Required(),
	// 		mcp.Description("Search query to find candidates. Can be keywords, phrases, or natural language queries."),
	// 	),
	// 	// mcp.WithString("query_by",
	// 	// 	mcp.Description("Fields to search in. Available fields: first_name, last_name, email, phone, skills, latest_experience, highest_education, description, if not provided, set *"),
	// 	// 	mcp.DefaultString("*"),
	// 	// ),
	// 	// mcp.WithString("filter_by",
	// 	// 	mcp.Description("Filter expressions in format: field:value. Example: skills:python, location:San Francisco, separate multiple filters with &&"),
	// 	// ),
	// 	mcp.WithNumber("page",
	// 		mcp.Description("Page number for pagination (1-based)"),
	// 		mcp.DefaultNumber(1),
	// 	),
	// 	mcp.WithNumber("per_page",
	// 		mcp.Description("Number of results per page (default: 10, max: 100)"),
	// 		mcp.DefaultNumber(10),
	// 	),
	// )
	// s.AddTool(searchTool, searchHandler.HandleSearch)

	// // Add search tool for attachments
	// attachmentsSearchTool := mcp.NewTool("mcp_tacitbase_search_attachments",
	// 	mcp.WithDescription("Search for candidate attachments (resumes, portfolios, etc.) in the candidates_candidate-attachments collection. "+
	// 		"Supports full-text search within document content and metadata."),
	// 	mcp.WithString("q",
	// 		mcp.Required(),
	// 		mcp.Description("Search query to find attachments. Can search in file names and content."),
	// 	),
	// 	// mcp.WithString("query_by",
	// 	// 	mcp.Description("Fields to search in. Available fields: name, content, record_id, if not provided, set *"),
	// 	// 	mcp.DefaultString("*"),
	// 	// ),
	// 	// mcp.WithString("filter_by",
	// 	// 	mcp.Description("Filter expressions in format: field:value. Example: record_id:candidate_id, created_at:>2024-01-01"),
	// 	// ),
	// 	mcp.WithNumber("page",
	// 		mcp.Description("Page number for pagination (1-based)"),
	// 	),
	// 	mcp.WithNumber("per_page",
	// 		mcp.Description("Number of results per page (default: 10, max: 50)"),
	// 	),
	// )
	// s.AddTool(attachmentsSearchTool, searchHandler.HandleAttachmentsSearch)

	// Add staging search tool for candidates
	stagingSearchTool := mcp.NewTool("mcp_tacitbase_staging_search_candidates",
		mcp.WithDescription("Search for candidates directly in Tacitbase staging environment"),
		mcp.WithString("q",
			mcp.Required(),
			mcp.Description("Search query to find candidates. Can be keywords, phrases, or natural language queries."),
		),
		// mcp.WithString("query_by",
		// 	mcp.Description("Fields to search in. Available fields: first_name, last_name, email, phone, skills, latest_experience, highest_education, description, if not provided, set *"),
		// 	mcp.DefaultString("*"),
		// ),
		// mcp.WithString("filter_by",
		// 	mcp.Description("Filter expressions in format: field:value. Example: years_of_experience:>5, location:San Francisco"),
		// 	mcp.Description("Only consider when filtering search results based on the specific field"),
		// ),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (1-based)"),
			mcp.DefaultNumber(1),
			mcp.Required(),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of results per page (default: 10, max: 100)"),
			mcp.DefaultNumber(10),
			mcp.Required(),
		),
	)
	s.AddTool(stagingSearchTool, searchHandler.HandleStagingSearch)
}
