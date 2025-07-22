// Package errors defines custom error types for the application.
package errors

import (
	"fmt"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "validation"
	ErrorTypeConfiguration ErrorType = "configuration"
	ErrorTypeService       ErrorType = "service"
	ErrorTypeClient        ErrorType = "client"
	ErrorTypeInternal      ErrorType = "internal"
	ErrorTypeNotFound      ErrorType = "not_found"
	ErrorTypeTimeout       ErrorType = "timeout"
	ErrorTypeUnauthorized  ErrorType = "unauthorized"
)

// Error represents a structured error with additional context
type Error struct {
	Type        ErrorType              `json:"type"`
	Message     string                 `json:"message"`
	Code        string                 `json:"code,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Cause       error                  `json:"-"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target error
func (e *Error) Is(target error) bool {
	if targetErr, ok := target.(*Error); ok {
		return e.Type == targetErr.Type && e.Code == targetErr.Code
	}
	return false
}

// WithDetail adds a detail field to the error
func (e *Error) WithDetail(key string, value interface{}) *Error {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// NewError creates a new Error with the specified type and message
func NewError(errorType ErrorType, message string) *Error {
	return &Error{
		Type:    errorType,
		Message: message,
	}
}

// NewErrorWithCode creates a new Error with type, code, and message
func NewErrorWithCode(errorType ErrorType, code, message string) *Error {
	return &Error{
		Type:    errorType,
		Code:    code,
		Message: message,
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, errorType ErrorType, message string) *Error {
	return &Error{
		Type:    errorType,
		Message: message,
		Cause:   err,
	}
}

// Predefined error codes
const (
	CodeInvalidCollection     = "INVALID_COLLECTION"
	CodeInvalidQuery          = "INVALID_QUERY"
	CodeInvalidPagination     = "INVALID_PAGINATION"
	CodeInvalidConfiguration  = "INVALID_CONFIGURATION"
	CodeClientNotInitialized  = "CLIENT_NOT_INITIALIZED"
	CodeServiceUnavailable    = "SERVICE_UNAVAILABLE"
	CodeTimeout               = "TIMEOUT"
	CodeUnauthorized          = "UNAUTHORIZED"
	CodeInternalError         = "INTERNAL_ERROR"
)

// Validation errors
func NewValidationError(message string) *Error {
	return NewError(ErrorTypeValidation, message)
}

func NewInvalidCollectionError(collection string) *Error {
	return NewErrorWithCode(ErrorTypeValidation, CodeInvalidCollection, "invalid collection").
		WithDetail("collection", collection)
}

func NewInvalidQueryError(query string) *Error {
	return NewErrorWithCode(ErrorTypeValidation, CodeInvalidQuery, "invalid query").
		WithDetail("query", query)
}

func NewInvalidPaginationError(page, perPage int) *Error {
	return NewErrorWithCode(ErrorTypeValidation, CodeInvalidPagination, "invalid pagination parameters").
		WithDetail("page", page).
		WithDetail("per_page", perPage)
}

// Configuration errors
func NewConfigurationError(message string) *Error {
	return NewError(ErrorTypeConfiguration, message)
}

func NewInvalidConfigurationError(field string, value interface{}) *Error {
	return NewErrorWithCode(ErrorTypeConfiguration, CodeInvalidConfiguration, "invalid configuration").
		WithDetail("field", field).
		WithDetail("value", value)
}

// Service errors
func NewServiceError(message string) *Error {
	return NewError(ErrorTypeService, message)
}

func NewClientNotInitializedError() *Error {
	return NewErrorWithCode(ErrorTypeService, CodeClientNotInitialized, "client is not initialized")
}

func NewServiceUnavailableError(service string) *Error {
	return NewErrorWithCode(ErrorTypeService, CodeServiceUnavailable, "service unavailable").
		WithDetail("service", service)
}

// Client errors
func NewClientError(message string) *Error {
	return NewError(ErrorTypeClient, message)
}

func NewTimeoutError(operation string) *Error {
	return NewErrorWithCode(ErrorTypeTimeout, CodeTimeout, "operation timed out").
		WithDetail("operation", operation)
}

func NewUnauthorizedError() *Error {
	return NewErrorWithCode(ErrorTypeUnauthorized, CodeUnauthorized, "unauthorized access")
}

// Internal errors
func NewInternalError(message string) *Error {
	return NewError(ErrorTypeInternal, message)
}

func NewInternalErrorWithCause(message string, cause error) *Error {
	return WrapError(cause, ErrorTypeInternal, message)
}
