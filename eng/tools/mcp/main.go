// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server instance
	s := server.NewMCPServer(
		"azure-sdk-go-mcp",
		version,
		server.WithToolCapabilities(false),
	)

	// Register all tools
	RegisterAllTools(s)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
