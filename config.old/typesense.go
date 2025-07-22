package config

import (
	"fmt"
	"os"
	"strconv"
)

type TypesenseConfig struct {
	Host     string
	Port     int
	Protocol string
	APIKey   string
}

func NewTypesenseConfig() *TypesenseConfig {
	config := &TypesenseConfig{
		Host:     getEnvOrDefault("TYPESENSE_HOST", "localhost"),
		Port:     getEnvIntOrDefault("TYPESENSE_PORT", 8108),
		Protocol: getEnvOrDefault("TYPESENSE_PROTOCOL", "http"),
		APIKey:   getEnvOrDefault("TYPESENSE_API_KEY", "xyz"),
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("invalid Typesense configuration: %v", err))
	}

	return config
}

// Validate validates the Typesense configuration
func (c *TypesenseConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}
	
	if c.Protocol != "http" && c.Protocol != "https" {
		return fmt.Errorf("protocol must be 'http' or 'https', got '%s'", c.Protocol)
	}
	
	if c.APIKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	
	return nil
}

func (c *TypesenseConfig) URL() string {
	return fmt.Sprintf("%s://%s:%d", c.Protocol, c.Host, c.Port)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
