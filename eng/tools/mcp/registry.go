// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"github.com/Azure/azure-sdk-for-go/eng/tools/mcp/tools"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllTools registers all available tools with the MCP server
func RegisterAllTools(s *server.MCPServer) {
	// Register environment checker tool
	s.AddTool(tools.EnvironmentCheckerTool(), tools.EnvironmentCheckerHandler)

	// Register SDK generator tool
	s.AddTool(tools.SDKGeneratorTool(), tools.SDKGeneratorHandler)
}
