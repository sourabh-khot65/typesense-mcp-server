// Package handlers provides HTTP and MCP request handlers.
package handlers

import (
	"context"
	"encoding/json"
	"time"

	"typesense-mcp-server/internal/constants"
	"typesense-mcp-server/internal/errors"
	"typesense-mcp-server/internal/logger"
	"typesense-mcp-server/internal/services"
	"typesense-mcp-server/internal/validation"

	"github.com/mark3labs/mcp-go/mcp"
)

// SearchHandler handles search-related MCP requests
type SearchHandler struct {
	typesenseService services.TypesenseService
	validator        *validation.Validator
	logger           *logger.Logger
}

// NewSearchHandler creates a new SearchHandler instance
func NewSearchHandler(typesenseService services.TypesenseService, log *logger.Logger) *SearchHandler {
	if log == nil {
		log = logger.NewDefault()
	}

	return &SearchHandler{
		typesenseService: typesenseService,
		validator:        validation.New(),
		logger:           log,
	}
}

// SearchInTypesenseCollection handles search requests for Typesense collections
func (h *SearchHandler) SearchInTypesenseCollection(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()
	
	entry := h.logger.WithComponent(constants.ComponentHandler).
		WithField(constants.LogFieldOperation, "search_collection")

	// Validate request and extract parameters
	params, err := h.validator.ValidateSearchRequest(ctx, request)
	if err != nil {
		duration := time.Since(start)
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).WithError(err).Error("search request validation failed")
		return nil, h.handleValidationError(err)
	}

	entry = entry.WithField(constants.LogFieldCollection, params.Collection).
		WithField(constants.LogFieldQuery, params.Query).
		WithField(constants.LogFieldPage, params.Page).
		WithField(constants.LogFieldPageSize, params.PerPage)

	// Perform search
	result, err := h.typesenseService.Search(ctx, params)
	if err != nil {
		duration := time.Since(start)
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).WithError(err).Error("search operation failed")
		return nil, h.handleServiceError(err)
	}

	// Format response
	response, err := h.formatSearchResponse(result)
	if err != nil {
		duration := time.Since(start)
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).WithError(err).Error("search response formatting failed")
		return nil, h.handleInternalError(err)
	}

	duration := time.Since(start)
	entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).
		WithField(constants.LogFieldResultCount, result.Found).
		Info("search request completed successfully")

	return response, nil
}

// GetTypesenseCollections handles requests to retrieve all Typesense collections
func (h *SearchHandler) GetTypesenseCollections(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	start := time.Now()
	
	entry := h.logger.WithComponent(constants.ComponentHandler).
		WithField(constants.LogFieldOperation, "get_collections")

	// Validate request (collections request has no specific parameters to validate)
	if err := h.validator.ValidateCollectionsRequest(ctx, request); err != nil {
		duration := time.Since(start)
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).WithError(err).Error("collections request validation failed")
		return nil, h.handleValidationError(err)
	}

	// Get collections
	result, err := h.typesenseService.GetCollections(ctx)
	if err != nil {
		duration := time.Since(start)
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).WithError(err).Error("get collections operation failed")
		return nil, h.handleServiceError(err)
	}

	// Format response
	response, err := h.formatCollectionsResponse(result)
	if err != nil {
		duration := time.Since(start)
		entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).WithError(err).Error("collections response formatting failed")
		return nil, h.handleInternalError(err)
	}

	duration := time.Since(start)
	entry.WithField(constants.LogFieldDuration, duration.Milliseconds()).
		WithField("collection_count", result.Count).
		Info("collections request completed successfully")

	return response, nil
}

// formatSearchResponse formats a search result for MCP response
func (h *SearchHandler) formatSearchResponse(result *services.SearchResult) (*mcp.CallToolResult, error) {
	if result == nil {
		return nil, errors.NewInternalError("search result is nil")
	}

	// Convert to JSON with proper indentation
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrorTypeInternal, "failed to marshal search result")
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

// formatCollectionsResponse formats a collections result for MCP response
func (h *SearchHandler) formatCollectionsResponse(result *services.CollectionsResult) (*mcp.CallToolResult, error) {
	if result == nil {
		return nil, errors.NewInternalError("collections result is nil")
	}

	// Convert to JSON with proper indentation
	resultJSON, err := json.MarshalIndent(result.Collections, "", "  ")
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrorTypeInternal, "failed to marshal collections result")
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleValidationError handles validation errors with appropriate logging
func (h *SearchHandler) handleValidationError(err error) error {
	h.logger.WithComponent(constants.ComponentHandler).
		WithError(err).
		Error("validation error occurred")
	
	// Return the validation error as-is since it contains appropriate user-facing messages
	return err
}

// handleServiceError handles service errors with appropriate logging and transformation
func (h *SearchHandler) handleServiceError(err error) error {
	h.logger.WithComponent(constants.ComponentHandler).
		WithError(err).
		Error("service error occurred")
	
	// Check if it's already a structured error
	if structuredErr, ok := err.(*errors.Error); ok {
		return structuredErr
	}
	
	// Wrap unknown service errors
	return errors.WrapError(err, errors.ErrorTypeService, "service operation failed")
}

// handleInternalError handles internal errors with appropriate logging
func (h *SearchHandler) handleInternalError(err error) error {
	h.logger.WithComponent(constants.ComponentHandler).
		WithError(err).
		Error("internal error occurred")
	
	// Check if it's already a structured error
	if structuredErr, ok := err.(*errors.Error); ok {
		return structuredErr
	}
	
	// Wrap unknown internal errors
	return errors.WrapError(err, errors.ErrorTypeInternal, "internal operation failed")
}
