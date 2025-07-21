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
	GetCollection(ctx context.Context, name string) (*api.CollectionResponse, error)
	CreateCollection(ctx context.Context, schema *api.CollectionSchema) (*api.CollectionResponse, error)
	UpdateCollection(ctx context.Context, name string, schema *api.CollectionUpdateSchema) (*api.CollectionResponse, error)
	DeleteCollection(ctx context.Context, name string) (*api.CollectionResponse, error)
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

// GetCollections gets all collections from Typesense
func (s *typesenseService) GetCollections(ctx context.Context) ([]*api.CollectionResponse, error) {
	result, err := s.client.Collections().Retrieve()
	if err != nil {
		logrus.Errorf("failed to get collections: %v", err)
		return nil, fmt.Errorf("failed to get collections: %v", err)
	}
	return result, nil
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

func (s *typesenseService) GetCollection(ctx context.Context, name string) (*api.CollectionResponse, error) {
	result, err := s.client.Collection(name).Retrieve()
	if err != nil {
		logrus.Errorf("failed to get collection %s: %v", name, err)
		return nil, fmt.Errorf("failed to get collection %s: %v", name, err)
	}
	return result, nil
}

func (s *typesenseService) CreateCollection(ctx context.Context, schema *api.CollectionSchema) (*api.CollectionResponse, error) {
	result, err := s.client.Collections().Create(schema)
	if err != nil {
		logrus.Errorf("failed to create collection: %v", err)
		return nil, fmt.Errorf("failed to create collection: %v", err)
	}
	return result, nil
}

func (s *typesenseService) UpdateCollection(ctx context.Context, name string, schema *api.CollectionUpdateSchema) (*api.CollectionResponse, error) {
	result, err := s.client.Collection(name).Update(schema)
	if err != nil {
		logrus.Errorf("failed to update collection %s: %v", name, err)
		return nil, fmt.Errorf("failed to update collection %s: %v", name, err)
	}
	return result, nil
}

func (s *typesenseService) DeleteCollection(ctx context.Context, name string) (*api.CollectionResponse, error) {
	result, err := s.client.Collection(name).Delete()
	if err != nil {
		logrus.Errorf("failed to delete collection %s: %v", name, err)
		return nil, fmt.Errorf("failed to delete collection %s: %v", name, err)
	}
	return result, nil
}
