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
	GetCollections(ctx context.Context) ([]*api.CollectionResponse, error)
}

// typesenseService is a service that provides a client for Typesense
type typesenseService struct {
	client *typesense.Client
}

// NewTypesenseService creates a new Typesense service
func NewTypesenseService(config *config.TypesenseConfig) TypesenseService {
	if config == nil {
		logrus.Error("typesense config is nil")
		return &typesenseService{client: nil}
	}

	if config.APIKey == "" {
		logrus.Error("typesense API key is empty")
		return &typesenseService{client: nil}
	}

	client := typesense.NewClient(
		typesense.WithServer(config.URL()),
		typesense.WithAPIKey(config.APIKey),
	)

	logrus.Infof("initialized Typesense client for %s", config.URL())
	
	return &typesenseService{
		client: client,
	}
}

// GetCollections gets all collections from Typesense
func (s *typesenseService) GetCollections(ctx context.Context) ([]*api.CollectionResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("typesense client is not initialized")
	}

	result, err := s.client.Collections().Retrieve()
	if err != nil {
		logrus.Errorf("failed to retrieve collections from Typesense: %v", err)
		return nil, fmt.Errorf("failed to retrieve collections: %w", err)
	}

	if result == nil {
		logrus.Warn("received nil result from Typesense collections API")
		return []*api.CollectionResponse{}, nil
	}

	return result, nil
}

// Search searches the collection for the given request
func (s *typesenseService) Search(ctx context.Context, collection string, request *api.SearchCollectionParams) (*api.SearchResult, error) {
	if s.client == nil {
		return nil, fmt.Errorf("typesense client is not initialized")
	}

	if collection == "" {
		return nil, fmt.Errorf("collection name cannot be empty")
	}

	if request == nil {
		return nil, fmt.Errorf("search request cannot be nil")
	}

	result, err := s.client.Collection(collection).Documents().Search(request)
	if err != nil {
		logrus.Errorf("failed to search documents in collection '%s': %v", collection, err)
		return nil, fmt.Errorf("failed to search collection '%s': %w", collection, err)
	}

	return result, nil
}
