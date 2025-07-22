package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"typesense-mcp-server/services"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"github.com/typesense/typesense-go/typesense/api"
)

// SearchHandler handles Typesense search requests
type SearchHandler struct {
	typesenseService services.TypesenseService
}

// NewSearchHandler creates a new instance of SearchHandler
func NewSearchHandler(typesenseService services.TypesenseService) *SearchHandler {
	return &SearchHandler{
		typesenseService: typesenseService,
	}
}

// SearchResponse represents the formatted search response
type SearchResponse struct {
	Found      int                      `json:"found"`
	Page       int                      `json:"page"`
	PerPage    int                      `json:"per_page"`
	Documents  []map[string]interface{} `json:"documents"`
	FacetCount []api.FacetCounts        `json:"facet_counts,omitempty"`
}

// GetTypesenseCollections handles the request to get all Typesense collections
func (h *SearchHandler) GetTypesenseCollections(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collections, err := h.typesenseService.GetCollections(ctx)
	if err != nil {
		return h.handleError("failed to get collections", err)
	}

	collectionsJSON, err := json.MarshalIndent(collections, "", "  ")
	if err != nil {
		return h.handleError("failed to marshal collections", err)
	}

	return mcp.NewToolResultText(string(collectionsJSON)), nil
}

// SearchInTypesenseCollection handles the search request for any Typesense collection
func (h *SearchHandler) SearchInTypesenseCollection(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract and validate collection name
	collection, err := h.extractCollectionName(request)
	if err != nil {
		return nil, err
	}

	// Build search parameters from request
	searchReq, err := h.buildSearchRequest(request)
	if err != nil {
		return h.handleError("failed to build search request", err)
	}

	// Perform search
	response, err := h.performSearch(ctx, collection, searchReq)
	if err != nil {
		return h.handleError(fmt.Sprintf("search failed for collection %s", collection), err)
	}

	// Format and return response
	return h.formatSearchResponse(response)
}

// extractCollectionName extracts and validates the collection name from the request
func (h *SearchHandler) extractCollectionName(request mcp.CallToolRequest) (string, error) {
	collection, ok := request.Params.Arguments["collection"].(string)
	if !ok || collection == "" {
		return "", fmt.Errorf("collection name is required and must be a non-empty string")
	}
	return collection, nil
}

// buildSearchRequest creates a SearchCollectionParams from the MCP request
func (h *SearchHandler) buildSearchRequest(request mcp.CallToolRequest) (*api.SearchCollectionParams, error) {
	// Remove the collection field from arguments as it's not part of the search params
	searchArgs := make(map[string]interface{})
	for key, value := range request.Params.Arguments {
		if key != "collection" {
			searchArgs[key] = value
		}
	}

	// Convert to SearchCollectionParams
	jsonData, err := json.Marshal(searchArgs)
	if err != nil {
		logrus.Errorf("failed to marshal search arguments: %v", err)
		return nil, fmt.Errorf("failed to marshal search arguments: %w", err)
	}

	var searchReq api.SearchCollectionParams
	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		logrus.Errorf("failed to unmarshal search request: %v", err)
		return nil, fmt.Errorf("failed to parse search parameters: %w", err)
	}

	return &searchReq, nil
}

// performSearch executes the search against Typesense
func (h *SearchHandler) performSearch(ctx context.Context, collection string, searchReq *api.SearchCollectionParams) (*api.SearchResult, error) {
	response, err := h.typesenseService.Search(ctx, collection, searchReq)
	if err != nil {
		logrus.Errorf("failed to search documents in collection %s: %v", collection, err)
		return nil, err
	}
	return response, nil
}

// formatSearchResponse formats the search result and returns it as a tool result
func (h *SearchHandler) formatSearchResponse(response *api.SearchResult) (*mcp.CallToolResult, error) {
	// Format the response
	result := formatTypesenseResults(response)

	// Convert to JSON string with indentation for better readability
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return h.handleError("failed to format search results", err)
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

// formatTypesenseResults formats the Typesense search response
func formatTypesenseResults(response *api.SearchResult) *SearchResponse {
	if response == nil {
		return &SearchResponse{
			Found:     0,
			Page:      1,
			PerPage:   0,
			Documents: make([]map[string]interface{}, 0),
		}
	}

	// Extract and format documents
	documents := extractDocuments(response)

	found := 0
	if response.Found != nil {
		found = *response.Found
	}

	var facetCounts []api.FacetCounts
	if response.FacetCounts != nil {
		facetCounts = *response.FacetCounts
	}

	return &SearchResponse{
		Found:      found,
		Page:       1, // Typesense uses offset-based pagination
		PerPage:    len(documents),
		Documents:  documents,
		FacetCount: facetCounts,
	}
}

// handleError provides consistent error handling and logging
func (h *SearchHandler) handleError(message string, err error) (*mcp.CallToolResult, error) {
	logrus.Errorf("%s: %v", message, err)
	return nil, fmt.Errorf("%s: %w", message, err)
}

// extractDocuments safely extracts and enriches documents from search hits
func extractDocuments(response *api.SearchResult) []map[string]interface{} {
	documents := make([]map[string]interface{}, 0)
	
	if response.Hits == nil {
		return documents
	}

	for _, hit := range *response.Hits {
		doc := extractDocument(hit)
		if doc != nil {
			documents = append(documents, doc)
		}
	}

	return documents
}

// extractDocument safely extracts a single document from a search hit
func extractDocument(hit api.SearchResultHit) map[string]interface{} {
	if hit.Document == nil {
		return nil
	}

	doc := make(map[string]interface{})
	
	// Copy original document data
	for key, value := range *hit.Document {
		doc[key] = value
	}

	// Add search metadata if available
	if hit.TextMatch != nil && *hit.TextMatch > 0 {
		doc["_text_match"] = *hit.TextMatch
	}
	
	if hit.Highlights != nil && len(*hit.Highlights) > 0 {
		doc["_highlights"] = *hit.Highlights
	}

	return doc
}
