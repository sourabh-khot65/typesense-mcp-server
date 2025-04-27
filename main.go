package main

import (
	"fmt"
	"log"
	"os"
	"typesense-mcp-server/tools"

	"github.com/mark3labs/mcp-go/server"
)

func init() {
	// Set up logging to file
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("ERROR: Failed to open log file: %v", err)
	} else {
		log.SetOutput(logFile)
	}
}

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Typesense MCP Server",
		"1.0.0",
		// server.WithResourceCapabilities(true, true),
		// server.WithToolCapabilities(true),
	)

	// Register tools
	tools.RegisterTools(s)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
