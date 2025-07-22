// Package services provides service implementations.
package services

import (
	"context"
	"encoding/json"
	"time"

	"typesense-mcp-server/internal/config"
	"typesense-mcp-server/internal/constants"
	"typesense-mcp-server/internal/errors"
	"typesense-mcp-server/internal/logger"
	"typesense-mcp-server/internal/validation"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

// typesenseService implements the TypesenseService interface
type typesenseService struct {
	client *typesense.Client
	config *config.TypesenseConfig
	logger *logger.Logger
}

// NewTypesenseService creates a new Typesense service instance
func NewTypesenseService(cfg *config.TypesenseConfig, log *logger.Logger) (TypesenseService, error) {
	if cfg == nil {
		return nil, errors.NewConfigurationError("typesense config cannot be nil")
	}

	if log == nil {
		log = logger.NewDefault()
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, errors.WrapError(err, errors.ErrorTypeConfiguration, "invalid typesense configuration")
	}

	// Create Typesense client
	client := typesense.NewClient(
		typesense.WithServer(cfg.URL()),
		typesense.WithAPIKey(cfg.APIKey),
	)

	service := &typesenseService{
		client: client,
		config: cfg,
		logger: log,
	}

	// Log successful initialization
	log.WithComponent(constants.ComponentService).
		WithField("url", cfg.URL()).
		Info("typesense service initialized")

	return service, nil
}

// Search performs a search operation on the specified collection
func (s *typesenseService) Search(ctx context.Context, params *validation.SearchParams) (*SearchResult, error) {
	if params == nil {
		return nil, errors.NewValidationError("search parameters cannot be nil")
	}

	start := time.Now()
	
	entry := s.logger.WithComponent(constants.ComponentService).
		WithField(constants.LogFieldOperation, "search").
		WithField(constants.LogFieldCollection, params.Collection).
		WithField(constants.LogFieldQuery, params.Query).
		WithField(constants.LogFieldPage, params.Page).
		WithField(constants.LogFieldPageSize, params.PerPage)

	// Convert params to API request
	searchReq, err := s.buildSearchRequest(params)
	if err != nil {
		duration := time.Since(start)
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).WithError(err).Error("failed to build search request")
		return nil, errors.WrapError(err, errors.ErrorTypeValidation, "failed to build search request")
	}

	// Perform search
	response, err := s.performSearch(ctx, params.Collection, searchReq)
	if err != nil {
		duration := time.Since(start)
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).WithError(err).Error("search operation failed")
		return nil, err
	}

	// Format response
	result := s.formatSearchResponse(response, params.Page, params.PerPage)
	
	duration := time.Since(start)
	entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).
		WithField(constants.LogFieldResultCount, result.Found).
		Info("search operation completed")

	return result, nil
}

// GetCollections retrieves all collections from Typesense
func (s *typesenseService) GetCollections(ctx context.Context) (*CollectionsResult, error) {
	start := time.Now()
	
	entry := s.logger.WithComponent(constants.ComponentService).
		WithField(constants.LogFieldOperation, "get_collections")

	if s.client == nil {
		err := errors.NewClientNotInitializedError()
		entry.WithError(err).Error("client not initialized")
		return nil, err
	}

	result, err := s.client.Collections().Retrieve()
	if err != nil {
		duration := time.Since(start)
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).WithError(err).Error("failed to retrieve collections")
		return nil, errors.WrapError(err, errors.ErrorTypeService, "failed to retrieve collections from typesense")
	}

	if result == nil {
		entry.Warn("received nil collections result from typesense")
		result = []*api.CollectionResponse{}
	}

	collectionsResult := &CollectionsResult{
		Collections: result,
		Count:       len(result),
	}

	duration := time.Since(start)
	entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).
		WithField("collection_count", collectionsResult.Count).
		Info("collections retrieval completed")

	return collectionsResult, nil
}

// Health checks the health of the Typesense service
func (s *typesenseService) Health(ctx context.Context) error {
	if s.client == nil {
		return errors.NewClientNotInitializedError()
	}

	// Try to retrieve collections as a health check
	_, err := s.client.Collections().Retrieve()
	if err != nil {
		return errors.WrapError(err, errors.ErrorTypeService, "typesense health check failed")
	}

	return nil
}

// Close closes the service and cleans up resources
func (s *typesenseService) Close() error {
	s.logger.WithComponent(constants.ComponentService).Info("closing typesense service")
	// Typesense client doesn't require explicit closing, but this is for interface compliance
	return nil
}

// buildSearchRequest converts validation params to Typesense API request
func (s *typesenseService) buildSearchRequest(params *validation.SearchParams) (*api.SearchCollectionParams, error) {
	searchArgs := params.ToMap()

	// Convert to SearchCollectionParams via JSON marshaling/unmarshaling
	// This ensures proper type conversion and validation
	jsonData, err := json.Marshal(searchArgs)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrorTypeInternal, "failed to marshal search arguments")
	}

	var searchReq api.SearchCollectionParams
	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		return nil, errors.WrapError(err, errors.ErrorTypeInternal, "failed to unmarshal search request")
	}

	return &searchReq, nil
}

// performSearch executes the search request against Typesense
func (s *typesenseService) performSearch(ctx context.Context, collection string, searchReq *api.SearchCollectionParams) (*api.SearchResult, error) {
	if s.client == nil {
		return nil, errors.NewClientNotInitializedError()
	}

	response, err := s.client.Collection(collection).Documents().Search(searchReq)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrorTypeService, "typesense search failed")
	}

	return response, nil
}

// formatSearchResponse converts Typesense response to our standard format
func (s *typesenseService) formatSearchResponse(response *api.SearchResult, page, perPage int) *SearchResult {
	if response == nil {
		return &SearchResult{
			Found:     0,
			Page:      page,
			PerPage:   perPage,
			Documents: make([]map[string]interface{}, 0),
		}
	}

	found := 0
	if response.Found != nil {
		found = *response.Found
	}

	documents := s.extractDocuments(response)

	var facetCounts []api.FacetCounts
	if response.FacetCounts != nil {
		facetCounts = *response.FacetCounts
	}

	queryTime := float64(0)
	if response.SearchTimeMs != nil {
		queryTime = float64(*response.SearchTimeMs)
	}

	return &SearchResult{
		Found:      found,
		Page:       page,
		PerPage:    perPage,
		Documents:  documents,
		FacetCount: facetCounts,
		QueryTime:  queryTime,
	}
}

// extractDocuments safely extracts and enriches documents from search hits
func (s *typesenseService) extractDocuments(response *api.SearchResult) []map[string]interface{} {
	documents := make([]map[string]interface{}, 0)

	if response.Hits == nil {
		return documents
	}

	for _, hit := range *response.Hits {
		doc := s.extractDocument(hit)
		if doc != nil {
			documents = append(documents, doc)
		}
	}

	return documents
}

// extractDocument safely extracts a single document from a search hit
func (s *typesenseService) extractDocument(hit api.SearchResultHit) map[string]interface{} {
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
