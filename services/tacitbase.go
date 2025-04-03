package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"tb-mcp-server/models"
)

// TacitbaseService defines the interface for Tacitbase operations
type TacitbaseService interface {
	Search(ctx context.Context, request *models.SearchRequest) (*models.SearchResponse, error)
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
func (s *tacitbaseService) Search(ctx context.Context, request *models.SearchRequest) (*models.SearchResponse, error) {
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
	var searchResp models.SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &searchResp, nil
}

// FormatSearchResults formats the search results into a human-readable string
func FormatSearchResults(resp *models.SearchResponse) string {
	if resp == nil || len(resp.Hits) == 0 {
		return "No results found."
	}

	var results []string
	for _, hit := range resp.Hits {
		// Extract fields from the hit map
		firstName := getStringValue(hit, "first_name")
		lastName := getStringValue(hit, "last_name")
		email := getStringValue(hit, "email")
		location := getStringValue(hit, "location")
		skills := getStringSlice(hit, "skills")
		latestExperience := getStringValue(hit, "latest_experience")
		highestEducation := getStringValue(hit, "highest_education")

		result := fmt.Sprintf("Name: %s %s\nEmail: %s\nLocation: %s\nSkills: %s\nLatest Experience: %s\nHighest Education: %s\n---",
			firstName,
			lastName,
			email,
			location,
			strings.Join(skills, ", "),
			latestExperience,
			highestEducation,
		)
		results = append(results, result)
	}

	return fmt.Sprintf("Found %d candidates (Page %d, %d per page):\n\n%s",
		resp.Found,
		resp.Page,
		resp.PerPage,
		strings.Join(results, "\n"),
	)
}

// Helper function to safely get string value from map
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// Helper function to safely get string slice from map
func getStringSlice(m map[string]interface{}, key string) []string {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case []string:
			return v
		case []interface{}:
			result := make([]string, 0, len(v))
			for _, item := range v {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		case string:
			return strings.Split(v, ",")
		}
	}
	return []string{}
}
