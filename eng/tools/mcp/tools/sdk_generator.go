// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/link"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
	"github.com/google/go-github/v62/github"
	"github.com/mark3labs/mcp-go/mcp"
)

// SDKGeneratorTool creates and returns the generate-go-sdk tool
func SDKGeneratorTool() mcp.Tool {
	return mcp.NewTool("generate-go-sdk",
		mcp.WithDescription("Generates Azure Go SDK from TypeSpec configuration. Accepts either a direct path to tspconfig.yaml or a GitHub PR link."),
		mcp.WithString("sdk_repo_path",
			mcp.Description("Local path to the Azure SDK for Go repository"),
			mcp.Required(),
		),
		mcp.WithString("spec_repo_path",
			mcp.Description("Local path to the Azure REST API Specs repository"),
			mcp.Required(),
		),
		mcp.WithString("tsp_config_path",
			mcp.Description("Direct path to tspconfig.yaml file (relative to spec repo root) OR GitHub PR link"),
			mcp.Required(),
		),
	)
}

// SDKGeneratorRequest represents the input for SDK generation
type SDKGeneratorRequest struct {
	SDKRepoPath   string `json:"sdk_repo_path"`
	SpecRepoPath  string `json:"spec_repo_path"`
	TspConfigPath string `json:"tsp_config_path"`
}

// SDKGeneratorResult represents the result of SDK generation
type SDKGeneratorResult struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	PackageName    string `json:"package_name,omitempty"`
	PackagePath    string `json:"package_path,omitempty"`
	Version        string `json:"version,omitempty"`
	ChangelogMD    string `json:"changelog_md,omitempty"`
	HasBreaking    bool   `json:"has_breaking_changes,omitempty"`
	GenerationType string `json:"generation_type,omitempty"`
}

// SDKGeneratorHandler handles the generate-go-sdk tool requests
func SDKGeneratorHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse arguments using MCP's built-in parsing functions
	req := SDKGeneratorRequest{
		SDKRepoPath:   mcp.ParseString(request, "sdk_repo_path", ""),
		SpecRepoPath:  mcp.ParseString(request, "spec_repo_path", ""),
		TspConfigPath: mcp.ParseString(request, "tsp_config_path", ""),
	}

	// Validate required parameters
	if req.SDKRepoPath == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Error: sdk_repo_path is required"),
			},
			IsError: true,
		}, nil
	}
	if req.SpecRepoPath == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Error: spec_repo_path is required"),
			},
			IsError: true,
		}, nil
	}
	if req.TspConfigPath == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Error: tsp_config_path is required"),
			},
			IsError: true,
		}, nil
	}

	// Validate paths
	if err := validatePaths(req.SDKRepoPath, req.SpecRepoPath); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Path validation error: %v", err)),
			},
			IsError: true,
		}, nil
	}

	// Resolve tspconfig path
	tspConfigPath, err := resolveTspConfigPath(ctx, req.TspConfigPath, req.SpecRepoPath)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error resolving TypeSpec config path: %v", err)),
			},
			IsError: true,
		}, nil
	}

	// Generate SDK
	result, err := generateSDK(ctx, req, tspConfigPath)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("SDK generation failed: %v", err)),
			},
			IsError: true,
		}, nil
	}

	// Format result as JSON
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Error formatting result: %v", err)),
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(resultJSON)),
		},
	}, nil
}

// validatePaths validates that the provided paths exist and are directories
func validatePaths(sdkPath, specPath string) error {
	if info, err := os.Stat(sdkPath); err != nil || !info.IsDir() {
		return fmt.Errorf("SDK repository path '%s' does not exist or is not a directory", sdkPath)
	}

	if info, err := os.Stat(specPath); err != nil || !info.IsDir() {
		return fmt.Errorf("spec repository path '%s' does not exist or is not a directory", specPath)
	}

	return nil
}

// resolveTspConfigPath resolves the TypeSpec config path from either a direct path or PR link
func resolveTspConfigPath(ctx context.Context, input, specRepoPath string) (string, error) {
	// Check if input looks like a GitHub PR link
	if isGitHubPRLink(input) {
		return extractTspConfigFromPR(ctx, input, specRepoPath)
	}

	// Treat as direct path - make it absolute relative to spec repo
	tspConfigPath := filepath.Join(specRepoPath, input)
	if !strings.HasSuffix(tspConfigPath, "tspconfig.yaml") {
		tspConfigPath = filepath.Join(tspConfigPath, "tspconfig.yaml")
	}

	// Verify the file exists
	if _, err := os.Stat(tspConfigPath); err != nil {
		return "", fmt.Errorf("TypeSpec config file not found at '%s'", tspConfigPath)
	}

	return tspConfigPath, nil
}

// isGitHubPRLink checks if the input string is a GitHub PR link
func isGitHubPRLink(input string) bool {
	prPattern := regexp.MustCompile(`^https://github\.com/[^/]+/[^/]+/pull/\d+`)
	return prPattern.MatchString(input)
}

// extractTspConfigFromPR extracts the tspconfig.yaml path from a GitHub PR
func extractTspConfigFromPR(ctx context.Context, prLink, specRepoPath string) (string, error) {
	// Extract PR number from the link
	prPattern := regexp.MustCompile(`https://github\.com/([^/]+)/([^/]+)/pull/(\d+)`)
	matches := prPattern.FindStringSubmatch(prLink)
	if len(matches) != 4 {
		return "", fmt.Errorf("invalid GitHub PR link format: %s", prLink)
	}

	owner := matches[1]
	repo := matches[2]
	prNumberStr := matches[3]

	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		return "", fmt.Errorf("invalid PR number in link: %s", prNumberStr)
	}

	// Create GitHub client to fetch changed files
	client := query.NewClient()

	// Get changed files from the PR using pagination
	opt := &github.ListOptions{
		PerPage: 100,
	}
	var allFiles []*github.CommitFile
	for {
		files, resp, err := client.PullRequests.ListFiles(ctx, owner, repo, prNumber, opt)
		if err != nil {
			return "", fmt.Errorf("failed to get changed files from PR %d: %v", prNumber, err)
		}
		allFiles = append(allFiles, files...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// Convert to file paths
	var filePaths []string
	for _, file := range allFiles {
		if file.GetFilename() != "" {
			filePaths = append(filePaths, file.GetFilename())
		}
	}

	// Use the existing logic to extract tspconfig path from changed files
	tspConfigRelativePath, err := link.GetTspConfigPathFromChangedFiles(ctx, client, filePaths)
	if err != nil {
		return "", fmt.Errorf("failed to extract tspconfig path from PR: %v", err)
	}

	// Convert to absolute path
	tspConfigPath := filepath.Join(specRepoPath, string(tspConfigRelativePath))

	// Verify the file exists locally
	if _, err := os.Stat(tspConfigPath); err != nil {
		return "", fmt.Errorf("TypeSpec config file not found locally at '%s'. Make sure you have pulled the latest changes from the spec repository", tspConfigPath)
	}

	return tspConfigPath, nil
}

// getSpecCommitHash gets the current Git commit hash from the spec repository
func getSpecCommitHash(specRepoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = specRepoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Git commit hash from spec repository: %v", err)
	}

	commitHash := strings.TrimSpace(string(output))
	if commitHash == "" {
		return "", fmt.Errorf("empty commit hash returned from spec repository")
	}

	return commitHash, nil
}

// generateSDK performs the actual SDK generation using the GenerateFromTypeSpec method
func generateSDK(ctx context.Context, req SDKGeneratorRequest, tspConfigPath string) (*SDKGeneratorResult, error) {
	// Create SDK repository reference
	sdkRepo, err := repo.OpenSDKRepository(req.SDKRepoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SDK repository: %v", err)
	}

	// Get the current commit hash from the spec repository
	specCommitHash, err := getSpecCommitHash(req.SpecRepoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get spec commit hash: %v", err)
	}

	// Create generation context
	generateCtx := common.GenerateContext{
		SDKPath:        utils.NormalizePath(req.SDKRepoPath),
		SDKRepo:        &sdkRepo,
		SpecPath:       utils.NormalizePath(req.SpecRepoPath),
		SpecCommitHash: specCommitHash,
		SpecRepoURL:    "Azure/azure-rest-api-specs",
	}

	// Create generation parameters
	generateParam := &common.GenerateParam{
		GoVersion:        MinGoVersion,
		TspClientOptions: []string{"--debug"},
		ApiVersion:       "",
		SdkReleaseType:   "",
	}

	// Perform the generation
	result, err := generateCtx.GenerateFromTypeSpec(tspConfigPath, generateParam)
	if err != nil {
		return &SDKGeneratorResult{
			Success: false,
			Message: fmt.Sprintf("Generation failed: %v", err),
		}, nil
	}

	if result == nil {
		return &SDKGeneratorResult{
			Success: false,
			Message: "Generation completed but no result was returned. This might indicate that the TypeSpec config doesn't have @azure-tools/typespec-go emitter configured.",
		}, nil
	}

	return &SDKGeneratorResult{
		Success:        true,
		Message:        "SDK generation completed successfully",
		PackageName:    result.PackageName,
		PackagePath:    result.PackageRelativePath,
		Version:        result.Version,
		ChangelogMD:    result.ChangelogMD,
		HasBreaking:    result.Changelog.HasBreakingChanges(),
		GenerationType: result.GenerationType,
	}, nil
}
