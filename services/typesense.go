package services

import (
	"context"
	"fmt"

	"typesense-mcp-server/config"

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
	client *typesense.Client
}

// NewTypesenseService creates a new Typesense service
func NewTypesenseService(config *config.TypesenseConfig) TypesenseService {
	client := typesense.NewClient(
		typesense.WithServer(config.URL()),
		typesense.WithAPIKey(config.APIKey),
	)

	return &typesenseService{
		client: client,
	}
}

// Search searches the collection for the given request
func (s *typesenseService) Search(ctx context.Context, collection string, request *api.SearchCollectionParams) (*api.SearchResult, error) {
	result, err := s.client.Collection(collection).Documents().Search(request)
	if err != nil {
		logrus.Errorf("failed to search documents in collection %s: %v", collection, err)
		return nil, fmt.Errorf("failed to search documents in collection %s: %v", collection, err)
	}

	return result, nil
}
