// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generate

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

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
	"github.com/spf13/cobra"
)

// SDKGeneratorResult represents the result of SDK generation
type SDKGeneratorResult struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	PackageName    string `json:"package_name,omitempty"`
	PackagePath    string `json:"package_path,omitempty"`
	SpecFolderPath string `json:"spec_folder_path,omitempty"`
	Version        string `json:"version,omitempty"`
	ChangelogMD    string `json:"changelog_md,omitempty"`
	HasBreaking    bool   `json:"has_breaking_changes,omitempty"`
	GenerationType string `json:"generation_type,omitempty"`
}

// Command returns the generate command
func Command() *cobra.Command {
	var outputFormat string
	var tspConfigPath string
	var githubPRLink string
	var debug bool

	generateCmd := &cobra.Command{
		Use:   "generate <sdk-repo-path> <spec-repo-path>",
		Short: "Generate Azure Go SDK from TypeSpec configuration",
		Long: `Generate Azure Go SDK from TypeSpec configuration.

This command generates an Azure Go SDK package from TypeSpec specifications. 
You must provide either a direct path to tspconfig.yaml OR a GitHub PR link (exactly one is required).

Examples:
  # Generate from direct TypeSpec config path
  generator generate /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs --tsp-config specification/cognitiveservices/OpenAI.Inference/tspconfig.yaml

  # Generate from GitHub PR
  generator generate /path/to/azure-sdk-for-go /path/to/azure-rest-api-specs --github-pr https://github.com/Azure/azure-rest-api-specs/pull/12345

The command will:
1. Validate the provided repository paths
2. Resolve the TypeSpec configuration (checking out PR branch if needed)
3. Generate the Go SDK using the TypeSpec-Go emitter
4. Output generation results in the specified format`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sdkRepoPath := args[0]
			specRepoPath := args[1]

			// Validate that exactly one of tsp_config_path or github_pr_link is provided
			if tspConfigPath == "" && githubPRLink == "" {
				return fmt.Errorf("either --tsp-config or --github-pr must be provided")
			}
			if tspConfigPath != "" && githubPRLink != "" {
				return fmt.Errorf("only one of --tsp-config or --github-pr should be provided, not both")
			}

			// Validate paths
			if err := validatePaths(sdkRepoPath, specRepoPath); err != nil {
				return fmt.Errorf("path validation error: %v", err)
			}

			// Resolve tspconfig path
			var inputPath string
			if tspConfigPath != "" {
				inputPath = tspConfigPath
			} else {
				inputPath = githubPRLink
			}

			resolvedTspConfigPath, err := resolveTspConfigPath(context.Background(), inputPath, specRepoPath)
			if err != nil {
				return fmt.Errorf("error resolving TypeSpec config path: %v", err)
			}

			// Generate SDK
			result, err := generateSDK(context.Background(), sdkRepoPath, specRepoPath, resolvedTspConfigPath, debug)
			if err != nil {
				return fmt.Errorf("SDK generation failed: %v", err)
			}

			// Output result
			switch outputFormat {
			case "json":
				jsonResult, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal result: %v", err)
				}
				fmt.Println(string(jsonResult))
			default:
				// Human-readable output
				if result.Success {
					fmt.Printf("✓ SDK generation completed successfully!\n\n")
					fmt.Printf("Package Name: %s\n", result.PackageName)
					fmt.Printf("Package Path: %s\n", result.PackagePath)
					fmt.Printf("Spec Folder: %s\n", result.SpecFolderPath)
					fmt.Printf("Version: %s\n", result.Version)
					fmt.Printf("Generation Type: %s\n", result.GenerationType)
					if result.HasBreaking {
						fmt.Printf("⚠ Has Breaking Changes: Yes\n")
					} else {
						fmt.Printf("✓ Has Breaking Changes: No\n")
					}
					if result.ChangelogMD != "" {
						fmt.Printf("\nChangelog:\n%s\n", result.ChangelogMD)
					}
				} else {
					fmt.Printf("✗ SDK generation failed: %s\n", result.Message)
					return fmt.Errorf("generation failed")
				}
			}

			return nil
		},
	}

	generateCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	generateCmd.Flags().StringVar(&tspConfigPath, "tsp-config", "", "Direct path to tspconfig.yaml file (relative to spec repo root)")
	generateCmd.Flags().StringVar(&githubPRLink, "github-pr", "", "GitHub PR link to extract TypeSpec configuration from")
	generateCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	return generateCmd
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

// checkoutPRBranch checks out the PR branch in the spec repository using GitHub CLI
func checkoutPRBranch(ctx context.Context, specRepoPath string, prNumber int, owner, repo string) error {
	fmt.Printf("Checking out PR #%d using GitHub CLI\n", prNumber)

	// Use GitHub CLI to checkout the PR
	checkoutCmd := exec.Command("gh", "pr", "checkout", strconv.Itoa(prNumber))
	checkoutCmd.Dir = specRepoPath

	output, err := checkoutCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout PR #%d using GitHub CLI: %v\nOutput: %s", prNumber, err, string(output))
	}

	fmt.Printf("Successfully checked out PR #%d\n", prNumber)
	fmt.Printf("GitHub CLI output: %s\n", string(output))

	return nil
}

// extractTspConfigFromPR extracts the tspconfig.yaml path from a GitHub PR and checks out the PR branch using GitHub CLI
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

	// Checkout the PR branch using GitHub CLI
	if err := checkoutPRBranch(ctx, specRepoPath, prNumber, owner, repo); err != nil {
		return "", fmt.Errorf("failed to checkout PR branch: %v", err)
	}

	// Get changed files from the PR using GitHub CLI
	getFilesCmd := exec.Command("gh", "pr", "view", strconv.Itoa(prNumber), "--json", "files", "--jq", ".files[].path")
	getFilesCmd.Dir = specRepoPath

	output, err := getFilesCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get changed files from PR #%d using GitHub CLI: %v", prNumber, err)
	}

	// Parse the file paths from the output
	filePaths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(filePaths) == 1 && filePaths[0] == "" {
		return "", fmt.Errorf("no changed files found in PR #%d", prNumber)
	}

	// Use the existing logic to extract tspconfig path from changed files
	tspConfigRelativePath, err := extractTspConfigPathFromFiles(filePaths)
	if err != nil {
		return "", fmt.Errorf("failed to extract tspconfig path from PR: %v", err)
	}

	// Convert to absolute path
	tspConfigPath := filepath.Join(specRepoPath, tspConfigRelativePath)

	// Verify the file exists locally (should exist now that we've checked out the PR branch)
	if _, err := os.Stat(tspConfigPath); err != nil {
		return "", fmt.Errorf("TypeSpec config file not found locally at '%s' after checking out PR branch", tspConfigPath)
	}

	return tspConfigPath, nil
}

// extractTspConfigPathFromFiles finds the tspconfig.yaml path from a list of changed files
func extractTspConfigPathFromFiles(filePaths []string) (string, error) {
	// Look for tspconfig.yaml files in the changed files
	for _, filePath := range filePaths {
		if strings.HasSuffix(filePath, "tspconfig.yaml") {
			return filePath, nil
		}
	}

	// If no tspconfig.yaml found directly, look for TypeSpec files and try to find the containing service directory
	for _, filePath := range filePaths {
		if strings.Contains(filePath, "specification/") && (strings.HasSuffix(filePath, ".tsp") || strings.Contains(filePath, "/")) {
			// Extract the service directory path
			parts := strings.Split(filePath, "/")
			for i, part := range parts {
				if part == "specification" && i+2 < len(parts) {
					// Pattern: specification/{service}/{namespace}/...
					servicePath := strings.Join(parts[:i+3], "/")
					tspConfigPath := servicePath + "/tspconfig.yaml"
					return tspConfigPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no tspconfig.yaml file found in changed files and unable to infer from TypeSpec files")
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
func generateSDK(ctx context.Context, sdkRepoPath, specRepoPath, tspConfigPath string, debug bool) (*SDKGeneratorResult, error) {
	// Create SDK repository reference
	sdkRepo, err := repo.OpenSDKRepository(sdkRepoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SDK repository: %v", err)
	}

	// Get the current commit hash from the spec repository
	specCommitHash, err := getSpecCommitHash(specRepoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get spec commit hash: %v", err)
	}

	// Create generation context
	generateCtx := common.GenerateContext{
		SDKPath:        utils.NormalizePath(sdkRepoPath),
		SDKRepo:        &sdkRepo,
		SpecPath:       utils.NormalizePath(specRepoPath),
		SpecCommitHash: specCommitHash,
		SpecRepoURL:    "Azure/azure-rest-api-specs",
	}

	// Create generation parameters
	tspClientOptions := []string{}
	if debug {
		tspClientOptions = append(tspClientOptions, "--debug")
	}

	generateParam := &common.GenerateParam{
		TspClientOptions: tspClientOptions,
		SkipUpdateDep:    true,
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

	specFolderPath := filepath.Dir(tspConfigPath)
	packagePath := filepath.Join(sdkRepoPath, result.PackageRelativePath)

	return &SDKGeneratorResult{
		Success:        true,
		Message:        "SDK generation completed successfully",
		PackageName:    result.PackageName,
		PackagePath:    packagePath,
		SpecFolderPath: specFolderPath,
		Version:        result.Version,
		ChangelogMD:    result.ChangelogMD,
		HasBreaking:    result.Changelog.HasBreakingChanges(),
		GenerationType: result.GenerationType,
	}, nil
}
