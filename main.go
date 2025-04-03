package main

import (
	"fmt"
	"tb-mcp-server/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Tacitbase MCP Server",
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
