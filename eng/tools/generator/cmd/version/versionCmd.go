// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package version

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/version"
	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
)

// VersionResult represents the result of a version update operation
type VersionResult struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	PackagePath     string `json:"package_path,omitempty"`
	PreviousVersion string `json:"previous_version,omitempty"`
	NewVersion      string `json:"new_version,omitempty"`
}

// Command returns the version command
func Command() *cobra.Command {
	var outputFormat string
	var verbose bool
	var sdkReleaseType string
	var sdkVersion string

	versionCmd := &cobra.Command{
		Use:   "version <package-path>",
		Short: "Update package version files",
		Long: `Update package version files by calculating new version or using a specified version.

This command will:
1. If sdkversion is specified: update all version files with the provided version
2. If sdkversion is not specified: calculate new version based on package changes and sdkreleasetype, then update all version files

The package path should be an absolute path to a Go module (containing both go.mod and version.go files).

Examples:
  # Update version files with a specific version
  generator version /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --sdkversion 1.2.0

  # Calculate and update version based on changes
  generator version /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute

  # Calculate and update version as stable release
  generator version /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --sdkreleasetype stable

  # Calculate and update version with JSON output
  generator version /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --output json`,
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

			// Process version update
			result, err := processVersionUpdate(sdkRoot, packagePath, sdkVersion, sdkReleaseType, verbose)
			if err != nil {
				return fmt.Errorf("version update failed: %v", err)
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
					fmt.Printf("✓ Version updated successfully!\n\n")
					fmt.Printf("Package: %s\n", result.PackagePath)
					if result.PreviousVersion != "" {
						fmt.Printf("Previous Version: %s\n", result.PreviousVersion)
					}
					fmt.Printf("New Version: %s\n", result.NewVersion)
				} else {
					fmt.Printf("✗ Version update failed: %s\n", result.Message)
					return fmt.Errorf("version update failed")
				}
			}

			return nil
		},
	}

	versionCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	versionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	versionCmd.Flags().StringVar(&sdkReleaseType, "sdkreleasetype", "", "SDK release type (beta|stable), only used when sdkversion is not specified")
	versionCmd.Flags().StringVar(&sdkVersion, "sdkversion", "", "Specific SDK version to set (e.g., 1.2.0 or 1.2.0-beta.1)")

	return versionCmd
}

// validatePackagePath validates that the provided package path exists and contains necessary files
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

	// Check if directory contains version.go file
	versionGoPath := filepath.Join(packagePath, "version.go")
	if _, err := os.Stat(versionGoPath); os.IsNotExist(err) {
		return fmt.Errorf("package path '%s' does not contain a version.go file", packagePath)
	} else if err != nil {
		return fmt.Errorf("failed to check for version.go file: %v", err)
	}

	return nil
}

// processVersionUpdate processes the version update for the given package
func processVersionUpdate(sdkRoot, packagePath, sdkVersion, sdkReleaseType string, verbose bool) (*VersionResult, error) {
	result := &VersionResult{
		PackagePath: packagePath,
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

	var newVersion *semver.Version

	if sdkVersion != "" {
		// Use specified version
		if verbose {
			log.Printf("Using specified version: %s", sdkVersion)
		}

		// Parse the provided version
		if newVersion, err = semver.NewVersion(sdkVersion); err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Invalid version format '%s': %v", sdkVersion, err)
			return result, nil
		}
		result.NewVersion = newVersion.String()
	} else {
		// Calculate new version based on changes
		if verbose {
			log.Printf("Calculating new version with sdkreleasetype: %s", sdkReleaseType)
		}

		override := new(bool)
		if sdkReleaseType == "stable" {
			*override = false
		} else if sdkReleaseType == "beta" {
			*override = true
		}
		// Determine if this should be a preview version
		isCurrentPreview, err := version.IsCurrentPreviewVersion(packagePath, sdkRepo, override)
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Failed to determine if preview: %v", err)
			return result, nil
		}

		if verbose {
			log.Printf("Is current preview: %v", isCurrentPreview)
		}

		// Calculate new version
		changelogResult, err := changelog.GenerateChangelog(packagePath, sdkRepo, isCurrentPreview)
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Failed to generate changelog: %v", err)
			return result, nil
		}

		newVersion, _, err = version.CalculateNewVersion(changelogResult.ChangelogData, changelogResult.PreviousVersion, isCurrentPreview)
		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Failed to calculate new version: %v", err)
			return result, nil
		}

		result.NewVersion = newVersion.String()
		if verbose {
			log.Printf("Calculated new version: %s", newVersion.String())
		}
	}

	// Update all version files
	if verbose {
		log.Printf("Updating all version files...")
	}

	if err = version.UpdateAllVersionFiles(packagePath, newVersion, sdkRepo); err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Failed to update version files: %v", err)
		return result, nil
	}

	if verbose {
		log.Printf("Successfully updated all version files")
	}

	result.Success = true
	result.Message = "Version updated successfully"
	return result, nil
}
