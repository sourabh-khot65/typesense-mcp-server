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

type CollectionsHandler struct {
	typesenseService services.TypesenseService
}

func NewCollectionsHandler(typesenseService services.TypesenseService) *CollectionsHandler {
	return &CollectionsHandler{typesenseService: typesenseService}
}

func (h *CollectionsHandler) GetCollection(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("collection name is required")
	}
	result, err := h.typesenseService.GetCollection(ctx, name)
	if err != nil {
		return nil, err
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal collection: %v", err)
		return nil, fmt.Errorf("failed to format result: %v", err)
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (h *CollectionsHandler) CreateCollection(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var schema api.CollectionSchema
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}
	if err := json.Unmarshal(jsonData, &schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal collection schema: %v", err)
	}
	result, err := h.typesenseService.CreateCollection(ctx, &schema)
	if err != nil {
		return nil, err
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal collection: %v", err)
		return nil, fmt.Errorf("failed to format result: %v", err)
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (h *CollectionsHandler) UpdateCollection(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("collection name is required")
	}
	var schema api.CollectionUpdateSchema
	jsonData, err := json.Marshal(request.Params.Arguments["schema"])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %v", err)
	}
	if err := json.Unmarshal(jsonData, &schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal collection update schema: %v", err)
	}
	result, err := h.typesenseService.UpdateCollection(ctx, name, &schema)
	if err != nil {
		return nil, err
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal collection: %v", err)
		return nil, fmt.Errorf("failed to format result: %v", err)
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (h *CollectionsHandler) DeleteCollection(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("collection name is required")
	}
	result, err := h.typesenseService.DeleteCollection(ctx, name)
	if err != nil {
		return nil, err
	}
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal collection: %v", err)
		return nil, fmt.Errorf("failed to format result: %v", err)
	}
	return mcp.NewToolResultText(string(resultJSON)), nil
}