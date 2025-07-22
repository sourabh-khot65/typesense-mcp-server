// Package services defines service interfaces and implementations.
package services

import (
	"context"

	"typesense-mcp-server/internal/validation"

	"github.com/typesense/typesense-go/typesense/api"
)

// TypesenseService defines the interface for Typesense operations
type TypesenseService interface {
	// Search performs a search operation on the specified collection
	Search(ctx context.Context, params *validation.SearchParams) (*SearchResult, error)
	
	// GetCollections retrieves all collections from Typesense
	GetCollections(ctx context.Context) (*CollectionsResult, error)
	
	// Health checks the health of the Typesense service
	Health(ctx context.Context) error
	
	// Close closes the service and cleans up resources
	Close() error
}

// SearchResult represents the result of a search operation
type SearchResult struct {
	Found      int                      `json:"found"`
	Page       int                      `json:"page"`
	PerPage    int                      `json:"per_page"`
	Documents  []map[string]interface{} `json:"documents"`
	FacetCount []api.FacetCounts        `json:"facet_counts,omitempty"`
	QueryTime  float64                  `json:"query_time_ms,omitempty"`
}

// CollectionsResult represents the result of a collections operation
type CollectionsResult struct {
	Collections []*api.CollectionResponse `json:"collections"`
	Count       int                       `json:"count"`
}

// Document represents a search result document with metadata
type Document struct {
	Data       map[string]interface{} `json:"data"`
	TextMatch  *float64               `json:"text_match,omitempty"`
	Highlights *[]api.SearchHighlight `json:"highlights,omitempty"`
}
