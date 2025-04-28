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
		logrus.Errorf("failed to get collections: %v", err)
		return nil, fmt.Errorf("failed to get collections: %v", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("%v", collections)), nil
}

// Search handles the search request for any Typesense collection
func (h *SearchHandler) SearchInTypesenseCollection(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract collection name from arguments
	collection, ok := request.Params.Arguments["collection"].(string)
	if !ok {
		return nil, fmt.Errorf("collection name is required")
	}

	// Create search parameters
	var searchReq api.SearchCollectionParams
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		logrus.Errorf("failed to marshal search arguments: %v", err)
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		logrus.Errorf("failed to unmarshal search request: %v", err)
		return nil, fmt.Errorf("failed to unmarshal search request: %v", err)
	}

	// Perform search using Typesense
	response, err := h.typesenseService.Search(ctx, collection, &searchReq)
	if err != nil {
		logrus.Errorf("failed to search documents in collection %s: %v", collection, err)
		return nil, fmt.Errorf("search failed: %v", err)
	}

	// Format the response
	result := formatTypesenseResults(response)

	// Convert to JSON string with indentation for better readability
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal search results: %v", err)
		return nil, fmt.Errorf("failed to format results: %v", err)
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

// formatTypesenseResults formats the Typesense search response
func formatTypesenseResults(response *api.SearchResult) *SearchResponse {
	if response == nil || response.Found == nil {
		return &SearchResponse{
			Found:     0,
			Page:      1,
			PerPage:   10,
			Documents: make([]map[string]interface{}, 0),
		}
	}

	// Format documents
	documents := make([]map[string]interface{}, 0)
	if response.Hits != nil {
		for _, hit := range *response.Hits {
			if hit.Document != nil {
				doc := *hit.Document
				// Add search score to document
				if hit.TextMatch != nil {
					doc["_text_match"] = *hit.TextMatch
				}
				if hit.Highlights != nil {
					doc["_highlights"] = hit.Highlights
				}
				documents = append(documents, doc)
			}
		}
	}

	return &SearchResponse{
		Found:      *response.Found,
		Page:       1, // Typesense uses offset-based pagination
		PerPage:    len(documents),
		Documents:  documents,
		FacetCount: *response.FacetCounts,
	}
}
