package services

import (
	"context"
	"fmt"

	"tb-mcp-server/config"
	"tb-mcp-server/models"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

// TypesenseService is an interface for the Typesense service
type TypesenseService interface {
	Search(ctx context.Context, request *models.SearchRequest) (*models.SearchResponse, error)
}

// typesenseService is a service that provides a client for Typesense
type typesenseService struct {
	client             *typesense.Client
	allowedCollections map[string]bool
}

// NewTypesenseService creates a new Typesense service
func NewTypesenseService(config *config.TypesenseConfig) TypesenseService {
	client := typesense.NewClient(
		typesense.WithServer(config.URL()),
		typesense.WithAPIKey(config.APIKey),
	)

	return &typesenseService{
		client: client,
		allowedCollections: map[string]bool{
			"candidates_candidates":            true,
			"candidates_candidate-attachments": true,
		},
	}
}

// Search searches the collection for the given request
func (s *typesenseService) Search(ctx context.Context, request *models.SearchRequest) (*models.SearchResponse, error) {
	if !s.allowedCollections[request.Collection] {
		return nil, fmt.Errorf("invalid collection: %s. Only 'candidates_candidates' and 'candidates_candidate-attachments' collections are allowed", request.Collection)
	}

	result, err := s.client.Collection(request.Collection).Documents().Search(&api.SearchCollectionParams{
		Q:        request.Query,
		QueryBy:  request.QueryBy,
		FilterBy: &request.FilterBy,
		SortBy:   &request.SortBy,
		Page:     &request.Page,
		PerPage:  &request.PerPage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %v", err)
	}

	return convertToSearchResponse(result)
}

// convertToSearchResponse converts the search result to a search response
func convertToSearchResponse(result *api.SearchResult) (*models.SearchResponse, error) {
	if result == nil {
		return &models.SearchResponse{}, nil
	}

	found := 0
	if result.Found != nil {
		found = *result.Found
	}

	page := 1
	if result.Page != nil {
		page = *result.Page
	}

	perPage := 10
	if result.RequestParams != nil {
		perPage = result.RequestParams.PerPage
	}

	response := &models.SearchResponse{
		Found:   found,
		Page:    page,
		PerPage: perPage,
		Hits:    make([]map[string]interface{}, 0),
	}

	if result.Hits != nil {
		for _, hit := range *result.Hits {
			if hit.Document != nil {
				response.Hits = append(response.Hits, *hit.Document)
			}
		}
	}

	return response, nil
}
