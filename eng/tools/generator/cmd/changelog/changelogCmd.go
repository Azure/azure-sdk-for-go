// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package changelog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/version"
	"github.com/spf13/cobra"
)

// ChangelogResult represents the result of a changelog operation
type ChangelogResult struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	PackagePath   string `json:"package_path,omitempty"`
	PackageStatus string `json:"package_status,omitempty"`
	Version       string `json:"version,omitempty"`
	ReleaseDate   string `json:"release_date,omitempty"`
	ChangelogMD   string `json:"changelog_md,omitempty"`
}

// Command returns the changelog command
func Command() *cobra.Command {
	var outputFormat string
	var verbose bool

	changelogCmd := &cobra.Command{
		Use:   "changelog <package-path>",
		Short: "Update changelog content for a package",
		Long: `Update changelog content for a package by determining the package status and generating appropriate changelog entries.

This command will:
1. Determine the package status (new package, existing package with new preview version, existing package with new stable version)
2. For new packages: generate changelog according to the template
3. For existing packages: compare current package exports with previous released version and calculate the changelog
4. Put the changelog entry into CHANGELOG.md file, replacing existing version entry if it exists

The package path should be an absolute path to a Go module (containing go.mod file).

Examples:
  # Update changelog for an existing package
  generator changelog /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute

  # Update changelog with verbose output
  generator changelog /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --verbose

  # Generate changelog with JSON output
  generator changelog /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --output json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packagePath := args[0]

			// Validate package path
			if err := validatePackagePath(packagePath); err != nil {
				return fmt.Errorf("package path validation error: %v", err)
			}

			// Get SDK root path
			sdkRoot, err := utils.FindSDKRoot(packagePath)
			if err != nil {
				return fmt.Errorf("failed to find SDK root: %v", err)
			}

			// Process changelog
			result, err := processChangelog(sdkRoot, packagePath, verbose)
			if err != nil {
				return fmt.Errorf("changelog operation failed: %v", err)
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
					fmt.Printf("✓ Changelog updated successfully!\n\n")
					fmt.Printf("Package: %s\n", result.PackagePath)
					fmt.Printf("Status: %s\n", result.PackageStatus)
					if result.Version != "" {
						fmt.Printf("Version: %s\n", result.Version)
					}
					fmt.Printf("Release Date: %s\n", result.ReleaseDate)
					if verbose && result.ChangelogMD != "" {
						fmt.Printf("\nGenerated Changelog:\n%s\n", result.ChangelogMD)
					}
				} else {
					fmt.Printf("✗ Changelog update failed: %s\n", result.Message)
					return fmt.Errorf("changelog update failed")
				}
			}

			return nil
		},
	}

	changelogCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	changelogCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	return changelogCmd
}

// validatePackagePath validates that the provided package path exists and contains a go.mod file
func validatePackagePath(packagePath string) error {
	if info, err := os.Stat(packagePath); err != nil || !info.IsDir() {
		return fmt.Errorf("package path '%s' does not exist or is not a directory", packagePath)
	}

	// Check if directory contains go.mod file
	goModPath := filepath.Join(packagePath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("package path '%s' does not contain a go.mod file", packagePath)
	} else if err != nil {
		return fmt.Errorf("failed to check for go.mod file: %v", err)
	}

	return nil
}

// processChangelog processes the changelog for the given package
func processChangelog(sdkRoot, packagePath string, verbose bool) (*ChangelogResult, error) {
	result := &ChangelogResult{
		PackagePath: packagePath,
		ReleaseDate: time.Now().Format("2006-01-02"),
	}

	// Initialize SDK repo
	sdkRepo, err := common.GetSDKRepo(sdkRoot, "")
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Failed to initialize SDK repository: %v", err)
		return result, nil
	}

	if verbose {
		log.Printf("Processing package at: %s", packagePath)
	}

	// Determine package status
	status, err := changelog.DetermineModuleStatus(packagePath, sdkRepo)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Failed to determine package status: %v", err)
		return result, nil
	}

	switch status {
	case utils.PackageStatusNew:
		result.PackageStatus = "New Package"
		err = handleNewPackage(packagePath, sdkRepo, result, verbose)
	case utils.PackageStatusExisting:
		result.PackageStatus = "Existing Package"
		err = handleExistingPackage(packagePath, sdkRepo, result, verbose)
	default:
		result.Success = false
		result.Message = fmt.Sprintf("Unknown package status: %v", status)
		return result, nil
	}

	if err != nil {
		result.Success = false
		result.Message = err.Error()
		return result, nil
	}

	result.Success = true
	result.Message = "Changelog updated successfully"
	return result, nil
}

func handleNewPackage(modulePath string, sdkRepo repo.SDKRepository, result *ChangelogResult, verbose bool) error {
	if verbose {
		log.Printf("Handling new package...")
	}

	result.Version = "0.1.0"

	err := changelog.CreateNewChangelog(modulePath, sdkRepo, result.Version, result.ReleaseDate)
	if err != nil {
		return fmt.Errorf("failed to create new package changelog: %v", err)
	}

	result.ChangelogMD = "New package"

	if verbose {
		log.Printf("Created new package changelog for %s", modulePath)
	}

	return nil
}

func handleExistingPackage(modulePath string, sdkRepo repo.SDKRepository, result *ChangelogResult, verbose bool) error {
	if verbose {
		log.Printf("Handling existing package...")
	}

	// Generate changelog using the new GenerateChangelog function
	if verbose {
		log.Printf("Generating changelog...")
	}

	// Generate changelog
	isCurrentPreview, err := version.IsCurrentPreviewVersion(modulePath, sdkRepo, nil)
	if err != nil {
		return fmt.Errorf("failed to determine if current version is preview: %v", err)
	}
	changelogResult, err := changelog.GenerateChangelog(modulePath, sdkRepo, isCurrentPreview)
	if err != nil {
		return fmt.Errorf("failed to generate changelog: %v", err)
	}

	// Calculate new version
	newVersion, _, err := version.CalculateNewVersion(changelogResult.ChangelogData, changelogResult.PreviousVersion, isCurrentPreview)
	if err != nil {
		return fmt.Errorf("failed to calculate new version: %v", err)
	}

	if verbose {
		log.Printf("Calculated new version: %s", newVersion.String())
	}

	result.Version = newVersion.String()

	// Update changelog file
	changelogMd, err := changelog.AddChangelogToFileWithReplacement(changelogResult.ChangelogData, newVersion, modulePath, result.ReleaseDate)
	if err != nil {
		return fmt.Errorf("failed to update changelog file: %v", err)
	}

	result.ChangelogMD = changelogMd

	if verbose {
		log.Printf("Successfully updated changelog for version %s", newVersion.String())
	}

	return nil
}
