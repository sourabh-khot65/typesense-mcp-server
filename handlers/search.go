package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"tb-mcp-server/models"
	"tb-mcp-server/services"

	"github.com/mark3labs/mcp-go/mcp"
)

// SearchHandler handles candidate search requests
type SearchHandler struct {
	tacitbaseService services.TacitbaseService
	typesenseService services.TypesenseService
}

// NewSearchHandler creates a new instance of SearchHandler
func NewSearchHandler(tacitbaseService services.TacitbaseService, typesenseService services.TypesenseService) *SearchHandler {
	return &SearchHandler{
		tacitbaseService: tacitbaseService,
		typesenseService: typesenseService,
	}
}

// HandleSearch handles the basic search request
func (h *SearchHandler) HandleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var searchReq models.SearchRequest
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search request: %v", err)
	}

	// Set default collection
	searchReq.Collection = "candidates_candidates"

	// Set default per_page if not specified
	if searchReq.PerPage == 0 {
		searchReq.PerPage = 10
	} else if searchReq.PerPage > 100 {
		searchReq.PerPage = 100
	}

	// Ensure page is at least 1
	if searchReq.Page < 1 {
		searchReq.Page = 1
	}

	// Perform search using Typesense
	response, err := h.typesenseService.Search(ctx, &searchReq)
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}

	return mcp.NewToolResultText(FormatSearchResults(response)), nil
}

// HandleAttachmentsSearch handles attachment search requests
func (h *SearchHandler) HandleAttachmentsSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var searchReq models.SearchRequest
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attachments search request: %v", err)
	}

	// Set collection to attachments
	searchReq.Collection = "candidates_candidate-attachments"

	// Perform search using Typesense
	response, err := h.typesenseService.Search(ctx, &searchReq)
	if err != nil {
		return nil, fmt.Errorf("attachments search failed: %v", err)
	}
	// Debug: Print raw response
	if response != nil && len(response.Hits) > 0 {
		fmt.Printf("First hit: %+v\n", response.Hits[0])
	}

	// Convert hits to attachments
	attachments := make([]models.Attachment, 0, len(response.Hits))
	for _, hit := range response.Hits {
		var attachment models.Attachment
		hitBytes, err := json.Marshal(hit)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(hitBytes, &attachment); err != nil {
			continue
		}
		attachments = append(attachments, attachment)
	}

	formattedResults := FormatSearchResults(attachments)
	return mcp.NewToolResultText(formattedResults), nil
}

// HandleStagingSearch handles search requests specifically for staging environment
func (h *SearchHandler) HandleStagingSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var searchReq models.SearchRequest
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search request: %v", err)
	}

	// Set default collection
	searchReq.Collection = "candidates_candidates"

	// Perform search directly using Tacitbase staging service
	response, err := h.tacitbaseService.Search(ctx, &searchReq)
	if err != nil {
		return nil, fmt.Errorf("staging search failed: %v", err)
	}

	return mcp.NewToolResultText(FormatSearchResults(response)), nil
}

// HandleAttachmentsSearchTool handles attachment search requests as a standalone tool
func (h *SearchHandler) HandleAttachmentsSearchTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var searchReq models.SearchRequest
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attachments search request: %v", err)
	}

	searchReq.Collection = "candidates_candidate-attachments"

	response, err := h.typesenseService.Search(ctx, &searchReq)
	if err != nil {
		return nil, fmt.Errorf("attachments search failed: %v", err)
	}

	attachments := make([]models.Attachment, 0, len(response.Hits))
	for _, hit := range response.Hits {
		var attachment models.Attachment
		hitBytes, err := json.Marshal(hit)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(hitBytes, &attachment); err != nil {
			continue
		}
		attachments = append(attachments, attachment)
	}

	formattedResults := FormatSearchResults(attachments)
	return mcp.NewToolResultText(formattedResults), nil
}

func FormatSearchResults(results interface{}) string {
	var formattedResults strings.Builder

	switch v := results.(type) {
	case []models.Attachment:
		formattedResults.WriteString(fmt.Sprintf("Found %d attachments:\n\n", len(v)))
		for i, attachment := range v {
			formattedResults.WriteString(fmt.Sprintf("%d. ", i+1))
			if attachment.Name != "" {
				formattedResults.WriteString(fmt.Sprintf("Name: %s\n", attachment.Name))
			}
			if attachment.RecordID != "" {
				formattedResults.WriteString(fmt.Sprintf("   Record ID: %s\n", attachment.RecordID))
			}
			if attachment.ModelName != "" {
				formattedResults.WriteString(fmt.Sprintf("   Model: %s\n", attachment.ModelName))
			}
			if attachment.Content != "" {
				formattedResults.WriteString(fmt.Sprintf("   Content Preview: %s\n", attachment.Content))
			}
			if attachment.CreatedAt != "" {
				formattedResults.WriteString(fmt.Sprintf("   Created: %s\n", attachment.CreatedAt))
			}
			formattedResults.WriteString("\n")
		}

	case []models.Candidate:
		formattedResults.WriteString(fmt.Sprintf("Found %d candidates:\n\n", len(v)))
		for i, candidate := range v {
			formattedResults.WriteString(fmt.Sprintf("%d. ", i+1))
			if candidate.FirstName != "" || candidate.LastName != "" {
				formattedResults.WriteString(fmt.Sprintf("Name: %s %s\n", candidate.FirstName, candidate.LastName))
			}
			if candidate.Email != "" {
				formattedResults.WriteString(fmt.Sprintf("   Email: %s\n", candidate.Email))
			}
			if len(candidate.Skills) > 0 {
				formattedResults.WriteString(fmt.Sprintf("   Skills: %s\n", strings.Join(candidate.Skills, ", ")))
			}
			if candidate.LatestExperience != "" {
				formattedResults.WriteString(fmt.Sprintf("   Latest Experience: %s\n", candidate.LatestExperience))
			}
			if candidate.HighestEducation != "" {
				formattedResults.WriteString(fmt.Sprintf("   Highest Education: %s\n", candidate.HighestEducation))
			}
			if candidate.Description != "" {
				formattedResults.WriteString(fmt.Sprintf("   Description: %s\n", candidate.Description))
			}
			formattedResults.WriteString("\n")
		}
	}

	return formattedResults.String()
}
