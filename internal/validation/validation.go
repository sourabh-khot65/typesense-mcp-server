// Package validation provides input validation functionality.
package validation

import (
	"context"
	"strings"

	"typesense-mcp-server/internal/constants"
	"typesense-mcp-server/internal/errors"

	"github.com/mark3labs/mcp-go/mcp"
)

// Validator provides validation methods for different types of input
type Validator struct{}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{}
}

// ValidateSearchRequest validates a search request and extracts parameters
func (v *Validator) ValidateSearchRequest(ctx context.Context, request mcp.CallToolRequest) (*SearchParams, error) {
	params := &SearchParams{}

	// Validate and extract collection
	collection, err := v.extractAndValidateCollection(request)
	if err != nil {
		return nil, err
	}
	params.Collection = collection

	// Validate and extract query
	query, err := v.extractAndValidateQuery(request)
	if err != nil {
		return nil, err
	}
	params.Query = query

	// Extract optional parameters with defaults
	params.QueryBy = v.extractStringWithDefault(request, constants.ParamQueryBy, constants.DefaultQueryBy)
	params.FilterBy = v.extractStringWithDefault(request, constants.ParamFilterBy, "")

	// Validate and extract pagination
	page, perPage, err := v.extractAndValidatePagination(request)
	if err != nil {
		return nil, err
	}
	params.Page = page
	params.PerPage = perPage

	return params, nil
}

// SearchParams holds validated search parameters
type SearchParams struct {
	Collection string
	Query      string
	QueryBy    string
	FilterBy   string
	Page       int
	PerPage    int
}

// ToMap converts SearchParams to a map for API calls
func (p *SearchParams) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		constants.ParamQuery:   p.Query,
		constants.ParamQueryBy: p.QueryBy,
		constants.ParamPage:    p.Page,
		constants.ParamPerPage: p.PerPage,
	}

	if p.FilterBy != "" {
		result[constants.ParamFilterBy] = p.FilterBy
	}

	return result
}

// extractAndValidateCollection extracts and validates the collection parameter
func (v *Validator) extractAndValidateCollection(request mcp.CallToolRequest) (string, error) {
	collection, ok := request.Params.Arguments[constants.ParamCollection].(string)
	if !ok {
		return "", errors.NewInvalidCollectionError("collection parameter must be a string")
	}

	collection = strings.TrimSpace(collection)
	if collection == "" {
		return "", errors.NewInvalidCollectionError("collection name cannot be empty")
	}

	// Validate collection name format (basic alphanumeric and underscores)
	if !isValidCollectionName(collection) {
		return "", errors.NewInvalidCollectionError("collection name contains invalid characters")
	}

	return collection, nil
}

// extractAndValidateQuery extracts and validates the query parameter
func (v *Validator) extractAndValidateQuery(request mcp.CallToolRequest) (string, error) {
	query, ok := request.Params.Arguments[constants.ParamQuery].(string)
	if !ok {
		return "", errors.NewInvalidQueryError("query parameter must be a string")
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return "", errors.NewInvalidQueryError("query cannot be empty")
	}

	return query, nil
}

// extractAndValidatePagination extracts and validates pagination parameters
func (v *Validator) extractAndValidatePagination(request mcp.CallToolRequest) (int, int, error) {
	// Extract page
	page := constants.DefaultPage
	if pageValue, exists := request.Params.Arguments[constants.ParamPage]; exists {
		if pageFloat, ok := pageValue.(float64); ok {
			page = int(pageFloat)
		} else {
			return 0, 0, errors.NewInvalidPaginationError(0, 0).
				WithDetail("error", "page parameter must be a number")
		}
	}

	// Extract per_page
	perPage := constants.DefaultPageSize
	if perPageValue, exists := request.Params.Arguments[constants.ParamPerPage]; exists {
		if perPageFloat, ok := perPageValue.(float64); ok {
			perPage = int(perPageFloat)
		} else {
			return 0, 0, errors.NewInvalidPaginationError(page, 0).
				WithDetail("error", "per_page parameter must be a number")
		}
	}

	// Validate pagination values
	if err := v.validatePaginationValues(page, perPage); err != nil {
		return 0, 0, err
	}

	return page, perPage, nil
}

// validatePaginationValues validates the pagination parameter values
func (v *Validator) validatePaginationValues(page, perPage int) error {
	if page < 1 || page > constants.MaxPage {
		return errors.NewInvalidPaginationError(page, perPage).
			WithDetail("error", "page must be between 1 and 10000")
	}

	if perPage < constants.MinPageSize || perPage > constants.MaxPageSize {
		return errors.NewInvalidPaginationError(page, perPage).
			WithDetail("error", "per_page must be between 1 and 100")
	}

	return nil
}

// extractStringWithDefault extracts a string parameter with a default value
func (v *Validator) extractStringWithDefault(request mcp.CallToolRequest, key, defaultValue string) string {
	if value, exists := request.Params.Arguments[key]; exists {
		if strValue, ok := value.(string); ok {
			strValue = strings.TrimSpace(strValue)
			if strValue != "" {
				return strValue
			}
		}
	}
	return defaultValue
}

// isValidCollectionName checks if a collection name is valid
func isValidCollectionName(name string) bool {
	if len(name) == 0 || len(name) > 100 {
		return false
	}

	// Allow alphanumeric characters, underscores, and hyphens
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}

	return true
}

// ValidateCollectionsRequest validates a collections request
func (v *Validator) ValidateCollectionsRequest(ctx context.Context, request mcp.CallToolRequest) error {
	// Collections request doesn't require any specific parameters
	// This is a placeholder for future validation if needed
	return nil
}
