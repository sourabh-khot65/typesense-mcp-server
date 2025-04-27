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
	return &TypesenseConfig{
		Host:     getEnvOrDefault("TYPESENSE_HOST", "localhost"),
		Port:     getEnvIntOrDefault("TYPESENSE_PORT", 8108),
		Protocol: getEnvOrDefault("TYPESENSE_PROTOCOL", "http"),
		APIKey:   getEnvOrDefault("TYPESENSE_API_KEY", "xyz"),
	}
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
