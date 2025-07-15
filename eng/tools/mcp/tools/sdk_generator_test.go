// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package tools

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestValidatePaths(t *testing.T) {
	// Create temporary directories for testing
	tempDir := t.TempDir()
	sdkPath := filepath.Join(tempDir, "sdk")
	specPath := filepath.Join(tempDir, "spec")

	// Test with non-existent paths
	err := validatePaths("/non/existent/path", "/another/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent paths")
	}

	// Create the directories
	if err := os.MkdirAll(sdkPath, 0755); err != nil {
		t.Fatalf("Failed to create SDK directory: %v", err)
	}
	if err := os.MkdirAll(specPath, 0755); err != nil {
		t.Fatalf("Failed to create spec directory: %v", err)
	}

	// Test with valid paths
	err = validatePaths(sdkPath, specPath)
	if err != nil {
		t.Errorf("Expected no error for valid paths, got: %v", err)
	}

	// Test with file instead of directory
	filePath := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = validatePaths(filePath, specPath)
	if err == nil {
		t.Error("Expected error when SDK path is a file, not directory")
	}
}

func TestIsGitHubPRLink(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"https://github.com/Azure/azure-rest-api-specs/pull/12345", true},
		{"https://github.com/owner/repo/pull/123", true},
		{"https://github.com/Azure/azure-rest-api-specs/pull/12345/files", true},
		{"https://example.com/pull/123", false},
		{"github.com/Azure/azure-rest-api-specs/pull/123", false},
		{"https://github.com/Azure/azure-rest-api-specs/issues/123", false},
		{"not-a-url", false},
		{"", false},
	}

	for _, test := range tests {
		result := isGitHubPRLink(test.input)
		if result != test.expected {
			t.Errorf("isGitHubPRLink(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestResolveTspConfigPath(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	specPath := filepath.Join(tempDir, "spec")
	tspDir := filepath.Join(specPath, "specification", "contosowidgetmanager", "Contoso.Management")
	tspConfigPath := filepath.Join(tspDir, "tspconfig.yaml")

	// Create the directory structure
	if err := os.MkdirAll(tspDir, 0755); err != nil {
		t.Fatalf("Failed to create tsp directory: %v", err)
	}

	// Create a test tspconfig.yaml file
	tspContent := `
emit:
  - "@azure-tools/typespec-go"
options:
  "@azure-tools/typespec-go":
    package-dir: "sdk/resourcemanager/contosowidgetmanager/armcontosowidgetmanager"
`
	if err := os.WriteFile(tspConfigPath, []byte(tspContent), 0644); err != nil {
		t.Fatalf("Failed to create tspconfig.yaml: %v", err)
	}

	ctx := context.Background()

	// Test with direct path to directory (should append tspconfig.yaml)
	result, err := resolveTspConfigPath(ctx, "specification/contosowidgetmanager/Contoso.Management", specPath)
	if err != nil {
		t.Errorf("Expected no error for valid directory path, got: %v", err)
	}
	if result != tspConfigPath {
		t.Errorf("Expected resolved path '%s', got '%s'", tspConfigPath, result)
	}

	// Test with direct path to file
	result, err = resolveTspConfigPath(ctx, "specification/contosowidgetmanager/Contoso.Management/tspconfig.yaml", specPath)
	if err != nil {
		t.Errorf("Expected no error for valid file path, got: %v", err)
	}
	if result != tspConfigPath {
		t.Errorf("Expected resolved path '%s', got '%s'", tspConfigPath, result)
	}

	// Test with non-existent path
	_, err = resolveTspConfigPath(ctx, "specification/nonexistent/service", specPath)
	if err == nil {
		t.Error("Expected error for non-existent path")
	}

	// Test with GitHub PR link
	prLink := "https://github.com/Azure/azure-rest-api-specs/pull/30040"
	result, err = resolveTspConfigPath(ctx, prLink, specPath)
	if err != nil {
		t.Errorf("Expected no error for valid PR link, got: %v", err)
	}
	if result != tspConfigPath {
		t.Errorf("Expected resolved path '%s', got '%s'", tspConfigPath, result)
	}
}

func TestSDKGeneratorHandler_ArgumentParsing(t *testing.T) {
	ctx := context.Background()

	// Test with missing required parameters
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"sdk_repo_path": "",
			},
		},
	}

	result, err := SDKGeneratorHandler(ctx, request)
	if err != nil {
		t.Errorf("Expected no error from handler, got: %v", err)
	}

	if !result.IsError {
		t.Error("Expected error result for missing required parameters")
	}

	// Test with valid parameters but invalid paths
	request = mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"sdk_repo_path":   "/non/existent/sdk/path",
				"spec_repo_path":  "/non/existent/spec/path",
				"tsp_config_path": "some/config/path",
			},
		},
	}

	result, err = SDKGeneratorHandler(ctx, request)
	if err != nil {
		t.Errorf("Expected no error from handler, got: %v", err)
	}

	if !result.IsError {
		t.Error("Expected error result for invalid paths")
	}
}

func TestSDKGeneratorHandler_ValidPaths(t *testing.T) {
	ctx := context.Background()

	// Create temporary directory structure
	tempDir := t.TempDir()
	sdkPath := filepath.Join(tempDir, "sdk")
	specPath := filepath.Join(tempDir, "spec")
	tspDir := filepath.Join(specPath, "specification", "test", "service")
	tspConfigPath := filepath.Join(tspDir, "tspconfig.yaml")

	// Create the directory structure
	if err := os.MkdirAll(sdkPath, 0755); err != nil {
		t.Fatalf("Failed to create SDK directory: %v", err)
	}
	if err := os.MkdirAll(tspDir, 0755); err != nil {
		t.Fatalf("Failed to create tsp directory: %v", err)
	}

	// Create a minimal tspconfig.yaml file
	tspContent := `
emit:
  - "@azure-tools/typespec-go"
options:
  "@azure-tools/typespec-go":
    package-dir: "sdk/resourcemanager/test/armtest"
`
	if err := os.WriteFile(tspConfigPath, []byte(tspContent), 0644); err != nil {
		t.Fatalf("Failed to create tspconfig.yaml: %v", err)
	}

	// Test with valid parameters (this will likely fail at SDK generation due to missing dependencies)
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"sdk_repo_path":   sdkPath,
				"spec_repo_path":  specPath,
				"tsp_config_path": "specification/test/service",
			},
		},
	}

	result, err := SDKGeneratorHandler(ctx, request)
	if err != nil {
		t.Errorf("Expected no error from handler, got: %v", err)
	}

	// The result should contain valid JSON even if generation fails
	if len(result.Content) == 0 {
		t.Error("Expected result content")
	}

	// Try to parse the result as JSON
	var generatorResult SDKGeneratorResult
	if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
		if err := json.Unmarshal([]byte(textContent.Text), &generatorResult); err != nil {
			t.Errorf("Expected valid JSON result, got unmarshal error: %v", err)
		}
	}
}

func TestGetSpecCommitHash(t *testing.T) {
	// Create a temporary directory and initialize it as a Git repository
	tempDir := t.TempDir()

	// Initialize Git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Skip("Git not available, skipping test")
	}

	// Configure Git user (required for commits)
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user name: %v", err)
	}

	// Create a test file and commit it
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file to git: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Test getting the commit hash
	commitHash, err := getSpecCommitHash(tempDir)
	if err != nil {
		t.Errorf("Expected no error getting commit hash, got: %v", err)
	}

	if len(commitHash) != 40 {
		t.Errorf("Expected commit hash to be 40 characters long, got %d characters: %s", len(commitHash), commitHash)
	}

	// Test with non-git directory
	nonGitDir := t.TempDir()
	_, err = getSpecCommitHash(nonGitDir)
	if err == nil {
		t.Error("Expected error for non-git directory")
	}
}

func TestSDKGeneratorResult_JSON(t *testing.T) {
	// Test JSON marshaling/unmarshaling of result struct
	result := &SDKGeneratorResult{
		Success:     true,
		Message:     "Test message",
		PackageName: "armtest",
		PackagePath: "sdk/resourcemanager/test/armtest",
		Version:     "v1.0.0",
		ChangelogMD: "## Features\n- Initial release",
		HasBreaking: false,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Errorf("Failed to marshal result to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaled SDKGeneratorResult
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal result from JSON: %v", err)
	}

	// Verify fields
	if unmarshaled.Success != result.Success {
		t.Errorf("Success field mismatch: expected %v, got %v", result.Success, unmarshaled.Success)
	}
	if unmarshaled.Message != result.Message {
		t.Errorf("Message field mismatch: expected %q, got %q", result.Message, unmarshaled.Message)
	}
	if unmarshaled.PackageName != result.PackageName {
		t.Errorf("PackageName field mismatch: expected %q, got %q", result.PackageName, unmarshaled.PackageName)
	}
}
