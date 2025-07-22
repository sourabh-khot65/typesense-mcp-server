// Package main provides the main entry point for the Typesense MCP Server.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"typesense-mcp-server/internal/config"
	"typesense-mcp-server/internal/constants"
	"typesense-mcp-server/internal/health"
	"typesense-mcp-server/internal/logger"
	"typesense-mcp-server/internal/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Initialize logger first
	log := logger.NewDefault()
	
	entry := log.WithComponent(constants.ComponentMain)
	entry.WithField("version", constants.AppVersion).
		Info("starting typesense mcp server")

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		entry.WithError(err).Fatal("failed to load configuration")
	}

	entry.WithField("typesense_url", cfg.Typesense.URL()).
		Info("configuration loaded successfully")

	// Create tool registry
	registry, err := tools.New(cfg, log)
	if err != nil {
		entry.WithError(err).Fatal("failed to create tool registry")
	}

	// Ensure cleanup on exit
	defer func() {
		if err := registry.Close(); err != nil {
			entry.WithError(err).Error("failed to close registry")
		}
	}()

	// Create MCP server
	mcpServer := server.NewMCPServer(
		cfg.Server.Name,
		cfg.Server.Version,
	)

	// Register tools
	if err := registry.RegisterTools(mcpServer); err != nil {
		entry.WithError(err).Fatal("failed to register tools")
	}

	// Setup health checker
	healthChecker := health.New(registry.GetService(), cfg.Server.Version, log)
	
	// Perform initial health check
	healthResult := healthChecker.Check(context.Background())
	if !healthResult.IsHealthy() {
		entry.WithError(healthResult.HealthError()).
			Error("initial health check failed")
		// Continue anyway, but log the issue
	} else {
		entry.Info("initial health check passed")
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		entry.Info("shutdown signal received, initiating graceful shutdown")
		cancel()
	}()

	// Start the server
	entry.Info("starting MCP server")
	
	serverErr := make(chan error, 1)
	go func() {
		if err := server.ServeStdio(mcpServer); err != nil {
			serverErr <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case <-ctx.Done():
		entry.Info("shutting down gracefully")
		// Give the server a moment to clean up
		time.Sleep(100 * time.Millisecond)
	case err := <-serverErr:
		entry.WithError(err).Error("server error occurred")
		os.Exit(1)
	}

	entry.Info("server shutdown complete")
}
