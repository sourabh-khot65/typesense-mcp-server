package tools

import (
	"typesense-mcp-server/config"
	"typesense-mcp-server/handlers"
	"typesense-mcp-server/services"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all tools with the server
func RegisterTools(s *server.MCPServer) {
	// Initialize configuration
	typesenseConfig := config.NewTypesenseConfig()

	// Initialize services
	typesenseService := services.NewTypesenseService(typesenseConfig)

	// Initialize handlers
	searchHandler := handlers.NewSearchHandler(typesenseService)
	collectionsHandler := handlers.NewCollectionsHandler(typesenseService)

	collectionsTool := mcp.NewTool("typesense_collections",
		mcp.WithDescription("Get all collections with their details such as schema etc. from Typesense"),
	)
	s.AddTool(collectionsTool, searchHandler.GetTypesenseCollections)

	// Add search tool for Typesense collections
	searchTool := mcp.NewTool("typesense_search",
		mcp.WithDescription("Search documents in a Typesense collection using powerful search capabilities. "+
			"Supports typo-tolerant search, filtering, faceting, and more."),
		mcp.WithString("collection",
			mcp.Required(),
			mcp.Description("Name of the Typesense collection to search in."),
		),
		mcp.WithString("q",
			mcp.Required(),
			mcp.Description("Search query. Can be keywords, phrases, or natural language queries."),
		),
		mcp.WithString("query_by",
			mcp.Description("Comma-separated list of fields to search in."),
			mcp.DefaultString("*"),
		),
		mcp.WithString("filter_by",
			mcp.Description("Filter expressions. Example: field:value, num_field:>100"),
		),
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
	s.AddTool(searchTool, searchHandler.SearchInTypesenseCollection)

	getCollectionTool := mcp.NewTool("typesense_get_collection",
		mcp.WithDescription("Get a Typesense collection by name."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the collection.")),
	)
	s.AddTool(getCollectionTool, collectionsHandler.GetCollection)

	createCollectionTool := mcp.NewTool("typesense_create_collection",
		mcp.WithDescription("Create a new Typesense collection."),
		mcp.WithObject("schema", mcp.Required(), mcp.Description("Collection schema as per Typesense API.")),
	)
	s.AddTool(createCollectionTool, collectionsHandler.CreateCollection)

	updateCollectionTool := mcp.NewTool("typesense_update_collection",
		mcp.WithDescription("Update a Typesense collection by name."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the collection.")),
		mcp.WithObject("schema", mcp.Required(), mcp.Description("Collection update schema as per Typesense API.")),
	)
	s.AddTool(updateCollectionTool, collectionsHandler.UpdateCollection)

	deleteCollectionTool := mcp.NewTool("typesense_delete_collection",
		mcp.WithDescription("Delete a Typesense collection by name."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the collection.")),
	)
	s.AddTool(deleteCollectionTool, collectionsHandler.DeleteCollection)
}
