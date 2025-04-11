package services

import (
	"context"
	"fmt"

	"tb-mcp-server/config"

	"github.com/sirupsen/logrus"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

// TypesenseService is an interface for the Typesense service
type TypesenseService interface {
	Search(ctx context.Context, collection string, request *api.SearchCollectionParams) (*api.SearchResult, error)
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
func (s *typesenseService) Search(ctx context.Context, collection string, request *api.SearchCollectionParams) (*api.SearchResult, error) {
	if !s.allowedCollections[collection] {
		logrus.Errorf("invalid collection: %s. Only 'candidates_candidates' and 'candidates_candidate-attachments' collections are allowed", collection)
		return nil, fmt.Errorf("invalid collection: %s. Only 'candidates_candidates' and 'candidates_candidate-attachments' collections are allowed", collection)
	}

	result, err := s.client.Collection(collection).Documents().Search(request)
	if err != nil {
		logrus.Errorf("failed to search documents: %v", err)
		return nil, fmt.Errorf("failed to search documents: %v", err)
	}

	return result, nil
}
