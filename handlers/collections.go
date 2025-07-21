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

// CollectionHandler handles Typesense collection operations.
type CollectionHandler struct {
	typesenseService services.TypesenseService
}

// NewCollectionHandler creates a new CollectionHandler.
func NewCollectionHandler(typesenseService services.TypesenseService) *CollectionHandler {
	return &CollectionHandler{typesenseService: typesenseService}
}

// GetByName retrieves a collection by name.
func (h *CollectionHandler) GetByName(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionName, ok := req.Params.Arguments["collection_name"].(string)
	if !ok || collectionName == "" {
		return nil, fmt.Errorf("collection_name is required")
	}
	result, err := h.typesenseService.GetCollection(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection '%s': %w", collectionName, err)
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal collection: %v", err)
		return nil, fmt.Errorf("failed to format result: %w", err)
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// Create creates a new collection.
func (h *CollectionHandler) Create(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var schema api.CollectionSchema
	jsonData, err := json.Marshal(req.Params.Arguments["schema"])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}
	if err := json.Unmarshal(jsonData, &schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal collection schema: %w", err)
	}
	result, err := h.typesenseService.CreateCollection(ctx, &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal collection: %v", err)
		return nil, fmt.Errorf("failed to format result: %w", err)
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// UpdateByName updates a collection by name.
func (h *CollectionHandler) UpdateByName(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionName, ok := req.Params.Arguments["collection_name"].(string)
	if !ok || collectionName == "" {
		return nil, fmt.Errorf("collection_name is required")
	}
	var schema api.CollectionUpdateSchema
	jsonData, err := json.Marshal(req.Params.Arguments["schema"])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}
	if err := json.Unmarshal(jsonData, &schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal collection update schema: %w", err)
	}
	result, err := h.typesenseService.UpdateCollection(ctx, collectionName, &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to update collection '%s': %w", collectionName, err)
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal collection: %v", err)
		return nil, fmt.Errorf("failed to format result: %w", err)
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// DeleteByName deletes a collection by name.
func (h *CollectionHandler) DeleteByName(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collectionName, ok := req.Params.Arguments["collection_name"].(string)
	if !ok || collectionName == "" {
		return nil, fmt.Errorf("collection_name is required")
	}
	result, err := h.typesenseService.DeleteCollection(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to delete collection '%s': %w", collectionName, err)
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal collection: %v", err)
		return nil, fmt.Errorf("failed to format result: %w", err)
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}