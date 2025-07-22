// Package constants defines application-wide constants.
package constants

import "time"

const (
	// Application metadata
	AppName    = "Typesense MCP Server"
	AppVersion = "1.0.0"
	
	// Default configuration values
	DefaultHost         = "localhost"
	DefaultPort         = 8108
	DefaultProtocol     = "http"
	DefaultAPIKey       = "xyz"
	DefaultPageSize     = 10
	DefaultPage         = 1
	
	// Limits and constraints
	MaxPageSize         = 100
	MinPageSize         = 1
	MaxPage             = 10000
	MinPort             = 1
	MaxPort             = 65535
	
	// Timeouts and retries
	DefaultTimeout      = 30 * time.Second
	DefaultRetryCount   = 3
	DefaultRetryDelay   = 1 * time.Second
	
	// Logging
	LogFieldComponent   = "component"
	LogFieldOperation   = "operation"
	LogFieldCollection  = "collection"
	LogFieldRequestID   = "request_id"
	LogFieldDuration    = "duration_ms"
	LogFieldError       = "error"
	LogFieldQuery       = "query"
	LogFieldPage        = "page"
	LogFieldPageSize    = "per_page"
	LogFieldResultCount = "result_count"
	
	// Environment variables
	EnvTypesenseHost     = "TYPESENSE_HOST"
	EnvTypesensePort     = "TYPESENSE_PORT"
	EnvTypesenseProtocol = "TYPESENSE_PROTOCOL"
	EnvTypesenseAPIKey   = "TYPESENSE_API_KEY"
	EnvLogLevel          = "LOG_LEVEL"
	EnvLogFormat         = "LOG_FORMAT"
	
	// MCP Tool names
	ToolTypesenseSearch      = "typesense_search"
	ToolTypesenseCollections = "typesense_collections"
	
	// Parameter names
	ParamCollection = "collection"
	ParamQuery      = "q"
	ParamQueryBy    = "query_by"
	ParamFilterBy   = "filter_by"
	ParamPage       = "page"
	ParamPerPage    = "per_page"
	
	// Default values for search parameters
	DefaultQueryBy = "*"
	
	// Component names for logging
	ComponentMain     = "main"
	ComponentConfig   = "config"
	ComponentService  = "service"
	ComponentHandler  = "handler"
	ComponentHealth   = "health"
	ComponentMetrics  = "metrics"
)

// Supported protocols
var SupportedProtocols = []string{"http", "https"}

// Supported log levels
var SupportedLogLevels = []string{"debug", "info", "warn", "error", "fatal", "panic"}

// Supported log formats
var SupportedLogFormats = []string{"json", "text"}
