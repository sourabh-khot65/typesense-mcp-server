package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"tb-mcp-server/models"
	"tb-mcp-server/services"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"github.com/typesense/typesense-go/typesense/api"
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
func (h *SearchHandler) HandleCandidateSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var searchReq api.SearchCollectionParams
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		logrus.Errorf("ERROR: Failed to marshal search arguments: %v", err)
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		logrus.Errorf("ERROR: Failed to unmarshal search request: %v", err)
		return nil, fmt.Errorf("failed to unmarshal search request: %v", err)
	}

	// Perform search using Typesense
	response, err := h.typesenseService.Search(ctx, "candidates_candidates", &searchReq)
	if err != nil {
		logrus.Errorf("ERROR: Search operation failed: %v", err)
		return nil, fmt.Errorf("search failed: %v", err)
	}

	// Log the response for debugging
	logrus.Infof("DEBUG: Search response: %+v", response)

	return mcp.NewToolResultText(FormatTypesenseCandidateResults(response)), nil
}

// HandleAttachmentsSearch handles attachment search requests
func (h *SearchHandler) HandleAttachmentsSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var searchReq api.SearchCollectionParams
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		logrus.Errorf("ERROR: Failed to marshal attachment search arguments: %v", err)
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		logrus.Errorf("ERROR: Failed to unmarshal attachments search request: %v", err)
		return nil, fmt.Errorf("failed to unmarshal attachments search request: %v", err)
	}

	// Perform search using Typesense
	response, err := h.typesenseService.Search(ctx, "candidates_candidate-attachments", &searchReq)
	if err != nil {
		logrus.Errorf("ERROR: Attachments search operation failed: %v", err)
		return nil, fmt.Errorf("attachments search failed: %v", err)
	}

	// Debug: Print raw response
	if response != nil && response.Hits != nil && len(*response.Hits) > 0 {
		logrus.Infof("DEBUG: First attachment search hit: %+v", (*response.Hits)[0])
	}

	formattedResults := FormatTypesenseAttachmentResults(response)
	return mcp.NewToolResultText(formattedResults), nil
}

// HandleStagingSearch handles search requests specifically for staging environment
func (h *SearchHandler) HandleStagingSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var searchReq api.SearchCollectionParams
	jsonData, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		logrus.Errorf("ERROR: Failed to marshal staging search arguments: %v", err)
		return nil, fmt.Errorf("failed to marshal arguments: %v", err)
	}

	if err := json.Unmarshal(jsonData, &searchReq); err != nil {
		logrus.Errorf("ERROR: Failed to unmarshal staging search request: %v", err)
		return nil, fmt.Errorf("failed to unmarshal search request: %v", err)
	}

	// Perform search using Typesense
	response, err := h.tacitbaseService.SearchCandidates(ctx, "candidates_candidates", &searchReq)
	if err != nil {
		logrus.Errorf("ERROR: Staging search operation failed: %v", err)
		return nil, fmt.Errorf("staging search failed: %v", err)
	}

	logrus.Infof("DEBUG: Staging search response: %+v", response)

	return mcp.NewToolResultText(FormatCandidateResults(response)), nil
}

// FormatCandidateResults formats the candidate search response
func FormatCandidateResults(response *models.CandidateSearchResponse) string {
	formattedResults := fmt.Sprintf("Found %d candidates\n", len(response.Candidates.Items))
	for _, candidate := range response.Candidates.Items {
		formattedResults += fmt.Sprintf("Candidate: %s %s\n", candidate.FirstName, candidate.LastName)
		formattedResults += fmt.Sprintf("Email: %s\n", candidate.Email)
		formattedResults += fmt.Sprintf("Phone: %s\n", candidate.Phone)
		formattedResults += fmt.Sprintf("Location: %s\n", candidate.Location)
		formattedResults += fmt.Sprintf("Highest Education: %s\n", candidate.HighestEducation)
		formattedResults += fmt.Sprintf("Latest Experience: %s\n", candidate.LatestExperience)
		formattedResults += fmt.Sprintf("Description: %s\n", candidate.Description)
		formattedResults += fmt.Sprintf("LinkedIn: %s\n", candidate.LinkedIn)
		formattedResults += fmt.Sprintf("GitHub: %s\n", candidate.GitHub)
	}
	return formattedResults
}

// FormatTypesenseCandidateResults formats the Typesense response
func FormatTypesenseCandidateResults(response *api.SearchResult) string {
	formattedResults := fmt.Sprintf("Found %d candidates\n", *response.Found)
	for _, hit := range *response.Hits {
		formattedResults += fmt.Sprintf("Candidate: %s\n", hit.Document)
	}
	return formattedResults
}

// FormatTypesenseAttachmentResults formats the Typesense response
func FormatTypesenseAttachmentResults(response *api.SearchResult) string {
	formattedResults := fmt.Sprintf("Found %d attachments\n", *response.Found)
	for _, hit := range *response.Hits {
		jsonData, err := json.Marshal(hit.Document)
		if err != nil {
			logrus.Errorf("ERROR: Failed to marshal attachment result: %v", err)
			continue
		}
		var attachment models.Attachment
		if err := json.Unmarshal(jsonData, &attachment); err != nil {
			logrus.Errorf("ERROR: Failed to unmarshal attachment result: %v", err)
			continue
		}
		formattedResults += fmt.Sprintf("Attachment: %s\n", attachment.Name)
	}
	return formattedResults
}
