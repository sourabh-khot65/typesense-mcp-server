package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"tb-mcp-server/models"

	"github.com/typesense/typesense-go/typesense/api"
)

// TacitbaseService defines the interface for Tacitbase operations
type TacitbaseService interface {
	SearchCandidates(ctx context.Context, collection string, request *api.SearchCollectionParams) (*models.CandidateSearchResponse, error)
}

// tacitbaseService implements the TacitbaseService interface
type tacitbaseService struct {
	baseURL    string
	httpClient *http.Client
}

// NewTacitbaseService creates a new instance of TacitbaseService
func NewTacitbaseService() TacitbaseService {
	return &tacitbaseService{
		baseURL:    "https://staging.local.tacitbase.com/v1",
		httpClient: &http.Client{},
	}
}

// Search performs a search operation using Tacitbase's API
func (s *tacitbaseService) SearchCandidates(ctx context.Context, collection string, request *api.SearchCollectionParams) (*models.CandidateSearchResponse, error) {
	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/search/documents/find/candidates", s.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Get authorization token from environment
	if authToken := os.Getenv("TACITBASE_AUTH_TOKEN"); authToken != "" {
		httpReq.Header.Set("Authorization", authToken)
	} else {
		return nil, fmt.Errorf("TACITBASE_AUTH_TOKEN environment variable is not set")
	}

	// Send request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var searchResp models.CandidateSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &searchResp, nil
}
