// Package config provides configuration management functionality.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"typesense-mcp-server/internal/constants"
	"typesense-mcp-server/internal/errors"
)

// Config holds all application configuration
type Config struct {
	Typesense TypesenseConfig `json:"typesense"`
	Server    ServerConfig    `json:"server"`
	Logging   LoggingConfig   `json:"logging"`
}

// TypesenseConfig holds Typesense-specific configuration
type TypesenseConfig struct {
	Host     string        `json:"host"`
	Port     int           `json:"port"`
	Protocol string        `json:"protocol"`
	APIKey   string        `json:"api_key"`
	Timeout  time.Duration `json:"timeout"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// LoggingConfig holds logging-specific configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// New creates a new Config instance with values from environment variables
func New() (*Config, error) {
	config := &Config{
		Typesense: TypesenseConfig{
			Host:     getEnvOrDefault(constants.EnvTypesenseHost, constants.DefaultHost),
			Port:     getEnvIntOrDefault(constants.EnvTypesensePort, constants.DefaultPort),
			Protocol: getEnvOrDefault(constants.EnvTypesenseProtocol, constants.DefaultProtocol),
			APIKey:   getEnvOrDefault(constants.EnvTypesenseAPIKey, constants.DefaultAPIKey),
			Timeout:  constants.DefaultTimeout,
		},
		Server: ServerConfig{
			Name:    constants.AppName,
			Version: constants.AppVersion,
		},
		Logging: LoggingConfig{
			Level:  getEnvOrDefault(constants.EnvLogLevel, "info"),
			Format: getEnvOrDefault(constants.EnvLogFormat, "text"),
			Output: "stdout",
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates the entire configuration
func (c *Config) Validate() error {
	if err := c.Typesense.Validate(); err != nil {
		return errors.WrapError(err, errors.ErrorTypeConfiguration, "invalid typesense configuration")
	}

	if err := c.Server.Validate(); err != nil {
		return errors.WrapError(err, errors.ErrorTypeConfiguration, "invalid server configuration")
	}

	if err := c.Logging.Validate(); err != nil {
		return errors.WrapError(err, errors.ErrorTypeConfiguration, "invalid logging configuration")
	}

	return nil
}

// Validate validates the Typesense configuration
func (t *TypesenseConfig) Validate() error {
	if t.Host == "" {
		return errors.NewInvalidConfigurationError("host", t.Host)
	}

	if t.Port < constants.MinPort || t.Port > constants.MaxPort {
		return errors.NewInvalidConfigurationError("port", t.Port)
	}

	if !isValidProtocol(t.Protocol) {
		return errors.NewInvalidConfigurationError("protocol", t.Protocol)
	}

	if t.APIKey == "" {
		return errors.NewInvalidConfigurationError("api_key", "empty")
	}

	if t.Timeout <= 0 {
		return errors.NewInvalidConfigurationError("timeout", t.Timeout)
	}

	return nil
}

// Validate validates the Server configuration
func (s *ServerConfig) Validate() error {
	if s.Name == "" {
		return errors.NewInvalidConfigurationError("name", s.Name)
	}

	if s.Version == "" {
		return errors.NewInvalidConfigurationError("version", s.Version)
	}

	return nil
}

// Validate validates the Logging configuration
func (l *LoggingConfig) Validate() error {
	if !isValidLogLevel(l.Level) {
		return errors.NewInvalidConfigurationError("level", l.Level)
	}

	if !isValidLogFormat(l.Format) {
		return errors.NewInvalidConfigurationError("format", l.Format)
	}

	return nil
}

// URL returns the full Typesense server URL
func (t *TypesenseConfig) URL() string {
	return fmt.Sprintf("%s://%s:%d", t.Protocol, t.Host, t.Port)
}

// isValidProtocol checks if the protocol is supported
func isValidProtocol(protocol string) bool {
	for _, p := range constants.SupportedProtocols {
		if strings.EqualFold(protocol, p) {
			return true
		}
	}
	return false
}

// isValidLogLevel checks if the log level is supported
func isValidLogLevel(level string) bool {
	for _, l := range constants.SupportedLogLevels {
		if strings.EqualFold(level, l) {
			return true
		}
	}
	return false
}

// isValidLogFormat checks if the log format is supported
func isValidLogFormat(format string) bool {
	for _, f := range constants.SupportedLogFormats {
		if strings.EqualFold(format, f) {
			return true
		}
	}
	return false
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault returns environment variable as int or default
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
