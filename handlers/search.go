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

	// Set default collection if not specified
	if searchReq.Collection == "" {
		searchReq.Collection = "candidates_candidates"
	}

	// Perform search using Typesense
	response, err := h.typesenseService.Search(ctx, &searchReq)
	if err != nil {
		// Fallback to Tacitbase if Typesense is unavailable
		response, err = h.tacitbaseService.Search(ctx, &searchReq)
		if err != nil {
			return nil, fmt.Errorf("search failed: %v", err)
		}
	}

	return mcp.NewToolResultText(FormatSearchResults(response)), nil
}

// HandleVectorSearch handles vector search requests
func (h *SearchHandler) HandleVectorSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var searchReq models.VectorSearchRequest
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vector search request: %v", err)
	}

	// Perform vector search using Typesense
	response, err := h.typesenseService.VectorSearch(ctx, &searchReq)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %v", err)
	}

	return mcp.NewToolResultText(FormatSearchResults(response)), nil
}

// HandleSemanticSearch handles semantic search requests
func (h *SearchHandler) HandleSemanticSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var searchReq models.SemanticSearchRequest
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal semantic search request: %v", err)
	}

	// Perform semantic search using Typesense
	response, err := h.typesenseService.SemanticSearch(ctx, &searchReq)
	if err != nil {
		return nil, fmt.Errorf("semantic search failed: %v", err)
	}

	return mcp.NewToolResultText(FormatSearchResults(response)), nil
}

type AttachmentsSearchRequest struct {
	Query        string `json:"query"`
	SearchFields string `json:"search_fields,omitempty"`
	FilterFields string `json:"filter_fields,omitempty"`
	SortBy       string `json:"sort_by,omitempty"`
	RecordID     string `json:"record_id,omitempty"`
	Page         int    `json:"page,omitempty"`
	PerPage      int    `json:"per_page,omitempty"`
}

func (h *SearchHandler) HandleAttachmentsSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var attachmentsReq AttachmentsSearchRequest
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &attachmentsReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attachments search request: %v", err)
	}

	// Create search request
	searchReq := models.SearchRequest{
		Collection: "candidates_candidate-attachments",
		Query:      attachmentsReq.Query,
		Page:       attachmentsReq.Page,
		PerPage:    attachmentsReq.PerPage,
	}

	// Convert comma-separated strings to slices
	if attachmentsReq.SearchFields != "" {
		searchReq.SearchFields = strings.Split(attachmentsReq.SearchFields, ",")
	}
	if attachmentsReq.FilterFields != "" {
		searchReq.FilterFields = strings.Split(attachmentsReq.FilterFields, ",")
	}
	if attachmentsReq.SortBy != "" {
		searchReq.SortBy = strings.Split(attachmentsReq.SortBy, ",")
	}

	// Add record_id filter if provided
	if attachmentsReq.RecordID != "" {
		if searchReq.FilterFields == nil {
			searchReq.FilterFields = make([]string, 0)
		}
		searchReq.FilterFields = append(searchReq.FilterFields, fmt.Sprintf("record_id:%s", attachmentsReq.RecordID))
	}

	// Set default search fields if not specified
	if len(searchReq.SearchFields) == 0 {
		searchReq.SearchFields = []string{"name", "content"}
	}

	// Perform search using Typesense
	response, err := h.typesenseService.Search(ctx, &searchReq)
	if err != nil {
		// Fallback to Tacitbase if Typesense is unavailable
		response, err = h.tacitbaseService.Search(ctx, &searchReq)
		if err != nil {
			return nil, fmt.Errorf("attachments search failed: %v", err)
		}
	}

	return mcp.NewToolResultText(FormatSearchResults(response)), nil
}

func FormatSearchResults(response *models.SearchResponse) string {
	if response == nil {
		return "No results found"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d results (Page %d)\n\n", response.Found, response.Page))

	for i, hit := range response.Hits {
		var candidate models.Candidate
		candidateBytes, err := json.Marshal(hit)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(candidateBytes, &candidate); err != nil {
			// Try to unmarshal as attachment if candidate unmarshal fails
			var attachment models.Attachment
			if err := json.Unmarshal(candidateBytes, &attachment); err != nil {
				continue
			}
			sb.WriteString(fmt.Sprintf("%d. [Attachment] %s\n", i+1, attachment.Name))
			if attachment.Content != "" {
				sb.WriteString(fmt.Sprintf("   Content: %s\n", truncateString(attachment.Content, 100)))
			}
			sb.WriteString(fmt.Sprintf("   Record ID: %s\n", attachment.RecordID))
		} else {
			// Format candidate information
			sb.WriteString(fmt.Sprintf("%d. %s %s\n", i+1, candidate.FirstName, candidate.LastName))
			if candidate.Email != "" {
				sb.WriteString(fmt.Sprintf("   📧 Email: %s\n", candidate.Email))
			}
			if candidate.Phone != "" {
				sb.WriteString(fmt.Sprintf("   📱 Phone: %s\n", candidate.Phone))
			}
			if candidate.Skills != "" {
				sb.WriteString(fmt.Sprintf("   💪 Skills: %s\n", candidate.Skills))
			}
			if candidate.LatestExperience != "" {
				sb.WriteString(fmt.Sprintf("   💼 Latest Experience: %s\n", candidate.LatestExperience))
			}
			if candidate.HighestEducation != "" {
				sb.WriteString(fmt.Sprintf("   🎓 Education: %s\n", candidate.HighestEducation))
			}
			if candidate.LinkedIn != "" {
				sb.WriteString(fmt.Sprintf("   🔗 LinkedIn: %s\n", candidate.LinkedIn))
			}
			if candidate.GitHub != "" {
				sb.WriteString(fmt.Sprintf("   💻 GitHub: %s\n", candidate.GitHub))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
